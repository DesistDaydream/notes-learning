---
title: FTP
linkTitle: FTP
date: 2024-02-23T12:26
weight: 20
---

# 概述

> 参考：
>
> - [RFC 959，FILE TRANSFER PROTOCOL (FTP)](https://datatracker.ietf.org/doc/html/rfc959)
> - [Wiki，File_Transfer_Protocol](https://en.wikipedia.org/wiki/File_Transfer_Protocol)

**File Transfer Protocol(文件传输协议，简称 FTP)** 是因特网网络上历史最悠久的网络工具，从 1971 年由 A KBHUSHAN 提出第一个 FTP 的 RFC（RFC114）至今近半个世纪来，FTP 凭借其独特的优势一直都是因特网中最重要、最广泛的服务之一。

FTP 的目标是提高文件的共享性，提供非直接使用远程计算机，使存储介质对用户透明和可靠高效地传送数据。它能操作任何类型的文件而不需要进一步处理，就像 MIME 或 Unicode 一样。但是，FTP 有着极高的延时，这意味着，从开始请求到第一次接收需求数据之间的时间，会非常长；并且不时的必须执行一些冗长的登录进程。

# FTP 协议的实现工具

vsftp

# vsftp - 实现 FTP 协议的工具

## vsftp 关联文件与配置

**/etc/vsftpd.conf** # vsftpd 程序的配置文件

**/etc/ftpusers** # 此文件包含*禁止*FTP 登录的用户名，通常有"root"，"uucp"，"news" 之类，因为这些用户权限太高，登录 FTP 误操作危险性大。

- User NAME
- User NAME
- User NAME
- .......
- User NAME

配置文件 keywords 说明

最小基本配置

- write_enable=YES # 对 ftp 服务器是否有写的权限
- local_enable=YES # 是否允许本地用户登录(本地用户即 ftp 服务器自身的用户)
- anonymous_enable=NO # 是否允许匿名登录

扩展配置

- chroot_local_user=YES # 是否启动本地用户 chroot 规则，chroot 改变登录 ftp 的本地用户根的目录位置
- allow_writeable_chroot=YES # 允许在限定目录有写权限
- chroot_list_enable=YES # 是否启动不受 chroot 规则的用户名单
- chroot_list_file=/etc/vsftpd.chroot_list # 定义不受限制的用户名单在哪个文件中
- pam_service_name=vsftpd 改为 pam_service_name=ftp # 如果始终登录时候提示密码错误，则修改此项

vsftp 可以使用 [Chroot](/docs/1.操作系统/Linux%20管理/Linux%20系统管理工具/Chroot.md) 功能。比如：下面第一个图是不启动 Chroot 规则的情况，第二张图是启用 Chroot 规则的情况，可以看到当使用 Chroot 时，`/srv/ftp/` 目录对于 ftp 程序来说是作为 `/` 存在的。由于这个原因，所以启动 Chroot 的时候，ftp 工具无法访问所设定的 `/` 目录以外的其他目录

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pvqe8m/1616165219993-ce6cd857-e9ba-4af0-b7fc-7d77cf547d84.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pvqe8m/1616165220004-51f8038e-598e-427a-9b04-8f1987475f04.jpeg)

- /etc/vsftpd.chroot_list # 此文件包含对服务器上所有 FTP 内容有权限的用户名，该文件中每个用户名写一行

## FTP 命令行工具说明

ftp \[\[UserName@]目标 IP] # 可以使用 UserName 用户来登录目标 IP 的 ftp 服务器或直接进入 ftp 客户端

- open \[IP] # 连接某个 FTP 服务器。
- close # 关闭连接
- user # 对已经连接的 FTP 再输入一次用户名和密码，类似于更改登陆账户
- ？ # 查看可以使用的命令
- get \[需要下载的文件] \[下载路径以及定义文件名] # 下载文件
- put \[需要上传的文件] \[上传路径以及定义文件名] # 上传文件

注意：这些命令可以直接输入然后根据提示再输入想要执行的内容；也可直接在命令后接需要完成的动作，直接执行

注意文件夹权限以免无法上传下载

# ftp 被动模式注意事项

如果 iptables 仅仅放了 21 端口，在启动被动模式后，需要给 iptables 添加模块以便被动模式正常运行

在 /etc/sysconfig/iptalbes-conf 文件中，修改成 IPTABLES_MODULES="nf_conntrack_ftp nf_nat_ftp" 这个内容

# SFTP

> 参考：
>
> - [Wiki，SSH_File_Transfer_Protocol](https://en.wikipedia.org/wiki/SSH_File_Transfer_Protocol)

**SSH File Transfer Protocol(SSH 文件传输协议，简称 SFTP)** 也称为 **Secure File Transfer Protocol(安全文件传输协议)**，是一种网络协议，可通过任何可靠的数据流提供文件访问、文件传输和文件管理。它由 [IETF](/docs/Standard/Internet/IETF.md) 设计，作为 [Secure Shell Protocol](/docs/4.数据通信/通信协议/Secure%20Shell%20Protocol.md)(SSH) 2.0 版的扩展，提供安全文件传输功能，并且由于卓越的安全性而被视为文件传输协议 (FTP) 的替代品。

IETF 互联网草案指出，尽管该协议是在 SSH-2 协议的上下文中描述的，但它可以用于许多不同的应用程序，例如通过传输层安全性 (TLS) 进行安全文件传输和管理传输VPN 应用程序中的信息。

该协议假定它在安全通道（例如 SSH）上运行，服务器已经对客户端进行了身份验证，并且客户端用户的身份可供该协议使用。

也就是说，想要使用 SFTP，通常是先建立安全的连接（e.g. SSH、etc.），然后基于该安全的连接实现 FTP 相关能力。

> SFTP 不一定是通过 SSH 运行的 FTP（绝大部分都是通过 SSH），而是由 IETF SECSH 工作组从头开始设计的新协议。它有时会与[简单文件传输协议](https://en.wikipedia.org/wiki/Simple_File_Transfer_Protocol)混淆。

SFTP 尝试比 SCP 更加独立于平台；例如，对于 SCP，客户端指定的通配符扩展由服务器决定，而 SFTP 的设计避免了这个问题。虽然 SCP 最常在 [Unix-like OS](/docs/1.操作系统/Operating%20system/Unix-like%20OS/Unix-like%20OS.md) 上实现，但 SFTP 服务器通常在大多数平台上可用。在 SFTP 中，可以轻松终止文件传输，而无需像其他机制那样终止会话。