---
title: Containerd Image
---

# 概述

> 参考：
>
> - <https://blog.frognew.com/tags/containerd.html>
> - [重学容器 09: Containerd 是如何存储容器镜像和数据的](https://blog.frognew.com/2021/06/relearning-container-09.html)

这是一个 /var/lib/containerd 目录的最基本组成：

```bash
~]# tree -L 1
.
├── io.containerd.content.v1.content
├── io.containerd.metadata.v1.bolt
├── io.containerd.runtime.v2.task
├── io.containerd.snapshotter.v1.overlayfs
├── ...... 等
└── tmpmounts
```

初始情况，Containerd 会加载部分插件，对应了 content、snapshot、metadata、runtime 等等插件。通过 `ctr plugin ls` 命令可以发现，目录名称与插件名称是一致的。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/abf521b3-7aea-4cfc-a5b4-343877ff4ccb/1637717831444-46d3e321-18da-456c-a543-7df7415d3689.png)

这是一个只有一个 lchdzh/k8s-debug:v1 镜像的 /var/lib/containerd 目录：

```bash
~]# tree /var/lib/containerd -L 5
.
├── io.containerd.content.v1.content
│   ├── blobs
│   │   └── sha256
│   │       ├── 02daccf1684b499e99c258348d492c5f0ea086174d2f0d430791d4f902ae4f71
│   │       ├── 188c0c94c7c576fff0792aca7ec73d67a2f7f4cb3a6e53a84559337260b36964
│   │       ├── 5f9b9d9c910519d9a4b1e06f031672e14acf9bcc288ed7e3ed3842916ed4394d
│   │       ├── c690d4fd64d6622c3721a1db686c2e4cfb559dd1d9f9ff825584a8f56ec02c7f
│   │       ├── df727e3daae2c57da7071b4056d328d4bbb9d6a913e469d8f07b58e35a5cff96
│   │       └── ee24b921ba004624b350e7f140e68c6a7d8297bb815b4ca526979a7e66cec15a
│   └── ingest
├── io.containerd.metadata.v1.bolt
│   └── meta.db
├── io.containerd.runtime.v1.linux
├── io.containerd.runtime.v2.task
├── io.containerd.snapshotter.v1.aufs
│   └── snapshots
├── io.containerd.snapshotter.v1.btrfs
├── io.containerd.snapshotter.v1.native
│   └── snapshots
├── io.containerd.snapshotter.v1.overlayfs
│   ├── metadata.db
│   └── snapshots
│       ├── 1
│       │   ├── fs
│       │   │   ├── bin
│       │   │   ├── dev
│       │   │   ├── etc
│       │   │   ├── home
│       │   │   ├── lib
│       │   │   ├── media
│       │   │   ├── mnt
│       │   │   ├── opt
│       │   │   ├── proc
│       │   │   ├── root
│       │   │   ├── run
│       │   │   ├── sbin
│       │   │   ├── srv
│       │   │   ├── sys
│       │   │   ├── tmp
│       │   │   ├── usr
│       │   │   └── var
│       │   └── work
│       ├── 2
│       │   ├── fs
│       │   │   ├── bin
│       │   │   ├── etc
│       │   │   ├── lib
│       │   │   ├── run
│       │   │   ├── sbin
│       │   │   ├── usr
│       │   │   └── var
│       │   └── work
│       ├── 3
│       │   ├── fs
│       │   │   ├── etc
│       │   │   ├── root
│       │   │   └── usr
│       │   └── work
│       └── 4
│           ├── fs
│           │   └── usr
│           └── work
└── tmpmounts
```

这些插件具有对应的目录：

- **io.containerd.content.v1.content/** #
  - **./blobs/sha256/** # OCI Image 的 blob 文件存放路径。其中 tar.gzip 类型的 blob 文件(即.镜像层的压缩文件)将会被解压到 io.containerd.snapshotter.v1.overlayfs/snapshots/ 目录中
  - **./ingest/** # 当 pull 一个镜像时，会在该目录形成缓存，并逐渐在 blobs 目录中生成 blob 文件，pull 完之后，该目录将会清空。
- **io.containerd.grpc.v1.introspection/** #
- **io.containerd.metadata.v1.bolt/** #
  - **./meta.db** # 这是一个 boltdb 的持久化文件。保存了 OCI Image 标准中 bolts 目录下的文件的组织信息。
- **io.containerd.runtime.v2.task/** #
  - .**/default/** # default 名称空间中运行的容器
  - **./moby/** # moby 名称空间中运行的容器
- **io.containerd.snapshotter.v1.overlayfs/** #
  - **./metadata.db** #
  - .**/snapshots/INT/** # OCI Image 的 Layers 文件解压后的文件系统存放路径。每个镜像层都使用一个 INT 类型的数字作为目录名，目录中即是文件系统内容。
    - 当运行一个容器时，就是挂载的这些目录
- **tmpmounts/** #

## meta.db 文件解析

meta.db 是 boltdb 文件。通过 `go.etcd.io/bbolt` 库，使用 go 语言获取一下 meta.db 中的数据

```go
package main

import (
 "fmt"
 "log"

 bolt "go.etcd.io/bbolt"
)

func iteratingAll(bucket *bolt.Bucket, space string) {
 space = space + "  "
 bucket.ForEach(func(k, v []byte) error {
  if v == nil {
   fmt.Printf("%s[%s]: \n", space, k)
   // 嵌套迭代
   iteratingAll(bucket.Bucket([]byte(k)), space)
  } else {
   fmt.Printf("%s%s=%s\n", space, k, v)
  }
  return nil
 })
}

func iteratingBucket(bucket *bolt.Bucket, space string) {
 space = space + "  "
 bucket.ForEach(func(k, v []byte) error {
  if v == nil {
   fmt.Printf("%s[%s]\n", space, k)
   // 嵌套迭代
   iteratingBucket(bucket.Bucket([]byte(k)), space)
  }
  return nil
 })
}

func r_transactions(tx *bolt.Tx) {
 // 迭代根中的所有 Bucket
 tx.ForEach(func(bucket_name []byte, bucket *bolt.Bucket) error {
  fmt.Printf("root_bucket=%v\n", string(bucket_name))

  space := ""
  // 迭代所有 Bucket
  iteratingBucket(bucket, space)

  // 迭代所有 Bucket 及其 K/V。
  // 可以指定名称空间，以跌点特定名称空间写下的 Bucket 及其 K/V
  // bucket = tx.Bucket(bucket_name).Bucket([]byte("default"))
  iteratingAll(bucket, space)
  return nil
 })
}

func main() {
 db, err := bolt.Open("meta.db", 0600, nil)
 if err != nil {
  log.Fatal(err)
 }
 defer db.Close()

 db.View(func(tx *bolt.Tx) error {
  r_transactions(tx)
  return nil
 })
}
```

首先先看 Containerd 设计的所有 Bucket

```bash
root_bucket=v1
  [default]
    [content]
      [blob]
        [sha256:02daccf1684b499e99c258348d492c5f0ea086174d2f0d430791d4f902ae4f71]
          [labels]
        [sha256:188c0c94c7c576fff0792aca7ec73d67a2f7f4cb3a6e53a84559337260b36964]
          [labels]
        [sha256:5f9b9d9c910519d9a4b1e06f031672e14acf9bcc288ed7e3ed3842916ed4394d]
          [labels]
        [sha256:c690d4fd64d6622c3721a1db686c2e4cfb559dd1d9f9ff825584a8f56ec02c7f]
          [labels]
        [sha256:df727e3daae2c57da7071b4056d328d4bbb9d6a913e469d8f07b58e35a5cff96]
          [labels]
        [sha256:ee24b921ba004624b350e7f140e68c6a7d8297bb815b4ca526979a7e66cec15a]
          [labels]
      [ingests]
    [images]
      [docker.io/lchdzh/k8s-debug:v1]
        [target]
    [leases]
    [snapshots]
      [overlayfs]
        [sha256:9fd95275d19bae7e959c6299f4fcf7b98c2dea021efcc2ad8fe4aab2130423cd]
          [children]
          [labels]
        [sha256:ace0eda3e3be35a979cec764a3321b4c7d0b9e4bb3094d20d3ff6782961a8d54]
          [children]
          [labels]
        [sha256:e75c38ede0d09abef411d0e8c438542f98f0114babb32577380f45898849f4b4]
          [children]
          [labels]
        [sha256:ef906066c23fd6da8b188fa18250cdf76030d529d1a9581f8af4bb2c4f775fa8]
          [labels]
  [moby]
    [containers]
      [4c5ec4bc9717bb9fd2a2ea7b507ac3c0e16da95fa87974152f0fe3b3a653cef9]
        [labels]
        [runtime]
      [697033ad3cc1a7c1d3e8d87c89ced26436e0c6e738524b107dceed17bb3d77e8]
        [labels]
        [runtime]
    [leases]
```

在 Containerd 源码 [./metdata/buckets.go](https://github.com/containerd/containerd/blob/main/metadata/buckets.go) 的开头注释中，描述了 meta.db 文件中的数据结构，与上面通过代码获取到的内容是一致的。这个结构总结一下，就是：

```bash
Version/Namespace/Object/Key -> Value
```

- **Version** # 当前版本始终是 v1
- **Namespace** # Object 所属的名称空间
- **Object** # Object(对象) 在数据库中就是具有 K/V 的 Bucket。由于 Containerd 的插件机制，不同类型的 blob 文件，是由不同插件管理的。所以，Object 由可以由 Plugin/Type 组成。Containerd 的对象就是其所管理的原子单位，也就是各种元数据信息：
  - **content/blob/** # OCI Image 规范的 blob 文件信息
  - **image/IMAGE** # 镜像名称
  - **snapshots/overlayfs/** # 镜像层解压后的文件系统信息
  - ......等等
  - 这里有两个特殊的 Object
    - labels # 用来存储其 父 Bucket 的额外属性。比如 content/blob/DIGEST/ 下就有 labels，用来描述 blob 数据的额外属性。
    - indexes # 暂时用不上。为将来的扩展预留。
- **key** # 特定于 Object 的键，用来描述 Object 的属性。比如这个对象的 创建时间、更新时间、媒体类型、大小 等等
  - 其中 containerd.io/uncompressed 的值 和 snapshots.overlayfs 字段下的很多内容 与 metadata.db 中的数据互相关联
  - 在下面的代码块中可以看看这些 Bucket 中的所有 K/V：

```bash
root_bucket=v1
  [default]:
    [content]:
      [blob]:
        [sha256:02daccf1684b499e99c258348d492c5f0ea086174d2f0d430791d4f902ae4f71]:
          createdat=ׯ
          [labels]:
            containerd.io/distribution.source.docker.io=lchdzh/k8s-debug
            containerd.io/uncompressed=sha256:6e63a43fa96c6ea85d34c23db4c28b76ecda01c03aa721f6a3355b04501bdc58
          size=ª
          updatedat=ׯ
        [sha256:188c0c94c7c576fff0792aca7ec73d67a2f7f4cb3a6e53a84559337260b36964]:
          createdat=ׯṉ퀿
          [labels]:
            containerd.io/distribution.source.docker.io=lchdzh/k8s-debug
            containerd.io/uncompressed=sha256:ace0eda3e3be35a979cec764a3321b4c7d0b9e4bb3094d20d3ff6782961a8d54
          size=
          updatedat=򿀿
        [sha256:5f9b9d9c910519d9a4b1e06f031672e14acf9bcc288ed7e3ed3842916ed4394d]:
          createdat=ׯ󁀿
          [labels]:
            containerd.io/distribution.source.docker.io=lchdzh/k8s-debug
            containerd.io/uncompressed=sha256:6f5211c02ff0b7e40b9ca7c5f62cc8732647b046e22cc5046053412d1fef97f6
          size=ꍂ
          updatedat=ׯÿÿ
        [sha256:c690d4fd64d6622c3721a1db686c2e4cfb559dd1d9f9ff825584a8f56ec02c7f]:
          createdat=ׯ󨅁
          [labels]:
            containerd.io/distribution.source.docker.io=lchdzh/k8s-debug
            containerd.io/gc.ref.snapshot.overlayfs=sha256:ef906066c23fd6da8b188fa18250cdf76030d529d1a9581f8af4bb2c4f775fa8
          size=!
          updatedat=ׯ
        [sha256:df727e3daae2c57da7071b4056d328d4bbb9d6a913e469d8f07b58e35a5cff96]:
          createdat=ׯ򟡥&ÿÿ
          [labels]:
            containerd.io/distribution.source.docker.io=lchdzh/k8s-debug
            containerd.io/uncompressed=sha256:4041c1a8637589d2c872e14d1068376c5e21bf96a837fa2225f91066e84b1e55
          size=쁷
          updatedat=ׯ󰔾fÿÿ
        [sha256:ee24b921ba004624b350e7f140e68c6a7d8297bb815b4ca526979a7e66cec15a]:
          createdat=ׯ񮶝뀿
          [labels]:
            containerd.io/distribution.source.docker.io=lchdzh/k8s-debug
            containerd.io/gc.ref.content.config=sha256:c690d4fd64d6622c3721a1db686c2e4cfb559dd1d9f9ff825584a8f56ec02c7f
            containerd.io/gc.ref.content.l.0=sha256:188c0c94c7c576fff0792aca7ec73d67a2f7f4cb3a6e53a84559337260b36964
            containerd.io/gc.ref.content.l.1=sha256:df727e3daae2c57da7071b4056d328d4bbb9d6a913e469d8f07b58e35a5cff96
            containerd.io/gc.ref.content.l.2=sha256:5f9b9d9c910519d9a4b1e06f031672e14acf9bcc288ed7e3ed3842916ed4394d
            containerd.io/gc.ref.content.l.3=sha256:02daccf1684b499e99c258348d492c5f0ea086174d2f0d430791d4f902ae4f71
          size=
          updatedat=ׯ񰿖Kÿÿ
      [ingests]:
    [images]:
      [docker.io/lchdzh/k8s-debug:v1]:
        createdat=ׯ
        [target]:
          digest=sha256:ee24b921ba004624b350e7f140e68c6a7d8297bb815b4ca526979a7e66cec15a
          mediatype=application/vnd.docker.distribution.manifest.v2+json
          size=
        updatedat=ׯ
    [leases]:
    [snapshots]:
      [overlayfs]:
        [sha256:9fd95275d19bae7e959c6299f4fcf7b98c2dea021efcc2ad8fe4aab2130423cd]:
          [children]:
            sha256:ef906066c23fd6da8b188fa18250cdf76030d529d1a9581f8af4bb2c4f775fa8=
          createdat=ׯ󵭁3ÿÿ
          [labels]:
            containerd.io/snapshot.ref=sha256:9fd95275d19bae7e959c6299f4fcf7b98c2dea021efcc2ad8fe4aab2130423cd
          name=default/6/sha256:9fd95275d19bae7e959c6299f4fcf7b98c2dea021efcc2ad8fe4aab2130423cd
          parent=sha256:e75c38ede0d09abef411d0e8c438542f98f0114babb32577380f45898849f4b4
          updatedat=ׯ󵭁3ÿÿ
        [sha256:ace0eda3e3be35a979cec764a3321b4c7d0b9e4bb3094d20d3ff6782961a8d54]:
          [children]:
            sha256:e75c38ede0d09abef411d0e8c438542f98f0114babb32577380f45898849f4b4=
          createdat=ׯ񿯸ÿÿ
          [labels]:
            containerd.io/snapshot.ref=sha256:ace0eda3e3be35a979cec764a3321b4c7d0b9e4bb3094d20d3ff6782961a8d54
          name=default/2/sha256:ace0eda3e3be35a979cec764a3321b4c7d0b9e4bb3094d20d3ff6782961a8d54
          updatedat=ׯ񿯸ÿÿ
        [sha256:e75c38ede0d09abef411d0e8c438542f98f0114babb32577380f45898849f4b4]:
          [children]:
            sha256:9fd95275d19bae7e959c6299f4fcf7b98c2dea021efcc2ad8fe4aab2130423cd=
          createdat=ׯ󮏻ÿÿ
          [labels]:
            containerd.io/snapshot.ref=sha256:e75c38ede0d09abef411d0e8c438542f98f0114babb32577380f45898849f4b4
          name=default/4/sha256:e75c38ede0d09abef411d0e8c438542f98f0114babb32577380f45898849f4b4
          parent=sha256:ace0eda3e3be35a979cec764a3321b4c7d0b9e4bb3094d20d3ff6782961a8d54
          updatedat=ׯ󮏻ÿÿ
        [sha256:ef906066c23fd6da8b188fa18250cdf76030d529d1a9581f8af4bb2c4f775fa8]:
          createdat=ׯ
          [labels]:
            containerd.io/snapshot.ref=sha256:ef906066c23fd6da8b188fa18250cdf76030d529d1a9581f8af4bb2c4f775fa8
          name=default/8/sha256:ef906066c23fd6da8b188fa18250cdf76030d529d1a9581f8af4bb2c4f775fa8
          parent=sha256:9fd95275d19bae7e959c6299f4fcf7b98c2dea021efcc2ad8fe4aab2130423cd
          updatedat=ׯ
```

## metadata.db 文件解析

使用相同的代码，可以获取文件内容如下：

```bash
root_bucket=v1
  [parents]
  [snapshots]
    [default/10/e4cf8bc98463185da7c6468e8c575e8022b6ad9bea4d3d0022017689fde711f8]
    [default/2/sha256:ace0eda3e3be35a979cec764a3321b4c7d0b9e4bb3094d20d3ff6782961a8d54]
      [labels]
    [default/4/sha256:e75c38ede0d09abef411d0e8c438542f98f0114babb32577380f45898849f4b4]
      [labels]
    [default/6/sha256:9fd95275d19bae7e959c6299f4fcf7b98c2dea021efcc2ad8fe4aab2130423cd]
      [labels]
    [default/8/sha256:ef906066c23fd6da8b188fa18250cdf76030d529d1a9581f8af4bb2c4f775fa8]
      [labels]

```

```yaml
root_bucket=v1
  [parents]:
    =default/4/sha256:e75c38ede0d09abef411d0e8c438542f98f0114babb32577380f45898849f4b4
    =default/6/sha256:9fd95275d19bae7e959c6299f4fcf7b98c2dea021efcc2ad8fe4aab2130423cd
    =default/8/sha256:ef906066c23fd6da8b188fa18250cdf76030d529d1a9581f8af4bb2c4f775fa8
    =default/10/e4cf8bc98463185da7c6468e8c575e8022b6ad9bea4d3d0022017689fde711f8
  [snapshots]:
    [default/10/e4cf8bc98463185da7c6468e8c575e8022b6ad9bea4d3d0022017689fde711f8]:
      createdat=װ<¢/1ÿÿ
      id=
      kind=
      parent=default/8/sha256:ef906066c23fd6da8b188fa18250cdf76030d529d1a9581f8af4bb2c4f775fa8
      updatedat=װ<¢/1ÿÿ
    [default/2/sha256:ace0eda3e3be35a979cec764a3321b4c7d0b9e4bb3094d20d3ff6782961a8d54]:
      createdat=ׯ򜃵ÿÿ
      id=
      inodes=
      kind=
      [labels]:
        containerd.io/snapshot.ref=sha256:ace0eda3e3be35a979cec764a3321b4c7d0b9e4bb3094d20d3ff6782961a8d54
      size=r
      updatedat=ׯ򜃵ÿÿ
    [default/4/sha256:e75c38ede0d09abef411d0e8c438542f98f0114babb32577380f45898849f4b4]:
      createdat=ׯÿÿ
      id=
      inodes=¶4
      kind=
      [labels]:
        containerd.io/snapshot.ref=sha256:e75c38ede0d09abef411d0e8c438542f98f0114babb32577380f45898849f4b4
      parent=default/2/sha256:ace0eda3e3be35a979cec764a3321b4c7d0b9e4bb3094d20d3ff6782961a8d54
      size=j
      updatedat=ׯÿÿ
    [default/6/sha256:9fd95275d19bae7e959c6299f4fcf7b98c2dea021efcc2ad8fe4aab2130423cd]:
      createdat=ׯ󵰲‿
      id=
      inodes=
      kind=
      [labels]:
        containerd.io/snapshot.ref=sha256:9fd95275d19bae7e959c6299f4fcf7b98c2dea021efcc2ad8fe4aab2130423cd
      parent=default/4/sha256:e75c38ede0d09abef411d0e8c438542f98f0114babb32577380f45898849f4b4
      size=²
      updatedat=ׯ󵰲‿
    [default/8/sha256:ef906066c23fd6da8b188fa18250cdf76030d529d1a9581f8af4bb2c4f775fa8]:
      createdat=ׯ
      id=
      inodes=
      kind=
      [labels]:
        containerd.io/snapshot.ref=sha256:ef906066c23fd6da8b188fa18250cdf76030d529d1a9581f8af4bb2c4f775fa8
      parent=default/6/sha256:9fd95275d19bae7e959c6299f4fcf7b98c2dea021efcc2ad8fe4aab2130423cd
      size=
      updatedat=ׯ
```

这文件是干啥的？？？没看懂，居然还有空键。看着非常像 meta.db 文件中的 snapshots 这个桶下的内容，但是又有一点点区别

## 镜像的使用

当一个容器运行起来后，查看其联合挂载信息：

```bash
~]# mount
overlay on /run/containerd/io.containerd.runtime.v2.task/default/ff15a5b9cfcfe53bb0a960c0fa76c4488f240d120555771d3f4309417567480d/rootfs type overlay (rw,relatime,lowerdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/18/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/17/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/16/fs,upperdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/31/fs,workdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/31/work,xino=off)
```

可以看到这些都是将 /var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/INT/\* 目录中镜像文件系统联合挂载的。最终都被联合挂载到容器运行时的 rootfs 目录 /run/containerd/io.containerd.runtime.v2.task/default/ff15a5b9cfcfe53bb0a960c0fa76c4488f240d120555771d3f4309417567480d/rootfs 中
