class SkyBox{
	var type = "content";
	var _parent = null;
	var list = [];
	function init(value){
		var t = this;
		console.log(">>",value);
	}
	set parent(obj:td){
		_parent = obj;
		//天空盒
		var path = "img/box/";       //设置路径
		var format = '.png';                        //设定格式
		var urls = [
		   path + 'Side'+ format,     
		   path + 'Side'+ format,
		   path + 'top' + format,
		   path + 'bottom' + format,
		   path + 'Side' + format,
		   path + 'Side' + format
		];
		var textureCube = new ? THREE.CubeTextureLoader().load( urls );
		obj.scene.background = textureCube; //作为背景贴图	
		//scene.background = new ? THREE.Color(0x112233);
		for(var i in list){
			_parent.addChild(list[i]);
		}
	}
	
	set children(value){
		list.push(value);
	}
}