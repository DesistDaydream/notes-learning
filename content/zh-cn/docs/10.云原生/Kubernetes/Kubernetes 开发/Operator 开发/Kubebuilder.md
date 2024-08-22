---
title: Kubebuilder
linkTitle: Kubebuilder
date: 2024-08-22T09:10
weight: 20
---

# 概述

> 参考：
>
> -

Kubebuilder 代码示例详见 [GitHub 上我的 kubernetesAPI 仓库](https://github.com/DesistDaydream/kubernetesAPI/tree/master/operator)

# kubebuilder 命令行工具

用于构建 Kubernetes 扩展和工具的开发工具包。 提供用于创建新项目，API 和控制器的库和工具。 包括用于将工件打包到安装程序容器中的工具。

典型的项目生命周期：

- 初始化项目：

  - kubebuilder init --domain example.com --license apache2 --owner "The Kubernetes authors"

- 创建一个或多个新资源 API 并将代码添加到其中：

  - kubebuilder create api --group \<group> --version \<version> --kind \<Kind>

Create resource will prompt the user for if it should scaffold the Resource and / or Controller. To only scaffold a Controller for an existing Resource, select "n" for Resource. To only define the schema for a Resource without writing a Controller, select "n" for Controller.

After the scaffold is written, api will run make on the project.

**kubebuilder \[COMMAND]**

Available Commands:

## create Scaffold a Kubernetes API or webhook

## edit This command will edit the project configuration

## init Initialize a new project

Initialize a new project including vendor / directory and Go package directories.

Writes the following files:

- a boilerplate license file
- a PROJECT file with the project configuration
- a Makefile to build the project
- a go.mod with project dependencies
- a Kustomization.yaml for customizating manifests
- a Patch file for customizing image for manager manifests
- a Patch file for enabling prometheus metrics
- a cmd/manager/main.go to run

项目将在写入项目文件后提示用户运行 "dep sure"。

**kubebuilder init \[FLAGS]**

FLAGS

1. --domain STRING # domain for groups (default "my.domain")
2. --fetch-deps # 确保下载依赖项。(default true)
3. --license STRING # license to use to boilerplate, may be one of 'apache2', 'none' (default "apache2")
4. --owner STRING # 在每个代码文件的开头添加 Cpoyright
5. --project-version STRING # project version (default "2")
6. --repo STRING # 用于 go 模块的名称(例如 github.com/user/repo)，默认为当前工作目录的 go 包名。
7. --skip-go-version-check # 如果指定，请跳过检查 Go 版本

Examples:

- 使用 "The Kubernetes authors" 作为所有者的 apache2 许可证来搭建项目
  - kubebuilder init --domain example.org --license apache2 --owner "The Kubernetes authors"
