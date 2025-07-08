---
title: SNMP(传统监控标准)
linkTitle: SNMP(传统监控标准)
weight: 20
---

# 概述

> 参考：
>
> - [RFC 1157，A Simple Network Management Protocol (SNMP)](https://datatracker.ietf.org/doc/html/rfc1157)
> - [RFC 1156，Management Information Base for Network Management of TCP/IP-based internets](https://datatracker.ietf.org/doc/html/rfc1156)
> - [Wiki, SNMP](https://en.wikipedia.org/wiki/Simple_Network_Management_Protocol)

**Simple Network Management Protocol(简单网络管理协议，简称 SNMP)**。想实现该协议，通常需要由两部分完成(监控端和被监控端)，是在两端的两个进程之间进行通信，该进程都需要占用一个 socket

- 监控端：通常称为 NMS 端，管理端
- 被监控端：通常称为 Agent 端，NMS 要去收集被监控端的数据的时候，可能收集到的是一些很敏感的数据(CPU 使用率，带宽占用率等，这些一般是不公开的)，所以需要在被监控节点上部署一个专门的程序，这个程序能在本地取得一定的管理权限，然后接受监控端发送的数据收集指令，代为在被监控节点本地完成数据收集，所以被称为 Agent 端，代理端

SNMP 的工作模式，使用 udp 协议发送报文

- 监控端主动发送请求到被监控端的 agent 去收集数据
- 被监控节点主动向监控端报告自己所采集的数据
- 当监控端发现被监控端发生异常时，可以发送一些控制指令，将被监控端修改一些参数

# 实现 SNMP 的组件

- **Management Information Base(管理信息库，简称 MIB)** # 用来定义所有监控端的 objects，其中包括 objects 的名称、OID、数据类型、描述(干什么用的)。MIB 也可以看作是 SNMP 的服务端与代理端的沟通桥梁，只有具有统一的格式，才能确定数据。
  - **Object(对象)** # 这个对象可以是一个具体需要采集到的数据，比如 内存、CPU、磁盘、网络接口等等，也可以是一种抽象的集合，比如地区、硬件、系统、硬件、网络等等。上面说的所有事物，每一个都是一个 Object。所以，Object 可以包含另一个 Object，这也是人称常常将 MIB 称为**树状**的原因
    - **Object Identifier(对象标识符，简称 OID)** # 每一个 Object 都有一个 OID
    - 数据存取格式：即每个 object 除了 OID 用作标示以外，还有数据内容需要遵循一定个格式规范
- **Structure of Managerment Intormation(管理信息结构,简称 SMI)** # 是 [ASN.1](/docs/2.编程/无法分类的语言/ASN.1.md) 的子集
- **SNMP 本身** # 一般通过 Net-SNMP 中的工具实现。

所谓的 **MIB**，其实主要是通过文件记录的内容。与其说是用文件记录，不如说 MIB 就是使用 ASN.1(标准的接口描述语言) 编写的代码。ASN.1 语言同样有类似 import、 function 这类的东西。只不过，记录 MIB 文件的语言，又与 ASN.1 有一些细微的区别，我们暂且称为 **MIB 语言** 吧~

可以这么说，**MIB 就是一门描述 OID 的编程语言**。

# SNMP 安全

<https://www.dnsstuff.com/snmp-community-string#what-is-an-snmp-community-string>

## Community String(团体字符串)

SNMP 的 **Community String(团体字符串)** 是一种访问存储在路由器或其他设备中的统计信息的方法。有时简称为社区字符串或 SNMP 字符串，它包含与 GET 请求一起提供的用户凭据（ID 或密码）。

# 协议规范

https://datatracker.ietf.org/doc/html/rfc1157#section-4

所有 SNMP 实现都必须支持这五种 PDU：

1. GetRequest-PDU
2. GetNextRequest-PDU
3. GetResponse-PDU
4. SetRequest-PDU
5. Trap-PDU

