# istio-ui

## &nbsp;&nbsp;&nbsp;&nbsp;istio-ui用于管理istio配置文件，目的是减轻运维的配置工作。主要实现：注入，istio配置和模板（还在开发中）等功能。
## 为了保证注入和配置的原生性，参考和使用了istio的源码。


### 三种注入方式
* 一键注入
  > 基于运行中的服务```Deployment:apps/v1```进行注入，使用这种方式服务会被重新部署
* 文件上传注入
  > 将你需要注入的文件发送到远程api接口
  
  > kubectl apply -f <(curl -F "config=@samples/bookinfo/platform/kube/bookinfo.yaml" http://localhost:9100/inject/file)
* 内容注入
  > 将你需要注入的内容发送到远程api接口
  
  > kubectl apply -f <(curl -X POST --data-binary @samples/bookinfo/platform/kube/bookinfo.yaml -H "Content-type: text/yaml" http://localhost:9100/inject/context)


### 安装
> 安装前先确认已安装k8s
* docker
> &nbsp;&nbsp;&nbsp;&nbsp;设置 [KUBECONFIG](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/#create-a-second-configuration-file)

> &nbsp;&nbsp;&nbsp;&nbsp;docker run -itd -v $KUBECONFIG:$HOME/.kube/config -p9100:9100 --name istio-ui --env KUBECONFIG=$HOME/.kube/config registry.cn-shenzhen.aliyuncs.com/jukylin/istio-ui:latest

* k8s
> &nbsp;&nbsp;&nbsp;&nbsp;kubectl apply -f https://raw.githubusercontent.com/jukylin/istio-ui/master/istio-ui.yaml


### 配置
> 使用环境变量

* ISTIO_CONFIG_DIR  配置文件存放目录，默认：/data/www/istio_config

* INJECT_UPLOAD_TMP_FILE_DIR  文件上传临时存放目录，默认：/data/www/istio_upload

* FILTER_NAMESPACES  被过滤的namespaces，默认：kube-public,kube-system,istio-system

* FILTER_NAME  不进行注入的name，默认：redis,mysql,istio-ui
