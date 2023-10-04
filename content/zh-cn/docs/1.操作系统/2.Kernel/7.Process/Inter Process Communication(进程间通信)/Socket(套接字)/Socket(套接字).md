---
title: Socket(套接字)
---

# 概述

> 参考：
>
> - [Wiki，Socket](https://en.wikipedia.org/wiki/Socket)
> - [Wiki，Unix domain Socket](https://en.wikipedia.org/wiki/Unix_domain_socket)
> - [Wiki，Network Scoket](https://en.wikipedia.org/wiki/Network_socket)

**Socket(套接字)** **是数据通信的基石**。是计算机领域中数据通信的一种约定，或者说是一种方法。通过 Socket 这种方法，计算机内的进程可以互相交互数据，不同计算机之间也可以互相交互数据。

Socket(套接字) 原意是`插座`，所以 Socket 就像插座的作用一样，只要把插头插上，就能让设备获得电力。同理，只要两个程序通过 Socket 互相套接，也就是说两个程序都插在同一个 Socket 上，那么这两个程序就能交互数据。

在计算机领域，Socket 有多种类型

- **Unix Domain Socket(简称 UDS)** # 用于同一台设备的不同进程间互相通信
- **Network Socket** # 用于进程在网络间互相通信
- **Berkeley Sockets API** # Unix Domain Socket 与 Network Socket 的 API

在软件上，Socket 负责套接计算机中的数据(可以想象成套接管，套接管即为套管，是用来把两个管连接起来的东西，套接字就是把计算机中的字(即最小数据)连接起来，且只把头部连接起来，套管也是，只把两根很长的管的头端套起来接上)

1. 在系统层面，socket 可以连接系统中的两个进程，进程与进程本身是互相独立的，如果需要传递消息，那么就需要两个进程各自打开一个接口(API)，socket 把两个进程的 api 套住使之连接起来，即可实现进程间的通信。该 socket 是抽象的，虚拟的，只是通过编程函数来实现进程的 API 功能，如果进程没有 API，那么就无法通过 socket 与其余进程通信。
2. 当然，一个进程也可以监听一个名为 _.scok 的文件，这个文件就像 API 一样，其他程序想与该进程交互，只要指定该_.sock 文件，然后对这个 sock 文件进行读写即可。
3. 在网络层面，socket 负责把不在同一主机上的进程(比如主机 A 的进程 C 和主机 B 的进程 D)连接起来，而两个不同主机上的进程如何被套接起来呢，套接至少需要提供一个头端来让套接管(字)包裹住才行。这时候(协议，IP，端口,例如：ftp://192.168.0.1:22)共同组成了网络上的进程标示，该进程逻辑上的头端即为紫色部分的端口号，不同主机的两个进程可以通过套接字把端口号套起来连接，来使两个网络上不同主机的进程进行通信，该同能同样是在程序编程的时候用函数写好的，程序启动为进程的时候，则该接口会被拿出来监听。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nlg3b5/1619421243110-2db70bc6-f358-459c-b9a9-e199658b151a.png)
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nlg3b5/1619421247179-b40abf99-2621-4f4e-aa6e-1d68bfe9e74b.png)

## Unix Domain Socket

**Unix Domain Socket** 是 **IPC** 的一种实现方式。Socket 原本是为了网络通信设计的，但后来在 Socket 的框架上发展出一种 IPC 机制，就是 Unix Domain Socket。虽然 Netork Socket 也可用于统一台主机的进程间通信(通过 loopback 地址 127.0.0.1)，但是 Unix Domain Socket 用于 IPC 更有效率，因为不需要经过网络协议栈，不需要打包拆包、计算校验和、维护序号和应答等，只是将应用层数据从一个进程拷贝到另一个进程。这是因为 IPC 机制本质上是可靠的通讯，而网络协议是为不可靠通讯设计的。

Unix Domain Socket 是全双工的，API 接口语义丰富，相比其它 IPC 机制有明显的优越性，目前已成为使用最广泛的 IPC 机制，比如 X Window 服务器和 GUI 程序之间就是通过 UNIX domain socket 通讯的。

Unix domain socket 是 POSIX 标准中的一个组件，所以不要被名字迷惑，linux 系统也是支持它的。

## Network Socket

详见 [Network Socket](/docs/4.数据通信/数据通信/Network%20Socket.md Socket.md)

## Berkeley Sockets API

**Berkeley Sockets**是 Network Socket 和 Unix Domain Sockets 的 应用程序编程接口（API），用于进程间通信（IPC）。通常将其实现为可链接模块的库。它起源于 1983 年发布的 4.2BSD Unix 操作系统。

套接字是网络通信路径的本地终结点的抽象表示（句柄）。Berkeley 套接字 API 将其表示为 Unix 哲学中的文件描述符（文件句柄），该描述符为输入和输出到数据流提供通用接口。

伯克利套接字几乎没有任何改动，从\_事实上的\_标准演变为 POSIX 规范的组件。术语 **POSIX 套接字**在本质上是\_Berkeley 套接字的\_同义词，但是它们也称为\_BSD 套接字\_，这是对 Berkeley Software Distribution 中的第一个实现的认可。
