---
title: Docker Network
linkTitle: Docker Network
date: 2024-07-05T08:39
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

- **bridge** # 默认驱动程序。当应用程序在需要与同一主机上的其他容器通信的容器中运行时，通常会使用 [Bridge(桥接)](#Bridge) 网络。
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