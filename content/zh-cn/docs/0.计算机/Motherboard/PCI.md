---
title: PCI
linkTitle: PCI
date: 2024-06-15T09:24
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, Peripheral_Component_Interconnect](https://en.wikipedia.org/wiki/Peripheral_Component_Interconnect)

**Peripheral Component Interconnect(外围组件互连，简称 PCI)**

## PCI-E

> 参考:
>
> - [Wiki, PCI_Express](https://en.wikipedia.org/wiki/PCI_Express)

**Peripheral Component Interconnect Express(简称 PCI-E)**


# PCI 规范

> 参考:
>
> - [pci设备身份识别码介绍说明](https://www.twblogs.net/a/5eee10c7264079afec950f51)

由 [PCI-SIG](/docs/Standard/IT/PCI-SIG.md) 制定规范。一个 PCI 设备的通常由下面几个 ID 进行唯一识别

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

详见 [PCI](/docs/1.操作系统/Kernel/Hardware/PCI.md)


