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
	extvb1 "k8s.io/api/extensions/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	appsvb1 "k8s.io/api/apps/v1beta1"
	appsvb2 "k8s.io/api/apps/v1beta2"
)

/**
注入信息并把数据更新到k8s
 */
func InjectData(raw []byte, namespace string) error {
	meshConfig, err := GetMeshConfigFromConfigMap()
	if err != nil {
		return err
	}

	injectConfig, err := GetInjectConfigFromConfigMap()
	if err != nil {
		return err
	}

	resource, err := IntoResource(injectConfig, meshConfig, raw)
	if err != nil {
		return err
	}
	obj, err := fromRawToObject(raw)
	group := obj.GetObjectKind().GroupVersionKind().GroupKind().Group

	ann, err := applyLastConfig(resource)
	if err != nil {
		return err
	}

	if obj.GetObjectKind().GroupVersionKind().Version == "v1beta1" && group == "extensions"{
		var deploy *extvb1.Deployment
		yaml.Unmarshal(resource, &deploy);
		deploy.GetObjectMeta().SetAnnotations(ann)
		_, err = GetKubeClent().ExtensionsV1beta1().Deployments(namespace).Update(deploy)
	}else if obj.GetObjectKind().GroupVersionKind().Version == "v1" && group == "apps"{
		var deploy *appsv1.Deployment
		yaml.Unmarshal(resource, &deploy)
		deploy.GetObjectMeta().SetAnnotations(ann)
		_, err = GetKubeClent().AppsV1().Deployments(namespace).Update(deploy)
	}else if obj.GetObjectKind().GroupVersionKind().Version == "v1beta1" && group == "apps"{
		var deploy *appsvb1.Deployment
		yaml.Unmarshal(resource, &deploy)
		deploy.GetObjectMeta().SetAnnotations(ann)
		_, err = GetKubeClent().AppsV1beta1().Deployments(namespace).Update(deploy)
	}else if obj.GetObjectKind().GroupVersionKind().Version == "v1beta2" && group == "apps"{
		var deploy *appsvb2.Deployment
		yaml.Unmarshal(resource, &deploy)
		deploy.GetObjectMeta().SetAnnotations(ann)
		_, err = GetKubeClent().AppsV1beta2().Deployments(namespace).Update(deploy)
	}

	return err
}

func applyLastConfig(resource []byte) (map[string]string, error)  {
	ann := make(map[string]string)
	json, err := yaml.YAMLToJSON(resource)

	if err != nil {
		return nil, err
	}
	ann[LastAppliedConfigAnnotation] = string(json)

	return ann, nil
}

/**
从k8s获取mesh配置信息
 */
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

/**
从k8s获取配置信息
 */
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
