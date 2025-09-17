---
title: kubectl CLI
linkTitle: kubectl CLI
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，参考 - kubectl](https://kubernetes.io/docs/reference/kubectl/)
> - [官方文档，任务 - 安装工具 - kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
> - [官方推荐常用命令备忘录](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)

kubectl 所用的 kubeconfig 文件，默认在 `~/.kube/confg`，该文件用于定位 Kubernetes 集群以及与 API Server 交互时进行认证，如果没有认证文件则 API Server 无法处理 kubectl 发出的任何指令并返回错误信息。

如果该文件不存在或配置不全(比如没有指定 current-context 字段)，kubectl 则会向 localhost:8080 发起请求(该端口是 API Server 默认监听的不安全端口，该端口不需要认证即可对集群执行所有操作)。

由于 API Server 默认不开启不安全端口，所以在没有配置文件时，就会报如下错误：`The connection to the server localhost:8080 was refused - did you specify the right host or port?`

如果 kubectl 使用的 KubeConfig 文件中，没有集群的 ca 信息，则会报如下错误：`Error from server (BadRequest): the server rejected our request for an unknown reason`

# kubeclt 安装

## 在 Linux 上安装 kubectl

**Ubuntu**

```bash
sudo apt-get update && sudo apt-get install -y apt-transport-https gnupg2 curl
curl https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | sudo apt-key add -
echo 'deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main' | sudo tee /etc/apt/sources.list.d/kubernetes.list
sudo apt-get update
sudo apt-get install -y kubectl
```

**CentOS**

```bash
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
EOF
yum install -y kubectl
```

## 直接下载二进制文件

```bash
export RELEASE="v1.26.1"
export ARCH="amd64"
curl -LO https://dl.k8s.io/release/${RELEASE}/bin/linux/${ARCH}/kubectl
```

## 下载后处理

安装完成后，使用 `kubectl completion bash > /etc/bash_completion.d/kubectl` 生成自动补全功能。

# kubectl 管理文件与配置

**~/.kube/config** # kubeclt 使用的 kubeconfig 文件的默认路径。kubectl 工具运行时将会使用该文件作为连接 kubernetes 集群的信息

kubeamd 部署的集群一般直接使用 /etc/kubernetes/admin.conf 文件拷贝到 ~/.kube/ 目录下并改名为 config

环境变量

- KUBECONFIG # kubectl 命令加载 kubeconfig 文件的路径

# Syntax(语法)

> 参考：
>
> - [官方文档，参考-kubectl 命令](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands)

**kubectl COMMAND \[TYPE] \[NAME] \[FLAGS]**

- **COMMAND** # 指定要在一个或多个资源进行操作，例如 create，get，describe，delete。
- **TYPE** # 指定资源类型。资源类型不区分大小写，您可以指定单数，复数或缩写形式。
- **NAME** # 指定资源的名称。名称区分大小写。如果省略名称，则显示所有资源的详细信息$ kubectl get pods。
- **FLAGS** # 指定全局命令行标志。例如，可以使用--kubeconfig 指定 kubectl 命令执行所使用的配置文件。

## Global FLAGS(全局标志)

https://kubernetes.io/docs/reference/kubectl/kubectl/#options

- **--kubeconfig=/PATH/TO/FILE** # 指定 kubectl 所要使用的配置文件(需要使用绝对路径)
- **-v, --v NUM** # 指定 kubectl 命令执行的 debug 级别，`默认值: 0`。如果使用高级别，可以看到 RESTful 风格请求 APIServer 时的请求头以及响应头信息。打开调试日志也可以看到每个 API 调用的格式。

# Basic Commands (Beginner)(基本命令(初学者))

create # 从文件或者 stdin 上创建一个资源

expose # 创建一个新的 service 资源

- kubectl expose deployment nginx --name nginx-svc --port 80 --type=NodePort

run # 在集群上创建并运行一个特定的镜像
基于 deployment 或 job 来管理和创建容器

set # 配置应用程序资源，用法详见单独章节

# Basic Commands (Intermediate)(基本命令(中级))

## explain - 解释。列出资源所支持的字段

kubectl explain RESOURCE\[.FIELD1.FELD2...FIELDn] \[options] # 每个 FIELD(字段)都可以用.后面跟字段名来查询这个字段下的描述信息，以及该字段下还可以声明什么字段

- EXAMPLE
  - kubectl explain pods.spec.containers

字段说明：

- -required- # 表示该字段为其父字段的必备字段
- <\[]Object> # 表示该字段下的子字段可以以列表形式定义，使用-符号定义多个该字段
- \#表示该字段需要加字符串来定义该字段，不再包含子字段
- <\[]string> # 表示该字段的字符串以列表形式，前面每个参数都要加-符号，依然要使用子字段来写这些字符串

## get - 显示一个或多个资源

详见：[get](docs/10.云原生/Kubernetes/Kubernetes%20管理/kubectl%20CLI/get.md)

## edit - 编辑服务器上的资源

详见：[对象的创建与修改](docs/10.云原生/Kubernetes/Kubernetes%20管理/kubectl%20CLI/对象的创建与修改.md)

## delete - 通过文件名、标准输入、资源名或者资源表删除资源

EXAMPLE

- kubectl delete deployment nginx-deployment
- kubectl delete -f nginx.yaml
- kubectl delete pods nginx --grace-period=0 --force # 强制删除 nginx 这个 pod

# Deploy Commands(部署命令)

## rollout - 管理资源的滚动更新，用法详见 set,rollout 更新资源命令.note

scale # 为 Deployment, ReplicaSet, Replication Controller, or Job 设置新的容量大小

autoscale # Auto-scale a Deployment, ReplicaSet, or ReplicationController

# Cluster Management Commands(集群管理命令)

certificate # 修改证书资源。Modify certificate resources.

approve # 批准一个证书请求 Approve a certificate signing request

deny # 拒绝一个证书请求。Deny a certificate signing request

## cluster-info - 展示 kubernetes 集群的信息，默认展示 master 运行的位置和 DNS 运行的位置

**kubectl cluster-info SubCommand \[flags] \[OPTIONS]**

```bash
root@desistdaydream:~# kubectl cluster-info
Kubernetes control plane is running at https://k8s-api.bj-net.ehualu.local:6443
KubeDNS is running at https://k8s-api.bj-net.ehualu.local:6443/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

SubCommand

- dump # 为调试和诊断倾倒大量相关信息

EXAMPLE

- kubectl cluster-info # 显示集群信息，效果如图
- kubectl cluster-info dump # 显示集群的 dbug 信息

## top - 显示硬件资源(CPU/内存/存储)的用量

该命令只有在集群部署玩 metrics-server 或者 kube-state-metrics 等资源后，才可以获得数据。显示每个 Node 或者每个 Pod 使用的硬件资源情况，效果如图
**kubectl top \[flags] \[options]**
EXAMPLE

- kubectl top node # 显示所有 Node 的硬件资源使用量
- kubectl top pod --all-namespaces # 显示所有名称空间下的 Pod 对硬件资源的使用量

## cordon - 将指定节点标记为不可调度

## uncordon - 将指定节点标记为可调度

## drain - 排空指定的节点，为维护做准备

给定节点将被标记为不可调度(就是 `cordon` 子命令)，以防止新 Pod 被调度到该节点。如果 APIServer 支持 <http://kubernetes.io/docs/admin/disruptions/>，则 `drain` 会 evicts(驱逐) Pod。否则，它将使用普通的 DELETE 请求删除 Pod。`drain` 会驱逐或删除除 mirror pods (不能通过 API 服务器删除) 之外的所有 pod。如果存在 DaemonSet 管理的 Pod，则不会在没有 --ignore-daemonsets 标志的情况下进行，并且无论如何也不会删除任何 DaemonSet 管理的 Pod，因为这些 Pod 将立即被 DaemonSet 控制器替换，该控制器忽略不可调度的标记。如果有任何 Pod 既不是 mirror pods，也不是由 replicationcontrol,replicaset，DaemonSet，statprit set 或 Job 管理的，则除非使用 --force，否则不会删除任何 Pod。-- force 还将允许在一个或多个 pod 的管理资源丢失时继续删除。

`drain` 命令等待优雅的终止。在命令完成之前，不应在计算机上进行操作。

当您准备好将节点重新投入服务时，请使用 `kubectl uncordon`，这将使节点再次可调度。

## taint - 在一个或多个 node 上更新污点

**kubectl taint NODE NAME KEY_1=VAL_1:TAINT_EFFECT_1 ... KEY_N=VAL_N:TAINT_EFFECT_N \[OPTIONS]**
定义的时候要指明 key，val 以及 effect，注意格式

EXAMPLE

- 删除 master 节点上 dedicated:NoSchedule 这个污点
  - k**ubectl taint nodes master dedicated:NoSchedule-**
- 给 master 节点加一个污点，key 为 node-type，val 为 qa，effect 为 NoExecut
  - **kubectl taint nodes master node-type=qa:NoExecute**

# Troubleshooting and Debugging Commands(故障排除和调试命令)

## debug - 创建调试 Pod 以便对工作负载或节点进行故障排除

详见：[故障处理技巧章节](/docs/10.云原生/Kubernetes/Kubernetes%20管理/性能优化与故障处理/故障处理技巧/故障处理技巧.md)

## describe - 显示特定资源或资源组的详细信息

**kubectl describe (-f FILENAME | TYPE \[NAME_PREFIX | -l label] | TYPE/NAME) \[OPTIONS]**

EXAMPLE

- kubectl describe node
- kubectl describe pod kubernetes-dashboard-87f58dc9-j244f --namespace=kube-system

## logs - 打印出在一个 pod 中的一个 container 的日志

详见：[logs](/docs/10.云原生/Kubernetes/Kubernetes%20管理/kubectl%20CLI/logs.md)

## attach - 连接到一个正在运行的容器上(进入容器)

EXAMPLE

- kubectl attach client-7c9999bd74-76s4t -it # 进入该 pod 中

## exec - 在一个容器中执行一条命令

可执行 `/bin/sh` 命令来进入容器当中

**kubectl exec POD \[-c CONTAINER] -- COMMAND \[args...] \[options]**

OPTIONS

- -i, --stdin=false # 传递 STDIN(标准输入)到这个容器
- -t, --tty=false # STDIN(标准输入)是一个 TTY 终端

EXAMPLE

- kubectl exec -it httpd-79c4f99955-2s8rw -- /bin/sh # 以 TTY 终端的形式传递/bin/sh 命令到容器中

## port-forward - 转发一个或多个本地端口到一个 pod 上

OPTIONS

- **--address IP** # 要监听的地址（逗号分隔），默认为 localhost。 仅接受 IP 或 localhost 为值。 提供 localhost 时，kubectl 将尝试同时绑定 127.0.0.1 和:: 1。

EXAMPLE

- kubectl port-forward -n monitoring prometheus-k8s-0 9090
- 将名为 traefik 的 service 的 8080 和 443 端口，进行端口转发暴露出来，监听的地址是本地 0.0.0.0
  - kubectl port-forward --address 0.0.0.0 service/traefik 8080:8080 443:4443

## proxy - 运行一个到 kubernetes 的 API 服务器的代理程序

在服务器和 Kubernetes API Server 之间创建代理服务器或应用程序级网关。 它还允许在指定的 HTTP 路径上保留静态内容。 所有传入数据都通过一个端口进入，并转发到远程 kubernetes API 服务器端口，但与静态内容路径匹配的路径除外

**kubectl proxy \[--port=PORT] \[--www=static-dir] \[--www-prefix=prefix] \[--api-prefix=prefix] \[options]**

OPTIONS

- --accept-hosts='EXPRESSION' # 代理应接受的主机的正则表达式，每个匹配项以逗号分隔。默认为’localhost$,^127.0.0.1$,\[::1]$'
- --accept-paths='^.\*': Regular expression for paths that the proxy should accept.
- --address='IP' # 代理监听的 IP，默认 127.0.0.1
- --api-prefix='/': Prefix to serve the proxied API under.
- --disable-filter=false: If true, disable request filtering in the proxy. This is dangerous, and can leave you vulnerable to XSRF attacks, when used with an accessible port.
- --keepalive=0s: keepalive specifies the keep-alive period for an active network connection. Set to 0 to disable keepalive.
- -p, --port=8001 # 代理监听的端口， 设置为 0 则选择一个随机端口。默认 8001
- --reject-methods='^$': Regular expression for HTTP methods that the proxy should reject (example
- --reject-methods='POST,PUT,PATCH').
- --reject-paths='/api/._/pods/._/exec,/api/._/pods/._/attach': Regular expression for paths that the proxy should reject. Paths specified here will be rejected even accepted by --accept-paths.
- -u, --unix-socket='': Unix socket on which to run the proxy.
- -w, --www='': Also serve static files from the given directory under the specified prefix.
- -P, --www-prefix='/static/': Prefix to serve static files under, if static file directory is specified.

EXAMPLE

- 在本地 8080 端口上启动 API Server 的一个代理网关，以便使用 curl 直接访问 api server 并获取数据
  - kubectl proxy --port=8080
- kubectl proxy --port=8080 --address=0.0.0.0 --accept-hosts='localhost$,^127.0.0.1$,\[::1]$,172.38.40.212' #

cp             Copy files and directories to and from containers.

auth           Inspect authorization

# Advanced Commands - 高级命令

## diff - Diff live version against would-be applied version

## apply - 通过文件或标准输入将配置应用到资源

详见《[对象的创建与修改](docs/10.云原生/Kubernetes/Kubernetes%20管理/kubectl%20CLI/对象的创建与修改.md)》

## patch - 用 strategic merge、JSON merge、JSON，更新一个资源的字段

**kubectl patch (-f FILENAME | TYPE NAME) -p PATCH \[options]**

## replace - 替换。使用文件或标准输入替换一个资源

详见《[对象的创建与修改](docs/10.云原生/Kubernetes/Kubernetes%20管理/kubectl%20CLI/对象的创建与修改.md)》

## wait - 在一个或多个资源上等待指定的条件达成

### EXAMPLE

Wait for the pod "busybox1" to contain the status condition of type "Ready"

```bash
kubectl wait --for=condition=Ready pod/busybox1
```

The default value of status condition is true; you can wait for other targets after an equal delimiter (compared after Unicode simple case folding, which is a more general form of case-insensitivity):

```bash
kubectl wait --for=condition=Ready=false pod/busybox1
```

Wait for the pod "busybox1" to contain the status phase to be "Running".

```bash
kubectl wait --for=jsonpath='{.status.phase}'=Running pod/busybox1
```

Wait for the pod "busybox1" to be deleted, with a timeout of 60s, after having issued the "delete" command

```bash
kubectl delete pod/busybox1
kubectl wait --for=delete pod/busybox1 --timeout=60s
```

## convert - Convert config files between different API versions

# Settings Commands - 设置命令

## label - 更新对象上的标签

详见 [标签与选择器 文章中的 通过 kubectl 命令设置标签](Label%20and%20Selector(标签和选择器).md and Selector(标签和选择器).md) 章节

annotate       Update the annotations on a resource

completion     Output shell completion code for the specified shell (bash or zsh)

# Other Commands - 其他命令

## api-resources - 显示所支持的所有 API 资源(即对象)

显示的信息包括: NAME(对象名), SHORTNAMES(短名称), APIGROUP(API 组), NAMESPACED, KIND(所属种类), VERBS(动作，即该对象可以执行的命令)

**kubectl api-resources \[OPTIONS]**

OPTIONS

- --namespaced=true|false # 显示所有<是 namesapce 的对象|不是 namespace 的对象>
- -o wide|name # 显示更多信息|只显示对象的名称

EXAMPLE

## api-versions - 以“组/版本”的方式在服务器上显示所支持的所有 API 版本

在编写 yaml 文件中的“apiVersion”字段时，可以使用该命令显示出的组/版本

## config - 使用子命令修改 kubeconfig 文件

用法详见 [config](docs/10.云原生/Kubernetes/Kubernetes%20管理/kubectl%20CLI/config.md)

plugin         Runs a command-line plugin

version        Print the client and server version information
