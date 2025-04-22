---
title: Compose 文件规范
linkTitle: Compose 文件规范
weight: 2
---

# 概述

> 参考：
>
> - [GitHub 项目，compose-spec/compose-spec](https://github.com/compose-spec/compose-spec)
> - [Compose 规范](https://compose-spec.io/)
> - [Docker 官方文档，参考 - compose 文件](https://docs.docker.com/reference/compose-file/)

Compose 文件是一个 [YAML](/docs/2.编程/无法分类的语言/YAML.md) 格式的配置文件，Compose 将每个容器抽象为一个 service。顶层字段 service 的下级字段，用来定义该容器的名称。

一个 Docker Compose 文件中通常包含如下顶级字段：

- **version** # **必须的**。
- **services**(map\[STRING][services](#services))
- **networks**([networks](#networks))
- **volumes**([volumes](#volumes))
- **secrets**([secrets](#secrets))

# version

指定本 yaml 依从的 compose 哪个版本制定的。

# services

详见 [services](/docs/10.云原生/Containerization%20implementation/Docker/Compose/services.md)

# networks

> 参考：
>
> - [官方文档，参考 - Compose 文件参考 - Networks 顶级元素](https://docs.docker.com/reference/compose-file/networks/)

**attachable: BOOLEAN** # 该网络是否可以被其他容器加入

**external: BOOLEAN** # 该网络是否由外部维护。若为 true，则该网络不受本 Compose 的管理。`默认值：false`

**name: STRING** # 指定网络名称

# volumes

# configs

# secrets
