define(function (require, exports, module) {
    module.exports = function (cb) {
        const routes = [
          { path: '/', redirect: '/deploy' },
          { path: '/deploy', component: require("/js/deploy.js")()},
          { path: '/global-istio-config', component: require("/js/global-istio-config.js")()},
        ];
        const router = new VueRouter({
           routes: routes,
        });
        
        cb?cb(router):'';
    }
});
