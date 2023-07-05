---
title: Libvirt
weight: 1
---

# 概述

> 参考：
>
> - [GitHub 组织，libvirt](https://github.com/libvirt)
> - [官网](https://libvirt.org/)

Libvirt 项目是用于管理虚拟化平台的工具包，包括开源的 API，后台程序和管理工具。它可以用于管理 KVM、Xen、VMware ESX，QEMU 和其他虚拟化技术。Libvirt 将虚拟机统一称为 **Domain**。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gglb2f/1616123800173-58542239-2205-4586-bcc0-4edde6579a3f.png)

Libvirt 提供了管理虚拟机和其它虚拟化功能（如：存储和网络接口等）的便利途径。这些软件包括：一个长期稳定的 C 语言 API、一个守护进程（libvirtd）和一个命令行工具（virsh）。Libvirt 的主要目标是提供一个单一途径以管理不同类型的虚拟化环境(也称为 drivers 或者 hypervisors )，包括：KVM/QEMU，Xen，VMware， VirtualBox hypervisors，LXC，OpenVZ

Libvirt 包含 3 个东西：

- **libvirtd** # 是守护进程，服务程序，接收和处理 API 请求
- **API** # API 库使得其他人可以开发基于 Libvirt 的高级工具，比如 virt-manager、virt-install、virt-viewer 等。
- **virsh** # 是我们经常要用的命令行工具

Note：其实 libvirtd 在绝大部分情况下是与 qemu/kvm 相搭配来使用，都是开源的，并且 redhat 官方推荐的也是使用 libvirt 管理 kvm 虚拟机

## Libvirt 原理

libvirt 支持不同的虚拟化类型，所以需要一种方法来指定所要连接的虚拟化驱动。

libvirt 使用 URI 来与各种类型的虚拟化程序连接。[Libvirt 对接 Hypervisor](docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/Libvirt%20对接%20Hypervisor.md)

## Libvirt 支持的虚拟化平台版本

> 参考：
> 
> - [官方文档，支持的主机平台](https://libvirt.org/platforms.html)

这些平台用作决定 libvirt 所依赖的第 3 方软件所需的最低版本的基础。如果此处未列出平台，并不意味着 libvirt 无法工作。如果未列出的平台具有与已列出的平台相当的软件版本，则完全可以期望它能够正常工作。对于在未列出的平台上遇到的问题，欢迎报告错误，除非它们明显比此处描述的版本更旧。

请注意，当将发行版中提供的软件版本视为支持目标时，libvirt 仅考虑版本号，并假设该发行版中的功能与具有相同版本的上游版本相匹配。换句话说，如果发行版向后移植额外的功能到其发行版中的软件，libvirt 上游代码将不会添加对这些向后移植的显式支持，除非该功能可以以也适用于上游版本的方式自动检测。

Repology站点是一个有用的资源，可用于识别各种操作系统中当前发布的软件版本，但它并不涵盖下面列出的所有发行版[。](https://repology.org/)

- [libvirt 上的 Repology](https://repology.org/metapackage/libvirt/versions)
- [qemu 上 Repology](https://repology.org/metapackage/qemu/versions)

# Libvirt 关联文件与配置

**/etc/libvirt/** # 配置保存路径

- **./libvirt.conf** # libvirt 客户端的配置文件。用于配置连接虚拟化驱动的 URI 别名、默认 URI、等等。
- .**/libvirtd.conf** # libvirtd 守护进程的配置文件。
- .**/qemu/** # xml 格式的配置文件存放路径，配置文件包括该 VM 的元数据(名字，uuid，内存，cpu 等)，设备配置(包括使用的硬盘文件的路径，网络类型等)，配置文件为 xml 格式。创建完一台 VM 后，会在该目录下生成对应 VM 名字的 xml 文件
- **./network/** #

**/var/lib/libvirt/** # 数据保存路径

- **./images/** # 所有通过 libvirt 创建的虚拟机所生成的 images 都保存在该目录下
- **./qemu/snapshot/** # 创建快照 xml 文件都保存在该目录下

**/run/libvirt/** # 

# XML 文件

Libvirt 管理的虚拟机都可以通过 XML 文件来描述其所应该模拟的硬件设备、状态等等。详见 [XML 文件详解](/docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/XML%20文件详解/XML%20文件详解.md)

我们甚至可以通过直接编写 XML 文件，以便 Libvirt 相关工具直接读取 XML 并启动 VM。

# nvram 文件

TODO: 这是什么文件？？？？
