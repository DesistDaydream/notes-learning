---
title: nmcli 命令行工具
---

# 概述

> 参考：
> - [红帽官方文档](https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/7/html/networking_guide/sec-using_the_networkmanager_command_line_tool_nmcli)
> - [Manual(手册)，nmcli(1)](https://networkmanager.dev/docs/api/latest/nmcli.html)

nmcli 用于 NetworkManager 的命令行工具

# Syntax(语法)

**nmcli \[OPTIONS] OBJECT { COMMAND | help }**

OBJECT 和 COMMAND 可以用全称也可以用简称，最少可以只用一个字母

OPTIONS

- -a, --ask ask for missing parameters
- -c, --colors auto|yes|no whether to use colors in output
- -e, --escape yes|no escape columns separators in values
- **-f, --fields \<FIELD,...>|all|common** # 指定要输出的字段，FIELD 可以是 任意 setting
- -g, --get-values \<field,...>|all|common shortcut for -m tabular -t -f
- -h, --help print this help
- **-m, --mode \<tabular|multiline>** # 指定输出模式,tabular 输出为表格样式，multiline 是多行样式。
  - nmcli con show # 默认为表格样式
  - nmcli con show DEV # 默认为多行样式
- -o, --overview overview mode
- **-p, --pretty** # 美化输出，以连接中的 setting 分段落展示
- -s, --show-secrets allow displaying passwords
- **-t, --terse** # 简洁的输出
- -v, --version show program version
- -w, --wait \<seconds> set timeout waiting for finishing operations

OBJECT

- **g\[eneral]** # NetworkManager's general status and operations
- **n\[etworking]** # overall networking control
- **r\[adio]** # NetworkManager radio switches
- **c\[onnection]** # NetworkManager's connections
- **d\[evice]** # devices managed by NetworkManager
- **a\[gent]** # NetworkManager secret agent or polkit agent
- **m\[onitor]** # monitor NetworkManager changes

# g\[eneral] # NetworkManager 的一般状态和操作

# n\[etworking] # overall networking control

# r\[adio] # NetworkManager radio switches

# c\[onnection] # Connections 的管理命令，常用命令

详见：[connection 子命令章节](https://www.yuque.com/go/doc/33221854)

## clone 克隆连接

clone \[--temporary] \[id | uuid | path ] \<ID> \<new name> # 克隆连接

## edit 在交互模式的编辑器中修改连接

edit \[id | uuid | path] \<ID> # 进入交互编辑器修改连接

edit \[type \<new_con_type>] \[con-name \<new_con_name>]

## monitor 监控连接

monitor \[id | uuid | path] \<ID> ... # 监控连接

## reload、load 加载连接信息

从磁盘重新加载所有连接文件。 默认情况下，NetworkManager 不会监视对连接文件的更改。 因此，您需要使用此命令来告诉 NetworkManager 在对它们进行更改时从磁盘重新读取连接配置文件。 但是，可以启用自动加载功能，然后 NetworkManager 会在每次更改连接文件时重新加载它们（NetworkManager.conf（5）中的 monitor-connection-files = true）。

从磁盘加载/重新加载一个或多个连接文件。 手动编辑连接文件后使用此命令可确保 NetworkManager 知道其最新状态。

reload # 从 /etc/sysconfig/network-scripts/ 目录重新加载连接文件

load \<filename> \[ \<filename>... ] #

## 其他

import \[--temporary] type \<type> file \<file to import>

export \[id | uuid | path] \<ID> \[\<output file>]

# d\[evice] # 通过 NetworkManager 来管理网络设备

# a\[gent] # NetworkManager secret agent or polkit agent

# m\[onitor] # monitor NetworkManager changes
