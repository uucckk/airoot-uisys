class Fog{
	var _near = 0,_far=200;
	var _color = 0xffffff;
	var _fog = null;
	function init(){
	
	}
	
	set parent(value){
		_fog = new ? THREE.Fog(_color, _near, _far);
		value.scene.fog = _fog;
	}
	
	/**
	 *雾气最近距离
	 */
	set near(value){
		_near = value;
		render();
	}
	/**
	 * 可见范围最大
	 */
	set far(value){
		_far = value;
		render();
	}
	
	set color(value){
		_color = value;
		render();
	}
	
	private function render(){
		if(!_fog){
			return;
		}
		_fog.near = _near;
		_fog.far = _far;
		_fog.color = new ? THREE.Color(_color);
	}
}