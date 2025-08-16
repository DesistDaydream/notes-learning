---
title: Custom Resource Definitions(CRD)
linkTitle: Custom Resource Definitions(CRD)
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，概念 - 扩展 Kubernetes - 扩展 API - 自定义资源，CRD](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/#customresourcedefinitions)

**Custom Resource Definitions 自定义资源定义，简称 CRD**。是一个 Kubernetes 的 [API Resource](/docs/10.云原生/Kubernetes/API%20Resource%20与%20Object/API%20Resource.md)(API 资源)。其前身是 Kubernetes1.6 版本中一个叫做 ThirdPartyResource(第三方资源，简称 TPR) 的内建对象，可以用它来创建自定义资源，但该对象在 kubernetes1.7 中版本已被 CRD（CustomResourceDefinition）取代。**CRD 的目的是让 Kubernetes API 能够认识新对象(就是 yaml 中的 kind)**。所以通过 CRD 创建的对象可以跟 kubernetes 中内建对象一样使用 kubectl 操作，就像 kubectl 操作 pod 一样，如果我定义了一个名为 lch 的 crd ，那么我可以使用 kubectl get lch 命令来操作 lch 这个资源

注意：单纯设置了 CRD，并没有什么用，只有跟自定义控制器(controller)结合起来，才能将资源对象中的声明式 API 翻译成用户所期望的状态。自定义控制器可以用来管理任何资源类型，但是一般是跟 CRD 结合使用。自定义控制器称为 [Operator](/docs/10.云原生/Kubernetes/Kubernetes%20扩展/Operator%20模式.md)。

为什么这么说呢？

比如，在一个 pod 的 yaml 文件里每一个字段就是对该资源的定义，name 字段定义了该资源名字，image 字段定义了 pod 这个资源的 container 所使用的镜像等等。但是，既然有字段可以定义一个 k8s 资源，那么谁又来决定这个字段就是其所描述的功能呢？~答案当然是 kubernetes 主程序(也就是 kube-controller)。

所以上文“注意”里描述的也是同样的道理，自己定义了一个资源的的各个字段，那么谁又能实现这个自定义资源各个字段所描述的功能呢(既然叫自定义，那么自定义的东西一般主程序是不认识的)？这就需要自定义控制器了(作为一个翻译器，把自定义的资源翻译成 k8s 主程序可以理解的东西)，所以 CRD 单独使用别没有什么作用，你想啊，咱自己定义了一个资源的一个字段，但是 kubernetes 没法识别这个字段是干什么用的，所以才需要一个第三方的工具，把这个字段的作用翻译给 kubernetes。一般情况下，自定义控制器是开发人员开发的一款程序，可以通过 Deployment 的方式部署在 k8s 集群上。与此同时，自定义控制器只需要持续监听关联的的 CRD ，让其保持声明的状态一致即可。(自定义控制器 与 controlller-manager 监控 deployment、job 等资源的行为是类似的)

# CRD 的使用方式

官方文档：<https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/>

Kubernetes 中的 CustomResourceDefinition API 资源允许你定义定制资源。 定义 CRD 对象的操作会使用你所设定的名字和模式定义（Schema）创建一个新的定制资源， Kubernetes API 负责为你的定制资源提供存储和访问服务。 CRD 对象的名称必须是合法的 [DNS 子域名](https://kubernetes.io/zh/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)。

CRD 使得你不必编写自己的 API 服务器来处理定制资源，不过其背后实现的通用性也意味着 你所获得的灵活性要比 [API 服务器聚合](https://kubernetes.io/zh/docs/concepts/extend-kubernetes/api-extension/custom-resources/#api-server-aggregation)少很多。

关于如何注册新的定制资源、使用新资源类别的实例以及如何使用控制器来处理事件， 相关的例子可参见[定制控制器示例](https://github.com/kubernetes/sample-controller)。

参考下面的 CRD，resourcedefinition.yaml：

```yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  # 名称必须符合下面的格式：<plural>.<group>
  name: crontabs.stable.example.com
spec:
  # REST API使用的组名称：/apis/<group>/<version>
  group: stable.example.com
  # REST API使用的版本号：/apis/<group>/<version>
  version: v1
  # Namespaced或Cluster
  scope: Namespaced
  names:
    # URL中使用的复数名称: /apis/<group>/<version>/<plural>
    plural: crontabs
    # CLI中使用的单数名称
    singular: crontab
    # CamelCased格式的单数类型。在清单文件中使用
    kind: CronTab
    # CLI中使用的资源简称
    shortNames:
    - ct
```

创建该 CRD：

kubectl create -f resourcedefinition.yaml

访问 RESTful API 端点如 <http://172.20.0.113:8080> 将看到如下 API 端点已创建：

/apis/stable.example.com/v1/namespaces/\*/crontabs/...

**创建自定义对象**

如下所示：

```yaml
apiVersion: "stable.example.com/v1"
kind: CronTab
metadata:
  name: my-new-cron-object
spec:
  cronSpec: "* * * * /5"
  image: my-awesome-cron-image
```

引用该自定义资源的 API 创建对象。

**终止器**

可以为自定义对象添加一个终止器，如下所示：

```yaml
apiVersion: "stable.example.com/v1"
kind: CronTab
metadata:
  finalizers:
  - finalizer.stable.example.com
```

删除自定义对象前，异步执行的钩子。对于具有终止器的一个对象，删除请求仅仅是为 metadata.deletionTimestamp 字段设置一个值，而不是删除它，这将触发监控该对象的控制器执行他们所能处理的任意终止器。

详情参考：[Extend the Kubernetes API with CustomResourceDefinitions](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/)
