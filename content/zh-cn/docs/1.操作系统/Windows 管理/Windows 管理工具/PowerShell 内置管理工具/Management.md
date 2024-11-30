---
title: "Management"
linkTitle: "Management"
weight: 20
---

# 概述

> 参考：
>
> - [官方文档-PowerShell，模块 - Management](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.management)

管理模块可以让我们在 PowerShell 中管理系统中的 进程、服务、Item 等。

# Item 管理工具

- [Clear-Item](https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.management/clear-item?view=powershell-7.3)
- [Copy-Item](https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.management/copy-item?view=powershell-7.3)
- [Get-Item](https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.management/get-item?view=powershell-7.3)
- [Invoke-Item](#Invoke-Item)
- [Move-Item](https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.management/move-item?view=powershell-7.3)
- [New-Item](#New-Item)
- [Remove-Item](https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.management/remove-item?view=powershell-7.3)
- [Rename-Item](https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.management/rename-item?view=powershell-7.3)
- [Set-Item](https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.management/set-item?view=powershell-7.3)
- [Get-ChildItem](#Get-ChildItem) # 获取指定位置的 Item 和 子Item。类似 ls 命令

## New-Item

https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.management/new-item

创建一个新的 Item 并设置它的值。可以创建的 Item 类型取决于当前所使用的 ProwerShell 提供程序。例如，在文件系统中，创建文件、目录、符号链接、等等；在注册表中，创建注册表条目

### Syntax(语法)

**OPTIONS**

- **-ItemType STRING** # 指定新 Item 类型，可用的类型取决于 PowerShell 的[提供程序](/docs/1.操作系统/Terminal%20与%20Shell/WindowsShell/PowerShell/提供程序.md)
  - 由于不同提供可用的类型非常多，笔记里就不写了，具体还是看 Net-Item 官方文档吧

### EXAMPLE

创建符号链接(软连接)

- 创建 C:/Users/DesistDaydream/AppData/Roaming/yuzu 符号链接文件，指向 E:/emulator/yuzu_data/user
  - `New-Item -ItemType SymbolicLink -Path "C:/Users/DesistDaydream/AppData/Roaming/yuzu" -Target "E:/emulator/yuzu_data/user"`
  - Notes: 这种用法可以代替 mklink 命令
- 查看符号链接文件所指向的原始文件路径
  - `(Get-Item ${PathToFile}).Target`



## Invoke-Item

https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.management/invoke-item?view=powershell-7.3

对指定的 Item 执行默认操作，默认操作取决于 Item 的类型。比如 目录类型的 Item，则使用默认的资源管理器打开、.docs 文件类型的 Item，则使用 .docs 默认的程序打开、等等

`ii` 是 `Invoke-Item` 的别名，我们可以直接使用 `ii .` 命令使用资源管理器打开当前目录。

## Get-ChildItem

> 参考：
>
> - [官方文档-PowerShell，模块-管理-Get-ChildItem](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.management/get-childitem)

ls 命令是 Get-ChildItem 的别名。

`Get-ChildItem` 获取一个或多个指定位置中的所有 Item。

默认情况下 `Get-ChildItem` ，列出 mode(模式) 、 LastWriteTime(最后编辑时间)、Length(长度，即.文件大小) 、Name(即. Item 的名称 。 **Mode** 属性中的字母可以解释如下：

- `l` # 链接
- `d` # 目录
- `a` # 存档
- `r` # 只读
- `h` # 隐藏
- `s` # 系统
- `-` # 普通文件

### Syntax(语法)

**OPTIONS**

- **-Path \<STRING>** # 指定一个或多个位置的路径，可以使用通配符。`默认值：.`

### EXAMPLE

实现类似 tree 命令的效果

- Get-ChildItem -Path D:/Tools -Recurse -Depth 2 | Select-Object FullName

# 服务管理工具

> 参考：
>
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

**OPTIONS**

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

[Get-Process](#Get-Process)
[Start-Process](#Start-Process)
[Stop-Process](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.management/stop-process?view=powershell-7.3)
[Wait-Process](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.management/wait-process?view=powershell-7.3) # 等待一个或多个正在运行的进程停止，然后再接受输入。

[Debug-Process](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.management/debug-process?view=powershell-7.3)

## Start-Process

https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.management/start-process?view=powershell-7.3

在本地计算机上启动一个或多个进程。

Start-Process 在本地计算机上启动一个或多个进程。默认情况下，Start-Process 创建一个新进程，该进程继承当前进程中定义的所有环境变量。

若要指定进程中运行的程序，请输入可执行文件或脚本文件，或者可以使用计算机上的程序打开的文件。 如果指定非可执行文件， `Start-Process` 则启动与该文件关联的程序，类似于 `Invoke-Item` cmdlet。

可以使用 的参数 `Start-Process` 指定选项，例如加载用户配置文件、在新窗口中启动进程或使用备用凭据。

### Syntax(语法)

**OPTIONS**

- **-ArgumentList** #
- **-WindowStyle** # 指定新进程的窗口样式。`默认值：Normal`
  - Hideen # 隐藏窗口
  - Minimized # 最小化窗口
  - Maximized # 最大化窗口
- **-RedirectStandardOutput** # 将进程产生的输出发送到指定的文件中。默认输出到控制台。
- **-RedirectStandardError** # 将进程产生的所有错误发送到指定的文件中。默认输出到控制台。

### EXAMPLE

Rclone 使用示例

```powershell
Start-Process "alist.exe" -ArgumentList "server --data D:\appdata\alist" -WindowStyle Hidden -RedirectStandardOutput "D:\Tools\Scripts\log\alist.log" -RedirectStandardError "D:\Tools\Scripts\log\alist-err.log"

Start-Process "rclone.exe" `
-ArgumentList "mount alist-net:/ Z: --cache-dir D:\appdata\rclone-cache --vfs-cache-mode full --vfs-cache-max-age 24h --header Referer:" `
-WindowStyle Hidden `
-RedirectStandardOutput "D:\Tools\Scripts\log\rclone.log" -RedirectStandardError "D:\Tools\Scripts\log\rclone-err.log"
```

## Get-Process

https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.management/get-process?view=powershell-7.3

获取本地计算机上正在运行的进程信息。如果没有参数，此 cmdlet 将获取本地计算机上的所有进程。 还可以按进程名称或进程 ID (PID) 指定特定进程，或通过管道将进程对象传递到此 cmdlet。

默认情况下，此 cmdlet 返回一个进程对象，该对象包含有关进程的详细信息，并支持可用于启动和停止进程的方法。 还可以使用 cmdlet 的参数 `Get-Process` 获取进程中运行的程序的文件版本信息，并获取进程加载的模块。

### Syntax(语法)

**OPTIONS**

- **-Id** #

