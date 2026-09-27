---
title: Debian 与 Ubuntu
weight: 60
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

## 更新 LTS 版本

```bash
# 1. 先把当前系统更新到最新
sudo apt update && sudo apt upgrade -y

# 2. 确保升级工具已安装
sudo apt install update-manager-core -y

# 3. 检查升级策略
cat /etc/update-manager/release-upgrades
# Prompt=lts    → 只能升级到下一个 LTS 版本
# Prompt=normal → 也可以升级到非 LTS 版本

# 4. 执行升级
sudo do-release-upgrade
```

只检查不执行

```bash
sudo do-release-upgrade -c
```

查看当前的版本

```bash
lsb_release -a
```

# 关联文件与配置

# 版本生命周期

https://ubuntu.com/about/release-cycle

https://ubuntu.com/project/docs/release-team/list-of-releases/
