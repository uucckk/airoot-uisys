class Text{
	public var type = "custom";
	var _text:String = "";
	var _x = 0,_y = 0,_z = 0;
	var _color:int = "#000000";
	var _mesh= null;
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
	public function set text(value:String):void{
		_text = value
	}
	
	public function get text():Array{
		return _text;
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
		_color = value;
		
	}
	
	public function get color():uint{
		return _color;
	}
	
	
	function isNumber(obj) {  
		return obj === +obj  
	}

	public get mesh(){
		let canvas = document.createElement('canvas');
			//画布建议使用正方形，长宽像素值使用2的次方数
			let ctx = canvas.getContext('2d');
			ctx.font = "bolder 24px Arial ";
			var width = ctx.measureText(text).width;
			console.log(width);
			canvas.width = width;
			canvas.height = width;
			//ctx.fillStyle = "rgb(58,255,250)";//设定画布颜色
			//ctx.fillRect(0,0,width,width);
			ctx.font = "bolder 24px Arial ";//设定字体属性，和css语法一致
			ctx.fillStyle = _color;//设定画布颜色
			ctx.fillText(text,0,width/2);//添加文字，并设置文字的位置
			ctx.globalAlpha =1;//设定canvas透明
			//将画布生成的图片作为贴图给精灵使用，并将精灵创建在设定好的位置
			let texture = new ?THREE.Texture(canvas);
			texture.needsUpdate = true; //告诉threejs材质需要更新
			let spriteMaterial = new ?THREE.PointsMaterial({ 
			  map:texture, //贴图
			  sizeAttenuation:true,  //开启尺寸衰减
			  size:width > 100 ? 2  : 1, //衰减尺寸
			  transparent:true, //允许透明
			  opacity:1, //设置不透明度
			});
			
			//创建点3D对象
			let geometry = new ?THREE.BufferGeometry();
			let vertices = [_x,_y,_z];
			geometry.addAttribute('position',new ?THREE.Float32BufferAttribute(vertices,3));
			let _mesh = new ?THREE.Points(geometry,spriteMaterial);//将材质赋给点实体
			return _mesh;
	}
}