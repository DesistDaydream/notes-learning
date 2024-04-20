---
title: WebDriver
linkTitle: WebDriver
date: 2023-11-23T21:49
weight: 1
---

# 概述

> 参考：
>
> - [MDN，WebDriver](https://developer.mozilla.org/en-US/docs/Web/WebDriver)

WebDriver 是远程控制接口，可以对用户代理进行控制。它提供了一个平台和语言中立的协议，作为浏览器自身进程外的程序远程控制 web 浏览器行为的方法。

WebDriver 符合 [W3C](/docs/Standard/Internet/W3C.md) 标准。

- https://www.w3.org/TR/webdriver1/ 2018 年标准
- https://www.w3.org/TR/webdriver2/ 2023 年草案

# ChromeDriver

> 参考：
>
> - [官网](https://chromedriver.chromium.org/)

`ChromeDriver` 是一个开源项目，由 Chrome 开发团队开发和维护。它是 Chrome 浏览器的一个驱动程序，用于自动化控制和与 Google Chrome 浏览器进行交互。ChromeDriver 允许开发人员使用编程语言（如 Python、Java、C# 等）编写脚本，以控制 Chrome 浏览器的行为。

在 Windows 可以直接下载一个 chromedriver.exe 文件，即可与 Python 库对接，由 Python 控制浏览器。

如果想使用 浏览器自动化程序（e.g. [Selenium](/docs/Web/WebDriver/Selenium.md)） 来模拟用户在浏览器中的各种操作，包括但不限于点击、复制、填写等，那么 ChromeDriver 就是 Chrome 浏览器与 Selenium 进行通信的载体之一

`Selenium` 是一个用于 Web 应用程序测试的工具。它可以直接运行在浏览器中，模拟用户在浏览器中的各种操作，包括但不限于点击、复制、填写等。Selenium 支持市场上所有主流浏览器的自动化，包括 Chrome、Firefox、Safari 等]

Selenium 通过使用 `WebDriver` 支持市场上所有主流浏览器的自动化。WebDriver 是一个 API 和协议，它定义了一个语言中立的接口，用于控制 web 浏览器的行为。每个浏览器都有一个特定的 WebDriver 实现，称为驱动程序。驱动程序是负责委派给浏览器的组件，并处理与 Selenium 和浏览器之间的通信。这种分离是有意识地努力让浏览器供应商为其浏览器的实现负责的一部分。Selenium 在可能的情况下使用这些第三方驱动程序，但是在这些驱动程序不存在的情况下，它也提供了由项目自己维护的驱动程序。

Python 有个 selenium 库用于与 ChromeDriver 交互

# 实现 WebDriver 的程序

## DrissionPage

> 参考：
>
> - [Gitee 项目，g1879/DrissionPage](https://gitee.com/g1879/DrissionPage)
>   - [GitHub 备份](https://github.com/g1879/DrissionPage)

用 requests 做数据采集面对要登录的网站时，要分析数据包、JS 源码，构造复杂的请求，往往还要应付验证码、JS 混淆、签名参数等反爬手段，门槛较高，开发效率不高。 使用浏览器，可以很大程度上绕过这些坑，但浏览器运行效率不高。

因此，DrissionPage 设计初衷，是将它们合而为一，同时实现“写得快”和“跑得快”。能够在不同需要时切换相应模式，并提供一种人性化的使用方法，提高开发和运行效率。  除了合并两者，本库还以网页为单位封装了常用功能，提供非常简便的操作和语句，使用户可减少考虑细节，专注功能实现。 以简单的方式实现强大的功能，使代码更优雅。

以前的版本是对 selenium 进行重新封装实现的。从 3.0 开始，作者另起炉灶，对底层进行了重新开发，摆脱对 selenium 的依赖，增强了功能，提升了运行效率。

## Selenium

[Selenium](/docs/Web/WebDriver/Selenium.md)

## Playwright

Playwright 是一个由 Microsoft 开发的用于浏览器测试和网页抓取的开源自动化库，于 2020 年 1 月 31 日推出

## 其他

[Cypress](https://github.com/cypress-io/cypress) # 多用于测试场景。
