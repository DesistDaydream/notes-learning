---
title: X.Linux 管理
---

# 概述

> - [GNU Manual(手册)](https://www.gnu.org/manual/) — Linux 中很多核心程序，都是 GNU 组织下的软件。

系统管理员可以通过 一系列用户空间的二进制应用程序来管理 Linux 操作系统。Linux 内核自带了一个名为 coreutils 包，包含了很多最基本的管理工具。

除了 Coreutils 包，还有很多很多的应用程序，一起组成了一套工具栈，系统管理员可以根据自身的需求，有选择得安装并使用它们。

# Coreutils

> 参考：
> - [Wiki,GNU Core Utilies](https://en.wikipedia.org/wiki/GNU_Core_Utilities)
> - [官方文档](https://www.gnu.org/software/coreutils/manual/)

GNU Core Utilities 是 GNU 操作系统的基本文件、Shell、文本操作的实用程序。同时，也是现在绝大部分 Linux 发行版内置的实用程序。

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

这个包中，通常包含如下应用程序

```bash
arch base64 basename cat chcon chgrp chmod chown chroot cksum comm cp  csplit cut date dd df dir dircolors dirname du echo env expand expr  factor false flock fmt fold groups head hostid id install join link ln  logname ls md5sum mkdir mkfifo mknod mktemp mv nice nl nohup nproc numfmt  od paste pathchk pinky pr printenv printf ptx pwd readlink realpath rm  rmdir runcon sha*sum seq shred sleep sort split stat stty sum sync tac  tail tee test timeout touch tr true truncate tsort tty uname unexpand  uniq unlink users vdir wc who whoami yes
```

可以发现，这些命令就是我们日常经常使用那些~

# Util-linux

> 参考：
> - [GitHub 项目，util-linux/util-linux](https://github.com/util-linux/util-linux)
> - [Wiki,Util-linux](https://en.wikipedia.org/wiki/Util-linux)

util-linux 是由 Linux 内核组织分发的标准软件包，用作 Linux 操作系统的一部分。一个分支 util-linux-ng（ng 的意思是“下一代”）是在开发停滞时创建的，但截至 2011 年 1 月，它已重命名为 util-linux，并且是该软件包的正式版本。

Util-linux 包中通常包含如下程序：

- addpart
- [agetty](<https://en.wikipedia.org/wiki/Getty_(Unix)>)
- blkdiscard
- [blkid](https://en.wikipedia.org/wiki/Blkid)
- blkzone
- [blockdev](https://en.wikipedia.org/w/index.php?title=Blockdev&action=edit&redlink=1)
- [cal](<https://en.wikipedia.org/wiki/Cal_(command)>)
- [cfdisk](https://en.wikipedia.org/wiki/Cfdisk)
- chcpu
- chfn
- chmem
- choom
- chrt
- [chsh](https://en.wikipedia.org/wiki/Chsh)
- col (legacy)
- colcrt
- colrm
- column
- ctrlaltdel
- delpart
- [dmesg](https://en.wikipedia.org/wiki/Dmesg)
- eject
- fallocate
- [fdformat](https://en.wikipedia.org/wiki/Fdformat)
- [fdisk](https://en.wikipedia.org/wiki/Fdisk)
- fincore
- findfs
- findmnt
- [flock](https://en.wikipedia.org/wiki/File_locking)
- [fsck](https://en.wikipedia.org/wiki/Fsck)
- [fsck](https://en.wikipedia.org/wiki/Fsck).[cramfs](https://en.wikipedia.org/wiki/Cramfs)
- [fsck](https://en.wikipedia.org/wiki/Fsck).[minix](https://en.wikipedia.org/wiki/MINIX_file_system)
- fsfreeze
- fstrim
- [getopt](https://en.wikipedia.org/wiki/Getopt)
- hardlink
- [hexdump](https://en.wikipedia.org/wiki/Hex_dump#od_and_hexdump)
- [hwclock](https://en.wikipedia.org/w/index.php?title=Hwclock&action=edit&redlink=1) (query and set the hardware clock (RTC))
- [ionice](https://en.wikipedia.org/wiki/Ionice)
- ipcmk
- [ipcrm](https://en.wikipedia.org/wiki/Ipcrm)
- [ipcs](https://en.wikipedia.org/wiki/Ipcs)
- irqtop
- isosize
- [kill](<https://en.wikipedia.org/wiki/Kill_(Unix)>)
- last
- ldattach
- line (legacy)
- logger
- login
- look
- [losetup](https://en.wikipedia.org/wiki/Losetup)
- lsblk
- lscpu[\[8\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-8)
- lsfd
- lsipc
- lsirq[\[9\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-9)
- lslocks
- lslogins
- lsmem
- lsns
- mcookie
- [mesg](https://en.wikipedia.org/wiki/Mesg)
- [mkfs](https://en.wikipedia.org/wiki/Mkfs) (legacy)
- [mkfs](https://en.wikipedia.org/wiki/Mkfs).[bfs](https://en.wikipedia.org/wiki/Boot_File_System)
- [mkfs](https://en.wikipedia.org/wiki/Mkfs).[cramfs](https://en.wikipedia.org/wiki/Cramfs)
- [mkfs](https://en.wikipedia.org/wiki/Mkfs).[minix](https://en.wikipedia.org/wiki/MINIX_file_system)
- mkswap
- [more](<https://en.wikipedia.org/wiki/More_(command)>)
- [mount](<https://en.wikipedia.org/wiki/Mount_(Unix)>)
- mountpoint
- namei
- newgrp
- nologin
- nsenter
- partx
- [pg](<https://en.wikipedia.org/wiki/Pg_(Unix)>) (legacy)
- pivot_root
- prlimit[\[10\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-10)
- raw
- readprofile
- rename
- [renice](<https://en.wikipedia.org/wiki/Nice_(Unix)>)
- reset (legacy)
- resizepart
- rev
- rfkill
- [rtcwake](https://en.wikipedia.org/wiki/RTC_Alarm)
- runuser
- [script](<https://en.wikipedia.org/wiki/Script_(Unix)>)
- scriptlive
- scriptreplay
- setarch (including architecture symlinks such as i386, linux32, linux64, x86_64, etc.)
- setpriv
- setsid
- setterm
- [sfdisk](https://en.wikipedia.org/wiki/Sfdisk)
- su
- sulogin
- swaplabel
- swapoff
- swapon
- switch_root
- taskset
- tunelp (deprecated)[\[11\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-11)
- ul
- umount
- unshare
- utmpdump
- uuidd
- uuidgen
- uuidparse
- [vipw](https://en.wikipedia.org/wiki/Vipw) (including symlink to vigr)
- wall
- wdctl
- [whereis](https://en.wikipedia.org/wiki/Whereis)
- wipefs
- [write](<https://en.wikipedia.org/wiki/Write_(Unix)>)
- zramctl

### Removed

Utilities formerly included, but removed as of 1 July 2015:

- arch[\[12\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-12)
- chkdupexe[\[13\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-13)
- clock[\[14\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-14)
- cytune[\[15\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-15)
- [ddate](https://en.wikipedia.org/wiki/Ddate) (removed from default build[\[16\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-16) before being removed[\[17\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-17) altogether)
- elvtune[\[18\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-18)
- fastboot[\[19\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-simpleinit-19)
- fasthalt[\[19\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-simpleinit-19)
- halt[\[19\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-simpleinit-19)
- initctl[\[19\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-simpleinit-19)
- ramsize (formerly a symlink to rdev)[\[20\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-rdev-20)
- rdev[\[20\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-rdev-20)
- reboot[\[19\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-simpleinit-19)
- rootflags (formerly a symlink to rdev)[\[20\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-rdev-20)
- [shutdown](<https://en.wikipedia.org/wiki/Shutdown_(command)>)[\[19\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-simpleinit-19)
- simpleinit[\[19\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-simpleinit-19)
- tailf[\[21\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-tailf-21)
- vidmode (formerly a symlink to rdev)[\[20\]](https://en.wikipedia.org/wiki/Util-linux#cite_note-rdev-20)
