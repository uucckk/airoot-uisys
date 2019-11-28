// str.go
package str

import (
	"fmt"
	"testing"

	"github.com/satori/go.uuid"
)

func TestIndex(t *testing.T) {
	a := FmtCmdAdv("vp node -host :10001 -forward :10002")
	fmt.Println(a)
	for _, v := range a.Cmds {
		fmt.Print(v, ",")
	}
	fmt.Println()
	for k, v := range a.Attr {
		fmt.Println(k, ":", v)
	}
	uid, _ := uuid.NewV4()
	fmt.Println(len(uid.String()))
}
