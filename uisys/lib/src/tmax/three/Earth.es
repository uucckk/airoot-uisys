class Earth{
	public var type = "custom";
	var _mesh = null;
	var _geometry = null;
	var _material = null;
	var _x = 0,_y = 0,_z = 0,_r = 3,_w = 8,_h = 6;
	var _color:int = 0xaaaaaa;
	
	function init(r,w,h){
		_r = r;
		_w = w;
		_h = h;
		_geometry = new ? THREE.SphereGeometry(_r,_w,_h);
		_material = new ? THREE.MeshLambertMaterial({
			color: color,
			//emissive: color,
			transparent:true,
			opacity: 0.5
		});
		_material.needsUpdate=true;
		_mesh = new ? THREE.Mesh(_geometry, _material);
		window.g = _geometry;
	}
	
	
	/**
	 * 设置父亲级别
	 */
	set parent(value:TD):void{
		
	}
	
	get mesh(){
		window.mesh = _mesh;
		_mesh.position.z = z;
		_mesh.position.x = x;
		_mesh.position.y = y;
		_mesh.rotation.x = -Math.PI / 2;
		return _mesh;
	}
	
	
	set x(value:int):void{
		_x = parseInt(value);
		_mesh.position.x = x
	}
	
	get x():int{
		return _x;
	}
	
	set y(value:int):void{
		_y = parseInt(value);
		_mesh.position.y = y
	}
	
	get y():int{
		return _y;
	}
	
	set z(value:int):void{
		_z = parseInt(value);
		_mesh.position.z = z;
	}
	
	get z():int{
		return _z;
	}
	
	set color(value:String):void{
		if(isNumber(value)){
			_color = value;
		}else if(value.length>0 && value.charAt(0) == "#"){
			value = "0x" + value.substring(1);
			_color = parseInt(value);
		}
		_mesh.material.color = new ? THREE.Color(_color);
	}
	
	get color():uint{
		return _color;
	}
	
	set r(value){
		_r = value
		render();
	}
	
	get r(){
		return _r;
	}
	
	set w(value){
		_w = value;
		render();
	}
	
	set h(value){
		_h = value
		render();
	}
	
	function render(){
		_geometry = new ? THREE.SphereGeometry(_r,_w,_h);
		_mesh.geometry = _geometry;
	}
	
	//返回求替
	get geometry(){
		return _geometry;
	}
	
	function isNumber(obj) {  
		return obj === +obj  
	} 
}