---
title: "WindowsShell变量"
linkTitle: "WindowsShell变量"
weight: 20
---

# 概述
> 参考：
> - https://sysin.org/blog/windows-env/


# 变量
引用方式：`%VAR%`

输出变量的值：`$env:<VarName>`
- `(type env:path) -split ';'` # 切割字符串，将 ; 替换为换行符。方便查看

## 常用环境变量
%USERPROFILE% # 用户家目录。默认值：`C:/Users/${USERNAME}/`

%TMP% # 临时目录。默认值：
- 系统级 `C:/WINDOWS/TEMP`
- 用户级 `%USERPROFILE%/AppData/Local/Temp`

%APPDATA% # 用作应用程序特定数据的通用存储库的文件系统目录。默认值：`%USERPROFILE%/AppData/Roaming/`