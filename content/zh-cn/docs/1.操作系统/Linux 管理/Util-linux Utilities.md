---
title: Util-linux Utilities
linkTitle: Util-linux Utilities
date: 2024-03-18T00:01
weight: 2
---

# 概述

> 参考：
>
> - [GitHub 项目，util-linux/util-linux](https://github.com/util-linux/util-linux)
> - [Wiki, Util-linux](https://en.wikipedia.org/wiki/Util-linux)

util-linux 是由 Linux 内核组织分发的标准软件包，用作 Linux 操作系统的一部分。一个分支 util-linux-ng（ng 的意思是“下一代”）是在开发停滞时创建的，但截至 2011 年 1 月，它已重命名为 util-linux，并且是该软件包的正式版本。

可以在 [这里](https://en.wikipedia.org/wiki/Util-linux#Included) 找到 Util-linux 包中通常包含的所有程序。这些程序可以分为几大类

- Namespace 管理，包括 unshare, nsenter, lsns, etc. 。详见 [容器运行时管理](/docs/10.云原生/Containerization%20implementation/容器管理/容器运行时管理/容器运行时管理.md)
- etc.

还有一部分已经弃用的程序可以在[这里](https://en.wikipedia.org/wiki/Util-linux#Removed)找到列表
