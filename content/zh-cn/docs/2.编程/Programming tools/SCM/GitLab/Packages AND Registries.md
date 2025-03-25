---
title: Packages AND Registries
linkTitle: Packages AND Registries
weight: 20
---

# 概述

> 参考:
>
> - [官方文档，部署并发布你的应用 - Packages 和 Registries](https://docs.gitlab.com/ee/user/packages/package_registry/)

与 GitHub 不同，GitLab 有一个单独的存放 Assets 的地方，而不是放在 Release Assets 中，这个地方就是 Packages registry。

在 WebUI 点击左侧导航栏中 **Deploy > Package Registry** 标签可以查看项目中所有 Packages

GitLab 将多种实体或抽象概念抽象为 Packages

- [Package registry](#Package%20registry) # 储存各种二进制文件
- Container registry # 储存容器镜像
- etc.

一个项目产生的交付物都可以称为 Packages。可以在项目左侧导航栏点击 **Deploy > Package Registry** 查看项目的 Packages。

# Package registry

GitLab 将如下的实体抽象为 Package:

- **Generic(通用)** # 二进制文件
- **Helm** # Helm 文件
- etc.

# Generic package

> 参考:
>
> - [官方文档，通用 Packages](https://docs.gitlab.com/ee/user/packages/generic_packages)

使用 [GitLab CI](/docs/2.编程/Programming%20tools/SCM/GitLab/GitLab%20CI/GitLab%20CI.md) 推送通用的 Package

- https://docs.gitlab.com/ee/user/packages/generic_packages/#publish-a-generic-package-by-using-cicd

## Generic package 管理

> 参考:
>
> - [官方文档，减少 Package registry 存储用量](https://docs.gitlab.com/ee/user/packages/package_registry/reduce_package_registry_storage.html)

### 上传与下载 Package

```bash
export MY_TOKEN="XXXXXXXX"
export UPLOAD_FILE="node_exporter"
export UPLOAD_URL="${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/packages/generic/node_exporter/v1.0.0"
# CI_API_V4_URL 通常为 https://gitlab.example.com/api/v4
# CI_PROJECT_ID 可以从项目左侧导航栏点击 **Settings > General** 查看到
curl -k --header "PRIVATE-TOKEN: ${MY_TOKEN}" --upload-file bin/node_exporter "${UPLOAD_URL}"
# 下载时使用 UPLOAD_URL 即可匿名下载 Package 而不必更改项目的可见性。
```

> [!Tip] Package 的可见性
> 从项目左侧导航栏点击 **Setting > General > Visibility, project features, permissions** 打开 **Package registry > Allow anyone to pull from Package Registry** 的开关。
>
> <font color="#ff0000">此时可以匿名下载 Package</font>

> [!Notes] 相同 Package 的处理
> 当我们上传相同的 Package 时，GitLab 会保留所有 Package，并不会覆盖原始的包。使用 API 下载时，会下载最后上传的版本。
>
> 若想删除相同 Pacakage 的旧版本，可以配置自动删除策略或手动删除，参考下文[删除 Package](#删除%20Package)。

### 删除 Package

自动删除根据策略删除相同的 Package，手动删除则任意

手动删除有两种途径，利用 [API](https://docs.gitlab.com/ee/api/packages.html#delete-a-project-package) 或者在 Web UI 操作。

```bash
curl --request DELETE --header "PRIVATE-TOKEN: <your_access_token>" \
  "https://gitlab.example.com/api/v4/projects/:id/packages/:package_id"
```

要配置自动删除策略，点击左侧导航栏 **Settings > Packages and registries** 配置 **Manage storage used by package assets**，这里包含如下规则

- 清理重复的 asset。执行周期为 12 小时。可以配置保留多少个重复的 asset。主要用在
- Notes: 现阶段暂无其他规则
