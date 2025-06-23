---
title: Debian 与 Ubuntu
---

# 概述

> 参考：
>
> - [Debian 官方 Manual(手册)](https://manpages.debian.org/)

Debian 与 Ubuntu 是 [Unix-like OS](/docs/1.操作系统/Operating%20system/Unix-like%20OS/Unix-like%20OS.md) 发行版

```bash
groupadd wheel
usermod -G wheel desistdaydream
tee /etc/sudoers.d/desistdaydream > /dev/null <<EOF
%wheel        ALL=(ALL)       NOPASSWD: ALL
EOF
```

\~/.bashrc

```bash
if [ "$color_prompt" = yes ]; then
    # PS1='${debian_chroot:+($debian_chroot)}\[\033[01;32m\]\u@\h\[\033[00m\]:\[\033[01;34m\]\w\[\033[00m\]\$ '
    PS1='${debian_chroot:+($debian_chroot)}[\[\e[34;1m\]\u@\[\e[0m\]\[\e[32;1m\]\H\[\e[0m\] \[\e[31;1m\]\w\[\e[0m\]]\\$ '
else
    # PS1='${debian_chroot:+($debian_chroot)}\u@\h:\w\$ '
    PS1='${debian_chroot:+($debian_chroot)}[\[\e[34;1m\]\u@\[\e[0m\]\[\e[32;1m\]\H\[\e[0m\] \[\e[31;1m\]\w\[\e[0m\]]\\$ '
fi

```

# Ubuntu

> 参考：
>
> - [官网](https://ubuntu.com/)
> - [Wiki, Ubuntu](https://en.wikipedia.org/wiki/Ubuntu)
> - [Ubuntu Manual(手册)](https://manpages.ubuntu.com/)

Ubuntu 是一个基于 Debian 的 Linux 发行版，主要由 [FOSS](https://en.wikipedia.org/wiki/Free_and_open-source_software) 组成。

Ubuntu 由英国公司 [Canonical](https://en.wikipedia.org/wiki/Canonical_(company)) 和其他开发者社区共同开发的，采用了一种精英治理模式。Canonical为每个Ubuntu版本提供安全更新和支持，从发布日期开始，直到该版本达到其指定的寿命终点(EOL)日期为止。Canonical 通过销售与 Ubuntu 相关的高级服务以及下载 Ubuntu 软件的人的捐赠来获得收入。

## 其他

Ubuntu Server 安装完成后，通常需要关闭自动更新，详见 [Debian 包管理](/docs/1.操作系统/Package%20管理/Debian%20包管理.md#包的自动更新)

# 关联文件与配置
