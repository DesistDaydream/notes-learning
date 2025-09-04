---
title: Kernel 参数
linkTitle: Kernel 参数
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，Linux 内核用户和管理员指南 - /proc/sys 文档](https://www.kernel.org/doc/html/latest/admin-guide/sysctl/index.html)

内核参数是以 key/value 的方式储存在 [sysfs](/docs/1.操作系统/Kernel/Filesystem/特殊文件系统/sysfs.md) 文件系统中。key 就是 **/proc/sys/** 目录下的某个文件，value 就是该文件的内容。

比如 net.ipv4.ip_forward 这个 key，就在 **/proc/sys/net/ipv4/** 目录下。以`.`分隔的就是字符就是目录名，最后一个字段就是某某目录下的文件名。

可以通过修改 **/proc/sys/** 目录下的文件中的值来修改内核的参数。每个文件就是内核的一种功能，文件中的内容就是该内核功能的参数。

注意事项：

- 一般情况下，内核参数的 0 和 1 这两个值表示如下含义：
  - **0** 表示 **否**，即错误、拒绝、关闭等等
  - **1** 表示 **是**，即正确、允许、开启等等
- sysctl 工具用来配置与显示在 /proc/sys 目录中的内核参数．如果想使参数长期保存，可以通过编辑 /etc/sysctl.conf 文件来实现。
- 修改 /proc 下内核参数文件内容，不能使用编辑器来修改内核参数文件，理由是由于内核随时可能更改这些文件中的任意一个，另外，这些内核参数文件都是虚拟文件，实际中不存在，因此不能使用编辑器进行编辑，而是使用 echo 命令，然后从命令行将输出重定向至 /proc 下所选定的文件中。参数修改后立即生效，但是重启系统后，该参数又恢复成默认值。因此，想永久更改内核参数，需要修改 /etc/sysctl.conf 文件。
  - `echo 1 > /proc/sys/net/ipv4/ip_forward`
  - `sysctl -w net.ipv4.ip_forward=1`
  - 永久的方法：
    - `echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf`
- 如果想使参数马上生效，也可以执行如下命令 `sysctl  -p`

# /proc/sys 目录的组成

**/proc/sys/** 目录下的每一个子目录，都表示一种内核参数的分类，大体可以分为如下几类：

- **./abi/** # execution domains & personalities
  - [Documentation for /proc/sys/abi/](https://www.kernel.org/doc/html/latest/admin-guide/sysctl/abi.html)
- **./debug/** # 空
- **./dev/** # device specific information (eg dev/cdrom/info)
- **./fs/** # specific filesystems filehandle, inode, dentry and quota tuning binfmt_misc \<Kernel Support for miscellaneous Binary Formats (binfmt_misc)>
  - 详见 [fs(文件系统相关参数)](/docs/1.操作系统/Kernel/Linux%20Kernel/Kernel%20参数/fs(文件系统相关参数).md)
- **./kernel/** # global kernel info / tuning miscellaneous stuff
  - [Documentation for /proc/sys/kernel/](https://www.kernel.org/doc/html/latest/admin-guide/sysctl/kernel.html)
- **./net/** # 网络相关的内核参数
  - 详见 [net(网络相关参数)](/docs/1.操作系统/Kernel/Linux%20Kernel/Kernel%20参数/net(网络相关参数)/net(网络相关参数).md)
- **./proc/** # 空
- **./sunrpc/** # SUN Remote Procedure Call (NFS)
  - [Documentation for /proc/sys/sunrpc/](https://www.kernel.org/doc/html/latest/admin-guide/sysctl/sunrpc.html)
- **./vm/** # 内存管理调整缓冲区和高速缓存管理。memory management tuning buffer and cache management
  - 详见 [vm(内存相关参数)](/docs/1.操作系统/Kernel/Linux%20Kernel/Kernel%20参数/vm(内存相关参数).md)
- **./user/** # Per user per user namespace limits
  - [Documentation for /proc/sys/user/](https://www.kernel.org/doc/html/latest/admin-guide/sysctl/user.html)
