---
title: "Object 管理"
weight: 1
---

# 概述

> 参考：
> 
> - [官方文档，概念-概述-使用 Kubernetes 对象-Kubernetes 对象管理](https://kubernetes.io/docs/concepts/overview/working-with-objects/object-management/)
> - [官方文档，任务-管理 Kubernetes 对象](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/)
> - [公众号-k8s 技术圈，理解 K8s 中的 Client-Side Apply 和 Server-Side Apply](https://mp.weixin.qq.com/s/EYtMO9KGRK_lHS2IW-mZug)
>   - [原文：掘金](https://juejin.cn/post/7173328614644006942)

使用 kubectl 等是传统的 **Client-Side Apply(简称 CSA)**，添加 --server-side 标志后，为 **Server-Side Apply(简称 SSA)。**

如果你经常与 kubectl 打交道，那相信你一定见过 kubectl.kubernetes.io/last-applied-configuration annotation，以及那神烦的 managedFields，像这样：

```yaml
$ kubectl get pods hello -oyaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","kind":"Pod","metadata":{"annotations":{},"creationTimestamp":null,"labels":{"run":"hello"},"name":"hello","namespace":"default"},"spec":{"containers":[{"image":"nginx","name":"hello","resources":{}}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always"},"status":{}}
  creationTimestamp: "2022-05-28T07:28:51Z"
  labels:
    run: hello
  managedFields:
  - apiVersion: v1
    fieldsType: FieldsV1
    fieldsV1:
      f:metadata:
        f:annotations:
          .: {}
          f:kubectl.kubernetes.io/last-applied-configuration: {}
        f:labels:
          .: {}
          f:run: {}
....
    manager: kubectl
    operation: Update
    time: "2022-05-28T07:28:51Z"
....
```

由这两个字段，引出本文的两位主角，Client-Side Apply（以下简称**CSA**）和 Server-Side Apply（以下简称**SSA**）

- kubectl.kubernetes.io/last-applied-configuration 是使用 kubectl apply 进行 Client-Side Apply 时，由 kubectl 自行填充的。
- managedFields 则是由 kubectl apply 的增强功能—— Server-Side Apply 的引入而添加。

## Client-Side Apply(客户端应用)

CSA 是 `kubectl apply` 早期(v1.14 版本之前)唯一的对象管理手段。

需要特别指出的是，`kubectl apply` 声明的仅仅是它关心的字段的状态，而不是整个对象的真实状态。apply 表达的意思是：“我”管理的字段应该和我 apply 的配置文件一致(但我不关心其他字段)。

什么是“我”管理的字段，什么又是其他的字段呢？举个例子，当我们希望使用 HPA 管理应用副本数时，[Kubernetes 推荐的做法](https://link.juejin.cn?target=https%3A%2F%2Fkubernetes.io%2Fdocs%2Ftasks%2Frun-application%2Fhorizontal-pod-autoscale%2F%23migrating-deployments-and-statefulsets-to-horizontal-autoscaling)是在 apply 的配置文中不指定具体 replicas 副本数。首次部署时，K8S 会将 replicas 值设置为默认 1，随后由 HPA 控制器扩容到合适的副本数。

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: nginx
  name: nginx
spec:
  # replicas: 1 不要设置replicas
  selector:
    matchLabels:
      app: nginx
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: nginx
    spec:
      containers:
        - image: nginx:latest
          name: nginx
          resources: {}
```

当升级应用时（修改镜像版本），修改配置文件中的 image 字段，再次执行 kubectl apply。此时 kubectl apply 只会影响镜像版本(因为他是“我”管理的字段)，而不会影响 HPA 控制器设置的副本数。在这个例子中，replicas 字段不是 kubectl apply 管理的字段，因此更新镜像时不会被删除，避免了每次应用升级时，副本数都会被重置。

在上述例子中，为了能识别出 replicas 不是 kubectl 管理的字段，kubectl 需要一个标识，用来追踪对象中哪些字段是由 kubectl apply 管理的，而这个标识就是 last-applied-configuration。 该 annotation 是在 kubectl apply 时，由 kubectl 客户端自行填充——每次执行 kubectl apply 时（未启用 SSA），kubectl 会将本次 apply 的配置文件全量的记录在 last-applied-configurationannotation 中，用于追踪哪些字段由 kubectl apply 管理。

CSA 的工作工作机制大致如下：当 apply 一个对象，如果该对象不存在，则创建它（同时写入 last-applied-configuration）。如果对象已经存在，则 kubectl 需要根据以下三个状态：

- 当前配置文件所表示的对象在集群中的真实状态。（修改对象前先 Get 一次）
- 当前 apply 的配置。
- 以及上次 apply 的配置。 （在 last-applied-configuration 里）

计算出 patch 报文，通过 patch 方式进行更新（而不是将配置文件全量的发送到服务端）。 patch 报文的计算方法如下：

- 计算需要被删除的字段。如果字段存在在 last-applied-configuration 中，但配置文件中没有，将删除它们。
- 计算需要修改或添加的字段。如果配置文件中的字段与真实状态不一致，则添加或修改它们。
- 特别的，对于那些 last-applied-configuration 中不存在的字段，不要修改它们（例如上述示例中的 replicas 字段）

详细的 patch 计算示例可参考 [K8S 文档中给出的详细示例](https://link.juejin.cn?target=https%3A%2F%2Fkubernetes.io%2Fzh-cn%2Fdocs%2Ftasks%2Fmanage-kubernetes-objects%2Fdeclarative-config%2F%23apply-%25E6%2593%258D%25E4%25BD%259C%25E6%2598%25AF%25E5%25A6%2582%25E4%25BD%2595%25E8%25AE%25A1%25E7%25AE%2597%25E9%2585%258D%25E7%25BD%25AE%25E5%25B7%25AE%25E5%25BC%2582%25E5%25B9%25B6%25E5%2590%2588%25E5%25B9%25B6%25E5%258F%2598%25E6%259B%25B4%25E7%259A%2584)。

由此可见，last-applied-configuration 体现的是一种 ownership 的关系，表示哪些字段是由 kubectl 管理，它是 kubectl apply 时，计算 patch 报文的依据。

## Server-Side Apply(服务端应用)

> 参考：
> - [官方文档，参考-API 概述-服务端 Apply](https://kubernetes.io/zh-cn/docs/reference/using-api/server-side-apply/)
> - <https://cloud.tencent.com/developer/article/1610073>
> - <https://yanhang.me/post/2021-ssa/>

**SSA 是另一种声明式的对象管理方式，和 CSA 的作用是基本一致的**。SSA 始于从 1.14 开始发布 alpha 版本，到 1.16 beta，到 1.18 beta2，终于在 1.22 升级为 GA。

从 Kubernetes 1.18 开始可以看到一个明显的变化就是资源的 YAML 在 metadata 字段多了很多一个 `managedFields` 字段，该字段用来声明一个资源的各个字段的具体的管理者是谁。

### 字段管理

> 参考：
> - [官方文档，参考-API 概述-服务端 Apply-字段管理](https://kubernetes.io/docs/reference/using-api/server-side-apply/#field-management)

# 管理对象的方式

## 使用配置文件对 Kubernetes 对象进行声明式管理

https://kubernetes.io/zh-cn/docs/tasks/manage-kubernetes-objects/declarative-config/

## 使用 Kustomize 对 Kubernetes 对象进行声明式管理

https://kubernetes.io/zh-cn/docs/tasks/manage-kubernetes-objects/kustomization/

## 使用指令式命令管理 Kubernetes 对象

https://kubernetes.io/zh-cn/docs/tasks/manage-kubernetes-objects/imperative-command/

## 使用配置文件对 Kubernetes 对象进行命令式管理

https://kubernetes.io/zh-cn/docs/tasks/manage-kubernetes-objects/imperative-config/

## 使用 kubectl patch 更新 API 对象

https://kubernetes.io/zh-cn/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/

使用 kubectl patch 更新 Kubernetes API 对象。做一个策略性的合并 patch 或 JSON 合并 patch。
