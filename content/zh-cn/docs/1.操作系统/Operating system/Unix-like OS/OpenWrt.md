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

OpenWrt 项目是一个针对嵌入式设备的 [Unix-like OS](/docs/1.操作系统/Operating%20system/Unix-like%20OS/Unix-like%20OS.md)。与尝试创建单一、静态固件不同，OpenWrt 提供了一个完全可写的文件系统，并配备了软件包管理。这使您摆脱了供应商提供的应用程序选择和配置，并允许您通过使用软件包来定制设备，以适应任何应用程序。对于开发人员来说，OpenWrt 是构建应用程序的框架，无需在其周围构建完整固件；对于用户来说，这意味着完全定制的能力，可以以前所未想象的方式使用设备。

OpenWrt 的包管理器是 OPKG。

据说，爱快(ikuai) 是基于 OpenWrt 的二次封装系统。

ikuai 与 OpenWrt 大部分出现在软路由场景。

iStoreOS 是 koolshare 团队基于OpenWrt定制的软路由系统

OpenWrt 生态项目

- [luci](https://github.com/openwrt/luci) # [Lua](/docs/2.编程/高级编程语言/Lua/Lua.md) 开发的配置接口。名称是 lua +  uci 的组合。
- [uhttpd](https://github.com/openwrt/uhttpd) # HTTP 服务器。通常与 LuCI 组合使用
- etc.

# 关联文件与配置

# 安装

与 [安装操作系统](/docs/1.操作系统/安装操作系统/安装操作系统.md) 的逻辑类似，大体分如下几步

- 下载 Release
- 将 Release 制作到 U 盘中
- 在目标机器上插入 U 盘并写入 OpenWrt 系统

# Release

## 官方

https://downloads.openwrt.org/

https://downloads.openwrt.org/releases/

### 简单的虚拟机测试与体验

使用 [openwrt-24.10.1-x86-64-generic-ext4-combined-efi.img.gz](https://downloads.openwrt.org/releases/24.10.1/targets/x86/64/openwrt-24.10.1-x86-64-generic-ext4-combined-efi.img.gz) Release 作为示例（带有 LuCI 的 Web 页面），下载并解压

将解压后的 img 转成 qcow2

```bash
qemu-img convert -f raw -O qcow2 openwrt-24.10.0-x86-64-generic-ext4-combined-efi.img openwrt-original.qcow2
```

安装时直接导入转换后的 qcow2 文件即可

```bash
virt-install --name openwrt-original \
--memory 2048 --vcpus 1 \
--os-variant linux2024 \
--disk /var/lib/libvirt/images/openwrt-original.qcow2,bus=virtio \
--network bridge=br0,model=virtio \
--graphics vnc,listen=0.0.0.0,port=5911 \
--noautoconsole \
--import
```

默认启动了虚拟机的 Console 可直接进入

> 一般出现 “Please press Enter to activate this console.” 即表示正常

```
]# virsh console openwrt-original
Connected to domain 'openwrt-original'
Escape character is ^] (Ctrl + ])



BusyBox v1.36.1 (2025-02-03 23:09:37 UTC) built-in shell (ash)

  _______                     ________        __
 |       |.-----.-----.-----.|  |  |  |.----.|  |_
 |   -   ||  _  |  -__|     ||  |  |  ||   _||   _|
 |_______||   __|_____|__|__||________||__|  |____|
          |__| W I R E L E S S   F R E E D O M
 -----------------------------------------------------
 OpenWrt 24.10.0, r28427-6df0e3d02a
 -----------------------------------------------------
=== WARNING! =====================================
There is no root password defined on this device!
Use the "passwd" command to set up a new password
in order to prevent unauthorized SSH logins.
--------------------------------------------------
root@OpenWrt:~#
```

修改一下地址让自己可以访问到，即可直接访问 OpenWRT 的 LuCI 图形界面了，效果如下

![800](https://notes-learning.oss-cn-beijing.aliyuncs.com/openwrt/demo_luci.png)

## eSir

[GitHub 项目，esirplayground/AutoBuild-OpenWrt](https://github.com/esirplayground/AutoBuild-OpenWrt) # 只有 [GitHub Actions](/docs/2.编程/Programming%20tools/SCM/GitHub/GitHub%20Actions/GitHub%20Actions.md)

[Telegram, eSir PlayGround固件发布频道](https://t.me/esirplayground)

https://drive.google.com/drive/folders/1MIzj4Hn9hdUZ3K8oksl2Efqs5inrBUQ7 # 各种编译好的 img 文件发布在 Google 云盘中

> .img 文件命中的 gdq = 高大全 ~~~ 表示这个镜像包含的内容多，占用空间大。。。。

虚拟机测试做法：

- 将下载的 img 转成 qcow2
  - qemu-img convert -f raw -O qcow2 openwrt-gdq-version-v1-2025-x86-64-generic-squashfs-legacy.img op.qcow2
- 安装时使用 "Import existing disk image"，启动后即可直接进入系统
- 修改一下地址让自己可以访问到，即可直接访问 Web 页面
