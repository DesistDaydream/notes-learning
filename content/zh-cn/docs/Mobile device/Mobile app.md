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

`WeChat Files\Applet\` # 该目录是小程序所在文件夹。每个小程序文件都是一个独立的文件夹，以 wx 开头，像 `wx64479c83c7630409` 这样

- 想要找到对应的小程序，可以把所有 wx 开头的文件夹都删除，然后打开小程序，就会生成一个信息的 wx 开头的文件夹。

## 小程序

调试 微信小程序大体有两种办法，官方的和非官方的

**直接使用官方的 devtools** - 必须要有小程序源码才可以

- https://github.com/Tencent/vConsole # 一个轻量、可拓展、针对手机网页的前端开发者调试面板。现在 vConsole 是微信小程序的官方调试工具。
- https://github.com/weimobGroup/WeConsole # 替代了 vConsole。功能全面、界面与体验对标 Chrome devtools 的可定制化的小程序开发调试面板。个人开发，后被腾讯收编，项目转移到 微盟技术中心
  - https://github.com/weimob-tech/WeConsole # 收编了 weimobGroup/WeConsole 后的项目所在仓库，位于 weimob-tech(微盟技术中心) 组织下
- 除了在小程序中内嵌 devtools 并在实体机器或 PC 微信中进行调试外，还可以使用微信官方的 [微信开发者工具](https://developers.weixin.qq.com/miniprogram/dev/devtools/devtools.html) 一边写代码一边调试。

**间接使用官方的 devtools** - 这类项目通常需要使用类似 hook、拦截等方式。TODO: 到底是使用了原本就存在的官方 devtools 还是

- https://github.com/x0tools/WeChatOpenDevTools # 通过 Frida Hook 注入开启小程序 DeBug。作者唐志远，因部分原因已删除库
  - https://github.com/shuaibibobo/WeChatOpenDevTools # WeChatOpenDevTool 原仓库的 fork 备份
  - https://github.com/JaveleyQAQ/WeChatOpenDevTools-Python # # WeChatOpenDevTool 的 Python 实现
  - https://github.com/yuweiping/WeChatOpenDevTools # 另一个 Python 实现

# 手机上的APP都是用什么编程语言写的？

https://zhuanlan.zhihu.com/p/444481759

## **第一类：针对单一APP开发的语言，即开发一套代码只能运行在一个平台上。** 

### 1、开发Android的：Java和Kotlin

Kotlin是一种在Java虚拟机上运行的静态类型编程语言，被称之为Android世界的Swift。Kotlin可以编译成Java字节码，也可以编译成JavaScript，方便在没有 JVM 的设备上运行。运行效率提高很多，并且语法更加简洁好用。

如果是与系统底层进行交互则需要使用JNI技术，通过和C或者C++结合实现相应的业务逻辑，比如美颜或者直播类型的APP。直播类型的APP采用的ffmpeg技术，其中ffmpeg就是用C语言实现的。

![](https://pic2.zhimg.com/v2-b6697b88d959003af2515ab4bc2092e9_b.jpg)

### 2、开发IOS的：Swift和Object-C**

Swift 结合了 C 和 Objective-C 的优点并且不受 C 兼容性的限制。

![](https://pic1.zhimg.com/v2-c988362c2aedf592ad366a335a53ec38_b.jpg)



## 第二类：可以针对多个APP端的编程语言，即只需开发出一套代码，就可在多个平台上运行。

### 1、第一种是Flutter技术

基于Dart语言，比如现在闲鱼APP就是基于flutter开发的。并且有着多年经验，而且闲鱼免费开源了框架。

Flutter是谷歌的移动UI框架，可以快速在iOS和Android上构建高质量的原生用户界面。Flutter可以与现有的代码一起工作。在全世界，Flutter正在被越来越多的开发者和组织使用，并且Flutter是完全免费、开源的。

Flutter的热重载可帮助您快速地进行测试、构建UI、添加功能并更快地修复错误。在iOS和Android模拟器或真机上可以在亚秒内重载，并且不会丢失状态。

尤其是在UI上使用Flutter的现代、响应式框架，和一系列基础widget，轻松构建您的用户界面。使用功能强大且灵活的API（针对2D、动画、手势、效果等）解决艰难的UI挑战。

![](https://pic4.zhimg.com/v2-3363e67d08bc43662392a39a0d18ca1b_b.jpg)

### 2、第二种是Uni-app框架

基于 Vue.js。俗称一套代码编到8个平台上。

uni-app是一个使用Vue.js开发所有前端应用的框架，开发者编写一套代码，可发布到iOS、Android、H5、以及各种小程序：微信/支付宝/百度/头条/QQ/钉钉等多个平台。

![](https://pic1.zhimg.com/v2-b915ba2b373b76680ef60ab18489db30_b.jpg)

## 小结：

第二类开发起来更显效率，节约很多时间成本。当然，如果加上大型APP这个限定条件的话，就不能用单一的某一种编程语言来说了，一般情况下大型APP必然会用到的编程语言有三种：

**第一种：** 平台原生推荐语言，如：Android平台的Java和Kotlin，ios平台的object-c和Swift（swift版本之间的差异比较大，要学习的话还是要注意版本选择）。这部分一般用于处理核心业务、权限请求以及高性能要求页面。

**第二种：** web语言，这里主要指的是h5相关的技术栈。这部分主要处理非核心业务逻辑，以及需要动态更新的页面。

**第三种：** NDK相关，这里一般用到的就是C、C++。大部分写业务逻辑的同学用的比较少，不过大型APP一般还是会用到一些。主要应用于安全性要求高，高性能算法以及跨平台算法实现。

素材源于:文章来源，C语言与程序设计；直接来源:嵌入式ARM
