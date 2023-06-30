---
title: Open Containers Initiative(开放容器倡议)
---

# 概述

> 参考：
> 
> - [OCI 官网](https://opencontainers.org/)
> - [GitHub 账户,OpenContainers](https://github.com/opencontainers)
> - [segmentfault,《走进 docker》系列文章](https://segmentfault.com/u/public0821/articles)

**Open Containers Initiative(开放容器倡议，简称 OCI)** 是一个轻量级的，开放的治理结构（项目），由 Linux Foundation 主持成立，其明确目的是围绕 Container 镜像格式和运行时创建 开放的行业标准。OCI 由 Docker，CoreOS 和其他容器行业领导者于 2015 年 6 月 22 日启动。

OCI 公有如下几个个规范：
一开始有两个

- **Image-spec(镜像规范)** # 容器镜像所包含的内容以及格式都遵循统一的格式标准，由 OCI 负责维护，官方详解地址为：image-spec
- **Runtime-spec(运行时规范)** # 容器运行时的内容以及格式都遵循统一的格式标准，由 OCI 负责维护，官方详解地址为：runtime-spec

后来新加的一个

- **Distribution-spec(分发规范)** #

在所有企业、各人在构建镜像、运行容器时，都应该遵守 OCI 标准，比如想用 docker 工具构建一个镜像，那么构建出来的镜像规范，必须符合 OCI 标准。其他类似 docker 的工具同理。如果想自己开发一个构建镜像的工具或者运行容器的运行时，都需要符合 OCI 的标准。这样大家都遵守同一套规范，才有利于技术的发展。
