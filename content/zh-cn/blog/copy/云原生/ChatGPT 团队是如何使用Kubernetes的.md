---
title: "ChatGPT 团队是如何使用Kubernetes的"
linkTitle: "ChatGPT 团队是如何使用Kubernetes的"
weight: 20
---

[原文链接](https://mp.weixin.qq.com/s/u7zibC7UmMSAYotWZ81eYg)

https://openai.com/research/scaling-kubernetes-to-7500-nodes

在本文中，OpenAI 的工程师团队分享了他们在 Kubernetes 集群扩展过程中遇到的各种挑战和解决方案，以及他们取得的性能和效果。

我们已经将 Kubernetes 集群扩展到 7500 个节点，为大型模型（如 GPT-3、 CLIP 和 DALL·E）创建了可扩展的基础设施，同时也为快速小规模迭代研究（如 神经语言模型的缩放定律）创建了可扩展的基础设施。

将单个 Kubernetes 集群扩展到这种规模很少见，但好处是能够提供一个简单的基础架构，使我们的机器学习研究团队能够更快地推进并扩展，而无需更改代码。

自上次发布关于扩展到 2500 个节点的帖子以来，我们继续扩大基础设施以满足研究人员的需求，在此过程中学到了许多的经验教训。本文总结了这些经验教训，以便 Kubernetes 社区里的其他人也能从中受益，并最后会介绍下我们仍然面临的问题，我们也将继续解决这些问题。

![](https://mmbiz.qpic.cn/mmbiz_png/YriaiaJPb26VOiciaV2ibeVb4gXQtYocv76XKdUh15yLPEl9eIzUgtNqcpB4tp85WFRKYQGYlv5dlPjib6OTbqCqeUdg/640?wx_fmt=png&wxfrom=13&tp=wxpic)

我们的工作负载

在深入探讨之前，我们着重描述一下我们的工作负载。我们在 Kubernetes 上运行的应用程序和硬件与大家在普通公司遇到的可能相当不同。因此，我们的问题及解决方案可能与你自己的设置匹配，也可能不匹配！

一个大型的机器学习作业跨越许多节点，当它可以访问每个节点上的所有硬件资源时，运行效率最高。这允许 GPU 直接使用 NVLink 进行交叉通信，或者 GPU 使用 GPUDirect 直接与 NIC 进行通信。因此，对于我们的许多工作负载，单个 Pod 占用整个节点。任何 NUMA、CPU 或 PCIE 资源争用都不是调度的因素。装箱或碎片化不是常见的问题。我们目前的集群具有完全的二分带宽，因此我们也不考虑机架或网络拓扑。所有这些都意味着，虽然我们有许多节点，但调度程序的负载相对较低。

话虽如此，kube-scheduler 的负载是有波动的。一个新的作业可能由许多数百个 Pod 同时创建组成，然后返回到相对较低的流失率。

![](https://mmbiz.qpic.cn/mmbiz_png/YriaiaJPb26VOiciaV2ibeVb4gXQtYocv76XKLpfylgOVT7BzEulB0dicicTPY64JIp4CzqozxGqiaibbxiawSQliaFeicVhWA/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1&tp=wxpic)

我们最大的作业运行 MPI，作业中的所有 Pod 都参与一个单一的 MPI 通信器。如果任何一个参与的 Pod 挂掉，整个作业就会停止，需要重新启动。作业会定期进行检查点，当重新启动时，它会从上一个检查点恢复。因此，我们认为 Pod 是半有状态的——被删掉的 Pod 可以被替换并且工作可以继续，但这样做会造成干扰，应该尽量减少发生。

我们并不太依赖 Kubernetes 的负载均衡。我们的 HTTPS 流量非常少，不需要进行 A/B 测试、蓝 / 绿或金丝雀部署。Pod 使用 SSH 直接通过 Pod IP 地址与 MPI 进行通信，而不是通过服务端点。服务“发现”是有限的；我们只在作业启动时进行一次查找，查找哪些 Pod 参与 MPI。

大多数作业与某种形式的 Blob 存储进行交互。它们通常会直接从 Blob 存储流式传输一些数据集的分片或检查点，或将其缓存到快速的本地临时磁盘中。我们有一些 PersistentVolumes，用于那些需要 POSIX 语义的情况，但 Blob 存储更具可扩展性，而且不需要缓慢的分离 / 附加操作。

最后，我们的工作性质本质上是研究，这意味着工作负载本身是不断变化的。虽然超级计算团队努力提供我们认为达到“生产”质量水平的计算基础架构，但在该集群上运行的应用程序寿命很短，它们的开发人员会快速迭代。因此，随时可能出现新的使用模式，这些模式会挑战我们对趋势和适当权衡的设定。我们需要一个可持续的系统，同时也能让我们在事情发生变化时快速做出响应。

网   络

随着集群内节点和 Pod 数量的增加，我们发现 Flannel 难以满足所需的吞吐量。因此，我们转而使用 Azure VMSS 的本地 Pod 网络技术和相关 CNI 插件来配置 IP。这使我们的 Pod 能够获得主机级别的网络吞吐量。

我们转而使用别名 IP 地址的另一个原因是，在我们最大的集群中，可能会同时使用约 20 万个 IP 地址。在测试了基于路由的 Pod 网络后，我们发现能够使用的路由数明显存在限制。

避免封装会增加底层 SDN 或路由引擎的需求，虽然这使我们的网络设置变得简单。添加 VPN 或隧道可以在不需要任何其他适配器的情况下完成。我们不需要担心由于某部分网络具有较低的 MTU 而导致的分组分段。网络策略和流量监控很简单；没有关于数据包源和目的地的歧义。

我们使用主机上的 iptables 标记来跟踪每个 Namespace 和 Pod 的网络资源使用情况，这使研究人员可以可视化他们的网络使用模式。特别是，由于我们的许多实验具有不同的 Internet 和 Pod 内通信模式，因此能够调查任何瓶颈发生的位置通常是非常有意义的。

可以使用 iptables `mangle` 规则任意标记符合特定条件的数据包。以下是我们用来检测流量是内部流量还是 Internet 流量的规则。`FORWARD` 规则涵盖了来自 Pod 的流量，而 `INPUT` 和 `OUTPUT` 规则涵盖了主机上的流量：

```nginx
iptables -t mangle -A INPUT ! -s 10.0.0.0/8 -m comment --comment "iptables-exporter openai traffic=internet-in"
iptables -t mangle -A FORWARD ! -s 10.0.0.0/8 -m comment --comment "iptables-exporter openai traffic=internet-in"
iptables -t mangle -A OUTPUT ! -d 10.0.0.0/8 -m comment --comment "iptables-exporter openai traffic=internet-out"
iptables -t mangle -A FORWARD ! -d 10.0.0.0/8 -m comment --comment "iptables-exporter openai traffic=internet-out"
```

一旦标记，iptables 将开始计数以跟踪匹配该规则的字节数和数据包数。你可以使用 `iptables` 本身来查看这些计数器：

```bash
% iptables -t mangle -L -v
Chain FORWARD (policy ACCEPT 50M packets, 334G bytes)
 pkts bytes target     prot opt in     out     source               destination
....
1253K  555M            all  --  any    any     anywhere            !10.0.0.0/8           /* iptables-exporter openai traffic=internet-out */
1161K 7937M            all  --  any    any    !10.0.0.0/8           anywhere             /* iptables-exporter openai traffic=internet-in */
```

我们使用名为 iptables-exporter 的开源 Prometheus 导出器将这些数据追踪到我们的监控系统中。这是一种简单的方法，可以跟踪与各种不同类型的条件匹配的数据包。

![](https://mmbiz.qpic.cn/mmbiz_png/YriaiaJPb26VOiciaV2ibeVb4gXQtYocv76XKd8NicaEsHsoGOvcBJMt435PARbypib0ENCBwpeXTibnrp9mkYENFG0Hfg/640?wx_fmt=png)

![](https://mmbiz.qpic.cn/mmbiz_png/YriaiaJPb26VOiciaV2ibeVb4gXQtYocv76XKQOOXeibMYDia4Eoic4E7iaNI2uPCa1HOVeY75vj4D2SfSnGVm5wVmDx9Ug/640?wx_fmt=png)

我们网络模型中比较独特的一点是，我们完全向研究人员公开节点、Pod 和 Service 网络 CIDR 范围。我们采用集线器和分支的网络模型，并使用本机节点和 Pod CIDR 范围路由该流量。研究人员连接到中心枢纽，然后可以访问任何一个单独的集群（分支）。但是这些集群本身无法相互通信。这确保了集群保持隔离、没有跨集群依赖，可以防止故障隔离中的故障传播。

我们使用一个“NAT”主机来翻译从集群外部传入的服务网络 CIDR 范围的流量。这种设置为我们的研究人员提供了很大的灵活性，他们可以选择各种不同类型的网络配置进行实验。

API 服务器

Kubernetes 的 API Server 和 etcd 是保持集群健康运行的关键组件，因此我们特别关注这些系统的压力。我们使用 kube-prometheus 提供的 Grafana 仪表板以及额外的内部仪表板。我们发现，将 HTTP 状态码 429（请求太多）和 5xx（服务器错误）的速率作为高级信号警报是有用的。

![](https://mmbiz.qpic.cn/mmbiz_png/YriaiaJPb26VOiciaV2ibeVb4gXQtYocv76XKDtNcvSjnOqmZ6iccOSmwmricspB55500xvHpBMAI6HiaGfcCtSjqq3ljA/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1&tp=wxpic)

虽然有些人在 kube 内部运行 API 服务器，但我们一直在集群外运行它们。etcd 和 API 服务器都在它们自己的专用节点上运行。我们的最大集群运行 5 个 API 服务器和 5 个 etcd 节点，以分散负载并尽可能减少发生故障后带来的影响。自从我们在 上一篇博文 中提到的将 Kubernetes 事件拆分到它们自己的 etcd 集群中以来，我们没有遇到 etcd 的任何值得注意的问题。API 服务器是无状态的，通常很容易在自我修复的实例组或扩展集中运行。我们尚未尝试构建任何自我修复 etcd 集群的自动化，因为发生事故非常罕见。

API 服务器可能会占用相当多的内存，并且往往会与集群中的节点数量成线性比例。对于我们有 7500 个节点的集群，我们观察到每个 API 服务器使用高达 70GB 的堆内存，因此幸运地是，未来这应该仍然在硬件能力范围之内。

![](https://mmbiz.qpic.cn/mmbiz_png/YriaiaJPb26VOiciaV2ibeVb4gXQtYocv76XKh2rebxQzCsojXTUN9nDTDZbYWPMpq9ORPlVQfzHwh4sbM4VCztKibWQ/640?wx_fmt=png&tp=wxpic&wxfrom=5&wx_lazy=1&wx_co=1)

API Servers 受到压力的主要来源之一就是对 Endpoints 的 WATCH。在整个集群中有一些服务，如“kubelet”和“node-exporter”，其中每个节点都是成员。当一个节点被添加或从集群中删除时，这个 WATCH 将被触发。通常，由于每个节点本身都通过 kube-proxy 监视 `kubelet` 服务，因此这些响应中所需的数量和带宽将是 N2 非常大，有时会达到 1GB/s 或更高。Kubernetes 1.17 中推出的 EndpointSlices 大大降低了这种负载，减少达 1000 倍。

![](https://mmbiz.qpic.cn/mmbiz_png/YriaiaJPb26VOiciaV2ibeVb4gXQtYocv76XKwwo1NJytItYTcNF6oeOvs5TeLaBNqZuJa7LKFysickwuoBcHEiacZltQ/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1&tp=wxpic)

总的来说，我们会非常注意随着集群规模增大而增加的 API Server 请求。我们尽量避免任何 DaemonSets 与 API Server 交互。在需要每个节点监视更改的情况下，引入缓存服务（例如 Datadog Cluster Agent）作为中介，似乎是避免集群范围瓶颈的良好模式。

随着集群的增长，我们对集群的实际自动伸缩越来越少。但当一次性自动扩展太多时，我们偶尔会遇到问题。当新节点加入集群时会生成大量请求，一次性添加数百个节点可能会超过 API 服务器容量的负荷。稍微平滑一下这个过程，即使只有几秒钟也有助于避免宕机。

时间序列度量与 Prometheus 和 Grafana

我们使用 Prometheus 收集时间序列度量数据，并使用 Grafana 进行图形、仪表板和警报。我们从 kube-prometheus 部署开始收集了各种各样的度量数据，并使用了一些良好的仪表板进行可视化。随着节点数量的不断增加，我们开始遇到 Prometheus 收集的度量数据数量过多的问题。尽管 kube-prometheus 公开了许多有用的数据，但我们实际上并没有查看所有的度量数据，一些数据也过于细化，无法有效地进行收集、存储和查询。因此，我们使用 Prometheus 规则从被摄入的度量数据中“删掉”一些数据。

有一段时间，我们遇到了 Prometheus 消耗越来越多的内存问题，最终导致容器崩溃并出现 Out-Of-Memory 错误（OOM）。即使为应用程序分配了大量的内存容量，这种情况似乎仍然会发生。更糟糕的是，它在崩溃时会花费很多时间在启动时回放预写日志文件，直到它再次可用。最终，我们发现了这些 OOM 的来源是 Grafana 和 Prometheus 之间的交互，其中 Grafana 使用 `/api/v1/series` API 查询 `{le!=""}`（基本上是“给我所有的直方图度量”）。`/api/v1/series` 的实现在时间和空间上没有限制，对于具有大量结果的查询，这将不断消耗更多的内存和时间。即使请求者已经放弃并关闭了连接，它也会继续增长。对于我们来说，内存永远不够，而 Prometheus 最终会崩溃。因此，我们修补了 Prometheus，将此 API 包含在上下文中以强制执行超时，从而完全解决了问题。

虽然 Prometheus 崩溃的次数大大减少，但在我们需要重新启动它的时候，WAL 回放仍然是一个问题。通常需要多个小时来回放所有 WAL 日志，直到 Prometheus 开始收集新的度量数据并提供服务。在 Robust Perception 的帮助下，我们发现将 `GOMAXPROCS=24` 应用于服务器可以显著提高性能。在 WAL 回放期间，Prometheus 尝试使用所有核心，并且对于具有大量核心的服务器，争用会降低所有性能。

我们正在探索新的选项来增加我们的监控能力，下面“未解决的问题”部分将对此进行描述。

健康检查

对于如此庞大的集群，我们当然依赖自动化来检测并从集群中移除行为不当的节点。随着时间的推移，我们建立了许多健康检查系统。

被动健康检查

某些健康检查是被动的，总是在所有节点上运行。这些检查监视基本的系统资源，例如网络可达性、坏盘或满盘，或者 GPU 错误。GPU 以许多不同的方式出现问题，但一个容易出现的常见问题是“不可纠正的 ECC 错误”。Nvidia 的数据中心 GPU 管理器（DCGM）工具使查询这个问题和许多其他“Xid”错误变得容易。我们跟踪这些错误的一种方式是通过 dcgm-exporter 将指标收集到我们的监控系统 Prometheus 中。这将出现为 `DCGM_FI_DEV_XID_ERRORS` 指标，并设置为最近发生的错误代码。此外，NVML 设备查询 API 公开了有关 GPU 的健康和操作的更详细信息。

一旦检测到错误，它们通常可以通过重置 GPU 或系统来修复，但在某些情况下确实需要更换基础 GPU。

另一种健康检查是跟踪来自上游云提供商的维护事件。每个主要的云提供商都公开了一种方式来了解当前 VM 是否需要进行会最终导致中断的、即将发生的维护事件。VM 可能需要重新启动以应用底层的超级管理程序补丁，或者将物理节点替换为其他硬件。

这些被动健康检查在所有节点上不断运行。如果健康检查开始失败，节点将自动划分，因此不会在节点上安排新的 Pod。对于更严重的健康检查失败，我们还将尝试 Pod 驱逐，以要求当前运行的所有 Pod 立即退出。这仍然取决于 Pod 本身，可通过 Pod 故障预算进行配置来决定是否允许此驱逐发生。最终，无论是在所有 Pod 终止之后，还是在 7 天过去之后（我们的服务级别协议的一部分），我们都将强制终止 VM。

活动 GPU 测试

不幸的是，并非所有 GPU 问题都会通过 DCGM 可见的错误代码表现出来。我们建立了自己的测试库，通过对 GPU 进行测试来捕捉其他问题，并确保硬件和驱动程序的行为符合预期。这些测试无法在后台运行 - 它们需要独占 GPU 运行数秒钟或数分钟。

我们首先在节点启动时运行这些测试，使用我们称之为“预检（preflight）”的系统。所有节点都会附带一个“预检”污点和标签加入集群。这个污点会阻止普通 Pod 被调度到节点上。我们配置了一个 DaemonSet，在所有带有此标签的节点上运行预检测试 Pod。测试成功完成后，测试本身将删除污点和标签，然后该节点就可供一般使用。

我们还定期在节点的生命周期中运行这些测试。我们将其作为 CronJob 运行，允许它着陆在集群中的任何可用节点上。哪些节点会被测试到可能有些随机和不受控制，但我们发现随着时间的推移，它提供了足够的覆盖率，并且最小化了协调或干扰。

配额和资源使用

随着集群规模的扩大，研究人员开始发现他们难以获取分配给他们的全部容量。传统的作业调度系统有许多不同的功能，可以公平地在竞争团队之间运行工作，而 Kubernetes 没有这些功能。随着时间的推移，我们从这些作业调度系统中汲取灵感，并以 Kubernetes 原生的方式构建了几个功能。

团队污点

我们在每个集群中都有一个服务，称为“team-resource-manager”，具有多个功能。它的数据源是一个 ConfigMap，为在给定集群中具有容量的所有研究团队指定了 （节点选择器、应用的团队标签、分配数量） 元组。它会将当前节点与这些元组进行对比，并使用 `openai.com/team=teamname:NoSchedule` 的污点对适当数量的节点进行标记。

“team-resource-manager”还有一个入站的 webhook 服务，因此在提交每个作业时会根据提交者的团队成员身份应用相应的容忍度。使用污点使我们能够灵活地限制 Kubernetes Pod 调度程序，例如允许较低优先级的 Pod 具有 "any" 容忍度，这样团队可以借用彼此的容量，而无需进行大量协调。

CPU 和 GPU  balloons

除了使用集群自动缩放器动态扩展我们基于虚拟机的集群之外，我们还使用它来纠正（删除和重新添加）集群中的不健康成员。我们通过将集群的 "最小值" 设置为零、"最大值" 设置为可用容量来实现这一点。然而，如果 cluster-autoscaler 发现有空闲节点，它将尝试缩小到只需要的容量。由于多种原因（VM 启动延迟、预分配成本、上面提到的 API 服务器影响），这种空闲缩放并不理想。

因此，我们为 CPU 和 GPU 主机都引入了“球形”部署。这个部署包含一个具有 "最大值" 数量的低优先级 Pod 副本集。这些 Pod 占用节点内的资源，因此自动缩放器不会将它们视为空闲。但由于它们是低优先级的，调度程序可以立即将它们驱逐出去，以腾出空间进行实际工作。（我们选择使用 Deployment 而不是 DaemonSet，以避免将 DaemonSet 视为节点上的空闲工作负载。）

需要注意的是，我们使用 pod 反亲和性（anti-affinity）来确保 pod 在节点之间均匀分布。Kubernetes 调度器的早期版本存在一个 O(N^2) 的性能问题，与 pod 反亲和性有关。自 Kubernetes 1.18 版本以后，这个问题已经得到了纠正。

Gang 调度

我们的实验通常涉及一个或多个 StatefulSets，每个 StatefulSet 操作不同部分的训练任务。对于优化器，研究人员需要在进行任何训练之前调度 StatefulSet 的所有成员（因为我们通常使用 MPI 在优化器成员之间协调，而 MPI 对组成员变化很敏感）。

然而再默认情况下，Kubernetes 不一定会优先满足某个 StatefulSet 的所有请求。例如，如果两个实验都请求 100％的集群容量，那么 Kubernetes 可能只会调度给每个实验需要的一半 Pod，这会导致死锁，使两个实验都无法进行。

我们尝试了一些需要自定义调度程序的方法，但遇到了一些与正常 Pod 调度方式冲突的边缘情况。Kubernetes 1.18 引入了核心 Kubernetes 调度程序的插件体系结构，使本地添加此类功能变得更加容易。我们最近选择了 Coscheduling 插件作为解决此问题的方法。

未解决的问题

随着 Kubernetes 集群规模的扩大，我们仍有许多问题需要解决。其中一些问题包括：

指标

在如今的规模下，Prometheus 内置的 TSDB 存储引擎很难压缩，并且每次重新启动时需要长时间回放 WAL（预写式日志）。查询还往往会导致“查询处理会加载过多样本”的错误。我们正在迁移到不同的、与 Prometheus 兼容的存储和查询引擎。大家可以期待下我们未来的博客文章，看看它的表现如何！

Pod 网络流量整形

随着集群规模的扩大，每个 Pod 的互联网带宽量被计算了出来。每个人的聚合互联网带宽需求变得非常大，我们的研究人员现在有能力会意外地对互联网上的其他位置施加重大资源压力，例如要下载的数据集和要安装的软件包。

结   论

Kubernetes 是一个非常灵活的平台，可以满足我们的研究需求。它具有满足我们所面临的最苛刻工作负载的能力。尽管它仍有许多需要改进的地方，但 OpenAI 的超级计算团队将继续探索 Kubernetes 的可扩展性。
