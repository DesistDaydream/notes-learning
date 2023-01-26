---
title: OpenStack Networking 介绍
---

# 概述

> 参考：

传统的网络管理方式很大程度上依赖于管理员手工配置和维护各种网络硬件设备；而云环境下的网络已经变得非常复杂，特别是在多租户场景里，用户随时都可能需要创建、修改和删除网络，网络的连通性和隔离不已经太可能通过手工配置来保证了。

如何快速响应业务的需求对网络管理提出了更高的要求。传统的网络管理方式已经很难胜任这项工作，而“软件定义网络（software-defined networking, SDN）”所具有的灵活性和自动化优势使其成为云时代网络管理的主流。

Neutron 的设计目标是实现“网络即服务（Networking as a Service）”。为了达到这一目标，在设计上遵循了基于 SDN 实现网络虚拟化的原则，在实现上充分利用了 Linux 系统上的各种网络相关的技术。

所以！在学习和理解 OpenStack Network 的概念时，要与虚拟机的概念分开。Neutron 仅仅负责 SDN，也就是创建网络拓扑，拓扑中有多少交换机、路由器、都是怎么连接的。网络拓扑构建完毕后，再决定 VM 如何接入该 SDN 网络中，以及使用什么方式来接入。

# OpenStack Networking 介绍

官方文档：<https://docs.openstack.org/neutron/latest/admin/index.html>

名为 Neutron 的服务套件实现了 OpenStack 的网络功能。可以通过 Neutron 创建和管理其他 OpenStack 组件可以使用的 Network Object(网路对象)。例如 networks、subnets、ports 等。

## Neutron 管理的 Network 对象

openstack 网络管理着以下几个核心对象。这几个对象在 web 界面创建网络时也会用到。这些对象从上到下是包含的概念，Network 包含 subnet，subnet 里包含 port

1. Network # 网络。是一个隔离的二层广播域，不同的 network 之间在二层上是隔离的。network 必须属于某个 Project(有时候也叫租户 tenant)。network 支持多种类型每种网络类型，由 ML2 中的 Type Drivers 管理。
   1. local # local 网络中的实例智能与同一节点上同一网络中的实例通信
   2. flat # flat 网络中的实例能与位于同一网络的实例通信，并且可以跨越多个节点
   3. vlan # vlan 是一个二层的广播域，同一 vlan 中的实例可以通信，不同 vlan 种的实例需要通过 router 通信。vlan 中的实例可以跨节点通信、
   4. vxlan # 比 vlan 更好的技术
   5. gre # 与 vxlan 类似的一种 overlay 网络。主要区别在于使用 IP 包而非 UDP 进行封装。
2. SubNet # 子网。是一个 IPv4 或者 IPv6 地址段，创建的 instance 就从 subnet 中获取 IP，每个 subnet 需要定义 IP 地址范围和掩码
   1. 注意：
      1. 在不同 network 中的 subnet 可以一样
      2. 在相同 network 中的 subnet 不可以一样
      3. DHCP # 子网中可以创建 DHCP 服务，当启用 DHCP 服务时，会创建一个 tap 设备连接到某 bridge 上，来与子网所在的网络通信
3. Port # 端口 可以当做虚拟交换机的一个端口，创建 port 时，会给 port 分配 MAC 和 IP，当 instance 绑定到 port 时，会自动在 instance 上创建一个网卡，并获取 port 的 MAC 和 IP。如果不启用 DHCP 服务，则仅能获取 MAC，而无法获取 IP，instace 中的网卡 ip 还需要手动添加
   1. 注意：openstack 创建的 instance 本身并没有网卡，instance 中的网卡是通过 neutron 来添加的，而添加方式就是绑定某个 network 中 subnet 里的 port。绑定成功后，在 instance 中即可看到网卡设备。

Project，Network，Subnet，Port 和 VIF 之间关系。(VIF 指的是 instance 的网卡)

Project 1 : m Network 1 : m Subnet 1 : m Port 1 : 1 VIF m : 1 Instance

注意：

1. 上述核心对象是由 ML2 Plugin 负责管理。详见：Neutron 架构。所以，在这个文章中讲到的 ML2 规范，必须要求接入的插件至少可以满足对这些核心对象的管理。否则网络都无法使用。
2. 在这里，network、subnet 都是非常抽象的概念。仅仅创建一个 network 的话，在某些时候(比如 OVS 模式)，不一定会创建出来网桥设备。所有的虚拟网络设备都是在需要的时候才会自动创建出来。

OpenStack 网络分为两种类型：

1. Provider Network # 提供者网络
   1. Provider Network 由管理员进行创建。用来隔离各租户网络
2. Project Network # 项目网络
   1. Project Network 由租户自己创建。用来隔离租户自己的网络

作为网络创建过程的一部分，可以在项目之间共享所有这些类型的网络。特别是，OpenStack Networking 支持具有多个专用网络的每个项目，并使项目能够选择自己的 IP 寻址方案，即使这些 IP 地址与其他项目使用的 IP 地址重叠也是如此。

比如我现在管理着云计算平台，并提供服务。租户 1，创建账号，并申请 192.168.0.0/24 网段。租户 1，创建账号，也申请 192.168.0.0/24 网段。这时候，就需要将租户之间的网络隔离，否则无法正常通信。

同理，租户自己也有可能会出现隔离网段的需求，比如 租户 1 想让自己公司的各研发组，使用相同的网段，但是又互不影响，这时候，也需要进行隔离。

这种隔离的特性，也正是 SDN 的特点之一。其实，SDN 中的隔离，主要是靠 VLAN 来实现的，由于云计算的迅猛发展，用户不断增多，VLAN 上限 4000 的数量已经无法满足各云计算提供商了，毕竟十万、百万的用户都需要隔离。所以这时候需要使用 VxLAN

## Provider Network

官方文档：<https://docs.openstack.org/neutron/latest/admin/intro-os-networking.html#provider-networks>

## Project Network

官方文档：<https://docs.openstack.org/neutron/latest/admin/intro-os-networking.html#self-service-networks>

# Neutron 的功能

Neutron 为整个 OpenStack 环境提供网络支持，包括二层交换，三层路由，负载均衡，防火墙和 VPN 等。Neutron 提供了一个灵活的框架，通过配置，无论是开源还是商业软件都可以被用来实现这些功能。

二层交换 Switching

Nova 的 Instance 是通过虚拟交换机连接到虚拟二层网络的。Neutron 支持多种虚拟交换机，包括 Linux 原生的 Linux Bridge 和 Open vSwitch。 Open vSwitch（OVS）是一个开源的虚拟交换机，它支持标准的管理接口和协议。

利用 Linux Bridge 和 OVS，Neutron 除了可以创建传统的 VLAN 网络，还可以创建基于隧道技术的 Overlay 网络，比如 VxLAN 和 GRE（Linux Bridge 目前只支持 VxLAN）。在后面章节我们会学习如何使用和配置 Linux Bridge 和 Open vSwitch。

三层路由 Routing

Instance 可以配置不同网段的 IP，Neutron 的 router（虚拟路由器）实现 instance 跨网段通信。router 通过 IP forwarding，iptables 等技术来实现路由和 NAT。我们将在后面章节讨论如何在 Neutron 中配置 router 来实现 instance 之间，以及与外部网络的通信。

负载均衡 Load Balancing

Openstack 在 Grizzly 版本第一次引入了 Load-Balancing-as-a-Service（LBaaS），提供了将负载分发到多个 instance 的能力。LBaaS 支持多种负载均衡产品和方案，不同的实现以 Plugin 的形式集成到 Neutron，目前默认的 Plugin 是 HAProxy。我们会在后面章节学习 LBaaS 的使用和配置。

防火墙 Firewalling

Neutron 通过下面两种方式来保障 instance 和网络的安全性。

Security Group

通过 iptables 限制进出 instance 的网络包。

Firewall-as-a-Service

FWaaS，限制进出虚拟路由器的网络包，也是通过 iptables 实现。
