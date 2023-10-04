---
title: psmisc 工具集
---

# 概述

> 参考：
>
> - 项目地址：<https://gitlab.com/psmisc/psmisc>

psmisc 是一个小型的应用程序集合，与 procps-ng 工具集类似，都是使用伪文件系统(/proc)内的信息来展示系统信息

该工具集包含包含以下程序(每个标题都是一个)

# fuser - 显示正在使用指定 文件 或 sockets 的进程

fuser 使用指定的文件或文件系统显示进程的 PID。 在默认的显示模式下，每个文件名后都有一个字母，表示访问类型：

- c # 当前目录。
- e # 一个可执行文件正在运行。
- f # 打开文件。 在默认显示模式下省略 f。
- F # 打开文件进行写入。 默认显示模式下省略 F。
- r # 根目录。
- m # 映射文件或共享库。

## fuser \[OPTIONS]

OPTIONS

- **-m** # 显示使用命名文件系统或块设备的所有进程

# killall - 通过进程名称向指定进程发送信号

与 kill 命令类似，但是不像 kill 只能指定进程的 PID，而是可以通过进程的名称来发送信号

EXAMPLE：

- killall -0 nginx # 向名为 nginx 的进程发送 0 信号

# peekfd - shows the data travelling over a file descriptor

# prtstat - 输出一个进程的统计信息

# pslog - prints log path(s) of a process

# pstree - 以树状显示当前正在运行的进程

该命令与 ps 类似，但是不会像 ps -ef 一样显示内核态进程

## Syntax(语法)

pstree \[OPTIONS] \[USER or PID]

OPTIONS:

- -p # 显示进程的 PID
- -a # 显示进程运行的命令行参数
- -c # 禁用相同分支的合并
- -h # 高亮显示当前进程及其父进程
- -H PID # 高亮显示指定进程
- -t # 显示完整的进程名称
- -s # 显示指定进程的父进程
- -n # 按 PID 排序
- -g # 显示 PGID。i.e 一个或多个进程组 ID

## EXAMPLE

- pstree -n

以最简单的形式调用时没有任何选项或参数，`pstree` 命令将显示所有正在运行的进程的分层树结构。

```bash
$ pstree
systemd─┬─VBoxService───7*[{VBoxService}]
  ├─accounts-daemon───2*[{accounts-daemon}]
  ├─2*[agetty]
  ├─atd
  ├─cron
  ├─dbus-daemon
  ├─irqbalance───{irqbalance}
  ├─2*[iscsid]
  ├─lvmetad
  ├─lxcfs───2*[{lxcfs}]
  ├─networkd-dispat───{networkd-dispat}
  ├─nginx───2*[nginx]
```
