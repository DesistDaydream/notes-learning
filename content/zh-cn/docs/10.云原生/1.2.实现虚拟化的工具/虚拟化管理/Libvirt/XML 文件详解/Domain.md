---
title: "Domain"
linkTitle: "Domain"
weight: 2
---

# 概述

> 参考：
> 
> - [官方文档，Domain XML 格式](https://libvirt.org/formatdomain.html)

Domain 对象对应 `<domain>` 根元素，该元素中有如下属性：

- **type** # 指定用于运行域的管理程序。允许的值是特定于驱动程序的，但包括“xen”、“kvm”、“hvf”（自 8.1.0 和 QEMU 2.12 起）、“qemu”和“lxc”。
- **id** # 它是正在运行的客户机的唯一整数标识符。非活动机器没有 id 值。

下列元素都属于 `<domain>` 这个根元素的子元素

- TODO:
- ......
- name
- uuid
- metadata
- memory
- [os](#os)
- [devices](#devices)
- ......

上面这些元素可以控制整个 Domain，通常可以分为几大类，下面的笔记将以这些类别进行记录

- 通用元数据
- [系统引导](#系统引导)
- SMBIOS 系统信息
- CPU 分配
- ....... 等等
- [设备](#设备)

# 系统引导

有很多种方式可以引导 Domain，不过个人使用中最常见的就是使用 BIOS 引导。官方文档中还有 Host 引导、直接内核引导、容器引导。等有需要了再记录。

支持完全虚拟化的虚拟机管理程序可以通过 BIOS 启动。在这种情况下，BIOS 具有引导顺序优先级 (软盘，硬盘，光盘，网络)，以确定在何处获取/查找所需的引导镜像。

**os** # 配置操作系统相关信息

- **type** # 操作系统的类型。可用的值有 hvm、linux。
  - 属性：
    - arch # CPU 的架构。比如 x86_64
    - machine # 机器类型
- **boot** # 指定 Domain 下次如何引导启动。
  - 属性：
    - dev # 指定 Domain 下次启动时的引导设备，可指定多次设置多个引导设备。可用的值有: fd、hd、cdrom、network。（fd 指软盘，hd 指硬盘）

## 配置示例

```xml
  <os>
    <type arch='x86_64' machine='pc-q35-6.2'>hvm</type>
    <boot dev='cdrom'/>
  </os>
```

# 设备

https://libvirt.org/formatdomain.html#devices

## Network interfaces

https://libvirt.org/formatdomain.html#network-interfaces

### Quality of service(QOS)

**bandwidth**

- bandwidth 元素可以为虚拟机中的每张网卡设置 QOS(服务质量)，可以为流量的出与入分别设置。QOS 下可用的元素详见 [Network XML](/docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/XML%20文件详解/Network.md) 的 [bandwidth 元素](/docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/XML%20文件详解/Network.md#bandwidth)。

# Graphical framebuffers

https://libvirt.org/formatdomain.html#graphical-framebuffers

Graphical 设备允许与 Domain 进行图形交互。我们通常会配置一个帧缓冲区或一个文本控制台来允许与管理员交互。

下面是几种配置示例

```xml
...
<devices>
  <!-- 这是啥 -->
  <graphics type='sdl' display=':0.0'/>
  <!-- 使用 VMC 作为图形交互。配置 VNC 的监听端口等 -->
  <graphics type='vnc' port='5904' sharePolicy='allow-exclusive'>
    <listen type='address' address='1.2.3.4'/>
  </graphics>
  <!-- 使用 RDP  -->
  <graphics type='rdp' autoport='yes' multiUser='yes' />
  <!-- 这是啥 -->
  <graphics type='desktop' fullscreen='yes'/>
  <!-- 这是啥 -->
  <graphics type='spice'>
    <listen type='network' network='rednet'/>
  </graphics>
</devices>
...
```