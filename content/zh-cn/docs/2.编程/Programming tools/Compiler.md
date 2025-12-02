---
title: Compiler
linkTitle: Compiler
weight: 10
---

# 概述

> 参考：
>
> - [Wiki, Compiler](https://en.wikipedia.org/wiki/Compiler)
> - [YouTube, Why Some Projects Use Multiple Programming Languages](https://www.youtube.com/watch?v=XJC5WB2Bwrc)

在计算机中，**Compiler(编译器)** 是一个计算机程序，它实现了一种转换功能：将一种编程语言（源语言）编写的计算机代码翻译成另一种语言（目标语言）。术语 “编译器” 主要用于指将源代码从高级编程语言翻译成低级编程语言（例如汇编语言、目标代码或机器代码）以创建可执行程序。

> [!Tip] 口语、实践 中常挂在嘴边的编译器通常并不是一个独立的程序，而是一组工具的集合。e.g. [GCC](/docs/2.编程/Programming%20tools/GCC.md), [LLVM](/docs/2.编程/Programming%20tools/LLVM.md), etc.

这组编译器工具通常会经历 4 个主要阶段

- **Pre-Processor(预处理)** # 通过删除注释、展开宏定义、处理条件编译等操作来准备源代码，最关键的是各种语言解析某些关键字
- **Compiler(编译器)** # 将预处理后的代码转译为汇编语言，这是计算机将要执行的指令，但仍以人类可读的语言形式存在。
- **Assembler(汇编器)** # （汇编器在技术上也是另一种编译器）会接收前一阶段生成的人类可读汇编代码，并将其翻译为机器码，也就是 CPU 能理解的 0 和 1 指令。
- **Linker(连接器)** # 将所有目标文件整合起来，形成一个独立可执行文件。

# 编译器的三阶段

https://en.wikipedia.org/wiki/Compiler#Three-stage_compiler_structure

无论编译器设计中究竟有多少个阶段，这些阶段都可以归为下面这三个阶段其中之一：

- **Frontend(前端)** # 前端扫描输入，并根据特定的源语言验证语法和语义。对于静态类型语言 ，它通过收集类型信息来执行类型检查 。如果输入程序语法错误或存在类型错误，它会生成错误和/或警告消息，通常会指出源代码中检测到问题的位置；在某些情况下，实际错误可能位于程序中更早的位置。前端的功能包括词法分析、语法分析和语义分析。前端将输入程序转换为中间表示 （IR），以便中间端进行进一步处理。该中间表示通常是相对于源代码的低级表示。
- **Middleend(中端)** # 中间层对中间表示 (IR) 进行优化，这些优化与目标 CPU 架构无关。这种源代码/机器代码独立性旨在使支持不同语言和目标处理器的编译器版本之间能够共享通用优化。中间层优化的示例包括：移除无用代码（ 死代码消除 ）或不可达代码（ 可达性分析 ）、发现并传播常量值（ 常量传播 ）、将计算移至执行频率较低的位置（例如，循环外），或根据上下文对计算进行专门化，最终生成后端使用的“优化”IR。
- **Backend(后端)** # 后端从中间端获取优化后的中间表示（IR）。它可能会针对目标 CPU 架构执行更多分析、转换和优化。后端生成目标相关的汇编代码，并在过程中进行寄存器分配 。后端执行指令调度 ，通过填充延迟槽来重新排列指令，使并行执行单元保持忙碌状态。尽管大多数优化问题都是 NP 难问题 ，但解决这些问题的启发式技术已经发展成熟，并在生产级编译器中得到实现。后端的输出通常是针对特定处理器和操作系统的机器代码。

这种前端/中间/后端架构使得将不同语言的前端与不同 CPU 的后端相结合成为可能，同时还能共享中间端的优化。 这种方法的实际例子包括 [GCC](/docs/2.编程/Programming%20tools/GCC.md) 、 [Clang](/docs/2.编程/Programming%20tools/Clang.md) （基于 LLVM 的 C/C++ 编译器）和 Amsterdam 编译器工具包 ，它们都具有多个前端、共享的优化和多个后端。

## Frontend

## Middleend

## Backend

# IR

https://en.wikipedia.org/wiki/Intermediate_representation

**Intermediate representation(中间表示，简称 IR)**
