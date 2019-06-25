// mscript.go
package util

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

//---------------------------UName-----------------------------

var __COUNT__ int = -1

func getName() string {
	__COUNT__++
	return "a" + strconv.Itoa(__COUNT__)
}

//----------------------------Var-------------------------------

type Var struct {
	Area, Name, Value, Note, VarType string
	IsStatic                         bool
}

//-------------------------Function-----------------------------

/**
 * @param area		函数可见范围
 * @param isStatic	是否静态
 * @param name		函数名
 * @param param		参数
 * @param note		注释
 */
type Function struct {
	Name, Note, Area string
	Param            []*Var
	IsPublic         bool
	IsStatic         bool
	IsSet            bool
	IsGet            bool
	FunctionType     string
}

//----------------------------Tag-------------------------------
type Tag struct {
	Level        int    //标签域名级别
	Domain       string //所属域级别
	Value        string
	TagType      int //-4:大注释,-3：中注释,-2:小注释,-1:隐含字符,0：字符，1：字符串,2：运算符,3:域操作符,4:语句结束符,5:换行符,6:数字,7:正则表达式,8:元素自动转换符,9点,10类型声明符,11:三目运算符,12:XML对象
	Cls          string
	PType        int  //参数域
	IsClass      bool //是否为类
	IsInnerClass bool //是否为内部类
	IsKeyWord    bool //是否为关键字
	IsStatic     bool //是否为全局变量
	IsPublic     bool //是否为公开
	IsFunction   bool //是否为声明函数
	IsVar        bool //是否为声明变量
	IsAttr       bool //是否为前引用的属性
	IsParameter  bool //是否为参数
	IsType       bool //是声明变量类型
	IsParamValue bool //是否为函数参数默认值
	IsObjectAttr bool //是否为JSON OBject
	IsSet        bool //是否为Setter
	IsGet        bool //是否为Getter
	IsAnonymous  bool //是否为匿名函数
	Note         *Tag //是否有注释
}

//-------------------------Class----------------------
type Class struct {
	Name    string
	Note    string
	IsInner bool
}

func (t *Tag) SetAttr(attr int, boolean bool) *Tag {
	switch attr {
	case -1:
		t.IsVar = boolean
		break
	case 0:
		t.IsFunction = boolean
		break
	case 1:
		t.IsKeyWord = boolean
		break
	case 2:
		t.IsStatic = boolean
		break
	case 3:
		t.IsPublic = boolean
		break
	case 4:
		t.IsParameter = boolean
		break
	case 5:
		t.IsParamValue = boolean
		break

	}
	return t
}

func (t *Tag) SetDomain(domain string) *Tag {
	t.Domain = domain
	return t
}

//----------------------------TagSet------------------------------
type TagSet struct {
	List map[string]*Tag
}

func (ts *TagSet) Get(key string) *Tag {
	return ts.List[key]
}

//----------------------------MScript------------------------------
/**
 *
 * @author sun
 * JScript 分析器
 * 兼容ECMAScript 5
 * 兼容ECMAScript 6 部分属性
 */
type MScript struct {
	domainList map[string]*TagSet
	kMap       map[string]bool
	lst        []*Tag
	position   int
	code       []rune
	tag        *Tag
	defNode    *Tag //被定义文本注释
	fc         int  //匿名函数递增变量
}

/**
 * 从字符串中读取Javascript
 * @param value
 * @throws Exception
 */
func (m *MScript) ReadFromString(js string) {
	m.tag = &Tag{Value: "", TagType: -99}
	m.domainList = make(map[string]*TagSet, 10)

	//00装入关键字
	keyWord := [...]string{"public", "private", "super", "var", "let", "function", "func", "if", "else", "switch", "case", "while", "for", "in", "do", "static", "import", "new", "include", "return", "class", "extends", "implements", "interface", "this", "@this", "@res", "set", "get", "try", "catch", "finally"}
	m.kMap = make(map[string]bool, 10)
	for _, v := range keyWord {
		m.kMap[v] = true
	}
	//01单字解析
	m.position = 0
	m.code = []rune(js)
	codeLength := len(m.code)
	var ch rune
	var p *Tag = nil
	var tp *Tag = nil
	tag := make([]rune, 0, 1000)
	tType := 0
	for m.position < codeLength {
		ch = m.code[m.position]
		m.position++
		if p != nil {
			if m.kMap[p.Value] {
				p.IsKeyWord = true
			}
		}

		//判断是否是注释
		if ch == '/' {
			nt, err := m.isCareNote()
			if err == nil {
				if nt != 2 {
					tp = &Tag{Value: m.readCareNote(nt), TagType: nt}
					m.lst = append(m.lst, tp)
					continue
				}
			} else {
				fmt.Println("MScript: read note has errors.")
				continue
			}

		}

		if ch == '/' && m.isObj() {
			if len(tag) > 0 {
				if m.isNum(tag) {
					tType = 6
				} else {
					tType = 0
				}
				p = &Tag{Value: string(tag), TagType: tType}

				m.lst = append(m.lst, p)
				tag = tag[0:0]
			}
			m.position--
			tp = &Tag{Value: m.ReadString(), TagType: 7}
			m.lst = append(m.lst, tp)
			continue
		}

		if ch == '<' && m.isObj() { //判断是否为XML结构
			if len(tag) > 0 {
				if m.isNum(tag) {
					tType = 6
				} else {
					tType = 0
				}
				p = &Tag{Value: string(tag), TagType: tType}

				m.lst = append(m.lst, p)
				tag = tag[0:0]
			}
			m.position--
			var tXML *HTML = nil
			tXML, m.position = (&HTML{}).ReadOneBlock(m.code, m.position)
			tp = &Tag{Value: tXML.At(0).ToString(), TagType: 12} //XML对象
			m.lst = append(m.lst, tp)
			continue
		}

		if ch == '"' || ch == '\'' {
			if len(tag) > 0 {
				if m.isNum(tag) {
					tType = 6
				} else {
					tType = 0
				}
				p = &Tag{Value: string(tag), TagType: tType}

				m.lst = append(m.lst, p)
				tag = tag[0:0]
			}
			m.position--
			m.lst = append(m.lst, &Tag{Value: m.ReadString(), TagType: 1})
			continue
		}

		if ch == '!' || ch == '.' || ch == '\r' || ch == '\n' || ch == ' ' || ch == '\t' || ch == '{' || ch == '}' || ch == '[' || ch == ']' || ch == '(' || ch == ')' || ch == ';' || ch == ':' || ch == ',' || ch == '?' || ch == '>' || ch == '=' || ch == '<' || ch == '&' || ch == '|' || ch == '%' || ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '#' {
			if len(tag) > 0 {
				if m.isNum(tag) {
					tType = 6
				} else {
					tType = 0
				}
				p = &Tag{Value: string(tag), TagType: tType}
				m.lst = append(m.lst, p)
				tag = tag[0:0]
			}
			if ch == '!' || ch == '.' || ch == '\r' || ch == '\n' || ch == '{' || ch == '}' || ch == '(' || ch == ')' || ch == '[' || ch == ']' || ch == ';' || ch == ':' || ch == ',' || ch == '?' || ch == '>' || ch == '=' || ch == '<' || ch == '&' || ch == '|' || ch == '%' || ch == ' ' || ch == '\t' || ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '#' {
				if ch == ' ' || ch == '\t' {
					m.lst = append(m.lst, &Tag{Value: string(ch), TagType: -1})
				} else if ch == ';' {
					m.lst = append(m.lst, &Tag{Value: string(ch), TagType: 4})
				} else if ch == '{' || ch == '}' || ch == '(' || ch == ')' || ch == '[' || ch == ']' {
					m.lst = append(m.lst, &Tag{Value: string(ch), TagType: 3})
				} else if ch == '\r' || ch == '\n' {
					m.lst = append(m.lst, &Tag{Value: string(ch), TagType: 5})
				} else if ch == '#' {
					m.lst = append(m.lst, &Tag{Value: string(ch), TagType: 8})
				} else if ch == '.' {
					m.lst = append(m.lst, &Tag{Value: string(ch), TagType: 9})
				} else {
					m.lst = append(m.lst, &Tag{Value: string(ch), TagType: 2})
				}

			}

			continue
		}
		tag = append(tag, ch)
	}

	if len(tag) > 0 {
		if m.isNum(tag) {
			tType = 6
		} else {
			tType = 0
		}
		p = &Tag{Value: string(tag), TagType: tType}
		m.lst = append(m.lst, p)
		tag = tag[0:0]
	}

	//02归拢数字
	tlst := make([]*Tag, 0, 1000)
	tag = tag[0:0]
	for _, p := range m.lst {
		if p.TagType == 9 || p.TagType == 6 {
			tag = appendRunes(tag, []rune(p.Value))
		} else {
			if len(tag) != 0 {
				if len(tag) == 1 {
					tlst = append(tlst, &Tag{Value: string(tag), TagType: 9})
				} else {
					tlst = append(tlst, &Tag{Value: string(tag), TagType: 6})
				}
				tag = tag[0:0]
			}
			tlst = append(tlst, p)
		}
	}
	m.lst = tlst
	tlst = tlst[0:0]

	//03归拢操作符
	tag = tag[0:0]
	for _, p := range m.lst {
		if p.TagType == 2 {
			tag = appendRunes(tag, []rune(p.Value))
		} else {
			if len(tag) != 0 {
				tlst = append(tlst, &Tag{Value: string(tag), TagType: 2})
				tag = tag[0:0]
			}
			tlst = append(tlst, p)
		}
	}

	if len(tag) != 0 {
		tlst = append(tlst, &Tag{Value: string(tag), TagType: 2})
	}
	tag = tag[0:0]
	m.lst = tlst
	tlst = tlst[0:0]
	for _, p := range m.lst {
		if p.TagType == 5 {
			tag = appendRunes(tag, []rune(p.Value))
		} else {
			if len(tag) != 0 {
				tlst = append(tlst, &Tag{Value: string(tag), TagType: 5})
				tag = tag[0:0]
			}
			tlst = append(tlst, p)
		}
	}
	if len(tag) != 0 {
		tlst = append(tlst, &Tag{Value: string(tag), TagType: 5})
	}
	m.lst = tlst
	tlst = tlst[0:0]

	//04重新归拢语句
	var note *Tag = nil

	i := 0
	for i < len(m.lst) {
		i = m.readArea(i, "class", 0) + 1
	}

	//分配变量作用域
	var tagSet *TagSet = nil
	for _, p := range m.lst {
		if p.TagType == 0 && !p.IsKeyWord && !p.IsClass && !p.IsFunction && !p.IsVar && !p.IsAttr && !p.IsParameter && !p.IsObjectAttr {
			for p.Domain != "" {
				tagSet = m.domainList[p.Domain]
				if tagSet == nil {
					p.Domain = m.getParentDomain(p.Domain)
					continue
				}
				note = tagSet.Get(p.Value)

				if note != nil {
					p.IsPublic = note.IsPublic
					p.IsStatic = note.IsStatic
					p.IsGet = note.IsGet
					p.IsSet = note.IsSet
					break
				} else {
					p.Domain = m.getParentDomain(p.Domain)
				}
			}
		}

	}

} //ReadFromString

/**
 * 添加数据
 * @param tag
 * @return
 */
func (m *MScript) Push(tag *Tag) *MScript {
	m.lst = append(m.lst, tag)
	return m
}

/**
 *
 * @param domain
 * @return
 */
func (m *MScript) getParentDomain(domain string) string {
	index := strings.LastIndex(domain, ".")
	if index == -1 {
		return ""
	}
	return string([]rune(domain)[0:index])
}

/**
 * 获取所有标签数据
 * @return
 */
func (m *MScript) GetData() []*Tag {
	return m.lst
}

/**
 * 判断是否使注释类型
 */
func (m *MScript) isCareNote() (int, error) {
	var ch rune
	offset := m.position
	if offset < len(m.code) {
		ch = m.code[offset]
		offset++
		if ch == '/' {
			return -2, nil
		}
		if ch == '*' {
			if offset < len(m.code) {
				ch = m.code[offset]
				if ch == '*' {
					return -4, nil
				}
				return -3, nil
			} else {
				return -3, fmt.Errorf("%s", "The Note isn't over.")
			}

		}
	}
	return 2, nil
}

/**
 * 读取注释信息
 */
func (m *MScript) readCareNote(noteType int) string {
	var ch rune
	result := ""
	sb := make([]rune, 0)
	offset := 0
	if noteType == -2 {
		offset = m.position + 1
		for offset < len(m.code) {
			ch = m.code[offset]
			offset++
			if ch == '\r' || ch == '\n' {
				result = string(sb)
				//if offset < len(m.code) && m.code[offset] == '\n' {
				//	offset++
				//}
				break
			}
			sb = append(sb, ch)
		}
	}
	if noteType == -3 || noteType == -4 {
		end := [2]rune{'*', '/'}
		pos := 0
		isNewLine := true
		offset = m.position - noteType - 2
		for offset < len(m.code) {
			ch = m.code[offset]
			offset++
			if ch == '\n' {
				isNewLine = true
				continue
			}
			if ch == end[pos] {
				pos++
				if pos == 2 {
					result = string(sb)
					break
				}
				if isNewLine {
					isNewLine = false
					continue
				}
			} else {
				pos = 0
			}
			if isNewLine && (ch == ' ' || ch == '\t') {
				continue
			} else {
				isNewLine = false
			}
			sb = append(sb, ch)
		}
	}
	m.position = offset
	return result
}

/**
 * 是否为正则表达式的开端
 * @return
 */
func (m *MScript) isObj() bool {
	offset := m.position - 2
	var ch rune
	var p *Tag = nil
	for offset >= 0 {
		ch = m.code[offset]
		offset--
		if ch == '\r' || ch == '\n' || ch == '\t' || ch == ' ' {
			continue
		}
		if ch == ',' || ch == '(' || ch == '=' {
			return true
		} else {
			offset = len(m.lst) - 1
			for offset >= 0 {
				p = m.lst[offset]
				offset--
				if p.TagType < 0 {
					continue
				}
				if (p.IsKeyWord && "return" == p.Value) || "?" == p.Value || ":" == p.Value {
					return true
				} else {
					return false
				}
			}
			return false
		}
	}
	return false
}

/**
 * 返回所有行
 */
func (m *MScript) ToString() string {
	return m.ToStringFrom(m.lst)
}

func (m *MScript) ToStringFrom(lst []*Tag) string {
	sb := bytes.NewBufferString("")
	for _, t := range lst {
		if t.TagType < -1 {
			continue
		}
		sb.WriteString(t.Value)
	}

	return sb.String()
}

/**
 * 读取字符串或正则表达式
 * @param code
 * @return
 */
func (m *MScript) ReadString() string {
	sb := make([]rune, 0)
	var t = m.code[m.position]
	m.position++
	var ch rune
	r := false
	sb = append(sb, t)
	for m.position < len(m.code) {
		ch = m.code[m.position]
		m.position++
		sb = append(sb, ch)
		if ch == t && !r {
			break
		}
		if ch == '\\' {
			if r {
				r = false
			} else {
				r = true
			}
		} else {
			r = false
		}
	}

	return strings.Replace(string(sb), "\r\n", string(t)+" + "+string(t), -1)
}

/**
 * 读取数字
 * @return
 */
func (m *MScript) isNum(sb []rune) bool {
	for _, ch := range sb {
		if !(ch >= '0' && ch <= '9') {
			return false
		}
	}
	return true
}

/**
 *
 * @param lst
 * @param index
 * @param domain
 * @param type		类型域 ：10代表参数类型域
 * @return
 * @throws Exception
 */
func (m *MScript) readArea(i int, domain string, paramType int) int {
	m.clearTag()
	var p *Tag = nil
	hLevel := 0
	zLevel := 0
	xLevel := 0
	for i < len(m.lst) {
		p = m.lst[i]
		i++
		p.PType = paramType
		p.Domain = domain
		if p.TagType < 0 || p.TagType == 5 {
			if p.TagType < -1 { //&& "class" == p.Domain
				m.defNode = p
			}
			continue
		}
		if p.TagType == 3 {
			if "(" == p.Value {
				xLevel++
			} else if ")" == p.Value {
				xLevel--
			} else if "[" == p.Value {
				zLevel++
			} else if "]" == p.Value {
				zLevel--
			} else if "{" == p.Value {
				hLevel++
			} else if "}" == p.Value {
				hLevel--
			}
		}

		if xLevel < 0 || zLevel < 0 || hLevel < 0 {
			i--
			break
		}
		if p.TagType == 2 && "," == p.Value && (xLevel == 0 && zLevel == 0 && hLevel == 0) {
			i--
			break
		}

		if p.TagType == 4 {
			i--
			break
		}

		if p.TagType == 9 && "." == p.Value {
			i = m.attrMethod(i, domain, paramType)
			continue
		}

		if p.TagType == 3 && "{" == p.Value {
			i = m.jsonMethod(i, domain, paramType)
			continue
		}

		if p.IsKeyWord {
			if "private" == p.Value {
				m.tag.IsPublic = false
			} else if "class" == p.Value {
				m.tag.IsClass = true
				i = m.classMethod(i, domain, paramType)
			} else if "static" == p.Value {
				m.tag.IsStatic = true
			} else if "if" == p.Value {
				i = m.logicMethod(i, domain, paramType)
			} else if "else" == p.Value {
				i = m.elseLogicMethod(i, domain, paramType)
			} else if "for" == p.Value {
				i = m.logicMethod(i, domain, paramType)
			} else if "while" == p.Value {
				i = m.logicMethod(i, domain, paramType)
			} else if "do" == p.Value {
				i = m.elseLogicMethod(i, domain, paramType)
			} else if "switch" == p.Value {
				i = m.logicMethod(i, domain, paramType)
			} else if "try" == p.Value {
				i = m.elseLogicMethod(i, domain, paramType)
			} else if "catch" == p.Value {
				i = m.logicMethod(i, domain, paramType)
			} else if "finally" == p.Value {
				i = m.logicMethod(i, domain, paramType)
			} else if "var" == p.Value || "let" == p.Value {
				m.tag.IsVar = true
				i = m.varMethod(i, domain, m.tag, false, paramType)
			} else if "function" == p.Value {
				m.tag.IsFunction = true
				i = m.funcMethod(i, domain, m.tag, paramType)
			} else if "func" == p.Value {
				p.Value = "function"
				m.tag.IsFunction = true
				i = m.funcMethod(i, domain, m.tag, paramType)
			} else if "set" == p.Value {
				m.tag.IsFunction = true
				m.tag.IsSet = true
				i = m.funcMethod(i, domain, m.tag, paramType)
			} else if "get" == p.Value {
				m.tag.IsFunction = true
				m.tag.IsGet = true
				i = m.funcMethod(i, domain, m.tag, paramType)
			}
			continue
		}
	}
	return i
}

func (m *MScript) clearTag() {
	m.tag.IsAttr = false
	m.tag.IsFunction = false
	m.tag.IsVar = false
	m.tag.IsStatic = false
	m.tag.IsSet = false
	m.tag.IsGet = false
	m.tag.IsPublic = true
	m.tag.IsClass = false
}

/**
 * 是否为属性
 * @param lst2
 * @param i
 * @param domain
 * @param type
 * @return
 */
func (m *MScript) attrMethod(i int, domain string, paramType int) int {
	var p *Tag = nil
	for i < len(m.lst) {
		p = m.lst[i]
		i++
		if p.TagType < 0 || p.TagType == 5 {
			continue
		}
		if p.TagType == 0 && !p.IsKeyWord {
			p.IsAttr = true
			break
		}
		break
	}
	return i
}

func (m *MScript) elseLogicMethod(i int, domain string, paramType int) int {
	var p *Tag = nil
sg:
	for i < len(m.lst) {
		p = m.lst[i]
		i++
		if p.TagType < 0 || p.TagType == 5 {
			continue
		}

		if p.IsKeyWord && "if" == p.Value {
			return m.logicMethod(i, domain, paramType)
		}

		if p.TagType == 3 && "{" == p.Value {
			for i < len(m.lst) {
				i = m.readArea(i, domain, paramType) + 1
				if "}" == m.lst[i-1].Value {
					break sg
				}
			}
		}

	}
	return i
}

/**
 * 读取Class 内容
 * @param lst
 * @param i
 * @param domain
 * @param type
 * @return
 * @throws Exception
 */
func (m *MScript) classMethod(i int, domain string, paramType int) int {
	var p *Tag = nil
	isClass := false
sg:
	for i < len(m.lst) {
		p = m.lst[i]
		i++
		if p.TagType < 0 || p.TagType == 5 {
			continue
		}

		if p.TagType == 0 {
			m.copyTag(p, m.tag)
			p.Domain = domain
			p.Note = m.defNode
			if domain == "class" {
				p.IsInnerClass = true
			}
			m.defNode = nil
			domain = p.Value
			isClass = true
		}
		if isClass {
			if p.TagType == 3 && "{" == p.Value {
				for i < len(m.lst) {

					i = m.readArea(i, domain, paramType) + 1
					if "}" == m.lst[i-1].Value {
						break sg
					}
				}
			}
		}
	}
	return i
}

/**
 * 逻辑判断方法
 * @param lst2
 * @param i
 * @param domain
 * @param type
 * @return
 * @throws Exception
 */
func (m *MScript) logicMethod(i int, domain string, paramType int) int {
	var p *Tag = nil
	isArea := false //逻辑判断域名
sg:
	for i < len(m.lst) {
		p = m.lst[i]
		p.Domain = domain
		i++
		if p.TagType < 0 || p.TagType == 5 {
			continue
		}
		if p.TagType == 3 && "(" == p.Value {
			i = m.readArea(i, domain, paramType)
			continue
		}

		if p.TagType == 3 && ")" == p.Value {
			isArea = true
			continue
		}
		if isArea {
			if p.TagType == 3 && "{" == p.Value {
				for i < len(m.lst) {
					i = m.readArea(i, domain, paramType) + 1
					if "}" == m.lst[i-1].Value {
						break sg
					}
				}
			} else {
				if p.TagType >= 0 && p.TagType != 3 {
					break sg
				}
			}
		}
	}
	return i
}

/**
 * 读取JSON格式
 * @param lst
 * @param i
 * @param domain
 * @return
 * @throws Exception
 */
func (m *MScript) jsonMethod(i int, domain string, paramType int) int {

	var p *Tag = nil
	isValue := false
	for i < len(m.lst) {
		p = m.lst[i]
		i++
		p.PType = paramType
		if p.TagType < 0 || p.TagType == 5 {
			continue
		}

		if p.TagType == 3 && "}" == p.Value {
			i--
			break
		}

		if p.TagType == 2 && "," == p.Value {
			continue
		}

		if p.TagType == 0 && !p.IsKeyWord {

			p.IsObjectAttr = true
			continue
		}

		if p.TagType == 2 && ":" == p.Value {
			isValue = true
		}

		if isValue {
			i = m.readArea(i, domain, paramType)
			isValue = false
		}

	}
	return i
}

/**
 * 变量数据
 * @param lst
 * @param index
 * @param domain
 * @param tag
 * @param isFuncParam
 * @param type
 * @return
 * @throws Exception
 */
func (m *MScript) varMethod(i int, domain string, tag *Tag, isFuncParam bool, paramType int) int {
	var name *Tag = nil
	var p *Tag = nil
	isF := false
	for i < len(m.lst) {
		p = m.lst[i]
		i++
		p.PType = paramType
		if p.TagType == 3 && (")" == p.Value || "]" == p.Value || "}" == p.Value) {
			i--
			break
		}
		if p.TagType == 4 {
			break
		}

		if p.TagType == 2 && "," == p.Value {
			name = nil
			p = nil
			isF = false
			continue
		}

		if p.TagType < 0 || p.TagType == 5 {
			continue
		}

		if p.TagType == 0 {
			if isF {
				p.IsType = true
				name.Cls = p.Value
				isF = false
			} else {
				name = p
				p.IsVar = true
				p.IsParameter = isFuncParam
				p.Domain = domain
				m.copyTag(p, tag)
				p.Note = m.defNode
				m.addVar(domain, p)
			}
		}

		if p.TagType == 2 && ":" == p.Value {
			isF = true
			p.TagType = 10
		}

		if p.TagType == 2 && "=" == p.Value {
			if isFuncParam {
				i = m.readArea(i, domain, paramType+1)
			} else {
				i = m.readArea(i, domain, paramType)
			}
		}

		if p.IsKeyWord && "in" == p.Value {
			if isFuncParam {
				i = m.readArea(i, domain, paramType+1)
			} else {
				i = m.readArea(i, domain, paramType)
			}
		}

	}
	return i

}

func (m *MScript) addVar(domain string, p *Tag) {
	tagSet := m.domainList[p.Domain]
	if tagSet == nil {
		tagSet = &TagSet{List: make(map[string]*Tag, 10)}
		m.domainList[p.Domain] = tagSet
	}

	note := tagSet.List[p.Value]
	if note != nil {
		//fmt.Println(p.Value + "变量重复")
	}

	tagSet.List[p.Value] = p
}

func (m *MScript) copyTag(p *Tag, tag *Tag) {
	if tag == nil {
		return
	}
	p.IsPublic = tag.IsPublic
	p.IsSet = tag.IsSet
	p.IsGet = tag.IsGet
	p.IsStatic = tag.IsStatic
	p.IsFunction = tag.IsFunction
	p.IsVar = tag.IsVar
	p.IsClass = tag.IsClass
}

func (m *MScript) getName() string {
	m.fc++
	return "f" + strconv.Itoa(m.fc)
}

/**
 * 函数
 * @param lst
 * @param i
 * @return
 * @throws Exception
 */
func (m *MScript) funcMethod(i int, domain string, tag *Tag, paramType int) int {
	var p *Tag = nil
	var funcName *Tag = nil
	isName := false
	isParam := false
	isF := false
sg:
	for i < len(m.lst) {
		p = m.lst[i]
		i++
		p.PType = paramType
		if p.TagType < 0 && p.TagType == 5 {
			continue
		}
		if p.IsKeyWord {
			p.Domain = domain
			if "set" == p.Value {
				tag.IsSet = true
			} else if "get" == p.Value {
				tag.IsGet = true
			}
		}

		if p.TagType == 0 && !p.IsKeyWord {
			if !isName {
				funcName = p
				p.Domain = domain
				domain = domain + "." + funcName.Value
				isName = true
				m.copyTag(p, tag)
				p.Note = m.defNode
				m.addVar(domain, p)
				continue
			}
		}

		if p.TagType == 3 {
			if "{" == p.Value {
				for i < len(m.lst) {
					i = m.readArea(i, domain, paramType) + 1
					if i > len(m.lst) || "}" == m.lst[i-1].Value {
						break sg
					}
				}

				continue
			}

			if "(" == p.Value {
				if funcName == nil {
					funcName = &Tag{Value: " " + m.getName(), TagType: 0}
					funcName.SetAttr(0, true)
					funcName.PType = paramType
					funcName.Domain = domain
					funcName.IsAnonymous = true
					domain = domain + "." + funcName.Value
					isName = true
					m.copyTag(funcName, tag)
					m.lst = insertElementAt(m.lst, funcName, i-1)
					i++
				}
				isParam = true
			} else if ")" == p.Value {
				isParam = false
			}
		}

		if isParam {
			i = m.varMethod(i, domain, nil, isParam, paramType)
		}

		if !isParam && p.TagType == 2 && ":" == p.Value {
			p.TagType = 10
			isF = true
			continue
		}

		if isF && p.TagType == 0 {
			p.IsType = true
			funcName.Cls = p.Value
			isF = false
		}

	}
	m.defNode = nil
	return i
}

/**
 * 获取JUIS预定义JavaScript数据Tag
 * @return
 */
func (m *MScript) GetJUIScriptData() []*Tag {
	out := make([]*Tag, 0)
	m.initFuncParam(m.lst, &out)
	return out
}

func (m *MScript) ToECSMAScript5() string {
	return m.ToStringFrom(m.GetJUIScriptData())
}

/**
 * 整理函数中的参数默认赋值
 * @param lst
 * @param out
 */
func (m *MScript) initFuncParam(lst []*Tag, out *[]*Tag) {
	mainLst := make([]*Tag, 0)
	defLst := make([]*Tag, 0)
	var p *Tag = nil
	var funcP *Tag = nil
	i := 0
	isArea := false
	for i < len(lst) {
		p = lst[i]
		i++
		if p.IsFunction && p.PType == 0 {
			funcP = p
		}
		if !isArea && p.PType > 0 {

			isArea = true
			uname := "_" + getName()
			t := &Tag{Value: uname, TagType: 0}
			t.SetAttr(5, true)
			mainLst = append(mainLst, t)

			t = &Tag{Value: "var", TagType: 0, IsVar: true, IsKeyWord: true}
			//t.SetAttr(1, true).SetDomain("class")
			defLst = append(defLst, t)
			defLst = append(defLst, &Tag{Value: " ", TagType: -1})

			t = &Tag{Value: uname, TagType: 0}
			t.SetAttr(-1, true)
			if funcP != nil {
				t.SetAttr(2, funcP.IsStatic)
			} else {
				t.SetAttr(2, false)
			}
			defLst = append(defLst, t)
			defLst = append(defLst, &Tag{Value: "=", TagType: 2})
		}

		if isArea {
			if p.PType <= 0 {
				isArea = false
				defLst = append(defLst, &Tag{Value: ";", TagType: 4})
				defLst = append(defLst, &Tag{Value: "\r\n", TagType: 5})
				m.initFuncParam(defLst, &mainLst)
				defLst = defLst[0:0]
				funcP = nil
			} else {
				p.PType--
				defLst = append(defLst, p)
				continue
			}
		}
		mainLst = append(mainLst, p)
	}

	*out = insertAllBefore(*out, mainLst)

}

/**
 * 获取指定域的定义列表
 * @param domain
 * @return
 */
func (m *MScript) GetDefine(domain string) *TagSet {
	return m.domainList[domain]
}

/**
 * 获取函数内部数据
 * @return
 */
func (m *MScript) GetFunctionContent(functionName string) []*Tag {
	ls := make([]*Tag, 0)
	var p *Tag = nil
	i := 0
	for i < len(m.lst) {
		p = m.lst[i]
		i++
		if p.IsFunction && p.Value == functionName {

		}
	}
	return ls
}

/**
 * 获取所有共有属性
 * @return
 */
func (m *MScript) GetVar(isPublic bool, isStatic bool) []*Var {
	arr := make([]*Var, 0)
	var v *Var = nil
	for _, t := range m.lst {
		if t.TagType < 0 || t.TagType == 5 {
			continue
		}
		if "class" == t.Domain && t.IsVar && t.IsPublic == isPublic && t.IsStatic == isStatic {
			v = &Var{Area: t.Domain, IsStatic: t.IsStatic, Name: t.Value, Value: "", VarType: t.Cls, Note: isNilNote(t.Note)}
			arr = append(arr, v)
		}
	}
	return arr
}

/**
 * 获取构造函数
 */
func (m *MScript) GetConstructor() *Function {
	var v *Function = nil
	var t *Tag = nil
	var p *Tag = nil
	level := 0
	isParam := false
	for i := 0; i < len(m.lst); i++ {
		t = m.lst[i]
		if t.TagType < 0 || t.TagType == 5 {
			continue
		}
		if "class" == t.Domain && t.IsFunction && t.Value == "init" {
			p = t
			vls := make([]*Var, 0)
			for ; i < len(m.lst); i++ {
				t = m.lst[i]
				if t.TagType == 3 && "(" == t.Value {
					level++
					isParam = true
					continue
				}
				if t.TagType == 3 && ")" == t.Value {
					level--
					continue
				}

				if isParam && level == 0 {
					isParam = false
					break
				}

				if isParam && t.IsVar {
					vls = append(vls, &Var{Area: t.Domain, IsStatic: t.IsStatic, Name: t.Value, Value: m.readParam(m.lst, i), VarType: t.Cls, Note: isNilNote(t.Note)})
				}
			}

			v = &Function{Area: p.Domain, IsPublic: p.IsPublic, IsStatic: p.IsStatic, IsSet: p.IsSet, IsGet: p.IsGet, Name: p.Value, Param: vls, FunctionType: p.Cls, Note: isNilNote(p.Note)}
			return v
		}
	}
	return nil
}

/**
 * 获取构造函数
 */
func (m *MScript) GetConstructorByClassName(className string) *Function {
	var v *Function = nil
	var t *Tag = nil
	var p *Tag = nil
	level := 0
	isParam := false
	for i := 0; i < len(m.lst); i++ {
		t = m.lst[i]
		if t.TagType < 0 || t.TagType == 5 {
			continue
		}
		if className == t.Domain && t.IsFunction && t.Value == "init" {
			p = t
			vls := make([]*Var, 0)
			for ; i < len(m.lst); i++ {
				t = m.lst[i]
				if t.TagType == 3 && "(" == t.Value {
					level++
					isParam = true
					continue
				}
				if t.TagType == 3 && ")" == t.Value {
					level--
					continue
				}

				if isParam && level == 0 {
					isParam = false
					break
				}

				if isParam && t.IsVar {
					vls = append(vls, &Var{Area: t.Domain, IsStatic: t.IsStatic, Name: t.Value, Value: m.readParam(m.lst, i), VarType: t.Cls, Note: isNilNote(t.Note)})
				}
			}

			v = &Function{Area: p.Domain, IsPublic: p.IsPublic, IsStatic: p.IsStatic, IsSet: p.IsSet, IsGet: p.IsGet, Name: p.Value, Param: vls, FunctionType: p.Cls, Note: isNilNote(p.Note)}
			return v
		}
	}
	return nil
}

/**
 * 获取所有共有属性
 * @return
 */
func (m *MScript) GetFunctionAndStatic(isPublic bool, isStatic bool) []*Function {
	arr := make([]*Function, 0)
	funcL := m.GetFunction(isPublic)
	for _, f := range funcL {
		if f.Name == "init" {
			continue
		}
		if isStatic != f.IsStatic {
			continue
		}
		arr = append(arr, f)
	}
	return arr
}

/**
 * 获取所有共有属性
 * @return
 */
func (m *MScript) GetFunctionAndStaticByClassName(className string, isPublic bool, isStatic bool) []*Function {
	arr := make([]*Function, 0)
	funcL := m.GetFunctionByClassName(className, isPublic)
	for _, f := range funcL {
		if f.Name == "init" {
			continue
		}
		if isStatic != f.IsStatic {
			continue
		}
		arr = append(arr, f)
	}
	return arr
}

func (m *MScript) GetFunction(isPublic bool) []*Function {
	arr := make([]*Function, 0)
	var v *Function = nil
	var t *Tag = nil
	var p *Tag = nil
	level := 0
	isParam := false
	for i := 0; i < len(m.lst); i++ {
		t = m.lst[i]
		if t.TagType < 0 || t.TagType == 5 || t.Value == "init" {
			continue
		}
		if "class" == t.Domain && t.IsFunction && t.IsPublic == isPublic {
			p = t
			vls := make([]*Var, 0)
			for ; i < len(m.lst); i++ {
				t = m.lst[i]
				if t.TagType == 3 && "(" == t.Value {
					level++
					isParam = true
					continue
				}
				if t.TagType == 3 && ")" == t.Value {
					level--
					continue
				}

				if isParam && level == 0 {
					isParam = false
					break
				}

				if isParam && t.IsVar {
					vls = append(vls, &Var{Area: t.Domain, IsStatic: t.IsStatic, Name: t.Value, Value: m.readParam(m.lst, i), VarType: t.Cls, Note: isNilNote(t.Note)})
				}
			}

			v = &Function{Area: p.Domain, IsPublic: p.IsPublic, IsStatic: p.IsStatic, IsSet: p.IsSet, IsGet: p.IsGet, Name: p.Value, Param: vls, FunctionType: p.Cls, Note: isNilNote(p.Note)}
			arr = append(arr, v)
		}
	}
	return arr
}

/**
 * 获取Class
 */
func (m *MScript) GetClass() []*Class {
	arr := make([]*Class, 0)

	var t *Tag = nil

	for i := 0; i < len(m.lst); i++ {
		t = m.lst[i]
		if t.TagType < 0 || t.TagType == 5 || t.Value == "init" {
			continue
		}
		if t.IsClass {
			arr = append(arr, &Class{Name: t.Value, Note: m.nil2Str(t.Note), IsInner: t.IsInnerClass})
		}
	}
	return arr
}

func (m *MScript) nil2Str(tag *Tag) string {
	if tag == nil {
		return ""
	} else {
		return tag.Value
	}
}

/**
 * 获取常量信息
 */
func (m *MScript) GetVarByClassName(className string, isPublic bool, isStatic bool) []*Var {
	arr := make([]*Var, 0)
	var v *Var = nil
	for _, t := range m.lst {
		if t.TagType < 0 || t.TagType == 5 {
			continue
		}
		if className == t.Domain && t.IsVar && t.IsPublic == isPublic && t.IsStatic == isStatic {
			v = &Var{Area: t.Domain, IsStatic: t.IsStatic, Name: t.Value, Value: "", VarType: t.Cls, Note: isNilNote(t.Note)}
			arr = append(arr, v)
		}
	}
	return arr
}

func (m *MScript) GetFunctionByClassName(className string, isPublic bool) []*Function {
	arr := make([]*Function, 0)
	var v *Function = nil
	var t *Tag = nil
	var p *Tag = nil
	level := 0
	isParam := false
	for i := 0; i < len(m.lst); i++ {
		t = m.lst[i]
		if t.TagType < 0 || t.TagType == 5 || t.Value == "init" || t.IsPublic != isPublic {
			continue
		}
		if className == t.Domain && t.IsFunction {
			p = t
			vls := make([]*Var, 0)
			for ; i < len(m.lst); i++ {
				t = m.lst[i]
				if t.TagType == 3 && "(" == t.Value {
					level++
					isParam = true
					continue
				}
				if t.TagType == 3 && ")" == t.Value {
					level--
					continue
				}

				if isParam && level == 0 {
					isParam = false
					break
				}

				if isParam && t.IsVar {
					vls = append(vls, &Var{Area: t.Domain, IsStatic: t.IsStatic, Name: t.Value, Value: m.readParam(m.lst, i), VarType: t.Cls, Note: isNilNote(t.Note)})
				}
			}

			v = &Function{Area: p.Domain, IsPublic: p.IsPublic, IsStatic: p.IsStatic, IsSet: p.IsSet, IsGet: p.IsGet, Name: p.Value, Param: vls, FunctionType: p.Cls, Note: isNilNote(p.Note)}
			arr = append(arr, v)
		}
	}
	return arr
}

func isNilNote(value *Tag) string {
	if value != nil {
		return value.Value
	}
	return ""
}

func (m *MScript) readParam(tag []*Tag, i int) string {
	sb := bytes.NewBufferString("")
	var t *Tag = nil
	isParam := false
	xl := 0
	zl := 0
	dl := 0
	for ; i < len(m.lst); i++ {
		t = m.lst[i]
		if t.TagType == 3 {
			if "(" == t.Value {
				xl++
			} else if ")" == t.Value {
				xl--
			} else if "[" == t.Value {
				zl++
			} else if "]" == t.Value {
				zl--
			} else if "{" == t.Value {
				dl++
			} else if "}" == t.Value {
				dl--
			}
		}
		if xl == -1 {
			break
		}
		if t.TagType == 2 && "=" == t.Value {
			isParam = true
		}

		if xl == 0 && zl == 0 && dl == 0 && "," == t.Value {
			break
		}
		if isParam {
			sb.WriteString(t.Value)
		}
	}

	return sb.String()
}

func appendRunes(old []rune, arr []rune) []rune {
	tmp := make([]rune, len(old)+len(arr))
	copy(tmp[copy(tmp, old):], arr)
	return tmp
}

func appendArray(old []*Tag, arr []*Tag) []*Tag {
	tmp := make([]*Tag, len(old)+len(arr))
	copy(tmp[copy(tmp, old):], arr)
	return tmp
}

//在序列中插入指定元素
func insertElementAt(lst []*Tag, tag *Tag, index int) []*Tag {
	tmp := make([]*Tag, len(lst)+1)
	copy(tmp, lst)
	copy(tmp[index+1:], tmp[index:])
	tmp[index] = tag
	return tmp

}

func insertElementAt0(lst []rune, tag rune, index rune) []rune {
	tmp := make([]rune, len(lst)+1)
	copy(tmp, lst)
	copy(tmp[index+1:], tmp[index:])
	tmp[index] = tag
	return tmp

}

func insertAllBefore(lst []*Tag, tags []*Tag) []*Tag {
	tmp := make([]*Tag, len(lst)+len(tags))
	copy(tmp[copy(tmp, tags):], lst)
	return tmp
}

func insertAllBefore0(lst []rune, tags []rune) []rune {
	tmp := make([]rune, len(lst)+len(tags))

	copy(tmp[copy(tmp, tags):], lst)
	return tmp
}
