// str.go
package str

import (
	"fmt"
	"testing"
)

func TestIndex(t *testing.T) {
	a := FmtCmd("a 'c b' c")
	fmt.Println(a, len(a))
}
