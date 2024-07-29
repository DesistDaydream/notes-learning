---
title: Tailscale DERP
linkTitle: Tailscale DERP
date: 2024-03-21T23:14
weight: 20
---

# 概述

> 参考：
> 
> - [公众号，自建 DERP 中继服务器，从此 Tailscale 畅通无阻](https://mp.weixin.qq.com/s/r5OQi9YreI-BFnClKhBR0w)
> - [官方文档，自定义 DERP 服务器](https://tailscale.com/kb/1118/custom-derp-servers)

[👉 上篇文章](https://mp.weixin.qq.com/s?__biz=MzU1MzY4NzQ1OA==&mid=2247504037&idx=1&sn=b059e0ed24be4ae39a25e5724700ff54&scene=21#wechat_redirect)介绍了如何使用 `Headscale` 替代 Tailscale 官方的控制服务器，并接入各个平台的客户端。本文将会介绍如何让 Tailscale 使用自定义的 DERP Servers。可能很多人都不知道 DERP 是个啥玩意儿，没关系，我先从 **中继服务器** 开始讲起。

## STUN 是什么

Tailscale 的终极目标是让两台**处于网络上的任何位置**的机器建立**点对点连接**（直连），但现实世界是复杂的，大部份情况下机器都位于 NAT 和防火墙后面，这时候就需要通过打洞来实现直连，也就是 NAT 穿透。

NAT 按照 **NAT 映射行为**和**有状态防火墙行为**可以分为多种类型，但对于 NAT 穿透来说根本不需要关心这么多类型，只需要看 **NAT 或者有状态防火墙是否会严格检查目标 Endpoint**，根据这个因素，可以将 NAT 分为  **Easy NAT** 和 **Hard NAT**。

- **Easy NAT** 及其变种称为 “Endpoint-Independent Mapping” (**EIM，终点无关的映射**) 这里的 Endpoint 指的是目标 Endpoint，也就是说，有状态防火墙只要看到有客户端自己发起的出向包，就会允许相应的入向包进入，**不管这个入向包是谁发进来的都可以**。
- **hard NAT** 以及变种称为 “Endpoint-Dependent Mapping”（**EDM，终点相关的映射**） 这种 NAT 会针对每个目标 Endpoint 来生成一条相应的映射关系。在这样的设备上，如果客户端向某个目标 Endpoint 发起了出向包，假设客户端的公网 IP 是 2.2.2.2，那么有状态防火墙就会打开一个端口，假设是 4242。那么只有来自该目标 Endpoint 的入向包才允许通过 `2.2.2.2:4242`，其他客户端一律不允许。这种 NAT 更加严格，所以叫 Hard NAT。

对于 Easy NAT，我们只需要提供一个第三方的服务，它能够告诉客户端“它看到的客户端的公网 ip:port 是什么”，然后将这个信息以某种方式告诉通信对端（peer），后者就知道该和哪个地址建连了！这种服务就叫 **STUN** (Session Traversal Utilities for NAT，NAT 会话穿越应用程序)。它的工作流程如下图所示：

- 笔记本向 STUN 服务器发送一个请求：“从你的角度看，我的地址什么？”
- STUN 服务器返回一个响应：“我看到你的 UDP 包是从这个地址来的：`ip:port`”。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pq65ei/1648866197865-44dd313f-e0cb-4108-b4ca-52e285e1dac5.jpeg)

## 中继是什么

对于 **Hard NAT** 来说，STUN 就不好使了，即使 STUN 拿到了客户端的公网 `ip:port` 告诉通信对端也于事无补，因为防火墙是和 STUN 通信才打开的缺口，这个缺口只允许 STUN 的入向包进入，其他通信对端知道了这个缺口也进不来。通常企业级 NAT 都属于 Hard NAT。

这种情况下打洞是不可能了，但也不能就此放弃，可以选择一种折衷的方式：创建一个中继服务器（relay server），客户端与中继服务器进行通信，中继服务器再将包中继（relay）给通信对端。

至于中继的性能，那要看具体情况了：

- 如果能直连，那显然没必要用中继方式；
- 但如果无法直连，而中继路径又非常接近双方直连的真实路径，并且带宽足够大，那中继方式并不会明显降低通信质量。延迟肯定会增加一点，带宽会占用一些，但**相比完全连接不上，还是可以接受的**。

事实上对于大部分网络而言，Tailscale 都可以通过各种黑科技打洞成功，只有极少数情况下才会选择中继，中继只是一种 fallback 机制。

## 中继协议简介

中继协议有多种实现方式。

### TURN

TURN 即 Traversal Using Relays around NAT，这是一种经典的中继实现方式，核心理念是：

- **用户**（人）先去公网上的 TURN 服务器认证，成功后后者会告诉你：“我已经为你分配了 ip:port，接下来将为你中继流量”，
- 然后将这个 ip:port 地址告诉对方，让它去连接这个地址，接下去就是非常简单的客户端/服务器通信模型了。

与 STUN 不同，这种协议没有真正的交互性，不是很好用，因此 Tailscale 并没有采用 TURN 作为中继协议。

### DERP

DERP 即 Detoured Encrypted Routing Protocol，这是 Tailscale 自研的一个协议：

- 它是一个**通用目的包中继协议，运行在 HTTP 之上**，而大部分网络都是允许 HTTP 通信的。
- 它根据目的公钥（destination’s public key）来中继加密的流量（encrypted payloads）。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pq65ei/1648866197927-9bf9875f-5ee8-474f-ab5a-d6c8448b2697.svg)

Tailscale 会自动选择离目标节点最近的 DERP server 来中继流量

Tailscale 使用的算法很有趣，**所有客户端之间的连接都是先选择 DERP 模式（中继模式），这意味着连接立即就能建立（优先级最低但 100% 能成功的模式），用户不用任何等待**。然后开始并行地进行路径发现，通常几秒钟之后，我们就能发现一条更优路径，然后将现有连接透明升级（upgrade）过去，变成点对点连接（直连）。

因此，DERP 既是 Tailscale 在 NAT 穿透失败时的保底通信方式（此时的角色与 TURN 类似），也是在其他一些场景下帮助我们完成 NAT 穿透的旁路信道。换句话说，它既是我们的保底方式，也是有更好的穿透链路时，帮助我们进行连接升级（upgrade to a peer-to-peer connection）的基础设施。

Tailscale 的私钥只会保存在当前节点，因此 DERP server 无法解密流量，它只能和互联网上的其他路由器一样，呆呆地将加密的流量从一个节点转发到另一个节点，只不过 DERP 使用了一个稍微高级一点的协议来防止滥用。

# 私有 DERP Server

> 参考：
>
> - [GitHub 项目，tailscale/tailscale - cmd/derper](https://github.com/tailscale/tailscale/tree/main/cmd/derper) derper 入口
>   - https://github.com/tailscale/tailscale/tree/main/derp DERP 的实现

Tailscale 官方[内置了很多 DERP 服务器](https://controlplane.tailscale.com/derpmap/default)，分步在全球各地，**惟独不包含中国大陆**，原因你懂得。这就导致了一旦流量通过 DERP 服务器进行中继，延时就会非常高。而且官方提供的 DERP 服务器是万人骑，存在安全隐患。

Tailscale 开源了 DERP Server 的代码，将其称为 **DERPer**。我们可以自己部署 derper 以便让所有流量都通过自己的程序。

> [!Tip] [Headscale](docs/4.数据通信/通信协议/Tunneling%20Protocol/Tailscale/Headscale.md) 包含了内置的 DERPer

## 部署 derper

为了实现低延迟、高安全性，可以参考 [Tailscale 官方文档](https://tailscale.com/kb/1118/custom-derp-servers)自建私有的 DERP 服务器。有两种部署模式，一种是基于域名，另外一种不需要域名，可以直接使用 IP，不过需要一点黑科技。

### 使用域名

这种方案需要满足以下几个条件：

- 要有自己的域名，并且申请了 SSL 证书
- 需要准备一台或多台云主机
- 如果服务器在国内，域名需要备案
- 如果服务器在国外，则不需要备案

如果以上条件都俱备，就可以按照下面的步骤开始部署了。

推荐直接使用 Docker 来部署，我已经构建好了 Docker 镜像，直接部署就可以了：

```
docker run --restart always \  
  --name derper -p 12345:12345 -p 3478:3478/udp \  
  -v /root/.acme.sh/xxxx/:/app/certs \  
  -e DERP_CERT_MODE=manual \  
  -e DERP_ADDR=12345 \  
  -e DERP_DOMAIN=xxxx \  
  -d ghcr.io/yangchuansheng/derper:latest
```

有几点需要注意：

- 能用 443 端口尽量用 443 端口，实在不行再用别的端口；
- 默认情况下也会开启 STUN 服务，UDP 端口是 `3478`；
- 防火墙需要放行端口 12345 和 3478；
- 准备好 SSL 证书；
- 域名部分我打了码，请换成你自己的域名。

关于证书部分需要重点说明：**假设你的域名是 `xxx.com`**，那么证书的名称必须是 **`xxx.com.crt`**，一个字符都不能错！同理，私钥名称必须是 **`xxx.com.key`**，一个字符都不能错！

查看容器日志：

```bash
🐳  → docker logs -f derper  
2022/03/26 11:36:28 no config path specified; using /var/lib/derper/derper.key  
2022/03/26 11:36:28 derper: serving on :12345 with TLS  
2022/03/26 11:36:28 running STUN server on [::]:3478
```

目前 derper 运行一段时间就会崩溃，暂时还没有更好的解决方案，只能通过定时重启来解决，比如通过 crontab 来设置每两小时重启一次容器：

```
0 */2 * * * docker restart derper &> /dev/null
```

具体可参考这个 issue：[Derper TLS handshake error: remote error: tls: internal error](https://github.com/tailscale/tailscale/issues/4082)

---

部署好 derper 之后，就可以修改 Headscale 的配置来使用自定义的 DERP 服务器了。Headscale 可以通过两种形式的配置来使用自定义 DERP：

- 一种是在线 URL，格式是 `JSON`，与 Tailscale 官方控制服务器使用的格式和语法相同。
- 另一种是本地文件，格式是 `YAML`。

我们可以直接使用本地的 YAML 配置文件，内容如下：

```yaml
# /etc/headscale/derp.yaml  
regions:  
  900:  
    regionid: 900  
    regioncode: thk   
    regionname: Tencent Hongkong   
    nodes:  
      - name: 900a  
        regionid: 900  
        hostname: xxxx  
        ipv4: xxxx  
        stunport: 3478  
        stunonly: false  
        derpport: 12345  
      - name: 900b  
        regionid: 900  
        hostname: xxxx  
        ipv4: xxxx  
        stunport: 3478  
        stunonly: false  
        derpport: 12345  
  901:  
    regionid: 901  
    regioncode: hs   
    regionname: Huawei Shanghai   
    nodes:  
      - name: 901a  
        regionid: 901  
        hostname: xxxx  
        ipv4: xxxx  
        stunport: 3478  
        stunonly: false  
        derpport: 12345
```

配置说明：

- `regions` 是 YAML 中的**对象**，下面的每一个对象表示一个**可用区**，每个**可用区**里面可设置多个 DERP 节点，即 `nodes`。
- 每个可用区的 `regionid` 不能重复。
- 每个 `node` 的 `name` 不能重复。
- `regionname` 一般用来描述可用区，`regioncode` 一般设置成可用区的缩写。
- `ipv4` 字段不是必须的，如果你的域名可以通过公网解析到你的 DERP 服务器地址，这里可以不填。如果你使用了一个二级域名，而这个域名你并没有在公共 DNS server 中添加相关的解析记录，那么这里就需要指定 IP（前提是你的证书包含了这个二级域名，这个很好支持，搞个泛域名证书就行了）。
- `stunonly: false` 表示除了使用 STUN 服务，还可以使用 DERP 服务。
- 上面的配置中域名和 IP 部分我都打码了，你需要根据你的实际情况填写。

接下来还需要修改 Headscale 的配置文件，引用上面的自定义 DERP 配置文件。需要修改的配置项如下：

```yaml
# /etc/headscale/config.yaml  
derp:  
  # List of externally available DERP maps encoded in JSON  
  urls:  
  #  - https://controlplane.tailscale.com/derpmap/default  
  
  # Locally available DERP map files encoded in YAML  
  #  
  # This option is mostly interesting for people hosting  
  # their own DERP servers:  
  # https://tailscale.com/kb/1118/custom-derp-servers/  
  #  
  # paths:  
  #   - /etc/headscale/derp-example.yaml  
  paths:  
    - /etc/headscale/derp.yaml  
  
  # If enabled, a worker will be set up to periodically  
  # refresh the given sources and update the derpmap  
  # will be set up.  
  auto_update_enabled: true  
  
  # How often should we check for DERP updates?  
  update_frequency: 24h
```

可以把 Tailscale 官方的 DERP 服务器禁用，来测试自建的 DERP 服务器是否能正常工作。

修改完配置后，重启 headscale 服务：

```
$ systemctl restart headscale
```

在 Tailscale 客户端上使用以下命令查看目前可以使用的 DERP 服务器：

```bash
$ tailscale netcheck  
  
Report:  
        * UDP: true  
        * IPv4: yes, xxxxx:57068  
        * IPv6: no  
        * MappingVariesByDestIP: false  
        * HairPinning: false  
        * PortMapping:   
        * Nearest DERP: Tencent Hongkong  
        * DERP latency:  
                - thk: 39.7ms (Tencent Hongkong)
```

`tailscale netcheck` 实际上只检测 `3478/udp` 的端口， 就算 netcheck 显示能连，也不一定代表 12345 端口可以转发流量。最简单的办法是直接打开 DERP 服务器的 URL：https://xxxx:12345，如果看到如下页面，且地址栏的 SSL 证书标签显示正常可用，那才是真没问题了。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pq65ei/1648866197919-6ec2816f-e8f4-4be6-993d-c14fe47aa72a.png)

查看与通信对端的连接方式：

```bash
$ tailscale status  
10.1.0.5        coredns              default      linux   -  
                carsondemacbook-pro  default      macOS   active; direct xxxx:2756; offline, tx 50424 rx 34056  
                oneplus-8t           default      android active; relay "thk"; offline, tx 1608 rx 1552  
                openwrt              default      linux   active; direct xxxx:2834; offline, tx 1403688 rx 1217620
```

这个客户端是一台云主机，有 3 个通信对端，分别是 macOS、OpenWRT 与 Android 手机，macOS 和 OpenWRT 都处于电信家庭内网中，Android 手机使用的是电信流量。可以看到只有 Android 手机是通过自定义的 DERP 服务器来中继流量的，打洞成功率相当高。使用 ping 来测试连通性：

```bash
$ ping 10.1.0.8  
PING 10.1.0.8 (10.1.0.8) 56(84) bytes of data.  
64 bytes from 10.1.0.8: icmp_seq=1 ttl=64 time=150 ms  
64 bytes from 10.1.0.8: icmp_seq=2 ttl=64 time=131 ms  
64 bytes from 10.1.0.8: icmp_seq=3 ttl=64 time=161 ms  
64 bytes from 10.1.0.8: icmp_seq=4 ttl=64 time=137 ms  
64 bytes from 10.1.0.8: icmp_seq=5 ttl=64 time=156 ms  
64 bytes from 10.1.0.8: icmp_seq=6 ttl=64 time=169 ms  
^C  
--- 10.1.0.8 ping statistics ---  
6 packets transmitted, 6 received, 0% packet loss, time 5005ms  
rtt min/avg/max/mdev = 131.728/151.154/169.627/13.193 ms
```

也可以使用 Tailscale 命令行工具来测试：

```bash
$ tailscale ping 10.1.0.8  
pong from oneplus-8t (10.1.0.8) via DERP(thk) in 104ms  
pong from oneplus-8t (10.1.0.8) via DERP(thk) in 111ms  
pong from oneplus-8t (10.1.0.8) via DERP(thk) in 105ms
```

这个更加友好一点，会直接告诉你是通过 DERP 中继服务器来和对方通信的。

如果当前 Tailscale 客户端所在主机开启了 IPv6，那么与手机便可以直接通过 IPv6 点对点连接：

```bash
$ /Applications/Tailscale.app/Contents/MacOS/Tailscale status  
                coredns              default      linux   active; direct xxxx:45986; offline, tx 124352 rx 185736  
                oneplus-8t           default      android active; direct [240e:472:da0:24a2:a07f:2a67:2a1e:4475]:37237; offline, tx 125216 rx 20052  
                openwrt              default      linux   active; direct [240e:390:caf:1870:c02c:e8ff:feb9:b0b]:41641; offline, tx 181992 rx 3910120  
  
$ /Applications/Tailscale.app/Contents/MacOS/Tailscale ping 10.1.0.8  
pong from oneplus-8t (10.1.0.8) via [240e:472:da0:24a2:a07f:2a67:2a1e:4475]:37237 in 62ms
```

所以如果你开启了 IPv6，可以大大增加**点对点连接**的成功率。

### 使用纯 IP

我知道，大部分人是没有自己的域名的。再退一步，就算有自己的域名，如果没有备案，也是没办法部署在国内服务器上使用的。

这个时候我们就只能从 derper 源码上动手脚了，找到 tailscale 仓库中的 `cmd/derper/cert.go` 文件，将与域名验证相关的内容删除或注释：

```go
func (m *manualCertManager) getCertificate(hi *tls.ClientHelloInfo) (*tls.Certificate, error) {  
 //if hi.ServerName != m.hostname {  
 // return nil, fmt.Errorf("cert mismatch with hostname: %q", hi.ServerName)  
 //}  
 return m.cert, nil  
}
```

还需要创建自签名证书，可以通过脚本来创建：

```bash
# build_cert.sh  
  
#!/bin/bash  
  
CERT_HOST=$1  
CERT_DIR=$2  
CONF_FILE=$3  
  
echo "[req]  
default_bits  = 2048  
distinguished_name = req_distinguished_name  
req_extensions = req_ext  
x509_extensions = v3_req  
prompt = no  
  
[req_distinguished_name]  
countryName = XX  
stateOrProvinceName = N/A  
localityName = N/A  
organizationName = Self-signed certificate  
commonName = $CERT_HOST: Self-signed certificate  
  
[req_ext]  
subjectAltName = @alt_names  
  
[v3_req]  
subjectAltName = @alt_names  
  
[alt_names]  
IP.1 = $CERT_HOST  
" > "$CONF_FILE"  
  
mkdir -p "$CERT_DIR"  
openssl req -x509 -nodes -days 730 -newkey rsa:2048 -keyout "$CERT_DIR/$CERT_HOST.key" -out "$CERT_DIR/$CERT_HOST.crt" -config "$CONF_FILE"
```

重新编写 Dockerfile，将 derper 的域名设置为 `127.0.0.1`：

```dockerfile
FROM golang:latest AS builder  
  
WORKDIR /app  
  
# ========= CONFIG =========  
# - download links  
ENV MODIFIED_DERPER_GIT=https://github.com/yangchuansheng/ip_derper.git  
ENV BRANCH=ip_derper  
# ==========================  
  
# build modified derper  
RUN git clone -b $BRANCH $MODIFIED_DERPER_GIT tailscale --depth 1 && \  
    cd /app/tailscale/cmd/derper && \  
    /usr/local/go/bin/go build -ldflags "-s -w" -o /app/derper && \  
    cd /app && \  
    rm -rf /app/tailscale  
  
FROM ubuntu:20.04  
WORKDIR /app  
  
# ========= CONFIG =========  
# - derper args  
ENV DERP_HOST=127.0.0.1  
ENV DERP_CERTS=/app/certs/  
ENV DERP_STUN true  
ENV DERP_VERIFY_CLIENTS false  
# ==========================  
  
# apt  
RUN apt-get update && \  
    apt-get install -y openssl curl  
  
COPY build_cert.sh /app/  
COPY --from=builder /app/derper /app/derper  
  
# build self-signed certs && start derper  
CMD bash /app/build_cert.sh $DERP_HOST $DERP_CERTS /app/san.conf && \  
    /app/derper --hostname=$DERP_HOST \  
    --certmode=manual \  
    --certdir=$DERP_CERTS \  
    --stun=$DERP_STUN  \  
    --verify-clients=$DERP_VERIFY_CLIENTS
```

构建好镜像后，就可以在你想部署 derper 的主机上直接通过该镜像启动 derper 容器了，命令如下：

```
🐳  → docker run --restart always --net host --name derper -d ghcr.io/yangchuansheng/ip_derper
```

和使用域名的方案一样，防火墙需要放行相应端口（12345 与 3478）。

查看容器日志：

```bash
🐳  → docker logs -f derper  
Generating a RSA private key  
.......................................+++++  
..............+++++  
writing new private key to '/app/certs//127.0.0.1.key'  
-----  
2022/03/26 14:30:31 no config path specified; using /var/lib/derper/derper.key  
2022/03/26 14:30:31 derper: serving on :443 with TLS  
2022/03/26 14:30:31 running STUN server on [::]:3478
```

如果你想自己构建 derper 镜像，可以参考我的 GitHub 仓库\[4]。

---

下面就是骚操作了，我们在 Headscale 的配置中需要**将 DERP 的域名设置为 IP**！不理解的可以再消化一下，然后继续往下看哈哈~~

除了 derper 之外，Tailscale 客户端还需要**跳过域名验证**，这个需要在 DERP 的配置中设置。而 Headscale 的本地 YAML 文件目前还不支持这个配置项，所以没办法，咱只能使用在线 URL 了。JSON 配置内容如下：

```
{  
  "Regions": {  
    "901": {  
      "RegionID": 901,  
      "RegionCode": "ali-sh",  
      "RegionName": "Aliyun Shanghai",  
      "Nodes": [  
        {  
          "Name": "901a",  
          "RegionID": 901,  
          "DERPPort": 443,  
          "HostName": "xxxx",  
          "IPv4": "xxxx",  
          "InsecureForTests": true  
        }  
      ]  
    }  
  }  
}
```

配置解析：

- `HostName` 直接填 derper 的公网 IP，即和 `IPv4` 的值相同。
- `InsecureForTests` 一定要设置为 true，以跳过域名验证。

你需要把这个 JSON 文件变成 Headscale 服务器可以访问的 URL，比如在 Headscale 主机上搭个 Nginx，或者上传到对象存储（比如阿里云 OSS）。

接下来还需要修改 Headscale 的配置文件，引用上面的自定义 DERP 的 URL。需要修改的配置项如下：

```bash
# /etc/headscale/config.yaml  
derp:  
  # List of externally available DERP maps encoded in JSON  
  urls:  
  #  - https://controlplane.tailscale.com/derpmap/default  
    - https://xxxxx/derp.json  
  
  # Locally available DERP map files encoded in YAML  
  #  
  # This option is mostly interesting for people hosting  
  # their own DERP servers:  
  # https://tailscale.com/kb/1118/custom-derp-servers/  
  #  
  # paths:  
  #   - /etc/headscale/derp-example.yaml  
  paths:  
    - /etc/headscale/derp.yaml  
  
  # If enabled, a worker will be set up to periodically  
  # refresh the given sources and update the derpmap  
  # will be set up.  
  auto_update_enabled: true  
  
  # How often should we check for DERP updates?  
  update_frequency: 24h
```

修改完配置后，重启 headscale 服务：

`$ systemctl restart headscale`

在 Tailscale 客户端上使用以下命令查看目前可以使用的 DERP 服务器：

```bash
$ tailscale netcheck  
  
Report:  
 * UDP: true  
 * IPv4: yes, 192.168.100.1:49656  
 * IPv6: no  
 * MappingVariesByDestIP: true  
 * HairPinning: false  
 * PortMapping: UPnP  
 * Nearest DERP: Home Hangzhou  
 * DERP latency:  
  - home: 9.7ms   (Home Hangzhou)  
  -  hs: 25.2ms  (Huawei Shanghai)  
  - thk: 43.5ms  (Tencent Hongkong)
```

再次查看与通信对端的连接方式：

```bash
$ tailscale status  
                coredns              default      linux   active; direct xxxx:45986; offline, tx 131012 rx 196020  
                oneplus-8t           default      android active; relay "home"; offline, tx 211900 rx 22780  
                openwrt              default      linux   active; direct 192.168.100.254:41641; offline, tx 189868 rx 4074772
```

可以看到这一次 Tailscale 自动选择了一个线路最优的**国内的** DERP 服务器作为中继，可以测试一下延迟：

```bash
$ tailscale ping 10.1.0.8  
pong from oneplus-8t (10.1.0.8) via DERP(home) in 30ms  
pong from oneplus-8t (10.1.0.8) via DERP(home) in 45ms  
pong from oneplus-8t (10.1.0.8) via DERP(home) in 30ms
```

完美！这里的 home 当然是我的家庭宽带，部署方式与上面所说的国内云主机类似，你需要额外开启公网的端口映射（12345/tcp, 3478/udp）。还有一点需要注意的是配置内容：
```json
{  
  "Regions": {  
    "901": {  
      "RegionID": 901,  
      "RegionCode": "ali-sh",  
      "RegionName": "Aliyun Shanghai",  
      "Nodes": [  
        {  
          "Name": "901a",  
          "RegionID": 901,  
          "DERPPort": 443,  
          "HostName": "xxxx",  
          "IPv4": "xxxx",  
          "InsecureForTests": true  
        }  
      ]  
    },  
    "902": {  
      "RegionID": 902,  
      "RegionCode": "home",  
      "RegionName": "Home Hangzhou",  
      "Nodes": [  
        {  
          "Name": "902a",  
          "RegionID": 902,  
          "DERPPort": 12345,  
          "HostName": "xxxx",  
          "InsecureForTests": true  
        }  
      ]  
    }  
  }  
}
```
与国内云主机相比，家庭宽带的配置有两点不同：

- 需要删除 `IPv4` 配置项。因为家用宽带的公网 IP 是动态变化的，所以你需要使用 **DDNS** 来动态解析公网 IP。
- `HostName` 最好填域名，因为你的公网 IP 是动态变化的，没法填写 IP，除非你不停地修改配置文件。填域名也没关系啦，反正不会验证域名的，也不用关心证书的事情，**只要域名能解析到你的公网 IP 即可。**

## 防止 DERP 被白嫖

默认情况下 DERP 服务器是可以被白嫖的，只要别人知道了你的 DERP 服务器的地址和端口，就可以为他所用。如果你的服务器是个小水管，用的人多了可能会把你撑爆，因此我们需要修改配置来防止被白嫖。

> 特别声明：只有使用域名的方式才可以通过认证防止被白嫖，使用纯 IP 的方式无法防白嫖，你只能小心翼翼地隐藏好你的 IP 和端口，不能让别人知道。

只需要做两件事情：

**1、在 DERP 服务器上安装 Tailscale。**

第一步需要在 DERP 服务所在的主机上安装 Tailscale 客户端，**启动 tailscaled 进程**。

**2、derper 启动时加上参数 **`--verify-clients`**。**

本文推荐的是通过容器启动，Dockerfile 内容\[5]如下：

```dockerfile
FROM golang:latest AS builder

LABEL org.opencontainers.image.source <https://github.com/yangchuansheng/docker-image>

WORKDIR /app

\# <https://tailscale.com/kb/1118/custom-derp-servers/>
RUN go install tailscale.com/cmd/derper@main

FROM ubuntu
WORKDIR /app

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \<br />    apt-get install -y --no-install-recommends apt-utils && \<br />    apt-get install -y ca-certificates && \<br />    mkdir /app/certs

ENV DERP_DOMAIN your-hostname.com
ENV DERP_CERT_MODE letsencrypt
ENV DERP_CERT_DIR /app/certs
ENV DERP_ADDR :443
ENV DERP_STUN true
ENV DERP_HTTP_PORT 80
ENV DERP_VERIFY_CLIENTS false

COPY --from=builder /go/bin/derper .

CMD /app/derper --hostname=![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pq65ei/latex)DERP_HTTP_PORT \<br />    --verify-clients=$DERP_VERIFY_CLIENTS

```

默认情况下 `--verify-clients` 参数设置的是 `false`。我们不需要对 Dockerfile 内容做任何改动，只需在容器启动时加上环境变量即可，将之前的启动命令修改一下：

```bash
docker run --restart always \  
  --name derper -p 12345:12345 -p 3478:3478/udp \  
  -v /root/.acme.sh/xxxx/:/app/certs \  
  -e DERP_CERT_MODE=manual \  
  -e DERP_ADDR=12345 \  
  -e DERP_DOMAIN=xxxx \  
  -e DERP_VERIFY_CLIENTS=true \  
  -d ghcr.io/yangchuansheng/derper:latest
```

这样就大功告成了，别人即使知道了你的 DERP 服务器地址也无法使用，但还是要说明一点，即便如此，你也应该尽量不让别人知道你的服务器地址，防止别人有可趁之机。

# 参考资料

- NAT 穿透是如何工作的：技术原理及企业级实践\[6]
- Encrypted TCP relays (DERP)\[8]

## 引用链接

\[4]我的 GitHub 仓库: [_https://github.com/yangchuansheng/ip_derper_](https://github.com/yangchuansheng/ip_derper)

\[5]Dockerfile 内容: [_https://github.com/yangchuansheng/docker-image/blob/master/derper/Dockerfile_](https://github.com/yangchuansheng/docker-image/blob/master/derper/Dockerfile)

\[6]NAT 穿透是如何工作的：技术原理及企业级实践: [_https://arthurchiao.art/blog/how-nat-traversal-works-zh/_](https://arthurchiao.art/blog/how-nat-traversal-works-zh/)

\[7]Custom DERP Servers: [_https://tailscale.com/kb/1118/custom-derp-servers/_](https://tailscale.com/kb/1118/custom-derp-servers/)

\[8]Encrypted TCP relays (DERP): [_https://tailscale.com/blog/how-tailscale-works/#encrypted-tcp-relays-derp_](https://tailscale.com/blog/how-tailscale-works/#encrypted-tcp-relays-derp)
