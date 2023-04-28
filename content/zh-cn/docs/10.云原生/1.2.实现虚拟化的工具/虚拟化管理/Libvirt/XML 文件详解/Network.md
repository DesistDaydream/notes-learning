---
title: "Network"
linkTitle: "Network"
weight: 3
---

# 概述

> 参考：
> 
> - [官方文档，Network XML 格式](https://libvirt.org/formatnetwork.html)

使用 `virsh net-list` 命令列出所有网络

# bandwidth

`<bandwidth>`元素允许为特定网络设置服务质量（自0.9.4版起）。只支持为具有 `<forward>` 模式为route、nat、bridge或没有模式的网络（即“隔离”网络）设置带宽。不支持为forward模式为passthrough、private或hostdev的网络设置带宽。尝试这样做将导致无法定义网络或创建临时网络。

average

peak

burst

# 最佳实践

限制虚拟机网卡的网宿

```xml
    <interface type='bridge'>
      <source bridge='br0'/>
      <model type='virtio'/>
      <driver name='vhost' queues='8'/>
      <bandwidth>
        <inbound average='1000' peak='5000' burst='1024'/>
        <outbound average='128' peak='256' burst='256'/>
      </bandwidth>
      <address type='pci' domain='0x0000' bus='0x01' slot='0x00' function='0x0'/>
    </interface>

```