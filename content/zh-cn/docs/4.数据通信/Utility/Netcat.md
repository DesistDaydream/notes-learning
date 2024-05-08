---
title: Netcat
linkTitle: Netcat
date: 2024-03-20T08:57
weight: 20
tags:
  - Network_analyzer
---

# 概述

> 参考：
> 
> - [Wike，Netcat](https://en.wikipedia.org/wiki/Netcat)
> - [Ncat Manual(手册)](https://nmap.org/book/ncat-man.html)
> - <https://zhuanlan.zhihu.com/p/83959309>

Netcat 是一个简单的实用程序通过 TCP 或 UDP 网络连接读取和写入数据。它旨在成为一个可靠的后端工具，可直接使用或轻松由其他程序和脚本驱动。同时，它还是一个功能丰富的网络调试和探索工具，因为它几乎可以创建您需要的任何类型的连接，包括端口绑定以接受传入连接。

由于 Netcat 的设计理念和功能，被人亲切的亲切得称为 **网络工具中的瑞士军刀**

最初的 Netcat 是由 hobbit 于 1995 年[发布](http://seclists.org/bugtraq/1995/Oct/0028.html)的，尽管它很受欢迎，但它并没有得到维护。有时甚至很难找到[v1.10 源代码的副本](http://download.insecure.org/stf/nc110.tgz)。该工具的灵活性和实用性促使 Nmap 项目产生 [Ncat](#Ncat))，这是一种支持 SSL、IPv6、SOCKS 和 http 代理、连接代理等的现代重新实现。除了 Nmap 项目重新了 Netcat，还有很多重写甚至扩展了 Netcat 的工具

- [Socat](https://sectools.org/tool/socat/)
- [OpenBSD 的 Netcat](http://www.openbsd.org/cgi-bin/cvsweb/src/usr.bin/nc/)
- [Cryptcat](http://cryptcat.sourceforge.net/)
- [Netcat6](http://www.deepspace6.net/projects/netcat6.html)
- [pnetcat](http://stromberg.dnsalias.org/~strombrg/pnetcat.html)
- [SBD](http://cycom.se/dl/sbd)
- 所谓的[GNU Netcat](http://netcat.sourceforge.net/)

如需下载和更多信息， [访问 Netcat 主页](http://en.wikipedia.org/wiki/Netcat)。

## Ncat

> 参考：
> 
> - [Nmap，ncat](http://nmap.org/ncat/)

Ncat 是一个功能丰富的网络实用程序，它可以从命令行跨网络读取和写入数据。Ncat 由 Nmap 项目编写的，是对古老的 [Netcat](http://sectools.org/tool/netcat/) 的大大改进的重新实现。它同时使用 TCP 和 UDP 进行通信，并被设计为一种可靠的后端工具，可立即为其他应用程序和用户提供网络连接。Ncat 不仅适用于 IPv4 和 IPv6，还为用户提供了几乎无限的潜在用途。

在 Ncat 的众多功能中，包括将 Ncat 链接在一起、将 TCP 和 UDP 端口重定向到其他站点、SSL 支持以及通过 SOCKS4 或 HTTP（CONNECT 方法）代理（以及可选的代理身份验证）进行代理连接的能力。一些通用原则适用于大多数应用程序，从而使您能够立即向通常不支持它的软件添加网络支持。

Ncat 与 Nmap 集成，可在 Nmap 下载页面提供的标准 Nmap 下载包（包括源代码和 Linux、Windows 和 Mac 二进制文件）中找到。。也可以再 [SVN 源代码存储库中](http://nmap.org/book/install.html#inst-svn)找到它。

许多用户要求提供一个静态编译的 ncat.exe 版本，他们可以将其放在 Windows 系统上并使用，而无需运行任何安装程序或复制额外的库文件。我们已经构建了一个静态编译的 Windows 二进制版本。可以在[此处](http://nmap.org/dist/ncat-portable-5.59BETA1.zip)下载 zip 文件。为确保文件未被篡改，可以检查[加密签名](http://nmap.org/book/install.html#inst-integrity)。如果您需要更新的 Ncat 版本的便携版本，请参阅 [Ncat 便携编译说明](https://secwiki.org/w/Nmap/Ncat_Portable)。

该 [NCAT 用户指南](https://nmap.org/ncat/guide/index.html)包含完整的文档，包括很多技巧，窍门和实用现实生活的例子！还有一个[Ncat 手册页](https://nmap.org/book/ncat-man.html)用于快速使用摘要。

## OpenBSD Netcat

> 参考：
>
> - [OpenBSD-nc Manual(手册)](https://man.openbsd.org/nc)

由 OpenBSD 对原始 Netcat 的更新支持

# Netcat 安装

## Ubuntu

Ubuntu 使用 OpenBSD 的 Netcat 作为 Netcat 的替代品

安装 netcat-openbsd 包即可，安装完成后，nc 命令本质上是 nc.openbsd 命令的软链接

```bash
 ~]# ll /usr/bin/nc
lrwxrwxrwx 1 root root 20 Aug 10  2023 /usr/bin/nc -> /etc/alternatives/nc*
~]# ll /etc/alternatives/nc
lrwxrwxrwx 1 root root 15 Aug 10  2023 /etc/alternatives/nc -> /bin/nc.openbsd*
```

## CentOS

CentOS 使用 Nmap 的 Ncat 作为 Netcat 的替代品

安装 nmap-ncat 包即可，安装完成后，nc 命令本质上是 ncat 命令的软链接

# Syntax(语法)

**nc \[\<OPTIONS> ...] \[ \<hostname> ] \[ \<port> ]**

## OPTIONS

- **-k, --keep-open** # 通常与 -l 选项配合使用。在监听模式下接受多个连接。若不使用 -k 选项，则第一个连接断开后，监听也就结束了。
- **-l, --listen \<PORT>** # 让程序监听指定的端口
- **-u, --udp** # 使用 UDP，而不是默认的 TCP

Ncat 与 OpenBSD-nc 这两个程序的选项有不同的地方

### Ncat OPTIONS

- **--proxy <ADDRESS:\[PORT]>** # 连接目的地时所使用代理 IP 和 PORT。
- **--proxy-type \<STRING>** # 连接目的地时所使用的代理类型(也就是代理协议)。可用的值有：
  - socks4 # 表示 SOCKS v.4
  - socks5 # 表示 SOCKS v.5（默认值）
  - http # 表示 HTTP

### OpenBSD-nc OPTIONS

- **-x \<ADDRESS\[:PORT]>** # 连接目的地时所使用代理 IP 和 PORT。代理不能与 -LsuU 这些选项一起使用。
- **-X \<PROXY_PROTOCOL>** # 连接目的地时所使用的代理协议。可用的值有：
  - 4 # 表示 SOCKS v.4
  - 5 # 表示 SOCKS v.5（默认值）
  - connect # 表示 HTTP

# 应用示例

- 测试本地 323/udp 端口
  - nc -uvz localhost 323
- 测试本地 22/tcp 端口
  - nc -vz localhost 22

执行效果如下：

```bash
~]# nc -vz localhost 22
Connection to localhost 22 port [tcp/ssh] succeeded!

~]# nc -uvz localhost 323
Connection to localhost 323 port [udp/*] succeeded!
```

测试两台机器之间的 8080/udp 连接是否正常

- 在主机 A 上监听 8080/udp
  - nc -u -l 8080
- 在主机 B 上测试主机 A 的 8080/udp 是否正常
  - nc -u 172.19.42.248 8080
- 然后在任意主机输入任意内容，只要另一个主机能看到相同内容，即表示连接正常

# 常见问题

在 Windows 使用 ncat 通过代理连接 ssh，报错 `ssh_exchange_identification: Connection closed by remote host`

问题 issue：<https://github.com/nmap/nmap/issues/2149>

解决：下载 7.80 版本即可，将下载连接的版本号改为 7.80 即可下载。

