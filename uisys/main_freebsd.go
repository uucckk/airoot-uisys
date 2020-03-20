package main

import (
	"fmt"
	"strconv"
)

//
func DevPrint(i int, value ...interface{}) string {
	out := ""
	switch i {
	case 0:
		out = "\033[32;1m" // Blue
	case 1:
		out = "\033[31;1m" // Red
	case 2:
		out = "\033[33;1m" // Yellow
	case 3:
		out = "\033[34;1m" // Green
	default:
		out = "\033[0m" // Default
	}
	str := fmt.Sprintf(out+value[0].(string)+"\033[0m", value[1:]...)
	fmt.Print(str)
	return "<span class='t" + strconv.Itoa(i) + "'>" + str + "</span>"
}

func SetConsoleTitle(title string) int {
	return 0
}
