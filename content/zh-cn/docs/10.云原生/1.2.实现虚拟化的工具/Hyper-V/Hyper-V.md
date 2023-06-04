---
title: "Hyper-V"
linkTitle: "Hyper-V"
weight: 20
---

# 概述

> 参考：
> -

virtmgmt 命令行工具打开图形化的 Hyper-V 图形界面

# 关联文件与配置

**编辑会话设置** # 用以设置连接信息。基本都是 tsclient 的设置。


# 共享

> 参考：
> - [官方文档-虚拟化，与你的虚拟机共享设备](https://learn.microsoft.com/zh-cn/virtualization/hyper-v-on-windows/user-guide/enhanced-session-mode)

增强会话模式可通过 RDP（远程桌面协议）将 Hyper-V 与虚拟机连接起来。 这不仅会改善你的整体虚拟机查看体验，而且使用 RDP 连接还可以使虚拟机与你的计算机共享设备。 由于 RDP 在 Windows 10 中默认打开，所以与 Windows 虚拟机连接时，你可能已经在使用 RDP。 本文着重介绍了一些好处和连接设置对话框中的隐藏选项。

RDP/增强会话模式：

-   使虚拟机实现可调整大小和高 DPI 感知。
-   改进虚拟机集成
    -   共享的剪贴板
    -   通过拖放和复制粘贴进行文件共享
-   允许设备共享
    -   麦克风/扬声器
    -   USB 设备
    -   数据磁盘（包括 C:）
    -   打印机

想要共享宿主机的磁盘，最好的方法是使用 RDP 的 tsclient 功能，效果如下图：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/hyper-v/20230209155050.png)


# 使用 PowerShell 模块管理 Hyper-V

> 参考：
> - [官方文档-PowerShell，模块-hyper-v](https://learn.microsoft.com/en-us/powershell/module/hyper-v/index)
> - https://learn.microsoft.com/zh-cn/virtualization/hyper-v-on-windows/quick-start/try-hyper-v-powershell


## New-VM

> 参考：
> - https://learn.microsoft.com/en-us/powershell/module/hyper-v/new-vm

创建一个新的虚拟机

## Get-VM

> 参考：
> - https://learn.microsoft.com/en-us/powershell/module/hyper-v/get-vm

从一个或多个 Hyper-V 主机获取虚拟机

```powershell
PS ~> get-vm

Name  State CPUUsage(%) MemoryAssigned(M) Uptime   Status   Version
----  ----- ----------- ----------------- ------   ------   -------
win10 Off   0           0                 00:00:00 正常运行 11.0
```