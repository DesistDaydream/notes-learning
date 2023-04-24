---
title: "XDG"
linkTitle: "XDG"
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

https://specifications.freedesktop.org/basedir-spec/latest/ar01s03.html

`XDG_CACHE_HOME` 定义了应该存储用户特定的非必要数据文件的基本目录。如果 `$XDG_CACHE_HOME` 未设置或为空，则应使用等于 `$HOME/.cache` 的默认值。这是一个 Linux 和 Unix 操作系统环境变量，Windows 系统中并没有这个环境变量。
