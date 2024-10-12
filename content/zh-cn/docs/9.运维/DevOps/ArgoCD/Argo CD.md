---
title: Argo CD
---

# 概述

> 参考：
>
> - [GitHub 项目，argoproj/argo-cd](https://github.com/argoproj/argo-cd)
> - [Argo CD 保姆级入门教程](https://mp.weixin.qq.com/s/r1DnnHptOTaS_Gp8tpPWdg)

在上一篇『👉[GitOps 介绍\[1\]](http://mp.weixin.qq.com/s?__biz=MzU1MzY4NzQ1OA==&mid=2247509873&idx=1&sn=dd6daee66f39a965e4680ecd91f884d7&chksm=fbede7bccc9a6eaa6ddf5d082ed20a31d1956425eb8129428b2e2701fb216ef331a706f6d008&scene=21#wechat_redirect)』中，我介绍了什么是 GitOps，包括 GitOps 的原则和优势，以及 GitOps 与 DevOps 的区别。本文将介绍用于实施 GitOps 的工具 Argo CD。

Argo CD 是以 Kubernetes 作为基础设施，遵循声明式 GitOps 理念的持续交付（continuous delivery, CD）工具，支持多种配置管理工具，包括 ksonnet/jsonnet、kustomize 和 Helm 等。它的配置和使用非常简单，并且自带一个简单易用的可视化界面。

按照官方定义，Argo CD 被实现为一个 Kubernetes 控制器，它会持续监控正在运行的应用，并将当前的实际状态与 Git 仓库中声明的期望状态进行比较，如果实际状态不符合期望状态，就会更新应用的实际状态以匹配期望状态。

在正式开始解读和使用 Argo CD 之前，我们需要先搞清楚为什么需要 Argo CD？它能给我们带来什么价值？

## 传统 CD 工作流

从上篇文章『👉[GitOps 介绍\[2\]](http://mp.weixin.qq.com/s?__biz=MzU1MzY4NzQ1OA==&mid=2247509873&idx=1&sn=dd6daee66f39a965e4680ecd91f884d7&chksm=fbede7bccc9a6eaa6ddf5d082ed20a31d1956425eb8129428b2e2701fb216ef331a706f6d008&scene=21#wechat_redirect)』可以知道，目前大多数 CI/CD 工具都使用基于 Push 的部署模式，例如 Jenkins、CircleCI 等。这种模式一般都会在 CI 流水线运行完成后执行一个命令（比如 kubectl）将应用部署到目标环境中。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/argocd/1662076442944-8d47843e-5af5-4d31-ba22-3e8e6194bc88.jpeg)

这种 CD 模式的缺陷很明显：

- 需要安装配置额外工具（比如 kubectl）；
- 需要 Kubernetes 对其进行授权；
- 需要云平台授权；
- 无法感知部署状态。也就无法感知期望状态与实际状态的偏差，需要借助额外的方案来保障一致性。

下面以 Argo CD 为例，来看看遵循声明式 GitOps 理念的 CD 工具是怎么实现的。

## 使用 Argo CD 的 CD 工作流

和传统 CI/CD 工具一样，CI 部分并没有什么区别，无非就是测试、构建镜像、推送镜像、修改部署清单等等。重点在于 CD 部分。

Argo CD 会被部署在 Kubernetes 集群中，使用的是基于 Pull 的部署模式，它会周期性地监控应用的实际状态，也会周期性地拉取 Git 仓库中的配置清单，并将实际状态与期望状态进行比较，如果实际状态不符合期望状态，就会更新应用的实际状态以匹配期望状态。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/argocd/1662076443010-40f73633-ab15-4158-8e71-f792a2b8fe11.png)

无论是通过 CI 流水线触发更新 K8s 编排文件，还是 DevOps 工程师直接修改 K8s 编排文件，Argo CD 都会自动拉取最新的配置并应用到 K8s 集群中。

最终会得到一个相互隔离的 CI 与 CD 流水线，CI 流水线通常由研发人员（或者 DevOps 团队）控制，CD 流水线通常由集群管理员（或者 DevOps 团队）控制。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/argocd/1662076442866-9f57ee65-a081-468a-8013-cb6319c73a51.png)

## Argo CD 的优势

下面我们来看看 Argo CD 相较于传统 CD 工具有哪些比较明显的优势。

### Git 作为应用的唯一真实来源

所有 K8s 的声明式配置都保存在 Git 中，并把 Git 作为应用的唯一事实来源，我们不再需要手动更新应用（比如执行脚本，执行 kubectl apply 或者 helm install 命令），只需要通过统一的接口（Git）来更新应用。

此外，Argo CD 不仅会监控 Git 仓库中声明的期望状态，还会监控集群中应用的实际状态，并将两种状态进行对比，只要实际状态不符合期望状态，实际状态就会被修正与期望状态一致。所以即使有人修改了集群中应用的状态（比如修改了副本数量），Argo CD 还是会将其恢复到之前的状态。**这就真正确保了 Git 仓库中的编排文件可以作为集群状态的唯一真实来源。**

当然，有时候我们需要快速更新应用并进行调试，通过 Git 来触发更新还是慢了点，这也不是没有办法，我们可以修改 Argo CD 的配置，使其不对手动修改的部分进行覆盖或者回退，而是直接发送告警，提醒管理员不要忘了将更新提交到 Git 仓库中。

### 快速回滚

Argo CD 会定期拉取最新配置并应用到集群中，一旦最新的配置导致应用出现了故障（比如应用启动失败），我们可以通过 Git History 将应用状态快速恢复到上一个可用的状态。

如果你有多个 Kubernetes 集群使用同一个 Git 仓库，这个优势会更明显，因为你不需要分别在不同的集群中通过 `kubectl delete` 或者 `helm uninstall` 等手动方式进行回滚，只需要将 Git 仓库回滚到上一个可用的版本，Argo CD 便会自动同步。

### 集群灾备

如果你在青云\[3]北京 3 区中的 KubeSphere\[4] 集群出现故障，且短期内不可恢复，可以直接创建一个新集群，然后将 Argo CD 连接到 Git 仓库，这个仓库包含了整个集群的所有配置声明。最终新集群的状态会与之前旧集群的状态一致，完全不需要人工干预。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/argocd/1662076442877-2851ab43-7392-4cf9-8328-2403d8d2c586.png)

### 使用 Git 实现访问控制

通常在生产环境中是不允许所有人访问 Kubernetes 集群的，如果直接在 Kubernetes 集群中控制访问权限，必须要使用复杂的 RBAC 规则。在 Git 仓库中控制权限就比较简单了，例如所有人（DevOps 团队，运维团队，研发团队，等等）都可以向仓库中提交 Pull Request，但只有高级工程师可以合并 Pull Request。

这样做的好处是，除了集群管理员和少数人员之外，其他人不再需要直接访问 Kubernetes 集群，只需访问 Git 仓库即可。对于程序而言也是如此，类似于 Jenkins 这样的 CI 工具也不再需要访问 Kubernetes 的权限，因为只有 Argo CD 才可以 apply 配置清单，而且 Argo CD 已经部署在 Kubernetes 集群中，必要的访问权限已经配置妥当，这样就不需要给集群外的任意人或工具提供访问的证书，可以提供更强大的安全保障。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/argocd/1662076443389-772dc87a-8c2c-4ba9-aad2-310ec93fcb50.png)

### 扩展 Kubernetes

虽然 Argo CD 可以部署在 Kubernetes 集群中，享受 Kubernetes 带来的好处，但这不是 Argo CD 专属的呀！Jenkins 不是也可以部署在 Kubernetes 中吗？Argo CD 有啥特殊的吗？

那当然有了，没这金刚钻也不敢揽这瓷器活啊，Argo CD 巧妙地利用了 Kubernetes 集群中的很多功能来实现自己的目的，例如所有的资源都存储在 Etcd 集群中，利用 Kubernetes 的控制器来监控应用的实际状态并与期望状态进行对比，等等。

这样做最直观的好处就是**可以实时感知应用的部署状态**。例如，当你在 Git 仓库中更新配置清单中的镜像版本后，Argo CD 会将集群中的应用更新到最新版本，你可以在 Argo CD 的可视化界面中实时查看更新状态（比如 Pod 创建成功，应用成功运行并且处于健康状态，或者应用运行失败需要进行回滚操作）。

## Argo CD 架构

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/argocd/1662076443403-3ee90747-ee9b-44c4-9990-0ec25bff6e51.png)

从功能架构来看，Argo CD 主要有三个组件：API Server、Repository Server 和 Application Controller。从 GitOps 工作流的角度来看，总共分为 3 个阶段：检索、调谐和呈现。

### 检索 -- Repository Server

检索阶段会克隆应用声明式配置清单所在的 Git 仓库，并将其缓存到本地存储。包含 Kubernetes 原生的配置清单、Helm Chart 以及 Kustomize 配置清单。履行这些职责的组件就是 **Repository Server**。

### 调谐 -- Application Controller

调谐（Reconcile）阶段是最复杂的，这个阶段会将 **Repository Server** 获得的配置清单与反映集群当前状态的实时配置清单进行对比，一旦检测到应用处于 `OutOfSync` 状态，**Application Controller** 就会采取修正措施，使集群的实际状态与期望状态保持一致。

### 呈现 -- API Server

最后一个阶段是呈现阶段，由 Argo CD 的 **API Server** 负责，它本质上是一个 gRPC/REST Server，提供了一个无状态的可视化界面，用于展示调谐阶段的结果。同时还提供了以下这些功能：

- 应用管理和状态报告；
- 调用与应用相关的操作（例如同步、回滚、以及用户自定义的操作）；
- Git 仓库与集群凭证管理（以 Kubernetes Secret 的形式存储）；
- 为外部身份验证组件提供身份验证和授权委托；
- RBAC 增强；
- Git Webhook 事件的监听器/转发器。

## 部署 Argo CD

Argo CD 有两种不同的部署模式：

### 多租户

Argo CD 最常用的部署模式是多租户，一般如果组织内部包含多个应用研发团队，就会采用这种部署模式。用户可以使用可视化界面或者 argocd CLI 来访问 Argo CD。argocd CLI 必须先通过 `argocd login <server-host>` 来获取 Argo CD 的访问授权。

```c
$ argocd login SERVER [flags]

## Login to Argo CD using a username and password
$ argocd login cd.argoproj.io

## Login to Argo CD using SSO
$ argocd login cd.argoproj.io --sso

## Configure direct access using Kubernetes API server
$ argocd login cd.argoproj.io --core
```

多租户模式提供了两种不同的配置清单：

#### 非高可用

推荐用于测试和演示环境，不推荐在生产环境下使用。有两种部署清单可供选择：

- install.yaml\[5] - 标准的 Argo CD 部署清单，拥有集群管理员权限。可以使用 Argo CD 在其运行的集群内部署应用程序，也可以通过接入外部集群的凭证将应用部署到外部集群中。
- namespace-install.yaml\[6] - 这个部署清单只需要 namespace 级别的权限。如果你不需要在 Argo CD 运行的集群中部署应用，只需通过接入外部集群的凭证将应用部署到外部集群中，推荐使用此部署清单。还有一种花式玩法，你可以为每个团队分别部署单独的 Argo CD 实例，但是每个 Argo CD 实例都可以使用特殊的凭证（例如 `argocd cluster add <CONTEXT> --in-cluster --namespace <YOUR NAMESPACE>`）将应用部署到同一个集群中（即 `kubernetes.svc.default`，也就是内部集群）。

> ⚠️ 注意：namespace-install.yaml 配置清单中并不包含 Argo CD 的 CRD，需要自己提前单独部署：`kubectl apply -k https://github.com/argoproj/argo-cd/manifests/crds\?ref\=stable`。

#### 高可用

与非高可用部署清单包含的组件相同，但增强了高可用能力和弹性能力，推荐在生产环境中使用。

- ha/install.yaml\[7] - 与上文提到的 install.yaml 的内容相同，但配置了相关组件的多个副本。
- ha/namespace-install.yaml\[8] - 与上文提到的 namespace-install.yaml 相同，但配置了相关组件的多个副本。

### Core

Core 模式也就是最精简的部署模式，不包含 API Server 和可视化界面，只部署了每个组件的轻量级（非高可用）版本。

用户需要 Kubernetes 访问权限来管理 Argo CD，因此必须使用下面的命令来配置 argocd CLI：

`$ kubectl config set-context --current --namespace=argocd # change current kube context to argocd namespace $ argocd login --core`

也可以使用命令 `argocd admin dashboard` 手动启用可视化界面。

具体的配置清单位于 Git 仓库中的 core-install.yaml\[9]。

---

除了直接通过原生的配置清单进行部署，Argo CD 还支持额外的配置清单管理工具。

### Kustomize

Argo CD 配置清单也可以使用 Kustomize 来部署，建议通过远程的 URL 来调用配置清单，使用 patch 来配置自定义选项。

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: argocd
resources:
  - https://raw.githubusercontent.com/argoproj/argo-cd/v2.0.4/manifests/ha/install.yaml
```

### Helm

Argo CD 的 Helm Chart 目前由社区维护，地址：https://github.com/argoproj/argo-helm/tree/master/charts/argo-cd\[10]。
接下来开始部署 Argo CD：

```bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

查看部署结果：

```bash
$ kubectl -n argocd get pod
argocd-applicationset-controller-69879c47c-pcbkg   1/1     Running   0          26m
argocd-notifications-controller-6b4b74d8d8-s7mrz   1/1     Running   0          26m
argocd-redis-65596bf87-2hzcv                       1/1     Running   0          26m
argocd-dex-server-78c9764884-6lcww                 1/1     Running   0          26m
argocd-repo-server-657d46f8b-87rzq                 1/1     Running   0          26m
argocd-application-controller-0                    1/1     Running   0          26m
argocd-server-6b48df79dd-b7bkw                     1/1     Running   0          26m
```

## 访问 Argo CD

部署完成后，可以通过 Service `argocd-server` 来访问可视化界面。

```bash
$ kubectl -n argocd get svc
NAME                                      TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
argocd-applicationset-controller          ClusterIP   10.105.250.212   <none>        7000/TCP,8080/TCP            5m10s
argocd-dex-server                         ClusterIP   10.108.88.97     <none>        5556/TCP,5557/TCP,5558/TCP   5m10s
argocd-metrics                            ClusterIP   10.103.11.245    <none>        8082/TCP                     5m10s
argocd-notifications-controller-metrics   ClusterIP   10.98.136.200    <none>        9001/TCP                     5m9s
argocd-redis                              ClusterIP   10.110.151.108   <none>        6379/TCP                     5m9s
argocd-repo-server                        ClusterIP   10.109.131.197   <none>        8081/TCP,8084/TCP            5m9s
argocd-server                             ClusterIP   10.98.23.255     <none>        80/TCP,443/TCP               5m9s
argocd-server-metrics                     ClusterIP   10.103.184.121   <none>        8083/TCP                     5m8s
```

如果你的客户端可以直连 Service IP，那就直接可以通过 argocd-server 的 Cluster IP 来访问。或者可以直接通过本地端口转发来访问：

```bash

$ kubectl port-forward svc/argocd-server -n argocd 8080:443
Forwarding from 127.0.0.1:8080 -> 8080
Forwarding from [::1]:8080 -> 8080
```

初始密码以明文形式存储在 Secret `argocd-initial-admin-secret` 中，可以通过以下命令获取：

```bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo
```

也可以通过以下命令来修改登录密码：

```bash
argocd account update-password --account admin --current-password xxxx --new-password xxxx
```

登录后的界面：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/argocd/1662076443531-40c539e0-8797-4c46-b641-f433feeb83b8.png)

## Argo CD 核心概念

在正式开始使用 Argo CD 之前，需要先了解两个基本概念。

### Argo CD Application

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/argocd/1662076443778-4339e877-616e-464e-a6ef-b2f02ff400e6.png)

Argo CD 中的 Application 定义了 Kubernetes 资源的**来源**（Source）和**目标**（Destination）。来源指的是 Git 仓库中 Kubernetes 资源配置清单所在的位置，而目标是指资源在 Kubernetes 集群中的部署位置。

来源可以是原生的 Kubernetes 配置清单，也可以是 Helm Chart 或者 Kustomize 部署清单。

目标指定了 Kubernetes 集群中 API Server 的 URL 和相关的 namespace，这样 Argo CD 就知道将应用部署到哪个集群的哪个 namespace 中。

简而言之，**Application 的职责就是将目标 Kubernetes 集群中的 namespace 与 Git 仓库中声明的期望状态连接起来**。

Application 的配置清单示例：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/argocd/1662076443806-967865f5-cc3a-4f10-b687-16c0c479a069.jpeg)

如果有多个团队，每个团队都要维护大量的应用，就需要用到 Argo CD 的另一个概念：**项目**（Project）。

### Argo CD Project

Argo CD 中的项目（Project）可以用来对 Application 进行分组，不同的团队使用不同的项目，这样就实现了多租户环境。项目还支持更细粒度的访问权限控制：

- 限制部署内容（受信任的 Git 仓库）；
- 限制目标部署环境（目标集群和 namespace）；
- 限制部署的资源类型（例如 RBAC、CRD、DaemonSets、NetworkPolicy 等）；
- 定义项目角色，为 Application 提供 RBAC（与 OIDC group 或者 JWT 令牌绑定）。

## Demo 演示

最后通过一个简单的示例来展示 Argo CD 的工作流程。

### 准备 Git 仓库

在 GitHub 上创建一个项目，取名为 argocd-lab\[13]，为了方便实验将仓库设置为公共仓库。在仓库中新建 dev 目录，在目录中创建两个 YAML 配置清单，分别是 `deployment.yaml` 和 `service.yaml`。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/argocd/1662076443997-1e374e40-cdf9-49dc-b9a1-68269c22df40.png)
配置清单内容如下：

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  selector:
    matchLabels:
      app: myapp
  replicas: 2
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: myapp
        image: nginx:latest
        ports:
        - containerPort: 80

# service.yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
spec:
  selector:
    app: myapp
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
```

接下来在仓库根目录中创建一个 Application 的配置清单：

```yaml
# application.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: myapp-argo-application
  namespace: argocd
spec:
  project: default

source:
    repoURL: https://github.com/yangchuansheng/argocd-lab.git
    targetRevision: HEAD
    path: dev
  destination:
    server: https://kubernetes.default.svc
    namespace: myapp

syncPolicy:
    syncOptions:
    - CreateNamespace=true

automated:
      selfHeal: true
      prune: true
```

参数解释：

- **syncPolicy** : 指定自动同步策略和频率，不配置时需要手动触发同步。
- **syncOptions** : 定义同步方式。
- **CreateNamespace=true** : 如果不存在这个 namespace，就会自动创建它。
- **automated** : 检测到实际状态与期望状态不一致时，采取的同步措施。
- **selfHeal** : 当集群世纪状态不符合期望状态时，自动同步。
- **prune** : 自动同步时，删除 Git 中不存在的资源。

Argo CD 默认情况下**每 3 分钟**会检测 Git 仓库一次，用于判断应用实际状态是否和 Git 中声明的期望状态一致，如果不一致，状态就转换为 `OutOfSync`。默认情况下并不会触发更新，除非通过 `syncPolicy` 配置了自动同步。

如果嫌周期性同步太慢了，也可以通过设置 Webhook 来使 Git 仓库更新时立即触发同步。具体的使用方式会放到后续的教程中，本文不再赘述。

### 创建 Application

现在万事具备，只需要通过 application.yaml 创建 Application 即可。

```bash
$ kubectl apply -f application.yaml
application.argoproj.io/myapp-argo-application created
```

在 Argo CD 可视化界面中可以看到应用已经创建成功了。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/argocd/1662076444145-4526c86e-f5b5-43c3-8737-ec29479adfdf.png)
点进去可以看到应用的同步详情和各个资源的健康状况。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/argocd/1662076444129-228cbe84-33c0-4d98-ad77-d0e76e9292b5.png)
**如果你更新了 deployment.yaml 中的镜像，Argo CD 会自动检测到 Git 仓库中的更新，并且将集群中 Deployment 的镜像更新为 Git 仓库中最新设置的镜像版本。**

## 总结

本文介绍了 Argo CD 的优势、架构和工作原理，并通过一个简单的示例对其功能进行演示，比如修改 Git 仓库内容后，可以自动触发更新。还可以通过 Event Source 和 Trigger 实现更多自动化部署的需求。

在部署 Kubernetes 资源时，Argo CD 还支持 Kustomize、Helm、Ksonnet 等资源描述方式，包括其他更高级的使用方式都会在后续的教程中为大家一一道来，敬请期待。

### 引用链接

\[1]

GitOps 介绍: [_https://icloudnative.io/_](https://icloudnative.io/)

\[2]

GitOps 介绍: [_https://icloudnative.io/_](https://icloudnative.io/)

\[3]

青云: [_https://www.qingcloud.com/_](https://www.qingcloud.com/)

\[4]

KubeSphere: [_https://kubesphere.com.cn_](https://kubesphere.com.cn)

\[5]

install.yaml: [_https://github.com/argoproj/argo-cd/blob/master/manifests/install.yaml_](https://github.com/argoproj/argo-cd/blob/master/manifests/install.yaml)

\[6]

namespace-install.yaml: [_https://github.com/argoproj/argo-cd/blob/master/manifests/namespace-install.yaml_](https://github.com/argoproj/argo-cd/blob/master/manifests/namespace-install.yaml)

\[7]

ha/install.yaml: [_https://github.com/argoproj/argo-cd/blob/master/manifests/ha/install.yaml_](https://github.com/argoproj/argo-cd/blob/master/manifests/ha/install.yaml)

\[8]

ha/namespace-install.yaml: [_https://github.com/argoproj/argo-cd/blob/master/manifests/ha/namespace-install.yaml_](https://github.com/argoproj/argo-cd/blob/master/manifests/ha/namespace-install.yaml)

\[9]

core-install.yaml: [_https://github.com/argoproj/argo-cd/blob/master/manifests/core-install.yaml_](https://github.com/argoproj/argo-cd/blob/master/manifests/core-install.yaml)

\[10]

<https://github.com/argoproj/argo-helm/tree/master/charts/argo-cd:> [_https://github.com/argoproj/argo-helm/tree/master/charts/argo-cd_](https://github.com/argoproj/argo-helm/tree/master/charts/argo-cd)

\[11]

KubeSphere Cloud 托管集群服务: [_https://kubesphere.cloud/console/resource/_](https://kubesphere.cloud/console/resource/)

\[12]

<https://kubesphere.cloud>: [_https://kubesphere.cloud/_](https://kubesphere.cloud/)

\[13]

argocd-lab: [_https://github.com/yangchuansheng/argocd-lab_](https://github.com/yangchuansheng/argocd-lab)

\[14]

KubeSphere 从 3.3.0 开始: [_https://kubesphere.com.cn/news/kubesphere-3.3.0-ga-announcement/_](https://kubesphere.com.cn/news/kubesphere-3.3.0-ga-announcement/)
