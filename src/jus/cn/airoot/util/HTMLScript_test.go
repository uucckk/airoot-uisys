// str.go
package util

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestIndex(t *testing.T) {
	path := `k/d/c/../a///b///c///d.txt`
	//path = filepath.FromSlash(path) // 平台处理

	d1 := filepath.Clean(path)

	fmt.Println(d1, filepath.Dir(path), filepath.Ext(path), filepath.Base(path))
}
