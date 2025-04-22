---
title: Network Socket
linkTitle: Network Socket
weight: 20
tags:
  - socket
  - IPC
---

# 概述

> 参考：
>
> - [Wiki, Network Scoket](https://en.wikipedia.org/wiki/Network_socket)

Network Socket(网络套接字) 是网络域的 [Socket](/docs/1.操作系统/Kernel/Process/Inter%20Process%20Communication/Socket/Socket.md)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zvw5dg/1616161399168-31d4bf21-49d1-45dc-993e-83ea35d7a7f2.jpeg)

面向连接服务（TCP 协议）：是电话系统服务模式的抽象，即每一次完整的数据传输都要经过建立连接，使用连接，终止连接的过程。在数据传输过程中，各数据分组不携带目的地址，而使用连接号（connect ID）。本质上，连接是一个管道，收发数据不但顺序一致，而且内容相同。TCP 协议提供面向连接的虚电路。例如

1. 文件传输 FTP
2. 远程登录 SSH
3. 数字语音
4. 等

无连接服务（UDP 协议）：是邮政系统服务的抽象，每个分组都携带完整的目的地址，各分组在系统中独立传送。无连接服务不能保证分组的先后顺序，不进行分组出错的恢复与重传，不保证传输的可靠性。UDP 协议提供无连接的数据报服务。例如：

1. 电子邮件
2. 电子邮件中的挂号信
3. 网络数据库查询
4. 等

TCP/IP 的 Socket 提供下列三种类型套接字。

1. 流式套接字（SOCK_STREAM）：TCP Socket
2. 提供了一个面向连接、可靠的数据传输服务，数据无差错、无重复地发送，且按发送顺序接收。内设流量控制，避免数据流超限；数据被看作是字节流，无长度限制。文件传送协议（FTP）即使用流式套接字。
3. 数据报式套接字（SOCK_DGRAM）：UDP Socket
4. 提供了一个无连接服务（UDP）。数据包以独立包形式被发送，不提供无错保证数据可能丢失或重复，并且接收顺序混乱。网络文件系统（NFS）使用数据报式套接字。
5. 原始式套接字（SOCK_RAW） ：裸 Socket
6. 从应用层直接封装网络层报文，跳过传输层的协议.该接口允许对较低层协议，如 IP、ICMP 直接访问。常用于检验新的协议实现或访问现有服务中配置的新设备。

Socket Domain 套接字域，根据其所有使用的地址进行分类

1. AF_INET：Address Family，IPv4
2. AF_INET6：同上，IPv6
3. AF_UNIX：同一主机上不同进程之间通信时使用

每类套接字至少提供了两种 socket：流，数据报

[TCP](/docs/4.数据通信/Protocol/TCP_IP/TCP/TCP.md)：传输控制协议，面向连接的协议，通信钱需要建立虚拟链路，结束后拆除链路

[UDP](/docs/4.数据通信/Protocol/UDP/UDP.md)，无连接的协议

1. 流：可靠的传递，面向连接，无边界
2. 数据报：不可靠的传递，有边界，无连接

## Socket 通信流程

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zvw5dg/1616161399174-a6270b36-8bb5-48a4-ba3c-1ee812d450fe.png)
