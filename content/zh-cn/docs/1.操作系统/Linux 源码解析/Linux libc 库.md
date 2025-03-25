---
title: Linux libc 库
linkTitle: Linux libc 库
weight: 2
---

# 概述

> 参考：
>
> - [Manual(手册)，libc(7)](https://man7.org/linux/man-pages/man7/libc.7.html)
> - [Wiki, glibc](https://en.wikipedia.org/wiki/Glibc)

术语 `libc` 通常用作“标准 C 库”的简写，这是所有 C 程序（以及其他语言的程序）可以使用的标准函数库。由于一些历史（见下文），使用术语 `libc` 来指代标准 C 库在 Linux 上是不够严谨的。

## glibc

到目前为止，Linux 上使用最广泛的 C 库是 GNU C 库。这是目前在所有主要 Linux 发行版中使用的 C 库。它也是 C 库，其详细信息记录在手册页项目的相关页面中 (主要在手册的第 3 节中)。glibc 的文档也可在 glibc 手册中找到，可通过命令信息 libc 获得。glibc 的版本 1.0 制作于 1992 年 9 月。(之前有 0.x 版本。)glibc 的下一个主要版本是 1997 年初的 2.0。

路径名 **/lib/libc.so.6** (或类似内容) 通常是指向 glibc 库位置的符号链接，执行此路径名将导致 glibc 显示有关系统上安装的版本的各种信息。

CentOS7：

```bash
~]# ls -l /lib64/libc.so.6
lrwxrwxrwx. 1 root root 19 Mar  7 16:53 /lib64/libc.so.6 -> /lib64/libc-2.17.so
```

Ubuntu20.04：

```bash
~]# ls -l /lib/x86_64-linux-gnu/libc.so.6
lrwxrwxrwx 1 root root 12 Dec 16  2020 /lib/x86_64-linux-gnu/libc.so.6 -> libc-2.31.so
```

## Linux libc

在 20 世纪 90 年代初期至中期，有一段时间 Linux libc，这是由 Linux 开发人员创建的 glibc 1.x 分支，他们认为当时的 glibc 开发不足以满足 Linux 的需求。通常，该库被 (模棱两可地) 称为 “libc”。Linux libc 发布了主要版本 2、3、4 和 5，以及这些版本的许多次要版本。Linux libc4 是最后一个使用 a.out 二进制格式的版本，也是第一个提供 (原始) 共享库支持的版本。Linux libc 5 是第一个支持 ELF 二进制格式的版本; 该版本使用共享库 soname libc.So.5。有一段时间，Linux libc 是许多 Linux 发行版中的标准 C 库。

然而，尽管 Linux libc 努力的最初动机，到 glibc 2.0 发布时 (1997 年)，它显然优于 Linux libc，并且所有使用 Linux libc 的主要 Linux 发行版很快又切换回 glibc。为了避免与 Linux libc 版本混淆，glibc 2.0 和更高版本使用了共享库 soname **/PATH/TO/libc.so.6**。

由于从 Linux libc 到 glibc 2.0 的切换发生在很久以前，手册页不再需要记录 Linux libc 的详细信息。尽管如此，有关 Linux libc 的信息的痕迹仍可见，这些信息保留在一些手册页面中，特别是对 libc4 和 libc5 的引用

由于上述原因，CentOS 与 Ubuntu 系统的 glibc 包的名字不太一样~

- CentOS 中叫 glibc，glibc 工具在 glibc-common 包中。
- Ubuntu 中叫 libc6，glibc 工具在 libc-bin、libc-dev-bin 包中。

**Ubuntu 中 Linux libc**

```bash
~]# apt show libc6
Package: libc6
Version: 2.31-0ubuntu9.2
Priority: required
Section: libs
Source: glibc
Origin: Ubuntu
Maintainer: Ubuntu Developers <ubuntu-devel-discuss@lists.ubuntu.com>
Original-Maintainer: GNU Libc Maintainers <debian-glibc@lists.debian.org>
Bugs: https://bugs.launchpad.net/ubuntu/+filebug
Installed-Size: 13.6 MB
Depends: libgcc-s1, libcrypt1 (>= 1:4.4.10-10ubuntu4)
Recommends: libidn2-0 (>= 2.0.5~)
Suggests: glibc-doc, debconf | debconf-2.0, locales
Conflicts: openrc (< 0.27-2~)
Breaks: hurd (< 1:0.9.git20170910-1), iraf-fitsutil (< 2018.07.06-4), libtirpc1 (< 0.2.3), locales (< 2.31), locales-all (< 2.31), nocache (< 1.1-1~), nscd (< 2.31), r-cran-later (< 0.7.5+dfsg-2), wcc (< 0.0.2+dfsg-3)
Replaces: libc6-amd64
Homepage: https://www.gnu.org/software/libc/libc.html
Task: minimal
Original-Vcs-Browser: https://salsa.debian.org/glibc-team/glibc
Original-Vcs-Git: https://salsa.debian.org/glibc-team/glibc.git
Download-Size: 2,715 kB
APT-Manual-Installed: no
APT-Sources: http://mirrors.aliyun.com/ubuntu focal-updates/main amd64 Packages
Description: GNU C Library: Shared libraries
 Contains the standard libraries that are used by nearly all programs on
 the system. This package includes shared versions of the standard C library
 and the standard math library, as well as many others.
```

**CentOS 中的 Linux libc**

```bash
~]# yum info glibc
Installed Packages
Name        : glibc
Arch        : x86_64
Version     : 2.17
Release     : 324.el7_9
Size        : 13 M
Repo        : installed
From repo   : updates
Summary     : The GNU libc libraries
URL         : http://www.gnu.org/software/glibc/
License     : LGPLv2+ and LGPLv2+ with exceptions and GPLv2+
Description : The glibc package contains standard libraries which are used by
            : multiple programs on the system. In order to save disk space and
            : memory, as well as to make upgrading easier, common system code is
            : kept in one place and shared between programs. This particular package
            : contains the most important sets of shared libraries: the standard C
            : library and the standard math library. Without these two libraries, a
            : Linux system will not function.

```

虽然这两个包中关于 URL 的描述不同，其实访问 <http://www.gnu.org/software/glibc/> 会重定向到 <https://www.gnu.org/software/libc/libc.html>

## 其他 C 库

还有其他各种不太广泛用于 Linux 的 C 库。就功能和内存占用而言，这些库通常比 glibc 小，并且通常用于构建小型二进制文件，可能针对嵌入式 Linux 系统的开发。在这些库中，有 uClibc 的 http:www.uclibc.org/ase，dietlibc 的 http: www.fefe.de/dietlibc/ 该库和 musl libc 的 http:www.uclibc.org/musl-libc.org。这些库的详细信息由已知的手册项目涵盖。

## 非常重要的 libc.so.6 文件

Linux 操作系统通过一个 libc.so.6 这个软链接文件，指向本机的 Libc 文件，并且，系统中绝大部分指令都要依靠这个软链接才能执行。想要升级是比较麻烦的，因为它可以说是整个系统的基础，所以在系统安装好后基本上是不会动这个文件的。但是如果要升级，就一定会修改这个 libc.so.6 所链接的文件，可以使用 `ln -sf` 这样直接覆盖。

有些人习惯先删除再重建，但是一旦这个链接和原来的 Glibc 文件断开，在 Shell 中除了便几乎无法执行指令。连 ls、mv、cp 都用不了的，更不用说 ln、ssh 等命令了，报错如下：

```shell
~]# ls
/bin/ls: error while loading shared libraries: libc.so.6: cannot open shared object file: No such file or directory
```

此时如果想要修复，可以通过 `${LD_PRELOAD}` 变量，在执行命令前添加 `LD_PRELOAD=<glibc 实际文件绝对路径>`，例如这样：

```shell
~]# LD_PRELOAD=/lib/x86_64-linux-gnu/libc-2.17.so ls -al
```

像这样手动加载好动态链接库，就可以正常运行命令了。
