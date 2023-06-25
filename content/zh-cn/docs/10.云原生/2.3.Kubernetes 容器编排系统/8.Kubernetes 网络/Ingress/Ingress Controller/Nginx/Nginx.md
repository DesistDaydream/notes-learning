---
title: Nginx
---

# 概述

分为两个版本

- K8S 社区版 Nginx Ingress Controller # <https://github.com/kubernetes/ingress-nginx>
- Nginx 官网版 Nginx Ingress Controller # <https://github.com/nginxinc/kubernetes-ingress>

# 部署

版本支持矩阵: https://github.com/kubernetes/ingress-nginx#supported-versions-table

## k8s 社区版部署方式

> 参考：
>
> - [官方文档，部署-安装指南-裸金属集群](https://kubernetes.github.io/ingress-nginx/deploy/#bare-metal-clusters)(就是通过纯 Manifests 文件部署)
> - [官方文档，部署-安装指南-](https://kubernetes.github.io/ingress-nginx/deploy/#using-helm)快速开始(直接就是 Helm 安装)

注意：

- 从 v1.0.0 版本开始，仅支持 Kubernetes 版本 >= v1.19 ，因为从 v1.0.0 版本开始，删除了对 `networking.k8s.io/v1beta` 资源的支持。
  - 详见：[公众号-CNCF，更新 NGINX-Ingress 以使用稳定的 Ingress API](https://mp.weixin.qq.com/s/hVTWlfrqmjZRrb0KTsDrZA)
- 从 v1.3.0 版本开始，仅支持 Kubernetes 版本 >= v1.20。
  - [详见：公众号-MoeLove，K8S 生态周报| Kubernetes Ingress-NGINX 功能冻结前最后一个版本发布](https://mp.weixin.qq.com/s/7vOTDqpi4tg-AEzP_YAWAQ)
  - 为了能兼容 Kubernetes 的更高版本，所以我们将 controller 中用于选举的机制修改成了使用 Lease API 的方式，而不再是原先的 configmap 的方式。其实在 Kubernetes Ingress-NGINX v1.3.0 版本中，我增加了往 Lease API 平滑迁移的逻辑，在使用 v1.3.0 版本的时候，可以自动的完成 ConfigMap 往 Lease API 迁移的逻辑。 **所以，如果是想要从旧版本进行平滑升级，建议先升级到 v1.3.0，待完成自动的迁移后，再往更新的版本升级**。

其中[Bare-metal 段落](https://kubernetes.github.io/ingress-nginx/deploy/#bare-metal)为通过 yaml 文件部署方式

- 部署一个 Nginx 的 IngressController
- Nginx-IngressController 的配置详见：<https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/>

### 使用 helm 部署

- helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
- helm pull ingress-nginx/ingress-nginx
- tar -zxvf ingress-nginx-X.X.X.tgz
- cd ingress-nginx
- helm install nginx --namespace ingress-controller --create-namesapce .

部署完成后，即可创建 Ingress 对象，关联后端 Service

- 配置好域名解析，直接使用域名访问即可，不用域名会出现问题。(如果 ingress 中没有配置 host 字段，则无需解析)
  - 因为在 《[HTTP](/docs/4.数据通信/通信协议/7.HTTP/7.HTTP.md)》 文章中，http 协议的 Request-URL 中包含了 client 访问时输入的网址，可以是 IP 或者域名，而 ingress 对象的配置中，host 配置的都是域名，如果是 ip 的话， ingress 是无法识别出带有 IP 的 URL 的。

## Nginx 官方版部署方式

- <https://docs.nginx.com/nginx-ingress-controller/installation/installation-with-manifests/>
  - 直接创建一个 daemonset 类型的 nginx-ingress(就是 ingress-controller)
- 创建 Ingress 对象，关联后端 Service
- 配置好域名解析，直接使用域名访问即可。(如果 ingress 中没有配置 host 字段，则无需解析)
  - Note：不用域名会出现问题。因为在 《[HTTP](/docs/4.数据通信/通信协议/7.HTTP/7.HTTP.md)》 文章中，http 协议的 Request-URL 中包含了 client 访问时输入的网址，可以是 IP 或者域名，而 ingress 对象的配置中，host 配置的都是域名，如果是 ip 的话， ingress 是无法识别出带有 IP 的 URL 的。

# 关联文件与配置

可以通过如下几种方式来配置 Nginx Ingress Controller

- **Command Line Flags** # Controller 程序的命令行标志，可以定义程序自身的运行时行为。
- **ConfigMap**# 通过 `--configmap` 命令行标志指定的名称空间下的 ConfigMap 资源，仅支持一个 ConfigMap 对象。
  - Nginx Ingress Controller 会读取该 ConfigMap 对象下的内容，并与 Custom template 一起生成 Nginx 的 nginx.conf 文件
- **Annotations**# 通过为 Ingress 对象添加 Annotations 字段下的内容来定义 Nginx Ingress Controller 的运行时行为。
- **Custom template** # Nginx Ingress Controller 会使用模板文件生成 nginx.conf 文件。模板文件可以在 [GitHub 代码](https://github.com/kubernetes/ingress-nginx/blob/master/rootfs/etc/nginx/template/nginx.tmpl)中找到

**ConfigMap、Annotations、Custom template 都可以用来定义 nginx.conf 这个配置文件**，不同的是：

- Annotations 主要用于定义每一个 Virtual Server 中的指令(比如 server{}、location{} 中的指令)
- 而 ConfigMap 更多是定义 http{}、stream{} 这种配置环境中的指令。
- Custom template 则是生成 ngxin.conf 的模板文件，Annotations 与 ConfigMap 中定义的内容，都会通过值传递的方式，传递到模板文件中，然后生成 ngxin.conf 文件。

[**/etc/nginx/template/nginx.tmpl**](https://github.com/kubernetes/ingress-nginx/blob/main/rootfs/etc/nginx/template/nginx.tmpl) # Nginx 的 nginx.conf 文件的模板文件

# Grafana 面板

## 12559

https://grafana.com/grafana/dashboards/12559

该面板曾经名称为：Loki v2 Web Analytics Dashboard，在 Grafana 官方面板首页最顶部推荐了很久；现在更名为：

Grafana Loki Dashboard for NGINX Service Mesh
