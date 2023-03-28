---
title: kube-scheduler 实现调度器的程序
---

# kube-scheduler 是实现 kuberntes Scheduler 的应用程序

kube-scheduler 启动后监听两个端口：

1. 10251 端口为无需身份验证和授权即可不安全地为 HTTP 服务的端口。(1.18 版本后将要弃用)
2. 10259 端口为需要身份验证和授权为 HTTPS 服务的端口。

## kube-scheduler 高科用

与 [kube-controller-manager 高可用](/docs/10.云原生/2.3.Kubernetes%20 容器编排系统/4.Controller(控制器)/kube-controller-manager%20 实现控制器的程序.md 实现控制器的程序.md) 原理相同。

## kube-scheduler 监控指标

详见：[kubernetes 监控](/docs/10.云原生/2.3.Kubernetes%20 容器编排系统/Kubernetes%20 管理/Kubernetes%20 监控/Kubernetes%20 系统组件指标.md 管理/Kubernetes 监控/Kubernetes 系统组件指标.md)

# Kube-scheduler 参数详解

> 参考：
>
> - [官方文档，参考-组件工具-kube-scheduler](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-scheduler/)

## 默认的 manifest 示例

```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    component: kube-scheduler
    tier: control-plane
  name: kube-scheduler
  namespace: kube-system
spec:
  containers:
    - command:
        - kube-scheduler
        - --authentication-kubeconfig=/etc/kubernetes/scheduler.conf
        - --authorization-kubeconfig=/etc/kubernetes/scheduler.conf
        - --bind-address=0.0.0.0
        - --kubeconfig=/etc/kubernetes/scheduler.conf
        - --leader-elect=true
      image: k8s.gcr.io/kube-scheduler:v1.16.3
      imagePullPolicy: IfNotPresent
      livenessProbe:
        failureThreshold: 8
        httpGet:
          host: 127.0.0.1
          path: /healthz
          port: 10251
          scheme: HTTP
        initialDelaySeconds: 15
        timeoutSeconds: 15
      name: kube-scheduler
      resources:
        requests:
          cpu: 100m
      volumeMounts:
        - mountPath: /etc/kubernetes/scheduler.conf
          name: kubeconfig
          readOnly: true
  hostNetwork: true
  priorityClassName: system-cluster-critical
  volumes:
    - hostPath:
        path: /etc/kubernetes/scheduler.conf
        type: FileOrCreate
      name: kubeconfig
```
