---
title: ADB
linkTitle: ADB
date: 2024-01-13T15:19
weight: 4
---

# 概述

> 参考：
>
> - [官方文档，SDK 工具-adb](https://developer.android.com/tools/adb)
>   - [官方文档-中文](https://developer.android.com/studio/command-line/adb?hl=zh-cn)
> - [Wiki，Android_Debug_Bridge](https://en.wikipedia.org/wiki/Android_Debug_Bridge)

**Android Debug Bridge(安卓调试桥，简称 ADB)** 是一种功能多样的命令行工具，可让您与设备进行通信。`adb` 命令可用于执行各种设备操作，例如安装和调试应用。`adb` 提供对 Unix shell（可用来在设备上运行各种命令）的访问权限。它是一种 C/S 架构程序，包括以下三个组件：

- **adb** 命令行工具，在开发机器上运行
  - **客户端**：用于发送命令。客户端在开发机器上运行。您可以通过发出 `adb` 命令从命令行终端调用客户端。
  - **服务端**：用于管理客户端与守护程序之间的通信。服务端在开发机器上作为后台进程运行。默认监听 5037 端口
- **adbd** 守护程序，在设备上运行。守护程序在每个设备上作为后台进程运行，以接收 adb 服务端发来的各种指令。

adb 命令行工具作为客户端运行时，会先检查是否有服务端在运行，如果没有，则会执行 `adb -L tcp:5037 fork-server server --reply-fd 4` 命令以启动 adb 服务端，默认监听在 5037 端口，并接收 adb 客户端发出的命令。

而 adbd 守护程序，通常是设备上默认自带的进程，但是默认并没有启动，如果想要启用设备上的 adbd，需要开启 **USB 调试**，该功能通常存在于**开发者选项**中，参考[这里](https://developer.android.com/studio/debug/dev-options?hl=zh-cn#enable)来启用开发者选项。

adb 服务端启动后会自动发现 USB 连接的设备、 安卓 Studio 模拟的设备，然后通过 `adb devices -l` 可以列出这些设备。

### ADB 连接设备

https://blog.51cto.com/u_15549234/5139197

我们可以通过两种方式让 adb 连接到设备

- 本地连接
- 无线连接

本利连接一般是通过 USB 连接真实设备或连接本地 安卓 Studio 模拟的设备。打开设备的 USB 调试并插上线，一般电脑都会自动发现设备。

无线连接则可以通过 Wi-Fi 连接到设备。通常手机的开发者选项中，有一个 **无线调试** 能力，启用后会让 adbd 监听在某个 IP:PORT 上，然后使用 `adb connect IP:PORT` 即可通过网络连接到设备。

### ADB 使用

adb 服务端连接到设备之后，就可以使用 adb <font color="#ff0000">客户端向设备发送命令了</font>。最简单直接的方式是使用 `adb shell` 命令打开一个 Shell，这就像使用 [OpenSSH](/docs/1.操作系统/5.登录%20Linux%20与%20访问控制/Secure%20Shell(SSH)%20安全外壳协议/OpenSSH.md) 类似，可以通过一个 Shell 访问 Android 系统。

> 获取 Root 权限的设备（比如通过 Magisk 授予 Shell 超级用户权限）可以通过 `su` 命令直接切换到 root 用户。

# Syntax(语法)

> 参考：
> 
> - [Manual(手册)，adb(1)](https://android.googlesource.com/platform/packages/modules/adb/+/refs/heads/master/docs/user/adb.1.md)

**adb [GLOBAL_OPTIONS] command [COMMAND_OPTIONS]**

## GLOBAL_OPTIONS

- **-a** # 监听所有网络接口，而不仅仅是本地主机
- **-s SERIAL** # 指定要操作的设备（覆盖 $ANDROID_SERIAL），后续所有子命令都会对该选项指定的 device 进行控制。

# 通用命令

## devices - 列出已连接的设备

# 网络命令

## connect - 通过 TCP/IP 连接到设备

## disconnect - 断开通过 TCP/IP 连接的设备

### Syntax(语法)

disconnect [HOST[:PORT]]

# 文件传输命令

  - push
  - pull
  - sync

# Shell 命令

  - shell
  
# 应用安装命令

  - install # 将单个包推送到设备并安装

# 内部调试命令


# 最佳实践

代理配置

```bash
# 进入 shell
adb shell
# 查看代理
settings get global http_proxy
# 设置代理
settings put global http_proxy 192.168.254.254:8888
# 取消代理
settings put global http_proxy :0
# 或者注意删除代理信息
settings delete global http_proxy
settings delete global global_http_proxy_host
settings delete global global_http_proxy_port
```

