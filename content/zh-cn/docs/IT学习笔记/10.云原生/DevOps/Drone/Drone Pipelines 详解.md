---
title: Drone Pipelines 详解
---

# 概述

> 参考：
>
> - 官方文档：<https://docs.drone.io/pipeline/overview/>

Drone 默认通过 yaml 格式的名为 .drone.yml 的文件来指定 Pipelines 的行为。yaml 语法易于阅读和表达，将 .drone.yml 提交到代码仓库的根目录后，如果代码仓库配置了 webhook 并指定 Drone Server，那么当代码更改并触发 webhook 后，Drone Server 首先就会 clone 该仓库，并根据其中的 .drone.yml 文件中的内容，来执行后续 pipeline 动作。

Note:Pipelines 从逻辑上看，就是一个或多个 Drone Plugins 的合集，每一个步骤都使用一个插件。具体原因详见下文 Pipeline 插件

## .drone.yml 文件基本样例

```yaml
kind: pipeline
type: docker
name: default
steps:
  - name: greeting
    image: alpine
    commands:
      - echo hello
      - echo world
```

1. kind # 指定本次任务的种类，该示例任务为 pipeline 种类。还有 secret 与 signature 种类
2. type # 指定 pipeline 的类型，不同的类型调用不同的 Runner 来执行任务。如果不填 type，则默认为 docker 类型。该示例任务为 docker 类型(docker 类型意味着后续的所有步骤都会启动一个容器，然后再容器中执行步骤的内容)。
3. name # The name attribute defines a name for your pipeline. You can define one or many pipelines for your project.
4. steps # The steps section defines an array of pipeline steps that are executed serially. If any step in the pipeline fails, the pipeline exits immediately.
5. name # The name attribute defines the name of the pipeline step.
6. image # 指定执行该步骤所使用的插件的镜像。
7. commands # 指定执行该步骤的具体内容，commands 表示在容器中执行 shell 命令。如果任何命令的退出码非 0，则 pipeline 将会失败。

### 该 Pipelines 样例的流程详解

1. 通过 drone/git 镜像启动容器，下载代码仓库中的代码
2. 通过 alpine 镜像启动容器，执行 commands 字段下给定的两条命令。

Note：

1. 这其中第二步生成的容器，其工作目录中会包含第一步从代码仓库中 clone 下来的所有文件。
2. 如果在 steps 中有新的文件生成，比如构建代码生成了二进制文件，那么这个新生成的文件连带着 clone 下来的文件，都会带到后面所有步骤启动的容器中。

## Pipeline 的类型

Drone 支持不同类型的 pipeline，上面的示例使用的是 docker 类型。每种类型的 pipeline 都对不同的环境行了优化，现阶段 Drone 包括以下几种 pipeline：

Note:Pipeline 的类型与 Runner 的类型相对应，比如 kubernetes 类型的 Pipeline 就需要一个在 kubernetes 类型的 Runner，否则当 Pipeline 触发时，其他类型的 Runner 无法连接 k8s 集群，也就无法继续处理后续任务。

1. Docker Pipelines
2. Kubernetes Pipelines
3. Exec Pipelines
4. SSH Pipelines
5. Digital Ocean Pipelines
6. MacStadium Pipelines
7. 等等等

## Pipeline 插件

官方文档：<https://docs.drone.io/plugins/overview/>

在 Drone Pipeline 中可以通过 Plugins(插件) 扩展功能，以实现更复杂的 CI/CD 需求。

Plugins 的特点：

1. 所谓的插件其实就是一个一个容器，这个容器就是 Pipeline 中每个 steps 中指定镜像启动的。
2. 每种类型的 Pipeline 都有自己可用的插件，各类型 Pipeline 下的 Plugins 的使用方法详见该类型的文档。
3. 插件可以用于执行命令、部署代码、发布版本、发送通知等等等。
4. 比如上文中 .drone.yml 样例中 steps 下的 alpine ，就可以称之为一个插件。
5. 由于插件只是一个 Docker 容器，这意味着我可以在任何运行的容器中使用任何变成语言来编写一个具有某些功能的插件。然后将这个容器提交成一个 image，在 Pipeline 中使用。

比如在 Pipeline 中有这么一个步骤：

    - name: publish
      image: plugins/docker
      settings:
        username: kevinbacon
        password: pa55word
        repo: lchdzh/drone-docker-build-test
        tags:
        - 1.0.0
        - 1.0

这个步骤使用了名为 plugins/docker 的插件，这个插件的作用就是获取其工作目录下的 Dockerfile 文件，并使用该文件构建 docker image，构建完成后，推送到指定的镜像仓库中。

通过 .drone.yml 文件中的 settings 字段可以对该插件的运行行为进行定义(作为环境变量传递到容器中)。这个例子就是将 username 与 password 作为登陆镜像仓库的用户名和密码；repo 是构建完成后，将要把镜像推送到哪个仓库；tags 是构建完成后镜像的 tag。settings 中会进行如下转换并传递到容器中：

    PLUGIN_USERNAME=kevinbacon
    PLUGIN_PASSWORD=pa55word
    PLUGIN_REPO=lchdzh/drone-docker-build-test
    PLUGIN_TAGS=1.0.0,1.0

# Docker 类型的 Pipeline

官方文档：<https://docs.drone.io/pipeline/docker/overview/>

# Kubernetes 类型的 Pipelin

官方文档：<https://docs.drone.io/pipeline/kubernetes/overview/>
