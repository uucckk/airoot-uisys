// str.go
package str

import (
	"bytes"
	"strings"
)

/**
 * 获取字符串，字符穿从头的指定位置
 */
func Index(s string, value string) int {
	if len(value) == 0 {
		return -1
	}
	ch := []rune(value)
	l := len(ch)
	p := 0
	r := []rune(s)
	for i, v := range r {
		if v == ch[p] {
			p++
			if p == l {
				return i - p + 1
			}
		} else {
			p = 0
		}
	}
	return -1
}

func IndexRune(r []rune, value string) int {
	ch := []rune(value)
	l := len(ch)
	p := 0
	for i, v := range r {
		if v == ch[p] {
			p++
			if p == l {
				return i - p + 1
			}
		} else {
			p = 0
		}
	}
	return -1
}

/**
 * 从末尾获取某个位置
 */
func LastIndex(s string, value string) int {
	ch := []rune(value)
	l := len(ch)
	p := l - 1
	r := []rune(s)
	for i := len(r) - 1; i >= 0; i-- {
		if r[i] == ch[p] {
			p--
			if p < 0 {
				return i
			}
		} else {
			p = l - 1
		}
	}
	return -1
}

/**
 * 是否这个扩展名
 */
func IsType(s string, value string) bool {
	ch := []rune(value)
	l := len(ch)
	p := l - 1
	r := []rune(s)
	for i := len(r) - 1; i >= 0; i-- {
		if r[i] == ch[p] {
			p--
			if p < 0 {
				return true
			}
		} else {
			break
		}
	}
	return false
}

func LastPath(s string) string {
	r := []rune(s)
	for i := len(r) - 1; i >= 0; i-- {
		if r[i] == '\\' || r[i] == '/' {
			return string(r[(i + 1):])
		}
	}
	return s
}

func FormatSimplePath(s string) string {
	r := []rune(s)
	sb := make([]rune, 0)
	f := false
	var ch rune
	for i := 0; i < len(r); i++ {
		ch = r[i]
		if ch == '\\' || ch == '/' {
			if f == false {
				f = true
			} else {
				continue
			}
		} else {
			f = false
		}
		sb = append(sb, ch)
	}
	return string(sb)
}

/**
 * 截取字符串
 */
func Substring(s string, start int, end int) string {
	if start == -1 {
		return ""
	}
	ch := []rune(s)
	if end > -1 {
		return string(ch[start:end])
	}
	return string(ch[start:])

}

/**
 * 替换字符串
 */
func Replace(s string, old string, rep string) string {
	out := bytes.NewBufferString("")
	r := []rune(s)
	p := []rune(old)
	l := len(p)
	i := IndexRune(r, old)
	//t := 0
	if i < 0 {
		return s
	}
	for {
		if i >= 0 {
			//t = i
			out.WriteString(string(r[0:i]))
			out.WriteString(rep)
			r = r[(i + l):]
		} else {
			out.WriteString(string(r))
			break
		}
		i = IndexRune(r, old)

	}
	return out.String()
}

/**
 * 去除多余className外的内容，例如 "."
 */
func TrimClassName(s string) string {
	r := []rune(s)
	var ch rune
	i := 0
	for i = 0; i < len(r); i++ {
		ch = r[i]
		if ch == '.' {
			continue
		}
		break
	}
	return Substring(s, i, -1)
}

/**
 * 获取字符串指定索引内容
 */
func CharAt(s string, i int) string {
	c := []rune(s)
	if i < len(c) {
		return string(c[i])
	}
	return ""
}

/**
 * 获取字符串长度
 */
func StringLen(s string) int {
	return len([]rune(s))
}

/**
 * 格式化命令行
 */
func FmtCmd(s string) []string {
	tmp := make([]rune, 0)
	lst := make([]string, 0, 1)
	code := []rune(s)
	i := 0
	var ch rune
	str := ""
	f := false
	for i < len(code) {
		ch = code[i]
		if ch == '"' || ch == '\'' {
			if len(tmp) > 0 {
				lst = append(lst, string(tmp))
				tmp = tmp[0:0]
			}
			if f {
				str, i = readSimpleString(code, i)
				f = false
			} else {
				str, i = readString(code, i)
			}

			lst = append(lst, str)
			continue
		}
		if ch == ' ' || ch == '\t' {
			if len(tmp) > 0 {
				lst = append(lst, string(tmp))
				tmp = tmp[0:0]
			}
		} else {
			if ch == '@' && len(tmp) == 0 {
				f = true
			} else {
				tmp = append(tmp, ch)
			}

		}
		i++
	}
	if len(tmp) > 0 {
		lst = append(lst, string(tmp))
		tmp = tmp[0:0]
	}
	return lst
}

func FmtCmdList(s string) [][]string {
	lst := make([][]string, 0)
	if s == "" {
		return lst
	}
	arr := strings.Split(s, "\r\n")
	for _, v := range arr {
		lst = append(lst, FmtCmd(v))
	}
	return lst
}

/**
 * ReadString
 */
func readString(code []rune, position int) (string, int) {
	sb := make([]rune, 0)
	var t = code[position]
	position++
	var ch rune
	r := false
	for position < len(code) {
		ch = code[position]
		position++

		if ch == t && !r {
			break
		}
		if ch == '\\' {
			if r {
				r = false
			} else {
				r = true
				continue
			}
		} else {
			r = false
		}
		sb = append(sb, ch)
	}

	return string(sb), position
}
func readSimpleString(code []rune, position int) (string, int) {
	sb := make([]rune, 0)
	var t = code[position]
	position++
	var ch rune
	for position < len(code) {
		ch = code[position]
		position++
		if ch == t {
			break
		}
		sb = append(sb, ch)
	}
	return string(sb), position
}

/**
 * 转化成jus识别的$属性值
 */
func ToJUSString(value string) string {
	lst := bytes.NewBufferString("")
	r := []rune(value)
	l := len(r)
	p := 0
	fc := r[0]
	var ch rune
	flag := false
	for p < l {
		ch = r[p]
		p++
		if ch != '\\' && ch != '$' {
			flag = false
		} else if ch == '\\' {
			flag = !flag
		}
		if ch == '$' && flag == false {
			lst.WriteRune(fc)
			lst.WriteString("+__NAME__+")
			lst.WriteRune(fc)
			continue
		}
		lst.WriteRune(ch)
	}
	return lst.String()
}

type Cmd struct {
	Cmds []string          //命令list
	Attr map[string]string //命令挂载属性
}

/**
 * 格式化命令行
 */
func FmtCmdAdv(s string) *Cmd {
	cmd := &Cmd{}
	l := FmtCmd(s)
	c := make([]string, 0, 3)
	m := make(map[string]string)
	f := false
	var k, v string
	for _, t := range l {
		if t[0] == '-' {
			if f {
				v = ""
				m[k] = v
			}
			f = true
			k = t[1:]
		} else {
			if f {
				v = t
				f = false
				m[k] = v
			} else {
				c = append(c, t)
			}
		}
	}
	cmd.Cmds = c
	cmd.Attr = m
	return cmd

}

///代码注释工具
///用于代码测试
///@param code	输入代码
///@param flag	是否显示被注释的吗
func CodeFx(code string, flag bool) string {
	arr := strings.Split(code, "\n")
	s := bytes.NewBufferString("")
	f := true
	if flag {
		for _, v := range arr {
			if len(v) != 0 && (v[0] == '#' || v[0] == '*') {
				v = v[1:]
			}
			if f {
				s.WriteString(v)
				s.WriteByte('\n')
			}

		}
	} else {
		for _, v := range arr {
			if len(v) != 0 {
				if v[0] == '*' {
					f = !f
					continue
				}
				if v[0] == '#' {
					continue
				}
			}

			if f {
				s.WriteString(v)
				s.WriteByte('\n')
			}

		}
	}
	return s.String()
}
