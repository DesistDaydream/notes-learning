---
title: Wireguard 流量伪装
linkTitle: Wireguard 流量伪装
date: 2024-04-19T23:33
weight: 20
---

# 概述

> 参考：
>
> -

WireGuard 在国内网络环境下会遇到一个致命的问题：**UDP 封锁/限速**。虽然通过 WireGuard 可以在隧道内传输任何基于 IP 的协议（TCP、UDP、ICMP、SCTP、IPIP、GRE 等），但 WireGuard 隧道本身是通过 UDP 协议进行通信的，而国内运营商根本没有能力和精力根据 TCP 和 UDP 的不同去深度定制不同的 QoS 策略，几乎全部采取一刀切的手段：对 UDP 进行限速甚至封锁。

虽然运营商对 UDP 不友好，但却无力深度检测 TCP 连接的真实性。既然对 TCP 连接睁一只眼闭一只眼，那我将 UDP 连接伪装成 TCP 连接不就蒙混过关了。目前支持将 UDP 流量伪装成 TCP 流量的主流工具是 [udp2raw](https://github.com/wangyu-/udp2raw-tunnel)，另一款比它更强大的新工具：[Phantun](https://github.com/dndx/phantun)。

# udp2raw

> 参考：
> 
> - [GitHub 项目，wangyu-/udp2raw](https://github.com/wangyu-/udp2raw)

使用原始套接字将 UDP 流量转换为加密的 UDP/FakeTCP/ICMP 流量的隧道，帮助您绕过 UDP 防火墙（或不稳定的 UDP 环境）

## 部署

这里我使用 docker 部署的，实体进程和相关文档见 [udp2raw 运行](https://github.com/wangyu-/udp2raw-tunnel/blob/master/doc/README.zh-cn.md#%E8%BF%90%E8%A1%8C)

Linux server 端 :

```bash
# 监听86的tcp端口，把86端口收到的伪装成tcp的udp报文转发到 127.0.0.1:16000 上
docker run \
  -d --name udp2raw  \
  --restart always \
  --net host \
  --cap-add NET_RAW \
  --cap-add NET_ADMIN  \
  -v /run/xtables.lock:/run/xtables.lock \
  zhangguanzhang/udp2raw \
  -s -l 0.0.0.0:86 \
  -r 127.0.0.1:16000 \
  -k passwd123 \
  --raw-mode faketcp  \
  --cipher-mode xor  -a
```

Linux 或者软路由系统 client 端，软路径 openwrt 的话 iptables 的锁文件是位于 /var/run/xtables.lock，常规系统是 \`/run/xtables.lock :

```bash
# 监听16000的 udp 端口，把16000端口收到的udp报文伪装成tcp发到 <public_ip>:86 上
docker run --net host \
  -d --name udp2raw  \
  --restart always \
  --cap-add NET_RAW \
  --cap-add NET_ADMIN    \
  -v /var/run/xtables.lock:/run/xtables.lock \
  zhangguanzhang/udp2raw \
  -c -l 0.0.0.0:16000 \
  -r <public_ip>:86 \
  -k passwd123 \
  --raw-mode faketcp   \
  --cipher-mode xor  -a
```

[windows 客户端下载](https://github.com/wangyu-/udp2raw-multiplatform)，运行命令参考:

```bash
./udp2raw_mp.exe -c -l 0.0.0.0:16000 -r <public_ip>:86 -k passwd123 --raw-mode faketcp --cipher-mode xor
```

所有 client 端的 \[peer] 部分里之前连云主机的 ip 都写成 127.0.0.1:16000，这样 wg 客户端是先向本地的 udp2raw 客户端发 udp 报文，然后报文被封装成 tcp 发往云主机上的 udp2raw server，再到 wg server 上。

**客户端和云主机上** 的 wg 的 mtu 设置成 1280(网上有写 1200 的，但是 windows 的 wg 客户端无法启动，邮件询问作者说最小 1280 才能启动)。例如我路由器配置

```bash
[Interface]
...
MTU = 1280
[Peer]
...
Endpoint = 127.0.0.1:16000
PersistentKeepalive = 10
```

windows 的 wg 目前 Endpoint 必须写本机的 ip（ipconfig 命令查看），不能写 127.0.0.1，否则无法连 peer（日志会一直刷 Failed to send handshake initiation write udp4 0.0.0.0:xxx->127.0.0.1:16000: wsasendto: The requested address is not valid in its context），这个 bug 已经反馈给作者了。
udp2raw 的 client 连上 server 后，双方都会打印下面日志:

```bash
# server
changed state to server_ready
# client
changed state from to client_handshake2 to client_ready
```

## qos

不同运营商可能不一样，比如你 A 和 B 同时 udp2raw 你云主机，A 可以 B 不可以，可以考虑换下 --raw-mode 和 --seq-mode ，有的可能 faketcp，有的可能 udp ，有的可能 icmp

## 一个注意点

openwrt 上在接口 添加 wireguard 接口，然后 peer 那里的 ip 写 127.0.0.1(也就是 openwrt 上的 udp2raw 的 ip)可能不行，换成 openwrt 的 局域网 ip 试下

# Phantun

> 参考：
>
> - [GitHub 项目，dndx/phantun](https://github.com/dndx/phantun)
> - [公众号-米开朗基杨，突破运营商 QoS 封锁，WireGuard 真有“一套”！](https://mp.weixin.qq.com/s/lgOqd7cMy8hFqST9uSih6w)

Phantun 整个项目**完全使用 Rust 实现**，性能吊打 udp2raw。它的初衷和 udp2raw 类似，都是为了实现一种简单的用户态 TCP 状态机来对 UDP 流量做伪装。主要的目的是希望能让 UDP 流量看起来像是 TCP，又不希望受到 TCP retransmission 或者 congestion control 的影响。

需要申明的是，**Phantun 的目标不是为了替代 udp2raw**，从一开始 Phantun 就希望设计足够的简单高效，所以 udp2raw 支持的 **ICMP 隧道，加密，防止重放**等等功能 Phantun 都选择不实现。

Phantun 假设 UDP 协议本身已经解决了这些问题，所以整个转发过程就是简单的明文换头加上一些必要的 TCP   状态控制信息。对于我日常使用的 WireGuard 来说，Phantun 这种设计是足够安全的，因为 WireGuard   的协议已经更好的实现了这些安全功能。

Phantun 使用 TUN 接口来收发 3 层数据包，udp2raw 使用 Raw Socket + BFP 过滤器。个人感觉基于 TUN 的实现要稍微的优雅一点，而且跨平台移植也要更容易。

Phantun 的 TCP 连接是按需创建的，只启动 Client 不会主动去连接服务器，需要第一个数据包到达了后才会按需创建。每个 UDP   流都有自己独立的 TCP 连接。这一点跟 udp2raw 很不一样，udp2raw 所有的 UDP 连接共用一个 TCP 连接。这样做的坏处就是 udp2raw 需要额外的头部信息来区分连接，更加增加了头部的开销。跟纯 UDP 比较，Phantun 每个数据包的额外头部开销是 12  byte，udp2raw 根据我的测试达到了 44 bytes 。

## Phantun 工作原理

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireguard/1669018288492-b85eb3fd-aeea-49b6-b81a-5ccdf8976c28.jpeg)

Phantun 分为服务端和客户端，服务端会监听一个端口，比如 4567（通过 `--local` 参数指定），并将 UDP 数据包转发到 UDP 服务（这里指的就是服务端 WireGuard 的监听端口和地址，通过 `--remote` 参数指定）。

客户端也会监听一个端口，比如 `127.0.0.1:4567`（通过 `--local` 参数指定），并且通过 `--remote` 参数与服务端（比如 `10.0.0.1:4567`）建立连接。

客户端与服务端都会创建一个 TUN 网卡，客户端 TUN 网卡默认分配的 IPv4/IPv6 地址分别是 `192.168.200.2` 和 `fcc8::2`，服务端 TUN 网卡默认分配的 IPv4/IPv6 地址分别是 `192.168.201.2` 和 `fcc9::2`。

客户端与服务端都需要开启 IP forwarding，并且需要创建相应的 NAT 规则。客户端在流量离开物理网卡之前，需要对 IP `192.168.200.2` 进行 SNAT；服务端在流量进入网卡之前，需要将 IP DNAT 为 `192.168.201.2`。

## Phantun 配置步骤

接下来我会通过一个示例来演示如何使用 Phantun 将 WireGuard 的 UDP 流量伪装成 TCP。我们需要在服务端和客户端分别安装 phantun，可以到 release 页面\[3]下载，推荐下载静态编译版本 `phantun_x86_64-unknown-linux-musl.zip`。

### 服务端

假设服务端的公网 IP 地址是 `121.36.134.95`，WireGuard 监听端口是 `51822`。首先修改配置文件 `/etc/wireguard/wg0.conf`，在 `[Interface]` 中添加以下配置：

```properties
MTU = 1300
PreUp = iptables -t nat -A PREROUTING -p tcp -i eth0 --dport 4567 -j DNAT --to-destination 192.168.201.2
PreUp = RUST_LOG=info phantun_server --local 4567 --remote 127.0.0.1:51822 &> /var/log/phantun_server.log &
PostDown = iptables -t nat -D PREROUTING -p tcp -i eth0 --dport 4567 -j DNAT --to-destination 192.168.201.2
PostDown = killall phantun_server || true
```

你需要将 eth0 替换为你服务端的物理网卡名。MTU 值先不管，后面再告诉大家调试方法。

```properties
PreUp = iptables -t nat -A PREROUTING -p tcp -i eth0 --dport 4567 -j DNAT --to-destination 192.168.201.2
```

这条 iptables 规则表示将 `4567` 端口的入站流量 DNAT 为 TUN 网卡的 IP 地址。

```properties
PreUp = RUST_LOG=info phantun_server --local 4567 --remote 127.0.0.1:51822 &> /var/log/phantun_server.log &
```

这里会启动 phantun_server，监听在 `4567` 端口，并将 UDP 数据包转发到 WireGuard。
服务端完整的 WireGuard 配置：

```properties
# local settings for Endpoint B
[Interface]
PrivateKey = QH1BJzIZcGo89ZTykxls4i2DKgvByUkHIBy3BES2gX8=
Address = 10.0.0.2/32
ListenPort = 51822
MTU = 1300
PreUp = iptables -t nat -A PREROUTING -p tcp -i eth0 --dport 4567 -j DNAT --to-destination 192.168.201.2
PreUp = RUST_LOG=info phantun_server --local 4567 --remote 127.0.0.1:51822 &> /var/log/phantun_server.log &
PostDown = iptables -t nat -D PREROUTING -p tcp -i eth0 --dport 4567 -j DNAT --to-destination 192.168.201.2
PostDown = killall phantun_server || true

# remote settings for Endpoint A
[Peer]
PublicKey = wXtD/VrRo92JHc66q4Ypmnd4JpMk7b1Sb0AcT+pJfwY=
AllowedIPs = 10.0.0.1/32
```

最后重启 WireGuard 即可：`systemctl restart wg-quick@wg0`

### 客户端

假设客户端的 WireGuard 监听端口是 `51821`。首先修改配置文件 `/etc/wireguard/wg0.conf`，在 `[Interface]` 中添加以下配置：

```properties
MTU = 1300
PreUp = iptables -t nat -A POSTROUTING -o eth0 -s 192.168.200.2 -j MASQUERADE
PreUp = RUST_LOG=info phantun_client --local 127.0.0.1:4567 --remote 121.36.134.95:4567 &> /var/log/phantun_client.log &
PostDown = iptables -t nat -D POSTROUTING -o eth0 -s 192.168.200.2 -j MASQUERADE
PostDown = killall phantun_client || true
```

你需要将 eth0 替换为你服务端的物理网卡名。

```properties
PreUp = iptables -t nat -A POSTROUTING -o eth0 -s 192.168.200.2 -j MASQUERADE
```

这条 iptables 规则表示对来自 `192.168.200.2`（TUN 网卡） 的出站流量进行 MASQUERADE。

```properties
PreUp = RUST_LOG=info phantun_client --local 127.0.0.1:4567 --remote 121.36.134.95:4567 &> /var/log/phantun_client.log &
```

这里会启动 phantun_client，监听在 `4567` 端口，并与服务端建立连接，将伪装的 TCP 数据包传送给服务端。
除此之外还需要修改 WireGuard peer 的 Endpoint，将其修改为 127.0.0.1:4567。

```properties
Endpoint = 127.0.0.1:4567
```

客户端完整的 WireGuard 配置：

```properties
# local settings for Endpoint A
[Interface]
PrivateKey = 0Pyz3cIg2gRt+KxZ0Vm1PvSIU+0FGufPIzv92jTyGWk=
Address = 10.0.0.1/32
ListenPort = 51821
MTU = 1300
PreUp = iptables -t nat -A POSTROUTING -o eth0 -s 192.168.200.2 -j MASQUERADE
PreUp = RUST_LOG=info phantun_client --local 127.0.0.1:4567 --remote 121.36.134.95:4567 &> /var/log/phantun_client.log &
PostDown = iptables -t nat -D POSTROUTING -o eth0 -s 192.168.200.2 -j MASQUERADE
PostDown = killall phantun_client || true

# remote settings for Endpoint B
[Peer]
PublicKey = m40NDb5Cqtb78b1DVwY1+kxbG2yEcRhxlrLm/DlPpz8=
Endpoint = 127.0.0.1:4567
AllowedIPs = 10.0.0.2/32
PersistentKeepalive = 25
```

最后重启 WireGuard 即可：`systemctl restart wg-quick@wg0`

查看 phantun_client 的日志：

```bash
$ tail -f /var/log/phantun_client.log
 INFO  client > Remote address is: 121.36.134.95:4567
 INFO  client > 1 cores available
 INFO  client > Created TUN device tun0
 INFO  client > New UDP client from 127.0.0.1:51821
 INFO  fake_tcp > Sent SYN to server
 INFO  fake_tcp > Connection to 121.36.134.95:4567 established
```

查看 wg0 接口：

```bash
$ wg show wg0
interface: wg0
  public key: wXtD/VrRo92JHc66q4Ypmnd4JpMk7b1Sb0AcT+pJfwY=
  private key: (hidden)
  listening port: 51821

peer: m40NDb5Cqtb78b1DVwY1+kxbG2yEcRhxlrLm/DlPpz8=
  endpoint: 127.0.0.1:4567
  allowed ips: 10.0.0.2/32
  latest handshake: 1 minute, 57 seconds ago
  transfer: 184 B received, 648 B sent
  persistent keepalive: every 25 seconds
```

测试连通性：

```bash
$ ping 10.0.0.2 -c 3
PING 10.0.0.2 (10.0.0.2) 56(84) bytes of data.
64 bytes from 10.0.0.2: icmp_seq=1 ttl=64 time=13.7 ms
64 bytes from 10.0.0.2: icmp_seq=2 ttl=64 time=14.4 ms
64 bytes from 10.0.0.2: icmp_seq=3 ttl=64 time=15.0 ms

--- 10.0.0.2 ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 2005ms
rtt min/avg/max/mdev = 13.718/14.373/15.047/0.542 ms
```

### 客户端（多服务端）

如果客户端想和多个服务端建立连接，则新增的服务端配置如下：

```properties
PreUp = RUST_LOG=info phantun_client --local 127.0.0.1:4568 --remote xxxx:4567 --tun-local=192.168.202.1 --tun-peer=192.168.202.2 &> /var/log/phantun_client.log &
PostDown = iptables -t nat -D POSTROUTING -o eth0 -s 192.168.202.2 -j MASQUERADE
```

本地监听端口需要选择一个与之前不同的端口，同理，TUN 网卡的地址也需要修改。最终的配置如下：

```properties
# local settings for Endpoint A
[Interface]
PrivateKey = 0Pyz3cIg2gRt+KxZ0Vm1PvSIU+0FGufPIzv92jTyGWk=
Address = 10.0.0.1/32
ListenPort = 51821
MTU = 1300
PreUp = iptables -t nat -A POSTROUTING -o eth0 -s 192.168.200.2 -j MASQUERADE
PreUp = RUST_LOG=info phantun_client --local 127.0.0.1:4567 --remote 121.36.134.95:4567 &> /var/log/phantun_client.log &
PreUp = RUST_LOG=info phantun_client --local 127.0.0.1:4568 --remote xxxx:4567 --tun-local=192.168.202.1 --tun-peer=192.168.202.2 &> /var/log/phantun_client.log &
PostDown = iptables -t nat -D POSTROUTING -o eth0 -s 192.168.200.2 -j MASQUERADE
PostDown = iptables -t nat -D POSTROUTING -o eth0 -s 192.168.202.2 -j MASQUERADE
PostDown = killall phantun_client || true

# remote settings for Endpoint B
[Peer]
PublicKey = m40NDb5Cqtb78b1DVwY1+kxbG2yEcRhxlrLm/DlPpz8=
Endpoint = 127.0.0.1:4567
AllowedIPs = 10.0.0.2/32
PersistentKeepalive = 25
```

## MTU 调优

如果你使用 ping 或者 dig 等工具（小数据包）测试 WireGuard 隧道能够正常工作，但浏览器或者远程桌面（大数据包）却无法正常访问，很有可能是 MTU 的问题，你需要将 MTU 的值调小一点。

Phantun 官方建议将 MTU 的值设为 `1428`（假设物理网卡的 MTU 是 1500），但经我测试是有问题的。建议直接将 MTU 设置为最低值 `1280`，然后渐渐增加，直到无法正常工作为止，此时你的 MTU 就是最佳值。
