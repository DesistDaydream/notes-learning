---
title: ipvsadm 命令行工具
---

# 概述

# Syntax(语法)

处理 Virtual Service 的命令语法：

- **ipvsadm COMMAND VirtualService \[-s Scheduler] \[Persistence OPTIONS]**

处理指定 Virtual Service 下的 Real Server 的命令语法：

- **ipvsadm COMMAND VirtualService ServerAddress \[PacketForwardingMethod] \[Weight Options]**

\*\*
命令语法中各参数的含义

- **VirtuslService** # 用于指定基于协议或者地址或者端口号的虚拟服务，通过三元组定义：Protocol、IP、PORT。
  - 格式：`-PROTOCOL IP:PORT`
    - -PROTOCOL 分两种
      - -t, --tcp-service
      - -u, --udp-service
- **Scheduler** # Director 的调度方法
  - 详见[ LVS 介绍](https://www.yuque.com/go/doc/33184072) 文章中描述的调度方法，使用其中 10 种任意一种的英文简称来写该参数，注意：是小写字母
- **ServerAddress** # RS 的 IP
- **PacketForwardingMethod** # 该位置指明 LVS 的工作模式，不写该参数表明默认 DR 类型
  - -g, --gatewaying # 网关，表示 DR 模式
  - -i, --ipip # IP 封装 IP，表示 TUN 模式
  - -m, --masquerading # 伪装，表示 NAT 模式
- **Weight Options** # 权重选项

## COMMAND

管理集群服务的 COMMAND

- ** -A, --add-service** # 创建一个 VirtualService . 服务地址由 IP，PORT，protocol(协议)组成
- **-E, --edit-service** # 修改一个 VirtualService
- **-D, --delete-service** # 删除一个虚拟服务,以及相关的 RS
- **-C, --clear** # 清空这个虚拟服务器的表
- **-R, --restore** # 从标准输出还原 ipvs 规则
- **-S, --save **# 保存 ipvs 规则到标准输出

管理集群服务中的 RS 的 COMMAND

- **-a, --add-server** # 向指定的 VirtualService 中添加一个 RS
- **-e, --edit-server** # 修改指定 VirtualService 中的 RS
- **-d, --delete-server** # 从指定的 VirtualService 中移除一个 RS

通用的 COMMAND

- **-Z, --zero** # 清空一个或所有 Virtual Service 下的数据包
- **-L, -l, --list** # 如果没有指定参数，则列出虚拟服务器表。如果输入了 service-address，只列出该服务。
- **--set TCP TCPFIN UDP **# 更改用于 IPVS 连接的超时值。 此命令始终使用 3 个参数，分别表示 TCP 会话，接收到 FIN 数据包后的 TCP 会话和 UDP 数据包的超时值（以秒为单位）。 超时值 0 表示保留了相应条目的当前超时值。`默认值：900 120 300`。

## OPTIONS

- **-p, --persistent \[TIME]** # 指定持久连接的超时时间

### 与 -L, -l, --list 命令搭配使用的选项

- **-n** # 当配合 -l 查询的时候，显示 IP 地址而不显示主机名，即不把 IP 解析成主机名。与 -l 命令配合时，效果如下：

```bash
[root@dr-01 ~]# ipvsadm -Ln
IP Virtual Server version 1.2.1 (size=4194304)
Prot LocalAddress:Port Scheduler Flags
  -> RemoteAddress:Port           Forward Weight ActiveConn InActConn
TCP  10.10.9.60:80 rr
  -> 10.10.9.61:80                Route   1      0          0
  -> 10.10.9.62:80                Route   1      0          0
TCP  10.10.9.60:30000 rr persistent 30
  -> 10.10.9.69:30000             Route   1      0          0
  -> 10.10.9.70:30000             Route   1      0          0
```

- Forward # 当前 lvs 的模型，其中 Masq 为 nat 模型，route 为 dr 模型
- Weight # 当前 rs 的权重
- ActiveConn # 活动连接数,也就是 tcp 连接状态的 ESTABLISHED;
- InActConn # 指除了 ESTABLISHED 以外的,所有的其它状态的 tcp 连接.
- **-c, --connection** # 输出 ipvs 当 前的连接状态信息。效果如下：

<!---->

    [root@node-1 ~]# ipvsadm --list -nc
    IPVS connection entries
    pro expire state       source             virtual            destination
    TCP 14:52  ESTABLISHED 10.244.3.238:60996 10.96.0.1:443      172.38.40.214:6443
    UDP 03:39  UDP         10.244.3.244:39398 10.96.0.10:53      10.244.0.39:53
    TCP 14:48  ESTABLISHED 10.244.3.238:58412 10.96.0.1:443      172.38.40.214:6443

- **--timeout** # 输出 TCP 会话，接收到 FIN 数据包后的 TCP 会话和 UDP 数据包的超时值（以秒为单位）。效果如下：

<!---->

    [root@node-1 ~]# ipvsadm -ln --timeout
    Timeout (tcp tcpfin udp): 900 120 300

- **--daemon **# Daemon information output. The list command with this option will display the daemon status and its multicast interface.
- **--stats** # 显示统计数据，效果如下

<!---->

    [root@node-1 ~]# ipvsadm -ln --stats
    IP Virtual Server version 1.2.1 (size=4096)
    Prot LocalAddress:Port               Conns   InPkts  OutPkts  InBytes OutBytes
      -> RemoteAddress:Port
    TCP  172.38.40.216:30080                 0        0        0        0        0
      -> 10.244.5.86:80                      0        0        0        0        0
    TCP  10.96.0.1:443                      66     4327     4056   500036  4558037
      -> 172.38.40.212:6443                 22      835      702   147508   743066
      -> 172.38.40.213:6443                 22     2485     2432   240938  2923352
      -> 172.38.40.214:6443                 22     1007      922   111590   891619

- **--rate** # 速率信息输出。显示服务及其服务器的速率信息（例如，连接/秒，字节/秒和数据包/秒）。
- **--thresholds** # 输出阈值信息。显示服务列表中每个服务器的上/下连接阈值信息。
- **--persistent-conn** # 持久连接信息的输出。在服务列表中显示每个服务器的持久连接计数器信息。持久连接用于将实际连接从同一客户端/网络转发到同一服务器。

# EXAMPLE

- 管理集群服务
  - ipvsadm -A|E virtual-service \[-s SCHEDULER] #增加修改
    - ipvsadm -A -t 192.168.0.63:80 -s rr #添加一个虚拟服务，调度模式为轮询
    - ipvsadm -D virtual-service #删除
- 管理集群服务中的 RS
  - ipvsadm -a|e virtual-service -r server-address \[-g|i|m] \[-w weight] \[-x upper] \[-y lower]
  - ipvsadm -d virtual-service -r server-address
    - ipvsadm -a -t 192.168.0.60:80 -r 192.168.0.62 -g #添加一个 IP 为 0.62 的 RS 到 0.60 的 LVS 中，LVS 类型为-g,dr 类型
- 通用
  - ipvsadm -C # 清空
  - ipvsadm -L|l \[virtual-service] \[options] # 查询
  - ipvsadm -R
  - ipvsadm -S \[-n]
- ipvsadm -Ln # 查询，直接显示 IP 不显示主机名，信息如下所示

那既然这样,为什么从 lvs 里看的 ActiveConn 会比在真实机上通过 netstats 看到的 ESTABLISHED 高很多呢?

原来 lvs 自身也有一个默认超时时间.可以用 ipvsadm -L --timeout 查看,默认是 900 120 300,分别是 TCP TCPFIN UDP 的时间.也就是说一条 tcp 的连接经过 lvs 后,lvs 会把这台记录保存 15 分钟,而不管这条连接是不是已经失效(哪怕这次 http 请求已经完成切断开连接，但是在 lvs 中有保存记录)!所以如果你的服务器在 15 分钟以内有大量的并发请求连进来的时候,你就会看到这个数值直线上升.

其实很多时候,我们看 lvs 的这个连接数是想知道现在的每台机器的真实连接数吧?怎么样做到这一点呢?其实知道现在的 ActiveConn 是怎样产生的,做到这一点就简单了.举个例子:比如你的 lvs 是用来负载网站,用的模式是 dr,后台的 web server 用的 nginx 这时候一条请求过来,在程序没有问题的情况下,一条连接最多也就五秒就断开了.这时候你可以这样设置 ipvsadm --set 5 10 300.设置 tcp 连接只保持 5 秒中.如果现在 ActiveConn 很高你会发现这个数值会很快降下来,直到降到和你用 nginx 的 status 看当前连接数的时候差不多.你可以继续增加或者减小 5 这个数值,直到真实机的 status 连接数和 lvs 里的 ActiveConn，一般保证该数值与后端 RS 中代理服务的连接保持时间相当即可
