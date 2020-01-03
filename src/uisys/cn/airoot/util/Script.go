// Script.go

package util

import (
	"fmt"
	"path/filepath"
	"strings"
	. "uisys"
	. "uisys/str"
	. "uisys/tool"
)

//--------------------------------GSetter---------------------------------------
/**
 * Getter Setter 方法
 * @author sun
 *
 */
type GSetter struct {
	Setter *Tag
	Getter *Tag
}

//--------------------------------Script----------------------------------------
type Script struct {
	jus          *UI
	root         string
	hMap         []*Attr  //导入的类文件
	eMap         []string //集成的类文件
	iMap         []string //接口文件
	gsMap        map[string]*GSetter
	domain       string
	value        *Attr
	extendScript string
	mjs          *MScript
	className    string
	isScript     bool
	fromBuf      string //通过from导入的类
}

func (s *Script) CreateFrom(jus *UI, root string, domain string, value *Attr, extendScript string, className string) *Script {
	s.jus = jus
	s.root = root
	s.domain = domain
	s.value = value
	s.extendScript = extendScript
	s.hMap = make([]*Attr, 0)
	s.gsMap = make(map[string]*GSetter)
	s.className = className
	s.isScript = true
	return s
}

func (s *Script) initScript(js *MScript) string {
	return s.initScriptFrom(js, "__OBJECT__", "__this_", "__pri")

}

/**
 * 初始化Script语句
 * @param script
 * @return
 * @throws Exception
 */
func (s *Script) initScriptFrom(js *MScript, _global string, _this string, _pri string) string {
	out := ""
	tmp := ""
	var hObj *HTMLObject = nil
	var tjs *MScript = nil
	lst := js.GetJUIScriptData()
	tl := make([]*Tag, 0)
	tlt := make([]*Tag, 0)

	p := 0
	var t *Tag = nil
	var f *Tag = nil
	var param *Tag = nil
	level := 0
	var isExtends bool = false
	var isImpl bool = false

	//00.将import 的类记录下来
	for p < len(lst) {
		t = lst[p]
		p++
		// 3.import
		if t.IsKeyWord && "import" == t.Value {
			tmp = ""
			point := -1
			at := 0
			isFrom := false
			for p < len(lst) {
				f = lst[p]
				p++
				if f.TagType == 1 {
					f.Value = Replace(f.Value, "\"", "")
					f.Value = Replace(f.Value, "'", "")
				}
				if f.TagType == -1 {
					continue
				}
				if f.TagType == 9 || f.Value == "\\" || f.Value == "/" {
					point = p
				}
				if f.TagType == 4 || f.TagType == 5 {
					if point == -1 {
						point = at
					}
					break
				}
				if f.IsKeyWord { //说明要用commandJS规范导读
					if f.Value == "from" {
						f.Value = "\002"
						isFrom = true
					} else if f.Value == "@root" {
						f.Value = "/index.res"
					} else if f.Value == "@lib" {
						t.Value = "\"" + IfStr(s.jus.IsSysLib, "index.src/", "") + "/" + s.jus.relativePath + ".lib/\""
						if s.jus.IsSysLib {
							s.jus.GetRoot().PushSysLibDirs(s.jus.className, s.jus.relativePath+".lib")
						}

					}
				}
				tmp += f.Value
				at = p - 1
			}
			if Index(tmp, ".") == 0 { //说明是获取自己本地路径
				tmp = Substring(s.jus.dirPath, StringLen(s.jus.root), -1) + "/" + tmp
				tmp = filepath.Clean(tmp)
				tmp = Replace(tmp, "\\", ".")
				tmp = Replace(tmp, "/", ".")
				if CharAt(tmp, 0) == "." {
					tmp = Substring(tmp, 1, -1)
				}
			}
			if Index(tmp, "/") != -1 || Index(tmp, "\\") != -1 {
				Single(&s.hMap, &Attr{tmp, tmp})
			} else {
				Single(&s.hMap, &Attr{lst[point].Value, tmp})
			}
			if isFrom {
				s.fromBuf = ImportFrom(s.jus.className, tmp)
				//tl = append(tl, &Tag{Value: ImportFrom(s.jus.className, tmp), TagType: 1})
			}
			continue
		}
		tl = append(tl, t)
	}

	lst = lst[0:0]
	lst = appendArray(lst, tl)
	tl = tl[0:0]
	p = 0
	level = 0
	//00.将内部class解析出来
	for p < len(lst) {
		t = lst[p]
		p++
		if (t.IsKeyWord && "class" == t.Value) || t.TagType < -1 {
			continue
		}
		if t.IsClass {
			f = t
			tjs = &MScript{}
			for p < len(lst) {
				t = lst[p]
				p++
				if t.TagType == 3 && "{" == t.Value {
					level++
					break
				}
				if isExtends == false && t.IsKeyWord && t.Value == "extends" {
					isExtends = true
					continue
				}
				if isImpl == false && t.IsKeyWord && t.Value == "implements" {
					isExtends = false
					isImpl = true
					continue
				}
				if isExtends && t.TagType == 0 {
					s.eMap = append(s.eMap, t.Value)
				}
				if isImpl && t.TagType == 0 {
					s.iMap = append(s.iMap, t.Value)
				}
			}
			isExtends = false
			for p < len(lst) {
				t = lst[p]
				p++
				if t.TagType == 3 && "{" == t.Value {
					level++
				} else if t.TagType == 3 && "}" == t.Value {
					level--
				}
				if level == 0 {
					tl = append(tl, &Tag{Value: s.initClass(f.Value, tjs.ToECSMAScript5()), TagType: 1})
					tjs = nil
					break
				}
				tjs.Push(t)
			}
			continue
		}
		tl = append(tl, t)
	}

	lst = lst[0:0]
	lst = appendArray(lst, tl)
	tl = tl[0:0]
	//01.去掉js语言不能分析的部分;02.整理js $符号部分;
	p = 0
	level = 0
	for p < len(lst) {
		t = lst[p]
		p++

		if t.TagType == 10 || t.IsType {
			continue
		}
		if t.IsKeyWord && "super" == t.Value {
			t.Value = "__UP__"
		} else if t.TagType == 1 { //初始化$符号
			t.Value = ScriptInitD(t.Value, s.domain)
		}

		tl = append(tl, t)
	}

	lst = lst[0:0]
	lst = appendArray(lst, tl)
	tl = tl[0:0]

	//02.开始执行分析
	p = 0
	for p < len(lst) {
		t = lst[p]
		p++
		//02.01处理静态数据
		if t.IsKeyWord && "static" == t.Value {
			for p < len(lst) {
				t = lst[p]
				p++
				if t.TagType >= 0 {
					break
				}
			}

			if t.TagType == 3 {
				tmp = ""
				p--
				for p < len(lst) {
					f = lst[p]
					p++
					if f.TagType < -1 {
						continue
					}
					tmp += f.Value
					if t.TagType == 3 && "{" == f.Value {
						level++
					} else if f.TagType == 3 && "}" == f.Value {
						level--
						if level == 0 {
							break
						}
					}
				}
				s.jus.AddStaticCode(s.jus.className, "__STATIC__", " = function()"+tmp+";")
				continue
			}

			tl = append(tl, t)
			continue
		}

		//02.02整理内部作用域
		/**
		 * JavaScript 在编写完毕之后，由于其语言的特殊原因（按照面向对象的编写Function方法和缺少静态类的原因），面向对象函数不能区分内部函数和挂在函数，
		 * 因此为了实现短缺的功能，此处必须由人工实现函数域的自定义判断，并指定变量到合适的域内。
		 */
		if t.TagType == 0 && !t.IsKeyWord && !t.IsFunction && !t.IsVar && !t.IsAttr && !t.IsObjectAttr {
			if "class" == t.Domain {
				newString := ""
				if t.IsPublic {
					if t.IsStatic && !t.IsGet && !t.IsSet {
						newString = "__WINDOW__[__APPDOMAIN__]['" + s.className + "']."
					} else {
						newString = _this + "."
					}
					newString += t.Value

				} else {
					if t.IsStatic && !t.IsGet && !t.IsSet {
						newString = "__WINDOW__[__APPDOMAIN__]['" + s.className + "']."
					} else {
						newString = _pri + "."
					}
					newString += t.Value
				}
				tl = append(tl, &Tag{Value: newString, TagType: 0})
			} else if t.Domain == "" {
				if s.jus != nil {
					hObj = s.jus.GetDefine(t.Value)
				}

				if hObj != nil {
					t.Value = hObj.Name
				}
				tl = append(tl, t)
			} else {
				tl = append(tl, t)
			}
			continue
		}

		//02.03.解析关键字
		/**
		 * 这里的关键字包含，#、import、include、new,this等关键字，实际上大部分还是自定义的关键字，
		 * 这里解析的做法是吧关键字转换为实际的JavaScript函数，例如#id转换为$("#id")
		 */
		// 2.#
		if t.TagType == 8 && "#" == t.Value {
			for p < len(lst) {
				f = lst[p]
				p++
				if f.TagType == 0 {
					param = f
					break
				}
			}
			if s.jus != nil {
				hObj = s.jus.GetDefine(param.Value)
			}

			if hObj != nil {
				param.Value = "__NAME__ + '" + param.Value //hObj.Name
				tl = append(tl, &Tag{Value: "$JGID(" + param.Value + "')", TagType: 0})
			} else {
				tl = append(tl, &Tag{Value: "$JGID('" + param.Value + "')", TagType: 0})
			}

			continue
		}

		if t.TagType == 12 {
			//tj := &JUS{SYSTEM_PATH: s.jus.SYSTEM_PATH, CLASS_PATH: s.jus.CLASS_PATH}
			//tj.CreateFromString(s.root, "", nil, t.Value, "temp", s.jus)
			//tl = append(tl, &Tag{Value: "Module(\"" + Escape(tj.ReadHTML().ToString()) + "\",\f)", TagType: 0})
			tl = append(tl, &Tag{Value: "\"" + Escape(t.Value) + "\"", TagType: 1})
			continue
		}

		// 5.include
		if t.IsKeyWord && "include" == t.Value {
			tmp = ""
			for p < len(lst) {
				f = lst[p]
				p++
				if f.TagType == 1 {
					tmp += f.Value
					break
				}
			}

			tl = append(tl, &Tag{Value: s.includeJs(tmp), TagType: 0})
			tl = append(tl, f)
			continue
		}
		tl = append(tl, t)
	} //02开始解析END.

	p = 0
	for p < len(tl) {
		t = tl[p]
		p++
		if "class" == t.Domain && t.IsKeyWord && ("public" == t.Value || "private" == t.Value || "static" == t.Value || "function" == t.Value || "var" == t.Value || "set" == t.Value || "get" == t.Value) {
			continue
		}
		// 4.new
		if t.IsKeyWord && "new" == t.Value {
			tmp = ""
			for p < len(tl) {
				f = tl[p]
				p++
				if "(" == f.Value {
					break
				}
				tmp += f.Value
			}

			newTmp := s.loadClass(tmp)
			if newTmp != "" {
				tlt = append(tlt, &Tag{Value: newTmp, TagType: 0})
			} else {
				tlt = append(tlt, t)
				tlt = append(tlt, &Tag{Value: " ", TagType: -1})
				tlt = append(tlt, &Tag{Value: Replace(tmp, "?", " "), TagType: 1})
			}
			tlt = append(tlt, f)
			continue
		}
		if t.IsKeyWord && "@this" == t.Value {
			t.Value = _this
		} else if t.IsKeyWord && "@global" == t.Value {
			t.Value = _global
		} else if t.IsKeyWord && "@lib" == t.Value {
			t.Value = "\"" + "./" + s.jus.relativePath + ".lib/\""
		} else if t.IsKeyWord && "this" == t.Value {
			tlt = append(tlt, t)
			if s.getLevel(t) == 1 {
				//t.Value = _pri
				param = t
				for p < len(tl) {
					t = tl[p]
					p++
					tlt = append(tlt, t)
					if t.TagType < 0 || t.TagType == 5 {
						continue
					}
					if t.TagType == 9 {
						for p < len(tl) {
							t = tl[p]
							p++
							tlt = append(tlt, t)
							if t.IsAttr {
								set := js.GetDefine("class")
								a := set.Get(t.Value)
								if a != nil {
									if a.IsPublic {
										param.Value = _this
									} else {
										param.Value = _pri
									}
								}
								if s.mjs != js && a == nil {
									c := s.mjs.GetDefine("class")
									if c != nil {
										a = c.Get(t.Value)
										if a.IsPublic {
											param.Value = "__this_"
										} else {
											param.Value = "__pri"
										}
									}
								}
								break
							}
						}
					}
					break
				}
			}
			continue
		}
		tlt = append(tlt, t)
	}

	tl = tl[0:0]
	tl = appendArray(tl, tlt)
	tlt = tlt[0:0]

	//05.处理静态函数
	p = 0
	var paramVar *Tag = nil
	var paramValue *Tag = nil
	var buffer []*Tag = make([]*Tag, 0)

	//06.将函数转义
	p = 0
	isStatic := false
	for p < len(tl) {
		t = tl[p]
		p++
		if t.Domain == "" && t.TagType == 0 && !t.IsAttr {
			he := GetSingle(s.hMap, t.Value)
			if he != nil {
				tlt = append(tlt, &Tag{Value: "__WINDOW__[__APPDOMAIN__]['" + he.Name + "']", TagType: 0})
				continue
			}
		}
		if !t.IsSet && !t.IsGet && "class" == t.Domain && (t.IsVar || t.IsFunction) {
			if t.IsFunction {
				isStatic = t.IsStatic
				if t.IsStatic {
					tlt = append(tlt, t)
				} else {
					if t.IsAnonymous {
						tlt = append(tlt, &Tag{Value: "function", TagType: 0})
					} else {
						tlt = append(tlt, &Tag{Value: IfStr(s.isScript, IfStr(t.IsPublic, _this+".", _pri+".")+t.Value+" = function", "__MODULE_METHOD__['"+s.domain+"']."+t.Value+" = "+IfStr(t.IsPublic, _this+".", _pri+".")+t.Value+" = function"), TagType: 0})
					}
				}

				//读参
				for p < len(tl) {
					t = tl[p]
					p++
					if t.TagType < 0 || t.TagType == 5 {
						continue
					}
					if t.TagType == 3 || (t.TagType == 2 && "," == t.Value) {
						tlt = append(tlt, t)
						if t.TagType == 3 && "{" == t.Value {
							tlt = append(tlt, &Tag{Value: "\r\n", TagType: 5})
							for len(buffer) > 0 {
								tlt = append(tlt, buffer[0])
								buffer = buffer[1:]
							}
							break
						}
					}
					if t.IsVar {
						paramVar = t
						tlt = append(tlt, t)
					}

					if t.IsParamValue {
						paramValue = t
					}

					if paramVar != nil && paramValue != nil {
						buffer = append(buffer, &Tag{Value: paramVar.Value + "=" + paramVar.Value + " || " + IfStr(isStatic, "__WINDOW__[__APPDOMAIN__]['"+s.jus.className+"']."+paramValue.Value, paramValue.Value) + ";\r\n", TagType: 0})
						paramVar = nil
						paramValue = nil
					}
				}
			} else if t.IsVar {
				isStatic = t.IsStatic
				if t.IsStatic {
					tlt = append(tlt, t)
					continue
				}

				tlt = append(tlt, &Tag{Value: IfStr(t.IsPublic, _this+".", _pri+".") + t.Value + " ", TagType: 0})
				//去除属性
				for p < len(tl) {
					t = tl[p]
					p++
					if t.TagType < 0 || t.TagType == 5 {
						continue
					}
					if (t.TagType == 2 && "=" == t.Value) || t.TagType == 4 || t.TagType == 5 {
						tlt = append(tlt, t)
						break
					}
				}
			}
			continue
		}

		if t.IsGet {
			s.pushGSetter(0, t)
			if !t.IsStatic {
				tlt = append(tlt, &Tag{Value: "function " + t.Value, TagType: 0})
				continue
			}

		}

		if t.IsSet {
			s.pushGSetter(1, t)
			if !t.IsStatic {
				tlt = append(tlt, &Tag{Value: "function " + t.Value, TagType: 0})
				continue
			}
		}

		tlt = append(tlt, t)
	}

	tl = tl[0:0]
	tl = appendArray(tl, tlt)
	tlt = tlt[0:0]

	/**
	 * 处理静态函数
	 */
	p = 0
	for p < len(tl) {
		t = tl[p]
		p++
		if t.IsStatic {
			//
			tmp = ""
			if t.IsFunction {
				tmp += " = function"
				for p < len(tl) {
					f = tl[p]
					p++
					if f.TagType < -1 {
						continue
					}
					tmp += f.Value
					if f.TagType == 3 && "{" == f.Value {
						level++
					} else if f.TagType == 3 && "}" == f.Value {
						level--
						if level == 0 {
							break
						}
					}
				}
				//if t.IsPublic {
				//tlt = append(tlt, &Tag{Value: IfStr(t.IsSet, "var", _this+".") + t.Value + " = __WINDOW__[__APPDOMAIN__]['" + s.className + "']." + t.Value + ";", TagType: 0})
				//}
			} else if t.IsVar {
				level = 0
				for p < len(tl) {
					f = tl[p]
					p++
					if f.TagType == 3 {
						if "(" == f.Value || "{" == f.Value {
							level++
						} else if ")" == f.Value || "}" == f.Value {
							level--
						}
					}

					if (f.TagType == 4 || f.TagType == 5) && level == 0 { //;
						break
					}
					tmp += f.Value
				}
			}
			if s.jus != nil {
				s.jus.AddStaticScript(s.className, t.Value, tmp)
			}
			continue
		}
		tlt = append(tlt, t)
	}
	tl = tl[0:0]
	tl = appendArray(tl, tlt)
	tlt = tlt[0:0]

	out += js.ToStringFrom(tl)
	//处理Getter Setter
	var pgs *GSetter = nil
	tsb := ""
	for name, value := range s.gsMap {
		pgs = value
		tsb += "Object.defineProperty(" + _this + ",'" + name + "',{"
		if pgs.Setter != nil {
			tsb += "set:"
			if pgs.Setter.IsStatic {
				tsb += "__WINDOW__[__APPDOMAIN__]['" + s.className + "']." + pgs.Setter.Value
			} else {
				tsb += pgs.Setter.Value
			}
		}
		if pgs.Getter != nil {
			if pgs.Setter != nil {
				tsb += ","
			}
			tsb += "get:"
			if pgs.Getter.IsStatic {
				tsb += "__WINDOW__[__APPDOMAIN__]['" + s.className + "']." + pgs.Getter.Value
			} else {
				tsb += pgs.Getter.Value
			}
		}
		tsb += ",enumerable:true});\r\n"
		delete(s.gsMap, name)
	}

	out += tsb

	return strings.Replace(out, "@this", "__MODULE_METHOD__."+s.domain, -1)

}

/**
 * 初始化$符号
 * @param value
 * @return
 */
func ScriptInitD(value string, domain string) string {
	sb := ""
	code := []rune(value)
	l := len(code)
	var ch rune
	for i := 0; i < l; i++ {
		ch = code[i]
		if ch == '\\' {
			//再读一个
			i++
			ch = code[i]
			if ch == '$' {
				sb += string(ch)
			} else {
				sb += "\\" + string(ch)
			}
			continue
		}
		if ch == '$' {
			//在读一个
			i++
			ch = code[i]
			if (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' || ch == '$' {
				sb += domain
			} else {
				sb += "$"
			}
			sb += string(ch)
			continue
		}

		sb += string(ch)

	}

	return sb
}

/**
 * 获取变量是第几层
 * @param t
 * @return
 */
func (s *Script) getLevel(t *Tag) int {
	value := t.Domain
	if value == "" {
		return 0
	}

	code := []rune(value)
	l := len(code)
	count := 0
	for i := 0; i < l; i++ {
		if code[i] == '.' {
			count++
		}
	}
	return count
}

func (s *Script) pushGSetter(i int, tag *Tag) {
	var p *GSetter = s.gsMap[tag.Value]
	if p == nil {
		p = &GSetter{}
		s.gsMap[tag.Value] = p
	}

	if i == 0 { //getter
		tag.Value = "get_" + tag.Value
		p.Getter = tag
	} else { //setter
		tag.Value = "set_" + tag.Value
		p.Setter = tag
	}
}

/**
 * 初始化名称
 * @param name
 * @param lst
 * @return
 * @throws Exception
 */
func (s *Script) initClass(name string, data string) string {
	ms := &MScript{}
	ms.ReadFromString(data)
	code := ""
	if len(s.eMap) > 0 {
		for _, value := range s.eMap {
			if Index(value, ".") == -1 {
				he := GetSingle(s.hMap, value)
				if he != nil {
					value = he.Name
				} else {
					fmt.Println("Script.go", "未找到对应类", value)
				}

			}
			ft := &UI{SYSTEM_PATH: s.jus.SYSTEM_PATH, CLASS_PATH: s.jus.CLASS_PATH}
			if err := ft.CreateFromParent(s.root, "", nil, strings.TrimSpace(value), s.jus); err == nil {
				if ft.IsScript() {
					code += "var __UP__ = new " + ft.ReadHTML().ToString() + ";\r\n"
				} else {
					tHTML := ft.ReadHTML()
					script := &HTMLScript{}
					script.CreateFrom(s.jus, s.jus.root, s.jus.domain, nil, "", "")
					scriptHTML := &HTML{}
					scriptHTML.ReadFromString("<script>" + script.ReadFromString(data) + "</script>")
					tHTML.Append(scriptHTML)
					return "\"" + Escape(tHTML.ToString()) + "\";\r\n"
				}
			} else {
				fmt.Println("Script.go ->", strings.TrimSpace(value), err)
			}
		}
	}
	return "function " + name + "(__FLAG__,__VALUE__){\r\n" +
		s.fromBuf +
		code +
		"var __inthis__ = this,__inpri__ = {};" +
		IfStr(len(code) == 0, "", "__EXTEND__(__inthis__,__UP__);") +
		s.initScriptFrom(ms, "__OBJECT__", "__inthis__", "__inpri__") + "\r\n" +
		"var __init__ = this.init || __inpri__.init;" +
		"if(__init__){" +
		"__init__.apply(this,__FLAG__ == __\\$__ ? __VALUE__ : arguments);" +
		"}" +
		"}"

}

/**
 * 读Script内容
 *
 * @param data
 * @throws IOException
 * @throws Exception
 */
func (s *Script) ReadFromString(script string) string {
	if len(script) == 0 {
		return ""
	}
	out := ""
	templ, err := GetCode(s.jus.SYSTEM_PATH + "/batch/j.ms")
	templ = strings.Replace(templ, "{@CLASS_NAME}", "//@ sourceURL=[ES]"+s.className+"\r\n", -1)
	if err != nil {
		return ""
	}

	s.mjs = &MScript{}
	s.mjs.ReadFromString(script)
	templ = strings.Replace(templ, "{@jscode}", s.initScript(s.mjs), -1)

	out += templ
	tmp := templ

	if s.extendScript != "" {
		s.mjs = &MScript{}
		s.mjs.ReadFromString(s.extendScript)
		templ = strings.Replace(tmp, "{@jscode}", s.initScript(s.mjs), -1)
		out += templ
	}

	if len(s.hMap) > 0 {
		for _, v := range s.hMap {
			s.jus.PushImportScript(v)
		}
	}

	return out
}

func (s *Script) loadClass(path string) string {
	className := strings.TrimSpace(Substring(path, 0, Index(path, "(")))
	if className[0] == '?' {
		return ""
	}
	tmpName := ""

	if Index(className, ".") == -1 {
		he := GetSingle(s.hMap, className)
		if he == nil {
			tmpName = ""
		} else {
			tmpName = he.Value
		}
	} else {
		//s.hMap[Substring(className, LastIndex(className, ".")+1, -1)] = &Attr{className, ""}
		s.jus.PushImportScript(&Attr{className, ""})
		tmpName = className
	}

	return IfStr(tmpName != "", "getModule('"+tmpName+"',__APPDOMAIN__)", "")
}

/**
 * 导入js数据
 *
 * @return
 */
func (s *Script) includeJs(path string) string {
	return ""
}

/**
 * 编译
 */
func Escape(value string) string {

	code := []rune(value)
	tmp := make([]rune, 0, len(code))

	for _, v := range code {
		if v == '"' {
			tmp = append(tmp, '\\')
		} else if v == '\\' {
			tmp = append(tmp, '\\')
		}
		if v == '\r' {
			tmp = append(tmp, '\\')
			tmp = append(tmp, 'r')
			continue
		}
		if v == '\n' {
			tmp = append(tmp, '\\')
			tmp = append(tmp, 'n')
			continue
		}
		if v == '\b' {
			tmp = append(tmp, '\\')
			tmp = append(tmp, 'b')
			continue
		}
		if v == 1 {
			tmp = append(tmp, '\\')
			tmp = append(tmp, '1')
			continue
		}
		tmp = append(tmp, v)
	}
	return strings.Replace(string(tmp), "</script", "<\" + \"/script", -1)
}
