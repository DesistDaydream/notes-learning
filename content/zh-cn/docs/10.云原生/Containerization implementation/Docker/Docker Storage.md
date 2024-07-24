---
title: Docker Storage
linkTitle: Docker Storage
date: 2024-07-05T08:39
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，在生产环境运行你的应用-管理应用数据-存储概述](https://docs.docker.com/storage/)

当关闭并重启 Container 的时候，其内的数据不受影响；但删除 Docker 容器后，则其内对最上面的可写层操作的内容则全部丢失，这时候会存在几个问题

- 存储于联合文件系统中，不易于宿主机访问
- 容器间数据共享不便
- 删除容器会使数据丢失

为了解决这些问题，可以通过三种 Storage 方式来将文件存储于宿主机中

- **volume**
  - volume 类型的 storage 是通过 docker volume 命令显式得创建一个抽象的内容，创建完一个 volume 会，会在 /var/lib/docker/volumes/\* 目录下生成与 volume 同名的目录，在将 volume 挂载进 Container 中时，也就是将 /var/lib/docker/volmes/XXX 目录挂载进去。非 Docker 进程不应修改文件系统的这一部分。卷是在 Docker 中持久保存数据的最佳方法。
- **bind mounts**
  - 可以存储在主机系统上的任何位置。它们甚至可能是重要的系统文件或目录。Docker 主机或 Docker 容器上的非 Docker 进程可以随时对其进行修改。
- **tmpfs mount**
  - 仅存储在主机系统的内存中，并且永远不会写入主机系统的文件系统中。

无论使用哪种方式，目的都是让宿主机上的某个“目录或者文件”绕过联合文件系统，与 Container 中的一个或多个“目录或文件”绑定，对目录中的操作，在 Container 和 Host 中都能看到(i.e.在宿主机目录中创建一个文件，Container 中对应的目录也会看到这个文件)，一个 Volume 可以绑定到多个 Container 上去。

这三种方式唯一的差异就是数据在 docker 宿主机上的位置，bind mount 和 volume 会在宿主机的文件系统中、而 tmpfs mount 则在宿主机的内存中。如下图所示：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/docker/storage/202407241948752.png)

Note：merged(可读写层) 的目录内无法看到这三种存储方式关联到容器中的任何数据，只能从这些存储类型的源目录看到。

数据卷挂载之后的可见性

## Volume

> 参考：
>
> - [官方文档，在生产环境运行你的应用-管理应用数据-卷](https://docs.docker.com/storage/volumes/)

卷是一个逻辑概念，需要手动创建出来之后才能使用，创建出来的卷会关联到一个目录上，对卷的操作就是对该目录的操作，创建出来的卷数据保存在 `/var/lib/docker/volumes/` 路径下，其内目录目录名为卷的名字。

我们可以通过 docker volume 子命令创建卷

### Syntax(语法)

**docker volume COMMAND \[OPTIONS]**

COMMAND：

- create # 创建一个 volume,若不指定 VolumeName，则会自动生成一串一堆字符为名的 volume name
- inspect # 显示一个或多个卷的详细信息
- ls # 列出所有 volume
- prune # 删除所有未使用的 volumes
- rm # 删除一个或多个 volumes

### inspect

**docker volume inspect \[OPTIONS] \[VoumeName]**
显示的信息示例如下：

```json
~]# docker volume inspect volume-test
[
    {
        "CreatedAt": "2020-06-29T13:40:14+08:00",
        "Driver": "local",
        "Labels": {},
        # 卷的挂载点。也就是在容器引用该卷时，所能使用的目录
        "Mountpoint": "/var/lib/docker/volumes/volume-test/_data",
        "Name": "volume-test", # 卷的名称
        "Options": {},
        "Scope": "local"
    }
]
```

## bind mount

> 官方文档：<https://docs.docker.com/storage/bind-mounts/>

bind mount 可以将宿主机上任意目录或者文件，与 Container 共享

就像这个类型的名字一样，将 docker 目录与宿主机目录绑定

## tmpfs mount

> 官方文档：<https://docs.docker.com/storage/tmpfs/>

# 三种类型的 Storage 应用在容器中的方法

使用 docker run|create 命令时，使用 -v 或者 --mount ，则可以将 3 种 Storage 应用在 Container 中

**docker run --mount type=TYPE,src=SRC,dst=DST\[,OPTIONS] ...** # 以 IMAGE 启动运行一个容器，并将 SRC 指定的目录绑定或者挂载到容器内的 DST 目录上。Note：src 还可以写成 source，dst 还可以写成 destination、target。

- TYPE # 挂载类型。可用的有 volume、bind、tmpfs 三种
- SRC # 宿主机路径(/HOST/PATH) 或者容器卷名称(VolumeName)。
  - 若不指定 SRC，则 docker 会自动创建一个。若指定路径在宿主机上不存在，则默认在 /var/lib/docker/volumes/ 路径下创建 volume 所用的目录
  - 挂载文件时，若指定的 SRC 不存在，则会在本机自动创建一个同名目录。
  - tmpfs 不用指定 SRC
- DST # 容器内的路径 /CONTAINER/PATH

OPTIONS # 以逗号分隔的选项列表，以 KEY=VALUE 的方式表示。

- ro=true | false(default) # 指定该 Stroage 是否是只读模式(默认为 false。i.e. rw(读写)模式)。Note: setting readonly for a bind mount does not make its submounts。read-only on the current Linux implementation. See also bind-nonrecursive.
- 适用于 bind 类型 Storage 的选项
  - bind-propagation=shared | slave | private | rshared | rslave | rprivate(default) #
  - consistency=consistent(default) | cached | delegated # 该选项目前仅对 mac 版 docker 有用。
  - bind-nonrecursive=true | false(default). 如果为 true，则子挂载不会递归挂载。这个 OPT 对于只读模式很有用
- 适用于 volume 类型 Storage 的选项
  - volume-driver: Name of the volume-driver plugin.
  - volume-label: Custom metadata.
  - volume-nocopy=true(default) | false. If set to false, the Engine copies existing files and directories under the mount-path into the volume,allowing the host to access them.如果设置为 false，则引擎会将安装路径下的现有文件和目录复制到该卷中，使主机可以访问它们。
  - volume-opt: specific to a given volume driver.
- 适用于 tmpfs 类型 Storage 的选项
  - tmpfs-size=SIZE # tmpfs 类型 Storage 的大小。在 Linux 中默认是无限的。
  - tmpfs-mode=MODE # tmpfs 类型 Storage 的文件模式(e.g. 700 or 0700.)。Linux 中默认为 1777

EXAMPLE

- type=bind,source=/path/on/host,destination=/path/in/container
- type=volume,source=my-volume,destination=/path/in/container,volume-label="color=red",volume-label="shape=round"
- type=tmpfs,tmpfs-size=512M,destination=/path/in/container

docker run -v \[SRC:]DST\[:OPTS] ... # 原始应用 Storage 的语法，新版本推荐使用 --mount 选项。

-v 和 --mount 行为之间的区别

由于 -v 和 --volume 标志已经很长时间成为 Docker 的一部分，因此它们的行为无法更改。这意味着 -v 和 --mount 之间存在一种不同的行为。

如果使用 -v 或 --volume 绑定安装 Docker 主机上尚不存在的文件或目录，请 -v 为您创建端点。始终将其创建为目录。

如果使用 --mount 绑定贴装尚不泊坞窗主机上存在的文件或目录，码头工人也不会自动为您创建它，但会产生一个错误。

Note：在使用 -v 方式应用 Storage，并且指定 SRC 为宿主机上的具体路径 ，则该 Storage 信息会记录在 /var/lib/docker/containers/ID/hostconfig.json 文件的 `.HostConfig.Binds` 字段下，通过 docker inspect ID 命令查看的容器信息中，也会显示在 .HostConfig.Binds 字段下面。

而使用 --mount 方式应用的 Storage，则会将信息记录在 hostconfig.json 文件的 `.HostConfig.Mounts` 下。

挂载一个单独文件的方法

docker run -v /HOST/PATH/FILE:/CONTAINER/PATH/FILE # 这样就把 host 上的 FILE 挂载到 container 指定路径下的 FILE 上去了。FILE 可以不同名

# 修改 docker -v 挂载的文件时遇到的问题

在启动 docker 容器时，为了保证一些基础配置与宿主机保持同步，通常需要将这些配置文件挂载进 docker 容器，例如/etc/resolv.conf、/etc/hosts、/etc/localtime 等文件。

当这些配置变化时，我们通常会修改这些文件。但是此时遇到了一个问题：

- 当在宿主机上修改这些文件后，docker 容器内查看时，这些文件并未发生对应的修改。

然后通过查阅相关资料，发现该问题是由 docker -v 挂载文件和某些编辑器存储文件的行为共同导致 的。

- docker 挂载文件时，并不是挂载了某个文件的路径，而是实打实的挂载了对应的文件，即挂载了某 个指定的 inode 文件。
- 某些编辑器(vi)在编辑保存文件时，采用了备份、替换的策略，即编辑过程中，将变更写入新文件， 保存时，再将备份文件替换原文件，此时会导致文件的 inode 发生变化。
- 原 inode 对应的文件其实并没有发生修改。

因此，我们从宿主机上修改这些文件时，应该采用 echo 重定向等操作，避免文件的 inode 发生变化。
