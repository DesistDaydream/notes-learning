---
title: Prometheus-adapter
linkTitle: Prometheus-adapter
weight: 20
---

# 概述

> 参考:
>
> - [GitHub 项目，kubernetes-sigs/prometheus-adapter](https://github.com/DirectXMan12/k8s-prometheus-adapter)
>   - 该项目从 DirectXMan12/k8s-prometheus-adapter 移动到 kubernetes-sigs/prometheus-adapter

**主要特性：**
**一、adapter 在成功注册 API 之后，可以通过 Prometheus 实现 custom.metrics.k8s.io API 和 metrics.k8s.io API 的功能**
adaper 可以替换掉 metrics server 来实现其功能。adapter 要想实现 kubectl top node/pod 命令的功能，则需要 adapter 通过查询 Prometheus 来获取数据完成，这需要 prometheus 提前获取某些数据来支撑 adapter 得查询，而查询语句则是根据 adapter 的配置文件中 resourceRules 配置环境中的规则来指定。

- 其中 kubectl top node 如果查询语句查询结果为空，则在执行命令查询时会报错：error: metrics not available yet
- 其中 kubectl top pod 如果查询语句查询结果为空，则在执行命令查询时会报错：No resources found

**二、adapter 可以根据 prometheus 提供的核心 metrics 数据(比如 CPU 使用率等)或者自定义 metrics 数据，来自动实现**[**HPA**](4.Controller(控制器).md 容器编排系统/4.Controller(控制器).md)**功能。**

> HPA 的概念详见《Controller 控制器》章节中的 HPA 控制器介绍

## adapter 工作流程

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rd9psx/1616068726920-7b7a454e-7853-49a7-b56e-f8fcc594ad0e.png)
一、注册指标。adpater 启动后根据配置中的 seriesQuery 和 name 字段的规则，将匹配到的 metrics 注册到 Metrics API 中
二、查询 metrics 的值并返回。当 hpa 或者 kubectl top 命令想要获取信息时，会对 Metrics API 发起请求，请求交给 adapter 进行处理。adapter 根据 resources 和 metricsQuery 字段的规则，向 prometheus 发起查询请求，并将结果返回给 hap 或者 kubectl top 命令。

adapter 不像 metrics server 直接与 kubelet 交互，然后从 kubelet 的 cAdvisor 中获取 Core Metrics 和其中的值。而是与 prometheus 交互，通过 PromQL 查询语句来获取想要的 Metrics 和 Metircs 的值。通过配置文件中的规则，发现可以使用的 Core Metrics 或者 Custom Metrics，并将其注册到 Metrics API 中，而不像 metrics server 可以直接将 cpu 和 memory 的 metrics 注册到 Metrics API 。

## apapter 处理 MetricsAPI 接收到的请求的方式

比如一个 MetricsAPI 请求是这样的：/apis/custom.metrics.k8s.io/v1beta1/namespaces/monitoring/pods/grafana-5c55845445-q2p9l/http_request_per_second，adapter 会提取其中的字段，将其分为三个部分 MetricsName、Resource、Resource 的 Objects。Resource 概念详见：Kubernetes API 介绍。Object 概念详见：Kubernets Object 对象。

- MetricsName 值为 http_request_per_second
- Resource 值为 namesapces 和 pods
- Resource 的 Object 值分为两部分
  - namespaces 资源的 object 值为 monitoring
  - pods 资源的 object 值为 grafana-5c55845445-q2p9l

adapter 将这三个内容填充到配置文件中 metricsQuery 关键字定义的 Go 模板中，生成 PromQL，并向 prometheus 发起查询。

# Adapter 配置

> - 官方文档：
>   - 配置文件说明：<https://github.com/DirectXMan12/k8s-prometheus-adapter/blob/master/docs/config.md>
>   - 配置文件示例：<https://github.com/DirectXMan12/k8s-prometheus-adapter/blob/master/docs/config-walkthrough.md>
>   - 官方默认的配置文件：<https://github.com/DirectXMan12/k8s-prometheus-adapter/blob/master/deploy/manifests/custom-metrics-config-map.yaml>

Adapter 的配置文件，是用来通过其中定义的 rules 来确定公开哪些 metrics ，以及如何公开它们。还有通过对 prometheus 查询获取样本值

Note：
如果配置文件有问题，adapter 无法获取到 metrics 并注册到 api，则该 api 会报错：FailedDiscoveryCheck
在配置文件里会经常看到 series 这个单词，表示序列的意思，应该时 time-series 的简称，也有 metrics name 的意思

    [root@master adapter]# kubectl get apiservices.apiregistration.k8s.io
    NAME                                   SERVICE                         AVAILABLE                      AGE
    ......
    v1beta1.custom.metrics.k8s.io          monitoring/prometheus-adapter   False (FailedDiscoveryCheck)   3s
    ......

配置文件有两大部分

1. rules 配置环境 # 用于 Custom Metrics
2. resourceRules 配置环境 # 用于 Core Metrics

## rules 配置环境

    rules:
    # 指定 PromQL 从 Prometheus 查找时间序列，然后 adapter 会剥离获取到的时间序列中的 label ，只保留 metrics name ，然后根据配置文件中 name 字段重命名之后，将新命名的 metrics name 注册到 Metrics API 中
    # 对 seriesQuery 中获取到的时间序列进行过滤，可以使用 is 或者 isNot 关键字
    - seriesQuery: 'PromQL'
      seriesFilters:
      - is: <RegEx> #保留 RegEx 中匹配到的时间序列
        isNot: <RegEx> #丢弃 RegEx 中匹配到的时间序列
      # 关联 metrics 中的 label 与 k8s 资源
      resources:
        overrides:
          # LabelName为填到 metricsQuery 中 PromQL 里的标签名
          # ResourceName为 k8s 中的资源，根据对Metrics API 请求的url，找到对应的资源，并将资源具体的对象填入 metricsQuery 中 PromQL 的标签值
          LabelName: {[group: GroupName,]resource: "ResourceName"} # 核心组可以省略组名
        template: TEMPLATE
      # 根据 matches 关键字中的正则对发现中获取的 metrics 重命名(如果不需要重命名，可以省略该配置段)。
      name:
        matches: "RegEx"
        as: "$1" # 可省略，默认值为 $1 。i.e.保留 metches 字段的正则中第一个 () 中的内容。也可以显式得添加 as 字段，来指定重命名结果
      # 该字段使用 Go模板，adapter 处理 MetricsAPI 的请求，(处理方式详见上文)
      # 处理结果的 MetricsName、Resource、Objects 这三个内容，会填充到模板中，生成 PromQL，然后向 prometheus 发起查询。
      metricsQuery: "sum(rate(<.Series>>{<<.LabelMatchers>>,container!="POD"}[2m])) by (<.GroupBy>)"

rules 大致可以分为四个部分，下面对四个部分的关键字进行说明：

一、Discovery(发现) # 根据 rules 中定义的 Prometheus QL 查找 metircs 。

1. seriesQuery(PromQL) # 查询 Prometheus 的语句，通过这个查询语句查询到的所有指标都可以用于 HPA
2. seriesFilters \[] # 查询到的指标可能会存在不需要的，可以通过它过滤掉。
3. is(RegEx) #
4. isNot(RegEx) #

二、Association(关联) # 关联 metrics 的 label 与 k8s resource。adapter 向 prometheus 发起查询时，会将关联规则中指定的名字作为标签名。adapter 接收到 MetricsAPI 请求中的 k8s object 作为标签值，将两者填充到 PromQL 中

1. resources # 通过 seriesQuery 查询到的只是指标，如果需要查询某个 Pod 的指标，肯定要将它的名称和所在的命名空间作为指标的标签进行查询， resources 就是将指标的标签和 k8s 的资源类型关联起来，最常用的就是 pod 和 namespace。有两种添加标签的方式，一种是 overrides，另一种是 template。
2. overrides # 它会将指标中的标签和 k8s 资源关联起来。上面示例中就是将指标中的 pod 和 namespace 标签和 k8s 中的 pod 和 namespace 关联起来，因为 pod 和 namespace 都属于核心 api 组，所以不需要指定 api 组。当我们查询某个 pod 的指标时，它会自动将 pod 的名称和名称空间作为标签加入到查询条件中。比如 nginx:{group:"apps",resource:"deployment"} 这么写表示的就是将指标中 nginx 这个标签和 apps 这个 api 组中的 deployment 资源关联起来；
   1. LabelName: {\[group: Group,]resource: "RESOURCE"}
3. template # 通过 go 模板的形式。比如 template:"kube\_<<.Group>>\_<<.Resource>>" 这么写表示，假如 <<.Group>> 为 apps， <<.Resource>> 为 deployment，那么它就是将指标中 kube_apps_deployment 标签和 deployment 资源关联起来。

三、Naming(命名) #(如果不需要重命名配置可省略)根据表达式匹配规则，把发现配置中查找到的 metrics 重命名。
比如某些以 total 结尾的指标，这些指标拿来做 HPA 是没有意义的，我们需要对这些指标进行速率计算，比如这种语句 sum(rate(http_requests_total{}\[2m])) by (pod,namespace) ，在进行计算后，使用 total 来命名没意义了，需要赋予一个新的名字来表示。

1. name # 用来给指标重命名的，之所以要给指标重命名是因为有些指标是只增的，比如以 total 结尾的指标。这些指标拿来做 HPA 是没有意义的，我们一般计算它的速率，以速率作为值，那么此时的名称就不能以 total 结尾了，所以要进行重命名。
2. matches(RegEx) # 通过正则表达式来匹配指标名，可以进行分组
3. as(STRING) # 默认值为$1。也就是第一个分组。 as 为空就是使用默认值的意思。

四、Querying(查询) # adapter 在向 prometheus 查询数据时，根据该规则发送 PromQL 。

1. metricsQuery(GoTemplate) # 模板中，adapter 处理 MetricsAPI 的请求转换后的三个部分，在模板中变为以下字段。(这就是 Prometheus 的查询语句了，前面的 seriesQuery 查询是获得 HPA 指标。当我们要查某个指标的值时就要通过它指定的查询语句进行了。可以看到查询语句使用了速率和分组，这就是解决上面提到的只增指标的问题。)
   1. <<.Series>> # PromQL 的指标名。根据 Discovery 和 Naming 配置部分结合获取
   2. <<.LabelMatchers>> # PromQL 的标签集合。根据 Association 配置部分获取。Name 与 Value 的对应关系是根据 Association 配置部分的规则实现的。Name 为 Association 配置段中的 LabelName；Value 为 Resource 的 Object。
   3. <<.GroupBy>> # 以逗号分隔的标签名的集合，该值是 .LabelMatchers 中的标签名集合
   4. 这三个字段的内容，需要根据 MetircsAPI 收到的请求，以及上面三个配置部分的规则，综合起来获取

### 总结

比如我现在做了如下配置
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rd9psx/1616068726798-f8716c85-b38c-4bc8-acd3-10d934933312.png)
当对 Metrics API 的请求为 kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1/namespaces/monitoring/pods/\*/http_request_per_second" 时，这个请求会被 adapter 拆分成如下三部分

1. MetricsName：http_requests_total
2. Resource：namesapces,pods
3. Object 分两部分
   1. namespaces 的 Object：monitoring
   2. pods 的 Object：grafana 和 prometheus #(Note:这里假定该时间序列的 monitoring 名称空间里两个 pods，prometheus 和 grafana)

这时，如果想要生成 PromQL 语句向 prometheus 发起查询，需要根据配置文件的 Association 和 Naming 配置段的内容来向 metricsQuery 中的模板填入信息

1. <<.Series>>：http_requests_total # 因为 MetricsAPI 请求中的指标名(http_request_per_second)是重命名之后的，所以 adapter 会根据这个指标名，找到配置文件中 name 字段中匹配到这个指标名的 rule 配置环境，将其中 seriesQuery 字段中的原始指标名填入模板 .Series 中
2. <<.LabelMatchers>>：namespace="monitoring",pod=~"grafana|prometheus"
   1. LabelMatchers 中的内容是根据 Association 配置段的内容取得的，namespace 和 pod 就是的'标签名'。
   2. 根据绑定规则，‘namesapce 和 pod 标签’分别绑定了‘namespaces 和 pods 资源’，那么根据 MetricsAPI 中的 url，获取到指定资源的 object，会作为对应标签的值。
   3. namespace 标签对应的 k8s 资源是 namesapces ，MetricsAPI 请求中该资源的 object 是 monitoring ，则 monitoring 就是 LabelValue
   4. pod 标签对应的 k8s 资源是 pods ，MetricsAPI 请求中该资源的 object 是全部对象，也就是说 grafana 和 prometheus ，则 grafana|prometheus 就是 LabelValue
3. <<.GroupBy>>：pod

那么，该模板会被转换为这样的 PromQL：sum(rate(http_requests_total{namespace="monitoring", pod=~"grafana|prometheus"}\[2m])) by (namesapce,pod)

Note:如果在关联时，没有关联 namespace ，那么在向 Metrics API 发起上述相同请求时会报错 Error from server (NotFound): the server could not find the metric http_request_per_second for pods。因为没有关联，所以 adapter 也就无法处理 url 中 namesapce 的信息。

配置文件中 Naming 加 Querying 的组合有点类似于 prometheus 的 Recording Rules ，Naming 是 record 字段，Querying 是 expr 字段。等于是通过一个表达式，生成了一个新的 Metrics 名称

### 配置示例

    rules:
    - seriesQuery: 'http_request_total{namespace!="",pod!=""}'
      resources:
        overrides:
          namespace: {resource: "namespace"}
          pod: {resource: "pod"}
      name:
        matches: "^(.*)_total"
        as: "${1}_per_second"
      metricsQuery: 'sum(rate(http_request_total{<<.LabelMatchers>>}[2m])) by (<.GroupBy>)'

通过该示例配置， adapter 会向 Metrics API 注册如下两个 metrics，namespaces/http_request_per_second 和 pods/http_request_per_second

    [root@master-1 custom]# kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1/" | jq .
    {
      "kind": "APIResourceList",
      "apiVersion": "v1",
      "groupVersion": "custom.metrics.k8s.io/v1beta1",
      "resources": [
        {
          "name": "namespaces/http_request_per_second",
          "singularName": "",
          # 此字段指明该 metric 是否需要在发起url请求时指定具体的namespace
          # 因为这里是namespaces资源，所以可能会有误解，如果此处node(i.e."name": "nodes/http_request_per_second")资源，则很好理解了，node是不需要namesacpe的
          "namespaced": false,
          "kind": "MetricValueList",
          "verbs": [
            "get"
          ]
        },
        {
          "name": "pods/http_request_per_second",
          "singularName": "",
          "namespaced": true,
          "kind": "MetricValueList",
          "verbs": [
            "get"
          ]
        }
      ]
    }

通过 MetricsAPI 获取指标 http_request_per_second 值时，会返回如下内容

    [root@master-1 custom]# kubectl get --raw /apis/custom.metrics.k8s.io/v1beta1/namespaces/monitoring/pods/*/http_request_per_second | jq .
    {
      "kind": "MetricValueList",
      "apiVersion": "custom.metrics.k8s.io/v1beta1",
      "metadata": {
        "selfLink": "/apis/custom.metrics.k8s.io/v1beta1/namespaces/monitoring/pods/%2A/http_request_per_second"
      },
      "items": [
        {
          "describedObject": {
            "kind": "Pod",
            "namespace": "monitoring",
            "name": "grafana-5c55845445-q2p9l",
            "apiVersion": "/v1"
          },
          "metricName": "http_request_per_second",
          "timestamp": "2020-06-18T07:47:26Z",
          "value": "9m",
          "selector": null
        }
      ]
    }

## resourceRules 配置环境

示例如下：

    apiVersion: v1
    data:
      config.yaml: |-
        resourceRules:
          cpu:
            # 该 PromQL 用于 kubectl top pod 命令
            containerQuery: sum(rate(container_cpu_usage_seconds_total{<<.LabelMatchers>>}[1m])) by (<.GroupBy>)
            # 该 PromQL 用于 kubectl top node 命令
            nodeQuery: sum(rate(container_cpu_usage_seconds_total{<<.LabelMatchers>>, id='/'}[1m])) by (<.GroupBy>)
            resources:
              overrides:
                node:
                  resource: node
                namespace:
                  resource: namespace
                pod:
                  resource: pod
            containerLabel: container
          memory:
            containerQuery: sum(container_memory_working_set_bytes{<<.LabelMatchers>>}) by (<.GroupBy>)
            nodeQuery: sum(container_memory_working_set_bytes{<<.LabelMatchers>>,id='/'}) by (<.GroupBy>)
            resources:
              overrides:
                node:
                  resource: node
                namespace:
                  resource: namespace
                pod:
                  resource: pod
            containerLabel: container
          window: 1m
    kind: ConfigMap
    metadata:
      name: adapter-config
      namespace: monitoring

# 问题示例

如果部署完成后，无法通过 Metrics API 获取指标，则可能的原因有以下几点

1. api 无法关联到后端 adapter
2. adapter 的 config 文件配置有问题
