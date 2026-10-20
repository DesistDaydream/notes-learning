---
title: "扩大 KVM 虚拟机 image 镜像"
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

## 不进虚拟机扩展非 LVM 磁盘

```bash
cd /var/lib/libvirt/images

cp --sparse=always hdf-01.qcow2 hdf-01.qcow2.pre-resize-$(date +%Y%m%d)

# 镜像变为 V3 格式（V2 不支持带快照扩容）
qemu-img amend -o compat=1.1 hdf-01.qcow2

qemu-img resize hdf-01.qcow2 +100G

# 加载 nbd 模块
modprobe nbd max_part=16
# 把虚拟机磁盘挂载到宿主机的 /dev/nbd0 上
qemu-nbd -c /dev/nbd0 hdf-01.qcow2

parted -s /dev/nbd0 resizepart 4 100%; partprobe /dev/nbd0
```

> [!Attention] 下面的 op4 中的 4 是分区号，由于要扩展 / 分区，所以这里的示例要保证根分区在最后一个盘符，且是第 4 个分区。

```bash
e2fsck -f /dev/nbd0p4 && resize2fs /dev/nbd0p4

# 将虚拟机镜像与 /dev/nbd0 断开
qemu-nbd -d /dev/nbd0
virsh start hdf-01
```

---

若 resizepart 时报错

```bash
警告: 并非所有可用于 /dev/nbd0 的空间都被用到了，您可以修正 GPT 以使用所有的空间 (额外的 209715200 个区块)，还是说要继续使用目前的设置？
错误: 无法满足分区的所有限制。
Warning: 并非所有可用于 /dev/nbd0 的空间都被用到了，您可以修正 GPT 以使用所有的空间 (额外的 209715200 个区块)，还是说要继续使用目前的设置？
```

> [!Quote] 原因：磁盘文件已经扩展到 200G，但 **GPT 备份头还停在旧的磁盘末尾**（100G 处）。使用 `-s` 非交互模式时，它无法回答那句 “Fix/Ignore?” 提示，便按“不修”继续；于是它认为可用末端只到旧的磁盘尾部，`100%` 落在该范围之外，从而报出 “无法满足分区的所有限制”

则执行如下方式修复后继续

```bash
~]# sgdisk -e /dev/nbd0
The operation has completed successfully.
~]# partprobe /dev/nbd0
~]# parted -s /dev/nbd0 resizepart 4 100%
~]# partprobe /dev/nbd0
~]# parted -s /dev/nbd0 unit s print
# ...略
```
