# AIroot UISYS
> A powerful webui tool, so we call it uisys!
[http://www.airoot.cn/](http://www.airoot.cn/)
## For example
~~~html
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
> 你需要使用uisys,[](download here)
