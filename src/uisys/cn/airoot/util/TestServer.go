// TestServer.go
//测试服务器连接的内容
package util

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TestServer struct {
	sync.RWMutex
	Time       int64
	useCount   int
	downSize   int
	upSize     int
	totalSize  int
	Name       string
	listen     net.Listener
	From       string
	To         string
	FromClient net.Conn
	ToClient   net.Conn
	log        string
	count      int
}

/**
 * IP地址
 */
func (t *TestServer) FromIPAddress() string {
	return t.From
}

func (t *TestServer) Running() bool {
	if t.listen == nil {
		return false
	} else {
		return true
	}
}

/**
 * IP地址
 */
func (t *TestServer) ToIPAddress() string {
	return t.To
}

/**
 * 启动监听程序
 */
func (t *TestServer) Start(src string, dest string) bool {
	t.From = src
	t.To = dest
	t.totalSize = 1024 * 1024 * 1024
	return t.initSocket()
}

func (t *TestServer) initSocket() bool {
	listen, err := net.Listen("tcp", t.From)
	if err != nil {
		fmt.Println(">>", err)
	} else {
		t.listen = listen
		go func() {

			for {
				socket, err := listen.Accept()
				if err != nil {
					fmt.Println(">>", err)
					t.listen = nil
					if e := listen.Close(); e != nil {
						fmt.Println(t.Name+" Close havs error: ", e)
					}
					continue
				}
				t.Client(socket, t.To)
			}
		}()
		return true
	}
	return false
}

/**
 * 重新启动
 */
func (t *TestServer) Restart() bool {
	if t.listen == nil {
		return t.initSocket()
	}
	return false
}

/**
 *
 */
func (t *TestServer) Client(socket net.Conn, dest string) {
	os.MkdirAll(t.log, 777)
	var logDown *os.File
	var logUp *os.File
	var err error
	ip := socket.RemoteAddr().String()
	vip := strings.ReplaceAll(ip, ":", "-P")
	if t.log != "" {
		date := time.Now().Format("2006-01-02_150405") + "_" + vip + "_" + strconv.Itoa(t.getCount())
		logDown, err = os.Create(t.log + "/" + date + "_" + "down.log")
		if err != nil {
			fmt.Println(err)
		}
		logUp, err = os.Create(t.log + "/" + date + "_" + "up.log")
		if err != nil {
			fmt.Println(err)
		}
	}
	t.useCount += 2
	count := 1024
	if dest != "" {
		conn, error := net.Dial("tcp", dest)
		if error == nil {
			t.FromClient = socket
			t.ToClient = conn
			go func() {
				data := make([]byte, count)
				flag := true
				for flag {
					n, err := conn.Read(data)
					if err != nil {
						break
					}
					if t.log != "" {
						if t.downSize+n > t.totalSize {
							n = t.totalSize - t.downSize
							flag = false
						}
						t.downSize += n
						if n > 0 {
							logDown.Write(data[0:n])
						}
					}
					socket.Write(data[0:n])
				}
				if flag {
					fmt.Println(t.From + ">" + t.To + " @IP:" + ip + " is release")
				} else {
					fmt.Println("The downSize more than totalSize and auto release")
				}
				logDown.Close()
				t.useCount--
				socket.Close()
			}()

			go func() {
				data := make([]byte, count)
				flag := true
				for flag {
					n, err := socket.Read(data)
					if err != nil {
						break
					}
					if t.log != "" {
						if t.upSize+n > t.totalSize {
							n = t.totalSize - t.upSize
							flag = false
						}
						t.upSize += n
						if n > 0 {
							logUp.Write(data[0:n])
						}
					}
					conn.Write(data[0:n])
				}
				if flag {
					fmt.Println(t.From + ">" + t.To + " @IP:" + ip + " is release")
				} else {
					fmt.Println("The upSize more than totalSize and auto release")
				}
				logUp.Close()
				t.useCount--
				conn.Close()
			}()

		} else {
			fmt.Println("Connect Error:", error)
			socket.Close()
		}
	} else {
		go func() {
			data := make([]byte, count)
			for {
				n, err := socket.Read(data)
				if err != nil {
					break
				}
				fmt.Println(string(data[0:n]))
			}
			socket.Close()
			fmt.Println("socket over.")
		}()
	}

}

func (t *TestServer) SetLog(path string, size int) {
	t.log = path
	t.totalSize = size

}

/**
 * 关闭程序
 */
func (t *TestServer) Stop() {
	if t.listen != nil {
		t.listen.Close()
	}

}

func (t *TestServer) getCount() int {
	t.Lock()
	t.count++
	t.Unlock()
	return t.count
}

func (t *TestServer) GetLogPath() string {
	return t.log
}

func (t *TestServer) ConnectStatus() string {
	return strconv.Itoa(t.useCount / 2)
}

func (t *TestServer) LogStatus() string {
	if t.log == "" {
		return ""
	} else {
		return strconv.Itoa(t.downSize) + "/" + strconv.Itoa(t.totalSize) + " | " + strconv.Itoa(t.upSize) + "/" + strconv.Itoa(t.totalSize)
	}

}
