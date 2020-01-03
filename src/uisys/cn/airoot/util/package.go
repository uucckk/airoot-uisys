package util

import (
	"strings"
)

/**
 * 信息包
 */
type Package struct {
	from string
	data []byte
}

func (p *Package) router() string {
	for i, v := range p.data {
		if v == 0 {
			return string(p.data[0:i])
		}
	}
	return ""
}

func (p *Package) uuid() string {
	t := 0
	s := 0
	for i, v := range p.data {
		if v == 0 {
			if s == 1 {
				return string(p.data[t:i])
			}
			t = i
			s++
		}
	}
	return ""
}

func (p *Package) frame() string {
	t := 0
	s := 0
	for i, v := range p.data {
		if v == 0 {
			if s == 2 {
				return string(p.data[t:i])
			}
			t = i
			s++
		}
	}
	return ""
}

func (p *Package) value() string {
	s := 0
	for i, v := range p.data {
		if v == 0 {
			if s == 2 {
				if i == len(p.data)-1 {
					return ""
				}
				return string(p.data[i+1:])
			}
			s++
		}
	}
	return ""
}

/**
 * 或取信息主体
 */
func (p *Package) getDat() string {
	for i, v := range p.data {
		if v == 0 {
			return string(p.data[i:])
		}
	}
	return ""
}

/**
 * 转换给指定用户
 */
func (p *Package) ToUser(m map[string]*connectElement) int {
	r := p.router()
	if len(r) == 0 {
		return 0
	}
	d := []byte(p.from + p.getDat())

	if r[len(r)-1] == '*' { //批量广播
		//fmt.Println("批量广播")
		n := r[:len(r)-1]
		for k, v := range m {
			if k != p.from {
				if strings.Index(k, n) == 0 {
					v.Conn.Write(d)
				}
			}
		}

	} else {
		client := m[p.router()]
		if client != nil {
			client.Conn.Write(d)
		} else {
			return -1
		}

	}
	return 1
}
