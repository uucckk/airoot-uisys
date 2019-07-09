// Script.go

package util

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	. "jus"
	. "jus/str"
	. "jus/tool"
	"strings"
)

//--------------------------------Script----------------------------------------
type HTMLScript struct {
	jus              *JUS
	root             string
	hMap             []*Attr //导入的类文件
	gsMap            map[string]*GSetter
	domain           string
	constructorValue *Attr
	innerValue       string
	extendScript     string
	mjs              *MScript
	isScript         bool
}

func (s *HTMLScript) CreateFrom(jus *JUS, root string, domain string, constructorValue *Attr, innerValue string, extendScript string) *HTMLScript {
	s.jus = jus
	s.root = root
	s.domain = domain
	s.constructorValue = constructorValue
	s.innerValue = innerValue
	s.hMap = make([]*Attr, 0)
	s.extendScript = extendScript
	s.gsMap = make(map[string]*GSetter, 10)
	return s
}

func (s *HTMLScript) initScript(js *MScript) string {
	return s.initScriptFrom(js, "____", "____")

}

/**
 * 初始化Script语句
 * @param script
 * @return
 * @throws Exception
 */
func (s *HTMLScript) initScriptFrom(js *MScript, _this_ string, _pri_ string) string {
	out := bytes.NewBufferString("")
	tmp := bytes.NewBufferString("")
	newString := bytes.NewBufferString("")
	var hObj *HTMLObject = nil
	var tjs *MScript = nil
	lst := js.GetJUIScriptData()
	//for k, v := range lst {
	//	fmt.Println(k, v.Domain, ">>", v.Value)
	//}

	tl := make([]*Tag, 0, 1000)
	tlt := make([]*Tag, 0, 1000)

	p := 0
	var t *Tag = nil
	var f *Tag = nil
	var param *Tag = nil
	level := 0
	//00.转化字符串到JUS的$形式
	for p < len(lst) {
		t = lst[p]
		p++
		if t.TagType == 1 { //如果是字符串
			t.Value = ToJUSString(Replace(t.Value, "@this", "$"))
		}

	}
	p = 0
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
			}

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
				tmp.Reset()
				p--
				for p < len(lst) {
					f = lst[p]
					p++
					if f.TagType < -1 {
						continue
					}
					tmp.WriteString(f.Value)
					if t.TagType == 3 && "{" == f.Value {
						level++
					} else if f.TagType == 3 && "}" == f.Value {
						level--
						if level == 0 {
							break
						}
					}
				}
				s.jus.AddStaticCode(s.jus.className, "__STATIC__", " = function()"+tmp.String()+";")
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
				newString.Reset()
				if t.IsPublic {
					if t.IsStatic && !t.IsGet && !t.IsSet {
						newString.WriteString("__WINDOW__[__APPDOMAIN__]['")
						newString.WriteString(s.jus.className)
						newString.WriteString("'].")
					} else {
						newString.WriteString(_this_)
						newString.WriteString(".")
					}
					newString.WriteString(t.Value)

				} else {
					if t.IsStatic && !t.IsGet && !t.IsSet {
						newString.WriteString("__WINDOW__[__APPDOMAIN__]['")
						newString.WriteString(s.jus.className)
						newString.WriteString("'].")
					} else {
						newString.WriteString(_pri_)
						newString.WriteString(".")
					}
					newString.WriteString(t.Value)
				}
				tl = append(tl, &Tag{Value: newString.String(), TagType: 0})
			} else if t.Domain == "" {
				if s.jus != nil {
					hObj = s.jus.GetDefine(t.Value)
				}

				if hObj != nil {
					if t.Value[0] == '$' {
						t.Value = t.Value[1:]
					}
					if hObj.HTMLObjectType == -1 {
						t.Value = "dom"
					} else {
						t.Value = "window[__NAME__+'" + t.Value + "']" //hObj.Name
					}

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
			md5Ctx := md5.New()
			md5Ctx.Write([]byte(t.Value))
			cipherStr := md5Ctx.Sum(nil)
			bs := hex.EncodeToString(cipherStr)
			ft := &JUS{SYSTEM_PATH: s.jus.SYSTEM_PATH, CLASS_PATH: s.jus.CLASS_PATH}
			for _, v := range s.jus.pkgMap {
				t.Value = "<@import value='" + v + "'/>" + t.Value
			}
			ft.CreateFromString(s.root, "", nil, t.Value, bs, nil)
			tl = append(tl, &Tag{Value: "getModule(\"" + bs + "\",__APPDOMAIN__)", TagType: 0})
			//s.jus.ToFormatLine("I", bs, "H"+ft.ToFormatString(), sb)
			s.jus.GetRoot().scriptElementBuffer = append(s.jus.GetRoot().scriptElementBuffer, &ScriptElement{"I", bs, "H", ft.ToFormatString()})
			continue
		}

		// 3.import
		if t.IsKeyWord && "import" == t.Value {
			tmp.Reset()
			point := -1 //类文件名
			at := 0
			isFrom := false
			for p < len(lst) {
				f = lst[p]
				p++

				if f.TagType == -1 || f.TagType == 5 {
					continue
				}
				if f.TagType == 9 {
					point = p
				}
				if f.TagType == 4 {
					if point == -1 {
						point = at
					}
					break
				}
				if f.IsKeyWord && f.Value == "from" { //说明要用commandJS规范导读
					f.Value = "\002"
					isFrom = true
				}
				tmp.WriteString(f.Value)
				at = p - 1
			}
			Single(&s.hMap, &Attr{lst[point].Value, tmp.String()})
			s.jus.PushImportScript(&Attr{lst[point].Value, tmp.String()})
			if isFrom {
				tl = append(tl, &Tag{Value: ImportFrom(s.jus.className, tmp.String()), TagType: 1})
			}
			continue
		}

		// 5.include
		if t.IsKeyWord && "include" == t.Value {
			tmp.Reset()
			for p < len(lst) {
				f = lst[p]
				p++
				if f.TagType == 1 {
					tmp.WriteString(f.Value)
					break
				}
			}

			tl = append(tl, &Tag{Value: s.includeJs(tmp.String()), TagType: 0})
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
			tmp.Reset()
			for p < len(tl) {
				f = tl[p]
				p++
				if "(" == f.Value {
					break
				}
				tmp.WriteString(f.Value)
			}

			newTmp := s.loadClass(tmp.String())
			if newTmp != "" {
				tlt = append(tlt, &Tag{Value: newTmp, TagType: 0})
			} else {
				tlt = append(tlt, t)
				tlt = append(tlt, &Tag{Value: Replace(tmp.String(), "?", ""), TagType: 1})
			}
			tlt = append(tlt, f)
			continue
		}
		if t.IsKeyWord && "@this" == t.Value {
			t.Value = _this_
		} else if t.IsKeyWord && "@res" == t.Value {
			t.Value = "\"" + s.jus.resPath + "/" + s.jus.relativePath + ".RES/\""
		} else if t.Value[0] == '@' {
			t.Value = s.jus.SERVER.GetServerVar(t.Value)
		} else if t.IsKeyWord && "this" == t.Value {
			tlt = append(tlt, t)
			if s.getLevel(t) == 1 {
				//t.Value = _pri_
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
										param.Value = _this_
									} else {
										param.Value = _pri_
									}
								}
								if s.mjs != js && a == nil {
									a = s.mjs.GetDefine("class").Get(t.Value)
									if a != nil {
										if a.IsPublic {
											param.Value = _this_
										} else {
											param.Value = _pri_
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
	var buffer []*Tag = make([]*Tag, 0, 1000)

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
						tlt = append(tlt, &Tag{Value: IfStr(s.isScript, IfStr(t.IsPublic, _this_+".", _pri_+".")+t.Value+" = function", IfStr(t.IsPublic, _this_+".", _pri_+".")+t.Value+" = function"), TagType: 0})
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

				tlt = append(tlt, &Tag{Value: IfStr(t.IsPublic, _this_+".", _pri_+".") + t.Value + " ", TagType: 0})
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
			tmp.Reset()
			if t.IsFunction {
				tmp.WriteString("=function")
				for p < len(tl) {
					f = tl[p]
					p++
					if f.TagType < -1 {
						continue
					}
					tmp.WriteString(f.Value)
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
				//tlt = append(tlt, &Tag{Value: IfStr(t.IsSet, "var", _this_+".") + t.Value + " = __WINDOW__[__APPDOMAIN__]['" + s.jus.className + "']." + t.Value + ";", TagType: 0})
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
					tmp.WriteString(f.Value)
				}
			}
			if s.jus != nil {
				s.jus.AddStaticScript(s.jus.className, t.Value, tmp.String())
			}
			continue
		}
		tlt = append(tlt, t)
	}
	tl = tl[0:0]
	tl = appendArray(tl, tlt)
	tlt = tlt[0:0]

	out.WriteString(js.ToStringFrom(tl))
	//处理Getter Setter
	var pgs *GSetter = nil
	tsb := bytes.NewBufferString("")
	for name, value := range s.gsMap {
		pgs = value
		tsb.WriteString("Object.defineProperty(" + _this_ + ",'" + name + "',{")
		if pgs.Setter != nil {
			tsb.WriteString("set:")
			if pgs.Setter.IsStatic {
				tsb.WriteString("__WINDOW__[__APPDOMAIN__]['" + s.jus.className + "']." + pgs.Setter.Value)
			} else {
				tsb.WriteString(pgs.Setter.Value)
			}
		}
		if pgs.Getter != nil {
			if pgs.Setter != nil {
				tsb.WriteString(",")
			}
			tsb.WriteString("get:")
			if pgs.Getter.IsStatic {
				tsb.WriteString("__WINDOW__[__APPDOMAIN__]['" + s.jus.className + "']." + pgs.Getter.Value)
			} else {
				tsb.WriteString(pgs.Getter.Value)
			}
		}
		tsb.WriteString(",enumerable:true});\r\n")
	}
	out.WriteString(tsb.String())

	return out.String()

}

/**
 * 获取变量是第几层
 * @param t
 * @return
 */
func (s *HTMLScript) getLevel(t *Tag) int {
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

func (s *HTMLScript) pushGSetter(i int, tag *Tag) {
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
func (s *HTMLScript) initClass(name string, data string) string {
	ms := &MScript{}
	ms.ReadFromString(data)
	if s.isScript {

		return "function " + name + "(__VALUE__){var __inthis__ = this,__inpri__ = {};" + s.initScriptFrom(ms, "__inthis__", "__inpri__") + "\r\n" +
			"var __init__ = this.init || __inpri__.init;" +
			"if(__init__){" +
			"__init__.apply(this,__VALUE__);" +
			"}" +
			"}"
	}

	return "function " + name + "(){var __inthis__ = this,__inpri__ = {};" + s.initScriptFrom(ms, "__inthis__", "__inpri__") + "\r\n" +
		"var __init__ = this.init || __inpri__.init;" +
		"if(__init__){" +
		"__init__.apply(this,arguments);" +
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
func (s *HTMLScript) ReadFromString(script string) string {
	if len(script) == 0 {
		return ""
	}
	out := bytes.NewBufferString("")
	msPath := ""
	if s.isScript {
		msPath = "/batch/j.ms"
	} else {
		msPath = "/batch/m.ms"
	}
	templ, err := GetCode(s.jus.SYSTEM_PATH + msPath)
	tmp := templ
	if err != nil {
		return ""
	}
	templ = strings.Replace(templ, "{@CLASS_NAME}", "//"+s.jus.className+"\r\n", -1)
	templ = strings.Replace(templ, "{@domain}", s.jus.domain, -1)
	templ = strings.Replace(templ, "{@Base}", "\b", -1)

	if s.constructorValue != nil {
		templ = strings.Replace(templ, "{@value}", s.constructorValue.Value, -1)
	} else {
		templ = strings.Replace(templ, "{@value}", "", -1)
	}

	s.mjs = &MScript{}
	s.mjs.ReadFromString(script)
	templ = strings.Replace(templ, "{@jscode}", "var context = {value:\""+Escape(s.innerValue)+"\"}\r\n"+s.initScript(s.mjs), -1)

	s.jus.ToFormatLine("M", s.jus.className, templ, out)
	//加入执行列表
	s.jus.AddRun(&RunElem{Type: "S", Name: s.jus.domain, Value: s.jus.className})

	if s.extendScript != "" {
		s.mjs = &MScript{}
		s.mjs.ReadFromString(s.extendScript)
		templ = strings.Replace(tmp, "{@CLASS_NAME}", "//"+s.jus.className, -1)
		templ = strings.Replace(templ, "{@jscode}", s.initScript(s.mjs), -1)
		E := s.jus.ToFormatLine("E", s.jus.className, templ, out) //E代表扩展代码
		//加入执行列表
		s.jus.AddRun(&RunElem{Type: "E", Name: s.jus.domain, Value: E})
	}
	if len(s.jus.CommandCode) > 0 {
		for _, v := range s.jus.CommandCode {
			s.jus.AddRun(&RunElem{Type: "C", Name: s.jus.domain, Value: v.Value})
		}
	}
	return out.String()
}

func (s *HTMLScript) FormatString(script string) string {
	if len(script) == 0 {
		return ""
	}
	s.mjs = &MScript{}
	s.mjs.ReadFromString(script)
	return s.initScript(s.mjs)
}

func (s *HTMLScript) loadClass(path string) string {
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
func (s *HTMLScript) includeJs(path string) string {
	return ""
}
