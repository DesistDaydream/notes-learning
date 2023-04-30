---
title: Device File System
---

# 概述

> 参考：
> - [Manual(手册),devfs(5)](https://man.cx/devfs)
> - [IBM，开发工作-Linux-DM 机制](https://www.ibm.com/developerworks/cn/linux/l-devmapper/index.html)

**Device File System(设备文件系统，简称 devfs)**，提供对全局文件系统名称空间中内核设备名称空间的访问。`一般挂载到 /dev 目录`。

这个文件系统包含一些目录、链接、符号链接和设备，其中一些是可写的。在 chroot 环境中，可以使用 devfs 创建一个新的/dev 挂载点。

The [mknod(8)](<https://man.cx/mknod(8)>) 工具可用于恢复 devfs 下已删除的设备。

# /dev/dm-\*

> 参考：
> - [Wiki,Device mapper](https://en.wikipedia.org/wiki/Device_mapper)

**Device Mapper(设备映射，简称 dm)**，是一个由 Linux 内核提供的框架，用于将物理块设备映射到更高级别的虚拟块设备。dm 是 LVM、软 Raid、dm-crypt 磁盘加密的基础。

dm 通过将将数据从虚拟块设备传递到另一个块设备来工作。数据也可以再过渡中进行修改，例如，在设备映射器提供磁盘加密或模拟不可靠硬件行为的情况下，可以执行此操作。

### dmsetup 命令行工具

dmsetup ls # 列出 dm 设备

```bash
~]# dmsetup ls
vg1-swap	(253:1)
vg1-root	(253:0)
```

其中 253 后面的数字，就是 dm-X 那个 X。所以 dm-0 对应 vg1-root 这个设备。使用 lsblk 命令可以看到 dm 与 块设备的关联关系。
