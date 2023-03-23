---
title: firewalld(Iptables 的管理工具)
---

# 概述

firewalld 与 iptabels 一样，是管理 Linux 内核中的 Netfilter 功能的工具。

FirewallD 使用两个配置模式：“runtime 运行时”和“permanent 持久”。

1. runtime 模式：默认模式。所有配置即时生效，在重启系统、重新启动 FirewallD 时、使用--reload 重载配置等操作是，在该模式下的配置都将被清除。
2. permanent 模式：需要使用 --permanent 选项生效，配置才会永久保存。如果想让 permanetn 模式下的配置立即生效，需要使用--reload 命令或者重启 firewalld 服务。

## firewalld 中 zone(区域)的概念

“区域”是针对给定位置或场景（例如家庭、公共、受信任等）可能具有的各种信任级别的预构建规则集。不同的区域允许不同的网络服务和入站流量类型，而拒绝其他任何流量。 首次启用 FirewallD 后，public 将是默认区域。

区域也可以用于不同的网络接口。例如，要分离内部网络和互联网的接口，你可以在 internal 区域上允许 DHCP，但在 external 区域仅允许 HTTP 和 SSH。未明确设置为特定区域的任何接口将添加到默认区域。

所以，一般情况下，所有区域都是拒绝所有新的入站流量，对已经建立连接的不再阻止。在使用 firewall-cmd 命令添加某 service、port、ip 等属性时，相当于允许对应属性的流量入站。i.e.添加进去就表示允许

## zone 的种类与说明

1. public（公共） # 默认的 zone。在公共区域内使用，不能相信网络内的其他计算机不会对您的计算机造成危害，只能接收经过选取的连接。
2. block（限制） # 任何接收的网络连接都被 IPv4 的 icmp-host-prohibited 信息和 IPv6 的 icmp6-adm-prohibited 信息所拒绝。
3. dmz（非军事区） # 用于您的非军事区内的电脑，此区域内可公开访问，可以有限地进入您的内部网络，仅仅接收经过选择的连接。
4. drop（丢弃） # 任何接收的网络数据包都被丢弃，没有任何回复。仅能有发送出去的网络连接。
5. external（外部） # 特别是为路由器启用了伪装功能的外部网。您不能信任来自网络的其他计算，不能相信它们不会对您的计算机造成危害，只能接收经过选择的连接。
6. home（家庭） # 用于家庭网络。您可以基本信任网络内的其他计算机不会危害您的计算机。仅仅接收经过选择的连接。
7. internal（内部） # 用于内部网络。您可以基本上信任网络内的其他计算机不会威胁您的计算机。仅仅接受经过选择的连接。
8. trusted（信任） # 可接受所有的网络连接。
9. work（工作） # 用于工作区。您可以基本相信网络内的其他电脑不会危害您的电脑。仅仅接收经过选择的连接。

用实际举例：将设备某个网卡放在区域中，则流经该网卡的流量都会遵循该区域中所定义的规则。

# Firewalld 关联文件与配置

/usr/lib/firewalld # 保存默认配置，如默认区域和公用服务。 避免修改它们，因为每次 firewall 软件包更新时都会覆盖这些文件。
/etc/firewalld # 保存系统配置文件。这些文件将覆盖默认配置。

- firewalld.conf #

# firewall 安装完成后的 iptables 模式的默认配置

以下是 public 区域下 filter 表的默认配置，大部分都是对于自定义链的规则

- # 设置 3 个基本链的默认 target

- -P INPUT ACCEPT
- -P FORWARD ACCEPT
- -P OUTPUT ACCEPT

- # 默认会创建多个自定义链

- -N FORWARD_IN_ZONES
- -N FORWARD_IN_ZONES_SOURCE
- -N FORWARD_OUT_ZONES
- -N FORWARD_OUT_ZONES_SOURCE
- -N FORWARD_direct
- -N FWDI_public
- -N FWDI_public_allow
- -N FWDI_public_deny
- -N FWDI_public_log
- -N FWDO_public
- -N FWDO_public_allow
- -N FWDO_public_deny
- -N FWDO_public_log
- -N INPUT_ZONES
- -N INPUT_ZONES_SOURCE
- -N INPUT_direct
- -N IN_public
- -N IN_public_allow
- -N IN_public_deny
- -N IN_public_log
- -N OUTPUT_direct

- # 设置 INPUT 链基本规则，所有流量直接交给 INPUT_direct、INPUT_ZONES_SOURCE、INPUT_ZONES 这 3 个自定义链来继续进行规则匹配

- -A INPUT -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
- -A INPUT -i lo -j ACCEPT
- -A INPUT -j INPUT_direct
- -A INPUT -j INPUT_ZONES_SOURCE
- -A INPUT -j INPUT_ZONES # 流量转给 INPUT 的 ZONES
- -A INPUT -m conntrack --ctstate INVALID -j DROP
- -A INPUT -j REJECT --reject-with icmp-host-prohibited # 在 INPUT 链上拒绝所有流量，并通过 icmp 协议提示客户端 prohibited

- # 设置 FORWARD 链基本规则，所有流量直接交给

- -A FORWARD -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
- -A FORWARD -i lo -j ACCEPT
- -A FORWARD -j FORWARD_direct
- -A FORWARD -j FORWARD_IN_ZONES_SOURCE
- -A FORWARD -j FORWARD_IN_ZONES
- -A FORWARD -j FORWARD_OUT_ZONES_SOURCE
- -A FORWARD -j FORWARD_OUT_ZONES
- -A FORWARD -m conntrack --ctstate INVALID -j DROP
- -A FORWARD -j REJECT --reject-with icmp-host-prohibited # 在 FORWARD 链上拒绝所有流量，并通过 icmp 协议提示客户端 prohibited

- # 设置 OUTPUT 链基本规则，直接把后续检查转交给 OUTPUT_direct 这个自定义 chain 进行规则匹配

- -A OUTPUT -j OUTPUT_direct

- # FORWARD 相关的自定义 chain 规则

- -A FORWARD_IN_ZONES -i bond0 -g FWDI_public
- -A FORWARD_IN_ZONES -i eth2 -g FWDI_public
- -A FORWARD_IN_ZONES -i eth1 -g FWDI_public
- -A FORWARD_IN_ZONES -i eth0 -g FWDI_public
- -A FORWARD_IN_ZONES -g FWDI_public
- -A FORWARD_OUT_ZONES -o bond0 -g FWDO_public
- -A FORWARD_OUT_ZONES -o eth2 -g FWDO_public
- -A FORWARD_OUT_ZONES -o eth1 -g FWDO_public
- -A FORWARD_OUT_ZONES -o eth0 -g FWDO_public
- -A FORWARD_OUT_ZONES -g FWDO_public
- -A FWDI_public -j FWDI_public_log
- -A FWDI_public -j FWDI_public_deny
- -A FWDI_public -j FWDI_public_allow
- -A FWDI_public -p icmp -j ACCEPT
- -A FWDO_public -j FWDO_public_log
- -A FWDO_public -j FWDO_public_deny
- -A FWDO_public -j FWDO_public_allow

- # INPUT 相关的自定义 chain 规则

- # INPUT_ZONES 用来将各个网络设备区分到指定的 ZONE 中

- -A INPUT_ZONES -i bond0 -g IN_public
- -A INPUT_ZONES -i eth2 -g IN_public
- -A INPUT_ZONES -i eth1 -g IN_public
- -A INPUT_ZONES -i eth0 -g IN_public # 将 eth0 的流量放到 public 区域中继续进行匹配
- -A INPUT_ZONES -g IN_public

- # INPUT 链上的 public 区域的规则

- -A IN_public -j IN_public_log
- -A IN_public -j IN_public_deny
- -A IN_public -j IN_public_allow
- -A IN_public -p icmp -j ACCEPT
- -A IN_public_allow -p tcp -m tcp --dport 22 -m conntrack --ctstate NEW -j ACCEPT

- # XXXX 区域的规则

- 。。。。每当一个网络设备被放到某个区域中，这个区域就会激活，会在整个 iptables 表中显示，可以使用 firewall-cmd --zone=drop --change-interface=eth1 进行验证

从下往上看的话，firewalld 会默认方通 22 端口(i.e.方通 sshd 服务)和 icmp 协议，并且自定义规则链的数据结构详见脑图 firewalld 之 filter 表基本配置图.mindmap
