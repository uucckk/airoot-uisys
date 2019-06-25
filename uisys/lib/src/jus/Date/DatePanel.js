public var onClick:Function = null;
private var showYear:int = 0.0;
private var showMonth:int = 0.0;
private var DATE:String = null;
private var fatherDate:String = null;
function  init(){
	dateObj = new Date(); 
	nowDate = dateObj.getDate();       // 当前几号
	nowMonth = dateObj.getMonth()+1;   // 当前月
	nowYear = dateObj.getFullYear();   // 当前年	
	 
	_self.find(".top-change-year").html(nowYear);
	nowMonth = extra(nowMonth);
	_self.find(".top-change-month").html(nowMonth);
	showMonth = _self.find(".top-change-month").html();
	showYear = _self.find(".top-change-year").html();
	Day();
	_self.find(".change-bg").html(showMonth);
	#closeIcon.hover(function(){
		closeIcon.src = "img/Login/ico_close_P.png";
	},function(){
		closeIcon.src = "img/Login/ico_close_popup.png";
	})
	_self.find(".left").click(function(){
		var showMonth = _self.find(".top-change-month").html();
		var showYear = _self.find(".top-change-year").html();
		if(showMonth == 1){
			showMonth = 12;
			showYear = parseInt(showYear) - 1;
			_self.find(".top-change-year").html(showYear);
		}else{
			showMonth = showMonth - 1;
		}
		_self.find(".change-bg").html(showMonth);
		showMonth = extra(showMonth);
		_self.find(".top-change-month").html(showMonth);
		panel();
		YearChange();   
		MonthChange();  
		Day();
	});
	
	_self.find(".right").click(function(){
		var showMonth = _self.find(".top-change-month").html();
		var showYear = _self.find(".top-change-year").html();
		if(showMonth == 12){
			showMonth = 1;
			showYear =  parseInt(showYear) + 1;
			_self.find(".top-change-year").html(showYear);
		}else{
			showMonth = parseInt(showMonth) + 1;
		}
		_self.find(".change-bg").html(showMonth);
		showMonth = extra(showMonth);
		_self.find(".top-change-month").html(showMonth);
		panel();
		YearChange();
		MonthChange();
		Day();
	});
	_self.find(".top-change-year").click(function(){
		panel();
		YearChange();
		_self.find(".month-panel").hide();
		_self.find(".detailChangePanel").hide();
		if(_self.find(".year-panel").is(":visible")){		
			_self.find(".contain").fadeIn();
			_self.find(".detailTime").fadeIn();
		}else{
			_self.find(".contain").fadeOut();
			_self.find(".detailTime").fadeOut();
		}
		_self.find(".year-panel").fadeToggle();
	});
	_self.find(".top-change-month").click(function(){
		panel();
		MonthChange();
		_self.find(".year-panel").hide();
		_self.find(".detailChangePanel").hide();
		if(_self.find(".month-panel").is(":hidden")){	
			_self.find(".contain").fadeOut();
			_self.find(".detailTime").fadeOut();
		}else{
			_self.find(".contain").fadeIn();
			_self.find(".detailTime").fadeIn();
		}		
		_self.find(".month-panel").fadeToggle();
	});
	// clear time
	_self.find(".clear").click(function(){
		_self.find(".detailTime").find("input[type='text']").val("00");
	})
	// hour min sec
	_self.find(".nowTime").click(function(){
		var changeDetailTime = new Date();
		nowHour = changeDetailTime.getHours();         //当前小时
		nowMin = changeDetailTime.getMinutes();        //当前分钟
		nowSec = changeDetailTime.getSeconds();        //当前秒
		_self.find(".top-change-month").html(nowMonth);
		_self.find(".top-change-year").html(nowYear);

		var nowMonthInt = _self.find(".top-change-month").text().replace(/[^0-9]/ig,"");  // 去0
		_self.find(".change-bg").html(nowMonthInt);
		Day();
		var H = extra(nowHour);
		var M = extra(nowMin);
		var S = extra(nowSec);
		#hour.val(H);
		#min.val(M);
		#sec.val(S);
	})
	#hour.click(function(){
		#setTime.toggle();
		setTime.trigger();
	})
	#min.click(function(){
		#setTime.toggle();
		setTime.trigger();
	})
	#sec.click(function(){
		#setTime.toggle();
		setTime.trigger();
	})
	#closeIcon.click(function(){
		#setTime.toggle();
	})
	$(".date-panel").click(function(){
		event.stopPropagation();       
	})
	//可编辑“十分秒”动态填充
	var changeHourListStr = "";
	var changeMinListStr = "";
	var changeSecListStr = "";
	for(var i=23;i>=0;i--){
		var zeroI = extra(i);
		changeHourListStr += "<li class='changeHourLi'>"+ zeroI +"</li>";
	}
	#changeHourList.find("ul").html(changeHourListStr); 
	for(var j=59;j>=0;j--){
		var zeroJ = extra(j);
		changeMinListStr += "<li class='changeMinLi'>"+ zeroJ +"</li>";
	}
	#changeMinList.find("ul").html(changeMinListStr); 
	for(var k=59;k>=0;k--){
		var zeroK = extra(k);
		changeSecListStr += "<li class='changeSecLi'>"+ zeroK +"</li>";
	}
	#changeSecList.find("ul").html(changeSecListStr); 
	_self.find(".changeHourLi").click(function(){
		#hour.val($(this).html());
		_self.find(".detailChangePanel").css("display","none");
	})
	_self.find(".changeMinLi").click(function(){
		#min.val($(this).html());
		_self.find(".detailChangePanel").css("display","none");
	})
	_self.find(".changeSecLi").click(function(){
		#sec.val($(this).html());
		_self.find(".detailChangePanel").css("display","none");
	});
	
	_self.on("click",".nowDay",function(){
		DATE = $(this).text();   // 日
		if(showMonth<10){
			showMonth = showMonth.substr(-1,1);
		}else{
			showMonth = showMonth;
		}
		if(onClick){
			onClick(@this);
		}
	});
}
/**
 * 月份补0
 */
function extra(x){  
    if(x < 10){ 
    	x = "0" + x
     }else{ 
     	x = x;
     }  
     return x;
} 
/**
 * 自定义年 + 月 面板填充
 */
function panel(){
	var showMonth = _self.find(".top-change-month").html();
	var showYear = _self.find(".top-change-year").html();
	var monthPanelStr = "";
	for(var i=1;i<13;i++){
		var titleShow = showYear+"-"+i;
		monthPanelStr +=   "<li><a href='#' class='monthChange' title='"+titleShow+"'>" + i + "月</a></li>";
	}
	_self.find(".monthPanelUl").html(monthPanelStr);      //自定义修改月份面板填充
	_self.find(".monthChange").each(function(){
		var monthInt = $(this).text().replace(/[^0-9]/ig,"");
		monthInt = extra(monthInt);
		if(monthInt == showMonth){
			$(this).addClass('changeMonthActive');
		}
	})
	var yearPanelStr = "";
	for(var i=1;i<6;i++){
		yearPanelStr +=   "<li><a href='#' class='yearChange'>" + (showYear-6+i) + "</a></li>";
	}
	yearPanelStr += "<li><a href='#' class='changeYearActive yearChange'>" + showYear + "</a></li>";
	for(var j=1;j<7;j++){
		yearPanelStr +=   "<li><a href='#' class='yearChange'>" + (parseInt(showYear)+j) + "</a></li>";
	}
	_self.find(".yearPanelUl").html(yearPanelStr);      //自定义修改年份面板填充
}
/**
 * 自定义年面板显隐控制和动态改变
 */
function YearChange(){
	_self.find(".yearChange").each(function(){
		$(this).click(function(){
			_self.find(".top-change-year").html($(this).text());
			_self.find(".year-panel").hide();
			Day();
			_self.find(".contain").fadeIn();
			_self.find(".detailTime").fadeIn();
		})
	})
}
/**
 * 自定义月面板显隐控制和动态改变
 */
function MonthChange(){
	_self.find(".monthChange").each(function(){
		$(this).click(function(){
			var monthInt = $(this).text().replace(/[^0-9]/ig,"");
			_self.find(".change-bg").html(monthInt);   
			monthInt = extra(monthInt);
			_self.find(".top-change-month").html(monthInt);
			_self.find(".month-panel").hide();
			Day();
			_self.find(".contain").fadeIn();
			_self.find(".detailTime").fadeIn();
		})
	})
}
/**
 * 动态填充当前日期面板
 */
function Day(){
	var showMonth = _self.find(".top-change-month").html();
	var showYear = _self.find(".top-change-year").html();
	var dateStr = "";                                                                         //天
	var new_date = new Date(showYear,showMonth-1,1);                                          //取当年当月中的第一天  
	var firstDay = new_date.getDay();                                                         //本月第一天是星期几 (月份 - 1)
	var new_date_next = new Date(showYear,showMonth,1);
	var prev_new_date_next = new Date(showYear,showMonth-1,1);
	var date_count = (new Date(new_date_next.getTime()-1000*60*60*24)).getDate();             //获取当月的天数 
	var prev_date_count = (new Date( prev_new_date_next.getTime()-1000*60*60*24)).getDate();  //获取上一个月的天数 
	var last_date = new Date(new_date_next.getTime()-1000*60*60*24);                          //获得当月最后一天的日期
	var lastDay = last_date.getDay();                                                         //最后一天星期几
    var beginDay = prev_date_count-firstDay+1;                                                //获取当年当月的第一个显示的天
    var lastShowDay = 42 - date_count - firstDay;                                             //计算当前月份有几个空需要下一个月来补
    var showMonthTitle = _self.find(".top-change-month").html(); 
	var showYearTitle = _self.find(".top-change-year").html();
	function prevDateFunction(){
		for(var j=0;j<parseInt(firstDay);j++){
			if(showMonthTitle == 1){
				showMonthTitle = 13;
				showYearTitle = parseInt(showYearTitle) - 1;
			}
			var title = showYearTitle + "-" + (parseInt(showMonthTitle)-1) + "-" + (beginDay + j);
			dateStr += "<li><a href='#' class='prevDay' title='"+title+"'>" + (beginDay + j) + "</a></li>";
		}
	}
	prevDateFunction();
	var showMonthTitle = _self.find(".top-change-month").html();
	var showYearTitle = _self.find(".top-change-year").html();
	for(var i=1;i<date_count+1;i++){
		var title = showYearTitle + "-" + showMonthTitle + "-" + i;
		dateStr += "<li><a href='#' class='nowDay' title='"+title+"'>"+ i +"</a></li>";
	}
	if(showMonthTitle == 12){
		showMonthTitle = 0;
		showYearTitle = parseInt(showYearTitle) + 1;
	}
	for(var k=1;k<lastShowDay+1;k++){
		var title = showYearTitle + "-" + (parseInt(showMonthTitle)+1) + "-" + k;
		dateStr += "<li><a href='#' class='nextDay' title='"+title+"'>"+ k +"</a></li>";
	}
	_self.find(".date-Ul").html(dateStr); 
	/**
	 * 激活当前天
	 */
	if(nowYear==showYear&&nowMonth==showMonth){
		_self.find(".nowDay").each(function(){
			if($(this).text()==nowDate){
				$(this).addClass('active');
			}
		})
	}
	
}  
public function remove(){
	_self.remove();
}

public function get value(){ 
	var hourVal = #hour.val();
	var minVal = #min.val();
	var secVal = #sec.val();
	showMonth = _self.find(".top-change-month").html();
	showYear = _self.find(".top-change-year").html();
	return showYear+"-"+showMonth+"-"+DATE;//+" "+hourVal+":"+minVal+":"+secVal;

}
public function set value(value:String):void{
	//TODO
	// return this.value;
}