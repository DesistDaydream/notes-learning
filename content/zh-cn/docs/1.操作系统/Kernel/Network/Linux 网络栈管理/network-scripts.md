---
title: network-scripts
---

# 概述

> 参考：
>
> - [RedHat 官方文档，生产文档-RedHatEnterpriseLinux-6-部署指南-11.2.接口配置文件](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/6/html/deployment_guide/s1-networkscripts-interfaces)
> - [/usr/share/doc/initscripts-XX/sysconfig.txt](https://github.com/OpenMandrivaSoftware/initscripts/blob/master/sysconfig.txt)中的 ifcfg-\<interface-name> 部分
> - [Manual(手册),nm-settings-ifcfg-rh(5)](https://networkmanager.dev/docs/api/latest/nm-settings-ifcfg-rh.html)

RedHad 相关发行版的网络配置通过一系列脚本实现，随着时代的发展，已经逐渐启用，并由 NetworkManager 取代，NetworkManager 还单独出了一个适用于 RedHad 发行版的插件，名为 nm-setting-ifcfg-rh。这样，NetworkManager 可以将原本的配置目录中文件的格式，转变为适应 RedHad 的格式，并将配置文件保存到 /etc/sysconfig/network-scripts/ 目录下。

# 关联文件

**/etc/sysconfig/\*** # 全局

- **./network** # 全局网络配置
- **./network-scripts/\*** # 曾经是网络配置脚本文件所在目录。CentOS 8 以后，移除了所有脚本，只用来为网络配置程序提供网络设备的配置文件
  - **./ifcfg-INTERFACE** # 名为 INTERFACE 网络设备配置文件。通常情况下，INTERFACE 的值通常与配置文件中 DEVICE 指令的值相同。
  - **./route-INTERFACE** # IPv4 静态路由配置文件。INTERFACE 为网络设备名称，该路由条目仅对名为 INTERFACE 的网络设备起作用
  - **./route6-INTERFACE** # IPv6 静态路由配置文件。INTERFACE 为网络设备名称，该路由条目仅对名为 INTERFACE 的网络设备起作用
  - **./rule-INTERFACE** # 定义内核将流量路由到特定路由表的 IPv4 源网络规则。
  - **./rule6-INTERFACE** # 定义内核将流量路由到特定路由表的 IPv6 源网络规则。
- **./networking/\** * #
  - 注意：在 [RedHat 6 文档](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/6/html/deployment_guide/ch-network_interfaces#s1-networkscripts-files)中表示，/etc/sysconfig/networking/ 目录由现在已经弃用的网络管理工具(system-config-network) 管理，这个内容不应该手动编辑。推荐使用 NetworkManager。并且在后续的版本中， NetworkManager 也接管了这些文件

**/etc/iproute2/rt_tables** # 如果您想要使用名称而不是数字来引用特定的路由表，这个文件会定义映射映射。

## rule-INTERFACE 文件

```bash
from 192.0.2.0/24 lookup 1
from 203.0.113.0/24 lookup 2
```

来自 192.0.2.0/24 的流量根据 1 号路由表规则进行路由

来自 203.0.113.0/24 的流量根据 2 号路由表规则进行路由

## route-INTERFACE 文件

第一种格式：

```bash
# 默认路由下一条是 192.168.1.1，从 eth0 网络设备发出
default via 192.168.1.1 dev eth0
# 目的网段是 10.0.0.1 且掩码是 255.255.255.0，从 eth1 网络设备发出数据包，下一跳为 192.168.0.1
10.0.0.1 192.168.0.1 255.255.255.0 eth1
```

第二种格式：每一个路由用 0,1,2....等表示

```bash
ADDRESS0=10.10.10.0 #目的网段IP
NETMASK0=255.255.255.0 #目的网段掩码
GATEWAY0=192.168.0.10 #下一跳IP
ADDRESS1=172.16.1.10 #目的网段IP
NETMASK1=255.255.255.0 #目的网段掩码
GATEWAY1=192.168.0.10 #下一跳IP
```

## ifcfg-INTERFACE 文件

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ggdfnf/1616165813293-bc7fde9f-4810-4fc8-b6b8-67282aaef73d.jpeg)

ifcfg-INTERFACE 文件由多个 **Items(条目)** 组成，每个 Item 类似于 键值对，以 `=` 分割，Item 分为很多类型：

- Base items:
- Mandriva specific items for DHCP clients:
- Mandriva items:
- Base items being deprecated:
- Alias specific items:
- IPv6-only items for real interfaces:
- IPv6-only items for static tunnel interface:
- Ethernet-only items:
- Ethernet 802.1q VLAN items:
- PPP/SLIP items:
- PPP-specific items
- IPPP-specific items (ISDN)
- ippp0 items being deprecated:
- Wireless-specific items:
- IPSEC specific items
- Bonding-specific items
- Tunnel-specific items:
- Bridge-specific items:
- TUN/TAP-specific items:

### Base items(基本条目)

### Bonding-specific items(特定于 Bonding 的条目)

**SLAVE={yes|no}** # 指定设备是否为 slave 设备。`默认值：no`
**MASTER=bondXX** # 指定要绑定的主设备
**BONDING_OPTS=\<OPTS>** # Bonding 驱动运行时选项，多个选项以空格分割

- "mode=active-backup arp_interval=60 arp_ip_target=192.168.1.1,192.168.1.2"

#### 配置示例

```bash
DEVICE=bond0
IPADDR=192.168.1.1
NETMASK=255.255.255.0
ONBOOT=yes
BOOTPROTO=none
USERCTL=no
NM_CONTROLLED=no
BONDING_OPTS="bonding parameters separated by spaces"
```
