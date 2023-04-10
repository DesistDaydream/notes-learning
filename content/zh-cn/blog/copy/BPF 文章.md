---
title: "BPF 文章"
linkTitle: "BPF 文章"
slug: "bpf_article"
weight: 20
---

# 概述


[颠覆传统、应用大爆发，eBPF何以改变Linux？](https://mp.weixin.qq.com/s/ejm2wPAJiSfTTyyTIsBWvQ)

开发者如何更好地学习 eBPF 技术？

InfoQ：对于开发者来说，如何才能更好地学习 eBPF？

陈莉君教授：首先要掌握 Linux 系统操作系统知识，再了解 eBPF 的大概原理，它是一种 Linux 内核可扩展性技术，它允许用户探索系统行为，以便在内核中运行自己的程序。此外，还要了解 eBPF 的关键概念，包括 BPF 程序、套接字程序、文件和目录程序、日志程序、安全程序、跟踪程序、XDP 程序，等等。

熟悉 eBPF 的编程语言，eBPF 程序使用 C 来编写（用户空间程序可以采用 Go、Rust 等语言），并且可以通过用户空间库调用它们的特定功能。可以使用 LLVM、clang、gcc 等编译器来编译 eBPF 程序。目前，龙蜥社区推出的 Coolbpf 项目，可以简化这部分的开发工作。


