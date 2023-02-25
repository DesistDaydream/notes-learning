---
title: Debian 与 Ubuntu
---

# 概述

> 参考：
> - [Debian 官方 Manual(手册)](https://manpages.debian.org/)

```bash
groupadd wheel
usermod -G wheel lichenhao
tee /etc/sudoers.d/lichenhao > /dev/null <<EOF
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

# Ubuntu 发行版

Ubuntu Server 安装完成后，通常需要关闭自动更新~~