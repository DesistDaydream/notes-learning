---
title: IPVS
linkTitle: IPVS
weight: 3
---

# 概述

> 参考：
>
> - [Wiki, IPVS](https://en.wikipedia.org/wiki/IP_Virtual_Server)
> - [官方文档](http://www.linuxvirtualserver.org/software/ipvs.html)

**IP Virtual Service(IP 虚拟服务，简称 IPVS)** 是基于 [Netfilter](/docs/1.操作系统/Kernel/Network/Linux%20网络流量控制/Netfilter/Netfilter.md) 的 Linux 内核模块，用来实现 [LVS](/docs/3.集群与分布式/LVS/LVS.md) 集群中的 **Scheduler(调度器)** 功能。启动这个模块的 Linux 服务器就变成了 LVS 系统中的 **Director**，此时，这个服务器可以看作是一种高效的 Layer-4(四层) 交换机。在 Director 上运行 IPVS 代码是 LVS 的基本要素。

IPVS 在服务器上运行，并充当 RS 集群前面的负载均衡器。IPVS 可以将基于 TCP 和 UDP 的服务请求定向到真实服务器，并使真实服务器的服务在单个 IP 地址上表现为虚拟服务。当一个 TCP 连接的初始 SYN 报文到达时，IPVS 就选择一台服务器，将报文转发给它。此后通过查发报文的 IP 和 TCP 报文头地址，保证此连接的后继报文被转发到相同的服务器。这样，IPVS 无法检查到请求的内容再选择服务器，这就要求后端的服务器组是提供相同的服务，不管请求被送到哪一台服务器，返回结果都应该是一样的。但是在有一些应用中后端的服务器可能功能不一，有的是提供 HTML 文档的 Web 服务器，有的是提供图片的 Web 服务器，有的是提供 CGI 的 Web 服务器。这时，就需要基于内容请求分发 (Content-Based Request Distribution)，同时基于内容请求分发可以提高后端服务器上访问的局部性。

- 一个 ipvs 主机可以同时定义多个 cluster service
- 一个 cluster service 上至少应该定义一个 real server，定义时指明 lvs-type，以及 lvs scheduler

用白话理解 IPVS：
IPVS 就是包括 Director 和 RS 在内的所有设备上的 IP，统一虚拟成一个 IP，这个 IP 就是面向用户的唯一 IP，用户通过这个 IP，就可以访问集群，让集群为其提供服务，这也是负载均衡的体现，也是集群的体现，把很多设备当做一个整体来看。

## IPVS 与 LVS 的关系

LVS 更偏向于描述一个概念，而 IPVS 程序则是实现 LVS 的最核心部分。通过 IPVS 以及其管理工具 ipvsadm，可以实现 LVS 中的 Director(指挥器)。而 RS，本质上并不需要 LVS 或者 IPVS 代码支持，只需要在 DR 模式下，配置一些内核参数即可。

而随着发展，IPVS 已经存单独的程序，被包含在 Linux 内核中，成了了默认自带的模块。

可以这么说，IPVS 就是 LVS；也可以说，LVS 包含 ipvs 与 ipvsadm。

# IPVS 配置

ipvs 可以通过两种方式进行配置：

- ipvsadm 命令
- ipvs 模块参数

ipvs 是一个内核模块，所以，想要配置 ipvs 则需要以[内核模块](/docs/1.操作系统/Kernel/Linux%20Kernel/Module(模块).md)的配置方式来进行配置。可以通过 modinfo -p ip_vs 命令查看该模块可以配置的参数

```bash
~]# modinfo -p ip_vs
conn_tab_bits:Set connections' hash size (int)
```

现阶段，可用的参数只有一个：

- **conn_tab_bits** # 设置连接表的大小。`默认值：12`。

- 该参数控制下面示例中 size 的大小，2 的 12 次方，4096

```bash
root@desistdaydream:~# ipvsadm -ln
IP Virtual Server version 1.2.1 (size=4096)
Prot LocalAddress:Port Scheduler Flags
  -> RemoteAddress:Port           Forward Weight ActiveConn InActConn
```

IPVS connection hash table size，该表用于记录每个进来的连接及路由去向的信息（这个和 iptables 跟踪表类似）。连接的 Hash 表要容纳几百万个并发连接，任何一个报文到达都需要查找连接 Hash 表。Hash 表的查找复杂度为 O(n/m)，其中 n 为 Hash 表中对象的个数，m 为 Hash 表的桶个数。当对象在 Hash 表中均匀分布和 Hash 表的桶个数与对象个数一样多时，Hash 表的查找复杂度可以接近 O(1)

连接跟踪表中，每行称为一个 hash bucket（hash 桶），桶的个数是一个固定的值 CONFIG_IP_VS_TAB_BITS，默认为 12（2 的 12 次方，4096）。这个值可以调整，该值的大小应该在 8 到 20 之间。

LVS 的调优建议将 hash table 的值设置为不低于并发连接数。例如，并发连接数为 200，Persistent 时间为 200S，那么 hash 桶的个数应设置为尽可能接近 200x200=40000，2 的 15 次方为 32768 就可以了。当 ip_vs_conn_tab_bits=20 时，哈希表的的大小（条目）为 pow(2,20)，即 1048576。

这里的 hash 桶的个数，并不是 LVS 最大连接数限制。LVS 使用哈希链表解决“哈希冲突”，当连接数大于这个值时，必然会出现哈稀冲突，会（稍微）降低性能，但是并不对在功能上对 LVS 造成影响。

修改模块参数：echo "options ip_vs conn_tab_bits=22" > /etc/modprobe.d/ip_vs.conf，效果如下：

```bash
IP Virtual Server version 1.2.1 (size=4194304)
Prot LocalAddress:Port Scheduler Flags
  -> RemoteAddress:Port           Forward Weight ActiveConn InActConn
TCP  10.10.9.60:30000 rr persistent 30
  -> 10.10.9.69:30000             Route   1      0          0
  -> 10.10.9.70:30000             Route   1      0          0
```
