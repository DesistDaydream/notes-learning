---
title: SSL/TLS Pinning
linkTitle: SSL/TLS Pinning
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, HTTP_Public_Key_Pinning](https://en.wikipedia.org/wiki/HTTP_Public_Key_Pinning)
> - [OWASP，Certificate and Public Key Pinning](https://owasp.org/www-community/controls/Certificate_and_Public_Key_Pinning)
> - [知乎，证书锁定SSL Pinning简介及用途](https://zhuanlan.zhihu.com/p/58204817)

**SSL/TLS Pinning** 也可以称为 **Public Key Pinning**、**Certificate Pinning**。顾名思义，将服务器提供的 SSL/TLS 证书内置到移动端开发的 APP 客户端中，当客户端发起请求时，通过比对内置的证书和服务器端证书的内容，以确定这个连接的合法性。

在公共网络中通知我们使用安全的 SSL/TLS 通信协议来进行通信，并且使用数字证书来提供加密和认证，在《[HTTPS入门, 图解SSL从回车到握手](https://link.zhihu.com/?target=https%3A//www.infinisign.com/faq/ssl-hello-process)》过程中我们知道握手环节仍然面临（MIM中间人）攻击的可能性，因为CA证书签发机构也存在被黑客入侵的可能性，同时移动设备也面临内置证书被篡改的风险。

# SSL/TLS Pinning 原理

证书锁定（SSL/TLS Pinning）提供了两种锁定方式： **Certificate Pinning** 和 **Public Key Pinning**，文头和概述描述的实际上是 Certificate Pinning（证书锁定）。

HTTP Public Key Pinning(简称 HPKP) 在 RFC 7469 中标准化。扩展了 Certificate Pinning，Certificate Pinning 对 Web 浏览器和应用程序中的知名网站或服务的公钥哈希进行硬编码。

## 证书锁定

我们需要将APP代码内置仅接受指定域名的证书，而不接受操作系统或浏览器内置的CA根证书对应的任何证书，通过这种授权方式，保障了APP与服务端通信的唯一性和安全性，因此我们移动端APP与服务端（例如API网关）之间的通信是可以保证绝对安全。但是CA签发证书都存在有效期问题，所以缺点是在证书续期后需要将证书重新内置到APP中。

## 公钥锁定

公钥锁定则是提取证书中的公钥并内置到移动端APP中，通过与服务器对比公钥值来验证连接的合法性，我们在制作证书密钥时，公钥在证书的续期前后都可以保持不变（即密钥对不变），所以可以避免证书有效期问题。

# 证书锁定指纹(Hash)

## 获取移动端所需证书

如果采用证书锁定方式，则获取证书的摘要hash，以 [infinisign.com](https://link.zhihu.com/?target=https%3A//www.infinisign.com/infinisign.com) 为例

```text
## 在线读取服务器端.cer格式证书
openssl s_client -connect infinisign.com:443 -showcerts < /dev/null | openssl x509 -outform DER > infinisign.der
## 提取证书的摘要hash并查看base64格式
openssl dgst -sha256 -binary infinisign.der | openssl enc -base64
wLgBEAGmLltnXbK6pzpvPMeOCTKZ0QwrWGem6DkNf6o=
```

所以其中的`wLgBEAGmLltnXbK6pzpvPMeOCTKZ0QwrWGem6DkNf6o=`就是我们将要进行证书锁定的指纹(Hash)信息。

## 获取移动端所需公钥

如果采用公钥锁定方式，则获取证书公钥的摘要hash，以[infinisign.com](https://link.zhihu.com/?target=https%3A//www.infinisign.com/infinisign.com)为例

```bash
## 在线读取服务器端证书的公钥
openssl x509 -pubkey -noout -in infinisign.der -inform DER | openssl rsa -outform DER -pubin -in /dev/stdin 2>/dev/null > infinisign.pubkey
## 提取证书的摘要hash并查看base64格式
openssl dgst -sha256 -binary infinisign.pubkey | openssl enc -base64
bAExy9pPp0EnzjAlYn1bsSEGvqYi1shl1OOshfH3XDA=
```

所以其中的`bAExy9pPp0EnzjAlYn1bsSEGvqYi1shl1OOshfH3XDA=`就是我们将要进行证书锁定的指纹(Hash)信息。

# 总结

证书锁定旨在解决移动端APP与服务端通信的唯一性，实际通信过程中，如果锁定过程失败，那么客户端APP将拒绝针对服务器的所有 SSL/TLS 请求，FaceBook/Twitter则通过证书锁定以防止Charles/Fiddler等抓包工具中间人攻击，关于如何在Android、IOS的各类开发工具中设置证书锁定，请参照《[Android SSL证书锁定(SSL/TLS Pinning)](https://link.zhihu.com/?target=https%3A//www.infinisign.com/faq/what-is-ssl-pinning%23)》、《[IOS SSL证书锁定(SSL/TLS Pinning)](https://link.zhihu.com/?target=https%3A//www.infinisign.com/faq/ios-ssl-pinning-three-method)》。
