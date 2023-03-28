---
title: Service(服务)
---

# 概述

> 参考：
>
> - [官方文档，概念-服务，负载均衡，网络-服务](https://kubernetes.io/docs/concepts/services-networking/service/)

Service 资源可以将一组运行在 Pods 上的应用程序暴露为网络服务，这是一种抽象的功能，说白了 Service 资源实现 服务发现，执行访问 POD 的任务、4 层代理 等功能

为什么要使用 Service？

Deployment 可以部署多个副本，每个 Pod 都有自己的 IP，外界如何访问这些副本呢？通过 Pod 的 IP 吗？要知道 Pod 很可能会被频繁地销毁和重启，它们的 IP 会发生变化，用 IP 来访问不太现实。答案是 Service。Service 作为访问 Pod 的接入层来使用。service 就像 lvs 的 director 一样，充当一个调度器的作用。

Service 定义了外界访问一组特定 Pod 的方式。Service 有自己的 IP 和 PORT ，Service 为 Pod 提供了负载均衡。

可以把 Service 想象成负载均衡功能的前端，该 Service 下的 Pod 是负载均衡功能的后端。通过其自动创建的 ipvs 或者 iptables 的规则，访问 Service 的 IP:PORT，然后转发数据到后端的 Pod

## Endpoints

注意：在 service 与 pod 中间，还有一个中间层，这个中间层就是 Endpoints 资源。

Endpoints 是一个由 IP 和 PORT 组成的 endpoint 列表，Endpoints 的这些 IP 和 PORT 来自于由 Service 的标签选择器匹配到的 pod 资源(也可在 service 不适用标签选择器的时候，手动指定 endpoints 的 IP 与 PORT)。默认情况下，创建 Service 资源时，会自动创建同名的 Endpoints 资源。从抽象角度看，service 所关联的每一个 pod 其实都是 Endpoints 资源列表中的一个 endpoint

# Service 的实现

Service 是 k8s 中的一种资源，但是如果想要实现 Service 资源定义的那些内容，则需要 [kube-proxy](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/8.Kubernetes%20网络/kube-proxy(实现%20Service%20 功能的组件).md)) 这个程序，实际上，创建一个 service，就是让 kube-proxy 创建一系列 iptables 或者 ipvs 规则

- userspace：在 1.1.0 版本之前使用该模型，由于需要把报文转到内核空间再回到用户空间过于低效
- iptables：通过 Kube-proxy 监听 Pod 变化在宿主机上生成并维护。由于 kube-proxy 需要在总之循环里不断得刷新 iptables 规则来确保它们始终是正确的，这样当 Host 上有大量 Pod 时，会产生极多 Iptables 规则，大量占用 Host 的 CPU 资源，这时候 ipvs 模型就可以解决该问题
- ipvs：将负载均衡与代理功能从 iptables 手中接过来，一些辅助性并且数量不多的(比如包过滤、SNAT 等)操作依然由 iptables 完成。如果想要启用 ipvs 工作模型，那么需要在/etc/sysconfig/kubelet 该配置文件中加入 KUBE_PROXY_MODE=ipvs 这一行，且给 linux 装载 ipvs 模块和连接跟踪模块

# Publishing Service(发布 Service)

发布 Service 是指将 Service 暴露出去以供其他客户端访问他，并将请求转给其所关联的后端 Pod。Service 有多种类型，不同的类型对应不同的发布 Service 方式，默认的类型为 `ClusterIP`

- ClusterIP # 通过集群的内部 IP 暴露 Service。该 Service 只能被暴露在进群内部，集群外部无法访问。
  - Service 暴露的集群内部 IP 由 kube-controller-manager 程序的 `--service-cluster-ip-range` 标志控制。
- NodePort # 在集群中每个节点 IP 的静态端口上暴露 Service 给集群外部，每个节点上将会创建到 Service 的 CluseterIP 的路由条目，以便我们从集群外部访问 NodePort 类型的 Service
- LoadBalancer # 使用云提供商的负载均衡器暴露 Service 给集群外部。
- ExternalName #

我们还可以使用 Kubernetes 的 Ingress 资源暴露 Service。

## ClusterIP：仅用于 kubernetes 集群内通信

每个 Service 创建完成后一般都会有一个 cluster-ip(headless 类型的 service 没有)，这个 IP 是 Kubernetes 集群的专用 IP，是一种虚拟 IP，可以把它当做 lvs 中的 vip。只不过这些 IP 并不能直接访问到，而是在 Service 创建完成后，在 iptables 或者 ipvs 规则中所使用的 IP。Kubernetes 创建完成后，cluster-ip 默认的使用范围是 10.96.0.0/12

- headless：无头服务，当不需要使用负载均衡和单一服务 IP 的时候，可以给 ClusterIP 设为 None。kube-proxy 不使用这些服务并且平台(platform)没有负载均衡和代理
  - headless 由于没有 cluster-ip，所以是通过域名的方式来让外部访问到该 service 的 endpoint 的，如果有 cluster-ip 的话，则 ServiceName.NameSpaceName.svc.cluster.local 的域名会解析到该 service 的 cluster-ip 上，如果是 headless 的话，域名解析的结果则是所有 endpoint 的 ip，客户端每次向此 headless 类型的 service 发起的请求，将直接接入到各个 endpoint 上，不再由 service 资源进行代理转发，而是由 DNS 服务器收到查询请求时以轮训的方式返回各个 endpoint 的 IP。

## NodePort：用于从集群外部访问 Service

通过 kube-proxy 添加 iptables 规则，把流量通过主机的 port 转发到 Service 的 port 上。

NodePort 建立在 ClusterIP 类型之上，NodePort 会将宿主机的 port 与 service 的 port 所关联，这样就可以将 service 乃至其 endpoint 都可以让集群外部直接访问。如果定义 NodePort 时不指定，则会随机选择宿主机上的 30000 至 32767 之间的一个端口作为 NodePort 的 PORT。

比如下图画红框的部分就表示冒号前是 service 的 port，冒号后是宿主机上的 port，当访问宿主机的 port 的时候，该访问请求会被 iptables 或者 ipvs 规则转发到 service 的 port 上，然后转交给其 endpoint
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fd3b0g/1616118462642-7b49b026-38af-4641-84db-d58f9cbfbed8.jpeg)
但是，请注意！NodePort 有个很致命的确定，详见 [kube-proxy 无法绑定 NodePort 案例](https://www.yuque.com/go/doc/44843491)

## LoadBalancer：一般当 kubernetes 部署在云上时使用

LoadBalancer 建立在 NodePort 之上，可以将实现 Service 资源的默认负载均衡器(i.e.kube-proxy)替换为其他的负载均衡器。这时可以直接通过负载均衡器暴露的 IP + PORT 直接访问 Service 后端关联的 Pod。

现阶段，想让 Service 对接外部负载均衡器，在指定 Service 的 `spec.type` 字段的值为 LoadBalancer 以外，还需要配置 `metadata.annotations` 字段，以便让外部负载均衡器的控制器获取 Service 信息以便对自己进行配置。

当我们创建了一个 LoadBalancer 类型的 Service 后，该 Sevice 对象将会获得一个 External-IP，通常这个 IP 是由外部负载均衡器的控制器提供的：

```bash
kubectl get svc nginx
NAME    TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)        AGE
nginx   LoadBalancer   10.100.126.91   192.168.0.101   80:31555/TCP   4s
```

此时，通过 LB 控制器，可以让 LB 自动关联到 Service 后端的 Pod 上，通常来说 `192.168.0.101` 是 LB 设备的 VIP，访问 `192.168.0.101` 的 31555 端口会自动负载均衡到 nginx Service 后端的 Pod

常见的负载均衡器：

- 各大公有云厂商的 LB。比如华为的 ELB、等
- [MetalLB](https://github.com/metallb/metallb) #
  - [公众号-运维开发故事，Kubernetes 开源 LoadBalancer-Metallb(BGP)](https://mp.weixin.qq.com/s/BY6hrLjaWfPYJzYmpbl1fQ)
- [OpenELB](https://github.com/openelb/openelb) #

## ExternalName：把集群外部的服务引入到集群内部

## External IPs(外部 IP)

为 Service 对象配置 `spec.externalIPs` 字段后，会在节点上创建一个 ipvs 条目，Director 为 externalIPs 的值，RealServer 为 Service 关联的后端 Pod 的 IP。此时，只要为客户端配置一条路由规则，目的地址是 ExternalIP 的包都转发给 K8S 的节点，就可以从集群外部访问 Service 了。

> 注意：externalIPs 不受 Kubernetes 管理

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/fd3b0g/1649947301381-c990dbc0-9232-4171-ac34-905219da5d8a.png)
External IP 还会在使用 LoadBalancer 类型的 Service 时，自动被公有云厂商的 LB 填充，通常都是将 ingress-controller-nginx 的 Service 配置为 LoadBalancer 类型以对接公有云厂商的 LB。

## 手动指定 Endpoints

不指定 selector，手动创建一个与 Service 同名的 Endpoints，这样就能实现手动指定该 Service 所关联的 Endpoints
