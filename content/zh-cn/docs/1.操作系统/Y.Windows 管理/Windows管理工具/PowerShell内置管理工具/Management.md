---
title: "Management"
linkTitle: "Management"
weight: 20
---

# 概述
> 参考：
> - [官方文档-PowerShell，模块-Management](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.management)

管理模块可以让我们在 PowerShell 中管理系统中的进程、服务等。

# Item 管理工具

## Get-ChildItem
> 参考：
> - [官方文档-PowerShell，模块-管理-Get-ChildItem](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.management/get-childitem)

### Syntax(语法)
**OPTIONS**
- **-Path \<STRING>** # 指定一个或多个位置的路径，可以使用通配符。`默认值：.`

### EXAMPLE

实现类似 tree 命令的效果

- Get-ChildItem -Path D:/Tools -Recurse -Depth 2 | Select-Object FullName


# 服务管理工具

> 参考：
> - [官方文档，PowerShell-脚本示例-管理进程和服务-管理服务](https://learn.microsoft.com/en-us/powershell/scripting/samples/managing-services)
> - [官方文档，.Net-开发 Windows 服务应用](https://learn.microsoft.com/zh-cn/dotnet/framework/windows-services/)

## 启动、停止、重启、暂停服务

- Start-Service # 启动服务
- Stop-Service # 停止服务
- Suspend-Service # 暂停服务
- Resume-Service # 恢复服务
- Restart-Service # 重启服务

## Get-Service

列出服务。Get-Service 获取代表计算机上服务的对象，包括正在运行和已停止的服务。默认情况下，当不带参数运行 Get-Service 时，将返回本地计算机的所有服务。

可以通过指定服务名称或服务的显示名称来指示此 cmdlet 仅获取特定服务，或者您可以将服务对象通过管道传递给此 cmdlet。

默认显示三个字段

- Status # 服务状态
- Name  # 服务名称
- DisplayName  # 服务的显示名称

服务名称与显示名称可以在窗口页面看到效果，显示名称有点类似于简短的描述信息。


### Syntax(语法)

Get-Service \[OPTIONS] [-Name] \<PATTERN>

PATTERN 支持通配符，前面的 -Name 可以省略，该命令默认通过**服务名称**进行匹配，将会列出所有匹配到的服务。

OPTIONS

- **-DependentServices** # 列出指定服务**被哪些服务依赖**。
- **-RequiredServices** # 列出指定服务**依赖于哪些服务**。即.若想该服务正常运行则必须要提前运行的其他服务
- **-Include <String[]>** # 
- **-Exclude <String[]>** # 

### EXAMPLE

列出服务名 s 开头的所有服务并按照状态排序

- Get-Service "s*" | Sort-Object status


## New-Service

创建服务

## Remove-Service

移除服务

## Set-Service

设置服务

# 进程管理工具


[Get-Process](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.management/get-process?view=powershell-7.3)
[Start-Process](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.management/start-process?view=powershell-7.3)
[Stop-Process](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.management/stop-process?view=powershell-7.3)
[Wait-Process](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.management/wait-process?view=powershell-7.3)

[Debug-Process](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.management/debug-process?view=powershell-7.3)

