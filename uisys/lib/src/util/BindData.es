class BindData{
	static var count = 1;
	static func ufo(p){
		if(!p){
			return;
		}
		if(p.constructor == Object){
			if(!p._u$){
				Object.defineProperty(p, '_u$', {
					writeable:false,
					configurable:false,
					enumerable: false,
					value:count++
				})
			}
			
			for(var i in p){
				ufo(p[i]);
			}
		}else if(p.constructor == Array){
			for(var j = 0;j<p.length;j++){
				ufo(p[j]);
			}
		}
		
	}
	static func eufo(p,f){
		if(p.constructor == Object){
			f(p,1)
			for(var i in p){
				eufo(p[i],f);
			}
		}else if(p.constructor == Array){
			for(var j = 0;j<p.length;j++){
				eufo(p[j],f);
			}
		}
		
	}
	var _data = null;
	var _array = null;
	static var _fList = [];
	var _put = Map ? new Map() : {};
	var put:Func;
	function init(value){
		data = value;
		put = Map ? putM : putO;
	}
	
	//添加监听
	function \$addListener(f,n){
		if(f){
			for(var i = 0;i<_fList.length;i++){
				if(_fList[i].f == f && _fList[i].n == n){
					return;
				}
			}
			_fList.push({f:f,n:n});
		}
	}
	
	//移除监听
	function \$removeListener(f){
		if(f){
			for(var i = 0;i<_fList.length;i++){
				if(_fList[i].f == f){
					_fList.splice(i,1);
				}
			}
		}
	}
	/**
	 * cls		事件类型
	 * obj 		变更对象
	 * new_id	新ID
	 * a		更改的属性名
	 * r		是否需要刷新
	 */
	function event(cls,obj,new_id,a,r){
		var p = null;
		for(var i = 0;i<_fList.length;i++){
			p = _fList[i];
			p.f(p.n,cls,obj,new_id,a,r);
		}
	}
	
	set data(value){
		ufo(value);
		if(_data){
			var t = _data;
			_data = value;
			eufo(t,event);
		}else{
			_data = value;
		}
		
	}
	
	get data(){
		return _data;
	}
	
	func g(attr,a){
		if(!a){
			a = _data;
		}
		var d;
		if(attr.indexOf){
			if(attr.indexOf(".") != -1 || attr.indexOf("[") != -1){
				d = eval("a." + attr);
			}else{
				d = a[attr];
			}
		}else{
			d = a[attr];
		}
		
		var o = {
			s:s,
			g:function(v){
				return @this.g(v,d);
			},
			v:d,
			toString:function(){
				return this.v;
			},
			_p:a,
			_t:attr,
			
		}
		if(d instanceof Array){
			o.push = function(v){
				var t = this.v.push.apply(this.v,arguments);
				event('p',t);//尾新增元素
				return t;
			}
			o.pop = function(){
				var t = this.v.pop();
				event('d',t);//删除t
				return t;
			}
			o.unshift = function(){
				var t = this.v.unshift.apply(this.v,arguments);
				event('u',t);//头新增元素
				return t;
			}
			o.shift = function(){
				var t = this.v.shift();
				event('d',t);//删除t
				return t;
			}
			
			o.splice = function(){
				var t = this.v.splice.apply(this.v,arguments);
				event('s',t);//删除[t]
				return t;
			}
		}
		
		return o;
	}
	
	func s(attr,value){
		var a = this.v ? this.v : _data;
		var d;
		var t = this._t;
		var p = this._p;
		put(attr,value);
		
		return {
			s:s,
			g:g,
			v:a,
			toString:function(){
				return this.v;
			},
			_p:p,
			_t:t,
		}
	}
	private var _t = -1;
	private func putO(a,v){
		_put[a] = v;
		if(_t == -1){
			_t = setTimeout(rdO,0);
		}
		
	}
	
	private func putM(a,v){
		console.log(a,v);
		_put.set(a,v);
		if(_t == -1){
			_t = setTimeout(rdM,0);
		}
		
	}
	
	private func rdO(){
		clearTimeout(_t);
		_t = -1;
		var a,t = _data;
		for(var k in _put){
			a = k;
			a = a.replace(/\]/g,'').replace(/\[/g,'.').split('.');
			for(var i = 0;i<a.length - 1;i++){
				t = t[a[i]];
			}
			t[a[a.length - 1]] = _put[k]
			if(v instanceof Array){
				ufo(v);
				event('c',t,null,null,true);
			}else{
				event('c',t);
			}
		}
		_put = {};
	}
	
	private func rdM(){
		clearTimeout(_t);
		_t = -1;
		var a,t = _data;
		_put.forEach(func(v, k){
			a = k;
			a = a.replace(/\]/g,'').replace(/\[/g,'.').split('.');
			for(var i = 0;i<a.length - 1;i++){
				t = t[a[i]];
			}
			console.log("chg",a[a.length - 1]);
			t[a[a.length - 1]] = v;
			if(v instanceof Array){
				ufo(v);
				event('c',t,null,null,true);
			}else{
				event('c',t);
			}
		})
		_put.clear();
	}
	
}