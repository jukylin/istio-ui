package controllers

import (
	"strings"
	"github.com/astaxie/beego"
	"gopkg.in/yaml.v2"
	"istio.io/istio/pilot/pkg/config/kube/crd"
	"github.com/jukylin/istio-ui/pkg"
	istiomodel "istio.io/istio/pilot/pkg/model"
)


type Istio_ConfigController struct {
	beego.Controller
}


func (c *Istio_ConfigController) Get() {
	name := c.Input().Get("name")
	//namespace := c.Input().Get("namespace")

	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "success", "data" : name}
	c.ServeJSON()
}


func (c *Istio_ConfigController) Save() {
	name := c.Input().Get("name")
	namespace := c.Input().Get("namespace")
	configStr := c.Input().Get("config")

	configStrArr := strings.Split(configStr, "---")

	var configs []istiomodel.Config
	istioKind := &crd.IstioKind{}

	for _,v := range  configStrArr{
		err := yaml.Unmarshal([]byte(v), istioKind)
		if err != nil{
			c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
			c.ServeJSON()
		}
		schema, exists := istiomodel.IstioConfigTypes.GetByType(crd.CamelCaseToKabobCase(istioKind.Kind))
		if !exists {
			c.Data["json"] = map[string]interface{}{"code": -1, "msg": "unrecognized type " + istioKind.Kind, "data": nil}
			c.ServeJSON()
		}

		config, err := crd.ConvertObject(schema, istioKind, "")
		if err != nil {
			c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
			c.ServeJSON()
		}

		if err := schema.Validate(config.Name, config.Namespace, config.Spec); err != nil {
			c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
			c.ServeJSON()
		}

		config, err = crd.ConvertObject(schema, istioKind, "")
		if err != nil {
			c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
			c.ServeJSON()
		}
		configs = append(configs, *config)
	}

	// write to istio_config_dir
	err := pkg.WriteIstioConfig([]byte(configStr), name + ".yaml", namespace)
	if err != nil{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
		c.ServeJSON()
	}

	//post to k8s
	err = pkg.PostIstioConfig(configs)
	if err != nil{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": nil}
		c.ServeJSON()
	}

	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "success", "data" : ""}
	c.ServeJSON()
}



func (c *Istio_ConfigController) Verify() {
	name := c.Input().Get("name")
	//namespace := c.Input().Get("namespace")
	//是否yaml

	//是否istio_config
	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "success", "data" : name}
	c.ServeJSON()
}