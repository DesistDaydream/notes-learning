---
title: Kubernetes 1.22 版本发布2021年8月5日
---

<https://mp.weixin.qq.com/s/TQoyU3S2q4Q-yCK0zGHo9g>
来源：CNCF
作者：**Kubernetes 1.22 发布团队\[1]**
我们很高兴地宣布 Kubernetes 1.22 的发布，这是 2021 年的第二个版本！
这个版本包含 53 个增强功能：13 个增强功能已经升级到稳定版，24 个增强功能正在进入 beta 版，16 个增强功能正在进入 alpha 版。另外，有三个特性已被弃用。
今年 4 月，Kubernetes 的发布节奏正式从每年 4 个版本改为 3 个版本。这是第一个与该变更相关的长周期版本。随着 Kubernetes 项目的成熟，每个周期的增强数量也在增长。对于贡献者社区和版本工程团队来说，这意味着从一个版本到另一个版本要做更多的工作，这也会给最终用户社区带来压力，要求他们不断更新包含越来越多特性的版本。
将发布周期从每年 4 个版本更改为 3 个版本平衡了项目的许多方面，包括如何管理贡献和发布，以及社区计划升级和保持最新的能力。
你可以在 Kubernetes 官方博客[Kubernetes 发布节奏改变：这是你需要知道的](https://mp.weixin.qq.com/s?__biz=MzI5ODk5ODI4Nw==&mid=2247504376&idx=1&sn=28157c22be9ff28d6af4498fe73b8fdd&scene=21#wechat_redirect)阅读更多信息。

## 主题

### 服务器端应用毕业到 GA

服务器端应用（Server-side Apply）是 Kubernetes API 服务器上运行的一个新的字段所有权和对象合并算法。服务器端应用程序通过声明性配置帮助用户和控制器管理其资源。它允许他们通过声明的方式创建和/或修改他们的对象，只需要发送他们完全指定的意图。经过几个测试版本后，服务器端应用程序现在普遍可用。

### 外部凭据提供者现在稳定了

对 Kubernetes 客户端证书插件的支持从 1.11 开始就处于测试阶段，随着 Kubernetes 1.22 的发布，现在已经稳定下来了。GA 特性集包括对提供交互式登录流的插件的改进支持，以及一些 bug 修复。有兴趣的插件作者可以看看**sample-exec-plugin\[2]**入门。

### etcd 移到 3.5.0

Kubernetes 的默认后端存储 etcd 有了新的版本：3.5.0。新版本对安全性、性能、监控和开发人员体验进行了改进。有许多 bug 修复和一些关键的新特性，比如迁移到结构化日志记录和内置日志旋转。该版本附带了一个详细的未来路线图，以实现解决交通超载问题的解决方案。你可以在**3.5.0 发布公告\[3]**中阅读完整而详细的更改列表。

### 内存资源的服务质量

最初，Kubernetes 使用的是 v1 cgroups API。通过这种设计，Pod 的 QoS 类只应用于 CPU 资源（如 cpu_shares）。作为一个 alpha 特性，Kubernetes v1.22 现在可以使用 cgroups v2 API 来控制内存分配和隔离。此特性旨在在存在内存资源争用时提高工作负载和节点可用性，并提高容器生命周期的可预测性。

### 节点系统 swap 支持

在设置和使用 Kubernetes 时，每个系统管理员或 Kubernetes 用户的处境都是一样的：禁用 swap 空间。随着 Kubernetes 1.22 的发布，alpha 支持可用来运行带有 swap 内存的节点。此更改允许管理员选择在 Linux 节点上配置 swap，将块存储的一部分作为额外的虚拟内存处理。

### Windows 增强和功能

为了继续支持不断增长的开发人员社区，SIG Windows 发布了他们的**开发环境\[4]**。这些新工具支持多个 CNI 提供程序，可以在多个平台上运行。还有一种新的方法可以从头开始运行最新的 Windows 特性，即编译 Windows kubelet 和 kube-proxy，然后将它们与其他 Kubernetes 组件的日常构建一起使用。
CSI 对 Windows 节点的支持在 1.22 版本中转移到了 GA。在 Kubernetes v1.22 中，Windows 特权容器是一个 alpha 特性。为了允许在 Windows 节点上使用 CSI 存储，**CSIProxy\[5]**支持将 CSI 节点插件部署为非特权 pod，使用代理在节点上执行特权存储操作。

### seccomp 的默认配置文件

kubelet 中增加了默认 seccomp 配置文件的 alpha 特性，以及新的命令行标志和配置。在使用时，这个新特性提供了集群范围的 seccomp 缺省值，默认情况下使用 RuntimeDefault seccomp 配置文件而不是 Unconfined。这增强了 Kubernetes 部署的默认安全性。安全管理员知道工作负载在默认情况下更安全，现在可以睡得更好。要了解更多的特性，请参考官方**seccomp 教程\[6]**。

### 使用 kubeadm 更安全的控制平面

一个新的 alpha 特性允许以非 root 用户的身份运行 kubeadm 控制平面组件。这是 kubeadm 长期以来要求的安全措施。要尝试它，你必须启用 kubeadm 特定的 RootlessControlPlane 功能门。当你使用这个 alpha 特性部署集群时，你的控制平面以较低的权限运行。
Kubernetes 1.22 还给 kubeadm 带来了新的**v1beta3 配置 API\[7]**。这个迭代添加了一些长期要求的特性，并弃用了一些现有的特性。v1beta3 版本现在是首选的 API 版本；v1beta2 API 仍然可用，尚未被弃用。

## 重大变化

### 移除几个已弃用的 beta api

一些被弃用的 beta API 已经在 1.22 中被移除，以支持这些 API 的 GA 版本。所有现有对象都可以通过稳定的 API 进行交互。该移除包括 Ingress、IngressClass、Lease、APIService、ValidatingWebhookConfiguration、MutatingWebhookConfiguration、CustomResourceDefinition、TokenReview、SubjectAccessReview 和 CertificateSigningRequest API 的 beta 版本。
要了解完整的列表，请查看**已弃用 API 迁移指南\[8]**以及博文**Kubernetes 1.22 版本将删除的 API 和特性：这里是你需要知道的\[9]**。

### 临时容器的 API 更改和改进

用于创建**临时容器\[10]**的 API 在 1.22 中有更改。临时容器特性是 alpha 特性，默认情况下是禁用的，并且新 API 不能与试图使用旧 API 的客户端一起工作。
对于稳定特性，kubectl 工具遵循 Kubernetes**版本倾斜策略\[11]**；然而，kubectl v1.21 及更早版本不支持临时容器的新 API。如果你计划使用 kubectl debug 来创建临时容器，并且你的集群正在运行 Kubernetes v1.22，那么你不能使用 kubectl v1.21 或更早的版本来这样做。请将 kubectl 升级到 1.22，如果你希望将 kubectl 调试与集群版本混合使用。

## 其他的更新

### 毕业到稳定

- Bound Service Account Token Volumes
- CSI Service Account Token
- Windows Support for CSI Plugins
- Warning mechanism for deprecated API use
- PodDisruptionBudget Eviction

### 显著特点更新

- 引入了一个新的 PodSecurity 准入 alpha 特性，旨在替代 PodSecurityPolicy
- 内存管理器进入测试版
- 启用 API Server Tracing 的新 alpha 特性
- kubeadm 配置格式的一个新的 v1beta3 版本
- 用于 PersistentVolumes 的通用数据填充器现在在 alpha 版本中可用
- Kubernetes 控制平面现在将始终使用 CronJobs v2 控制器
- 作为 alpha 特性，所有 Kubernetes 节点组件（包括 kubelet、kube-proxy 和容器运行时）都可以作为非 root 用户运行

## 版本说明

你可以在**版本说明\[12]**中查看 1.22 发布的完整细节。

## 下载

Kubernetes 1.22 可供**下载\[13]**，也可以在**GitHub\[14]**项目中获得。
有一些很好的资源可以帮助你开始使用 Kubernetes。你可以在 Kubernetes 主站点上查看一些**交互式教程\[15]**，或者通过**kind\[16]**使用 Docker 容器在你的机器上运行一个本地集群。如果你想尝试从头开始构建集群，请参阅 Kelsey Hightower 编写的**Kubernetes the Hard Way\[17]**教程。

## 发布团队

这个发行版是由一群非常敬业的个人组成的，他们组成了一个团队，交付 Kubernetes 发行版中的技术内容、文档、代码和许多其他组件。
非常感谢发行负责人 Savitha Raghunathan 带领我们度过了一个成功的发行周期，也感谢发行团队中所有人的相互支持，以及为向社区发布 1.22 版本所付出的努力。
我们也想借此机会纪念 Peeyush Gupta，他是我们今年早些时候失去的团队成员。Peeyush 积极参与 SIG ContribEx 和 Kubernetes 发布团队，最近担任 1.22 通信负责人。他的贡献和努力将继续反映在他帮助建立的社区中。已经创建了一个**CNCF 纪念页面\[18]**，社区可以在这里分享想法和记忆。

## 版本标志

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/om7q43/1628135795073-ad4f9a51-8a78-4918-9d08-f169af41606d.png)
在持续的流行病、自然灾害和始终存在的倦怠阴影中，Kubernetes 1.22 版本包括 53 个增强。这是迄今为止最大的版本。这一成就的实现要归功于努力工作、充满激情的发布团队成员以及 Kubernetes 生态系统中令人惊叹的贡献者。版本标志是我们不断达到新的里程碑和创造新的记录的提醒。它是献给所有的发行团队成员，徒步旅行者和天文爱好者的！
该标志由**Boris Zotkin\[19]**设计。Boris 是 MathWorks 的 Mac/Linux 管理员。他喜欢生活中简单的事情，喜欢和家人在一起。这个精通技术的人总是准备好迎接挑战，并乐于帮助朋友！

## 用户亮点

- 今年 5 月，CNCF 在全球欢迎了 27 个新组织，成为多样的云原生生态系统的成员。这些新成员将参加 CNCF 的活动，包括即将于 2021 年 10 月 12 日至 15 日在洛杉矶举行的北美 KubeCon + CloudNativeCon。
- CNCF 在 2021 年欧洲 KubeCon + CloudNativeCon 虚拟大会上授予 Spotify 最高最终用户奖。

## 项目速度

**CNCF K8s DevStats 项目\[20]**聚集了许多与 Kubernetes 的开发速度和各种子项目相关的有趣数据点。这包括从个人贡献到贡献的公司数量的所有内容，并说明了这个生态系统的发展所付出的努力的深度和广度。
在运行了 15 周（4 月 26 日至 8 月 4 日）的 v1.22 发布周期中，我们看到了来自 1063 家公司和 2054 个人的贡献。

## 生态系统更新

- **2021 年欧洲 KubeCon + CloudNativeCon 虚拟大会\[21]**已于 5 月举行，这是第三次虚拟活动。所有的讲座现在都可以**按需提供\[22]**给任何想赶上进度的人！
- [春季 LFX 项目](https://mp.weixin.qq.com/s?__biz=MzI5ODk5ODI4Nw==&mid=2247504107&idx=2&sn=4789ee5a0d156af0fefc4165f4434435&scene=21#wechat_redirect)拥有最大的毕业班，有 28 名成功的 CNCF 实习生！
- CNCF 在今年年初在 Twitch 上推出了直播，目标是为任何想要在云原生社区学习、成长和与他人合作的人提供明确的互动媒体体验。

## 活动更新

- **2021 年北美 KubeCon + CloudNativeCon 大会\[23]**将于 2021 年 10 月 12 日至 15 日在洛杉矶举行！你可以在活动网站上找到关于会议和注册的更多信息。
- **Kubernetes 社区日\[24]**即将在意大利、英国和华盛顿特区举行活动。

## 1.22 版本网络研讨会

请于 2021 年 9 月 7 日加入 Kubernetes 1.22 发布团队，了解该版本的主要特性，以及弃用和删除，以帮助制定升级计划。更多信息和注册请访问 CNCF 在线项目网站的**活动页面\[25]**。

## 加入参与

如果你有兴趣为 Kubernetes 社区做出贡献，那么特殊兴趣小组（SIG）是一个很好的起点。他们中的许多人可能与你的兴趣一致！如果你有什么想与社区分享的东西，你可以参加每周社区会议，或使用以下任何渠道：

- 在**Kubernetes Contributor\[26]**网站了解更多关于为 Kubernetes 做贡献的信息。
- 关注我们的推特\[]@Kubernetesio]\(https://twitter.com/kubernetesio "]@Kubernetesio")获取最新消息
- 加入**Discuss\[27]**社区的讨论
- 加入**Slack\[28]**社区
- 分享你的 Kubernetes**故事\[29]**
- 在**博客\[30]**上阅读更多关于 Kubernetes 的情况
- 了解更多关于**Kubernetes 发布团队\[31]**的信息

### 参考资料

\[1]\[2]\[3]\[4]\[5]\[6]\[7]\[8]\[9]\[10]\[11]\[12]\[13]\[14]\[15]\[16]\[17]\[18]\[19]\[20]\[21]\[22]\[23]\[24]\[25]\[26]\[27]\[28]\[29]\[30]\[31]
Kubernetes 1.22 发布团队: *https://github.com/kubernetes/sig-release/blob/master/releases/release-1.22/release-team.md*
sample-exec-plugin: *https://github.com/ankeesler/sample-exec-plugin*
3.5.0 发布公告: *https://etcd.io/blog/2021/announcing-etcd-3.5/*
开发环境: *https://github.com/kubernetes-sigs/sig-windows-dev-tools/*
CSIProxy: *https://github.com/kubernetes-csi/csi-proxy*
seccomp 教程: *https://kubernetes.io/docs/tutorials/clusters/seccomp/#enable-the-use-of-runtimedefault-as-the-default-seccomp-profile-for-all-workloads*
v1beta3 配置 API: *https://kubernetes.io/docs/reference/config-api/kubeadm-config.v1beta3/*
已弃用 API 迁移指南: *https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-22*
Kubernetes 1.22 版本将删除的 API 和特性：这里是你需要知道的: *https://blog.k8s.io/2021/07/14/upcoming-changes-in-kubernetes-1-22/*
临时容器: *https://kubernetes.io/docs/concepts/workloads/pods/ephemeral-containers/*
版本倾斜策略: *https://kubernetes.io/releases/version-skew-policy/*
版本说明: *https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.22.md*
下载: *https://kubernetes.io/releases/download/*
GitHub: *https://github.com/kubernetes/kubernetes/releases/tag/v1.22.0*
交互式教程: *https://kubernetes.io/docs/tutorials/*
kind: *https://kind.sigs.k8s.io/*
Kubernetes the Hard Way: *https://github.com/kelseyhightower/kubernetes-the-hard-way*
CNCF 纪念页面: *https://github.com/cncf/memorials/blob/main/peeyush-gupta.md*
Boris Zotkin: *https://www.instagram.com/boris.z.man/*
CNCF K8s DevStats 项目: *https://k8s.devstats.cncf.io/*
2021 年欧洲 KubeCon + CloudNativeCon 虚拟大会: *https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/*
按需提供: *https://www.youtube.com/playlist?list=PLj6h78yzYM2MqBm19mRz9SYLsw4kfQBrC*
2021 年北美 KubeCon + CloudNativeCon 大会: *https://events.linuxfoundation.org/kubecon-cloudnativecon-north-america/*
Kubernetes 社区日: *https://community.cncf.io/kubernetes-community-days/about-kcd/*
活动页面: *https://community.cncf.io/events/details/cncf-cncf-online-programs-presents-cncf-live-webinar-kubernetes-122-release/*
Kubernetes Contributor: *https://www.kubernetes.dev/*
Discuss: *https://discuss.kubernetes.io/*
Slack: *http://slack.k8s.io/*
故事: *https://docs.google.com/a/linuxfoundation.org/forms/d/e/1FAIpQLScuI7Ye3VQHQTwBASrgkjQDSS5TP0g3AXfFhwSM9YpHgxRKFA/viewform*
博客: *https://kubernetes.io/blog/*
Kubernetes 发布团队: *https://github.com/kubernetes/sig-release/tree/master/release-team*
