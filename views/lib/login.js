/**
 * Created by Sever on 2017/7/3.
 */
require(['jquery'],
    function ($) {
        var app = new Vue({
            el:"#app",
            data:{
                username:utils.getCookie('login_username') || '',
                password:'',
            },
            beforeCreate: function () {

            },
            mounted: function () {
                var _self = this;
                _self.init();
            },
            methods:{
                init: function(){
                    if(utils.getCookie('admin_token')){
                        location.href="/home.html";
                    }
                },
                //登录
                toLogin: function(){
                    var _self = this;
                    if(!_self.username){
                        utils.alert('账号不能为空');
                        return;
                    }
                    if(!_self.password){
                        utils.alert('密码不能为空');
                        return;
                    }

                    utils.commonAjax({
                        url:'privilege/front/login',
                        data:{
                            username:_self.username,
                            password:_self.password,
                        },
                        type:'POST',
                        success:function(json){
                            if(json.code==0){
                                if(json.data.first==1){
                                    
                                }else{
                                    utils.setCookie('login_username',_self.username,1);
                                    location.href="/home.html";
                                }
                            }
                            else{
                                
                                utils.alert(json.msg);
                            }
                        },
                    });
                },
                clearData: function () {
                    this.username = "";
                    this.password = "";
                },
            }
        });
    }
);
