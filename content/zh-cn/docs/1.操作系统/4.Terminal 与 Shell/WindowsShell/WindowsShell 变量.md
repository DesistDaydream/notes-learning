---
title: "WindowsShell 变量"
linkTitle: "WindowsShell 变量"
weight: 20
---

# 概述

> 参考：
> 
> - https://sysin.org/blog/windows-env/
> - https://ss64.com/nt/syntax-variables.html

赋值方式与引用方式详见各自 Shell 章节

- 在 CMD 中：

引用方式：

- 在 CMD 和资源管理器中：`%VAR%`

TODO: Windows 中的变量好像不区分大小写？

# 常用环境变量

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



