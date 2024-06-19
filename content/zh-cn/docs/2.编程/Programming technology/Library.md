---
title: Library
linkTitle: Library
date: 2024-06-20T12:31
weight: 20
---

# 概述

> 参考：
>
> - [Wiki，Library(computing)](https://en.wikipedia.org/wiki/Library_(computing))
> - [公众号 - 码农的荒岛求生，动态库和静态库有什么区别？](https://mp.weixin.qq.com/s/9pavORd5qjqEaKN7G_NBmw)
>   - [B 站视频](https://www.bilibili.com/video/BV1fb421q7gc)

**Library(库)** 是一个只读资源的集合(collection)，用来实现计算机程序。这个 collection 中通常是很多已经写好的**可复用**的代码，类似于代码中的 [Function](docs/2.编程/计算机科学/Function/Function.md)。相对代码文件中的 Function、Library 则更像是存在于代码文件外部的 Function，表现一种可执行代码的 **Binary(二进制)** 文件、纯文本代码文件，甚至随着发展，还可能包括图像。

> e.g. 程序可以使用 Library 来间接进行 [System Call](docs/1.操作系统/Kernel/System%20Call/System%20Call.md)，而不是直接在程序中编写系统调用的相关代码。

一个 Library 可以被多个独立的使用者（程序、其他 Library）使用以，modular(模块化) 的方式 code reuse(重用代码)。

## Static lib 与 Dynamic lib

Library 通常分为两大类

- **Static lib(静态库)** # 编译程序时，将 Library 打包进最终的可执行文件
  - 通常可以通过代码中使用 import、etc. 关键字在文件开头导入，各种语言和环境的 Library 形态不太一样。有的语言直接把代码源码放到本地，编译时引用；有的语言会把代码编译成 .a、.lib 文件，在编译时引用；etc.
- **Dynamic lib(动态库)** # 程序运行时，将 Library 加载到内存中
  - 通常以文件的形式。Linux 中文件后缀为 .so，Windows 中文件后缀为 .dll。

# Linker

> 参考:
>
> - [Wiki，Linker_(computing)](https://en.wikipedia.org/wiki/Linker_(computing))

Library 在程序链接或绑定过程中很重要，它解决了对库模块的称为链接或符号的引用。链接过程通常由 **Linker(链接器)** 或 **Bunder(绑定器)** 程序自动完成，该程序按给定顺序搜索一组库和其他模块。如果在给定的一组库中可以多次找到链接目标，通常不会将其视为错误。链接可以在创建可执行文件时完成（静态链接），也可以在运行时使用程序时完成（动态链接）。
