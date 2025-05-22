---
title: Headscale
linkTitle: Headscale
weight: 2
---

# 概述

> 参考：
>
> - [GitHub 项目，juanfont/headscale](https://github.com/juanfont/headscale)
> - [公众号 - 云原声实验室，Tailscal 开源版本让你的 WireGuard 直接起飞](https://mp.weixin.qq.com/s/Y3z5RzuapZc8jS0UuHLhBw)
> - [馆长博客，headscale 搭建和应用场景](https://zhangguanzhang.github.io/2024/07/25/headscale/)

Tailscale 的控制服务器是不开源的，而且对免费用户有诸多限制，这是人家的摇钱树，可以理解。好在目前有一款开源的实现叫 Headscale，这也是唯一的一款，希望能发展壮大。

Headscale 由欧洲航天局的 Juan Font 使用 Go 语言开发，在 BSD 许可下发布，实现了 Tailscale 控制服务器的所有主要功能，可以部署在企业内部，没有任何设备数量的限制，且所有的网络流量都由自己控制。

目前 Headscale 还没有可视化界面，期待后续更新吧。

# Headscale 部署

> 理论上来说只要你的 Headscale 服务可以暴露到公网出口就行，但最好不要有 NAT，所以推荐将 Headscale 部署在有公网 IP 的云主机上。

## 安装程序与应用配置

**准备一些环境变量**

```bash
export HeadscaleVersion="0.22.3"
export HeadscaleArch="amd64"
# 各个 Tailscale 节点与 Headscale 通信的 IP
export HeadscaleAddr="https://X.X.X.X:YYY"
```

**准备 Headscale 相关文件及目录**。从 [GitHub 仓库的 Release 页面](https://github.com/juanfont/headscale/releases)下载最新版的二进制文件。

```bash
wget --output-document=/usr/local/bin/headscale \
  https://github.com/juanfont/headscale/releases/download/v${HeadscaleVersion}/headscale_${HeadscaleVersion}_linux_${HeadscaleArch}

chmod +x /usr/local/bin/headscale
```

创建相关目录及文件

```bash
mkdir -p /etc/headscale
mkdir -p /var/lib/headscale
touch /var/lib/headscale/db.sqlite
```

创建 headscale 用户，并修改相关文件权限：

```bash
useradd headscale -d /home/headscale -m
chown -R headscale:headscale /var/lib/headscale
```

**创建 Headscale 配置文件**

下载配置文件

```bash
wget https://raw.githubusercontent.com/juanfont/headscale/v${HeadscaleVersion}/config-example.yaml -O /etc/headscale/config.yaml
```

- 修改配置文件
  - **server_url** # 改为公网 IP 或域名。**如果是国内服务器，域名必须要备案**。我的域名无法备案，所以我就直接用公网 IP 了。
  - **magic_dns** # 如果暂时用不到 DNS 功能，该值设为 false
  - **prefixes** # 可自定义私有网段

**创建 Systemd Unit 文件**

```bash
curl -o /etc/systemd/system/headscale.service -LO https://github.com/juanfont/headscale/raw/refs/heads/main/packaging/systemd/headscale.service
```

**启动 Headscale 服务**

```bash
systemctl daemon-reload
systemctl enable headscale --now
```

## 创建 Headscale User

> [!Note] 老版本是创建 Namesapce
>
> - Tailscale 中有一个概念叫 tailnet，可以理解成租户， Tailscale 与 Tailscale 之间是相互隔离的，具体看参考 Tailscale 的官方文档：[What is a tailnet](https://tailscale.com/kb/1136/tailnet/)。
> - Headscale 也有类似的实现叫 namespace，即命名空间。Namespace 是一个实体拥有的机器的逻辑组，这个实体对于 Tailscale 来说，通常代表一个用户。
> - 我们需要先创建一个 namespace，以便后续客户端接入，例如：

```bash
# 这部分命令不用再执行了
~]# headscale namespaces create desistdaydream
Namespace created
~]# headscale namespaces list
ID | Name      | Created
1  | desistdaydream | 2022-03-24 04:23:04
```

注意：

- 从 v0.15.0 开始，Namespace 之间的边界已经被移除了，所有节点默认可以通信，如果想要限制节点之间的访问，可以使用 [ACL](https://github.com/juanfont/headscale/blob/v0.15.0/docs/acls.md)。在配置文件中只用 `acl_policy_path` 字段指定 ACL 配置文件路径，文件配置方式详见: https://tailscale.com/kb/1018/acls/

```bash
~]# headscale users create desistdaydream
~]# headscale user list
An updated version of Headscale has been found (0.23.0-alpha5 vs. your current v0.22.3). Check it out https://github.com/juanfont/headscale/releases
ID | Name           | Created
1  | desistdaydream | 2024-03-20 14:29:47
```

# Tailscale 客户端部署与接入 Headscale

https://headscale.net/stable/about/clients/

Headscale 只是实现了 Tailscale 的控制台，想要接入，依然需要使用 Tailscale 客户端。Headscale 默认支持最新 10 个版本的 Tailscale

| OS      | 是否支持 Headscale                                                                                                                |
| ------- | ----------------------------------------------------------------------------------------------------------------------------- |
| Linux   | Yes                                                                                                                           |
| OpenBSD | Yes                                                                                                                           |
| FreeBSD | Yes                                                                                                                           |
| macOS   | Yes                                                                                                                           |
| Windows | Yes (see [docs](https://headscale.net/stable/usage/connect/windows/) and `/windows` on your headscale for more information)   |
| Android | Yes (see [docs](https://headscale.net/stable/usage/connect/android/) for more information)                                    |
| iOS     | Yes (see [docs](https://headscale.net/stable/usage/connect/apple/#macos) and `/apple` on your headscale for more information) |

想要让 Tailscale 客户端接入 Headscale，大体分为两个部分

- 下载并配置 Tailscale 客户端，获取加入节点的指令
- 在 Headscale 上执行加入节点的指令

## 部署 Tailscale 客户端

### Linux

在 Tailscale 部署的节点准备环境变量

```bash
export TailscaleVersion="1.78.1"
export TailscaleArch="amd64"
export HeadscaleAddr="https://X.X.X.X:YYY"
```

Tailscale 官方提供了各种 Linux 发行版的软件包，但在国内由于网络原因，这些软件源基本用不了。所以我们可以在[这里](https://pkgs.tailscale.com/stable/#static)可以找到所有 Tailscale 的二进制文件，下载，并解压

```bash
wget https://pkgs.tailscale.com/stable/tailscale_${TailscaleVersion}_${TailscaleArch}.tgz
tar -zxvf tailscale_${TailscaleVersion}_${TailscaleArch}.tgz
```

这个包里包含如下文件：

```text
tailscale_${TailscaleVersion}_${TailscaleArch}/tailscale
tailscale_${TailscaleVersion}_${TailscaleArch}/tailscaled
tailscale_${TailscaleVersion}_${TailscaleArch}/systemd/
tailscale_${TailscaleVersion}_${TailscaleArch}/systemd/tailscaled.defaults
tailscale_${TailscaleVersion}_${TailscaleArch}/systemd/tailscaled.service
```

将文件复制到对应路径下(这里的路径其实就是通过各种软件源安装的路径)：

```bash
cp tailscale_${TailscaleVersion}_${TailscaleArch}/tailscaled /usr/sbin/tailscaled
cp tailscale_${TailscaleVersion}_${TailscaleArch}/tailscale /usr/bin/tailscale
cp tailscale_${TailscaleVersion}_${TailscaleArch}/systemd/tailscaled.service /lib/systemd/system/tailscaled.service
cp tailscale_${TailscaleVersion}_${TailscaleArch}/systemd/tailscaled.defaults /etc/default/tailscaled
```

启动 tailscaled.service 并设置开机自启：

```bash
systemctl enable tailscaled --now
```

Tailscale 接入 Headscale：

> 这里推荐将 DNS 功能关闭，因为它会覆盖系统的默认 DNS。关闭接收路由功能，有需要再打开，这样可以让机器只能访问到各 Tailscale

```bash
tailscale up --login-server=${HeadscaleAddr} --accept-routes=false --accept-dns=false
```

执行完上面的命令后，会出现下面的信息：

```bash
Warning: IP forwarding is disabled, subnet routing/exit nodes will not work.
See https://tailscale.com/kb/1104/enable-ip-forwarding/

To authenticate, visit:

 http://X.X.X.X:8080/register?key=30e9c9c952e2d66680b9904eb861e24a595e80c0839e3541142edb56c0d43e16

Success.
```

在浏览器中打开该链接，就会出现如下的界面：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/headscale/1648097896734-2252f545-7f43-46fb-ad3a-1e1c85ce3d08.png)

最后，根据 [在 Headscale 中添加节点](#Headscale%20中添加节点) 部分的文档，将节点接入到 Headscale 中。

### Windows

https://headscale.net/windows-client/

Windows Tailscale 客户端想要使用 Headscale 作为控制服务器，只需在浏览器中打开 `${HeadscaleAddr}/windows`，根据页面提示，本质上是执行下面这些操作

- 添加注册表信息（两种方式）（在 `HKEY_LOCAL_MACHINE\SOFTWARE\Tailscale IPN` 位置生成信息）
  - 点击页面中的 `Windows registry file`，下载注册表文件，并运行
  - 或者执行下面的 PowerShell 命令添加注册表信息

```powershell
$headscale_server="https://DOMAIN:PORT"
New-Item -Path "HKLM:\SOFTWARE\Tailscale IPN"
New-ItemProperty -Path 'HKLM:\Software\Tailscale IPN' -Name UnattendedMode -PropertyType String -Value always
New-ItemProperty -Path 'HKLM:\Software\Tailscale IPN' -Name LoginURL -PropertyType String -Value ${headscale_server}
```

- 在[这里](https://pkgs.tailscale.com/stable/#windows)下载 Windows 版的 Tailscale 客户端并安装
- 在 Powershell 执行

```bash
tailscale up --login-server=http://${headscale_server} --accept-routes=true --accept-dns=false
```

- 右键点击任务栏中的 Tailscale 图标，再点击 `Log in` 获取接入 Headscale 的命令
- ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/headscale/1648104111282-99d562e1-d7d9-4ea5-9943-f16861efe87e.png)
- 此时会自动在浏览器中出现接入 Headscale 的页面，记录下注册命令，去 Headscale 所在设备上执行命令添加节点。
  - 注册命令示例：

```bash
headscale nodes register --user USERNAME --key nodekey:75b424a753067b906fee373411743bf34264a9c40a4f7798836bf28bb24d2467
```

- 根据 [在 Headscale 中添加节点](#Headscale%20中添加节点) 部分的文档，使用注册命令将节点接入到 Headscale 中。

### 其他 Linux 发行版

除了常规的 Linux 发行版之外，还有一些特殊场景的 Linux 发行版，比如 OpenWrt、威联通（QNAP）、群晖等，这些发行版的安装方法已经有人写好了，这里就不详细描述了，相关的 GitHub 仓库：

- OpenWrt：<https://github.com/adyanth/openwrt-tailscale-enabler>
- 群晖：<https://github.com/tailscale/tailscale-synology>
- 威联通：<https://github.com/ivokub/tailscale-qpkg>

## Headscale 中添加节点

将其中的命令复制粘贴到 Headscale 所在机器的终端中，并将 NAMESPACE 替换为前面所创建的 namespace。

```bash
export HeadscaleUser="desistdaydream"
# 上面例子中的 Linux 客户端
headscale -n ${HeadscaleUser} nodes register --key 30e9c9c952e2d66680b9904eb861e24a595e80c0839e3541142edb56c0d43e16
# 上面例子中的 Windows 客户端
headscale -n ${HeadscaleUser} nodes register --key 105363c37b5449b85bb3e4107b6f6bbd3a2bb379dcf731bc98f979584740644a
```

注册成功，查看注册的节点：

> 这里可以看到，已经注册的节点将会分配一个 IP，这里是 100.64.0.1，其他注册的节点可以通过这个 IP 访问该节点。

```bash
~]# headscale nodes list
ID | Hostname        | Name            | MachineKey | NodeKey | User           | IP addresses                  | Ephemeral | Last seen           | Expiration          | Online | Expired
1  | HOME-WUJI       | home-wuji       | [fqHlf]    | [dbQkp] | desistdaydream | 100.64.0.1, fd7a:115c:a1e0::1 | false     | 2024-03-20 15:55:30 | 0001-01-01 00:00:00 | online | no
2  | DESKTOP-R02G6RP | desktop-r02g6rp | [zMy/C]    | [Utjz0] | desistdaydream | 100.64.0.2, fd7a:115c:a1e0::2 | false     | 2024-03-20 15:55:36 | 0001-01-01 00:00:00 | online | no
```

## 检查 Tailscale

回到 Tailscale 客户端所在的 Linux 主机，可以看到 Tailscale 会自动创建相关的路由表和 iptables 规则。路由表可通过以下命令查看：

```shell
~]# ip rule show
0: from all lookup local
5210: from all fwmark 0x80000 lookup main
5230: from all fwmark 0x80000 lookup default
5250: from all fwmark 0x80000 unreachable
5270: from all lookup 52
32766: from all lookup main
32767: from all lookup default
~]# ip route show table 52
100.64.0.2 dev tailscale0 # 这就是那个 Windows 节点
100.100.100.100 dev tailscale0
```

查看 iptables 规则：

```bash
~]# iptables -S
-P INPUT ACCEPT
-P FORWARD ACCEPT
-P OUTPUT ACCEPT
-N ts-forward
-N ts-input
-A INPUT -j ts-input
-A FORWARD -j ts-forward
-A FORWARD -i company -j ACCEPT
-A FORWARD -o company -j ACCEPT
-A ts-forward -i tailscale0 -j MARK --set-xmark 0x40000/0xffffffff
-A ts-forward -m mark --mark 0x40000 -j ACCEPT
-A ts-forward -s 100.64.0.0/10 -o tailscale0 -j DROP
-A ts-forward -o tailscale0 -j ACCEPT
-A ts-input -s 100.64.0.1/32 -i lo -j ACCEPT
-A ts-input -s 100.115.92.0/23 ! -i tailscale0 -j RETURN
-A ts-input -s 100.64.0.0/10 ! -i tailscale0 -j DROP
~]# iptables -t nat -S
-P PREROUTING ACCEPT
-P INPUT ACCEPT
-P OUTPUT ACCEPT
-P POSTROUTING ACCEPT
-N ts-postrouting
-A POSTROUTING -j ts-postrouting
-A POSTROUTING -o ens192 -j MASQUERADE
-A ts-postrouting -m mark --mark 0x40000 -j MASQUERADE
```

# 利用 Tailscale 打通局域网

到目前为止我们只是打造了一个点对点的 Mesh 网络，各个节点之间都可以通过 WireGuard 的私有网络 IP 进行直连(就是部署时默认使用的 100.64.0.0/10 网段中的 IP)。

但我们可以更大胆一点，还记得我在文章开头提到的访问家庭内网的资源吗？我们可以通过适当的配置让每个节点都能访问其他节点的局域网 IP。这个使用场景就比较多了，你可以直接访问家庭内网的 NAS，或者内网的任何一个服务，**更高级的玩家可以使用这个方法来访问云上 Kubernetes 集群的 Pod IP 和 Service IP。**

假设你的家庭内网有一台 Linux 主机（比如 OpenWrt）安装了 Tailscale 客户端，我们希望其他 Tailscale 客户端可以直接通过家中的局域网 IP（例如 **192.168.100.0/24**） 访问家庭内网的任何一台设备。

配置方法很简单，首先需要设置 IPv4 与 IPv6 路由转发：

```bash
tee /etc/sysctl.d/ipforwarding.conf > /dev/null <<EOF
net.ipv4.ip_forward = 1
net.ipv6.conf.all.forwarding = 1
EOF

sysctl -p /etc/sysctl.d/ipforwarding.conf
```

客户端修改注册节点的命令，在原来命令的基础上加上参数 `--advertise-routes=192.168.100.0/24`。多个 CIDR 以 `,` 分割

```bash
tailscale up --login-server=${HeadscaleAddr} --accept-routes=true --accept-dns=false  --advertise-routes=172.38.40.0/24,192.168.88.0/24
```

或通过 tailscale set 命令直接增加路由

```bash
tailscale set --advertise-routes=172.38.40.0/24,192.168.88.0/24
```

在 Headscale 端查看路由，可以看到相关路由是关闭的。

```bash
~]# headscale nodes list
ID | Hostname        | Name            | MachineKey | NodeKey | User           | IP addresses                  | Ephemeral | Last seen           | Expiration          | Online | Expired
1  | HOME-WUJI       | home-wuji       | [fqHlf]    | [dbQkp] | desistdaydream | 100.64.0.1, fd7a:115c:a1e0::1 | false     | 2024-03-20 15:55:30 | 0001-01-01 00:00:00 | online | no
2  | DESKTOP-R02G6RP | desktop-r02g6rp | [zMy/C]    | [Utjz0] | desistdaydream | 100.64.0.2, fd7a:115c:a1e0::2 | false     | 2024-03-20 15:55:36 | 0001-01-01 00:00:00 | online | no
~]# headscale routes list
ID | Machine         | Prefix           | Advertised | Enabled  | Primary
1  | desktop-r02g6rp | 172.38.40.0/24   | true       | false    | false
2  | desktop-r02g6rp | 192.168.88.0/24  | true       | false    | false
```

开启路由（0.26.0+ 版本）：

```bash
$ headscale nodes list-routes
ID | Hostname           | Approved | Available       | Serving (Primary)
1  | ts-head-ruqsg8     |          | 0.0.0.0/0, ::/0 |
2  | ts-unstable-fq7ob4 |          | 0.0.0.0/0, ::/0 |

$ headscale nodes approve-routes --identifier 1 --routes 0.0.0.0/0,::/0
Node updated

$ headscale nodes list-routes
ID | Hostname           | Approved        | Available       | Serving (Primary)
1  | ts-head-ruqsg8     | 0.0.0.0/0, ::/0 | 0.0.0.0/0, ::/0 | 0.0.0.0/0, ::/0
2  | ts-unstable-fq7ob4 |                 | 0.0.0.0/0, ::/0 |
```

开启路由：

```bash
~]# headscale routes enable -r 1
~]# headscale routes enable -r 2
ID | Machine         | Prefix           | Advertised | Enabled | Primary
1  | desktop-r02g6rp | 172.38.40.0/24   | true       | true    | true
2  | desktop-r02g6rp | 192.168.88.0/24  | true       | true    | true

```

# Headscale 内嵌 DERPer

https://github.com/juanfont/headscale/issues/1326

https://github.com/juanfont/headscale/pull/388

修改配置文件中的如下几个字段

```yaml
derp:
  server:
    enable: true
tls_cert_path: "/PATH/TO/FILE"
tls_key_path: "/PATH/TO/FILE"
```

stun 能力监听 3478 端口，且可以通过 https 访问

# Headscale 关联文件与配置

**config.{yaml,json}** # Headscale 配置文件。Headscale 启动时从 /etc/headscale/、~/.headscale/、当前工作目录 这三个地方查找文件以加载配置。

- [GitHub 项目，juanfont/headscale - config-example.yaml](https://github.com/juanfont/headscale/blob/main/config-example.yaml) 中是配置文件的示例

**/var/lib/headscale/** # Headscale 运行时数据保存路径。包括 数据库文件、证书 等

- **./db.sqlite** # Headscale 运行后数据持久化的 Sqlite3 存储
- **./private.key** # 用于加密 Headscale 和 Tailscale 客户端之间流量的私钥。如果私钥文件丢失，将自动生成。

[这里](https://github.com/juanfont/headscale/blob/main/config-example.yaml)是配置文件示例

## 配置详解

**tls_cert_path**(STRING) # 证书路径

**tls_key_path**(STRING) # 证书私钥路径

**randomize_client_port**(BOOLEAN) # 启用此选项使设备使用随机端口来传输 WireGuard 流量，而不是默认 41641 端口。`默认值: false`。此选项旨在作为某些有缺陷的防火墙设备的解决方法，详见: https://tailscale.com/kb/1181/firewalls

### derp

**server**(OBJECT)

- **enable**(BOOLEAN) # 是否开启 Headscale 的内嵌 DERP。`默认值: false`

**url**(\[]STRING) # 下发给 tailscale 的 DERP 节点。`默认值: https://controlplane.tailscale.com/derpmap/default`

- <font color="#ff0000">这些默认值是 Tailscale 提供的一些公共的 DERP 节点，全球都用，个人建议直接关了，用自己的 Headscale DERP</font>

**paths**(\[]STRING) # 与 URL 类似。不过是以文件形式定义要使用的 DERP。 `默认值: 空`。文件格式详见: https://tailscale.com/kb/1118/custom-derp-servers
