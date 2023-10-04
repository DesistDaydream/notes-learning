---
title: trace 子命令
---

## 概述

> 参考：
> 
> - [Manual(手册)，perf-trace(1)](https://man7.org/linux/man-pages/man1/perf-trace.1.html)

perf trace 是类似 strace 的追踪工具，可以追踪系统调用。

# Syntax(语法)

perf trace \[OPTIONS]

OPTIONS

- **-e, --expr, --event \<SYSCALL>** # 显示指定的系统调用
- **-s** # 统计每一次系统调用的执行时间、次数、错误次数
- **-F,--pf={all | min | maj}** #

