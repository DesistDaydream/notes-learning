---
title: sysctl
---

# 概述

sysctl 是 [procps 工具集](/docs/IT学习笔记/1.操作系统/X.Linux%20管理/Linux%20系统管理工具/procps%20工具集.md) 中的一个用于控制内核参数的工具

# Syntax(语法)

**sysctl \[OPTIONS] \[VARIABLE\[=VALUE]] \[...]**

在运行环境中配置内核参数。VARIABLE 为内核的一个变量

**OPTIONS**

- **-a** # 显示所有变量
- **-p \[/PATH/TO/FILE]** # 从文件中读取值,默认文件为/etc/sysctl.conf。可以指定从哪个文件来读取参数，可使用通配符。
- **-w** # 允许写一个值到变量中

# EXAMPLE

- sysctl -w net.ipv4.ip_forward=1 # 开启 IP 转发模式
- sysctl -p /etc/sysctl.d/\* # 从 sysctl.d 目录中读取所有文件的内容加载到内核中