---
title: "File System 管理"
linkTitle: "File System 管理"
weight: 1
---

# 概述
> 参考：
> - <https://www.howtogeek.com/318177/what-is-the-appdata-folder-in-windows/>



**%USERPROFILE%/AppData/\*** #
- **./Local/\*** #
- **./LocalLow/\*** #
- **./Roaming/\*** #

刚装完的 win10 专业版，用户的 AppData 中将会有如下结构：
```powershell
%USERPROFILE%\appdata\Local
%USERPROFILE%\appdata\LocalLow
%USERPROFILE%\appdata\Roaming
%USERPROFILE%\appdata\Local\Comms
%USERPROFILE%\appdata\Local\ConnectedDevicesPlatform
%USERPROFILE%\appdata\Local\D3DSCache
%USERPROFILE%\appdata\Local\Microsoft
%USERPROFILE%\appdata\Local\Packages
%USERPROFILE%\appdata\Local\Temp
%USERPROFILE%\appdata\LocalLow\Microsoft
%USERPROFILE%\appdata\LocalLow\MSLiveStickerWhiteList
%USERPROFILE%\appdata\Roaming\Adobe
%USERPROFILE%\appdata\Roaming\Microsoft
```
最主要的是这三个目录下的 Microsfot 目录，还有 Packages 目录。在整理 AppData 时，不要误删了。

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

**./System32/** # 类似于 Linux 中的 /usr/sbin/ 目录，系统自带的命令、服务、msc 的可执行文件都在这里。