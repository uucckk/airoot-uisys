package main

import (
	"fmt"
	"strconv"
	"syscall"
	"unsafe"
)

func DevPrint(i int, value ...interface{}) string {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("SetConsoleTextAttribute")            //SetConsoleTitle
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(i)) //12 Red light
	str := fmt.Sprintf(value[0].(string), value[1:]...)
	fmt.Print(str)
	CloseHandle := kernel32.NewProc("CloseHandle")
	CloseHandle.Call(handle)
	handle, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(7)) //White dark
	CloseHandle = kernel32.NewProc("CloseHandle")
	CloseHandle.Call(handle)
	return "<span class='t" + strconv.Itoa(i) + "'>" + str + "</span>"
}

func SetConsoleTitle(title string) int {
	kernel32, _ := syscall.LoadLibrary("kernel32.dll")
	_SetConsoleTitle, _ := syscall.GetProcAddress(kernel32, "SetConsoleTitleW")
	ret, _, callErr := syscall.Syscall(_SetConsoleTitle, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)

	if callErr != 0 {

		fmt.Println("callErr", callErr)

	}

	return int(ret)

}
