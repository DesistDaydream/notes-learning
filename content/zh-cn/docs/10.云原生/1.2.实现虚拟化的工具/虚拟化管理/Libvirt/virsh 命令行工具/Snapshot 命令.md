---
title: Snapshot 命令
---

# 概述

> 参考：
> - [官方 Manual(手册)，SNAPSHOT COMMANDS](https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-commands)

# [snapshot-create](https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-create) # 从 XML 文件中创建一个 domain 的快照

EXAMPLE

- 使用 base.xml 文件，为虚拟机 master 创建一个快照
  - virsh snapshot-create master base.xml

# [snapshot-create-as](https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-create-as) # 从一组参数中创建一个 domain 的快照

EXAMPLE

- 为虚拟机 master 创建一个当前状态的快照，名字为 base
  - virsh snapshot-create-as master --name base

# [snapshot-current](https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-current)

# [snapshot-edit](https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-edit) # 编辑指定的 domain 的快照的 XML 文件

# [snapshot-info](https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-info)

# [snapshot-list](https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-list) # 列出指定 domain 的快照

# [snapshot-dumpxml](https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-dumpxml)

# [snapshot-parent](https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-parent)

# [snapshot-revert](https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-revert) # 恢复一个 domain 到其上的一个快照

# [snapshot-delete](https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#snapshot-delete) # 删除指定的 domain 的快照
