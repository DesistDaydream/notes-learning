---
title: "Domains"
linkTitle: "Domains"
weight: 1
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


# Graphical framebuffers

https://libvirt.org/formatdomain.html#graphical-framebuffers

# Devices

https://libvirt.org/formatdomain.html#devices

## Network interfaces

https://libvirt.org/formatdomain.html#network-interfaces

