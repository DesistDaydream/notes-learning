---
title: Linux 网络栈管理
linkTitle: Linux 网络栈管理
weight: 1
---

# 概述

> 参考：
>
> - [Kernel 文档，Linux Networking Documentation](https://www.kernel.org/doc/html/latest/networking/index.html)
> - [Kernel 文档，Linux Networking and Network Devices APIs](https://www.kernel.org/doc/html/latest/networking/kapi.html)
> - [arthurchiao.art 的文章](http://arthurchiao.art/index.html)
>   - [[译] Linux 网络栈监控和调优：接收数据（2016）](http://arthurchiao.art/blog/tuning-stack-rx-zh/) — 老文章，作者不再推荐，查看 2022 之后的文章
>   - [[译] Linux 网络栈监控和调优：发送数据（2017）](http://arthurchiao.art/blog/tuning-stack-tx-zh/)
>   - [Linux 网络栈原理、监控与调优：前言（2022）](https://arthurchiao.art/blog/linux-net-stack-zh/)

和磁盘设备类似，Linux 用户想要使用网络功能，不能通过直接操作硬件完成，而需要直接或间接的操作一个 Linux 为我们抽象出来的设备，i.e. 通用的 **Linux 网络设备**来完成。一个常见的情况是，计算机里装有一个硬件网卡，Linux 会在系统里为其生成一个网络设备实例，e.g. eth0，用户需要对 eth0 发出命令以配置或使用它了。更多的硬件会带来更多的网络设备实例（另外，[虚拟化](/docs/10.云原生/Virtualization/Virtualization.md)的硬件也会带来更多的网络设备实例）。

网卡本身并不会连接连接任何网络，网卡需要相应的配置文件来告诉他们如何实现网络连接。而让网卡与配置文件关联的过程，就是 network.service 这类服务来实现的

在 [Linux Kernel](/docs/1.操作系统/Kernel/Linux%20Kernel/Linux%20Kernel.md) 中，一般使用“网络设备”这种称呼，来描述硬件物理网卡设备在系统中的实例。在不同的语境中，有时也简称为 “设备”、“DEV” 等等。网络设备可以是一块真实机器上的网卡，也可以是创建的虚拟的网卡。

而网络设备与网卡之间如何建立关系，就是网卡驱动程序的工作了，不同的网卡，驱动不一样，可以实现的功能也各有千秋。所以，想要系统出现 eth0 这种网络设备，网卡驱动程序是必须存在的，否则，没有驱动，也就无法识别硬件，无法识别硬件，在系统中也就不知道如何操作这个硬件。

## 常见术语

### DataPath(数据路径)

网络数据在内核中进行网络传输时，所经过的所有点组合起来，称为数据路径。

### Socket Buffer(简称 sk_buff 或 skb)

在内核代码中是一个名为 [**sk_buff**](https://www.kernel.org/doc/html/latest/networking/kapi.html#c.sk_buff) 的结构体。内核显然需要一个数据结构来储存报文的信息。这就是 skb 的作用。

sk_buff 结构自身并不存储报文内容，它通过多个指针指向真正的报文内存空间:

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/efrsi8/1617849698535-471768e0-dcf8-4471-8dd2-605a1bc4e020.png)

sk_buff 是一个贯穿整个协议栈层次的结构，在各层间传递时，内核只需要调整 sk_buff 中的指针位置就行。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/efrsi8/1617849692989-54095177-b85c-449e-8c66-3b026e4925da.png)

### DEVICE(设备)

在内核代码中，是一个名为 [**net_device**](https://www.kernel.org/doc/html/latest/networking/kapi.html#c.net_device) 的结构体。一个巨大的数据结构，描述一个网络设备的所有 属性、数据 等信息。

# Linux 网络功能的实现

# 数据包的 Transmit(发送) 与 Receive(接收) 过程概览

## Receive(接收) 过程

本文将拿 **Intel I350** 网卡的 `igb` 驱动作为参考，网卡的 data sheet 这里可以下 载 [PDF](http://www.intel.com/content/dam/www/public/us/en/documents/datasheets/ethernet-controller-i350-datasheet.pdf) （警告：文件很大）。
从比较高的层次看，一个数据包从被网卡接收到进入 socket 接收队列的整个过程如下：

1. 加载网卡驱动，初始化
2. 包从外部网络进入网卡
3. 网卡（通过 DMA）将包 copy 到内核内存中的 ring buffer
4. 产生硬件中断，通知系统收到了一个包
5. 驱动调用 NAPI，如果轮询（poll）还没开始，就开始轮询
6. `ksoftirqd` 进程调用 NAPI 的 `poll` 函数从 ring buffer 收包（`poll` 函数是网卡 驱动在初始化阶段注册的；每个 CPU 上都运行着一个 `ksoftirqd` 进程，在系统启动期 间就注册了）
7. ring buffer 里包对应的内存区域解除映射（unmapped）
8. （通过 DMA 进入）内存的数据包以 `skb` 的形式被送至更上层处理
9. 如果 packet steering 功能打开，或者网卡有多队列，网卡收到的包会被分发到多个 CPU
10. 包从队列进入协议层
11. 协议层处理包
12. 包从协议层进入相应 socket 的接收队列

接下来会详细介绍这个过程。

## Transmit(发送) 过程

本文将拿**Intel I350**网卡的 `igb` 驱动作为参考，网卡的 data sheet 这里可以下载 [PDF](http://www.intel.com/content/dam/www/public/us/en/documents/datasheets/ethernet-controller-i350-datasheet.pdf) （警告：文件很大）。
从比较高的层次看，一个数据包从用户程序到达硬件网卡的整个过程如下：

1. 使用 **系统调用**（如 `sendto`，`sendmsg` 等）写数据
2. 数据穿过 **socket 子系统**，进入**socket 协议族**（protocol family）系统（在我们的例子中为 `AF_INET`）
3. 协议族处理：数据穿过 **协议层**，这一过程（在许多情况下）会将 **数据**（data）转换成 **数据包**（packet）
4. 数据穿过 **路由层**，这会涉及路由缓存和 ARP 缓存的更新；如果目的 MAC 不在 ARP 缓存表中，将触发一次 ARP 广播来查找 MAC 地址
5. 穿过协议层，packet 到达 **设备无关层**（device agnostic layer）
6. 使用 XPS（如果启用）或散列函数 **选择发送队列**
7. 调用网卡驱动的 **发送函数**
8. 数据传送到网卡的 `qdisc`（queue discipline，排队规则）
9. qdisc 会直接 **发送数据**（如果可以），或者将其放到队列，下次触发 `**NET_TX**` **类型软中断**（softirq）的时候再发送
10. 数据从 qdisc 传送给驱动程序
11. 驱动程序创建所需的 **DMA 映射**，以便网卡从 RAM 读取数据
12. 驱动向网卡发送信号，通知 **数据可以发送了**
13. **网卡从 RAM 中获取数据并发送**
14. 发送完成后，设备触发一个 **硬中断**（IRQ），表示发送完成
15. **硬中断处理函数** 被唤醒执行。对许多设备来说，这会 **触发 `NET_RX` 类型的软中断**，然后 NAPI poll 循环开始收包
16. poll 函数会调用驱动程序的相应函数，**解除 DMA 映射**，释放数据

# Linux 网络栈关联文件与配置

不同的 Linux 发行版，所使用的上层网络配置程序各不相同，各种程序所读取的配置文件也各不相同。

- 对于 RedHat 相关的发行版，网络配置在 /etc/sysconfig/network-scripts/ 目录中
- 对于 Debian 相关的发行版，网络配置在 /etc/network/ 目录中

在这些目录中，其实都是通过脚本来实现的

后来随着时代的发展，涌现出很多通用的网络管理程序，比如 [Netplan](/docs/1.操作系统/Kernel/Network/Linux%20网络栈管理/Netplan/Netplan.md)、[NetworkManager](/docs/1.操作系统/Kernel/Network/Linux%20网络栈管理/NetworkManager/NetworkManager.md)、[systemd-networkd](/docs/1.操作系统/Kernel/Network/Linux%20网络栈管理/systemd-networkd.md)、etc.，这样就可以让各个发行版使用相同的程序来管理网络了，减少切换发行版而需要学习对应配置的成本，并且也更利于发展。

## /proc/net/

https://github.com/torvalds/linux/blob/v6.18/Documentation/filesystems/proc.rst#13-networking-info-in-procnet

该目录是 **软连接文件**，被连接到 **self/net** 目录。主要记录的是，链接到的进程的 network namespace 的信息

我们可以利用这些信息来查看系统中哪些网络设备可用，以及有多少流量通过这些设备路由，比如：

```
> cat /proc/net/dev
Inter-|Receive                                                    |Transmit
 face |bytes     packets errs drop fifo frame compressed multicast|bytes      packets errs drop fifo colls carrier compressed
    lo:908188    5596    0    0    0    0     0          0         908188     5596    0    0    0    0     0       0
  ppp0:15475140  20721   410  0    0    410   0          0         1375103    17405   0    0    0    0     0       0
  eth0:614530    7085    0    0    0    0     0          1         1703981    5535    0    0    0    3     0       0
```

此外，每个 Bond 接口都有自己的目录。e.g. bond0 设备会有一个名为 /proc/net/bond0/ 的目录。该目录包含特定于该绑定的信息，例如绑定的当前从设备、从设备的链路状态以及从设备的链路故障次数。

**./dev** # 网络设备的统计信息。包括 Interface(接口名称), Receive(接收的数据), Transmit(发送的数据) 三部分

**./nf_conntrack** # 链接跟踪表，该文件用于记录已跟踪的连接

**./tcp** # 所有的 TCP 连接信息。

**./tcp6** # 所有的基于 IPv6 的 TCP 连接信息。

参考：[GitHub Linux 项目文档](https://github.com/torvalds/linux/blob/master/Documentation/networking/proc_net_tcp.rst)

保存 TCP 套接字表的转储。除了调试之外，大部分信息都没有什么用。

- sl # 值是套接字的内核哈希槽位
- local_address # 是本地地址和端口号对
- rem_address # 是远程地址和端口号对(如果连接)
- St # 是套接字的内部状态。根据内核内存使用情况，
- tx_queue 和 rx_queue # 是传出和传入的数据队列。
- tr、tm->when 和 rexmits # 字段保存内核套接字状态的内部信息，仅在调试时有用。
- uid # 字段保存套接字创建者的有效 uid。
- inode # 该 socket 的 inode 号，后面一串 16 进制的字符是该 socket 在内存中的地址。

```bash
root@desistdaydream:~# cat /proc/net/tcp
  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:177A 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12975942 1 ffff923dd621a300 100 0 0 10 0
   1: 3500007F:0035 00000000:0000 0A 00000000:00000000 00:00000000 00000000   101        0 28017 1 ffff923ef9dd08c0 100 0 0 10 0
   2: 00000000:0016 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 33221 1 ffff923eecf6c600 100 0 0 10 0
   3: F82A13AC:0016 CB2A13AC:FD4C 01 00000000:00000000 02:00025EB2 00000000     0        0 12973284 4 ffff923dd621e900 20 4 29 10 20
   4: F82A13AC:0016 CB2A13AC:FD48 01 00000000:00000000 02:000A6D4A 00000000     0        0 12944563 2 ffff923dd621bd40 20 4 31 10 23
```

注意：

这里用 16 进制表示的 IP 有点奇葩比如 `F82A13AC` 转换成 IP 地址是 `248.42.19.172`，真实 IP 地址是 `172.19.42.248`，也就是说反过来了。。。。`F82A13AC` 应该是 `AC132AF8`

**./udp** # 所有 UDP 连接信息

**./udp6** # 所有基于 IPv6 的 UDP 连接信息

**./unix** # 所有 Unix Domain Socket 连接信息

## /sys/class/net/

详见 [Linux 网络设备](/docs/1.操作系统/Kernel/Network/Linux%20网络栈管理/Linux%20网络设备/Linux%20网络设备.md)
