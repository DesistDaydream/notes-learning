---
title: DOM
---

# 概述

> 参考：
> 
> - [MDN 官方文档，参考-WebAPIs-DOM](https://developer.mozilla.org/en-US/docs/Web/API/Document_Object_Model)

**Document Ojbect Model(文档对象模型，简称 DOM)** 是 Web 文档(HTML 和 XML)的编程接口(通常描述为 WebAPI 中的 DOM 接口)。应用程序可以通过该接口更改 Web 文档的 结构、样式、内容 等。DOM 将 Web 文档抽象为 **Node(节点)** 和 **Ojbect(对象，包含属性和方法的对象) **组成的结构集合。

一个 Web 页面即是一个文档，这个文档可以在浏览器中作为 HTML 源码展示出来。DOM 则可以将文档表示为另一种形式，以便 JavaScript 等编程语言可以对其进行修改。

比如：

```javascript
// 我们通过 document.querySelectorAll() 获取 Web 文档中所有 <p> 元素的列表
// 将所有 <p> 元素实例化为 paragraphs 变量
var paragraphs = document.querySelectorAll("p")
// 之后，通过代码对 paragraphs 的所有操作都会直接反应到前端 Web 页面上
// 这里表示将将会弹出提示框，并将其中第一个 <p> 元素的名称显示在提示框中
alert(paragraphs[0].nodeName)
```

从上面的示例中可以看到，JavaScript 中使用 `**document 类型的对象**`表示 Web 文档本身；document 对象里包含了非常多的方法来控制 Web 文档中的元素，在 [MDN 官方文档，WebAPIs-Document](https://developer.mozilla.org/en-US/docs/Web/API/Document) 中可以看到所有 document 对象下的属性、方法、事件。示例中的 [querySelectorAll()](https://developer.mozilla.org/zh-CN/docs/Web/API/Document/querySelectorAll) 方法将会返回匹配到的元素列表。

DOM 本身并不是一个编程语言，可以说是一种 规范、模型、接口；DOM 可以用任何语言实现，DOM 被内嵌在浏览器中，各种编程语言可以自己实现 DOM 库以便在浏览器中调用 DOM。

编程语言之于 DOM，有点类似于 runc 等运行时之于 OCI

我们甚至可以在 Python 中使用 DOM 来控制 Web 文档

```python
# Python DOM example
import xml.dom.minidom as m
doc = m.parse(r"C:\Projects\Py\chap1.xml")
doc.nodeName # DOM property of document object
p_list = doc.getElementsByTagName("para")
```
