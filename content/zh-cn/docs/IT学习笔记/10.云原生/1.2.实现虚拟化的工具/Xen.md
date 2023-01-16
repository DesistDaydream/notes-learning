---
title: Xen
---

在 Xen 使用的方法中，没有指令翻译。这是通过两种方法之一实现的。第一，使用一个能理解和翻译虚拟操作系统发出的未修改指令的 CPU（此方法称作完全虚拟化或 full virtualization）。另一种，修改操作系统，从而使它发出的指令最优化，便于在虚拟化环境中执行（此方法称作准虚拟化或 paravirtualization）。

在 Xen 环境中，主要有两个组成部分。一个是虚拟机监控器（VMM），也叫 hypervisor。Hypervisor 层在硬件与虚拟机之间，是必须最先载入到硬件的第一层。Hypervisor 载入后，就可以部署虚拟机了。在 Xen 中，虚拟机叫做“domain”。在这些虚拟机中，其中一个扮演着很重要的角色，就是 domain0，具有很高的特权。通常，在任何虚拟机之前安装的操作系统才有这种特权。

Domain0 要负责一些专门的工作。由于 hypervisor 中不包含任何与硬件对话的驱动，也没有与管理员对话的接口，这些驱动就由 domain0 来提供了。通过 domain0，管理员可以利用一些 Xen 工具来创建其它虚拟机（Xen 术语叫 domainU）。这些 domainU 也叫无特权 domain。这是因为在基于 i386 的 CPU 架构中，它们绝不会享有最高优先级，只有 domain0 才可以。

在 domain0 中，还会载入一个 xend 进程。这个进程会管理所有其它虚拟机，并提供这些虚拟机控制台的访问。在创建虚拟机时，管理员使用配置程序与 domain0 直接对话。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fpcdzk/1616124139353-6dccbc0f-787a-4305-8851-45f3e17c8162.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fpcdzk/1616124139376-0a4f6b44-7eef-419c-ad71-617061d9e6d2.jpeg)

# Xen 的组成部分：

1. Xen Hypervisor：分配 CPU、Memory、Interrupt
2. Domain0：特权域，安装完 Hypervisor 后，自动生成的一个特权 Virtual Machine，负责管理整个 Xen 的 VM，可以直接访问硬件 IO 资源，修改 Linux Kernel 以实现半虚拟化功能
   1. 提供管理 DomainU 的工具栈，用于实现对虚拟机进行添加，启动，快照，停止，删除等操作
3. DomainU：非特权域，U 是各个虚拟机的编号 1,2,3,4,.....，受 Domain0 管理
   1. PV：不依赖于 CPU 的 HVM 特性，但要求 GusetOS 的内核做出修改以直销自己运行于 PV 环境
   2. HVM：依赖于 Intel VT 或 AMD AMD-V，还要依赖于 Qemu 来模拟 IO 设备
   3. PV on HVM：CPU 为 HVM 模式运行，IO 设备为 PV 模式运行

Xen 的工具栈

1. xm/xend：在 Xen Hypervisor 的 Domain0 中要启动 xend 服务
2. xm：命令行管理工具，有诸多子命令
3. xl：基于 libxenlight 提供的轻量级的命令行工具栈
4. xe/xapi：提供了对 Xen 管理的 api，因此多用于 cloud 环境中
5. virsh/libvirt：

XenStore：为各 Domain 提供的共享信息存储空间，有着层级结构的名称空间，位于 Domain0 上
