---
title: NetworkManager
weight: 1
---

# 概述

> 参考：
>
> - [GitHub 项目，NetworkManager/NetworkManager](https://github.com/NetworkManager/NetworkManager)
> - [GitLab 项目，freedesktop-NetworkManager/NetworkManager](https://gitlab.freedesktop.org/NetworkManager/NetworkManager)
> - [Manual(手册),NetworkManager.conf(5)](https://networkmanager.dev/docs/api/latest/NetworkManager.conf.html)
> - [官网](https://networkmanager.dev/)

**NetworkManager daemon** 是管理网络的守护进程。该守护进程尝试通过管理主网络连接和其他网络接口（如以太网，WiFi 和移动宽带设备），使网络配置和操作尽可能轻松自动。 除非禁用该行为，否则 NetworkManager 将在该设备的连接可用时连接任何网络设备。 有关网络的信息通过 D-Bus 接口导出到任何感兴趣的应用程序，提供丰富的 API，用于检查和控制网络设置和操作。

# Connection

> 参考：
>
> - [Manual(手册),nm-settings-nmcli(5)](https://networkmanager.dev/docs/api/latest/nm-settings-nmcli.html)
> - [Manual(手册),nm-settings-dbus(5)](https://networkmanager.dev/docs/api/latest/nm-settings-dbus.html)
> - [Manual(手册),nm-settings-keyfile(5)](https://networkmanager.dev/docs/api/latest/nm-settings-keyfile.html)
> - [Manual(手册),nm-settings-ifcfg-rh(5)](https://networkmanager.dev/docs/api/latest/nm-settings-ifcfg-rh.html)

NetworkManager 将所有网络配置抽象成 **Connection(连接)**，这些 Connection 的配置中包含网络配置(比如 IP 地址、网关等)。当 NetworkManager 激活网络设备上的 Connection 时，将为这个网络设备应用配置文件中的内容，并建立活动的网络连接。所以，可以创建多个 Connection 来关联到一个网络设备上；这样，它们就可以灵活地具有用于不同网络需求的各种网络配置。

**用白话说就是：Connection 就是“网络配置”，网络设备(device)关联并使用“网络配置”来实现联网。而 NetworkManager 就是管理这些 Connection 的。Connection 可以表示一个概念，也可以表示一个配置文件。**

## Connection 插件

NetworkManager 通过 **Plugins(插件)** 的方式来管理 Connection 配置文件。在不同的 Linux 发行版中，所使用的插件各不相同，但是默认情况下，NetworkManager 使始终启用名为 **keyfile** 的插件，这是一个通用插件，当其他插件无法支持某些类型的 Connection 配置时，keyfile 插件将会自动提供支持。keyfile 插件会将 Connection 文件保存到 /etc/NetworkManager/system-connections/、/usr/lib/NetworkManager/system-connections/、/run/NetworkManager/system-connections/ 这三个目录中。

可以在 /etc/NetworkManager/NetworkManager.conf 文件中配置想要使用的插件，插件用于读写系统范围的连接配置文件。当指定多个插件时，将从所有列出的插件中读取 Connections。写入 Connections 时，会要求插件按照此处列出的顺序保存连接；如果第一个插件无法写出该连接类型（或无法写出任何连接），则尝试下一个插件。如果没有插件可以保存连接，则会向用户返回错误。

可用插件的数量是特定于发行版的。所有可用的插件详见 [Manual(手册) 中 Plugins 章节](https://networkmanager.dev/docs/api/latest/NetworkManager.conf.html#settings-plugins)

**keyfile**

- keyfile 插件是支持 NetworkManager 拥有的所有连接类型和功能的**通用插件**。它以 .ini 格式在 /etc/NetworkManager/system-connections 文件中写入连接配置。
  - 有关文件格式的详细信息，请参阅 nm-settings-keyfile(5)。
- keyfile 插件存储的连接文件可能包含纯文本形式的 passwords、secrets、private keys，因此它将仅对 root 用户可读，并且插件将忽略除 root 用户或组之外的任何用户或组可读或可写的文件。
  - 有关如何避免以纯文本形式存储密码，请参阅 nm-settings(5) 中的“秘密标志类型”。
- 此插件始终处于活动状态，并将自动用于存储其他插件不支持的连接。

**ifcfg-rh**

- 此插件用于 Fedora 和 Red Hat Enterprise Linux 发行版，用于从标准 /etc/sysconfig/network-scripts/ifcfg-\* 文件读取和写入配置。它目前支持读取 Ethernet, Wi-Fi, InfiniBand, VLAN, Bond, Bridge, Team 这几种类型的连接。启用 ifcfg-rh 隐式启用 ibft 插件(如果可用)。这可以通过添加 no-ibft 来禁用。
- 有关 ifcfg 文件格式的更多信息，请参见 /usr/share/doc/initscripts/sysconfig.txt 和 nm-settings-ifcfg-rh(5)。

**ifupdown**

- This plugin is used on the Debian and Ubuntu distributions, and reads Ethernet and Wi-Fi connections from /etc/network/interfaces.
- This plugin is read-only; any connections (of any type) added from within NetworkManager when you are using this plugin will be saved using the keyfile plugin instead.

**ibft, no-ibft**

- These plugins are deprecated and their selection has no effect. This is now handled by nm-initrd-generator.

**ifcfg-suse, ifnet**

- These plugins are deprecated and their selection has no effect. The keyfile plugin should be used instead.

## Connection D-Bus

NetworkManager 还会将这些 Connection 配置导出到 D-Bus 上，比如，通过 **busctl** 命令，可也获取 Connection 中的内容：

```bash
[root@ansible dispatcher.d]# busctl get-property org.freedesktop.NetworkManager /org/freedesktop/NetworkManager/Devices/2 org.freedesktop.NetworkManager.Device Interface
s "ens33"
```

所以，真正的底层实现，是通过 D-bus 中的网络设备配置文件来实现的

## Connection 关联文件

默认情况下，由 **keyfile 插件**管理 **INI 格式**的 Connection 配置文件。并默认保存在 /etc/NetworkManager/system-connections/ 目录中。

> 注意：
>
> - 在 RedHad 相关的发行版中，NetworkManager 会运行名为 ifcfg-rh 的插件，插件会将 /etc/NetworkManager/system-connections/ 目录中的 Connection 配置文件翻译成老式配置文件格式，并保存在 /etc/sysconfig/network-scripts/ 目录中
> - **所以，在 RedHad 中，是无法从 /etc/NetworkManager/system-connections/ 目录中找到连接配置文件**
> - 若想禁用 ifcfg-rh 插件，只需要在 /etc/NetworkManager/NetworkManager.conf 文件中的 main 部分添加 plugins=keyfile 即可

在 D-Bus API 上的 Connection 配置中，将 INI 中的 **Sections(部分) 称为 Settings(设置)**，Setting 即是 **Properties(属性)** 的集合。所以，很多文档，都将 Connection 表示为一组特定的、封装好的、独立的 **Settings(集合)** 集合。Connection 由一个或多个 Settings 组成。

**Settings**用于描述一个 Connection。每个 Setting 都具有一个或多个 `**Property(属性)**` 。Setting 与 Property 中间以点 `.` 连接。每个 Setting.Property 都会有一个值。

一个 Connection 有哪些 Settings，Setting 又有哪些 Property，以及这些 Property 都有什么作用，详见 [Connection 配置详解](/docs/1.操作系统/Kernel/Network/Linux%20网络栈管理/NetworkManager/Connection%20配置详解.md)

下面的命令，可以从 D-Bus API 中获取配置文件所在路径

```bash
# CentOS 中使用 ifcfg-rh 插件
~]# busctl get-property org.freedesktop.NetworkManager /org/freedesktop/NetworkManager/Settings/4 org.freedesktop.NetworkManager.Settings.Connection Filename
s "/etc/sysconfig/network-scripts/ifcfg-enp25s0f3"

# CentOS 中不使用 ifcfg-rh 插件
~]# busctl get-property org.freedesktop.NetworkManager /org/freedesktop/NetworkManager/Settings/4 org.freedesktop.NetworkManager.Settings.Connection Filename
s "/etc/NetworkManager/system-connections/eth1"
```

可以看到，使用不同的插件，配置文件所在路径是不同的

**用白话说：如果说 Connection 是一个配置文件的话，Setting 就是配置文件中的 `context(配置段，或称为"配置环境")`，`Property(属性)` 是该配置环境下的 `keyword(关键字,或称为"键"、"字段")`**。所以，一般情况下，Connection 也可以描述为由一个或多个 Property(属性) 组成。我们都把 Setting.Property 简称为 属性。**其实 Setting 就是很多产品的配置文件中的 Context**。

### 配置文件示例

```bash
~]# cat /etc/NetworkManager/system-connections/ens3.nmconnection
[connection]
id=ens3
uuid=8f8541bc-4893-418b-98d4-fbc7433747cf
type=ethernet
interface-name=ens3
permissions=

[ethernet]
mac-address-blacklist=

[ipv4]
address1=172.19.42.248/24,172.19.42.1
dns-search=
method=manual

[ipv6]
addr-gen-mode=stable-privacy
dns-search=
method=auto

[proxy]
```

如果通过 nmcli 命令查看这个 Connection，格式如下：

```bash
~]# nmcli connection show eth0
connection.id:                          ens3
connection.uuid:                        8f8541bc-4893-418b-98d4-fbc7433747cf
connection.type:                        802-3-ethernet
connection.interface-name:              eth0
.........
ipv4.method:                            manual
ipv4.dns:                               223.5.5.5
ipv4.addresses:                         172.19.42.248/24
ipv4.gateway:                           172.19.42.1
.......
```

第一列中的 connection 与 ipv4 就是 Setting。其中 id、uuid、type、interface-name 都是 connection 这个 Setting 的 Property，而 method、dns 等等都是 ipv4 这个 Setting 的 Property。第二列就是同一行 Property 对应的值。

## NetworkManager API

NetworkManager 提供了一个 API，用来管理 Connection、检查网络配置等。[nmcli 命令行工具](/docs/1.操作系统/Kernel/Network/Linux%20网络栈管理/NetworkManager/nmcli%20命令行工具.md) 是官方提供的用于使用 API 的客户端应用程序。

> 也可以手动管理 Connection 文件，就跟出现 NetworkManager 之前一样，手动配置 /etc/sysconfig/network-scripts 目录下的网络设备配置文件，然后重启 deamon 进程以便加载这些文件即可。

注意：

1. 一个网络设备(device)可以关联多个 connection，但是同一时间只能有一个与该网络设备(device)关联 connection 处于 active 状态。这就可以让一个网卡(device)同时具备多个配置，可以随时切换。
2. NetworkManager 默认不会识别到配置文件的更改 并会继续使用旧的配置数据。如果更改 /etc/NetworkManager/system-connections/ 目录下的配置文件，那么需要让 NetworkManager 再次读取已经改动过的配置文件，如果想要确保这件事，需要执行如下几条命令
   1. nmcli connection reload # 让 Connection 重新加载以读取配置文件
   2. nmcli connection up ConnectionName # 再次启动指定的 Connection，这里的 up 也有 restart 的意思

# NetworkManager 关联文件

**/etc/NetworkManager/** #

- **./conf.d/** # 类似 include 功能，是 NetworkManager.conf 文件的内容片段。
- **./NetworkManager.conf** # NetworkManager 程序的运行时配置文件
- **./system-connections/** # 每个 Connection 的配置文件保存路径。
  - 在 RedHad 中，该路径被修改到 /etc/sysconfig/network-scripts/ 上去了。

**/run/NetworkManager/** #

- **./system-connections/** # 自动生成的 Connection 的配置文件保存路径。

**/usr/lib/NetworkManager/** #

- **./system-connections/** #

# 常见问题

## LACP 在 NetworkManager 管理的 Bonding 不工作

<https://github.com/systemd/systemd/issues/15208>

当 systemd 版本在 242、243、245 时，NetworkManager 对于 802.3ad 模式的 Bonding 在发送 LACP 包是可能会产生异常

如果通过 NetworkManager 创建的 Bond 网络设备失效，有如下几种可用的解决方式：

- 通过 ip 命令先删除网络设备，再通过 ip 命令添加即可。
  - ip link set bond1 down
  - ip link del bond1
  - ip link add bond1 type bond mod 802.3ad
