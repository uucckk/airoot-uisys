# AIroot UISYS
> A powerful webui tool, so we call it uisys!
[http://www.airoot.cn/](http://www.airoot.cn/)
## For example
~~~html
<@pub/>
<!-- define a module -->
<@define name="MyBox">
  <div>Hello Baby!</div>
</@define>

<!-- you code here. -->
<div>
  <MyBox/>
  <MyBox/>
</div>
~~~
> You can also use JavaScript as follows:
~~~html
<@pub/>
<!-- define a module -->
<@define name="MyBox">
  <div>Hello Baby!</div>
</@define>

<!-- you code here. -->
<div>
  <!-- dom area -->
</div>
<script>
  function init(){
    var box = new MyBox();
    dom.appendChild(box);
  }
</script>
~~~
> mybe ...
~~~html
<@pub/>
<!-- define a module -->
<@define name="MyBox">
  <div>Hello Baby!</div>
</@define>

<!-- you code here. -->
<div id="ct"></div>
<script>
  function init(){
    var box = new MyBox();
    #ct.appendChild(box);
  }
</script>
~~~
## How to run it?
> 1. You need to use uisys, [https://github.com/uuckk/airoot-uisys/releases] (download here)
> 2. Then save the above file to D:\test\index.ui
> 3. Finally, run uisys.exe, and enter the following command in the terminal:
~~~linux
pub D:\test\ :90
~~~
> OK, open Chrome / Firefox and type http://localhost:90/index.ui
> you can se "Hello Baby!" in your broswer.
