// css.go
package util

import (
	"bytes"
	. "jus"
	"strings"
)

//CSS Class 元素
type ClassElement struct {
	Next        *ClassElement
	ElementType int
	Value       string
}

type Selecter struct {
	Element []*ClassElement
	Value   string
	index   int
	first   *ClassElement
	end     *ClassElement
}

func (s *Selecter) SetValue(value string) {
	s.Value = value
}

//新增一条规则
func (s *Selecter) AddRule(ce *ClassElement) {
	s.end = ce
	s.first = ce
	s.Element = append(s.Element, s.first)
}

//继续增加规则
func (s *Selecter) Push(ce *ClassElement) *Selecter {
	if s.first != nil {
		s.end.Next = ce
		s.end = ce
	} else {
		s.end = ce
		s.first = ce
		s.Element = append(s.Element, s.first)
	}
	return s
}

func (s *Selecter) NextRule() {
	s.first = nil
	s.end = nil
}

func (s *Selecter) At(position int) *Selecter {
	s.index = position
	s.first = s.Element[position]
	s.end = s.first.Next
	for s.end.Next != nil {
		s.end = s.end.Next
	}
	return s
}

func (s *Selecter) Length() int {
	return len(s.Element)
}

//-----------------------------------------------CSS----------------------------
var res [3]rune = [3]rune{'l', 'i', 'b'}

//转为JUS定制的CSS样式解析器
type CSS struct {
	selecter    []*Selecter
	jus         *JUS
	CurrentPath string
}

//从字符串里读CSS内容
func (c *CSS) ReadFromString(css string) {
	rp := 0
	code := []rune(css)
	position := 0
	var ch rune
	tag := make([]rune, 0)
	values := bytes.NewBufferString("")
	sel := &Selecter{index: -1}
	isValue := false
	lvl := 0 //中括号的数量
out:
	for position < len(code) {
		ch = code[position]
		position++
		if ch == '\r' || ch == '\n' || ch == '\t' {
			continue
		}
		if !isValue {
			if ch == '[' {
				lvl++
			} else if ch == ']' {
				lvl--
			}
		}

		if ch == ' ' {
			if len(tag) > 0 {
				sel.Push(&ClassElement{Value: string(tag), ElementType: 0})
				tag = tag[0:0]
			}
			continue
		}
		if ch == '>' {
			if len(tag) > 0 {
				sel.Push(&ClassElement{Value: string(tag), ElementType: 1})
				tag = tag[0:0]
			}
		}

		if ch == '.' {
			if len(tag) > 0 {
				sel.Push(&ClassElement{Value: string(tag), ElementType: 1})
				tag = tag[0:0]
			}
		}

		if ch == ':' {
			if len(tag) > 0 {
				sel.Push(&ClassElement{Value: string(tag), ElementType: 1})
				tag = tag[0:0]
			}
		}

		if ch == ',' {
			if len(tag) > 0 {
				sel.Push(&ClassElement{Value: string(tag), ElementType: -1}).NextRule()
				tag = tag[0:0]
			}
			continue
		}

		if ch == '{' && lvl == 0 {
			if len(tag) > 0 {
				sel.Push(&ClassElement{Value: string(tag), ElementType: -1})
				tag = tag[0:0]
			}
			lvl = 1
			values.Reset()
			values.WriteRune('{')
			for position < len(code) {
				ch = code[position]
				position++
				if ch == '{' {
					lvl++
				} else if ch == '}' {
					lvl--
				}
				if ch == '@' { //说明有特殊标识如lib
					rp = 0
					for position < len(code) {
						ch = code[position]
						position++

						if ch == res[rp] {
							rp++

							if rp == len(res) {
								ch = code[position]
								if !c.isChar(ch) {
									values.WriteString(c.CurrentPath)
									break
								}
							}
						} else {
							break
						}
					}
					continue
				}

				values.WriteRune(ch)
				if lvl == 0 {
					sel.Value = values.String()
					c.selecter = append(c.selecter, sel)
					sel = &Selecter{index: -1}
					continue out
				}
			}
		}

		tag = append(tag, ch)
	}

}

/**
 * 为所有类添加首域名
 * @param domain
 */
func (c *CSS) AddDomain(domain string) {
	var p []*ClassElement = nil
	var ce *ClassElement = nil
	for _, v := range c.selecter {
		p = v.Element
		for i := 0; i < len(p); i++ {
			ce = p[i]
			if strings.IndexRune(ce.Value, '$') != -1 || "body" == ce.Value {
				continue
			}
			nce := &ClassElement{Value: domain, ElementType: 0}
			nce.Next = ce
			p[i] = nce
		}
	}

}
func (c *CSS) isChar(ch rune) bool {
	if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '-' {
		return true
	}
	return false
}

/**
 * 替换属性
 * @param selector
 * @param value
 */
func (c *CSS) ReplaceSelecter(attr string, value string) {
	var p []*ClassElement = nil
	var ce *ClassElement = nil
	for _, v := range c.selecter {
		p = v.Element
		for i := 0; i < len(p); i++ {
			ce = p[i]
			if ce.Value == attr {
				ce.Value = value
			}
		}
	}

}

/**
 * 移除选择器
 * @param selecter
 */
func (c *CSS) RemoveSelecter(attr string) {

}

/**
 * 获取需要被转移控件Class
 * @return
 */
func (c *CSS) GetComponentClass() map[string]string {
	hm := make(map[string]string)
	var ce *ClassElement = nil
	var ch rune
	for _, s := range c.selecter {
		for _, a := range s.Element {
			ce = a
			for ce != nil {
				ch = []rune(ce.Value)[0]
				if (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '>' {

					if c.jus.GetPackageMap()[strings.ToLower(ce.Value)] != "" {
						hm[strings.ToLower(ce.Value)] = c.jus.GetDomain() + "-" + ce.Value
						ce.Value = "." + c.jus.GetDomain() + "-" + ce.Value

					}
				}
				ce = ce.Next
			}
		}
	}
	return hm
}

func (c *CSS) Length() int {
	return len(c.selecter)
}

/**
 * 转换为字符串
 */
func (c *CSS) ToString(tp int) string {
	sb := bytes.NewBufferString("")
	var ce *ClassElement = nil
	for _, s := range c.selecter {
		for _, a := range s.Element {
			ce = a
			for ce != nil {
				if c.jus == nil {
					sb.WriteString(ce.Value + IfStr(ce.ElementType == 0, " ", ""))
				} else {
					if ce.Value[0] == '#' {
						if tp == 0 { //共有属性
							sb.WriteString("[src_id='" + ce.Value[1:] + "']")
						} else { //私有属性
							if c.jus.GetDefine(ce.Value[1:]) == nil {
								sb.WriteString(ce.Value)
							} else {
								sb.WriteString("#" + c.jus.GetDefine(ce.Value[1:]).Name)
							}
						}

					} else {
						sb.WriteString(ce.Value)
					}
					sb.WriteString(IfStr(ce.ElementType == 0, " ", ""))
				}
				ce = ce.Next
			}
			sb.WriteString(",")
		}
		sb.Truncate(sb.Len() - 1) //sb = Substring(sb, 0, StringLen(sb)-1)
		sb.WriteString(s.Value)
		sb.WriteString("\n")
	}
	return sb.String()
}
