---
title: K3S 配置详解
---

# 概述

> 参考：
> 
> - [官方文档，安装-配置选项](https://docs.k3s.io/zh/installation/configuration)
> - [官方文档，CLI 工具-server](https://docs.k3s.io/zh/cli/server)
> - [官方文档，CLI 工具-agent](https://docs.k3s.io/zh/cli/agent)

k3s 可以通过如下如下几种方式配置运行时行为

- 命令行标志
- 环境变量
- 配置文件
  - k3s 运行时默认读取 `/etc/rancher/k3s/config.yaml` 文件中的值。

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

## 其他配置文件说明

**一、/etc/rancher/k3s/registries.yaml 文件**

k3s 通过 containerd 来控制容器，在 pull 镜像时，会默认指定 docker.io 为 registry 且不可改。

该配置用于将 docker.io 这个域名解析到指定的 私有镜像仓库地址，这样在使用 `crictl pull IMAGE` 时，会去私有镜像仓库拉取镜像。这其中需要指定登录私有镜像仓库的用户名和密码。

```yaml
mirrors:
  docker.io:
    endpoint:
     - "http://172.38.40.180:8080" # 这里可以直接使用 https://hub-mirror.c.163.com 等等域名以加速从 DockerHub 下载镜像的速度
configs:
  "172.38.40.180:8080":
    auth:
     username: admin
     password: Harbor12345
```

可以将 docker.io 修改为自己设定的域名，只不过这样需要在 pull 的时候，选择镜像的时候加上这个域名。

# k3s server 命令行参数详解

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

### 定制 Kubernetes 的组件

**--disable**([]STRING)	# 禁用 [K3S 封装的一些 K8S 之外的组件](https://docs.k3s.io/zh/installation/packaged-components)。多个组件以 `,` 分割。

**--disable-scheduler** # 	禁用 Kubernetes 默认调度程序

**--disable-cloud-controller** # 禁用 k3s 默认云 Controller Manager

**--disable-kube-proxy**	# 禁用运行 kube-proxy

**--disable-network-policy** # 禁用 K3s 默认网络策略控制器

**--disable-helm-controller** # 禁用 Helm 控制器

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

## 其他



# 最佳实践

## 为 Prometheus 提供 K8S 组件的监控指标

参考：<https://github.com/k3s-io/k3s/issues/3619#issuecomment-973188304>

以 kube-prometheus-stack 部署的 Prometheus 套件为例，修改 Helm 的 Value 文件

```yaml
kubeApiServer:
  enabled: true
kubeControllerManager:
  enabled: true
  endpoints:
  - 192.168.42.10
  - 192.168.42.11
  - 192.168.42.12
kubeScheduler:
  enabled: true
  endpoints:
  - 192.168.42.10
  - 192.168.42.11
  - 192.168.42.12
kubeProxy:
  enabled: true
  endpoints:
  - 192.168.42.10
  - 192.168.42.11
  - 192.168.42.12
kubeEtcd:
  enabled: true
  endpoints:
  - 192.168.42.10
  - 192.168.42.11
  - 192.168.42.12
  service:
    enabled: true
    port: 2381
    targetPort: 2381
```

在 K3S 的配置文件中添加如下内容：

```yaml
kube-controller-manager-arg:
- "bind-address=0.0.0.0"
kube-scheduler-arg:
- "bind-address=0.0.0.0"
kube-proxy-arg:
- "metrics-bind-address=0.0.0.0"
etcd-expose-metrics: true
```

## 个人推荐配置

```yaml
service-node-port-range: "30000-40000"
disable-cloud-controller: true
disable:
  - servicelb
  - traefik
kube-controller-manager-arg:
  - bind-address=0.0.0.0
kube-scheduler-arg:
  - bind-address=0.0.0.0
kube-proxy-arg:
  - proxy-mode=ipvs
  - masquerade-all=true
  - metrics-bind-address=0.0.0.0
etcd-expose-metrics: true
```