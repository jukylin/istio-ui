package controllers

import (
	"github.com/json-iterator/go"
	"github.com/astaxie/beego"
	"github.com/jukylin/istio-ui/models"
	"github.com/jukylin/istio-ui/pkg"
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	yaml2 "github.com/ghodss/yaml"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary


type ConfigController struct {
	beego.Controller
}


type listReturnItem struct{
	Name string `json:"name"`
	Namespace string `json:"namespace"`
	Version string `json:"version"`
	Create_time metav1.Time `json:"create_time"`
	IsInject string `json:"is_inject"`
}


func (c *ConfigController) List() {
	deploysList := models.DeploysList()
	var list []listReturnItem
	var version,isInject string

	for _, deployItem := range deploysList{
		deploy := deployItem.(*appv1.Deployment)

		if _, ok := deploy.Labels["version"]; ok {
			version = deploy.Labels["version"]
		} else if _, ok := deploy.Spec.Template.Labels["version"]; ok{
			version = deploy.Spec.Template.Labels["version"]
		} else {
			version = ""
		}

		if _, ok := deploy.Spec.Template.Annotations["sidecar.istio.io/status"] ; ok {
			isInject = "1"
		} else {
			isInject = "0"
		}

		lRI := listReturnItem{
			deploy.Name,
			deploy.Namespace,
			version,
			deploy.CreationTimestamp,
			isInject,
		}

		list = append(list, lRI)
	}

	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "success", "data" : list}
	c.ServeJSON()
}



func (c *ConfigController) Inject() {
	name := c.Input().Get("name")
	namespace := c.Input().Get("namespace")

	item, exists, err := models.GetByKey(namespace + "/" + name)
	if err != nil{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data" : nil}
		c.ServeJSON()
	}

	if exists != true{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": "not exists", "data" : nil}
		c.ServeJSON()
	}

	deploy := item.(*appv1.Deployment)

	if _, ok := deploy.Spec.Template.Annotations["sidecar.istio.io/status"] ; ok {
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": "has injected", "data" : nil}
		c.ServeJSON()
	}

	Anno := deploy.GetObjectMeta().GetAnnotations()
	if _, ok := Anno[pkg.LastAppliedConfigAnnotation]; !ok{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": "lost last configuration", "data" : nil}
		c.ServeJSON()
	}

	lastConfig := Anno[pkg.LastAppliedConfigAnnotation]
	yd, err := yaml2.JSONToYAML([]byte(lastConfig))

	err = pkg.InjectData(yd, namespace)
	if err != nil{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data" : nil}
		c.ServeJSON()
	}

	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "success", "data" : ""}
	c.ServeJSON()
}

