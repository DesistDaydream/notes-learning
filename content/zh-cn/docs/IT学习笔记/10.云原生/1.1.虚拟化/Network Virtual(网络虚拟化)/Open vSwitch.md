---
title: Open vSwitch
---

# Open vSwitch # 开放的虚拟交换机

特性：支持 802.1q，trunk，access；支持网卡绑定技术(NIC Bonding);支持 QoS 配置及策略；支持 GRE 通用路由封装；支持 VxLAN；等等

OVS 的组成部分：

1. ovs-vswitchd #守护进程,实现数据报文交换功能，和 Linux 内核兼容模块一同实现了基于流的交换技术

2. ovsdb-server #ovs 的数据库，轻量级的数据库服务器，主要保存了 OVS 的配置信息，EXP 接口、交换、VLAN 等，ovs-vswitchd 的交换功能基于此库实现,相关数据信息保存在这个文件中：/etc/openvswitch/conf.db

3. ovs-vsctl #命令行工具，用于获取或更改 ovs-vswitchd 的配置信息，其修改操作会保存至 ovsdb-server 中

4. ovs-dpctl #

5. ovs-appctl #

6. ovsdbmonitor #

7. ovs-controller #

8. ovs-ofctl #

9. ovs-pki #

一般情况下都是对于一台物理机上的几个 vSwithc 上的 VM 进行同行进行的配置，比如两个 VM 各连接一个 vSwitch，这时候可以对物理机使用 ip link add veth1.1 type veth peer name veth1.2 命令俩创建一对虚拟接口然后使用 ovs-vsctl add-port BRIDGE PORT 命令分别把这两个虚拟接口绑定在两个 vSwitch 上，实现俩个 vSwitch 之间互联并且能够通信。还有就是如图所示，由于 OVS 有 DB，各 NODE 之间的 OVS 数据都互相共享，那么可以直接把 VM 连接到 vSwitch 上，然后再连接到物理网络就可以互通了相当于只是几个交换机互联而已，如果进行隔离后，使得隔离的 VM 可以通信，那么使用 namespace 功能创建一个 vRouter，通过 vRouter 实现被隔离的网络间互相通信

安装 OVS：yum install openvswitch -y

### OVS 的几种 INTERFACE 类型

1. patch：用于两个 vSwitch 的互联，需要一对，在每一个上面的选项配置对端接口

### OVS 命令行工具

1. ovs-vsctl \[OPTIONS] COMMAND \[ARGS(arguments 参数).....] #Open vSwitch-vSwitch control，OVS 的虚拟交换机控制

2. EXAMPLE

3. ovs-vsctl show #显示 ovsdb 中的内容概况，即显示 ovs 创建的相关网络信息

创建命令

1. ovs-vsctl add-br BRIDGE #创建一个名为 BRIDGE 的桥设备(即创建一个 vSwitch)，创建 BRIDGE 完成后会自动创建一个名字一样的 PORT 和 INTERFACE

2. ovs-vsctl add-port BRIDGE PORT #给 BRIDGE 这个桥设备添加一个 PORT(PORT 也可以是物理机的网卡，相当于把物理机的网卡连了根虚拟的网线到这个 BRIDGE 的虚拟交换机上)，创建完 PORT 后会自动创建一个名字一样的 INTERFACE

删除命令

1. ovs-vsctl del-br BRIDGE #删除一个名为 BRIDGE 的桥设备

2. ovs-vsctl del-port \[BRIDGE] PORT #从指定 BRIDGE 删除指定的 PORT，由于一个 PORT 只能绑定在一个 BRIDGE 上，所有 BRIDGE 可省略

查询命令

1. ovs-vsctl list-br #显示所有创建了的桥设备的名字(仅显示名字)

2. ovs-vsctl list-ports BRIDGE #显示 BRIDGE 这个桥设备上添加的所有 PORT 接口

3. ovs-vsctl list

DB 命令，设置、更改等，由于有 ovs 信息全部写在自己的 DB 中，所以这个命令就是以数据库的模式按列显示信息

1. ovs-vsctl list br|port|interface \[NAME] #显示 BRIDGE 或者 PORT 或者 INTERFACE 的详细信息\[具体某个的]

2. ovs-vsctl find #查找

3. ovs-vsctl set port PORT tag=NUM #设定 PORT 的 vlan tag 号为 NUM

4. ovs-vsctl set port PORT vlan_mode=trunk|access #设定该 PORT 的 vlan 模式为 trunk 还是 access

5. ovs-vsctl set port PORT trunks=NUM,NUM.... #设定该 PORT 在 trunk 模式下允许通过的 vlan NUM 号有哪些

6. ovs-vsctl set interface INTERFACE type=TYPE options:OPT=VAL #设定 INTERFACE 的类型为 TYPE，某个选项=该选项值.TYPE 默认是 internal

   1. ovs-vsctl set interface gre0 type=gre options:remote_ip=192.268.20.2 #设定 gre0 接口的类型为 gre，gre 类型中有个选项是远端地址，IP 为 192.168.20.2(gre 为隧道技术，所以设定该接口的时候必须要指明对端 IP，否则无法进行 IP 上套 IP 的封装操作)

   2. ovs-vsctl set interface vx0 type=vxlan options:remote_ip=192.268.20.2 #设定 gre0 接口的类型为 vxlan,vxlan 类型中有个选项是远端地址，IP 为 192.168.20.2
