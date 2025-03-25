---
title: Kubernetes DNS
linkTitle: Kubernetes DNS
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，概念-服务,负载均衡,网络-service 与 pod 的 DNS](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/)

DNS 是 kubernetes 中 **Service Discovery(服务发现，简称 SD)** 的重要实现方式，虽然 K8S SD 还可以通过其他协议和机制提供，但 DNS 是非常常用的、并且是强烈建议的附加组件。

Kubernetes 集群中创建的每个 **service 对象** 和 **pod 对象** 都会被分配一个 **DNS 名称**，Kubernetes 中实现 DNS 功能的程序需要自动创建 **DNS Resource Records(域名解析服务的资源记录)**。基于此，我们可以通过 DNS 名称连接我们部署到集群中的服务，而不用通过 IP 地址。

kubernetes 实现 DNS 的方式：新版本默认使用 [CoreDNS](/docs/10.云原生/Kubernetes/Kubernetes%20网络/Kubernetes%20DNS/CoreDNS/CoreDNS.md)，1.11.0 之前使用的是 kube-dns。Kubernetes DNS 的实现必须符合既定的规范，规范详见 [基于 DNS 的 Kubernetes 服务发现的规范](/docs/10.云原生/Kubernetes/Kubernetes%20网络/Kubernetes%20DNS/基于%20DNS%20的%20Kubernetes%20服务发现的规范.md) 文章。

也就是说，任何可以用于实现 Kubernetes DNS 功能的应用程序，至少需要满足规范中描述的 Resource Records 格式标准。

# Service 对象的 DNS

## A/AAAA 记录

Normal(正常) Service(除了 Headless 类型以外的所有 Service) 会以 `my-svc.my-namespace.svc.cluster-domain.example` 这种名字的形式被分配一个 DNS A 或 AAAA 记录，取决于服务的 IP 协议族。 该名称会解析成对应 Service 对象的 CLUSTER-IP。

Headless(无头) Service(没有 CLUSTER-IP) 也会以 `my-svc.my-namespace.svc.cluster-domain.example` 这种名字的形式分配一个 DNS A 或 AAAA 记录。与 Normal Service 不同的是，这个记录会被解析成对应服务所选择的 Pod 的 IP 集合。

## SRV 记录

Kubernetes 会为命名端口创建 SRV 记录，这些端口是普通服务或 无头服务的一部分。

对每个具有名称端口，SRV 记录具有 `_my-port-name._my-port-protocol.my-svc.my-namespace.svc.cluster-domain.example` 这种形式。 对普通服务，该记录会被解析成端口号和域名：`my-svc.my-namespace.svc.cluster-domain.example`。 对无头服务，该记录会被解析成多个结果，服务对应的每个后端 Pod 各一个； 其中包含 Pod 端口号和形为 `auto-generated-name.my-svc.my-namespace.svc.cluster-domain.example` 的域名。

# Pod 对象的 DNS

## A/AAAA 记录

一般情况下，Pod 对象的域名格式如下：
`Pod-IP-Address.NAMESPACE.pod.ClusterDomain`

例如，如果有一个 Pod 对象，IP 为 172.17.0.3，在 default 名称空间，集群域名为默认的 cluster.local。则该 Pod 的域名为：
`172-17-0-3.default.pod.cluster.local`

## Pod 的 hostname 和 subdomain 字段

当前，创建 Pod 时其主机名取自 Pod 的 `metadata.name` 值。Pod 规约中包含一个可选的 `hostname` 字段，可以用来指定 Pod 的主机名。 当这个字段被设置时，它将优先于 Pod 的名字成为该 Pod 的主机名。 举个例子，给定一个 `hostname` 设置为 "`my-host`" 的 Pod， 该 Pod 的主机名将被设置为 "`my-host`"。Pod 规约还有一个可选的 `subdomain` 字段，可以用来指定 Pod 的子域名。 举个例子，某 Pod 的 `hostname` 设置为 “`foo`”，`subdomain` 设置为 “`bar`”， 在名字空间 “`my-namespace`” 中对应的完全限定域名（FQDN）为 “`foo.bar.my-namespace.svc.cluster-domain.example`”。

示例：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: default-subdomain
spec:
  selector:
    name: busybox
  clusterIP: None
  ports:
    - name: foo # 实际上不需要指定端口号
      port: 1234
      targetPort: 1234
---
apiVersion: v1
kind: Pod
metadata:
  name: busybox1
  labels:
    name: busybox
spec:
  hostname: busybox-1
  subdomain: default-subdomain
  containers:
    - image: busybox:1.28
      command:
        - sleep
        - "3600"
      name: busybox
---
apiVersion: v1
kind: Pod
metadata:
  name: busybox2
  labels:
    name: busybox
spec:
  hostname: busybox-2
  subdomain: default-subdomain
  containers:
    - image: busybox:1.28
      command:
        - sleep
        - "3600"
      name: busybox
```

如果某无头服务与某 Pod 在同一个名字空间中，且它们具有相同的子域名， 集群的 DNS 服务器也会为该 Pod 的全限定主机名返回 A 记录或 AAAA 记录。 例如，在同一个名字空间中，给定一个主机名为 “busybox-1”、 子域名设置为 “default-subdomain” 的 Pod，和一个名称为 “`default-subdomain`” 的无头服务，Pod 将看到自己的 FQDN 为 "`busybox-1.default-subdomain.my-namespace.svc.cluster-domain.example`"。 DNS 会为此名字提供一个 A 记录或 AAAA 记录，指向该 Pod 的 IP。 “`busybox1`” 和 “`busybox2`” 这两个 Pod 分别具有它们自己的 A 或 AAAA 记录。Endpoints 对象可以为任何端点地址及其 IP 指定 `hostname`。

> **说明：**因为没有为 Pod 名称创建 A 记录或 AAAA 记录，所以要创建 Pod 的 A 记录 或 AAAA 记录需要 `hostname`。没有设置 `hostname` 但设置了 `subdomain` 的 Pod 只会为 无头服务创建 A 或 AAAA 记录（`default-subdomain.my-namespace.svc.cluster-domain.example`） 指向 Pod 的 IP 地址。 另外，除非在服务上设置了 `publishNotReadyAddresses=True`，否则只有 Pod 进入就绪状态 才会有与之对应的记录。

## Pod 的 setHostnameAsFQDN 字段

**FEATURE STATE:** `Kubernetes v1.20 [beta]`
**前置条件**：`SetHostnameAsFQDN` 特性门控 必须在 API 服务器 上启用。当你在 Pod 规约中设置了 `setHostnameAsFQDN: true` 时，kubelet 会将 Pod 的全限定域名（FQDN）作为该 Pod 的主机名记录到 Pod 所在名字空间。 在这种情况下，`hostname` 和 `hostname --fqdn` 都会返回 Pod 的全限定域名。

> **说明：**在 Linux 中，内核的主机名字段（`struct utsname` 的 `nodename` 字段）限定 最多 64 个字符。
> 如果 Pod 启用这一特性，而其 FQDN 超出 64 字符，Pod 的启动会失败。 Pod 会一直出于 `Pending` 状态（通过 `kubectl` 所看到的 `ContainerCreating`）， 并产生错误事件，例如 "Failed to construct FQDN from pod hostname and cluster domain, FQDN `long-FQDN` is too long (64 characters is the max, 70 characters requested)." （无法基于 Pod 主机名和集群域名构造 FQDN，FQDN `long-FQDN` 过长，至多 64 字符，请求字符数为 70）。 对于这种场景而言，改善用户体验的一种方式是创建一个 准入 Webhook 控制器， 在用户创建顶层对象（如 Deployment）的时候控制 FQDN 的长度。

## Pod 的 DNS 策略

DNS 策略可以逐个 Pod 来设定。目前 Kubernetes 支持以下特定 Pod 的 DNS 策略。 这些策略可以在 Pod 规约中的 `dnsPolicy` 字段设置：

- "`Default`": Pod 从运行所在的节点继承名称解析配置。参考 相关讨论 获取更多信息。
- "`ClusterFirst`": 与配置的集群域后缀不匹配的任何 DNS 查询（例如 "www.kubernetes.io"） 都将转发到从节点继承的上游名称服务器。集群管理员可能配置了额外的存根域和上游 DNS 服务器。 参阅相关讨论 了解在这些场景中如何处理 DNS 查询的信息。
- "`ClusterFirstWithHostNet`"：对于以 hostNetwork 方式运行的 Pod，应显式设置其 DNS 策略 "`ClusterFirstWithHostNet`"。
- "`None`": 此设置允许 Pod 忽略 Kubernetes 环境中的 DNS 设置。Pod 会使用其 `dnsConfig` 字段 所提供的 DNS 设置。 参见 Pod 的 DNS 配置节。

> **说明：** "Default" 不是默认的 DNS 策略。如果未明确指定 `dnsPolicy`，则使用 "ClusterFirst"。

下面的示例显示了一个 Pod，其 DNS 策略设置为 "`ClusterFirstWithHostNet`"， 因为它已将 `hostNetwork` 设置为 `true`。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: busybox
  namespace: default
spec:
  containers:
    - image: busybox:1.28
      command:
        - sleep
        - "3600"
      imagePullPolicy: IfNotPresent
      name: busybox
  restartPolicy: Always
  hostNetwork: true
  dnsPolicy: ClusterFirstWithHostNet
```

## Pod 的 DNS 配置

Pod 的 DNS 配置可让用户对 Pod 的 DNS 设置进行更多控制。

`dnsConfig` 字段是可选的，它可以与任何 `dnsPolicy` 设置一起使用。 但是，当 Pod 的 `dnsPolicy` 设置为 "`None`" 时，必须指定 `dnsConfig` 字段。用户可以在 `dnsConfig` 字段中指定以下属性：

- `nameservers`：将用作于 Pod 的 DNS 服务器的 IP 地址列表。 最多可以指定 3 个 IP 地址。当 Pod 的 `dnsPolicy` 设置为 "`None`" 时， 列表必须至少包含一个 IP 地址，否则此属性是可选的。 所列出的服务器将合并到从指定的 DNS 策略生成的基本名称服务器，并删除重复的地址。
- `searches`：用于在 Pod 中查找主机名的 DNS 搜索域的列表。此属性是可选的。 指定此属性时，所提供的列表将合并到根据所选 DNS 策略生成的基本搜索域名中。 重复的域名将被删除。Kubernetes 最多允许 6 个搜索域。
- `options`：可选的对象列表，其中每个对象可能具有 `name` 属性（必需）和 `value` 属性（可选）。 此属性中的内容将合并到从指定的 DNS 策略生成的选项。 重复的条目将被删除。

以下是具有自定义 DNS 设置的 Pod 示例：
`service/networking/custom-dns.yaml`

```yaml
apiVersion: v1
kind: Pod
metadata:
  namespace: default
  name: dns-example
spec:
  containers:
    - name: test
      image: nginx
  dnsPolicy: "None"
  dnsConfig:
    nameservers:
      - 1.2.3.4
    searches:
      - ns1.svc.cluster-domain.example
      - my.dns.search.suffix
    options:
      - name: ndots
        value: "2"
      - name: edns0
```

创建上面的 Pod 后，容器 `test` 会在其 `/etc/resolv.conf` 文件中获取以下内容：

    nameserver 1.2.3.4
    search ns1.svc.cluster-domain.example my.dns.search.suffix
    options ndots:2 edns0

对于 IPv6 设置，搜索路径和名称服务器应按以下方式设置：

    kubectl exec -it dns-example -- cat /etc/resolv.conf

输出类似于

    nameserver fd00:79:30::a
    search default.svc.cluster-domain.example svc.cluster-domain.example cluster-domain.example
    options ndots:5
