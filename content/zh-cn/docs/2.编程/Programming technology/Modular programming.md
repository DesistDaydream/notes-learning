---
title: Modular programming
linkTitle: Modular programming
weight: 22
---

# 概述

> 参考：
>
> - [Wiki, Modular programming](https://en.wikipedia.org/wiki/Modular_programming)(模块化编程)
> - https://en.wikipedia.org/wiki/Code_reuse

**Modular programming(模块化编程)**

在计算机程序的开发过程中，随着程序代码越写越多，会产生很多问题：

- 想要使用已有源代码，而不是重新自己写一遍。需要 [代码复用](https://en.wikipedia.org/wiki/Code_reuse)
- 代码复用将会产生 **dependencies(依赖性)** 问题
- 在一个文件里代码会越来越长，越来越不容易维护，需要将代码拆分到多个文件中
- ...... etc.

为了编写可维护的代码，我们把很多代码分组，分别放到不同的文件里，这样，每个文件包含的代码就相对较少，很多编程语言都采用这种组织代码的方式。

此时，调用另一个文件中的源代码则需要先知道要调用的代码所在文件中的位置。

# 模块化编程的实践

[Library](/docs/2.编程/Programming%20technology/Library.md)(库)

Framework(框架)