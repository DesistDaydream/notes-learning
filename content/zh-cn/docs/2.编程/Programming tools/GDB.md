---
title: GDB
linkTitle: GDB
weight: 41
---

# 概述

> 参考：
>
> - [官网](https://www.sourceware.org/gdb/)

**GNU Project debugger(GNU 项目调试器，简称 GDB)** 允许您查看另一个程序在执行时“内部”发生了什么，或者另一个程序在崩溃时正在做什么。

GDB 是支持多语言的 [Debugger](docs/2.编程/Programming%20tools/Debugger.md)

- Ada
- Assembly
- [C](docs/2.编程/高级编程语言/C/C.md)
- C++
- D
- Fortran
- [Go](docs/2.编程/高级编程语言/Go/Go.md)
- Objective-C
- OpenCL
- Modula-2
- Pascal
- [Rust](docs/2.编程/高级编程语言/Rust/Rust.md)

GDB 可以执行四种主要操作（以及其他辅助操作），以帮助您及时发现错误：

- 启动程序，并指定任何可能影响其行为的因素。
- 让程序在满足特定条件时停止运行。
- 检查一下程序停止运行时发生了什么。
- 修改程序中的某些内容，以便您可以尝试纠正一个错误的影响，并继续了解另一个错误。

这些程序可能在与 GDB 相同的机器上（本地运行），也可能在另一台机器上（远程运行），或者在模拟器上运行。GDB 可以在大多数流行的 UNIX 和 Microsoft Windows 系统以及 macOS 上运行。
