class MapTool{
	private static var EARTH_RADIUS = 6378137;//地球半径米

		
	private static function rad(d:double):double
	{
	   return d * Math.PI / 180.0;
	} 
	public static function GetDistance(lng0,lat0,lng1,lat1):double
	{
		lng0 = rad(lng0);
		lat0 = rad(lat0);
		lng1 = rad(lng1);
		lat1 = rad(lat1);
	   	//先将经纬度转为空间坐标系
		//01.lng0,lat0
		var x0 = EARTH_RADIUS*Math.cos(lat0)*Math.sin(lng0);
		var y0 = EARTH_RADIUS*Math.sin(lat0);
		var z0 = - EARTH_RADIUS*Math.cos(lat0)*Math.cos(lng0);
		
		
		//02. lng1,lat1
		var x1 = EARTH_RADIUS*Math.cos(lat1)*Math.sin(lng1);
		var y1 = EARTH_RADIUS*Math.sin(lat1);
		var z1 = - EARTH_RADIUS*Math.cos(lat1)*Math.cos(lng1);
		//console.log("A",x0,y0,z0,"B",x1,y1,z1);
		
		//求弦长
		var d = Math.sqrt(Math.pow(x1 - x0,2) + Math.pow(y1 - y0,2) + Math.pow(z1 - z0,2));
		//求弦对角
		var a = Math.asin(d/2/EARTH_RADIUS);
		
	   return EARTH_RADIUS*a*2;
	}
	
	
	
	public static function latitude(data,level):Number{
		var data:Number = data*Math.PI/180;
		var CJ:Number = 85.05112877980659*Math.PI/180;//89.89265
		var A:Number = 6378137;//地球半径
		var B:Number = 6378137;
		var E:Number = Math.sqrt(1.0-(B/A)*(B/A));
		
		//这个是将地图考虑为一椭球球体的标准
		/*
		var k:Number = Math.pow((1.0-E*Math.sin(data))/(1.0+E*Math.sin(data)),E/2);
		var CJK:Number = Math.pow((1.0-E*Math.sin(CJ))/(1.0+E*Math.sin(CJ)),E/2);
		CJK = k=1;
		*/
		var tmp = Math.pow(2,level);
		return ((tmp/2) - (Math.log(Math.tan(data/2+Math.PI/4)))*(tmp/2)/(Math.log(Math.tan(CJ/2+Math.PI/4))))*256;

	}
	
	public static function longitude(data,level):Number{
		var CJ:Number = 180.0;
		return (data-0+CJ)/(CJ*2)*Math.pow(2,level)*256;

	};
	
	
	/**
	 * 横坐标转换精度
	 */
	public static function x2Lng(value:Number,level:int):Number{
		var CJ:Number = 180.0;
		return value/256/Math.pow(2,level)*(CJ*2) - CJ;
	}
	
	
	/*
	public static function y2Lat(y:int,level:int):Number{
		var wL:int = Math.pow(2,level)*256;
		var lat:Number = (Math.atan(Math.exp((wL/2-y)/(wL/2)*Math.PI))-(Math.PI/4))*2.0*180/Math.PI;
		return lat;
	}
	*/
	
	public static function y2Lat(y:int,level:int):Number{
		var wL:int = Math.pow(2,level)*256/4;
		var lat:Number = (wL - y)/wL*90;
		return lat;
	}
	
	
	
	
	
	/**
	 * 通过长度（米）和纬度偏移量求纬度变化角度
	 */
	private static function GetAnagleByDistanceAndLat(distance:Number):Lat{
		return distance*360/(2*EARTH_RADIUS*Math.PI);
	}
	
	
	
	
	
	/**
	 * 获取一个圆形
	 */
	public static function GetCircle(lng:Number,lat:Number,r:Number,angle:Number):Array{
		//console.log("-distance",GetDistance(lng,lat,lng,90));
		//console.log("-distance",GetDistance(lng,lat,126,90));
		var arr = new Array();
		var a = rad(GetAnagleByDistanceAndLat(r));
		if(!angle){
			angle = Math.PI/20;
		}
		var b = 0;
		while(b<Math.PI*2){
			var R = Math.sqrt(Math.pow(Math.cos(a),2) + Math.pow(Math.sin(a)*Math.cos(b),2));
			var A = Math.atan(-Math.tan(a)*Math.cos(b)) + rad(lat);
			var jd = Math.atan(Math.sin(a)*Math.sin(b)/R/Math.cos(A)) + rad(lng);
			var wd = Math.asin(R*Math.sin(A));
			//console.log("distance",GetDistance(lng,lat,jd*180/Math.PI,wd*180/Math.PI),"from",lng,lat,jd*180/Math.PI,wd*180/Math.PI);
			arr.push(jd*180/Math.PI,wd*180/Math.PI);
			b += angle;
		}
		return arr;
		
	}
	
	/** 
	 * 判断点在多边形内
	 */
	public static function inPoly(point, vs) {
		// ray-casting algorithm based on
		// http://www.ecse.rpi.edu/Homepages/wrf/Research/Short_Notes/pnpoly.html
		
		var x = point[0], y = point[1];
		
		var inside = false;
		for (var i = 0, j = vs.length - 1; i < vs.length; j = i++) {
			var xi = vs[i][0], yi = vs[i][1];
			var xj = vs[j][0], yj = vs[j][1];
			
			var intersect = ((yi > y) != (yj > y))
				&& (x < (xj - xi) * (y - yi) / (yj - yi) + xi);
			if (intersect) inside = !inside;
		}
		
		return inside;
	};
	
	

	
}
