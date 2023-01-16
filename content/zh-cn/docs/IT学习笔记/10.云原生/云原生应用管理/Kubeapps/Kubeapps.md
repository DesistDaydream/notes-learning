---
title: Kubeapps
---

# 概述

> 参考：
> - [官网](https://kubeapps.com/)
> - [GitHub](https://github.com/kubeapps/kubeapps)

Kubeapps 是一个基于 Web 的 UI，用于在 Kubernetes 集群中部署和管理应用程序。

Kubeapps 可以实现下述功能：

- 从 Chart 仓库浏览和部署 Helm
- 检查、升级、删除集群中已经安装的基于 Helm 的应用程序。
  - 可以看到，Kubeapps 与原生 Helm 结合比较紧密，这一点是 Rancher 做不到的。
- 添加自定义和自由 Chart 仓库
- 浏览和部署 Kubernetes Operator
- 使用 OAuth2/OIDC 提供程序对 Kubeapps 进行安全身份验证
- 基于 Kubernetes RBAC 的安全授权

## Kubeapps 组件

Kubeapps 抽象了一个 **Asset(资产)** 的概念，Asset 是多种事物的集合，比如一个 Chart 仓库就属于一个资产。

一个完整的 Kubeapps 服务，通常包含如下组件：

- **Apprepository-controller** # 应用仓库管理
- **Asset-syncer** # 扫描 Helm 仓库，并在 PostgreSQL 中填充 Chart 元数据的工具，然后 Assetsvc 组件将会提供这些元数据
- **Assetsvc** # 暴露 API，用于访问 PostgreSQL 中的 Chart 仓库中的 Chart 元数据
- **Dashboard** # Web UI
- **Kubeops** # 暴露 API 来访问 Helm API 和 Kubernetes 资源
- **PostgreSQL** # 存储 Chart 仓库的信息，其他组件都是无状态的

# Kubeapps 权限管理

> 参考：
> - [GitHub,docs-user-access-control.md](https://github.com/kubeapps/kubeapps/blob/master/docs/user/access-control.md)

Kubeapps 会创建如下几个 ClusterRole 对象，以便进行权限管理：

- kubeapps:controller:kubeops-ns-discovery-kubeapps # 查看名称空间
- kubeapps:controller:kubeops-operators-kubeapps # 对 packagemanifests/icon 自定义资源的查看权限
- APP 仓库管理权限
  - kubeapps:kubeapps:apprepositories-read # 查看仓库
  - kubeapps:kubeapps:apprepositories-refresh # 更新仓库
  - kubeapps:kubeapps:apprepositories-write # 增删改查仓库，即对 apprepositories.kubeapps.com 这个自定义资源的全部权限

## 权限管理示例

```bash
kubectl create -n kubeapps rolebinding lch-kubeapps-repositories-write  \
--clusterrole=kubeapps:kubeapps:apprepositories-write \
--serviceaccount user-sa-manage:lch
```

在 kubeapps 名称空间中创建一个 RoleBinding，将 `kubeapps:kubeapps:apprepositories-write` 这个角色与 user-sa-manage 名称空间下的 lch 这个 ServiceAccount 绑定。这样，使用 lch 的 token 登录时，就可以在全局名称空间(kubeapps)中操作镜像仓库了。

> 2.3.2 版本有个小 BUG，从 其他名称空间切换的 kubeapps 名称空间后，会提示没权限，需要刷新一下才行。
