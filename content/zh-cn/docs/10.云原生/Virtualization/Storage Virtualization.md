---
title: Storage Virtualization
linkTitle: Storage Virtualization
weight: 5
---

# 概述

> 参考：
>
> - [Wiki, Storage_virtualization](https://en.wikipedia.org/wiki/Storage_virtualization)

# 存储虚拟化(I/O)

# KVM 模式的存储虚拟化

## 第一种：存储虚拟化是通过存储池（Storage Pool）和卷（Volume）来管理的

1. Storage Pool 是宿主机上可以看到的一片存储空间，可以是多种类型。
2. 文件目录类型的 Storage Pool 。KVM 将宿主机目录 /var/lib/libvirt/images/ 作为默认的 Storage Pool。
3. KVM 将 HOST 目录 /var/lib/libvirt/images/ 作为默认的 Storage Pool
4. KVM 所有可以使用的 Storage Pool 都定义在宿主机的 /etc/libvirt/storage 目录下，每个 Pool 一个 xml 文件，默认有一个 default.xml
5. LVM 类型的 Storage Pool。宿主机上 VG 中的 LV 也可以作为虚拟磁盘分配给虚拟机使用。不过，LV 由于没有磁盘的 MBR 引导记录，不能作为虚拟机的启动盘，只能作为数据盘使用。
6. KVM 还支持 iSCSI，Ceph 等多种类型的 Storage Pool，最常用的就是目录类型，其他类型可以参考文档<http://libvirt.org/storage.html>
7. Volume 是在 Storage Pool 中划分出的一块空间，宿主机将 Volume 分配给虚拟机，Volume 在虚拟机中看到的就是一块硬盘。Volume 是 Storage Pool 目录下面的文件，一个文件就是一个 Volume(使用文件做 Volume 有很多优点：存储方便、移植性好、可复制)。Volume 分为几种类型，类型如下
8. qcow2 # 是推荐使用的格式，QEMU V2 磁盘镜像格式，cow 表示 copy on write，能够节省磁盘空间，支持 AES 加密，支持 zlib 压缩，支持多快照，功能很多。
9. raw # 是默认格式，即原始磁盘镜像格式，移植性好，性能好，但大小固定，不能节省磁盘空间。
10. vmdk # 是 VMWare 的虚拟磁盘格式，也就是说 VMWare 虚机可以直接在 KVM 上 运行。
