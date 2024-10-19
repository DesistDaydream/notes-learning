---
title: "Web"
linkTitle: "Web"
date: "2023-06-13T14:56"
weight: 1
---

# 概述

> 参考：
>
> - [Wiki, World_Wide_Web](https://en.wikipedia.org/wiki/World_Wide_Web)

Web 有很多种理解

Nginx 等软件可以提供 Web 服务。

# 学习资料

[MDN，Web 开发技术](https://developer.mozilla.org/en-US/docs/Web)(通常指的是网站首页的 References 标签中的文档)

# Web 开发技术

## Web APIs

不管用什么语言编写 Web 代码，通常都有一些标准的 APIs，有点类似于操作系统的 [POSIX](/docs/1.操作系统/Operating%20system/POSIX.md)。这些 Web API 的标准通常都是由 [W3C](/docs/Standard/Internet/W3C.md)、[IETF](/docs/Standard/Internet/IETF.md)、等多个组织和公司一起制定的，其中 W3C 和 IETF 占了很重要的地位。

详见 [WebAPIs](/docs/Web/WebAPIs/WebAPIs.md)

## 编程语言

[HTML](/docs/2.编程/标记语言/HTML.md)

[CSS](/docs/2.编程/标记语言/CSS.md)

[ECMAScript](/docs/2.编程/高级编程语言/ECMAScript/ECMAScript.md)

[XML](/docs/2.编程/标记语言/XML.md)

## WebAssembly

详见 [WebAssembly](/docs/Web/WebAssembly.md)

## WebDriver

WebDriver 是一种**浏览器自动化**机制，通过模拟真实的人使用浏览器的动作来远程控制浏览器。它被广泛用于网络应用的跨浏览器测试。

详见 [WebDriver](/docs/Web/WebDriver/WebDriver.md)

# Glossary(术语)

**Window(窗口)** # 打开浏览器就相当于打开了一个窗口，这个窗口是用户可以在显示器上直接看到的，可以最小化、最大化、移动、关闭。这是 Windows 系统常用的术语。

**User-Agent(用户代理)** # 在我们发送一个请求时，User-Agent 是可以表示用户的代理方（proxy）。大多数情况下，这个**用户代理都是一个网页浏览器**，不过它也可能是任何东西，比如一个爬取网页来充实、维护搜索引擎索引的机器 [Crawler](/docs/7.信息安全/Crawler/Crawler.md)(爬虫)（其实就是代码写的具有发起 HTTP 请求的程序，毕竟浏览器也是代码写的）。说白了，任何可以发起 HTTP 请求的都可以称为 User-Agent。

- https://en.wikipedia.org/wiki/User_agent
