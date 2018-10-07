package controllers

import (
	"github.com/astaxie/beego"
	"istio.io/istio/pilot/pkg/config/kube/crd"
	"github.com/jukylin/istio-ui/pkg"
	istiomodel "istio.io/istio/pilot/pkg/model"
)


type Istio_ConfigController struct {
	beego.Controller
}

/**
get istio config from local file
 */
func (c *Istio_ConfigController) Get() {
	var istio_config []byte
	name := c.Input().Get("name")
	namespace := c.Input().Get("namespace")
	exists, err := pkg.CheckIstioConfigIsExists(name + ".yaml", namespace)
	if err != nil{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
		c.ServeJSON()
	}

	if exists {
		istio_config, err = pkg.GetIstioConfig(name + ".yaml", namespace)
		if err != nil{
			c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
			c.ServeJSON()
		}
	}else{
		istio_config = nil
	}

	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "success", "data" : string(istio_config)}
	c.ServeJSON()
}

/**
write to local file and post to k8s
 */
func (c *Istio_ConfigController) Save() {
	name := c.Input().Get("name")
	namespace := c.Input().Get("namespace")
	configStr := c.Input().Get("config")

	var configs []istiomodel.Config

	configs, _, err := crd.ParseInputs(configStr)

	if err != nil{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
		c.ServeJSON()
	}

	// write to istio_config_dir
	err = pkg.WriteIstioConfig([]byte(configStr), name + ".yaml", namespace)
	if err != nil{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
		c.ServeJSON()
	}

	//post to k8s
	err = pkg.PostIstioConfig(configs, namespace)
	if err != nil{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
		c.ServeJSON()
	}

	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "success", "data" : ""}
	c.ServeJSON()
}

/**
delete local file and remote istio config
 */
func (c *Istio_ConfigController) Del() {
	name := c.Input().Get("name")
	namespace := c.Input().Get("namespace")
	configStr := c.Input().Get("config")

	var configs []istiomodel.Config

	configs, _, err := crd.ParseInputs(configStr)

	if err != nil{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
		c.ServeJSON()
	}

	err = pkg.DelLocalIstioConfig(name + ".yaml", namespace)
	if err != nil{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
		c.ServeJSON()
	}

	err = pkg.DelRemoteIstioConfig(configs, namespace)
	if err != nil{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
		c.ServeJSON()
	}

	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "success", "data" : ""}
	c.ServeJSON()
}