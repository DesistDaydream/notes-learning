---
title: pam_limits 模块
---

# 概述

> 参考：
> 
> - [Manual(手册)，pam_limits(8)](https://man7.org/linux/man-pages/man8/pam_limits.8.html)
> - [Manula(手册)，limits.conf(5)](https://www.man7.org/linux/man-pages/man5/limits.conf.5.html)

应用场景：我的 Linux 主机里面同时登入了十个人，这十个人不知怎么搞的， 同时开启了 100 个文件，每个文件的大小约 10M ，请问一下， 我的 Linux 主机的内存要有多大才够？ 10*100*10 = 10000M = 10G ... 老天爷，这样，系统不挂点才有鬼哩！为了要预防这个情况 的发生，所以我们是可以『限制用户的某些系统资源』的，包括可以开启的文件数量， 可以使用的 CPU 时间，可以使用的内存总量等等。除了这个模块有这个限制功能外，bash 还自带 ulimit 命令可以实现这个功能。

- 限制用户会话过程中对各种系统资源的使用情况，默认情况下该模块的配置文件是/etc/security/limits.conf
- 该配置文件中的参数，可以通过 `ulimit -a` 命令查看

# 关联文件与配置

**/etc/security/** #

- **/etc/security/limits.conf** # 限制用户会话过程中对各种系统资源的使用情况的配置文件

## 配置文件语法

**/etc/security/limit.conf 文件语法格式：**

    <domain>           <type>          <item>        <value>
      XXX              XXXX            XXXX           XXXX
    用户名/组名        限制类型       要限制的项目      具体值

- **domain** # 设置需要被限制的用户名或组名，组名前面加@和用户名区别。也可使用通配符 `*` 来表示所有用户
- **type** # 在设定上，通常 soft 会比 hard 小，举例来说，soft 可设定为 80 而 hard 设定为 100，那么你可以使用到 90 (因为没有超过 100)，但介于 80~100 之间时，系统会有警告讯息通知你！
    - hard # 严格的设定，指设定的 value 必定不能超过设定的数值
    - soft # 警告的设定，指设定的 value 可以超过设定值，但是若超过则有警告讯息。
- **item** # 指定要限制的项目
    - core # 限制内核文件的大小
    - 何谓 core 文件,当一个程序崩溃时，在进程当前工作目录的 core 文件中复制了该进程的存储图像。core 文件仅仅是一个内存映象（同时加上调试信息），主要是用来调试的。core 文件是个二进制文件，需要用相应的工具来分析程序崩溃时的内存映像，系统默认 core 文件的大小为 0，所以没有被创建。可以用 ulimit 命令查看和修改 core 文件的大小，例如：ulimit -c 1000 # 指定修改 core 文件的大小，1000 指定了 core 文件大小。也可以对 core 文件的大小不做限制，如： ulimit -c unlimited
    - date # 最大数据大小
    - fsize # 最大文件大小
    - memlock # 最大锁定内存地址空间
    - nofile # 打开文件的最大数目，默认为 1024
    - 对于需要做许多套接字连接并使它们处于打开状态的应用程序而言，最好通过使用 ulimit -n，或者通过设置 nofile 参数，为用户把文件描述符的数量设置得比默认值高一些
    - rss # 最大持久设置大小
    - stack # 最大栈大小
    - cpu # 以分钟为单位的最多 CPU 时间
    - nproc # 打开进程的最大数
    - as # 地址空间限制
    - maxlogins # 此用户允许登录的最大数目
- **value** # 指定 item 中具体项目的值
    - NUM # 可以是具体的数值
    - unlimited # 表示无限制的

## security 目录下文件的配置示例

/etc/security/limits.conf 配置文件示例

如图所示，左边的是 limit.conf 的配置，意思是所有用户的可以打开的文件最大数为 1000000，然后通过 bash 内嵌命令 ulimit -a 可以查看到右图的效果

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/xg1n26/1616166776345-13a61424-8c05-47a2-af4a-aac10869dcaa.png)

注意：在 CentOS 7 & Ubuntu 系统中，使用 Systemd 替代了之前的 SysV。/etc/security/limits.conf 文件的配置作用域缩小了。/etc/security/limits.conf 的配置，只适用于通过 PAM 认证登录用户的资源限制，它对 systemd 的 service 的资源限制不生效。因此登录用户的限制，通过 /etc/security/limits.conf 与 /etc/security/limits.d 下的文件设置即可。对于 systemd service 的资源设置，则需修改全局配置，全局配置文件放在 /etc/systemd/system.conf 和 /etc/systemd/user.conf。

DefaultLimitNOFILE=100000 # 对应 max openfile

DefaultLimitNPROC=65535 # 对应 max processes

# 命令行工具

## ulimit

> 参考：
> - [Manual(手册)，ulimit(1p)](https://man7.org/linux/man-pages/man1/ulimit.1p.html)

临时更改或者查看当前登录用户的资源限制情况。

注意：该命令只可在当前 shell 环境中生效，如果想要让配置永久生效，需要修改/etc/security/limit.conf 文件

OPTIONS

- **-a** # 显示所有限制 All current limits are reported
- -b # The maximum socket buffer size
- -c # core 文件大小的上限 The maximum size of core files created
- -d # 进程数据段大小的上限 The maximum size of a process's data segment
- -e # The maximum scheduling priority ("nice")
- -f # shell 所能创建的文件大小的上限 The maximum size of files written by the shell and its children
- -i # The maximum number of pending signals
- -l # The maximum size that may be locked into memory
- -m # 驻留内存大小的上限 The maximum resident set size (many systems do not honor this limit)
- **-n** # 打开文件数的上限 The maximum number of open file descriptors (most systems do not allow this value to be set)
- -p # 管道大小 The pipe size in 512-byte blocks (this may not be set)
- -q # The maximum number of bytes in POSIX message queues
- -r # The maximum real-time scheduling priority
- -s # 堆栈大小的上限 The maximum stack size
- -t # 每秒可占用的 CPU 时间上限 The maximum amount of cpu time in seconds
- -u # 进程数的上限 The maximum number of processes available to a single user
- -v # 虚拟内存的上限 The maximum amount of virtual memory available to the shell and, on some systems, to its children
- -x # The maximum number of file locks
- -T # The maximum number of threads

EXAMPLE

- ulimit -n 1000 #

# 优化实例

```bash
cat >> /etc/security/limits.d/max-openfile.conf  << EOF
* soft nproc 1000000
* hard nproc 1000000
* soft nofile 1000000
* hard nofile 1000000
EOF
```

cat /dev/urandom | head -1 | md5sum | head -c 32 >> /pass # 生成 32 位随机密码
