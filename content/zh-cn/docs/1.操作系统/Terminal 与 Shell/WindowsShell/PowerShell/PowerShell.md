---
title: "PowerShell"
linkTitle: "PowerShell"
weight: 1
---

# 概述

> 参考：
>
> - [官方文档](https://learn.microsoft.com/en-us/powershell)
> - [官方文档，关于](https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.core/about/about?view=powershell-7.3)

PowerShell 是一种跨平台的任务自动化解决方案，由命令行 shell、脚本语言和配置管理框架组成。 PowerShell 在 Windows、Linux 和 macOS 上运行。

PowerShell 的独特之处在于，它<font color="#ff0000">接受并返回 .NET 对象</font>，而<font color="#ff0000">非纯文本</font>。这个特点让 PowerShell 可以更轻松地在一个管道中串联不通的命令。

> [!Notes] 这里面所说的对象，就是[面向对象编程](/docs/2.编程/计算机科学/Object-oriented%20Programming/OOP.md)中常说的“对象”，就像 Go 语言中的 Struct 类似的东西，只不过是 .NET 语言中的对象。

这些对象在被接收后，再交给格式化函数处理，以人类可读的方式，输出出来。

我们可以在 [PowerShell 官方文档的参考 - 关于](https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.core/about/about) 部分找到对 PowerShell 的基本概念的描述。

## PowerShell 变量

详见 [PowerShell 变量](/docs/1.操作系统/Terminal%20与%20Shell/WindowsShell/PowerShell/PowerShell%20变量.md) 章节

## PowerShell 命令

> 参考：
>
> - [官方文档，PowerShell 命令是什么？](https://learn.microsoft.com/zh-cn/powershell/scripting/powershell-commands)
> - [about_Command_Precedence](https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.core/about/about_command_precedence) 介绍 PowerShell 如何确定要运行的命令。
> - [about_Command_Syntax](https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.core/about/about_command_syntax) 介绍 PowerShell 中使用的语法关系图
> - [about_Core_Commands](https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.core/about/about_core_commands) 列出设计用于 PowerShell 提供程序的 cmdlet

PowerShell 中可以执行的命令分两类

- 操作系统中的可执行文件
- cmdlet

PowerShell 内置了一组 **cmdlet(全称 command-lets)**，cmdlet 是一种命令的统称，并不是一个独立的可执行文件。cmdlet 被收集在 PowerShell 模块中，可以按需加载它们。可以用任何编译的 .NET 语言或 PowerShell 脚本语言来编写 cmdlet。

cmdlet 中每个命令的名称都是由 **Verb-Noun(动词-名词)** 组成，比如 Get-Command 命令用于获取在 CLI 中注册的所有 cmdlet。

我们可以通过如下几个命令来搜索可用的 cmdlet

- `Get-Verb` # 获取所有可用的动词
- `Get-Command` # 此命令会检索计算机上安装的所有命令的列表。
- `Get-Member` # 基于其他命令运行，可以获取 cmdlet 返回的对象信息，包括对象中的**属性、方法、等等**
- `Get-Help` # 以命令名称为参数调用此命令，将显示一个帮助页面，其中说明了命令的各个部分。

从本笔记的 [Windows 管理工具](/docs/1.操作系统/Windows%20管理/Windows%20管理工具/Windows%20管理工具.md) 目录查找所有可用的命令，以及查看命令的用法

我们可以通过 `$psversiontable` 和 `$host` 变量查看 PowerShell 版本信息

```powershell
PS C:\> $psversiontable

Name                           Value
----                           -----
PSVersion                      7.3.2
PSEdition                      Core
GitCommitId                    7.3.2
OS                             Microsoft Windows 10.0.19045
Platform                       Win32NT
PSCompatibleVersions           {1.0, 2.0, 3.0, 4.0…}
PSRemotingProtocolVersion      2.3
SerializationVersion           1.1.0.1
WSManStackVersion              3.0

PS C:\> $host

Name             : ConsoleHost
Version          : 7.3.2
InstanceId       : 518ca4c4-e959-4d51-b3bb-cdcb3d5a1484
UI               : System.Management.Automation.Internal.Host.InternalHostUserInterface
CurrentCulture   : zh-CN
CurrentUICulture : zh-CN
PrivateData      : Microsoft.PowerShell.ConsoleHost+ConsoleColorProxy
DebuggerEnabled  : True
IsRunspacePushed : False
Runspace         : System.Management.Automation.Runspaces.LocalRunspace
```

# 关联文件与配置

https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_profiles?view=powershell-7.5

PowerShell 的配置文件称为 Profiles。与 Bash 类似，根据文件所在路径，有优先级和生效范围的概念，不同路径下的文件可以针对 单个/所有 的 用户/主机 生效。下面列出的文件优先级从低到高，所有用户所有主机的配置文件优先级最低。

- `$PSHOME/Profile.ps1` # 所有用户，所有主机
- `$PSHOME/Microsoft.PowerShell_profile.ps1` # 所有用户，当前主机
- `$HOME/Documents/PowerShell/Profile.ps1` # 当前用户，所有主机
- `$HOME/Documents/PowerShell/Microsoft.PowerShell_profile.ps1` # 当前用户，当前主机

Notes: 这里面的 `$HOME/Documents/` 有可能因为重定向 Documents 而与实际不同，直接在 Powshell 输入 `echo env:$PROFILE` 可以查看当前用户当前主机配置文件路径：

```powershell
PS C:\Users\DesistDaydream> echo env:$PROFILE
env:D:\Documents\PowerShell\Microsoft.PowerShell_profile.ps1
```

使用 `$PROFILE | Select-Object *` 可以列出所有优先级的配置文件

```powershell
PS C:\Users\DesistDaydream> $PROFILE | Select-Object *

AllUsersAllHosts       : C:\Program Files\PowerShell\7\profile.ps1
AllUsersCurrentHost    : C:\Program Files\PowerShell\7\Microsoft.PowerShell_profile.ps1
CurrentUserAllHosts    : D:\Documents\PowerShell\profile.ps1
CurrentUserCurrentHost : D:\Documents\PowerShell\Microsoft.PowerShell_profile.ps1
Length                 : 56
```

# 安装与更新

# 使用 PowerShell

`powershell` 和 `pwsh` 这几个命令一般都是用来打开 PowerShell 的，同时也是执行 PowerShell 脚本的前置命令。就像执行 Bash 脚本前加个 `bash` 命令一样

# ITEM

> 参考：
>
> - [官方文档，脚本-基本概念-Items](https://learn.microsoft.com/en-us/powershell/scripting/lang-spec/chapter-03#33-items)

PowerShell 中会抽象出一个 **Item(项)** 的概念，Item 可以一个 **Alias(别名)**、**Variable(变量)**、**Function(函数)**、**EnvironmentVariable(环境变量)**、甚至可以是文件系统中的 **File(文件)** 或者 **Directory(目录)**。

我们常用的 `ls` 命令，在 PowerShell 中其实就是调用了 `Get-ChildItem` 命令

# Providers(提供程序) 和 Drives(驱动器)

在 PowerShell 中，**Providers** 和 **Drives** 是提供对不同数据源（如文件系统、注册表、Certificate 等）的访问的特定接口。Drives 则是实际代表特定数据源的容器，比如本地磁盘驱动器、注册表驱动器等。使用 PowerShell 可以对这些数据源进行管理和操作。

详见：[提供程序](/docs/1.操作系统/Terminal%20与%20Shell/WindowsShell/PowerShell/提供程序.md)

# 常见问题

[powershell、cmd终端修改编码](https://www.cnblogs.com/xututu6/p/16574454.html)
