---
title: "Cipher"
linkTitle: "Cipher"
date: "2023-08-09T13:35"
weight: 1
---

# 概述

> 参考：
>
> - [Wiki，Cipher](https://en.wikipedia.org/wiki/Cipher)

**Cipher** 在密码学中表示一种用于执行加密或解密的 **algorithm(算法)**。Cipher 不太好翻译中文，password 是密码，Cipher 可以理解为一套**密码系统**，i.e. 一系列可以作为过程遵循的明确定义的步骤

> 注意：cipher 和 cypher 是同一个意思，两种不同的拼写方法

## Block cipher 与 Stream cipher

随着时代的发展，曾经对每个字节进行加密的方式不再显示，一个动辄几个 G 的文件，如果使用与明文相同的密钥进行加密，那么密钥也需要几个 G，这给密钥的分发造成了苦难，这时，可以将原始明文划分成多个长度相同的小块（也就是分组），然后使用和这些小块长度相同改的密钥依次和所有分组中的明文进行异或运算以进行加密，这就是早期的 [Block cipher](/docs/7.信息安全/Cryptography/Cipher/Block%20cipher.md) 加密算法。解密时，先对密文进行同样大小的分组，然后用相同的密钥和所有的密文块异或，再合并得到明文。

TODO: Stream cipher 算法是基于什么痛点出现的呢？

# 分类

> #密码学
