---
title: GNU C Library 管理工具
weight: 8
---

# 概述

CentOS 与 Ubuntu 中关于 [Linux libc 库](docs/1.操作系统/Linux%20源码解析/Linux%20libc%20库.md)(glibc) 的管理工具包的名称不太一样

- CentOS 为 glibc-common
- Ubuntu 为 libc-bin、libc-dev-bin

# ldd

这个命令可以显示一个可执行文件所使用的动态链接库。如：

```bash
~]# ldd /usr/bin/ls
 linux-vdso.so.1 (0x00007ffd37562000)
 libselinux.so.1 => /lib/x86_64-linux-gnu/libselinux.so.1 (0x00007f76b6c46000)
 libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007f76b6a54000)
 libpcre2-8.so.0 => /lib/x86_64-linux-gnu/libpcre2-8.so.0 (0x00007f76b69c4000)
 libdl.so.2 => /lib/x86_64-linux-gnu/libdl.so.2 (0x00007f76b69be000)
 /lib64/ld-linux-x86-64.so.2 (0x00007f76b6ca2000)
 libpthread.so.0 => /lib/x86_64-linux-gnu/libpthread.so.0 (0x00007f76b699b000)
```

# ldconfig

lbconfig 是一个动态链接库的管理命令。可以创建必要的链接并缓存到指定文件中

## 关联文件与配置

**/lib64/ld-linux.so** # 运行时链接器/加载器

**/etc/ld.so.conf** # 从配置文件中指定的目录中搜索库

**/etc/ld.so.conf.d/** # 从配置文件中指定的目录中搜索库

**/etc/ld.so.cache** # 搜索到的库列表被缓存到该文件中

### Syntax(语法)

OPTIONS

- **-p, --print-cache** # 输出当前系统已经加载的动态库
