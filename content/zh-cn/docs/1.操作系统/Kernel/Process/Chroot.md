---
title: chroot
linkTitle: chroot
date: 2024-03-16T18:59
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, chroot](https://en.wikipedia.org/wiki/Chroot)
> - [Manual(手册)，chroot(2)](https://man7.org/linux/man-pages/man2/chroot.2.html)

**Change root(改变根，简称 Chroot)** 是 [Unix-like OS](/docs/1.操作系统/Operating%20system/Unix-like%20OS/Unix-like%20OS.md) 的一种操作，用于更改对当前正在运行的进程及其子进程展现出来的 `/` 目录。在这种修改过的环境中运行的程序无法访问指定目录之外的文件。

Chroot 的意思是改变根路径的位置(Linux 系统中以 `/` 为根目录位置，但是对于执行 Chroot 的用户或者程序来说，是 Chroot 后 PATH 的位置是新的根目录位置)，比如 Telnet，ssh，如果都定义了 Chroot(PATH)规则，那么远程登录的用户将无法访问到该 linux 系统中除了定义的 PATH 外的其余目录

```bash
]# pwd
/var/lib/docker/overlay2/72a3b770bf98493a90e2e335adbdc9f92eeb18f19044136f74c5c9138cb13304/merged
]# ls
bin  dev  etc  home  lib  LICENSE  NOTICE  npm_licenses.tar.bz2  proc  prometheus  root  sys  tmp  usr  var
]# ls /root
backup  downloads  go  nohup.out  p.pcap  projects  scripts  snap  tmp
]# chroot . /bin/sh
/ # pwd
/
/ # ls /root
/ #
```

上面例子中，我们通过 chroot 程序进入了以  `/var/lib/docker/overlay2/72a3b770bf98493a90e2e335adbdc9f92eeb18f19044136f74c5c9138cb13304/merged/` 目录作为 `/` 目录的空间中。这就像是将本地文件系统划分了一块空间给 Chroot 后的使用者。

Chroot 功能主要依赖于 chroot [System Call](/docs/1.操作系统/Kernel/System%20Call/System%20Call.md) 实现

# CLI

## chroot

### Syntax(语法)

> 参考：
>
> - [Manual(手册)，chroot(1)](https://man7.org/linux/man-pages/man1/chroot.1.html)

**chroot \[OPTION] NEWROOT \[COMMAND \[ARG]...]**

COMMAND 是指这次执行命令过程中，将要进入的目标后执行的命令，通常是 `$SHELL`，也就是当前环境所使用的 SHELL 的命令。

**OPTIONS**

- **--userspec=USER:GROUP** # 指定要使用的用户和组的 ID 或 NAME。
