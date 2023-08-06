---
title: "tmpfs"
linkTitle: "tmpfs"
date: "2023-08-06T11:06"
weight: 20
---

# 概述

> 参考：
> 
> - [Kernel 文档，文件系统-tmpfs](https://www.kernel.org/doc/html/latest/filesystems/tmpfs.html)

Tmpfs 是一个将所有文件保存在虚拟内存中的文件系统。tmpfs 中的所有内容都是 **temporary(临时)** 的，因为不会在硬盘上创建任何文件。如果卸载 tmpfs 实例，其中存储的所有内容都会丢失。

tmpfs 将所有内容放入内核内部缓存中，并增长和收缩以容纳其包含的文件，并且如果为 tmpfs 挂载启用了交换，则能够将不需要的页面交换到交换空间。 tmpfs 还支持 THP。


# 分类

#文件系统