---
title: Libvirt API
weight: 1
---

# 概述

> 参考：
> 
> - [官网，API 参考](https://libvirt.org/html/index.html)
> - [官网，使用 libvirt 的应用程序](https://libvirt.org/apps.html)

在官方的 API 参考中，包含所有 C 语言表示的 API 信息，这些 API 分为如下几类：

- [common](https://libvirt.org/html/libvirt-libvirt-common.html) # libvirt 和 libvirt-admin 库的常用宏和枚举
- [domain](https://libvirt.org/html/libvirt-libvirt-domain.html) # 用于管理 Domain 的 API
- [domain checkpoint](https://libvirt.org/html/libvirt-libvirt-domain-checkpoint.html) # 用于管理 Domain 检查点的 API
- [domain snapshot](https://libvirt.org/html/libvirt-libvirt-domain-snapshot.html) # 用于管理 Domain 快照的 API
- [error](https://libvirt.org/html/libvirt-virterror.html)
- [event](https://libvirt.org/html/libvirt-libvirt-event.html)
- [host](https://libvirt.org/html/libvirt-libvirt-host.html) # 用于管理主机的 API
- [interface](https://libvirt.org/html/libvirt-libvirt-interface.html)
- [network](https://libvirt.org/html/libvirt-libvirt-network.html)
- [node device](https://libvirt.org/html/libvirt-libvirt-nodedev.html)
- [network filter](https://libvirt.org/html/libvirt-libvirt-nwfilter.html)
- [secret](https://libvirt.org/html/libvirt-libvirt-secret.html)
- [storage](https://libvirt.org/html/libvirt-libvirt-storage.html)
- [stream](https://libvirt.org/html/libvirt-libvirt-stream.html) and [admin](https://libvirt.org/html/libvirt-libvirt-admin.html)
- [QEMU](https://libvirt.org/html/libvirt-libvirt-qemu.html)
- [LXC](https://libvirt.org/html/libvirt-libvirt-lxc.html) libs

在我们使用非 C 代码，比如 Python、Go 的时候，如果想要找到 API 的说明，可以参考官方的 C 语言的 API 文档，函数名基本差不多，用一个最常见的“列出所有 Domain”举例：

- C: [virConnectListAllDomains](https://libvirt.org/html/libvirt-libvirt-domain.html#virConnectListAllDomains)
- Go: [ListAllDomains](https://pkg.go.dev/libvirt.org/go/libvirt#Connect.ListAllDomains)
- Python: listAllDomains
  - Python 好像没有 Go Package 网站那种东西，没法找到在线的 API 文档
  - 有个类似这样的代码 `raise libvirtError("virConnectListAllDomains() failed")` 可以知道该方法对应的是哪个 C API 方法。

总的来说，还是 Go 好用，追踪代码后注释明确，主时钟也包含对应的 C API 的链接

# 使用 Libvirt API 的应用程序

除了 virsh 以外，还有很多使用 Libvirt API 的应用

## virt-host-validate # 虚拟化宿主机环境检验

## virt-manager # 图形化模式的虚拟机管理工具

virt-manager 是一个图形化的应用程序，通过 libvirt 管理虚拟机。

virt-manager 提供了多个配套的工具

- virt-install # 是一个命令行工具，它提供了一种将操作系统配置到虚拟机中的简单方法。
- virt-viewer # 是一个轻量级的 UI 界面，用于与虚拟客户操作系统的图形显示进行交互。它可以显示 VNC 或 SPICE，并使用 libvirt 查找图形连接详细信息。
- virt-clone # 是一个用于克隆现有非活动客户的命令行工具。它复制磁盘映像，并使用指向复制磁盘的新名称、UUID 和 MAC 地址定义配置。
- virt-xml # 是一个命令行工具，用于使用 virt-install 的命令行选项轻松编辑 libvirt 域 XML。
- virt-bootstrap # 是一个命令行工具，提供了一种简单的方法来为基于 libvirt 的容器设置根文件系统。

详见 [virt-manager](/docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/Libvirt%20API/virt-manager.md) 章节

## libguestfs：是一组用于访问和修改虚拟机（VM）磁盘映像的工具

> 参考：

可以使用它来查看和编辑 guest 虚拟机内的文件，更改 VM 的脚本， 监视磁盘使用/免费统计信息， 创建 guest 虚拟机，P2V， V2V，执行备份，克隆 VM，构建 VM，格式化磁盘，调整磁盘大小等等。

使用 yum 安装即可直接使用：yum install libguestfs-tools
