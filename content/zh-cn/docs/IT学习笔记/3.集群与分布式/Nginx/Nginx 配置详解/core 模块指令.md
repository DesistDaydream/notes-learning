---
title: core 模块指令
weight: 2
---

# 概述

> 参考：
> - [org 官方文档,核心功能](http://nginx.org/en/docs/ngx_core_module.html)

main 模块主要用来为 nginx 程序的运行方式进行定义，并不涉及流量处理相关工作。

# 指令详解

**user USERNAME \[GROUPNAME]; **#指定运行 work 线程的用户和组
**pid /PATH/PidFile;** # 指定 nginx 守护进程的 pid 文件
**work_rlimit_nofile NUMBER;** # 指定所有 work 线程加起来所能打开的最大文件句柄数
