---
title: lrzsz
linkTitle: lrzsz
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，UweOhse/lrzsz](https://github.com/UweOhse/lrzsz)
> - [官网](https://ohse.de/uwe/software/lrzsz.html)

> [!warning]
> 若想通过 [ssh](/docs/4.数据通信/Utility/OpenSSH/ssh.md) 使用 lrzsz 中的程序，必须要保证执行 ssh 命令的程序支持将 X/Y/Z Modem 协议封装到 SSH 中，比如 XShell、SecurityCRT、etc. 。如果是通过 [PowerShell](/docs/1.操作系统/Terminal%20与%20Shell/WindowsShell/PowerShell/PowerShell.md) 使用 ssh 命令连接到服务器，是无法使用 lrzsz 的。

lrzsz 是一个用在 [Unix-like OS](/docs/1.操作系统/Operating%20system/Unix-like%20OS/Unix-like%20OS.md) 上的 [File transfer](/docs/4.数据通信/File%20transfer/File%20transfer.md) 工具包，使用 X/Y/Z Modem 文件传输协议从本地直接上传/下载文件到操作系统中

lrzsz 中包含如下程序

> r 开头的是 receive，让 lrzsz 所在系统从外部接受文件；s 开头的是 send，从 lrzsz 所在系统往外发送文件。都是相对 lrzsz 程序所在系统所说，以 lrzsz 程序为主语。
>
> b 是 Ymodem、x 是 Xmodem、z 是 Zmodem

- rb
- rx
- rz
- sb
- sx
- sz

# Syntax(语法)

> 参考:
>
> - [Ubuntu 官网 Manual(手册)，rz(1)](https://manpages.ubuntu.com/manpages/focal/man1/rz.1.html)

**OPTIONS**

- **-b,** **--binary** # 二进制（按原样告诉它）文件传输覆盖。
- **-e, --escape** # 强制发送者转义所有控制字符；通常 XON、XOFF、DLE、CR-@-CR 和 Ctrl-X 会被转义。

# EXAMPLE

`rz -be`
