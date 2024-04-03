---
title: Glossary
linkTitle: Glossary
date: 2024-02-21T11:35
weight: 20
---

# 概述

> 参考：
>
> - [Wiki，Glossary](https://en.wikipedia.org/wiki/Glossary)

**Standardized Glossary(标准化术语)**

学习某项技术时，有些名词，比如某某可以是技术、规范、标准、行为、协议(协议其实从广义角度看也是标准)、等。

# IDC

https://en.wikipedia.org/wiki/Data_center

**Internet data center(互联网数据中心，简称 IDC)**，也可以简称为 Data center(数据中心)，并不用只限制在互联网。IDC 是一座建筑物、建筑物内的专用空间或一组建筑物，用于容纳计算机系统和相关设备。通常用于对外或对内提供 计算、存储、通信 这最基本的三大能力。

# ISP

https://en.wikipedia.org/wiki/Internet_service_provider

**Internet service provider(互联网服务提供商，简称 ISP)** 是提供访问、使用、管理或参与 Internet 服务的组织。 ISP 可以以多种形式组织，例如商业、社区所有、非营利或其他私人所有。比如 中国移动、中国联通、中国电信、etc. 都属于 ISP

# 版本信息

| 英文     | 中文           | 缩写 | 说明                                                                |
| -------- | -------------- | ---- | ------------------------------------------------------------------- |
| Portable | 便携式、可移植 |      | 一个程序如果不需要安装，直接使用二进制文件运行，通常称为 Portable。 |

# 全部

| 英文                                                 | 中文             | 缩写与简称   | 链接                                                                                       | 解释                                                                                                       |
| -------------------------------------------------- | -------------- | ------- | ---------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- |
| 5-tuple                                            | 五元组            |         | [RFC 6146](https://datatracker.ietf.org/doc/html/rfc6146#section-2)                      | 源 IP，源 PORT，目的 IP，目的 PORT，传输层协议，这五个量组成的一个集合                                                              |
| Advanced Telecommunications Computing Architecture | 高级电信计算架构       | ATCA    | [Wiki](https://en.wikipedia.org/wiki/Advanced_Telecommunications_Computing_Architecture) | atca架构本身就是一组工业标准框架，只要是基于这个国际统一标准做的板卡都可以集成到一起                                                             |
| Broadband Remote Access Server                     | 宽带远程接入服务器      | BRAS    | [Wiki](https://en.wikipedia.org/wiki/Broadband_remote_access_server)                     | 是一种用于管理和控制带宽接入用户的网络设备                                                                                    |
| Call detail record                                 | 通话详细记录         | CDR     | [Wiki](https://en.wikipedia.org/wiki/Call_detail_record)                                 | 中文常简称为 "话单"。随着发展，话单的含义也逐步扩展，包含了不止是通话的详细信息。有时候也用 xDR 描述。                                                  |
| Data Plane Development Kit                         | 数据平面开发套件       | DPDK    | [DPDK](/docs/4.数据通信/DPDK.md)                                                             |                                                                                                          |
| Deep packet inspection                             | 深度数据包检测        | DPI     | [DPI](/docs/7.信息安全/Network%20analysis/DPI.md)                                            |                                                                                                          |
| Call Detail Record                                 | 通话详细记录         | CDR(话单) | [CDR](https://en.wikipedia.org/wiki/Call_detail_record)                                  | 后期随着发展该名词逐渐包含了 网络、等 通信之间的详细记录，而不是单指通话。可以写为 **xDR**(wiki 上没有 xDR，自己造的)                                    |
| Cyberspace Situation Awareness                     | 网络态势感知         | CSA     |                                                                                          |                                                                                                          |
| Mellanox Technologies                              |                |         |                                                                                          | 一家以色列裔美国跨国供应商，提供基于 InfiniBand 和以太网技术的计算机网络产品。Mellanox 为高性能计算、数据中心、云计算、计算机数据存储和金融服务                       |
| Remote Authentication Dial-In User Service         | 远程用户拨号认证       | RADIUS  | [Wiki](https://en.wikipedia.org/wiki/RADIUS)                                             |                                                                                                          |
| Situational awareness                              | 态势感知           | SA      | [Wiki](https://en.wikipedia.org/wiki/Situation_awareness)                                |                                                                                                          |
| Switched Port Analyzer                             |                | SPAN    |                                                                                          |                                                                                                          |
| Transaction<br>                                    | 事务             |         | [Transaction](#transaction)                                                              |                                                                                                          |
| tap/terminal access point/test access point        | 窃听/终端接入点/测试接入点 | tap/TAP | [Wiki](https://en.wikipedia.org/wiki/Network_tap#Terminology)                            | tap 并不是什么单词的缩写，仅仅就是原本的窃听之类的意思。只不过有人觉得这词不好，而且也有一些特殊情况在，就通过 backronyms(把单词扩写，不是缩写的楞扩成好几个词) 把 tap 改成了其他的意思。 |

## Transaction

假设某个数据可能需要经过 A、B、C、D 几个步骤才能修改完毕，我们把这四个步骤打包放到事务中，那么事务就可以确保这四个步骤要么全部执行完毕，要么全部都不去执行。这样即使在任意一个步骤断电或者程序崩溃都不会影响到数据的一致性问题。
