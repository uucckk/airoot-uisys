<!-- 地图瓦力 -->
class AIMark{
	var map:MapView = null;
	
	public var level:int = 0;
	public var left:int = 0;
	public var top:int = 0;
	
	public var image:Image = null;
	public var enabled:Boolean = false;
	
	
	private var _self:AIMark = null;
	private var markWidth:int = 256;
	public var mapX:Number = 0;//实际的MAPX
	public var mapY:Number = 0;//实际的MAPY
	
	var _x:int = 0;
	var _y:int = 0;
	
	/**
	 * 
	 */
	function init(map:MapView,level:int,left:int,top:int,filter:Function){
		_self = __inthis__;
		map.hav[left + "_" + top] = _self;
		this.map = map;
		this.left = left;
		this.top = top;
		this.level = level;
		mapX = left*markWidth;
		mapY = top*markWidth;
		
		image = new Image();
		//image.src = map.server + level + "/" + left + "/" + top + ".png";
		//image.src = map.server + "&x=" + left + "&y=" + top + "&l=" + level;
		image.src = map.getServerImage(left,top,level);
		image.onload = function(){
			enabled = (image && map.level == level) ? true :false;
			if(enabled){
				map.drawMark(_self);
			}
			
		}
		image.onerror = function(e){
			
		}

		requestAnimationFrame(enterframe);
	}
	
	
	
	public function set x(value:int):void{
		_x = value;
	}
	
	public function set y(value:int):void{
		_y = value;
	}
	
	public function get x():int{
		return _x;
	}
	
	public function get y():int{
		return _y;
	}
	
	
	public function enterframe():void{
		var bl1 = 2;
		var bl2 = 2;
		if(map.projection.yFlag == 1){
			if(mapX>map.mapRect.right + 256*bl1){
				destroy();
				return;
			}
			if(mapX<map.mapRect.left - 256*bl1){
				destroy();
				return;
			}
			if(mapY >map.mapRect.top +256*bl1){
				destroy();
				return;
			}
			if(mapY<map.mapRect.bottom - 256*bl1){
				destroy();
				return;
			}
			
			//left;
			if(mapX >map.mapRect.left - 256*bl2 &&!map.hav[(left - 1) + "_" + top]){
				var mark = new component.map.AIMark(map,level,left-1,top);
			}
			//right;
			if(mapX <map.mapRect.right + 256*bl2 && !map.hav[(left + 1) + "_" + top]){
				var mark = new component.map.AIMark(map,level,left+1,top);
			}
			//top;
			if(mapY <map.mapRect.top +256*bl2 && !map.hav[left + "_" + (top + 1)]){
				var mark = new component.map.AIMark(map,level,left,(top + 1));
			}
			//bottom
			if(mapY  > map.mapRect.bottom - 256*bl2 && !map.hav[left + "_" + (top - 1)]){
				var mark = new component.map.AIMark(map,level,left,(top - 1));
			}
			
		}else{
			if(mapX>map.mapRect.right + 256*bl1){
				destroy();
				return;
			}
			if(mapX<map.mapRect.left - 256*bl1){
				destroy();
				return;
			}
			if(mapY<map.mapRect.top - 256*bl1){
				destroy();
				return;
			}
			if(mapY>map.mapRect.bottom + 256*bl1){
				destroy();
				return;
			}
			
			//left;
			if(mapX >map.mapRect.left - 256*bl2 &&!map.hav[(left - 1) + "_" + top]){
				var mark = new component.map.AIMark(map,level,left-1,top);
			}
			//right;
			if(mapX <map.mapRect.right + 256*bl2 && !map.hav[(left + 1) + "_" + top]){
				var mark = new component.map.AIMark(map,level,left+1,top);
			}
			//top;
			if(mapY >map.mapRect.top -256*bl2 && !map.hav[left + "_" + (top - 1)]){
				var mark = new component.map.AIMark(map,level,left,(top - 1));
			}
			//bottom
			if(mapY  < map.mapRect.bottom + 256*bl2 && !map.hav[left + "_" + (top + 1)]){
				var mark = new component.map.AIMark(map,level,left,(top + 1));
			}
			
		}
		
	}
	
	
	
	public function destroy(){
		enabled = false;
		delete map.hav[left + "_" + top];
		image = null;
		//console.log(left + "_" + top + " REMOVE");
	}
}