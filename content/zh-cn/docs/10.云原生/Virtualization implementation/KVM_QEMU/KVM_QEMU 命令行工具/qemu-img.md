---
title: qemu-img
---

# 概述

> 参考：
> - [官方文档，工具-QEMU 磁盘镜像工具](https://www.qemu.org/docs/master/tools/qemu-img.html)

管理 VM image(虚拟机的镜像) 文件。qemu-img 工具可以创建、转换和修改 image 文件。所有 QEMU 支持的格式都可以处理

> 注意：不要使用 qemu-img 修改正在运行的虚拟机或者任何其他进程正在使用的 image 文件。

# Syntax(语法)

**qemu-img [Standard OPTIONS] COMMAND [COMMAND OPTIONS]**

### Standard OPTIONS

- **-U, --force-share** #

### COMMAND

- check \[-q] \[-f fmt] \[--output=ofmt] \[-r \[leaks | all]] \[-T src_cache] filename
- create \[-q] \[-f fmt] \[-o options] filename \[size]
- commit \[-q] \[-f fmt] \[-t cache] filename
- compare \[-f fmt] \[-F fmt] \[-T src_cache] \[-p] \[-q] \[-s] filename1 filename2
- convert \[-c] \[-p] \[-q] \[-n] \[-f fmt] \[-t cache] \[-T src_cache] \[-O output_fmt] \[-o options] \[-s snapshot_name] \[-S sparse_size] filename \[filename2 \[...]] output_filename
- info \[-f fmt] \[--output=ofmt] \[--backing-chain] filename
- map \[-f fmt] \[--output=ofmt] filename
- snapshot \[-q] \[-l | -a snapshot | -c snapshot | -d snapshot] filename
- rebase \[-q] \[-f fmt] \[-t cache] \[-T src_cache] \[-p] \[-u] -b backing_file \[-F backing_fmt] filename
- resize \[-q] filename \[+ | -]size
- amend \[-q] \[-f fmt] \[-t cache] -o options filename

# check - 检查 VM 镜像文件

**qemu-img check \[-q] \[-f fmt] \[--output=ofmt] \[-r \[leaks | all]] \[-T src_cache] filename**

对 VM 的镜像文件进行检查，只有 qcow2、qed、vdi 格式支持一致性检查。如果正常的话，基本输出信息如下

```bash
~]# qemu-img check centos8-1911.qcow2
No errors were found on the image.
34666/8192000 = 0.42% allocated, 70.85% fragmented, 66.63% compressed clusters
Image end offset: 2148270080
```

对磁盘映像文件名执行一致性检查。 该命令可以以“人类”或“ json”的 mt 格式输出。

如果指定了“ -r”，则 qemu-img 尝试修复在检查过程中发现的所有不一致之处。 “ -r 泄漏”仅修复群集泄漏，而“ -r all”修复所有类型的错误，选择错误的修复或隐藏已发生的损坏的风险更高。

仅格式“ qcow2”，“ qed”和“ vdi”支持一致性检查。

# create - 创建虚拟机的镜像文件

Note：

- 增量镜像 # 如果指定了 backing_file 选项。则新镜像只会记录与 backing_file 指定的基础镜像的差异，在这种情况下，无需指定 SIZE。除非使用 commit 命令将增量镜像与基础镜像合并，否则使用该虚拟机时，其 backing_file 镜像永远不会被修改，只会在增量镜像上，有大小的变化。

### Syntax(语法)

**qemu-img create \[-q] \[-f FMT] \[-o OTPIONS] FileName \[SIZE]**

- FileName # 创建一个名为 FileName 的虚拟机镜像文件

OPTIONS

- **-f FMT** # 指定第一个镜像文件的格式为 FMT。默认为 raw 格式。
- **-F FMT** # 指定第二个镜像文件的格式为 FMT。默认为 raw 格式。
- **-o OPTIONS** # 指定参数。可以使用 -o ? 来查看支持的参数。Note：查看可用的参数会根据不同的-f FMT，有不同的显示。Note：各个不同参数可以使用简化 OPTIONS 来指定
  - **backing_file=BaseFILE** # 指定基础镜像为 BaseFILE。简化为 -b
  - **size=SIZE** # 指定新镜像文件的大小为 SIZE

### EXAMPLE

- 创建一个 1T 容量的 qcow2 格式的镜像文件
  - **qemu-img create -f qcow2 -o size=1Ti data.bj-cs.qcow2**
- 基于 centos8-2004.qcow2 镜像文件，创建一个名为 lichenhao.bj-net.qcow2 的增量镜像文件
  - **qemu-img create -f qcow2 -b /var/lib/libvirt/images/backingFile/centos8-2004.qcow2 -F qcow2 lichenhao.bj-net.qcow2**
  - **qemu-img create -f qcow2 -b /var/lib/libvirt/images/backingFile/centos8-2004.qcow2 -F qcow2 lichenhao.bj-net.qcow2** # 其中 -o backing_file 可以简写为 -b
  - **qemu-img create -f qcow2 -b /var/lib/libvirt/images/backingFile/centos8-2004.qcow2 -o size=1Ti -F qcow2 lichenhao.bj-net.qcow2** # 创建时指定新镜像文件的大小

# convert - 转换 VM 镜像文件的格式

### Syntax(语法)

**qemu-img convert \[-c] \[-p] \[-q] \[-n] \[-f fmt] \[-t cache] \[-T src_cache] \[-O output_fmt] \[-o options] \[-s snapshot_name] \[-S sparse_size] filename \[filename2 \[...]] output_filename**

### EXAMPLE

- 压缩 test.qcow2 镜像,生成新的名为 test.qcow2.new 的镜像，新的镜像文件大小更小
  - `qemu-img convert -c -O qcow2 test.qcow2 test.qcow2.new`

```
# 压缩前的镜像信息
~]# qemu-img check centos8-2004.qcow2.src
No errors were found on the image.
8192000/8192000 = 100.00% allocated, 0.00% fragmented, 0.00% compressed clusters
Image end offset: 536953094144
# 压缩后的镜像信息
~]# qemu-img check centos8-2004.qcow2
No errors were found on the image.
36060/8192000 = 0.44% allocated, 91.45% fragmented, 89.83% compressed clusters
Image end offset: 1085603840
```

# info - 显示 VM 镜像文件的信息

**info \[-f fmt] \[--output=ofmt] \[--backing-chain] filename**

基本信息如下所示

```bash
~]# qemu-img info master-1.kg.tjiptv.net.qcow2
image: master-1.kg.tjiptv.net.qcow2
file format: qcow2
virtual size: 550G (590558003200 bytes)
disk size: 6.6G
cluster_size: 65536
Snapshot list:
ID        TAG                 VM SIZE                DATE       VM CLOCK
1         1565061772                0 2019-08-06 11:22:52   00:00:00.000
Format specific information:
compat: 1.1
    lazy refcounts: false
```

若虚拟机正在运行，执行该命令会出现如下提示

```bash
~]# qemu-img info duanyunhulian.qcow2
qemu-img: Could not open 'XXXXX.qcow2': Failed to get shared "write" lock
```

EXAMPLE

- 查看 test.qcow2 镜像的信息
  - **qemu-img info test.qcow2**

# resize - 设置 VM 磁盘容量

**qemu-img resize [-q] filename [+ | -]size**

EXAMPLE

- qemu-img resize cirros-0.3.6-x86_64-disk.img +40G # 给 test.qcow2 镜像的磁盘量容量添加 40G

# snapshot - 为 VM 创建、删除、应用、列出快照

**qemu-img snapshot \[-q] \[-l | -a snapshot | -c snapshot | -d snapshot] filename**
