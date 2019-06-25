# 指令
## 定义指令
```javascript
class Red{
    function init(e){
        e.target.dom.style = "color:#ff0000";
    }
}
```

## 使用指令
```html
<div -red>
    Hello World!
</div>
```
## 指令接收参数值
```javascript
class Color{
    function init(e){
        e.target.dom.style = e.value;
    }
}
```

```html
<div -color="#FFaa00">
    Hello World!
</div>
```
<hr/>
下一页