package controllers

import (
	"github.com/jukylin/istio-ui/pkg"
	"github.com/astaxie/beego"
)


type InjectController struct {
	beego.Controller
}

func (c *InjectController) Index() {
	config := c.Input().Get("config")
	if config == ""{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": "config is empty", "data" : nil}
		c.ServeJSON()
	}
	deploy, err := pkg.InjectData([]byte(config))
	if err != nil{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data" : nil}
		c.ServeJSON()
	}

	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "", "data" : deploy}
	c.ServeJSON()
}
