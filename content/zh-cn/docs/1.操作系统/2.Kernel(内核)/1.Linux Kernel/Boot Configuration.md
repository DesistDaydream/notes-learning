---
title: Boot Configuration
---

# 概述

> 参考：
> - [Linux Kernel 官网文档，Linux 内核用户和管理员指南-内核引导配置](https://www.kernel.org/doc/html/latest/admin-guide/bootconfig.html)

Boot Configuration(引导配置) 扩展了当前内核命令行，在引导内核时，可以提供额外的运行时配置。该文件默认在 `/boot/config-$(uname -r)`，每个内核版本都有一个对应的文件。

该文件有几千行，每一行都是一个以 `=` 分割的键值对，用来在系统启动内核前的引导阶段，为内核配置运行时行为

# Linux Namespace 配置

CONFIG_CHECKPOINT_RESTORE=y
CONFIG_NAMESPACES=y # 是否启用 Linux Namespace
CONFIG_UTS_NS=y # 是否启用 UTS NS
CONFIG_IPC_NS=y # 是否启用 IPC NS
CONFIG_USER_NS=y
CONFIG_PID_NS=y
CONFIG_NET_NS=y # # 是否启用 NET NS
CONFIG_UIDGID_STRICT_TYPE_CHECKS=y
CONFIG_SCHED_AUTOGROUP=y
CONFIG_MM_OWNER=y
