package util

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	. "uisys"
	. "uisys/str"
	. "uisys/tool"

	"github.com/dop251/goja"
)

//-------------------------------HTMLObject-------------------------------------
type HTMLObject struct {
	Name           string
	HTMLObjectType int //-1代表容器节点，0代表普通元素，1代表对象
}

//----------------------------------JUS-----------------------------------------
type UI struct {
	IsPublic            bool              //JS方法是否公开，默认不公开，但是如果html使用了@this关键字，则默认公开且不可改变。
	IsTest              bool              //是否为测试模块，如果是则打开测试开关*，#，如果不是测不显示测试代码
	IsSysLib            bool              //是否为库中组件
	SysLibDirs          map[string]string //存储为库中组件的lib目录
	Debug               bool              //判断是否被测试
	SERVER              *UIServer         //服务器引用
	dirPath             string            //所在目录地址
	path                string            //记录类的文件夹路径
	htmlPath            string            //html模块的绝对路径
	jsPath              string            //js模块的绝对路径
	cssPath             string            //css模块路径
	SYSTEM_PATH         string            //系统路径
	CLASS_PATH          string            //
	root                string            //工程路径
	parent              *UI               //
	scanLevel           int               //递归层级
	domain              string            //
	className           string            //
	relativePath        string            //相对路径
	node                *HTML             //此HTML节点
	innerContent        []*HTML           //此HTML节点的子元素Child
	contentToList       []*HTML           //节点为变量的储存列表
	contentTo           string            //内容信息变量添加到
	paramValue          *Attr             //
	innerValue          string            //内部代码转string
	innerModule         string            //内部模块
	html                *HTML             //
	extendsScriptBuffer string            //
	scriptBuffer        bytes.Buffer      //
	styleBuffer         bytes.Buffer      //
	cssBuffer           bytes.Buffer      //全局css属性
	pkgMap              map[string]string
	idMap               map[string]*HTMLObject
	staticScript        map[string][]*Attr
	staticCode          map[string][]*Attr
	styleCode           map[string]string
	CommandCode         []*Attr //指令集合
	extendFlag          bool
	style               *CSS
	css                 *CSS
	cssTag              map[string]string
	scriptFile          bool             //判断是否为独立JavaScript文件
	componentParam      []*Attr          //控件初始化参数
	componentCode       []*Attr          //控件默认代码
	scriptElement       map[string]*Attr //需要导入的头文件，类似与import
	scriptElementBuffer []*ScriptElement
	componentParams     []string                       //所有编译时初始化集合，只有顶级的元素可以接受
	count               int                            //自动化数量
	moduleMap           map[string]*Attr               //模块地图
	runList             []*RunElem                     //run列表，用于记录模块的执行顺序，非常重要的一个字段
	IsImport            string                         //是否为导入类
	headBuffer          bytes.Buffer                   //html Head标签
	pub                 string                         //发布模板，没有模板就是空
	defineMap           map[string](map[string]string) //自定义模块类映射
	defineClassMap      map[string]string              //自定义模块对应实体
}

/**
 * @param root			编译工程目录
 * @param domain		文件作用域
 * @param value			传递的实际参数
 * @param file			读取文件路径
 * @throws IOException
 */
func (j *UI) CreateFromString(root string, domain string, node *HTML, code string, className string, parent *UI, defLib *UI) error {
	j.parent = parent

	j.pkgMap = make(map[string]string, 10)
	j.idMap = make(map[string]*HTMLObject, 10)
	j.root = filepath.Clean(root)
	j.className = className
	j.contentToList = make([]*HTML, 0)
	if j.parent == nil {
		j.moduleMap = make(map[string]*Attr, 10)
		j.defineMap = make(map[string](map[string]string))
		j.defineClassMap = make(map[string]string)
	}
	if node != nil {
		j.node = node
		j.innerContent = node.Child()
	} else {
		j.innerContent = nil
	}

	if className == "" {
		return errors.New("ui.go -> className is empty.")
	}

	j.html = &HTML{}

	_, err := j.html.ReadFromString(CodeFx(code, j.IsTest)) //j.html.ReadFromString(j.scanMedia(t))
	if err != nil {
		return err
	}

	if domain == "" {
		j.domain = "\b"
	} else {
		j.domain = domain
	}

	if defLib != nil {
		j.GetRoot().defineMap = defLib.GetRoot().defineMap
		j.GetRoot().defineClassMap = defLib.GetRoot().defineClassMap
	}

	return nil
}

/**
 * @param root			编译工程目录
 * @param domain		文件作用域
 * @param value			传递的实际参数
 * @param file			读取文件路径
 * @throws IOException
 */
func (j *UI) CreateFrom(root string, domain string, node *HTML, className string) error {
	root = filepath.Clean(root)
	className = Replace(className, "/", ".")
	className = Replace(className, "\\", ".")
	className = TrimClassName(className)
	if j.scanLevel > 100 { //最大深度100层级
		return errors.New(j.GetRoot().className + " : scan module over stack when on create \"" + className + "\"")
	}

	j.pkgMap = make(map[string]string, 10)
	j.idMap = make(map[string]*HTMLObject, 10)
	j.root = root
	j.className = className
	j.contentToList = make([]*HTML, 0)
	if j.parent == nil {
		j.moduleMap = make(map[string]*Attr, 10)
		j.defineMap = make(map[string](map[string]string))
		j.defineClassMap = make(map[string]string)
	}
	if node != nil {
		j.node = node
		j.innerContent = node.Child()
	} else {
		j.innerContent = nil
	}

	if className == "" {
		return errors.New("ui.go -> className is empty.")
	}
	j.relativePath = strings.Replace(className, ".", "/", -1)
	file := j.relativePath
	if root == "" {
		j.path = file
		j.htmlPath = JUSExist(file + ".ui")
		j.jsPath = JUSExist(file + ".es")
		j.cssPath = JUSExist(file + ".css")
	} else {
		if file[0] == '$' {
			j.path = j.CLASS_PATH + file[1:]
			j.htmlPath = JUSExist(j.path + ".ui")
			j.jsPath = JUSExist(j.path + ".es")
			j.cssPath = JUSExist(j.path + ".css")
			j.IsSysLib = true
		} else {
			j.path = root + "/" + file
			j.htmlPath = JUSExist(j.path + ".ui")
			j.jsPath = JUSExist(j.path + ".es")
			j.cssPath = JUSExist(j.path + ".css")
			if j.htmlPath == "" && j.jsPath == "" && j.cssPath == "" {
				j.path = j.CLASS_PATH + file
				j.htmlPath = JUSExist(j.path + ".ui")
				j.jsPath = JUSExist(j.path + ".es")
				j.cssPath = JUSExist(j.path + ".css")
				j.IsSysLib = true
			}
		}

	}
	j.dirPath = Substring(j.path, 0, LastIndex(j.path, "/"))

	if j.htmlPath != "" {
		j.html = &HTML{}
		t, err := GetCode(j.htmlPath)
		if err != nil {
			return err
		}
		_, err = j.html.ReadFromString(CodeFx(t, j.IsTest))
		if err != nil {
			return err
		}
	} else if j.jsPath != "" {
		j.scriptFile = true
		j.PushImportScript(&Attr{className, className}) //change by sunxy 2018-3-2
	} else {
		return errors.New("ui.go -> the file isn't exist.")

	}

	if domain == "" {
		j.domain = "\b"
	} else {
		j.domain = domain
	}

	return nil
}

///获取此类需要导入的头文件
func (j *UI) GetImportScript() map[string]*Attr {
	return j.GetRoot().scriptElement
}

func (j *UI) PushImportScript(value *Attr) {
	if j.GetRoot().scriptElement == nil {
		j.GetRoot().scriptElement = make(map[string]*Attr, 10)
	}
	if j.GetRoot().scriptElement[value.Value] == nil {
		j.GetRoot().scriptElement[value.Value] = value
		pos := Index(value.Value, "\002")
		if pos != -1 {
			if Index(value.Value, "/") != -1 || Index(value.Value, "\\") != -1 { //for example import vue from "lib/js/vue.min.js";
				tp := string(value.Value[pos+1:])
				if Index(tp, "index.res") == 0 {
					tp = j.SYSTEM_PATH + "/root/" + tp[10:]
				} else {
					tp = j.root + "/" + tp
				}
				if Exist(tp) {
					tp = filepath.Clean(tp)
					m5, _ := F2md5(tp)
					j.GetRoot().scriptElementBuffer = append(j.GetRoot().scriptElementBuffer, &ScriptElement{"I", value.Name, "U", m5 + value.Value}) //代表UMD规范的包导入
				}
				return
			} else { //for example: import dialog from "jus.Dialog";
				value.Value = Substring(value.Value, pos+1, -1)
			}
		}
		if Index(value.Value, "/") != -1 || Index(value.Value, "\\") != -1 { //for example: import /lib/js/jquery.min.js;
			tp := ""
			if Index(value.Value, "index.res") == 0 {
				tp = j.SYSTEM_PATH + "/root/" + value.Value[10:]
			} else {
				tp = j.root + "/" + value.Value
				//value.Value = Substring(tp, StringLen(j.root), -1)
			}
			if Exist(tp) {
				tp = filepath.Clean(tp)
				m5, _ := F2md5(tp)

				j.GetRoot().scriptElementBuffer = append(j.GetRoot().scriptElementBuffer, &ScriptElement{"I", value.Name, "P", m5 + value.Value})
			} else {
				j.GetRoot().scriptElementBuffer = append(j.GetRoot().scriptElementBuffer, &ScriptElement{"O", value.Value, "O", j.className + ": " + value.Value + " isn't Exist."})
			}
			return
		}
		pub := false
		if j.SERVER != nil {
			pub = j.SERVER.IsPublic
		}
		ft := &UI{SERVER: j.SERVER, IsPublic: pub, SYSTEM_PATH: j.SYSTEM_PATH, CLASS_PATH: j.CLASS_PATH}
		if err := ft.CreateFromParent(j.root, "", nil, strings.TrimSpace(value.Value), j); err == nil { //for example: import jus.Dialog;
			ft.IsImport = value.Value
			if ft.IsScript() {
				scriptObj := &Script{}
				scriptObj.CreateFrom(j, j.root, j.domain, j.paramValue, j.extendsScriptBuffer, strings.TrimSpace(value.Value))
				tpr, _ := ft.GetInitString()
				j.GetRoot().scriptElementBuffer = append(j.GetRoot().scriptElementBuffer, &ScriptElement{"I", value.Value, "S", scriptObj.ReadFromString(tpr)}) //j.GetRoot().scriptElementBuffer = append(j.GetRoot().scriptElementBuffer, "\t_MODULE_CONTENT_LIST_[\f]['"+strings.TrimSpace(value.Name)+"'] = "+scriptObj.ReadFromString(j.scanMedia(tpr))+";\r\n")
			} else {
				j.GetRoot().scriptElementBuffer = append(j.GetRoot().scriptElementBuffer, &ScriptElement{"I", value.Value, "H", ft.ToFormatString()})
			}
		} else {
			j.GetRoot().scriptElementBuffer = append(j.GetRoot().scriptElementBuffer, &ScriptElement{"O", value.Value, "", err.Error()})
		}
	}
}

func (j *UI) PushSysLibDirs(className, path string) {
	if j.SysLibDirs == nil {
		j.SysLibDirs = make(map[string]string)
	}
	j.SysLibDirs[className] = path
}

/**
 * 加入指令语言
 */
func (j *UI) PushCommandScript(value *Attr) {
	if j.CommandCode == nil {
		j.CommandCode = make([]*Attr, 0)
	}
	j.CommandCode = append(j.CommandCode, value)
}

/**
 * 获取初始化导入的数据，html，js文件
 */
func (j *UI) GetInitString() (string, bool) {
	if j.htmlPath != "" {
		t, err := GetCode(j.htmlPath)
		if err != nil {
			return "", false
		}
		return CodeFx(t, j.IsTest), true
	} else if j.jsPath != "" {
		t, err := GetCode(j.jsPath)
		if err != nil {
			return "", false
		}
		return CodeFx(t, j.IsTest), true
	} else {
		return "", false

	}
}

/**
 *
 * @param root
 * @param domain
 * @param value
 * @param className
 * @param parent
 * @throws IOException
 */
func (j *UI) CreateFromParent(root string, domain string, node *HTML, className string, parent *UI) error {
	j.parent = parent
	if parent != nil {
		j.scanLevel = parent.scanLevel + 1
	}
	return j.CreateFrom(root, domain, node, className)

}

func (j *UI) SetConstructor(value *Attr) *UI {
	j.paramValue = value
	return j
}

func (j *UI) SetValue(value string) {
	j.innerValue = value
}

func (j *UI) GetDomain() string {
	return j.domain
}

func (j *UI) GetClassName() string {
	return j.className
}

func (j *UI) GetStaticMap() map[string][]*Attr {
	if j.parent != nil {
		return j.parent.GetStaticMap()
	}

	if j.staticScript == nil {
		j.staticScript = make(map[string][]*Attr, 10)
	}
	return j.staticScript
}

func (j *UI) GetConstructorCode() *[]*Attr {
	if j.parent != nil {
		return j.parent.GetConstructorCode()
	}
	return &(j.componentCode)
}

/**
 * 获取参数集合
 */
func (j *UI) GetComponentParamSet() *[]string {
	if j.parent != nil {
		return j.parent.GetComponentParamSet()
	}
	if j.componentParams == nil {
		j.componentParams = make([]string, 0, 10)
	}
	return &j.componentParams

}

func (j *UI) GetStaticCodeMap() map[string][]*Attr {
	if j.parent != nil {
		return j.parent.GetStaticCodeMap()
	}
	if j.staticCode == nil {
		j.staticCode = make(map[string][]*Attr, 10)
	}
	return j.staticCode
}

func (j *UI) GetStyleCodeMap() map[string]string {
	//if j.parent != nil {
	//	return j.parent.GetStyleCodeMap()
	//}
	if j.styleCode == nil {
		j.styleCode = make(map[string]string, 10)
	}
	return j.styleCode
}

/**
 * 获取顶级
 */
func (j *UI) GetRoot() *UI {
	if j.parent != nil {
		return j.parent.GetRoot()
	}
	return j
}

/**
 * 添加静态函数表达式
 * @param className
 * @param func
 */
func (j *UI) AddStaticScript(className string, funcName string, value string) {
	j.staticScript = j.GetStaticMap()
	fun := j.staticScript[className]
	if fun == nil {
		fun = make([]*Attr, 0, 100)
	}
	attr := &Attr{Name: funcName, Value: value}
	for _, a := range fun {
		if a.Name == funcName {
			return
		}
	}
	fun = append(fun, attr)
	j.staticScript[className] = fun
}

/**
 * 添加静态代码
 * @param className
 * @param func
 */
func (j *UI) AddStaticCode(className string, funcName string, value string) {
	j.staticCode = j.GetStaticCodeMap()
	fun := j.staticCode[className]
	if fun == nil {
		fun = make([]*Attr, 0, 1000)
	}
	attr := &Attr{Name: funcName, Value: value}
	for _, a := range fun {
		if a.Name == funcName {
			return
		}
	}
	fun = append(fun, attr)
	j.staticCode[className] = fun
}

func (j *UI) AddStyleCode(className string, value string) {
	if j.styleCode == nil {
		j.styleCode = make(map[string]string, 10)
	}
	fun := j.styleCode[className]
	if fun == "" {
		j.styleCode[className] = value
	}

}

func (j *UI) overHTML(node []*HTML) {
	child := &HTML{}
	child.InsertList(node, 0)
	overList := child.Filter("@override")
	var lst []*HTML = nil
	var t *HTML = nil
	var p *HTML = nil

	//Override
	if len(overList) != 0 {
		for k := 0; k < len(overList); k++ {
			lst = overList[k].Child()
			for i := 0; i < len(lst); i++ {
				p = lst[i]
				if "script" == p.TagName() && p.GetAttr("type") == "" {
					j.extendsScriptBuffer += ListToHTMLString(p.Child())
					continue
				}
				t = j.html.GetElementById(Replace(p.GetAttr("id"), "\b", j.domain))
				if t != nil {
					t.ReplaceWith(p)
				}
			}
		}
	}

	child.RemoveChildByTagName("@override")

	//----开始替换----
	pList := j.html.GetElementsByTagName("@content")
	if len(pList) > 0 {
		var t []*HTML = nil
		for _, h := range pList {
			if h.GetAttr("to") != "" {
				j.contentTo = h.GetAttr("to")
			}
			if child.IsEmpty() { //采用默认自带方案
				t = h.Child()
			} else {
				t = child.Child()
			}
			for _, v := range t {
				if v.tag != "" && h.GetAttr("to") != "" {
					v.SetAttr("____CONTENT____", h.GetAttr("to")) //如果@content有to表示把这条数据添加到变量里
					j.contentToList = append(j.contentToList, v)
				}
			}
			h.ReplaceWithFromList(t)
		}
	}

	pList = j.html.GetElementsByTagName("@value")
	if len(pList) > 0 {
		if j.innerValue == "" {
			if child.IsEmpty() { //文字
				j.clearMark(pList[0].Child())
				j.innerValue = ListToHTMLString(pList[0].Child())

			} else {
				j.clearMark(child.Child())
				j.innerValue = ListToHTMLString(child.Child())
			}
		}

		for _, h := range pList {
			h.Remove()
		}
	}

	pList = j.html.GetElementsByTagName("@component")
	cp := ""
	if len(pList) > 0 {
		if child.IsEmpty() { //文字
			j.clearMark(pList[0].Child())
			cp = ListToHTMLString(pList[0].Child())

		} else {
			j.clearMark(child.Child())
			cp = ListToHTMLString(child.Child())
		}
		md5Ctx := md5.New()
		md5Ctx.Write([]byte(cp))
		cipherStr := md5Ctx.Sum(nil)
		bs := hex.EncodeToString(cipherStr)
		ft := &UI{SERVER: j.SERVER, SYSTEM_PATH: j.SYSTEM_PATH, CLASS_PATH: j.CLASS_PATH}
		ft.CreateFromString(j.root, "", nil, cp, bs, nil, j)
		j.GetRoot().scriptElementBuffer = append(j.GetRoot().scriptElementBuffer, &ScriptElement{"I", bs, "H", ft.ToFormatString()})
		j.innerModule = bs
		for _, h := range pList {
			h.Remove()
		}
	}

} //overHTML

func (j *UI) clearMark(child []*HTML) []*HTML {
	var p *HTML = nil
	for i := 0; i < len(child); i++ {
		p = child[i]
		p.RemoveAttr("domain")
		p.RemoveAttr("____format____")
		if p.TagName() == "@uncare" {
			p.ReplaceWithFromList(p.Child())
			continue
		}
		j.clearMark(p.Child())
	}
	return child
}

func (j *UI) scanHTML(child []*HTML) {
	tagName := ""
	attrName := ""
	attrValue := ""
	arr := make([]string, 0)
	var tHTML *HTML = nil
	for _, p := range child {
		tagName = p.TagName()
		if "@uncare" == p.TagName() {
			continue
		}
		if "module" == p.TagName() {
			tagName = "core.module"
		}

		if "script" == p.TagName() && p.GetAttr("type") == "" {
			continue
		}

		if p.GetAttr("isroot") != "" {
			p.SetAttr("id", j.domain)
		}

		//解读指令
		c := p.GetAttrCmd()
		for _, v := range c {
			attrValue = p.GetAttr(v)
			if Index(v, ".") == -1 {
				p.RemoveAttr(v)
				v = v[1:]
				v = strings.ToLower(v)
				attrName = j.pkgMap[v]
				if attrName != "" {
					p.SetAttrName(v, attrName)
					j.PushCommandScript(&Attr{attrName, "-" + v + "\002" + attrName + "\002" + p.GetAttr("id") + "\002" + attrValue})
					j.PushImportScript(&Attr{attrName, attrName})
				}
			} else {
				v = v[1:]
				attrName = strings.ToLower(v)
				j.PushCommandScript(&Attr{attrName, "-" + v + "\002" + attrName + "\002" + p.GetAttr("id") + "\002" + attrValue})
				j.PushImportScript(&Attr{attrName, attrName})
			}
		}

		if Index(tagName, ".") != -1 {
			arr = strings.Split(tagName, ":")
			if len(arr) > 1 {
				tagName = arr[1]
			}
			pub := false
			if j.SERVER != nil {
				pub = j.SERVER.IsPublic
			}
			var tFunc *UI = &UI{SERVER: j.SERVER, IsPublic: pub, SYSTEM_PATH: j.SYSTEM_PATH, CLASS_PATH: j.CLASS_PATH, IsImport: j.IsImport, Debug: j.Debug}
			var err error
			if Index(tagName, ".$") == 0 {
				err = tFunc.CreateFromString(j.root, p.GetAttr("id"), p, j.GetRoot().defineClassMap[tagName], tagName, j, nil)
			} else {
				err = tFunc.CreateFromParent(j.root, p.GetAttr("id"), p, tagName, j)
			}
			if err == nil {
				if tFunc.IsScript() {
					tFunc.SetConstructor(&Attr{tagName, p.GetConstructerParameter()}).setExtend(p.GetAttr("id") == j.domain)
					// if p.GetConstructerCode() != "" {
					// 	j.idMap[p.GetAttr("src_id")] = &HTMLObject{Name: p.GetAttr("id"), HTMLObjectType: 1}
					// }
					th := tFunc.ReadHTML()
					j.AddRun(&RunElem{"L", j.domain, th.ToString()})
					tHTML = &HTML{}
				} else {
					if p.GetConstructerParameter() != "" {
						j.componentParam = append(j.componentParam, &Attr{p.GetAttr("id"), p.GetConstructerParameter()})
						j.AddRun(&RunElem{Type: "P", Name: p.GetAttr("id"), Value: "[" + j.componentInitParam(&Attr{Name: p.GetAttr("id"), Value: p.GetConstructerParameter()}) + "]"})
					}
					if p.GetConstructerCode() != "" {
						j.componentCode = append(j.componentCode, &Attr{p.GetAttr("id"), p.GetConstructerCode()})
					}
					tFunc.SetConstructor(&Attr{tagName, p.GetConstructerParameter()}).setExtend(p.GetAttr("id") == j.domain)
					tHTML = tFunc.ReadHTML()
					//clsTmp := tHTML.GetAttr("class")
					tHTML.CopyFrom(p)
					//tHTML.SetAttr("class", clsTmp+" &"+tHTML.GetAttr("class"))
					if len(arr) > 1 {
						tHTML.SetTagName(arr[0])
					}
					if p.GetAttr("____CONTENT____") != "" {
						j.AddRun(&RunElem{"T", p.GetAttr("id"), ""}) //统一L，但是是系统初始化最后加上的。
					}
				}
			} else {
				tHTML, _ = (&HTML{}).ReadFromString("<div style='font-size:14px;font-weight:bold;background-color: #E91E63;color: #fefefe;padding: 5px;border-radius: 5px;display: inline;'>" + err.Error() + "</div>")
			}
			if p == j.html {
				j.html = tHTML
			} else {
				p.ReplaceWith(tHTML)
			}
			continue
		}
		j.scanHTML(p.Child())
	}
}

func (j *UI) componentInitParam(value *Attr) string {
	s := j.componentInitCode(value)
	script := &HTMLScript{}
	script.CreateFrom(j, j.root, j.domain, j.paramValue, j.innerValue, j.extendsScriptBuffer)
	str := script.FormatString(s)
	return str
}

/**
 * 设置扩展
 * @param flag
 * @return
 */
func (j *UI) setExtend(flag bool) *UI {
	j.extendFlag = flag
	return j
}

/**
 * 对所有控件的ID进行记录
 * @param html
 */
func (j *UI) componentId(child []*HTML) {
	for _, p := range child {
		if p.GetAttr("domain") != "" && p.GetAttr("domain") == j.domain {
			if p.GetAttr("class_id") != "" || Index(p.TagName(), ".") != -1 {
				j.idMap[p.GetAttr("src_id")] = &HTMLObject{Name: p.GetAttr("id"), HTMLObjectType: 1}
			} else {
				j.idMap[p.GetAttr("src_id")] = &HTMLObject{Name: p.GetAttr("id"), HTMLObjectType: 0}
			}
		}

		j.componentId(p.Child())
	}

}

func (j *UI) cssComponent(child []*HTML) {
	var tagName string
	c := 0
	for _, p := range child {
		if p.GetAttr("domain") != "" && p.GetAttr("domain") == j.domain {
			tagName = p.GetAttr("class_id")
			if tagName != "" {
				c = LastIndex(tagName, ".")
				if c != -1 {
					tagName = Substring(tagName, c+1, len(tagName))
				}
				tmp := ""
				if j.cssTag[strings.ToLower(tagName)] != "" {
					tmp = " " + j.cssTag[strings.ToLower(tagName)]
				}

				p.SetAttr("class", p.GetAttr("class")+" "+tmp)

			}
		}

		j.cssComponent(p.Child())
	}

}

func (j *UI) styleFormat() string {
	j.style.AddDomain("." + j.domain)
	j.style.ReplaceSelecter("body", "."+j.domain)
	j.cssTag = j.style.GetComponentClass(false)
	j.cssComponent([]*HTML{j.html})
	return ScriptInitD(j.style.ToString(1), j.domain)

}

/**
 * 公共css属性，也可以认为某个控件的全局css样式
 */
func (j *UI) cssFormat() string {
	j.css.AddDomain(".-" + Replace(j.className, ".", "-"))
	j.css.ReplaceSelecter("body", ".-"+Replace(j.className, ".", "-"))
	j.cssTag = j.css.GetComponentClass(true)
	j.cssComponent([]*HTML{j.html})
	return ScriptInitD(j.css.ToString(0), j.domain)
}

/**
 * 加载网页配置信息
 * 例如判断网页是否可以发布，发布的方式和模板是什么
 */
func (j *UI) loadSetting() {
	sets := j.html.GetElementsByTagName("@pub")
	for _, v := range sets {
		j.pub = v.GetAttr("value")
		if j.pub == "" {
			j.pub = "default"
		}
	}
	j.html.RemoveChildByTagName("@pub")
}

func (j *UI) importHTML() {
	sets := j.html.GetElementsByTagName("@import")
	sets = append(sets, &HTML{tagData: map[string]string{"value": ""}})
	sets = append(sets, &HTML{tagData: map[string]string{"value": ".$" + j.className}})

	p := 0
	value := ""
	path := ""
	fileName := ""
	key := ""
	cls := "" //文件类型

	for i := 0; i < len(sets); i++ {
		value = sets[i].GetAttr("value")
		if Index(value, ".$") == 0 { //说明是获取模块内部自定义类的类
			m := j.GetRoot().defineMap[value]
			if m != nil {
				for k, v := range m {
					j.pkgMap[k] = v
				}
			}
		} else {
			if Index(value, ".") == 0 { //说明是获取自己本地路径
				value = Substring(j.dirPath, StringLen(j.root), -1) + value
				value = filepath.Clean(value)
				value = Replace(value, "\\", ".")
				value = Replace(value, "/", ".")
			}
			value = strings.TrimLeft(value, ".")
			value = strings.Replace(value, ";", "", -1)
			value = strings.Replace(value, " ", "", -1)
			p = LastIndex(value, ".")
			if p != -1 {
				path = Substring(value, 0, p)
				if CharAt(value, p+1) != "*" {
					fileName = Substring(value, p+1, -1) + ".ui"
				}
			}

			fl := j.CLASS_PATH + "/" + strings.Replace(path, ".", "/", -1)
			var lst []os.FileInfo
			var err error
			if Exist(fl) {
				lst, err = ioutil.ReadDir(fl)
				if err == nil {
					for _, f := range lst {
						if !f.IsDir() && (fileName == "" || fileName == f.Name()) {
							cls = filepath.Ext(f.Name())
							if cls == ".ui" || cls == ".es" {
								key = Substring(f.Name(), 0, LastIndex(f.Name(), "."))
								j.pkgMap[strings.ToLower(key)] = path + "." + key
							}
						}
					}
				} else {
					fmt.Println(err)
				}
			}

			fl = j.root + "/" + strings.Replace(path, ".", "/", -1)
			if Exist(fl) {
				lst, err = ioutil.ReadDir(fl)
				if err == nil {
					for _, f := range lst {
						if !f.IsDir() && (fileName == "" || fileName == f.Name()) {
							cls = filepath.Ext(f.Name())
							if cls == ".ui" || cls == ".es" {
								key = Substring(f.Name(), 0, LastIndex(f.Name(), "."))
								j.pkgMap[strings.ToLower(key)] = path + "." + key
							}
						}
					}
				} else {
					fmt.Println(err)
				}
			}
		}

		path = ""
		fileName = ""
	}
	j.html.RemoveChildByTagName("@import")
}

//自定义HTML模块
func (j *UI) defineHTML() {
	sets := j.html.GetElementsByTagName("@define")
	var name string
	for _, v := range sets {
		name = ".$" + j.className + "." + v.GetAttr("name")
		if j.GetRoot().defineMap[".$"+j.className] == nil {
			j.GetRoot().defineMap[".$"+j.className] = make(map[string]string)
		}
		j.GetRoot().defineMap[".$"+j.className][v.GetAttr("name")] = name
		j.GetRoot().defineClassMap[name] = "<@import value='.$" + j.className + "'/>\r\n" + v.GetInnerHTML()
	}
	j.html.RemoveChildByTagName("@define")
}

/**
 * 获取注释信息
 */
func (j *UI) rootHTML() {
	child := j.html.Filter("!")
	for _, v := range child {
		v.Remove()
	}
	child = j.html.Filter("head")
	for _, v := range child {
		j.headBuffer.Write(ListToHTMLStringBytes(v.Child()))
		v.Remove()
	}
	child = j.html.Filter("style")
	for _, v := range child {
		j.styleBuffer.Write(ListToHTMLStringBytes(v.Child()))
		v.Remove()
	}
	child = j.html.Filter("css")
	for _, v := range child {
		j.cssBuffer.Write(ListToHTMLStringBytes(v.Child()))
		v.Remove()
	}

}

/**
 *
 */
func (j *UI) packageHTML(child []*HTML) {
	tagName := ""
	extName := ""
	var arr []string
	for _, p := range child {
		tagName = strings.ToLower(p.TagName())
		arr = strings.Split(tagName, ":")
		if len(arr) > 1 {
			tagName = arr[0]
			extName = arr[1]
		}

		//替换Module TagName 变为真是的tagName

		if s := j.pkgMap[tagName]; s != "" {
			tagName = s
			p.SetTagName(s)
		} else if s := j.pkgMap[p.TagName()]; s != "" {
			tagName = s
			p.SetTagName(s)
		}

		if extName != "" && j.pkgMap[extName] != "" {
			extName = j.pkgMap[extName]
			p.SetTagName(tagName + ":" + extName)
		}
		extName = ""
		j.packageHTML(p.Child())

	}
}

/**
 * 对Module内部的innerHTML 做提前本域名下id绑定
 * @param html
 */
func (j *UI) domainHTML(child []*HTML) {
	tagName := ""
	for _, p := range child {
		if p.tagType == -1 {
			tagName = "\b"
		} else {
			tagName = strings.ToLower(p.TagName())
		}

		if "@override" == tagName {
			continue
		}

		if "@uncare" == tagName {
			continue
		}

		if ("script" == tagName || "~script" == tagName) && p.GetAttr("type") == "" {
			j.scriptBuffer.Write(ListToHTMLStringBytes(p.Child()))
			p.Remove()
			continue
		}

		if "style" == tagName {
			j.styleBuffer.Write(ListToHTMLStringBytes(p.Child()))
			p.Remove()
			continue
		}

		if "css" == tagName {
			j.cssBuffer.Write(ListToHTMLStringBytes(p.Child()))
			p.Remove()
			continue
		}

		j.domainHTML(p.Child())
	}
}

func (j *UI) GetPackageMap() map[string]string {
	return j.pkgMap
}

/**
 * 获取HTML定义的ID内容
 * @param name
 * @return
 */
func (j *UI) GetDefine(name string) *HTMLObject {
	if name[0] == '$' && len(name) > 1 {
		name = string([]rune(name)[1:])
		if j.idMap[name] == nil {
			return &HTMLObject{Name: j.domain + name, HTMLObjectType: 0}
		}
	}
	return j.idMap[name]
}

/**
 * 替换所有@lib:变量
 * @param value
 * @return
 */
func (j *UI) scanMedia(value string) string {
	data := []rune(value)
	sb := bytes.NewBufferString("")
	tmp := make([]rune, 0, 1000)
	position := 0
	k := 0
	var ch rune
	keys := [...]rune{'@', 'l', 'i', 'b', '('}
	xl := 0
	for position < len(data) {
		ch = data[position]
		position++
		if ch == keys[k] {
			k++
			if k == len(keys) {
				k = 0
				xl = 1
				for position < len(data) {
					ch = data[position]
					position++
					if ch == '(' {
						xl++
					} else if ch == ')' {
						xl--
					}
					if xl == 0 {
						break
					}
					if xl == 1 && ch == '(' {
						continue
					}

					tmp = append(tmp, ch)
				}
				file := string(tmp)
				file = Replace(Replace(file, "'", ""), "\"", "")
				path, err := filepath.Abs(j.htmlPath)
				if err != nil {
					fmt.Println(path, err)
				}

				f := Substring(path, 0, LastIndex(path, ".")) + ".lib/" + file

				if Exist(f) {
					data, _ := GetBytes(f)
					switch filepath.Ext(file) {
					case ".svg":
						sb.WriteString("data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString(data))
					case ".jpg":
						fallthrough
					case ".jpeg":
						sb.WriteString("data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(data))
					default:
						sb.WriteString("data:image/png;base64," + base64.StdEncoding.EncodeToString(data))
					}

				}

				tmp = tmp[0:0]
			}
			continue
		} else {
			if k > 0 {
				for i := 0; i < k; i++ {
					sb.WriteRune(keys[i])
				}
			}
			sb.WriteRune(ch)
			k = 0
		}
	}
	sb.WriteString(" ")
	return sb.String()
}

/**
 * 获取ClassPath 位置路径
 */
func (j *UI) GetRootRealPath() string {
	return ""
}

/**
 * 判断事否为纯Script文件
 * @return
 */
func (j *UI) IsScript() bool {
	return j.scriptFile
}

/**
 * 初始化Attr里的@use
 * @param html
 */
func (j *UI) useHTML(html *HTML) {
	arr := make([]*HTML, 0)
	useFunc(html, &arr)
	for _, p := range arr {
		str := p.GetAttr("@use")
		if str != "" {
			path := j.root + str
			code := ""
			if Exist(path) { //先查下用户.settings目录下有没有
				code, _ = GetCode(path)
			} else if Exist(j.SYSTEM_PATH + "/core/use/" + str) { //如果没有查下SYS lib下有没有
				path = j.SYSTEM_PATH + "/core/use/" + str
				code, _ = GetCode(path)
			}
			if code != "" {
				vm := goja.New()
				console := &Console{Name: str}
				vm.Set("__goja_log__", console.Log)
				vm.RunString(`var console = {};console.log = __goja_log__;`)
				vm.Set("__ARG__", p)
				//首先先加载这个地址执行的文件
				vm.RunString(code) //执行str内部的函数
				//将这个文件里的内容，交给js执行。
				_, re := vm.RunString("main(__ARG__)")
				if re != nil {
					console.Log(re.Error())
				}

			}
		}
	}

}

func useFunc(html *HTML, arr *[]*HTML) {
	for _, p := range html.Child() {
		if p.GetAttr("@use") != "" {
			*arr = append(*arr, p)
		}
		useFunc(p, arr)
	}
}

/**
 * 初始化Attr里的@this
 * @param html
 */
func (j *UI) initObj(html *HTML) {
	for _, p := range html.Child() {
		if p.TagName() != "" && !("\b" == p.TagName()) {
			if p.GetAttr("domain") == "" {
				p.SetAttr("domain", j.domain)
			}
		}
		if p.GetAttr("id") == "" {
			p.SetAttr("id", p.GetAttr("domain")+j.getName())
		} else {
			if p.GetAttr("id")[0] == '$' {
				p.SetAttr("src_id", p.GetAttr("id")[1:])
				p.SetAttr("id", p.GetAttr("domain")+p.GetAttr("id")[1:])
			} else {
				p.SetAttr("src_id", p.GetAttr("id"))
				p.SetAttr("id", p.GetAttr("domain")+p.GetAttr("id"))
			}
		}
		if p.TagName() == "@uncare" {
			continue
		}
		if p.TagName() == "script" {
			if p.GetAttr("type") == "" {
				j.scriptBuffer.Write(ListToHTMLStringBytes(p.Child()))
				p.Remove()
			}
			continue
		}
		for _, attr := range p.Attrs() {
			if "id" == strings.ToLower(attr.Name) {
				continue
			}
			if Index(p.GetAttr(attr.Name), "@this") != -1 {
				p.SetAttr(attr.Name, ScriptInitD(strings.Replace(p.GetAttr(attr.Name), "@this", j.domain, -1), j.domain))
				j.IsPublic = true
			}
			p.SetAttr(attr.Name, ScriptInitD(strings.Replace(p.GetAttr(attr.Name), "@lib", IfStr(j.IsSysLib, "index.src/", "./")+j.relativePath+".lib", -1), j.domain))
		}
		j.initObj(p)
	}
}

func (j *UI) testHTML() *HTML {
	return j.html
}

func (j *UI) ReadHTML() *HTML {
	if j.scriptFile {
		tHTML := &HTML{}
		//
		if j.parent == nil {
			sb := bytes.NewBufferString("<script>")
			for _, v := range j.scriptElementBuffer {
				j.ToFormatLine(v.Cls, v.ModuleName, v.Header+v.Value, sb)
			}
			sb.WriteString("</script>")
			tHTML.ReadFromString(sb.String())

			return tHTML
		}

		tps := bytes.NewBufferString("")
		if j.paramValue != nil && j.paramValue.Value != "" {
			tps.WriteString(j.paramValue.Value)
			tps.WriteRune(',')
		}
		if tps.Len() > 0 {
			tps.Truncate(tps.Len() - 1)
		}
		tst := bytes.NewBufferString("")

		//解析节点属性值
		if j.node != nil {
			for _, va := range j.node.Attrs() {
				tst.WriteString("$$.")
				tst.WriteString(va.Name)
				tst.WriteString("=\"")
				if Index(va.Value, "\b") == -1 {
					tst.WriteString(Escape(va.Value))
				} else {
					tst.WriteString(va.Value)
				}
				tst.WriteString("\";\r\n")
			}
			tst.WriteString("__OBJECT__.")
			tst.WriteString(j.node.GetAttr("id"))
			tst.WriteString("= $$;\r\n")
		}

		//解析内部设置值
		for _, v := range j.innerContent {
			for _, v2 := range v.Child() {
				if v2.IsText() {
					if strings.TrimSpace(v2.Text()) != "" {
						tst.WriteString(j.domain)
						tst.WriteRune('.')
						tst.WriteString(v.TagName())

						tst.WriteString("=\"")
						tst.WriteString(Escape(v2.Text()))
						tst.WriteString("\";\r\n")
					}
				} else {
					if v2.GetAttr("id") == "" {
						v2.SetAttr("id", v2.GetAttr("domain")+j.getName())
					} else {

						if v2.GetAttr("id")[0] == '$' {
							v2.SetAttr("src_id", v2.GetAttr("id")[1:])
							v2.SetAttr("id", v2.GetAttr("domain")+v2.GetAttr("id")[1:])
						} else {
							v2.SetAttr("src_id", v2.GetAttr("id"))
							v2.SetAttr("id", v2.GetAttr("domain")+v2.GetAttr("id"))
						}
					}
					pub := false
					if j.SERVER != nil {
						pub = j.SERVER.IsPublic
					}
					var tFunc *UI = &UI{SERVER: j.SERVER, IsPublic: pub, SYSTEM_PATH: j.SYSTEM_PATH, CLASS_PATH: j.CLASS_PATH}
					j.idMap[v2.GetAttr("src_id")] = &HTMLObject{Name: v2.GetAttr("id"), HTMLObjectType: 1}
					if err := tFunc.CreateFromParent(j.root, v2.GetAttr("id"), v2, v2.TagName(), j); err == nil {
						if tFunc.IsScript() {
							tFunc.SetConstructor(&Attr{v2.TagName(), v2.GetConstructerParameter()}).setExtend(v2.GetAttr("id") == j.domain)
							if v2.GetConstructerCode() != "" {
								var arr *[]*Attr = j.GetConstructorCode()
								*arr = append(*arr, &Attr{v2.GetAttr("id"), v2.GetConstructerCode()})
							}
							tst2 := bytes.NewBufferString("var " + v2.GetAttr("id") + "=" + tFunc.ReadHTML().Text())
							tst2.WriteString("()\r\n")
							tst2.Write(tst.Bytes())
							tst = tst2
						} else {
							//如果是HTML文件
							htps := ""
							if v2.GetConstructerParameter() != "" {
								htps = "," + v2.GetConstructerParameter()
							}
							if v2.GetConstructerCode() != "" {
								var arr *[]*Attr = j.GetConstructorCode()
								*arr = append(*arr, &Attr{v2.GetAttr("id"), v2.GetConstructerCode()})
							}

							j.PushImportScript(&Attr{v2.TagName(), v2.TagName()})
							tst2 := bytes.NewBufferString(tFunc.ReadHTML().Text())
							tst2.WriteString("var ")
							tst2.WriteString(v2.GetAttr("id"))
							tst2.WriteString(" = getModule(\"")
							tst2.WriteString(v2.TagName())
							tst2.WriteString("\",__APPDOMAIN__")
							tst2.WriteString(htps)
							tst2.WriteString(");\r\n")
							tst2.Write(tst.Bytes())
							tst = tst2
						}
					} else {
						tst.WriteString("__ERROR__('" + err.Error() + "')")
					}
					//tst.WriteString(Replace(j.domain, "\b", "____"))
					tst.WriteString("$$")
					tst.WriteRune('.')
					tst.WriteString(v.TagName())
					tst.WriteString("=")
					tst.WriteString(v2.GetAttr("id"))
					tst.WriteString(";\r\n")
				}
			}
		}
		tHTML.ReadFromString("(function(){ var $$ = getModule(\"" + j.className + "\",__APPDOMAIN__)(" + tps.String() + ");\r\n" + tst.String() + "return $$;})")
		return tHTML
	}
	j.loadSetting()
	j.useHTML(j.html)
	j.rootHTML()
	j.defineHTML()
	j.importHTML()
	j.initObj(j.html)
	htmls := j.html.GetUnTextChild()
	if len(htmls) == 1 {
		j.html = htmls[0]
	} else {
		j.html = &HTML{}
		j.html.ReadFromString("<div></div>")
		j.html = j.html.At(0)
		j.html.InsertList(htmls, 0)
	}
	//加载外部CSS
	if j.cssPath != "" {
		css := &HTML{}
		tpr, _ := GetCode(j.cssPath)
		css.ReadFromString("<style>" + CodeFx(tpr, j.IsTest) + "</style>")
		j.html.Append(css)
	}

	j.overHTML(j.innerContent)
	j.packageHTML([]*HTML{j.html})
	j.componentId([]*HTML{j.html})
	j.domainHTML([]*HTML{j.html})
	if j.styleBuffer.Len() > 0 {
		j.style = &CSS{Root: &Attr{"#" + j.html.GetAttr("src_id"), j.html.TagName()}, Class: j.html.GetAttr("class"), jus: j, CurrentPath: IfStr(j.IsSysLib, "index.src/", "./") + j.relativePath + ".lib"}
		j.style.ReadFromString(j.scanMedia(j.styleBuffer.String()))
	}
	j.idMap[j.html.GetAttr("src_id")] = &HTMLObject{Name: j.domain, HTMLObjectType: -1} //代表容器节点
	j.html.SetAttr("id", j.domain)
	j.html.SetAttr("isroot", "true")
	headCss := ""
	headCss = Replace(j.className, ".", "-")
	headCss = "-" + headCss + " " + j.domain
	j.scanHTML([]*HTML{j.html})

	if j.contentTo != "" {
		j.scriptBuffer.WriteString("if(_MODULE_INNER_[__DOMAIN__]){____." + j.contentTo + "=_MODULE_INNER_[__DOMAIN__];}")
	}
	j.html.SetAttr("class", headCss+" "+j.html.GetAttr("class"))
	// if Index(j.html.GetAttr("class"), j.domain) == -1 {
	// 	fmt.Println(">>>")
	// 	j.html.SetAttr("class", IfStr(j.html.GetAttr("class") != "", j.html.GetAttr("class")+" ", "")+j.domain)
	// }

	if j.style != nil {
		style := &HTML{}
		style.ReadFromString("<style>" + j.styleFormat() + "</style>")
		j.html.Insert(style, 0)
	}

	//开始组装参数
	script := &HTMLScript{}
	script.CreateFrom(j, j.root, j.domain, j.paramValue, j.innerValue, j.extendsScriptBuffer)
	scriptCodeString := script.ReadFromString(j.scriptBuffer.String())

	if len(scriptCodeString) != 0 {
		node := &HTML{}
		node.AppendNode("script", scriptCodeString)
		j.html.Append(node)
	}

	if len(j.componentCode) > 0 {
		code := "function init(){"
		for _, v := range j.componentCode {
			code += j.componentInitCode(v)
		}
		code += "}"
		node := &HTML{}
		node.AppendNode("script", script.ReadFromString(code))
		j.html.Append(node)
	}

	if j.jsPath != "" {
		script = &HTMLScript{}
		script.CreateFrom(j, j.root, j.domain, j.paramValue, j.innerValue, j.extendsScriptBuffer)
		tpr, _ := GetCode(j.jsPath)
		scriptString := script.ReadFromString(CodeFx(tpr, j.IsTest)) //scriptString = script.ReadFromString(j.scanMedia(tpr))

		if len(scriptString) != 0 {
			node := &HTML{}
			node.AppendNode("script", script.ReadFromString(scriptString))
			j.html.Append(node)
		}
	}
	if j.cssBuffer.Len() > 0 {
		j.css = &CSS{Root: &Attr{"#" + j.html.GetAttr("src_id"), j.html.TagName()}, Class: j.html.GetAttr("class"), jus: j, CurrentPath: IfStr(j.IsSysLib, "index.src/", "./") + j.relativePath + ".lib"}
		j.css.ReadFromString(j.scanMedia(j.cssBuffer.String()))
		j.AddStyleCode(j.className, j.cssFormat())
	}
	sb := bytes.NewBufferString("<css>")
	j.styleCode = j.GetStyleCodeMap()
	for name, value := range j.styleCode {
		j.ToFormatLine("A", name, value, sb)
	}
	if sb.Len() > 5 {
		cssHTML := &HTML{}
		sb.WriteString("</css>")
		cssHTML.ReadFromString(sb.String())
		j.html.Insert(cssHTML, 0)
	}

	if len(j.CommandCode) > 0 {
		for _, v := range j.CommandCode {
			j.AddRun(&RunElem{Type: "C", Name: j.domain, Value: v.Value})
		}
	}
	//最终加入静态函数变量
	if j.parent == nil {
		st := bytes.NewBufferString("")
		sb := bytes.NewBufferString("")
		for _, v := range j.scriptElementBuffer {
			if v.Header == "P" || v.Header == "S" || v.Header == "U" || v.Header == "O" {
				j.ToFormatLine(v.Cls, v.ModuleName, v.Header+v.Value, sb)
				continue
			}
			sb.WriteString(v.Value)
		}
		j.staticCode = j.GetStaticCodeMap()
		for name, value := range j.staticCode {
			for _, attr := range value {
				st.WriteString("__POS_VALUE__")
				st.WriteString(attr.Value)
				j.ToFormatLine("S", name, attr.Name+" "+st.String(), sb)
				st.Reset()
			}
		}
		j.staticScript = j.GetStaticMap()
		for name, value := range j.staticScript {
			for _, attr := range value {
				st.WriteString("__POS_VALUE__" + attr.Value + ";\r\n")
				j.ToFormatLine("S", name, attr.Name+" "+st.String(), sb)
				st.Reset()
			}
		}

		if sb.Len() > 0 {
			node := &HTML{}
			node.AppendNode("script", sb.String())
			j.html.Insert(node, 0)
		}

		sb.Reset()
		//本类是否有参数
		if len(j.componentParams) > 0 {
			for _, v := range j.componentParams {
				sb.WriteString(v)
				//sb.WriteString("\r\n")
			}
		}
		if sb.Len() > 0 {
			node := &HTML{}
			node.AppendNode("script", sb.String())
			j.html.Append(node)
		}

		pList := j.html.GetElementsByTagName("@uncare")
		for _, v := range pList {
			v.ReplaceWithFromList(v.Child())
		}

		sb.Reset()
		sb = nil
	}

	j.html.SetAttr("class_id", j.className)
	return j.testHTML()
}

/**
 * 初始化空间默认代码
 */
func (j *UI) componentInitCode(value *Attr) string {
	ms := &MScript{}
	ms.ReadFromString(value.Value)
	sb := bytes.NewBufferString("")
	for _, v := range ms.GetJUIScriptData() {
		if v.Value == "this" && v.Domain == "class" {
			v.Value = value.Name
		}
		if v.Value == "@this" {
			v.Value = "__OBJECT__['" + j.domain + "']"
		}
		sb.WriteString(v.Value)
	}

	return sb.String()
}

///获取模块地图
///md5码表，区分是否有不同模块
func (j *UI) GetModuleMap() map[string]*Attr {
	if j.parent != nil {
		return j.parent.GetModuleMap()
	}
	return j.moduleMap
}

/**
 * 添加执行命令
 */
func (j *UI) AddRun(attr *RunElem) {
	if j.IsImport != "" {
		if j.parent != nil && j.parent.IsImport != "" && j.parent.IsImport == j.IsImport {
			j.parent.AddRun(attr)
			return
		}
	} else {
		if j.parent != nil {
			j.parent.AddRun(attr)
			return
		}
	}
	j.runList = append(j.runList, attr)
}

func (j *UI) getName() string {
	if j.parent != nil {
		return j.parent.getName()
	}
	j.count++
	return "a" + strconv.Itoa(j.count)
}

/**
 * 格式化输出内容
 * @param cls		信息类型
 * @param moduleName	模块名称
 * @param value		内容
 */
func (j *UI) ToFormatLine(cls string, moduleName string, value string, data *bytes.Buffer) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(value))
	cipherStr := md5Ctx.Sum(nil)
	bs := hex.EncodeToString(cipherStr)
	m := j.GetModuleMap()[bs]
	if cls != "L" && m != nil && m.Name == moduleName {
		return bs
	}
	j.GetModuleMap()[bs] = &Attr{Name: moduleName}
	data.WriteString(cls)
	data.WriteString(moduleName)
	data.WriteByte(' ')
	data.WriteString(bs)
	data.WriteByte(' ')
	data.WriteString(value)
	data.WriteString("\r\n\x01")
	return bs
}

/**
 * 格式化输出内容
 * @param cls		信息类型
 * @param moduleName	模块名称
 * @param value		内容
 */
func (j *UI) ToFormatRun(cls string, domain string, value string, data *bytes.Buffer) {
	data.WriteString(cls)
	data.WriteByte(' ')
	data.WriteString(j.className)
	data.WriteByte(' ')
	data.WriteString(domain)
	data.WriteByte(' ')
	data.WriteString(value)
	data.WriteString("\x01")
}

/**
 * JUServer使用的字节流
 * cls 模块生成类型，有可预览和不可预览的
 */
func (j *UI) ToFormatBytes() []byte {
	result := j.ReadHTML()
	stls := result.GetElementsByTagName("css") //获取公共css属性
	json := bytes.NewBufferString("\x01")
	if j.pub != "" && j.headBuffer.Len() > 0 {
		j.ToFormatLine("T", j.className, j.headBuffer.String(), json)
	}
	for _, v := range stls {
		json.WriteString(ListToHTMLString(v.Child()))
		v.Remove()
	}
	stls = result.GetElementsByTagName("style") //获取公共css属性
	for _, v := range stls {
		j.ToFormatLine("B", j.className, ListToHTMLString(v.Child()), json)
		v.Remove()
	}

	spts := result.GetElementsByTagName("script") //获取Script属性
	for _, v := range spts {
		if v.GetAttr("type") == "" {
			json.Write(ListToHTMLStringBytes(v.Child()))
			v.Remove()
		}

	}

	if j.IsScript() {

	} else {
		j.ToFormatLine("H", j.className, result.ToString(), json)
	}

	for _, v := range j.runList {
		j.ToFormatRun("R"+v.Type, v.Name, v.Value, json)
	}
	j.ToFormatRun("RO", j.className, "", json)
	return json.Bytes()
}

func (j *UI) ToFormatString() string {
	return string(j.ToFormatBytes())
}

/**
 * 对外暴露接口
 */
func (j *UI) GetCode(path string) (string, error) {
	return GetCode(path)
}

func (j *UI) ToFormatHTMLString(result string) string {
	vm := goja.New()
	console := &Console{}
	vm.Set("__goja_log__", console.Log)
	vm.RunString(`var console = {};console.log = __goja_log__;`)
	vm.Set("code", result)
	vm.Set("UI", j)
	v, e := GetCode(j.SYSTEM_PATH + "/core/pub/" + j.pub + "/decode.js")
	if e != nil {
		return "O" + j.className + " 00000000000000000000000000000000 " + e.Error()
	}
	var r goja.Value
	var err error
	_, err = vm.RunString(v)
	if err != nil {
		return err.Error()
	}
	r, err = vm.RunString("main()")
	if err != nil {
		return err.Error()
	}
	return r.String()
}

/**
 * 生成最终字节流
 */
func (j *UI) Bytes() []byte {
	bs := j.ToFormatBytes()
	if j.pub != "" {
		return []byte(j.ToFormatHTMLString(string(bs)))
	} else {
		return bs
	}
}
