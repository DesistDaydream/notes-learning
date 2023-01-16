---
title: XFS 文件系统
---

# 概述

> 参考：
> - [Wiki,XFS](https://en.wikipedia.org/wiki/XFS)

XFS 是由 Silicon Graphics 于 1993 年创建的高性能 64 位 **Journaling File System(日志文件系统)**。

# xfs_repair 命令行工具

官方文档：[xfs_repair man 手册](https://man.cx/xfs_repair)

`**xfs_repair**` 用于修复 XFS 文件系统。

    [root@lichenhao ~]# xfs_info /dev/mapper/vg1-root
    meta-data=/dev/mapper/vg1-root   isize=512    agcount=4, agsize=32735744 blks
             =                       sectsz=512   attr=2, projid32bit=1
             =                       crc=1        finobt=1, sparse=1, rmapbt=0
             =                       reflink=1
    data     =                       bsize=4096   blocks=130942976, imaxpct=25
             =                       sunit=0      swidth=0 blks
    naming   =version 2              bsize=4096   ascii-ci=0, ftype=1
    log      =internal log           bsize=4096   blocks=63937, version=2
             =                       sectsz=512   sunit=0 blks, lazy-count=1
    realtime =none                   extsz=4096   blocks=0, rtextents=0
