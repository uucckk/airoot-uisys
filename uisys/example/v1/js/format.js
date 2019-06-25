var keyWords = new Object();
var buf = "";
var code = null;
var position = 0;
addEventListener("message",function(e){
	var lineNum = 0;
	if(e.data instanceof Array){
		keyWords = e.data;
		return;
	}
	
	code = e.data;
	var tag = ""
	var ch;
	var length = code.length;
	while(position<length){
		ch = code.charAt(position ++);
		
		if(ch == '"' || ch == "'"){
			position --;
			buf += readString();
			continue;
		}
		
		if(ch == '!' || ch == '.' || ch == '\r' || ch == '\n' || ch == ' ' || ch == '\t' || ch == '{' || ch == '}' || ch == '[' || ch == ']' || ch == '(' || ch == ')' || ch == ';' || ch == ':' || ch == ',' || ch == '?' || ch == '>' || ch == '=' || ch == '<' || ch == '&' || ch == '|' || ch == '%' || ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '#' ){
			if(tag.length>0){
				buf += formatWord(tag);
				tag = "";
			}
			if(ch == '!' || ch == '.' || ch == '\r' || ch == '\n' || ch == '{' || ch == '}' || ch == '(' || ch == ')' || ch == '[' || ch == ']' || ch == ';' || ch == ':' || ch == ',' || ch == '?' || ch == '>' || ch == '=' || ch == '<' || ch == '&' || ch == '|' || ch == '%' || ch == ' ' || ch == '\t' || ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '#'){
				if(ch == '\n'){
					buf += "<br/>";
					++ lineNum;
				}else if(ch == '\t'){
					buf += "<span islock='true' type='TAB' style='color: transparent;'>____</span>"
					
				}else if(ch == ' '){
					buf += "<span islock='true' type='SPACE' style='color: transparent;'>_</span>"
				}
				
				else{
					if(ch == '\r'){
						continue;
					}
					buf += "<span class='word' type='SYMBOL' style='color:#234678;font-weight:bold;'>" + ch + "</span>";
					
				}
				
			}
			
			continue;
		}
		
		
		tag += ch;
		
	}
	
	if(tag.length>0){
		buf += "<span class='word'>" + tag + "</span>";
	}
	
	
	postMessage({data:buf,count:lineNum});
});


function formatWord(value){
	for(var i = 0;i<keyWords.length;i++){
		if(keyWords[i].name == value){
			return "<span class='word' style='color:" + keyWords[i].color + "' type='KEYWORD'>" + value + "</span>";
		}
	}
	return "<span class='word'>" + value + "</span>";
	
}


/**
 * 读取字符串或正则表达式
 * @param code
 * @return
 */
function readString(){
	var sb = "";
	var t = code.charAt(position ++);
	var ch;
	var r = false;
	sb += t;
	while(position<code.length){
		ch = code.charAt(position ++);
		sb += ch;
		if(ch == t && !r){
			break;
		}
		if(ch == '\\'){
			if(r){
				r = false;
			}else{
				r = true;
			}
		}else{
			r = false;
		}
	}
	
	return "<span class='word' style='color:#ffaa00'>" + sb.replace("\r\n", t + " + " + t) + "</span>";
}