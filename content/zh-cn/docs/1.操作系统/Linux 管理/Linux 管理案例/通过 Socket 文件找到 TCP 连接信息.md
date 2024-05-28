---
title: 通过 Socket 文件找到 TCP 连接信息
---

# 概述

> 参考：[原文链接](https://www.cnblogs.com/web21/p/6520164.html)

# 进程的打开文件描述符表

[Linux Kernel](docs/1.操作系统/Kernel/Linux%20Kernel/Linux%20Kernel.md) 的三个系统调用：open，socket，pipe 返回的都是一个描述符。不同的进程中，他们返回的描述符可以相同。那么，在同一个进程中，他们可以相同吗？或者说，普通文件、套接字和管道，这三者的描述符属于同一个集合吗？

在内核源码中，三个系统调用声明如下：

```c
SYSCALL_DEFINE3(open, const char __user *, filename, int, flags, umode_t, mode);
SYSCALL_DEFINE3(socket, int, family, int, type, int, protocol);
SYSCALL_DEFINE1(pipe, int __user *, fildes);
```

他们都会先后调用函数

- get_unused_fd_flags：获取当前进程打开文件描述符表中的空闲描述符；
- fd_install：安装新描述符到当前进程打开文件描述符表。

内核为每个进程维护了一个结构体 struct task_struct，可称为进程表项、进程控制块（PCB: Process Control Block）或者进程描述符，定义如下：

```c
struct task_struct {
    volatile long state;  /* -1 unrunnable, 0 runnable,>0 stopped */
    …
    pid_t pid;
    …
    struct files_struct *files;
    …
};
```

其中 files 成员成为打开文件描述符表，定义如下：

```c
struct files_struct {
    …
    struct fdtable fdtab;
    …
    struct file __rcu * fd_array[NR_OPEN_DEFAULT];
};
```

其成员 fdtab 为关键数据成员，定义如下：

```c
struct fdtable {
    unsigned int max_fds;
    struct file __rcu **fd;      /* current fd array */
    unsigned long *close_on_exec;
    unsigned long *open_fds;
    struct rcu_head rcu;
};
```

这说明普通的文件、套接字、管道等，都被抽象为文件，共同占用进程的打开文件描述符。

http://blog.csdn.net/aprilweet/article/details/53482930

/Proc 目录下面有许多数字命名的子目录，这些数字表示系统当前运行的进程号；

其中/proc/N/fd 目录下面保存了打开的文件描述符，指向实际文件的一个链接。如下：

```bash
root@yang-ubuntu:/proc/4810/fd# ll
总用量 0
dr-x------ 2 root root 0 3月 8 16:07 ./
dr-xr-xr-x 8 root root 0 3月 8 16:07 ../
lrwx------ 1 root root 64 3月 8 16:08 0 -> /dev/pts/1
lrwx------ 1 root root 64 3月 8 16:08 1 -> /dev/pts/1
lrwx------ 1 root root 64 3月 8 16:09 10 -> socket:[21190]
lrwx------ 1 root root 64 3月 8 16:07 2 -> /dev/pts/1
lrwx------ 1 root root 64 3月 8 16:08 3 -> /tmp/ZCUDY7QsPB (deleted)
lrwx------ 1 root root 64 3月 8 16:08 4 -> /tmp/sess_0fpvhvcpftcme834e1l4beo2i6
lrwx------ 1 root root 64 3月 8 16:08 5 -> socket:[20625]
lrwx------ 1 root root 64 3月 8 16:08 6 -> anon_inode:[eventpoll]
lrwx------ 1 root root 64 3月 8 16:08 7 -> socket:[20626]
lrwx------ 1 root root 64 3月 8 16:08 8 -> socket:[20627]
lrwx------ 1 root root 64 3月 8 16:09 9 -> socket:[21189]
```

我们想查看 9 -> Socket 文件描述符的链接状态该怎么看呢？聪明的注意到后面有个数字\[21189]，这个数字又是哪儿来的呢？

在 `/proc/net/tcp` 目录下面保存了所有 TCP 链接的状态信息。

```bash
root@yang-ubuntu:/proc/net# vim /proc/net/tcp
sl local_address rem_address st tx_queue rx_queue tr tm->when retrnsmt uid timeout inode
0: 00000000:0CEA 00000000:0000 0A 00000000:00000000 00:00000000 00000000 1001 0 9482 1 ffff88001a501a00 100 0 0 10 -1
1: 00000000:008B 00000000:0000 0A 00000000:00000000 00:00000000 00000000 0 0 8916 1 ffff88001a501380 100 0 0 10 -1
2: 00000000:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000 0 0 11440 1 ffff88001a502080 100 0 0 10 -1
3: 0100007F:0035 00000000:0000 0A 00000000:00000000 00:00000000 00000000 0 0 12333 1 ffff88001a502700 100 0 0 10 -1
4: 00000000:0016 00000000:0000 0A 00000000:00000000 00:00000000 00000000 0 0 7922 1 ffff88001a500000 100 0 0 10 -1
5: 0100007F:0277 00000000:0000 0A 00000000:00000000 00:00000000 00000000 0 0 13302 1 ffff88001a500680 100 0 0 10 -1
6: 00000000:01BD 00000000:0000 0A 00000000:00000000 00:00000000 00000000 0 0 8914 1 ffff88001a500d00 100 0 0 10 -1
7: 00000000:0929 00000000:0000 0A 00000000:00000000 00:00000000 00000000 0 0 20625 1 ffff88001a504100 100 0 0 10 -1
8: 8064A8C0:01BD 0164A8C0:C26A 01 00000000:00000000 02:00030E57 00000000 0 0 13216 2 ffff88001a503a80 22 4 1 10 18
9: 8064A8C0:0929 0164A8C0:F4B5 01 00000000:00000000 02:00097B3E 00000000 0 0 21189 2 ffff88001a505b00 24 4 28 10 -1
10: 8064A8C0:0016 0164A8C0:CD9C 01 00000000:00000000 02:0000B4B4 00000000 0 0 17721 2 ffff88001a503400 22 4 20 10 -1
11: 8064A8C0:0016 0164A8C0:CDAE 01 00000000:00000000 02:0000DB1B 00000000 0 0 18130 2 ffff88001a504e00 24 4 31 10 -1
12: 8064A8C0:0929 0164A8C0:F4B6 01 00000000:00000000 02:00097B3E 00000000 0 0 21190 2 ffff88001a506800 24 4 24 10 -1
13: 8064A8C0:0016 0164A8C0:CDAC 01 00000000:00000000 02:0000DB1B 00000000 0 0 18074 2 ffff88001a502d80 21 4 24 10 -1
14: 8064A8C0:0016 0164A8C0:F3FC 01 00000000:00000000 02:00089B3B 00000000 0 0 20675 2 ffff88001a506180 24 4 25 10 -1
15: 8064A8C0:0016 0164A8C0:CDAD 01 00000080:00000000 01:00000018 00000000 0 0 18102 4 ffff88001a504780 24 4 21 10 -1
```

看上数字【21189 】没有，就是这儿来的，到此我们可以找出链接的 IP、PORT 链接四元组【8064A8C0:0929 0164A8C0:F4B5】这个地方是用十六进制保存的，换算成十进制方式【192.168.100.128:2345            192.168.100.1:62645】；

去网络连接状态里面看一下：

```bash
root@yang-ubuntu:/proc/4275/fd# netstat -antp
激活 Internet 连接 (服务器和已建立连接的)
Proto Recv-Q Send-Q Local Address Foreign Address State PID/Program name
tcp 0 0 0.0.0.0:3306 0.0.0.0:_ LISTEN 1710/mysqld
tcp 0 0 0.0.0.0:139 0.0.0.0:_ LISTEN 1062/smbd
tcp 0 0 0.0.0.0:80 0.0.0.0:_ LISTEN 1736/nginx.conf
tcp 0 0 127.0.0.1:53 0.0.0.0:_ LISTEN 1925/dnsmasq
tcp 0 0 0.0.0.0:22 0.0.0.0:_ LISTEN 628/sshd
tcp 0 0 127.0.0.1:631 0.0.0.0:_ LISTEN 709/cupsd
tcp 0 0 0.0.0.0:445 0.0.0.0:_ LISTEN 1062/smbd
tcp 0 0 0.0.0.0:2345 0.0.0.0:_ LISTEN 4809/start.php
tcp 0 0 192.168.100.128:445 192.168.100.1:49770 ESTABLISHED 2514/smbd
tcp 0 0 192.168.100.128:2345 192.168.100.1:62645 ESTABLISHED 4810/0.0.0.0:2345
tcp 0 0 192.168.100.128:22 192.168.100.1:52636 ESTABLISHED 3565/sshd: root@not
tcp 0 0 192.168.100.128:22 192.168.100.1:52654 ESTABLISHED 3718/3
tcp 0 0 192.168.100.128:22 192.168.100.1:52652 ESTABLISHED 3714/1
tcp 0 0 192.168.100.128:22 192.168.100.1:62460 ESTABLISHED 4817/4
tcp 0 0 192.168.100.128:22 192.168.100.1:52653 ESTABLISHED 3716/2
tcp6 0 0 :::139 :::_ LISTEN 1062/smbd
tcp6 0 0 :::22 :::_ LISTEN 628/sshd
tcp6 0 0 ::1:631 :::_ LISTEN 709/cupsd
tcp6 0 0 :::445 :::_ LISTEN 1062/smbd
```

回到开始的问题：9 Socket 文件描述符代表的是本地【192.168.100.128:2345】到【192.168.100.1:62645】的一条 TCP 连接！

为什么往 socket 中写数据，就会发送到对端(只针对 tcp 协议的研究)?  举例：浏览器请求服务器？

客户端首先发起建立与服务器 TCP 连接。一旦建立连接，浏览器进程和服务器进程就可以通过各自的套接字来访问 TCP。

客户端套接字是客户进程和 TCP 连接之间的“门”，服务器端套接字是服务器进程和同一 TCP 连接之间的“门”。

客户往自己的套接字发送 HTTP 请求消息，也从自己的套接字接收 HTTP 响应消息。

类似地，服务器从自己的套接字接收 HTTP 请求消息，也往自己的套接字发送 HTTP 响应消息。

客户端或服务器一旦把某个消息送入各自的套接字，这个消息就完全落入 TCP 的控制之中。

所以说底层是基于 tcp 提供的可靠的消息传输机制

TCP 给 HTTP 提供一个可靠的数据传输服务;这意味着由客户发出的每个 HTTP 请求消息最终将无损地到达服务器，由服务器发出的每个 HTTP 响应消息最终也将无损地到达客户。


