---
title: DMTF
linkTitle: DMTF
date: 2024-01-07T21:58
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, Distributed_Management_Task_Force](https://en.wikipedia.org/wiki/Distributed_Management_Task_Force)

**Distributed Management Task Forc(分布式管理任务组，简称 DMTF)** 是一个 501(c)(6) 非营利行业标准组织，它创建涵盖各种新兴和传统 IT 基础设施（包括云、虚拟化、网络、服务器和存储）的开放可管理性标准。成员公司和联盟伙伴就标准进行协作，以改进信息技术的互操作管理。

DMTF 总部位于俄勒冈州波特兰，由代表科技公司的董事会领导，这些公司包括：Broadcom Inc.、Cisco、Dell Technologies、Hewlett Packard Enterprise、Intel Corporation、Lenovo、NetApp、Positivo Tecnologia S.A. 和 Verizon。

该组织成立于 1992 年，当时名为 **Desktop Management Task Force(桌面管理任务组)**，其第一个标准是现在的 Desktop Management Interface(DMI)。随着该组织发展到通过附加标准（例如通用信息模型 (CIM)）解决分布式管理问题，它于 1999 年更名为分布式管理任务组，但现在称为 DMTF。

# DMI

> 参考：
>
> - [Wiki, Desktop_Management_Interface](https://en.wikipedia.org/wiki/Desktop_Management_Interface)

**Desktop Management Interface(桌面管理接口，简称 DMI)** 是帮助收集电脑系统信息的管理系统，DMI 信息的收集必须在严格遵照 SMBIOS 规范的前提下进行。SMBIOS（System Management BIOS）是主板或系统制造者以标准格式显示产品管理信息所需遵循的统一规范。SMBIOS 和 DMI 是由行业指导机构 Distributed Management Task Force(DMTF) 起草的开放性的技术标准，其中 DMI 设计适用于任何的平台和操作系统。

DMI 充当了管理工具和系统层之间接口的角色。它建立了标准的可管理系统更加方便了电脑厂商和用户对系统的了解。DMI 的主要组成部分是 Management Information Format(MIF) 数据库。这个数据库包括了所有有关电脑系统和配件的信息。通过 DMI，用户可以获取序列号、电脑厂商、串口信息以及其它系统配件信息。

# CIM

> 参考：
>
> - [Wiki, Common_Information_Model](https://en.wikipedia.org/wiki/Common_Information_Model_(computing))

**Common Information Model(通用信息模型，简称 CIM)** 是一个开放标准，定义了 IT 环境中被管理的元素如何表示为一组公共对象及其之间的关系。

> Notes: 这里的对象是指 [计算机科学](/docs/2.编程/计算机科学/计算机科学.md) 中的对象，可以是一个变量、数据结构、函数、方法、等等等等。
