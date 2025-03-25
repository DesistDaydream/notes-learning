---
title: DevTools
linkTitle: DevTools
weight: 20
---

# 概述

> 参考：
>
> - [Chrome 开发者工具官方文档](https://developer.chrome.com/docs/devtools/)
> - [Wiki, Web_development_tools](https://en.wikipedia.org/wiki/Web_development_tools)

Chrome DevTools 是一组内置于 Chrome 浏览器中的 Web 开发者工具（通常使用 F12 快捷键打开；`Ctrl + Shift +i` 也是常见的默认快捷键）。

- **元素** # 当前网页被渲染的数据
- **控制台** # 可以输入 JS 代码
- **源代码/来源** # HTML、JS、CSS、字体等静态资源文件
- **网络** # 数据传输内容
- .....
- **应用** #

# 元素

使用 `Ctrl+Shift+C` 快捷键可以选定页面上的元素以定位到 HTML 中的位置。

# 网络

# 控制台

在控制台中，通过查看日志消息和运行 JS 代码以控制页面的控制。

- 可以直接输出当前页面 js 代码中的各种变量、函数等等
- 等等

# 源代码

> 参考：
>
> - [官方文档，源代码](https://developer.chrome.com/docs/devtools/sources/)

在源代码中可以检查静态资源、打断点、等等。

当浏览器的请求被断点暂停时，我们可以将鼠标悬停后，看到 FunctionLocation 关键字，从这里追踪函数的位置。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/devtools/202311202258018.png)

> 在控制台输出该函数，也可以通过点击跳转到函数位置。

## 代码段

https://developer.chrome.com/docs/devtools/javascript/snippets?hl=zh-cn

若我们需要在[控制台](#控制台)反复运行相同的代码，可以将代码保存为**代码段**，以便随时使用

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/devtools/202312090821304.png)

最常见的用法是替换 JS 中的函数以 Debug，比如下面这个，可以通过 `JSON.stringify = function (params) {}` 和 `JSON.parse = function (params) {}` 修改原始 `JSON.stringify` 和 `JSON.parse` 方法的逻辑：在其中加入文本输出和 debug 暂停能力。

```js
(function () {
    var my_stringify = JSON.stringify;
    JSON.stringify = function (params) {
        console.log("Hook 字符串化", params);
        debugger
        return my_stringify(params);
    }
    var my_parse = JSON.parse;
    JSON.parse = function (params) {
        console.log("Hook 解析", params);
        debugger
        return my_parse(params);
    }
})();
```

这时，凡是页面的 JS 代码中调用了 JSON.stringify 和 JSON.parse 这俩方法的地方，都会输出参数，并被 debugger 关键字暂停以进行断点检查。然后可以在右侧 *调用堆栈* 中点击直接跳转到网页的代码中，对应的位置。

Notes: 仅对已经加载完成的页面有效，i.e. 发起 Fetch/XHR 之类的请求时才会被拦截。

## 断点

可以通过多种方式添加断点

- 断点 # 与 IDE 类型，直接点击代码的行号能添加。
- XHR/提取 断点
- DOM 断点
- 事件监听器断点
- CSP 违规断点

### 最佳实践

从 `网络 - XHR/Fetch` 找到的 请求，可以提取出其中部分路径当做要断点的检查点。比如 Grafana 点击首页是，会向 `https://DOMAIN/api/dashboards/home` 发起请求，那么在断点处添加 `/api/dashboards/home`

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/devtools/202311201206181.png)

此时访问首页时，即可让页面在发起这个请求时暂停并定位到源码位置

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/devtools/202311201207139.png)

## Scope(作用域)

查看当前代码执行的地方都有哪些函数、变量及其对应的值。

## Call Stack(调用堆栈)

https://developer.chrome.com/docs/devtools/javascript/reference#call-stack

**Call Stack(调用堆栈)** 可以查看当前断点位置的调用堆栈，i.e. 从访问网页开始到断点停止时，所有的函数调用，每行都是一个调用，也可以称为 Frame(帧)

- 左侧是函数名称。若函数是匿名的则显示`（匿名）`，而且浏览器中的代码基本都是打包过甚至混淆过的，函数名不一定与源代码显示一样。
- 右侧是函数所在文件及所在文件的行号。

整个 Call Stack **由下到上**的顺序，是浏览器从访问网站开始，到断点位置为止的所有函数调用。

> 若想从断点位置反查代码，只需要从上往下逐一查看，即可找到想要的网页中某个数据生成的最初位置。

以百度为例

![image.png|1000](https://notes-learning.oss-cn-beijing.aliyuncs.com/browser_devtools/202401182141362.png)

当前断点的函数名 `XMLHttpRequest.send` 的，那么执行到 `u.apply()` 代码之前的函数就是第二行的 `send`；在执行到 `send` 代码之前的函数就是第三行的 ajax；以此类推。

我们可以**右键点击**函数名，选择 **重启帧**，这样可以让代码直接重新运行并停顿到该函数位置。这样就无需重启整个调试流程（刷新网页）。

假如我们右键点击 send，并选择重启帧，然后再点击一下执行下一步，可以看到上述关于 Call Stack 笔记的具体应用效果。在 `send` 执行完成后，会执行到 `XMLHttpRequest.send`

![recording.gif|1000](https://notes-learning.oss-cn-beijing.aliyuncs.com/browser_devtools/restart_frame_1.gif)

DevTools 源码中的调用栈，用来调试找代码是非常方便且好用的！
