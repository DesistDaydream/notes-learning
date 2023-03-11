---
title: Prometheus Manifest 详解
---

# 概述

> 参考：
>
> - [官方文档](https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/api.md#prometheus)

# apiVersion: monitoring.coreos.com/v1

# kind: Prometheus

# metadata

# spec

## additionalScrapeConfigs: \<Object> # 额外的抓取配置

该字段可以通过 [additional 功能](docs/IT学习笔记/6.可观测性/监控系统/Prometheus/Prometheus%20衍生品/Prometheus%20Operator/CR%20详解/Prometheus/Prometheus.md#additionalScrapeConfigs) 为 Prometheus Server 创建额外的 Scrape 配置。这种方式常用来为 Prometheus Server 创建静态的 Scrape 配置。

- **key: \<STRING>** # 要引用的 secret 对象中 .data 字段下，指定的 key 的值
- **name: \<STRING>** # 要使用的 secret 对象名称

## containers: <\[]Object> # 注入其他容器或修改 Operator 生成的容器

这可用于允许将身份验证代理添加到 Prometheus pod 或更改 Operator 生成的容器的行为。 如果此处描述的容器共享相同的名称，则它们将修改操作员生成的容器，并且通过战略合并补丁进行修改。 当前的容器名称为：“ prometheus”，“ config-reloader”和“ thanos-sidecar”。 覆盖容器完全不在维护人员支持的范围之内，因此，您接受此行为可能随时中断，恕不另行通知。

- **name: \<STRING>** # 指定要修改的容器名。支持的容器名为：prometheus、config-reloader、thanos-sidecar。若指定的名称不存在，则创建新的容器

## externalUrl: \<STRING> # 为 prometheus Server 指定 --web.external-url 命令行标志的值

## logFormat: \<STRING> # 为 Prometheus Server 指定 --log.format 命令行标志的值

## logLevel: \<STRING> # 为 Prometheus Server 指定 --log.level 命令行标志的值

## probeNamespaceSelector: Object> # 选择指定名称空间下的 Probe 资源

效果与 serviceMonitorNamespaceSelector 字段一样，只不过是与 probeSelector 字段配合使用。
**matchExpressions: <\[]Object>**#
**matchLabels: \<map\[string]string>** #

## probeSelector: Object> # 通过 Probe 资源发现待采集目标

为 Prometheus Server 发现想要抓取指标的目标。效果与 serviceMonitorSelector 字段一样，只不过是发现 Probe 资源。
**matchExpressions: <\[]Object>**#
**matchLabels: \<map\[string]string>** #

## resources: # 与 pod 资源下的同名字段功能一样。Note:该字段下的内容仅对 prometheus 容器生效

## retention: STRING> # 为 Prometheus Server 指定 --storage.tsdb.retention.time 命令行标志

## serviceMonitorNamespaceSelector: Object> # 选择指定名称空间下的 ServiceMonitoring 资源

通过[标签选择器](Label%20and%20Selector(标签和选择器).md 容器编排系统/1.API、Resource(资源)、Object(对象)/Label and Selector(标签和选择器).md)，匹配出指定的名称空间，该名称空间将会被 serviceMonitorSelector 字段使用，serviceMonitorSelector 将会从匹配到的名称空间中发现 ServiceMonitor 资源。
若该字段值为 `nil`，则仅从 Prometheus 对象所在名称空间中发现 ServiceMonitor 资源。
**matchExpressions: <\[]Object>**#
**matchLabels: \<map\[string]string>** #

## serviceMonitorSelector: Object> # 通过 ServiceMonitor 资源发现待采集目标

为 Prometheus Server 发现想要抓取指标的目标。

通过[标签选择器](Label%20and%20Selector(标签和选择器).md 容器编排系统/1.API、Resource(资源)、Object(对象)/Label and Selector(标签和选择器).md)，匹配出指定的 ServiceMontior 资源。serviceMonitorSelector 会从 probeNamespaceSelector 字段定义的名称空间中，查找 Service Monitor 资源，并获取其中的信息，以便转换为 Prometheus Server 的配置文件中 scrape_config 字段的内容。

若该字段值为 `{}`，则发现所有 ServiceMonitor 资源。否则可以根据匹配规则，选择指定的 ServiceMonitor。
**matchExpressions: <\[]Object>** # matchExpressions is a list of label selector requirements. The requirements are ANDed.
**matchLabels: \<map\[STRING]STRING>** # matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.

## storage: Object> # 定义 Prometheus 的存储方式

## volumeMounts: <\[]Object> # 与 pod 资源下的同名字段功能一样

用于指定 volume 的挂载路径。Note：该字段内容只对 prometheus 容器生效

## volmues: <\[]Object> # 与 pod 资源下的同名字段功能一样。用于指定一个 volume
