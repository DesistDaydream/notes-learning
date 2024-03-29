---
title: BPF
linkTitle: BPF
date: 2023-11-02T23:51
weight: 1
---

# 概述

> 参考：
>
> - [Wiki，BPF](https://en.wikipedia.org/wiki/Berkeley_Packet_Filter)
> - [GitHub 项目,bcc-BPF 特性与 LInux 内核版本对照表](https://github.com/iovisor/bcc/blob/master/docs/kernel-versions.md)

《Linux 内核观测技术 BPF》

**Berkeley Packet Filter(伯克利包过滤器，简称 BPF)**，是类 Unix 系统上数据链路层的一种原始接口，提供原始链路层封包的收发。在 Kernel 官方文档中，BPF 也称为 **Linux Socket Filtering(LInux 套接字过滤，简称 LSF)**。BPF 有时也只表示 **filtering mechanism(过滤机制)**，而不是整个接口。

**注意：不管是后面描述的 eBPF 还是 BPF，这个名字或缩写，其本身所表达的含义，其实已经没有太大的意义了，因为这个项目的发展远远超出了它最初的构想。**

在 BPF 之前，如果想做数据包过滤，则必须将所有数据包复制到用户空间中，然后在那里过滤它们，这种方式意味着必须将所有数据包复制到用户空间中，复制数据的开销很大。当然可以通过将过滤逻辑转移到内核中解决开销问题，我们来看 BPF 做了什么工作。

实际上，BPF 最早称为 **BSD Packet Filter**，是很早就有的 Unix 内核特性，最早可追溯到 1992 年发表在 USENIX Conference 上的一篇论文[《BSD 数据包过滤：一种新的用户级包捕获架构》](http://www.tcpdump.org/papers/bpf-usenix93.pdf)，这篇文章作者描述了他们如何在 Unix 内核实现网络数据包过滤，这种技术比当时最先进的数据包过滤技术快了 20 倍。这篇文章描述的 BPF 在数据包过滤上引入了两大革新：

- 一个新的虚拟机设计，可以有效得工作在基于寄存器结构的 CPU 之上。
- 应用程序使用缓存只复制与过滤数据包相关的数据，不会复制数据包的宿友信息。这样可以最大程度得减少 BPF 处理的数据。

随后，得益于如此强大的性能优势，所有 Unix 系统都将 BPF 作为网络包过滤的首选技术，抛弃了消耗更多内存和性能更差的原有技术实现。后来由于 BPF 的理念逐渐成为主流，为各大操作系统所接受，这样早期 "B" 所代表的 BSD 便渐渐淡去，最终演化成了今天我们眼中的 BPF(Berkeley Packet Filter)。

比如我们熟知的 Tcpdump 程序，其底层就是依赖 BPF 实现的包过滤。我们可以在命令后面增加 ”-d“ 来查看 tcpdump 过滤条件的底层汇编指令。

```bash
~]# tcpdump -d 'ip and tcp port 8080'
(000) ldh      [12]
(001) jeq      #0x800           jt 2 jf 12
(002) ldb      [23]
(003) jeq      #0x6             jt 4 jf 12
(004) ldh      [20]
(005) jset     #0x1fff          jt 12 jf 6
(006) ldxb     4*([14]&0xf)
(007) ldh      [x + 14]
(008) jeq      #0x1f90          jt 11 jf 9
(009) ldh      [x + 16]
(010) jeq      #0x1f90          jt 11 jf 12
(011) ret      #262144
(012) ret      #0
```

-dd 可以打印字节码

```bash
~]# tcpdump -dd 'ip and tcp port 8080'
{ 0x28, 0, 0, 0x0000000c },
{ 0x15, 0, 10, 0x00000800 },
{ 0x30, 0, 0, 0x00000017 },
{ 0x15, 0, 8, 0x00000006 },
{ 0x28, 0, 0, 0x00000014 },
{ 0x45, 6, 0, 0x00001fff },
{ 0xb1, 0, 0, 0x0000000e },
{ 0x48, 0, 0, 0x0000000e },
{ 0x15, 2, 0, 0x00001f90 },
{ 0x48, 0, 0, 0x00000010 },
{ 0x15, 0, 1, 0x00001f90 },
{ 0x6, 0, 0, 0x00040000 },
{ 0x6, 0, 0, 0x00000000 },
```

## BPF 的进化

得益于 BPF 在包过滤上的良好表现，Alexei Starovoitov 对 BPF 进行彻底的改造，并增加了新的功能，改善了它的性能，这个新版本被命名为 **extended BPF(扩展的 BPF，简称 eBPF)**，新版本的 BPF 全面兼容并扩充了原有 BPF 的功能。因此，将传统的 BPF 重命名为 **classical BPF(传统的 BPF，简称 cBPF)**，相对应的，新版本的 BPF 则命名为 eBPF 或直接称为 BPF(**所以，我们现在所说的 BPF，大部分情况下就是指 eBPF**)。Linux Kernel 3.18 版本开始实现对 eBPF 的支持。

# eBPF 概述

> 参考：
>
> - [eBPF 官网](https://ebpf.io/)
> - [某网站系列文章](http://kerneltravel.net/categories/ebpf/)

**extended Berkeley Packet Filter(扩展的 BPF，简称 eBPF)** 起源于 BPF，是对 BPF 的扩展。eBPF 针对现代硬件进行了优化和全新的设计，使其生成的指令集比 cBPF 解释器生成的机器码更快。这个扩展版本还将 cBPF VM 中的寄存器数量从两个 32 位寄存器增加到 10 个 64 位寄存器。寄存器数量和寄存器宽度的增加为编写更复杂的程序提供了可能性，开发人员可以自由的使用函数参数交换更多的信息。这些改进使得 eBPF 比原来的 cBPF 快四倍。这些改进，主要还是对网络过滤器内部处理的 eBPF 指令集进行优化，仍然被限制在内核空间中，只有少数用户空间中的程序可以编写 BPF 过滤器供内核处理，比如 Tcpdump 和 Seccomp。

除了上述的优化之外，eBPF 最让人兴奋的改进，是其向用户空间的开放。开发者可以在用户空间，编写 eBPF 程序，并将其加在到内核空间执行。虽然 eBPF 程序看起来更像内核模块，但与内核模块不同的是，eBPF 程序不需要开发者重新编译内核，而且保证了在内核不崩溃的情况下完成加载操作，着重强调了安全性和稳定性。BPF 代码的主要贡献单位主要包括 Cilium、Facebook、Red Hat 以及 Netronome 等。

Linux Kernel 一直是实现 可观察性、网络、安全性 的理想场所。不幸的是，想要自定义这些实现通常是不切实际的，因为它需要更改内核源代码或加载内核模块，并导致彼此堆叠的抽象层。而 eBPF 的出现，让这一切成为可能，**eBPF 可以在 Linux 内核中运行沙盒程序，而无需更改内核源代码或加载内核模块**。通过**使 Linux Kernel 可编程**，基础架构软件可以利用现有的层，从而使它们更加智能和功能丰富，而无需继续为系统增加额外的复杂性层。

也正由于此，eBPF 不再局限于网络的过滤，而且 eBPF 就相当于内核本身的代码，想象空间无限，并且热加载到内核，换句话说，一旦加载到内核，内核的行为就变了。所以，eBPF 带动了 安全性、应用程序配置/跟踪、性能故障排除 等等领域的新一代工具的开发，这些工具不再依赖现有的内核功能，而是在不影响执行效率或安全性的情况下主动重新编程运行时行为：

- **Networking(网络)**
- **Observability(可观测性)**
  - **Monitoring(监控)**
  - **Tracing(跟踪)**
  - **Profiling(分析)**
- **Security(安全)**
- **等等，随着发展，eBPF 还可以实现更多!**

可以这么说，BPF 的种种能力，实现了 **Software Define Kernel(软件定义内核)**。

## eBPF 与 内核模块 的对比

在 Linux 观测方面，eBPF 总是会拿来与 kernel 模块方式进行对比，eBPF 在安全性、入门门槛上比内核模块都有优势，这两点在观测场景下对于用户来讲尤其重要。

| 维度                | Linux 内核模块                       | eBPF                                           |
| ------------------- | ------------------------------------ | ---------------------------------------------- |
| kprobes/tracepoints | 支持                                 | 支持                                           |
| 安全性              | 可能引入安全漏洞或导致内核 Panic     | 通过验证器进行检查，可以保障内核安全           |
| 内核函数            | 可以调用内核函数                     | 只能通过 BPF Helper 函数调用                   |
| 编译性              | 需要编译内核                         | 不需要编译内核，引入头文件即可                 |
| 运行                | 基于相同内核运行                     | 基于稳定 ABI 的 BPF 程序可以编译一次，各处运行 |
| 与应用程序交互      | 打印日志或文件                       | 通过 perf_event 或 map 结构                    |
| 数据结构丰富性      | 一般                                 | 丰富                                           |
| 入门门槛            | 高                                   | 低                                             |
| 升级                | 需要卸载和加载，可能导致处理流程中断 | 原子替换升级，不会造成处理流程中断             |
| 内核内置            | 视情况而定                           | 内核内置支持                                   |

# eBPF 发展史

eBPF 是如何诞生的呢？我最初开始讲起。这里“最初”我指的是 2013 年之前。

## 2013

### 传统的流量控制工具和系统

回顾一下当时的 “SDN” 蓝图。

1. 当时有 OpenvSwitch（OVS）、`tc`（Traffic control），以及内核中的 Netfilter 子系 统（包括 `iptables`、`ipvs`、`nftalbes` 工具），可以用这些工具对 datapath 进行“ 编程”：。
2. BPF 当时用于 `tcpdump`，**在内核中尽量前面的位置抓包**，它不会 crash 内核；此 外，它还用于 seccomp，**对系统调用进行过滤**（system call filtering），但当时 使用的非常受限，远不是今天我们已经在用的样子。
3. 此外就是前面提到的 feature creeping 问题，以及 **tc 和 netfilter 的代码重复问题，因为这两个子系统是竞争关系**。
4. **OVS 当时被认为是内核中最先进的数据平面**，但它最大的问题是：与内核中其他网 络模块的集成不好【译者注 1】。此外，很多核心的内核开发者也比较抵触 OVS，觉得它很怪。

> 【译者注 1】例如，OVS 的 internal port、patch port 用 tcpdump 都是 [抓不到包的](http://arthurchiao.art/blog/ovs-deep-dive-4-patch-port/)，排障非常不方便。

### eBPF 与 传统流量控制 的区别

对比 eBPF 和这些已经存在很多年的工具：

1. tc、OVS、netfilter 可以对 datapath 进行“编程”：但前提是 datapath 知道你想做什 么（but only if the datapath knows what you want to do）。
   - 只能利用这些工具或模块提供的既有功能。
2. eBPF 能够让你**创建新的 datapath**（eBPF lets you create the datapath instead）。

> - eBPF 就是内核本身的代码，想象空间无限，并且热加载到内核；换句话说，一旦加 载到内核，内核的行为就变了。
> - 在 eBPF 之前，改变内核行为这件事情，只能通过修改内核再重新编译，或者开发内 核模块才能实现。

译者注

### eBPF：第一个（巨型）patch

- 描述 eBPF 的 RFC 引起了广泛讨论，但普遍认为侵入性太强了（改动太大）。
- 另外，当时 nftables (inspired by BPF) 正在上升期，它是一个与 eBPF 有点类似的 BPF 解释器，大家不想同时维护两个解释器。

最终这个 patch 被拒绝了。
被拒的另外一个原因是前面提到的，没有遵循“大改动小提交”原则，全部代码放到了一个 patch。Linus 会疯的。

## 2014

### 第一个 eBPF patch 合并到内核

- 用一个**扩展（extended）指令集**逐步、全面替换原来老的 BPF 解释器。
- **自动新老 BPF 转换**：in-kernel translation。
- 后续 patch 将 eBPF 暴露给 UAPI，并添加了 verifier 代码和 JIT 代码。
- 更多后续 patch，从核心代码中移除老的 BPF。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ksq56w/1617847999424-2c5e5133-2862-4574-bffb-0776c6b0aa4b.png)

我们也从那时开始，顺理成章地成为了 eBPF 的 maintainer。

### Kubernetes 提交第一个 commit

巧合的是，**对后来影响深远的 Kubernetes，也在这一年提交了第一个 commit**：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ksq56w/1617847999414-121919ac-67f3-4841-b2a0-c0bd1668c7bd.png)

## 2015

### eBPF 分成两个方向：networking & tracing

到了 2015 年，eBPF 开发分成了两个方向：

- networking
- tracing

### eBPF backend 合并到 LLVM 3.7

这一年的一个重要里程碑是 eBPF backend 合并到了 upstream LLVM 编译器套件，因此你 现在才能用 clang 编译 eBPF 代码。

### 支持将 eBPF attach 到 kprobes

这是 tracing 的第一个使用案例。
Alexei 主要负责 tracing 部分，他添加了一个 patch，支持加载 eBPF 用来做 tracing， 能获取系统的观测数据。

### 通过 cls_bpf，tc 变得完全可编程

我主要负责 networking 部分，使 tc 子系统可编程，这样我们就能用 eBPF 来灵活的对 datapath 进行编程，获得一个高性能 datapath。

### 为 tc 添加了一个 lockless ingress & egress hook 点

> 译注：可参考：
>
> - [深入理解 tc ebpf 的 direct-action (da) 模式（2020）](http://arthurchiao.art/blog/understanding-tc-da-mode-zh/)
> - [为容器时代设计的高级 eBPF 内核特性（FOSDEM, 2021）](http://arthurchiao.art/blog/advanced-bpf-kernel-features-for-container-age-zh/)

### 添加了很多 verifer 和 eBPF 辅助代码（helper）

使用更方便。

### bcc 项目发布

作为 tracing frontend for eBPF。

## 2016

### eBPF 添加了一个新 fast path：XDP

- XDP 合并到内核，支持在驱动的 ingress 层 attach BPF 程序。
- nfp 最为第一家网卡及驱动，支持将 eBPF 程序 offload 到 cls_bpf & XDP hook 点。

### Cilium 项目发布

Cilium 最开始的目标是 **docker 网络解决方案**。

- 通过 eBPF 实现高效的 label-based policy、NAT64、tunnel mesh、容器连通性。
- 整个 datapath & forwarding 逻辑全用 eBPF 实现，不再需要 Docker 或 OVS 桥接设备。

## 2017

### eBPF 开始大规模应用于生产环境

2016 ~ 2017 年，eBPF 开始应用于生产环境：

1. Netflix on eBPF for tracing: ‘Linux BPF superpowers’
2. Facebook 公布了生产环境 XDP+eBPF 使用案例（DDoS & LB）
   - 用 XDP/eBPF 重写了原来基于 IPVS 的 L4LB，性能 `10x`。
   - **eBPF 经受住了严苛的考验**：从 2017 开始，每个进入 facebook.com 的包，都是经过了 XDP & eBPF 处理的。
3. Cloudflare 将 XDP+BPF 集成到了它们的 DDoS mitigation 产品。
   - 成功将其组件从基于 Netfilter 迁移到基于 eBPF。
   - 到 2018 年，它们的 XDP L4LB 完全接管生产环境。
   - 扩展阅读：[(译) Cloudflare 边缘网络架构：无处不在的 BPF（2019）](http://arthurchiao.art/blog/cloudflare-arch-and-bpf-zh/)

> 译者注：基于 XDP/eBPF 的 L4LB 原理都是类似的，简单来说，
>
> 1. 通过 BGP 宣告 VIP
> 2. 通过 ECMP 做物理链路高可用
> 3. 通过 XDP/eBPF 代码做重定向，将请求转发到后端（VIP -> Backend）

对此感兴趣可参考入门级介绍：[L4LB for Kubernetes: Theory and Practice with Cilium+BGP+ECMP](http://arthurchiao.art/blog/k8s-l4lb/)

## 2017 ~ 2018

### eBPF 成为内核独立子系统

随着 eBPF 社区的发展，feature 和 patch 越来越多，为了管理这些 patch，Alexei、我和 networking 的一位 maintainer David Miller 经过讨论，决定将 eBPF 作为独立的内核子 系统。

- eBPF patch 合并到 `bpf` & `bpf-next` kernel trees on git.kernel.org
- 拆分 eBPF 邮件列表：`bpf@vger.kernel.org` (archive at: `lore.kernel.org/bpf/`)
- eBPF PR 经内核网络部分的 maintainer David S. Miller 提交给 Linus Torvalds

### kTLS & eBPF

> kTLS & eBPF for introspection and ability for in-kernel TLS policy enforcement

kTLS 是**将 TLS 处理 offload 到内核**，例如，将加解密过程从 openssl 下放到内核进 行，以**使得内核具备更强的可观测性**（gain visibility）。
有了 kTLS，就可以用 eBPF 查看数据和状态，在内核应用安全策略。 **目前 openssl 已经完全原生支持这个功能**。

### bpftool & libbpf

为了检查内核内 eBPF 的状态（introspection）、查看内核加载了哪些 BPF 程序等， 我们添加了一个新工具 bpftool。现在这个工具已经功能非常强大了。
同样，为了方便用户空间应用使用 eBPF，我们提供了**用户空间 API**（user space API for applications） `libbpf`。这是一个 C 库，接管了所有加载工作，这样用户就不需要 自己处理复杂的加载过程了。

### BPF to BPF function calls

增加了一个 BPF 函数调用另一个 BPF 函数的支持，使得 BPF 程序的编写更加灵活。

## 2018

### Cilium 1.0 发布

这标志着 **BPF 革命之火燃烧到了 Kubernetes networking & security 领域**。
Cilium 此时支持的功能：

- K8s CNI
- Identity-based L3-L7 policy
- ClusterIP Services

### BTF（Byte Type Format）

内核添加了一个称为 BTF 的组件。这是一种元数据格式，和 DWARF 这样的 debugging data 类似。但 BTF 的 size 要小的多，而更重要的是，有史以来，**内核第一次变得可自 描述了**（self-descriptive）。什么意思？
想象一下当前正在运行中的内核，它**内置了自己的数据格式**（its own data format） 和**内部数据结构**（internal structures），你能用工具来查看这些东西（you can introspect them）。还是不太懂？这么说吧，**BTF 是后来的 “一次编译、到处运行”、 热补丁（live-patching）、BPF global data 处理等等所有这些 BPF 特性的基础**。
新的特性不断加入，它们都依赖 BTF 提供富元数据（rich metadata）这个基础。

> 更多 BTF 内容，可参考 [(译) Cilium：BPF 和 XDP 参考指南（2019）](http://arthurchiao.art/blog/cilium-bpf-xdp-reference-guide-zh/)
> 译者注

### Linux Plumbers 会议开辟 BPF/XDP 主题

这一年，Linux Plumbers 会议第一次开辟了专门讨论 BPF/XDP 的微型分会，我们 一起组织这场会议。其中，Networking Track 一半以上的议题都涉及 BPF 和 XDP 主题，因为这是一个非常振奋人心的特性，越来越多的人用它来解决实际问题。

### 新 socket 类型：AF_XDP

内核添加了一个**新 socket 类型** `AF_XDP`。它提供的能力是：**在零拷贝（ zero-copy）的前提下将包从网卡驱动送到用户空间**。

> 回忆前面的内容，数据包到达网卡后，先经过 XDP，然后才为这个包分配内存。 因此在 XDP 层直接将包送到用户态是无需拷贝的。
> 译者注

`AF_XDP` 提供的能力与 [DPDK](/docs/4.数据通信/DPDK.md) 有点类似，不过

- DPDK 需要**重写网卡驱动**，需要额外维护**用户空间的驱动代码**。
- `AF_XDP` 在**复用内核网卡驱动**的情况下，能达到与 DPDK 一样的性能。

而且由于**复用了内核基础设施，所有的网络管理工具还都是可以用的**，因此非常方便， 而 DPDK 这种 bypass 内核的方案导致绝大大部分现有工具都用不了了。

由于所有这些操作都是发生在 XDP 层的，因此它称为 `AF_XDP`。插入到这里的 BPF 代码 能直接将包送到 socket。

### bpffilter

开始了 bpffilter prototype，作用是通过用户空间驱动（userspace driver），**将 iptables 规则转换成 eBPF 代码**。

这是将 iptables 转换成 eBPF 的第一次尝试，整个过程对用户都是无感知的，其中的某些 组件现在还在用，用于在其他方面扩展内核的功能。

## 2018 ~ 2019

### bpftrace

Brendan 发布了 bpftrace 工具，作为 DTrace 2.0 for Linux。

### BPF 专著《BPF Performance Tools》

Berendan 写了一本 800 多页的 BPF 书。

### Cilium 1.6 发布

第一次支持完全干掉基于 iptables 的 kube-proxy，全部功能基于 eBPF。

> 这个版本其实是有问题的，例如 1.6 发布之后我们发现 externalIPs 的实现是有问题 ，社区在后面的版本修复了这个问题。在修复之前，还是得用 kube-proxy： <https://github.com/cilium/cilium/issues/9285>
> 译者注

### BPF live-patching

添加了一些内核新特性，例如尾调用（tail call），这使得 **eBPF 核心基础 设施第一次实现了热加载**。这个功能帮我们极大地优化了 datapath。
另一个重要功能是 BPF trampolines，这里就不展开了，感兴趣的可以搜索相关资料，我只 能说这是另一个振奋人心的技术。

### 第一次 bpfconf：受邀请才能参加的 BPF 内核专家会议

如题，这是 BPF 内核专家交换想法和讨论问题的会议。与 Linux Plumbers 会议互补。

### BPF backend 合并到 GCC

前面提到，BPF backend 很早就合并到 LLVM/Clang，现在，它终于合并到 GCC 了。 至此，**GCC 和 LLVM 这两个最主要的编译器套件都支持了 BPF backend**。
此外，BPF 开始支持有限循环（bounded loops），在此之前，是不支持循环的，以防止程 序无限执行。

## 2019 ~ 2020

### 不知疲倦的增长和 eBPF 的第三个方向：Linux security modules

- Google 贡献了 [BPF LSM](https://www.kernel.org/doc/html/latest/bpf/bpf_lsm.html)（安全），部署在了他们的数据中心服务器上。
- BPF verifier 防护 Spectre 漏洞（2018 年轰动世界的 CPU bug）：even verifying safety on speculative program paths。
- **主流云厂商开始通过 SRIOV 支持 XDP**：AWS (ena driver), Azure (hv_netvsc driver), …
- Cilium 1.8 支持基于 XDP 的 Service 负载均衡和 host network policies。
- Facebook 开发了基于 BPF 的 TCP 拥塞控制模块。
- Microsoft 基于 BPF 重写了将他们的 Windows monitoring 工具。
