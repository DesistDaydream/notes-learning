---
title: 证书 与 PKI
linkTitle: 证书 与 PKI
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, Public Key Certificate](https://en.wikipedia.org/wiki/Public_key_certificate)
> - [Wiki, PKI](https://en.wikipedia.org/wiki/Public_key_infrastructure)
> - [Wiki, CSR](https://en.wikipedia.org/wiki/Certificate_signing_request)
> - [Wiki, CA](https://en.wikipedia.org/wiki/Certificate_authority)
> - [Wiki, Root Certificate](https://en.wikipedia.org/wiki/Root_certificate)
> - [RFC,5280](https://datatracker.ietf.org/doc/html/rfc5280)
> - [公众号,云原生生态圈-白话文说 CA 原理](https://mp.weixin.qq.com/s/E-aU-lbieGLokDKbjdGc3g)
> - [Arthurchiao 博客,\[译\] 写给工程师：关于证书（certificate）和公钥基础设施（PKI）的一切（SmallStep, 2018）](https://arthurchiao.art/blog/everything-about-pki-zh/)
>   - [公众号-云原生实验室，搬运了上面的文章](https://mp.weixin.qq.com/s/li3ZjfNgX5nh7AKjyyzt5A)

Certificate 与 PKI 的目标很简单：Bind names to Public Keys(将名字关联到公钥)。这是关于 Certificate 与 PKI 的最高抽象，其他都是属于实现细节

# Certificate

**Certificate(证书)** 在密码学中，是指 [公开密钥加密](/docs/7.信息安全/Cryptography/公开密钥加密/公开密钥加密.md) 中完善其签名[缺点](/docs/7.信息安全/Cryptography/公开密钥加密/公开密钥加密.md#缺点)的 **Public Key Certificate(公钥证书)**。在公开密钥加密的介绍中，我们看到了公钥加密的特点，并且也发现了缺点，公钥容易被劫持。那么为了解决这个问题，就需要一个东西可以**验证公钥的真实性**。公钥证书也就由此而来。

**Public Key Certificate(公钥证书，简称 PKC)** 也称为 **Digital Certifacte(数字证书)** 或 **Identity Certificate(身份证书)**，是一种用于证明公钥的所有权的电子文档。

假设有这么一种场景：公钥加密系统使我们能知道和谁在通信，但这个事情的前提是：我们必须要有对方的公钥

那么，如果我们不知道对方的公钥，那么该怎么办呢？这时候 Certificate 就出现了。

- 首先，我需要从对方手里拿到公钥和其拥有者的信息
- 那么我如何相信我拿到的信息是真实有效的呢？~可以请一个双方都信任的权威机构，对我拿到的信息做出证明
  - 而这个权威机构用来证明信息有效的东西，就是 Certificate

公钥证书通常应该包含如下内容：

- 密钥的信息
- 有关其所有者的身份信息，称为 Subject(主体)
- 验证证书内容的实体的数字签名，这个实体称为 Issuer(发行人)
- 权威机构对证书的签名，签名的大概意思就是：`Public key XXX 关联到了 name XXX`，这就对应了文章开头的那句话：Certificate 与 PKI 的目标很简单：Bind names to Public Keys(将名字关联到公钥)
  - 对证书的签名的实体称为 **Certificate Authority(简称 CA)**，也可以称为 **Issuer(签发者)**。
  - 被签名的实体称为 **Subject(主体)**。

举个例子，如果某个 Issuer 为 Bob 签发了一张证书，其中的内容就可以解读如下：

_Some Issuer_ says _Bob_’s public key is 01:23:42…

![image.png|800](https://notes-learning.oss-cn-beijing.aliyuncs.com/wlyw54/1634110798410-4fe856d6-2d02-43a9-b233-229b8d48fa51.png)

证书是权威机构颁发的身份证明，并没有什么神奇之处

其中 `Some Issuer` 是证书的签发者(CA)，证书是为了证明这是 Bob 的公钥， Some Issuer 也是这个声明的签字方。

如果签名有效，并且检查证书的软件信任发行者，那么它可以使用该密钥与证书的主题安全地通信。在电子邮件加密，代码签名和电子签名系统中，证书的主体通常是一个人或组织。然而，在传输层安全性（TLS）中，证书的主体通常是计算机或其他设备，但除了在识别设备中的核心作用之外，TLS 证书还可以识别组织或个人。 TLS 有时被其旧的名称安全套接字层（SSL）调用，对于 HTTPS 的一部分是值得注意的，该协议是安全浏览 Web 的协议。

公钥证书最常用的格式是 X.509 标准。

## 证书的类型

1. 自签名证书：一般都是 CA 机构使用 CA 自己的的公钥签署的证书，这样别人拿到该证书后，才可以去找 CA 验证这个证书是不是可信的。
2. 根证书：根证书是标识根证书颁发机构(CA)的公钥证书。根证书是自签名的，并构成基于 X.509 的公钥基础结构（PKI）的基础。
3. TLS / SSL 服务器证书
4. TLS / SSL 客户端证书
5. 电子邮件证书
6. 代码签名证书
7. 合格证书
8. 中级证书
9. 终端实体或叶子证书

# Public Key Infrastructure

**Public Key Infrastructure(公钥基础设施)** 是用来管理数字证书和管理公钥加密所需的一组 Entity(实体，包括但不限于 角色、策略、硬件、软件、程序 等等)。PKI 的目的是为了促进一系列网络活动的安全传输，例如电子商务、网上银行、敏感电子邮件。

PKI 是通用的，与厂商无关的概念，适用于任何地方，因此及时系统分布在世界各地，彼此之间也能安全地通信；如果我们使用 TLS everywhere 模型，甚至连 VPN 都不需要。

PKI 的标准由 RFC 5280 定义，然后 [CA/Browser Forum](https://cabforum.org/) (a.k.a., CA/B or CAB Forum) 对其进行了进一步完善。PKI 的标准称为 **Public-Key Infrastructure(X.509)(公要基础设施(X.509)，简称 PKIX)**，由于 PKIX 是围绕 X.509 证书标准定义的 PKI 标准，所以 PKI 后面就加了一个 X~~~~~o(╯□╰)o ~~~从 RFC 5280 也可以看到每页文档的页眉写的是 `PKIX Certificate and CRL Profile`。并且 IETF 与 1995 年秋季成了了 [PKIX 工作组](https://datatracker.ietf.org/wg/pkix/about/)。

## Certificate Authority(证书权威)

**Certificate Authority(证书权威，简称 CA)** 是拥有公信力的颁发数字证书的实体，通常称为**证书颁发机构**。

CA 自身的证书，通常称为 **Root Certificate(根证书)**，根证书是使用 CA 的私钥 **Self-signed(自签名的)**，且是 CA 的唯一标识。

**CA 使用自己的私钥为其他实体签名并颁发证书**。就像文章开头提到的一样，Bind names to Public Keys，CA 为公钥和名字之间的绑定关系做担保。

### Trust Stores(信任仓库)，即操作系统、浏览器等保存证书的地方

那么，当我们访问一个网站时，是如何验证其证书的真实性呢?~其实，通常都会有一个 CA 仓库，用来保存一些预配置的可信的根证书。凭借这些预配置的根证书，当我们访问互联网上绝大部分网站时，就可以验证其身份。

这个信任仓库又是如何来的呢？

- 浏览器 # 浏览器默认使用的信任仓库以及其他任何使用 TLS 的东西，都是由 4 个组织维护的
  - [Apple’s root certificate](http://www.apple.com/certificateauthority/ca_program.html)：iOS/macOS 程序
  - [Microsoft’s root certificate program](https://social.technet.microsoft.com/wiki/contents/articles/31633.microsoft-trusted-root-program-requirements.aspx)：Windows 使用
  - [Mozilla’s root certificate program](https://www.mozilla.org/en-US/about/governance/policies/security-group/certs/)： Mozilla 产品使用，由于其开放和透明，也作为其他一些信任仓库从基础 (e.g., for many Linux distributions)
  - Google [未维护 root certificate program](https://www.chromium.org/Home/chromium-security/root-ca-policy) （Chrome 通常使用所在计算的操作系统的信任仓库），但 [维护了自己的黑名单](https://chromium.googlesource.com/chromium/src/+/master/net/data/ssl/blacklist/README.md)， 列出了自己不信任的根证书或特定证书。 ([ChromeOS builds off of Mozilla’s certificate program](https://chromium.googlesource.com/chromiumos/docs/+/master/ca_certs.md))
- 操作系统 # 操作系统中的信任仓库通常是各自发行版自带的。不同的发行版，保存路径不同，保存方式也不同：
  - CentOS 发行版
    - /etc/pki/ca-trust/extracted/openssl/ca-bundle.trust.crt # 包含所有证书，每个证书前有注释
  - Ubuntu 发行版
    - /etc/ssl/certs/\* # 该目录中一个证书一个文件
  - Windows，证书位置如图
    - ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wlyw54/1662898111701-e46d8a99-c518-48fa-8eb2-9a37448d3df3.png)

信任仓库中通常包含了超过 100 个由这些程序维护的常见 CA。比如：

- Let’s Encrypt
- Symantec
- DigiCert
- Entrust

#### 示例

比如在 Linux 各种发行版中上述目录中有一个名为 GlobalSign_Root_CA 的根证书，百度就是使用这个证书签名的。

这是通过浏览器访问 http://www.baidu.com 获取到的证书信息

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wlyw54/1638255619893-aaa4aaf0-0b19-4aab-94c7-ea6d52c40e8b.png)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wlyw54/1638255970307-44dc01bc-7d7d-4adb-94c7-bb7ec2ce1636.png)

这是从服务器的 CA 仓库中获取的 GlobalSign_Root_CA 这个 CA 的证书信息

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wlyw54/1638255663132-94adcb92-b634-4f45-936b-8f51245f7558.png)

```bash
[root@hw-cloud-xngy-jump-server-linux-2 ~]# openssl x509 -text -noout -in /etc/ssl/certs/GlobalSign_Root_CA.pem
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            04:00:00:00:00:01:15:4b:5a:c3:94
        Signature Algorithm: sha1WithRSAEncryption
        Issuer: C = BE, O = GlobalSign nv-sa, OU = Root CA, CN = GlobalSign Root CA
        Validity
            Not Before: Sep  1 12:00:00 1998 GMT
            Not After : Jan 28 12:00:00 2028 GMT
        Subject: C = BE, O = GlobalSign nv-sa, OU = Root CA, CN = GlobalSign Root CA
......
```

可以看到，两个证书的信息是一样的，时间不一致是由于时区设置的问题，浏览器直接打开证书变成了东八区。所以是 20:00:00；如果从窗口导出证书成文件，再使用 openssl 命令查看，就可以发现，两个证书是一模一样的。

### Chain of trust(信任链)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wlyw54/1638252283465-124564d9-4f52-4812-9cfb-484fb54b599b.png)

CA 使用自己的私钥签一个根证书，然后再为下级 Issuer 签署证书，下级 Issuer 还可以为其自身的下级 Issuer 签署证书。这么层层签署，可以形成一个树形结构的信任链。

现在假如 A 是 CA，签署证书给 B，C 想访问 B 提供的服务，C 在访问时，如何确保 B 就是 B 呢？~这就是 TLS/SSL 协议所要做的事情。总结一下就是 C 首先要获取 A 的证书，这是访问 B 时，就可以使用 A 的证书验证 B 的证书。并且，由于 A 的证书是用其私钥签名的，只要 A 签证书的私钥(即 CA 的私钥)不泄露，整个信任链就是可信的。

### 保证 CA 私钥的安全

CAB Forum Baseline Requirements 4.3.1 明确规定：一个 Web PKI CA 的 root private key 只能通过 issue a direct command 来签发证书。

- 换句话说，Web PKI root CA 不能自动化证书签名（certificate signing）过程。
- 对于任何大的 CA operation 来说，无法在线完成都是一个问题。 不可能每次签发一个证书时，都人工敲一个命令。

这样规定是出于安全考虑。

- Web PKI root certificates 广泛存在于信任仓库中，很难被撤回。截获一个 root CA private key 理论上将影响几十亿的人和设备。
- 因此，最佳实践就是，确保 root private keys 是离线的（offline），理想情况下在一些 [专用硬件](https://en.wikipedia.org/wiki/Hardware_security_module) 上，连接到某些物理空间隔离的设备上，有很好的物理安全性，有严格的使用流程。

一些 internal PKI 也遵循类似的实践，但实际上并没有这个必要。

- 如果能自动化 root certificate rotation （例如，通过配置管理或编排工具，更新信任仓库）， 你就能轻松地 rotate 一个 compromised root key。
- 由于人们如此沉迷于 internal PKI 的根秘钥管理，导致 internal PKI 的部署效率大大 降低。你的 AWS root account credentials 至少也是机密信息，你又是如何管理它的呢？

### 其他

可以通过工具来创建私有 CA 证书(即自己创造一个 CA 所用的证书，等于是自己的其中一台设备当做 CA 来给自己所用的其余设备颁发证书)(自签名的)以便让个人或公司内部使用。比如 openssl 工具就可以实现。openssl 即可创建私有 CA 证书。

## Certificate Signing Request(证书签名请求)

在 PKI 系统中，**Certificate Signing Request(证书签名请求，简称 CSR)** 是申请人发送到 PKI 的注册机构，用来申请**公钥证书**的一种消息。

CSR 最常见的格式是 [PKCS](https://en.wikipedia.org/wiki/PKCS)＃10 规范；另一个是某些[浏览器](https://en.wikipedia.org/wiki/Web_browser)生成的签名公钥和质询 [SPKAC](https://en.wikipedia.org/wiki/SPKAC) 格式。

在创建 CSR 之前，申请人首先需要生成一个密钥对，并将私钥保密。实际上，CSR 也可以称为证书，想要创建一个 CSR，则需要使用申请人的密钥中的私钥进行签名。

CSR 中应包含

- **名字** # 申请人的识别信息(比如 X.509 规范中的 Subject 字段)
- **公钥** # 从申请人密钥中提取出的公钥
- **签名** #
- **其他信息** #

# 证书的验证过程

CA 收到一个 CSR 并验证签名之后，接下来需要确认证书中绑定的 name 是否真的 是这个 subscriber 的 name。这项工作很棘手。 证书的核心功能是**能让 RP 对 subscriber 进行认证**。因此， 如果一个**证书都还没有颁发，CA 如何对这个 subscriber 进行认证呢**？

答案是：分情况。

## Web PKI 证明身份过程

Web PKI 有三种类型的证书，它们**最大的区别就是如何识别 subscriber**， 以及它们所用到的 **identity proofing 机制**。

这三种证书是：

1. domain validation (DV，域验证)DV 证书绑定的是 **DNS name**，CA 在颁发时需要验证的这个 domain name 确实是由该 subscriber 控制的。证明过程通常是通过一个简单的流程，例如
   1. 给 WHOIS 记录中该 domain name 的管理员发送一封确认邮件。
   2. [ACME protocol](https://ietf-wg-acme.github.io/acme/draft-ietf-acme-acme.html) （最初由 Let’s Encrypt 开发和使用）改进了这种方式，更加自动化：不再用邮件验证 ，而是由 ACME CA 提出一个 challenge，该 subscriber 通过完成这个问题来证明它拥有 这个域名。challenge 部分属于 ACME 规范的扩展部门，常见的包括：
      - 在指定的 URL 上提供一个随机数（HTTP challenge）
      - 在 DNS TXT 记录中放置一个随机数（DNS challenge）
2. organization validation (OV，组织验证)
   - OV 和下面将介绍的 EV 证书构建在 DV 证书之上，它们包括了 name 和域名 **所属组织的位置信息（location）**。
   - OV 和 EV 证书不仅仅将证书关联到域名，还关联到控制这个域名的法律实体（legal entity）。
   - OV 证书的验证过程，不同的 CA 并不统一。为解决这个问题，CAB Forum 引入了 EV 证书。
3. **extended validation** (EV，扩展验证)这些完成之后，当相应网站时，**某些浏览器会在 URL 栏中显示该组织的名称**。例如：![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wlyw54/1638261597799-3bbc2a87-727d-4c25-8492-f993a2e94ead.png)但除了这个场景之外，EV certificates 并未得到广泛使用，Web PKI RP 也未强依赖它。
   - EV 证书包含的基本信息与 OV 是一样的，但强制要求严格验证（identity proofing）。
   - EV 过程需要几天或几个星期，其中可能包括公网记录搜索（public records searches）和公司人员（用笔）签署的（纸质）证词。

**本质上来说，每个 Web PKI RP 只需要 DV 级别的 assurance** 就行了， 也就是确保域名是被该 subscriber 控制的。重要的是能理解一个 DV 证书在设计上的意思和在实际上做了什么：

- 在设计上，希望通过它证明：请求这个证书的 entity 拥有对应的域名；
- 在实际上，真正完成的操作是：在某个时间，请求这个证书的 entity 能读一封邮件，或配置一条 DNS 记录，或能通过 HTTP serve 一个指定随机数等等。

但话说回来，DNS、电子邮件和 BGP 这些底层基础设施本身的安全性也并没有做到足够好， 针对这些基础设施的攻击还是 [时有发生](https://doublepulsar.com/hijack-of-amazons-internet-domain-service-used-to-reroute-web-traffic-for-two-hours-unnoticed-3a6f0dda6a6f)， 目的之一就是获取证书。

## Internal PKI 证明身份过程

上面是 Web PKI 的身份证明过程，再来看 internal PKI 的身份证明过程。

实际上，用户可以使用**任何方式**来做 internal PKI 的 identity proofing， 并且效果可能比 Web PKI 依赖 DNS 或邮件方式的效果更好。

乍听起来好像很难，但其实不难，因为可以**利用已有的受信基础设施**： 用来搭建基础设施的工具，也能用来为这些基础设施之上的服务创建和证明安全身份。

- 如果用户已经信任 Chef/Puppet/Ansible/Kubernetes，允许它们将代码放到服务器上， 那也应该信任它们能完成 identity attestations
- 如果在 AWS 上，可以用 [instance identity documents](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-identity-documents.html)
- 如果在 GCP：[GCP](https://cloud.google.com/compute/docs/instances/verifying-instance-identity)
- [Azure](https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/how-to-use-vm-token)

provisioning infrastructure 必须理解 identity 的概念，这样才能将正确的代码放到正确的位置。 此外，用户必须信任这套机制。基于这些知识和信任，才能配置 RP 信任仓库、将 subscribers 纳入你的 internal PKI 管理范围。 而完成这些功能全部所需做的就是：设计和实现某种方式，能让 provisioning infrastructure 在每个服务启动时，能将它们的 identity 告诉你的 CA。 顺便说一句，这正是 [step certificates](https://smallstep.com/certificates/) 解决的事情。

# 证书的生命周期

## Expiration（过期）

证书通常都会过期。虽然这不是强制规定，但一般都这么做。设置一个过期时间非常重要，

- **证书都是分散在各处的**：通常 RP 在验证一个证书时，并没有某个中心式权威能感知到（这个操作）。
- 如果没有过期时间，证书将永久有效。
- 安全领域的一条经验就是：**时间过的越久，凭证被泄露的概率就越接近 100%**。

因此，设置过期时间非常重要。具体来说，X.509 证书中包含一个有效时间范围：

1. _issued at_
2. _not before_
3. _not after_：过了这个时间，证书就过期了。

这个机制看起来设计良好，但实际上也是有一些不足的：

- 首先，**没有什么能阻止 RP** 错误地（或因为糟糕的设计）**接受一个过期证书**；
- 其次，证书是分散的。验证证书是否过期是每个 RP 的责任，而有时它们会出乱子。例如，**RP 依赖的系统时钟不对**时。 **最坏的情况就是系统时钟被重置为了 unix epoch**（1970.1.1），此时它无法信任任何证书。

在 subscriber 侧，证书过期后，私钥要处理得当：

- 如果一个密钥对之前是**用来签名/认证**的（例如，基于 TLS），
  - 应该在不需要这个密钥对之后，**立即删除私钥**。
  - 保留已经失效的签名秘钥（signing key）会导致不必要的风险：对谁都已经没有用处，反而会被拿去仿冒签名。
- 如果密钥对是**用来加密的**，情况就不同了。
  - 只要还有数据是用这个加密过的，就需要**留着这个私钥**。

这就是为什么很多人会说，**不要用同一组秘钥来同时做签名和加密**（signing and encryption）。 因为当一个用于签名的私钥过期时，**无法实现秘钥生命周期的最佳管理**： 最终不得不保留着这个私钥，因为解密还要用它。

## Renewal（续期）

证书快过期时，如果还想继续使用，就需要续期。

### Web PKI 证书续期

Web PKI 实际上并**没有标准的续期过期**：

- 没有一个标准方式来延长证书的合法时间，
- 一般是**直接用一个新证书替换过期的**。
- 因此续期过程和 issuance 过程是一样的：**生成并提交一个 CSR**，然后完成 identity proofing。

### Internal PKI 证书续期

对于 internal PKI 我们能做的更好。
最简单的方式是：

- **用 mTLS 之类的协议对老证书续期**。
- CA 能对 subscriber 提供的客户端证书进行认证（authenticate），**重签一个更长的时间**，然后返回这个证书。
- 这使得续期过程**很容易自动化**，而且强制 subscriber 定期与中心权威保持沟通。
- 基于这种机制能轻松**构建一个证书的监控和撤销基础设施**。

### 小结

证书的续期过程其实并不是太难，**最难的是记得续期这件事**。
几乎每个管理过公网证书的人，都经历过证书过期导致的生产事故，[例如这个](https://expired.badssl.com/)。 我的建议是：

1. 发现问题之后，一定要全面排查，解决能发现的所有此类问题。
2. 另外，使用生命周期比较短的证书。这会反过来逼迫你们优化和自动化整个流程。

Let’s Encrypt 使自动化非常容易，它签发 90 天有效期的证书，因此对 Web PKI 来说非常合适。 对于 internal PKI，建议有效期签的更短：24 小时或更短。有一些实现上的挑战 —— [hitless certificate rotation](https://diogomonica.com/2017/01/11/hitless-tls-certificate-rotation-in-go/) 可能比较棘手 —— 但这些工作是值得的。

用 step 检查证书过期时间：

```bash
step certificate inspect cert.pem --format json | jq .validity.end
step certificate inspect https://smallstep.com --format json | jq .validity.end
```

将这种命令行封装到监控采集脚本，就可以实现某种程度的监控和自动化。

## Revocation（撤销）

如果一个私钥泄露了，或者一个证书已经不再用了，就需要撤销它。即希望：

1. 明确地将其标记为非法的，
2. 所有 RP 都不再信任这个证书了，即使它还未过期。

但实际上，**撤销证书过程也是一团糟**。

### 主动撤销的困难

- 与过期类似，**执行撤回的职责在 RP**。
- 与过期不同的是，**撤销状态无法编码在证书中**。RP 只能依靠某些带外过程（out-of-band process） 来判断证书的撤销状态。

除非显式配置，否则大部分 Web PKI TLS RP 并不关注撤销状态。换句话说，默认情况下， 大部分 TLS 实现都乐于接受已经撤销的证书。

### Internal PKI：被动撤销机制

Internal PKI 的趋势是接受这个现实，然后试图通过**被动撤销**（passive revocation）机制来弥补， 具体来说就是**签发生命周期很短的证书**，这样就使撤销过程变得不再那么重要了。 想撤销一个证书时，直接不给它续期就行了，过一段时间就会自动过期。
可以看到，**这个机制有效的前提**就是使用生命周期很短的证书。具体有多短？

1. 取决于你的威胁模型（安全专家说了算）。
2. 24 小时是很常见的，但也有短到 5 分钟的。
3. 如果生命周期太短，显然也会给可扩展性和可用性带来挑战：**每次续期都需要与 online CA 交互**， 因此 CA 有性能压力。
4. 如果缩短了证书的生命周期，记得**确保你的时钟是同步的**，否则就有罪受了。

对于 web 和其他的被动撤销不适合的场景，如果认真思考之后发现**真的** 需要撤销功能，那有两个选择：

1. CRL（，**证书撤销列表**，RFC 5280）
2. OCSP（Online Certificate Signing Protocol，**在线证书签名协议**，RFC 2560）

### 主动检查机制：CRL

**Certificate Revocation Lists(证书吊销列表，简称 CRL)** 定义在 RFC 5280 中，这是一个相当庞杂的 RFC，还定义了很多其他东西。 简单来是，CRL 是一个**有符号整数序列，用来识别已撤销的证书**。
这个维护在一个 **CRL distribution point** 服务中，每个证书中都包含指向这个服务的 URL。 工作流程：每个 RP 下载这个列表并缓存到本地，在对证书进行验证时，从本地缓存查询撤销状态。 但这里也有一些明显的问题：

1. **CRL 可能很大**，
2. distribution point 也可能失效。
3. RP 的 CRL 缓存同步经常是天级的，因此如果一个证书撤销了，可能要几天之后才能同步到这个状态。
4. 此外，RP _fail open_ 也很常见 —— CRL distribution point 挂了之后，就接受这个证书。 这显然是一个安全问题：只要对 CRL distribution point 发起 DDoS 攻击，就能让 RP 接受一个已经撤销的证书。

因此，即使已经在用 CRL，也应该考虑使用短时证书来保持 CRL size 比较小。 CRL 只需要包含**已撤销但还未过期的证书**的 serial numbers，因此 证书生命周期越短，CRL 越短。

### 主动检查机制：OCSP

主动检查机制除了 CRL 之外，另一个选择是 **Online Certificate Signing Protocol(简称 OCSP)**，它允许 RP 实时查询一个 _OCSP responder_： 指定证书的 serial number 来获取这个证书的撤销状态。
与 CRL distribution point 类似，OCSP responder URL 也包含在证书中。 这样看，OCSP 似乎更加友好，但实际上它也有自己的问题。对于 Web PKI，它引入了验证的隐私问题：

1. 每次查询 OCSP responder，使得它能看到我正在访问哪个网站。
2. 此外，它还增加了每个 TLS 连接的开销：需要一个额外请求来检查证实的撤销状态。
3. 与 CRL 一样，很多 RPs (including browsers) 会在 OCSP responder 失效时直接认为证书有效（未撤销）。

### 主动检查机制：OCSP stapling（合订，绑定）

OCSP stapling 是 OCSP 的一个变种，目的是解决以上提到的那些问题。

相比于让 RP 每次都去查询 OCSP responder，OCSP stapling 中让证书的 subscriber 来做这件事情。 OCSP response 是一个经过签名的、时间较短的证词（signed attestation），证明这个证书未被撤销。

attestation 包含在 subscriber 和 RP 的 TLS handshake (“stapled to” the certificate) 中。 这给 RP 提供了相对比较及时的撤销状态，而不用每次都去查询 OCSP responder。 subscriber 可以在 signed OCSP response 过期之前多次使用它。这减少了 OCSP 的负担，也解决了 OCSP 的隐私问题。

但是，所有这些东西其实最终都像是一个 **鲁布·戈德堡装置（Rube Goldberg Device）**，

> 鲁布·戈德堡机械（Rube Goldberg machine）是一种被设计得过度复杂的机械组合，以 迂回曲折的方法去完成一些其实是非常简单的工作，例如倒一杯茶，或打一只蛋等等。 设计者必须计算精确，令机械的每个部件都能够准确发挥功用，因为任何一个环节出错 ，都极有可能令原定的任务不能达成。
> 解释来自 [知乎](https://www.zhihu.com/topic/20017497/intro)。

如果让 subscribers 去 CA 获取一些生命周期很短的证词（signed attestation）来证明对应的证书并没有过期， 为什么不直接干掉中间环节，直接使用生命周期很短的证书呢？

# 证书申请及签署步骤

1. 生成申请请求
2. RA 核验你的申请信息
3. CA 签署
4. 获取证书(从证书存取库)

私有 CA 的创建以及签发证书步骤，详细命令详见 2.0.OpenSSL.note 命令说明

1. (可选)配置需要使用 CA 功能服务器的 CA 配置文件
2. 在 CA 功能服务器上创建自签证书以便给其余设备签证
3. 在需要签证的设备上创建密钥以及证书签署请求，并把请求文件发送给 CA 服务器
4. CA 服务器给该请求签证后，把生成的证书文件发还给需要签证的设备。

吊销证书

1. 获取证书的 serial
2. 根据用户提交的 serial 与 subject 信息，对比验证是否与 index.txt 文件中的信息一致，使用吊销命令吊销/etc/pki/CA/newcerts/目录下对应的证书文件
3. 生成吊销证书的编号
4. echo 01 > /etc/pki/CA/crl/NUM
5. 更新证书吊销列表

证书文件格式：

1. XXX.pem # 证书相关文件标准格式
2. XXX.key # 明确指明这是一个密钥文件
3. XXX.csr # Certificate signing request。证书签署请求文件
4. XXX.crt # 明确指明这是一个证书文件
