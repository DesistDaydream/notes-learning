---
title: Huge Pages
linkTitle: Huge Pages
date: 2024-05-22T14:26
weight: 20
---

# 概述

> 参考:
>
> - [Linux Kernel 文档，管理员指南 - 内存管理 - 概念 - Huge pages](https://www.kernel.org/doc/html/latest/admin-guide/mm/concepts.html#huge-pages)
> - [Wiki, Huge pages](https://en.wikipedia.org/wiki/Huge_pages)

**Huge Pages(大页)** 是指每个 Page 的容量都远超默认大小(4 KiB)的 Page。比如 2 MiB、1 GiB、etc. 都是常见的大页内存空间中的每页容量。Huge Pages Memory 则是指每个 Page 的容量都超过 4 KiB 的内存的统称。

在 x86 架构上，可以使用 第二级 和 第三级 页表中的条目来映射 2MiB 甚至 1GiB 的 Page。

HugePages 可以减少页表开销、减轻 TLB 压力并提高 TLB 的命中率、减轻内存数据查询压力、避免使用 Swap 降低性能。但是前提是保证使用大页的程序可以完善得利用大页，否则就会造成内存的极大浪费。

> [!Notes] 为什么已经分页了还要用大页？
> 如果一个程序（比如数据库），把大量数据加载到内存中，这时候其查询的数据量一定远超 TLB 的容量，这必然会导致 TLB 的未命中急速上升，严重影响性能。还有很多其他的方面就不一一举例了。
>
> 所以大页并不是所有程序都适用的，而是针对特定场景，需要处理大量数据，亲自管理内存的程序，才要配置大页。比如 [DPDK](/docs/4.数据通信/DPDK/DPDK.md) 处理流量也需要使用大页的内存空间

Linux Kernel 中有两种机制可以实现 物理内存 与 Huge Pages 的映射

- **HugeTLB filesystem** # 大页文件系统，简称 **[Hugetlbfs](#HugeTLB%20FS)**。在此文件系统中创建的文件，数据驻留在内存中并使用大页进行映射。
- **Transparent HugePages** # 透明大页，简称 [**THP**](#THP)。

# HugeTLB FS

> 参考:
>
> - [Linux Kernel 文档，管理员指南 - 内存管理 - HugeTLB Pages](https://www.kernel.org/doc/html/latest/admin-guide/mm/hugetlbpage.html)

HugeTLB Filesystem 是一种特殊的 [Filesystem](/docs/1.操作系统/Kernel/Filesystem/Filesystem.md)

```bash
~]# mount -t hugetlbfs
hugetlbfs on /dev/hugepages type hugetlbfs (rw,relatime,pagesize=1024M)
~]# ls /dev/hugepages/
libvirt       rtemap_16388  rtemap_16402  rtemap_17  rtemap_3      rtemap_65544  rtemap_65558  rtemap_73731  rtemap_73745  rtemap_8202
......略
~]# cat /proc/filesystems | grep huge
nodev   hugetlbfs
```

x86 CPU 通常支持 4K 和 2M（如果架构支持，则为 1G）PageSize，ia64 架构支持多种 PageSize 4K、8K、64K、256K、1M、4M、16M、256M，ppc64 支持 4K 和 16M。 TLB 是虚拟到物理转换的缓存。通常，这是处理器上非常稀缺的资源。操作系统尝试充分利用有限数量的 TLB 资源。现在，随着越来越大的物理内存（几 GB）变得越来越可用，这种优化变得更加重要。

## 关联文件与配置

**/dev/hugepages/** # HugeTLB 文件系统的默认挂载目录

有多种方式可以设置 HugePages（TODO: 优先级？）

- **内核引导参数** # 可以在系统引导配置中（e.g. grub2.cfg、etc.）中添加一些参数以设定大页的信息
  - hugepagesz # HugePages 的页容量
  - hugepages # HugePages 的页数量
  - default_hugepagesz # 默认的 HugePages 容量
  - etc.
- **/proc/sys/vm/\*hugepages\*** # 仅为了向后兼容，保留这些文件。这些 HugePage 的用户空间接口已在 [sysfs](/docs/1.操作系统/Kernel/Filesystem/特殊文件系统/sysfs.md) 中的 /sys/kernel/mm/hugepages/ 实现。
- **/sys/kernel/mm/hugepages/** # 由于要在运行时支持多个，该目录实现了 /proc/sys/vm/ 中大部分关于 HugePage 的用户空间接口
  - **`./hugepages-${SIZE}kB/`** # 详见下文 [通用目录](#通用目录) 中的 `hugepages-${SIZE}kB` 介绍
- **`/sys/devices/system/node/node[0-9]*/hugepages/`** # 针对 [NUMA](/docs/1.操作系统/Kernel/Memory/NUMA.md) 为每个 Node 设置 HugePages
  - **`./hugepages-${SIZE}kB/`** # 详见下文 [通用目录](#通用目录) 中的 `hugepages-${SIZE}kB` 介绍

Notes: 关于 NUMA 的 HugePages 的说明

- `/sys/bus/node/devices/node[0-9]*` 目录将会链接到 `/sys/devices/system/node/node[0-9]*` 目录
- sysfs 中根大页控制目录的内容子集将复制到内存位于以下位置的每个 NUMA 节点的每个系统设备下。

### 通用目录

**./hugepages-${SIZE}kB/** # 特定于指定 SIZE(HugePage 的 **页容量**) 的 HugePage 用户空间接口。SIZE 通常默认有两个值: 1048576 和 2048

- **nr_hugepages** # HugePage 的页数量。当任务释放时，“持久”大页将返回到大页池。具有 root 权限的用户可以通过增加或减少 nr_hugepages 的值来动态分配更多或释放一些持久性大页面。

## 使用 HugePages

https://www.kernel.org/doc/html/latest/admin-guide/mm/hugetlbpage.html#using-huge-pages

程序想要使用 HugePage 的内存，需要在申请内存的系统调用（e.g. mmap()）时添加大页相关参数（e.g. MAP_HUGETLB、etc.）即可。

也可以通过标准 SYSV 共享内存系统调用（shmget、shmat）使用 Linux 的 Hugepages。

如果用户应用程序要使用 mmap 系统调用请求大页，则需要先挂载 Hugetlbfs 类型的文件系统：

```bash
mount -t hugetlbfs \
      -o uid=<value>,gid=<value>,mode=<value>,pagesize=<value>,size=<value>,\
      min_size=<value>,nr_inodes=<value> none /mnt/huge
```

此命令在目录 /mnt/huge 上挂载一个 hubetlbfs 类型的文件系统。在 /mnt/huge 上创建的任何文件都使用大页面。

# THP

> 参考:
>
> - [Linux Kernel 文档，管理员指南 - 内存管理 - transhuge](https://www.kernel.org/doc/html/latest/admin-guide/mm/transhuge.html)


