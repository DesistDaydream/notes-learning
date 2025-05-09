---
title: link
linkTitle: link
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

# 应用示例

- 显示链路详细信息，包括接收与发送的数据包数以及错误数，丢弃数等
  - ip -s link show
- 查看所有 bridge 类型的网络设备
  - ip link show type bridge
- 启动或者停止 eth0 网卡，可以简写为 `ip l s eth0 up`
  - ip link set dev eth0 up|down
- 添加名字为 veth1.1 的链路，类型为 veth，veth 的另一半名字为 veth1.2
  - ip link add veth1.1 type veth peer name veth1.2
- 把 veth1.1 这个设备绑定到名为 r1 的 namespace 上(注意：一个网络设备只能绑定到一个 namespace 上，一个 namesapce 可以绑定多个网络设备)
  - ip link set veth1.1 netns r1
- 将 vnet0 设备绑定到 br0 桥上
  - ip link set dev vnet0 master br0
- 创建 Bond 类型网络设备
  - 创建 802.3ad 模式的 Bond 类型网络设备
    - ip link add bond1 type bond mod 802.3ad
  - 将物理网卡关联到的网络设备先关掉，再添加到 Bond 网络设备中
    - ip link set enp6s0f0 master down
    - ip link set enp6s0f0 down
    - ip link set enp6s0f0 master bond1
    - ip link set enp6s0f1 down
    - ip link set enp6s0f1 master bond1
