---
title: OpenSSH Utilities
linkTitle: OpenSSH Utilities
date: 2024-04-19T10:38
weight: 3
---

# 概述

> 参考：
>
> - [官方 Manual(手册)](http://www.openssh.com/manual.html)
> - <https://www.myfreax.com/how-to-setup-ssh-tunneling/>
> - <https://hellolzc.github.io/2020/04/port-forwarding-with-ssh/>
> - <http://tuxgraphics.org/~guido/scripts/ssh-socks5-port-forwarding.html>

- ssh、scp、sftp # 客户端管理工具
- ssh-add、ssh-keysign、ssh-keyscan、ssh-keygen # 密钥管理工具
- sshd、sftp-server、ssh-agent # 服务端管理工具

# ssh - OpenSSH 的 ssh 客户端工具(远程登录程序)

详见 [ssh](/docs/4.数据通信/Utility/OpenSSH/ssh.md)

# scp - 基于 ssh 协议的文件传输工具

scp 是基于 SSH 的 [File transfer](/docs/4.数据通信/File%20transfer/File%20transfer.md) 工具

## Syntax(语法)

**scp \[OPTIONS] SourceFILE DestinationFILE**

Note：远程 FILE 的格式为：USER@IP:/PATH/FILE)

OPTIONS：

- **-p** #
- **-r** # 以递归方式复制，用于复制整个目录

## EXAMPLE

把本地 nginx 文件推上去复制到以 root 用户登录的 10.10.10.10 这台机器的/opt/soft/scptest 目录下

- scp /opt/soft/nginx-0.5.38.tar.gz root@10.10.10.10:/opt/soft/scptest

把以 root 用户登录的 10.10.10.10 机器中的 nginx 文件拉下来复制到本地/opt/soft 目录下

- scp root@10.10.10.10:/opt/soft/nginx-0.5.38.tar.gz /opt/soft/

基于密钥的认证,当对方主机 ssh 登录的用户的家目录存在公钥，并且公钥设置密码为空，那么以后 ssh 协议登录传输都可以直接登录而不用密码

# ssh-keygen - 在客户端生成密钥对

**ssh-keygen -t rsa \[-P ''] \[-f ~/.ssh/id_rsa]**

EXAMPLE

- ssh-keygen -t rsa -P '' -f ~/.ssh/id_rsa

# ssh-copy-id - 把生成的公钥传输至远程服务器对应用户的家目录

**ssh-copy-id \[-i \[Identity_File]] \[User@]HostIP**

Identity_File(身份文件) # 一般为 /root/.ssh/id_rsa.pub

EXAMPLE

- 将公钥拷贝到服务端
  - ssh-copy-id -i /root/.ssh/id_rsa.pub root@192.168.0.10
- 若没有 ssh-copy-id 命令，则可以这么这么弄
  - cat ~/.ssh/id_rsa.pub | ssh root@192.168.0.10 'cat >> .ssh/authorized_keys'
