---
title: "Rsync"
linkTitle: "Rsync"
weight: 20
---

# 概述

> 参考：
> 
> - [GitHub 项目，WayneD/rsync](https://github.com/WayneD/rsync)
> - [官网](https://rsync.samba.org/)
> - [rsync+inotify 数据实时同步介绍与 K8s 实战应用](https://mp.weixin.qq.com/s/VxnDEQ8e3yQOLJi0JwtyjA)

rsync（remote sync） 远程同步，rsync 是 linux 系统下的数据镜像备份工具。使用快速增量备份工具 Remote Sync 可以远程同步，支持本地复制，或者与其他 SSH、rsync 主机同步。已支持跨平台，可以在 Windows 与 Linux 间进行数据同步。rsync 监听端口：873，rsync 运行模式：C/S。
