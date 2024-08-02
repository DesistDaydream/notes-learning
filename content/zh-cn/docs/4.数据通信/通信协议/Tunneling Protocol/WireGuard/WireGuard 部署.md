---
title: WireGuard 部署
---

# 概述

> 参考：
>
> - 原文链接：<https://mp.weixin.qq.com/s/vbt30eEGcp5JP5sHAPkwhw>
> - 英文原文链接：<https://github.com/pirate/wireguard-docs>

# 安装 WireGuard

## Unix-like OS

安装 WireGuard 本质是安装如下几个程序：

- **wireguard** # WireGuard 的实现。wireguard 程序有两种实现方式，分为是**用户态的实现**与**内核态的实现**。
- CLI
  - **wg** # 管理 WireGuard 网络。生成密钥、设置网络设备信息、查看状态、etc.
  - **wg-quick** # 用户友好的命令行工具，可以通过配置文件管理调用 wg 程序以启动或停止 WireGuard

> [!Tip] wireguard 用户态与内核态的实现
>
> 绝大多数 [Unix-like OS](/docs/1.操作系统/Operating%20system/Unix-like%20OS/Unix-like%20OS.md) 发行版的内核版本（[Linux Kernel 5.6 版本](https://github.com/torvalds/linux/blob/v5.6/drivers/net/wireguard/version.h)开始）通常都内置了 wireguard 模块。
>
> 有些设备内核没有 wireguard 模块，使用 `modprob wireguard` 命令加载模块时报错 `modprobe: FATAL: Module wireguard not found in directory /lib/modules/$(uanme -r)`，此时可以利用用户态的 Wireguard 程序实现 Wireguard 互联。
>
> - https://github.com/WireGuard/wireguard-go 是使用 Go 语言在用户态实现的 Wireguard
> - https://github.com/cloudflare/boringtun 是使用 Rust 语言在用户态实现的 Wireguard
> - etc.
>
> 当系统中存在多种 Wireguard 时，内核态优先级高，用户态优先级低。例如数据包发过来都是先让内核简单 hook 下 iptables 几个链，最后再发给进程。

在中继服务器上开启 IP 地址转发：

```bash
cat > /etc/sysctl.d/wireguard.conf <<EOF
net.ipv4.ip_forward = 1
net.ipv4.conf.all.proxy_arp = 1
EOF

sysctl -p /etc/sysctl.conf
```

### CentOS 7

```bash
yum install epel-release.noarch elrepo-release.noarch -y
yum install --enablerepo=elrepo-kernel kmod-wireguard wireguard-tools -y
```

```bash
# 如果使用的是非标准内核，需要安装 DKMS 包，待验证
yum install https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm
curl -o /etc/yum.repos.d/jdoss-wireguard-epel-7.repo https://copr.fedorainfracloud.org/coprs/jdoss/wireguard/repo/epel-7/jdoss-wireguard-epel-7.repo
yum install wireguard-dkms wireguard-tools

# CentOS Stream
yum install epel-release.noarch elrepo-release.noarch -y
yum install --enablerepo=elrepo-kernel kmod-wireguard wireguard-tools
```

### Ubuntu

```bash
# Ubuntu ≥ 18.04
apt install wireguard
```

### 其他

对于不支持内核模块的系统，可以使用用户态的实现，这里用 wireguard-go 举例

```bash
git clone https://git.zx2c4.com/wireguard-go
cd wireguard-go
make
```

在任意一台机器构建好二进制包，拷贝到需要运行 WireGuard 的机器放到 $PATH 下即可

安装 WireGuard 工具

```bash
git clone https://github.com/WireGuard/wireguard-tools.git
cd src
make install
```

wg 工具会自动找到 wireguard-go 程序。正常使用 wg 和 wg-quick 即可

### 利用 DKMS

对于默认不支持内核模块的系统，除了使用用户态的实现外，还可以通过 DKMS 动态加载模块。

略

## MicrosoftOS

```powershell
# Windows 客户端下载地址：
https://download.wireguard.com/windows-client/wireguard-amd64-0.1.1.msi
```

GUI 的程序，直接使用即可。

# 编写配置文件

配置文件可以放在任何路径下，但必须通过绝对路径引用。默认路径是 `/etc/wireguard/`。下面通过三个 Peer 组成网络拓扑的配置进行演示

## 生成密钥

```bash
# 生成中继服务器 Peer 的公钥与私钥
wg genkey | tee /etc/wireguard/key/relay-privatekey | wg pubkey > /etc/wireguard/key/relay-publickey
# 生成其他 Peer 的公钥与私钥
wg genkey | tee /etc/wireguard/key/peer-client-privatekey | wg pubkey > /etc/wireguard/key/peer-client-publickey
wg genkey | tee /etc/wireguard/key/peer-company-privatekey | wg pubkey > /etc/wireguard/key/peer-company-publickey
```

## 在中继服务器上，生成所有 Peer 的配置文件

### 配置中继服务器 Peer

```bash
cat > /etc/wireguard/wg-company.conf <<EOF
[Interface]
ListenPort = 16000
Address = 10.1.0.254/24
PrivateKey = $(cat /etc/wireguard/key/relay-privatekey)
PostUp   = iptables -A FORWARD -i %i -j ACCEPT; iptables -A FORWARD -o %i -j ACCEPT; iptables -t nat -A POSTROUTING -o ens3 -j MASQUERADE
PostDown = iptables -D FORWARD -i %i -j ACCEPT; iptables -D FORWARD -o %i -j ACCEPT; iptables -t nat -D POSTROUTING -o ens3 -j MASQUERADE

[Peer]
# 客户端
PublicKey = $(cat /etc/wireguard/key/peer-client-publickey)
AllowedIPs = 10.1.0.253/32

[Peer]
# 天津-公司
PublicKey = $(cat /etc/wireguard/key/peer-company-publickey)
AllowedIPs = 10.1.0.1/32, 172.38.0.0/16, 10.20.5.0/24
EOF
```

### 配置其他 Peer

生成配置后，将配置文件拷贝到对应 Peer 的 /etc/wireguard 目录下

```bash
cat > /etc/wireguard/client-company.conf <<EOF
[Interface]
PrivateKey = $(cat /etc/wireguard/key/peer-client-privatekey)
Address = 10.1.0.253/24

[Peer]
PublicKey = $(cat /etc/wireguard/key/relay-publickey)
AllowedIPs = 10.1.0.0/24, 10.20.5.0/24, 172.38.0.0/16
Endpoint = $(curl -4 -s ip.sb):16000
PersistentKeepalive = 30
EOF
```

```bash
cat > /etc/wireguard/peer-company.conf <<EOF
[Interface]
PrivateKey = $(cat /etc/wireguard/key/peer-company-privatekey)
Address = 10.1.0.1/24
# 由于需要通过家里的 Peer 访问公司内很多网段，所以公司内的 Peer 同样需要开启转发以访问其它网段
PostUp   = iptables -A FORWARD -i %i -j ACCEPT; iptables -A FORWARD -o %i -j ACCEPT; iptables -t nat -A POSTROUTING -o ens33 -j MASQUERADE
PostDown = iptables -D FORWARD -i %i -j ACCEPT; iptables -D FORWARD -o %i -j ACCEPT; iptables -t nat -D POSTROUTING -o ens33 -j MASQUERADE

[Peer]
PublicKey = $(cat /etc/wireguard/key/relay-publickey)
AllowedIPs = 10.1.0.0/24, 192.168.0.0/24
Endpoint = $(curl -4 -s ip.sb):16000
PersistentKeepalive = 30
EOF
```

# 启动与停止

启动中继服务器 Peer

```bash
wg-quick up /etc/wireguard/company-wg.conf
```

启动其他 Peer

```bash
wg-quick up /etc/wireguard/company-client.conf
wg-quick up /etc/wireguard/company.conf
```

服务启动时，本质是只是执行了几条命令，比如我现在使用 wg-company 配置启动 WireGuard：

```bash
~]# wg-quick up /etc/wireguard/company-wg.conf
[#] ip link add company-wg type wireguard
[#] wg setconf company-wg /dev/fd/63
[#] ip -4 address add 10.1.0.254/24 dev company-wg
[#] ip link set mtu 1420 up dev company-wg
[#] ip -4 route add 10.20.5.0/24 dev company-wg
[#] ip -4 route add 172.38.0.0/16 dev company-wg
[#] iptables -A FORWARD -i company-wg -j ACCEPT; iptables -A FORWARD -o company-wg -j ACCEPT; iptables -t nat -A POSTROUTING -o ens3 -j MASQUERADE

```

可以看到执行了如下几部操作

- 创建 wireguard 类型的 网络设备，并设置该网络设备
  - ip link add wg-company type wireguar
  - ip link set mtu 1420 up dev wg-company
- 根据 `[Interface]` 部分的配置，为新添加的网络设备添加 IP 地址
  - ip -4 address add 10.1.0.254/24 dev wg-company
- 根据所有 `[Peer]` 部分的配置，为主机添加路由条目
  - ip -4 route add 10.20.5.0/24 dev wg-company
  - ip -4 route add 172.38.0.0/16 dev wg-company
- 为了让中继服务器可以转发数据包，需要配置 Netfilter 规则。这个规则，可以在 `[Interface]` 部分的配置中通过 PostUp 字段定义
  - iptables -A FORWARD -i wg-company -j ACCEPT
  - iptables -A FORWARD -o wg-company -j ACCEPT
  - iptables -t nat -A POSTROUTING -o ens3 -j MASQUERADE

可以看到，即使我们不使用 wg-quick 命令，通过上述操作，同样可以激活 WireGuard，毕竟，WireGuard 已经被包含在 Linux Kernel 当中了，我们只需要创建出来 WireGuard 类型的网络设备，并配置好路由条目，即可转发数据包，实现 VPN 的功能。

## 通过 systemd 启动 Wireguard

WireGuard 安装完成后，会生成一个 `wg-qucik@.service` 的 Unit 文件：

```bash
~]# systemctl cat wg-quick@.service
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

从 `[Service]` 部分可以看到，启动和停止服务，就是通过 wg-quick 命令实现的。

所以想要在让 Peer 开机自启 WireGuard，只需要执行如下操作即可：

```bash
systemctl enable wg-quick@wg-company --now
systemctl enable wg-quick@client-company --now
systemctl enable wg-quick@company --now
```

# 查看信息

接口：

```bash
# 查看系统 VPN 接口信息
$ ip link show wg-company
# 查看 VPN 接口详细信息
$ wg show all
$ wg show wg-company
```

地址：

```bash
# 查看 VPN 接口地址
$ ip address show wg-company
```

路由

```bash
# 查看系统路由表
$ ip route show table main
$ ip route show table local
# 获取到特定 IP 的路由
$ ip route get 192.0.2.3
```
