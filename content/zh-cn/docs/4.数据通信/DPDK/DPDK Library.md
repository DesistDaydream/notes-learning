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

在 Telemetry Library 相关代码 [telemetry_v2_init](https://github.com/DPDK/dpdk/blob/v25.03/lib/telemetry/telemetry.c#L599) 进行初始化，之后即可使用。初始化之后会注册一些基本的命令，其他命令则取决于我们使用的 DPDK 都引用了哪些 Librairies，以及这些 Libraires 是否初始化了 telemetry

- 注册了几个基本的命令
    - `/` 命令可以列出所有可用命令
    - `/info` 命令显示基本信息
    - `/help,COMMAND` 命令显示 COMMAND 的帮助信息
- 其他注册的命令则需要到各种 Libraries 的代码中查看。可以通过搜索 [init_telemetry](https://github.com/search?q=repo%3ADPDK%2Fdpdk%20init_telemetry&type=code) 关键字找到各种 Library 注册到 Telemetry 的命令，e.g. [ethdev](https://github.com/DPDK/dpdk/blob/v25.03/lib/ethdev/rte_ethdev_telemetry.c#L1540), [mempool](https://github.com/DPDK/dpdk/blob/v25.03/lib/mempool/rte_mempool.c#L1600), etc.

简单示例如下:

```bash
~]# dpdk-telemetry.py 
Connecting to /var/run/dpdk/rte/dpdk_telemetry.v2
{
  "version": "DPDK 23.11.0",
  "pid": 4054033,
  "max_output_len": 16384
}
Connected to application: "usps"
--> /
{
  "/": [
    "/",
    "/cnxk/ethdev/info",
    "/eal/memseg_info",
    "/ethdev/info",
    "/ethdev/link_status",
    "/ethdev/list",
    "/eventdev/dev_dump",
    ......有很多，都省略了
  ]
}
--> /help,/ethdev/list
{
  "/help": {
    "/ethdev/list": "Returns list of available ethdev ports. Takes no parameters"
  }
}
--> /ethdev/list
{
  "/ethdev/list": [
    0,
    1
  ]
}
--> /help,/ethdev/stats
{
  "/help": {
    "/ethdev/stats": "Returns the common stats for a port. Parameters: int port_id"
  }
}
--> /ethdev/stats,0
{
  "/ethdev/stats": {
    "ipackets": 0,
    "opackets": 0,
    "ibytes": 0,
    "obytes": 0,
    "imissed": 0,
    "ierrors": 0,
    "oerrors": 0,
    "rx_nombuf": 0,
    "q_ipackets": [
      0,
......
```

> [!Tip] 在个人的[学习项目](https://github.com/DesistDaydream/go-dpdk/blob/main/cmd/telemetry/telemetry.go)中使用 Go 语言通过使用 unixpacket 与 Socket 文件建立连接后，也可以实现 dpdk-telemetry.py 的效果

从 DPDK 的 [API](https://doc.dpdk.org/api/) 也可以查看一些，各种命令返回信息的含义

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dpdk/20250419221259237.png)

**TODO**: 还有没有其它官方文档来代替 API 文档或者源码中来了解各种命令返回的信息？

**TODO**: 根据 https://github.com/search?q=repo%3ADPDK%2Fdpdk%20init_telemetry&type=code 总结一下都有哪些库注册了 Telemetry

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

### ethdev

**/ethdev/stats** # 端口（网卡）的普通统计信息

- 通常以 i, o 开头
- https://doc.dpdk.org/api/structrte__eth__stats.html
- https://github.com/DPDK/dpdk/blob/cd60dcd503b91956f966a1f6d595b35d256ac00f/lib/ethdev/rte_ethdev_telemetry.c#L85

**/ethdev/xstats** # 端口（网卡）的扩展统计信息

- 通常以 rx, tx 开头
- **TODO**: 这些扩展信息的详细字段在哪里？
