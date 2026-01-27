---
title: "Block cipher"
linkTitle: "Block cipher"
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, Block_cipher](https://en.wikipedia.org/wiki/Block_cipher)
> - [Wiki, Block_cipher_mode_of_operation](https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation)
> - [B 站，【计算机博物志】DES的生与死](https://www.bilibili.com/video/BV1qW4y1L7tN)

**Block cipher(分组密码)** 是一种确定性算法，对固定长度的 bits（也称为 blcok）进行操作。与 Block cipher 相关的另一种算法，是 [Stream cipher](/docs/7.信息安全/Cryptography/Cipher/Stream%20cipher.md)。

Block cipher 的基本逻辑：将待加密的明文，分为固定 bits 的块，比如 AES 的 128 bit block，然后使用相同长度的密钥（或指定倍数的长度），对每个 block 进行加密。

Block cipher 算法设计了多种模式以适应各种加密需求：

- ECB # 最简单的模式，不需要 IV
- CBC
- CFB
- CTR
- 等等

对于每个加密操作，大多数模式都需要唯一的二进制序列，通常称为初始化向量 (IV)。 IV 必须是非重复的，并且对于某些模式来说，也是随机的。初始化向量用于确保即使使用相同的密钥独立地对相同的明文进行多次加密，也会生成不同的密文。 分组密码可能能够对不止一种块大小进行操作，但在转换过程中块大小始终是固定的。分组密码模式对整个块进行操作，如果数据的最后部分小于当前块大小，则要求将其填充到整个块。然而，有些模式不需要填充，因为它们有效地使用分组密码作为 [Stream cipher](/docs/7.信息安全/Cryptography/Cipher/Stream%20cipher.md)。

在 Block cipher 算法中，首先要明确几个概念

- [Padding(填充)](#padding)
- [Initialization vector((初始化向量，简称 IV)](#iv)
  - ECB 模式不需要 IV

## Padding

Block cipher 算法在加密时将明文分成固定长度的块（128 bit，i.e. 16 Bytes）进行加密。如果明文的长度不是16 Bytes 的整数倍，就需要进行 **Padding(填充)** 来补齐。

填充方式：

- **NoPadding** # 不填充，要求明文是 16 Bytes 的整数倍。
- **PKCS5Padding/PKCS7Padding** # 按照 `PKCS#5` 或 `PKCS#7` 的标准进行填充，填充值为缺少的字节数，例如缺少 5 个字节，就填充 `05 05 05 05 05`。这两种方式在 [AES](/docs/7.信息安全/Cryptography/对称密钥加密/AES.md) 中没有区别。
  - 与非对称加密标准中的 PKCS 好像不是一个概念？？
- **ISO10126Padding** # 除了最后一个字节表示缺少的字节数外，其他字节用随机数填充，例如缺少 5 个字节，就填充 `XX XX XX XX 05`（XX 表示随机数）。
- **ZeroPadding** # 用 0 填充，如果明文刚好是 16 字节的整数倍，就再添加一个分组的 0。
- **其他** # 还有一些其他的填充方式，如 ISO7816-4Padding、X923Padding、TBCPadding 等，具体可以参考网上的资料。

不同的填充方式会影响加密和解密的效率和安全性，一般推荐使用 PKCS5Padding/PKCS7Padding。

## IV

**Initialization vector(初始化向量，简称 IV)** 是一个 block of bits，除了 ECB 以外的很多模式都是用 IV 来随机化加密（就像 Hash 加盐似的），从而让相同的密码文可以生成不同的密文。

# 加密模式

- ECB # 最简单，不再推荐使用。Go 的标准库中，没有 ECB 模式的实现方式。
- CBC
- CFB
- 等等

## ECB

https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Electronic_codebook_(ECB)

**Electronic codebook(电子密码本，简称 ECB)** 是最简单（不再使用）的加密模式是模式（以传统物理密码本命名）。消息被分成块，每个块都单独加密。

## CBC

https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Cipher_block_chaining_(CBC)

**Cipher block chaining(密码块链)**

