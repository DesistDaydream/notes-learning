---
title: DPDK
linkTitle: DPDK
date: 2024-03-18T11:28
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，DPDK/dpdk](https://github.com/DPDK/dpdk)
> - [官网](https://www.dpdk.org/)
> - [Wiki，Data Plane Development Kit](https://en.wikipedia.org/wiki/Data_Plane_Development_Kit)

**Data Plane Development Kit(数据平面开发套件，简称 DPDK)** 是一个由 [Linux 基金会](/docs/x_标准化/Foundation/Linux%20Foundation.md) 管理的开源软件项目。用于将 TCP 数据包的处理能力从内核空间移动到用户空间中的进程。主要是跳过了内核的 [Interrupts(中断)](/docs/1.操作系统/Kernel/CPU/Interrupts(中断)/Interrupts(中断).md) 逻辑。

处理数据包的传统方式是 CPU 中断方式，即网卡驱动接收到数据包后通过中断通知 CPU 处理，然后由 CPU 拷贝数据并交给协议栈。在数据量大时，这种方式会产生大量 CPU 中断，导致 CPU 无法运行其他程序。

而 DPDK 则采用轮询方式实现数据包处理过程：DPDK 程序加载了网卡驱动，该驱动在收到数据包后不中断通知 CPU，而是将数据包通过零拷贝技术存入内存，这时应用层程序就可以通过 DPDK 提供的接口，直接从内存读取数据包。

这种处理方式节省了 CPU 中断时间、内存拷贝时间，并向应用层提供了简单易行且高效的数据包处理方式，使得网络应用的开发更加方便。但同时，由于需要重载网卡驱动，因此该开发包目前只能用在部分采用 Intel 网络处理芯片的网卡中。

DPDK 主要包含如下几个部分（https://doc.dpdk.org/guides/prog_guide/source_org.html#libraries）

- **Environmental Abstraction Layer(环境抽象层，简称 EAL)** # 负责为应用间接访问底层的资源，比如内存空间、线程、设备、定时器等。如果把我们使用了 DPDK 的应用比作一个豪宅的主人的话，`EAL`就是这个豪宅的管家。
- **DPDK API Library** # DPDK 的 API 库
  - etc.
- **[NIC](/docs/4.数据通信/Networking%20device/NIC.md) Driver(网卡驱动程序)** # 如名，只不过是轮询模式的驱动。
  - etc.
- **DPDK APP** # 一些实用的程序

## EAL

https://doc.dpdk.org/guides/prog_guide/env_abstraction_layer.html

**Environment Abstraction Layer(环境抽象层，简称 EAL)** 负责访问低级资源，例如硬件和内存空间。它提供了一个通用接口，对应用程序和库隐藏了环境细节。初始化例程负责决定如何分配这些资源（即内存空间、设备、定时器、控制台等）。

## DPDK Library

TODO: 好多好多的库，功能非常全。

## Driver

Linux 驱动

- https://doc.dpdk.org/guides/linux_gsg/linux_drivers.html#binding-and-unbinding-network-ports-to-from-the-kernel-modules
- **vfio** # 使用 [VFIO](/docs/1.操作系统/Kernel/VFIO.md) 功能的驱动。依赖 `vfio-pci` 模块。[官方文档](https://doc.dpdk.org/guides/linux_gsg/linux_drivers.html#binding-and-unbinding-network-ports-to-from-the-kernel-modules)建议所有情况下都是用 vfio-pci 作为 DPDK 绑定端口的内核模块

TODO: 好多好多驱动。

# DPDK 部署

https://core.dpdk.org/doc/quick-start/

[这里](https://doc.dpdk.org/guides/linux_gsg/index.html) 是要在 Linux 安装 DPDK 的一些要求

DPDK 的部署依赖 meson 和 ninja 两个构建工具、vfio-pci 内核模块。

## DPDK 使用注意

当我们让 DPDK 程序接管网卡后，从系统中就无法通过常见的命令（e.g. ip, etc.）查看到系统中的网络设备，效果如下：

```bash
~]# lspci | grep Ethernet
01:00.0 Ethernet controller: Intel Corporation I350 Gigabit Network Connection (rev 01)
01:00.1 Ethernet controller: Intel Corporation I350 Gigabit Network Connection (rev 01)
02:00.0 Ethernet controller: Intel Corporation I350 Gigabit Network Connection (rev 01)
02:00.1 Ethernet controller: Intel Corporation I350 Gigabit Network Connection (rev 01)

~]# ip link show
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN mode DEFAULT group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
2: eno3: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP mode DEFAULT group default qlen 1000
    link/ether 44:a8:42:38:4e:2c brd ff:ff:ff:ff:ff:ff
5: enp2s0f1: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc mq state DOWN mode DEFAULT group default qlen 1000
    link/ether a0:36:9f:91:17:a7 brd ff:ff:ff:ff:ff:ff

~]# dpdk-devbind.py -s
Network devices using DPDK-compatible driver
============================================
0000:01:00.1 'I350 Gigabit Network Connection 1521' drv=vfio-pci unused=igb
0000:02:00.0 'I350 Gigabit Network Connection 1521' drv=vfio-pci unused=igb

Network devices using kernel driver
===================================
0000:01:00.0 'I350 Gigabit Network Connection 1521' if=eno3 drv=igb unused=vfio-pci *Active*
0000:02:00.1 'I350 Gigabit Network Connection 1521' if=enp2s0f1 drv=igb unused=vfio-pci

```

可以看到机器上有 4 个 Ethernet，但是通过 ip 命令仅能看到 2 个（ip 这类命令只能看到内核管理的），另外 2 个被 DPDK 使用了，只有使用 DPDK 提供的程序才能获取到非内核管理的网卡信息。

当我们让 DPDK 接管了网卡后，想要查看网卡当前所使用的驱动，可以使用 [lspci](/docs/1.操作系统/Linux%20管理/Linux%20硬件管理工具/Linux%20硬件管理工具.md#lspci) 命令查看到当前网卡所使用的驱动程序。（虽然其中显示的是 Kernel driver in use，但是实际上，驱动是 vfio-pci 的网卡已经被 DPDK 接管了）

```bash
]# lspci -k | grep -i "Ethernet controller" -A 3
01:00.0 Ethernet controller: Intel Corporation I350 Gigabit Network Connection (rev 01)
        DeviceName: NIC1
        Subsystem: Dell Gigabit 2P I350-t LOM
        Kernel driver in use: igb
--
01:00.1 Ethernet controller: Intel Corporation I350 Gigabit Network Connection (rev 01)
        DeviceName: NIC2
        Subsystem: Dell Gigabit 2P I350-t LOM
        Kernel driver in use: vfio-pci
--
02:00.0 Ethernet controller: Intel Corporation I350 Gigabit Network Connection (rev 01)
        Subsystem: Intel Corporation Ethernet Server Adapter I350-T2
        Kernel driver in use: vfio-pci
        Kernel modules: igb
02:00.1 Ethernet controller: Intel Corporation I350 Gigabit Network Connection (rev 01)
        Subsystem: Intel Corporation Ethernet Server Adapter I350-T2
        Kernel driver in use: igb
        Kernel modules: igb
```

# DPDK 的实现

> 参考：
>
> - [知乎，DPDK的整体工作原理](https://zhuanlan.zhihu.com/p/486288121)
> - https://doc.dpdk.org/guides/prog_guide/overview.html

在 [DPDK 官方文档 - 7.1](https://doc.dpdk.org/guides/linux_gsg/linux_drivers.html?highlight=vfio%20pci#binding-and-unbinding-network-ports-to-from-the-kernel-modules) 章节中，建议在所有情况下都使用 **vfio-pci** 作为 DPDK 绑定端口的内核模块。如果 IOMMU 不可用，则可以在 no-iommu 模式下使用 vfio-pci。如果由于某种原因 vfio 不可用，则可以使用基于 UIO 的模块 igb_uio 和 uio_pci_generic。详细信息请参见 UIO 部分。

## lcore & socket

**Logical core(逻辑核心)** DPDK 使用逻辑核心来管理并行处理，将任务分配给不同的逻辑核心以充分利用多核处理器的性能。通过使用逻辑核心，DPDK 应用程序可以更好地控制和利用系统中的处理器资源。狭义上 1 个 lcore 可以指 EAL 的 1 个线程。

这两个概念在 `DPDK`的代码中随处可见，**注意**这里的 **socket** 不是网络编程里面的那一套东西，而是 **CPU** 相关的东西。具体的概念可以参看[Differences between physical CPU vs logical CPU vs Core vs Thread vs Socket](https://link.zhihu.com/?target=https%3A//link.segmentfault.com/%3Fenc%3DkeRHFE71AwA4LEWK3gy%252F%252Bg%253D%253D.zL548YXL%252F1rT%252BQN9hhP0BSuhSZUszAZly2ULaOcHzihSAMFb3k6C8kBfLxFL65VqtdDdc0MigZcmMKHcWpmSUSTYHTTZ%252BgYz9XlsRQPKOK%252BWch%252FsoT6h%252BzR46e7YgN19TmdjGuy%252BwWL%252FfT2wdU6Q7Q%253D%253D) 或者其翻译版本[physical CPU vs logical CPU vs Core vs Thread vs Socket（翻译）](https://link.zhihu.com/?target=https%3A//link.segmentfault.com/%3Fenc%3DQvdG3%252BI75LlbmFlczo3WpQ%253D%253D.wu%252FAAZX7seAtEJjmuttDzYMu8zmGLKxX2fcFZ%252BZxKkIQvVQguoy2MSUnOVZPbaU6)。

对我们来说，只要知道可以`DPDK`可以运行在多个`lcore`上就足够了.

`DPDK` 如何知道有多少个`lcore`呢 ? 在启动时解析文件系统中的特定文件就可以了, 参考函数`eal_cpu_detected`

# DPDK 工具

https://doc.dpdk.org/guides/tools/index.html

## dpdk-devbind - 绑定和取消绑定设备与驱动程序，检查他们的状态

**dpdk-devbind \[OPTIONS] DEVICE1 DEVICE2 ....**

**OPTIONS**

- **-s, --status** # 打印所有已知网络接口的当前状态。对于每个设备，它显示 PCI 域、总线、插槽和功能，以及设备的文本描述。根据设备是否由内核驱动程序、vfio-pci 驱动程序或无驱动程序使用，将显示其他相关信息： - Linux 接口名称，例如if=eth0 - 正在使用的驱动程序，例如drv=vfio-pci - 当前未使用该设备的任何合适的驱动程序，例如used=vfio-pci 注意：如果此标志与绑定/取消绑定选项一起传递，则状态显示将始终在其他操作发生后发生。
- **-b, --bind DRIVER** # 选择绑定网卡要使用的驱动程序。可以使用 none 以解除绑定
- **-u, --unbind** # 接触网卡设备绑定。等价于 `-b none`
- **--force** # 强制绑定。默认情况下，若目标网卡已被内核启用（通常表现为已在路由表条目中），则无法被 DPDK 绑定。

EXAMPLE

从当前驱动程序绑定 eth1 并转而使用 vfio-pci：

- `dpdk-devbind --bind=vfio-pci eth1`

# DPDK 与 BPF 与 Netfilter

TODO
