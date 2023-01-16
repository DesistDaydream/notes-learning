---
title: Skaffold Pipeline 详解
---

# 概述

> 参考：
> - [官方文档](https://skaffold.dev/docs/design/config/)

Skaffold 默认通过 yaml 格式的名为 skaffold.yaml 的文件来指定 Pipelines 的行为。该文件应存放在项目目录的根目录中。当在项目的根目录中运行 skaffold 命令时，skaffold 将尝试从当前目录读取 skaffold.yaml 文件，来执行其内定义的 pipeline 任务。对于官方来说，skaffold.yaml 也称为 skaffold 的配置文件。

## Skaffold Pipeline 的工作流程

参考：[官方文档](https://skaffold.dev/docs/pipeline-stages/)

Skafflod 在执行 Pipeline 时，会经过 **5 个 stages(阶段)**。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/bxs72m/1616077584593-04a7412c-41f5-49e0-9eb5-bda9543222d7.jpeg)

> **artifacts(工件)** 大意就是指工具、物品等。比如 build artifacts(构建工件)，泛指可以执行构建功能的工具，比如 docker、Jib、Bazel 等等

运行 skaffold 时，skaffold 会在项目中收集源代码并使用配置文件中指定的构建工件来进行构建。构建成功后使用标记工件为镜像设置 tag，并推送到指定的镜像仓库中。在最后，skaffold 哈可以使用指定的部署工具(kubectl、helm 等)来将构建好的代码部署到 kubernetes 集群中。下面是各个阶段的详细介绍：

Init # generate a starting point for Skaffold configuration

Build # 根据指定的构建器来构建镜像

Tag # 根据指定的策略来为镜像设置 tag

Test #test images with structure tests

Deploy # 使用 kubectl、kustomize、helm 中的任意一种向 kubernetes 集群进行部署

File Sync #sync changed files directly to containers

Log Tailing #tail logs from workloads

Port Forwarding #forward ports from services and arbitrary resources to localhost

Cleanup # 清理执行 pipeline 中生成的 manifests 文件和 images。

# skaffold.yaml 详解

参考：[官方文档](https://skaffold.dev/docs/references/yaml/)

skaffold.yaml 文件由下面几个主要字段组成：

## apiVersion: skaffold/v2beta10 # 要使用的 Skaffold API 版本。当前的 API 版本是 skaffold/v2beta10。

## kind: Config # kind(类型) 始终为 Config。

## build # 指定 Skaffold 使用何种方式构建工件、标记工件以及推送工件来进行 Pipeline 任务。

Skaffold 支持使用本地 Docker 守护程序，Google Cloud Build，Kaniko 或 Bazel 来构建工件。

artifacts: <\[]Object> # 工件的信息

- image: STRING # 要构建的工件的镜像名
- docker
  - dockerfile

tagPolicy: # 确定如何为镜像添加 TAG。默认为：`gitCommit：{variant：Tags}`。

- envTemplate: # 使用可配置的模板字符串为镜像添加 TAG。
  - template: # 用于产生图像名称和标签。参阅 golang 的 [text/template](https://golang.org/pkg/text/template/)。针对当前环境执行模板，并注入这些变量。

local: # 描述如何在本地 docker 守护程序上进行构建以及如何选择推送到存储库

- useBuildkit: BOOL # 是否使用 BuildKit 构建 Docker 映像。`默认值：false`

## test # 指定 Skaffold 如何测试工件。

Skaffold 支持容器结构测试以测试构建的工件。有关更多信息，请参见测试人员。

## deploy # 指定 Skaffold 如何部署工件。

Skaffold 支持使用 kubectl，helm 或 kustomize 部署工件。有关更多信息，请参见部署者。

## profiles # 配置文件是一组设置，激活这些设置后，它们将覆盖当前配置。

您可以使用配置文件来覆盖 build，test 而 deploy 部分。
