---
title: 从PTTYPE="dos"到TYPE="LVM2_member"的救援 · zhangguanzhang's Blog
---

同事叫我救援一台云主机，虽说是虚拟机，但是类比到硬件服务器还是一样的操作，这里记录下给后来者查阅

## 故障信息

控制台进去看到 centos7 的背景虚化的数字 7 + 转圈，重启下看下完整的错误，重启选了内核然后进到图形界面的时候按下 ecs 取消，观察终端

    [  OK ] Started Show Plymouth Boot Screen.
    [  OK ] Reached target Paths.
    [  OK ] Reached target Basic System.
    [  124.522110] dracut-initqueue[240]: Warning: dracut-initqueue timeout - starting timeout scripts
    [  125.034736] dracut-initqueue[240]: Warning: dracut-initqueue timeout - starting timeout scripts
    [  125.542788] dracut-initqueue[240]: Warning: dracut-initqueue timeout - starting timeout scripts
    [  126.522110] dracut-initqueue[240]: Warning: dracut-initqueue timeout - starting timeout scripts
    [  127.068643] dracut-initqueue[240]: Warning: dracut-initqueue timeout - starting timeout scripts
    [  127.576830] dracut-initqueue[240]: Warning: dracut-initqueue timeout - starting timeout scripts
    ...
    [  185.082387] dracut-initqueue[240]: Warning: Could not boot.
    [  185.118736] dracut-initqueue[240]: Warning: /dev/centos/root does not exist
    [  185.119239] dracut-initqueue[240]: Warning: /dev/mapper/centos-root does not exist
              Starting Dracut Emergency Shell...
    Warning: /dev/centos/root does not exist
    Warning: /dev/mapper/centos-root does not exist
    Generating "/run/initramfs/rdsosreport.txt"

    Entering emergency mode. Exit the shell to continue.
    Type "journalctl" to view system logs.
    You might want to save "/run/initramfs/rdsosreport.txt" to a USB stack or /boot
    after mounting them and attach it to a bug report.

    dracut:/#

## 处理

### 挂载 iso 进救援模式

找不到根分区，关闭虚机，后台拷贝下系统盘的卷先备份下。然后给虚机的 IDE 光驱挂载了个 centos 的 iso，修改虚机启动顺序到 ISO，进`Troubleshooting` –> `Rescue a CentOS Linux system`

    1) Continue
    2) Read-only mount
    3) Skip to shell
    3) Quit (Reboot)
    Please make a selection from the above: 1

选择了 1 后提示没有任何 Linux 分区

    =====================================================================================
    =====================================================================================
    Rescue Mount

    You don't have any Linux partitions. The system will reboot automatically when you exit from the shell.
    Please press <return> to get a shell.

按下回车进入交互式 shell

    sh-4.2# lsblk
    NAME            MAJ:MIN RM   SIZE RO TYPE MOUNTPOINT
    vda             252:0    0    40G  0 disk
    ├─vda1          252:1    0     2M  0 part
    ├─vda2          252:2    0   200M  0 part
    └─vda3          252:3    0  39.8G  0 part
    vdb             252:16   0   400G  0 disk
    loop0             7:1    0 420.3M  1 loop
    loop1             7:1    0     2G  1 loop
    ├─live-rw       253:0    0     2G  0 dm    /
    └─live-base     253:1    0     2G  1 dm
    loop2             7:2    0   512M  1 loop
    └─live-rw       253:0    0     2G  0 dm    /
    sh-4.2# fdisk -l /dev/vda

    Disk /dev/vda: 42.9 GB, 42949672960 bytes, 83886080 sectors
    Units = sectors of 1 * 512 = 512 bytes
    Sector size (logical/physical): 512 bytes / 512 bytes
    I/O size (minimum/optimal): 512 bytes / 512 bytes
    Disk label type: dos
    Disk identifier: 0x000ad4f2

       Device Boot      Start         End      Blocks   Id  System
    /dev/vda1            2048        6143        2048   83  Linux
    /dev/vda2   *        6144      415743      204800   83  Linux
    /dev/vda3          415744    83886079    41735168   8e  Linux LVM
    sh-4.2# blkid
    /run/install/repo/LiveOS/squashfs.img: TYPE="squashfs"
    /dev/sr0: UUID="2018-05-03-20-55-23-00" LABEL="Centos 7 x86_64" TYPE="iso9660" PTTYPE="dos"
    /dev/sr1: UUID="2019-11-01-16-33-37-00" LABEL="config-2" TYPE="iso9660"
    /dev/vda1: UUID="e438c18a-c97d-432c-ae66-a538cd1fbb4d" TYPE="xfs"
    /dev/vda3: PTTYPE="dos"
    /dev/loop0: TYPE="squashfs"
    ...

查看下块，vda2 是 boot，vda3 是 lvm 也就是根所在，问题是`/dev/vda3: PTTYPE="dos"`不知道为何变成了 dos 类型，正常应该是`TYPE="LVM2_member"`

    /dev/vda3: UUID="xxxxxxxxxxxxxx" TYPE="LVM2_member"

看看 lvm 的状态

    mkdir /mnt/sysimage
    vgchange -a y
    mount /dev/centos/root /mnt/sysimage
    ls -l /mnt/sysimage

发现根分区的文件还在

### 重做 lvm 为了去掉 PTTYPE=”dos”

尝试着重做 pv 试试

    sh-4.2# vgremove centos
    Do you really wanto ro remove volume group "centos" containing 1 logical volumes? [y/n]: y
      Logical volume "root" successfully removed
      Volume group "centos" successfully removed
    sh-4.2# pvremove /dev/vda3
      Labels on physical volume "/dev/vda3" successfully wiped.
    sh-4.2# pvcreate /dev/vda3
    WARNING: dos signature detected on /dev/vda3 at offset 510. Wipe it? [y/n]: y
      Wiping dos signature on /dev/vda3.
      Physical volume "/dev/vda3" successfully created.
    sh-4.2# vgcreate centos /dev/vda3
      Volume group "centos" successfully created
    sh-4.2# lvcreate -n root -l 100%FREE centos
      Logical volume "root" created
    sh-4.2# mkdir /mnt/root
    sh-4.2# mount /dev/centos/root /mnt/root

被后面的`xfs_repair`输出滚动冲没了，大致就是 lvcreate 的时候提示有 xfs 标签，选择不擦除，最终得到了个残缺的的`/dev/centos/root`，然后`xfs_repair`它后重启也无法正常开机。再次进救援模式挂载了后，在 chroot 到故障的根进不去报错`/bin/bash no such file`，才意识到损坏了文件，很多 so 文件都丢了

最开始留有备份，打算在文件层面恢复

### 最后也应该是最开应该做的正确操作

下发了台不是 lvm 的 centos7.6 机器，然后给该机器挂载了数据盘为 vdb(50G，其实大于等于故障机器根的真实占用大小即可)，备份的卷挂载为 vdc。利用数据盘中转下原有的根的文件

打算把 lvm 的文件系统文件拷贝到数据盘 vdb 的文件系统上，然后在故障机器的救援模式下挂载这个数据盘，把数据盘的根文件拷贝到残缺的系统盘的根下

#### 格式化 vdb

    [root@fix-data ~]# lsblk
    NAME   MAJ:MIN RM  SIZE RO TYPE MOUNTPOINT
    sr0     11:0    1 1024M  0 rom
    sr1     11:1    1  464K  0 rom
    vda    253:0    0   40G  0 disk
    └─vda1 253:1    0   40G  0 part /
    vdb    253:16   0   50G  0 disk
    [root@fix-data ~]# parted /dev/vdb
    GNU Parted 3.1
    Using /dev/vdb
    Welcome to GNU Parted! Type 'help' to view a list of commands.
    (parted) mklabel gpt
    (parted) mkpart
    Partition name?  []? 1
    File system type?  [ext2]? xfs
    Start? 1
    End? 100%
    (parted) p
    Model: Virtio Block Device (virtblk)
    Disk /dev/vdb: 53.7GB
    Sector size (logical/physical): 512B/512B
    Partition Table: gpt
    Disk Flags:

    Number  Start   End     Size    File system  Name  Flags
     1      1049kB  53.7GB  53.7GB               1

    (parted) q
    Information: You may need to update /etc/fstab.

    [root@fix-data ~]# mkfs.xfs /dev/vdb1
    meta-data=/dev/vdb1              isize=512    agcount=4, agsize=3276672 blks
             =                       sectsz=512   attr=2, projid32bit=1
             =                       crc=1        finobt=0, sparse=0
    data     =                       bsize=4096   blocks=13106688, imaxpct=25
             =                       sunit=0      swidth=0 blks
    naming   =version 2              bsize=4096   ascii-ci=0 ftype=1
    log      =internal log           bsize=4096   blocks=6399, version=2
             =                       sectsz=512   sunit=0 blks, lazy-count=1
    realtime =none                   extsz=4096   blocks=0, rtextents=0
    [root@fix-data ~]# mkdir -p /mnt/{root,data}

#### 挂载需要修复的系统盘的克隆卷

后台挂载好后

    # lsblk
    NAME   MAJ:MIN RM  SIZE RO TYPE MOUNTPOINT
    sr0     11:0    1 1024M  0 rom
    sr1     11:1    1  464K  0 rom
    vda    253:0    0   40G  0 disk
    └─vda1 253:1    0   40G  0 part /
    vdb    253:16   0   50G  0 disk
    └─vdb1 253:17   0   50G  0 part
    vdc    253:32   0   40G  0 disk
    ├─vdc1 253:33   0    2M  0 part
    ├─vdc2 253:34   0  200M  0 part
    └─vdc3 253:35   0 39.8G  0 part

安装 lvm2 工具 (新机器因为不是 lvm 的根分区所以系统没有安装这个)

    yum install -y lvm2

激活 lvm 状态

    $ vgchange -a y
      1 logical volume(s) in volume group "centos" now active

根 –> /mnt/root/

/boot –> /mnt/root/boot

数据盘 –> /mnt/data/

    mount /dev/centos/root /mnt/root/
    mount /dev/vdc2 /mnt/root/boot
    mount /dev/vdb1 /mnt/data

然后拷贝

    cp -a /mnt/root/* /mnt/data/

拷贝完后取消挂载并关机

    umount /mnt/data/
    umount /mnt/root/boot/
    umount /mnt/root/
    poweroff

#### 拷贝

后台把该数据盘挂载到故障机器上，故障机器在救援模式里操作

    sh-4.2# lsblk
    NAME            MAJ:MIN RM   SIZE RO TYPE MOUNTPOINT
    sr0              11:0    1   4.2G  0 rom  /run/install/repo
    sr1              11:1    1   464K  0 rom
    vda             252:0    0    40G  0 disk
    ├─vda1          252:1    0     2M  0 part
    ├─vda2          252:2    0   200M  0 part /mnt/sysimage/boot
    └─vda3          252:3    0  39.8G  0 part
      └─centos-root 253:2    0  39.8G  0 lvm  /mnt/sysimage
    vdb             252:16   0   400G  0 disk
    vdc             252:32   0    50G  0 disk
    └─vdc1          252:33   0    50G  0 part
    loop0             7:1    0 420.3M  1 loop
    loop1             7:1    0     2G  1 loop
    ├─live-rw       253:0    0     2G  0 dm    /
    └─live-base     253:1    0     2G  1 dm
    loop2             7:2    0   512M  1 loop
    └─live-rw       253:0    0     2G  0 dm    /

可以看到数据盘为 vdc，挂载 vdc1 到`/mnt/data`，然后拷贝到`/mnt/sysimage`

    sh-4.2# mkdir /mnt/data
    sh-4.2# mount /dev/vdc1 /mnt/data
    sh-4.2# ls -l /mnt/data
    total 16
    lrwxrwxrwx    1 root root    7 Jul 24  2018 bin -> usr/bin
    dr-xr-xr-x.   5 root root 4096 Jul 24  2018 boot
    drwxr-xr-x.   2 root root   18 Dec  3 04:24 dev
    drwxr-xr-x. 143 root root 8192 Dec  3 03:54 etc
    drwxr-xr-x.   3 root root   20 Jul 25  2018 home
    lrwxrwxrwx    1 root root    7 Jul 24  2018 lib -> usr/lib
    lrwxrwxrwx    1 root root    9 Jul 24  2018 lib64 -> usr/lib64
    drwxr-xr-x.   2 root root    6 Apr 11  2018 media
    drwxr-xr-x.   2 root root    6 Apr 11  2018 mnt
    drwxr-xr-x.   4 root root   34 Nov  1 17:19 opt
    drwxr-xr-x.   2 root root    6 Jul 24  2018 proc
    dr-xr-x---.   9 root root  258 Nov  5 16:14 root
    drwxr-xr-x.   2 root root    6 Jul 24  2018 run
    lrwxrwxrwx    1 root root    8 Jul 24  2018 sbin -> usr/sbin
    drwxr-xr-x.   2 root root    6 Apr 11  2018 srv
    drwxr-xr-x.   2 root root    6 Jul 24  2018 sys
    drwxrwxrwt.   7 root root  114 Dec  3 04:55 tmp
    drwx------    7 root root   66 Jul 24  2018 usr
    drwxr-xr-x.  21 root root 4096 Jul 24  2018 var
    drwxr-xr-x.   2 root root    6 Nov  4 02:09 version
    sh-4.2# cd /mnt/sysimage
    sh-4.2# cp -a /mnt/data/* .

重做 grub.cfg，该系统不是 grub2，如果是 grub2 则 / boot 下有 grub2 目录

    mount -o bind /dev /mnt/sysimage/dev
    mount -o bind /proc /mnt/sysimage/proc
    mount -o bind /run /mnt/sysimage/run
    mount -o bind /sys /mnt/sysimage/sys
    mv boot/grub/grub.cfg{,.bak}

然后 chroot 进来生成`grub.cfg`

    chroot .
    grub-mkconfig -o /boot/grub2/grub.cfg

开机正常

<https://zhangguanzhang.github.io/2019/12/03/dos-to-gpt/>
