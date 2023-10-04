---
title: logs
linkTitle: logs
date: 2023-11-04T10:47
weight: 20
---

# 概述

> 参考：

logs 命令可以打印 pod 中的 container 的日志

## kubectl logs \[-f] \[-p] (POD | TYPE/NAME) \[-c CONTAINER] \[options]

语法结构

- kubectl logs \<PodName> # 查看指定 pod 的日志
- kubectl logs -f \<PodName> # 类似 tail -f 的方式查看(tail -f 实时查看日志文件 tail -f 日志文件 log)
- kubectl logs \<PodName> -c \<container_name> # 查看指定 pod 中指定容器的日志

OPTIONS

- **-f** # 实时查看日志文件，类似于 tailf
- **-p,--previous** # 输出 pod 中曾经运行过，但目前已终止的容器的日志。(i.e 查看一个 container 重启之前的日志，用于排障)

EXAMPLE

- kubectl logs --namespace=kube-system calico-node-krgz6 calico-node # 查看 calico-node-krgz6 这个 pod 的日志
