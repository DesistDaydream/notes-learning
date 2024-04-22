---
title: "iptables"
linkTitle: "iptables"
date: "2022-07-21T12:40"
weight: 1
---

# 概述

> 参考：
>
> - [Manual(手册)，iptables(8)](https://man7.org/linux/man-pages/man8/iptables.8.html)
> - [Netfilter 官方文档，iptables 教程](https://www.frozentux.net/iptables-tutorial/iptables-tutorial.html)

iptables 是 Netfilter 团队开发的一组用于与 netfilter 模块进行交互的 CLI 工具，其中包括 iptables、ip6tables、arptables、ebtables 等。

iptables 和 ip6tables 用于建立、维护和检查 Linux 内核中的 IPv4 和 IPv6 包过滤规则表。可以定义几个不同的表中的各种规则，也可以定义用户定义的链。并把已经定义的规则发送给 netfilter 模块。

## 四表(Table)

**注意：四表是 iptables 框架中的概念，不是 Netfilter 中的**

iptables 框架将流量抽象分为 4 类：过滤类、网络地址转换类、拆解报文类、原始类。每种类型的链作用在 Netfilter 系统中的 Hook 各不不相同，每种类型具有不同的功能。每一类都称为一张表。比如 fileter 表用来在指定链上检查流量是否可以通过，nat 表用来在指定链上检查流量是否可以进行地址转换，等等。Note：不是所有表都可以在所有链上具有规则，下表是 4 个表在 5 个 Hook 上的可用关系。

| 表名\链名 | PREROUTING | INPUT | FORWARD | OUTPUT | POSTROUTING |
| --------- | ---------- | ----- | ------- | ------ | ----------- |
| filter    |            | 可用  | 可用    | 可用   |             |
| nat       | 可用       | 可用  |         | 可用   | 可用        |
| mangle    | 可用       | 可用  | 可用    | 可用   | 可用        |
| raw       | 可用       |       |         | 可用   |             |

iptables 中有默认的内置 4 个表，每个表的名称就是其 chain 类型的名称

### filter - 过滤器

- 该类型的链可作用在以下几个 Hook 点上：INPUT、FORWARD、OUTPUT

### nat - 网络地址转换

- 该类型的链可作用在以下几个 Hook 点上：PREROUTING(DNAT)、INPUT、OUTPUT、POSTROUTING(SNAT)
  - 其中 PREROUTING 与 POSTROUTING 是流量经过物理网卡设备时做 nat 的地方
  - 其中 INPUT 与 OUTPUT 则是主机内部从网络栈直接下来的流量做 nat 的地方。e.g.从主机一个服务发送数据到同一个主机另一个服务的流量

### mangle - 拆解报文，做出修改，封装报文

- 该类型的链可作用在以下几个 Hook 点上：PREROUTING、INPUT、FORWARD、OUTPUT、POSTROUTING

### raw- 原始

用于跳过 nat 表以及连接追踪机制(ip_conntrack)的处理，详见 [连接跟踪系统](/docs/1.操作系统/Kernel/Network/Linux%20网络流量控制/Connnection%20Tracking(连接跟踪).md)

- 该类型的链可作用在以下几个 Hook 点上：PREROUTING、OUTPUT

整个表只有这一个作用，且仅有一个 target 就是 NOTRACK。具有最高优先级，所有流量先在两个 Hook 的 raw 功能上进行检查。一旦在 raw 上配置了规则，则 raw 表处理完成后，跳过 nat 表和 ip_conntrack 处理，i.e.不再做地址转换和数据包的链接跟踪处理，不把匹配到的数据包保存在“链接跟踪表”中。常用于那些不需要做 nat 的情况下以提高性能。e.g.大量访问的 web 服务器，可以让 80 端口不再让 iptables 做数据包的链接跟踪处理 ，以提高用户的访问速度。不过该功能不影响其余表的连接追踪追踪功能的正常使用，依然会有记录写到连接追踪文件中去

- 该功能的起源：iptables 表有大小的上限，如果每个数据包进来都要进行检查，会非常影响性能，可以对那些不用进行 nat 功能的数据进行放弃后面的检查，i.e.可以在 raw 配置然后直接让这些数据包跳过后面的表对该数据包的处理

Note：四表的优先级从高到低依次为：raw-mangle-nat-filter，i.e.数据到达某个 Hook 上，则会优先使用优先级最高类型的链来处理数据包。其实，iptables 表的作用更像是用来划分优先级的。

## iptables 处理链上规则的顺序以及规范

注意：每个数据包在 CHAIN 中匹配到适用于自己的规则之后，则直接进入下一个 CHAIN，而不会遍历 CHAIN 中每条规则去挨个匹配适用于自己的规则。比如下面两种情况

INPUT 链默认 DROP，匹配第一条：目的端口是 9090 的数据 DROP，然后不再检查下一项，那么 9090 无法访问

```bash
-P INPUT DROP
-A INPUT -p tcp -m tcp --dport 9090 -j DROP
-A INPUT -p tcp -m tcp --dport 9090 -j ACCEPT
```

INPUT 链默认 DROP，匹配第一条目的端口是 9090 的数据 ACCEPT，然后不再检查下一条规则，则 9090 可以访问

```bash
-P INPUT DROP
-A INPUT -p tcp -m tcp --dport 9090 -j ACCEPT
-A INPUT -p tcp -m tcp --dport 9090 -j DROP
```

# 安装 iptables

iptables 程序一般随系统安装自带（Minimal 也带），需要安装的通常是保证 iptables 规则可以在开机时启动的程序

## Ubuntu 安装 iptables

```bash
apt install netfilter-persistent iptables-persistent
```

netfilter-persistent 用来在保证在系统启动时加载 [Netfilter](/docs/1.操作系统/Kernel/Network/Linux%20网络流量控制/Netfilter/Netfilter.md) 规则；或者通过期内的一些 [Systemd](/docs/1.操作系统/Systemd/Systemd.md) 的 [Unit File](/docs/1.操作系统/Systemd/Unit%20File/Unit%20File.md) 和脚本自动加载保存好的 Netfilter 规则。iptables-persistent 算作 netfilter-persistent 包的插件，可以实现加载 iptables 规则效果

> [!Notes]
> Ubuntu 20.04 版本后，默认使用使用 nftables，安装 iptables-persistent 本质是 netfilter-persistent 包。iptables-persistent 作为 netfilter-persistent 的插件以兼容老的 iptables 功能。
>
> ```bash
> ~]# ll /lib/systemd/system/iptables.service
> lrwxrwxrwx 1 root root 34 Apr 20 11:36 /lib/systemd/system/iptables.service -> /etc/alternatives/iptables.service
> ~]# ll /etc/alternatives/iptables.service
> rwxrwxrwx 1 root root 48 Apr 20 11:36 /etc/alternatives/iptables.service -> /lib/systemd/system/netfilter-persistent.service
> ```

# iptables 关联文件与配置

**/run/xtables.lock** # 该文件在 iptables 程序启动时被使用，以获取排他锁

- 可以通过 `XTABLES_LOCKFILE` 环境变量修改 iptables 需要使用 xtalbes.lock 文件的路径

---

**RedHat 系特定的关联文件**

**/etc/sysconfig/iptables** # 存放用户定义的规则信息，每次重启 iptabels.service 服务后，都会读取该配置文件信息并应用到系统中

**/etc/sysconfig/iptables-conf** # 存放 iptables 工具的具体配置信息

---

**Debian 系特定的关联文件**

**/etc/iptables/rules.v4** # IPv4 版本的 iptables 规则保存文件，由 iptables-persistent.service 服务使用

**/etc/iptables/rules.v6** # IPv6 版本的 iptables 规则保存文件，由 iptables-persistent.service 服务使用
