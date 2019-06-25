// TestServer.go
//测试服务器连接的内容
package util

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type TestServer struct {
	Time       int64
	useCount   int
	loadSize   int
	totalSize  int
	Name       string
	listen     net.Listener
	From       string
	To         string
	FromClient net.Conn
	ToClient   net.Conn
	log        string
	logFrom    *os.File
	logTo      *os.File
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
					break
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
	t.useCount += 2
	count := 1024
	if dest != "" {
		conn, error := net.Dial("tcp", dest)
		if error == nil {
			t.FromClient = socket
			t.ToClient = conn
			go func() {
				data := make([]byte, count)
				for {
					fmt.Println("start...")
					n, err := conn.Read(data)
					if err != nil {
						break
					}
					if t.log != "" {
						if t.loadSize+n > t.totalSize {
							n = t.totalSize - t.loadSize
						}
						t.loadSize += n
						if n > 0 {
							t.logFrom.Write(data[0:n])
						}
					}
					socket.Write(data[0:n])
				}
				fmt.Println(t.From + ">" + t.To + ": [" + t.To + "] is release")
				t.useCount--
				socket.Close()
			}()

			go func() {
				data := make([]byte, count)
				for {
					n, err := socket.Read(data)
					if err != nil {
						break
					}
					if t.log != "" {
						if t.loadSize+n > t.totalSize {
							n = t.totalSize - t.loadSize
						}
						t.loadSize += n
						if n > 0 {
							t.logFrom.Write(data[0:n])
						}
					}
					conn.Write(data[0:n])
				}
				fmt.Println(t.From + ">" + t.To + ": [" + t.From + "] is release")
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
	os.MkdirAll(path, 777)
	var err error
	date := time.Now().Format("2006-01-02_150405")
	if t.logFrom == nil {
		t.logFrom, err = os.Create(path + "/" + date + "_" + "from.log")
		if err != nil {
			fmt.Println(err)
		}
	}

	if t.logTo == nil {
		t.logTo, err = os.Create(path + "/" + date + "_" + "to.log")
		if err != nil {
			fmt.Println(err)
		}
	}

}

/**
 * 关闭程序
 */
func (t *TestServer) Shutdown() {
	if t.listen != nil {
		t.listen.Close()
	}

	if t.logFrom != nil {
		t.logFrom.Close()
		t.logFrom = nil
	}
	if t.logTo != nil {
		t.logTo.Close()
		t.logTo = nil
	}
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
		return strconv.Itoa(t.loadSize) + "/" + strconv.Itoa(t.totalSize)
	}

}
