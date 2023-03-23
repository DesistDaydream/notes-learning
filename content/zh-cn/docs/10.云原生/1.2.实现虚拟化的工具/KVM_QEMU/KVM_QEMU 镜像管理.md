---
title: KVM/QEMU 镜像管理
weight: 4
---

# 概述

> 参考：
>
> - [官方文档，系统模拟-磁盘镜像](https://qemu-project.gitlab.io/qemu/system/images.html)
> - [GitHub 文档，qemu/qemu/docs/interop/qcow2.txt](https://github.com/qemu/qemu/blob/master/docs/interop/qcow2.txt)
> - [Wiki,qcow](https://en.wikipedia.org/wiki/Qcow)
> - 其他
>   - <https://www.unixmen.com/qemu-kvm-using-copy-write-mode/>
>   - <https://opengers.github.io/virtualization/kvm-libvirt-qemu-5/>

KVM/QEMU 通过 [qemu-img](/docs/10.云原生/1.2.实现虚拟化的工具/KVM_QEMU/KVM_QEMU%20命令行工具/qemu-img.md) 命令行工具管理虚拟机镜像。

# QEMU Copy On Write

注意：

- 使用该特性创建出来虚拟机之后，整个快照链的根节点(i.e.backingfile 虚拟机)一定要不有任何更改，否则会导致基于其创建的其他所有虚拟机文件系统出现问题。比如变成 read-only 状态。

**QEMU Copy On Write(QEMU 写时复制，简称 QCOW)** 是 QEMU 创建的虚拟机使用的磁盘镜像文件的文件格式。

当使用 QCOW 时，不会对原始磁盘映像应用任何更改。所有更改都记录在其他的 QCOW 文件中。多个 QCOW 文件可以指向同一个镜像，而不会危及基本系统。QEMU/KVM 允许将 QCOW 文件的更改合并到原始图像中。

# 修改 backingfile 后，导致依赖 backingfile 的虚拟机的文件系统崩溃

https://www.cnblogs.com/fengrenzw/p/3383773.html

<https://www.cnblogs.com/fengrenzw/p/3383773.html>

我们知道 qcow2 的磁盘格式可以带来很大的便利性，因为部署的时候可以减少大量的时间、空间，可以增量备份、快照等非常诱人的特性。

因为下边可能会有点绕：

backing_file：后端，母镜像

qcow2：前端，子镜像

在使用的时候可能会遇到一种情况，就是使用 backing_file 时，如果修改了 backing_file，“可能”会导致前端的 qcow2 的崩溃，出现这种问题个人觉得是很正常的，并且是可以完全避免的。所以，在 openstack 在使用 qcow2 的过程中会使用 glance 镜像管理来保证它的安全和完整性，我们在使用 qcow2 的时候也务必不回去修改它。

至于为什么会出现这种现象，下面简单分析一下，可能会有些纰漏、错误，但感觉整体思路上不会有太大的偏差。

什么是 qcow2？

之前的博客也讲述过，qcow2：就是 qemu 的 cow 磁盘的第 2 版。既然是 cow，必然是创建的前端磁盘内容是“空”的，即只有 qcow2 磁盘格式的数据结构（当然包含 backing_file 的指针），而不包含任何磁盘内应该存放的实际内容。

我们启动虚拟机的时候，指定的是 qcow2，但同时会加载 backing_file（使用 qmp，info block 可以查看，或者/proc/$pid/fd/）。当读取文件的时候，根据 qcow2 内部的指针指向 backing_file，读取 backing_file 磁盘块上的内容。如果需要修改文件，则修改后的文件会保存到 qcow2 文件上。

那我们在使用 backing_file 特性时，再修改 backing_file（后端）时可能就出现大概三种情况：

- backing_ing 删除或修改了 qcow2 中没修改过的内容

因为 qcow2 本来就没有什么数据，所有能查看到的数据都是通过 backing_file 的指针查看到的，所以当 backing_file 修改了，qcow2 还是直接去读 backing_file，就相当于同步了，并不会有冲突或腐败。

- backing_ing 删除或修改了 qcow2 中修改过的内容

因为 qcow2 的 cow 机制，修改后的文件会保存到 qcow2 文件中，所以 backing_file 修改不会对 qcow2 文件造成任何影响，因为压根就没去读 backing_file。（但有个前提，修改的幅度不应太大，如果文件的 inode 也变了，可能会造成冲突和错误（不过，这都不是重点！）

- backing_ing 创建了 qcow2 中没有的内容

这种情况有点复杂，因为创建文件肯定会影响到文件系统的 inode，如果这个 inode 没有在 qcow2 做修改的话，会直接读取 backing_file,自然能够找到新添加的文件，如果 qcow2 的文件系统内的 inode 做了修改，我按照我自己的 inode 去找 backing_file 中的 block 发现对应不上，就会照成文件系统的损坏甚至崩溃；另一种情况是在 qcow2 中删除了一些文件或者将某个目录清空，然而在 backing_file 又在这个目录写了一些东西，qcow2 中的 inode 中可能去查看到这个数据，但因为 qcow2 中的这块数据已经修改，不会去查看 backing_file，就会导致有了 inode，但找不到 block。无论是有 inode 查看不到 block 还是 inode 删除，但 block 还在，不是我们所希望的。

所以，如果使用 qcow2 的 backing_file，请务必保证其安全和完整性。
