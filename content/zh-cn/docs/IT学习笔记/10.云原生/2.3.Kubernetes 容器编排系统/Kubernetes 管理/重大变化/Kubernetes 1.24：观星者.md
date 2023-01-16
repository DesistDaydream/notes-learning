---
title: Kubernetes 1.24：观星者
---

原文链接：[公众号-CNCF，Kubernetes 1.24：观星者](https://mp.weixin.qq.com/s/Nd9UFyqRKS6qpUJGF9nAYg)

作者：Kubernetes 1.24 发布团队\[1]

我们很兴奋地宣布 Kubernetes 1.24 的发布，这是 2022 年的第一个版本！

这个版本包含 46 个增强功能：14 个增强功能已经升级到稳定版，15 个增强功能正在进入 beta 版，13 个增强功能正在进入 alpha 版。此外，两个功能已给弃用，两个功能给删除。

## 主要主题

### Dockershim 从 kubelet 中删除

在 1.20 版中给弃用之后，在 Kubernetes v1.24 版中，dockershim 组件已从 kubelet 中删除。从 1.24 版开始，你将需要使用其他受支持的运行时\[2]（如 containerd 或 CRI-O），如果你依赖 Docker Engine 作为容器运行时，可使用 cri-dockerd。有关确保你的群集已准备好进行此删除的更多信息，请参见本指南\[3]。

### 默认情况下关闭（新的）测试 API

默认情况下，新的测试版 API 不会在集群中启用\[4]。默认情况下，将继续启用现有的 beta API 和现有 beta API 的新版本。

### 签名发布工件

发布工件使用 cosign\[5]进行签名\[6]，并且有验证镜像签名\[7]的实验支持。发布工件的签名和验证是增加 Kubernetes 发布过程的软件供应链安全性\[8]的一部分。

### OpenAPI v3

Kubernetes 1.24 为以 OpenAPI v3 格式\[9]发布 API 提供了测试版支持。

### Storage Capacity 和 Volume Expansion 是 GA

存储容量跟踪\[10]支持通过 CSIStorageCapacity 对象\[11]显示当前可用的存储容量，并通过后期绑定（late binding）增强使用 CSI 卷的 pod 的编排。

卷扩展\[12]增加了对调整现有永久卷大小的支持。

### NonPreemptingPriority 到 Stable

此功能为 PriorityClasses 添加了新选项\[13]，可以启用或禁用 pod 抢占。

### 存储插件迁移

正在进行的工作是迁移树内存储插件的内部\[14]，以调用 CSI 插件，同时保持原来的 API。Azure Disk\[15]和 OpenStack Cinder\[16]插件都已迁移。

### gRPC 探针升级到 Beta

在 Kubernetes 1.24 中，gRPC probes 功能\[17]已经进入测试阶段，默认情况下是可用的。现在，你可以在 Kubernetes 中为 gRPC 应用程序原生配置启动、活动和就绪探测器\[18]，而无需暴露 HTTP 端点或使用额外的可执行文件。

### Kubelet 凭据提供程序升级到 Beta

最初在 Kubernetes 1.20 中作为 Alpha 版本发布，kubelet 对镜像凭证提供者\[19]的支持现在已经升级到 Beta 版本。这允许 kubelet 使用 exec 插件动态检索容器镜像注册中心的凭证，而不是将凭证存储在节点的文件系统上。

### Contextual Logging 处于 Alpha

Kubernetes 1.24 引入了上下文日志记录\[20]，使函数的调用者能够控制日志记录的所有方面（输出格式、详细程度、附加值和名称）。

### 避免服务的 IP 分配中的冲突

Kubernetes 1.24 引入了一个新的选择加入特性，允许你为服务的静态 IP 地址分配软保留一个范围\[21]。手动启用此功能后，群集将倾向于从服务 IP 地址池中自动分配，从而降低冲突风险。

可以分配 ClusterIP 服务：

- 动态地，这意味着群集将在配置的服务 IP 范围内自动选择一个空闲 IP。
- 静态，这意味着用户将在配置的服务 IP 范围内设置一个 IP。

ClusterIP 服务是唯一的；因此，尝试使用已经分配的 ClusterIP 创建服务将会返回错误。

### Dynamic Kubelet Configuration 已从 Kubelet 中移除

在 Kubernetes 1.22 中被弃用后，Dynamic Kubelet Configuration 已从 kubelet 中移除。在 Kubernetes 1.26 中，该特性将从 API 服务器中删除。

## CNI 版本相关的重大变更

在升级到 Kubernetes 1.24 之前，请验证你使用/升级到的容器运行时已经过测试，可以在此版本中正常工作。

例如，以下容器运行时正在为 Kubernetes 准备，或者已经为 Kubernetes 准备好了：

- containerd 1.6.4 版和更高版本、1.5.11 版和更高版本
- CRI-O 1.24 和更高版本

当 CNI 插件未升级和/或 CNI 配置版本未在 CNI 配置文件中声明时，containerd v1.6.0–1.6.3 中的 pod CNI 网络设置和拆除存在服务问题。containerd 团队报告说，“containerd v1.6.4 解决了这些问题。”

使用 containerd v1.6.0–1.6.3，如果你不升级 CNI 插件和/或声明 CNI 配置版本，你可能会遇到“Incompatible CNI versions”或“Failed to destroy network for sandbox”错误情况。

## 其他更新

### 升级到稳定

在此版本中，有 14 项增强功能升级为稳定版：

- Container Storage Interface (CSI) Volume Expansion\[22]
- Pod Overhead\[23]: Account for resources tied to the pod sandbox but not specific containers.
- Add non-preempting option to PriorityClasses\[24]
- Storage Capacity Tracking\[25]
- OpenStack Cinder In-Tree to CSI Driver Migration\[26]
- Azure Disk In-Tree to CSI Driver Migration\[27]
- Efficient Watch Resumption\[28]: Watch can be efficiently resumed after kube-apiserver reboot.
- Service Type=LoadBalancer Class Field\[29]: Introduce a new Service annotation service.kubernetes.io/load-balancer-class that allows multiple implementations of type: LoadBalancer Services in the same cluster.
- Indexed Job\[30]: Add a completion index to Pods of Jobs with a fixed completion count.
- Add Suspend Field to Jobs API\[31]: Add a suspend field to the Jobs API to allow orchestrators to create jobs with more control over when pods are created.
- Pod Affinity NamespaceSelector\[32]: Add a namespaceSelector field for to pod affinity/anti-affinity spec.
- Leader Migration for Controller Managers\[33]: kube-controller-manager and cloud-controller-manager can apply new controller-to-controller-manager assignment in HA control plane without downtime.
- CSR Duration\[34]: Extend the CertificateSigningRequest API with a mechanism to allow clients to request a specific duration for the issued certificate.

### 主要变化

这个版本有两个主要变化：

- Dockershim Removal\[35]
- 默认情况下，测试 API 是关闭的\[36]

### 发布说明

在我们的发行说明\[37]中查看 Kubernetes 1.24 版本的全部细节。

### 下载

Kubernetes 1.24 可以在 GitHub\[38]下载。要开始使用 Kubernetes，请查看这些交互式教程\[39]或通过 kind\[40]，使用容器作为“节点”运行本地 Kubernetes 集群。你也可以使用 kubeadm\[41]轻松安装 1.24。

### 发布团队

如果没有 Kubernetes 1.24 发布团队成员的共同努力，这次发布是不可能的。这个团队一起交付每个 Kubernetes 版本的所有组件，包括代码、文档、发行说明等等。

特别感谢我们的发布负责人 James Laverack，他指导我们完成了一个成功的发布周期，并感谢所有发布团队成员投入时间和精力为 Kubernetes 社区发布 1.24 版本。

### 发布主题和徽标

Kubernetes 1.24：观星者

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/da4dd63e-1e80-46ed-9523-b6f8566ff267/640)

Kubernetes 1.24 的主题是观星者。

从古代天文学家到建造詹姆斯·韦伯太空望远镜的科学家，世世代代的人都怀着敬畏和好奇的心情仰望星空。星星鼓舞了我们，点燃了我们的想象力，指引我们在艰难的海上度过漫漫长夜。

随着这一版本的发布，我们向上看，当我们的社区走到一起时，什么是可能的。Kubernetes 是全球数百名贡献者和数千名支持数百万应用程序的最终用户的作品。每个人都是我们天空中的一颗星，帮助我们指引方向。

版本标志由 Britnee Laverack\[42]制作，描绘了一个设置在星空和昴宿星\[43]上的望远镜，昴宿星在神话中通常被称为“七姐妹”。数字 7 对 Kubernetes 项目来说特别吉祥，是对我们最初的“Project Seven”名称的引用。

这次发布的 Kubernetes 是以那些仰望夜空并感到惊奇的人命名的——致所有的观星者。✨

### 用户亮点

- 了解领先的零售电子商务公司 La Redoute\[44]如何使用 Kubernetes 以及其他 CNCF 项目来转变和简化其软件交付生命周期——从开发到运营。
- 为了确保对 API 调用的更改不会导致任何中断，Salt Security\[45]完全在 Kubernetes 上构建了其微服务，它通过 gRPC 进行通信，而 Linkerd 则确保消息是加密的。
- 在从私有云迁移到公共云的过程中，Allainz Direct\[46]工程师仅用了三个月时间就重新设计了 CI/CD 流水线，同时成功将 200 个工作流精简到 10-15 个。
- 了解总部位于英国的金融科技公司 Bink\[47]如何使用 Linkerd 更新其内部 Kubernetes 发行版，以构建一个与云无关的平台，该平台可根据需要进行扩展，同时允许他们密切关注性能和稳定性。
- 使用 Kubernetes，荷兰 Stichting Open Nederland\[48]组织在短短一个半月内创建了一个测试门户，以帮助在荷兰安全地重新开放活动。Testen voor Toegang\[49]测试平台利用 Kubernetes 的性能和可扩展性，帮助个人每天预约超过 400，000 次新冠肺炎测试\[50]。
- Santagostino\[51]与 SparkFabrik 合作，利用 Backstage，创建了开发人员平台 Samaritan 来集中服务和文档，管理服务的整个生命周期，并简化 Santagostino 开发人员的工作。

### 生态系统更新

- 2022 年欧洲 KubeCon + CloudNativeCon 将于 2022 年 5 月 16 日至 20 日在西班牙巴伦西亚举行！你可以在活动网站\[52]上找到有关会议和注册的更多信息。
- 在 2021 年云原生调查\[53]中，CNCF 见证了创纪录的 Kubernetes 和容器采用。来看一下调查结果\[54]。
- Linux 基金会和 CNCF 宣布推出新的云原生开发者训练营\[55]，为参与者提供设计、构建和部署云原生应用的知识和技能。查看公告\[56]以了解更多信息。

### 项目速度

CNCF K8s DevStats\[57]项目汇总了许多与 Kubernetes 和各种子项目的速度相关的有趣数据点。这包括从个人贡献到做出贡献的公司数量的所有内容，体现了发展这一生态系统的深度和广度。

在为期 17 周\[58]（1 月 10 日至 5 月 3 日）的 v1.24 发布周期中，我们看到了来自 1029 家公司\[59]和 1179 名个人\[60]的贡献。

## 网络研讨会预告

在 2022 年 5 月 24 日星期二上午 9:45–11 点（太平洋时间），加入 Kubernetes 1.24 发布团队的成员，了解此版本的主要功能，以及弃用和删除内容，以帮助规划升级。了解更多信息和注册，请访问 CNCF 在线计划网站的活动页面\[61]。

## 参与

加入 Kubernetes 最简单的方法，是加入与你的兴趣一致的特别殊兴趣小组（SIG）\[62]。你有什么想对 Kubernetes 社区说的吗？在我们每周的社区会议\[63]上，以及通过以下渠道分享你的观点：

- 在 Kubernetes 贡献者\[64]网站上找到更多关于贡献给 Kubernetes 的信息
- 在 Twitter @Kubernetesio\[65]上关注我们的最新动态
- 加入 Discuss\[66]的社区讨论
- 加入 Slack\[67]社区
- 在 Server Fault\[68]发布问题（或回答问题）。
- 分享你的 Kubernetes 故事\[69]
- 在博客\[70]上阅读更多关于 Kubernetes 的信息
- 了解更多关于 Kubernetes 发布团队\[71]的信息

### 参考资料

\[1]

Kubernetes 1.24 发布团队: [_https://github.com/kubernetes/sig-release/blob/master/releases/release-1.24/release-team.md_](https://github.com/kubernetes/sig-release/blob/master/releases/release-1.24/release-team.md)

\[2]

受支持的运行时: [_https://kubernetes.io/docs/setup/production-environment/container-runtimes/_](https://kubernetes.io/docs/setup/production-environment/container-runtimes/)

\[3]

指南: [_https://kubernetes.io/blog/2022/03/31/ready-for-dockershim-removal/_](https://kubernetes.io/blog/2022/03/31/ready-for-dockershim-removal/)

\[4]

默认情况下，新的测试版 API 不会在集群中启用: [_https://github.com/kubernetes/enhancements/issues/3136_](https://github.com/kubernetes/enhancements/issues/3136)

\[5]

cosign: [_https://github.com/sigstore/cosign_](https://github.com/sigstore/cosign)

\[6]

签名: [_https://github.com/kubernetes/enhancements/issues/3031_](https://github.com/kubernetes/enhancements/issues/3031)

\[7]

验证镜像签名: [_https://kubernetes.io/docs/tasks/administer-cluster/verify-signed-images/_](https://kubernetes.io/docs/tasks/administer-cluster/verify-signed-images/)

\[8]

增加 Kubernetes 发布过程的软件供应链安全性: [_https://github.com/kubernetes/enhancements/issues/3027_](https://github.com/kubernetes/enhancements/issues/3027)

\[9]

OpenAPI v3 格式: [_https://github.com/kubernetes/enhancements/issues/2896_](https://github.com/kubernetes/enhancements/issues/2896)

\[10]

存储容量跟踪: [_https://github.com/kubernetes/enhancements/issues/1472_](https://github.com/kubernetes/enhancements/issues/1472)

\[11]

CSIStorageCapacity 对象: [_https://kubernetes.io/docs/concepts/storage/storage-capacity/#api_](https://kubernetes.io/docs/concepts/storage/storage-capacity/#api)

\[12]

卷扩展: [_https://github.com/kubernetes/enhancements/issues/284_](https://github.com/kubernetes/enhancements/issues/284)

\[13]

PriorityClasses 添加了新选项: [_https://github.com/kubernetes/enhancements/issues/902_](https://github.com/kubernetes/enhancements/issues/902)

\[14]

迁移树内存储插件的内部: [_https://github.com/kubernetes/enhancements/issues/625_](https://github.com/kubernetes/enhancements/issues/625)

\[15]

Azure Disk: [_https://github.com/kubernetes/enhancements/issues/1490_](https://github.com/kubernetes/enhancements/issues/1490)

\[16]

OpenStack Cinder: [_https://github.com/kubernetes/enhancements/issues/1489_](https://github.com/kubernetes/enhancements/issues/1489)

\[17]

gRPC probes 功能: [_https://github.com/kubernetes/enhancements/issues/2727_](https://github.com/kubernetes/enhancements/issues/2727)

\[18]

配置启动、活动和就绪探测器: [_https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#configure-probes_](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#configure-probes)

\[19]

镜像凭证提供者: [_https://kubernetes.io/docs/tasks/kubelet-credential-provider/kubelet-credential-provider/_](https://kubernetes.io/docs/tasks/kubelet-credential-provider/kubelet-credential-provider/)

\[20]

上下文日志记录: [_https://github.com/kubernetes/enhancements/issues/3077_](https://github.com/kubernetes/enhancements/issues/3077)

\[21]

静态 IP 地址分配软保留一个范围: [_https://kubernetes.io/docs/concepts/services-networking/service/#service-ip-static-sub-range_](https://kubernetes.io/docs/concepts/services-networking/service/#service-ip-static-sub-range)

\[22]

Container Storage Interface (CSI) Volume Expansion: [_https://github.com/kubernetes/enhancements/issues/284_](https://github.com/kubernetes/enhancements/issues/284)

\[23]

Pod Overhead: [_https://github.com/kubernetes/enhancements/issues/688_](https://github.com/kubernetes/enhancements/issues/688)

\[24]

Add non-preempting option to PriorityClasses: [_https://github.com/kubernetes/enhancements/issues/902_](https://github.com/kubernetes/enhancements/issues/902)

\[25]

Storage Capacity Tracking: [_https://github.com/kubernetes/enhancements/issues/1472_](https://github.com/kubernetes/enhancements/issues/1472)

\[26]

OpenStack Cinder In-Tree to CSI Driver Migration: [_https://github.com/kubernetes/enhancements/issues/1489_](https://github.com/kubernetes/enhancements/issues/1489)

\[27]

Azure Disk In-Tree to CSI Driver Migration: [_https://github.com/kubernetes/enhancements/issues/1490_](https://github.com/kubernetes/enhancements/issues/1490)

\[28]

Efficient Watch Resumption: [_https://github.com/kubernetes/enhancements/issues/1904_](https://github.com/kubernetes/enhancements/issues/1904)

\[29]

Service Type=LoadBalancer Class Field: [_https://github.com/kubernetes/enhancements/issues/1959_](https://github.com/kubernetes/enhancements/issues/1959)

\[30]

Indexed Job: [_https://github.com/kubernetes/enhancements/issues/2214_](https://github.com/kubernetes/enhancements/issues/2214)

\[31]

Add Suspend Field to Jobs API: [_https://github.com/kubernetes/enhancements/issues/2232_](https://github.com/kubernetes/enhancements/issues/2232)

\[32]

Pod Affinity NamespaceSelector: [_https://github.com/kubernetes/enhancements/issues/2249_](https://github.com/kubernetes/enhancements/issues/2249)

\[33]

Leader Migration for Controller Managers: [_https://github.com/kubernetes/enhancements/issues/2436_](https://github.com/kubernetes/enhancements/issues/2436)

\[34]

CSR Duration: [_https://github.com/kubernetes/enhancements/issues/2784_](https://github.com/kubernetes/enhancements/issues/2784)

\[35]

Dockershim Removal: [_https://github.com/kubernetes/enhancements/issues/2221_](https://github.com/kubernetes/enhancements/issues/2221)

\[36]

默认情况下，测试 API 是关闭的: [_https://github.com/kubernetes/enhancements/issues/3136_](https://github.com/kubernetes/enhancements/issues/3136)

\[37]

发行说明: [_https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.24.md_](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.24.md)

\[38]

GitHub: [_https://github.com/kubernetes/kubernetes/releases/tag/v1.24.0_](https://github.com/kubernetes/kubernetes/releases/tag/v1.24.0)

\[39]

交互式教程: [_https://kubernetes.io/docs/tutorials/_](https://kubernetes.io/docs/tutorials/)

\[40]

kind: [_https://kind.sigs.k8s.io/_](https://kind.sigs.k8s.io/)

\[41]

kubeadm: [_https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/_](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/)

\[42]

Britnee Laverack: [_https://www.instagram.com/artsyfie/_](https://www.instagram.com/artsyfie/)

\[43]

昴宿星: [_https://en.wikipedia.org/wiki/Pleiades_](https://en.wikipedia.org/wiki/Pleiades)

\[44]

La Redoute: [_https://www.cncf.io/case-studies/la-redoute/_](https://www.cncf.io/case-studies/la-redoute/)

\[45]

Salt Security: [_https://www.cncf.io/case-studies/salt-security/_](https://www.cncf.io/case-studies/salt-security/)

\[46]

Allainz Direct: [_https://www.cncf.io/case-studies/allianz/_](https://www.cncf.io/case-studies/allianz/)

\[47]

Bink: [_https://www.cncf.io/case-studies/bink/_](https://www.cncf.io/case-studies/bink/)

\[48]

Stichting Open Nederland: [_http://www.stichtingopennederland.nl/_](http://www.stichtingopennederland.nl/)

\[49]

Testen voor Toegang: [_https://www.testenvoortoegang.org/_](https://www.testenvoortoegang.org/)

\[50]

每天预约超过 400，000 次新冠肺炎测试: [_https://www.cncf.io/case-studies/true/_](https://www.cncf.io/case-studies/true/)

\[51]

Santagostino: [_https://www.cncf.io/case-studies/santagostino/_](https://www.cncf.io/case-studies/santagostino/)

\[52]

活动网站: [_https://events.linuxfoundation.org/archive/2021/kubecon-cloudnativecon-europe/_](https://events.linuxfoundation.org/archive/2021/kubecon-cloudnativecon-europe/)

\[53]

2021 年云原生调查: [_https://www.cncf.io/announcements/2022/02/10/cncf-sees-record-kubernetes-and-container-adoption-in-2021-cloud-native-survey/_](https://www.cncf.io/announcements/2022/02/10/cncf-sees-record-kubernetes-and-container-adoption-in-2021-cloud-native-survey/)

\[54]

调查结果: [_https://www.cncf.io/reports/cncf-annual-survey-2021/_](https://www.cncf.io/reports/cncf-annual-survey-2021/)

\[55]

云原生开发者训练营: [_https://training.linuxfoundation.org/training/cloudnativedev-bootcamp/_](https://training.linuxfoundation.org/training/cloudnativedev-bootcamp/)

\[56]

公告: [_https://www.cncf.io/announcements/2022/03/15/new-cloud-native-developer-bootcamp-provides-a-clear-path-to-cloud-native-careers/_](https://www.cncf.io/announcements/2022/03/15/new-cloud-native-developer-bootcamp-provides-a-clear-path-to-cloud-native-careers/)

\[57]

CNCF K8s DevStats: [_https://k8s.devstats.cncf.io/d/12/dashboards?orgId=1\&refresh=15m_](https://k8s.devstats.cncf.io/d/12/dashboards?orgId=1&refresh=15m)

\[58]

为期 17 周: [_https://github.com/kubernetes/sig-release/tree/master/releases/release-1.24_](https://github.com/kubernetes/sig-release/tree/master/releases/release-1.24)

\[59]

1029 家公司: [_https://k8s.devstats.cncf.io/d/9/companies-table?orgId=1\&var-period_name=v1.23.0 - now\&var-metric=contributions_](https://k8s.devstats.cncf.io/d/9/companies-table?orgId=1&var-period_name=v1.23.0%20-%20now&var-metric=contributions)

\[60]

1179 名个人: [_https://k8s.devstats.cncf.io/d/66/developer-activity-counts-by-companies?orgId=1\&var-period_name=v1.23.0 - now\&var-metric=contributions\&var-repogroup_name=Kubernetes\&var-country_name=All\&var-companies=All\&var-repo_name=kubernetes%2Fkubernetes_](https://k8s.devstats.cncf.io/d/66/developer-activity-counts-by-companies?orgId=1&var-period_name=v1.23.0%20-%20now&var-metric=contributions&var-repogroup_name=Kubernetes&var-country_name=All&var-companies=All&var-repo_name=kubernetes%2Fkubernetes)

\[61]

活动页面: [_https://community.cncf.io/e/mck3kd/_](https://community.cncf.io/e/mck3kd/)

\[62]

特别殊兴趣小组（SIG）: [_https://github.com/kubernetes/community/blob/master/sig-list.md_](https://github.com/kubernetes/community/blob/master/sig-list.md)

\[63]

社区会议: [_https://github.com/kubernetes/community/tree/master/communication_](https://github.com/kubernetes/community/tree/master/communication)

\[64]

Kubernetes 贡献者: [_https://www.kubernetes.dev/_](https://www.kubernetes.dev/)

\[65]

@Kubernetesio: [_https://twitter.com/kubernetesio_](https://twitter.com/kubernetesio)

\[66]

Discuss: [_https://discuss.kubernetes.io/_](https://discuss.kubernetes.io/)

\[67]

Slack: [_http://slack.k8s.io/_](http://slack.k8s.io/)

\[68]

Server Fault: [_https://serverfault.com/questions/tagged/kubernetes_](https://serverfault.com/questions/tagged/kubernetes)

\[69]

故事: [_https://docs.google.com/a/linuxfoundation.org/forms/d/e/1FAIpQLScuI7Ye3VQHQTwBASrgkjQDSS5TP0g3AXfFhwSM9YpHgxRKFA/viewform_](https://docs.google.com/a/linuxfoundation.org/forms/d/e/1FAIpQLScuI7Ye3VQHQTwBASrgkjQDSS5TP0g3AXfFhwSM9YpHgxRKFA/viewform)

\[70]

博客: [_https://kubernetes.io/blog/_](https://kubernetes.io/blog/)

\[71]

Kubernetes 发布团队: [_https://github.com/kubernetes/sig-release/tree/master/release-team_](https://github.com/kubernetes/sig-release/tree/master/release-team)

**\_\***CNCF\***\*\*\*\***（\_\_**_云原生计算基金会_**）致力于培育和维护一个厂商中立的开源生态系统，来推广云原生技术。我们通过将最前沿的模式民主化，让这些创新为大众所用。请长按以下二维码进行关注。\_\*\*
