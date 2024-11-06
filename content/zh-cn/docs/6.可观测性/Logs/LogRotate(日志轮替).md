---
title: LogRotate(日志轮替)
weight: 3
---

# 概述

> 参考：
>
> - [Wiki, Log_rotation](https://en.wikipedia.org/wiki/Log_rotation)
> - [公众号-马哥 Linux 运维，\[译\] 理解 logrotate 实用工具](https://mp.weixin.qq.com/s/b_CWt_ycvnbQG9TXPqRoCQ)

为了防止日志文件持续被写入文件导致过于庞大，那么就需要对日志进行拆分，每隔一段时间就把日志文件保存(打包压缩)起来，然后再创建一个新的空文件继续接收日志，来回循环该过程直到通过配置规定的保留日期，来清除存在过久的日志。通过这种方式来进行日志的归档，分类，清理。这就是 LogRotate 所做的事情。是否进行轮替会有一系列的配置，比如文件的大小达到 N 会轮替一次，每隔多少天轮替一次等等。

logrotate 只是一个命令行工具，不以守护进程的方式运行在后台，默认情况下，logrotate 命令作为放在 /etc/cron.daily 中的 cron 任务，每天运行一次，该任务会根据设置的策略进行日志文件的检查，其中达到设置中满足轮替配置的日志文件将被轮替。

# 关联文件与配置

**/etc/logrotate.conf** # logrotate 基本配置文件

**/etc/logrotate.d/** # 对基本文件的扩展，该目录下的文件的配置会被包含在基本配置文件中。该目录下一般是一个程序一个文件，每个程序都有自己的日志轮替配置。

**/etc/cron.daily/logrotate** # 该文件定义了 cron 定时任务执行日志轮替工作的时间

**/var/lib/logrotate.status** # logrotate 的执行历史

## logrotate.conf 配置文件详解

- /PATH/TO/FILES {...} # 指定想要轮替的日志文件，可以通过 `*` 通配指定多个文件名
  - **copytruncate** # 把正在输出的日志拷(copy)一份出来，再清空(trucate)原来的日志。
  - **compress** # 压缩日志文件的所有非当前版本
  - **dateext** # 切换后的日志文件会附加上一个短横线和 YYYYMMDD 格式的日期,
  - **daily** # 日志文件将每天轮替一次。其它可用值为 monthly(每月)，weekly(每周)、yearly(每年)
  - **delaycompress** # 在轮替任务完成后，已轮替的归档将使用 gzip 进行压缩
  - **errors \<EMAIL>** # 给指定邮箱发送错误通知
  - **missingok** # 如果日志文件丢失，不要显示错误
  - **notifempty** # 如果日志文件为空，则不轮换日志文件
  - **olddir \<DIR>** # 指定日志文件的旧版本放在 “DIR”目录 中
  - **postrotate 和 endscript** # 在所有其它指令完成后，postrotate 和 endscript 里面指定的命令将被执行。在这种情况下，rsyslogd 进程将立即再次读取其配置并继续运行。
  - **rotate N** # 共存储 N 个轮替后日志。当产生第 N+1 个轮替后的日志，时间最久的日志将被删除
  - **sharedscripts** # 有多个日志需要轮替时，只执行一次脚本
  - **size \<LogSize>** # 在日志文件大小大于 LogSize（例如 100K，4M）时进行轮替

配置样例

```text
/var/log/nginx/*log {
    daily
    rotate 10
    missingok
    notifempty
    compress
    dateext
    sharedscripts
    postrotate
        /bin/kill -USR1 $(cat /var/run/ngnix/nginx.pid 2>/dev/null) 2>/dev/null
    endscript
}
```

Note：关于 postrotate

postrotate 后面跟随的是一个命令行，一般是用来重新生成日志文件或者冲定义应用所指向的文件描述符（fd：file description），拿 nginx 和 uwsgi 为例：

完成日志切割后创建新的 nginx 日志文件，此时该文件的 fd 发生改变

nginx 中日志输出对应的文件 fd 未同步更新，nginx 会向原 fd 对应的日志文件写数据

“/bin/kill -USR1 cat /var/run/nginx.pid || true”，更新 nginx 默认日志文件的 fd 到新建的日志文件（该效果等同于 reload）。

关于/bin/kill -HUP

"/bin/kill -USR1 `cat /var/run/nginx.pid` || true"

看到这条命令很容易想到：

/bin/kill -HUP `cat /var/run/nginx.pid 2> /dev/null` 2> /dev/null || true

这两条命令的大概含义是重载 nginx 服务，目的是重新生成 nginx 的日志文件。

# 命令行工具

**logrotate \[OPTIONS...]**

OPTIONS

- **-f** # 告诉 logrotate 强制执行轮替，即使这不是必要的(i.e.测试轮替的配置文件是否可以正常运行)。 有时，在向 logrotate 配置文件添加新条目之后，或者如果已经手动删除旧的日志文件，这将是有用的，因为将创建新文件，并且日志记录将正常继续。

EXAMPLE

- logrotate -f /etc/logrotate.d/keepalived # 使用/etc/logrotate.d/keepalived 配置文件执行轮替

PS:

- 遇到不能记录日志的情况：kill -USR1 pid 重发信号量
