---
title: Tailscale
linkTitle: Tailscale
weight: 1
---

WireGuard 相比于传统 VPN 的核心优势是没有 VPN 网关，所有节点之间都可以点对点（P2P）连接，也就是我之前提到的全互联模式（full mesh），效率更高，速度更快，成本更低。

WireGuard 目前最大的痛点就是上层应用的功能不够健全，因为 WireGuard 推崇的是 Unix 的哲学，WireGuard 本身只是一个内核级别的模块，只是一个数据平面，至于上层的更高级的功能（比如秘钥交换机制，UDP 打洞，ACL 等），需要通过用户空间的应用来实现。

所以为了基于 WireGuard 实现更完美的 VPN 工具，现在已经涌现出了很多项目在互相厮杀。Netmaker 通过可视化界面来配置 WireGuard 的全互联模式，它支持 UDP 打洞、多租户等各种高端功能，几乎适配所有平台，非常强大。然而现实世界是复杂的，无法保证所有的 NAT 都能打洞成功，且 Netmaker 目前还没有 fallback 机制，如果打洞失败，无法 fallback 改成走中继节点。Tailscale 在这一点上比 Netmaker 高明许多，它支持 fallback 机制，可以尽最大努力实现全互联模式，部分节点即使打洞不成功，也能通过中继节点在这个虚拟网络中畅通无阻。

# 概述

> 参考：
>
> - [GitHub 项目，tailscale/tailscale](https://github.com/tailscale/tailscale)
> - [官网](https://tailscale.com/)
> - [公众号 - 云原声实验室，Tailscal 开源版本让你的 WireGuard 直接起飞](https://mp.weixin.qq.com/s/Y3z5RzuapZc8jS0UuHLhBw)
> - [Tailscale 的开源方案](https://tailscale.com/opensource)

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

简而言之，可以将 Tailscale 看成是更为易用、功能更完善的 [WireGuard](/docs/4.数据通信/Protocol/Tunneling%20Protocol/WireGuard/WireGuard.md)。

Tailscale 是一款商业产品，但个人用户是可以白嫖的，个人用户在接入设备不超过 20 台的情况下是可以免费使用的（虽然有一些限制，比如子网网段无法自定义，且无法设置多个子网）。除 Windows 和 macOS 的图形应用程序外，其他 Tailscale 客户端的组件（包含 Android 客户端）是在 BSD 许可下以开源项目的形式开发的，你可以在他们的 GitHub 仓库找到各个操作系统的客户端源码。

对于大部份用户来说，白嫖 Tailscale 已经足够了，如果你有更高的需求，比如自定义网段，可以选择付费。

**我就不想付费行不行？行，不过得看 [Headscale](/docs/4.数据通信/Protocol/Tunneling%20Protocol/Tailscale/Headscale.md)**

# Tailscale 架构概述

- **Tailscale 控制台** # 官方称为 **Coordination server(协调服务器)**。管理 Tailscale 客户端，向 Tailscale 客户端下发规则。
  - 可以通过 [Headscale](/docs/4.数据通信/Protocol/Tunneling%20Protocol/Tailscale/Headscale.md) 开源实现
- **Tailscale 客户端** # 主要是 [tailscale CLI](/docs/4.数据通信/Protocol/Tunneling%20Protocol/Tailscale/tailscale%20CLI.md)。windows 也有调用 tailscale 命令行工具的守护进程以右下角小图标的形式存在
  - Tailscale 客户端通常分为两部分，一部分是处理数据包的主程序（平时说的 Tailscale 客户端就是指这个主程序）；一部分类似 CLI 用以控制主程序。
  - e.g. Linux 的 Tailscale 客户端由两个程序组成: tailscale 和 tailscaled，tailscale 是 CLI，tailscaled 是守护程序用以处理数据包的路由。有点类似 docker 与 dockerd
- **Tailscale DERP** # 当两个节点第一次连接以及两个节点直连失败时，会切换到通过 DERP 来连接。DERP 是 Tailscale 自研的协议，也是一个中继程序，用以代理两个节点的访问请求。
  - Notes: 可以自行搭建 [DERP](/docs/4.数据通信/Protocol/Tunneling%20Protocol/Tailscale/Tailscale%20DERP.md)

# Tailscale 部署

在 https://pkgs.tailscale.com/stable/ 可以看到 Tailscale 的所有平台的安装包

# Tailscale 关联文件与配置

**/var/lib/tailscale/** # 运行时数据保存路径

- **./tailscaled.state** # 文本格式的配置文件。通过 [tailscale CLI](/docs/4.数据通信/Protocol/Tunneling%20Protocol/Tailscale/tailscale%20CLI.md) 指定的参数将会保存到该文件中
- **./derpmap.cached.json** # 可用的 DERP 缓存配置

---

Windows

**C:/ProgramData/Tailscale/** # 运行时数据保存路径

- **./Logs/** # 程序运行日志
- 其他文件及其功能与 Linux 下的 /var/lib/tailscale/ 目录相同

# Tailscale 工作逻辑

> 参考：
>
> [Tailscale blog, How Tailscale works](https://tailscale.com/blog/how-tailscale-works)
> - [Tailscale How it works](/docs/4.数据通信/Protocol/Tunneling%20Protocol/Tailscale/Tailscale%20How%20it%20works.md)

Tailscale 的 **所有客户端之间的连接都是先选择 DERP 模式（中继模式），这意味着连接立即就能建立（优先级最低但 100% 能成功的模式），用户不用任何等待**。然后开始并行地进行路径发现，通常几秒钟之后，就能发现一条更优路径，然后将现有连接透明升级（upgrade）成点对点连接（直连）。

可以通过 `tailscale ping ${HOST}` 命令查看到目标 HOST 的路由路径是否经过 DERP

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/tailscale/tailscale-cli-ping-example-1.png)

# 路由表

Tailscale 会自动创建相关的路由表和 iptables 规则。路由表可通过以下命令查看：

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
100.64.0.2 dev tailscale0
100.100.100.100 dev tailscale0
```

一般都是 52 表

# 访问控制

> 参考：
>
> - [官方文档，管理访问 - 管理访问控制](https://tailscale.com/kb/1393/access-control)

[What is a tailnet](https://tailscale.com/kb/1136/tailnet/)

## ACL

Tailscale 使用 [Tailscale ACL](/docs/4.数据通信/Protocol/Tunneling%20Protocol/Tailscale/Tailscale%20ACL.md) 管理权限，后来改为 Grant

```json
{
    "acls": [
    {
        "action": "accept",
        "src": [ <list-of-sources> ], // These sources (devices or users)
        "dst": [ <destination>:<port> ], // can access these destination devices on their defined ports
    }
  ]
}
```

## Grant


