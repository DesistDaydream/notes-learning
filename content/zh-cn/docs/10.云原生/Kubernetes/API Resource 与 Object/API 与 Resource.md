---
title: API 与 Resource
linkTitle: API 与 Resource
date: 2019-11-13T21:49:00
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，概念 - 概述 - Kubernetes API](https://kubernetes.io/docs/concepts/overview/kubernetes-api/)
> - [官方文档，参考 - API 概述](https://kubernetes.io/docs/reference/using-api/)

Kubernetes API 使我们可以查询和操纵 Kubernetes API 中资源的状态。Kubernetes API 符合 RESTful 规范。

Kubernetes 把自身一切抽象理解为 **Resource(资源)**，也叫 **API Resource**(有地方也叫 Group Resource)。对集群的所有操作都是通过对 Kubernetes API 的 HTTP(s) 请求来实现的。可以使用命令 `kubectl api-resources` 命令查看所有支持的资源。

kubernetes 控制平面的核心是 **API Server**。API Server 是实现了 Kubernets API 的应用程序，并为 Kubernetes 公开了一个 HTTP(s) 的 API，以供用户、集群中的不同部分和集群外部组件相互通信。

Kubernetes 中各种资源(对象)的数据都通过 API 接口被提交到后端的持久化存储（etcd）中，Kubernetes 集群中的各部件之间通过该 API 接口实现解耦合，同时 Kubernetes 集群中一个重要且便捷的管理工具 kubectl 也是通过访问该 API 接口实现其强大的管理功能的。

> Note：kubectl 就是代替用户执行各种 http 请求的工具

在 Kubernetes 系统中，在大多数情况下，API 定义和实现都符合标准的 HTTP REST 格式，比如通过标准的 HTTP 动词（POST、PUT、GET、DELETE）来完成对相关资源对象的查询、创建、修改、删除等操作。但同时，Kubernetes 也为某些非标准的 REST 行为实现了附加的 API 接口，例如 Watch 某个资源的变化、进入容器执行某个操作等。另外，某些 API 接口可能违背严格的 REST 模式，因为接口返回的不是单一的 JSON 对象，而是其他类型的数据，比如 JSON 对象流或非结构化的文本日志数据等。

另外，从另一个角度看，其实 kubernetes 就是提供了一个 web 服务，只是这个 web 服务不像传统的 B/S 架构那样，可以通过浏览器直接操作~kubernetes API 就是这个 web 服务的入口。

> 注意：Kubernetes 的 API 与传统意义上的 API 不太一样。传统 API，一个 API 就是一个功能；而 Kubernetes API 中，一个 API 实际上又可以当作功能，也可以当作一个资源。对 API 的操作，就是对 Kubernets 资源进行操作

## API Resource(资源) 分类

> 参考：
>
> - [官方文档，参考 - kubernetes API](https://kubernetes.io/docs/reference/kubernetes-api/)
> - [1.19 版本 API 参考(一页模式)](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/)(要查看其他版本，修改 URL 最后的版本号即可)。

资源大体可以分为下面几类：

- **workload(工作负载)** # 用于在集群上管理和运行容器
   - Pod，Deployment，StatefuSet，DaemonSet，Job 等
- **Discovery & LB(服务发现及均衡)** # 可以使用这些资源类型的对象将工作负载“缝合”到一个外部可访问的、负载均衡的服务中。
   - Service，Ingress 等
- **Config & Storage(配置与存储)** # 这种类型的资源是用于将初始化数据注入到应用程序中并保留容器外部数据的对象。
   - Volume，ConfigMap，secret 等
- **Cluster(集群级资源)** # 这种类型的资源对象定义了群集本身的配置方式。这些通常仅由集群运营商使用。
   - Namesapces,Node,Role,ClusterRole,RoleBinding,ClusterRoleBinding 等
- **Metadata(元数据型资源)** # 这种类型的资源是用于配置集群中其他资源行为的对象。
   - HPA，PodTemplate，LimitRange 等

各种资源所用的 manifest 文件中的各个字段的含义就可以参考该页面找到详解。

## API Resource(资源) 的 URL 结构

在 Kubernetes 中，资源的 URL 结构是由：Group（组）、Version（版本）和 Resource（资源种类）三个部分组成的。(还有一种 /metrics，/healthz 之类的结构，这里面的资源是系统自带的，不在任何组里)

通过这样的结构，整个 Kubernetes 里的所有资源，实际上就可以用如下图的树形结构表示出来：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/sz9hgm/1616120310758-dc53a2df-2a39-45e9-92e3-9beb5d9101f0.png)

比如，如果要创建一个 CronJob 资源

```yaml
apiVersion: batch/v2alpha1
kind: CronJob
```

在这个 YAML 文件中，“CronJob”就是资源的种类(Resource)，“batch”就是它的组(Group)，v2alpha1 就是它的版本(Version)。

现阶段，有两个 API Groups 正在使用

- **core group(核心组)** # 在/api/v1 路径下(由于某些历史原因而并没有在 `/apis/core/v1` 路径下)。核心组是不需要 Group 的（即：它们 Group 是 `""`）。URI 路径为 `/api/v1`，并且在定义资源的 manifest 文件中 apiVersion 字段的值不用包含组名，直接使用 v1 即可
- **named groups(已命名组)** # URI 路径为 `/apis/$GROUP_NAME/$VERSION`，在定义资源的 manifest 文件中 apiVersion 中省略 apis，使用 GroupName/Version

Notes:

- 有的资源是 cluster 级别的(比如 node)，有的资源是 namespace 级别的(比如 pod)，对于 namespace 级别的资源，可以在 Version 和 Resource 中间添加 namespace 字段以获取指定 namespace 下的资源。i.e.`/api/v1/namespaces/$NAMESPACE/pods/` (`${NAMESPACE}` 就是具体的 namesapce 的名称)。
- 所以 namesapce 级别资源的对象的 URI 应该像这样：`/api/v1/namespaces/kube-system/pods/coredns-5644d7b6d9-tw4rh`
- 而 cluster 级别资源的对象的 URI 则是：`/api/v1/nodes/master1`

**所有资源类型要么受集群范围限制（`/apis/GROUP/VERSION/_`），要么受命名空间限制（`/apis/GROUP/VERSION/namespaces/NAMESPACE/_`）**

集群范围的资源：

- GET /apis/GROUP/VERSION/RESOURCETYPE # 返回指定资源类型的资源集合(返回的是一个 list 列表，比如 NodeList 等)
- GET /apis/GROUP/VERSION/RESOURCETYPE/NAME # 返回指定资源类型下名为 NAME 的的资源

名称空间范围的资源：

- GET /apis/GROUP/VERSION/RESOURCETYPE # 返回所有名称空间指定资源类型的实例集合(返回的是一个 list 列表，比如 podList、serviceList 等)
- GET /apis/GROUP/VERSION/namespaces/NAMESPACE/RESOURCETYPE # 返回 NAMESPACE 下指定 ResourceType 的所有实例集合(返回的是一个 list 列表，比如 podList、serviceList 等)
- GET /apis/GROUP/VERSION/namespaces/NAMESPACE/RESOURCETYPE/NAME # 返回 NAMESPACE 下指定 ResourceType，名为 NAME 的实例

# Declarative API(声明式 API) 的特点：

- 首先，所谓 **Declarative(声明式)**，指的就是我只需要提交一个定义好的 API 对象来 **Declarative(声明)** 我所期望的状态是什么样子。
- 其次，“声明式 API”允许有多个 API 写端，以 PATCH 的方式对 API 对象进行修改，而无需关心本地原始 YAML 文件的内容。
- 最后，也是最重要的，有了上述两个能力，Kubernetes 项目才可以基于对 API 对象的增、删、改、查，在完全无需外界干预的情况下，完成对“实际状态”和“期望状态”的调谐（Reconcile）过程。

所以说，声明式 API，才是 Kubernetes 项目编排能力“赖以生存”的核心所在。而想要实现 声明式 API，离不开 Controller 控制器，K8S 的大脑 的工作。

# API URL 使用示例

下面是在 1.18.8 版本下获取到的 api 路径结构

根路径将列出所有可用路径

```json
~]# curl --cacert /etc/kubernetes/pki/ca.crt -H "Authorization: Bearer ${TOKEN}"  https://172.38.40.215:6443/ -s
{
  "paths": [
    "/api",
    "/api/v1",
    "/apis",
    "/apis/",
    "/apis/admissionregistration.k8s.io",
    "/apis/admissionregistration.k8s.io/v1",
    "/apis/admissionregistration.k8s.io/v1beta1",
    "/apis/apiextensions.k8s.io",
    "/apis/apiextensions.k8s.io/v1",
    "/apis/apiextensions.k8s.io/v1beta1",
    "/apis/apiregistration.k8s.io",
......
}
```

如果访问到错误的资源，还会返回 404 的响应码

```json
~]# curl -s --cacert /etc/kubernetes/pki/ca.crt -H "Authorization: Bearer ${TOKEN}"  https://172.38.40.215:6443/api/v1/service
{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {

  },
  "status": "Failure",
  "message": "the server could not find the requested resource",
  "reason": "NotFound",
  "details": {

  },
  "code": 404
}
```

在“组/版本”下面可以看到该“组/版本”下所包含的 API 资源列表

```json
~]# curl -s --cacert /etc/kubernetes/pki/ca.crt -H "Authorization: Bearer ${TOKEN}"  https://172.38.40.215:6443/api/v1/
{
  "kind": "APIResourceList",
  "groupVersion": "v1",
  "resources": [
.......
    {
      "name": "configmaps",
      "singularName": "",
      "namespaced": true,
      "kind": "ConfigMap",
      "verbs": [
        "create",
        "delete",
        "deletecollection",
        "get",
        "list",
        "patch",
        "update",
        "watch"
      ],
      "shortNames": [
        "cm"
      ],
      "storageVersionHash": "qFsyl6wFWjQ="
    },
    {
      "name": "endpoints",
      "singularName": "",
      "namespaced": true,
      "kind": "Endpoints",
      "verbs": [
        "create",
        "delete",
        "deletecollection",
        "get",
        "list",
        "patch",
        "update",
        "watch"
      ],
      "shortNames": [
        "ep"
      ],
      "storageVersionHash": "fWeeMqaN/OA="
    },
......
}
```

在“资源”下可以看到该“资源”下所包含的所有对象，下图是 pod 资源的列表，包含所有 pod 对象及其信息

```json
~]# curl -s --cacert /etc/kubernetes/pki/ca.crt -H "Authorization: Bearer ${TOKEN}"  https://172.38.40.215:6443/api/v1/pods | more
{
  "kind": "PodList",
  "apiVersion": "v1",
  "metadata": {
    "selfLink": "/api/v1/pods",
    "resourceVersion": "618871"
  },
  "items": [
    {
      "metadata": {
        "name": "cattle-cluster-agent-cc6ddc6dc-7f89l",
        "generateName": "cattle-cluster-agent-cc6ddc6dc-",
        "namespace": "cattle-system",
        "selfLink": "/api/v1/namespaces/cattle-system/pods/cattle-cluster-agent-cc6ddc6dc-7f89l",
        "uid": "72f4a825-feb2-416a-900d-d8401acc9a18",
        "resourceVersion": "452264",
        "creationTimestamp": "2020-09-13T09:59:49Z",
        "labels": {
          "app": "cattle-cluster-agent",
          "pod-template-hash": "cc6ddc6dc"
        },
        "ownerReferences": [
          {
            "apiVersion": "apps/v1",
            "kind": "ReplicaSet",
            "name": "cattle-cluster-agent-cc6ddc6dc",
            "uid": "7d4b6cbe-d6d1-46e3-99e5-8410095880c7",
            "controller": true,
            "blockOwnerDeletion": true
          }
        ],
        "managedFields": [
          {
......
}
```
