---
title: Neutron 工作原理
---

# 概述

neutron 是 openstack 的一个重要模块，也是比较难以理解和 debug 的模块之一。

我这里安装如图安装了经典的三个节点的 Havana 的 Openstack

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dk3eew/1616123209479-a2bcb576-f320-42d7-a4df-0a7a050420f7.jpeg)

图 1

分三个网络：

- External Network/API Network，这个网络是连接外网的，无论是用户调用 Openstack 的 API，还是创建出来的虚拟机要访问外网，或者外网要 ssh 到虚拟机，都需要通过这个网络

- Data Network，数据网络，虚拟机之间的数据传输通过这个网络来进行，比如一个虚拟机要连接另一个虚拟机，虚拟机要连接虚拟的路由都是通过这个网络来进行

- Management Network，管理网络，Openstack 各个模块之间的交互，连接数据库，连接 Message Queue 都是通过这个网络来。

将这三个网络隔离，一方面是安全的原因，在虚拟机里面，无论采取什么手段，干扰的都紧紧是 Data Network，都不可能访问到我的数据库，一方面是流量分离，Management Network 的流量不是很大的，而且一般都会比较优雅的使用，而 Data network 和 External Network 就需要进行流量控制的策略。

我的这个网络结构有些奇怪，除了 Controller 节点是两张网卡之外，其他的都多了一张网卡连接到 external network，这个网卡是用来做 apt-get 的，因为 Compute Node 按说是没有网卡连接到外网的，为了 apt-get 添加了 eth0，Network Node 虽然有一个网卡 eth1 是连接外网的，然而在 neutron 配置好之前，这个网卡通常是没有 IP 的，为了 apt-get 也添加了 eth0，有人说可以通过添加 route 规则都通过 Controller 连接外网，但是对于初学的人，这个样比较容易操作。

neutron 是用来创建虚拟网络的，所谓虚拟网络，就是虚拟机启动的时候会有一个虚拟网卡，虚拟网卡会连接到虚拟的 switch 上，虚拟的 switch 连接到虚拟的 router 上，虚拟的 router 最终和物理网卡联通，从而虚拟网络和物理网络联通起来。

neutron 分成多个模块分布在三个节点上。

Controller 节点：

- neutron-server，用于接受 API 请求创建网络，子网，路由器等，然而创建的这些东西仅仅是一些数据结构在数据库里面

Network 节点：

- neutron-l3-agent，用于创建和管理虚拟路由器，当 neutron-server 将路由器的数据结构创建好，它是做具体的事情的，真正的调用命令行将虚拟路由器，路由表，namespace，iptables 规则全部创建好

- neutron-dhcp-agent，用于创建和管理虚拟 DHCP Server，每个虚拟网络都会有一个 DHCP Server，这个 DHCP Server 为这个虚拟网络里面的虚拟机提供 IP

- neutron-openvswith-plugin-agent，这个是用于创建虚拟的 L2 的 switch 的，在 Network 节点上，Router 和 DHCP Server 都会连接到二层的 switch 上

Compute 节点：

- neutron-openvswith-plugin-agent，这个是用于创建虚拟的 L2 的 switch 的，在 Compute 节点上，虚拟机的网卡也是连接到二层的 switch 上

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dk3eew/1616123209494-7100becf-9bd4-4260-9859-6d78ebc03def.jpeg)

图 2

当我们搭建好了 Openstack，然后创建好了 tenant 后，我们会为这个 tenant 创建一个网络。

\#!/bin/bash

TENANT_NAME="openstack"

TENANT_NETWORK_NAME="openstack-net"

TENANT_SUBNET_NAME="${TENANT_NETWORK_NAME}-subnet"

TENANT_ROUTER_NAME="openstack-router"

FIXED_RANGE="192.168.0.0/24"

NETWORK_GATEWAY="192.168.0.1"

PUBLIC_GATEWAY="172.24.1.1"

PUBLIC_RANGE="172.24.1.0/24"

PUBLIC_START="172.24.1.100"

PUBLIC_END="172.24.1.200"

TENANT_ID=$(keystone tenant-list | grep " $TENANT_NAME " | awk '{print $2}')

(1) TENANT_NET_ID=$(neutron net-create --tenant_id $TENANT_ID $TENANT_NETWORK_NAME --provider:network_type gre --provider:segmentation_id 1 | grep " id " | awk '{print $4}')

(2) TENANT_SUBNET_ID=$(neutron subnet-create --tenant_id $TENANT_ID --ip_version 4 --name $TENANT_SUBNET_NAME $TENANT_NET_ID $FIXED_RANGE --gateway $NETWORK_GATEWAY --dns_nameservers list=true 8.8.8.8 | grep " id " | awk '{print $4}')

(3) ROUTER_ID=$(neutron router-create --tenant_id $TENANT_ID $TENANT_ROUTER_NAME | grep " id " | awk '{print $4}')

(4) neutron router-interface-add $ROUTER_ID $TENANT_SUBNET_ID

(5) neutron net-create public --router:external=True

(6) neutron subnet-create --ip_version 4 --gateway $PUBLIC\_GATEWAY public $PUBLIC\_RANGE --allocation-pool start=$PUBLIC_START,end=$PUBLIC_END --disable-dhcp --name public-subnet

(7) neutron router-gateway-set ${TENANT_ROUTER_NAME} public

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dk3eew/1616123209460-e8eb430a-0276-4596-b53b-7cb8faeeac66.jpeg)

图 3

1. 为这个 Tenant 创建一个 private network，不同的 private network 是需要通过 VLAN tagging 进行隔离的，互相之间 broadcast 不能到达，这里我们用的是 GRE 模式，也需要一个类似 VLAN ID 的东西，称为 Segment ID

2. 创建一个 private network 的 subnet，subnet 才是真正配置 IP 网段的地方，对于私网，我们常常用 192.168.0.0/24 这个网段

3. 为这个 Tenant 创建一个 Router，才能够访问外网

4. 将 private network 连接到 Router 上

5. 创建一个 External Network

6. 创建一个 Exernal Network 的 Subnet，这个外网逻辑上代表了我们数据中心的物理网络，通过这个物理网络，我们可以访问外网。因而 PUBLIC_GATEWAY 应该设为数据中心里面的 Gateway， PUBLIC_RANGE 也应该和数据中心的物理网络的 CIDR 一致，否则连不通，而之所以设置 PUBLIC_START 和 PUBLIC_END，是因为在数据中心中，不可能所有的 IP 地址都给 Openstack 使用，另外可能搭建了 VMware Vcenter，可能有物理机器，仅仅分配一个区间给 Openstack 来用。

7. 将 Router 连接到 External Network

经过这个流程，从虚拟网络，到物理网络就逻辑上联通了。

创建完毕网络，如果不创建虚拟机，我们还是发现 neutron 的 agent 还是做了很多工作的，创建了很多的虚拟网卡和 switch

在 Compute 节点上：

root@ComputeNode:~# ip addr

1: eth0: mtu 1500 qdisc pfifo_fast state UP qlen 1000

link/ether 08:00:27:49:5c:41 brd ff:ff:ff:ff:ff:ff

inet 172.24.1.124/22 brd 172.24.1.255 scope global eth0

2: eth2: mtu 1500 qdisc pfifo_fast state UP qlen 1000

link/ether 08:00:27:8e:42:2c brd ff:ff:ff:ff:ff:ff

inet 192.168.56.124/24 brd 192.168.56.255 scope global eth2

3: eth3: mtu 1500 qdisc pfifo_fast state UP qlen 1000

link/ether 08:00:27:68:92:ce brd ff:ff:ff:ff:ff:ff

inet 10.10.10.124/24 brd 10.10.10.255 scope global eth3

4: br-int: mtu 1500 qdisc noqueue state UNKNOWN

link/ether d6:2a:96:12:4a:49 brd ff:ff:ff:ff:ff:ff

5: br-tun: mtu 1500 qdisc noqueue state UNKNOWN

link/ether a2:ee:75:bd:af:4a brd ff:ff:ff:ff:ff:ff

6: qvof5da998c-82: mtu 1500 qdisc pfifo_fast state UP qlen 1000

link/ether c2:7e:50:de:8c:c5 brd ff:ff:ff:ff:ff:ff

7: qvbf5da998c-82: mtu 1500 qdisc pfifo_fast state UP qlen 1000

link/ether c2:33:73:40:8f:e0 brd ff:ff:ff:ff:ff:ff

root@ComputeNode:~# ovs-vsctl show

39f69272-17d4-42bf-9020-eecc9fe8cde6

Bridge br-int

Port patch-tun

Interface patch-tun

type: patch

options: {peer=patch-int}

Port br-int

Interface br-int

type: internal

Bridge br-tun

Port patch-int

Interface patch-int

type: patch

options: {peer=patch-tun}

Port "gre-1"

Interface "gre-1"

type: gre

options: {in_key=flow, local_ip="10.10.10.124", out_key=flow, remote_ip="10.10.10.121"}

Port br-tun

Interface br-tun

type: internal

ovs_version: "1.10.2"

在 Network Node 上：

root@NetworkNode:~# ip addr

1: eth0: mtu 1500 qdisc pfifo_fast state UP qlen 1000

link/ether 08:00:27:22:8a:7a brd ff:ff:ff:ff:ff:ff

inet 172.24.1.121/22 brd 172.24.1.255 scope global eth0

2: eth1: mtu 1500 qdisc pfifo_fast state UP qlen 1000

link/ether 08:00:27:f1:31:81 brd ff:ff:ff:ff:ff:ff

3: eth2: mtu 1500 qdisc pfifo_fast state UP qlen 1000

link/ether 08:00:27:56:7b:8a brd ff:ff:ff:ff:ff:ff

inet 192.168.56.121/24 brd 192.168.56.255 scope global eth2

4: eth3: mtu 1500 qdisc pfifo_fast state UP qlen 1000

link/ether 08:00:27:26:bc:84 brd ff:ff:ff:ff:ff:ff

inet 10.10.10.121/24 brd 10.10.10.255 scope global eth3

5: br-ex: mtu 1500 qdisc noqueue state UNKNOWN

link/ether 08:00:27:f1:31:81 brd ff:ff:ff:ff:ff:ff

inet 172.24.1.8/24 brd 172.24.1.255 scope global br-ex

6: br-int: mtu 1500 qdisc noqueue state UNKNOWN

link/ether 22:fe:f1:9b:29:4b brd ff:ff:ff:ff:ff:ff

7: br-tun: mtu 1500 qdisc noqueue state UNKNOWN

link/ether c6:ea:94:ff:23:41 brd ff:ff:ff:ff:ff:ff

root@NetworkNode:~# ip netns

qrouter-b2510953-1ae4-4296-a628-1680735545ac

qdhcp-96abd26b-0a2f-448b-b92c-4c98b8df120b

root@NetworkNode:~# ip netns exec qrouter-b2510953-1ae4-4296-a628-1680735545ac ip addr

8: qg-97040ca3-2c: mtu 1500 qdisc noqueue state UNKNOWN

link/ether fa:16:3e:26:57:e3 brd ff:ff:ff:ff:ff:ff

inet 172.24.1.100/24 brd 172.24.1.255 scope global qg-97040ca3-2c

11: qr-e8b97930-ac: mtu 1500 qdisc noqueue state UNKNOWN

link/ether fa:16:3e:43:ef:16 brd ff:ff:ff:ff:ff:ff

inet 192.168.0.1/24 brd 192.168.0.255 scope global qr-e8b97930-ac

root@NetworkNode:~# ip netns exec qdhcp-96abd26b-0a2f-448b-b92c-4c98b8df120b ip addr

9: tapde5739e1-95: mtu 1500 qdisc noqueue state UNKNOWN

link/ether fa:16:3e:19:8c:67 brd ff:ff:ff:ff:ff:ff

inet 192.168.0.2/24 brd 192.168.0.255 scope global tapde5739e1-95

inet 169.254.169.254/16 brd 169.254.255.255 scope global tapde5739e1-95

root@NetworkNode:~# ovs-vsctl show

d5d5847e-1c9e-4770-a68c-7a695b7b95cd

Bridge br-ex

Port "qg-97040ca3-2c"

Interface "qg-97040ca3-2c"

type: internal

Port "eth1"

Interface "eth1"

Port br-ex

Interface br-ex

type: internal

Bridge br-int

Port patch-tun

Interface patch-tun

type: patch

options: {peer=patch-int}

Port "tapde5739e1-95"

tag: 1

Interface "tapde5739e1-95"

type: internal

Port br-int

Interface br-int

type: internal

Port "qr-e8b97930-ac"

tag: 1

Interface "qr-e8b97930-ac"

type: internal

Bridge br-tun

Port patch-int

Interface patch-int

type: patch

options: {peer=patch-tun}

Port "gre-2"

Interface "gre-2"

type: gre

options: {in_key=flow, local_ip="10.10.10.121", out_key=flow, remote_ip="10.10.10.124"}

Port br-tun

Interface br-tun

type: internal

ovs_version: "1.10.2"

这时候如果我们创建一个虚拟机在这个网络里面，在 Compute Node 多了下面的网卡：

13: qvof5da998c-82: mtu 1500 qdisc pfifo_fast state UP qlen 1000

link/ether c2:7e:50:de:8c:c5 brd ff:ff:ff:ff:ff:ff

14: qvbf5da998c-82: mtu 1500 qdisc pfifo_fast state UP qlen 1000

link/ether c2:33:73:40:8f:e0 brd ff:ff:ff:ff:ff:ff

15: qbr591d8cc4-df: mtu 1500 qdisc noqueue state UP

link/ether f2:d9:f0:d5:48:c8 brd ff:ff:ff:ff:ff:ff

16: qvo591d8cc4-df: mtu 1500 qdisc pfifo_fast state UP qlen 1000

link/ether e2:58:d4:dc:b5:16 brd ff:ff:ff:ff:ff:ff

17: qvb591d8cc4-df: mtu 1500 qdisc pfifo_fast master qbr591d8cc4-df state UP qlen 1000

link/ether f2:d9:f0:d5:48:c8 brd ff:ff:ff:ff:ff:ff

18: tap591d8cc4-df: mtu 1500 qdisc pfifo_fast master qbr591d8cc4-df state UNKNOWN qlen 500

link/ether fe:16:3e:6e:ba:d0 brd ff:ff:ff:ff:ff:ff

如果我们按照 ovs-vsctl show 的网卡桥接关系，变可以画出下面的图

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dk3eew/1616123209536-dacb32ea-a9c6-4f21-bcfc-bfb8bfe2f36e.jpeg)

图 4

当然如果你配的不是 GRE 而是 VLAN 的话，便有下面这个著名的复杂的图。

在 network Node 上：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dk3eew/1616123209516-12e64548-54a7-479b-9098-39a2abdccf09.jpeg)

图 5

在 Compute Node 上：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dk3eew/1616123209481-28a3ced3-6a9f-47a4-bb44-24d4d1ceeb66.jpeg)

图 6

当看到这里，很多人脑袋就大了，openstack 为什么要创建这么多的虚拟网卡，他们之间什么关系，这些 dl_vlan, mod_vlan_vid 都是什么东东啊？

这些网卡看起来复杂，却是各有用处，这种虚拟网络拓扑，正是我们经常使用的物理网络的拓扑结构。

我们先来回到一个最最熟悉不过的场景，我们的大学寝室，当时我们还买不起路由器，所以一般采取的方法如下图所示：

寝室长的机器上弄两张网卡，寝室买一个 HUB，其他人的电脑都接到 HUB 上，寝室长的电脑的两张网卡一张对外连接网络，一张对内连接 HUB。寝室长的电脑其实充当的便是路由器的作用。

后来条件好了，路由器也便宜了，所以很多家庭也是类似的拓扑结构，只不过将 Computer1 和 switch 合起来，变成了一个路由器，路由器也是有多个口一个连接 WLAN，一个连接 LAN。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dk3eew/1616123209487-4740b1f4-55ca-4fef-893e-05869d01e536.jpeg)

图 7

现在我们想象一个寝室变成了一台物理机 Hypervisor，所有的电脑都变成了虚拟机，就成了下面的样子。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dk3eew/1616123209460-3d9b83f0-8ddc-4964-ab00-363e28034e44.jpeg)

图 8

我们先忽略 qbr 和 DHCP Server，以及 namespace。

br-int 就是寝室里面的 HUB，所有虚拟机都会连接到这个 switch 上，虚拟机之间的相互通信就是通过 br-int 来的。

Router 就是你们寝室长的电脑，一边接在 br-int 上，一边接在对外的网口上，br-ex/eth0 外面就是我们的物理网络。

为什么会有个 DHCP Server 呢，是同一个 private network 里的虚拟机得到 IP 都是通过这个 DHCP Server 来的，这个 DHCP Server 也是连接到 br-int 上和虚拟机进行通信。

为什么会有 qbr 呢，这是和 security group 的概念有关，openstack 中的 security group 开通哪些端口，屏蔽哪些端口是用 iptables 来实现的，然而 br-int 这些虚拟 bridge 都是 openvswitch 创建的，openvswitch 的 kernel mode 和 netfilter 的 kernel mode 不兼容，一个 IP 包进来要么走 iptables 的规则进行处理，要么走 openvswitch 的规则进行处理，通过上面的复杂的图我们可以看到，br-int 上面有很多 openvswitch 的规则，比如 vlan tag 等，所以 iptables 必须要另外建立一个 linux bridge 来做，因而有了 qbr，在了解拓扑结构的时候，可以将 qbr 忽略，看成 VM 直接连接到 br-int 上就可以了。

为什么会有 namespace 呢，java 的 namespace 是为了相同的类名，不同的 namespace，显然是不同的类。openstack 也想做到这一点，不同的 tenant 都想创建自己的 router 和 private network，彼此不知道别人指定了哪些网段，很有可能两个 tenant 都指定了 192.168.0.0/24，这样不同的 private network 的路由表，DHCP Server 就需要隔离，不然就乱了，因而有了 namespace。

上面的图其实就是单节点的 openstack 的网络结构，虽然复杂，但是就是把我们家里的，或者寝室里面的物理机搬到一个 Hypervisor 上了，其结构就不难理解了。

当然单节点的 openstack 不过是个测试环境，compute 节点和 network 节点也是要分开的，如图 4，每个机器上都有了自己的 br-int。

但是对于虚拟机和虚拟 Router 来讲，他们仍然觉得他们是连接到了一个大的 L2 的 br-int 上，通过这个 br-int 相互通信的，他们感受不到 br-int 下面的虚拟网卡 br-tun。所以对于多节点结构，我们可以想象 br-int 是一个大的，横跨所有的 compute 和 network 节点的二层 switch，虚拟机之间的通信以及虚拟机和 Router 的通信，就像在一个寝室一样的。

然而 br-int 毕竟被物理的割开了，需要有一种方式将他们串联起来，openstack 提供了多种方式，图 4 中是用 GRE tunnel 将不同机器的 br-int 连接起来，图 5 图 6 是通过 VLAN 将 br-int 连接起来，当然还可以使用 vxlan。

这就是为什么 openstack 有了 br-int 这个 bridge，但是不把所有的 openvswitch 的规则都在它上面实现。就是为了提供这种灵活性，对于虚拟机来讲，看到的是一大整个 br-int，不同机器的 br-int 可以有多种方式连接，这在 br-int 下面的网卡上面实现。

如果有不同的 Tenant，创建了不同的 private network，为了在 data network 上对包进行隔离，创建 private network 的时候，需要指定 vlanid 或者 segmentid。

从 ovs-vsctl show 我们可以看到，不同的 tenant 的 private network 上创建的虚拟机，连接到 br-int 上的时候是带 tag 的，所以不同 tenant 的虚拟机，即便连接到同一个 br-int 上，因为 tag 不同，也是不能相互通信的，然而同一个机器上的 tag 的计数是仅在本机有效的，并不使用我们创建 private network 的时候指定的全局唯一的 vlanid 或者 segmentid，一个 compute 节点上的 br-int 上的 tag 1 和另一台 compute 节点上的 br-int 的 tag1 很可能是两码事。全局的 vlanid 和 segmentid 仅仅在 br-int 以下的虚拟网卡和物理网络中使用，虚拟机所有能看到的东西，到 br-int 为止，看不到打通 br-int 所要使用的 vlanid 和 segmentid。

从局部有效的 taging 到全局有效的 vlanid 或者 segmentid 的转换，都是通过 openvswitch 的规则，在 br-tun 或者 br-eth1 上实现。

我们可以用下面的命令看一下这个规则:

在 Compute 节点上：

private network “openstack-net”的 tag 在这台机器上是 2，而我们创建的时候的 segmentid 设定的是 1

Bridge br-int

Port patch-tun

Interface patch-tun

type: patch

options: {peer=patch-int}

Port br-int

Interface br-int

type: internal

Port "qvo591d8cc4-df"

tag: 2

Interface "qvo591d8cc4-df"

root@ComputeNodeCliu8:~# ovs-ofctl dump-flows br-tun

NXST_FLOW reply (xid=0x4):

//in_port=1 是指包是从 patch-int，也即是从虚拟机来的，所以是发送规则，跳转到 table1

cookie=0x0, duration=77419.191s, table=0, n_packets=22, n_bytes=2136, idle_age=6862, hard_age=65534, priority=1,in_port=1 actions=resubmit(,1)

//in_port=2 是指包是从 GRE 来的，也即是从物理网络来的，所以是接收规则，跳转到 table2

cookie=0x0, duration=77402.19s, table=0, n_packets=3, n_bytes=778, idle_age=6867, hard_age=65534, priority=1,in_port=2 actions=resubmit(,2)

cookie=0x0, duration=77418.403s, table=0, n_packets=0, n_bytes=0, idle_age=65534, hard_age=65534, priority=0 actions=drop

//multicast，跳转到 table21

cookie=0x0, duration=77416.63s, table=1, n_packets=21, n_bytes=2094, idle_age=6862, hard_age=65534, priority=0,dl_dst=01:00:00:00:00:00/01:00:00:00:00:00 actions=resubmit(,21)

//broadcast，跳转到 table 20

cookie=0x0, duration=77417.389s, table=1, n_packets=1, n_bytes=42, idle_age=6867, hard_age=65534, priority=0,dl_dst=00:00:00:00:00:00/01:00:00:00:00:00 actions=resubmit(,20)

//这是接收规则的延续，如果接收的 tun_id=0x1 则转换为本地的 tag，mod_vlan_vid:2，跳转到 table 10

cookie=0x0, duration=6882.254s, table=2, n_packets=3, n_bytes=778, idle_age=6867, priority=1,tun_id=0x1 actions=mod_vlan_vid:2,resubmit(,10)

cookie=0x0, duration=77415.638s, table=2, n_packets=0, n_bytes=0, idle_age=65534, hard_age=65534, priority=0 actions=drop

cookie=0x0, duration=77414.432s, table=3, n_packets=0, n_bytes=0, idle_age=65534, hard_age=65534, priority=0 actions=drop

cookie=0x0, duration=77412.825s, table=10, n_packets=3, n_bytes=778, idle_age=6867, hard_age=65534, priority=1 actions=learn(table=20,hard_timeout=300,priority=1,NXM_OF_VLAN_TCI\[0..11],NXM_OF_ETH_DST\[]=NXM_OF_ETH_SRC\[],load:0->NXM_OF_VLAN_TCI\[],load:NXM_NX_TUN_ID\[]->NXM_NX_TUN_ID\[],output:NXM_OF_IN_PORT\[]),output:1

cookie=0x0, duration=77411.549s, table=20, n_packets=0, n_bytes=0, idle_age=65534, hard_age=65534, priority=0 actions=resubmit(,21)

//这是发送规则的延续，如果接收到的 dl_vlan=2，则转换为物理网络的 segmentid=1，set_tunnel:0x1

cookie=0x0, duration=6883.119s, table=21, n_packets=10, n_bytes=1264, idle_age=6862, priority=1,dl_vlan=2 actions=strip_vlan,set_tunnel:0x1,output:2

cookie=0x0, duration=77410.56s, table=21, n_packets=11, n_bytes=830, idle_age=6885, hard_age=65534, priority=0 actions=drop

在 Network 节点上：

Bridge br-int

Port patch-tun

Interface patch-tun

type: patch

options: {peer=patch-int}

Port "tapde5739e1-95"

tag: 1

Interface "tapde5739e1-95"

type: internal

Port br-int

Interface br-int

type: internal

Port "qr-e8b97930-ac"

tag: 1

Interface "qr-e8b97930-ac"

type: internal

非常相似的规则。

root@NetworkNodeCliu8:~# ovs-ofctl dump-flows br-tun

NXST_FLOW reply (xid=0x4):

//in_port=1 是指包是从 patch-int，也即是从虚拟机来的，所以是发送规则，跳转到 table1

cookie=0x0, duration=73932.142s, table=0, n_packets=12, n_bytes=1476, idle_age=3380, hard_age=65534, priority=1,in_port=1 actions=resubmit(,1)

//in_port=2 是指包是从 GRE 来的，也即是从物理网络来的，所以是接收规则，跳转到 table2

cookie=0x0, duration=73914.323s, table=0, n_packets=9, n_bytes=1166, idle_age=3376, hard_age=65534, priority=1,in_port=2 actions=resubmit(,2)

cookie=0x0, duration=73930.934s, table=0, n_packets=0, n_bytes=0, idle_age=65534, hard_age=65534, priority=0 actions=drop

//multicast，跳转到 table21

cookie=0x0, duration=73928.59s, table=1, n_packets=6, n_bytes=468, idle_age=65534, hard_age=65534, priority=0,dl_dst=01:00:00:00:00:00/01:00:00:00:00:00 actions=resubmit(,21)

//broadcast，跳转到 table20

cookie=0x0, duration=73929.695s, table=1, n_packets=3, n_bytes=778, idle_age=3380, hard_age=65534, priority=0,dl_dst=00:00:00:00:00:00/01:00:00:00:00:00 actions=resubmit(,20)

//这是接收规则的延续，如果接收的 tun_id=0x1 则转换为本地的 tag，mod_vlan_vid:1，跳转到 table 10

cookie=0x0, duration=73906.864s, table=2, n_packets=9, n_bytes=1166, idle_age=3376, hard_age=65534, priority=1,tun_id=0x1 actions=mod_vlan_vid:1,resubmit(,10)

cookie=0x0, duration=73927.542s, table=2, n_packets=0, n_bytes=0, idle_age=65534, hard_age=65534, priority=0 actions=drop

cookie=0x0, duration=73926.403s, table=3, n_packets=0, n_bytes=0, idle_age=65534, hard_age=65534, priority=0 actions=drop

cookie=0x0, duration=73925.611s, table=10, n_packets=9, n_bytes=1166, idle_age=3376, hard_age=65534, priority=1 actions=learn(table=20,hard_timeout=300,priority=1,NXM_OF_VLAN_TCI\[0..11],NXM_OF_ETH_DST\[]=NXM_OF_ETH_SRC\[],load:0->NXM_OF_VLAN_TCI\[],load:NXM_NX_TUN_ID\[]->NXM_NX_TUN_ID\[],output:NXM_OF_IN_PORT\[]),output:1

cookie=0x0, duration=73924.858s, table=20, n_packets=0, n_bytes=0, idle_age=65534, hard_age=65534, priority=0 actions=resubmit(,21)

//这是发送规则的延续，如果接收到的 dl_vlan=1，则转换为物理网络的 segmentid=1，set_tunnel:0x1

cookie=0x0, duration=73907.657s, table=21, n_packets=0, n_bytes=0, idle_age=65534, hard_age=65534, priority=1,dl_vlan=1 actions=strip_vlan,set_tunnel:0x1,output:2

cookie=0x0, duration=73924.117s, table=21, n_packets=6, n_bytes=468, idle_age=65534, hard_age=65534, priority=0 actions=drop
