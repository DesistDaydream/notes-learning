---
title: Archive file
weight: 1
---

# 概述

> 参考：
>
> - [Wiki, Archive file](https://en.wikipedia.org/wiki/Archive_file)
> - [Wiki, Tar](<https://en.wikipedia.org/wiki/Tar_(computing)>)

在计算机中，**Archive file(归档文件)** 是由一个或多个文件及元数据组成的一个计算机文件。归档文件用于将多个数据文件放在一起收集到一个文件中，以便于移植和存储。归档文件通常存储目录结构、错误检测和纠正信息、任意注释，有时还使用内置加密。

> [!Tip] Archive file 通常都会进行 [Data compression](/docs/5.数据存储/Archive%20file/Data%20compression.md)(数据压缩)

> [!Note] 有的时候归档文件也翻译成 “存档文件”
>
> archiver(归档程序)，好多时候中文称为压缩程序

## 归档与压缩的概念

> [!Question] 为什么要区分这两个概念？
> 把多个文件（文件类型可以不同）放在一个文件中称为 归档文件；本来 10KiB 的文件以 5KiB 存储的文件称为 压缩文件。
>
> 那把归档文件压缩后，因为称为什么呢？

所以我个人感觉，可以把多个文件放在一个文件中称为打包，那么，归档就是指 **打包 + 压缩**，也可以指经过压缩的归档文件（e.g. 一个单日志文件压缩后也可以称为归档文件）。很多时候，对于一个已经归档并压缩的文件，人们称为：打包文件、压缩文件、etc.。

> [!Note] 另外，是不是还有一种理解方向：归档文件就是指一个或多个文件通过某种格式合并成一个文件后，再进行压缩形成的文件。既然最后都是压缩，也就干脆叫压缩文件。

所以我看 7-Zip 官方页面就把 `7-Zip is a file archiver with a high compression ratio.` 翻译为 `7-Zip 是一款拥有极高压缩比的开源压缩软件。` i.e. archiver 翻译成压缩软件，而不是归档软件。

e.g.

Linux 中默认自带的  [tar与gzip](/docs/1.操作系统/Linux%20管理/Linux%20系统管理工具/文件与文件系统管理工具/tar与gzip.md) 通过打包与压缩，能生成一个归档文件。gzip 程序只能针对一个文件进行压缩，这样当你想要压缩一大堆文件时，先将这一大堆文件先打成一个包（tar 命令），然后再用压缩程序进行压缩（gzip, bzip2, etc. 命令）。i.e. `tar -z` 则会在打包后，自动调用 gzip 命令进行压缩。

跨平台通用的 7z 则自带打包与压缩，也能生成一个归档文件

# 归档程序

Tar 是一种计算机应用程序，用于将许多文件汇集到一个 **Archive file(归档文件)** 中，通常称为 **Tarball**。该名称源于 Tape Archive(磁带存档)，取 Tape 中的 t 和 Archive 中的 ar。

利用 [tar](/docs/1.操作系统/Linux%20管理/Linux%20系统管理工具/文件与文件系统管理工具/tar与gzip.md) 命令，可以把一大堆的文件和目录全部打包成一个文件，这对于备份文件或将几个文件组合成为一个文件以便于网络传输是非常有用的。利用 tar，可以为某一特定文件创建档案（备份文件），也可以在档案中改变文件，或者向档案中加入新的文件。tar 最初被用来在磁带上创建档案，现在，用户可以在任何设备上创建档案。

# 归档文件格式

仅归档: .tar, etc.

仅压缩: .gz, etc.

归档与压缩: .7z, .zip, .rar, .tar.gz etc.

各种归档格式没有一个统一的标准，但是通常包含如下两类元数据：

- 结构元数据 # 包含 成员清单、文件大小、校验、etc.
- 文件系统元数据 # 成员权限、时间戳、etc.

通常，所有归档格式都要包含结构元数据，不一定包含文件系统元数据。

# zip

> 参考：
>
> - [官方文档](https://infozip.sourceforge.net/)
> - [规范](https://www.pkware.com/documents/casestudies/APPNOTE.TXT)

unzip, zip, etc. 程序

# 7z

> 参考：
>
> - [SourceForge 项目，sevenzip](https://sourceforge.net/projects/sevenzip/)
> - [GitHub 项目，ip7z/7zip](https://github.com/ip7z/7zip)
> - [官网](https://www.7-zip.org/)

7-Zip 是一款拥有极高压缩比的开源归档软件。

使用 [.7z](https://github.com/ip7z/7zip/blob/main/DOC/7zFormat.txt) 归档格式，使用 **LZMA** 与 **LZMA2** 压缩算法。

可执行文件用途

- 7z # CLI 命令行工具
- 7zFM # 
- 7zG # 

## Syntax

