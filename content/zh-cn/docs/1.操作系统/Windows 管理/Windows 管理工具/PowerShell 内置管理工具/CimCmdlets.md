---
title: CimCmdlets
linkTitle: CimCmdlets
weight: 20
---

# 概述

> 参考：
>
> - [官方文档 - PowerShell，模块 - CimCmdlets](https://learn.microsoft.com/en-us/powershell/module/cimcmdlets)

https://learn.microsoft.com/zh-cn/powershell/scripting/learn/ps101/07-working-with-wmi

PowerShell 早期使用 [Windows Management Instrumentation(简称 WMI)](/docs/1.操作系统/Windows%20管理/Windows%20Management%20Instrumentation.md) cmdlet，后改用 CIM cmdlet。可以使用 `Get-Command -Module CimCmdlets` 命令查看所有可用的 CimCmdlets

```powershell
CommandType     Name                                               Version    Source
-----------     ----                                               -------    ------
Cmdlet          Get-CimAssociatedInstance                          7.0.0.0    CimCmdlets
Cmdlet          Get-CimClass                                       7.0.0.0    CimCmdlets
Cmdlet          Get-CimInstance                                    7.0.0.0    CimCmdlets
Cmdlet          Get-CimSession                                     7.0.0.0    CimCmdlets
Cmdlet          Invoke-CimMethod                                   7.0.0.0    CimCmdlets
Cmdlet          New-CimInstance                                    7.0.0.0    CimCmdlets
Cmdlet          New-CimSession                                     7.0.0.0    CimCmdlets
Cmdlet          New-CimSessionOption                               7.0.0.0    CimCmdlets
Cmdlet          Register-CimIndicationEvent                        7.0.0.0    CimCmdlets
Cmdlet          Remove-CimInstance                                 7.0.0.0    CimCmdlets
Cmdlet          Remove-CimSession                                  7.0.0.0    CimCmdlets
Cmdlet          Set-CimInstance                                    7.0.0.0    CimCmdlets
```

# Get-CimClass

https://learn.microsoft.com/en-us/powershell/module/cimcmdlets/get-cimclass

获取特定名称空间中的 CIM 类的列表。返回一个 [CimClass 类](https://learn.microsoft.com/en-us/dotnet/api/microsoft.management.infrastructure.cimclass)

# Get-CimInstance

https://learn.microsoft.com/en-us/powershell/module/cimcmdlets/get-ciminstance

从 CIM 服务器获取类的 CIM 实例。

## Syntax(语法)

**OPTIONS**

- **-ClassName**(STRING) # 指定要检索 CIM 实例的 CIM 类的名称。您可以使用Tab键自动补全来浏览类的列表，因为 PowerShell 会从本地 WMI 服务器获取类的列表，从而提供类名的列表。

## EXAMPLE

获取主板信息

- Get-CimInstance Win32_Baseboard | Select-Object Product, Manufacturer, Version, SerialNumber
