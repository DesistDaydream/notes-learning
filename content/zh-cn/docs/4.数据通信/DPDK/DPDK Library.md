---
title: Library
linkTitle: Library
weight: 101
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
- [Telemetry Library](docs/4.数据通信/DPDK/Telemetry%20Library.md)(遥测库) # 遥测库提供了一个接口，用于从各种 DPDK 库中检索信息。该库通过 Unix Socket 提供这些信息，接收来自客户端的请求，并回复包含所请求遥测信息的 JSON 响应。
- etc.
