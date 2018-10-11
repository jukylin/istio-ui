package controllers

import (
	"io"
	"os"
	"github.com/astaxie/beego"
	"github.com/jukylin/istio-ui/pkg"
)


type InjectController struct {
	beego.Controller
}

// curl -F "config=@bookinfo.yaml" http://localhost:9100/inject/file
func (this *InjectController) File() {

	f, h, err := this.GetFile("config")
	if err != nil {
		this.Ctx.ResponseWriter.Write([]byte(err.Error() + "\n"))
		this.StopRun()
	}
	config_path := beego.AppConfig.String("InjectUploadTmpFileDir") + h.Filename
	defer f.Close()
	defer func (){
		os.Remove(config_path)
	}()
	err = this.SaveToFile("config", config_path)
	if err != nil {
		this.Ctx.ResponseWriter.Write([]byte(err.Error() + "\n"))
		this.StopRun()
	}

	var in *os.File
	var reader io.Reader

	in, err = os.Open(config_path)
	if err != nil {
		this.Ctx.ResponseWriter.Write([]byte(err.Error() + "\n"))
		this.StopRun()
	}
	reader = in

	meshConfig, err := pkg.GetMeshConfigFromConfigMap()
	if err != nil {
		this.Ctx.ResponseWriter.Write([]byte(err.Error() + "\n"))
		this.StopRun()
	}

	injectConfig, err := pkg.GetInjectConfigFromConfigMap()
	if err != nil {
		this.Ctx.ResponseWriter.Write([]byte(err.Error() + "\n"))
		this.StopRun()
	}
	pkg.IntoResourceFile(injectConfig, meshConfig, reader, this.Ctx.ResponseWriter)
}


// curl -X POST --data-binary @bookinfo.yaml -H "Content-type: text/yaml" http://localhost:9100/inject/context
func (this *InjectController) Context() {

	meshConfig, err := pkg.GetMeshConfigFromConfigMap()
	if err != nil {
		this.Ctx.ResponseWriter.Write([]byte(err.Error() + "\n"))
		this.StopRun()
	}

	injectConfig, err := pkg.GetInjectConfigFromConfigMap()
	if err != nil {
		this.Ctx.ResponseWriter.Write([]byte(err.Error() + "\n"))
		this.StopRun()
	}

	pkg.IntoResourceFile(injectConfig, meshConfig, this.Ctx.Request.Body, this.Ctx.ResponseWriter)
}