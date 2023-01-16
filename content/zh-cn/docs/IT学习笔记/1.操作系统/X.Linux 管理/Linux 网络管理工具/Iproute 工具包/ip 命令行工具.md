---
title: ip 命令行工具
---

# 概述

> 参考：
> - [Manual(手册),ip(8)](https://man7.org/linux/man-pages/man8/ip.8.html)

ip 命令行工具可以控制各种 **Object(对象)**，这些对象包括：路由、网络设备、接口、隧道 等

# Syntax(语法)

**ip \[Global OPTIONS] OBJECT \[COMMAND]**

ip 程序的语法有点复杂，对 Object 控制的命令中，有非常多的参数，不像普通命令一样，把参数称为 FLAGS 或 OPTIONS，且都是以 `-` 或者 `--` 符号开头的。

这里我们使用 **大写字母 **来描述 **一个参数** 或 **一个具体的值**。参数中还可以包含一个或多个其他参数，每个参数的值，同样使用大写字母表示。

在后面的文章中，凡是这种复杂的参数，都使用这类格式表示：`参数 := 参数 | 值`，这就有点像编程中初始化**变量**一样。在这里就是等于是定义一个参数，并为参数赋值。比如 `ip link` 命令中，就有这种样子的写法：
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/us4bal/1638423450051-14c93955-fbe9-425a-9d96-eaf14b140241.png)
这里面有一个 IFADDR 表示一个参数，IFADDR 参数又是由 PREFIX、SCOPE-ID 等参数组成，而 SCOPE-ID 则表示有具体含义的值。其实，本质上，命令行工具的参数，就是要传入代码的 Function 中的的实际参数。

## Global OPITONS

注意：这里的 OPTIONS 是全局选项，要用在 ip 与 OBJECT 之间，比如：

```bash
root@lichenhao:~# ip -c route
default via 172.19.42.1 dev ens3 proto static metric 100
10.19.0.0/24 dev docker0 proto kernel scope link src 10.19.0.1 linkdown
172.19.42.0/24 dev ens3 proto kernel scope link src 172.19.42.248 metric 100

root@lichenhao:~# ip route -c
Command "-c" is unknown, try "ip route help".
```

可以看到，-c 选项用在 OBJECT 后面是无效的。

- **-d, -details** # 输出更详细的信息,主要显示该网络设备的类型等
- **-f, -family <FAMILY>** # 指定要使用的协议族。协议族标识符可以是 inet、inet6、bridge、mpls、link 中的一种。如果不存在此选项，则从其他参数中猜测协议族。如果命令行的其余部分没有提供足够的信息来猜测该系列，则 ip 会退回到默认值，通常是 inet 或任何其他值。链接是一个特殊的系列标识符，表示不涉及网络协议。
  - **-4** # `-family inet` 的简写
  - **-6** # `-family inet6` 的简写
  - **-B** # `-family bridge` 的简写
  - **-M** # `-family mpls` 的简写
  - **-0** # `-family link` 的简写
- **-o, -oneline **# 在一行中输出每条记录，并用''字符替换换行符。在使用 wc(1) 对记录进行计数 或 对输出进行 grep(1) 时，这非常方便。
  - 注意，使用 -o 选项时，不会打印没有 IP 地址的网络设备
- **-s, -stats** # 显示更详细的信息,主要显示该网络设备的接收、发送、错误、丢弃的数据包信息

## ARGUMENTS

**OBJECT := { link | address | addrlabel | route | rule | neigh | ntable | tunnel | tuntap | maddress | mroute | mrule | monitor | xfrm | netns | l2tp | tcp_metrics | token | macsec | vrf | mptcp }**

- 注意：OBJECT(对象)选项用来指定 ip 程序想要控制的网络栈中的实体。
- 比如：
  - link(链路)
  - address(地址)
  - route(路由条目)
  - 等
- ip 命令可以对这些网络对象进行相关操作，选定 object 后，后面执行相关 command 进行操作

## COMMAND

主要都是对各个 Object 的 add | delete | set | show | list 等类似增删改查的命令，还有各个 Object 独自的子命令

# OBJECT 命令详解

## link # 网络设备配置

详见：
[link](✏IT 学习笔记/📄1.操作系统/X.Linux%20 管理/Linux%20 网络管理工具/Iproute%20 工具包/ip%20 命令行工具/link.md 命令行工具/link.md)

## address # IPv4 或 IPv6 地址管理

### Syntax(语法)

**ip \[OPTIONS] address \[COMMAND]**

**COMMAND：**

- add | del | change | replace | show
- save | flush
- showdump | restore

### EXAMPLE

- 简略显示网络设备信息
  - ip -4 -o a s

<!---->

    root@lichenhao:~/projects/kubeappsops# ip -4 -o a s
    1: lo    inet 127.0.0.1/8 scope host lo\       valid_lft forever preferred_lft forever
    2: ens3    inet 172.19.42.248/24 brd 172.19.42.255 scope global ens3\       valid_lft forever preferred_lft forever
    3: docker0    inet 172.17.0.1/16 brd 172.17.255.255 scope global docker0\       valid_lft forever preferred_lft forever

- 筛选满足 IP 地址格式的网卡信息
  - ip a s up | egrep --color=auto -n '\[0-9]+.\[0-9]+.\[0-9]+.\[0-9]+'
- 显示\[已经启动的]网卡 ip 信息,类似于 ifconfig 命令,可简写为 ip a s up
  - ip address show \[up]
- 以详细信息显示 ens33 的网卡关于地址的信息，包括收发包的状态等
  - ip -s addr show ens33
- 给 eth0 网卡添加一个临时的 IP 地址
  - ip addr add 192.168.0.1/24 dev eth0
- ip addr add 10.0.0.101/24 broadcast 10.0.0.255 dev eth0 label eth0:0

Note：在 ip address show 命令中列出的网络设备，可能包含这样的名称 eth0.2@eth0 。其实设备名就是 eth0.2(其中 2 表示 vlan 号)，至于后面的`@`则是一个关联同等级网络设备的符号，表示 eth0.2 这个设备是与 eth0 关联的。如果是 bridge 和 bond 之类的网络设备，则没有`@`符号，因为和 bridge 或者 bond 关联的设备都是属于下级设备。

## route # 路由条目管理

详见：&#x20;
[rule 与 route](✏IT 学习笔记/📄1.操作系统/X.Linux%20 管理/Linux%20 网络管理工具/Iproute%20 工具包/ip%20 命令行工具/rule%20 与%20route.md 命令行工具/rule 与 route.md)

## rule # 路由策略数据库管理

详见：
[rule 与 route](✏IT 学习笔记/📄1.操作系统/X.Linux%20 管理/Linux%20 网络管理工具/Iproute%20 工具包/ip%20 命令行工具/rule%20 与%20route.md 命令行工具/rule 与 route.md)

## neighbor #管理 ARP 或 NDISC 缓存条目

**ip \[OPTIONS] neighbor \[COMMAND]**

EXAMPLE

1. ip neighbor list # 显示邻居和 arp 表，即学到的 IP 地址，可显示该 IP 是否可达等状态，以及是从哪个端口学到的
2. ip neigh flush dev eth0 # 移除 eth0 设备上的邻居条目（i.e.清空 arp）

## tuntap # tun/tap 网络设备的管理

**ip tuntap COMMAND mode { tun | tap } \[ dev PHYS_DEV ] \[ user USER ] \[ group GROUP ] \[ one_queue ] \[ pi ] \[ vnet_hdr ] \[ multi_queue ] \[ name NAME ]**

EXAMPLE

1. ip tuntap add dev vnet3 mode tun # 创建一个 tun 类型，名为 vnet3 的网络设备

## netns # 进程网络命名空间管理

**ip \[OPTIONS] netns \[COMMAND]**
是在 linux 中提供网络虚拟化的一个项目，使用 netns 网络空间虚拟化可以在本地虚拟化出多个网络环境，目前 netns 在 lxc 容器中被用来为容器提供网络(注意:可以查看 openstack、docker 各个节点上的虚拟网络设备并进行操作)。

COMMAND
add、delete、set、list 增删改查通用命令

EXAMPLE

- ip netns add r1 #创建名字为 r1 的 namespace
- ip netns list #列出 net namespace，显示名称

identify

pids

exec

**ip netns exec NAME COMMAND.... **# 对 NAME 这个 namesapce 执行 COMMAND 命令

EXAMPLE

- ip netns exec r1 bash #进入 r1 这个 namesapce 的空间中，相当于启动了一个单独的关于该 namespace 的 shell，可以使用 exit 退出
- ip netns exec r1 ip a s #显示 r1 这个 namespace 的网路信息
- ip netns exec r1 ip link set veth1.1 name eth0 #设定 r1 这个 namespace 中的 veth1.1 网卡的名字为 eth0

monitor
