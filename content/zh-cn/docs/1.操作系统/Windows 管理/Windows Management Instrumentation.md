---
title: Windows Management Instrumentation
linkTitle: Windows Management Instrumentation
date: 2024-01-07T22:42
weight: 20
---

# 概述

> 参考：
>
> - [官方文档](https://learn.microsoft.com/en-us/windows/win32/wmisdk/wmi-start-page)
> - [Wiki，Windows_Management_Instrumentation](https://en.wikipedia.org/wiki/Windows_Management_Instrumentation)

**Windows Management Instrumentation(简称 WMI)** 使用 [CIM](/docs/Standard/IT/DMTF.md#CIM)(通用信息模型) 行业标准来表示系统、应用程序、网络、设备和其他受管理的组件。

https://www.syscom.com.tw/ePaper_Content_EPArticledetail.aspx?id=76&EPID=159&j=4&HeaderName=%E7%A0%94%E7%99%BC%E6%96%B0%E8%A6%96%E7%95%8C

## 在 PowerShell 中使用 WMI

https://learn.microsoft.com/zh-cn/powershell/scripting/learn/ps101/07-working-with-wmi

Windows [PowerShell](/docs/1.操作系统/Terminal%20与%20Shell/WindowsShell/PowerShell/PowerShell.md) 早期 WMI cmdlet 已弃用，在 PowerShell 6+ 中不可用，请改用 CIM cmdlet。

PowerShell 中存在多个本机 WMI cmdlet，且无需安装任何其他软件或模块。 `Get-Command` 可用于确定 Windows PowerShell 中存在哪些 WMI cmdlet。 以下结果来自运行 5.1 版 PowerShell 的 Windows 10 实验环境计算机。 结果因运行的 PowerShell 版本而异。

```powershell
Get-Command -Noun WMI*
```

```
CommandType     Name                                               Version    Source
-----------     ----                                               -------    ------
Cmdlet          Get-WmiObject                                      3.1.0.0    Microsof...
Cmdlet          Invoke-WmiMethod                                   3.1.0.0    Microsof...
Cmdlet          Register-WmiEvent                                  3.1.0.0    Microsof...
Cmdlet          Remove-WmiObject                                   3.1.0.0    Microsof...
Cmdlet          Set-WmiInstance                                    3.1.0.0    Microsof...
```

PowerShell 版本 3.0 中引入了通用信息模型 (CIM) cmdlet。 CIM cmdlet 的设计目的是使其可以同时在 Windows 和非 Windows 计算机上使用。

所有 CIM cmdlet 都包含在一个模块中。 若要获取 CIM cmdlet 的列表，请结合使用 `Get-Command` 与 Module 参数，如以下示例中所示。

```powershell
Get-Command -Module CimCmdlets
```

```
CommandType     Name                                               Version    Source
-----------     ----                                               -------    ------
Cmdlet          Export-BinaryMiLog                                 1.0.0.0    CimCmdlets
Cmdlet          Get-CimAssociatedInstance                          1.0.0.0    CimCmdlets
Cmdlet          Get-CimClass                                       1.0.0.0    CimCmdlets
Cmdlet          Get-CimInstance                                    1.0.0.0    CimCmdlets
Cmdlet          Get-CimSession                                     1.0.0.0    CimCmdlets
Cmdlet          Import-BinaryMiLog                                 1.0.0.0    CimCmdlets
Cmdlet          Invoke-CimMethod                                   1.0.0.0    CimCmdlets
Cmdlet          New-CimInstance                                    1.0.0.0    CimCmdlets
Cmdlet          New-CimSession                                     1.0.0.0    CimCmdlets
Cmdlet          New-CimSessionOption                               1.0.0.0    CimCmdlets
Cmdlet          Register-CimIndicationEvent                        1.0.0.0    CimCmdlets
Cmdlet          Remove-CimInstance                                 1.0.0.0    CimCmdlets
Cmdlet          Remove-CimSession                                  1.0.0.0    CimCmdlets
Cmdlet          Set-CimInstance                                    1.0.0.0    CimCmdlets
```
