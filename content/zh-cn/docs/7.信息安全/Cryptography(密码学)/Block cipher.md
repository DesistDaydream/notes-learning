---
title: "Block cipher"
linkTitle: "Block cipher"
date: "2023-08-08T15:52"
weight: 20
---

# 概述

> 参考：
> 
> - [Wiki，Block_cipher](https://en.wikipedia.org/wiki/Block_cipher)
> - [WIki，Block_cipher_mode_of_operation](https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation)

**Block cipher(分组密码)** 是一种确定性算法，对固定长度的 bits（也称为 blcok）进行操作。

与 Block cipher 相关的另一种算法，是 [Stream cipher](docs/7.信息安全/Cryptography(密码学)/Stream%20cipher.md)

# 加密模式

- ECB # 最简单，不再推荐使用。Go 的标准库中，没有 ECB 模式的实现方式。
- CBC
- CFB
- 等等

## ECB

https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Electronic_codebook_(ECB)

**Electronic codebook(电子密码本，简称 ECB)** 是最简单（不再使用）的加密模式是模式（以传统物理密码本命名）。消息被分成块，每个块都单独加密。

# 分类

#密码学 