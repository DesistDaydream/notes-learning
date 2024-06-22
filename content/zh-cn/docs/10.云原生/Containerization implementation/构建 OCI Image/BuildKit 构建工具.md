---
title: BuildKit 构建工具
---

# 概述

> 参考：
> 
> - [GitHub 项目，moby/buildkit](https://github.com/moby/buildkit)
> - [官方文档](https://docs.docker.com/develop/develop-images/build_enhancements/)
> - [知乎，官方下一代Docker镜像构建神器 -- BuildKit](https://zhuanlan.zhihu.com/p/137261919)

BuildKit 是 [Docker](/docs/10.云原生/Containerization%20implementation/Docker/Docker%20介绍/Docker%20介绍.md) 上游社区(Moby)推出的下一代镜像构建神器，可以更加快速，有效，安全地构建容器镜像。Docker v18.06 已经集成了该组件。BuildKit 可用于多种导出格式（例如 OCI 或 Docker）以及前端支持（Dockerfile），并提供高效缓存和运行并行构建操作等功能。BuildKit 仅需要容器运行时就能执行，当前受支持的运行时包括 Containerd 和 Runc。

## 构建步骤优化

Docker 提供的原始构建最令人沮丧的问题之一是 Dockerfile 指令执行构建步骤的顺序性。在引入多阶段构建之后，可以将构建步骤分组为单独的逻辑构建任务在同一个 Dockerfile 中。

有时，这些构建阶段是彼此完全独立的，这意味着它们可以并行执行-或根本不需要执行。遗憾的是，传统的 Docker 镜像构建无法满足这种灵活性。这意味着构建时间通常会比绝对必要的时间更长。

相比之下，BuildKit 会创建一个构建步骤之间的依赖关系图，并使用该图来确定可以忽略构建的哪些元素;可以并行执行的元素;需要顺序执行的元素。这可以更有效地执行构建，这对开发人员来说很有价值，因为他们可以迭代其应用程序的镜像构建。

## 高效灵活的缓存

虽然在旧版 Docker 镜像构建中缓存构建步骤非常有用，但效率却不如预期。作为对构建后端的重写，BuildKit 在此方面进行了改进，并提供了更快，更准确的缓存机制。使用为构建生成的依赖关系图，并且基于指令定义和构建步骤内容。

BuildKit 提供的另一个巨大好处是以构建缓存导入和导出的形式出现，正如 Kaniko 和 Makisu 允许将构建缓存推送到远程注册表一样，BuildKit 也是如此，但是 BuildKit 使您可以灵活地将缓存嵌入到内部注册表中。镜像（内联）并将它们放在一起（虽然不是每个注册表都支持），或者将它们分开导入。也可以将缓存导出到本地目录以供以后使用。

当从头开始建立构建环境而没有任何先前的构建历史时，导入构建缓存的能力就发挥了自己的作用：导入“预热”缓存，对于临时 CI/CD 环境特别有用。

## 工件

当使用旧版 Docker 镜像构建器构建镜像时，将生成的镜像添加到 Docker 守护进程管理的本地镜像的缓存中。需要单独的`docker push`将该镜像上载到远程容器镜像注册表。新的工件构建工具通过允许您在构建调用时指定镜像推送来增强体验，BuildKit 也不例外，它还允许以几种不同格式输出镜像；本地目录中的文件，本地 tarball，一个本地 OCI 镜像 tarball，一个 Docker 镜像 tarball，一个存储在本地缓存中的 Docker 镜像以及一个推送到注册表的 Docker 镜像，有很多格式！

## 扩展语法

对于 docker 构建体验而言，经常重复出现的众多功能请求之一就是安全处理镜像构建过程中所需的机密信息。Moby 项目抵制了这一要求很多年了，但是，借助 BuildKit 灵活的“前端”定义，为 Buildkit 提供了一个实验性前端，它扩展了 Dockerfile 语法。扩展后的语法为 RUN Dockerfile 指令提供了有用的补充，其中包括安全性功能。

```dockerfile
RUN --mount=type=secret,id=top-secret-passwd my_command
```

引用实验性前端的 Dockerfile 可以为 RUN 指令临时挂载秘钥。使用 `--secret` 标志将秘钥提供给构建，用于 docker build。使用 ssh mount 类型可以转发 SSH 代理连接以实现安全 SSH 身份验证。

## BuildKit 使用场景

BuildKit 还有许多其他功能，可以极大地改善构建容器镜像的技巧。如果它是适用于许多不同环境的通用工具，那么如何使用它呢？

根据您工作的环境，这个问题的答案是多种多样的。让我们来看看。

### Docker

尽管目前 BuildKit 不是 Docker 的默认构建工具，但是完全可以考虑将其作为 Docker（v18.09 +）的首选构建工具。当然目前在 windows 平台是不支持的。

临时方案是设置环境变量`DOCKER_BUILDKIT=1`。如果是想永久生效的话，将`"features":{"buildkit": true}` 添加到 docker 守护进程的配置文件中。在此配置中，由于 Docker 守护程序中的当前限制，Docker 并未充分展现 BuildKit 的全部功能。因此，Docker 客户端 CLI 已扩展为提供插件框架，该框架允许使用插件扩展提供了可用的 CLI 功能。一个名为`Buildx`的实验性插件会绕过守护程序中的旧版构建函数，并使用 BuildKit 后端进行所有构建，它提供所有熟悉的镜像构建命令和功能，但通过一些特定于 BuildKit 的附加功能对其进行了扩充。

BuildKit 以及 Buildx 的都支持多个构建器实例，这是一项重要功能，这实际上意味着可以共享一个构建器实例场以进行构建;也许是一个项目被分配了一组构建器实例。

```bash
$ docker buildx ls
NAME/NODE DRIVER/ENDPOINT STATUS PLATFORMS
default * docker
  default default running linux/amd64, linux/386
```

默认情况下，Buildx 插件以 docker 驱动程序为目标，该驱动程序使用 Docker 守护程序提供的 BuildKit 库具有其固有的局限性。另一个驱动程序是 docker-container，它可以透明地在容器内启动 BuildKit 以执行构建。 BuildKit 中提供的功能 CLI：这是否是理想的工作流程，完全取决于个人或公司的选择。

### Kubernetes

越来越多的组织将构建放到 Kubernetes 当中，通常将容器镜像构建作为 CI/CD 工作流的一部分出现在 pod 中。在 Kubernetes 中运行 BuildKit 实例时，有一个每种部署策略都有其优缺点，每种策略都适合不同的目的。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/buildkit/kubernetes.png)

除了使用 Docker CLI 为 BuildKit 启动面向开发人员的构建之外，构建还可以通过多种 CI/CD 工具触发。使用 BuildKit 进行的容器镜像构建可以作为 Tekton Pipeline Task 执行。

## 结论

本文主要讲了 BuildKit 诸多特性和使用场景。

目前类似的工具不少，如 Redhat 的 Buildah，Google 的 Kaniko 或 Docker 的 BuildKit。

不过 BuildKit 是官方提供，和 docker 本身结合比较好。

# 部署 Buildkit

## Docker 启用 Buildkit

从 Docker 23.0 版本开始，Buildx 是默认的构建工具（Buildx 就是使用 BuildKit 作为构建工具），若使用传统的 `docker build` 命令，将会出现弃用提示：

```bash
DEPRECATED: The legacy builder is deprecated and will be removed in a future release.
            Install the buildx component to build images with BuildKit:
            https://docs.docker.com/go/buildx/
```

从 [buildx 的 Release 页面](https://github.com/docker/buildx/releases)下载二进制文件，添加可执行权限，并放到 `$HOME/.docker/cli-plugins/` 目录中。

```bash
export BuildxVersion="0.11.2"
wget https://github.com/docker/buildx/releases/download/v${BuildxVersion}/buildx-v${BuildxVersion}.linux-amd64
mkdir -p $HOME/.docker/cli-plugins
mv buildx-v${BuildxVersion}.linux-amd64 $HOME/.docker/cli-plugins/docker-buildx
chmod 755 $HOME/.docker/cli-plugins/docker-buildx
```

若是 Docker 23.0 版本之前，则还需要通过两种方式启用 Buildkit 功能

1. 构建之前添加环境变量 `export DOCKER_BUILDKIT=1`
2. 在 /etc/docker/deamon.json 文件中，添加 `"features": { "buildkit": true }`。

## Nerdctl 使用 Buildkit

正常部署 buildkitd 即可，使用 systemd 启动，下载 buildkitd 二进制文件，创建 [unit file](https://github.com/moby/buildkit/tree/master/examples/systemd/system)
