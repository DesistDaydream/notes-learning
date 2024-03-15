---
title: "Core"
linkTitle: "Core"
weight: 2
---

# 概述

> 参考：
> 
> - [官方文档-PowerShell，模块 - Core](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core)

核心模块包含管理 PowerShell 基本功能的 cmdlet 和提供程序。

# Get-Command

https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/get-command

获取所有命令

## Syntax(语法)

**Get-Command \[OPTIONS]**

**OPTIONS**

- **-Name**(PATTERN) # 列出匹配到名字的命令。支持通配符。`默认值：None`
- **-CommandType**(STRING) # 列出指定类型的命令。`默认值：cmdlet,function,alias`。可用的类型有：Alias、All、Application、Cmdlet、ExternalScript、Filter、Function、Script

## EXAMPLE
