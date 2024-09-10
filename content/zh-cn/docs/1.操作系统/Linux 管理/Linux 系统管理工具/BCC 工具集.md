---
title: BCC 工具集
---

# 概述

> 参考：
>
> - [官网](https://iovisor.github.io/bcc/)
> - [GitHub 项目, iovisor/bcc](https://github.com/iovisor/bcc)

**BPF Compiler Collection(BPF 编译器合集，简称 BCC)** 是用于创建有效的内核跟踪和操作程序的工具包。BCC 是 Linux 基金会旗下的 IO Visor 项目组做出来的基于 eBPF 的产品。BBC 主要用来为 Linux 提供 **Dynamic Tracing(动态追踪)** 功能的实现。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ng174l/1619076409581-b90915a0-9bcb-4aa7-8ea4-4f0d66048ddd.png)

# BCC 安装

## 通过 Linux 包管理器安装

### Ubuntu

标准的 Ubuntu Universe 仓库 与 iovisor 的 PPA 仓库中都可以用来安装 BCC 工具，但是包的名称不同。Ubuntu 安装完的程序，其名称会在最后加上 `-bpfcc`。

- 使用 Ubuntu 仓库安装

```bash
sudo apt-get install bpfcc-tools linux-headers-$(uname -r)
```

- 使用 iovisor 仓库安装

```bash
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 4052245BD4284CDD
echo "deb https://repo.iovisor.org/apt/$(lsb_release -cs) $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/iovisor.list
sudo apt-get update
sudo apt-get install bcc-tools libbcc-examples linux-headers-$(uname -r)
```

### CentOS

# BCC 工具概述

命名规则

XXXsnoop 这类工具的名字通常用来追踪指定对象，snoop 有窥探之意。比如 opensnoop 工具用来追踪 open() 系统调用、execsnoop 工具用来追踪 exec() 系统调用 等等。

- **syscount** # 追踪系统调用，并统计次数
- **tcpconnect** # 追踪活动的 TCP 连接，即 `connect()` 系统调用。
- **tcptracer** # 追踪
