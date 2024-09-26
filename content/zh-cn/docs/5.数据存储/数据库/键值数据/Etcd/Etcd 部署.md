---
title: Etcd 部署
---

# 概述

> 参考：
>
> - [官方文档](https://etcd.io/docs/latest/op-guide/container/)

Etcd 可以通过多种方式部署。如果要启动 etcd 集群，则每种部署方式，都需要配置最基本标志为以下几个：

- --name # etcd 集群中的节点名，这里可以随意，可区分且不重复就行
- --listen-peer-urls # 监听的用于节点之间通信的 url，可监听多个，集群内部将通过这些 url 进行数据交互(如选举，数据同步等)
- --initial-advertise-peer-urls # 建议用于节点之间通信的 url，节点间将以该值进行通信。
- --listen-client-urls # 监听的用于客户端通信的 url，同样可以监听多个。
- --advertise-client-urls # 建议使用的客户端通信 url，该值用于 etcd 代理或 etcd 成员与 etcd 节点通信。
- --initial-cluster-token etcd-cluster-1 # 节点的 token 值，设置该值后集群将生成唯一 id，并为每个节点也生成唯一 id，当使用相同配置文件再启动一个集群时，只要该 token 值不一样，etcd 集群就不会相互影响。
- --initial-cluster # 也就是集群中所有的 initial-advertise-peer-urls 的合集。
- --initial-cluster-state new # 新建集群的标志

如果是单节点部署，则直接启动即可。

# 使用二进制文件部署 etcd

直接使用 yum install etcd -y 命令即可安装

# 在容器内运行 etcd

运行一个单节点的 etcd

```bash
export NODE1=192.168.1.21

# 配置Docker卷以存储etcd数据：
docker volume create --name etcd-data
export DATA_DIR="etcd-data"

# 运行最新版本的etcd：
REGISTRY=quay.io/coreos/etcd
# available from v3.2.5
REGISTRY=gcr.io/etcd-development/etcd

docker run \
  -p 2379:2379 \
  -p 2380:2380 \
  --volume=${DATA_DIR}:/etcd-data \
  --name etcd ${REGISTRY}:latest \
  /usr/local/bin/etcd \
  --data-dir=/etcd-data --name node1 \
  --initial-advertise-peer-urls http://${NODE1}:2380 --listen-peer-urls http://0.0.0.0:2380 \
  --advertise-client-urls http://${NODE1}:2379 --listen-client-urls http://0.0.0.0:2379 \
  --initial-cluster node1=http://${NODE1}:2380
```

## 部署 3 节点 etcd 集群

```bash
REGISTRY=quay.io/coreos/etcd
# available from v3.2.5
REGISTRY=gcr.io/etcd-development/etcd

# For each machine
ETCD_VERSION=latest
TOKEN=my-etcd-token
CLUSTER_STATE=new
NAME_1=etcd-node-0
NAME_2=etcd-node-1
NAME_3=etcd-node-2
HOST_1=10.20.30.1
HOST_2=10.20.30.2
HOST_3=10.20.30.3
CLUSTER=${NAME_1}=http://${HOST_1}:2380,${NAME_2}=http://${HOST_2}:2380,${NAME_3}=http://${HOST_3}:2380
DATA_DIR=/var/lib/etcd

# For node 1
THIS_NAME=${NAME_1}
THIS_IP=${HOST_1}
docker run \
  -p 2379:2379 \
  -p 2380:2380 \
  --volume=${DATA_DIR}:/etcd-data \
  --name etcd ${REGISTRY}:${ETCD_VERSION} \
  /usr/local/bin/etcd \
  --data-dir=/etcd-data --name ${THIS_NAME} \
  --initial-advertise-peer-urls http://${THIS_IP}:2380 --listen-peer-urls http://0.0.0.0:2380 \
  --advertise-client-urls http://${THIS_IP}:2379 --listen-client-urls http://0.0.0.0:2379 \
  --initial-cluster ${CLUSTER} \
  --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}

# For node 2
THIS_NAME=${NAME_2}
THIS_IP=${HOST_2}
docker run \
  -p 2379:2379 \
  -p 2380:2380 \
  --volume=${DATA_DIR}:/etcd-data \
  --name etcd ${REGISTRY}:${ETCD_VERSION} \
  /usr/local/bin/etcd \
  --data-dir=/etcd-data --name ${THIS_NAME} \
  --initial-advertise-peer-urls http://${THIS_IP}:2380 --listen-peer-urls http://0.0.0.0:2380 \
  --advertise-client-urls http://${THIS_IP}:2379 --listen-client-urls http://0.0.0.0:2379 \
  --initial-cluster ${CLUSTER} \
  --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}

# For node 3
THIS_NAME=${NAME_3}
THIS_IP=${HOST_3}
docker run \
  -p 2379:2379 \
  -p 2380:2380 \
  --volume=${DATA_DIR}:/etcd-data \
  --name etcd ${REGISTRY}:${ETCD_VERSION} \
  /usr/local/bin/etcd \
  --data-dir=/etcd-data --name ${THIS_NAME} \
  --initial-advertise-peer-urls http://${THIS_IP}:2380 --listen-peer-urls http://0.0.0.0:2380 \
  --advertise-client-urls http://${THIS_IP}:2379 --listen-client-urls http://0.0.0.0:2379 \
  --initial-cluster ${CLUSTER} \
  --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
```

# 使用证书部署一个安全的 ETCD 集群

## 生成自签名证书

官方文档：

make 自动生成：<https://github.com/etcd-io/etcd/tree/master/hack/tls-setup>

使用 cfssl 工具：<https://github.com/coreos/docs/blob/master/os/generate-self-signed-certificates.md>

## 部署 etcd

在上面的部署示例中，每个 etcd 节点添加如下命令行标志即可

客户端到服务端通信认证所需配置

--client-cert-auth=true # 设置此选项后，etcd 将检查所有传入的 HTTPS 请求以查找由受信任的 CA 签名的客户端证书，未提供有效客户端证书的请求将失败。如果启用了[身份验证](https://github.com/etcd-io/etcd/blob/master/Documentation/op-guide/authentication.md)，则证书将为“公用名”字段提供的用户名提供凭据。

--cert-file=/etc/kubernetes/pki/etcd/server.crt # 用于与 etcd 的 SSL / TLS 连接的证书。设置此选项后，advertise-client-urls 可以使用 HTTPS 模式。

--key-file=/etc/kubernetes/pki/etcd/server.key # 证书密钥。必须未加密。

--trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt # 受信任的证书颁发机构。

服务端到服务端通信认证所需配置

--peer-client-cert-auth=true # 设置后，etcd 将检查来自集群的所有传入对等请求，以查找由提供的 CA 签名的有效客户端证书。

--peer-cert-file=/etc/kubernetes/pki/etcd/peer.crt # 对等体之间用于 SSL / TLS 连接的证书。这将用于侦听对等方地址以及向其他对等方发送请求。

--peer-key-file=/etc/kubernetes/pki/etcd/peer.key # 证书密钥。必须未加密。

--peer-trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt # 受信任的证书颁发机构。
