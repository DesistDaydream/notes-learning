---
title: systemd-resolved.service
---

# 概述

> 参考：
> 
> - [Manual(手册)，systemd-resolved.service(8)](https://man7.org/linux/man-pages/man8/systemd-resolved.service.8.html)
> - [金步国-system 中文手册，systemd-resolved.service](http://www.jinbuguo.com/systemd/systemd-resolved.service.html)

systemd-resolved.service 是一个类似于 [DNSmasq](/docs/4.数据通信/DNS/DNSmasq.md) 的域名解析服务，只不过这个服务只适用于 Linux 中，且被 systemd 所管理。

systemd-resolved 是一个为本地应用程序提供网络名称解析的系统服务。它实现了一个缓存和验证 DNS/DNSSEC 存根解析器，以及一个 LLMNR 和 MulticastDNS 解析器和响应器。

# 关联文件与配置

**/etc/systemd/resolved.conf** # 运行时配置文件。

**/run/systemd/resolve/** # /etc/resolv.conf 文件将会软链接到此目录下的某个文件。通常默认链接到 stub-resolv.conf 文件。

- **./stub-resolv.conf** # 通常只包含一个本地地址(127.0.0.53)，指向 systemd-resolved.service 提供的 DNS Stub Listener。
- **./resolv.conf** # 包含了系统实际使用的 DNS 服务器和搜索域的列表。systemd-resolved.service 的 DNS Stub Listener 将会读取该文件中的 DNS Server

**/usr/lib/systemd/resolv.conf** # 代替 /etc/resolv.conf 文件

# 命令行工具

## resolvectl

> 参考：
> 
> - [Manual(手册)，resolvectl(1)](https://man7.org/linux/man-pages/man1/resolvectl.1.html)

### Syntax(语法)

**resolvectl [OPTIONS] COMMAND [NAME...]**

**COMMAND**

- **query HOSTNAME | ADDRESS** # 解析域名以及 IPv4 和 IPv6 地址。当与 --type= 或 --class= （见下文）结合使用时，解析低级 DNS 资源记录。

### EXAMPLE

