---
title: virt-manager
---

# 概述

> 参考
>
> - [GitHub 项目，virt-manager/virt-manager](https://github.com/virt-manager/virt-manager)
> - [官网](https://virt-manager.org/)

virt-manager 是一个图形化的应用程序，通过 libvirt 管理虚拟机。

virt-manager 提供了多个配套的工具

- virt-manager # GUI 模式的 VM 管理程序
- virt-viewer # 是一个轻量级的 UI 界面，用于与虚拟客户操作系统的图形显示进行交互。它可以显示 VNC 或 SPICE，并使用 libvirt 查找图形连接详细信息。
- virt-install # 是一个命令行工具，它提供了一种将操作系统配置到虚拟机中的简单方法。
- virt-clone # 是一个用于克隆现有非活动客户的命令行工具。它复制磁盘映像，并使用指向复制磁盘的新名称、UUID 和 MAC 地址定义配置。
- virt-xml # 是一个命令行工具，用于使用 virt-install 的命令行选项轻松编辑 libvirt 域 XML。
- virt-bootstrap # 是一个命令行工具，提供了一种简单的方法来为基于 libvirt 的容器设置根文件系统。

virt-clone、virt-xml、virt-install 属于安装虚拟机的工具，通常都在 virtinst 包中。

virt-manager、virt-viewer 属于图形化管理虚拟机的工具，通常都在 virt-manager 包中。

# virt-manager

## 使用 virt-manager 管理多台虚拟机

在一台机器上的 virt-manager 可以通过 add connection 管理其它宿主机上的虚拟机，但是前提是建立 ssh 的密钥认证，因为在 virt-manager 在通过 ssh 连接的时候，需要使用窗口模式输入密码，而一般情况下 ssh 是默认不装该组件的。如果不想添加密钥认证，那么安装 `ssh-askpass-gnome` 组件即可。

输入 virt-manager 打开管理界面。选择 File—Add Connecttion.. 勾选 Connect to remote host

依次填入文本框中内容如下：

- Hypervisor: QEMU/KVM
- Method:SSH
- Username:root
- Hostname:192.168.0.123（需要被操作的服务器地址）

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/abyoqm/1616123543706-5c47d241-1780-40d5-b04e-1cfd4b802d6c.png)
然后点击 connect 连接即可，此时就会出现另一台服务器上的虚拟机供操作。

# virt-viewer

可以显示虚拟机的图形控制台

# virt-install

可以将 virt-install 理解为 CLI 的 virt-manager，可以放在脚本中创建虚拟机，内容非常多，详见 [《virt-install》](/docs/10.云原生/Virtualization%20implementation/虚拟化管理/Libvirt/Libvirt%20API/virt-install.md) 章节。

# virt-clone

> 参考：
>
> - [Manual(手册),virt-clone(1)](https://man.cx/virt-clone)

## Syntax(语法)

**virt-clone \[OPTION]...**
