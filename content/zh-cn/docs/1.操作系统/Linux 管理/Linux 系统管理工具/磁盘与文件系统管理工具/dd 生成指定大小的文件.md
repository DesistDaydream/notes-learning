---
title: dd 生成指定大小的文件
---

# 概述

# Syntax(语法)

> 参考:
>
> - https://blog.csdn.net/menogen/article/details/38059671

**dd \[OPTIONS\[=FLAGS]]**

OPTIONS

- **bs**=BYTES # 每次读取和写入的字节数
- **cbs**=BYTES # convert BYTES bytes at a time
- **conv**=CONVS # convert the file as per the comma separated symbol list
- **count**=N # 读取的 block 数，block 的大小由 ibs 指定（只针对输入参数）
- **ibs**=BYTES # read up to BYTES bytes at a time (default: 512)
- **if**=FILE # 指定输入文件。默认从标准输入读取。
  - /dev/zero 是 Linux 的一个伪文件，它可以产生连续不断的 null 流（二进制的 0）
- **iflag**([FLAGS](#FLAGS)) # 使用 FLAGS 来控制输入(读取数据)时的行为特征。多个 FLAG 以 , 分隔
- **obs**=BYTES # write BYTES bytes at a time (default: 512)
- **of**=FILE # 指定输出文件。默认输出到标准输出。
- **oflag**=FLAGS # 使用 iflag 来控制输出(写入数据)时的行为特征。多个 FLAG 以 , 分隔
- **seek**=N # skip N obs-sized blocks at start of output
- **skip**=N # skip N ibs-sized blocks at start of input
- **status**=LEVEL # The LEVEL of information to print to stderr; 'none' suppresses everything but error messages, 'noxfer' suppresses the final transfer statistics, 'progress' shows periodic transfer statistics

## FLAGS

- append # 追加模式(仅对输出有意义；隐含了 conv=notrunc)
- direct # 使用直接 I/O 存取模式，即跳过缓存，不操作内存，而是直接操作磁盘
- directory # 除非是目录，否则 directory 失败
- dsync # 使用同步 I/O 存取模式
- sync # 与上者类似，但同时也对元数据生效
- fullblock # 为输入积累完整块(仅 iflag)
- nonblock # 使用无阻塞 I/O 存取模式
- noatime # 不更新存取时间
- nocache # 丢弃缓存数据
- noctty # 不根据文件指派控制终端
- nofollow # 不跟随链接文件

# EXAMPLE

测试当前磁盘 写入文件 的速度

```bash
dd if=/dev/zero of=testdd bs=1M count=1000
```

测试当前磁盘 纯写入文件 的速度，即不使用缓存

```bash
dd if=/dev/zero of=testdd bs=1M count=1024 oflag=sync,direct,nonblock
```

测试当前磁盘 纯读取文件 的速度，即不使用缓存

```bash
dd if=testdd of=/dev/null bs=1M count=1024 iflag=sync,direct,nonblock
```

测试 sdb 磁盘 的 写入速度。<font color="#ff0000">注意：要使用一块空盘，否则数据没了</font>

```bash
dd if=/dev/urandom of=/dev/sdb1 bs=1M count=1024
```

测试 sdb 磁盘 的读取速度

```bash
dd if=/dev/sdb1 of=/dev/null bs=1M count=1024
```

