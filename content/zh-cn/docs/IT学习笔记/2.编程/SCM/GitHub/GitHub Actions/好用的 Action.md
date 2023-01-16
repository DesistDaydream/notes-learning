---
title: 好用的 Action
---

# 概述

# [docker/build-push-action](https://github.com/docker/build-push-action)

使用 Buildx 构建和推送 Docker 映像的 GitHub Action

name: Build and push Docker image
uses: docker/build-push-action@v2
with:

- **context: .** # 构建上下文
- **file: simulate_mysql_exporter/e37_exporter/Dockerfile** # 指定要使用的 Dockerfile 路径，`默认值：{context字段的值}/Dockerfile`
- **push: true** # 构建完成后，是否推送镜像
- **tags: ghcr.io/desistdaydream/e37-exporter:v0.2.0** # 指定要构建的镜像名称
