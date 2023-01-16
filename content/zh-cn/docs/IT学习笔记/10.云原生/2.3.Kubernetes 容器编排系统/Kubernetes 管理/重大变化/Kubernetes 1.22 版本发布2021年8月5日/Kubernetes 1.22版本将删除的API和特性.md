---
title: Kubernetes 1.22版本将删除的API和特性
---

<https://mp.weixin.qq.com/s/PTPTmR6Dprqsyd6E-hBvlw>
<https://kubernetes.io/blog/2021/07/14/upcoming-changes-in-kubernetes-1-22/>

作者：Krishna Kilari（AWS），Tim Bannister（The Scale Factory）
随着 Kubernetes API 的发展，API 会周期性地重新组织或升级。当 API 演进时，它们所取代的旧 API 将被弃用，并最终被移除。请参阅**Kubernetes API removals\[1]**阅读更多关于 Kubernetes 移除 API 的政策。
我们想确保你知道一些即将的删除。这些是测试版（beta）的 API，你可以在当前受支持的 Kubernetes 版本中使用，它们已经被弃用了。所有这些删除的原因是它们已经被一个更新的、稳定的（“GA”）API 所取代。
Kubernetes 1.22 将于 2021 年 8 月发布，它将移除一些已经弃用的 API。**Kubernetes 1.22 发布信息\[2]**详细介绍了 v1.22 的发布时间表。

## Kubernetes v1.22 移除的 API

v1.22 版本将停止提供下面列出的 API 版本。这些都是 beta 版本的 API，之前为了支持更新和更稳定的 API 版本而被弃用。

- Beta versions of the ValidatingWebhookConfiguration and MutatingWebhookConfiguration API (the admissionregistration.k8s.io/v1beta1 API versions)
- The beta CustomResourceDefinition API (apiextensions.k8s.io/v1beta1)
- The beta APIService API (apiregistration.k8s.io/v1beta1)
- The beta TokenReview API (authentication.k8s.io/v1beta1)
- Beta API versions of SubjectAccessReview, LocalSubjectAccessReview, SelfSubjectAccessReview (API versions from authorization.k8s.io/v1beta1)
- The beta CertificateSigningRequest API (certificates.k8s.io/v1beta1)
- The beta Lease API (coordination.k8s.io/v1beta1)
- All beta Ingress APIs (the extensions/v1beta1 and networking.k8s.io/v1beta1 API versions)

Kubernetes 文档涵盖了这些**1.22 版本中删除的 API\[3]**，并解释了这些 API 在 beta 和稳定版本之间是如何变化的。

## 要做什么

我们将浏览每一种受这些删除影响的资源，并解释你需要采取的步骤。

### Ingress

迁移至使用 networking.k8s.io/v1 Ingress API，从 v1.19 开始可用。相关的 API IngressClass 是为了补充 Ingress 概念而设计的，允许你在一个集群中配置多种 Ingress。如果你目前正在使用已弃用的 kubernetes.io/ingress.class 注释，请计划改用.spec.ingressClassName 字段。
在任何运行 Kubernetes v1.19 或更高版本的集群上，你都可以使用 v1 API 来检索或更新现有的 Ingress 对象，即使它们是使用较旧的 API 版本创建的。
当你将一个 Ingress 转换为 v1 API 时，你应该检查该 Ingress 中的每个规则。较老版本的 Ingress 使用遗留的 ImplementationSpecific 路径类型。切换路径匹配 Prefix 或 Exact，而不是 ImplementationSpecific。迁移到这些替代路径类型的好处之一是，在不同的 Ingress 类之间迁移变得更加容易。
除了作为客户端升级自己对 Ingress API 的使用之外，还要确保你使用的每个 Ingress 控制器都与 v1 Ingress API 兼容。阅读**Ingress 先决条件\[4]**了解更多关于 Ingress 和 Ingress 控制器的上下文。

#### 更新 NGINX-Ingress 以使用稳定的 Ingress API

<https://mp.weixin.qq.com/s/hVTWlfrqmjZRrb0KTsDrZA>
作者：James Strong、Ricardo Katz
对于所有 Kubernetes API，都有一个创建、维护和最终弃用它们的过程。networking.k8s.io API 组没有不同。即将发布的 Kubernetes 1.22 版本将删除几个与网络相关的弃用 API：

- **IngressClass\[1]**的 networking.k8s.io/v1beta1 API 版本
- **Ingress\[2]**的所有测试版本：extensions/v1beta1 和 networking.k8s.io/v1beta1

在 v1.22 Kubernetes 集群中，你可以通过稳定的（v1）API 访问 Ingress 和 IngressClass 对象，但无法通过 beta API 访问它们。自**2017 年\[3]**、**2019 年\[4]**1.16 Kubernetes API 弃用以来，以及最近的是 KEP-1453：**从 Ingress API 毕业到 GA\[5]**，一直在讨论中这种变化。
在社区会议上，SIG Network 决定使用 Ingress-NGINX 版本 0.47.0，继续支持 Kubernetes 1.22 之前的版本。Kubernetes 1.22 发布后，对 Ingress-NGINX 的支持将持续 6 个月。对于 Ingress-NGINX 的任何额外的 bug 修复和 CVE 都将根据需要进行处理。
Ingress-NGINX 将有独立的分支和 Ingress-NGINX 版本来支持这个模型，这反映了 Kubernetes 项目过程。Ingress-NGINX 项目的未来版本将跟踪并支持 Kubernetes 的最新版本。

| Kubernetes version | Ingress-NGINX version | Notes                                                                                                                            |
| ------------------ | --------------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| v1.22              | v1.0.0-alpha.2        | New features, plus bug fixes.                                                                                                    |
| v1.21              | v0.47.x               | Bugfixes only, and just for security issues or crashes. No end-of-support date announced.                                        |
| v1.20              | v0.47.x               | Bugfixes only, and just for security issues or crashes. No end-of-support date announced.                                        |
| v1.19              | v0.47.x               | Bugfixes only, and just for security issues or crashes. Fixes only provided until 6 months after Kubernetes v1.22.0 is released. |

由于 Kubernetes 1.22 中的更新，v0.47.0 将不能与 Kubernetes 1.22 一起工作。

#### 你需要做什么

团队目前正在升级 ingress-nginx 以支持 v1 迁移，你可以在这里**跟踪进度\[6]**。
在完成对 Ingress v1 的支持之前，我们不会对 ingress-nginx 进行特性改进。
同时，确保没有兼容性问题：

- 更新到最新版本的 Ingress-NGINX；目前是**v0.47.0\[7]**
- Kubernetes 1.22 发布后，确保你使用的是最新版本的 Ingress-NGINX，支持 Ingress 和 IngressClass 的稳定 API。
- 在集群版本>=1.19 的环境测试 Ingress-NGINX 版本 v1.0.0-alpha.2，并报告任何问题到项目 Github 页面。

欢迎社区对这项工作的反馈和支持。Ingress-NGINX 子项目定期举行社区会议，讨论这个项目面临的问题和其他问题。有关子项目的更多信息，请参阅**SIG Network\[8]**。

### 参考资料

\[1]\[2]\[3]\[4]\[5]\[6]\[7]\[8]
IngressClass: *https://kubernetes.io/docs/concepts/services-networking/ingress/#ingress-class*
Ingress: *https://kubernetes.io/docs/concepts/services-networking/ingress/*
2017 年: *https://github.com/kubernetes/kubernetes/issues/43214*
2019 年: *https://kubernetes.io/blog/2019/07/18/api-deprecations-in-1-16/*
从 Ingress API 毕业到 GA: *https://github.com/kubernetes/enhancements/tree/master/keps/sig-network/1453-ingress-api#122*
跟踪进度: *https://github.com/kubernetes/ingress-nginx/pull/7156*
v0.47.0: *https://github.com/kubernetes/ingress-nginx/releases/tag/controller-v0.47.0*
SIG Network: *https://github.com/kubernetes/community/tree/master/sig-network*

### ValidatingWebhookConfiguration 和 MutatingWebhookConfiguration

迁移以使用 admissionregistration.k8s.io/v1 API 的 ValidatingWebhookConfiguration 和 MutatingWebhookConfiguration，从 v1.16 开始提供。你可以使用 v1 API 来检索或更新现有对象，即使它们是使用较旧的 API 版本创建的。

### CustomResourceDefinition

迁移以使用 CustomResourceDefinition apiextensions.k8s.io/v1 API，从 v1.16 开始可用。你可以使用 v1 API 来检索或更新现有对象，即使它们是使用较旧的 API 版本创建的。如果你在集群中定义了任何自定义资源，那么在升级之后这些资源仍然会被使用。
如果你正在使用外部 CustomResourceDefinitions，则可以使用 kubectl convert 将现有清单转换为使用较新的 API。因为在 beta 版和稳定版 CustomResourceDefinitions 之间存在一些功能差异，所以我们的建议是对每一个进行测试，以确保它在升级后按照你期望的方式工作。

### APIService

迁移到使用 apiregistration.k8s.io/v1 APIService API，从 v1.10 开始可用。你可以使用 v1 API 来检索或更新现有对象，即使它们是使用较旧的 API 版本创建的。如果你已经有了使用 APIService 对象的 API 聚合，那么这个聚合在升级之后将继续工作。

### TokenReview

迁移以使用 authentication.k8s.io/v1 TokenReview API，从 v1.10 开始可用。除了通过 HTTP 提供这个 API 之外，Kubernetes API 服务器也使用相同的格式将 TokenReviews 发送到 webhook。v1.22 版本继续使用 v1beta1 API 发送 TokenReviews 给 webhook。关于切换到稳定 API 的一些具体技巧，请参阅展望未来。

### SubjectAccessReview、SelfSubjectAccessReview 和 LocalSubjectAccessReview

迁移以使用 authorization.k8s.io/v1 版本的 API，从 v1.6 开始提供。

### CertificateSigningRequest

迁移到使用 certificates.k8s.io/v1 CertificateSigningRequest API，从 v1.19 开始可用。你可以使用 v1 API 来检索或更新现有对象，即使它们是使用较旧的 API 版本创建的。现有颁发的证书在升级时保留其有效性。

### Lease

迁移以使用 coordination.k8s.io/v1 Lease API，从 v1.14 开始可用。你可以使用 v1 API 来检索或更新现有对象，即使它们是使用较旧的 API 版本创建的。

### kubectl convert

kubectl 有一个插件，提供 kubectl convert 子命令。这是一个官方插件，你可以下载作为 Kubernetes 的一部分。有关更多细节，请参阅**Download Kubernetes\[5]**。
你可以使用 kubectl convert 来更新清单文件，以使用不同的 API 版本。例如，如果你在源代码控制中有一个使用 beta Ingress API 的清单，你可以 check out 这个定义，并运行 kubectl convert -f--output-version/。你可以使用 kubectl convert 命令自动转换现有清单。
例如，将旧的 Ingress 定义转换为 networking.k8s.io/v1，可以运行：
kubectl convert -f ./legacy-ingress.yaml --output-version networking.k8s.io/v1
自动转换使用了一种类似于 Kubernetes 控制平面更新对象的技术，这些对象最初是使用旧 API 版本创建的。因为这是一个机械转换，你可能需要进去改变清单来调整默认值等等。

### 为升级进行排练

如果你管理集群的 API 服务器组件，那么你可以在升级到 Kubernetes v1.22 之前尝试删除这些 API。
为此，将以下内容添加到 kube-apiserver 命令行参数中：
\--runtime-config=admissionregistration.k8s.io/v1beta1=false,apiextensions.k8s.io/v1beta1=false,apiregistration.k8s.io/v1beta1=false,authentication.k8s.io/v1beta1=false,authorization.k9s.io/v1=false,certificates.k8s.io/v1beta=false,coordination.k8s.io/v1beta1=false,extensions/v1beta1/ingresses=false,networking.k8s.io/v1beta1=false
（作为一个副作用，这也关闭了 EndpointSlice 的 v1beta1——在测试时要注意。）
一旦你将集群中的所有 kube-apiserver 切换为使用该设置，这些 beta API 就会被删除。你可以测试 API 客户端（kubectl，部署工具，自定义控制器等）是否仍然按照你期望的方式工作，如果你需要，你可以恢复，而不必计划一个更具破坏性的降级。

### 给软件作者的建议

也许你读这篇文章是因为你是一个插件或其他集成 Kubernetes 组件的开发者？
如果你开发了 Ingress 控制器，webhook 验证器，API 聚合，或任何其他依赖于这些废弃 API 的工具，你应该已经开始转换你的软件了。
你可以使用预演中的升级技巧来运行你自己的只使用新 API 的 Kubernetes 集群，并确保你的代码工作正常。对于你的文档，请确保读者了解 Kubernetes v1.22 升级应该采取的任何步骤。
在可能的情况下，尽早帮助你的用户采用新的 API——也许是在测试环境中——这样他们就可以就任何问题向你提供反馈。
在 Kubernetes v1.25 中还会有更多的弃用之处，所以也计划将其囊括在内。

## Kubernetes API 的移除

以下是 Kubernetes 为什么删除一些 API 的背景信息，以及 Kubernetes 对稳定 API 的承诺。
Kubernetes 对其特性遵循一个已定义的**弃用策略\[6]**，包括 Kubernetes API。该策略允许替换 Kubernetes 的稳定（“GA”）API。重要的是，这个策略意味着只有当一个稳定的 API 有新的稳定版本可用时，才会弃用该 API。
这种稳定性保证很重要：如果你使用的是稳定的 Kubernetes API，就不会有新版本的发布迫使你切换到 alpha 或 beta 特性。
早期阶段是不同的。Alpha 特性还在测试中，可能不完整。alpha 特性在默认情况下几乎总是禁用的。Kubernetes 发布的版本可以并且确实删除了 alpha 版本中那些没有发挥作用的特性。
alpha 后面是 beta。这些特性通常是默认启用的；如果测试成功，该特性可以逐步趋于稳定。如果没有，可能需要重新设计。
去年，Kubernetes 正式对已经进入测试阶段的 API**采取\[7]**了一项政策：
对于 Kubernetes REST API 来说，当一个新特性的 API 达到测试版时，就开始倒计时了。beta 质量的 API 现在有三个版本：

- 到达 GA，弃用测试版，或者
- 拥有一个新的测试版（并弃用之前的测试版）。

_在写那篇文章的时候，Kubernetes 发布三次大约相当于 9 个日历月。同月晚些时候，Kubernetes 采用了新的发布节奏，即每日历年发布 3 个版本，所以现在的倒计时时间大约是 12 个日历月。_
不管 API 的删除是因为测试版特性已经趋于稳定，还是因为该 API 没有被证明是成功的，Kubernetes 将继续通过遵循其弃用策略并确保迁移选项被记录下来来删除 API。

### 展望未来

如果你使用 webhook 认证检查，有一个相关的设置。将来的 Kubernetes 版本将使用默认将 TokenReview 对象发送到 authentication.k8s.io/v1 API。目前，默认是发送 authentication.k8s.io/v1beta1 TokenReviews 到 webhooks，这仍然是 Kubernetes v1.22 版本的默认方式。不过，如果你想，现在可以切换到稳定的 API：在 kube-apiserver 的命令行选项中添加--authentication-token-webhook-version=v1，并检查用于身份验证的 webhook 是否仍按你期望的方式工作。
一旦你满意了，你就可以在控制平面上设置--authentication-token-webhook-version=v1 选项。
计划于明年发布的 v1.25 版本将停止提供几个 Kubernetes API 的 beta 版本，这些 API 目前已经稳定了一段时间。相同的 v1.25 版本将删除 PodSecurityPolicy，它已被弃用，不会升级到稳定状态。有关更多信息，请参阅**PodSecurityPolicy Deprecation: Past, Present, and Future\[8]**。
Kubernetes **1.25 计划移除的 API\[9]**的官方列表如下：

- The beta CronJob API (batch/v1beta1)
- The beta EndpointSlice API (networking.k8s.io/v1beta1)
- The beta PodDisruptionBudget API (policy/v1beta1)
- The beta PodSecurityPolicy API (policy/v1beta1)

## 想了解更多吗？

Kubernetes 发布说明中宣布了弃用的内容。你可以在**1.19\[10]**、**1.20\[11]**和**1.21\[12]**的发布说明中看到弃用声明。
有关弃用和移除过程的信息，请查看 Kubernetes 官方**弃用策略\[13]**文档。

### 参考资料

\[1]\[2]\[3]\[4]\[5]\[6]\[7]\[8]\[9]\[10]\[11]\[12]\[13]
Kubernetes API removals: *https://kubernetes.io/blog/2021/07/14/upcoming-changes-in-kubernetes-1-22/#kubernetes-api-removals*
Kubernetes 1.22 发布信息: *https://www.kubernetes.dev/resources/release/*
1.22 版本中删除的 API: _这些_
Ingress 先决条件: *https://kubernetes.io/docs/concepts/services-networking/ingress/#prerequisites*
Download Kubernetes: *https://kubernetes.io/releases/download/*
弃用策略: *https://kubernetes.io/docs/reference/using-api/deprecation-policy/*
采取: *https://kubernetes.io/blog/2020/08/21/moving-forward-from-beta/#avoiding-permanent-beta*
PodSecurityPolicy Deprecation: Past, Present, and Future: *https://kubernetes.io/blog/2021/04/06/podsecuritypolicy-deprecation-past-present-and-future/*
1.25 计划移除的 API: *https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-25*
1.19: *https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.19.md#deprecations*
1.20: *https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.20.md#deprecation*
1.21: *https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.21.md#deprecation*
弃用策略: *https://kubernetes.io/docs/reference/using-api/deprecation-policy/#deprecating-parts-of-the-api*

点击【阅读原文】阅读网站原文。
