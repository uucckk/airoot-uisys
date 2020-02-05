// socket_list
package tool

import (
	"net"
	"sync"
	"time"
)

type SocketMap struct {
	sync.RWMutex
	m map[net.Conn]*SocketInfo
}

func (s *SocketMap) Init() {
	s.Lock()
	s.m = make(map[net.Conn]*SocketInfo)
	s.Unlock()
}

func (s *SocketMap) Get(socket net.Conn) *SocketInfo {
	var info *SocketInfo
	s.Lock()
	info = s.m[socket]
	s.Unlock()
	return info
}

func (s *SocketMap) GetMap() map[net.Conn]*SocketInfo {
	return s.m
}

func (s *SocketMap) GetList() []*SocketInfo {
	s.Lock()
	list := make([]*SocketInfo, 0, 100)
	for _, v := range s.m {
		list = append(list, v)
	}
	s.Unlock()
	return list
}

func (s *SocketMap) Append(t *SocketInfo) {
	s.Lock()
	s.m[t.HostSocket] = t
	s.Unlock()
}

//移除
func (s *SocketMap) Remove(socket net.Conn) {
	s.Lock()
	delete(s.m, socket)
	s.Unlock()
}

type SocketList struct {
	sync.Mutex
	first *SocketElement
	last  *SocketElement
	pos   *SocketElement
}

func (s *SocketList) First() *SocketElement {
	s.pos = s.first
	return s.first
}

func (s *SocketList) Last() *SocketElement {
	s.pos = s.last
	return s.last
}

func (s *SocketList) Next() *SocketElement {
	if s.first == nil {
		return nil
	}
	if s.pos == nil {
		s.pos = s.first
	}
	if s.pos.next != nil {
		s.pos = s.pos.next
		return s.pos
	}

	return nil
}

func (s *SocketList) Prev() *SocketElement {
	if s.pos == nil {
		s.pos = s.last
	}
	if s.pos.prev != nil {
		s.pos = s.pos.prev
		return s.pos
	}
	return nil
}

func (s *SocketList) Get() *SocketElement {
	return s.pos
}

func (s *SocketList) Append(t *SocketElement) {
	if s.first == nil {
		s.first = t
		s.last = t
	} else {
		s.last.next = t
		t.prev = s.last
		s.last = t
		s.pos = t
	}
}

//移除
func (s *SocketList) Remove(t *SocketElement) *SocketElement {
	if t == nil {
		return nil
	}
	if t.prev == nil && t.next == nil {
		if s.first != nil {
			if s.first.Value == t.Value {

				s.first = nil
				return t
			} else {
				return nil
			}
		} else {
			return nil
		}
	} else if t.prev == nil { //说明是第一个
		s.first = nil
		t.prev = nil
		s.pos = nil
		return t
	} else if t.next == nil { //说明最后一个
		s.last = t.prev
		t.prev.next = nil
		t.prev = nil
		return t
	} else {
		t.prev.next = t.next
		t.next.prev = t.prev
		t.next = nil
		t.prev = nil
		return t
	}
}

type SocketElement struct {
	Value interface{}
	next  *SocketElement
	prev  *SocketElement
}

//ddns 记录各连接的状态
type SocketInfo struct {
	SendBytes    int       //发送数据
	ReceiveBytes int       //接收数据
	CreateTime   time.Time //创建时间
	HostSocket   net.Conn  //套接字
	Connect      bool
}
