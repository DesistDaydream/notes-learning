---
title: Docker
linkTitle: Docker
date: 2024-07-24T16:43
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，指南](https://docs.docker.com/guides/)
> - [官方文档，手册](https://docs.docker.com/manuals/)
> - [官方文档，参考](https://docs.docker.com/reference/)

Docker 是一个基于 [Containerization(容器化)](/docs/10.云原生/Containerization/Containerization.md) 的开放式平台，可以 开发、分享、运行应用程序。Docker 分为两个版本

- Docker-CE # 社区版
- Docker-EE # 商业版

Docker 为了解决 LXC 无法批量管理、复刻容器等问题应运而生，简化用户对容器的应用。Docker 是 Docker.inc 公司开源的一个基于 LXC 技术之上构建的 Container 引擎，不再使用模板技术，而是使用了 [Docker Image](/docs/10.云原生/Containerization%20implementation/Docker/Docker%20Image.md) 文件的方式来创建。Image 是放在统一的互联网仓库中，当需要使用 Container 的时候，直接 run 或者 creat 等即可从仓库中下载到该 Image，然后基于该 Image 再运行 Container。

Note：一开始，Docker 在 linux 上实现容器技术的后端使用的是 lxc，后来使用 runc 来代替。

# Docker 架构

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/docker/202407241618409.png)

Docker 对使用者来讲是一个 [C/S](/docs/Standard/B_S%20和%20C_S%20架构.md) 模式的架构，Client 和 Server 使用 REST API 通过 UNIX [Socket](/docs/1.操作系统/Kernel/Process/Inter%20Process%20Communication(进程间通信)/Socket(套接字)/Socket(套接字).md) 或者网络进行通信。[Compose](/docs/10.云原生/Containerization%20implementation/Docker/Compose/Compose.md) 同样也可以作为客户端。

官方将这种架构称为 [**Docker Engine(引擎)**](https://docs.docker.com/engine/)，通常这个引擎具有：

- 一个 Server 进程 dockerd，长时间以 Daemon 形式运行
- 与 dockerd 通信的 API
- 一个 CLI 程序 docker

dockerd 是实现容器能力的核心，用来管理 **Docker Objects(Docker 对象)**，e.g. [Docker Image](/docs/10.云原生/Containerization%20implementation/Docker/Docker%20Image.md)、[Docker Runtime](/docs/10.云原生/Containerization%20implementation/Docker/Docker%20Runtime.md)、[Docker Network](/docs/10.云原生/Containerization%20implementation/Docker/Docker%20Network.md)、[Docker Storage](/docs/10.云原生/Containerization%20implementation/Docker/Docker%20Storage.md)

![Docker Architecture](https://notes-learning.oss-cn-beijing.aliyuncs.com/docker/202407301251102.png "https://newsletter.iximiuz.com/posts/ivan-on-the-server-side-1")

## 运行逻辑概述

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qqh0gm/1616122015445-eda7a719-b2a0-4fd6-8c61-b8d450d2dc3d.png)

当利用 docker run 来创建容器时，Docker 在后台运行的标准操作包括：

1. 检查本地是否存在指定的镜像，不存在就从公有仓库下载
2. 利用镜像创建并启动一个容器
3. 分配一个文件系统，并在只读的镜像层外面挂在一层可读写层
4. 从宿主主机配置的网桥接口中桥接一个虚拟接口到容器中去
5. 从地址池配置一个 ip 地址给容器
6. 执行用户指定的应用程序
7. 执行完毕后 Container 被终止
8. docker 容器默认会把容器内部第一个进程，也就是 pid=1 的程序作为 docker 容器是否正在运行的依据，如果 docker 容器 pid 挂了，那么 docker 容器便会直接退出。
9. 如果不想让 Container 运行完程序就终止，那么需要让 PID 为 1 的程序始终运行，比如 nginx 使用 daemon off 选项，或者其余任何可以让程序运行在前台的方法

# Docker 关联文件

Note：目录名中的 overlay2 指的是 docker 当前 Storage Driver 类型，使用不同的存储驱动，则目录名字也不同，如果我使用 aufs 驱动，那么目录名就会变为 aufs
由于 Docker 为 C/S 架构，所以 客户端与服务端分别有各自的配置

- 客户端就是 docker 程序，也可以称为 docker-cli、docker 命令行工具
- 服务端就是 dockerd 程序

## dockerd 程序关联文件

**/etc/docker/daemon.json** # dockerd 服务运行时配置文件。该目录与文件需要自行创建，默认不存在，以 JSON 格式为守护程序设置任何配置选项。

**/run/docker/** # container 的状态文件(state.json)、IO 文件 、netns 文件保存路径。

- **./containerd/** # container 的 IO 文件(init-stdin、init-stdout)保存路径。其内目录名为 **ContainerID**。
- **./netns/** # 网络名称空间保存路径。
- **./runtime-runc/moby/** # container 的运行时状态文件保存路径。其内目录名为 **ContainerID**。

**/run/containerd/** # container 的 bundle 以及 containerd.sock 文件保存路径

- **./io.containerd.runtime.v1.linux/moby/** # 容器启动后生成的 bundle 文件保存路径，其内目录名为 **ContainerID**。

**/var/lib/docker/** # Docker 管理的 网络、镜像、容器 等信息的保存路径。该路径为默认路径，可以通过配置修改。

- **./containers/** # 所有 container 的元数据保存路径(其中包括容器日志文件、容器运行配置等)。其内目录名为 **ContainerID**
- **./image/overlay2/** # docker images 以及 所有 layers 的元数据保存路径。
  - **./imagedb/** # images 的元数据保存路径
    - **./content/sha256/** # 所有 images 的 Image Configuration 文件保存路径。其内文件名为 **ImageID**。
    - **./metadata/sha256/** # 所有 images 的 创建时间、更新时间、父镜像的 Image Configuration 文件名 等信息保存路径，其内目录名为 **ImageID**。
      - 注意：好像只有自己在本地构建的镜像才会在该目录中记录。
  - **./layerdb** # 所有 layers 的元数据保存路径。
    - **./mounts/** # container layers 元数据保存路径，其内目录名为 ContainerID。容器创建完后，该容器的可读写层的元数据保存在此。包括可读写层父层的 chainID、可读写层的 cacheID(目录内的 mount-id 文件内容就是 **cacheID**)。
    - **./sha256/** # images layers 元数据保存路径，其内目录名为 chainID。包括 layer 的 **cacheID**
- **./overlay2/** # 所有 layers 的数据保存路径，其内目录名为 cacheID。docker run 的时候，是通过该目录中镜像层来启动的。创建容器后生成的可写层，也会保存在该目录，直到容器被删除。
- **./volumes/** # docker 创建的 volume 信息保存在该目录，如果是自动自动创建的 volume 则名为一串随机数

## docker 程序关联文件

**~/.docker/** # docker 运行时数据文件保存路径

- **./config.json** # docker login 后的信息都保存在此处，用户名和密码通过 base64 格式保存在其中。
- **./cli-plugins/** # docker 命令行工具插件的保存路径。

# Docker 部署

[Docker 部署](/docs/10.云原生/Containerization%20implementation/Docker/Docker%20部署.md)

# Docker 日志介绍

容器日志指的是每个容器打到 stdout 和 stderr 的日志，而不是容器内部的日志文件。docker 管理所有容器打到 stdout 和 stderr 的日志，其他来源的日志不归 docker 管理。我们通过 docker logs 命令查看容器日志都是读取容器打到 stdout 和 stderr 的日志。

基于日志驱动（loging driver）的日志管理机制

Docker 提供了一套通用、灵活的日志管理机制，Docker 将所有容器打到 stdout 和 stderr 的日志都统一通过日志驱动重定向到某个地方。

Docker 支持的日志驱动有很多，比如 local、json-file、syslog、journald 等等，类似插件一样，不同的日志驱动可以将日志重定向到不同的地方，这体现了 Docker 日志管理的灵活性，以热插拔的方式实现日志不同目的地的输出。

Dokcer 默认的日志日志驱动是 json-file，该驱动将将来自容器的 stdout 和 stderr 日志都统一以 json 的形式存储到本地磁盘。日志存储路径格式为：/var/lib/docker/containers/<容器 id>/<容器 id>-json.log。所以可以看出在 json-file 日志驱动下，Docker 将所有容器日志都统一重定向到了 /var/lib/docker/containers/ 目录下，这为日志收集提供了很大的便利。

注意：只有日志驱动为：local、json-file 或者 journald 时，docker logs 命令才能查看到容器打到 stdout/stderr 的日志。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qqh0gm/1616122015494-8bc7a655-2804-40b9-b3d4-0e541a93359b.png)

下面为官方支持的日志驱动列表：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qqh0gm/1616122015440-06e63de1-bbb8-4d6c-9271-b37947d483ae.png)

## Docker 日志驱动（loging driver）配置

上面我们已经知道 Docker 支持多种日志驱动类型，我们可以修改默认的日志驱动配置。日志驱动可以全局配置，也可以给特定容器配置。

- 查看 Docker 当前的日志驱动配置

```bash
docker info |grep "Logging Driver"
```

- 查看单个容器的设置的日志驱动

```bash
docker inspect  -f '{{.HostConfig.LogConfig.Type}}' 容器id
```

- Docker 日志驱动全局配置，全局配置意味所有容器都生效，编辑 /etc/docker/daemon.json 文件（如果文件不存在新建一个），添加日志驱动配置。示例：配置 Docker 引擎日志驱动为 syslog

```json
{
  "log-driver": "syslog"
}
```

- 给特定容器配置日志驱动，在启动容器时指定日志驱动 --log-driver 参数。示例：启动 nginx 容器，日志驱动指定为 journald

```bash
docker  run --name nginx -d --log-driver journald nginx
```

## Docker 默认的日志驱动 json-file

json-file 日志驱动记录所有容器的 STOUT/STDERR 的输出 ，用 JSON 的格式写到文件中，每一条 json 日志中默认包含 log, stream, time 三个字段，示例日志如下：文件路径为： /var/lib/docker/containers/40f1851f5eb9e684f0b0db216ea19542529e0a2a2e7d4d8e1d69f3591a573c39/40f1851f5eb9e684f0b0db216ea19542529e0a2a2e7d4d8e1d69f3591a573c39-json.log

```json
 {"log":"14:C 25 Jul 2019 12:27:04.072 * DB saved on disk\n","stream":"stdout","time":"2019-07-25T12:27:04.072712524Z"}
```

那么打到磁盘的 json 文件该如何配置轮替，防止撑满磁盘呢？每种 Docker 日志驱动都有相应的配置项日志轮转，比如根据单个文件大小和日志文件数量配置轮转。json-file 日志驱动支持的配置选项如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qqh0gm/1616122015431-a74f2c52-7a4b-443a-b6ec-d7031a967089.png)
