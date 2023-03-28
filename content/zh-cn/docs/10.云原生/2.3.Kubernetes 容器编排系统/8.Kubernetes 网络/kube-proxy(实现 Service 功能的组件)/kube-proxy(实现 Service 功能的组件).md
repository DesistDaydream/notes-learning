---
title: kube-proxy(实现 Service 功能的组件)
---

# 概述

> 参考：
>
> - [官方文档，概念-概述-Kubernetes 组件-kube-proxy](https://kubernetes.io/docs/concepts/overview/components/#kube-proxy)

kube-proxy 可以转发 Service 的流量到 POD

kube-proxy 有三种模式，userspace、iptables、ipvs。

- service 在逻辑上代表了后端的多个 Pod，外界通过 service 访问 Pod。service 接收到的请求是如何转发到 Pod 的呢？这就是 kube-proxy 要完成的工作。接管系统的 iptables，所有到达 Service 的请求，都会根据 proxy 所定义的 iptables 的规则，进行 nat 转发
- 每个 Node 都会运行 kube-proxy 服务，它负责将访问 service 的 TCP/UPD 数据流转发到后端的容器。如果有多个副本，kube-proxy 会实现负载均衡。
- 每个 Service 的变动(创建，改动，摧毁)都会通知 proxy，在 proxy 所在的本节点创建响应的 iptables 规则，如果 Service 后端的 Pod 摧毁后重新建立了，那么就是靠 proxy 来把 pod 信息提供给 Service。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/cp8r8a/1616118387292-eec78059-6dc3-4131-a895-85ccae5711f3.jpeg)
Note:

- kube-proxy 的 ipvs 模式为 lvs 的 nat 模型
- 如果想要在 ipvs 模式下从 VIP:nodePort 去访问就请你暴露的服务的话，需要将 VIP 的掩码设置为 /32。
  - 参考 issue：<https://github.com/kubernetes/kubernetes/issues/75443>

## kube-proxy 监控指标

kube-proxy 在 10249 端口上暴露监控指标，通过 curl -s http://127.0.0.1:10249/metrics 命令即可获取 Metrics

# kube-proxy 配置

> 参考：
>
> - [官方文档,参考-组件工具-kube-proxy](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-proxy/)(这里是命令行标志)
> - [官方文档,参考-配置 APIs-kube-proxy 配置](https://kubernetes.io/docs/reference/config-api/kube-proxy-config.v1alpha1/)(v1alpha1)(这里是配置文详解)
> - [kube-proxy 代码(v1alpha1)](https://pkg.go.dev/k8s.io/kube-proxy/config/v1alpha1#KubeProxyConfiguration)

kube-proxy 可以通过 **命令行标志**和 **配置文件**来控制运行时行为。与 [kubelet 配置](/docs/10.云原生/2.3.Kubernetes%20 容器编排系统/2.Kubelet%20 节点代理/Kubelet%20 配置详解.md 节点代理/Kubelet 配置详解.md)一样，很多 命令行标志 与 配置文件 具有一一对应的关系。

## 命令行标志详解

**--config=<STRING>** # 加载配置文件的路径。

## 配置文件详解

kubectl get configmaps -n kube-system kube-proxy -o yaml # 在 kubeadm 安装的集群中，kube-proxy 的配置保存在 configmap 中，通过 kubectl 命令进行查看
