package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"errors"
	"strings"
	"io"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"istio.io/istio/pilot/pkg/kube/inject"
	meshconfig "istio.io/api/mesh/v1alpha1"
	"istio.io/istio/pilot/pkg/kube/inject"
	"istio.io/istio/pilot/pkg/serviceregistry/kube"
	"github.com/ghodss/yaml"
	"istio.io/istio/pkg/log"
	appsv1 "k8s.io/api/apps/v1"
	"github.com/astaxie/beego"
	istiomodel "istio.io/istio/pilot/pkg/model"
	corev1 "k8s.io/api/core/v1"
)

var DeployStore DeployIndexStore


/**
Inject the injectConfig and meshConfi to rwa，then post to k8s
 */
func InjectData(raw []byte) (*appsv1.Deployment, error) {
	resource, err := GetResource(raw)
	if err != nil {
		return nil, err
	}

	ann, err := applyLastConfig(resource)
	if err != nil {
		return nil, err
	}

	var deploy *appsv1.Deployment
	err = yaml.Unmarshal(resource, &deploy)
	if err != nil {
		return nil, err
	}
	deploy.GetObjectMeta().SetAnnotations(ann)

	return deploy, nil
}

func GetResource(raw []byte) ([]byte, error) {
	meshConfig, err := GetMeshConfigFromConfigMap()
	if err != nil {
		return nil, err
	}

	injectConfig, err := GetInjectConfigFromConfigMap()
	if err != nil {
		return nil, err
	}

	return IntoResource(injectConfig, meshConfig, raw)
}

func UpdateDeploy(deploy *appsv1.Deployment, namespace string) error {
	_, err := GetKubeClent().AppsV1().Deployments(namespace).Update(deploy)
	if err != nil {
		return err
	}

	return nil
}


func applyLastConfig(resource []byte) (map[string]string, error)  {
	ann := make(map[string]string)
	json, err := yaml.YAMLToJSON(resource)

	if err != nil {
		return nil, err
	}
	ann[corev1.LastAppliedConfigAnnotation] = string(json)

	return ann, nil
}


/**
get mesh's config from k8s
 */
func GetMeshConfig() (string, error) {
	client := GetKubeClent()

	config, err := client.CoreV1().ConfigMaps(kube.IstioNamespace).Get("istio", metav1.GetOptions{})

	if err != nil {
		return "", fmt.Errorf("could not read valid configmap %q from namespace  %q: %v - "+
			"Use --meshConfigFile or re-run kube-inject with `-i <istioSystemNamespace> and ensure valid MeshConfig exists",
			"istio", kube.IstioNamespace, err)
	}
	// values in the data are strings, while proto might use a
	// different data type.  therefore, we have to get a value by a
	// key
	configYaml, exists := config.Data["mesh"]
	if !exists {
		return "", fmt.Errorf("missing configuration map key %q", "mesh")
	}

	return configYaml, nil
}


/**
get mesh's config from k8s
 */
func GetMeshConfigFromConfigMap() (*meshconfig.MeshConfig, error) {

	configYaml, err := GetMeshConfig()
	if err != nil {
		return nil, err
	}

	return ApplyMeshConfigDefaults(configYaml)
}

/**
get inject's config from k8s
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

func CheckIsExists(fileOrDir string) (bool, error) {
	_, err := os.Stat(fileOrDir)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CheckIstioConfigIsExists(filename, namespace  string) (bool, error) {
	istioConfigDir, err := getIstioConfigDir(namespace)
	if err != nil{
		return false, err
	}

	return CheckIsExists(istioConfigDir + "/" + filename)
}


func getIstioConfigDir(namespace string) (string, error) {
	istioConfigDir := beego.AppConfig.String("istio_config_dir")
	if istioConfigDir == ""{
		return "", errors.New("istio_config_dir is empty")
	}

	return checkOrCreateDir(istioConfigDir, namespace)
}


func getBackUpIstioConfigDir(namespace string) (string, error) {
	istioBackUpDir := beego.AppConfig.String("back_up_dir")
	if istioBackUpDir == ""{
		return "", errors.New("istio_back_up_dir is empty")
	}

	return checkOrCreateDir(istioBackUpDir, namespace)
}

func checkOrCreateDir(dirName, namespace string) (string, error) {
	exists, err := CheckIsExists(dirName)
	if err != nil{
		return "", err
	}

	if !exists {
		return "", errors.New(dirName + " not exists")
	}

	if namespace != "" {
		dirName = dirName + "/" + namespace
		exists, err = CheckIsExists(dirName)
		if err != nil {
			return "", err
		}

		if !exists {
			err = os.Mkdir(dirName, os.ModePerm)
			if err != nil {
				return "", err
			}
		}
	}

	return dirName, nil
}


/**
write to local file
 */
func WriteIstioConfig(data []byte, filename, namespace  string) error {
	istioConfigDir, err := getIstioConfigDir(namespace)
	if err != nil{
		return err
	}

	BackUp(filename, namespace)

	err = ioutil.WriteFile(istioConfigDir + "/" + filename, []byte(data), 0644)
	if err != nil{
		return err
	}

	return nil
}


func DelLocalIstioConfig(filename, namespace  string) error {
	istioConfigDir, err := getIstioConfigDir(namespace)
	if err != nil{
		return err
	}

	BackUp(filename, namespace)

	err = os.Remove(istioConfigDir + "/" + filename)
	if err != nil{
		return err
	}

	return nil
}

//copy filename.yaml to backup/filename.yaml
func BackUp(filename, namespace  string) error {

	back_up, err := beego.AppConfig.Bool("back_up")
	if err != nil{
		beego.Warn(err.Error())
		back_up = false
	}

	if back_up == false{
		beego.Info("istio config not backing up")
		return nil
	}

	exists, err := CheckIstioConfigIsExists(filename, namespace)
	if err != nil{
		beego.Warn(err.Error())
	}

	if exists {
		beego.Info(namespace + "'s " + filename + " istio config exists，"+ "prepare for backup")
		//内容是否相同

		istioConfigDir, err := getIstioConfigDir(namespace)
		if err != nil{
			beego.Warn(err.Error())
			return err
		}

		srcFile ,err := os.Open(istioConfigDir + "/" + filename)
		if err != nil{
			beego.Warn(err.Error())
			return err
		}
		defer srcFile.Close()

		istioBackUpDir, err := getBackUpIstioConfigDir(namespace)
		if err != nil{
			beego.Warn(err.Error())
			return nil
		}

		destFile,err := os.OpenFile(istioBackUpDir  + "/" + filename, os.O_WRONLY|os.O_CREATE,0644)
		if err != nil{
			beego.Warn(err.Error())
			return err
		}
		defer destFile.Close()
		_, err = io.Copy(destFile,srcFile)
		if err != nil{
			beego.Error(namespace + "'s " + filename + " backup error：", err.Error())
		}else{
			beego.Info(namespace + "'s " + filename + " backup completed")
		}
		return err
	}else{
		beego.Warn(namespace + "'s " + filename + " istio config not exists，"+ " not backing up")
		return nil
	}
}

//check istio config' backup exists
func CheckBackUpIsExists(filename, namespace  string) (bool, error) {
	istioBackUpDir, err := getBackUpIstioConfigDir(namespace)
	if err != nil{
		return false, nil
	}

	return CheckIsExists(istioBackUpDir + "/" + filename)
}


/**
post to k8s
 */
func DelRemoteIstioConfig(configs []istiomodel.Config, namespace string) error {

	client := GetConfigClient()

	for _, config := range configs {
		config.Namespace = handleNamespaces(config.Namespace, namespace)
		exists := client.Get(config.ConfigMeta.Type, config.Name, config.Namespace)

		if  exists == nil {
			beego.Warn(config.ConfigMeta.Type + "'s " + config.Name + "not exists")
		}else{
			config.ResourceVersion = exists.ResourceVersion
			err := client.Delete(config.ConfigMeta.Type, config.Name, config.Namespace)
			if err != nil {
				return err
			}
			beego.Info("Delete config "+ config.ConfigMeta.Type +"'s " + config.Name)
		}
	}

	return nil
}


/**
get istio config from local file
 */
func GetIstioConfig(filename, namespace string)([]byte, error){
	istioConfigDir, err := getIstioConfigDir(namespace)
	if err != nil{
		return nil, err
	}

	data, err := ioutil.ReadFile(istioConfigDir + "/" + filename)
	if err != nil{
		return nil, err
	}

	return data, nil
}


/**
get istio config from local file
 */
func GetBackedUpIstioConfig(filename, namespace string)([]byte, error){
	backUpDir, err := getBackUpIstioConfigDir(namespace)
	if err != nil{
		return nil, err
	}

	data, err := ioutil.ReadFile(backUpDir + "/" + filename)
	if err != nil{
		return nil, err
	}

	return data, nil
}

/**
post istio config to k8s
If there is an update, otherwise create
 */
func PostIstioConfig(configs []istiomodel.Config, namespace string) error {
	client := GetConfigClient()

	for _, config := range configs {
		config.Namespace = handleNamespaces(config.Namespace, namespace)
		exists := client.Get(config.ConfigMeta.Type, config.Name, config.Namespace)

		if  exists == nil {
			rev, err := client.Create(config)
			if err != nil {
				return err
			}
			beego.Info("Create config "+ config.Key() +" at revision "+ rev +"\n")
		}else{
			config.ResourceVersion = exists.ResourceVersion
			rev, err := client.Update(config)
			if err != nil {
				return err
			}
			beego.Info("Update config "+ config.Key() +" at revision "+ rev +"\n")
		}
	}

	return nil
}


func handleNamespaces(objectNamespace , namespace string) string {
	if objectNamespace != "" {
		return objectNamespace
	}

	return namespace
}

func GetWorkNameSpace() []string {
	work_namespace := beego.AppConfig.String("work_namespace")
	nameSpaces := strings.Split(work_namespace, ",")
	return nameSpaces
}


func InitDeployIndexStore()  {
	DeployStore = NewDeployIndexStore()
}