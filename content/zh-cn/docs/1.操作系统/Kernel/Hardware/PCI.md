---
title: PCI
linkTitle: PCI
date: 2024-06-17T11:15
weight: 20
---

# 概述

> 参考:
>
> - [GitHub 项目，torvalds/linux - 通过 sysfs 访问 PCI 设备资源](https://github.com/torvalds/linux/blob/master/Documentation/PCI/sysfs-pci.rst)
> - [GitHub 项目，torvalds/linux - 如何编写 PCI 驱动](https://github.com/torvalds/linux/blob/master/Documentation/PCI/pci.rst)
>   - [博客园，如何编写Linux PCI驱动程序](https://www.cnblogs.com/wanglouxiaozi/p/15525726.html)
> - https://www.makelinux.net/ldd3/ - [12.1. The PCI Interface](https://www.makelinux.net/ldd3/chp-12-sect-1.shtml)

**[PCI](/docs/0.计算机/Motherboard/PCI.md) device resources(PCI 设备资源)** 由 Kernel 注册在 [sysfs](/docs/1.操作系统/Kernel/Filesystem/特殊文件系统/sysfs.md) 的 `/sys/devices/pci${DOMAIN:BUS}/` 目录。每个 PCI 设备资源在该目录下都有一个以 **唯一标识符(有的时候也称为 PCI Address)** 命名的目录，格式为: **`DOMAIN:BUS:SLOT.FUNC`**（e.g. 0000:17:00.0）

- **DOMAIN(域)** # 表示 PCI 域编号。用于识别主机系统中的不同 PCI 主机桥。在较早期的系统中，只有一个域编号为 0。随着系统规模扩大，可能存在多个 PCI 域。
- **BUS(总线)** # 表示 PCI 总线编号（16 进制）。一个 PCI 域中可能包含多条 PCI 总线，每条总线都有一个唯一编号。
- **SLOT(插槽)** # 表示 PCI 插槽编号。每条 PCI 总线上可以连接多个 PCI 设备，每个设备对应一个插槽编号。
  - 有的源码中也描述为 DEVICE。比如[这里](https://github.com/torvalds/linux/blob/master/Documentation/PCI/pci-iov-howto.rst)
- **FUNC(功能)** # 表示 PCI 功能编号。一些 PCI 设备可能包含多个独立的功能，每个功能都有一个编号,用于区分和访问。对于单功能设备,该编号通常为 0。
  - 有的源码中也描述为 FUNCTION

通常来说，目录可能是像这样的:

```bash
/sys/devices/pci0000:17
|-- 0000:17:00.0
|   |-- class
|   |-- config
|   |-- device
|   |-- enable
|   |-- irq
|   |-- local_cpus
|   |-- remove
|   |-- resource
|   |-- resource0
|   |-- resource1
|   |-- resource2
|   |-- revision
|   |-- rom
|   |-- subsystem_device
|   |-- subsystem_vendor
|   `-- vendor
`-- ...
```

上面这个目录的例子描述的 PCI 设备在 Domain(域) 编号为 0000，Bus(总线) 编号为 17（HEX）。该 Bus 在 Slot(插槽) 0 中包含了 1 个单一功能的 device(设备)。

## PCI 设备资源信息

> 参考:
>
> - [GitHub 项目，torvalds/linux - include/linux/pci_ids.h](https://github.com/torvalds/linux/blob/master/include/linux/pci_ids.h) 定义了所有 PCI 相关的 Class、Vendor、Device 的 ID。
>   - 比如该目录下的 class、device、vendor、etc. 文件中的值的含义，可以从这里找到定义
> - https://github.com/torvalds/linux/blob/master/Documentation/PCI/sysfs-pci.rst

**PCI device resouorces(PCI 设备的资源)** 数据以文件形式保存。*PCI 设备资源信息* 的储存目录中，通常包含如下文件

| 文件名                | 功能                                                                                      |
| ------------------ | --------------------------------------------------------------------------------------- |
| class              | PCI 设备的类型 ID (ascii, ro)。比如 网卡、USB、etc.                                                 |
| device             | PCI 设备的型号 ID (ascii, ro)。类似设备型号，比如 I350 Gigabit Network 这种。设备 ID 基于供应商 ID，是在供应商 ID 之下的。 |
| vendor             | PCI 设备的供应商  ID (ascii, ro)。各种厂家，比如 Intel、AMD、etc.                                       |
| subsystem_device   | PCI 设备的型号子系统 ID (ascii, ro)                                                             |
| subsystem_vendor   | PCI 设备的供应商子系统 ID (ascii, ro)                                                            |
| config             | PCI config space (binary, rw)                                                           |
| enable             | Whether the device is enabled (ascii, rw)                                               |
| irq                | IRQ number (ascii, ro)                                                                  |
| local_cpus         | nearby CPU mask (cpumask, ro)                                                           |
| remove             | remove device from kernel's list (ascii, wo)                                            |
| resource           | PCI resource host addresses (ascii, ro)                                                 |
| resource0..N       | PCI resource N, if present (binary, mmap, rw)                                           |
| resource0_wc..N_wc | PCI WC map resource N, if prefetchable (binary, mmap)                                   |
| revision           | PCI revision (ascii, ro)                                                                |
| rom                | PCI ROM resource, if present (binary, ro)                                               |

该目录中主要由如下几类文件

- **ro(只读) 类型文件** # **是信息性的**，对它们的写入将被忽略，“rom”文件除外。可写文件可用于在设备上执行操作（例如更改配置空间、分离设备）。 mmapable 文件可通过偏移量 0 处的文件 mmap 获得，并可用于从用户空间进行实际设备编程。请注意，某些平台不支持某些资源的映射，因此请务必检查任何尝试的映射的返回值。其中最值得注意的是 I/O 端口资源，它也提供读/写访问
  - **class、vendor、device、subsystem_vendor、subsystem_device** # PCI 设备信息文件，文件内容是 **HEX(16 进制)** 的数字（i.e. 0x8086，0x 是 16 进制标识，具体数值是 8086），数字对应的含义可以参考 [GitHub 项目，torvalds/linux - include/linux/pci_ids.h](https://github.com/torvalds/linux/blob/master/include/linux/pci_ids.h)
    - 这些文件中的值分别对应 CLASS_ID、VENDOR_ID、DEVICE_ID、SBUVENDOR_ID、SUBDEVICE_ID
    - **Tips**: DEVICE_ID 都是包含在 VENDOR_ID 下的，理解为某个供应商下的某个设备比如 1521 这个 DEVICE_ID，可以在 8086 这个 VENDOR_ID 下找到，表示英特尔公司的 1350 Gigabit Network Connection 型号的网卡
- **enable 文件** # 提供了一个计数器，指示设备已启用的次数。如果“enable”文件当前返回“4”，并且回显“1”，则它将返回“5”。回显“0”会减少计数。但是，即使返回到 0，某些初始化也可能无法逆转。
- **rom 文件** # 的特殊之处在于它提供对设备 ROM 文件（如果可用）的只读访问。但默认情况下它是禁用的，因此应用程序应在尝试读取调用之前将字符串“1”写入文件以启用它，并在访问后通过将“0”写入文件来禁用它。请注意，必须启用设备才能读取 ROM 才能成功返回数据。如果驱动程序未绑定到设备，则可以使用上面记录的“enable”文件来启用它。
- **remove 文件** # 用于通过向文件写入非零整数来删除 PCI 设备。这不涉及任何类型的热插拔功能，例如关闭设备电源。该设备将从内核的 PCI 设备列表中删除，其 sysfs 目录也将被删除，并且该设备将从附加到它的任何驱动程序中删除。不允许移除 PCI 根总线。

### 目录中的其他文件

**uevent** # userspace event(用户空间事件，简称 uevent)

# PCI 关联文件

> [!Question]
> 如果有 pci.ids 文件，为啥还要自己定义 class、device、vendor、etc. 这些文件中的数字对应关系呢？是通过代码中的定义运行更效率？但是代码中的定义好像少。
>
> 看过下文之后，假如 DEVICE_ID 是 1521，在 pid.ids 文件中 https://admin.pci-ids.ucw.cz/read/PC/8086/1521, 但是 Linux 源码 https://github.com/torvalds/linux/blob/master/include/linux/pci_ids.h 这里找不到对应 DEVICE_ID 是 1521。也就是说有部分 PCI 设备信息在 <font color="#ff0000">pci.ids 文件中有，但是 Linxu 代码中没有</font>。

**pci.ids** # [pci.ids 文件](#pci.ids%20文件)，包含所有已知 PCI ID（vendors(供应商)、devices(设备)、classes(类别)、subclasses(子类别)）列表。可以称为 **PCI ID 存储库**。

- 不同的 Linux 发行版该文件所在位置不同，通常都是 `/usr/share/XXX/pci.ids`

**.pciids-cache** # 通过 DNS 查询模式获取到的所有 PCI ID 缓存文件

**pci_ids.h** # [Linux 源码](https://github.com/torvalds/linux/blob/master/include/linux/pci_ids.h)中的 PCI ID 与 名称 的对应关系

**`/sys/devices/pci${DOMAIN:BUS}/${DOMAIN:BUS:SLOT.FUNC}/`** # [PCI 设备资源信息](#PCI%20设备资源信息)储存目录。包含该 PCI 设备的各种信息，e.g. 供应商、PCI 类型、etc.

- Note: 由于 PCI 是有层级的（e.g. `/sys/devices/pci0000:00/0000:00:1c.1/0000:05:00.0/0000:06:00.0/0000:07:00.0/0000:08:00.0`），所以每个层级的目录都有可能存在下面这些内容
- **/sys/bus/pci/devices/** # 可以在这里找到所有 PCI 设备的列表，该目录下的文件是以 PCI Addr 命名的软链接，指向 `/sys/devices/pci${DOMAIN:BUS}/` 目录下的某个子目录。
  - Notes: 有的 PCI 层级可能比较多，比如这种 `0000:08:00.0 -> ../../../devices/pci0000:00/0000:00:1c.1/0000:05:00.0/0000:06:00.0/0000:07:00.0/0000:08:00.0`

# pci.ids 文件

> 参考:
>
> - [官网](https://pci-ids.ucw.cz/)
>   - [GitHub 项目，pciutils/pciids](https://github.com/pciutils/pciids)
> - [Manual(手册)，pci.ids(5)](https://man7.org/linux/man-pages/man5/pci.ids.5.html)

pci.ids 是 PCI 设备中使用的所有已知 ID 的公共存储库：供应商、设备、子系统和设备类别的 ID。它在各种程序（例如 [PCI 实用程序](http://mj.ucw.cz/sw/pciutils/)）中用于显示完整的人类可读名称，而不是神秘的数字代码。

> Note: pci.ids 数字都是 **HEX(16 进制)** 的

该文件由 [pciutils](https://github.com/pciutils) 维护。在官方网站中可以在线查询，主要分为两大块查询

- 查询 [PCI Device ID](https://admin.pci-ids.ucw.cz/read/PC/)
  - 每个 Device 都在某个 Vendor 下，比如 1521 这个 DEVICE_ID，可以在 8086 这个 VENDOR_ID 下找到，表示英特尔公司的 1350 Gigabit Network Connection 型号的网卡
  - 使用在线库查询时，先找到 VENDOR ID 点进去，再找 DEVICE ID，即可找到 设备的型号命名。
- 查询 [PCI Device 类型 ID](https://admin.pci-ids.ucw.cz/read/PD/)

> [!Question]
> 该文件的内容如何生成的？比如 /sys/devices/pci0000:00/0000:00:00.0/vendor 文件的值怎么来的？从 [include/linux/pci.h](https://github.com/torvalds/linux/blob/master/include/linux/pci.h) 文件的提示中 `For more information, please consult the following manuals (look at http://www.pcisig.com/ for how to get them)` 合理猜测这些数据都是从 PCI SIG 定义的规范中获取到的？

下面的维护和更新由 AI 回答，待验证

## 维护和更新

**维护机构**：

- `pci.ids` 文件由 `pciutils` 项目维护，该项目的开发者和社区成员定期更新这个文件。
- `pciutils` 项目由 Martin Mares 维护，托管在 `https://github.com/pciutils/pciids`。

**数据来源**：

- 新的供应商 ID 和设备 ID 通常由硬件制造商提供。
- 制造商向 [PCI-SIG](/docs/Standard/IT/PCI-SIG.md) 注册他们的 ID。PCI-SIG 是管理和分配 PCI ID 的机构。
- 社区成员和开发者通过手动添加新设备的信息到 `pci.ids` 文件中。

**更新过程**：

- 新的设备信息通过提交补丁的方式添加到 `pciutils` 项目的仓库中。
- 开发者在 GitHub 上提交 Pull Request 或直接发送补丁给维护者。
- 维护者审核并合并这些更改。
- 合并后的更新会在 `pciutils` 项目的新版本中发布。

# PCI I/O 虚拟化

https://github.com/torvalds/linux/blob/master/Documentation/PCI/pci-iov-howto.rst

用 [DPDK](/docs/4.数据通信/DPDK/DPDK.md) 举例

```bash
~]# lshw -c net -businfo
Bus info          Device     Class          Description
=======================================================
pci@0000:02:00.0             network        Ethernet 10G 2P X520 Adapter
pci@0000:02:00.1             network        Ethernet 10G 2P X520 Adapter
pci@0000:01:00.0  eno3       network        I350 Gigabit Network Connection
pci@0000:01:00.1  eno4       network        I350 Gigabit Network Connection
~]#
~]# dpdk-devbind.py -s

Network devices using DPDK-compatible driver
============================================
0000:02:00.0 'Ethernet 10G 2P X520 Adapter 154d' drv=igb_uio unused=ixgbe,vfio-pci
0000:02:00.1 'Ethernet 10G 2P X520 Adapter 154d' drv=igb_uio unused=ixgbe,vfio-pci

Network devices using kernel driver
===================================
0000:01:00.0 'I350 Gigabit Network Connection 1521' if=eno3 drv=igb unused=igb_uio,vfio-pci
0000:01:00.1 'I350 Gigabit Network Connection 1521' if=eno4 drv=igb unused=igb_uio,vfio-pci
~]#
~]# find / -name 'sriov_numvfs'
/sys/devices/pci0000:00/0000:00:02.0/0000:02:00.0/sriov_numvfs
/sys/devices/pci0000:00/0000:00:02.0/0000:02:00.1/sriov_numvfs
```

通过 DPDK 管理的 PCI 设备可以找到 sriov_numvfs 文件

# PCI 源码与实现规范

https://github.com/torvalds/linux/blob/v6.9/include/linux/pci.h#L322 - `struct pci_dev` 是存储 PCI 信息的结构体

`include/linux/pci.h` 文件的开头注释中有这么一段

```c
 *	For more information, please consult the following manuals (look at
 *	http://www.pcisig.com/ for how to get them):
 *
 *	PCI BIOS Specification
 *	PCI Local Bus Specification
 *	PCI to PCI Bridge Specification
 *	PCI Express Specification
 *	PCI System Design Guide
```

可以查查这些规范里都有什么

- PCI BIOS Specification
- PCI Local Bus Specification
- PCI to PCI Bridge Specification
- PCI Express Specification
- PCI System Design Guide

# /sys/bus/ 下的 PCI

> 参考:
>
> - https://github.com/torvalds/linux/blob/master/Documentation/ABI/testing/sysfs-bus-pci
> - https://www.kernel.org/doc/html/latest/admin-guide/abi-testing.html#abi-file-testing-sysfs-bus-pci

**/sys/bus/pci/drivers/** # 管理 PCI 驱动。详见 [Driver](/docs/1.操作系统/Kernel/Hardware/Driver.md#PCI) 的 PCI 章节

**/sys/bus/pci/devices/** # 目录下是所有 PCI 设备指向 /sys/devices/ 的 [Symbolic link](/docs/1.操作系统/Kernel/Filesystem/文件管理/Symbolic%20link.md)。

**/sys/bus/pci/drivers_autoprobe** #

**/sys/bus/pci/drivers_probe** #

**/sys/bus/pci/rescan** # 向该文件写入非零值将强制重新扫描系统中的所有 PCI 总线，并重新发现以前删除的设备。

**/sys/bus/pci/resource_alignment** #

**/sys/bus/pci/slots/** #

**/sys/bus/pci/uevent** #

# TODO: WSL 中的 网卡 PCI 总结
