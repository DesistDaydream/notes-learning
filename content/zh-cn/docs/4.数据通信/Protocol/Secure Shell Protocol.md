---
title: Secure Shell Protocol
linkTitle: Secure Shell Protocol
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, Secure_Shell_Protocol](https://en.wikipedia.org/wiki/Secure_Shell_Protocol)
> - <https://www.digitalocean.com/community/tutorials/how-to-set-up-ssh-keys-on-centos-8>

**Secure Shell Protocol(安全外壳协议，简称 SSH)** 是一种加密的[Communication protocol](/docs/4.数据通信/Protocol/Communication%20protocol.md)，可在不安全的网络中为网络服务提供安全的传输环境。SSH 通过在网络中创建安全隧道来实现 SSH 客户端与服务器之间的连接。虽然任何网络服务都可以通过 SSH 实现安全传输，SSH 最常见的用途是远程登录系统，人们通常利用 SSH 来传输命令行界面和远程执行命令。使用频率最高的场合类 Unix 系统，但是 Windows 操作系统也能有限度地使用 SSH。2015 年，微软宣布将在未来的操作系统中提供原生 SSH 协议支持，Windows 10 1809 版本已提供可手动安装的 OpenSSH 工具

# SSH 的实现

[OpenSSH](/docs/4.数据通信/Utility/OpenSSH/OpenSSH.md)

Go 语言的 SSH 实现

- https://github.com/search?q=sshd+language%3AGo&ref=opensearch&type=repositories
- [GitHub 项目，jpillora/sshd-lite](https://github.com/jpillora/sshd-lite)
  - 对应博客 https://blog.gopheracademy.com/go-and-ssh/
- https://github.com/nwtgck/handy-sshd
- https://github.com/Matir/sshdog

# 其他

https://github.com/shazow/ssh-chat
