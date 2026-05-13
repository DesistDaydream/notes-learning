---
title: LVM
linkTitle: LVM
created: 2026-05-13T22:18
weight: 100
---

# 概述

> 参考：
>
> - [GitHub 项目,lvmteam/lvm2](https://github.com/lvmteam/lvm2)
> - [Manual(手册),lvm(8)](https://man7.org/linux/man-pages/man8/lvm.8.html)

**Logical Volume Manager(逻辑卷管理器，简称 LVM)** 提供了从物理设备创建 **Virtual Block Devices(虚拟块设备)** 的**工具**。虚拟设备比物理设备更易于管理，并且可以具有超出物理设备自身提供的功能的能力。

- **Volume Group(卷组，简称 VG)** 是一个或多个物理设备的集合，每个设备称为 **Physical Volume(物理卷，简称 PV)**
- **Logical Volume(逻辑卷，简称 LV)** 从 VG 中创建，是可由系统或应用直接使用的虚拟块设备。
  - 根据内核中 **Device Mapper(设备映射，简称 DM)** 实现的算法，LV 中的每个数据块都存储在 VG 中的一个或多个 PV 上。

LVM 创建的 LV 本质上是相当于创建了一个新的物理磁盘，这不说是逻辑上的磁盘。所以称为逻辑卷。在 Linux 中，LVM 表现为一种 device-mapper 类型的磁盘。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/zqtob4/1640420060047-c8503823-746d-41d5-8db8-3c4c5bd218da.png)

上图中，虽然 testvg1-testlv 磁盘是一个 LVM，并且是通过 /dev/sdb 磁盘的分区创建的，但是从整个磁盘管理的角度来看，一个 LV 就是一个物理磁盘，只不过类型不同罢了。

创建 LVM 顺序

1. **Physical Volume(物理卷，简称 PV)** # 相当于一块一块真实的物理磁盘
2. **Volume Group(卷组，简称 VG)** # 把多块物理磁盘组合在一起形成一个组，并在创建 VG 时把整个组的空间划分为一个个默认大小为 4M 的 PE，VG 就相当于一个存储池子，里面有好多好多 PE 块组成了一个大的磁盘空间供 LV 使用
    1. **Pysical Extent(物理区域，简称 PE)** # 把真实的物理磁盘的空间切分为一个个固定大小的块，每个块就是 PE
3. **Logical Volume(逻辑卷，简称 LV)** # 从 VG 中拿出多个 PE 组成逻辑上的磁盘空间，可以把 LV 当成一个分区

> [!Tip] 所以，操作 LV 本质就是操作 PE
>
> <font color="#ff0000">一个 LV 可用磁盘空间的大小取决于分配了多少块 PE</font>
>
> 各种 LVM 的控制命令，e.g. pvmove, lvextend, etc. 命令，本质上都是在操作 PE。哪怕命令选项中看似可以指定空间，但是实际上空间会折算成 PE 数量。

## Syntax(语法)

**lvm \[ COMMAND | FILE ]**

# pv - 物理卷管理命令

> 参考：
>
> - [Manual(手册), pvcreate(8)](https://man7.org/linux/man-pages/man8/pvcreate.8.html)
> - [Manual(手册), pvresize(8)](https://man7.org/linux/man-pages/man8/pvresize.8.html)

## pvcreat

pvcreat 创建 PV，初始化物理卷以供 PVM 使用

**pvcreate POSITION \[OPTIONS]**

OPTIONS

EXAMPLE

- pvcreate /dev/sdb1
- pvcreate /dev/sdb{1,2}

## pvremove - 删除 PV

## pvscan - 扫描 PV

## pvs - 简要显示 PV

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zqtob4/1616167609644-770f62f5-144e-44cb-8b06-feceb8309f0f.jpeg)

- PV # 物理卷名称
- VG # PV 所属的卷组
- Fmt # 物理卷格式
- Attr
- PSize # 物理卷大小
- PFree # 物理卷空余大小（并不是指真实空间，而是有多少富裕的 PE 还没分配）

## pvdisplay - 详细显示 PV

https://man7.org/linux/man-pages/man8/pvdisplay.8.html

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zqtob4/1616167609649-9587c471-5e11-4e45-a7b6-0046a8d2030b.jpeg)

- PV Name # 物理卷名称
- VG Name # PV 所属的卷组名称
- PV Size # 物理卷占用多少空间
- Allocatable # 是否可分配
- PE Size # 每个 PE 占用的空间容量
- Total PE：合计有多少个 PE
- Free PE：有多少空余的 PE
- Allocated PE：分配了多少 PE
- PV UUID：PV 的 UUID 号

**OPTIONS**

- **-m, --maps** # 显示物理扩展区到逻辑卷（LV）及逻辑扩展区的映射关系。

## pvmove - 把一个 pv 上的数据挪动到另一个 pv 上

https://man7.org/linux/man-pages/man8/pvmove.8.html

将 PE 从一个 PV 移动到另一个 PV

# vg - 卷组管理命令

vgcreat,vgs,vgdisplay,vgextend,vgreduce,vgremove,vgrename

## vgcreat

**vgcreat VgName PvPath**

创建 vg 自定义 vg 名 需要加入 VG 的 PV 路径

选项

- **-s** # 指定 PE 大小，默认 4M

应用范例

- vgcreate myvg /dev/sda{1,2} # 以 sdb1 和 sdb2 组合起来创建一个名叫 myvg 的 VG

## vgs

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zqtob4/1616167609659-c21d1b79-70bd-4975-a150-6348c5f0bd1e.jpeg)

- **\#PV** # 该卷组中物理卷的个数 # LV：逻辑卷的个数

## vgdisplay

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zqtob4/1616167609663-ebeba7a4-5acb-4d14-ba10-6eddf182d96e.jpeg)

## vgextend 扩大卷组

**vgextend VgName PvPath**

EXAMPLE

- vgextend myvg /dev/sdc1 # 为卷组 myvg 添加一个/dev/sdc1 的物理卷

## vgreduce 缩小卷组

>[!Attention] 缩减卷组前，要把要缩减掉的 pv 上的数据用 pvmove 挪走)

# lv - 逻辑卷管理命令

> 参考：
>
> - [Manual(手册),lvcreate(8)](https://man7.org/linux/man-pages/man8/lvcreate.8.html)

lvcreat, lvdisplay, lvs, lvextend, lvreduce, lvremove

## lvcreat

**lvcreat \[OPTIONS] VgName \[PvPath]**

创建 lv 从哪个 vg 中创建

OPTIONS

- **-n** # 指定 lv 名字
- **-L** # 指定 lv 的空间大小
- **-p** # 指定访问权限
- **-s** # 创建 lv 的快照 snapshot

EXAMPLE

- **lvcreate -L 3G -n lv1 myvg** # congmyvg 卷组中创建名为 Imyvg 大小为 3G 的 LV
- **lvcreate -s -L3G -n lv1snapshot -p r /dev/myvg/lv1** # 创建名为 lv1snapshot，只有读权限，3G 大小的 lv1 的快照

## lvdisplay

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zqtob4/1616167609656-34a3e1a7-1885-42e0-99db-234afc6e1061.jpeg)

## lvextend

https://man7.org/linux/man-pages/man8/lvextend.8.html

向 LV(逻辑卷) 添加空间

**lvextend \[OPTIONS] LvPath \[PvPath]**

增加 lv 增加哪个 LV \[可以指定从哪个 PV 里给该 LV 增加]

OPTIONS

- **-l, --extens**(`[+]Number[PERCENT]`) # 
- **-L, --size**(`[+]Size[m|UNIT]`) # 指定要增加的空间，增加到 5G 则是-L 5G，增加了 5G 则是-L +5G

EXAMPLE

- lvextend -L 5G /dev/myvg/lv1 # 扩展 lv1 这个逻辑卷到 5G 的容量
- lvextend -L +5G /dev/myvg/lv1 sdc1 # 扩展 lv1 逻辑卷多出 5G 的空间，从 sdc1 中给 lv 提供空间
- lvextend -l+100%FREE /dev/mapper/vg0-root # 扩展卷组中所有空间给 vg0-root 这个逻辑卷

# 其他

## snapshot 快照卷

注意：lvm 创建的快照是基于原始卷来运行的，原始卷修改了多少内容，其修改之前的内容就会被自动复制到快照卷里，快照是按照 PE 复制的，每个 PE 的改变，原始 PE 都会自动备份到快照当中，当快照卷的容量满了，就失去了快照卷了，注意 lvm 创建的快照与 openstack 的实例创建的快照还有 VMware 创建的快照的区别

# 最佳实践

创建逻辑卷：

1. 创建分区，并使之可用

```bash
pvcreate /dev/sda{1,2}
vgcreate testvg1 /dev/sda{1,2}
lvcreate -L 10G -n testlv testvg1
```

扩展逻辑卷：

1. 确定扩展的大小，并确保 lv 所属的 vg 有足够的剩余空间

```bash
pvcreate /DEV/PATH
vgextend VG /DEV/PATH
```

2. 扩展 lv 大小，执行 lvextend 命令，第二条为扩展所有剩余空间到逻辑卷

```bash
lvextend -L [+]SIZE LvPath
lvextend -l+100%FREE /dev/mapper/vg0-root
```

3. 扩展系统文件大小，执行 resize2fs 命令

```bash
xfs_growfs LvPath
```

缩减逻辑卷

1. 卸载卷 umount，并执行强制检测

```bash
e2fsck -f LvPath
```

2. 缩减文件系统大小

```bash
resize2fs LvPath SIZE
```

3. 缩减 lv 大小

```bash
lvreduce -L [-]Siza LvPath
```

## 去掉 lv 中某块物理磁盘，并将去掉的补充到另一个 lv 中

场景如下：由于缺乏思考，某人安装系统时，把本就容量不大的固态硬盘，分出去几百 G，给到了容量本就很大的机械硬盘中。最后导致根目录只有 70G 空间。

我们需要把 sdb 的空间都还给根分区

```bash
~]# lsblk -d -o name,rota,size,model
NAME ROTA   SIZE MODEL
sda     1   3.6T HGST HUS726T4TALE6L4
sdb     0 447.1G SAMSUNG MZ7LH480HAHQ-00005
~]# df -h
Filesystem                  Size  Used Avail Use% Mounted on
devtmpfs                    4.0M     0  4.0M   0% /dev
tmpfs                       126G     0  126G   0% /dev/shm
tmpfs                        51G   61M   51G   1% /run
tmpfs                       4.0M     0  4.0M   0% /sys/fs/cgroup
/dev/mapper/openeuler-root   69G  6.5G   59G  10% /
tmpfs                       126G  8.0K  126G   1% /tmp
/dev/sdb1                   974M  176M  731M  20% /boot
/dev/mapper/openeuler-home  4.0T   64G  3.8T   2% /home
~]# lsblk
NAME               MAJ:MIN RM   SIZE RO TYPE MOUNTPOINTS
sda                  8:0    1   3.6T  0 disk 
└─sda1               8:1    1   3.6T  0 part 
  └─openeuler-home 253:2    0     4T  0 lvm  /home
sdb                  8:16   1 447.1G  0 disk 
├─sdb1               8:17   1     1G  0 part /boot
└─sdb2               8:18   1 446.1G  0 part 
  ├─openeuler-root 253:0    0    70G  0 lvm  /
  ├─openeuler-swap 253:1    0     4G  0 lvm  [SWAP]
  └─openeuler-home 253:2    0     4T  0 lvm  /home
~]# lvs
  LV   VG        Attr       LSize  Pool Origin Data%  Meta%  Move Log Cpy%Sync Convert
  home openeuler -wi-ao----  4.00t                                                    
  root openeuler -wi-ao---- 70.00g                                                    
  swap openeuler -wi-ao----  4.00g                                                    
~]# pvs
  PV         VG        Fmt  Attr PSize    PFree
  /dev/sda1  openeuler lvm2 a--    <3.64t    0 
  /dev/sdb2  openeuler lvm2 a--  <446.13g    0 
~]# vgs
  VG        #PV #LV #SN Attr   VSize VFree
  openeuler   2   3   0 wz--n- 4.07t    0
```

确认 PE 数据分布

```bash
~]# pvdisplay -m /dev/sdb2 /dev/sda1
  --- Physical volume ---
  PV Name               /dev/sdb2
  VG Name               openeuler
  PV Size               <446.13 GiB / not usable 0   
  Allocatable           yes (but full)
  PE Size               4.00 MiB
  Total PE              114209
  Free PE               0
  Allocated PE          114209
  PV UUID               cG8yY3-43rd-C6a6-W6ae-JNeT-3BBt-1UWXNL
   
  --- Physical Segments ---
  Physical extent 0 to 1023:
    Logical volume      /dev/openeuler/swap
    Logical extents     0 to 1023
  Physical extent 1024 to 96288:
    Logical volume      /dev/openeuler/home
    Logical extents     953861 to 1049125
  Physical extent 96289 to 114208:
    Logical volume      /dev/openeuler/root
    Logical extents     0 to 17919
   
  --- Physical volume ---
  PV Name               /dev/sda1
  VG Name               openeuler
  PV Size               <3.64 TiB / not usable 2.00 MiB
  Allocatable           yes (but full)
  PE Size               4.00 MiB
  Total PE              953861
  Free PE               0
  Allocated PE          953861
  PV UUID               vnqB0B-IUJU-xcx3-nIRF-NKun-uYnm-8WPswx
   
  --- Physical Segments ---
  Physical extent 0 to 953860:
    Logical volume      /dev/openeuler/home
    Logical extents     0 to 953860
```

我们需要把 /dev/sdb2 中 1024 - 96288 这些 PE 取出来还给 /dev/sda1

---

**卸载分区**

检查目录是否有进程占用

```bash
lsof +D /home
```

干掉所有进程后，卸载 /home 目录

```bash
umount /home
```

**检查文件系统**

```bash
e2fsck -f /dev/mapper/openeuler-home
```

---

**缩减 LV，留出可以容纳 PE 的位置**

```bash
# 缩小文件系统到 400G（当前用了 64G，留足余量）
resize2fs /dev/mapper/openeuler-home 400G

# 缩小 LV 到 400G
lvreduce -L 400G /dev/openeuler/home
```

---

**迁移数据（i.e. 迁移 PE）**

将 /dev/sdb2 中的 1024-96288 这些 PE 移动到其他可用的 PV 中

```bash
pvmove /dev/sdb2:1024-96288
```

> [!Note] 有时候 lvreduce 之后，会发现 /dev/sdb2 上的这部分空间自动变为 Free 了，这是因为在缩减 LV 时，发现这些 PE 中没有数据，也就自动释放了。

检查 /dev/sdb2 上的 PV 可以看到 1024-96288 已经 FREE 了

```bash
~]# pvdisplay -m /dev/sdb2
  --- Physical volume ---
  PV Name               /dev/sdb2
  VG Name               openeuler
  PV Size               <446.13 GiB / not usable 0   
  Allocatable           yes 
  PE Size               4.00 MiB
  Total PE              114209
  Free PE               95265
  Allocated PE          18944
  PV UUID               cG8yY3-43rd-C6a6-W6ae-JNeT-3BBt-1UWXNL

  --- Physical Segments ---
  Physical extent 0 to 1023:
    Logical volume      /dev/openeuler/swap
    Logical extents     0 to 1023
  Physical extent 1024 to 96288:
    FREE
  Physical extent 96289 to 114208:
    Logical volume      /dev/openeuler/root
    Logical extents     0 to 17919
```

---

**扩展 `/` LV**

将 /dev/sdb2 中所有空闲的 PE 加入到 /dev/openeuler/root LV 中

```bash
lvextend -l+100%FREE /dev/openeuler/root /dev/sdb2
# 或者使用如下命令：
# lvextend -l +95265 /dev/openeuler/root
# 说明：95265 是 pvdisplay -m /dev/sdb2 命令中查到的 Free PE。该命令会就近分配，并不会把 /dev/sda1 中的 PE 也扩展进来。
```

**恢复文件系统空间**

```bash
resize2fs /dev/mapper/openeuler-root
```

---

**还原 /dev/openeuler/home LV 中的空间**

```bash
lvextend -l+100%FREE /dev/openeuler/home /dev/sda1
resize2fs /dev/mapper/openeuler-home
```

重新挂载

```bash
mount /home
```

---

**完成**

美中不足的是 /dev/openeuler/root LV 由两部分组成，这两部分没法合并了。

```bash
~]# df -h
Filesystem                  Size  Used Avail Use% Mounted on
devtmpfs                    4.0M     0  4.0M   0% /dev
tmpfs                       126G   12K  126G   1% /dev/shm
tmpfs                        51G   61M   51G   1% /run
tmpfs                       4.0M     0  4.0M   0% /sys/fs/cgroup
/dev/mapper/openeuler-root  435G  6.5G  410G   2% /
tmpfs                       126G  8.0K  126G   1% /tmp
/dev/sdb1                   974M  176M  731M  20% /boot
/dev/mapper/openeuler-home  3.7T   67G  3.4T   2% /home
~]# lsblk
NAME               MAJ:MIN RM   SIZE RO TYPE MOUNTPOINTS
sda                  8:0    1   3.6T  0 disk 
└─sda1               8:1    1   3.6T  0 part 
  └─openeuler-home 253:2    0   3.6T  0 lvm  /home
sdb                  8:16   1 447.1G  0 disk 
├─sdb1               8:17   1     1G  0 part /boot
└─sdb2               8:18   1 446.1G  0 part 
  ├─openeuler-root 253:0    0 442.1G  0 lvm  /
  └─openeuler-swap 253:1    0     4G  0 lvm  
~]# pvdisplay -m /dev/sda1 /dev/sdb2 
  --- Physical volume ---
  PV Name               /dev/sdb2
  VG Name               openeuler
  PV Size               <446.13 GiB / not usable 0   
  Allocatable           yes (but full)
  PE Size               4.00 MiB
  Total PE              114209
  Free PE               0
  Allocated PE          114209
  PV UUID               cG8yY3-43rd-C6a6-W6ae-JNeT-3BBt-1UWXNL
   
  --- Physical Segments ---
  Physical extent 0 to 1023:
    Logical volume      /dev/openeuler/swap
    Logical extents     0 to 1023
  Physical extent 1024 to 96288:
    Logical volume      /dev/openeuler/root
    Logical extents     17920 to 113184
  Physical extent 96289 to 114208:
    Logical volume      /dev/openeuler/root
    Logical extents     0 to 17919
   
  --- Physical volume ---
  PV Name               /dev/sda1
  VG Name               openeuler
  PV Size               <3.64 TiB / not usable 2.00 MiB
  Allocatable           yes (but full)
  PE Size               4.00 MiB
  Total PE              953861
  Free PE               0
  Allocated PE          953861
  PV UUID               vnqB0B-IUJU-xcx3-nIRF-NKun-uYnm-8WPswx
   
  --- Physical Segments ---
  Physical extent 0 to 953860:
    Logical volume      /dev/openeuler/home
    Logical extents     0 to 953860
```