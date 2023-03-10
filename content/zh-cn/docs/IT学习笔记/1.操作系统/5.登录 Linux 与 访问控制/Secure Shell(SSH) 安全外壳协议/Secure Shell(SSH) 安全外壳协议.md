---
title: Secure Shell(SSH) 安全外壳协议
---

# 概述

> 参考：
> - [Wiki-SSH](https://en.wikipedia.org/wiki/Secure_Shell_Protocol)
> - <https://www.digitalocean.com/community/tutorials/how-to-set-up-ssh-keys-on-centos-8>

**Secure Shell Protocol(安全外壳协议，简称 SSH) **是一种加密的[网络传输协议](https://zh.wikipedia.org/wiki/%E7%BD%91%E7%BB%9C%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE)，可在不安全的网络中为网络服务提供安全的传输环境\[1]。SSH 通过在网络中创建安全隧道来实现 SSH 客户端与服务器之间的连接\[2]。虽然任何网络服务都可以通过 SSH 实现安全传输，SSH 最常见的用途是远程登录系统，人们通常利用 SSH 来传输命令行界面和远程执行命令。使用频率最高的场合类 Unix 系统，但是 Windows 操作系统也能有限度地使用 SSH。2015 年，微软宣布将在未来的操作系统中提供原生 SSH 协议支持\[3]，Windows 10 1809 版本已提供可手动安装的 OpenSSH 工具

# OpenSSH # SSH 协议的实现

> 参考：
> - [官网](https://www.openssh.com/)
> - [官方文件,Manual(手册)](https://www.openssh.com/manual.html)

OpenSSH 是 ssh 协议的开源实现方案，该工具为 C/S 架构，服务端默认监听在 22/tcp 端口上。如果想要连接到服务端，同样需要一个客户端来进行连接。

比如，现在有两台主机，主机 A 和主机 B，如果想要在 B 上可以操作 A，那么就需要在 A 机上安装服务端工具(openssh-server)，在 B 机上安装客户端工具(openssh-client)，然后通过 ssh 工具进行互联

Note：现在 OpenSSH 一般作为 Linux 发行版的标准远程登录工具默认安装在系统中且开机自启动。

## 关联文件

**/var/log/secure** # 登录信息日志所在位置

- 可以通过该日志获取到尝试暴力破解的 IP

**/etc/ssh/ssh_config** # OpenSSH 的 client 端配置(ssh、scp 等程序)
**/etc/ssh/sshd_config** # OpenSSH 的 server 端配置(sshd 程序)
**~/.ssh/know_hosts** # 已知的曾经连接过的主机的信息。凡是使用 ssh 连接过该主机，都会将信息记录在其中
**~/.ssh/authorized_keys** # 已经认证的公钥。如果其他 client 想要连接服务端，凡是在该文件中的公钥，都可以免密连接。

- 注意：OpenSSH 使用非对称加密的方式，与传统互联网的 https 使用方式相反。https 的公钥是交给客户端，用来验证服务端返回的网页是否可信。而 OpenSSH 则是将公钥交给服务端，用来验证客户端发送的信息是否可信。
- 这也确实符合逻辑
  - ssh 是一个客户端需要登录多个服务端，服务端要验证客户端发送的信息的真实性，要是不验证，那么就可以随便在自己服务器上执行命令了，这是不对的~
  - 而互联网通过 https 访问，则是多个客户端对应一个服务端。

**~/.ssh/config** # OpenSSH 的 client 端配置，该配置文件主要针对不同用户来使用，默认不存在，需要手动创建。
&#x20;在客户端添加如下配置内容，就可以通过名字，而不是 IP 来 ssh 登录目标主机了，还不用改 hosts 文件，也不用配置域名解析
Host centos8 User root Hostname 10.10.100.249 # 效果如下

    [root@lichenhao ~]# ssh centos8
    root@10.10.100.249's password:
    Last login: Fri Jul 10 22:56:38 2020 from 10.10.100.200
    [root@centos-8-init ~]#

工作环境一般配置

- 不要使用默认端口
- 禁止使用 protocaol version 1
- 限制可登录用户
- 设定空闲会话超时时长
- 利用防火墙设置 ssh 访问策略
- 仅监听特定的 IP 地址
- 基于口令认证时，使用强密码策略
  - tr -dc A-Za-z0-9\_ < /dev/urandom | head -c 30 | xargs # 生成 30 位随机字符串
- 禁止使用空密码
- 禁止 root 用户直接登录
- 限制 ssh 的访问频度和并发(即同时)在线数
- 做好日志

# OpenSSH 优化

- 提高 ssh 连接速度
  - 修改 sshd_config 文件中 useDNS 为 no
  - 修改 ssh_config 文件中 GSSAPIAuthentication 为 no

# 问题实例

## 使用 Xshell 多次 ssh 跳转连接后 x11 无法转发

> 参考：
> - <https://serverfault.com/a/425413>

问题描述：
在 windows 上使用 xshell 连接到主机 A，再通过主机 A 连接到主机 B。在主机 B 上打开 x11 程序失败。(x11 程序必须 virt-manager、xclock 等具有图形化的程序)

> Note：如果使用 xshell 直接连接主机 B 则无此问题。

问题原因：
缺少 x11 认证程序

解决方式：
安装程序包 yum install -y xorg-x11-xauth

xorg-x11-xaut 软件包安装完成后退出当前 shell 重新登录，即可自动配置$DISPLAY 环境变量的值为 localhost:10.0

由于现阶段 sshd_config 文件中已经默认开启 X11 转发( i.e. X11Forwarding yes )，所以 X11 转发的相关配置已不用修改

最后在使用 ssh 登录的时候加上-Y 选项，i.e. ssh -Y root@主机 B 的 IP，然后 x11 的转发效果就实现了。

https://serverfault.com/a/425413
