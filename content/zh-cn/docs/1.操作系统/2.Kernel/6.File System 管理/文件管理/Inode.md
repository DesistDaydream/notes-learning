---
title: Inode
---

# 概述

> 参考：
>
> - [Wiki,Inode](https://en.wikipedia.org/wiki/Inode)
> - [知乎，本地文件系统小计（二）：inode](https://zhuanlan.zhihu.com/p/78724124)

**Index node(索引节点，简称 inode)** 是 Unix 风格的文件系统中的一种数据结构。每个索引节点保存了文件系统中的一个文件系统对象(i.e.文件、目录等)的元信息数据，但不包括数据内容或者文件名。

注：数据的内容存放在硬盘的一个 block(区域块) 中，通过 inode(索引节点) 来访问 block，每个索引节点都会命名一个文件名。文件的索引节点通过 `ls -i` 命令查看

```bash
~]# ls -i /
     12 bin         1 dev  6029313 home       14 lib32       16 libx32      2097153 media  1572865 opt   1048577 root       17 sbin  1966081 srv             1 sys  4325377 usr
      2 boot  1835009 etc       13 lib        15 lib64       11 lost+found  4063233 mnt          1 proc        1 run   4980737 snap       18 swap.img  3801089 tmp  1179649 var
```

所以，linux 里的所有文件，都相当于一个硬链接，链接到 inode 号上，展现在屏幕上的只是该文件内容的文件名。就算几个文件名字不一样的文件只要节点号相同，那么这几个文件的内容是就是相同的。想查看文件内容，就要找到该文件名对应的 inode 然后通过 inode 找到 block，找到 block 就能看到其中的内容了。

实际例子：我家里的内容有什么就是 block，我家的门牌号就是 inode，我家叫什么名字，比如叫“DesistDaydream 的家“这几个字就是文件名，通过文件名找到门牌号，找到了门牌号才能开门看到我家中的内容。（能不能开门就是权限决定的了）

## Inode 的计算

不同文件系统，有不同的计算方式，详情可以参考各类文件系统中的章节

- [EXT FileSystem](/docs/1.操作系统/2.Kernel/6.File%20System%20管理/磁盘文件系统/EXT%20FileSystem.md#块、块组、Inode%20计算)
- [XFS](/docs/1.操作系统/2.Kernel/6.File%20System%20管理/磁盘文件系统/XFS.md)

总得来说，所以我们可以<font color="#ff0000">通过降低 BytesPerInode 的值以提高 Inode 数量</font>

## 索引区和数据区

文件系统怎样从文件名索引到存放文件内容的 Block？

解答这个问题，先了解下 inode。

linux 文件系统在磁盘分区格式化后可以简单认为会划分**两个区域**

- 一个是**索引区**（inode 区），用来存放文件的 inode（inode 包含文件的各种属性，统称为元数据），索引区存放的 inode 个数是有上限的，在格式化的时候会分配好
- 另一个是**数据区**，存放文件数据，即文件里面的真实内容。

如果文件系统太大，将所有的 inode 与 block 放在一起很难管理，因此 Ext2 文件系统在格式化的时候基本上是区分为多个 **block group(块组)**，每个块组都有独立的 inode/block/super block 系统。

inode 的大小和 block 的大小可以使用命令来看，inode 的大小一般为 256 Bytes，block 的大小一般为 4096 Bytes。

```bash
export DEVICE="/dev/vdb"
dumpe2fs ${DEVICE} | egrep -i "Block size|Inode size"
```

inode 会有一个总量，如果 inode 数耗尽了（小文件很多），磁盘就使用不了了。inode 的总量和使用量可以用 `df -i` 命令查看

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/xz4nq5/1627368676221-0f8a3d6d-1159-4806-85d6-20fa8f3c9801.png)

/dev/sdb 中可以看到 Inodes 的具体总数，使用情况和可用个数

inode 总数有一个简单的计算公式，如下所示。假设 bytes-per-inode 为 1024，一个 1GB 的磁盘，inode 的总数可能有 1024 \* 1024 个，如果 inode 的大小为 256 字节，那么索引区的大小会达到 256M。

`inode_count = 磁盘大小 / bytes-per-inode`

其中 bytes-per-inode 在格式化的时候可以指定，例如

`mkfs.ext4 -i bytes-per-inode /dev/sda2`

不管是什么文件类型（包括目录），都会给它分配一个 inode，每个 inode 都有一个唯一的编号，文件系统就用 inode 号码来识别不同的文件。inode 结构体包含文件的元信息，简单来说有以下内容：文件大小、所属组、读写执行权限、时间戳、链接数（硬链接），**文件数据 block 的地址**。

inode 里面存有文件数据 block 的地址，所以要想获取文件的数据，必须要得到对应 inode。但是 inode 里面并没有保存文件名，文件名和 inode 号是存在于上一级目录的内容里面，所以得获取到目录的内容。目录的内容也至少占用一个 block，可以做一个实验来证明。分配一个 100M 的磁盘，然后格式化成 ext4 文件系统，block size 设置为 4K，并挂载。

`mkfs.ext4 -b 4096 /dev/sdk mount /dev/sdk /mnt/test4`

在往 /mnt/test3 目录下创建 20480 个目录前后，分别统计下 df 命令输出的结果如下

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/xz4nq5/1627368676228-626187af-52ac-4225-b43a-1f0b72b5b0c1.png)

创建 20480 个目录前后，blocks 已用的个数对比

可以看到已用的 blocks 个数 82444/4 = 20611 个，这里面有 20480 个是存的是 20480 个目录的内容。

# 源码解析

## inode 中判断是文件还是目录

在 linux 中一切皆文件，目录也是文件。虽然都是文件，但类型不同，可以通过 inode 结构体中的 `umode*t i_mode` 成员来判断是目录文件，还是普通文件。

`typedef unsigned short umode_t`

umode_t 是 16 位的 2 进制数，保存的就是文件类型及用户权限信息

- 第 0-3 位---文件类型位
- 第 4 位 -- suid 位
- 第 5 位 -- sgid 位
- 第 6 位 -- sticky 位
- 第 7-9 位 -- 文件所属主权限位
- 第 10-12 位 -- 文件所属组权限位
- 第 13-15 位 -- 其他用户权限位

1，文件类型位

用来判断文件类型

```c
#define         S_IFMT  0170000 /* type of file ，文件类型掩码*/
#define         S_IFREG 0100000 /* regular 普通文件*/
#define         S_IFBLK 0060000 /* block special 块设备文件*/
#define         S_IFDIR 0040000 /* directory 目录文件*/
#define         S_IFCHR 0020000 /* character special 字符设备文件*/
......
```

S_IFMT 是文件类型掩码，其值是 0170000，转换成二进制就是 1111 0000 0000 0000，

S_IFMT 就是用来取 mode 的 0--3 位

所以判断文件还是目录的方法就出来了，如下。ceph 文件系统的用户态客户端中 struc Inode 结构体中 uint32_t mode 成员来判断是目录还是文件：is_dir()函数判断是否是目录

```c
#define         S_ISDIR(m)      (((m) & S_IFMT) == S_IFDIR)    // 判断是否是目录文件
#define         S_ISREG(m)      (((m) & S_IFMT) == S_IFREG)    // 判断是否是普通文件
```

_2，SUID，SGID，sticky 位_

SUID 是 Set User ID, SGID 是 Set Group ID 的意思。具体怎么用，暂时不研究

_3，文件权限_

文件权限分为属主权限、属主组权限和其他用户权限，即我们所知的 rwxrwxrwx（777）之类

## 目录项：ext4_dir_entry

普通文件的 inode 指向的数据块存的是它自己的数据的，而目录的 inode 指向的数据块存的是位于该目录下的目录项的，即目录下的**子目录和文件的目录项**，**这里注意此目录项并不是 dentry**，dentry 是 vfs 层的，只存在于**内存**中，而各个实际的文件系统都有自己的目录项，这都是存在**磁盘**上的，比如 ext4 文件系统的目录项是**ext4_dir_entry**，里面有 inode 号，文件名。

```c
struct ext4_dir_entry {
 __le32 inode;   /* Inode number */
 __le16 rec_len;  /* Directory entry length */
 __le16 name_len;  /* Name length */
 char name[EXT4_NAME_LEN]; /* File name */
};
```

根据 inode 号，再找到存于索引区的 inode 数据。

inode 里面不会存在文件名，文件名存在上一级目录的 block 里面，删除文件实际上只删除了上一级目录中该文件名的记录，即对应的 `ext4*dir_entry` 目录项。如果文件在使用的情况下（被服务占用），删除文件，仅仅是删除了文件名，空间是没有释放掉（之前遇到过，ceph-mon 产生了很大的日志文件，将日志文件删掉后，发现空间仍然没有释放），可以通过命令

lsof | grep deleted

可以看到被占用的文件.

这里注意：vfs 和各个实际文件系统都有 inode，vfs 的是 struct inode，这部分数据只存在于内存中，而具体存于磁盘上的是具体文件系统的 inode，比如 ext4 文件系统中是 `struct ext4*inode`，ceph 文件系统中是 `struct ceph_inode_info`。

## 根据文件名索引到文件内容

表面上，用户通过文件名，打开文件。实际上，系统内部这个过程分成三步：首先，系统找到这个文件名对应的 inode 号码；其次，通过 inode 号码，获取 inode 信息；最后，根据 inode 信息，找到文件数据所在的 block，读出数据。

现在举一个具体的例子，来说明文件是怎么读取到的，比如读取/home/bzw/test 里的内容，目录结构如下图

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/xz4nq5/1627368676172-1255beef-a616-4a70-8edb-b34cf4038201.png)
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/xz4nq5/1627368676862-227a6bd6-c0a7-4a8e-8dae-bcdf9b042f6d.png)
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/xz4nq5/1627368676840-ce833059-9860-4974-a7d7-57704bac2ed6.png)

假设文件系统的的简单分区如下

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/xz4nq5/1627368677025-108015bd-75b2-4abe-98b1-bf8d6c967f91.png)

- 获取 home 对应的 inode 号：先找根目录'/'的 inode（不考虑缓存）：根目录的 inode 号可以从 super_block 中获取，ext4 文件系统的根目录 inode 号为 2（xfs 文件系统根目录 inode 号是 64，ceph 文件系统根目录 inode 号是 1），所以在索引区读取 inode 号为 2 存的 inode 内容。假如 inode 中存的 block 地址是\_1000\_，那么去数据区读取地址为\_1000\_的 block 存的内容，内容如下图所示。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/xz4nq5/1627368677348-ef597c30-1973-4696-baae-ebe253a1bb64.png)

地址为 1000 的 block 存的内容

地址为 1000 的 block 里面存了 20 个目录项（struct ext4\*dir_entry）,可以找到\*\*\*目录 home 对应的 inode 号为 100\_\*\*。

- 获取 bzw 对应的 inode 号：上一步获取到了目录 home 的 inode 号，在索引区读取 inode 号为\_100\_存的 inode 内容。假如 inode 中存的 block 地址为\_2000\_，那么去读地址为\_2000\_的 block 存的内容，内容如下图所示。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/xz4nq5/1627368677456-031fdb40-9076-4165-a80e-ebf6d422b349.png)

地址为 2000 的 block 存的内容

地址为 2000 的 block 里面存了 3 个目录项（struct ext4\*dir_entry）,可以找到\*\*\*目录 bzw 对应的 inode 号为 200\_\*\*。

- 获取 test 对应的 inode 号：上一步获取到了目录 bzw 的 inode 号，在索引区读取 inode 号为\_200\_存的 inode 内容。假如 inode 中存的 block 地址为\_3000\_，那么去读地址为\_3000\_的 block 存的内容，内容如下图所示。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/xz4nq5/1627368677627-211f81db-d5cb-4eb7-95bc-b682d43bbc54.png)

地址为 3000 的 block 存的内容

地址为 3000 的 block 里面存了 2 个目录项（struct ext4\*dir_entry）,可以找到\*\*\*文件 test 对应的 inode 号为 300\_\*\*。

- 获取 test 对应的内容：上一步获取到了文件 test 的 inode 号，在索引区读取 inode 号为\_300\_存的 inode 内容。假如 inode 中存的 block 地址为\_4000\_，那么去读地址为\_4000\_的 block 存的内容。这个时候就完成了操作。

这里注意如果 test 内容很大，那么在 inode 里面存的 block 地址就不止一个了。

可以以 ext4 中的 struct ext4_inode 为例

```c
struct ext4_inode {
 __le16 i_mode;  /* File mode */
 __le16 i_uid;  /* Low 16 bits of Owner Uid */
 __le32 i_size_lo; /* Size in bytes */
 __le32 i_atime; /* Access time */
 __le32 i_ctime; /* Inode Change time */
 __le32 i_mtime; /* Modification time */
 __le32 i_dtime; /* Deletion Time */
 __le16 i_gid;  /* Low 16 bits of Group Id */
 __le16 i_links_count; /* Links count */
 __le32 i_blocks_lo; /* Blocks count */
 __le32 i_flags; /* File flags */
        ......
 __le32 i_block[EXT4_N_BLOCKS];/* Pointers to blocks */ 这里面存的就是block的地址
 __le32 i_generation; /* File version (for NFS) */
 __le32 i_file_acl_lo; /* File ACL */
 __le32 i_size_high;
 __le32 i_obso_faddr; /* Obsoleted fragment address */
 ......
 __le16 i_extra_isize;
 __le16 i_checksum_hi; /* crc32c(uuid+inum+inode) BE */
 __le32  i_ctime_extra;  /* extra Change time      (nsec << 2 | epoch) */
 __le32  i_mtime_extra;  /* extra Modification time(nsec << 2 | epoch) */
 __le32  i_atime_extra;  /* extra Access time      (nsec << 2 | epoch) */
 __le32  i_crtime;       /* File Creation time */
 __le32  i_crtime_extra; /* extra FileCreationtime (nsec << 2 | epoch) */
 __le32  i_version_hi; /* high 32 bits for 64-bit version */
 __le32 i_projid; /* Project ID */
};
```

## inode_hashtable

文件系统中的位于内存中的所有 inode 存放在一个名为 inode_hashtable 的全局哈希表中（如果 inode 还在磁盘，尚未读入内存中，则不会加入到全局哈希表中）。另一方面，所有的 inode 还存放在超级块中的 s_inode 链表中。

`static struct hlist_head *inode_hashtable __read_mostly;`

在上面进行文件索引时，并没有讲根目录 inode 号为 2 是怎么获取的，这将下一章中讲解。

# 分类

#文件系统 #文件