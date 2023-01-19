---
title: Headscale(Tailscale开源版本)
---

# 概述

> 参考：
> - [公众号，云原声实验室-Tailscal 开源版本让你的 WireGuard 直接起飞](https://mp.weixin.qq.com/s/Y3z5RzuapZc8jS0UuHLhBw)
> - [GitHub 项目，juanfont/headscale](https://github.com/juanfont/headscale)

目前国家工信部在大力推动三大运营商发展 IPv6，对家用宽带而言，可以使用的 IPv4 公网 IP 会越来越少。有部分地区即使拿到了公网 IPv4 地址，也是个大内网地址，根本不是真正的公网 IP，访问家庭内网的资源将会变得越来越困难。

部分小伙伴可能会选择使用 frp 等针对特定协议和端口的内网穿透方案，但这种方案还是不够酸爽，无法访问家庭内网任意设备的任意端口。更佳的选择还是通过 VPN 来组建大内网。至于该选择哪种 VPN，毫无疑问肯定是 WireGuard，WireGuard 就是 VPN 的未来。

WireGuard 相比于传统 VPN 的核心优势是没有 VPN 网关，所有节点之间都可以点对点（P2P）连接，也就是我之前提到的全互联模式（full mesh），效率更高，速度更快，成本更低。

WireGuard 目前最大的痛点就是上层应用的功能不够健全，因为 WireGuard 推崇的是 Unix 的哲学，WireGuard 本身只是一个内核级别的模块，只是一个数据平面，至于上层的更高级的功能（比如秘钥交换机制，UDP 打洞，ACL 等），需要通过用户空间的应用来实现。

所以为了基于 WireGuard 实现更完美的 VPN 工具，现在已经涌现出了很多项目在互相厮杀。笔者前段时间一直在推崇 Netmaker，它通过可视化界面来配置 WireGuard 的全互联模式，它支持 UDP 打洞、多租户等各种高端功能，几乎适配所有平台，非常强大。然而现实世界是复杂的，无法保证所有的 NAT 都能打洞成功，且 Netmaker 目前还没有 fallback 机制，如果打洞失败，无法 fallback 改成走中继节点。Tailscale 在这一点上比 Netmaker 高明许多，它支持 fallback 机制，可以尽最大努力实现全互联模式，部分节点即使打洞不成功，也能通过中继节点在这个虚拟网络中畅通无阻。

## Tailscale 是什么

Tailscale 是一种基于 WireGuard 的虚拟组网工具，和 Netmaker 类似，**最大的区别在于 Tailscale 是在用户态实现了 WireGuard 协议，而 Netmaker 直接使用了内核态的 WireGuard**。所以 Tailscale 相比于内核态 WireGuard 性能会有所损失，但与 OpenVPN 之流相比还是能甩好几十条街的，Tailscale 虽然在性能上做了些许取舍，但在功能和易用性上绝对是完爆其他工具：

- 开箱即用
  - 无需配置防火墙
  - 没有额外的配置
- 高安全性/私密性
  - 自动密钥轮换
  - 点对点连接
  - 支持用户审查端到端的访问记录
- 在原有的 ICE、STUN 等 UDP 协议外，实现了 DERP TCP 协议来实现 NAT 穿透
- 基于公网的控制服务器下发 ACL 和配置，实现节点动态更新
- 通过第三方（如 Google） SSO 服务生成用户和私钥，实现身份认证

简而言之，我们可以将 Tailscale 看成是更为易用、功能更完善的 WireGuard。

光有这些还不够，作为一个白嫖党，咱更关心的是**免费**与**开源**。

Tailscale 是一款商业产品，但个人用户是可以白嫖的，个人用户在接入设备不超过 20 台的情况下是可以免费使用的（虽然有一些限制，比如子网网段无法自定义，且无法设置多个子网）。除 Windows 和 macOS 的图形应用程序外，其他 Tailscale 客户端的组件（包含 Android 客户端）是在 BSD 许可下以开源项目的形式开发的，你可以在他们的 GitHub 仓库\[3]找到各个操作系统的客户端源码。

对于大部份用户来说，白嫖 Tailscale 已经足够了，如果你有更高的需求，比如自定义网段，可以选择付费。

**我就不想付费行不行？行，不过得往下看。**

## Headscale 是什么

Tailscale 的控制服务器是不开源的，而且对免费用户有诸多限制，这是人家的摇钱树，可以理解。好在目前有一款开源的实现叫 Headscale，这也是唯一的一款，希望能发展壮大。

Headscale 由欧洲航天局的 Juan Font 使用 Go 语言开发，在 BSD 许可下发布，实现了 Tailscale 控制服务器的所有主要功能，可以部署在企业内部，没有任何设备数量的限制，且所有的网络流量都由自己控制。

目前 Headscale 还没有可视化界面，期待后续更新吧。

# Headscale 部署

Headscale 部署很简单，推荐直接在 Linux 主机上安装。

> 理论上来说只要你的 Headscale 服务可以暴露到公网出口就行，但最好不要有 NAT，所以推荐将 Headscale 部署在有公网 IP 的云主机上。

## 准备一些环境变量

```bash
export HeadscaleVersion="0.15.0"
export HeadscaleArch="amd64"
# Headscale 用于与各个节点通信的 IP
export HeadscaleIP="X.X.X.X"
```

## 准备 Headscale 相关文件及目录

从 [GitHub 仓库的 Release 页面](https://github.com/juanfont/headscale/releases)下载最新版的二进制文件。

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

### 创建 Headscale 配置文件

有两种方式

- 下载文件后修改内容
- 直接按照自己的要求创建

下载配置文件

```bash
wget https://github.com/juanfont/headscale/raw/main/config-example.yaml -O /etc/headscale/config.yaml
```

- 修改配置文件
  - **server_url** # 改为公网 IP 或域名。**如果是国内服务器，域名必须要备案**。我的域名无法备案，所以我就直接用公网 IP 了。
  - **magic_dns** # 如果暂时用不到 DNS 功能，该值设为 false
  - **unix_socket** # unix_socket: /var/run/headscale/headscale.sock
  - **ip_prefixes** # 可自定义私有网段

直接创建配置

```yaml
tee /etc/headscale/config.yaml > /dev/null <<EOF
server_url: http://${HeadscaleIP}:8080
listen_addr: 0.0.0.0:8080
metrics_listen_addr: 127.0.0.1:9090
grpc_listen_addr: 0.0.0.0:50443
grpc_allow_insecure: false
private_key_path: /var/lib/headscale/private.key
ip_prefixes:
  - fd7a:115c:a1e0::/48
  - 100.64.0.0/10
derp:
  server:
    enabled: false
    region_id: 999
    region_code: "headscale"
    region_name: "Headscale Embedded DERP"
    stun_listen_addr: "0.0.0.0:3478"
  urls:
    - https://controlplane.tailscale.com/derpmap/default
  paths: []
  auto_update_enabled: true
  update_frequency: 24h
disable_check_updates: false
ephemeral_node_inactivity_timeout: 30m
db_type: sqlite3
db_path: /var/lib/headscale/db.sqlite
acme_url: https://acme-v02.api.letsencrypt.org/directory
acme_email: ""
tls_letsencrypt_hostname: ""
tls_client_auth_mode: relaxed
tls_letsencrypt_cache_dir: /var/lib/headscale/cache
tls_letsencrypt_challenge_type: HTTP-01
tls_letsencrypt_listen: ":http"
tls_cert_path: ""
tls_key_path: ""
log_level: info
acl_policy_path: ""
dns_config:
  nameservers:
    - 1.1.1.1
  domains: []
  magic_dns: true
  base_domain: example.com
unix_socket: /var/run/headscale/headscale.sock
unix_socket_permission: "0770"
EOF
```

## 创建 Systemd Unit 文件

```bash
tee /etc/systemd/system/headscale.service > /dev/null <<EOF
[Unit]
Description=headscale controller
After=syslog.target
After=network.target

[Service]
Type=simple
User=headscale
Group=headscale
ExecStart=/usr/local/bin/headscale serve
Restart=always
RestartSec=5

# Optional security enhancements
NoNewPrivileges=yes
PrivateTmp=yes
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=/var/lib/headscale /var/run/headscale
AmbientCapabilities=CAP_NET_BIND_SERVICE
RuntimeDirectory=headscale

[Install]
WantedBy=multi-user.target
EOF
```

## 启动 Headscale 服务

```bash
systemctl daemon-reload
systemctl enable --now headscale
```

## 创建 Headscale Namespace

Tailscale 中有一个概念叫 tailnet，你可以理解成租户， Tailscale 与 Tailscale 之间是相互隔离的，具体看参考 Tailscale 的官方文档：[What is a tailnet](https://tailscale.com/kb/1136/tailnet/)。

Headscale 也有类似的实现叫 namespace，即命名空间。Namespace 是一个实体拥有的机器的逻辑组，这个实体对于 Tailscale 来说，通常代表一个用户。

我们需要先创建一个 namespace，以便后续客户端接入，例如：

```bash
~]# headscale namespaces create lichenhao
Namespace created
~]# headscale namespaces list
ID | Name      | Created
1  | lichenhao | 2022-03-24 04:23:04
```

注意：

- 从 v0.15.0 开始，Namespace 之间的边界已经被移除了，所有节点默认可以通信，如果想要限制节点之间的访问，可以使用 [ACL](https://github.com/juanfont/headscale/blob/v0.15.0/docs/acls.md)。在配置文件中只用 `acl_policy_path` 字段指定 ACL 配置文件路径，文件配置方式详见：<https://tailscale.com/kb/1018/acls/>

# Headscale 关联文件与配置

**/etc/headscale/config.yaml** # Headscale 运行时配置文件
**/var/lib/headscale/\* **# Headscale 运行时数据目录。包括 数据库文件、证书 等

- **./db.sqlite** # Headscale 使用 sqlite 作为数据库

[这里](https://github.com/juanfont/headscale/blob/main/config-example.yaml)是配置文件示例

# Tailscale 客户端部署与接入 Headscale

Headscale 只是实现了 Tailscale 的控制台，想要接入，依然需要使用 Tailscale 客户端。

目前除了 iOS 客户端，其他平台的客户端都有办法自定义 Tailscale 的控制服务器。

| OS      | 是否支持 Headscale              |
| ------- | ------------------------------- |
| Linux   | Yes                             |
| OpenBSD | Yes                             |
| FreeBSD | Yes                             |
| macOS   | Yes                             |
| Windows | Yes 参考 Windows 客户端文档\[6] |
| Android | 需要自己编译客户端\[7]          |
| iOS     | 暂不支持                        |

想要让 Tailscale 客户端接入 Headscale，大体分为两个部分

- 下载并配置 Tailscale 客户端，获取加入节点的指令
- 在 Headscale 上执行加入节点的指令

## 部署 Tailscale 客户端

### Linux

在 Tailscale 部署的节点准备环境变量

```bash
export TailscaleVersion="1.26.1"
export TailscaleArch="amd64"
export HeadscaleIP="X.X.X.X"
```

Tailscale 官方提供了各种 Linux 发行版的软件包，但在国内由于网络原因，这些软件源基本用不了。所以我们可以在[这里](https://pkgs.tailscale.com/stable/#static)可以找到所有 Tailscale 的二进制文件，下载，并解压

```bash
wget https://pkgs.tailscale.com/stable/tailscale_${TailscaleVersion}_${TailscaleArch}.tgz
tar -zxvf tailscale_${TailscaleVersion}_${TailscaleArch}.tgz
```

这个包里包含如下文件：

```bash
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

> 这里推荐将 DNS 功能关闭，因为它会覆盖系统的默认 DNS。

```bash
tailscale up --login-server=http://${HeadscaleIP}:8080 --accept-routes=true --accept-dns=false
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
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/1648097896734-2252f545-7f43-46fb-ad3a-1e1c85ce3d08.png)
最后，根据 [在 Headscale 中添加节点](#LWIp8) 部分的文档，将节点接入到 Headscale 中。

### macOS

macOS 客户端的安装相对来说就简单多了，只需要在应用商店安装 APP 即可，前提是你**需要一个美区 ID**。。。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/640)

安装完成后还需要做一些骚操作，才能让 Tailscale 使用 Headscale 作为控制服务器。当然，Headscale 已经给我们提供了详细的操作步骤，你只需要在浏览器中打开 URL：`http://<HEADSCALE_PUB_IP>:8080/apple`，便会出现如下的界面：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/640)

你只需要按照图中所述的步骤操作即可，本文就不再赘述了。

修改完成后重启 Tailscale 客户端，在 macOS 顶部状态栏中找到 Tailscale 并点击，然后再点击 `Log in`。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/640)

然后立马就会跳转到浏览器并打开一个页面。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/640)

接下来与之前 Linux 客户端相同，回到 Headscale 所在的机器执行浏览器中的命令即可，注册成功：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/640)

回到 Headscale 所在主机，查看注册的节点：

\`$ headscale nodes list

ID | Name | NodeKey | Namespace | IP addresses | Ephemeral | Last seen | Onlin

e | Expired

1 | coredns | \[Ew3RB] | default | 10.1.0.1 | false | 2022-03-20 09:08:58 | onlin

e | no
2 | carsondemacbook-pro | \[k7bzX] | default   | 10.1.0.2     | false     | 2022-03-20 09:48:30 | online  | no

\`

回到 macOS，测试是否能 ping 通对端节点：

\`$ ping -c 2 10.1.0.1
PING 10.1.0.1 (10.1.0.1): 56 data bytes
64 bytes from 10.1.0.1: icmp_seq=0 ttl=64 time=37.025 ms
64 bytes from 10.1.0.1: icmp_seq=1 ttl=64 time=38.181 ms

\--- 10.1.0.1 ping statistics ---
2 packets transmitted, 2 packets received, 0.0% packet loss
round-trip min/avg/max/stddev = 37.025/37.603/38.181/0.578 ms

\`

也可以使用 Tailscale CLI 来测试：

`$ /Applications/Tailscale.app/Contents/MacOS/Tailscale ping 10.1.0.1 pong from coredns (10.1.0.1) via xxxx:41641 in 36ms`

---

如果你没有美区 ID，无法安装 App，可以直接使用命令行版本，通过 Homebrew 安装即可：

`$ brew install tailscale`

### Android

Android 客户端就比较麻烦了，需要自己修改源代码编译 App，具体可参考这个 issue\[9]。编译过程还是比较麻烦的，需要先修改源码，然后构建一个包含编译环境的 Docker 镜像，最后再通过该镜像启动容器编译 apk。

我知道很多人一看麻烦就不想搞了，这个问题不大，我送佛送到西，提供了一条龙服务，你只需 fork 我的 GitHub 仓库 tailscale-android\[10]：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/640)

然后在你的仓库中点击 **Settings** 标签，找到 **Secrets** 下拉框中的 Actions 选项：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/640)

选择 **New repository secret** 添加一个 secret 叫 `HEADSCALE_URL`，将你的 Headscale 服务公网地址填入其中：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/640)

**添加在这里的配置，将只对你可见，不用担心会泄露给他人。**

然后点击 **Actions** 标签，选择 **Release** Workflow。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/640)

你会看到一个 **Run workflow** 按钮，点击它，然后在下拉框中点击 **Run workflow**。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/640)

流水线就会开始执行，执行成功后就会在 Release 页面看到编译好的 apk。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/640)

接下来的事情就简单了，下载这个 apk 到你的 Android 手机上安装就好了。安装完成后打开 Tailscale App，选择 **Sign in with other**。
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/1648104924984-783c04d7-d504-43dd-b1ad-98dcf7231783.png)

然后就会跳出这个页面：
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/1648103496461-34386c19-5d46-4ddf-88f3-198d0e463d8a.png)

将其中的命令粘贴到 Headscale 所在主机的终端，将 **NAMESPACE** 替换为之前创建的 namespace，然后执行命令即可。注册成功后可将该页面关闭，回到 App 主页，效果如图：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/640)

---

回到之前的 GitHub 仓库，刚才我们是通过手动触发 Workflow 来编译 apk 的，有没有办法自动编译呢？**只要 Tailscale 官方仓库有更新，就立即触发 Workflow 开始编译。**

那当然是可以实现的，而且我已经实现了，仔细看 GitHub Actions 的编排文件：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/640)

红框圈出来的部分表示只要仓库的 `main` 分支有更新，便会触发 Workflow。**现在的问题是如何让 main 分支和上游官方仓库一致，一直保持在最新状态。**

这个问题使用第三方 Github App 就可以解决，这个 App 名字简单粗暴，就叫 Pull\[11]，它的作用非也很简单粗暴：保持你的 Fork 在最新状态。

Pull 的使用方法很简单：

1. 打开 Pull App\[12] 页面
2. 点击右上角绿色的 install 按钮

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/640)

3. 在选项页面，使用默认的 **All repositories** 即可（你也可以选择指定的仓库，比如 tailscale-android），然后点击绿色的 **install** 按钮：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/640)

简单三步，Pull App 就安装好了。接下来 Pull App 会每天定时帮你更新代码库，使你 fork 的代码始终是最新版的。

### Windows

Windows Tailscale 客户端想要使用 Headscale 作为控制服务器，只需在浏览器中打开 `http://${HeadscaleIP}>:8080/windows`，便会出现如下的界面：
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/1648103689211-3b019793-7262-4a6c-a668-127c0f01c284.png)

- 下载页面中的 `Windows registry file`，这是一个注册表文件，用来配置 Tailscale 客户端中控制服务器的地址，让其指向自己部署的 Headscale
- 下载完成后运行 `tailscale.reg` 文件以编辑注册表信息。
- 在[这里](https://tailscale.com/download/windows)下载 Windows 版的 Tailscale 客户端并安装
- 右键点击任务栏中的 Tailscale 图标，再点击 `Log in` 获取接入 Headscale 的命令

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/1648104111282-99d562e1-d7d9-4ea5-9943-f16861efe87e.png)

- 此时会自动在浏览器中出现接入 Headscale 的页面

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/1648104342439-c7a4a6ba-8690-4883-bf2f-c324c8336607.png)

### 其他 Linux 发行版

除了常规的 Linux 发行版之外，还有一些特殊场景的 Linux 发行版，比如 OpenWrt、威联通（QNAP）、群晖等，这些发行版的安装方法已经有人写好了，这里就不详细描述了，我只给出相关的 GitHub 仓库，大家如果自己有需求，直接去看相关仓库的文档即可。

- OpenWrt：<https://github.com/adyanth/openwrt-tailscale-enabler>
- 群晖：<https://github.com/tailscale/tailscale-synology>
- 威联通：<https://github.com/ivokub/tailscale-qpkg>

### iOS

Tailscale iOS 客户端源代码没有开源，目前还无法破解使其使用第三方控制服务器，遗憾~~

## Headscale 中添加节点

将其中的命令复制粘贴到 Headscale 所在机器的终端中，并将 NAMESPACE 替换为前面所创建的 namespace。

```bash
export HeadscaleNamespace="lichenhao"
# 上面例子中的 Linux 客户端
headscale -n ${HeadscaleNamespace} nodes register --key 30e9c9c952e2d66680b9904eb861e24a595e80c0839e3541142edb56c0d43e16
# 上面例子中的 Windows 客户端
headscale -n ${HeadscaleNamespace} nodes register --key 105363c37b5449b85bb3e4107b6f6bbd3a2bb379dcf731bc98f979584740644a
```

注册成功，查看注册的节点：

> 这里可以看到，已经注册的节点将会分配一个 IP，这里是 100.64.0.1，其他注册的节点可以通过这个 IP 访问该节点。

```bash
~]# headscale nodes  list
ID | Name                 | NodeKey | Namespace | IP addresses                  | Ephemeral | Last seen           | Online | Expired
1  | tj-test-oc-lichenhao | [Bo2d3] | lichenhao | fd7a:115c:a1e0::1, 100.64.0.1 | false     | 2022-03-24 06:48:46 | online | no
2  | home-desktop         | [VZuAp] | lichenhao | fd7a:115c:a1e0::2, 100.64.0.2 | false     | 2022-03-24 06:49:31 | online | no
```

## 检查 Tailscale

回到 Tailscale 客户端所在的 Linux 主机，可以看到 Tailscale 会自动创建相关的路由表和 iptables 规则。路由表可通过以下命令查看：

```shell
~]# ip rule show
0:	from all lookup local
5210:	from all fwmark 0x80000 lookup main
5230:	from all fwmark 0x80000 lookup default
5250:	from all fwmark 0x80000 unreachable
5270:	from all lookup 52
32766:	from all lookup main
32767:	from all lookup default
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

# 打通局域网

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

客户端修改注册节点的命令，在原来命令的基础上加上参数 `--advertise-routes=192.168.100.0/24`。

- 多个 CIDR 以 `,` 分割

```bash
tailscale up --login-server=http://${HeadscaleIP}:8080 --accept-routes=true --accept-dns=false  --advertise-routes=172.38.40.0/24,192.168.88.0/24
```

在 Headscale 端查看路由，可以看到相关路由是关闭的。

```bash
~]# headscale nodes list
ID | Name                 | NodeKey | Namespace | IP addresses                  | Ephemeral | Last seen           | Online | Expired
1  | tj-test-oc-lichenhao | [Bo2d3] | lichenhao | fd7a:115c:a1e0::1, 100.64.0.1 | false     | 2022-03-24 05:08:46 | online | no
2  | home-desktop         | [qZVTo] | lichenhao | fd7a:115c:a1e0::2, 100.64.0.2 | false     | 2022-03-24 05:09:16 | online | no
~]# headscale routes list -i 1
Route           | Enabled
172.38.40.0/24  | false
192.168.88.0/24 | false
```

开启路由：

```bash
~]# headscale routes enable -i 1 -r "172.38.40.0/24,192.168.88.0/24"
Route           | Enabled
172.38.40.0/24  | true
192.168.88.0/24 | true

```

其他非 Headscale 节点查看路由结果：
`$ ip route show table 52|grep "172.38.40.0/24" 172.38.40.0/24 dev tailscale0`

# 总结

目前从稳定性来看，Tailscale 比 Netmaker 略胜一筹，基本上不会像 Netmaker 一样时不时出现 ping 不通的情况，这取决于 Tailscale 在用户态对 NAT 穿透所做的种种优化，他们还专门写了一篇文章介绍 NAT 穿透的原理\[13]，中文版\[14]由国内的 eBPF 大佬赵亚楠翻译，墙裂推荐大家阅读。放一张图给大家感受一下：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/640)

本文给大家介绍了 Tailscale 和 Headscale，包括 Headscale 的安装部署和各个平台客户端的接入，以及如何打通各个节点所在的局域网。下篇文章将会给大家介绍如何让 Tailscale 使用自定义的 DERP Servers（也就是中继服务器），See you~~

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/af37c7a5-9f8e-4905-a639-0a377cb44ee6/640)

随着云原生对 IT   产业的重新洗牌，很多传统的技术在云原生的场景下已经不再适用，譬如备份和容灾。传统的备份容灾还停留在数据搬运的层次上，备份机制比较固化，以存储为核心，无法适应容器化的弹性、池化部署场景；而云原生的核心是服务本身，不再以存储为核心，用户需要更贴合容器场景的备份容灾能力，利用云原生的编排能力，实现备份容灾的高度自动化，同时灵活运用云原生的弹性能力按需付费，降低成本。

为了适应云原生场景，众多 Kubernetes 备份容灾产品开始涌现，比如 Veeam 推出的 Kasten K10 以及 VMware 推出的 Velero。青云科技也推出了 Kubernetes 备份容灾即服务产品，基于原生的 Kubernetes API，提供了可视化界面，能够覆盖云原生数据保护的绝大多数重要场景，而且能够跨集群、跨云服务商、跨存储区域，轻松实现基础设施间多地、按需的备份恢复。目前该服务已正式上线，提供了 1TB 的免费托管仓库，感兴趣的可以前往试用 👇‍

# 引用链接

\[1]

全互联模式（full mesh）: *https://fuckcloudnative.io/posts/wireguard-full-mesh/#1-全互联模式架构与配置*

\[2]

Netmaker: [_https://fuckcloudnative.io/posts/configure-a-mesh-network-with-netmaker/_](https://fuckcloudnative.io/posts/configure-a-mesh-network-with-netmaker/)

\[3]

GitHub 仓库: [_https://github.com/tailscale/_](https://github.com/tailscale/)

\[4]

Headscale: [_https://github.com/juanfont/headscale_](https://github.com/juanfont/headscale)

\[6]

Windows 客户端文档: [_https://github.com/juanfont/headscale/blob/main/docs/windows-client.md_](https://github.com/juanfont/headscale/blob/main/docs/windows-client.md)

\[7]

需要自己编译客户端: [_https://github.com/juanfont/headscale/issues/58#issuecomment-950386833_](https://github.com/juanfont/headscale/issues/58#issuecomment-950386833)

\[9]

这个 issue: [_https://github.com/juanfont/headscale/issues/58#issuecomment-950386833_](https://github.com/juanfont/headscale/issues/58#issuecomment-950386833)

\[10]

tailscale-android: [_https://github.com/yangchuansheng/tailscale-android_](https://github.com/yangchuansheng/tailscale-android)

\[11]

Pull: [_https://github.com/apps/pull_](https://github.com/apps/pull)

\[12]

Pull App: [_https://github.com/apps/pull_](https://github.com/apps/pull)

\[13]

NAT 穿透的原理: [_https://tailscale.com/blog/how-nat-traversal-works/_](https://tailscale.com/blog/how-nat-traversal-works/)

\[14]

中文版: [_https://arthurchiao.art/blog/how-nat-traversal-works-zh/_](https://arthurchiao.art/blog/how-nat-traversal-works-zh/)
