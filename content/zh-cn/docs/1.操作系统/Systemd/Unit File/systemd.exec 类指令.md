---
title: "systemd.exec 类指令"
linkTitle: "systemd.exec 类指令"
date: "2023-08-01T11:29"
weight: 20
---

# 概述

> 参考：
>
> - [Manual(手册)，systemd.exec(5)](https://man7.org/linux/man-pages/man5/systemd.exec.5.html)

systemd.exec 类的指令是 [Unit File 指令中特殊部分的指令中的通用指令](/docs/1.操作系统/Systemd/Unit%20File/Unit%20File%20指令.md#通用指令) 的一种，可以配置进程执行时的环境，主要用于 service、socket、mount、swap 部分。

systemd.exec 包含很多很多指令，我们可以将其分为如下几大类：

- [PATHS](#paths) # 路径相关指令
- [USER/GROUP IDENTITY](#USER/GROUP%20IDENTITY) # 用户/组标识相关指令
- CAPABILITIES
- SECURITY
- MANDATORY ACCESS CONTROL
- PROCESS PROPERTIES
- SCHEDULING
- SANDBOXING
- SYSTEM CALL FILTERING
- [ENVIRONMENT](#environment) # 环境变量相关指令
- LOGGING AND STANDARD INPUT/OUTPUT
- CREDENTIALS
- SYSTEM V COMPATIBILITY

比如 环境变量、运行程序的用户和组、运行路径 等等

# Paths

https://man7.org/linux/man-pages/man5/systemd.exec.5.html#PATHS

Paths(路径) 相关的指令可用于更改文件系统的服务视图。请注意，路径必须是绝对路径，并且不得包含 `..` 路径组件。

**WorkingDirectory=\<STRING>** # 采用相对于由 RootDirectory 指令 或特殊值 `~` 指定的服务根目录的目录路径。

# User/Group Identity

https://man7.org/linux/man-pages/man5/systemd.exec.5.html#USER/GROUP_IDENTITY

- **User=\<STRING>** # 指定运行该 Unit 使用的用户。

# CAPABILITIES(能力)相关指令

https://man7.org/linux/man-pages/man5/systemd.exec.5.html#CAPABILITIES

# SECURITY(安全) 相关指令

https://man7.org/linux/man-pages/man5/systemd.exec.5.html#SECURITY

# MANDATORY ACCESS CONTROL(强制访问控制) 相关指令

https://man7.org/linux/man-pages/man5/systemd.exec.5.html#MANDATORY_ACCESS_CONTROL

# PROCESS PROPERITES(进程属性) 相关指令

https://man7.org/linux/man-pages/man5/systemd.exec.5.html#PROCESS_PROPERTIES

为执行的进程设置各种资源的软限制和硬限制。

# SCHEDULING(调度)相关指令

https://man7.org/linux/man-pages/man5/systemd.exec.5.html#SCHEDULING

# SANDBOXING(沙盒) 相关指令

https://man7.org/linux/man-pages/man5/systemd.exec.5.html#SANDBOXING

# SYSTEM CALL FILTERING(系统调用过滤) 相关指令

https://man7.org/linux/man-pages/man5/systemd.exec.5.html#SYSTEM_CALL_FILTERING

# Environment

https://man7.org/linux/man-pages/man5/systemd.exec.5.html#ENVIRONMENT

**Environment(STRING)** # 指定 Unit 所使用的环境变量。多个变量以空格分隔

- e.g. `Environment="VAR1=word1 word2" VAR2=word3 "VAR3=$word 5 6"`

**EnvironmentFile(STRING)** # 与 Environment 指令的逻辑类似，但是可以直接指定一个文件，在文件中设置环境变量，文件中的格式与 Environment 指令的值的格式保持一致。

- **`-` 符号** # 指令的值之前可以添加 `-` 符号，比如 `EnvironmentFile=-/etc/default/ssh`。`-` 符号作为前缀的话，表明如果文件不存在将不会被读取，并且不会记录任何错误或警告信息。
  - 默认情况下，若指定的文件不存在，服务将无法启动并报错。

# LOGGING AND STANDARD INPUT/OUTPUT(日志的标准输入/输出) 相关指令

https://man7.org/linux/man-pages/man5/systemd.exec.5.html#LOGGING_AND_STANDARD_INPUT/OUTPUT

# 分类

> #systemd #unit-file
