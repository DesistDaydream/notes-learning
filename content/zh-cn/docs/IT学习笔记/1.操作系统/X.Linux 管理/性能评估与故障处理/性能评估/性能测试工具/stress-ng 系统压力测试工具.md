---
title: stress-ng 系统压力测试工具
---

# 概述

# Syntax(语法)

OPTIONS

- -c, --cpu <NUM> # 指定要测试的 CPU 数量，测几个就起几个进程。
- -i <N>, --io <N> # 启动 N 个工作程序，连续调用 sync(2) 将缓冲区高速缓存提交到磁盘。 可以与 --hdd 选项结合使用。
- -d <>N, --hdd <N> # 开始 N 个工作人员不断写入，读取和删除临时文件。 默认模式是对顺序写入和读取进行压力测试。 如果启用了--ggressive 选项，而没有任何--hdd-opts 选项，则 hdd Stressor 将一个接一个地处理所有--hdd-opt 选项，以涵盖一系列 I / O 选项。
- --timeout <TIME> # 指定程序运行时间

EXAMPLE

- stress-ng -c 1 --timeout 600 # 模拟 CPU 使用，导致 us 升高
- stress-ng -i 1 --hdd 1 --timeout 600 # 模拟磁盘 io，会导致 wa 升高
