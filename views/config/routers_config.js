/**
 * 路由配置
 * Created by Sever on 2017/7/3.
 */


define(function (require, exports, module) {
    module.exports = function (cb) {
        
        //路由
        const routes = [
          { path: '/', redirect: '/deploy' },
          { path: '/deploy', component: require("/js/deploy.js")()},
          { path: '/global-istio-config', component: require("/js/global-istio-config.js")()},
          { path: '/service-entry', component: require("/js/service-entry.js")()},
        ];
        
        const router = new VueRouter({
           routes: routes,
        });
        
        
        cb?cb(router):'';
    }
});
