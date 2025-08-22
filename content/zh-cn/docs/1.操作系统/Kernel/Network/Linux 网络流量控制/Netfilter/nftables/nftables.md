---
title: Nftables
linkTitle: Nftables
weight: 1
---

# 概述

> 参考：
>
> - [官方介绍](https://www.netfilter.org/projects/nftables/index.html)
> - [官方 wiki](https://wiki.nftables.org/wiki-nftables/index.php/Main_Page)
> - https://mp.weixin.qq.com/s/VPBgIqlTiFe7KwqF1Q4pCA

nftables 是一个 [Netfilter](/docs/1.操作系统/Kernel/Network/Linux%20网络流量控制/Netfilter/Netfilter.md) 项目，旨在替换现有的 {ip,ip6,arp,eb}tables 框架，为 {ip,ip6}tables 提供一个新的包过滤框架、一个新的用户空间实用程序（nft）和一个兼容层。它使用现有的钩子、链接跟踪系统、用户空间排队组件和 netfilter 日志子系统。

nftables 主要由三个组件组成：内核实现、libnl netlink 通信、 nftables 用户空间。 其中内核提供了一个 [Netlink](/docs/1.操作系统/Kernel/Network/Linux%20网络栈管理/Netlink/Netlink.md) 配置接口以及运行时规则集评估，libnl 包含了与内核通信的基本函数，用户空间可以通过 nft 和用户进行交互。

nftables 与 iptables 的区别

nftables 和 iptables 一样，由 table(表)、chain(链)、rule(规则) 组成。nftables 中，表包含链，链包含规则，规则是真正的 action。与 iptables 相比，nftables 主要有以下几个变化：

- iptables 规则的布局是基于连续的大块内存的，即数组式布局；而 nftables 的规则采用链式布局。其实就是数组和链表的区别
- iptables 大部分工作在内核态完成，如果要添加新功能，只能重新编译内核；而 nftables 的大部分工作是在用户态完成的，添加新功能很 easy，不需要改内核。
- iptables 有内置的链，即使你只需要一条链，其他的链也会跟着注册；而 nftables 不存在内置的链，你可以按需注册。由于 iptables 内置了一个数据包计数器，所以即使这些内置的链是空的，也会带来性能损耗。
- 简化了 IPv4/IPv6 双栈管理
- 原生支持集合、字典和映射

> [!tip] nftables 没有任何默认规则（也没有任何内置表）。如果关闭了 firewalld 服务，则命令 `nft list ruleset` 输出结果为空。

## nftables table 与 nftables family

> 参考：
>
> - https://wiki.nftables.org/wiki-nftables/index.php/Nftables_families
> - https://www.mankier.com/8/nft#Address_Families
> - https://www.mankier.com/8/nft#Tables

Nftables Tables 是承载 Chain, rule, etc. 的容器。Tables 必须由 “族” 和 “表名” 作为标识。

> [!Tip] 与 Iptables 的表改变不同。Nftables 没有内置表，表的数量与名称可以自行设置，表本身并没有类型。只有表中的 Chain 需要设置类型

**family(族)** 是 nftables 技术引用的新概念。一共有 6 种族。不同的 family 可以处理不同 Hook 上的数据包。

> [!Note]
>
> - `族` 是具有相同属性的一类网络层级或者说网络类型，比如建立一个名为 test 的表，该表的族为 inet(i.e.表的类型是 inet)。
> - 在 [Iptables](/docs/1.操作系统/Kernel/Network/Linux%20网络流量控制/Netfilter/Iptables/Iptables.md) 中，每个网络层级是由单独的工具实现的：e.g. iptables, ip6tables, arptables, ebtables 。而 nftables 想要通过单个命令行工具提供这些层级网络的控制，则需要抽象出一个新的分类概念，i.e. family
> - 所以每个表应且只应指定一个族，且当表中的链被指定类型时，只能指定该族下可以处理的链类型，详情见本文《nftables chain 链》章节

nftables 中一同以下几种 family：

- **ip** # IPv4 地址族。对应 iptables 中 iptables 命令行工具所实现的效果。默认族，nft 命令的所有操作如果不指定具体的 family，则默认对 ip 族进行操作
- **ip6** # IPv6 地址族。对应 iptables 中 ip6tables 命令行工具所实现的效果
- **inet** # Internet (IPv4/IPv6)地址族。对应 iptables 中 iptables 和 ip6tables 命令行工具所实现的效果
- **arp** # ARP 地址族，处理 IPv4 ARP 包。对应 iptables 中 arptables 命令行工具所实现的效果
- **bridge** # 桥地址族。处理通过桥设备的数据包对应 iptables 中 ebtables 命令行工具所实现的效果
- **netdev** # Netdev address family, handling packets from ingress.

各 family 可处理的 Hook 矩阵图如下：

| family\hook | prerouting | input | forward | output | postrouting | ingress | egress |
| ----------- | :--------: | :---: | :-----: | :----: | :---------: | :-----: | :----: |
| **ip**      |     √      |   √   |    √    |   √    |      √      |    ×    |   ×    |
| **ip6**     |     √      |   √   |    √    |   √    |      √      |    ×    |   ×    |
| **inet**    |     √      |   √   |    √    |   √    |      √      |    √    |   ×    |
| **arp**     |     ×      |   √   |    ×    |   √    |      ×      |    ×    |   ×    |
| **bridge**  |     √      |   √   |    √    |   √    |      √      |    ×    |   ×    |
| **netdev**  |     ×      |   ×   |    ×    |   ×    |      ×      |    √    |   √    |

## nftables chain

nftables 中的 chain(链) 用于保存规则，分为两大 kind(种类)：

- **base chains(基本链)** # 数据包的入口点，需要指定该链的基本信息（type(类型)、作用的 Hook 点、priority(优先级)、policy(策略)、etc.）才可以让链中的规则生效。因为链中包含一条一条的规则，所以一个可以正常处理流量的链，需要指定其类型来区分该链上的规则干什么用的，还需要指定 Hook 来指明数据包到哪个 Hook 了来使用这个规则，还需要配置优先级来处理相同类型的规则，该规则应该先执行还是后执行。
- **regular chains(常规链)**(也叫自定义链) # 不需要指定钩子类型和优先级，可以用来做链与链之间的跳转，从逻辑上对规则进行分类。

> 具有类型的叫基本链，没有类型的叫常规链。

chain 分为下述三种 type(类型)：

- **filter 类型** # 用于过滤数据包所用
  - 允许定义在哪些 family 下：all family
  - 链中的规则会处理这些 Hook 点的数据包：all Hook
- **nat 类型** # 用于进行地址转换
  - 允许定义在哪些 family 下：ip、ip6
  - 链中的规则会处理这些 Hook 点的数据包：prerouting、input、output、postrouting
- **route 类型**
  - 允许定义在哪些 family 下：ip、ip6
  - 链中的规则会处理这些 Hook 点的数据包：output

每个基本链都需要指定该链 type、hook、priority、policy，比如下面：

```bash
table ip filter {
        chain FORWARD {
                type filter hook forward priority filter; policy accept;
        }

        chain INPUT {
                type filter hook input priority filter; policy accept;
        }
}
table ip nat {
        chain POSTROUTING {
                type nat hook postrouting priority srcnat; policy accept;
        }
}
```

## nftables rule

nftables 中的规则标识符有两种，一种 index，一种 handle

**index # 规则的索引。每条规则在其链中，从 0 开始计数(每条链中的规则，第一条规则的 index 为 0，第二条规则的 indext 为 2，依次类推)。**

```bash
chain DOCKER {
    tcp dport tcpmux accept # 规则 index 为0
    tcp dport 5 accept # 规则 index 为 1
    tcp dport 6 accept # 后续依次类推
    tcp dport 2 accept
    tcp dport 3 accept
    tcp dport afs3-fileserver accept
}
```

**handle # 规则的句柄。句柄对于整个 nftalbes 而言，不管添加在哪个链中，第一条规则的句柄为 1，第二条规则句柄为 2。如果规则句柄为 33 号被删除，则新添加的规则的句柄为 34**

```bash
chain DOCKER { # handle 4
    tcp dport tcpmux accept # handle 28
    tcp dport 5 accept # handle 32
    tcp dport 6 accept # handle 33
    tcp dport 2 accept # handle 29
    tcp dport 3 accept # handle 30
    tcp dport afs3-fileserver accept # handle 31
}
```

Note：对于每条规则而言，其 index 可以随时改变，当在多个规则中间插入新规则时，新插入规则下面的规则 index 则会改变。而 handle 则不会改变，除非删除后重新添加

## Object

Object 我自己根据官方文档总结的：

Nftables 创建的 tables, chains, rules, sets, etc. 都可以称为 Object(对象)，每个对象都有一个 handle(句柄) 作为唯一标识。

> 这个总结的灵感来源是 [-a, --handle](https://www.mankier.com/8/nft#--handle) 选项的解释：Show object handles in output. 当使用 `nft -a list ruleset` 时，创建的每个条目（tables, chains, rules, etc.）都有一个 handle 号，所以我将 Nftables 创建的每个条目都抽象为 Object

其中有些 Object 是 [Stateful Objects](https://www.mankier.com/8/nft#Stateful_Objects)(有状态对象)

## 总结

nftables 的结构为：表包含链，链包含规则，这个逻辑是非常清晰明了的。而 iptable 呢，则需要先指定什么类型的表，再添加规则，规则与链则互相存在，让人摸不清关系；其实也可以说，iptables 的表类型，就是 nftables 中的链的类型。

Nftables 除了基本的 table, chain, rule, sets，还可以使用很多对象来控制网络流量

- Flowtables(流表)
- Expressions(表达式)
- etc.

# 安装 Nftables

nftables 程序与 [Iptables](/docs/1.操作系统/Kernel/Network/Linux%20网络流量控制/Netfilter/Iptables/Iptables.md) 程序一样，一般随系统安装自带（Minimal 也带），需要安装的通常是保证 nftables 规则可以在开机时启动的程序（只不过 nftables 是新的程序，各类系统默认安装的是 iptables 还是 nftables，取决于自身的规划）。

各 [Unix-like OS](/docs/1.操作系统/Operating%20system/Unix-like%20OS/Unix-like%20OS.md) 默认使用 nftables 的版本说明

- Ubuntu # 从 20 及以后的版本，默认使用 nftables，所有 iptables 相关的文件，都是 nftables 相关文件的 [Symbolic link](/docs/1.操作系统/Kernel/Filesystem/文件管理/Symbolic%20link.md)(符号链接)
  - 主要通过 netfilter-persistent 包实现规则的持久化
- RedHat # TODO

# Nftables 关联文件与配置

---

**RedHat 系特定的关联文件**

**/etc/sysconfig/nftables.conf** # CentOS 8 中，nftables.service 的规则被存储在此目录中，其中 include 一些其他的示例规则

**/etc/sysconfig/nftables/** # nftables.conf 文件中 include 的文件，都在该目录下

---

**Debian 系特定的关联文件**

TODO

# Nftable 特性

## Set

> 参考：
>
> - [官方文档，Sets](https://wiki.nftables.org/wiki-nftables/index.php/Sets)

nftables 的语法原生支持 set(集合)，集合可以用来匹配多个 IP 地址、端口号、网卡或其他任何条件。类似于 ipset 的功能。

集合分为匿名集合与命名集合。

### 匿名集合

匿名集合比较适合用于未来不需要更改的规则

例如下面的两个示例，

- 该规则允许来自源 IP 处于 10.10.10.123 ~ 10.10.10.231 这个区间内的主机的流量通过。
  - `nft add rule inet my_table my_filter_chain ip saddr { 10.10.10.123, 10.10.10.231 } accept`
- 该规则允许来自目的端口是 http、nfs、ssh 的流量通过。
  - `nft add rule inet my_table my_filter_chain tcp dport { http, nfs, ssh } accept`

匿名集合的缺点是，如果需要修改集合中的内容，比如像 ipset 中修改 ip 似的，就得替换规则。如果后面需要频繁修改集合，推荐使用命名集合。

### 命令集合

iptables 可以借助 ipset 来使用集合，而 nftables 中的命名集合就相当于 ipset 的功能

命名集合需要使用 nft add set XXXX 命令进行创建，创建时需要指定族名、表名、以及 set 的属性

命名集合中包括以下几种属性（其中 type 为必须指定的属性，其余属性可选）：

- **type** # 集合中所有元素的类型，包括 ipv4_addr(ipv4 地址), ipv6_addr(ipv6 地址), ether_addr(以太网地址), inet_proto(网络协议), inet_service(网络服务), mark(标记类型) 这几类
- **flags** # 集合的标志。包括 constant、interval、timeout 。
  - interval # 让集合支持区间模式。默认集合中无法使用这种方式 `nft add element inet my_table my_set { 10.20.20.0-10.20.20.255 }` 来添加集合 。当给集合添加该 flag 时，就可以在给集合添加元素时，使用‘区间’的表示方法。因为内核必须提前确认该集合存储的数据类型，以便采用适当的数据结构。
- **timeout** #
- **gc-interval** #
- **elements** #
- **size** #
- **policy** #
- **auto-merge** #

像 ipset 一样，光创建完还没法使用，需要在 iptables 中添加规则引用 ipset 才可以。nftables 的 set 一样，创建完成后，需要在规则中引用，引用集合规则时使用 @ 并跟上集合的名字，即可引用指定的集合(e.g.`nft insert rule inet my_table my_filter_chain ip saddr @my_set drop`)这条命令即时引用了 my_set 集合中的内容

级联不同类型

命名集合也支持对不同类型的元素进行级联，通过级联操作符 . 来分隔。例如，下面的规则可以一次性匹配 IP 地址、协议和端口号。

```bash
$ nft add set inet my_table my_concat_set { type ipv4_addr . inet_proto . inet_service ; }

$ nft list set inet my_table my_concat_set

table inet my_table {
    set my_concat_set {
        type ipv4_addr . inet_proto . inet_service
    }
}
```

向集合中添加元素：

`$ nft add element inet my_table my_concat_set { 10.30.30.30 . tcp . telnet }`

在规则中引用级联类型的集合和之前一样，但需要标明集合中每个元素对应到规则中的哪个位置。

`$ nft add rule inet my_table my_filter_chain ip saddr . meta l4proto . tcp dport @my_concat_set accept`

这就表示如果数据包的源 IP、协议类型、目标端口匹配 10.30.30.30、tcp、telnet 时，nftables 就会允许该数据包通过。

匿名集合也可以使用级联元素，例如：

`$ nft add rule inet my_table my_filter_chain ip saddr . meta l4proto . udp dport { 10.30.30.30 . udp . bootps } accept`

nftables 级联类型的集合类似于 ipset 的聚合类型，例如 hash:ip,port。

# Nftable 规则文件

> 参考：
>
> - [Manual, nft - Input File Formats](https://www.mankier.com/8/nft#Input_File_Formats)

使用 define 定义变量，使用 `$` 引用变量，比如如下官方[示例](https://wiki.nftables.org/wiki-nftables/index.php/Sets#nftables.conf_syntax)：

```bash
define SIMPLE_SET = { 192.168.1.1, 192.168.1.2 }

define CDN_EDGE = {
    192.168.1.1,
    192.168.1.2,
    192.168.1.3,
    10.0.0.0/8
}

define CDN_MONITORS = {
    192.168.1.10,
    192.168.1.20
}

define CDN = {
    $CDN_EDGE,
    $CDN_MONITORS
}

# Allow HTTP(S) from approved IP ranges only
tcp dport { http, https } ip saddr $CDN accept
udp dport { http, https } ip saddr $CDN accept
```
