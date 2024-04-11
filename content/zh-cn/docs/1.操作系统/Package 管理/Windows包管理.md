---
title: "Windows包管理"
linkTitle: "Windows包管理"
weight: 20
---

# 概述

> 参考：
> 
> - 官方文档

**%LOCALAPPDATA%/Packages/** # 通过应用商店安装的程序会保存在这里？待确认

# AppX

> 参考：
> 
> - [官方文档](https://learn.microsoft.com/en-us/powershell/module/appx/)

Get-AppxPackage

Remove-AppxPackage

# MSIX

> 参考：
> 
> - 官方文档


# WinGet

> 参考：
> 
> - [官方文档-Windows，包管理器](https://learn.microsoft.com/en-us/windows/package-manager/)

winget 是一个 Windows Package Manager(Windows 包管理器)，由命令行工具 (WinGet) 和一组用于在 Windows 设备上安装应用程序的服务组成。

## 安装 winget

## Syntax(语法)

https://learn.microsoft.com/en-us/windows/package-manager/winget/#commands



## EXAMPLE

卸载 *Windows 小组件*

```powershell
winget uninstall "windows web experience pack"
```