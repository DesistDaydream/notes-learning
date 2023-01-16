---
title: docker-shim 何去何从
---

#

关于 dockershim 即将灭亡的传言无疑存在严重夸大。如果一直有关注 Kubernetes 生态系统，很多朋友一时之间可能确实被 Kubernetes 1.20 版本的发布公告弄得有点不知所措。从公告内容来看，自 1.20 版本开始 dockershim 将被全面弃用。但请不要恐慌，调整呼吸，一切都会好起来。

更重要的是，Mirantis 现已同意与 Docker 开展合作，在 Kubernetes 之外独立维护 shim 代码并将其作为 Docker Engine API 的统一 CRI 接口。对于 Mirantis 客户而言，这意味着 Docker Engine 的商业支持版本 Mirantis Container Runtime（MCR）也将提供 CRI 兼容能力。我们将从<https://github.com/dims/cri-d...>，并逐步将其转化为开源项目<https://github.com/Mirantis/c...>。换句话说，你可以像之前一样继续基于 Docker Engine 构建 Kubernetes，唯一的区别就是 dockershim 由内置方案变成了外部方案。我们将共同努力，保证它在保持原有功能的同时，顺利通过各类一致性测试并提供与此前内置版本相同的使用体验。Mirantis 将在 Mirantis Kubernetes Engine 中使用 dockershim，Docker 方面也将在 Docker Desktop 中继续提供 dockershim。

从头说起……

用过 Kubernetes 的朋友都清楚，它的最大作用就是编排各类容器。对不少用户来说，容器已经与 Docker 完全统一了起来。但这种说法并不准确，Docker 本身只是彻底改变了容器技术并将其推向了通用舞台，因此 Docker Engine 也成为 Kubernetes 所支持的第一种（也是最初唯一一种）容器运行时。

但 Kubernetes 社区并不打算长期保持这样的状态。

从长远来看，社区希望能够使用多种不同类型的容器，因此参与者们创建了容器运行时接口（CRI），也就是容器引擎与 Kubernetes 间进行通信的标准方式。如果容器引擎与 CRI 相兼容，即可轻松在 Kubernetes 当中运行。

第一款兼容 CRI 的容器引擎是 containerd，而它来自……好吧，还是来自 Docker。很明显，Docker 本身不仅仅是一种容器运行时，而且提供可供其他用户消费的种种部件，甚至包括用户界面。因此，Docker 提取出与容器实际相关的部分，并将其调整为第一种与 CRI 兼容的运行时，而后把它捐赠给了云原生计算基金会（CNCF）。由此衍生出的 cri-containerd 组件具有运行时中立特性，而且能够支持多种 Linux 与 Windows 操作系统。

但这还留下最后一个问题——Docker 本身仍然不兼容 CRI。

Dockershim 是什么？

正如 Kubernetes 最初对 Docker Engine 提供内置支持一样，其中同样包含对各类存储卷解决方案、网络解决方案甚至是云服务商的内置支持。但要不断维护这些支持功能实在是太过麻烦，因此社区决定将所有第三方解决方案从核心中剥离出来并创建相关接口，例如：

- 容器运行时接口（CRI）

- 容器网络接口（CNI）

- 容器存储接口（CSI）

其中的基本思路在于，只要兼容这些接口，那么任何供应商都可以创建出能自动与 Kubernetes 相对接的产品。

当然，这绝不是说不兼容的组件就没办法与 Kubernetes 配合使用；Kubernetes 可以使用正确的组件完成各类协同。换言之，不兼容的组件仅仅需要加上个“shim（意为垫片）”，由其在组件与相应的 Kubernetes 接口之间完成转换，即可轻松解决问题。例如，dockershim 会接收 CRI 命令并将其转换为 Docker Engine 能够理解的内容，反之亦然。但在第三方组件被从 Kubernetes 核心内剥离的背景之下，dockershim 自身也需要逐步退出。

虽然听起来好像事情不小，但实际上没那么严重。大家使用 docker build 构建起的 CRI 兼容型镜像，未来仍然可以与 Kubernetes 正常配套使用。

Kubernetes 放弃对 dockershim 的支持，会带来哪些影响？

对大多数人来说，弃用 dockershim 其实半点影响也没有。这是因为大部分用户既意识不到 dockershim 的存在，实际上使用的也不是 Docker 本体；相反，他们使用的是与 CRI 相兼容的 containerd。

当然，也有一部分用户（包括 Mirantis 的客户）在运行依赖于 dockershim 的工作负载，借此与 Kubernetes 实现无缝协作。

考虑到 dockershim 仍然是不少企业难以割舍的重要组件，Mirantis 与 Docker 达成协议，继续支持并开发 dockershim。只不过这一次，dockershim 将以独立开源组件的身份存在。

那么，这到底意味着什么？

简单来讲，如果你直接使用 containerd，则不必抱有任何担心；因为 containerd 能够与 CRI 相兼容。如果你身为 Mirantis 的客户，同样不用担心；因为 Mirantis 容器运行时将包含对 dockershim 的支持，确保其与 CRI 相兼容。

但如果你使用的是开源 Docker Engine，则 dockershim 项目将以开源组件的形式提供，您可以继续在 Kubernetes 上正常使用；唯一的区别就是需要对配置做出少量修改，具体请参见我们后续发布的说明文档。

所以，请大家不必惊异。Docker 还在，dockershim 还在，一切如常。

原文链接：<https://www.mirantis.com/blog...>
