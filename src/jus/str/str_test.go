// str.go
package str

import (
	"fmt"
	"testing"
)

func TestIndex(t *testing.T) {
	//fmt.Println(Index("ssss", ""))
	lst := make([]string, 0)
	lst = append(lst, "DD")
	tlst := lst[0:0]

	tlst = append(tlst, "AA")
	fmt.Println(len(lst), lst[0], len(tlst))
}
