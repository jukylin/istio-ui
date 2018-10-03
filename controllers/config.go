package controllers

import (
	//"encoding/json"
	"github.com/json-iterator/go"
	"k8s.io/api/core/v1"
	"github.com/astaxie/beego"
	"github.com/jukylin/istio-ui/models"
	"github.com/jukylin/istio-ui/pkg"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary


type ConfigController struct {
	beego.Controller
}


type listReturnItem struct{
	Name string `json:"name"`
	Namespace string `json:"namespace"`
	Status string `json:"status"`
	Labels string `json:"labels"`
	IsInject string `json:"is_inject"`
}


func (c *ConfigController) List() {
	podsList := models.PodsList()
	var list []listReturnItem
	//var item []listReturnItem
	var labels,isInject string

	for _, podItem := range podsList{
		pod := podItem.(*v1.Pod)

		if _, ok := pod.Labels["version"]; ok {
			labels = pod.Labels["version"]
		} else {
			labels = ""
		}
		if _, ok := pod.Annotations["sidecar.istio.io/status"]; ok {
			isInject = "1"
		} else {
			isInject = "0"
		}

		lRI := listReturnItem{
			pod.Name,
			pod.Namespace,
			string(pod.Status.Phase),
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

	data, err := pkg.GetInjectData(item.(*v1.Pod))
	if err != nil{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data" : nil}
		c.ServeJSON()
	}

	err = pkg.PatchData(namespace, name, data)
	if err != nil{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data" : nil}
		c.ServeJSON()
	}



	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "success", "data" : nil}
	c.ServeJSON()
}

