---
title: ctr
---

# 概述

> 参考：
> - 官方文档

# ctr \[GLOBAL OPTIONS] COMMAND \[OPTIONS] \[ARGs...]

GLOBAL OPTIONS:

COMMANDS:

- plugins, plugin provides information about containerd plugins
- version print the client and server versions
- containers, c, container manage containers
- content manage content
- events, event display containerd events
- images, image, i manage images
- leases manage leases
- namespaces, namespace, ns manage namespaces
- pprof provide golang pprof outputs for containerd
- run run a container
- **snapshots, snapshot** # manage snapshots
- **tasks, t, task** # manage tasks
- **install** # install a new package
- **oci **# OCI tools
- **shim** # interact with a shim directly

## tasks # 任务管理

`create` 的命令创建了容器后，并没有处于运行状态，只是一个静态的容器。一个 container 对象只是包含了运行一个容器所需的资源及配置的数据结构，这意味着 namespaces、rootfs 和容器的配置都已经初始化成功了，只是用户进程(这里是 `nginx`)还没有启动。然而一个容器真正的运行起来是由 task 对象实现的，`task` 代表任务的意思，可以为容器设置网卡，还可以配置工具来对容器进行监控等。

所以还需要通过 task 启动容器：

    ?  → ctr task start -d nginx
    ?  → ctr task ls
    TASK     PID       STATUS
    nginx    131405    RUNNING

当然，也可以一步到位直接创建并运行容器：

    ?  → ctr run -d docker.io/library/nginx:alpine nginx

进入容器：

    # 和 docker 的操作类似，但必须要指定 --exec-id，这个 id 可以随便写，只要唯一就行
    ?  → ctr task exec --exec-id 0 -t nginx sh

暂停容器：

    # 和 docker pause 类似
    ?  → ctr task pause nginx

容器状态变成了 PAUSED：

    ?  → ctr task ls
    TASK     PID       STATUS
    nginx    149857    PAUSED

恢复容器：

    ?  → ctr task resume nginx

**ctr 没有 stop 容器的功能，只能暂停或者杀死容器。**
杀死容器：

    ?  → ctr task kill nginx

获取容器的 cgroup 信息：

    # 这个命令用来获取容器的内存、CPU 和 PID 的限额与使用量。
    ?  → ctr task metrics nginx
    ID       TIMESTAMP
    nginx    2020-12-15 09:15:13.943447167 +0000 UTC
    METRIC                   VALUE
    memory.usage_in_bytes    77131776
    memory.limit_in_bytes    9223372036854771712
    memory.stat.cache        6717440
    cpuacct.usage            194187935
    cpuacct.usage_percpu     [0 335160 0 5395642 3547200 58559242 0 0 0 0 0 0 6534104 5427871 3032481 2158941 8513633 4620692 8261063 3885961 3667830 0 4367411 356280 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1585841 0 7754942 5818102 21430929 0 0 0 0 0 0 1811840 2241260 2673960 6041161 8210604 2991221 10073713 1111020 3139751 0 640080 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
    pids.current             97
    pids.limit               0

查看容器中所有进程的 `PID`：

    ?  → ctr task ps nginx
    PID       INFO
    149857    -
    149921    -
    149922    -
    149923    -
    149924    -
    149925    -
    149926    -
    149928    -
    149929    -
    149930    -
    149932    -
    149933    -
    149934    -
    ...

注意：这里的 PID 是宿主机看到的 PID，不是容器中看到的 PID。

# 应用示例

导入名为 123.tar 的镜像到 k8s.io 这个名称空间中

- **ctr -n k8s.io image import 123.tar**

# 其他

ctr 目前很多功能做的还没有 docker 那么完善，但基本功能已经具备了。下面将围绕**镜像**和**容器**这两个方面来介绍其使用方法。

### **镜像**

**镜像下载：**

    ?  → ctr i pull docker.io/library/nginx:alpine
    docker.io/library/nginx:alpine:                                                   resolved       |++++++++++++++++++++++++++++++++++++++|
    index-sha256:efc93af57bd255ffbfb12c89ec0714dd1a55f16290eb26080e3d1e7e82b3ea66:    done           |++++++++++++++++++++++++++++++++++++++|
    manifest-sha256:6ceeeab513f7d15cea38c1f8dfe5455323b5a1bfd568516b3b0ee70406f75247: done           |++++++++++++++++++++++++++++++++++++++|
    config-sha256:0fde4fb87e476fd1655b3f04f55aa5b4b3ef7de7c701eb46573bb5a5dcf66fd2:   done           |++++++++++++++++++++++++++++++++++++++|
    layer-sha256:abaddf4965e5e9ce9953f2e136b3bf9cc15365adbcf0c68b108b1cc26c12b1be:    done           |++++++++++++++++++++++++++++++++++++++|
    layer-sha256:05e7bc50f07f000e9993ec0d264b9ffcbb9a01a4d69c68f556d25e9811a8f7f4:    done           |++++++++++++++++++++++++++++++++++++++|
    layer-sha256:c78f7f670e47cf98494e7dbe08e463d34c160bf6a5939a2155ff4438cb8b0e80:    done           |++++++++++++++++++++++++++++++++++++++|
    layer-sha256:ce77cf6a2ede66c463dcdd39f1a43cfbac3723a99e94f697bc20faee0f7cce1b:    done           |++++++++++++++++++++++++++++++++++++++|
    layer-sha256:3080fd9f46494247c9298a6a3d9694f03f6a32898a07ffbe1c17a0752bae5c4e:    done           |++++++++++++++++++++++++++++++++++++++|
    elapsed: 17.3s                                                                    total:  8.7 Mi (513.8 KiB/s)
    unpacking linux/amd64 sha256:efc93af57bd255ffbfb12c89ec0714dd1a55f16290eb26080e3d1e7e82b3ea66...
    done

**本地镜像列表查询：**

    ?  → ctr i ls
    REF                                                               TYPE                                                      DIGEST                                                                  SIZE      PLATFORMS                                                                                LABELS
    docker.io/library/nginx:alpine                                    application/vnd.docker.distribution.manifest.list.v2+json sha256:efc93af57bd255ffbfb12c89ec0714dd1a55f16290eb26080e3d1e7e82b3ea66 9.3 MiB   linux/386,linux/amd64,linux/arm/v6,linux/arm/v7,linux/arm64/v8,linux/ppc64le,linux/s390x -

这里需要注意 PLATFORMS，它是镜像的能够运行的平台标识。
**将镜像挂载到主机目录：**

    ?  → ctr i mount docker.io/library/nginx:alpine /mnt
    ?  → tree -L 1 /mnt
    /mnt
    ├── bin
    ├── dev
    ├── docker-entrypoint.d
    ├── docker-entrypoint.sh
    ├── etc
    ├── home
    ├── lib
    ├── media
    ├── mnt
    ├── opt
    ├── proc
    ├── root
    ├── run
    ├── sbin
    ├── srv
    ├── sys
    ├── tmp
    ├── usr
    └── var
    18 directories, 1 file

**将镜像从主机目录上卸载：**

    ?  → ctr i unmount /mnt

**将镜像导出为压缩包：**

    ?  → ctr i export nginx.tar.gz docker.io/library/nginx:alpine

**从压缩包导入镜像：**

    ?  → ctr i import nginx.tar.gz

其他操作可以自己查看帮助：

    ?  → ctr i --help
    NAME:
       ctr images - manage images
    USAGE:
       ctr images command [command options] [arguments...]
    COMMANDS:
       check       check that an image has all content available locally
       export      export images
       import      import images
       list, ls    list images known to containerd
       mount       mount an image to a target path
       unmount     unmount the image from the target
       pull        pull an image from a remote
       push        push an image to a remote
       remove, rm  remove one or more images by reference
       tag         tag an image
       label       set and clear labels for an image
    OPTIONS:
       --help, -h  show help

对镜像的更高级操作可以使用子命令 `content`，例如在线编辑镜像的 `blob` 并生成一个新的 `digest`：

    ?  → ctr content ls
    DIGEST         SIZE AGE  LABELS
    ...
    ...
    sha256:fdd7fff110870339d34cf071ee90fbbe12bdbf3d1d9a14156995dfbdeccd7923 740B 7 days  containerd.io/gc.ref.content.2=sha256:4e537e26e21bf61836f827e773e6e6c3006e3c01c6d59f4b058b09c2753bb929,containerd.io/gc.ref.content.1=sha256:188c0c94c7c576fff0792aca7ec73d67a2f7f4cb3a6e53a84559337260b36964,containerd.io/gc.ref.content.0=sha256:b7199797448c613354489644be1f60aa2d8e9c2278989100c72ede3001334f7b,containerd.io/distribution.source.ghcr.fuckcloudnative.io=yangchuansheng/grafana-backup-tool
    ?  → ctr content edit --editor vim sha256:fdd7fff110870339d34cf071ee90fbbe12bdbf3d1d9a14156995dfbdeccd7923

### **容器**

创建容器：

    ?  → ctr c create docker.io/library/nginx:alpine nginx
    ?  → ctr c ls
    CONTAINER    IMAGE                             RUNTIME
    nginx        docker.io/library/nginx:alpine    io.containerd.runc.v2

查看容器的详细配置：

    # 和 docker inspect 类似
    ?  → ctr c info nginx

其他操作可以自己查看帮助：

    ?  → ctr c --help
    NAME:
       ctr containers - manage containers
    USAGE:
       ctr containers command [command options] [arguments...]
    COMMANDS:
       create           create container
       delete, del, rm  delete one or more existing containers
       info             get info about a container
       list, ls         list containers
       label            set and clear labels for a container
       checkpoint       checkpoint a container
       restore          restore a container from checkpoint
    OPTIONS:
       --help, -h  show help

### **命名空间**

除了 k8s 有命名空间以外，Containerd 也支持命名空间。

    ?  → ctr ns ls
    NAME    LABELS
    default

如果不指定，`ctr` 默认是 `default` 空间。目前 Containerd 的定位还是解决运行时，所以目前他还不能完全替代 `dockerd`，例如使用 `Dockerfile` 来构建镜像。其实这不是什么大问题，我再给大家介绍一个大招：**Containerd 和 Docker 一起用！**

### **Containerd + Docker**

事实上，Docker 和 Containerd 是可以同时使用的，只不过 Docker 默认使用的 Containerd 的命名空间不是 default，而是 `moby`。下面就是见证奇迹的时刻。

首先从其他装了 Docker 的机器或者 GitHub 上下载 Docker 相关的二进制文件，然后使用下面的命令启动 Docker：

    ?  → dockerd --containerd /run/containerd/containerd.sock --cri-containerd

1
2
Plain Text

接着用 Docker 运行一个容器：

    ?  → docker run -d --name nginx nginx:alpine

1
2
Plain Text

现在再回过头来查看 Containerd 的命名空间：

    ?  → ctr ns ls
    NAME    LABELS
    default
    moby

1
2
3
4
5
Plain Text

查看该命名空间下是否有容器：

    ?  → ctr -n moby c ls
    CONTAINER                                                           IMAGE    RUNTIME
    b7093d7aaf8e1ae161c8c8ffd4499c14ba635d8e174cd03711f4f8c27818e89a    -        io.containerd.runtime.v1.linux

1
2
3
4
Plain Text

我艹，还可以酱紫？看来以后用 Containerd 不耽误我 `docker build` 了~~最后提醒一句：Kubernetes 用户不用惊慌，Kubernetes 默认使用的是 Containerd 的 `k8s.io` 命名空间，所以 `ctr -n k8s.io` 就能看到 Kubernetes 创建的所有容器啦，也不用担心 `crictl` 不支持 load 镜像了，因为 `ctr -n k8s.io` 可以 load 镜像啊，嘻嘻?
