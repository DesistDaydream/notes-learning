---
title: su 与  sudo
---

# 概述

>

# su

> 参考：
> - [Manual(手册),su(1)](https://man7.org/linux/man-pages/man1/su.1.html)

## 总结

使用 su 命令切换用户身份然后执行命令虽然简单，但是，也有一些致命的缺点：

- 普通用户必须知道 root 密码才可以切换到 root，这样 root 密码就泄漏了。
- 使用 su 命令切换身份，无法对切换后的身份做精细的控制，拿到超级权限的人可以为所欲为。甚至可以改掉 root 密码，让真正的管理员无法再拥有 root 权限。

这时候，就可以使用 sudo 工具

# su 配置

/etc/pam.d/su #&#x20;
/etc/pam.d/su-l #&#x20;
/etc/default/su #&#x20;
/etc/login.defs #

# sudo

> 参考：
> - [Manual(手册),sudo(8)](https://man7.org/linux/man-pages/man8/sudo.8.html)
> - [Manual(手册),sudoers(5)](https://man7.org/linux/man-pages/man5/sudoers.5.html)
> - [如何改变 sudo 日志文件](https://ostechnix.com/how-to-change-default-sudo-log-file-in-linux/)

**sudo(substitute user \[或 superuser] do)** 程序可以让当前用户使用其他的用户的权限来执行指定的命令

通过 sudo 命令，我们可以把某些 root 权限(e.g.只有 root 用户才能执行的命令)分类有针对性授权给指定的普通用户，并且普通用户不需要知道 root 密码就可以使用得到的授权来管理。效果如下所示(配置好 sudo 之后，普通用户 lichenhao 也可以通过在命令前加 sudo 来执行 root 才能执行的命令)
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ucpd2y/1616166756129-0c242239-c867-4503-8d14-2b199cab5600.png)

sudo 通过各种插件实现功能。默认插件为 sudoers，用来确定用户的 sudo 权限，sudoers 的策略，通过 /etc/sudoers 文件进行配置，或者在 LDAP 中进行配置。

# sudo 配置

**/etc/sudo.conf **# sudo 程序的配置文件
**/etc/sudoers** # suders 插件的配置文件，可以配置哪些用户可以拥有哪些权限。注意：该文件是只读的，只能通过 visudo 命令进行编辑

- **/etc/sudoers.d/\* **# /etc/sudoers 的 include 配置的默认目录

## sudo 日志配置

使用

# sudo 命令行工具

- sudo # 可以让普通用户拥有 root 权限去执行命令，sudo 的配置文件是/etc/sudoers。
- visudo # 通过 visudo 编辑/etc/sudoers，可以检查语法。

## sudu # 使用其余用户的权限执行指定的命令

**sudo \[OPTIONS] \[COMMAND]**

OPTIONS

- **-l, --list** # 查看授权情况，列出用户在主机上可用的和被禁止的命令
- **-k, --reset-timestamp** # 删除时间戳，时间戳默认 5 分钟也会失效
- **-u，--user=<STRING>** # 以指定用户执行命令。STRING 可以 用户名 或 用户 ID
- **-s, --shell** # 以目标用户运行 shell。
  - 若直接使用 `sudo -s` 命令，相当于以 root 用户运行 shell，省去了 su - root 再输入密码的操作

EXAMPLE

- sudo -u lichenhao whoami #使用用户 lichenhao 来执行 whoami 命令


    [root@master ~]# whoami
    root
    [root@master ~]# sudo -u lichenhao whoami
    lichenhao

## visudo # 编辑/etc/sudoers 文件

使用 visudo 命令可直接进入编辑模式开始配置 /etc/sudoers 文件，配置 visudo 后，使用 sudo 命令，可以让非 root 用户在执行某些命令时不用 root 密码

OPTIONS：

- **-c **# 检查 /etc/sudoers 文件的语法
- **-f, --file=sudoers** # 指定 sudoers 文件的路径
- **-q, --quiet** # less verbose (quiet) syntax error messages
- **-s, --strict** # 严格的语法检查，在编辑 sudoers 文件并保存退出后，如果语法错误，则会弹出提示
