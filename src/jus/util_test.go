// str.go
package jus

import (
	"fmt"
	"testing"
)

func TestIndex(t *testing.T) {
	fmt.Println(ImportFrom("vue,{math,log}\002lib/js/vue.js;"))
	fmt.Println(ImportFrom("{math,log},vue\002lib/js/vue.js;"))
}
