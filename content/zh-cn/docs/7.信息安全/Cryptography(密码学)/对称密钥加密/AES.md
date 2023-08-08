---
title: "AES"
linkTitle: "AES"
date: "2023-08-08T23:25"
weight: 20
---

# 概述

> 参考：
> 
> - [Wiki，Advanced_Encryption_Standard](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard)

**Advanced Encryption Standard(高级加密标准，简称 AES)** 原名 **Rijndael**，是一种基于 **[Block cipher](docs/7.信息安全/Cryptography(密码学)/Block%20cipher.md)(块密码)** 变体算法的电子数据加密规范。

AES加密算法有以下应用：

- 无线安全：使用AES加密标准来验证路由器和设备之间的通信，防止数据被窃取或篡改。
- 加密浏览：AES在保护网站服务器认证方面发挥了巨大作用，可以在客户端和服务器端之间建立安全的连接。
- VPN服务：VPN用户可以使用AES加密算法来保护他们的网络流量，防止被监视或跟踪。VPN协议如OpenVPN也使用AES来加密数据包。

## 填充

您好，这是必应。我可以用中文回答您的问题。😊

AES（高级加密标准）是一种对称加密算法，它需要将明文分成固定长度的块（16字节）进行加密。如果明文的长度不是16字节的整数倍，就需要进行填充（padding）来补齐。PKCS7（公钥密码学标准7）是一种常用的填充方式，它的原理是：如果缺少n个字节，就在明文末尾添加n个n1。例如，如果明文是“Hello”，那么需要填充11个11，即“Hello\x0b\x0b\x0b\x0b\x0b\x0b\x0b\x0b\x0b\x0b\x0b”。如果明文刚好是16字节的整数倍，那么还需要再填充一个完整的块，即16个162 3 4。

因此，AES跟PKCS7有关系的地方就是：PKCS7是AES加密时常用的一种填充方式，它可以保证明文的长度总是16字节的整数倍，方便AES进行分块加密和解密1 2 3 4。

希望这个回答对您有帮助。如果您还有其他问题，请随时提问。😊