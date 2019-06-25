package main

import (
	"fmt"
	"image/color"
	"math"
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

func drawImage(a []int, width int) {
	//var handle uintptr
	//kernel32 := syscall.NewLazyDLL("kernel32.dll")
	//proc := kernel32.NewProc("SetConsoleTextAttribute") //SetConsoleTitle
	count := 0
	row := 0
	for i := 0; i < len(a); i += 4 {
		c := toOneColor(a[i], a[i+1], a[i+2], a[i+3])
		//handle, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(c)) //12 Red light
		fmt.Print(c)

		if count > 0 && (count+1)%width == 0 {
			fmt.Println()
			row++
			if row > 9 {
				break
			}
		}
		count++
	}

	//handle, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(7)) //White dark
	//CloseHandle := kernel32.NewProc("CloseHandle")
	//CloseHandle.Call(handle)
}
func toOneColor(R, G, B, A int) string {

	t := ""
	if R < 50 {
		t = "  "
	} else if R <= 255 {
		t = " \""
	}
	return t

}

func _drawImage(a []color.RGBA, width int) {
	var handle uintptr
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("SetConsoleTextAttribute") //SetConsoleTitle
	count := 0
	for i := 1; i < len(a); i++ {
		//c := _toOneColor(a[i].R, a[i].G, a[i].B, a[i].A)
		//handle, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(c)) //12 Red light
		fmt.Print("/")

		if count > 0 && (count+1)%width == 0 {
			fmt.Println()
		}
		count++
	}

	handle, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(7)) //White dark
	CloseHandle := kernel32.NewProc("CloseHandle")
	CloseHandle.Call(handle)
}

func _toOneColor(R, G, B, A uint8) uint8 {
	tp := 0
	var t uint8

	if math.Abs(float64(R-G)) < 10 && math.Abs(float64(R-B)) < 10 && math.Abs(float64(G-B)) < 10 { //rgb
		t = R
		if t < 50 {
			t = 0
		} else if t < 120 {
			t = 112
		} else if t < 200 {
			t = 128
		} else if t <= 255 {
			t = 240
		}
		return t
	}

	if R > G {
		t = R
		tp = 0
	} else {
		t = G
		tp = 1
	}

	if t < B {
		t = B
		tp = 2
	}
	if tp == 0 { //r
		if t < 64 {
			t = 112
		} else if t < 192 {
			t = 128
		} else if t < 255 {
			t = 204
		}
	}
	if tp == 1 { //g
		if t < 32 {
			t = 112
		} else if t < 200 {
			t = 160
		} else if t < 255 {
			t = 175
		}
	}
	if tp == 2 { //b
		if t < 31 {
			t = 112
		} else if t < 200 {
			t = 144
		} else if t < 255 {
			t = 159
		}
	}

	return t
}
