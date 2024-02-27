---
title: "扩大 KVM 虚拟机 image 镜像"
linkTitle: "扩大 KVM 虚拟机 image 镜像"
weight: 2
---

# 概述

直接使用 qemu-img resize 的时候一定要先备份镜像。先使用 virsh shutdown VM 关闭镜像再进行如下操作：

- 给原始系统文件添加磁盘空间
    - qemu-img resize centos7-baseImage-50G.qcow2 +500G
- 进入虚拟机操作新建硬盘分区
    - parted /dev/vda mkpart primary XXX 100%(XX 改为指定的最后一块扇区的容量)
- 扩容 LVM，注意/dev/vda4 设备为真实情况的设备名称，注意修改
    - partprobe(可能还需要重启)
    - pvcreate /dev/vda4
    - vgextend vg0 /dev/vda4
    - lvextend -l+100%FREE /dev/mapper/vg0-root
        - 这里可以使用 lvextend -l +单元数量 /dev/mapper/vg0-lv101
- 扩展文件系统大小，注意 /dev/mapper/vg0-root 修改为真实情况的分区路径
    - XFS 文件系统
        - xfs_growfs /dev/mapper/vg0-root
    - EXT 文件系统
        - resize2fs /dev/mapper/vg0-root (-f 强制扩展)