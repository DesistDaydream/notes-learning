---
title: Redis 高可用部署
---

# Docker 部署 Redis 高可用

Docker 部署 Redis Sentinel 模式

Sentinel 模式至少需要 3 个节点，所以这里假设有如下三个节点

- 172.19.42.231

- 172.19.42.232

- 172.19.42.233

### 创建配置文件与存储所在路径

    mkdir -p /opt/redis/config
    mkdir -p /opt/redis/data
    chmod 777 /opt/redis/data

### 启动 Redis

master 节点配置

    cat > /opt/redis/config/redis.conf <<EOF
    save 900 1
    maxmemory 1G
    EOF
    chmod 666 /opt/redis/config/redis.conf

replica 节点配置

    cat > /opt/redis/config/redis.conf <<EOF
    save 900 1
    maxmemory 1G
    replicaof 172.19.42.231 6379
    EOF
    chmod 666 /opt/redis/config/redis.conf

启动 Redis

    docker run -d --name redis \
      --network=host \
      -v /opt/redis/config:/etc/redis \
      -v /opt/redis/data:/data \
      redis:5.0.10-alpine \
      /etc/redis/redis.conf

### 启动 Redis Sentinel

所有节点配置

```bash
cat > /opt/redis/config/sentinel.conf <<EOF
sentinel monitor mymaster 172.19.42.231 6379 2
sentinel down-after-milliseconds mymaster 60000
sentinel failover-timeout mymaster 180000
sentinel parallel-syncs mymaster 1
EOF
chmod 666 /opt/redis/config/sentinel.conf
```

启动 Sentinel

```bash
docker run -d --name redis-sentinel \
  --network=host \
  -v /opt/redis/config:/etc/redis \
  redis:5.0.10-alpine \
  /etc/redis/sentinel.conf \
  --sentinel
```

# Kubernetes 中部署 Redis 高可用

## Helm 官方维护的 Redis-HA Chart

参考：[Helm 官方网站](https://github.com/helm/charts/tree/master/stable/redis-ha)、[GitHub](https://github.com/DandyDeveloper/charts)、[ArtifactHub](https://artifacthub.io/packages/helm/dandydev-charts/redis-ha)

Grafana Dashboard:11835

## 第三方 redis operator 部署 Cluster 模式 Redis

<https://github.com/ucloud/redis-cluster-operator> ucloud 出品

## 第三方 redis operator 部署 Sentinel 模式 redis

<https://github.com/spotahome/redis-operator>，通过 operator 可以简单得创建出 6 个 pod，3 个 redis 节点，3 个 sentinel 节点。

ucloud 基于该项目推出了一个类似的：<https://github.com/ucloud/redis-operator>

部署所需 yaml 在 github 上

创建 operator

    kubectl apply -f https://raw.githubusercontent.com/spotahome/redis-operator/master/example/operator/all-redis-operator-resources.yaml

配置 redis 密码认证

    # “密码”修改为自己想设置的密码
    echo -n "密码" > password
    kubectl create -n redis secret generic redis-auth --from-file=password

部署 redis

    kubectl create -f https://raw.githubusercontent.com/spotahome/redis-operator/master/example/redisfailover/basic.yaml

Bitnami 官方用于部署 redis 的 helm chart

<https://github.com/bitnami/charts/tree/master/bitnami/redis/>

获取 charts 文件

1. helm repo add bitnami <https://charts.bitnami.com/bitnami>

2. helm pull bitnami/redis

3. tar -zxvf redis-XX.X.X.tgz

4. 修改值文件，参考：

部署 redis

1. helm install redis -n redis --set password=oc123 .

Bitnami 版问题：

- 无法故障恢复，删除 pod 后， master 无法切换

- 是有了安全环境容器，导致容器内无法读取 /proc/sys/net/core/somaxconn 参数的值

- 问题跟踪：

  - <https://github.com/bitnami/charts/issues/3700>

  - <https://github.com/bitnami/charts/issues/4569>
