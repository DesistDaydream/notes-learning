---
title: vsftpd
linkTitle: vsftpd
weight: 20
---

# 概述

> 参考：
>
> - https://security.appspot.com/vsftpd.html
> - [Wiki, Vsftpd](https://en.wikipedia.org/wiki/Vsftpd)

vsftpd 是实现 [FTP](/docs/4.数据通信/Protocol/FTP.md) 协议的程序

## vsftp 关联文件与配置

**/etc/vsftpd.conf** # vsftpd 程序的配置文件

**/etc/ftpusers** # 此文件包含 *禁止* FTP 登录的用户名，通常有"root"，"uucp"，"news" 之类，因为这些用户权限太高，登录 FTP 误操作危险性大。

- User NAME
- User NAME
- User NAME
- .......
- User NAME

配置文件 keywords 说明

最小基本配置

- write_enable=YES # 对 ftp 服务器是否有写的权限
- local_enable=YES # 是否允许本地用户登录(本地用户即 ftp 服务器自身的用户)
- anonymous_enable=NO # 是否允许匿名登录

扩展配置

- chroot_local_user=YES # 是否启动本地用户 chroot 规则，chroot 改变登录 ftp 的本地用户根的目录位置
- allow_writeable_chroot=YES # 允许在限定目录有写权限
- chroot_list_enable=YES # 是否启动不受 chroot 规则的用户名单
- chroot_list_file=/etc/vsftpd.chroot_list # 定义不受限制的用户名单在哪个文件中
- pam_service_name=vsftpd 改为 pam_service_name=ftp # 如果始终登录时候提示密码错误，则修改此项

vsftp 可以使用 [Chroot](/docs/1.操作系统/Kernel/Process/Chroot.md) 功能。比如：下面第一个图是不启动 Chroot 规则的情况，第二张图是启用 Chroot 规则的情况，可以看到当使用 Chroot 时，`/srv/ftp/` 目录对于 ftp 程序来说是作为 `/` 存在的。由于这个原因，所以启动 Chroot 的时候，ftp 工具无法访问所设定的 `/` 目录以外的其他目录

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pvqe8m/1616165219993-ce6cd857-e9ba-4af0-b7fc-7d77cf547d84.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pvqe8m/1616165220004-51f8038e-598e-427a-9b04-8f1987475f04.jpeg)

- /etc/vsftpd.chroot_list # 此文件包含对服务器上所有 FTP 内容有权限的用户名，该文件中每个用户名写一行

## FTP 命令行工具说明

ftp \[\[UserName@]目标 IP] # 可以使用 UserName 用户来登录目标 IP 的 ftp 服务器或直接进入 ftp 客户端

- open \[IP] # 连接某个 FTP 服务器。
- close # 关闭连接
- user # 对已经连接的 FTP 再输入一次用户名和密码，类似于更改登陆账户
- ？ # 查看可以使用的命令
- get \[需要下载的文件] \[下载路径以及定义文件名] # 下载文件
- put \[需要上传的文件] \[上传路径以及定义文件名] # 上传文件

注意：这些命令可以直接输入然后根据提示再输入想要执行的内容；也可直接在命令后接需要完成的动作，直接执行

注意文件夹权限以免无法上传下载

# ftp 被动模式注意事项

如果 iptables 仅仅放了 21 端口，在启动被动模式后，需要给 iptables 添加模块以便被动模式正常运行

在 /etc/sysconfig/iptalbes-conf 文件中，修改成 IPTABLES_MODULES="nf_conntrack_ftp nf_nat_ftp" 这个内容
