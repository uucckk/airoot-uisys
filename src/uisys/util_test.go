// str.go
package jus

import (
	"flag"
	"fmt"
	"testing"
)

func TestIndex(t *testing.T) {
	//fmt.Println(ImportFrom("vue,{math,log}\002lib/js/vue.js;"))
	//fmt.Println(ImportFrom("{math,log},vue\002lib/js/vue.js;"))
}

var which = flag.Bool("which", true, "")
var path = flag.String("path", "", "")
var cnt = flag.Int("cnt", 100, "")

func TestIndex1(t *testing.T) {
	fmt.Println(f2md5("E:/juswork/highLight/lib/js/format.js"))
}
