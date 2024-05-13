---
title: sudoers 配置详解
---

# 概述

> 参考：
>
> - [Manual(手册),sudoers(5)-sudoers 文件格式](https://man7.org/linux/man-pages/man5/sudoers.5.html#SUDOERS_FILE_FORMAT)

使用 `visudo` 命令可直接进入编辑模式以编辑 /etc/sudoers 文件

sudoers 文件使用 [**Extended Backus-Naur Form(简称 EBNF)**](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form) 格式书写。

sudoers 文件由如下几类条目组成：

- **Aliases(别名)** #
- **User Specifications(用户规范)** # 即指定谁可以运行什么程序
- **Defaults** # Defaults 条目的配置，可以在运行时变更 sudo 的运行行为。
  - 比如配置日志输出路径等等。

## Aliases

用户规范条目中的每个字段都可以使用别名来表示，提前设置好别名，然后在这个字段使用别名来表示

sudoers 中可以设置 4 种别名：

- User_Alias(用户别名)
- Runas_Alias()
- Host_Alias(主机别名)
- Cmnd_Alias(命令别名)

别名语法：
在别名的设置中，可以使用 `*` 这种通配符来匹配所有

- User_Alias 别名名称 = 用户名,...
- Host_Alias 别名名称 = 主机名,...
- Cmnd_Alias 别名名称 = 命令,...

### 默认配置示例

```bash
## 主机别名
## 对于一组服务器，你可能会更喜欢使用主机名（可能是全域名的通配符）
## 或IP地址代替，这时可以配置主机别名
# Host_Alias     FILESERVERS = fs1, fs2
# Host_Alias     MAILSERVERS = smtp, smtp2

## 用户别名
## 这并不很常用，因为你可以通过使用组来代替一组用户的别名
# User_Alias ADMINS = jsmith, mikem


## 命令别名
## 指定一系列相互关联的命令（当然可以是一个）的别名，通过赋予该别名sudo权限，
## 可以通过sudo调用所有别名包含的命令，下面是一些示例

## 网络操作相关命令别名
Cmnd_Alias NETWORKING = /sbin/route, /sbin/ifconfig, /bin/ping, /sbin/dhclient,
 /usr/bin/net, /sbin/iptables, /usr/bin/rfcomm, /usr/bin/wvdial, /sbin/iwconfig,
 /sbin/mii-tool

## 软件安装管理相关命令别名
Cmnd_Alias SOFTWARE = /bin/rpm, /usr/bin/up2date, /usr/bin/yum

## 服务相关命令别名
Cmnd_Alias SERVICES = /sbin/service, /sbin/chkconfig

## 本地数据库升级命令别名
Cmnd_Alias LOCATE = /usr/sbin/updatedb

## 磁盘操作相关命令别名
Cmnd_Alias STORAGE = /sbin/fdisk, /sbin/sfdisk, /sbin/parted, /sbin/partprobe, /bin/mount, /bin/umount

## 代理权限相关命令别名
Cmnd_Alias DELEGATING = /usr/sbin/visudo, /bin/chown, /bin/chmod, /bin/chgrp

## 进程相关命令别名
Cmnd_Alias PROCESSES = /bin/nice, /bin/kill, /usr/bin/kill, /usr/bin/killall

## 驱动命令别名
Cmnd_Alias DRIVERS = /sbin/modprobe
```

## User Specifications

User Specifications(用户规范) 用来指定谁可以运行什么程序，可以理解为赋权。

用户规范语法：

- 用户 登录的主机=\[(可以变换的身份) ] 可以执行的命令,...

### 默认配置示例

```bash
## 下面是规则配置：什么用户在哪台服务器上可以执行哪些命令（sudoers文件可以在多个系统上共享）
## 语法:
##      user    MACHINE=COMMANDS
##  用户 登录的主机=（可以变换的身份） 可以执行的命令
##
## 命令部分可以附带一些其它的选项
##
## 允许root用户执行任意路径下的任意命令
root    ALL=(ALL)       ALL

# 允许sys用户组中的用户使用NETWORKING等所有别名中配置的命令
# %sys ALL = NETWORKING, SOFTWARE, SERVICES, STORAGE, DELEGATING, PROCESSES, LOCATE, DRIVERS

# 允许wheel用户组中的用户执行所有命令
# %wheel        ALL=(ALL)       ALL

## 允许wheel用户组中的用户在不输入该用户的密码的情况下执行所有命令
# %wheel        ALL=(ALL)       NOPASSWD: ALL

# 允许users用户组中的用户像root用户一样使用mount、unmount、chrom命令
# %users  ALL=/sbin/mount /mnt/cdrom, /sbin/umount /mnt/cdrom

# 允许users用户组中的用户像root用户一样使用shutdown命令
# %users  localhost=/sbin/shutdown -h now
```

## Defaults

Defaults 语法：

- Default_Type Parameter_List
  - Default_type 可用的值有如下几个：
    - Defautls
    - Defaults @ Host_List
    - Defaults : User_List
    - Defaults ! Cmnd_List
    - Defaults > Runas_List
  - Parameter_List 格式如下：
    - Parameter = Value
    - Parameter += Value
    - Parameter -= Value
    - ! Parameter
  - 可用的参数参考下文 Defautls 条目参数

要在参数中包含文字反斜杠字符，必须对反斜杠进行两次转义。例如，要匹配 “\ n” 作为命令行参数的一部分，必须在 sudoers 文件中使用 “\ \ n”。这是由于存在两个转义级别，一个在 sudoers 解析器本身中，另一个在 fnmatch(3) 函数匹配命令行参数时。

### 默认配置示例

```bash
# 如果无法在终端上禁用 echo，则拒绝运行
Defaults   !visiblepw

# 开启sudo日志，让sudo命令每次执行都写入到/var/log/sudo.log文件中
Defaults logfile=/var/log/sudo.log

# Preserving HOME has security implications since many programs
# use it when searching for configuration files. Note that HOME
# is already set when the the env_reset option is enabled, so
# this option is only effective for configurations where either
# env_reset is disabled or HOME is present in the env_keep list.
#
Defaults    always_set_home
Defaults    match_group_by_gid

# Prior to version 1.8.15, groups listed in sudoers that were not
# found in the system group database were passed to the group
# plugin, if any. Starting with 1.8.15, only groups of the form
# %:group are resolved via the group plugin by default.
# We enable always_query_group_plugin to restore old behavior.
# Disable this option for new behavior.
Defaults    always_query_group_plugin

Defaults    env_reset
Defaults    env_keep =  "COLORS DISPLAY HOSTNAME HISTSIZE KDEDIR LS_COLORS"
Defaults    env_keep += "MAIL PS1 PS2 QTDIR USERNAME LANG LC_ADDRESS LC_CTYPE"
Defaults    env_keep += "LC_COLLATE LC_IDENTIFICATION LC_MEASUREMENT LC_MESSAGES"
Defaults    env_keep += "LC_MONETARY LC_NAME LC_NUMERIC LC_PAPER LC_TELEPHONE"
Defaults    env_keep += "LC_TIME LC_ALL LANGUAGE LINGUAS _XKB_CHARSET XAUTHORITY"

# Adding HOME to env_keep may enable a user to run unrestricted
# commands via sudo.
#
# Defaults   env_keep += "HOME"

# 下面这个指令指定当用户执行 sudo 命令时在什么地方寻找二进制代码和命令。
# 这个选项的目的显然是要限制用户运行 sudo 命令的范围，这是一种好做法。
# 说白了，就是替换 $PATH 变量的值。
Defaults    secure_path = /sbin:/bin:/usr/sbin:/usr/bin
```

# [Defaults 条目参数](https://man7.org/linux/man-pages/man5/sudoers.5.html#SUDOERS_OPTIONS)

## 日志相关

**logfile = <FILE>** # 使用本地文件记录日志，并指定文件的绝对路径。默认情况下，sudo 使用 syslog 记录日志。
**syslog = <FACILITY>** # 使用 syslog 记录日志，并指定 syslog 的日志设施。`默认值：authpriv`。

- 可用的设施有：**authpriv**(if your OS supports it), **auth**, **daemon**,**user**, **local0**, **local1**, **local2**, **local3**, **local4**, **local5**,**local6**, and **local7**.

## 其他

**secure_path = <PATH:PATH:...>** # 如果设置了该选项，则通过 sudo 运行命令时，将会用该选项的值替代用户设置的 `$PATH` 变量。
**visiblewp** # 默认情况下，如果用户必须输入密码，但无法在终端上禁用 echo，sudo 将拒绝运行。如果设置了 visiblepw 标志，即使在屏幕上可见，sudo 也会提示输入密码。这使得运行 “ssh somehost sudo ls” 之类的东西成为可能，因为默认情况下，ssh(1) 在运行命令时不分配 tty。默认情况下，此标志处于关闭状态。

# 用户赋权的实用案例

EXAMPLE

- 别名
  - 设定别名 DOCKER，该 DOCKER 别名包括 docker 命令和 systemctl 中子命令对 docker 服务的操作
    - Cmnd_Alias DOCKER = /usr/bin/docker*, /usr/bin/systemctl* docker\*
- 赋权
  - 表示 desistdaydream 用户可以在所有主机执行所有命令
    - desistdaydream ALL=ALL #
  - 表示 desistdaydream 用户在所有主机，变换为 root 身份，可以执行所有命令
    - desistdaydream ALL=(root) ALL

赋予 developer 用户，可以操作 docker 和 nginx 的所有权限。i.e.通过 systemctl 控制 docker 和 nginx 服务，使用 docker 和 nginx 相关命令

- **Cmnd_Alias DOCKER = /usr/bin/systemctl * docker*, /usr/bin/docker*** # 为某些命令设置别名
- **Cmnd_Alias NGINX = /usr/bin/systemctl * nginx*, /usr/bin/nginx*** # 为某些命令设置别名
- **developer ALL=(root) NOPASSWD:DOCKER,NGINX** # 让 developer 用户可以以 root 用户且不使用密码执行 DOCKER 和 NGINX 别名中的所有命令
- **Defaults logfile = /var/log/sudo/sudo.log** # 将 sudo 的日志保存到其他目录，不通过 rsyslog 保存。
