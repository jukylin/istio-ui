define(function (require,exports,module) {
    module.exports = function() {
        var temp = require('text!/pages/global-istio-config.html');
        require("/lib/vue-resource.min.js")
        Vue.http.options.emulateJSON = true;
        return {
            template: temp,
            data: function (){
                return {
                    istio_config: "",
                    handle_name: "gateway",
                    namespace_options: [],
                    search_namespace: "default",
                    backedup: false
                }
            },
            mounted : function () {
                this.getHandleNmae();
                this.getIstioConfig();
                this.getWorkNameSpace();
            },
            methods : {
                getIstioConfig: function () {
                    var _self = this;
                    this.$http.get('/istio_config/get?name='+ _self.handle_name + '&namespace=' + _self.search_namespace)
                        .then(function (resp) {
                            if (resp.body.code === 0) {
                                _self.istio_config = resp.body.data.data
                                _self.backedup = resp.body.data.backedup
                            }else{
                                this.$message({
                                    message: resp.body.msg,
                                    type: 'error'
                                });
                            }
                        });
                },
                saveIstioConfig: function () {
                    var _self = this;
                    if(_self.istio_config == ""){
                        this.$message({
                            message: "请填写配置信息",
                            type: 'error'
                        });
                        return false
                    }
                    this.$http.post('/istio_config/save',
                        {name: _self.handle_name, namespace: _self.search_namespace, config: _self.istio_config})
                        .then(function (resp) {
                            if (resp.body.code === 0) {
                                this.$message({
                                    message: "保存成功",
                                    type: 'success'
                                });
                            }else{
                                this.$message({
                                    message: resp.body.msg,
                                    type: 'error'
                                });
                            }
                        });
                },
                getBackUp: function () {
                    var _self = this;
                    _self.backedup = false;
                    this.$http.get('/istio_config/getbackup?name='+ _self.handle_name + '&namespace=' + _self.search_namespace)
                        .then(function (resp) {
                            if (resp.body.code === 0) {
                                _self.istio_config = resp.body.data
                            }else{
                                this.$message({
                                    message: resp.body.msg,
                                    type: 'error'
                                });
                            }
                        });
                },
                delIstioConfig:function () {
                    var _self = this;

                    this.$http.post('/istio_config/del', {name: _self.handle_name, namespace: _self.search_namespace})
                        .then(function (resp) {
                            if (resp.body.code === 0) {
                                _self.istio_config = ""
                                this.$message({
                                    message: "删除成功",
                                    type: 'success'
                                });
                            }else{
                                this.$message({
                                    message: resp.body.msg,
                                    type: 'error'
                                });
                            }
                        });
                },
                getWorkNameSpace: function () {
                    var _self = this;
                    this.$http.get('/deploy/getworknamespaces?').then(function (resp) {
                        if (resp.body.code === 0) {
                            _self.namespace_options = resp.body.data;
                        }else{
                            this.$message({
                                message: resp.body.msg,
                                type: 'error'
                            });
                        }
                    });
                },
                getHandleNmae: function () {
                    var _self = this;
                    if(this.$route.query.type == ""){
                        this.$message({
                            message: "缺少type",
                            type: 'error'
                        });
                        return false;
                    }
                    _self.handle_name = this.$route.query.type
                },
                reloadIstioConfig:function () {
                    this.getHandleNmae()
                    this.getIstioConfig();
                }
            },
            watch:{
                '$route': 'reloadIstioConfig'
            }
        };
    };
});
