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

| 英文 | 中文 | 缩写 | 链接 | 解释 |
| ---- | ---- | ---- | ---- | ---- |
| Deep packet inspection | 深度数据包检测 | DPI | [DPI](docs/7.信息安全/DPI/DPI.md) |  |
| Transaction<br> | 事务 |  | [Transaction](#Transaction) |  |

## Transaction

假设某个数据可能需要经过A、B、C、D几个步骤才能修改完毕，我们把这四个步骤打包放到事务中，那么事务就可以确保这四个步骤要么全部执行完毕，要么全部都不去执行。这样即使在任意一个步骤断电或者程序崩溃都不会影响到数据的一致性问题。