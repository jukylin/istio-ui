package pkg

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/jukylin/istio-ui/pkg/config/kube/crd"
	istiomodel "istio.io/istio/pilot/pkg/model"
	"github.com/astaxie/beego"
)

var clientset *kubernetes.Clientset
var configClient *crd.Client

func InitKubeClient()  {
	kubeConfigPath := beego.AppConfig.String("kube_config_dir")
	if kubeConfigPath == ""{
		kubeConfigPath = clientcmd.RecommendedHomeFile
	}
	//NewNonInteractiveDeferredLoadingClientConfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
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
	kubeConfigPath := beego.AppConfig.String("kube_config_dir")
	if kubeConfigPath == ""{
		kubeConfigPath = clientcmd.RecommendedHomeFile
	}

	client, err := crd.NewClient(kubeConfigPath, "",
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


