// APIList.go
package util

import (
	. "jus"
	. "jus/str"
	"strings"
)

//修饰符
type Modifier struct {
	IsStatic     bool
	Setter       bool
	Getter       bool
	Description  string
	Param        []*Var
	Name         string
	FunctionType string
}

type Note struct {
	am    map[string]string
	hm    map[string]string
	code  []rune
	lst   []string
	_note string
}

func (a *Note) ReadFromString(note string) {
	a.am = make(map[string]string, 0)
	a.hm = make(map[string]string, 0)
	a.code = []rune(note)
	sb := make([]rune, 0)
	var ch rune
	position := 0
	for position < len(a.code) {
		ch = a.code[position]
		position++
		if ch == ' ' || ch == '\t' {
			if len(sb) > 0 {
				a.lst = append(a.lst, string(sb))
				sb = sb[0:0]
			}
			continue
		}
		sb = append(sb, ch)
	}
	if len(sb) > 0 {
		a.lst = append(a.lst, string(sb))
		sb = sb[0:0]
	}

	//01整理
	sb = sb[0:0]
	tsb := ""
	cmt := ""
	position = 0
	tag := ""
	name := ""
	value := ""
	for position < len(a.lst) {
		tag = a.lst[position]
		position++
		if "@param" == tag {
			name = a.lst[position]
			position++
			for position < len(a.lst) {
				tag = a.lst[position]
				position++
				if CharAt(tag, 0) == "@" {
					position--
					break
				}
				tsb += tag
			}
			value = tsb
			tsb = ""
			a.hm[name] = value
		} else if "@author" == tag {
			value = a.lst[position]
			position++
			a.am["author"] = value
		} else if "@version" == tag {
			value = a.lst[position]
			position++
			a.am["version"] = value
		} else if "@return" == tag {
			value = a.lst[position]
			position++
			a.am["return"] = value
		} else {
			cmt += tag + "\n"

		}
	}
	a._note = strings.TrimSpace(cmt)
}

func (a *Note) GetNote() string {
	return a._note
}

func (a *Note) GetParamNote(note string) string {
	note = a.hm[note]
	return IfStr(note != "", note, "无信息")
}

func (a *Note) GetAttr(name string) string {
	name = a.am[name]
	return IfStr(name != "", name, "无信息")
}

type APIlist struct {
	className string
	IsScript  bool //是Script文件
	script    *MScript
	sb        string
	html      *HTML
	style     string
	attr      []*Modifier
}

/**
 * 通过什么创建APIList
 */
func (a *APIlist) CreateFrom(jus *UIServer, className string) error {
	a.className = className
	a.style = `<style>
			table {
				font-family: verdana,arial,sans-serif;
				font-size:12px;
				color:#333333;
				border-width: 1px;
				border-color: #666666;
				border-collapse: collapse;
				width:100%;
			}
			table th {
				border-width: 1px;
				padding: 8px;
				border-style: solid;
				border-color: #666666;
				background-color: #dedede;
			}
			table td {
				border-width: 1px;
				padding: 8px;
				border-style: solid;
				border-color: #666666;
				background-color: #ffffff;
			}
			
			
			a{
				font-size:13px;
			    color: #0000CC;
			    text-decoration: none;
			    font-weight: bold;
			}
			
			ul{
				margin:0px;
				margin-top:10px;
				list-style-type: unset;
			}
			
			div.address{
				font-size:12pt;
				color:#dddddd;
				padding-left:10px;
				padding-top:5px;
			}
			
		</style>`
	a.IsScript = true
	script := ""
	path := strings.Replace(className, ".", "/", -1)
	path = IfStr(CharAt(path, 0) == "$", jus.SysPath+"/src/"+Substring(path, 1, -1), jus.RootPath+"/"+path)
	file := path + ".ui"
	if Exist(file) {
		a.IsScript = false
		code, err := GetCode(file)
		if err != nil {
			return err
		}
		//剖析script
		html := &HTML{}
		html.ReadFromString(code)

		for _, v := range html.GetElementsByTagName("script") {
			script += v.ToString()
		}
		a.html = html
	}

	file = path + ".es"
	if Exist(file) {
		code, err := GetCode(file)
		if err != nil {
			return err
		}
		script += "\n" + code
	}

	a.script = &MScript{}
	a.script.ReadFromString(script)
	if a.IsScript {
		a.initClass(path, a.script.GetClass())
	} else {
		a.init(path)
	}

	return nil
}

//添加属性
func (a *APIlist) appendAttr(f *Function) {
	var p *Modifier = nil
	for _, v := range a.attr {
		if v.Name == f.Name && v.IsStatic == f.IsStatic {
			p = v
		}
	}
	if p == nil {
		p = &Modifier{Name: f.Name, IsStatic: f.IsStatic, Param: f.Param}
		a.attr = append(a.attr, p)
	}
	if f.IsGet {
		p.Getter = true
		p.FunctionType = f.FunctionType
	}
	if f.IsSet {
		p.Setter = true
	}
	if f.Note != "" {
		p.Description += f.Note + "\r\n"
	}

}

func (a *APIlist) init(name string) {
	tsb := ""
	sb := a.style
	sb += "<div style='border-bottom:solid 1px #000000;font-size:26px;font-weight:bold;padding:10px;'>" + a.className + "</div>" //name
	sb += "<div class='address'>" + name + "</div>"
	sb += "<div style='padding:10px;padding-bottom:20px;'>"
	if a.html != nil {
		child := a.html.Filter("!")
		for _, v := range child {
			tsb += v.Text()
		}
		tsb = strings.Replace(tsb, "\n", "<br/>", -1)
		sb += tsb
	}
	sb += "</div>"
	//
	sb += "<div style='padding-left:10px;padding-right:10px'>"

	js := a.script
	svc := js.GetVar(true, true)
	if len(svc) > 0 {
		sb += "<b style='padding-bottom:5px;display:block;'>静态属性</b>"
	}
	sb += "<table>"
	for _, v := range svc {
		sb += ("<tr><td>")
		sb += IfStr(v.IsStatic, "static ", "") + "<a href='javascript:void(0);'>" + v.Name + "</a>" + IfStr(v.VarType != "", " : <a href='javascript:void(0);' style='color:#888888'>"+v.VarType+"</a>", "") + "<br/>"
		sb += v.Note
		sb += "</td></tr>"
	}
	sb += "</table><br/>"

	svc = js.GetVar(true, false)
	fc := js.GetFunctionAndStatic(true, false)
	if len(svc) > 0 || len(fc) > 0 {
		sb += "<b style='padding-bottom:5px;display:block;'>公共属性</b>"
	}
	sb += ("<table>")
	for _, v := range svc {
		sb += ("<tr><td>")
		sb += IfStr(v.IsStatic, "static ", "") + "<a href='javascript:void(0);'>" + v.Name + "</a>" + IfStr(v.VarType != "", " : <a href='javascript:void(0);' style='color:#888888'>"+v.VarType+"</a>", "") + "<br/>"
		sb += v.Note
		sb += "</td></tr>"
	}

	for _, f := range fc { //遍历所有setter getter 函数
		if f.IsSet || f.IsGet {
			a.appendAttr(f)
		}
	}
	for _, f := range a.attr {
		note := &Note{}
		note.ReadFromString(f.Description)
		sb += "<tr><td>"
		sb += ("<span>")
		sb += IfStr(f.IsStatic, "static ", "")
		ut := ""
		if f.Setter && !f.Getter {
			ut = "<b style='color:#cc0000'>[只写]</b>"
			sb += "<a href='javascript:void(0);' style='font-size:14px;'>" + f.Name + "</a> : <b style='color:#888888;'>" + f.Param[0].VarType + "</b>"
		} else if f.Getter && !f.Setter {
			ut = "<b style='color:#00cc00'>[只读]</b>"
			sb += "<a href='javascript:void(0);' style='font-size:14px;'>" + f.Name + "</a>" + IfStr(f.FunctionType != "", " : <b style='color:#888888;'>"+f.FunctionType+"</b>", "")
		} else {
			sb += "<a href='javascript:void(0);' style='font-size:14px;'>" + f.Name + "</a>" + IfStr(f.FunctionType != "", " : <b style='color:#888888;'>"+f.FunctionType+"</b>", "")
		}

		sb += "</span>\r\n"
		sb += "<div style='padding:10px 0px;padding-bottom:0px;font-size:13px;'>"
		sb += IfStr(ut == "", "", ut+" ") + Replace(note.GetNote(), "\n", "<br/>")
		sb += "</div>"
		sb += "</td></tr>"
	}

	sb += ("</table>")
	sb += ("<br/>")

	fc = js.GetFunctionAndStatic(true, true)

	if len(fc) > 0 {
		sb += "<b style='padding-bottom:5px;display:block;'>静态方法</b>"
	}
	sb += "<table>"

	for _, f := range fc {
		if f.IsSet || f.IsGet {
			continue
		}
		note := &Note{}
		note.ReadFromString(f.Note)
		sb += "<tr><td>"
		sb += ("<span>")
		sb += IfStr(f.IsStatic, "static ", "")

		sb += "<a href='javascript:void(0);' style='font-size:14px;'>" + f.Name + "</a> " + "(" + paramToString(f.Param) + ")" + IfStr(f.FunctionType != "", " : <b style='color:#888888;'>"+f.FunctionType+"</b>", "")

		sb += "</span>\r\n"
		sb += "<div style='padding:10px;padding-bottom:0px;font-size:13px;'>"
		sb += Replace(note.GetNote(), "\n", "<br/>")
		sb += "</div>"
		sb += "<ul>\r\n"
		param := f.Param
		for _, s := range param {
			sb += "<li><span style='font-weight:bold'>" + s.Name + "</span> " + note.GetParamNote(s.Name) + "</li>"
		}
		sb += "<li><span style='font-weight:bold;color:#992200'>返回值</span> " + note.GetAttr("return") + "</li>"
		sb += "</ul>"
		sb += "</td></tr>"
	}
	sb += "</table><br/>"

	cc := js.GetConstructor()
	if cc != nil {
		sb += "<b style='padding-bottom:5px;display:block;'>构造方法</b>"
		sb += "<table>"
		note := &Note{}
		note.ReadFromString(cc.Note)
		sb += "<tr><td>"
		sb += ("<span>")
		sb += IfStr(cc.IsStatic, "static ", "")
		if cc.IsSet {
			sb += "<b>set</b> <a href='javascript:void(0);' style='font-size:14px;'>" + cc.Name + "</a> <b>" + cc.Param[0].VarType + "</b>"
		} else if cc.IsGet {
			sb += "<b>get</b> <a href='javascript:void(0);' style='font-size:14px;'>" + cc.Name + "</a>" + IfStr(cc.FunctionType != "", " <b style='color:#888888;'>"+cc.FunctionType+"</b>", "")
		} else {
			sb += "<a href='javascript:void(0);' style='font-size:14px;'>" + cc.Name + "</a> " + "(" + paramToString(cc.Param) + ")" + IfStr(cc.FunctionType != "", " : <b style='color:#888888;'>"+cc.FunctionType+"</b>", "")
		}

		sb += "</span>\r\n"
		sb += "<div style='padding:10px;padding-bottom:0px;font-size:13px;'>"
		sb += Replace(note.GetNote(), "\n", "<br/>")
		sb += "</div>"
		sb += "<ul>\r\n"
		param := cc.Param
		for _, s := range param {
			sb += "<li><span style='font-weight:bold'>" + s.Name + "</span> " + note.GetParamNote(s.Name) + "</li>"
		}
		sb += "</ul>"
		sb += "</td></tr>"
		sb += "</table><br/>"
	}

	fc = js.GetFunctionAndStatic(true, false)

	if len(fc) > 0 {
		sb += "<b style='padding-bottom:5px;display:block;'>公共方法</b>"
	}
	sb += "<table>"
	for _, f := range fc {
		if f.IsGet || f.IsSet {
			continue
		}
		note := &Note{}
		note.ReadFromString(f.Note)
		sb += "<tr><td>"
		sb += ("<span>")
		sb += IfStr(f.IsStatic, "static ", "")
		sb += "<a href='javascript:void(0);' style='font-size:14px;'>" + f.Name + "</a> " + "(" + paramToString(f.Param) + ")" + IfStr(f.FunctionType != "", " : <b style='color:#888888;'>"+f.FunctionType+"</b>", "")
		sb += "</span>\r\n"
		sb += "<div style='padding:10px;padding-bottom:0px;font-size:13px;'>"
		sb += Replace(note.GetNote(), "\n", "<br/>")
		sb += "</div>"
		sb += "<ul>\r\n"
		param := f.Param
		for _, s := range param {
			sb += "<li><span style='font-weight:bold'>" + s.Name + "</span> " + note.GetParamNote(s.Name) + "</li>"
		}
		sb += "<li><span style='font-weight:bold;color:#992200'>返回值</span> " + note.GetAttr("return") + "</li>"
		sb += "</ul>"
		sb += "</td></tr>"
	}
	sb += "</table></div>"
	a.sb = sb
}

/**
 * 初始化类
 */
func (a *APIlist) initClass(path string, cls []*Class) {

	sb := a.style
	for _, v := range cls {
		sb += "<div style='border-bottom:solid 1px #000000;font-size:26px;font-weight:bold;padding:10px;'>" + IfStr(v.IsInner, v.Name, a.className+" : "+v.Name) + "</div>"
		sb += "<div class='address'>" + path + "</div>"
		sb += "<div style='padding:10px;padding-bottom:20px;'>"
		sb += Replace(v.Note, "\r", "<br/>")
		sb += "</div>"
		sb += "<div style='padding-left:10px;padding-right:10px'>"
		fv := a.script.GetVarByClassName(v.Name, true, true)
		if len(fv) > 0 {
			sb += "<b style='padding-bottom:5px;display:block;'>静态属性</b>"
		}
		sb += "<table>"
		for _, v := range fv {
			sb += ("<tr><td>")
			sb += IfStr(v.IsStatic, "static ", "") + "<a href='javascript:void(0);'>" + v.Name + "</a>" + IfStr(v.VarType != "", " : <a href='javascript:void(0);' style='color:#888888'>"+v.VarType+"</a>", "") + "<br/>"
			sb += v.Note
			sb += "</td></tr>"
		}
		sb += "</table><br/>"
		fv = a.script.GetVarByClassName(v.Name, true, false)
		fc := a.script.GetFunctionByClassName(v.Name, true)
		if len(fv) > 0 || len(fc) > 0 {
			sb += "<b style='padding-bottom:5px;display:block;'>公共属性</b>"
		}
		sb += ("<table>")

		for _, v := range fv {
			sb += ("<tr><td>")
			sb += IfStr(v.IsStatic, "static ", "") + "<a href='javascript:void(0);'>" + v.Name + "</a>" + IfStr(v.VarType != "", " : <a href='javascript:void(0);' style='color:#888888'>"+v.VarType+"</a>", "") + "<br/>"
			sb += v.Note
			sb += "</td></tr>"
		}

		for _, f := range fc { //遍历所有setter getter 函数
			if f.IsSet || f.IsGet {
				a.appendAttr(f)
			}
		}
		for _, f := range a.attr {
			note := &Note{}
			note.ReadFromString(f.Description)
			sb += "<tr><td>"
			sb += ("<span>")
			sb += IfStr(f.IsStatic, "static ", "")
			ut := ""
			if f.Setter && !f.Getter {
				ut = "<b style='color:#cc0000'>[只写]</b>"
				sb += "<a href='javascript:void(0);' style='font-size:14px;'>" + f.Name + "</a> : <b style='color:#888888;'>" + f.Param[0].VarType + "</b>"
			} else if f.Getter && !f.Setter {
				ut = "<b style='color:#00cc00'>[只读]</b>"
				sb += "<a href='javascript:void(0);' style='font-size:14px;'>" + f.Name + "</a>" + IfStr(f.FunctionType != "", " : <b style='color:#888888;'>"+f.FunctionType+"</b>", "")
			} else {
				sb += "<a href='javascript:void(0);' style='font-size:14px;'>" + f.Name + "</a>" + IfStr(f.FunctionType != "", " : <b style='color:#888888;'>"+f.FunctionType+"</b>", "")
			}

			sb += "</span>\r\n"
			sb += "<div style='padding:10px 0px;padding-bottom:0px;font-size:13px;'>"
			sb += IfStr(ut == "", "", ut+" ") + Replace(note.GetNote(), "\n", "<br/>")
			sb += "</div>"
			sb += "</td></tr>"
		}
		sb += "</table><br/>"

		fc = a.script.GetFunctionAndStaticByClassName(v.Name, true, true)

		if len(fc) > 0 {
			sb += "<b style='padding-bottom:5px;display:block;'>静态方法</b>"
		}
		sb += "<table>"
		for _, f := range fc {
			note := &Note{}
			note.ReadFromString(f.Note)
			sb += "<tr><td>"
			sb += ("<span>")
			sb += IfStr(f.IsStatic, "static ", "")
			if f.IsSet {
				sb += "<b>set</b> <a href='javascript:void(0);' style='font-size:14px;'>" + f.Name + "</a> <b>" + f.Param[0].VarType + "</b>"
			} else if f.IsGet {
				sb += "<b>get</b> <a href='javascript:void(0);' style='font-size:14px;'>" + f.Name + "</a>" + IfStr(f.FunctionType != "", " <b style='color:#888888;'>"+f.FunctionType+"</b>", "")
			} else {
				sb += "<a href='javascript:void(0);' style='font-size:14px;'>" + f.Name + "</a> " + "(" + paramToString(f.Param) + ")" + IfStr(f.FunctionType != "", " : <b style='color:#888888;'>"+f.FunctionType+"</b>", "")
			}

			sb += "</span>\r\n"
			sb += "<div style='padding:10px;padding-bottom:0px;font-size:13px;'>"
			sb += Replace(note.GetNote(), "\n", "<br/>")
			sb += "</div>"
			sb += "<ul>\r\n"
			param := f.Param
			for _, s := range param {
				sb += "<li><span style='font-weight:bold'>" + s.Name + "</span> " + note.GetParamNote(s.Name) + "</li>"
			}
			sb += "<li><span style='font-weight:bold;color:#992200'>返回值</span> " + note.GetAttr("return") + "</li>"
			sb += "</ul>"
			sb += "</td></tr>"
		}
		sb += "</table><br/>"

		cc := a.script.GetConstructorByClassName(v.Name)
		if cc != nil {
			sb += "<b style='padding-bottom:5px;display:block;'>构造方法</b>"
			sb += "<table>"
			note := &Note{}
			note.ReadFromString(cc.Note)
			sb += "<tr><td>"
			sb += ("<span>")
			sb += IfStr(cc.IsStatic, "static ", "")
			if cc.IsSet {
				sb += "<b>set</b> <a href='javascript:void(0);' style='font-size:14px;'>" + cc.Name + "</a> <b>" + cc.Param[0].VarType + "</b>"
			} else if cc.IsGet {
				sb += "<b>get</b> <a href='javascript:void(0);' style='font-size:14px;'>" + cc.Name + "</a>" + IfStr(cc.FunctionType != "", " <b style='color:#888888;'>"+cc.FunctionType+"</b>", "")
			} else {
				sb += "<a href='javascript:void(0);' style='font-size:14px;'>" + cc.Name + "</a> " + "(" + paramToString(cc.Param) + ")" + IfStr(cc.FunctionType != "", " : <b style='color:#888888;'>"+cc.FunctionType+"</b>", "")
			}

			sb += "</span>\r\n"
			sb += "<div style='padding:10px;padding-bottom:0px;font-size:13px;'>"
			sb += Replace(note.GetNote(), "\n", "<br/>")
			sb += "</div>"
			sb += "<ul>\r\n"
			param := cc.Param
			for _, s := range param {
				sb += "<li><span style='font-weight:bold'>" + s.Name + "</span> " + note.GetParamNote(s.Name) + "</li>"
			}
			sb += "</ul>"
			sb += "</td></tr>"
			sb += "</table><br/>"
		}

		fc = a.script.GetFunctionByClassName(v.Name, true)
		if len(fc) > 0 {
			sb += "<b style='padding-bottom:5px;display:block;'>公共方法</b>"
		}
		sb += "<table>"
		for _, f := range fc {
			if f.IsGet || f.IsSet {
				continue
			}
			note := &Note{}
			note.ReadFromString(f.Note)
			sb += "<tr><td>"
			sb += ("<span>")
			sb += IfStr(f.IsStatic, "static ", "")
			sb += "<a href='javascript:void(0);' style='font-size:14px;'>" + f.Name + "</a> " + "(" + paramToString(f.Param) + ")" + IfStr(f.FunctionType != "", " : <b style='color:#888888;'>"+f.FunctionType+"</b>", "")

			sb += "</span>\r\n"
			sb += "<div style='padding:10px;padding-bottom:0px;font-size:13px;'>"
			sb += Replace(note.GetNote(), "\n", "<br/>")
			sb += "</div>"
			sb += "<ul>\r\n"
			param := f.Param
			for _, s := range param {
				sb += "<li><span style='font-weight:bold'>" + s.Name + "</span> " + note.GetParamNote(s.Name) + "</li>"
			}
			sb += "<li><span style='font-weight:bold;color:#992200'>返回值</span> " + note.GetAttr("return") + "</li>"
			sb += "</ul>"
			sb += "</td></tr>"
		}
		sb += "</table></div>"
	}
	a.sb = sb
}

func paramToString(param []*Var) string {
	if len(param) > 0 {
		sb := ""
		for _, s := range param {
			sb += "<b>" + s.Name + "</b>" + IfStr(s.VarType != "", " : <span style='color:#888888;'>"+s.VarType+"</span>", "") + IfStr(s.Value != "", " "+s.Value, "") + ", "
		}
		sb = Substring(sb, 0, StringLen(sb)-2)
		sb += ("</b>")
		return sb
	}
	return ""
}

func (a *APIlist) ToString() string {
	return a.sb
}
