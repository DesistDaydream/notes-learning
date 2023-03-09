---
title: inode 已满解决方法
---

## 问题描述

在 Linux 系统的云服务器 ECS 实例内创建文件时，出现类似如下空间不足的提示。

    No space left on device …

## 问题原因

导致该问题的可能原因如下所示：

- 磁盘分区空间使用率达到百分之百。

- 磁盘分区 inode 使用率达到百分之百。

- 存在僵尸文件。

- 挂载点覆盖。

## 解决方案

> 阿里云提醒您：

要解决该问题，请根据不同的问题原因，通过以下方式进行处理：

- 分区容量满

- inode 容量满

- 修改 inode 数量

- 僵尸文件分析删除

- 挂载点覆盖

### 分区容量满的处理

- 登录服务器，使用`df -h`命令查看磁盘使用率，其中的**Mounted on**指挂载的目录。

- 进入根目录，执行 `du -sh *` 指令，逐级查看哪个目录占用磁盘空间较大，进入相应的目录，直到找到最精确的文件或目录。

- 最后，结合业务情况等判断对相关文件或目录进行删除，或者购买更大的数据盘分担处理。

### inode 容量满的处理

通过如下操作，解决 inode 容量满的问题。

**查询 inode 使用情况**

Linux 的 inode 节点中，记录了文件的类型、大小、权限、所有者、文件连接的数目、创建时间与更新时间等重要的信息，还有一个比较重要的内容就是指向数据块的指针。一般情况不需要特殊配置，如果存放文件很多，则需要配置。有时磁盘空间有剩余但是不能存放文件，可能是由于 inode 耗尽所致。

1. 执行`df -i`命令，可以查询 inode 的使用情况。

2. 如果 inode 使用率达到或者接近 100%，可以通过以下两种方式进行处理：

3. 清除 inode 占用高的文件或者目录。

4. 修改 inode 数量。

**清除 inode 占用高的文件或者目录**

如果不方便格式化磁盘以增加 inode 数量，可以参考以下步骤，清理 inode 占用量高的文件或者目录。

1. 登录服务器，执行以下命令，分析根目录下的每个二级目录下有多少个文件。`for i in /*; do echo $i; find $i | wc -l; done`。

2. 然后，逐层进入 inode 占用最高的目录，继续执行上述指令，逐步定位占用过高空间的文件或目录，最后进行相应清理。

**修改 inode 数量**

如果不允许清理磁盘中的文件，或者清理后 inode 使用率仍然较高，则需要通过以下步骤，增加 inode 节点数量。

> **说明**：inode 的调整需要重新格式化磁盘，请确保数据已经得到有效备份后，再进行以下操作。

1. 执行以下命令，卸载系统文件。`umount /home`

2. 执行以下命令，重新建立文件系统，指定 inode 节点数。`mkfs.ext3 /dev/xvdb -N 1638400`

3. 执行以下命令，修改 fstab 文件。`vim /etc/fstab`

4. 执行以下命令，查看修改后的 inode 节点数。`dumpe2fs -h /dev/xvdb | grep node`。

### 僵尸文件分析与删除

如果磁盘和 inode 都没有问题，则需要查看是否存在未被清除句柄的僵尸文件。这些文件实际上已经被删除，但是有服务程序在使用这些文件，导致这些文件一直被占用，无法释放磁盘空间。如果这些文件过多，会占用很大的磁盘空间。参考以下步骤查看并删除僵尸文件。

1. 远程登录服务器。

2. 执行以下命令，安装 lsof。`yum install lsof -y`

3. 执行以下命令，查看僵尸文件占用情况。`lsof |grep delete | more`系统显示类似如下。

4. 如果僵尸文件过多，会占用很大的磁盘空间。可以通过以下方法释放句柄，以清除僵尸文件。

5. 重启服务器，验证效果。重启服务器，系统会退出现有的进程，开机后重新加载。该过程会释放调用的 deleted 文件的句柄。

6. 根据`lsof`命令列出的 pid 进程号，使用`kill`命令正常停止或结束占用这些文件的服务进程。

### 挂载点覆盖

先取消磁盘挂载，再检查原挂载目录下的空间占用情况。

## 适用于

- 云服务器 ECS

find / -xdev -printf '%h\n' | sort | uniq -c | sort -k 1 -n

今天 login server 的一个网站，发现 login 后没有生成 session。根据以往经验，一般是空间已满导致 session 文件生成失败。
```bash
df -h
Filesystem                    Size  Used Avail Use% Mounted on
/dev/mapper/dev01-root         75G   58G   14G  82% /
udev                          2.0G  4.0K  2.0G   1% /dev
tmpfs                         396M  292K  396M   1% /run
none                          5.0M     0  5.0M   0% /run/lock
none                          2.0G  4.0K  2.0G   1% /run/shm
/dev/sda1                     228M  149M   68M  69% /boot
```

空间剩余 14G，可以排除空间已满的情况。导致文件生成失败还有另一个原因，就是文件索引节点 inode 已满。

```bash
df -i
Filesystem                    Inodes   IUsed  IFree IUse% Mounted on
/dev/mapper/dev01-root       4964352 4964352      0  100% /
udev                          503779     440 503339    1% /dev
tmpfs                         506183     353 505830    1% /run
none                          506183       5 506178    1% /run/lock
none                          506183       2 506181    1% /run/shm
/dev/sda1                     124496     255 124241    1% /boot
```

inodes 占用 100%，果然是这个问题。

解决方法：删除无用的临时文件，释放 inode。

查找发现 /tmp 目录下有很多 sess_xxxxx 的 session 临时文件。

ls -lt /tmp | wc -l4011517

进入/tmp 目录，执行 find -exec 命令

sudo find /tmp -type f -exec rm {} ;

如果使用 rm \*，有可能因为文件数量太多而出现 Argument list too long 错误，关于 Argument list too long 错误可以参考《linux Argument list too long 错误解决方法》

除了/tmp 的临时文件外，0 字节的文件也会占用 inode，应该也释放。

遍历寻找 0 字节的文件，并删除。

sudo find /home -type f -size 0 -exec rm {} ;

删除后，inode 的使用量减少为 19%，可以正常使用了。

```bash
df -i
Filesystem                    Inodes  IUsed   IFree IUse% Mounted on
/dev/mapper/dev01-root       4964352 940835 4023517   19% /
udev                          503779    440  503339    1% /dev
tmpfs                         506183    353  505830    1% /run
none                          506183      5  506178    1% /run/lock
none                          506183      2  506181    1% /run/shm
/dev/sda1                     124496    255  124241    1% /boot
```
