---
title: CRI 对比
---

原文：[阳明公众号](https://mp.weixin.qq.com/s/H3vUUvEiOfLkd_YEoo8sNg)

下面是我已经测试的几个 CRI，并进行一些基准测试来对他们进行了简单的对比，希望对你有所帮助：

- dockershim
- containerd
- crio

对于 cri-o，已经测试了 2 个后端：runc 和 crun，以测试对 `cgroupsv2` 的影响。

## 测试环境

我这里的测试环境是一个 1.19.4 版本的 kubernetes 集群，使用 ansible 进行创建（<https://gitlab.com/incubateur-pe>）。集群运行在 kvm 上，配置如下：

- master：Centos/7, 2vcpus/2G 内存。
- crio-crun 节点：Fedora-32, 2vcpus/4G 内存。
- 其他节点：Centos/7, 2vcpus/4G 内存.

底层是 i7-9700K ，64G 的内存和一个 mp510 nvme 硬盘。

## 创建集群

这里我直接使用 molecule 创建一个集群，并配置了它在每个 worker 节点上使用不同的 cri，对应的 ansible 源码位于：<https://gitlab.com/incubateur-pe/kubernetes-bare-metal/-/tree/dev/molecule/criBench>

使用上面的脚本，执行 `molecule converge` 命令后，大概 10 分钟左右，我们就可以得到一个如下所示的 kubernetes 集群。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/czfhxy/1616120034837-61ee3a87-ccca-4c3a-91ff-983ad9808b6d.png)

接下来我们就可以进行一些简单的基准测试了。

## 测试

### 1. bucketbench 测试

Bucketbench (https://github.com/estesp/bucketbench) 是一个可以对容器引擎执行一系列操作的测试工具，它非常适合于了解之前每个节点的性能。

这里我们的测试参数很简单：

- 3 个线程
- 15 次循环
- run/stop/delete 操作

对应的结果如下所示（ms 为单位）：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/czfhxy/1616120034822-55f1f7a0-b8cf-422b-9490-f0a7b5f8593d.png)

我们可以看到在性能上还是有相当大的差异的。但是需要注意的是我们这里为什么测试了 5 个实例呢？上面不是只有 4 个 worker 节点吗？

这里其实是因为这里我们使用的 docker 客户端并不是 kubernetes 使用的，事实上 docker 实现了 CRI，并提供了一个 socket，这个 socket 和其他 cri socket 一样可以调用。所以这里的区别是：

- docker-shim：是通过 cri 的 socket 来做测试
- docker-cli：是通过 docker 客户端来做测试

但是实际上 docker 并没有想象中那么差，在这个测试中我们可以看到他比 cri-o 要快点，当然这个测试中很明显 dockerd 是表现最好的。

### 2. kubernetes 测试

上面的测试并不能完整说明这几个 cri 之间的差距，当它们被 kubernetes 使用的时候，它们表现又如何呢？是否不止 `run/stop/delete` 这些操作？性能上的差异在真正的集群上又有什么意义吗?

下面我们就来深入了解下，这次我们使用集群中的 Prometheus、Grafana 来可视化监控指标，对应的自定义 dashboard 数据可以在 <https://gitlab.com/ulrich.giraud/bench-cri/-/blob/master/dashboard/dashboard_bench.json> 这里获取。由于只是测试容器运行时，不是工作负载，所以这里我们只是简单的在集群中部署的一个 busybox 镜像并一直 sleep 的 DaemonSet 应用。

    apiVersion: apps/v1kind: DaemonSetmetadata:  name: benchds-replaceme  namespace: benchds  labels:    k8s-app: benchdsspec:  selector:    matchLabels:      name: benchds  template:    metadata:      labels:        name: benchds    spec:      containers:      - name: benchds        image: busybox:latest        command:          - sleep          - infinity        resources:          limits:            memory: 20Mi          requests:            cpu: 10m            memory: 20Mi

该 DamonSet 将用唯一的名称进行部署：

- 100 次（两次创建之间有一定延迟）
- 批量 100 次
- 批量 1000 次

对应的 Grafana 展示图表信息如下所示：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/czfhxy/1616120034839-e6353428-f475-4c08-9be7-73a5d7b656c2.png)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/czfhxy/1616120034856-b41bdfbc-a4b4-477a-914c-95f009b88193.png)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/czfhxy/1616120034825-ef0d77c9-202f-45f6-aeb5-6282f10214e7.png)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/czfhxy/1616120034888-e1840303-2907-45d0-a780-05119ec475eb.png)

缓慢创建数百个 DaemonSets

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/czfhxy/1616120034850-86008be9-f893-47c9-a20b-d18523e12d20.png)

快速创建数百个 DaemonSets

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/czfhxy/1616120034854-a1a537e5-e71e-4313-aec4-42265497e49d.png)

快速创建数千个 DaemonSets

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/czfhxy/1616120034848-339697b7-b147-4af4-aae5-08728a93332c.png)

现在我们来分析下上面的测试结果。

- Cri-o/runc：令人惊讶的是，在所有 create/delete 中是最慢的，但在其他方面处于中等水平。
- Cri-o/crun：在 create/delete 方面不是很好，但是在其他方面表现是最好的。
- Containerd：表现非常好，几乎在所有情况下都可以快速响应。
- Docker：在 create/delete 方面比 cri-o 快，但在 status/list 请求方面是最慢的。

status/list 请求是 cri 上最频繁的请求，所以这也是性能最重要的地方，cri-o 在这里似乎是更好的选择，其次就是 containerd。

containerd 在所有指标上的表现都比较好，应该是最均衡的一个选择了。另外一方面，docker 并没有得到很好的测试结果，但是无论负载情况如何，它的表现基本上都是一致的。

## 总结

从纯性能角度来说，确实有比 docker 更好的替代品，我们的集群也不会替换 docker 产生什么影响。从另外一个角度来看，kubernetes 这次废弃 docker 的事情也算是一件好事，让更多的人意识到 docker 并不是唯一可用的 CRI，甚至不是唯一的构建镜像工具。

在我看来，docker 仍然是让整个容器化向前发展的一个伟大工具。但是好像我还没有回答我最初的问题，那就是：我应该为我的 k8s 集群使用什么 CRI？

从我个人角度考虑的话，我个人的选择是：containerd，他速度快，配置方便，相当可靠和安全，不过 cri-o 已经支持 cgroupsv2 了，所以如果我使用 fedora 或者 centos/8 的话我会优先选择 cri-o。

> 原文链接：<https://ulrich-giraud.medium.com/which-cri-should-i-use-to-replace-docker-for-my-kubernetes-cluster-14a45c080004>
