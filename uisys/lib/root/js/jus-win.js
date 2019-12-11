//弹出窗口管理
//唯一性句柄集合
var __MODULE_HANDLE__ = {};
//添加句柄
function AddHandle(objName,listener){
	if(!listener){
		alert("AddHandle: " + "please tell me handle listener");
		return;
	}
	if(__MODULE_HANDLE__[objName] == true){
		return;
	}
	
	if(__MODULE_HANDLE__[objName] && __MODULE_HANDLE__[objName].dom.parent().length != 0){
		if(listener){
			listener({target:__MODULE_HANDLE__[objName]});
		}
	}else{
		if(listener){
			var t = listener({target:null})
			if(t.listener){
				__MODULE_HANDLE__[objName] = true;
				t.listener(function(e){
					__MODULE_HANDLE__[objName] = e;
				});
			}else{
				__MODULE_HANDLE__[objName] = t;
			}
			
		}
	}
}
//弹出框管理
var PopManager = new function(){
	var __ZINDEX_CONTENT__ = [];
	//添加弹出框
	this.addPopUp = function(child,content){//弹出类，弹出容器
		if(!content){
			content = document.body;
		}
		if(child.dom){
			$(child.dom).bind("mousedown",function(){
				PopManager.bringToFront(child);
			});
			content.appendChild(child);
			var list = getList(content);
			if(list.length>0){
				for(var i in list){
					if(list[i].child == child){
						return;
					}
				}
				__ZINDEX_CONTENT__.push({content:content,child:child,index:999});
			}else{
				__ZINDEX_CONTENT__.push({content:content,child:child,index:999});
			}
			this.bringToFront(child);
		}
		
	}
	this.bringToFront = function(child){//显示在最前
		//找到自己所在对象
		var c = getChildData(child);
		//找到所有同级元素
		if(c){
			var pta = getList(c.content);
			//查看有没有999级别的
			var p = null;
			var f = false;
			for(var i in __ZINDEX_CONTENT__){
				p = __ZINDEX_CONTENT__[i];
				if(p != c && p.index == 999){
					f = true;
					break;
				}
			}
			if(f){
				//如果有999则降低所有级别
				for(var i in __ZINDEX_CONTENT__){
					__ZINDEX_CONTENT__[i].index --;
				}
			}
			c.index = 999;
			render();
		}
		
	}
	
	this.removePopUp = function(child){//删除窗口
		var p = null;
		for(var i in __ZINDEX_CONTENT__){
			p = __ZINDEX_CONTENT__[i];
			if(p.child == child && p.child.dom){
				__ZINDEX_CONTENT__.splice(i,1);
				var qt = p.child.dom;
				if(qt.attr("onRemove")){
					var clearFunc = "";
					qt.find("div[onRemove]").each(function(){
						clearFunc += this.getAttribute("onRemove") + ";\r\n";
					});
					if(clearFunc != ""){
						(new Function(clearFunc))();
					}
				}
				p.child.dom.remove();
				break;
			}
		}
	}
	
	function getList(content){//获取容器列表
		var list = [];
		for(var i in __ZINDEX_CONTENT__){
			if(__ZINDEX_CONTENT__[i].content == content){
				list.push(__ZINDEX_CONTENT__[i]);
			}
		}
		return list;	
	}
	function getChildData(child){
		var p = null;
		for(var i in __ZINDEX_CONTENT__){
			p = __ZINDEX_CONTENT__[i];
			if(p.child == child){
				return p;
			}
		}
		return null;
	}
	function render(){//渲染图层
		var p = null;
		var arr = [];
		for(var i in __ZINDEX_CONTENT__){
			p = __ZINDEX_CONTENT__[i];
			if(p.child.dom){
				if($(p.child.dom).parent().length != 0){
					$(p.child.dom).css({"position":"absolute","z-index":p.index});
					arr.push(p);
				}else{
					$(p.child.dom).unbind();
				}
				
			}else{
				p.child.css({"position":"absolute","z-index":p.index});
			}
		}
		__ZINDEX_CONTENT__.length = 0;
		__ZINDEX_CONTENT__ = arr;
	}
	
}