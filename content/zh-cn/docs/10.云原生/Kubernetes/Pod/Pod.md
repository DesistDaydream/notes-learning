---
title: Pod
linkTitle: Pod
date: 2024-04-02T10:32
weight: 1
---

# 概述

> 参考：
>
> - [官方文档,概念-工作负载-Pods](https://kubernetes.io/docs/concepts/workloads/pods/)

Pod 是 Kubernetes 集群内**最小的工作单元**，是一个逻辑概念。Kubernetes 真正处理的，还是通过 CRI 在 HostOS 上的 Namespace 和 Cgroups。所谓的 Pod 只是一组共享了某些资源的 Container，这一组 Container 共享同一个 NetworkNamespace 并且可以声明共享同一个 Volume。

**Infrastructure(基础设施，简称 Infra) 容器**：为了保证多个 Container 在共享的时候是对等关系(一般情况可以先启动 ContainerA，再启动 ContainerB 并共享 ContainerA 的资源，但是这样 A 与 B 不对等，A 是必须先启动才能启动 B)，需要一个中间 Container，即 **Infra 容器**，Infra 容器 永远是第一个被创建的 Container，想要共享某些资源的 Container 则通过加入 NetworkNamespce 的方式，与 Infra 容器 关联在一起。效果如图

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/iogldt/1616119861463-06b2877d-519d-43a9-a6e4-6fc743d6ee30.jpeg)

- Infra 类型的 Container 使用一个名为 pause 的镜像，就像它的名字一样，永远处于"暂停"状态
- Kubernetes 为每个 Pod 都附属了 k8s.gcr.ip/pause，这个 Container 只接管 Pod 的网络信息，业务 Container 通过加入网络 Container 的网络来实现网络共享。此容器随着 pod 创建而创建，随着 Pod 删除而删除。该容器是对业务 pod 的命名空间的解析。Note：如果想要更改该容器，则需要在 kubelet 中使用--pod-infra-container-image 参数进行配置
- 与 Infra 关联的 Container 的所有 NetworkNamespace 必然是完全一样的。
- 该链接有一种详细的解释
- Note：对于 kubelet 来说，这种容器称为 Sandbox。每次 kubelet 创建 pod 时，首先创建的也是 sandbox(i.e.pause)容器

一组 Container 共享 Infra 的 NetworkNamespace 意味着：

- 它们可以直接使用 localhost 进行通信
- 它们看到的网络设备跟 Infra 容器中看到的完全一样
- 一个 Pod 只能有有一个 IP 地址，就是这个 Pod 的 NetworkNamespace 对应的 IP 地址
- Pod 的生命周期只跟 Infra 容器一致，同与 Infra 关联的所有 Container 无关
- Pod 中的所有 Container 的进出流量都是通过 Infra 容器完成的，所以网络插件不必关心除 Infra 以外的容器的启动与否，只需关注如何配置 Pod(也就是 Infra 容器的 NetworkNamespace)即可

每个 Pod 包含一个或多个容器。Pod 中的容器会作为一个整体被 Master 调度到一个 Node 上运行。

如果把 Pod 想象成一台"服务器"，把 Container 想象成运行在这台服务器中的"用户程序"

- 凡是调度、网络、存储、以及安全相关的字段，基本都是 Pod 级别的，比如：
- 配置这台"服务器"的网卡(Pod 的网络)、配置这台“服务器”的磁盘(Pod 的存储，Volume)、配置这台”服务器“的防火墙(Pod 中的安全)、配置这台”服务器“运行在哪个机房(Pod 的调度)
- 凡是资源配额、所要使用的 port、探测该进程是否存活或就绪、需要使用"服务器"上的哪块 Volume 等等字段，都是 Container 级别的

## Kubernetes 引入 Pod 主要基于下面两个目的

- 可管理性。
  - 有些 Container 天生就是需要紧密联系，一起工作。Pod 提供了比容器更高层次的抽象，将它们封装到一个部署单元中。Kubernetes 以 Pod 为最小单位进行调度、扩展、共享资源、管理生命周期。
- 通信和资源共享。
  - Pod 中的所有 Container 使用同一个网络 namespace，即相同的 IP 地址和 Port 空间。它们可以直接用 localhost 通信。同样的，这些 Container 可以共享存储，当 Kubernetes 挂载 volume 到 Pod，本质上是将 volume 挂载到 Pod 中的每一个 Container。user,mnt,pnt。
- 使用户从传统虚拟机环境向容器环境迁移更加平滑，可以把 Pod 想象成 VM，Pod 中的 Container 是 VM 中的进程，甚至可以启动一个 systemd 的 Container
- 还可以把 Pod 理解为传统环境中的物理主机

## Container 设计模式

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/iogldt/1616119861478-cf678269-344a-4932-8448-c9eee14a8438.png)

- **sidecar** #(该英文的解释“跨斗”：一辆小而低的车辆，安装在摩托车旁边，用于载客，就像右图中的样子)，所以该模式就类似于这个，指可以再一个 Pod 中启动一个辅助 Container，来完成一些独立于主进程(主 Container)之外的工作。
  - 比如 Container 的日志收集：现在有一个 APP，需要不断把日志文件输出到 Container 的/var/log 目录中。这时，把一个 Pod 里的 Volume 挂载到应用 Container 的/var/log 目录上，然后在 Pod 中同时运行一个 sidecar 的 Container 也声明挂载同一个 Volume 到自己的/var/log 目录上，然后 sidecar 只需要不断得从自己的/var/log 目录里读取日志文件，转发到 MongoDB 或者 Elasticsearch 中存储起来即可。
  - Istio 项目也是使用 sidecar 模式的 Container 完成微服务治理的原理。

## Pod 的类型

- 动态 Pod：由 k8s 管理，网络组件，监控，等等，这些在 使用 kubeadm 初始化集群后才创建的 Pod 为动态 Pod
- 静态 Pod：由 kubelet 直接管理的，在 /etc/kubernetes/manifest/ 目录下的 yaml 文件

# Pod 使用方式

> 参考：
>
> - [官方文档，概念-工作敷在-Pod-初始化容器](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/)
> - [官方文档，概念-工作负载-Pod-临时容器](https://kubernetes.io/docs/concepts/workloads/pods/ephemeral-containers/)

- 运行单一 Container。
  - one-container-per-Pod 是 Kubernetes 最常见的模型，这种情况下，只是将单个 Container 简单封装成 Pod。
- 运行多个 Container。
  - 这些 Container 联系必须非常紧密，而且需要直接共享资源的应该放到一个 Pod 中(注意：当使用多 Container 的时候，其中一个 Container 要加上 command 的参数，否则其中一个起不来。因为 container 如果不执行某些命令，则启动后会自动结束，详见 docker 说明 1.LXC 与 Docker 入门最佳实践.note 里《Dokcer 的工作模式》章节)
  - 比如：File Puller 会定期从外部的 Content Manager 中拉取最新的文件，将其存放在共享的 volume 中。Web Server 从 volume 读取文件，响应 Consumer 的请求。这两个容器是紧密协作的，它们一起为 Consumer 提供最新的数据；同时它们也通过 volume 共享数据。所以放到一个 Pod 是合适的。

在 Pod 中，可运行的容器分为三类：

- **ephemeral_container(临时容器)** # 与 1.23 版本进入 beta，用来调试集群
- **init_container(初始化容器)** # 在应用容器启动前运行一次就结束的，常用来为容器运行初始化运行环境，比如设置权限等等
- **application_container(应用容器)** # 真正运行业务的容器。

这三类容器，可以在 kubelet 代码中找到运行逻辑，详见 [《kubelet 源码解析-PodWorker 模块》](/docs/10.云原生/2.3.Kubernetes%20 容器编排系统/Kubernetes%20 开发/源码解析/Kubelet%20 源码/PodWorker%20 模块.md 开发/源码解析/Kubelet 源码/PodWorker 模块.md)

## ephemeral_container(临时容器)

## init_container(初始化容器)

Pod 能够具有多个容器，应用运行在容器里面，但是它也可能有一个或多个先于应用容器启动的 Init 容器。Init 容器在所有容器运行之前执行（run-to-completion），常用来初始化配置。

如果为一个 Pod 指定了多个 Init 容器，那些容器会按顺序一次运行一个。 每个 Init 容器必须运行成功，下一个才能够运行。 当所有的 Init 容器运行完成时，Kubernetes 初始化 Pod 并像平常一样运行应用容器。

下面是一个 Init 容器的示例：

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: init-demo
spec:
  containers:
    - name: nginx
      image: nginx
      ports:
        - containerPort: 80
      volumeMounts:
        - name: workdir
          mountPath: /usr/share/nginx/html
  # These containers are run during pod initialization
  initContainers:
    - name: install
      image: busybox
      command:
        - wget
        - "-O"
        - "/work-dir/index.html"
        - http://kubernetes.io
      volumeMounts:
        - name: workdir
          mountPath: "/work-dir"
  dnsPolicy: Default
  volumes:
    - name: workdir
      emptyDir: {}
```

因为 Init 容器具有与应用容器分离的单独镜像，使用 init 容器启动相关代码具有如下优势：

- 它们可以包含并运行实用工具，出于安全考虑，是不建议在应用容器镜像中包含这些实用工具的。
- 它们可以包含使用工具和定制化代码来安装，但是不能出现在应用镜像中。例如，创建镜像没必要 FROM 另一个镜像，只需要在安装过程中使用类似 sed、 awk、 python 或 dig 这样的工具。
- 应用镜像可以分离出创建和部署的角色，而没有必要联合它们构建一个单独的镜像。
- 它们使用 Linux Namespace，所以对应用容器具有不同的文件系统视图。因此，它们能够具有访问 Secret 的权限，而应用容器不能够访问。
- 它们在应用容器启动之前运行完成，然而应用容器并行运行，所以 Init 容器提供了一种简单的方式来阻塞或延迟应用容器的启动，直到满足了一组先决条件。

Init 容器的资源计算，选择一下两者的较大值：

- 所有 Init 容器中的资源使用的最大值
- Pod 中所有容器资源使用的总和

Init 容器的重启策略：

- 如果 Init 容器执行失败，Pod 设置的 restartPolicy 为 Never，则 pod 将处于 fail 状态。否则 Pod 将一直重新执行每一个 Init 容器直到所有的 Init 容器都成功。
- 如果 Pod 异常退出，重新拉取 Pod 后，Init 容器也会被重新执行。所以在 Init 容器中执行的任务，需要保证是幂等的。

## container(容器) # 也称为 application_container(应用容器)

# Pod 名字的命名规范

一般情况都不会直接使用 Pod，而是通过 [Controller(控制器)](/docs/10.云原生/Kubernetes/Controller/Controller.md) 来创建。通过 Controller 创建一个 POD 的流程为，以及 POD 名字的命名方式每个对象的命名方式是：子对象的名字 = 父对象名字 + 随机字符串或数字。如图所示

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/iogldt/1616119861434-f2c4735b-c549-40a4-ab70-67217755ed3f.png)

- 用户通过 kubectl 创建 Deployment。
- Deployment 创建 ReplicaSet。
- ReplicaSet 创建 Pod。
