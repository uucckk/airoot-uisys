package jus

import (
	"fmt"
)

func trace(value ...interface{}) {
	if len(value) != 1 {
		fmt.Println(value[0])
	} else {
		for i := 0; i < len(value); i++ {
			fmt.Println(value[i])
		}

	}

}

/**
 * 三目运算符
 */
func IfStr(isTrue bool, True string, False string) string {
	if isTrue {
		return True
	} else {
		return False
	}
}
