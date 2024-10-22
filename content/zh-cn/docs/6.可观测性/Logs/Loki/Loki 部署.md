---
title: Loki 部署
linkTitle: Loki 部署
date: 2022-10-23T13:44:00
weight: 2
---

# 概述

> 参考：
>
> - [官方文档，安装](https://grafana.com/docs/loki/latest/installation/)
> - [官方文档，基础知识-架构-部署模式](https://grafana.com/docs/loki/latest/fundamentals/architecture/deployment-modes)

# 使用 docker 运行 Loki

```bash
docker run -d --rm --name loki \
  --network host \
  -v /opt/loki/config:/etc/loki \
  -v /opt/loki/data:/loki \
  -v /etc/localtime:/etc/localtime:ro \
  grafana/loki
```

注意：与 Prometheus 类似，需要修改 /opt/loki 目录权限为 777，否则 pod 内进程对该目录无操作权限

# 在 Kubernets 集群中部署

添加 loki 的 helm chart 仓库

- **helm repo add grafana https://grafana.github.io/helm-charts**
- **helm repo update**

## Helm 部署 Loki 套件

> 参考：
> - [官方文档 2.4.x，安装-helm](https://grafana.com/docs/loki/v2.4.x/installation/helm/)

部署 Loki 栈

- kubectl create ns loki # 创建名称空间
- helm pull grafana/loki-stack # 获取 loki-stack 的 charts 压缩包
- tar -zxvf loki-stack-X.XX.X.tgz # 解压 charts
- cd loki-stack # 进入目录，根据需求修改模板或 values.yaml 文件
- helm upgrade --install loki --namespace=loki . # 使用默认配置在 loki 名称空间中部署 loki 栈 。该方式会部署 loki 与 promtail

在 grafana 中添加 loki 数据源，如图所示

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/vg0v2e/1616129749320-bb4bc4c9-2acb-460f-a655-5ff76766eb24.jpeg)

## Helm 部署 Simple scalable 架构 Loki

部署 Loki

- **helm install -n logging loki --create-namespace grafana/loki-simple-scalable**

注意：可扩展模式会部署 `loki-gateway` 用以接收请求，并分离 读/写 请求，所有 Promtail 用来向 Loki 发起写请求的采集客户端和 Grafana 这种用来向 Loki 发起读请求的展示客户端，指定 Loki 端点时，都要指定 `loki-gateway`。

这里为什么要自带 grafana-agent？这个 grafana-agent 是通过 grafana/agent-operator 拉起来的。

## Helm 部署 Microservices 架构 Loki

部署 Loki

- **helm install -n logging loki --create-namespace grafana/loki-distributed**

# 日志测试容器

这俩容器会频繁刷新各种类型日志

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dummylogs
spec:
  replicas: 3
  selector:
    matchLabels:
      app: dummylogs
  template:
    metadata:
      labels:
        app: dummylogs
        logging: "true" # 要采集日志需要加上该标签
    spec:
      containers:
        - name: dummy
          image: cnych/dummylogs:latest
          args:
            - msg-processor
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dummylogs2
spec:
  replicas: 3
  selector:
    matchLabels:
      app: dummylogs2
  template:
    metadata:
      labels:
        app: dummylogs2
        logging: "true" # 要采集日志需要加上该标签
    spec:
      containers:
        - name: dummy
          image: cnych/dummylogs:latest
          args:
            - msg-receiver-api
```

# 其他

> 参考；
> 
> - [GitHub 项目，grafana/loki，production 目录](https://github.com/grafana/loki/tree/main/production)
> - [公众号，Loki 生产环境集群方案](https://mp.weixin.qq.com/s/qnt7JUzHLUU6Qs_tv5V0Hw)

