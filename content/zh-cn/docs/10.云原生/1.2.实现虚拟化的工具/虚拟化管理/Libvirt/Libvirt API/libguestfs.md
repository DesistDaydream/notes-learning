---
title: libguestfs
---

# 概述

> 参考：
> - [GitHub 项目，libguestfs/libguestfs](https://github.com/libguestfs/libguestfs)
> - [官网](https://libguestfs.org/)

Libguestfs 是用于访问和修改虚拟机磁盘映像的库和工具。

## 常见问题

> 参考：
> - <https://access.redhat.com/solutions/4073061>
> - <https://wandering-wolf.tistory.com/entry/Centos-7-KVM-%EC%97%90%EC%84%9C-rhel-8-vm-virt-sysperp-error>
> - <https://dovangiang.wordpress.com/2021/08/06/errorcentos-mount-mount-exited-with-status-32-mount-wrong-fs-type-bad-option-bad-superblock/>

注意：CentOS7 宿主机上想要挂载 Ubuntu 20.04 虚拟机的 xfs 格式的文件系统是不行的，内核版本低不支持，报错如下：

```bash
~]# guestmount --rw -a /var/lib/libvirt/images/common-ubuntu-test.bj-test.qcow2 -m /dev/ubuntu-vg/lv-0 /mnt/test
libguestfs: error: mount_options: mount exited with status 32: mount: wrong fs type, bad option, bad superblock on /dev/mapper/ubuntu--vg-lv--0,
       missing codepage or helper program, or other error

       In some cases useful info is found in syslog - try
       dmesg | tail or so.
guestmount: ‘/dev/ubuntu-vg/lv-0’ could not be mounted.
guestmount: Did you mean to mount one of these filesystems?
guestmount: 	/dev/sda1 (unknown)
guestmount: 	/dev/sda2 (xfs)
guestmount: 	/dev/ubuntu-vg/lv-0 (xfs)
```

其他依赖挂载的命令都是这种报错，比如 virt-sysprep、virt-edit 等等。

如果在 CentOS 7 宿主机上挂载 CentOS 7 虚拟机中的磁盘(或者挂载 ext4 的 Ubuntu)，则是可以的：

```bash
~]# guestmount --rw -a /var/lib/libvirt/images/common-centos-test.bj-test.qcow2 -m /dev/vg1/root /mnt/test
~]# ls /mnt/test/
bin  boot  dev  etc  home  lib  lib64  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
```

# guestmount

> 参考：
> - [官方 Manual(手册)，guestmount](https://libguestfs.org/guestmount.1.html)

# guestunmount

<https://libguestfs.org/guestunmount.1.html>

# virt-edit

> 参考：
> - [Manual(手册),virt-edit(1)](https://libguestfs.org/virt-edit.1.html)

# virt-sysprep

> 参考：
> - [Manual(手册)，virt-sysprep(1)](https://libguestfs.org/virt-sysprep.1.html)
> - <https://www.cnblogs.com/qiuhom-1874/p/13547752.html>

## Syntax(语法)

**virt-sysprep \[OPTIONS] -d DomainName**
**virt-sysprep \[OPTIONS] -a DISK.img \[-a DISK.img ...]**

OPTIONS

- **--copy-in <LocalPath:RemoteDir>** # 将本地宿主机上的 LocalPath 文件拷贝到虚拟机的 RemoteDir 目录中。
  - 注意：RemoteDir 必须已存在。不能使用通配符

EXAMPLE

- 将宿主机上的当前目录中的 test.log 文件拷贝到 ubuntu-2004 虚拟机的 /root/ 目录下
  - `virt-sysprep --copy-in test.log:/root/ -d ubuntu-2004`

# virt-cat

EXAMPLE

- virt-cat -d test /etc/sysconfig/network-scripts/ifcfg-eth0 # 查看名为 test 的 VM 的 ifcfg-eth0 这个文件
