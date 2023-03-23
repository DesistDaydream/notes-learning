---
title: Docker 权限管理
---

# 概述

# Docker Capabilities

Capabilities 详见 [Linux Capabilities 权限管理章节](/docs/IT学习笔记/1.操作系统/5.登录%20Linux%20 与%20 访问控制/Access%20Control(访问控制)/Capabilities(能力)%20 管理.md Linux 与 访问控制/Access Control(访问控制)/Capabilities(能力) 管理.md)

我们说 Docker 容器本质上就是一个进程，所以理论上容器就会和进程一样会有一些默认的开放权限，默认情况下 Docker 会删除必须的 `capabilities` 之外的所有 `capabilities`，因为在容器中我们经常会以 root 用户来运行，使用 `capabilities` 现在后，容器中的使用的 root 用户权限就比我们平时在宿主机上使用的 root 用户权限要少很多了，这样即使出现了安全漏洞，也很难破坏或者获取宿主机的 root 权限，所以 Docker 支持 `Capabilities` 对于容器的安全性来说是非常有必要的。
不过我们在运行容器的时候可以通过指定 `--privileded` 参数来开启容器的超级权限，这个参数一定要慎用，因为他会获取系统 root 用户所有能力赋值给容器，并且会扫描宿主机的所有设备文件挂载到容器内部，所以是非常危险的操作。
但是如果你确实需要一些特殊的权限，我们可以通过 `--cap-add` 和 `--cap-drop` 这两个参数来动态调整，可以最大限度地保证容器的使用安全。下面表格中列出的 `Capabilities` 是 Docker 默认给容器添加的，我们可以通过 `--cap-drop` 去除其中一个或者多个：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hlragi/1621522556084-0aa763d8-6f2d-4e2f-8f69-1953f75511be.png)docker capabilities
下面表格中列出的 `Capabilities` 是 Docker 默认删除的，我们可以通过`--cap-add`添加其中一个或者多个：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hlragi/1621522556093-902979da-99b9-479e-8b11-db55a3d83605.png)docker drop capabilities

> `--cap-add`和`--cap-drop` 这两参数都支持`ALL`值，比如如果你想让某个容器拥有除了`MKNOD`之外的所有内核权限，那么可以执行下面的命令： `$ sudo docker run --cap-add=ALL --cap-drop=MKNOD ...`

比如现在我们需要修改网络接口数据，默认情况下是没有权限的，因为需要的 `NET_ADMIN` 这个 `Capabilities` 默认被移除了：

    $ docker run -it --rm busybox /bin/sh
    / # ip link add dummy0 type dummy
    ip: RTNETLINK answers: Operation not permitted
    / #

所以在不使用 `--privileged` 的情况下（不建议）我们可以使用 `--cap-add=NET_ADMIN` 将这个 `Capabilities` 添加回来：

    $ docker run -it --rm --cap-add=NET_ADMIN busybox /bin/sh
    / # ip link add dummy0 type dummy
    / #

可以看到已经 OK 了。
