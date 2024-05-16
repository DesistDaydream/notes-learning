---
title: ss
linkTitle: ss
date: 2024-05-16T21:40
weight: 20
---

# 概述

> 参考：
>
> - [Manual(手册)，ss(8)](https://man7.org/linux/man-pages/man8/ss.8.html)

ss 用于转储 socket 统计信息。 它允许显示类似于 netstat 工具的信息。 它可以显示比其他工具更多的 TCP 和状态信息

# Syntax(语法)

**ss \[OPTIONS] \[FILTER]**

## OPTIONS

如果不使用任何选项，则 ss 命令将显示已建立连接的打开的非监听套接字（例如 TCP / UNIX / UDP）的列表。

```bash
~]# ss
Netid            State            Recv-Q            Send-Q                                                               Local Address:Port                           Peer Address:Port
......
u_str            ESTAB            0                 0                                                                                * 22170                                     * 22171
u_str            ESTAB            0                 0                                                      /run/systemd/journal/stdout 23554                                     * 23347
u_str            ESTAB            0                 0                            /var/lib/sss/pipes/private/sbus-dp_implicit_files.850 22246                                     * 22240
tcp              ESTAB            0                 0                                                                    172.19.42.248:ssh                           172.19.42.203:63482
```

- **-A QUERY, --query=QUERY, --socket=QUERY** # 要转储的套接字表的列表，用逗号分隔。可以理解以下标识符：all，inet，tcp，udp，raw，unix，packet，netlink，unix_dgram，unix_stream，unix_seqpacket，packet_raw，packet_dgram，dccp，sctp，vsock_stream，vsock_dgram，xdp 列表中的任何项目都可以选择添加前缀带有感叹号`!`，以防止该套接字表被转储。
- **-a, --all** # 显示所有已监听和未监听的 Sockets。
- **-e, --extended** # 显示详细的 Socket 信息。对于 TCP 连接来说，相当于 -o 选项。
  - ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/olkcwd/1626247848384-bfa595e0-6f72-49f3-b73c-665c07b5dfa1.png)
  - timer # TCP 连接的 keepalive 计时器
  - uid # 该 Socket 所属的用户 ID
  - ino # 该 Socket 的 inode 编号
  - sk # 该 Socket 的 UUID
- **-n, --numeric** # 直接使用 ip 地址，而不通过域名服务器(常用，节省资源，界面明确）
- **-o, --options** # 显示定时器信息。格式为：`timer:(timer_name>,<expire_time>,<retrans)`
  - timer_name # 计时器的名称。一共有 5 种：
    - on # 以下三种计时器之一: TCP retrans timer,TCP early retrans timer and tail loss probe timer
    - keepalive # TCP 的保持活动计时器。Linux Kernel 默认 7200 秒。
    - timewait # 当 TCP 连接进入 TIME-WAIT 状态时，将会触发该计时器。
    - persist # zero window probe timer
    - unknown # 未知计时器
  - expire_time # 计时器的过期时间
  - retrans # 重新传输的次数。即 TCP keepalive 探测是重试的次数。
- **-s, --summary**# 输出摘要统计信息. 此选项不解析套接字列表，以从各种来源获取摘要。 常用于套接字数量太大以至于无法解析 /proc/net/tcp 时。输出信息如下：

```bash
~]# ss -s
Total: 170
TCP:   5 (estab 1, closed 0, orphaned 0, timewait 0)
Transport Total     IP        IPv6
RAW    0         0         0
UDP    5         3         2
TCP    5         3         2
INET   10        6         4
FRAG   0         0         0
```

## FILTER

**FILTER = \[ state STATE-FILTER ] \[ EXPRESSION ]**
**STATE-FILTER(状态过滤)** # 指定要过滤的连接状态，可用的状态有 established, syn-sent, syn-recv, fin-wait-1, fin-wait-2, time-wait, closed, close-wait, last-ack, listen 和 closing.

- 除了以上状态，还可以用其他描述方式代替一个或多个状态
- **all** # for all the states
- **connected** # all the states except for **listening** and **closed**
- **synchronized** # **syn-sent** 状态之外的所有**已经连接**的状态
- **bucket** # states, which are maintained as minisockets, i.e. **time-wait** and **syn-recv**
- **big** # opposite to **bucket**

**EXPRESSION(表达式)**#

### IP 地址筛选

ss src ADDRESS_PATTERN

- src：表示来源
- ADDRESS_PATTERN：表示地址规则

如下：

```bash
ss src 120.33.31.1
# 列出来之20.33.31.1的连接
```

```bash
＃ 列出来至120.33.31.1,80端口的连接
ss src 120.33.31.1:http
ss src 120.33.31.1:8
```

### 端口筛选

ss dport OP PORT

- OP # 运算符
- PORT # 表示端口
- dport # 表示过滤目标端口、相反的有 sport

OP 运算符如下：

```bash
<= 或 le : 小于等于 >= or ge : 大于等于
== 或 eq : 等于
!= 或 ne : 不等于端口
< 或 lt : 小于这个端口 > or gt : 大于端口
```

OP 实例

```bash
ss sport = :http 也可以是 ss sport = :80
ss dport = :http
ss dport \> :1024
ss sport \> :1024
ss sport \< :32000
ss sport eq :22
ss dport != :22
ss state connected sport = :http
ss \( sport = :http or sport = :https \)
ss -o state fin-wait-1 \( sport = :http or sport = :https \) dst 192.168.1/24
```

# EXAMPLE

显示状态是 established 的所有 tcp 连接

`ss -nta state established`

显示所有已建立的 SMTP 连接

`ss -o state established '( dport = :smtp or sport = :smtp )'`

显示所有已建立的 HTTP 连接

`ss -o state established '( dport = :http or sport = :http )'`

为我们的网络 193.233.7/24 列出所有处于状态 FIN-WAIT-1 的 tcp socket，并查看它们的计时器。

`ss -o state fin-wait-1 '( sport = :http or sport = :https )' dst 193.233.7/24`

列出所有 socket(TCP 除外) 中所有状态的套接字。

`ss -a -A 'all,!tcp'`

# netstat 命令行工具

与 ss 命令类似，OPTIONS 都差不多，就是显示的信息格式还有内容有细微差别

OPTIONS

- -a 或--all：显示所有状态的连线中的端口
- -p 或--programs：显示正在使用端口的程序识别码和程序名称；
- -t 或--tcp：显示 TCP 传输协议的连线状况
- -u 或--udp：显示 UDP 传输协议的连线状况
- -x 或--unix：显示 Unix 相关选项
- -e 或--extend：显示网络其他相关信息；
- -l 或--listening：显示监听中的服务端口
- -r 或--route：显示 Routing Table；（与 route 命令显示信息相同）
- -s 或--statistice：显示网络工作信息统计表；
- -c 或--continuous：持续列出网络状态；

输出列表内容解析：分为两部分

第一部分：Active Internet connections 活跃的互联网连接

Proto 显示 socket 使用的协议(tcp,udp,raw)

Recv-Q 和 Send-Q 指的是接收队列和发送队列(这些数字一般都应该是 0,如果不是则表示软件包正在队列

中堆)

Local Address 显示在本地哪个地址和端口上监听

Foreign Address 显示接收外部哪些地址哪个端口的请求

State 显示 socket 的状态(通常只有 tcp 有状态信息)

PID/Program name 显示 socket 进程 id 和进程名

第二部分：Active UNIX domain sockets 活跃的 UNIX 域套接口，只能用于本机通信的，性能可以提高一倍

RefCnt 表示连接到本套接口上的进程号

Types 显示套接口的类型

Path 表示连接到套接口的其它进程使用的路径名

常用例子：

不进行域名解析的查看所有状态的端口协议为 tcp 且显示该端口的进程号和进程名的相关网络信息

`netstat -atnp`

- netstat -n | awk '/^tcp/ {print $5}' | awk -F: '{print $1}'| sort | uniq -c | sort -rn # 统计当前系统中所有 tcp 和 udp 的连接，根据状态统计
- netstat -n | awk '/^tcp/ {++S\[$NF]} END {for(a in S) print a, S\[a]}' # 统计当前系统中所有 tcp 和 udp 的连接，根据状态统计
