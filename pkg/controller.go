// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pkg

import (
	"errors"
	//"fmt"
	"reflect"
	//"strconv"
	"time"

	//"k8s.io/api/core/v1"
	//"istio.io/istio/pilot/pkg/serviceregistry/kube"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/log"
	"k8s.io/api/apps/v1"
)


// ControllerOptions stores the configurable attributes of a Controller.
type ControllerOptions struct {
	// Namespace the controller watches. If set to meta_v1.NamespaceAll (""), controller watches all namespaces
	WatchedNamespace string
	ResyncPeriod     time.Duration
	DomainSuffix     string
}

// Controller is a collection of synchronized resource watchers
// Caches are thread-safe
type Controller struct {
	domainSuffix string

	client    kubernetes.Interface
	queue     Queue
	services  cacheHandler
	endpoints cacheHandler
	nodes     cacheHandler
	deploy     cacheHandler

	pods *PodCache

	// Env is set by server to point to the environment, to allow the controller to
	// use env data and push status. It may be null in tests.
	Env *model.Environment
}

type cacheHandler struct {
	informer cache.SharedIndexInformer
	handler  *ChainHandler
}

// NewController creates a new Kubernetes controller
func NewController(client kubernetes.Interface, options ControllerOptions) *Controller {
	log.Infof("Service controller watching namespace %q for Deployments, refresh %s",
		options.WatchedNamespace, options.ResyncPeriod)

	// Queue requires a time duration for a retry delay after a handler error
	out := &Controller{
		domainSuffix: options.DomainSuffix,
		client:       client,
		queue:        NewQueue(1 * time.Second),
	}

	sharedInformers := informers.NewSharedInformerFactoryWithOptions(client, options.ResyncPeriod,
		informers.WithNamespace(options.WatchedNamespace))
	deployInformer := sharedInformers.Apps().V1().Deployments().Informer()

	out.deploy = out.createCacheHandler(deployInformer, "Deployments")

	return out
}

// notify is the first handler in the handler chain.
// Returning an error causes repeated execution of the entire chain.
func (c *Controller) notify(obj interface{}, event model.Event) error {
	if !c.HasSynced() {
		return errors.New("waiting till full synchronization")
	}
	return nil
}

// createCacheHandler registers handlers for a specific event.
// Current implementation queues the events in queue.go, and the handler is run with
// some throttling.
// Used for Service, Endpoint, Node and Pod.
// See config/kube for CRD events.
// See config/ingress for Ingress objects
func (c *Controller) createCacheHandler(informer cache.SharedIndexInformer, otype string) cacheHandler {
	handler := &ChainHandler{[]Handler{c.notify}}

	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			// TODO: filtering functions to skip over un-referenced resources (perf)
			AddFunc: func(obj interface{}) {
				deploy := obj.(*v1.Deployment)
				SetDeployIndex(deploy.Namespace+"/"+deploy.Name, deploy.Namespace)
				c.queue.Push(Task{handler.Apply, obj, model.EventAdd})
			},
			UpdateFunc: func(old, cur interface{}) {
				if !reflect.DeepEqual(old, cur) {
					c.queue.Push(Task{handler.Apply, cur, model.EventUpdate})
				}
			},
			DeleteFunc: func(obj interface{}) {
				deploy := obj.(*v1.Deployment)
				DelDeployIndex(deploy.Namespace+"/"+deploy.Name, deploy.Namespace)
				c.queue.Push(Task{handler.Apply, obj, model.EventDelete})
			},
		})

	return cacheHandler{informer: informer, handler: handler}
}

// HasSynced returns true after the initial state synchronization
func (c *Controller) HasSynced() bool {
	if !c.deploy.informer.HasSynced() {
		return false
	}
	return true
}

// Run all controllers until a signal is received
func (c *Controller) Run(stop <-chan struct{}) {
	go c.deploy.informer.Run(stop)

	<-stop
	log.Infof("Controller terminated")
}

func (c *Controller) GetDeployList(deployIndexs []string) []interface{} {
	list := make([]interface{}, len(deployIndexs))
	for k, index := range deployIndexs {
		item, exists, err := c.deploy.informer.GetIndexer().GetByKey(index)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		if !exists {
			log.Error(index + "not exists")
			continue
		}
		list[k] = item
	}

	return list
}

func (c *Controller) ListKeys() []string {
	list := c.deploy.informer.GetStore().ListKeys()
	return list
}

func (c *Controller) GetByKey(key string) (item interface{}, exists bool, err error) {
	return c.deploy.informer.GetStore().GetByKey(key)
}