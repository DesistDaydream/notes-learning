---
title: 1.Linux Kernel
---

# 概述

> 参考：
> - [Linus Torvalds](https://github.com/torvalds)
> - [GitHub,Linux 内核项目](https://github.com/torvalds/linux)
> - [官网](https://www.kernel.org/)
> - [官方文档](https://www.kernel.org/doc/html/latest/)
> - [官方 Manual(手册)](https://www.kernel.org/doc/man-pages/index.html)
> - [Wiki,Kernel](<https://en.wikipedia.org/wiki/Kernel_(operating_system)>)
> - [Wiki,/boot](https://en.wikipedia.org/wiki//boot/)
> - [Wiki,vmlinux](https://en.wikipedia.org/wiki/Vmlinux)
> - [Wiki,Initial ramdisk](https://en.wikipedia.org/wiki/Initial_ramdisk)
> - [Wiki,System.map](https://en.wikipedia.org/wiki/System.map)
> - [树莓派 Linux](https://github.com/raspberrypi/linux)
> - [RedHat 官方文档,8-管理、监控和更新内核](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html/managing_monitoring_and_updating_the_kernel/index)
> - <http://www.linfo.org/vmlinuz.html>
> - [知乎,initrd 和 initramfs 的区别](https://www.zhihu.com/question/22045825)

**Kernel(内核) **是一个作为操作系统核心的计算机程序，对系统中的一切具有完全控制权。它负责管理系统的进程、内存、设备驱动程序、文件和网络系统，决定着系统的性能和稳定性。

Kernel 是计算器启动时首先加载程序之一(在 [Bootloader](/docs/IT学习笔记/1.操作系统/1.Bootloader/1.Bootloader.md)并处理硬件和软件之间的交互。并且处理启动过程的其余部分、以及内存、外设、和来自软件的输入/输出请求，将他们转换为 CPU 的数据处理指令。

## Kernel 组成 及 系统调用

Linux 内核由如下几部分组成：内存管理、进程管理、设备驱动程序管理、文件系统管理、网络管理等。如图：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fkp6xi/1616168349819-c21dd43c-79b7-4ec2-abd4-c8bb0e3c7686.jpeg)
**System Call Interface(系统调用接口，简称 SCI) **层提供了某些机制执行从用户空间到内核的函数调用。这个接口依赖于体系结构，甚至在相同的处理器家族内也是如此。SCI 实际上是一个非常有用的函数调用多路复用和多路分解服务。

系统调用介绍详见[ system call(系统调用)](/docs/IT学习笔记/1.操作系统/2.Kernel(内核)/3.System%20Call(系统调用).md Call(系统调用).md) 章节

## Linux man 手册使用说明

在 Linux Kernel 的官方 man 手册中，记录了用户空间程序使用 Linux 内核 和 C 库的接口。对于 C 库，主要聚焦于 GUN C(glibc)，尽管在已知的情况下，还包括可用于 Linux 的其他 C 库中的变体文档。在这个 man 手册中，分为如下几部分

- **[User commands](https://man7.org/linux/man-pages/dir_section_1.html)(用户命令)** # 介绍一些用户空间的应用程序。
- **[System calls](https://man7.org/linux/man-pages/dir_section_2.html)(系统调用)** # Linux Kernel 可以提供的所有 System Calls(系统调用)
- **[Library functions](https://man7.org/linux/man-pages/dir_section_3.html)(库函数)** # C 标准库可以提供的函数。
- **[Devices](https://man7.org/linux/man-pages/dir_section_4.html)(设备)** # 各种设备的详细信息，这些设备大多都在 /dev 目录中。
- **[Files](https://man7.org/linux/man-pages/dir_section_5.html)(文件)** # 各种文件格式和文件系统
- **[Overviews, conventions, and miscellaneous](https://man7.org/linux/man-pages/dir_section_7.html)(概述、约定 和 其他)** #
- **[Superuser and system administration commands](https://man7.org/linux/man-pages/dir_section_8.html)(超级用户和系统管理员命令)** # 介绍一些 GUN C 库提供的程序。

在 Linux man 手册中，可以找到 Linux 系统中的一切使用说明。Linux 操作系统围绕 Linux Kernel 构建了一套高效、健壮的应用程序运行环境

# intird.img、vmlinuz、System.map 文件

Kernel 会被安装到 /boot 目录中，并生成 **config、initrd.img、System.map、vmlinuz** 这几个文件

## vmlinuz

\_vmlinuz \_是 [Linux](http://www.linfo.org/linuxdef.html) [_内核_](http://www.linfo.org/kernel.html) \_可执行文件\_的名称。

vmlinuz 是一个压缩的 Linux 内核，它是\_可引导的\_。可引导意味着它能够将操作系统加载到内存中，以便计算机变得可用并且可以运行应用程序。

vmlinuz 不应与\_vmlinux\_混淆，后者是非压缩和不可引导形式的内核。vmlinux 通常只是生成 vmlinuz 的中间步骤。

vmlinuz 位于\_/boot\_目录中，该目录包含开始引导系统所需的文件。名为\_vmlinuz\_的文件可能是实际的内核可执行文件本身，也可能是内核可执行文件的链接，该链接可能带有诸如\_/boot/vmlinuz-2.4.18-19.8.0\_之类的名称（即特定内核的名称）内核版本）。这可以通过使用\_ls\_ [命令](http://www.linfo.org/command.html)（其目的是列出指定目录的内容）及其\_-l\_选项（它告诉 ls 提供有关指定目录中每个对象的详细信息）来轻松确定，如下所示：

> `ls -l /boot`

如果 vmlinuz 是一个普通文件（包括可执行文件），则第一列中有关它的信息将以连字符开头。如果是链接，它将以字母\_l\_开头。
通过发出以下命令   来\_编译\_Linux 内核：

> `make bzImage`

这会在 \_/usr/src/linux/arch/i386/linux/boot/ \_等目录中创建名为\_bzImage\_的文件。

编译是将内核的[_源代码_](http://www.linfo.org/source_code.html)（即内核由人类编写的原始形式）转换为\_目标代码\_（计算机处理器可以直接理解）。它由称为[_编译器_](http://www.linfo.org/compiler.html)的专门程序执行，通常是[_GCC_](http://www.linfo.org/gcc.html)（[GNU](http://www.linfo.org/gnu.html)编译器集合）中的一个。

然后使用 \_cp _命令将 bzImage 复制到 /boot 目录，同时使用诸如以下命令   重命名\_vmlinuz_

> `cp /usr/src/linux/arch/i386/linux/boot/bzImage /boot/vmlinuz`

vmlinuz 不仅仅是一个压缩图像。它还内置了\_gzip\_解压缩器代码。gzip 是[类 Unix](http://www.linfo.org/unix-like.html)操作系统上最流行的压缩实用程序之一。

一个名为\_zImage\_文件的编译内核是在一些较旧的系统上创建的，并保留在较新的系统上以实现向后兼容性。zImage 和 bzImage 都是用 gzip 压缩的。区别在于 zImage 解压到\_低内存\_（即前 640kB），bzImage 解压到\_高内存\_（1MB 以上）。有一个常见的误解，认为 bzImage 是使用\_bzip2\_实用程序压缩的。实际上，_b\_只代表\_big_。

_vmlinuz_ 这个名字很大程度上是历史的偶然。在贝尔实验室开发的原始 UNIX 上的内核二进制文件称为\_unix\_。当后来在加州大学伯克利分校 (UCB) 编写包含支持[_虚拟内存_](http://www.linfo.org/virtual_memory.html)的新内核时，内核二进制文件更名为\_vmunix\_。

虚拟内存是使用硬盘驱动器 (HDD) 上的空间来模拟额外的 RAM（随机存取内存）容量。与当时使用的其他一些流行操作系统（例如[MS-DOS）](http://www.linfo.org/ms-dos.html)相比，Linux 内核几乎从 Linux 一开始就支持它。

因此，Linux 内核很自然地被称为\_vmlinux\_。由于 Linux 内核可执行文件被制作成压缩文件，并且压缩文件在类 Unix 系统上通常具有\_z\_或\_gz\_扩展名，因此压缩内核可执行文件的名称变为\_vmlinuz\_。

## initrd

**Initial RAM Disk(初始内存磁盘，简称 initrd)** 是一种将临时根文件系统加载到内存中的方案，可以作为 Linux 启动过程的一部分。有两种方法来实现这种方案：

- **initrd # Initial RAM Disk。**就是把一块内存（ram）当做磁盘（disk）去挂载，然后找到 ram 里的 init 执行。
- **initramfs # Initial RAM Filesystem。**直接在 ram 上挂载文件系统，执行文件系统中的 init。

这两者通常用于在挂载真正的根文件系统之前执行一些准备工作。

> 不要被文件名迷惑，kernel 2.6 以来都是 initramfs 了，只是很多还沿袭传统使用 initrd 的名字
> initramfs 的工作方式更加简单直接一些，启动的时候加载内核和 initramfs 到内存执行，内核初始化之后，切换到用户态执行 initramfs 的程序/脚本，加载需要的驱动模块、必要配置等，然后加载 rootfs 切换到真正的 rootfs 上去执行后续的 init 过程。

> initrd 是 2.4 及更早的用法（现在你能见到的 initrd 文件实际差不多都是 initramfs 了），运行过程大概是内核启动，执行一些 initrd 的内容，加载模块啥的，然后交回控制权给内核，最后再切到用户态去运行用户态的启动流程。

> 从格式看，老的 initrd 是一个压缩的内存文件系统，具体是啥忘了，年月太久了。现在的 initramfs 是一个 gzip 压缩的 cpio 文件系统打包，如果遇到什么紧急情况需要处理的时候，你可以建立一个临时目录，把 initramfs 解压之后，直接 cpio -idv 解压出来，改之后再用 cpio 和 gzip 封上即可。虽然大家都喜欢用 tar 打包，但掌握点 cpio 在关键时刻还是可以救命的。

在早期的 Linux 系统中，一般就只有软盘或者硬盘被用来作为 Linux 的根文件系统，因此很容易把这些设备的驱动程序集成到内核中。但是现在根文件系统可能保存在各种存储设备上，包括 SCSI、SATA、U 盘等等。总不能每出一个，就要重新编译一遍内核吧？~这样不但麻烦，也不实用，所以后来 Linux 就提供了一个灵活的方法来解决这些问题。就是 **initrd**。

> 可以把 initrd 当做 WinPE。当使用 WinPE 启动后会发现你的计算机就算没有硬盘也能在正常运行，其中有个文件系统 B:/ 分区，这个分区就是内存模拟的磁盘。

initrd.img 文件就是个 ram disk 的映像文件。ramdisk 是用一部分内存模拟成磁盘，让操作系统访问。ram disk 是标准内核文件认识的设备(/dev/ram0)文件系统也是标准内核认识的文件系统。内核加载这个 ram disk 作为根文件系统并开始执行其中的"某个文件"（2.6 内核是 init 文件）来加载各种模块，服务等。经过一些配置和运行后，就可以去物理磁盘加载真正的 root 分区了，然后又是一些配置等，最后启动成功。

也就是你只需要定制适合自己的 initrd.img 文件就可以了。这要比重编内核简单多了，省时省事低风险。

### 查看 initrd 文件

我们可以通过如下方式，解压出 initrd.img 文件，下面分别以 Ubuntu 20.04 TLS 系统和 CentOS Stream 8 系统为例：

> 解压方法来源：<https://unix.stackexchange.com/questions/163346/why-is-it-that-my-initrd-only-has-one-directory-namely-kernel>

Ubuntu 20.04 TLS

```bash
root@lichenhao:~/test_dir# mkdir -p /root/test_dir/root
root@lichenhao:~/test_dir# cp /boot/initrd.img /root/test_dir/
root@lichenhao:~/test_dir/root# (cpio -id; cpio -i; unlz4 | cpio -id) < ../initrd.img
62 blocks
9004 blocks
450060 blocks
root@lichenhao:~/test_dir/root# ls
bin  conf  cryptroot  etc  init  kernel  lib  lib32  lib64  libx32  run  sbin  scripts  usr  var
```

CentOS Stream 8

```bash
[root@master-2 root]# mkdir -p /root/test_dir/root
[root@master-2 root]# cp /boot/initramfs-4.18.0-294.el8.x86_64.img /root/test_dir
[root@master-2 root]# cd /root/test_dir/root
[root@master-2 root]# (cpio -id; zcat | cpio -id) < ../initramfs-4.18.0-294.el8.x86_64.img

[root@master-2 root]# ls
bin  dev  early_cpio  etc  init  kernel  lib  lib64  proc  root  run  sbin  shutdown  sys  sysroot  tmp  usr  var
```

可以看到，initrd.img 中包含了一个系统最基本的目录结构

# Kernel 关联文件

**/boot/\* **#

- **./config-$(uname -r) **# Kernel 的扩展配置文件。Kernel 文档中，将该文件称为 **Boot Configuration**。
- **./initrd.img** # 在内核挂载真正的根文件系统前使用的临时文件系统
- **./vmlinuz** # Linux 内核

**/etc/sysctl.conf **# 系统启动时读取的内核参数文件

- **/etc/sysctl.d/\* **# 系统启动时时读取的内核参数目录

**/usr/lib/sysctl.d/\* **#
**/proc/sys/\* **# 内核参数(也称为内核变量)所在路径。该目录(从 1.3.57 版本开始)包含许多与内核变量相对应的文件和子目录。 可以使用 [proc 文件系统](https://www.yuque.com/go/doc/33222789) 以及 sysctl(2) 系统读取或加载这些变量，有时可以对其进行修改。
