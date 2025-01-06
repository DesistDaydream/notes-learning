---
title: kubeadm 实现细节
linkTitle: kubeadm 实现细节
date: 2024-06-14T08:56
weight: 4
---

# 概述

> 参考：
>
> -


# kubeadm 核心设计原则

> 参考：
>
> - [官方文档，参考 - kubeadm - 实现细节](https://kubernetes.io/docs/reference/setup-tools/kubeadm/implementation-details/)

`kubeadm init` 和 `kubeadm join` 结合在一起提供了良好的用户体验。`kubeadm init` 和 `kubeadm join` 设置的集群应为：

- Secure 安全——应采用最新的最佳做法
  - 加强 RBAC
  - 使用节点授权器
  - 在控制平面组件之间使用安全通信
  - 在 API 服务器和 kubelet 之间使用安全通信
  - 锁定 kubelet API
  - 锁定对系统组件（例如 kube-proxy 和 CoreDNS）的 API 的访问
  - 锁定引导令牌可以访问的内容
- Easy to use 易于使用——用户只需要运行几个命令即可
  - kubeadm init
  - export KUBECONFIG=/etc/kubernetes/admin.conf
  - kubectl apply -f \<network-of-choice.yaml>
  - kubeadm join --token  :
- Extendable 可扩展
  - 它应该不偏袒任何特定的网络提供商。配置群集网络超出范围
  - 它应该提供使用配置文件来自定义各种参数的可能性

## kubeadm 默认所需使用的值和目录

为了降低复杂性并简化基于 kubeadm 的高级工具的开发，它使用一组有限的常量值来存储众所周知的路径和文件名。

Kubernetes 目录 /etc/kubernetes 在应用程序中是一个常量，因为在大多数情况下，它显然是给定的路径，并且是最直观的位置；其他常量路径和文件名是：

**/etc/kubernetes/manifests/** # 作为 kubelet 在其中查找静态 Pod 清单的路径。静态 Pod 清单的名称为：

- etcd.yaml
- kube-apiserver.yaml
- kube-controller-manager.yaml
- kube-scheduler.yaml

**/etc/kubernetes/** # 作为存储带有控制平面组件标识的 kubeconfig 文件的路径。kubeconfig 文件的名称为：

- kubelet.conf（bootstrap-kubelet.conf 在 TLS 引导期间）
- controller-manager.conf
- scheduler.conf
- admin.conf 用于集群管理员和 kubeadm 本身

**/etc/kubernetes/pki/** # 证书和密钥文件的名称：

- ca.crt，ca.key # 用于 Kubernetes 证书颁发机构
- apiserver.crt，apiserver.key # 用于 API 服务器证书
- apiserver-kubelet-client.crt，apiserver-kubelet-client.key # 用于 API 服务器用于安全连接到 kubelet 的客户端证书
- sa.pub，sa.key # 用于控制器管理员在签署 ServiceAccount 时使用的密钥
- front-proxy-ca.crt，front-proxy-ca.key # 用于前端代理证书颁发机构
- front-proxy-client.crt，front-proxy-client.key # 用于前端代理客户端

# kubeadm init 工作流程内部设计

kubeadm init 的[内部的工作流程](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init/#init-workflow)由一系列要执行工作任务组成

[kubeadm init phase](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init-phase/) 命令允许用户单独调用每个任务，并最终提供可重用和可组合的 API /工具箱，其他 Kubernetes 引导工具，任何 IT 自动化工具或高级用户都可以使用该 API /工具来创建自定义群集。

## Preflight checks(事前准备核对清单)

Kubeadm 在启动 init 之前执行一组 prefight checks(预检检查)，目的是验证先决条件并避免常见的集群启动问题。用户可以选择跳过特定的 prefight checks 或全部跳过--ignore-preflight-errors。

- \[警告]如果要使用的 Kubernetes 版本（带有--kubernetes-version 标志指定）至少比 kubeadm CLI 版本高一个次要版本。
- Kubernetes 系统要求：
  - 如果在 Linux 上运行：
    - \[错误]如果内核早于最低要求的版本
    - \[错误]如果未设置所需的 cgroups 子系统
  - 如果使用 docker：
    - \[警告/错误]如果 Docker 服务不存在，则被禁用，如果它不处于活动状态。
    - \[错误]如果 Docker 端点不存在或不起作用
    - \[警告]如果 docker 版本不在经过验证的 docker 版本列表中
  - 如果使用其他 cri 引擎：
    - \[错误]如果 crictl 套接字未回答
- \[错误]如果用户不是 root 用户
- \[错误]如果计算机主机名不是有效的 DNS 子域
- \[警告]如果无法通过网络查找访问主机名
- \[错误]如果 kubelet 版本低于 kubeadm 支持的最小 kubelet 版本（当前次要-1）
- \[错误]如果 kubelet 版本比所需的控制面板版本至少高一个小（不支持的版本偏斜）
- \[警告]如果 kubelet 服务不存在或已被禁用
- \[警告]如果 firewalld 处于活动状态
- \[错误]如果使用 API 服务器 bindPort 或端口 10250/10251/10252
- \[错误]如果/etc/kubernetes/manifest 文件夹已经存在并且不为空
- \[错误]如果/proc/sys/net/bridge/bridge-nf-call-iptables 文件不存在/不包含 1
- \[错误]如果广告地址是 ipv6，/proc/sys/net/bridge/bridge-nf-call-ip6tables 并且不存在/不包含 1。
- \[错误]如果启用了交换
- \[错误]如果 conntrack，ip，iptables，mount，nsenter 的命令不存在于该命令路径
- \[警告]如果 ebtables，ethtool，socat，tc，touch，crictl 的命令不存在于该命令路径
- \[警告]如果 API 服务器，控制器管理器，调度程序的额外 arg 标志包含一些无效选项
- \[警告]如果与 https://API.AdvertiseAddress:API.BindPort 的连接通过代理
- \[警告]如果服务子网的连接通过代理进行（仅检查第一个地址）
- \[警告]如果 Pods 子网的连接通过代理进行（仅检查第一个地址）
- 如果提供了外部 etcd：
  - \[错误]如果 etcd 版本早于最低要求版本
  - \[错误]如果指定了 etcd 证书或密钥，但未提供
- 如果未提供外部 etcd（因此将安装本地 etcd）：
  - \[错误]如果使用端口 2379
  - \[错误]如果 Etcd.DataDir 文件夹已经存在并且不为空
- 如果授权方式为 ABAC：
  - \[错误]如果 abac_policy.json 不存在
- 如果授权方式为 WebHook
  - \[错误]如果 webhook_authz.conf 不存在

请注意：

- 可以使用以下 [kubeadm init phase preflight](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init-phase/#cmd-phase-preflight) 命令分别调用预检检查

## 生成必要的证书

Kubeadm 为不同目的生成证书和私钥对：

- Kubernetes 集群的自签名证书颁发机构已保存到 ca.crt 文件和 ca.key 私钥文件中
- API 服务器的服务证书，使用 ca.crtCA 生成，并 apiserver.crt 使用其私钥保存到文件中 apiserver.key。该证书应包含以下备用名称：
  - Kubernetes 服务的内部 clusterIP（服务 CIDR 中的第一个地址，例如，10.96.0.1 如果服务子网为 10.96.0.0/12）
  - Kubernetes DNS 名称，例如，kubernetes.default.svc.cluster.local 如果--service-dns-domain 标志值 cluster.local，再加上默认的 DNS 名称 kubernetes.default.svc，kubernetes.default，kubernetes
  - 节点名称
  - 的 --apiserver-advertise-address
  - 用户指定的其他备用名称
- API 服务器安全连接到 kubelet 的客户端证书，使用 ca.crtCA 作为证书生成，并 apiserver-kubelet-client.crt 使用其私钥保存到 文件中 apiserver-kubelet-client.key。该证书应在 system:masters 组织中
- 用于签名保存到 sa.key 文件中的 ServiceAccount 令牌的私钥及其公钥 sa.pub
- 前代理的证书颁发机构 front-proxy-ca.crt 及其密钥保存在文件中 front-proxy-ca.key
- 前代理客户端的客户端证书，使用 front-proxy-ca.crtCA 生成并 front-proxy-client.crt 使用其私钥保存到文件中 front-proxy-client.key

证书默认情况下存储在中/etc/kubernetes/pki，但是可以使用该--cert-dir 标志配置该目录。

请注意：

1. 如果给定证书和私钥对都存在，并且其内容经过评估符合上述规范，则将使用现有文件，并跳过给定证书的生成阶段。这意味着用户可以例如将现有的 CA 复制到 /etc/kubernetes/pki/ca.{crt,key}，然后 kubeadm 将使用这些文件对其余证书进行签名。另请参阅[使用自定义证书](https://kubernetes.io/docs/tasks/administer-cluster/kubeadm/kubeadm-certs#custom-certificates)
2. 仅对于 CA，如果所有其他证书和 kubeconfig 文件已经到位，则可以提供 ca.crt 文件，而不提供文件 ca.keykubeadm 可以识别这种情况并激活 ExternalCA，这也意味着 csrsignercontroller-manager 中的控制器不会开始了
3. 如果 kubeadm 在[外部 CA 模式下运行](https://kubernetes.io/docs/tasks/administer-cluster/kubeadm/kubeadm-certs#external-ca-mode) ; 所有证书必须由用户提供，因为 kubeadm 无法自行生成它们
4. 如果以该--dry-run 模式执行 kubeadm ，则证书文件将写入一个临时文件夹中
5. 证书生成可以使用以下[kubeadm init phase certs all](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init-phase/#cmd-phase-certs)命令单独调用

## 为控制平面组件生成 kubeconfig 文件

Kubeadm 生成具有用于控制平面组件标识的 kubeconfig 文件：

- 在 TLS 引导期间要使用的 kubelet 的 kubeconfig 文件-/etc/kubernetes/bootstrap-kubelet.conf。在此文件中，有一个引导令牌或嵌入式客户端证书，用于通过群集验证此节点。该客户证书应：
  - system:nodes 根据[节点授权](https://kubernetes.io/docs/reference/access-authn-authz/node/)模块的要求在组织中
  - 具有通用名称（CN） system:node:
- 控制器管理器的 kubeconfig 文件/etc/kubernetes/controller-manager.conf；在此文件中嵌入了具有控制器管理员身份的客户端证书。此客户端证书应具有 CN system:kube-controller-manager，如默认[RBAC 核心组件角色](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#core-component-roles)所定义
- 用于调度程序的 kubeconfig 文件/etc/kubernetes/scheduler.conf；在此文件中嵌入了具有调度程序标识的客户端证书。此客户端证书应具有 CN system:kube-scheduler，如默认[RBAC 核心组件角色](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#core-component-roles)所定义

此外，还将生成 kubeadm 本身和 admin 的 kubeconfig 文件，并将其保存到该 /etc/kubernetes/admin.conf 文件中。此处的“管理员”定义为正在管理集群并希望对集群具有完全控制权（root）的实际人员。管理员的嵌入式客户端证书应位于 system:masters 组织中，这是默认的[RBAC 用户面对的角色绑定](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles)定义的 。它还应包含一个 CN。Kubeadm 使用 kubernetes-adminCN。

请注意：

1. ca.crt 证书嵌入在所有 kubeconfig 文件中。
2. 如果存在给定的 kubeconfig 文件，并且其内容经过评估符合上述规范，则将使用现有文件，并跳过给定 kubeconfig 的生成阶段
3. 如果 kubeadm 在[ExternalCA 模式下](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init/#external-ca-mode)运行，则所有必需的 kubeconfig 也必须由用户提供，因为 kubeadm 无法自行生成任何一个
4. 如果以该--dry-run 模式执行 kubeadm ，则 kubeconfig 文件将写入一个临时文件夹中
5. 可以使用以下[kubeadm init phase kubeconfig all](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init-phase/#cmd-phase-kubeconfig)命令分别调用 Kubeconfig 文件的生成

## 生成控制平面组件的静态 Pod 清单

Kubeadm 将用于控制平面组件的静态 Pod 清单文件写入/etc/kubernetes/manifests。Kubelet 会在此目录中监视 Pods 在启动时创建的目录。

静态 Pod 清单共享一组公共属性：

- 所有静态 Pod 都部署在 kube-system 名称空间上
- 所有静态 Pod 获取 tier:control-plane 并 component:{component-name}标记
- 所有静态 Pod 使用 system-node-critical 优先级类别
- hostNetwork: true 在所有静态 Pod 上设置，以允许在配置网络之前启动控制平面；作为结果：
  - 该 address 控制器经理和调度使用，指的 API 服务器 127.0.0.1
  - 如果使用本地 etcd 服务器，则 etcd-servers 地址将设置为 127.0.0.1:2379
- 同时为控制器管理器和调度程序启用了领导者选举
- 控制器管理器和调度器将使用各自的唯一身份引用 kubeconfig 文件
- 如[将自定义参数传递到控制平面组件中](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/control-plane-flags/)所述，所有静态 Pod 都会获得用户指定的任何其他标志
- 所有静态 Pod 都会获得用户指定的任何其他卷（主机路径）

请注意：

1. 默认情况下，将从 k8s.gcr.io 中提取所有图像。请参阅[使用自定义图像](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init/#custom-images)自定义图像存储库
2. 如果以该--dry-run 模式执行 kubeadm，则将静态 Pods 文件写入一个临时文件夹中
3. 可以使用以下[kubeadm init phase control-plane all](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init-phase/#cmd-phase-control-plane)命令分别调用主组件的静态 Pod 清单生成

### API Server

API 服务器的静态 Pod 清单受用户提供的以下参数影响：

- 的 apiserver-advertise-address 和 apiserver-bind-port，以结合; 如果未提供，则这些值默认为计算机和端口 6443 上默认网络接口的 IP 地址
- 将 service-cluster-ip-range 用于服务
- 如果指定的外部服务器 ETCD，所述 etcd-servers 地址和相关的 TLS 设置（etcd-cafile，etcd-certfile，etcd-keyfile）; 如果未提供外部 etcd 服务器，则将使用本地 etcd（通过主机网络）
- 如果指定了云提供商，则将--cloud-provider 配置相应的云提供商，以及--cloud-config 路径（如果存在该文件的话）（这是实验性的 Alpha，将在以后的版本中删除）

其他无条件设置的 API 服务器标志是：

- --insecure-port=0 避免与 api 服务器的连接不安全
- --enable-bootstrap-token-auth=true 启用 BootstrapTokenAuthenticator 身份验证模块。有关更多详细信息，请参见[TLS 引导](https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet-tls-bootstrapping/)
- --allow-privileged 到 true（例如，kube 代理要求）
- --requestheader-client-ca-file 至 front-proxy-ca.crt
- --enable-admission-plugins 至：
  - [NamespaceLifecycle](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#namespacelifecycle) 例如，避免删除系统保留的名称空间
  - [LimitRanger](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#limitranger)并[ResourceQuota](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#resourcequota)限制名称空间
  - [ServiceAccount](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#serviceaccount) 强制执行服务帐户自动化
  - [PersistentVolumeLabel](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#persistentvolumelabel)将区域或区域标签附加到由云提供程序定义的 PersistentVolumes（不推荐使用此准入控制器，并将在以后的版本中将其删除。默认情况下，不明确选择使用 gce 或 aws 作为云使用时，kubeadm 不会从 v1.9 开始部署该控件）提供者）
  - [DefaultStorageClass](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#defaultstorageclass)对 PersistentVolumeClaim 对象实施默认存储类
  - [DefaultTolerationSeconds](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#defaulttolerationseconds)
  - [NodeRestriction](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#noderestriction) 限制 kubelet 可以修改的内容（例如，仅此节点上的 pod）
- --kubelet-preferred-address-types 到 InternalIP,ExternalIP,Hostname;这使得 kubectl logs 在环境和其它 API 服务器 kubelet 沟通工作，其中节点的主机名不解析
- 使用先前步骤中生成的证书的标志：
  - --client-ca-file 至 ca.crt
  - --tls-cert-file 至 apiserver.crt
  - --tls-private-key-file 至 apiserver.key
  - --kubelet-client-certificate 至 apiserver-kubelet-client.crt
  - --kubelet-client-key 至 apiserver-kubelet-client.key
  - --service-account-key-file 至 sa.pub
  - --requestheader-client-ca-file 至 front-proxy-ca.crt
  - --proxy-client-cert-file 至 front-proxy-client.crt
  - --proxy-client-key-file 至 front-proxy-client.key
- 用于保护前端代理（[API Aggregation](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/api-machinery/aggregated-api-servers.md)）通信的其他标志：
  - --requestheader-username-headers=X-Remote-User
  - --requestheader-group-headers=X-Remote-Group
  - --requestheader-extra-headers-prefix=X-Remote-Extra-
  - --requestheader-allowed-names=front-proxy-client

### Controller manager

API 服务器的静态 Pod 清单受用户提供的以下参数影响：

- 如果通过指定 a 调用 kubeadm --pod-network-cidr，则可以通过以下设置启用某些 CNI 网络插件所需的子网管理器功能：
  - --allocate-node-cidrs=true
  - --cluster-cidr 并--node-cidr-mask-size 根据给定的 CIDR 进行标记
- 如果指定了云提供商，则将指定相应的云提供商，--cloud-provider 并指定--cloud-config 路径（如果存在此配置文件）（这是实验性的 Alpha，将在以后的版本中删除）

其他无条件设置的标志是：

- --controllers 为 TLS 引导启用所有默认控制器 plus BootstrapSigner 和 TokenCleanercontrollers。有关更多详细信息，请参见[TLS 引导](https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet-tls-bootstrapping/)
- --use-service-account-credentials 至 true
- 使用先前步骤中生成的证书的标志：
  - --root-ca-file 至 ca.crt
  - --cluster-signing-cert-file 到 ca.crt，如果外部 CA 模式被禁用，否则""
  - --cluster-signing-key-file 到 ca.key，如果外部 CA 模式被禁用，否则""
  - --service-account-private-key-file 至 sa.key

### Scheduler

调度程序的静态 Pod 清单不受用户提供的参数的影响。

### 为本地 etcd 生成静态 Pod 清单

如果用户指定了外部 etcd，则将跳过此步骤，否则 kubeadm 会生成静态 Pod 清单文件，以创建在 Pod 中运行的具有以下属性的本地 etcd 实例：

- 聆听 localhost:2379 和使用 HostNetwork=true
- 使 hostPath 从坐骑 dataDir 到主机的文件系统
- 用户指定的任何其他标志

请注意：

1. k8s.gcr.io 默认情况下，etcd 图像将被拉出。请参阅[使用自定义图像](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init/#custom-images)自定义图像存储库
2. 如果在该 --dry-run 模式下执行 kubeadm ，则将 etcd 静态 Pod 清单写入一个临时文件夹中
3. 可以使用以下[kubeadm init phase etcd local](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init-phase/#cmd-phase-etcd)命令分别调用本地 etcd 的静态 Pod 清单生成

### 可选的动态 Kubelet 配置

要使用此功能，请致电 kubeadm alpha kubelet config enable-dynamic。它将 kubelet 初始化配置写入/var/lib/kubelet/config/init/kubelet 文件。

init 配置用于在此特定节点上启动 kubelet，从而为 kubelet 插入文件提供了一种替代方法。这样的配置将由 kubelet 基本配置代替，如以下步骤中所述。有关其他信息，请参见[通过配置文件设置 Kubelet 参数](https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file)。

请注意：

1. 要使动态 kubelet 配置工作，--dynamic-config-dir=/var/lib/kubelet/config/dynamic 应在以下位置指定标志/etc/systemd/system/kubelet.service.d/10-kubeadm.conf
2. 可以通过将 KubeletConfiguration 对象传递到配置文件 kubeadm init 或 kubeadm join 使用配置文件来更改 kubelet 配置--config some-file.yaml。该 KubeletConfiguration 对象可以与其他对象分开，例如 InitConfiguration 使用---分隔符。有关更多详细信息，请查看 kubeadm config print-default 命令。

## 等待控制平面启动

kubeadm 等待（最多 4 分钟）直到 localhost:6443/healthz（kube-apiserver liveness）返回 ok。但是，为了检测死锁条件，如果 localhost:10255/healthz（小方块活跃度）或 localhost:10255/healthz/syncloop（小方块就绪状态）ok 分别不在 40s 和 60s 之内返回，则 kubeadm 将快速失败。

kubeadm 依靠 kubelet 提取控制平面图像并将其作为静态 Pod 正确运行。控制平面启动后，kubeadm 将完成以下段落中描述的任务。

## (可选)编写基本的 kubelet 配置

功能状态： Kubernetes v1.9 \[alpha]
如果使用以下命令调用 kubeadm --feature-gates=DynamicKubeletConfig：

1. 将 kubelet 基本配置写入名称空间中的 kubelet-base-config-v1.9ConfigMap 中 kube-system
2. 创建 RBAC 规则，以授予对所有引导令牌和所有 kubelet 实例（即 system:bootstrappers:kubeadm:default-node-token 和 system:nodes 组）对该 ConfigMap 的读取访问权限
3. 通过指向 Node.spec.configSource 新创建的 ConfigMap 为初始控制平面节点启用动态 kubelet 配置功能

## 将 kubeadm ClusterConfiguration 保存在 ConfigMap 中以供以后参考

kubeadm 节约传递给配置 kubeadm init 在一个名为 ConfigMap kubeadm-config 下 kube-system 的命名空间。

这将确保将来执行的 kubeadm 动作（例如 kubeadm upgrade）将能够确定实际/当前集群状态并基于该数据做出新的决策。

请注意：

1. 在保存 ClusterConfiguration 之前，从配置中剥离了诸如令牌之类的敏感信息。
2. 可以使用以下[kubeadm init phase upload-config](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init-phase/#cmd-phase-upload-config)命令分别调用主配置的上传

## 将节点标记为控制平面

一旦控制平面可用，kubeadm 将执行以下操作：

- 使用以下命令将节点标记为控制平面 node-role.kubernetes.io/master=""
- 污染节点 node-role.kubernetes.io/master:NoSchedule

请注意：

- 标记控制平面相位可以通过以下[kubeadm init phase mark-control-plane](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init-phase/#cmd-phase-mark-master)命令单独调用

## 为节点加入配置 TLS 引导

Kubeadm 使用[Bootstrap Token](https://kubernetes.io/docs/reference/access-authn-authz/bootstrap-tokens/)进行[身份验证将](https://kubernetes.io/docs/reference/access-authn-authz/bootstrap-tokens/)新节点加入现有集群；有关更多详细信息，请参见[设计建议](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/cluster-lifecycle/bootstrap-discovery.md)。

kubeadm init 确保为该过程正确配置所有内容，这包括以下步骤以及设置 API 服务器和控制器标志，如前几段所述。请注意：

- 可以使用以下[kubeadm init phase bootstrap-token](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init-phase/#cmd-phase-bootstrap-token) 命令配置用于节点的 TLS 引导，执行以下各段中描述的所有配置步骤；或者，每个步骤都可以单独调用

### 创建一个引导令牌

kubeadm init 创建第一个引导令牌，该令牌是自动生成的或由用户提供的带有--token 标志的；如引导令牌规范中所述，令牌应另存为名称空间 bootstrap-token-下的秘密 kube-system。请注意：

- 创建的默认令牌 kubeadm init 将在 TLS 引导过程中用于验证临时用户；这些用户将成为该 system:bootstrappers:kubeadm:default-node-token 组的成员
- 令牌的有效期有限，默认为 24 小时（该间隔可以用—token-ttl 标志更改）
- 可以使用该[kubeadm token](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-token/)命令创建其他令牌，这些令牌还提供了其他有用的令牌管理功能

### 允许加入节点调用 CSR API

Kubeadm 确保 system:bootstrappers:kubeadm:default-node-token 组中的用户能够访问证书签名 API。

这是通过 kubeadm:kubelet-bootstrap 在上述组和默认 RBAC 角色之间创建一个 ClusterRoleBinding 来实现的 system:node-bootstrapper。

### 为新的引导令牌设置自动批准

Kubeadm 确保引导令牌将获得 csrapprover 控制器自动批准的 CSR 请求。

这是通过创建 kubeadm:node-autoapprove-bootstrap 在 system:bootstrappers:kubeadm:default-node-token 组和默认角色之间命名的 ClusterRoleBinding 来实现的 system:certificates.k8s.io:certificatesigningrequests:nodeclient。

system:certificates.k8s.io:certificatesigningrequests:nodeclient 还应创建该角色，向授予 POST 权限/apis/certificates.k8s.io/certificatesigningrequests/nodeclient。

### 通过自动批准设置节点证书轮换

Kubeadm 确保为节点启用证书轮换，并且对节点的新证书请求将获得 csrapprover 控制器自动批准的 CSR 请求。

这是通过创建 kubeadm:node-autoapprove-certificate-rotation 在 system:nodes 组和默认角色之间命名的 ClusterRoleBinding 来实现的 system:certificates.k8s.io:certificatesigningrequests:selfnodeclient。

### 创建公开的名为 cluster-info 的 ConfigMap

此阶段 cluster-info 在 kube-public 名称空间中创建 ConfigMap 。

此外，它还创建了 Role 和 RoleBinding 来授予未经身份验证的用户（即 RBAC 组中的用户 system:unauthenticated）对 ConfigMap 的访问。

注意：

- 对 cluster-infoConfigMap 的访问没有速率限制。如果您将母版暴露在互联网上，则可能会或可能不会出现问题。最糟糕的情况是 DoS 攻击，攻击者使用 kube-apiserver 可以处理的所有运行中请求来为 cluster-info ConfigMap 提供服务

对 kube-public 存在的解释详见官方文档：<https://github.com/kubernetes/community/blob/master/contributors/design-proposals/cluster-lifecycle/bootstrap-discovery.md#new-kube-public-namespace>

## 安装插件

Kubeadm 通过 API 服务器安装内部 DNS 服务器和 kube-proxy 附加组件。请注意：

- 此阶段可以使用[kubeadm init phase addon all](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init-phase/#cmd-phase-addon)命令单独调用。

### Proxy

kube-proxy 在 kube-system 名称空间中创建了一个 ServiceAccount ；然后将 kube-proxy 部署为 DaemonSet：

- 母版的凭据（ca.crt 和 token）来自 ServiceAccount
- 主服务器的位置来自 ConfigMap
- 该 kube-proxyServiceAccount 被绑定到特权 system:node-proxierClusterRole

### DNS

- 在 Kubernetes 版本 1.18 中，不赞成在 kubeadm 中使用 kube-dns，并将在以后的版本中将其删除
- CoreDNS 服务的名称为 kube-dns。这样做是为了防止当用户将群集 DNS 从 kube-dns 切换到 CoreDNS 或反之时，[此处](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init-phase/#cmd-phase-addon)--config 描述的方法引起的服务中断。
- 在 kube-system 名称空间中创建了 CoreDNS / kube-dns 的 ServiceAccount 。
- 该 kube-dnsServiceAccount 被绑定到特权 system:kube-dnsClusterRole

# kubeadm join 阶段内部设计

与相似 kubeadm init，kubeadm join 内部工作流程也包括一系列要执行的原子工作任务。

这分为发现（使节点信任 Kubernetes 主节点）和 TLS 引导程序（使 Kubernetes 主节点信任节点）。

请参阅[使用 Bootstrap 令牌进行身份验证](https://kubernetes.io/docs/reference/access-authn-authz/bootstrap-tokens/)或相应的[设计建议](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/cluster-lifecycle/bootstrap-discovery.md)。

## 飞行前检查

kubeadm 在开始加入之前执行一组预检检查，以验证先决条件并避免常见的集群启动问题。

请注意：

1. kubeadm join 飞行前检查基本上是 kubeadm init 飞行前检查的子集
2. 从 1.9 开始，kubeadm 为 CRI 通用功能提供了更好的支持。在这种情况下，docker 特定的控件将被跳过或替换为 crictl 的类似控件。
3. 从 1.9 开始，kubeadm 支持加入在 Windows 上运行的节点。在这种情况下，将跳过 Linux 特定的控件。
4. 在任何情况下，用户都可以使用该--ignore-preflight-errors 选项跳过特定的飞行前检查（或最终跳过所有飞行前检查）。

## 发现集群信息

有两种主要的发现方案。第一种是使用共享令牌以及 API 服务器的 IP 地址。第二个是提供一个文件（它是标准 kubeconfig 文件的子集）。

## 共享令牌发现

如果使用 kubeadm join 调用--discovery-token，则使用令牌发现；否则，使用令牌发现。在这种情况下，节点基本上是从名称空间中的 cluster-infoConfigMap 检索群集 CA 证书 kube-public。

为了防止“中间人”攻击，采取了以下步骤：

- 首先，通过不安全的连接检索 CA 证书（这是可能的，因为已 kubeadm init 授予的 cluster-info 用户访问权限 system:unauthenticated）
- 然后，CA 证书将通过以下验证步骤：
  - 基本验证：针对 JWT 签名使用令牌 ID
  - 发布密钥验证：使用提供的--discovery-token-ca-cert-hash。此值在 kubeadm init 标准工具的输出中可用，或可以使用标准工具计算（该哈希值是按 RFC7469 中的主题公共密钥信息（SPKI）对象的字节计算的）。该--discovery-token-ca-cert-hash flag 可重复多次，以允许多于一个公钥。
  - 作为附加验证，通过安全连接检索 CA 证书，然后将其与最初检索的 CA 进行比较

请注意：

1. 可以通过发布--discovery-token-unsafe-skip-ca-verification 标志跳过发布密钥验证；这会削弱 kubeadm 安全模型，因为其他人可能会冒充 Kubernetes Master。

## 文件/ https 发现

如果用 kubeadm join 调用--discovery-file，则使用文件发现；否则，使用文件发现。该文件可以是本地文件，也可以通过 HTTPS URL 下载；如果是 HTTPS，则使用主机安装的 CA 捆绑包来验证连接。

通过文件发现，群集 CA 证书将提供到文件本身中。事实上，发现文件是一个文件 kubeconfig 只 server 和 certificate-authority-data 属性设置，如在[kubeadm join](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-join/#file-or-https-based-discovery)参考文档; 与群集建立连接后，kubeadm 尝试访问 cluster-infoConfigMap，并在可用时使用它。

## TLS 引导程序

知道集群信息后，即 bootstrap-kubelet.conf 会写入文件，从而允许 kubelet 执行 TLS 引导（相反，直到由 kubeadm 管理 v.1.7 TLS 引导）。

TLS 引导机制使用共享令牌对 Kubernetes Master 进行临时身份验证，以提交本地创建的密钥对的证书签名请求（CSR）。

然后自动批准该请求，并且该操作完成了保存 ca.crt 文件以及 kubelet.confkubelet 用来加入集群的文件，而该文件 bootstrap-kubelet.conf 已被删除。

请注意：

- 相对于 kubeadm init 过程中保存的令牌（或使用创建的其他令牌 kubeadm token）对临时身份验证进行了验证
- 临时身份验证解析为 system:bootstrappers:kubeadm:default-node-token 在该 kubeadm init 过程中被授予访问 CSR api 权限的组的用户成员
- 自动 CSR 批准由 csrapprover 控制器管理，并根据配置完成 kubeadm init 流程

## (可选)编写 init kubelet 配置

功能状态： Kubernetes v1.9 \[alpha]
如果使用以下命令调用 kubeadm --feature-gates=DynamicKubeletConfig：

1. 使用引导令牌令牌从命名空间中的 kubelet-base-config-v1.9ConfigMap 中读取 kubelet 基本配置 kube-system，并将其作为 kubelet 初始配置文件写入磁盘。/var/lib/kubelet/config/init/kubelet
2. 一旦 kubelet 以 Node 自己的凭据（/etc/kubernetes/kubelet.conf）开始，请更新当前节点配置，指定该节点/ kubelet 配置的源是上述 ConfigMap。

请注意：

1. 要使动态 kubelet 配置工作，--dynamic-config-dir=/var/lib/kubelet/config/dynamic 应在以下位置指定标志/etc/systemd/system/kubelet.service.d/10-kubeadm.conf
