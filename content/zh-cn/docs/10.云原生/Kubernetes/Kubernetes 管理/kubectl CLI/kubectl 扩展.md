---
title: kubectl 扩展
---

# 概述

> 参考：
>
> - [官方文档，任务 - 扩展 kubectl - kubectl 插件](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/)

kubectl 有一个称为 **Plugins(插件)** 机制，可以扩展 kubectl 工具的能力。通过插件，就相当于为 kubectl 添加了子命令。

# 安装 kubectl 插件

插件是一个独立的可执行文件，名称以 `kubectl-` 开头。 要安装插件，只需将此可执行文件移动到 `$PATH` 中的任何位置。

Kubernetes SIG 研发了一款名为 **Krew** 的插件，这是一个可以管理插件的插件，Krew 之于 Kubectl，就好像 yum/apt 之于 CentOS/Ubuntu。可以使用 [Krew](https://krew.dev/) 来发现和安装开源的 kubectl 插件。

> **注意：** Krew [插件索引](https://krew.sigs.k8s.io/plugins/) 所维护的 kubectl 插件并未经过安全性审查。 你要了解安装和运行第三方插件的安全风险，因为它们本质上时是一些在你的机器上 运行的程序。

# Krew

> 参考：
>
> - [GitHub 项目，kubernetes-sigs/krew](https://github.com/kubernetes-sigs/krew/)

## Krew 配置

**~/.krew/** # Krew 配置文件与存储路径。

- **./bin/** # Krew 安装的插件的软连接
- **./index/default/plugins/** # Krew 发现的插件元数据，想要安装插件，就会通过这里面的元数据信息进行。
- **./receipts/** # 已安装的插件的元数据。
- **./store/** # Krew 存储路径，所有安装的插件的二进制文件都会在该目录下。

# 常见 kubectl 插件

kubectl 插件管理工具，项目地址：<https://github.com/kubernetes-sigs/krew-index/blob/master/plugins.md>

- **neat** # 让 kubectl get -o yaml 的输出更简洁
  - 项目地址：<https://github.com/itaysk/kubectl-neat>
- **node-shell** # 通过 kubectl 命令直接进入 node 的 shell 中
  - 项目地址：[GitHub 项目](https://github.com/kvaps/kubectl-node-shell)

# 好用的 kubectl 扩展工具

## kubecm

项目地址：<https://github.com/sunny0826/kubecm>

整合 kubectl 的 config 文件，并可以简单得切换 kubectl 所要操作的目标集群
