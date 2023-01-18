---
title: Android 操作系统
---

# 概述

> ## 参考：

旧版本微信、QQ 接收文件存储路径为存储根目录的 Tencent 目录下，而新版微信接收文件路径切换到了 `Android/data/com.tencent.mm/MicroMsg/Download` 目录下
QQ 接收的文件也切换到了 `Android/data/com.tencent.mobileqq/Tencent/QQfile_recv` 目录下

# 目录结构

安卓的目录结构与 [Linux 内核的目录结构](/docs/IT学习笔记/1.操作系统/2.Kernel(内核)/6.File%20System%20 管理/FHS(文件系统层次标准).md System 管理/FHS(文件系统层次标准).md)类似，但是有一些约定俗成的用于保存各类数据的目录

## /data # ？

- /data/app/ # ?
- /data/data/ # ?

## /sdcard # 与 /storage/emulated/0 目录一样，这俩谁是谁的软链接？

## /storage/emulated/0/

- .**/Android/data/${应用的包名}/** # 应用的缓存和临时目录？
