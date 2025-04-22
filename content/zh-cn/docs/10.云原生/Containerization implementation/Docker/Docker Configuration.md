---
title: Docker Configuration
linkTitle: Docker Configuration
weight: 3
---

# 概述

> 参考：
>
> - [官方文档，参考 - 命令行参考 - Daemon CLI](https://docs.docker.com/engine/reference/commandline/dockerd/)
> - https://blog.csdn.net/u013948858/article/details/79974796

Docker 的守护进程为 dockerd，dockerd 可以通过两种方式配置运行时行为

- /etc/docker/daemon.json 配置文件
- dockerd 命令的 flags ，可以将 flags 添加到 dockerd.service 中。

> [!Note]
> 配置文件中的配置，也可以通过 dockerd 的命令行参数(也就是 flags)指定，比如配置文件中的 data-root 字段，对应的 dockerd flags 为 --data-root

## 配置文件示例

dockerd 配置文件是 JSON 格式，基本常用的配置内容如下。

```json
{
  // 指定 docker pull 时，首先去连接的 registry。
  "registry-mirrors": [
    "http://172.38.40.180",
    "https://ac1rmo5p.mirror.aliyuncs.com"
  ],
  // 指定运行 docker 操作的不安全的 registry 的列表
  "insecure-registries": ["http://172.38.40.180"],
  // 指定 docker 运行时其他的选项，这里面指定 docker 的 cgroupdriver 为 systemd
  "exec-opts": ["native.cgroupdriver=systemd"],
  // 指定 docker 的日志驱动为 json-file
  "log-driver": "json-file",
  // 指定 docker 记录容器日志的参数，这里指定容器日志文件大小最大为100m
  "log-opts": {
    "max-size": "100m"
  },
  // 指定 docker 的存储驱动类型为 overlay2
  "storage-driver": "overlay2",
  // 指定 docker 存储驱动的其他选项
  "storage-opts": ["overlay2.override_kernel_check=true"]
}
```

# 命令行标志详解

绝大部分的命令行标志都可以参考[配置文件详解](#配置文件详解)进行使用

**-H, --host \<LIST>** # 要连接的守护进程 socket。官方给的 systemd 的 service 中，使用的是 `docker -H fd://`，这样的用法必须要先启动 docker.socket

# 配置文件详解

[官方文档，参考 - CLI 参考 - dockerd - daemon 配置文件](https://docs.docker.com/reference/cli/dockerd/#daemon-configuration-file) 是官网给的配置文件中所有可用字段的列表

**data-root**(STRING) # 配置 docker info 命令中的 Docker Root Dir，也就是 docker 存储数据的路径。

**features**(OBJECT) # 一些新的特性可以通过配置该字段来启动或停止

**hosts**(\[]STRING) # 指定 docker 守护进程监听的端口

- 可以从其他机器使用 `docker -H URL` 命令对该设备进行 docker 操作

**live-restore**(BOOLEAN) # 在 docker.service 守护程序停止期间，保持容器状态，说白了就是重启 docker 的时候 Containers 不重启。

- 开启该参数后，就算重启 dockerd 服务也不会更改 default-address-pools 参数执行的地址范围

**registry-mirrors**(\[]STRING) # 指定 pull、push 镜像时候的加速器地址

- 可用的地址参考 [容器镜像管理](/docs/10.云原生/Containerization%20implementation/容器管理/容器镜像管理.md)

**insecure-registries**(\[]STRING) # 指定不安全仓库，Docker 默认无法连接 HTTP 协议的仓库，将仓库的 URL 添加到该字段后，docker 即可连接

## 日志配置

**log-driver(STRING)** # 指定 docker 的日志驱动

**log-opts(OBJECT)** # 指定 docker 记录容器日志的参数

## 存储驱动配置

**storage-driver(STRING)** # 在 Linux 上，Dockerd 进程支持多种不同的镜像层存储驱动程序: `overlay2`, `fuse-overlayfs`, `btrfs`, `zfs`, and `devicemapper`.

**storage-opts(\[]STRING)** # 与存储驱动相关的选项。

- 很多选项随着版本的更新，会弃用，比如 `overlay2.override_kernel_check` 已于 [24.0+ 版本弃用](https://docs.docker.com/engine/deprecated/#support-for-the-overlay2override_kernel_check-storage-option)

## 网络配置

**ip-forward-no-drop**(BOOLEAN) # 禁止 Docker 将 iptables 的 filter 表中 FORWARD 链的默认行为改为 DROP。`默认值: false`

**bip** # 指定 docker0 桥的 IP

# 代理配置

从 Docker 的 23.0 版本开始，可以在 daemon.json 文件中配置 dockerd 的代理行为：

```json
{
  "proxies": {
    "http-proxy": "http://192.168.254.254:7890",
    "https-proxy": "http://192.168.254.254:7890",
    "no-proxy": "*.test.example.com,.example.org,127.0.0.0/8"
  }
}
```

如果 docker 所在的环境是通过代理服务器和互联网连通的，那么需要一番配置才能让 docker 正常从外网正常拉取镜像。然而仅仅通过配置环境变量的方法是不够的。本文结合已有文档，介绍如何配置代理服务器能使docker正常拉取镜像。

本文使用的 docker 版本是19.03

下面的文章权当记录留个念想了。

## 问题现象

> 参考：
>
> - [如何配置docker通过代理服务器拉取镜像](https://www.lfhacks.com/tech/pull-docker-images-behind-proxy/)

如果不配置代理服务器就直接拉镜像，docker 会直接尝试连接镜像仓库，并且连接超时报错。如下所示：

```shell
$ docker pull busybox
Using default tag: latest
Error response from daemon: Get https://registry-1.docker.io/v2/: net/http: request canceled
while waiting for connection (Client.Timeout exceeded while awaiting headers)
```

## 容易误导的官方文档

有这么一篇关于 docker 配置代理服务器的 [官方文档](https://docs.docker.com/network/proxy/#configure-the-docker-client) （[新链接](https://docs.docker.com/engine/cli/proxy/#configure-the-docker-client)），如果病急乱投医，直接按照这篇文章配置，是不能成功拉取镜像的。

我们来理解一下这篇文档，文档关键的原文摘录如下：

> If your container needs to use an HTTP, HTTPS, or FTP proxy server, you can configure it in different ways: Configure the Docker client On the Docker client, create or edit the file ~/.docker/config.json in the home directory of the user that starts containers.
>
> When you create or start new containers, the environment variables are set automatically within the container.

这篇文档说：如果你的容器或者使用 `docker build` 构建镜像时需要使用代理服务器，那么可以以如下方式配置： 在运行容器的用户 home 目录下，配置 `~/.docker/config.json` 文件。重新启动容器后，这些环境变量将自动设置进容器，从而容器内的进程可以使用代理服务。

所以这篇文章是讲如何配置运行容器的环境，与如何拉取镜像无关。如果按照这篇文档的指导，如同南辕北辙。

要解决问题，我们首先来看一般情况下命令行如何使用代理。

## 环境变量

常规的命令行程序如果要使用代理，需要设置两个环境变量：`HTTP_PROXY` 和 `HTTPS_PROXY` ，设置环境变量的方法见 [这篇文章](https://www.lfhacks.com/test/cypress-download-failure#env) 。但是仅仅这样设置环境变量，也不能让 docker 成功拉取镜像。

我们仔细观察 [上面的报错信息](https://www.lfhacks.com/tech/pull-docker-images-behind-proxy/#problem)，有一句说明了报错的来源：

> Error response from daemon:

因为镜像的拉取和管理都是 docker daemon 的职责，所以我们要让 docker daemon 知道代理服务器的存在。而 docker daemon 是由 systemd 管理的，所以我们要从 systemd 配置入手。

## 正确的官方文档

关于 systemd 配置代理服务器的 [官方文档在这里](https://docs.docker.com/config/daemon/systemd/#httphttps-proxy)（[新链接](https://docs.docker.com/engine/daemon/proxy/#httphttps-proxy)），原文说：

> The Docker daemon uses the HTTP_PROXY, HTTPS_PROXY, and NO_PROXY environmental variables in its start-up environment to configure HTTP or HTTPS proxy behavior. You cannot configure these environment variables using the daemon.json file.
>
> This example overrides the default docker.service file.
>
> If you are behind an HTTP or HTTPS proxy server, for example in corporate settings, you need to add this configuration in the Docker systemd service file.

这段话的意思是，docker daemon 使用 `HTTP_PROXY`, `HTTPS_PROXY`, 和 `NO_PROXY` 三个环境变量配置代理服务器，但是你需要在 systemd 的文件里配置环境变量，而不能配置在 `daemon.json` 里。

## 具体操作

下面是来自官方文档的操作步骤和详细解释：

1. 创建 dockerd 相关的 systemd 目录，这个目录下的配置将覆盖 dockerd 的默认配置

```shell
sudo mkdir -p /etc/systemd/system/docker.service.d
```

2. 新建配置文件 `/etc/systemd/system/docker.service.d/http-proxy.conf`，这个文件中将包含环境变量

```ini
[Service]
Environment="HTTP_PROXY=http://192.168.254.254:7890"
Environment="HTTPS_PROXY=http://192.168.254.254:7890"
```

3. 如果你自己建了私有的镜像仓库，需要 dockerd 绕过代理服务器直连，那么配置 `NO_PROXY` 变量：

```ini
[Service]
Environment="HTTP_PROXY=http://192.168.254.254:7890"
Environment="HTTPS_PROXY=http://192.168.254.254:7890"
Environment="NO_PROXY=your-registry.com,10.10.10.10,*.example.com"
```

多个 `NO_PROXY` 变量的值用逗号分隔，而且可以使用通配符（*），极端情况下，如果 `NO_PROXY=*`，那么所有请求都将不通过代理服务器。

4. 重新加载配置文件，重启 dockerd

```shell
sudo systemctl daemon-reload
sudo systemctl restart docker
```

5. 检查确认环境变量已经正确配置：

```shell
sudo systemctl show --property=Environment docker
```

这样配置后，应该可以正常拉取 docker 镜像。

## 结论

docker 镜像由 docker daemon 管理，所以不能用修改 shell 环境变量的方法使用代理服务，而是从 systemd 角度设置环境变量。
