---
title: "Linux 硬件管理工具"
linkTitle: "Linux 硬件管理工具"
weight: 1
---

# 概述

> 参考：
>
> -

# lshw

> 参考：
>
> - [GitHub 项目，lyonel/lshw](https://github.com/lyonel/lshw)
> - [官网](http://lshw.ezix.org/)

TODO: 这命令哪来的？看的信息还挺全乎

# lspci

> 参考：
>
> - [Manual(手册)，lspci(8)](https://man7.org/linux/man-pages/man8/lspci.8.html)
> - [GitHub 项目，pciutils/pciutils](https://github.com/pciutils/pciutils)
> - [官网](https://mj.ucw.cz/sw/pciutils/)

列出所有 PCI 设备。

CentOS 包：pciutils

Ubuntu 包：pciutils

其中还有 setpci 工具用来配置 PCI 设备。

## Syntax(语法)

**lspci [OPTIONS]**

### OPTIONS

展示内容相关选项

- **-k** # 显示处理每个设备的内核驱动程序以及能够处理它的内核模块。在正常输出模式下给出 -v 时默认打开。 （目前仅适用于内核为 2.6 或更新版本的 Linux。）

选择指定设备选项

- **`-s [[[[<DOMAIN>]:]<BUS>]:][<DEVICE>][.[<FUNC>]]`** # 仅显示指定域中的设备（如果您的机器有多个主机桥，它们可以共享一个公共总线编号空间，或者它们中的每一个都可以寻址自己的 PCI 域；域编号从 0 到 ffff），bus （0 到 ff）、设备（0 到 1f）和功能（0 到 7）。设备地址的每个组成部分都可以省略或设置为“*”，均表示“任意值”。所有数字都是十六进制的。例如，“0：”表示总线 0 上的所有设备，“0”表示任何总线上设备 0 的所有功能，“0.3”选择所有总线上设备 0 的第三个功能，“.4”仅显示每个总线上的第四个功能设备。
  - 注意：-s 的值可以通过 uevent 文件中的 PCI_SLOT_NAME 字段的值获取

# smartctl

smartctl -a /dev/sda

# 网卡

## mii-tool



# USB 管理工具

## usbutils

> 参考：
> - [GitHub 项目，gregkh/usbutils](https://github.com/gregkh/usbutils)
> - [官网](http://www.linux-usb.org/)
> - [Manual(手册)，lsusb(8)](https://man7.org/linux/man-pages/man8/lsusb.8.html)

适用于 Linux 的 USB 实用程序，包括 lsusb。这是在 Linux 和 BSD 系统上使用的 USB 工具的集合，用于查询连接到系统的 USB 设备类型。这将在 USB 主机 (即您插入USB设备的机器) 上运行，而不是在 USB 设备 (即您插入USB主机的设备) 上运行。

包括如下几个工具

- lsusb
- usb-devices
- usbhid-dump
- usbreset

### lsusb Syntax(语法)

列出系统上的USB总线和USB设备的详细信息。在输出中，您将看到USB控制器的制造商、型号和当前的状态。

- `sudo lshw -class bus -class usb`

只查看有关USB设备的更详细信息

- `lshw -class usb`