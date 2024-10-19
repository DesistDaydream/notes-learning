---
title: Archive File(归档文件)
linkTitle: Archive File(归档文件)
date: 2024-04-20T13:51
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, Tar](<https://en.wikipedia.org/wiki/Tar_(computing)>)
> - [Wiki, Archive_file](https://en.wikipedia.org/wiki/Archive_file)

在计算机中，**Archive File(归档文件)** 是由一个或多个文件及元数据组成的一个计算机文件。归档文件用于将多个数据文件放在一起收集到一个文件中，以便于移植和存储。归档文件通常存储目录结构、错误检测和纠正信息、任意注释，有时还使用内置加密。

> 有的时候归档文件也翻译成 “存档文件”

## 归档与压缩的概念

首先要弄清两个概念：打包和压缩。打包是指将一大堆文件或目录变成一个总的文件；压缩则是将一个大的文件通过一些压缩算法变成一个小文件。

为什么要区分这两个概念呢？这源于 Linux 中很多压缩程序只能针对一个文件进行压缩，这样当你想要压缩一大堆文件时，你得先将这一大堆文件先打成一个包（tar 命令），然后再用压缩程序进行压缩（gzip bzip2 命令）。

## 归档程序

Tar 是一种计算机应用程序，用于将许多文件汇集到一个 **Archive file(归档文件)** 中，通常称为 **Tarball**。该名称源于 Tape Archive(磁带存档)，取 Tape 中的 t 和 Archive 中的 ar。

利用 tar 命令，可以把一大堆的文件和目录全部打包成一个文件，这对于备份文件或将几个文件组合成为一个文件以便于网络传输是非常有用的。利用 tar，可以为某一特定文件创建档案（备份文件），也可以在档案中改变文件，或者向档案中加入新的文件。tar 最初被用来在磁带上创建档案，现在，用户可以在任何设备上创建档案。

# 数据压缩

> 参考：
>
> - [Wiki, Data compression](https://en.wikipedia.org/wiki/Data_compression)(数据压缩)
