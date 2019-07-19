# UI-SYSTEM
[http://www.airoot.cn/](http://www.airoot.cn/)
# 说明
- UI-SYSTEM 是一套轻巧、高性能、静态前端系统，可以快速构建稳健的WEB服务。  
整套系统采用了混合式的开发思想，即“选择合适的技术做擅长的事，各尽其职，优势融合”。  
- UI-SYSTEM 设计思想是 “简单明了，直截了当”，让前端工程师直接进入开发状态。  
- UI-SYSTEM 在模块化设计上，采用微模块方案，你可以把他理解为DNA和氨基酸的配合，一切展现模块都是即时组合，并且在渲染上做了大量优化，能提供原生性能的渲染体验。  
- UI-SYSTEM 提供了命令行控制服务和自动配置文件两种方式部署，服务采用热更新方式，动态切换服务参数。  
- UI-SYSTEM 采用Go作为开发语言，充分利用了Go的高并发能力，很高兴选择了Go作为服务开发（之前一直使用Java，实际上GO在复杂业务上处理速度已经远远高于Java这门语言，网上很多人用JIT命中高例子和Go对比是不正确的。）

<table>
    <tr>
        <td><img src='doc/images/h.png' /></td>
        <td valign='top'>
            <h2>支持自定义HTML标签</h2>
可以自己动手编写一个标签，完全由你定义， 支持使用HTML、CSS、JavaScript混合编写。
        </td>
    </tr>
    <tr>
        <td><img src='doc/images/package.png' /></td>
        <td valign='top'>
           <h2>支持HTML封装为控件</h2>
HTML页就是控件(模块页)，你可以发布这个模块到任何人的Super HTML平台上使用。
        </td>
    </tr>
    <tr>
        <td><img src='doc/images/as.png' /></td>
        <td valign='top'>
            <h2>支持高级Script语法糖</h2>
支持更好的面向对象特性，同时兼容JavaScript的写法。
        </td>
    </tr>
    <tr>
        <td><img src='doc/images/kz.png' /></td>
        <td valign='top'>
            <h2>支持继承和扩展</h2>
支持继承HTML页，同时可以对HTML页元素扩展。 支持public、private、static等关键字。
        </td>
    </tr>
    <tr>
        <td><img src='doc/images/comment.png' /></td>
        <td valign='top'>
            <h2>自动生成API文档</h2>
完成一个HTML模块页，平台会自动生成模块的API说明文档。不用人工干涉。
        </td>
    </tr>
</table>
 

## 前端渲染平台
UI-SYSTEM (代号), 是一款编译系统,可以将WEB页封装和发布。它不是JS框架、也不是DOM框架。
其内建的编译引擎可以升级当前WEB平台，使它具备面向对象的特性。
UI-SYSTEM的优势是抛开了以往语言级别的面向对象设计，采用全面向对象设计（包含HTML、CSS、JavaScript等）。
UI-SYSTEM的开发和普通的WEB开发没有什么区别，几乎零学习成本。
## 与Node.js区别
首先，UI-SYSTEM不负责后台开发，也就是说它只面向前端开发，你不可以把它作为动态服务器来学习。
### 优势
UI-SYSTEM 不需要安装、即学即用。程序体积小，在性能及易用性上都做了深度考量。    
UI-SYSTEM 即面向普通WEB开发，同时也面向模块化开发。   
UI-SYSTEM 即适合初学者人群，同时也适合高级JavaScript开发人群。

## 面向对象的模块化能力
UI-SYSTEM 相比现在的WEB 模块化框架,其在设计之初就是以模块化思想为基础，UI-SYSTEM 的开发体验更加友好，拥有高效的编程方案，同时兼备了友好的开发体验。
UI-SYSTEM 在设计之初，就将模块化、资源管理、整合其他框架作为其创建目标，所以你会发现，当你开启UI-SYSTEM编程第一个Hello World的页面时，其背后已经帮你解决了许多关键问题。不需要什么全家桶、不需要什么机制，加载、渲染、卸载、资源管理、网络资源、外部包等等，都已经考虑了。

## UI-SYSTEM 是重量级框架吗？
首先，UI-SYSTEM 不是框架，他是一个面向前端的渲染平台，采用Google Golang 开发。UI-SYSTEM的默认发布框架包仅仅只有20KB。UI-SYSTEM的设计初衷就是用最简单的代码解决复杂的问题，平台发布代码均以高扩展性、高复用性、高维护性作为思想基础，你可以想象到，一套轻量级的代码量去构建重量级别的项目。  
因此，UI-SYSTEM是轻量级的平台，但是可以承接复杂的项目。UI-SYSTEM发布的项目是轻量级的，但是具备重量级项目的表现力。 


# 如何使用
> WINDOWS  
[点击这里下载 UI-SYSTEM 的演示包](https://github.com/uucckk/AIroot-UISYS-LIB/archive/master.zip)
- <b>方式1</b> 
1. window系统运行uisys.exe  
2. 从example里面添加一个工程，写法如下：

  ```linux
  pub example/v1 :80
  ```
其中:80是开启服务的端口好，也可以指定发布地址，如127.0.0.1:80  
3. 打开chrome浏览器，输入：http://127.0.0.1/
- <b>方式2</b>  
1. 也可以直接将您的工程目录拖动到uisys.exe图标上。  
![conv_ops](README/img/dragstart.gif)
2. 确保弹出的控制台没有端口占用错误错误。
3. 打开chrome浏览器，输入：http://127.0.0.1/
> LINUX 和 CENTOS

> DARWIN

> ARM

# 服务运行方式
UI-SYS 的服务节点运行方式，默认是以 开发(Development) 模式运行的，在此模式下，UISYS的WEB SERVER 会对扩展名为\*.ui,\*.es 两种文件进行实时解析，当我们要发布的时候，需要切换到静态发布模式，届时，所有的*.ui,*.es都会变成*.ui.hmtl这样的扩展名。
如果想了解UISYS的*.ui,*.es扩展名文件的概念可以点击<a href='README/module.md'>这里</a>。
静态模式命令如下：
```linux
$> pub example/v1 :80 -s
```
命令的格式是在原有格式后加上 -s 即可（-s 代表 static）。

我们介绍下 UISYS 平台的基本命令，包括<b>服务控制命令</b>和<b>项目参数设置</b>命令。
# 命令解析
## 服务控制命令
### 1. <b>help</b>  
获取帮助信息
```linux
$> --help
---------------------------------------------
    lang Language Setting.
    COMMAND: lang <zh/cn>

    pub Publishing websites.
    COMMAND: pub <path> [HTTP Service IP:PORT]

    ls Show services list.
    COMMAND: ls

    add Add services and don't use command as services name.
    COMMAND: add <Service Name> [Project Path] [HTTP Service IP:PORT]
...
```

### 2. <b>version</b>  
获取软件版本
```linux
$> version
---------------------------------------------
    AIroot UI-SYSTEM 0.9.5beta
```

### 3. <b>pub</b>  
发布指定目录为网站  
命令格式: pub <path> [HTTP Service IP:PORT]
```linux
$> pub example/v1
---------------------------------------------
  The [p0] setted in [E:\UISYS-RELEASE\example\v1].
  The [p0] starting at  [:80]
  WEB Server Started At: [:80]. Use protocol http
```
可以指定端口：
```linux
$> pub example/v1 :8888
---------------------------------------------
  The [p0] setted in [E:\UISYS-RELEASE\example\v1].
  The [p0] starting at  [:8888]
  WEB Server Started At: [:8888]. Use protocol http
```
可以指定绝对路径：
```linux
$> pub E:\UISYS-RELEASE\example\v1 :8888
---------------------------------------------
  The [p0] setted in [E:\UISYS-RELEASE-0.9.5BETA\example\v1].
  The [p0] starting at  [:8888]
  WEB Server Started At: [:8888]. Use protocol http
```
对于带空格的路径可以用引号括起来，如下：
```linux
$> pub "E:\UISYS RELEASE\example\v1" :8888
---------------------------------------------
  The [p0] setted in [E:\UISYS RELEASE\example\v1\example\v1].
  The [p0] starting at  [:8888]
  WEB Server Started At: [:8888]. Use protocol http
```
可以指定https服务
```linux
$> pub "E:\UISYS-RELEASE\example\v1" https://:80
---------------------------------------------
  The [p0] starting at  [https://:80]
  WEB Server Started At: [:80]. Use protocol https
```
可以全部制定：
```linux
$> pub "E:\UISYS-RELEASE\example\v1" https://10.110.10.34:80
---------------------------------------------
  The [p0] starting at  [https://10.110.10.34:80]
  WEB Server Started At: [10.110.10.34:80]. Use protocol https
```

### 4. <b>ls</b>  
列出当前存在的服务节点。
```linux
$> ls
---------------------------------------------
  0. p1 Running 2019-07-10 23:43:28     D:\UISYS-RELEASE\example\v1       http://0.0.0.0:80/
  1. a2 Stopping        2019-07-10 23:43:43     D:\UISYS-RELEASE\example\v2     http:///
  ----list over----
```
### 5. <b>add</b> 
Add services and don't use command as services name.  
添加服务节点，用于挂在被发布的工程。  
注意：服务名称不能使用add作为服务的名字。
- 命令格式: add \<Service Name> [Project Path] [HTTP Service IP:PORT]
```linux
$> add a0 example/v1 :80
---------------------------------------------
  The [a0] setted in [E:\UISYS-RELEASE\example\v1].
  The [a0] starting at  [:80]
  WEB Server Started At: [:80]. Use protocol http
```
也可以只创建服务节点，但是不挂在项目：
```linux
$> add a0
---------------------------------------------
  The [a0] added successfully.
```
如果需要挂在节点，可以通过 <b>stp</b>（set project）命令挂在项目目录：
```linux
$> a0 stp example/v1
---------------------------------------------
  The [a0] setted in [C:\Users\Administrator\Desktop\UISYS-RELEASE-0.9.5BETA\example\v1].
```
此时，我们只是挂在了项目，如果要运行需要使用 <b>run</b> 命令：
```linux
$> run a0 :80
---------------------------------------------
  The [a0] starting at  [:80]
  WEB Server Started At: [:80]. Use protocol http
```
### 6. <b>run</b> 
Start service.  
运行服务节点  
- 命令格式: run \<Service Name> [IP:PORT], For Example:run test 127.0.0.1:1511
```linux
$> run a0
---------------------------------------------
  The [a2] starting at  [:80]
  WEB Server Started At: [:80]. Use protocol http
```
可以指定端口
```linux
$> run a0 :80
---------------------------------------------
  The [a2] starting at  [:80]
  WEB Server Started At: [:80]. Use protocol http
```
可以指定IP
```linux
$> run a0 10.110.10.34:80
---------------------------------------------
  The [a2] starting at  [10.110.10.34:80]
  WEB Server Started At: [10.110.10.34:80]. Use protocol http
```
可以指定https服务
```linux
$> run a0 https://:80
---------------------------------------------
  The [a2] starting at  [https://:80]
  WEB Server Started At: [:80]. Use protocol https
```
可以全部制定：
```linux
$> run a0 https://10.110.10.34:80
---------------------------------------------
  The [a2] starting at  [https://10.110.10.34:80]
  WEB Server Started At: [10.110.10.34:80]. Use protocol https
```
### 7. <b>stop</b> 
停止服务 
- 命令格式: stop \<Service Name>
```linux
$> stop a0
---------------------------------------------
  a0 Stop [a0]
  status: [:80]http: Server closed.
  [:80]JUS Server END.
```
### 8. <b>rm</b> 
移除服务
- 命令格式: rm \<Service Name>
```linux
$> rm a0
---------------------------------------------
  [a0] remove success.
  status: [:80]http: Server closed.
  [:80]JUS Server END.
```


### 9. <b>nat</b> 
实现端口穿透功能。  
- 命令格式: nat <-add/-remove> <Nat Name> <本机端口> <映射机器IP:端口号>  
例如，我们希望将本机的12000端口映射到10.110.10.28的3389端口。  
3389 端口是window服务器的远程桌面服务端口，这样就可以暴露本机的12000端口来对外提供服务。    
写法如下：
```linux
$> nat -add desktop :12000 10.110.10.28:3389
---------------------------------------------
  The [desktop] starting at  [:12000-->10.110.10.28:3389]
  ----list over----
```

查看本平台用了多少个对外映射可以：
```linux
$> nat
---------------------------------------------
  desktop       [:12000-->127.0.0.1:3389]       Running 0
  ----list over----
```
如果要删除这个映射服务可以用一下命令：
```linux
$> nat -remove desktop
---------------------------------------------
>> accept tcp [::]:12000: use of closed network connection
desktop Close havs error:  close tcp [::]:12000: use of closed network connection
  ----list over----
```


### 10. <b>-c</b>
关闭控制台（Console）的输入功能。
```linux
$> -c
---------------------------------------------
  Console Input Method Unabled.
```
### 11. <b>webc</b> 
启动Web版的服务器命令窗口。该功能默认以https发布。
- 命令格式: webc [HTTP Service IP:PORT]
```linux
$> webc
---------------------------------------------
Web Control Server Started At: [:3690]. Use protocol https
```
如果自己设定服务端口，可以用：
```linux
$> webc :10000
---------------------------------------------
Web Control Server Started At: [:10000]. Use protocol https
```
### 12. <b>bat</b> 
执行批处理命令，可以指定多个批处理文件。
- 命令格式：bat <batch file Name> [batch file Name...]  
UI-SYSTEM 可以运行多个WEB服务，因此如果每次服务重启都要手工重新敲击一边太慢了。
我们可以将经常重服务输入的命令写在一个文件或多个文件里。

例如，我们编写一个“config.conf”，如下：
```txt
#发布example/v1工程到80端口
pub example/v1 :80
#发布example/v2工程到90端口
pub example/v2 :90
```
然后保存到uisys.exe 可以访问的目录，例如，放到uisys目录下。
然后再uisys控制台输入命令：
```linux
$> bat config.conf
---------------------------------------------
#发布example/v1工程到80端口
  The [p0] setted in [C:\UISYS-RELEASE\example\v1].
  The [p0] starting at  [:80]
#发布example/v2工程到90端口
  WEB Server Started At: [:80]. Use protocol http
  The [p1] setted in [C:\UISYS-RELEASE\example\v2].
  The [p1] starting at  [:90]
  WEB Server Started At: [:90]. Use protocol http
```
bat 多个执行多个命令文件，如下：
```linux
$> bat config.conf config1.conf "E:/uisys conf/config2.conf"
---------------------------------------------
  ...
```
### 13. <b>stat</b> 
get application status，for example time and so on.
获取平台的运行状态，用以显示当前平台的起始时间和运行时间。
```linux
$> stat
---------------------------------------------
    Date 2019-07-13 23:27:45
    Now  2019-07-13 23:33:20
```

### 14. <b>exit</b> 
Exit.
退出服务



## 项目参数设置
### 1. <b>ctp</b> 
create project dir.
创建一个UI交互工程目录，ctp 是 <b style='color:#aa0000'>c</b>rea<b style='color:#aa0000'>t</b>e <b style='color:#aa0000'>p</b>roject 的缩写。 
COMMAND: ctp \<Project Path>
命令格式：ctp \<项目路径>  
说明：被创建的工程平台会直接帮您挂在到一个临时服务节点上。
```linux
$> ctp D:\uisys\project01
---------------------------------------------
  create project [D:\uisys\project01].

  The [a0] added successfully.
  The [a0] setted in [D:\uisys\project01].
  The project mount at[a0] server.
```

### 2. <b>stp</b> 
set project dir.  
重新设置一个服务节点的工程目录,stp 是 <b style='color:#aa0000'>s</b>e<b style='color:#aa0000'>t</b> <b style='color:#aa0000'>p</b>roject 的缩写。   
COMMAND: stp \<Service Name> <Project Path>
例如，如果平台已有一个服务节点a0，可以重新让其指向"D:\uisys\project01"的路径。  
```linux
$> stp a0 D:\uisys\project01
---------------------------------------------
  The [a0] setted in [D:\uisys\project01].
```

### 3. <b>ctf</b>
create module file.  
COMMAND: ctf [-Create Method(-h|m|s|r)] \<Service Name> <Project Path>
For Example:ctf test component.Test
ctf test -hr component.Test

### 4. <b>release</b>  
release project.
发布工程为原生工程，以便其他服务器可以使用。  
COMMAND: release \<Service Name> [Project Path]  
```linux
$> release a0 D:\uisys\project-release\
---------------------------------------------
  ----Release Complete----
```
### 5. <b>send</b> 
push data to Service by websocket.  
通过UI-System自建的websocket数据服务，推送数据到WEB客户端。  
COMMAND: send \<Service Name> \<User ID> \<UUID> \<Value>

### 6. <b>lw</b> 
display websocket list of Service  
查看服务节点提供的websocket服务被多少个WEB客户端连接。  
COMMAND: lw \<Service Name> [-h]

### 7. <b>info</b> 
The project infomation  
显示项目信息  
COMMAND: rm \<Service Name>

### 8. <b>set</b> 
Set project attributes.  
设置WEB工程的属性  
COMMAND: set \<Service Name> \<AttributeName> \<Value> [Value...]

### 9. <b>ret</b> 
Remove project attributes.  
COMMAND: ret \<Service Name> \<AttributeName>
