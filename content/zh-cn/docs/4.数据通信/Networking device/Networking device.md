---
title: Networking device
linkTitle: Networking device
date: 2024-02-22T13:29
weight: 20
---

# 概述

> 参考：
>
> - [Wiki，Networking hardware](https://en.wikipedia.org/wiki/Networking_hardware)

**Networking device(网络设备)** 也成为 Networking hardware(网络硬件) 是计算机网络上的设备之间进行通信和交互所需的电子设备。具体来说，它们调解计算机网络中的数据传输。最后接收或产生数据的单元称为主机、端系统或数据终端设备。

[Router(路由器) 与 Switch(交换机)](/docs/4.数据通信/Networking%20device/Router%20And%20Switch.md) 是常见的基础网络设备。

# 辅助材料

## 以太网双绞线

> 参考：
> - [Wiki，Ethernet over twisted pair](https://en.wikipedia.org/wiki/Ethernet_over_twisted_pair)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/networking_device/202402271413818.png)

## 光纤

> 参考：
>
> - [Wiki，Optical fiber](https://en.wikipedia.org/wiki/Optical_fiber)(光纤)
> - [Wiki，Fiber-optic communication](https://en.wikipedia.org/wiki/Fiber-optic_communication)(光纤通信)

![image.png|300](https://notes-learning.oss-cn-beijing.aliyuncs.com/networking_device/202402271411895.png)

光功率的单位是 dbm，光衰的单位是 db。发送端光功率(大)-接收端光功率（小）=光衰（正直）

光功率值分大小，越小信号越弱。通常发光小于 0dbm(负值)。

接收端能够接收的最小光功率称为灵敏度，发光功率减去接收灵敏度是允许的光纤衰耗值（dbm-dbm=db）.测试时实际的发光功率减去实际接收到的光功率的值就是光纤衰耗(db)

举例说明：一段光纤能接受的最小光功率（即接受灵敏度）是-20dbm(低于-20 光纤点不亮)。

发送端光功率是-5dbm，则最大光衰为-5-（-20）=15db,即允许的最大衰耗为 15db，大于则点不亮。 如接收端测出来的是-10dmb，那么光纤衰耗是-5-（-10）=5db,小于最大衰耗 15db，此时光纤能正常点亮。

光纤电缆，光缆 Fiber-optic cable https://en.wikipedia.org/wiki/Fiber-optic_cable

TODO: Single-mode fiber(单模) 和 Multi-mode fiber(多模)？

- SC 接口
- LC 接口
- FC 接口

## Optical module(光模块)

**光电收发器** 即光模块，光转电，由 光电子器件、功能电路、光接口 组成。光电子器件包括发射和接收两部分。

- 简单的说，光模块的作用就是光电转换，发送端把电信号转换成光信号，通过光纤传送后，接收端再把光信号转换成电信号。

