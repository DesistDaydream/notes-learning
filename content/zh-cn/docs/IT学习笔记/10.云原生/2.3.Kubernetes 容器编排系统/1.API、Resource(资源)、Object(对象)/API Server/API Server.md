---
title: API Server
---

# 概述

> 参考：
> - [官方文档，概念-概述-Kubernetes 组件-kube-apiserver](https://kubernetes.io/docs/concepts/overview/components/#kube-apiserver)
> - [官方文档，参考-通用组件-kube-apiserver](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/)

**API Server 是实现 kubernetes API 的应用程序，它是 Kubernetes 控制平面的一个组件，用以对外暴露 Kubernetes API。**Kubernetes API Server 验证和配置 API 对象的数据，包括 pod、service、replicationcontroller 等。 API Server 为 REST 操作提供服务，并为集群的共享状态提供前端，所有其他组件通过该前端进行交互。

如果是通过 kubeadm 安装的 k8s 集群，那么 API Server 的表现形式就是一个名为 **kube-apiserver 的静态 pod。**kube-apiserver 可以水平扩展，i.e.部署多个 kube-apiserver 以实现高可用，应对高并发请求，到达 kube-apiserver 的流量可以在这些实例之间平衡。

API Server 启动后，默认监听在 6443 端口(http 默认监听在 8080 上)。API Server 是 Kubernetes 集群的前端接口 ，各种客户端工具（CLI 或 UI）以及 Kubernetes 其他组件可以通过它管理集群的各种资源。kubectl 就是 API Server 的客户端程序，实现对 k8s 各种资源的增删改查的功能。各个 node 节点的 kubelet 也通过 master 节点的 API Server 来上报本节点的 Pod 状态。

- 提供集群管理的 REST 风格 API 接口，包括认证授权、数据校验以及集群状态变更等
- 提供其他模块之间的数据交互和通信的枢纽（其他模块通过 API Server 查询或修改数据，只有 API Server 才可以直接操作 etcd）

# API Server 的访问方式：

> 参考：
> - [官方文档，任务-管理集群-使用 Kubernetes API 访问集群](https://kubernetes.io/docs/tasks/administer-cluster/access-cluster-api/)

注意：

1. API Server 默认是安全的，在访问时，应使用 https 协议来操作。
2. 参考 [K8S 认证与授权介绍](7.API%20 访问控制.md 访问控制.md) 文章，学习在访问 API Server 时所遇到的验证问题。

## 使用 kubectl 访问 API

现阶段有 kubectl 工具可以实现对 API Server 的访问

使用 kubectl get --raw / 命令让 kubectl 不再输出标准格式的数据，而是直接向 api server 请求原始数据

## 直接访问 REST API(e.g.使用 curl、浏览器 等方式访问 API)

kubectl 处理对 API 服务器的定位和身份验证。如果你想通过 http 客户端（如 curl 或 wget，或浏览器）直接访问 REST API，你可以通过多种方式对 API 服务器进行定位和身份验证：

1. 以代理模式运行 kubectl(推荐)。 推荐使用此方法，因为它用存储的 apiserver 位置并使用自签名证书验证 API 服务器的标识。 使用这种方法无法进行中间人（MITM）攻击。
2. 另外，你可以直接为 HTTP 客户端提供位置和身份认证。 这适用于被代理混淆的客户端代码。 为防止中间人攻击，你需要将根证书导入浏览器。

比如 curl --request DELETE -cacert ${CAPATH} -H "Authorization: Bearer ${TOKEN}"  https://${IP}:6443/api/v1/namespaces/monitoring/pods/prometheus-k8s-0 -k 这样一个请求就可以将集群内 monitoring 空间下的 prometheus-k8s-0 这个 pod 删除

### 获取认证所需信息

**方法一：使用 kubectl 的配置文件中的证书与私钥**
想要访问 https 下的内容，首先需要准备证书与私钥或者 ca 与 token 等等。

1. 首先获取 kubeclt 工具配置文件中的证书与私钥
   1. cat /etc/kubernetes/admin.conf | grep client-certificate-data | awk '{print $2}' | base64 -d > /root/certs/admin.crt
   2. cat /etc/kubernetes/admin.conf | grep client-key-data | awk '{print $2}' | base64 -d > /root/certs/admin.key
2. 确定 CA 文件位置(文件一般在 /etc/kubernetes/pki/ca.crt)
   1. CAPATH=/etc/kubernetes/pki/ca.crt
3. 确定要访问组件的的 IP
   1. IP=172.38.40.212

**方法二：使用拥有最高权限 ServiceAccount 的 Token 访问 https**

- (可选)创建一个专门存放 SA 的名称空间
  - kubectl create namespace user-sa-manage
- 创建一个 ServiceAccount
  - kubectl create -n user-sa-manage serviceaccount test-admin
- 将该 ServiceAccount 绑定到 cluster-admin 这个 clusterrole，以赋予最高权限
  - kubectl create clusterrolebinding test-admin --clusterrole=cluster-admin --serviceaccount=user-sa-manage:test-admin
- 将该 ServiceAccount 的 Token 的值注册到变量中
  - TOKEN=$(kubectl get -n user-sa-manage secrets -o jsonpath="{.items\[?(@.metadata.annotations\['kubernetes.io/service-account.name']=='test-admin')].data.token}"|base64 -d)
- 确定 CA 文件位置(文件一般在 /etc/kubernetes/pki/ca.crt)
  - CAPATH=/etc/kubernetes/pki/ca.crt
- 确定要访问组件的的 IP
  - IP=172.38.40.212
- 使用令牌玩转 API
  - curl -k $IP/api -H "Authorization: Bearer $TOKEN"

Note：也可以从一个具有权限的 ServiceAccount 下的 secret 获取，可以使用现成的，也可以手动创建。比如下面用 promtheus 自带的 token。

1. 如果权限不足，那么访问的时候会报错，比如权限不够，或者认证不通过等等。报错信息有如下几种
   1. no kind is registered for the type v1.Status in scheme "k8s.io/kubernetes/pkg/api/legacyscheme/scheme.go:30"
   2. Unauthorized
2. TOKEN=$(kubectl get secrets -n monitoring prometheus-k8s-token-q5hm4 --template={{.data.token}} | base64 -d)

**方法三：官方推荐，类似方法二**
官方文档：<https://kubernetes.io/docs/tasks/administer-cluster/access-cluster-api/>

```bash
# 查看所有的集群，因为你的 .kubeconfig 文件中可能包含多个上下文
kubectl config view -o jsonpath='{"Cluster name\tServer\n"}{range .clusters[*]}{.name}{"\t"}{.cluster.server}{"\n"}{end}'
# 从上述命令输出中选择你要与之交互的集群的名称
export CLUSTER_NAME="some_server_name"
# 指向引用该集群名称的 API 服务器
APISERVER=$(kubectl config view -o jsonpath="{.clusters[?(@.name==\"${CLUSTER_NAME}\")].cluster.server}")
# 获得令牌
TOKEN=$(kubectl get secrets -o jsonpath="{.items[?(@.metadata.annotations['kubernetes\.io/service-account\.name']=='default')].data.token}"|base64 -d)
# 使用令牌玩转 API
curl -X GET $APISERVER/api --header "Authorization: Bearer $TOKEN" --insecure
```

### 访问 API Server

1. 执行访问 https 前准备方法一
   1. 通过证书与私钥访问
      1. curl --cacert ${CAPATH} --cert /root/certs/admin.crt --key  /root/certs/admin.key  https://${IP}:6443/
2. 执行访问 https 前准备方法二
   1. 通过 https 的方式访问 API
      1. curl --cacert ${CAPATH} -H "Authorization: Bearer ${TOKEN}"  https://${IP}:6443/
3. kubeclt
   1. kubectl get --raw / #让 kubectl 不再输出标准格式的数据，而是直接向 api server 请求原始数据
4. kubectl proxy，一般监听在 6443 端口的 api server 使用该方式，监听在 8080 上的为 http，可直接访问
   1. kubectl proxy --port=8080 --accept-hosts='^localhost$,^127.0.0.1$,^\[::1]$,10.10.100.151' --address='0.0.0.0' #在本地 8080 端口上启动 API Server 的一个代理网关，以便使用 curl 直接访问 api server 并使用命令 curl localhost:8080/获取数据
      1. 直接访问本地 8080 端口，即可通过 API Server 获取集群所有数据

## 编程方式访问 API

Kubernetes 官方支持  [Go](https://kubernetes.io/zh/docs/tasks/administer-cluster/access-cluster-api/#go-client)、[Python](https://kubernetes.io/zh/docs/tasks/administer-cluster/access-cluster-api/#python-client)、[Java](https://kubernetes.io/zh/docs/tasks/administer-cluster/access-cluster-api/#java-client)、 [dotnet](https://kubernetes.io/zh/docs/tasks/administer-cluster/access-cluster-api/#dotnet-client)、[Javascript](https://kubernetes.io/zh/docs/tasks/administer-cluster/access-cluster-api/#javascript-client)  和  [Haskell](https://kubernetes.io/zh/docs/tasks/administer-cluster/access-cluster-api/#haskell-client)  语言的客户端库。还有一些其他客户端库由对应作者而非 Kubernetes 团队提供并维护。 参考[客户端库](https://kubernetes.io/zh/docs/reference/using-api/client-libraries/)了解如何使用其他语言 来访问 API 以及如何执行身份认证。

### Go 客户端介绍

> 参考：官方文档：<https://github.com/kubernetes/client-go/#compatibility-matrix>
> 详见 [Client Libraries](https://www.yuque.com/go/doc/33161293)

版本控制策略：k8s 版本 1.18.8 对应 client-go 版本 0.18.8，其他版本以此类推。

使用前注意事项：
使用 client-go 之前，需要手动获取对应版本的的 client-go 库。根据版本控制策略，使用如下命令进行初始化

    go mod init client-go-test
    go get k8s.io/client-go@kubernetes-1.19.2

这是一个使用 client-go 访问 API 的基本示例

```go
package main

import (
    "context"
    "fmt"

    v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    // 根据指定的 kubeconfig 创建一个用于连接集群的配置，/root/.kube/config 为 kubectl 命令所用的 config 文件
    config, _ := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
    // 根据 BuildConfigFromFlags 创建的配置，返回一个可以连接集群的指针
    clientset, _ := kubernetes.NewForConfig(config)
    // 根据 NewForConfig 所创建的连接集群的指针，来访问 API，并对集群进行操作
    pods, _ := clientset.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
    fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
}
```

## 从 Pod 中访问 API

从 Pod 内部访问 API 时，定位 API 服务器和向服务器认证身份的操作 与上面描述的外部客户场景不同。

从 Pod 使用 Kubernetes API 的最简单的方法就是使用官方的 客户端库。 这些库可以自动发现 API 服务器并进行身份验证。

### 使用官方客户端库

从一个 Pod 内部连接到 Kubernetes API 的推荐方式为：

- 对于 Go 语言客户端，使用官方的 Go 客户端库。 函数 `rest.InClusterConfig()` 自动处理 API 主机发现和身份认证。 参见这里的一个例子。
- 对于 Python 客户端，使用官方的 Python 客户端库。 函数 `config.load_incluster_config()` 自动处理 API 主机的发现和身份认证。 参见这里的一个例子。
- 还有一些其他可用的客户端库，请参阅客户端库页面。

在以上场景中，客户端库都使用 Pod 的服务账号凭据来与 API 服务器安全地通信。

### 直接访问 REST API

在运行在 Pod 中时，可以通过 `default` 命名空间中的名为 `kubernetes` 的服务访问 Kubernetes API 服务器。也就是说，Pod 可以使用 `kubernetes.default.svc` 主机名 来查询 API 服务器。官方客户端库自动完成这个工作。

向 API 服务器进行身份认证的推荐做法是使用 [服务账号](https://kubernetes.io/zh/docs/tasks/configure-pod-container/configure-service-account/) 凭据。 默认情况下，每个 Pod 与一个服务账号关联，该服务账户的凭证（令牌）放置在此 Pod 中 每个容器的文件系统树中的 `/var/run/secrets/kubernetes.io/serviceaccount/token` 处。

如果由证书包可用，则凭证包被放入每个容器的文件系统树中的 `/var/run/secrets/kubernetes.io/serviceaccount/ca.crt` 处， 且将被用于验证 API 服务器的服务证书。

最后，用于命名空间域 API 操作的默认命名空间放置在每个容器中的 `/var/run/secrets/kubernetes.io/serviceaccount/namespace` 文件中。

### 使用 kubectl proxy

如果你希望不实用官方客户端库就完成 API 查询，可以将 `kubectl proxy` 作为 command 在 Pod 启动一个边车（Sidecar）容器。

这样，`kubectl proxy` 自动完成对 API 的身份认证，并将其暴露到 Pod 的 `localhost` 接口，从而 Pod 中的其他容器可以 直接使用 API。

### 不使用代理

通过将认证令牌直接发送到 API 服务器，也可以避免运行 kubectl proxy 命令。 内部的证书机制能够为链接提供保护。

```bash
# 指向内部 API 服务器的主机名
APISERVER=https://kubernetes.default.svc

# 服务账号令牌的路径
SERVICEACCOUNT=/var/run/secrets/kubernetes.io/serviceaccount

# 读取 Pod 的名字空间
NAMESPACE=$(cat ${SERVICEACCOUNT}/namespace)

# 读取服务账号的持有者令牌
TOKEN=$(cat ${SERVICEACCOUNT}/token)

# 引用内部整数机构（CA）
CACERT=${SERVICEACCOUNT}/ca.crt

# 使用令牌访问 API
curl --cacert ${CACERT} --header "Authorization: Bearer ${TOKEN}" -X GET ${APISERVER}/api
```

输出类似于：

```json
{
  "kind": "APIVersions",
  "versions": ["v1"],
  "serverAddressByClientCIDRs": [
    {
      "clientCIDR": "0.0.0.0/0",
      "serverAddress": "10.0.1.149:443"
    }
  ]
}
```

# API Server 健康检查点

> 参考：[官方文档](https://kubernetes.io/docs/reference/using-api/health-checks/)

Kubernetes API 服务器 提供 API 端点以指示 API 服务器的当前状态。 本文描述了这些 API 端点，并说明如何使用。

### API 健康检查点

Kubernetes API 服务器提供 3 个 API 端点（`healthz`、`livez` 和 `readyz`）来表明 API 服务器的当前状态。 `healthz` 端点已被弃用（自 Kubernetes v1.16 起），你应该使用更为明确的 `livez` 和 `readyz` 端点。 `livez` 端点可与 `--livez-grace-period` 标志一起使用，来指定启动持续时间。 为了正常关机，你可以使用 `/readyz` 端点并指定 `--shutdown-delay-duration` 标志。 检查 API 服务器的 `health`/`livez`/`readyz` 端点的机器应依赖于 HTTP 状态代码。 状态码 `200` 表示 API 服务器是 `healthy`、`live` 还是 `ready`，具体取决于所调用的端点。 以下更详细的选项供操作人员使用，用来调试其集群或专门调试 API 服务器的状态。

以下示例将显示如何与运行状况 API 端点进行交互。

对于所有端点，都可以使用 `verbose` 参数来打印检查项以及检查状态。 这对于操作人员调试 API 服务器的当前状态很有用，这些不打算给机器使用：

    curl -k https://localhost:6443/livez?verbose

或从具有身份验证的远程主机：

    kubectl get --raw='/readyz?verbose'

输出将如下所示：

    [+]ping ok
    [+]log ok
    [+]etcd ok
    [+]poststarthook/start-kube-apiserver-admission-initializer ok
    [+]poststarthook/generic-apiserver-start-informers ok
    [+]poststarthook/start-apiextensions-informers ok
    [+]poststarthook/start-apiextensions-controllers ok
    [+]poststarthook/crd-informer-synced ok
    [+]poststarthook/bootstrap-controller ok
    [+]poststarthook/rbac/bootstrap-roles ok
    [+]poststarthook/scheduling/bootstrap-system-priority-classes ok
    [+]poststarthook/start-cluster-authentication-info-controller ok
    [+]poststarthook/start-kube-aggregator-informers ok
    [+]poststarthook/apiservice-registration-controller ok
    [+]poststarthook/apiservice-status-available-controller ok
    [+]poststarthook/kube-apiserver-autoregistration ok
    [+]autoregister-completion ok
    [+]poststarthook/apiservice-openapi-controller ok
    healthz check passed

Kubernetes API 服务器也支持排除特定的检查项。 查询参数也可以像以下示例一样进行组合：

    curl -k 'https://localhost:6443/readyz?verbose&exclude=etcd'

输出显示排除了 `etcd` 检查：

    [+]ping ok
    [+]log ok
    [+]etcd excluded: ok
    [+]poststarthook/start-kube-apiserver-admission-initializer ok
    [+]poststarthook/generic-apiserver-start-informers ok
    [+]poststarthook/start-apiextensions-informers ok
    [+]poststarthook/start-apiextensions-controllers ok
    [+]poststarthook/crd-informer-synced ok
    [+]poststarthook/bootstrap-controller ok
    [+]poststarthook/rbac/bootstrap-roles ok
    [+]poststarthook/scheduling/bootstrap-system-priority-classes ok
    [+]poststarthook/start-cluster-authentication-info-controller ok
    [+]poststarthook/start-kube-aggregator-informers ok
    [+]poststarthook/apiservice-registration-controller ok
    [+]poststarthook/apiservice-status-available-controller ok
    [+]poststarthook/kube-apiserver-autoregistration ok
    [+]autoregister-completion ok
    [+]poststarthook/apiservice-openapi-controller ok
    [+]shutdown ok
    healthz check passed

### 独立健康检查

**FEATURE STATE:** `Kubernetes v1.19 [alpha]`每个单独的健康检查都会公开一个 http 端点，并且可以单独检查。 单个运行状况检查的模式为 `/livez/<healthcheck-name>`，其中 `livez` 和 `readyz` 表明你要检查的是 API 服务器是否存活或就绪。 `<healthcheck-name>` 的路径可以通过上面的 `verbose` 参数发现 ，并采用 `[+]` 和 `ok` 之间的路径。 这些单独的健康检查不应由机器使用，但对于操作人员调试系统而言，是有帮助的：

    curl -k https://localhost:6443/livez/etcd

# API Server 与 Etcd 的交互方式

数据通过 API Server 时，一般是进行序列化后保存到 etcd 中的，可以使用参数 --etcd-prefix 来指定数据保存在 etcd 中后的地址前缀，默认为 `/registry`

一般情况，保存到 etcd 中后，会省略 Group 与 Version，直接使用 Resource 来作为 etcd 中的路径。比如：URI 为 /api/v1/namespaces/kube-system/pods/kube-apiserver-master1 的 pod 资源，在 etcd 中的存储路径为 /registry/pods/kube-system/kube-apiserver-master1。

而序列化的方式可以通过 --storage-media-type 来指定，默认为 protobuf 。使用这种方式将数据序列化之后，得出来的将会有很多乱码，详见 [Etcd 数据探秘章节](https://www.yuque.com/go/doc/33166015) 中的说明

# kube-apiserver Manifests 示例

    apiVersion: v1
    kind: Pod
    metadata:
      annotations:
        kubeadm.kubernetes.io/kube-apiserver.advertise-address.endpoint: 172.19.42.231:6443
      creationTimestamp: null
      labels:
        component: kube-apiserver
        tier: control-plane
      name: kube-apiserver
      namespace: kube-system
    spec:
      containers:
      - command:
        - kube-apiserver
        - --advertise-address=172.19.42.231
        - --allow-privileged=true
        - --authorization-mode=Node,RBAC
        - --client-ca-file=/etc/kubernetes/pki/ca.crt
        - --enable-admission-plugins=NodeRestriction
        - --enable-bootstrap-token-auth=true
        - --etcd-cafile=/etc/kubernetes/pki/etcd/ca.crt
        - --etcd-certfile=/etc/kubernetes/pki/apiserver-etcd-client.crt
        - --etcd-keyfile=/etc/kubernetes/pki/apiserver-etcd-client.key
        - --etcd-servers=https://127.0.0.1:2379
        - --insecure-port=0
        - --kubelet-client-certificate=/etc/kubernetes/pki/apiserver-kubelet-client.crt
        - --kubelet-client-key=/etc/kubernetes/pki/apiserver-kubelet-client.key
        - --kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname
        - --proxy-client-cert-file=/etc/kubernetes/pki/front-proxy-client.crt
        - --proxy-client-key-file=/etc/kubernetes/pki/front-proxy-client.key
        - --requestheader-allowed-names=front-proxy-client
        - --requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.crt
        - --requestheader-extra-headers-prefix=X-Remote-Extra-
        - --requestheader-group-headers=X-Remote-Group
        - --requestheader-username-headers=X-Remote-User
        - --secure-port=6443
        - --service-account-key-file=/etc/kubernetes/pki/sa.pub
        - --service-cluster-ip-range=10.96.0.0/12
        - --service-node-port-range=30000-60000
        - --tls-cert-file=/etc/kubernetes/pki/apiserver.crt
        - --tls-private-key-file=/etc/kubernetes/pki/apiserver.key
        image: registry.aliyuncs.com/k8sxio/kube-apiserver:v1.19.2
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 8
          httpGet:
            host: 172.19.42.231
            path: /livez
            port: 6443
            scheme: HTTPS
          initialDelaySeconds: 10
          periodSeconds: 10
          timeoutSeconds: 15
        name: kube-apiserver
        readinessProbe:
          failureThreshold: 3
          httpGet:
            host: 172.19.42.231
            path: /readyz
            port: 6443
            scheme: HTTPS
          periodSeconds: 1
          timeoutSeconds: 15
        resources:
          requests:
            cpu: 250m
        startupProbe:
          failureThreshold: 24
          httpGet:
            host: 172.19.42.231
            path: /livez
            port: 6443
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
        - mountPath: /etc/localtime
          name: host-time
          readOnly: true
        - mountPath: /etc/kubernetes/pki
          name: k8s-certs
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
          path: /etc/localtime
          type: ""
        name: host-time
      - hostPath:
          path: /etc/kubernetes/pki
          type: DirectoryOrCreate
        name: k8s-certs
