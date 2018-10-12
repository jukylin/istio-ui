/**
 * Created by Sever on 2017/7/3.
 */
define(function (require, exports, module) {
    module.exports = function() {
        // console.log('加载组件-侧边栏');
        var temp = require('text!/components/side_bar.html');
        return {
            template: temp,
            data: function () {
                return {
                    defaultActive: "/deploy",
                    isCollapse: true,
                }
            },
            mounted: function () {
                var _self = this;

                _self.initSideBar();
                
                //监听路由变化
                this.$router.afterEach(function () {
                    _self.initSideBar();
                    _self.defaultActive = _self.$route.path;
                })
            },
            methods: {
                initSideBar: function(){
                    var _self = this;
                    //console.log(_self.$route.path);
                    _self.defaultActive = _self.$route.path;
                },
                
                handleOpen: function(key, keyPath) {
                    //console.log(key, keyPath);
                },
                handleClose: function(key, keyPath) {
                    //console.log(key, keyPath);
                },
                menuSelect: function (key,keyPath) {
                    if (key) {
                        //console.log(key, keyPath);
                        // location.href='#'+key;
                        this.$router.push({
                            path: key
                        });
                    }
                    
                }
            },
        };
    };
});
