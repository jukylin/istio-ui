
define(function (require, exports, module) {
    module.exports = function () {
        // console.log('注册组件');
        Vue.component('side_bar', require("/components/side_bar.js")());
    }
});