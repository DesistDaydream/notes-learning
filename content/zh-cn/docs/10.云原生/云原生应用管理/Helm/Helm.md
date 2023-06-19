---
title: Helm
---

# 概述

> 参考：
> - [GitHub 项目，helm/helm](https://github.com/helm/helm)
> - [官方文档](https://helm.sh/docs/)
> - 其他后期发现的文章
>   - <https://www.cnblogs.com/liugp/p/16659802.html>

**Helm** 是 Kubernetes 的 **Package Manager(包管理器)**。Kubernetes 在希腊语中，意为舵手或飞行员，是一个蓝色的舵轮图标。所以 Helm 就以类似的概念命名，Helm 称为舵柄，图标与 Kubernetes 类似，寓意把握着 Kubernetes 航行的方向。

## 主要概念

> 参考：
> - [官方文档，介绍-使用 Helm](https://helm.sh/docs/intro/using_helm/)

Helm 与 Kubernetes 的关系，就好比 yum 与 RedHat，apt 与 Ubuntu 一样，是一个 Kubernetes 专用的包管理器，安装专用于 k8s 集群之上的软件包。Helm 使用 Chart 帮助我们管理应用，Chart 就像 RPM 一样，里面描述了应用及其依赖关系。

**Chart(图表)** 是由 Helm 管理的应用部署包。Chart 是在一个结构相对固定的目录中，包含用于描述一个应用的一组 manifests 文件。

1. Chart Archive(图标归档) # 是一个将 Chart 打包成 .tgz 格式的压缩文件。
2. 实际上，Chart 就是很多 [**manifests**](https://kubernetes.io/docs/reference/glossary/?all=true#term-manifest) 的集合，里面有一个应用程序所需的 yaml 文件。而对于 kubernetes 来说，所谓的应用程序(软件包)也就是一堆 manifests，每个 manifest 代表一种资源(比如 deployment、service、ingress、configmap 等等)，这些 manifests 组合起来，就构成了一个应用。
3. **Chart 包** 就像 **RPM 包**一样。这不过没有类似 rpm 的命令，而是直接使用类似 yum 的 helm 命令来管理这些包。并且，Chart 包 也和 RPM 包一样，具有依赖关系。

**Release(发布)** 是将 Chart 部署到 Kubernets 集群中运行的实例，每一次 helm install CHART 就会生成一个 Release

1. Chart 是基础包，通过 config 赋值，生成 release。一般 config 来自于 chart archive 包中的 values.yaml 文件

chart 与 release 的关系就像 windows 中的 .exe 的安装文件与安装完成后在添加删除程序中看到的应用程序。chart 就是安装文件，release 就是程序。只不过 chart 可以是压缩包或者文件夹或者 url；helm list 命令就相当于打开了添加删除程序，可以看到已经安装好的 release。

**Repository(仓库)** 是存放 Charts 的地方，就是类似于 yum 源的概念。Helm 添加一个仓库 URL，就可以查看或安装该仓库下的 Charts。

基于上述概念，可以这么描述他们：Helm 安装 Charts 到 kubernetes 中，并为每个安装创建一个新的 Release。如果想要找到新的 Charts，可以使用 Helm 搜索 Repository

### [**在安装前自定义 Chart**](https://helm.sh/docs/intro/using_helm/#customizing-the-chart-before-installing)

在我们使用 yum 安装应用时，一般都是安装完成后，再对配置文件进行编辑，以改变应用的运行行为。

但是使用 Helm 则不能这么做，因为 Helm 直接将应用部署到 Kubernetes 集群中，部署完成后，再改变运行行为的方式是不优雅，且不方便的。所以，我们需要 **Customizing the Chart Before Installing(在安装前自定义图表)**。

要查看图表上可配置的选项，使用`helm show values`命令。

## Charts Repository(图表仓库)

通常来说，我们可以直接使用别人已经做好的 Chart，就跟使用 Docker 镜像，或者使用 yum 安装一样。并且，Helm 社区在早期已经维护了一个 [**Helm Charts Hub**](https://github.com/helm/charts)，这个 Hub 里包含丰富的 Charts。随着云原生应用的发展，这个仓库需要处理的 PR 越来越多，维护非常困难，所以 Helm 逐步把 Charts Hub 中的内容移动到 [Artifact Hub](https://artifacthub.io/) 中(Charts Hub 维护期持续 1 年)。[**Arifact Hub**](https://artifacthub.io/) 是一个基于 Web 的应用程序，可用于查找，安装和发布 CNCF 项目的软件包和配置。例如，这可能包括 Helm 图表，Falco 配置，开放策略代理（OPA）策略和 OLM 运算符。

Artifact Hub 中，除了包含 Helm Charts Hub 中的各种资源外，还有各种开源软件官方维护的 Charts，以及 [**Bitnami 中适用于 kubernetes 的 Charts**](https://bitnami.com/stacks/helm)。

[**Bitnami**](https://bitnami.com/) 使我们可以轻松地在任何平台上启动并运行您我们喜爱的开源软件，包括笔记本电脑，Kubernetes 和所有主要云。除了流行的社区产品之外，Bitnami 现在是 VMware 的一部分，它为 IT 组织提供安全，合规，连续维护和可根据组织策略自定义的企业产品。

# Helm 的安装

注意：Helm 使用时，会读取 /root/.kube/config 文件来连接 Kubernetes 集群。

- 下载 [helm 的 linux 压缩文件](https://github.com/helm/helm/releases)，解压后把 helm 的二进制文件移动到 /usr/bin 目录下，即可直接使用 helm 命令
- 配置命令补全
  - echo "source <(helm completion bash)" >> /root/.bashrc
- 安装 push 插件
  - helm plugin install https://github.com/chartmuseum/helm-push.git
- 当创建一个 Release 的时候，会先把 Chart 的数据存档(.tgz 格式)文件下载到 Helm 配置目录的 archive 目录下，然后再安装

# Helm 关联文件与配置

**~/.cache/helm/** # helm 缓存路径

- **./plugin/** # helm 插件安装路径

**/root/.local/share/helm/plugins/** #

# Helm 安装资源的顺序

- Namespace
- NetworkPolicy
- ResourceQuota
- LimitRange
- PodSecurityPolicy
- PodDisruptionBudget
- ServiceAccount
- Secret
- SecretList
- ConfigMap
- StorageClass
- PersistentVolume
- PersistentVolumeClaim
- CustomResourceDefinition
- ClusterRole
- ClusterRoleList
- ClusterRoleBinding
- ClusterRoleBindingList
- Role
- RoleList
- RoleBinding
- RoleBindingList
- Service
- DaemonSet
- Pod
- ReplicationController
- ReplicaSet
- Deployment
- HorizontalPodAutoscaler
- StatefulSet
- Job
- CronJob
- Ingress
- APIService
