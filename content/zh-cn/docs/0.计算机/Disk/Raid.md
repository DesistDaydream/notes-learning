---
title: Raid
---

# 概述

> 参考：
>
> - [Wiki, Raid](https://en.wikipedia.org/wiki/RAID)
> - [Wiki 中文，Raid](https://zh.wikipedia.org/wiki/RAID)
> - [Wiki, Erasure code](https://en.wikipedia.org/wiki/Erasure_code)(纠删码)
> - [Wiki, Parity bit](https://en.wikipedia.org/wiki/Parity_bit)(奇偶校验)
> - <https://support.huawei.com/enterprise/zh/doc/EDOC1000163568/26751928>

**Redundant array of independent disks(独立磁盘冗余阵列，简称 RAID)** 是一种存储虚拟化技术，可以将多个 [Disk](/docs/0.计算机/Disk/Disk.md)(物理磁盘) 组合成一个或多个 Logical unit(逻辑单元) 以达到数据冗余、提高性能或两者兼得的目的。

数据以多种方式之一分布在驱动器上，称为 RAID 级别，具体取决于所需的冗余和性能级别。不同的方案或数据分布布局由单词“RAID”命名，后跟一个数字，例如 RAID 0 或 RAID 1。每个方案或 RAID 级别在关键目标之间提供不同的平衡：可靠性、可用性、性能和容量。高于 RAID 0 的 RAID 级别可针对不可恢复的扇区读取错误以及整个物理驱动器的故障提供保护。

Raid 5 等的奇偶校验机制，是纠删码的最佳实践

# 虚拟磁盘的读写策略

在创建虚拟磁盘时，会需要对其数据读写策略进行定义，以规范后续虚拟磁盘运行过程中数据的读写方式。

## 数据读策略

在配置界面中一般体现为“Read Policy”。RAID 卡支持如下两种数据读策略：

- **Read-ahead(预读取)** # 在配置界面中一般有“Always Read Ahead”、“Read Ahead”、“Ahead”等配置选项。使用此策略后，从虚拟磁盘中读取所需数据时，会把后续数据同时读出放在 [Cache](/docs/8.通用技术/Cache.md) 中，用户随后访问这些数据时可以直接在 Cache 中命中，将减少磁盘寻道操作，节省响应时间，提高了数据读取速度。要使用该策略，要求 RAID 控制卡支持数据掉电保护功能，但如果此时超级电容异常，可能导致数据丢失。
- **No-Read-Ahead(非预读取)** # 使用此策略后，RAID 卡接收到数据读取命令时，才从虚拟磁盘读取数据，不会做预读取的操作。

## 数据写策略

在配置界面中一般体现为“Write Policy”。RAID 卡支持如下三种数据写策略：

- **Write Back(回写)** # 在配置界面中一般体现为“Write Back”等字样。使用此策略后，需要向虚拟磁盘写数据时，会直接写入 Cache 中，当写入的数据积累到一定程度，RAID 卡才将数据刷新到虚拟磁盘，这样不但实现了批量写入，而且提升了数据写入的速度。当控制器 Cache 收到所有的传输数据后，将给主机返回数据传输完成信号。要使用该策略，要求 RAID 卡支持数据掉电保护功能，但如果此时超级电容异常，可能导致数据丢失。
- **Write Through(直写)** # 在配置界面中一般有“Write Through”等选项。使用此策略后，RAID 卡向虚拟磁盘直接写入数据，不经过 Cache。当磁盘子系统接收到所有传输数据后，控制器将给主机返回数据传输完成信号。此种方式不要求 RAID 卡支持数据掉电保护功能，即使超级电容故障，也无影响。该写策略的缺点是写入速度较低。。
- 与 BBU 相关的回写 # 在配置界面中一般有“Write Back with BBU”等选项。使用此策略后，当 RAID 卡 BBU 在位且状态正常时，RAID 卡到虚拟磁盘的写操作会经过 Cache 中转（即回写方式）；当 RAID 卡 BBU 不在位或 BBU 故障时，RAID 卡到虚拟磁盘的写操作会自动切换为不经过 Cache 的直接写入（即写通方式）。

# 阵列卡

## JBOD

> 参考：
>
> - [Wiki, Non-RAID_drive_architectures](https://en.wikipedia.org/wiki/Non-RAID_drive_architectures)

**Just a Bunch Of Disks(只是一堆磁盘，简称 JBOD)** 是阵列卡可以提供的一种透传模式，让硬盘可以直接透过阵列卡对系统直接提供本身的存储能力。还可以将多个物理硬盘组成单个逻辑磁盘，但是**这并不是一种 RAID**，也不提供数据冗余，仅仅是将这些硬盘从头到尾串联起来而已。

JBOD 在很多厂家对外的口语描述中，称为**直通模式**。意味着硬盘直通到系统，中间不做任何处理。

## H3C 阵列配置

[H460&P460&P2404&P4408 系列](http://www.h3c.com/cn/d_202201/1526857_30005_0.htm#_Toc92721209) 阵列卡
