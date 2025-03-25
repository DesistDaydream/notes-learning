---
title: EXT FileSystem
linkTitle: EXT FileSystem
weight: 20
---

# 概述

> 参考：
>
> - [公众号 - 小林 coding，一口气搞懂「文件系统」，就靠这 25 张图了](https://mp.weixin.qq.com/s/qJdoXTv_XS_4ts9YuzMNIw)
> - [骏马金龙，第4章 ext文件系统机制原理剖析](https://www.junmajinlong.com/linux/ext_filesystem/)

## block 的出现

硬盘最底层的读写 I/O 一次是一个扇区 512 字节，如果要读写大量文件，以扇区为单位肯定很慢很消耗性能，所以硬盘使用了一个称作逻辑块的概念。逻辑块是逻辑的，由磁盘驱动器负责维护和操作，它并非是像扇区一样物理划分的。一个逻辑块的大小可能包含一个或多个扇区，每个逻辑块都有唯一的地址，称为 LBA。有了逻辑块之后，磁盘控制器对数据的操作就以逻辑块为单位，一次读写一个逻辑块，磁盘控制器知道如何将逻辑块翻译成对应的扇区并读写数据。

到了 Linux 操作系统层次，通过文件系统提供了一个也称为块的读写单元，文件系统数据块的大小一般为 1024bytes (1KiB) 或 2048bytes (2KiB) 或 4096bytes (4KiB)。文件系统数据块也是逻辑概念，是文件系统层次维护的，而磁盘上的逻辑数据块是由磁盘控制器维护的，文件系统的 IO 管理器知道如何将它的数据块翻译成磁盘维护的数据块地址 LBA。

对于使用文件系统的 IO 操作来说，比如读写文件，这些 **IO 的基本单元**是**文件系统上的数据块**，一次读写一个文件系统数据块。比如需要读一个或多个块时，文件系统的 IO 管理器首先计算这些文件系统块对应在哪些磁盘数据块，也就是计算出 LBA，然后通知磁盘控制器要读取哪些块的数据，硬盘控制器将这些块翻译成扇区地址，然后从扇区中读取数据，再通过硬盘控制器将这些扇区数据重组写入到内存中去。

**Block(块)**，存放数据的最小单位，假如每个块为 4KiB，那大于 5KiB 的块就需要两个块，并且逻辑上占用了 8KiB 的空间。

**Block Group(块组)**，多个 Block 的集合

Ext 预留了一些 Inode 做特殊特性使用，如下：某些可能并非总是准确，具体的 inode 号对应什么文件可以使用 `find /-inum NUM` 查看

```bash
Ext4 的特殊 inode
Inode号    用途
0         不存在0号inode，可用于标识目录data block中已删除的文件
1         虚拟文件系统，如/proc和/sys
2         根目录         # 注意此行
3         ACL索引
4         ACL数据
5         Boot loader
6         未删除的目录
7         预留的块组描述符inode
8         日志inode
11        第一个非预留的inode，通常是 lost+found 目录
```

所以在 ext4 文件系统的 dumpe2fs 信息中，能观察到 fisrt inode 号可能为 11 也可能为 12。

## 块、块组、Inode 计算

> 参考：
>
> - 参考哪里？我也想知道真实的计算逻辑。。。

这些计算的结果通常与下列设置有关

- **DiskSize** # 磁盘空间
- **BlockSize** # 通常为 4096 Bytes
  - 可使用 mke2fs -b 手动指定
- **BlocksPerGroup** # 通常为 32768。每个块组中块的数量。
  - 可使用 mke2fs -g 手动指定
- **BytesPerInode** # 通常为 16384 Bytes。创建文件系统时，为每块 BytesPerInode 大小的空间创建一个 Inode。
  - BytesPerInode 也称为 inode_ratio，即.Inode 比率，全称应该是 Inode 分配比率，即每多少空间分配一个 Inode。
  - 可使用 mke2fs -i 手动指定
- **InodeSize** # 通常为 256 Bytes。大小是 128 的倍数，最小为 128 Bytes。

其中 BlocksPerGroup(每个块组中块的数量)、BytesPerInode(每个Inode负责的空间大小) 这种值是后续计算的基础。固定下来这些，就算分区空间自动扩容/缩容，也可以根据这种数据自动增加/删除块的数量和 Inode 的数量。

其中块大小为 4K，Inode 比率为 16K，也就是说，至少每 4 个块分配一个 Inode。当然分配的这些 inode 号只是预分配，并不真的代表会全部使用，毕竟每个文件才会分配一个 inode 号。

这些数据将会计算出：

- **BlockCount** # 块总数
- **BlockGroup** # 块组总数
- **InodeCount** # Inode 总数
- **InodePreGroup** # 每个块组中包含的 Inode 数量
- **InodeUseage** # 所有 Inode 占用的空间

假如现在有一块 35GiB 的磁盘，需要先转为 Bytes。然后根据给定的 BlockSize(块大小) 和 BlocksPerGroup(块组中块的数量)，计算出需要创建创建的**块数量**和**块组数量**。

**Block 与 BlockGroup 的计算**

- BlockCount = DiskSize / BlockSize = 37580963840 / 4096 = 9175040
- BlockGroupCount = BlockCount / BlocksPerGroup = 9175040 / 32768 = 280

**Inode 的计算**

- InodeCount = DiskSize / BytesPerInode = 37580963840 / 16384 = 2293760

> 由于之前已经知道了“每4个块分配一个Inode”，那么用“总块数/4”也是可以得到这个 Inode 总数的。

计算出的 Inode 数量将会平均分配到每个块组中

- InodePreGroup = InodeCount / BlockGroupCount = 2293760 / 280 = 8192

> 由于之前已经知道了“每4个块分配一个Inode”，那么用“每个块组中块的数量/4”也是可以得到每个块组中 Inode 的数量。

计算所有 Inode 需要占用的磁盘空间

- InodeUseage = InodeCount *InodeSize = 2293760* 256 = 587202560 Bytes = 560 MiB

也就是说，一块 35 G 的硬盘，需要拿出来至少 560 MiB 的空间来存放 Inode 数据。

**TODO:**

- **Inode 还有一个最低数量？就算执行 `mke2fs -N 1` 命令，最后也不会只有一个 Inode，而是有 4480 个 Inode，这个数是怎么来的？**
- **试了下 `mke2fs -N 1000000` 最后生成的 Inode 数为 1003520，正好是 4480 的倍数
- **如果是 `mke2fs -N 4481`，则生成的 Inode 数为 8960，也是 4480 的倍数。。。这个值到底咋来的。。。o(╯□╰)o****
- 好像是根据 Inodes per group 的值来的，这个值好像必须要是 16 的倍数，并且最低值是 16，可是这个说法的来源是在哪呢？

上述计算的结果可以通过 dumpe2fs 命令查看

```bash
~]# dumpe2fs -h ${DEVICE} | egrep -i "block|inode"
dumpe2fs 1.45.5 (07-Jan-2020)
Filesystem features:      ext_attr resize_inode dir_index filetype sparse_super large_file
Inode count:              2293760
Block count:              9175040
Reserved block count:     458752
Free blocks:              9018814
Free inodes:              2293749
First block:              0
Block size:               4096
Reserved GDT blocks:      1021
Blocks per group:         32768
Inodes per group:         8192
Inode blocks per group:   512
Reserved blocks uid:      0 (user root)
Reserved blocks gid:      0 (group root)
First inode:              11
Inode size:           256
```

### TODO: 最低的 Inode

假如我们需要最少 1 个 Inode

此时已知

- BlockCount = 9175040
- BlockGroupCount  = 280
- BlockSize = 4096 Bytes
- InodeSize = 256 Bytes
- InodeRatio = 16384 Bytes

每 4 个块给 1 个 Inode，但肯定不是这个思路，因为照着这种算法，那就是有 2293760 个。所以肯定不是每 4 个块给 1 个 Inode。

真实场景：

- 现在是每 2048 个块给 一个 Inode
- 280 个块组

```bash
~]# mke2fs -N 1 /dev/vdb
mke2fs 1.45.5 (07-Jan-2020)
/dev/vdb contains a ext2 file system
 last mounted on Sat Mar 11 16:14:14 2023
Proceed anyway? (y,N) y
Creating filesystem with 9175040 4k blocks and 4480 inodes
Filesystem UUID: acebc9ab-c53e-4f74-bd6b-443343a76bab
Superblock backups stored on blocks:
 32768, 98304, 163840, 229376, 294912, 819200, 884736, 1605632, 2654208,
 4096000, 7962624

Allocating group tables: done
Writing inode tables: done
Writing superblocks and filesystem accounting information: done

~]# dumpe2fs -h ${DEVICE} | egrep -i "block|inode"
dumpe2fs 1.45.5 (07-Jan-2020)
Filesystem features:      ext_attr resize_inode dir_index filetype sparse_super large_file
Inode count:              4480
Block count:              9175040
Reserved block count:     458752
Free blocks:              9161894
Free inodes:              4469
First block:              0
Block size:               4096
Reserved GDT blocks:      1021
Blocks per group:         32768
Inodes per group:         16
Inode blocks per group:   1
Reserved blocks uid:      0 (user root)
Reserved blocks gid:      0 (group root)
First inode:              11
Inode size:           256
```

# EXT 文件系统详解

> 参考:
>
> - [公众号 -  无聊的闪客，你管这破玩意叫文件系统](https://mp.weixin.qq.com/s/q6OjwCXSk05TvX_BIu1M0g)

**你手里有一块硬盘，大小为 1T**

![](https://mmbiz.qpic.cn/mmbiz_png/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDfZqVQibyBs0wVPoFMcGCmqwVCiaBK4j30rciagooOJK38S0Gk3Tb2udw0g/640?wx_fmt=png)

**你还有一堆文件**

![](https://mmbiz.qpic.cn/mmbiz_jpg/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDfDf52Y0RCfHGsgYDk3yP8RXAjlFgFaPwPyNPVqmiaxrCDY2y2hHDRuzQ/640?wx_fmt=jpeg)

这些文件在硬盘看来，就是一堆二进制数据而已

![](https://mmbiz.qpic.cn/mmbiz_jpg/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDfmENc0IKA0Kd8ITOfDe5D0z3wlTsyvb0iamvbzjy3icUU5uKHPb8icklIg/640?wx_fmt=jpeg)

你准备把这些文件存储在硬盘上，并在需要的时候读取出来。

要设计怎样的软件，才能更方便地在硬盘中读写这些文件呢？

首先我不想和复杂的扇区，设备驱动等细节打交道，因此我先实现了一个简单的功能，将硬盘按逻辑分成一个个的**块**，并可以以块为单位进行读写。

每个块就定义为两个物理扇区的大小，即 1024 字节，就是 1KB 啦。

硬盘太大不好分析，我们就假设你的硬盘只有 1MB，那么这块硬盘则有 1024 个块。

![](https://mmbiz.qpic.cn/mmbiz_jpg/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDf56b6Iia6zq3Yw6XeosYOV7Rds82Xn2sxuV75Iaib9Qb72yE30zNqXOicg/640?wx_fmt=jpeg)

OK，我们开始存文件啦！

准备一个文件

![](https://mmbiz.qpic.cn/mmbiz_png/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDftL60S3R7ftR7ygNTQOHRAmPayVlLj4At7uJKoAyu6icibxIXzIWAKsHg/640?wx_fmt=png)

随便选个块放进去，3 号块吧！

![](https://mmbiz.qpic.cn/mmbiz_jpg/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDfkHibk00PpYTMhgbwniamDKiarXxAU0oibicgXTjeznytfHdJZJR7a2PTl2w/640?wx_fmt=jpeg)

成功！首战告捷！

## 再存一个文件

诶？发现问题了，万一这个文件也存到了 3 号块，不是把原来的文件覆盖了么？不行，得有一个地方记录，现在可使用的块有哪些，像这样。

块 0：未使用

块 1：未使用

块 2：未使用

块 3：已使用

块 4：未使用

...

块 1023：未使用

那我们就用 0 号块，来记录所有块的使用情况吧！怎么记录呢？

**位图！**

![](https://mmbiz.qpic.cn/mmbiz_gif/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDfRruNoC8sYKTcib1ibOYFzGlLtYBphS1U3fnavQY1vasqjw4EG6IkrGfQ/640?wx_fmt=gif)

那我们给块 0 起个名字，叫**块位图**，之后这个块 0 就专门用来记录所有块的使用情况，不再用来存具体文件了。

![](https://mmbiz.qpic.cn/mmbiz_jpg/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDfN7kfnVIs7NVmZVxcBoTyK2aoa24OOxu56VMDEYXpeibswiagvGMZqB7A/640?wx_fmt=jpeg)

当我们再存入一个新文件时，只需要在块位图中找到第一个为 0 的位，就可以找到第一个还未被使用的块，将文件存入。同时，别忘了把块位图中的相应位置 1。

完美！

## 尝试读取刚刚的文件

咦？又遇到问题了，我怎么找到刚刚的文件呢？根据块号么？这也太蠢了，就像你去书店找书，店员让你提供书的编号，而不是书名，显然不合理。

因此我们给每个文件起一个名字，叫**文件名**，通过它来寻找这个文件。

那必然就要有一个地方，记录文件名与块号的对应关系，像这样。

葵花宝典.txt：3 号块

数学期末复习资料.mp4：5 号块

无聊的闪客.pdf：10 号块

...

别急，既然都要选一个地方记录文件名称了，不妨多记录一点我们关心的信息吧，比如文件大小、文件创建时间、文件权限等。

这些东西自然也要保存在硬盘上，我们选择用一个固定大小的空间，来表示这些信息，多大空间呢？128 字节吧。

为啥是 128 字节呢？我乐意。

![](https://mmbiz.qpic.cn/mmbiz_png/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDfUHOgLG93SuaRKVQx7KwvwEMkibXibUBiahk2zPGHxJEllw9yGIiaEHOZ7w/640?wx_fmt=png)

我们将这 128 字节的结构体，叫做一个 **inode**。

之后，我们每存入一个新的文件，不但要占用一个块来存放这个文件本身，还要占用一个 inode 来存放文件的这些**元信息**，并且这个 inode 的**所在块号**这个字段，就指向这个文件所在的块号。

![](https://mmbiz.qpic.cn/mmbiz_jpg/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDfXh1ic04TbyrJJ8tjpSWGDnsStYz6Diazzibg6keU2nJdODg5XAq6IZ8eQ/640?wx_fmt=jpeg)

如果一个 inode 为 128 字节，那么一个块就可以容纳 8 个 inode，我们可以将这些 inode 编上号。

![](https://mmbiz.qpic.cn/mmbiz_jpg/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDf4yV2MR2qXibFaibLRaT23S0FicZyuFAbUtCQnylRic8KdQL6zS04KvzM2Q/640?wx_fmt=jpeg)

如果你觉得 inode 数不够，也可以用两个或者多个块来存放 inode 信息，但这样用于存放数据的块就少了，这就看你自己的平衡了。

![](https://mmbiz.qpic.cn/mmbiz_jpg/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDffrPoFxjuBZdpk3FFcrBricxWA8Iur8vSianfxgOFbVSxA1qjkMF9AugQ/640?wx_fmt=jpeg)

同样，和块位图管理块的使用情况一样，我们也需要一个 **inode 位图**，来管理 inode 的使用情况。我们就把 inode 位图，放在 1 号块吧！

同时，我们把 inode 信息，放在 2 号块，一共存 8 条 inode，这样我们的 2 号块就叫做 **inode 表**。

现在，我们的文件系统结构，变成了下面这个样子。

![](https://mmbiz.qpic.cn/mmbiz_jpg/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDf7uy19fsibjjvkCBD82o2LzYoiaXMOUCoH5tswW0HVTAIxmdsbyH0oVDQ/640?wx_fmt=jpeg)

注意：块位图是管理可用的块，每一位代表一个块的使用与否。inode 位图管理的是一条一条的 inode，并不是 inode 所占用的块，比如上图中有 8 条 inode，则 inode 位图中就有 8 位是管理他们的使用与否。

## 多个块

现在，我们的文件很小，一个块就能容下。

但如果需要两个块、三个块、四个块呢？

很简单，我们只需要采用**连续存储法**，而 inode 则只记录文件的第一个块，以及后面还需要多少块，即可。

这种办法的缺点就是：容易留下大大小小的**空洞**，新的文件到来以后，难以找到合适的空白块，空间会被浪费。

![](https://mmbiz.qpic.cn/mmbiz_png/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDfdjiaibCH3Amv5m7eslaImM2Z1HgnobqPdicCdczx9yd01ZkkykHxiaB4ug/640?wx_fmt=png)

看来这种方式不行，那怎么办呢？

既然在 inode 中记录了文件所在的块号，为什么不扩展一下，多记录几块呢？

![](https://mmbiz.qpic.cn/mmbiz_png/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDfUHOgLG93SuaRKVQx7KwvwEMkibXibUBiahk2zPGHxJEllw9yGIiaEHOZ7w/640?wx_fmt=png)

原来在 inode 中只记录了一个块号，现在扩展一下，记录 8 个块号！而且这些块**不需要连续**。

![](https://mmbiz.qpic.cn/mmbiz_png/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDfkJRKnxe1tV7DsOYnibcqnx06huWZHgJZS8qG5kKG3UPOfAlvGF5oL1A/640?wx_fmt=png)

嗯，这是个可行的办法！

但是这也仅仅能表示 8 个块，能记录的最大文件是 8K（记住，一个块是 1K）, 现在的文件轻松就超过这个限制了，这怎么办？

很简单，我们可以让其中一个块，作为**间接索引**。

![](https://mmbiz.qpic.cn/mmbiz_png/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDf78H2eJA2Hsib2kDzzibPdaM37BII5aJOkIbemBSWlfKL9VpaLmLemYkA/640?wx_fmt=png)

这样瞬间就有 263 个块（多了 256 -1 个块）可用了，这种索引叫**一级间接索引**。

如果还嫌不够，就再弄一个块做一级间接索引，或者做二级间接索引（二级间接索引则可以多出 256 * 256 - 1 个块）。

我们的文件系统，暂且先只弄一个一级间接索引。硬盘一共才 1024 个块，一个文件 263 个块够大了。再大了不允许，就这么任性，爱用不用。

好了，现在我们已经可以保存很大的文件了，并且可以通过文件名和文件大小，将它们准确读取出来啦！

## 元数据记录

但我们得精益求精，我们再想想看这个文件系统有什么毛病。

比如，inode 数量不够时，我们是怎么得知的呢？是不是需要在 inode 位图中找，找不到了才知道不够用了？

同样，对于块数量不够时，也是如此。

要是有个全局的地方，来记录这一切，就好了，也方便随时调整，比如这样

inode 数量

空闲 inode 数量

块数量

空闲块数量

那我们就再占用一个块来存储这些数据吧！由于他们看起来像是站在上帝视角来描述这个文件系统的，所以我们把它放在最开始的块上，并把它叫做**超级块**，现在的布局如下。

![](https://mmbiz.qpic.cn/mmbiz_jpg/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDfc7McjdEjQfcUbDDxOzVw6ZDqLLrJDAsY8IgdXN5lmsodZ6aHCxpbjg/640?wx_fmt=jpeg)

我们继续精益求精。

现在，**块位图**、**inode 位图**、**inode 表**，都是是固定地占据这块 1、块 2、块 3 这三个位置。

假如之后 inode 的数量很多，使得 inode 表或者 inode 位图需要占据多个块，怎么办？

或者，块的数量增多（硬盘本身大了，或者每个块变小了），块位图需要占据多个块，怎么办？

程序是死的，你不告诉它哪个块表示什么，它可不会自己猜。

很简单，与超级块记录信息一样，这些信息也选择一个块来记录，就不怕了。那我们就选择紧跟在超级块后面的 1 号块来记录这些信息吧，并把它称之为**块描述符**。

![](https://mmbiz.qpic.cn/mmbiz_jpg/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDffia4t0XszAj6icMpTAiaJY4d5cgN0Ka0DnULFZZZlCEdK3UmanFbqcSEA/640?wx_fmt=jpeg)

当然，这些所在块号只是记录起始块号，块位图、inode 位图、inode 表分别都可以占用多个块。

好了，大功告成！

## 我们再尝试存入一批文件

- 葵花宝典.txt
- 数学期末复习资料.mp4
- 赘婿1.mp4
- 赘婿2.mp4
- 赘婿3.mp4
- 赘婿4.mp4
- 无聊的闪客.pdf

诶？这看着好不爽，所有的文件都是平铺开的，能不能拥有**层级关系**呢？比如这样

- 葵花宝典.txt
- 数学期末复习资料.mp4
- 赘婿
- 赘婿1.mp4
- 赘婿2.mp4
- 赘婿3.mp4
- 赘婿4.mp4
- 无聊的闪客.pdf

我们将葵花宝典.txt 这种称为**普通文件**，将赘婿这种称为**目录文件**，如果要访问赘婿1.mp4，那全文件名要写成：赘婿/赘婿1.mp4。

如何做到这一点呢？那我们又得把 inode 结构拿出来说事了。

![](https://mmbiz.qpic.cn/mmbiz_png/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDfUHOgLG93SuaRKVQx7KwvwEMkibXibUBiahk2zPGHxJEllw9yGIiaEHOZ7w/640?wx_fmt=png)

此时需要一个属性来区分这个文件是普通文件，还是目录文件。

缺什么就补什么嘛，我们已经很熟悉了，专门加一个 4 字节，来表示**文件类型**。

![](https://mmbiz.qpic.cn/mmbiz_png/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDf9C2jf6WySMyJTvESoChM9aOCWDibOys76I4hufo45zu0WPWWMnea0xw/640?wx_fmt=png)

如果是**普通文件**，则这个 inode 所指向的数据块仍然和之前一样，就是文件本身原封不动的内容。

但如果是**目录文件**，则这个 inode 所指向的数据块，就需要重新规划了。

这个数据块里应该是什么样子呢？可以是一个一个指向不同 inode 的紧挨着的结构体，比如这样。

![](https://mmbiz.qpic.cn/mmbiz_png/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDf4ia6OOsneVH1xWqnibiboqPCbDvZetEmTJ15oRATib4NicpqHibhPPp83BMw/640?wx_fmt=png)

这样先通过 **赘婿**这个目录文件，找到所在的数据块。再根据这个数据块里的一个个带有 **inode** 信息的结构体，找到这个目录下的所有文件。

完美！

## 目录

不过这样的话，你想想看，如果想要查看一下赘婿**这个目录下的所有文件**（比如 ll 命令），将文件名和文件类型都展示出来，怎么办呢？

就需要把一个个结构体指向的 inode 从 inode 表中取出，再把文件名和文件类型取出，这很是浪费时间。

而让用户看到一个目录下的所有文件，又是一个极其常见的操作。

所以，不如把文件名和文件类型这种常见的信息，放在数据块中的结构体里吧。

![](https://mmbiz.qpic.cn/mmbiz_png/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDfaCmiaOAibU7q8kxf8CJdJRYYdqGdmnGAvb2pN5IhKENAON2GicIjAAAww/640?wx_fmt=png)

同时，inode 结构中的文件名，好像就没啥用了，这种变长的东西放在这种定长的结构中本身就很讨厌，早就想给它去掉了。而且还能给其他信息省下空间，比如文件所在块的数组，就能再多几个了。

太好了，去掉它！

![](https://mmbiz.qpic.cn/mmbiz_png/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDfJSNp4Dxibd7IFba5DoicWdgkHiba1khAdAibwU1vmpDjN8ia3DS1ibwMklew/640?wx_fmt=png)

OK，大功告成，现在我们就可以给文件分门别类放进不同目录下了，还可以在目录下创建目录，无限套娃！

现在的文件系统，已经比较完善了，只是还有一点不太爽。

## 最后

我们访问到一个目录下，可以很舒服地看到目录里的文件，然后再根据名称访问这个目录下的文件或者目录，整个过程都是一个套路。

但是，最上层的目录下的所有文件，即**根目录**，现在仍然需要通过遍历所有的 inode 来获得，能不能和上面的套路统一呢？

答案非常简单，我们规定，**inode 表中的 0 号 inode，就表示根目录**，一切的访问，就从这个根目录开始！

![](https://mmbiz.qpic.cn/mmbiz_png/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDf8t0qNVGLR6cKLv4SR0SGywsX1b5lPcu3AcGzUNjnCQ7V7QA7YE5IWg/640?wx_fmt=png)

好了，这回没有然后了！

我们最后来欣赏下我们的文件系统架构。

![](https://mmbiz.qpic.cn/mmbiz_png/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDfnk5ApwlxjBJSF2rSoxyhVFDcmktSrrah5Mj8iakhm4raOciaT4PLqQRg/640?wx_fmt=png)

你是不是觉得这没啥了不起的。

**但这个破玩意，它就叫文件系统**

## 后记

这个文件系统，和 linux 上的经典文件系统 **ext2** 基本相同。

下面是我画的 ext2 文件系统的结构（字段部分只画了核心字段）

![](https://mmbiz.qpic.cn/mmbiz_png/GLeh42uInXRPpGdb1rLVOuU9fgMuLiaDfP9ia54kT3wLekE3Id4bAib4EYTHVWP7p0PX1VQun1789BYzJNgTq2uUw/640?wx_fmt=png)

如果你想了解更多的细节，可以参考官方说明文档: https://www.nongnu.org/ext2-doc/ext2.pdf

你也可以用 linux 的 mke2fs 命令生成一个 ext2 文件系统的磁盘镜像，然后一个字节一个字节地对照这官方说明文档拆解，这种方式其实是最直接的。
