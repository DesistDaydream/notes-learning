---
title: revision概念
---

#

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dv1a6u/1616136169287-784743e4-15d5-4b53-81a3-4ba3477a27d4.jpeg)

每次 key 的 value 改变，version 都会+1

revision 概念

Etcd 存储数据时，并不是像其他的 KV 存储那样，存放数据的键做为 key，而是以数据的 revision 做为 key，键值做为数据来存放。如何理解 revision 这个概念，以下面的例子来说明。

比如通过批量接口两次更新两对键值，第一次写入数据时，写入和，在 Etcd 这边的存储看来，存放的数据就是这样的：

    revision={1,0}, key=key1, value=value1
      revision={1,1}, key=key2, value=value2

而在第二次更新写入数据和后，存储中又记录（注意不是覆盖前面的数据）了以下数据：

    revision={2,0}, key=key1, value=update1
    revision={2,1}, key=key2, value=update2

其中 revision 有两部分组成，第一部分成为 main revision，每次事务递增 1；第二部分称为 sub revision，一个事务内的一次操作递增 1。 两者结合，就能保证每次 key 唯一而且是递增的。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dv1a6u/1616136169539-c45ecbb0-f59a-48f4-b901-9cb6b923f169.jpeg)

但是，就客户端看来，每次操作的时候是根据 Key 来进行操作的，所以这里就需要一个 Key 映射到当前 revision 的操作了，为了做到这个映射关系，Etcd 引入了一个内存中的 Btree 索引，整个操作过程如下面的流程所示。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dv1a6u/1616136169298-3f7e3fa8-e446-400e-8830-f3fdaf07e6b5.jpeg)

查询时，先通过内存中的 btree 索引来查询该 key 对应的 keyIndex 结构体，然后再根据这个结构体才能去 boltdb 中查询真实的数据返回。

所以，下面先展开讨论这个 keyIndex 结构体和 btree 索引。

keyIndex 结构

keyIndex 结构体有以下成员：

- key：存储数据真实的键。

- modified：最后一次修改该键对应的 revision。

- generations：generation 数组。

如何理解 generation 结构呢，可以认为每个 generation 对应一个数据从创建到删除的过程。每次删除 key 的操作，都会导致一个 generation 最后添加一个 tombstone 记录，然后创建一个新的空 generation 记录添加到 generations 数组中。

generation 结构体存放以下数据：

- ver：当前 generation 中存放了多少次修改，其实就是 revs 数组的大小-1（因为需要去掉 tombstone）。

- created：创建该 generation 时的 revision。

- revs：存放该 generation 中存放的 revision 数组。

以下图来说明 keyIndex 结构体：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dv1a6u/1616136169269-c36dd19f-ee7b-40db-873d-c6da6b81da6e.jpeg)

如上图所示，存放的键为 test 的 keyIndex 结构。

它的 generations 数组有两条记录，其中 generations\[0]在 revision 1.0 时创建，当 revision2.1 的时候进行 tombstone 操作，因此该 generation 的 created 是 1.0；对应的 generations\[1]在 revision3.3 时创建，紧跟着就做了 tombstone 操作。

所以该 keyIndex.modifiled 成员存放的是 3.3，因为这是这条数据最后一次被修改的 revision。

一个已经被 tombstone 的 generation 是可以被删除的，如果整个 generations 数组都已经被删除空了，那么整个 keyIndex 记录也可以被删除了。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dv1a6u/1616136169303-1a5e569c-0985-495b-b196-ab09de02c4ad.jpeg)

如上图所示，keyIndex.compact(n)函数可以对 keyIndex 数据进行压缩操作，将删除满足 main revision < n 的数据。

- compact(2)：找到了 generations\[0]的 1.0 revision 的数据进行了删除。

- compact(3)：找到了 generations\[0]的 2.1 revision 的数据进行了删除，此时由于 generations\[0]已经没有数据了，所以这一整个 generation 被删除，原先的 generations\[1]变成了 generations\[0]。

- compact(4)：找到了 generations\[0]的 3.3 revision 的数据进行了删除。由于所有的 generation 数据都被删除了，此时这个 keyIndex 数据可以删除了。

treeIndex 结构

Etcd 中使用 treeIndex 来在内存中存放 keyIndex 数据信息，这样就可以快速的根据输入的 key 定位到对应的 keyIndex。

treeIndex 使用开源的 github.com/google/btree 来在内存中存储 btree 索引信息，因为用的是外部库，所以不打算就这部分做解释。而如果很清楚了前面 keyIndex 结构，其实这部分很好理解。

所有的操作都以 key 做为参数进行操作，treeIndex 使用 btree 根据 key 查找到对应的 keyIndex，再进行相关的操作，最后重新写入到 btree 中。
