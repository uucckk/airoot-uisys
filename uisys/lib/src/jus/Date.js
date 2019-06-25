import component.Date.DatePanel;
var _self = $(dom);
public var onChange = null;
private var dp = null;
private var _dateFormat:String = null;
function init(yyyy:int,mm:int,dd:int){
	if(_self.attr("value")){
		ipt.value = _self.attr("value");
	}
	if(_self.attr("format")){
		_dateFormat = _self.attr("format");
	}
   /**
	 * 日期控件的显隐控制
	 */
	_self.on("click","#$ipt,#$img",function(){
		if(dp){
			dp.remove();
			_self.find(".inp-date").css("border","1px solid #ccc")
			_self.find(".DataIcon").attr("src","img/Common/ico_calendar_N.png");
			dp = null;
			return;
		}else{
			_self.find(".inp-date").css("border","1px solid #3685EA")
			_self.find(".DataIcon").attr("src","img/Common/ico_calendar_P.png");

		}

		dp = new DatePanel();      // 引入DatePanel文件
	 	dp.x = _self.offset().left;
	 	dp.y = _self.offset().top + _self.outerHeight(true);
	 	$("body").addChild(dp);
	 	

	 	dp.onClick = function(e){
	 		var  selfValue = _dateFormat;
	 		var dataArr = e.value.split("-");
	 		if(dataArr[1] < 10){
	 			dataArr[1] = dataArr[1].substr(1,1);
	 		}
	 		// year
	 		var selfValue = selfValue.replace(/yyyy|YYYY/,dataArr[0]);   
	 		var twoYear = dataArr[0].substr(2,2);
    		selfValue = selfValue.replace(/yy|YY/,(twoYear % 100)>9?(twoYear % 100).toString():'0' + (twoYear % 100));   

    		// month
 			selfValue = selfValue.replace(/MM/,dataArr[1]>9?dataArr[1].toString():'0' + dataArr[1]);   
    		selfValue = selfValue.replace(/M|m/g,dataArr[1]);   

    		// date
    		selfValue = selfValue.replace(/dd|DD/,dataArr[2]>9?dataArr[2].toString():'0' + dataArr[2]);   
    		selfValue = selfValue.replace(/d|D/g,dataArr[2]);   

    		// hour
    		// str=str.replace(/hh|HH/,this.getHours()>9?this.getHours().toString():'0' + this.getHours());   
		    // str=str.replace(/h|H/g,this.getHours()); 
		    // //min 
		    // str=str.replace(/mm/,this.getMinutes()>9?this.getMinutes().toString():'0' + this.getMinutes());   
		    // str=str.replace(/m/g,this.getMinutes());   
		    // //sec
		    // str=str.replace(/ss|SS/,this.getSeconds()>9?this.getSeconds().toString():'0' + this.getSeconds());   
		    // str=str.replace(/s|S/g,this.getSeconds());   

	 		ipt.value = selfValue;    // 日期填入Input控件
	 		e.remove();
	 		_self.find(".inp-date").css("border","1px solid #ccc");
			_self.find(".DataIcon").attr("src","img/Common/ico_calendar_N.png");
	 		dp = null;
	 		
	 	}
	});

	$(window).on("click",clickEvt);
}

private function clickEvt(e){
	var $this = $(e.target);

	if($this.is("input") && $this.attr("id") == #ipt.attr("id")){
              // 当前点击的是input 并且当前点击的是自己而不是前一个Input控件
	}else{    // 当前点击的input不是此刻展开的input
		if(dp){
			dp.remove();
			_self.find(".inp-date").css("border","1px solid #ccc");
			_self.find(".DataIcon").attr("src","img/Common/ico_calendar_N.png");
			dp = null;	
		}
	}
	
}


public function get value():String{
	return ipt.value;
}
public function set value(value:String):void{
	#ipt.val(value);
}

/**
 * 设置日期格式
 */
public function set format(value:String):void{
	_dateFormat = value;
}
public function get format():void{
	return _dateFormat;
}


public function finalize():void{
	$(window).off("click",clickEvt);    // 解绑绑定在window上的事件
}


