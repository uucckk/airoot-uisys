/**
 * DOM绑定控制
 * @author sunxy
 * @type 功能包
 */
import util.BindData;
class DOMBinding{
	private var node = null;
	private static var a = /\{.+?\}/g;
	private var map = {};
	private var aList = [];
	private var uList = [];//没有被渲染的
	private var _bindData = null;
	private var _set = {};
	//过滤器对象
	public var filter = {};
	
	
	function init(node:HTML){
		if(!node){
			alert("DOMBind need a HTMLNode");
		}
		if(!(node instanceof HTMLElement)){
			node = node.dom;
		}
		defer(node,@this);
		read(this.node = node,aList);
		initListener();
	}
	
	/**
	 * 设置数据
	 * 数据值为对象
	 */
	set data(value){
		if(@type(value) != 'util.BindData'){
			value = new BindData(value);
			
		}
		_set['#'] = value;
		render(aList);
		rebind();
	}
	
	/**
	 * 获取默认绑定数据对象的值
	 */
	function get data(){
		return _set['#']? _set['#'].data : null;
	}
	
	private func render(aList,name){
		var p = null;
		var q = null;
		var a = null;
		var v = null;
		var t = null;
		var fo = null;
		var pa = null;
		var pb = null;
		var h = -1;
		for(var i = 0;i<aList.length;i++){//{属性}列表
			p = aList[i];
			q = p.attr;
			pa = [];
			pb = [];
			for(var j = 0;j<q.length;j++){// | 列表
				a = q[j];
				for(var k = 0;k<a.length;k++){//多少个参属性或参数
					v = a[k];
					if(j == 0){
						switch(v.t){
							case 0 :
								let m = gAttr(v.v);
								switch(h){
									case 0 : //+
										pa[pa.length - 1] += m
										//alert(a[k - 1])
									break;
									case 1 : //-
										pa[pa.length - 1] -= m
									break;
									case 2 : //*
										pa[pa.length - 1] *= m
									break;
									case 3 : //\/
										pa[pa.length - 1] /= m
									break;
									case 4 : //%
										pa[pa.length - 1] %= m
									break;
									case 5 : //&
										pa[pa.length - 1] = pa[pa.length - 1] & m
									break;
									case 6 : //|
										pa[pa.length - 1] %= pa[pa.length - 1] | m
									break;
									case 7 : //！
										pa.push(!m);
									break;
								}
								if(h == -1){
									pa.push(m);
								}
								h = -1;
								break;
							case -1:
							case 1 :
								switch(h){
									case 0 : //+
										pa[pa.length - 1] += v.v
										//alert(a[k - 1])
									break;
									case 1 : //-
										pa[pa.length - 1] -= v.v
									break;
									case 2 : //*
										pa[pa.length - 1] *= v.v
									break;
									case 3 : //\/
										pa[pa.length - 1] /= v.v
									break;
									case 4 : //%
										pa[pa.length - 1] %= v.v
									break;
									case 5 : //&
										pa[pa.length - 1] = pa[pa.length - 1] & v.v
									break;
									case 6 : //|
										pa[pa.length - 1] %= pa[pa.length - 1] | v.v
									break;
									case 7 : //！
										pa.push(!v.v);
									break;
								}
								if(h == -1){
									pa.push(v.v);
								}
								h = -1;
								
								break;
							case 2 ://计算符
								h = v.v
						}
					}else{
						if(v.t == 0 && v.p == 0){//函数名
							fo = v.v;
						}else{//参数
							pb.push(v.t == 0 ? gAttr(v.v) : v.v);
						}
					}
				}
				if(fo){
					if(filter[fo]){
						t = filter[fo].apply(this,pb.concat(pa));
						pb = [t];
						pa = [];
					}else{
						console.error("AIroot-UISYS: util.DomBinding: filter[" + fo + "] isn't exist.");
					}
					fo = null;
				}else{
					t = pa[0];
				}	
			}
			pa = null;
			if(p.node.constructor == Text){
				p.node.nodeValue = t;
			}else{
				if(@global[p.node.id]){
					if(@global[p.node.id][p.name]){
						@global[p.node.id][p.name] = t;
					}else{
						let n = @global[p.node.id];
						for(let k in n){
							if(k.toLowerCase() == p.name){
								n[k] = t;
							}
						}
					}							
				}else{
					if(t != null){
						if(t.constructor == Object){
						
							let n = p.node[p.name];
							for(let k in t){
								n[k] = t[k];
							}
						}else if(t.constructor == Array){
							if(p.name == "class"){
								p.node.classList = "";
								let n = p.node.classList;
								for(let i = 0;i< t.length;i++){
									n.add(t[k]);
								}
							}else{
								p.node[p.name] = t.join(" ");
							}
							
						}else{
							p.node[p.name] = t;
						}
					}
				}
			};
		}
	}
	
	
	/**
	 * 记录 binding 数据项
	 */
	private func rebind(){
		map = {};
		var t = null;
		var p = null;
		var f = null;
		var a,j,n;
		var o;
		for(var i = 0;i<aList.length;i++){
			p = aList[i];
			f = p.father;
			s:for(var k = 0;k<f.length;k++){
				o = gAttrD(f[k]);
				if(o.v == undefined){//说明访问的属性没有值
					continue;
				}
				if(!map[o.d]){
					map[o.d] = {};
				}
				t = map[o.d];
				n = o.v._u$ + "." + o.a;
				if(!t[n]){
					t[n] = [];
				}else{
					a = t[n];
					for(j = 0;i<a.length;i++){
						if(a[i] == p){
							continue s;
						}
					}
				}
				
				t[n].push(p);
			}
		}
		
	}
	
	
	/**
	 * @param n		数据域名
	 * @param cls	操作类型
	 * @param obj	数据对象
	 * @param c		新数据ID
	 * @param a		更改的属性名称
	 * @param r		是否要刷新绑定
	 */
	private func _chg(n,cls,obj,c,a,r){
		var v = obj._u$;
		var d;
		if(v>-1 && map[n]){
			d = v + "." + a;
			var l = map[n][d];
			if(l){
				if(c){
					var j = l[0].father[0];
					map[n][gAttr( j , "_u$")] = map[n][d];
					delete map[n][d];
				}		
				render(l,a);
			}
		}
		if(r){
			render(aList);
			rebind();
		}
	}
	
	//通过绑定方式设置的数据
	set bindData(value){
		if(value){
			_set['#'] = value;
			value.\$addListener(_chg,'#');
			data = value;	
		}
	}
	
	get bindData(value){
		return _set['#']
	}
	
	
	/**
	 * @param name
	 * @param value 
	 * 增加绑定数据
	 */
	func addBindData(name,value){
		name = "#" + name;
		_set[name] = value;
		value.\$addListener(_chg,name);
		render(aList);
		rebind();
	}
	
	func addData(name,value){
		addBindData(name,new BindData(value));
	}
	
	
	func removeBindData(name){
		delete _set["#" + name];
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
		var p = null;
		var a = null;
		var data = null;
		for(var i = 0;i<aList.length;i++){
			p = aList[i];
			if(p.node == e.target){
				a = p.attr;
				sAttr(a[0][0].v, p.node.value);
				let j = null;
				for(let i = 0;i<p.father.length;i++){
					j = p.father[i];
					if(j[0][0] == '#'){
						data = _set[j[0]]
					}else{
						data = _set['#'];
					}
					data.event('c',gAttrD(j).v,null,j[j.length - 1]);
				}
			}
		}
	}
	
	
	
	
	
	private static function read(node:HTML,aList){
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
					start = 0;
					while((m = a.exec(p.nodeValue)) != null){
						name = p.nodeValue.substring(m.index + 1,a.lastIndex - 1).trim();
						if(start != m.index){
							txt = document.createTextNode(p.nodeValue.substring(start,m.index));//插入前面的值
							p.parentNode.insertBefore(txt,p);
							i ++;
						}
						
						txt = document.createTextNode("{" + name +  "}");//插入当前匹配项目
						aList.push(fx(name,txt));
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
						if(nm.value.charAt(0) == "{"){
							var n = nm.value.substring(1,nm.value.length - 1);
							aList.push(fx(n,p,nm.name));
							nm.value = "";
						}
					}
					read(p,aList);
				}
				
			}
		}
	}
	
	
	/**
	 * 分析{}内的内容
	 */
	private static function fx(v,node,name){
		var attr = [];
		var a = [];
		var c;
		var t = null;
		var s = "";
		var l = 0;
		var z = false;
		var father = [];
		var cnt = null;//当前处理的属性值
		for(var i = 0;i<v.length;i++){
			c = v[i];
			if(c == '"' || c == "'"){
				t = readStr(v,i);
				if(cnt){
					cnt.push(t.v);
				}else{
					a.push({v:t.v,t:1,p:l});
				}
				i = t.i;
				continue;
			}else if(c == ' ' || c == ',' || c == '(' || c == ')'){
				if(s.length > 0){
					cnt = cnt ? cnt : [];
					
					cnt.push(s)
					if(cnt.length == 1 && !isNaN(cnt[0])){//说明数字
						a.push({v:+cnt[0],t:-1,p:l});//v 名称，t 是否为字符串, p：是否为参数
					}else{
						if(attr.length>0 && a.length == 0){//说明是过滤函数了
							a.push({v:cnt[0],t:0,p:l});
						}else{
							father.push(cnt);
							a.push({v:cnt,t:0,p:l});//v 名称，t 是否为字符串, p：是否为参数
						}
						
					}
					s = "";
					cnt = null;
					
				}
				
				if(c == '('){
					l ++;
				}else if (c == ')'){
					l --;
				}
				
				continue;
			}else if(c == '.' || c == '[' || c == ']'){
				if(s.length != 0){
					cnt = cnt ? cnt : [];
					cnt.push(s);
					s = "";
				}
				
				if(c == '['){
					z = true
				}else if(c == ']'){
					z = false;
				}
				
				continue;
			}else if(l == 0 && c == '|'){
				
				if(s.length>0){
					cnt = cnt ? cnt : [];
					cnt.push(s);
					if(cnt.length == 1 && !isNaN(cnt[0])){//说明数字
						a.push({v:+cnt[0],t:-1,p:l});//v 名称，t 是否为字符串, p：是否为参数
					}else{
						a.push({v:cnt,t:0,p:l});
					}
					
					
				}
				attr.push(a);
				a = [];
				cnt = null;
				s = "";
				continue;
			}else if(!z && l == 0 && (c == '+' || c == '-' || c == '*' || c == '/' || c == '%' || c == '&' || c == '|' || c == '!')){
				if(s.length>0){
					cnt = cnt ? cnt : [];
					cnt.push(s);
					if(cnt.length == 1 && !isNaN(cnt[0])){//说明数字
						a.push({v:+cnt[0],t:-1,p:l});//v 名称，t 是否为字符串, p：是否为参数
					}else{
						father.push(cnt);
						a.push({v:cnt,t:0,p:l});
					}
					
					cnt = null;
					s = "";
				}
				switch(c){
					case '+' :
					a.push({v:0,t:2,p:0});//t = 2 代表运算符
					break;
					case '-' :
					a.push({v:1,t:2,p:0});//t = 2 代表运算符
					break;
					case '*' :
					a.push({v:2,t:2,p:0});//t = 2 代表运算符
					break;
					case '/' :
					a.push({v:3,t:2,p:0});//t = 2 代表运算符
					break;
					case '%' :
					a.push({v:4,t:2,p:0});//t = 2 代表运算符
					break;
					case '&' :
					a.push({v:5,t:2,p:0});//t = 2 代表运算符
					break;
					case '|' :
					a.push({v:6,t:2,p:0});//t = 2 代表运算符
					break;
					case '!' :
					a.push({v:7,t:2,p:0});//t = 2 代表运算符
					break;
				}
				
				continue;
			}
			s += c;
		}
		
		
		if(s.length > 0){//如果结束为变量名
			cnt = cnt ? cnt : [];
			
			cnt.push(s);
			if(cnt.length == 1 && !isNaN(cnt[0])){//说明数字
				a.push({v:+cnt[0],t:-1,p:l});//v 名称，t 是否为字符串, p：是否为参数
			}else{
				father.push(cnt);
				a.push({v:cnt,t:0,p:l});//v 名称，t 是否为字符串, p：是否为参数
			}
			
			attr.push(a);
			s = "";
			a = null
			cnt = null;
		}
		if(cnt){//如果结束为字符串
			father.push(cnt);
			a.push({v:cnt,t:0,p:l});//v 名称，t 是否为字符串, p：是否为参数
			attr.push(a);
			a = null;
			cnt = null;
		}
		
		if(a && a.length>0){//如果多个|后的结束
			attr.push(a);
		}
		if(name){
			return {father:father,attr:attr,node:node,name:name};
		}else{
			return {father:father,attr:attr,node:node};
		}
		
		
	}
	
	//读取字符串
	private static func readStr(v,i):int{
		var f = v[i++];
		var t = i;
		var b = false;
		var c = null;
		for(;i<v.length;i++){
			c = v[i];
			if(c == '\\'){
				b != b;
				continue;
			}
			if(!b && c == f){
				return {v:v.substring(t,i),i:i};
			}
		}
		throw new Error("string isn't over");
	}
	
	//读取属性
	private func gAttr(cmt,value){
		var i,j;
		var data;
		if(cmt[0][0] == '#'){
			data = _set[cmt[0]];
			i = 1;
		}else{
			data = _set['#'];
			i = 0;
		}
		if(data){
			data = data.data;
			for(;i<cmt.length - 1;i++){
				j = cmt[i];
				j = j[0] == '-' ?cmt.length + (j - 1) : j;
				data = data[j];
				if(data == undefined){
					return "";
				}
			}
			return value ? data[value] : data[cmt[cmt.length - 1]];
		}
		return "";
	}
	
	//读取属性和数据作用域
	private func gAttrD(cmt){
		var d;
		var i,j;
		var data;
		if(cmt[0][0] == '#'){
			data = _set[d = cmt[0]];
			i = 1;
		}else{
			d = '#'
			data = _set[d];
			i = 0;
		}
		if(data){
			data = data.data;
			for(;i<cmt.length - 1;i++){
				j = cmt[i];
				j = j[0] == '-' ? cmt.length + (j - 1) : j;
				
				if(data[j] == undefined){
					return {v:data,d:d}
				}
				data = data[j];
			}
			return {v:data,d:d,a:cmt[cmt.length - 1]};
		}
		return {};
	}
	
	
	//设置属性
	private func sAttr(cmt,value){
		var i,j;
		var data;
		if(cmt[0][0] == '#'){
			data = _set[cmt[0]];
			i = 1;
		}else{
			data = _set['#'];
			i = 0;
		}
		data = data.data;
		for(;i<cmt.length - 1;i++){
			j = cmt[i];
			j = j[0] == '-' ?cmt.length + (j - 1) : j;
			data = data[j];
		}
		return data[cmt[cmt.length - 1]] = value ;
	}
	
	
	//销毁这个模板
	public function destroy(){
		removeListener();
		node = null;
		map = null;
		aList.length = 0;
		aList = null;
		if(_bindData){
			_bindData.\$removeListener(_chg);
			_bindData = null;
		}
		var b = null;
		for(var k in _set){
			b = _set[k];
			if(b.\$removeListener){
				b.\$removeListener(_chg);
			}
			_set[k] = null;
		}
		_set = null;
		console.log("clear");
	}
	
}