---
title: Standardized Glossary
linkTitle: Standardized Glossary
date: 2024-01-14T20:30
weight: 20
---

# 概述

> 参考：
> 
> - [RFC 4949(互联网安全术语)](https://datatracker.ietf.org/doc/html/rfc4949)

**Standardized Glossary(标准化术语)**

**Crypto(密码学)** # 算是 **Cryptographic** 的前缀简写吧

**Entity(实体)** # 是任何存在的东西（anything that exists） —— 即使 只在逻辑或概念上存在（even if only exists logically or conceptually）。 例如，

- 你的计算机是一个 entity，
- 你写的代码也是一个 entity，
- 你自己也是一个 entity，
- 你吃的杂粮饼也是一个 entity，
- 你六岁时见过的幽灵也是一个 entity —— 即使你妈告诉你幽灵并不存在，这只是你的臆想。
- 所以
  - CA 也算一个实体

**Identity(身份)** # 每个 entity（实体）都有一个 identity（身份）。 要精确定义这个概念比较困难，这么来说吧：identity 是使你之所以为你 （what makes you you）的东西。

- 具体到计算机领域，identity 通常用一系列属性来表示，描述某个具体的 entity， 这里的属性包括 group、age、location、favorite color、shoe size 等等。

**Identifier(身份标识符)** # Identifier 跟 identity 还不是一个东西：每个 identifier 都是一个唯一标识符， 也唯一地关联到某个有 identity 的 entity。

- 例如，我是 Mike，但 Mike 并不是我的 identity，而只是个 name —— 虽然二者在我们 小范围的讨论中是同义的。

**Claim(声明) & Authentication(认证)** #

- 一个 entity 能 claim（声明）说，它拥有某个或某些 name。
- 其他 entity 能够对这个 claim 进行认证（authenticate），以确认这份声明的真假。一般来说，认证的目的是确认某些 claim 的合法性。
- Claim 不是只能关联到 name，还可以关联到别的东西。例如，我能 claim 任何东西： my age, your age, access rights, the meaning of life 等等。

**Subscriber & CA & relying party (RP)** #

- 能作为一个证书的 subject 的 entity，称为 subscriber（证书 owner）或 end entity。对应地，subscriber 的证书有时也称为 end entity certificates 或 leaf certificates， 原因在后面讨论 certificate chains 时会介绍。
- CA（certificate authority，证书权威）是给 subscriber 颁发证书的 entity，是一种 certificate issuer（证书颁发者）。CA 的证书，通常称为 root certificate 或 intermediate certificate，具体取决于 CA 类型。
- Relying party 是 使用证书的用户（certificate user），它验证由 CA 颁发（给 subscriber）的证书是否合法。一个 entity 可以同时是一个 subscriber 和一个 relying party。 也就是说，单个 entity 既有自己的证书，又使用其他证书来认证 remote peers， 例如双向 TLS（mutual TLS，mTLS）场景。

**key(密钥)** # 在密码学中，是指某个用来完成加密、解密、完整性验证等密码学应用的秘密信息。对于加密算法，key 指定明文转换成密文；对于解密算法，key 指定密文转换成明文

- **Plaintext/Cleartext(明文)** # 在密码学中，明文是未加密的信息，可以供人类和计算机读取的信息
- **Ciphertext/Cyphertext(密文)** # 在密码学中，密文是明文通过加密算法计算后生成的人类或计算器无法读取的一种信息
  - cipher 和 cypher 是同一个意思，两种不同的拼写方法

**Key Generation(密钥生成)** # [详见 Wiki](https://en.wikipedia.org/wiki/Key_generation)。密钥一般都是各种程序根据指定算法生成的。

**Password(密码)** 与 **Key(密钥) 的区别** # 详见 [Wiki Key，Key vs Password](<https://en.wikipedia.org/wiki/Key_(cryptography)>)。

- 对于大多数计算机安全目的和大多数用户而言，“密钥”与“密码”（或“密码短语”）并不相同，尽管实际上可以将密码用作密钥。密钥和密码之间的主要实际区别在于，密码和密码旨在由人类用户生成，读取，记住和再现（尽管用户可以将这些任务委托给密码管理软件）。相反，密钥旨在由实现密码算法的软件使用，因此不需要人类可读性等。实际上，大多数用户在大多数情况下甚至都不知道其日常软件应用程序的安全组件正在使用代表他们的密钥。
- 如果 [密码](https://en.wikipedia.org/wiki/Password) 被用作加密密钥，然后在精心设计的密码系统就不会这样使用它自己。这是因为密码往往是人类可读的，因此可能不是特别强。作为补偿，一个好的加密系统将不使用\_密码作为密钥\_来执行主要的加密任务本身，而是充当[密钥派生功能](https://en.wikipedia.org/wiki/Key_derivation_function)（KDF）的输入。该 KDF 使用密码作为起点，然后它将从该起点本身生成实际的安全加密密钥。世代可以使用各种方法，例如添加盐和拉伸键。

**Encoding(编码)** # 将数据的原始格式，转换为便于存储的格式

**Decoding(解码)** # 将存储的数据转换为原始格式以便使用

**Encrypt(加密)** # 使用 Key(密钥) 对信息进行编码的过程。

**Decrypt(解密)** # 使用 Key(密钥) 对信息进行解码的过程

**Encoding(编码) 与 Encrypt(加密) 的区别**

- 编码使用公开的方案，将数据转换为另一种格式，便于维护数据与传播。任何人都可以使用相同的编码规范，解码数据
- 加密使用私密的方法，将数据抓换为另一种格式，着重于数据的保密。只有拥有相同 Key 的人才可以使用相同的加密规范，解密数据
- 总结：编码和加密都是对格式的一种转换，但是它们是有区别的。**编码是公开的**，比如 Base 64 编码，任何人都可以解码；**而加密则相反，你只希望自己或者特定的人才可以对内容进行解密。**

**Signature(签名)** # 非对称加密中，使用私钥进行数字签名的行为。

**Verifying(验证)** # 非对称加密中，使用公钥验证数字签名的行为。

**Digital Signature(数字签名)** # [详见 Wiki](https://en.wikipedia.org/wiki/Digital_signature)。用于检验数字消息或文件的真实性的数学方案

**Public Key Cryptography Standards(非对称密钥加密标准，简称 PKCS)** # [详见 Wiki](https://en.wikipedia.org/wiki/PKCS)。该标准指定了使用公开密钥加密技术时所应该遵守的标准

**Public Key Infrastructure(非对称密钥基础设施，简称 PKI)** # [详见 Wiki](https://en.wikipedia.org/wiki/Public_key_infrastructure)。一个包括硬件、软件、人员、策略和规程的集合，用来实现基于公钥密码体制的密钥和证书的产生、管理、存储、分发和撤销等功能。

**Secure Hash Algorithm(安全哈希算法，简称 SHA)**

**Personal Identification Number(个人识别码，简称 PIN)**