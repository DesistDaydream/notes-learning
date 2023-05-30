---
title: Device 命令
---

# 概述

> 参考：
> - [官方 Manual(手册)，DEVICE COMMANDS](https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#device-commands)

# attach-device # attach device from an XML file

# attach-disk # 将一个新的磁盘设备添加到 domian 中

https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#attach-disk

## Syntax(语法)

**virsh attach-disk DOMAIN SOURCE TARGET [OPTIONS]**

DOMAIN 可以使用 --domain 选项指定，SOURCE 可以使用 --source 选项指定，TARGET 可以使用 --target 选项指定

**OPTIONS**

- **--driver DRIVER** # 可用的值有：qemu
- **--subdriver STRING** # `默认值：raw`。可用的值有：对于 QEMU 来说可以是 raw、qcow2；对于 Xen 来说可以是 aio
- **--target STING** # 指定暴露给操作系统的的总线或设备。如果是硬盘的话，通常的值是 vdb、vdc、vdd 这种。
- **--targetbus STRING** # 指定要模拟的设备类型。`默认值：从设备名称的样式中推断出总线类型`
    - 可用的值：virtio
- **--cache STRING** # 可用的值：none

## EXAMPLE

先创建一个 qcow2 文件，然后使用这个文件进行测试

- qemu-img create -f qcow2 -o size=1G test-data.qcow2

TODO: virsh attach-disk tj-test-spst-common-ubuntu2204 /var/lib/libvirt/images/test-data.qcow2 vdb 这种命令添加的设备没有容量显示有问题

# attach-interface # 附加网络接口(i.e.给 VM 添加一个网卡)

# detach-device # detach device from an XML file

# detach-disk # detach disk device

# detach-interface # 分离网络接口(i.e.删除 VM 的一个网卡)

# update-device # update device from an XML file

# update-memory-device

# change-media # Change media of CD or floppy drive
