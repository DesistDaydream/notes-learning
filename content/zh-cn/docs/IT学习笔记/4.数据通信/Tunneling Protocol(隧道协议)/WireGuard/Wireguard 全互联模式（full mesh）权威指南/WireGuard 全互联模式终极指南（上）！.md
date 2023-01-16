---
title: WireGuard 全互联模式终极指南（上）！
---

[WireGuard 全互联模式终极指南（上）！](https://mp.weixin.qq.com/s/s6eIoxaVNXVHRnWBnugylg)
<https://mp.weixin.qq.com/s/KrDJs3e6JjKgCADNigPUJA>

大家好，我是米开朗基杨。

关注我的读者应该都还记得我之前写过一篇 [👉WireGuard 全互联模式 (full mesh) 的配置指南](https://mp.weixin.qq.com/s?__biz=MzU1MzY4NzQ1OA==&mid=2247492833&idx=1&sn=642025bef0c6d400bc3f1cab9906a8e6&chksm=fbeda42ccc9a2d3a00711e3c79d0c2dc50935b5139e89032537daf2b8af544a148a740ff12c1&scene=21&cur_album_id=1612086810350829568#wechat_redirect)，限于当时还没有成熟的产品来帮助我们简化全互联模式的配置，所以我选择了使用可视化界面 [👉wg-gen-web](https://mp.weixin.qq.com/s?__biz=MzU1MzY4NzQ1OA==&mid=2247491998&idx=1&sn=840c87b4ecc2108d4a121aa26825ce65&chksm=fbeda153cc9a284516e177a6bdbfc90e57a4f253beb5f2d1abaa9bca54a388e1fc60a5b61b2c&scene=21&cur_album_id=1612086810350829568#wechat_redirect) 来达成目的。但 [👉wg-gen-web](https://mp.weixin.qq.com/s?__biz=MzU1MzY4NzQ1OA==&mid=2247491998&idx=1&sn=840c87b4ecc2108d4a121aa26825ce65&chksm=fbeda153cc9a284516e177a6bdbfc90e57a4f253beb5f2d1abaa9bca54a388e1fc60a5b61b2c&scene=21&cur_album_id=1612086810350829568#wechat_redirect) 的缺陷也很明显，它生成的每一个客户端的配置都要手动调整，终究还是不够便利。

今天我将为大家介绍一种更加完美的工具来配置 WireGuard 的全互联模式，这个工具就是 Netmaker\[1]。

**由于篇幅原因，本系列文章将会分成两篇进行介绍。本篇文章介绍 Netmaker 的工作原理和功能解读；下一篇文章将会介绍如何使用 Netmaker 来配置 WireGuard 全互联模式。**

## Netmaker 介绍

Netmaker 是一个用来配置 WireGuard 全互联模式的可视化工具，它的功能非常强大，不仅支持 UDP 打洞、NAT 穿透、多租户，还可以使用 Kubernetes 配置清单来部署，客户端几乎适配了所有平台，包括 Linux, Mac 和 Windows，还可以通过 WireGuard 原生客户端连接 iPhone 和 Android，真香！

其最新版本的基准测试结果显示，基于 Netmaker 的 WireGuard 网络速度比其他全互联模式的 VPN（例如 Tailscale 和 ZeroTier）网络速度快 50% 以上。

## Netmaker 架构

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/9125c963-31dc-45d8-ae30-c608056cb4ee/640)

Netmaker 使用的是 C/S 架构，即客户端 / 服务器架构。Netmaker Server 包含两个核心组件：用来管理网络的可视化界面，以及与客户端通信的 gRPC Server。你也可以可以选择部署 DNS 服务器（CoreDNS）来管理私有 DNS。

客户端（netclient）是一个二进制文件，可以在绝大多数 Linux 客户端以及 macOS 和 Windows 客户端运行，它的功能就是自动管理 WireGuard，动态更新 Peer 的配置。

> **注意**：这里不要将 Netmaker 理解成我之前的文章所提到的[👉 中心辐射型网络拓扑](https://mp.weixin.qq.com/s?__biz=MzU1MzY4NzQ1OA==&mid=2247491469&idx=1&sn=a5a1be7c9f4d2cf1f5f071a6f26f9856&chksm=fbee5f40cc99d6560cd756f6c710c066ae3605bad366b9b01006358542fd8374bb97ce50ac4f&scene=21&cur_album_id=1612086810350829568#wechat_redirect)。Netmaker Server 只是用来存储虚拟网络的配置并管理各个 Peer 的状态，Peer 之间的网络流量并不会通过 Netmaker Server。

Netmaker 还有一个重要的术语叫**签到**，客户端会通过定时任务来不断向 Netmaker Server 签到，以动态更新自身的状态和 Peer 的配置，它会从 Netmaker Server 检索 Peer 列表，然后与所有的 Peer 建立点对点连接，即全互联模式。所有的 Peer 通过互联最终呈现出来的网络拓扑结构就类似于本地子网或 VPC。

## Netmaker 部署

Netmaker 支持多种部署方式，包括二进制部署和容器化部署，容器化部署还支持 docker-compose 和 Kubernetes。如果你没有可以暴露服务到公网的 Kubernetes 集群，我推荐还是直接通过 docker-compose 来部署，简单粗暴。

官方推荐的做法是使用 Caddy 或 Nginx 来反向代理 Netmaker UI、API Server 和 gRPC Server，但我的域名没有备案，我也怕麻烦，就直接通过公网 IP 来提供服务了。如果你也想通过公网域名来暴露 Netmaker 的服务，可以参考 Netmaker 的官方文档\[2]。

本文的部署方案将直接通过公网 IP 来提供服务，首先需要安装 docker-compose，安装方法可以参考 Docker 官方文档\[3]。

安装完 docker-compose 后，再下载 docker-compose 的 YAML 配置清单：

`$ wget https://cdn.jsdelivr.net/gh/gravitl/netmaker@master/compose/docker-compose.yml`

现在还不能直接部署，需要根据自己的实际环境对配置清单进行修改。例如，我修改后的配置清单内容如下：

\`version: "3.4"

services:
  netmaker:
    container_name: netmaker
    image: gravitl/netmaker:v0.8.2
    volumes:
      - /etc/netclient/config:/etc/netclient/config
      - dnsconfig:/root/config/dnsconfig
      - /usr/bin/wg:/usr/bin/wg
      - /data/sqldata/:/root/data
    cap_add:
      - NET_ADMIN
    restart: always
    network_mode: host
    environment:
      SERVER_HOST: "\<public_ip>"
      COREDNS_ADDR: "\<public_ip>"
      GRPC_SSL: "off"
      DNS_MODE: "on"
      CLIENT_MODE: "on"
      API_PORT: "8081"
      GRPC_PORT: "50051"
      SERVER_GRPC_WIREGUARD: "off"
      CORS_ALLOWED_ORIGIN: "\*"
      DATABASE: "sqlite"
  netmaker-ui:
    container_name: netmaker-ui
    depends_on:
      - netmaker
    image: gravitl/netmaker-ui:v0.8
    links:
      - "netmaker:api"
    ports:
      - "80:80"
    environment:
      BACKEND_URL: "http://\<public_ip>:8081"
    restart: always
    network_mode: host
  coredns:
    depends_on:
      - netmaker
    image: coredns/coredns
    command: -conf /root/dnsconfig/Corefile
    container_name: coredns
    restart: always
    network_mode: host
    volumes:
      - dnsconfig:/root/dnsconfig
volumes:
  dnsconfig: {}

\`

总共有以下几处改动：

- 删除了不必要的环境变量，并修改了其中一部分环境变量，比如关闭 SSL 模式，将域名替换为公网 IP。你需要根据自己的实际环境将 `<public_ip>` 替换为你的公网 IP。
- 将所有容器的网络模式都改为 host 模式，即 `network_mode: host`。
- 将 sqlite 的数据存储改为 hostpath，即 `/data/sqldata/:/root/data`。

其中 `CLIENT_MODE: "on"` 表示将 Netmaker Server 所在的节点也作为 Mesh Network 的 Peer 节点。

最后我们就可以通过配置清单来部署容器了：

`$ docker-compose up -d`

查看是否部署成功：

\`$ docker ps

CONTAINER ID    IMAGE                                 COMMAND                   CREATED       STATUS    PORTS    NAMES
0daf3a35f8ce    docker.io/coredns/coredns:latest      "/coredns -conf /roo…"    7 days ago    Up                 coredns
0dbb0158e821    docker.io/gravitl/netmaker-ui:v0.8    "/docker-entrypoint.…"    7 days ago    Up                 netmaker-ui
bd39ee52013e    docker.io/gravitl/netmaker:v0.8.2     "./netmaker"              7 days ago    Up                 netmaker

\`

部署成功后，就可以在浏览器的地址栏输入你的公网 IP 来访问 Netmaker UI 了。

## Netmaker 功能解读

我们先通过 UI 来看看 Netmaker 都有哪些功能。

### 网络（Networks）

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/9125c963-31dc-45d8-ae30-c608056cb4ee/640)

Netmaker 允许创建任意数量的私有网络，可以设置任意地址范围。你只需要给这个网络起个名字，设置一个地址范围，并选择想要启用的选项。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/9125c963-31dc-45d8-ae30-c608056cb4ee/640)

目前总共包含三个可选项：

- **Dual Stack** : 双栈，即开启 IPv6。
- **Local Only** : 各个 Peer 之间只会通过内网地址来互联，即 Endpoint 皆为内网地址。适用于数据中心、VPC 或家庭 / 办公网络的内部。
- **Hole Punching** : 动态发现和配置 Endpoint 和端口，帮助 Peer 轻松穿透 NAT 进行 UDP 打洞。

管理员拥有对网络的最高控制器，例如，更改私有网络的网段，Peer 便会自动更新自身的 IP。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/9125c963-31dc-45d8-ae30-c608056cb4ee/640)

如果发现网络被入侵，也可以让网络中的所有节点刷新公钥。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/9125c963-31dc-45d8-ae30-c608056cb4ee/640)

### 节点（Nodes）

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/9125c963-31dc-45d8-ae30-c608056cb4ee/640)

Node 表示节点，通常是运行 Linux 的服务器，安装了 netclient 和 WireGuard。这个节点会通过 WireGuard 私有网络和其他所有节点相连。一但节点被添加到私有网络中，Netmaker 管理员就可以操控该节点的配置，例如：

- 私有网络地址
- 过期时间
- WireGuard 相关设置

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/9125c963-31dc-45d8-ae30-c608056cb4ee/640)

管理员也可以将该节点从私有网络中完全删除，让其无法连接其他所有 Peer 节点。

Node 还有两个比较重要的功能，就是将自身设置为 Ingress Gateway（入口网关）或者 Egress Gateway（出口网关）。Ingress Gateway 允许外部客户端的流量进入内部网络，Egress Gateway 允许将内部网络的流量转发到外部指定的 IP 范围。这两项功能对全互联模式进行了扩展，比如手机客户端就可以通过 Ingress Gateway 接入进来。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/9125c963-31dc-45d8-ae30-c608056cb4ee/640)

### 访问秘钥（Access Keys）

一个节点想要加入到私有网络，需要获取访问秘钥进行授权，当然你也可以选择手动批准。

一个访问秘钥可以被多个节点重复使用，你只需修改 Number 数量就可以实现这个目的。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/9125c963-31dc-45d8-ae30-c608056cb4ee/640)

访问秘钥创建后只会显示一次，展示了三个选项：

1. 原始访问秘钥
2. 访问令牌（access token），它将访问密钥与用于加入网络的参数（例如地址、端口和网络名称）包装在一起。当你运行 `netclient join -t <token>` 时，netclient 会对该令牌进行解码，并解析参数。
3. 安装脚本，用于在标准 Linux 服务器上首次安装 netclient。它只是简单地下载 netclient 并为你运行 "join" 命令。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/9125c963-31dc-45d8-ae30-c608056cb4ee/640)

### DNS

如果启用了 DNS 组件，Netmaker 就会通过 CoreDNS 来维护私有 DNS，它会为私有网络中的每个节点创建一个默认的 DNS 条目。你也可以创建自定义的 DNS 条目。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/9125c963-31dc-45d8-ae30-c608056cb4ee/640)

### 外部客户端（External Clients）

Netclient 目前只支持 Linux、macOS 和 Windows，如果 Android 和 iOS 端想要加入 VPN 私有网络，只能通过 WireGuard 原生客户端来进行连接。要想做到这一点，需要管理员事先创建一个 External Client，它会生成一个 WireGuard 配置文件，WireGuard 客户端可以下载该配置文件或者扫描二维码进行连接。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/9125c963-31dc-45d8-ae30-c608056cb4ee/640)

当然，在创建 External Client 之前，需要先设置其中一个节点为 Ingress Gateway。

需要说明的是，目前移动设备通过 External Client 接入只是权宜之计，随着 Netclient 对更多操作系统的支持，最终所有的客户端都应该使用 netclient 来连接。

## Netclient 介绍

netclient 是一个非常简单的 CLI，用于创建 WireGuard 配置和接口，将节点加入到 Netmaker 的私有网络中。netclient 可以管理任意数量的 Netmaker 私有网络，所有的网络都由同一个 netclient 实例管理。

\`$ netclient --help
NAME:
   Netclient CLI - Netmaker's netclient agent and CLI. Used to perform interactions with Netmaker server and set local WireGuard config.

USAGE:
   netclient \[global options] command \[command options]\[arguments...]

VERSION:
   v0.8.1

COMMANDS:
   join       Join a Netmaker network.
   leave      Leave a Netmaker network.
   checkin    Checks for local changes and then checks into the specified Netmaker network to ask about remote changes.
   push       Push configuration changes to server.
   pull       Pull latest configuration and peers from server.
   list       Get list of networks.
   uninstall  Uninstall the netclient system service.
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)

\`

### Netclient 工作原理

使用 netclient 可以加入某个网络，拉取或推送变更，以及离开某个网络。同时 netclient 还有几个辅助命令用于其他场景。

使用 netclient 加入某个网络时，它会创建一个目录 `/etc/netclient`，并将 netclient 二进制文件本身复制到该目录下。

`$ ls -lh /etc/netclient/netclient -rwxr-xr-x 1 root root 12M Oct  8 23:08 /etc/netclient/netclient`

同时会在该目录下创建一个子目录 `config`，并在子目录下创建相应的配置文件。比如你加入的网络名称是 default，那么配置文件名称就是 `netconfig-default`。

`$ ls -lh /etc/netclient/config/ total 32K -rwxr-xr-x 1 root root 1.8K Oct 17 16:23 netconfig-default -rw-r--r-- 1 root root  176 Oct  8 23:08 nettoken-default -rw-r--r-- 1 root root   16 Oct  8 23:08 secret-default -rw-r--r-- 1 root root   44 Oct  8 23:08 wgkey-default`

如果第一次使用 netclient 加入某个网络，它会尝试将自己设置为当前节点的守护进程，以 Linux 为例，它会创建一个 systemd 服务：

\`$ cat /etc/systemd/system/netclient.service
\[Unit]
Description=Network Check
Wants=netclient.timer

\[Service]
Type=simple
ExecStart=/etc/netclient/netclient checkin -n all

\[Install]
WantedBy=multi-user.target

\`

该 systemd 服务的作用是向 Netmaker Server **签到**，并将本地的配置与 Netmaker Server 托管的配置进行比较，根据比较结果进行适当修改，再拉取所有的 Peer 列表，最后重新配置 WireGuard。

同时还会设置一个计划任务，来定期（每 15 秒执行一次）启动守护进程同步本地和远程 Netmaker Server 的配置。

\`$ cat /etc/systemd/system/netclient.timer
\[Unit]
Description=Calls the Netmaker Mesh Client Service
Requires=netclient.service

\[Timer]
Unit=netclient.service

OnCalendar=_:_:0/15

\[Install]
WantedBy=timers.target

\`

对于不支持 systemd 的 Linux 发行版，我们可以采取其他方式来执行守护进程和计划任务。我们也可以把 netclient 作为调试工具，执行 `netclient pull` 从 Netmaker Server 获取最新配置，执行 `netclient push` 将本地变更推送到 Netmaker Server，等等。

## 总结

本文在讲解过程中略过了很多功能和选项的细节，如果你有兴趣了解某个特定的功能或者选项，可以查阅 Netmaker 的官方文档\[4]。下一篇文章将会介绍如何使用 Netmaker 来配置 WireGuard 全互联模式，我会详细介绍 Linux、macOS 和手机客户端分别该如何配置，敬请期待！

### 引用链接

\[1]

Netmaker: [_https://github.com/gravitl/netmaker_](https://github.com/gravitl/netmaker)

\[2]

Netmaker 的官方文档: [_https://docs.netmaker.org/quick-start.html_](https://docs.netmaker.org/quick-start.html)

\[3]

Docker 官方文档: [_https://docs.docker.com/compose/install/_](https://docs.docker.com/compose/install/)

\[4]

Netmaker 的官方文档: [_https://docs.netmaker.org/_](https://docs.netmaker.org/)
