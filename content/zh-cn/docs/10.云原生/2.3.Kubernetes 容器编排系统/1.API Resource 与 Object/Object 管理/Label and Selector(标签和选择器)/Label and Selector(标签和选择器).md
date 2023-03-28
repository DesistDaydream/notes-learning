---
title: Label and Selector(标签和选择器)
weight: 1
---

# 概述

> 参考：
>
> - [官方文档,概念-使用 Kubernetes 对象-标签和选择器](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/)
> - [官方文档,概念-使用 Kubernetes 对象-推荐的标签](https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/)
> - [官方文档,参考-常见的内置标签、注释、污点](https://kubernetes.io/docs/reference/labels-annotations-taints/)
> - [官方文档,参考-Kubernets API-通用定义-标签选择器](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/label-selector/)
> - [官方文档,参考-Kubernets API-通用定义-节点选择器请求](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/node-selector-requirement/)

**Label(标签)** 是 `键/值对` 的集合，在 Kubernetes 中，每一个对象都可以具有一个或多个 **Label(标签)**。Label 主要用来让用户定义对象的属性，以便为所有对象进行分类，并且还可以组织和选择对象的子集。标签可以在创建对象的同时添加，也可以随时修改对象上的标签。

Kubernetes 中的 Label 功能与 时间序列数据 中标签功能有异曲同工之妙，说白了，就是用来描述一个东西的。而且通过 Label，我们可以以松耦合的方式将我们自己想要的组织方式组织集群中的 Pod，而并不需要自己维护这些。

**Kubernetes 中标签概念的重要性不亚于 API 资源和对象的概念。Pod 要运行在哪个 Node 上、下文将会提到的标签选择器，以及 Kubernetes 的调度系统等等等等，想要实现这些功能，都要依赖于标签。**

_标签_ 是键值对。有效的标签键有两个段：可选的前缀和名称，用斜杠 `/` 分隔。 名称段是必需的，必须小于等于 63 个字符，以字母数字字符（`[a-z0-9A-Z]`）开头和结尾， 带有破折号 `-`，下划线 `_`，点  `.` 和之间的字母数字。 前缀是可选的。如果指定，前缀必须是 DNS 子域：由点 `.` 分隔的一系列 DNS 标签，总共不超过 253 个字符， 后跟斜杠 `/`。

如果省略前缀，则假定标签键对用户是私有的。 向最终用户对象添加标签的自动系统组件（例如 `kube-scheduler`、`kube-controller-manager`、 `kube-apiserver`、`kubectl` 或其他第三方自动化工具）必须指定前缀。

有效标签值：

- 必须为 63 个字符或更少（可以为空）
- 除非标签值为空，必须以字母数字字符（`[a-z0-9A-Z]`）开头和结尾
- 包含破折号 `-`、下划线 `_`、点 `.` 和字母或数字。

## 集群中特殊的标签

`kubernetes.io/`  和 `k8s.io/` 前缀是为 Kubernetes 核心组件保留的。

在每个 master 上都会有这么一个标签 `node-role.kubernetes.io/XXXXX=`，在使用 `kubeclt get nodes` 命令时， ROLES 列的值，就是根据该标签的 key 来决定的，key 中 XXXX 的值，会填写到 ROLES 列中。

## 一个对象中标签示例

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: label-demo
  labels:
    environment: production
    app.kubernetes.io/name: nginx
spec:
  containers:
    - name: nginx
      image: nginx:1.14.2
```

# Label Selector(标签选择器，简称 Selector)

可以给 kubernetes 中的对象打上标签，然后让某个对象使用 Selector 来选择具有相同标签的对象成为同一组来协调工作或者进行各种限定

比如具有相同标签的 Pod 和 Node，该 Pod 会使用 Selector 选择在该 Node 上运行，该 Pod 对该 Node 具有倾向性；或者把具有相同标签的 Service 和 Pod 关联起来，使 Service 使用 Selector 知道可以选择哪些 Pod 来进行调度

说白了，所谓的 Selector 就是根据给定的规则对标签进行匹配，凡是带有匹配到的标签的资源，都会被选择器选中。

标签选择器可以用在下面这些资源的字段中：

- 各种控制器
    - .spec.selector
- pod
    - .spec.affinity.所有亲和类型.软/硬规则.nodeSelectorTerms
    - .spec.nodeSelector
- service
    - .spec.selector
- 等等

## Selector Manifest

选择器的 Manifest 字段及其写法详见 [API 参考-LabelSelector](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/1.API%20Resource%20与%20Object/API%20参考/Common%20Definitions(通用定义)/LabelSelector.md)

# Label 的使用方式

## API

### LIST 和 WATCH 过滤

LIST 和 WATCH 操作可以使用查询参数指定标签选择算符过滤一组对象。 两种需求都是允许的。（这里显示的是它们出现在 URL 查询字符串中）

- _基于等值_ 的需求: `?labelSelector=environment%3Dproduction,tier%3Dfrontend`
- _基于集合_ 的需求: `?labelSelector=environment+in+%28production%2Cqa%29%2Ctier+in+%28frontend%29`

两种标签选择算符都可以通过 REST 客户端用于 list 或者 watch 资源。 例如，使用 `kubectl` 定位 `apiserver`，可以使用 _基于等值_ 的标签选择算符可以这么写：

```bash
kubectl get pods -l environment=production,tier=frontend
```

或者使用 _基于集合的_ 需求：

```bash
kubectl get pods -l 'environment in (production),tier in (frontend)'
```

正如刚才提到的，_基于集合_ 的需求更具有表达力。例如，它们可以实现值的 _或_ 操作：

```bash
kubectl get pods -l 'environment in (production, qa)'
```

或者通过 _exists_ 运算符限制不匹配：

```bash
kubectl get pods -l 'environment,environment notin (frontend)'
```

### 在 API 对象中设置引用

一些 Kubernetes 对象，例如 `[services](https://kubernetes.io/zh/docs/concepts/services-networking/service/)` 和 `[replicationcontrollers](https://kubernetes.io/zh/docs/concepts/workloads/controllers/replicationcontroller/)` ， 也使用了标签选择算符去指定了其他资源的集合，例如 [pods](https://kubernetes.io/zh/docs/concepts/workloads/pods/)。

#### Service 和 ReplicationController

一个 `Service` 指向的一组 Pods 是由标签选择算符定义的。同样，一个 `ReplicationController` 应该管理的 pods 的数量也是由标签选择算符定义的。
两个对象的标签选择算符都是在 `json` 或者 `yaml` 文件中使用映射定义的，并且只支持 _基于等值_ 需求的选择算符：

```json
"selector": {
    "component" : "redis",
}
```

或者

```yaml
selector:
  component: redis
```

这个选择算符(分别在 `json` 或者 `yaml` 格式中) 等价于 `component=redis` 或 `component in (redis)` 。

#### 支持基于集合需求的资源

比较新的资源，例如 `[Job](https://kubernetes.io/zh/docs/concepts/workloads/controllers/job/)`、 `[Deployment](https://kubernetes.io/zh/docs/concepts/workloads/controllers/deployment/)`、 `[Replica Set](https://kubernetes.io/zh/docs/concepts/workloads/controllers/replicaset/)` 和 `[DaemonSet](https://kubernetes.io/zh/docs/concepts/workloads/controllers/daemonset/)` ， 也支持 _基于集合的_ 需求。

```yaml
selector:
  matchLabels:
    component: redis
  matchExpressions:
    - { key: tier, operator: In, values: [cache] }
    - { key: environment, operator: NotIn, values: [dev] }
```

`matchLabels` 是由 `{key,value}` 对组成的映射。 `matchLabels` 映射中的单个 `{key,value }` 等同于 `matchExpressions` 的元素， 其 `key` 字段为 "key"，`operator` 为 "In"，而 `values` 数组仅包含 "value"。 `matchExpressions` 是 Pod 选择算符需求的列表。 有效的运算符包括 `In`、`NotIn`、`Exists` 和 `DoesNotExist`。 在 `In` 和 `NotIn` 的情况下，设置的值必须是非空的。 来自 `matchLabels` 和 `matchExpressions` 的所有要求都按逻辑与的关系组合到一起 -- 它们必须都满足才能匹配。

#### 选择节点集

通过标签进行选择的一个用例是确定节点集，方便 Pod 调度。 有关更多信息，详见 [让 Pod 运行在指定 Node](https://www.yuque.com/go/doc/33166071) 章节。

# 使用 kubectl 命令控制标签

## Syntax(语法)

**kubectl label \[--overwrite] (-f FILENAME | TYPE NAME) KEY_1=VAL_1 ... KEY_N=VAL_N \[--resource-version=version] \[options]**

在 get 子命令中，`--show-label` 标志还可以显示获取到的对象的所有标签；`-l` 标签可以根据表达式来过滤想要获取的对象

## EXAMPLE

- 获取所有 node 的标签
  - **kubectl get node --show-labels**
- 为 node-1.bj-test 节点添加名为 `node-role.kubernetes.io/proxy` 的标签，标签值为 `ingress-controller`
  - **kubectl label nodes test-node-4 node-role.kubernetes.io/proxy=ingress-controller**
- 将 node-1.bj-test 节点上的 `node-role.kubernetes.io/proxy` 标签删除
  - **kubectl label nodes node-1.bj-test node-role.kubernetes.io/proxy-**
- 给 k8s-node1 节点添加 disktype=ssd 这个标签
  - **kubectl label node k8s-node1 disktype=ssd**
- 删除 k8s-node1 节点上的 disktype 标签
  - **kubectl label node k8s-node1 disktype-**
