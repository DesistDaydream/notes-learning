---
title: "Linux 管理"
linkTitle: "Linux 管理"
weight: 1
---

# 概述

> 参考：
>
> - [GNU Manual(手册)](https://www.gnu.org/manual/) — Linux 中很多核心程序，都是 GNU 组织下的软件。

系统管理员可以通过 一系列用户空间的二进制应用程序来管理 Linux 操作系统。Linux 内核自带了一个名为 coreutils 包，包含了很多最基本的管理工具。

除了 Coreutils 包，还有很多很多的应用程序，一起组成了一套工具栈，系统管理员可以根据自身的需求，有选择得安装并使用它们。

# Coreutils

> 参考：
>
> - [Wiki，GNU Core Utilies](https://en.wikipedia.org/wiki/GNU_Core_Utilities)
> - [官方文档](https://www.gnu.org/software/coreutils/manual/)

GNU Core Utilities 是 GNU/Linux 操作系统的基本文件、Shell、文本操作的实用程序。同时，也是现在绝大部分 Linux 发行版内置的实用程序。

Coreutils 通常可以通过各种 Linux 发行版的包管理器直接安装。

```bash
root@lichenhao:~/downloads# apt-cache show coreutils
Package: coreutils
Architecture: amd64
Version: 8.30-3ubuntu2
Multi-Arch: foreign
Priority: required
Essential: yes
Section: utils
Origin: Ubuntu
Maintainer: Ubuntu Developers <ubuntu-devel-discuss@lists.ubuntu.com>
Original-Maintainer: Michael Stone <mstone@debian.org>
Bugs: https://bugs.launchpad.net/ubuntu/+filebug
Installed-Size: 7196
Pre-Depends: libacl1 (>= 2.2.23), libattr1 (>= 1:2.4.44), libc6 (>= 2.28), libselinux1 (>= 2.1.13)
Filename: pool/main/c/coreutils/coreutils_8.30-3ubuntu2_amd64.deb
Size: 1249368
MD5sum: e8e201b6d1b7f39776da07f6713e1675
SHA1: 1d4ab60c729a361d46a90d92defaca518b2918d2
SHA256: 99aa50af84de1737735f2f51e570d60f5842aa1d4a3129527906e7ffda368853
Homepage: http://gnu.org/software/coreutils
Description-en: GNU core utilities
 This package contains the basic file, shell and text manipulation
 utilities which are expected to exist on every operating system.
 .
 Specifically, this package includes:
 arch base64 basename cat chcon chgrp chmod chown chroot cksum comm cp
 csplit cut date dd df dir dircolors dirname du echo env expand expr
 factor false flock fmt fold groups head hostid id install join link ln
 logname ls md5sum mkdir mkfifo mknod mktemp mv nice nl nohup nproc numfmt
 od paste pathchk pinky pr printenv printf ptx pwd readlink realpath rm
 rmdir runcon sha*sum seq shred sleep sort split stat stty sum sync tac
 tail tee test timeout touch tr true truncate tsort tty uname unexpand
 uniq unlink users vdir wc who whoami yes
Description-md5: d0d975dec3625409d24be1238cede238
Task: minimal
```

这个包中包含的程序可以在上面的 Specifically, this package includes 部分看到。

这些命令就是我们日常经常使用那些~

# Util-linux

详见 [Util-linux Utilities](docs/1.操作系统/Linux%20管理/Util-linux%20Utilities.md)
