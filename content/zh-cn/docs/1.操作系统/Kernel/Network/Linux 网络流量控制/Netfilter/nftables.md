---
title: nftables
linkTitle: nftables
date: 2024-04-20T10:06
weight: 4
---

# 概述

> 参考：
>
> - [官方 wiki](https://wiki.nftables.org/wiki-nftables/index.php/Main_Page)

nftables 是一个 [Netfilter](/docs/1.操作系统/Kernel/Network/Linux%20网络流量控制/Netfilter/Netfilter.md) 项目，旨在替换现有的 {ip,ip6,arp,eb}tables 框架，为 {ip,ip6}tables 提供一个新的包过滤框架、一个新的用户空间实用程序（nft）和一个兼容层。它使用现有的钩子、链接跟踪系统、用户空间排队组件和 netfilter 日志子系统。

nftables 主要由三个组件组成：内核实现、libnl netlink 通信、 nftables 用户空间。 其中内核提供了一个 netlink 配置接口以及运行时规则集评估，libnl 包含了与内核通信的基本函数，用户空间可以通过 nft 和用户进行交互。

nftables 与 iptables 的区别

nftables 和 iptables 一样，由 table(表)、chain(链)、rule(规则) 组成。nftables 中，表包含链，链包含规则，规则是真正的 action。与 iptables 相比，nftables 主要有以下几个变化：

- iptables 规则的布局是基于连续的大块内存的，即数组式布局；而 nftables 的规则采用链式布局。其实就是数组和链表的区别
- iptables 大部分工作在内核态完成，如果要添加新功能，只能重新编译内核；而 nftables 的大部分工作是在用户态完成的，添加新功能很 easy，不需要改内核。
- iptables 有内置的链，即使你只需要一条链，其他的链也会跟着注册；而 nftables 不存在内置的链，你可以按需注册。由于 iptables 内置了一个数据包计数器，所以即使这些内置的链是空的，也会带来性能损耗。
- 简化了 IPv4/IPv6 双栈管理
- 原生支持集合、字典和映射

nftables 没有任何默认规则，如果关闭了 firewalld 服务，则命令 nft list ruleset 输出结果为空。意思就是没有任何内置链或者表

## nftables table 表 与 nftables family 簇

nftables 没有内置表，表的数量与名称由用户决定。

**family(簇)** 是 nftables 技术引用的新概念。一共有 6 种簇。不同的 family 可以处理不同 Hook 上的数据包。

Note：

- `簇` 可以当做 `类型` 来理解，比如建立一个名为 test 的表，该表的簇为 inet(i.e.表的类型是 inet)。
- 所以每个表应且只应指定一个簇，且当表中的链被指定类型时，只能指定该簇下可以处理的链类型，详情见本文《nftables chain 链》章节

nftables 中一同以下几种 family：

- ip # IPv4 地址簇。对应 iptables 中 iptables 命令行工具所实现的效果。默认簇，nft 命令的所有操作如果不指定具体的 family，则默认对 ip 簇进行操作
  - 可处理流量的 Hook：与 inet 簇相同
- ip6 # IPv6 地址簇。对应 iptables 中 ip6tables 命令行工具所实现的效果
  - 可处理流量的 Hook：与 inet 簇相同
- inet # Internet (IPv4/IPv6)地址簇。对应 iptables 中 iptables 和 ip6tables 命令行工具所实现的效果
  - 可处理流量的的 Hook：prerouting、input、forward、output、postrouting。ip 与 ip6 簇与 inet 簇所包含的 Hook 相同
- arp # ARP 地址簇，处理 IPv4 ARP 包。对应 iptables 中 arptables 命令行工具所实现的效果
  - 可处理流量的 Hook：input、output。
- bridge # 桥地址簇。处理通过桥设备的数据包对应 iptables 中 ebtables 命令行工具所实现的效果
  - 可处理流量的 Hook：与 inet 簇相同
- netdev # Netdev address family, handling packets from ingress.
  - 可处理流量的 Hook：ingress

基本效果示例如下：

```bash
~]# nft add table test # 创建名为test的表，簇为默认的ip簇
~]# nft list ruleset # 列出所有规则
table ip test { # 仅有一个名为test的表，簇为ip，没有任何规则
}
~]# nft add table inet test # 创建名为test的表，使用inet簇
~]# nft list ruleset
table ip test {
}
table inet test {
}
```

## nftables chain(链)

在 nftables 中，链是用来保存规则的。链在逻辑上被分为下述三种类型：

1. filter 类型的链 # 用于过滤数据包所用
   1. 允许定义在哪些 family 下：all family
   2. 链中的规则会处理这些 Hook 点的数据包：all Hook
2. nat 类型的链 # 用于进行地址转换
   1. 允许定义在哪些 family 下：ip、ip6
   2. 链中的规则会处理这些 Hook 点的数据包：prerouting、input、output、postrouting
3. route 类型的链
   1. 允许定义在哪些 family 下：ip、ip6
   2. 链中的规则会处理这些 Hook 点的数据包：output

在创建 nftables 中的链时，通常有两种叫法，没有类型的叫常规链，含有类型的叫基本链：

常规链(也叫自定义链) : 不需要指定钩子类型和优先级，可以用来做链与链之间的跳转，从逻辑上对规则进行分类。

基本链 : 数据包的入口点，需要指定该链的基本信息(类型、作用的 Hook 点、优先级、默认策略等)才可以让链中的规则生效(在链管理命令的 {} 中添加链信息)。因为链中包含一条一条的规则，所以一个可以正常处理流量的链，需要指定其类型来区分该链上的规则干什么用的，还需要指定 Hook 来指明数据包到哪个 Hook 了来使用这个规则，还需要配置优先级来处理相同类型的规则，该规则应该先执行还是后执行。

## nftables rule(规则)

nftables 中的规则标识符有两种，一种 index，一种 handle

**index # 规则的索引。每条规则在其链中，从 0 开始计数(每条链中的规则，第一条规则的 index 为 0，第二条规则的 indext 为 2，依次类推)。**

```bash
 chain DOCKER {
  tcp dport tcpmux accept # 规则index为0
  tcp dport 5 accept # 规则index为1
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

## 总结

nftables 的结构为：表包含链，链包含规则，这个逻辑是非常清晰明了的。而 iptable 呢，则需要先指定什么类型的表，再添加规则，规则与链则互相存在，让人摸不清关系；其实也可以说，iptables 的表类型，就是 nftables 中的链的类型。

# 安装 Nftables

nftables 程序与 [iptables](/docs/1.操作系统/Kernel/Network/Linux%20网络流量控制/Netfilter/iptables/iptables.md) 程序一样，一般随系统安装自带（Minimal 也带），需要安装的通常是保证 nftables 规则可以在开机时启动的程序（只不过 nftables 是新的程序，各类系统默认安装的是 iptables 还是 nftables，取决于自身的规划）。

各 [Unix-like OS](/docs/1.操作系统/Operating%20system/Unix-like%20OS/Unix-like%20OS.md) 默认使用 nftables 的版本说明

- Ubuntu # 从 20 及以后的版本，默认使用 nftables，所有 iptables 相关的文件，都是 nftables 相关文件的 [Symbolic link](/docs/1.操作系统/Kernel/Filesystem/文件管理/Symbolic%20link.md)(符号链接)
  - 主要通过 netfilter-persistent 包实现规则的持久化
- RedHat # TODO

# Nftables 关联文件与配置

---

**RedHat 系特定的关联文件**

**/etc/sysconfig/nftables.conf** # CentOS 8 中，nftables.service 的规则被存储在此目录中，其中 include 一些其他的示例规则

**/etc/sysconfig/nftables/** # nftables.conf 文件中 include 的文件，都在该目录下

备份规则：

$ nft list ruleset > /root/nftables.conf

---

**Debian 系特定的关联文件**

TODO

# nftable 的 set(集合)与 map(字典) 特性介绍

nftables 的语法原生支持集合，集合可以用来匹配多个 IP 地址、端口号、网卡或其他任何条件。类似于 ipset 的功能。

集合分为匿名集合与命名集合。

## 匿名集合

匿名集合比较适合用于未来不需要更改的规则

例如下面的两个示例，

- 该规则允许来自源 IP 处于 10.10.10.123 ~ 10.10.10.231 这个区间内的主机的流量通过。
  - nft add rule inet my_table my_filter_chain ip saddr { 10.10.10.123, 10.10.10.231 } accept
- 该规则允许来自目的端口是 http、nfs、ssh 的流量通过。
  - nft add rule inet my_table my_filter_chain tcp dport { http, nfs, ssh } accept

匿名集合的缺点是，如果需要修改集合中的内容，比如像 ipset 中修改 ip 似的，就得替换规则。如果后面需要频繁修改集合，推荐使用命名集合。

## 命令集合

iptables 可以借助 ipset 来使用集合，而 nftables 中的命名集合就相当于 ipset 的功能。

命名集合需要使用 nft add set XXXX 命令进行创建，创建时需要指定簇名、表名、以及 set 的属性

命名集合中包括以下几种属性，其中 type 为必须指定的属性，其余属性可选。

- type # 集合中所有元素的类型，包括 ipv4_addr(ipv4 地址), ipv6_addr(ipv6 地址), ether_addr(以太网地址), inet_proto(网络协议), inet_service(网络服务), mark(标记类型) 这几类
- flags # 集合的标志。包括 constant、interval、timeout 。
  - interval # 让集合支持区间模式。默认集合中无法使用这种方式 nft add element inet my_table my_set { 10.20.20.0-10.20.20.255 } 来添加集合 。当给集合添加类型 flag 时，就可以在给集合添加元素时，使用‘区间’的表示方法。因为内核必须提前确认该集合存储的数据类型，以便采用适当的数据结构。
- timeout #
- gc-interval #
- elements #
- size #
- policy #
- auto-merge #

像 ipset 一样，光创建完还没法使用，需要在 iptables 中添加规则引用 ipset 才可以。nftables 的 set 一样，创建完成后，需要在规则中引用，引用集合规则时使用 @ 并跟上集合的名字，即可引用指定的集合(e.g.nft insert rule inet my_table my_filter_chain ip saddr @my_set drop)这条命令即时引用了 my_set 集合中的内容

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

$ nft add element inet my_table my_concat_set { 10.30.30.30 . tcp . telnet }

在规则中引用级联类型的集合和之前一样，但需要标明集合中每个元素对应到规则中的哪个位置。

$ nft add rule inet my_table my_filter_chain ip saddr . meta l4proto . tcp dport @my_concat_set accept

这就表示如果数据包的源 IP、协议类型、目标端口匹配 10.30.30.30、tcp、telnet 时，nftables 就会允许该数据包通过。

匿名集合也可以使用级联元素，例如：

$ nft add rule inet my_table my_filter_chain ip saddr . meta l4proto . udp dport { 10.30.30.30 . udp . bootps } accept

现在你应该能体会到 nftables 集合的强大之处了吧。

nftables 级联类型的集合类似于 ipset 的聚合类型，例如 hash:ip,port。

# nft 命令行工具介绍

**nft \[OPTIONS] \[COMMANDS]**

COMMANDS 包括：

- ruleset # 规则集管理命令
- table # 表管理命令
- chain # 链管理命令
- rule # 规则管理命令
- set # 集合管理命令
- map # 字典管理命令
- NOTE：
  - 该 COMMANDS 与后面子命令中的 COMMAND 不同，前者是 nft 命令下的子命令，后者是 nft 命令下子命令的子命令
  - nft 子命令默认对 ip 簇进行操作，当指定具体的 FAMILY 时，则对指定的簇进行操作

OPTIONS

- -a,--handle # 在使用命令获得输出时，显示每个对象的句柄
  - Note：handle(句柄)在 nftables 中，相当于标识符，nftables 中的每一行内容都有一个 handle。
- -e,--echo # 回显已添加、插入或替换的内容
- -f,--file FILE # 从指定的文件 FILE 中读取 netfilter 配置加载到内核中

EXAMPLE：

- nft -f /root/nftables.conf # 从 nftables.conf 文件中，将配置规则加载到系统中

Note：下面子命令中的 FAMILY 如果不指定，则所有命令默认都是对 ip 簇进行操作。

## table - 表管理命令

nft COMMAND table \[FAMILY] TABLE # FAMILY 指定簇名，TABLE 为表的名称

nft list tables # 列出所有的表，不包含表中的链和规则

COMMAND

- add # 添加指定簇下的表。
- create # 与 add 命令类似，但是如果表已经存在，则返回错误信息。
- delete # 删除指定的表。不管表中是否有内容都一并删除
- flush # 清空指定的表下的所有规则，保留链
- list # 列出指定的表的所有链，及其链中的规则

EXAMPLE

- nft add table my_table # 创建一个 ip 簇的，名为 my_table 的表
- nft add table inet my_table # 创建一个 inet 簇的，名为 my_table 的表
- nft list table inet my_table # 列出 inet 簇的名为 my_table 的表及其链和规则

## chains - 链管理命令

nft COMMAND chain \[FAMILY] TABLE CHAIN \[{ type TYPE hook HOOK \[device DEVICE] priority PRIORITY; \[policy POLICY;] }] # FAMILY 指定簇名，TABLE 指定表名，CHAIN 指定链名，TYPE 指定该链的类型，HOOK 指定该链作用在哪个 hook 上，DEVICE 指定该链作用在哪个网络设备上，PRIORITY 指定该链的优先级，POLICY 指定该链的策略(i.e.该链的默认策略，accept、drop 等等。)

nft list chains # 列出所有的链

Note:

- 在输入命令时，使用反斜线 \ 用来转义分号 ; ，这样 shell 就不会将分号解释为命令的结尾。如果是直接编辑 nftables 的配置文件则不用进行转义
- PRIORITY 采用整数值，可以是负数，值较小的链优先处理。

COMMAND

- add # 在指定的表中添加一条链
- create # 与 add 命令类似，但是如果链已经存在，则返回错误信息。
- delete # 删除指定的链。该链不能包含任何规则，或者被其它规则作为跳转目标，否则删除失败。
- flush #
- list # 列出指定表下指定的链，及其链中的规则
- rename #

EXAMPLE

- nft add chain inet my_table my_utility_chain # 在 inet 簇的 my_table 表上创建一个名为 my_utility_chain 的常规链，没有任何参数
- nft add chain inet my_table my_filter_chain{type filter hook input priority 0;} # 在 inet 簇的 my_table 表上创建一个名为 my_filter_chain 的链，链的类型为 filter，作用在 input 这个 hook 上，优先级为 0
- nft list chain inet my_table my_filter_chain # 列出 inet 簇的 my_table 表下的 my_filter_chain 链的信息，包括其所属的表和其包含的规则

## rule, ruleset - 规则管理命令

**nft COMMAND rule \[FAMILY] TABLE CHAIN \[handle HANDLE|index INDEX] STATEMENT...**

- FAMILY 指定簇名
- HANDLE 和 INDEX 指定规则的句柄值或索引值
- STATEMENT 指明该规则的语句

**nft list ruleset \[FAMILY]** # 列出所有规则，包括规则所在的链，链所在的表。i.e. 列出 nftables 中的所有信息。可以指定 FAMILY 来列出指定簇的规则信息

- \[FAMILY] # 清除所有规则，包括表。i.e.清空 nftables 中所有信息。可以指定 FAMILY 来清空指定簇的规则信息

COMMAND

- add # 将规则添加到链的末尾，或者指定规则的 handle 或 index 之后
- insert # 将规则添加到链的开头，或者指定规则的 handle 或 index 之前
- delete # 删除指定的规则。Note:只能通过 handle 删除
- replace # 替换指定规则为新规则

EXAMPLE

- nft add rule inet my_table my_filter_chain tcp dport ssh accept # 在 inet 簇的 my_table 表中的 my_filter_chain 链中添加一条规则，目标端口是 ssh 服务的数据都接受
- nft add rule inet my_table my_filter_chain ip saddr @my_set drop # 创建规则时引用 my_set 集合

## set - 集合管理命令

COMMAND set \[FAMILY] table set { type TYPE; \[flags FLAGS;] \[timeout TIMEOUT ;] \[gc-interval GC-INTERVAL ;] \[elements = { ELEMENT\[,...] } ;] \[size SIZE;] \[policy POLICY;] \[auto-merge AUTO-MERGE ;] } # 各字段解释详见上文 nftables 的 set 与 map 特性介绍

list sets # 列出所有结合

{add | delete} element \[family] table set { element\[,...] } # 在指定集合中添加或删除元素

Note:

- 在输入命令时，使用反斜线 \ 用来转义分号 ; ，这样 shell 就不会将分号解释为命令的结尾。如果是直接编辑 nftables 的配置文件则不用进行转义

COMMAND

- add
- delete # 通过 handle 删除指定的集合
- flush #
- list #

EXAMPLE

- nft add set inet my_table my_set {type ipv4_addr; } # 在 inet 簇的 my_table 表中创建一个名为 my_set 的集合，集合的类型为 ipv4_addr
- nft add set my_table my_set {type ipv4_addr; flags interval;} # 在默认 ip 簇的 my_table 表中创建一个名为 my_set 的集合，集合类型为 ipv4_addr ，标签为 interval。让该集合支持区间
- nft add element inet my_table my_set { 10.10.10.22, 10.10.10.33 } # 向 my_set 集合中添加元素，一共添加了两个元素，是两个 ipv4 的地址

## 字典管理命令

字典

字典是 nftables 的一个高级特性，它可以使用不同类型的数据并将匹配条件映射到某一个规则上面，并且由于是哈希映射的方式，可以完美的避免链式规则跳转的性能开销。

例如，为了从逻辑上将对 TCP 和 UDP 数据包的处理规则拆分开来，可以使用字典来实现，这样就可以通过一条规则实现上述需求。

$ nft add chain inet my_table my_tcp_chain

$ nft add chain inet my_table my_udp_chain

$ nft add rule inet my_table my_filter_chain meta l4proto vmap { tcp : jump my_tcp_chain, udp : jump my_udp_chain }

$ nft list chain inet my_table my_filter_chain

```bash
table inet my_table {
    chain my_filter_chain {
    ...
    meta nfproto ipv4 ip saddr . meta l4proto . udp dport { 10.30.30.30 . udp . bootps } accept
    meta l4proto vmap { tcp : jump my_tcp_chain, udp : jump my_udp_chain }
    }
}
```

和集合一样，除了匿名字典之外，还可以创建命名字典：

$ nft add map inet my_table my_vmap { type inet_proto : verdict ; }

向字典中添加元素：

$ nft add element inet my_table my_vmap { 192.168.0.10 : drop, 192.168.0.11 : accept }

后面就可以在规则中引用字典中的元素了：

$ nft add rule inet my_table my_filter_chain ip saddr vmap @my_vmap

表与命名空间

在 nftables 中，每个表都是一个独立的命名空间，这就意味着不同的表中的链、集合、字典等都可以有相同的名字。例如：

$ nft add table inet table_one

$ nft add chain inet table_one my_chain

$ nft add table inet table_two

$ nft add chain inet table_two my_chain

$ nft list ruleset

```bash
...
table inet table_one {
    chain my_chain {
    }
}
table inet table_two {
    chain my_chain {
    }
}
```

有了这个特性，不同的应用就可以在相互不影响的情况下管理自己的表中的规则，而使用 iptables 就无法做到这一点。

当然，这个特性也有缺陷，由于每个表都被视为独立的防火墙，那么某个数据包必须被所有表中的规则放行，才算真正的放行，即使 table_one 允许该数据包通过，该数据包仍然有可能被 table_two 拒绝。为了解决这个问题，nftables 引入了优先级，priority 值越高的链优先级越低，所以 priority 值低的链比 priority 值高的链先执行。如果两条链的优先级相同，就会进入竞争状态。

总结

希望通过本文的讲解，你能对 nftables 的功能和用法有所了解，当然本文只涉及了一些浅显的用法，更高级的用法可以查看 nftables 的官方 Wiki, 或者坐等我接下来的文章。相信有了本文的知识储备，你应该可以愉快地使用 nftables 实现 Linux 的智能分流了，具体扫一扫下方的二维码参考这篇文章：Linux 全局智能分流方案。
