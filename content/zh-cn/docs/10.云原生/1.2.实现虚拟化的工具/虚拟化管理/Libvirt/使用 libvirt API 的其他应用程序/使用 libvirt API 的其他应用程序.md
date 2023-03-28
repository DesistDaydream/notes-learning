---
title: 使用 libvirt API 的其他应用程序
weight: 1
---

# 概述

> 参考：
>
> - [官网，使用 libvirt 的应用程序](https://libvirt.org/apps.html)

# virt-host-validate # 虚拟化宿主机环境检验

# virt-manager # 图形化模式的虚拟机管理工具

virt-manager 是一个图形化的应用程序，通过 libvirt 管理虚拟机。

virt-manager 提供了多个配套的工具

- virt-install # 是一个命令行工具，它提供了一种将操作系统配置到虚拟机中的简单方法。
- virt-viewer # 是一个轻量级的 UI 界面，用于与虚拟客户操作系统的图形显示进行交互。它可以显示 VNC 或 SPICE，并使用 libvirt 查找图形连接详细信息。
- virt-clone # 是一个用于克隆现有非活动客户的命令行工具。它复制磁盘映像，并使用指向复制磁盘的新名称、UUID 和 MAC 地址定义配置。
- virt-xml # 是一个命令行工具，用于使用 virt-install 的命令行选项轻松编辑 libvirt 域 XML。
- virt-bootstrap # 是一个命令行工具，提供了一种简单的方法来为基于 libvirt 的容器设置根文件系统。

详见 [《virt-manager》](/docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/使用%20libvirt%20API%20 的其他应用程序/virt-manager.md libvirt API 的其他应用程序/virt-manager.md) 章节

# libguestfs：是一组用于访问和修改虚拟机（VM）磁盘映像的工具

> 参考：

可以使用它来查看和编辑 guest 虚拟机内的文件，更改 VM 的脚本， 监视磁盘使用/免费统计信息， 创建 guest 虚拟机，P2V， V2V，执行备份，克隆 VM，构建 VM，格式化磁盘，调整磁盘大小等等。

使用 yum 安装即可直接使用：yum install libguestfs-tools
