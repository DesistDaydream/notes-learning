---
title: Traefik
---

# Traefik 介绍

traefik 在 2.0 之后的版本，不再使用 ingress 资源作为入口，而是一个名为 IngressRoute 的 CRD 资源来作为入口

# Traefik 安装

k8s 的 yaml 文件：<https://docs.traefik.io/user-guides/crd-acme/>

# Yaml 样例

```
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: simpleingressroute

spec:
  entryPoints:
    - web
  routes:
  - match: Host(`your.example.com`) && PathPrefix(`/notls`) #Host可省略，这样就跟ingress的host字段省略、istio的VirtualService资源的hosts字段写*，是一个效果
    kind: Rule
    services:
    - name: whoami
      port: 80

```

https 的 yaml 配置

    apiVersion: traefik.containo.us/v1alpha1
    kind: IngressRoute
    metadata:
      name: ingressroutetls

    spec:
      entryPoints:
        - websecure
      routes:
      - match: Host(`your.example.com`) && PathPrefix(`/tls`)
        kind: Rule
        services:
        - name: whoami
          port: 80
      tls:
        certResolver: myresolver
