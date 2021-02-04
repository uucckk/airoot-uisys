"use strict";
String.prototype.trim=function(){
	return this.replace(/(^\s*)|(\s*$)/g, "");
}
function trace(){
	var value = "";
	var arr = null;
	for(var i = 0;i<arguments.length;i++){
		value += "[" + i + "]";
		if(typeof(arguments[i]) == 'array'){
			arr = arguments[i];
			for(var j = 0;j<arr.length;j++){
				value += arr[j] + ',';
			}
			value = value.subtring(0,value.length - 1);
			continue;
		}
		value += arguments[i] + " ";
	}
	alert(value);
}

/**
 * asjs API 封装包
 */
var asjs = new function(){
	var __handle__ = [];
	
	/**
	 * 获取域名表示
	 * @param name 编译后的元素id字符串名字
	 */
	this.getDomain = function(name){
		var len = name.length;
		var tmp = null;
		for(var i = len - 1;i>=1;i--){
			tmp = name.substring(0,i);
			if(window[tmp]){
				return tmp;
			}
		}
		return name;
	}
	
	/**
	 * 获取name的原始名字
	 */
	this.getMName = function(name){
		var len = name.length;
		var tmp = null;
		for(var i = len - 1;i>=1;i--){
			tmp = name.substring(0,i);
			//if(window[tmp]){
				//alert(window[tmp] + "");
			//}
			if(window[tmp] && typeof(window[tmp]) == 'object' && (window[tmp] + '').indexOf('[object HTML') == -1){
				return name.substring(i); 
			}
		}
		return name;
	}
	
	
	//Replaced from the original function to leverage the built in methods in
	//JavaScript. Thanks to Robert Kieffer for pointing this one out
	var returnBase = function(number, base){
	 return (number).toString(base).toUpperCase();
	};
	
	 
	
	//pick a random number within a range of numbers
	//int b rand(int a); where 0 <= b <= a
	var rand = function(max){
	 return Math.floor(Math.random() * (max + 1));
	};
	
	/**
	 * uuid 
	 */
	this.uuid = function(len, radix) { 
	    var dg = new Date(1582, 10, 15, 0, 0, 0, 0);
	    var dc = new Date();
	    var t = dc.getTime() - dg.getTime();
	    var tl = getIntegerBits(t,0,31);
	    var tm = getIntegerBits(t,32,47);
	    var thv = getIntegerBits(t,48,59) + '1'; // version 1, security version is 2
	    var csar = getIntegerBits(rand(4095),0,7);
	    var csl = getIntegerBits(rand(4095),0,7);
	    var n = getIntegerBits(rand(8191),0,7) +
	            getIntegerBits(rand(8191),8,15) +
	            getIntegerBits(rand(8191),0,7) +
	            getIntegerBits(rand(8191),8,15) +
	            getIntegerBits(rand(8191),0,15); // this last number is two octets long
	    return tl + tm  + thv  + csar + csl + n;
	};
	
	
	var getIntegerBits = function(val,start,end){
	 var base16 = returnBase(val,16);
	 var quadArray = new Array();
	 var quadString = '';
	 var i = 0;
	 for(i=0;i<base16.length;i++){
	     quadArray.push(base16.substring(i,i+1));   
	 }
	 for(i=Math.floor(start/4);i<=Math.floor(end/4);i++){
	     if(!quadArray[i] || quadArray[i] == '') quadString += '0';
	     else quadString += quadArray[i];
	 }
	 return quadString;
	};
	

	
	
	
	/**
	 * 将元素节点复制一份
	 */
	this.copy = function(obj){
		return obj.cloneNode(true);
	}
	
	
	/**
	 * 加载函数
	 */
	this.load = function(url,compEvt,data,dataType){
		compEvt = compEvt ? compEvt :function(e){};
		data = data ? data : null;
		var ul = new URLLoader();
		var req = new URLRequest(url);
		req.method = URLRequestMethod.POST;
		req.data= data;
		req.dataType = dataType;//json,text,or null
		ul.addEventListener(Event.COMPLETE,compEvt);
		ul.addEventListener(IOErrorEvent.IO_Error,compEvt);
		ul.load(req);
		if(data){data.url = url};
		return data;
	};
	
	/**
	 * 加载函数
	 */
	this.get = function(url,compEvt,data){
		compEvt = compEvt ? compEvt :function(e){};
		data = data ? data : null;
		var ul = new URLLoader();
		var req = new URLRequest(url);
		req.method = URLRequestMethod.GET;
		req.data= data;
		ul.addEventListener(Event.COMPLETE,compEvt);
		ul.addEventListener(IOErrorEvent.IO_Error,compEvt);
		ul.load(req);
		if(data){data.url = url};
		return data;
	};
	
	/**
	 * 返回URLLoader
	 */
	this.url = function(url,compEvt,data){
		compEvt = compEvt ? compEvt :function(e){};
		data = data ? data : null;
		var ul = new URLLoader();
		var req = new URLRequest(url);
		req.method = URLRequestMethod.GET;
		req.data= data;
		ul.addEventListener(Event.COMPLETE,compEvt);
		ul.addEventListener(IOErrorEvent.IO_Error,compEvt);
		ul.load(req);
		return ul;
	};
	
	
	/**
	 * @param 	domain
	 * @url		域名
	 * @compEvt	回调函数
	 * @data	数据
	 */
	this.handle = function(domain,url,compEvt,data){
		compEvt = compEvt ? compEvt :function(e){};
		data = data ? data : null;
		var ul = new URLLoader(domain);
		var req = new URLRequest(url);
		req.method = URLRequestMethod.GET;
		req.data= data;
		ul.addEventListener(Event.COMPLETE,compEvt);
		ul.load(req);
		__handle__.push({domain:domain,ul:ul});
		return ul;
	};
	
	/**
	 * 关闭所有域名下的链接
	 */
	this.closeHandle = function(domain){
		var p = null;
		for(var i = __handle__.length - 1;i>=0;i--){
			p = __handle__[i];
			if(!domain || p.domain == domain){
				p.ul.close();
				__handle__.splice(i,1);
			}
		}
	}
	
	this.send = this.load;
	this.post = this.load;
	
	

}();












var URLRequestMethod = {GET:"get",POST:"post"};
var Event = {COMPLETE:"complete"};
var IOErrorEvent = {IO_Error:"ioError"};
/**
 *
 */
function URLRequest(url){
	//url
	this.URL = url;
	//执行方法，例如使用post方法还是用get方法
	this.method = URLRequestMethod.GET;
	//需要传递的数据，如果method == get ，此方法不生效。
	this.data = null;
}

/**
 * URLLoader 加载数据的请求
 */
function URLLoader(id){
	var req = null;
	var COMP_FUN = null;
	var IOERROR_FUN = null;
	//常量参数
	var READY_STATE_UNINITIALIZED = 0;
	var READY_STATE_LOADING = 1;
	var READY_STATE_LOADED = 2;
	var READY_STATE_INTERACTIVE = 3;
	var READY_STATE_COMPLETE = 4;
	var target = this;
	//最终获得的数据
	this.data = null;
	//加载函数
	this.load = function(urlRequest){
		if(urlRequest instanceof URLRequest){
			req = getRequest();
			if(req){
				req.onreadystatechange = onReadyState;
				var tmp = "";
				if(urlRequest.data != null && urlRequest.dataType != "json" && urlRequest.dataType != "text"){//既不是json，也不是text
					if(typeof(urlRequest.data) == 'string'){
						tmp = urlRequest.data;
						
					}else if(typeof(urlRequest.data) == 'object'){
						for(var p in urlRequest.data){
							var lst = urlRequest.data[p];
							if(lst instanceof Array){
								for(var t = 0;t<lst.length;t++){
									tmp += (p + '=' + encodeURIComponent(lst[t]) + '&');
								}
								continue;
							}
							tmp += (p + '=' + encodeURIComponent(lst) + '&');
						}
					}
				}
				if(urlRequest.method == URLRequestMethod.POST){
					req.open("POST",urlRequest.URL,true);
					req.withCredentials = true;
					//req.setRequestHeader("Content-Length",tmp.length);	
					if(urlRequest.dataType == "json"){
						req.setRequestHeader("Content-Type","application/json;charset=UTF-8");
						tmp = JSON.stringify(urlRequest.data);
					}else if(urlRequest.dataType == "text"){
						req.setRequestHeader("Content-Type","text/plain;charset=UTF-8");
					}else{
						req.setRequestHeader("Content-Type","application/x-www-form-urlencoded;charset=UTF-8");
					}
					
					req.send(tmp);
				}
				else if(urlRequest.method == URLRequestMethod.GET){
					var urlTmp = urlRequest.URL;
					if(urlTmp.indexOf('?') != -1){
						urlTmp += '&' + tmp;
					}else{
						urlTmp += '?' + tmp;
					}
					req.open(urlRequest.method,urlTmp,true);
					req.withCredentials = true;
					req.setRequestHeader("Content-Type","application/x-www-form-urlencoded;charset=UTF-8");
					req.send();
				}
				
			}
		}else{
			throw "URLLoader::load(): The Value isn't URLRequest";
		}
	}
	//时间侦听器 
	this.addEventListener = function(type,listener,useCapture){
		switch(type){
			case Event.COMPLETE :
				COMP_FUN = listener;
			break;
			case IOErrorEvent.IO_Error :
				IOERROR_FUN = listener;
			break;
			default :
				throw "URLLoader::addEventListner(): no this type";
		}
	};//addEventListener
	
	this.close = function(){
		if(req){
			if(req.readyState != READY_STATE_COMPLETE){
				req.abort();
			}
			req = null;
			//console.log("...................close");
		}

	};
	
	/**
	 * 获取连接异常
	 */
	function getRequest(){
		var xRequest = null;
		if(window.XMLHttpRequest){
			xRequest = new XMLHttpRequest();
		}else if(window.ActiveXObject){
			xRequest = new ActiveXObject("Microsoft.XMLHTTP");
		}
		return xRequest;
	}
	
	function onReadyState(){
		var ready = req.readyState;
		if(ready == READY_STATE_COMPLETE){
			if(req.status == "200"){
				target.data = req.responseText;
				if(COMP_FUN != null){
					req = null;
					COMP_FUN({type:Event.COMPLETE,target:target,id:id});
				}
			}else{
				if(IOERROR_FUN != null){
					req = null;
					IOERROR_FUN({type:IOErrorEvent.IO_Error,target:target,id:id});
				}
			}
			
		}
	}
}//URLLoader

//url#hash
window.addEventListener("hashchange", function(e){
	var hash = __Hash__();
	UI.loadModule(document.body,hash["loc"]);
}, false);
//获取Hash
function __Hash__(){
	var hash =  window.location.hash;
	var d = hash.indexOf("#");
	if(d != -1){
		hash = hash.substring(d + 1);
		var arr = hash.split("&");
		var obj = {};
		var t = null;
		for(var i = 0;i<arr.length;i++){
			t = arr[i].split("=");
			obj[t[0]] = t[1];
		}
		return obj;
	}else{
		return null
	}
}

//获取query
function getQueryVariable(q)
{
       var query = window.location.search.substring(1);
       var vars = query.split("&");
       for (var i=0;i<vars.length;i++) {
               var pair = vars[i].split("=");
               if(pair[0] == q){return pair[1];}
       }
       return false;
}
