---
title: Kubernetes 1.23：下一个前沿
---

原文链接：<https://mp.weixin.qq.com/s/SiFhhDHHjHNbg8QPX4fOBw>
作者：Kubernetes 1.23 发布团队\[1]

我们很高兴地宣布发布 Kubernetes 1.23，2021 年的最后一个版本！

这个版本包含 47 个增强：11 个增强升级到稳定版，17 个增强升级到 beta 版，19 个增强进入 alpha 版。另外，有一个特性已经被弃用。

## 版本主题

### 弃用 FlexVolume

FlexVolume 已被弃用。out-of-tree 驱动程序是在 Kubernetes 中编写卷驱动程序的推荐方法。更多信息请参见此文档\[2]。FlexVolume 驱动程序的维护者应该实现一个 CSI 驱动程序，并将 FlexVolume 的用户转移到 CSI。FlexVolume 的用户应该将他们的工作负载转移到 CSI 驱动程序。

### 弃用 klog 特定标志

为了简化代码库，一些日志标记\[3]在 Kubernetes 1.23 中被标记为弃用。实现它们的代码将在未来的版本中删除，因此用户需要开始用一些替代解决方案替换已弃用的标志。

### Kubernetes 发布过程中的软件供应链 SLSA 1 级遵从性

Kubernetes 版本现在生成描述发布过程的登台和发布阶段的出处认证文件。工件现在在从一个阶段移交到下一个阶段时得到验证。这最后一部分完成了遵守 SLSA 安全框架\[4]的第 1 级（软件工件的供应链级别）所需的工作。

### IPv4/IPv6 双栈组网向 GA 过渡

IPv4/IPv6 双栈网络\[5]演化为 GA。从 1.21 开始，Kubernetes 集群默认支持双栈网络。在 1.23 中，IPv6DualStack 特性门被移除。并不是必须使用双栈组网。虽然集群支持双堆栈网络，但 Pods 和 Services 继续默认为单堆栈。要使用双栈组网，Kubernetes 节点必须有可路由的 IPv4/IPv6 网络接口，必须使用一个双栈功能的 CNI 网络插件，pod 必须配置为双栈，Services 必须有它们的.spec.ipFamilyPolicy 字段设置为 PreferDualStack 或 RequireDualStack。

### HorizontalPodAutoscaler v2 毕业到 GA

在 1.23 版本中，HorizontalPodAutscaler autoscaling/v2 稳定 API 移动到了 GA。HorizontalPodAutoscaler autoscaling/v2beta2 API 已被弃用。

### 通用临时卷特性毕业到 GA

通用临时卷特性在 1.23 中转移到了 GA。该特性允许任何支持动态供应的现有存储驱动程序作为临时卷使用，将卷的生命周期绑定到 Pod。支持卷创建的所有 StorageClass 参数和 PersistentVolumeClaims 支持的所有功能。

### 跳过卷所有权毕业到 GA

配置 pod 的卷权限和所有权变更策略的特性在 1.23 迁移到 GA。这允许用户跳过挂载时的递归权限更改，并加快了 pod 的启动时间。

### 允许 CSI 驱动程序选择加入卷所有权和许可毕业到 GA

允许 CSI 驱动程序声明支持基于 fsGroup 的权限的特性在 1.23 中毕业到 GA。

### PodSecurity 进入测试阶段

PodSecurity\[6]升级到测试版。PodSecurity 替换已弃用的 PodSecurityPolicy 接纳控制器。PodSecurity 是一个允许控制器，它根据设置实施级别的特定名称空间标签对命名空间中的 Pods 实施 Pod 安全标准。在 1.23 中，PodSecurity 特性门默认是启用的。

### 容器运行时接口（CRI） v1 是默认值

Kubelet 现在支持 CRI v1 API，它现在是项目范围内的默认 API。如果容器运行时不支持 v1 API，Kubernetes 将退回到 v1alpha2 实现。用户不需要操作，因为 v1 和 v1alpha2 在实现上没有区别。很有可能 v1alpha2 将在未来 Kubernetes 的某个版本中被删除，以便能够开发 v1。

### 结构化日志到 Beta 版本

结构化日志达到了 Beta 里程碑。来自 kubelet 和 kube-scheduler 的大多数日志消息已被转换。鼓励用户尝试 JSON 输出或结构化文本格式的解析，并就开放问题的可能解决方案提供反馈，例如在日志值中处理多行字符串。

### 简化多点插件配置的调度程序

kube-scheduler 为插件添加了一个新的、简化的配置字段，以允许在一个点上启用多个扩展点。新的 multiPoint 插件字段旨在为管理员简化大多数调度程序的设置。通过 multiPoint 启用的插件将自动为它们实现的每个扩展点注册。例如，一个实现了 Score 和 Filter 扩展的插件可以同时为两者启用。这意味着可以启用或禁用整个插件，而无需手动编辑单个扩展点设置。这些扩展点现在可以被抽象掉，因为它们与大多数用户无关。

### CSI 迁移更新

CSI Migration 支持替换现有的树内存储插件，如 kubernetes.io/gce-pd 或 kubernetes.io/aws-ebs 与相应的 CSI 驱动程序。如果 CSI Migration 正常工作，Kubernetes 的最终用户应该不会注意到任何区别。迁移之后，Kubernetes 用户可能继续依赖于使用现有接口的树内存储插件的所有功能。

- CSI 迁移功能在默认情况下是开启的，但在 GCE PD、AWS EBS 和 Azure Disk 1.23 版本中仍处于测试状态。
- CSI 迁移在 1.23 中对 Ceph RBD 和 Portworx 引入 Alpha 特性。

### 用于 CRD 的表达式语言验证是 alpha

用于 CRD 的表达式语言验证的 alpha 版本从 1.23 开始。如果启用了 CustomResourceValidationExpressions 特性门，则自定义资源将通过使用公共表达式语言（CEL）\[7]的验证规则进行验证。

### 服务器端字段验证是 Alpha

从 1.23 开始启用了 ServerSideFieldValidation 特性门，当用户在请求中发送包含未知或重复字段的 Kubernetes 对象时，将收到来自服务器的警告。以前未知的字段和除最后一个重复字段外的所有字段都将被服务器删除。

启用了特性门后，我们还引入了 fieldValidation 查询参数，以便用户可以根据每个请求指定服务器的期望行为。验证字段查询参数的有效值是：

- Ignore（特性门被禁用时的默认值，与 1.23 之前删除/忽略未知字段的行为相同）
- Warn（当启用特性门时的默认值）
- Strict（这将使请求失败并显示无效请求错误）

### OpenAPI v3 是 Alpha

OpenAPIV3 特性门从 1.23 开始启用，用户将能够为所有 Kubernetes 类型请求 OpenAPI v3.0 规范。OpenAPI v3 的目标是完全透明，并包括对一组字段的支持，这些字段在发布 OpenAPI v2 时被删除：default、nullable、oneOf、anyOf。每个 Kubernetes 组版本（在 cluster/openapi/v3 路径中找到。

## 毕业到稳定

- IPv4/IPv6 Dual-Stack Support
- Skip Volume Ownership Change
- TTL After Finished Controller
- Config FSGroup Policy in CSI Driver object
- Generic Ephemeral Inline Volumes
- Defend Against Logging Secrets via Static Analysis
- Namespace Scoped Ingress Class Parameters
- Reducing Kubernetes Build Maintenance
- Graduate HPA API to GA

## 重大变化

- Priority and Fairness for API Server Requests

## 发布说明

在我们的发布说明\[8]中查看 Kubernetes 1.23 版本的全部细节。

## 下载

Kubernetes 1.23 可以在 GitHub\[9]下载。要开始使用 Kubernetes，请查看这些交互式教程\[10]或通过 kind\[11]使用 Docker 容器“节点”运行本地 Kubernetes 集群。您还可以使用 kubeadm\[12]轻松安装 1.23。

## 发布团队

这个版本是由一群非常敬业的人组成的，他们组成了一个团队，交付技术内容、文档、代码和许多其他组件，这些组件包含在每个 Kubernetes 版本中。

非常感谢发行团队的负责人 Rey Lejano 带领我们走过了一个成功的发行周期，也非常感谢发行团队中所有人相互支持，如此努力地为社区交付 1.23 发行版。

## 发布主题及标志

### Kubernetes 1.23：下一个前沿

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/a65b48a0-c265-4605-bdd2-d124ba827ade/640)

“下一个前沿（The Next Frontier）”主题代表了 1.23 中新的和逐步增强的功能，Kubernetes 关于 Star Trek 参考文献的历史，以及发行团队中社区成员的成长。

Kubernetes 有 Star Trek 的历史。在谷歌中 Kubernetes 最初的代号是 Project 7，参考了 Star Trek Voyager 中的 Seven of Nine。当然 Borg 是 Kubernetes 的前代的名字。“The Next Frontier”主题延续了 Star Trek 的风格。“The Next Frontier”是两个 Star Trek 标题的融合，Star Trek V: The Final Frontier 和 Star Trek the Next Generation。

“下一个前沿”代表了 SIG 发布章程中的一条线，“确保有一组一致的社区成员在适当的地方支持跨时间的发布过程。”在每个发布团队中，我们都有新的发布团队成员来发展社区，对于许多人来说，这是他们在开源领域的第一次贡献。

参考：

- <https://kubernetes.io/blog/2015/04/borg-predecessor-to-kubernetes/>
- <https://github.com/kubernetes/community/blob/master/sig-release/charter.md>

Kubernetes 1.23 版本的 logo 沿用了 Star Trek 的主题。每一颗星都是 Kubernetes 标志上的头盔。这艘船代表了释放团队的集体合作。

Rey Lejano 设计了标志。

## 用户亮点

- 最新的 CNCF 最终用户技术雷达围绕 DevSecOps\[13]展开。请查看雷达页面\[14]，以获得完整的细节和发现。
- 了解最终用户 Aegon Life India\[15]如何将其核心流程从传统的整体架构迁移到基于微服务的架构，以努力转型为领先的数字服务公司。
- 利用多个云原生项目，Seagate\[16]设计了 edgerX 在边缘上运行实时分析。
- 看看 Zambon\[17]是如何与 SparkFabrik 合作开发 16 个网站的，使用云原生技术，使利益相关者能够轻松更新内容，同时保持一致的品牌标识。
- 使用 Kubernetes，InfluxData\[18] 能够通过创建真正的云抽象层来实现多云、多区域服务可用性的承诺，该层允许将 InfluxDB 作为单个应用程序无缝交付到三个主要云提供商的多个全球集群。

## 生态系统更新

2021 年北美 KubeCon + CloudNativeCon\[19]于 2021 年 10 月举行，包括在线和现场。所有的演讲现在都可以按需提供\[20]给任何想要赶上的人！

- Kubernetes 和云原生要素培训和 KCNA 认证现在一般可供报名和安排\[21]。此外，一个新的在线培训课程，Kubernetes 和 Cloud Native Essentials（LFS250）\[22]，已经发布，以帮助个人准备入门级的云角色和参加 KCNA 考试。
- 包容性命名计划提供了新的资源\[23]，包括开放源代码包容性战略（LFC103）课程、语言评估框架和实现路径。

## 项目速度

CNCF K8s 的 DevStats 项目\[24]汇集了许多与 Kubernetes 和各种子项目的速度相关的有趣数据点。这包括了从个人贡献到正在贡献的公司数量的方方面面，这也说明了在这个生态系统的发展过程中所付出的努力的深度和广度。

在持续 16 周（8 月 23 日至 12 月 7 日）的 v1.23 发布周期中，我们看到了来自 1032 家公司\[25]和 1084 位个人\[26]的贡献。

## 活动更新

2021 年中国 KubeCon + CloudNativeCon\[27]将于本月 9 日至 10 日举行。继去年中断后，今年的活动将采用虚拟形式，包括 105 个环节。在这里\[28]查看活动时间表。

- 2022 年欧洲 KubeCon + CloudNativeCon 将于 2022 年 5 月 4 日至 7 日在西班牙巴伦西亚举行！您可以在活动网站\[29]上找到更多关于会议和注册的信息。
- Kubernetes 社区日即将在巴基斯坦、巴西、成都和澳大利亚举行活动。

## 版本网络研讨会

在 2022 年 1 月 4 日加入 Kubernetes 1.23 发布团队的成员，了解该发布的主要特性，以及弃用和移除，以帮助计划升级。欲了解更多信息和注册信息，请访问 CNCF 在线项目网站的活动页面\[30]。

## 参与

加入 Kubernetes 最简单的方法是加入许多与你的兴趣相一致的特殊兴趣小组（SIGs）之一。你有什么想向 Kubernetes 社区广播的吗？在我们每周的社区会议上分享你的声音，并通过以下渠道：

- 在 Kubernetes 贡献者\[31]网站上找到更多关于为 Kubernetes 贡献的信息
- 关注我们的 Twitter @Kubernetesio，了解最新消息
- 加入 Discuss\[32]的社区讨论
- 在 Slack\[33]上加入社区
- 在 Stack Overflow\[34]上发布问题（或回答问题）
- 分享你的 Kubernetes 故事\[35]
- 在博客\[36]上阅读更多关于 Kubernetes 发生的事情
- 了解更多关于 Kubernetes 发布团队\[37]的信息

### 参考资料

\[1]

Kubernetes 1.23 发布团队: [_https://github.com/kubernetes/sig-release/blob/master/releases/release-1.23/release-team.md_](https://github.com/kubernetes/sig-release/blob/master/releases/release-1.23/release-team.md)

\[2]

文档: [_https://github.com/kubernetes/community/blob/master/sig-storage/volume-plugin-faq.md#kubernetes-volume-plugin-faq-for-storage-vendors_](https://github.com/kubernetes/community/blob/master/sig-storage/volume-plugin-faq.md#kubernetes-volume-plugin-faq-for-storage-vendors)

\[3]

一些日志标记: [_https://kubernetes.io/docs/concepts/cluster-administration/system-logs/#klog_](https://kubernetes.io/docs/concepts/cluster-administration/system-logs/#klog)

\[4]

SLSA 安全框架: [_https://slsa.dev/_](https://slsa.dev/)

\[5]

IPv4/IPv6 双栈网络: [_https://github.com/kubernetes/enhancements/tree/master/keps/sig-network/563-dual-stack_](https://github.com/kubernetes/enhancements/tree/master/keps/sig-network/563-dual-stack)

\[6]

PodSecurity: [_https://kubernetes.io/docs/concepts/security/pod-security-admission/_](https://kubernetes.io/docs/concepts/security/pod-security-admission/)

\[7]

公共表达式语言（CEL）: [_https://github.com/google/cel-spec_](https://github.com/google/cel-spec)

\[8]

发布说明: [_https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.23.md_](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.23.md)

\[9]

GitHub: [_https://github.com/kubernetes/kubernetes/releases/tag/v1.23.0_](https://github.com/kubernetes/kubernetes/releases/tag/v1.23.0)

\[10]

交互式教程: [_https://kubernetes.io/docs/tutorials/_](https://kubernetes.io/docs/tutorials/)

\[11]

kind: [_https://kind.sigs.k8s.io/_](https://kind.sigs.k8s.io/)

\[12]

kubeadm: [_https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/_](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/)

\[13]

DevSecOps: [_https://www.cncf.io/announcements/2021/09/22/cncf-end-user-technology-radar-provides-insights-into-devsecops/_](https://www.cncf.io/announcements/2021/09/22/cncf-end-user-technology-radar-provides-insights-into-devsecops/)

\[14]

雷达页面: [_https://radar.cncf.io/_](https://radar.cncf.io/)

\[15]

Aegon Life India: [_https://www.cncf.io/case-studies/aegon-life-india/_](https://www.cncf.io/case-studies/aegon-life-india/)

\[16]

Seagate: [_https://www.cncf.io/case-studies/seagate/_](https://www.cncf.io/case-studies/seagate/)

\[17]

Zambon: [_https://www.cncf.io/case-studies/zambon/_](https://www.cncf.io/case-studies/zambon/)

\[18]

InfluxData: [_https://www.cncf.io/case-studies/influxdata/_](https://www.cncf.io/case-studies/influxdata/)

\[19]

2021 年北美 KubeCon + CloudNativeCon: [_https://www.cncf.io/events/kubecon-cloudnativecon-north-america-2021/_](https://www.cncf.io/events/kubecon-cloudnativecon-north-america-2021/)

\[20]

按需提供: [_https://www.youtube.com/playlist?list=PLj6h78yzYM2Nd1U4RMhv7v88fdiFqeYAP_](https://www.youtube.com/playlist?list=PLj6h78yzYM2Nd1U4RMhv7v88fdiFqeYAP)

\[21]

Kubernetes 和云原生要素培训和 KCNA 认证现在一般可供报名和安排: [_https://www.cncf.io/announcements/2021/11/18/kubernetes-and-cloud-native-essentials-training-and-kcna-certification-now-available/_](https://www.cncf.io/announcements/2021/11/18/kubernetes-and-cloud-native-essentials-training-and-kcna-certification-now-available/)

\[22]

Kubernetes 和 Cloud Native Essentials（LFS250）: [_https://www.cncf.io/announcements/2021/10/13/entry-level-kubernetes-certification-to-help-advance-cloud-careers/_](https://www.cncf.io/announcements/2021/10/13/entry-level-kubernetes-certification-to-help-advance-cloud-careers/)

\[23]

包容性命名计划提供了新的资源: [_https://www.cncf.io/announcements/2021/10/13/inclusive-naming-initiative-announces-new-community-resources-for-a-more-inclusive-future/_](https://www.cncf.io/announcements/2021/10/13/inclusive-naming-initiative-announces-new-community-resources-for-a-more-inclusive-future/)

\[24]

CNCF K8s 的 DevStats 项目: [_https://k8s.devstats.cncf.io/d/12/dashboards?orgId=1\&refresh=15m_](https://k8s.devstats.cncf.io/d/12/dashboards?orgId=1&refresh=15m)

\[25]

1032 家公司: [_https://k8s.devstats.cncf.io/d/9/companies-table?orgId=1\&var-period_name=v1.22.0 - now\&var-metric=contributions_](https://k8s.devstats.cncf.io/d/9/companies-table?orgId=1&var-period_name=v1.22.0%20-%20now&var-metric=contributions)

\[26]

1084 位个人: [_https://k8s.devstats.cncf.io/d/66/developer-activity-counts-by-companies?orgId=1\&var-period_name=v1.22.0 - now\&var-metric=contributions\&var-repogroup_name=Kubernetes\&var-country_name=All\&var-companies=All\&var-repo_name=kubernetes%2Fkubernetes_](https://k8s.devstats.cncf.io/d/66/developer-activity-counts-by-companies?orgId=1&var-period_name=v1.22.0%20-%20now&var-metric=contributions&var-repogroup_name=Kubernetes&var-country_name=All&var-companies=All&var-repo_name=kubernetes%2Fkubernetes)

\[27]

2021 年中国 KubeCon + CloudNativeCon: [_https://www.lfasiallc.com/kubecon-cloudnativecon-open-source-summit-china/_](https://www.lfasiallc.com/kubecon-cloudnativecon-open-source-summit-china/)

\[28]

这里: [_https://www.lfasiallc.com/kubecon-cloudnativecon-open-source-summit-china/program/schedule/_](https://www.lfasiallc.com/kubecon-cloudnativecon-open-source-summit-china/program/schedule/)

\[29]

活动网站: [_https://events.linuxfoundation.org/archive/2021/kubecon-cloudnativecon-europe/_](https://events.linuxfoundation.org/archive/2021/kubecon-cloudnativecon-europe/)

\[30]

活动页面: [_https://community.cncf.io/e/mrey9h/_](https://community.cncf.io/e/mrey9h/)

\[31]

Kubernetes 贡献者: [_https://www.kubernetes.dev/_](https://www.kubernetes.dev/)

\[32]

Discuss: [_https://discuss.kubernetes.io/_](https://discuss.kubernetes.io/)

\[33]

Slack: [_http://slack.k8s.io/_](http://slack.k8s.io/)

\[34]

Stack Overflow: [_http://stackoverflow.com/questions/tagged/kubernetes_](http://stackoverflow.com/questions/tagged/kubernetes)

\[35]

故事: [_https://docs.google.com/a/linuxfoundation.org/forms/d/e/1FAIpQLScuI7Ye3VQHQTwBASrgkjQDSS5TP0g3AXfFhwSM9YpHgxRKFA/viewform_](https://docs.google.com/a/linuxfoundation.org/forms/d/e/1FAIpQLScuI7Ye3VQHQTwBASrgkjQDSS5TP0g3AXfFhwSM9YpHgxRKFA/viewform)

\[36]

博客: [_https://kubernetes.io/blog/_](https://kubernetes.io/blog/)

\[37]

Kubernetes 发布团队: [_https://github.com/kubernetes/sig-release/tree/master/release-team_](https://github.com/kubernetes/sig-release/tree/master/release-team)

**\_\***CNCF\***\*\*\*\***（\_\_**_云原生计算基金会_**）致力于培育和维护一个厂商中立的开源生态系统，来推广云原生技术。我们通过将最前沿的模式民主化，让这些创新为大众所用。请长按以下二维码进行关注。\_\*\*
