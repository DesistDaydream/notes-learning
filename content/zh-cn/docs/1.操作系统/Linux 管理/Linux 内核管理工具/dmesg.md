---
title: "dmesg"
linkTitle: "dmesg"
weight: 20
---

# 概述

> 参考：
>
> - [Manual(手册)，dmesg(1)](https://man7.org/linux/man-pages/man1/dmesg.1.html)

dmesg 命令是用来在 Unix-like 系统中显示内核的相关信息的。**dmesg** 全称是 **display message (or display driver)**，即显示信息。默认操作是显示来自内核环形缓冲区的所有消息。

实际上，dmesg 命令是从内核环形缓冲区中获取数据的。当我们在 Linux 上排除故障时，dmesg 命令会十分方便，它能很好地帮我们鉴别硬件相关的 error 和 warning。除此之外，dmesg 命令还能打印出守护进程相关的信息，已帮助我们 debug。

## dmesg \[OPTIONS]

**OPTIONS**

- **-L, --color** # 输入内容带上颜色。
- **-l, --level LIST** # 指定输出的级别，多个级别以逗号分隔。可用的级别有以下几种
  - emerg - system is unusable
  - alert - action must be taken immediately
  - crit - critical conditions
  - err - error conditions
  - warn - warning conditions
  - notice - normal but significant condition
  - info - informational
  - debug - debug-level messages
- **-f, --facility LIST** # 指定要输出的 [Facility(设施)](/docs/6.可观测性/Logs/日志系统.md#Facility(设施))，多个设施以逗号分隔。可用的 Facility 有[Facility(设施)](/docs/6.可观测性/Logs/日志系统.md)
  - ser - random user-level messages
  - mail - mail system
  - daemon - system daemons
  - auth - security/authorization messages
  - syslog - messages generated internally by syslogd
  - lpr - line printer subsystem
  - news - network news subsystem
- **-H, --human** # 启用人类可读的输出。是 `--color、--reltime、--nopager` 这三个选项的结合
- **-T, --ctime** # 打印人类可读的时间戳。
  - 请注意，时间戳记可能不正确！挂起/恢复系统后，用于日志的时间源不会更新。根据引导时间和单调时钟之间的当前增量调整时间戳，这仅适用于上次恢复后打印的消息。

**EXAMPLE**
