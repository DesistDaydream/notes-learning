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

https://doc.dpdk.org/api/rte__ethdev_8h.html

**/ethdev/stats** # 端口（网卡）的普通统计信息

- 通常以 i, o 开头
- https://doc.dpdk.org/api/structrte__eth__stats.html
- https://github.com/DPDK/dpdk/blob/cd60dcd503b91956f966a1f6d595b35d256ac00f/lib/ethdev/rte_ethdev_telemetry.c#L85
- **ipackets** # 成功接收到的数据包总数
- **opackets** # 发送成功的数据包总数
- **imissed** # 由于没有可用缓冲区（i.e. Rx 队列已满）而被硬件丢弃的 Rx 数据包总数。
    - Tips: 人话: 数据包进入缓冲区后，使用 DPDK 的上层应用程序处理不过来，导致缓冲区的队列一直是满的，所以就把无法进入队列的数据包扔掉了。
- ...... etc.

**/ethdev/xstats** # 端口（网卡）的扩展统计信息

- 通常以 rx, tx 开头
- **TODO**: 这些扩展信息的详细字段在哪里？

# imissed

## DPDK 的 ethdev 中有个统计信息 imissed。怎么理解这个？

`imissed` 是 DPDK ethdev 统计信息中的一个重要字段，表示 **由于接收队列满或其他硬件资源不足而被丢弃的数据包数量**。

具体来说：

1. **定义位置**：在 `rte_eth_stats` 结构体中，`imissed` 表示因硬件接收资源不足导致的丢包
2. **丢包原因**：
    - RX 队列（ring buffer）已满，新到的包无法入队
    - 网卡的接收描述符用尽
    - 硬件缓冲区不足
3. **与其他丢包统计的区别**：
    - `imissed`：硬件层面的丢包（包已到达网卡但无法接收）
    - `ierrors`：接收错误（如 CRC 错误、对齐错误等）
    - `rx_nombuf`：因 mbuf 不足导致的丢包（软件层面）
4. **常见触发场景**：
    - 应用处理速度跟不上网卡接收速度
    - 分配给 RX 队列的描述符数量太少
    - CPU 核心被其他任务占用，来不及处理收包

如果看到 `imissed` 持续增长，通常需要：

- 增加 RX 描述符数量
- 优化应用处理逻辑
- 检查 CPU 亲和性设置
- 考虑使用多队列分流

不同网卡驱动对 `imissed` 的实现可能略有差异，但核心含义是硬件层面因资源不足导致的丢包。

## 这个 rx 队列是做什么用的？是网卡接收包的时候满了，还是程序接收包的时候满了？

这个 RX 队列是 **网卡和应用程序之间共享的数据结构**，用于传递接收到的数据包。理解它需要看完整的数据包接收流程：

### 数据包接收流程

1. **数据包到达网卡**
    - 网卡硬件接收到物理层的数据包
    - 网卡将数据包通过 DMA 写入内存中的 RX 描述符指向的缓冲区
2. **RX 队列（Ring Buffer）的作用**
    - 这是一个环形队列，由应用程序在初始化时创建
    - 队列中每个描述符指向一个 mbuf（数据包缓冲区）
    - 网卡通过 DMA 将数据包内容写入这些 mbuf
3. **网卡写入 → 程序读取**
    - **网卡侧**：不断将新到的包写入队列，移动写指针
    - **程序侧**：调用 `rte_eth_rx_burst()` 从队列读取，移动读指针

### imissed 发生的时机

**`imissed` 是网卡写入时发现队列满了**：

```
网卡接收包 → 尝试写入 RX 队列 → 发现队列满（所有描述符都被占用）→ imissed++
```

具体来说：

- 网卡已经将包 DMA 到内存了
- 但发现 RX 队列的描述符都被占用（应用还没来得及处理）
- 网卡无法将新包的描述符放入队列
- 只能丢弃这个包，`imissed` 计数增加

### 为什么会满？

RX 队列满的根本原因是 **应用程序处理速度跟不上网卡接收速度**：

- 应用调用 `rte_eth_rx_burst()` 取包太慢
- 取出的包还没处理完（或处理太慢）
- 导致 RX 队列中的描述符无法被释放回来给网卡继续使用

所以虽然队列是在网卡写入时满的，但根源是程序处理不过来。

## DMA 是什么？

DMA 是 **Direct Memory Access（直接内存访问）** 的缩写，是一种硬件机制，允许外设（如网卡、磁盘）直接访问系统内存，而不需要 CPU 参与数据搬运。

### 传统方式 vs DMA

**传统方式（没有 DMA）**：

```
网卡 → CPU 读取 → CPU 写入内存
```

- CPU 必须参与每个字节的搬运
- CPU 被占用，无法做其他事情
- 效率低，尤其是大量数据时

**DMA 方式**：

```
网卡 ← 告诉地址 ← CPU
网卡 → 直接写入内存（CPU 可以做其他事）
网卡 → 完成后中断通知 CPU
```

- CPU 只需告诉网卡"把数据放到这个内存地址"
- 网卡自己完成数据传输
- CPU 被解放出来做其他计算

### 在 DPDK 中的应用

当网卡接收数据包时：

1. **初始化阶段**：
    - 应用程序分配 mbuf（内存缓冲区）
    - 将这些 mbuf 的物理地址填入 RX 描述符
    - 告诉网卡：收到包就放到这些地址
2. **收包阶段**：
    - 网卡收到数据包
    - **通过 DMA 直接将包内容写入 mbuf 的内存地址**
    - 不需要 CPU 搬运数据
3. **通知阶段**：
    - 网卡更新描述符状态
    - 应用程序轮询（或中断）发现有新包
    - 直接读取 mbuf 中的数据

### DMA 的优势

- **高效**：CPU 不参与数据搬运，减少 CPU 开销
- **快速**：硬件直接操作内存，速度快
- **并行**：CPU 可以在 DMA 传输时做其他计算

这就是为什么高性能网络处理（如 DPDK）严重依赖 DMA —— 它让数据传输几乎不消耗 CPU 资源。