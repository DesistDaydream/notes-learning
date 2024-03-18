---
title: "Libvirt 守护进程"
linkTitle: "Libvirt 守护进程"
date: "2023-07-06T08:50"
weight: 2
---

# 概述

> 参考：
>
> - [官方文档，Libvirt 守护进程](https://libvirt.org/daemons.html)
> - [官方文档，手册-libvirtd](https://libvirt.org/manpages/libvirtd.html)

传统上，Libvirt 项目提供了一个名为 **libvirtd** 的单一守护进程，它暴露了对所有有状态 Driver 的支持，包括主要虚拟机管理程序驱动程序和辅助支持驱动程序。它还支持监听在 TCP/IP 上，以便主机外运行的客户端进行安全远程访问。

**未来，想要将一个整体的 libvirtd 守护进程替换为一组模块化的守护进程，以 `virt${DRIVER}d` 命令，就是每个驱动都有自己独立的守护进程。还有一个 virtproxyd 守护进程可以提供安全的远程访问。** —— 2023.7.5

# libvirtd

libvirtd 程序是 libvirt 虚拟化管理系统的服务器端守护进程组件。包含了部分 Hypervisor 驱动，并暴露了 [Libvirt API](/docs/10.云原生/Virtualization%20implementation/虚拟化管理/Libvirt/Libvirt%20API/Libvirt%20API.md)

该守护进程在主机服务器上运行，并为虚拟来宾执行所需的管理任务。这包括启动、停止和在主机服务器之间迁移来宾、配置和操作网络以及管理供来宾使用的存储等活动。

libvirt 客户端库和实用程序连接到此守护进程以发出任务并收集有关主机系统和来宾的配置和资源的信息。

默认情况下，libvirtd 守护进程侦听本地 Unix 域套接字上的请求。使用 -l | --listen 命令行选项，可以指示 libvirtd 守护进程另外侦听 TCP/IP 套接字。要使用的 TCP/IP 套接字在 libvirtd 配置文件中定义。

> 这里官方说的其实有一些问题，参考 https://stackoverflow.com/questions/65663825/could-not-add-the-parameter-listen-to-open-tcp-socket ，并且在下面关于启动模式中也有说明。

重新启动 libvirtd 不会影响正在运行的 guest 虚拟机。如果定义了 XML 配置，来宾将继续操作并将被自动接听。任何尚未定义 XML 配置的来宾都将从配置中丢失。

libvirtd 守护进程能够以两种模式启动

- **传统模式** # 它将自行创建并侦听 UNIX 套接字。
  - 如果给出了 --listen 参数，它还将根据 /etc/libvirt/libvirtd.conf 中的 listen_tcp 和 listen_tls 选项监听 TCP/IP 套接字
- **套接字激活模式** # 它将依靠 systemd 在 UNIX 和可选的 TCP/IP 套接字上创建和侦听，并将它们作为预打开的文件描述符传递。
  - <font color="#ff0000">注意</font>：在这种模式下，不允许传递 --listen 参数，并且 /etc/libvirt/libvirtd.conf 中大多数与套接字相关的配置选项将不再起作用。
  - 如果想要启用 TCP 或 TLS 套接字，可以开启 libvirtd-tcp.socket 或 libvirtd-tls.socket 这两个 Unit。

在使用 systemd 的主机操作系统上运行时，套接字激活模式通常是默认模式。要恢复到传统模式，必须使用如下命令屏蔽所有套接字单元文件：

```bash
$ systemctl mask libvirtd.socket libvirtd-ro.socket \
   libvirtd-admin.socket libvirtd-tls.socket libvirtd-tcp.socket
```

最重要的是，请确保 --timeout 参数不用于守护进程，因为它不会在以后的任何连接中启动。

## 远程访问

https://libvirt.org/remote.html
