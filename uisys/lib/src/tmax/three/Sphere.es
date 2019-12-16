class Sphere{
	public var type:String = "custom";
	private var _mesh = null;
	private var _geometry = null;
	private var _material = null;
	private var _x = 0,_y = 0,_z = 0,_r = 3,_w = 8,_h = 6;
	private var _color:int = 0xaaaaaa;
	
	function init(r,w,h){
		_r = r;
		_w = w;
		_h = h;
		_geometry = new ? THREE.SphereGeometry(_r,_w,_h);
		_material = new ? THREE.MeshLambertMaterial({
			color: color,
			//emissive: color,
			//transparent:true,
			//opacity: 0.5
		});
		_material.needsUpdate=true;
		_mesh = new ? THREE.Mesh(_geometry, _material);
	}
	
	
	/**
	 * 设置父亲级别
	 */
	set parent(value:TD):void{
		
	}
	
	set alpha(value){
		_material.transparent = value < 1 ? true : false;
		_material.opacity = value;
	}
	
	get mesh():Mesh{
		_mesh.position.z = z;
		_mesh.position.x = x;
		_mesh.position.y = y;
		return _mesh;
	}
	
	
	set x(value:int):void{
		_x = value - 0;
		_mesh.position.x = x
	}
	
	get x():int{
		return _x;
	}
	
	set y(value:int):void{
		_y = value - 0;
		_mesh.position.y = y
	}
	
	get y():int{
		return _y;
	}
	
	set z(value:int):void{
		_z = value - 0;
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
	//设置半径
	//@param value	整数
	set r(value){
		_r = value
		render();
	}
	
	get r():int{
		return _r;
	}
	
	set w(value:int){
		_w = value;
		render();
	}
	
	set h(value:int){
		_h = value
		render();
	}
	
	/**
	 * 设置材质
	 */
	set texture(value){
		_material.map = value.texture;
	}
	
	
	//返回求替
	get geometry(){
		return _geometry;
	}
	
	function isNumber(obj) {  
		return obj === +obj  
	} 
	
	function render(){
		_geometry = new ? THREE.SphereGeometry(_r,_w,_h);
		_mesh.geometry = _geometry;
	}
}