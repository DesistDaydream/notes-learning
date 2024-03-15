---
title: opensuse的一次救援
---

昨晚吃完晚饭回到办公室，右边同事在控制台看着一个 suse 起不来一直启动的时候卡在 suse 的蜥蜴 logo 背景图那。见我来了叫我看下，他们已经尝试过恢复快照，但是还不行，应该是很久之前损坏的，只不过因为没重启没发现，我叫他重启下看看卡在哪，重启后进入内核后显示背景图那按下 esc 然后看卡在 / sysroot 挂载那。目测分区损坏了，经历了 ubuntu 的安装 iso 的 rescue mode 就是渣渣后，我还是信任 centos 的 iso。

## 处理

### 先备份和准备工作

关闭虚机，后台拷贝下系统盘的卷先备份下。然后给虚机的 IDE 光驱挂载了个 centos 7.5 DVD 的 iso，修改虚机启动顺序到 ISO，进`Troubleshooting` –> `Rescue a CentOS Linux system`
一般损坏的都不建议选 1，因为挂载不上，所以是选 3 手动处理

### Device or resource busy

    1) Continue
    2) Read-only mount
    3) Skip to shell
    4) Quit (Reboot)
    Please make a selection from the above: 3

最开始我 lsblk 和看了下硬盘的分区表，最后`vgchange -a y`激活 lvm 后`xfs_repair /dev/mapper/suse-lv_root`的时候提示该设备繁忙，遂查看了下

    sh-4.2# lsof /dev/mapper/suse-lv_root
    sh-4.2# ps aux | less

lsof 和 fuser 都是返回空的，最后就`ps aux`一个个看，发现了个 mount 进程一直 hung 在那

    sh-4.2# ps aux | grep moun[t]
    root       6126  0.0  0.0  19940    840 pts/0      D+   11:56   0:00 /usr/bin/mount -t xfs -o defaults,ro /dev/mapper/suse-lv_root /mnt/sysimage

这个进程尝试过了，死活杀不掉，进 rescue mode 的时候选的`Skip to shell`，以为是 iso 的版本 bug，换了一个 7.6 minimal 的 iso 进入 rescue mode 后不选择直接`ctrl+alt+F2`进到 tty 还是一样会自动挂载，于是想下从父进程的角度上看看能不能处理

    [anaconda root@localhost /]# ps aux | grep moun[t]
    root       6113  0.0  0.0  19940    840 pts/0      D+   12:02   0:00 /usr/bin/mount -t xfs -o defaults,ro /dev/mapper/suse-lv_root /mnt/sysimage
    [anaconda root@localhost /]# ps -Al | grep mount
    4 D     0   6113   5862  0  80   0 -  4985 xfs_bu pts/0    00:00:00 mount
    [anaconda root@localhost /]# ps aux | grep 586[2]
    root       5862  0.0  0.0  19940    840 pts/0      D+   12:02   0:00 python anaconda

找到了该 mount 的父进程是 anaconda，也就是我们进入 rescue mode 的第一个 tty 提供交互，直接杀掉它，mount 终止

    [anaconda root@localhost /]# kill -9 5862;kill 6113
    bash: kill: (6113) - No such process

但是还是 busy，发现该 mount 又 tm 的起来了，最终想了个骚套路，进入 rescue mode，然后 4 个选项不选择，直接`ctrl+alt+F2`进到 tty 后杀掉 mount 的进程后把 mount 命名改名，执行下面

    [anaconda root@localhost /]# mv /usr/bin/mount{,.bak}  #先改名，再杀进程
    [anaconda root@localhost /]# ps aux | grep moun[t]
    root       6128  0.3  0.0  19940    844 pts/0      D+   12:06   0:00 /usr/bin/mount -t xfs -o defaults,ro /dev/mapper/suse-lv_root /mnt/sysimage
    [anaconda root@localhost /]# ps -Al | grep mount
    4 D     0   6128   5877  0  80   0 -  4985 xfs_bu pts/0    00:00:00 mount
    [anaconda root@localhost /]# kill -9 5877;kill 6128
    bash: kill: (6128) - No such process
    [anaconda root@localhost /]# ps aux | grep moun[t] #然后再也没mount进程

### 修复

接着前面的，激活 lvm，这里不详细说，可以自己去`lsblk`和`fdisk -l /dev/vdx`去看相关分区信息

    [anaconda root@localhost /]# vgchange -a y
      1 logical volume(s) in volume group "suse" now active
      4 logical volume(s) in volume group "vgsap" now active


    [anaconda root@localhost /]# xfs_repair /dev/mapper/suse-lv_root
    ERROR: The filesystem has valuable metadata changes in a log which needs to
    be replayed.  Mount the filesystem to replay the log, and unmount it before
    re-running xfs_repair.  If you are unable to mount the filesystem, then use
    the -L option to destroy the log and attempt a repair.
    Note that destroying the log may cause corruption -- please attempt a mount
    of the filesystem before doing this.

该报错大致意思是: 文件系统的 log 需要在 repair 之前先 mount 它来触发回放 log，如果无法挂载，使用`xfs_repair`带上`-L`选项摧毁 log 强制修复
正确姿势是先使用`xfs_metadump`导出 metadata，见文章 <https://serverfault.com/questions/777299/proper-way-to-deal-with-corrupt-xfs-filesystems>
这里因为已经损坏了，没必要测试 mount 了，并且我未导出 metadata，直接 - L 修复的，下次遇到了相似场景可以试下`xfs_metadump`

    [anaconda root@localhost /]# xfs_repair -L /dev/mapper/suse-lv_root

漫长的等待修复，然后卡在了一个 inode 那，等待了 20 分钟直接`ctrl c`取消再来，然当这次不需要带`-L`选项

    [anaconda root@localhost /]# xfs_repair /dev/mapper/suse-lv_root
    ...
    ...
    resetting inode 15847758 nlinks from 0 to 2
    resetting inode 16180728 nlinks from 0 to 2
    resetting inode 16500950 nlinks from 0 to 2
    resetting inode 17347042 nlinks from 0 to 2
    resetting inode 19414733 nlinks from 0 to 2
    Metadata corruption detected at xfs_dir3_block block 0x2a09ba8/0x1000
    libxfs_writebufr: write verufer failed on xfs_dur3_block bno 0x2a09ba8/0x1000
    Metadata corruption detected at xfs_dir3_block block 0x145ce28/0x1000
    libxfs_writebufr: write verufer failed on xfs_dur3_block bno 0x145ce28/0x1000
    ...
    ...
    Maximum metadata LSN (1919513701:1600352110) is ahead of log (1:2).
    Format log to cycle 1919513704.
    releasing dirty buffer (bulk) to free list!releasing dirty buffer (bulk) to free list!releasing dirty buffer (bulk) to free list!releasing dirty buffer (bulk) to free list!releasing dirty buffer (bulk) to free list!releasing dirty buffer (bulk) to free list!releasing dirty buffer (bulk) to free list!releasing dirty buffer (bulk) to free list!releasing dirty buffer (bulk) to free list!releasing dirty buffer (bulk) to free list!done

然后再次修复

    [anaconda root@localhost /]# xfs_repair /dev/mapper/suse-lv_root
    Phase 1 - find and verify superblock...
    Phase 2 - using internal log
            - zero log...
            - scan filesystem freespace and inode maps...
            - found root inode chunk
    Phase 3 - for each AG...
            - scan and clear agi unlinked lists...
            - process known inodes and perfrom inode discovery...
            - agno = 0
    ...
            - setting up duplicate extent list...
            - check for inodes claiming duplicate blocks...
            - agno = 1
            - agno = 0
            - agno = 2
            - agno = 3
    Phase 5 - rebuild AG headers and trees...
            - reset superblock...
    Phase 6 - check inode connectivity...
            - resetting contents of realtime bitmap and summary inodes
            - traversing filesystem ...
            - traversal finished ...
            - moveing disconnected inodes to lost_found ...
    Phase 7 - verify and correct link counts...
    resetting inode 64 nlinks from 25 to 24
    done

然后挂载试试

    [anaconda root@localhost /]# mount /dev/mapper/suse-lv_root /mnt/sysimage
    bash： mount: command not found
    [anaconda root@localhost /]# mount.bak /dev/mapper/suse-lv_root /mnt/sysimage #漫长等待，大概30多秒
    [anaconda root@localhost /]# ls -l /mnt/sysimage
    total 40296
    drwxr-xr-x    2 root root       4096 Aug  6  2018 bin
    drwxr-xr-x    3 root root          6 Aug  6  2018 boot
    drwxr-xr-x   22 root root          6 Aug  6 22:19 dev
    drwxr-xr-x  131 root root       8192 Nov 30 04:05 etc
    drwxr-xr-x    5 root root         46 Oct 18  2018 home
    drwxr-xr-x   12 root root       8192 Nov 30 04:04 lib
    drwxr-xr-x    7 root root       8192 Aug  6  2018 lib64
    drwxr-xr-x 2270 root root   27242496 Dec  4 20:47 lost+found
    drwxr-xr-x    2 root root          6 Jun 27  2017 mnt
    drwxr-xr-x    2 root root          6 Jun 27  2017 opt
    dr-xr-xr-x  190 root root          6 Aug  6  2018 proc
    drwx------   21 root root       4096 Nov 30 04:39 root
    drwxr-xr-x   31 root root          6 Aug  6  2018 run
    drwxr-xr-x    4 root sapsys        6 Oct 11  2018 sapmnt
    drwxr-xr-x    2 root root       8192 Oct 11  2018 sbin
    drwxr-xr-x    2 root root          6 Jun 27  2017 selinux
    drwxr-xr-x    9 root root       4096 Oct 11  2018 software
    drwxr-xr-x    4 root root         28 Aug  6  2018 srv
    dr-xr-xr-x   13 root root          0 Dec  4 22:16 sys
    drwxrwxrwt   31 root root       4096 Dec  3 22:18 tmp
    drwxr-xr-x   14 root root        182 Nov 30 04:37 usr
    drwxr-xr-x   13 root root        201 Nov 30 04:37 var

然后取消光驱挂载，修改启动顺序重启，能进到登陆，直接`ctrl+alt+F2`进到 tty 登陆，发现没有网络，查看了下失败的启动，控制台观察的，输出不能被复制，下面命令输出大致的写下

    $ systemctl --failed
      UNIT                  LOAD   ACTIVE SUB    DESCRIPTION
    ● cryptctl-auto-unlock@sys-devices-pci0000:00-0000:00:08.0-virtio3-block-vda.service
    ● cryptctl-auto-unlock@aD7Wov-Krfg-KPbq-Dnf6-1dAj-e9dM-N7dUir.service
    ● cryptctl-auto-unlock@abd69e01-d874-4658-b738-1107d33cd84c.service
    ● cryptctl-auto-unlock@abd69e01-d874-4658-b738-1107d33cd84c.service
    ● cryptctl-auto-unlock@0zSmA1-nPGR-FuVE-ZIvq-vxhl-2WdX-eh58e2.service
    ● postfix.service       loaded failed failed Postfix Mail Transport Agent
    ● wicked.service        loaded failed failed wicked managed network interfaces
    ● wickedd-auto4.service loaded failed failed wicked AutoIPv4 supplicant service
    ● wickedd-dhcp4.service loaded failed failed wicked DHCPv4 supplicant service
    ● wickedd-dhcp6.service loaded failed failed wicked DHCPv6 supplicant service
    ● wickedd.service       loaded failed failed wicked network management service daemon

    LOAD   = Reflects whether the unit definition was properly loaded.
    ACTIVE = The high-level unit activation state, i.e. generalization of SUB.
    SUB    = The low-level unit activation state, values depend on unit type.

查看了下系统日志

    vi /var/log/messages
    2019-12-04T21:42:56.401177+08:00 bpcprdascs1 cryptctl[1506]: open /etc/sysconfig/cryptctl-client: no such file or directory
    2019-12-04T21:42:56.403110+08:00 bpcprdascs1 cryptctl[1515]: open /etc/sysconfig/cryptctl-client: no such file or directory
    2019-12-04T21:42:56.408475+08:00 bpcprdascs1 cryptctl[1495]: open /etc/sysconfig/cryptctl-client: no such file or directory
    2019-12-04T21:42:56.408634+08:00 bpcprdascs1 cryptctl[1529]: open /etc/sysconfig/cryptctl-client: no such file or directory
    2019-12-04T21:42:56.408807+08:00 bpcprdascs1 cryptctl[1508]: open /etc/sysconfig/cryptctl-client: no such file or directory
    2019-12-04T21:42:56.414068+08:00 bpcprdascs1 systemd[1]: cryptctl-auto-unlock@sys-devices-pci0000:00-0000:00:08.0-virtio3-block-vda.service: Main process exited, code=exited, status=1/FAILURE
    2019-12-04T21:42:56.414237+08:00 bpcprdascs1 systemd[1]: cryptctl-auto-unlock@sys-devices-pci0000:00-0000:00:08.0-virtio3-block-vda.service: Unit entered failed state.
    2019-12-04T21:42:56.414319+08:00 bpcprdascs1 systemd[1]: cryptctl-auto-unlock@sys-devices-pci0000:00-0000:00:08.0-virtio3-block-vda.service: Failed with result 'exit-code'.
    2019-12-04T21:42:56.414403+08:00 bpcprdascs1 systemd[1]: cryptctl-auto-unlock@aD7Wov-Krfg-KPbq-Dnf6-1dAj-e9dM-N7dUir.service: Main process exited, code=exited, status=1/FAILURE
    2019-12-04T21:42:56.414467+08:00 bpcprdascs1 systemd[1]: cryptctl-auto-unlock@aD7Wov-Krfg-KPbq-Dnf6-1dAj-e9dM-N7dUir.service: Unit entered failed state.
    2019-12-04T21:42:56.414528+08:00 bpcprdascs1 systemd[1]: cryptctl-auto-unlock@aD7Wov-Krfg-KPbq-Dnf6-1dAj-e9dM-N7dUir.service: Failed with result 'exit-code'.
    2019-12-04T21:42:56.414596+08:00 bpcprdascs1 systemd[1]: cryptctl-auto-unlock@abd69e01-d874-4658-b738-1107d33cd84c.service: Main process exited, code=exited, status=1/FAILURE
    2019-12-04T21:42:56.414657+08:00 bpcprdascs1 systemd[1]: cryptctl-auto-unlock@abd69e01-d874-4658-b738-1107d33cd84c.service: Unit entered failed state.
    2019-12-04T21:42:56.414735+08:00 bpcprdascs1 systemd[1]: cryptctl-auto-unlock@abd69e01-d874-4658-b738-1107d33cd84c.service: Failed with result 'exit-code'.
    2019-12-04T21:42:56.414794+08:00 bpcprdascs1 systemd[1]: cryptctl-auto-unlock@0zSmA1-nPGR-FuVE-ZIvq-vxhl-2WdX-eh58e2.service: Main process exited, code=exited, status=1/FAILURE
    2019-12-04T21:42:56.414851+08:00 bpcprdascs1 systemd[1]: cryptctl-auto-unlock@0zSmA1-nPGR-FuVE-ZIvq-vxhl-2WdX-eh58e2.service: Unit entered failed state.
    2019-12-04T21:42:56.414907+08:00 bpcprdascs1 systemd[1]: cryptctl-auto-unlock@0zSmA1-nPGR-FuVE-ZIvq-vxhl-2WdX-eh58e2.service: Failed with result 'exit-code'.

发现该文件丢失，同样系统的机器上去把内容手动创建，然后重启只剩下这些

    ● wicked.service        loaded failed failed wicked managed network interfaces
    ● wickedd-auto4.service loaded failed failed wicked AutoIPv4 supplicant service
    ● wickedd-dhcp4.service loaded failed failed wicked DHCPv4 supplicant service
    ● wickedd-dhcp6.service loaded failed failed wicked DHCPv6 supplicant service
    ● wickedd.service       loaded failed failed wicked network management service daemon

找到相关日志，或者手动启动 wicked 或者网卡也报错下面类似

    2019-12-04T21:54:50.170654+08:00 bpcprdascs1 wickedd[1399]: Failed to register dbus bus name "org.opensuse.Network" (Connection ":1.2" is not allowed to own the service "org.opensuse.Network" due to security policies in the configuration file)
    2019-12-04T21:54:50.170657+08:00 bpcprdascs1 wickedd[1399]: unable to initialize dbus service
    2019-12-04T21:54:50.170659+08:00 bpcprdascs1 systemd[1]: wickedd.service: Main process exited, code=exited, status=1/FAILURE
    2019-12-04T21:54:50.170661+08:00 bpcprdascs1 systemd[1]: Failed to start wicked network management service daemon.
    ...
    2019-12-04T22:02:05.868058+08:00 bpcprdascs1 wicked: /org/opensuse/Network/Interface.getManagedObjects failed. Server responds:
    2019-12-04T22:02:05.868883+08:00 bpcprdascs1 wicked: org.freedesktop.DBus.Error.ServiceUnknown: The name org.opensuse.Network was not provided by any .service files

这个错误找了一圈都没正确的解决办法，还是自己突发奇想在`/etc/dbus-1/`对比了下发现文件丢失
正常机器上

    bpcprdascs2:/etc/dbus-1/system.d # find /etc/dbus-1/ -type f
    /etc/dbus-1/system.d/org.opensuse.Snapper.conf
    /etc/dbus-1/system.d/org.freedesktop.hostname1.conf
    /etc/dbus-1/system.d/org.freedesktop.import1.conf
    /etc/dbus-1/system.d/org.freedesktop.locale1.conf
    /etc/dbus-1/system.d/org.freedesktop.login1.conf
    /etc/dbus-1/system.d/org.freedesktop.machine1.conf
    /etc/dbus-1/system.d/org.freedesktop.systemd1.conf
    /etc/dbus-1/system.d/org.freedesktop.timedate1.conf
    /etc/dbus-1/system.d/com.redhat.PrinterDriversInstaller.conf
    /etc/dbus-1/system.d/org.freedesktop.UPower.conf
    /etc/dbus-1/system.d/org.freedesktop.GeoClue2.Agent.conf
    /etc/dbus-1/system.d/org.freedesktop.GeoClue2.conf
    /etc/dbus-1/system.d/bluetooth.conf
    /etc/dbus-1/system.d/com.redhat.tuned.conf
    /etc/dbus-1/system.d/org.freedesktop.PolicyKit1.conf
    /etc/dbus-1/system.d/org.freedesktop.UDisks2.conf
    /etc/dbus-1/system.d/org.freedesktop.RealtimeKit1.conf
    /etc/dbus-1/system.d/org.freedesktop.Accounts.conf
    /etc/dbus-1/system.d/org.opensuse.Network.AUTO4.conf
    /etc/dbus-1/system.d/org.opensuse.Network.DHCP4.conf
    /etc/dbus-1/system.d/org.opensuse.Network.DHCP6.conf
    /etc/dbus-1/system.d/org.opensuse.Network.Nanny.conf
    /etc/dbus-1/system.d/org.opensuse.Network.conf
    /etc/dbus-1/system.d/pulseaudio-system.conf
    /etc/dbus-1/system.d/org.freedesktop.PackageKit.conf
    /etc/dbus-1/system.d/cups.conf
    /etc/dbus-1/system.d/org.opensuse.CupsPkHelper.Mechanism.conf
    /etc/dbus-1/system.d/gdm.conf
    /etc/dbus-1/session.conf
    /etc/dbus-1/system.conf

该故障机器上

    bpcprdascs1:/var/log # find /etc/dbus-1/ -type f
    /etc/dbus-1/system.d/org.opensuse.Snapper.conf
    /etc/dbus-1/system.d/org.freedesktop.hostname1.conf
    /etc/dbus-1/system.d/org.freedesktop.import1.conf
    /etc/dbus-1/system.d/org.freedesktop.locale1.conf
    /etc/dbus-1/system.d/org.freedesktop.login1.conf
    /etc/dbus-1/system.d/org.freedesktop.machine1.conf
    /etc/dbus-1/system.d/org.freedesktop.systemd1.conf
    /etc/dbus-1/system.d/org.freedesktop.timedate1.conf
    /etc/dbus-1/system.d/com.redhat.PrinterDriversInstaller.conf
    /etc/dbus-1/system.d/org.freedesktop.UPower.conf
    /etc/dbus-1/system.d/org.freedesktop.GeoClue2.Agent.conf
    /etc/dbus-1/system.d/org.freedesktop.GeoClue2.conf
    /etc/dbus-1/system.d/bluetooth.conf
    /etc/dbus-1/system.d/com.redhat.tuned.conf
    /etc/dbus-1/system.d/org.freedesktop.PolicyKit1.conf
    /etc/dbus-1/system.d/org.freedesktop.RealtimeKit1.conf
    /etc/dbus-1/session.conf
    /etc/dbus-1/system.conf

因为故障机器的网络无法启动，即使手动`ip addr add`也报错 dbus，所以无法通过网络 scp。于是在后台正常机器给添加了一个数据盘，把该目录的文件拷贝到数据盘上，再把数据盘挂载到故障机器上。然后 cp 拷贝完重启，然后网络起来了
只剩下故障

    $ systemctl --failed
      UNIT                  LOAD   ACTIVE SUB    DESCRIPTION
    ● wickedd-auto4.service loaded failed failed wicked AutoIPv4 supplicant service
    ● wickedd-dhcp4.service loaded failed failed wicked DHCPv4 supplicant service
    ● wickedd-dhcp6.service loaded failed failed wicked DHCPv6 supplicant service

上面三个通过系统日志可以找到是文件丢失，其他机器上去拷贝就行了，当然也不是只有这三个服务的文件丢失，其他服务的文件也丢失了，自行看下系统日志处理下

    2019-12-04T21:54:50.170648+08:00 bpcprdascs1 display-manager[1429]: /usr/lib/X11/display-manager: line 17: /etc/sysconfig/displaymanager: No such file or directory
    2019-12-04T22:18:03.689878+08:00 bpcprdascs1 systemd[1395]: wickedd-dhcp6.service: Failed at step EXEC spawning /usr/lib/wicked/bin/wickedd-dhcp6: No such file or directory
    ...
    2019-12-04T22:18:03.689884+08:00 bpcprdascs1 systemd[1396]: wickedd-dhcp4.service: Failed at step EXEC spawning /usr/lib/wicked/bin/wickedd-dhcp4: No such file or directory
    ...
    2019-12-04T22:18:03.689897+08:00 bpcprdascs1 systemd[1403]: wickedd-auto4.service: Failed at step EXEC spawning /usr/lib/wicked/bin/wickedd-auto4: No such file or directory

<https://zhangguanzhang.github.io/2019/12/05/suse-fix-data-but-device-busy/>
