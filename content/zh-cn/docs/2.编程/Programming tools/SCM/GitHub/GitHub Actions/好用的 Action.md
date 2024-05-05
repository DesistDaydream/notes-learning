---
title: 好用的 Action
linkTitle: 好用的 Action
date: 2024-04-22T22:43
weight: 20
---

# 概述

> 参考：
>
> - [GitHub Marketplace，Action](https://github.com/marketplace?type=actions)

# docker/build-push-action

https://github.com/docker/build-push-action

使用 Buildx 构建和推送 Docker 映像的 GitHub Action

```yaml
name: Build and push Docker image
uses: docker/build-push-action@v2
with:
- **context: .** # 构建上下文
- **file: simulate_mysql_exporter/e37_exporter/Dockerfile** # 指定要使用的 Dockerfile 路径，`默认值：{context字段的值}/Dockerfile`
- **push: true** # 构建完成后，是否推送镜像
- **tags: ghcr.io/desistdaydream/e37-exporter:v0.2.0** # 指定要构建的镜像名称
```

# action-gh-release

https://github.com/softprops/action-gh-release

用于创建 Realease 的 GitHub Action

> Notes: action-gh-release 需要上传文件，所以需要将仓库的 [Actions 配置](/docs/2.编程/Programming%20tools/SCM/GitHub/GitHub%20Actions/Actions%20配置.md) 中的 Workflow 权限设置为<font color="#ff0000">读写</font>

