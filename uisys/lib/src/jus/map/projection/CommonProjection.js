/**
 * 经纬度直投
 */
class CommonProjection{
	public var tileOrginX:int = 0;
	public var tileOrginY:int = 0;
	private var left:int = 0;
	private var right:int = 0;
	private var top:int = 0;
	private var bottom:int = 0;
	public var xFlag:int = 1;//Y周向上
	public var yFlag:int = 1;//X周向右
	
	
	/**
	 * @param tileOrginX	瓦片起始X
	 * @param tileOrginY	瓦片起始Y
	 * @param left			投影起始经度
	 * @param top			投影起始纬度
	 * @param right			投影终止经度
	 * @param bottom		投影起始纬度
	 */
	function init(tileOrginX:int,tileOrginY:int,left:Number,top:Number,right:Number,bottom:Number,xFlag:int,yFlag:int):void{
		this.tileOrginX = tileOrginX;
		this.tileOrginY = tileOrginY;
		this.left = left;
		this.right = right;
		this.top = top
		this.bottom = bottom;
		this.xFlag = xFlag;
		this.yFlag = yFlag;
	}
	
	
	
	/**
	 *
	 */
	public function latitude(data,layer):Number{
		data = parseFloat(data);
		if(yFlag == 1){//Y上
			return data/layer.extent.maxY*layer.height;
		}else{//Y下
			return (90 - data)/90*layer.height/2;
		}
		

	}
	
	public function longitude(data,layer):Number{
		data = parseFloat(data);
		var CJ:Number = 180;
		if(yFlag == 1){//Y上
			CJ = layer.extent.maxX;
			return (data)/CJ*layer.width;
		}else{//Y下
			return (data + 180)/360*layer.width;
		}

	};
	
	
	/**
	 * 横坐标转换精度
	 */
	public function x2Lng(value:Number,layer:Object):Number{
		if(yFlag == 1){
			var CJ:Number = Math.abs(layer.extent.minX) + layer.extent.maxX;
			var v = value/layer.width*CJ;
			return v;
		}else{//Y下
			var CJ:Number = Math.abs(layer.extent.minX) + layer.extent.maxX;
			var v = value/layer.width*CJ - layer.extent.maxX;
			return v;
		}
		
	}
	
	
	
	public function y2Lat(y:int,layer:Object):Number{
		var wL:int = layer.height;
		if(yFlag == 1){//Y上
			var lat:Number = y/wL*layer.extent.maxY;
			return lat;
		}else{//Y下
			wL /= 2;
			var lat:Number = (wL - y)/wL*90;
			return lat;
		}
	}
	
}
