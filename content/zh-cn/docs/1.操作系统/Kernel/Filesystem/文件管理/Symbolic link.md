---
title: Symbolic link
linkTitle: Symbolic link
date: 2024-04-20T12:27
weight: 3
tags:
  - 符号链接
  - 软链接
---

# 概述

> 参考：
>
> - [Wiki，Symbolic link](https://en.wikipedia.org/wiki/Symbolic_link)

**Symbolic link(符号链接，也称为 symlink 或 soft link(软链接))** 是一种文件，其目的是通过指定路径来指向文件或目录（称为“目标”），这也是一种 [文件管理](/docs/1.操作系统/Kernel/Filesystem/文件管理/文件管理.md) 的方式。

Symbolic link 与目标文件本质上是两个文件，两者的 [Inode](/docs/1.操作系统/Kernel/Filesystem/文件管理/Inode.md) 不一样。
