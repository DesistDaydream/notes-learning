---
title: Device 命令
---

# 概述

> 参考：
>
> - [官方 Manual(手册)，DEVICE COMMANDS](https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#device-commands)

注意：由于 Libvirt 适用于多种虚拟化后端，所以每个选项可用的值会分别适用于多种虚拟化后端，但是个人一般只用 QEMU/KVM，所以笔记中的选项一般也只记录 QEMU/KVM 相关的值。

这部分命令主要是用来为虚拟机添加/移除各种设备，比如 网卡、硬盘 等等。这些添加/移除的命令有一些通用的选项可用，为了记笔记方便，在开头这统一记录一下，添加设备的命令通常以 **attach** 开头，移除设备的命令通常以 **detach** 开头。

生效策略配置选项，用于配置添加/移除的行为在什么时候生效：

  - **--config** # 影响已关机的 Domain，将会在下一次启动时添加/移除设备
  - **--live** # 影响运行中的 Domain，立刻为虚拟机添加/移除设备
  - **--current** # 等效于 --live 或 --config，具体取决于虚拟机当前的状态
  - **--persistent** # 处于兼容的目的，该命令对关机或者开机状态的虚拟机都有效，相当于为当前运行中的虚拟机以及以后启动后的虚拟机都添加/移除设备。

# attach-device - attach device from an XML file

# 添加与移除磁盘设备

在 Libvirt 的[最佳实践](/docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/最佳实践.md)中有命令使用示例

## attach-disk

将一个新的磁盘设备添加到 domian 中

https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#attach-disk

### Syntax(语法)

**virsh attach-disk DOMAIN SOURCE TARGET [OPTIONS]**

将 SOURCE 添加到 DOMAIN 中，作为 TARGET 磁盘设备。

- SOURCE 是本地的 qcow2、raw 这种格式的文件。如果指定了 --source-protocol 选项，则 SOURCE 可以是网络磁盘。
- TARGET 是虚拟机中的设备，比如 vdb、vdc 这种。
  - 可以使用 --target 选项指定 TARGET

**OPTIONS**

- **--driver DRIVER** # 指定要使用的磁盘驱动程序。
  - 可用的值有：对于 QEMU 来说可以是 qemu；对于 Xen 来说可以是 file、tap、phy
- **--subdriver STRING** # 为 --driver 选项提供更多的详细信息。
  - 可用的值有：对于 QEMU 来说可以是 raw、qcow2；对于 Xen 来说可以是 aio
- **--target STING** # 指定暴露给操作系统的的总线或设备。如果是硬盘的话，通常的值是 vdb、vdc、vdd 这种。
- **--targetbus STRING** # 指定要模拟的设备类型。`默认值：从设备名称的样式中推断出总线类型`
  - 可用的值：virtio、sata、scsi、usb
- **--cache STRING** # 可用的值：none

## detach-disk

将虚拟机的磁盘设备与虚拟机分离

### Syntax(语法)

**detach-disk DOMAIN TARGET [OPTIONS]**

# attach-interface - 附加网络接口(i.e.给 VM 添加一个网卡)

# detach-device - detach device from an XML file

# detach-interface - 分离网络接口(i.e.删除 VM 的一个网卡)

# update-device - update device from an XML file

# update-memory-device

# change-media - Change media of CD or floppy drive
