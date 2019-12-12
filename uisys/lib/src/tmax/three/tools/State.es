import js/libs/stats.min.js;
class State{
	var stats = null;
	function init(){
		stats = new Stats();
	}
	
	function set parent(obj){
		obj.container.appendChild(stats.dom);
		obj.enterframe(this,function(){
			stats.update();
		});
	}
	
	
}