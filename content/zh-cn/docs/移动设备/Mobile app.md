---
title: Mobile app
linkTitle: Mobile app
date: 2024-01-17T14:16
weight: 20
---

# 概述

> 参考：
> 
> - [Wiki，Mobile_app](https://en.wikipedia.org/wiki/Mobile_app)

**Mobile application(移动应用程序，简称 APP)** 是设计用于在手机、平板电脑或手表等移动设备上运行的计算机程序或软件应用程序。移动应用程序通常与设计用于在台式计算机上运行的桌面应用程序以及在移动 Web 浏览器中运行而不是直接在移动设备上运行的 Web 应用程序形成对比。

> 电脑上的应用程序通常都称为 Application software。


# WeChat

> 参考：
> 
> - [微信开发者官方文档](https://developers.weixin.qq.com/doc/)

[微信浏览器：简单三步打开调试工具](https://www.cnblogs.com/conne/p/15884968.html)

## 微信开发者工具

> 参考：
> 
> - [官方文档](https://developers.weixin.qq.com/miniprogram/dev/devtools/devtools.html)

从 https://developers.weixin.qq.com/miniprogram/dev/devtools/download.html 下载

这本质就是一个类似 VSCode 的编辑器。

## 关联文件与配置

`WeChat Files\Applet\` # 该目录为小程序所在文件夹。每个小程序文件都是一个独立的文件夹，以 wx 开头，像 `wx64479c83c7630409` 这样

- 想要找到对应的小程序，可以把所有 wx 开头的文件夹都删除，然后打开小程序，就会生成一个信息的 wx 开头的文件夹。

## 小程序

调试 微信小程序大体有两种办法，官方的和非官方的

**官方的 devtools** - 必须要有小程序源码才可以

- https://github.com/Tencent/vConsole # 一个轻量、可拓展、针对手机网页的前端开发者调试面板。现在 vConsole 是微信小程序的官方调试工具。
- https://github.com/weimobGroup/WeConsole # 替代了 vConsole。功能全面、界面与体验对标 Chrome devtools 的可定制化的小程序开发调试面板。个人开发，后被腾讯收编，项目转移到 微盟技术中心
  - https://github.com/weimob-tech/WeConsole # 收编了 weimobGroup/WeConsole 后的项目所在仓库，位于 weimob-tech(微盟技术中心) 组织下
- 除了在小程序中内嵌 devtools 并在实体机器或 PC 微信中进行调试外，还可以使用微信官方的 [微信开发者工具](https://developers.weixin.qq.com/miniprogram/dev/devtools/devtools.html) 一边写代码一边调试。

**非官方的 devtools** - 这类项目通常需要使用类似 hook 的方式来为本身没有 devtools 的小程序添加 devtools。

- https://github.com/x0tools/WeChatOpenDevTools # 通过 Frida Hook 注入开启小程序 DeBug，唐志远作者。因部分原因已删除库
  - https://github.com/shuaibibobo/WeChatOpenDevTools # WeChatOpenDevTool 原仓库的 fork 备份
  - https://github.com/JaveleyQAQ/WeChatOpenDevTools-Python # # WeChatOpenDevTool 的 Python 实现
  - https://github.com/yuweiping/WeChatOpenDevTools # 另一个 Python 实现