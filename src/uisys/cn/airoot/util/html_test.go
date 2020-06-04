// str.go
package util

import (
	"fmt"
	"testing"
)

func TestIndex2(t *testing.T) {
	str := `
	<!-- Input -->
	<div>
		<span id="s" x="100" -Cmd></span>
		<Func("nihao",1234) data={x:10,y:10} />
	</div>
	<script>
		var list = [
			{name:"sunxy",pass:"123456"}
		];
		func init(){
			console.log(#s.getAttribute("x"));
		}
	</script>
	`
	html := &HTML{}
	html.ReadFromString(str)
	html.At(3).SetAttr("", "\"hello\",1234")

	html1 := &HTML{}
	html1.ReadFromString(str)

	html.At(3).Child()[1].ReplaceWith(html1.At(3))
	fmt.Println(">>", html.At(1).GetConstructerParameter())
	fmt.Println(">>", html.At(3).GetAttrCmd())
	fmt.Println(html.ToString())
	fmt.Println("------------------")
	fmt.Println(html.ToXHTML())
}
