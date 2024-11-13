---
title: Kubelet 配置
linkTitle: Kubelet 配置
date: 2019-11-04T09:09:00
weight: 2
---

# 概述

> 参考：
>
> - [官方文档，入门-生产环境-使用工具安装 Kubernetes-使用 kubeadm 引导集群-使用 kubeadm 配置集群中每个 kubelet](https://kubernetes.io/zh/docs/setup/production-environment/tools/kubeadm/kubelet-integration/)
> - [官方文档，参考-配置 APIs-Kubelet 配置(v1beta1)](https://kubernetes.io/docs/reference/config-api/kubelet-config.v1beta1/)

可以通过两种方式配置 kubelet 运行时行为

1. **config.yaml 配置文件** # config.yaml 文件默认路径为 /var/lib/kubelet/config.yaml ，可以通过 --config \<FILE> 来指定其他的文件。
   1. [这里](https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/)是官方文档对于配置文件的概述。在[章节中间部分](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/kubelet/config/v1beta1/types.go)，可以直接看到配置文件对应的代码中结构体，也就是配置文件详细内容
   2. [这里](https://kubernetes.io/docs/reference/config-api/kubelet-config.v1beta1/)是配置文件中每个字段的详解，与代码中的结构体互相对应，只不过是整理后，可以直接在网页上查看，更清晰。
2. **kubelet 命令行标志**
   1. [这里](https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/)是官方文档对命令行标志的详解

官方更推荐使用第一种方式，通过 config.yaml 的文件修改，来改变 kubelet 的运行时参数。

很多配置文件的内容与命令行标志具有一一对应的关系，比如：

| 配置文件                                                     | 命令行标志                                  |
| ------------------------------------------------------------ | ------------------------------------------- |
| cgroupDriver: systemd                                        | --cgroup-driver=systemd                     |
| clusterDNS: \[10.96.0.10,...]                                | --cluster-dns=10.96.0.10,...                |
| authentication.x509.clientCAFile: /etc/kubernetes/pki/ca.crt | --client-ca-file=/etc/kubernetes/pki/ca.crt |
| 等等                                                         | 等等                                        |

但是也有一些是没有对应关系的，只能通过配置文件，或者命令行标志配置。比如命令行标志的 `--container-runtime` 就无法在配置文件中配置。在命令行标志官方文档中，凡是标着 `DEPRECATED` 的命令行标志，都是可以在配置文件中配置的。

# 命令行标志详解

> 参考：
>
> - [官方文档,参考-组件工具-kubelet](https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/)

**--cni-conf-dir \<STRING>** # Warning：Alpha 功能。指定 STRING 目录中搜索 CNI 配置文件。 `默认值：/etc/cni/net.d`

- 仅当 CRI 为 docker 时此标志才有效

**--config=\<STRING>** # 加载配置文件的路径。kubelet 将从该标志指定的文件中加载其初始配置。

**--pod-infra-container-image \<STRINIG>** # 指定在每个 pod 中将会使用的 network/ipc 名称空间的基础容器。`默认值：k8s.gcr.io/pause:3.1`

- 这个就是用来指定 infra 基础设施容器。

**--container-runtime \<STRING>** # kubelet 要使用的容器运行时，也就是要对接的 CRI。`默认值：docker`。

- remote # 表示使用其他运行时。需要配合 `--container-runtime-endpoint` 标志一起使用。

**--container-runtime-endpoint \<STRING>** # kubelet 要使用的运行时的路径。`默认值：unix:///var/run/dockershim.sock`

- STRING 是 socket 路径，现阶段只支持 UNIX sock，后面还可以支持远程，比如通过 http 来连接运行时。

**--image-service-endpoint \<STRING>** # kubelet 处理镜像所用的后端路径。若未指定，则于 `--container-runtime-endpoint` 标志的值相同

# 配置文件详解

> 参考：
>
> - [官方文档,任务-通过配置文件设置 kubelet 参数](https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/)
> - [官方文档,参考-配置 APIs-Kubelet 配置](https://kubernetes.io/docs/reference/config-api/kubelet-config.v1beta1/)
> - [代码](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/kubelet/config/v1beta1/types.go) # 该代码是 Go 结构体与 JSON 格式的解析对应关系。其中的注释就是配置文件各字段的含义

**apiVersion:** kubelet.config.k8s.io/v1beta1

**kind:** KubeletConfiguration

**address(STRING)** # kubelet 服务的 IP。默认为 0.0.0.0

**cgroupDriver(cgroupfs|systemd)** # kubelet 用于操纵主机上 cgroup 的驱动程序。`默认值：cgroupfs`

**imageMinimumGCAge(DURATION)** # 未使用的 image 进行垃圾回收之前的最小期限。`默认值：2m`

**nodeStatusReportFrequency(DURATION)** # 节点状态报告频率。`默认值：10s`

**nodeStatusUpdateFrequency(DURATION)** # 节点状态更新频率。`默认值：5m`

**resolvConf(STRING)** # kubelet 启动的容器所使用的解析器的配置文件。`默认值：/etc/resolv.conf`

- Ubuntu 中，配置则会被改为 `/run/systemd/resolve/resolv.conf`

## 配置文件示例

下面是 kubelet 1.18.8 版本的基本示例，其中通过 kubeadm 传递了 cgroupDriver 的值。

```yaml
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
authentication:
  anonymous:
    enabled: false
  webhook:
    cacheTTL: 0s
    enabled: true
  x509:
    clientCAFile: /etc/kubernetes/pki/ca.crt
authorization:
  mode: Webhook
  webhook:
    cacheAuthorizedTTL: 0s
    cacheUnauthorizedTTL: 0s
cgroupDriver: systemd #kubelet 用于操纵主机上 cgroup 的驱动程序。  (默认为 cgroupfs )
clusterDNS:
  - 10.96.0.10
clusterDomain: cluster.local
cpuManagerReconcilePeriod: 0s
evictionPressureTransitionPeriod: 0s
fileCheckFrequency: 0s
healthzBindAddress: 127.0.0.1
healthzPort: 10248
httpCheckFrequency: 0s
imageMinimumGCAge: 0s
nodeStatusReportFrequency: 0s
nodeStatusUpdateFrequency: 0s
rotateCertificates: true
runtimeRequestTimeout: 0s
staticPodPath: /etc/kubernetes/manifests
streamingConnectionIdleTimeout: 0s
syncFrequency: 0s
volumeStatsAggPeriod: 0s
```
