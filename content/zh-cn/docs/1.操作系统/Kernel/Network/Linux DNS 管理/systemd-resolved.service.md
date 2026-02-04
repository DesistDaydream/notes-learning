---
title: systemd-resolved.service
linkTitle: systemd-resolved.service
weight: 100
---

# 概述

> 参考：
>
> - [Linux man pages，systemd-resolved.service(8)](https://man7.org/linux/man-pages/man8/systemd-resolved.service.8.html)
> - [金步国 - system 中文手册，systemd-resolved.service](http://www.jinbuguo.com/systemd/systemd-resolved.service.html)

systemd-resolved.service 是一个类似于 [DNSmasq](/docs/4.数据通信/DNS/DNSmasq.md) 的域名解析服务，只不过这个服务只适用于 Linux 中，且被 systemd 所管理。

systemd-resolved 是一个为本地应用程序提供网络名称解析的系统服务。它实现了一个缓存和验证 DNS/DNSSEC 存根解析器，以及一个 LLMNR 和 MulticastDNS 解析器和响应器。

# 关联文件与配置

**resolved.conf** # 运行时配置文件。

- 程序会从如下位置逐一加载，并且只使用找到的第一个文件
    - /etc/systemd/resolved.conf
    - /run/systemd/resolved.conf
    - /usr/local/lib/systemd/resolved.conf
    - /usr/lib/systemd/resolved.conf
- 除了主配置文件之外，还会从下面的目录读取可替换的配置片段。这些可替换的配置片段优先级更高，会覆盖主配置文件。`resolved.conf.d/` 配置子目录中的文件会按文件名字典顺序排序，而不管它们位于哪个子目录中。当多个文件指定相同的选项时，对于只接受单个值的选项，排序在最后一个文件中的条目优先级最高；对于接受值列表的选项，则会按照排序后的文件顺序收集条目。
    - `/usr/lib/systemd/resolved.conf.d/*.conf`
    - `/usr/local/lib/systemd/resolved.conf.d/*.conf`
    - `/run/systemd/resolved.conf.d/*.conf`
    - `/etc/systemd/resolved.conf.d/*.conf`

**resolv.conf** # TODO: systemd-resolved 对于 resolve.conf 的文件好像有比较繁琐的规则，而且我能在多个目录下找到多个名为 resolv.conf 的文件。需要梳理

- TODO: 好像是，systemd-resolved 根据 resolved.conf 文件的配置生成 resolv.conf 文件

**/run/systemd/resolve/** # /etc/resolv.conf 文件将会软链接到此目录下的某个文件。通常默认链接到 stub-resolv.conf 文件。

- **./stub-resolv.conf** # 通常只包含一个本地地址(127.0.0.53)，指向 systemd-resolved.service 提供的 DNS Stub Listener。
- **./resolv.conf** # 包含了系统实际使用的 DNS 服务器和搜索域的列表。systemd-resolved.service 的 DNS Stub Listener 将会读取该文件中的 DNS Server

**/usr/lib/systemd/resolv.conf** # 代替 /etc/resolv.conf 文件

# 命令行工具

## resolvectl

> 参考：
> 
> - [Linux man pages，resolvectl(1)](https://man7.org/linux/man-pages/man1/resolvectl.1.html)

### Syntax(语法)

**resolvectl [OPTIONS] COMMAND [NAME...]**

**COMMAND**

- **query HOSTNAME | ADDRESS** # 解析域名以及 IPv4 和 IPv6 地址。当与 --type= 或 --class= （见下文）结合使用时，解析低级 DNS 资源记录。

### EXAMPLE

# resolved.conf

> 参考：
>
> - [Linux man pages, resolved.conf(5)](https://man7.org/linux/man-pages/man5/resolved.conf.5.html)

这是类似 [INI](docs/2.编程/无法分类的语言/INI.md) 格式的配置文件。具体的解析配置在 Resolve.DNS 键上
