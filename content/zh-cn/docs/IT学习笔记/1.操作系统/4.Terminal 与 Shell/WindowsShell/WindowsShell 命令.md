---
title: "WindowsShell 命令"
linkTitle: "WindowsShell 命令"
weight: 20
---

# 概述
> 参考：
> - [官方文档，Windows Server-命令](https://learn.microsoft.com/en-us/windows-server/administration/windows-commands/windows-commands)
> - [官方文档，PowerShell-模块参考](https://learn.microsoft.com/en-us/powershell/module)

所有受支持的 Windows 和 Windows Server 版本都内置了一组 Win32 控制台命令。同时，[PowerShell](docs/IT学习笔记/1.操作系统/4.Terminal%20与%20Shell/WindowsShell/PowerShell.md) 也内置了一组 cmdlet

这些命令可以用来管理 Windows，很多命令的详解可以参见 [Windows 管理工具](/docs/IT学习笔记/1.操作系统/Y.Windows%20管理/Windows管理工具/_index.md) 目录。

本质上，内置的命令就两类：

- Win32 控制台命令。一般在保存 `C:/Windows/System32/` 目录中，就像 Unix 的 `/usr/bin` 这种目录似的，都是一些可执行文件。
- PowerShell 中的 cmdlet。也就是 PowerShell 的各种模块。这些 cmdlet 虽然不是可见的可执行文件，但是也可以实现类似命令的效果。