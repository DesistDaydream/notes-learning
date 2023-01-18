---
title: Rsyslog
---

# 概述

> 参考：
> - [官网](https://www.rsyslog.com/)
> - [官方文档,配置-模块](https://www.rsyslog.com/doc/v8-stable/configuration/modules/index.html)
> - [GitHub 项目](https://github.com/rsyslog/rsyslog)
> - [Wiki,Rsyslog](https://en.wikipedia.org/wiki/Rsyslog)
> - [Manual(手册),syslog(3)](https://man7.org/linux/man-pages/man3/syslog.3.html)
> - [Manual(手册),rsyslogd(8)](https://man7.org/linux/man-pages/man8/rsyslogd.8.html)
> - [Arch 文档,Systemd-Journal-配合 syslog 使用](<https://wiki.archlinux.org/title/Systemd_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87)/Journal_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87)>)

**Rocket-fast system for log processing(像火箭一样快的日志处理系统，简称 rsyslog) **是一款开源应用程序，用于 UNIX 和 类 Unix 操作系统，可以在 IP 网络中转发日志消息。Rsyslog 实现了基本的 Syslog 协议，并扩展了丰富的功能，比如基于内容的过滤、排队处理离线输出、支持模块、灵活的配置、使用 TCP 传输 等等.

RsysLog 是一个日志统一管理的程序。通过 rsyslogd 这个守护进程提供服务，rsyslogd 程序是对 syslogd 的扩展，提供了更多的功能和可靠性。

Rsyslog 提供了一个符合 [RFC 5424](https://datatracker.ietf.org/doc/html/rfc5424) 标准的日志消息系统。

RsysLog 的特点：

- 可以监听在某个端口上作为日志服务器，来手机多个主机的日志
- RsysLog 自带多个模块，可以通过模块来实现更多的功能。以 im 开头的是在收集日志时候所用到的，以 om 开头的是在输出日志时用到的(比如把收集到的日志保存在某一文件中)。

## Moules(模块)

Rsyslog 采用模块化设计，可以通过加载模块来动态加载功能，模块也可以由任何第三方编写，只要符合 Rsyslog 规范即可。

每个模块都有参数可以配置。

## Rsyslog 日志处理

Rsyslog 使用 [**imuxscok**](https://www.rsyslog.com/doc/v8-stable/configuration/modules/imuxsock.html) 模块监听本地 Unix Socket(`默认为 /dev/log`) 以接收本地系统上运行的应用程序产生的 syslog 格式的日志消息。当该 Socket 收到消息时，会通过 syslog(3) 这里面所描述的系统调用将日志消息传递给 Rsyslog。

这个模块在 Rsyslog 的配置文件中必须进行配置，因为没有它，本地日志记录将无法进行，因为没有监听任何 Unix Socket，任何程序发往 /dev/log 的消息，也就无法接收了。

实际上，由于 syslog 协议是标准的，绝大部分编程语言，都可以通过自身或第三方实现的 syslog 库，将日志数据直接写到 syslog 的 Socket 中。下面用 go 举例

```go
package main

import (
        "log"
        "log/syslog"
)

func main() {
        sysLog, err := syslog.Dial("", "",syslog.LOG_ERR, "lichenhao")
        if err != nil {
                log.Fatal(err)
        }

        sysLog.Emerg("Hello world!")
}
```

运行一下，查看日志，可以看到，进程名 lichenhao，输出了一条日志消息

```bash
[root@hw-cloud-xngy-jump-server-linux-2 /var/log]# tail -n 1 syslog
Oct 19 23:46:03 hw-cloud-xngy-jump-server-linux-2 lichenhao[3283]: Hello world!
```

此时如果同时在查看 /dev/log 文件，也可以看到同样的内容

```bash
[root@hw-cloud-xngy-jump-server-linux-2 /home/lichenhao/test_dir]# socat - /dev/log

Broadcast message from systemd-journald@hw-cloud-xngy-jump-server-linux-2 (Tue 2021-10-19 23:49:24 HKT):

lichenhao[3820]: Hello world!
```

### 验证 Rsyslog 接收日志

看一下 rsyslog 进程打开的文件描述符

```bash
[root@hw-cloud-xngy-jump-server-linux-2 ~]# ll /proc/$(pgrep rsyslog)/fd
total 0
dr-x------ 2 root   root    0 Oct 19 21:16 ./
dr-xr-xr-x 9 syslog syslog  0 Oct 19 21:16 ../
lr-x------ 1 root   root   64 Oct 19 21:16 0 -> /dev/null
l-wx------ 1 root   root   64 Oct 19 21:16 1 -> /dev/null
l-wx------ 1 root   root   64 Oct 19 21:16 2 -> /dev/null
lrwx------ 1 root   root   64 Oct 19 21:16 3 -> 'socket:[15881]'
lr-x------ 1 root   root   64 Oct 19 21:16 4 -> /dev/urandom
lr-x------ 1 root   root   64 Oct 19 21:16 5 -> /proc/kmsg
lrwx------ 1 root   root   64 Oct 19 21:16 6 -> 'socket:[21398]'
l-wx------ 1 root   root   64 Oct 19 21:16 7 -> /var/log/syslog
l-wx------ 1 root   root   64 Oct 19 21:16 8 -> /var/log/kern.log
l-wx------ 1 root   root   64 Oct 19 21:16 9 -> /var/log/auth.log
```

追踪一下进程的系统调用(这里是执行了一下 `su - root` 命令产生的日志)

```bash
[root@hw-cloud-xngy-jump-server-linux-2 ~]# strace -p 595 -f -e recvmsg -s 1000
strace: Process 595 attached with 4 threads
[pid   626] recvmsg(3, {msg_name=NULL, msg_namelen=0, msg_iov=[{iov_base="<37>Oct 19 22:01:40 su: (to root) lichenhao on pts/1", iov_len=8096}], msg_iovlen=1, msg_control=[{cmsg_len=32, cmsg_level=SOL_SOCKET, cmsg_type=SO_TIMESTAMP_OLD, cmsg_data={tv_sec=1634652100, tv_usec=471671}}, {cmsg_len=28, cmsg_level=SOL_SOCKET, cmsg_type=SCM_CREDENTIALS, cmsg_data={pid=2549, uid=0, gid=0}}], msg_controllen=64, msg_flags=0}, MSG_DONTWAIT) = 52
[pid   626] recvmsg(3, {msg_name=NULL, msg_namelen=0, msg_iov=[{iov_base="<86>Oct 19 22:01:40 su: pam_unix(su-l:session): session opened for user root by lichenhao(uid=0)", iov_len=8096}], msg_iovlen=1, msg_control=[{cmsg_len=32, cmsg_level=SOL_SOCKET, cmsg_type=SO_TIMESTAMP_OLD, cmsg_data={tv_sec=1634652100, tv_usec=471786}}, {cmsg_len=28, cmsg_level=SOL_SOCKET, cmsg_type=SCM_CREDENTIALS, cmsg_data={pid=2549, uid=0, gid=0}}], msg_controllen=64, msg_flags=0}, MSG_DONTWAIT) = 96
```

可以发现，从 fd 为 3 的 `socket:[15881]` 这个文件接收到了日志信息

看看这个文件是个啥

```bash
[root@hw-cloud-xngy-jump-server-linux-2 ~]# lsof -p 595 -a -d 3
COMMAND  PID   USER   FD   TYPE             DEVICE SIZE/OFF  NODE NAME
rsyslogd 595 syslog    3u  unix 0xffff99c534a7d800      0t0 15881 /run/systemd/journal/syslog type=DGRAM

[root@hw-cloud-xngy-jump-server-linux-2 ~]# cat /proc/net/unix | grep 15881
ffff99c534a7d800: 00000002 00000000 00000000 0002 01 15881 /run/systemd/journal/syslog
```

两种方式都指向了同一个文件 /run/systemd/journal/syslog

```bash
[root@hw-cloud-xngy-jump-server-linux-2 ~]# ll /run/systemd/journal/syslog
srw-rw-rw- 1 root root 0 Oct 19 21:16 /run/systemd/journal/syslog=
```

这是一个 socket 文件，从 lsof 命令中可以看到是是一个用于实现 本地数据报通信的 [DGRAM 类型的 Unix Socket](/docs/IT学习笔记/1.操作系统/2.Kernel(内核)/7.Process%20 管理/Inter%20Process%20Communication(进程间通信).md 管理/Inter Process Communication(进程间通信).md)。

这个文件替代了传统的 /dev/log 文件，/dev/log 变成了指向 /run/systemd/journal/dev-log 的软链接

```bash
[root@hw-cloud-xngy-jump-server-linux-2 ~]# ll /dev/log
lrwxrwxrwx 1 root root 28 Oct 19 21:16 /dev/log -> /run/systemd/journal/dev-log=
```

但是在 CentOS 7 中，Rsyslog 依然直接使用的 /dev/log 这个 Socket。

## RsysLog 的的规范

RsysLog 使用 **Facility(设施)** 来对各个程序产生的日志进行分类好便于管理，每个 Facility 包含 1 个或多个程序，Facility 用于约束多个程序所产生的日志数据流到同一个管道内，默认有以下几个，括号中的数字与名称相对应

- **kern(0)** # 内核的日志。
- **user(1)** # 用户层日志，比如用户使用 logger 命令来记录日志功能。
- **mail(2)** # 邮件相关的日志。
- **daemon(3)** # 系统服务产生的信息，比如 systemd 管理的服务的信息。
- **authpriv(4)** # 认证相关的日志，比如 login、ssh、su 等需要账号密码的。
- **syslog(5)** # 由 syslog 相关协议产生的信息，就是 rsyslog 程序本身的日志信息。
- **lpr(6)** #打印相关的日志。
- **news(7)** # 新闻组服务器有关的日志。
- **uucp(8)**：
- **cron(9)** # 定时任务产生的日志。
- **authpriv(10)** # 与 auth 类似，更多的是记录账号私人的日志，包括 pam 模块运作的日志。
- **ftp(11) **# 与 ftp 相关的信息。
- **16 到 23.local0 到 local7** # 保留给本机用户自定义设施。比如可以把某些设施设置成 local0，然后供 RsysLog 收集日志

日志的级别：

- emerg(0)：错误信息。最严重日志等级，意味着系统将要宕机
- alert(1)：错误信息。比 emerg 等级轻
- crit(2)：错误信息。
- err(3)：错误信息。err 就是 error
- warn(4)：警告信息。可能有问题，但是还不至于影响到程序的运行。warn 就是 warnning
- notice(5)：基本信息。
- info(6)：基本信息。
- debug(7)：特殊的等级，用来 troubleshooting 时产生的日志
- none：特殊的等级。表示某个 Facility 不需要执行 Action。i.e.即不记录的级别

RsysLog 默认把日志保存在 /var/log/ 目录下的文件中，该目录下常见的日志文件有：

- messages # 几乎所有系统发生的信息都会记录在这个文件中
  - Ubuntu 发型版中是 syslog 文件
- boot.log #
- cron # 记录 crontab 执行的信息
- dmesg # 系统开机时内核检查过程所产生的各项信息
- maillog 与 mail/\* # 记录邮件的往来日志主要是 postfix(SMTP)与 dovecot(POP3)所产生的信息
- secure # 只要涉及到需要输入账号密码的软件，那么当登录时，会被记录在这个文件中。包括系统的 login 程序、su 和 sudo、ssh 等
  - Ubuntu 发型版中是 auth.og 文件
- lastlog # 记录系统上所有账号最近一次登录系统时的相关信息。lastlog 命令就是利用这个文件记录的信息来展示的
- wtmp 与 faillog # 记录正确登录系统的账号信息与错误登录时所使用的账号信息。last 命令时读取的 wtmp 中的内容

## 日志的格式

Linux 相关的日志格式一般为：

月 日 时:分:秒 主机名 程序名:事件内容

## 总结

随着时代的发展，各个应用程序大部分都通过各自的日志库，将日志直接写到磁盘上了~~

# Rsyslog 关联文件

**/etc/rsyslog.conf **# rsyslog 程序的基础配置文件

- **/etc/rsyslog.d/\*.conf **# rsyslog.conf 可以包含该目录下的配置文件。常用于定义单独程序的日志配置，以便日后方便管理

**/etc/sysconfig/rsyslog **# rsyslogd 运行时参数配置
**/dev/log** # 一个 Unix Domain Socket，rsyslogd 从这个 Socket 中读取日志消息。这是传统的日志服务 Socket。在 CentOS 8 及以后的版本中，该文件是一个指向 /run/systemd/journal/syslog 文件的软链接

- **/run/systemd/journal/syslog** # rsyslogd 会持续监听该 Socket，当有数据传入时，使用 recvmsg() 调用获取日志数据。
  - 这个文件是由 Systemd 提供的 Socket 文件，用以兼容传统日志服务。在 /etc/systemd/journald.conf 配置文件中，可以看到默认 ForwardToSyslog=yes 设置，即表示将自己的日志转发到 syslog 中。

**/var/log/\*** # 日志记录的位置。根据 rsyslog 程序的基础配置文件，各个 Linux 发行版的文件名也许不同，但是大体都差不多

- ./message # CentOS 发行版的绝大部分日志文件
  - ./syslog # Ubuntu 发型版的绝大部分日志文件
- ./secure # 所有 authpriv 设施的日志，比如 su、sudo、sshd 的登录信息等等。
- /var/log/maillog mail 记录
- /var/log/utmp
- /var/log/wtmp 登陆记录信息（last 命令即读取此日志）
