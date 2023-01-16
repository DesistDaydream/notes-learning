---
title: 深入理解Kubernetes Operator
---

#

James Laverack InfoQ _11 月 18 日_

作者 | James Laverack 译者 | 王者**本文要点：**

- Kubernetes API 为所有云资源提供了单个集成点，以此来促进云原生技术的采用。

- 有一些框架和库可以用来简化 Operator 的编写。支持多种语言，其中 Go 生态系统是最为成熟的。

- 你可以为非自有的软件创建 Operator。DevOps 团队可能会通过这种方式来管理数据库或其他外部产品。

- 难点不在于 Operator 本身，而是要学会理解它的行为。

多年来，Operator 一直是 Kubernetes 生态系统的重要组成部分。通过将管理界面移动到 Kubneretes API 中，带来了“单层玻璃”的体验。对于希望简化 kuberentes 原生应用程序的开发人员或者希望降低现有系统复杂性的 DevOps 工程师来说，Operator 可能是一个非常有吸引力的选择。但如何从头开始创建一个 Operator 呢？

深入理解 OperatorOperator 是什么？

如今，Operator 无处不在。数据库、云原生项目、任何需要在 Kubernetes 上部署或维护的复杂项目都用到了 Operator。CoreOS 在 2016 年首次引入了 Operator，将运维关注点转移到软件系统中。Operator 自动执行操作，例如，Operator 可以部署数据库实例、升级数据库版本或执行备份。然后，这些系统可以被测试，响应速度比人类工程师更快。

Operator 还通过使用自定义资源定义对 Kubenretes API 进行了扩展，将工具配置转移到了 API 中。这意味着 Kubenretes 本身就变成了“单层玻璃”。DevOps 工程师可以利用围绕 Kubernetes API 资源而构建的工具生态系统来管理和监控他们部署的应用程序：

- 使用 Kubernetes 内置的基于角色的访问控制 (RBAC) 来修改授权和身份验证。

- 使用“git ops”对生产变更进行可复制的部署和代码审查。

- 使用基于开放策略代理 (OPA) 的安全工具在自定义资源上应用策略。

- 使用 Helm、Kustomize、ksonnet 和 Terraform 等工具简化部署描述。

这种方法还可以确保生产、测试和开发环境之间的一致性。如果每个集群都是 Kubernetes 集群，则可以使用 Operator 在每个集群中部署相同的配置。

为什么要使用 Operator？

使用 Operator 有很多理由。通常情况下，要么是开发团队为他们的产品创建 Operator，要么是 DevOps 团队希望对第三方软件管理进行自动化。无论哪种方式，都应该从确定 Operator 应该负责哪些东西开始。

最基本的 Operator 用于部署，使用 kubectl apply 就可以创建一个用于响应 API 资源的数据库，但这比内置的 Kubernetes 资源 (如 StatefulSets 或 Deployments) 好不了多少。复杂的 Operator 将提供更大的价值。如果你想要对数据库进行伸缩该怎么办？

如果是 StatefulSet，你可以执行 kubectl scale statefulset my-db --replicas 3，这样就可以得到 3 个实例。但如果这些实例需要不同的配置呢？是否需要指定一个实例为主实例，其他实例为副本？如果在添加新副本之前需要执行设置步骤，那该怎么办？在这种情况下，可以使用 Operator。

更高级的 Operator 可以处理其他一些特性，如响应负载的自动伸缩、备份和恢复、与 Prometheus 等度量系统的集成，甚至可以进行故障检测和自动调优。任何具有传统“运行手册”文档的操作都可以被自动化、测试和依赖，并自动做出响应。

被管理的系统甚至不需要部署在 Kubernetes 上也能从 Operator 中获益。例如，主要的云服务提供商（如 Amazon Web Services、微软 Azure 和谷歌云）提供 Kubenretes Operator 来管理其他云资源，如对象存储。用户可以通过配置 Kubernetes 应用程序的方式来配置云资源。运维团队可能对其他资源也采取同样的方法，使用 Operator 来管理任何东西——从第三方软件服务到硬件。

Operator 示例

在本文中，我们将重点关注 etcd-cluster-operator。这是我和一些同事共同开发的 Operator，用于管理 Kubernetes 内部的 etcd。本文不是专门介绍 Operator 或 etcd 本身，所以我不会太过详细介绍 etcd 的细节，只要能够让你了解 etcd 的用途即可。

简单地说，etcd 是一个分布式键值数据存储。它有能力管理自己的稳定性，只要：

- 每个 etcd 实例都有一个用于计算、网络和存储的独立故障域。

- 每个 etcd 实例都有一个唯一的网络名称。

- 每个 etcd 实例都可以连接到其他实例。

- 每个 etcd 实例都知道其他实例的存在。

此外：

- etcd 集群的增长或缩小需要使用 etcd 管理 API 进行特定的操作，在添加或删除实例之前声明集群要发生的变化。

- 可以使用 etcd 管理 API 上的“快照”端点进行备份。通过 gRPC 调用它，你将得到一个备份文件。

- 使用 etcdctl 工具操作备份文件和 etcd 主机上的数据目录来实现恢复。这在真实的机器上很容易，但在 Kubernetes 上需要做一些协调。

正如你所看到的，这比 Kubernetes StatefulSet 能做更多的事情，所以我们使用 Operator。我们不会深入讨论 etcd-cluster-operator 的机制，但在本文的其余部分，我们都将引用这个 Operator 示例。

Operator 剖析

Operator 由两部分组成：

一个或多个 Kubernetes 自定义资源定义 (CRD)，它们描述了一种新的资源，包括应该具有哪些字段。CRD 可能会有多个，例如 etcd-cluster-operator 同时使用 EtcdCluster 和 EtcdPeer 来封装不同的概念。

一个运行中的软件，读取自定义资源并作出响应。

通常，Operator 被包含并部署在 Kubernetes 集群中，通常使用一个简单的 Deployment 资源。理论上，只要 Operator 能够与集群的 Kubernetes API 通信，它就可以在任何地方运行。但是，在集群中运行 Operator 通常更容易。通常情况下会使用自定义 Namespace 将 Operator 与其他资源分隔开来。

如果我们使用这种方法来运行 Operator，还需要做一些事情：

- 一个容器镜像，其中包含 Operator 可执行文件。

- 一个 Namespace。

- Operator 的 ServiceAccount，授予读取自定义资源的权限，并配置它要管理的资源 (例如 Pod)。

- 用于 Operator 容器的 Deployment。

- ClusterRoleBinding 和 ClusterRole 资源，绑定到 ServiceAccount。

- Webhook 配置。

稍后我们将详细讨论权限模型和 Webhook。

软件和工具

第一个问题是编程语言和生态系统。从理论上讲，几乎任何能够进行 HTTP 调用的语言都可以使用 Operator。假设 Operator 部署在与资源相同的集群中，那么只需要在集群容器中运行它即可。通常是 linux/x86_64，这也是 etcd-cluster-operator 的目标平台，但 Operator 也可以被编译成 arm64 或其他架构，甚至是 Windows 容器。

Go 语言拥有最成熟的工具。用于构建 Kubernetes 控制器的框架 controller-runtime 可以作为一个独立的工具。此外，Kubebuilder 和 Operator SDK 等项目都构建在控制器运行时之上，目的是提供一种流线化的开发体验。

除了 Go 语言，其他语言 (如 Java、Rust、Python 和其他语言) 通常会提供用于连接 Kubernetes API 或者专门用于构建 Operator 的工具。这些工具的成熟度和支持水平各有差别。

另一种选择是通过 HTTP 直接与 Kubernetes API 交互。这种方式所需的工作量最大，好处是团队可以使用他们最熟悉的编程语言。

最终，这种选择取决于负责构建和维护 Operator 的团队。如果团队已经习惯使用 Go，那么 Go 生态系统丰富的工具显然是最佳的选择。如果团队还没有使用 Go，那么就需要做出权衡，要么在学习和培训更成熟的生态系统工具方面付出代价，要么选择不成熟但团队熟悉其底层语言的生态系统。

对于 etcd-cluster-operator 来说，开发团队已经非常精通 Go，因此 Go 对我们来说是一个很明智的选择。我们还选择使用 Kubebuilder 而不是 Operator SDK，但这只是因为我们对它比较熟悉。我们的目标平台是 linux/x86_64，但如果需要的话，也可以以其他平台为目标。

自定义资源和目标状态

我们为我们的 etcd Operator 创建了一个叫作 EtcdCluster 的自定义资源定义。安装好 CRD 后，用户就可以创建 EtcdCluster 资源。EtcdCluster 资源描述了 etcd 集群的需求，并给出了它的配置。

PlainTextapiVersion:etcd.improbable.io/v1alpha1kind:EtcdClustermetadata: name: my-first-etcd-clusterspec: replicas: 3 version: 3.2.28

apiVersion 指定这是哪个版本的 API，在本例中是 v1alpha1。kind 声明这是一个 EtcdCluster。与其他类型的资源一样，我们有一个 metadata，它必须包含一个 name，也可能包含一个 namespace、labels、annotations 和其他标准项。这样我们就可以像对待 Kubernetes 中的其他资源一样对待 EtcdCluster。例如，我们可以使用一个标签来标识哪个团队负责哪一个集群，然后通过 kubectl get etcdcluster -l team=foo 搜索这些集群，就像使用其他标准资源一样。

spec 字段包含了有关这个 etcd 集群的运维信息。还有很多其他字段，但这里我们只介绍最基本的字段。version 字段描述要部署的 etcd 版本，replicas 字段描述有多少个实例。

还有一个 status 字段 (在示例中不可见)，运维人员用这个字段来描述集群的当前状态。spec 和 status 是 Kubernetes API 提供的标准字段，可以很好地与其他资源和工具集成。

因为我们使用了 Kubebuilder，所以可以借助工具生成这些自定义资源定义。我们写了一个 Go 结构体，定义了 spec 和 status 字段：

    type EtcdClusterSpec struct {
        Version     string               `json:"version"`
        Replicas    *int32               `json:"replicas"`
        Storage     *EtcdPeerStorage     `json:"storage,omitempty"`
        PodTemplate *EtcdPodTemplateSpec `json:"podTemplate,omitempty"`
    }

1
2
3
4
5
6
Go

基于这个 Go 结构体（和一个类似的 status 结构体），Kubebuilder 会生成我们的自定义资源定义，我们只需要编写代码处理调解逻辑即可。

其他语言提供的支持可能有所不同。如果你使用的是专为 Operator 设计的框架，那么可能会生成这个，例如 Rust 库 kube-derive 的生成方式就跟这个差不多。如果有团队直接使用 Kubernetes API，那么他们就必须分别编写 CRD 和用于解析数据的代码。

调解循环

现在我们已经有了描述 etcd 集群的方式，可以构建 Operator 来管理集群资源。Operator 可以以任何方式运行，而几乎所有 Operator 都可以使用控制器模式。

控制器是一种简单的程序循环，通常被称为“Reconcile(调解) 循环”，它可以执行以下逻辑：

1. 观察期望的状态。

2. 观察所管理资源的当前状态。

3. 采取行动，使托管的资源处在期望的状态。

对于 Kubernetes 中的 Operator，目标状态就是资源（示例中是 EtcdCluster 的 spec 字段指定的值）。我们的托管资源可以是集群内部或外部的任何资源。在我们的示例中，我们将创建其他 Kubneretes 资源，如 ReplicaSets、PersistentVolumeClaims 和 Services。

对于 etcd，我们直接连接到 etcd 进程，使用管理 API 来获取它的状态。这种“非 kubernetes”的访问方式需要小心一点，因为它可能会受到网络中断的影响，所以对于这种情况，并不一定是因为服务被关闭了。我们不能将无法连接到 etcd 作为 etcd 没有在运行的信号 (如果我们这么认为了，那么重启 etcd 实例只会加重网络中断的发生)。

通常，在与非 Kubernetes API 服务通信时，最重要的是要考虑可用性或一致性。对于 etcd 来说，如果我们获得响应，那它们一定是一致的，但其他系统可能不是这样。关键要避免由于信息过时而导致错误操作，从而使中断变得更糟。

控制器的特性

对于控制器来说，最简单的就是定时运行调解循环，比如每 30 秒一次。这样做是可以的，但有很多缺点。例如，它必须能够检测上一次循环是否还在运行，这样就不会同时运行两个循环。此外，这意味着每 30 秒会对 Kubernetes 进行一次完整的扫描来获得相关的资源，然后，对于 EtcdCluster 的每个实例，需要运行调解函数来获得相关 Pod 和其他资源。这种方式给 Kubernetes API 造成大量的负载。

这也导致出现了一种非常“程序性”的方法，因为在下一次协调之前可能需要很长时间才能尽可能快地执行每个循环。例如，一次性创建多个资源。这可能会导致一种非常复杂的状态，运维人员需要进行很多检查才能知道要做什么，而且很有可能会出错。

为了解决这个问题，控制器提供了一些特性：

- Kubernetes API 监听。

- API 缓存。

- 批量更新。

所有这些都可以有效地减少要执行的任务，因为运行单个循环的成本和需要等待的时间都减少了，协调逻辑的复杂性也就降低了。

API 监听

Kubernetes API 支持“监听”，而不是定时扫描。API 使用者可以对感兴趣的资源或资源类别进行注册，并在匹配的资源发生变更时收到通知。因为请求负载减少了，所以 Operator 大部分时间处于空闲状态，而且几乎可以立即对变更做出响应。Operator 框架通常会为你处理监听所需的注册和管理操作。

这种设计的另一个结果是你还需要监听你所创建的资源。例如，如果我们创建了 Pods，那么也必须监听我们创建的 Pod。如果它们被删除或修改，导致与我们想要的状态不一致，我们就可以收到通知，并纠正它们。

我们现在可以进一步简化调解程序。例如，为了响应 EtcdCluster，Operator 希望创建一个 Service 和一些 EtcdPeer 资源。它不是一次性创建好它们，而是先创建 Service，然后退出。但因为我们关注了自己的 Services，我们会收到通知，并立即重新进行调解。这样我们就可以创建对等资源了。否则，我们将创建大量的资源，然后为每个资源重新调解一次，这可能会触发更多的重新调解。

这种设计有助于保持调解器循环的简单，因为只需要执行一个操作就退出，开发人员不需要处理复杂的状态。

这样做的一个主要后果是可能会错过更新。网络中断、Pod 重启和其他问题在某些情况下可能导致错过事件。为了解决这个问题，关键在于 Operator 的运行方式应该“基于条件”而不是“基于边缘”。

这些术语来自信号控制软件，是指基于信号电压做出响应。在软件领域，当我们说“基于边缘”时，意思是“对事件做出反应”，当我们说“基于条件”时，意思是“对观察到的状态做出反应”。

例如，如果一个资源被删除，我们可以观察到删除事件并选择重新创建。但是，如果我们错过了删除事件，就可能永远不会尝试重新创建。或者，更糟糕的是，我们认为它还在，导致后续出现问题。相反，“基于条件”的方法将触发器简单地视为应该重新进行调解。它将再次观察外部状态，丢弃触发它的变更。

API 缓存

控制器的另一个主要特性是缓存请求。如果我们请求 Pods，并且会在 2 秒后再次触发，那么我们可能会为第二个请求保留缓存结果。这减少了 API 服务器的负载，但也给开发人员带来了一些需要注意的问题。

由于资源请求可能过期，我们必须处理这个问题。资源创建没有被缓存，因此可能出现这种情况：

- 调解 EtcdCluster 资源

- 搜索 Service，没有找到。

- 创建 Service 并退出。

- 对创建的 Service 做出响应。

- 搜索 Service，缓存过期，找不到。

- 创建 Service。

我们错误地创建了一个相同的 Service。Kubernetes API 将会处理这个问题，并给出一个错误，说明 Service 已经存在。因此，我们必须处理这个问题。一般来说，最好的做法是在以后的某个时间进行重新调解。在 Kubebuilder 中，只是简单地在 reconcile 函数中返回一个错误就会导致这种情况发生，但不同的框架可能会有所不同。当稍后重新运行时，缓存最终会保持一致，并可能发生下一阶段的调解。

这样做的一个副作用是所有资源都必须有确定的名称。否则，如果我们创建了一个重复的资源，可能会使用不同的名称，导致真正的资源重复。

批量更新

在某些情况下，我们可能会同时进行很多个调解。例如，如果我们正在监听大量的 Pod 资源，其中有很资源同时处于停止状态 (例如，由于节点故障、管理员操作错误，等等)，那么我们希望得到多次通知。然而，在第一次调解触发并观察到集群状态时，所有的 Pod 都已经消失了，那么后续的调解就是没有必要的。

如果数量很小，这就不是一个问题。但在较大的集群中，当一次处理数百或数千个更新时，这样做有可能会导致调解循环慢得像爬行一样，因为它一次性重复 100 次相同的操作，甚至会导致队列超载，并最终导致 Operator 崩溃。

因为我们的调解函数是“基于条件”的，所以我们可以对其加以优化来解决这个问题。当我们将特定资源的更新操作放入队列时，如果队列中已经有该资源的更新操作，那么就将其删除。在从队列读取数据之前先等待一下，我们就可以有效地进行“批量”操作。因此，如果 200 个 Pod 同时停止，我们可能只需要进行一次调解，具体取决于 Operator 及其队列的配置情况。

权 限

访问 Kubernetes API 必须提供凭证。在集群中，这是由 ServiceAccount 负责处理的。我们可以使用 ClusterRole 和 ClusterRoleBinding 资源将权限与 ServiceAccount 关联起来。对于 Operator 来说，这很关键。Operator 必须拥有权限来 get、list 和 watch 它在整个集群中管理的资源。此外，对于它创建的任何资源，都需要权限。例如，Pods、StatefulSets、Services 等。

Kubebuilder 和 Operator SDK 等框架可以为你提供这些权限。例如，Kubebuilder 采用了注解为每个控制器分配权限。如果多个控制器合并为一个二进制文件 (就像我们对 etcd-cluster-operator 所做的那样)，那么权限也将合并在一起。

    //+kubebuilder:rbac:groups=etcd.improbable.io,resources=etcdpeers,verbs=get;list;watch
    //+kubebuilder:rbac:groups=etcd.improbable.io,resources=etcdpeers/status,verbs=get;update;patch
    //+kubebuilder:rbac:groups=apps,resources=replicasets,verbs=list;get;create;watch
    //+kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=list;get;create;watch;delete

1
2
3
4
Plain Text

这是 EtcdPeer 资源的调解器权限。可以看到，我们 get、list 和 watch 自己的资源，并且可以 update 和 patch 状态子资源。我们可以只更新状态，将信息显示给其他用户。最后，我们对所管理的资源具有广泛的权限，可以根据需要创建和删除它们。

验证和默认值

虽然自定义资源本身提供了一定级别的验证和默认值，但更复杂的检查操作需要由 Operator 来执行。最简单的方法是在 Operator 读取资源时执行这些操作，无论是 watch 返回的，还是手动读取后。但是，这意味着默认值将永远不会被应用到 Kubernetes 中，这种行为会让管理员感到困惑。

更好的方法是使用验证和可变的 Webhook 配置。这些资源告诉 Kubernetes，当一个资源被创建、更新或者在持久化之前被删除时，必须使用 Webhook。

例如，可变 Webhook 可以用来设置默认值。在 Kubebuilder 中，我们提供了一些额外的配置来创建 MutatingWebhookConfiguration，Kubebuilder 负责提供 API 端点。我们只需要在 spec 结构体中设置 Default 值。然后，当资源被创建时，Webhook 在持久化资源之前被调用，就会应用默认值。

不过，我们仍然要在读取资源时应用默认值。Operator 不能假设已经知道平台是否启用了 Webhook。即使启用了，也可能配置错误，或者因为网络中断导致 Webhook 被跳过，或者资源可能在配置 Webhook 之前就已经被应用过了。所有这些问题都意味着，虽然 Webhook 提供了更好的用户体验，但 Operator 代码不能完全依赖它们，必须再次应用默认值。

测 试

任何一个单独的逻辑单元都可以使用编程语言的常规工具进行单元测试，但是，在进行集成测试时会出现一些特定的问题。我们可能会把 API 服务器当成可以被 mock 的数据库。但在真实的系统中，API 服务器会执行大量的验证和默认操作。这意味着测试和现实之间的行为可能是不一样的。

一般来说，主要有两种方式：

第一种方法，下载测试工具并执行 kube-apiserver etcd 可执行文件，创建一个真正的 API 服务器。当然，虽然你可以创建一个 ReplicaSet，但缺少了可以创建 Pods 的 Kubernetes 组件，所以我们看不到有东西真正在运行。

第二种方法更加全面一些，它使用一个真正的 Kubernetes 集群，可以运行 Pods，并能准确做出响应。通过使用 kind，这种集成测试变得更加容易。kind 是“Kubernetes in Docker”的缩写，它可以在任何可以运行 Docker 容器的地方运行一个完整的 Kubernetes 集群。它提供了一个 API 服务器，可以运行 Pods，并运行 Kubernetes 所有主要的组件。因此，使用了 kind 的测试可以在笔记本电脑上或 CI 中运行，并提供近乎完美的 Kubernetes 体验。

总 结

在这篇文章中，我们谈到了很多想法：

- 将 Operator 作为 Pods 部署在集群中。

- 可以支持任何一种编程语言，所以请选择最适合团队的那一种。不过，Go 语言拥有最成熟的生态系统。

- 小心使用非 kubernetes 资源，特别是在网络中断或上游 API 发生故障时，它们可能会导致更严重的中断。

- 在每个调解周期中执行一个操作，然后退出，并允许 Operator 重新将其放入队列。

- 使用“基于条件”的方法，忽略触发调解的事件的内容。

- 为新资源使用确定性的命名。

- 为你的服务帐户提供最小权限。

- 在 Webhook 和代码中应用默认值。

- 使用 kind 进行集成测试。

有了这些工具，你就可以构建 Operator 来简化部署，并减轻运维团队的负担，无论是你所拥有的应用程序，还是你自己开发的应用程序。

关于作者

James Laverack 是英国 Kubernetes 专业服务公司 Jetstack 的一名解决方案工程师。凭借超过 7 年的行业经验，他的大部分时间都在帮助企业实现云原生化。他也是一个 Kubernetes 贡献者，并且从 1.18 版本开始就加入了 Kubernetes 发布团队。

原文链接

Kubernetes Operators in Depth

<https://www.infoq.com/articles/kubernetes-operators-in-depth/>
