---
title: Mount(挂载)
---

# 概述

> 参考：
> - [Manual(手册),fstab(5)](https://man7.org/linux/man-pages/man5/fstab.5.html)
> - [Manual(手册),mount(8)](https://man7.org/linux/man-pages/man8/mount.8.html)

注意：

mount 命令无法列出 bind 类型的挂载(比如 Docker 和 Containerd 的 bind 类型挂载，不知道如何列出)。

不过 findmnt 命令可以列出使用 `mount --bind XX XX` 挂载的目录，效果如下

```bash
~]# findmnt
TARGET                                SOURCE                                  FSTYPE      OPTIONS
/                                     /dev/vda3                               ext4        rw,relatime
......
└─/mnt/cdrom                          /dev/vda3[/root/downloads/webvirtcloud] ext4        rw,relatime
```

# 挂载关联文件与配置

/etc/fstab # 包含各种 file systems 的描述性信息。系统启动时，根据该文件配置挂载分区到指定路径。

/etc/mtab # 是一个软连接，连接到 /proc/self/mounts

XXX.mount # 以 .mount 为后缀的 unit 文件，是由 systemd 管理的文件系统描述信息。systemd 将根据这些 unit 文件，将指定的文件系统挂载到系统中。

## fstab 文件详解

**File System Table(文件系统表，简称 fstab)**是一个关于文件系统的静态信息文件。默认路径为 **/etc/fstab**

fstab 包含操作系统中可以挂载的文件系统的描述性信息。每个描述信息占用单独一行，每行的各个字段由制表符或空格分隔。fsck、mount、umount 命令在执行某些操作时将会顺序读取该文件的每一行。

下面是一个 fstab 文件中关于文件系统描述的典型示例(一共具有 6 各字段)：

```bash
LABEL=t-home2 /home ext4 defaults,auto_da_alloc 0 2
```

下面各字段的名称中 fs 就是 file system 的简称。

### 第一字段(fs_spec) # 要挂载的分区或存储设备

该字段用于指定该文件系统所用存储设备是什么，可以是块设备、远程文件系统、文件系统镜像、交换分区 等等。

对于普通的挂载，它将为要挂载的设备保存（链接到）块专用设备节点（由 mknod（2）创建），例如“ /dev/cdrom”或“ /dev/sdb7”。 对于 NFS 挂载，此字段为：，例如\`knuth.aeb.nl:/'。 对于没有存储空间的文件系统，可以使用任何字符串，例如，这些字符串将显示在 df（1）输出中。 procfs 的典型用法是“ proc”； tmpfs 的“ mem”，“ none”或“ tmpfs”。 其他特殊文件系统（例如 udev 和 sysfs）通常不在 fstab 中列出。

可以给出 LABEL = <标签>或 UUID = 代替设备名称。这是推荐的方法，因为设备名称通常是硬件检测顺序的重合，并且在添加或删除其他磁盘时可能会更改。例如，“ LABEL = Boot”或“ UUID = 3e6be9de-8139-11d1-9106-a43f08d823a6”。 （使用特定于文件系统的工具，例如 e2label（8），xfs_admin（8）或 fatlabel（8）在文件系统上设置 LABEL）。

### 第二字段(fs_file) # 挂载的位置

该字段描述文件系统的挂载点（目标）。 使用 绝对路径 来描述一个挂载点.

### 第三字段(fs_vfstype) # 文件系统的类型

要挂载设备或是分区的文件系统类型，支持许多种不同的文件系统：ext2, ext3, ext4, reiserfs, xfs, jfs, smbfs, iso9660, vfat, ntfs, swap 及 auto。 设置成 auto 类型，mount 命令会猜测使用的文件系统类型，对 CDROM 和 DVD 等移动设备是非常有用的。

### 第四字段(fs_mntopts) # 与文件系统关联的挂载选项

详见下文 mount 命令中关于挂载选项的详解 [https://www.yuque.com/desistdaydream/learning/hla4gc#PfMRm](#PfMRm)

### 第五字段(fs_freq) # 确定需要转储哪些文件系统

dump 工具通过该字段的值决定何时作备份。 允许的数字是 0 和 1 。

- 0 # (默认值)表示忽略
- 1 # 表示进行备份。

大部分的用户是没有安装 dump 的 ，对这些用户而言该字段值应设为 0。

### 第六字段：fs_passno # 确定引导时执行文件系统检查的顺序。

fsck 工具读取该字段的值来决定需要检查的文件系统的检查顺序。允许的数字是 0, 1, 和 2。

- 0 # 表示设备不会被 fsck 所检查
- 1 # 根目录应当获得最高的优先权 1,
- 2 # 其它所有需要被检查的设备设置为 2.

# mount/umount # 文件系统挂载/卸载工具

> 参考：
> - [Manual(手册)，mount(8)](https://man7.org/linux/man-pages/man8/mount.8.html)

## Syntax(语法)

**mount \[-l] \[-t fstype]**
**mount -a \[-fFnrsvw] \[-t fstype] \[-O optlist]**
**mount \[-fnrsvw] \[-o options] \<DEVICE | MountPoint>** # 从 /etc/fstab 文件读取 Device 或 MountPoint 的信息后执行对应的挂载操作
**mount \[-fnrsvw] \[-t fstype] \[-o options] DEVICE MountPoint** # 将 DEVICE 设备挂载到 MountPoint 上

- MountPoint 通常是一个绝对路径，/PATH/DIR，即将 DEVICE 设备挂载在 PATH 路径下的 DIR 目录上

**mount --bind|--rbind|--move \<OldDIR> \<NewDIR>** # 执行 bind 挂载，将 OldDIR 目录挂载到 NewDIR 目录上
**mount --make-\[shared|slave|private|unbindable|rshared|rslave|rprivate|runbindable] MountPoint**

OPTIONS

- **-o, --options \<OPTS>** # 使用指定选项挂载，OPTS 是一个逗号分割的列表，可以从下面的 [Mount OPTIONS](#Mount%20OPTIONS(挂载选项)) 中查看所有可用的选项。
- **-r, --read-only** # 以只读方式挂载。可以被 `-o ro` 选项替代
- **-t \<STRING>** # 指明文件系统的类型。当前内核支持的文件系统类型可以从 /proc/filesystems 或 /lib/modules/$(uname -r)/kernel/fs 文件中获取
- **--bind|--rbind|--move** # 进行 bind 模式挂载

### Mount OPTIONS(挂载选项)

挂载选项分为两类：

- [FileSystem-Independent Mount Options(适用于所有文件系统的选项)](https://man7.org/linux/man-pages/man8/mount.8.html#FILESYSTEM-INDEPENDENT_MOUNT_OPTIONS)
- [FileSystem-Specific Mount Options(只适用于特定文件系统的选项)](https://man7.org/linux/man-pages/man8/mount.8.html#FILESYSTEM-SPECIFIC_MOUNT_OPTIONS)

- async # I/O 异步进行。
- auto # 在启动时或键入了 mount -a 命令时自动挂载。
- defaults # 使用文件系统的默认挂载参数，例如 ext4 的默认参数为:rw, suid, dev, exec, auto, nouser, async.
- dev - 解析文件系统上的块特殊设备。
- exec # 允许执行此分区的二进制文件。
- flush - vfat 的选项，更频繁的刷新数据，复制对话框或进度条在全部数据都写入后才消失。
- **noatime** # 不更新文件系统上 inode 访问记录。**可以提升性能**。
- noauto # 只在你的命令下被挂载。
- nodev - 不解析文件系统上的块特殊设备。
- nodiratime - 不更新文件系统上的目录 inode 访问记录，可以提升性能(参见 atime 参数)。
- noexec # 不允许执行此文件系统上的二进制文件。
- nosuid - 禁止 suid 操作和设定 sgid 位。
- nouser - 只能被 root 挂载。
- owner - 允许设备所有者挂载.
- **relatime** # 实时更新 inode access 记录。只有在记录中的访问时间早于当前访问才会被更新。（与 - noatime 相似，但不会打断如 mutt 或其它程序探测文件在上次访问后是否被修改的进程。），可以提升性能(参见 atime 参数)。
- **ro** # 以只读模式挂载文件系统
- **rw** # 以读写模式挂载文件系统
- suid - 允许 suid 操作和设定 sgid 位。这一参数通常用于一些特殊任务，使一般用户运行程序时临时提升权限。
- **sync** # I/O 同步进行。
- user - 允许任意用户挂载此文件系统，若无显示定义，隐含启用 noexec, nosuid, nodev 参数。
- users - 允许所有 users 组中的用户挂载文件系统.

## EXAMPLE

- mount -a # 会将 /etc/fstab 中定义的所有挂载点都挂上(一般是在系统启动时的脚本中调用，自己最好别用！)。
- mount /dev/sdb1 /mnt # /把/dev/sdb1 分区挂载到/mnt 目录
- mount --bind /root/tmp_one /root/tmp_two # 将 /root/tmp_one 目录挂载到 /root/tmp_two 目录下。

mount \[-t ] # 查看当前系统下的挂载信息\[查看指定的类型]
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hla4gc/1616167790128-b08b713d-147f-4c61-8289-47133e8124cf.jpeg)
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hla4gc/1616167790134-5f30f99c-3a70-44b7-b9c4-6c1cc3f9429f.jpeg)
各段落含义与 fstab 文件相同

# systemd 管理 mount

详见 [mount Unit](/docs/IT学习笔记/1.操作系统/3.Systemd%20 系统守护程序/Unit%20File/mount%20Unit.md 系统守护程序/Unit File/mount Unit.md)

https://wiki.archlinux.org/index.php/Fstab_(%25E7%25AE%2580%25E4%25BD%2593%25E4%25B8%25AD%25E6%2596%2587
