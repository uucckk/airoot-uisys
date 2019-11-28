// client_cmd_test.go
package tool

import (
	"fmt"
	"testing"
)

func TestIndex(test *testing.T) {
	fmt.Println("start")
	list := &SocketList{}
	s1 := &SocketElement{Value: "1"}
	s2 := &SocketElement{Value: "2"}
	s3 := &SocketElement{Value: "3"}
	s4 := &SocketElement{Value: "4"}
	s5 := &SocketElement{Value: "5"}
	s6 := &SocketElement{Value: "6"}
	list.Append(s1)
	list.Append(s2)
	list.Append(s3)
	list.Append(s4)
	list.Append(s5)
	list.Append(s6)
	list.First()
	for {
		fmt.Println(list.Get().Value)
		if list.Next() == nil {
			break
		}
	}
	fmt.Println("==")
	list.Remove(s5)
	list.Remove(s3)
	list.Remove(s4)
	list.Remove(s2)
	list.Remove(s1)
	list.Remove(s6)
	fmt.Println("==")
	list.First()
	for {
		if list.Get() != nil {
			fmt.Println(list.Get().Value)
		}

		if list.Next() == nil {
			break
		}
	}

}
