---
title: Kubelet
linkTitle: Kubelet
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，参考-组件工具-kubelet](https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/)

Kubelet 是在每个节点上运行的主要 **节点代理**。它可以使用以下之一向 APIServer 注册节点：用于覆盖主机名的标志；或云提供商的特定逻辑。

kubelet 根据 PodSpec 起作用。 PodSpec 是一个描述 Pod 的 YAML 或 JSON 对象。 kubelet 接受通过各种机制（主要是通过 apiserver）提供的一组 PodSpec，并确保这些 PodSpec 中描述的容器正在运行且运行状况良好。 Kubelet 不管理不是 Kubernetes 创建的容器。一般情况， PodSpec 都是由在 k8s 对象的 yaml 文件中定义的。

kubelet 负责维护容器(CNI)的生命周期，同时也负责 Volume（CVI）和 Network（CNI）的管理。kubernetes 集群的宿主机上，启动的每一个 pod 都有由 kubelet 这个组件管理的。

kubelet 在每个 Node 上都会启动一个 kubelet daemon 进程，默认监听在 **10250** 端口。该进程用于处理 Master 节点(主要是 apiserver)下发到本节点的任务，管理 Pod 以及 Pod 中的容器。每个 kubelet 进程会在 APIServer 上注册节点自身信息，定期向 Master 节点汇报节点资源的使用情况，并通过 cAdvisor(kubelet 内部功能) 监控容器和节点资源。10248 为 kubelet 健康检查的 healthz 端口

> Note: 如果 master 节点不运行 Pod 的话，是不用部署 kubelet 的。

kubelet 使用 PodSpec 来对其所在节点的 Pod 进行管理，PodSpec(Pod Specification) 是描述 pod 的 yaml 或者 json 对象（PodSpec 一般是 yaml 或者 json 格式的文本文件）。这些 PodSpecs 有多个个来源

- **apiserver** # 使用最多的方式，通过 kubectl 命令向 apiserver 提交 PodSpec 文件，然后 apiserver 再下发给相应节点的 node。还有一个通过 APIServer 监听 etcd 目录，同步 PodSpec
- **File** # kubelet 定期监控某个路径(默认路径为/etc/kubernetes/manifests)下所有文件，把这些文件当做 PodSpec，这种方式也就是所谓的 [StaticPod(静态 Pod)](Static Pod)。默认情况下每 20 秒监控一下，可以通过 flag 进行配置，配置时可以指定具体的路径以及监控周期
- **HTTP endpoint(URL)** # 使用--manifest-url 参数，让 kubelet 每 20 秒检查一次 URL 指定的 endpoint(端点)
- **HTTP server** # kubelet 监听 HTTP 请求，并响应简单的 API 以提交新的 Pod 清单

## Static Pod

所有以非 API Server 方式创建的 Pod 都叫 Static Pod(静态 Pod)。

kubelet 的工作核心，就是一个控制循环。驱动这个控制循环运行的实践，包括四种

1. Pod 更新事件
2. Pod 生命周期变化
3. kubelet 本身设置的执行周期
4. 定时的清理事件

注意：kubelet 调用下层容器运行时的执行过程，并不会直接调用 Docker 的 API，而是通过一组叫作 CRI（Container Runtime Interface，容器运行时接口）的 gRPC 接口来间接执行的。gRPC 接口规范详见官网：<https://grpc.io/docs/>

kubelet 是集群的基础设施，其他主要组件如果不以 daemon 形式运行，则依赖 kubelet 以 pod 方式启动，这样，才可以组成集群最基本的形态。

## Kubelet Metrics

详见：k8s 主要组件 metrics 获取指南

# Kubelet 部署

```bash
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64/
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF
yum install -y kubelet
```

# Kubelet 关联文件与配置

kubelet 可以通过多个地方读取其自身的配置并更改自己的行为方式，可以通过指定的 yaml 格式的文件读取配置信息，也可以直接指定命令行参数传递到 kubelet 程序中。

**/var/lib/kubelet/** # kubelet 配置文件目录、以及运行时数据目录，包含基础配置文件、证书、通过 kubelet 启动的容器信息等等

- **./config.yaml** # kubelet 基础配置文件。一般在 kubelet 启动时使用 --cofnig 参数指定该读取该文件的路径进行加载。
  - Note：该文件内容与 kubectl get configmap -n kube-system kubelet-config-X.XX -o yaml 命令所得结果一样
    - 如果想要在 kubelet 运行时动态得更改其配置，则可以修改 configmap 中的内容，详见：<https://kubernetes.io/docs/tasks/administer-cluster/reconfigure-kubelet/>
- **./kubeadm-flags.env** # 该文件将内容作为 kubelet 参数，在 kubelet 启动时加载，常用来在 kubeadm 初始化时使用
- **./pods/** # kubelet 启动的 Pod 的数据保存路径，其内目录名为 Pod 的 uid 。
  - ./${POD_UID}/volumes/ # 对应 pod 挂载的 volume 保存路径，其内目录为 `kubernetes.io~TYPE/` ，其中 TYPE 为 volume 的类型。
- **./pki/** # kubelet 与 apiserver 交互时所用到的证书存放目录。
  - ./kubelet.crt # 在 kubelet 完成 TLS bootstrapping 后并且没有配置 --feature-gates=RotateKubeletServerCertificate=true 时生成；这种情况下该文件为一个独立于 apiserver CA 的自签 CA 证书，有效期为 1 年；被用作 kubelet 10250 api 端口。当其他东西需要访问 kubelet 的 api 时，需要使用该证书作为认证。
  - ./kubelet-client-current.pem # 与 API server 通讯所用到的证书，与 apiserver 交互后生成。

**/etc/kubernetes/** # Kubernetes 系统组件运行时目录。

- **./manifests/** # Kubelet 默认从该目录中读取 Pod 的 Manifests 文件，以运行静态类型 Pod。
- kubelet 在 k8s 集群交互时认证文件所在目录，kubelet 需要读取认证配置，用来与 apiserver 进行交互。
  - **./bootstrap-kubelet.conf** # 用于 TLS 引导程序的 KubeConfig 文件。该 kubeconfig 文件的用户信息为 token。该文件用于 kubelet 所在节点不在集群中时，向集群发起注册请求所用，如果节点已在集群中，则会自动生成 kubelet.conf 文件
  - **./kubelet.conf** # 具有唯一 kubelet 标识的 KubeConfig 文件(与 kubectl 的 config 文件一样，用于 kubelet 与 apiserver 交互时提供认证信息)。该 kubeconfig 文件的用户信息为客户端证书和私钥，一般在 kubelet 启动时由 bootstrap-kubelet.conf 文件生成。
    - 当该文件不存在时，会在 kubelet 启动时，由 bootstrap-kubelet.confg 生成

**/etc/sysconfig/kubelet** # 与 /var/lib/kubelet/kubeadm-flags.env 文件作用一样 ，将内容作为 kubelet 参数，在 kubelet 启动时加载。一般用于让用户指定 kubelet 的运行时参数 。 KUBELET_EXTRA_ARGS 在标志链中排在最后，并且在其他设置冲突时具有最高优先级。

- Note：对于 DEB 系统，配置文件位于：/etc/default/kubelet

**/usr/lib/systemd/system/kubelet.service.d/10-kubeadm.confg** # 与 kubelet 守护进程运行参数

# Kubelet 的启动过程 && Kubelet 与 APIServer 的交互说明

## kubelet 启动过程

> 参考：
>
> - kubelet 启动流程源码分析
>   - https://xiaohanliang.gitbook.io/notes/k8s/zhi-shi-yu-ding-yi/jie-dian-zu-jian-kubelet/kubelet-qi-dong-liu-cheng
>   - https://www.jianshu.com/p/e07d84cce9f9

1. 读取配置 kubelet 配置文件。kubelet 启动时首先会根据 `--config=PATH` 参数指定路径(默认 /var/lib/kubelet/config.yaml)读取配置文件，并根据其内配置加载 kubelet 相关参数，
   1. 根据 ca 配置路径加载 ca.crt 文件，并在 /var/lib/kubelet/pki 目录下生成关于 kubelet 的 10250 私有 api 端口所需的 crt 与 key 文件。
   2. 如果该文件不存在或有问题，则启动失败。
2. 配置与 apiserver 通信的 kubeconfig 文件。根据 --bootstrap-kubeconfig=PATH 参数加载 /etc/kubernetes/bootstrap-kubelet.conf 文件(如果不存在则根据 --kubeconfig=PATH 参数加载 /etc/kubernetes/kubelet.conf 文件)，两个文件都不存在则报错
   1. 如果当前节点不在集群中，则会执行证书申请操作
      1. 首先 kubelet 向 bootstrap-kubelet.conf 文件内配置的 apiserver(文件中 server 的配置) 发送加入集群的申请申，同时 kubelet 会根据 bootstrap-kubelet.conf 文件生成 kubelet.conf 文件，将该 kubeconfig 文件中的 user 认证方式改为证书认证方式，并指定证书路径为/var/lib/kubelet/pki/kubelet-client-current.pem。
      2. 在集群 master 上执行 kubectl get csr 命令获取当前申请列表，并使用 kubectl certificate approve XXXX 命令通过该节点的申请
      3. master 节点的 controller-manager 处理完该请求后，本地 kubelet 将生成 kubelet.conf 文件所需证书，并保存在 /var/lib/kubelet/pki/ 目录下，以 /var/lib/kubelet/pki/kubelet-client-current.pem-TIME 命名，并建立软件链接指向该文件。
   2. 如果当前节点已在集群中，则会根据 bootstrap-kubelet.conf 文件生成 kubelet.conf 文件，且根据 kubelet.conf 文件中的证书信息与 apiserver 进行通信。
   3. 如果当前节点已在集群且存在 kubelet.conf 文件，则使用该 kubeconfig 文件与 apiserver 进行交互后执行后续操作
3. 初始化 kubelet 组件内部模块
   1. 会在 /var/lib/kubelet 目录下添加相关文件，用以驱动 pod 及 kubelet 相关功能。
   2. 检查 cgroup 驱动设置与 CRI 的 cgroup 驱动设置是否一致，如果不一致则启动失败
   3. 等等操作，后续有发现再补充
4. 启动 kubelet 内部模块及服务(Note：当 csr 申请发送后，即可执行后续工作，无需等待申请审批通过)

## kubelet 启动后的工作原理

kubelet 的工作核心就是在围绕着不同的生产者生产出来的不同的有关 pod 的消息来调用相应的消费者（不同的子模块）完成不同的行为(创建和删除 pod 等)，即图中的控制循环（SyncLoop），通过不同的事件驱动这个控制循环运行。如下图所示：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kubernetes/kubelet/1616120074647-12f084c2-e91c-475d-93d4-f4e3fd3e61ef.png)

### kubelet 创建 pod 流程

参考：

<https://www.jianshu.com/p/5e0c9d1dbe95>
<https://www.kubernetes.org.cn/6766.html>

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kubernetes/kubelet/1616120074625-d25d5fad-d58b-4142-958d-b1aee16d8e71.png)

Note：注意 14，,15 步，kubelet 会先将生成配置(volume 挂载、配置主机名等等)，才会去启动 pod，哪怕 pod 启动失败，挂载依然存在。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kubernetes/kubelet/1616120074636-5f1f460b-6cf2-4190-ab01-6284178805e6.png)

# Kubelet 所管理三大板块

## kubelet 负责所在节点 Container(CRI) 的管理

CRI 的起源：

官方介绍：<https://kubernetes.io/blog/2016/12/container-runtime-interface-cri-in-kubernetes/>

Kubernetes 项目之所以要在 kubelet 中引入这样一层单独的抽象，是为了对 Kubernetes 屏蔽下层容器运行时的差异。因为 Kubernets 是一个编排工具，不只可以编排 Docker，还可以编排除 Docker 以外的其余的容器项目。而每种容器项目的实现方式不尽相同，为了解决这个问题，那么可以把 kubelet 对容器的操作，统一抽象成一个借口，这样，kubelet 就只需要跟这个借口打交道，而作为具体的容器项目，它们只需要自己提供一个该接口的实现，然后对 kubelet 暴露出 gRPC 服务即可(docker shim 现在集成在了 kubelet 中，以后会单独拿出来甚至废弃)

CRI 的运作方式：

- 当 kubernetes 通过编排能力声明一个 Pod 后，调度器会为这个 pod 选择一个具体的 Node 来运行，这时候，该 Node 上的 kubelet 会通过 SyncLoop 判断需要执行的具体操作，这个时候 kubelet 会调用一个叫做 GenericRuntime 的通用组件来发起创建 Pod 的 CRI 请求。
- CRI shim(CRI 垫片，宿主机与 kubelet 之间的东西)来响应 CRI 请求，然后把请求“翻译”成对后端容器项目的请求或者操作。
- 每个容器项目都会自己实现一个 CRI shim，然后 CRI shim 收到的请求会转给对应的容器守护进程(e.g.docker 项目里的 dockerd)，由该守护进程进行容器的创建、更改、删除、exec 等操作

当前主流的 CRI 有如下几种：

- Docker(kubelet 默认的 CRI)
  - Note：kubelet 内置 dockershim，在启动或，会生成 dockersim.sock 文件，kubelet 与 crictl 都会默认与该文件关联
- CRI-O
- Containerd(通常与 docker 同时安装)
- Frakti(kata OCI 的实现)

kubelet 与 CRI 对接的方式

- kubelet 根据参数 --container-runtime-endpoint 来决定其所绑定的 CRI sock。默认使用 docker 的 sock，路径为：/var/run/dockershim.sock
  - 可以通过 crictl 工具来测试目标 sock 是否可用，crictl 用法详见：crictl 命令行工具
- 如果使用不同的 CRI 运行时，则需要为 kubelet 指定不同的 flag。 例如，当使用非 docker CRI 时， 则需要使用 --container-runtime=remote 与 --container-runtime-path-endpoint=\<PATH> 指定 CRI 端点。endpoint 的值为指定 CRI 的 sock 文件。
- 各个 CRI 需要进行配置才可与 kubelet 对接成功，如果不进行初始化配置，则 kubelet 无法获取到该 CRI 对于 k8s 的相关配置参数

kubelet 通过如下两个命令行标志来指定要使用的 CRI

```bash
--container-runtime=remote
--container-runtime-endpoint=unix:///run/containerd/containerd.sock
```

## kubelet 负责所在节点 Network(CNI) 的管理

kubelet 对 cni(Container Network Interface 容器网络接口)的调用详见另一片关于网络介绍的文章：[Kubernetes 网络](/docs/10.云原生/Kubernetes/Kubernetes%20网络/Kubernetes%20网络.md)

kubelet 配置 pod 网络时，首先会读取下 /etc/cni/net.d/_目录下的配置，查看当前所使用的 CNI 插件及插件参数，比如现在是 flannel ，那么 flannel 会将 /run/flannel/subnet.env 文件的配置信息传递给 kubelet ，然后 kubelet 使用 /opt/cni/bin/_ 目录中的二进制文件，来处理处理 pod 的网络信息。

## kubelet 负责所在节点 Volume(CVI) 的管理
