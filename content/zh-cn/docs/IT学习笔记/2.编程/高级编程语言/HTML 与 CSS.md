---
title: HTML 与 CSS
---

# 概述

> 参考：
> - [Wiki，HTML](https://en.wikipedia.org/wiki/HTML)

**HyperText Markup Lanugage(超文本标记语言，简称 HTML)** 是构成 Web 世界的一砖一瓦。它定义了网页内容的含义和结构。除 HTML 以外的其它技术则通常用来描述一个网页的表现与展示效果（如 [CSS](https://developer.mozilla.org/zh-CN/docs/Web/CSS)），或功能与行为（如 [JavaScript](https://developer.mozilla.org/zh-CN/docs/Web/JavaScript)）。
“超文本”（hypertext）是指连接单个网站内或多个网站间的网页的链接。链接是网络的一个基本方面。只要将内容上传到互联网，并将其与他人创建的页面相链接，你就成为了万维网的积极参与者。
HTML 使用“标记”（markup）来注明文本、图片和其他内容，以便于在 Web 浏览器中显示。HTML 标记包含一些特殊“元素”如 [<head>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/head)、[<title>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/title)、[<body>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/body)、[<header>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/header)、[<footer>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/footer)、[<article>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/article)、[<section>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/section)、[<p>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/p)、[<div>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/div)、[<span>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/span)、[<img>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/img)、[<aside>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/aside)、[<audio>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/audio)、[<canvas>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/canvas)、[<datalist>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/datalist)、[<details>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/details)、[<embed>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/embed)、[<nav>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/nav)、[<output>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/output)、[<progress>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/progress)、[<video>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/video)、[<ul>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/ul)、[<ol>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/ol)、[<li>](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/li) 等等。
HTML **Element(元素)** 通过 **Tag(标签)** 将文本从文档中引出，标签由在 `<` 和 `>` 中包裹的元素名组成，HTML 标签里的元素名不区分大小写。也就是说，它们可以用大写，小写或混合形式书写。例如，<title> 标签可以写成 <Title>，<TITLE> 或以任何其他方式。然而，习惯上与实践上都推荐将标签名全部小写。

# 学习资料

[MDN 官方文档，Web 开发技术](https://developer.mozilla.org/en-US/docs/Web)(通常指的是网站首页的 References 标签中的文档)

- [HTML](https://developer.mozilla.org/en-US/docs/Web/HTML)

[W3schools，HTML 教程](https://www.w3schools.com/html/default.asp)
[网道，HTML](https://wangdoc.com/html/)
[菜鸟教程，HTML](https://www.runoob.com/html/html-tutorial.html)

## 各种 HTML、CSS、Vue 等代码示例

- [菜鸟教程](https://www.runoob.com/)里对应代码的的示例非常多
- <https://gitee.com/wyanhui02/html_css_demo>

# Hello World

# HTML 语言关键字

> 参考：
> - [MDN 学习 Web 开发，HTML-HTML 基础](https://developer.mozilla.org/zh-CN/docs/Learn/Getting_started_with_the_web/HTML_basics)
> - [MDN Web 开发技术，HTML-参考-HTML 元素](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/uw7agz/1666090441255-222a4602-e9ed-43f6-979c-944836075c4f.png)
HTML 的关键字又称为 **Tag(标签)**，有时候又称为 **Element(元素)**。

- Tag 是由尖括号 `<KeyWord>` 包围的关键字，比如 `<html>`；Tag 总是成对出现，比如 `<html></html>`，称为 **OpeningTag**(**起始标签) **和 **ClosingTag**(**结束标签)**(有时候也成为开放标签和闭合标签)

起始标签 与 结束标签 中包含的就是 **Content(内容)，**也可以称为元素的内容。**起始标签**、**结束标签**以及**内容** 组合在一起，称之为 **Element(元素)。**Element 严格来讲不能称为关键字，但是有时候人们会经常把两个词混用。
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/uw7agz/1666090500626-571d651a-a81a-404d-963c-2a414eb6466b.png)
可以在**起始标签**中定义 **Attribute(属性)** 以改变元素的表现形式和行为

- 在属性与元素名称（或上一个属性，如果有超过一个属性的话）之间的空格符。
- 属性的名称，并接上一个等号。
- 由引号所包围的属性值。

# HTML 语言规范
