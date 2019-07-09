// jus project main.go
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"

	//_ "image/jpeg"
	//_ "image/png"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	. "jus/str"

	. "jus/cn/airoot/util"

	. "jus"

	"golang.org/x/net/websocket"
)

var version string = "AIroot UI-SYSTEM 0.9.5beta"
var lang map[string]string

var zhCN = make(map[string]string, 0)
var enCH = make(map[string]string, 0)

func init() {
	zhCN["文件不存在"] = "\"%s\" 文件不存在"
	zhCN["添加成功"] = "[%s] 添加成功"
	zhCN["已经添加"] = "[%s] 已经添加"
	zhCN["项目已存在"] = "项目[%s]已存在"
	zhCN["建立项目"] = "建立项目[%s]"
	zhCN["项目挂载在"] = "项目挂载在[%s]服务上."
	zhCN["遍历结束"] = "----遍历结束----"
	zhCN["加载系统路径错误"] = "加载系统路径错误."
	zhCN["遍历运行"] = "%s. %s\t运行\t%s\t%s\t%s"
	zhCN["遍历停止"] = "%s. %s\t停止\t%s\t%s\t%s"
	zhCN["WS运行"] = "%s. %s\t%s\t%s\t%s\t%s"
	zhCN["WS停止"] = "%s. %s\t%s\t%s\t%s\t%s"
	zhCN["遍历未初始化"] = "<未初始化>"
	zhCN["不存在工程"] = "不存在[%s],设置工程目录失败."
	zhCN["不存在服务"] = "不存在[%s]服务，或者此服务已经被移除."
	zhCN["移除失败"] = "[%s] 移除失败."
	zhCN["移除成功"] = "[%s] 移除成功."
	zhCN["设置成功"] = "[%s] 设置成功."
	zhCN["属性移除成功"] = "%s 移除成功."
	zhCN["属性移除失败"] = "%s 移除失败,可能不存在此[%s]属性."
	zhCN["服务关闭失败"] = "%s 服务停止失败"
	zhCN["工程设置成功"] = "[%s] 的工程路径 [%s] 设置成功."
	zhCN["工程设置失败"] = "[%s] 的工程路径 [%s] 设置失败."
	zhCN["服务正在启动"] = "[%s] 服务在[%s]"
	zhCN["关闭服务"] = "%s 服务关闭[%s]"
	zhCN["发布完成"] = "----发布完成----"
	zhCN["添加WEB用户成功"] = "添加WEB用户成功."
	zhCN["移除WEB用户成功"] = "移除WEB用户成功."
	zhCN["模块创建成功"] = "模块创建成功."
	zhCN["Test运行在"] = "%s 运行在 [%s]"
	zhCN["ls"] = "ls 服务列表\r\n命令格式: ls\r\n"
	zhCN["add"] = "add 添加服务，注意服务名称不能使用命令关键字\r\n命令格式: add <服务名称> [工程路径] [HTTP服务IP:端口]\r\n"
	zhCN["ctp"] = "ctp 创建工程目录\r\n命令格式: ctp <工程路径>\r\n例如:ctp C:/jus/project/\r\n"
	zhCN["stp"] = "stp 设置工程目录\r\n命令格式: stp <服务名称> <工程路径>\r\n例如:stp test C:/jus/project/\r\n"
	zhCN["ctf"] = "ctf 创建模块页\r\n命令格式: ctf [-创建方式(-h|m|s|r)] <服务名称> <模块全路径>\r\n例如:ctf test component.Test\r\nctf test -hr component.Test\r\n"
	zhCN["release"] = "release 发布工程\r\n命令格式: release <服务名称> [工程路径]\r\n例如:release test C:/jus/project/\r\n"
	zhCN["run"] = "run 启动服务\r\n命令格式: run <服务名称> [IP:端口], 例如:run test 127.0.0.1:1511\r\n"
	zhCN["run_not_set"] = "未设置发布目录，不能启动服务器;\r\n您可以使用stp命令设置发布目录。"
	zhCN["shutdown"] = "shutdown 停止服务\r\n命令格式: shutdown <服务名称>\r\n"
	zhCN["rm"] = "rm 移除服务\r\n命令格式: rm <服务名称>\r\n"
	zhCN["lw"] = "lw 显示指定服务节点下Websocket连接用户\r\n命令格式: lw <服务名称> [-h]\r\n"
	zhCN["info"] = "info 项目信息\r\n命令格式: rm <服务名称>\r\n"
	zhCN["set"] = "set 设置项目信息\r\n命令格式: set <服务名称> <属性名称> <属性值> [属性值...]\r\n"
	zhCN["ret"] = "ret 移除项目信息\r\n命令格式: ret <服务名称> <属性名称>\r\n"
	zhCN["send"] = "send 通过Websocket向指定节点发送数据\r\n命令格式: send <服务名称> <用户ID> <UUID> <内容>\r\n"
	zhCN["exit"] = "exit 退出\r\n命令格式: exit\r\n"
	zhCN["lang"] = "lang 语言设置.\r\n命令格式: lang <zh/cn>\r\n"
	zhCN["version"] = "version 软件版本号.\r\n命令格式: version\r\n"
	zhCN["nat"] = "nat 它可以测试HTTP客户端请求的内容代码，并将其打印到屏幕上\r\n"
	zhCN["-c"] = "-c 关闭控制台输入功能\r\n命令格式: -c\r\n"
	zhCN["webc"] = "webc 启动远程HTTP控制端通讯功能\r\n命令格式: webc [HTTP服务IP:端口]\r\n"
	zhCN["bat"] = "bat 执行本程序的批处理文件，您可以执行多套批处理命令\r\n命令格式：bat <文件名称> [文件名称...]\r\n"
	zhCN["stat"] = "stat 获取当前文件执行状态，例如时间等\r\n命令格式：stat\r\n"
	enCH["文件不存在"] = "The '%s' file isn't exist. \r\n"
	enCH["添加成功"] = "The [%s] added successfully."
	enCH["已经添加"] = "[%s] was added.\r\n"
	enCH["项目已存在"] = "The project [%s] is exist.\r\n"
	enCH["建立项目"] = "create project [%s].\r\n"
	enCH["项目挂载在"] = "The project mount at[%s] server."
	enCH["遍历结束"] = "----list over----"
	enCH["加载系统路径错误"] = "load sys path has errors."
	enCH["遍历运行"] = "%s. %s\tRunning\t%s\t%s\t%s"
	enCH["遍历停止"] = "%s. %s\tStopping\t%s\t%s\t%s"
	enCH["WS运行"] = "%s. %s\t%s\t%s\t%s\t%s"
	enCH["WS停止"] = "%s. %s\t%s\t%s\t%s\t%s"
	enCH["遍历未初始化"] = "<Uninitialized>"
	enCH["不存在工程"] = "The [%s] isn't exist,so set project dir is error."
	enCH["不存在服务"] = "The [%s] services isn't exits,or the services was be removed."
	enCH["移除失败"] = "[%s] remove failed."
	enCH["移除成功"] = "[%s] remove success."
	enCH["设置成功"] = "[%s] setted."
	enCH["属性移除成功"] = "%s remove success."
	enCH["属性移除失败"] = "%s remove failed,mybe the attributes of [%s] isn't exist."
	enCH["服务关闭失败"] = "%s Services stoped failed."
	enCH["工程设置成功"] = "The [%s] setted in [%s]."
	enCH["工程设置失败"] = "The [%s] setted in [%s] failed."
	enCH["服务正在启动"] = "The [%s] starting at  [%s]"
	enCH["关闭服务"] = "%s Stop [%s]"
	enCH["发布完成"] = "----Release Complete----"
	enCH["添加WEB用户成功"] = "Add Web Controller [%s] Success."
	enCH["移除WEB用户成功"] = "Remove Web Controller [%s] Success."
	enCH["模块创建成功"] = "Create Module Success."
	enCH["Test运行在"] = "%s Running at [%s]"
	enCH["-c"] = "Close Console Input Method."
	enCH["ls"] = "ls Show services list.\r\nCOMMAND: ls\r\n"
	enCH["add"] = "add Add services and don't use command as services name.\r\nCOMMAND: add <Service Name> [Project Path] [HTTP Service IP:PORT]\r\n"
	enCH["ctp"] = "ctp create project dir.\r\nCOMMAND: ctp <Project Path>\r\nFor Example:ctp C:/jus/project/\r\n"
	enCH["stp"] = "stp set project dir.\r\nCOMMAND: stp <Service Name> <Project Path>\r\nFor Example:stp test C:/jus/project/\r\n"
	enCH["ctf"] = "ctf create module file.\r\nCOMMAND: ctf [-Create Method(-h|m|s|r)] <Service Name> <Project Path>\r\nFor Example:ctf test component.Test\r\nctf test -hr component.Test\r\n"
	enCH["release"] = "release release project.\r\nCOMMAND: release <Service Name> [Project Path]\r\nFor Example:release test C:/jus/project/\r\n"
	enCH["run"] = "run Start service.\r\nCOMMAND: run <Service Name> [IP:PORT], For Example:run test 127.0.0.1:1511\r\n"
	enCH["run_not_set"] = "The publishing directory is not set, the server cannot be started; \r\nYou can use the STP command to set up the publish directory."
	enCH["shutdown"] = "shutdown Shutdown Service.\r\nCOMMAND: shutdown <Service Name>\r\n"
	enCH["rm"] = "rm Remove Service.\r\nCOMMAND: rm <Service Name>\r\n"
	enCH["lw"] = "lw display websocket list of Service\r\nCOMMAND: lw <Service Name> [-h]\r\n"
	enCH["info"] = "info The project infomation\r\nCOMMAND: rm <Service Name>\r\n"
	enCH["set"] = "set Set project attributes.\r\nCOMMAND: set <Service Name> <AttributeName> <Value> [Value...]\r\n"
	enCH["ret"] = "ret Remove project attributes.\r\nCOMMAND: set <Service Name> <AttributeName>\r\n"
	enCH["send"] = "send push data to Service by websocket.\r\nCOMMAND: send <Service Name> <User ID> <UUID> <Value>\r\n"
	enCH["exit"] = "exit Exit.\r\nCOMMAND: exit\r\n"
	enCH["lang"] = "lang Language Setting.\r\nCOMMAND: lang <zh/cn>\r\n"
	enCH["version"] = "version Software Version.\r\nCOMMAND: version\r\n"
	enCH["nat"] = "nat It's can test http request medhod and print request code.\r\n"
	enCH["-c"] = "-c Close Console Input Method.\r\nCOMMAND: -c\r\n"
	enCH["webc"] = "webc Start HTTP client server to this.\r\nCOMMAND: webc [HTTP Service IP:PORT]\r\n"
	enCH["bat"] = "bat Execute local batch file,you can execute manay batch files.\r\n\r\nCOMMAND：bat <batch file Name> [batch file Name...]\r\n"
	enCH["stat"] = "stat get application status，for example time and so on.\r\nCOMMAND：stat\r\n"
}

/**
 * 服务器列表
 */
var serverList map[string]*UIServer = make(map[string]*UIServer)
var testHandle map[string]*TestServer = make(map[string]*TestServer)
var webCList map[string]*websocket.Conn = make(map[string]*websocket.Conn) //webControl用户列表
var SysPath string
var SysLibPath string
var SysStartDate string
var _Count_ int = 0

/**
 * 创建工程目录
 * @param path 	目录路径
 */
func CreateProjectDir(path string) string {
	abs, _ := filepath.Abs(path)
	if Exist(path) {
		DevPrintln(335, lang["项目已存在"], abs)
		return lang["项目已存在"]
	}
	os.MkdirAll(path, 0777)
	os.MkdirAll(path+"/lib/img", 0777)   //图片
	os.MkdirAll(path+"/lib/css", 0777)   //css
	os.MkdirAll(path+"/lib/js", 0777)    //javascript库
	os.MkdirAll(path+"/lib/font", 0777)  //字体
	os.MkdirAll(path+"/lib/wasm", 0777)  //web汇编文件
	os.MkdirAll(path+"/.serv", 0777)     //配置文件项
	os.MkdirAll(path+"/.serv/use", 0777) //格式化命令
	os.MkdirAll(path+"/.serv/pub", 0777) //发布配置

	s, _ := filepath.Abs("lib/js")
	Copy(s, path+"/lib/js", "")
	s, _ = filepath.Abs("lib/core/icon/")
	Copy(s, path+"/lib/img", "")
	f, e := os.Create(path + "/index.html")
	defer f.Close()
	if e == nil {
		data, _ := GetBytes("./lib/core/template/index.template")
		f.Write(data)
	} else {
		return e.Error()
	}

	f, e = os.Create(path + "/Index.ui")
	defer f.Close()
	if e == nil {
		data, _ := GetCode("./lib/core/template/codeIndex.template")
		data = strings.Replace(data, "{@content}", abs, -1)
		f.Write([]byte(data))
	} else {
		return e.Error()
	}

	f, e = os.Create(path + "/.uisys")
	defer f.Close()
	if e == nil {
		f.WriteString("release-path " + filepath.Dir(abs) + "/" + filepath.Base(abs) + "-release/")
	} else {
		return e.Error()
	}
	DevPrintln(2, lang["建立项目"], abs)
	tName := GetName()
	commandEvt("add " + tName + " " + abs)
	DevPrintln(240, lang["项目挂载在"], tName)
	_Count_++
	return ""
}

func GetName() string {
	name := ""
	for true {
		name = "a" + strconv.Itoa(_Count_)
		_Count_++
		if serverList[name] == nil {
			return name
		}
	}
	return ""
}

func ColorList() {
	for i := 0; i <= 99; i++ {
		DevPrintln(i, "%s,Color", strconv.Itoa(i))
	}
}

func DevPrintln(i int, value ...interface{}) string {
	value[0] = "  " + value[0].(string) + "\r\n"
	return DevPrint(i, value...)
}

/**
 * 控制服务器
 */
var webc *http.Server = nil //web服务控制指针
func webControl(addr string) {
	go func() {
		fmt.Println("Web Control Server Started At: [" + addr + "]. Use protocol https")
		handler := http.NewServeMux()
		handler.HandleFunc("/", root)
		handler.Handle("/ws", websocket.Handler(wsHandler))
		webc = &http.Server{Addr: addr, Handler: handler}
		var err error = nil
		err = webc.ListenAndServeTLS("lib/manager/ssl/cert.pem", "lib/manager/ssl/key.pem")
		if err != nil {
			fmt.Println("status:", err)
		}
		fmt.Println("JUS Server END.")
	}()
}

//控制首页
func root(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		http.Redirect(w, req, "/webc.html", http.StatusFound)
		return
	}
	if IsType(req.URL.Path, ".ui") {
		req.URL.Path = Substring(req.URL.Path, 0, StringLen(req.URL.Path)-5)
		jusEvt(w, req, ".ui")
		return
	}
	if req.URL.Path == "/ui-sys.js" {
		b, e := GetCode("lib/core/parser/module.tpl")
		b0, e0 := GetCode("lib/core/parser/module_base.tpl")
		b1, e1 := GetCode("lib/core/parser/module_manager.tpl")
		b = Replace(b, "{@base}", b0)
		b = Replace(b, "{@manager}", b1)
		w.Header().Add("Content-Type", "text/javascript; charset=utf-8")
		if e != nil || e0 != nil || e1 != nil {
			be := []byte("alert('UI-SYS Error.');")
			w.Header().Add("Content-Length", strconv.Itoa(len(be)))
			w.Write(be)
		} else {
			be := []byte(b)
			w.Header().Add("Content-Length", strconv.Itoa(len(be)))
			w.Write([]byte(be))
		}

		return
	}
	path := "lib/manager/" + req.URL.Path
	path, _ = filepath.Abs(path)
	value, err := GetBytes(path)
	if err != nil {
		value = []byte("404")
	}
	w.Write(value)
}

var jusDirName string = "/juis/"

func jusEvt(w http.ResponseWriter, req *http.Request, ext string) {
	path := "lib/manager" + req.URL.Path
	if Exist(path) {
		root(w, req)
	} else {
		jus := &JUS{SYSTEM_PATH: "lib", CLASS_PATH: "lib/src/"}
		className := Substring(req.RequestURI, 0, LastIndex(req.RequestURI, ext))
		className = Replace(className, "/", ".")
		if jus.CreateFrom("lib/manager/", "", nil, className) {
			b := jus.ToFormatBytes()
			w.Header().Add("Content-Length", strconv.Itoa(len(b)))
			w.Write(b)
		} else {
			fmt.Println("不存在", className)
		}
	}
}

//广播到客户
//name 	向指定用户名广播，如果填写为空，则向所有用户广播
//value	广播内容
func BroadCast(name string, cls int, value string) {
	d := []byte(DevPrintln(cls, value))
	if name == "" {
		for _, v := range webCList {
			v.Write(d)
		}
	} else {
		if webCList[name] != nil {
			webCList[name].Write(d)
		}
	}

}

func wsHandler(ws *websocket.Conn) {
	user := ""
	msg := make([]byte, 512)
	n, err := ws.Read(msg)
	if err != nil {
		fmt.Println("error.")
		return
	}
	flag := true
	cmdstr := string(msg[:n])
	cmds := FmtCmd(cmdstr)
	if len(cmds) == 3 {
		if cmds[0] == "login" && havUser(cmds) {
			user = cmds[1]
			webCList[user] = ws
			defer func() {
				delete(webCList, user)
			}()
		} else {
			ws.Write([]byte("The Name or Password is wrong."))
			ws.Close()
			fmt.Println("The Name or Password is wrong.")
			return
		}
	} else {
		return
	}

	//hostName, _ := os.Hostname()
	str := DevPrint(11, version+" ")
	//str += DevPrint(880, " Running at "+runtime.GOARCH+" "+runtime.GOOS+" "+hostName+" ")
	ws.Write([]byte(str))
	for {
		n, err := ws.Read(msg)
		if err != nil {
			break
		}
		cmdstr = string(msg[:n])
		DevPrintln(240, cmds[1]+": %s\n", cmdstr)
		flag, cmdstr = commandEvt(cmdstr)
		if !flag {
			ws.Write([]byte("The client will over."))
			break
		}
		_, err = ws.Write([]byte(cmdstr))
		if err != nil {
			fmt.Println(err)
		}
	}
	ws.Close()
	//fmt.Println("一个连接已经结束")
	BroadCast("", 11, "["+user+"] login out.")
}

/**
 * 是否存在此用户名
 */
func havUser(cmds []string) bool {
	if !Exist("conf/") {
		os.Mkdir("conf", 0777)
	}
	if !Exist("conf/conf.mg") {
		f, _ := os.Create("conf/conf.mg")
		defer f.Close()
	}
	data, err := GetCode("conf/conf.mg")
	if err != nil {
		fmt.Println("Get conf/conf.mg has error.", err)
		return false
	}
	cmdLst := FmtCmdList(data)
	for _, v := range cmdLst {
		if len(v) > 1 {
			if v[0] == cmds[1] && v[1] == cmds[2] {
				return true
			}
		}
	}
	time.Sleep(2 * time.Second)
	return false
}

/**
 * 添加属性
 */
func SetData(cmds []string) {
	data, err := GetCode("conf/conf.mg")
	if err != nil {
		fmt.Println(err)
		return
	}
	var pos int = 0
	var obj []string = nil
	command := FmtCmdList(data)

	for i, v := range command {
		if len(v) > 0 && cmds[0] == v[0] {
			pos = i
			obj = v
			break
		}
	}

	if obj == nil {
		command = append(command, cmds)
	} else {
		command[pos] = cmds
	}

	//对源文件备份
	os.Rename("conf/conf.mg", "conf/conf.mgb")
	//生成新文件
	f, e := os.Create("conf/conf.mg")
	defer f.Close()
	if e == nil {
		sb := ""
		for i, l := range command {
			for _, n := range l {
				if Index(n, " ") != -1 {
					n = "\"" + n + "\""
				}
				sb += n + " "

			}
			if i+1 != len(command) {
				sb += "\r\n"
			}

		}
		f.WriteString(sb)
		os.Remove("conf/conf.mgb")
	}
}

/**
 * 移除属性
 */
func RetData(cmds []string) bool {
	success := false
	data, err := GetCode("conf/conf.mg")
	if err != nil {
		fmt.Println(err)
		return success
	}

	lst := make([][]string, 0)
	command := FmtCmdList(data)

	for _, v := range command {
		if len(v) == 0 {
			continue
		}

		if cmds[0] == v[0] {
			success = true
			continue
		}

		lst = append(lst, v)
	}

	if len(command) != len(lst) {
		command = lst
		//对源文件备份
		os.Rename("conf/conf.mg", "conf/conf.mgb")
		//生成新文件
		f, e := os.Create("conf/conf.mg")
		defer f.Close()
		if e == nil {
			sb := ""
			for i, l := range command {
				for _, n := range l {
					if Index(n, " ") != -1 {
						n = "\"" + n + "\""
					}
					sb += n + " "

				}
				if i+1 != len(command) {
					sb += "\r\n"
				}

			}
			f.WriteString(sb)
			os.Remove("conf/conf.mgb")
		}
	}

	return success

}

/**
 * 执行批处理文件
 */
func BatCode(value string, w bool) (bool, string) {
	str := ""
	tmp := ""
	_run := false
	running := true
	if code, err := GetCode(value); err == nil {
		list := strings.Split(code, "\n")
		for _, v := range list {
			v = strings.TrimSpace(v)
			if CharAt(v, 0) == "#" || (StringLen(v) >= 2 && Substring(v, 0, 2) == "//") { //包含注释
				if w {
					str += tmp + "\r\n"
					fmt.Println(v)
				}
				continue
			}
			if len(v) > 0 {
				_run, tmp = commandEvt(v)
				if !_run {
					running = _run
				}
				str += tmp
			}
		}
	}

	return running, str
}

func commandEvt(value string) (bool, string) {
	return command(FmtCmd(value))
}

/**
 * 命令代码
 */
func command(cmds []string) (bool, string) {
	str := ""
	if len(cmds) > 0 {
		switch cmds[0] {
		case "-c": //退出命令行
			str += DevPrintln(2, "Console Input Method Unabled.")
			return false, str
		case "bat": //批处理文件
			if len(cmds) > 1 {
				for i := 1; i < len(cmds); i++ {
					if Exist(cmds[i]) {
						_, v := BatCode(cmds[i], true)
						str += v
					} else {
						str += DevPrintln(335, lang["文件不存在"], cmds[i]) //文件不存在
					}

				}
			} else {
				str = DevPrintln(8, lang["bat"])
			}
			return true, str
		case "ls":
			if len(cmds) > 1 {
				str += "<table class='list'>"
				str += "<tr><th>ID</th><th>Name</th><th>Create Time</th><th>Init Status</th><th>Index Address</th></tr>"
				i := 0
				for key, value := range serverList {
					str += "<tr>"
					if value.Status { //Connect.
						str += "<td>" + strconv.Itoa(i) + "</td><td>" + key + "</td><td>" + value.Datetime.Format("2006-01-02 15:04:05") + "</td><td>" + IfStr(value.RootPath == "", lang["遍历未初始化"], value.RootPath) + "</td><td>" + value.GetProtocol() + "://" + IfStr(Index(value.Addr, ":") == 0, "0.0.0.0"+value.Addr, value.Addr) + "/" + "</td>"
					} else {
						str += "<td>" + strconv.Itoa(i) + "</td><td>" + key + "</td><td>" + value.Datetime.Format("2006-01-02 15:04:05") + "</td><td>" + IfStr(value.RootPath == "", lang["遍历未初始化"], value.RootPath) + "</td><td>" + value.GetProtocol() + "://" + IfStr(Index(value.Addr, ":") == 0, "0.0.0.0"+value.Addr, value.Addr) + "/" + "</td>"
					}
					str += "</tr>"
					i++
				}
				str += "</table>"
			} else {
				i := 0
				for key, value := range serverList {
					if value.Status {
						str += DevPrintln(7, lang["遍历运行"], strconv.Itoa(i), key, value.Datetime.Format("2006-01-02 15:04:05"), IfStr(value.RootPath == "", lang["遍历未初始化"], value.RootPath), value.GetProtocol()+"://"+IfStr(Index(value.Addr, ":") == 0, "0.0.0.0"+value.Addr, value.Addr)+"/")
					} else {
						str += DevPrintln(8, lang["遍历停止"], strconv.Itoa(i), key, value.Datetime.Format("2006-01-02 15:04:05"), IfStr(value.RootPath == "", lang["遍历未初始化"], value.RootPath), value.GetProtocol()+"://"+IfStr(Index(value.Addr, ":") == 0, "0.0.0.0"+value.Addr, value.Addr)+"/")
					}

					i++
				}
				str += DevPrintln(8, lang["遍历结束"])
			}

			return true, str
		case "add": //创建服务
			if len(cmds) > 1 && (zhCN[cmds[1]] == "" || len(cmds[1]) != len([]rune(cmds[1]))) {
				if serverList[cmds[1]] == nil {
					serverList[cmds[1]] = &UIServer{}
					serverList[cmds[1]].CreateServer(SysLibPath, "", "", "/")
					str = DevPrintln(2, lang["添加成功"], cmds[1]) //添加成功
				} else {
					str = DevPrintln(335, lang["已经添加"], cmds[1]) //已经添加
				}
				if len(cmds) > 2 {
					if Exist(cmds[2]) {
						_, str = commandEvt("stp " + cmds[1] + " " + cmds[2])
					} else {
						str = DevPrintln(335, lang["不存在工程"], cmds[2])
						return true, str
					}

				}
				if len(cmds) > 3 {
					_, str = commandEvt("run " + cmds[1] + " " + cmds[3])
				}

			} else {
				str = DevPrintln(8, lang["add"])
			}
			return true, str
		case "nat": //添加请求测试服务
			if len(cmds) > 1 {
				if cmds[1] == "-add" {
					if len(cmds) > 4 {
						testHandle[cmds[2]] = &TestServer{Name: cmds[2], Time: time.Now().Unix()}
						if testHandle[cmds[2]].Start(cmds[3], cmds[4]) {
							str = DevPrintln(2, lang["服务正在启动"], cmds[2], testHandle[cmds[2]].FromIPAddress()+"-->"+testHandle[cmds[2]].ToIPAddress()) //添加成功
						}
					} else {
						testHandle[cmds[2]] = &TestServer{Name: cmds[2], Time: time.Now().Unix()}
						if testHandle[cmds[2]].Start(cmds[3], "") {
							str = DevPrintln(2, lang["Test运行在"], cmds[2], testHandle[cmds[2]].FromIPAddress()+"-->"+testHandle[cmds[2]].ToIPAddress()) //添加成功
						}
					}

				} else if cmds[1] == "-remove" {
					if testHandle[cmds[2]] != nil {
						testHandle[cmds[2]].Shutdown()
						delete(testHandle, cmds[2])
					}
				} else if cmds[1] == "-restart" {
					if testHandle[cmds[2]] != nil {
						testHandle[cmds[2]].Restart()
					}
				} else if cmds[1] == "-stop" {
					if testHandle[cmds[2]] != nil {
						testHandle[cmds[2]].Shutdown()
					}
				} else if cmds[1] == "-log" {
					if testHandle[cmds[2]] != nil {
						var size = 1024 * 1024 * 1024
						var err error
						if len(cmds) > 4 {
							size, err = strconv.Atoi(cmds[4])
							if err != nil {
								str += DevPrintln(337, "文件大小输入有误") //添加成功
							}
						}
						testHandle[cmds[2]].SetLog(cmds[3], size)
					}
				} else if cmds[1] == "-h" {
					str += "<table class='list'>"
					str += "<tr><th>Name</th><th>From IP Address</th><th>To IP Address</th><th>Status</th><th>Connect Time</th><th>Connect Count</th><th>Log Path</th><th>Log Size</th></tr>"
					for _, v := range testHandle {
						str += "<tr>"
						if v.Running() { //Connect.
							str += "<td>" + v.Name + "</td><td>" + v.FromIPAddress() + "</td><td>" + v.ToIPAddress() + "</td><td>Running.</td><td>" + time.Unix(v.Time, 0).Format("2006-01-02 15:04:05") + "</td><td>" + v.ConnectStatus() + "</td><td>" + v.GetLogPath() + "</td><td>" + v.LogStatus() + "</td>"
						} else {
							str += "<td>" + v.Name + "</td><td>" + v.FromIPAddress() + "</td><td>" + v.ToIPAddress() + "</td><td>Stopping.</td><td>" + time.Unix(v.Time, 0).Format("2006-01-02 15:04:05") + "</td><td>" + v.ConnectStatus() + "</td><td>" + v.GetLogPath() + "</td><td>" + v.LogStatus() + "</td>"
						}
						str += "</tr>"
					}
					str += "</table>"
				}

			} else { //显示nat列表

				for key, value := range testHandle {
					if value.Running() {
						str += DevPrintln(7, key+"\t["+value.FromIPAddress()+"-->"+value.ToIPAddress()+"]\tRunning\t"+value.ConnectStatus()+"\t"+value.LogStatus()+"\t"+value.GetLogPath()) //添加成功
					} else {
						str += DevPrintln(8, key+"\t["+value.FromIPAddress()+"-->"+value.ToIPAddress()+"]\tStopping\t"+value.ConnectStatus()+"\t"+value.LogStatus()+"\t"+value.GetLogPath()) //添加成功
					}

				}
			}
			str += DevPrintln(8, lang["遍历结束"])
			return true, str
		case "stp": //设置工程目录
			if len(cmds) == 2 {
				if serverList[cmds[1]] == nil {
					str = DevPrintln(335, lang["不存在服务"], cmds[1])
				} else {
					str = DevPrintln(8, cmds[1]+" "+serverList[cmds[1]].RootPath)
				}
			} else if len(cmds) > 2 {
				if serverList[cmds[1]] == nil {
					str = DevPrintln(335, lang["不存在服务"], cmds[1])
				} else {
					result := serverList[cmds[1]].SetProject(cmds[2])
					if result == 2 { //module 版本不一致
						str = DevPrintln(110, lang["modulev"], cmds[1])
						str += DevPrintln(2, lang["工程设置成功"], cmds[1], serverList[cmds[1]].RootPath)
					} else if result == 1 {
						str = DevPrintln(2, lang["工程设置成功"], cmds[1], serverList[cmds[1]].RootPath)
					} else if result == 0 {
						str = DevPrintln(2, lang["工程设置失败"], cmds[1], serverList[cmds[1]].RootPath)
					}

				}
			} else {
				str = DevPrintln(8, lang["stp"])
			}
			return true, str
		case "ctf": //创建模块文件
			if len(cmds) == 2 {
				if serverList[cmds[1]] == nil {
					str = DevPrintln(335, lang["不存在服务"], cmds[1])
				} else {
					str = DevPrintln(8, cmds[1]+" "+serverList[cmds[1]].RootPath)
				}

			} else if len(cmds) > 3 {
				if serverList[cmds[1]] == nil {
					str = DevPrintln(335, lang["不存在服务"], cmds[1])
				} else {
					serverList[cmds[1]].CreateModule(cmds[2], cmds[3])
					str = DevPrintln(2, lang["模块创建成功"])
				}

			} else if len(cmds) > 2 {
				if serverList[cmds[1]] == nil {
					str = DevPrintln(335, lang["不存在服务"], cmds[1])
				} else {
					serverList[cmds[1]].CreateModule("-h", cmds[2])
					str = DevPrintln(2, lang["模块创建成功"])
				}
			} else {
				str = DevPrintln(8, lang["ctf"])
			}
			return true, str
		case "send": //向服务器的WebSocket用户发送信息
			if len(cmds) > 1 {
				if serverList[cmds[1]] == nil {
					str = DevPrintln(335, lang["不存在服务"], cmds[1])
				} else if len(cmds) > 4 {
					serverList[cmds[1]].Send(cmds[2], cmds[3], cmds[4])
				} else if len(cmds) > 3 {
					serverList[cmds[1]].Send(cmds[2], "UUID", cmds[3])
				}
			} else {
				str = DevPrintln(8, lang["send"])
			}
			return true, str
		case "run": //运行工程
			if len(cmds) > 1 {
				port := ":80"
				if len(cmds) > 2 {
					port = cmds[2]
				}
				if serverList[cmds[1]] == nil {
					str = DevPrintln(335, lang["不存在服务"], cmds[1])
				} else {
					if serverList[cmds[1]].RootPath != "" {
						str = DevPrintln(2, lang["服务正在启动"], cmds[1], port)
						serverList[cmds[1]].Start(port, BroadCast)
					} else {
						str = DevPrintln(335, lang["run_not_set"])
					}
				}
			} else {
				str = DevPrintln(8, lang["run"])
			}
			return true, str
		case "shutdown":
			if len(cmds) > 1 {
				if serverList[cmds[1]] == nil {
					str = DevPrintln(335, lang["不存在服务"], cmds[1])
				} else {
					if serverList[cmds[1]].Close() != nil {
						str = DevPrintln(8, lang["服务关闭失败"], cmds[1])
					}
				}
			}
			if str == "" {
				str = DevPrintln(2, lang["关闭服务"], cmds[1], cmds[1])
			}
			return true, str
		case "rm":
			if len(cmds) > 1 {
				if serverList[cmds[1]] == nil {
					str = DevPrintln(335, lang["不存在服务"], cmds[1])
				} else {
					if serverList[cmds[1]].Destroy() == nil {
						delete(serverList, cmds[1])
						str = DevPrintln(2, lang["移除成功"], cmds[1])
					} else {
						str = DevPrintln(8, lang["移除失败"], cmds[1])
					}
				}
			}
			return true, str
		case "ctp": //增加一个项目
			if len(cmds) > 1 {
				str = DevPrintln(2, CreateProjectDir(cmds[1]))
			} else {
				str = DevPrintln(8, lang["ctp"])
			}
			return true, str
		case "release": //发布项目
			if len(cmds) > 1 {
				if serverList[cmds[1]] == nil {
					str = DevPrintln(335, lang["不存在服务"], cmds[1])
				} else {
					serverList[cmds[1]].Release()
					str = DevPrintln(8, lang["发布完成"])
				}
			} else {
				str = DevPrintln(8, lang["release"])
			}
			return true, str
		case "info": //查看项目设置
			if len(cmds) > 1 {
				if serverList[cmds[1]] == nil {
					str = DevPrintln(335, lang["不存在服务"], cmds[1])
				} else {
					for j, v := range serverList[cmds[1]].GetData() {
						for i, n := range v {
							if len(v) == 0 {
								continue
							}
							if len(v) > 1 && i == 0 {
								if j%2 == 0 {
									str += DevPrint(63, n)
								} else {
									str += DevPrint(111, n)
								}
							} else {
								str += DevPrint(7, " "+n)
							}
						}
						DevPrintln(7, "")
						str += "\r\n"
					}
				}

			} else {
				str = DevPrintln(8, lang["info"])
			}
			return true, str
		case "color":
			for i := 0; i < 256; i++ {
				str += DevPrintln(i, "     "+strconv.Itoa(i))
			}
			return true, str

		case "set": //设置项目变量
			if len(cmds) > 3 {
				if serverList[cmds[1]] == nil {
					str = DevPrintln(335, lang["不存在服务"], cmds[1])
				} else {
					serverList[cmds[1]].SetData(cmds[2:])
					str = DevPrintln(2, lang["设置成功"], cmds[1])
				}

			} else {
				str = DevPrintln(8, lang["set"])
			}

			return true, str
		case "ret": //移除项目变量
			if len(cmds) > 2 {
				if serverList[cmds[1]] == nil {
					str = DevPrintln(335, lang["不存在服务"], cmds[1])
				} else {
					if serverList[cmds[1]].RetData(cmds[2:]) {
						str = DevPrintln(2, lang["属性移除成功"], cmds[1])
					} else {
						str = DevPrintln(2, lang["属性移除失败"], cmds[1], cmds[2])
					}

				}

			} else {
				str = DevPrintln(8, lang["ret"])
			}

			return true, str
		case "lang":
			if len(cmds) > 1 {
				if cmds[1] == "en" {
					lang = enCH
				} else if cmds[1] == "zh" {
					lang = zhCN
				}
			} else {
				str = DevPrintln(7, "您可以输入 <zh> 或者 <en> 来选取中文或者英文.")
				str = DevPrintln(7, "You can chosen 'zh' for china or 'en' for english.")
			}
			return true, str
		case "webc": //WEB程序控制
			if len(cmds) == 2 { //默认端口3690
				webControl(cmds[1])
			} else if len(cmds) == 3 { //-del
				if cmds[1] == "-del" {
					if RetData(cmds[2:]) {
						str = DevPrintln(2, lang["移除WEB用户成功"], cmds[2])
					}
				} else {
					str = DevPrintln(2, lang["移除失败"], cmds[2])
				}

			} else if len(cmds) == 4 { //-add
				if cmds[1] == "-add" {
					SetData(cmds[2:])
					str = DevPrintln(2, lang["添加WEB用户成功"], cmds[2])
				} else {
					str = DevPrintln(4, lang["移除失败"], cmds[2])
				}
			} else {
				webControl(":3690")
			}
			return true, str
		case "lw": //显示目前socket链接用户
			if len(cmds) > 1 {
				if serverList[cmds[1]] == nil {
					str = DevPrintln(335, lang["不存在服务"], cmds[1])
				} else {
					if len(cmds) > 2 && cmds[2] == "-h" {
						str += "<table class='list'>"
						str += "<tr><th>ID</th><th>Name</th><th>IP Address</th><th>Remote Addr</th><th>Local Addr</th><th>Connect Time</th></tr>"
						for i, v := range serverList[cmds[1]].WebsocketList() {
							str += "<tr>"
							if v.Connected { //Connect.
								str += "<td>" + strconv.Itoa(i) + "</td><td>" + v.Name + "</td><td>" + v.IP_Address + "</td><td>" + v.RemoteAddr + "</td><td>" + v.LocalAddr + "</td><td>" + time.Unix(v.Time, 0).Format("2006-01-02 15:04:05") + "</td>"
							} else {
								str += "<td>" + strconv.Itoa(i) + "</td><td>" + v.Name + "</td><td>" + v.IP_Address + "</td><td>" + v.RemoteAddr + "</td><td>" + v.LocalAddr + "</td><td>" + time.Unix(v.Time, 0).Format("2006-01-02 15:04:05") + "</td>"
							}
							str += "</tr>"
						}
						str += "</table>"
					} else {
						for i, v := range serverList[cmds[1]].WebsocketList() {
							if v.Connected { //Connect.
								str += DevPrintln(7, lang["WS运行"], strconv.Itoa(i), v.Name, v.IP_Address, v.RemoteAddr, v.LocalAddr, time.Unix(v.Time, 0).Format("2006-01-02 15:04:05"))
							} else {
								str += DevPrintln(8, lang["WS运行"], strconv.Itoa(i), v.Name, v.IP_Address, v.RemoteAddr, v.LocalAddr, time.Unix(v.Time, 0).Format("2006-01-02 15:04:05"))
							}
						}
						str += DevPrintln(8, lang["遍历结束"])
					}

				}
			} else {
				str = DevPrintln(8, lang["lw"])
			}
			return true, str
		case "version":
			str = DevPrintln(496, version)
			return true, str
		case "stat": //转换地址，获取当前程序
			str = DevPrintln(7, SysStartDate+"\r\nNow  "+time.Now().Format("2006-01-02 15:04:05"))
			return true, str

		case "--help":
			str += DevPrintln(7, lang["lang"])
			str += DevPrintln(7, lang["ls"])
			str += DevPrintln(7, lang["add"])
			str += DevPrintln(7, lang["ctp"])
			str += DevPrintln(7, lang["stp"])
			str += DevPrintln(7, lang["ctf"])
			str += DevPrintln(7, lang["release"])
			str += DevPrintln(7, lang["update"])
			str += DevPrintln(7, lang["send"])
			str += DevPrintln(7, lang["run"])
			str += DevPrintln(7, lang["shutdown"])
			str += DevPrintln(7, lang["rm"])
			str += DevPrintln(7, lang["lw"])
			str += DevPrintln(7, lang["info"])
			str += DevPrintln(7, lang["set"])
			str += DevPrintln(7, lang["ret"])
			str += DevPrintln(7, lang["version"])
			str += DevPrintln(7, lang["nat"])
			str += DevPrintln(7, lang["-c"])
			str += DevPrintln(7, lang["webc"])
			str += DevPrintln(7, lang["bat"])
			str += DevPrintln(7, lang["stat"])
			str += DevPrintln(7, lang["exit"])

			return true, str
		case "exit":
			return false, "quit"
		default:
			if len(cmds) > 1 {
				if (serverList[cmds[0]]) != nil {
					t := cmds[1]
					cmds[1] = cmds[0]
					cmds[0] = t
					return command(cmds)
				} else {
					str = DevPrintln(335, lang["不存在服务"], cmds[0])
				}
			}
			return true, str
		}
	}

	return true, ""

}

var exitFlag bool = true

func httpPost() {
	resp, err := http.Post("http://www.airoot.cn/_version", "application/x-www-form-urlencoded", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	mess := string(body)
	if strings.TrimSpace(mess) != version {
		DevPrintln(14, "Download Latest Version: http://www.airoot.cn/")
	}
}

/**
 *
 */
func main() {
	lang = enCH
	SetConsoleTitle(version)
	/*
		arr := make([]color.RGBA, 0)

		fImg1, _ := os.Open("aa1.png")
		defer fImg1.Close()

		img, str, e := image.Decode(fImg1)
		if e != nil {
			fmt.Println("---." + str)
		} else {
			b := img.Bounds()
			f, _ := os.Create("data.dat")
			for y := 0; y < b.Size().Y; y++ {
				for x := 0; x < b.Size().X; x++ {
					r, g, b, a := img.At(x, y).RGBA()
					c := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
					f.WriteString(strconv.Itoa(int(r>>8)) + "," + strconv.Itoa(int(g>>8)) + "," + strconv.Itoa(int(b>>8)) + "," + strconv.Itoa(int(a>>8)) + ",")
					arr = append(arr, c)
				}

			}
			f.Close()
			_drawImage(arr, b.Size().X)
		}
	*/
	arri := []int{
		0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255}
	drawImage(arri, 36)
	fmt.Println()
	DevPrint(11, "  "+version+" ")
	DevPrintln(7, "http://www.airoot.cn/")
	DevPrintln(3, "Bootup successful ")
	SysStartDate = "Date " + time.Now().Format("2006-01-02 15:04:05")
	DevPrintln(3, SysStartDate+" ")
	fmt.Println()

	var err error
	s, err := os.Executable()
	if err != nil {
		fmt.Println(lang["加载系统路径错误"])
		return
	}
	SysPath = filepath.Dir(s)
	SysLibPath = SysPath + "/lib"
	if len(os.Args) > 1 {
		fmt.Println("  Sys: " + SysPath)
		fmt.Println("  Sys Library: " + SysLibPath)
	}

	//默认传入参数
	args := ""
	for _, v := range os.Args[1:] {
		args += v + " "
	}
	if args != "" {
		fmt.Println("ARGS", args)
	}
	running := true
	if len(os.Args) == 2 && (Index(args, "/") != -1 || Index(args, "\\") != -1) {
		args = "add d0 " + args + " :80"
	} else {
		if !Exist("boot.conf") {
			f, _ := os.Create("boot.conf")
			defer f.Close()
		}
		running, _ = BatCode("boot.conf", true) //程序默认执行一个控制类
	}
	go httpPost()
	//键盘输入
	quit := ""
	if running {
		running, quit = commandEvt(args)
		if running {
			reader := bufio.NewReader(os.Stdin)
			for running {
				data, _, _ := reader.ReadLine()
				running, quit = commandEvt(string(data))
			}
		}
	}
	for exitFlag && quit != "quit" {
		time.Sleep(1 * time.Second)
	}
	fmt.Println("End.")
	os.Exit(0)

}
