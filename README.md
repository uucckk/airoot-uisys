# AIroot UISYS
[www.airoot.cn](http://www.airoot.cn/)
> A powerful web ui tool, so we call it uisys!

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
## How to run it?
> 1. You need to use uisys, [download](https://github.com/uucckk/airoot-uisys/releases)
> 2. Then save the above file to D:\test\index.ui
> 3. Finally, run uisys.exe, and enter the following command in the terminal:
~~~linux
pub D:\test\ :90
~~~
> OK, open Chrome / Firefox and type http://localhost:90/index.ui

> you can see double **"Hello Baby!"** in your broswer.
## Other examples
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

