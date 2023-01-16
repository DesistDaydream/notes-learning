---
title: Service Manifests 详解
---

# 概述

> 参考：
> - [API 文档单页](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#service-v1-core)
> - [API 参考-Service](https://kubernetes.io/docs/reference/kubernetes-api/service-resources/service-v1/)
> - [API 参考-Endpoint](https://kubernetes.io/docs/reference/kubernetes-api/service-resources/endpoints-v1/)

# Service Manifests 详解

## apiVersion: v1 # API 版本，基础字段必须要有

## kind: Service # 指明要创建的资源类型为 Service，基础字段必须要有

## metadata: # 指明该 Service 资源的元数据，其中 name 是必须要写明的元数据项

name: NAME # 指定该资源的名字

## sepc:  #指明该 Service 的规格(specification)

clusterIP: XXX.XXX.XXX.XXX #手动给该 Service 分配 IP，该 IP 在服务创建后无法手动修改，可以设置为 None，变成无头 service，这时候请求不由 service 处理，直接通过 service 名称转发到后端的 Pod
&#x20; ports:
&#x20; - protocol: TCP #将 service 的端口映射到 pod 的端口，使用 TCP 协议
&#x20;   nodePort: NUM #指明 Service 通过 k8s 集群中的那个端口对外提供服务，默认随机从 30000-32767 中随机分配(注：该字段只有 type 为 NodePort 的时候才有作用)
&#x20;   port: NUM #指明该 service 所使用的端口
&#x20;   targetPort: XXX #指明后端 Pod 的端口
&#x20; selector: #标签选择器 label selector 简称 selector，选择哪些 Pod 是该 Service 的后端
&#x20;   run: httpd #指明挑选那些 label 为 run:httpd 的 pod 作为该 service 的后端
&#x20; sessionAffinity: \<ClientIP|None> #设置会话亲和度，当为 None 的时候为同一个客户端的访问都会指向同一个 Pod，ClientIP 为进行负载调度
&#x20; type: TYPE

# Endpoints Manifests 详解

## apiVersion: v1

## kind: Endpoints

## metadata:

name: NAME #与 Endpoints 所关联的 Service 的 name 想同

## subsets: #指定子集

- addresses:
  &#x20; - ip: 10.10.100.101 #指定其中一个 endpoint 的 IP
  &#x20;   hostname: lch-test1 #指定该 endpoint 所在的主机的主机名
  &#x20; - ip: 10.10.100.102 #指定第二个 endpoint 的 IP
  &#x20;   hostname: lch-test2
  &#x20; ports:
  &#x20; - port: 9100 #指定 IP 所使用的 PORT

# Manifests 样例

    apiVersion: v1
    kind: Service
    metadata:
      labels:
        name: myapp
      name: myapp
    spec:
      ports:
      - name: http
        port: 80
        targetPort: 80
        nodePort: 30080
      type: NodePort
      selector:
        name: myapp

## service 绑定集群外部设备的 endpoints

需要手动添加 endpoints

    apiVersion: v1
    kind: Service
    metadata:
      name: external-metrics-service
      namespace: monitoring
      labels:
        prometheus: external-metrics
    spec:
      ports:
      - port: 9100

    ---
    apiVersion: v1
    kind: Endpoints
    metadata:
      name: external-metrics-service
      namespace: monitoring
      labels:
        prometheus: external-metrics
    subsets:
    - addresses:
      - ip: 10.10.100.101
        nodeName: lch-test
      - ip: 10.10.100.171
        nodeName: nfs-storage
      ports:
      - port: 9100
