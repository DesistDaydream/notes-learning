---
title: Docker 命令行工具
linkTitle: Docker 命令行工具
date: 2023-11-03T22:25
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，参考 - CLI 参考 - docker](https://docs.docker.com/engine/reference/commandline/docker/)

# Syntax(语法)

**docker \[OPTIONS] COMMAND \[ARG...]**

## OPTIONS

- --config=~/.docker # Location of client config files # 客户端配置文件的位置
- -D, --debug=false # Enable debug mode # 启用 Debug 调试模式
- -H, --host=\[] # Daemon socket(s) to connectto # 守护进程的套接字（Socket）连接
- -l, --log-level=info # Set the logging level # 设置日志级别
- --tls=false # Use TLS; implied by--tlsverify #
- --tlscacert=~/.docker/ca.pem # Trust certs signed only by this CA # 信任证书签名 CA
- --tlscert=~/.docker/cert.pem # Path to TLS certificate file # TLS 证书文件路径
- --tlskey=~/.docker/key.pem # Path to TLS key file # TLS 密钥文件路径
- --tlsverify=false # Use TLS and verify theremote # 使用 TLS 验证远程

# Management Commands

management command 在使用的时候，当后面还需要跟其子命令的时候，是可省的。直接使用子命令就表示对其执行，但是有的管理命令不行，比如 create，对于 container 可省，对于 network 不可省

## container - 容器管理

attach cp diff export kill ls port rename rm start stop unpause wait

commit create exec inspect logs pause prune restart run stats top update

对容器的操作说明：docker 的 container 相关操作命令有一些关于如何进入容器操作的命令

其更本思想为：

- 连接标准输入，输入到 host 上的内容同样输入到 container 中
- 连接标准输出，输出到 container 中的同样输出到 host 上
- 可以分配一个终端(shell)给 container，以便操作更便捷

注意：有的 Container 在启动的时候，会自带命令去执行一些操作，该操作会自动输出一些内容，当连接到该 Container 的标准输入和输出上之后，可能没法对其进行输入操作，因为该 Container 正在其前台运行某程序(就好像平时用 linux 的 tailf 命令似的)(有的程序有输出内容，有的程序没有输出内容)，前台运行程序的时候，是没法输入的。

**docker container \[OPTIONS]**

EXAMPLE

- docker container prune -a # 清理所有已经停止的 container

## image - Docker 镜像的管理命令

详见 [image](/docs/10.云原生/Containerization%20implementation/Docker/Docker%20命令行工具/image.md)

## network - Docker 网络的管理命令

详见 [network](/docs/10.云原生/Containerization%20implementation/Docker/Docker%20命令行工具/network.md)

## plugin Manage plugins

## secret Manage Docker secrets

## service Manage services

## stack Manage Docker stacks

## system - docker 系统管理

**docker system COMMAND**

COMMAND:

- df # 显示 docker 系统的磁盘使用情况，效果如下：

```bash
~]# docker system df
TYPE                TOTAL               ACTIVE              SIZE                RECLAIMABLE
Images              11                  10                  1.043GB             96.12MB (9%)
Containers          34                  18                  4.588kB             2.294kB (50%)
Local Volumes       2                   2                   152.8MB             0B (0%)
Build Cache         0                   0                   0B                  0B
```

- events # Get real time events from the server
- info # 等同于 docker info 命令
- prune # 删除未使用的数据。删除内容如下

```bash
~]# docker system prune
WARNING! This will remove:
  - all stopped containers
  - all networks not used by at least one container
  - all dangling images
  - all dangling build cache
```

## volume - 管理 docker 的卷

详见 《[Docker 存储](/docs/10.云原生/2.2.实现容器的工具/Docker/Docker%20 存储.md 存储.md)》

# COMMANDS

## attach - 当前 shell 下 attach(连接)到指定运行中的镜像

Attach local standard input, output, and error streams to a running container # 把本地终端上的标准输入、输出和错误数据流连接到一个运行中的容器(即从一个运行中的容器剥离了其终端，再重新连接到其终端)

## build - 通过 Dockerfile 定制镜像

docker build \[OPTIONS] PATH # 使用 dockerfile 文件自动创建镜像

注：PATH 是 build context 创建环境的位置，从创建环境中搜索 Dockerfile 文件来使用。PATH 也可以是一个 URL，通过网上下载镜像。

OPTIONS：

- -t \<NAME> # 自己定义一个创建完成后的镜像名 NAME
- -f \<NAME> # 指定使用创建环境中哪个 NAME 文件作为 Dockerfile,默认使用文件名为 Dockerfile 的文件
- --no-cache # 创建镜像时不使用缓存

EXAMPLE

- docker build ./ # 从当前目录下查找名为 Dockerfile 的文件进行 Image 的创建
- docker build -t ubuntu-vi -f test /dockerfile/ # 使用 dockerfile 目录，并使用该目录中的 test 文件作为 dockerfile 文件，创建一个名为 ubuntu-vi 的 Image

## commit - 从容器的变化中创建一个新的 image。提交当前容器为新的镜像

docker commit \[OPTIONS] CONTAINER \[REPOSITORY\[:TAG]]

使用 docker commit ContainerName NewName 命令创建新镜像

(ContainerName 为正在运行的容器名 NewName 为需要创建的镜像名，自己定)

每一次 commit 就相当于把当前的可写入层变成 image 的一层

## cp - 从容器中拷贝指定文件或者目录到宿主机中

## create - 创建一个新的容器，同 run 但不启动容器

## diff - Inspect changes on a container's filesystem 查看 docker 容器变化

## events - Get real time events from the server 从 docker 服务获取容器实时事件

## exec - 在运行中的容器上执行命令

## export - 导出容器的文件系统为一个 tar 归档文件(对应 import)

## history - 展示一个镜像形成历史

docker history \[OPTIONS] IMAGE

OPTIONS

- --no-trunc # 不要截断输出 i.e.每列显示的内容都是完整内容，不会被截断

EXAMPLE

## images - 列出系统当前镜像

OPTIONS

- --all , -a # Show all images (default hides intermediate images)
- --digests # Show digests
- --filter , -f # Filter output based on conditions provided
- --format # Pretty-print images using a Go template
  - 可用的 Go 模板占位符详见：https://docs.docker.com/engine/reference/commandline/images/#format-the-output
- --no-trunc # Don’t truncate output
- --quiet , -q # Only show numeric IDs

EXAMPLE

- docker images # 查看本地 images，效果如图

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/urb4r9/1616121613241-36e0f6eb-bee7-4db6-a9ed-d02ca5cd679d.png)

## import - Import the contents from a tarball to create a filesystem image 从 tar 包中的内容创建一个新的文件系统映像(对应 export)

## info - 与 docker system info 命令效果相同

## inspect - 返回有关容器或镜像的底层信息

显示 docker 所能管理的所有 object 的详细信息，object 包括 image，container，network 等等

**docker \[OBJECT] inspect \[OPTIONS]**

OPTIONS

EXAMPLE

- 获取 snmp_exporter 容器 merged 目录的绝对路径
  - **docker inspect snmp_exporter | jq .\[0].GraphDriver.Data.MergedDir | tr -d """**
  - **docker inspect snmp_exporter --format='{{.GraphDriver.Data.MergedDir}}'**
- 获取容器的 PID
  - docker inspect mysql-pxc1-1 --format='{{.State.Pid}}'

## kill - kill 一个运行中的容器

## load - 从 tar 包 或标准输入中加载一个镜像(对应 save)

**docker load \[OPTIONS]FILE**

OPTIONS

- -i # 从 tar 存档文件读取，而不是 STDIN

EXAMPLE

- docker load -i kubernetes.tar

## login - 注册或者登录到一个 Docker Registry

## logout - 从当前 Docker Registry 登出

## logs - 获取容器得日志

**docker logs \[OPTIONS] CONTAINER**

OPTIONS

- --details # Show extra details provided to logs
- **-f, --follow**# 跟踪日志的输出
- --since string Show logs since timestamp (e.g. 2013-01-02T13:23:37Z) or relative (e.g. 42m for 42 minutes)
- **-n, --tail STRING**# 从日志末尾开始显示日志的指定行数。`默认值：all`
- **-t, --timestamps** # 在每行日志行首显示时间戳
- --until string # Show logs before a timestamp (e.g. 2013-01-02T13:23:37Z) or relative (e.g. 42m for 42 minutes)

## pause - 暂停一个 Container 中的所有进程

## port - 查看映射端口对应的容器内部源端口

docker port CONTAINER \[PRIVATE_PORT\[/PROTO]]

EXAMPLE

- docker port nginx1 # 查看名为 nginx1 这个 Container 的端口映射情况

## ps - 列出容器

详见：docker ps 命令，可以查看很多容器信息

## pull - 从 Registry 拉取指定镜像或者镜像仓库

**docker pull \[REGISTRY]\[:Port]/\[NAMESPACE/]\<NAME>:\[TAG]**

如果不加 registry，则默认从 hub.docker.com 拉取 image；如果不设置 namespace，则默认从指定的 registry 中的顶层仓库拉取镜像，如果使用了 namespace，则从该用户仓库拉取镜像；如果不指定 TAG，则默认拉取 lastest 版的 image

EXAMPLE

- docker pull quay.io/coreos/flannel

## push - 推送指定镜像或者库镜像至 docker 源服务器

## rename - 重命名容器

## restart - 重启运行的容器

## rm - 移除一个或者多个容器

EXAMPLE

- docker rm `docker ps -a | grep "Exited" | awk '{print $NF}'` # 移除所以已经停止的容器

## rmi - 移除一个或多个镜像

Notes: 无容器使用该镜像才可以删除，否则需要删除相关容器才可以继续或者-f 强制删除

## run - 创建一个新的容器并运行一个命令

详见：[docker run 运行容器](/docs/10.云原生/Containerization%20implementation/Docker/Docker%20命令行工具/run.md)

## save - 保存一个或多个镜像为一个 tar 包(对应 load)

docker save \[OPTIONS] Image1 Image2 ... ImageN /PATH/FILE

OPTIONS

- -o # save 的时候写入文件，而不是 STDOUT

EXAMPLE

- docker save k8s.gcr.io/kube-proxy:v1.12.1 -o kubernetes.tar # 保存 k8s.gcr.io/kube-proxy:v1.12.1 这个 image 到 kubernetes.tar 这个文件中
- docker save -o XXXX.tar $(docker images | awk '{print $1,$2}' OFS=":" | awk 'NR!=1{print}') # 保存全部镜像到 XXX.tar 文件中

## search - 在 dockerhub 中搜索镜像

docker search \[OPTIONS] TERM

EXAMPLE

`docker search centos` # 搜索所有 centos 的 docker 镜像

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/urb4r9/1616121613281-5ce787a3-8986-43c9-926a-680de555c36e.png)

## start - 启动容器

## stats - 显示实时的容器资源使用情况统计流

命令用法详见：容器状态查看命令

## stop - 停止容器

## tag - 在一个 repostiory 中标记一个 image

**docker tag SOURCE_IMAGE\[:TAG] TARGET_IMAGE\[:TAG]**

## top - 查看容器中运行的进程信息

EXAMPLE

`docker top prometheus` # 查看 prometheus 这个 container 运行的程序，该信息的格式为 ps 命令所输出的内容

## unpause - 取消暂停容器

## wait - 截取容器停止时的退出状态值
