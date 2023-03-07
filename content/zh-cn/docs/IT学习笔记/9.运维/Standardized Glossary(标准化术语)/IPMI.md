---
title: IPMI
---

# 概述

> 参考：
> [Wiki，IPMI](https://en.wikipedia.org/wiki/Intelligent_Platform_Management_Interface)
> [Wiki，BMC](https://en.wikipedia.org/wiki/Intelligent_Platform_Management_Interface#Baseboard_management_controller)

**IntelligentPlatformManagement Interface(智能平台管理接口，简称 IPMI) **是一组自动计算机子系统的计算机接口规范，可提供管理和监视功能，独立于主机系统的 CPU，固件（BIOS 或 UEFI）和操作系统。

IPMI 定义了系统管理员使用的一组接口，用于计算机系统的带外管理和监控其操作。例如，IPMI 提供了一种方法来管理可以通过使用与硬件的网络连接而不是对操作系统或登录外部关闭或以其他方式无响应的方式。另一个用例可以远程安装自定义操作系统。如果没有 IPMI，安装自定义操作系统可能需要管理员在计算机附近物理存在，请插入 DVD 或包含 OS 安装程序的 USB 闪存驱动器，并使用监视器和键盘完成安装过程。使用 IPMI，管理员可以安装 ISO 映像，模拟安装程序 DVD，并远程执行安装。

## BMC

**Baseboard Management Controller(主板管理控制器)** 提供 IPMI 架构中的智能。它是一个专门的微控制器，嵌入计算机主板上 - 通常是服务器。 BMC 管理系统管理软件和平台硬件之间的接口。 BMC 有自己的固件和 RAM。

计算机系统内置的不同类型的传感器对 BMC 的参数，如温度，冷却风扇速度，电源状态，操作系统（OS）状态等。 BMC 监视传感器，如果任何参数在预设限制内，则可以通过网络向系统管理员发送警报，指示系统的潜在故障。管理员还可以远程与 BMC 通信，采取一些纠正措施 - 例如重置或电源循环系统以获得再次运行的挂起操作系统。这些能力降低了系统的总体拥有成本。

符合 IPMI 版本 2.0 的系统也可以通过串行通信 LAN，从而可以通过 LAN 远程查看串行控制台输出。实现 IPMI 2.0 的系统通常还包括 KVM OVER IP，远程虚拟媒体和带外嵌入式 Web 服务器界面功能，虽然严格来说，这些位于 IPMI 接口标准的范围之外。

## IMPI 的实现

Dell iDrac
HuaWei iBmc
H3C HDM
浪潮 未知
HP iLO
