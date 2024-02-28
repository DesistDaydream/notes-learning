---
title: Connnection Tracking(连接跟踪)
---

# 概述

> 参考：
>
> - [arthurchiao，连接跟踪：原理、应用及 Linux 内核实现](http://arthurchiao.art/blog/conntrack-design-and-implementation-zh/)

**Connection Tracking(连接跟踪系统，简称 ConnTrack、CT)**，用于跟踪并且记录连接状态。CT 是 [Linux 网络流量控制](/docs/1.操作系统/2.Kernel/8.Network/Linux%20网络流量控制/Linux%20网络流量控制.md)的基础

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ynfo7m/1617860207674-43ea3c6d-0d0f-4fac-bccb-90e752e75a47.png)
例如，上图是一台 IP 地址为 `10.1.1.2` 的 Linux 机器，我们能看到这台机器上有三条 连接：

1. 机器访问外部 HTTP 服务的连接（目的端口 80）
2. 外部访问机器内 FTP 服务的连接（目的端口 21）
3. 机器访问外部 DNS 服务的连接（目的端口 53）

连接跟踪所做的事情就是发现并跟踪这些连接的状态，具体包括：

- 从数据包中提取**元组**（tuple）信息，辨别**数据流**（flow）和对应的**连接**（connection）
- 为所有连接维护一个**状态数据库**（conntrack table），例如连接的创建时间、发送 包数、发送字节数等等
- 回收过期的连接（GC）
- 为更上层的功能（例如 NAT）提供服务

需要注意的是，**连接跟踪中所说的“连接”，概念和 TCP/IP 协议中“面向连接”（ connection oriented）的“连接”并不完全相同**，简单来说：

- TCP/IP 协议中，连接是一个四层（Layer 4）的概念。
  - TCP 是有连接的，或称面向连接的（connection oriented），发送出去的包都要求对端应答（ACK），并且有重传机制
  - UDP 是无连接的，发送的包无需对端应答，也没有重传机制
- CT 中，一个元组（tuple）定义的一条数据流（flow ）就表示一条连接（connection）。
  - 后面会看到 UDP 甚至是 **ICMP 这种三层协议在 CT 中也都是有连接记录的**
  - 但**不是所有协议都会被连接跟踪**

本文中用到“连接”一词时，大部分情况下指的都是后者，即“连接跟踪”中的“连接”。

# 原理

要跟踪一台机器的所有连接状态，就需要：

- **拦截（或称过滤）流经这台机器的每一个数据包，并进行分析**。
- 根据这些信息**建立**起这台机器上的**连接信息数据库**（conntrack table）。
- 根据拦截到的包信息，不断更新数据库

例如

- 拦截到一个 TCP SYNC 包时，说明正在尝试建立 TCP 连接，需要创建一条新 conntrack entry 来记录这条连接
- 拦截到一个属于已有 conntrack entry 的包时，需要更新这条 conntrack entry 的收发包数等统计信息

除了以上两点功能需求，还要考虑**性能问题**，因为连接跟踪要对每个包进行过滤和分析 。性能问题非常重要，但不是本文重点，后面介绍实现时会进一步提及。
之外，这些功能最好还有配套的管理工具来更方便地使用。

## Connection Tracking 的实现

现在(2021 年 4 月 9 日)提到连接跟踪（conntrack），可能首先都会想到 Netfilter。但由上节讨论可知， 连接跟踪概念是独立于 Netfilter 的，**Netfilter 只是 Linux 内核中的一种连接跟踪实现**。

换句话说，**只要具备了 hook 能力，能拦截到进出主机的每个包，完全可以在此基础上自 己实现一套连接跟踪**。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ynfo7m/1617861067581-3b23cb80-cd1f-4d7d-9767-57581c62233b.png)

云原生网络方案 Cilium 在 `1.7.4+` 版本就实现了这样一套独立的连接跟踪和 NAT 机制 （完备功能需要 Kernel `4.19+`）。其基本原理是：

1. 基于 BPF hook 实现数据包的拦截功能（等价于 netfilter 里面的 hook 机制）
2. 在 BPF hook 的基础上，实现一套全新的 conntrack 和 NAT

因此，即便[卸载 Netfilter](https://github.com/cilium/cilium/issues/12879) ，也不会影响 Cilium 对 Kubernetes ClusterIP、NodePort、ExternalIPs 和 LoadBalancer 等功能的支持。

由于这套连接跟踪机制是独立于 Netfilter 的，因此它的 conntrack 和 NAT 信息也没有 存储在内核的（也就是 Netfilter 的）conntrack table 和 NAT table。所以常规的 `conntrack/netstats/ss/lsof` 等工具是看不到的，要使用 Cilium 的命令，例如：

```bash
cilium bpf nat list
cilium bpf ct list global
```

配置也是独立的，需要在 Cilium 里面配置，例如命令行选项 `--bpf-ct-tcp-max`。

另外，本文会多次提到连接跟踪模块和 NAT 模块独立，但**出于性能考虑，具体实现中 二者代码可能是有耦合的**。例如 Cilium 做 conntrack 的垃圾回收（GC）时就会顺便把 NAT 里相应的 entry 回收掉，而非为 NAT 做单独的 GC。

### Netfilter 中的 Connection Tracking

详见 Netifilter 中 [Connection Tracking(连接跟踪)机制](/docs/1.操作系统/2.Kernel/8.Network/Linux%20网络流量控制/Netfilter%20流量控制系统/Connection%20Tracking(连接跟踪)机制.md)

### BPF 中的 Connection Tracking

## Connection Tracking 的应用实践

来看几个 conntrack 的具体应用。

### 网络地址转换（NAT）

网络地址转换（NAT），名字表达的意思也比较清楚：对（数据包的）网络地址（`IP + Port`）进行转换。

![NAT 及其内核位置示意图](https://notes-learning.oss-cn-beijing.aliyuncs.com/ynfo7m/1617861112761-712ea43a-0ed4-4c94-8758-3f593fc3a1b6.png)

例如上图中，机器自己的 IP `10.1.1.2` 是能与外部正常通信的，但 `192.168` 网段是私有 IP 段，外界无法访问，也就是说源 IP 地址是 `192.168` 的包，其**应答包是无 法回来的**。因此

- 当源地址为 `192.168` 网段的包要出去时，机器会先将源 IP 换成机器自己的 `10.1.1.2` 再发送出去；
- 收到应答包时，再进行相反的转换。

这就是 NAT 的基本过程。

Docker 默认的 `bridge` 网络模式就是这个原理。每个容器会分一个私有网段的 IP 地址，这个 IP 地址可以在宿主机内的不同容器之间通信，但容器流量出宿主机时要进行 NAT。
NAT 又可以细分为几类：

- SNAT：对源地址（source）进行转换
- DNAT：对目的地址（destination）进行转换
- Full NAT：同时对源地址和目的地址进行转换

以上场景属于 SNAT，将不同私有 IP 都映射成同一个“公有 IP”，以使其能访问外部网络服 务。这种场景也属于正向代理。

NAT 依赖连接跟踪的结果。连接跟踪**最重要的使用场景**就是 NAT。

#### 四层负载均衡（L4LB）

再将范围稍微延伸一点，讨论一下 NAT 模式的四层负载均衡。

四层负载均衡是根据包的四层信息（例如 `src/dst ip, src/dst port, proto`）做流量分发。

VIP（Virtual IP）是四层负载均衡的一种实现方式：

- 多个后端真实 IP（Real IP）挂到同一个虚拟 IP（VIP）上
- 客户端过来的流量先到达 VIP，再经负载均衡算法转发给某个特定的后端 IP

如果在 VIP 和 Real IP 节点之间使用的 NAT 技术（也可以使用其他技术），那客户端访 问服务端时，L4LB 节点将做双向 NAT（Full NAT），数据流如下图所示：

![L4LB: Traffic path in NAT mode](https://notes-learning.oss-cn-beijing.aliyuncs.com/ynfo7m/1617861112756-297f87b7-f40e-4886-899a-65629964fd2c.png)

### 1.5.2 有状态防火墙

有状态防火墙（stateful firewall）是相对于早期的**无状态防火墙**（stateless firewall）而言的：早期防火墙只能写 `drop syn` 或者 `allow syn` 这种非常简单直接 的规则，**没有 flow 的概念**，因此无法实现诸如 **“如果这个 ack 之前已经有 syn， 就 allow，否则 drop”** 这样的规则，使用非常受限。

显然，要实现有状态防火墙，就必须记录 flow 和状态，这正是 conntrack 做的事情。

来看个更具体的防火墙应用：OpenStack 主机防火墙解决方案 —— 安全组（security group）。

#### OpenStack 安全组

简单来说，安全组实现了**虚拟机级别**的安全隔离，具体实现是：在 node 上连接 VM 的 网络设备上做有状态防火墙。在当时，最能实现这一功能的可能就是 Netfilter/iptables。
回到宿主机内网络拓扑问题： OpenStack 使用 OVS bridge 来连接一台宿主机内的所有 VM。 如果只从网络连通性考虑，那每个 VM 应该直接连到 OVS bridge `br-int`。但这里问题 就来了：

- OVS 没有 conntrack 模块，
- Linux 中有 conntrack 模块，但基于 conntrack 的防火墙**工作在 IP 层**（L3），通过 iptables 控制，
- 而 **OVS 是 L2 模块**，无法使用 L3 模块的功能，

因此无法在 OVS （连接虚拟机）的设备上做防火墙。

所以，2016 之前 OpenStack 的解决方案是，在每个 OVS 和 VM 之间再加一个 Linux bridge ，如下图所示，

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ynfo7m/1617861113322-0f9c4ca7-ffca-43ab-840d-78db13d23008.png)

Fig 1.6. Network topology within an OpenStack compute node, picture from [Sai's Blog](https://thesaitech.wordpress.com/2017/09/24/how-to-trace-the-tap-interfaces-and-linux-bridges-on-the-hypervisor-your-openstack-vm-is-on/)

Linux bridge 也是 L2 模块，按道理也无法使用 iptables。但是，**它有一个 L2 工具 ebtables，能够跳转到 iptables**，因此间接支持了 iptables，也就能用到 Netfilter/iptables 防火墙的功能。

这种 workaround 不仅丑陋、增加网络复杂性，而且会导致性能问题。因此， RedHat 在 2016 年提出了一个 OVS conntrack 方案 \[7]，从那以后，才有可能干掉 Linux bridge 而仍然具备安全组的功能。
