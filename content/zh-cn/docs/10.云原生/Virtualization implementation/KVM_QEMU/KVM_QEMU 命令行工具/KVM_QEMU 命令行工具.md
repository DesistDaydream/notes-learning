---
title: KVM/QEMU 命令行工具
linkTitle: KVM/QEMU 命令行工具
date: 2019-10-17T11:02:00
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，QEMU 用户文档](https://www.qemu.org/docs/master/system/qemu-manpage.html)
> - [官方文档，系统模拟-Invocation](https://www.qemu.org/docs/master/system/invocation.html)
> - [官方文档，工具](https://www.qemu.org/docs/master/tools/index.html)

KVM/QEMU 的虚拟机的生命周期是通过一系列 KVM/QEMU 工具集管理的，包括

- **qemu-img** # 虚拟机镜像管理工具
- **qemu-system-$ARCH** # 虚拟机运行时管理工具。
  - 注意：在 CentOS 系统中，该二进制文件的名字是 qemu-kvm，并且是一个在 /usr/local/bin/qemu-kvm 指向 /usr/libexec/qemu-kvm 的软链接
- 等等

通常情况下，我们不管是通过 virt-manager 程序创建的虚拟机、还是使用 Libvirt 工具包创建的虚拟机，本质上，都是调用的 **qemu-img、qemu-system-x86_64** 等工具。

如果用容器比较的话

- qemu-img 像各种容器镜像管理工具
- qemu-system-x86_64 像 runc

# qemu-img

[qemu-img](/docs/10.云原生/Virtualization%20implementation/KVM_QEMU/KVM_QEMU%20命令行工具/qemu-img.md)

# qemu-system

[qemu-system](/docs/10.云原生/Virtualization%20implementation/KVM_QEMU/KVM_QEMU%20命令行工具/qemu-system.md)
