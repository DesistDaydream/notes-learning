---
title: Neutron 架构
---

# 概述

> 参考：
> - 官方文档：<https://docs.openstack.org/neutron/latest/admin/intro-os-networking.html>

1. Neutron Server
   1. 接收请求。对外提供 OpenStack 网络 API，接收请求，并调用 Plugin 处理请求。
2. Plugins 插件/Agent 代理
   1. 实现请求。
   2. 实现 OpenStack 网络的主要组件。用来创建各种网络设备和配置规则
   3. Plugins 用来处理 Neutron Server 发来的请求，维护 OpenStack 逻辑网络状态， 并调用 Agent 处理请求。
   4. Agent 用来处理对应 Plugin 的请求，并在宿主机上创建相应的网络设备以及生成网络规则。
   5. Plugins 与 Agent 一般都是配套使用。比如 OVS Plugin 需要 OVS Agent。
3. Queue 队列
   1. 组件间通信。Neutron Server，Plugin 和 Agent 之间通过 Messaging Queue 通信和调用。
4. Database 数据库
   1. 保存网络状态。接收 Plugins 的信息，保存 OpenStack 的网络状态信息，包括 Network, Subnet, Port, Router 等。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/isxvn1/1616123261056-8e665764-b4d5-4d71-9756-af5d1fde6121.jpeg)

额外说明：

- plugin 解决的是 What 的问题，即网络要配置成什么样子？而至于如何配置 How 的工作则交由 agent 完成。
- plugin 的一个主要的职责是在数据库中维护 Neutron 网络的状态信息，这就造成一个问题：所有 network provider 的 plugin 都要编写一套非常类似的数据库访问代码。为了解决这个问题，Neutron 在 H 版本实现了一个 ML2（Modular Layer 2）plugin，对 plugin 的功能进行抽象和封装。有了 ML2 plugin，各种 network provider 无需开发自己的 plugin，只需要针对 ML2 开发相应的 driver 就可以了，工作量和难度都大大减少。

# Neutron Server

Neutron Server 提供了一个公开 Neutron API 的 Web 服务器，并将所有 Web 服务调用传递给 Neutron 插件进行处理。

# Plugins 插件

plugin 按照功能分为两类：core plugin 和 service plugin。core plugin 维护 Neutron 的 netowrk, subnet 和 port 相关资源的信息，与 core plugin 对应的 agent 包括 linux bridge, OVS 等；service plugin 提供 routing, firewall，load balance 等服务，也有相应的 agent。

## Core Plugins 核心插件

官方文档：<https://docs.openstack.org/neutron/latest/admin/config-ml2.html>

ML2(Modular Layer 2) 是 Neutron 的一个 Core Plugin(核心插件)，该插件提供了一个框架，可也可以称为 ML2 Framework(ML2 框架)。

ML2 框架允许在 OpenStack 网络中可以使用多种 Layer 2(二层) 网络技术(只需要更改配置文件即可)，不同的节点可以使用不同的网络实现机制。这就好像是一种规范，只要符合 ML2 规范，都可以作为插件，接入到 Openstack 中，提供网络服务。跟 K8S 的 CNI 有异曲同工之妙。

ML2 对 二层网络 进行抽象，引入了 type driver 和 mechansim driver 这两个概念。

1. Type Drivers # OpenStack 网络类型驱动。定义底层如何实现 OpenStack 网络。比如 VXLAN、Flat 等
   1. 就是 OpenStack 的 Nework Type，详见：OpenStack Networking 介绍 中的基本概念。不同类型的驱动用来维护该类型的网络。
2. Mechanism Drivers # OpenStack 网络机制驱动。定义访问某种 OpenStack 网络类型的机制。比如 Open vSwitch、Linux Bridge 等。
   1. 就是 Plugins 本身，不同机制的插件，可以管理的 OpenStack Nework Type 不同。详情见下：
      1. Open vSwitch 支持：Flat、VLAN、VXLAN、GRE
      2. Linux Bridge 支持：Flat、VLAN、VXLAN
      3. SRIOV 支持：Flat、VLAN
      4. MacVTap 支持：Flat、VLAN
      5. L2 Population 支持：VXLAN、GRE
   2. Mechanism Drivers 可以利用 L2 Agent(通过 RPC 调用)与设备或控制器进行交互，也可以直接与设备或控制器进行交互。

只要实现了上述两个概念的插件，皆可接入 ML2，为 OpenStack 提供网络服务。

### Network Type Drivers

在 ML2 插件中启用指定的 Type Drivers。编辑 /etc/neutron/plugins/ml2/ml2_conf.ini 文件，示例：

                [ml2]type_drivers = flat,vlan,vxlan

更多配置信息参考 Networking configuration options

可以使用以下类型的网络驱动

- local # 本地网络。无法与宿主机外部胡同
- Flat # 平面网络。无法为网络添加 vlan 标签
- VLAN # VLAN 网络。可以为网络中任何一个端口添加 vlan 标签
- GRE
- VXLAN
- geneve

注意：Provider Network 与 Project Network 可用的 Type Drivers 是不同的，在 Project Network 中是无法使用 flat 这个 Type Driver，因为 flat 类型的网络与物理网卡一一对应

### Mechanism Drivers

要在 ML2 插件中启用 Mechanism Drivers，在 neutron 服务器上编辑 /etc/neutron/plugins/ml2/ml2_conf.ini 文件，示例：

                [ml2]mechanism_drivers = linuxbridge

可以使用以下机制的网络驱动，更多配置参考： Configuration Reference.

1. Linux Bridge
   1. 这个 Mechanism Driver 不需要其他配置。但是需要代理配置。有关详细信息，请参阅下面的 L2 Agent 相关部分。
2. Open vSwitch
   1. 这个 Mechanism Driver 不需要其他配置。但是需要代理配置。有关详细信息，请参阅下面的 L2 Agent 相关部分。
3. SRIOV
   1. SRIOV 驱动程序接受所有 PCI 供应商设备。
4. MacVTap
   1. 这个 Mechanism Driver 不需要其他配置。但是需要代理配置。请参阅相关部分。
5. L2 population
   1. 管理员可以配置一些可选的配置选项。有关更多详细信息，请参阅《配置参考》中的相关部分。
6. Specialized
   1. 开源的
      1. 存在外部开源机制驱动程序以及中子集成参考实现。这些驱动程序的配置不是本文档的一部分。例如：
      2. OpenDaylight
      3. OpenContrail
   2. 专有（供应商）
7. 存在来自各种供应商的外部机制驱动程序以及中子集成参考实现。

## Service Plugins 服务插件

Core [Plugin/Agent 负责管理核心实体](http://mp.weixin.qq.com/s?__biz=MzIwMTM5MjUwMg==&mid=2653587219&idx=1&sn=e14476e223c7a2743ce9efacdf2f020c&scene=21#wechat_redirect)：net, subnet 和 port。而对于更高级的网络服务，则由 Service Plugin 管理。

Service Plugin 及其 Agent 提供更丰富的扩展功能，包括 Virtual Routers(虚拟路由)，load balance，firewall 等

# Agent

官方文档：<https://docs.openstack.org/neutron/latest/admin/config-ml2.html#agents>

Agent 用来处理对应 Plugin 的请求，并在宿主机上创建相应的网络设备以及生成网络规则。提供与 instances 的 2 层和 3 层 连接。处理物理网络—虚拟网络的过渡。处理元数据。等等。

每种 Plugins 都有其对应的 Agent 来处理请求。

在 /etc/neutron/plugins/ml2/ml2_conf.ini 配置文件中配置了指定的 mechanism driver 后。启动 Neutron 服务，则所有节点上都会运行配置中指定 driver 的对应 Agent

比如我指定 mechanism driver 为 linuxbridge

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/isxvn1/1616123261031-37e35b59-a381-4c74-87ad-705aafc88a0d.jpeg)

那么就会启动一个 linuxbridge 的 agent 进程。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/isxvn1/1616123261049-0a087308-40d4-40dc-8342-4c8887c15ba2.jpeg)

下面是各种类型的 Agent 的介绍

1. L2 Agent
   1. L2 Agent 提供 2 层网络。可用的 L2 Agent 有 Linux Bridge、OVS 等。
2. L3 Agent
   1. L3 Agent 提供高级的三层网络功能，比如 Virtual Routers(虚拟路由)、Floating IPs(弹性 IP)等等。L3 Agent 依赖并行运行的 L2 Agent。
   2. L3 agent 需要正确配置才能工作，配置文件为 /etc/neutron/l3_agent.ini，位于控制节点或网络节点上。
      1. interface_driver 是最重要的选项，
         1. 如果 mechanism driver 是 linux bridge，则：
            1. interface_driver = neutron.agent.linux.interface.BridgeInterfaceDriver
         2. 如果选用 open vswitch，则：
            1. interface_driver = neutron.agent.linux.interface.OVSInterfaceDriver
3. DHCP Agent
   1. DHCP Agent 负责 DHCP 和 RADVD 服务。它需要在同一节点上运行的 L2 代理。
4. Metadata Agent
   1. instance 在启动时需要访问 nova-metadata-api 服务获取 metadata 和 userdata，这些 data 是该 instance 的定制化信息，比如 hostname, ip， public key 等。
   2. 但 instance 启动时并没有 ip，那如何通过网络访问到 nova-metadata-api 服务呢？
   3. 答案就是 neutron-metadata-agent 该 agent 让 instance 能够通过 dhcp-agent 或者 l3-agent 与 nova-metadata-api 通信
5. L3 metering Agent
   1. L3 metering Agent 启用第 3 层流量计量。它需要在同一节点上运行的 L3 代理。
6. Security
   1. L3 agent 可以在 router 上配置防火墙策略，提供网络安全防护。另一个与安全相关的功能是 Security Group，也是通过 IPtables 实现。 Firewall 与 Security Group 的区别在于：
      1. Firewall 安全策略位于 router，保护的是某个 project 的所有 network。
      2. Security Group 安全策略位于 instance，保护的是单个 instance。

# Database 数据库

# Reference Implementations 参考实现

mechanism driver 和 L2 agent 的组合称为“参考实现”。下表列出了这些实现：

|                  |                                        |
| ---------------- | -------------------------------------- |
| Mechanism Driver | L2 agent                               |
| Open vSwitch     | Open vSwitch agent                     |
| Linux bridge     | Linux bridge agent                     |
| SRIOV            | SRIOV nic switch agent                 |
| MacVTap          | MacVTap agent                          |
| L2 population    | Open vSwitch agent, Linux bridge agent |

下表显示了哪些参考实现支持哪些非 L2 的 Agent

|                                   |          |            |                |                   |
| --------------------------------- | -------- | ---------- | -------------- | ----------------- |
| 参考实现                          | L3 agent | DHCP agent | Metadata agent | L3 Metering agent |
| Open vSwitch & Open vSwitch agent | yes      | yes        | yes            | yes               |
| Linux bridge & Linux bridge agent | yes      | yes        | yes            | yes               |
| SRIOV & SRIOV nic switch agent    | no       | no         | no             | no                |
| MacVTap & MacVTap agent           | no       | no         | no             | no                |
