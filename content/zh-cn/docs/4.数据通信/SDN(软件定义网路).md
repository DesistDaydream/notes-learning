---
title: SDN(软件定义网路)
---

# 概述

> 参考：
>
> - [Wiki, Software-defined networking](https://en.wikipedia.org/wiki/Software-defined_networking)

**Software Defined Networking(软件定义网络，简称 SDN)** 技术是一种[网络管理](https://en.wikipedia.org/wiki/Network_management)方法，它支持动态的、以编程方式高效的网络配置，以提高网络性能和监控，使其更像[云计算，而](https://en.wikipedia.org/wiki/Cloud_computing)不是传统的网络管理。[\[1\]](https://en.wikipedia.org/wiki/Software-defined_networking#cite_note-ReferenceA-1) SDN 旨在解决传统网络的静态架构分散且复杂的事实，而当前网络需要更多的灵活性和易于故障排除。SDN 试图通过将[网络数据包](https://en.wikipedia.org/wiki/Network_packet)的转发过程（数据平面）与路由过程（控制平面）分离，将网络智能集中在一个网络组件中。[\[2\]](https://en.wikipedia.org/wiki/Software-defined_networking#cite_note-2)该[控制平面](https://en.wikipedia.org/wiki/Control_plane)由一个或多个控制器组成，这些控制器被认为是包含整个智能的 SDN 网络的大脑。然而，智能中心化在安全性、[\[1\]](https://en.wikipedia.org/wiki/Software-defined_networking#cite_note-ReferenceA-1)可扩展性和弹性[\[1\]](https://en.wikipedia.org/wiki/Software-defined_networking#cite_note-ReferenceA-1)方面有其自身的缺点，这是 SDN 的主要问题。

自[OpenFlow](https://en.wikipedia.org/wiki/OpenFlow)协议于 2011 年出现以来，SDN 通常与[OpenFlow](https://en.wikipedia.org/wiki/OpenFlow)协议（用于与网络平面元素进行远程通信，以确定[网络数据包](https://en.wikipedia.org/wiki/Network_packet)通过[网络交换机](https://en.wikipedia.org/wiki/Network_switch)的路径）相关联。然而，自 2012 年[\[3\] ](https://en.wikipedia.org/wiki/Software-defined_networking#cite_note-TechTarget:_SDN_is_not_OpenFlow-3)[\[4\]](https://en.wikipedia.org/wiki/Software-defined_networking#cite_note-TechTarget:_OpenFlow_not_the_only_show_in_town-4) OpenFlow 对于许多公司不再是独家解决方案，他们增加了专有技术。其中包括[Cisco Systems](https://en.wikipedia.org/wiki/Cisco_Systems)的开放网络环境和[Nicira](https://en.wikipedia.org/wiki/Nicira)的[网络虚拟化平台](https://en.wikipedia.org/wiki/Network_virtualization_platform)。
[SD-WAN](https://en.wikipedia.org/wiki/SD-WAN)将类似技术应用于[广域网](https://en.wikipedia.org/wiki/Wide_area_network)(WAN)。[\[5\]](https://en.wikipedia.org/wiki/Software-defined_networking#cite_note-5)

SDN 技术目前可用于需要极快故障转移的工业控制应用，称为操作技术 (OT) 软件定义网络 (SDN)。OT SDN 技术是一种在关键基础设施网络的环境强化硬件上管理网络访问控制和以太网数据包交付的方法。OT SDN 将控制平面的管理从集中在流控制器中的交换机抽象出来，并将 SDN 应用为交换机中的底层控制平面。去除了传统控制平面，简化了交换机，同时集中控制平面管理。OT SDN 中使用的通用控制平面标准是 OpenFlow，使其可与其他 SDN 解决方案互操作，不同之处在于 OpenFlow 是交换机中唯一的控制平面，并且交换机在电源循环期间保留流量，并且所有流量和冗余都经过主动流量工程设计因此交换机可以执行转发，它们被配置为在有或没有在线流量控制器的情况下执行。OT SDN 在性能、网络安全和态势感知方面为工业网络提供了优势。性能优势是通过使用 OpenFlow 中的快速故障转移组的主动流量工程意外事件实现的，从而在微秒内从链路或交换机故障中恢复网络，而不是像生成树技术那样的毫秒级。另一个性能优势是环路缓解是通过流量工程路径规划完成的，而不是阻塞端口，允许系统所有者主动使用所有端口。OT SDN 的网络安全优势在于交换机默认拒绝，流是允许流量转发的规则。这提供了强大的网络访问控制，可以在每一跳从 OSI 模型的第 1 层到第 4 层检查数据包。由于旧控制平面不再存在，因此移除了旧控制平面安全漏洞。MAC 表欺骗和 BPDU 欺骗不再可能，因为两者都不存在于 OT SDN 交换机中。旋转和网络侦察不再适用于适当的流编程，因为仅允许转发结合物理位置和路径与虚拟数据包过滤的流量。OT SDN 的态势感知优势使网络所有者能够了解其网络上有哪些设备，哪些对话可以和正在发生，以及这些对话可以在谁之间发生。OT SDN 网络技术允许以太网满足关键基础设施测量和控制的苛刻通信消息交换要求，并简单地为系统所有者提供对哪些设备可以连接到网络、这些设备可以连接到哪里以及每个设备可以进行哪些对话的控制有。OT SDN 的态势感知优势使网络所有者能够了解其网络上有哪些设备，哪些对话可以和正在发生，以及这些对话可以在谁之间发生。OT SDN 网络技术允许以太网满足关键基础设施测量和控制的苛刻通信消息交换要求，并简单地为系统所有者提供对哪些设备可以连接到网络、这些设备可以连接到哪里以及每个设备可以进行哪些对话的控制有。OT SDN 的态势感知优势使网络所有者能够了解其网络上有哪些设备，哪些对话可以和正在发生，以及这些对话可以在谁之间发生。OT SDN 网络技术允许以太网满足关键基础设施测量和控制的苛刻通信消息交换要求，并简单地为系统所有者提供对哪些设备可以连接到网络、这些设备可以连接到哪里以及每个设备可以进行哪些对话的控制有。

SDN 的研究仍在继续，因为许多[仿真器](https://en.wikipedia.org/wiki/Emulator)正在开发用于研究目的，例如 vSDNEmul、[\[6\]](https://en.wikipedia.org/wiki/Software-defined_networking#cite_note-6) EstiNet、[\[7\]](https://en.wikipedia.org/wiki/Software-defined_networking#cite_note-7) Mininet [\[8\]](https://en.wikipedia.org/wiki/Software-defined_networking#cite_note-8)等。
