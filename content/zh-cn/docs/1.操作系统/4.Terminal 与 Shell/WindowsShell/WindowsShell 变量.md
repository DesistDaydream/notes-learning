---
title: "WindowsShell 变量"
linkTitle: "WindowsShell 变量"
weight: 20
---

# 概述

> 参考：
> 
> - [官方文档-PowerShell，脚本-基本概念-环境变量](https://learn.microsoft.com/en-us/powershell/scripting/lang-spec/chapter-03?view=powershell-7.3#312-environment-variables)
> - https://sysin.org/blog/windows-env/
> - https://ss64.com/nt/syntax-variables.html

赋值方式：

- 在 CMD 中：
- 在 PowerShell 中：`$env:VAR=VALUE`

引用方式：

- 在 CMD 和资源管理器中：`%VAR%`
- 在 PowerShell 中：`$env:VAR`
	- `(type env:path) -split ';'` # 切割字符串，将 ; 替换为换行符。方便查看

TODO: Windows 中的变量好像不区分大小写？

## 常用环境变量

**COMPUTERNAME** # 主机名

**USERNAME** # 用户名

**USERPROFILE** # 用户家目录。默认值：`C:/Users/${USERNAME}/`

**TMP** # 临时目录。默认值：

- 系统级 `C:/WINDOWS/TEMP`
- 用户级 `%USERPROFILE%/AppData/Local/Temp`

**APPDATA** # 应用程序的数据保存路径。默认值：`%USERPROFILE%/AppData/Roaming/`

- 这个目录下的数据通常可以随着网络连接同步到其他电脑。比如用户的配置、插件等等。当然，很多时候，应用程序也会将这些可以在网络同步的数据保存到 文档、家目录 等等地方中。

**LOCALAPPDATA** # 应用程序的本地数据保存路径。默认值：`%USERPROFILE%/AppData/Local/`

**ProgramData** # 指定程序数据文件夹的路径。与 Program Files 文件夹不同，应用程序可以使用此文件夹为标准用户存储数据，因为它不需要提升的权限。`默认值：C:/ProgramData`

**ProgramFiles** # `默认值：C:/Program Files`

注意：
- Windows 中没有指向 “文档”、“视频” 等等目录的变量，可以在 PowerShell 中使用 `[environment]::getfolderpath("mydocuments")` 获取。
    - 参考：https://stackoverflow.com/questions/3492920/is-there-a-system-defined-environment-variable-for-documents-directory

想要获取变量，有几下几种方式：

- Get-Item Env:*
- ls env:
- Get-Variable

# 变量管理工具

在 PowerShell 中，变量也可以称为一个 Item，因此可以由大多数与 Item 相关的 cmdlet 命令控制。
-   [New-Variable](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/new-variable?view=powershell-7.3): Creates a variable
-   [Set-Variable](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/set-variable?view=powershell-7.3): Creates or changes the characteristics of one or more variables
-   [Get-Variable](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/get-variable?view=powershell-7.3): Gets information about one or more variables
-   [Clear-Variable](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/clear-variable?view=powershell-7.3): Deletes the value of one or more variables
-   [Remove-Variable](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/remove-variable?view=powershell-7.3): Deletes one or more variables

# 最佳实践

## 设置代理

PowerShell

- `$env:HTTPS_PROXY="http://127.0.0.1:7890"`
