---
title: K3S 配置详解
---

# 概述

> 参考：
> - [官方文档，安装-配置选项](https://docs.k3s.io/zh/installation/configuration)
> - [官方文档，参考-K3S Server 配置](https://docs.k3s.io/reference/server-config)
> - [官方文档，参考-K3S Agent 配置](https://docs.k3s.io/reference/agent-config)

k3s 可以通过如下如下几种方式配置运行时行为
- 命令行标志
- 环境变量
- 配置文件

k3s 运行时默认读取 `/etc/rancher/k3s/config.yaml` 文件中的值。

## 配置文件与命令行标志之间的对应关系

直接使用 `k3s server` 命令并配置如下配置文件：

```yaml
write-kubeconfig-mode: "0644"  
tls-san:  
- "foo.local"  
node-label:  
- "foo=bar"  
- "something=amazing"  
cluster-init: true
```

等效于：

```bash
k3s server \  
--write-kubeconfig-mode "0644" \  
--tls-san "foo.local" \  
--node-label "foo=bar" \  
--node-label "something=amazing" \  
--cluster-init
```

## 其他说明

/etc/rancher/k3s/registries.yaml

k3s 通过 containerd 来控制容器，在 pull 镜像时，会默认指定 docker.io 为 registry 且不可改。

该配置用于将 docker.io 这个域名解析到指定的 私有镜像仓库地址，这样在 crictl pull IMAGE 时，会去私有镜像仓库拉取镜像。这其中需要指定登录私有镜像仓库的用户名和密码。

mirrors: docker.io: endpoint: - "http://172.38.40.180"configs: "172.38.40.180": auth: username: admin password: Harbor12345

可以将 docker.io 修改为自己设定的域名，只不过这样需要在 pull 的时候，选择镜像的时候加上这个域名。

# k3s server 命令行参数详解

-v value (logging) Number for the log level verbosity (default: 0)
--vmodule value (logging) Comma-separated list of pattern=N settings for file-filtered logging
--log value, -l value (logging) Log to file
--alsologtostderr (logging) Log to standard error as well as file (if set)
--bind-address value (listener) k3s bind address (default: 0.0.0.0)
--https-listen-port value (listener) HTTPS listen port (default: 6443)
--advertise-address value (listener) IP address that apiserver uses to advertise to members of the cluster (default: node-external-ip/node-ip)
--advertise-port value (listener) Port that apiserver uses to advertise to members of the cluster (default: listen-port) (default: 0)
--tls-san value (listener) Add additional hostname or IP as a Subject Alternative Name in the TLS cert
--data-dir value, -d value (data) Folder to hold state default /var/lib/rancher/k3s or ${HOME}/.rancher/k3s if not root
--token value, -t value                    (cluster) Shared secret used to join a server or agent to a cluster \[$K3S_TOKEN]
--token-file value (cluster) File containing the cluster-secret/token \[$K3S\_TOKEN\_FILE]
--write-kubeconfig value, -o value         (client) Write kubeconfig for admin client to this file \[$K3S_KUBECONFIG_OUTPUT]
--write-kubeconfig-mode value (client) Write kubeconfig with this mode \[$K3S_KUBECONFIG_MODE]
--default-local-storage-path value (storage) Default local storage path for local provisioner storage class
--no-deploy value (components) Do not deploy packaged components (valid items: coredns, servicelb, traefik, local-storage, metrics-server)
--disable-scheduler (components) Disable Kubernetes default scheduler
--disable-cloud-controller (components) Disable k3s default cloud controller manager
--disable-network-policy (components) Disable k3s default network policy controller
--node-name value (agent/node) Node name \[$K3S\_NODE\_NAME]
--with-node-id                             (agent/node) Append id to node name
--node-label value                         (agent/node) Registering kubelet with set of labels
--node-taint value                         (agent/node) Registering kubelet with set of taints
--container-runtime-endpoint value         (agent/runtime) Disable embedded containerd and use alternative CRI implementation
--pause-image value                        (agent/runtime) Customized pause image for containerd sandbox
--private-registry value                   (agent/runtime) Private registry configuration file (default: "/etc/rancher/k3s/registries.yaml")
--node-ip value, -i value                  (agent/networking) IP address to advertise for node
--node-external-ip value                   (agent/networking) External IP address to advertise for node
--resolv-conf value                        (agent/networking) Kubelet resolv.conf file \[$K3S_RESOLV_CONF]
--flannel-iface value (agent/networking) Override default flannel interface
--flannel-conf value (agent/networking) Override default flannel config file
--rootless (experimental) Run rootless
--agent-token value (experimental/cluster) Shared secret used to join agents to the cluster, but not servers \[$K3S\_AGENT\_TOKEN]
--agent-token-file value                   (experimental/cluster) File containing the agent secret \[$K3S_AGENT_TOKEN_FILE]
--server value, -s value (experimental/cluster) Server to connect to, used to join a cluster \[$K3S\_URL]
--cluster-init                             (experimental/cluster) Initialize new cluster master \[$K3S_CLUSTER_INIT]
--cluster-reset (experimental/cluster) Forget all peers and become a single cluster new cluster master \[$K3S\_CLUSTER\_RESET]

## 常规选项

### 数据库

**--datastore-endpoint value** #  指定 etcd、Mysql、Postgres 或 Sqlite（默认）数据源名称
- $K3S_DATASTORE_ENDPOINT

**--datastore-cafile value** # (db) TLS Certificate Authority file used to secure datastore backend communication \[$K3S_DATASTORE_CAFILE]

**--datastore-certfile value** # (db) TLS certification file used to secure datastore backend communication \[$K3S\_DATASTORE\_CERTFILE]

**--datastore-keyfile value** # (db) TLS key file used to secure datastore backend communication \[$K3S_DATASTORE_KEYFILE]

## 高级选项

### 网络

**--cluster-cidr value** # 用于 pod IP 的 IPv4/IPv6 网络 CIDR。`默认值：10.42.0.0/16`

**--service-cidr value** # 用于 service IP 的 IPv4/IPv6 网络 CIDR。`默认值：10.43.0.0/16`

**--service-node-port-range VALUE** # 为具有 NodePort 可见性的服务保留的端口范围

--cluster-dns value                        (networking) Cluster IP for coredns service. Should be in your service-cidr range (default: 10.43.0.10)

--cluster-domain value                     (networking) Cluster Domain (default: "cluster.local")

--flannel-backend value                    (networking) One of 'none', 'vxlan', 'ipsec', or 'flannel' (default: "vxlan")

### 定制 Kubernetes 进程的标志

**--etcd-arg value** # etcd 进程的运行时 Flags

**--kubelet-arg value** # 指定 kubelet 的运行时 Flags。

- 比如 k3s server --docker --kubelet-arg cgroup-driver=systemd
- Note：在指定 kubelet 参数时，不用加 -- ，k3s 会自动添加。如果加了-- ，那么就会变成 ----cgroup-driver=systemd。这样 kubelet 是无法启动的。
- 上面仅仅举个例子，k3s 内嵌的 kubelet 不支持 systemd 类型的 cgroup-driver

**--kube-apiserver-arg VALUE** # 自定义 kube-apiserver 进程的命令行 Flags。

**--kube-controller-manager-arg value** # 自定义 kube-controller-manager 进程的命令行 Flags。

**--kube-scheduler-arg value** # (自定义 kube-scheduler 进程的命令行 Flags。

**--kube-proxy-arg value** # kube-proxy 进程的运行时 Flags

--kube-cloud-controller-manager-arg value (flags) Customized flag for kube-cloud-controller-manager process

## 实验选项

**--docker** # 指定 k3s 的 CRI 为 docker。默认为 containerd。
