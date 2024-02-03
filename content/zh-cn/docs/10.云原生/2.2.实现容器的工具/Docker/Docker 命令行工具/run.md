---
title: run
linkTitle: run
date: 2023-11-03T22:39
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，参考-命令行参考-docker-Docker run 参考](https://docs.docker.com/engine/reference/run/)
> - [官方文档，参考-命令行参考-docker-docker run](https://docs.docker.com/engine/reference/commandline/run/)

run 命令可以启动容器

# Syntax(语法)

**docker run \[OPTIONS] ImageName \[COMMAND] \[ARG...]**

## OPTIONS

- **-d, --detach** # 让容器运行在后台并打印出容器的 ID
- **--entrypoint**(STRING) #  覆盖容器镜像的的默认 ENTRYPOINT。即 [Dockerfile 指令](docs/10.云原生/2.2.实现容器的工具/构建%20OCI%20Image/Dockerfile%20指令.md) 中的 ENTRYPOINT。
- **-e, --env**(LIST) # 设定容器内的环境变量。LIST 格式为 `VAR=VALUE`，若要指定多个变量，则使用多次 --env 选项。
- --env-file list Read in a file of environment variables
- **--expose**(LIST) # 等效于 Dockerfile 中的 EXPOSE 指令，仅暴露容器端口，不在宿主机暴露。
- **-h, --hostname**(STRING) # 指定容器内的 hostname
- --init Run an init inside the container that forwards signals and reaps processes
- **-i, --interactive** # 即使没有 attach 到容器，也保持 STDIN(标准输入)开启。通常与 -t 一起使用
- **--name**(STRING) # 为容器分配一个名称。默认为随机字符串
- **--network**(STRING) # 连接一个容器到一个容器网络(default "default")，可以是 docker network ls 列出的网络，也可以是其余 container 的网络。STRING 包括下面几种
  - none # 容器使用自己的网络（类似--net=bridge），但是不进行配置
  - bridge # 通过 veth 接口将容器连接到默认的 Docker 桥(默认为 docker0 的网桥).
  - host # 直接使用宿主机的网络而不是独立的 network namespace
  - ContainerName # 连接到指定 container 的网络中
  - NetworkName # 连接到 docker network ls 所列出的其中一个 docker 网络上
- **-p, --publish \[HostIP:]\[HostPort:]\<ContainerPort>**# 指明 Container 要映射到 Host 上的 IP 和端口。若只指明 HostIP 和 ContainerPort 则中间俩个冒号不可省。若不指定 HostIP，则第一个冒号可不写。要暴露多个端口则多次使用 -p 即可。
- **-P, --publish-all** # 将 Image 定义的 EXPOSE 要暴露的端口暴露给 host，随机分配 host 上的端口与之建立映射关系。一般从 10000 端口开始
- **--read-only** # 将容器的根文件系统挂载为只读模式
- **--rm** # 当容器退出时，删除它。包括创建的 volume 等一并删除
- **-t, --tty** # 为此命令分配一个 pseudo-TTY(伪终端)，可以支持终端登录，通常与-i 一起使用。
- **-u, --user**(STRING) # 为容器进程指定运行的用户名/UID
  - STRING 格式：`<NAME|UID>[:<GROUP|GID>])`
- **-v, --volume \[SRC:]DST** # 为容器创建一个 Volume 并挂载到其中的目录上。若指定的 host 上的路径不存在，则自动创建这个目录；若不指定 SRC 则 docker 会自动创建一个。默认在 /var/lib/docker/volumes/ 目录下创建 volume 所用的目录
  - Note：使用 /HOST/PATH 与 VolumeName 的区别详见：《[Docker 存储](/docs/10.云原生/2.2.实现容器的工具/Docker/Docker%20存储.md)》
- --volume-driver string Optional volume driver for the container
- **--volumes-from \<ContainerName>** # 运行的新容器从 ContainerName 这个容器复制存储卷来使用
- **-w, --workdir \<STRING>** # 指定容器内的工作目录，让指定的目录执行当前命令

### 资源配置

- --cpu-period int Limit CPU CFS (Completely Fair Scheduler) period
- --cpu-quota int Limit CPU CFS (Completely Fair Scheduler) quota
- --cpu-rt-period int Limit CPU real-time period in microseconds
- --cpu-rt-runtime int Limit CPU real-time runtime in microseconds
- -c, --cpu-shares int CPU shares (relative weight)
- **--cpus \<INT>** # 容器可使用的最大 CPU 资源
- --cpuset-cpus string CPUs in which to allow execution (0-3, 0,1)
- --cpuset-mems string MEMs in which to allow execution (0-3, 0,1)
- **-m, --memory \<BYTES>** # 内存限制。容器能使用的最大内存
- --mem ory-reservation bytes Memory soft limit
- --memory-swap bytes Swap limit equal to memory plus swap: '-1' to enable unlimited swap
- --memory-swappiness int Tune container memory swappiness (0 to 100) (default -1)
- **--restart \<string>** # 容器的重启策略。`默认值：0`
- **--ulimit \<UlimitDesc>** # 为容器配置 Ulimit。`默认值：[]`
  - 比如：
    - --ulimit nofile=1000 # 限制容器最多能打开 1 万 个文件描述符
    - --ulimit nproc=10 # 限制容器最多能打开 10 个进程

### 特权 与 Linux Capabilities

> 参考：
>
> - [官方文档，参考-命令行参考 - docker-Docker run 参考 - 特权 与 Linux Capabilities](https://docs.docker.com/engine/reference/run/)

默认情况下，Docker 容器是 “unprivileged”，即无特权的。

另外，在官方文档中，还可以找到所有可以通过 --cap-add 和 --cap-drop 控制的容器内的 [Linux Capabilities(能力)](/docs/1.操作系统/5.登录%20Linux%20与%20访问控制/Access%20Control(访问控制)/Capabilities(能力)%20管理.md)。

- **--cap-add < STRING | ALL>** # 添加 Linux 能力。可以多次指定该选项，或使用 ALL 添加所有。
- **--cap-drop < STRING | ALL>** # 禁用 Linux 能力。可以多次指定该选项，或使用 ALL 禁用所有。
- **--privileged** # 特权模式

# 最佳实践

`docker run -d -p 80:80 httpd`

- 其过程可以简单的描述为
  - 从 Docker Hub 下载 httpd 镜像。镜像中已经安装好了 Apache HTTP Server。
  - 以后台启动 httpd 容器，并将容器的 80 端口映射到 host 的 80 端口。

以后台运行镜像 nginx:latest

- docker run -p 80:80 -v /data:/data -d nginx:latest

在运行 centos 容器的时候，执行 tail 命令。该命令是为了让容器启动后不自动关闭

- docker run -d centos tail -f /dev/null
