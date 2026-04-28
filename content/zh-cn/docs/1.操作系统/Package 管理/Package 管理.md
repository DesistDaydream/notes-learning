---
title: Package 管理
linkTitle: Package 管理
weight: 1
---

# 概述

> 参考：
>
> -

在 Linux 操作系统中，Package 就是指应用程序的安装包。保存 Package 的地方(网站、ISO 等)称为 **Repository(简称 Repo)**，我们可以从各种 Linux 发行版的官方 Repo 中下载对应的可用的 Package，以安装到这些发行版的 Linux 系统中。

> [!Attention]
> 哪怕两个发行版的包管理器相同，也不代表他们的 Package 是可以公用的，比如 CentOS 和 OpenEuler 都用 yum，但是 CentOS 的 Package 是无法装在 OpenEuler 上的，安装时将会报错(比如包与包之间 **conflict(冲突)**)

[Debian 包管理](/docs/1.操作系统/Package%20管理/Debian%20包管理.md)

[Redhat 包管理](/docs/1.操作系统/Package%20管理/Redhat%20包管理.md)

[Snap](/docs/1.操作系统/Package%20管理/Snap.md)

[Windows包管理](/docs/1.操作系统/Package%20管理/Windows包管理.md)

# Linux 各发行版的官方 Repo 站点

- 包含很多发行版的 Repo 站点: https://pkgs.org/
- OpenEuler: https://repo.openeuler.org/
- CentOS: https://centos.pkgs.org/
- Ubuntu: https://packages.ubuntu.com/
 	- 在这里可以找到 jammy 版本(20.04 TLS)的所有软件包列表: https://packages.ubuntu.com/jammy/allpackages

