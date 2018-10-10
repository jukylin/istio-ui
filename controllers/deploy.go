package controllers

import (
	"strconv"
	"strings"
	"github.com/json-iterator/go"
	"github.com/astaxie/beego"
	"github.com/jukylin/istio-ui/models"
	"github.com/jukylin/istio-ui/pkg"
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	yaml2 "github.com/ghodss/yaml"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary


type DeployController struct {
	beego.Controller
}


type listReturnItem struct{
	Name string `json:"name"`
	Namespace string `json:"namespace"`
	Version string `json:"version"`
	Create_time metav1.Time `json:"create_time"`
	IsInject string `json:"is_inject"`
}


func (c *DeployController) List() {
	namespace := c.Input().Get("namespace")
	get_page := c.Input().Get("page")
	get_pagesize := c.Input().Get("page_size")
	name := c.Input().Get("name")

	if get_page == ""{
		get_page = "1"
	}

	if get_pagesize == ""{
		get_pagesize = "10"
	}

	if namespace == ""{
		namespace = "default"
	}

	var deployIndexs []string
	var total int
	if name != ""{
		deployIndexs = []string{name}
		total = 1
	}else {
		page, err := strconv.Atoi(get_page)
		if err != nil {
			c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": ""}
			c.ServeJSON()
		}

		pagesize, err := strconv.Atoi(get_pagesize)
		if err != nil {
			c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data": ""}
			c.ServeJSON()
		}

		total = pkg.DeployIndexLen(namespace)

		start := (page - 1) * 10
		end := pagesize

		if total < start+end {
			end = 0
		}

		deployIndexs = pkg.GetDeployIndexLimit(start, end, namespace)
	}

	deploysList := models.DeploysList(deployIndexs)
	var list []listReturnItem
	var version,isInject string

	for _, deployItem := range deploysList{
		deploy := deployItem.(*appv1.Deployment)

		if _, ok := deploy.Labels["version"]; ok {
			version = deploy.Labels["version"]
		} else if _, ok := deploy.Spec.Template.Labels["version"]; ok{
			version = deploy.Spec.Template.Labels["version"]
		} else {
			version = ""
		}

		if _, ok := deploy.Spec.Template.Annotations["sidecar.istio.io/status"] ; ok {
			isInject = "1"
		} else {
			isInject = "0"
		}

		lRI := listReturnItem{
			deploy.Name,
			deploy.Namespace,
			version,
			deploy.CreationTimestamp,
			isInject,
		}

		list = append(list, lRI)
	}

	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "success", "data" : map[string]interface{}{"list":list, "total":total}}
	c.ServeJSON()
}



func (c *DeployController) Inject() {
	name := c.Input().Get("name")
	namespace := c.Input().Get("namespace")

	item, exists, err := models.GetByKey(namespace + "/" + name)
	if err != nil{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data" : nil}
		c.ServeJSON()
	}

	if exists != true{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": "not exists", "data" : nil}
		c.ServeJSON()
	}

	deploy := item.(*appv1.Deployment)

	if _, ok := deploy.Spec.Template.Annotations["sidecar.istio.io/status"] ; ok {
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": "has injected", "data" : nil}
		c.ServeJSON()
	}

	Anno := deploy.GetObjectMeta().GetAnnotations()
	if _, ok := Anno[pkg.LastAppliedConfigAnnotation]; !ok{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": "lost last configuration", "data" : nil}
		c.ServeJSON()
	}

	lastConfig := Anno[pkg.LastAppliedConfigAnnotation]
	yd, err := yaml2.JSONToYAML([]byte(lastConfig))

	deploy, err = pkg.InjectData(yd)
	if err != nil{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data" : nil}
		c.ServeJSON()
	}

	err = pkg.UpdateDeploy(deploy, namespace)
	if err != nil{
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": err.Error(), "data" : nil}
		c.ServeJSON()
	}

	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "success", "data" : ""}
	c.ServeJSON()
}

/**
filter namespace eg:kube-public,kube-system
 */
func (c *DeployController) GetWorkNameSpaces()  {
	filterNamespace := beego.AppConfig.String("filter_namespace")
	filterNamespaces := strings.Split(filterNamespace, ",")
	nameSpaceList := models.NameSpacesList()

	var flag bool
	var allowNameSpace []string
	for _, nameSpace := range nameSpaceList{
		flag = false
		for _, filterNameSpace := range filterNamespaces{
			if filterNameSpace == nameSpace{
				flag = true
				continue
			}
		}

		if !flag {
			allowNameSpace = append(allowNameSpace, nameSpace)
		}
	}

	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "success", "data" : allowNameSpace}
	c.ServeJSON()
}

func (c *DeployController) GetDeployIndex()  {
	deployIndexs := pkg.GetAllDeployIndex("default")
	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "success", "data" : deployIndexs}
	c.ServeJSON()
}

