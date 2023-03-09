---
title: KVM/QEMU
weight: 1
---

# 概述

> 参考：
> - [KVM 官网](https://www.linux-kvm.org/page/Main_Page)
> - [Ubuntu 官方文档，虚拟化-qemu](https://ubuntu.com/server/docs/virtualization-qemu)

## KVM 背景

**Kernel-based Virtual Machine(基于内核的虚拟化机器，简称 KVM)**， 是 Linux 的一个内核模块，就叫 **kvm**，只用于管理虚拟 CPU 和内存。该内核模块使得 Linux 变成了一个 Hypervisor。

- 它由 Quramnet 开发，该公司于 2008 年**被 Red Hat 收购**。
- 它支持 x86 (32 and 64 位), s390, Powerpc 等 CPU。
- 它从 Linux 2.6.20 起就作为一模块被包含在 Linux 内核中。
- 它需要支持虚拟化扩展的 CPU。
- 它是完全开源的。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zuowkm/1616124035086-2c826a6e-2fd2-402b-babd-06bfe2380e3d.png)

KVM 实际是 Linux 内核提供的虚拟化架构，可将内核直接充当 Hypervisor 来使用。KVM 需要宿主机的 CPU 本身支持虚拟化扩展，如 intel VT 和 AMD AMD-V 技术。KVM 自 2.6.20 版本后已合入主干并发行。除了支持 x86 的处理器，同时也支持 S/390,PowerPC,IA-61 以及 ARM 等平台。

KVM 包含包含两个内核模块

1. kvm 用来实现核心虚拟化功能
2. kvm-intel # 与处理器强相关的模块

KVM 本身只提供了 CPU 和 Memory 的虚拟化，并暴露了一个 **`/dev/kvm`** 设备，以供宿主机上的用户空间的程序访问(比如 下文提到的 QEMU)。用户空间的程序通过 **/dev/kvm** 接口可以实现多种功能

## QEMU 背景

> 参考：
>
> - [QEMU 官网](https://www.qemu.org/)

QEMU 是一个通过软件实现的完全虚拟化程序，通过动态二进制转换来模拟 CPU，并模拟一系列的硬件，使虚拟机认为自己和硬件直接打交道，其实是同 QEMU 模拟出来的硬件打交道，QEMU 再将这些指令翻译给真正硬件进行操作。通过这种模式，虚拟机可以和主机上的硬盘，网卡，CPU，CD-ROM，音频设备和 USB 设备进行交互。但由于所有指令都需要经过 QEMU 来翻译，因而性能会比较差

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zuowkm/1616124035102-78899618-45f3-4dcc-9de1-9c80ecd532cb.jpeg)

## KVM/QEMU 诞生

![图片来源：RedHat8 Virtualization Architecture 在 KVM/QEMU Storage Stack Performance Discussion 这篇文章中，作者还画了一个非常形象的图，可以作为参考，下面缩小的图就是](https://notes-learning.oss-cn-beijing.aliyuncs.com/zuowkm/1616124035098-45602829-f8bc-4f49-b65d-56b8dba6c466.png "图片来源：RedHat8 Virtualization Architecture 在 KVM/QEMU Storage Stack Performance Discussion 这篇文章中，作者还画了一个非常形象的图，可以作为参考，下面缩小的图就是")

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zuowkm/1616124035057-cdeb1319-ee83-4674-99c3-70a16da96211.jpeg)

从前面的背景介绍可知，KVM 实现了 CPU 和 Memory 的虚拟化，但 KVM 并不能模拟其他设备，所以需要其他东西来支持其他设备的模拟；而 QEMU 是通过纯软件实现的一套完整的虚拟化，但是性能非常低下。所以 KVM 与 QEMU 天然得相辅相成，KVM 的开发者选择了比较成熟的开源虚拟化软件 QEMU 来模拟 I/O 设备(网卡，磁盘等)，最后形成了 KVM/QEMU。

在 KVM/QEMU 中，KVM 运行在内核空间，QEMU 运行在用户空间，实际模拟创建、管理各种虚拟硬件，QEMU 将 KVM 整合了进来，通过 ioctl() 系统调用来调用 /dev/kvm 设备，从而将 CPU 指令的部分交给内核模块来做，KVM 实现了 CPU 和 Memory 的虚拟化，QEMU 模拟 IO 设备(磁盘，网卡，显卡等)，KVM 加上 QEMU 后就是完整意义上的服务器虚拟化。

综上所述，QEMU-KVM 具有两大作用：

- KVM 负责 cpu，内存 的虚拟
- QEMU 负责 I/O 设备 的模拟。比如显卡、PCI、USB、声卡、网卡、存储设备等等。

## 结语

QEMU-KVM，是 QEMU 的一个特定于 KVM 加速模块的分支，里面包含了很多关于 KVM 的特定代码，与 KVM 模块一起配合使用。

目前 QEMU-KVM 已经与 QEMU 合二为一，所有特定于 KVM 的代码也都合入了 QEMU，当需要与 KVM 模块配合使用的时候，只需要在 QEMU 命令行加上 --enable-kvm 就可以。

# KVM/QEUM 虚拟化实现原理

KVM/QEMU 主要通过以下组件来实现完整的虚拟化功能

- **kvm.ko.xz** # kvm 内核模块。用来模拟 CPU 与 RAM。
- **/dev/kvm** # 一个字符设备(也是一个接口)。供用户空间的程序使用 `ioctl()` 系统调用来访问 kvm 模块
- **qemu-system** # 一个二进制文件。用来调用 `/dev/kvm` 设备，并为虚拟机模拟各种 I/O 设备。qemu-system 也是最基本的用于创建虚拟机的命令行工具。

KVM/QEMU 通过 [qemu-img 命令行工具](docs/IT学习笔记/10.云原生/1.2.实现虚拟化的工具/KVM_QEMU/KVM_QEMU%20命令行工具/qemu-img.md)来创建 VM 文件。[qemu-system](docs/IT学习笔记/10.云原生/1.2.实现虚拟化的工具/KVM_QEMU/KVM_QEMU%20命令行工具/qemu-system.md) 程序使用 /dev/kvm 接口来调用 kvm 模块，以运行 VM。qemu-kvm 和 qemu-system 两个工具也是创建与运行 VM 的最基础工具。

使用 /dev/kvm 接口的 qemu-system 程序可以提供如下能力：

- 设置 VM 的地址空间。宿主机必须提供固件镜像(通常为模拟出来的 BIOS)以便让 VM 可以引导到 操作系统中
- 为 VM 模拟 I/O 设备。
- 将 VM 的视频显示映射回宿主机上。

## Virtualization CPU && Memory

### CPU 虚拟化

QEMU 创建 CPU 线程，在初始化的时候设置好相应的虚拟 CPU 寄存器的值，然后调用 KVM 的接口，运行虚拟机，在物理 CPU 上执行虚拟机代码。

在虚拟机运行时，KVM 会截获虚拟机中的敏感指令，当虚拟机中的代码是敏感指令或者满足了一定的退出条件时，CPU 会从 VMX non-root 模式退出到 KVM，这就是下图的 VM exit。虚拟机的退出首先陷入到 KVM 进行处理，但是如果遇到 KVM 无法处理的事件，比如虚拟机写了设备的寄存器地址，那么 KVM 就会将这个操作交给 QEMU 处理。当 QEMU/KVM 处理好了退出事件后，又会将 CPU 置于 VMX non-root 模式，也就是下图的 VM Entry。

KVM 使用 VMCS 结构来保存 VM Exit 和 VM Entry
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zuowkm/1616124035076-0336490e-7922-482c-91be-ee6b0a2ba562.png)

### Memory 虚拟化

QEMU 初始化时调用 KVM 接口告知 KVM，虚拟机所需要的物理内存，通过 mmap 分配宿主机的虚拟内存空间作为虚拟机的物理内存，QEMU 在更新内存布局时会持续调用 KVM 通知内核 KVM 模块虚拟机的内存分布。

在 CPU 支持 EPT（拓展页表）后，CPU 会自动完成**虚拟机物理地址**到**宿主机物理地址**的转换。虚拟机第一次访问内存的时候会陷入 KVM，KVM 逐渐建立起 EPT 页面。这样后续的虚拟机的虚拟 CPU 访问虚拟机**虚拟内存地址**时，会先被转换为**虚拟机物理地址**，接着查找 EPT 表，获取宿主机物理地址
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zuowkm/1616124035074-4de2d638-b8fb-499c-a1de-92ce5e6a10b3.png)

## Paravirtualized Devices(半虚拟化设备)

> 参考：
> - [官方文档，系统模拟-设备模拟](https://www.qemu.org/docs/master/system/device-emulation.html)
> - [RedHat7 虚拟化硬件设备章节中的半虚拟化章节](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/virtualization_getting_started_guide/sec-virtualization_getting_started-products-virtualized-hardware-devices#sec-Virtualization_Getting_Started-Products-paravirtdevices)

半虚拟化设备，就是 Qemu 模拟的各种 I/O 设备

在 QEMU/KVM 早期模拟其他的硬件(如存储、网络设备)性能不足。为了提高 IO 设备性能，所以产生了 **Paravirtualized Devices(半虚拟化设备)**，Paravirtualized(半虚拟化) 为 VM 使用宿主机上的设备提供了**快速且高效的通讯方式**。KVM/QEMU 使用 **Virtio API** 作为 VM 与 Hypervisor 的中间层，以便为 VM 提供 Paravirtualized Devices(半虚拟化设备)。

> 一些半虚拟化设备可以有效减少 I/O 的延迟，并把 I/O 的吞吐量提高至接近裸机的水平。

所有 **Virtio** 设备都**由两部分组成**：

- **Host Device** # 宿主机设备
  - 也称为 **Device Back End(后端设备)**
- **Guest Device** # 虚拟机设备
  - 也称为 **Device Front End(前端设备)**

Paravirtualizd device driver(半虚拟化设备驱动) 可以让 VM 直接访问宿主机上的物理硬件设备。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zuowkm/1616124035106-08438d88-b937-43b9-af1c-01ddb6771941.jpeg)
现阶段有多种半虚拟化设备可供使用

- virtio-net(半虚拟化网络设备) # 半虚拟化网络设备是一种虚拟网络设备，可通过增加的 I/O 性能和较低的延迟为虚拟机提供网络访问。
- virtio-blk(半虚拟化块设备) # 半虚拟化块设备是一种高性能的虚拟存储设备，可为虚拟机提供更高的 I / O 性能和更低的延迟。 虚拟机管理程序支持半虚拟化的块设备，该设备已连接到虚拟机（必须仿真的软盘驱动器除外）。
- virtio-scsi(半虚拟化控制器设备) # 半虚拟化 SCSI 控制器设备是一种更为灵活且可扩展的 virtio-blk 替代品。virtio-scsi 客机能继承目标设备的各种特征，并且能操作几百个设备，相比之下，virtio-blk 仅能处理 28 台设备。
- 半虚拟化时钟
- virtio-serial(半虚拟化串口设备) #
- virtio-ballon(气球设备) # 气球（ballon）设备可以指定虚拟机的部分内存为没有被使用（这个过程被称为气球“_充气_ ” — inflation），从而使这部分内存可以被主机（或主机上的其它虚拟机）使用。当虚拟机这部分内存时，气球可以进行“_放气_ ”（deflated），主机就会把这部分内存重新分配给虚拟机。
- virtio-rng(半虚拟化随机数生成器)
- QXL(半虚拟化显卡) # 半虚拟化显卡与 QXL 驱动一同提供了一个有效地显示来自远程主机的虚拟机图形界面。SPICE 需要 QXL 驱动。

### 创建虚拟机示例

qemu-kvm 命令中的 `-device` 选项用于指定前端设备，比如 网卡、磁盘、usb 等等。而 `-XXXXdev`等选项则是为了指定宿后端设备。比如宿主机上的文件、socket 等等

qemu-kvm 使用 `-device` 选项指定的参数将这些模拟出来的硬件设备，通过 ID 关联到`-XXXdev`定的宿主机文件上。

而在新版，则使用了更简单的方式，通过一个选项，来直接指定半虚拟化的两端设备，比如使用 `-drive` 代替 `-device` 和 `-blockdev`、使用 `-nic` 代替 `-device` 和 `-netdev` 等等。

比如下面的示例：

```bash
qemu-system-x86_64 -m 4096 -smp 2 -name test \
# 模拟块设备
# -drive 使新版选项，是 -blockdev 和 -device 两个参数的集合，可以模拟一个块设备
# Host Device 为 /var/lib/libvirt/images/test-2.bj-net.qcow2
# Guest Device 为 virtio-blk 设备
-drive file=/var/lib/libvirt/images/test-2.bj-net.qcow2,format=qcow2,if=virtio \
-vnc :3 \
# 模拟网卡
# -netdev 指定 Host Device(宿主机设备) 为 tap 类型的网络设备
# -device 指定 Guest Device(虚拟机设备) 为 virtio-net 设备
-netdev tap,id=n1 \
-device virtio-net,netdev=n1 \
# 模拟串口
# -chardev 指定 Host Device(宿主机设备)为 socket
# -device 指定 Guest Device(虚拟机设备)为 virtio-serial
# virtserialport 是 virtio-serial-port 的意思
-chardev socket,path=/tmp/qga.sock,server,nowait,id=qga0 \
-device virtio-serial \
-device virtserialport,chardev=qga0,name=org.qemu.guest_agent.0
```

可以看到，所有通过 -device 选项在 VM 中模拟的硬件设备，都会根据 ID 关联到宿主机的某个文件或者设备上。

而且还有一个 -drive 选项这种更简单的使用方式，来免去设定 ID 的困扰，并且输入的字符更少。

## 总结一下

KVM/QEMU 虚拟化环境中，除了 CPU 与 Memory 是通过 KVM 虚拟化的，其他所有硬件设备，都是通过 QEMU 模拟出来，并且，要想让模拟出来的硬件设备能正常工作(模拟的硬件与宿主机交互)，则还需要在宿主机上创建与之关联的文件。所以，一共两部分来实现 QEMU 的模拟功能。

1. **一部分是 QEMU 在 VM 中模拟出来的各种硬件**
2. **另一部分是在宿主机中与 VM 中模拟出来的硬件对应的各种文件或设备**。

这两部分共同实现了 VM 中模拟的硬件与宿主机交互的能力。如果 VM 中的硬件与宿主机无法交互，那么是无法使用滴~~~

# 不同 CPU 架构的模拟

> 参考：
> - [官方文档，系统模拟-QEMU 系统模拟目标](https://www.qemu.org/docs/master/system/targets.html)

QEMU 是一个通用的模拟器，它可以模拟很多机器。qemu-system 工具的大多数选项对于所有机器都是相似的。以下各节将提及有关各种目标的具体信息。

[Arm 系统模拟](https://www.qemu.org/docs/master/system/target-arm.html)

## TODO: 不同 CPU 架构之间的互相模拟

这个可以实现吗？amd 机器上开 arm 的虚拟机，或者 arm 机器上开 amd 的虚拟机