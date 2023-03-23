---
title: 各 CRD 的 yaml 样例
---

# 应用实例

## prometheus 使用 storageclass 实现持久存储

可以在 prometheus 这个 CRD 的定义中找到关于 storage 的字段<https://github.com/coreos/kube-prometheus/blob/master/manifests/0prometheus-operator-0prometheusCustomResourceDefinition.yaml#L3633>

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/yifyaq/1616068754603-7b67500b-0560-447e-a3c4-bedcdf7e8d28.png)

注解中写到如果不指定 storage 的话，则默认使用 emptydir 类型作为存放监控数据的 volume 类型

如果想要修改成持久存储，则只需要在声明 prometheus 资源的 yaml 文件中加入如下 storage 字段即可，下面的实例是让 prometheus 使用名为 managed-nfs-storage 的 StorageClass

      storage:
        volumeClaimTemplate:
          spec:
            storageClassName: managed-nfs-storage
            resources:
              requests:
                storage: 10Gi

当加入该字段后，prometheus 资源生成的 statefulset 就会多出来一个字段，如图所示。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/yifyaq/1616068754580-603a0098-31fc-47b1-b008-911c582efdef.png)

## prometheus 使用 hostPath 实现持久存储

如下示例，使用本地 pod 所在节点的本地目录/root/prometheus-k8s-db 来作为数据存储目录

注意，最好使用 nodeSelector 让 pod 始终调度到同一个节点。

    apiVersion: monitoring.coreos.com/v1
    kind: Prometheus
    metadata:
      name: test
    spec:
      replicas: 2
      containers:
      - name: prometheus
        volumeMounts:
        - mountPath: /prometheus
          name: prometheus-k8s-db
      volumes:
      - name: prometheus-k8s-db
        hostPath:
          path: /root/prometheus-k8s-db
      nodeSelector:
        monitor: prometheus

## 修改 prometheus 的启动参数

可以在 prometheus 这个 CRD 的定义中找到关于 container 的字段<https://github.com/coreos/kube-prometheus/blob/master/manifests/0prometheus-operator-0prometheusCustomResourceDefinition.yaml#L806>

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/yifyaq/1616068754607-380f9402-777a-46f0-a081-4e61120fb81b.png)

注解中写到 container 字段可以在生成 statefulset 的时候注入额外的 container 或者修改自动生成的 container(比如修改 args、volumemount 字段等)，

注意：

1. 在修改已经存在的容器时，需要指定要修改的 container 的 name
2. 如果要修改某个字段，需要全部重新填写(e.g.修改 arg 字段，默认有 7 个 arg，如果只是修改其中一个，那么在修改 yaml 的时候，所有的 arg 都要填上，否则，最后生成的 statefulset 就会只有 1 个 arg)。

下面是修改 storage.tsdb.retention.time 这个参数以便让 prometheus 可以保存数据更久的时间的样例(默认是 24 小时，我现在想保存 7 天)

      containers:
      - args:
        - --web.console.templates=/etc/prometheus/consoles
        - --web.console.libraries=/etc/prometheus/console_libraries
        - --config.file=/etc/prometheus/config_out/prometheus.env.yaml
        - --storage.tsdb.path=/prometheus
        - --storage.tsdb.retention.time=7d
        - --web.enable-lifecycle
        - --storage.tsdb.no-lockfile
        - --web.route-prefix=/
        name: prometheus

加入上述字段后，在 prometheus 资源生成的 statefulset 就会多出来上述内容
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/yifyaq/1616068754606-245c3379-04db-42b6-84d9-cd03b06bd28f.png)

## 监控集群外部设备

创建一个关联外部设备的同名的 service 和 endpoint，然后再创建关联该 service 的 servicemonitor 即可，service 样例详见 service 后端绑定集群外部设备.note，下面是 servicemonitor 的样例

其中 sourceLabel 和 targetLabel 可以将默认的 instance 名修改为 endpoint 的 nodeName，默认的是 IP:PORT 格式，可以在定义 endpoint 的时候，给每个 ip 配上一个 nodeName

    apiVersion: monitoring.coreos.com/v1
    kind: ServiceMonitor
    metadata:
      name: external-metrics
      namespace: monitoring
      labels:
        prometheus: external-metrics
    spec:
      endpoints:
      - interval: 15s
        relabelings:
        - action: replace
          regex: (.*)
          replacment: $1
          sourceLabels:
          - __meta_kubernetes_endpoint_hostname
          targetLabel: instance
      selector:
        matchLabels:
          prometheus: external-metrics
    ---
    apiVersion: v1
    kind: Service
    metadata:
      name: external-metrics
      namespace: monitoring
      labels:
        prometheus: external-metrics
    spec:
      ports:
      - port: 9100
        nodePort: 30005 #测试用，暴露出来看看能不能获取到metric的信息
      type: NodePort

    ---
    apiVersion: v1
    kind: Endpoints
    metadata:
      name: external-metrics
      namespace: monitoring
      labels:
        prometheus: external-metrics
    subsets:
    - addresses:
      - ip: 10.10.100.101
        hostname: lch-test
      - ip: 10.10.100.171
        hostname: nfs-storage
      ports:
      - port: 9100

# prometheus 资源的配置样例

    apiVersion: monitoring.coreos.com/v1
    kind: Prometheus
    metadata:
      labels:
        app: k8s-prometheus
      name: monitor-bj-net-k8s-prometheus
      namespace: monitoring
    spec:
      alerting:
        alertmanagers:
        - apiVersion: v2
          name: monitor-bj-net-k8s-alertmanager
          namespace: monitoring
          pathPrefix: /
          port: web
      externalUrl: http://prometheus.desistdaydream.ltd/
      image: quay.io/prometheus/prometheus:v2.22.1
      logFormat: logfmt
      logLevel: info
      portName: web
      replicas: 1
      resources:
        limits:
          cpu: "2"
          memory: 2Gi
        requests:
          cpu: 500m
          memory: 400Mi
      retention: 10d
      routePrefix: /
      ruleSelector:
        matchLabels:
          app: k8s
          release: monitor-bj-net
      securityContext:
        fsGroup: 2000
        runAsGroup: 2000
        runAsNonRoot: true
        runAsUser: 1000
      serviceAccountName: monitor-bj-net-k8s-prometheus
      storage:
        volumeClaimTemplate:
          metadata:
            name: prometheus
          spec:
            resources:
              requests:
                storage: 10Gi
            storageClassName: managed-nfs-storage
      version: v2.22.1
