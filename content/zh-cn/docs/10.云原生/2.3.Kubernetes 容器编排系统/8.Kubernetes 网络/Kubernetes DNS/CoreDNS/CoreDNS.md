---
title: CoreDNS
---

# 概述

> 参考：
> 
> - [官网](https://coredns.io/)
> - [官方手册](https://coredns.io/manual/toc/)
> - [CoreDNS 所有插件详解](https://coredns.io/plugins/)
> - [k8s 中的 CoreDNS 的配置示例](https://kubernetes.io/docs/tasks/administer-cluster/dns-custom-nameservers/)

CoreDNS 是一个用 Go 编写的 DNS 服务器，目标是要称为云原生环境下的 DNS 服务器和服务发现解决方案。

CoreDNS 与 [BIND](https://zh.wikipedia.org/zh-cn/BIND) 这类 DNS 服务器不同，CoreDNS 非常灵活，几乎所有功能都由插件来实现

插件可以独立或者一起运行，以便执行一个 **DNS Function(DNS 功能)。**可以说 CoreDNS 是由插件驱动的

那么什么是“ DNS 功能”？ 出于 CoreDNS 的目标，我们将其定义为实现 CoreDNS Plugin API 的软件。 实现的功能可能会大相径庭。 有些插件本身并不会创建响应，例如指标或缓存，但会添加功能。 然后有一些插件确实会产生响应。 这些也可以做任何事情：有一些与 Kubernetes 通信以提供服务发现的插件，有一些从文件或数据库中读取数据的插件。

So what’s a “DNS function”? For the purpose of CoreDNS, we define it as a piece of software that implements the CoreDNS Plugin API. The functionality implemented can wildly deviate. There are plugins that don’t themselves create a response, such as [metrics](https://coredns.io/plugins/metrics) or [cache](https://coredns.io/plugins/cache), but that add functionality. Then there are plugins that _do_ generate a response. These can also do anything: There are plugins that communicate with [Kubernetes](https://coredns.io/plugins/kubernetes) to provide service discovery, plugins that read data from a [file](https://coredns.io/plugins/file) or a [database](https://coredns.io/explugins/pdsql).

用白话说：

它有以下几个特性：

- 插件化（Plugins）基于 Caddy 服务器框架，CoreDNS 实现了一个插件链的架构，将大量应用端的逻辑抽象成 plugin 的形式（如 Kubernetes 的 DNS 服务发现，Prometheus 监控等）暴露给使用者。CoreDNS 以预配置的方式将不同的 plugin 串成一条链，按序执行 plugin 的逻辑。从编译层面，用户选择所需的 plugin 编译到最终的可执行文件中，使得运行效率更高。CoreDNS 采用 Go 编写，所以从具体代码层面来看，每个 plugin 其实都是实现了其定义的 interface 的组件而已。第三方只要按照 CoreDNS Plugin API 去编写自定义插件，就可以很方便地集成于 CoreDNS。
- 配置简单化引入表达力更强的 DSL\[2]，即 `Corefile` 形式的配置文件（也是基于 Caddy 框架开发）。
- 一体化的解决方案区别于 `kube-dns`，CoreDNS 编译出来就是一个单独的二进制可执行文件，内置了 cache，backend storage，health check 等功能，无需第三方组件来辅助实现其他功能，从而使得部署更方便，内存管理更为安全。

其实从功能角度来看，CoreDNS 更像是一个通用 DNS 方案（类似于 `BIND`），然后通过插件模式来极大地扩展自身功能，从而可以适用于不同的场景（比如 Kubernetes）。正如官方博客所说：

## CoreDNS Metrics

在 K8S 中，CoreDNS 不在宿主机暴露端口，需要通过其 service 来访问

指标获取路径：CoreDNS_SVC_IP:9153/metrics

# CoreDNS 关联文件与配置

**./Corefile** # CoreDNS 运行所需配置文件，参考：[**Corefile 解释**](https://coredns.io/2017/07/23/corefile-explained/)。

## kubeadm 安装的 k8s 集群中 coredns 的默认配置文件

在 kubeadm 安装的集群中，coredns 的配置保存在 configmap 中，通过 kubectl 命令进行查看

```yaml
[root@master myapp]# kubectl get configmaps -n kube-system coredns -o yaml
apiVersion: v1
kind: ConfigMap
data:
  Corefile: |
    .:53 {
        errors
        health
        ready
        kubernetes cluster.local in-addr.arpa ip6.arpa {
           pods insecure
           fallthrough in-addr.arpa ip6.arpa
           ttl 30
        }
        prometheus :9153
        forward . /etc/resolv.conf
        cache 30
        loop
        reload
        loadbalance
    }
```

# CoreDNS Plugins(插件)

参考：[官方 Plugins 手册](https://coredns.io/manual/toc/#plugins)

## CoreDNS Plugins 的工作模式

当 CoreDNS 启动后，它将根据配置文件启动不同 server ，每台 server 都拥有自己的插件链。当有 DNS 请求时，它将依次经历如下 3 步逻辑：

- 如果有当前请求的 server 有多个 zone，将采用贪心原则选择最匹配的 zone；
- 一旦找到匹配的 server，按照 [**plugin.cfg**](https://github.com/coredns/coredns/blob/master/plugin.cfg) 定义的顺序执行插件链上的插件；
  - plugin.cfg 是定义在代码中的，在处理请求时，总是根据 plugin.cfg 中定义的顺序加载插件
- 每个插件将判断当前请求是否应该处理，将有以下几种可能：
  - **请求被当前插件处理**插件将生成对应的响应并回给客户端，此时请求结束，下一个插件将不会被调用，如 whoami 插件；
  - **请求被当前插件以 Fallthrough 形式处理**如果请求在该插件处理过程中有可能将跳转至下一个插件，该过程称为 fallthrough，并以关键字 `fallthrough` 来决定是否允许此项操作，例如 host 插件，当查询域名未位于 /etc/hosts，则调用下一个插件；
  - **请求在处理过程被携带 Hint**请求被插件处理，并在其响应中添加了某些信息（hint）后继续交由下一个插件处理。这些额外的信息将组成对客户端的最终响应，如 `metric` 插件；

## CoreDNS 的常用插件

参考：[官方 Plugins 列表文档](https://coredns.io/plugins/)

- [errors](https://coredns.io/plugins/errors/) # coredns 查询处理过程中遇到的任何错误都将被打印到标准输出。Note：每个配置文件中，仅能使用一次 errors 插件
- [health](https://coredns.io/plugins/health/) # 该插件用来将 CoreDNS 的运行状态暴露在 http://localhost:8080/health。
- [ready](https://coredns.io/plugins/ready/) # 当所有插件都就绪之后，会在 8181 端口上返回 http 的 200 状态码。
- [kubernetes](https://coredns.io/plugins/kubernetes/) # CoreDNS 将基于 Kubernetes 的服务和 Pod 的 IP 答复 DNS 查询。
- [prometheus](https://coredns.io/plugins/metrics/) # 在 9153 端口上以 OpenMetrics 格式暴露 coredns 指标。
- [forward](https://coredns.io/plugins/forward/) # 任何不在 Kubernetes 集群域内的查询都将转发到预定义的解析器（/etc/resolv.conf）。
- [cache](https://coredns.io/plugins/cache/) # 启用前端缓存。
- [loop](https://coredns.io/plugins/loop/) # 检测简单的转发循环，如果发现循环，则中止 CoreDNS 进程。
- [reload](https://coredns.io/plugins/reload) # 允许自动重新加载已更改的 Corefile。 编辑 ConfigMap 配置后，等待两分钟，以使更改生效。
- [loadbalance](https://coredns.io/plugins/loadbalance) # 一个轮询 DNS 负载均衡器，它随机分配 dns 响应中的 A，AAAA 和 MX 记录的顺序。
- 等等等

## kubernetes 插件

kubernetes \[ZONES...] {    endpoint URL    tls CERT KEY CACERT    kubeconfig KUBECONFIG CONTEXT    namespaces NAMESPACE...    labels EXPRESSION    pods POD-MODE    endpoint_pod_names    ttl TTL # 设置自定义 TTL。 默认值为 5 秒。 允许的最小 TTL 为 0 秒，最大为 3600 秒。 将 TTL 设置为 0 将防止记录被缓存。    noendpoints    transfer to ADDRESS...    fallthrough \[ZONES...]    ignore empty_service}

pods POD-MODE # 提供了 pods insecure 选项，以便与 kube-dn 向后兼容。 您可以使用经过验证的 Pod 选项，只有在相同名称空间中存在具有匹配 IP 的 Pod 时，该选项才返回 A 记录。 如果您不使用广告连播记录，则可以使用“广告连播禁用”选项。

# CoreDNS 的工作模式

如果 Corefile 为：

```bash
coredns.io:5300 {
    file db.coredns.io
}

example.io:53 {
    log
    errors
    file db.example.io
}

example.net:53 {
    file db.example.net
}

.:53 {
    kubernetes
    proxy . 8.8.8.8
    log
    health
    errors
    cache
}
```

从配置文件来看，我们定义了两个 server（尽管有 4 个区块），分别监听在 `5300` 和 `53` 端口。其逻辑图可如下所示：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gggg2s/1622809781189-17c30eb0-ece4-4e93-b40a-34714856af7f.png)
每个进入到某个 server 的请求将按照 [**plugin.cfg**](https://github.com/coredns/coredns/blob/master/plugin.cfg) 定义顺序执行其已经加载的插件。

从上图，我们需要注意以下几点：

- 尽管在 `.:53` 配置了 `health` 插件，但是它并为在上面的逻辑图中出现，原因是：该插件并未参与请求相关的逻辑（即并没有在插件链上），只是修改了 server 配置。更一般地，我们可以将插件分为两种：
  - **Normal 插件**：参与请求相关的逻辑，且插入到插件链中；
  - **其他插件**：不参与请求相关的逻辑，也不出现在插件链中，只是用于修改 server 的配置，如 `health`，`tls` 等插件；
