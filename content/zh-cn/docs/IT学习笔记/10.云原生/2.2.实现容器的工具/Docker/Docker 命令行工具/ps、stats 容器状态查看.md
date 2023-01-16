---
title: ps、stats 容器状态查看
---

# docker ps

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

<!---->

    [root@k8s-monitor-agent ~]# docker ps --format "table {{.Names}}\t{{.Size}}"
    NAMES               SIZE
    pushgateway         46B (virtual 19.4MB)
    node_exporter       16B (virtual 22.9MB)

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

# docker stats \[OPTIONS] \[CONTAINER....]

显示效果如下，可以显示容器的 CPU、内存的使用率，和磁盘的 IO。并实时刷新。

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

- docker stats --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}" # 使用 go 模板输出指定内容
