---
title: Unix Domain Socket
linkTitle: Unix Domain Socket
weight: 20
---

# 概述

> 参考：
>
> - [Manual, unix(7)](https://man7.org/linux/man-pages/man7/unix.7.html)
> - [Wiki, Unix domain Socket](https://en.wikipedia.org/wiki/Unix_domain_socket)

**Unix Domain Socket** 是 **IPC** 的一种实现方式。Socket 原本是为了网络通信设计的，但后来在 Socket 的框架上发展出一种 IPC 机制，就是 Unix Domain Socket。虽然 Netork Socket 也可用于同一台主机的进程间通信(通过 loopback 地址 127.0.0.1)，但是 Unix Domain Socket 用于 IPC 更有效率，因为不需要经过网络协议栈，不需要打包拆包、计算校验和、维护序号和应答等，只是将应用层数据从一个进程拷贝到另一个进程。这是因为 IPC 机制本质上是可靠的通讯，而网络协议是为不可靠通讯设计的。

Unix Domain Socket 是全双工的，API 接口语义丰富，相比其它 IPC 机制有明显的优越性，目前已成为使用最广泛的 IPC 机制，比如 X Window 服务器和 GUI 程序之间就是通过 UNIX domain socket 通讯的。

Unix domain socket 是 POSIX 标准中的一个组件，所以不要被名字迷惑，linux 系统也是支持它的。