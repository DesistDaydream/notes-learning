---
title: LLVM
linkTitle: LLVM
date: 2024-09-12T09:57
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，llvm/llvm-project](https://github.com/llvm/llvm-project)
> - [Wiki, LLVM](https://en.wikipedia.org/wiki/LLVM)

LLVM 是一组 [Compiler](/docs/2.编程/Programming%20tools/Compiler.md)(编译器) 和 工具链技术，可以用于为任何编程语言开发前端（Notes: 不是 Web 开发中的前端），为任何指令集架构开发后端。LLVM 围绕一种与语言无关的中间表示（IR）设计，这种表示作为一种可移植的高级汇编语言，可以通过多次转换进行优化。

> [!Tip]
> 然而，LLVM项目逐渐演变为一个涵盖多个子项目的综合性项目，这与大多数当前开发者所理解的虚拟机关系不大。这使得该首字母缩略词变得“令人困惑”和“不合适”，因此自2011年以来，LLVM 官方不再是一个缩略词，而是一个适用于 LLVM 综合项目的品牌。该项目包括 LLVM [intermediate representation](https://en.wikipedia.org/wiki/Intermediate_representation)(中间表示，简称 IR)、LLVM 调试器、LLVM 对 C++ 标准库的实现（完全支持 C++ 11 和C++ 14）等。LLVM 由 LLVM 基金会管理。编译工程师 Tanya Lattner 于 2014 年成为其主席。
