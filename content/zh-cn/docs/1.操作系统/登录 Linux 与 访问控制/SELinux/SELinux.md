---
title: SELinux
---

# 概述

> 参考：
>
> - [Wiki, Security-Enhanced Linux](https://en.wikipedia.org/wiki/Security-Enhanced_Linux)

**Security Enhanced Linux(安全强化的 Linux，简称 SELinux)** 是由美国国家安全局 (NSA) 开发的，当初开发这玩意儿的目的是因为很多企业界发现， 通常系统出现问题的原因大部分都在于『内部员工的资源误用』所导致的，实际由外部发动的攻击反而没有这么严重。 那么什么是『员工资源误用』呢？举例来说，如果有个不是很懂系统的系统管理员为了自己设定的方便，将网页所在目录 /var/www/html/ 的权限设定为 drwxrwxrwx 时，你觉得会有什么事情发生？

当初设计的目标：避免资源的误用

现在我们知道所有的系统资源都是透过进程来进行存取的，那么 /var/www/html/ 如果设定为 777 ，代表所有进程均可对该目录存取，万一你真的有启动 WWW 服务器软件，那么该软件所触发的进程将可以写入该目录， 而该进程却是对整个 Internet 提供服务的！只要有心人接触到这支进程，而且该进程刚好又有提供用户进行写入的功能， 那么外部的人很可能就会对你的系统写入些莫名其妙的东西！那可真是不得了！一个小小的 777 问题可是大大的！

为了控管这方面的权限与进程的问题，所以美国国家安全局就着手处理操作系统这方面的控管。 由于 Linux 是自由软件，程序代码都是公开的，因此她们便使用 Linux 来作为研究的目标， 最后更将研究的结果整合到 Linux 核心里面去，那就是 SELinux 啦！所以说， SELinux 是整合到核心的一个模块喔！ 更多的 SELinux 相关说明可以参考：

这也就是说：其实 SELinux 是在进行进程、文件等细部权限设定依据的一个核心模块！ 由于启动网络服务的也是进程，因此刚好也能够控制网络服务能否存取系统资源的一道关卡！ 所以，在讲到 SELinux 对系统的访问控制之前，我们得先来回顾一下之前谈到的系统文件权限与用户之间的关系。因为先谈完这个你才会知道为何需要 SELinux 的啦！

目前 SELinux 依据启动与否，共有三种模式，分别如下：

1. enforcing：强制模式，代表 SELinux 运作中，且已经正确的开始限制 domain/type 了；
2. permissive：宽容模式：代表 SELinux 运作中，不过仅会有警告讯息并不会实际限制 domain/type 的存取。这种模式可以运来作为 SELinux 的 debug 之用
3. disabled：关闭，SELinux 并没有实际运作。

# SELinux 关联文件与配置

/etc/selinux/config #

> `sed -i 's@^\(SELINUX=\).*@\1disabled@' /etc/selinux/config` 关闭 SELinux

# 命令行工具

## setenforce {0|1} - 设定 selinux 模式

0 为 permissive 宽容模式

1 为 Enforcing 强制模式

## getenforce - 查看当前 selinux 模式

## sestatus  查看 selinux 状态
