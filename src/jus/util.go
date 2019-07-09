package jus

import (
	"bytes"
	"fmt"
	. "jus/str"
	. "jus/tool"
)

func trace(value ...interface{}) {
	if len(value) != 1 {
		fmt.Println(value[0])
	} else {
		for i := 0; i < len(value); i++ {
			fmt.Println(value[i])
		}

	}

}

/**
 * 三目运算符
 */
func IfStr(isTrue bool, True string, False string) string {
	if isTrue {
		return True
	} else {
		return False
	}
}

/**
 * 向数组中添加去重
 */
func Single(attr *[]*Attr, value *Attr) {
	p := Index(value.Value, "\002")
	if p != -1 {
		value.Value = Substring(value.Value, p+1, -1)
	}
	a := *attr
	for _, v := range a {
		if v.Name == value.Name {
			return
		}
	}
	*attr = append(*attr, value)
}

/**
 * 从数组里取数据
 */
func GetSingle(attr []*Attr, name string) *Attr {
	for _, v := range attr {
		if v.Name == name {
			return v
		}
	}
	return nil
}

/**
 * UMD 导入方法
 */
func ImportFrom(name string, value string) string {
	i := Index(value, "\002")
	buf := bytes.NewBufferString("")
	var k = Substring(value, 0, i)
	key := ""
	keys := make([]string, 0, 1)
	var isAttr = false
	for _, c := range k {
		if c == '{' {
			isAttr = true
			continue
		}
		if c == ',' || c == '}' {
			if isAttr {
				keys = append(keys, buf.String())
			} else {
				key = buf.String()
			}
			buf.Reset()
			if c == '}' {
				isAttr = false
			}
			continue
		}

		buf.WriteRune(c)
	}
	var p = Substring(value, i+1, -1)
	if buf.Len() != 0 {
		key = buf.String()
		buf.Reset()
	}
	buf.WriteString("var _UMD_ = __GET_UMD_LIB__('" + p + "','local');\n")
	buf.WriteString("var " + key + "=_UMD_['default'] ? _UMD_['default']  : _UMD_;\n")
	for _, v := range keys {
		buf.WriteString("var " + v + "=__FMT_UMD__(_UMD_,'" + v + "','" + name + "');\n")
	}

	return buf.String()
}
