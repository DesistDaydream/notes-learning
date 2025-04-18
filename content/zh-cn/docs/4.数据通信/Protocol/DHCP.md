---
title: DHCP
linkTitle: DHCP
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, DHCP](https://en.wikipedia.org/wiki/Dynamic_Host_Configuration_Protocol)

**Dynamic Host Configuration Protocol(动态主机设置协议，简称 DHCP)** 主要用于由一部主机来自动的分配所有的网络参数给指定网段内的所有设备。

DHCP 可以分配的网络参数有 IP 地址、掩码、网关、DNS 等

DHCP 的运作方式：
他主要由客户端发送广播包给整个物理网段内的所有主机， 若该网段内有 DHCP 服务器时，就会响应客户端的 DHCP 请求。所以，DHCP 服务器与客户端是应该要在同一个物理网段内的，如果想跨网段提供 DHCP 服务，需要在对应网段启用 dhcrelay 服务。

至于整个 DHCP 封包在服务器与客户端的来来回回情况如右图，具体有 4 个步骤

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/sc7lh2/1616161468691-c82332ec-254e-485f-861d-9186e60a3dd8.jpeg)

- **DISCOVER(发现)** # 客户端利用广播封包发送搜索 DHCP 服务器的封包：
  - 若客户端网络设定使用 DHCP 协议取得 IP (在 Windows 内为『自动取得 IP』)，则当客户端开机或者是重新启动网络卡时， 客户端主机会发送出搜寻 DHCP 服务器的 UDP 封包给所有物理网段内的计算机。此封包的目标 IP 会是 255.255.255.255， 所以一般主机接收到这个封包后会直接予以丢弃，但若局域网络内有 DHCP 服务器时，则会开始进行后续行为。
- **OFFER(提供)** # DHCP 服务端提供客户端网络相关的租约以供选择
  - DHCP 服务器在接收到这个客户端的要求后，会针对这个客户端的硬件地址 (MAC) 与本身的设定数据来进行下列工作：
    - 到服务器的登录文件中寻找该用户之前是否曾经用过某个 IP ，若有且该 IP 目前无人使用，则提供此 IP 给客户端；
    - 若配置文件针对该 MAC 提供额外的固定 IP (static IP) 时，则提供该固定 IP 给客户端；
    - 若不符合上述两个条件，则随机取用目前没有被使用的 IP 参数给客户端，并记录下来。
    - 总之，服务器端会针对客户端的要求提供一组网络参数租约给客户端选择，由于此时客户端尚未有 IP ，因此服务器端响应的封包信息中， 主要是针对客户端的 MAC 来给予回应。此时服务器端会保留这个租约然后开始等待客户端的回应。
- **REQUEST(请求)** # 客户端决定选择的 DHCP 服务器提供的网络参数租约并回报服务器
  - 由于局域网络内可能并非仅有一部 DHCP 服务器，但客户端仅能接受一组网络参数的租约。 因此客户端必需要选择是否要认可该服务器提供的相关网络参数的租约。当决定好使用此服务器的网络参数租约后， 客户端便开始使用这组网络参数来设定自己的网络环境。此外，客户端也会发送一个广播封包给所有物理网段内的主机， 告知已经接受该服务器的租约。此时若有第二台以上的 DHCP 服务器，则这些没有被接受的服务器会收回该 IP 租约。至于被接受的 DHCP 服务器会继续进行底下的动作。
- **Acknowledge(确认，ACK)** # 服务端记录该次租约行为并回报客户端已确认的响应封包信息
  \- 当服务器端收到客户端的确认选择后，服务器会回传确认的响应封包，并且告知客户端这个网络参数租约的期限， 并且开始租约计时喔！那么该次租约何时会到期而被解约，可以这样想：
  \- 客户端脱机：不论是关闭网络接口 (ifdown)、重新启动 (reboot)、关机 (shutdown) 等行为，皆算是脱机状态，这个时候 Server 端就会将该 IP 回收，并放到 Server 自己的备用区中，等待未来的使用；
  \- 客户端租约到期：前面提到 DHCP server 端发放的 IP 有使用的期限，客户端使用这个 IP 到达期限规定的时间，而且没有重新提出 DHCP 的申请时，就需要将 IP 缴回去！这个时候就会造成断线。但用户也可以再向 DHCP 服务器要求再次分配 IP 啰。

这四个步骤也称为 DHCP 分配地址时的需要进行的 DORA 进程

## DHCP 租约

每次 DHCP 服务给客户端提供网络参数时，同时会提供一个有效时间，该有效时间表示客户端使用这些参数的有效时间，当有效时间过了之后，客户端则不再拥有这些网络参数。租约是由客户端发起请求具体租用时间，然后服务器回应是否可以给客户端租用这些时间。

DHCP 提供的租约信息保存在 /var/lib/dhcpd/dhcpd.lesses 文件中，租约中包括这些参数：提供的 IP，客户端 MAC，租赁时长

# ISC-DHCP

> 参考：
>
> - [ISC DHCP 官网](https://www.isc.org/dhcp/)
> - <https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/8/html/configuring_and_managing_networking/providing-dhcp-services_configuring-and-managing-networking#setting-up-the-dhcp-service-for-subnets-directly-connected-to-the-dhcp-server_providing-dhcp-services>

ISC-DHCP 提供了一个完整的开源解决方案，用来实现 DHCP 服务端、DHCP 客户端、DHCP 中继代理 这几种 DHCP 服务，可以将 ISC-DHCP 看作一个**程序的集合**。ISC DHCP 支持 IPv4 和 IPv6，适用于大批量和高可靠性的应用

CentOS 和 Ubuntu 安装的都是 ISC 分发的 DCHP 程序，不过 ISC 官方已经推荐使用 Kea DHCP 替代 ISC DHCP

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/sc7lh2/1643021582679-beb6a411-3391-4577-a658-295acbd0f2a5.png)

注意：

- 在为 IPv6 分配地址时，通常需要与 radvd 程序一起运行，详见 [DHCPv6 与 radvd 章节](#NMGBa)

## 安装

### CentOS

```bash
yum install dhcp
```

CentOS 的包是真的不更新。。。囧。。。。好多这种情况了。。。包中写的网址都是错误的。。。

```bash
~]# yum info dhcp
Loaded plugins: fastestmirror, product-id, search-disabled-repos, subscription-manager

This system is not registered with an entitlement server. You can use subscription-manager to register.

Loading mirror speeds from cached hostfile
 * base: mirrors.tuna.tsinghua.edu.cn
 * epel: mirrors.tuna.tsinghua.edu.cn
 * extras: mirrors.tuna.tsinghua.edu.cn
 * updates: mirrors.tuna.tsinghua.edu.cn
Available Packages
Name        : dhcp
Arch        : x86_64
Epoch       : 12
Version     : 4.2.5
Release     : 83.el7.centos.1
Size        : 515 k
Repo        : updates/7/x86_64
Summary     : Dynamic host configuration protocol software
URL         : http://isc.org/products/DHCP/   // 在 ISC 官网。。。这个页面已经没了。。。。。
License     : ISC
Description : DHCP (Dynamic Host Configuration Protocol) is a protocol which allows
            : individual devices on an IP network to get their own network
            : configuration information (IP address, subnetmask, broadcast address,
            : etc.) from a DHCP server. The overall purpose of DHCP is to make it
            : easier to administer a large network.
            :
            : To use DHCP on your network, install a DHCP service (or relay agent),
            : and on clients run a DHCP client daemon.  The dhcp package provides
            : the ISC DHCP service and relay agent.
```

### Ubunt

```bash
apt install isc-dhcp-server
```

包的信息

```bash
~]# apt-cache show isc-dhcp-server
Package: isc-dhcp-server
Architecture: amd64
Version: 4.4.1-2.1ubuntu5.20.04.2
Priority: optional
Section: net
Source: isc-dhcp
Origin: Ubuntu
Maintainer: Ubuntu Developers <ubuntu-devel-discuss@lists.ubuntu.com>
Original-Maintainer: Debian ISC DHCP Maintainers <isc-dhcp@packages.debian.org>
Bugs: https://bugs.launchpad.net/ubuntu/+filebug
Installed-Size: 1503
Depends: debconf (>= 0.5) | debconf-2.0, libc6 (>= 2.15), libdns-export1109, libirs-export161, libisc-export1105, debianutils (>= 2.8.2), lsb-base, adduser
Recommends: isc-dhcp-common
Suggests: policykit-1, isc-dhcp-server-ldap, policycoreutils
Breaks: isc-dhcp-common (= 4.3.3-1), logcheck-database (= 1.3.17~)
Replaces: isc-dhcp-common (= 4.3.3-1)
Filename: pool/main/i/isc-dhcp/isc-dhcp-server_4.4.1-2.1ubuntu5.20.04.2_amd64.deb
Size: 454712
MD5sum: 92a2ec90073e62f5fe65ecb840e2fcd9
SHA1: 1f1253f7bcd4a8d3bec5d9736d4238115635df21
SHA256: b537b40e5c35054d8d3f82060936de737b08f1daa48e731652aa230170f0b21a
SHA512: f00f30b52b085dcf9737c7cd5d37095ecdbcf9750e7bf1e5ca1f80591d1ad9f24d5b29b2c89d40dee639803b7e1e90b92f8b3396d8976a1171321eef5d960559
Homepage: http://www.isc.org
Description-en: ISC DHCP server for automatic IP address assignment
 This is the Internet Software Consortium's DHCP server.
 .
 Dynamic Host Configuration Protocol (DHCP) is a protocol like BOOTP
 (actually dhcpd includes much of the functionality of bootpd). It
 gives client machines "leases" for IP addresses and can
 automatically set their network configuration.
 .
 This server can handle multiple ethernet interfaces.
Description-md5: 38647f497f13c9a0a99f9d9cf772d70d
```

## dhcp 程序关联文件

/etc/default/isc-dhcp-server #

/etc/dhcp/dhcpd.conf # dhcpd 程序运行时行为的配置文件

/var/lib/dhcpd/dhcpd.leases 与 /var/lib/dhcpd/dhcpd.leases~ # DHCP 的客户端租约的缓存文件

- 所有从 DHCP 服务器获取 IP 的客户端信息都会记录在这个文件里
- 并且 dhcp 服务会探测所能分配的 IP 是否被占用，所有被占用的记录也会记在这个文件里
- 注意：当 DHCP 无法按照预想的样子分配 IP 的时候，可以尝试清除该文件里的内容，让其重新获取信息

### 配置文件样例

dhcp4

```nginx
ddns-update-style interim;

allow booting;
allow bootp;

ignore client-updates;
set vendorclass = option vendor-class-identifier;

option pxe-system-type code 93 = unsigned integer 16;

#以下是dhcp服务所能租赁的具体网段的IP信息
subnet 192.168.10.0 netmask 255.255.255.0 {
    # 指定要分配的网关
     option routers             192.168.10.1;
     # 指定要分配的DNS地址
     option domain-name-servers 114.114.114.114;
     # 指定要分配的子网掩码
     option subnet-mask         255.255.255.0;
     # 指定可分配的IP地址范围是从哪到哪
     range dynamic-bootp        192.168.10.100 192.168.10.254;
     # 默认租赁时间，如果客户端没有请求租约，则提供默认时间，数值单位为秒
     default-lease-time         21600;
     # 最大租赁时间，可以提供给客户端租用网络参数的最大时间，数值单位为秒
     max-lease-time             43200;
     # PXE环境下指定的提供引导程序的服务器。i.e.该参数的值为PEX服务端的设备IP或HostName。
     # 比如cobbler服务所在的服务器地址，或者tftp服务所在服务地址等等
     next-server                10.10.17.15;
     class "pxeclients" {
          match if substring (option vendor-class-identifier, 0, 9) = "PXEClient";
          if option pxe-system-type = 00:02 {
                  filename "ia64/elilo.efi";
          } else if option pxe-system-type = 00:06 {
                  filename "grub/grub-x86.efi";
          } else if option pxe-system-type = 00:07 {
                  filename "grub/grub-x86_64.efi";
          } else if option pxe-system-type = 00:09 {
                  filename "grub/grub-x86_64.efi";
          } else {
                  filename "pxelinux.0";
          }
     }

}

# group for Cobbler DHCP tag: default
group {
}

#当有其余网段需要分配时，可以继续添加相关信息
subnet 10.10.17.0 netmask 255.255.255.0 {
 option domain-name-servers 114.114.114.114;
 option routers 10.10.17.1;
 range dynamic-bootp 10.10.17.15 10.10.17.22;
 option subnet-mask 255.255.255.0;
 next-server 10.10.17.15;
 default-lease-time 21600;
 max-lease-time 43200;
}
```

dhcp6

```properties
default-lease-time 600;
max-lease-time 7200;
log-facility local7;
subnet6 2001:db8:0:1::/64 {
    # 可以分配的地址范围
    range6 2001:db8:0:1::129 2001:db8:0:1::254;

    # 为客户端分配的临时地址
    range6 2001:db8:0:1::/64 temporary;

    # Prefix range for delegation to sub-routers
    prefix6 2001:db8:0:100:: 2001:db8:0:f00:: /56;

    # 一些可选的配置
    option dhcp6.name-servers fec0:0:0:1::1;
    option dhcp6.domain-search "domain.example";

    # 将 IP 地址与 MAC 绑定的示例
    host specialclient {
        host-identifier option dhcp6.client-id 00:01:00:01:4a:1f:ba:e3:60:b9:1f:01:23:45;
        fixed-address6 2001:db8:0:1::127;
    }
}
```

# Kea-DHCP

> 参考：
>
> - [Kea DHCP 官网](https://www.isc.org/kea/)

# DhcRelay(DHCP 中继器)简介

DHCP Relay，DHCP 中继器，用于从没有 DHCP 服务的子网中连接到其他子网上的一个或多个 DHCP 服务器，来中继(代理、转发)DHCP 和 BOOTP 请求。它支持 DHCPv4 / BOOTP 和 DHCPv6 协议。实际应用：e.g.我有两个网段 192.168.10.0/24 和 192.168.20.0/24，在 10 网段有 DHCP 服务器，但是 20 网段没有，这时候就可以在 20 网段的设备上开启 dhcrelay 服务，然后指定一个 DHCP 服务器，来为 20 网段进行 DHCP 服务。

## dhcrelay 的使用方法

## dhcrelay 命令

语法格式：

- dhcrelay \[ -4 ] \[ -dqaD ] \[ -p port ] \[ -c count ] \[ -A length ] \[ -pf pid-file ] \[ --no-pid ] \[ -m append | replace | forward | dis‐card ] \[ -i interface0 \[ ... -i interfaceN ] ] server0 \[ ...serverN ]
- dhcrelay -6 \[ -dqI ] \[ -p port ] \[ -c count ] \[ -pf pid-file ] \[ --no-pid ] -l lower0 \[ ... -l lowerN ] -u upper0 \[ ... -u upperN ]

OPTIONS

- -i # 在指定网络接口上监听 DHCPv4 / BOOTP 查询。

EXAMPLE

- dhcrelay 192.168.10.12 # 在本机开启 dhcp 中继代理，指定目的 DHCP 服务器 IP 为 192.168.10.12

# DHCPv6 与 radvd

> 参考：
>
> - <https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/8/html/configuring_and_managing_networking/providing-dhcp-services_configuring-and-managing-networking#comparison-of-dhcpv6-to-radvd_providing-dhcp-services>

在 IPv6 网络中，只有路由器广告信息在 IPv6 默认网关上提供信息。因此，如果您要在需要默认网关设置的子网中使用 DHCPv6，还必须配置路由器广告服务，如 Router Advertisement Daemon（radvd）。

radvd 服务使用路由器广告数据包中的标记声明 DHCPv6 服务器可用。

路由器广告守护进程（radvd）发送路由器公告信息，这是 IPv6 无状态自动配置所需的。这可让用户根据这些公告自动配置其地址、设置、路由和选择默认路由器。
本节中的步骤解释了如何配置 radvd。
**先决条件**

- 您以 root 用户身份登录。

**流程**

1. 安装 radvd 软件包：# yum install radvd
2. 编辑 /etc/radvd.conf 文件并添加以下配置：

```properties
interface ens3
{
  AdvSendAdvert on;
  AdvManagedFlag on;
  AdvOtherConfigFlag on;

  prefix 2001:db8:0:1::/64 {
  };
};
```

3. 这些设置将 radvd 配置为在 enp1s0 设备中为 2001:db8:0:1::/64 子网发送路由器广告信息。AdvManagedFlag on 设置定义客户端应该从 DHCP 服务器接收 IP 地址，AdvOtherConfigFlag 参数设置为 on 定义客户端也应该从 DHCP 服务器接收非地址信息。
4. （可选）配置 radvd 会在系统引导时自动启动：# systemctl enable radvd
5. 启动 radvd 服务：# systemctl start radvd
6. 另外，还可显示路由器公告软件包的内容和配置的值 radvd 发送：# radvdump

**其它资源**

- 有关配置 radvd 的详情，请查看 radvd.conf(5) man page。
- 如需 radvd 的示例配置，请参阅 /usr/share/doc/radvd/radvd.conf.example 文件。
