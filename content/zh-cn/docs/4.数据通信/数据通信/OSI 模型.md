---
title: OSI 模型
linkTitle: OSI 模型
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, OSI model](https://en.wikipedia.org/wiki/OSI_model)

# OSI

**Open System Interconnection Model(开放式系统互联模型，简称 OSI 模型)** 是一种概念模型，由国际标准化组织提出，一个试图使各种计算机在世界范围内互连为网络的标准框架。定义于 ISO/IEC 7498-1。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tudqyb/1616161546435-bc8ba6d0-bb9a-4fc8-819e-6108035f2d1a.jpeg)

# OSI 模型背景

在制定计算机网络标准方面，起着重大作用的两大国际组织是：国际电信联盟电信标准化部门，与国际标准组织（ISO），虽然它们工作领域不同，但随着科学技术的发展，通信与信息处理之间的界限开始变得比较模糊，这也成了国际电信联盟电信标准化部门和 ISO 共同关心的领域。1984 年，ISO 发布了著名的 ISO/IEC 7498 标准，它定义了网络互联的 7 层框架，也就是开放式系统互联参考模型。

层次划分

根据建议 X.200，OSI 将计算机网络体系结构划分为以下七层，标有 1 ～ 7，第 1 层在底部。 现“OSI/RM”是英文“Open Systems Interconnection Reference Model”的缩写。

第 7 层 应用层

主条目：应用层

应用层（Application Layer）提供为应用软件而设的接口，以设置与另一应用软件之间的通信。例如: HTTP，HTTPS，FTP，TELNET，SSH，SMTP，POP3.HTML.等。

第 6 层 表达层

主条目：表达层

表达层（Presentation Layer）把数据转换为能与接收者的系统格式兼容并适合传输的格式。

第 5 层 会话层

主条目：会话层

会话层（Session Layer）负责在数据传输中设置和维护电脑网络中两台电脑之间的通信连接。

第 4 层 传输层

主条目：传输层

传输层（Transport Layer）把传输表头（TH）加至数据以形成数据包。传输表头包含了所使用的协议等发送信息。例如:传输控制协议（TCP）等。

第 3 层 网络层

主条目：网络层

网络层（Network Layer）决定数据的路径选择和转寄，将网络表头（NH）加至数据包，以形成报文。网络表头包含了网络数据。例如:互联网协议（IP）等。

第 2 层 数据链路层

主条目：数据链路层

数据链路层（Data Link Layer）负责网络寻址、错误侦测和改错。当表头和表尾被加至数据包时，会形成帧。数据链表头（DLH）是包含了物理地址和错误侦测及改错的方法。数据链表尾（DLT）是一串指示数据包末端的字符串。例如以太网、无线局域网（Wi-Fi）和通用分组无线服务（GPRS）等。

分为两个子层：逻辑链路控制（logical link control，LLC）子层和介质访问控制（Medium access control，MAC）子层。

第 1 层 物理层

主条目：物理层

物理层（Physical Layer）在局部局域网上传送数据帧（data frame），它负责管理电脑通信设备和网络媒体之间的互通。包括了针脚、电压、线缆规范、集线器、中继器、网卡、主机接口卡等。

| OSI 七层网络模型        | TCP/IP 四层概念模型 | 对应网络协议                                             |
| ----------------- | ------------- | -------------------------------------------------- |
| 应用层（Application）  | 应用层           | HTTP、TFTP, FTP, NFS, WAIS、SMTP                     |
| 表示层（Presentation） |               | Telnet, Rlogin, SNMP, Gopher                       |
| 会话层（Session）      |               | SMTP, DNS                                          |
| 传输层（Transport）    | 传输层           | TCP, UDP                                           |
| 网络层（Network）      | 网络层           | IP, ICMP, ARP, RARP, AKP, UUCP                     |
| 数据链路层（Data Link）  | 数据链路层         | FDDI, Ethernet, Arpanet, PDN, SLIP, PPP、VLAN、VxLAN |
| 物理层（Physical）     |               | IEEE 802.1A, IEEE 802.2 到 IEEE 802.11              |
