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
