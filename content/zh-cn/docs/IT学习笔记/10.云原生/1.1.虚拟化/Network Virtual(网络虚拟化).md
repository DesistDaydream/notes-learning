---
title: Network Virtual(网络虚拟化)
---

# 概述

> ## 参考：

在传统网络环境中，一台物理主机包含一个或多个网卡（NIC），要实现与其他物理主机之间的通信，需要通过自身的 NIC 连接到外部的网络设施，如交换机上；为了对应用进行隔离，往往是将一个应用部署在一台物理设备上，这样会存在两个问题：

1. 是某些应用大部分情况可能处于空闲状态，
2. 是当应用增多的时 候，只能通过增加物理设备来解决扩展性问题。不管怎么样，这种架构都会对物理资源造成极大的浪费。

为了解决这个问题，可以借助虚拟化技术对一台物理资源进行抽象，将一张物理网卡虚拟成多张虚拟网卡（vNIC），通过虚拟机来隔离不同的应用。

1. 针对问题 1），可以利用虚拟化层 Hypervisor 的调度技术，将资源从空闲的应用上调度到繁忙的应用上，达到资源的合理利用；
2. 针对问题 2），可以根据物理设备的资源使用情况进行横向扩容，除非设备资源已经用尽，否则没有必要新增设备。

综上所述：SDN 主要是通过系统的功能，模拟出网络设备中的路由器，交换机，端口，网线等等，这些现实中的数通设备都可以通过软件来模拟实现

网络虚拟化的几种最基础模型：

1. 隔离模型：在 host 上创建一个 vSwitch(bridge device)：每个 VM 的 TAP 设备直接添加至 vswitch 上，VM 通过 vSwitch 互相通信，与外界隔离
2. 路由模型：基于隔离模型，在 vSwitch 添加一个端口，作为 host 上的虚拟网卡使用(就是 VMware workstation 中创建的那些虚拟网卡，其中的 IP 作为虚拟机的网关)，并打开 host 的核心转发功能，使数据从 VM 发送到 host；该模型数据包可以从 VM 上出去，但是外界无法回到 VM，如果想让外部访问 VM，需要添加 NAT 功能，变成 NAT 模型
3. NAT 模型：配置 Linux 自带的 NAT(可通过 iptables 定义)功能，所有 VM 的 IP 被 NAT 成物理网卡 IP，这是一种常用的虚拟网络模型
4. 桥接模型：可以想象成把物理机网卡变成一台 vSwitch，然后给物理机创建一个虚拟网卡，虚拟机和物理机都连接到 vSwitch，相当于把虚拟机直接接入到网络中，从网络角度看，VM 相当于同网段的一台 host
5. 隧道模型：VM 的数据包在经过某个具备隧道功能的虚拟网络设备时，可以在数据包外层再封装一层 IP，以 IP 套 IP 的隧道方式，与对方互通

网络虚拟化术语

1. Network Stack：网络栈，包括网卡（Network Interface）、回环设备（LoopbackDevice）、路由表（Routing Table）和 iptables 规则。对于一个进程来说，这些要素，其实就构成了它发起和响应网络请求的基本环境。
2. port：当成虚拟交换机上的端口，可以打 VLAN TAG
3. interface：当成接口，类似于连到端口的网线，可以设置 TYPE。[注意端口与接口的概念](https://blog.csdn.net/number1killer/article/details/79226772)
4. bridge，port，interface：一个 BRIDGE 上可以配置多个 PORT，一个 PORT 上可以配置多个 INTERFACE

# 网路虚拟化的解决方案

基于 Linux 本身的网络虚拟化方案

1. Linux Bridge # 虚拟网络基础资源，用于二层网络虚拟化
2. Namespace # 网络名称空间，用于三层网络虚拟化

高级虚拟化方案

1. Open vSwitch # 开源的虚拟交换机，用于二层网络虚拟化

## Linux Bridge(Linux 网桥)

> 参考：
> - [Linux 上抽象网络设备的原理及使用](https://www.teambition.com/project/5f90e312755d8a00446050eb/app/5eba5fba6a92214d420a3219/workspaces/5f90e312c800160016ea22fb/docs/5f9d2b14eaa119000192206f)
> - [云计算底层技术-虚拟网络设备(Bridge,VLAN)](https://www.teambition.com/project/5f90e312755d8a00446050eb/app/5eba5fba6a92214d420a3219/workspaces/5f90e312c800160016ea22fb/docs/5f9d2a2712d5ba000162d3e2)
> - [云计算底层技术-虚拟网络设备(tun tap,veth)](https://www.teambition.com/project/5f90e312755d8a00446050eb/app/5eba5fba6a92214d420a3219/workspaces/5f90e312c800160016ea22fb/docs/5f9d294eeaa11900019215fc)

每台 VM 都有一套独立的网络栈，如果想让两台 VM 互相通信怎么办呢？最直接的办法就是把两台 VM 用一根网线连接起来，而如果想要实现多台 VM 通信，那就需要用把它们连接到一台交换机上。这个**简单的虚拟交换机**的功能就叫 `Bridge(网桥)`，而**虚拟交换机想要实现普通物理交换机的**上的功能，比如接口，网线，端口，Vlan 等 都有专门对应的虚拟网络设备来模拟实现。而交换机中的 VLAN 则是 Linux 自己本身实现的。

### 虚拟网络设备的类型

**Veth Pair：**(Virtual Ethernet Pair)虚拟以太网设备对，可以理解为一跟网线。是 Linux 内核由模块实现的虚拟网络设备，该设备被创建出来后，总是以两张虚拟网卡（Veth Peer）的形式成对出现的。并且，从其中一个“网卡”发出的数据包，可以直接出现在与它对应的另一张“网卡”上，哪怕这两个“网卡”在不同的 Network Namespace 里。这就使得 Veth Pair 常常被用作连接不同 Network Namespace 的“网线”。创建方式详见 ip link add 命令。可以把 Veth 连接到 vSwich 与 vSwitch、vSwitch 与 nameSpace、namespace 与 namespace 上。

veth 的作用是反转数据流量，从一段接收到数据后，会把该数据流进行反转变成发送，发送到对端后，对端接收到数据流之后，再次反转会发送给绑定到该端的 namesapce 中。ethtool -S VethNAME 可以使用该命令通过 Veth 的一半查看另一半网卡的序号

**TAP/TUN**：该设备会创建一个 /dev/tunX 的文件并作用在内核空间，与用户空间的 APP 相连(比如 VM)，当这个 VM 通过其 Hypervisor 通信时，会把数据写入该/dev/tunX 文件，并交给内核，内核会处理这些数据并通过网卡发送。TAP 工作在二层，TUN 工作在三层。详见 TUN TAP 设备浅析(一) -- 原理浅析

- Veth Pari 与 TAP/TUN 设备在 VM 与 Container 中的使用注意事项以及原因
  - 为什么 VM 要使用 TAP/TUN,而 Container 不用？因为 VM 数据包在从其进程发送到 Host 的时候，由于 VM 有自己的内核，那么这个数据包相当于已经经过了一个 VM 的网络栈，这时候就不能直接发送给 Host 的网络栈再进行处理了，所以需要 TAP/TUN 设备来作为一个转折点，接收 VM 的数据包，并以一个已经经过网络栈处理过的姿态直接进入内核的网络设备。

**Bridge(网桥)**：在 Linux 中能够起到虚拟交换机作用的网络设备，但不同于 TAP/TUN 这种单端口的设备，Bridge 实现虚拟为多端口，本质上是一个虚拟交换机，具备和物理交换机类似的功能。Bridge 可以绑定其他 Linux 网络设备作为从设备，并将这些从设备虚拟化为端口，当一个从设备被绑定到 Bridge 上时，就相当于真实网络中的交换机端口上插入了一根连有终端的网线。

注意：一旦一块虚拟网卡被连接到 Bridge 上，该设备会变成 该 Bridge 的“从设备”。从设备会被剥夺调用网络协议栈处理数据包的资格，从而降级成网桥上的一个端口。而这个端口的唯一作用就是接收流入的数据包，然后把这些数据包的“生杀大权”(比如转发或者丢弃等)全部交给对应的 Bridge 进行处理

### Linux 如何实现 VLAN

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kinqyh/1616124322231-b84ec223-407c-41fd-b0bf-0a1c61faa7c9.jpeg)

1. eth0 是宿主机上的物理网卡，有一个命名为 eth0.10 的子设备与之相连。 eth0.10 就是 VLAN 设备了，其 VLAN ID 就是 VLAN 10。 eth0.10 挂在命名为 brvlan10 的 Linux Bridge 上，虚机 VM1 的虚拟网卡 vent0 也挂在 brvlan10 上。
2. 这样的配置其效果就是： 宿主机用软件实现了一个交换机（当然是虚拟的），上面定义了一个 VLAN10。 eth0.10，brvlan10 和 vnet0 都分别接到 VLAN10 的 Access 口上。而 eth0 就是一个 Trunk 口。VM1 通过 vnet0 发出来的数据包会被打上 VLAN10 的标签。
3. eth0.10 的作用是：定义了 VLAN10
4. brvlan10 的作用是：Bridge 上的其他网络设备自动加入到 VLAN10 中
5. 再增加一个 VLAN20
6. 样虚拟交换机就有两个 VLAN 了，VM1 和 VM2 分别属于 VLAN10 和 VLAN20。 对于新创建的虚机，只需要将其虚拟网卡放入相应的 Bridge，就能控制其所属的 VLAN。
7. VLAN 设备总是以母子关系出现，母子设备之间是一对多的关系。 一个母设备（eth0）可以有多个子设备（eth0.10，eth0.20 ……），而一个子设备只有一个母设备。

Linux Bridge + VLAN = 虚拟交换机

1. 物理交换机存在多个 VLAN，每个 VLAN 拥有多个端口。 同一 VLAN 端口之间可以交换转发，不同 VLAN 端口之间隔离。 所以交换机其包含两层功能：交换与隔离。
2. Linux 的 VLAN 设备实现的是隔离功能，但没有交换功能。 一个 VLAN 母设备（比如 eth0）不能拥有两个相同 ID 的 VLAN 子设备，因此也就不可能出现数据交换情况。
3. Linux Bridge 专门实现交换功能。 将同一 VLAN 的子设备都挂载到一个 Bridge 上，设备之间就可以交换数据了。

总结起来，Linux Bridge 加 VLAN 在功能层面完整模拟现实世界里的二层交换机。eth0 相当于虚拟交换机上的 trunk 口，允许 vlan10 和 vlan20 的数据通过。

eth0.10，vent0 和 brvlan10 都可以看着 vlan10 的 access 口。

eth0.20，vent1 和 brvlan20 都可以看着 vlan20 的 access 口。

## Open vSwitch # 开放的虚拟交换机

详见 [Open vSwitch](https://www.yuque.com/go/doc/33175901)

## Network Namespace # 网络名称空间

Network Namespace 可以简单得理解为 Linux 上的 **虚拟路由器(vRouter)。**详见：[Network Namespace 详解](https://www.yuque.com/go/doc/33173252)

Network Namespace 在逻辑上是网络栈的另一个副本，具有自己的路由、防火墙规则、网络设备等功能。默认情况下，每个进程从其父进程继承其网络名称空间。在最初的时候，所有进程都共享来自系统启动的 PID 为 1 的父进程的名称空间。整个系统 PID 为 1 的进程的 Network Namespace 就是整台机器和系统的网络栈。Linux 内核可以通过 clone()函数的功能在默认的网络名称空间中，克隆出来一个具备相同功能的网络栈，该克隆出来的 Network Namespace 为绑定上来的进程提供一个完全独立的网络协议栈，多个进程也可同时共享同一个 Network Namespace。Host 上一个网络设备只能存在于一个 Network Namespace 当中，而不同的 Network Namespace 之间要想通信，可以通过虚拟网络设备来提供一种类似于管道的抽象，在不同的 Network Namespace 中建立隧道实现通信。当一个 Network Namespace 被销毁时，插入到该 Netwrok Namespace 上的虚拟网络设备会被自动移回最开始默认整台设备的 Network Namespace

Network Namespace 的应用场景

1. 可以把 Net Namepace 就相当于在物理机上创建了一台 vRouter，这台 vRouter 就是一块 namespace，把与 VM 连接的 vSwitch 连接到这台 vRouter，然后 VM 通过 vRouter 与外部或者另一部分被隔离的网络通信，这样即可实现对这台 vSwitch 以及与之关联的 VM 进行网络隔离（如果要与外部通信，那么需要使用桥接模型，把物理网卡模拟成 vSwitch，然后把该 vSwitch 关联到该 vRouter）
2. Network Namespace 还可用于承载 Container 技术的网络功能，一个 Container 占据一个 Namespace，通过使用 Veth 设备连接 Namespace 与 Bridge 相连来实现各 Namespace 中的 Container 之间互相通信。具体详见 [Network Namespace 详解](https://www.yuque.com/go/doc/33173252)

管理 Network Namesapce 的方式：

1. 通过 ip netns 命令来管理，该命令的用法详见[ Iproute2 命令行工具](https://www.yuque.com/go/doc/33221906) 中的 netns 子命令

# Overlay Network 叠加网络

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kinqyh/1616124322243-b9ecf172-ef2b-4452-8c3c-364af95526aa.jpeg)
Overlay Network 产生的原因

在网桥的概念中，各个 VM 可以通过 Host 上的 vSwitch 来进行通信，那么当需要访问另外一台 Host 上不同网段的 VM 的时候呢？

实现 overlay network 的方式有 gre 等

VXLAN

概念详见 flannel.note 中的 vxlan 模型
