// message.go
//序列通话
//仅限于文本传输混淆
package nets

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"
)

//怎加了\0 和 \1 两个字段，一个负责结束，一个负责混淆
//例如 MA 1231\1INDEX\1123132123 sfsdfsfsdfsfddsf\0
//例如 MA 1231123132123 sfsdfsfs\1INDEX\1dfsfddsf\0
//例如 M\1INDEX\1A 1231123132123 sfsdfsfsdfsfddsf\0
type Message struct {
	Socket net.Conn //socket入口
	data   []byte
	buf    *bytes.Buffer
	sIndex int64 //本地序列
	dIndex int64 //对方序列
}

func (m *Message) Send(value string) (n int, err error) {
	if m.Socket != nil {
		return m.Socket.Write([]byte(value + "\x00"))
	} else {
		fmt.Println(value)

		return 0, errors.New("socket is nil")
	}

}

func (m *Message) Read() (string, error) {

	if m.data == nil {
		m.data = make([]byte, 1)
	}
	if m.buf == nil {
		m.buf = bytes.NewBufferString("")
	}
	m.buf.Reset()

	var c byte
	var count = 0
	for true {
		_, e := m.Socket.Read(m.data)
		if e != nil {
			return "", e
		}
		c = m.data[0]
		if c == 0 {
			break
		}
		count++
		if count > 65536 {
			fmt.Println("Message -> fulldeep")
		}
		m.buf.WriteByte(c)
	}
	return m.buf.String(), nil
}

func (m *Message) SetReadDeadline(t time.Time) error {
	return m.Socket.SetReadDeadline(t)
}

func (m *Message) Close() error {
	return m.Socket.Close()
}
