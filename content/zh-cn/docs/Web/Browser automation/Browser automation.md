---
title: Browser automation
linkTitle: Browser automation
weight: 1
---

# 概述

> 参考：
>
> -

实现浏览器自动化机制的 Awesome 项目

- https://github.com/angrykoala/awesome-browser-automation

## WebDriver

WebDriver 是一种**浏览器自动化**机制，通过模拟真实的人使用浏览器的动作来远程控制浏览器。它被广泛用于网络应用的跨浏览器测试。

详见 [WebDriver](/docs/Web/Browser%20automation/WebDriver.md)

## 其他

[Cypress](https://github.com/cypress-io/cypress) # 多用于测试场景。

### Playwright

https://github.com/microsoft/playwright

Playwright 是一个由 Microsoft 开发的用于浏览器测试和网页抓取的开源自动化库，于 2020 年 1 月 31 日推出

基于 [DevTools](/docs/Web/Browser/DevTools.md) 的协议

![https://playwright.dev/python/docs/selenium-grid#introduction|800](https://notes-learning.oss-cn-beijing.aliyuncs.com/browser_automation/20250709102946555.png)

使用 [connectOverCDP](https://playwright.dev/docs/api/class-browsertype#browser-type-connect-over-cdp) 方法通过 Chrome DevTools Protocol(CDP) 连接
