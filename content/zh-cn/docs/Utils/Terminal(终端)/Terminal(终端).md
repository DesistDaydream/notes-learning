---
title: Terminal(终端)
linkTitle: Terminal(终端)
date: 2024-05-06T09:04
weight: 1
---

# 概述

> 参考：
>
> -

这里说的 Terminal(终端) 工具是一种可以连接 [Shell](/docs/1.操作系统/Terminal%20与%20Shell/Terminal%20与%20Shell.md) 的图形话工具。

# Terminal 增强工具

## tmux

> 参考：
>
> - [GitHub 项目，tmux/tmux](https://github.com/tmux/tmux)

**terminal multiplexer(终端多路复用，简称 tmux)**，它允许从单个屏幕创建、访问和控制多个终端。 tmux 可能会从屏幕上分离并继续在后台运行，然后再重新连接。

# GUI 终端工具

## Xmanager

https://blog.csdn.net/zhouchen1998/article/details/103424698

Xshell 没有自带的 x11 能力

## SecureCRT

> 参考：
>
> - [官网](https://www.vandyke.com/products/securecrt/index.html)

SecureCRT 是 VanDyke Software 开发的商业终端产品。初始发行于 1995 年 10 月 4 日，没有任何免费版可用，且界面样式非常老旧。

## MobaXterm

> 参考：
>
> - [官网](https://mobaxterm.mobatek.net/)

可以转发 x11

### 关联文件与配置

**MobaXterm.ini** # 所有程序配置、会话信息、etc. 都保存在该文件中。

### 隧道

启动隧道时，若隧道中的 SSH Server 是之前没有登录过或者没有在 User sessions 中创建，则会有新的提示

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/terminal/202403042144907.png)

若 SSH Server 已经在 User seesions 中创建且登录过，则隧道会自动读取这些信息并连接，并不需要再次输入认证信息。

开启动态隧道后，在 Network setting 中设置 Socks5 代理，并指向开启隧道时本地监听的端口即可通过隧道访问 SSH Server 另一侧的设备

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/terminal/202403042218800.png)

## WindTerm

> 参考：
>
> - [GitHub 项目，kingToolbox/WindTerm](https://github.com/kingToolbox/WindTerm)
> - [公众号-k8s 技术圈，C 语言编写的超好用的新一代 SSH 终端 - WindTerm](https://mp.weixin.qq.com/s/2KJi7frtKYExkyBuM5K2hw)

由于作者工作原因，没有更多的时间维护 https://github.com/kingToolbox/WindTerm/issues/1596


快捷键

- Alt + w, Alt + m # 打开/关闭菜单栏

工具 - 同步输入 可以实现多终端同时响应用户输入（Notes: ctrl + l 快捷键无法同时识别）

### WindTerm 关联文件与配置

**${WindTermInstalledDir}/tmp/** # 从 WindTerm 程序的文件管理器直接打开的文件将会在该目录下载作为缓存。当结束 WindTerm 程序时，该目录将会清空

**.wind/** # 数据存储路径。可以指定保存在程序安装目录、用户家目录、自定义目录。

## Tabby

> 参考：
>
> - [GitHub 项目，Eugeny/tabby](https://github.com/Eugeny/tabby)

## Warp

> 参考：
>
> - [GitHub 项目，warpdotdev/Warp](https://github.com/warpdotdev/Warp)
> - [官网](https://www.warp.dev/)

