---
title: Android 开发
linkTitle: Android 开发
date: 2023-06-18T16:39
weight: 3
---

# 概述

> 参考：
>
> - [安卓开发者资源](https://developer.android.com/)

# AOSP

https://source.android.com/

**Android Open Source Project(安卓开源项目，简称 AOSP)** 提供了创建 Android 操作系统的自定义变体、将设备和配件移植到 Android 平台所需的信息和源代码，并确保设备满足兼容性要求，从而保持 Android 生态系统的健康。为数百万用户提供健康稳定的环境。

# Android Studio

**Andriid Studio(安卓工作室)** 是一款用于 Android 应用程序开发的官方 **Integrated Development Environment(集成开发环境，简称 IDE)**，包含代码编辑器、构建工具、签名工具、SDK 工具等。

# SDK 工具

> 参考：
>
> - [官方文档，SDK 工具](https://developer.android.com/tools)

## ADB

详见 [ADB](/docs/Mobile%20device/Android%20MGMT/ADB.md)

# 最佳实践

## 获取 Root 权限

Magisk [GitHub 项目，topjohnwu/Magisk](https://github.com/topjohnwu/Magisk)

https://github.com/topjohnwu/Magisk/ 对应用隐藏 Root 信息

## 其他

使用 Magisk 的模块为系统添加 CA 证书

- https://github.com/NVISOsecurity/MagiskTrustUserCerts 2 年没更新了
- https://github.com/nccgroup/ConscryptTrustUserCerts # 适用于 andriod 14 ？
  - https://github.com/nccgroup/ConscryptTrustUserCerts/issues/3 有了下面那个仓库
- https://github.com/lupohan44/TrustUserCertificates
