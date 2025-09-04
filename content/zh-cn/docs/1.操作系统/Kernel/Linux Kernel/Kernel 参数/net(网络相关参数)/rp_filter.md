---
title: rp_filter
linkTitle: rp_filter
weight: 20
---

# 概述

> 参考：
>
> - 

rp_filter 参数应用实例

rp_filter 参数示例

假设机器有 2 个网口:

- eth0: 192.168.1.100

- eth1：200.153.1.122

数据包源 IP：10.75.153.98，目的 IP：200.153.1.122

系统路由表配置为：

\[root@localhost ~]# route -n

Kernel IP routing table

Destination Gateway Genmask Flags Metric Ref Use Iface

default 192.168.1.234 0.0.0.0 　　 UG 0 0 0 eth0

192.168.120.0 0.0.0.0 255.255.255.0 U 0 0 0 eth0

10.75.153.98 0.0.0.0 255.255.255.0 U 0 0 0 eth0

系统 rp_filter 参数的配置为：

\[root@localhost ~]# sysctl -a | grep rp_filter

net.ipv4.conf.all.rp_filter=1

net.ipv4.conf.default.rp_filter=1

如上所示，数据包发到了 eth1 网卡，如果这时候开启了 rp_filter 参数，并配置为 1，则系统会严格校验数据包的反向路径。从路由表中可以看出，返回响应时数据包要从 eth0 网卡出，即请求数据包进的网卡和响应数据包出的网卡不是同一个网卡，这时候系统会判断该反向路径不是最佳路径，而直接丢弃该请求数据包。（业务进程也收不到该请求数据包）

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zhv2lk/1616168317314-27937a08-a8c2-41a5-b15c-2cd36722d4ca.png)

解决办法：

1.修改路由表，使响应数据包从 eth1 出，即保证请求数据包进的网卡和响应数据包出的网卡为同一个网卡。

2.关闭 rp_filter 参数。（注意 all 和 default 的参数都要改）

1\)修改/etc/sysctl.conf 文件，然后 sysctl -p 刷新到内存。

2\)使用 sysctl -w 直接写入内存：sysctl -w net.ipv4.conf.all.rp_filter=0

3\)修改/proc 文件系统： echo "0">/proc/sys/net/ipv4/conf/all/rp_filter

## 开启 rp_filter 参数的作用

1. 减少 DDoS 攻击

校验数据包的反向路径，如果反向路径不合适，则直接丢弃数据包，避免过多的无效连接消耗系统资源。

2. 防止 IP Spoofing

校验数据包的反向路径，如果客户端伪造的源 IP 地址对应的反向路径不在路由表中，或者反向路径不是最佳路径，则直接丢弃数据包，不会向伪造 IP 的客户端回复响应。

### Ps：两种常见的非法攻击手段：

1. DDos 攻击(Distribute Deny of Service)

分布式拒绝服务攻击。通过构造大量的无用数据包向目标服务发起请求，占用目标服务主机大量的资源，还可能造成网络拥塞，进而影响到正常用户的访问。

2. IP Spoofing（IP 欺骗）

IP Spoofing 指一个客户端通过伪造源 IP，冒充另外一个客户端与目标服务进行通信，从而达到某些不可告人的秘密。

## 另一种说法：

比如一台设备安装两~~个 ~~haproxy，有多个网~~卡~~，多个 ip，不同 haproxy 代理不同网段，给不同后端。那么就需要关闭校验。

## 另一种说法：

在 Linux 中用于对 网卡的反向路由过滤策略进行配置的内核参数是 rp_filter，有关此参数的详细介绍以及配置方式请参见 Linux 内核参数 rp_filter。

LVS 在 VS/TUN 模式下，需要对 tunl0 虚拟网卡的反向路由过滤策略进行配置。最直接的办法是把其值设置为 0。

```bash
net.ipv4.conf.tunl0.rp_filter=0
net.ipv4.conf.all.rp_filter=0
```

因为 Linux 系统在对网卡应用反向路由过滤策略时，除了检查本网卡的 rp_filter 参数外，还会检查 all 配置项上的 rp_filter 参数，并使用这两个值中较大的值作为应用到当前网卡的反向路由过滤策略。所以需要同时把 `net.ipv4.conf.all.rp_filter` 参数设置为 0。
