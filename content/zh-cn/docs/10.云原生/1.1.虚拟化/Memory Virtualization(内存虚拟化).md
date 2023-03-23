---
title: Memory Virtualization(内存虚拟化)
---

# 概述

为了在一台机器上运行多个虚拟机，KVM 需要实现 VA（虚拟内存） -> PA（物理内存） -> MA（机器内存）之间的地址转换。虚机 OS 控制虚拟地址到客户内存物理地址的映射 （VA -> PA），但是虚机 OS 不能直接访问实际机器内存，因此 KVM 需要负责映射客户物理内存到实际机器内存 （PA -> MA）。内存也是可以 overcommit 的，即所有虚机的内存之和可以超过宿主机的物理内存。但使用时也需要充分测试，否则性能会受影响。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tqgo41/1616124379716-c4fa6e3c-3050-4d41-935c-8cd582fe465e.png)
