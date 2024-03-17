---
title: systemctl 命令行工具
---

# 概述

> 参考：
>
> - [Manual(手册)，systemctl(1)](https://man7.org/linux/man-pages/man1/systemctl.1.html)

systemctl 命令用来对整个“systemd”的系统和服务进行管理

# Syntax(语法)

**systemctl \[OPTIONS] COMMAND \[UNIT...]**

UNIT 为 Unit 名称，如果指定了 UNIT 则只对这个 Unit 执行 COMMAND，如果不指定则对全局 Unit 进行操作

## OPTIONS

- **-t** # 对指定类型的 unit 进行操作
- **--all** #
- **--now** # 该选项可以与 enable、disable、mask 命令一起使用。
  - 与 enable 命令一起使用时，将同时启动该 Unit
  - 与 disable 和 mask 命令一起使用时，将同时停止该 Unit、
  - 注意：只有当 enable 或 disable 命令成功时，才会执行启动或停止操作。加了该选项就类似于执行了 `systemctl enable UNIT && systemctl start UNIT` 命令

## COMMAND 分类

- [Unit Command](#Unit%20Command) # 对 unit 执行操作的命令
- [Unit File Commands](#Unit%20File%20Commands) # 对 Unit 文件执行操作的命令
- Machine Commands
- Job Commands
- Snapshot Commands
- Environment Commands
- [Manager Lifecycle Commands](#Manager%20Lifecycle%20Commands) # 生命周期管理器的命令
- System Commands

注意：OBJECT 可以使用 PATTERN(模式)来进行匹配，i.e.使用正则表达式来通过关键字查找 unit 来显示包含这些关键字的 unit 的状态

# Unit Command

## list-units

**默认命令**，当 COMMAND 为空时，默认执行该命令列出已加载(已启动)的 UNIT

### Syntax(语法)

**systemctl list-units \[PATTERN]**

### EXAMPLE

- systemctl -t service # 查看所有 service unit 的信息，systemctl 命令默认列出所有 unit
- systemctl list-units --failed # 列出所有失败的 unit

list-sockets \[PATTERN] List loaded sockets ordered by address

list-timers \[PATTERN] List loaded timers ordered by next elapse

{start | stop | restart} UnitName # 立刻启动或者关闭或者重启某个 Unit

## reload NAME

不关闭 UNIT 的情况下重新载入配置文件，让配置生效，只重新加载.conf 类的文件

try-restart NAME... Restart one or more units if active

reload-or-restart NAME... Reload one or more units if possible, otherwise start or restart

reload-or-try-restart NAME... Reload one or more units if possible,otherwise restart if active

## isolate NAME

启动一个 unit 并关闭其他的。如果指定的 Unit 没写扩展名，则默认 target。

这个命令的作用类似于老的 init 系统中修改运行级别的效果

EXAMPLE

- systemctl isolate multi-user.target # 启动 multi-user.target 这个 unit 并关闭其他(类似于切换成纯文本运行方式)
- systemctl isolate graphical.target # 类似于切换成图形模式

kill NAME... Send signal to processes of a unit

is-active PATTERN... Check whether units are active

is-failed PATTERN... Check whether units are failed

## status [PATTERN|PID]

显示整个系统的 Unit 状态信息包括树状关联信息,如果指定了 `[]` 中的内容,则显示指定 Unit 运行时的状态信息

EXAMPLE

- systemctl status ssh.service # 查看 ssh.service 这个 unit 的状态
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/iqtd0r/1616167368492-74c581c3-e6a6-48e6-b6db-49e0fd799a63.jpeg)

show \[PATTERN...|JOB...] Show properties of one or more units/jobs or the manager

## cat PATTERN

显示一个或多个 unit 的文件及其内容

查看 sshd 这个服务的配置文件路径以及配置文件的内容，效果如下

```bash
~]# systemctl cat sshd
# /lib/systemd/system/ssh.service
[Unit]
Description=OpenBSD Secure Shell server
Documentation=man:sshd(8) man:sshd_config(5)
After=network.target auditd.service
ConditionPathExists=!/etc/ssh/sshd_not_to_be_run

[Service]
EnvironmentFile=-/etc/default/ssh
ExecStartPre=/usr/sbin/sshd -t
ExecStart=/usr/sbin/sshd -D $SSHD_OPTS
ExecReload=/usr/sbin/sshd -t
ExecReload=/bin/kill -HUP $MAINPID
KillMode=process
Restart=on-failure
RestartPreventExitStatus=255
Type=notify
RuntimeDirectory=sshd
RuntimeDirectoryMode=0755

[Install]
WantedBy=multi-user.target
Alias=sshd.service
```

set-property NAME ASSIGNMENT... Sets one or more properties of a unit

help PATTERN...|PID... Show manual for one or more units

reset-failed \[PATTERN...] Reset failed state for all, one, or more units

## list-dependencies

列出服务的依赖关系，树状显示。默认列出 default.target 的依赖树，即 default.target **被哪些服务依赖**。

### Syntax(语法)

**systemctl list-dependencies \[OPTIONS] \[UnitName]**

可以指定 Unit，以列出该 UNIT **被哪些服务依赖**

**OPTIONS**

- **--reverse** # 反向追踪，列出该 UNIT **依赖于哪些服务**。即该 UNIT 可以给谁提供依赖，即被谁需要，即启动哪些 UNIT 需要以这个 UNIT 启动为前提

### EXAMPLE

列出 sshd.service 这个 unit 可以给谁提供依赖

- systemctl list-dependencies sshd.service --reverse

# Unit File Commands

## list-unit-files [PATTERN...]

列出所有已经安装的 Unit 的配置文件。（目录为/usr/lib/systemd/system/下的所有文件）

## {enable|disable} NAME

启用或禁用一个或多个 Unit 文件

reenable NAME... Reenable one or more unit files

preset NAME... Enable/disable one or more unit files based on preset configuration

preset-all Enable/disable all unit files based on preset configuration

is-enabled NAME... Check whether unit files are enabled

mask NAME... Mask one or more units

unmask NAME... Unmask one or more units

link PATH... Link one or more units files into the search path

add-wants TARGET NAME... Add 'Wants' dependency for the target on specified one or more units

add-requires TARGET NAME... Add 'Requires' dependency for the target on specified one or more units

edit NAME... Edit one or more unit files

## get-default # 获取 default.target 的名字

获取引导进入的 default.target。获取的 TargetUnit 名字是(会通过符号链接的方式)default.target 的别名
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/iqtd0r/1616167368520-867ed609-6df3-41f4-8a42-48d7b7497340.jpeg)

## set-default NAME # 设置 default.target

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/iqtd0r/1616167368506-ab34c24b-0f8f-4a0f-8618-8d6ad5a75673.jpeg)
设置引导进入的 default.target。这个设置(会通过符号链接的方式)会用给定的 TargetUnit 给 default.target 起一个别名。相当于给一个 target 类型的 unit 建立了一个名为 default.target 的软链接

1. EXAMPLE
   1. systemctl set-defult graphical.target # 给 graphical.target 创建一个名为 default.target 的软连接到 /etc/systemd/system/ 目录下

# Machine Commands

list-machines \[PATTERN...] List local containers and host

# Job Commands

list-jobs \[PATTERN...] List jobs

cancel \[JOB...] Cancel all, one, or more jobs

# Snapshot Commands

snapshot \[NAME] Create a snapshot

delete NAME... Remove one or more snapshots

# Environment Commands

show-environment Dump environment

set-environment NAME=VALUE... Set one or more environment variables

unset-environment NAME... Unset one or more environment variables

import-environment \[NAME...] Import all or some environment variables

# Manager Lifecycle Commands

## daemon-reload

重新加载所有 daemon 的配置文件，包括.service 等文件一起重新加载

daemon-reexec Reexecute systemd manager

# System Commands

is-system-running Check whether system is fully running

default Enter system default mode

rescue Enter system rescue mode

emergency Enter system emergency mode

halt Shut down and halt the system

poweroff Shut down and power-off the system

**reboot \[ARG]** 关闭并重启系统

- [Util-linux Utilities](docs/1.操作系统/Linux%20管理/Util-linux%20Utilities.md) 包中的 reboot 命令不再维护后，reboot 命令的软链接就是 systemctl。

kexec Shut down and reboot the system with kexec

exit Request user instance exit

switch-root ROOT \[INIT] Change to a different root file system

suspend Suspend the system

hibernate Hibernate the system

hybrid-sleep Hibernate and suspend the system
