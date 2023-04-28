---
title: "Domain"
linkTitle: "Domain"
weight: 2
---

# 概述

> 参考：
> 
> - [官方文档，Domain XML 格式](https://libvirt.org/formatdomain.html)

Domain 元素中有如下属性：

- **type** # 指定用于运行域的管理程序。允许的值是特定于驱动程序的，但包括“xen”、“kvm”、“hvf”（自 8.1.0 和 QEMU 2.12 起）、“qemu”和“lxc”。
- **id** # 它是正在运行的客户机的唯一整数标识符。非活动机器没有 id 值。

下列元素都属于 `<domain>` 这个根元素的子元素

- TODO:
- ......
- [devices](#devices)
- ......

# Devices

https://libvirt.org/formatdomain.html#devices

## Network interfaces

https://libvirt.org/formatdomain.html#network-interfaces

### Quality of service(QOS)

**bandwidth**

- bandwidth 元素可以为虚拟机中的每张网卡设置 QOS(服务质量)，可以为流量的出与入分别设置。QOS 下可用的元素详见 [Network XML](/docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/XML%20文件详解/Network.md) 的 [bandwidth 元素](/docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/XML%20文件详解/Network.md#bandwidth)。

# Graphical framebuffers

https://libvirt.org/formatdomain.html#graphical-framebuffers

