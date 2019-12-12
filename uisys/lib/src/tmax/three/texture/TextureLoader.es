class TextureLoader{
	var _src = null;
	
	set src(value){
		_src = value;
	}
	
	get texture(){
		return new ? THREE.TextureLoader().load(_src);
	}
}