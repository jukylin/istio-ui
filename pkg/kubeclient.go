package pkg

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var clientset *kubernetes.Clientset

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

func GetKubeClent() *kubernetes.Clientset {
	return clientset
}
