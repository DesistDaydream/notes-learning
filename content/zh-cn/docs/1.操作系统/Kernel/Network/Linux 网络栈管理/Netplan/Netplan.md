---
title: Netplan
weight: 1
---

# 概述

> 参考：
>
> - [GitHub 项目，canonical/netplan](https://github.com/canonical/netplan)
> - [官网](https://netplan.io/)

**Netplan** 是一个网络配置抽象渲染器。属于 netplan.io 包，通过 yaml 文件来管理 Linux 的网络配置。

**Netplan** 是用于在 Linux 系统上轻松配置网络的实用程序。只需为每个网络设备应该具有的配置，创建一个 YAML 格式的描述文件。 Netplan 将根据此描述为指定的 **Renderer(渲染器)** 生成所有必要的配置。剩下的工作，就是由这些 Renderer 来处理配置，并配置网络了。

# Netplan 的工作方式

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/vv96im/1616165548496-6a738611-4db5-4f06-9cfe-ce0c82d9cf01.png)

Netplan 从 /etc/netplan/\*.yaml 文件中读取配置信息。Netplan 启动初期，在 /run 目录中生成特定于后端的配置文件，以便让这些后端的网络守护程序根据这些配置文件管理网络设备。在 Netplan 中，这些特定的 **后端**被称为 **Renderers(渲染器)**。

Netplan 当前支持如下 Renderers(渲染器)：

- **networkd** # 默认 Renderer。该 Renderers 是 systemd 管理的网络管理程序 systemd-networkd，它属于 systemd 包
- **Network Manager** # 详见：[NetworkManager](/docs/1.操作系统/Kernel/Network/Linux%20网络栈管理/NetworkManager/NetworkManager.md)

注意：殊途同归，就算是 systemd-networkd，同样是会在 d-bus 中保存信息的

```bash
~]# busctl get-property org.freedesktop.network1 /org/freedesktop/network1/network/_310_2dnetplan_2dens3 org.freedesktop.network1.Network MatchName
as 1 "ens3"
~]# busctl get-property org.freedesktop.network1 /org/freedesktop/network1/network/_310_2dnetplan_2dens3 org.freedesktop.network1.Network MatchDriver
as 0
~]# busctl get-property org.freedesktop.network1 /org/freedesktop/network1/network/_310_2dnetplan_2dens3 org.freedesktop.network1.Network SourcePath
s "/run/systemd/network/10-netplan-ens3.network"
```

# Netplan 关联文件与配置

**/etc/netplan/** # netplan 读取 yaml 格式配置文件的路径
