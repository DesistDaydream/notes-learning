---
title: Memory 的缓存
linkTitle: Memory 的缓存
weight: 3
---

# 概述

> 参考：
>
> -

linux 中每个程序启动之后都会占用内存，一般情况下是不会把内存全部占满。那么空闲的这部分内存用来干什么呢？~

Linux 会充分利用这些空闲的内存，设计思想是：内存空闲还不如拿来多缓存一些数据，等下次程序再次访问这些数据速度就快了，而如果程序要使用内存而系统中内存又不足时，这时不是使用交换分区，而是快速回收部分 cached，将它们留给用户程序使用。

比如说：当使用该程序时，就会从硬盘中读取该程序的内容，这一部分读取的内容会加载到内存中(caceh 中)，因为内存比硬盘的读写速度更快，所以当下次再使用该程序需要读取该程序的内容时，直接从内存就能读取了。而 buff 中的数据一般是程序运行中产生的(比如玩游戏的存档)，当程序结束之前，需要把再 buff 中的数据写入到硬盘中以便永久保存(要不再运行这游戏，不就没存档了么~哈哈)。

内存管理做了很多精心的设计，除了对 dentry 进行缓存（用于 VFS，加速文件路径名到 inode 的转换），还采取了两种主要 Cache 方式：Buffer Cache 和 Page Cache，目的就是为了提升磁盘 IO 的性能。从低速的块设备(e.g.硬盘)上读取数据会暂时保存在内存中，即使数据在当时已经不再需要了，但在应用程序下一次访问该数据时，它可以从内存中直接读取，从而绕开低速的块设备，从而提高系统的整体性能。

因此，可以看出，buffers/cached 真是百益而无一害，真正的坏处可能让用户产生一种错觉——Linux 耗内存！其实不然，Linux 并没有吃掉你的内存，只要还未使用到交换分区，你的内存所剩无几时，你应该感到庆幸，因为 Linux 缓存了大量的数据，也许某一次应用程序读取数据时，就会用到这些已经缓存了的数据！

## Buffer 与 Cache

在 [Memory 管理工具](/docs/1.操作系统/Linux%20管理/Linux%20系统管理工具/Memory%20管理工具.md) 的 free 命令中，可以看到的 buff 与 cache 是从虚拟文件系统 /proc/meminfo 中获取的数据

- **buff** # 内核缓冲器用到的内存，对应的是 /proc/meminfo 中 Buffers 的值
  - Buffers 是对原始磁盘块的临时存储，也就是来**缓存从磁盘块读写的数据**，通常不会特别大(20MB 左右)。
    - 这样，内核就可以把分散的写集中起来，统一优化磁盘的写入，比如可以把多次小的写合并成单次大的写，等等。
- **cache** # 内核 PageCache(页缓存) 和 Slab 用到的内存，对应的是 /proc/meminfo 中的 Cached 与 SReclaimable 之和
  - Cached 是从磁盘读写文件时的页缓存，也就是用来**缓存从文件读写的数据**。
    - 读取过一个文件后会缓存，下次访问过这些文件数据时，就可以直接从内存中快速获取，而不需要再次访问缓慢的磁盘。
    - 或者在向文件写入数据时，先将数据写入到内存中，然后系统后台再逐步将数据从内存写入到磁盘中。
  - SReclaimable 是 Slab 的一部分。Slab 包括两部分，其中的可回收部分，用 SReclaimable 记录；而不可回收部分，用 SUnreclaim 记录。

### 实验示例

写文件时会用到 Cache 缓存数据，写磁盘时则会用到 Buffer 缓存数据

读文件时数据会缓存到 Cache 中，而读磁盘时数据会缓存到 Buffer 中。

**简单来说，Buffer 是对磁盘数据的缓存，而 Cache 是文件数据的缓存，他们既会用在读请求中，也会用在写请求中。**

上述结论的实验是通过 dd 命令分别读写磁盘和文件(参考 dd 命令用法)，并使用 vmstat 命令观察缓存变化情况得出的结论。

**通过缓存提高读取速度的实验示例**

创建一个大文件，读取(以便让文件缓存到内存中)该文件查看消耗的时间，再次读取(已经缓存到内存中)该文件查看消耗的时间

1.首先生成一个 1G 的大文件

```bash
~]# dd if=/dev/zero of=bigfile bs=1M count=1000
1000+0 records in
1000+0 records out
1048576000 bytes (1.0 GB) copied, 0.687036 s, 1.5 GB/s
```

2.清空缓存

```bash
~]# echo 1 > /proc/sys/vm/drop_caches
~]# free -wh
               total        used        free      shared     buffers       cache   available
Mem:           3.8Gi       300Mi       3.6Gi       812Ki        48Ki        53Mi       3.5Gi
Swap:          2.4Gi          0B       2.4Gi
```

3.读入这个文件，测试消耗的时间

```bash
~]# time cat bigfile > /dev/null
real 0m0.614s
user 0m0.011s
sys 0m0.412s
bonding]# free -wh
               total        used        free      shared     buffers       cache   available
Mem:           3.8Gi       375Mi       2.6Gi       812Ki        56Ki       1.0Gi       3.5Gi
Swap:          2.4Gi          0B       2.4Gi
```

4.再次读入该文件，测试消耗的时间

```
~]# time cat bigfile > /dev/null
real 0m0.217s
user 0m0.011s
sys 0m0.206s
```

从上面看出，第一次读这个 1G 的文件大约耗时 0.6s，而第二次再次读的时候，只耗时 0.2s，提升了 3 倍。如果是 centos6 的话，会足足提升 60 倍！

## TODO: Cached

Cached 的来源都是什么？什么情况下 available 的值会小于 cache 的值？下面这个为什么会这样？

```bash
~]# free -hw && cat /proc/meminfo | grep -i -E 'MemTotal|MemFree|MemAvailable|Buffers|Cached|Slab|SReclaimable|SUnreclaim'
               total        used        free      shared     buffers       cache   available
Mem:           377Gi       366Gi       937Mi       2.5Gi       1.1Mi        10Gi       6.3Gi
Swap:          4.0Gi       4.0Gi          0B
MemTotal:       395528228 kB
MemFree:          959756 kB
MemAvailable:    6614588 kB
Buffers:            1104 kB
Cached:          7516320 kB
SwapCached:       294916 kB
Slab:            4020388 kB
SReclaimable:    3042324 kB
SUnreclaim:       978064 kB
```

## Swap

在物理磁盘上划分的一块空间，用于当做内存使用，称为 swap，一般情况用不到。但是当连续内存空间不足的时候，需要使用 Swap 对内存中的数据进行换入和换出，以腾出来一个连续的大空间供新启动的程序使用。

Swap 分区（也称交换分区）是硬盘上的一个区域，被指定为操作系统可以临时存储数据的地方，这些数据不能再保存在 RAM 中。 基本上，这使您能够增加服务器在工作“内存”中保留的信息量，但有一些注意事项，主要是当 RAM 中没有足够的空间容纳正在使用的应用程序数据时，将使用硬盘驱动器上的交换空间。

写入磁盘的信息将比保存在 RAM 中的信息慢得多，但是操作系统更愿意将应用程序数据保存在内存中，并使用交换旧数据。 总的来说，当系统的 RAM 耗尽时，将交换空间作为回落空间可能是一个很好的安全网，可防止非 SSD 存储系统出现内存不足的情况。

调整 swap 的内核参数

vm.swappiness 这个内核参数可以用来调整系统使用 swap 的时机。默认值为 60，i.e.当内存中空闲空间低于 60%的时候，就会开始使用 swap 空间(也就是说系统使用了 40%的内存之后，就开始使用 swap) ，一般可以将该值设置为 10 来进行优化。

```bash
cat >>/etc/sysctl.conf <<END
vm.swappiness = 10
END
```

swappiness 参数配置您的系统将数据从 RAM 交换到交换空间的频率, 值介于 0 和 100 之间，表示百分比。如果 swappiness 值接近 0，内核将不会将数据交换到磁盘，除非绝对必要。要记住一点，与 swap 文件的交互是“昂贵的”，因为与 swap 交互花费的时间比与 RAM 的交互更长，并且会导致性能的显著下降。系统更少依赖 swap 分区通常会使你的系统更快。swappiness 接近 100 的值将尝试将更多的数据放入交换中，以保持更多的 RAM 空间。根据您的应用程序的内存配置文件或您使用的服务器，这可能会在某些情况下更好。

### 使用 Swap

- https://www.myfreax.com/how-to-add-swap-space-on-ubuntu-22-04/
- https://www.yangbolin.com/?id=296

```bash
export SWAP_FILE="/swapfile.img"
dd if=/dev/zero of=${SWAP_FILE} bs=1024 count=2000000
chmod 600 ${SWAP_FILE}
mkswap -f ${SWAP_FILE}

# 激活 Swap
swapon ${SWAP_FILE}
```

可以在 /etc/fstab 文件中添加配置以实现开启自动激活 swap

```
/swapfile.img none swap sw 0 0
```

# 缓存的清理

可以通过写入 /proc/sys/vm/drop_caches 这个内核参数来让内核从内存中清理 cache、dentries 和 inode。

> [!Attention] 由于写入该文件是一个非破坏性操作，而且脏对象是不可释放的，所以应该首先运行 sync 命令。
> 释放内存前先使用 sync 命令做同步，以确保文件系统的完整性，将所有未写的系统缓冲区写到磁盘中，包含已修改的 inode、已延迟的块 I/O 和读写映射文件。否则在释放缓存的过程中，可能会丢失未保存的文件。

## 仅清理 PageCache(页面缓存)

sync; echo 1 > /proc/sys/vm/drop_caches

## 清理 Dentries 和 Inode

sync; echo 2 > /proc/sys/vm/drop_caches

## 清理 PageCache、Dentries 和 Inode

sync; echo 3 > /proc/sys/vm/drop_caches

> [!Note]
> 如果必须清除磁盘高速缓存，第一个命令在企业和生产环境中是最安全，`...echo 1> ...` 只会清除页面缓存。 在生产环境中不建议使用上面的第三个选项 `...echo 3 > ...` ，除非明确自己在做什么，因为它会清除缓存页，目录项和 inodes。

## swap 清理

swapoff -a && swapon -a

注意：这样清理有个前提条件，空闲的内存必须比已经使用的 swap 空间大
