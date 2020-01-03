// str.go
package util

import (
	"fmt"
	"testing"
)

func TestIndex01(t *testing.T) {
	m := &MScript{}
	m.ReadFromString("{p:-100}")
	fmt.Println(m.ToString())
}
