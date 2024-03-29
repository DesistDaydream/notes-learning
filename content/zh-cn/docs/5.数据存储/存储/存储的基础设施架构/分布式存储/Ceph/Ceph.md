---
title: Ceph
weight: 1
---

# 概述

> 参考：
>
> - [官网](https://ceph.io/)
> - [官方文档](https://docs.ceph.com/en/latest/)
> - [Wiki，Ceph](<https://en.wikipedia.org/wiki/Ceph_(software)>)
> - <https://blog.csdn.net/younger_china/article/details/73410727>

Ceph 是一个开源的分布式存储系统，可以提供 对象存储、快存储、文件存储 能力。是一个 Software Defined Storage(软件定义存储) 的代表性产品。

一个 Ceph 存储集群至少需要 Ceph Monitor、Ceph Manager、Ceph OSD 这三个组件；如果要运行 Ceph 文件系统客户端，则也需要 Ceph MDS。

- **Monitor** # **Ceph Monitor(Ceph 监视器，简称 ceph-mon)** 负责维护集群状态的映射关系。通常至少需要 3 个 ceph-mon 以实现高可用，多节点使用 Paxos 算法达成共识。
  - 可以这么说，Ceph 集群就是指 ceph-mon 集群。ceph-mon 负责维护的集群状态，就是用来提供存储服务的。
  - ceph-mon 映射、ceph-mgr 映射、ceph-osd 映射、ceph-mds 映射、ceph-crush 映射。这些映射是 Ceph 守护进程相互协调所需的关键集群状态，说白了，就是**映射关系**。
    - 这里的映射，英文用的是 Map，其实也有地图的意思，就是表示这个集群有多少个 ceph-mon、有多少个 ceph-mgr 等等，还有底层对象属于哪个 PG，等等等等，这些东西构成了一副 Ceph 的运行图。
  - ceph-mon 还负责管理守护进程和客户端之间的身份验证。
- **Manager** # **Ceph Manager(Ceph 管理器，简称 ceph-mgr)** 负责跟踪运行时指标和 Ceph 集群的当前状态，包括存储利用率、性能、系统负载等。通常至少需要 2 个 ceph-mgr 以实现高可用。
  - ceph-mgr 可以提供 Web 管理页面、关于 Ceph 集群的 Prometheus 格式的监控指标
- **OSD Daemon** # **Ceph OSD Daemon(Ceph OSD 守护进程，简称 ceph-osd)** 负责向 OSD 读写数据、处理数据复制、恢复、重新平衡，并通过检查其他 ceph-osd 的心跳向 ceph-mon 和 ceph-mgr 提供一些监控信息。通常至少需要 3 个 ceph-osd 以实现高科用。
  - **Object Storage Device(对象存储设备，简称 OSD)** 是一个物理或逻辑上的存储单元(比如一块硬盘)，这是 Ceph 得以运行的最基本的存储单元。
    - 有的时候，人们容易把 OSD 理解为 Ceph OSD Daemon，这俩是有本质区别的。因为在最早的时候，OSD 有两种含义，一种是 `Object Storage Device` 另一种是 `Object Storage Daemon`。由于这种称呼的模糊性，后来就将 Object Storage daemon 扩展为 OSD Daemon。OSD 则仅仅代表 Object Storage Device。只不过运行 OSD Daemon 的程序名称，依然沿用了 osd 的名字。
  - 注意，为了让每一个 OSD 都可以被单独使用并管理，所以每个 OSD 都有一个对应的 ceph-osd 进程来管理。一般情况，Ceph 集群中每个节点，除了系统盘做 Raid 以外，其他硬盘都会单独作为 OSD 使用，且一个节点会有大量磁盘来对应 OSD。
- **MDS** # **Ceph Metadata Server(Ceph 元数据服务器，简称 ceph-mds)** 代表 Ceph 文件系统元数据。ceph-mds 允许 POSIX 文件系统用户执行基本命令(比如 ls、find 等)，而不会给 Ceph 集群带来巨大负担。
  - 注意，Ceph 提供的 块存储 和 对象存储 功能并不使用 ceph-mds。

## 架构

> 参考：
>
> - [官方文档，架构](https://docs.ceph.com/en/latest/architecture/)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/sakrws/1630769971104-82bcc0c6-1dbd-4c47-b986-3e5b8321aac0.png)

**其实 Ceph 本身就是一个对象存储**，基于**RADOS** 实现，并通过 Ceph Client 为上层应用提供了通用的 块存储、文件存储、对象存储 的调用接口。

- **RADOS** # **Reliable Autonomic Distributed Object Store(可靠的、自动化的分布式对象存储，简称 RADOS)** 是一种由多个主机组成、由 CRUSH 算法实现数据路由的，分布式对象存储系统。是 Ceph 的底层存储系统。
  - OSD 是组成 RADOS 的基本存储单元。
- **Ceph Client** # **Ceph 客户端**。是可以访问 Ceph 存储集群(即 RADOS) 的 Ceph 组件的集合。
  - **LIBRADOS** # **Library RADOS(RADOS 库，简称 librados)**。应用程序可以调用 librados 以直接访问 RADOS。当我们使用 Ceph 时，Ceph 实际上是调用 librados 的 API(这是一个 rpc 接口)，将提交的文件切分为固定大小的数据，存放到 RADOS 中。
    - 同时，我们自己也可以使用 librados 开发出类似 ceph-rgw、ceph-rbd 这种应用程序以实现个性化需求。
  - **RADOSGW** # **RADOS Gateway(RADOS 网关，简称 radosgw)**。使用 librados 实现的应用程序，可以提供兼容 S3 和 Swift 对象存储的接口
  - **RBD** # **RADOS Block Device(RADOS 块设备，简称 RBD)**。使用 librados 实现的应用程序，为 Linux 内核 和 QEMU/KVM 提供一个可靠且完全分布式的块存储设备。
  - **CEPH FS** # **Ceph File System(Ceph 文件系统，简称 CFS)**。直接使用 RADOS 实现一个符合 POSIX 的分布式文件系统，带有 Linux 内核客户端并支持 FUSE，可以直接挂载使用。甚至可以进一步抽象，实现 NFS 功能。

## Ceph 数据写入流程

Ceph 集群从 Ceph 的客户端接收到的数据后，将会切分为一个或多个固定大小的 **RADOS Object(RADOS 对象)**。Ceph 使用 **Controlled Replication Under Scalable Hashing(简称 CRUSH)** 算法计算出 RADOS 对象应该放在哪个 **Placement Group(归置组，简称 PG)**，并进一步计算出，应该由哪个 ceph-osd 来处理这个 PG 并将 PG 存储到指定的 OSD 中。ceph-osd 会通过存储驱动器处理 RADOS 对象的 读、写 和 复制操作。

> 注意：当创建完 Ceph 集群后，会有一个默认的 Pool，Pool 是用来对 PG 进行分组的，且 PG 必须属于一个组，不可独立存在。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/sakrws/1630834243384-b650e1e5-1c84-4846-bdc5-9180a361fb09.png)

RADOS 对象有如下几个部分组成

- **Object Identify(对象标识符，简称 OID)** # OID 在整个 Ceph 集群中是唯一。
- **Binary Data(二进制数据数据)** # 对象的数据
- **Metadata(元数据)** # 元数据的语义完全取决于 Ceph 客户端。例如，CephFS 使用元数据来存储文件属性，如文件所有者、创建日期、上次修改日期等。

ceph-osd 将数据作为对象存储在平坦的命名空间中 (例如，没有目录层次结构)。对象具有标识符，二进制数据和由一组名称/值对组成的元数据。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/sakrws/1630808425695-75766062-7570-47f0-9ae4-916c7819d113.png)

# RADOS

与传统分布式存储不同，传统分布式存储中的 NameNode 极易形成性能瓶颈。基于此，RADOS 设计了一种新的方式来快速找到对象数据。RADOS 中并不需要 NameNode 来存储每个对象的元数据，RADOS 中的对象，都是通过 **Controlled Replication Under Scalable Hashing(简称 CRUSH)** 算法来快速定位的。

## bluestore

这是 Ceph 所管理的 OSD 的文件系统类型

# Ceph 的存储能力

## 块存储

Ceph 通过 RDB 提供块存储能力

## 文件存储

Ceph 通过 CEPHFS 提供文件存储能力

## 对象存储

RADOS Gateway 简称 radosgw，Ceph 通过 radosgw 程序，可以对外提供标准的 S3 或 swift 接口，以实现主流对象存储功能。很多时候，radosgw 程序运行的进程称为 ceph-rgw

# Ceph Manager

> 参考：
>
> - [官方文档,Ceph 管理器](https://docs.ceph.com/en/latest/mgr/)

Ceph Manager 是通过一个名为 ceph-mgr 的二进制程序以守护进程运行的管理器。ceph-mgr 可以向外部监控和管理系统提供额外的监控和接口。

ceph-mgr 曾经是 ceph-mon 的一部分，自 luinous(12.x) 版本依赖，ceph-mgr 独立出来，成为 Ceph 集群的必选组件。

## Dashboard 模块

Dashboard 模块是一个内置的基于 Web 的 Ceph 管理和监控程序，通过它可以检查和管理 Ceph 集群中的各个方面和资源。默认监听 `8443` 端口

在 Dashboard 模块中，提供了一组用于管理集群的 RESTful 风格的 API 接口。这组 API 位于 `/api` 路径下。详见《[API](</docs/5.数据存储/1.存储/存储的基础设施架构/Distributed%20Storage(分布式存储)/Ceph/API.md>>)》章节

## Prometheus 模块

启动 Prometheus 模块后，ceph-mgr 默认在 `9283` 端口上暴露 Prometheus 格式的监控指标。

# Ceph RADOSGW

默认监听 `7480` 端口
