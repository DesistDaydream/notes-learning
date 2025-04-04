---
title: TUN and TAP
linkTitle: TUN and TAP
weight: 20
tags:
  - virtualization
---

# 概述

> 参考：
>
> - [网络虚拟化技术（二）: TUN/TAP MACVLAN MACVTAP](https://blog.kghost.info/2013/03/27/linux-network-tun/)
> - [Wiki, TUN/TAP](https://en.wikipedia.org/wiki/TUN/TAP)

TUN/TAP 类型的设备会利用 /dev/net/tun 文件, 基于 tun.ko 模块实现(TODO).

## TUN 设备

TUN 设备是一种虚拟网络设备，通过此设备，程序可以方便地模拟网络行为。TUN 模拟的是一个三层设备,也就是说,通过它可以处理来自网络层的数据，更通俗一点的说，通过它，我们可以处理 IP 数据包。

先来看看物理设备是如何工作的：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/xu2gy9/1616124241972-bd04293d-efc8-45da-8088-314f83bf47f9.png)

上图中的 eth0 表示我们主机已有的真实的网卡接口 (interface)。

网卡接口 eth0 所代表的真实网卡通过网线(wire)和外部网络相连，该物理网卡收到的数据包会经由接口 eth0 传递给内核的网络协议栈(Network Stack)。然后协议栈对这些数据包进行进一步的处理。

对于一些错误的数据包,协议栈可以选择丢弃；对于不属于本机的数据包，协议栈可以选择转发；而对于确实是传递给本机的数据包,而且该数据包确实被上层的应用所需要，协议栈会通过 Socket API 告知上层正在等待的应用程序。

下面看看 TUN 的工作方式：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/xu2gy9/1616124241917-45f4416e-f176-4d02-8ad5-0a31cc478397.png)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/xu2gy9/1616124241956-edea8b9d-7c2c-463d-9206-f508ac7e8241.png)

我们知道，普通的网卡是通过网线来收发数据包的话，而 TUN 设备比较特殊，它通过一个文件收发数据包。

如上图所示，tunX 和上面的 eth0 在逻辑上面是等价的， tunX 也代表了一个网络接口,虽然这个接口是系统通过软件所模拟出来的.

网卡接口 tunX 所代表的虚拟网卡通过文件 /dev/tunX 与我们的应用程序(App) 相连，应用程序每次使用 write 之类的系统调用将数据写入该文件，这些数据会以网络层数据包的形式，通过该虚拟网卡，经由网络接口 tunX 传递给网络协议栈，同时该应用程序也可以通过 read 之类的系统调用，经由文件 /dev/tunX 读取到协议栈向 tunX 传递的所有数据包。

此外，协议栈可以像操纵普通网卡一样来操纵 tunX 所代表的虚拟网卡。比如说，给 tunX 设定 IP 地址，设置路由，总之，在协议栈看来，tunX 所代表的网卡和其他普通的网卡区别不大，当然，硬要说区别，那还是有的,那就是 tunX 设备不存在 MAC 地址，这个很好理解，tunX 只模拟到了网络层，要 MAC 地址没有任何意义。当然，如果是 tapX 的话，在协议栈的眼中，tapX 和真是网卡没有任何区别。

如果我们使用 TUN 设备搭建一个基于 UDP 的 VPN ，那么整个处理过程可能是这幅样子：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/xu2gy9/1616124241974-9a5fa9c2-7881-4919-b735-9ad361d649df.png)

首先，我们的应用程序通过 eth0 和远程的 UDP 程序相连,对方传递过来的 UDP 数据包经由左边的协议栈传递给了应用程序，UDP 数据包的内容其实是一个网络层的数据包，比如说 IP 数据报，应用程序接收到该数据包的数据（剥除了各种头部之后的 UDP 数据）之后，然后进行一定的处理，处理完成后将处理后的数据写入文件 /dev/tunX，这样，数据会第二次到达协议栈。需要注意的是，上图中绘制的两个协议栈其实是同一个协议栈，之所以这么画是为了叙述的方便。

## TAP 设备

TAP 设备与 TUN 设备工作方式完全相同，区别在于：

- TUN 设备是一个三层设备，它只模拟到了 IP 层，即网络层 我们可以通过 /dev/tunX 文件收发 IP 层数据包，它无法与物理网卡做 bridge，但是可以通过三层交换（如 ip_forward）与物理网卡连通。可以使用 ifconfig 之类的命令给该设备设定 IP 地址。
- TAP 设备是一个二层设备，它比 TUN 更加深入，通过 /dev/tapX 文件可以收发 MAC 层数据包，即数据链路层，拥有 MAC 层功能，可以与物理网卡做 bridge，支持 MAC 层广播。同样的，我们也可以通过 ifconfig 之类的命令给该设备设定 IP 地址，你如果愿意，我们可以给它设定 MAC 地址。

最后，关于文章中出现的二层，三层，我这里说明一下，第一层是物理层，第二层是数据链路层，第三层是网络层，第四层是传输层。

参考文章:

\[1]. <https://blog.kghost.info/2013/03/27/linux-network-tun/>
