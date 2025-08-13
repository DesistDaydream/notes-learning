---
title: "nft CLI"
linkTitle: "nft CLI"
weight: 20
---

# 概述

> 参考：
>
> - [Manual, nftables](https://www.netfilter.org/projects/nftables/manpage.html)
>     - https://www.mankier.com/8/nft

nft 是一个命令行工具，用于在 Linux 内核中设置、维护和检查 nftables 框架中的数据包过滤和分类规则。Linux 内核子系统被称为 nf_tables，其中 "nf" 代表 Netfilter。

**nft \[OPTIONS] \[VERB] \[COMMANDS]**

VERB 是动词，用来表示要执行的操作。不同 COMMAND 可用的动词不太一样，常见的包括：

- add # 添加对象。可以添加到指定的对象之后
- create # 与 add 命令类似，但是如果对象已经存在，则返回错误信息。
- insert # 插入对象。可以添加到指定的对象之前
- replace # 替换对象
- delete # 删除对象。不管表中是否有内容都一并删除
- flush # 清空对象
- list # 列出对象
- etc.

COMMANDS 包括：

- ruleset # 规则集管理命令
- table # 表管理命令
- chain # 链管理命令
- rule # 规则管理命令
- set # 集合管理命令
- map # 字典管理命令

> [!Notes]
> - 该 COMMANDS 与后面子命令中的 COMMAND 不同，前者是 nft 命令下的子命令，后者是 nft 命令下子命令的子命令
> - nft 子命令默认对 ip 族进行操作，当指定具体的 FAMILY 时，则对指定的族进行操作

## OPTIONS

如何加载规则集的输入选项

- **-f,--file**(STRING) # 从指定的文件 FILE 中读取 netfilter 配置加载到内核中

修改 list ruleset 命令输出的规则集列表格式

- **-a,--handle** # 在使用命令获得输出时，显示每个对象的句柄
  - Note：handle(句柄)在 nftables 中，相当于标识符，nftables 中的每个对象都有一个 handle。

命令输出格式

- **-j, --json** # 以 [JSON](/docs/2.编程/无法分类的语言/JSON.md) 格式化输出内容。从 [libnftables-json(5)](https://www.mankier.com/5/libnftables-json) 这里参考结构描述
- **-e,--echo** # 回显已添加、插入或替换的内容

# table - 表管理命令

**nft VERB table \[FAMILY] TABLE**

- FAMILY 指定族名
- TABLE 为表的名称

VERB

- add # 添加指定族下的表。
- create # 与 add 命令类似，但是如果表已经存在，则返回错误信息。
- delete # 删除指定的表。不管表中是否有内容都一并删除
- flush # 清空指定的表下的所有规则，保留链
- list # 列出指定的表的所有链，及其链中的规则

EXAMPLE

- nft add table my_table # 创建一个 ip 族的，名为 my_table 的表
- nft add table inet my_table # 创建一个 inet 族的，名为 my_table 的表
- nft list tables # 列出所有的表，不包含表中的链和规则
- nft list table inet my_table # 列出 inet 族的名为 my_table 的表及其链和规则

# chains - 链管理命令

**nft VERB chain \[FAMILY] TABLE CHAIN \[{ type TYPE hook HOOK \[device DEVICE] priority PRIORITY; \[policy POLICY;] }]**

- FAMILY 指定族名
- TABLE 指定表名
- CHAIN 指定链名
- TYPE 指定该链的类型
- HOOK 指定该链作用在哪个 hook 上
- DEVICE 指定该链作用在哪个网络设备上
- PRIORITY 指定该链的优先级
- POLICY 指定该链的策略(i.e.该链的默认策略，accept、drop 等等。)

> [!Notes]
> - 在输入命令时，使用反斜线 `\` 用来转义分号 `;` ，这样 Shell 就不会将分号解释为命令的结尾。如果是直接编辑 nftables 的配置文件则不用进行转义
> - PRIORITY 采用整数值，可以是负数，值较小的链优先处理。

VERB

- add # 在指定的表中添加一条链
- create # 与 add 命令类似，但是如果链已经存在，则返回错误信息。
- delete # 删除指定的链。该链不能包含任何规则，或者被其它规则作为跳转目标，否则删除失败。
- flush #
- list # 列出指定表下指定的链，及其链中的规则
- rename #

EXAMPLE

- nft list chains # 列出所有的链
- nft add chain inet my_table my_utility_chain # 在 inet 族的 my_table 表上创建一个名为 my_utility_chain 的常规链，没有任何参数
- nft add chain inet my_table my_filter_chain{type filter hook input priority 0;} # 在 inet 族的 my_table 表上创建一个名为 my_filter_chain 的链，链的类型为 filter，作用在 input 这个 hook 上，优先级为 0
- nft list chain inet my_table my_filter_chain # 列出 inet 族的 my_table 表下的 my_filter_chain 链的信息，包括其所属的表和其包含的规则

# rule - 规则管理命令

**nft VERB rule \[FAMILY] TABLE CHAIN \[handle HANDLE|index INDEX] STATEMENT...**

- FAMILY 指定族名
- HANDLE 和 INDEX 指定规则的句柄值或索引值
- STATEMENT 指明该规则的语句

VERB

- add # 将规则添加到链的末尾，或者指定规则的 handle 或 index 之后
- insert # 将规则添加到链的开头，或者指定规则的 handle 或 index 之前
- delete # 删除指定的规则。Note:只能通过 handle 删除
- replace # 替换指定规则为新规则

EXAMPLE

- nft add rule inet my_table my_filter_chain tcp dport ssh accept # 在 inet 族的 my_table 表中的 my_filter_chain 链中添加一条规则，目标端口是 ssh 服务的数据都接受
- nft add rule inet my_table my_filter_chain ip saddr @my_set drop # 创建规则时引用 my_set 集合

# ruleset - 规则集管理命令

**nft list ruleset \[FAMILY]** # 列出所有规则，包括规则所在的链，链所在的表。i.e. 列出 nftables 中的所有信息。可以指定 FAMILY 来列出指定族的规则信息

**nft flush ruleset \[FAMILY]** # 清除所有规则，包括表。i.e.清空 nftables 中所有信息。可以指定 FAMILY 来清空指定族的规则信息

# set - 集合管理命令

**nft VERB set \[FAMILY] table set { type TYPE; \[flags FLAGS;] \[timeout TIMEOUT ;] \[gc-interval GC-INTERVAL ;] \[elements = { ELEMENT\[,...] } ;] \[size SIZE;] \[policy POLICY;] \[auto-merge AUTO-MERGE ;] }**

> 在输入命令时，使用反斜线 `\` 用来转义分号 `;` ，这样 shell 就不会将分号解释为命令的结尾。如果是直接编辑 nftables 的配置文件则不用进行转义
>
> 各字段解释详见 [Nftables](/docs/1.操作系统/Kernel/Network/Linux%20网络流量控制/Netfilter/Nftables/Nftables.md#Set) 中的 Set 章节

list sets # 列出所有结合

{add | delete} element \[family] table set { element\[,...] } # 在指定集合中添加或删除元素


VERB

- add
- delete # 通过 handle 删除指定的集合
- flush #
- list #

EXAMPLE

- nft add set inet my_table my_set {type ipv4_addr;} # 在 inet 族的 my_table 表中创建一个名为 my_set 的集合，集合的类型为 ipv4_addr
- nft add set my_table my_set {type ipv4_addr; flags interval;} # 在 ip 族的 my_table 表中创建一个名为 my_set 的集合，集合类型为 ipv4_addr ，标签为 interval。让该集合支持区间
- nft add element inet my_table my_set { 10.10.10.22, 10.10.10.33 } # 向 my_set 集合中添加元素，一共添加了两个元素，是两个 ipv4 的地址
- 删除元素。删除 my_table 表中，ssh_allowed_nets 集合内的 183.192.0.0/10 元素
    - `nft delete element inet my_table ssh_allowed_nets { 183.192.0.0/10 }`

# 字典管理命令 TODO

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

# 主表达式


# Payload 表达式

## Conntrack 表达式

https://www.mankier.com/8/nft#Payload_Expressions-Conntrack_Expressions

Conntrack 表达式是指与数据包关联的连接跟踪条目的元数据。用人话说：对 [Connection tracking for Netfilter](/docs/1.操作系统/Kernel/Network/Linux%20网络流量控制/Netfilter/Connection%20tracking%20for%20Netfilter.md) 状态进行匹配

# 最佳实践

## 基本示例

从 0 开始创建规则，实现只允许 172.16.50.0/24 访问本设备 22 端口。保存规则后，再删除规则。通过文件加载规则

**创建**

```bash
# 创建表。名为 filter_demo，属于 inet 族
nft add table inet filter_demo
# 创建链
nft add chain inet filter_demo input_demo { type filter hook input priority 0 \; }
# 创建规则
nft add rule inet filter_demo input_demo ip saddr 172.16.50.0/24 tcp dport 22 accept
nft add rule inet filter_demo input_demo tcp dport 22 drop
```

**查看**

```bash
# 检查规则及每个对象的句柄
nft -a list ruleset inet
```

结果如下

```bash
table inet filter_demo { # handle 10
        chain input_demo { # handle 1
                type filter hook input priority filter; policy accept;
                ip saddr 172.16.50.0/24 tcp dport 22 accept # handle 2
                tcp dport 22 drop # handle 3
        }
}
```

此时，除了 172.16.50.0/24，其他 IP 将无法访问本机的 22 端口，<font color="#ff0000">规则是即时生效</font>的

**保存**当前配置

```bash
nft list ruleset > /etc/nftables.conf
```

**删除**。<font color="#ff0000">注意删除顺序</font>。若递归删除上层对象（e.g. 链、表）则不需要注意删除顺序

```bash
# 仅删除规则。删除 filter_demo 表下的 input_demo 链中 3 号和 2 号句柄规则
nft delete rule inet filter_demo input_demo handle 3
nft delete rule inet filter_demo input_demo handle 2
```

> [!Note] 示例是只删除了具体的一条规则，但是可以通过删除上级对象同时删除该对象包含的所有内容，比如删除 链，也会对应删除链中的规则，以此类推。甚至可以直接清空所有内容
> - `nft delete chain inet filter_demo handle 1` # 删除链。
> - `nft delete table inet handle 10` # 删除表
> - `nft flush ruleset` 清空全部内容（虽然用的是 ruleset，但是会连带着表、链一起删了）

**加载**。从 nftables.conf 文件中，将配置规则加载到系统中

```bash
nft -f /etc/nftables.conf
```

## 配置 IP 的集合

继续上面的[基本示例](#基本示例)实践

**创建集合并插入元素**

```bash
# 1. 创建一个名为 ssh_whitelist 的 set
nft add set inet filter_demo ssh_whitelist { type ipv4_addr \; flags interval \; }

# 2. 向 set 中添加 IP 和网段
nft add element inet filter_demo ssh_whitelist { 10.10.11.54 }
nft add element inet filter_demo ssh_whitelist { 172.16.50.0/24 }
```

**变更新规则**

有两种变更规则的方式

1. 插入新规则后删除老规则
```bash
nft insert rule inet filter_demo input_demo position 2 ip saddr @ssh_whitelist tcp dport 22 accept
```
2. 替换老规则
```bash
nft replace rule inet filter_demo input_demo handle 2 ip saddr @ssh_whitelist tcp dport 22 accept
```

这里使用替换的方式，效果如下：

```bash
~]# nft -a list ruleset inet
table inet filter_demo { # handle 10
        set ssh_whitelist { # handle 4
                type ipv4_addr
                flags interval
                elements = { 10.10.11.54, 172.16.50.0/24 }
        }

        chain input_demo { # handle 1
                type filter hook input priority filter; policy accept;
                ip saddr @ssh_whitelist tcp dport 22 accept # handle 2
                tcp dport 22 drop # handle 3
        }
}
```

## 其他

对 related,established 状态的连接全部放通

`nft add rule ip my_table my_chain ct state related,established accept`