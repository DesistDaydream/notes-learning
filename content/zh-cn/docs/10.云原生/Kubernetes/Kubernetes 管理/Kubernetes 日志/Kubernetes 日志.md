---
title: Kubernetes 日志
linkTitle: Kubernetes 日志
weight: 1
---

# 概述

> 参考
>
> - [官方文档，概念 - 集群管理 - 日志架构](https://kubernetes.io/docs/concepts/cluster-administration/logging/)

集群级日志架构需要一个单独的后端来存储、分析和查询日志。Kubernetes 不提供日志数据的原生存储解决方案。相反，有许多与 Kubernetes 集成的日志记录解决方案。

# Kubernetes 日志管理机制

在 Kubernetes 中日志也主要有两大类：

- 应用 Pod 日志；
- Kuberntes 集群组件日志；

## 应用 Pod 日志

Kubernetes Pod 的日志管理是基于 Docker 引擎的，Kubernetes 并不管理日志的轮转策略，日志的存储都是基于 Docker 的日志管理策略。k8s 集群调度的基本单位就是 Pod，而 Pod 是一组容器，所以 k8s 日志管理基于 Docker 引擎这一说法也就不难理解了，最终日志还是要落到一个个容器上面。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dstb2v/1616116985979-e135fbea-a7b9-4da3-81ab-ae6e56dc742b.png)

假设 Docker 日志驱动为 json-file，那么在 k8s 每个节点上，kubelet 会为每个容器的日志创建一个软链接，软连接存储路径为：/var/log/containers/，软连接会链接到 /var/log/pods/ 目录下相应 pod 目录的容器日志，被链接的日志文件也是软链接，最终链接到 Docker 容器引擎的日志存储目录：/var/lib/docker/container 下相应容器的日志。另外这些软链接文件名称含有 k8s 相关信息，比如：Pod id，名字空间，容器 ID 等信息，这就为日志收集提供了很大的便利。

举例：我们跟踪一个容器日志文件，证明上述的说明，跟踪一个 kong Pod 日志，Pod 副本数为 1

/var/log/containers/kong-kong-d889cf995-2ntwz_kong_kong-432e47df36d0992a3a8d20ef6912112615ffeb30e6a95c484d15614302f8db03.log

------->

/var/log/pods/kong_kong-kong-d889cf995-2ntwz_a6377053-9ca3-48f9-9f73-49856908b94a/kong/0.log

------->

/var/lib/docker/containers/432e47df36d0992a3a8d20ef6912112615ffeb30e6a95c484d15614302f8db03/432e47df36d0992a3a8d20ef6912112615ffeb30e6a95c484d15614302f8db03-json.log

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dstb2v/1616116985855-728179e1-7d6a-48da-a79e-d0e4c15be782.png)

## Kuberntes 集群组件日志

Kuberntes 集群组件日志分为两类：

- 运行在容器中的 Kubernetes scheduler 和 kube-proxy。
- 未运行在容器中的 kubelet 和容器 runtime，比如 Docker。

在使用 systemd 机制的服务器上，kubelet 和容器 runtime 写入日志到 journald。如果没有 systemd，他们写入日志到 /var/log 目录的 .log 文件。容器中的系统组件通常将日志写到 /var/log 目录，在 kubeadm 安装的集群中它们以静态 Pod 的形式运行在集群中，因此日志一般在 /var/log/pods 目录下。

# Kubernetes 集群日志收集方案

Kubernetes 本身并未提供集群日志收集方案，k8s 官方文档给了三种日志收集的建议方案：

- 使用运行在每个节点上的节点级的日志代理
- 在应用程序的 pod 中包含专门记录日志 sidecar 容器
- 应用程序直接将日志传输到日志平台

## 节点级日志代理方案

从前面的介绍我们已经了解到，k8s 每个节点都将容器日志统一存储到了 /var/log/containers/ 目录下，因此可以在每个节点安装一个日志代理，将该目录下的日志实时传输到日志存储平台。

由于需要每个节点运行一个日志代理，因此日志代理推荐以 DaemonSet 的方式运行在每个节点。官方推荐的日志代理是 fluentd，当然也可以使用其他日志代理，比如 filebeat，logstash 等。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dstb2v/1616116985903-131ba156-b0a4-448f-a2f0-474461f0d060.png)

## sidecar 容器方案

有两种使用 sidecar 容器的方式：

- sidecar 容器重定向日志流
- sidecar 容器作为日志代理

### sidecar 容器重定向日志流

这种方式基于节点级日志代理方案，sidecar 容器和应用容器在同一个 Pod 运行，这个容器的作用就是读取应用容器的日志文件，然后将读取的日志内容重定向到 stdout 和 stderr，然后通过节点级日志代理统一收集。这种方式不推荐使用，缺点就是日志重复存储了，导致磁盘使用会成倍增加。比如应用容器的日志本身打到文件存储了一份，sidecar 容器重定向又存储了一份（存储到了 /var/lib/docker/containers/ 目录下）。这种方式的应用场景是应用本身不支持将日志打到 stdout 和 stderr，所以才需要 sidecar 容器重定向下。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dstb2v/1616116985876-0caf78da-d1b0-481d-b01a-7889a5a39e42.png)

### sidecar 容器作为日志代理

这种方式不需要节点级日志代理，和应用容器在一起的 sidecar 容器直接作为日志代理方式运行在 Pod 中，sidecar 容器读取应用容器的日志，然后直接实时传输到日志存储平台。很显然这种方式也存在一个缺点，就是每个应用 Pod 都需要有个 sidecar 容器作为日志代理，而日志代理对系统 CPU、和内存都有一定的消耗，在节点 Pod 数很多的时候这个资源消耗其实是不小的。另外还有个问题就是在这种方式下由于应用容器日志不直接打到 stdout 和 stderr，所以是无法使用 kubectl logs 命令查看 Pod 中容器日志。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dstb2v/1616116985873-cb51bd59-c44d-4460-aa4e-35d6507e4fff.png)

## 应用程序直接将日志传输到日志平台

这种方式就是应用程序本身直接将日志打到统一的日志收集平台，比如 Java 应用可以配置日志的 appender，打到不同的地方，很显然这种方式对应用程序有一定的侵入性，而且还要保证日志系统的健壮性，从这个角度看应用和日志系统还有一定的耦合性，所以个人不是很推荐这种方式。

总结：综合对比上述三种日志收集方案优缺点，更推荐使用节点级日志代理方案，这种方式对应用没有侵入性，而且对系统资源没有额外的消耗，也不影响 kubelet 工具查看 Pod 容器日志。
