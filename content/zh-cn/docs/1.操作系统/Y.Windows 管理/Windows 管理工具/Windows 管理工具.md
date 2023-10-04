---
title: Windows 管理工具
linkTitle: Windows 管理工具
date: 2024-01-07T09:57
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，Windows Server-命令](https://learn.microsoft.com/en-us/windows-server/administration/windows-commands/windows-commands)
> - [官方文档，PowerShell-模块参考](https://learn.microsoft.com/en-us/powershell/module)

所有受支持的 Windows 和 Windows Server 版本都内置了一组 Win32 控制台命令。同时，[PowerShell](/docs/1.操作系统/4.Terminal%20与%20Shell/WindowsShell/PowerShell/PowerShell.md) 也内置了一组 cmdlet

本质上，内置的命令就两类：

- Win32 控制台命令。一般在保存 `C:/Windows/System32/` 目录中，就像 Unix 的 `/usr/bin` 这种目录似的，都是一些可执行文件。
- PowerShell 中的 cmdlet。也就是 PowerShell 的各种模块。这些 cmdlet 虽然不是可见的可执行文件，但是也可以实现类似命令的效果。
