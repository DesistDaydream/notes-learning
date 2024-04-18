---
title: sysbench
---

# 概述

> 参考
>
> - [Github,akopytov/sysbench](https://github.com/akopytov/sysbench)

脚本数据库和系统性能基准

# Syntax(语法)

**sysbench \[OPTIONS] \[TestName] \[COMMAND]**

TestName

- **fileio **# 文件 I/O 测试
- **cpu** # CPU 性能测试
- **memory** # 内存功能速度测试
- **threads** # 线程子系统性能测试
- **mutex** # 互斥体性能测试

## General OPTIONS

- **--threads \<INT>** # 要创建的工作线程总数
- **--time \<DURATION>** # 运行持续时间。单位：秒。`默认值：0`。0 表示没有限制。

# Example

- 以 10 个线程运行 5 分钟的线程子系统性能测试。常用来模拟多线程切换的问题
  - **sysbench --threads=10 --time=300 threads run**
