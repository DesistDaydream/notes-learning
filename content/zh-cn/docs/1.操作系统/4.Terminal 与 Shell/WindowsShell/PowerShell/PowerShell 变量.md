---
title: "PowerShell 变量"
linkTitle: "PowerShell 变量"
weight: 20
---

# 概述

> 参考：
> 
> - [官方文档-PowerShell，关于-关于变量](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_variables)
> - [官方文档-PowerShell，关于-关于自动变量](https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.core/about/about_variables)
> - [官方文档-PowerShell，关于-关于首选项变量](https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.core/about/about_preference_variables)
> - [官方文档-PowerShell，脚本-基本概念-环境变量](https://learn.microsoft.com/en-us/powershell/scripting/lang-spec/chapter-03?view=powershell-7.3#312-environment-variables)

PowerShell 变量名称不区分大小写，可以包含空格和特殊字符。但是官方推荐尽量避免使用空格和特殊字符，使用起来很麻烦，详见[包含特殊字符的变量名称](https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.core/about/about_variables#variable-names-that-include-special-characters)

PowerShell 中的**环境变量**与**普通变量**在声明方式和引用方式上有**明显不同**，这与 Linux Shell 的 [变量](/docs/1.操作系统/4.Terminal%20与%20Shell/Shell%20编程语言/变量.md) 不太一样。举一个非常简单的例子：

```powershell
PS C:\Users\DesistDaydream> $test_var="这是一个普通变量"
PS C:\Users\DesistDaydream> $test_var
这是一个普通变量
PS C:\Users\DesistDaydream> $env:test_var

PS C:\Users\DesistDaydream> $env:test_env_var="这是一个环境变量"
PS C:\Users\DesistDaydream> $test_env_var
PS C:\Users\DesistDaydream> $env:test_env_var
这是一个环境变量
```

## 普通变量

PowerShell 中有几种不同类型的变量：

- **User-created variables(用户创建的变量)** # 用户创建的变量由用户创建和维护。 默认情况下，仅在 PowerShell 窗口打开时，在 PowerShell 命令行中创建的变量才存在。 关闭 PowerShell 窗口时，将删除变量。 若要保存变量，请将其添加到 PowerShell 配置文件。 还可以在具有全局、脚本或本地范围的脚本中创建变量。
- **Automatic variables(自动变量)** # 自动变量存储 PowerShell 的状态。 这些变量由 PowerShell 创建，PowerShell 会根据需要更改其值，以保持其准确性。 用户无法更改这些变量的值。 例如，变量 `$PSHOME` 存储 PowerShell 安装目录的路径。
- **Preference variables(首选项变量)** # 首选项变量存储 PowerShell 的用户首选项。 这些变量由 PowerShell 创建，并使用默认值填充。 用户可以更改这些变量的值。 例如，变量 `$MaximumHistoryCount` 确定会话历史记录中的最大条目数。

### 自动变量

描述存储 PowerShell 的状态信息的变量。 这些变量由 PowerShell 创建和维护。

`$?` # 最后一个命令的执行状态。如果最后一个命令成功，值为 True，如果失败，值为 False

`$HOME` # 用户家目录的绝对路径。此变量使用 `$env:USERPROFILE` 环境变量的值。

`$PSHOME` # PowerShell 安装目录的绝对路径。

`$PWD` # 当前 PowerShell 运行时所在目录位置的绝对路径。每次执行 cd 命令都会更新该变量的值。

### 首选项变量

自定义 PowerShell 行为的变量。

## 环境变量

> 参考：
> 
> - [官方文档-PowerShell，关于-关于环境变量](https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.core/about/about_environment_variables)

环境变量存储操作系统和其他应用程序使用的数据。

在 Windows 中，可以在 3 个作用域中定义环境变量

- 系统
- 用户
- 进程

我们在 PowerShell 声明的环境变量通常都是进程作用范围的，对于系统和用户作用范围，通常就是指永久声明环境变量。

当我们想要在系统和用户作用域中永久声明环境变量时，可以使用 `Machine` 表示系统作用域，使用 `User` 表示用户作用域；也可以在 GUI 上找到`控制面板-系统-高级系统设置-高级-环境变量` 处修改。

在 PowerShell 中，我们可以引用 [WindowsShell 变量](docs/1.操作系统/4.Terminal%20与%20Shell/WindowsShell/WindowsShell%20变量.md) 中的环境变量。

### PowerShell 使用的其他环境变量

- **PATH** # 包含操作系统搜索可执行文件的文件夹位置的列表。
- **PATHEXT** # 包含 Windows 视为可执行文件的文件扩展名列表。
- **XDG** # XDG [基本目录规范定义的 XDG](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) 环境变量
  - XDG_CONFIG_HOME
  - XDG_DATA_HOME
  - XDG_CACHE_HOME

# 声明变量

- 普通变量: `$VAR_NAME=VALUE`
- 环境变量: `$env:VAR_NAME=VALUE`

# 引用变量

- 普通变量: `$VAR_NAME`
- 环境变量: `$env:VAR_NAME`

想要获取变量，有几下几种方式：

- `Get-Item Env:*`
- `ls env:`
- `Get-Variable`

# 变量管理工具

在 PowerShell 中，变量也可以称为一个 Item，因此可以由大多数与 Item 相关的 cmdlet 命令控制

-   [New-Variable](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/new-variable?view=powershell-7.3): Creates a variable
-   [Set-Variable](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/set-variable?view=powershell-7.3): Creates or changes the characteristics of one or more variables
-   [Get-Variable](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/get-variable?view=powershell-7.3): Gets information about one or more variables
-   [Clear-Variable](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/clear-variable?view=powershell-7.3): Deletes the value of one or more variables
-   [Remove-Variable](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/remove-variable?view=powershell-7.3): Deletes one or more variables

# 最佳实践

## 切割 PATH 变量

- `(type env:path) -split ';'` # 切割字符串，将 ; 替换为换行符。方便查看

## 永久设置系统或用户范围的环境变量

设置和取消系统或用户范围的环境变量：

```powershell
# 设置
[Environment]::SetEnvironmentVariable('Foo', 'Bar', 'Machine')
# 取消
[Environment]::SetEnvironmentVariable('Foo', '', 'Machine')
# 设置
[Environment]::SetEnvironmentVariable('Foo', 'Bar', 'User')
# 取消
[Environment]::SetEnvironmentVariable('Foo', '', 'User')
```


## 设置代理

PowerShell

- `$env:HTTPS_PROXY="http://127.0.0.1:7890"`

