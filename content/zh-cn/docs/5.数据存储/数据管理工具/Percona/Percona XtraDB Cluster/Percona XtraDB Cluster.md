---
title: "Percona XtraDB Cluster"
linkTitle: "Percona XtraDB Cluster"
weight: 1
---

# 概述

> 参考：
> - [GitHub 项目，percona/percona-xtradb-cluster](https://github.com/percona/percona-xtradb-cluster)
> - [官网](https://www.percona.com/software/mysql-database/percona-xtradb-cluster)


# 部署 PXC

> 参考：
> - [官方文档，安装 PXC 集群](https://docs.percona.com/percona-xtradb-cluster/latest/install/index.html)


## 使用 Docker 部署

> 参考：
> - [Running Percona XtraDB Cluster in a Docker Container - Percona XtraDB Cluster](https://docs.percona.com/percona-xtradb-cluster/8.0/install/docker.html#docker)


## 使用 PXC Operator 在 Kubernetes 中部署

> 参考：
> - [官方文档，快速开始指南-使用 kubectl 安装 PXC Operator](https://docs.percona.com/percona-operator-for-mysql/pxc/kubectl.html)
> - [官方文档，高级安装指南-通用 Kubernetes 安装](https://docs.percona.com/percona-operator-for-mysql/pxc/kubernetes.html)

### 快速体验

部署 Operator

kubectl apply -f https://raw.githubusercontent.com/percona/percona-xtradb-cluster-operator/v1.12.0/deploy/bundle.yaml

部署 PXC

kubectl apply -f https://raw.githubusercontent.com/percona/percona-xtradb-cluster-operator/v1.12.0/deploy/cr.yaml

### 高级安装

创建名称空间

kubectl create namespace pxc

kubectl apply -f crd.yaml

kubectl apply -n pxc -f rbac.yaml

kubectl apply -n pxc -f operator.yaml

kubectl apply -n pxc -f secrets.yaml

kubectl apply -n pxc -f cr.yaml

### manifests 配置详解

> 参考：
> - [官方文档，参考-CR选项](https://docs.percona.com/percona-operator-for-mysql/pxc/operator.html)

#### HAProxy 部分配置

haproxy.readinessProbes