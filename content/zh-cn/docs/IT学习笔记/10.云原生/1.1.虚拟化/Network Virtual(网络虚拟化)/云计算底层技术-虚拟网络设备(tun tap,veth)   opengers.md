---
title: 云计算底层技术-虚拟网络设备(tun tap,veth)   opengers
---

# 虚拟网络设备(tun/tap,veth)介绍

Posted on September 30, 2017 by opengers in openstack
[openstack 底层技术-各种虚拟网络设备一(Bridge,VLAN)](https://opengers.github.io/openstack/openstack-base-virtual-network-devices-bridge-and-vlan/)
[openstack 底层技术-各种虚拟网络设备二(tun/tap,veth)](https://opengers.github.io/openstack/openstack-base-virtual-network-devices-tuntap-veth/)

> 第一篇文章介绍了 Bridge 和 VLAN，本文继续介绍 tun/tap，veth 等虚拟设备，除了 tun，其它设备都能在 openstack 中找到应用，这些各种各样的虚拟网络设备使网络虚拟化成为了可能

- [tun/tap](https://opengers.github.io/openstack/openstack-base-virtual-network-devices-tuntap-veth/#tuntap)
- [tap 设备作为虚拟机网卡](https://opengers.github.io/openstack/openstack-base-virtual-network-devices-tuntap-veth/#tap%E8%AE%BE%E5%A4%87%E4%BD%9C%E4%B8%BA%E8%99%9A%E6%8B%9F%E6%9C%BA%E7%BD%91%E5%8D%A1)
- [openvpn 中使用的 tun 设备](https://opengers.github.io/openstack/openstack-base-virtual-network-devices-tuntap-veth/#openvpn%E4%B8%AD%E4%BD%BF%E7%94%A8%E7%9A%84tun%E8%AE%BE%E5%A4%87)
- [veth 设备](https://opengers.github.io/openstack/openstack-base-virtual-network-devices-tuntap-veth/#veth%E8%AE%BE%E5%A4%87)
- [veth 设备在 openstack 中的应用](https://opengers.github.io/openstack/openstack-base-virtual-network-devices-tuntap-veth/#veth%E8%AE%BE%E5%A4%87%E5%9C%A8openstack%E4%B8%AD%E7%9A%84%E5%BA%94%E7%94%A8)
  tun/tap

我们知道 KVM 虚拟化中单个虚拟机是主机上的一个普通`qemu-kvm`进程，虚拟机当然也需要网卡，最常见的虚拟网卡就是使用主机上的 tap 设备。那从主机的角度看，这个`qemu-kvm`进程是如何使用 tap 设备呢，下面先介绍下`tun/tap`设备概念，然后分别用一个实例来解释`tun/tap`的具体用途`tun/tap`是操作系统内核中的虚拟网络设备，他们为用户层程序提供数据的接收与传输。实现`tun/tap`设备的内核模块为`tun`，其模块介绍为`Universal TUN/TAP device driver`，该模块提供了一个设备接口`/dev/net/tun`供用户层程序读写，用户层程序通过读写`/dev/net/tun`来向主机内核协议栈注入数据或接收来自主机内核协议栈的数据，**可以把 tun/tap 看成数据管道，它一端连接主机协议栈，另一端连接用户程序**

    [root@compute01 ~]# modinfo tun
    filename:       /lib/modules/3.10.0-514.16.1.el7.x86_64/kernel/drivers/net/tun.ko
    alias:          devname:net/tun
    ...
    description:    Universal TUN/TAP device driver
    ...

    [root@compute01 ~]# ls /dev/net/tun
    /dev/net/tun

为了使用`tun/tap`设备，用户层程序需要通过系统调用打开`/dev/net/tun`获得一个读写该设备的文件描述符(FD)，并且调用 ioctl()向内核注册一个 TUN 或 TAP 类型的虚拟网卡(实例化一个 tun/tap 设备)，其名称可能是`tap7b7ee9a9-c1/vnetXX/tunXX/tap0`等。此后，用户程序可以通过该虚拟网卡与主机内核协议栈交互。当用户层程序关闭后，其注册的 TUN 或 TAP 虚拟网卡以及路由表相关条目(使用 tun 可能会产生路由表条目，比如 openvpn)都会被内核释放。可以把用户层程序看做是网络上另一台主机，他们通过 tap 虚拟网卡相连 TUN 和 TAP 设备区别在于他们工作的协议栈层次不同，TAP 等同于一个以太网设备，用户层程序向 tap 设备读写的是二层数据包如以太网数据帧，tap 设备最常用的就是作为虚拟机网卡。TUN 则模拟了网络层设备，操作第三层数据包比如 IP 数据包，`openvpn`使用 TUN 设备在 C/S 间建立 VPN 隧道

# tap 设备关联虚拟机网卡

tap 设备最常见的用途就是关联虚拟机网卡，这里以一个具体的虚拟机 P2 为例，来看看它与其所在主机，及网络上其它主机通信的数据流向

本文依然使用下面这张图作为参考

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/casm7q/1616124199627-2d2cd67b-8e31-44f4-a932-6e1925d09376.png)

虚拟机 P2 使用桥接模式，网桥为 br0，先启动虚拟机

    #virsh domiflist P2
    Interface  Type       Source     Model       MAC
    -------------------------------------------------------
    tap0       bridge     br0        virtio      fa:16:3e:b1:70:52
    ...

虚拟机启动后，主机上多了一张虚拟网卡`tap0`，`tap0`桥接在 br0 网桥上，查看此虚拟机进程

    #ps -ef |grep P2
    qemu      7748     1  0 Nov07 ?        00:22:09 /usr/libexec/qemu-kvm -name guest=P2 ... -netdev tap,fd=26,id=hostnet0,vhost=on,vhostfd=28 ...

进程 PID 为 7748，其网络部分参数中，`-netdev tap,fd=26` 表示连接主机上的 tap 设备，`fd=26`为读写`/dev/net/tun`的文件描述符。使用`lsof -p 7748`也可以验证，如下，PID 为 7748 的进程打开了文件/dev/net/tun，分配的文件描述符为 26，打开设备文件类型为 CHR

    # lsof -p 7748
    COMMAND    PID USER   FD      TYPE             DEVICE    SIZE/OFF     NODE NAME
    ...
    qemu-kvm 7748  qemu   26u      CHR             10,200         0t0    17439 /dev/net/tun
    ...

因此，在虚拟机 P2 启动时，其打开了设备文件`/dev/net/tun`并获得了读写该文件的文件描述符(FD)26，同时向内核注册了一个 tap 类型虚拟网卡`tap0`，`tap0`与 FD 26 关联，虚拟机关闭时`tap0`设备会被内核释放。此虚拟网卡`tap0`一端连接用户空间程序`qemu-kvm`，另一端连接主机链路层

**从虚拟机 P2 发送数据到外部网络**

1. 虚拟机通过其网卡 eth0 向外发送数据，从主机角度看，就是用户层程序`qemu-kvm`进程使用文件描述符(FD)26 向字符设备`/dev/net/tun`写入数据 `P2 --> write(fd,...) --> N`
2. 文件描述符 26 与虚拟网卡`tap0`关联，也就是说主机从`tap0`网卡收到数据 `N --> E`
3. `tap0`为网桥`br0`上一个接口，需要进行`Bridging decision`以决定数据包如何转发 `E --> D`
4. P2 是与外部网络其它主机通信，因此`br0`转发该数据到`em2.100`，最后从物理网卡 em2 发出 `D --> C --> B -- > A`

这个过程中，虚拟机发出的数据通过`tap0`虚拟网卡直接注入主机链路层网桥处理逻辑中，然后被转发到外部网络。可以看出数据包没有穿过主机协议栈上层，主机仅仅起了类似物理二层交换机的数据转发功能，这是当然的也是必须的，因为虚拟机通信对象不是此主机，明白这点之后，可以对 tap 设备有更深的认识如下图，Linux 上的各种网络应用程序基本上都是通过 Linux Socket 编程接口来和内核空间的网络协议栈通信的。那么问题来了，为何虚拟机进程`qemu-kvm`不能使用 Linux Socket 与外部网络通信呢？

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/casm7q/1616124199610-b6861cf1-bdbb-4002-8507-8de6b84cf44d.png)

需要注意的是，虽然`qemu-kvm`只是主机上的一个进程，但它实现的是一台虚拟机，虚拟机和主机是两台机器，他们都桥接在主机上的软件交换机上。明白了这点，我们就知道主机上普通的网络程序发出的网络通信是属于”主机自身发出的数据包”，而虚拟机发出的数据包当然不可能使用另一台机器(主机)的 Linux Socket 来通信了。

主机上能看到的是虚拟机发到 tap 设备上的二层以太网帧，因此主机上工作在内核协议栈 IP 层的 iptables 是无法过滤虚拟机数据包的，当然这也有解决方法，本系列第一篇文章有详细说明

# openvpn 中使用的 tun 设备

TUN 设备与上面介绍的 TAP 设备很类似，只是 TUN 设备连接的是主机内核协议栈 IP 层

openvpn 是使用 TUN 设备的一个常见例子，搭建好 openvpn server 后启动 openvpn，主机中多了一个虚拟网卡`tun0`

    #ifconfig tun0
    tun0      Link encap:UNSPEC  HWaddr 00-00-00-00-00-00-00-00-00-00-00-00-00-00-00-00
              inet addr:10.5.0.1  P-t-P:10.5.0.2  Mask:255.255.255.255
    ...

查看该虚拟网卡驱动信息，可以看到`tun0`网卡使用内核模块`tun`，类型为`tun`设备(区别于 tap 设备的`bus-info: tap`)

\#ethtool -i tun0 driver: tun … bus-info: tun …

    使用`lsof`查看openvpn进程打开的所有文件
    ``` shell
    # ps -ef |grep openvpn
    root      3056     1  4 Nov14 ?        03:18:28 /usr/sbin/openvpn --daemon --config /etc/openvpn/config/server.conf
    #lsof -p 3056
    COMMAND  PID USER   FD   TYPE DEVICE SIZE/OFF     NODE NAME
    ...
    openvpn 3056 root    6u   CHR 10,200      0t0     8590 /dev/net/tun
    openvpn 3056 root    7u  IPv4  26744      0t0      UDP *:9112
    ...

因此，openvpn 进程在启动时，打开了字符设备`/dev/net/tun`并获得了文件描述符 6，同时向内核注册了一个虚拟网卡`tun0`。openvpn 进程还打开了一个 udp 的 socket，监听在`9112`端口。 下面我们来看下用户通过 vpn 访问 web 服务的数据包流向

**客户端使用 openvpn 访问 web 服务**

1. 客户端启动 openvpn client 进程连接 openvpn server
2. server 下发路由条目到客户端机器路由表中，同时生成虚拟网卡`tun1`(tun 设备，openvpn client 进程与 openvpn server 一样会注册 tun 虚拟网卡)
3. 客户端通过浏览器访问 web 服务
4. 浏览器生成的数据包在协议栈 IP 层进行路由选择，决定通过虚拟网卡`tun1`发出
5. 虚拟网卡`tun1`另一端连接用户层 openvpn client 进程
6. openvpn client 进程收到原始请求数据包
7. openvpn client 封装原始请求数据包，通过 udp 协议发送 vpn 封包到 openvpn server 上的 9112 端口 `A -- > T --> K --> R --> P1`
8. openvpn server 上的 openvpn 进程收到 vpn 封包，解包，使用文件描述符 6 写数据到`/dev/net/tun` `P1 --> write(fd,...) --> N`
9. 文件描述符 6 与虚拟网卡`tun0`关联，主机从`tun0`网卡收到数据包 `N --> H ---> I`
10. 主机进行`Routing decision`，根据数据包目的 IP(用户访问 web 网站 IP 地址)从相应网卡发出 `I --> K --> T --> M --> A`

**如何创建 tun 设备**

    #添加一个tun设备，并配置ip
    ip tuntap add mode tun2
    ip link set tun2 up
    ip addr add 172.16.1.2/24 dev tun2

    #删除tun2
    ip link del dev tun2

# veth 设备

veth 也是 Linux 实现的虚拟网络设备，veth 设备总是成对出现，其作用是反转数据流方向。例如如果 veth-a 和 veth-b 是一对 veth 设备，veth-a 收到的数据会从 veth-b 发出。相反，veth-b 收到的数据会从 veth-a 发出。一个常见用途是连接两个 netwok namespace。openstack Neutron 网络中，dhcp agent 和 l3 agent 都用到了 veth 设备对。拿 dhcp agent 来说，openstack 中每个网络都需要起一个 dhcp 服务器用于给此网络虚拟机分配 ip 地址，每个 openstack 网络都使用一个单独的 network namespace，每个 network namespace 和网络节点上 Bridge 通过 veth 设备对连接，这样多个 openstack 网络才不会引起冲突。理解 veth 设备对和 network namespace 是解决 openstack 中虚拟机网络故障问题的关键

使用 ip 命令新建一对 veth 设备 veth-a veth-b

    ip link add veth-a type veth peer name veth-b

如何确定一个网络设备是 veth 设备，如下`driver: veth`

    ethtool -i veth-a
    driver: veth
    ...

如何查找`veth-a`的对端设备呢

    ip -d link show
    7: veth-b@veth-a: <BROADCAST,MULTICAST,M-DOWN> mtu 1500 qdisc noop state DOWN mode DEFAULT qlen 1000
        link/ether 82:e6:4f:09:b2:df brd ff:ff:ff:ff:ff:ff promiscuity 0
        veth addrgenmode eui64
    8: veth-a@veth-b: <BROADCAST,MULTICAST,M-DOWN> mtu 1500 qdisc noop state DOWN mode DEFAULT qlen 1000
        link/ether d6:4b:5e:69:ff:2d brd ff:ff:ff:ff:ff:ff promiscuity 0
        veth addrgenmode eui64

若 veth 设备对其中一个设备位于网络命令空间中，可以这样查找

    #查看网络命名空间qdhcp-ba48a3fc-e9e8-4ce0-8691-3d35b6cca80a中设备ns-45400fa0-9c的对端设备序号
    ip netns exec qdhcp-ba48a3fc-e9e8-4ce0-8691-3d35b6cca80a ethtool -S ns-45400fa0-9c
    NIC statistics:
         peer_ifindex: 6
    #查看序号为6的网络设备
    ip a |grep '^6:'
    6: tap45400fa0-9c@if2: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue master brqba48a3fc-e9 state UP qlen 1000

因此网络命名空间中的`ns-45400fa0-9c`与主机上的设备`tap45400fa0-9c`是一对 veth 设备

# veth 设备在 openstack 中的应用

如下，我们使用 veth 设备对连接 Bridge 和 network namespace，凡是桥接到 Bridge`brqba48a3fc-e9`上的虚拟机，其发出的 dhcp 请求都会进入网络命名空间`qdhcp-ba48a3fc-e9e8-4ce0-8691-3d35b6cca80a`中的 dhcp 服务器。

    #新建veth设备veth-a,veth-b
    ip link add veth-a type veth peer name veth-b

    #新建一个network namespace
    ip netns add qdhcp-ba48a3fc-e9e8-4ce0-8691-3d35b6cca80a

    #将 veth-b 添加到 network namespace
    ip link set veth-b netns qdhcp-ba48a3fc-e9e8-4ce0-8691-3d35b6cca80a

    #给命名空间中的veth设备 veth-b 设置ip地址， 此IP地址做为同网段虚拟机的dhcp服务器
    ip netns exec qdhcp-ba48a3fc-e9e8-4ce0-8691-3d35b6cca80a ip addr add 10.0.0.2/32 dev veth-b
    ip netns exec qdhcp-ba48a3fc-e9e8-4ce0-8691-3d35b6cca80a ip link set dev veth-b up

    #在nstest命名空间中起一个dnsmasq服务器
    ip netns exec qdhcp-ba48a3fc-e9e8-4ce0-8691-3d35b6cca80a /usr/sbin/dnsmasq --no-hosts --no-resolv --except-interface=lo  --bind-interfaces --interface=veth-b ...

    #查看此命令空间中服务监听端口
    ip netns exec qdhcp-ba48a3fc-e9e8-4ce0-8691-3d35b6cca80a netstat -tnpl
    Active Internet connections (only servers)
    Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name
    tcp        0      0 10.0.0.2:53          0.0.0.0:*               LISTEN      28795/dnsmasq

    #新建网桥brqba48a3fc-e9
    brctl add brqba48a3fc-e9
    #将veth-a添加到网桥 brqba48a3fc-e9
    brctl addif brqba48a3fc-e9 veth-a

    #新建kvm虚拟机过程忽略，查看此虚拟机网卡设备信息
    virsh domiflist testvm
    Interface  Type       Source     Model       MAC
    -------------------------------------------------------
    tap796091c0-07 bridge     brqba48a3fc-e9 virtio      fa:16:3e:85:f8:c

如上步骤，这样虚拟机发出的 dhcp 请求包会被网桥`brqba48a3fc-e9`转发到接口`veth-a`，由于 veth-a 与 veth-b 是一对 veth 设备，因此数据包会到达网络命名空间中的 veth-b，也即虚拟机 dhcp 请求能够成功我们上面是直接在网桥`brqba48a3fc-e9`上建立虚拟机测试 dhcp 请求，也可以把主机上一块物理网卡添加到此网桥，这样此物理网卡所在网络上的其它所有主机都能使用此 dhcp 服务器。并且，仿照上面步骤，我们可以建多个 Bridge，每个 Bridge 都建立有自己的 network namespace 和 dhcp 服务器。这其实就是 openstack Bridge+Vlan(无 l3 agent)网络模式下，虚拟机自动获取 IP 的机制，如下图

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/casm7q/1616124199638-9318f4cd-d1a2-4ff2-ae9b-4fc0cbdfb562.png)

(本文完)
