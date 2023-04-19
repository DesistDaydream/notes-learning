---
title: Ingress Manifest 详解
---

# 概述

> 参考：
>
> - [API 文档,单页](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#ingress-v1-networking-k8s-io)
> - [官方文档,参考-Kubernetes API-服务资源-Ingress](https://kubernetes.io/docs/reference/kubernetes-api/service-resources/ingress-v1/)

# apiVersion: networking/v1

# kind: Ingress

# metadata

- **name:**STRING # Ingress 对象名称。必须名称空间唯一。
- **annotations:** # Ingress 的控制器将会根据 annotations，以自定义其行为。这些注释下的 kev/value 对 可以通过 ingress 传递给 controller ，然后 controller 根据这些信息进行更详细的配置，比如 url rewrite、代理超时时间等等。
  - 注意：不同的 controller 对 annotaions 中定义的内容有不同的处理。
    - [nginx ingress controller 社区版的 annotaions 说明](https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/)
    - nginx ingress controller 官方版的 annotaions 说明
  - 如果是在公有云上，公有云的各种 LB，也会读取 annotations 中的内容，以便将自家的 LB 与 ingress 关联

# spec

## defaultBackend: \<Object>

## ingressClassName: \<STRING>

## rules: <\[]Object> # 是一个类似于 nginx 的 7 层反向代理的配置

Ingress 资源最重要的字段，主要的实现逻辑都在这里了

- **host: \<STRING>** # (可省略，省略之后，可以使用 ip 来访问而不用使用域名)指定用户访问的域名，必须是 FQDN(完全限定域名)，不能是 IP，类似于 nginx 的 service 字段
- **http: \<Object>** #
  - **paths: <\[]Object>** # 用来定义，当用户访问不同资源时，把用户请求代理到不同的后端 Service，然后 Service 再把请求交给 Pod
    - **backend: \<Object> # 必须的**。定义要把哪些后端 Pod 的信息发送给 IngressController，如果是 IngressController 是 nginx，那么就是 upstream 的 IP 和 Port。可以通过两种方式来获取后端 Pod 的信息。
    - **resource: \<Object>** #
    - **service: \<Object>** # 指定后端类型为 service，就是通过 service 资源关联的 pod 来获取后端 pod。service 与 resource 字段互斥。
      - **name: \<STRING> # 必须的**。指明用于采集后端后端 Pod 信息的 service 的名称
      - **port: \<Object>** #
        - **number: \<INTEGER>**# 指定 service 上的端口号
  - **path \<STRING>** #
  - pathType

## tls: <\[]Object> # 启用 https 需要配置该字段

- **hosts: <\[]STRING>**#
- **secretName: \<STRING>** # 导入指定的 secret 对象内的数据。该 secret 中包括两个数据：证书和密钥

# 简单示例

```yaml
apiVersion: networking/v1
kind: Ingress
metadata:
  name: cafe-ingress
  annotations: # 将这些 annotations 添加到特定的 Ingress 对象，以自定义其行为。
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  tls: #启用https需要配置该字段
  - hosts:
    - cafe.example.com
    secretName: cafe-secret #把secret对象内的数据导入，该secret中包括两个数据：自签的证书和密钥，自签方式详见马哥的书《Kubernetes进阶实战》166页
  rules: #是一个类似于nginx的7层反向代理的配置
  - host: cafe.example.com #(可省略，省略之后，可以使用ip来访问而不用使用域名)指定用户访问的域名，必须是FQDN(完全限定域名)，不能是IP，类似于nginx的service字段
    http:
   paths: # 用来定义，当用户访问不同资源时，把用户请求代理到不同的后端Service，然后Service再把请求交给Pod
      #访问cafe.example.com/tea，则把请求代理到tea-svc这个service所拥有的后端pod上；访问cafe.example.com/coffee则把请求代理到名为coffee-svc这个service所拥有的后端pod上
   - path: /tea
     backend: #定义要把哪些后端Pod的IP信息发送给IngressController，如果是IngressController是nginx，那么这几个IP就是upstream的IP
       service:
            name: tea-svc # tes-svc 指明用于采集后端后端Pod信息的 service 的名称
         port:
              number: 80
   - path: /coffee
     backend:
       service:
            name: coffee-svc # coffee-svc 指明用于采集后端后端Pod信息的 service 的名称
         port:
              number: 80
  - host: bakery.example.com
    http:
   paths:
   - path: /
     .....
```

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: myapp
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
    - host: myapp.example.com
      http:
        paths:
          - path: /myapp
            backend:
              serviceName: myapp
              servicePort: 80
```
