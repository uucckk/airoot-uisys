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
> If you want the **MyBox** module to be a file, you can save it, for example:

> Create the file **MyBox.ui**, and the code is as follows:
~~~html
<!-- module MyBox code in D:\test\MyBox.ui -->
<div>Hello Baby!</div>
~~~
> Then create the file **Index.ui** with the following code:
~~~html
<!-- module Index code in D:\test\Index.ui -->
<@pub/>
<div>
  <MyBox/>
  <MyBox/>
</div>
~~~
> open the browser to see the results,
> You can still see double "**Hello Baby!**" words.

## use package
> If you want to put files in a folder, you can do this:

> For example, take the **D:\test\MyBox.ui**  move to **D:\test\mymod\mybox.ui**.

> Then **Index.ui** import **MyBox.ui** should be changed to the following code:
~~~html
<!-- module Index code in D:\test\Index.ui -->
<@pub/>
<div>
  <mymod.MyBox/>
  <mymod.MyBox/>
</div>
~~~
> Or it can be imported globally：
~~~html
<!-- module Index code in D:\test\Index.ui -->
<@pub/>
<@import value="mymod.MyBox" />
<div>
  <MyBox/>
  <MyBox/>
</div>
~~~
> open the browser to see the results,
> You can still see double "**Hello Baby!**" words.
