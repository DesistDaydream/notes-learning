---
title: Android
---

# 概述

> 参考：
> 
> - [官网](https://www.android.com/)

Android 是一种 [类 Unix 操作系统](/docs/1.操作系统/0.操作系统/类%20Unix%20操作系统/类%20Unix%20操作系统.md)，主要设计用于触摸屏移动设备，如智能手机和平板电脑。Android由一个名为开放手持设备联盟的开发者联盟开发，但其最广泛使用的版本主要由Google开发。它于2007年11月公布，第一款商用Android设备HTC Dream于2008年9月发布。

**Device(设备)** 通常指 手机、平板、手表、等等，甚至可以是安卓 Studio 模拟的设备。

# 目录结构

安卓的目录结构与 [Linux 内核的目录结构](/docs/1.操作系统/2.Kernel(内核)/6.File%20System%20 管理/FHS(文件系统层次标准).md System 管理/FHS(文件系统层次标准).md)类似，但是有一些约定俗成的用于保存各类数据的目录

## /data # ？

- /data/app/ # ?
- /data/data/ # ?

## /sdcard # 与 /storage/emulated/0 目录一样，这俩谁是谁的软链接？

## /storage/emulated/0/

- .**/Android/data/${应用的包名}/** # 应用的缓存和临时目录？

旧版本微信、QQ 接收文件存储路径为存储根目录的 Tencent 目录下，而新版微信接收文件路径切换到了 `Android/data/com.tencent.mm/MicroMsg/Download` 目录下

QQ 接收的文件也切换到了 `Android/data/com.tencent.mobileqq/Tencent/QQfile_recv` 目录下