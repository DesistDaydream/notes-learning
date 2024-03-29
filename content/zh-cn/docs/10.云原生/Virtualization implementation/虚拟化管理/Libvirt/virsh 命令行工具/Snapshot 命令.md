---
title: Snapshot 命令
---

# 概述

> 参考：
>
> - [官方 Manual(手册)，SNAPSHOT COMMANDS](https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-commands)

# snapshot-create - 从 XML 文件中创建一个 domain 的快照

https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-create

这个命令可以在[虚拟机迁移](/docs/10.云原生/Virtualization%20implementation/虚拟化管理/虚拟化管理案例/虚拟机迁移.md)时为虚拟机还原快照的元数据。

## Syntax(语法)

snapshot-create DOMAIN [OPTIONS]

**OPTIONS**

为虚拟机还原名为 base 快照的元数据

- virsh snapshot-create --redefine --xmlfile /var/lib/libvirt/qemu/snapshot/tj-test-spst-node-1/base.xml tj-test-spst-node-1

## EXAMPLE

- 使用 base.xml 文件，为虚拟机 master 创建一个快照
  - virsh snapshot-create master base.xml

# snapshot-create-as - 从一组参数中创建一个 domain 的快照

https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-create-as

EXAMPLE

- 为虚拟机 master 创建一个当前状态的快照，名字为 base
  - virsh snapshot-create-as master --name base

# snapshot-current

https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-current

# snapshot-edit - 编辑指定的 domain 的快照的 XML 文件

https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-edit

# snapshot-info

https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-info

# snapshot-list - 列出指定 domain 的快照

https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-list

# snapshot-dumpxml

https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-dumpxml

# snapshot-parent

https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-parent

# snapshot-revert - 恢复一个 domain 到其上的一个快照

https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-revert

# snapshot-delete - 删除指定的 domain 的快照

https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-delete
