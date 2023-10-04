---
title: net(网络相关参数)
weight: 1
---

# 概述

> 参考：
>
> - [Manual(手册)，/proc/sys 部分](https://man7.org/linux/man-pages/man5/proc.5.html#DESCRIPTION)
> - [官方文档，Linux 内核用户和管理员指南-/proc/sys 文档-/proc/sys/net 文档](https://www.kernel.org/doc/html/latest/admin-guide/sysctl/net.html)
> - [官方文档，内核子系统文档-Networking-IP Sysctl](https://www.kernel.org/doc/html/latest/networking/ip-sysctl.html)
>   - 这里包含 net 下的 ipv4、ipv6、bridge 等相关参数，但是没有 netfilter 等的相关参数

**/proc/sys/net/** 目录下通常包含下面的一个或多个子目录

| 目录名    | 用处               |     | 目录名    | 用处                |
| --------- | ------------------ | --- | --------- | ------------------- |
| 802       | E802 protocol      |     | mptcp     | Multipath TCP       |
| appletalk | Appletalk protocol |     | netfilter | Network Filter      |
| ax25      | AX25               |     | netrom    | NET/ROM             |
| bridge    | Bridging           |     | rose      | X.25 PLP layer      |
| core      | General parameter  |     | tipc      | TIPC                |
| ethernet  | Ethernet protocol  |     | unix      | Unix domain sockets |
| ipv4      | IP version 4       |     | x25       | X.25 protocol       |
| ipv6      | IP version 6       |     |           |                     |

# ipv4 参数

> 参考：
>
> - [官方文档，内核子系统文档-Networking-IP Sysctl-/proc/sys/net/ipv4/* 变量](https://www.kernel.org/doc/html/latest/networking/ip-sysctl.html#proc-sys-net-ipv4-variables)

## net.ipv4.ip_forward(0 | 非 0)

在多个网络设备间是否允许转发数据包。`默认值：0`

这个变量很特殊，它的改变会将所有配置参数重置为默认状态([RFC1122](https://www.rfc-editor.org/rfc/rfc1122.html) 用于主机，[RFC1812](https://www.rfc-editor.org/rfc/rfc1812.html) 用于路由器)

## TCP 相关参数

### net.core.somaxconn(INTEGER)

服务端所能 accept 即处理数据的最大客户端数量，即完成连接上限，`默认值：128`。

推荐 65535

### net.ipv4.tcp_fin_timeout(INTEGER)

表示如果 socket 由本端要求关闭，值为整数，单位为秒。这个参数决定了它保持在 FIN-WAIT-2 状态的时间。

此参数确定孤立(未引用)连接在本地端终止之前等待的时间长度。当发生在远程对等端上的事件阻止或过度延迟响应时，此参数尤其有用。由于用于连接的每个套接字大约消耗 1.5K 字节的内存，内核必须主动中止和清除死的或陈旧的资源。

推荐 5

### net.ipv4.tcp_max_syn_backlog(INTEGER)

所能接受 SYN 报文段的最大客户端数量，即 SYN_REVD 状态(人们常说的半连接)的连接数。`默认值：128`。

设置小一点，可以用来防止 SYN 攻击。设置大一点，可以扩大并发连接

我们知道，服务器端一般使用 mq 来减轻高并发下的洪峰冲击，将暂时不能处理的请求放入队列，后续再慢慢处理。其实操作系统已经帮我们做了一些类似的东西了，这个东西就是 backlog。服务端一般通过 accept 调用，去获取 socket。但是假设我们的程序处理不过来（比如因为程序 bug，或者设计问题，没能及时地去调用 accept），那么此时的网络请求难道就直接丢掉吗？

当然不会！这时候，操作系统会帮我们放入 accept 队列，先暂存起来。等我们的程序缓过来了，直接调用 accept 去队列取就行了，这就达到了类似 mq 的效果。

而 backlog，和另一个参数 /proc/sys/net/core/somaxconn 一起，决定了队列的容量，算法为：min(/proc/sys/net/core/somaxconn, backlog) 。

**打个简单的比方：**

某某发布公告要邀请四海之内若干客人到场参加酒席。客人参加酒席分为两个步骤：

1. 到大厅；
2. 找到座位(吃东西，比如糖果、饭菜、酒等)。

tcp_max_syn_backlog 用于指定酒席现场面积允许容纳多少人进来；

somaxconn 用于指定有多少个座位。

显然 tcp_max_syn_backlog>=somaxconn。

如果要前来的客人数量超过 tcp_max_syn_backlog，那么多出来的人虽然会跟主任见面握手，但是要在门外等候；

如果到大厅的客人数量大于 somaxconn，那么多出来的客人就会没有位置坐(必须坐下才能吃东西)，只能等待有人吃完有空位了才能吃东西。.

### net.ipv4.tcp_max_tw_buckets = 5000

服务器 TIME-WAIT 状态连接的数量上限。`默认值：262144`。

如果超过这个数量， 新来的 TIME-WAIT 套接字会直接释放 (过多的 TIME-WAIT 套接字很影响服务器性能)。

### net.ipv4.tcp_synack_retries = 2

表示回应第二个握手包（SYN+ACK 包）给客户端 IP 后，如果收不到第三次握手包（ACK 包），进行重试的次数。`默认值：5`

修改这个参数为 0，可以加快回收半连接，减少资源消耗，但是有一个副作用：网络状况很差时，如果对方没收到第二个握手包，可能连接服务器失败，但对于一般网站，用户刷新一次页面即可。根据抓包经验，这种情况很少，但为了保险起见，可以只在被 tcp 洪水攻击时临时启用这个参数。之所以可以把 tcp_synack_retries 改为 0，因为客户端还有 tcp_syn_retries 参数，默认是 5，即使服务器端没有重发 SYN+ACK 包，客户端也会重发 SYN 握手包。

### net.ipv4.tcp_syncookies = INTEGER

当出现 SYN 等待队列溢出时，是否开启 cookies 来处理，可防范少量 SYN 攻击。`默认值：1`。

处在 SYN_RECV 的 TCP 连接称为半连接，存储在 SYN 队列。大量 SYN_RECV 会导致队列溢出，后续请求将被内核直接丢弃，也就是 SYN Flood 攻击。开启 syncookies 后，当 SYN 队列满了后，TCP 会通过原地址端口，目的地址端口和时间戳打造一个特别的 Sequence Number(又叫 cookie 发回去，如果是攻击者则不会有响应，如果是正常连接则把这个 SYNCookie 发回来，然后服务器端可以通过 cookie 建立连接(即使不在 SYN 队列)。

### ~~net.ipv4.tcp_tw_recycle = 1~~(高内核版本已移除该参数)

~~是否快速回收 TCP 连接中 TIME-WAIT sockets ，~~`默认值：0`

### net.ipv4.tcp_tw_reuse = 1

是否允许将 TIME-WAIT 状态的 sockets 重新用于新的 TCP 连接。`默认值：0`

### TCP keepalive 相关参数

注意：网络 keepalive 有两种

- TCP 层 keepalive
- HTTP 层 Keep-Alive

这里内核参数调整的是 TCP 层的 keepalive，常用于修复 ipvs 模式下长连接 timeout 问题 小于 900 即可

#### net.ipv4.tcp_keepalive_time = 7200

当一个 TCP 连接不再收到数据包后，经过 7200 秒后将当前连接标记为 keepalive 状态，并开始发送探测信息。将连接标记为需要保持活动状态后，将不再使用此计数器。`默认值：7200`。单位秒，即 2 小时

#### net.ipv4.tcp_keepalive_probes = 9

TCP 连接在 keepalive 状态下的探测次数。`默认值：9`

#### net.ipv4.tcp_keepalive_intvl = 75

TCP 连接在 keepalive 状态下的探测间隔。`默认值：75`。单位秒

该参数的值乘以 tcp_keepalive_probes 参数的值所得结果，是探测启动后，TCP 断开没有响应的连接的时间。默认情况下探测间隔是 75 秒，那么就说明，11 分钟后，没有响应的 TCP 连接将会断开

#### 总结

也就是说，当一个连接不再收到数据包后，经过 7200 秒后，开始发送探测，每隔 75 秒探测一次，当探测 9 次都失败时，将会断开 TCP 连接。

> 而且，7200 计时器是从连接建立时就开始计算的，不管有没有数据包，这个时间都会持续减小。

比如，我通过 ssh 与 sshd 建立了连接，就会发现有一个 120min 的计时器，这就是 7200 秒

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/qwuk68/1626272435911-652ae22c-fd68-4d2f-b12c-38e57f111971.png)

当我将 `net.ipv4.tcp_keepalive_time` 参数的值改为 60 时，再次登录就会发现只有 60 秒的计时器

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/qwuk68/1626272377586-25617789-baf8-4c2f-a992-d34eccc2f074.png)

**为什么需要 TCP keepalive？**

设想有一种场景：A 和 B 两边通过三次握手建立好 TCP 连接，然后突然间 B 就宕机了，之后时间内 B 再也没有起来。如果 B 宕机后 A 和 B 一直没有数据通信的需求，A 就永远都发现不了 B 已经挂了，那么 A 的内核里还维护着一份关于 A\&B 之间 TCP 连接的信息，浪费系统资源。于是在 TCP 层面引入了 keepalive 的机制，A 会定期给 B 发空的数据包，通俗讲就是心跳包，一旦发现到 B 的网络不通就关闭连接。这一点在 LVS 内尤为明显，因为 LVS 维护着两边大量的连接状态信息，一旦超时就需要释放连接。

## VS 相关参数

> 参考：
>
> - [官方文档，网络-IPvs Sysctl](https://www.kernel.org/doc/html/latest/networking/ipvs-sysctl.html)

### net.ipv4.vs.conn_reuse_mode = 1

值为 `0` 时，ipvs 不会对新连接进行重新负载，而是复用之前的负载结果，将新连接转发到原来的 rs 上；
值为 `1` 时，ipvs 则会对新连接进行重新调度。

### net.ipv4.vs.expire_nodest_conn = <0 | 非 0>

`默认值：0`
于控制连接的 rs 不可用时的处理。在开启时，如果后端 rs 不可用，会立即结束掉该连接，使客户端重新发起新的连接请求；否则将数据包**silently drop**，也就是 DROP 掉数据包但不结束连接，等待客户端的重试。

# bridge 参数

> 参考：
>
> - [官方文档，内核子系统文档-Networking-IP Sysctl-/proc/sys/net/bridge/* 变量](https://www.kernel.org/doc/html/latest/networking/ip-sysctl.html#proc-sys-net-bridge-variables)

### net.bridge.bridge-nf-call-iptables = 1

经过 bridge 网络设备的数据包，是否会被 iptables 进行过滤处理。`默认值：1`。若参数无法设置，加载 br_netfilter 模块即可

### net.bridge.bridge-nf-call-ip6tables = 1

经过 bridge 网络设备的数据包，是否会被 ip6tables 进行过滤处理。`默认值：1`。若参数无法设置，加载 br_netfilter 模块即可

### net.ipv6.conf.all.disable_ipv6 = 1

是否关闭所有网络设备(lo 除外)的 ipv6 协议。`默认值：0`。

### net.ipv6.conf.lo.disable_ipv6 = 1

是否关闭 lo 网络设备的 ipv6 协议。`默认值：0`。

### net.ipv4.neigh.default.gc_stale_time = 120

检查一次相邻层记录的有效性的周期(单位是秒)。当相邻层记录失效时，将在给它发送数据前，再解析一次。`默认值：60`。

### net.ipv4.conf.all.rp_filter = 0

是否校验所有网络设备上收到的数据包的源地址。`默认值：1`。

- 0：不开启源地址校验。
- 1：开启严格的反向路径校验。对每个进来的数据包，校验其反向路径是否是最佳路径。如果反向路径不是最佳路径，则直接丢弃该数据包。
- 2：开启松散的反向路径校验。对每个进来的数据包，校验其源地址是否可达，即反向路径是否能通（通过任意网口），如果反向路径不同，则直接丢弃该数据包。

rp_filter 参数详解见 [rp_filter](/docs/1.操作系统/2.Kernel/1.Linux%20Kernel/Kernel%20参数/net(网络相关参数)/rp_filter.md)
