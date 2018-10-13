define(function (require,exports,module) {
    module.exports = function() {
        var temp = require('text!/pages/deploy.html');
        require("https://cdn.jsdelivr.net/npm/vue-resource@1.5.1")
        Vue.http.options.emulateJSON = true;

        return {
            template: temp,
            data: function (){
                return {
                    tableData: [],
                    is_show_istio_config: false,
                    sure_inject: false,
                    istio_config: "",
                    handle_name: "",
                    handle_namespace: "",
                    handle_raw : [],
                    search_name: "",
                    search_namespace: "default",
                    pagination: {
                        currentPage: 1,
                        pageSize: 10,
                        total: 0
                    },
                    namespace_options: []
                }
            },
            mounted : function () {
                this.getList();
                this.getWorkNameSpace();
            },
            methods : {
                getIsInject: function (row) {
                    return row.is_inject == 1 ? '是' : '否';
                },
                handleInject: function () {
                    var _self = this;
                    this.$http.post('/deploy/inject?name='+ _self.raw.name + '&namespace='+_self.raw.namespace).then(function (resp) {
                        if (resp.body.code === 0) {
                            _self.sure_inject = false
                            this.$message({
                                message: '注入成功',
                                type: 'success'
                            });
                            this.getList();
                        }else{
                            this.$message({
                                message: resp.body.msg,
                                type: 'error'
                            });
                        }
                    });
                },
                getIstioConfig: function (row) {
                    var _self = this;
                    _self.is_show_istio_config = true;
                    _self.handle_name = row.name;
                    _self.handle_namespace = row.namespace;
                    this.$http.get('/istio_config/get?name='+ _self.handle_name + '&namespace=' + _self.handle_namespace)
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
                sureInject: function (raw) {
                    var _self = this;
                    _self.sure_inject = true
                    _self.raw = raw
                },
                saveIstioConfig: function () {
                    var _self = this;
                    _self.sure_inject = false
                    _self.is_show_istio_config = true;
                    if(_self.istio_config == "") {
                        this.$message({
                            message: "请填写配置信息",
                            type: 'error'
                        });
                        return false
                    }

                    this.$http.post('/istio_config/save',
                        {name: _self.handle_name, namespace: _self.handle_namespace, config: _self.istio_config})
                        .then(function (resp) {
                            if (resp.body.code === 0) {
                                this.$message({
                                    message: "保存成功",
                                    type: 'success'
                                });
                            } else {
                                this.$message({
                                    message: resp.body.msg,
                                    type: 'error'
                                });
                            }
                        });
                },
                delIstioConfig:function () {
                    var _self = this;
                    _self.is_show_istio_config = true;
                    this.$http.post('/istio_config/del', {name: _self.handle_name, namespace: _self.handle_namespace, config: _self.istio_config})
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
                getList: function () {
                    var _self = this;
                    this.$http.get('/deploy/list?' , {params: {page: _self.pagination.currentPage,
                            page_size: _self.pagination.pageSize,
                            name: _self.search_name,
                            namespace: _self.search_namespace}}).then(function (resp) {
                        if (resp.body.code === 0) {
                            _self.tableData = resp.body.data.list;
                            _self.pagination.total = resp.body.data.total;
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
                search:function () {
                    this.getList();
                }
            }
        };
    };
});
