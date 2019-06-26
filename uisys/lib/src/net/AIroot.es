/**
 * 神经元连接
 * @author sunxy
 * @type 通讯类
 */
class AIRoot{
	var _userName = null;
	var hid = -1;//心跳id
	var ws;
	var msgFrame:int = 0;
	var host:String;
	var port:int;
	var protocol:String;
	var _onopen:Function = null,_onmessage:Function = null,_onclose:Function = null,_onerror:Function = null;
	function init(host:String,port:int,protocol:String){
		this.host = host;
		this.port = port;
		this.protocol = protocol;
	}
	
	public get userName(){
		return _userName;
	}
	
	/**
	 * 登录
	 * @param userName 	用户名
	 * @param passwd	密码
	 */
	public function login(userName:String,passwd:String){
		this._userName = userName;
		console.log(protocol + "://" + host + ":" + port + "/ws");
		ws = new WebSocket(protocol + "://" + host + ":" + port + "/ws");
		ws.onopen = function(){
			ws.send("login " + userName + " " + passwd);
		}
		ws.onmessage = function(e){
			console.log(">>",e.data);
			if(_onopen){
				var p:int = e.data.indexOf(" ");
				_onopen({status:e.data.substring(0,p),data:e.data.substring(p + 1)});
				
			}
			if(e.data.length >6 && e.data.substring(0,7) == "accept "){
				@this._userName = e.data.substring(7);
				hid = setInterval(function(){
					ws.send("\0");//心跳
				},10000);
				ws.onmessage = function(e){
					if(_onmessage){
						_onmessage(formatPackage(e.data));
					}
				}
			}else{
				ws.close();
				trace("登陆失败",e.data);
			}
			
		}
		ws.onclose = function(e){
			clearInterval(hid);
			if(_onclose){
				_onclose(e)
			}
		}
		ws.onerror = function(e){
			clearInterval(hid);
			if(_onerror){
				_onerror(e)
			}
		}
		
	}
	
	var key:Array = ["router","uuid","frame","data"];
	function formatPackage(str:String):Object{
		var obj:Object = {};
		var arr:Array = str.split("\0");
		for(var i:int = 0;i<arr.length;i++){
			obj[key[i]] = arr[i];
		}
		return obj;
		
	}
	
	/**
	 * 返回信息
	 */
	public set onMessage(f:Function){
		_onmessage = f;
	}
	
	public set onOpen(f:Function){
		_onopen = f;
		
	}
	
	public set onClose(f:Function){
		_onclose = f;
	}
	
	public get onClose(){
		return _onclose;
	}
	
	public set onError(f:Function){
		_onerror = f;
	}
	
	public get onError(){
		return _onerror;
	}
	
	/**
	 * 发送信息
	 * @param router 	路由
	 * @param uuid	 	信息唯一标识
	 * @param value		实际的传递值
	 */
	public function send(router:String,uuid:String,value:String){
		msgFrame++;
		ws.send(router + "\0" + uuid + "\0" + msgFrame + "\0" +value + "\0");
	}
}//