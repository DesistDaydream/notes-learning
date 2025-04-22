---
title: Docker Network
linkTitle: Docker Network
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，引擎 - 网络](https://docs.docker.com/engine/network/)

# DNS

容器连接到默认网络时，会使用与主机相同的 DNS Server，继承主机的 /etc/resolv.conf。

容器连接到[自定义网络](https://docs.docker.com/engine/network/tutorials/standalone/#use-user-defined-bridge-networks)时，会使用 **<font color="#ff0000">Docker 内嵌的 DNS Server</font>**（127.0.0.11:53），将 容器 ID、容器名、别名 注册为 DNS 解析记录，以便其他容器可以通过这些域名访问到自己。示例可以参考下文 ”Bridge“ 驱动程序中的示例

# Docker 网络的驱动程序

https://docs.docker.com/engine/network/drivers/

Docker 使用驱动程序实现网络子系统的核心能力，网络子系统是可插拔的。默认存在如下几种驱动程序：

- **bridge** # 默认驱动程序。当应用程序在需要与同一主机上的其他容器通信的容器中运行时，通常会使用 [Bridge(桥接)](#bridge) 网络。
- **host** # 取消容器与 Docker 主机之间的网络隔离，直接使用主机的网络
- **overlay**
- **ipvlan**
- **macvlan**
- **none** # 容器与宿主机及其他容器完全隔离。 none 不适用于 Swarm 服务。用人话说就是不使用任何网络驱动程序。
- [网络插件](https://docs.docker.com/engine/extend/plugins_network/)：我们可以安装和使用第三方网络插件，以实现 Docker 自身无法实现的网络功能

## Bridge

https://docs.docker.com/engine/network/drivers/bridge

当 Docker 启动时，会自动创建一个名为 brdige，Bridge 驱动类型的网络。默认情况下新启动的容器都会连接到 bridge 网络。

```bash
~]# docker network ls
NETWORK ID     NAME         DRIVER    SCOPE
dadd048eefa0   bridge       bridge    local
84cab5ef9276   host         host      local
4718cdfcb116   monitoring   bridge    local
4d68f227ca5d   none         null      local
```

> Notes: 名为 bridge 的是默认 Bridge 类型网络；名为 monitoring 的是用户自定义的 Bridge 类型网络。

连接到 Docker 默认 Bridge 类型网络上的容器只能通过 IP 地址互相访问。但是如果容器连接到用户自定义的 Bridge 类型网络，这些容器将会记录在 [Docker 内嵌的 DNS](#内嵌的%20DNS)，后续可以通过名称或别名访问这些容器。

使用默认 Bridge 网络的效果：

```bash
~]# docker run -d --rm --name nginx nginx:1.27
55cd45af85d9e2be4dd505b3236120130155b66e5636960b9980ff39a0c0858c
~]# docker run --rm -it  nicolaka/netshoot nslookup nginx
Server:         172.16.11.222
Address:        172.16.11.222#53

** server can't find nginx: SERVFAIL
```

使用自定义 Bridge 网络的效果：

```bash
~]# docker network create --driver bridge test
fd603b1caf6cbbd55c57542851eeeaeb63c3fc622899206eb9909527a03fe7d4
~]# docker run -d --rm --name nginx --network test nginx:1.27
7d000f6870a41ac3a3c5cd44e4e4bb47c27f027748cea869659add90aecc1a54
~]# docker run --rm -it --network test  nicolaka/netshoot sh -c "nslookup nginx && nslookup 7d000f6870a4"
Server:         127.0.0.11
Address:        127.0.0.11#53

Non-authoritative answer:
Name:   nginx
Address: 10.38.9.2

Server:         127.0.0.11
Address:        127.0.0.11#53

Non-authoritative answer:
Name:   7d000f6870a4
Address: 10.38.9.2
```

本质上，使用 默认 或 自定义 Bridge 类型网络的区别在于是否会使用内嵌 DNS，从 /etc/resolv.conf 文件也可以看出来，DNS 地址也是不一样的。

# 数据包处理与防火墙

> 参考：
>
> - [官方文档，引擎 - 网络 - 包过滤与防火墙](https://docs.docker.com/engine/network/packet-filtering-firewalls)

在 [Unix-like OS](docs/1.操作系统/Operating%20system/Unix-like%20OS/Unix-like%20OS.md) 上，Docker 使用 [iptables](docs/1.操作系统/Kernel/Network/Linux%20网络流量控制/Netfilter/iptables/iptables.md) 来实现 网络隔离、端口发布、端口过滤。

Netfilter 规则保证 Docker 的 bridge 网络驱动正常运行，不应素以修改 Docker 创建的规则，如果想要通过 iptables 规则来控制外部网络访问本机的容器，有一些特殊的地方需要注意

## 端口发布与隔离

默认情况下，对于 IPv4 和 IPv6，守护进程都会阻止对未发布端口的访问（所有）。已发布的容器端口将映射到主机 IP 地址。为此，它使用 iptables 执行网络地址转换 (NAT)、端口地址转换 (PAT) 和伪装。

例如， `docker run -p 8080:80 [...]` 在 Docker 主机上任意地址的端口 8080 和容器的端口 80 之间创建映射。来自容器的传出连接将使用 Docker 主机的 IP 地址进行伪装。

<font color="#ff0000">所有已发布的端口</font>，默认情况下，<font color="#ff0000">允许所有外部的 IP 访问</font>。若想为这些暴露出来的端口添加白名单机制，可以通过类似下面这种方式添加 iptables 规则：

> Notes: 因为 nat 表中包含 `-A DOCKER ! -i br-51dda9852a61 -p tcp -m tcp --dport 10443 -j DNAT --to-destination 10.38.1.2:10443` 这种规则，这种规则可以让数据包从 PREROUTING 直接进入 filter 表的 FORWARD 链，从而忽略 INPUT 链的各种规则

```bash
export ext_if="ens7f0np0" # 与外部通信的网络设备名称
iptables -I DOCKER-USER -i ${ext_if} ! -s 100.64.0.0/24 -j DROP
```

> 上面这个规则可以让所有源地址不是 100.64.0.0/24 的请求全部丢弃，让它们无法访问 Docker 暴露出来的端口。

但是，<font color="#ff0000">请注意</font>：官方文档的规则会导致本身可以访问互联网的容器<font color="#ff0000">无法访问互联网</font>，因为官方文档给的规则会让所有除了 100.64.0.0/24 端的数据包都丢弃，那么访问出去后回来的包也是被丢弃的。从下面的抓包可以看出来，veth5f39ff0(容器网络设备), br-df92ba64bbdc(Docker 网桥) 这两个只有出没有回，但是宿主机上的 ens3 设备包含来/回的包。说明从 223.5.5.5 回来的数据包从 ens3 转发(FORWARD) 时被丢弃了。

```bash
~]# tcpdump -i any host 223.5.5.5 -nn
tcpdump: data link type LINUX_SLL2
tcpdump: verbose output suppressed, use -v[v]... for full protocol decode
listening on any, link-type LINUX_SLL2 (Linux cooked v2), snapshot length 262144 bytes
13:31:51.418398 veth5f39ff0 P   IP 10.38.1.2 > 223.5.5.5: ICMP echo request, id 7, seq 279, length 64
13:31:51.418398 br-df92ba64bbdc In  IP 10.38.1.2 > 223.5.5.5: ICMP echo request, id 7, seq 279, length 64
13:31:51.418440 ens3  Out IP 10.10.4.90 > 223.5.5.5: ICMP echo request, id 7, seq 279, length 64
13:31:51.436852 ens3  In  IP 223.5.5.5 > 10.10.4.90: ICMP echo reply, id 7, seq 279, length 64
13:31:52.442411 veth5f39ff0 P   IP 10.38.1.2 > 223.5.5.5: ICMP echo request, id 7, seq 280, length 64
13:31:52.442411 br-df92ba64bbdc In  IP 10.38.1.2 > 223.5.5.5: ICMP echo request, id 7, seq 280, length 64
13:31:52.442462 ens3  Out IP 10.10.4.90 > 223.5.5.5: ICMP echo request, id 7, seq 280, length 64
13:31:52.461303 ens3  In  IP 223.5.5.5 > 10.10.4.90: ICMP echo reply, id 7, seq 280, length 64
```

综上，假如我们想要限制 10443 和 19093 两个端口不可以被随便访问，则要使用下面这种规则

```bash
-A DOCKER-USER -s 100.64.0.0/24 -i ens7f0np0 -j ACCEPT
-A DOCKER-USER -i ${ext_if} -p tcp -m multiport --dports 10443,19093 -j DROP
```
