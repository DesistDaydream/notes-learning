---
title: Clang
linkTitle: Clang
weight: 13
---

# 概述

> 参考：
>
> - [Wiki, Clang](https://en.wikipedia.org/wiki/Clang)

Clang 是 C、C++、Objective-C 和 Objective-C++ 编程语言以及 OpenMP、OpenCL、RenderScript、CUDA、SYCL 和 HIP 的 [Compiler](/docs/2.编程/Programming%20tools/Compiler.md)(编译器) 前端框架。

对于 LLVM 来说，Clang 充当 [GCC](/docs/2.编程/Programming%20tools/GCC.md)(GNU 编译器集合) 的直接替代品，支持其大部分编译标志和非官方语言扩展。它包括一个静态分析器和几个代码分析工具。

Clang 与 [LLVM](/docs/2.编程/Programming%20tools/LLVM.md) 编译器后端协同工作，并且一直是 LLVM 2.6 及更高版本的子项目。与 LLVM 一样，它是根据 Apache 2.0 软件许可证发布的自由开源软件 。其贡献者包括 Apple 、 IBM 、 Microsoft 、 Google 、 ARM 、 Sony 、 Intel 和 AMD 。

> [!Attention] 上面描述中的前端、后端是指[编译器的三阶段](/docs/2.编程/Programming%20tools/Compiler.md#编译器的三阶段)中的概念

## 背景

LLVM 项目最初计划使用 GCC 的前端。然而，GCC 的源代码庞大且略显繁琐；正如一位长期从事 GCC 开发的开发者在谈到 LLVM 时所说：“试图让河马跳舞可不是什么有趣的事。”此外，苹果软件使用 Objective-C ，而 Objective-C 对 GCC 开发者而言优先级较低。因此，GCC 无法顺利集成到苹果的集成开发环境中。 最后，GCC 的许可协议 ——GNU 通用公共许可证 (GPL) 第 3 版 ——要求分发 GCC 扩展或修改版本的开发者必须公开其源代码 ，而 LLVM 的宽松软件许可则没有这项要求。[ 6 ]

基于这些原因，苹果公司开发了 Clang，这是一个支持 C、Objective-C 和 C++ 的新型编译器前端。 2007 年 7 月，该项目获得批准成为开源项目。

