---
title: virsh 命令行工具
---

# 概述

> 参考：
> - [官方文档，手册-virsh](https://libvirt.org/manpages/virsh.html)
>   - [GitHub 位置，libvirt/libvirt/docs/manpages/virsh.rst](https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst)

virsh 是 libvirt 核心发行版的一部分，通过 libvirt 的 API 管理虚拟机的命令行工具。
virsh 有两种使用方式

- virsh + 子命令
- 交互式 shell。不加任何子命令时，进入交互式 shell

# Syntax(语法)

**virsh \[OPTION]... \[COMMAND_STRING]**
**virsh \[OPTION]... COMMAND \[ARG]...**
注意：

- DOMAIN 的说明：libvirt 使用 domain 代指 VM，所有关于 domain 的描述都可以理解为 VM 或者 instance
- 当指定要操作某个特定 DOMAIN 的时候，可以使用该 DOMAIN 的 ID 号、NAME、UUID，三者任选其一
- 想要查看 VM 的信息，可以在 virsh 的 shell 中使用 help 命令，找 list 和 info 的关键字，help 中是以不同类型的命令进行分类的比如 DOMAIN 相关的，监控信息，网络存储等

## OPTIONS

- -c | --connect=URI hypervisor connection URI
- -d | --debug=NUM debug level \[0-4]
- -e | --escape <char> set escape sequence for console
- -k | --keepalive-interval=NUM #keepalive interval in seconds, 0 for disable
- -K | --keepalive-count=NUM #number of possible missed keepalive messages
- -l | --log=FILE output logging to file
- -q | --quiet quiet mode
- -r | --readonly connect readonly
- -t | --timing print timing information

# COMMAND

Note：其中各种命令用法，详见 virsh 命令行工具目录下每个子命令的专题文章

## [Generic commands](/docs/IT学习笔记/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/virsh%20 命令行工具/通用命令.md 命令行工具/通用命令.md)(通用命令)

通用命令与 domain 无关

## [Domain commands](/docs/IT学习笔记/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/virsh%20 命令行工具/Domain%20 命令.md 命令行工具/Domain 命令.md)(虚拟机命令)

The following commands manipulate domains directly, as stated previously most commands take domain as the first parameter. The domain can be specified as a short integer, a name or a full UUID.

## [Device commands](/docs/IT学习笔记/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/virsh%20 命令行工具/Device%20 命令.md 命令行工具/Device 命令.md)(设备命令)

Device 命令用以控制与 domains 关联的设备。The domain can be specified as a short integer, a name or a full UUID. To better understand the values allowed as options for the command reading the documentation at <https://libvirt.org/formatdomain.html> on the format of the device sections to get the most accurate set of accepted values.

## [NodeDev commands](https://libvirt.org/manpages/virsh.html#nodedev-commands)()

## [Virtual Network commands](/docs/IT学习笔记/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/virsh%20 命令行工具/Virtual%20Network%20 命令.md 命令行工具/Virtual Network 命令.md)(虚拟网络命令)

The following commands manipulate networks. Libvirt has the capability to define virtual networks which can then be used by domains and linked to actual network devices. For more detailed information about this feature see the documentation at <https://libvirt.org/formatnetwork.html> . Many of the commands for virtual networks are similar to the ones used for domains, but the way to name a virtual network is either by its name or UUID.

## NETWORK PORT COMMANDS

## [Interface commands](/docs/IT学习笔记/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/virsh%20 命令行工具/Interface%20 命令.md 命令行工具/Interface 命令.md)(接口命令)

The following commands manipulate host interfaces. Often, these host interfaces can then be used by name within domain <interface> elements (such as a system-created bridge interface), but there is no requirement that host interfaces be tied to any particular guest configuration XML at all.
Many of the commands for host interfaces are similar to the ones used for domains, and the way to name an interface is either by its name or its MAC address. However, using a MAC address for an _iface_ argument only works when that address is unique (if an interface and a bridge share the same MAC address, which is often the case, then using that MAC address results in an error due to ambiguity, and you must resort to a name instead).

## STORAGE POOL COMMANDS

1. 一个 storage pool 包括很多 storage volume，storage volume 有单独的一组命令进行管理
2. storage pool 就是存放 s torage volume 的地方，是一个目录，详见 1.5.Storage Virtualization.note 中的 kvm 的存储虚拟化
3. 存储池默认目录为/var/lib/libvirt/images/，这个目录会存放所有 VM 的文件，通过 libvirt 创建完虚拟机后生成的 image 都会放在 default 目录下

## VOLUME COMMANDS

1. 存储卷就是 VM 上的一块物理硬盘，一个物理硬盘是通过一个文件的形式表现的，修改这个文件，就可以修改这个硬盘的内容。详见 1.5.Storage Virtualization.note 中的 kvm 的存储虚拟化
2. storage volume 管理命令可以调整硬盘大小、类型等，还能增删改查指定的硬盘

## SECRET COMMANDS

## Snapshot commands(快照命令)

The following commands manipulate domain snapshots. Snapshots take the disk, memory, and device state of a domain at a point-of-time, and save it for future use. They have many uses, from saving a "clean" copy of an OS image to saving a domain's state before a potentially destructive operation. Snapshots are identified with a unique name. See <https://libvirt.org/formatsnapshot.html> for documentation of the XML format used to represent properties of snapshots.

## CHECKPOINT COMMANDS

## NWFILTER COMMANDS

## NWFILTER BINDING COMMANDS

## HYPERVISOR-SPECIFIC COMMANDS
