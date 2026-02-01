---
title: ISO文件与可启动U盘
linkTitle: ISO文件与可启动U盘
weight: 61
---

# 概述

> 参考：
>
> - [Wiki, ISO 9660](https://en.wikipedia.org/wiki/ISO_9660)

ISO 9660（也称为 ECMA-119） 是光盘介质的文件系统。.iso 文件格式即是 ISO 9660 标准

由于光盘的逐渐淘汰，当代直接把这种文件系统拷贝到 U 盘中

# Rufus

> 参考：
>
> - [GitHub 项目，pbatard/rufus](https://github.com/pbatard/rufus)
> - [官网](https://rufus.ie/)

引导类型选择

- 本地
- 下载

> Tip: 将 Win11 制作到 U 盘之前，可以选择去掉 TPM2.0 检查相关限制。

# UltraISO

> 参考：
>
> - [官网](https://www.ultraiso.com/)

UltraISO 是一个运行在 Microsoft Windows 平台上的用来创建、修改和转换 ISO 文件的软件。

官方似乎没有 Portable(便携版)~~ o(╯□╰)o

## 将 ISO 文件写到 U 盘

- 插入 U 盘
- 打开 iso 文件。
- “写入硬盘映像”
  - ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/iso/20230214180238.png)
- 在”硬盘驱动器“栏选择想要写入数据的 U 盘，点击“写入“
