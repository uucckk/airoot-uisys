package util

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"
	. "uisys/str"
	. "uisys/tool"
)

//特殊关键字
var keyWords = [...]string{"script", "style", "css", "~script", "@uncare"}
var closeTag = [...]string{"input", "img", "br", "meta", "hr", "source", "link"}

//判断是不是关键字
func isKeyWord(value string) bool {
	for _, v := range keyWords {
		if v == value {
			return true
		}
	}
	return false
}

//判断是否为自关闭标签
func isCloseTag(value string) bool {
	for _, v := range closeTag {
		if v == value {
			return true
		}
	}
	return false
}

type StringBuffer []rune

//返回此对象的String
func (p StringBuffer) toString() string {
	return string(p)
}

/*
	HTML
*/
type HTML struct {
	parent  *HTML          //父级元素
	tag     string         //标签类型
	value   string         //字符串实际值
	param   string         //节点构造函数参数
	code    string         //节点初始化代码
	tagData map[string]*Ch //内部属性
	tagList []string       //内部属性列表，方便排序
	tagType int            //HTML结束类型
	list    []*HTML        //内部的HTML列表
}

//初始化Tag
func (h *HTML) fx(code []rune, position int) (map[string]*Ch, []string, int) {
	lst := []*Ch{}
	var key string
	eq := false
	var ch rune
	var block string
	var tagData = make(map[string]*Ch)
	var tagList = make([]string, 0)
	tp := 0
	tmp := []rune{}
	for position < len(code) { //整理元素，去掉不必要的空格
		ch = code[position]
		if ch == '/' {
			if eq == true {
				block, position = h.fxr(code, position)
				tagData[key] = &Ch{block, 4}
				key = ""
				eq = false
			}
			position++
			continue
		}
		if ch == '>' {
			if len(tmp) > 0 {
				if key == "" {
					tagData[string(tmp)] = nil
				} else {
					tagData[key] = &Ch{string(tmp), 0}
				}

				tmp = tmp[0:0]
			} else if key != "" {
				tagData[key] = nil
			}

			return tagData, tagList, position
		}
		if ch == '(' || ch == '[' || ch == '{' {
			if len(tmp) > 0 {
				lst = append(lst, &Ch{string(tmp), 0})
				tmp = tmp[0:0]
			}
			switch ch {
			case '(':
				block, position = h.fxa(code, position, ch, ')')
				tp = 1
			case '[':
				block, position = h.fxa(code, position, ch, ']')
				tp = 2
			case '{':
				block, position = h.fxa(code, position, ch, '}')
				tp = 3
			}
			if key == "" {
				tagData[key] = &Ch{block, -1}
			} else {
				tagData[key] = &Ch{block, tp}
			}
			key = ""
			eq = false
			position++
			continue
		}

		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '=' {
			if len(tmp) > 0 {
				if eq {
					tagData[key] = &Ch{string(tmp), 0}
					key = ""
					eq = false
				} else {
					if key != "" {
						tagData[key] = nil
						tagList = append(tagList, key)
					}
					key = string(tmp)
					tagList = append(tagList, key)

				}
				tmp = tmp[0:0]
			}

			if ch == '=' {
				eq = true
			}

			position++
			continue
		}

		if ch == '"' || ch == '\'' {
			block, position = h.fxs(code, position, ch)
			tagData[key] = &Ch{block, 0}
			key = ""
			eq = false
			position++
			continue
		}
		tmp = append(tmp, ch)
		position++
	}

	if len(tmp) > 0 {
		lst = append(lst, &Ch{string(tmp), 0})
	}

	return tagData, tagList, position
}

///处理正则表达式
func (h *HTML) fxr(code []rune, position int) (string, int) {
	f := false
	var v rune
	position++
	sb := bytes.NewBufferString("/")
	for position < len(code) {
		v = code[position]
		sb.WriteRune(v)
		if !f && v == '/' {
			break
		}
		if v == '\\' {
			f = !f
		} else {
			f = false
		}

		position++
	}
	position++
	for position < len(code) {
		v = code[position]
		if v == ' ' || v == '/' || v == '>' {
			break
		}
		sb.WriteRune(v)
		position++
	}
	return sb.String(), position
}

///分析括号变量作用域
func (h *HTML) fxa(code []rune, position int, s rune, e rune) (string, int) {
	position++
	var block string
	lvl := 1
	sb := bytes.NewBufferString("")
	var v rune
	for position < len(code) {
		v = code[position]
		if v == s {
			lvl++
		} else if v == e {
			lvl--
		}
		if v == e && lvl == 0 {
			break
		}
		if v == '"' || v == '\'' {
			block, position = h.fxs(code, position, v)
			sb.WriteRune('"')
			sb.WriteString(block)
			sb.WriteRune('"')
			position++
			continue
		}
		sb.WriteRune(v)
		position++
	}
	return sb.String(), position
}

///分析字符串作用域
func (h *HTML) fxs(code []rune, position int, ch rune) (string, int) {
	position++
	tch := ch
	var zy bool = false //转义符号
	tmp := []rune{}
	for position < len(code) {
		ch = code[position]

		if ch == tch && zy == false {
			break
		}

		if ch == '"' && zy == false {
			tmp = append(tmp, '\\')
			tmp = append(tmp, '"')
		} else {
			tmp = append(tmp, ch)
		}

		if ch == '\\' {
			zy = !zy
		} else {
			zy = false
		}
		position++
	}
	return string(tmp), position
}

/**
 * 从字符串中获取HTML
 */
func (h *HTML) ReadFromString(value string) (*HTML, error) {
	html, _, err := h.read([]rune(value), -1)
	return html, err
}

/**
 * 从文字序列里读取一块HTML内容
 * 内容是以<标签开始的一段HTML内容
 */
func (h *HTML) ReadOneBlock(code []rune, index int) (*HTML, int, error) {
	return h.read(code, index)
}

func (h *HTML) read(code []rune, index int) (*HTML, int, error) {
	h.list = make([]*HTML, 0, 100)
	h.tagData = make(map[string]*Ch, 20)
	position := 0
	if index != -1 {
		position = index
	}
	sb := StringBuffer{}
	var ch rune
	var tag *HTML
	var tagName string
	var parent *HTML = h
	var tagType int = 0 //HTML的类型
	var block int = 0
	tagTemp := make([]string, 0, 100) //tag临时储存位置，用于记录标签配对问题

	var tagData map[string]*Ch
	var tagList []string

m:
	for position < len(code) {
		ch = code[position]
		position++
		if ch == '<' {
			//tagName
			for position < len(code) {
				ch = code[position]
				position++
				if ch == '(' || ch == '{' || ch == ' ' || ch == '!' || ch == '>' || (ch == '/' && code[position-2] != '<') {
					if ch == '!' {
						k := 0
						sb = sb[0:0]
						var keys []rune
						if code[position] == '[' {
							keys = []rune("]]>")
						} else {
							keys = []rune("-->")
						}

						for position < len(code) {
							ch = code[position]
							position++
							sb = append(sb, ch)
							if keys[k] == ch {
								k++
								if k == len(keys) {
									sb = sb[:(len(sb) - k + 2)]
									parent.list = append(parent.list, &HTML{tag: "!", value: sb.toString(), parent: parent, tagType: 0, tagData: make(map[string]*Ch, 20)})
									sb = sb[0:0] //清除
									tagName = ""
									block--
									break
								}
							} else {
								k = 0
							}
						}
						continue m
					} else {
						tagName = string(sb)
						tagData, tagList, position = h.fx(code, position-1)
					}

					sb = sb[0:0]
					break
				}

				sb = append(sb, ch)

			}
			if tagName == "" {
				return h, position, errors.New("html.go -> tagName is empty.")
			}
			for position < len(code) {
				ch = code[position]
				position++
				if ch == '>' {
					if code[position-2] == '/' || tagName == "!" {
						tagType = 0 //in one
					} else {
						tagType = 1 //start
					}
					break
				}
				sb = append(sb, ch)

			}

			if tagName[0] == '/' {
				if parent.parent != nil {
					parent = parent.parent
					sb = sb[0:0] //清除
					block--
					if block == 0 && index != -1 {
						return h, position, nil
					}
				}

			} else {
				if isCloseTag(tagName) { //判断是否为自关闭标签
					tagType = 0
				}
				tagTemp = append(tagTemp, tagName)
				tag = &HTML{tag: tagName, value: sb.toString(), tagData: tagData, tagList: tagList, parent: parent, tagType: tagType}
				parent.list = append(parent.list, tag)
				parent = tag
				sb = sb[0:0] //清除
				block++
				if isKeyWord(tagName) {
					k := 0
					keys := []rune("</" + tagName + ".")
					for position < len(code) {
						ch = code[position]
						position++
						sb = append(sb, ch)
						if keys[k] == ch || (keys[k] == '.' && ch == ' ') {
							k++
						} else {
							if k > 1 && ch == '>' {
								sb = sb[:(len(sb) - k - 1)]
								parent.list = append(parent.list, &HTML{value: sb.toString(), parent: parent, tagType: -1, tagData: make(map[string]*Ch, 20)})
								sb = sb[0:0] //清除
								parent = parent.parent
								block--
								break
							}
							k = 0
						}
					}
				} else if tagType == 0 {
					parent = parent.parent
					block--
					if block == 0 && index != -1 {
						return h, position, nil
					}
				}
			}

		} else { //文字
			sb = sb[0:0]
			sb = append(sb, ch)
			for position < len(code) {
				ch = code[position]
				if ch == '<' {
					break
				}
				position++
				sb = append(sb, ch)
			}
			if len(sb) != 0 {
				parent.list = append(parent.list, &HTML{value: sb.toString(), parent: parent, tagType: -1, tagData: make(map[string]*Ch, 20)})
				sb = sb[0:0] //清除

			}

		}

	}

	//变换为HTML

	return h, position, nil
}

//返回标签名称
func (h *HTML) TagName() string {
	return h.tag
}

//设置节点名称
func (h *HTML) SetTagName(value string) {
	h.tag = value
}

//返回HTML的属性值
func (h *HTML) GetAttr(attrName string) string {
	v := h.tagData[attrName]
	if v != nil {
		return v.Value
	} else {
		return ""
	}

}

//返回HTML的属性值
func (h *HTML) GetAttrCmd() []string {
	attrs := make([]string, 0)
	for k, _ := range h.tagData {
		if k == "" {
			continue
		}
		if k != "" && k[0] == '-' {
			attrs = append(attrs, k)
		}
	}
	return attrs
}

func (h *HTML) GetConstructerParameter() string {
	v := h.tagData[""]
	if v != nil {
		return v.Value
	} else {
		return ""
	}
}

func (h *HTML) GetConstructerCode() string {
	return h.code
}

//设置HTML的属性
func (h *HTML) SetAttr(attrName string, attrValue string) string {
	h.tagData[attrName] = &Ch{Type: 0, Value: attrValue}
	return attrValue
}

//设置HTML的属性名
func (h *HTML) SetAttrName(attrName string, nAttrName string) {
	h.tagData[nAttrName] = h.tagData[attrName]
	h.RemoveAttr(attrName)
}

func (h *HTML) RemoveAttr(attrName string) *HTML {
	delete(h.tagData, attrName)
	return h
}

/**
 * 遍历自己内容
 * @param tagName
 * @return
 */
func (h *HTML) Filter(tagName string) []*HTML {
	filter := make([]*HTML, 0, 100)
	for _, v := range h.list {
		if v.TagName() == tagName {
			filter = append(filter, v)
		}
	}
	return filter
}

func (h *HTML) Child() []*HTML {
	return h.list
}

/**
 * 插入指定HTML
 * @param html
 * @param index
 * @return
 */
func (h *HTML) Insert(html *HTML, index int) *HTML {
	tmp := make([]*HTML, len(h.list)+1)
	copy(tmp, h.list)
	copy(tmp[index+1:], tmp[index:])
	tmp[index] = html
	h.list = tmp
	return h
}

/**
 * 插入指定HTML 节点上 插入多个HTML标签
 * @param html
 * @param index
 * @return
 */
func (h *HTML) InsertList(list []*HTML, index int) *HTML {
	ln := len(list)
	tmp := make([]*HTML, len(h.list)+ln)
	copy(tmp, h.list)
	copy(tmp[index+ln:], tmp[index:])
	for i := 0; i < ln; i++ {
		list[i].parent = h
		tmp[index+i] = list[i]
	}
	h.list = tmp
	return h
}

/**
 * 插入指定HTML 字符串
 * @param html
 * @param index
 * @return
 */
func (h *HTML) InsertFromString(value string, index int) (*HTML, error) {
	html := &HTML{}
	_, err := html.ReadFromString(value)
	if err != nil {
		return h, err
	}
	return h.Insert(html, index), nil
}

//获取指定索引节点的HTML
func (h *HTML) At(index int) *HTML {
	return h.list[index]
}

//将此节点删除
func (h *HTML) Remove() {

	if h.parent == nil {
		list := h.list
		for i := len(list) - 1; i >= 0; i-- {
			list[i].Remove()
		}
	} else {
		for i, v := range h.parent.list {
			if v == h {
				if i > 0 {
					if t := h.parent.list[i-1]; t.IsText() {
						t.value = strings.TrimSpace(t.value)
					}
				}
				if i < len(h.parent.list)-1 {
					if t := h.parent.list[i+1]; t.IsText() {
						t.value = strings.TrimSpace(t.value)
					}
				}
				h.parent.list = deleteHTML(h.parent.list, i)
				break
			}
		}
	}
}

//通过标签名删除
func (h *HTML) RemoveChildByTagName(tagName string) {
	var v *HTML
	for i := len(h.list) - 1; i >= 0; i-- {
		v = h.list[i]
		if v.TagName() == tagName {
			h.list = deleteHTML(h.list, i)
		}
	}
}

/**
 * 替换现有HTML
 * @param value
 * @return
 */
func (h *HTML) ReplaceWith(html *HTML) *HTML {
	t := &HTML{}
	t.ReadFromString(html.ToString())
	l := t.Child()
	if len(l) == 0 {
		return html
	}
	t = l[0]
	if h.parent != nil {
		for i, v := range h.parent.list {
			if v == h {
				h.parent.list[i] = t
				t.parent = h.parent
				break
			}
		}
	}
	return t
}

/**
 * 替换现有HTML，通过String
 * @param value
 * @return
 */
func (h *HTML) ReplaceWithFromList(list []*HTML) []*HTML {
	if h.parent != nil {
		for i, v := range h.parent.list {
			if v == h {
				h.Remove()
				h.parent.InsertList(list, i)
				break
			}
		}
	} else {
		fmt.Println("=nil")
	}
	return list
}

/**
 * 替换现有HTML，通过String
 * @param value
 * @return
 */
func (h *HTML) ReplaceWithFromString(value string) (*HTML, error) {
	html := &HTML{}
	_, err := html.ReadFromString(value)
	if err != nil {
		return h, err
	}
	return h.ReplaceWith(html), nil
}

/**
 * 将子节点全部替换点
 * InnerHTML
 */
func (h *HTML) ReplaceInnerWidthHTML(html *HTML) *HTML {
	h.list = h.list[0:0]
	for i := 0; i < html.Length(); i++ {
		html.At(i).parent = h
		h.list = append(h.list, html.At(i))
	}
	h.tagType = 1
	return h
}

/**
 * 将子节点全部替换点
 * 通过String
 * InnerHTML
 */
func (h *HTML) SetInnerHTML(value string) (*HTML, error) {
	if len(value) == 0 {
		h.list = make([]*HTML, 0)
		return h, nil
	}
	html := &HTML{}
	_, err := html.ReadFromString(value)
	if err != nil {
		return h, err
	}
	return h.ReplaceInnerWidthHTML(html), nil
}

func (h *HTML) GetInnerHTML() string {
	return ListToHTMLString(h.Child())
}

/**
 * 插入特殊字符串
 * @param value
 * @return
 */
func (h *HTML) SetInnerString(value string) *HTML {
	if len(value) == 0 {
		h.list = make([]*HTML, 0)
		return h
	}
	h.list = h.list[0:0]
	html := &HTML{value: value, tagType: -1}
	h.list = append(h.list, html)
	return h
}

/**
 * 复制标签所有属性
 * @param value
 * @return
 */
func (h *HTML) CopyFrom(html *HTML) {
	if h.tag != "" {
		for _, v := range html.Attrs() {
			h.SetAttr(v.Name, v.Value)
		}
	}
}

/**
 * 属性列表
 * @return
 */
func (h *HTML) Attrs() []*Attr {
	arr := make([]*Attr, 0, 20)
	for name, value := range h.tagData {
		if value == nil {
			arr = append(arr, &Attr{Name: name, Value: ""})
		} else {
			if value.Type == 0 {
				arr = append(arr, &Attr{Name: name, Value: value.Value})
			}
		}
	}
	return arr
}

/// 高级属性
func (h *HTML) AdvanceAttrs() []*AdvAttr {
	arr := make([]*AdvAttr, 0, 20)
	for name, value := range h.tagData {
		if value != nil {
			arr = append(arr, &AdvAttr{name, value})
		}
	}
	return arr
}

//在指定节点追加HTML
func (h *HTML) Append(list *HTML) {
	if list.tag == "" {
		for _, v := range list.list {
			v.parent = h
			h.list = append(h.list, v)
		}
	} else {
		list.parent = h
		h.list = append(h.list, list)
	}
}

//在指定节点名的文本
func (h *HTML) AppendNode(tagName string, value string) {
	tag := &HTML{tag: tagName, tagType: 1}
	tag.list = append(tag.list, &HTML{value: value, tagType: -1, parent: tag})
	tag.parent = h
	h.list = append(h.list, tag)
}

//按标签获取元素列表
func (h *HTML) GetElementsByTagName(tagName string) []*HTML {
	tmp := make([]*HTML, 0, 100)
	return h.getElementsByTagName(tagName, &tmp)
}

func (h *HTML) getElementsByTagName(tagName string, buffer *[]*HTML) []*HTML {
	for _, v := range h.list {
		if v.tag == tagName {
			*buffer = append(*buffer, v)
		}
		v.getElementsByTagName(tagName, buffer)
	}
	return *buffer
}

/**
 * 获取非Text的Child Element
 */
func (h *HTML) GetUnTextChild() []*HTML {
	arr := make([]*HTML, 0, 100)
	for _, p := range h.list {
		if p.tagType != -1 {
			arr = append(arr, p)
		}
	}
	return arr
}

//按元素id获取HTML
func (h *HTML) GetElementById(id string) *HTML {
	var p *HTML = nil
	if id == "" {
		return nil
	}
	for _, v := range h.list {
		if v.GetAttr("id") == id {
			return v
		}
		p = v.GetElementById(id)
		if p != nil {
			return p
		}
	}
	return nil
}

//获取当点子节点的长度
func (h *HTML) Length() int {
	return len(h.list)
}

/**
 * 如果内容为空
 */
func (h *HTML) IsEmpty() bool {
	if len(h.list) == 0 {
		return true
	}
	return false
}

/**
 * 获取这个节点下的所有文本
 */
func (h *HTML) Text() string {
	if h.tag == "!" {
		return Substring(h.value, 2, StringLen(h.value)-2)
	} else if isKeyWord(h.tag) {
		return ListToHTMLString(h.Child())
	} else {
		sb := make([]rune, 0, 1000)
		code := []rune(h.ToString())
		p := 0
		var ch rune
		for p < len(code) {
			ch = code[p]
			p++
			if ch == '<' {
				for p < len(code) {
					ch = code[p]
					p++
					if ch == '>' {
						break
					}
				}
			} else {
				sb = append(sb, ch)
			}
		}
		return string(sb)
	}
}

func (h *HTML) IsText() bool {
	if h.tagType == -1 {
		return true
	}
	return false
}

func (h *HTML) ToXHTML() string {
	if h.tag == "!" {
		return "<!" + h.value + ">"
	}
	if h.tagType == -1 {
		return h.value
	}
	sb := bytes.NewBufferString("")
	if h.parent != nil {
		sb.WriteString("<")
		sb.WriteString(h.tag)
		var keys []string
		for k := range h.tagData {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		for _, v := range keys {
			if h.tagData[v] == nil {
				sb.WriteString(" " + v)
				continue
			}
			switch h.tagData[v].Type {
			case 0:
				sb.WriteString(" " + v + "=" + "\"" + h.tagData[v].Value + "\"")
			}
		}
		if h.tagType == 0 {
			sb.WriteString("/>")
		} else {
			sb.WriteString(">")
		}

	}
	list := h.list
	for _, v := range list {
		sb.Write(v.ToStringBytes(true))
	}
	if h.parent != nil && h.tagType == 1 {
		sb.WriteString("</" + h.tag + ">")
	}
	return sb.String()
}

/**
 * 将HTML转换为字符串
 */
func (h *HTML) ToString() string {
	if h.tag == "!" {
		return "<!" + h.value + ">"
	}
	if h.tagType == -1 {
		return h.value
	}
	sb := bytes.NewBufferString("")
	if h.parent != nil {
		sb.WriteString("<")
		sb.WriteString(h.tag)
		var keys []string
		for k := range h.tagData {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		for _, v := range keys {
			if h.tagData[v] == nil {
				sb.WriteString(" " + v)
				continue
			}
			switch h.tagData[v].Type {
			case -1:
				sb.WriteString("(" + h.tagData[v].Value + ")")
			case 0:
				sb.WriteString(" " + v + "=" + "\"" + h.tagData[v].Value + "\"")
			case 1:
				sb.WriteString(" " + v + "=" + "(" + h.tagData[v].Value + ")")
			case 2:
				sb.WriteString(" " + v + "=" + "[" + h.tagData[v].Value + "]")
			case 3:
				sb.WriteString(" " + v + "=" + "{" + h.tagData[v].Value + "}")
			case 4:
				sb.WriteString(" " + v + "=" + h.tagData[v].Value)
			}
		}
		if h.tagType == 0 {
			sb.WriteString("/>")
		} else {
			sb.WriteString(">")
		}

	}
	list := h.list
	for _, v := range list {
		sb.Write(v.ToStringBytes(false))
	}
	if h.parent != nil && h.tagType == 1 {
		sb.WriteString("</" + h.tag + ">")
	}
	return sb.String()
}

func (h *HTML) ToTextStringBytes() []byte {
	if h.tag == "!" {
		return []byte("&lt;!" + h.value + "&gt;<br/>")
	}
	if h.tagType == -1 {
		return []byte(h.value)
	}
	sb := bytes.NewBufferString("")
	if h.parent != nil {
		sb.WriteString("&lt;")
		if strings.ToLower(h.tag) == "style" || strings.ToLower(h.tag) == "css" {
			sb.WriteString("<span style='font-weight:bold;color:#7f0096'>" + h.tag + "</span>")
		} else {
			sb.WriteString("<span style='font-weight:bold;color:#009688'>" + h.tag + "</span>")
		}

		for _, v := range h.tagList {
			sb.WriteString(" <span style='color: #FF5722;font-weight: bold;'>" + v + "</span>=" + "\"<span style='color:#888888'>" + h.tagData[v].Value + "</span>\"")
		}

		if h.tagType == 0 {
			sb.WriteString("/&gt;<br/>")
		} else {
			sb.WriteString("&gt;<br/>")
		}

	}
	list := h.list
	sb.WriteString("<div style='padding:4px 20px'>")
	if strings.ToLower(h.tag) == "style" || strings.ToLower(h.tag) == "css" {
		sb.WriteString("<pre>")
		for _, v := range list {
			sb.Write(v.ToTextStringBytes())
		}
		sb.WriteString("</pre>")
	} else {
		for _, v := range list {
			sb.Write(v.ToTextStringBytes())
		}
	}
	sb.WriteString("</div>")
	if h.parent != nil && h.tagType == 1 {
		if strings.ToLower(h.tag) == "style" || strings.ToLower(h.tag) == "css" {
			sb.WriteString("&lt;/<span style='font-weight:bold;color:#7f0096'>" + h.tag + "</span>&gt;<br/>")
		} else {
			sb.WriteString("&lt;/<span style='font-weight:bold;color:#009688'>" + h.tag + "</span>&gt;<br/>")
		}
	}
	return sb.Bytes()
}

/**
 * 将HTML转换为字符串
 */
func (h *HTML) ToTextString() string {
	if h.tag == "!" {
		return "&lt;!" + h.value + "&gt;<br/>"
	}
	if h.tagType == -1 {
		return h.value
	}
	sb := bytes.NewBufferString("")
	if h.parent != nil {
		sb.WriteString("&lt;")
		if strings.ToLower(h.tag) == "style" || strings.ToLower(h.tag) == "css" {
			sb.WriteString("<span style='font-weight:bold;color:#7f0096'>" + h.tag + "</span>")
		} else {
			sb.WriteString("<span style='font-weight:bold;color:#009688'>" + h.tag + "</span>")
		}

		for i, v := range h.tagData {
			sb.WriteString(" <span style='color: #FF5722;font-weight: bold;'>" + i + "</span>=" + "\"<span style='color:#888888'>" + v.Value + "</span>\"")
		}
		if h.tagType == 0 {
			sb.WriteString("/&gt;<br/>")
		} else {
			sb.WriteString("&gt;<br/>")
		}

	}
	list := h.list
	sb.WriteString("<div style='padding:4px 20px'>")
	if strings.ToLower(h.tag) == "style" || strings.ToLower(h.tag) == "css" {
		sb.WriteString("<pre>")
		for _, v := range list {
			sb.Write(v.ToTextStringBytes())
		}
		sb.WriteString("</pre>")
	} else {
		for _, v := range list {
			sb.Write(v.ToTextStringBytes())
		}
	}
	sb.WriteString("</div>")
	if h.parent != nil && h.tagType == 1 {
		if strings.ToLower(h.tag) == "style" || strings.ToLower(h.tag) == "css" {
			sb.WriteString("&lt;/<span style='font-weight:bold;color:#7f0096'>" + h.tag + "</span>&gt;<br/>")
		} else {
			sb.WriteString("&lt;/<span style='font-weight:bold;color:#009688'>" + h.tag + "</span>&gt;<br/>")
		}
	}
	return sb.String()
}
func (h *HTML) ToStringBytes(xf bool) []byte {
	if h.tag == "!" {
		return []byte("<!" + h.value + ">")
	}
	if h.tagType == -1 {
		return []byte(h.value)
	}
	sb := bytes.NewBufferString("")
	if h.parent != nil {
		sb.WriteString("<")
		sb.WriteString(h.tag)
		var keys []string
		for k := range h.tagData {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, v := range keys {
			if h.tagData[v] == nil {
				sb.WriteString(" " + v)
				continue
			}
			switch h.tagData[v].Type {
			case -1:
				if !xf {
					sb.WriteString("(" + h.tagData[v].Value + ")")
				}
			case 0:
				sb.WriteString(" " + v + "=" + "\"" + h.tagData[v].Value + "\"")
			case 1:
				if !xf {
					sb.WriteString(" " + v + "=" + "(" + h.tagData[v].Value + ")")
				}
			case 2:
				if !xf {
					sb.WriteString(" " + v + "=" + "[" + h.tagData[v].Value + "]")
				}
			case 3:
				if !xf {
					sb.WriteString(" " + v + "=" + "{" + h.tagData[v].Value + "}")
				}
			case 4: //正则表达
				if !xf {
					sb.WriteString(" " + v + "=" + h.tagData[v].Value)
				}
			}
		}
		if h.tagType == 0 {
			sb.WriteString(" />")
		} else {
			sb.WriteString(" >")
		}

	}
	list := h.list
	for _, v := range list {
		sb.Write(v.ToStringBytes(xf))
	}
	if h.parent != nil && h.tagType == 1 {
		sb.WriteString("</" + h.tag + ">")
	}
	return sb.Bytes()
}

func deleteHTML(a []*HTML, index int) []*HTML {
	copy(a[index:], a[(index+1):])
	a = a[:(len(a) - 1)]
	return a
}

/**
 * 列表转换为HTMLString
 */
func ListToHTMLString(l []*HTML) string {
	str := bytes.NewBufferString("")
	for _, v := range l {
		str.WriteString(v.ToString())
	}
	return str.String()
}

/**
 * 列表转换为HTMLString
 */
func ListToHTMLStringBytes(l []*HTML) []byte {
	str := bytes.NewBufferString("")
	for _, v := range l {
		str.WriteString(v.ToString())
	}
	return str.Bytes()
}
