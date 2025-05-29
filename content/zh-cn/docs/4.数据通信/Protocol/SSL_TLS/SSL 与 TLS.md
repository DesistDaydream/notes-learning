---
title: SSL 与 TLS
linkTitle: SSL 与 TLS
weight: 1
tags:
  - Information_security
---

# 概述

> 参考：
>
> - [Wiki, TLS](https://en.wikipedia.org/wiki/Transport_Layer_Security)

为了解决人类在互联网世界信息的安全性，所研究出来的相关技术

安全机制：加密、数字签名、访问控制、数据完整性、认证交换、流量填充、路由控制、公证

安全服务：认证、访问控制、数据保密性(连接保密性、无连接保密性、选择与保密性、流量保密性)、数据完整性、不可否认性

## SSL/TLS 介绍

**Secure Socket Layer(安全的套接字层，简称 SSL)** # 一个安全协议

**Transport Layer Security(传输层安全，简称 TLS)**# SSL3.0 的升级版

SSL/TLS 就是在应用层与传输层中间又加了半层，应用层协议可以自行决定改层的功能，比如 http 协议用了这半层，就是 https。

SSL/TLS 的分层设计

1. 最底层，基础算法原语的实现，比如 aes，rsa，md5 等
2. 各种算法的实现
3. 组合算法实现的半成品
4. 用各种组件拼装而成的各种成品密码学协议/软件，tls,ssh 等 openssh 也是用 openssl 实现的软件

**key(密钥)** # 在密码学中，是指某个用来完成加密、解密、完整性验证等密码学应用的秘密信息。对于加密算法，key 指定明文转换成密文；对于解密算法，key 指定密文转换成明文

- **Plaintext or Cleartext(明文)** # 在密码学中，明文是未加密的信息，可以供人类和计算机读取的信息
- **Ciphertext or Cyphertext(密文)**# 在密码学中，密文是明文通过加密算法计算后生成的人类或计算器无法读取的一种信息

**PKI：Public Key Infrastructure(公开密钥基础建设，简称 PKI)**，又称公开密钥基础架构、公钥基础建设、公钥基础设施、公开密码匙基础建设或公钥基础架构，是一组由硬件、软件、参与者、管理政策与流程组成的基础架构，其目的在于创造、管理、分配、使用、存储以及撤销数字证书。

PKI 是借助 CA（权威数字证书颁发/认证机构）将用户的个人身份跟公开密钥链接在一起，它能够确保每个用户身份的唯一性，这种链接关系是通过注册和发布过程实现，并且根据担保级别，链接关系可能由 CA 和各种软件或在人为监督下完成。PKI 用来确定链接关系的这一角色称为 RA（Registration Authority, 注册管理中心），RA 能够确保公开密钥和个人身份链接，可以防抵赖，防篡改。在微软的公钥基础建设下，RA 又被称为 CA，目前大多数称为 CA。

PKI 组成要素
从上面可以得知 PKI 的几个主要组成要素，用户（使用 PKI 的人或机构），认证机构（CA，颁发证书的人或机构），仓库（保存证书的数据库）等。

1. 签证机构：CA(Certificate authority)证书权威机构
2. 注册机构：RA
3. 证书吊销列表:CRL
4. 证书存取库：
5. x.509：一种证书格式规范

# 通信加密安全实例

甲要发送数据给乙，甲为了只让乙看到

1. 首先，甲用单向加密算法提取数据的特征码，然后用自己的私钥加密这段特征码，并附加在这段数据的后面。
2. 甲用对称密钥，把整个数据加密。再用乙的公钥加密这个对称密钥，并附加在特征码后面
3. 乙先用自己的私钥解密出来对称密钥是什么。再使用对称加密机制，用解密出来的对称密钥解密整个数据和加密的特征码。
4. 乙再对方的公钥解密特征码，得到特征码，使用同样的单向加密算法计算特征码是否一样，则说明数据完整

密钥交换：IKE，DH

openSSL 与 gpg(pgp 协议)

# 关联文件

Linux 发行版中，有一个目录会保存一些常见的 CA 证书，称之为[信任仓库](/docs/7.信息安全/Cryptography(密码学)/公开密钥加密/证书%20 与%20PKI.md 与 PKI.md)：

- CentOS 发行版
  - **/etc/pki/ca-trust/extracted/openssl/ca-bundle.trust.crt** # 包含所有证书，每个证书前一行有注释
- Ubuntu 发行版
  - **/etc/ssl/certs/** # 该目录中一个证书一个文件

# TLS 扩展

> 参考：
>
> - [RFC 3546](https://datatracker.ietf.org/doc/html/rfc3546)

## SNI

> 参考：
>
> - [RFC 3546, section-3.1](https://datatracker.ietf.org/doc/html/rfc3546#section-3.1)
> - [Wiki, Server Name Indication](https://en.wikipedia.org/wiki/Server_Name_Indication)

**Server Name Indication(服务器名称指示，简称 SNI)** 是 TLS 协议的扩展， [客户端](https://en.wikipedia.org/wiki/Client_\(computing\) "Client (computing)")通过该协议在握手过程开始时指示其尝试连接的[主机名](https://en.wikipedia.org/wiki/Hostname "Hostname") 。

SNI 的出现主要是为了解决这么一个问题：

在 HTTP 的 TLS 握手过程中，客户端首先需要向服务端请求证书，如果一个服务器上只部署了一个网站还比较好办，但同一个服务器上部署了多个网站，就必须有一个机制来区分客户端要访问哪一个网站

客户端在 TLS 握手过程中，主动携带了一个 SNI 的数据字段，指明自己要连接的网站，服务端选择 SNI 中指定的证书发送给客户端，从而完成HTTPS的握手

![1000](https://notes-learning.oss-cn-beijing.aliyuncs.com/ssl_tls/sni_demo.png)


不幸的是，在 HTTPS 协议中，TLS握手发生在数据的正式传递之前，因此其中的数据都是明文传递的，这就意味着，SNI 信息可能会被网络传说中的中间人监听，中间人可以轻易的知道用户准备访问哪个网站，用户因此就会泄露隐私，中间人还可以通过识别 SNI 的信息，阻断一部分TLS握手的建立

这种攻击也就是 <font color="#ff0000">SNI 阻断</font>

