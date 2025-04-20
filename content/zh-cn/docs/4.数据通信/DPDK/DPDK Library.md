---
title: Library
linkTitle: Library
weight: 20
date: 2025-04-18T18:24:00
---

# 概述

> 参考：
>
> - [官方文档，开发者指南 - XXX 库](https://doc.dpdk.org/guides/prog_guide/index.html)

DPDK 的主要对外函数接口通常以 `rte_`(runtime environment) 作为前缀。

# 内存管理

- [Memory Pool Library](https://doc.dpdk.org/guides/prog_guide/mempool_lib.html)
- etc.

# CPU 管理

- etc.

# CPU 包处理

- etc.

# Device Libraries

- etc.

# Protocol Processing Libraries

- etc.

# High-Level Libraries

- [Graph Library](/docs/4.数据通信/DPDK/Graph%20Library.md)
- etc.

# Utility Libraries

- [Metrics Library](https://doc.dpdk.org/guides/prog_guide/metrics_lib.html)
- [Telemetry Library](#Telemetry%20Library)(遥测库) # 遥测库提供了一个接口，用于从各种 DPDK 库中检索信息。该库通过 Unix Socket 提供这些信息，接收来自客户端的请求，并回复包含所请求遥测信息的 JSON 响应。
- etc.

## Telemetry Library

> 参考：
>
> - [官方文档，开发者指南 - 遥测库](https://doc.dpdk.org/guides/prog_guide/telemetry_lib.html)

在 Telemetry Library 相关代码 [telemetry_v2_init](https://github.com/DPDK/dpdk/blob/v25.03/lib/telemetry/telemetry.c#L599) 进行初始化，注册了几个基本的命令（`/`, `/info`, `/help`）。其他注册的命令则需要到各种 Libraries 的代码中查看。可以通过搜索 [init_telemetry](https://github.com/search?q=repo%3ADPDK%2Fdpdk%20init_telemetry&type=code) 关键字找到各种 Library 注册到 Telemetry 的命令，比如 [ethdev](https://github.com/DPDK/dpdk/blob/v25.03/lib/ethdev/rte_ethdev_telemetry.c#L1540), [mempool](https://github.com/DPDK/dpdk/blob/v25.03/lib/mempool/rte_mempool.c#L1600), etc.

从 DPDK 的 [API](https://doc.dpdk.org/api/) 也可以查看一些，各种命令返回信息的含义

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dpdk/20250419221259237.png)

TODO: 还有没有其它官方文档来代替 API 文档或者源码中来了解各种命令返回的信息？

TODO: 根据 https://github.com/search?q=repo%3ADPDK%2Fdpdk%20init_telemetry&type=code 总结一下都有哪些库注册了 Telemetry

cnxk_mempool_init_telemetry

cnxk_ethdev_init_telemetry

ring_init_telemetry

librawdev_init_telemetry

dmadev_init_telemetry

security_init_telemetry

mempool_init_telemetry

eventdev_init_telemetry

cryptodev_init_telemetry

cnxk_ipsec_init_telemetry

ta_init_telemetry

ethdev_init_telemetry

rxa_init_telemetry 

