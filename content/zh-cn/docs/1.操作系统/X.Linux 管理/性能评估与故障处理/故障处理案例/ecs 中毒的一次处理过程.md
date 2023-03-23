---
title: ecs 中毒的一次处理过程
---

原文：[张馆长博客，ecs 中毒的一次处理过程](https://zhangguanzhang.github.io/2022/04/21/ecs-xmrig/)
一次客户 ecs 中毒的处理过程，可以给读者参考下中毒的处理过程。

## 由来

客户机器中毒了，pm 找我来让处理下，记录下，给其他人做个处理过程的参考。

## 处理过程

机器是 centos ，先利用 `rpm -V <pkg_name>` 确认基础的排查命令没被修改过：

```bash
$ rpm -qf `which ps`
procps-ng-3.3.10-23.el7.x86_64
$ rpm -V procps-ng
$ rpm -qf `which top`
procps-ng-3.3.10-23.el7.x86_64
```

top 看到异常 cpu 的进程占用 cpu 很高：

    $ top
    top - 19:44:29 up 34 days,  5:08,  4 users,  load average: 612.03, 617.15, 482.75
    Tasks: 2014 total,  66 running, 1946 sleeping,   0 stopped,   2 zombie
    %Cpu(s): 96.6 us,  3.1 sy,  0.0 ni,  0.0 id,  0.0 wa,  0.0 hi,  0.3 si,  0.0 st
    KiB Mem : 13186040+total,  2722452 free, 48820448 used, 80317504 buff/cache
    KiB Swap:        0 total,        0 free,        0 used. 78946784 avail Mem

      PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND
     1206 root      20   0 5251748   2.3g   3584 S  2956  1.8 465:37.77 ld-linux-x86-64

给它 `STOP` 信号不让 cpu 切换到它，而不是直接 kill 掉它：

    $ kill -STOP 1206

查看来源和清理：

    $ ll /proc/1206/exe
    lrwxrwxrwx 1 root root 0 Apr 21 19:44 /proc/1206/exe -> /dev/shm/.x/stak/ld-linux-x86-64.so.2

### 清理定时任务

排查定时任务，发现有内容，清理掉， crond 的子目录也看下，文件内容和多了的子文件也处理下

    $ crontab -l
    * * * * * /dev/shm/.x/upd >/dev/null 2>&1
    @reboot /dev/shm/.x/upd >/dev/null 2>&1

    [root@p-96b7-dndl etc]# find /etc/cron.* -type f
    /etc/cron.d/0hourly
    /etc/cron.d/sysstat
    /etc/cron.daily/logrotate
    /etc/cron.daily/man-db.cron
    /etc/cron.deny
    /etc/cron.hourly/0anacron

查看下进程树，是否有父进程拉起 `1206`:

    $ pstree -sp 1206
    systemd(1)───ld-linux-x86-64(1206)─┬─{ld-linux-x86-64}(1209)
                                       ├─{ld-linux-x86-64}(1211)
                                       ├─{ld-linux-x86-64}(1216)
                                       ├─{ld-linux-x86-64}(1217)
                                       ├─{ld-linux-x86-64}(1218)
                                       ├─{ld-linux-x86-64}(6436)
                                       ├─{ld-linux-x86-64}(6437)
                                       ├─{ld-linux-x86-64}(6439)
                                       ├─{ld-linux-x86-64}(6440)
                                       ├─{ld-linux-x86-64}(6441)
                                       ├─{ld-linux-x86-64}(6443)
                                       ├─{ld-linux-x86-64}(6471)
                                       ├─{ld-linux-x86-64}(6472)
                                       ├─{ld-linux-x86-64}(6476)
                                       ├─{ld-linux-x86-64}(6484)
                                       ├─{ld-linux-x86-64}(6489)
                                       ├─{ld-linux-x86-64}(6495)
                                       ├─{ld-linux-x86-64}(6501)
                                       ├─{ld-linux-x86-64}(6504)
                                       ├─{ld-linux-x86-64}(6505)
                                       ├─{ld-linux-x86-64}(6508)
                                       ├─{ld-linux-x86-64}(6509)
                                       ├─{ld-linux-x86-64}(6511)
                                       ├─{ld-linux-x86-64}(6523)
                                       ├─{ld-linux-x86-64}(6527)
                                       ├─{ld-linux-x86-64}(6529)
                                       ├─{ld-linux-x86-64}(6531)
                                       ├─{ld-linux-x86-64}(6535)
                                       ├─{ld-linux-x86-64}(6547)
                                       ├─{ld-linux-x86-64}(6554)
                                       ├─{ld-linux-x86-64}(6563)
                                       ├─{ld-linux-x86-64}(6567)
                                       ├─{ld-linux-x86-64}(6568)
                                       ├─{ld-linux-x86-64}(6569)
                                       ├─{ld-linux-x86-64}(6572)
                                       ├─{ld-linux-x86-64}(6579)
                                       └─{ld-linux-x86-64}(6580)

发现并没有，查看下进程的 `cmdline`:

    $ xargs -0 < /proc/1206/cmdline
    xmrig
        --library-path stak stak/xmrig -o 185.82.200.52:443 -k

### 检查系统的 so 和开机启动项

搜了下这个 ip 是外国的，查看下 ld 的 so 导入配置文件，看看是否有被加入额外的 so 导入：

    $ rpm -qf /etc/ld.so.conf
    glibc-2.17-260.el7.x86_64
    # glibc 也提供了很多基础的 so，这部同时也可以看出来
    # 基础的 so 有被替换不
    $ rpm -V glibc

同理查看下 systemd 的

    $ rpm -V systemd
    .M.......  c /etc/machine-id
    SM5....T.  c /etc/rc.d/rc.local # 这个文件也记得查下
    S.5....T.  c /etc/systemd/system.conf
    .M.......  g /etc/udev/hwdb.bin
    .M.......  g /var/lib/systemd/random-seed
    .M....G..  g /var/log/journal
    .M....G..  g /var/log/wtmp
    .M....G..  g /var/run/utmp
    # 查看下有没有被添加 systemd 的开机启动任务
    $ systemctl list-units

### 清理进程相关

我们环境是 k8s 和 docker 的，etcd 没证书，kubelet 的 http 可写，docker 开网络端口不 tls，redis 无密码这种现象是不存在的。看了下我们配置的部署配置文件，初步怀疑是一个有 sudo 的弱密码用户被爆破导致的中毒，查看了具有 sudo 权限和 root 的 `~/.ssh/authorized_keys` 也没被添加别人的公钥（有的话记得清理下），开始删除挖矿进程的目录：

    rm -rf /dev/shm/.x/
    kill -9 1206

### 排查网络

看看是否还有其他后台进程上报或者下载的，看了下 udp 的正常，所以提取所有活跃的 tcp 连接 ip 看看有异常的 IP 没：

    $ netstat -ant |& grep -Po '(\d{1,3}\.){3}\d{1,3}' | sort | grep -v 10.187.0 | uniq -c
         49 0.0.0.0
          2 119.82.135.65
        111 127.0.0.1
          4 169.254.169.254
       2271 192.168.0.235
         13 2xx.1xx.15.161
          1 3x.1xx.2x.7
         27 4x.x.1xx.x3
          1 xx.1xx.6x.x54

    $ netstat -ant | grep 119.82.135.65
    tcp        0   1281 192.168.0.235:22        119.82.135.65:38525     LAST_ACK
    tcp        0      1 192.168.0.235:22        119.82.135.65:54598     LAST_ACK
    $ lsof -nPi :38525
    $ lsof -nPi :54598
    $

看了下只有不断被外国 IP 暴力 ssh 的 IP，其余几个 IP 是我和客户那边的人员 IP。让客户改密码后再观察下。
