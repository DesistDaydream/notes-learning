---
title: OpenWrt
linkTitle: OpenWrt
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，openwrt/openwrt](https://github.com/openwrt/openwrt)
> - [官网](https://openwrt.org/)

OpenWrt 项目是一个针对嵌入式设备的 [Unix-like OS](docs/1.操作系统/Operating%20system/Unix-like%20OS/Unix-like%20OS.md)。与尝试创建单一、静态固件不同，OpenWrt 提供了一个完全可写的文件系统，并配备了软件包管理。这使您摆脱了供应商提供的应用程序选择和配置，并允许您通过使用软件包来定制设备，以适应任何应用程序。对于开发人员来说，OpenWrt 是构建应用程序的框架，无需在其周围构建完整固件；对于用户来说，这意味着完全定制的能力，可以以前所未想象的方式使用设备。

OpenWrt 的包管理器是 OPKG。

据说，爱快(ikuai) 是基于 OpenWrt 的二次封装系统。

ikuai 与 OpenWrt 大部分出现在软路由场景。

iStoreOS 是 koolshare 团队基于OpenWrt定制的软路由系统

# 关联文件与配置


# Release

https://downloads.openwrt.org/

https://downloads.openwrt.org/releases/

TODO: 官方各种 Image 文件没法直接用吗？原因是啥？

## 官方

## eSir

[GitHub 项目，esirplayground/AutoBuild-OpenWrt](https://github.com/esirplayground/AutoBuild-OpenWrt) # 只有 [GitHub Actions](docs/2.编程/Programming%20tools/SCM/GitHub/GitHub%20Actions/GitHub%20Actions.md)

[Telegram, eSir PlayGround固件发布频道](https://t.me/esirplayground)

https://drive.google.com/drive/folders/1MIzj4Hn9hdUZ3K8oksl2Efqs5inrBUQ7 # 各种编译好的 img 文件发布在 Google 云盘中

> .img 文件命中的 gdq = 高大全 ~~~

虚拟机测试做法：

- 将下载的 img 转成 qcow2
    - qemu-img convert -f raw -O qcow2 openwrt-gdq-version-v1-2025-x86-64-generic-squashfs-legacy.img op.qcow2
- 安装时使用 "Import existing disk image"，启动后即可直接进入系统
- 修改一下地址让自己可以访问到，即可直接访问 Web 页面