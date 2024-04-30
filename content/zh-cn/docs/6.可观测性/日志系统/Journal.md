---
title: Journal
weight: 4
---

# 概述

> 参考：
>
> - [Manual,systemd-journald.service(8)](https://man7.org/linux/man-pages/man8/systemd-journald.service.8.html)

相关服务说明

- systemd-Journald.service # 日志功能通过该 Unit 来实现，是一个用于收集和存储日志数据的系统服务，是系统启动前要启动的第一个进程，Journald 会把所有收集到的信息保存在内存中。
- rsyslog.service # 另一种日志数据持久化，Journald 会把日志信息转发给 rsyslog.service 进行处理和保存，如果没有 Journald，rsyslog 也可以自动生成日志而不用从 journald 去获取
- logrotate # logrotate 会对日志文件进行轮替操作，i.e.把已经非常大的日志文件改名后，创建一个新的日志文件，新产生的日志会保存在新文件中，老文件保留一定时期后会自动清除

# Journald 关联文件与配置

**/etc/systemd/journal.conf**

## 日志存放路径

**/run/log/journal/${MACHINE-ID}/**

**/var/log/journal/${MACHINE-ID}/**

默认情况下，journald 的日志保存在 /run/log/journal 中，系统重启就会清除。通过新建 /var/log/journal 目录，日志会自动记录到这个目录中，并永久存储。

路径中的 MACHINE-ID 的值，可以通过 `cat /etc/machine-id` 命令获取

```bash
~]# ls /run/log/journal
c14766a3e9ae49a3872fb9b7e2583710
~]# cat /etc/machine-id
c14766a3e9ae49a3872fb9b7e2583710
```

所有 journal 程序生成的日志，都会存在 MACHIN-ID 目录下

```bash
~]# ll -h /var/log/journal/c14766a3e9ae49a3872fb9b7e2583710
total 153M
drwxr-sr-x+ 2 root systemd-journal 4.0K Feb 21 23:15  ./
drwxr-sr-x+ 3 root systemd-journal   46 Dec  9 17:19  ../
-rw-r-----+ 1 root systemd-journal  40M Dec 28 16:23 'system@aa6b2b3f8f9d46fdb169f9d8aaab56c3-0000000000000001-0005b6048d0b7824.journal'
-rw-r-----+ 1 root systemd-journal  32M Jan 22 12:39 'system@aa6b2b3f8f9d46fdb169f9d8aaab56c3-00000000000080e9-0005b781fc8c48d9.journal'
-rw-r-----+ 1 root systemd-journal  32M Feb 21 23:10 'system@aa6b2b3f8f9d46fdb169f9d8aaab56c3-000000000000df93-0005b975c74c3caf.journal'
-rw-r-----+ 1 root systemd-journal  40M Mar 12 15:25  system.journal
-rw-r-----+ 1 root systemd-journal 8.0M Dec 28 16:23 'user-1000@571778ddc0db463990a85592631fa5e8-0000000000000496-0005b6049323448d.journal'
......
```

# journalctl 命令行工具

> 参考：
>
> - [Manual(手册)，journalctl(1)](https://man7.org/linux/man-pages/man1/journalctl.1.html)

[Systemd](docs/1.操作系统/Systemd/Systemd.md) 统一管理所有 Unit 的启动日志。带来的好处就是，可以只用 journalctl 命令，查看所有日志（内核日志和应用日志）。日志的配置文件是 /etc/systemd/journald.conf。journalctl 功能强大，用法非常多。

## Syntax(语法)

**journalctl \[OPTIONS] \[MATCHES]**

### OPTIONS

- **--disk-usage** # 显示所有日志文件的磁盘使用情况，包括持久化和临时的日志。
- **-f, --follow** # 实时更新
- **--file=FILE** # 查看指定文件中的日志信息，FILE 可以使用"?"与"\*"进行匹配。常用于查看从别的设备上拷贝过来的日志文件
- **-p UNM** # 指定要显示的日志级别(NUM 为 0-7 级)
- **-u UNIT** # 显示指定的 unit 的日志信息

格式选项

- **--no-pager** # 在单一页面显示信息，不分页。默认情况下，若日志过长，需要使用 → 方向键翻页才能查看后面的日志内容。
- **-o, --output=\<STRING>** # 指定输出格式。`默认值：short`
  - 可用的格式有：json、json-pretty、verbose、export、with-unit 等等

过滤选项

- **--output-fields=FIELD** # 显示指定字段的日志，多个字段以 `,` 分割。
  - 字段筛选仅对 -o 选项指定的 verbose、export、json、json-pretty、json-sse、json-seq 这几个输出格式有效
- **-S, --since TIME** 与 **-U, --unitl TIME** # 设置输出日志信息的开始与结束时间

## EXAMPLE

- 查看指定用户的日志
  - sudo journalctl \_UID=33 --since today
- 查看指定进程的日志
  - sudo journalctl \_PID=1
- 以 JSON 格式仅输出 MESSAGE 与 \_CMDLINE 字段的消息
  - journalctl -u docker -ojson-pretty --output-fields=MESSAGE,\_CMDLINE
- 查看指定时间的日志
  - sudo journalctl --since="2012-10-30 18:17:16"
  - sudo journalctl --since "20 min ago"
  - sudo journalctl --since yesterday
  - sudo journalctl --since 09:00 --until "1 hour ago"
  - journalctl --since "2018-11-13" --until "2018-11-14 03:00"

# 查看所有日志（默认情况下 ，只保存本次启动的日志）

$ sudo journalctl

# 查看内核日志（不显示应用日志）

$ sudo journalctl -k

# 查看系统本次启动的日志

$ sudo journalctl -b

$ sudo journalctl -b -0

# 查看上一次启动的日志（需更改设置）

$ sudo journalctl -b -1

# 显示尾部的最新 10 行日志

$ sudo journalctl -n

# 显示尾部指定行数的日志

$ sudo journalctl -n 20

# 查看指定服务的日志

$ sudo journalctl /usr/lib/systemd/systemd

# 查看某个路径的脚本的日志

$ sudo journalctl /usr/bin/bash

# 查看某个 Unit 的日志

$ sudo journalctl -u nginx.service

$ sudo journalctl -u nginx.service --since today

# 实时滚动显示某个 Unit 的最新日志

$ sudo journalctl -u nginx.service -f

# 合并显示多个 Unit 的日志

$ journalctl -u nginx.service -u php-fpm.service --since today

# 查看指定优先级（及其以上级别）的日志，共有 8 级

- 0: emerg
- 1: alert
- 2: crit
- 3: err
- 4: warning
- 5: notice
- 6: info
- 7: debug

$ sudo journalctl -p err -b

# 指定日志文件占据的最大空间

$ sudo journalctl --vacuum-size=1G

# 指定日志文件保存多久

$ sudo journalctl --vacuum-time=1years
