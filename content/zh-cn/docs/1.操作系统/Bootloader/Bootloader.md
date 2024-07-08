---
title: Bootloader
linkTitle: Bootloader
date: 2024-03-15T20:59
weight: 1
---

# 概述

> 参考：
>
> - [Wiki，Bootloader](https://en.wikipedia.org/wiki/Bootloader)

**Bootloader(引导加载程序)** 是负责 [booting(引导)](https://en.wikipedia.org/wiki/Booting) 计算器的软件。通常也被称为 Bootstrap Loader、Bootstrap。

当计算机关闭时，操作系统、应用程序代码和数据 ‍‌ 仍存储在非易失性存储器中。当计算机开机时，它通常没有操作系统或其随机存取存储器 (RAM) 中的加载程序。计算机首先执行存储在只读存储器（ROM，以及后来的 EEPROM、NOR 闪存）中的相对较小的程序以及一些需要的数据，以初始化 RAM（特别是在 x86 系统上），访问非易失性设备（通常是块设备，例如 NAND 闪存）或可以将操作系统程序和数据加载到 RAM 中的设备。

# 关联文件与配置

**/boot/** # 所有关于系统引导启动的配置信息，都在该目录下

**/boot/grub2/** #

**/etc/default/grub** # TODO: 好像不同系统路径不同？这是啥？

# 引导管理命令行工具

grub2-\*

grubby

grub2-mkconfig
