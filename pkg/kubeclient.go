package pkg

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"istio.io/istio/pilot/pkg/config/kube/crd"
	istiomodel "istio.io/istio/pilot/pkg/model"
)

var clientset *kubernetes.Clientset
var configClient *crd.Client

func InitKubeClient()  {

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/apple/.kube/config")
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

func InitConfigClient()  {
	client, err := crd.NewClient("/Users/apple/.kube/config", "",
		istiomodel.IstioConfigTypes, "")
	if err != nil {
		panic(err.Error())
	}
	configClient = client
}

func GetKubeClent() *kubernetes.Clientset {
	return clientset
}


func GetConfigClient() *crd.Client {
	return configClient
}


