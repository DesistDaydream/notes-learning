---
title: Kubernetes 系统组件指标
linkTitle: Kubernetes 系统组件指标
weight: 20
---

# 概述

> 参考：
>
> - [官方文档](https://kubernetes.io/docs/concepts/cluster-administration/system-metrics/)

Kubernetes 系统组件以 Prometheus 格式暴露监控所需的指标。这种格式是结构化的纯文本，人类和机器都可以很方便得阅读。

Kubernetes 的下面几个系统组件默认都会在 `/metrics` 端点暴露指标信息：

- kubelet
  - kubelet 除了基本 /metrics 端点还会在 /metrics/cadvisor、/metrics/resource、/metrics/probes 这几个端点暴露指标
- kube-apiserver
- kube-controller-manager
- kube-scheduler
- kube-proxy

想要采集这些组件的指标，通常需要 Prometheus 或类似的程序，配置抓取程序，以便定期收集，并将指标存储在时间序列数据库中。

## 访问 https 前准备，获取认证所需信息

与[访问 API Server 的 HTTPS](API%20Server.md Server.md) 的方式一样

### 方法一：使用 kubectl 的配置文件中的证书与私钥

想要访问 https 下的内容，首先需要准备证书与私钥或者 ca 与 token 等等。

1. 首先获取 kubeclt 工具配置文件中的证书与私钥
   1. cat /etc/kubernetes/admin.conf | grep client-certificate-data | awk '{print $2}' | base64 -d > /root/certs/admin.crt
   2. cat /etc/kubernetes/admin.conf | grep client-key-data | awk '{print $2}' | base64 -d > /root/certs/admin.key
2. 确定 CA 文件位置(文件一般在 /etc/kubernetes/pki/ca.crt)
   1. CAPATH=/etc/kubernetes/pki/ca.crt
3. 确定要访问组件的的 IP
   1. IP=172.38.40.212

### 方法二：使用拥有最高权限 ServiceAccount 的 Token 访问 https

1. 创建一个 ServiceAccount
   1. kubectl create serviceaccount test-admin
2. 将该 ServiceAccount 绑定到 cluster-admin 这个 clusterrole，以赋予最高权限
   1. kubectl create clusterrolebinding test-admin --clusterrole=cluster-admin --serviceaccount=default:test-admin
3. 将该 ServiceAccount 的 Token 的值注册到变量中
   1. TOKEN=$(kubectl get secrets test-admin-token-599qd -o jsonpath={.data.token} | base64 -d)
4. 确定 CA 文件位置(文件一般在 /etc/kubernetes/pki/ca.crt)
   1. CAPATH=/etc/kubernetes/pki/ca.crt
5. 确定要访问组件的的 IP
   1. IP=172.38.40.212

Note：也可以从一个具有权限的 ServiceAccount 下的 secret 获取，可以使用现成的，也可以手动创建。比如下面用 promtheus 自带的 token。

1. 如果权限不足，那么访问的时候会报错，比如权限不够，或者认证不通过等等。报错信息有如下几种
   1. no kind is registered for the type v1.Status in scheme "k8s.io/kubernetes/pkg/api/legacyscheme/scheme.go:30"
   2. Unauthorized

### 方法三：官方推荐，类似方法二

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

# kubelet 指标

kubelet 在 10205 端口上的多个端点暴露指标

- **/metrics** # kubelet 程序本身运行情况的指标
- **/metrics/cadvisor** # 容器的各种资源使用情况指标，比如容器的 memory、cpu 使用 等等
- **/metrics/resource** # 容器的各种资源使用情况的总和，只有个别几个指标
- **/metrics/probes** # \[ALPHA]实验性质的端点，统计 kubelet 对容器的探针

## 获取指标

通过 https 接口获取 metrics

- 执行访问 https 前准备方法一
  - 通过证书与私钥访问
    - `curl -k --cert /root/cert/admin.crt --key /root/cert/admin.key https://${IP}:10250/metrics`
    - 在 10250 端口的 /metrics/cadvisor 路径下具有 cadvisor 相关的 metrics
      - `curl -k --cert /root/cert/admin.crt --key /root/cert/admin.key https://${IP}:10250/metrics/cadvisor`
- 执行访问 https 前准备方法二
  - 通过 token 访问
    - `curl --cacert ${CAPATH} -H "Authorization: Bearer ${TOKEN}" https://${IP}:10250/metrics`
    - `curl --cacert ${CAPATH} -H "Authorization: Bearer ${TOKEN}" https://${IP}:10250/metrics/cadvisor`

## 常用指标

主要是 [cadvisor](/docs/10.云原生/Containerization%20implementation/容器管理/观测容器.md) 暴露的指标

# kube-apiserver 指标

## 获取指标

- 执行访问 https 前准备方法一
   - 通过证书与私钥访问
      - `curl --cacert ${CAPATH} --cert /root/certs/admin.crt --key  /root/certs/admin.key  https://${IP}:6443/`
- 执行访问 https 前准备方法二
   - 通过 https 的方式访问 API
      - `curl --cacert ${CAPATH} -H "Authorization: Bearer ${TOKEN}"  https://${IP}:6443/`
- kubeclt
   - `kubectl get --raw /` # 让 kubectl 不再输出标准格式的数据，而是直接向 api server 请求原始数据
- kubectl proxy，一般监听在 6443 端口的 api server 使用该方式，监听在 8080 上的为 http，可直接访问
   - `kubectl proxy --port=8080 --accept-hosts='^localhost$,^127.0.0.1$,^\[::1]$,10.10.100.151' --address='0.0.0.0'` # 在本地 8080 端口上启动 API Server 的一个代理网关，以便使用 curl 直接访问 api server 并使用命令 curl localhost:8080/获取数据
      - 直接访问本地 8080 端口，即可通过 API Server 获取集群所有数据

## 常用指标

# kube-controller-manager 指标

kube-controller-manager 在 10257 端口的 `/metrics` 端点上暴露指标数据

kube-controller-manager 指标提供了控制器内部逻辑的性能度量，如 Go 语言运行时度量、etcd 请求延时、云服务商 API 请求延时、云存储请求延时等。Prometheus 格式的性能度量数据，可以通过 curl http://localhost:10252/metrics 来访问。1.18 版本后，10252 端口被弃用，kube-controller-manager 不再提供不安全的 http 端口。

## 获取指标

通过 https 接口获取 metrics

1. 执行访问 https 前准备方法一
   1. 通过证书与私钥访问
      1. curl -k --cert /root/cert/admin.crt --key /root/cert/admin.key https://${IP}:10257/metrics
2. 执行访问 https 前准备方法二
   1. 通过 token 访问
      1. curl --cacert ${CAPATH} -H "Authorization: Bearer ${TOKEN}"  https://${IP}:10257/metrics

## 常用指标

# kube-scheduler 指标

kube-scheduler 在 10259 端口的 `/metrics` 端点上暴露指标数据

## 获取指标

通过 https 接口获取 metrics

1. 执行访问 https 前准备方法一
   1. 通过证书与私钥访问
      1. curl -k --cert /root/cert/admin.crt --key /root/cert/admin.key https://${IP}:10259/metrics
2. 执行访问 https 前准备方法二
3. 通过 token 访问
   1. curl --cacert ${CAPATH} -H "Authorization: Bearer ${TOKEN}"  https://${IP}:10259/metrics

## 常用指标

# kube-proxy 指标

## 常用指标

# Etcd 指标

Etcd 在 2381 端口的 /metrics 端点暴露指标

## 常用指标
