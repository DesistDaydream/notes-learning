---
title: lspci
linkTitle: lspci
date: 2024-06-15T11:36
weight: 20
---

# 概述

> 参考：
>
> - [Manual(手册)，lspci(8)](https://man7.org/linux/man-pages/man8/lspci.8.html)
> - [GitHub 项目，pciutils/pciutils](https://github.com/pciutils/pciutils)
> - [官网](https://mj.ucw.cz/sw/pciutils/)

列出所有 [PCI](docs/0.计算机/Motherboard/PCI.md) 设备。在列出的设备信息中，包含了一些供应商的名称、分类名称、etc. 信息。这些信息统一从 [pci.dis](https://pci-ids.ucw.cz/) 文件中获取。

CentOS 包：pciutils

Ubuntu 包：pciutils

其中还有 setpci 工具用来配置 PCI 设备。

TODO: lspci 是如何从 Linux 中拿到 PCI 设备列表的？

# Syntax(语法)

**lspci [OPTIONS]**

## OPTIONS

**基本显示模式**

- **-mm** # 以机器可读的形式转储 PCI 设备数据，以便于脚本解析。  显示内容详见 [详细格式](#详细格式)，通常与 -v 一起使用
- **-m** # 与 -mm 一样，但是向后兼容。不要在新代码中使用

**展示内容相关选项**

- **-v** # 显示 PCI 设备的详细信息。多次使用 -v 可以显示更多信息，最多支持 3 个 v。
- **-k** # 显示用于处理 PCI 设备的内核驱动以及能够处理ta的内核模块。在正常输出模式下若使用 -v 选项时默认打开 -k。
- **-D** # 始终显示 PCI 的 domain 部分。

**选择指定设备选项**

- **`-s [[[[<DOMAIN>]:]<BUS>]:][<DEVICE>][.[<FUNC>]]`** # 仅显示指定域中的设备（如果您的机器有多个主机桥，它们可以共享一个 [Bus](docs/0.计算机/Motherboard/Bus.md) 编号空间，或者它们中的每一个都可以寻址自己的 PCI 域；域编号从 0 到 ffff），bus （0 到 ff）、设备（0 到 1f）和功能（0 到 7）。设备地址的每个组成部分都可以省略或设置为 `*`，均表示“任意值”。所有数字都是十六进制的。例如，“0：”表示总线 0 上的所有设备，“0”表示任何总线上设备 0 的所有功能，“0.3”选择所有总线上设备 0 的第三个功能，“.4”仅显示每个总线上的第四个功能设备。
  - 注意：-s 的值可以通过 uevent 文件中的 PCI_SLOT_NAME 字段的值获取
- **`-d [<VENDOR>]:[<DEVICE>][:<CLASS>[:<PROG-IF>]]`**

ID **与名称的解析行为选项**

- **-n** 显示数字化ID,而不是名称。
- **-nn** 显示数字化供应商和设备ID。


TODO

# EXAMPLE

`lspci -Dvmmnnk` 显示效果如下

```bash
~]# lspci -Dvmmnnk | more
Slot:   0000:00:00.0
Class:  Host bridge [0600]
Vendor: Intel Corporation [8086]
Device: Xeon E7 v3/Xeon E5 v3/Core i7 DMI2 [2f00]
SVendor:        Intel Corporation [8086]
SDevice:        Device [0000]
Rev:    02
ProgIf: 00
NUMANode:       0

Slot:   0000:00:02.0
Class:  PCI bridge [0604]
Vendor: Intel Corporation [8086]
Device: Xeon E7 v3/Xeon E5 v3/Core i7 PCI Express Root Port 2 [2f04]
SVendor:        Intel Corporation [8086]
SDevice:        Device [0000]
Rev:    02
ProgIf: 00
Driver: pcieport
NUMANode:       0

......略
```

## 详细格式

使用 -vmm 选项输出下面格式的内容

```bash
~]# lspci -vmm  | more
Slot:   00:00.0
Class:  Host bridge
Vendor: Intel Corporation
Device: Xeon E7 v3/Xeon E5 v3/Core i7 DMI2
SVendor:        Intel Corporation
SDevice:        Device 0000
Rev:    02
ProgIf: 00
NUMANode:       0

Slot:   00:02.0
Class:  PCI bridge
Vendor: Intel Corporation
Device: Xeon E7 v3/Xeon E5 v3/Core i7 PCI Express Root Port 2
SVendor:        Intel Corporation
SDevice:        Device 0000
Rev:    02
ProgIf: 00
NUMANode:       0

......略
```

-vmm 输出的格式是由**空行**分隔的 PCI 设备信息。每个 PCI 设备信息称为 record(记录)，每条有多行，每行是一个 `TAG: VALUE` 对。标记和值由单个制表符分隔。  

有如下 TAG 可用

- **Slot** # 设备所在插槽的名称（`[domain:]bus:device.function`）。  该 TAG 始终是记录中的第一个。
- **Class** # 类别
- **Vendor** # 供应商
- **Device** # 设备名称或编号