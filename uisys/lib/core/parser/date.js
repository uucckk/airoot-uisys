
Date.prototype._METHOD_ = {
	_label:"Date 方法说明",
	"format" : {value:"设置日期显示格式.例如 format('yyyy-MM-dd') 其显示内容为：2015-09-07"},
	"offsetDate" : {value:"设置日偏移量.例如 offsetDate(+1) 表示显示当天日期的下一天"},
	"offsetMonth" : {value:"设置日期显示格式.例如 format('yyyy-MM-dd') 其显示内容为：2015-09-07"},
	"offsetDay" : {value:"设置日期显示格式.例如 format('yyyy-MM-dd') 其显示内容为：2015-09-07"},
};
Date.prototype.format = function (fmt) { //author: meizz 
    var o = {
        "M+": this.getMonth() + 1, //月份 
        "d+": this.getDate(), //日 
        "h+": this.getHours(), //小时 
        "m+": this.getMinutes(), //分 
        "s+": this.getSeconds(), //秒 
        "q+": Math.floor((this.getMonth() + 3) / 3), //季度 
        "S": this.getMilliseconds() //毫秒 
    };
    if (/(y+)/.test(fmt)) fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
    for (var k in o)
    if (new RegExp("(" + k + ")").test(fmt)) fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
    return fmt;
}
Date.prototype.offsetDate = function(offset){
	this.setDate(this.getDate() + offset);
	return this;
}

Date.prototype.offsetYear = function(offset){
	this.setYear(this.getFullYear() + offset);
	return this;
}

Date.prototype.offsetMonth = function(offset){
	this.setMonth(this.getMonth() + offset);
	return this;
}

Date.prototype.offsetDay = function(offset){
	this.setDay(this.getDay() + offset);
	return this;
}


Date.prototype.Date = function(date){
	this.setDate(date);
	return this;
}


this.decodeScript = function(codeValue){
    	var code = '';
        jsdecoder = new JsDecoder();
        jscolorizer = new JsColorizer();

        jsdecoder.s = codeValue;

        code = jsdecoder.decode();
        
		
		code = code.replace(/&/g, "&amp;");
        code = code.replace(/</g, "&lt;");
        code = code.replace(/>/g, "&gt;");
        jscolorizer.s = code;
        try {
            code = jscolorizer.colorize();
        } catch (e) {
            $('msg').innerHTML += 'error<br><br>'+new String(e).replace(/\n/g, '<br>');
            return;
        }
       
        code = new String(code);
        code = code.replace(/(\r\n|\r|\n)/g, "<br>\n");
        code = code.replace(/<font\s+/gi, '<font@@@@@');
        code = code.replace(/( |\t)/g, '&nbsp;');
        code = code.replace(/<font@@@@@/gi, '<font ');

        code = code.replace(/\n$/, '');

        var count = 0;
        var pos = code.indexOf("\n");
        while (pos != -1) {
           count++;
           pos = code.indexOf("\n", pos+1);
        }
        count++;

        pad = new String(count).length;
        var lines = '';

        for (var i = 0; i < count; i++) {
            var p = pad - new String(i+1).length;
            var no = new String(i+1);
            for (k = 0; k < p; k++) { no = '&nbsp;'+no; }
            no += '&nbsp;';
            lines += '<div style="background: ' + '#333333' + '; color: #f0f0f0;margin-right:5px;">'+no+'</div>';
        }


        return "<table><tr><td>" + lines + "</td><td>" + code + "</td></tr></table>"; 
      
	}
	/**
 	 *	看看存不存在此对象
 	 * @param name		判断标示
 	 * @param obj		判断对象
 	 * @return			如果存在返回true,否则返回false;
 	 */
	this.HAVE = function (name,obj){
		if(name != null && document.getElementById(name) != undefined){
			alert("The Id [" + name + "] is exist.");
			return false;
		}
		try{
			obj = typeof(eval(obj))
			if(obj == "undefined"){
				return false;
			}
			alert("The Object " + name + "[" + obj + "] is exist.");
			return true;
		}catch(e){
			return false;
		}
	};
	
	
	/**
	 * 继承操作函数代码
	 */
	this.extendCode = function(father){
		var code = father.constructor.toString();
		var start = code.indexOf('{') + 1;
		var end = code.lastIndexOf('}');
		return code.substring(start,end);
	}//extendCode
	
	/**
	 * 继承父类的属性
	 */
	this.extend = function (dest,src){
		for (var p in src) { 
	        dest[p] = src[p]; 
	    } 
	    return dest;
	};
	
	//Pull out only certain bits from a very large integer, used to get the time
	//code information for the first part of a UUID. Will return zero's if there
	//aren't enough bits to shift where it needs to.
	var getIntegerBits = function(val,start,end){
	 var base16 = returnBase(val,16);
	 var quadArray = new Array();
	 var quadString = '';
	 var i = 0;
	 for(i=0;i<base16.length;i++){
	     quadArray.push(base16.substring(i,i+1));   
	 }
	 for(i=Math.floor(start/4);i<=Math.floor(end/4);i++){
	     if(!quadArray[i] || quadArray[i] == '') quadString += '0';
	     else quadString += quadArray[i];
	 }
	 return quadString;
	};
	
