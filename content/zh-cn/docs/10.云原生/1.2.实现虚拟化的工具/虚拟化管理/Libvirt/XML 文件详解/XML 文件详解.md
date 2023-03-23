---
title: "XML 文件详解"
linkTitle: "XML 文件详解"
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，XML 格式](https://libvirt.org/format.html)

Libvirt API 中的**对象**使用 **XML 格式**的文档进行配置，以便在未来的版本中轻松扩展。每个 XML 文档都有一个关联的 Relax-NG 模式，可用于在使用前验证文档。

这里面的 Libvirt API 对象指的就是 Domain(虚拟机)、存储、快照、网络 等等。对于 Libvirt，所有 VM 相关的资源都会抽象为对象，这样也利于代码编写。

下面是**根元素**的名称，每个根元素都对应一个 **Libvirt 对象**。

- [Domain](/docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/XML%20文件详解/Domain%20XML.md) # 虚拟机相关的 XML 配置，可以直接使用 Domains XML 创建、启动、管理虚拟机。
- [Network](https://libvirt.org/formatnetwork.html)
- [Network filtering](https://libvirt.org/formatnwfilter.html)
- [Network ports](https://libvirt.org/formatnetworkport.html)
- [Storage](https://libvirt.org/formatstorage.html)
- [Storage encryption](https://libvirt.org/formatstorageencryption.html)
- [Capabilities](https://libvirt.org/formatcaps.html)
- [Domain capabilities](https://libvirt.org/formatdomaincaps.html)
- [Storage Pool capabilities](https://libvirt.org/formatstoragecaps.html)
- [Node devices](https://libvirt.org/formatnode.html)
- [Secrets](https://libvirt.org/formatsecret.html)
- [Snapshots](https://libvirt.org/formatsnapshot.html)
- [Checkpoints](https://libvirt.org/formatcheckpoint.html)
- [Backup jobs](https://libvirt.org/formatbackup.html)

# virt-xml-validate

virt-xml-validate 工具是一个简单的检验 XML 文档的工具，直接在命令后面加上 XML 文档的 /PATH/FILE 即可对该文件进行检验，确保其在传递给 libvirt 时是正确的。
