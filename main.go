package main

import (
	_ "github.com/jukylin/istio-ui/routers"
	"github.com/astaxie/beego"
	"github.com/jukylin/istio-ui/models"
	"github.com/jukylin/istio-ui/pkg"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	stop := make(chan struct{})
	pkg.InitKubeClient()
	pkg.InitConfigClient()
	models.InitController()
	models.Run(stop)

	beego.Run()
	close(stop)
}
