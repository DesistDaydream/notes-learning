---
title: kubectl port-forward 工作原理
---

原文链接：[公众号-CNCF，源码解析 kubectl port-forward 工作原理](https://mp.weixin.qq.com/s/cFxs8cseuXfO7llM4KAoVg)

本文的源码基于 Kubernetes v1.24.0，容器运行时使用 Containerd 1.5，从源码来分析 kubectl port-forward 的工作原理。

通过 port-forward 流程的分析，梳理出 kubectl -> api-server -> kubelet -> 容器运行时 的交互，了解 cri 的工作方式。

![kubectl-port-forward](https://notes-learning.oss-cn-beijing.aliyuncs.com/kubernetes/source/port-forward.png)
## kubectl

简单创建个 pod：

`kubectl run pipy --image flomesh/pipy:latest -n default`

在执行  `kubectl forward`  时添加参数  `-v 9`  打印日志。

```bash
kubectl port-forward pipy 8080 -v 9
...
I0807 21:45:58.457986   14495 round_trippers.go:466] curl -v -XPOST  -H "User-Agent: kubectl/v1.24.3 (darwin/arm64) kubernetes/aef86a9" -H "X-Stream-Protocol-Version: portforward.k8s.io" 'https://192.168.1.12:6443/api/v1/namespaces/default/pods/pipy/portforward'
I0807 21:45:58.484013   14495 round_trippers.go:553] POST https://192.168.1.12:6443/api/v1/namespaces/default/pods/pipy/portforward 101 Switching Protocols in 26 milliseconds
I0807 21:45:58.484029   14495 round_trippers.go:570] HTTP Statistics: DNSLookup 0 ms Dial 0 ms TLSHandshake 0 ms Duration 26 ms
I0807 21:45:58.484035   14495 round_trippers.go:577] Response Headers:
I0807 21:45:58.484040   14495 round_trippers.go:580]     Upgrade: SPDY/3.1
I0807 21:45:58.484044   14495 round_trippers.go:580]     X-Stream-Protocol-Version: portforward.k8s.io
I0807 21:45:58.484047   14495 round_trippers.go:580]     Date: Sun, 07 Aug 2022 13:45:58 GMT
I0807 21:45:58.484051   14495 round_trippers.go:580]     Connection: Upgrade
Forwarding from 127.0.0.1:8080 -> 8080
Forwarding from [::1]:8080 -> 8080
```

从日志可以看到请求的地址为  `/api/v1/namespaces/default/pods/pipy/portforward`，其中  `portforward`  为 pod 资源的子资源。

> 这里使用的协议是 spdy。

`kubectl`  此时会监听本地端口，同时使用 pod 子资源 portforward 的 url 创建到 api-server 的连接。

当本地端口有连接接入时，`kubectl`  会**不断地在两个连接间拷贝数据**。

### 参考源码：

- staging/src/k8s.io/kubectl/pkg/cmd/portforward/portforward.go:389\[1]
- staging/src/k8s.io/client-go/tools/portforward/portforward.go:242\[2]
- staging/src/k8s.io/client-go/tools/portforward/portforward.go:330\[3]

## api-server

pod 的三个子资源 exec、attach 和 portforward，对这三个资源的操作都会代理有对应 node 的 kubetlet server 进行处理。

api-server 在接收到访问 pod 子资源 portforward 的请求后，通过 pod 及其所在 node 的信息，获取访问该 node 上 kubelet server 的 url。

然后将访问 pod 的 portforward 的请求，代理到 kubelet server。

### 参考源码

- pkg/registry/core/pod/rest/subresources.go:185\[4]

## kubelet

portforward 请求来到了 pod 所在节点的 kubelet server，在 kubelet server 中，有几个用于调试的 endpoint，portforward 便是其中之一：

- `/run/{podNamespace}/{podID}/{containerName}`
- `/exec/{podNamespace}/{podID}/{containerName}`
- `/attach/{podNamespace}/{podID}/{containerName}`
- `/portforward/{podNamespace}/{podID}`
- `/containerLogs/{podNamespace}/{podID}/{containerName}`
- `/runningpods/`

kubelet server 收到请求后，首先会通过  `RuntimeServiceClient`  发送 gRCP 请求到容器运行时的接口（`/runtime.v1alpha2.RuntimeService/PortForward`）获取容器运行时 streaming server 处理 pordforward 请求的 url。

拿到 portforward streaming 的 url 之后，kubelet server 将请求代理到该 url。

### 参考源码

- pkg/kubelet/server/server.go:463\[5]
- pkg/kubelet/server/server.go:873\[6]
- pkg/kubelet/cri/streaming/portforward/portforward.go:46\[7]
- pkg/kubelet/cri/streaming/server.go:111\[8]

## cri

这里以 Containerd 为例。

Containerd 在启动时会启动 runtime service 和 image service。前者是负责容器相关的操作，后者负责镜像相关的操作。

**kubelet 获取用于端口转发的 streaming url，就是调用了 runtime service 的 gRPC 接口完成的。**

除了两个 gRPC service 以外，还加载了一系列插件。这些插件中，其中有一个是 cri service。

cri service 会启动 streaming server。这个 server 会响应  `/exec`、`/attach`  和  `/portforward`  的 stream 请求。

portforward 支持两种操作系统 linux 和 windows：`sandbox_portforward_linux.go`  和  `sandbox_portforward_windows.go`。

在 linux 上，在 pod 所在的 network namespace 中使用地址  `localhost`  创建到目标端口的连接。然后在 streaming server 的连接和该连接之间拷贝数据，完成数据的传递。

在 windows 上，是通过  `wincat.exe`  使用地址  `127.0.0.1`  创建到目标端口的连接。

### 参考源码

- pkg/cri/streaming/server.go:149\[9]
- pkg/cri/server/streaming.go:69\[10]
- pkg/cri/server/service.go:138\[11]
- pkg/cri/server/sandbox_portforward_linux.go:34\[12]

## 总结

结合源码分析对 port-foward 工作原理的梳理，相信对 cri 的工作方式也有了一定的了解。本文是以容器运行时 Containerd 为例，不同的容器运行时虽然实现了 cri，但是实现的细节上也会有所差异。

比如在 port-forward 的实现上，Kubernetes v1.23.0 版本中的 docker shim（1.24 中被移除）\[13]  中，是使用`nsenter`  进入 pod 所在的 network namespace 中通过  `socat`  完成的端口转发。

### 参考资料

\[1]staging/src/k8s.io/kubectl/pkg/cmd/portforward/portforward.go:389: [_https://github.com/kubernetes/kubernetes/tree/release-1.24/staging/src/k8s.io/kubectl/pkg/cmd/portforward/portforward.go#L389_](https://github.com/kubernetes/kubernetes/tree/release-1.24/staging/src/k8s.io/kubectl/pkg/cmd/portforward/portforward.go#L389)
\[2]staging/src/k8s.io/client-go/tools/portforward/portforward.go:242: [_https://github.com/kubernetes/kubernetes/tree/release-1.24/staging/src/k8s.io/client-go/tools/portforward/portforward.go#L242_](https://github.com/kubernetes/kubernetes/tree/release-1.24/staging/src/k8s.io/client-go/tools/portforward/portforward.go#L242)
\[3]staging/src/k8s.io/client-go/tools/portforward/portforward.go:330: [_https://github.com/kubernetes/kubernetes/tree/release-1.24/staging/src/k8s.io/client-go/tools/portforward/portforward.go#L330_](https://github.com/kubernetes/kubernetes/tree/release-1.24/staging/src/k8s.io/client-go/tools/portforward/portforward.go#L330)
\[4]pkg/registry/core/pod/rest/subresources.go:185: [_https://github.com/kubernetes/kubernetes/tree/release-1.24/pkg/registry/core/pod/rest/subresources.go#L185_](https://github.com/kubernetes/kubernetes/tree/release-1.24/pkg/registry/core/pod/rest/subresources.go#L185)
\[5]pkg/kubelet/server/server.go:463: [_https://github.com/kubernetes/kubernetes/tree/release-1.24/pkg/kubelet/server/server.go#L463_](https://github.com/kubernetes/kubernetes/tree/release-1.24/pkg/kubelet/server/server.go#L463)
\[6]pkg/kubelet/server/server.go:873: [_https://github.com/kubernetes/kubernetes/tree/release-1.24/pkg/kubelet/server/server.go#L873_](https://github.com/kubernetes/kubernetes/tree/release-1.24/pkg/kubelet/server/server.go#L873)
\[7]pkg/kubelet/cri/streaming/portforward/portforward.go:46: [_https://github.com/kubernetes/kubernetes/tree/release-1.24/pkg/kubelet/cri/streaming/portforward/portforward.go#L46_](https://github.com/kubernetes/kubernetes/tree/release-1.24/pkg/kubelet/cri/streaming/portforward/portforward.go#L46)
\[8]pkg/kubelet/cri/streaming/server.go:111: [_https://github.com/kubernetes/kubernetes/tree/release-1.24/pkg/kubelet/cri/streaming/server.go#L111_](https://github.com/kubernetes/kubernetes/tree/release-1.24/pkg/kubelet/cri/streaming/server.go#L111)
\[9]pkg/cri/streaming/server.go:149: [_https://github.com/containerd/containerd/tree/release/1.5/pkg/cri/streaming/server.go#L149_](https://github.com/containerd/containerd/tree/release/1.5/pkg/cri/streaming/server.go#L149)
\[10]pkg/cri/server/streaming.go:69: [_https://github.com/containerd/containerd/tree/release/1.5/pkg/cri/server/streaming.go#L69_](https://github.com/containerd/containerd/tree/release/1.5/pkg/cri/server/streaming.go#L69)
\[11]pkg/cri/server/service.go:138: [_https://github.com/containerd/containerd/tree/release/1.5/pkg/cri/server/service.go#L138_](https://github.com/containerd/containerd/tree/release/1.5/pkg/cri/server/service.go#L138)
\[12]pkg/cri/server/sandbox*portforward_linux.go:34: [\_https://github.com/containerd/containerd/tree/release/1.5/pkg/cri/server/sandbox_portforward_linux.go#L34*](https://github.com/containerd/containerd/tree/release/1.5/pkg/cri/server/sandbox_portforward_linux.go#L34)
\[13]Kubernetes v1.23.0 版本中的 docker shim（1.24 中被移除）: [_https://github.com/kubernetes/kubernetes/blob/release-1.23/pkg/kubelet/dockershim/docker_streaming_others.go#L43_](https://github.com/kubernetes/kubernetes/blob/release-1.23/pkg/kubelet/dockershim/docker_streaming_others.go#L43)
