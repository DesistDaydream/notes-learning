---
title: Graph Library
linkTitle: Graph Library
date: 2024-05-21T17:13
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，Graph Library and Inbuilt Nodes](https://doc.dpdk.org/guides/prog_guide/graph_lib.html)
> - https://zhuanlan.zhihu.com/p/604202266
> - https://zhuanlan.zhihu.com/p/613233087

DPDK 的 **Graph(图)** Library 将数据处理功能抽象为 **Node(节点)**，并将这些 Node links(链接) 在一起以创建一张大型的数据处理流程图，以实现可重用的/模块化的数据处理能力。一个 Node 中可以有一个或多个流量处理 **Function(功能)**，一个 Node 处理完成后，交给下一个或几个其他 Node 继续处理流量数据

![image.png|800](https://notes-learning.oss-cn-beijing.aliyuncs.com/dpdk/graph_library_1.png)

假如我们设计了如下一系列功能

- decode
- flow
- reassemble
- resource
- control
- asset
- security
- record
- drop
- ......略

可以将这些功能分散到多个 Node 中，每个 Node 又可以规划如何如何调用这些功能

```json
    "node_m": [
      {
        "name": "sink",
        "next": [
          "flow"
        ]
      },
      ...... 略
    ],
    "node_n": [
      {
        "name": "rx",
        "next": [
          "decode"
        ]
      },
      {
        "name": "decode",
        "next": [
          "traffic_filter"
        ]
      },
      {
        "name": "traffic_filter",
        "next": [
          "capture",
          "drop"
        ]
      },
      ...... 略
    ],
    "node_o": [
      {
        "name": "rx",
        "next": [
          "decode"
        ]
      },
      {
        "name": "decode",
        "next": [
          "flow"
        ]
      },
      ...... 略
    ],
    ......略
```

# Graph 架构

Graph Library 的设计思想源自于开源项目 [Vector Packet Processor(VPP)](https://github.com/FDio/vpp)
