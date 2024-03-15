---
title: Connection 配置详解
---

# 概述

> 参考：
>
> - [Manual(手册),nm-settings-nmcli(5)](https://networkmanager.dev/docs/api/latest/nm-settings-nmcli.html) # 这个 man 手册中，可以看到每个 Setting 中都有哪些 Property 以及这些 Property 的作用。
> - [Manual(手册),nm-settings-dbus(5)](https://networkmanager.dev/docs/api/latest/nm-settings-dbus.html) # 这里有 Property 的默认值
> - [Manual(手册),nm-settings-keyfile(5)](https://networkmanager.dev/docs/api/latest/nm-settings-keyfile.html)
> - [Manual(手册),nm-settings-ifcfg-rh(5)](https://networkmanager.dev/docs/api/latest/nm-settings-ifcfg-rh.html)
> - 在 [GNOME 开发者中心官网](https://developer-old.gnome.org/NetworkManager/)中，也可以查到 Manual

Connection 配置文件默认由 keyfile 插件管理，是类似 **INI 格式**的。同时配置文件还会保存在 D-Bus 中。

在 D-Bus 中，NetworkManager 将 INI 中的 Sections(部分) 称为 **Settings(设置)**，Setting 多个是 **Properties(属性)** 的集合。所以，很多文档，都将 Connection 表示为一组特定的、封装好的、独立的 **Settings(集合)**。Connection 由一个或多个 Settings 组成。

```bash
# 我启动了一个 连接
~]# nmcli con up bridge-slave-bond0
Connection successfully activated (master waiting for slaves) (D-Bus active path: /org/freedesktop/NetworkManager/ActiveConnection/16)
# 从 D-Bus 的路径中可以看到这些信息
~]# busctl get-property org.freedesktop.NetworkManager /org/freedesktop/NetworkManager/ActiveConnection/16 org.freedesktop.NetworkManager.Connection.Active
Connection      Default6        Dhcp4Config     Id              Ip6Config       SpecificObject  StateFlags      Uuid
Default         Devices         Dhcp6Config     Ip4Config       Master          State           Type            Vpn
# 上面的信息是按 TAB 补全出来的
~]# busctl get-property org.freedesktop.NetworkManager /org/freedesktop/NetworkManager/ActiveConnection/16 org.freedesktop.NetworkManager.Connection.Active  Id
s "bridge-slave-bond0"
```

注意：

- 在 [nmcli 命令行工具](/docs/1.操作系统/Kernel/Network/Linux%20网络栈管理/NetworkManager/nmcli%20命令行工具.md) 命令中使用 SETTING.PROPERTY 时，如果 SETTING 和 PROPERTY 是 唯一的，则可以使用 `Alias(别名)`。
  - 比如 connection.typ 的别名为 type。其实就是缩写，简化操作。
- 802-3-ethernet 类型就是一般物理网卡，通常使用 ethernet 别名表示。不管 ethernet 是作为 SETTING 还是作为 VALUE，都可以使用别名。但是在 -f 选项中，不能使用别名。
- **在 RedHad 中，是无法从 /etc/NetworkManager/system-connections/ 目录中找到连接配置文件，这是因为 RedHad 系发行版使用的是 ifcfg-rh 插件**
- 配置文件可以通过命令行修改，比如 `nmcli con add connection.type XXX`，其中 connection 就是 connection 这个 SETTING，type 就是 该 SETTING 下的一个属性
  - 所以，命令行中的 **SETTING.PROPERTY** 与配置文件是**一一对应的**。

# connection SETTING

通用的 connection 配置设置

**autoconnect=\<BOOLEAN>** # 别名 autoconnect。该连接是否自动连接。

**id=\<STRING>** # 别名 con-name。该连接的名称。若不指定，则会默认生成一个

**interface-name STRING** # 别名 ifname。该连接绑定的网络设备名称。

**master=\<STRING>** # 别名 master。该连接的主设备的 name 或 UUID。具有 master 属性的连接将会降级为从设备.常用于向 bond 或者 brdige 设备中添加从设备时使用。

- 在使用 nmcli 命令时如果使用 master 别名，则会自动为连接添加 slave-type 属性，属性根据主设备的类型决定。如果不使用别名，则需要显式得使用 connection.slave-type 来指定该连接的从属类型。
- 若主设备状态 down，则该从设备状态变为 lowerlayerdown

**type=\<STRING>** # 别名 type。连接类型，常用的有 ethernet、bridge、bridge-slave、bond、bond-slave、tun 等等。其实就是要添加的连接的网络设备的类型。**必选，每个连接必须有一个 type**

# ipv4 SETTING

用于配置一个 connection 的 IPv4 信息

**address=<\[]UNIT32>** # 别名 ip4。指定该连接的 IP 地址。可以使用 192.168.0.0/24 这种格式

- 若有多个 IP 地址，则

**dns=<\[]UINT32>** # DNS 服务端 IP 地址列表。多个地址以 `;` 分割。

**dns-search=\<STRING>** # 待补充

**gateway=\<IP>** # 别名 gw4。指定该连接的网关地址。

**method=\<METHOD>** # 该连接的 ipv4 获取方法。即通过 dhcp 获取还是手动指定等等

- auto # ipv4 地址可以自动获取(通过 dhcp、pp 等)。默认值
- disabled # 在此连接上不使用 ipv4。Note：如果在配置文件中不指定 ipv4 信息，那么在 reload 时，连接的 ipv4.method 也会变成该值
- link-local # 为该连接分配在 169.254/16 范围的本地地址
- manual # 手动指定 IP。如果使用此方法，则必须指定 ipv4.address 属性。
- shared # 表示此连接将提过对其他计算机的网络访问权限，为接口分配一个 10.42.x.1/24 的地址，并启动 dhcp 和 dns 转发服务，并且该接口为 Nat-ed 到当前的默认网络连接

**routes=\<ROUTEs>** # 设定作用在该连接上的路由条目，多个条目以逗号分割

**routing-rule=\<RULEs>** # 设定路由策略。

# ipv6 SETTING

用于配置一个 connection 的 IPv6 信息

**method=\<METHOD>** # 与 ipv4.method 基本相同

# bond SETTING

**options=\<map\[STRING]STRING>** # bond 选项。以逗号分隔的键值对，每个键值对都对应 bond 的一个选项和其值。`默认值：{'mode':'balance-rr'}`

# vlan SETTING

**id=\<NUM>** # 别名 id。此连接关联的网络设备的 VLAN 标识符。有效范围是 0 到 4094

**parent=\<STRING>** # 别名 dev。VLAN 网络设备的父设备的 name 或 UUID。即指定在哪个网络设备上附加 VLAN 标识符。

```bash
~]# nmcli con add type vlan con-name vlan2-eth0 dev eth0  id 2 ipv4.method disabled ipv6.method disabled
Connection 'vlan2-eth0' (b54bd8d9-de8a-4e26-a579-8e9ff95f126f) successfully added.
~]# nmcli c s
NAME        UUID                                  TYPE      DEVICE
....
vlan2-eth0  0a557c85-f98a-44c3-a2a5-31645efb98b9  vlan      eth0.2
~]# ip a s
....
165: eth0.2@eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default qlen 1000
    link/ether 52:54:00:6a:86:89 brd ff:ff:ff:ff:ff:ff
```

# 802-3-ethernet SETTING

别名 ethernet

有线以太网的配置

**mtu=\<UINT32>** # 别名 mtu。连接关联的物理设备的 MTU 的值。`默认值：auto`

# RedHad 中 Connection 与 老式配置 对应关系

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/qpdcnf/1624352815778-2fd45e29-fc7b-495f-90e7-0bc0c2691a27.png)

## 看一下老式配置的内容

DEVICE=eth0 # 指明与此配置文件相关联的网络设备名称

TYPE=Ethernet # 指明该网络设备的类型

BOOTPROTO=none # 指明设备获取 ip 的方法，有 atuo、none、disabled 等

- auto # 自动获取 ip(通过 dhcp 等方法)
- none # 没有方法，也就是说手动指定 IP。可以手动指定 IPADDR、PREFIX、GATEWAY 等
- disabled # 不获取 ip，如果不指定 IPADDR、NETMASK、GATEWAY 的话，默认就是 disabled

IPADDR=192.168.10.22 # 指定该设备 IP

NETMASK=255.255.255.0 # 指定该设备的掩码

- PREFIX=24 # 也可是使用这种方式来表示掩码

GATEWAY=192.168.10.2 # 指定该设备的网关

DEFROUTE=yes # 该网卡的路由是否设置为默认路由

ONBOOT=yes|no # 启动网卡时是否自动加载该配置文件

USERCTL={yes|no} # 非 root 用户是否可以控制该设备

PROXY_METHOD=none

BROWSER_ONLY=no

IPV4_FAILURE_FATAL=no

IPV6INIT=yes

IPV6_AUTOCONF=yes

IPV6_DEFROUTE=yes

IPV6_FAILURE_FATAL=no

IPV6_ADDR_GEN_MODE=stable-privacy

NAME=eth0 # 用于 NetworkManager 服务中，作为连接名字。

UUID=74b13ff5-697f-43cf-abdc-27389b57ecbe

ZONE=public

# keyfile 配置文件示例

## 一个以太网的配置

```bash
[connection]
id=Main eth0
uuid=27afa607-ee36-43f0-b8c3-9d245cdc4bb3
type=802-3-ethernet
autoconnect=true

[ipv4]
method=auto

[802-3-ethernet]
mac-address=00:23:5a:47:1f:71
```

## 简单的桥设备配置

桥设备

```bash
[connection]
id=MainBridge
uuid=171ae855-a0ab-42b6-bd0c-60f5812eea9d
interface-name=MainBridge
type=bridge

[bridge]
interface-name=MainBridge
```

桥的从设备

```bash
[connection]
id=br-port-1
uuid=d6e8ae98-71f8-4b3d-9d2d-2e26048fe794
interface-name=em1
type=ethernet
master=MainBridge
slave-type=bridge
```

VLAN 设备

```bash
[connection]
id=VLAN for building 4A
uuid=8ce1c9e0-ce7a-4d2c-aa28-077dda09dd7e
interface-name=VLAN-4A
type=vlan

[vlan]
interface-name=VLAN-4A
parent=eth0
id=4
```
