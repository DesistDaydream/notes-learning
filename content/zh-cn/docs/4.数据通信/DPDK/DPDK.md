---
title: DPDK
linkTitle: DPDK
weight: 1
data: 2024-03-01T00:00:00
---

# 概述

> 参考：
>
> - [GitHub 项目，DPDK/dpdk](https://github.com/DPDK/dpdk)
> - [官网](https://www.dpdk.org/)
> - [官方文档，API](https://doc.dpdk.org/api/)
> - [Wiki, Data Plane Development Kit](https://en.wikipedia.org/wiki/Data_Plane_Development_Kit)
> - [DPDK 开发中文网](https://sdn.0voice.com/)

**Data Plane Development Kit(数据平面开发套件，简称 DPDK)** 是一个由 [Linux 基金会](/docs/Standard/Foundation/Linux%20Foundation.md) 管理的开源软件项目。用于将 TCP 数据包的处理能力从内核空间移动到用户空间中的进程。主要是跳过了内核的 [Interrupts(中断)](/docs/1.操作系统/Kernel/CPU/Interrupts(中断)/Interrupts(中断).md) 逻辑。

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

[DPDK Library](/docs/4.数据通信/DPDK/DPDK%20Library.md)

## Driver

Linux 驱动

- https://doc.dpdk.org/guides/linux_gsg/linux_drivers.html#binding-and-unbinding-network-ports-to-from-the-kernel-modules
- **vfio** # 使用 [VFIO](/docs/1.操作系统/Kernel/VFIO.md) 功能的驱动。依赖 `vfio-pci` 模块。[官方文档](https://doc.dpdk.org/guides/linux_gsg/linux_drivers.html#binding-and-unbinding-network-ports-to-from-the-kernel-modules)建议所有情况下都是用 vfio-pci 作为 DPDK 绑定端口的内核模块

TODO: 好多好多驱动。

## DPDK APP

https://doc.dpdk.org/guides/tools/index.html

# DPDK 部署

> 参考:
>
> - [官方文档，快速开始](https://core.dpdk.org/doc/quick-start/)
> - [知乎，KVM虚拟机安装DPDK(vfio+igb_uio)](https://zhuanlan.zhihu.com/p/613149235) TODO: 待验证

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

此时想要查看网卡当前所使用的驱动，可以使用 [Linux 硬件管理工具](/docs/1.操作系统/Linux%20管理/Linux%20硬件管理工具/Linux%20硬件管理工具.md) 查看到当前网卡所使用的驱动程序（lspci、lshw、etc.）。（虽然其中显示的是 Kernel driver in use，但是实际上，驱动是 vfio-pci 的网卡已经被 DPDK 接管了）

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

这两个概念在 `DPDK`的代码中随处可见，**注意**：这里的 **socket** 不是网络编程里面的那一套东西，而是 **CPU** 相关的东西。具体的概念可以参看 [Differences between physical CPU vs logical CPU vs Core vs Thread vs Socket](https://www.daniloaz.com/en/differences-between-physical-cpu-vs-logical-cpu-vs-core-vs-thread-vs-socket/) 或者其翻译版本 [physical CPU vs logical CPU vs Core vs Thread vs Socket（翻译）](https://www.cnblogs.com/zh1164/p/9883852.html)。

对我们来说，只要知道 `DPDK` 可以运行在多个 `lcore` 上就足够了.

`DPDK` 如何知道有多少个 `lcore` 呢 ? 在启动时解析文件系统中的特定文件就可以了, 参考函数 `eal_cpu_detected`

## 内存的使用

TODO: https://xie.infoq.cn/article/a7c83189a19387018b8595e98

调用 DPDK 的程序自己管理如何使用内存，需要提前预留一块空间给 DPDK 程序自身，并保证不被其他程序使用，通常这块内存空间的的分页是 1GiB、512MiB 等大小的 [Huge Pages](/docs/1.操作系统/Kernel/Memory/Huge%20Pages.md)。

e.g. 一个 256 GiB 的系统，可以留出来 192 个 1GiB 页大小的内存空间，总计 192 GiB，这一部分内存空间就是大页空间，因为每个分页的大小高达 1GiB

# DPDK 工具

https://doc.dpdk.org/guides/tools/index.html

## dpdk-devbind - 绑定和取消绑定设备与驱动程序，检查他们的状态

https://github.com/DPDK/dpdk/blob/main/usertools/dpdk-devbind.py

https://doc.dpdk.org/guides/tools/devbind.html

### Syntax(语法)

**dpdk-devbind \[OPTIONS] DEVICE1 DEVICE2 ....**

**OPTIONS**

- **-s, --status** # 打印所有已知网络接口的当前状态。对于每个设备，它显示 PCI 域、总线、插槽和功能，以及设备的文本描述。根据设备是否由内核驱动程序、vfio-pci 驱动程序或无驱动程序使用，将显示其他相关信息：
  - Linux 接口名称，例如 if=eth0
  - 正在使用的驱动程序，例如drv=vfio-pci
  - 当前未使用该设备的任何合适的驱动程序，例如used=vfio-pci 注意：如果此标志与绑定/取消绑定选项一起传递，则状态显示将始终在其他操作发生后发生。
- **-b, --bind DRIVER** # 选择绑定网卡要使用的驱动程序。可以使用 none 以解除绑定
- **-u, --unbind** # 接触网卡设备绑定。等价于 `-b none`
- **--force** # 强制绑定。默认情况下，若目标网卡已被内核启用（通常表现为已在路由表条目中），则无法被 DPDK 绑定。

-s 选项输出内容说明，包含如下几部分

- Network devices using DPDK-compatible driver # 由 DPDK 管理的网卡
- Network devices using kernel driver # 由内核管理的网卡
- Other Network devices # 既不由 DPDK 管理也不由内核管理的网卡。（使用 -u 解绑后的网卡）

### EXAMPLE

从当前驱动程序绑定 eth1 并转而使用 vfio-pci：

- `dpdk-devbind --bind=vfio-pci eth1`

## dpdk-telemetry.py

https://github.com/DPDK/dpdk/blob/main/usertools/dpdk-telemetry.py

用于通过 DPDK Telemetry Library 暴露的 **SOCK_SEQPACKET** 类型的 Unix Socket（通常为 `/var/run/dpdk/*/dpdk_telemetry.v2`） 查询 DPDK 中的遥测信息。目前包括 ethdev 状态、ethdev 端口列表、eal 参数、etc.。

`/help,COMMAND` # 获取命令的帮助信息

# DPDK 与 BPF 与 Netfilter

TODO

# FAQ

https://doc.dpdk.org/guides/faq/index.html

## 将 DPDK 管理的网卡还给内核

一、查看绑定的网卡

```bash
dpdk-devbind.py -s
```

二、将 `Network devices using DPDK-compatible driver` 部分的网卡从 DPDK 中解绑

> Note: 多个网卡以空格分割。同时记录下

```bash
dpdk-devbind.py -u 0000:21:00.0 0000:21:00.1
```

解绑后，通过 `dpdk-devbind.py -s` 可以在 `Other Network devices` 部分查看到未被 Kernel 和 DPDK 管理的网卡。其中最后一部分 `unused=i40e,igb_uio,vfio-pci` unused 的值的第一段（逗号分割）可以作为内核所使用的驱动（e.g. i40e）

三、将 `Other Network devices` 部分的网卡绑定给内核

> Note: 这里使用上面查看到的 i40e 作为驱动示例

```bash
dpdk-devbind.py -b i40e 0000:21:00.0 0000:21:00.1
```

四、总结

这里面的绑定与解绑从 dpdk-devbind.py 源码中可以看到，主要依赖于 [sysfs](/docs/1.操作系统/Kernel/Filesystem/特殊文件系统/sysfs.md) 中 PCI 的 [Driver](/docs/1.操作系统/Kernel/Hardware/Driver.md#PCI) 管理能力。i.e. /sys/bus/pci/drivers/.../bind 与 /sys/bus/pci/drivers/.../unbind 文件

## CPU 使用率 100%

DPDK 的 UIO 驱动屏蔽了硬件发出中断，然后在用户态采用主动轮询的方式，这种模式被称为 [PMD](http://doc.dpdk.org/guides/prog_guide/poll_mode_drv.html)（Poll Mode Driver）。

UIO 旁路了内核，主动轮询去掉硬中断，DPDK 从而可以在用户态做收发包处理。带来 Zero Copy、无系统调用的好处，同步处理减少上下文切换带来的 Cache Miss。

运行在 PMD 的 Core 会处于用户态 CPU100% 的状态

![](https://ask.qcloudimg.com/draft/1141755/v43qilsryd.png?imageView2/2/w/1620)

网络空闲时CPU长期空转，会带来能耗问题。所以，DPDK推出Interrupt DPDK模式。

![](https://ask.qcloudimg.com/draft/1141755/pg3d428gpr.png?imageView2/2/w/1620)

它的原理和 [NAPI](https://www.ibm.com/developerworks/cn/linux/l-napi/index.html) 很像，就是没包可处理时进入睡眠，改为中断通知。并且可以和其他进程共享同个CPU Core，但是DPDK进程会有更高调度优先级。
