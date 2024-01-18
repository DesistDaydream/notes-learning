---
title: WSA
linkTitle: WSA
date: 2023-11-29T09:39
weight: 20
---

# 概述

> 参考：
> 
> - [GitHub 项目，microsoft/WSA](https://github.com/microsoft/WSA)
> - [官方文档，windows-android-wsa](https://learn.microsoft.com/zh-cn/windows/android/wsa/)

Windows Subsystem for Android with Amazon Appsotre

WSA 管理器若没有打开任何应用、文件管理等功能，过一会会自动关闭 WSA，此时 adb 工具也连接不上，要想保持连接，至少要开着一个 WSA 系统中的功能。

# 关联文件与配置

`%LOCALAPPDATA%/Packages/MicrosoftCorporationII.WindowsSubsystemForAndroid_8wekyb3d8bbwe/` # 安装目录？数据保存目录？

在浏览器或资源管理器的导航栏中输入 `wsa://com.android.settings` 即可打开像手机设置一样的 WSA 安卓设置。

# 安装 WSA

> 参考：
> 
> - [秋风于渭水，win11 安卓子系统（WSA）ROOT安装面具（Magisk）与谷歌框架（Google Apps）](https://www.tjsky.net/tutorial/384)
> - [bitxeno's notes，通过 WSA 抓取 android 的 https 网络请求包](https://blog.xenori.com/2023/05/capture-android-https-network-packet-with-wsa/)
> - [吾爱破解，[Android Tools] WSA with Magisk Root安装配置教程(2023.5)](https://www.52pojie.cn/thread-1790633-1-1.html)

> Notes: 截至 2024.1.13，微软官方还未向中国地区推送 WSA，在商店搜索不到，就算通过网页上的商店连接打开电脑上的 Microsoft Store，一会提示所在地去不可用，所以需要先修改系统中的 **国际或地区**

“系统设置”→“时间和语言”→“语言和区域”→“区域”→“国家或地区”，选择「美国」

Notes: 如果系统中的 Microsoft Store 搜不到 WSA，可以通过下面的网页链接打开 Microsoft Store 对应的页面进行安装

- https://apps.microsoft.com/detail/9P3395VX91NR?hl=en-us&gl=US WSA 本体
- WSA 增强工具（非官方）
  - https://apps.microsoft.com/detail/9PPSP2MKVTGT?hl=zh-cn&gl=cn WSA 工具箱

## 安装已 root 的 WSA

卸载官方的 WSA，后使用下面的项目构建安装包后安装。

 [GitHub 项目，LSPosed/MagiskOnWSALocal ](https://github.com/LSPosed/MagiskOnWSALocal ) 已 root 带 Magisk、Google app 的 WSA
 
- https://github.com/MustardChef/WSABuilds # 好像是把 LSPosed/MagiskOnWSALocal  项目的内容构建出来放到 release 里了，不用自己再 clone 项目后运行。

安装完成后，启用开发者模式，可以看到启动的监听端口，使用 adb 工具可以控制 WSA。

### 问题

有个 BUG，在 Create system images 时报错: ERROR: Not yet implemented。用 https://github.com/sn-o-w/MagiskOnWSALocal 这个构建能解决。

 - https://github.com/LSPosed/MagiskOnWSALocal/issues/754

Magisk 安装完模块后，重启 WSA 后模块不显示

- https://github.com/MustardChef/WSABuilds/issues/154
- 临时解决办法: https://github.com/MustardChef/WSABuilds/issues/154#issuecomment-1729105000

# 最佳实践

## 为 WSA 配置代理

[秋风于渭水，Windows Android 子系统 WSA 代理设置方法教程](https://www.tjsky.net/tutorial/391)

方法一、通过 [ADB](docs/Mobile%20device/ADB.md) 进入 shell 使用 settings 设置代理。

- 有局限性，有些 APP 的包在 Windows 的 Charles 上收不到包，WSA 中浏览器访问的包很全。

方法二、（已作废）因为 WSA 系统网络机制的更新，在较新的 WSA 上，方法三已经作废。  如果找不到名为 VirtWifi 的 wifi 说明这个方法已经不适合你的 WSA 了。
