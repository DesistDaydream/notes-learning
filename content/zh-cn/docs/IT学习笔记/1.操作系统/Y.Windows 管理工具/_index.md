---
title: Y.Windows 管理工具
weight: 10
---

# 概述

>

# 服务管理工具

> <https://docs.microsoft.com/zh-cn/dotnet/framework/windows-services/> > <https://docs.microsoft.com/zh-cn/powershell/scripting/samples/managing-services?view=powershell-7.2>

Get-Service # 列出所有服务

# 管理链接文件

## mklink

### Syntax(语法)

mklink \[\[/d] | \[/h] | \[/j]] \<link> \<target>
为 target 创建一个名为 link 的链接文件。即 link 是要创建的新文件

EXAMPLE

- mklink /D C:/Users/DesistDaydream/AppData/Roaming/yuzu E:/yuzu/user
