---
title: UTS Namespace
---

# 概述

**UTS(UNIX Time-Sharing System)** Namespace 可隔离 hostname 和 NIS Domain name 资源，使得一个宿主机可拥有多个主机名或 Domain Name。换句话说，可让不同 namespace 中的进程看到不同的主机名。

例如，使用 unshare 命令(较新版本 Linux 内核还支持 nscreate 命令)创建一个新的 uts namespace：

```bash
# -u或--uts表示创建一个uts namespace
# 这个namespace中运行/bin/bash程序
$ hostname
longshuai-vm      # 当前root namespace的主机名为longshuai-vm
$ sudo unshare -u /bin/bash
root@longshuai-vm:/home/longshuai#   # 进入了新的namespace中的shell
# 其主机名初始时也是longshuai-vm，
# 其拷贝自上级namespace资源
```

上面指定运行的是/bin/bash 程序，这会进入交互式模式，当执行 exit 时，bash 退出，回到当前的 namespace 中。也可以指定在 namespace 中运行其他程序，例如 unshare -u sleep 3 表示在 uts namespace 中睡眠 3 秒后退出并回到当前 namespace。

因为是 uts namespace，所以可在此 namespace 中修改主机名：

```bash
# 修改该namespace的主机名为ns1
# 修改后会立即生效，但不会显示在当前Shell提示符下
# 需重新加载Shell环境
root@longshuai-vm:/home/longshuai# hostname ns1
root@longshuai-vm:/home/longshuai# hostname
ns1
root@longshuai-vm:/home/longshuai# exec $SHELL
root@ns1:/home/longshuai#
```

namespace 中修改的主机名不会直接修改主机名配置文件(如/etc/hostname)，而是修改内核属性/proc/sys/kernel/hostname：

    root@ns1:/home/longshuai# cat /proc/sys/kernel/hostname
    ns1
    root@ns1:/home/longshuai# cat /etc/hostname
    longshuai-vm

创建了新的 namespace 并在其中运行/bin/bash 进程后，再去关注一下进程关系：

    # ns1中的bash进程PID
    root@ns1:/home/longshuai# echo $$
    14279

    # bash进程(PID=14279)和grep进程运行在ns1 namespace中，
    # 其父进程sudo(PID=14278)运行在ns1的上级namespace即root namespace中
    root@ns1:/home/longshuai# pstree -p | grep $$
        |-sshd(10848)---bash(10850)---sudo(14278)---bash(14279)-+-grep(14506)

    # 运行在ns1中当前bash进程(PID=14279)的namespace
    root@ns1:/home/longshuai# ls -l /proc/14279/ns
    lrwxrwxrwx ... cgroup -> 'cgroup:[4026531835]'
    lrwxrwxrwx ... ipc -> 'ipc:[4026531839]'
    lrwxrwxrwx ... mnt -> 'mnt:[4026531840]'
    lrwxrwxrwx ... net -> 'net:[4026531992]'
    lrwxrwxrwx ... pid -> 'pid:[4026531836]'
    lrwxrwxrwx ... pid_for_children -> 'pid:[4026531836]'
    lrwxrwxrwx ... user -> 'user:[4026531837]'
    lrwxrwxrwx ... uts -> 'uts:[4026532588]'  # 注意这一行，和sudo进程的uts inode不同

    # 父进程sudo(PID=14278)不在ns1中，它的namespace信息
    root@ns1:/home/longshuai# ls -l /proc/14278/ns
    lrwxrwxrwx ... cgroup -> 'cgroup:[4026531835]'
    lrwxrwxrwx ... ipc -> 'ipc:[4026531839]'
    lrwxrwxrwx ... mnt -> 'mnt:[4026531840]'
    lrwxrwxrwx ... net -> 'net:[4026531992]'
    lrwxrwxrwx ... pid -> 'pid:[4026531836]'
    lrwxrwxrwx ... pid_for_children -> 'pid:[4026531836]'
    lrwxrwxrwx ... user -> 'user:[4026531837]'
    lrwxrwxrwx ... uts -> 'uts:[4026531838]'   # 注意这一行，和PID=1的uts inode相同

回到创建 uts namespace 时敲下的 unshare 命令：

    sudo unshare -u /bin/bash

1
Shell

从进程关系...---sudo(14278)---bash(14279)可知两个进程 PID 是连续的，说明 unshare 程序对应的进程被/bin/bash 程序通过 execve()替换了。

详细的过程如下：「sudo 进程运行在当前 namespace 中，它将 fork 一个新进程来运行 unshare 程序，unshare 程序加载完成后，将创建一个新的 uts namespace，unshare 进程自身将加入到这个 uts namespace 中，unshare 进程内部再 exec 加载/bin/bash，于是 unshare 进程被替换为/bin/bash 进程，/bin/bash 进程也将运行在 uts namespace 中」。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/lu9a00/1616122843124-f5e9da25-edb1-45f4-8eb3-58a15fb18d89.png)

当 namespace 中的/bin/bash 进程退出，该 namespace 中将没有任何进程，该 namespace 将自动销毁。注意，在默认情况下，namespace 中必须要有至少一个进程，否则将被自动被销毁。但也有一些手段可以让 namespace 持久化，即使已经没有任何进程在其中运行。

如果在 ns1 中再创建一个 namespace ns2，这个 ns2 初始时将共享 ns1 的其他资源并拷贝 ns1 的主机名资源，其初始主机名也为 ns1。

    $ sudo unshare -u /bin/bash    # 在root namespace环境下创建一个namespace
    root@longshuai-vm:/home/longshuai# hostname ns1 # 修改主机名为ns1
    root@longshuai-vm:/home/longshuai# hostname
    ns1
    # 在ns1中创建一个namespace
    ############ 注意没有sudo
    root@longshuai-vm:/home/longshuai# unshare -u /bin/bash
    root@ns1:/home/longshuai# hostname    # 初始主机名拷贝自上级namespace的主机名ns1
    ns1
    root@ns1:/home/longshuai# hostname ns2
    root@ns1:/home/longshuai# hostname  # 修改主机名为ns2
    ns2
    root@ns1:/home/longshuai# exit
    exit
    root@longshuai-vm:/home/longshuai# hostname  # ns2修改主机名不影响ns1
    ns1
    root@longshuai-vm:/home/longshuai# exit
    exit
    [~]->$ hostname      # ns1修改主机名不影响root namespace
    longshuai-vm

注意，即使 root namespace 当前用户为 longshuai，但因为使用了 sudo 创建 ns1，进入 ns1 后其用户名为 root，所以在 ns1 中执行 unshare 命令创建新的 namespace 不需要再使用 sudo。

    $ echo $USER      # 当前root namespace的用户为longshuai
    longshuai
    $ sudo unshare -u /bin/bash
    root@longshuai-vm:/home/longshuai# echo $USER  # ns中的用户名变为root
    root
    root@longshuai-vm:/home/longshuai# id;echo $HOME;echo ~
    uid=0(root) gid=0(root) groups=0(root)
    /root
    /root
