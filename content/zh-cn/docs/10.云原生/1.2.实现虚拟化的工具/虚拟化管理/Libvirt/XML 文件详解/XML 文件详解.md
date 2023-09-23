---
title: "XML 文件详解"
linkTitle: "XML 文件详解"
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，XML 格式](https://libvirt.org/format.html)

Libvirt API 中的**对象**使用 [**XML 格式**](/docs/2.编程/无法分类的语言/XML.md) 的文档进行配置，以便在未来的版本中轻松扩展。每个 XML 文档都有一个关联的 Relax-NG 模式，可用于在使用前验证文档。

这里面的 Libvirt API 对象指的就是 Domain(虚拟机)、存储、快照、网络 等等。对于 Libvirt，所有 VM 相关的资源都会抽象为对象，这样也利于代码编写。

> Kubernetes 的 API 对象跟这个有点像，只不过 Kubernetes 中，使用 YAML 格式来声明对象，而不是 XML 格式来配置对象。

下面是所有可用的 Libvirt API 对象，每个 **Libvirt 对象** 都对应一个 **根元素**。

- [Domain](/docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/XML%20文件详解/Domain.md) # 虚拟机相关的 XML 配置，可以直接使用 Domain XML 文件创建、启动、管理虚拟机。
  - 根元素名称: `<domain>`
- [Network](/docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/XML%20文件详解/Network.md) # 虚拟网络相关的 XML 配置。
  - 根元素名称: `<network>`
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
