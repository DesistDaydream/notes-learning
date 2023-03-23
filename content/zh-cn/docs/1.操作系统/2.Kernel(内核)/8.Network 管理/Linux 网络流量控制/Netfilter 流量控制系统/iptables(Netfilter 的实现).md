---
title: iptables(Netfilter 的实现)
---

# 概述

> 参考：
>
> - [Manual(手册),iptables(8)](https://man7.org/linux/man-pages/man8/iptables.8.html)
> - [Netfilter 官方文档，iptables 教程](https://www.frozentux.net/iptables-tutorial/iptables-tutorial.html)

iptables 是一组工具合集的统称，其中包括 iptables、ip6tables、arptables、ebtables 等，用于与 netfilter 模块进行交互的 CLI 工具

iptables 和 ip6tables 用于建立、维护和检查 Linux 内核中的 IPv4 和 IPv6 包过滤规则表。可以定义几个不同的表中的各种规则，也可以定义用户定义的链。并把已经定义的规则发送给 netfilter 模块。

## 四表(Table) Note:四表是 iptables 框架中的概念，不是 netfilter 中的

iptables 框架将流量抽象分为 4 类：过滤类、网络地址转换类、拆解报文类、原始类。每种类型的链作用在 Netfilter 系统中的 Hook 各不不相同，每种类型具有不同的功能。每一类都称为一张表。比如 fileter 表用来在指定链上检查流量是否可以通过，nat 表用来在指定链上检查流量是否可以进行地址转换，等等。Note：不是所有表都可以在所有链上具有规则，下表是 4 个表在 5 个 Hook 上的可用关系。

| 表名\链名 | PREROUTING | INPUT | FORWARD | OUTPUT | POSTROUTING |
| --------- | ---------- | ----- | ------- | ------ | ----------- |
| filter    |            | 可用  | 可用    | 可用   |             |
| nat       | 可用       | 可用  |         | 可用   | 可用        |
| mangle    | 可用       | 可用  | 可用    | 可用   | 可用        |
| raw       | 可用       |       |         | 可用   |             |

iptables 中有默认的内置 4 个表，每个表的名称就是其 chain 类型的名称

### filter(过滤器) # 过滤，防火墙

- 该类型的链可作用在以下几个 Hook 点上：INPUT、FORWARD、OUTPUT

### nat(网络地址转换) # 网络地址转换

- 该类型的链可作用在以下几个 Hook 点上：PREROUTING(DNAT)、INPUT、OUTPUT、POSTROUTING(SNAT)
  - 其中 PREROUTING 与 POSTROUTING 是流量经过物理网卡设备时做 nat 的地方
  - 其中 INPUT 与 OUTPUT 则是主机内部从网络栈直接下来的流量做 nat 的地方。e.g.从主机一个服务发送数据到同一个主机另一个服务的流量

### mangle # 拆解报文，做出修改，封装报文

- 该类型的链可作用在以下几个 Hook 点上：PREROUTING、INPUT、FORWARD、OUTPUT、POSTROUTING

### raw(原始) # 用于跳过 nat 表以及连接追踪机制(ip_conntrack)的处理，详见 [连接跟踪系统](/docs/IT学习笔记/1.操作系统/2.Kernel(内核)/8.Network%20 管理/Linux%20 网络流量控制/Connnection%20Tracking(连接跟踪).md Tracking(连接跟踪).md)

- 该类型的链可作用在以下几个 Hook 点上：PREROUTING、OUTPUT

整个表只有这一个作用，且仅有一个 target 就是 NOTRACK。具有最高优先级，所有流量先在两个 Hook 的 raw 功能上进行检查。一旦在 raw 上配置了规则，则 raw 表处理完成后，跳过 nat 表和 ip_conntrack 处理，i.e.不再做地址转换和数据包的链接跟踪处理，不把匹配到的数据包保存在“链接跟踪表”中。常用于那些不需要做 nat 的情况下以提高性能。e.g.大量访问的 web 服务器，可以让 80 端口不再让 iptables 做数据包的链接跟踪处理 ，以提高用户的访问速度。不过该功能不影响其余表的连接追踪追踪功能的正常使用，依然会有记录写到连接追踪文件中去

1. 该功能的起源：iptables 表有大小的上限，如果每个数据包进来都要进行检查，会非常影响性能，可以对那些不用进行 nat 功能的数据进行放弃后面的检查，i.e.可以在 raw 配置然后直接让这些数据包跳过后面的表对该数据包的处理

Note：四表的优先级从高到低依次为：raw-mangle-nat-filter，i.e.数据到达某个 Hook 上，则会优先使用优先级最高类型的链来处理数据包。其实，iptables 表的作用更像是用来划分优先级的。

## iptables 处理链上规则的顺序以及规范

注意：每个数据包在 CHAIN 中匹配到适用于自己的规则之后，则直接进入下一个 CHAIN，而不会遍历 CHAIN 中每条规则去挨个匹配适用于自己的规则。比如下面两种情况

INPUT 链默认 DROP，匹配第一条：目的端口是 9090 的数据 DROP，然后不再检查下一项，那么 9090 无法访问

    -P INPUT DROP
    -A INPUT -p tcp -m tcp --dport 9090 -j DROP
    -A INPUT -p tcp -m tcp --dport 9090 -j ACCEPT

INPUT 链默认 DROP，匹配第一条目的端口是 9090 的数据 ACCEPT，然后不再检查下一条规则，则 9090 可以访问

    -P INPUT DROP
    -A INPUT -p tcp -m tcp --dport 9090 -j ACCEPT
    -A INPUT -p tcp -m tcp --dport 9090 -j DROP

# iptables 关联文件与配置

/etc/sysconfig/iptables # 该文件存放用户定义的规则信息，每次重启 iptabels 后，都会读取该配置文件信息并应用到系统中
/etc/sysconfig/iptables-conf # 该文件存放 iptables 工具的具体配置信息
/run/xtables.lock # 该文件在 iptables 程序启动时被使用，以获取排他锁

- 可以通过 `XTABLES_LOCKFILE` 环境变量修改 iptables 需要使用 xtalbes.lock 文件的路径

# iptables 命令行工具详解

## Syntax(语法)

**iptables \[-t TABLE] \[OPTIONS] SubCOMMAND CHAIN \[RuleSpecifitcation]**

- subCommand # 指对命令进行什么操作，CHAIN 指定要执行操作的链
- RuleSpecifitcation=MATCHES TARGET # 由两部分组成 \[MATCHES...] 和 \[TARGET]
  - matches 由一个多个参数组成。

### OPTIONS

- **-t TALBLE** # 指定 iptables 命令要对 TABLE 这个表进行操作，默认 filter 表
- **-n** # 不显示域名，直接显示 IP
- **--line-numbers** # 显示每个 chain 中的行号
- **-v** # 显示更详细的信息，vv 更详细，vvv 再详细一些
  - pkts # 报文数
  - bytes # 字节数
  - target #
  - prot #
  - in/out # 显示要限制的具体网卡，\*为所有
  - source/destination #
- **-S** # 以人类方便阅读的方式打印出来 iptables 规则

### SubCOMMAND

- 增
  - **-I \<CHAIN> \[RuleNum] \<RuleSpecification>** # 在规则链开头加入规则详情，也可以指定添加到指定的规则号
  - **-A \<CHAIN> \<RuleSpecification>** # 在规则连末尾加入规则详情
  - **-N ChainName** # 创建名为 ChainName 的自定义规则链
- 删
  - **-F \[CHAIN \[RuleNum]]** # 删除所有 chain 下的所有规则，也可删除指定 chain 下的指定的规则
  - **-D \<CHAIN> \<RULE>** # 删除一个 chain 中规则，RULE 可以是该 chain 中的行号，也可以是规则具体配置
  - **-X \[CHAIN]** # 删除用户自定义的空的 chain
- 改
  - **-P \<CHAIN> \<TARGET>** # 设置指定的规则链(CHAIN)的默认策略为指定目标(Targe)
  - **-E \<OldChainName> \<NewChainName>** # 重命名自定义 chain，引用计数不为 0 的自定义 chain，无法改名也无法删除
  - **-R** # 替换指定链上的指定规则
- 查
  - **-L \[CHAIN \[RuleNum]]** # 列出防火墙所有 CHAIN 的配置，可以列出指定的 CHAIN 的配置

### MATCHES

一个或多个 parameters(参数) 构成 MATCHES # i.e.Match(匹配)相关的 OPTIONS。在每个 parameters 之前添加!，即可将该匹配取反(比如：不加叹号是匹配某个 IP 地址，加了叹号表示除了这个地址外都匹配)

MATCHES=\[-m] MatchName \[Per-Match-Options]

基本匹配规则

- **-s IP/MASK** # 指定规则中要匹配的来源地址的 IP/MASK
- **-d IP/MASK** # 指定规则中要匹配的目标地址的 IP/MASK
- **-i 网卡名称** # 指定数据流入规则中要匹配的网卡，仅用于 PREROUTING、INPUT、FORWARD 链
- **-o 网卡名称** # 指定数据流出规则中要匹配的网卡，仅用于 FORWARD、OUTPUT、POSTROUTING 链
- **-p tcp|udp|icmp** # 指定规则中要匹配的协议，即 ip 首部中的 protocols 所标识的协议

扩展匹配规则(ExtendedMatch)

使用格式：-m ExtendedMatchName --MatchRule

通用的扩展匹配，指定具体的扩展匹配名以及该扩展匹配的匹配规则

- **-m conntrack --ctstate CTState1\[,CTState2...]** # 匹配指定的名为 CTState 的[连接追踪](/docs/IT学习笔记/1.操作系统/2.Kernel(内核)/8.Network%20 管理/Linux%20 网络流量控制/Netfilter%20 流量控制系统/Connection%20Tracking(连接跟踪)机制.md Tracking(连接跟踪)机制.md)状态。CTState 为 conntrack State，可用的状态有{INVALID|ESTABLISHED|NEW|RELATED|UNTRACKED|SNAT|DNAT}
  - -m state --state STATE1\[,STATE2,....] # conntrack 的老式用法，慢慢会被淘汰
- **-m set --match-set SetName {src|dst}..**. # 匹配指定的{源|目标}IP 是名为 SetName 的 ipset 集合
  - 其中 FLAG 是逗号分隔的 src 和 dst 规范列表，其中不能超过六个。
- -**m iprange {--src-range|--dst-range} IP1-IP2** # 匹配的指定的{源|目标}IP 范围
- **-m sting --MatchRule** # 指明要匹配的字符串，用于检查报文中出现的字符串(e.g.某个网页出现某个字符串则拒绝)
  - OPTIONS
    - **--algo {bm|kmp}** # 指明使用哪个搜索算法

基于基本匹配的扩展匹配

- -p tcp 的扩展匹配
  - **-m \[tcp] --dport NUM** # 指定规则中要匹配的目标端口号
  - **-m \[tcp] --sport NUM** # 指定规则中要匹配的来源端口号
  - **-m multiport {--dport|--sport} NUM** # 让 tcp 匹配多个端口，可以是目标端口(dport)或者源端口(sport)
  - **-m \[tcp] --tcp-flags LIST1 LIST2** # 检查 LIST1 所指明的所有标志位，且这其中 LIST2 所表示出的所有标志位必须为 1，而余下的必须为 0,；没有 LIST1 中指明的，不做检查(e.g.--tcp-flags SYN,ACK,FIN,RST SYN)。LIST 包括“SYN ACK FIN RST URG PSH ALL NONE”
- -p udp 的扩展匹配：--dport、--sport 与 tcp 的扩展匹配用法一样
- -p icmp 的扩展匹配
  - **-m \[icmp] --icmp-type TYPE** # 指定 icmp 的类型，具体类型可以搜 icmp type 获得，可以是数字

### TARGET

**TARGET =-j|g TARGET \[Per-Target-Options]** # 指定规则中的目标(target)是什么。即“如果数据包匹配上规则之后应该做什么”。如果目标是自定义链，则指明具体的自定义 Chain 的名称。TARGET 都有哪些详见 Netfilter 流量控制系统。下面是各种 TARGET 的类型：

- ACCEPT # 允许流量通过
- REJECT # 拒绝流量通过
  - OPTIONS
    - --reject-with icmp-host-prohibited # 通过 icmp 协议显示给客户机一条消息:主机拒绝(icmp-host-prohibited)
- DROP # 丢弃，不响应，发送方无法判断是被拒绝
- RETURN # 返回调用链
- MARK # 做防火墙标记
- 用于 nat 表的 target
  - DNAT|SNAT # {目的|源}地址转换
  - REDIRECT # 端口重定向
  - MASQUERADE # 地址伪装类似于 SNAT，但是不用指明要转换的地址，而是自动选择要转换的地址，用于外部地址不固定的情况
- 用于 raw 表的 target
  - NOTRACK # raw 表专用的 target，用于对匹配规则进行 notrack(不跟踪)处理
- LOG # 记录日志信息
- 引用自定义链 # 直接使用“-j|-g 自定义链的名称”即可，让基本 5 个 Chain 上匹配成功的数据包继续执行自定义链上的规则。

## EXAMPLE

注意：在使用 iptables 命令的时候，为了防止配置问题导致网络不通，可以先设置一个定时任务来修改 iptables 规则或者 iptables XX && sleep 5 && iptables -P INPUT ACCEPT && iptables -P OUTPUT ACCEPT

### Fileter 表的配置

#### INPUT 默认为 DROP 情况下

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

#### INPUT 默认为 ACCEPT 情况下

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

#### 其他

显示 INPUT 链中的所有规则，并显示规则行号

- iptables -L INPUT --line-numbers

删除 INPUT 链中的第 11 条规则

- iptables -D INPUT 11

实现我可以 ping 别人，别人不能 ping 我的方法：

- iptables -A INPUT -p icmp --icmp-type 8 -s 0/0 -j DROP # 默认 INPUT 链的策略为 ACCEPT 的时候用
- iptables -A INPUT -p icmp --icmp-type 0 -s 0/0 -j ACCEPT # 默认 INPUT 链的策略为 DROP 的时候用
- iptables -A OUTPUT -p icmp --icmp-type 0 -s LOCALIP -j DROP # 默认 OUTPUT 链的策略为 ACCEPT 的时候用，注意把 Localip 改为本机 IP
- iptables -A OUTPUT -p icmp --icmp-type 8 -s LOCALIP -j ACCEPT # 默认 OUTPUT 链的策略为 DROP 的时候用，注意把 Localip 改为本机 IP

### NAT 表的配置

- 凡是基于 tcp 协议访问本机 80 端口，且目的地址是 110.119.120.1 的数据包。全部把目的地址转变为 192.168.20.2，且目的端口转换为 8080。
  - iptables -t nat -A PREROUTING -d 110.119.120.1 -p tcp --dport 80 -j DNAT --to-destination 192.168.20.2:8080
- 将源地址网段是 192.168.122.0/24 并且目的地址网段不是 192.168.122.0/24 的全部转换成 123.213.1.5 这个地址。常用于公司内使用私网 IP 的设备访问互联网使用
  - iptables -t nat -A POSTROUTING -s 192.168.122.0/24 ! -d 192.168.122.0/24 -j SNAT --to-source 123.213.1.5
- 将所有到达 eth0 网卡 25 端口的流量转发到 2525 端口。也叫端口转发
  - iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 25 -j REDIRECT --to-port 2525

### RAW 表配置

- 所有来自 10.0.9.0/24 网段的数据包，都不跟踪。
  - iptables -t raw -A PREROUTING -s 10.0.9.0/24 -j NOTRACK

## iptables-save - 将 iptables 规则转储到标准输出

该命令输出的内容更容易被人类阅读，可以用重定向把内容保存到文件中

该命令显示出的信息说明
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fadrg5/1616165483878-a0be9032-9c98-4971-922e-592c61529d86.jpeg)
**Syntax(语法)**
iptables-save \[-M,--modprobe modprobe] \[-c] \[-t table]

**EXAMPLE**

1. iptables-save > /etc/sysconfig/iptables.rules

## iptables-restore - 从标准输入恢复 iptables 规则，可以视同重定向通过文件来读取到标准输入

**EXAMPLE**

1. iptables-restore < /etc/sysconfig/iptables.rules

# ipset - IP 集合的管理工具

ipset 是 iptables 的一个协助工具。可以通过 ipset 设置一组 IP 并对这一组 IP 统一命名，然后在 iptables 里的匹配规则里通过名字直接引用该组 IP，并对这组 IP 进行流量控制。注意：由于是 iptables 规则引用，所以我直接修改 ipset 集合里的 IP，并不用重启 iptables 服务，就可以直接生效。这类似于域名解析，我的机器指定访问 baidu.com，至于百度公司他们怎么更改 IP 与域名绑定的关系，作为用户都不用更更改 baidu.com 这个域名。

## ipset \[OPTIONS] COMMAND \[COMMAND-OPTIONS]

COMMANDS：Note：ENTRY 指的就是 ip 地址

1. create SETNAME TYPENAME \[type-specific-options] # 创建一个新的集合。Create a new set
2. add SETNAME ENTRY # 向指定集合中添加条目。i.e.添加 ip。Add entry to the named set
3. del SETNAME ENTRY # 从指定集合中删除条目 Delete entry from the named set
4. test SETNAME ENTRY # 测试指定集合中是否包含该条目 Test entry in the named set
5. destroy \[SETNAME] # 摧毁全部或者指定的集合 Destroy a named set or all sets
6. list \[SETNAME] # 列出全部或者指定集合中的条目 List the entries of a named set or all sets
7. save \[SETNAME] # 将指定的集合或者所有集合保存到标准输出
8. restore # 还原保存的 ipset 信息
9. flush \[SETNAME] # 删除全部或者指定集合中的所有条目 Flush a named set or all sets
10. rename FROM-SETNAME TO-SETNAME # Rename two sets
11. swap FROM-SETNAME TO-SETNAME # 交换两个集合中的内容 Swap the contect of two existing sets

OPTIONS

1. **-exist** # 在 create 已经存在的 ipset、add 已经存在的 entry、del 不存在的 entry 时忽略错误。
2. **-f** # 在使用 save 或者 restore 命令时，可以指定文件，而不是从标准输出来保存或者还原 ipset 信息

{ -exist | -output { plain | save | xml } | -quiet | -resolve | -sorted | -name | -terse | -file filename }

EXAMPLE

1. ipset list # 列出 ipset 所设置的所有 IP 集合
2. ipset create lichenhao hash:net # 创建一个 hash:net 类型的名为 lichenhao 的 ipset
3. ipset add lichenhao 1.1.1.0/24 # 将 1.1.1.0/24 网段添加到名为 lichenhao 的 ipset 中
4. ipset flush # 清空所有 ipset 下的 ip
5. ipset restore -f /etc/sysconfig/ipset # 从/etc/sysconfig/ipset 还原 ipset 的集合和条目信息

9、屏蔽 HTTP 服务 Flood×××

有时会有用户在某个服务，例如 HTTP 80 上发起大量连接请求，此时我们可以启用如下规则：

iptables -A INPUT -p tcp --dport 80 -m limit --limit 100/minute --limit-burst 200 -j ACCEPT

上述命令会将连接限制到每分钟 100 个，上限设定为 200。

11、允许访问回环网卡

环回访问（127.0.0.1）是比较重要的，建议大家都开放：

iptables -A INPUT -i lo -j ACCEPT

iptables -A OUTPUT -o lo -j ACCEPT

12、屏蔽指定 MAC 地址

使用如下规则可以屏蔽指定的 MAC 地址：

iptables -A INPUT -m mac --mac-source 00:00:00:00:00:00 -j DROP

13、限制并发连接数

如果你不希望来自特定端口的过多并发连接，可以使用如下规则：

iptables -A INPUT -p tcp --syn --dport 22 -m connlimit --connlimit-above 3 -j REJECT

以上规则限制每客户端不超过 3 个连接。

17、允许建立相关连接

随着网络流量的进出分离，要允许建立传入相关连接，可以使用如下规则：

iptables -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT

允许建立传出相关连接的规则：

iptables -A OUTPUT -m conntrack --ctstate ESTABLISHED -j ACCEPT

18、丢弃无效数据包

很多网络 ××× 都会尝试用 ××× 自定义的非法数据包进行尝试，我们可以使用如下命令来丢弃无效数据包：

iptables -A INPUT -m conntrack --ctstate INVALID -j DROP

19、IPtables 屏蔽邮件发送规则

如果你的系统不会用于邮件发送，我们可以在规则中屏蔽 SMTP 传出端口：

iptables -A OUTPUT -p tcp --dports 25,465,587 -j REJECT
