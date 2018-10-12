/**
 * js工具
 * Created by Sever on 2017/7/3.
 */
Date.prototype.Format = function (fmt) { //author: meizz
    var o = {
        "M+": this.getMonth() + 1, //月份
        "d+": this.getDate(), //日
        "h+": this.getHours(), //小时
        "m+": this.getMinutes(), //分
        "s+": this.getSeconds(), //秒
        "q+": Math.floor((this.getMonth() + 3) / 3), //季度
        "S": this.getMilliseconds() //毫秒
    };
    if (/(y+)/.test(fmt))
        fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
    for (var k in o)
        if (new RegExp("(" + k + ")").test(fmt))
            fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
    return fmt;
};
/**
 * js中更改日期
 * @param part y年， m月， d日， h小时， n分钟，s秒
 * @param value
 */
Date.prototype.add = function (part, value) {
    value *= 1;
    if (isNaN(value)) {
        value = 0;
    }
    switch (part) {
        case "y":
            this.setFullYear(this.getFullYear() + value);
            break;
        case "m":
            this.setMonth(this.getMonth() + value);
            break;
        case "d":
            this.setDate(this.getDate() + value);
            break;
        case "h":
            this.setHours(this.getHours() + value);
            break;
        case "n":
            this.setMinutes(this.getMinutes() + value);
            break;
        case "s":
            this.setSeconds(this.getSeconds() + value);
            break;
    }
    return this;
};
Date.prototype.getLastDay = function (){
    var year = this.getFullYear();
    var month = this.getMonth()+1;
    // var firstdate = year + '-' + month + '-01';
    var day = new Date(year,month,0);
    var lastdate = day.getDate();//获取当月最后一天日期
    return lastdate;
};

var utils = {
    //是否是开发环境
    isDev: function () {
        //return true;
        return false;
    },
    //获取url查询字符串
    queryString: function (key) {
        return (document.location.search.match(new RegExp("(?:^\\?|&)" + key + "=(.*?)(?=&|$)")) || ['', null])[1];
    },
    queryHash: function (key) {
        return (document.location.hash.match(new RegExp("(?:\\?|&)" + key + "=(.*?)(?=&|$)")) || ['', null])[1];
    },
    //获取url锚点字符串
    hashString: function () {
        return document.location.hash.replace("#","");
    },
    //通用ajax
    commonAjax: function (obj) {
        var loadingInstance = null;
        // obj.url = utils.isDev()?'':obj.url;
        // obj.type = utils.isDev()?'GET':obj.type;
        obj.beforeSend = obj.beforeSend || function () {
            loadingInstance = ELEMENT.Loading.service({fullscreen: true});
        };
        obj.complete = obj.complete || function () {
            loadingInstance.close();
        };
        obj.error = obj.error || function () {
            utils.alert('网络错误，请重试');
        };

        /*var success = function(json){
            //统一对code进行异常处理
            if(json.code==0){
                obj.success.call(this,json);
            }else{
                utils.alert(json.msg);
            }
        }
        //深度复制
        var clone_obj = utils.deepClone(obj);
        clone_obj.success = success;
        $.ajax(clone_obj);*/

        $.ajax(obj);
    },
    //警告
    alert: function (obj,cb) {
        var content = obj.content || String(obj) || "";
        var btnText = obj.btnText || "确定";
        var title = obj.title || "提示";

        ELEMENT.MessageBox.alert(utils.html2vnode(content),title,{
            // closeOnClickModal:false,//点击遮罩关闭窗口
            confirmButtonText: btnText,
        }).then(function () {
            cb ? cb() : '';
        }).catch(function () {
            // utils.alert(obj,cb);
        });
    },
    //确认
    confirm: function (obj, cb) {
        var content = obj.content || String(obj) || "";
        var btnText = obj.btnText || "确定";
        var cancelText = obj.cancelText || "取消";
        var title = obj.title || "确认";

        ELEMENT.MessageBox.confirm(utils.html2vnode(content), title, {
            confirmButtonText: btnText,
            cancelButtonText: cancelText,
        }).then(function () {
            cb ? cb(true) : '';
        }).catch(function () {
            cb ? cb(false) : '';
        });
    },
    //设置cookie（名称，值，时间[单位：天]）
    setCookie: function (name, value, time) {
        var oDate = new Date();
        oDate.setDate(oDate.getDate() + time);
        document.cookie = name + "=" + encodeURIComponent(value) + ";expires=" + oDate;
    },
    //获取cookie
    getCookie: function (name) {
        var arr = document.cookie.split("; ");
        for (var i = 0; i < arr.length; i++) {
            var arr2 = arr[i].split("=");
            if (arr2[0] == name) {
                return decodeURIComponent(arr2[1]);
            }
        }
        return "";
    },
    //移除cookie
    removeCookie: function (name) {
        this.setCookie(name, "", -1);
    },
    //json对象转数组
    jsonObjToArray:function(obj){
        var arr = [];
        for(var key in obj){
            arr.push(obj[key]);
        }
        return arr;
    },
    //深度复制json
    deepClone: function(obj){
        var result={};
        for(key in obj){
            result[key]=obj[key];
        }
        return result;
    },
    //格式化时间（毫秒）
    formatTime: function (ms) {
        var ss = 1000;
        var mi = ss * 60;
        var hh = mi * 60;
        var dd = hh * 24;

        var day = parseInt(ms / dd);
        var hour = parseInt((ms - day * dd) / hh);
        var minute = parseInt((ms - day * dd - hour * hh) / mi);
        var second = parseInt((ms - day * dd - hour * hh - minute * mi) / ss);
        var milliSecond = parseInt(ms - day * dd - hour * hh - minute * mi - second * ss);

        var result = {
            d: day,
            h: hour,
            mi: minute,
            s: second,
            ms: milliSecond,
            toString:function(){
                var str = "";
                str += this.d?this.d+'天':'';
                str += this.h?this.h+'小时':'';
                str += this.mi?this.mi+'分钟':'';
                str += this.s?this.s+'秒':'';
                str += this.ms?this.ms+'毫秒':'';
                return str;
            }
        };

        return result;

    },
    //格式化时间_补零（毫秒）
    formatTime_zero: function (ms) {
        var ss = 1000;
        var mi = ss * 60;
        var hh = mi * 60;
        var dd = hh * 24;

        var day = utils.prefixInteger(ms / dd,2);
        var hour = utils.prefixInteger((ms - day * dd) / hh,2);
        var minute = utils.prefixInteger((ms - day * dd - hour * hh) / mi,2);
        var second = utils.prefixInteger((ms - day * dd - hour * hh - minute * mi) / ss,2);
        var milliSecond = utils.prefixInteger(ms - day * dd - hour * hh - minute * mi - second * ss,3);

        var result = {
            d: day,
            h: hour,
            mi: minute,
            s: second,
            ms: milliSecond,
            toString:function(){
                var str = "";
                str += this.d!='00'?this.d+'天':'';
                str += this.h!='00'?this.h+'小时':'';
                str += this.mi!='00'?this.mi+'分钟':'';
                str += this.s!='00'?this.s+'秒':'';
                str += this.ms!='000'?this.ms+'毫秒':'';
                return str;
            }
        };
        return result;
    },
    //补零
    prefixInteger:function (num, n) {
        num = parseInt(num);
        if(num.toString().length>n){
            return num.toString();
        }else{
            return (Array(n).join(0) + num).slice(-n);
        }
    },
    //html转vnode（vnode参数传入vue创建方法this.$createElement）
    html2vnode:function(html){
        // console.log(Vue,window.$createElement,Vue);
        var vnode = new Vue().$createElement;
        this.vnode_result = null;
        //生成
        var objE = document.createElement("div");
        objE.innerHTML = html;
        this.traDom2Vnode(objE.childNodes);
        // console.log(this.vnode_result);
        return this.vnode_result;
    },
    //遍历dom生成vnode
    traDom2Vnode:function(node,parent){
        var vnode = new Vue().$createElement;
        if(!node){
            return;
        }
        if(!parent){
            parent = vnode('p', null,[]);
            this.vnode_result = parent;
        }
        node.forEach(function(dom){
            // console.log(dom,dom.nodeName,dom.nodeValue,dom.innerHTML,dom.attributes,dom.childNodes);
            var style_obj = dom.attributes?{}:null;
            if(style_obj){
                for(var key=0;key<dom.attributes.length;key++){
                    // console.log(key,dom.attributes[key]);
                    // console.log(dom.attributes[key]);
                    style_obj[dom.attributes[key].nodeName] = dom.attributes[key].nodeValue;
                }
                // console.log(style_obj);
            }

            if(dom.nodeName!=='#text'){
                // console.log(dom.childNodes);
                if(dom.childNodes.length>0){
                    var child = vnode(dom.nodeName,style_obj,[]);
                    parent.children.push(child);
                    utils.traDom2Vnode(dom.childNodes,child);
                }else{
                    parent.children.push(vnode(dom.nodeName, style_obj, dom.innerHTML));
                }
            }else{
                parent.children.push(vnode('span', null, dom.nodeValue));
            }
        })
    },
    
    //返回当前 url
    getCurrentHtml:function(){
        var current_link = location.href;
        var index = current_link.lastIndexOf('/');
        var index_file = current_link.lastIndexOf('.');
        var current_html = current_link.substr(index+1,30);
        model_js = current_link.substr(index+1,(index_file-index-1));
        
        return [current_html,model_js+".js"];
    },
    
    //返回拼装url
    toQueryString: function (obj) {
        var ret = []; 
        for(var key in obj){ 
            key = encodeURIComponent(key); 
            var values = obj[key]; 
            if(values && values.constructor == Array){//数组 
                var queryValues = []; 
                for (var i = 0, len = values.length, value; i < len; i++) { 
                    value = values[i]; 
                    queryValues.push(this.toQueryPair(key, value)); 
                } 
                ret = ret.concat(queryValues); 
            }else{ //字符串 
                ret.push(this.toQueryPair(key, values)); 
            } 
        } 
        return ret.join('&'); 
    },
    toQueryPair: function (key, value) { 
        if (typeof value == 'undefined'){ 
            return key; 
        } 
        return key + '=' + encodeURIComponent(value === null ? '' : String(value)); 
    },
};