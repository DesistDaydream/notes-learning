---
title: 登录 Linux 与 访问控制
linkTitle: 登录 Linux 与 访问控制
weight: 1
---

# 概述

> 参考：

想要登录 Linux，必须通过 [**Terminal(终端)**](/docs/1.操作系统/Terminal%20与%20Shell/Terminal%20与%20Shell.md)，我们才可以与操作系统进行交互。

本质上，想要登录 Linux，必然需要调用某些程序(比如 Shell)，以便分配一个终端。通常，我们有多种方式可供选择：

- 本地命令行登录
- 远程命令行登录
- 图形界面登录

Linux 操作系统是一个多用户操作系统，所以除了 **Terminal(终端)** 以外，还需 **Account(账户)** 才可以登录上去，Linux 操作系统允许多个用户访问安装在一台机器上的单个系统。每个 **User(用户)** 都在自己的 **Account(账户)** 下操作。因此，Account Manager 代表了 Linux 系统管理的核心要素。

# 登录 Linux

我们可以通过多种方式登录 Linux

- 本地登录
- 远程登录

## 通过本地 TTY 登陆 Linux 系统

登录 Linux 最基本的方式，就是使用 `login` 程序。

### login 程序

由于历史原因，`login` 可能被包含在两个包中：

- util-linux
- shadow-utils

#### login 的登录行为

当我们刚刚安装完操作系统，systemd-logind.service 服务会让我们看到这样的画面

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/linux_login/1634785246289-3a353c73-2899-4b6c-8341-ffc4a02008ef.png)

想要在服务器本地登录系统，则需要进行认证，在输入用户名之后，实际上是调用了 `login` 这个二进制程序，看到：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/linux_login/1634785206973-885fa7fb-3dfb-4fb1-80c8-1c56cb903974.png)

此时我们通过远程方式(如果可以的话)登录服务器，查看进程，将会看到名为 login 的进程

```bash
~]$ pgrep login -alf
600 /lib/systemd/systemd-logind
1476 /bin/login -p --
```

当我们输入完密码，经过 [Access Control(访问控制)](/docs/1.操作系统/登录%20Linux%20与%20访问控制/Access%20Control(访问控制)/Access%20Control(访问控制).md) 相关程序的认证之后，login 工具会为我们分配一个 ttyX 的终端设备，然后我们就可以通过 tty 所关联的 Shell(通常是 bash)，与系统进行交互

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/linux_login/1634785329507-0cb1fcec-8c6e-4fd0-a99f-005a2b19807e.png)

#### login 关联文件与配置

**/etc/login.defs** # shadow 与 password 套件的配置文件。

**/etc/pam.d/login** #

## 通过远程的方式来登陆 Linux 系统

### ssh 程序

详见：[OpenSSH](/docs/4.数据通信/Utility/OpenSSH/OpenSSH.md)

```bash
root       981     1  0 Jul08 ?        Ss     0:00 /usr/sbin/sshd -D
root      1947   981  0 09:05 ?        Ss     0:00  \_ sshd: root@pts/0
root      1949  1947  1 09:05 pts/0    Ss     0:00      \_ -bash
root      1970  1949  0 09:05 pts/0    R+     0:00          \_ ps -ef f
```

OpenSSH 会为用户分配一个 Pseudoterminal(伪终端，即 pts) 以便用户可以与操作系统进行交互。

# 登录提示信息

通过 `touch ~/.hushlogin` 命令可以为当前用户禁用欢迎信息。

## MOTD

> 参考：
>
> - [Wiki, MOTD](https://en.wikipedia.org/wiki/Motd_(Unix))
> - [Manual(手册)，MOTD](https://man7.org/linux/man-pages/man5/motd.5.html)

**Message of the day(每日消息，简称 MOTD)** 是一种比向所有用户发送一个邮件更有效的发送共同的信息的方式

#### MOTD 关联文件

- **/etc/default/motd-news** # 动态 MOTD 新闻信息配置，新闻信息主要是互联网相关的
- **/etc/update-motd.d/** # MOTD 执行脚本保存路径
- **~/.hushlogin** # 该文件存在时，将为当前用户禁用 MOTD 消息
- **/etc/pam.d/sshd** # 使用 ssh 登录后显示的 MOTD，可以通过 sshd 的 [PAM 配置文件](/docs/1.操作系统/登录%20Linux%20与%20访问控制/PAM/PAM%20配置文件.md#/etc/pam.d/sshd)进行配置。

## 示例

`There were 9 failed login attempts since the last successful login.`

TODO: 这段信息来源未知。CentOS 系自带；Ubuntu 无法通过配置 faillock PAM 模块显示出来，也不知道该如何显示出来。

`Last login: Sun May  5 22:16:35 2024 from 192.168.254.254`

这段信息由 sshd_config 文件中的 PrintLastLog yes 指令配置

`其他`

其他的信息可能是由 /etc/pam.d/sshd 配置中的 pam_motd.so 模块产生的

也有可能是 /etc/profile 相关的脚本文件输出的

# 访问控制

Linux 的登录与访问控制是相辅相成的，一个用户想要登录 Linux，通常来说都需要经过访问控制系统对其所使用的账户进行认证，只有认证通过后，才可以正常登录。

一个正常的 Linux 发行版操作系统，通常都提供了多种方式

- 密码
- 会话
- 账户锁定
- 等等......

## Account Manager(账户管理)

详见 [Account Manager(账户管理)](/docs/1.操作系统/登录%20Linux%20与%20访问控制/Account%20Manager(账户管理)/Account%20Manager(账户管理).md)

# 多窗口操作

登录服务器后，我们可以重复登录，以便在多个窗口执行不同的操作以观察服务器状态或排查问题。

但是当我们在机房通过显示器连接到服务器时，是不像使用 ssh 命令一样方便的，但是依然可以实现多窗口操作。

使用 `Ctrl + Alt + F<X>` 快捷键，即可打开其他窗口，`Ctrl + Alt + F2` 切换到第二个窗口，F1 可以切回第一个默认窗口。
