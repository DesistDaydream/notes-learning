---
title: CNI
---

# 概述

> 参考：
> - [GitHub 项目](https://github.com/containernetworking/cni)
> - [GitHub,containernetworking-cni-规范](https://github.com/containernetworking/cni/blob/master/SPEC.md)

CNI 与 [OCI](✏IT 学习笔记/☁️10.云原生/2.1.容器/Open%20Containers%20Initiative(开放容器倡议).md Containers Initiative(开放容器倡议).md) 是类似的东西，都是一种规范。

**Container Network Interface(容器网络接口，简称 CNI)** 是一个 CNCF 项目，用于编写为 Linux 容器配置网络接口的插件。CNI 由两部分组成：

- CNI Specification(规范)
- CNI Libraries(库)

由于 CNI 仅仅关注在容器的网络连接以及在删除容器时移出通过 CNI 分配的网络资源。所以，CNI 具有广泛的支持，并且该规范易于实现。

## CNI Specification(规范)

每个 CNI 插件必须由 二进制文件 来实现，且这些文件应该可以被容器管理系统(比如 Kubernetes)调用。

CNI 插件负责将网络接口插入容器网络名称空间(例如 veth 对的一端)中，并在主机上进行任何必要的更改(例如将 veth 的另一端连接到网桥)。然后通过调用适当的 IPAM 插件，将 IP 分配给接口并设置与 IP 地址管理部分一致的路由。

## CNI Libraries(库)

任何程序都可以调用 CNI 库来实现容器网络，比如 [nerdctl](https://github.com/containerd/nerdctl)、kubelet 等

# CNI 的部署和使用方式

> 官方文档：<https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/network-plugins/#installation>

CNI 规范与编程语言无关，并且 CNI 自身仅仅维护标准配置文件和基础插件，想要使用 CNI 来实现容器网络，只需根据标准，调用 CNI 库，即可在程序中实现(比如 nerdctl、kubelet 等)。这些通过 CNI 库实现了容器网络的程序，通过 **CNI 插件 **为其所启动的容器，创建关联网络。

就拿 kubelet 举例,<https://github.com/kubernetes/kubernetes/blob/release-1.22/cmd/kubelet/app/options/container_runtime.go>

```go
const (
	// When these values are updated, also update test/utils/image/manifest.go
	defaultPodSandboxImageName    = "k8s.gcr.io/pause"
	defaultPodSandboxImageVersion = "3.5"
)

var (
	defaultPodSandboxImage = defaultPodSandboxImageName +
		":" + defaultPodSandboxImageVersion
)

// NewContainerRuntimeOptions will create a new ContainerRuntimeOptions with
// default values.
func NewContainerRuntimeOptions() *config.ContainerRuntimeOptions {
	dockerEndpoint := ""
	if runtime.GOOS != "windows" {
		dockerEndpoint = "unix:///var/run/docker.sock"
	}

	return &config.ContainerRuntimeOptions{
		ContainerRuntime:          kubetypes.DockerContainerRuntime,
		DockerEndpoint:            dockerEndpoint,
		DockershimRootDirectory:   "/var/lib/dockershim",
		PodSandboxImage:           defaultPodSandboxImage,
		ImagePullProgressDeadline: metav1.Duration{Duration: 1 * time.Minute},

		CNIBinDir:   "/opt/cni/bin",
		CNIConfDir:  "/etc/cni/net.d",
		CNICacheDir: "/var/lib/cni/cache",
	}
}
```

可以看到，kubelet 在启动时通过 /etc/cni/net.d/ 目录下的配置文件来加载网络插件。当启动一个 Pod 时，kubelet 调用该目录下的网络插件配置后，由网络插件代为给 Pod 地址分配，接口创建，网络创建等

在部署 Kubernetes 的 CNI 插件时，有一个步骤是安装 kubernetes-cni 包，其目的就是在宿主机上安装 CNI 插件所需的基础二进制文件。这些文件一般保存在 /opt/cni/bin/ 目录中。

kubelet 配置 pod 网络时，首先会读取下 /etc/cni/net.d/_ 目录下的配置，查看当前所使用的 CNI 插件及插件参数，比如现在是 flannel ，那么 flannel 会将 /run/flannel/subnet.env 文件的配置信息传递给 kubelet ，然后 kubelet 使用 /opt/cni/bin/_ 目录中的二进制文件，来处理处理 pod 的网络信息。

注意：各种 CNI 的 cidr 配置由 controller-manager 维护，`--cluster-cidr=10.244.0.0/16` 与 `--node-cidr-mask-size=24` 这俩参数用来指定 cidr 的范围。

同时 CNI 还可以被 nerdctl 工具使用，作为直接使用 Containerd 启动容器的网络接口，让容器附加在通过 CNI 的 plugin 创建出来的网络设备上。nerdctl 同样会读取 CNI 配置文件，并通过 CNI 插件创建网络设备之后，nerdctl 再将容器关联到网络设备上。

## CNI 插件说的通俗易懂点，其实就是两个主要功能

1. 路由表维护
2. 发现路由规则
3. 生成路由表
4. 流量处理(如果需要的话)

# CNI 插件列表

> 官方文档：<https://github.com/containernetworking/plugins>

下面仅列出由 CNI 团队维护的一些参考和示例插件。CNI 的基础可执行文件一般分为三类：

第一类：Main 插件，用于创建具体的网络设备。
比如，bridge、ipvlan、loopback、macvlan、ptp(Veth Pair)、vlan 等。都属于“网桥”类型的 CNI 插件，所以在具体实现中，往往会调用 bridge 这个二进制文件。

- **bridge** # 创建一个桥设备，向其中添加主机和容器。
- **ipvlan** # Adds an ipvlan interface in the container.
- **loopback** # Set the state of loopback interface to up.
- **macvlan** # Creates a new MAC address, forwards all traffic to that to the container.
- **ptp** # 创建一对 veth 设备
- **vlan** # Allocates a vlan device.
- **host-device** # Move an already-existing device into a container.

第二类：IPAM 插件(IP Address Management)，用于负责分配 IP 地址。
比如，dchp、host-local

- dhcp：这个文件会向 DHCP 服务器发起请求；
- host-local，会使用预先配置的 IP 地址段来进行分配

第三类：其他插件
比如 flannel、tuning、portmap、bandwidth

- flannel：专门为 Flannel 项目提供的 CNI 插件。早期的默认插件，叠加网络，不支持网络策略(即定义哪个 Pod 访问哪个 Pod 等策略)
- tuning：是一个通过 sysctl 调整网络设备参数的二进制文件
- portmap：是一个通过 iptables 配置端口映射的二进制文件
- bandwidth：是一个使用 Token Bucket Filter(TBF)来进行限流的二进制文件

### 第三方 plugin

> 可用的第三方插件列表在这里：<https://github.com/containernetworking/cni#3rd-party-plugins>

下面是一些常用的第三方 CNI 简介

- **calico** # 三层隧道网络，基于 BGP 协议，即支持网络配置也支持网络策略
- **Cilium - BPF & XDP for containers **# 基于 eBPF 实现的，性能很好
- 等

## 各种 CNI 的对比

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/afl2qt/1616118670278-ee95f4dd-7834-4fbc-82e2-e8b9ddddcf33.jpeg)

# CNI 关联文件

**/etc/cni/net.d/\* **# 默认配置文件保存目录
**/opt/cni/bin/\*** # 默认 CNI 插件保存目录
**/var/lib/cni/\*** # 默认 CNI 运行时产生的数据目录
