---
title: Virtualization
linkTitle: Virtualization
date: 2024-02-27T08:48
weight: 1
---

# 概述

> 参考：
> 
> - [Wiki，Virtualization](https://en.wikipedia.org/wiki/Virtualization)
> - [RedHat 7 虚拟化入门指南](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/virtualization_getting_started_guide/index)
> - [Redhat 8 官方对虚拟化的定义](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html/configuring_and_managing_virtualization/virtualization-in-rhel-8-an-overview_configuring-and-managing-virtualization#what-is-virtualization-in-rhel-8-virt-overview)
> - [RedHat 7 对“虚拟化性能不行”这个误区的辟谣](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/virtualization_getting_started_guide/index)
> - [Ubuntu 官方文档，虚拟化-介绍](https://ubuntu.com/server/docs/virtualization-introduction)

**Virtualization(虚拟化)** 是用于运行软件的广义的计算机术语。通常情况下，**Virtualization(虚拟化)** 体现在让单个可以运行多个操作系统，这些操作系统同时运行，而又是互相独立的。

虚拟化是云计算的基础。简单的说，虚拟化使得在一台物理的服务器上可以跑多台虚拟机，虚拟机共享物理机的 CPU、内存、IO 硬件资源，但逻辑上虚拟机之间是相互隔离的。物理机我们一般称为 **Host(宿主机)**，宿主机上面的虚拟机称为 **Guest(客户机)**。那么 Host 是如何将自己的硬件资源虚拟化，并提供给 Guest 使用的呢？这个主要是通过一个叫做 Hypervisor 的程序实现的。

## Hypervisor

参考：<https://www.redhat.com/zh/topics/virtualization/what-is-a-hypervisor>

Hypervisor 是用来创建与运行虚拟机的软件、固件或硬件。被 Hypervisor 用来运行一个或多个虚拟机的设备称为 Host Machine(宿主机)，这些虚拟机则称为 Guest Machine(客户机)。**Hypervisor 有时也被称为 Virtual Machine Monitor (虚拟机监视器，简称 VMM)**

# 虚拟化技术的分类

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ihdpea/1616124416735-5e89f29f-21cd-4fed-af5e-194227de3048.png)

根据 Hypervisor 的实现方式和所处的位置，虚拟化又分为两种：1 型虚拟化和 2 型虚拟化

1. 半虚拟化（para-virtualization）：TYPE1，也叫裸金属虚拟化比如 Vmware ESXi、Xen 等是一款类似于操作系统的 Hypervisor，直接运行在硬件之上，需要修改 Guest OS 的内核，让 VM 知道自己是虚拟机
2. 完全虚拟化（full-virtualization）：TYPE2，物理机上首先安装常规的操作系统，比如 Redhat、Ubuntu 和 Windows。Hypervisor 作为 OS 上的一个程序模块运行，并对管理虚拟机进行管理。比如 Vmware Workstation、KVM 等是一款类似于软件的 Hypervisor，运行于操作系统之上，VM 不知道自己是虚拟机
   1. BT：软件，二进制翻译。性能很差
   2. HVM：硬件，硬件辅助的虚拟化。性能很好。现阶段 KVM 主要基于硬件辅助进行虚拟化
      1. **硬件辅助全虚拟化主要使用了支持虚拟化功能的 CPU 进行支撑，CPU 可以明确的分辨出来自 GuestOS 的特权指令，并针对 GuestOS 进行特权操作，而不会影响到 HostOS。**
      2. 从更深入的层次来说，虚拟化 CPU 形成了新的 CPU 执行状态 —— `Non-Root Mode& Root Mode`。从上图中可以看见，GuestOS 运行在 Non-Root Mode 的 Ring 0 核心态中，这表明 GuestOS 能够直接执行特却指令而不再需要 _特权解除_ 和 _陷入模拟_ 机制。并且在硬件层上面紧接的就是虚拟化层的 VMM，而不需要 HostOS。这是因为在硬件辅助全虚拟化的 VMM 会以一种更具协作性的方式来实现虚拟化 —— _将虚拟化模块加载到 HostOS 的内核中_，例如：KVM，KVM 通过在 HostOS 内核中加载 **KVM Kernel Module** 来将 HostOS 转换成为一个 VMM。所以此时 VMM 可以看作是 HostOS，反之亦然。这种虚拟化方式创建的 GuestOS 知道自己是正在虚拟化模式中运行的 GuestOS，KVM 就是这样的一种虚拟化实现解决方案。
3. OS 级别虚拟化：容器级虚拟化，准确来说不能叫虚拟化了，只能叫容器技术无 Hypervisor，将用户空间分隔为多个，彼此互相隔离，每个 VM 中没有独立内核，OpenVZ、LXC(Linux container)、libcontainer 等，比如 Docker，Docker 的基础是 LXC。
4. 模拟(Emulation)：比如 QEMU，PearPC，Bochs
5. 库虚拟化：WINE
6. 应用程序虚拟化：JVM
7. 理论上 Type1 和 Typ2 之间的区别
   1. 1 型虚拟化一般对硬件虚拟化功能进行了特别优化，性能上比 2 型要高；
   2. 2 型虚拟化因为基于普通的操作系统，会比较灵活，比如支持虚拟机嵌套。嵌套意味着可以在 KVM 虚拟机中再运行 KVM。

# 虚拟化总结(云计算基础，实现云功能的灵活调度)

所谓的云计算：当一台虚拟机需要跨越多个物理机进行数据交互，比如拿来运行 VM 的物理主机不止一台，在每台物理机上按需启动既定数量的 VM，每个 VM 有多少 CPU 和 MEM，每个 VM 启动在哪个物理机上，启动 VM 需要的存储设备在什么地方，存储设备中的系统是临时安装，还是通过一个已经装好的系统模板直接使用，还有多个 VM 跨物理主机进行网络通信等等一系列工作，可以使用一个虚拟化管理工具(VM Manager)来实现，这个管理器的功能即可称为云计算。在没有这个管理器的时候，人们只能人为手工从把 VM 从一台物理机移动到另一台物理机，非常不灵活。

计算机五大部件：运算器(cpu)，控制器(cpu)，存储器(memory)，输入与输出设备(磁盘 I/O，网络 I/O)。

一般情况，VM 的 CPU 与 Memory 无法跨主机使用；但是磁盘 I/O 与网络 I/O 则可以跨主机使用。云计算的灵活性（即 VM 或者单个云计算节点挂了但是不影响数据，可以重新启动在任一一个节点等类似的功能）

磁盘 I/O 的灵活调度

所以，在启动一个 VM 的时候，分为这么几个启动步骤，模拟 CPU 和内存，模拟存储，模拟网络。当在多个 node 的虚拟化集群中创建完一个 VM 并想启动的时候，又分为两种情况：

1. 当该 VM 的虚拟存储放在某个节点上的时候，则该 VM 只能启动在该节点上，因为没有存储就没法加载系统镜像，何谈启动呢
2. 当该 VM 的虚拟存储放在虚拟化集群的后端存储服务器或者共享存储空间的时候，则该 VM 可以根据调度策略在任一节点启动,然后把该 VM 对应的虚拟存储挂载或下载到需要启动的节点上即可（这个所谓的虚拟存储，可以称为模板，每次 VM 启动的时候，都可以通过这个模板直接启动而不用重新安装系统了）

这种可以灵活调度 VM，而不让 VM 固定启动在一个虚拟机上的机制，这就是云功能的基础，用户不用关心具体运行在哪个节点上，都是由系统自动调度的。

网络 I/O 的灵活调度

同样的，在一个 VM 从 node1 移动到 node2 的时候，除了存储需要跟随移动外，还需要网络也跟随移动，移动的前提是所有 node 的网络配置是一样的，不管是隔离模型，还是路由模型，还是 nat 模型，还是桥接模型，都需要给每个 node 进行配置，但是，会有这么几个情况，

1. 一个公司，有 2 个部门，有两台物理 server，node1 最多有 4 个 VM，node2 最多有 4 个 VM，其中一个部门需要 5 台 VM，另一个部门需要 3 台 VM，而两个部门又要完全隔离，这时候可以通过对 vSwitch 进行 vlan 划分来进行隔离，。这时候就一个开源的软件应运而生，就是 Open vSwtich，简称为 OVS。
2. 普通 VLAN 只有 4096 个，对于公有云来说，该 vlan 数量远远不够，这时候，vxlan 技术应运而生
3. 每个公司有多个部门,每个部门有的需要连接公网，有的不需要连接公网,如果想隔离开两个公司，仅仅依靠虚拟交换机从二层隔离，无法隔离全面，这时候 vRouter 虚拟路由器技术应运而生，通过路由来隔离，并通过路由来访问，而这台 vRouter 就是由 linux 的 net namespace 功能来创建的

Openstack 中创建的每个 network 就相当于一个 vSwitch，创建的每个 route 就相当于一个 vRoute 即 net namespace，然后把 network 绑定到 route 上，就相当于把 vSwitch 连接到了 vRoute，所以，在绑定完成之后，会在 route 列表中看到个端口的 IP，这个 IP 就是 vSwitch 子网(在创建 vSwitch 的时候会设置一个可用的网段)中的一个 IP，就相当于交换机连到路由器后，路由器上这个端口的 IP

# 实际上，一个虚拟机就是宿主机上的一个文件，虚拟化程序可以通过这个文件来运行虚拟机
