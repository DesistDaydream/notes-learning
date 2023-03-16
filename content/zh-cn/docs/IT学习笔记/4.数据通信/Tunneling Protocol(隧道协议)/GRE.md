---
title: GRE
---

# 概述

**Generic Routing Encapsulation(通用路由封装，简称 GRE)**是一种隧道协议，在数据两端，对数据进行封装和解封装。是 Cisco Systems 开发的隧道协议，可以通过 Internet 协议网络将各种网络层协议封装在虚拟点对点链路中。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qp0rg3/1616160946635-68a0d422-0333-48eb-aa84-f8f03f1c76a0.jpeg)
如上图，当从本机想要通过隧道发送数据时，会通过 GRE 模块进行封装，然后把对外通信的 IP 地址当做 GRE 的外部地址，封装在最外成变成新的 IP。此时，GRE 还会有一个内部地址用来与隧道的对端进行协商，以便识别公网上的隧道两端的设备

## Linux 下实现 GRE

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qp0rg3/1616160946644-61636d63-9002-4967-b9e5-98fb974f4fb7.jpeg)
场景一：
如果 Linux1 想要访问 10.10.2.0/24 网段的设备，这时候可以使用 GRE 建立隧道来实 现。比如在 Linux1 上想访问 10.10.2.100 这台设备，把数据包通过手动添加路由的方式，直接送给 tun 设备，tun 设备会直接把数据包发送给其隧道的对端(i.e.Linux2)，当 Linux2 收到这个数据包时，会发现 GRE 的包头，并解开 GRE 后发现真实的目的地址(10.10.2.100)，这时候可以在 Linux2 上做一个 snat，指明源地址 172.16.0.1/32 且目的地址是 10.10.2.0/24 网段的数据包全部把源地址转换成 10.10.2.1(如果不做 snat，源地址是 172.16.0.1，在回包的时候，是无法回去的，否则再添加其余路由条目)

场景二：
10.10.0.1/24 与 10.10.1.1/24 作为网段的网关使用，可以让两边的内网机器，直接互相访问。同样需要手动添加路由，目的地址是对端内网网段的 IP 的数据包送给 tun 设备。

注意：在 Linux1 上 ping172.16.0.2 的时候，在 Linux2 上抓包的话，会显示源地址是 172.16.0.1。

## 配置方式：

通过配置文件来进行配置

| Linux1 的 tun 设备的配置文件：
DEVICE=tun0
ONBOOT=yes
TYPE=GRE
PEER_OUTER_IPADDR=100.0.2.100
PEER_INNER_IPADDR=172.16.0.2/30
MY_OUTER_IPADDR=100.0.1.100
MY_INNER_IPADDR=172.16.0.1
KEY=lichenhao
BOOTPROTO=none
配置完成后会自动下面的路由信息
172.16.0.2/30 dev tun0 proto kernel scope link src 172.16.0.1/30 | Linux2 的 tun 设备的配置文件：
DEVICE=tun0
ONBOOT=yes
TYPE=GRE
PEER_OUTER_IPADDR=100.0.1.100
PEER_INNER_IPADDR=172.16.0.1/30
MY_OUTER_IPADDR=100.0.2.100
MY_INNER_IPADDR=172.16.0.2
KEY=lichenhao
BOOTPROTO=none
配置完成后会自动下面的路由信息
172.16.0.1/30 dev tun0 proto kernel scope link src 172.16.0.2/30 |
| --- | --- |

如果想要实现场景一的功能,可以在 Linux1 上添加路由条目：ip route add 10.10.2.0/24 dev tun0。在 Linux 添加 snat 规则：iptables -t nat -A POSTROUTING -s 1772.16.0.1 -d 10.10.2.0/24 -j SNAT --to-source 10.10.2.1

配置文件说明：

- PEER_OUTER_IPADDR=100.0.2.100 # 隧道外部对端地址。i.e.能让两台服务器互相访问的对端地址
- PEER_INNER_IPADDR=172.16.0.2/30 # 隧道内部对端地址。i.e.对端服务器的 tun 设备的 IP
- MY_OUTER_IPADDR=100.0.1.100 # 隧道外部本地地址。i.e.能让两台服务器互相访问的本地地址
- MY_INNER_IPADDR=172.16.0.1 # 隧道内部本地地址。i.e.本服务器的 tun 设备的 ip

从命令行来体现的话，PEER 就是命令行中的 remote 和 peer。MY 就是 local

通过 ip 命令来进行配置

Linux1 上的命令
ip tunnel add tun0 mode gre remote 100.0.2.100 local 100.0.1.100
ip addr add 172.16.0.1/30 peer 172.16.0.2/30 dev tun0
ip route add 10.10.2.0/24 dev tun0

Linux2 上的命令
ip tunnel add tun0 mode gre remote 100.0.2.100 local 100.0.1.100
ip addr add 172.16.0.1/30 peer 172.16.0.2/30 dev tun0
ip route add 10.10.1.0/24 dev tun0

通过 NetworkManager 来进行配置

Linux1 上的命令
nmcli connection add con-name tun0 ifname tun0 type ip-tunnel mode gre ip-tunnel.remote 100.0.2.100 ip-tunnel.local 100.0.1.100 ipv4.method manual ipv4.addresses 172.16.0.1/30
ip route add 10.10.2.0/24 dev tun0

Linux2 上的命令
nmcli connection add con-name tun0 ifname tun0 type ip-tunnel mode gre ip-tunnel.remote 100.0.1.100 ip-tunnel.local 100.0.2.100 ipv4.method manual ipv4.addresses 172.16.0.2/30
ip route add 10.10.1.0/24 dev tun0
