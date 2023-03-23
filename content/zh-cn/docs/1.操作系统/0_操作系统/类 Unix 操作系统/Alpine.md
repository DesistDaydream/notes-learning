---
title: Alpine
---

# 概述

> 参考：
>
> - [官网](https://www.alpinelinux.org/)
> - [GitHub 项目,docker-alpine](https://github.com/alpinelinux/docker-alpine)
> - [DockerHub](https://hub.docker.com/_/alpine)
> - <https://mp.weixin.qq.com/s/Qt8ASPefVG-9bZe6FO_YQw>

# APK # 包管理器

> 参考：
>
> - [官方文档](https://docs.alpinelinux.org/user-handbook/0.1a/Working/apk.html)

**Alpine Package Keeper(Alpine 包管理圆，简称 APK)** 是 Alpine 发行版的包管理工具。

## 关联文件

**/etc/apk/repositories** # 包仓库的配置文件

- 阿里仓库
  - sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
- 中科大仓库
  - sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

**/var/cache/apk/\*** # APK 程序运行时产生的缓存文件保存路径

## apk 命令行工具

### Syntax(语法)

**apk \[Global OPTIONS] COMMAND \[COMMAND OPTIONS]**
Global OPTIONS

- --no-cache # 不在 /var/cache/apk/ 目录下生成缓存，并且也不使用该目录下的缓存。

COMMAND

- 安装和移除包命令
  - **add** # 为正在运行的系统添加新包或升级包
  - **del** # 从正在运行的系统中删除包
- 系统维护命令(管理包的元数据)
  - cache Maintenance operations for locally cached package repository
  - **fix** # 尝试修复或升级已安装的包
  - update Update the index of available packages
  - upgrade Upgrade the currently installed packages
- 查询包的信息
  - dot Create a graphviz graph description for a given package
  - info # Prints information about installed or available packages
  - list # List packages by PATTERN and other criteria
  - policy Display the repository that updates a given package, plus repositories that also offer the package
  - search Search for packages or descriptions with wildcard patterns
- 仓库管理命令(管理包源)
  - **fetch**# 下载包，但是不安装它。
  - index create a repository index from a list of packages
  - verify Verify a package signature
  - stats Display statistics, including number of packages installed and available, number of directories and files, etc.
  - manifest Display checksums for files contained in a given package

del OPTIONS

- **-r, --rdepends** # 递归删除所有顶级反向依赖项

# Alpine 容器镜像

alpine 是基于 Alpine Linux 发型版的最小容器映像，具有完整的软件包索引，大小仅为 5 MB！Alpine 采用了 musl libc 和 busybox 以减小系统的体积和运行时的资源消耗。同时，Alpine 具有自己的**包管理器 APK**。可以在 <https://pkgs.alpinelinux.org/packages> 查询到所有包的信息，并且可以直接通过 `apk` 命令查询和安装各种软件。

## 添加时区

安装 tzdata 包，后，配置 TZ 环境变量或者创建 /etc/localtime 软链接文件即可。示例如下：

```bash
FROM alpine

# 设置时区为上海
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk add --no-cache tzdata && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
```

## 设置系统语言，防止中文乱码

```bash
FROM alpine:3.11.5
ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US.UTF-8
```

## 常见问题

### golang 在 alpine 镜像下 hosts 定义的域名不生效-及解决 x509certificates

> 参考：
>
> - [Go Issue #22846](https://github.com/golang/go/issues/22846)
> - [CSDN](https://blog.csdn.net/huangruifeng/article/details/96594065)
> - [腾讯云文章](https://cloud.tencent.com/developer/article/1756065)

golang 在 alpine 镜像下 hosts 定义的域名不生效

解决方案

    echo "hosts: files dns" > /etc/nsswitch.conf

> 参考：<https://github.com/golang/go/issues/22846>

以下为调整后的 Dockerfile

    FROM alpine
    RUN apk update #解决 apk下载失败问题 ERROR: unsatisfiable constraints
    RUN apk add --no-cache ca-certificates # 在go程序中无法访问https链接，解决x509certificates
    RUN echo "hosts: files dns" > /etc/nsswitch.conf #go程序在alpine下不解析hosts文件
