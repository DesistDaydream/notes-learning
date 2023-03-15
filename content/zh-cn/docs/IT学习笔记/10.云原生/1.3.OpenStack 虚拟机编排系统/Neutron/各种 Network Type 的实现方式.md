---
title: 各种 Network Type 的实现方式
---

# Linux Bridge 可用的网络类型详解

Linux bridge 技术非常成熟，而且高效，所以业界很多 OpenStack 方案选择 linux bridge，比如 Rackspace 的 private cloud。

open vswitch 实现的 Neutron 虚拟网络较为复杂，不易理解；而 linux bridge 方案更直观。先理解 linux bridge 方案后再学习 open vswitch 方案会更容易。并且可以通过两种方案的对比更加深入地理解 Neutron 网络。

所谓的 Linux Bridge Provider 就是 Neutron 直接利用 Linux 中的 Bridge 来实现自身的网络功能，而没有任何额外的功能。就是类似直接在系统中使用 ip bridge、brctl 之类的命令。

Linux Bridge 环境中，一个数据包从 instance 发送到物理网卡会经过下面几个网络设备：

1. tap interface 命名为 tapN (N 为 0, 1, 2, 3......)
2. linux bridge 命名为 brqXXXX。
3. vlan interface 命名为 ethX.Y（X 为 interface 的序号，Y 为 vlan id）
4. vxlan interface 命名为 vxlan-Z（z 是 VNI）
5. 物理 interface 命名为 ethX（X 为 interface 的序号）

## Local Network 本地网络，无法与宿主机之外通信

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oo0pgf/1616123235077-f61119cd-ce2f-4f23-8343-9c6ba9696e26.jpeg)

## Flat Network 平面网络，通过桥接与宿主机之外通信

VM——Bridge——网卡(纯二层实现，适合私有云)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oo0pgf/1616123235091-3b8a5417-47cd-437a-82f7-3a59c44dd7a9.jpeg)

该方法最为简单直接有效，通过纯二层来实现 VM 与外部设备的互通，通过 Linux Bridge，并手动管理一些虚拟网络设备即可

环境描述，host1，host2，...，hostN 等设备有两个网卡，eth0 为管理 192.168.0.0/24，eth1 不配 IP，但是 eth1 连接的物理网络网段为 10.0.0.0/24。这时候如果想让 host 上所创建的虚拟机的 IP 在 10.0.0.0/24 网段中，且可以直接与该网段通信，那么就可以下面描述的方法。这种方法，不用 nat，直接通过桥接方式让虚拟机与外部网络 IP 保持一致，可以减少网络开销。

1. 在 host 上创建一个 bridge 设备 br1。
   1. brctl addbr br1
2. 在 network 中创建 (与所需连通的外部网络 IP 相同的)subnet 以及 port。
3. 将创建完的 port 关联到虚拟机中，这时候会在后台生成一个对应的 tap 设备
4. 通过后台，把 eth1 与 tap 设备关联到 br1 上。
   1. brctl addif tap0 br1 && brctl addif tap1 br1 && brctl addif tap2 br1 && brctl addif eth1 br1
5. 这时候的情况就像上图一样，eth1 和 tap 都是一个虚拟交换机上的端口，通过二层进入外部物理网络中，只要虚拟机上配的 10.0.0.0/24 网段的 IP 网关在外部网络中可以找到，那么这时候虚拟机就可以直接以 10.0.0.0/24 网段的 IP 连接到外部网络中了。

## VLAN Network 使用 VLAN 隔离 VM，并通过路由连通

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oo0pgf/1616123235065-7ade4af6-857d-40d6-958d-debb0b4fc779.jpeg)

## 租户网络与外部网络互相访问

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oo0pgf/1616123235150-1331b664-155f-4f8b-b75e-974da25b7bf3.jpeg)

## VxLAN Network

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oo0pgf/1616123235118-913da696-7a1c-4bb5-9d82-b5783b0e86d7.jpeg)

# Open vSwitch 可用的网络类型详解

Neutron 会在 OVS 模式下会自动创建三种虚拟网络设备

1. br-ex # 连接外部（external）网络的网桥。(计算节点无该网络设备)
   1. br-ex 一般用网络设备名命名，比如写成 br-eth0、br-bond0 等等
2. br-int # 聚合 (integration) 多个实例网的网桥。所有 instance 在创建后都会自动创建一个 LinuxBridge，这个 LinuxBridge 通过 veth 对 连接到 br-int 上。
   1. 也可以称为 连接内部（internal）网络的网桥。
3. br-tun # 隧道（tunnel）网桥，基于隧道技术的 VxLAN 和 GRE 网络将使用该网桥进行通信。

注意：上述设备仅仅是自动创建出来，在实际情况中，并不一定会真正使用到

Open vSwitch 环境中，数据包从 instance 到物理网卡的路径

1. tap interface # 命名为 tapXXXX
2. linux bridge # 命名为 qbrXXXX。
3. veth pair # 命名为 qvbXXXX，qvoXXXX
4. OVS integration bridge # 命名为 br-int
5. OVS patch ports # 命名为 int-br-exX 和 phy-br-ethX(X 为 interface 序号)
6. OVS provider bridge # 命名为 br-exX(X 为 interface 序号)
7. vlan interface # 命名为 ethX.Y（X 为 interface 的序号，Y 为 vlan id）
8. vxlan interface # 命名为 vxlan-Z（z 是 VNI）
9. 物理 interface # 命名为 ethX（X 为 interface 的序号）
10. OVS tunnel bridge # 命名为 br-tun

OVS provider bridge 会在 flat 和 vlan 网络中使用；OVS tunnel bridge 则会在 vxlan 和 gre 网络中使用

instance—TAP—Linux Bridge—Veth Pair—OVS Bridge—其它

这是是 OVS 在 openstack 的最基本数据流走向

## Local Network 本地网络，无法与宿主机之外通信

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oo0pgf/1616123235087-35c044d8-beda-40c9-b946-6a4d04d81968.jpeg)

## Flat Network 平面网络，通过桥接与宿主机之外通信

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oo0pgf/1616123235109-8220d6a7-1e96-4a48-9612-7049716ce29b.jpeg)

如图所示，在 OVS 模式下，新的虚拟机实例创建出来后，需要先连接到 LInux Bridge 上，再连接到 OVS Bridge 上。Linux Bridge 与 OVS Bridge 上绑定了一对 veth 设备作为二者的连接。其中一共创建了 4 个虚拟网络设备。

1. TAP # VM 启动后生成的网络设备，用于与虚拟机内部网卡连接
2. Linux Bridge # 作为 VM 与 br-int 的桥梁中转。每为 VM 添加一个网卡，就会生成一个 Linux Bridge。为什么要多做这么一层中转而不让 VM 直接连到 br-int 的理由详见下文。
3. Veth Pair # 用于连接 Linux Bridge 与 OVS Bridge。
4. OVS Bridge # 所有 VM 通过其自身的 Linux Bridge 连接到 OVS Bridge

可以发现，对于同一个 VM 所使用的相关虚拟网络设备的名称，有一部分是相同的。qvofc1c6ebb-71、qvbfc1c6ebb-71、tapfc1c6ebb-71。其中 fc1c6ebb-71 这部分是相同的，不同的的是 qvo、qvb、tap。其中 qvo 与 qvb 是 一对 veth 设备，用来连接两个桥设备，tap 则用来连接 VM 与 桥设备

注意：

1. 上述 4 个虚拟网络设备与 dhcp 无关，不要混淆。这其中 DHCP 则是在创建 network 中的 subnet 的时候，决定是否创建的。如果不创建 dhcp，则上图中的 dhcp 以及其对应的 tap 设备则都不存在

通过几个命令可以看到在实际环境中的虚拟网络设备之间的关系

1. ovs-vsctl show # 在 OVS 中的网桥 br-int 上有一个 interface，名为 qvofc1c6ebb-71

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oo0pgf/1616123235082-6c539d65-6dc4-49d1-8a79-168795c33eb7.jpeg)

1. brctl show # 在 LinuxBridge 这个网桥 qbrfc1c6ebb-71 上，有一个 interface，名为 qvbfc1c6ebb-71
   1. 在这里还能看到实例网卡所关联的名为 tapfc1c6ebb-71 这个 TAP 设备也在 Linux Bridge 上

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oo0pgf/1616123235113-3f709aac-0438-47ca-98f8-bd82eb31b287.jpeg)

1. ethtools -S qvofc1c6ebb-71 && ethtools qvbfc1c6ebb-71 # 这是 一对 veth 设备 ，用来连接 Linux Bridge 与 OVS Bridge

                root@devstack-controller:~# ethtool -S qvbfc1c6ebb-71NIC statistics:    peer_ifindex: 12root@devstack-controller:~# ethtool -S qvofc1c6ebb-71NIC statistics:    peer_ifindex: 13

那问题来了，为什么 tapfc1c6ebb-71 不能像左边的 DHCP 设备 tap7970bdcd-f2 那样直接连接到 br-int 呢？

1. 其原因是： Open vSwitch 目前还不支持将 iptables 规则放在与它直接相连的 tap 设备上。如果做不到这一点，就无法实现 Security Group 功能。为了支持 Security Group，不得不多引入一个 Linux Bridge 支持 iptables。这样的后果就是网络结构更复杂了，路径上多了一个 linux bridge 和 一对 veth pair 设备。

而且，就算在同一个 network 中的 instance，每一个 instance 都会有一个对应的 LinuxBridge，这些 LinuxBridge 再与 OVS Brigde 相连，甚至一个 instance 上如果关联了多个 port，那么这个 instance 上则会连接多个 LinuxBridge(一个 port 对应一个 LinuxBridge)。那么所有 instance 都会间接连接到 OVS Bridge 这同一个 bridge 设备上，如何进行二层隔离呢？这个问题可以在下一段《高级应用，创建多个实例，且不同 VLAN 通过路由连通》文章中得到答案，提前说一声，是给每个 port 打上 tag，通过 tag 来进行二层隔离，相同 tag 的 port 在同一个 network 中

## VLAN Network 使用 VLAN 隔离 VM，并通过路由连通

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oo0pgf/1616123235119-d55a337b-ffbb-43e9-b9b4-f8e5b4301428.jpeg)

OVS 中的 VLAN tag

VM1，VM2，VM3 都在同一个网桥(br-int)上，为什么 VM1、VM2 与 VM3 无法互通呢，因为在创建 network 的时候，对他们进行了隔离，VM1 和 VM2 属于同一个 network。隔离的方式则是通过图中的 tag 标签，在使用 ovs-vsctl show 命令查看网络设备的时候，会看到 Port 下面有个一 tag 字段，这个 tag 字段表示了这个接口属于哪个 network，不同 tag 之间的接口上关联的 VM 是无法进行二层通信的，需要借助路由才可以。

br-int 上多了两个 port，如图所示：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oo0pgf/1616123235125-69e65c47-190b-4469-9b9d-99fb8dededcf.jpeg)

1. qr-d295b258-45，从命名上可以推断该 interface 对应 router_100_101 的 interface (d295b258-4586)，是 subnet_172_16_100_0 的网关。
2. qr-2ffdb861-73，从命名上可以推断该 interface 对应 router_100_101 的 interface (2ffdb861-731c)，是 subnet_172_16_101_0 的网关。

同时 route_100_101 运行在自己的 namespace 中，并将 br-int 上的两个 port 也添加进该 namespace 中。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oo0pgf/1616123235112-f13b4a1d-dbf7-4813-82e9-503af0a75729.jpeg)

如上所示，qrouter-a81cc110-16f4-4d6c-89d2-8af91cec9714 是 router 的 namespace，两个 Gateway IP 分别配置在 qr-2ffdb861-73 和 qr-d295b258-45 上。

## 租户网络与外部网络互相访问

VM——br-int——router——br-ex——网卡

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oo0pgf/1616123235113-884e5f64-9739-42ab-aa33-e017f6c8ed40.jpeg)

如图所示，如果 VM 想与宿主机相通，那么数据流应该如此走：VM——br-int——router——br-ex——eth2——外部物理网络。具体描述如下：

1. VM 到 br-int 比较简单，正常的二层
2. 当 VM 想要访问非本网段的网络时，br-int 无法处理三层网络请求，需要去网关找对方网段的路由。即 br-int——router 的过程。
3. 路由器经过路由，让数据包跳转到自己的另一个端口 92.168.0.123 上。即 router——br-ex 的过程
4. eth2 网卡作为从设备，被附加到 br-ex 网桥上(eth2 上的 IP 也无作用)，仅作为 br-ex 上的一个端口使用(原理详见 Network Virtual)。所以从路由器过来的数据流经过 eth1 直接流向外部物理网络，经过物理网络的交换机和路由器，到达 eth0 上的 IP。即 br-ex——eth1——外部物理网络 的过程

### Floating IP

如果外部网络想要访问内部 VM，那么则需要给 VM 绑定一个 floatIP，这个 floatIP 会出现在 router 上的网卡中，在 router 这个 namespace 中，会有 iptables 规则对 10.0.0.10 与 92.168.0.10 进行 nat 转换，这样在访问 92.168.0.10 的时候，数据包进入 eth2 后，还昂目的地址是 92.168.0.10，则执行 dnat 操作，将目的地址修改为 10.0.0.10 后再将数据包交给 br-int，这样从外面访问 92.168.0.10 就相当于是访问 10.0.0.10 了

## VxLAN Network

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oo0pgf/1616123235138-4cc1c71d-b4c9-48b5-83c2-4e2c81cf0048.jpeg)

# Namespace 的作用,以及浮动(弹性)IP(Floating IP)

namespace 功能详见 Namespace && CGroups 中关于网络 namespace 介绍。

[那么，为什么要在 openstack 中使用 namesapce 呢？](https://mp.weixin.qq.com/s?__biz=MzIwMTM5MjUwMg==&mid=2653587558&idx=1&sn=e4381acfef2030b74f870e4f7c3547c6&chksm=8d30807fba4709693926bc849dd4736b47e92e23c3b9d1af215466f524d21c7ec8865a859999&scene=21#wechat_redirect)

为了要实现多租户(tenant)，同样也是为了实现叠加网络(overlay network)。

有这么一种场景，如下图

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oo0pgf/1616123235138-6da55154-0ec2-4a5d-b64f-87af8be17357.jpeg)

左侧为租户 1，右侧为租户 2，都在同一 host 上开了虚拟机。其特征是网关 IP 配置在 TAP interface 上。因为没有 namespace 隔离，router_100_101 和 router_102_103 的路由条目都只能记录到控制节点操作系统（root namespace）的路由表中，内容如下：

    Destination Gateway Genmask Flags Metric Use Iface
      10.10.1.0  * 255.255.255.0   U     0      0      tap1
      10.10.2.0  * 255.255.255.0   U     0      0      tap2
      10.10.1.0  * 255.255.255.0   U     0      0      tap3
      10.10.2.0  * 255.255.255.0   U     0      0      tap4

这样的路由表是无法工作的。

按照路由表优先匹配原则，Tenant B 的数据包总是错误地被 Tenant A 的 router 路由。例如 vlan102 上有数据包要发到 vlan103。选择路由时，会匹配路由表的第二个条目，结果数据被错误地发到了 vlan101。

这时候，如果不使用 namespace 就会出现一个问题，两边的路由器中的路由表，是 host 上的路由表，出现相同的内容，这时候，租户 1 访问的网段有可能就会到租户 2 去，为了解决这个问题，则需要使用 namespace，让租户 1 与租户 2 拥有各自的网络栈，而不用共享 host 上的网络栈，这样，两边的路由器都维护自己的路由表，实现了租户隔离。

如果使用 namespace，网络结构如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oo0pgf/1616123235113-28148189-91f4-4ffa-9a8b-4a65cc919591.png)

其特征是网关 IP 配置在 namespace 中的 veth interface 上。每个 namespace 拥有自己的路由表。

router_100_101 的路由表内容如下：

    Destination Gateway Genmask Flags Metric Use Iface
    10.10.1.0 * 255.255.255.0   U     0      0    qr-1
    10.10.2.0 * 255.255.255.0   U     0      0    qr-2

router_102_103 的路由表内容如下：

    Destination Gateway Genmask Flags Metric Use Iface
    10.10.1.0 * 255.255.255.0   U     0      0     qr-3
    10.10.2.0 * 255.255.255.0   U     0      0     qr-4

这样的路由表是可以工作的。

例如 vlan102 上有数据包要发到 vlan103。选择路由时，会查看 router_102_103 的路由表, 匹配第二个条目，数据通过 qr-4

被正确地发送到 vlan103。

同样当 vlan100 上有数据包要发到 vlan101 时，会匹配 router_100_101 路由表的第二个条目，数据通过 qr-2 被正确地发送到 vlan101。

可见，namespace 使得每个 router 有自己的路由表，而且不会与其他 router 冲突，所以能很好地支持网络重叠。

# 其他

在 neutron 中，网络虚拟化(NFV)是必须的存在，Namespace 使得 openstack 中的每个 route 都有自己的路由表，都是一个网络隔离域，在各自的区域中就算配置相同的 IP 也不会冲突，很好的支持网络重叠。

网络虚拟化原理

- 在 openstack 中创建一个 Namespace，该命名空间相当于把虚拟网络放在一个独立的虚拟化环境中，可以在 linux 中通过 ip netns exec NAME COMMAND 进行操作，一个 Namespace 就是一整套完整独立的网络，包括路由交换等(创建的 route 就是放在 namespace 中)可以为系统之上的虚拟机提供全套网络服务
  - ![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oo0pgf/1616123235147-59333d41-b2a7-47d9-9fb2-298f751aff54.jpeg)
-

- 如图所示中间的路由器就是由 neutron(linux 也具备此功能)模拟出来的虚拟路由器，该虚拟路由器就存在于 Namespace 为 XXXXXX 的虚拟环境中(可以通过 neutron router-list 命令查看 Namespace 的 ID)，可以通过 ip nets exec NAME bash 打开该虚拟化网络环境，之后所有的操作，都相当于是对这个环境启动了一个 bash，所有操作都是对这个环境下的虚拟设备进行的
- 不同 VLAN 之间的通信，就是依靠 namespace 这个路由器来进行三层通信的
- 当内网虚拟机需要访问外网的时候，还可以创建一个外部网络放到该 Namespace，如图所示，当访问 vlan100 与 101 以外的所有流量都转发给 10.10.10.1(该网络相当于公网 IP)
- NAT，当数据到达 qg-b8b32a88-03 这个接口的时候，虚拟路由器会进行 NAT 转换，把 172 的地址段全部转换成 10 网段
- Floating IP，当外网想要访问内网机器的时候，就需要用到这功能，创建一个 Floating 后从外网地址池分配一个 IP 给他，然后可以把 Floating IP 与内网中的任意主机绑定，当 route 收到访问该 Floating IP 的数据包时，则进行 NAT 转换成改 FIP 对应的内网机器的 IP，然后把数据包转发给内网机器，实现外网访问内网，公有云中个人买的云主机想访问，就是这个原理
- linux Namespace：在所有的 linux 系统都一个 init 进程（即初始化进程），其 PID=1,而由于在不同的 namespace 中的进程都是彼此透明的，因此在不同的 namespace 中都可以有自己的 PID=1 的 init 进程，相应 Namespace 内的孤儿进程都将以该进程为父进程，当该进程被结束时该 Namespace 内所有的进程都会被结束。换句话说，在同一个 linux 系统中由于 namespace 的存在，可以允许 n 个相同的进程存在并互不干扰的运行。
