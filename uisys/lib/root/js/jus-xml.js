

/**
 * 将字符串数据进行整理,变成XML的可读类
 */
function XML(data){
	var _root = this;
	var xml = null;
	var arr = new Array();
	var btype 
	if(data != null){
		if(data instanceof Array){
			arr = data;
		}else{
			if(window.DOMParser){
				xml = (new DOMParser()).parseFromString(data, "text/xml").childNodes; 
			}else{
				xml = new ActiveXObject("Microsoft.XMLDOM"); 
				xml.loadXML(data);
				xml  = xml.childNodes;
			}
			for(var n = 0;n<xml.length;n++){
				arr.push(xml[n]);
			}
		}
	}
	
	
	
	
	/**
	 * 获取元素的属性
	 * @param name	属性名称
	 * @param value	属性值，可以不填写
	 */
	this.qname = function(name,value){//setAttribute,getAttribute
		if(!name){
			return arr[0].attributes;
		}
		var outStr = "";
		for(var i = 0;i<arr.length;i++){
			if(value != null){
				arr[i].setAttribute(name,value);
			}
			outStr += arr[i].getAttribute(name) + ",";
		}
		return outStr.substr(0,outStr.length - 1);
	}
	
	
	
	
	/**
	 * 长度
	 */
	this.length = function (){
		return arr.length;
	}
	
	
	/**
	 * 查找节点内容
	 */
	this.child = function (nodeName){
		
		if(nodeName && nodeName.charAt(0) == '@'){
			return this.qname(nodeName.substring(1));
		}else if(nodeName && nodeName.charAt(0) == '['){
			nodeName = nodeName.substr(1,nodeName.length - 2);
			var values = nodeName.split(".");
			var p = this.child(values[0]);
			var t = null;
			for(var n = 1;n<values.length;n++){
				t = values[n]; 
				if(t.charAt(t.length - 1) == ')'){
					p = p.at(t.substring(3,t.length - 1));
					continue;
				}
				p = p.child(t);
			}
			return p;
		}
		var a = new Array();
		var ch = null;
		for(var i = 0;i<arr.length;i++){
			ch = arr[i].childNodes;
			for(var j = 0;j<ch.length;j++){
				if(nodeName == null || ch[j].nodeName == nodeName){
					a.push(ch[j]);
				}
			}
		}
		return new XML(a);
	}
	
	
	/**
	 * 制定具体位置
	 */
	this.at = function (pos){
		var a = new Array();
		a.push(arr[pos]);
		return new XML(a);
	}
	
	/**
	 * 节点赋值
	 */
	this.node = function(value){
		if(value){
			for(var i = 0;i<arr.length;i++){
				if(typeof(value) == "string"){
					arr[i].parentNode.replaceChild(new XML(value)._nodeValue(0),arr[i]);
				}
				
			}
		}
	}
	
	this._nodeValue = function(p){
		return arr[p];
	}
	
	/**
	 * 获取XML格式内容
	 */
	this.toXMLString = function(){
		var outStr = "";
		var i = 0;
		switch(browser){
			case "ie5+":
				for(i = 0;i<arr.length;i++){
					outStr += arr[i].xml;
				}
			break;
			case "other" :
				for(i = 0;i<arr.length;i++){
					outStr += (new XMLSerializer()).serializeToString(arr[i]); 
				}
			break;
		}
		
		return outStr;
	}
	
	/**
	 * 获取JSON数据
	 */
	this.toJSON = function(){
		var obj = {};
		var arr = null;
		for(var i = 0;i<this.length();i++){
			arr = this.at(i).child('@');
			for(var j = 0;j<arr.length;j++){
				obj[arr[j].name] = arr[j].value;
			}
		}
		return obj;
	};
	
	/**
	 * 获取JSON数组
	 */
	this.toJSONArray = function(){
		var list = [];
		var obj = null;
		var arr = null;
		for(var i = 0;i<this.length();i++){
			arr = this.at(i).child('@');
			obj = {};
			for(var j = 0;j<arr.length;j++){
				obj["@"+arr[j].name] = arr[j].value;
			}
			var child = this.at(i).child();
			for(j = 0;j<child.length();j++){
				obj[child.at(j).getName()] = child.at(j).toString();
			}
			
			list.push(obj);
		}
		return list;
		
	};
	
	this.appendChild = function (data){
		var child = new XML("<response>" + data + "</response>").child();
		var len = child.length();
		if(arr.length == 1){
			for(var i = 0;i<len;i++){
				arr[0].appendChild(child._nodeValue(i));
			}
		}
	}
	
	/**
	 * 删除指定元素
	 * @param nodeName
	 * @return
	 */
	this.removeChild = function(nodeName){
		var len = 0;
		
		var pos = null;
		var child = null;
		
		for(var i = 0;i<arr.length;i++){
			pos = arr[i];
			if(nodeName instanceof XML){
				child = nodeName;
			}else{
				child = new XML([arr[i]]).child(nodeName);
			}
			
			len = child.length();
			for(var j = 0;j<len;j++){
				pos.removeChild(child._nodeValue(j));
			}
		}
	}
	
	
	this.getName = function(){
		return arr.length>0 ? arr[0].nodeName : null;
	}
	
	
	
	/**
	 * 重写toString方法
	 */
	this.toString = function(){
		var outStr = "";
		for(var i = 0;i<arr.length;i++){
			if(arr[i].childNodes.length != 0){
				outStr += arr[i].childNodes[0].wholeText;
			}
		}
		return outStr;
	}
	
}//XML