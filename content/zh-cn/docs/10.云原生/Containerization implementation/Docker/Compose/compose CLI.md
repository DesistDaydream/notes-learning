---
title: compose CLI
linkTitle: compose CLI
weight: 3
---

# 概述

> 参考：
>
> - [官方文档](https://docs.docker.com/reference/cli/docker/compose/)

通过 run 命令，可以在容器启动失败时，进行调试

docker-compose -f docker-compose-backup.yaml run keepalived bash

这样会启动 keepalived 容器，并分配一个终端。
