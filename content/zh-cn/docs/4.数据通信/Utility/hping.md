---
title: hping
linkTitle: hping
date: 2024-04-18T09:11
weight: 20
tags:
  - Network_analyzer
---

# 概述

> 参考：
>
> - [GitHub 项目，antirez/hping](https://github.com/antirez/hping)
>   - https://www.kali.org/tools/hping3/
>   - https://salsa.debian.org/debian/hping3
> - [Wiki, Hping](https://en.wikipedia.org/wiki/Hping)

原始 GitHub 项目截至 2024-04-18 已经有 10 年没更新了。

hping 不支持 IPv6

https://github.com/TeddyGuo/tping 像 hping，但支持 IPv6，Python 编写

[nping](/docs/4.数据通信/Utility/nping.md) 由 Nmap

hping 是由 Salvatore Sanfilippo（也称为 Antirez）创建的 TCP/IP 协议的开源数据包生成器和分析器。它是用于防火墙和网络安全审计和测试的常用工具之一，用于利用空闲扫描扫描技术（也是由 hping 作者发明的），现在在 Nmap Security Scanner 中实现。新版本的 hping，hping3，可以使用 Tcl 语言编写脚本，并实现一个基于字符串的、人类可读的 TCP/IP 数据包描述引擎，以便程序员可以编写与低级 TCP/IP 数据包操作和分析相关的脚本。很短的时间。

# Syntax(语法)



# EXAMPLE

https://www.cnblogs.com/Higgerw/p/16469371.html

通过 eth1，发送 SYN 报文到 172.24.194.179 的 80 端口。源地址伪造为 192.168.180.131，时间间隔 1000us

- `hping3 -I eth0 -S 172.24.194.179 -p 80 -a 192.168.180.131 -i u1000`

通过 eth1, 发送 SYN  报文到 192.168.180.133:80, 使用随机源地址, 时间间隔 1000us

- `hping3 -I eth1 -S 192.168.180.133 -p 80 --rand-source -i u1000`