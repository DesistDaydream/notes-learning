---
title: mount Unit
---

# 概述

> 参考：
> - [Manual(手册),systemd-mount(5)](https://man7.org/linux/man-pages/man5/systemd.mount.5.html)
> - [张馆长博客,fstab 与 systemd.mount 自动挂载的一点研究和见解](https://zhangguanzhang.github.io/2019/01/30/fstab/)

所有以 `.mount` 结尾的 Unit 都是由 Systemd 控制和监督的文件系统挂载功能。该功能可以代替传统的 /etc/fstab 文件。

# 张馆长文章

x-systemd.automount 属于 systemd.mount，systemd 造了好多轮子，什么 crontab、网络管理器、日志服务 它都想给接替了。fstab 也是这样，systemd 引入了 .mount 单元这么个东西，用于控制文件系统挂载。

`defaults`下有`auto`会被开机挂载，`noauto`一般是和`x-systemd.automount`配合使用。而`x-systemd.automount`属于 systemd.mount，systemd 造了好多轮子，什么 crontab、网络管理器、日志服务 它都想给接替了。fstab 也是这样，systemd 引入了 .mount 单元这么个东西，用于控制文件系统挂载。

其实现在很多发行版都开始慢慢抛弃 fstab 了，先看一个 centos6 在 init 下的 fstab

    [root@APP-SRV-001 ~]# cat /etc/fstab
    #
    # /etc/fstab
    # Created by anaconda on Mon Nov 26 22:13:02 2018
    #
    # Accessible filesystems, by reference, are maintained under '/dev/disk'
    # See man pages fstab(5), findfs(8), mount(8) and/or blkid(8) for more info
    #
    UUID=19698973-561a-4b5b-aded-f6092bd1f341 /                       ext4    defaults        1 1
    UUID=1f808272-972d-416a-af7f-d3c88b16f434 /boot                   ext4    defaults        1 2
    UUID=C6D4-B036          /boot/efi               vfat    umask=0077,shortname=winnt 0 0
    UUID=7877eff3-a174-4cb5-9024-43be4fab35b2 swap                    swap    defaults        0 0
    tmpfs                   /dev/shm                tmpfs   defaults        0 0
    devpts                  /dev/pts                devpts  gid=5,mode=620  0 0
    sysfs                   /sys                    sysfs   defaults        0 0
    proc                    /proc                   proc    defaults        0 0

下面是 centos7 在 systemd 下的 fstab

    [root@CentOS76 ~]# cat /etc/fstab
    #
    # /etc/fstab
    # Created by anaconda on Mon Dec 17 01:38:08 2018
    #
    # Accessible filesystems, by reference, are maintained under '/dev/disk'
    # See man pages fstab(5), findfs(8), mount(8) and/or blkid(8) for more info
    #
    /dev/mapper/centos-root /                       xfs     defaults        0 0
    UUID=cec0b2aa-7695-4d24-a641-5b3ae111500a /boot                   xfs     defaults        0 0
    #/dev/mapper/centos-swap swap                    swap    defaults        0 0

我们发现很多东西在 fstab 里消失了，但是 centos7 上 mount 看实际上还是挂载了

    [root@CentOS76 ~]# mount
    sysfs on /sys type sysfs (rw,nosuid,nodev,noexec,relatime)
    proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
    devtmpfs on /dev type devtmpfs (rw,nosuid,size=995764k,nr_inodes=248941,mode=755)
    securityfs on /sys/kernel/security type securityfs (rw,nosuid,nodev,noexec,relatime)
    tmpfs on /dev/shm type tmpfs (rw,nosuid,nodev)
    devpts on /dev/pts type devpts (rw,nosuid,noexec,relatime,gid=5,mode=620,ptmxmode=000)
    tmpfs on /run type tmpfs (rw,nosuid,nodev,mode=755)
    tmpfs on /sys/fs/cgroup type tmpfs (ro,nosuid,nodev,noexec,mode=755)
    configfs on /sys/kernel/config type configfs (rw,relatime)
    /dev/mapper/centos-root on / type xfs (rw,relatime,attr2,inode64,noquota)
    systemd-1 on /proc/sys/fs/binfmt_misc type autofs (rw,relatime,fd=22,pgrp=1,timeout=0,minproto=5,maxproto=5,direct,pipe_ino=21472)
    hugetlbfs on /dev/hugepages type hugetlbfs (rw,relatime,pagesize=2M)
    debugfs on /sys/kernel/debug type debugfs (rw,relatime)
    mqueue on /dev/mqueue type mqueue (rw,relatime)
    /dev/sda1 on /boot type xfs (rw,relatime,attr2,inode64,noquota)
    tmpfs on /run/user/0 type tmpfs (rw,nosuid,nodev,relatime,size=201784k,mode=700)

其实是已经被 systemd 挂载了，我们可以查看 systemd 的挂载单元 tmp

    [root@CentOS76 ~]# systemctl cat tmp.mount
    # /usr/lib/systemd/system/tmp.mount
    #  This file is part of systemd.
    #
    #  systemd is free software; you can redistribute it and/or modify it
    #  under the terms of the GNU Lesser General Public License as published by
    #  the Free Software Foundation; either version 2.1 of the License, or
    #  (at your option) any later version.
    [Unit]
    Description=Temporary Directory
    Documentation=man:hier(7)
    Documentation=http://www.freedesktop.org/wiki/Software/systemd/APIFileSystems
    ConditionPathIsSymbolicLink=!/tmp
    DefaultDependencies=no
    Conflicts=umount.target
    Before=local-fs.target umount.target
    [Mount]
    What=tmpfs
    Where=/tmp
    Type=tmpfs
    Options=mode=1777,strictatime
    # Make 'systemctl enable tmp.mount' work:
    [Install]
    WantedBy=local-fs.target

内容也很好看懂。Requires、After 这和 systemd service 的写法基本一致，\[Mount] 下就是挂载的说明了，What 后是要挂载的文件系统，Where 是挂载到的地方，Type 是文件系统类型，Options 是挂载参数。相当于：

    mount -t <Type> -o <Options> <What> <Where>

我们看看目录下的挂载单元

    $ ll /usr/lib/systemd/system/*.mount
    -rw-r--r--. 1 root root 670 Oct 31 07:31 /usr/lib/systemd/system/dev-hugepages.mount
    -rw-r--r--. 1 root root 590 Oct 31 07:31 /usr/lib/systemd/system/dev-mqueue.mount
    -rw-r--r--. 1 root root 615 Oct 31 07:31 /usr/lib/systemd/system/proc-sys-fs-binfmt_misc.mount
    -rw-r--r--. 1 root root 681 Oct 31 07:31 /usr/lib/systemd/system/sys-fs-fuse-connections.mount
    -rw-r--r--. 1 root root 719 Oct 31 07:31 /usr/lib/systemd/system/sys-kernel-config.mount
    -rw-r--r--. 1 root root 662 Oct 31 07:31 /usr/lib/systemd/system/sys-kernel-debug.mount
    -rw-r--r--. 1 root root 703 Oct 31 07:31 /usr/lib/systemd/system/tmp.mount

我们看一个文件

    $ cat  /usr/lib/systemd/system/proc-sys-fs-binfmt_misc.mount
    #  This file is part of systemd.
    #
    #  systemd is free software; you can redistribute it and/or modify it
    #  under the terms of the GNU Lesser General Public License as published by
    #  the Free Software Foundation; either version 2.1 of the License, or
    #  (at your option) any later version.
    [Unit]
    Description=Arbitrary Executable File Formats File System
    Documentation=https://www.kernel.org/doc/Documentation/admin-guide/binfmt-misc.rst
    Documentation=http://www.freedesktop.org/wiki/Software/systemd/APIFileSystems
    DefaultDependencies=no
    [Mount]
    What=binfmt_misc
    Where=/proc/sys/fs/binfmt_misc
    Type=binfmt_misc

实际上挂载单元的文件名和挂载点是关联的，挂载点路径转换为去掉以一个斜线，所有斜线转成横线

例如它这个`/proc/sys/fs/binfmt_misc`转成`proc-sys-fs-binfmt_misc.mount`，注意，mount 单元不能从模版实例化而来， 也不能通过创建软连接的方法给同一个 mount 单元赋予多个别名。

在系统运行时创建的挂载点(独立于单元文件与 /etc/fstab 之外)将同样被 systemd 监控， 并且看上去与其他常规的 mount 单元没啥差别。详见 proc(5) 手册中对 /proc/self/mountinfo 的解释。

注意，某些 虚拟文件系统 拥有特别的功能，例如： /sys, /proc, /dev, /tmp, sys/fs/cgroup, /dev/mqueue, /proc/sys/fs/binfmt_misc … 无法通过 mount 单元对其进行修改或禁用。

mount 单元既可以通过单元文件进行配置， 也可以通过 /etc/fstab 文件(参见 fstab(5) 手册)进行配置。 /etc/fstab 中的挂载点将在每次重新加载 systemd 配置时(包括系统启动时) 动态的自动转化为 mount 单元。 一般来说，通过 /etc/fstab 配置挂载点是首选的方法， 详见 systemd-fstab-generator(8) 手册。

可以在 /etc/fstab 中使用一些无法被其他程序识别的 systemd 专用挂载选项， 以帮助创建与挂载点相关的依赖关系。 对于本地文件系统挂载点， systemd 将会自动在 local-fs.target 中创建指向此挂载点的 Wants= 或 Requires= 依赖； 对于网络文件系统挂载点， systemd 将会自动在 remote-fs.target 中创建指向此挂载点的 Wants= 或 Requires= 依赖； 至于究竟是 Wants= 还是 Requires= 依赖， 则取决于是否设置了下面的 nofail 选项。

- x-systemd.requires=设置一个到其他单元(例如 device 或 mount 单元)的 Requires= 与 After= 依赖(详见 systemd.unit(5) 手册)， 参数必须是以下三者之一： 一个单独的单元名称、 一个以绝对路径表示的设备节点、 一个以绝对路径表示的挂载点。 可以多次使用此选项以指定对多个单元的依赖。此选项对于如下两种挂载点特别有用： (1)需要额外辅助设备的，例如将日志存储在其他设备上的日志文件系统 (2)需要事先存在其他挂载点的，例如可以融合多个挂载点的叠合文件系统 (Overlay Filesystem)。
- x-systemd.before=, x-systemd.after=设置一个到其他单元(例如 mount 单元)的 Before= 或 After= 依赖。 参数必须是一个单独的单元名称或者一个以绝对路径表示的挂载点。 可以多次使用这些选项以指定对多个单元的依赖。 这些选项对于同时具备如下特征的挂载点特别有用： 挂载点带有 nofail 标记、 以异步方式挂载(async)(默认即是异步)、 需要在特定单元启动之前或之后才能挂载(例如必须在启动 local-fs.target 之前挂载)。 关于 Before= 与 After= 的解释，详见 systemd.unit(5) 手册。
- x-systemd.requires-mounts-for=设置一个到其他挂载点的 RequiresMountsFor= 依赖(详见 systemd.unit(5) 手册)。 参数必须是一个以绝对路径表示的挂载点， 可以多次使用此选项以指定对多个挂载点的依赖。
- x-systemd.device-bound 将文件系统所在的块设备升级为 BindsTo= 依赖。 此选项仅在使用 mount(8) 手动挂载时才有意义，因为此时默认为 Requires= 依赖。 注意，对于 /etc/fstab 中的条目或 mount 单元来说， 已经自动隐含的设置了此选项。
- x-systemd.automount 同时创建一个对应的 automount 单元， 详见 systemd.automount(5) 手册。
- x-systemd.idle-timeout=设置对应的 automount 单元的最大闲置时长。详见 systemd.automount(5) 手册中的 TimeoutIdleSec= 选项。
- x-systemd.device-timeout=设置等候所依赖的设备进入可用状态的最大时长，若超时则放弃挂载。 可以使用 “ms”, “s”, “min”, “h” 这样的时间单位后缀。 若省略后缀则表示单位是秒。

注意，此选项仅可用于 /etc/fstab 文件， 不可用于单元文件中的 Options= 选项。

- x-systemd.mount-timeout=设置 /etc/fstab 中挂载点的超时时长。 如果超过指定的时间仍未挂载成功，那么将会放弃该挂载点。 可以使用 “s”, “min”, “h”, “ms” 这样的时间单位后缀。 若省略后缀则表示单位是秒。

注意，此选项仅可用于 /etc/fstab 文件， 不可用于单元文件中的 Options= 选项。

参见下文的 TimeoutSec= 选项。

- \_netdev 通常根据文件系统的类型来判断是否为网络文件系统(例如 xfs, ext4 是本地，而 cifs, nfs 则是网络)， 以决定是否必须在网络可用之后才能启动。 但是在某些情况下这种判断是不可靠的(例如对于 iSCSI 这样基于网络的块设备)， 使用此选项之后就可以强迫将此文件系统标记为网络文件系统。

网络文件系统所对应的 mount 单元将被安排在 remote-fs-pre.target 与 remote-fs.target 之间启动(而不是在 local-fs-pre.target 与 local-fs.target 之间)， 并且自动获得 After=network-online.target, Wants=network-online.target, After=network.target 依赖。

- noauto, autonoauto 表示不将此挂载点加入到 local-fs.target/remote-fs.target 的依赖中， 也就是不在系统启动时自动挂载(除非为了满足其他单元的依赖而被挂载)。 auto(默认值) 表示自动将此挂载点加入到 local-fs.target/remote-fs.target 的依赖中， 也就是在系统启动时自动挂载。
- nofail 表示仅在 local-fs.target/remote-fs.target 中对此挂载点使用 Wants= 依赖(而不是默认的 Requires=)。 也就是即使此挂载点挂载失败，也不会中断系统的启动流程。
- x-initrd.mount 要在 initramfs 中额外挂载的文件系统，参见 systemd.special(7) 手册中对 initrd-fs.target 的解释。

如果一个挂载点既被封装到了一个 mount 单元中， 又被配置到了 /etc/fstab 文件中，那么： (1)如果单元文件位于 /usr 中， 那么以 /etc/fstab 文件为准(无视单元文件)。 (2)如果单元文件位于 /etc 中， 那么以单元文件为准(无视 /etc/fstab 文件)。 这意味着对于同一个挂载点来说， /etc 中的单元文件优先级最高、/etc/fstab 文件次之、 /usr 中的单元文件优先级最低。

选项

每个 mount 单元文件都必须包含一个 \[Mount] 小节， 用于包含该单元封装的挂载点信息。 可在 \[Mount] 小节中使用的选项， 有许多是与其他单元共享的，详见 systemd.exec(5) 与 systemd.kill(5) 手册。 这里只列出仅能用于 \[Mount] 小节的选项(亦称”指令”或”属性”)：

- What=绝对路径形式表示的被挂载对象：设备节点、LOOP 文件、其他资源(例如网络资源)。 详见 mount(8) 手册。 如果是一个设备节点，那么将会自动添加对此设备节点单元的依赖(参见 systemd.device(5) 手册)。 这是一个必需的设置。注意，因为可以在此选项中使用 “%” 系列替换标记， 所以百分号(%)应该使用 “%%” 表示。
- Where=绝对路径形式表示的挂载点目录。 注意，不可设为一个软连接(即使它实际指向了一个目录)。 如果挂载时此目录不存在，那么将尝试创建它。 注意，这里设置的绝对路径必须与单元文件的名称相对应(见上文)。 这是一个必需的设置。
- Type=字符串形式表示的文件系统类型。详见 mount(8) 手册。这是一个可选的设置。
- Options=一组逗号分隔的挂载选项。详见 mount(8) 手册。 这是一个可选的设置。注意，因为可以在此选项中使用 “%” 系列替换标记， 所以百分号(%)应该使用 “%%” 表示。
- SloppyOptions=设为 yes 表示允许在 Options= 中使用文件系统不支持的挂载选项， 且不会导致挂载失败(相当于使用了 mount(8) 的 -s 命令行选项)。 默认值 no 表示禁止在 Options= 中使用文件系统不支持的挂载选项(会导致挂载失败)。
- LazyUnmount=设置是否使用延迟卸载。 设为 yes 表示立即将文件系统从当前的挂载点分离， 但是一直等待到设备不再忙碌的时候， 才会清理所有对此文件系统的引用(也就是真正完成卸载)。 这相当于使用 umount(8) 的 -l 选项进行卸载。 默认值为 no
- ForceUnmount=设置是否使用强制卸载。 设为 yes 表示使用强制卸载(仅建议用于 NFS 文件系统)。 这相当于使用 umount(8) 的 -f 选项进行卸载。 默认值为 no
- DirectoryMode=自动创建挂载点目录(包括必要的上级目录)时， 所使用的权限模式(八进制表示法)。 默认值是 0755
- TimeoutSec=最大允许使用多长时间以完成挂载动作。 若超时则被视为挂载失败， 并且所有当前正在运行的命令都将被以 SIGTERM 信号终止； 若继续等待相同的时长之后命令仍未终止， 那么将使用 SIGKILL 信号强制终止 (详见 systemd.kill(5) 中的 KillMode= 选项)。 可以使用 “ms”, “s”, “min”, “h” 这样的时间单位后缀。 若省略后缀则表示单位是秒。 设为零表示永不超时。 默认值为 DefaultTimeoutStartSec= 选项的值 (参见 systemd-system.conf(5) 手册)。

参见 systemd.exec(5) 与 systemd.kill(5) 以了解更多设置。

按照我找的实际观察下,开机后

    $ grep mnt /etc/fstab
    //10.0.23.85/test  /mnt  cifs  x-systemd.automount,noauto,username=zhangguanzhang@xxxx.com,password=xxxxx,iocharset=utf8,x-systemd.device-timeout=20s   0       0
    $ mount | grep mnt
    systemd-1 on /mnt type autofs (rw,relatime,fd=26,pgrp=1,timeout=0,minproto=5,maxproto=5,direct,pipe_ino=21338)
    $ ll /mnt
    total 51288
    -rwxr-xr-x 1 root root 52517368 Jul 23  2018 XshellPlus-6.0.0006r.exe
    $ mount | grep mnt
    systemd-1 on /mnt type autofs (rw,relatime,fd=26,pgrp=1,timeout=0,minproto=5,maxproto=5,direct,pipe_ino=21338)
    //10.0.23.85/test on /mnt type cifs (rw,relatime,vers=default,cache=strict,username=zhangguanzhang@xxxx.com,domain=,uid=0,noforceuid,gid=0,noforcegid,addr=10.0.23.85,file_mode=0755,dir_mode=0755,soft,nounix,serverino,mapposix,rsize=1048576,wsize=1048576,echo_interval=60,actimeo=1)

从上面可以看出只有我们访问的时候 systemd 收到了访问的请求才会去挂载的，也就是下面选项生效了

    noauto,x-systemd.automount

我们看下生成的挂载单元

    systemctl cat mnt.mount
    # /run/systemd/generator/mnt.mount
    # Automatically generated by systemd-fstab-generator
    [Unit]
    SourcePath=/etc/fstab
    Documentation=man:fstab(5) man:systemd-fstab-generator(8)
    [Mount]
    What=//10.0.23.85/test
    Where=/mnt
    Type=cifs
    Options=x-systemd.automount,noauto,username=zhangguanzhang@xxxx.com,password=xxxxxx,iocharset=utf8

我们还可以通过`systemctl status mnt.mount`查看啥时候挂载的，可以查看`systemctl status mnt.automount`查看是被哪个进程 pid 触发的挂载

    [root@CentOS76 ~]# systemctl status  mnt.mount
    ● mnt.mount - /mnt
       Loaded: loaded (/etc/fstab; bad; vendor preset: disabled)
       Active: active (mounted) since Wed 2019-04-24 18:46:19 CST; 3min 7s ago
        Where: /mnt
         What: //10.0.23.85/test
         Docs: man:fstab(5)
               man:systemd-fstab-generator(8)
      Process: 6901 ExecMount=/bin/mount //10.0.23.85/test /mnt -t cifs -o x-systemd.automount,username=zhangguanzhang@outlook.com,password=hj945417,_netdev,iocharset=utf8 (code=exited, status=0/SUCCESS)
        Tasks: 0
       Memory: 580.0K
    Apr 24 18:46:19 CentOS76 systemd[1]: Mounting /mnt...
    Apr 24 18:46:19 CentOS76 systemd[1]: Mounted /mnt.
    [root@CentOS76 ~]# date
    Wed Apr 24 18:49:33 CST 2019
    [root@CentOS76 ~]# systemctl status  mnt.automount
    ● mnt.automount
       Loaded: loaded (/etc/fstab; bad; vendor preset: disabled)
       Active: active (running) since Wed 2019-04-24 18:44:43 CST; 5min ago
        Where: /mnt
         Docs: man:fstab(5)
               man:systemd-fstab-generator(8)
    Apr 24 18:46:19 CentOS76 systemd[1]: Got automount request for /mnt, triggered by 6900 (ls)

参考
<http://www.jinbuguo.com/systemd/systemd.mount.html>
<http://www.jinbuguo.com/systemd/systemd.automount.html#>
<https://blog.csdn.net/richerg85/article/details/17917129>
[https://wiki.archlinux.org/index.php/Fstab\_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87)](<https://wiki.archlinux.org/index.php/Fstab_(%25E7%25AE%2580%25E4%25BD%2593%25E4%25B8%25AD%25E6%2596%2587)>)
