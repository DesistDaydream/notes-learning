---
title: RabbitMQ Operator
---

# 概述

> 参考：
>
> - 官方文档，
> - <https://www.rabbitmq.com/kubernetes/operator/operator-overview.html>
> - <https://www.rabbitmq.com/kubernetes/operator/using-operator.html#override>

# 通过 Operator 部署一套生产可用的 RabbitMQ 集群

### 创建 namespace

    kubectl create ns rabbitmq

### 部署 operator

    kubectl apply -f "https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml"

### 创建 PV

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
    ---
    apiVersion: v1
    kind: PersistentVolume
    metadata:
      name: rabbitmq-node-2
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
              - node-2.tj-test
    ---
    apiVersion: v1
    kind: PersistentVolume
    metadata:
      name: rabbitmq-node-3
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
              - node-3.tj-test
    EOF

    kubectl apply -f rabbitmq-pv.yaml

### 部署 RabbitMQ 集群

若使用默认 image ，则镜像中无 python 环境，也无 rabbitmqadmin 工具。

    cat > rabbitmq.yaml <<EOF
    apiVersion: rabbitmq.com/v1beta1
    kind: RabbitmqCluster
    metadata:
      name: tj-test
      namespace: rabbitmq
    spec:
      rabbitmq:
        additionalPlugins:
        - rabbitmq_delayed_message_exchange
      image: registry.cn-zhangjiakou.aliyuncs.com/ehl_common/rabbitmq:3.8.3-xdelay
      replicas: 3
      service:
        type: NodePort
      override:
        # clientService 再新版crd中变成了 service
        clientService:
          spec:
            ports:
            - name: http
              protocol: TCP
              nodePort: 45672
              port: 15672
            - name: amqp
              protocol: TCP
              nodePort: 35672
              port: 5672

    EOF

    kubectl apply -f rabbitmq.yaml

### 创建用户并配置权限

    kubectl exec -it -n rabbitmq rabbitmq-tj-test-server-0 -- /bin/bash
    rabbitmqctl add_user admin admin && rabbitmqctl set_user_tags admin administrator && rabbitmqctl add_vhost test  && rabbitmqctl set_permissions -p test admin '.*' '.*' '.*' && rabbitmqctl set_permissions -p / admin '.*' '.*' '.*' && exit

### 测试是否可以创建队列

    rabbitmqadmin --vhost=test declare queue name=test-1 queue_type=quorum
