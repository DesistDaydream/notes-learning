---
title: Samba 配置详解
---

# 概述

> 参考：
> [Manual(手册)](https://www.systutorials.com/docs/linux/man/5-smb.conf/)

Samba 通过配置文件来配置运行时行为。

## 文件格式

Samba 的配置文件格式类似于 [INI](/docs/2.编程/无法分类的语言/INI.md) 格式的配置，由 **Sections(部分)** 和 **Parameters(参数)**组成。参数是以 `=` 分割的 **键/值对**。

配置文件中的每个 Sections(部分) 用于描述共享资源(Note：global 部分除外)。Section 的名称就是共享资源的名称，Section 中的 Parameters 定义共享资源的属性。

> 用白话说就是，加入有一个名为 share_dir 的 Section，则 Windows 中兴就会看到一个名为 share_dir 的目录

有 3 个特殊的 Sections，global、homes、printers，这 3 个特殊部分，是配置文件预留的，用于定义一些通用的运行时配置。

> 注意：其他的 Sections 不能以这 3 个特殊的 Sections 的名字命名。

# 配置文件详解

## global 部分

## homes 部分

## printers 部分

## 其他部分
