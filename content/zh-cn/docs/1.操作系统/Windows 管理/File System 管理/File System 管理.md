---
title: "File System 管理"
linkTitle: "File System 管理"
weight: 1
---

# 概述

> 参考：
>
> - <https://www.howtogeek.com/318177/what-is-the-appdata-folder-in-windows/>

# APPDATA

**${USERPROFILE}/AppData/** #

- **./Local/** #
- **./LocalLow/** #
- **./Roaming/** #

刚装完的 win10 专业版，用户的 AppData 中将会有如下结构：

```powershell
$USERPROFILE/AppData/Local
$USERPROFILE/AppData/LocalLow
$USERPROFILE/AppData/Roaming
$USERPROFILE/AppData/Local/Comms
$USERPROFILE/AppData/Local/ConnectedDevicesPlatform
$USERPROFILE/AppData/Local/D3DSCache
$USERPROFILE/AppData/Local/Microsoft
$USERPROFILE/AppData/Local/Packages
$USERPROFILE/AppData/Local/Temp
$USERPROFILE/AppData/LocalLow/Microsoft
$USERPROFILE/AppData/LocalLow/MSLiveStickerWhiteList
$USERPROFILE/AppData/Roaming/Adobe
$USERPROFILE/AppData/Roaming/Microsoft
```

最主要的是这三个目录下的 Microsfot 目录，还有 Packages 目录。在整理 AppData 时，不要误删了。

## LocalAppData

`${LocalAppData}` 在 `${UserProfile}/AppData/Local/`

**IconCache.db** # 图标缓存数据。若图标变白板，可删除缓存，并重启资源管理。

# Program Files

该目录存储安装在计算机上的大多数应用程序的执行文件。

# Program Files(x86)

该目录存储在 64 位 Windows 系统上安装的 32 位应用程序的执行文件。

# ProgramData

该目录存储全局数据，包括应用程序的配置文件，以及系统的安装和更新信息。

# Users

该目录存储在 Windows 系统上创建的每个用户的个人文件夹，如桌面、文档和图片。

# Windows

Windows 操作系统的核心文件和 DLL 文件都存储在此目录中。

**./System32/** # 类似于 Linux 中的 /usr/sbin/ 目录，系统自带的命令、服务、[msc](/docs/1.操作系统/Windows%20管理/Microsoft%20Management%20Console/Microsoft%20Management%20Console.md#MSC) 的可执行文件都在这里。

# 资源管理器

资源管理器左侧的 “导航窗格” 有时候会加载缓慢，有可能是因为快速访问里有非本地的链接，将这些非本地链接取消固定试试。比如 把 WSL 中的目录固定到快速访问、某种网络文件系统固定到快速访问、etc.