package tool

import (
	"fmt"
)

type Console struct {
	Name string
}

func (c *Console) Log(value ...interface{}) {
	fmt.Println(c.Name+">>", value)
}
