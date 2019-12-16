class Shape{
	public var type = "custom";
	var _data:Array = [];
	var _depth:int = 10;//厚度
	var _x = 0,_y = 0,_z = 0;
	var _color:int = 0xaaaaaa;
	var _mesh = null;
	var _material = null;
	var _alpha = 1;
	function init(value:Array){
		data = value;
		_material = new ? THREE.MeshLambertMaterial({	//MeshLambertMaterial,MeshBasicMaterial
			color: color,
			//emissive: color,
			//transparent:true,
			//opacity: 0.5
		});
		_material.needsUpdate=true;
		var shape = createShape(data);
		var shape3d = new ? THREE.ExtrudeBufferGeometry( shape, {
			amount: depth,
			bevelEnabled: false,
			bevelSize:1,
			bevelThickness:1,
			curveSegments:12
		} );
		
		_mesh = new ? THREE.Mesh( shape3d, _material );
	}
	
	
	/**
	 * 设置父亲级别
	 */
	public function set parent(value:TD):void{
		
	}
	/**
	 * 设置数据
	 */
	public function set data(data:Array):void{
		var i:int = 0;
		var p = null;
		for(i = 0;i<data.length;i++){
			p = data[i];
			_data.push({x:p[0] - 0,y:p[1] -0});
		}
	}
	
	public function get data():Array{
		return _data;
	}
	
	public function set depth(value:String):void{
		_depth = value;
		_material = new ? THREE.MeshLambertMaterial({	//MeshLambertMaterial,MeshBasicMaterial
			color: color,
			//emissive: color,
			//transparent:true,
			//opacity: 0.5
		});
		_material.needsUpdate=true;
		var shape = createShape(data);
		var shape3d = new ? THREE.ExtrudeBufferGeometry( shape, {
			amount: depth,
			bevelEnabled: false,
			bevelSize:1,
			bevelThickness:1,
			curveSegments:12
		} );
		
		_mesh = new ? THREE.Mesh( shape3d, _material );
	}
	
	public function get depth():int{
		return _depth;
	}
	
	
	public function set x(value:int):void{
		_x = value - 0;
		render();
	}
	
	public function get x():int{
		return _x;
	}
	
	public function set y(value:int):void{
		_y = value - 0;
		render();
	}
	
	public function get y():int{
		return _y;
	}
	
	public function set z(value:int):void{
		_z = value - 0;
		render();
	}
	
	public function get z():int{
		return _z;
	}
	
	public set alpha(value){
		_alpha = value;
		_material.transparent = value < 1 ? true : false;
		_material.opacity = value;
	}
	
	public get alpha(){
		return _alpha;
	}
	
	public get mesh(){
		
		//mesh.receiveShadow = true;
		//mesh.castShadow = true;
		//mesh.rotation.x = Math.PI;
		
		//mesh.rotation.x = -Math.PI / 2;
		return _mesh;
	}
	
	public function set color(value:String):void{
		if(isNumber(value)){
			_color = value;
		}else if(value.length>0 && value.charAt(0) == "#"){
			value = "0x" + value.substring(1);
			_color = parseInt(value);
		}
		
	}
	
	public function get color():uint{
		return _color;
	}
	
	function isNumber(obj) {  
		return obj === +obj  
	}
	
	/**
	 * 创建shape
	 */
	function createShape(points:Array):Shape{
		var shape = new ? THREE.Shape();
		var p = null;
		if(points.length>0){
			p = points[0];
			shape.moveTo(p.x,p.y);
		}
		for(var i:int = 1;i<points.length;i++){
			p = points[i];
			shape.lineTo(p.x,p.y);
		}
		return shape;
	}
	public get base(){
		return _mesh;
	}
	function render(){
		_mesh.position.x = x;
		_mesh.position.y = y;
		_mesh.position.z = z;
	}
}