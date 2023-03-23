---
title: Kubernetes 1.25：合作者
---

<https://mp.weixin.qq.com/s/6nhv2zQIAOAfUJ661YmDsQ>
<https://mp.weixin.qq.com/s/6yrd1Dtf0wA0Ixu-Fr4MNQ>
<https://mp.weixin.qq.com/s/Rn0A8SzLJvPJbwAIhHuRbQ>
删除 PodSecurityPolicy
弃用 GlusterFS、Portworx In-Tree 类型的内置卷，正在删除 Flocker、Quobyte、StorageOS
Kubelet 将逐渐走向不在 nat 表中创建以下 iptables 链：

- KUBE-MARK-DROP
- KUBE-MARK-MASQ
- KUBE-POSTROUTING

原文链接：<https://mp.weixin.qq.com/s/6yrd1Dtf0wA0Ixu-Fr4MNQ>
作者：Kubernetes 1.25 发布团队\[1]

宣布发布 Kubernetes v1.25！

这个版本总共包括 40 个增强功能。这些增强功能中有 15 个正在进入 Alpha，10 个正在升级到 Beta，13 个正在升级到 Stable。我们也有两个功能被弃用或移除。

## 版本主题和徽标

Kubernetes 1.25：合作者

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/5b4ad302-0e91-4420-88da-e875d0d0f25f/640)

Kubernetes v1.25 的主题是合作者。

Kubernetes 项目本身由许多许多单独的组件组成，当这些组件组合在一起时，就形成了你今天看到的项目形式。它也是由许多个人构建和维护的，他们都有不同的技能、经验、历史和兴趣，他们不仅作为发布团队，还作为许多 SIG 全年支持项目和社区。

在这个版本中，我们希望向协作、开放的精神致敬，这种精神将我们从分散在全球的孤立的开发者、作者和用户带到了一个能够改变世界的联合力量中。Kubernetes v1.25 包含了惊人的 40 项增强功能，如果没有我们共同努力所获得的惊人力量，任何一项功能都不会存在。

受我们发布负责人的儿子 Albert Song 的启发，Kubernetes v1.25 以你们每一个人的命名，无论你们如何选择为成为 Kubernetes 的合力而贡献自己独特的力量。

## 新功能（主要主题）

### PodSecurityPolicy 已移除；Pod Security Admission 升级到 Stable

PodSecurityPolicy 最初在 1.21 版中被弃用，随着 1.25 版的发布，它已被移除。提高其可用性所需的更新会带来突破性的变化，因此有必要移除它以支持更友好的替代。该替换是 Pod Security Admission，随着此版本升级到 Stable。如果你当前依赖于 PodSecurityPolicy，请按照说明迁移到 Pod Security Admission\[2]。

### 临时容器升级到 Stable

Ephemeral Containers\[3]（临时容器）是在现有 pod 中仅存在有限时间的容器。当你需要检查另一个容器，但由于该容器已经崩溃或其镜像缺少调试工具而无法使用 kubectl exec 时，这对于故障排除特别有用。在 Kubernetes v1.23 中，临时容器升级到了 Beta，而在这个版本中，这个特性升级到了 Stable。

### cgroups v2 的支持升级到 Stable

Linux 内核 cgroups v2 API 宣布稳定已经两年多了。现在一些发行版默认使用这个 API，Kubernetes 必须支持它才能继续在这些发行版上运行。cgroups v2 提供了对 cgroups v1 的几项改进，有关更多信息，请参见 cgroups v2 文档。虽然 cgroups v1 将继续得到支持，但这一增强让我们为其最终被弃用和取代做好了准备。

### 改进的 Windows 支持

- 性能仪表盘\[4]增加了对 Windows 的支持
- 单元测试\[5]增加了对 Windows 的支持
- 一致性测试\[6]增加了对 Windows 的支持
- 为 Windows 操作就绪\[7]创建新的 GitHub 存储库

### 将容器注册服务从 k8s.gcr.io 移动到 registry.k8s.io

将容器注册中心从 k8s.gcr.io 移动到 registry.k8s.io 给合并。有关更多详细信息，请参见 wiki 页面，公告已发送至 kubernetes 开发邮件列表。

### 已将 SeccompDefault 升级为 Beta

SeccompDefault 升级为 beta，更多细节请参见教程用 seccomp 限制容器的系统调用\[8]。

### 将网络策略中的 endPort 升级为 Stable

将 Network Policy\[9]（网络策略）中的 endPort 升级为 GA。支持 endPort 字段的网络策略提供程序现在可以使用它来指定应用网络策略的端口范围。以前，每个网络策略只能针对一个端口。

请注意，网络策略提供程序必须支持 endPort 字段。如果你的提供商不支持 endPort，并且在网络策略中指定了此字段，则创建的网络策略将仅包含端口字段（单个端口）。

### 将本地临时存储容量隔离升级为 Stable

本地临时存储容量隔离\[10]（Local Ephemeral Storage Capacity Isolation）功能已升级为 GA。这是在 1.8 中作为 alpha 引入的，在 1.10 中移到了 beta，现在是一个 Stable 的特性。它支持 pod 之间的本地临时存储的容量隔离，例如 EmptyDir，这样，如果 pod 的本地临时存储的消耗超过该限制，则可以通过驱逐 pod 来硬限制 pod 对共享资源的消耗。

### 将核心 CSI Migration 升级到 Stable

CSI Migration\[11]是 SIG Storage 在几个版本中一直在进行的工作。目标是将树内卷插件移动到树外 CSI 驱动程序，并最终移除树内卷插件。核心 CSI Migration\[12]功能已升级为 GA。GCE PD 和 AWS EBS 的 CSI Migration 也升级到 GA。vSphere 的 CSI Migration 仍处于 beta 阶段（但默认情况下处于开启状态）。Portworx 的 CSI Migration 已到 beta（但默认关闭）。

### 将 CSI 临时卷升级到 Stable

CSI 临时卷\[13]（CSI Ephemeral Volume）特性允许在 pod 规范中为临时用例直接指定 CSI 卷。使用挂载卷的方式直接在 pod 内部注入任意状态，例如配置、秘密、身份、变量或类似的信息。这最初是在 1.15 中作为 alpha 特性引入的，现在升级为 GA。此功能由一些 CSI 驱动程序使用，如 secret-store CSI 驱动程序。

### 将 CRD 验证表达式语言升级为 Beta

CRD 验证表达式语言\[14]（CRD Validation Expression Language）升级为 beta，这使得可以声明如何使用 CEL\[15]（Common Expression Language，公共表达式语言）验证自定义资源。请参阅验证规则指南\[16]。

### 将 ServerSideFieldValidation 升级为 Beta

将 ServerSideFieldValidation 功能门升级为 beta（默认情况下打开）。这允许选择性地在 API 服务器上触发模式验证，当检测到未知字段时会出错。这允许从 kubectl 中移除客户端验证，同时保持相同的核心功能，即在包含未知或无效字段的请求中出错。

### 引入 KMS v2 API

引入 KMS v2alpha1 API 以增加性能、旋转和可观察性改进。对于 kms 数据加密，使用 AES-GCM 而不是 AES-CBC，通过 DEK 对静态数据（即 Kubernetes Secrets）进行加密。不需要用户操作。将继续允许使用 AES-GCM 和 AES-CBC 进行读取。有关详细信息，请参阅 KMS provider\[17]指南。

## 其他更新

### 升级到 Stable

此版本总共包括 13 项升级为 Stable 的增强功能：

- Ephemeral Containers
- Local Ephemeral Storage Resource Management
- CSI Ephemeral Volumes
- CSI Migration - Core
- Graduate the kube-scheduler ComponentConfig to GA
- CSI Migration - AWS
- CSI Migration - GCE
- DaemonSets Support MaxSurge
- NetworkPolicy Port Range
- cgroups v2
- Pod Security Admission
- Add minReadySeconds to Statefulsets
- Identify Windows pods at API admission level authoritatively

### 弃用和移除

在此版本中，Kubernetes 弃用或移除了两个特性。

- PodSecurityPolicy 已移除
- GlusterFS 插件已从可用的树内驱动程序中弃用

### 发布说明

Kubernetes v1.25 版本的完整细节可以在我们的版本说明\[18]中找到。

### 下载

Kubernetes v1.25 可以在 GitHub 下载\[19]。要开始使用 Kubernetes，请查看交互式教程\[20]或通过 kind\[21]使用容器作为“节点”运行本地 Kubernetes 集群。你也可以使用 kubeadm\[22]轻松安装 1.25。

### 发布团队

Kubernetes 只有在社区的支持、承诺和努力下才有可能实现。每个发布团队都由专门的社区志愿者组成，他们一起工作来构建许多部分，当这些部分结合在一起时，就构成了你所依赖的 Kubernetes 发布。这需要来自我们社区各个角落的人的专业技能，从代码本身到它的文档和项目管理。

我们要感谢整个发布团队，感谢他们为确保我们为我们的社区发布一个可靠的 Kubernetes v1.25 版本而付出的辛勤工作。你们每一个人都在建设这个过程中发挥了作用，你们都做得很好。我们要特别感谢无所畏惧的发行负责人 Cici Huang，为保证我们获得成功所做的一切。

### 用户亮点

- Finleap Connect\[23]在高度监管的环境中运营。2019 年，他们有五个月的时间在其集群的所有服务中实施 mTLS（mutual TLS），以使其业务代码符合新的欧洲 PSD2 支付指令。
- PNC\[24]寻求开发一种方法来确保新代码自动满足安全标准和审计合规性要求，从而取代他们现有的繁琐的 30 天手动流程。使用 Knative，PNC 开发了内部工具来自动检查新代码和对现有代码的更改。
- Nexxiot\[25]需要高度可靠、安全、高性能且经济高效的 Kubernetes 集群。他们将 Cilium 作为 CNI 来锁定他们的集群，并获得可靠的第二天操作实现弹性网络。
- 由于创建网络保险单的过程是一个复杂的多步骤过程，At-Bay\[26]试图通过使用基于异步消息的通信模式/设施来改善运营。他们认为 Dapr 满足了它所期望的一系列要求，甚至更多。

### 生态系统更新

- 2022 年北美 KubeCon + CloudNativeCon 将于 2022 年 10 月 24 日至 28 日在密歇根州底特律举行！你可以在活动网站上找到有关会议和注册的更多信息。
- KubeDay 系列活动将于 12 月 7 日在日本 KubeDay 拉开帷幕！在活动网站上注册或提交提案。
- 在 2021 年的云原生调查中，CNCF 见证了创纪录的 Kubernetes 和容器采用。看一下调查结果\[27]。

### 项目速度

CNCF K8s DevStats\[28]项目汇总了许多与 Kubernetes 和各子项目的速度相关的有趣数据点。这包括从个人贡献到做出贡献的公司数量的所有内容，体现了发展这一生态系统的深度和广度。

在为期 15 周（5 月 23 日至 8 月 23 日）的 v1.25 发布周期中，我们看到了来自 1065 家公司和 1620 名个人的贡献。

## 即将举办的发布网络研讨会

加入 Kubernetes v1.25 发布团队的成员，于 2022 年 9 月 22 日星期四上午 10：00–11：00 PT（太平洋时间）了解此版本的主要功能，以及反对意见和移除内容，以帮助规划升级。有关更多信息和注册，请访问活动页面\[29]。

## 加入

加入 Kubernetes 最简单的方法，是加入与你的兴趣一致的 SIG\[30]（Special Interest Group，特殊兴趣小组）。你有什么想对 Kubernetes 社区说的吗？在我们每周的社区会议\[31]上，以及通过以下渠道分享你的观点：

- 在 Kubernetes 贡献者\[32]网站上找到更多关于贡献给 Kubernetes 的信息
- 在 Twitter @Kubernetesio\[33]上关注我们的最新动态
- 加入 Discuss\[34]的社区讨论
- 加入 Slack\[35]社区
- 在 Server Fault\[36]提出问题（或回答问题）
- 分享你的 Kubernetes 故事\[37]
- 在博客\[38]上阅读更多关于 Kubernetes 的信息
- 了解更多关于 Kubernetes 发布团队\[39]的信息

### 参考资料

\[1]

Kubernetes 1.25 发布团队: [_https://github.com/kubernetes/sig-release/blob/master/releases/release-1.25/release-team.md_](https://github.com/kubernetes/sig-release/blob/master/releases/release-1.25/release-team.md)

\[2]

迁移到 Pod Security Admission: [_https://kubernetes.io/docs/tasks/configure-pod-container/migrate-from-psp/_](https://kubernetes.io/docs/tasks/configure-pod-container/migrate-from-psp/)

\[3]

Ephemeral Containers: [_https://kubernetes.io/docs/concepts/workloads/pods/ephemeral-containers/_](https://kubernetes.io/docs/concepts/workloads/pods/ephemeral-containers/)

\[4]

性能仪表盘: [_http://perf-dash.k8s.io/#/?jobname=soak-tests-capz-windows-2019_](http://perf-dash.k8s.io/#/?jobname=soak-tests-capz-windows-2019)

\[5]

单元测试: [_https://github.com/kubernetes/kubernetes/issues/51540_](https://github.com/kubernetes/kubernetes/issues/51540)

\[6]

一致性测试: [_https://github.com/kubernetes/kubernetes/pull/108592_](https://github.com/kubernetes/kubernetes/pull/108592)

\[7]

Windows 操作就绪: [_https://github.com/kubernetes-sigs/windows-operational-readiness_](https://github.com/kubernetes-sigs/windows-operational-readiness)

\[8]

用 seccomp 限制容器的系统调用: [_https://kubernetes.io/docs/tutorials/security/seccomp/#enable-the-use-of-runtimedefault-as-the-default-seccomp-profile-for-all-workloads_](https://kubernetes.io/docs/tutorials/security/seccomp/#enable-the-use-of-runtimedefault-as-the-default-seccomp-profile-for-all-workloads)

\[9]

Network Policy: [_https://kubernetes.io/docs/concepts/services-networking/network-policies/#targeting-a-range-of-ports_](https://kubernetes.io/docs/concepts/services-networking/network-policies/#targeting-a-range-of-ports)

\[10]

本地临时存储容量隔离: [_https://github.com/kubernetes/enhancements/tree/master/keps/sig-storage/361-local-ephemeral-storage-isolation_](https://github.com/kubernetes/enhancements/tree/master/keps/sig-storage/361-local-ephemeral-storage-isolation)

\[11]

CSI Migration: [_https://kubernetes.io/blog/2021/12/10/storage-in-tree-to-csi-migration-status-update/#quick-recap-what-is-csi-migration-and-why-migrate_](https://kubernetes.io/blog/2021/12/10/storage-in-tree-to-csi-migration-status-update/#quick-recap-what-is-csi-migration-and-why-migrate)

\[12]

核心 CSI Migration: [_https://github.com/kubernetes/enhancements/tree/master/keps/sig-storage/625-csi-migration_](https://github.com/kubernetes/enhancements/tree/master/keps/sig-storage/625-csi-migration)

\[13]

CSI 临时卷: [_https://github.com/kubernetes/enhancements/tree/master/keps/sig-storage/596-csi-inline-volumes_](https://github.com/kubernetes/enhancements/tree/master/keps/sig-storage/596-csi-inline-volumes)

\[14]

CRD 验证表达式语言: [_https://github.com/kubernetes/enhancements/blob/master/keps/sig-api-machinery/2876-crd-validation-expression-language/README.md_](https://github.com/kubernetes/enhancements/blob/master/keps/sig-api-machinery/2876-crd-validation-expression-language/README.md)

\[15]

CEL: [_https://github.com/google/cel-spec_](https://github.com/google/cel-spec)

\[16]

验证规则指南: [_https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#validation-rules_](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#validation-rules)

\[17]

KMS provider: [_https://kubernetes.io/docs/tasks/administer-cluster/kms-provider/_](https://kubernetes.io/docs/tasks/administer-cluster/kms-provider/)

\[18]

版本说明: [_https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.25.md_](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.25.md)

\[19]

下载: [_https://github.com/kubernetes/kubernetes/releases/tag/v1.25.0_](https://github.com/kubernetes/kubernetes/releases/tag/v1.25.0)

\[20]

交互式教程: [_https://kubernetes.io/docs/tutorials/_](https://kubernetes.io/docs/tutorials/)

\[21]

kind: [_https://kind.sigs.k8s.io/_](https://kind.sigs.k8s.io/)

\[22]

kubeadm: [_https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/_](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/)

\[23]

Finleap Connect: [_https://www.cncf.io/case-studies/finleap-connect/_](https://www.cncf.io/case-studies/finleap-connect/)

\[24]

PNC: [_https://www.cncf.io/case-studies/pnc-bank/_](https://www.cncf.io/case-studies/pnc-bank/)

\[25]

Nexxiot: [_https://www.cncf.io/case-studies/nexxiot/_](https://www.cncf.io/case-studies/nexxiot/)

\[26]

At-Bay: [_https://www.cncf.io/case-studies/at-bay/_](https://www.cncf.io/case-studies/at-bay/)

\[27]

调查结果: [_https://www.cncf.io/reports/cncf-annual-survey-2021/_](https://www.cncf.io/reports/cncf-annual-survey-2021/)

\[28]

CNCF K8s DevStats: [_https://k8s.devstats.cncf.io/d/12/dashboards?orgId=1\&refresh=15m_](https://k8s.devstats.cncf.io/d/12/dashboards?orgId=1&refresh=15m)

\[29]

活动页面: [_https://community.cncf.io/events/details/cncf-cncf-online-programs-presents-cncf-live-webinar-kubernetes-v125-release/_](https://community.cncf.io/events/details/cncf-cncf-online-programs-presents-cncf-live-webinar-kubernetes-v125-release/)

\[30]

SIG: [_https://github.com/kubernetes/community/blob/master/sig-list.md_](https://github.com/kubernetes/community/blob/master/sig-list.md)

\[31]

社区会议: [_https://github.com/kubernetes/community/tree/master/communication_](https://github.com/kubernetes/community/tree/master/communication)

\[32]

Kubernetes 贡献者: [_https://www.kubernetes.dev/_](https://www.kubernetes.dev/)

\[33]

@Kubernetesio: [_https://twitter.com/kubernetesio_](https://twitter.com/kubernetesio)

\[34]

Discuss: [_https://discuss.kubernetes.io/_](https://discuss.kubernetes.io/)

\[35]

Slack: [_http://slack.k8s.io/_](http://slack.k8s.io/)

\[36]

Server Fault: [_https://serverfault.com/questions/tagged/kubernetes_](https://serverfault.com/questions/tagged/kubernetes)

\[37]

故事: [_https://docs.google.com/a/linuxfoundation.org/forms/d/e/1FAIpQLScuI7Ye3VQHQTwBASrgkjQDSS5TP0g3AXfFhwSM9YpHgxRKFA/viewform_](https://docs.google.com/a/linuxfoundation.org/forms/d/e/1FAIpQLScuI7Ye3VQHQTwBASrgkjQDSS5TP0g3AXfFhwSM9YpHgxRKFA/viewform)

\[38]

博客: [_https://kubernetes.io/blog/_](https://kubernetes.io/blog/)

\[39]

Kubernetes 发布团队: [_https://github.com/kubernetes/sig-release/tree/master/release-team_](https://github.com/kubernetes/sig-release/tree/master/release-team)

**\_\***CNCF\***\*\*\*\***（\_\_**_云原生计算基金会_**）致力于培育和维护一个厂商中立的开源生态系统，来推广云原生技术。我们通过将最前沿的模式民主化，让这些创新为大众所用。请长按以下二维码进行关注。\_\*\*
