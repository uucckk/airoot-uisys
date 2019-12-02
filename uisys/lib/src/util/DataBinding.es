/**
 * 数据绑定控制
 * @author sunxy
 * @type 功能包
 */
class DataBinding{
	public var \$ID = __UUID__();
	private var node = null;
	private var _data = {};
	private var _onChange = null;
	private var watchLst = {};//查看关注数据项目
	
	
	function init(node:Object){
		read(this.node = node,"");
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
		Event(\$ID,{name:name,value:newName});//发送变更广播
		var p = null;
		for(var k in watchLst){
			console.log(k,name);
			if(k == name){
				p = watchLst[k];
				for(var i = 0;i<p.length;i++){
					p[i](name,newName,oldName);
				}
			}
		}
	}
	
	
	
	/**
	 * 数据变更通知
	 */
	public function set onDataChange(value:Function){
		_onChange = value;
	}
	
	private function read(node:Object,parent:String){
		if(typeof node == "object"){
			for(var k in node){
				pushText(fx(parent + k,node[k]),node[k]);
				read(node[k],parent + k + ".")
			}
		}		
	}
	
	private function pushText(obj,node){
		//先把所有属性转化下
		var arr = obj.arr;
		var c = __inthis__;
		var p = null;
		for(var i = 0;i<arr.length;i++){
			p = arr[i];
			if(i == arr.length - 1){
				if(!c.hasOwnProperty(p)){
					setAttribute(c,p,obj);
				}
				if(typeof obj.value == "object"){
					if(obj.value instanceof Array){
						c[p] = [];
					}else{
						c[p] = {};
					}
					
				}else{
					c[p] = obj.value;
				}
			}
			c = c[p];
		}
		
	}
	
	/**
	 * obj
	 */
	private function fx(path,value){
		return {stat:path,arr:path.split("."),value:value};
	}
	
	
	
	private function setAttribute(obj,name,info){
		console.log(">>",name);
		var index = info.arr;
		var _value = null;
		var V = null;
		Object.defineProperty(obj,name,{
			set:function(value){
				var tValue = _value;
				if(_value && typeof value == "object"){
					if(V != value){
						V = value;
					}
					for(var k in _value){
						_value[k] = value[k];
					}
					return;
				}
				if(_value != value){
					var p = __inthis__;
					for(var i = 0;i<index.length - 1;i++){
						p = p[index[i]];
						if(!p){
							break;
						}
					}
					_value = value;
					notify(index,value,tValue);
					p[name] = value;
				}
				
				
			},
			get:function(){
				return _value;
			}
		,enumerable:true});
	}
	
	
	public function destroy(){
		//removeListener();
	}
	
}