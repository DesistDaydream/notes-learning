---
title: "BIOS"
linkTitle: "BIOS"
weight: 20
---

# 概述

> 参考：
> 
> - [Wiki，BIOS](https://en.wikipedia.org/wiki/BIOS)

**Basic Input/Output System(基本输出输出系统，简称 BIOS)** 是

# UEFI

> 参考：
> - [Wiki，UEFI](https://en.wikipedia.org/wiki/UEFI)

**Unified Extensible Firmware Interface(统一可扩展接口，简称 UEFI)** 用以替代 BIOS

# SMBIOS

> 参考：
> - [Wiki，System_Management_BIOS](https://en.wikipedia.org/wiki/System_Management_BIOS)

**System Management BIOS (系统管理 BIOSS，简称 MBIOS)** 是一种规范，该规范定义了一计算机的 BIOS 中产生的管理信息的**数据结构**以及读取这些数据的**访问方式**。这消除了操作系统直接探测硬件以发现计算机中存在哪些设备的需要。SMBIOS 规范由非营利性标准开发组织分布式管理任务组 ([DMTF](https://en.wikipedia.org/wiki/Distributed_Management_Task_Force)) 制定。

SMBIOS 最初被称为桌面管理 BIOS (DMIBIOS)，因为它与桌面管理界面 (DMI) 交互。

# DMI

> 参考：
> - [Wiki，Desktop_Management_Interface](https://en.wikipedia.org/wiki/Desktop_Management_Interface)

桌面管理接口 (DMI) 通过从管理它们的软件中抽象出这些组件，生成用于管理和跟踪台式机、笔记本电脑或服务器计算机中的组件的标准框架。 1998 年 6 月 24 日开发的 DMI 2.0 版标志着分布式管理任务组 (DMTF) 向桌面管理标准迈出的第一步。在引入 DMI 之前，没有标准化的信息来源可以提供有关个人计算机组件的详细信息。

从 1999 年开始，Microsoft 要求 OEM 和 BIOS 供应商支持 DMI 接口/数据集才能获得 Microsoft 认证

DMI 将系统数据（包括系统管理 BIOS (SMBIOS) 数据）公开给管理软件，但两种规范独立运行。

DMI 通常与 SMBIOS 混淆，SMBIOS 在其第一次修订时实际上称为 DMIBIOS。