---
title: "XDG"
weight: 20
---

# 概述

> 参考：
>
> - [freedesktop 规范](https://www.freedesktop.org/wiki/Specifications/)

freedesktop.org 制定互操作性规范，但我们不是官方标准机构。项目不需要实施所有这些规范，也不需要认证。

这些规范许多都在 **X Desktop Group(简称 XDG)** 的旗帜下。（Cross-Desktop Group 代表跨桌面组）

其中一些规范正在（非常）活跃地使用，并且有大量感兴趣的开发人员。其中许多被认为是稳定的，不需要进一步开发，并且可能没有积极的发展。其中一些未被使用或广泛实施。

# 常见变量

> 参考：
>
> - [freedesktop 规范，基础目录 - 变量](https://specifications.freedesktop.org/basedir-spec/latest/ar01s03.html)
>     - https://specifications.freedesktop.org/basedir/latest/#variables

`$XDG_DATA_HOME` # 定义了存放用户专属**数据文件**的基准目录。`默认值: $HOME/.local/share`

`$XDG_CONFIG_HOME` # 定义了存放用户专属**配置文件**的基准目录。`默认值: $HOME/.config`

`$XDG_STATE_HOME` # 定义了存放用户专属**状态文件**的基准目录。`默认值: $HOME/.local/state`

`$XDG_CACHE_HOME` # 定义了存放用户专属**非必要数据文件（i.e. 缓存）**的基准目录。`默认值: $HOME/.cache`
