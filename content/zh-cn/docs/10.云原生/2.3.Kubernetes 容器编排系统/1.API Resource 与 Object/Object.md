---
title: "Object"
weight: 2
---

# 概述

> 参考：
> 
> - [官方文档,概念-使用 Kubernetes 对象](https://kubernetes.io/docs/concepts/overview/working-with-objects/)

从某些角度看来，Kubernetes 里的一切介 **Object(对象)**。就像 Linux 里，一切介文件的角度一样。

**API Resource(资源)** 用于表示 **Objectk(对象)** 的集合。例如 pod 资源可以用于描述所有 pod 资源类型的对象，比如我创建一个 pod 资源，生成了一个名为 test1 的 pod 类型的对象，如果创建了多个 pod 对象，那么每个对象都是 pod 类型的资源。

Kubernetes **Objectk(对象)** 是 Kubernetes 系统中，**Resource(资源)** 的持久化实体。Kubernetes 使用这些实体来表示集群的状态。具体来说，他们可以描述：

- 哪些容器化应用程序正在运行（以及在哪些节点上）
- 这些应用可用的资源
- 有关这些应用程序的行为的策略，例如重新启动策略，升级和容错

Kubernetes 对象是“record of intent(目标性记录)”：即,一旦创建了对象，Kubernetes 系统会确保对象存在。通过创建对象，本质上时告诉 Kubernetes 系统你希望集群的工作负载是什么样的，这就是 kubernetes 集群的 desired state(期望状态)。

要操作 Kubernetes 对象(无论是创建，修改还是删除)，都需要使用 [Kubernetes API](https://kubernetes.io/docs/concepts/overview/kubernetes-api/) 。例如，当使用 kubectl 命令管理工具时，CLI 会执行必要的 KubernetesAPI 调用。也可以直接在自己的程序中使用 [Client Libraries](https://kubernetes.io/docs/reference/using-api/client-libraries/) 来调用 KubernetesAPI。[Client Libraries](https://kubernetes.io/docs/reference/using-api/client-libraries/)(客户端库)可以理解为编程语言的一个第三方库，通过这个库中的方法，可以直接调用 KubernetesAPI。

用白说描述：每个已经启动的 pod 就是一个 object(对象)，每个已经创建的 namesapce 也是一个 object。而 pod、namespace 本身称为 resource(资源)。所以 object 就叫 kubernetes 系统中持久化的实体。

# 描述 kubernetes 对象

在 Kubernetes 中创建对象时，必须提供描述其所需状态的 **spec(规范)** 以及有关该对象的 **基本信息**(例如名称)。当我们使用 KubernetesAPI 创建对象（直接向 API 发请求 或者 通过 kubectl 向 API 发请求）时，这个 API 请求必须在 body 中包含 JSON 格式的信息。通常，我们会在 .yaml 文件中将信息提供给 kubectl， kubectl 发出 API 请求时将信息转换为 JSON 格式。一个 YAML 文件中定义了一个我要如何运行一个对象。这个 yaml 文件其实就相当于该对象的配置文件。

**通常，我们把描述一个对象的 .yaml 或 .json 格式的文件，称为** [**Manifest(清单)**](https://kubernetes.io/docs/reference/glossary/?all=true#term-manifest)。Manifest 指定了在应用该 Manifest 时，Kubrenetes 将维护的对象的期望状态。

> 比如，交流时经常这么说 manifest of resource(资源的清单)，就是指创建该资源对象时的 yaml 文件

注意：

- kubernetes 把 YAML 语言里的 mapping 称为 field(字段) 下文都用 field 来进行描述，每个 field 与 yaml 语言都是，都是一组 key/value pair(键值对)。一般情况下，field 的 key 是一个 k8s 对象的属性，一般是不变的；field 的 value 是这个属性的具体描述内容。
- 一个 Manifest 文件最大为 262144 Bytes，即 256 KiB，如果文件超过了 256 KiB，会报错提示，无法创建对象，效果如下：
  - ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/qbvmvb/1620550376332-3ee2445f-75f9-437a-9f1b-c175fd6c24d1.png)

这是一个简单的 .yaml 文件示例，显示 kubernetes 中一个 deployment 类型的对象所需的 field 和对象规范

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
            - containerPort: 80
```

想要创建一个 k8s 对象，该对象的 yaml 中必须包含这么几个 Field(字段)

1. apiVersion: GROUP/VERSION # 您正在使用哪个版本的 Kubernetes API 创建该对象（组/API 版本，可以通过 kubectl api-version 查询所支持的组/版本）
2. kind: ResourcesType # 您要创建什么类型的对象(指明要创建的资源类型,可以通过 kubectl api-resources 查询所有支持的资源类型)
3. metadata: # 有助于唯一标识对象的数据，包括 name 字符串 UID，和可选 namespace（指明该资源的元数据，其中 name 是必须要写明的元数据项）
4. sepc: # 对象的期望状态（指明该资源的规格(specification)）
5. status: # (特殊的字段) 指明当前的状态，本字段由 kubernetes 集群自己维护，用户无法自己定义，status 用于确保现有状态无限接近于目标状态(即 sepc 中的信息)

其中 metadata 与 spec 是对一个对象描述的最重要的内容。

使用命令 kubectl explain pods.spec.containers。。。。可以查看每个 field 的意义以及其下还有什么 field 可用

## 对象的 Spec(规格) 与 Status(状态)

几乎每个 Kubernetes 对象包含两个嵌套的对象字段，它们负责管理对象的配置： 对象 spec（规约） 和 对象 status（状态） 。 对于具有 spec 的对象，你必须在创建对象时设置其内容，描述你希望对象所具有的特征： 期望状态（Desired State） 。

status 描述了对象的 当前状态（Current State），它是由 Kubernetes 系统和组件 设置并更新的。在任何时刻，Kubernetes 控制平面 都一直积极地管理着对象的实际状态，以使之与期望状态相匹配。

例如，Kubernetes 中的 Deployment 对象能够表示运行在集群中的应用。 当创建 Deployment 时，可能需要设置 Deployment 的 spec，以指定该应用需要有 3 个副本运行。 Kubernetes 系统读取 Deployment 规约，并启动我们所期望的应用的 3 个实例 —— 更新状态以与规约相匹配。 如果这些实例中有的失败了（一种状态变更），Kubernetes 系统通过执行修正操作 来响应规约和状态间的不一致 —— 在这里意味着它会启动一个新的实例来替换。

关于对象 spec、status 和 metadata 的更多信息，可参阅 [Kubernetes API 约定](https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md)

## Object 在现阶段有三大类别

1. 对象(object)：如上文所述，代表系统中的一个永久资源(实体)，例如 pod、service、namespace、node 等等。通过操作这些资源的属性，客户端可以对该对象进行创建、修改、删除、获取操作
2. 列表(list)：一个或多个资源类别的集合。例如 PodList、ServiceList、NodeList 等。kubectl get pod 命令实际上就是在获取 pod 资源的 list
3. 简单类别(simple)：该类别包含作用在对象上的特殊行为和非持久实体。 该类别限制了使用范围， 它有一个通用元数据的有限集合， 例如 Binding、 Status。
