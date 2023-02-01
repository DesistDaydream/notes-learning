---
title: "Security"
linkTitle: "Security"
weight: 20
---

# 概述
> 参考：
> - [官方文档，PowerShell-参考-Security 模块](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.security)
> -


# Get-ExecutionPolicy

获取当前会话的执行策略。
- Restricted 执行策略不允许任何脚本运行。  
- AllSigned 和 RemoteSigned 执行策略可防止 Windows PowerShell 运行没有数字签名的脚本。

默认使用 Restricted 策略，此时当我们执行脚本时将会失败，并报错：
```
无法加载文件 XXXXX，因为在此系统上禁止运行脚本。有关详细信息，请参阅 https:/go.microsof
t.com/fwlink/?LinkID=135170 中的 about_Execution_Policies。
```

## Syntax(语法)


## EXAMPLE

# Set-ExecutionPolicy

为 Windows 计算机设置 PowerShell 执行策略

## Syntax(语法)

Get-Command

## EXAMPLE

设置策略为 RemoteSigned
- Set-ExecutionPolicy RemoteSigned