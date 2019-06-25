
对WEB的开发，了解一点的人可能觉得那是前端内容，懂一些的人感觉Web也是一门复杂的学科，但是对于初学者来说，Web本应该没那么复杂，因为单独从面向对象角度来讲，Web才是刚刚起步。

开发过Java，C# WPF的朋友们，在看HTML代码的时候，感觉代码不成体系，没错实际上，确实这样，Web由于语言平台的限制，在使用上很难实现传统的面向对象编程。比如说现在还没有真正一款Web框架可以实现WPF那样的MVVM编程。有的同学说react、angularJS、Vue不都实现了MVVM吗。但实际上，这些框架或者平台只是实现了部分的逻辑，如果从代码解耦、易用、写法上来讲，那是不能比拟的。

即便是这一样，Web的发展如今还是空前繁荣的，不得不说Web的发展势头非常强悍。为什么会这样呢，我分析有以下几点原因：

1. Web开发触手可得、学习成本低、IDE丰富、产品易于实现。
2. Web跨平台，PC、平板、手机等主要平台都支持。
3. Web调试方便、在浏览器端、移动端都有简单易用的调试工具。
4. Web是动态解析语言、修正错误、维护成本都很低。

如果，Web能够占用更小的内存、速度和原生一样好，那我估计Web就一统天下了。我记得之前同事调试android程序，修改一行就看下效果，改来改去一上午就过去了，他们看到我直接就在浏览器里面改，然后直接生效，感觉羡慕不已。

但是Web开发的确定是什么呢？

1. 缺少完备的面向对象、模块化、插件化支持，项目越庞大、约难以维护。
2. 由于编写的不当会造成内存不会收，JS、CSS之前产生冲突。
3. 业内流行的框架，引入了更多的思想来解决Web开发问题，但是随着项目累积，框架学习成本增加，还需要填更多坑。

很多人都不了解AIroot-JUS所提倡的”友好、易用、稳定、强大“到底指的是什么?

JUS就是为了解决上面三点问题出现的。

我们来看看，JUS编写一个模块。
编写一个Red.html页面。
```html
<span style="color:#F00">
  <@content/>
</span>
```
编写一个Main.html主页面
```html
<div>
  <Red>I'm Red.</Red>
</div>
```
看一下JUS的模板使用
```html
<div>Hello {name}!</div>
<script>
  var name = "World";
  dom.dataContext = @this;
</script>
```

当然，你也可以这样写：
```html
<div>Hello {name}!</div>
<script>
  dom.dataContext = {name:"World"};
</script>
```
## 看下JUS的数据绑定
```html
<div>
  <div>请输入：<input type="text" bind="name" /></div>
  您输入的是 {name}!
</div>
<script>
  dom.dataContext;
</script>
```

## JUS用语言创建对象

  创建一个模块，命名类Label.html
```html
<h1>我是一个标签</h1>
```
```html
<div>
  <div id="list"></div>
</div>
<script>
  import Label;
  function init(){
    for(var i = 0;i<10;i++){
      list.append(new Label());
    }
  }
</script>
```

## 判断module被挂在到节点上
```html
<div>
  <div id="list"></div>
</div>
<script>
  import Label;
  function init(){
    for(var i = 0;i<10;i++){
      list.append(new Label());
    }
  }

  function mount(){
    alert("我要被挂了");
  }

  function mounted(){
    alert("我被挂在到节点上了");
  }

  function remove(){
    alert("我被移除了");
  }

  function finalize(){
    alert("我进入垃圾回收状态了");
  }
</script>
```
## 为模块对象添加 setter,getter 属性
```html
<div>
  你好 {name} ！
</div>
<script>
  function init(){
    dom.dataContext = @this;
  }

  set name(value){
    dom.dataContext.name = value;
  }

  get name(){
    return dom.dataContext.name;
  }
</script>
```
## 多个模块绑定一个数据源

### 模块01 Li.html
```html
<div>
  你好 {name} ！我是李雷！
</div>
<script>
  set data(value){
    dom.dataContext = value;
  }
</script>
```
### 模块02 Han.html
```html
<div>
  你好 {name} ！我是韩梅梅！
</div>
<script>
  set data(value){
    dom.dataContext = value;
  }
</script>
```
### 主模块 Main.html
```html
<div>
  <Li id="li" />
  <Han id="han" />
</div>
<script>
  var data = {name:"Sunxy"};
  function init(value){
    li.data = data;
    han.data = data;
    //两秒后改为中文名
    setInterval(function(){
      data.name = "孙晓玉";
      update(data);//更新数据
    },2000);
  }
</script>
```
## 为单个微模块添加样式
```html
<style>
  body{
    color:#456789;
    font-weight:bold;
  }
</style>
<div>
  Hello World!
</div>
```
## 为本类全局模块添加样式
```html
<css>
  body{
    color:#456789;
    font-weight:bold;
  }
</css>
<div>
  Hello World!
</div>
```

## 用Jus实现纯MVVM逻辑

### 定义一个Module
```javascript
class Module{
  public var onComplete = null;
  function init(){
    //TODO Sth. get data;
    if(onComplete){
      onComplete(data);
    }
  }
  
}
```

### 定义一个ViewModule
```javascript
class ViewModule{
  public var name = null;
  public var sex = null;
  public function showName(){
    alert(name);
  }
  public set data(value){
    name = value.name;
    sex = value.sex;
  }
}
```

### 定义一个View,名字叫做：View.html
```html
<div>
  data:{value}
</div>
<script>
  set data(value){
    dom.dataContext = value;
  }
</script>
```

### 现在开始组装MVVM，我们建立一个Main.html
```html
<div>
  <View id="view" />
</div>
<script>
  import Module;
  import ViewModule;
  function init(){
    var module = new Module();
    module.onComplete = function(data){
      var vm = new ViewModule();
      vm.data = JSON.parse(data);
      view.data = vm;
    }
  }
</script>
```

JUS 单点测试用例
做