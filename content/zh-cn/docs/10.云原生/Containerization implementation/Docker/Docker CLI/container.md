---
title: container
linkTitle: container
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，参考 - CLI - docker - container](https://docs.docker.com/reference/cli/docker/container)


# run

> 参考：
>
> - [官方文档，参考 - 命令行参考 - docker-Docker run 参考](https://docs.docker.com/engine/reference/run/)
> - [官方文档，参考 - 命令行参考 - docker-docker run](https://docs.docker.com/engine/reference/commandline/run/)

run 命令可以启动容器

## Syntax(语法)

**docker run \[OPTIONS] ImageName \[COMMAND] \[ARG...]**

### OPTIONS

- **-d, --detach** # 让容器运行在后台并打印出容器的 ID
- **--entrypoint**(STRING) #  覆盖容器镜像的的默认 ENTRYPOINT。即 [Dockerfile 指令](/docs/10.云原生/Containerization%20implementation/构建%20OCI%20Image/Dockerfile%20指令.md) 中的 ENTRYPOINT。
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
  - Note：使用 /HOST/PATH 与 VolumeName 的区别详见：《[Docker Storage](/docs/10.云原生/Containerization%20implementation/Docker/Docker%20Storage.md)》
- --volume-driver string Optional volume driver for the container
- **--volumes-from \<ContainerName>** # 运行的新容器从 ContainerName 这个容器复制存储卷来使用
- **-w, --workdir \<STRING>** # 指定容器内的工作目录，让指定的目录执行当前命令

#### 资源配置

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

#### 特权 与 Linux Capabilities

> 参考：
>
> - [官方文档，参考-命令行参考 - docker-Docker run 参考 - 特权 与 Linux Capabilities](https://docs.docker.com/engine/reference/run/)

默认情况下，Docker 容器是 “unprivileged”，即无特权的。

另外，在官方文档中，还可以找到所有可以通过 --cap-add 和 --cap-drop 控制的容器内的 [Linux Capabilities(能力)](/docs/1.操作系统/登录%20Linux%20与%20访问控制/Access%20Control(访问控制)/Capabilities(能力)%20管理.md)。

- **--cap-add < STRING | ALL>** # 添加 Linux 能力。可以多次指定该选项，或使用 ALL 添加所有。
- **--cap-drop < STRING | ALL>** # 禁用 Linux 能力。可以多次指定该选项，或使用 ALL 禁用所有。
- **--privileged** # 特权模式

## Example

`docker run -d -p 80:80 httpd`

- 其过程可以简单的描述为
  - 从 Docker Hub 下载 httpd 镜像（镜像中已经安装好了 Apache HTTP Server）
  - 以后台启动 httpd 容器，并将容器的 80 端口映射到 host 的 80 端口。

以后台运行镜像 nginx:latest

- docker run -p 80:80 -v /data:/data -d nginx:latest

在运行 centos 容器的时候，执行 tail 命令。该命令是为了让容器启动后不自动关闭

- docker run -d centos tail -f /dev/null

# docker ps

https://docs.docker.com/reference/cli/docker/container/ls/

## Syntax(语法)

**docker ps \[OPTIONS]**

以列表的形式显示容器，包括以下几个字段 CONTAINER ID(容器 ID 号)、IMAGE(启动该容器所用的 image)、COMMAND(该容器运行的命令)、CREATED(该容器被创建了多久)、STATUS(容器当前状态)、PORTS(容器所用端口)、NAMES(容器名，随机生成)，效果如图所示：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tuw6e1/1616121626421-52b41c7f-068e-4c87-9a54-8aab4d638bb0.png)

还可以通过 -s 选项，来输出容器占用的磁盘空间大小。

**OPTIONS**

- **-a, --all** # 显示所有容器(默认只显示正在 running 状态的)
- **-f, --filter FILTER** # 根据提供的条件过滤输出内容。
  - 可用的过滤条件详见：<https://docs.docker.com/engine/reference/commandline/ps/#filtering>
  - 比较常见的是根据 volume 进行过滤，可以根据指定的 volume 来过滤，从而发现哪些容器正在使用哪些 volume。
- **--format STRING** # 使用 Go 模板漂亮得输出容器信息。
  - 可用的 Go 模板占位符详见：<https://docs.docker.com/engine/reference/commandline/ps/#formatting>
  - 可以使用 table 指令，让输出内容以表格的方式呈现，效果如下(如果没有 table 指令，那么输出内容将会扎堆)。

```bash
~]# docker ps --format "table {{.Names}}\t{{.Size}}"
NAMES               SIZE
pushgateway         46B (virtual 19.4MB)
node_exporter       16B (virtual 22.9MB)
```

- **-n, --last INT** # Show n last created containers (includes all states) (default -1)
- **-l, --latest** # 显示最后创建的容器(所有状态)
- **--no-trunc** # 不要截断输出 i.e.每列显示的内容都是完整内容，不会被截断
- **-q, --quiet** # 仅输出 CONTAINER ID
- **-s, --sizes** # 显示容器所用磁盘容量。一个是可写层的数据量，还有一个是只读镜像数据的磁盘空间总量。

## EXAMPLE

- 显示所有容器的 CONTAINER ID 与 COMMAND 字段，且不截断输出
  - docker ps --format "table {{.ID}}\t{{.Command}}" -a --no-trunc
- 查看容器所占磁盘空间大小，并按照所占空间大小排序
  - docker ps --format "{{.ID}}\t{{.Size}}" | sort -k 4 -h

### 过滤器示例

只显示状态为 restarting 的容器

- docker ps -a --filter status=restarting

只显示状态为 exited 的容器

- docker ps -a --filter status=exited

查看那个容器使用了指定的 Volume

```bash
~]# docker volume ls
DRIVER    VOLUME NAME
local     87e775bf78c42bc70b63f49f5495081d835d4571a922b2c5400371456fb9fbd1
~]# docker ps -a -f volume=87e775bf78c42bc70b63f49f5495081d835d4571a922b2c5400371456fb9fbd1
CONTAINER ID   IMAGE     COMMAND                  CREATED       STATUS                    PORTS     NAMES
1b857d27d391   mysql:8   "docker-entrypoint.s…"   7 weeks ago   Exited (0) 14 hours ago             mysql

```

# docker stats

https://docs.docker.com/reference/cli/docker/container/stats/

显示效果如下，可以显示容器的 CPU、内存的使用率，和磁盘的 I/O。并实时刷新。

```bash
CONTAINER ID        NAME                CPU %               MEM USAGE / LIMIT     MEM %               NET I/O             BLOCK I/O           PIDS
4a12a78282a5        pushgateway         0.00%               8.383MiB / 7.638GiB   0.11%               656B / 0B           0B / 0B             9
0a5fde8051fd        node_exporter       0.00%               4.312MiB / 7.638GiB   0.06%               0B / 0B
```

**docker stats \[OPTIONS] \[CONTAINER...]**

OPTIONS

- **-a, --all** # Show all containers (default shows just running)
- **--format string** # 使用 Go 模板漂亮得输出容器信息。
  - 可用的 Go 模板占位符详见：<https://docs.docker.com/engine/reference/commandline/stats/#formatting>
- **--no-stream** # 禁用流信息，仅显示第一次请求的结果。i.e.不实时刷新
- **--no-trunc** # Do not truncate output

EXAMPLE

- `docker stats --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"` # 使用 go 模板输出指定内容
