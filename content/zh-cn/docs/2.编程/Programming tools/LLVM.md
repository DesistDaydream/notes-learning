---
title: LLVM
linkTitle: LLVM
weight: 10
---

# 概述

> 参考：
>
> - [GitHub 项目，llvm/llvm-project](https://github.com/llvm/llvm-project)
> - [Wiki, LLVM](https://en.wikipedia.org/wiki/LLVM)

LLVM 是一组 [Compiler](/docs/2.编程/Programming%20tools/Compiler.md)(编译器) 和 工具链技术，可以用于开发任何编程语言的前端，以及任何指令集架构的后端。LLVM 围绕一种与语言无关的 IR(中间表示) 设计，这种表示作为一种可移植的高级汇编语言，可以通过多次转换进行优化。

> [!Attention] 上面描述中的前端、后端是指[编译器的三阶段](/docs/2.编程/Programming%20tools/Compiler.md#编译器的三阶段)中的概念

LLVM 的架构允许为**任何编程语言**开发相应的前端，这些前端都将源代码转换成统一的 LLVM IR，然后由 LLVM 的后端将 IR 转换成特定处理器架构的机器码。这种设计实现了前端和后端的解耦，使得一个新语言只需要开发前端即可利用 LLVM 的所有优化和后端支持。

> [!Tip]
> LLVM 最初是 Low Level Virtual Machine(低级虚拟机) 的缩写。然而，LLVM 项目逐渐演变为一个涵盖多个子项目的综合性项目，这与大多数当前开发者所理解的虚拟机关系不大。这使得该首字母缩略词变得“令人困惑”和“不合适”，因此自2011年以来，LLVM 官方不再是一个缩略词，而是一个适用于 LLVM 综合项目的品牌。该项目包括 LLVM [intermediate representation](https://en.wikipedia.org/wiki/Intermediate_representation)(中间表示，简称 IR)、LLVM 调试器、LLVM 对 C++ 标准库的实现（完全支持 C++ 11 和C++ 14）等。LLVM 由 LLVM 基金会管理。编译工程师 Tanya Lattner 于 2014 年成为其主席。
