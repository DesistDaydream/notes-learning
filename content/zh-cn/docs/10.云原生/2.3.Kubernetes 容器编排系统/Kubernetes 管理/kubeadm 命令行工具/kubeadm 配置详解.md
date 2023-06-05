---
title: kubeadm 配置详解
---

# 概述

> 参考：
>
> - [官方文档,参考-安装工具-kubeadm-kubeadmin init-结合配置文件使用 kubeadm init](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init/#config-file)
> - [官方文档,参考-配置 APIs-kubeadm 配置(v1beta3)](https://kubernetes.io/docs/reference/config-api/kubeadm-config.v1beta3/)
> - [kubeadm 库](https://pkg.go.dev/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm)
> - [v1beta2 版本的 kubeadm 包的配置文件字段详解](https://pkg.go.dev/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta2)

由于配置文件还在 beta 阶段，但是官方又推荐使用，所以很是纠结，我也不知道为啥。。。文档只有这种通过代码注释生成的内容~~~

kubeadm 库中的 [**Type**](https://pkg.go.dev/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm#pkg-types) 实际上就是配置文件中可用的字段，其实就是 Go 语言中 struct 与 yaml 的对应，配置文件都是 yaml 格式的，kind 中值，其实就是代码中的一个 struct。

kubeadm init 命令初始化集群时，对集群配置的首选方法是使用 `--config=FILE` 标志传递 YAML 格式的配置文件。

> kubeadm 配置文件中定义的某些配置选项在 kubeamd init 的命令行标志中也有对应标志，但是这些标志仅支持最常见/最简单的用例。

kubeadm 的配置文件可以看成是 kubeadm 几个资源的 Manifests 文件的集合。kubeadm 其中包括多种资源，一个文件中可以使用三个破折号 `---` 分隔的多种资源(其实就是 yaml 的语法)。每个资源就是一种配置类型。现阶段支持以下配置类型：

- **InitConfiguration** # 初始化集群配置
- **ClusterConfiguration** # 集群通用配置
- **KubeletConfiguration** # 覆盖 kubelet 运行时配置文件
- **KubeProxyConfiguration** # 覆盖 kube-proxy 运行时配置文件
- **JoinConfiguration** # 加入集群配置

# InitConfiguration Manifest 详解

参考：[**v1beta3 版本**](https://pkg.go.dev/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta3#InitConfiguration)
应该使用 InitConfiguration 类型来配置运行时设置，在 kubeadm init 情况下，是 bootstrap Tokens 的配置以及所有特定于执行 kubeadm 的节点的设置，包括以下几个字段

## apiVersion: kubeadm.k8s.io/v1beta3

## kind: InitConfiguration

## bootstrapTokens: <Object>

**groups: \[]STRING**

- system:bootstrappers:kubeadm:default-node-token

**ttl: TIME # 该 Token 的存活时间。**
0s 为永不过期。`默认值: 24h`

**usages: \[]STRING**

- signing
- authentication

## nodeRegistration: <Object>

其中包含与将新节点注册到集群有关的字段；使用它来自定义节点名称，要使用的 CRI 套接字或仅应应用于该节点的任何其他设置（例如，节点 ip）。
**name: STRING**
该字段的信息将会写入到 Node API 对象 的 .Metadata.Name 字段中。
此字段还用于 kubelet 到 API Server 的客户端证书的 CommonName 字段中。`默认值：节点主机名`。

**KubeletExtraArgs: map\[string]string # 通过额外的参数传递给 kubelet。**
此处的参数通过环境文件传递到 kubelet 命令行。kubeadm 在运行时将 kubelet 写入源。这将覆盖 kubelet-config-1.X ConfigMap 中的通用基本级别配置。解析时，标志具有更高的优先级。这些值是本地的，特定于正在执行 kubeadm 的节点。

**criSocket: STRING # kubelet 要使用的 runtime 的 Socket 文件的绝对路径**
CRISocket 用于检索容器运行时信息。此信息将注释到节点 API 对象，以便稍后重用

## LocalAPIEndpoint: <Object> # API Server 暴露的 IP 和 Port

该字段通常不用设置，直接设置 ClusterConfiguration 资源中的 controlPlaneEndpoint 字段即可。
**advertiseAddress: STRING** # API Server 暴露的 IP 地址
**bindPort: INT32** # API Server 暴露的安全端口。`默认值：6443`。

# ClusterConfiguration Manifest 详解

参考：[kubeadm 代码(v1beta2)](https://pkg.go.dev/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta2#ClusterConfiguration)
ClusterConfiguration 类型应用于配置群集范围的设置，包括以下设置：

- 网络，其中包含集群网络拓扑的配置；使用它例如定制节点子网或服务子网。
- Etcd 配置；使用它例如自定义本地 etcd 或配置 API 服务器以使用外部 etcd 集群。
- kube-apiserver，kube-scheduler，kube-controller-manager 程序的运行时配置；通过添加自定义设置或覆盖 kubeadm 默认设置，使用它来自定义控制平面组件。
  - 官方文档：<https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/control-plane-flags/>

## apiVersion: kubeadm.k8s.io/v1beta2

## kind: ClusterConfiguration

## etcd: <Object> # 集群中 etcd 配置

## networking: <Object> # 集群中网络拓扑的配置

**dnsDomain: STRING #**`默认值:cluster.local`。
**serviceSubnet: STRING #**`默认值:10.96.0.0/12`。

## controlPlaneEndpoint: STRING # 为控制平面设置一个 IP 或域名

`默认值：InitConfiguration 资源中 localAPIEndpoint.advertiseAddress + localAPIEndpoint.bindPort 两个字段的值`。
该字段就是设置访问 Kubernetes API 时，所要使用的端点，通过访问 Endpoint 就应该可以访问 [Kubernetes 的 API Server 程序](https://www.yuque.com/go/doc/33168516)。同时，各种与 API Server 交互时所用到的证书，其中也会包含该字段的值。

## apiServer: <Object> # 配置 apiserver 程序

**certSANs([]STRING)** # 为 API Server 的证书中的 Subject Alternative Name 字段设置额外的名称。
**extraArgs: <Object>** # 设定 apiserver 程序的命令行标志
**extravolumes([]Object)** # 设定 apiserver 程序的卷，以及挂载卷

## controllerManager: <Object> # 配置 controller-manager 程序

**extraArgs: <Object>** # 设定 controller-manager 程序的命令行标志
**extravolumes([]Object)** # 设定 controller-manager 程序的卷，以及挂载卷

## scheduler: <Object> # 配置 scheduler 程序

**extraArgs: <Object>** # 设定 scheduler 程序的命令行标志
**extravolumes([]Object)** # 设定 scheduler 程序的卷，以及挂载卷

## dns: <Object> # 配置 DNS 插件

## certificateDir: <STRING> # 指定 kubeadm 生成和读取证书的路径。`默认值：/etc/kubernetes/pki`

## imageRepository: STRING # 部署集群时拉取所需镜像的仓库。`默认值:k8s.gcr.io`

## clusterName: STRING # 集群的名称。`默认值:kubernetes`

# KubeletConfiguration Manifest 详解

> 参考：
>
> - [官方文档,入门-生产环境-使用部署工具安装 Kubernetes-使用 kubeadm 配置集群中的每个 kubelet](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/kubelet-integration/)
> - [官方文档,参考-配置 APIs-kubelet 配置(v1beta1)](https://kubernetes.io/docs/reference/config-api/kubelet-config.v1beta1/)
> - [kubelet 代码中 struct 与 yaml 字段对应(v1beta1)](https://pkg.go.dev/k8s.io/kubelet/config/v1beta1#KubeletConfiguration)
> - [kubelet 配置详解](/docs/10.云原生/2.3.Kubernetes%20 容器编排系统/2.Kubelet%20 节点代理/Kubelet%20 配置详解.md 节点代理/Kubelet 配置详解.md)

KubeletConfiguration 类型的配置中的字段，将会覆盖 kubelet 的配置文件(默认路径为 /var/lib/kubelet/config.yaml)中的字段

说白了，这些配置其实就跟直接修改 kubelet 运行时使用 --config 标志指定的文件是一样，就像下图一样：
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wm87kv/1618021089988-6fef9049-1249-4a85-b3f7-ce59968a35ec.png)

# KubeProxyConfiguration Manifest 详解

> 参考：
>
> - [kube-proxy 代码(v1alpha1)](https://pkg.go.dev/k8s.io/kube-proxy/config/v1alpha1#KubeProxyConfiguration)
> - [kube-proxy 配置详解](/docs/10.云原生/2.3.Kubernetes%20 容器编排系统/8.Kubernetes%20 网络/kube-proxy(实现%20Service%20 功能的组件).md 网络/kube-proxy(实现 Service 功能的组件).md)

与 KubeletConfiguration 类型配置一样，将会覆盖 kubeproxy 的配置。可以根据 [**kube-proxy 命令行工具官方文档**](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-proxy/) 或 [**kubeproxy 代码**](https://pkg.go.dev/k8s.io/kube-proxy/config/v1alpha1#KubeProxyConfiguration) 参考这个类型配置应如何配置

# JoinConfiguration Manifest 详解

# 配置示例

在[**这里**](https://pkg.go.dev/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta2#hdr-Basics)有 [**v1beta2 版本的 kubeadm 包**](https://pkg.go.dev/k8s.io/kubernetes@v1.19.4/cmd/kubeadm/app/apis/kubeadm/v1beta2) 的 kubeadm-config.yaml 文件的完整配置，这个示例包含了所有可用字段。

> 在 `Here is a fully populated example of a single YAML file containing multiple configuration types to be used during a`kubeadm init`run.` 这段描述中

## 最简单的配置

```yaml
cat > kubeadm-config.yaml << EOF
apiVersion: kubeadm.k8s.io/v1beta2
kind: InitConfiguration
bootstrapTokens:
- groups:
  - system:bootstrappers:kubeadm:default-node-token
  ttl: 0s
  usages:
  - signing
  - authentication
# 可选，当使用其他 runtime 时，在此指定，这里使用了 containerd
# nodeRegistration:
#   criSocket: /run/containerd/containerd.sock
---
apiVersion: kubeadm.k8s.io/v1beta2
kind: ClusterConfiguration
kubernetesVersion: v1.19.4
controlPlaneEndpoint: "X.X.X.X:6443"
imageRepository: registry.aliyuncs.com/k8sxio
networking:
  podSubnet: 10.244.0.0/16
---
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
mode: "ipvs"
---
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
cgroupDriver: systemd
EOF
```

## 复杂配置

```yaml
apiVersion: kubeadm.k8s.io/v1beta2
kind: InitConfiguration
bootstrapTokens:
  - groups:
      - system:bootstrappers:kubeadm:default-node-token
    ttl: 0s
    usages:
      - signing
      - authentication
---
apiVersion: kubeadm.k8s.io/v1beta2
kind: ClusterConfiguration
kubernetesVersion: v1.19.2
controlPlaneEndpoint: k8s-api.bj-test.desistdaydream.ltd:6443
imageRepository: registry.aliyuncs.com/k8sxio
networking:
  podSubnet: 10.244.0.0/16
  serviceSubnet: 10.96.0.0/12
etcd:
  local:
    extraArgs:
      listen-metrics-urls: http://0.0.0.0:2381
apiServer:
  certSANs:
    - localhost
    - 127.0.0.1
    - k8s-api.bj-test.desistdaydream.ltd
    - 172.19.42.234
    - master-3.bj-test
    - 172.19.42.233
    - master-2.bj-test
    - 172.19.42.232
    - master-1.bj-test
    - 172.19.42.231
  extraArgs:
    service-node-port-range: 30000-60000
  extraVolumes:
    - name: host-time
      hostPath: /etc/localtime
      mountPath: /etc/localtime
      readOnly: true
controllerManager:
  extraArgs:
    bind-address: 0.0.0.0
  extraVolumes:
    - name: host-time
      hostPath: /etc/localtime
      mountPath: /etc/localtime
      readOnly: true
scheduler:
  extraArgs:
    bind-address: 0.0.0.0
  extraVolumes:
    - name: host-time
      hostPath: /etc/localtime
      mountPath: /etc/localtime
      readOnly: true
---
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
mode: "ipvs"
---
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
cgroupDriver: systemd
kubeReserved:
  cpu: 200m
  memory: 250Mi
systemReserved:
  cpu: 200m
  memory: 250Mi
evictionHard:
  memory.available: 5%
evictionSoft:
  memory.available: 10%
evictionSoftGracePeriod:
  memory.available: 2m
```
