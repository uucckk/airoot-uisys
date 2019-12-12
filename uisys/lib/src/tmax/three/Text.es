class Text{
	public var type = "text";
	var _text:String = "";
	var _x = 0,_y = 0,_z = 0;
	var _color:int = 0xaaaaaa;
	
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
		_x = parseInt(value);
	}
	
	public function get x():int{
		return _x;
	}
	
	public function set y(value:int):void{
		_y = parseInt(value);
	}
	
	public function get y():int{
		return _y;
	}
	
	public function set z(value:int):void{
		_z = parseInt(value);
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
	
	function isNumber(obj) {  
		return obj === +obj  
	} 
}