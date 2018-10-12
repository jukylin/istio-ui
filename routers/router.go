// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/jukylin/istio-ui/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.SetStaticPath("/static","views")
	beego.SetStaticPath("/lib","views/lib")
	beego.SetStaticPath("/css","views/css")
	beego.SetStaticPath("/js","views/js")
	beego.SetStaticPath("/config","views/config")
	beego.SetStaticPath("/components","views/components")
	beego.SetStaticPath("/pages","views/pages")

	beego.AutoRouter(&controllers.DeployController{})
	beego.AutoRouter(&controllers.Istio_ConfigController{})
	beego.AutoRouter(&controllers.InjectController{})
}
