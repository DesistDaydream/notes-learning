---
title: Login info mgmt
linkTitle: Login info mgmt
date: 2024-04-19T11:08
weight: 20
---

# 概述

> 参考：
>
> -

查看登录日志

- tail /var/log/secure

[procps](/docs/1.操作系统/Linux%20管理/Linux%20系统管理工具/procps%20工具集.md) 工具包中，有一些可以处理登录信息的工具，比如 pkill、w 等

# last - 用来显示所有登录的信息

last 是 [Util-linux Utilities](/docs/1.操作系统/Linux%20管理/Util-linux%20Utilities.md)

**EXAMPLE**

- last reboot # 查看系统重新引导的时间。i.e.客户查看设备什么时候关机再开机过

# lastb - 查看登录失败的用户信息

> Notes: 这是一个指向 last 程序的软连接

lastb 程序在不同的发行版中，所属的 Package 不同。

lastb 工具会读取 /var/log/btmp 文件，并把文件内容中记录登录失败的信息列出来。b 是 bad 的含义，lastb 就是指现实坏的登录信息

```bash
~]# lastb
root     ssh:notty    172.19.42.203    Mon Jun  7 21:49 - 21:49  (00:00)
root     ssh:notty    172.19.42.203    Mon Jun  7 21:24 - 21:24  (00:00)
root     ssh:notty    172.19.42.203    Mon Jun  7 21:22 - 21:22  (00:00)
......

btmp begins Mon Jun  7 21:11:54 2021
```

# faillog - 查看用户登录失败信息并处理 /var/log/faillog 文件有信息

**faillog \[OPTIONS]**

OPTIONS

- **-u** # 指定用户名
- **-r** # 删除失败信息

EXAMPLE

- faillog -u oracle -r # 删除 oracle 用户的登录失败信息，以便解锁 oracle 用户

# lastlog - 用来显示系统中‘所有用户’最近一次登陆的信息

# who - 通过查询 /var/run/utmp 文件来显示系统中当前登录的每个用户

# users - 用单独的一行打印出当前登录的用户，每个显示的用户名对应一个登录会话

