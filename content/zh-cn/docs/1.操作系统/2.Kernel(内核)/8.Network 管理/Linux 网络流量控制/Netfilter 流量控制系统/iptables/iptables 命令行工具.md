---
title: "iptables 命令行工具"
linkTitle: "iptables 命令行工具"
date: "2023-07-21T12:40"
weight: 20
---

# 概述

> 参考：
> 
> - [Manual(手册)，iptables(8)](https://man7.org/linux/man-pages/man8/iptables.8.html)
> - [Manual(手册)，iptables-extensions(8)](https://man7.org/linux/man-pages/man8/iptables-extensions.8.html)

Man 手册中，将 iptables 分为两部分，基本的 iptables 和用于描述扩展规则的 iptables-extensions，当使用扩展规则时，需要通过 `-m, --match ModuleName` 选项指定要使用的扩展模块的名称。

另外，所谓的扩展规则，其实是对原始 iptables 的扩展，并不仅仅扩展了匹配数据包的规则，还有很多其他的功能。

# Syntax(语法)

**iptables \[OPTIONS] COMMAND \[CHAIN] \[RuleSpecifitcation]**

- **Command** # 指定要执行的具体操作。比如 *增删改查规则/链* 等等。
- **CHAIN** # 指定要执行操作的链。在不指定的时候，默认对所有链进行操作。
  - CHAIN 其实不应该放在这，一般都是 COMMAND 中的组成部分。
- **RuleSpecifitcation = MATCHES TARGET** # 通常用在增加规则时，指定规则的具体规范。由两部分组成：\[MATCHES...] 和 \[TARGET]
  - **MATCHES = \[基本匹配规则] \[扩展匹配规则]** # 匹配条件，可以指定多个。用以筛选出要执行 TARGET 的数据包的条件
  - **TARGET = -j TargetName \[Per-Target-Options]** # 指定匹配到规则的数据包的 Target(目标) 是什么。

## OPTIONS

- **-t, --table TALBLE** # 指定 iptables 命令要对 TABLE 这个表进行操作。`默认值: filter`。省略该选项时，表示默认对 filter 表进行操作。
- **-n, --numeric** # 所有输出以数字的形式展示。IP 地址、端口号等都以数字输出。默认情况下一般是显示主机名、网络名称、服务。
- **--line-numbers** # 显示每个 chain 中的行号
- **-v** # 显示更详细的信息，vv 更详细，vvv 再详细一些
  - pkts # 报文数
  - bytes # 字节数
  - target #
  - prot #
  - in/out # 显示要限制的具体网卡，`*` 为所有
  - source/destination #

## COMMAND

- 增
  - **-A, --append \<CHAIN> \<RuleSpecification>** # 在规则连末尾添加规则
  - **-I, --insert \<CHAIN> \[RuleNum] \<RuleSpecification>** # 在规则链开头添加规则，也可以指定添加到指定的规则号
  - **-N, --new-chain CHAIN** # 创建名为 CHAIN 的自定义规则链
- 删
  - **-F, --flush \[CHAIN \[RuleNum]]** # 删除所有 chain 下的所有规则，也可删除指定 chain 下的指定的规则
  - **-D, --delete \<CHAIN> \<RULE>** # 删除一个 chain 中规则，RULE 可以是该 chain 中的行号，也可以是规则具体配置
  - **-X, --delete-chain \[CHAIN]** # 删除用户自定义的空的 chain
- 改
  - **-P, --policy \<CHAIN> \<TARGET>** # 设置指定的规则链(CHAIN)的默认策略为指定目标(Targe)
  - **-E, --rename-chain \<OldChainName> \<NewChainName>** # 重命名自定义 chain，引用计数不为 0 的自定义 chain，无法改名也无法删除
  - **-R, --replace CHAIN \[RuleNum] \<RuleSpecification>** # 替换指定链上的指定规则
- 查
  - **-L, --list \[CHAIN \[RuleNum]]** # 列出指定 CHAIN 的规则。`默认值: 不指定。i.e.列出所有 CHAIN 的规则`
  - **-S, --list-rules \[CHAIN]** # 以 iptables-save 命令的方式列出指定 CHAIN 的规则。`默认值: 不指定。i.e.列出所有 CHAIN 的规则`
- 其他
  - **-Z, --zero \[CHAIN \[RULE_NUM]]** # 将所有链中的数据包和字节计数器归零，或者仅将给定链归零，或者仅将给定规则归零。

## MATCHES

iptabes 语法中的 MATCHES 部分是整个语法中非常重要且复杂的。该部分用于指定要匹配的数据包条件，只有符合条件的数据包，才会被 Netfilter 框架处理。

MATCHES 规则可以分为两类：

- 基本匹配规则
- 扩展匹配规则

我们可以指定一个或多个匹配规则，很多匹配规则前可以添加 `!` 符号，即可将该匹配取反（比如：不加叹号是匹配某个 IP 地址，加了叹号表示除了这个地址外都匹配）

MATCHES 的语法是一种类似命令行参数的写法

### 基本匹配规则

> 在 Man 手册中，这部分的介绍在 OPTIONS 下的 PARAMETERS 章节。

基本匹配规则语法

- **-s, --source IP/MASK** # 指定规则中要匹配的来源地址的 IP/MASK
- **-d, --destination IP/MASK** # 指定规则中要匹配的目标地址的 IP/MASK
- **-i, --in-interface 网卡名称** # 指定数据流入规则中要匹配的网卡，仅用于 PREROUTING、INPUT、FORWARD 链
- **-o, --out-interface 网卡名称** # 指定数据流出规则中要匹配的网卡，仅用于 FORWARD、OUTPUT、POSTROUTING 链
- **-p, --protocol PROTOCOL** # 指定规则中要匹配的协议，即 ip 首部中的 protocols 所标识的协议
  - 可用的协议有 tcp、udp、icmp、等等
- **-m, --match MATCH** # 使用[扩展匹配语法](#扩展匹配语法)指定扩展匹配规则
- **-j, --jump TARGET** # 指定 [TARGET](#TARGET)

### 扩展匹配规则

扩展匹配语法：`-m ModuleName Per-Match-Options`，即指定要使用的扩展模块以及该模块可用的选项。

详见下文 [扩展匹配语法](#扩展匹配语法)

## TARGET

**TARGET = -j TargetName \[Per-Target-Options]**

指定规则中的 target(目标)是什么。即“如果数据包匹配上规则之后应该做什么”。如果目标是自定义链，则指明具体的自定义 Chain 的名称。TARGET 都有哪些详见 Netfilter 流量控制系统。下面是各种 TARGET 的类型：

- **ACCEPT** # 允许流量通过
- **REJECT** # 拒绝流量通过
  - --reject-with icmp-host-prohibited # 通过 icmp 协议显示给客户机一条消息:主机拒绝(icmp-host-prohibited)
- **DROP** # 丢弃，不响应，发送方无法判断是被拒绝
- **RETURN** # 返回调用链
- **MARK** # 做防火墙标记
  - 用于 nat 表的 target
    - DNAT|SNAT # {目的|源}地址转换
    - REDIRECT # 端口重定向
    - MASQUERADE # 地址伪装类似于 SNAT，但是不用指明要转换的地址，而是自动选择要转换的地址，用于外部地址不固定的情况
  - 用于 raw 表的 target
    - NOTRACK # raw 表专用的 target，用于对匹配规则进行 notrack(不跟踪)处理
- **LOG** # 记录日志信息
- **引用自定义链** # 直接使用“-j|-g 自定义链的名称”即可，让基本 5 个 Chain 上匹配成功的数据包继续执行自定义链上的规则。

# 扩展匹配语法

> 参考：
> 
> - [Manual(手册)，iptables-extensions(8)](https://man7.org/linux/man-pages/man8/iptables-extensions.8.html)

**iptables \[OPTIONS] COMMAND \[CHAIN] \[基本匹配规则] -m ExtendedModuleName --ExtendedModuleArgs -j TARGET**

官方文档的描述有点问题，-m 后面应该是指扩展模块名称

- ExtendedModuleName # 扩展模块名称
- ExtendedModuleArgs # 扩展模块对应的参数

**注意，所有的扩展匹配语法必须使用 `-m` 选项指定要使用的扩展模块**，下文的描述也是直接对扩展模块名称进行记录。

iptables 可以使用带有 -m 或 --match 选项的扩展的数据包匹配模块，这些模块可以提供更加强大的匹配数据包的能力，而不仅仅只是源目IP。之后，根据特定模块的不同，可以使用各种额外的命令行选项。我们可以在一行中指定多个扩展匹配模块。扩展匹配模块按照规则中指定的顺序进行评估。

- 比如使用 `-p tcp`，即可在语句后面添加 `-m tcp` 以调用 tcp 扩展模块。
- 然后还可以使用 tcp 模块的专属参数，比如 --dport、--sport 等等。

我们可以在指定模块后使用 -h 或 --help 选项以获取特定模块的帮助信息，比如 `iptables -m tcp --help` 将会列出 tcp 模块所支持的所有选项。

> 备注：使用 -p, --protocol 选项时， iptables 将尝试加载与协议同名的模块，以使用该模块的扩展匹配选项。比如 `-p tcp -m tcp --dport 22` 可以简写为 `-p --dport 22`，因为默认加载同名的扩展模块

扩展模块分为如下几类：

- 通用扩展模块
- 协议相关的模块
- 其他模块

## 通用扩展模块

指定具体的扩展匹配名以及该扩展匹配的匹配规则

**conntrack 模块**

- **--ctstate CTState1\[,CTState2...]** # 匹配指定的名为 CTState(conntrack State) 的[连接追踪](/docs/1.操作系统/2.Kernel(内核)/8.Network%20管理/Linux%20网络流量控制/Netfilter%20流量控制系统/Connection%20Tracking(连接跟踪)机制.md)状态。可用的状态有{INVALID|ESTABLISHED|NEW|RELATED|UNTRACKED|SNAT|DNAT}

> -m state --state STATE1\[,STATE2,....] # conntrack 的老式用法，慢慢会被淘汰

**comment 模块**

- **--comment STRING** # 向规则添加注释，最多 256 个字符。比如：`iptables -A INPUT -i eth1 -m comment --comment "my local LAN"`

**set 模块**

**--match-set SetName {src|dst}..**. # 匹配指定的{源|目标}IP 是名为 SetName 的 ipset 集合
  - 其中 FLAG 是逗号分隔的 src 和 dst 规范列表，其中不能超过六个。

**iprange 模块**

- **{--src-range|--dst-range} IP1-IP2** # 匹配的指定的{源|目标}IP 范围

**sting 模块**

- **--algo {bm|kmp}** # 指明使用哪个搜索算法
- **--string PATTERN** # 指明要匹配的字符串，用于检查报文中出现的字符串(e.g.某个网页出现某个字符串则拒绝)

## 协议相关的扩展模块

**multiport** # 仅可使在 **tcp**, **udp**, **udplite**, **dccp**, **sctp** 这几个协议中使用

- **{--dports|--sports} NUM** # 让 tcp 匹配多个端口，可以是目标端口(dport)或者源端口(sport)

**tcp 模块**

  - **--sport, --source-port NUM** # 指定规则中要匹配的来源端口号
  - **--dport, --destination-port NUM** # 指定规则中要匹配的目标端口号
  - **--tcp-flags LIST1 LIST2** # 检查 LIST1 所指明的所有标志位，且这其中 LIST2 所表示出的所有标志位必须为 1，而余下的必须为 0,；没有 LIST1 中指明的，不做检查(e.g.--tcp-flags SYN,ACK,FIN,RST SYN)。LIST 包括“SYN ACK FIN RST URG PSH ALL NONE”

**udp 模块**

  - **--dport NUM** # 指定规则中要匹配的目标端口号
  - **--sport NUM** # 指定规则中要匹配的来源端口号

**icmp 模块**

  - **-m \[icmp] --icmp-type TYPE** # 指定 icmp 的类型，具体类型可以搜 icmp type 获得，可以是数字

# EXAMPLE

注意：在使用 iptables 命令的时候，为了防止配置问题导致网络不通，可以先设置一个定时任务来修改 iptables 规则或者 iptables XX && sleep 5 && iptables -P INPUT ACCEPT && iptables -P OUTPUT ACCEPT

## Fileter 表的配置

### INPUT 默认为 DROP 情况下

- iptables -P INPUT DROP

给 INPUT 链添加一条接受 192.168.19.64/27 这个网段数据包的规则

- iptables -I INPUT -s 192.168.19.64/27 -j ACCEPT

允许指定的网段访问本机的 22 号端口

- iptables -A INPUT -s 192.168.1.0/24 -p tcp -m tcp --dport 22 -j ACCEPT

允许所有机器进入本机的 1000 到 1100 号端口

- iptables -A INPUT -p tcp -m tcp --dport 1000:1100 -j ACCEPT

允许源地址是 10.10.100.4 到 10.10.100.10 这 7 个 ip 的流量进入本机

- iptables -A INPUT -m iprange --src-range 10.10.100.4-10.10.100.10 -j ACCEPT

允许源地址是 10.10.100.4 和 10.10.100.8 且目的端口是 1935 和 4000 的流量进入本机

- iptables -A INPUT -m iprange --src-range 10.10.100.4,10.10.100.8 -p tcp -m multiport --dports 1935,4000 -j ACCEPT

允许源地址是 ipset(名为 cdn2) 中的所有 IP，且目标端口为 80 的所有数据包通过

- iptables -A INPUT -p tcp -m set --match-set cdn2 src --dports 80 -j ACCEPT

允许源地址是 ipset(名为 cdn1) 中的所有 IP 通过

- iptables -A INPUT -m set --match-set cdn1 src -j ACCEPT

### INPUT 默认为 ACCEPT 情况下

- iptables -P INPUT DROP

只放开部分 IP 的 22 端口

> 注意：下面的顺序不可变，由于是顺序执行，指定源地址的先允许之后就不会再拒绝了，凡是没有允许过的源 IP，都会被 DROP

```bash
-P INPUT ACCEPT
-P FORWARD ACCEPT
-P OUTPUT ACCEPT
-A INPUT -s 116.63.160.74/32 -p tcp -m tcp --dport 22 -j ACCEPT
-A INPUT -s 122.9.154.106/32 -p tcp -m tcp --dport 22 -j ACCEPT
-A INPUT -p tcp -m tcp --dport 22 -j DROP
```

所有 ping 本机的数据包全部丢弃，禁 ping。给 INPUT 链添加一条禁止 icmp 协议的规则

- iptables -I INPUT -p icmp -j DROP

### 其他

显示 INPUT 链中的所有规则，并显示规则行号

- iptables -L INPUT --line-numbers

删除 INPUT 链中的第 11 条规则

- iptables -D INPUT 11

## NAT 表的配置

凡是基于 tcp 协议访问本机 80 端口，且目的地址是 110.119.120.1 的数据包。全部把目的地址转变为 192.168.20.2，且目的端口转换为 8080

  - iptables -t nat -A PREROUTING -d 110.119.120.1 -p tcp --dport 80 -j DNAT --to-destination 192.168.20.2:8080

将源地址网段是 192.168.122.0/24 并且目的地址网段不是 192.168.122.0/24 的全部转换成 123.213.1.5 这个地址。常用于公司内使用私网 IP 的设备访问互联网使用

  - iptables -t nat -A POSTROUTING -s 192.168.122.0/24 ! -d 192.168.122.0/24 -j SNAT --to-source 123.213.1.5

将所有到达 eth0 网卡 25 端口的流量转发到 2525 端口。也叫端口转发

  - iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 25 -j REDIRECT --to-port 2525

## RAW 表配置

所有来自 10.0.9.0/24 网段的数据包，都不跟踪

  - iptables -t raw -A PREROUTING -s 10.0.9.0/24 -j NOTRACK

# 其他关联命令行工具

## iptables-save - 将 iptables 规则转储到标准输出

该命令输出的内容更容易被人类阅读，可以用重定向把内容保存到文件中

该命令显示出的信息说明

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fadrg5/1616165483878-a0be9032-9c98-4971-922e-592c61529d86.jpeg)

**Syntax(语法)**

iptables-save \[-M,--modprobe modprobe] \[-c] \[-t table]

**EXAMPLE**

- iptables-save > /etc/sysconfig/iptables.rules

## iptables-restore - 从标准输入恢复 iptables 规则，可以视同重定向通过文件来读取到标准输入

**EXAMPLE**

- iptables-restore < /etc/sysconfig/iptables.rules

# ipset - IP 集合的管理工具

ipset 是 iptables 的一个协助工具。可以通过 ipset 设置一组 IP 并对这一组 IP 统一命名，然后在 iptables 里的匹配规则里通过名字直接引用该组 IP，并对这组 IP 进行流量控制。注意：由于是 iptables 规则引用，所以我直接修改 ipset 集合里的 IP，并不用重启 iptables 服务，就可以直接生效。这类似于域名解析，我的机器指定访问 baidu.com，至于百度公司他们怎么更改 IP 与域名绑定的关系，作为用户都不用更更改 baidu.com 这个域名。

## Syntax(语法)

**ipset \[OPTIONS] COMMAND \[COMMAND-OPTIONS]**

COMMANDS：Note：ENTRY 指的就是 ip 地址

- create SETNAME TYPENAME \[type-specific-options] # 创建一个新的集合。Create a new set
- add SETNAME ENTRY # 向指定集合中添加条目。i.e.添加 ip。Add entry to the named set
- del SETNAME ENTRY # 从指定集合中删除条目 Delete entry from the named set
- test SETNAME ENTRY # 测试指定集合中是否包含该条目 Test entry in the named set
- destroy \[SETNAME] # 摧毁全部或者指定的集合 Destroy a named set or all sets
- list \[SETNAME] # 列出全部或者指定集合中的条目 List the entries of a named set or all sets
- save \[SETNAME] # 将指定的集合或者所有集合保存到标准输出
- restore # 还原保存的 ipset 信息
- flush \[SETNAME] # 删除全部或者指定集合中的所有条目 Flush a named set or all sets
- rename FROM-SETNAME TO-SETNAME # Rename two sets
- swap FROM-SETNAME TO-SETNAME # 交换两个集合中的内容 Swap the contect of two existing sets

OPTIONS

- **-exist** # 在 create 已经存在的 ipset、add 已经存在的 entry、del 不存在的 entry 时忽略错误。
- **-f** # 在使用 save 或者 restore 命令时，可以指定文件，而不是从标准输出来保存或者还原 ipset 信息

{ -exist | -output { plain | save | xml } | -quiet | -resolve | -sorted | -name | -terse | -file filename }

EXAMPLE

- **ipset list** # 列出 ipset 所设置的所有 IP 集合
- **ipset create lichenhao hash:net** # 创建一个 hash:net 类型的名为 lichenhao 的 ipset
- **ipset add lichenhao 1.1.1.0/24** # 将 1.1.1.0/24 网段添加到名为 lichenhao 的 ipset 中
- **ipset flush** # 清空所有 ipset 下的 ip
- **ipset restore -f /etc/sysconfig/ipset** # 从/etc/sysconfig/ipset 还原 ipset 的集合和条目信息
