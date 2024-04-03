---
title: Bypass tap
linkTitle: Bypass tap
date: 2024-02-26T23:09
weight: 20
---

# 概述

> 参考：
>
> - [Wiki，Bypass tap](https://en.wikipedia.org/wiki/Bypass_switch)

**Bypass tap(旁路分路器)** 也称为 Bypass switch(旁路交换机)。Bypass tap 是一种硬件设备，可以与安全设备并联并串联到网络链路中，为安全设备提供 **fail-safe access(故障时可安全访问)** 的能力。

Bypass tap 通常至少有 4 个端口。

- A 和 B 两个端口串联，且中间不经过任何电路，Bypass tap 本身不运行的情况下也可以保证 A 到 B 之间的链路畅通。当安全设备正常运行时，A 到 B 的连接是断开的；
- C 和 D 是用来监控安全设备的端口，安全设备正常运行时，流量通过 C 和 D 端口，相当于将安全设备串联进网络中。
- 当检测到安全设备出现异常时，将会切断 C 和 D 的端口，将流量转交给 A 和 B 以保证网络链路上的数据不间断。

![bypass.drawio.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/information_security/bypass_and_dpi_1.png)

通常来说，这两种情况可以用两种模式来概括

- **Normal mode** # 流量经过 Bypass 连接的其他设备后，在发送到下一跳的设备
  - 有的地方也称为 **控制模式**
- **Bypass mode** # 流量直接经过 Bypass tap 到下一跳的设备
  - 有的地方也称为 **直通模式**

> 在中文环境中也有的将 Bypass tap 称为 Optical swap(光切换设备，Optiswap)、光开关、etc.
