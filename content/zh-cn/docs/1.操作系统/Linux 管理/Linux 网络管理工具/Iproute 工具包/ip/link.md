---
title: link
weight: 20
---

# 概述

> 参考：
>
> - [Manual(手册)，ip-link(8)](https://man7.org/linux/man-pages/man8/ip-link.8.html)

一个 **link** 代表一个 **network device(网络设备)**。link 对象及其相应的命令集，可以查看和操纵网络设备(增删改查等)。主要通过其自身的子命令来实现本身的功能。

网络设备配置

# Syntax(语法)

**ip \[OPTIONS] link \[COMMAND]**

**COMMAND：**

- **add|delete|set|show** # 增|删|改|查 一个虚拟 link
- **xstats** #
- **afstats** #
- **property** #

**OPTIONS：**

- **-4** # 指定使用的网络层协议是 IPV4
- **-r** # 显示主机时，不使用 IP，而是使用主机的域名

# add - 添加网络设备的虚拟链接

> Notes: 真实物理网卡对应的网络设备无法通过 ip link add 命令添加

**ip link add \[link DEVICE] \[ name ] NAME \[ARGS] type TYPE \[ ARGS ]**

- **DEVICE** # 要操作的物理设备
- **NAME** # 要操作的设备的名称
- **ARGS** # 这个参数可以设定设备的 IP 地址、网络地址、MTU 等
- **TYPE** # 设备类型
    - **bridge** # 以太网网桥设备
    - **bond** # Bonding(绑定)设备
    - **dummy** # 虚拟网络接口
    - **veth** # Virtual ethernet interface(虚拟以太网接口)设备
    - **vlan** # 802.1q tagged virtual LAN interface
    - **vxlan** # Virtual eXtended LAN
    - **ipip** # Virtual tunnel interface IPv4 over IPv4
    - 等等...... 所有可用的设备类型详见 Man 手册的 [Description 部分](https://man7.org/linux/man-pages/man8/ip-link.8.html#DESCRIPTION)

# set - 改变设备属性

> 注意: 如果请求多个参数更改，则在任何更改失败后，ip 立即中止。当 ip 可以将系统移动到不可预测的状态时，这是唯一的情况。解决方案是避免使用一个 ip 链路集调用更改几个参数。修饰符更改等效于 set。

```bash
ip link set { DEVICE | group GROUP } [ { up | down } ]
[ type ETYPE TYPE_ARGS ]
[ arp { on | off } ]
[ dynamic { on | off } ]
[ multicast { on | off } ]
[ allmulticast { on | off } ]
[ promisc { on | off } ]
[ protodown { on | off } ]
[ trailers { on | off } ]
[ txqueuelen PACKETS ]
[ name NEWNAME ]
[ address LLADDR ]
[ broadcast LLADDR ]
[ mtu MTU ]
[ netns { PID | NETNSNAME } ]
[ link-netnsid ID ]
[ alias NAME ]
[ vf NUM [ mac LLADDR ]
[ VFVLAN-LIST ]
[ rate TXRATE ]
[ max_tx_rate TXRATE ]
[ min_tx_rate TXRATE ]
[ spoofchk { on | off } ]
[ query_rss { on | off } ]
[ state { auto | enable | disable } ]
[ trust { on | off } ]
[ node_guid eui64 ]
[ port_guid eui64 ] ]
[ { xdp | xdpgeneric | xdpdrv | xdpoffload } { off |
object FILE [ section NAME ] [ verbose ] |
pinned FILE } ]
[ master DEVICE ]
[ nomaster ]
[ vrf NAME ]
[ addrgenmode { eui64 | none | stable_secret | random } ]
[ macaddr [ MACADDR ]
[ { flush | add | del } MACADDR ]
[ set MACADDR ] ]
```

# show - 显示设备属性

>[!Notes]
> show 命令无法显示网络设备的类型。想要查看网络设备的类型，可以通过与 [ethtool](/docs/1.操作系统/Linux%20管理/Linux%20网络管理工具/ethtool.md) 工具配置实现
>

```bash
for i in $(ip link show | awk -F: '$0 !~ "lo|vir|wl|^[^0-9]"{print $2;getline}'); do
  printf "%s\t%s\n" "${i}" "$(ethtool -i $i | grep driver)";
done | column -t
```

**ip link show \[ DEVICE | group GROUP ] \[ up ] \[ master DEVICE ] \[ type ETYPE ] \[ vrf NAME ]**
