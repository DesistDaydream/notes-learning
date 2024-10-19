---
title: "XML"
linkTitle: "XML"
weight: 20
---

# 概述

> 参考：
>
> - [MDN，XML](https://developer.mozilla.org/en-US/docs/Web/XML)
> - [Wiki, XML](https://en.wikipedia.org/wiki/XML)
> - [W3C 官网，XML 标准](https://www.w3.org/TR/xml/)

**Extensible Markup Language(可扩展标记语言，简称 XML)** 是一种用于存储、传输和重建任意数据的标记语言和文件格式，常用来作为配置文件使用。它定义了一组规则，用于以人类可读和机器可读的格式对文档进行编码。万维网联盟 1998 年的 XML 1.0 规范和其他几个相关规范——它们都是免费的开放标准——定义了 XML。

XML 语言由 [DOM](/docs/Web/WebAPIs/DOM.md) 严格序列化，XML 只是一种没有预定义 tags(标签) 的 [HTML](/docs/2.编程/标记语言/HTML.md)（人话: XML 中的 `<dev>、<p>、等等` 标签没有特殊含义）。所有 tag 都像关键字一样，

# XML 标准

XML 使用了与 HTML 相似的术语

- Element(元素)
- Tag(标签)
- Attribute(属性)

整个 XML 是由一个元素的集合体，由根元素开头。通过缩进控制层级，每个层级都表示是上层元素的子元素。

# XPath

> 参考：
>
> - [MDN，XPath](https://developer.mozilla.org/zh-CN/docs/Web/XPath)
> - [Wiki, XPath](https://en.wikipedia.org/wiki/XPath)
> - [菜鸟教程，XPath](https://www.runoob.com/xpath/xpath-tutorial.html)

**XML Path Language(XML 路径语言，简称 XPath)** 是一种表达语言，它使用非 XML 语法来提供一种灵活地定位（指向）[XML](https://developer.mozilla.org/zh-CN/docs/Web/XML) 文档的不同部分的方法。它也可以用于检查文档中某个定位节点是否与某个模式（pattern）匹配。它由万维网联盟 (W3C) 于 1999 年定义，可用于根据 XML 文档的内容计算值（例如字符串、数字或布尔值）。支持 XML 的应用程序（例如 Web 浏览器）和许多编程语言都支持 XPath。

用人话说: XPath 类似于 XML 中的元素的唯一标识符，通过 XPath 可以定位到 XML 中唯一的一个元素。
