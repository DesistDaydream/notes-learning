---
title: Kubernetes v1.26：振奋人心
---

作者：Kubernetes 1.26 发布团队\[1]

我们怀着无比的喜悦宣布 Kubernetes v1.26 版的发布！

这个版本总共包括 37 个增强：其中 11 个升级到 Stable，10 个升级到 Beta 版，16 个进入 Alpha 版。我们还有 12 个功能被弃用或删除，其中三个我们在本次宣布中会详细介绍。

## 发布主题和徽标

### Kubernetes 1.26：振奋人心（Electrifying）

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/a8966cdb-f942-4f4b-846f-0223ed6e9a9b/640)

Kubernetes v1.26 版的主题是振奋人心（Electrifying）。

Kubernetes 的每一个版本，都是敬业的志愿者们共同努力的结果，只有通过使用分散在全球多个数据中心和地区的各种复杂的计算资源，才有可能实现。发布的最终结果——二进制文件、镜像容器、文档——随后被部署在越来越多的个人、内部和云计算资源上。

在此版本中，我们希望认识到开发和使用 Kubernetes 所基于的所有这些构建模块的重要性，同时提高对考虑能耗足迹的重要性的认识：环境可持续性是任何软件解决方案的创作者和用户不可避免的关注，以及像 Kubernetes 这样的软件的环境足迹，我们相信这一领域将在未来的版本中发挥重要作用。

作为一个社区，我们总是努力使每个新的发布过程比以前更好（例如，在这个版本中，我们已经开始使用 Projects 来跟踪增强\[2]）。如果[1.24 版“观星者”](https://mp.weixin.qq.com/s?__biz=MzI5ODk5ODI4Nw==&mid=2247520103&idx=1&sn=ff9a8fbe322bf8fa4c15abbdc597dcde&scene=21#wechat_redirect)让我们向上看，当我们的社区走到一起时会有什么可能，而[1.25 版“合作者”](https://mp.weixin.qq.com/s?__biz=MzI5ODk5ODI4Nw==&mid=2247524895&idx=1&sn=69e2122e3dfaecc6291137b866536d70&scene=21#wechat_redirect)是我们社区的共同努力所能做到的，那么 1.26 版“振奋人心的”也是献给所有那些个人行动，整合到发布流中，使所有这些成为可能的人。

## 主要主题

Kubernetes v1.26 由许多变化组成，由全球志愿者团队为你带来。对于这个版本，我们已经确定了几个主要主题。

### 容器镜像注册表的更改

在之前的版本中，Kubernetes 改变了容器注册表\[3]，允许在多个云提供商和地区之间传播负载，这一改变减少了对单个实体的依赖，并为大量用户提供了更快的下载体验。

Kubernetes 的这个版本是第一个专门在新的 registry.k8s.io 容器镜像注册表中发布的版本。在（现在遗留的）k8s.gcr.io 镜像注册表中，不会发布 1.26 版的容器镜像标签，只有 1.26 版之前的版本的标签会继续更新。请参考[k8s.gcr.io -> registry.k8s.io：更快、更便宜，且普遍可用（GA）](https://mp.weixin.qq.com/s?__biz=MzI5ODk5ODI4Nw==&mid=2247528136&idx=2&sn=97332f8ca328619b5670a0a210129b5d&scene=21#wechat_redirect)，了解有关这一重大变化的动机、优势和影响的更多信息。

### CRI v1alpha2 已删除

随着 1.24 版中 CRI\[4]\(Container Runtime Interface，容器运行时接口)的采用和 dockershim 的删除\[5]，CRI 成为 Kubernetes 与不同容器运行时交互的唯一受支持和有文档记录的方式。每个 kubelet 协商在该节点上的容器运行时使用哪个版本的 CRI。

在之前的版本中，Kubernetes 项目推荐使用 CRI 版本 v1，但 kubelet 仍然可以协商使用 CRI 版本 1alpha2，但该版本已被弃用。

Kubernetes v1.26 不再支持 CRI v1alpha2。如果容器运行时不支持 CRI v1，那么这删除\[6]将导致 kubelet 不注册节点。这意味着在 Kubernetes 1.26 中不支持 containerd minor 1.5 及更早的版本；如果你使用 containerd，在将该节点升级到 Kubernetes v1.26 之前，你需要升级到 containerd 版本 1.6.0 或更高版本。这同样适用于任何其他仅支持 v1alpha2 的容器运行时：如果这影响到你，你应该联系容器运行时供应商以获得建议，或者查看他们的网站以获得如何向前发展的其他说明。

### 存储改进

继上一版本中的\[核心 CSI（Container Storage Interface，容器存储接口）迁移]\(core Container Storage Interface "核心 CSI（Container Storage Interface，容器存储接口）迁移")功能正式发布后，CSI 迁移是一项持续的工作，我们现在已经为几个版本做了工作，此版本继续添加（和删除）与迁移目标一致的功能，以及对 Kubernetes 存储的其他改进。

### Azure File 和 vSphere 的 CSI 迁移升级到 Stable

vSphere\[7]和 Azure\[8]的 in-tree 驱动程序到 CSI 的迁移都已升级到 Stable。你可以在\[vSphere CSI 驱动程序]\( "vSphere CSI 驱动程序")https://github.com/kubernetes-sigs/vsphere-csi-driver和Azure File CSI 驱动程序\[9]仓库中找到有关它们的更多信息。

### 将 FSGroup 委托给 CSI 驱动程序升级到 Stable

该功能允许 Kubernetes 在挂载卷时向 CSI 驱动程序提供 pod 的 fsGroup\[10]，以便驱动程序可以利用挂载选项来控制卷权限。以前，kubelet 总是根据 Pod 的.spec.securityContext.fsGroupChangePolicy 字段中指定的策略，将 fsGroup 所有权和权限更改应用于卷中的文件。从此版本开始，CSI 驱动程序可以选择在卷的连接或挂载期间应用 fsGroup 设置。

### 树内 GlusterFS 驱动程序删除

在 1.25 版中已被弃用，在此版本中删除\[11]了 In-tree GlusterFS 驱动程序。

### 树内 OpenStack Cinder 驱动程序删除

此版本删除了已弃用的 OpenStack 树内存储集成（cinder 卷类型）。你应该从https://github.com/kubernetes/cloud-provider-openstack迁移到外部云提供商和CSI驱动程序。有关更多信息，请访问Cinder in-tree to CSI 驱动程序迁移\[12]。

### 签名 Kubernetes 发布工件升级到 Beta

这个特性\[13]是在 Kubernetes v1.24 中引入的，是提高 Kubernetes 发布过程安全性的一个重要里程碑。所有发布工件都使用 cosign\[14]进行无密钥签名，二进制工件和镜像都可以被验证\[15]。

### 对 Windows 特权容器的支持升级到 Stable

特权容器支持允许容器以与直接在主机上运行的进程相似的对主机的访问来运行。Windows 节点中对此功能的支持称为 HostProcess 容器\[16]，现在将升级到 Stable\[17]，允许从特权容器访问主机资源（包括网络资源）。

### Kubernetes 指标的改进

这个版本在指标上有几个值得注意的改进。

### 指标框架扩展升级到 Alpha

指标框架扩展升级到 Alpha\[18]，Kubernetes 代码库中的每个指标都发布了文档\[19]。这一增强为 Kubernetes 指标增加了两个额外的元数据字段：Internal 和 Beta，代表指标成熟度的不同阶段。

### 健康服务水平指标升级到 Alpha

此外，在使用 Kubernetes 指标的能力方面，组件健康服务水平指标\[20]（SLI，Service Level Indicator）已经升级到 Alpha\[21]：通过启用 ComponentSLIs 功能标志，将有一个额外的指标端点，允许从转换为指标格式的原始健康检查数据计算服务水平目标（SLO，Service Level Objective）。

### 功能指标现已可用

现在每个 Kubernetes 组件都有了特性指标，通过检查 kubernetes_feature_enabled 的组件指标端点，可以跟踪每个活跃功能门是否被启用\[22]。

### 动态资源分配升级到 Alpha

动态资源分配\[23]是一项新功能\[24]，它将资源调度交给了第三方开发人员：它为请求访问资源的有限“countable”接口（如 nvidia.com/gpu：2）提供了一种替代方案，提供了一种更类似于持久性卷的 API。在底层，它使用 CDI\[25]（Container Device Interface，容器设备接口）进行设备注入。此功能被 DynamicResourceAllocation 功能门阻止。

### 准入控制的 CEL 升级到 Alpha

这个特性\[26]引入了一个 v1alpha1 API 来验证准入策略\[27]，通过 CEL\[28]（Common Expression Language，公共表达式语言）表达式实现可扩展的接纳控制。目前，自定义策略是通过准入 webhook\[29]来执行的，虽然这种方式很灵活，但与进程内策略执行相比，还是有一些缺点。要使用，请通过--runtime-config 启用 ValidatingAdmissionPolicy 功能门和 admissionregistration.k8s.io/v1alpha1 API。

### Pod 调度改进

Kubernetes v1.26 引入了一些相关的增强功能，可以更好地控制调度行为。

### PodSchedulingReadiness 升级到 Alpha

该特性\[30]在 Pod 的 API 中引入了一个.spec.schedulingGates 字段，以指示是否允许调度 Pod\[31]。外部用户/管理器可以使用此字段根据他们的策略和需求来阻止 Pod 的调度。

### NodeInclusionPolicyInPodTopologySpread 升级到 Beta

通过在 topologySpreadConstraints 中指定 nodeInclusionPolicy，你可以控制在计算 Pod 拓扑分布偏差时是否考虑污点/容差\[32]。

## 其他更新

### 升级到 Stable

此版本总共包括 11 项升级为 Stable 的增强功能：

-

Support for Windows privileged containers\[33]

-

vSphere in-tree to CSI driver migration\[34]

-

Allow Kubernetes to supply pod's fsgroup to CSI driver on mount\[35]

-

Azure file in-tree to CSI driver migration\[36]

-

Job tracking without lingering Pods\[37]

-

Service Internal Traffic Policy\[38]

-

Kubelet Credential Provider\[39]

-

Support of mixed protocols in Services with type=LoadBalancer\[40]

-

Reserve Service IP Ranges For Dynamic and Static IP Allocation\[41]

-

CPUManager\[42]

-

DeviceManager\[43]

### 弃用和移除

在此版本中，Kubernetes 弃用或删除\[44]了 12 个功能。

-

CRI `v1alpha2` API is removed\[45]

-

Removal of the `v1beta1` flow control API group\[46]

-

Removal of the `v2beta2` HorizontalPodAutoscaler API\[47]

-

GlusterFS plugin removed from available in-tree drivers\[48]

-

Removal of legacy command line arguments relating to logging\[49]

-

Removal of `kube-proxy` userspace modes\[50]

-

Removal of in-tree credential management code\[51]

-

The in-tree OpenStack cloud provider is removed\[52]

-

Removal of dynamic kubelet configuration\[53]

-

Deprecation of non-inclusive `kubectl` flag\[54]

-

Deprecations for `kube-apiserver` command line arguments\[55]

-

Deprecations for `kubectl run` command line arguments\[56]

## 发布说明

Kubernetes v1.26 版本的完整细节可以在我们的发行说明\[57]中找到。

## 下载

Kubernetes v1.26 版可以在 Kubernetes 网站\[58]下载。要开始使用 Kubernetes，请查看这些交互式教程\[59]或通过 kind\[60]，使用容器作为“节点”运行本地 Kubernetes 集群。你也可以使用 kubeadm\[61]轻松安装 1.26 版。

## 发布团队

Kubernetes 只有在社区的支持、承诺和努力下才有可能实现。每个发布团队都是由专门的社区志愿者组成的，他们一起工作来构建构成你所依赖的 Kubernetes 版本的许多部分。这需要来自我们社区各个角落的人的专业技能，从代码本身到它的文档和项目管理。

我们要感谢整个发布团队，感谢他们为确保我们为我们的社区发布一个可靠的 Kubernetes v1.26 版本而付出的辛勤工作。

特别感谢我们的发布负责人 Leonard Pahlke，他在整个发布周期中成功地指导了整个发布团队，通过不断的支持和关注构成成功发布之路的众多不同细节，确保我们都能以最佳方式为此次发布做出贡献。

## 用户亮点

-

Wortell\[62]面临着越来越多的开发人员专业知识和日常基础设施管理时间。他们使用 Dapr 来降低复杂性，减少所需的基础设施相关代码的数量，使他们能够将更多的时间放在新功能上。

-

Utmost\[63]处理敏感的个人数据，需要 SOC 2 Type II 认证、ISO 27001 认证和零信任网络。使用 Cilium，他们创建了允许开发人员创建新策略的自动化管道，支持每秒超过 4000 个流量。

-

全球网络安全公司 Ericom\[64]的解决方案依赖于超低延迟和数据安全。借助 Ridge 的托管 Kubernetes 服务，他们能够通过单一 API 部署到全球服务提供商网络。

-

斯堪的纳维亚在线银行 Lunar\[65]希望实施季度生产集群故障转移测试，为灾难恢复做准备，并且需要一种更好的方法来管理他们的平台服务。他们首先集中管理日志管理系统，然后集中所有平台服务，使用 Linkerd 连接集群。

-

Datadog\[66]跨多个云提供商运行 10 多个集群，包含 10,000 多个节点和 100,000 多个 pod。他们转向 Cilium 作为他们的 CNI 和 kube-proxy 的替代品，以利用 eBPF 的强大功能，并为跨任何云的用户提供一致的网络体验。

-

Insiel\[67]希望更新他们的软件生产方法，并在其软件生产中引入云原生范式。他们与 Kiratech 和微软 Azure 的数字化转型项目，使他们能够发展云优先的文化。

## 生态系统消息

-

2023 年欧洲 KubeCon + CloudNativeCon，将于 2023 年 4 月 17 日至 21 日在荷兰阿姆斯特丹举行！你可以在活动网站\[68]上找到有关会议和注册的更多信息。

-

北美 CloudNativeSecurityCon 将于 2023 年 2 月 1 日至 2 日在美国华盛顿州西雅图举办，为期两天，旨在促进云原生安全项目的协作、讨论和知识共享，以及如何最好地利用这些项目来应对安全挑战和机遇。有关更多信息，请参见活动页面\[69]。

-

CNCF 宣布了 2022 年社区奖获奖者\[70]：社区奖旨在表彰 CNCF 社区成员在推进云原生技术方面做出的卓越贡献。

## 项目速度

CNCF K8s DevStats\[71]项目汇总了许多与 Kubernetes 和各种子项目的速度相关的有趣数据点。这包括从个人贡献到做出贡献的公司数量的所有内容，体现了发展这一生态系统的深度和广度。

在为期 14 周\[72]（9 月 5 日至 12 月 9 日）的 v1.26 发布周期中，我们看到了来自 976 家公司\[73]和 6877 名个人\[74]的贡献。

## 即将举办的网络研讨会

加入 Kubernetes v1.26 发布团队的成员，于 2023 年 1 月 17 日星期二美国东部时间 EST 上午 10:00-11:00（世界协调时间 UTC 下午 3:00-4:00）了解此版本的主要功能，以及弃用和删除的内容，以帮助规划升级。有关更多信息和注册，请访问活动页面\[75]。

## 加入

加入 Kubernetes 最简单的方法，是加入与你的兴趣一致的 SIG\[76]（Special Interest Group，特殊兴趣组）。

你有什么想对 Kubernetes 社区说的吗？在我们每周的社区会议\[77]上，以及通过以下渠道分享你的观点：

-

在 Kubernetes 贡献者\[78]网站上找到更多关于贡献给 Kubernetes 的信息

-

在 Twitter @Kubernetesio\[79]上关注我们的最新动态

-

在 Discuss\[80]上加入社区讨论

-

在 Slack\[81]上加入社区

-

在 Server Fault\[82]上发布（或回答）问题

-

分享你的 Kubernetes 故事\[83]

-

在博客\[84]上阅读更多关于 Kubernetes 的信息

-

了解更多关于 Kubernetes 发布团队\[85]的信息

### 参考资料

\[1]

Kubernetes 1.26 发布团队: _<https://github.com/kubernetes/sig-release/blob/master/releases/release-1.26/release-team.md>_

\[2]

开始使用 Projects 来跟踪增强: _<https://github.com/orgs/kubernetes/projects/98/views/1>_

\[3]

Kubernetes 改变了容器注册表: _<https://github.com/kubernetes/kubernetes/pull/109938>_

\[4]

CRI: _<https://kubernetes.io/docs/concepts/architecture/cri/>_

\[5]

dockershim 的删除: _<https://kubernetes.io/blog/2022/02/17/dockershim-faq/>_

\[6]

删除: _<https://github.com/kubernetes/kubernetes/pull/110618>_

\[7]

vSphere: _<https://github.com/kubernetes/enhancements/issues/1491>_

\[8]

Azure: _<https://github.com/kubernetes/enhancements/issues/1885>_

\[9]

Azure File CSI 驱动程序: _<https://github.com/kubernetes-sigs/azurefile-csi-driver>_

\[10]

挂载卷时向 CSI 驱动程序提供 pod 的 fsGroup: _<https://github.com/kubernetes/enhancements/issues/2317>_

\[11]

删除: _<https://github.com/kubernetes/enhancements/issues/3446>_

\[12]

Cinder in-tree to CSI 驱动程序迁移: _<https://github.com/kubernetes/enhancements/issues/1489>_

\[13]

特性: _<https://github.com/kubernetes/enhancements/issues/3031>_

\[14]

cosign: _<https://github.com/sigstore/cosign/>_

\[15]

可以被验证: _<https://kubernetes.io/docs/tasks/administer-cluster/verify-signed-artifacts/>_

\[16]

HostProcess 容器: _<https://kubernetes.io/docs/tasks/configure-pod-container/create-hostprocess-pod/>_

\[17]

升级到 Stable: _<https://github.com/kubernetes/enhancements/issues/1981>_

\[18]

升级到 Alpha: _<https://github.com/kubernetes/enhancements/issues/3498>_

\[19]

Kubernetes 代码库中的每个指标都发布了文档: _<https://kubernetes.io/docs/reference/instrumentation/metrics/>_

\[20]

组件健康服务水平指标: _<https://kubernetes.io/docs/reference/instrumentation/slis/>_

\[21]

升级到 Alpha: _<https://github.com/kubernetes/kubernetes/pull/112884>_

\[22]

跟踪每个活跃功能门是否被启用: _<https://github.com/kubernetes/kubernetes/pull/112690>_

\[23]

动态资源分配: _<https://kubernetes.io/docs/concepts/scheduling-eviction/dynamic-resource-allocation/>_

\[24]

新功能: _<https://github.com/kubernetes/enhancements/blob/master/keps/sig-node/3063-dynamic-resource-allocation/README.md>_

\[25]

CDI: _<https://github.com/container-orchestrated-devices/container-device-interface>_

\[26]

特性: _<https://github.com/kubernetes/enhancements/issues/3488>_

\[27]

验证准入策略: _<https://kubernetes.io/docs/reference/access-authn-authz/validating-admission-policy/>_

\[28]

CEL: _<https://github.com/google/cel-spec>_

\[29]

准入 webhook: _<https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/>_

\[30]

特性: _<https://github.com/kubernetes/enhancements/issues/3521>_

\[31]

指示是否允许调度 Pod: _<https://kubernetes.io/docs/concepts/scheduling-eviction/pod-scheduling-readiness/>_

\[32]

考虑污点/容差: _<https://kubernetes.io/docs/concepts/scheduling-eviction/topology-spread-constraints/>_

\[33]

Support for Windows privileged containers: _<https://github.com/kubernetes/enhancements/issues/1981>_

\[34]

vSphere in-tree to CSI driver migration: _<https://github.com/kubernetes/enhancements/issues/1491>_

\[35]

Allow Kubernetes to supply pod's fsgroup to CSI driver on mount: _<https://github.com/kubernetes/enhancements/issues/2317>_

\[36]

Azure file in-tree to CSI driver migration: _<https://github.com/kubernetes/enhancements/issues/1885>_

\[37]

Job tracking without lingering Pods: _<https://github.com/kubernetes/enhancements/issues/2307>_

\[38]

Service Internal Traffic Policy: _<https://github.com/kubernetes/enhancements/issues/2086>_

\[39]

Kubelet Credential Provider: _<https://github.com/kubernetes/enhancements/issues/2133>_

\[40]

Support of mixed protocols in Services with type=LoadBalancer: _<https://github.com/kubernetes/enhancements/issues/1435>_

\[41]

Reserve Service IP Ranges For Dynamic and Static IP Allocation: _<https://github.com/kubernetes/enhancements/issues/3070>_

\[42]

CPUManager: _<https://github.com/kubernetes/enhancements/issues/3570>_

\[43]

DeviceManager: _<https://github.com/kubernetes/enhancements/issues/3573>_

\[44]

弃用或删除: _<https://kubernetes.io/blog/2022/11/18/upcoming-changes-in-kubernetes-1-26/>_

\[45]

CRI `v1alpha2` API is removed: _<https://github.com/kubernetes/kubernetes/pull/110618>_

\[46]

Removal of the `v1beta1` flow control API group: _<https://kubernetes.io/docs/reference/using-api/deprecation-guide/#flowcontrol-resources-v126>_

\[47]

Removal of the `v2beta2` HorizontalPodAutoscaler API: _<https://kubernetes.io/docs/reference/using-api/deprecation-guide/#horizontalpodautoscaler-v126>_

\[48]

GlusterFS plugin removed from available in-tree drivers: _<https://github.com/kubernetes/enhancements/issues/3446>_

\[49]

Removal of legacy command line arguments relating to logging: _<https://github.com/kubernetes/kubernetes/pull/112120>_

\[50]

Removal of `kube-proxy` userspace modes: _<https://github.com/kubernetes/kubernetes/pull/112133>_

\[51]

Removal of in-tree credential management code: _<https://github.com/kubernetes/kubernetes/pull/112341>_

\[52]

The in-tree OpenStack cloud provider is removed: _<https://github.com/kubernetes/enhancements/issues/1489>_

\[53]

Removal of dynamic kubelet configuration: _<https://github.com/kubernetes/kubernetes/pull/112643>_

\[54]

Deprecation of non-inclusive `kubectl` flag: _<https://github.com/kubernetes/kubernetes/pull/113116>_

\[55]

Deprecations for `kube-apiserver` command line arguments: _<https://github.com/kubernetes/kubernetes/pull/38186>_

\[56]

Deprecations for `kubectl run` command line arguments: _<https://github.com/kubernetes/kubernetes/pull/112261>_

\[57]

发行说明: _<https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.26.md>_

\[58]

Kubernetes 网站: _<https://k8s.io/releases/download/>_

\[59]

交互式教程: _<https://kubernetes.io/docs/tutorials/>_

\[60]

kind: _<https://kind.sigs.k8s.io/>_

\[61]

kubeadm: _<https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/>_

\[62]

Wortell: _<https://www.cncf.io/case-studies/wortell/>_

\[63]

Utmost: _<https://www.cncf.io/case-studies/utmost/>_

\[64]

Ericom: _<https://www.cncf.io/case-studies/ericom/>_

\[65]

Lunar: _<https://www.cncf.io/case-studies/lunar/>_

\[66]

Datadog: _<https://www.cncf.io/case-studies/datadog/>_

\[67]

Insiel: _<https://www.cncf.io/case-studies/insiel/>_

\[68]

活动网站: _<https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/>_

\[69]

活动页面: _<https://events.linuxfoundation.org/cloudnativesecuritycon-north-america/>_

\[70]

2022 年社区奖获奖者: _<https://www.cncf.io/announcements/2022/10/28/cloud-native-computing-foundation-reveals-2022-community-awards-winners/>_

\[71]

CNCF K8s DevStats: _<https://k8s.devstats.cncf.io/d/12/dashboards?orgId=1&refresh=15m>_

\[72]

为期 14 周: _<https://github.com/kubernetes/sig-release/tree/master/releases/release-1.26>_

\[73]

976 家公司: _[https://k8s.devstats.cncf.io/d/9/companies-table?orgId=1\&var-period_name=v1.25.0 - v1.26.0\&var-metric=contributions](https://k8s.devstats.cncf.io/d/9/companies-table?orgId=1&var-period_name=v1.25.0%20-%20v1.26.0&var-metric=contributions)_

\[74]

6877 名个人: _[https://k8s.devstats.cncf.io/d/66/developer-activity-counts-by-companies?orgId=1\&var-period_name=v1.25.0 - v1.26.0\&var-metric=contributions\&var-repogroup_name=Kubernetes\&var-country_name=All\&var-companies=All\&var-repo_name=kubernetes%2Fkubernetes](https://k8s.devstats.cncf.io/d/66/developer-activity-counts-by-companies?orgId=1&var-period_name=v1.25.0%20-%20v1.26.0&var-metric=contributions&var-repogroup_name=Kubernetes&var-country_name=All&var-companies=All&var-repo_name=kubernetes%2Fkubernetes)_

\[75]

活动页面: _<https://community.cncf.io/events/details/cncf-cncf-online-programs-presents-cncf-live-webinar-kubernetes-v126-release/>_

\[76]

SIG: _<https://github.com/kubernetes/community/blob/master/sig-list.md>_

\[77]

社区会议: _<https://github.com/kubernetes/community/tree/master/communication>_

\[78]

Kubernetes 贡献者: _<https://www.kubernetes.dev/>_

\[79]

@Kubernetesio: _<https://twitter.com/kubernetesio>_

\[80]

Discuss: _<https://discuss.kubernetes.io/>_

\[81]

Slack: _<http://slack.k8s.io/>_

\[82]

Server Fault: _<https://serverfault.com/questions/tagged/kubernetes>_

\[83]

Kubernetes 故事: _<https://docs.google.com/a/linuxfoundation.org/forms/d/e/1FAIpQLScuI7Ye3VQHQTwBASrgkjQDSS5TP0g3AXfFhwSM9YpHgxRKFA/viewform>_

\[84]

博客: _<https://kubernetes.io/blog/>_

\[85]

Kubernetes 发布团队: _<https://github.com/kubernetes/sig-release/tree/master/release-team>_

**\_\***CNCF\***\*\*\*\***（\_\_**_云原生计算基金会_**）致力于培育和维护一个厂商中立的开源生态系统，来推广云原生技术。我们通过将最前沿的模式民主化，让这些创新为大众所用。请长按以下二维码进行关注。\_\*\*
