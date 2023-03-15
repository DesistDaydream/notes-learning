---
title: pptp(VPN)协议及其实现
---

# 概述

Note：pptpd 存在安全隐患，详情可参考[这里](http://pptpclient.sourceforge.net/protocol-security.phtml)

The PPTP Server for Linux

# Poptop 介绍

官方网址：<http://poptop.sourceforge.net/>

## Poptop 的安装

需要 epel 源

yum install pptpd

## Poptop 关联文件与配置

/etc/pptpd.conf #pptpd 主程序的主配置文件

localip IP #当远程客户端连接到本地服务时，服务端建立虚拟网卡所用的 ip

remoteip IP-Range #远程客户端连接到本地服务后，可以分配的给客户端的 ip 范围

Note：这俩配置完成后，当客户端连接时，就会出现下面虚拟网卡。

    4: ppp0: <POINTOPOINT,MULTICAST,NOARP,UP,LOWER_UP> mtu 1396 qdisc pfifo_fast state UNKNOWN group default qlen 3
        link/ppp
        inet 192.168.0.1 peer 192.168.0.234/32 scope global ppp0
           valid_lft forever preferred_lft forever

/etc/ppp/options.pptpd #pptpd 主程序运行时参数

ms-dns IP #配置 DNS

/etc/ppp/chap-secrets #配置客户端连接时所用的用户名、密码、协议

    # Secrets for authentication using CHAP
    # client        server  secret       IP addresses
    lichenhao       pptpd   lichenhao       *

client # 指定用户名

server # 指定该用户连接时所用的协议

secret # 指定该用户密码

IP addresses #所用的 ip,ip 可以用\*表示该用户可用所有 ip

## 其他系统配置

添加路由转发功能

1. cat >> /etc/sysctl.conf << EOF
2. net.ipv4.ip_forward = 1
3. EOF
4. sysctl -p

配置防火墙

iptables -t nat -A POSTROUTING -s 192.168.0.0/24 -o eth0 -j MASQUERADE

VPN iptables NAT 配置

iptables -t nat -A POSTROUTING -s 10.254.13.0/24 ! -d 10.254.13.0/24 -j SNAT --to-source 123.151.161.92

iptables -t nat -A POSTROUTING -s 192.168.16.0/20 -j MASQUERADE

windows 配置
