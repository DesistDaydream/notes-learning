---
title: 为 Pod 注入数据
---

# 概述

> 参考：
>
> - [官方文档, 任务 - 给应用注入数据](https://kubernetes.io/docs/tasks/inject-data-application/)

# 将 Pod 的 Manifests 信息映射到容器中的环境变量上

## 用 Pod 字段作为环境变量的值

在这个练习中，你将创建一个包含一个容器的 Pod。这是该 Pod 的配置文件：

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: dapi-envars-fieldref
spec:
  containers:
    - name: test-container
      image: k8s.gcr.io/busybox
      command: ["sh", "-c"]
      args:
        - while true; do
          echo -en '\n';
          printenv MY_NODE_NAME MY_POD_NAME MY_POD_NAMESPACE;
          printenv MY_POD_IP MY_POD_SERVICE_ACCOUNT;
          sleep 10;
          done;
      env:
        - name: MY_NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
        - name: MY_POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: MY_POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: MY_POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        - name: MY_POD_SERVICE_ACCOUNT
          valueFrom:
            fieldRef:
              fieldPath: spec.serviceAccountName
  restartPolicy: Never
```

这个配置文件中，你可以看到五个环境变量。`env` 字段是一个 [EnvVars](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#envvar-v1-core). 对象的数组。 数组中第一个元素指定 `MY_NODE_NAME` 这个环境变量从 Pod 的 `spec.nodeName` 字段获取变量值。 同样，其它环境变量也是从 Pod 的字段获取它们的变量值。

> **说明：** 本示例中的字段是 Pod 字段，不是 Pod 中 Container 的字段。

创建 Pod：

    kubectl apply -f https://k8s.io/examples/pods/inject/dapi-envars-pod.yaml

验证 Pod 中的容器运行正常：

    kubectl get pods

查看容器日志：

    kubectl logs dapi-envars-fieldref

输出信息显示了所选择的环境变量的值：

    minikube
    dapi-envars-fieldref
    default
    172.17.0.4
    default

要了解为什么这些值在日志中，请查看配置文件中的`command` 和 `args`字段。 当容器启动时，它将五个环境变量的值写入 stdout。每十秒重复执行一次。
接下来，通过打开一个 Shell 进入 Pod 中运行的容器：

    kubectl exec -it dapi-envars-fieldref -- sh

在 Shell 中，查看环境变量：

    /# printenv

输出信息显示环境变量已经设置为 Pod 字段的值。

    MY_POD_SERVICE_ACCOUNT=default
    ...
    MY_POD_NAMESPACE=default
    MY_POD_IP=172.17.0.4
    ...
    MY_NODE_NAME=minikube
    ...
    MY_POD_NAME=dapi-envars-fieldref

## 用 Container 字段作为环境变量的值

前面的练习中，你将 Pod 字段作为环境变量的值。 接下来这个练习中，你将用 Container 字段作为环境变量的值。这里是包含一个容器的 Pod 的配置文件：
[`pods/inject/dapi-envars-container.yaml`](https://notes-learning.oss-cn-beijing.aliyuncs.com/ooyi9u/1621520643090-327ae7d2-cb76-4240-b963-9070371cdaca.svg)

    apiVersion: v1
    kind: Pod
    metadata:
      name: dapi-envars-resourcefieldref
    spec:
      containers:
        - name: test-container
          image: k8s.gcr.io/busybox:1.24
          command: [ "sh", "-c"]
          args:
          - while true; do
              echo -en '\n';
              printenv MY_CPU_REQUEST MY_CPU_LIMIT;
              printenv MY_MEM_REQUEST MY_MEM_LIMIT;
              sleep 10;
            done;
          resources:
            requests:
              memory: "32Mi"
              cpu: "125m"
            limits:
              memory: "64Mi"
              cpu: "250m"
          env:
            - name: MY_CPU_REQUEST
              valueFrom:
                resourceFieldRef:
                  containerName: test-container
                  resource: requests.cpu
            - name: MY_CPU_LIMIT
              valueFrom:
                resourceFieldRef:
                  containerName: test-container
                  resource: limits.cpu
            - name: MY_MEM_REQUEST
              valueFrom:
                resourceFieldRef:
                  containerName: test-container
                  resource: requests.memory
            - name: MY_MEM_LIMIT
              valueFrom:
                resourceFieldRef:
                  containerName: test-container
                  resource: limits.memory
      restartPolicy: Never

这个配置文件中，你可以看到四个环境变量。`env` 字段是一个 [EnvVars](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#envvar-v1-core). 对象的数组。数组中第一个元素指定 `MY_CPU_REQUEST` 这个环境变量从 Container 的 `requests.cpu` 字段获取变量值。同样，其它环境变量也是从 Container 的字段获取它们的变量值。

> **说明：** 本例中使用的是 Container 的字段而不是 Pod 的字段。

创建 Pod：

    kubectl apply -f https://k8s.io/examples/pods/inject/dapi-envars-container.yaml

验证 Pod 中的容器运行正常：

    kubectl get pods

查看容器日志：

    kubectl logs dapi-envars-resourcefieldref

输出信息显示了所选择的环境变量的值：

    1
    1
    33554432
    67108864
