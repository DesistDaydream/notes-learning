---
title: Wireguard配置详解
---

# 概述

> 参考：
> 
> - <https://zhangguanzhang.github.io/2020/08/05/wireguard-for-personal/>
> - <https://fuckcloudnative.io/posts/wireguard-docs-practice/#peer>

WireGuard 使用 [INI](/docs/2.编程/无法分类的语言/INI.md) 作为其配置文件格式。配置文件可以放在任何路径下，但必须通过绝对路径引用。默认路径是 `/etc/wireguard/*.conf`。

配置文件的命名形式必须为 `${WireGuard_Interface_Name}.conf`。通常情况下 WireGuard 接口名称以 `wg` 为前缀，并从 `0` 开始编号，但你也可以使用其他名称，只要符合正则表达式 `^[a-zA-Z0-9_=+.-]{1,15}$` 就行。当启动时，如果配置文件中有 wg0.conf 文件，则会创建一个名为 wg0 的网络设备。效果如下

注意：`${WireGuard_Interface_Name}` 不能过长，否则将会报错：`wg-quick: The config file must be a valid interface name, followed by .conf`

## 基本配置示例

```ini
[Interface]
# Name = node1.example.tld
Address = 192.0.2.3/32
ListenPort = 51820
PrivateKey = localPrivateKeyAbcAbcAbc=
DNS = 1.1.1.1,8.8.8.8
Table = 12345
MTU = 1500
PreUp = /bin/example arg1 arg2 %i
PostUp = /bin/example arg1 arg2 %i
PreDown = /bin/example arg1 arg2 %i
PostDown = /bin/example arg1 arg2 %i

[Peer]
# Name = node2-node.example.tld
AllowedIPs = 192.0.2.1/24
Endpoint = node1.example.tld:51820
PublicKey = remotePublicKeyAbcAbcAbc=
PersistentKeepalive = 25
```

# \[Interface]

Interface 部分定义本地 VPN 配置。例如：

- 本地节点是客户端，只路由自身的流量，只暴露一个 IP。

```ini
[Interface]
# Name = phone.example-vpn.dev
Address = 192.0.2.5/32
PrivateKey = <private key for phone.example-vpn.dev>
```

- 本地节点是中继服务器，它可以将流量转发到其他对等节点（peer），并公开整个 VPN 子网的路由。

```ini
[Interface]
# Name = public-server1.example-vpn.tld
Address = 192.0.2.1/24
ListenPort = 51820
PrivateKey = <private key for public-server1.example-vpn.tld>
DNS = 1.1.1.1
```

## Address

定义本地节点应该对哪个地址范围进行路由。如果是常规的客户端，则将其设置为节点本身的单个 IP（使用 CIDR 指定，例如 192.0.2.3/32）；如果是中继服务器，则将其设置为可路由的子网范围。
例如：

- 常规客户端，只路由自身的流量：`Address = 192.0.2.3/32`
- 中继服务器，可以将流量转发到其他对等节点（peer）：`Address = 192.0.2.1/24`
- 也可以指定多个子网或 IPv6 子网：`Address = 192.0.2.1/24,2001:DB8::/64`

## ListenPort

当本地节点是中继服务器时，需要通过该参数指定端口来监听传入 VPN 连接，默认端口号是 `51820`。常规客户端不需要此选项。

## PrivateKey

本地节点的私钥，所有节点（包括中继服务器）都必须设置。不可与其他服务器共用。

私钥可通过命令 `wg genkey > example.key` 来生成。

## DNS

通过 DHCP 向客户端宣告 DNS 服务器。客户端将会使用这里指定的 DNS 服务器来处理 VPN 子网中的 DNS 请求，但也可以在系统中覆盖此选项。例如：

- 如果不配置则使用系统默认 DNS
- 可以指定单个 DNS：`DNS = 1.1.1.1`
- 也可以指定多个 DNS：`DNS = 1.1.1.1,8.8.8.8`

## Table

定义 VPN 子网使用的路由表，默认不需要设置。该参数有两个特殊的值需要注意：

- **Table = off** : 禁止创建路由
- **Table = auto（默认值）** : 将路由添加到系统默认的 table 中，并启用对默认路由的特殊处理。

例如：`Table = 1234`

## MTU

定义连接到对等节点（peer）的 `MTU`（Maximum Transmission Unit，最大传输单元），默认不需要设置，一般由系统自动确定。

## PreUp

启动 VPN 接口之前运行的命令。这个选项可以指定多次，按顺序执行。
例如：

- 添加路由：`PreUp = ip rule add ipproto tcp dport 22 table 1234`

## PostUp

启动 VPN 接口之后运行的命令。这个选项可以指定多次，按顺序执行。

例如：

- 从文件或某个命令的输出中读取配置值：
  - PostUp = wg set %i private-key /etc/wireguard/wg0.key <(some command here)
- 添加一行日志到文件中：
  - PostUp = echo "$(date +%s) WireGuard Started" >> /var/log/wireguard.log
- 调用 WebHook：
  - PostUp = curl https://events.example.dev/wireguard/started/?key=abcdefg
- 添加路由：
  - PostUp = ip rule add ipproto tcp dport 22 table 1234
- 添加 iptables 规则，启用数据包转发：
  - PostUp = iptables -A FORWARD -i %i -j ACCEPT; iptables -A FORWARD -o %i -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
- 强制 WireGuard 重新解析对端域名的 IP 地址：
  - PostUp = resolvectl domain %i "~."; resolvectl dns %i 192.0.2.1; resolvectl dnssec %i yes

## PreDown

停止 VPN 接口之前运行的命令。这个选项可以指定多次，按顺序执行。
例如：

- 添加一行日志到文件中：
  - PreDown = echo "$(date +%s) WireGuard Going Down" >> /var/log/wireguard.log
- 调用 WebHook：
  - PreDown = curl https://events.example.dev/wireguard/stopping/?key=abcdefg

## PostDown

停止 VPN 接口之后运行的命令。这个选项可以指定多次，按顺序执行。
例如：

- 添加一行日志到文件中：
  - PostDown = echo "$(date +%s) WireGuard Going Down" >> /var/log/wireguard.log
- 调用 WebHook：
  - PostDown = curl https://events.example.dev/wireguard/stopping/?key=abcdefg
- 删除 iptables 规则，关闭数据包转发：
  - PostDown = iptables -D FORWARD -i %i -j ACCEPT; iptables -D FORWARD -o %i -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE

# \[Peer]

定义 Peer，用来将一个或多个地址路由流量的 Peer。Peer 可以是将流量转发到其他 Peer 的中继服务器，也可以是通过公网或内网直连的客户端。

中继服务器必须将所有的客户端定义为对等节点（peer），除了中继服务器之外，其他客户端都不能将位于 NAT 后面的节点定义为对等节点（peer），因为路由不可达。对于那些只为自己路由流量的客户端，只需将中继服务器作为对等节点（peer），以及其他需要直接访问的节点。

举个例子，在下面的配置中，`public-server1` 作为中继服务器，其他的客户端有的是直连，有的位于 NAT 后面：

- `public-server1`（中继服务器）
  \[peer] : `public-server2`, `home-server`, `laptop`, `phone`
- `public-server2`（直连客户端）
  \[peer] : `public-server1`
- `home-server`（客户端位于 NAT 后面）
  \[peer] : `public-server1`, `public-server2`
- `laptop`（客户端位于 NAT 后面）
  \[peer] : `public-server1`, `public-server2`
- `phone`（客户端位于 NAT 后面）
  \[peer] : `public-server1`, `public-server2`

配置示例：

- 对等节点（peer）是路由可达的客户端，只为自己路由流量

```ini
[Peer]
# Name = public-server2.example-vpn.dev
Endpoint = public-server2.example-vpn.dev:51820
PublicKey = <public key for public-server2.example-vpn.dev>
AllowedIPs = 192.0.2.2/32
```

- 对等节点（peer）是位于 NAT 后面的客户端，只为自己路由流量

```ini
[Peer]
# Name = home-server.example-vpn.dev
Endpoint = home-server.example-vpn.dev:51820
PublicKey = <public key for home-server.example-vpn.dev>
AllowedIPs = 192.0.2.3/32
```

- 对等节点（peer）是中继服务器，用来将流量转发到其他对等节点（peer）

```ini
[Peer]
# Name = public-server1.example-vpn.tld
Endpoint = public-server1.example-vpn.tld:51820
PublicKey = <public key for public-server1.example-vpn.tld>
# 路由整个 VPN 子网的流量
AllowedIPs = 192.0.2.1/24
PersistentKeepalive = 25
```

## PublicKey

Peer 的公钥，所有节点（包括中继服务器）都必须设置。可与其他对等节点（peer）共用同一个公钥。

公钥可通过命令 `wg pubkey <example.key> example.key.pub` 来生成，其中 `example.key` 是上面生成的私钥。

例如：`PublicKey = somePublicKeyAbcdAbcdAbcdAbcd=`

## Endpoint

指定 其他 Peer 的公网地址。如果 Peer 位于 NAT 后面或者没有稳定的公网访问地址，则忽略这个字段。通常只需要指定**中继服务器**的 `Endpoint`，当然有稳定公网 IP 的节点也可以指定。例如：

- 通过 IP 指定：

Endpoint = 123.124.125.126:51820

- 通过域名指定：

Endpoint = public-server1.example-vpn.tld:51820

NAT 后的任何 Peer 都会将 AllowedIPs 指定网段的数据包，发送到 Endpoint。

## AllowedIPs

<font color="#ff0000">核心配置</font>

AllowedIPs 有两层含义：

- 其他 Peer 向本 Peer 发送数据包时，只有源地址在该字段指定的地址范围内时，才会处理这些数据包，否则直接丢弃。
- 本 Peer 会根据该字段的地址范围，在路由表中添加路由条目。凡是发送给 AllowedIPs 字段中指定地址的数据包，都会通过 WireGuard 管理的网络设备处理后发往该 Peer。比如：

```bash
[Interface]
Address = 10.1.0.254/24
[Peer]
AllowedIPs = 10.1.0.1/32, 172.16.0.0/24
```

则路由表会出现如下条目

```bash
10.1.0.0/24 dev wg0 proto kernel scope link src 10.1.0.254
172.16.0.0/24 dev wg0 scope link
```

**如果 Peer 是常规的客户端，则 AllowedIPs 其设置为节点本身的单个 IP；如果 Peer 是 relay(中继) 服务器，则将 AllowedIPs 设置为可路由的子网范围。可以使用 `,` 来指定多个 IP 或子网范围。该字段也可以指定多次。**

本质上 **Endpoint 与 AllowedIPs 两个字段将会组成路由条目，可以这么描述：`目的地址是 AllowedIPs 的数据包，下一跳是 Endpoint`。也可以不指定 Endpoint，而是仅仅将数据包送入 WireGuard 创建的网络设备中**。

当决定如何对一个数据包进行路由时，系统首先会选择最具体的路由，如果不匹配再选择更宽泛的路由。例如，对于一个发往 `192.0.2.3` 的数据包，系统首先会寻找地址为 `192.0.2.3/32` 的对等节点（peer），如果没有再寻找地址为 `192.0.2.1/24` 的对等节点（peer），以此类推。

例如：

- 对等节点（peer）是常规客户端，只路由自身的流量：

```ini
AllowedIPs = 192.0.2.3/32
```

- 对等节点（peer）是中继服务器，可以将流量转发到其他对等节点（peer）：

```ini
AllowedIPs = 192.0.2.1/24
```

- 对等节点（peer）是中继服务器，可以转发所有的流量，包括外网流量和 VPN 流量，可以用来干嘛你懂得：

```ini
AllowedIPs = 0.0.0.0/0,::/0
```

- 对等节点（peer）是中继服务器，可以路由其自身和其他对等节点（peer）的流量：

```ini
AllowedIPs = 192.0.2.3/32,192.0.2.4/32
```

- 对等节点（peer）是中继服务器，可以路由其自身的流量和它所在的内网的流量：

```ini
AllowedIPs = 192.0.2.3/32,192.168.1.1/24
```

## PersistentKeepalive

如果连接是从一个位于 NAT 后面的对等节点（peer）到一个公网可达的对等节点（peer），那么 NAT 后面的对等节点（peer）必须定期发送一个出站 ping 包来检查连通性，如果 IP 有变化，就会自动更新 `Endpoint`。

例如：

- 本地节点与对等节点（peer）可直连：该字段不需要指定，因为不需要连接检查。
- 对等节点（peer）位于 NAT 后面：该字段不需要指定，因为维持连接是客户端（连接的发起方）的责任。
- 本地节点位于 NAT 后面，对等节点（peer）公网可达：需要指定该字段 `PersistentKeepalive = 25`，表示每隔 `25` 秒发送一次 ping 来检查连接。

# 配置示例

## 最简单的配置

在带公网的 Peer 上执行如下指令生成三对公私钥、以及 3 个配置文件

```bash
wg genkey | tee gw-privatekey | wg pubkey > gw-publickey
wg genkey | tee peer1-privatekey | wg pubkey > peer1-publickey
wg genkey | tee peer2-privatekey | wg pubkey > peer2-publickey
wg genkey | tee peer3-privatekey | wg pubkey > peer3-publickey
```

带公网的 Peer 配置文件

```bash
cat > wg0.conf <<EOF
[Interface]
# 除了本 Peer 以外的其他 Peer 连接本 Peer 所使用的端口，也表示本 Peer 监听的 UDP 端口号。
ListenPort = 16000
# 用来表示 WireGuard 在本 Peer 上创建的网络设备的 IP，每个 Peer 都是独立的
# 说白了，就是本 Peer 用来与其他 Peer 通信的网络设备的 IP 地址和网段
Address = 10.1.0.254/24
PrivateKey = $(cat gw-privatekey)
# 用来处理流量转发的 iptables 规则
PostUp   = iptables -A FORWARD -i %i -j ACCEPT; iptables -A FORWARD -o %i -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
PostDown = iptables -D FORWARD -i %i -j ACCEPT; iptables -D FORWARD -o %i -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE

# 定义一些 Peer，用来定义将目的地址是哪些的数据包发给哪个 Peer（AllowedIPs 决定的）。
# peer1
[Peer]
PublicKey = $(cat peer1-publickey)
AllowedIPs = 10.1.0.1/32, 172.16.0.0/24

# peer2
[Peer]
PublicKey = $(cat peer2-publickey)
AllowedIPs = 10.1.0.2/32, 172.19.42.0/24

# peer3
[Peer]
PublicKey = $(cat peer3-publickey)
AllowedIPs = 10.1.0.3/32, 192.168.31.0/24
EOF
```

Peer1 配置文件

```bash
cat > peer1.conf <<EOF
[Interface]
PrivateKey = $(cat peer1-privatekey)
# 用来表示 WireGuard 在本 Peer 上创建的网络设备的 IP，每个 Peer 都是独立的
Address = 10.1.0.1/24

# 定义一个 Peer（i.e.公网的 Peer）。
[Peer]
# 该 Peer 的公钥
PublicKey = $(cat gw-publickey)
# 目的地址是下面这些 IP 或 IP 段的请求发给这个 Peer
AllowedIPs = 10.1.0.0/24, 172.19.42.0/24, 192.168.31.0/24
# 该 Peer 的公网 IP 和 PORT
Endpoint = $(curl -s ip.sb):16000
PersistentKeepalive = 60
EOF
```

Peer2 配置文件

```bash
cat > peer2.conf <<EOF
[Interface]
PrivateKey = $(cat peer2-privatekey)
Address = 10.1.0.2/24

[Peer]
PublicKey = $(cat gw-publickey)
AllowedIPs = 10.1.0.0/24, 172.16.0.0/24, 192.168.31.0/24
Endpoint = $(curl -s ip.sb):16000
PersistentKeepalive = 10
EOF
```

Peer3 配置文件

```bash
cat > peer3.conf <<EOF
[Interface]
PrivateKey = $(cat peer3-privatekey)
Address = 10.1.0.3/24

[Peer]
PublicKey = $(cat gw-publickey)
AllowedIPs = 10.1.0.0/24, 172.19.42.0/24, 172.16.0.0/24
Endpoint = $(curl -s ip.sb):16000
PersistentKeepalive = 10
EOF
```

此时，在 Peer1 上使用 peer1.conf 文件启动 WireGuard；在 Peer2 上使用 peer2.conf 文件启动 WireGuard；在 Peer3 上使用 peer3.conf 文件启动 WireGuard。然后 Peer{1,2,3} 这三个节点就可以互通了~

# Unit 文件

WireGuard 使用包管理器安装后，会自动创建一个由 Systemd 管理的 Unit 文件：

```bash
[root@hw-cloud-xngy-jump-server-linux-1 ~]# systemctl cat wg-quick@.service
# /lib/systemd/system/wg-quick@.service
[Unit]
Description=WireGuard via wg-quick(8) for %I
After=network-online.target nss-lookup.target
Wants=network-online.target nss-lookup.target
PartOf=wg-quick.target
Documentation=man:wg-quick(8)
Documentation=man:wg(8)
Documentation=https://www.wireguard.com/
Documentation=https://www.wireguard.com/quickstart/
Documentation=https://git.zx2c4.com/wireguard-tools/about/src/man/wg-quick.8
Documentation=https://git.zx2c4.com/wireguard-tools/about/src/man/wg.8

[Service]
Type=oneshot
RemainAfterExit=yes
ExecStart=/usr/bin/wg-quick up %i
ExecStop=/usr/bin/wg-quick down %i
Environment=WG_ENDPOINT_RESOLUTION_RETRIES=infinity

[Install]
WantedBy=multi-user.target
```

我们可以通过 Unit 来配置 WireGuard 的开机自启，将 WireGuard 的维护工作交给 Systemd。
