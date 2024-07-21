---
title: Android MGMT
linkTitle: Android MGMT
date: 2024-07-21T11:42
weight: 20
---

# 概述

> 参考：
>
> -

# ADB

[ADB](/docs/Mobile%20device/Android%20MGMT/ADB.md)

# Scrcpy

> 参考:
>
> - [GitHub 项目，Genymobile/scrcpy](https://github.com/Genymobile/scrcpy)

Scrcpy 可以显示和控制 Android 设备。

该应用程序镜像通过 USB 或 TCP/IP 连接的 Android 设备（视频和音频），并允许使用计算机的键盘和鼠标控制设备。它不需要任何根访问权限。它适用于 Linux、Windows 和 macOS。

> [!Tip]
>
> 使用 Scrcpy 连接运行在 Linux 容器中的 Android 系统，配合 [Yolo](https://github.com/ultralytics/yolov5)。可以实现在电脑上自动玩手机游戏，比如 https://www.bilibili.com/video/BV1mihkeTEDc

Android 容器镜像

- https://hub.docker.com/r/redroid/redroid
- https://github.com/budtmo/docker-android 好像跑不起来？