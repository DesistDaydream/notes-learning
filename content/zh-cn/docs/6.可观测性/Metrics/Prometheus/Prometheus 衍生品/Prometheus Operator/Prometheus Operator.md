---
title: Prometheus Operator
---

# 概述

> 参考：
>
> - [官网](https://prometheus-operator.dev/)
> - [GitHub 项目](https://github.com/prometheus-operator/prometheus-operator)

该项目曾经在 [coreos/prometheus-operator](https://github.com/coreos/prometheus-operator) 仓库中，后来移动到 prometheus-operator/prometheus-operator

## 背景

为什么会需要 prometheus-operator(后文简称 operator)

当 prometheus 需要监控 kubernetes 集群时，要手动修改配置文件中的 scrape 配置段是非常复杂且繁琐的。每启动一个新 pod 就要新加配置，并手动更新 prometheus 配置文件，有没有一种办法可以在新增 pod 时，让 prometheus 自动更新其配置文件呢？这就是 operator 的作用。

Prometheus Operator 通过数个 CRD 资源来控制 Prometheus 监控套件的运行，并作为这几个 CRD 的 controller(类似于 kube-controller-manager，只不过这个 Controller 只维护几个自定义的资源)来维护其正常运行，这些 CRD 就可以实现这样的功能：自动添加配置文件中 scrape 配置段的 job，并且自动执行热更新来加载配置文件等等。下面是这几个 CRD 的简介

## CRD 介绍

> 参考：
>
> - [官方文档](https://prometheus-operator.dev/docs/operator/design)
> - https://github.com/coreos/prometheus-operator/blob/master/Documentation/design.md

Prometheus Operator 现阶段引入了如下几种自定义资源：

- [Prometheus](https://prometheus-operator.dev/docs/operator/design/#prometheus) # 它定义了所需的 Prometheus 主程序。Operator 始终确保正在运行与资源定义匹配的 prometheus 主程序。
- [Alertmanager](https://prometheus-operator.dev/docs/operator/design/#alertmanager) # 它定义了所需的 Alertmanager 主程序。Operator 始终确保正在运行与资源定义匹配的 Alertmanager 主程序。
- [ThanosRuler](https://prometheus-operator.dev/docs/operator/design/#thanosruler) #
- [ServiceMonitor](https://prometheus-operator.dev/docs/operator/design/#servicemonitor) # 为 Prometheus Server 配置文件中的 scrape_config 配置段生成配置内容。以声明方式指定应如何监控服务组。
- [PodMonitor](https://prometheus-operator.dev/docs/operator/design/#podmonitor) # 为 Prome theus Server 配置文件中的 scrape_config 配置段生成配置内容。与 ServiceMonitor 类型类似，只不过是从指定的 pod 中，发现待抓去的目标。
- [Probe](https://prometheus-operator.dev/docs/operator/design/#probe) # 为 Prometheus Server 配置文件中的 scrape_config 配置段生成配置内容。只会生成 blackbox-exporter 程序所需的配置。
- [PrometheusRule](https://prometheus-operator.dev/docs/operator/design/#prometheusrule) # 它定义了一个所需的 Prometheus 规则文件，该文件可以由包含 Prometheus 警报和记录规则的 Prometheus 实例加载。
- [AlertmanagerConfig](https://prometheus-operator.dev/docs/operator/design/#alertmanagerconfig) #

随着发展，也许还会有其他的 CR 产生

其中 ServiceMonitor、PodMonitor、Probe、PrometheusRule 这几个资源，会被 Operator 监听，并通知配置换换程序将其转换为 Prometheus Server 的配置文件中的内容

### Prometheus

详见：[Prometheus](/docs/6.可观测性/Metrics/Prometheus/Prometheus%20衍生品/Prometheus%20Operator/CR%20详解/Prometheus/Prometheus.md) CR 详解

### Alertmanager

### ThanosRuler

### ServiceMonitor(简称 SM。。。囧)

详见：[Service Monitor](/docs/6.可观测性/Metrics/Prometheus/Prometheus%20衍生品/Prometheus%20Operator/CR%20详解/Service%20Monitor.md)

### PodMonitor

详见：[Pod Monitor](/docs/6.可观测性/Metrics/Prometheus/Prometheus%20衍生品/Prometheus%20Operator/CR%20详解/Pod%20Monitor.md)

### Probe

Probe CRD 定义应如何监视分组和静态目标。除目标外，该`Probe`对象还需要一个`prober`服务，该服务可监视目标并提供 Prometheus 进行刮擦的度量。例如，可以使用 [blackbox-exporter](https://github.com/prometheus/blackbox_exporter/) 来实现。

### PrometheusRule

它定义了一个所需的 Prometheus 规则文件，该文件可以由包含 Prometheus 警报和记录规则的 Prometheus 实例加载。

### Alertmanager

它定义了所需的 Alertmanager 部署。operator 始终确保正在运行与资源定义匹配的部署。

PrometheusRule:对于 Prometheus 而言，在原生的管理方式上，我们需要手动创建 Prometheus 的告警文件，并且通过在 Prometheus 配置中声明式的加载。而在 Prometheus Operator 模式中，告警规则也编程一个通过 Kubernetes API 声明式创建的一个资源.告警规则创建成功后，通过在 Prometheus 中使用想 servicemonitor 那样用 ruleSelector 通过 label 匹配选择需要关联的 PrometheusRule 即可

# Prometheus Operator 部署

    curl -LO https://raw.githubusercontent.com/coreos/prometheus-operator/master/bundle.yaml

该文件会在 default 名称空间里创建 operator。如果要放在其他 namespace 中，需要修改一下 bundle.yaml 文件中 namespace 的值，并修改 clusterrolebinding 中引用的 ServiceAccount 的名称空间。

    [root@master-1 prometheus-operator]# kubectl apply -f bundle.yaml
    customresourcedefinition.apiextensions.k8s.io/alertmanagerconfigs.monitoring.coreos.com created
    customresourcedefinition.apiextensions.k8s.io/alertmanagers.monitoring.coreos.com created
    customresourcedefinition.apiextensions.k8s.io/podmonitors.monitoring.coreos.com created
    customresourcedefinition.apiextensions.k8s.io/probes.monitoring.coreos.com created
    customresourcedefinition.apiextensions.k8s.io/prometheuses.monitoring.coreos.com created
    customresourcedefinition.apiextensions.k8s.io/prometheusrules.monitoring.coreos.com created
    customresourcedefinition.apiextensions.k8s.io/servicemonitors.monitoring.coreos.com created
    customresourcedefinition.apiextensions.k8s.io/thanosrulers.monitoring.coreos.com created
    clusterrolebinding.rbac.authorization.k8s.io/prometheus-operator created
    clusterrole.rbac.authorization.k8s.io/prometheus-operator created
    deployment.apps/prometheus-operator created
    serviceaccount/prometheus-operator created
    service/prometheus-operator created

部署成功后会有一个名为 prometheus-operator 的 deployment、相关的 RBAC(ServiceAccount、ClusterRole、ClusterRoleBinding)、一个 service。还有几个 CRD。

```bash
[root@master-1 prometheus-operator]# kubectl get -n monitor all
NAME                                       READY   STATUS    RESTARTS   AGE
pod/prometheus-operator-6cdb7d79fb-mgv97   1/1     Running   0          35s

NAME                          TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)    AGE
service/prometheus-operator   ClusterIP   None         <none>        8080/TCP   36s

NAME                                  READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/prometheus-operator   1/1     1            1           36s

NAME                                             DESIRED   CURRENT   READY   AGE
replicaset.apps/prometheus-operator-6cdb7d79fb   1         1         1       36s

[root@master-1 prometheus-operator]# kubectl get clusterrole,clusterrolebindings | grep prometheus
clusterrole.rbac.authorization.k8s.io/prometheus-operator                                                    2021-01-22T15:43:34Z
clusterrolebinding.rbac.authorization.k8s.io/prometheus-operator                                    ClusterRole/prometheus-operator                                                    58s
```

## 使用 helm 快速部署一个 prometheus operator 套件

在 [Artifact Hub 上有官方发布的 chart 包](https://artifacthub.io/packages/helm/prometheus-community/kube-prometheus-stack)

为适应 eHualu 生成部署，添加了一个名为 custom 的 subchart 。具体详见 GitHub

**其他**

其他的在安装时使用 -f 参数使用自定义的值文件覆盖即可。

# prometheus-operator 的局限

[Vermouth 博客，高可用 Prometheus 问题集锦](http://www.xuyasong.com/?p=1921)

如果你是在 K8S 集群内部署 Prometheus，那大概率会用到 prometheus-operator，他对 Prometheus 的配置做了 CRD 封装，让用户更方便的扩展 Prometheus 实例，同时 prometheus-operator 还提供了丰富的 Grafana 模板，包括上面提到的 master 组件监控的 Grafana 视图，operator 启动之后就可以直接使用，免去了配置面板的烦恼。

operator 的优点很多，就不一一列举了，只提一下 operator 的局限：

- 因为是 operator，所以依赖 K8S 集群，如果你需要二进制部署你的 Prometheus，如集群外部署，就很难用上 prometheus-operator 了，如多集群场景。当然你也可以在 K8S 集群中部署 operator 去监控其他的 K8S 集群，但这里面坑不少，需要修改一些配置。
- operator 屏蔽了太多细节，这个对用户是好事，但对于理解 Prometheus 架构就有些 gap 了，比如碰到一些用户一键安装了 operator，但 Grafana 图表异常后完全不知道如何排查，record rule 和 服务发现还不了解的情况下就直接配置，建议在使用 operator 之前，最好熟悉 prometheus 的基础用法。
- operator 方便了 Prometheus 的扩展和配置，对于 alertmanager 和 exporter 可以很方便的做到多实例高可用，但是没有解决 Prometheus 的高可用问题，因为无法处理数据不一致，operator 目前的定位也还不是这个方向，和 Thanos、Cortex 等方案的定位是不同的，下面会详细解释。
