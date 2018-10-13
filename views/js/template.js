define(function (require,exports,module) {
    module.exports = function() {
        var temp = require('text!/pages/template.html');
        require("https://cdn.jsdelivr.net/npm/vue-resource@1.5.1")
        Vue.http.options.emulateJSON = true;
        return {
            template: temp,
            data: function (){
                return {
                    config: "",
                    tmp_name: "mesh-config",
                }
            },
            mounted : function () {
                var _self = this;
                _self.getTmpName();
                if(_self.tmp_name == "mesh-config"){
                    _self.getMeshConfig();
                }else if(_self.tmp_name == "inject-config"){
                    _self.getInjectConfig();
                }
            },
            methods : {
                getMeshConfig: function () {
                    var _self = this;
                    this.$http.get('/template/getmeshconfig')
                        .then(function (resp) {
                            if (resp.body.code === 0) {
                                _self.config = resp.body.data
                            }else{
                                this.$message({
                                    message: resp.body.msg,
                                    type: 'error'
                                });
                            }
                        });
                },
                getInjectConfig:function () {
                    var _self = this;
                    this.$http.get('/template/getinjectconfig')
                        .then(function (resp) {
                            if (resp.body.code === 0) {
                                _self.config = resp.body.data
                            }else{
                                this.$message({
                                    message: resp.body.msg,
                                    type: 'error'
                                });
                            }
                        });
                },
                getTmpName: function () {
                    var _self = this;
                    if(this.$route.query.type == ""){
                        this.$message({
                            message: "缺少type",
                            type: 'error'
                        });
                        return false;
                    }
                    _self.tmp_name = this.$route.query.type
                    console.log(_self.tmp_name)
                },
                reloadTmpConfig: function () {
                    var _self = this;
                    _self.getTmpName()
                    if(_self.tmp_name == "mesh-config"){
                        _self.getMeshConfig();
                    }else if(this.tmp_name == "inject-config"){
                        _self.getInjectConfig();
                    }
                }
            },
            watch:{
                '$route': 'reloadTmpConfig'
            }
        };
    };
});
