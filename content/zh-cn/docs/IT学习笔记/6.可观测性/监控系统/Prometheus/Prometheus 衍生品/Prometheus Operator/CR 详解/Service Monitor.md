---
title: Service Monitor
---

# Service Monitor 介绍

注意：ServiceMonitor 资源本身无法直接为目标 job 添加 label，所有 label 只能从关联的 Service 中获取，然后再通过 ServiceMonitor 资源的 spec.endpoints.relabelings 字段(就是使用 Prometheus 的 relabel 功能)，将获取到的 label 改为自己想要的

operator 根据 ServiceMonitor 的定义自动生成 Prometheus 配置文件中的 scrape 配置段中的内容。创建一个 SM，就代表要给 prometheus 配置中 scrape 配置段中加入内容。

ServiceMonitor 资源描述了 Prometheus Server 的 Target 列表，Operator 会监听这个资源的变化来动态的更新 Prometheus Server 的 Scrape Targets 并让 Prometheus Server 去 reload 配置。而该资源主要通过 Selector 根据 Labels 选取对应 Service 的 endpoints，并让 Prometheus Server 通过 Service 进行拉取 Metrics,Metrics 信息要在 http 的 url 输出符合 metrics 格式的信息,ServiceMonitor 也可以定义目标的 metrics 的 url.

## ServiceMonitor 是如何自动生成 Prometheus Server 配置文件的 scrape 段的呢？

ServiceMonitor 关联好 service 后，选取 service 对应的 endpoints 作为 target。然后自动生成的 scrape 配置中，job 的名字是以关联的 service 名来命名；targets 就是该 service 所关联的每一个 endpoint；每个 target 的 instance 一般情况都是 endpoint 的 IP:PORT

下面是一个获取 k8s 系统资源 controller-manager 的 metrics 的样例

    apiVersion: monitoring.coreos.com/v1
    kind: ServiceMonitor
    metadata:
      labels:
        k8s-app: kube-controller-manager
      name: kube-controller-manager
      namespace: monitor
    spec:
    # 指定从哪个namesapce中关联service。默认情况是从prometheus资源所在的namespace中关联service
      namespaceSelector:
        matchNames:
        - kube-system
    # 如果想要关联所有namespace中的service，则不用进行matchNames，使用any: true即可
        # any: true
    # 指定要匹配的service的label
      selector:
        matchLabels:
          k8s-app: kube-controller-manager
    # 选择endpoints中指定的端口获取metrics
      endpoints:
      - port: http-metrics
        path: "/snmp" #指定从该endpoints的哪个路径获取metrics，默认路径为/metrics

声明完成后，还需要修改 prometheus 的 yaml 文件，以便让其可以匹配指定的 ServiceMonitor，否则默认 prometheus 资源是不匹配任何 ServiceMonitor 的，在 spec 键下添加字段如下：

    # 指定prometheus匹配哪个名称空间的ServiceMonitor。{}表示匹配所有名称空间的ServiceMonitor
    serviceMonitorNamespaceSelector: {}
    # 指定prometheus选择ServiceMonitor时的label。{}表示匹配所有label的ServiceMonitor。
    serviceMonitorSelector: {}

现在 prometheus 已经与 ServiceMonitor 关联上了，等待 operator 将 ServiceMonitor 获取的配置处理之后，添加进 prometheus 并更新配置。则可以在配置文件的 scrape_configs 字段看到新的配置，说明 prometheus 已经正确获取了 ServiceMonitor 传递的信息。

但是此时会有一个问题就是，prometheus 已经有了 scrape 的配置，知道了 target，但是在 kubernetes 中，prometheus 作为一个 pod，target 也是一个 pod，这就等于是一个 pod 要去 get 另一个 pod 的信息，这明显是 RBAC 和 SA 相关的事宜(概念详见 7.AuthenticationAndAuthorization.note)，并且现在也还没进行任何 RBAC 相关的配置，所以虽然 prometheus 有配置，但是无法获取 target 上的任何 metrics 信息。如果查看 prometheus 的日志，会看到如下类似的报错

    level=error ts=2019-09-19T15:44:30.735Z caller=klog.go:94 component=k8s_client_runtime func=ErrorDepth msg="/app/discovery/kubernetes/kubernetes.go:265: Failed to list *v1.Pod: pods is forbidden: User \"system:serviceaccount:monitor:default\" cannot list resource \"pods\" in API group \"\" in the namespace \"kube-system\""
    level=error ts=2019-09-19T15:44:30.735Z caller=klog.go:94 component=k8s_client_runtime func=ErrorDepth msg="/app/discovery/kubernetes/kubernetes.go:264: Failed to list *v1.Service: services is forbidden: User \"system:serviceaccount:monitor:default\" cannot list resource \"services\" in API group \"\" in the namespace \"kube-system\""
    level=error ts=2019-09-19T15:44:31.735Z caller=klog.go:94 component=k8s_client_runtime func=ErrorDepth msg="/app/discovery/kubernetes/kubernetes.go:263: Failed to list *v1.Endpoints: endpoints is forbidden: User \"system:serviceaccount:monitor:default\" cannot list resource \"endpoints\" in API group \"\" in the namespace \"kube-system\""

报错解读：endpoints is forbidden: User "system:serviceaccount:monitor:default" cannot list resource "endpoints" in API group "" in the namespace "nginx-ingress""

这一段的大意就是在 monitor 这个 namespace 下的 default 这个用户(i.e.ServiceAccount)，不能对 kube-system 这个 namespace 里的 endpoints 资源使用 list 命令。也就是说 prometheus 这个 pod 里的进程，想要对 endpoints 执行 list 命令，但是被禁止了

其原因就是由于 prometheus 资源的 pod 默认使用的是其所在 namespace 下的 default 这个 ServiceAccount，这个 SA 不具备任何可以操作 pod 的权限，所以需要自己创建一个 SA 和具有特定操作权限的 ClusterRole，再将二者绑定，则 prometheus 的 pod 就可以获取其余 pod 上的信息了，示例如下：

    apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: prometheus
      namespace: monitor
    ---
    apiVersion: rbac.authorization.k8s.io/v1beta1
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
    apiVersion: rbac.authorization.k8s.io/v1beta1
    kind: ClusterRoleBinding
    metadata:
      name: prometheus
    roleRef:
      apiGroup: rbac.authorization.k8s.io
      kind: ClusterRole
      name: prometheus
    subjects:
    - kind: ServiceAccount
      name: prometheus
      namespace: monitor

然后此时又要修改 prometheus 的 yaml 文件了，因为要让 prometheus 使用自己创建的 SA，所以还需要在 spec 键下中加入如下内容

    serviceAccountName: prometheus

这时 operator 才算真正完成了相关任务，此时查看 prometheus，可以看到采集到 kube-controller-manager 的 metrics 了

# ServiceMonitor CRD Manifest 详解

# ServiceMonitor 样例

## snmp_exporter 样例

    apiVersion: monitoring.coreos.com/v1
    kind: ServiceMonitor
    metadata:
      name: snmp-metrics
      namespace: monitoring
      labels:
        prometheus: snmp-metrics
    spec:
      jobLabel: snmp
      selector:
        matchLabels:
          prometheus: snmp-metrics
      endpoints:
      - interval: 30s
        scrapeTimeout: 120s
        port: snmp
        params:
          module:
          - if_mib
          target:
          - 10.10.100.254
        targetPort: 9116
        path: "/snmp"
        relabelings:
        - action: replace
          sourceLabels:
          - __param_target
          targetLabel: instance

    ---
    kind: Service
    apiVersion: v1
    metadata:
      name: snmp-metrics
      namespace: monitoring
      labels:
        prometheus: snmp-metrics
    spec:
      ports:
      - port: 9116
        nodePort: 30006
        name: snmp
      type: NodePort
    ---
    apiVersion: v1
    kind: Endpoints
    metadata:
      name: snmp-metrics
      namespace: monitoring
      labels:
        prometheus: snmp-metrics
    subsets:
    - addresses:
      - ip: 10.10.100.12
        hostname: storage-1
      ports:
      - port: 9116
        name: snmp
