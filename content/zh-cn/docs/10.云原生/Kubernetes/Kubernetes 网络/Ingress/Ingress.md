---
title: Ingress
weight: 1
---

# 概述

> 参考：
> 
> - [官方文档，概念-服务，负载均衡和网络-Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/)
> - 参考：<https://zhangguanzhang.github.io/2018/10/06/IngressController/>

Ingress 可以简单理解为是 Service 的 Service，是 Kubernetes 对“反向代理”概念的抽象。是一个专门给 kubernetes 用的 haproxy

举个例子，假如我现在有这样一个站点：<https://cafe.example.com>。其中，<https://cafe.example.com/coffee>，对应的是“咖啡点餐系统”。而，<https://cafe.example.com/tea>，对应的则是“茶水点餐系统”。这两个系统，分别由名叫 coffee 和 tea 这样两个 Deployment 来提供服务。

那么现在，我如何能使用 Kubernetes 来创建一个代理系统，从而实现当用户访问不同的域名时，能够访问到不同的 Deployment 呢？

上述功能，在 Kubernetes 里就需要通过 Ingress 对象来描述

Service 都是工作在 4 层模型上的，如果在 k8s 上的应用基于 https 来提供服务，那么在调度到 pod 上的时候就需要使用 7 层调度，这时候可以创建一个独特的 pod，略过 service，直接通过这个独特的 pod 进行反向代理把请求调度给用户，把 service 放在这个特殊的 pod 前端，但是这样经过的调度算法过多，导致性能过差；这时候可以把整个独特的 Pod 通过设置，把端口直接暴露，作为 node 上的一个进程来占用一个端口使用，然后通过 daemonset 给集群中某些需要的节点各自部署一个该 Pod(可以给一部分 node 加污点不让 pod 调度到此，并让该独特的 pod 容忍这个污点并调度上来)

这种独特的 Pod 统称 IngressController，由于 Pod 都是无状态的，随时可能会被摧毁后重建，这时候 HAProxy 和 Nginx 基于配置文件中 IP 地址的方式，在云环境下就没法用了。这时候可以创建一个 Service 关联上后端 Pod 与 Ingress Controller，该 Service 不做代理，仅作为分类来用，可以让 Ingress Controller 来正确找到自己所管理提供服务的 Pod。

为了让 Service 把后端 Pod 的信息能正常反馈给 Ingress Controller，k8s 有一种专门的对象叫做 Ingress，就是做这事的。Ingress 这个对象作为 Pod 与 IngressController 之间的桥梁，可以把 service 分好类的 Pod 识别出来，把 Pod 生成的 IP 地址变成配置信息注入到 Ingress Controller 中，这时 Ingress Controller 就可以动态变化自己所管理的后端 Pod。

# Ingress Controller

Ingress Controller 就是 Ingress 资源的控制器，用来让 Ingress 可以按照设定的状态来运行。同时可以当做具有类似 ngxin、haproxy 等的反向代理程序。

可以实现 IngressController 功能的工具有 HAProxy，Nginx，Traefik，Envoy 等等；HAProxy 是最不好用的，Nginx 还可以，Traefik 适应性很好，Envoy 为微服务而生

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/yq5gfx/1616117925361-31b9eb01-ac78-46e0-b300-1c179dfa9300.png)

工作模式：

1. IngressController 以 Pod 方式运行于 Kubernetes 集群之上，并时刻监听 Ingress 对象的状态
2. 创建完 Ingress 对象后，Ingress 会通过 `.spec.rules.host` 字段定义用户的访问入口，然后通过 `.spec.rules.http.paths.path` 字段进行 7 层转发
   1. host 字段类似 ngxin 配置里 ngx_http_core_module 模块的 service 字段 ，path 字段类似 nginx 配置里的 location,详见 2.nginx.note。
3. Ingress 对象会关联后端 Service，并把所有字段下的信息以及采集到的 Service 的信息传递给 IngressControlle
4. IngressController 会根据 Ingress 传递的信息更新自己的配置文件
   1. IngressController 的配置文件可以理解为 Nginx 或者 HAproxy 里的配置文件，所以可以把 IngressController 理解为 Nginx 或者 HAProxy。也可以直接把 Ingress 这个对象直接理解为一个配置文件，只不过这个配置文件是动态变化的。
5. client 通过 IngressController 所监听的端口进行访问，IngressController 根据从 Ingress 拿到的配置信息把用户的访问请求直接转发到信息中的 pod 上(这种转发会绕过 service，service 仅仅用于给 Ingress 提供其所关联的 endpoint 的 pod 的信息)
   1. IngressController 监听的是 80 端口，那么可以使用 Ingress 中 host 字段传递过来的 cafe.example.com 域名进行访问
      1. 比如：curl --resolve cafe.example.com:80:IngressControllerIP <http://cafe.example.com:80/tea>
   2. 如果访问到了在 Ingress 的 role 字段中未定义的域名，则一般会返回 404 错误，可以在 IngressController 的启动命令加一条 --default-backend-service=SERVICE 参数来设定一个当匹配域名失败的时候，会为用户返回定义的 SERVICE 下的 Pod 中的页面
6. 为了保证 IngressController 的高可用，可以创建多个 replicas，并在 IngressController 前面放置负载均衡设备，既然 IngressController 是以 Pod 形式出现的，那么就可以创建一个 Service 关联到这些 Pod 上来，或者直接使用其余的负载均衡功能。然后用户直接访问负载均衡的 VIP 即可。

## Ingress Controller 高可用

Ingress Controller 到集群内的路径这部分都有负载均衡了,我们比较关注部署了 Ingress Controller 后,外部到它这段的流量路径怎么高可用?

流量从入口到 Ingress Controller 的 Pod 有下面几种方式：

- type 为 LoadBalancer 的时候手写 externalIPs 很鸡肋,后面会再写文章去讲它
- type 为 LoadBalancer 的时候只有云厂商支持分配公网 ip 来负载均衡,LoadBalancer 公开的每项服务都将获得自己的 IP 地址,但是需要收费,自己建立集群想使用它的话得部署 metaLB。
- 不创建 svc,pod 直接用 hostport,效率等同于 hostNetwork,如果不代理四层端口还好,代理了的话每增加一个四成端口都需要修改 pod 的 template 来滚动更新来让 nginx bind 的四层端口能映射到宿主机上
- Nodeport,端口不是 web 端口(但是可以修改 Nodeport 的范围改成 web 端口),如果进来流量负载到 Nodeport 上可能某个流量路线到某个 node 上的时候因为 Ingress Controller 的 pod 不在这个 node 上,会走这个 node 的 kube-proxy 转发到 Ingress Controller 的 pod 上,多走一趟路
- 不创建 svc,效率最高,也能四层负载的时候不修改 pod 的 template,唯一要注意的是 hostNetwork: true 下 pod 会继承宿主机的网络协议,也就是使用了主机的 dns,会导致 svc 的请求直接走宿主机的上到公网的 dns 服务器而非集群里的 dns server,需要设置 pod 的 dnsPolicy: ClusterFirstWithHostNet 即可解决
- 可以使用 daemonset 的方式部署 IngressController 的 pod，然后配置 hostNetwork，直接让 IngressController 的 pod 的 ip 直接暴露在所在 node 上，然后前端配置外部负载均衡设备负载 ingressController 所在 node，不用 service 资源来关联 pod，以便节省网络开销。
- 当然，如果 node 很多，也可以通过 label 让 IngressController 运行在指定的 node 上，然后只负载指定的几个 node。

```yaml
---
spec:
  containers:
    - xxx
  dnsPolicy: ClusterFirstWithHostNet
  hostNetwork: true
```

已经部署的 deploy 的话使用 patch 修改为 hostNetwork

```bash
kubectl -n ingress-nginx \
  patch deploy nginx-ingress-controller \
  -p '{"spec":{"template":{"spec":{"dnsPolicy":"ClusterFirstWithHostNet","hostNetwork":true}}}}'
```

部署方式没多大区别开心就好

- daemonSet + nodeSeletor
- deploy 设置 replicas 数量 + nodeSeletor + pod 互斥

client 通过域名，最后导向到 SLB，负载到三个 node 上，三个 node 都部署了 hostNetwork 的 ingress controller，然后 ingress controller 根据 server_name 反向代理相关的 svc 的 http 服务

- 所以可以一个 vip 飘在拥有存活的 controller 的宿主机上,云上的话就用 slb 来负载代替 vip,自己有条件有 F5 之类的硬件 LB 一样可以代替 VIP
- 最后说说域名请求指向它,如果部署在内网或者办公室啥的,内网有 dns server 的话把 ing 的域名全部解析到 ingress controller 的宿主机 ip(或者 VIP,LB 的 ip)上,否则要有人访问每个人设置/etc/hosts 才能把域名解析来贼麻烦,如果没有 dns server 可以跑一个 external-dns,它的上游 dns 是公网的 dns 服务器,办公网内机器的 dns server 指向它即可,云上的话把域名请求解析到对应 ip 即可
- traefik 和 ingress nginx 类似,不过它用 go 实现的并且好像它不支持四层代理,如果上微服务可以上 istio,没接触过它,不知道原理是否如此
- ingress nginx 的 log 里会一直刷找不到 ingress-nginx 的 svc 不处理的话会狂刷 log 导致机器 load 过高,创建一个同名的 svc 即可解决,例如创建一个不带选择器 clusterip 为 null 的
- get ing 输出的时候 ADDRESS 一栏会为空，ingress-nginx 加参数--report-node-internal-ip-address 即可解决
- 使用了 rancher 在负载均衡也就是 ingress 页面，ingress 状态不为 Active 的话在 ingress-nginx 的参数配置--publish-service 和--publish-status-address

## 多 Ingress Controllers

当然，一个集群可以多组 ingress controller，或者不同的 ingress controller。拿 ingress nginx 举例。

假如公司的内网组了 vpn，办公网和机房的 node 的网络打通，我们就需要两组 ingress controller 了(不一定需要一样，例如 ingress nginx 和 traefik)，例如下面

- 硬件 F5 对外暴露了一组业务的 ingress 给用户
- 而 it 的研发听说 k8s 部署方便，想把服务部署到 k8s 上，然后这些服务想对办公网的同事让 http 访问，压力不高。我们使用低成本的 keepalived 飘在内网组的 node 上
- 两组 deploy 使用 nodeSelector 固定在两组不同的 label 的 node 上

```bash
                                            +---------------------------+
                                            |                           |
                               -+---------->+ hostNetwork的ingress nginx|
                              /             +---------------------------+ node1
                             /
                            /
                           /                +---------------------------+
+--------+                /                 |                           |
| client +----------->    F5 -------------->+ hostNetwork的ingress nginx|
+--------+                \                 +---------------------------+ node1
                           \
                            \
                             \              +---------------------------+
                              ------------->+hostNetwork的ingress nginx |
                                            |                           |
                                            +---------------------------+ node3

                                            +---------------------------+
                                            |                           |
                               -+---------->+ hostNetwork的ingress nginx|
                              /             +---------------------------+ node4
                             /
                            /
                           /                +---------------------------+
+--------+                /                 |                           |
| client +-----------> SLB(or VIP) -------->+ hostNetwork的ingress nginx|
+--------+                \                 +---------------------------+ node5
                           \
                            \
                             \              +---------------------------+
                              ------------->+ hostNetwork的ingress nginx|
                                            |                           |
                                            +---------------------------+ node6
```

我们需要一个 svc 供内网服务，只需要创建该 ingress 的时候添加一个 annotation

    apiVersion: networking.k8s.io/v1beta1
    kind: Ingress
    metadata:
      name: it-work
      annotations:
        kubernetes.io/ingress.class: "nginx-internal"
    spec:
      tls:
      - secretName: tls-secret
      rules:
      - http:
          paths:
          - backend:
              serviceName: it-work-svc
              servicePort: 8080

而两组的 ingress controller 的该选项的值不同即可

    $ docker run --rm --entrypoint /nginx-ingress-controller \
      quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.26.1 --help |& grep -A2 -- --ingress-class
          --ingress-class string
    Name of the ingress class this controller satisfies.
    The class of an Ingress object is set using the annotation "kubernetes.io/ingress.class".
    All ingress classes are satisfied if this parameter is left empty.

这个参数内网的那组我们就写 --ingress-class=nginx-internal，面向公网的不能写和这个值一样，定义成另一个即可

# Ingress Class

Ingress 可以由不同的控制器实现，每个 Ingress 对象都应该通过 `spec.ingressClassName` 字段指定一个控制器

说白了，就是为 Ingress 分类，不通类的 Ingress 可以被不同的 Ingress Controller 引用。

Ingress Class 也是一个资源，当我们部署 Ingress Controller 时，程序通常都会指定一个想要生成的 IngressClass。

比如 Nginx Ingress Controller，通过 `--controller-class` 命令行标志来生成一个名为 k8s.io/ingress-nginx 的 nginx 类的 IngressClass，效果如下：

```bash
~]# kubectl get ingressclasses.networking.k8s.io -A
NAME    CONTROLLER             PARAMETERS   AGE
nginx   k8s.io/ingress-nginx   <none>       114m
```
