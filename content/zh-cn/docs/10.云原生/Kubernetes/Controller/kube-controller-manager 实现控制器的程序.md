---
title: kube-controller-manager 实现控制器的程序
linkTitle: kube-controller-manager 实现控制器的程序
weight: 2
---

# 概述

> 参考：
>
> -

kube-controller-manager 是实现 Kubernetes [Controller(控制器)](/docs/10.云原生/Kubernetes/Controller/Controller.md) 的程序。一般在集群启动之前，由 kubelet 使用静态 Pod 以容器方式运行；或者使用 systemd 以 daemon 方式运行。

kube-controller-manager 启动后监听两个端口。

1. 10257 端口是需要身份验证和授权的 https 服务端口。
2. 10252 为不安全的 http 服务端口。

## kube-controller-manager 高可用

> 参考：[Leader Election(领导人选举)](/docs/10.云原生/Kubernetes/Kubernetes%20机制与特性/Leader%20Election(领导人选举).md)

我们都知道 k8s 核心组件，其中 apiserver 只用于接收 api 请求，不会主动进行各种动作，所以他们在每个节点都运行并且都可以接收请求，不会造成异常；kube-proxy 也是一样，只用于做端口转发，不会主动进行动作执行。
但是 scheduler, controller-manager 不同，他们参与了 Pod 的调度及具体的各种资源的管控，如果同时有多个 controller-manager 来对 Pod 资源进行调度，结果太美不敢看，那么 k8s 是如何做到正确运转的呢？
k8s 所有功能都是通过 `services` 对外暴露接口，而 `services` 对应的是具体的 `endpoints` ，那么来看下 scheduler 和 controller-manager 的 `endpoints` 是什么：

```bash
[root@node70 21:04:46 ~]$kubectl -n kube-system describe endpoints kube-scheduler
Name:         kube-scheduler
Namespace:    kube-system
Labels:       <none>
Annotations:  control-plane.alpha.kubernetes.io/leader:
                {"holderIdentity":"node70_ed12bf09-7aa3-47d6-9546-97752bb589b5","leaseDurationSeconds":15,"acquireTime":"2019-09-11T05:31:58Z","renewTime"...
Subsets:
Events:  <none>
[root@node70 21:05:25 ~]$kubectl -n kube-system describe endpoints kube-controller-manager
Name:         kube-controller-manager
Namespace:    kube-system
Labels:       <none>
Annotations:  control-plane.alpha.kubernetes.io/leader:
                {"holderIdentity":"node71_c8deeaea-2d66-4459-90ee-65c28563062f","leaseDurationSeconds":15,"acquireTime":"2019-09-12T12:44:15Z","renewTime"...
Subsets:
Events:
  Type    Reason          Age   From                     Message
  ----    ------          ----  ----                     -------
  Normal  LeaderElection  22m   kube-controller-manager  node71_c8deeaea-2d66-4459-90ee-65c28563062f became leader
```

可以看到关键字 `[control-plane.alpha.kubernetes.io/leader](http://control-plane.alpha.kubernetes.io/leader)` ，这两个组件是通过 leader 选举来从集群中多个节点选择一个执行具体动作

其实，kube-controller-manager 默认使用 leases 资源来实现领导者选举功能：

```bash
~]# kubectl get leases.coordination.k8s.io -n kube-system
NAME                      HOLDER                                                 AGE
kube-controller-manager   master-2.bj-net_62b724de-66a3-4aff-9a7c-e7c3d66555d1   176d
kube-scheduler            master-2.bj-net_50df0a21-f59a-48de-98a5-93ab4a0ddf3b   176d
~]# kubectl get leases.coordination.k8s.io -n kube-system kube-controller-manager -oyaml  | neat
apiVersion: coordination.k8s.io/v1
kind: Lease
metadata:
  name: kube-controller-manager
  namespace: kube-system
spec:
  acquireTime: "2021-03-28T13:22:03.000000Z"
  holderIdentity: master-2.bj-net_62b724de-66a3-4aff-9a7c-e7c3d66555d1
  leaseDurationSeconds: 15
  leaseTransitions: 9
  renewTime: "2021-04-07T08:36:21.251656Z"
```

这里的关键字段是 holderIdentity，可以看到，现在 master-2 上的 kube-controller-manager 是领导者

如果我们去看 `/etc/kubernetes/manifests/` 下的配置文件，会看到这行配置：

```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    component: kube-controller-manager
    tier: control-plane
  name: kube-controller-manager
  namespace: kube-system
spec:
  containers:
    - command:
        - kube-controller-manager
---
- --leader-elect=true
```

通过在 YAML 中添加 `--leader-elect=true` 来决定是否进行选主逻辑。而这个参数也是在执行 `kubeadm` 部署集群时就自动配置好了，无需手动配置。

## kube-controller-manager 监控指标

详见：[kubernetes 监控](/docs/10.云原生/2.3.Kubernetes%20 容器编排系统/Kubernetes%20 管理/Kubernetes%20 监控/Kubernetes%20 系统组件指标.md 管理/Kubernetes 监控/Kubernetes 系统组件指标.md)

# kube-controller-manager 配置

> 参考：
>
> - [官方文档,参考-组件工具-kube-controller-manager](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-controller-manager/)

kube-controller-manager 主要通过命令行标志来控制运行时行为

- **--bind-address**(STRING) # 指定监听在 --secure-port 参数设定的端口的 IP。集群的其余部分以及 CLI 或者 web 客户端必须可以访问关联的接口。
- **--cluster-cidr**(STRING) # 集群中 Pod 的 CIDR 范围。
  - 用白话说：可以给 Pod 分配的 IP 范围。
- **--controllers**(STRING) # 要启动的控制器列表。`默认值：*`。
  - `*` 表示默认的控制器，比如 deployment、etc. 。
- **--leader-elect**(BOOLEAN) # 在程序开始循环监控之前，是否要启用领导者选举功能。`默认值：true`。
  - 若集群中有多个 kube-controller-manager，则必须要启用领导者选举功能。
- **--leader-elect-resource-lock**(STRING) # 在领导者选举期间用于获取锁的资源。`默认值：leases`。
  - 其他可以获取锁的资源有 endpoints 和 configmaps。
- **--node-cidr-mask-size**(INT32) # 集群中 Node 的 CIDR 的掩码。`IPv4 的默认值：24`。`IPv6的默认值：64`。
  - 用白话说：与 --cluster-cidr 标志互相配合，以确定 CNI 可以为每个节点分配的 IP 范围。该标志的值必须要大于 --cluster-cidr 中掩码的值。
- **--node-monitor-grace-period**(DURATION) # 将一个 Running 状态节点标记为不健康(NotReady、Unkonw 等)状态之前，允许节点处于不健康状态的时间上限。
- **--node-monitor-period**(DURATION) # 节点控制器同步节点状态的周期。`默认值：5s`。
  - kube-controller-manager 每隔 --node-monitor-period 时间就会去检查所有节点上 kubelet 的状态。如果持续 --node-monitor-grace-period 时间之后，被检查的节点依然不健康，则会将该节点标记为不健康
  - 并且，当节点被标记为不健康时，所有节点的 Pod 都会从 Endpoint 中踢出。
- **--node-startup-grace-period**(DURATION) # 节点启动期间可以处于无响应状态，但是超出 --node-startup-grace-period 时间后依然无响应，则将节点标记为不健康(NotRead、Unknow 等)状态。
- **--pod-eviction-timeout**(DURATION) # 节点被标记为不健康状态(NotReady、Unkonw 等)后，等待 DURATION 时间后驱逐故障节点上所有 Pod。`默认值：5m0s`
  - 这个标志的用法，可以通过[官方文档, 概念 - 集群架构 - 节点 章节中节点状态](https://kubernetes.io/docs/concepts/architecture/nodes/#condition)小节获得更详细的说明。
  - 其实，由于节点不可用，kubelet 无法接收到消息，说是删除 Pod，其实故障节点上 Pod 只会一直处于 Terminating 状态，因为故障节点的 kubelet 不可用，无法真正完成删除操作。
  - 但是在驱逐之前，如果节点状态不健康，则 service 管理的 endpoint 中，所有属于该节点的 Pod 都会被踢出，防止异常节点上的 Pod 处理请求。
- **--secure-port**(INT) # 指定通过身份验证和授权为 HTTPS 服务的端口。`默认值：10257`。
- **--use-service-account-credentials**(BOOLEAN) # 是否为每个控制器使用单独的 service account。`默认值：无`。

## 默认的 manifest 示例

```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    component: kube-controller-manager
    tier: control-plane
  name: kube-controller-manager
  namespace: kube-system
spec:
  containers:
    - command:
        - kube-controller-manager
        - --allocate-node-cidrs=true
        - --authentication-kubeconfig=/etc/kubernetes/controller-manager.conf
        - --authorization-kubeconfig=/etc/kubernetes/controller-manager.conf
        - --bind-address=0.0.0.0
        - --client-ca-file=/etc/kubernetes/pki/ca.crt
        - --cluster-cidr=10.244.0.0/16
        - --cluster-name=kubernetes
        - --cluster-signing-cert-file=/etc/kubernetes/pki/ca.crt
        - --cluster-signing-key-file=/etc/kubernetes/pki/ca.key
        - --controllers=*,bootstrapsigner,tokencleaner
        - --kubeconfig=/etc/kubernetes/controller-manager.conf
        - --leader-elect=true
        - --node-cidr-mask-size=24
        - --port=10252
        - --requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.crt
        - --root-ca-file=/etc/kubernetes/pki/ca.crt
        - --service-account-private-key-file=/etc/kubernetes/pki/sa.key
        - --service-cluster-ip-range=10.96.0.0/12
        - --use-service-account-credentials=true
      image: registry.aliyuncs.com/k8sxio/kube-controller-manager:v1.19.2
      imagePullPolicy: IfNotPresent
      livenessProbe:
        failureThreshold: 8
        httpGet:
          path: /healthz
          port: 10257
          scheme: HTTPS
        initialDelaySeconds: 10
        periodSeconds: 10
        timeoutSeconds: 15
      name: kube-controller-manager
      resources:
        requests:
          cpu: 200m
      startupProbe:
        failureThreshold: 24
        httpGet:
          path: /healthz
          port: 10257
          scheme: HTTPS
        initialDelaySeconds: 10
        periodSeconds: 10
        timeoutSeconds: 15
      volumeMounts:
        - mountPath: /etc/ssl/certs
          name: ca-certs
          readOnly: true
        - mountPath: /etc/pki
          name: etc-pki
          readOnly: true
        - mountPath: /usr/libexec/kubernetes/kubelet-plugins/volume/exec
          name: flexvolume-dir
        - mountPath: /etc/localtime
          name: host-time
          readOnly: true
        - mountPath: /etc/kubernetes/pki
          name: k8s-certs
          readOnly: true
        - mountPath: /etc/kubernetes/controller-manager.conf
          name: kubeconfig
          readOnly: true
  hostNetwork: true
  priorityClassName: system-node-critical
  volumes:
    - hostPath:
        path: /etc/ssl/certs
        type: DirectoryOrCreate
      name: ca-certs
    - hostPath:
        path: /etc/pki
        type: DirectoryOrCreate
      name: etc-pki
    - hostPath:
        path: /usr/libexec/kubernetes/kubelet-plugins/volume/exec
        type: DirectoryOrCreate
      name: flexvolume-dir
    - hostPath:
        path: /etc/localtime
        type: ""
      name: host-time
    - hostPath:
        path: /etc/kubernetes/pki
        type: DirectoryOrCreate
      name: k8s-certs
    - hostPath:
        path: /etc/kubernetes/controller-manager.conf
        type: FileOrCreate
      name: kubeconfig
```
