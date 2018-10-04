package controllers

import (
	//"encoding/json"
	"github.com/json-iterator/go"
	//"k8s.io/api/core/v1"
	"github.com/astaxie/beego"
	"github.com/jukylin/istio-ui/models"
	"github.com/jukylin/istio-ui/pkg"
	//"github.com/judwhite/go-svc/svc"
	appv1 "k8s.io/api/apps/v1"
	//"github.com/ghodss/yaml"
	//cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	//"fmt"
	yaml2 "github.com/ghodss/yaml"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary


type ConfigController struct {
	beego.Controller
}


type listReturnItem struct{
	Name string `json:"name"`
	Namespace string `json:"namespace"`
	Labels string `json:"labels"`
	IsInject string `json:"is_inject"`
}


func (c *ConfigController) List() {
	deploysList := models.DeploysList()
	var list []listReturnItem
	var labels,isInject string

	for _, deployItem := range deploysList{
		deploy := deployItem.(*appv1.Deployment)

		if _, ok := deploy.Labels["version"]; ok {
			labels = deploy.Labels["version"]
		} else {
			labels = ""
		}

		if _, ok := deploy.Spec.Template.Annotations["sidecar.istio.io/status"] ; ok {
			isInject = "1"
		} else {
			isInject = "0"
		}

		lRI := listReturnItem{
			deploy.Name,
			deploy.Namespace,
			labels,
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
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": "不存在", "data" : nil}
		c.ServeJSON()
	}

	deploy := item.(*appv1.Deployment)

	Anno := deploy.GetObjectMeta().GetAnnotations()
	lastConfig := Anno["kubectl.kubernetes.io/last-applied-configuration"]

	yd, err := yaml2.JSONToYAML([]byte(lastConfig))
	err = pkg.InjectData(yd, namespace)
	if err != nil{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data" : nil}
		c.ServeJSON()
	}

	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "success", "data" : ""}
	c.ServeJSON()
}

