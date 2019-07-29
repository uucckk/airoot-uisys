class Drag{
	var dom = null;
	var cx = 0;
	var cy = 0;
	var mx = 0;
	var my = 0;
	var dx = 0;
	var dy = 0;
	var df = false;
	function init(obj){
		dom = obj.dom;
		dom.style.position = "absolute";
		dom.style.left = "0px";
		dom.style.top = "0px";
		dom.addEventListener("mousedown",mouseEvt);
		window.addEventListener("mouseup",mouseEvt);
		initWin();
	}
	
	/**
	 * 
	 */
	function dframes(value){
		var f = document.body.getElementsByTagName("iframe");
		var p = null;
		if(value){
			for(var i = 0;i<f.length;i++){
				p = f[i];
				if(!p.style){
					p.style = "";
				}
				p.style["pointer-events"] = "unset";
			}
		}else{
			for(var i = 0;i<f.length;i++){
				p = f[i];
				if(!p.style){
					p.style = "";
				}
				p.style["pointer-events"] = "none";
			}
		}
	}
	
	function mouseEvt(e){
		switch(e.type){
			case "mousedown":
				document.onselectstart = new Function("return false;");
				dframes(false);
				mx = e.clientX;
				my = e.clientY;
				cx = dom.offsetLeft;
				cy = dom.offsetTop;
				df = true;
			break;
			case "mousemove":
				if(df){
					dx = e.clientX - mx;
					dy = e.clientY - my;
					cx += dx;
					cy += dy;
					dom.style.left = cx + "px";
					dom.style.top = cy + "px";
					mx = e.clientX;
					my = e.clientY;
					
				}
			break;
			case "mouseup":
				df = false;
				document.onselectstart = null;
				dframes(true);
			break;
		}
	}
	
	function initWin(){
		window.addEventListener("mousemove",mouseEvt);
	}
}