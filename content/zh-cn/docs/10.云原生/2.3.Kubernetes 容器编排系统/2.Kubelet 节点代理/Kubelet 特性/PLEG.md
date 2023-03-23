---
title: PLEG
---

# 概述

> 参考：
> - [公众号,运维开发故事-PLEG is not healthy？幕后黑手居然是它！](https://mp.weixin.qq.com/s/lPYd9tNQyjidJ-sLt2sDLg)
> - [公众号,云原生实验室-Kubelet 中的 “PLEG is not healthy” 到底是个什么鬼？](https://mp.weixin.qq.com/s/t7H2MQ2429LQB9XfrB23YA)
> - <https://developers.redhat.com/blog/2019/11/13/pod-lifecycle-event-generator-understanding-the-pleg-is-not-healthy-issue-in-kubernetes#>

[
](https://mp.weixin.qq.com/s/lPYd9tNQyjidJ-sLt2sDLg)

## 问题描述

环境 ：ubuntu18.04，自建集群 k8s 1.18 ，容器运行时 docker。

现象：某个 Node 频繁 NotReady，kubectl describe 该 Node，出现如下报错日志：

```bash
PLEG is not healthy: pleg was last seen active 3m46.752815514s ago; threshold is 3m0s
```

频率在 5-10 分钟就会出现一次。

## 我们首先要明白 PLEG 是什么？

**Pod Lifecycle Event Generator(Pod 生命周期事件生成器，简称 PLEG)** 是 Kubelet 中的一个模块，主要职责就是通过每个匹配的 Pod 级别事件来调整容器运行时的状态，并将调整的结果写入缓存，使 Pod 的缓存保持最新状态。先来聊聊 PLEG 的出现背景。在 Kubernetes 中，每个节点上都运行着一个守护进程 Kubelet 来管理节点上的容器，调整容器的实际状态以匹配 spec 中定义的状态。具体来说，Kubelet 需要对两个地方的更改做出及时的回应：

1. Pod spec 中定义的状态
2. 容器运行时的状态

对于 Pod，Kubelet 会从多个数据来源 watch Pod spec 中的变化。对于容器，Kubelet 会定期（例如，10s）轮询容器运行时，以获取所有容器的最新状态。随着 Pod 和容器数量的增加，轮询会产生不可忽略的开销，并且会由于 Kubelet 的并行操作而加剧这种开销（为每个 Pod 分配一个 goruntine，用来获取容器的状态）。轮询带来的周期性大量并发请求会导致较高的 CPU 使用率峰值（即使 Pod 的定义和容器的状态没有发生改变），降低性能。最后容器运行时可能不堪重负，从而降低系统的可靠性，限制 Kubelet 的可扩展性。为了降低 Pod 的管理开销，提升 Kubelet 的性能和可扩展性，引入了 PLEG，改进了之前的工作方式：

- 减少空闲期间的不必要工作（例如 Pod 的定义和容器的状态没有发生更改）。
- 减少获取容器状态的并发请求数量。

所以我们看这一切都离不开 kubelet 与 pod 的容器运行时。![](https://notes-learning.oss-cn-beijing.aliyuncs.com/35160200-5f9f-4cfd-911c-afa960062a5c/640)
一方面，kubelet 扮演的是集群控制器的角色，它定期从 API Server 获取 Pod 等相关资源的信息，并依照这些信息，控制运行在节点上 Pod 的执行；

另外一方面，kubelet 作为节点状况的监视器，它获取节点信息，并以集群客户端的角色，把这些状况同步到 API Server。在这个问题中，kubelet 扮演的是第二种角色。Kubelet 会使用上图中的 NodeStatus 机制，定期检查集群节点状况，并把节点状况同步到 API Server。而 NodeStatus 判断节点就绪状况的一个主要依据，就是 PLEG。

PLEG 是 Pod Lifecycle Events Generator 的缩写，基本上它的执行逻辑，是定期检查节点上 Pod 运行情况，如果发现感兴趣的变化，PLEG 就会把这种变化包装成 Event 发送给 Kubelet 的主同步机制 syncLoop 去处理。但是，在 PLEG 的 Pod 检查机制不能定期执行的时候，NodeStatus 机制就会认为，这个节点的状况是不对的，从而把这种状况同步到 API Server。

整体的工作流程如下图所示，虚线部分是 PLEG 的工作内容。![](https://notes-learning.oss-cn-beijing.aliyuncs.com/35160200-5f9f-4cfd-911c-afa960062a5c/640)

### 以 node notready 这个场景为例，来讲解 PLEG：

Kubelet 中的 NodeStatus 机制会定期检查集群节点状况，并把节点状况同步到 API Server。而 NodeStatus 判断节点就绪状况的一个主要依据，就是 PLEG。

PLEG 定期检查节点上 Pod 运行情况，并且会把 pod 的变化包装成 Event 发送给 Kubelet 的主同步机制 syncLoop 去处理。但是，在 PLEG 的 Pod 检查机制不能定期执行的时候，NodeStatus 机制就会认为这个节点的状况是不对的，从而把这种状况同步到 API Server，我们就会看到 not ready 。

PLEG 有两个关键的时间参数，一个是检查的执行间隔，另外一个是检查的超时时间。以默认情况为准，PLEG 检查会间隔一秒，换句话说，每一次检查过程执行之后，PLEG 会等待一秒钟，然后进行下一次检查；而每一次检查的超时时间是三分钟，如果一次 PLEG 检查操作不能在三分钟内完成，那么这个状况，会被 NodeStatus 机制当做集群节点 NotReady 的凭据，同步给 API Server。

PLEG Start 就是启动一个协程，每个 relistPeriod(1s) 就调用一次 relist，根据最新的 PodStatus 生成 PodLiftCycleEvent。relist 是 PLEG 的核心，它从 container runtime 中查询属于 kubelet 管理 containers/sandboxes 的信息，并与自身维护的 pods cache 信息进行对比，生成对应的 PodLifecycleEvent，然后输出到 eventChannel 中，通过 eventChannel 发送到 kubelet syncLoop 进行消费，然后由 kubelet syncPod 来触发 pod 同步处理过程，最终达到用户的期望状态。

### PLEG is not healthy 的原因

这个报错清楚地告诉我们，容器 runtime 是不正常的，且 PLEG 是不健康的。这里容器 runtime 指的就是 docker daemon 。Kubelet 通过操作 docker daemon 来控制容器的生命周期。而这里的 PLEG，指的是 pod lifecycle event generator。PLEG 是 kubelet 用来检查 runtime 的健康检查机制。这件事情本来可以由 kubelet 使用 polling 的方式来做。但是 polling 有其高成本的缺陷，所以 PLEG 应用而生。PLEG 尝试以一种 “中断” 的形式，来实现对容器 runtime 的健康检查，虽然实际上，它同时用了 polling 和”中断”这样折中的方案。

从 Docker 1.11 版本开始，Docker 容器运行就不是简单通过 Docker Daemon 来启动了，而是通过集成 containerd、runc 等多个组件来完成的。虽然 Docker Daemon 守护进程模块在不停的重构，但是基本功能和定位没有太大的变化，一直都是 CS 架构，守护进程负责和 Docker Client 端交互，并管理 Docker 镜像和容器。现在的架构中组件 containerd 就会负责集群节点上容器的生命周期管理，并向上为 Docker Daemon 提供 gRPC 接口。![](https://notes-learning.oss-cn-beijing.aliyuncs.com/35160200-5f9f-4cfd-911c-afa960062a5c/640)

PLEG 在每次迭代检查中会调用 runc 的 relist() 函数干的事情，是定期重新列出节点上的所有容器，并与上一次的容器列表进行对比，以此来判断容器状态的变换。相当于 docker ps 来获取所有容器，在通过 docker Inspect 来获取这些容器的详细信息。在有问题的节点上，通过 docker ps 命令会没有响应，这说明上边的报错是准确的。

### 经常出现的场景

出现 pleg not healthy，一般有以下几种可能：

- 容器运行时无响应或响应超时，如 docker 进程响应超时（比较常见）
- 该节点上容器数量过多，导致 relist 的过程无法在 3 分钟内完成
- relist 出现了死锁，该 bug 已在 Kubernetes 1.14 中修复。
- 网络

### 排查处理过程描述

1. 我们在问题节点上执行 top，发现有进程名为 scope 的进程 cpu 占用率一直是 100%。通过翻阅资料得知 systemd.scope：范围 (scope) 单元并不通过单元文件进行配置， 而是仅能以编程的方式通过 systemd D-Bus 接口创建。范围单元的名称都以 ".scope" 作为后缀。与服务 (service) 单元不同，范围单元用于管理 一组外部创建的进程， 它自身并不派生 (fork) 任何进程。范围 (scope) 单元的主要目的在于以分组的方式管理系统服务的工作进程。2. 在继续执行在有问题的节点上，通过 docker ps 命令会没有响应。说明容器 runtime 也是有问题的。那容器 runtime 与 systemd 有不有关系呢？3. 我们通过查阅到阿里的一篇文章，阿里巴巴 Kubernetes 集群问题排查思路和方法。找到了关系，有兴趣的可以根据文末提供的链接去细致了解。以下是在该文章中截取的部分内容。

#### 什么是 D-Bus 呢？

通过阿里巴巴 Kubernetes 集群问题排查思路和方法\[1]中如下描述：在 Linux 上，dbus 是一种进程间进行消息通信的机制。![](https://notes-learning.oss-cn-beijing.aliyuncs.com/35160200-5f9f-4cfd-911c-afa960062a5c/640)

#### RunC 请求 D-Bus

容器 runtime 的 runC 命令，是 libcontainer 的一个简单的封装。这个工具可以用来管理单个容器，比如容器创建和容器删除。在上节的最后，我们发现 runC 不能完成创建容器的任务。我们可以把对应的进程杀掉，然后在命令行用同样的命令启动容器，同时用 strace 追踪整个过程。![](https://notes-learning.oss-cn-beijing.aliyuncs.com/35160200-5f9f-4cfd-911c-afa960062a5c/640)
分析发现，runC 停在了向带有 org.free 字段的 dbus socket 写数据的地方。

#### 解决问题

最后可以断定是 systemd 的问题，我们用 systemctl daemon-reexec 来重启 systemd，问题消失了。所以更加确定是 systemd 的问题。

具体原因大家可以参考：[https://www.infoq.cn/article/t_ZQeWjJLGWGT8BmmiU4 这篇文章。](https://www.infoq.cn/article/t_ZQeWjJLGWGT8BmmiU4%E8%BF%99%E7%AF%87%E6%96%87%E7%AB%A0%E3%80%82)

根本上解决问题是：将 systemd 升级到 v242-rc2，升级后需要重启操作系统。（[https://github.com/lnykryn/systemd-rhel/pull/322）](https://github.com/lnykryn/systemd-rhel/pull/322%EF%BC%89)

## 总结

PLEG is not healthy 的问题居然是因为 systemd 导致的。最后通过将 systemd 升级到 v242-rc2，升级后需要重启操作系统。（[https://github.com/lnykryn/systemd-rhel/pull/322）](https://github.com/lnykryn/systemd-rhel/pull/322%EF%BC%89) 参考资料

- Kubelet: Pod Lifecycle Event Generator (PLEG)
- Kubelet: Runtime Pod Cache
- relist() in kubernetes/pkg/kubelet/pleg/generic.go
- Past bug about CNI — PLEG is not healthy error, node marked NotReady
- <https://www.infoq.cn/article/t_ZQeWjJLGWGT8BmmiU4>
- <https://cloud.tencent.com/developer/article/1550038>
