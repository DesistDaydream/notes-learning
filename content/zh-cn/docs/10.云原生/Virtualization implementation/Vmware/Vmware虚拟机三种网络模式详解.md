---
title: Vmware虚拟机三种网络模式详解
---

# 概述

vmware 为我们提供了三种网络工作模式，它们分别是：Bridged（桥接模式）、NAT（网络地址转换模式）、Host-Only（仅主机模式）。

打开 vmware 虚拟机，我们可以在选项栏的“编辑”下的“虚拟网络编辑器”中看到 VMnet0（桥接模式）、VMnet1（仅主机模式）、VMnet8（NAT 模式），那么这些都是有什么作用呢？其实，我们现在看到的 VMnet0 表示的是用于桥接模式下的虚拟交换机；VMnet1 表示的是用于仅主机模式下的虚拟交换机；VMnet8 表示的是用于 NAT 模式下的虚拟交换机。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098737-7a8bd1ed-da33-4bdd-b482-e95a4c8f7a45.jpeg)
同时，在主机上对应的有 VMware Network Adapter VMnet1 和 VMware Network Adapter VMnet8 两块虚拟网卡，它们分别作用于仅主机模式与 NAT 模式下。在“网络连接”中我们可以看到这两块虚拟网卡，如果将这两块卸载了，可以在 vmware 的“编辑”下的“虚拟网络编辑器”中点击“还原默认设置”，可重新将虚拟网卡还原。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098765-2a469c9a-504f-4da3-bbe8-4b02178b6264.jpeg)
小伙伴看到这里，肯定有疑问，为什么在真机上没有 VMware Network Adapter VMnet0 虚拟网卡呢？那么接下来，我们就一起来看一下这是为什么。

一、Bridged（桥接模式）

什么是桥接模式？桥接模式就是将主机网卡与虚拟机虚拟的网卡利用虚拟网桥进行通信。在桥接的作用下，类似于把物理主机虚拟为一个交换机，所有桥接设置的虚拟机连接到这个交换机的一个接口上，物理主机也同样插在这个交换机当中，所以所有桥接下的网卡与网卡都是交换模式的，相互可以访问而不干扰。在桥接模式下，虚拟机 ip 地址需要与主机在同一个网段，如果需要联网，则网关与 DNS 需要与主机网卡一致。其网络结构如下图所示：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098764-ba15bb94-065e-48e3-a6de-68ff4f51b149.jpeg)

接下来，我们就来实际操作，如何设置桥接模式。

首先，安装完系统之后，在开启系统之前，点击“编辑虚拟机设置”来设置网卡模式。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098816-121d68c2-fed6-42b6-9167-3262e7d6e8e2.jpeg)

点击“网络适配器”，选择“桥接模式”，然后“确定”

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098751-19705c4a-85aa-4053-b70f-80ea1d845c9e.jpeg)

在进入系统之前，我们先确认一下主机的 ip 地址、网关、DNS 等信息。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098777-4efa0137-e475-46d6-ae38-37c6ebdf0d0a.jpeg)

然后，进入系统编辑网卡配置文件，命令为 vi /etc/sysconfig/network-scripts/ifcfg-eth0

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098745-c54868f9-0bf4-4789-9865-9845099db9b3.jpeg)

添加内容如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098736-d44d13d1-dcb7-487f-8b33-2e432f83c795.jpeg)

编辑完成，保存退出，然后重启虚拟机网卡，使用 ping 命令 ping 外网 ip，测试能否联网。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098775-e91b9d9e-7109-4acd-99e3-9f52d808336c.jpeg)

能 ping 通外网 ip，证明桥接模式设置成功。

那主机与虚拟机之间的通信是否正常呢？我们就用远程工具来测试一下。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098792-180b3615-fdce-4841-8e94-60d0f7083497.jpeg)

主机与虚拟机通信正常。

这就是桥接模式的设置步骤，相信大家应该学会了如何去设置桥接模式了。桥接模式配置简单，但如果你的网络环境是 ip 资源很缺少或对 ip 管理比较严格的话，那桥接模式就不太适用了。如果真是这种情况的话，我们该如何解决呢？接下来，我们就来认识 vmware 的另一种网络模式：NAT 模式。

二、NAT（地址转换模式）

刚刚我们说到，如果你的网络 ip 资源紧缺，但是你又希望你的虚拟机能够联网，这时候 NAT 模式是最好的选择。NAT 模式借助虚拟 NAT 设备和虚拟 DHCP 服务器，使得虚拟机可以联网。其网络结构如下图所示：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098751-091c3b48-fec8-4bb2-916e-b16f0634b516.jpeg)

在 NAT 模式中，主机网卡直接与虚拟 NAT 设备相连，然后虚拟 NAT 设备与虚拟 DHCP 服务器一起连接在虚拟交换机 VMnet8 上，这样就实现了虚拟机联网。那么我们会觉得很奇怪，为什么需要虚拟网卡 VMware Network Adapter VMnet8 呢？原来我们的 VMware Network Adapter VMnet8 虚拟网卡主要是为了实现主机与虚拟机之间的通信。在之后的设置步骤中，我们可以加以验证。

首先，设置虚拟机中 NAT 模式的选项，打开 vmware，点击“编辑”下的“虚拟网络编辑器”，设置 NAT 参数及 DHCP 参数。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098787-fa5cc218-8d03-45e9-84c5-92a6a05f223c.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098802-f73a4965-095f-4c3f-a0e8-0cc6d27f0256.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098785-98bd34de-6d9f-4f23-ac2d-7f0c9ee82973.jpeg)

将虚拟机的网络连接模式修改成 NAT 模式，点击“编辑虚拟机设置”。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098821-a45bd243-4c21-4b71-87ef-dceacf19eda0.jpeg)

点击“网络适配器”，选择“NAT 模式”

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098805-7dd6b8a0-e5a9-449d-89a3-8b135222b028.jpeg)

然后开机启动系统，编辑网卡配置文件，命令为 vi /etc/sysconfig/network-scripts/ifcfg-eth0

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098807-65f0682b-015e-4df9-af4d-bca2b5a52025.jpeg)

具体配置如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098832-1ac99ee5-6257-4373-ae61-1131adf1afcf.jpeg)

编辑完成，保存退出，然后重启虚拟机网卡，动态获取 ip 地址，使用 ping 命令 ping 外网 ip，测试能否联网。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098824-6010fc01-167e-4313-993f-d35130c52095.jpeg)

之前，我们说过 VMware Network Adapter VMnet8 虚拟网卡的作用，那我们现在就来测试一下。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098821-4f3b67fb-a159-46c2-892b-1ddde535f5a6.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098854-9a7e1d84-723f-4816-abe8-f010ecdbd292.jpeg)

如此看来，虚拟机能联通外网，确实不是通过 VMware Network Adapter VMnet8 虚拟网卡，那么为什么要有这块虚拟网卡呢？

之前我们就说 VMware Network Adapter VMnet8 的作用是主机与虚拟机之间的通信，接下来，我们就用远程连接工具来测试一下。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098882-2551f3bc-4f38-4a75-affe-91fd509f3176.jpeg)

然后，将 VMware Network Adapter VMnet8 启用之后，发现远程工具可以连接上虚拟机了。

那么，这就是 NAT 模式，利用虚拟的 NAT 设备以及虚拟 DHCP 服务器来使虚拟机连接外网，而 VMware Network Adapter VMnet8 虚拟网卡是用来与虚拟机通信的。

三、Host-Only（仅主机模式）

Host-Only 模式其实就是 NAT 模式去除了虚拟 NAT 设备，然后使用 VMware Network Adapter VMnet1 虚拟网卡连接 VMnet1 虚拟交换机来与虚拟机通信的，Host-Only 模式将虚拟机与外网隔开，使得虚拟机成为一个独立的系统，只与主机相互通讯。其网络结构如下图所示：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098860-63c786a8-1067-4d03-b193-06c79c19ab2a.jpeg)

通过上图，我们可以发现，如果要使得虚拟机能联网，我们可以将主机网卡共享给 VMware Network Adapter VMnet1 网卡，从而达到虚拟机联网的目的。接下来，我们就来测试一下。

首先设置“虚拟网络编辑器”，可以设置 DHCP 的起始范围。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098846-49e70c14-f019-4672-8e48-8ea869759cc0.jpeg)

设置虚拟机为 Host-Only 模式。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098956-aa07d2b6-f9c9-400a-ade8-943ed0788396.jpeg)

开机启动系统，然后设置网卡文件。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098908-dde060c3-c79d-4618-8a4b-95638c9c8ee8.jpeg)

保存退出，然后重启网卡，利用远程工具测试能否与主机通信。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098909-7e89ba50-87c8-4771-af21-917bd3e42c02.jpeg)

主机与虚拟机之间可以通信，现在设置虚拟机联通外网。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098877-7b27d77f-e9ba-4af5-bcb2-4ae4032c7c59.jpeg)

我们可以看到上图有一个提示，强制将 VMware Network Adapter VMnet1 的 ip 设置成 192.168.137.1，那么接下来，我们就要将虚拟机的 DHCP 的子网和起始地址进行修改，点击“虚拟网络编辑器”

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098932-a9a17fc2-4c2f-4fdb-a1ea-4638acdfdc0a.jpeg)

重新配置网卡，将 VMware Network Adapter VMnet1 虚拟网卡作为虚拟机的路由。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098911-64b955dd-fd20-4732-a2c3-e3d3586cd77e.jpeg)

重启网卡，然后通过 远程工具测试能否联通外网以及与主机通信。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/udwtxy/1616124098891-ee532e62-4fbd-4470-afd4-d396521565d6.jpeg)

测试结果证明可以使得虚拟机连接外网。

以上就是关于 vmware 三种网络模式的工作原理及配置详解。
