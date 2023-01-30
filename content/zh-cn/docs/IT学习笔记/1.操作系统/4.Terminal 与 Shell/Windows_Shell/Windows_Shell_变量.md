---
title: "Windows Shell 变量"
linkTitle: "Windows Shell 变量"
weight: 20
---

# 概述
> 参考：
> - [官方文档，公认的环境变量](https://learn.microsoft.com/en-us/windows/deployment/usmt/usmt-recognized-environment-variables)
> - https://sysin.org/blog/windows-env/
> - https://ss64.com/nt/syntax-variables.html

赋值方式：
- 在 CMD 中：
- 在 PowerShell 中：`$env:VAR=VALUE`

引用方式：
- 在 CMD 和资源管理器中：`%VAR%`
- 在 PowerShell 中：`$env:VAR`
	- `(type env:path) -split ';'` # 切割字符串，将 ; 替换为换行符。方便查看


## 常用环境变量

UserName # 

USERPROFILE # 用户家目录。默认值：`C:/Users/${USERNAME}/`

TMP # 临时目录。默认值：
- 系统级 `C:/WINDOWS/TEMP`
- 用户级 `%USERPROFILE%/AppData/Local/Temp`

APPDATA # 用作应用程序特定数据的通用存储库的文件系统目录。默认值：`%USERPROFILE%/AppData/Roaming/`

Windows 中没有指向 “文档”、“视频” 等等目录的变量，可以使用 `[environment]::getfolderpath("mydocuments")` 命令获取。参考：https://stackoverflow.com/questions/3492920/is-there-a-system-defined-environment-variable-for-documents-directory


# 最佳实践

设置代理
- `$env:HTTPS_PROXY="http://127.0.0.1:7890"`