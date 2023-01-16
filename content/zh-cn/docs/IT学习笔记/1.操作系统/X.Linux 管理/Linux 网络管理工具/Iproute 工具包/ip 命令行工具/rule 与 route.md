---
title: rule 与 route
---

# 概述

> 参考：
> - [Manual(手册),ip-route(8)](https://man7.org/linux/man-pages/man8/ip-route.8.html)
> - [Manual(手册),ip-rule(8)](https://man7.org/linux/man-pages/man8/ip-rule.8.html)

# route # 路由条目管理

route 可以操作内核路由表中的**条目**。直接使用命令可以列出 main 路由表中的条目：

```bash
root@lichenhao:~# ip route
default via 172.19.42.1 dev ens3 proto static
172.19.42.0/24 dev ens3 proto kernel scope link src 172.19.42.248
```

> 注意：如果默认网关已由 DHCP 分配，并且配置文件中指定了具有相同度量的同一网关，则在启动或启动接口时将发生错误。可能会显示以下错误消息：`RTNETLINK answers:File exists`。可以忽略此错误。

| 目的地址 | via
下一跳 | dev
网络设备 | proto
生成路由条目的方式 | scope
覆盖范围 | src
源地址 |
| --- | --- | --- | --- | --- | --- |
| default | 172.19.42.1 | ens3 | static | | |
| 172.19.42.0/24 | | ens3 | kernel | link | 172.19.42.248 |

**Route Type(路由类型)**

- unicast
- unreachable
- blackhole
- prohibit
- local
- broadcast
- throw
- nat
- anycast
- multicast

**Route Tables(路由表)**
Linux-2.x 版本内核以后，可以根据 **SELECTOR(选择器)** 将数据包交给不同的路由表进行路由处理。这些路由表由 1 到 232 范围内的数字表示(/etc/iproute2/rt_tables 文件中可以为数字起一个别名)。默认情况下，所有普通路由规则都会插入名为 main 的路由表中(main 路由表的 ID 为 254)。ip rule 命令可以对路由表进行控制。

## Syntax(语法)

**ip \[Global OPTIONS] route \[COMMAND]**
**ip \[Global OPTIONS] route { show | flush } \[to] SELECTOR**
**ip \[Global OPTIONS] route get ROUTE_GET_FLAGS ADDRESS \[ from ADDRESS iif STRING ] \[ oif STRING ] \[ mark MARK ] \[ tos TOS ] \[ vrf NAME ] \[ ipproto PROTOCOL ] \[ sport NUMBER ] \[ dport NUMBER ]**
**ip \[Global OPTIONS] route {add|del|change|append|replace} ROUTE**

### ARGUMENTS

**SELECTOR = \[ root PREFIX ] \[ match PREFIX ] \[ exact PREFIX ] \[ table TABLE_ID ] \[ vrf NAME ] \[ proto RTPROTO ] \[ type TYPE ] \[ scope SCOPE ]**

- **table \<TABLE_ID> **# 指定路由表的标识符

**ROUTE = NODE_SPEC \[ INFO_SPEC ]**

- NODE_SPEC = \[ TYPE ] PREFIX \[ tos TOS ] \[ table TABLE_ID ] \[ proto RTPROTO ] \[ scope SCOPE ] \[ metric METRIC ] \[ ttl-propagate { enabled | disabled } ]
- INFO_SPEC = { NH | nhid ID } OPTIONS FLAGS \[ nexthop NH ] ...
  - NH := \[ encap ENCAP ] \[ via \[ FAMILY ] ADDRESS ] \[ dev STRING ] \[weight NUMBER ] NHFLAGS

**FAMILY := \[ inet | inet6 | mpls | bridge | link ]**

**RTPROTO = \[ STRING | NUMBER ]** # 生成本路由条目时所使用的 Routing Protocol(路由协议，简称 RTPROTO) 的标识符。RTPROTO 的值来自于 /etc/iproute2/rt_protos 文件中数字或字符串。

> 注意：此协议不是指传统意义上的 http、tcp 这种协议，而是指，生成路由条目的方式、或者说生成路由条目的实体。比如我们可以说内核自己生成了一个路由条目；也可以说通过 dhcp 获取 IP 时生成了路由条目；还可以说通过人为手动创建了一个路由条目；等等。

- boot # 默认 RTPROTO。该路由条目在 bootup sequence 期间生成的。且路由守护进程启动时，这些条目将被删除
  - 不太理解这个官方的解释？？用人话说，就是 ip 命令默认添加的路由条目在机器重启后会被删除。但是这个类型兜底是啥子意思哦？~
- kernel # 该路由条目在内核自动配置期间生成。
- dhcp #
- static # 该路由条目由管理员手动添加以覆盖动态路由。路由守护进程会尊重它们，甚至可能将它们通告给它的 peers。
- ra # 该路由条目由 Router Discovery Protocol(路由发现协议) 生成。通常只出现在 IPv6 中
- 等等

**TABLE_ID := \[ local| main | default | all | NUMBER ]** #

**SCOPE := \[ STRING | NUMBER ]** # 目的地址覆盖的范围。即路由数据包之前，从哪些地方找目的地址。SCOPE 的值来自于 /etc/iproute2/rt_scopes 文件中的数字或字符串。如果省略此参数，则 ip 程序默认 unicast(单播) 类型的路由范围为 global、local 类型的路由范围为 host、unicast 和 broadcst 类型的路由范围为 link。
用人话说：为数据包选择路由条目前，还需要判断目的地址的有效性。也就是说，目的地址在哪里才是可以被路由的。

- host # 目的地址仅在本主机上有效
- link # 目的地址仅在本网络设备上有效
- global # 目的地址全局有效

## ip route show

显示路由表的内容或按某些标准选择的路由。

- **to <SELECTOR>** #
- **protocol <RTPROTO>** # 显示指定协议标识符的路由条目

## EXAMPLE

- 查看名为 local 的路由表的条目
  - ip route show table local
- 查看该 IP 从哪里过来
  - ip route get 192.168.0.1
- 添加默认路由条目，经过 ens3 网络设备，下一跳是 172.19.42.1
  - ip route add default via 172.19.42.1 dev ens3
- 添加路由条目，目的地址是 10.10.10.0/24 网段的数据包的下一跳地址是 192.168.0.2 使用 eth0 网络设备
  - ip route add 10.10.10.0/24 via 192.168.0.2 dev eth0

# rule # 路由策略数据库管理

rule 可以操作路由策略数据库中的规则，控制路由选择算法。说白了就是可以**控制路由表，**而 ip route 则是控制**路由表的条目。**

在互联网上，传统的路由算法仅基于数据包的目标地址做出路由选择。但是在某些情况下，我们希望路由数据包的策略，而不仅仅取决于目标地址，还可以通过源地址、IP 协议、传输协议、端口、甚至数据包的 payload 等等信息来对数据包的路由进行选择。这种方式，称为 **Policy Routing(策略路由)**。

为了解决上面的问题，传统的基于目的地址的 **Routing table(路由表)** 被替换为 **Routing policy database(路由策略数据库，简称 RPDB)**。RPDB 通过执行一些 **Rule(规则) **来选择路由表。

每个路由策略规则由以下两部分组成：

- **Selector(选择器)** # 通过一些规则，对数据包进行匹配，匹配到的数据包，将会执行 Action 定义的动作。
- **Action(动作) **# 匹配到的数据包将要执行的动作。
  - 比如有一个动作叫 lookup，用来指定要查找路由条目的路由表 ID。意思就是指，根据指定路由表中的路由条目，来决定 Selector 匹配到的数据包应该被路由到哪里

RPDB 按照优先级递减的顺序注意扫描这些规则(数字越小，优先级越高)。

在启动时，内核将会配置三个规则组成默认的 RPDB 条目：

```bash
root@lichenhao:~# ip rule
0:	from all lookup local
32766:	from all lookup main
32767:	from all lookup default
```

- 0: from all lookup local
  - **local 路由表(ID 255) # **是包含用于本地和广播地址的高优先级控制路由的特殊路由表。
  - Priority(优先级) # 0
  - Selector(选择器) # 匹配所有
  - Action(动作) # 查找名为 local 的路由表。
- 32766: from all lookup main
  - **main 路由表(ID 254)** # 是包含所有非策略路由的正常路由表。可以通过管理员删除和/或覆盖此规则。我们平时配置的路由条目都是在这个表中配置的。
  - Priority(优先级) # 32766
  - Selector(选择器) # 匹配所有
  - Action(动作) # 查找名为 local 的路由表。
- 32767: from all lookup default
  - **default 路由表(ID 253)** # 为空。如果未选择先前的默认规则，则保留某些后处理。也可以删除此规则。
  - Priority(优先级) # 327667
  - Selector(选择器) # 匹配所有
  - Action(动作) # 查找名为 local 的路由表。

每个 RPDB 条目都有附加属性。每个规则都有一个指向某个路由表的指针。 NAT 和伪装规则有一个属性来选择要转换/伪装的新 IP 地址。除此之外，规则还有一些可选属性，路由也有，即领域。这些值不会覆盖路由表中包含的值。它们仅在路由未选择任何属性时使用。

## Syntax(语法)

**ip \[OPTIONS] rule { COMMAND | help }**
**ip \[OPTIONS] rule \[ list \[ SELECTOR ]]**
**ip \[OPTIONS] rule { add | del } SELECTOR ACTION**

**SELECTOR**

- **from PREFIX** # 选择要匹配的源地址。`默认值：all`
- **to PREFIX **# 选择要匹配的目的地址。`默认值：all`
- **priority NUM** # 策略规则的优先级。`默认值:当前数字最大的优先级逐一减 1`
- **ACTION**

- **lookup TABLEID** # 根据 SELECTOR 匹配到的查找路由表，根据指定的 TABLEID 路由表来处理数据包的路由。`默认值：254`

## EXAMPLE

- 添加优先级为 1，ID 为 2 的路由表，所有源地址是 192.168.0.0/24 网段的数据包都根据该路由表的规则进行路由。
  - **ip rule add priority 1  from 192.168.0.0/24 table 2**
