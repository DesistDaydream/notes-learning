---
title: 类 Unix 操作系统
---

# 概述

> 参考：
> - [Wiki，Unix](https://en.wikipedia.org/wiki/Unix)
> - [Manual(手册)，os-release](https://man7.org/linux/man-pages/man5/os-release.5.html)
> - [DistroWatch](https://distrowatch.com/)(类 UNIX 系统的资讯网站)

Unix 是一系列多任务、多用户计算机操作系统的统称。

最初打算在[贝尔系统](https://en.wikipedia.org/wiki/Bell_System)内部使用，AT\&T 在 1970 年代后期[将](https://en.wikipedia.org/wiki/License)Unix[授权](https://en.wikipedia.org/wiki/License)给外部各方，导致来自[加利福尼亚大学伯克利分校](https://en.wikipedia.org/wiki/University_of_California,_Berkeley)( [BSD](https://en.wikipedia.org/wiki/Berkeley_Software_Distribution) )、[微软](https://en.wikipedia.org/wiki/Microsoft)( [Xenix](https://en.wikipedia.org/wiki/Xenix) )、[Sun Microsystems 等](https://en.wikipedia.org/wiki/Sun_Microsystems)供应商的各种学术和商业 Unix 变体( [SunOS](https://en.wikipedia.org/wiki/SunOS) / [Solaris](<https://en.wikipedia.org/wiki/Solaris_(operating_system)>) )、[HP](https://en.wikipedia.org/wiki/Hewlett-Packard) / [HPE](https://en.wikipedia.org/wiki/Hewlett_Packard_Enterprise) ( [HP-UX](https://en.wikipedia.org/wiki/HP-UX) ) 和[IBM](https://en.wikipedia.org/wiki/IBM) ( [AIX](https://en.wikipedia.org/wiki/AIX) )。在 1990 年代初期，AT\&T 将其在 Unix 中的权利出售给了[Novell](https://en.wikipedia.org/wiki/Novell)，[Novell](https://en.wikipedia.org/wiki/Novell)随后将其 Unix 业务出售给了[Santa Cruz Operation](https://en.wikipedia.org/wiki/Santa_Cruz_Operation) (SCO) 于 1995 年。[\[4\]](https://en.wikipedia.org/wiki/Unix#cite_note-4) UNIX 商标转让给了[The Open Group](https://en.wikipedia.org/wiki/The_Open_Group)，这是一个成立于 1996 年的行业联盟，允许将该标志用于符合[单一 UNIX 规范](https://en.wikipedia.org/wiki/Single_UNIX_Specification)(SUS) 的认证操作系统。但是，Novell 继续拥有 Unix 版权，[SCO Group, Inc. 诉 Novell, Inc.](https://en.wikipedia.org/wiki/SCO_Group,_Inc._v._Novell,_Inc.)法庭案件 (2010) 证实了这一点。

Unix 系统的特点是[模块化设计](https://en.wikipedia.org/wiki/Modular_design)，有时被称为“ [Unix 哲学](https://en.wikipedia.org/wiki/Unix_philosophy)”。根据这一理念，操作系统应该提供一组简单的工具，每个工具都执行有限的、定义明确的功能。[\[5\]](https://en.wikipedia.org/wiki/Unix#cite_note-5)统一的[文件系统](https://en.wikipedia.org/wiki/Filesystem)（[Unix 文件系统](https://en.wikipedia.org/wiki/Unix_filesystem)）和称为“[管道](<https://en.wikipedia.org/wiki/Pipeline_(Unix)>)”[的进程间通信](https://en.wikipedia.org/wiki/Inter-process_communication)机制作为主要通信手段，[\[3\]](https://en.wikipedia.org/wiki/Unix#cite_note-Ritchie-3)和[shell](<https://en.wikipedia.org/wiki/Shell_(computing)>)脚本和命令语言（[Unix shell](https://en.wikipedia.org/wiki/Unix_shell)）用于结合执行复杂工作流程的工具。

作为第一个[可移植](https://en.wikipedia.org/wiki/Software_portability)操作系统，Unix 与其前辈不同：几乎整个操作系统都是用[C 编程语言](<https://en.wikipedia.org/wiki/C_(programming_language)>)编写的，这使得 Unix 可以在众多平台上运行。

## 类 Unix 操作系统

1983 年，[Richard Stallman](https://en.wikipedia.org/wiki/Richard_Stallman)宣布了[GNU](https://en.wikipedia.org/wiki/GNU)（“GNU's Not Unix”的缩写）项目，这是一项雄心勃勃的努力，旨在创建一个[类似 Unix](https://en.wikipedia.org/wiki/Unix-like)的[自由软件](https://en.wikipedia.org/wiki/Free_software) 系统；“免费”是指每个收到副本的人都可以免费使用、研究、修改和重新分发它。GNU 项目自己的内核开发项目[GNU Hurd](https://en.wikipedia.org/wiki/GNU_Hurd)尚未生产出可运行的内核，但在 1991 年，[Linus Torvalds](https://en.wikipedia.org/wiki/Linus_Torvalds)在[GNU 通用公共许可证](https://en.wikipedia.org/wiki/GNU_General_Public_License)下发布了内核[Linux](https://en.wikipedia.org/wiki/Linux_kernel)作为自由软件。除了在[GNU](https://en.wikipedia.org/wiki/GNU)操作系统中使用之外，许多 GNU 软件包——例如[GNU Compiler Collection](https://en.wikipedia.org/wiki/GNU_Compiler_Collection)（以及其余的[GNU 工具链](https://en.wikipedia.org/wiki/GNU_toolchain)）、[GNU C 库](https://en.wikipedia.org/wiki/Glibc)和[GNU 核心实用程序](https://en.wikipedia.org/wiki/Coreutils) ——也继续在其他自由 Unix 系统中发挥核心作用。

由 Linux 内核和大量兼容软件组成的[Linux 发行版](https://en.wikipedia.org/wiki/Linux_distribution)在个人用户和企业中都很受欢迎。流行的发行版包括 [Red Hat Enterprise Linux](https://en.wikipedia.org/wiki/Red_Hat_Enterprise_Linux)、[Fedora](<https://en.wikipedia.org/wiki/Fedora_(operating_system)>)、[SUSE Linux Enterprise](https://en.wikipedia.org/wiki/SUSE_Linux)、[openSUSE](https://en.wikipedia.org/wiki/OpenSUSE)、[Debian GNU/Linux](https://en.wikipedia.org/wiki/Debian)、[Ubuntu](<https://en.wikipedia.org/wiki/Ubuntu_(operating_system)>)、[Linux Mint](https://en.wikipedia.org/wiki/Linux_Mint)、[Mandriva Linux](https://en.wikipedia.org/wiki/Mandriva_Linux)、[Slackware Linux](https://en.wikipedia.org/wiki/Slackware_Linux)、[Arch Linux](https://en.wikipedia.org/wiki/Arch_Linux)和[Gentoo](https://en.wikipedia.org/wiki/Gentoo_Linux)。

[BSD](https://en.wikipedia.org/wiki/BSD) Unix 的免费衍生产品[386BSD](https://en.wikipedia.org/wiki/386BSD)于 1992 年发布，引发了[NetBSD](https://en.wikipedia.org/wiki/NetBSD)和[FreeBSD](https://en.wikipedia.org/wiki/FreeBSD)项目。1994 年，[Unix 系统实验室](https://en.wikipedia.org/wiki/Unix_System_Laboratories)对加州大学和伯克利软件设计公司（[_USL 诉 BSDi_](https://en.wikipedia.org/wiki/USL_v._BSDi)）提起的诉讼达成和解，澄清了伯克利有权免费分发 BSD Unix，如果它愿意的话。从那时起，BSD Unix 已经在几个不同的产品分支中开发，包括[OpenBSD](https://en.wikipedia.org/wiki/OpenBSD)和[DragonFly BSD](https://en.wikipedia.org/wiki/DragonFly_BSD)。

Linux 和 BSD 越来越多地满足传统上由专有 Unix 操作系统提供服务的市场需求，并扩展到新市场，如消费桌面和移动和嵌入式设备。由于 Unix 模型的模块化设计，共享组件比较常见；因此，大多数或所有 Unix 和类 Unix 系统至少包含一些 BSD 代码，一些系统还在其发行版中包含 GNU 实用程序。

在 1999 年的一次采访中，Dennis Ritchie 表达了他的观点，即 Linux 和 BSD 操作系统是 Unix 设计基础的延续，是 Unix 的衍生物：[\[27\]](https://en.wikipedia.org/wiki/Unix#cite_note-Interview_1999-27)

> 我认为 Linux 现象非常令人愉快，因为它强烈地依赖于 Unix 提供的基础。Linux 似乎是最健康的直接 Unix 衍生产品之一，尽管也有各种 BSD 系统以及来自工作站和大型机制造商的更多官方产品。

在同一次采访中，他表示他认为 Unix 和 Linux 都是“多年前由 Ken 和我以及许多其他人发起的想法的延续”。[\[27\]](https://en.wikipedia.org/wiki/Unix#cite_note-Interview_1999-27)

[OpenSolaris](https://en.wikipedia.org/wiki/OpenSolaris)是[Sun Microsystems](https://en.wikipedia.org/wiki/Sun_Microsystems)开发的[Solaris](<https://en.wikipedia.org/wiki/Solaris_(operating_system)>)的[免费软件](https://en.wikipedia.org/wiki/Free_software)对应物，其中包括[CDDL](https://en.wikipedia.org/wiki/CDDL)许可的内核和主要的[GNU](https://en.wikipedia.org/wiki/GNU)用户空间。然而，[甲骨文](https://en.wikipedia.org/wiki/Oracle_Corporation)在收购 Sun 后停止了该项目，这促使一群前 Sun 员工和 OpenSolaris 社区成员将 OpenSolaris 分叉到[illumos](https://en.wikipedia.org/wiki/Illumos)内核中。截至 2014 年，illumos 仍然是唯一活跃的开源 System V 衍生产品。

# 关联文件

**/etc/os-release** # 操作系统标识。该文件是 /usr/lib/os-release 文件的软链接

注意：

- 不同的 Linux 发行版，都会有一些自身特有的配置，比如 /etc/sysconfg 目录，只会在 RedHat 相关发行版(比如.CentOS)中出现，Ubuntu 并没有这个目录。

## os-release 详解

/etc/os-release 是本身是 [Systemd](/docs/1.操作系统/Systemd/Systemd.md) 的一部分，包含了操作系统的识别数据。在[这篇文章](http://0pointer.de/blog/projects/os-release)里，详解描述了为什么需要这个文件。该文件通常是操作系统供应商定义的，不应该手动修改。

os-release 是一个以换行符分隔的类似环境的 shell 兼容变量赋值列表。示例如下：

```bash
~# cat /etc/os-release
NAME="Ubuntu"
VERSION="20.04.2 LTS (Focal Fossa)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 20.04.2 LTS"
VERSION_ID="20.04"
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
VERSION_CODENAME=focal
UBUNTU_CODENAME=focal
```

这个文件可以很方便得让各种 shell 脚本获取到操作系统的信息，只需要脚本中执行 `source /etc/os-release` 命令即可。

# Ubuntu 与 CentOS 的异同

## 安装镜像

CentOS 的 iso 中包含了所有基础环境所需的软件包，但是 iso 文件过大，7 是 4G 多，8 有 9G 多
Ubuntu 的 iso 中只有一点软件包，虚拟化环境的都没有，但是 iso 文件很小，只有不到 1G

## 网络配置

CentOS 对 NetworkManager 改动较大
Ubuntu 对 NetworkManager 改动几乎没有
