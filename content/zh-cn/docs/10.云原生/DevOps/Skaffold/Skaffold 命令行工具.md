---
title: Skaffold 命令行工具
---

# 概述

> 参考：
> - [官方文档](https://skaffold.dev/docs/references/cli/)

Skaffold 命令行工具提供以下命令

End-to-end pipelines:

- [skaffold run](https://skaffold.dev/docs/references/cli/#skaffold-run) # 构建并运行一次
- [skaffold dev](https://skaffold.dev/docs/references/cli/#skaffold-dev) # 进入`dev`模式，Skaffold 将监视应用程序的源文件，并在检测到更改时将重建您的映像（或将文件同步到正在运行的容器中），推送所有新映像并将该应用程序重新部署到群集中。
- [skaffold debug](https://skaffold.dev/docs/references/cli/#skaffold-debug) - to run a pipeline in debug mode

Pipeline building blocks for CI/CD:

- [skaffold build](https://skaffold.dev/docs/references/cli/#skaffold-build) - to just build and tag your image(s)
- [skaffold deploy](https://skaffold.dev/docs/references/cli/#skaffold-deploy) - to deploy the given image(s)
- [skaffold delete](https://skaffold.dev/docs/references/cli/#skaffold-delete) - to cleanup the deployed artifacts
- [skaffold render](https://skaffold.dev/docs/references/cli/#skaffold-render) - build and tag images, and output templated Kubernetes manifests

Getting started with a new project:

- [skaffold init](https://skaffold.dev/docs/references/cli/#skaffold-init) # 初始化一个 skaffold.yaml 文件
- [skaffold fix](https://skaffold.dev/docs/references/cli/#skaffold-fix) - to upgrade from

Other Commands:

- [skaffold help](https://skaffold.dev/docs/references/cli/#skaffold-help) - print help
- [skaffold version](https://skaffold.dev/docs/references/cli/#skaffold-version) - get Skaffold version
- [skaffold completion](https://skaffold.dev/docs/references/cli/#skaffold-completion) - setup tab completion for the CLI
- [skaffold config](https://skaffold.dev/docs/references/cli/#skaffold-config) - manage context specific parameters
- [skaffold credits](https://skaffold.dev/docs/references/cli/#skaffold-credits) - export third party notices to given path (./skaffold-credits by default)
- [skaffold diagnose](https://skaffold.dev/docs/references/cli/#skaffold-diagnose) - diagnostics of Skaffold works in your project
- [skaffold schema](https://skaffold.dev/docs/references/cli/#skaffold-schema) - list and print json schemas used to validate skaffold.yaml configuration

# skaffold dev \[FLAGS]

FLAGS:

- --cache-artifacts=BOOL # 是否使用工件的缓存。`默认值：true`
- --no-prun=BOOL # 是否清理由 Skaffold 购机爱你的镜像和容器。`默认值：false`
- --render-only=BOOL # 仅打印渲染后的 manifests 文件，而不部署。

EXAMPLE

# skaffold run \[FLAGS]
