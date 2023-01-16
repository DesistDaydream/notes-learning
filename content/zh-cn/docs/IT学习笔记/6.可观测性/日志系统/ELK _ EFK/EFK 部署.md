---
title: EFK 部署
---

# 在 Kubernetes 中部署 EFK 套件

官方 addons：<https://github.com/kubernetes/kubernetes/tree/release-1.18/cluster/addons/fluentd-elasticsearch>

官方 addons 问题：<https://github.com/kubernetes/kubernetes/issues/94429>

参考：<https://www.qikqiak.com/post/install-efk-stack-on-k8s/>

es-service.yaml

```
kind: Service
apiVersion: v1
metadata:
  name: elasticsearch
  namespace: logging
  labels:
    app: elasticsearch
spec:
  selector:
    app: elasticsearch
  clusterIP: None
  ports:
    - port: 9200
      name: rest
    - port: 9300
      name: inter-node

```

es-statefulset.yaml

    apiVersion: apps/v1
    kind: StatefulSet
    metadata:
      name: es
      namespace: logging
    spec:
      serviceName: elasticsearch
      replicas: 3
      selector:
        matchLabels:
          app: elasticsearch
      template:
        metadata:
          labels:
            app: elasticsearch
        spec:
          containers:
          - name: elasticsearch
            image: docker.elastic.co/elasticsearch/elasticsearch:7.6.2
            ports:
            - name: rest
              containerPort: 9200
            - name: inter
              containerPort: 9300
            resources:
              limits:
                cpu: 1000m
              requests:
                cpu: 1000m
            volumeMounts:
            - name: data
              mountPath: /usr/share/elasticsearch/data
            env:
            - name: cluster.name
              value: k8s-logs
            - name: node.name
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: cluster.initial_master_nodes
              value: "es-0,es-1,es-2"
            - name: discovery.zen.minimum_master_nodes
              value: "2"
            - name: discovery.seed_hosts
              value: "elasticsearch"
            - name: ES_JAVA_OPTS
              value: "-Xms512m -Xmx512m"
            - name: network.host
              value: "0.0.0.0"
          volumes:
          - name: data
            emptyDir: {}

fluentd-configmap.yaml

    kind: ConfigMap
    apiVersion: v1
    metadata:
      name: fluentd-config
      namespace: logging
    data:
      system.conf: |-
        <system>
          root_dir /tmp/fluentd-buffers/
        </system>
      containers.input.conf: |-
        <source>
          @id fluentd-containers.log
          @type tail                              # Fluentd 内置的输入方式，其原理是不停地从源文件中获取新的日志。
          path /var/log/containers/*.log          # 挂载的服务器Docker容器日志地址
          pos_file /var/log/es-containers.log.pos
          tag raw.kubernetes.*                    # 设置日志标签
          read_from_head true
          <parse>                                 # 多行格式化成JSON
            @type multi_format                    # 使用 multi-format-parser 解析器插件
            <pattern>
              format json                         # JSON解析器
              time_key time                       # 指定事件时间的时间字段
              time_format %Y-%m-%dT%H:%M:%S.%NZ   # 时间格式
            </pattern>
            <pattern>
              format /^(?<time>.+) (?<stream>stdout|stderr) [^ ]* (?<log>.*)$/
              time_format %Y-%m-%dT%H:%M:%S.%N%:z
            </pattern>
          </parse>
        </source>
        # 在日志输出中检测异常，并将其作为一条日志转发
        # https://github.com/GoogleCloudPlatform/fluent-plugin-detect-exceptions
        <match raw.kubernetes.**>           # 匹配tag为raw.kubernetes.**日志信息
          @id raw.kubernetes
          @type detect_exceptions           # 使用detect-exceptions插件处理异常栈信息
          remove_tag_prefix raw             # 移除 raw 前缀
          message log
          stream stream
          multiline_flush_interval 5
          max_bytes 500000
          max_lines 1000
        </match>

        <filter **>  # 拼接日志
          @id filter_concat
          @type concat                # Fluentd Filter 插件，用于连接多个事件中分隔的多行日志。
          key message
          multiline_end_regexp /\n$/  # 以换行符“\n”拼接
          separator ""
        </filter>

        # 添加 Kubernetes metadata 数据
        <filter kubernetes.**>
          @id filter_kubernetes_metadata
          @type kubernetes_metadata
        </filter>

        # 修复 ES 中的 JSON 字段
        # 插件地址：https://github.com/repeatedly/fluent-plugin-multi-format-parser
        <filter kubernetes.**>
          @id filter_parser
          @type parser                # multi-format-parser多格式解析器插件
          key_name log                # 在要解析的记录中指定字段名称。
          reserve_data true           # 在解析结果中保留原始键值对。
          remove_key_name_field true  # key_name 解析成功后删除字段。
          <parse>
            @type multi_format
            <pattern>
              format json
            </pattern>
            <pattern>
              format none
            </pattern>
          </parse>
        </filter>

        # 删除一些多余的属性
        <filter kubernetes.**>
          @type record_transformer
          remove_keys $.docker.container_id,$.kubernetes.container_image_id,$.kubernetes.pod_id,$.kubernetes.namespace_id,$.kubernetes.master_url,$.kubernetes.labels.pod-template-hash
        </filter>

        # 只保留具有logging=true标签的Pod日志
        <filter kubernetes.**>
          @id filter_log
          @type grep
          <regexp>
            key $.kubernetes.labels.logging
            pattern ^true$
          </regexp>
        </filter>

      ###### 监听配置，一般用于日志聚合用 ######
      forward.input.conf: |-
        # 监听通过TCP发送的消息
        <source>
          @id forward
          @type forward
        </source>

      output.conf: |-
        <match **>
          @id elasticsearch
          @type elasticsearch
          @log_level info
          include_tag_key true
          host elasticsearch
          port 9200
          logstash_format true
          logstash_prefix k8s  # 设置 index 前缀为 k8s
          request_timeout    30s
          <buffer>
            @type file
            path /var/log/fluentd-buffers/kubernetes.system.buffer
            flush_mode interval
            retry_type exponential_backoff
            flush_thread_count 2
            flush_interval 5s
            retry_forever
            retry_max_interval 30
            chunk_limit_size 2M
            queue_limit_length 8
            overflow_action block
          </buffer>
        </match>

fluentd-daemonset.yaml

    apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: fluentd-es
      namespace: logging
      labels:
        k8s-app: fluentd-es
        kubernetes.io/cluster-service: "true"
        addonmanager.kubernetes.io/mode: Reconcile
    ---
    kind: ClusterRole
    apiVersion: rbac.authorization.k8s.io/v1
    metadata:
      name: fluentd-es
      labels:
        k8s-app: fluentd-es
        kubernetes.io/cluster-service: "true"
        addonmanager.kubernetes.io/mode: Reconcile
    rules:
    - apiGroups:
      - ""
      resources:
      - "namespaces"
      - "pods"
      verbs:
      - "get"
      - "watch"
      - "list"
    ---
    kind: ClusterRoleBinding
    apiVersion: rbac.authorization.k8s.io/v1
    metadata:
      name: fluentd-es
      labels:
        k8s-app: fluentd-es
        kubernetes.io/cluster-service: "true"
        addonmanager.kubernetes.io/mode: Reconcile
    subjects:
    - kind: ServiceAccount
      name: fluentd-es
      namespace: logging
      apiGroup: ""
    roleRef:
      kind: ClusterRole
      name: fluentd-es
      apiGroup: ""
    ---
    apiVersion: apps/v1
    kind: DaemonSet
    metadata:
      name: fluentd-es
      namespace: logging
      labels:
        k8s-app: fluentd-es
        kubernetes.io/cluster-service: "true"
        addonmanager.kubernetes.io/mode: Reconcile
    spec:
      selector:
        matchLabels:
          k8s-app: fluentd-es
      template:
        metadata:
          labels:
            k8s-app: fluentd-es
            kubernetes.io/cluster-service: "true"
          # 此注释确保如果节点被驱逐，fluentd不会被驱逐，支持关键的基于 pod 注释的优先级方案。
          annotations:
            scheduler.alpha.kubernetes.io/critical-pod: ''
        spec:
          serviceAccountName: fluentd-es
          containers:
          - name: fluentd-es
            image: quay.io/fluentd_elasticsearch/fluentd:v3.0.1
            env:
            - name: FLUENTD_ARGS
              value: --no-supervisor -q
            resources:
              limits:
                memory: 500Mi
              requests:
                cpu: 100m
                memory: 200Mi
            volumeMounts:
            - name: varlog
              mountPath: /var/log
            - name: varlibdockercontainers
              mountPath: /var/lib/docker/containers
              readOnly: true
            - name: config-volume
              mountPath: /etc/fluent/config.d
          tolerations:
          - operator: Exists
          terminationGracePeriodSeconds: 30
          volumes:
          - name: varlog
            hostPath:
              path: /var/log
          - name: varlibdockercontainers
            hostPath:
              path: /var/lib/docker/containers
          - name: config-volume
            configMap:
              name: fluentd-config

kibana.yaml

    apiVersion: v1
    kind: Service
    metadata:
      name: kibana
      namespace: logging
      labels:
        app: kibana
    spec:
      ports:
      - port: 5601
        nodePort: 30003
      type: NodePort
      selector:
        app: kibana

    ---
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: kibana
      namespace: logging
      labels:
        app: kibana
    spec:
      selector:
        matchLabels:
          app: kibana
      template:
        metadata:
          labels:
            app: kibana
        spec:
          containers:
          - name: kibana
            image: docker.elastic.co/kibana/kibana:7.6.2
            resources:
              limits:
                cpu: 1000m
              requests:
                cpu: 1000m
            env:
            - name: ELASTICSEARCH_HOSTS
              value: http://elasticsearch:9200
            ports:
            - containerPort: 5601

# 常见问题

## elasticsearch_logging_discovery.go:142] Found \[]

<https://github.com/kubernetes/kubernetes/issues/94429>
