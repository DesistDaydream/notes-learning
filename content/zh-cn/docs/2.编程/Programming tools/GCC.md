---
title: GCC
linkTitle: GCC
weight: 11
---

# 概述

> 参考：
>
> - [官网](https://gcc.gnu.org/)
> - [Wiki, GNU Compiler Collection](https://en.wikipedia.org/wiki/GNU_Compiler_Collection)

**GNU Compiler Collection(GUN 编译器集合，简称 GCC)** 包括了 C、C++、Objective-C、Fortran、Ada、Go 以及 D 等语言的前端，同时也包括了这些语言的库（如 libstdc++ 等），是为 GNU 操作系统编写的编译器。GNU 系统是开发成 100% 自由软件的，这里的自由是指它遵循用户的自由原则。

> [!Attention] 上面描述中的前端、后端是指[编译器的三阶段](/docs/2.编程/Programming%20tools/Compiler.md#编译器的三阶段)中的概念

> [!Tip] GCC 最初单指 **GNU C compiler(GNU C 编译器)**

# 安装 GCC

Linux 内核本身就是 C 写的，所以一般都自带 GCC，我们安装的通常都是适用于 Windows 的 GCC，一般是 [MinGW-w64](#MinGW-w64)

## Linux

yum install gcc make -y

## MinGW-w64

> 参考：
>
> - [SourceForge 项目，mingw-w64](https://sourceforge.net/projects/mingw-w64/)
> - [GitHub 项目，mingw-w64/mingw-w64](https://github.com/mingw-w64/mingw-w64)
> - [官网](https://www.mingw-w64.org/)

mingw-w64 项目是完整的运行时环境，支持 gcc 编译生成本地运行于 Windows 64 位和 32 位操作系统的二进制文件。

打开 sourceforge 中的 [MinGW-w64](https://sourceforge.net/projects/mingw-w64/) 页面，在 [file 标签页](https://sourceforge.net/projects/mingw-w64/files/)中，下载 [x86_64-win32-seh](https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-win32/seh/x86_64-8.1.0-release-win32-seh-rt_v6-rev0.7z) 这个版本并安装即可。这是一个 tar 包，解压完成后，需要在 Windows 的 ${PATH} 环境变量中，添加解压出来的 bin 目录，通常都在 `PATH\TO\x86_64-8.1.0-release-win32-seh-rt_v6-rev0\mingw64\bin` 这里

### 其他 GCC 整合

[GitHub 项目，skeeto/w64devkit](https://github.com/skeeto/w64devkit)

- 这是 Portable(便携版) 的
- 包含 OpenMP
