// str.go
package util

import (
	"fmt"
	"testing"
)

func TestIndex2(t *testing.T) {
	str := `
	<!-- 
		类注释
		@author sunxy
		@version 0.0
	 -->
	<@pub/>
	<@import value="jus.*" />
	<head>
		<style>
			body{
				margin:0px;
			}
		</style>
	</head>
	<desktop>
		<style>
			body{
				text-align:center;
	          	
			}
	      	.content{
	        	background-color:#fefefe;  
	      	}
	      	.login{
	        	color:#f0f0f0;
	          	font-size:64px;
	          	font-weight:bold;
	      	}
	      
		</style>
	  	<position class="content" width="100%" height="100%">
			<hbox left="100" right="100" top="100" bottom="100">
				<jus.Coder width="100%" height="100%" debug="true" left="100" top="100" style="border:1px solid #dddddd">
					<@uncare>
						<style type="text/css">
							table {
								font-family: verdana,arial,sans-serif;
								font-size:11px;
								color:#333333;
								border-width: 1px;
								border-color: #a9c6c9;
								border-collapse: collapse;
							}
							table th {
								border-width: 1px;
								padding: 8px;
								border-style: solid;
								border-color: #a9c6c9;
							}
							table td {
								border-width: 1px;
								padding: 8px;
								border-style: solid;
								border-color: #a9c6c9;
							}
						</style>
						<div>
							<table>
								<@content/>
							</table>
						</div>
					</@uncare>
				</jus.Code>
				<jus.Coder width="100%" height="100%" debug="true" left="600" top="100" style="border:1px solid #dddddd">
					<@uncare>
						<DataGrid>
							<tr>
								<th>用户名</th><th>密码</th>
							</tr>
							<tr>
								<td>用户名</td><td>密码</td>
							</tr>
						</DataGrid>
						我爱背景天安们
						<div>11</div>
						<script>
							function init(){
							
							}
						</script>
						dddd
					</@uncare>
				</Code>
				<jus.Coder width="100%" height="100%" debug="true" left="600" top="100" style="border:1px solid #dddddd">
					<@uncare>
						private var _value = 0;
						/**
						 * 构造函数
						 */
						function init(){
							var i = 0;
							console.log(i*5);
						}
						//比例设置
						set value(v:int):void{
							_value = v;
						}
						//获得比例
						get value():int{
							return _value;
						}
					</@uncare>
				</Code>
				<!--
				<HTMLLoader url="http://www.baidu.com" width="100%" height="100%" left="1200" top="100" style="border:1px solid #dddddd"/>
				-->
			</hbox>
	  	</position>
	  
		<script>
	
			/**
			 * 默认初始化函数
			 */
			function init():void{
	
			}
	
			/**
			 * 析构函数
			 * 默认情况下可以不使用
			 */
			public function finalize():void{
	
			}
		</script>
	</desktop>
	`
	html := &HTML{}
	html.ReadFromString(str)
	ss := html.GetElementsByTagName("script")
	for _, v := range ss {
		v.Remove()
	}
	fmt.Println(html.ToString())
}
