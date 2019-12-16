class Line{
	public var type = "custom";
	var _data:Array = [];
	var _x = 0,_y = 0,_z = 0;
	var _color:int = 0xaaaaaa;
	var _mesh = null;
	function init(){
		//alert("A");
	}
	
	
	/**
	 * 设置父亲级别
	 */
	public function set parent(value:TD):void{
		
	}
	/**
	 * 设置数据
	 */
	public function set data(arr:String):void{
		var i:int = 0;
		var p = null;
		for(i = 0;i<arr.length;i++){
			p = arr[i];
			_data.push({x:p[0] - 0,y:p[1] -0});
		}
	}
	
	public function get data():Array{
		return _data;
	}
	
	
	
	
	public function set x(value:int):void{
		_x = value - 0;
	}
	
	public function get x():int{
		return _x;
	}
	
	public function set y(value:int):void{
		_y = value - 0;
	}
	
	public function get y():int{
		return _y;
	}
	
	public function set z(value:int):void{
		_z = value - 0;
	}
	
	public function get z():int{
		return _z;
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
	
	public get mesh(){
		 //定义材质THREE.LineBasicMaterial . MeshBasicMaterial...都可以
		var material = new ?THREE.LineBasicMaterial({color:_color});
		// 空几何体，里面没有点的信息,不想BoxGeometry已经有一系列点，组成方形了。
		var geometry = new ?THREE.Geometry();
		// 给空白几何体添加点信息，这里写3个点，geometry会把这些点自动组合成线，面。
		var paths = data;
		var p = null;
		for(var i:int = 0;i<paths.length;i++){
			p = paths[i];
			geometry.vertices.push(new ?THREE.Vector3(p.x,p.y,0));
		}
		//线构造
		_mesh = new ?THREE.Line(geometry,material);
		_mesh.position.z = z;
		_mesh.position.x = x;
		_mesh.position.y = y;
		return _mesh;
	}
	
	function isNumber(obj) {  
		return obj === +obj  
	} 
}