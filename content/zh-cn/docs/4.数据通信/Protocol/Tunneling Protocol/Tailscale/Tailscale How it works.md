---
title: "Tailscale: How it works"
source: "https://tailscale.com/blog/how-tailscale-works#nat-traversal"
author:
  - "[[Avery Pennarun]]"
published: 2020-03-21
created: 2025-05-23
description: "Understand the entire Tailscale system, how it works, how we built it, and its benefits compared to legacy VPNs. Use this article as a guide to quickly build your own Tailscale replacement."
tags:
  - "clippings"
---
[Blog 博客](https://tailscale.com/blog) | 三月 20, 2020

## How Tailscale works Tailscale 工作原理

People often ask us for an overview of how Tailscale works. We’ve been putting off answering that, because we kept changing it! But now things have started to settle down.  
人们经常要求我们概述 Tailscale 的工作原理。我们一直推迟回答这个问题，因为我们的架构在不断演进！但现在系统已逐渐趋于稳定。

Let’s go through the entire Tailscale system from bottom to top, the same way we built it (but skipping some zigzags we took along the way). With this information, you should be able to build your own Tailscale replacement… except you don’t have to, since our [node software is open source](https://github.com/tailscale/tailscale) and we have a [flexible free plan](https://tailscale.com/pricing).  
让我们从底层到顶层完整梳理 Tailscale 系统架构，这与我们的构建思路一致（但会略过期间走过的弯路）。了解这些信息后，您完全可以构建自己的 Tailscale 替代方案...不过其实无需重复造轮子，因为我们的节点软件是开源的，并且提供灵活的免费方案。

### The data plane: WireGuard®数据平面：WireGuard®

Our base layer is the increasingly popular and excellent open source [WireGuard](https://www.wireguard.com/) package (specifically the userspace Go variant, [wireguard-go](https://git.zx2c4.com/wireguard-go/about/)). WireGuard creates a set of extremely lightweight encrypted tunnels between your computer, VM, or container (which WireGuard calls an “endpoint” and we’ll call a “node”), and any other nodes in your network.  
我们的基础层是日益流行的优秀开源 WireGuard 组件（具体采用用户态 Go 语言实现 wireguard-go）。WireGuard 能在您的计算机、虚拟机或容器（WireGuard 称之为"端点"，我们称为"节点"）与网络中的其他节点之间，创建一组超轻量级的加密隧道。

Let’s clarify that: most of the time, VPN users, including WireGuard users, build a “hub and spoke” architecture, where each client device (say, your laptop) connects to a central “concentrator” or VPN gateway.  
需要澄清的是：大多数情况下，VPN 用户（包括 WireGuard 用户）构建的是"中心辐射型"架构，每个客户端设备（例如您的笔记本电脑）会连接到中央"集线器"或 VPN 网关。

#### Hub-and-spoke networks 中心辐射型网络

![Figure 1. A traditional hub-and-spoke VPN.](https://cdn.sanity.io/images/w77i7m8x/production/3cbc3fa27f0b798d3a0bc98f57829a9083dad769-1400x1080.svg?w=3840&q=75&fit=clip&auto=format)

Figure 1. A traditional hub-and-spoke VPN. 图 1. 传统中心辐射型 VPN 示意图。

This is the easiest way to set up WireGuard, because each node in the network needs to know the public key, public IP address, and port number of each other node it wants to connect directly to. If you wanted to fully connect 10 nodes, then that would be 9 peer nodes that each node has to know about, or 90 separate tunnel endpoints. In a hub-and-spoke, you just have one central hub node and 9 outliers (with “spokes” connecting them), which is much simpler. The hub node usually has a static IP address and a hole poked in its firewall so it’s easy for everyone to find. Then it can accept incoming connections from nodes at other IP addresses, even if those nodes are themselves behind a firewall, in the usual way that client-server Internet protocols do.  
这是配置 WireGuard 最简单的方式，因为网络中每个节点都需要知道其需要直连的其他节点的公钥、公网 IP 地址和端口号。如果要完全连接 10 个节点，意味着每个节点都需要了解其他 9 个对等节点，即总共需要维护 90 个独立的隧道端点。在中心辐射型架构中，只需 1 个中心枢纽节点和 9 个外围节点（通过"辐射"连接），结构更加简洁。中心节点通常配置静态 IP 地址并在防火墙开放端口，便于所有节点发现连接。随后它就能以常规客户端-服务器模式的互联网协议，接受来自其他动态 IP 地址节点的入站连接请求，即使这些节点自身处于防火墙保护之下。

Hub-and-spoke works well, but it has some downsides. First of all, most modern companies don’t have a single place they want to designate as a hub. They have multiple offices, multiple cloud datacenters or regions or VPCs, and so on. In traditional VPN setups, companies configure a single VPN concentrator, and then set up secondary tunnels (often using IPsec) between locations. So remote users arrive at the VPN concentrator in one place, then have their traffic forwarded to its final destination in another place.  
中心辐射型架构效果显著，但也存在一些缺陷。首先，大多数现代企业并不希望将单一地点指定为中心枢纽。它们往往拥有多个办公场所、多个云数据中心或区域、多个虚拟私有云（VPC）等等。在传统 VPN 方案中，企业会配置单个 VPN 集中器，然后在不同位置之间建立次级隧道（通常使用 IPsec 协议）。这意味着远程用户会首先连接到某一地点的 VPN 集中器，随后其流量会被转发至位于另一地点的最终目的地。

This traditional setup can be hard to scale. First of all, remote users might or might not be close to the VPN concentrator; if they’re far away, then they incur high latency connecting to it. Secondly, the datacenter they want to reach might not be close to the VPN concentrator either; if it’s far away, they incur high latency again. Imagine this case, with a worker in New York trying to reach a server in New York, through the company’s VPN concentrator or BeyondCorp proxy at their head office in San Francisco, not that I’m bitter:  
这种传统架构难以扩展。首先，远程用户与 VPN 集中器之间可能存在地理距离问题——若距离过远就会产生高延迟。其次，用户需要访问的数据中心与 VPN 集中器之间同样可能存在距离问题，此时又会再次产生高延迟。试想这种情况：纽约员工试图访问纽约本地的服务器，却必须通过公司位于旧金山总部的 VPN 集中器或 BeyondCorp 代理进行连接——当然，我并不是在抱怨这种情况。

![hub](https://cdn.sanity.io/images/w77i7m8x/production/d0363ebfb736fa6e394aef3cb26585cecd842cd2-1320x980.svg?w=3840&q=75&fit=clip&auto=format)

Figure 2(a). Inefficient routing through a traditional VPN concentrator. 图 2(a). 通过传统 VPN 集中器的低效路由。

Luckily, WireGuard is different. Remember how we said it creates a “set of extremely lightweight tunnels” above? The tunnels are light enough that we can create a multi-hub setup without too much trouble:  
幸运的是，WireGuard 有所不同。还记得我们前文提到的它能创建"一组极其轻量化的隧道"吗？这些隧道的轻量化程度足以让我们无需耗费太多精力就能构建多中心架构：

![spoke](https://cdn.sanity.io/images/w77i7m8x/production/66c858e17ca544818791e6d23e6085fa2e21b5e4-1320x980.svg?w=3840&q=75&fit=clip&auto=format)

Figure 2(b). A WireGuard multipoint VPN routes traffic more efficiently. 图 2(b). WireGuard 多点 VPN 可更高效地路由流量。

The only catch is that now each of the datacenters needs a static IP address, an open firewall port, and a set of WireGuard keys. When we add a new user, we’ll have to distribute the new key to all five servers. When we add a new server, we’ll have to distribute its key to every user. But we can do it, right? It’s only about 5 times as much work as one hub, which is not that much work.  
唯一的麻烦在于现在每个数据中心都需要一个静态 IP 地址、一个开放的防火墙端口以及一组 WireGuard 密钥。当我们新增用户时，必须将新密钥分发到全部五台服务器；当我们新增服务器时，又需要将其密钥分发给每位用户。但这能搞定吧？不过是一个中心节点五倍的工作量，也不算太多。

#### Mesh networks 网状网络

So that’s hub-and-spoke networks. They’re not too hard to set up with WireGuard, although a bit tedious, and we haven’t really talked about safe practices for key management yet (see below).  
这就是中心辐射型网络（hub-and-spoke）。虽然设置过程略显繁琐，但使用 WireGuard 建立这样的网络并不算太难，不过我们尚未深入讨论密钥管理的安全实践（见下文）。

Still, there’s something fundamentally awkward about hub-and-spoke: they don’t let your individual nodes talk to each other. If you’re old like me, you [remember when computers could just exchange files directly without going to the cloud and back](https://tailscale.com/blog/remembering-the-lan). That was how the Internet all used to work, believe it or not! Sadly, developers have stopped building peer-to-peer apps because the modern Internet’s architecture has evolved, almost by accident, entirely into this kind of hub-and-spoke design, usually with the major cloud providers in the center charging rent.  
不过，中心辐射型架构本质上存在一个根本性缺陷：它不允许各个节点直接通信。如果你和我一样经历过旧时代，就会记得电脑曾无需经由云端就能直接交换文件。无论你信不信，这就是互联网最初的工作方式！遗憾的是，由于现代互联网架构几乎意外地完全演变为这种中心辐射型设计（通常以大型云服务商为核心坐收租金），开发者已不再构建点对点应用。

Remember those “extremely lightweight” WireGuard tunnels? Wouldn’t it be nice if you could directly connect all the nodes to all the other nodes? That’s called a mesh network:  
还记得那些"极其轻量级"的 WireGuard 隧道吗？如果能将所有节点直接连接到其他所有节点，岂不是很好？这就是所谓的网状网络：

![mesh](https://cdn.sanity.io/images/w77i7m8x/production/e989a4a69acd182abbd662d0de93cb31c4c4d210-1600x1080.svg?w=3840&q=75&fit=clip&auto=format)

Figure 3. A Tailscale point-to-point mesh network minimizes latency. 图 3. Tailscale 点对点网状网络可最大程度减少延迟。

That would be very elegant—at least, it would let you design elegant peer-to-peer apps—but it would be tricky. As mentioned above, a 10-node network would require 10 x 9 = 90 WireGuard tunnel endpoint configurations; every node would need to know its own key plus 9 more, and each node would have to be updated every time you rotate a key or add/remove a user.  
这将非常优雅——至少能让你设计出优雅的点对点应用——但实现起来会相当棘手。如前所述，一个 10 节点网络需要 10×9=90 个 WireGuard 隧道端点配置；每个节点都需要知晓自身的密钥外加另外 9 个密钥，且每次轮换密钥或增删用户时，所有节点都必须同步更新。

And the nodes would all have to find each other somehow—user devices rarely have static IP addresses—and reconnect whenever one of them moves around.  
且节点之间必须设法相互发现——用户设备很少拥有静态 IP 地址——并在任一节点移动时重新建立连接。

Plus, you can’t easily open a firewall port for every single node to allow incoming connections in, say, a cafe, hotel, or airport.  
此外，在咖啡馆、酒店或机场等场所，你无法轻易为每个节点单独开放防火墙端口以允许传入连接。

And then, after you’ve done all that, if your company has compliance requirements, you need to be able to somehow block and audit traffic between all the nodes, even though it’s no longer all going through a central location that you can clamp down.  
即使流量不再经过可集中管控的中央枢纽，在完成上述所有配置后，若公司存在合规要求，你仍需设法拦截并审计所有节点间的流量。

Tailscale makes all that work too! Let’s talk about how. This is where things get a bit hairy.  
Tailscale 也能搞定这一切！我们来聊聊具体实现方式。事情就变得有点棘手了。

### The control plane: key exchange and coordination控制平面：密钥交换与协调

Okay, we’ve made it this far. We got WireGuard connecting. We got it connecting to multiple things at a time. We’ve marvelled at its unparalleled reliability and efficiency. Great! Now we want to build a mesh network—everything connected to everything else.  
好的，我们已经取得了长足进展。成功建立了 WireGuard 连接，实现了同时连接多个节点，更惊叹于其无与伦比的可靠性和效率。太棒了！现在我们希望构建一个网状网络——让每个节点都彼此互联。

All those firewalls and dynamic IP addresses are tricky, so let’s leave them aside for the moment. (See further below.) For now, pretend the whole world uses static IPs (dare we say… IPv6?) and that firewalls don’t exist or that it’s easy to open ports.  
所有的防火墙和动态 IP 地址都很棘手，我们暂且搁置不谈（详见下文）。现在假设全世界都使用静态 IP（或许我们可以称之为……IPv6？），并且防火墙不存在或可以轻松开放端口。

How do we get all the WireGuard encryption keys (a simplified and more secure form of “certificates”) onto every device?  
我们如何将所有 WireGuard 加密密钥（一种简化且更安全的“证书”形式）部署到每台设备上？

For that, we use the [open source Tailscale node software](https://github.com/tailscale/tailscale). It talks to what we call a “coordination server” (in our case, login.tailscale.com)—essentially, a shared drop box for public keys.  
为此，我们使用开源的 Tailscale 节点软件。该软件会与所谓的“协调服务器”（即 login.tailscale.com）通信——本质上这是一个用于存储公钥的共享存储库。

![server](https://cdn.sanity.io/images/w77i7m8x/production/dbba97845c1ad1955669cc6a84c94f9d5fb78ade-1600x1080.svg?w=3840&q=75&fit=clip&auto=format)

Figure 4. Tailscale public keys and metadata are shared through a centralized coordination server. 图 4. Tailscale 公钥及元数据通过集中式协调服务器进行共享。

Hold on, are we back to hub-and-spoke again? Not exactly. The so-called “control plane” is hub and spoke, but that doesn’t matter because it carries virtually no traffic. It just exchanges a few tiny encryption keys and sets policies. The data plane is a mesh.  
等等，我们是不是又回到了中心辐射架构？并不尽然。所谓的“控制平面”确实是中心辐射架构，但这无关紧要，因为它几乎不承载任何流量，仅负责交换少量微型加密密钥并设置策略。而数据平面则是网状结构。

(We like this hybrid centralized-distributed model. Companies and teams usually want to have central control, but they don’t want a central bottleneck for their data. Traditional VPN concentrators centralize both, limiting performance, especially under stress. Conversely, Tailscale’s data processing power scales up with the number of nodes.)  
我们青睐这种混合式中心化-分布式架构。企业与团队通常希望拥有中央管控能力，但不愿让数据受制于中心化瓶颈。传统 VPN 集中器将控制与数据流双重中心化，导致性能受限（尤其在压力场景下）。而 Tailscale 的数据处理能力会随节点数量线性扩展。

Here’s what happens:以下是具体过程：

1. Each node generates a random public/private keypair for itself, and associates the public key with its identity (see login, below).  
	每个节点会为自身生成随机公钥/私钥对，并将公钥与其身份标识进行绑定（参见下方登录机制）。
2. The node contacts the coordination server and leaves its public key and a note about where that node can currently be found, and what domain it’s in.  
	该节点联系协调服务器，并留下其公钥及包含该节点当前所在位置和所属域的说明。
3. The node downloads a list of public keys and addresses in its domain, which have been left on the coordination server by other nodes.  
	节点会下载其所在域内其他节点留在协调服务器上的公钥和地址列表。
4. The node configures its WireGuard instance with the appropriate set of public keys.  
	节点使用相应的公钥集合配置其 WireGuard 实例。

Note that the private key never, ever leaves its node. This is important because the private key is the only thing that could potentially be used to impersonate that node when negotiating a WireGuard session. As a result, only that node can encrypt packets addressed from itself, or decrypt packets addressed to itself. It’s important to keep that in mind: Tailscale node connections are end-to-end encrypted (a concept called “ [zero trust networking](https://tailscale.com/kb/1123/zero-trust) ”).  
请注意，私钥永远不会离开其所在节点。这一点非常重要，因为私钥是唯一可能被用来在协商 WireGuard 会话时冒充该节点的凭证。因此，只有该节点本身才能加密发自该节点的数据包，或解密发送至该节点的数据包。需谨记：Tailscale 节点间的连接采用端到端加密（这一概念被称为"zero trust networking"）。

Unlike a hub-and-spoke network, unencrypted packets never need to be sent over a wire, and no intermediary can inspect them. (The exception is if you use [subnet routes](https://tailscale.com/kb/1019/subnets), which are useful for incremental deployment, which Tailscale supports. You can create a hybrid network that combines new-style mesh connections with zero or more old-style “hubs” that decrypt packets, then forward them to legacy physical networks.)  
与中心辐射型网络不同，未加密数据包无需通过物理线路传输，且没有中间方能够检查这些数据包。（唯一例外是使用子网路由时——该模式适用于渐进式部署场景，Tailscale 对此提供了支持。您可以创建混合网络，将新型网状连接与零个或多个旧式"中心节点"结合使用，这些中心节点会解密数据包并将其转发至传统物理网络。）

#### Login and 2-factor auth (2FA)登录与双因素认证（2FA）

But we skipped over a detail. How does the coordination server know which public keys should be sent to which nodes? Public keys are just that—public—so they are harmless to leak to anyone, or even post on a public web site. This is exactly the same situation as an ssh server with an authorized\_keys file; you don’t have to keep your public ssh key secret, but you still have to be careful which public keys you put in authorized\_keys.  
但我们忽略了一个细节。协调服务器如何知道该将哪些公钥发送给哪些节点？公钥——公钥本就是公开的——因此即使泄露给任何人，甚至发布到公共网站上也并无危害。这与带有 authorized\_keys 文件的 ssh 服务器情形完全相同：你无需对自己的 ssh 公钥保密，但仍需谨慎决定将哪些公钥放入 authorized\_keys 文件中。

There are many ways to make the authentication decision. An obvious way would be to build a username+password system, also known as PSK (pre-shared keys). To set up your node, connect to the server, enter your username and password, then upload your public key and download other public keys posted by either your account or other accounts in your domain. If you want to get fancy, you can add two-factor authentication (2FA, also known as MFA) such as SMS, Google Authenticator, Microsoft Authenticator, and so on.  
有多种方式可进行身份验证决策。一种显而易见的方式是构建用户名+密码系统，也称为 PSK（预共享密钥）。设置节点时需连接服务器，输入用户名和密码，随后上传公钥并下载您账户或所在域其他账户发布的公钥。若需增强安全性，可添加双因素认证（2FA，亦称 MFA），例如短信验证、Google Authenticator、Microsoft Authenticator 等。

A system administrator could also set up your machine with a “ [machine certificate](https://tailscale.com/kb/1010/machine-certs) ” — a key that belongs permanently (or semi-permanently) to the device rather than to the user account. It could use this to ensure that, even with the right username and password, an untrusted device can never publish new keys to the coordination server.  
系统管理员还可通过“机器证书”（即永久或半永久归属于设备而非用户账户的密钥）对设备进行配置。系统可利用此证书确保即使拥有正确的用户名和密码，不受信任的设备也永远无法向协调服务器发布新密钥。

Tailscale operates a coordination server based around these concepts. However, we don’t handle user authentication ourselves. Instead, we always outsource authentication to an [OAuth2, OIDC (OpenID Connect), or SAML provider](https://tailscale.com/kb/1013/sso-providers). Popular ones include Gmail, GSuite, and Office365.  
Tailscale 基于上述概念运营着一个协调服务器。不过，我们自身并不处理用户身份验证，而是始终将验证工作交由 OAuth2、OIDC（OpenID Connect）或 SAML 供应商完成，常见的包括 Gmail、GSuite 和 Office365。

![2fa](https://cdn.sanity.io/images/w77i7m8x/production/b86c5249f27186ceb06bd69c852a45627bbca6e5-1600x900.svg?w=3840&q=75&fit=clip&auto=format)

Figure 5. Tailscale 2FA authentication flow in the control plane. 图 5. 控制平面中的 Tailscale 双因素认证流程。

The identity provider maintains a list of users in your domain, passwords, 2FA settings, and so on. This avoids the need to maintain a separate set of user accounts or certificates for your VPN — you can use the authentication you already have set up for use with Google Docs, Office 365, and other web apps. Also, because all your private user account and login data is hosted on another service, Tailscale is able to operate a highly reliable central coordination service while holding a minimum of your users’ personally identifiable information (PII). (Read our [privacy policy](https://tailscale.com/privacy-policy) for more.)  
身份提供商负责维护您域内的用户列表、密码、双因素认证（2FA）设置等信息。这避免了为 VPN 维护独立用户账户或证书体系的需求——您可直接复用已为 Google Docs、Office 365 及其他网络应用配置的认证体系。此外，由于所有用户账户与登录数据均托管在第三方服务，Tailscale 在运营高可靠性中心协调服务的同时，仅保留用户最低限度的个人身份信息（PII）。（详见我们的隐私政策。）

You also don’t have to change how you do single sign-on, account creation and removal, password recovery, and 2FA setup. It’s all exactly like what you already have.  
您也无需更改单点登录、账户创建与删除、密码恢复以及双因素认证(2FA)设置的实现方式。所有流程均与现有方式完全一致。

And because of all that, Tailscale domains can activate instantly the moment you login. If you download our macOS or iOS app from the App Store, for example, and login to your account, this immediately creates a secure key drop box for your domain and lets you exchange public keys with other devices in that account or domain, like a Windows or Linux server.  
正因如此，Tailscale 域能在您登录的瞬间即刻激活。例如，如果您从 App Store 下载我们的 macOS 或 iOS 应用并登录账户，系统会立即为您的域创建一个安全密钥保管箱，并允许您与该账户或域中的其他设备（如 Windows 或 Linux 服务器）交换公钥。

And then, as we just saw, the public keys get downloaded to each node, each node configures WireGuard with a super-lightweight tunnel to each other node, and ta da! A mesh network in two minutes!  
接着，正如我们刚才所见，公钥会被下载到每个节点，每个节点通过 WireGuard 与其他节点之间配置一条超轻量级隧道，瞧！两分钟内就建成了网状网络！

#### NAT traversal NAT 穿透

…but we’re not quite done yet. Recall that up above, we decided to pretend that every node has a static IP address and an open firewall port for incoming WireGuard traffic. In real life, that’s not too likely. In real life, some of your nodes are in cafes or on airplanes or LTE on the highway with one bar of signal, battery drained, fleeing desperately from the cops across state lines… oh, sorry, wrong movie.  
……但我们还没完全搞定。回想前文，我们曾假设每个节点都拥有静态 IP 地址并开放了防火墙端口用于接收 WireGuard 流量。现实中这种情况不太可能发生——你的某些节点可能位于咖啡馆、飞机机舱、高速公路上只有一格信号的 LTE 网络，设备电量耗尽，正在州际公路上拼命逃离警察的追捕……哦抱歉，我串戏了。

Anyway, in the worst case what you have is something like this:  
总之，最坏情况下你得到的结果大致如下：

![remote access](https://cdn.sanity.io/images/w77i7m8x/production/cdd753282121c1279a392aca86adcb6d0a5d1c80-1320x775.svg?w=3840&q=75&fit=clip&auto=format)

Figure 6. Tailscale can connect even when both nodes are behind separate NAT firewalls. 图 6. 即使两个节点位于不同的 NAT 防火墙之后，Tailscale 仍可建立连接。

That’s two NATs, no open ports. Historically, people would ask you to enable uPnP on your firewall, but that rarely works and even when it does work, it usually works [dangerously well](https://www.howtogeek.com/122487/htg-explains-is-upnp-a-security-risk/) until administrators turn it off.  
这里存在双重 NAT 且未开放任何端口。传统做法是建议用户在防火墙上启用 uPnP 协议，但这种方式往往难以奏效，即便偶尔成功，其过度智能化的特性也常导致安全隐患，最终迫使管理员不得不将其关闭。

For now, suffice it to say that Tailscale uses several very advanced techniques, based on the Internet [STUN](https://tools.ietf.org/html/rfc5389) and [ICE](https://tools.ietf.org/html/rfc8445) standards, to make these connections work even though you wouldn’t think it should be possible. This avoids the need for firewall configurations or any public-facing open ports, and thus greatly reduces the potential for human error.  
简而言之，Tailscale 运用基于互联网 STUN 和 ICE 标准的多种尖端技术，使得看似无法建立的连接成为可能。这种方法无需进行防火墙配置或开放任何对外暴露的端口，从而大幅降低了人为错误的发生几率。

For all the gory details, see my teammate Dave Anderson’s post: [How NAT traversal works](https://tailscale.com/blog/how-nat-traversal-works).  
如需了解所有技术细节，请参阅我的队友 Dave Anderson 的文章：NAT 穿透的工作原理。

#### Encrypted TCP relays (DERP)加密 TCP 中继（DERP）

Just one more thing! Some especially cruel networks block UDP entirely, or are otherwise so strict that they simply cannot be traversed using STUN and ICE. For those situations, Tailscale provides a network of so-called DERP (Designated Encrypted Relay for Packets) servers. These fill the same role as [TURN servers](https://tools.ietf.org/html/rfc5766) in the ICE standard, except they use HTTPS streams and WireGuard keys instead of the obsolete TURN recommendations.  
最后再补充一点！某些特别严苛的网络会完全屏蔽 UDP 协议，或者设置如此严格的限制以致根本无法通过 STUN 和 ICE 技术实现穿透。针对这类极端情况，Tailscale 部署了一套名为 DERP（指定加密数据包中继）的服务器网络。其功能相当于 ICE 标准中的 TURN 服务器，区别在于 DERP 使用 HTTPS 数据流与 WireGuard 密钥进行传输，而非过时的 TURN 协议方案。

Relaying through DERP looks like this:  
通过 DERP 进行中继的流程如下：

![relay](https://cdn.sanity.io/images/w77i7m8x/production/0e7f059799b6ba76cfc1df8e7c103d67620d8226-1320x900.svg?w=3840&q=75&fit=clip&auto=format)

Figure 7. Tailscale asymmetrically routes traffic through the DERP nearest to each recipient. 图 7. Tailscale 通过距离每个接收端最近的 DERP 非对称路由流量。

(Despite how it might look above, we really do have a globally distributed network of relay servers. It’s not just the United States. In the words of our intrepid designer, “I was hoping to include London in there, but finding two SVGs with compatible non-Mercator projections was a disaster.” This is how much we care about technical accuracy. Also, stop using Mercator projections. They’re super misleading.)  
（尽管上图可能看起来并非如此，但我们确实拥有一个全球分布式的中继服务器网络。不仅仅是美国。用我们无畏设计师的话来说，“我原本希望把伦敦也加进去，但找到两个使用兼容非墨卡托投影法的 SVG 文件简直是场灾难。”这就是我们对技术准确性的执着程度。另外，请停止使用墨卡托投影法。它们非常误导人。）

Remember that Tailscale private keys never leave the node where they were generated; that means there is never a way for a DERP server to decrypt your traffic. It just blindly forwards already-encrypted traffic from one node to another. This is like any other router on the Internet, albeit using a slightly fancier protocol to prevent abuse.  
请记住，Tailscale 私钥永远不会离开生成它们的节点；这意味着 DERP 服务器永远无法解密您的流量。它只是盲目地将已经加密的流量从一个节点转发到另一个节点。这与互联网上的其他路由器类似，尽管采用了略微复杂的协议来防止滥用。

Many people on the Internet are curious about whether WireGuard supports TCP, like OpenVPN does. Currently this is not built into WireGuard itself, but the open source Tailscale node software includes DERP support, which adds this feature. Over time, it’s possible the code will be refactored to include this feature in WireGuard itself. (Tailscale has already contributed several fixes and improvements to WireGuard-Go.) We have also open sourced [the (quite simple) DERP relay server code](https://github.com/tailscale/tailscale/tree/main/derp) so you can see exactly how it works.  
互联网上许多人好奇 WireGuard 是否像 OpenVPN 那样支持 TCP 协议。目前该功能并未原生集成在 WireGuard 中，但开源项目 Tailscale 的节点软件通过 DERP 支持实现了这一特性。未来代码可能会被重构以原生支持该功能（Tailscale 团队已为 WireGuard-Go 项目贡献了多项修复与改进）。我们同时开源了（非常简洁的）DERP 中继服务器代码，您可查看其具体实现细节。

#### Bonus: ACLs and security policies附加项：访问控制列表(ACLs)和安全策略

People often think of VPNs as “security” software. This categorization seems obvious because VPNs are filled with cryptography and public/private keys and certificates. But VPNs really fit in the category of “connectivity” software: they increase the number of devices that can get into your private network. They don’t decrease it, like security software would.  
人们通常认为 VPN 属于“安全”软件。这种归类似乎显而易见，因为 VPN 充斥着密码学、公钥/私钥和证书体系。但实际上 VPN 应归类为"连接性"软件：它们增加了可接入您私有网络的设备数量，而非像安全软件那样减少接入可能。

As a result, VPN concentrators (remember the hub-and-spoke model from up above) are usually coupled with firewalls; sometimes they’re even sold together. All the traffic comes from individual client devices, flows through the VPN concentrator, then into the firewall, which applies access controls based on IP addresses, and finally into the network itself.  
因此，VPN 集中器（还记得前文提到的中心辐射型模型吗）通常与防火墙配合使用，有时甚至捆绑销售。所有流量均来自独立客户端设备，流经 VPN 集中器后进入防火墙，防火墙根据 IP 地址实施访问控制，最终进入网络内部。

That model works, but it can become a pain. First of all, firewall rules are usually based on IP addresses, not users or roles, so they can be very awkward to configure safely. As a result, you end up having to add more layers of authentication, at the transport or application layers. Why do you need ssh or HTTPS? Because the network layer is too insecure to be trusted.  
该模型有效，但可能带来困扰。首先，防火墙规则通常基于 IP 地址而非用户或角色，配置起来可能非常棘手且难以保证安全。因此最终不得不在传输层或应用层叠加更多认证机制。为何需要 ssh 或 HTTPS？因为网络层安全性过低，完全不可信任。

Second, firewalls are often scattered around the organization and need to be configured individually. If you have a multi-hub network (for example, with different VPN concentrators in different geographic locations), you have to make sure to configure the firewall correctly on each one, or risk a security breach.  
其次，防火墙通常分散在组织各处，需要单独配置。如果您拥有多中心网络（例如在不同地理位置部署了不同的 VPN 集中器），则必须确保在每个节点上正确配置防火墙，否则将面临安全漏洞风险。

Finally, it’s usually a pain to configure some particular vendor’s VPN/firewall device to authenticate VPN connections against some other vendor’s identity system. This is why some identity system vendors will try to sell you a VPN, and some VPN vendors will try to sell you an identity system.  
最后，针对特定厂商的 VPN/防火墙设备配置与其他厂商身份系统进行 VPN 连接认证通常非常麻烦。这正是为什么有些身份系统厂商会试图向你推销 VPN 产品，而有些 VPN 厂商又会试图向你兜售身份系统。

When you switch to a mesh network, this problem seems to get even worse — there is no central point at all, so where do you even put the firewall rules? In Tailscale, the answer is: in every single node. Each node is responsible for blocking incoming connections that should not be allowed, at decryption time.  
当转向使用网状网络时，这个问题似乎变得更加棘手——完全不存在中心节点，防火墙规则该部署在何处？在 Tailscale 中，答案是：部署在每个节点上。每个节点都负责在数据解密阶段拦截本应被禁止的入站连接。

To make that easier, your company’s security policy is stored on the Tailscale coordination server, all in one place, and automatically distributed to each node. This way you have central control over policy, but efficient, distributed enforcement.  
为简化操作，贵公司的安全策略集中存储在 Tailscale 协调服务器上，并自动分发至各节点。由此实现策略的集中管控与高效的分布式执行。

![acl](https://cdn.sanity.io/images/w77i7m8x/production/36ca9e4e723ea056d45695a8da0f56c9094a3e70-1320x900.svg?w=3840&q=75&fit=clip&auto=format)

Figure 8. Central ACL policies are enforced by each Tailscale node’s incoming packet filter. If an ‘accept’ rule doesn’t exist, the traffic is rejected. 图 8. 中央 ACL 策略由每个 Tailscale 节点的入站数据包过滤器强制执行。如果不存在"接受"规则，流量将被拒绝。

At a less granular level, the coordination server (key drop box) protects nodes by giving each node the public keys of only the nodes that are supposed to connect to it. Other Internet computers are unable to even request a connection, because without the right public key in the list, their encrypted packets cannot be decoded. It’s like the unauthorized machines don’t even exist. This is a very powerful protection model; it prevents virtually any kind of protocol-level attack. As a result, Tailscale is especially good at protecting legacy, non-web based services that are no longer maintained or receiving updates.  
在较粗粒度层面，协调服务器（密钥保管箱）通过仅向每个节点提供应连接节点的公钥来实施保护。其他互联网计算机甚至无法发起连接请求，因为若其公钥未在列表内，其加密数据包将无法被解密。未经授权的设备仿佛根本不存在。这是一种极其强大的防护模型，可有效防御几乎所有协议层攻击。因此，Tailscale 特别擅长保护已停止维护或更新的遗留非 Web 服务。

#### Bonus: Audit logs 附加：审计日志

Another concern of companies with tight compliance requirements is audit trails. Modern firewalls don’t just block traffic — they log it. Since Tailscale doesn’t have a central traffic funnel, what can you do?  
合规要求严格的公司的另一个担忧是审计跟踪。现代防火墙不仅阻止流量——还会记录日志。由于 Tailscale 没有中央流量通道，您能采取哪些措施？

The answer is: we allow you to log all internal network connections from each node, asynchronously, to a central logging service. An interesting side effect of this design is that every connection is logged twice: on the source machine and the destination machine. As a result, log tampering is easy to detect, because it would require simultaneous tampering on two different nodes.  
答案是：我们允许您从每个节点异步记录所有内部网络连接至中央日志服务。该设计的一个有趣副作用是每个连接会被记录两次：一次在源机器，一次在目标机器。因此日志篡改很容易被检测到，因为这需要同时在两个不同节点上进行篡改。

Because logs are streamed in real time from each node, rather than batched, the window for local log tampering on a node is extremely short, in the order of dozens of milliseconds, making even a single-node tampering attack difficult to achieve.  
由于日志是从每个节点实时流式传输而非批量传输，因此节点本地篡改日志的时间窗口极短，仅为数十毫秒量级，这使得即便针对单节点的篡改攻击也难以实现。

Tailscale’s central logging service has controllable retention periods and each node just logs some metadata about how your internal mesh is established, not Internet usage or personal information. (Read our [privacy policy](https://tailscale.com/privacy-policy) for more.)  
Tailscale 中央日志服务具有可控的保留周期，每个节点仅记录有关内部网格建立方式的元数据，不涉及互联网使用记录或个人隐私信息。（详见我们的隐私政策。）

Rather than providing a complete logs and metrics pipeline, the Tailscale logging service is intended as a real-time streaming data collector that can then stream data out into other systems for further analysis. (You can [read my early blog post that evolved into Tailscale’s logs collector architecture](https://apenwarr.ca/log/20190216).)  
Tailscale 日志服务并非旨在提供完整的日志和指标处理管道，而是作为实时流式数据收集器，可将数据实时传输至其他系统进行深度分析（可参阅我早期关于 Tailscale 日志收集器架构演进的博文）。

#### Bonus: Incremental deployment额外优势：增量部署

There is one last question that comes up a lot: given that Tailscale creates a mesh “overlay” network (a VPN that parallels a company’s internal physical network), does a company have to switch to it all at once? Many BeyondCorp and zero-trust style products work that way. Or can it be deployed incrementally, starting with a small proof of concept?  
最后一个经常被提及的问题是：既然 Tailscale 创建的是网状"覆盖"网络（一种与公司内部物理网络并行的 VPN），那么公司是否必须一次性全面切换？许多 BeyondCorp 和零信任架构产品都采用这种全量切换方式。或者可以分阶段部署，从一个小型概念验证开始？

Tailscale is uniquely suited to incremental deployments. Since you don’t need to install any hardware or any servers at all, you can get started in two minutes: just install the Tailscale node software onto two devices (Linux, Windows, macOS, iOS), login to both devices with the same user account or auth domain, and that’s it! They’re securely connected, no matter how the devices move around. Tailscale runs on top of your existing network, so you can safely deploy it without disrupting your existing infrastructure and security settings.  
Tailscale 特别适合逐步部署。由于无需安装任何硬件或服务器，您可在两分钟内快速上手：只需在两台设备（Linux、Windows、macOS、iOS）上安装 Tailscale 节点软件，使用同一用户账户或认证域登录这两台设备即可！无论设备如何移动，它们都能安全连接。Tailscale 运行在现有网络之上，因此可安全部署而不会干扰现有基础设施和安全设置。

You can then extend the network by adding [subnet routes](https://tailscale.com/kb/1019/subnets) to one or more offices or datacenters, building up a traditional hub-and-spoke or multi-hub VPN.  
随后，您可以通过向一个或多个办公室或数据中心添加子网路由来扩展网络，构建起传统的中心辐射型或多中心 VPN 架构。

As you gain confidence, you can install Tailscale on more servers or containers, which allows point-to-point fully encrypted connections (no unencrypted traffic carried over any LAN). This final configuration is called “ [zero trust networking](https://www.amazon.com/dp/B072WD347M/ref=dp-kindle-redirect?_encoding=UTF8&btkr=1),” which is gaining fame lately, but so far has been hard to attain.  
随着信心的增强，您可以在更多服务器或容器上安装 Tailscale，这将允许建立点对点全加密连接（不通过任何局域网传输未加密流量）。这种最终配置被称为「零信任网络」，该概念近年来声名鹊起，但迄今为止仍难以真正实现。

With Tailscale, you can build up to “zero trust,” one employee device and one server at a time. After the incremental deployment is done, you can safely shut down your legacy or unencrypted access methods.  
通过 Tailscale，您可以逐步构建「零信任」体系，每次仅需处理一名员工设备与一台服务器。完成增量部署后，即可安全关闭遗留或未加密的访问方式。

Thanks to [Ross Zurowski](https://rosszurowski.com/), our intrepid designer, for the illustrations:)  
感谢我们无畏的设计师 Ross Zurowski 为我们创作了这些插图:)

Share

Author

Avery Pennarun
