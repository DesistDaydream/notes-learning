---
title: OCI Image 规范
---

# 概述

> 参考：
> - [GitHub 项目,opencontainers/image-spec/spec.md](https://github.com/opencontainers/image-spec/blob/main/spec.md)
> - [思否大佬](https://segmentfault.com/a/1190000009309347)
> - <https://mp.weixin.qq.com/s/8wAv87DkJjE6fVEEmoQ60Q>
> - <https://blog.k8s.li/Exploring-container-image.html>

OCI Image 规范的目的，是为了让其他人按照规范创建交互工具，这个工具应该可以 **building(构建)**、**transporting(传输)**、**running(运行)** 一个容器镜像。

一个 OCI Image 应该由一个 Image Manifest、一个 Image Index(可选)、一组文件系统层、一个配置文件 组成。

本质上，镜像的每一层就是一个一个的 tar.gz 的文件，当各种容器工具 pull 镜像时，会根据各种元数据文件，获取到这些 tar.gz 文件，下载到本地，并根据自身的实现解压他们。

## OCI Image 规范的组件

前文所描述的组成 OCI Image 规范的多个组成部分，又被细分为如下 **Components(组件)**：

- [Image Layout](https://github.com/opencontainers/image-spec/blob/master/image-layout.md) # 镜像内容的文件系统布局。说白了，镜像的主要内容就在这里。
  - [Image Manifest](https://github.com/opencontainers/image-spec/blob/master/manifest.md) # 描述构成容器镜像所具有的组件的文件。比如这个镜像有哪些 layer，额外的 annotation 信息。manifest 文件中保存了很多和当前平台有关的信息
  - [Image Configuration](https://github.com/opencontainers/image-spec/blob/master/config.md) # 一个文档，该文档确定适用于转换为 [runtime bundle](https://github.com/opencontainers/runtime-spec) 运行时包的映像的层顺序和配置。保存了文件系统的层级信息（每个层级的 hash 值，以及历史信息），以及容器运行时需要的一些信息（比如环境变量、工作目录、命令参数、mount 列表），指定了镜像在某个特定平台和系统的配置。比较接近我们使用 docker inspect   看到的内容
  - [Image Index](https://github.com/opencontainers/image-spec/blob/master/image-index.md) # 带注释的图像清单索引。指向不同平台的 manifest 文件，这个文件能保证一个镜像可以跨平台使用，每个平台拥有不同的 manifest 文件，使用 index 作为索引
  - [Filesystem Layer changeset](https://github.com/opencontainers/image-spec/blob/main/layer.md) # 描述容器文件系统的变更集。以 layer 保存的文件系统，每个 layer 保存了和上层之间变化的部分，layer 应该保存哪些文件，怎么表示增加、修改和删除的文件等
- [Conversion](https://github.com/opencontainers/image-spec/blob/master/conversion.md) # 描述此翻译应如何发生。a document describing how this translation should occur
- [Descriptor](https://github.com/opencontainers/image-spec/blob/master/descriptor.md) # 描述所引用内容的类型，元数据和内容地址的引用。a reference that describes the type, metadata and content address of referenced content

Future versions of this specification may include the following OPTIONAL features:

- Signatures that are based on signing image content address
- Naming that is federated based on DNS and can be delegated

OCI Image 的所有组件其实都是一个个的文件，这些文件的名称都是其内容执行 sha256 后的值。每个文件都有一个 [**OCI Image Media Types**](https://github.com/opencontainers/image-spec/blob/master/media-types.md)**(OCI 镜像媒体类型)**，在官方文档中详细介绍了规范中各个组件的媒体类型

### Media Types(媒体类型) # OCI Image 组件文件的打包格式

OCI Image 中的每个组件都会打包成一个文件。做过 web 开发的程序员对 media type 应该比较熟悉，简单点说，就是当客户端用 http 协议下载一个文件的时候，需要在 http 的首部带上 Accept 字段，告诉服务器端它支持哪些类型的文件，服务器返回文件的时候，需要在 http 的首部带上 Content-Type 字段，告诉客户端返回文件的类型，如 Accept: text/html,application/xml 和 Content-Type: text/html

**OCI Media Type 文件类型:**

| Media Type                                                   | 说明                                  |
| ------------------------------------------------------------ | ------------------------------------- |
| application/vnd.oci.descriptor.v1+json                       | Content Descriptor 内容描述文件       |
| application/vnd.oci.layout.header.v1+json                    | OCI Layout 布局描述文件               |
| application/vnd.oci.image.index.v1+json                      | Image Index 高层次的镜像元信息文件    |
| application/vnd.oci.image.manifest.v1+json                   | Image Manifest 镜像元信息文件         |
| application/vnd.oci.image.config.v1+json                     | Image Config 镜像配置文件             |
| application/vnd.oci.image.layer.v1.tar                       | Image Layer 镜像层文件                |
| application/vnd.oci.image.layer.v1.tar+gzip                  | Image Layer 镜像层文件 gzip 压缩      |
| application/vnd.oci.image.layer.nondistributable.v1.tar      | Image Layer 非内容寻址管理            |
| application/vnd.oci.image.layer.nondistributable.v1.tar+gzip | Image Layer, gzip 压缩 非内容寻址管理 |

# Image Layout

在开始介绍 layout 之前，先来回顾一下上一篇介绍 hello-world 时提到的从 register 服务器拉 image 的过程：

- 首先获取 image 的 manifests
- 根据 manifests 文件中 config 的 sha256 码，得到 image config 文件
- 遍历 manifests 里面的所有 layer，根据其 sha256 码在本地找，如果找到对应的 layer，则跳过，否则从服务器取相应 layer 的压缩包
- 等上面的所有步骤完成后，就会拼出完整的 image

从上面的过程中可以看出，我们从服务器上取 image 的时候不需要知道 image manifests 和 config 文件的名字，也不需要知道 layer 压缩包的名字。

那么 image 从服务器拉下来后，在本地应该怎么存储呢？文件名称和目录结构应该是怎样的呢？OCI 也有相应的标准，名字叫 [image layout](https://github.com/opencontainers/image-spec/blob/master/image-layout.md)，有了这样的标准之后，我们就可以将整个 image 打成一个包，方便的在不同机器，不同容器平台之间导入导出。

不过遗憾的是，OCI 的这个标准还在变化中，根据 github 上所看到的，v1.0.0-rc5 在 v1.0.0-rc4 上就有较大的修改，并且现在 docker 也不支持该标准。

> docker 对 OCI image layout 的支持还在开发中

下面是 [v1.0 版本的镜像布局示例](https://github.com/opencontainers/image-spec/blob/v1.0/image-layout.md#example-layout)：

> 这里以通过 `nerdctl image save lchdzh/k8s-debug:v1 -o k8s-debug.tar` 命令将 lchdzh/k8s-debug 镜像打包，打包后再通过 `tar` 命令解包获取 OCI 格式的镜像文件。该镜像的构建详见[ kubernetes 的故障处理技巧](/docs/IT学习笔记/10.云原生/2.3.Kubernetes%20 容器编排系统/Kubernetes%20 管理/性能优化%20 与%20 故障处理/故障处理技巧.md 容器编排系统/Kubernetes 管理/性能优化 与 故障处理/故障处理技巧.md)

```bash
~]# tree
.
├── blobs
│   └── sha256
│       ├── 02daccf1684b499e99c258348d492c5f0ea086174d2f0d430791d4f902ae4f71
│       ├── 188c0c94c7c576fff0792aca7ec73d67a2f7f4cb3a6e53a84559337260b36964
│       ├── 5f9b9d9c910519d9a4b1e06f031672e14acf9bcc288ed7e3ed3842916ed4394d
│       ├── c690d4fd64d6622c3721a1db686c2e4cfb559dd1d9f9ff825584a8f56ec02c7f
│       ├── df727e3daae2c57da7071b4056d328d4bbb9d6a913e469d8f07b58e35a5cff96
│       └── ee24b921ba004624b350e7f140e68c6a7d8297bb815b4ca526979a7e66cec15a
├── index.json
├── manifest.json # 不用关注这里的这个文件，版本问题导致还有，其实应该没有
└── oci-layout
```

> 注意，这里 Image Layout 的根目录中多了一个 manifest.json 的目录，这是一个历史遗留问题，这里根目录中的 manifest.json 文件是旧 OCI 版本的 Image Manifest 文件。但是实际上，新版 OCI 标准的 Image Layout 的根目录中，不必包含 maniest.json 文件，该文件已经在 blobs 中了，其中的 Tags 信息则在 index.json 文件中。

可以看到，一个镜像通常由这几部分组成：

- **blobs 目录**
- **oci-layout 文件**
- **index.json 文件(单平台是可选的)**

## blobs 目录

blobs 目录可以说是一个镜像的核心内容，这些内容被组织成一个一个的 [Binary Large Object(二进制大对象，简称 blob)](https://en.wikipedia.org/wiki/Binary_large_object)，这些 blob 文件分为三类

- Image Manifest 文件
- Image Configuration 文件
- Image Filesystem Layer 文件
  - 这是一个被打包、压缩后的镜像文件系统。也就是镜像中的所有实体文件。

每个文件名都是其内容的 sha256 码

```bash
~]# sha256sum ./blobs/sha256/*
02daccf1684b499e99c258348d492c5f0ea086174d2f0d430791d4f902ae4f71  ./blobs/sha256/02daccf1684b499e99c258348d492c5f0ea086174d2f0d430791d4f902ae4f71
188c0c94c7c576fff0792aca7ec73d67a2f7f4cb3a6e53a84559337260b36964  ./blobs/sha256/188c0c94c7c576fff0792aca7ec73d67a2f7f4cb3a6e53a84559337260b36964
5f9b9d9c910519d9a4b1e06f031672e14acf9bcc288ed7e3ed3842916ed4394d  ./blobs/sha256/5f9b9d9c910519d9a4b1e06f031672e14acf9bcc288ed7e3ed3842916ed4394d
c690d4fd64d6622c3721a1db686c2e4cfb559dd1d9f9ff825584a8f56ec02c7f  ./blobs/sha256/c690d4fd64d6622c3721a1db686c2e4cfb559dd1d9f9ff825584a8f56ec02c7f
df727e3daae2c57da7071b4056d328d4bbb9d6a913e469d8f07b58e35a5cff96  ./blobs/sha256/df727e3daae2c57da7071b4056d328d4bbb9d6a913e469d8f07b58e35a5cff96
ee24b921ba004624b350e7f140e68c6a7d8297bb815b4ca526979a7e66cec15a  ./blobs/sha256/ee24b921ba004624b350e7f140e68c6a7d8297bb815b4ca526979a7e66cec15a
```

看一下每个文件的类型：

```bash
~]# file debian/blobs/sha256/*
./blobs/sha256/02daccf1684b499e99c258348d492c5f0ea086174d2f0d430791d4f902ae4f71: gzip compressed data, original size modulo 2^32 30614528
./blobs/sha256/188c0c94c7c576fff0792aca7ec73d67a2f7f4cb3a6e53a84559337260b36964: gzip compressed data, original size modulo 2^32 5843456
./blobs/sha256/5f9b9d9c910519d9a4b1e06f031672e14acf9bcc288ed7e3ed3842916ed4394d: gzip compressed data, original size modulo 2^32 6679552
./blobs/sha256/c690d4fd64d6622c3721a1db686c2e4cfb559dd1d9f9ff825584a8f56ec02c7f: JSON data
./blobs/sha256/df727e3daae2c57da7071b4056d328d4bbb9d6a913e469d8f07b58e35a5cff96: gzip compressed data, original size modulo 2^32 105835008
./blobs/sha256/ee24b921ba004624b350e7f140e68c6a7d8297bb815b4ca526979a7e66cec15a: JSON data
```

其中 gzip 类型的就是镜像层的真实数据。另外两个一个是 Image Manifest 文件，一个是 Image Configuration 文件

### Image Manifest 文件

符合 [Image Manifest 标准](#9b1dd83b)的文件。这个入口文件描述了 OCI 镜像的实际配置和其中的 Layer 配置。如果有多层那 layers 数组中的元素也会相应增加。

```json
~]# cat  blobs/sha256/ee24b921ba004624b350e7f140e68c6a7d8297bb815b4ca526979a7e66cec15a | jq .
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
  "config": {
    "mediaType": "application/vnd.docker.container.image.v1+json",
    "size": 2145,
    "digest": "sha256:c690d4fd64d6622c3721a1db686c2e4cfb559dd1d9f9ff825584a8f56ec02c7f"
  },
  "layers": [
    {
      "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
      "size": 2796860,
      "digest": "sha256:188c0c94c7c576fff0792aca7ec73d67a2f7f4cb3a6e53a84559337260b36964"
    },
    {
      "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
      "size": 33480758,
      "digest": "sha256:df727e3daae2c57da7071b4056d328d4bbb9d6a913e469d8f07b58e35a5cff96"
    },
    {
      "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
      "size": 2226823,
      "digest": "sha256:5f9b9d9c910519d9a4b1e06f031672e14acf9bcc288ed7e3ed3842916ed4394d"
    },
    {
      "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
      "size": 12125653,
      "digest": "sha256:02daccf1684b499e99c258348d492c5f0ea086174d2f0d430791d4f902ae4f71"
    }
  ]
}

```

### Image Configuration 文件

符合 [Image Configuration 标准](#iJRmF)的文件

```json
~]# cat  blobs/sha256/c690d4fd64d6622c3721a1db686c2e4cfb559dd1d9f9ff825584a8f56ec02c7f | jq .
{
  "architecture": "amd64",
  "config": {
    "Env": [
      "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
      "TZ=Asia/Shanghai",
      "LC_ALL=C.UTF-8",
      "LANG=C.UTF-8",
      "LANGUAGE=C.UTF-8"
    ],
    "Entrypoint": [
      "/bin/bash"
    ],
    "ArgsEscaped": true,
    "OnBuild": null
  },
  "created": "2021-08-17T16:38:31.174884751+08:00",
  "history": [
    {
      "created": "2020-10-22T02:19:24.33416307Z",
      "created_by": "/bin/sh -c #(nop) ADD file:f17f65714f703db9012f00e5ec98d0b2541ff6147c2633f7ab9ba659d0c507f4 in / "
    },
    {
      "created": "2020-10-22T02:19:24.499382102Z",
      "created_by": "/bin/sh -c #(nop)  CMD [\"/bin/sh\"]",
      "empty_layer": true
    },
    {
      "created": "2021-08-17T16:04:08.177778645+08:00",
      "created_by": "RUN /bin/sh -c sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories &&     apk update &&     apk add --no-cache vim bash tcpdump curl wget strace mysql-client iproute2 redis jq iftop tzdata tar nmap bind-tools htop &&     ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime # buildkit",
      "comment": "buildkit.dockerfile.v0"
    },
    {
      "created": "2021-08-17T16:04:29.900220581+08:00",
      "created_by": "RUN /bin/sh -c wget -O /usr/bin/httpstat https://github.com/davecheney/httpstat/releases/download/v1.0.0/httpstat-linux-amd64-v1.0.0 &&     chmod +x /usr/bin/httpstat # buildkit",
      "comment": "buildkit.dockerfile.v0"
    },
    {
      "created": "2021-08-17T16:38:31.174884751+08:00",
      "created_by": "COPY /go/bin/grpcurl /usr/bin/grpcurl # buildkit",
      "comment": "buildkit.dockerfile.v0"
    },
    {
      "created": "2021-08-17T16:38:31.174884751+08:00",
      "created_by": "ENV TZ=Asia/Shanghai LC_ALL=C.UTF-8 LANG=C.UTF-8 LANGUAGE=C.UTF-8",
      "comment": "buildkit.dockerfile.v0",
      "empty_layer": true
    },
    {
      "created": "2021-08-17T16:38:31.174884751+08:00",
      "created_by": "ENTRYPOINT [\"/bin/bash\"]",
      "comment": "buildkit.dockerfile.v0",
      "empty_layer": true
    }
  ],
  "os": "linux",
  "rootfs": {
    "type": "layers",
    "diff_ids": [
      "sha256:ace0eda3e3be35a979cec764a3321b4c7d0b9e4bb3094d20d3ff6782961a8d54",
      "sha256:4041c1a8637589d2c872e14d1068376c5e21bf96a837fa2225f91066e84b1e55",
      "sha256:6f5211c02ff0b7e40b9ca7c5f62cc8732647b046e22cc5046053412d1fef97f6",
      "sha256:6e63a43fa96c6ea85d34c23db4c28b76ecda01c03aa721f6a3355b04501bdc58"
    ]
  }
}
```

### Layers 文件

符合 [Filesystem Layers 标准](#Qy2Oy)的文件。对于 Layers 类型的 blob 来说，这个文件的格式可以是 application/vnd.oci.image.layer.v1.tar 和 application/vnd.oci.image.layer.v1.tar+gzip 两种中的一种。

同时标准还定义了 application/vnd.oci.image.layer.nondistributable.v1.tar 和 application/vnd.oci.image.layer.nondistributable.v1.tar+gzip 这两种对应于 nondistributable 的格式，其实这两种格式和前两种格式包含的内容是一样的，只是用不同的类型名称来区分它们的用途，对于名称中有 nondistributable 的 layer，标准要求这种类型的 layer 不能上传，只能下载。

解压后会得到一个 rootfs。这就是这些镜像层的最真实内容

```bash
~]# mkdir -p layers/{1,2,3,4}
~]# tar -xf blobs/sha256/02daccf1684b499e99c258348d492c5f0ea086174d2f0d430791d4f902ae4f71 -C layers/1/
~]# tar -xf blobs/sha256/188c0c94c7c576fff0792aca7ec73d67a2f7f4cb3a6e53a84559337260b36964 -C layers/2/
~]# tar -xf blobs/sha256/5f9b9d9c910519d9a4b1e06f031672e14acf9bcc288ed7e3ed3842916ed4394d -C layers/3/
~]# tar -xf blobs/sha256/df727e3daae2c57da7071b4056d328d4bbb9d6a913e469d8f07b58e35a5cff96 -C layers/4/
~]# tree -L 2 layers/
layers/
├── 1
│   └── usr
├── 2
│   ├── bin
│   ├── dev
│   ├── etc
│   ├── home
│   ├── lib
│   ├── media
│   ├── mnt
│   ├── opt
│   ├── proc
│   ├── root
│   ├── run
│   ├── sbin
│   ├── srv
│   ├── sys
│   ├── tmp
│   ├── usr
│   └── var
├── 3
│   ├── etc
│   ├── root
│   └── usr
└── 4
    ├── bin
    ├── etc
    ├── lib
    ├── run
    ├── sbin
    ├── usr
    └── var
```

解压到 layer/2 目录中的内容，就是 Dockerfile 中的：

```dockerfile
FROM alpine:latest
```

解压到 layer/4 目录中的内容，就是 Dockerfile 中的：

```dockerfile
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    apk add --no-cache vim bash tcpdump curl wget strace mysql-client iproute2 redis jq iftop tzdata tar nmap bind-tools htop && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
```

解压到 layer/3 目录中的内容，就是 Dockerfile 中的：

```dockerfile
RUN wget -O /usr/bin/httpstat https://github.com/davecheney/httpstat/releases/download/v1.0.0/httpstat-linux-amd64-v1.0.0 && \
    chmod +x /usr/bin/httpstat
```

解压到 layer/1 目录中的内容，就是 Dockerfile 中的：

```dockerfile
COPY --from=grpcurl  /go/bin/grpcurl /usr/bin/grpcurl
```

#### 通过镜像层文件启动容器

在 OCI Runtime 规范中的 [Filesystem Bundle 示例](/docs/IT学习笔记/10.云原生/2.1.容器/Open%20Containers%20Initiative(开放容器倡议)/OCI%20Runtime%20 规范.md Runtime 规范.md)中，我们可以直接通过 Layers 文件以及 runc 工具，直接启动一个标准的符合 OCI 规范的简单容器。

## oci-layout 文件

一个 JSON 格式的文件，包含 image 标准的版本信息。也就是用于说明本镜像遵循的 OCI 标准的版本号

```bash
~]# cat oci-layout
{"imageLayoutVersion":"1.0.0"}
```

## index.json 文件

一个符合 [Image Index 标准](#jCDjJ)的 JSON 格式的文件，关于本镜像布局的索引信息。

index.json 文件更多的是用来定义一个镜像的多平台信息，比如，一个镜像在 Linux、Windows 或者 amd64、arm64 这种平台上运行时，应该使用的 Image Manifest 文件。同时还包含该镜像的名称以及 Tag。

```json
~]# cat index.json  | jq .
{
  "schemaVersion": 2,
  "manifests": [
    {
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
      "digest": "sha256:ee24b921ba004624b350e7f140e68c6a7d8297bb815b4ca526979a7e66cec15a",
      "size": 1163,
      "annotations": {
        "io.containerd.image.name": "docker.io/lchdzh/k8s-debug:v1",
        "org.opencontainers.image.ref.name": "v1"
      }
    }
  ]
}
```

并不像示例中的文件格式一模一样，具体原因未知，估计是版本问题。也有可能是只有 Linxu amd64 这一个平台的版本。

## 总结

- 各种实现容器的工具自身是否还维护了一个镜像名称与 Index 文件的对应关系？待确认
- 根据 Index 文件中的 `manifests.digest` 字段找到本平台的 Manifest 文件
- 根据 Manifest 文件中的 `config.digest` 字段找到 Configuration 文件
- 根据 Configuration，找到 Image Configuration，即可将 Manifest 中指定的各种 layers 组合起来，形成一个镜像，并根据其中的环境变量等信息，使用该镜像启动一个容器。

# Filesystem Layers

Filesystem Layer 包含了文件系统的信息，即该 image 包含了哪些文件，以及它们的属性和数据。比如在某一层增加了一个文件，那么这一层所包含的内容就是增加的这个文件的数据以及它的属性

每个 filesystem layer 都包含了在上一个 layer 的改动情况，主要包含三方面的内容：

- 变化类型：是增加、修改还是删除了文件
- 文件类型：每个变化发生在哪种文件类型上
- 文件属性：文件的修改时间、用户 ID、组 ID、RWX 权限等

最终，每个 Layers 都会被打包成如下几种媒体类型其中之一：

- application/vnd.oci.image.layer.v1.tar # 通常都是这种媒体类型
- application/vnd.oci.image.layer.v1.tar+gzip # 通常都是这种媒体类型
- application/vnd.oci.image.layer.v1.tar+zstd
- application/vnd.oci.image.layer.nondistributable.v1.tar
- application/vnd.oci.image.layer.nondistributable.v1.tar+gzip
- application/vnd.oci.image.layer.nondistributable.v1.tar +zstd

也就是说，一个 OCI Image 中每个镜像层的“文件系统”或“文件系统的更改”层在打包完成后，会被序列化为一个 Layers 类型的 blob 文件

# Image Manifest

Image Manifest 是一个 json 文件，media type 为 application/vnd.oci.image.manifest.v1+json，这个文件包含了对前面 filesystem layers 和 image config 的描述

manifest 文件中 config 的 sha256 就是 image 的 ID，即上面 image config 文件的 sha256 值，这里是 b5b2b2c507a0944348e0303114d8d93aaaa081732b86451d9bce1f432a537bc7

```json
{
  "schemaVersion": 2,
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "size": 7023,
    "digest": "sha256:b5b2b2c507a0944348e0303114d8d93aaaa081732b86451d9bce1f432a537bc7"
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 32654,
      "digest": "sha256:9834876dcfb05cb167a5c24953eba58c4ac89b1adf57f28f2f9d09af107ee8f0"
    },
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 16724,
      "digest": "sha256:3c3a4604a545cdc127456d94e421cd355bca5b528f4a9c1905b15da2eb4a4c6b"
    },
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 73109,
      "digest": "sha256:ec4b8955958665577945c89419d1af06b5f7636b4ac3da7f12184802ad867736"
    }
  ],
  "annotations": {
    "com.example.key1": "value1",
    "com.example.key2": "value2"
  }
}
```

- config 里面包含了对 image config 文件的描述，有 media type，文件大小，以及 sha256 码
- layers 包含了对每一个 layer 的描述，和对 config 文件的描述一样，也包含了 media type，文件大小，以及 sha256 码

这里 layer 的 sha256 和 image config 文件中的 diff_ids 有可能不一样，比如这里的 layer 文件格式是 tar+gzip，那么这里的 sha256 就是 tar+gzip 包的 sha256 码，而 diff_ids 是 tar+gzip 解压后 tar 文件的 sha256 码

# Image Configuration

image config 就是一个 json 文件，它的 media type 是 application/vnd.oci.image.config.v1+json，这个 json 文件包含了对这个 image 的描述。

```json
{
  "created": "2015-10-31T22:22:56.015925234Z",
  "author": "Alyssa P. Hacker <alyspdev@example.com>",
  "architecture": "amd64",
  "os": "linux",
  "config": {
    "User": "alice",
    "ExposedPorts": {
      "8080/tcp": {}
    },
    "Env": [
      "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
      "FOO=oci_is_a",
      "BAR=well_written_spec"
    ],
    "Entrypoint": ["/bin/my-app-binary"],
    "Cmd": ["--foreground", "--config", "/etc/my-app.d/default.cfg"],
    "Volumes": {
      "/var/job-result-data": {},
      "/var/log/my-app-logs": {}
    },
    "WorkingDir": "/home/alice",
    "Labels": {
      "com.example.project.git.url": "https://example.com/project.git",
      "com.example.project.git.commit": "45a939b2999782a3f005621a8d0f29aa387e1d6b"
    }
  },
  "rootfs": {
    "diff_ids": [
      "sha256:c6f988f4874bb0add23a778f753c65efe992244e148a1d2ec2a8b664fb66bbd1",
      "sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef"
    ],
    "type": "layers"
  },
  "history": [
    {
      "created": "2015-10-31T22:22:54.690851953Z",
      "created_by": "/bin/sh -c #(nop) ADD file:a3bc1e842b69636f9df5256c49c5374fb4eef1e281fe3f282c65fb853ee171c5 in /"
    },
    {
      "created": "2015-10-31T22:22:55.613815829Z",
      "created_by": "/bin/sh -c #(nop) CMD [\"sh\"]",
      "empty_layer": true
    },
    {
      "created": "2015-10-31T22:22:56.329850019Z",
      "created_by": "/bin/sh -c apk add curl"
    }
  ]
}
```

这里只介绍几个比较重要的属性，其它的请参考[**标准文档**](https://github.com/opencontainers/image-spec/blob/main/config.md)

- **architecture** # CPU 架构类型，现在大部分都是 amd64，不过 arm64 估计会慢慢多起来
- **os** # 操作系统，本人只用过 linux
- **config** # 当根据这个 image 启动 container 时，config 里面的配置就是运行 container 时的默认参数，在后续介绍 runtime 的时候再仔细介绍每一项的意义
- **rootfs** # 指定了 image 所包含的 filesystem layers，type 的值必须是 layers，diff_ids 包含了 layer 的列表（顺序排列），每一个 sha256 就是每层 layer 对应 tar 包的 sha256 码

## OCI Image 中的各种标识符(XXXID)

有多个 XXXID 来标识 OCI Image 的各种信息

- **ImageID **# 镜像的唯一标志，值为 Image Configuration 文件通过 sha256 计算的结果
  - imageID 对于 Docker 来说一般可以在 ${DockerRootDir}/image/${StorageDriver}/repositories.json 文件中找到
  - 镜像的 configuration 文件就是以 imageID 命名，对于 Docker 来说 一般保存在 ${DockerRootDir}/image/${StorageDriver}/imagedb/content/sha256/ 目录下
- **Layer DiffID **# 镜像层的校验 ID，根据该镜像层的打包文件校验获得
  - diffID 一般在 configuration 文件的 `.rootfs.diff_ids` 字段中找到
- **Layer ChainID **# docker 内容寻址机制采用的索引 ID，其值根据当前层和所有父层的 diffID(或父层的 chainID) 计算获得
  - chainID 计算完成后，对于 Docker 来说一般可以在 ${DockerRootDir}/image/${StorageDriver}/layerdb/sha256/ 目录中找到 chainID 的同名目录
- **digest** # 对于某些 image 来说，可能在发布之后还会做一些更新，比如安全方面的，这时虽然镜像的内容变了，但镜像的名称和 tag 没有变，所以会造成前后两次通过同样的名称和 tag 从服务器得到不同的两个镜像的问题，于是 docker 引入了镜像的 digest 的概念，一个镜像的 digest 就是镜像的 manifes 文件的 sha256 码，当镜像的内容发生变化的时候，即镜像的 layer 发生变化，从而 layer 的 sha256 发生变化，而 manifest 里面包含了每一个 layer 的 sha256，所以 manifest 的 sha256 也会发生变化，即镜像的 digest 发生变化，这样就保证了 digest 能唯一的对应一个镜像

# Image Index(可选)

[image index](https://github.com/opencontainers/image-spec/blob/master/image-index.md) 也是个 json 文件，media type 是 application/vnd.oci.image.index.v1+json。

```json
{
  "schemaVersion": 2,
  "manifests": [
    {
      "mediaType": "application/vnd.oci.image.index.v1+json",
      "size": 7143,
      "digest": "sha256:0228f90e926ba6b96e4f39cf294b2586d38fbb5a1e385c05cd1ee40ea54fe7fd",
      "annotations": {
        "org.opencontainers.image.ref.name": "stable-release"
      }
    },
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "size": 7143,
      "digest": "sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f",
      "platform": {
        "architecture": "ppc64le",
        "os": "linux"
      },
      "annotations": {
        "org.opencontainers.image.ref.name": "v1.0"
      }
    },
    {
      "mediaType": "application/xml",
      "size": 7143,
      "digest": "sha256:b3d63d132d21c3ff4c35a061adf23cf43da8ae054247e32faa95494d904a007e",
      "annotations": {
        "org.freedesktop.specifications.metainfo.version": "1.0",
        "org.freedesktop.specifications.metainfo.type": "AppStream"
      }
    }
  ],
  "annotations": {
    "com.example.index.revision": "r124356"
  }
}
```

其实到 manifest 为止，已经有了整个 image 的完整描述，为什么还需要 image index 这个文件呢？主要原因是 manifest 描述的 image 只能支持一个平台，也没法支持多个 tag，加上 index 文件的目的就是让这个 image 能支持多个平台和多 tag。

image index 是 v1.0.0-rc5 才加进来的一个文件，并且 docker 到现在也不支持该文件

index 文件包含了对 image 中所有 manifest 的描述，相当于一个 manifest 列表，包括每个 manifest 的 media type，文件大小，sha256 码，支持的平台以及平台特殊的配置。

比如 ubuntu 想让它的 image 支持 amd64 和 arm64 平台，于是它在两个平台上都编译好相应的包，然后将两个平台的 layer 都放到这个 image 的 filesystem layers 里面，然后写两个 config 文件和两个 manifest 文件，再加上这样一个描述不同平台 manifest 的 index 文件，就可以让这个 image 支持两个平台了，两个平台的用户可以使用同样的命令得到自己平台想要的那些 layer。

image index 最新的标准里面并没有涉及到 tag，不过估计后续会加上。

# 容器镜像总结

Image Layout 规范是容器镜像的最主要部分，其他的规范，都是对 Image Layout 规范的扩展和补充，不管是 docker、containerd 还是什么其他的容器化实现程序，在执行 pull 命令时，都是获取 OCI 标准的 Image Layout，这样就等于是获取到了镜像文件系统的每一层，以及索引数据。

然后将其中 gzip 类型的镜像层文件解压缩到本地文件系统上，这些程序再自己实现一套适用于自己的关联系统，将镜像层组合起来。

当我们使用镜像运行容器时，就是根据索引数据，找到每一层，通过联合挂载的方式，形成容器的文件系统，再为其分配一个 mount namespace。然后再根据具体配置，为进程分配其他所需的名称空间即可。

# 通用规范

## Annotations(注释)

在 Image Manifests 和 Descriptors 规范中具有 annotations 字段，可用于为镜像添加一些标识符以便其他服务识别。比如 GitHub Package 将会根据 `org.opencontainers.image.source` 注释将镜像归属到指定的仓库中。

> 在 Dockerfile 中可以通过 LABEL 关键字为镜像添加注释

OCI 中预定义了一些常用注释以用于镜像索引或识别镜像作者

- **org.opencontainers.image.created** date and time on which the image was built (string, date-time as defined by [RFC 3339](https://tools.ietf.org/html/rfc3339#section-5.6)).
- **org.opencontainers.image.authors** contact details of the people or organization responsible for the image (freeform string)
- **org.opencontainers.image.url** URL to find more information on the image (string)
- **org.opencontainers.image.documentation** URL to get documentation on the image (string)
- **org.opencontainers.image.source** URL to get source code for building the image (string)
- **org.opencontainers.image.version** version of the packaged software
  - The version MAY match a label or tag in the source code repository
  - version MAY be [Semantic versioning-compatible](https://semver.org/)
- **org.opencontainers.image.revision** Source control revision identifier for the packaged software.
- **org.opencontainers.image.vendor** Name of the distributing entity, organization or individual.
- **org.opencontainers.image.licenses** License(s) under which contained software is distributed as an [SPDX License Expression](https://spdx.org/spdx-specification-21-web-version#h.jxpfx0ykyb60).
- **org.opencontainers.image.ref.name** Name of the reference for a target (string).
  - SHOULD only be considered valid when on descriptors on index.json within [image layout](https://github.com/opencontainers/image-spec/blob/main/image-layout.md).
  - Character set of the value SHOULD conform to alphanum of A-Za-z0-9 and separator set of -.\_:@/+
  - The reference must match the following [grammar](https://github.com/opencontainers/image-spec/blob/main/considerations.md#ebnf):ref ::= component ("/" component)_ component ::= alphanum (separator alphanum)_ alphanum ::= \[A-Za-z0-9]+ separator ::= \[-.\_:@+] | "--"
- **org.opencontainers.image.title** Human-readable title of the image (string)
- **org.opencontainers.image.description** Human-readable description of the software packaged in the image (string)
- **org.opencontainers.image.base.digest** [Digest](https://github.com/opencontainers/image-spec/blob/main/descriptor.md#digests) of the image this image is based on (string)
  - This SHOULD be the immediate image sharing zero-indexed layers with the image, such as from a Dockerfile FROM statement.
  - This SHOULD NOT reference any other images used to generate the contents of the image (e.g., multi-stage Dockerfile builds).
- **org.opencontainers.image.base.name** Image reference of the image this image is based on (string)
  - This SHOULD be image references in the format defined by [distribution/distribution](https://github.com/distribution/distribution/blob/d0deff9cd6c2b8c82c6f3d1c713af51df099d07b/reference/reference.go).
  - This SHOULD be a fully qualified reference name, without any assumed default registry. (e.g., registry.example.com/my-org/my-image:tag instead of my-org/my-image:tag).
  - This SHOULD be the immediate image sharing zero-indexed layers with the image, such as from a Dockerfile FROM statement.
  - This SHOULD NOT reference any other images used to generate the contents of the image (e.g., multi-stage Dockerfile builds).
  - If the image.base.name annotation is specified, the image.base.digest annotation SHOULD be the digest of the manifest referenced by the image.ref.name annotation.
