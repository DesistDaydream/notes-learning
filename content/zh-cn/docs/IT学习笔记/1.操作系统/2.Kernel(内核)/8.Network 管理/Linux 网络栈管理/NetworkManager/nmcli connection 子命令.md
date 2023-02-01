---
title: nmcli connection 子命令
---

# 概述

> 参考：
> - [红帽官方文档,RedHat7-网络指南-使用 nmcli 创建带有 VLAN 的 bond 并作为 Bridge 的从设备](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/networking_guide/sec-vlan_on_bond_and_bridge_using_the_networkmanager_command_line_tool_nmcli)

**nmcli connection {show | up | down | modify | add | edit | clone | delete | monitor | reload | load | import | export} \[ARGUMENTS...]**

# up | down # 启动 | 停止连接

nmcli connection up \[\[id | uuid | path] <ID>] \[ifname <ifname>] \[ap <BSSID>] \[passwd-file <file with passwords>] # 启动连接

nmcli connection down \[id | uuid | path | apath] <ID> ... # 停止连接

EXAMPLE

- 启动名为 eth0 的 connection。i.e.把配置应用到指定的网络设备上，并且会自动重启网络设备
  - **nmcli con up eth0**

# add # 增加连接

使用指定的 PROPERTY(属性) 创建一个新的连接

Note：

1. add 与 modify 命令的参数用法基本一致。delete 连接没有多少复杂的参数，直接指定连接的标识符，即可将连接删除
2. 如果想对该连接进行更详细的配置，比如配置 ip、网关、bond 参数等等。就需要指定具体的 PROPETY 和对应的 VALUE。

## Syntax(语法)

**nmcli connection add \[save BOOLEAN] SETTING.PROPERTY VALUE ...**

- **save** # 指定该连接创建完成后，是否以文件形式保存到本地磁盘。默认为 true。
- **SETTING.PROPERTY(设置.属性)** # 该连接包含的属性。SETTING.PROPERTY 简称 property(属性) 用来指定要增加的连接的配置信息。
  - 如果 SETTING 和 PROPERTY 是 唯一的，则可以使用缩写(比如 connection.type 缩写为 type)。不同的 SETTING 中有不同的 PROPERTY。并非所有属性都适用于所有类型的连接(type 是创建连接时必须指定的一个属性)。也就是说 property(属性) 分为两种，一种适用于全局的通用属性，另一种是只对特定类型的连接生效。比如我当前创建一个 ethernet 类型的连接，那么就不能使用 bond 属性。
  - 注意：如果要在脚本中使用 nmcli 命令，最好不要使用别名
  - 可用的 SETTING.PROPERTY 详见 [Connection 配置详解](https://www.yuque.com/go/doc/33221861)，在命令行中的写法与配置文件中是一致的，SETTING 就是 配置文件中的 SETTING，PROPERTY 就是配置文件中所属 SETTING 的 PROPERTY。
  - 特殊连接类型的属性
    - 在命令使用中，特殊类型的属性可能并不在文档中，比如 bond 类型中，可以使用 mode 属性指定 bond 类型。
    - bond 类型的连接可以使用 mode 属性，该属性的值会添加到 bond.options 属性的值中，作为 mode=MODE 这种键值对存在。
- **VALUE** # SETTING.PROPERTY 的值

## EXAMPLE

- 创建名为 eth0 的 connection 并关联到 eth0 上，手动指定 IP，并设置开机自动启动网络,关闭 ipv6 网络。
  - nmcli con add type ethernet con-name eth0 ifname eth0 ipv4.method manual ip4 10.10.10.10/24 gw4 10.10.10.1 autoconnect yes ipv6.method disabled

# delete 删除连接

## Syntax(语法)

**nmcli connection delete \[id | uuid | path] <ID> # 删除连接**

# modify 修改连接

## Syntax(语法)

**nmcli connection modify \[+|-]PROPETY VALUE # 修改连接。**

Note：

- 在使用 modify 变更网络配置时，可以使用+或者-来实现"增加"或者"删除"某项配置的功能

## EXAMPLE

- 给名为 bond0 的 Connection 添加两个参数，使用加号可以不负载之前配置的参数而添加新的参数
  - nmcli con modify bond0 +bond.options miimon=200 +bond.options xmit_hash_policy=layer3+4

# show 查看连接

## Syntax(语法)

**nmcli connection show \[--active] \[--order <order spec>] \[id | uuid | path | apath] <ID> ...**

## EXAMPLE

- nmcli con show eth0 # 查看 eth1 这个 connectin 的所有状态，该命令会列出该 connection 的全部属性
- nmcli -f bond con show bridge-slave-bond0 # 查看 bridge-slave-bond0 这个连接配置中，bond 这个 SETTING 的所有属性及其值。效果如下


    [root@lichenhao ~]# nmcli -f bond con show bridge-slave-bond0
    bond.options:                           mode=active-backup

# 应用示例

## 静态路由配置

- 在 eth0 网卡上添加静态路由，目的网段是 192.168.122.0/24 的流量下一跳是 10.10.10.1
  - nmcli connection modify eth0 +ipv4.routes "192.168.122.0/24 10.10.10.1"
- 删除 eth0 网卡上的静态路由
  - nmcli connection modify eth0 -ipv4.routes "192.168.122.0/24 10.10.10.1"

## 路由策略配置

- 添加一条路由策略：优先级为 5，从 10.0.0.0/24 网段来的数据包，都通过 5000 路由表处理
  - nmcli connection add type ethernet con-name eth0 ifname eth0 ipv4.routing-rules "priority 5 from 10.0.0.0/24 table 5000"

## 多路由表，双网卡，双网关

- 添加 ens9 连接，配置 IP 地址，不分配默认网关，而是在 3 号路由表中配置一条路由条目：任意目的地址的下一跳是 192.168.122.1
  - nmcli con add type ethernet con-name ens9 ifname ens9 ipv4.method manual ipv4.addresses 192.168.122.2 ipv4.routes "0.0.0.0/0 192.168.122.1 table=3"
- 为 ens9 连接添加一个路由策略，将源地址是 192.168.122.0/24 网段的数据包，都交给 3 号路由表处理
  - nmcli con mod ens9 ipv4.routing-rules "priority 5 from 192.168.122.0/24 table 3"

## Bond 配置

- 添加一个 Bond 类型的连接
  - 使用 bond0 网络设备，bond 模式为主备
    - nmcli con add type bond con-name bond0 ifname bond0 mode active-backup
  - 使用 bond1 网络设备，bond 模式为 802.3ad
    - nmcli con add type bond con-name bond1 ifname bond1 mode 802.3ad
  - 创建一个 bond 类型的 connection，名字叫 bond0 且与 bond0 网络设备绑定(若没有 bond0 网络设备则自动创建)；手动设定 ip 并指定 ip、prefix；指定该 bond 的 3 个参数（bond 模式、检测时间、hash 算法）
    - nmcli con add type bond con-name bond0 ifname bond0 ipv4.method manual ipv4.addr 192.168.20.22/24 bond.options "mode=802.3ad,miimon=100,xmit_hash_policy=layer3+4"
- 添加 eth0 网络设备到 bond0 中
  - nmcli con add type ethernet master bond0 ifname eth0

## Bridge 配置

- 创建一个 Bridge 类型的 connection，名字叫 br0 且与 br0 网络设备绑定(若没有 br0 网络设备则自动创建)，手动获取 ip 并设定 ip、prefix、gateway。
  - nmcli con add type bridge con-name br0 ifname br0 ipv4.method manual ip4 192.168.10.10/24 gw4 192.168.10.1
- 添加 eth0 网络设备到 br0 中
  - nmcli con add type ethernet ifname eth0 master br0

注意：若不为 bridge 类型的网络设备配置 IP，比如在虚拟化环境中，需要关闭 STP 功能

## Vlan 配置

创建一个 Blan 类型的连接，意思就是为指定的网络设备 DEV 划分 vlan。创建一个名为 DEV.VLANID 的新网络设备，凡是通过该设备发送的数据包都会添加上 VLANID。

为 Bond 配置 VLAN 标签

- 创建 bond 类型的连接，名为 bond1，关闭 IPv4 和 IPv6
  - nmcli con add type bond con-name bond1 ifname bond1 ipv4.method disabled ipv6.method ignore bond.options "mode=802.3ad,miimon=100,xmit_hash_policy=layer3+4"
- 创建 vlan 类型的连接，绑定到 bond1 设备上，vlan 号为 2409，配置 IP
  - nmcli con add type vlan con-name vlan2409-bond1 dev bond1 id 2409 ipv4.method manual ipv4.addresses 100.75.9.17/24 ipv4.gateway 100.75.9.254

为 Ethernet 配置 VLAN 标签

- 创建 ethernet 类型的连接，名为 eth0，关闭 IPv4 和 IPv6
  - nmcli con add ethernet con-name eth0 ifname eth0 ipv4.method disabled ipv6.method ignore
- 创建一个 vlan 类型的连接，连接名为 vlan2-bond0，为 eth0 划分 vlan，vlan 号为 2
  - nmcli con add type vlan con-name vlan2-eth0 dev eth0 id 2 ipv4.method manual ipv4.addresses 100.75.9.17/24

## 具有 VLAN TAG 的 Ethernet 绑定到 Bridge 上

```bash
nmcli connection add type ethernet con-name ens1f0 ifname ens1f0 ipv4.method disabled ipv6.method ignore
nmcli connection add type bridge con-name br0 ifname br0 ipv4.method manual ipv4.address 10.253.26.242/24 ipv4.gateway 10.253.26.254
nmcli connection add type vlan con-name ens1f0.1251 ifname ens1f0.1251 dev ens1f0 id 1251 master br0 slave-type bridge
```

这种配置方法待验证，使用红帽官方的方法配完了起不来，bridge 设备是 down 的状态，ens1f0.1251 设备是 lowerlayerdown 状态。主要问题出在 bridge 设备的配置上，通过简单的传统配置，只要不由 NetworkManager 管理 ifcfg-br0，即可正常使用。

> 在寻找该问题的解决方法时，发现了一个相关 BUG，详见：<https://serverfault.com/questions/682183/bridge-on-vlan-on-teaming-for-kvm/861450>。不过这个连接的解决方案并不适合我

最主要的还是 Bridge 设备的配置，**在创建 Bridge 设备的时候，关闭 STP 即可解决**。可能与交换机那边的设置有关。

```bash
nmcli connection add type ethernet con-name ens1f0 ifname ens1f0 ipv4.method disabled ipv6.method ignore
nmcli connection add type bridge con-name br0 ifname br0 ipv4.method manual ipv4.address 10.253.26.242/24 ipv4.gateway 10.253.26.254 bridge.stp no
nmcli connection add type vlan con-name ens1f0.1251 ifname ens1f0.1251 dev ens1f0 id 1251 master br0
```

## 具有 VLAN TAG 的 Bond 并绑定到 Bridge 上

若出现问题无法启动，解决方法与上面的《配置 VLAN TAG 到 Ethernet 绑定到 Bridge 上》的例子解决方法一样。

```bash
nmcli connection add type bond con-name Bond0 ifname bond0 bond.options "mode=active-backup,miimon=100" ipv4.method disabled ipv6.method ignore
nmcli connection add type ethernet con-name Slave1 ifname ens1f0 master bond0 slave-type bond
nmcli connection add type ethernet con-name Slave2 ifname ens1f1 master bond0 slave-type bond
nmcli connection add type bridge con-name Bridge0 ifname br0 ip4 10.253.26.242/24
nmcli connection add type vlan con-name Vlan1251 ifname bond0.1251 dev bond0 id 1251 master br0 slave-type bridge

nmcli connection show
 NAME     UUID                                  TYPE            DEVICE
 Bond0    f05806fa-72c3-4803-8743-2377f0c10bed  bond            bond0
 Bridge0  22d3c0de-d79a-4779-80eb-10718c2bed61  bridge          br0
 Slave1   e59e13cb-d749-4df2-aee6-de3bfaec698c  802-3-ethernet  em1
 Slave2   25361a76-6b3c-4ae5-9073-005be5ab8b1c  802-3-ethernet  em2
 Vlan2    e2333426-eea4-4f5d-a589-336f032ec822  vlan            bond0.2
```
