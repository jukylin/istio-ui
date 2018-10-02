package pkg

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"istio.io/istio/pilot/pkg/kube/inject"
	meshconfig "istio.io/api/mesh/v1alpha1"
	"istio.io/istio/pilot/pkg/kube/inject"
	"istio.io/istio/pilot/pkg/serviceregistry/kube"
	"github.com/ghodss/yaml"
	"istio.io/istio/pkg/log"
	"k8s.io/api/core/v1"
)

func GetInjectData(pod *v1.Pod) []byte {
	meshConfig, _ := GetMeshConfigFromConfigMap()
	injectConfig, _ := GetInjectConfigFromConfigMap()
fmt.Printf("%+v \n", meshConfig.DefaultConfig)
	//spec, status, err := injectionData(wh.sidecarConfig.Template, wh.sidecarTemplateVersion,
	// &pod.ObjectMeta, &pod.Spec, &pod.ObjectMeta, wh.meshConfig.DefaultConfig, wh.meshConfig) // nolint: lll
	spec, status, _ := injectionData(injectConfig, sidecarTemplateVersionHash(injectConfig), &pod.ObjectMeta,
		&pod.Spec,
			&pod.ObjectMeta,
				meshConfig.DefaultConfig,
					meshConfig)

	annotations := map[string]string{annotationStatus.name: status}

	patchBytes, _ := createPatch(pod, injectionStatus(pod), annotations, spec)

	return patchBytes
}

func GetMeshConfigFromConfigMap() (*meshconfig.MeshConfig, error) {
	client := GetKubeClent()

	config, err := client.CoreV1().ConfigMaps(kube.IstioNamespace).Get("istio", metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not read valid configmap %q from namespace  %q: %v - "+
			"Use --meshConfigFile or re-run kube-inject with `-i <istioSystemNamespace> and ensure valid MeshConfig exists",
			"istio", kube.IstioNamespace, err)
	}
	// values in the data are strings, while proto might use a
	// different data type.  therefore, we have to get a value by a
	// key
	configYaml, exists := config.Data["mesh"]
	if !exists {
		return nil, fmt.Errorf("missing configuration map key %q", "mesh")
	}

	return ApplyMeshConfigDefaults(configYaml)
}


func GetInjectConfigFromConfigMap() (string, error) {
	client := GetKubeClent()

	config, err := client.CoreV1().ConfigMaps(kube.IstioNamespace).Get("istio-sidecar-injector", metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("could not find valid configmap %q from namespace  %q: %v - "+
			"Use --injectConfigFile or re-run kube-inject with `-i <istioSystemNamespace> and ensure istio-inject configmap exists",
			"istio-sidecar-injector", kube.IstioNamespace, err)
	}
	// values in the data are strings, while proto might use a
	// different data type.  therefore, we have to get a value by a
	// key
	injectData, exists := config.Data["config"]
	if !exists {
		return "", fmt.Errorf("missing configuration map key %q in %q",
			"config", "istio-sidecar-injector")
	}
	var injectConfig inject.Config
	if err := yaml.Unmarshal([]byte(injectData), &injectConfig); err != nil {
		return "", fmt.Errorf("unable to convert data from configmap %q: %v",
			"istio-sidecar-injector", err)
	}
	log.Debugf("using inject template from configmap %q", "istio-sidecar-injector")
	return injectConfig.Template, nil
}
