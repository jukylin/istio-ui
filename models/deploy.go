package models

import (
	"time"
	"github.com/jukylin/istio-ui/pkg"
    "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var controller *pkg.Controller

func InitController() {

	options := pkg.ControllerOptions{
		DomainSuffix : "cluster.local",
		ResyncPeriod : 60*time.Second,
		WatchedNamespace : "",
	}

	controller = pkg.NewController(pkg.GetKubeClent(), options)
}

func Run(stop <-chan struct{})  {
	go controller.Run(stop)
}

func DeploysList(deployIndexs []string) []interface{} {
	return controller.GetDeployList(deployIndexs)
}


func ListKeys() []string {
	return controller.ListKeys()
}

func GetByKey(key string) (item interface{}, exists bool, err error) {
	return controller.GetByKey(key)
}