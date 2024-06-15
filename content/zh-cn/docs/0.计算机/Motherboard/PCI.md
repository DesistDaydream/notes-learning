---
title: PCI
linkTitle: PCI
date: 2024-06-15T09:24
weight: 20
---

# 概述

> 参考：
>
> - [Wiki，Peripheral_Component_Interconnect](https://en.wikipedia.org/wiki/Peripheral_Component_Interconnect)

**Peripheral Component Interconnect(外围组件互连，简称 PCI)**

## PCI-E

> 参考:
>
> - [Wiki，PCI_Express](https://en.wikipedia.org/wiki/PCI_Express)

**Peripheral Component Interconnect Express(简称 PCI-E)**


# PCI 规范

> 参考:
>
> - [pci设备身份识别码介绍说明](https://www.twblogs.net/a/5eee10c7264079afec950f51)

由 [PCI-SIG](docs/Standard/IT/PCI-SIG.md) 制定规范。一个 PCI 设备的通常由下面几个 ID 进行唯一识别

- **VID** # Vendor ID。成为 PCI-SIG 会员的公司会获得 Vendor ID，可以从[这里](https://pcisig.com/membership/member-companies)检索 Vendor ID 对应的公司名称（e.g. 0x8086 表示英特尔公司）
- **DID** # Device ID
- **SID** # Subsystem ID（有时候可以理解为 SDID）
- **SVID** # Subsystem-Vendor ID
- **RID** # Revision ID，也称 Rev ID，i.e. 版本号
- **CC** # Class-Code 类型代码。
- etc.

TODO:

- https://stackoverflow.com/questions/49050847/how-is-pci-segmentdomain-related-to-multiple-host-bridgesor-root-bridges
- https://pcisig.com/
- https://pcisig.com/specifications
- https://pcisig.com/specifications/conventional-pci/

# Linux 中的 PCI

> 参考:
>
> - https://github.com/torvalds/linux/blob/master/Documentation/PCI/sysfs-pci.rst
> - [How To Write Linux PCI Drivers(如何编写 PCI 驱动)](https://github.com/torvalds/linux/blob/master/Documentation/PCI/pci.rst)
> - [博客园，如何编写Linux PCI驱动程序](https://www.cnblogs.com/wanglouxiaozi/p/15525726.html)

https://github.com/torvalds/linux/blob/v6.9/include/linux/pci.h#L322 - `struct pci_dev`

https://github.com/torvalds/linux/blob/master/include/linux/pci_ids.h 定义了所有 PCI 类别、供应商和设备 ID。如果有 pci.ids 文件，为啥还要自己定义呢？

https://www.makelinux.net/ldd3/ - [12.1. The PCI Interface](https://www.makelinux.net/ldd3/chp-12-sect-1.shtml)

TODO: 下列说法需要确认来源

PCI 地址的命名规范 `domain:bus:slot.func` （e.g. 0000:00:15.1）是由 PCI 本地总线规范所定义的。这个规范由 PCI 特别兴趣小组 (PCI SIG) 制定和维护。

这个规范中,每个部分的含义如下:

1. **domain**: 表示 PCI 域编号,用于识别主机系统中的不同 PCI 主机桥。在较早期的系统中,只有一个域编号为 0。随着系统规模扩大,可能存在多个 PCI 域。
2. **bus**: 表示 PCI 总线编号。一个 PCI 域中可能包含多条 PCI 总线,每条总线都有一个唯一编号。
3. **slot**: 表示 PCI 插槽编号。每条 PCI 总线上可以连接多个 PCI 设备,每个设备对应一个插槽编号。
4. **func**: 表示 PCI 功能编号。一些 PCI 设备可能包含多个独立的功能,每个功能都有一个编号,用于区分和访问。对于单功能设备,该编号通常为 0。

该文件的开头注释中有这么一段

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

可以查查这三个贵方

- PCI BIOS Specification
- PCI Local Bus Specification
- PCI to PCI Bridge Specification
- PCI Express Specification
- PCI System Design Guide

## 关联文件

**pci.ids** # 所有已知 PCI ID（vendors(供应商)、devices(设备)、classes(类别)、subclasses(子类别)）列表。可以称为 **PCI ID 存储库**。

**.pciids-cache** # 通过 DNS 查询模式获取到的所有 PCI ID 缓存文件

### pci.ids

> 参考:
>
> - https://pci-ids.ucw.cz/
>   - [GitHub 项目，pciutils/pciids](https://github.com/pciutils/pciids)

pci.ids 是 PCI 设备中使用的所有已知 ID 的公共存储库：供应商、设备、子系统和设备类别的 ID。它在各种程序（例如 [PCI 实用程序](http://mj.ucw.cz/sw/pciutils/)）中用于显示完整的人类可读名称，而不是神秘的数字代码。

该文件由 [pciutils](https://github.com/pciutils) 维护

- 不同的 Linux 发行版该文件所在位置不同，通常都是 `/usr/share/XXX/pci.ids`
- TODO: 该文件的内容如何生成的？比如 /sys/devices/pci0000:00/0000:00:00.0/vendor 文件的值怎么来的？从 [include/linux/pci.h](https://github.com/torvalds/linux/blob/master/include/linux/pci.h) 文件的提示中 `For more information, please consult the following manuals (look at http://www.pcisig.com/ for how to get them)` 合理猜测这些数据都是从 PCI SIG 定义的规范中获取到的？

AI 回答，待验证

### 维护和更新

**维护机构**：

- `pci.ids` 文件由 `pciutils` 项目维护，该项目的开发者和社区成员定期更新这个文件。
- `pciutils` 项目由 Martin Mares 维护，托管在 `https://github.com/pciutils/pciids`。

**数据来源**：

- 新的供应商 ID 和设备 ID 通常由硬件制造商提供。
- 制造商向 [PCI-SIG](docs/Standard/IT/PCI-SIG.md) 注册他们的 ID。PCI-SIG 是管理和分配 PCI ID 的机构。
- 社区成员和开发者通过手动添加新设备的信息到 `pci.ids` 文件中。

**更新过程**：

- 新的设备信息通过提交补丁的方式添加到 `pciutils` 项目的仓库中。
- 开发者在 GitHub 上提交 Pull Request 或直接发送补丁给维护者。
- 维护者审核并合并这些更改。
- 合并后的更新会在 `pciutils` 项目的新版本中发布。



## 其他

