---
title: crictl 命令行工具
---

# 概述

> 参考：
> - [项目地址](https://github.com/kubernetes-sigs/cri-tools)
> - [使用 crictl 对 kubernetes 进行调试](https://kubernetes.io/docs/tasks/debug-application-cluster/crictl/)

crictl 用于为 kubelet CRI 进行调试的 命令行工具

cri-tools 旨在为 Kubelet CRI 提供一系列调试和验证工具，其中包括：

- crictl: kubelet 的 CRI 命令行工具
- critest:kubelet CRI 的验证测试套件

用白话说就是：kubelet 如果要与 CRI 对接，那么如何检测对接成功呢，就是使用 crictl 工具来测试。还可以对已经与 kubelet 建立连接的 CRI 执行相关操作，比如启停容器等。

Note：要想使用 crictl 命令行工具，必须要先进行配置，指定好要操作的 CRI 的 endpoint，才可以正常使用

# crictl 配置

**/etc/crictl.yaml** # crictl 命令行工具运行时配置文件

## 基本配置文件示例

    runtime-endpoint: unix:///run/containerd/containerd.sock
    image-endpoint: unix:///run/containerd/containerd.sock
    timeout: 10
    debug: true

# crictl 命令行工具

**crictl \[Global OPTIONS] COMMAND \[COMMAND OPETIONS] \[ARGUMENTS...]**
COMMMAND

- attach Attach to a running container
- create Create a new container
- exec Run a command in a running container
- version Display runtime version information
- images List images
- inspect Display the status of one or more containers
- inspecti Return the status of one or more images
- inspectp Display the status of one or more pods
- logs Fetch the logs of a container
- port-forward Forward local port to a pod
- ps List containers
- pull Pull an image from a registry
- runp Run a new pod
- rm Remove one or more containers
- rmi Remove one or more images
- rmp Remove one or more pods
- pods List pods
- start Start one or more created containers
- **info** #显示与 crictl 对接的 CRI 信息
- stop Stop one or more running containers
- stopp Stop one or more running pods
- update Update one or more running containers
- config Get and set crictl options
- stats List container(s) resource usage statistics
- completion Output bash shell completion code
- help, h Shows a list of commands or help for one command

OPTIONS

EXAMPLE

- crictl info #显示与 crictl 对接的 CRI 信息
- crictl pull docker.io/lchdzh/pause:3.1 #拉取 dockerhub 上的 lchdzh 中的 pause 镜像
