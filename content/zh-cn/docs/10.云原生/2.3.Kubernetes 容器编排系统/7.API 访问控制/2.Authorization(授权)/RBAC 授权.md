---
title: RBAC 授权
---

概述

> 参考：
>
> - [官方文档,参考-API 访问控制-使用 RBAC 授权](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)
> - [RBAC](/docs/7.信息安全/Access%20Control/RBAC.md) 概念

基于 **Role(角色)** 的访问控制(RBAC) 是一种根据各个用户的角色来控制对 Kubernetes 内资源的访问权限的方法。

在 Kubernetes 的 RBAC 机制中有如下标准化术语：

- **Role(角色)**# 是一组规则的集合，这些规则定义了对 Kubernetes 集群(即 APIserver)的操作权限
  - 权限包括：get、list、watch、create、update、patch、delete
- **Subject(主体)**# 即把 Role 的规则作用于 Subject 上。Subject 就是本文开头讲的 Accounts
  - Subject 类型(kind)包括：User、Group、ServiceAccount
    - 其中 User 就是 认证里的 User Account。User 的名字可以是字符串，也可以是邮件风格的名称，或者以字符串形式表达的数字 ID。
    - Group 的概念是什么还不知道，也没找到参考文档。不过有一个可能应该是这样描述的：
      - Group 与 User 有关系，在创建 User 的证书时，在 subjct 中，O 的值就是表示 Kubernetes RBAC 机制中 Group 的概念。这么看，其实这个 User 与 Group 的概念与 Linux 中用户与组的概念一样。
    - ServiceAccount 详见 [Service Account 详解](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/7.API%20访问控制/1.Authenticating(认证)/Service%20Account%20详解.md)
- **RoleBinding**：定义了 Role 与 Subject 的绑定关系
  - **rules**# 规则，i.e.当前 role 所拥有的权限，其中有 3 个关键字。
  - **apiGroups** # 指定该 role 可以操作那些 api 组。使用 '\*' 表示对所有组具有操作权限
  - **resources**# 指定该 role 可以操作的资源。使用 '\*' 表示对指定组下的所有资源具有操作权限
  - **verbs** # 指定该 role 可以对组内的资源执行什么操作。e.g. get、watch、list 等。使用 '\*' 表示对指定资源具有所有的操作权限

Kubernetes 的 RBAC 机制通过 `rbac.authorization.k8s.io` 这个 API 组实现，该组下一共有 4 个资源可用:

```bash
~]# kubectl get --raw /apis/rbac.authorization.k8s.io/v1 | jq . | grep kind
      "kind": "ClusterRoleBinding",
      "kind": "ClusterRole",
      "kind": "RoleBinding",
      "kind": "Role",
```

1. Role 和 ClusterRole 就是 RBAC 中的 Role(角色) 概念
2. RoleBinding 和 ClusterRoleBinding 用来将 角色 绑定到某些 Subject(主体) 上。

这里将 RBAC 的 Role 概念分为两部分的原因在于，要区分 Role 是作用于集群上，还是作用在指定的名称空间中。比如 Role 可以定义指定名称空间下的权限，而 ClusterRole 则定义所有名称空间下的权限。

# Role 与 ClusterRole

Role 与 ClusterRole 包含一组描述访问权限的规则。凡是被包含的规则即表示允许(不存在拒绝某操作的规则)

1. Role # 名称空间角色，该对象创建后，即表示该角色可以允许对哪个对象执行哪些操作；没有拒绝权限，定义了就是允许，所以称为许可(Permission)授权
2. Cluster Role # 集群角色，该资源定义后，表示集群资源这个角色允许对哪个集群执行哪些操作

## Role Manifests 示例

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: default #定义该Role作用于哪个namespace
  name: pod-reader # 指定该Role的名字
rules: #定义该role的权限规则，即允许Subject对该文件定义的namespace下的pods资源进行GET、WATCH、LIST操作
  # 双引号表示核心组 API，也就是说，改 Role 定义了可以在核心组 API 下可以操作的资源
  # 若想对所有 API 组授权，则使用 * 符号
  - apiGroups: [""]
    resources: ["pods"] #定义规则允许生效的资源都有哪些
    verbs: ["get", "watch", "list"] #定义规则允许进行的动作是哪些
```

上述角色的规则为：允许对 default 名称空间下 核心组 的 Pods 资源，执行 get、watch、list 操作。

## ClusterRole Manifests 示例

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  # "namespace"字段省略因为ClusterRoles不用namespace
  name: secret-reader
rules:
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["get", "watch", "list"]
```

# RoleBinding 与 ClusterRoleBinding

1. RoleBinding # 绑定 Subject 与 Role(也可以是 ClusterRole)，绑定后该 Subject 就有了与之绑定的这个角色的相关权限。
2. ClusterRoleBinding # 绑定 Subject 与 ClusterRole，绑定后该 Subject 就有了与之绑定的这个角色的相关权限

## RoleBinding Manifests 示例

1. RoleBinding 引用（reference 一般简写为 REF）了一个角色，但不包含角色。它可以引用同一命名空间中的角色或全局命名空间中的一个 CR。它通过主题和命名空间信息来添加 WHO 信息，命名空间中存在命名空间。给定命名空间中的角色绑定在该命名空间中仅起作用。

下面的例子是将 k8s 集群中 default 这个 namespace 的用户 lichenhao 绑定到名为 pod-reader 这个 role 上，并具备该 role 上定义的对集群操作的相关权限

如果 subjects.kind 是 ServiceAccount，则表示凡是使用该 ServcieAccount 的 pod，都可以对该名称空间中其他 pod 具有要绑定的 role 中指定的动作。

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: read-pods
  namespace: default
subjects: # 指定“被作用者”的相关信息
  - kind: User # 指定Subject的类型，有 User、ServiceAccount、Group
    name: lichenhao # 指定Subject类型的名字是什么
    apiGroup: rbac.authorization.k8s.io # 指定
roleRef: #定义角色引用，即绑定Role与Subject，让Subject具有Role中定义的规则权限
  kind: Role #定义要绑定的是Role还是ClusterRole
  name: pod-reader # 指定要绑定的Role或ClusterRole的的名字
  apiGroup: rbac.authorization.k8s.io
```

## ClusterRoleBinding Manifests 示例

1. CRB references a CR, but not contain it. It can reference a ClusterRole in the global namespace, and adds who information via Subject.
2. CRB 引用 CR，但不包含 CR。它可以引用全局命名空间中的簇角色，并通过主题添加 WHO 信息。

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: read-secrets-global
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: manager # Name is case sensitive
roleRef:
  kind: ClusterRole
  name: secret-reader
  apiGroup: rbac.authorization.k8s.io
```

注意：kubernetes 的资源分别属于两个级别，一是集群，二是名称空间，一个集群包含多个名称空间，所以 ClusterRole 可以具备其内所有名称空间的权限

# Aggregated ClusterRoles(聚合集群角色)

你可以将若干 ClusterRole **聚合（Aggregate）** 起来，形成一个复合的 ClusterRole。 某个控制器作为集群控制面的一部分会监视带有 `aggregationRule` 的 ClusterRole 对象集合。`aggregationRule` 为控制器定义一个标签 [选择算符](https://kubernetes.io/zh/docs/concepts/overview/working-with-objects/labels/)供后者匹配 应该组合到当前 ClusterRole 的 `roles` 字段中的 ClusterRole 对象。
下面是一个聚合 ClusterRole 的示例：

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: monitoring
aggregationRule:
  clusterRoleSelectors:
    - matchLabels:
        rbac.example.com/aggregate-to-monitoring: "true"
rules: [] # 控制面自动填充这里的规则
```

如果你创建一个与某现有聚合 ClusterRole 的标签选择算符匹配的 ClusterRole， 这一变化会触发新的规则被添加到聚合 ClusterRole 的操作。 下面的例子中，通过创建一个标签同样为 `rbac.example.com/aggregate-to-monitoring: true` 的 ClusterRole，新的规则可被添加到 "monitoring" ClusterRole 中。

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: monitoring-endpoints
  labels:
    rbac.example.com/aggregate-to-monitoring: "true"
# 当你创建 "monitoring-endpoints" ClusterRole 时，
# 下面的规则会被添加到 "monitoring" ClusterRole 中
rules:
  - apiGroups: [""]
    resources: ["services", "endpoints", "pods"]
    verbs: ["get", "list", "watch"]
```

# 默认角色

> 参考：
>
> - [官方文档](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#default-roles-and-role-bindings)，该文档包含一些默认自带的角色的作用详解

API Server 会创建一组默认的 Role、ClusterRole、RoleBinding、ClusterRoleBinding 对象。这些对象的名称大多都是以 `system:` 为前缀，用以标识该资源是直接呦集群控制平面的组件管理的。所有默认的 RBAC 对象都具有 `kubernetes.io/bootstrapping=rbac-defaults` 标签

```bash
~]# kubectl get role -A --show-labels
NAMESPACE     NAME                                             CREATED AT             LABELS
kube-public   system:controller:bootstrap-signer               2021-03-02T07:04:50Z   kubernetes.io/bootstrapping=rbac-defaults
kube-system   extension-apiserver-authentication-reader        2021-03-02T07:04:50Z   kubernetes.io/bootstrapping=rbac-defaults
kube-system   system::leader-locking-kube-controller-manager   2021-03-02T07:04:50Z   kubernetes.io/bootstrapping=rbac-defaults
kube-system   system::leader-locking-kube-scheduler            2021-03-02T07:04:50Z   kubernetes.io/bootstrapping=rbac-defaults
kube-system   system:controller:bootstrap-signer               2021-03-02T07:04:50Z   kubernetes.io/bootstrapping=rbac-defaults
kube-system   system:controller:cloud-provider                 2021-03-02T07:04:50Z   kubernetes.io/bootstrapping=rbac-defaults
kube-system   system:controller:token-cleaner                  2021-03-02T07:04:50Z   kubernetes.io/bootstrapping=rbac-defaults
~]# kubectl get clusterrole --show-labels
NAME                                                                   CREATED AT             LABELS
admin                                                                  2021-03-02T07:04:50Z   kubernetes.io/bootstrapping=rbac-defaults
cloud-controller-manager                                               2021-03-02T07:04:53Z   objectset.rio.cattle.io/hash=5089468545c5482413c7f05e837e9b88f02ad052
cluster-admin                                                          2021-03-02T07:04:50Z   kubernetes.io/bootstrapping=rbac-defaults
edit                                                                   2021-03-02T07:04:50Z   kubernetes.io/bootstrapping=rbac-defaults,rbac.authorization.k8s.io/aggregate-to-admin=true
local-path-provisioner-role                                            2021-03-02T07:04:53Z   objectset.rio.cattle.io/hash=183f35c65ffbc3064603f43f1580d8c68a2dabd4
system:aggregate-to-admin                                              2021-03-02T07:04:50Z   kubernetes.io/bootstrapping=rbac-defaults,rbac.authorization.k8s.io/aggregate-to-admin=true
system:aggregate-to-edit                                               2021-03-02T07:04:50Z   kubernetes.io/bootstrapping=rbac-defaults,rbac.authorization.k8s.io/aggregate-to-edit=true
system:aggregate-to-view                                               2021-03-02T07:04:50Z   kubernetes.io/bootstrapping=rbac-defaults,rbac.authorization.k8s.io/aggregate-to-view=true
system:aggregated-metrics-reader                                       2021-03-02T07:04:53Z   objectset.rio.cattle.io/hash=d795b01b9d5cf4d3744e28995d3303a815726cae,rbac.authorization.k8s.io/aggregate-to-admin=true,rbac.authorization.k8s.io/aggregate-to-edit=true,rbac.authorization.k8s.io/aggregate-to-view=true
system:auth-delegator                                                  2021-03-02T07:04:50Z   kubernetes.io/bootstrapping=rbac-defaults
system:basic-user                                                      2021-03-02T07:04:50Z   kubernetes.io/bootstrapping=rbac-defaults
system:certificates.k8s.io:certificatesigningrequests:nodeclient       2021-03-02T07:04:50Z   kubernetes.io/bootstrapping=rbac-defaults
system:certificates.k8s.io:certificatesigningrequests:selfnodeclient   2021-03-02T07:04:50Z   kubernetes.io/bootstrapping=rbac-defaults
system:certificates.k8s.io:kube-apiserver-client-approver              2021-03-02T07:04:50Z   kubernetes.io/bootstrapping=rbac-
..... 还有很多很多
```

## 默认角色的故障恢复

API Server 每次启动时，会进行 **Auto-reconciliation(自动协商)** 行为，该行为将会更新默认的 RBAC 对象，以修复一些由于意外导致的删除这些 RBAC 对象的故障。比如删除了某些 role 或 clusterrole 对象，重启 API Server 的话，则会恢复这些被删除的对象。

## 面向用户的默认角色

还有一些默认的 ClusterRole 不以 `system:` 为前缀。这些是面向用户的角色，我们可以直接使用这些 ClusterRole 绑定到需要的 Subjects 上，从而直接为其赋权，省去了自己创建 ClusterRole 的烦恼。

1. **cluster-admin**# 具有最高权限的 clusterrole，与此集群角色绑定的 Subject(主体) 具有对集群操作的最高权限
   1. 默认情况下，与 ClusterRole 同名的 ClusterRoleBinding，将该 ClusterRole 绑定到了 `system:masters` 组上，所有隶属于此组的用户都将具有集群的超级管理权限，比如 /etc/kubernets/pki/ 目录下的相关 \*.conf 文件里的证书通过 base64 解码并查看证书时，其中 subject 的信息有一条 `O=system:masters`. ，在进行认证的时候，CN 的值可作为用户名使用，O 的值将作为用户所属组名使用，因此所有使用这类证书的事物在访问集群时，都具有最高权限
   2. kubectl 命令所用的 kubeconfig 文件中的客户端证书的 subject 字段就有 O，并且值为 `system:masters`
2. **admin**# 允许管理员访问权限，旨在使用  **RoleBinding**  在名字空间内执行授权。 如果在  **RoleBinding**  中使用，则可授予对名字空间中的大多数资源的读/写权限， 包括创建角色和角色绑定的能力。 但是它不允许对资源配额或者名字空间本身进行写操作。
3. **edit**# 允许对名字空间的大多数对象进行读/写操作。 它不允许查看或者修改角色或者角色绑定。 不过，此角色可以访问 Secret，以名字空间中任何 ServiceAccount 的身份运行 Pods， 所以可以用来了解名字空间内所有服务账户的 API 访问级别。
4. **view** # 允许对名字空间的大多数对象有只读权限。 它不允许查看角色或角色绑定。此角色不允许查看 Secrets，因为读取 Secret 的内容意味着可以访问名字空间中 ServiceAccount 的凭据信息，进而允许利用名字空间中任何 ServiceAccount 的 身份访问 API（这是一种特权提升）。

核心组件的角色

其他组件的角色

内置控制器的角色

# 通过 kubectl 创建 RBAC

详见：[kubectl 命令行工具的 create 子命令](/docs/10.云原生/2.3.Kubernetes%20 容器编排系统/Kubernetes%20 管理/kubectl%20 命令行工具/对象的创建与修改命令.md 管理/kubectl 命令行工具/对象的创建与修改命令.md)
