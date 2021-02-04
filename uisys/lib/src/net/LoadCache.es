/**
 * 加载资源缓存
 */
class LoaderCache{
	var urls = [];
	
	func init(){
		
	}
	//当加载完成时候
	set onComplete(v){
		
	}
	
	//当有某个加载项加载完成时候
	set onLoad(v){
		
	}
	
	//增加加载URL连接
	func add(url){
		urls.push({url:url});
	}
	
	//开始加载
	func load(){
		for(var i = 0;i<urls.length;i++){
			
		}
	}
}