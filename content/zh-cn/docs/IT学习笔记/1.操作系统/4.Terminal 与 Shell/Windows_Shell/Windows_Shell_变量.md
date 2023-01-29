---
title: "Windows Shell 变量"
linkTitle: "Windows Shell 变量"
weight: 20
---

# 概述
> 参考：
> - [官方文档，公认的环境变量](https://learn.microsoft.com/en-us/windows/deployment/usmt/usmt-recognized-environment-variables)
> - https://sysin.org/blog/windows-env/



# 变量
引用方式：`%VAR%`

输出变量的值：`$env:<VarName>`
- `(type env:path) -split ';'` # 切割字符串，将 ; 替换为换行符。方便查看

Windows 中没有指向 “文档”、“视频” 等等目录的变量，可以使用 `[environment]::getfolderpath("mydocuments")` 命令获取。参考：https://stackoverflow.com/questions/3492920/is-there-a-system-defined-environment-variable-for-documents-directory

## 常用环境变量
%USERPROFILE% # 用户家目录。默认值：`C:/Users/${USERNAME}/`

%TMP% # 临时目录。默认值：
- 系统级 `C:/WINDOWS/TEMP`
- 用户级 `%USERPROFILE%/AppData/Local/Temp`

%APPDATA% # 用作应用程序特定数据的通用存储库的文件系统目录。默认值：`%USERPROFILE%/AppData/Roaming/`