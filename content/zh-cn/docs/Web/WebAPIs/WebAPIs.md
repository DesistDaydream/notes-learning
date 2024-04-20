---
title: WebAPIs
linkTitle: WebAPIs
date: 2024-04-15T20:08
weight: 1
---

# 概述

> 参考：
>
> - [MDN，参考 - Web API](https://developer.mozilla.org/en-US/docs/Web/API)

在编写 Web 代码时，有许多 Web APIs 可供调用。下面是开发 Web 应用程序或网站时可能使用的所有 API 和接口（对象类型）的列表。

Web APIs 主要用于 JavaScript，但也可能有例外。

# Window

https://developer.mozilla.org/zh-CN/docs/Web/API/Window

Window 接口是各种函数、对象、等等的家。`window` 对象表示一个包含 [DOM](/docs/Web/WebAPIs/DOM.md) 的窗口（通常来说都是具有浏览器功能的窗口），document 属性指向窗口中载入的 DOM 文档。

Navigator 接口表示用户代理的状态和标识。它允许脚本查询它和注册自己进行一些活动。用白话说：这里面包含了浏览器相关的信息。

- appVersion # 浏览器的版本号
- appName # 浏览器的名称
- language # 浏览器使用的语言
- platform # 浏览器所使用的平台
- userAgent # 浏览器的 user-agent 信息。常用来区分浏览网站的人使用了什么设备
- webdriver # 当前窗口是否使用了 [WebDriver](/docs/Web/WebDriver/WebDriver.md)。在爬虫技术中，移除这个是很重要的一点避免被网站识别成 WebDriver。
