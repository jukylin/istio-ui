package controllers

import (
	"github.com/astaxie/beego"
	"istio.io/istio/pilot/pkg/config/kube/crd"
	"github.com/jukylin/istio-ui/pkg"
	istiomodel "istio.io/istio/pilot/pkg/model"
	"fmt"
)

type TemplateController struct {
	beego.Controller
}


/**
get istio config from local file
 */
func (this *TemplateController) GetMeshConfig() {

	config, err := pkg.GetMeshConfig()
	if err != nil{
		this.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
		this.ServeJSON()
	}

	fmt.Println(config)

	this.Data["json"] = map[string]interface{}{"code": 0, "msg": "success", "data" : config}
	this.ServeJSON()
}

/**
write to local file and post to k8s
 */
func (this *TemplateController) GetInjectConfig() {

	config, err := pkg.GetInjectConfigFromConfigMap()
	if err != nil{
		this.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
		this.ServeJSON()
	}

	this.Data["json"] = map[string]interface{}{"code": 0, "msg": "success", "data" : config}
	this.ServeJSON()
}

/**
delete local file and remote istio config
 */
func (this *TemplateController) Del() {
	name := this.Input().Get("name")
	namespace := this.Input().Get("namespace")

	fileName := name + ".yaml"
	configData, err := pkg.GetIstioConfig(fileName, namespace)
	if err != nil{
		this.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
		this.ServeJSON()
	}

	if configData == nil{
		this.Data["json"] = map[string]interface{}{"code": -1, "msg": "config is empty", "data": nil}
		this.ServeJSON()
	}

	var configs []istiomodel.Config
	configs, _, err = crd.ParseInputs(string(configData))

	if err != nil{
		this.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
		this.ServeJSON()
	}

	err = pkg.DelLocalIstioConfig(fileName, namespace)
	if err != nil{
		this.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
		this.ServeJSON()
	}

	err = pkg.DelRemoteIstioConfig(configs, namespace)
	if err != nil{
		this.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
		this.ServeJSON()
	}

	this.Data["json"] = map[string]interface{}{"code": 0, "msg": "success", "data" : ""}
	this.ServeJSON()
}