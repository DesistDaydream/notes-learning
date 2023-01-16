---
title: Docker Runtime
---

# 概述

Runtime 和 Image 一样，也有标准，也由 OCI 维护，官方详解地址为：runtime-spec。现阶段 Docker 的 1.19 版本使用 runc 作为 Runtime

OCI Runtime 旨在指定 Container 的配置、执行环境和生命周期。

容器的配置被指定为 config.json ，并详细说明了可以创建容器的字段。指定执行环境是为了确保容器内运行的应用程序在运行时之间具有一致的环境，以及为容器的生命周期定义的常见操作。

runtime 规范有如下几个，所有人必须遵守该规范来使用 runtime 。

1. Filesystem Bundle # 文件系统捆绑。bundle 是以某种方式组织的一组文件，包含了容器所需要的所有信息，有了这个 bundle 后，符合 runtime 标准的程序(e.g.runc)就可以根据 bundle 启动容器了(哪怕没有 docker，也可以启动一个容器)。
2. Runtime and Lifecycle #
3. Linux-specific Runtime and Lifecycle # 这是关于 linux 平台的 Runtime 与 Lifecycle
4. Configuration # Configuration 包含对容器执行标准操作(比如 create、start、stop 等)所必须的元数据。这包括要运行的过程、要注入的环境变量、要使用的沙盒功能等等。不同平台(linux、window 等)，有不同的规范。
5. Linux-specific configuration #这是关于 linux 平台的 Configuration

# Docker create

有了 image 之后，就可以使用 image 来创建并启动 container 了。

docker run 命令直接创建并运行一个容器，它的背后其实包含独立的两步，一步是 docker create 创建容器，另一步是 docker start 启动容器，先介绍在 docker create 这一步中，docker 做了哪些事情。

简单点来说，dockerd 在收到客户端的创建容器请求后，做了两件事情

1. 准备容器所需的 layer
2. 检查客户端传过来的参数，并和 image 配置文件中的参数进行合并，然后存储成容器的配置文件。

<!---->

    # 创建容器前的 layers
    [root@lichenhao overlay2]# ls
    113a9d8407c2db3892944c17beba7a635ea39aa5108c7f716088466ea302a7e3  7704e53a9392b092479707d38b2b183b17bbe2cc220e2283cead9493e19aa651  l
    5de7ac8af2fb0a5fb0be4244aa07685bfcfcfc4c4b1c149bc753eb044d7f4a12  8f377ae99a442b37f5a831724951ce1cf8bfc7b874843c97d09e8027c3dd19e6
    # 创建容器后的 layers，多了两个
    [root@lichenhao overlay2]# docker create -it --name docker_runtime_test ubuntu:latest
    28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10
    [root@lichenhao overlay2]# ls
    113a9d8407c2db3892944c17beba7a635ea39aa5108c7f716088466ea302a7e3  8f377ae99a442b37f5a831724951ce1cf8bfc7b874843c97d09e8027c3dd19e6       l
    5de7ac8af2fb0a5fb0be4244aa07685bfcfcfc4c4b1c149bc753eb044d7f4a12  d976eddf7575a3464486d92539229146f3df66080a3265195791ebb0d24b24dd
    7704e53a9392b092479707d38b2b183b17bbe2cc220e2283cead9493e19aa651  d976eddf7575a3464486d92539229146f3df66080a3265195791ebb0d24b24dd-init
    # 这俩 layers 的元数据在 ${DockerRootDir}/image/${StorageDriver}/layerdb/mounts目录中
    [root@lichenhao 28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10]# pwd
    /var/lib/docker/image/overlay2/layerdb/mounts/28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10
    [root@lichenhao 28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10]# ls
    init-id  mount-id  parent
    # init-id 文件包含了 init layer 的 cacheID
    # init layer 的 cacheid 就是在 mount layer 的 cacheID 后面加上了一个“-init”
    [root@lichenhao 28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10]# cat init-id
    d976eddf7575a3464486d92539229146f3df66080a3265195791ebb0d24b24dd-init
    # mount-id 文件包含了 mount layer 的 cacheID
    [root@lichenhao 28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10]# cat mount-id
    d976eddf7575a3464486d92539229146f3df66080a3265195791ebb0d24b24dd
    # parent 里面包含的是 image 的最上 layer 的 chainID
    # 表示这个容器的 init layer 的父 layer 是 image 的最顶层 layer
    [root@lichenhao 28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10]# cat parent
    sha256:8a8d1f0b34041a66f09e49bdc03e75c2190f606b0db7e08b75eb6747f7b49e11

Note:

1. 新加的这两层 layer 比较特殊，元数据只保存在 layerdb/mounts 下面，在 layerdb/sha256 目录下没有相关信息，说明 docker 将 container 的 layer 和 image 的 layer 的元数据放在了不同的两个目录中。
2. 根据元数据中的信息，就可以通过 mount 信息中的 cacheID 来查找该 mount 信息来源于哪个容器，从而定位问题，参考：最后 docker 使用技巧 中 mount 信息章节

从上面的文章可以看到，每个创建完的容器都会新增两个层

1. mount layer #供容器写数据的层，如果容器仅创建而没运行的话，那么该层的目录中，没有 merged 目录，并且其余目录也是空的

<!---->

    [root@lichenhao d976eddf7575a3464486d92539229146f3df66080a3265195791ebb0d24b24dd]# tree
    .
    ├── diff
    ├── link
    ├── lower
    └── work
        └── work

1. init layer #包含了 docker 为容器所预先准备的文件

<!---->

    [root@lichenhao d976eddf7575a3464486d92539229146f3df66080a3265195791ebb0d24b24dd-init]# tree
    .
    ├── committed
    ├── diff
    │   ├── dev
    │   │   ├── console
    │   │   ├── pts
    │   │   └── shm
    │   └── etc
    │       ├── hostname
    │       ├── hosts
    │       ├── mtab -> /proc/mounts
    │       └── resolv.conf
    ├── link
    ├── lower
    └── work
        └── work

init layer 里面的文件有什么作用呢？从下面的结果可以看出，除了 mtab 文件是指向/proc/mounts 的软连接之外，其他的都是空的普通文件。

这几个文件都是 Linux 运行时必须的文件，如果缺少的话会导致某些程序或者库出现异常，所以 docker 需要为容器准备好这些文件：

- /dev/console: 在 Linux 主机上，该文件一般指向主机的当前控制台，有些程序会依赖该文件。在容器启动的时候，docker 会为容器创建一个 pts，然后通过 bind mount 的方式将 pts 绑定到容器里面的/dev/console 上，这样在容器里面往这个文件里面写东西就相当于往容器的控制台上打印数据。这里创建一个空文件相当于占个坑，作为后续 bind mount 的目的路径。
- hostname，hosts，resolv.conf：对于每个容器来说，容器内的这几个文件内容都有可能不一样，这里也只是占个坑，等着 docker 在外面生成这几个文件，然后通过 bind mount 的方式将这些文件绑定到容器中的这些位置，即这些文件都会被宿主机中的文件覆盖掉。
- /etc/mtab：这个文件在新的 Linux 发行版中都指向/proc/mounts，里面包含了当前 mount namespace 中的所有挂载信息，很多程序和库会依赖这个文件。

注意： 这里 mtab 指向的路径是固定的，但内容是变化的，取决于你从哪里打开这个文件，当在宿主机上打开时，是宿主机上/proc/mounts 的内容，当启动并进入容器后，在容器中打开看到的就是容器中/proc/mounts 的内容。

## 容器的元数据

容器创建完成后，就会生成容器的元数据信息，包括默认配置、运行时配置等等，文件在 ${DockerRootDir}/containers/ContainerID/\* 目录下

Note：容器启动后，该目录还会有新的文件产生。

    [root@lichenhao containers]# pwd
    /var/lib/docker/containers
    [root@lichenhao containers]# ls
    28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10
    [root@lichenhao containers]# tree
    .
    ├── checkpoints
    ├── config.v2.json # 通用的配置，如容器名称、启动后要执行的命令等等
    └── hostconfig.json # 该容器关于docker 宿主机的配置，日志驱动、是否自动删除、cgroup的配置等等

# Docker start

容器启动简单流程

- docker（client）发送启动容器命令给 dockerd
- dockerd 收到请求后，准备好 rootfs，以及一些其它的配置文件，然后通过 grpc 的方式通知 containerd 启动容器
- containerd 根据收到的请求以及配置文件位置，创建容器运行时需要的 bundle，然后启动 shim 进程，让它来启动容器
- shim 进程启动后，做一些准备工作，然后调用 runc 启动容器

容器启动后，会在下面几个目录中生成容器运行所需的内容：

- /run/docker/runtime-runc/ # 容器状态 json 文件
- /run/docker/containerd/ # 容器的 IO 文件
- /run/containerd/io.containerd.runtime.v1.linux/moby/ # 容器的 bundle 文件、pid 号
- /var/lib/containerd/io.containerd.runtime.v1.linux/moby/ # shim.stderr.log 与 shim.stdout.log

Note：这些目录在容器停止后，会自动删除

    [root@lichenhao containerd]# find / -name "28f5bed704dc*"
    /run/docker/runtime-runc/moby/28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10
    /run/docker/containerd/28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10
    /run/containerd/io.containerd.runtime.v1.linux/moby/28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10
    /var/lib/containerd/io.containerd.runtime.v1.linux/moby/28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10
    /var/lib/docker/containers/28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10
    /var/lib/docker/containers/28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10/28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10-json.log
    /var/lib/docker/image/overlay2/layerdb/mounts/28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10
    [root@lichenhao docker]# docker stop docker_runtime_test
    docker_runtime_test
    [root@lichenhao docker]# find / -name "28f5bed704dc*"
    /var/lib/docker/containers/28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10
    /var/lib/docker/containers/28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10/28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10-json.log
    /var/lib/docker/image/overlay2/layerdb/mounts/28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10

## 准备 rootfs

    # 容器运行后，在没有 mount layer 中，会多出来一个 merged 的目录，这就是当前已经启动容器的可读写层，所有变化都会在这里。
    # 并且当容器停止后，merged 目录也会随之消失
    [root@lichenhao overlay2]# ls d976eddf7575a3464486d92539229146f3df66080a3265195791ebb0d24b24dd
    diff  link  lower  merged  work
    [root@lichenhao overlay2]# ls d976eddf7575a3464486d92539229146f3df66080a3265195791ebb0d24b24dd-init/
    committed  diff  link  lower  work

    # 当在容器中创建一个文件时，该变化会同时应用到 mount layer 的 diff 和 merged 目录
    [root@lichenhao overlay2]# docker start docker_runtime_test
    docker_runtime_test
    [root@lichenhao overlay2]# docker exec -it docker_runtime_test /bin/bash
    root@28f5bed704dc:/# ls
    bin  boot  dev  etc  home  lib  lib32  lib64  libx32  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
    root@28f5bed704dc:/# touch home/lichenhao

    [root@lichenhao overlay2]# tree d976eddf7575a3464486d92539229146f3df66080a3265195791ebb0d24b24dd/diff/
    d976eddf7575a3464486d92539229146f3df66080a3265195791ebb0d24b24dd/diff/
    ├── home
    │   └── lichenhao
    └── root
    [root@lichenhao overlay2]# tree d976eddf7575a3464486d92539229146f3df66080a3265195791ebb0d24b24dd/merged/home/
    d976eddf7575a3464486d92539229146f3df66080a3265195791ebb0d24b24dd/merged/home/
    └── lichenhao
    # 容器停止后，merged 的目录消失，但是 diff 目录还在，所有对容器的操作产生的变化，都会在diff目录中永久保存，直到该容器被删除。

可以通过 mount 命令看到如下信息

    overlay on /var/lib/docker/overlay2/d976eddf7575a3464486d92539229146f3df66080a3265195791ebb0d24b24dd/merged type overlay
    (rw,relatime,lowerdir=
    /var/lib/docker/overlay2/l/QNYNLXQAPEKTOMU3TO27ITE3YO:
    /var/lib/docker/overlay2/l/2EP6BMP6AI5RGGBLLTGZURP72X:
    /var/lib/docker/overlay2/l/OLJPRTJOMYVHG3OZOMYZMBAEMQ:
    /var/lib/docker/overlay2/l/FWHKA7CXM7LSCGUQLDTAJSZFPE:
    /var/lib/docker/overlay2/l/KVQ7AO63OIRUAHZTZAS474Y3VT,
    upperdir=/var/lib/docker/overlay2/d976eddf7575a3464486d92539229146f3df66080a3265195791ebb0d24b24dd/diff,
    workdir=/var/lib/docker/overlay2/d976eddf7575a3464486d92539229146f3df66080a3265195791ebb0d24b24dd/work)

通过 overlay 联合挂载技术，将多个 layer 挂载到该容器的 mount layer 的 merged 目录中

1. lower(下层挂载) #用启动该容器的 image 的所有 layers 作为 lowerdir
2. upperdir(上层挂在) #用容器 mount layer 的 diff 目录作为 upperdir

所有在 merged 目录的变化，会同步到 diff 目录中，这样在容器停止，merged 目录消失后，所有变化依然得以保留在 diff 目录中，这样后续再启动容器的时候，上次的操作还能看到。

## 容器元数据目录的变化

rootfs 准备好之后，dockerd 接着会准备一些容器里面需要用的配置文件，下面是容器元数据目录的变化

    [root@lichenhao 28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10]# tree
    .
    ├── 28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10-json.log
    ├── checkpoints
    ├── config.v2.json
    ├── hostconfig.json
    ├── hostname
    ├── hosts
    ├── mounts
    ├── resolv.conf
    └── resolv.conf.hash

容器启动后多了几个文件，这几个文件时 docker 动态生成的

- 28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10-json.log #容器的日志文件，后续容器的 stdout 和 stderr 都会输出到这个目录。当然如果配置了其它的日志插件的话，日志就会写到别的地方。
- hostname：里面是容器的主机名，来自于 config.v2.json，由 docker create 命令的-h 参数指定，如果没指定的话，就是容器 ID 的前 12 位，这里即为 28f5bed704dc
- resolv.conf：里面包含了 DNS 服务器的 IP，来自于 hostconfig.json，由 docker create 命令的--dns 参数指定，没有指定的话，docker 会根据容器的网络类型生成一个默认的，一般是主机配置的 DNS 服务器或者是 docker bridge 的 IP。
- resolv.conf.hash：resolv.conf 文件的校验码

Note：除了日志文件外，其它文件在每次容器启动的时候都会自动生成，所以修改他们的内容后只会在当前容器运行的时候生效，容器重启后，配置又都会恢复到默认的状态

## 准备 OCI 所需的 bundle

bundle 被 docker 放在 /run/containerd/ 目录下，展示如下：

    [root@lichenhao containerd]# pwd
    /run/containerd
    [root@lichenhao containerd]# tree
    .
    ├── containerd.sock
    ├── io.containerd.runtime.v1.linux
    │   └── moby
    │       └── 28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10
    │           ├── config.json
    │           ├── init.pid
    │           ├── log.json
    │           └── rootfs
    └── io.containerd.runtime.v2.task

## 准备 IO 文件

容器运行所需的 IO 文件被 docker 放在 /run/docker/containerd/\* 目录下

    [root@lichenhao containerd]# pwd
    /run/docker/containerd
    [root@lichenhao containerd]#  tree
    .
    └── 28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10
        ├── init-stdin
        └── init-stdout

init-stdin 文件用来向容器的 stdin 中写数据，init-stdout 用来接受容器的 stdout。如果使用 echo "XXX" > init-stdin 向容器的标准输入写入内容，则容器会接收该命令，并返回执行结果给 init-stdout。与此同时 cat init-stdout 的话，在宿主机就可以显示容器内在标准输出的内容。

docker exec 命令就是通过这两个文件，来让宿主机与容器进行交互，效果如下：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/xgxt2w/1616121764249-a0867491-440e-4f2d-a52d-186aceee3136.png)

## 正常启动容器

容器正常启动后，会在 /run/docker/runtime-runc/moby/\* 目录中创建该容器的状态文件 state.json 。该文件包含当前容器详细的配置及状态信息。其中也包括 bundle 路径等等。

    [root@lichenhao moby]# pwd
    /run/docker/runtime-runc/moby
    [root@lichenhao moby]# tree
    .
    └── 28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10
        └── state.json

<!---->

    [root@lichenhao moby]# cat 28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10/state.json  | jq .
    {
      "id": "28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10",
      "init_process_pid": 32221,
      "init_process_start": 355109292,
      "created": "2020-06-24T04:03:19.300399652Z",
      "config": {
        "no_pivot_root": false,
        "parent_death_signal": 0,
        "rootfs": "/var/lib/docker/overlay2/d976eddf7575a3464486d92539229146f3df66080a3265195791ebb0d24b24dd/merged",
        "readonlyfs": false,
        "rootPropagation": 0,
        "mounts": [
          {
    .......
          },
    .....
          "bundle=/run/containerd/io.containerd.runtime.v1.linux/moby/28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10"
    ......
    }
