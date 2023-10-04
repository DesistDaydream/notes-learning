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

**Advanced Encryption Standard(高级加密标准，简称 AES)** 原名 **Rijndael**，是一种基于 **[Block cipher](docs/7.信息安全/Cryptography/Cipher/Block%20cipher.md)(块密码)** 变体算法的电子数据加密规范。取代了 1977 年发布的 DES(数据加密标准)

AES 规范中，将待加密的明文分为多个块，每块 128 bit，可以使用具有三种不同的密钥长度：128 bit、192 bit 和 256 bit，若我们提供的密钥不足长度，则可以自定规则将密钥补足，若超过长度，则截断超过的部分。

> - 这三种密钥长度换算为字节的话，分别是 16 Byte、24 Byte、32 Byte，也就是说，一般设置 16、24、32 个字符刚刚好。
> - 密钥的不足规则，可以是直接在末尾补 0、补充随机数、甚至可以直接生成 16、24、32 个随机字符。只要满足密钥长度即可，剩下的由个人掌握。
> - 其实常见的做法是丢进哈希函数里，比如你 rar 加密用的 AES-256。在加密之前把你给的密码进行 sha256 哈希（输出刚好是 256 bit），然后当成 AES-256 的密钥使用

## 描述用语

AES 128 ECB PKCS7 表示采用 128 bit 长度的密钥进行 AES 加密；使用 Block cipher 算法中的 ECB 模式，填充块的标准为 PKCS7，没有初始化向量
