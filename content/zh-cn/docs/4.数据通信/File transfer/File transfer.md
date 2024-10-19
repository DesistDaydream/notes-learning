---
title: File transfer
linkTitle: File transfer
date: 2024-05-08T12:31
weight: 1
---

# 概述

> 参考：
>
> - [Wiki, File_transfer](https://en.wikipedia.org/wiki/File_transfer)
> - [Wiki, Protocol for file transfer](https://en.wikipedia.org/wiki/Protocol_for_file_transfer)

**File transfer(文件传输)** 是通过[数据通信](/docs/4.数据通信/数据通信/数据通信.md)的通道将计算机文件从一个计算机系统传输到另一个计算机系统的行为。文件传输是通过[通信协议](/docs/4.数据通信/Protocol/通信协议.md)协调的。许多文件传输协议是针对不同环境而设计的

**File transfer protocol(文件传输协议)** 是描述如何在两个计算 endpoint 之间传输文件的约定。除了作为单个单元存储在文件系统中的文件的比特流之外，有些还可能发送相关元数据，例如文件名、文件大小和时间戳，甚至文件系统权限和文件属性。

文件传输协议分为两大种

- **Packet switched Protocol(分组交换网络协议)**
  - [FTP](/docs/4.数据通信/Protocol/FTP.md)
- **Serial Protocol(串行协议)**
  - Modems(拨号调制解调器) 使用 XMODEM、YMODEM、ZMODEM 和类似的空调制解调器链接。TODO
- 其他
  - USB 等外部存储设备与计算机之间文件互传
- etc.

# X/Y/Z Modem

> 参考:
>
> - [Wiki, XMODEM](https://en.wikipedia.org/wiki/XMODEM)
> - [Wiki, YMODEM](https://en.wikipedia.org/wiki/YMODEM)
> - [Wiki, ZMODEM](https://en.wikipedia.org/wiki/ZMODEM)
> - https://pauillac.inria.fr/~doligez/zmodem/ymodem.txt
> - [知乎，Xmodem 协议介绍及应用（基于 ESP-IDF）](https://zhuanlan.zhihu.com/p/349921713)

X/Y/Z Modem 并不依赖于 TCP/IP 进行传输，早期是用来在串行通信（比如调制解调器）中点对点传输文件的协议

# 下载

## Aria2

> 参考：
>
> - [GitHub 项目，aria2/aria2](https://github.com/aria2/aria2)

aria2 是一个轻量级的多协议和多源跨平台下载实用程序，在命令行中运行。它支持 HTTP/HTTPS、FTP、SFTP、BitTorrent 和 Metalink。

可以下载 BT 种子

`aria2c xxx.torrnet`

使用 Aria2 的客户端

- https://github.com/agalwood/Motrix/
