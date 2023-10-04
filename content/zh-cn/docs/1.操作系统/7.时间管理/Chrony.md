---
title: Chrony
---

# 概述

> 参考：
> 
> - [GitLab，chrony/chrony](https://gitlab.com/chrony/chrony)
> - [官网](https://chrony.tuxfamily.org/index.html)
> - [官方文档](https://chrony.tuxfamily.org/documentation.html)

Chrony 是 NTP(网络时间协议) 的通用实现。它可以将系统时钟与 NTP 服务器，参考时钟（例如 GPS 接收器）以及使用手表和键盘进行的手动输入进行同步。它还可以充当 NTPv4（RFC 5905）服务器并与之对等，以向网络中的其他计算机提供时间服务。

`yum install chrony` 即可安装该工具

# Chrony 配置

**/etc/chrony.conf** # chronyd 守护进程运行的配置文件。该文件官方说明：<https://chrony.tuxfamily.org/doc/3.4/chrony.conf.html>

## 基础配合文件示例

```bash
# 指定 NTP 服务器。使用指定的 NTP 服务来同步本地时间。可以使用本机 ip，表示与本机同步时间。
server 172.40.0.3 iburst
# Record the rate at which the system clock gains/losses time.
driftfile /var/lib/chrony/drift
# Allow the system clock to be stepped in the first three updates
# if its offset is larger than 1 second.
makestep 1.0 3
# Enable kernel synchronization of the real-time clock (RTC).
rtcsync
# Enable hardware timestamping on all interfaces that support it.
#hwtimestamp *
# Increase the minimum number of selectable sources required to adjust
# the system clock.
#minsources 2
# 使用 allow 指令，即表明该服务器作为 NTP 服务端，来为其他 NTP 客户端来提供时钟服务
# allow 指令用于指定允许哪些 NTP 客户端来连接该服务器来校准时间。
#allow 172.40.0.0/16
# 即使不与时间源同步，也可以提供时间同步服务。在局域网内，常与 allow 指令同时使用，作为局域网内的时间服务器。
#local stratum 10
# Specify file containing keys for NTP authentication.
#keyfile /etc/chrony.keys
# Specify directory for log files.
logdir /var/log/chrony
# Select which information is logged.
#log measurements statistics tracking
```

## Chrony 配置详解

> 参考：
> - [官方文档,Manual(手册)](https://chrony.tuxfamily.org/doc/4.1/chrony.conf.html)

chrony.conf 配置文件与 nginx 的配置文件类似，由指令和指令的参数组成。每条指令放在单独的一行上。虽然受支持的指令数量很多，但通常仅需要其中几个就可以满足日常所需。下面介绍几个常用的指令

**server** # 该指令用来指定 NTP 服务器(i.e.使用指定的 NTP 服务来同步本地时间)。可以使用本机 ip，表示与本机同步时间。
该指令支持的 OPTIONS:

- iburst # 使用此选项，发送到服务器的前四个请求之间的间隔将为 2 秒或更短，而不是 minpoll 选项指定的间隔，这使 chronyd 在启动后不久即可进行第一次时钟更新

**allow \<IPRANGE>** # 该指令用于指定允许哪些 NTP 客户端来连接该服务器来校准时间。
使用 allow 指令，即表明该服务器作为 NTP 服务端，来为其他 NTP 客户端来提供时钟服务。

> 可以使用 all 参数来允许所有服务器

**local \<ARGS>** # 该指令标名 chrony 启动本地引用模式。
本地引用模式允许 chronyd 作为一个 NTP 服务器运行，即使它从来没有同步过，或者时钟的最后一次更新发生在很久以前，它也能实时同步(从客户轮询的角度来看)。

# chronyc 命令行工具

> 参考：
> 
> - [官方文档,Manual(手册)](https://chrony.tuxfamily.org/doc/4.1/chronyc.html)

**chronyc \[OPTIONS] \[COMMAND]**

chronyc 通过子命令来实现各种功能

**OPTIONS**

## System Clock COMMAND

### tracking # 显示有关系统时钟性能的参数

```bash
[lichenhao@hw-cloud-xngy-jump-server-linux-2 ~]$ chronyc tracking
Reference ID    : 647D0028 (100.125.0.40)
Stratum         : 4
Ref time (UTC)  : Sat Oct 09 08:49:47 2021
System time     : 0.000041988 seconds slow of NTP time
Last offset     : -0.000040913 seconds
RMS offset      : 0.000098497 seconds
Frequency       : 9.403 ppm slow
Residual freq   : -0.004 ppm
Skew            : 0.156 ppm
Root delay      : 0.184026539 seconds
Root dispersion : 0.026236653 seconds
Update interval : 256.5 seconds
Leap status     : Normal
```

**makestep**

**maxupdateskew**

**waitsync**

## Time Sources COMMAND

### sources # 显示 chronyd 进程访问的当前时间源的信息

```bash
[lichenhao@hw-cloud-xngy-jump-server-linux-2 ~]$ chronyc sources -v
210 Number of sources = 1

  .-- Source mode  '^' = server, '=' = peer, '#' = local clock.
 / .- Source state '*' = current synced, '+' = combined , '-' = not combined,
| /   '?' = unreachable, 'x' = time may be in error, '~' = time too variable.
||                                                 .- xxxx [ yyyy ] +/- zzzz
||      Reachability register (octal) -.           |  xxxx = adjusted offset,
||      Log2(Polling interval) --.      |          |  yyyy = measured offset,
||                                \     |          |  zzzz = estimated error.
||                                 |    |           \
MS Name/IP address         Stratum Poll Reach LastRx Last sample
===============================================================================
^* 100.125.0.40                  3   8   377   250  +5730ns[  +16us] +/-  133ms
```

### sourcestats # 显示有关 chronyd 进程所使用的每个时间源的状态信息

```bash
[lichenhao@hw-cloud-xngy-jump-server-linux-2 ~]$ chronyc sourcestats -v
210 Number of sources = 1
                             .- Number of sample points in measurement set.
                            /    .- Number of residual runs with same sign.
                           |    /    .- Length of measurement set (time).
                           |   |    /      .- Est. clock freq error (ppm).
                           |   |   |      /           .- Est. error in freq.
                           |   |   |     |           /         .- Est. offset.
                           |   |   |     |          |          |   On the -.
                           |   |   |     |          |          |   samples. \
                           |   |   |     |          |          |             |
Name/IP Address            NP  NR  Span  Frequency  Freq Skew  Offset  Std Dev
==============================================================================
100.125.0.40               10   7   38m     +0.006      0.157   +913ns    63us
```

**selectdata**

**reselect**

**reselectdist**

## NTP Sources COMMAND

### activity # 报告在线与离线的服务端和对等体的数量

```bash
[lichenhao@hw-cloud-xngy-jump-server-linux-2 ~]$ chronyc activity
200 OK
1 sources online
0 sources offline
0 sources doing burst (return to online)
0 sources doing burst (return to offline)
0 sources with unknown address
```

**authdata**

### ntpdata \[ADDRESS] # 显示 NTP 源的信息

**add peer**

**add pool**

**add server**

**delete ADDRESS
**
**burst**

## Manual Time Input COMMAND

## NTP Access COMMAND

## Monitoring Access COMMAND

## Real-time Clock COMMAND

## Other Daemon COMMAND
