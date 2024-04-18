---
title: 虚拟网络设备(Bridge,VLAN)
linkTitle: 虚拟网络设备(Bridge,VLAN)
date: 2024-04-19T13:03
weight: 20
---



# 概述

> 参考：
>
> -


# 云计算底层技术-虚拟网络设备(Bridge,VLAN)

Posted on September 24, 2017 by opengers in openstack

原文链接：[openstack 底层技术-各种虚拟网络设备一(Bridge,VLAN)](https://opengers.github.io/openstack/openstack-base-virtual-network-devices-bridge-and-vlan/)

> IBM 网站上有一篇高质量文章。本文会参考文章部分内容，本系列介绍 OpenStack 使用的这些网络设备包括 Bridge，VLAN，tun/tap, veth，vxlan/gre。本篇先介绍 Bridge 和 VLAN 相关，其它在下一篇中介绍

OpenStack 一般分为计算，存储，网络三部分。考虑构建一个灵活的可扩展的云网络环境，而物理网络架构一般是固定和难于扩展的，因此虚拟网络将更有优势。Linux 平台上实现了各种不同功能的虚拟网络设备，包括`Bridge,Vlan,tun/tap,veth pair,vxlan/gre，...`，这些虚拟设备就像一个个积木块一样，被 OpenStack 组合用于构建虚拟网络。 还有火热的 Docker，docker 容器的隔离技术实现脱胎于 Linux 平台上的`namspace`,以及更早的`chroot`。

文中会牵涉虚拟机，所以文中出现的”主机”一词明确表示一台物理机，”接口”指挂载到网桥上的网络设备，环境如下：

```
CentOS Linux release 7.3.1611 (Core)
Linux controller 3.10.0-514.16.1.el7.x86_64 #1 SMP Wed Apr 12 15:04:24 UTC 2017 x86_64 x86_64 x86_64 GNU/Linux

OpenStack社区版 Newton
```

Linux Bridge

内核模块`bridge`

```bash
[root@controller ~]# modinfo bridge
filename:       /lib/modules/3.10.0-514.16.1.el7.x86_64/kernel/net/bridge/bridge.ko
```

Bridge 是 Linux 上工作在内核协议栈二层的虚拟交换机，虽然是软件实现的，但它与普通的二层物理交换机功能一样。可以添加若干个网络设备(em1,eth0,tap,..)到 Bridge 上(`brctl addif`)作为其接口，添加到 Bridge 上的设备被设置为只接受二层数据帧并且转发所有收到的数据包到 Bridge 中(bridge 内核模块)，在 Bridge 中会进行一个类似物理交换机的查 MAC 端口映射表，转发，更新 MAC 端口映射表这样的处理逻辑，从而数据包可以被转发到另一个接口/丢弃/广播/发往上层协议栈，由此 Bridge 实现了数据转发的功能。如果使用`tcpdump`在 Bridge 接口上抓包，是可以抓到桥上所有接口进出的包跟物理交换机不同的是，运行 Bridge 的是一个 Linux 主机，Linux 主机本身也需要 IP 地址与其它设备通信。但被添加到 Bridge 上的网卡是不能配置 IP 地址的，他们工作在数据链路层，对路由系统不可见。不过 Bridge 本身可以设置 IP 地址，可以认为当使用`brctl addbr br0`新建一个`br0`网桥时，系统自动创建了一个同名的隐藏`br0`网络设备。`br0`一旦设置 IP 地址，就意味着`br0`可以作为路由接口设备，参与 IP 层的路由选择(可以使用`route -n`查看最后一列`Iface`)。因此只有当`br0`设置 IP 地址时，Bridge 才有可能将数据包发往上层协议栈。

根据下图来具体分析下 Bridge 工作过程

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rg120l/1616124222069-97b6dbd5-5adf-4cd6-811e-db15669b7bba.png)

上图主机有 em1 和 em2 两块网卡，有网桥`br0`。用户空间进程有 app1，app2 等普通网络应用，还有 OpenVPN 进程 P1，以及一台或多台 kvm 虚拟机 P2(kvm 虚拟机实现为主机上的一个`qemu-kvm`进程，下文用`qemu-kvm`进程表示虚拟机)。此主机上使用到了多种虚拟网络设备，在具体介绍某个虚拟网络设备时，我们可以忽略其它网络设备工作细节，只专注于当前网络设备。下面来具体分析网桥`br0`

**Bridge 处理数据包流程**

图中可以看到`br0`有 N 个`tap`类型接口(tap0,..,tapN)，tap 设备名称可能不同，例如`tap45400fa0-9c`或`vnet*`，但都是 tap 设备。一个”隐藏”的`br0`接口(可设置 IP)，以及物理网卡 em2 的一个 VLAN 子设备`em2.100`(这里简单看作有一个网卡桥接到 br0 上即可，VLAN 下面会讲)，他们都工作在链路层(Link Layer)。来看数据从外部网络(A)发往虚拟机(P2)`qemu-kvm`这一过程，首先数据包从 em2(B)物理网卡进入，之后 em2 将数据包转发给其 vlan 子设备 em2.100，经过`Bridge check`(L)发现子设备`em2.100`属于网桥接口设备，因此数据包不会发往协议栈上层(T),而是进入 bridge 代码处理逻辑，从而数据包从`em2.100`接口(C)进入`br0`，经过`Bridging decision`(D)发现数据包应当从`tap0`(E)接口发出，此时数据包离开主机网络协议栈(G)，发往被用户空间进程`qemu-kvm`打开的字符设备`/dev/net/tun`(N)，`qemu-kvm`进程执行系统调用`read(fd,...)`从字符设备读取数据。 这个过程中，外部网络 A 发出的数据包是不会也没必要进入主机上层协议栈的，因为 A 是与主机上的 P2 虚拟机通信，**主机只是起到一个网桥转发的作用**作为网桥的对比，如果是从网卡 em1(M)进入主机的数据包，经过`Bridge check`(L)后，发现 em1 非网桥接口，则数据包会直接发往(T)协议栈 IP 层,从而在`Routing decision`环节决定数据包的去向(A –> M –> T –> K)

**Bridging decision**

上图中网桥`br0`收到数据包后，根据数据包目的 MAC 的不同，`Bridging decision`环节(D)对数据包的处理有以下几种：

- 包目的 MAC 为 Bridge 本身 MAC 地址(当`br0`设置有 IP 地址)，从 MAC 地址这一层来看，收到发往主机自身的数据包，交给上层协议栈(D –> J)

- 广播包，转发到 Bridge 上的所有接口(br0,tap0,tap1,tap…)

- 单播&&存在于 MAC 端口映射表，查表直接转发到对应接口(比如 D –> E)

- 单播&&不存在于 MAC 端口映射表，泛洪到 Bridge 连接的所有接口(br0,tap0,tap1,tap…)

- 数据包目的地址接口不是网桥接口，桥不处理，交给上层协议栈(D –> J)

# Bridge 与 netfilter

Linux 防火墙是通过`netfiler`这个内核框架实现，`netfiler`用于管理网络数据包。不仅具有网络地址转换(NAT)的功能，也具备数据包内容修改、以及数据包过滤等防火墙功能。利用运作于用户空间的应用软件，如 iptables/firewalld/ebtables 等来控制`netfilter`。Netfilter 在内核协议栈中指定了五个处理数据包的钩子(hook)，分别是 PRE_ROUTING、INPUT、OUTPUT、FORWARD 与 POST_ROUTING，通过 iptables/firewalld/ebtables 等用户层工具向这些 hook 点注入一些数据包处理函数，这样当数据包经过相应的 hook 时，处理函数就被调用，从而实现包过滤功能。这些用户层工具中，iptables 工作在 IP 层，只能过滤 IP 数据包；ebtables 工作在数据链路层，只能过滤以太网帧(比如更改源或目的 MAC 地址)当主机上没有 Bridge 存在时，从网卡进入主机的数据包会依次穿过主机内核协议栈，最后到达应用层交给某个应用程序处理。这样我们可以很方便的使用 iptables 设置本主机的防火墙规则。进入数据包流向对应下图路径`A --> L --> T --> ...`

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rg120l/1616124222054-993994b9-971c-4bdd-96ef-ad7ee716c329.png)

Bridge 的出现使 Linux 上设置防火墙变得复杂，因为此时从物理网卡进入主机的数据包目的地可能是其上运行的一台虚拟机。上图是上面介绍的`数据从外部网络(A)发往虚拟机(P2)`这一过程中数据包所经过的防火墙链(文中的两张图可以对比来看)。物理网卡 em2 子设备 em2.100 从外部网络 A 收到二层数据包，经过`bridge check`后进入 br0 并穿越一系列防火墙链`L --> D --> E`，最终从 Bridge 上的另一个接口`tap0`发出。上图红色导向线可以很清楚看到整个过程中数据包是没有进入主机内核协议栈的，因此位于主机 IP 层(Network Layer)的 iptables 根本无法过滤`L --> D --> E`这一路径的数据包。那有没有办法**使 iptables 能够过滤 Bridge 中的数据包呢?**

- ebtables 只可以简单过滤二层以太网帧，无法过滤 ipv4 数据包。

- 当然也可以在虚拟机内使用 iptables，但是一般不会这么玩，特别是在云平台环境。 一个原因是如果一台主机上运行有 10 台虚拟机，那就需要分别登录这 10 台虚拟机设置其 iptables 规则，工作量会多很多。而且这么做意味着云平台必须要能够登录用户的虚拟机来设置 iptables 规则。 OpenStack 中安全组也是 iptables 实现，我们在虚拟机内部并没有发现有 iptables 规则存在。

- 解决办法就是下文要讲的`bridge_netfilter`

来自：[Bridge-nf Frequently Asked Questions](http://ebtables.netfilter.org/misc/brnf-faq.html)

为了解决上面提到的问题，Linux 内核引入了`bridge_netfilter`，简称`bridge_nf`。`bridge_netfilter`在链路层 Bridge 代码中插入了几个能够被 iptables 调用的钩子函数，Bridge 中数据包在经过这些钩子函数时，iptables 规则被执行(上图中最下层 Link Layer 中的绿色方框即是 iptables 插入到链路层的 chain,蓝色方框为 ebtables chain)。这就使得{ip,ip6,arp}tables 能够”看见”Bridge 中的 IPv4,ARP 等数据包。这样不管此数据包是发给主机本身，还是通过 Bridge 转发给虚拟机，iptables 都能完成过滤。

**如何使用 bridge_nf**

从 Linux 2.6.1 内核开始，可以通过设置内核参数开启`bridge_netfilter`机制。看名字就很容易知道具体作用

    [root@controller ~]# sysctl -a |grep 'bridge-nf-'
    net.bridge.bridge-nf-call-arptables = 1
    net.bridge.bridge-nf-call-ip6tables = 1
    net.bridge.bridge-nf-call-iptables = 1
    ...

然后在 iptables 中使用`-m physdev`引入相应模块，以文中第一张图上的虚拟机 P2 为例，它的虚拟网卡`tap0`桥接在`br0`上。我们在主机上设置如下 iptables 规则：丢弃从网桥`br0`的`tap0`接口进入的数据包。

    #查看网桥
    # brctl show
    bridge name     bridge id               STP enabled     interfaces
    br0             8000.f8bc1212c3a0       no              em1
                                                            tap
    #操作对象是tap0
    #ptables -t raw -A PREROUTING -m physdev --physdev-in tap0  -j DROP

注意 iptables `-m physdev`操作对象是 Bridge 上的某个接口，因此规则的有效范围是针对从此接口进出 Bridge 的数据包还需要注意一点，这条命令是在主机上执行的，从主机角度看，主机收到从`tap0`接口进入的数据包，因此使用`--physdev-in`。但是从虚拟机 P2 角度来看，它发出了数据包，发给主机上的网桥 br0。因此上面这条 iptables 命令实际的作用是丢弃虚拟机 P2 发出的数据包，也就是禁止虚拟机 P2 访问外网。方向要分清，后面讲到/tun/tap 设备时会细说

上面介绍了使用 iptables 过滤 Bridge 中数据包的方法，实际中如果直接使用 iptables 命令显然太繁琐，而且如果主机上有多台虚拟机的话，网桥接口就会变多导致容易出错，需要依靠工具去做这些。

**使用 libvirt 提供的 virsh 工具**

libvirt 的`virsh nwfilter-*`系列命令提供了设置虚拟机防火墙的功能，它其实是封装了 iptables 过滤 Bridge 中数据包的命令(`-m physdev`)。它使用多个 xml 文件，每个 xml 文件中都可以定义一系列防火墙规则，然后把某个 xml 文件应用到某虚拟机的某张网卡(Bridge 中的接口)，这样就完成了对此虚拟机的这张网卡的防火墙设置。当然可以把一个定义好防火墙规则的 xml 文件应用到多台虚拟机。

    #查看用于设置
    # virsh --help |grep nwfilter
        nwfilter-define                define or update a network filter from an XML file
        nwfilter-dumpxml               network filter information in XML
        nwfilter-edit                  edit XML configuration for a network filter
        nwfilter-list                  list network filters
        nwfilter-undefine              undefine a network filter

    # 定义有防火墙规则的xml文件
    # virsh nwfilter-dumpxml centos6.3_filter
    <filter name='centos6.3_filter' chain='root'>
      <uuid>b1fdd87c-a44c-48fb-9a9d-e30f1466b720</uuid>
      <rule action='accept' direction='in' priority='400'>
        <tcp dstportstart='8000' dstportend='8002'/>
      </rule>
      <rule action='accept' direction='in' priority='400'>
        <tcp srcipaddr='172.16.1.0' srcipmask='24'/>
      </rule>
    </filter>

    #查看定义的防火墙xml文件
    # virsh nwfilter-list
     UUID                                  Name
    ------------------------------------------------------------------
     69754f43-0325-453f-bd53-4a6e3ab5d456  centos6.3_filter

    # 在虚拟机xml文件中应用centos6.3_filter
        <interface type='bridge'>
          <mac address='f8:c9:79:ce:60:01'/>
          <source bridge='br0'/>
          <model type='virtio'/>
          <filterref filter='centos6.3_filter'/>
          <address type='pci' domain='0x0000' bus='0x00' slot='0x03' function='0x0'/>
        </interface>

**openstack 中的安全组**

像 openstack 等很多云平台在 web 控制台中都会提供有设置虚拟机防火墙功能(安全组)，可以很方便的添加应用防火墙规则到云主机

在 OpenStack 部署中，若使用 Bridge 实现虚拟网络，其安全组功能就是依靠`bridge_nf`实现，计算节点上 iptables 才能”看见”并过滤发往其上 instance 的数据包，如下是 OpenStack 计算节点上部分 iptables 规则，当启用安全组时，OpenStack 会自动设置`net.bridge.bridge-nf-call-iptables = 1`等内核参数，不用再明确设置。 `tap10f15e45-aa`为该计算节点上某 instance 网卡(tap 设备)

    ...
    #针对虚拟网卡tap10f15e45-aa的部分规则
    -A neutron-filter-top -j neutron-linuxbri-local
    -A neutron-linuxbri-FORWARD -m physdev --physdev-out tap10f15e45-aa --physdev-is-bridged -m comment --comment "Direct traffic from the VM interface to the security group chain." -j neutron-linuxbri-sg-chain
    -A neutron-linuxbri-FORWARD -m physdev --physdev-in tap10f15e45-aa --physdev-is-bridged -m comment --comment "Direct traffic from the VM interface to the security group chain." -j neutron-linuxbri-sg-chain
    ...

Bridge+netfilter 内容很多，下次有时间会专门用一篇文章介绍 OpenStack 中的安全组实现，关键字 `iptables+bridge+netfilter+conntrack`

Linux 上还有一款虚拟交换机 OVS，主要区别是 OVS 支持 vlan tag 以及流表(例如 openflow)等一些高级特性，Bridge 只是单纯二层交换机也不支持 vlan tag，OVS 具体介绍参考这里[openstack 底层技术-使用 openvswitch](https://opengers.github.io/openstack/openstack-base-use-openvswitch/)

# VLAN

上面简单介绍过 Bridge 和 OVS 区别，要在 Linux 上实现一个带 VLAN 功能的虚拟交换机，OVS 可以通过给不同 port 打不同 tag 实现 vlan 功能，而 Bridge 需要结合 VLAN 设备才能实现。

这部分先介绍 VLAN 设备原理及配置，然后介绍 VLAN 在 openstack 中的应用

### VLAN 设备原理及配置

VLAN 又称虚拟网络，其基本原理是在二层协议里插入额外的 VLAN 协议数据（称为 802.1.q VLAN Tag)，同时保持和传统二层设备的兼容性。Linux 里的 VLAN 设备是对 802.1.q 协议的一种内部软件实现，模拟现实世界中的 802.1.q 交换机。详细介绍参考文章开头给出的 IBM 文章中”VLAN device for 802.1.q”部分，这里不再重复

下面使用 VLAN 结合 Bridge 实现文中第一张图上的多个 VLAN 子设备`VLAN 100, VLAN X, VLAN Y, ...`以及多个网桥`br0, brX, brY, ...`(X Y 都是小于 2048 的正整数)。前提是此主机上有一块网卡设备比如 em2，不管 em2 为物理网卡或虚拟网卡，em2 所连接的交换机端口必须设置为 trunk。添加 VLAN 子设备`VLAN 100, VLAN X, VLAN Y, ...`

    #cat /etc/sysconfig/network-scripts/ifcfg-em2
    TYPE=Ethernet
    BOOTPROTO=none
    IPV4_FAILURE_FATAL=no
    NAME=em2
    UUID=4f2cfd28-ba78-4f25-afa1-xxxxxxxxxxxxx
    DEVICE=em2
    ONBOOT=yes

    #添加vlan子设备em2.100
    #cat /etc/sysconfig/network-scripts/ifcfg-em2.100
    DEVICE=em2.100
    BOOTPROTO=static
    ONBOOT=yes
    VLAN=yes

    #可以继续添加多个带不同vlan tag的子设备,比如VLAN X, VLAN Y, VLAN ...
    #vlan子设备em2.X
    #cat /etc/sysconfig/network-scripts/ifcfg-em2.X
    #DEVICE=em2.X
    #BOOTPROTO=static
    #ONBOOT=yes
    #VLAN=yes

    #vlan子设备em2.Y
    #cat /etc/sysconfig/network-scripts/ifcfg-em2.Y
    #DEVICE=em2.Y
    #BOOTPROTO=static
    #ONBOOT=yes
    #VLAN=yes

    #...

    #重启网络
    service network restart

查看子设备 em2.100，可以看到，其 driver 为 VALN

    #ethtool -i eth1.101
    driver: 802.1Q VLAN Support
    version: 1.8
    ...

添加网桥`br0, brX, brY, ...`

    #添加br0网桥
    brctl addbr br0
    # br0添加em2.100子设备 凡是桥接到br0上的数据包都会带上tag 100
    brctl addif br0 em2.100

    #可以继续添加多个网桥brX, brY, ...
    # brctl addbr brX
    #凡是桥接到brX上的数据包都会带上tag X
    # brctl addif brX em2.X

    # brctl addbr brY
    # brctl addif brY em2.Y

    #...

    # brctl show
    bridge name     bridge id               STP enabled     interfaces
    br0         8000.525400315e23       no                em2.100
                   tap0

VLAN 设备的作用是建立一个个带不同 vlan tag 的子设备，它并不能建立多个可以交换转发数据的接口，因此需要借助于 Bridge，把 VLAN 建立的子设备例如 em2.100 桥接到网桥例如 br0 上，这样凡是桥接到 br0 上的设备就自动加入了 vlan 100 子网。对比一台带有两个 vlan 100，X 的物理交换机，这里 br0 网桥上所连接的接口相当于物理交换机上那些划分到 vlan 100 的端口，而 brX 所连接的接口相当于物理交换机上那些划分到 vlan X 的端口。因此 Bridge 加 VLAN 能在功能层面完整模拟现实世界里的 802.1.q 交换机。

参考文中第一张图，我们从网桥 br0 上 tap0 接口(E)角度，来看下具体的数据收发流程：

- 数据从 tap0 接口进入，发往外部网络 A

`tap0收到的数据被发送给br0(E) --> D --> br0把数据从em2.100接口发出(C) --> 母设备em2收到em2.100发来的数据(B) --> 母设备em2给数据打上100的vlan tag(因为来自em2.100) --> em2将带有100 tag的数据发出到外部网络(A) --> em2所连接的交换机收到数据(trunk口)`

- 数据从 em2 网卡进入，发往 tap0

`em2从所连接的交换机收到tag 100的数据 --> em2发现此数据带有tag 100,移除数据包中tag --> 不带tag的数据发给em2.100子设备 --> br0收到从em2.100进入的数据包 --> D --> br0转发数据到tap0`

上面忽略了对于数据包是否带 tag，以及数据包所带 tag 的子设备是否存在的检查。这些属于 vlan 基础知识

### VLAN 在 openstack 中的应用

openstack 中虚拟机网络使用 VLAN 模式的话，就会用到 VLAN 设备。openstack 中配置 vlan+bridge 模式如下

    #neutron-server节点(网络节点)配置
    #cat /etc/neutron/plugins/ml2/ml2_conf.ini
    [ml2]
    #neutron-server启动时，加载flat，vlan两种网络类型驱动
    type_drivers = flat,vlan
    #vlan模式不需要tenant_network，留空
    tenant_network_types =
    #neutron-server启动时加载linuxbridge和openvswitch网桥驱动
    mechanism_drivers = linuxbridge,openvswitch

    [ml2_type_flat]
    # 在命令行或控制台新建flat类型网络时需要指定的名称，此名称会配置映射到计算节点上某块网卡，下面会设置
    flat_networks = proext

    [ml2_type_vlan]
    # 在命令行或控制台新建vlan类型网络时需要指定的名称，此名称会配置映射到计算节点上某块网卡，下面会设置
    network_vlan_ranges = provlan

    #重启neutron-server服务

    # 使用Bridge+vlan网络模式的nova-compute节点(计算节点)配置
    #cat /etc/neutron/plugins/ml2/linuxbridge_agent.ini
    [linux_bridge]
    #provlan名称映射到此计算节点eth2网卡，因为使用vlan模式，eth2需要设置为trunk
    #proext名称映射到此计算节点eth3网卡，我这个环境下eth3网卡为虚拟机连接外网接口
    physical_interface_mappings = provlan:eth2,proext:eth3

    #重启此计算节点nova-compute服务
    #配置中只需要指定vlan要用的母设备eth2，后续控制台新建带tag的网络时，neutron会自动建立eth2.{TAG}子设备并加入到网桥

在控制台新建一个 vlan tag 为 1023 的 Provider network: subvlan-1023，使用此 subvlan-1023 网络新建几台虚拟机，看下计算节点上网桥配置

    [root@compute03 neutron]# brctl show
    bridge name     bridge id               STP enabled     interfaces
    brq82405415-7a          8000.52540048b1a9       no              eth2.1023
                                                            tap10f15e45-aa
                                                            tapa659a214-b1
    brqf5808b72-44          8000.5254001ac83d       no              eth3
                                                            tapd3388a60-ae

    [root@compute03 neutron]# virsh domiflist instance-00000145
    Interface  Type       Source     Model       MAC
    -------------------------------------------------------
    tapa659a214-b1 bridge     brq82405415-7a virtio      fa:16:3e:bc:c9:e0

虚拟机`instance-00000145`的网卡`tapa659a214-b1`桥接到`brq82405415-7a`。跟上面介绍的类似，桥接到`brq82405415-7a`上的接口设备就自动加入了 vlan 1023 子网，因此从`instance-00000145`发出的数据包会带有 tag 1023(eth2.1023 的母设备 eth2 负责添加或移除 tag)

假如在控制台或命令行再新建一个 tag 为 1024 的子网，则网桥配置如下

    [root@compute03 neutron]# brctl show
    bridge name     bridge id               STP enabled     interfaces
    brq82405415-7a          8000.52540048b1a9       no              eth2.1023
                                                            tap10f15e45-aa
                                                            tapa659a214-b1
    brq7d59440b-cc          8000.525400aabbcc       no              eth2.1024
                                                            tap20ffafb2-1b
    brqf5808b72-44          8000.5254001ac83d       no              eth3
                                                            tapd3388a60-ae

    [root@compute03 neutron]# virsh domiflist instance-00000147
    Interface  Type       Source     Model       MAC
    -------------------------------------------------------
    tap20ffafb2-1b bridge     brq7d59440b-cc virtio      fa:16:3e:bd:12:40

虚拟机`instance-00000147`的网卡`tap20ffafb2-1b`桥接到`brq7d59440b-cc`,`instance-00000147`属于 vlan 1024 子网,这就实现了属于不同 vlan 子网的`instance-00000145`与`instance-00000147`的隔离性。他们虽然在同一台计算节点上，但彼此不互通，除非设置为两个 VLAN 可以互通感谢 Bridge 和 VLAN 设备，他们让 openstack 配置 vlan 网络成了可能，BUT!, Bridge+VLAN 不是唯一的选择，openstack 也支持 OVS，OVS 中是靠给不同 instance 接口打不同 tag 来实现 instance 的多 vlan 环境，OVS 模式除了配置部分跟 Bridge+VLAN 不同之外，使用上并没有什么区别，这里的设置`mechanism_drivers = linuxbridge,openvswitch`加载相应驱动，屏蔽掉了底层操作的差别与 Bridge 中`provlan，proext`映射到计算节点网卡的配置不同，OVS 配置文件中映射关系为 vlan 类型网络`provlan`映射到网桥`br-vlan`，flat 类型网络`proext`映射到网桥`br-ext`。至于`br-vlan`桥接 eth2 网卡，`br-ext`桥接 eth3 网卡则需要预先手动配置好，来看一个使用 OVS 的计算节点网桥

    [root@compute01 neutron]# ovs-vsctl show
    dd7ccaae-6a24-4d28-8577-9e5e6b5dfbd3
        Manager "ptcp:6640:127.0.0.1"
            is_connected: true
        Bridge br-ext
            Controller "tcp:127.0.0.1:6633"
                is_connected: true
            fail_mode: secure
            Port phy-br-ext
                Interface phy-br-ext
                    type: patch
                    options: {peer=int-br-ext}
            Port br-ext
                Interface br-ext
                    type: internal
            Port "eth3"
                Interface "eth3"
        Bridge br-vlan
            Controller "tcp:127.0.0.1:6633"
                is_connected: true
            fail_mode: secure
            Port br-vlan
                Interface br-vlan
                    type: internal
            Port phy-br-vlan
                Interface phy-br-vlan
                    type: patch
                    options: {peer=int-br-vlan}
            Port "eth2"
                Interface "eth2"
        Bridge br-int
            Controller "tcp:127.0.0.1:6633"
                is_connected: true
            fail_mode: secure
            Port int-br-ext
                Interface int-br-ext
                    type: patch
                    options: {peer=phy-br-ext}
            Port br-int
                Interface br-int
                    type: internal
            Port int-br-vlan
                Interface int-br-vlan
                    type: patch
                    options: {peer=phy-br-vlan}
            Port "qvo7d59440b-cc"
                tag: 1
                Interface "qvo7d59440b-cc"
        ovs_version: "2.5.0"
    [root@compute01 neutron]# brctl show
    bridge name     bridge id               STP enabled     interfaces
    qbr7d59440b-cc          8000.26d03016fcf6       no              qvb7d59440b-cc
                                                            tap7d59440b-cc

    #查看虚拟机网卡
    [root@compute01 neutron]# virsh domiflist instance-00000149
    Interface  Type       Source     Model       MAC
    -------------------------------------------------------
    tap7d59440b-cc bridge     qbr7d59440b-cc virtio      fa:16:3e:12:ba:e6

`instance-00000149`出数据流向为`tap7d59440b-cc --> qbr7d59440b-cc --> qvo7d59440b-cc(tag 1) --> br-int --> br-vlan --> eth2`。qvb7d59440b-cc 与 qvo7d59440b-cc 为一对 veth 设备这其中牵涉 OVS 流表和 OVS 内外部 tag 转换问题，又足够写一篇文章来介绍了，本文暂不继续介绍。还有一点，在使用 OVS 做网桥的同时又开启安全组功能时，会多出一个 Bridge 网桥用于设置安全组，如上面的`qbr7d59440b-cc`, 因为目前 iptables 不支持 OVS，只能在虚拟机与 OVS 网桥之间加进一个 Bridge 网桥用于设置 iptables 规则

其它网络设备会在另一篇文章介绍，本文完
