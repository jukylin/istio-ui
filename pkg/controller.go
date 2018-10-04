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
				c.queue.Push(Task{handler.Apply, obj, model.EventAdd})
			},
			UpdateFunc: func(old, cur interface{}) {
				if !reflect.DeepEqual(old, cur) {
					c.queue.Push(Task{handler.Apply, cur, model.EventUpdate})
				}
			},
			DeleteFunc: func(obj interface{}) {
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

func (c *Controller) GetDeployList() []interface{} {
	list := c.deploy.informer.GetIndexer().List()
	return list
}

func (c *Controller) ListKeys() []string {
	list := c.deploy.informer.GetStore().ListKeys()
	return list
}

func (c *Controller) GetByKey(key string) (item interface{}, exists bool, err error) {
	return c.deploy.informer.GetStore().GetByKey(key)
}



//// WorkloadHealthCheckInfo implements a service catalog operation
//func (c *Controller) WorkloadHealthCheckInfo(addr string) model.ProbeList {
//	pod, exists := c.pods.getPodByIP(addr)
//	if !exists {
//		return nil
//	}
//
//	probes := make([]*model.Probe, 0)
//
//	// Obtain probes from the readiness and liveness probes
//	for _, container := range pod.Spec.Containers {
//		if container.ReadinessProbe != nil && container.ReadinessProbe.Handler.HTTPGet != nil {
//			p, err := convertProbePort(container, &container.ReadinessProbe.Handler)
//			if err != nil {
//				log.Infof("Error while parsing readiness probe port =%v", err)
//			}
//			probes = append(probes, &model.Probe{
//				Port: p,
//				Path: container.ReadinessProbe.Handler.HTTPGet.Path,
//			})
//		}
//		if container.LivenessProbe != nil && container.LivenessProbe.Handler.HTTPGet != nil {
//			p, err := convertProbePort(container, &container.LivenessProbe.Handler)
//			if err != nil {
//				log.Infof("Error while parsing liveness probe port =%v", err)
//			}
//			probes = append(probes, &model.Probe{
//				Port: p,
//				Path: container.LivenessProbe.Handler.HTTPGet.Path,
//			})
//		}
//	}
//
//	// Obtain probe from prometheus scrape
//	if scrape := pod.Annotations[PrometheusScrape]; scrape == "true" {
//		var port *model.Port
//		path := PrometheusPathDefault
//		if portstr := pod.Annotations[PrometheusPort]; portstr != "" {
//			portnum, err := strconv.Atoi(portstr)
//			if err != nil {
//				log.Warna(err)
//			} else {
//				port = &model.Port{
//					Port: portnum,
//				}
//			}
//		}
//		if pod.Annotations[PrometheusPath] != "" {
//			path = pod.Annotations[PrometheusPath]
//		}
//		probes = append(probes, &model.Probe{
//			Port: port,
//			Path: path,
//		})
//	}
//
//	return probes
//}
//
//// InstancesByPort implements a service catalog operation
//func (c *Controller) InstancesByPort(hostname model.Hostname, reqSvcPort int,
//	labelsList model.LabelsCollection) ([]*model.ServiceInstance, error) {
//	// Get actual service by name
//	name, namespace, err := parseHostname(hostname)
//	if err != nil {
//		log.Infof("parseHostname(%s) => error %v", hostname, err)
//		return nil, err
//	}
//
//	item, exists := c.serviceByKey(name, namespace)
//	if !exists {
//		return nil, nil
//	}
//
//	// Locate all ports in the actual service
//
//	svc := convertService(*item, c.domainSuffix)
//	if svc == nil {
//		return nil, nil
//	}
//
//	svcPortEntry, exists := svc.Ports.GetByPort(reqSvcPort)
//	if !exists && reqSvcPort != 0 {
//		return nil, nil
//	}
//
//	for _, item := range c.endpoints.informer.GetStore().List() {
//		ep := *item.(*v1.Endpoints)
//		if ep.Name == name && ep.Namespace == namespace {
//			var out []*model.ServiceInstance
//			for _, ss := range ep.Subsets {
//				for _, ea := range ss.Addresses {
//					labels, _ := c.pods.labelsByIP(ea.IP)
//					// check that one of the input labels is a subset of the labels
//					if !labelsList.HasSubsetOf(labels) {
//						continue
//					}
//
//					pod, exists := c.pods.getPodByIP(ea.IP)
//					az, sa, uid := "", "", ""
//					if exists {
//						az, _ = c.GetPodAZ(pod)
//						sa = kubeToIstioServiceAccount(pod.Spec.ServiceAccountName, pod.GetNamespace(), c.domainSuffix)
//						uid = fmt.Sprintf("kubernetes://%s.%s", pod.Name, pod.Namespace)
//					}
//
//					// identify the port by name. K8S EndpointPort uses the service port name
//					for _, port := range ss.Ports {
//						if port.Name == "" || // 'name optional if single port is defined'
//							reqSvcPort == 0 || // return all ports (mostly used by tests/debug)
//							svcPortEntry.Name == port.Name {
//							out = append(out, &model.ServiceInstance{
//								Endpoint: model.NetworkEndpoint{
//									Address:     ea.IP,
//									Port:        int(port.Port),
//									ServicePort: svcPortEntry,
//									UID:         uid,
//								},
//								Service:          svc,
//								Labels:           labels,
//								AvailabilityZone: az,
//								ServiceAccount:   sa,
//							})
//						}
//					}
//				}
//			}
//			return out, nil
//		}
//	}
//	return nil, nil
//}
//
//// GetProxyServiceInstances returns service instances co-located with a given proxy
//func (c *Controller) GetProxyServiceInstances(proxy *model.Proxy) ([]*model.ServiceInstance, error) {
//	var out []*model.ServiceInstance
//	proxyIP := proxy.IPAddress
//	for _, item := range c.endpoints.informer.GetStore().List() {
//		ep := *item.(*v1.Endpoints)
//
//		svcItem, exists := c.serviceByKey(ep.Name, ep.Namespace)
//		if !exists {
//			continue
//		}
//		svc := convertService(*svcItem, c.domainSuffix)
//		if svc == nil {
//			continue
//		}
//
//		for _, ss := range ep.Subsets {
//			for _, port := range ss.Ports {
//				svcPort, exists := svc.Ports.Get(port.Name)
//				if !exists {
//					continue
//				}
//
//				out = append(out, getEndpoints(ss.Addresses, proxyIP, c, port, svcPort, svc)...)
//				nrEP := getEndpoints(ss.NotReadyAddresses, proxyIP, c, port, svcPort, svc)
//				out = append(out, nrEP...)
//				if len(nrEP) > 0 && c.Env != nil {
//					c.Env.PushContext.Add(model.ProxyStatusEndpointNotReady, proxy.ID, proxy, "")
//				}
//			}
//		}
//	}
//	if len(out) == 0 {
//		if c.Env != nil {
//			c.Env.PushContext.Add(model.ProxyStatusNoService, proxy.ID, proxy, "")
//			status := c.Env.PushContext
//			if status == nil {
//				log.Infof("Empty list of services for pod %s %v", proxy.ID, c.Env)
//			}
//		} else {
//			log.Infof("Missing env, empty list of services for pod %s", proxy.ID)
//		}
//	}
//	return out, nil
//}
//
//func getEndpoints(addr []v1.EndpointAddress, proxyIP string, c *Controller,
//	port v1.EndpointPort, svcPort *model.Port, svc *model.Service) []*model.ServiceInstance {
//
//	var out []*model.ServiceInstance
//	for _, ea := range addr {
//		if proxyIP != ea.IP {
//			continue
//		}
//		labels, _ := c.pods.labelsByIP(ea.IP)
//		pod, exists := c.pods.getPodByIP(ea.IP)
//		az, sa := "", ""
//		if exists {
//			az, _ = c.GetPodAZ(pod)
//			sa = kubeToIstioServiceAccount(pod.Spec.ServiceAccountName, pod.GetNamespace(), c.domainSuffix)
//		}
//		out = append(out, &model.ServiceInstance{
//			Endpoint: model.NetworkEndpoint{
//				Address:     ea.IP,
//				Port:        int(port.Port),
//				ServicePort: svcPort,
//			},
//			Service:          svc,
//			Labels:           labels,
//			AvailabilityZone: az,
//			ServiceAccount:   sa,
//		})
//	}
//	return out
//}
//
//// GetIstioServiceAccounts returns the Istio service accounts running a serivce
//// hostname. Each service account is encoded according to the SPIFFE VSID spec.
//// For example, a service account named "bar" in namespace "foo" is encoded as
//// "spiffe://cluster.local/ns/foo/sa/bar".
//func (c *Controller) GetIstioServiceAccounts(hostname model.Hostname, ports []int) []string {
//	saSet := make(map[string]bool)
//
//	// Get the service accounts running the service, if it is deployed on VMs. This is retrieved
//	// from the service annotation explicitly set by the operators.
//	svc, err := c.GetService(hostname)
//	if err != nil {
//		// Do not log error here, as the service could exist in another registry
//		return nil
//	}
//	if svc == nil {
//		// Do not log error here as the service could exist in another registry
//		return nil
//	}
//
//	instances := make([]*model.ServiceInstance, 0)
//	// Get the service accounts running service within Kubernetes. This is reflected by the pods that
//	// the service is deployed on, and the service accounts of the pods.
//	for _, port := range ports {
//		svcinstances, err := c.InstancesByPort(hostname, port, model.LabelsCollection{})
//		if err != nil {
//			log.Warnf("InstancesByPort(%s:%d) error: %v", hostname, port, err)
//			return nil
//		}
//		instances = append(instances, svcinstances...)
//	}
//
//	for _, si := range instances {
//		if si.ServiceAccount != "" {
//			saSet[si.ServiceAccount] = true
//		}
//	}
//
//	for _, serviceAccount := range svc.ServiceAccounts {
//		sa := serviceAccount
//		saSet[sa] = true
//	}
//
//	saArray := make([]string, 0, len(saSet))
//	for sa := range saSet {
//		saArray = append(saArray, sa)
//	}
//
//	return saArray
//}




