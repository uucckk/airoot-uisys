//花生壳客户端
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

//动态转发服务器
type DDNSClient struct {
	Name   string
	Host   string
	Client string
}

//初始化客户端
func (d *DDNSClient) connectHost() {
	conn, error := net.Dial("tcp", "127.0.0.1:8080")
	if error != nil {
		fmt.Println("Connect Error:", error)
	}
	go func() {
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil {
			return
		}
		fmt.Println("01")
		client, err := d.getClient(conn)
		fmt.Println("02")
		if err != nil {
			return
		}
		client.Write(data[0:n])
		for {
			fmt.Println("start...")
			n, err := conn.Read(data)
			if err != nil {
				break
			}
			client.Write(data[0:n])
		}
	}()
}

func (d *DDNSClient) getClient(host net.Conn) (net.Conn, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:80")
	go func() {
		b := make([]byte, 1024)
		for true {
			n, e := conn.Read(b)
			if e != nil {
				break
			}
			host.Write(b[0:n])
		}
	}()
	return conn, err
}

//开始服务
func (d *DDNSClient) Start() {
	d.connectHost()
}
func main() {
	fmt.Println("hehe")
	ddns := &DDNSClient{}
	ddns.Start()
	buf := bufio.NewReader(os.Stdin)
	buf.ReadByte()
}
