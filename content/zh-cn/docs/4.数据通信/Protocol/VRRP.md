---
title: VRRP
linkTitle: VRRP
weight: 3
---

# 概述

> 参考：
>
> -

**Virtual Router Redundancy Protocol(虚拟路由冗余协议，简称 VRRP)** 是一种容错协议，其主要目的是解决路由单点故障的问题。VRRP 协议将局域网内的一组路由器虚拟为单个路由，通常将其称为一个路由备份组， 而这组路由器内包括一个 Master 路由（ 即活动路由器）和若干个 Backup 路由（即备份路由器）， VRRP 虚拟路由示意图如图 3-3 所示。在图 3-3 中 RouterA 、RouterB 和 RouterC 属于同一个 VRRP 组，组成一个虚拟路由器，而由 VRRP 协议虚拟出来的路由器拥有自己的 IP 地址 10.110.10.1 ，而备份组内的路由器也有自己的 IP 地址（如 Master 的 IP 地址为 10.110.10.5, Backup 的 IP 地址为 10.110.10.6 和 10.110.10.7）。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/cv6gcm/1616161487047-bb1bc9ee-9e3e-40a2-9a92-403553eb3520.jpeg)

1. 虚拟 IP：VIP，Virtual IP Address，在实际使用中，局域网内的主机仅仅知道这命虚拟路由器的 IP 地址 10 .110.10.1，而并不知道具体的 Master 路由器的 IP 地址以及 Backup 路由器的 IP 地址。局域网内的主机将自己的默认路由下一跳地址设置为该虚拟路由器的 IP 地址 10.110.10.1 之后，网络内的主机就通过这个虚拟的路由器来与其他网络进行通信。在通信过程中，如果备份组内的 Master 路由器故障，则 Backup 路由器将会通过选举机制重新选出一个新的 Master 路由器，从而继续向网络内的主机提供路由服务，最终实现了路由功能的高可用。

2. 此外，虚拟 IP 可以配置多个，比如 RouteA,B,C 的各个端口 1 绑定成一个 VIP1，RouteA,B,C 的各个端口 2 绑定成另一个 VIP2，RouteA,B,C 的各个端口 3 绑定成另一个 VIP3 以此类推，一组 VRRP 中，可以有多个 VIP，VIP1 的主路由是 Route1，VIP2 的主路由是 Route2，VIP3 的主路由是 Route3，各个 VIP 坏了，都可以有另外两台 Route 来代替工作，这样也解决了资源闲置问题，所有设备都是主，另外两台也是各个主设备的备设备

3. 虚拟 MAC 地址：如果某一时刻 VIP 在 RouteA 上，当 A 坏了之后，VIP 自动转移到了 RouteB 上，但是对于后端的 Host 来说，该 VIP 所对应的 MAC 地址已经改变了，是该 VIP 所配置的物理网卡的 MAC 地址，这意味着之前 Host 已经缓存下来了这个对应规则，则再发送数据的时候，还会使用原来的 MAC，发送给坏了的 RouteA。为了解决这个情况，则可以使用虚拟 MAC 地址的方案,使用一个虚拟的 MAC 与 VIP 绑定

4. ARP 欺骗：当一台新设备接入网络后，可以使用自问自答方式，广播问一下“网关在哪里”，然后自己广播回答“网关在这里”，这时候，这台设备就可以拿到网关的 IP 所对应的自己的端口的 MAC 地址了

5. 优先级：路由器开启 VRRP 功能后，会根据设定的优先级确定自己在备份组中的角色。优先级高的路由器成为 Master 路由器，优先级低的成为 Backup 路由器，并且 Master 路由器定期发送 VRRP 通告报文，通知备份组内的其他 Backup 路由器自己工作正常， 而备用路由器则启动定时器等待通告报文的到来。（如何判断 Master 路由器是否正常工作？）如果 Backup 路由器的定时器超时后仍未收到 Master 路由器发送来的 VRRP 通告报文， 则认为 Master 路由器已经故障，此时 Backup 路由器会认为自己是主用路由器（备份组内的路由器会根据优先级选举出新的 Master 路由器），并对外发送 VRRP 通告报文。

6. 此外， VRRP 在提高路由可靠性的同时，还简化了主机的路由配置， 在具有多播或广播能力的局域网中，借助 VRRP 能在某台路由器出现故障时仍然提供高可靠的默认链路，有效避免单一链路发生故障后网络中断的问题，并且无需修改主机动态路由协议、路由发现协议等配置信息。

# 试验

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/cv6gcm/1616161487021-bdc76b76-08cd-4480-b281-3230f6e829a3.jpeg)

PC11 在 VLAN10 中 192.168.10.11，PC21 在 VLAN20 中 192.168.20.21，如果让两台设备互通，那么需要通过三层路由实现

为了实现路由冗余，在网关所在设备启 VRRP 协议，保证一个路由（网关）坏掉的情况，可以有另一个代替。

那么首先需要在 VRRP master 和 VRRPbackup 设备上配置 VLAN10 和 VLAN20 的虚拟网关 IP，以及给每个 VLAN 配置一个 IP 以启动该 VLAN 接口，如图所示 （需要设置优先级 priority 以确保其中一台设备始终为 MASTER，优先级默认为 100）

VRRPmaster 的配置

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/cv6gcm/1616161487070-34b33cb9-f96b-4142-934b-ad1261217e5d.jpeg)

VRRPbackup 的配置

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/cv6gcm/1616161487028-8f6c1dca-c7d6-4dda-92cb-d00a1251a495.jpeg)

其余基础配置：

1. port link-type trunk

           port trunk allow-pass vlan 10 20       交换机互联端口为trunk切允许相应vlan通过

2. port link-type access

           port default vlan 10                           指定接入层交换机接入设备端口的VLAN以实现隔离

两个路由之间互相发送报文保证双方是主还是备（注：当 VRRP 两个设备互相无法收到验证报文的时候，会出现两台都是 MASTER 的情况）

为了实现更高的冗余效果，可以在两台 VRRP 设备上增加级联线（即互联两台设备）该极连线采取链路聚合原则（上图的 23，24 口），具体配置如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/cv6gcm/1616161487040-d85f0bfa-26c8-45f4-85a8-a7ab3a98234f.jpeg)

并且使用 interface eth-trunk 1

              trunkport gigabitethernet 0/0/x to 0/0/y      该命令，创建聚合端口并且把X到Y的端口都加到聚合组里去
