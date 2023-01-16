---
title: VNC
---

# 概述

> 参考：
> - [Wiki,VNC](https://en.wikipedia.org/wiki/Virtual_Network_Computing)
> - [GitHub 项目，rfbproto/rfbproto](https://github.com/rfbproto/rfbproto)(RFB 协议规范)

**Virtual Network Computing(虚拟网络计算，简称 VNC)** 是一种图形桌面共享系统，VNC 使用 **RFB Protocol(远程帧缓冲协议)** 以控制另一台计算机；它将键盘和鼠标输入从一台计算机传输到另一台计算机，通过网络中继更新显示器上的信息。

VNC 是独立于平台的——有许多基于 GUI 的操作系统和 Java 的客户端和服务器。多个客户端可以同时连接到 VNC 服务器。该技术的流行用途包括远程技术支持和从家庭计算机访问工作计算机上的文件，反之亦然。

为什么没法通过 VNC 连接物理机的 CLI，但是可以连接虚拟机的 CLI 呢？！在 [rfbproto/rfbproto 的 #18 issue](https://github.com/rfbproto/rfbproto/issues/18) 中倒是说了 QEMU 中内置了 VNC/RFB，并对该协议进行了一些修改以支持一些额外的功能

# TigerVNC

> 参考：
> - [GitHub 项目，TigerVNC/tigernvc](https://github.com/TigerVNC/tigervnc)
> - [官网](https://tigervnc.org/)

TigerVNC 是 VNC 的高性能、平台无关的实现，是一个 C/S 架构应用程序，允许用户在远程机器上启动图形应用程序并与之交互。

## 安装 TigerVNC

CentOS

```bash
yum install tigervnc-server
```

Ubuntu

# RealVNC

> 参考：
> - [官网](https://www.realvnc.com/en/)
