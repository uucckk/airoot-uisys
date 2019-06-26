/**
 * 模板控制
 * @author sunxy
 * @type 功能包
 */
class Template{
	var node = null;
	var a = /\{.+?\}/g;
	var map = {};
	var attrLst = [];
	var _data = {};
	var _pdata = null;
	var _onChange = null;
	var watchLst = {};//查看关注数据项目
	public var filter = {};
	
	
	function init(node:HTML){
		if(!(node instanceof HTMLElement)){
			node = node.dom;
		}
		defer(node,@this);
		read(this.node = node);
		initListener();
	}
	
	
	/**
	 * 监听输入元素
	 */
	private function initListener(){
		node.addEventListener("input",listener);
	}
	
	private function removeListener(){
		node.removeEventListener("input",listener);
	}
	
	private function listener(e){
		var target = e.target;
		var p = null;
		var index = null;
		f:for(var k in map){
			p = map[k];
			for(var i in p){
				if(p[i].node == target){
					index = k.split(".");
					break f;
				}
			}
		}
		if(index && dataContext){
			p = dataContext;
			for(var i = 0;i<index.length - 1;i++){
				p = p[index[i]];
			}
			var name = target.getAttribute(":value");
			if(name){				
				var n = name.lastIndexOf(".");
				name = n == -1 ? name : name.substring(n + 1);
				p[name] = e.target.value;
			}
		}
	}
	
	/**
	 * 数据提供
	 */
	public function set dataContext(value:Object){
		if(typeof value == "object"){
			if(_pdata != value){
				__TPL_MAP_PUSH__(value,{target:_data,domain:@this});
				_pdata = value;
			}
			
			var p = null;
			for(var i in _data){
				p = value[i]
				if(p != undefined && p != null){
					_data[i] = p;
				}
				
			}
		}else{
			if(window.console){
				console.error("dataContext input isn't object");
			}else{
				alert(value);
			}
		}
		
	}
	
	public function get dataContext():Object{
		return _data;
	}
	
	/**
	 * 监听数据变化
	 */
	public function watch(value,listener){
		if(!watchLst[value]){
			if(listener){
				watchLst[value] = [listener];
			}
			return;
		}
		var lst = watchLst[value];
		for(var i = 0;i<lst.length;i++){
			if(listener == lst[i]){
				return;
			}
		}
		lst.push(listener);
	}
	
	
	/**
	 * 通知
	 */
	private function notify(index,newName,oldName){
		var name = "";
		for(var j = 0;j<index.length - 1;j++){
			name += index[j] + ".";
		}
		name += index[index.length - 1];
		var p = null;
		for(var k in watchLst){
			if(k == name){
				p = watchLst[k];
				for(var i = 0;i<p.length;i++){
					p[i](name,newName,oldName);
				}
			}
		}
	}
	
	/**
	 * 插入模版
	 */
	public function innerTemplate(id,dom){
		if(typeof dom == "string"){
			id.innerHTML = dom;
		}else if(dom instanceof HTMLElement){
			id.appendChild(dom);
		}else{
			alert("未识别元素");
			return;
		}
		read(id);
		var value = _pdata ? _pdata :_data;
		var p = null;
		for(var i in _data){
			p = value[i]
			if(p != undefined && p != null){
				_data[i] = p;
			}
			
		}
		
		
	}
	
	
	/**
	 * 数据变更通知
	 */
	public function set onDataChange(value:Function){
		_onChange = value;
	}
	
	private function read(node:HTML){
		var node = node.childNodes;
		var p = null;
		var m = null;
		var start = 0;
		var name = null;
		var txt = 0;
		for(var i = 0;i<node.length;i++){
			p = node[i];
			if(p instanceof Text){
				if(p.length>2 && p.nodeValue.indexOf("{") != -1){
					while((m = a.exec(p.nodeValue)) != null){
						name = p.nodeValue.substring(m.index + 1,a.lastIndex - 1).trim();
						if(start != m.index){
							txt = document.createTextNode(p.nodeValue.substring(start,m.index));//插入前面的值
							p.parentNode.insertBefore(txt,p);
							i ++;
						}
						txt = document.createTextNode("{" + name +  "}");//插入当前匹配项目
						pushText(fx(name),txt);
						p.parentNode.insertBefore(txt,p);
						start = a.lastIndex;
						i ++;
					};
					txt = document.createTextNode(p.nodeValue.substring(start));//插入前面的值
					p.parentNode.insertBefore(txt,p);
					p.parentNode.removeChild(p);
					
					
				}
			}else{
				var atl = p.attributes;
				if(atl){
					var nm = null;
					for(var t = 0;t<atl.length;t++){
						nm = atl[t];
						if(nm.name.charAt(0) == ":"){
							console.log(">>",nm,nm.value);
							pushText(fx(nm.value),p,nm.name.substring(1));
						}
					}
					read(p);
				}
				
			}
		}
	}
	
	private function pushText(obj,node,attr){
		if(!map[obj.stat]){
			map[obj.stat] = [];
		}
		map[obj.stat].push({node:node,attr:attr,filter:obj.filter});
		//先把所有属性转化下
		var arr = obj.arr;
		var c = _data;
		var p = null;
		for(var i = 0;i<arr.length;i++){
			p = arr[i];
			if(c[p]){
				if(i == arr.length - 1 && obj.value != null){
					c[p] = obj.value;
				}else if(!(typeof c[p] == "object")){
					c[p] = {};
				}
			}else{
				if(!c.hasOwnProperty(p)){
					setAttribute(c,p,obj);
				}
				if(i == arr.length - 1){
					c[p] = obj.value;
				}else{
					c[p] = {};
				}
			}
			c = c[p];
		}
		
	}
	
	/**
	 * obj
	 */
	private function fx(value){
		var stat = null;
		var arr = null;
		var val = null;
		var filter = [];
		/*
		var i = value.indexOf("=");
		if(i == -1){
			stat = value;
			arr = value.split(".");
		}else{
			stat = value.substring(0,i).trim();
			var b = value.substring(i+1).trim();
			arr = stat.split(".");
			val = eval(b);
		}
		return {stat:stat,arr:arr,value:val};
		*/
		var sb = "";
		var p = null;
		var E = false;
		for(var i = 0;i<value.length;i++){
			p = value.charAt(i);
			if(p == "'" || p == '"'){
				if(E){
					var b = readString(i,value);
					val = eval(b.value);
					i = b.index;
				}else{
					throw new Error("read: " + value + "hav error.");
				}
				continue;
			}
			if(p == '='){
				E = true;
				continue;
			}
			if(p == '|'){
				readFilter(i,value,filter);
				break;
			}
			sb += p;
		}
		stat = sb.trim();
		//console.log({stat:stat,arr:stat.split("."),value:val,filter:filter});
		return {stat:stat,arr:stat.split("."),value:val,filter:filter};
	}
	
	private function readString(index,value){
		var p = null;
		var c = value.charAt(index);
		var str = c;
		var f = false;
		for(var i = index + 1;i<value.length;i++){
			p = value.charAt(i);
			str += p;
			if(p == c && f == false){
				return {index:i,value:str};
			}
			if(p == '\\'){
				f = !f;
			}
		}
		throw new Error("read: " + value + "hav error.");
	}
	
	private function readFilter(index,value,filter){
		var str = "";
		for(var i = index + 1;i<value.length;i++){
			p = value.charAt(i);
			if(p == '|'){
				filter.push(str.trim());
				str = "";
				continue;;
			}
			str += p;
			
		}
		filter.push(str.trim());
		return {};
	}
	
	private function setAttribute(obj,name,info){
		var index = info.arr;
		var _value = null;
		var V = null;
		Object.defineProperty(obj,name,{
			set:function(value){
				var tValue = _value;
				if(_value && typeof value == "object"){
					if(V != value){
						__TPL_MAP_REMOVE__(V,_value);
						__TPL_MAP_PUSH__(value,{target:_value,domain:@this});
						V = value;
					}
					for(var k in _value){
						_value[k] = value[k];
					}
					return;
				}
				if(_value != value){
					var p = _pdata ? _pdata : _data;
					for(var i = 0;i<index.length - 1;i++){
						p = p[index[i]];
						if(!p){
							break;
						}
					}
					_value = value;
					if(p){
						
						if(_pdata){
							notify(index,value,tValue);
							p[name] = value;
						}
						update(p);
						
					}
					
					render(info.stat,value);
				}
				
				
			},
			get:function(){
				return _value;
			}
		,enumerable:true});
	}
	
	//渲染元素
	private function render(stat,value){
		var arr = map[stat];
		var f = null;
		var p = null;
		var k = null;
		var v = null;
		for(var i = 0;i<arr.length;i++){
			p = arr[i];
			if(p.filter){
				f = p.filter;
				v = value;
				for(var j = 0;j<f.length;j++){
					k = filter[f[j]];
					if(k){
						if(p){
							v = k(v);
						}
					}else{
						console.error("Template Filter:[" + f[j] + "] isn't exist.");
					}
				}
			}else{
				v = value;
			}
			if(p.attr){//HTMLElement
				if(p.attr.charAt(0) == "+"){
					p.node.setAttribute(p.attr.substring(1),v);
				}else{
					if(p.node.type){
						if(p.node.type == "radio"){
							if(p.node.value == v){
								p.node.checked = true;
							}
						}else{
							p.node[p.attr] = v;
						}
					}else{
						p.node[p.attr] = v;
					}
					
				}
				
			}else{//Text
				p.node.nodeValue = v;
			}
		}
	}	
	
	public function destroy(){
		removeListener();
		var p = null;
		for(var i = 0;i<tplMap.length;i++){
			p = tplMap[i];
			if(p.domain == @this){
				p.domain = null;
			}
		}
	}
	
}