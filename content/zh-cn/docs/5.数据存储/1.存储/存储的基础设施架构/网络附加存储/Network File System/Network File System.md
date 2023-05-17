---
title: Network File System
---

# 概述

> 参考：
> - [Wiki，Nework_File_System](https://en.wikipedia.org/wiki/Network_File_System)

**Network File System(网络文件系统，简称 NFS)** 是让客户端通过网络访问不同主机上磁盘里的数据，主要用在类 Unix 系统上实现文件共享的一种方法。 本例演示 CentOS 7 下安装和配置 NFS 的基本步骤。

# nfs-ganesha

> 参考：
> - [GitHub 项目，nfs-ganesha/nfs-ganesha](https://github.com/nfs-ganesha/nfs-ganesha)

NFS-Ganesha 是一个 NFSv3、v4、v4.1 文件服务器，在大多数 UNIX/Linux 系统上以用户模式运行。

常见用法：
- 把 Ceph 对象存储转成符合 POSIX 规范的文件系统挂载到 Linux 里

已知问题：
- 麒麟系统+长城服务器挂载 ceph rgw 后，系统严重卡顿
