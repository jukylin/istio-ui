package controllers

import (
	"github.com/astaxie/beego"
	"istio.io/istio/pilot/pkg/config/kube/crd"
	"github.com/jukylin/istio-ui/pkg"
	istiomodel "istio.io/istio/pilot/pkg/model"
	"fmt"
)

type Istio_ConfigController struct {
	beego.Controller
}


/**
get istio config from local file
 */
func (this *Istio_ConfigController) Get() {
	var istio_config []byte
	name := this.Input().Get("name")
	namespace := this.Input().Get("namespace")
	exists, err := pkg.CheckIstioConfigIsExists(name + ".yaml", namespace)
	if err != nil{
		this.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
		this.ServeJSON()
	}

	if exists {
		istio_config, err = pkg.GetIstioConfig(name + ".yaml", namespace)
		if err != nil{
			this.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
			this.ServeJSON()
		}
	}else{
		istio_config = nil
	}

	this.Data["json"] = map[string]interface{}{"code": 0, "msg": "success", "data" : string(istio_config)}
	this.ServeJSON()
}

/**
write to local file and post to k8s
 */
func (this *Istio_ConfigController) Save() {
	name := this.Input().Get("name")
	namespace := this.Input().Get("namespace")
	configStr := this.Input().Get("config")

	var configs []istiomodel.Config
	configs, _, err := crd.ParseInputs(configStr)

	if err != nil{
		this.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
		this.ServeJSON()
	}

	// write to istio_config_dir
	err = pkg.WriteIstioConfig([]byte(configStr), name + ".yaml", namespace)
	if err != nil{
		this.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
		this.ServeJSON()
	}

	//post to k8s
	err = pkg.PostIstioConfig(configs, namespace)
	if err != nil{
		this.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
		this.ServeJSON()
	}

	this.Data["json"] = map[string]interface{}{"code": 0, "msg": "success", "data" : ""}
	this.ServeJSON()
}

/**
delete local file and remote istio config
 */
func (this *Istio_ConfigController) Del() {
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