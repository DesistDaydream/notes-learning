---
title: RabbitMQ 部署
---

# 概述

> 参考：
> - [官方文档，安装与配置](https://www.rabbitmq.com/download.html)

## 使用 docker 启动单节点 RabbitMQ

```bash
docker run -d --hostname my-rabbit --name rabbit \
-p 15672:15672 -p 5672:5672 \
rabbitmq:3-management
```

## 在 kubernetes 集群中使用 Operator 部署 RabbitMQ

> 参考：
> - [官方文档，安装和配置-Kubernetes Operator](https://www.rabbitmq.com/kubernetes/operator/operator-overview.html)
>   - [安装](https://www.rabbitmq.com/kubernetes/operator/install-operator.html)
>   - [通过 Operator 使用 RabbitMQ 集群](https://www.rabbitmq.com/kubernetes/operator/using-operator.html)

注意：RabbitMQ Operator 会为每一个被其创建的 `rabbitmqclusters.rabbitmq.com` 资源的对象添加 `[finalizers](https://www.yuque.com/desistdaydream/learning/md69s3)` 字段，效果如下：

```yaml
apiVersion: rabbitmq.com/v1beta1
kind: RabbitmqCluster
metadata:
  finalizers:
    - deletion.finalizers.rabbitmqclusters.rabbitmq.com
```

基于此，若删除 RabbitmqCluster 对象前删除了 RabbitMQ Operator，那么 RabbitmqCluster 将无法被删除，除非手动删除 `finalizers` 字段。

### 部署 operator

> 这里会自动创建 rabbitmq-system 名称空间

```bash
kubectl apply -f "https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml"
```

### 创建 pv

```yaml
cat > rabbitmq-pv.yaml << EOF
apiVersion: v1
kind: PersistentVolume
metadata:
  name: rabbitmq-node-1
spec:
  accessModes:
  - ReadWriteOnce
  capacity:
    storage: 10Gi
  local:
    path: /opt/rabbitmq
  nodeAffinity:
    required:
      nodeSelectorTerms:
      - matchExpressions:
        - key: kubernetes.io/hostname
          operator: In
          values:
          - node-1.tj-test
EOF

kubectl apply -f rabbitmq-pv.yaml
```

### 部署 rabbitmq

```yaml
cat > rabbitmq.yaml << EOF
apiVersion: rabbitmq.com/v1beta1
kind: RabbitmqCluster
metadata:
  name: tj-test
  namespace: rabbitmq
EOF

kubectl apply -f rabbitmq.yaml
```

### 创建用户并配置权限

```bash
kubectl exec -it -n rabbitmq rabbitmq-bj-cs-server-0 -- /bin/bash
rabbitmqctl add_user admin admin && rabbitmqctl set_user_tags admin administrator && rabbitmqctl add_vhost test  && rabbitmqctl set_permissions -p test admin '.*' '.*' '.*' && exit
```

# RabbitMQ 部署后验证

```bash
rabbitmqadmin declare vhost name=test
rabbitmqadmin declare queue name=test queue_type=quorum
rabbitmqadmin publish routing_key=test payload="hello world"
rabbitmqadmin get queue=test
rabbitmqadmin get queue=test ackmode=ack_requeue_false
```

```bash
rabbitmqadmin declare exchange name=test.topic type=topic
rabbitmqadmin declare binding source=test.topic destination=test routing_key=my.#
rabbitmqadmin publish routing_key=my.test exchange=my.topic  payload="hello world by my.test"
rabbitmqadmin publish routing_key=my.test.test exchange=my.topic  payload="hello world by my.test.test"
rabbitmqadmin get queue=test count=2
```

# RabbitMQ 监控部署

> 参考：
> - [官方文档，安装与配置-Kubernetes Operator-在 Kubernetes 上监控 RabbitMQ 集群-使用 Prometheus Operator 监控 RabbitMQ](https://www.rabbitmq.com/kubernetes/operator/operator-monitoring.html#prom-operator)
> - [官方文档，安装与配置-Kubernetes Operator-在 Kubernetes 上监控 RabbitMQ 集群-导入 Grafana 面板](https://www.rabbitmq.com/kubernetes/operator/operator-monitoring.html#grafana)

为 prometheus-operator 赋予更多权限

```yaml
export PrometheusOperatorSA="bj-test-k8s-operator"

cat > prometheus-cluster-role-and-binding.yaml <<EOF
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: prometheus
rules:
- apiGroups: [""]
  resources:
  - nodes
  - services
  - endpoints
  - pods
  verbs: ["get", "list", "watch"]
- apiGroups: [""]
  resources:
  - configmaps
  verbs: ["get"]
- nonResourceURLs: ["/metrics"]
  verbs: ["get"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: prometheus
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: prometheus
subjects:
- kind: ServiceAccount
  name: ${PrometheusOperatorSA}
  namespace: monitoring
EOF
```

添加 rabbitmq 的监控配置

```yaml
cat > rabbitmq-podmonitor.yaml <<EOF
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
  name: rabbitmq
  namespace: monitoring
spec:
  podMetricsEndpoints:
  - interval: 15s
    port: prometheus
  selector:
    matchLabels:
      app.kubernetes.io/component: rabbitmq
  namespaceSelector:
    any: true
EOF
```

导入 grafana 面板(面板 ID：10991)
