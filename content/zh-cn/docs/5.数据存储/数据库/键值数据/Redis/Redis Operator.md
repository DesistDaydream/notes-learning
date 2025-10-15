---
title: Redis Operator
---

# 通过 Operator 部署一套生产可用的 Redis 集群

该部署方式基于 [spotahome 的 redis-operator 项目](https://github.com/spotahome/redis-operator)

[这里是](https://github.com/spotahome/redis-operator/blob/master/api/redisfailover/v1/types.go) crd 可用字段的信息

### 创建名称空间

```bash
kubectl create ns redis
```

### 创建 operator

```bash
curl -LO https://raw.githubusercontent.com/spotahome/redis-operator/master/example/operator/all-redis-operator-resources.yaml

# 修改 ClusterRoleBinding 的 namespace 为 redis
sed -i 's/namespace: default/namespace: redis/g' all-redis-operator-resources.yaml
# 添加 ServiceAccount namespace 为 redis
sed -i '/^kind: ServiceAccount/{N;a\  namespace: redis
}' all-redis-operator-resources.yaml
# 添加 Deployment namespace 为 redis
sed -i '/^kind: Deployment/{N;a\  namespace: redis
}' all-redis-operator-resources.yaml

kubectl apply -f all-redis-operator-resources.yaml
```

### 配置 redis 密码认证

```bash
# “密码”修改为自己想设置的密码
echo -n "密码" > password
kubectl create -n redis secret generic redis-auth --from-file=password
```

### 部署 redis

```bash
cat > redis.yaml <<EOF
apiVersion: databases.spotahome.com/v1
kind: RedisFailover
metadata:
  name: tj-test
  namespace: redis
spec:
  sentinel:
    replicas: 3
    resources:
      requests:
        cpu: 100m
      limits:
        memory: 100Mi
  redis:
    hostNetwork: true
    dnsPolicy: ClusterFirstWithHostNet
    replicas: 3
    customConfig:
    - "maxmemory 1073741824"
    resources:
      requests:
        cpu: 500m
        memory: 500Mi
      limits:
        cpu: 1000m
        memory: 1000Mi
  auth:
    secretPath: redis-auth
EOF

kubectl apply -f redis.yaml
```

### 暴露 sentinel

```bash
cat > sentinel-external-service.yaml <<EOF
apiVersion: v1
kind: Service
metadata:
  labels:
    app.kubernetes.io/component: sentinel
    app.kubernetes.io/managed-by: redis-operator
    app.kubernetes.io/name: tj-test
    app.kubernetes.io/part-of: redis-failover
    redisfailovers.databases.spotahome.com/name: tj-test
  name: rfs-tj-test-external
  namespace: redis
spec:
  ports:
  - name: sentinel
    port: 26379
    nodePort: 36379
  type: NodePort
  selector:
    app.kubernetes.io/component: sentinel
    app.kubernetes.io/name: tj-test
    app.kubernetes.io/part-of: redis-failover
EOF

kubectl apply -f sentinel-external-service.yaml
```

### 配置可观察性

部署 exporter 并配置 podmonitor，创建 redis svc 以供

```yaml
cat > redis-metrics-service.yaml <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-exporter
  namespace: monitoring
  labels:
    k8s-app: redis-exporter
spec:
  selector:
    matchLabels:
      k8s-app: redis-exporter
  template:
    metadata:
      labels:
        k8s-app: redis-exporter
    spec:
      containers:
      - name: redis-exporter
        image: oliver006/redis_exporter:latest
        args:
        - --redis.addr=redis://rfr-tj-test-external.redis.svc.cluster.local:6379
        - --redis.password=修改为Redis密码
        ports:
        - containerPort: 9121
          name: http
---
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
  name: redis-exporter
  namespace: monitoring
spec:
  podMetricsEndpoints:
  - interval: 30s
    port: http
    relabelings:
    - sourceLabels:
      - __meta_kubernetes_pod_name
      action: replace
      targetLabel: instance
      regex: (.*redis.*)
  selector:
    matchLabels:
      k8s-app: redis-exporter
  namespaceSelector:
    any: true
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app.kubernetes.io/component: redis
  name: rfr-tj-test-external
  namespace: redis
spec:
  clusterIP: None
  ports:
  - name: sentinel
    port: 6379
  selector:
    app.kubernetes.io/component: redis
EOF

kubectl apply -f redis-metrics-service.yaml
```
