---
title: Docker Image
linkTitle: Docker Image
date: 2023-11-03T22:25
weight: 4
---

# 概述

> 参考：
>
> - [思否，走进docker(02)：image(镜像)是什么？](https://segmentfault.com/a/1190000009309347)

在虚拟化中，运行程序的地方是一个虚拟的操作系统。而容器技术中，运行程序的地方是一个容器 image(镜像)。

容器的镜像与虚拟机的虚拟系统有异曲同工之妙，基本原理相似，只不过虚拟系统可以像正常安装系统一样进行安装，而容器镜像，则是一个已经打包好的操作系统，可以开箱即用。

由于这种构造，docker 公司研究出一种技术，就是联合文件系统(UnionFS(Union File System))。可以将镜像分为多层(layers)，每层附加一个功能。

实现联合文件系统的驱动程序：docker 本身支持 overlay，overlay2，aufs，btrfs，devicemapper，vfs 等

## Container Image 的分层结构(联合文件系统)

http://www.cnblogs.com/CloudMan6/p/6806193.html

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ma1cb7/1616121959962-5b74016a-126f-4a55-a37f-f93631fd335c.png)

Container Image 采用分层结构，最底层为 bootfs，其它为 rootfs

1. bootfs：用于系统引导的文件系统，包括 bootloader 和 kernel，Container 启动完成后会被卸载以节约内存资源
2. rootfs：位于 bootfs 之上，表现为 Container 的根文件系统；
   1. 传统模式中，系统启动时，内核挂载 rootfs 时会先将其挂载为只读模式，完整性自检完成后将其重新挂载为读写模式。
   2. docker 中，rootfs 由内核挂载为“只读”模式，然后通过“联合挂载”技术额外挂载一个“可写(writable)”层。

可以这么理解：

1. 通过一个 Base Image 启动了一个 Container，然后安装一个 vim 编辑器，commit 这个 container，生成的新镜像就是两层，第一层是系统，第二层是 vim。
2. 这时候用这个这个新的 image 启动一个 Container 后，再安装一个 Nginx，然后 commit 这个 Container，生成的新镜像就是三层，1 系统、2vim、3nginx。
3. 以此类推，每一次 Container 的变化被 commit 后都可以当作一层。
4. 对 Container 的操作产生的变化，是在可写的容器层中进行的

## 可写(writable)的层(layers)

![600](https://notes-learning.oss-cn-beijing.aliyuncs.com/ma1cb7/1616121959993-37ba6cb0-18ec-495b-84c0-72e941c5a240.png)

![800](https://notes-learning.oss-cn-beijing.aliyuncs.com/ma1cb7/1616121960020-2c73e0a5-a7e5-4e71-9907-d098c5233d30.png)

当容器启动时，一个新的“可读写”层被加载到镜像的顶部。这一层通常被称作“容器层”，“容器层”之下的都叫“镜像层”。位于下层的 image 称为父镜像(parent image)，最底层的称为基础镜像(base image)

所有对容器的改动 - 无论添加、删除、还是修改文件都只会发生在容器层中。

只有容器层是“可读写”的，容器层下面的所有镜像层都是“只读”的。

镜像层数量可能会很多，所有镜像层会联合在一起组成一个统一的文件系统。如果不同层中有一个相同路径的文件，比如 /a，上层的 /a 会覆盖下层的 /a，也就是说用户只能访问到上层中的文件 /a。在容器层中，用户看到的是一个叠加之后的文件系统。

1. 添加文件：在容器中创建文件时，新文件被添加到容器层中。
2. 读取文件 ：在容器中读取某个文件时，Docker 会从上往下依次在各镜像层中查找此文件。一旦找到，打开并读入内存。
3. 修改文件 ：在容器中修改已存在的文件时，Docker 会从上往下依次在各镜像层中查找此文件。一旦找到，立即将其复制到容器层，然后修改之。
4. 删除文件 ：在容器中删除文件时，Docker 也是从上往下依次在镜像层中查找此文件。找到后，会在容器层中记录下此删除操作。

只有当需要修改时才复制一份数据，这种特性被称作 Copy-on-Write。可见，容器层保存的是镜像变化的部分，不会对镜像本身进行任何修改；所以，如果多个容器共享一份基础镜像，当某个容器修改了基础镜像的内容，比如 /etc 下的文件，这时其他容器的 /etc 不会被修改，因为修改的只是容器层那个变化的地方，底下的镜像层是不变的

为什么 Docker Image 要采用这种分层结构呢？

最大的一个好处就是 - 共享资源。

比如：有多个镜像都从相同的 base 镜像构建而来，那么 Docker Host 只需在磁盘上保存一份 base 镜像；同时内存中也只需加载一份 base 镜像，就可以为所有容器服务了。而且镜像的每一层都可以被共享。

比如当我获取一个镜像时，可以看到下面的信息

```bash
~]# docker pull lchdzh/network-test:v2.0
v2.0: Pulling from lchdzh/network-test
f34b00c7da20: Pull complete # 镜像第一层
b248a5455a16: Pull complete # 镜像第二层
beaf4c6c50c6: Pull complete # 镜像第三层
Digest: sha256:b27d98887f62c0cf28bc8707ee2de39f8c753afbd047e910e6f1cf2670ae141b
Status: Downloaded newer image for lchdzh/network-test:v2.0
docker.io/lchdzh/network-test:v2.0
```

在获取镜像时可以看到，一共有三个层的镜像需要下载。

而当容器运行时，我们通过 mount 命令可以看到如下内容

```bash
overlay on /var/lib/docker/overlay2/XXXXXXX/merged type overlay (rw,relatime,lowerdir=/var/lib/docker/overlay2/l/XFVMRG3WBD4RLHJ73V5DQHSCZ5:/var/lib/docker/overlay2/l/PZG7BURXDF2DQ4FF54XUEGZE6A:/var/lib/docker/overlay2/l/L4GF2CGQ7LA6LAACXSSE5CMC4P:/var/lib/docker/overlay2/l/6XYKKSJHBNHYSTQMCETDRWL2AA,upperdir=/var/lib/docker/overlay2/7b64f08bef3ca5ab8a2aa0fa0b124b4e55f3f98f421d0cfe7dab271447cb77a2/diff,workdir=/var/lib/docker/overlay2/7b64f08bef3ca5ab8a2aa0fa0b124b4e55f3f98f421d0cfe7dab271447cb77a2/work)
```

这就是那个读写层，这就说明当一个容器运行的时候，会在 docker 的存储类型(这里是 overlay2)目录中，创建一个 XXX/merged 的目录，如果通过 docker exec 进入容器的话，会发现容器中的目录内容与该目录一模一样，并且如果在宿主机上的挂载目录修改文件，同样会影响到容器中。

# Docker Image 的管理标准(OCI 标准)介绍

OCI 规范详见：[OCI Image 规范](/docs/10.云原生/Containerization/Open%20Containers%20Initiative(开放容器倡议)/OCI%20Image%20规范.md)

## docker pull 命令的大概过程

如果对 Image manifest，Image Config 和 Filesystem Layers 等概念不是很了解，请先参考 image(镜像)是什么。

拉取 image 的大概过程如下：

- docker 发送 image 的名称+tag（或者 digest）给 registry 服务器，服务器根据收到的 image 的名称+tag（或者 digest），找到相应 image 的 manifest，然后将 manifest 返回给 docker
- docker 得到 manifest 后，读取里面 image 配置文件的 digest(sha256)，这个 sha256 码就是 image 的 ID
- 根据 ID 在本地找有没有存在同样 ID 的 image，有的话就不用继续下载了
- 如果没有，那么会给 registry 服务器发请求（里面包含配置文件的 sha256 和 media type），拿到 image 的配置文件（Image Config）
- 根据配置文件中的 diff_ids（每个 diffid 对应一个 layer tar 包的 sha256，tar 包相当于 layer 的原始格式），在本地找对应的 layer 是否存在
- 如果 layer 不存在，则根据 manifest 里面 layer 的 sha256 和 media type 去服务器拿相应的 layer（相当去拿压缩格式的包）。
- 拿到后进行解压，并检查解压后 tar 包的 sha256 能否和配置文件（Image Config）中的 diff_id 对的上，对不上说明有问题，下载失败
- 根据 docker 所用的后台文件系统类型，解压 tar 包并放到指定的目录
- 等所有的 layer 都下载完成后，整个 image 下载完成，就可以使用了

注意： 对于 layer 来说，config 文件中 diffid 是 layer 的 tar 包的 sha256，而 manifest 文件中的 digest 依赖于 media type，比如 media type 是 tar+gzip，那 digest 就是 layer 的 tar 包经过 gzip 压缩后的内容的 sha256，如果 media type 就是 tar 的话，diffid 和 digest 就会一样。

dockerd 和 registry 服务器之间的协议为 Registry HTTP API V2。

## OCI Image 中的各种标识符(XXXID)

在上述四部分中，有多个 XXXID 来标识 docker image 的各种信息

1. **imageID** # 镜像的唯一标志，根据镜像的元数据配置文件采用 sha256 算法计算获得
   1. imageID 一般可以在 `${DockerRootDir}/image/${StorageDriver}/repositories.json` 文件中找到
   2. 镜像的 configuration 文件就是以 imageID 命名，一般保存在 `${DockerRootDir}/image/${StorageDriver}/imagedb/content/sha256/` 目录下
2. **diffID** # 镜像层的校验 ID，根据该镜像层的打包文件校验获得
   1. diffID 一般在 configuration 文件的 .rootfs.diff_ids 字段中找到
3. **chainID** # docker 内容寻址机制采用的索引 ID，其值根据当前层和所有父层的 diffID(或父层的 chainID) 计算获得
   1. chainID 计算完成后，一般可以在 `${DockerRootDir}/image/${StorageDriver}/layerdb/sha256/` 目录中找到 chainID 的同名目录
4. **cacheID** # 下载 layer 时、创建容器后产生可写 layers 时，随机生成的 uuid，用于索引镜像层
   1. 在 chainID 的目录中，可以找到 image 的 cache-id 文件，文件内容就是 cacheID。
   2. 然后在 `${DockerRootDir}/${StorageDriver}/` 目录中找到与 cacheID 同名的目录，这些目录中存储了镜像层的所有数据
   3. 至于创建容器后生成的可写 layers 的 cacheID 信息，一般是保存在容器相关的信息文件中的，比如容器的状态文件、容器的可写层的信息、系统的 mount 信息等等地方，都会有相关记录
5. **digest** # 对于某些 image 来说，可能在发布之后还会做一些更新，比如安全方面的，这时虽然镜像的内容变了，但镜像的名称和 tag 没有变，所以会造成前后两次通过同样的名称和 tag 从服务器得到不同的两个镜像的问题，于是 docker 引入了镜像的 digest 的概念，一个镜像的 digest 就是镜像的 manifes 文件的 sha256 码，当镜像的内容发生变化的时候，即镜像的 layer 发生变化，从而 layer 的 sha256 发生变化，而 manifest 里面包含了每一个 layer 的 sha256，所以 manifest 的 sha256 也会发生变化，即镜像的 digest 发生变化，这样就保证了 digest 能唯一的对应一个镜像

# Docker Image 本地路径存放规则

现在通过命令 docker pull ubuntu 获取一个镜像(官方提供的最新 ubuntu 镜像，对应的完整名称为 docker.io/library/ubuntu:latest)

```bash
~]# docker pull ubuntu
Using default tag: latest
latest: Pulling from library/ubuntu
a4a2a29f9ba4: Pull complete
127c9761dcba: Pull complete
d13bf203e905: Pull complete
4039240d2e0b: Pull complete
Digest: sha256:35c4a2c15539c6c1e4e5fa4e554dac323ad0107d8eb5c582d6ff386b383b7dce
Status: Downloaded newer image for ubuntu:latest
docker.io/library/ubuntu:latest
```

在 /var/lib/docker/ 目录中，存放着关于 docker image 的所有信息

首先需要先在 docker 保存 image 的元数据目录(/var/lib/docker/image/)中查找相关信息

## 查找镜像的基础信息，IAMGE ID 与 DIGEST

/var/lib/docker/image/overlay2/repositories.json 文件中记录了和本地 image 相关的 repository 信息，主要是和 image 的名字和 ID 的对应关系。当 image 从 registry 上被 pull 下来之后，就会更新该文件

```bash
~]# cat repositories.json | jq .
{
  "Repositories": {
    # 仓库名
    "ubuntu": {
      # IMAGE NAME ，以及 IMAGE ID(在 docker images 命令中，只显示ID的前12位)
      "ubuntu:latest": "sha256:74435f89ab7825e19cf8c92c7b5c5ebd73ae2d0a2be16f49b3fb81c9062ab303",
      # IMAGE 的 DIGEST，以及 IMAGE ID
      "ubuntu@sha256:35c4a2c15539c6c1e4e5fa4e554dac323ad0107d8eb5c582d6ff386b383b7dce": "sha256:74435f89ab7825e19cf8c92c7b5c5ebd73ae2d0a2be16f49b3fb81c9062ab303"
    }
  }
}
~]# docker images --digests
REPOSITORY          TAG                 DIGEST                                                                    IMAGE ID            CREATED             SIZE
ubuntu              latest              sha256:35c4a2c15539c6c1e4e5fa4e554dac323ad0107d8eb5c582d6ff386b383b7dce   74435f89ab78        4 days ago          73.9MB
~]# docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
ubuntu              latest              74435f89ab78        4 days ago          73.9MB
```

## 根据基础信息查找镜像层(layer)的元数据

根据 repositories.json 文件中的 IMAGE ID ，可以在如下路径中找到该 image 的 Image Configuration 文件

/var/lib/docker/image/overlay2/imagedb/content/sha256/74435f89ab7825e19cf8c92c7b5c5ebd73ae2d0a2be16f49b3fb81c9062ab303 文件就是 OCI 标准的 Image Configuration 文件，从该配置文件的 rootfs 字段中，可以看到该 image 包含 4 个镜像层，从上到下依次是底层到顶层的 diffID

```json
 .....
   "rootfs": {
    "type": "layers",
    "diff_ids": [
      "sha256:e1c75a5e0bfa094c407e411eb6cc8a159ee8b060cbd0398f1693978b4af9af10",
      "sha256:9e97312b63ff63ad98bb1f3f688fdff0721ce5111e7475b02ab652f10a4ff97d",
      "sha256:ec1817c93e7c08d27bfee063f0f1349185a558b87b2d806768af0a8fbbf5bc11",
      "sha256:05f3b67ed530c5b55f6140dfcdfb9746cdae7b76600de13275197d009086bb3d"
    ]
  }
....
```

Note：该对应关系对查找镜像的存储路径没有绝对关系。在 /var/lib/docker/image/overlay2/distribution/ 目录中，保存了镜像层的 digest 与 diffid 对应关系。

- diffid-by-digest： 存放 digest 到 diffid 的对应关系
- v2metadata-by-diffid： 存放 diffid 到 digest 的对应关系

/var/lib/docker/image/overlay2/layerdb/ 目录下存放了所有镜像层的元数据信息，在 /var/lib/docker/image/overlay2/layerdb/sha256/\* 目录下，以镜像层的 chaind 命名。根据 OCI 默认规则，最底下的镜像层 chainid 与 diffid 相同

```bash
~]# pwd
/var/lib/docker/image/overlay2/layerdb/sha256
~]# ll -h
total 0
drwx------ 2 root root 85 Jun 21 19:26 27d46ebb54384edbc8c807984f9eb065321912422b0e6c49d6a9cd8c8b7d8ffc
drwx------ 2 root root 85 Jun 21 19:26 8a8d1f0b34041a66f09e49bdc03e75c2190f606b0db7e08b75eb6747f7b49e11
drwx------ 2 root root 71 Jun 21 19:26 e1c75a5e0bfa094c407e411eb6cc8a159ee8b060cbd0398f1693978b4af9af10
drwx------ 2 root root 85 Jun 21 19:26 f1b8f74eff975ae600be0345aaac8f0a3d16680c2531ffc72f77c5e17cbfeeee
```

```bash
# echo -n "sha256:父层chainID sha256:本层diffID" | sha256sum -
# 根据该命令得出本层的 chainID
# 下面根据最低镜像层计算第二层镜像的 chainid
~]# echo -n "sha256:e1c75a5e0bfa094c407e411eb6cc8a159ee8b060cbd0398f1693978b4af9af10 sha256:9e97312b63ff63ad98bb1f3f688fdff0721ce5111e7475b02ab652f10a4ff97d" | sha256sum -
27d46ebb54384edbc8c807984f9eb065321912422b0e6c49d6a9cd8c8b7d8ffc  -
```

在每个镜像层的目录中，包含了如下几种信息

```bash
~]# cd /var/lib/docker/image/overlay2/layerdb/sha256 && ls
27d46ebb54384edbc8c807984f9eb065321912422b0e6c49d6a9cd8c8b7d8ffc

~]# cd 27d46ebb54384edbc8c807984f9eb065321912422b0e6c49d6a9cd8c8b7d8ffc

27d46ebb54384edbc8c807984f9eb065321912422b0e6c49d6a9cd8c8b7d8ffc]# ls
cache-id  diff  parent  size  tar-split.json.gz

# cache-id 是 docker 在下载镜像层时随机生成的 uuid。也就是说，该文件内容指向了真正存放镜像层文件的地方。
27d46ebb54384edbc8c807984f9eb065321912422b0e6c49d6a9cd8c8b7d8ffc]# cat cache-id
5de7ac8af2fb0a5fb0be4244aa07685bfcfcfc4c4b1c149bc753eb044d7f4a12

# diff 文件存放镜像层的 diffid
27d46ebb54384edbc8c807984f9eb065321912422b0e6c49d6a9cd8c8b7d8ffc]# cat diff
sha256:9e97312b63ff63ad98bb1f3f688fdff0721ce5111e7475b02ab652f10a4ff97d

# parent 文件存放 当前layer 的 父layer 的 diffid。注意：对于最底层的 layer 来说，由于没有 父layer，所以没有这个文件
27d46ebb54384edbc8c807984f9eb065321912422b0e6c49d6a9cd8c8b7d8ffc]# cat parent
sha256:e1c75a5e0bfa094c407e411eb6cc8a159ee8b060cbd0398f1693978b4af9af10

# size 文件存放当前镜像层的大小，单位是字节
27d46ebb54384edbc8c807984f9eb065321912422b0e6c49d6a9cd8c8b7d8ffc]# cat size
1006717
# tar-split.json.gz，layer 压缩包的 split 文件，通过这个文件可以还原 layer 的 tar 包，
# 在 docker save 导出 image 的时候会用到
# 详情可参考 https://github.com/vbatts/tar-split
```

## 根据 cache-id 查找到存放镜像层的路径

/var/lib/docker/overlay2/ 目录存放了所有镜像层的数据。

```bash
~]# pwd
/var/lib/docker/overlay2
~]# ll -h
total 0
drwx------ 4 root root  72 Jun 21 19:26 113a9d8407c2db3892944c17beba7a635ea39aa5108c7f716088466ea302a7e3
# 根据 cache-id 中显示的信息，ubuntu 第二层镜像就是在这个目录中
drwx------ 4 root root  72 Jun 21 19:26 5de7ac8af2fb0a5fb0be4244aa07685bfcfcfc4c4b1c149bc753eb044d7f4a12
drwx------ 4 root root  72 Jun 21 19:28 7704e53a9392b092479707d38b2b183b17bbe2cc220e2283cead9493e19aa651
# 根据 cache-id 中显示的信息，ubuntu 最底层的镜像就是在这个目录中
drwx------ 3 root root  47 Jun 21 19:26 8f377ae99a442b37f5a831724951ce1cf8bfc7b874843c97d09e8027c3dd19e6
drwx------ 2 root root 142 Jun 21 21:01 l
```

# 怎样修改 docker 容器 hosts 文件的内容？

这就要了解 docker 镜像的分层结构了，其中有一个叫 Init 的层，该层专门用来存储一些配置文件，比如：/etc/hosts、/etc/resolv.conf 等信息的，该层并不会跟随镜像一起提交，所以如果我们直接在 Dockerfile 中去覆盖 /etc/hosts 文件的话是不会生效的，要解决这个问题可以有几种方法：
1\. 启动容器的时候(docker run)添加参数—add-host machine:ip 可以实现 hosts 修改，缺点就是如果很多个节点的话命令会很长
2\. 修改容器 hosts de 查找目录，我们可以让容器启动的时候不去找 /etc/hosts 文件，而是去查找我们自己定义的 hosts 文件，下面是一个 Dockerfile 实例：

```dockerfile
FROM ubuntu:14.04
RUN cp /etc/hosts /tmp/hosts  # 路径长度最好保持一致
RUN mkdir -p -- /lib-override && cp /lib/x86_64-linux-gnu/libnss_files.so.2 /lib-override
RUN sed -i 's:/etc/hosts:/tmp/hosts:g' /lib-override/libnss_files.so.2
ENV LD_LIBRARY_PATH /lib-override
RUN echo "192.168.0.1 node1" &gt;&gt; /tmp/hosts  # 可以随意修改/tmp/hosts了
...
```

3. 在 dockerfile 中，使用脚本作为镜像入口，然后利用脚本运行修改 hosts 文件的命令以及真正的应用程序入口，下面是一个 Dockerfile 实例：

```dockerfile
FROM centos:6
RUN mkdir /data
COPY run.sh /data/
COPY myhosts /data/
RUN chmod +x /data/run.sh
ENTRYPOINT /bin/sh -c /data/run.sh
```

其中 run.sh 示例：

```bash
#!/bin/bash
cat /data/myhosts >> /etc/hosts # 向 hosts 文件追加内容
```

其他命令

/bin/bash # 保留终端，防止容器自动退出

然后在 myhosts 文件中添加上你需要添加的 hosts 映射，然后镜像构建完成后，执行 docker run 指令运行容器，查看 /etc/hosts 配置是否生效。

这个问题最重要的就是要理解 docker 镜像的分层结构，由只读层+可读写层+ init 层组成。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ma1cb7/1616121959946-c34b9fbf-7490-4ef8-919e-433e1c41f5b8.png)
