---
title: FTP
linkTitle: FTP
date: 2024-02-23T12:26
weight: 20
---

# 概述

> 参考：
>
> - [RFC 959, FILE TRANSFER PROTOCOL (FTP)](https://datatracker.ietf.org/doc/html/rfc959)
> - [Wiki, File_Transfer_Protocol](https://en.wikipedia.org/wiki/File_Transfer_Protocol)

**File Transfer Protocol(文件传输协议，简称 FTP)** 是因特网网络上历史最悠久的网络工具，从 1971 年由 A KBHUSHAN 提出第一个 FTP 的 RFC（RFC114）至今近半个世纪来，FTP 凭借其独特的优势一直都是因特网中最重要、最广泛的服务之一。

FTP 的目标是提高文件的共享性，提供非直接使用远程计算机，使存储介质对用户透明和可靠高效地传送数据。它能操作任何类型的文件而不需要进一步处理，就像 MIME 或 Unicode 一样。但是，FTP 有着极高的延时，这意味着，从开始请求到第一次接收需求数据之间的时间，会非常长；并且不时的必须执行一些冗长的登录进程。

# SFTP

> 参考：
>
> - [Wiki, SSH_File_Transfer_Protocol](https://en.wikipedia.org/wiki/SSH_File_Transfer_Protocol)

**SSH File Transfer Protocol(SSH 文件传输协议，简称 SFTP)** 也称为 **Secure File Transfer Protocol(安全文件传输协议)**，是一种网络协议，可通过任何可靠的数据流提供文件访问、文件传输和文件管理。它由 [IETF](/docs/Standard/Internet/IETF.md) 设计，作为 [Secure Shell Protocol](/docs/4.数据通信/Protocol/Secure%20Shell%20Protocol.md)(SSH) 2.0 版的扩展，提供安全文件传输功能，并且由于卓越的安全性而被视为文件传输协议 (FTP) 的替代品。

IETF 互联网草案指出，尽管该协议是在 SSH-2 协议的上下文中描述的，但它可以用于许多不同的应用程序，例如通过传输层安全性 (TLS) 进行安全文件传输和管理传输VPN 应用程序中的信息。

该协议假定它在安全通道（例如 SSH）上运行，服务器已经对客户端进行了身份验证，并且客户端用户的身份可供该协议使用。

也就是说，想要使用 SFTP，通常是先建立安全的连接（e.g. SSH、etc.），然后基于该安全的连接实现 FTP 相关能力。

> SFTP 不一定是通过 SSH 运行的 FTP（绝大部分都是通过 SSH），而是由 IETF SECSH 工作组从头开始设计的新协议。它有时会与[简单文件传输协议](https://en.wikipedia.org/wiki/Simple_File_Transfer_Protocol)混淆。

SFTP 尝试比 SCP 更加独立于平台；例如，对于 SCP，客户端指定的通配符扩展由服务器决定，而 SFTP 的设计避免了这个问题。虽然 SCP 最常在 [Unix-like OS](/docs/1.操作系统/Operating%20system/Unix-like%20OS/Unix-like%20OS.md) 上实现，但 SFTP 服务器通常在大多数平台上可用。在 SFTP 中，可以轻松终止文件传输，而无需像其他机制那样终止会话。

# FTP 的实现

OpenSSH 的 [SFTP Subsystem](/docs/4.数据通信/Utility/OpenSSH/SFTP%20Subsystem.md) # 可以实现 SFTP

[vsftpd](/docs/4.数据通信/Utility/vsftpd.md)
