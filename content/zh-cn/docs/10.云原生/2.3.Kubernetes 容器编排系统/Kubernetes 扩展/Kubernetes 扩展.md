---
title: Kubernetes 扩展
---

# 扩展 Kubernetes 集群概述

> 参考：
> - 概念：<https://kubernetes.io/docs/concepts/extend-kubernetes/extend-cluster/>
> - 任务：<https://kubernetes.io/docs/tasks/extend-kubernetes/>

这里引用张磊大佬的一张图来开篇![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eo9qpz/1619144340194-3bfc699a-f6ca-4732-b4f2-0d49c7512b1a.jpeg)

Kubernetes 是高度可配置和可扩展的。因此，极少需要分发或提交补丁代码给 Kubernetes 项目。通过对 Kubernetes 的扩展，可以将中间狭窄的部分扩大。

本文档介绍自定义 Kubernetes 集群的方式。本文档的目标读者包括希望了解如何使 Kubernetes 集群满足其业务环境需求的 [集群运维人员](https://kubernetes.io/zh/docs/reference/glossary/?all=true#term-cluster-operator)、 Kubernetes 项目的[贡献者](https://kubernetes.io/zh/docs/reference/glossary/?all=true#term-contributor)。 或潜在的[平台开发人员](https://kubernetes.io/zh/docs/reference/glossary/?all=true#term-platform-developer) 也可以从本文找到有用的信息，如对已存在扩展点和模式的介绍，以及它们的权衡和限制。

自定义方法可以大致分为两类

1. Configuration (配置) # 配置只涉及更改标志参数、本地配置文件或 API 资源；
2. Extension (扩展) # 扩展涉及运行额外的程序或服务。

扩展的实现依赖于某些特殊配置，这些配置可能在某些 k8s 发行版中并不自带。需要自行修改程序参数来让集群支撑某些扩展方式。

## 配置 Kubernetes 集群

关于 配置文件 和 标志 的说明文档位于在线文档的"参考"部分，按照可执行文件组织：

- [kubelet](https://kubernetes.io/zh/docs/reference/command-line-tools-reference/kubelet/)
- [kube-apiserver](https://kubernetes.io/zh/docs/reference/command-line-tools-reference/kube-apiserver/)
- [kube-controller-manager](https://kubernetes.io/zh/docs/reference/command-line-tools-reference/kube-controller-manager/)
- [kube-scheduler](https://kubernetes.io/zh/docs/reference/command-line-tools-reference/kube-scheduler/).

在托管的 Kubernetes 服务或受控安装的 Kubernetes 版本中，标志和配置文件可能并不总是可以更改的。而且当它们可以进行更改时，它们通常只能由集群管理员进行更改。此外，标志和配置文件在未来的 Kubernetes 版本中可能会发生变化，并且更改设置后它们可能需要重新启动进程。出于这些原因，只有在没有其他选择的情况下才使用它们。

内置策略 API ，例如 [ResourceQuota](https://kubernetes.io/zh/docs/concepts/policy/resource-quotas/)、 [PodSecurityPolicy](https://kubernetes.io/zh/docs/concepts/policy/pod-security-policy/)、 [NetworkPolicy](https://kubernetes.io/zh/docs/concepts/services-networking/network-policies/) 和基于角色的权限控制 ([RBAC](https://kubernetes.io/zh/docs/reference/access-authn-authz/rbac/))， 是内置的 Kubernetes API。API 通常与托管的 Kubernetes 服务和受控的 Kubernetes 安装一起使用。 它们是声明性的，并使用与其他 Kubernetes 资源（如 Pod ）相同的约定，所以新的集群配置可以重复使用， 并以与应用程序相同的方式进行管理。 而且，当它们变稳定后，也遵循和其他 Kubernetes API 一样的 [支持政策](https://kubernetes.io/docs/reference/using-api/deprecation-policy/)。 出于这些原因，在合适的情况下它们优先于 配置文件 和 标志 被使用。

## 扩展 Kubernetes 集群

Kubernetes 的设计是通过编写客户端程序来实现自动化的。 任何读和（或）写 Kubernetes API 的程序都可以提供有用的自动化工作。 自动化 程序可以运行在集群之中或之外。按照本文档的指导，你可以编写出高可用的和健壮的自动化程序。 自动化程序通常适用于任何 Kubernetes 集群，包括托管集群和受管理安装的集群。

控制器（Controller） 模式是编写适合 Kubernetes 的客户端程序的一种特定模式。 控制器通常读取一个对象的 .spec 字段，可能做出一些处理，然后更新对象的 .status 字段。

一个控制器是 Kubernetes 的一个客户端。 当 Kubernetes 作为客户端调用远程服务时，它被称为 Webhook ， 远程服务称为 Webhook 后端。 和控制器类似，Webhooks 增加了一个失败点。

在 webhook 模型里，Kubernetes 向远程服务发送一个网络请求。 在 可执行文件插件 模型里，Kubernetes 执行一个可执行文件（程序）。 可执行文件插件被 kubelet（如 [Flex 卷插件](https://github.com/kubernetes/community/blob/master/contributors/devel/flexvolume.md)和 [网络插件](https://kubernetes.io/zh/docs/concepts/extend-kubernetes/compute-storage-net/network-plugins/)和 kubectl 所使用。

# 扩展 Kubernetes 集群的方式：

## Kubernetes API 扩展

1. [**Custom Resources Definitions(自定义资源定义)**](https://www.teambition.com/project/5f90e312755d8a00446050eb/app/5eba5fba6a92214d420a3219/workspaces/5f90e312c800160016ea22fb/docs/5f9a605812d5ba00014a7368)
2. [**Aggregation API(聚合 API)**](https://www.teambition.com/project/5f90e312755d8a00446050eb/app/5eba5fba6a92214d420a3219/workspaces/5f90e312c800160016ea22fb/docs/5f9a602837398300016bc061)
3. [**Operator 模式**](https://www.teambition.com/project/5f90e312755d8a00446050eb/app/5eba5fba6a92214d420a3219/workspaces/5f90e312c800160016ea22fb/docs/5f9a5fe537398300016bbb9a)

## 计算、存储、网络扩展。

## Service Catalog 服务目录。详见：Service Catalog 服务目录

# 扩展点

扩展点的意思就是在 kubernetes 可以实现扩展的位置。“点” 这个字的意思类似与这种语境(我下面要说的几点注意事项。某件事物的触发点。兴奋点。G 点哈哈囧)

下图显示了各种扩展点如何与 Kubernetes 控制平面进行交互。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eo9qpz/1619144357206-5faf80b1-f5f1-4daa-bc39-26ab29976d5e.jpeg)

下图显示了 Kubernetes 系统的扩展点。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eo9qpz/1619144364591-352e19a9-91b8-4a9d-a257-13dc02b6008e.jpeg)

1. 用户通常使用 kubectl 与 Kubernetes API 进行交互。 [kubectl 插件](https://kubernetes.io/zh/docs/tasks/extend-kubectl/kubectl-plugins/)扩展了 kubectl 可执行文件。 它们只影响个人用户的本地环境，因此不能执行站点范围的策略。
2. API 服务器处理所有请求。API 服务器中的几种类型的扩展点允许对请求进行身份认证或根据其内容对其进行阻止、 编辑内容以及处理删除操作。这些内容在 [API 访问扩展](https://kubernetes.io/zh/docs/concepts/extend-kubernetes/#api-access-extensions)小节中描述。
3. API 服务器提供各种 资源（Resource） 。 内置的资源种类（Resource Kinds） ，如 pods， 由 Kubernetes 项目定义，不能更改。你还可以添加你自己定义的资源或其他项目已定义的资源， 称为 自定义资源（Custom Resource），如[自定义资源](https://kubernetes.io/zh/docs/concepts/extend-kubernetes/#user-defined-types) 部分所述。自定义资源通常与 API 访问扩展一起使用。
4. Kubernetes 调度器决定将 Pod 放置到哪个节点。有几种方法可以扩展调度器。 这些内容在[调度器扩展](https://kubernetes.io/zh/docs/concepts/extend-kubernetes/#scheduler-extensions) 小节中描述。
5. Kubernetes 的大部分行为都是由称为控制器（Controllers）的程序实现的，这些程序是 API 服务器的客户端。 控制器通常与自定义资源一起使用。
6. kubelet 在主机上运行，并帮助 Pod 看起来就像在集群网络上拥有自己的 IP 的虚拟服务器。 [网络插件](https://kubernetes.io/zh/docs/concepts/extend-kubernetes/#network-plugins/)让你可以实现不同的 pod 网络。
7. kubelet 也负责为容器挂载和卸载卷。新的存储类型可以通过 [存储插件](https://kubernetes.io/docs/concepts/extend-kubernetes/#storage-plugins)支持。

如果你不确定从哪里开始扩展，下面流程图可以提供一些帮助。请注意，某些解决方案可能涉及多种类型的扩展。
