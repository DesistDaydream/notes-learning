---
title: "PowerShell 内置管理工具"
linkTitle: "PowerShell 内置管理工具"
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，PowerShell 模块参考](https://learn.microsoft.com/en-us/powershell/module)

[PowerShell](/docs/1.操作系统/Terminal%20与%20Shell/WindowsShell/PowerShell/PowerShell.md) 内置的管理工具都是 cmdlet，以模块形式提供，这个目录下记录的笔记中，每个文件或目录名都是一个模块的名称

# CimCmdlets

https://learn.microsoft.com/en-us/powershell/module/cimcmdlets

# 通用模块

[Core](/docs/1.操作系统/Windows%20管理/Windows%20管理工具/PowerShell%20内置管理工具/Core.md)

[Management](/docs/1.操作系统/Windows%20管理/Windows%20管理工具/PowerShell%20内置管理工具/Management.md)

[Security](/docs/1.操作系统/Windows%20管理/Windows%20管理工具/PowerShell%20内置管理工具/Security.md)

[CimCmdlets](/docs/1.操作系统/Windows%20管理/Windows%20管理工具/PowerShell%20内置管理工具/CimCmdlets.md)

[Utility](/docs/1.操作系统/Windows%20管理/Windows%20管理工具/PowerShell%20内置管理工具/Utility.md)

etc.

# 特定于 Windows 的模块

https://learn.microsoft.com/en-us/powershell/module/?view=windowsserver2025-ps

## NetTCPIP

> 参考：
>
> - [官方文档 - PowerShell，参考 - NetTCPIP](https://learn.microsoft.com/en-us/powershell/module/nettcpip)


### Get-NetTCPConnection

https://learn.microsoft.com/en-us/powershell/module/nettcpip/get-nettcpconnection

#### Syntax(语法)

**OPTIONS**

- **-LocalPort**(\[]INT) # 查看指定的端口，多个端口以 `,` 分割。
- **-State**(\[]STRING) # 查看指定 TCP 状态的端口。可用的值有: Bound, Closed, CloseWait, Closing, DeleteTCB, Established, FinWait1, FinWait2, LastAck, Listen, SynReceived, SynSent, TimeWait

#### Example

利用该模块可以比 netstat 命令更方便得获取各种基于网络连接的信息以及监听该端口的进程信息

获取监听在 1080 端口上的程序的路径

- `(get-process -id (Get-NetTCPConnection -LocalPort 10800 -State Listen).OwningProcess).path`

甚至可以像这样组合出人类可读的信息

```powershell
PS C:\Users\DesistDaydream> Get-NetTCPConnection -LocalPort 1080 -State Listen | Select-Object LocalAddress, LocalPort, @{Name="PID";Expression={$_.OwningProcess}}, @{Name="Path";Expression={(Get-Process -Id $_.OwningProcess -FileVersionInfo).FileName}}

LocalAddress LocalPort   PID Path
------------ ---------   --- ----
::1               1080 21804 D:\Tools\VanDyke Software\SecureCRT\SecureCRT.exe
127.0.0.1         1080 21804 D:\Tools\VanDyke Software\SecureCRT\SecureCRT.exe
```