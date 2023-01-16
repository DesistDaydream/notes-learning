---
title: X.509
---

# 概述

> 参考：
> - [Wiki,X.509](https://en.wikipedia.org/wiki/X.509)
> - [RFC,5280](https://datatracker.ietf.org/doc/html/rfc5280)
> - [Arthurchiao 博客,\[译\] 写给工程师：关于证书（certificate）和公钥基础设施（PKI）的一切（SmallStep, 2018）](https://arthurchiao.art/blog/everything-about-pki-zh/)
> - [RFC 5280,Internet X.509 PKI 证书和 CRL 配置文件](https://datatracker.ietf.org/doc/html/rfc5280)

X.509 是[密码学](https://www.yuque.com/go/doc/33178368)里定义公钥证书格式的**标准**。X.509 格式的证书已应用在包括 TSL/SSL 在内的众多网络协议里，它是 HTTPS 的基础。

在大部分时候，人们提到证书而没有加额外的限定词时，通常都是指 X.509 v3 证书。

- 更准确的说，是 RFC 5280 中描述、 CA/Browser Forum [Baseline Requirements](https://cabforum.org/baseline-requirements-documents/)中进一步完善的 PKIX 变种。
- 也可以说，指的是浏览器理解并用来做 HTTPS 的那些证书。
- 也是那些具有通过 HTTP + TLS 协议交互的程序们所使用的证书

当然，全世界并不是只有 X.509 这一种格式，SSH 和 PGP 都有其各自的格式。

X.509 在 1988 年作为 ITU(国际电信联盟) X.500 项目的一部分首次标准化。 这是 telecom(通信) 领域的标准，想通过它构建一个 global telephone book(全球电话簿)。 虽然这个项目没有成功，但却留下了一些遗产，X.509 就是其中之一。如果查看 X.509 的证书，会看到其中包含了 locality、state、country 等信息， 之前可能会有疑问为什么为 web 设计的证书会有这些东西，现在应该明白了，因为 X.509 并不是为 web 设计的。
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/puor14/1635944301557-e8774c02-d1c8-4e0f-9f7a-a2c3a7180ce0.png)
图片补充：可以说 Subject 其实就是符合 [Distinguished Name(专有名称，简称 DN)](https://en.wikipedia.org/wiki/Lightweight_Directory_Access_Protocol#Directory_structure) 的规范，只不过 Subject 只是包含了 DN 中的部分字段罢了。也可以说，**Subject 是符合 X.509 标准的 DN。**

# X.509 证书的格式

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/puor14/1635931450920-fd8cad72-9ee7-476a-96ef-5e6ed60cc52b.png)
**Certificate: <OBJECT>**

- **Data: <OBJECT> **# 证书的数据
  - **Version: STRING** # 版本号
  - **Serial Number: STRING **# 序列号
  - **Signature Algorithm: STRING** # 签名算法
  - **Issuer: STRING** # 发行人名称，也就是这个证书的签发者。
  - **Validity: <OBJECT> **# 有效期
    - **Not Before: STRING **# 不能早于该日期。即证书从本日期开始生效
    - **Not After: STRING **# 不能晚于该日期。即证书到本日期为止失效
  - **Subject: STRING **# 主体信息。如何 X.509 规范的 Distinguished Name。
    - 对于 CA 证书来说， Subject 与 Issuer 的值相同。
  - **Subject Public Key Info: **# 主体的公钥信息
    - **Public Key Algorithm: STRING **# 公钥算法
    - **主体的公钥**
  - **Issuer Unique Identifier:** # 颁发者唯一身份信息（可选项）
  - **Subject Unique Identifier:** # 主体唯一身份信息（可选项）
  - **Extensions: <OBJECT>** # 扩展信息（可选项）
    - ......
    - **X509v3 Subject Alternative Name:** # SAN 信息。常用来作为该证书的名称。
    - ......
- **Signature Algorithm: STRING **# 证书签名算法
- 证书签名

证书的 Issuer 和证书的 Subject 用 X.509 DN 表示，DN 是由 RDN 构成的序列。RDN 用“属性类型=属性值”的形式表示。常用的属性类型名称以及简写如下：

| 属性类型名称             | 含义         | 简写 |
| ------------------------ | ------------ | ---- |
| Common Name              | 通用名称     | CN   |
| Organizational Unit name | 机构单元名称 | OU   |
| Organization name        | 机构名       | O    |
| Locality                 | 地理位置     | L    |
| State or province name   | 州/省名      | S    |
| Country                  | 国名         | C    |

## DN 与 SAN 命名的历史包袱

### DN (distinguished names)

历史上，X.509 使用 X.500 distinguished names (DN) 来命名证书的使用者（name the subject of a certificate），即 subscriber。 一个 DN 包含了一个 common name （对作者我来说，就是 “Mike Malone”），此外还可以包含 locality、country、organization、organizational unit 及其他一些东西（数字电话簿相关）。

- 没人理解 DN，它在互联网上也没什么意义。
- 应该避免使用 DN。如果真的要用，也要尽量保持简单。
- 无需使用全部字段，实际上，也不应该使用全部字段。
- common name 可能就是需要用到的全部字段了，如果你是一个 thrill seeker ，可以在用上一个 organization name。

PKIX 规定一个网站的 DNS hostname 应该关联到 DN common name。最近，CAB Forum 已 经废弃了这个规定，使整个 DN 字段变成可选的（Baseline Requirements, sections 7.1.4.2）。

### SAN (subject alternative name)

在 [RFC 5280 的 4.2.1.6 部分](https://tools.ietf.org/html/rfc5280#section-4.2.1.6)中，推荐的现代最佳实践是使用 **证书扩展中的 subject alternative name(SAN)** 来绑定证书中的 name。

常用的 SAN 有四种类型，绑定的都是广泛使用的名字：

- domain names (DNS)
- email addresse
- IP addresse
- URI

在我们讨论的上下文中，这些都是唯一的，而且它们能很好地映射到我们想识别的东西：

- email addresses for people
- domain names and IP addresses for machines and code,
- URIs if you want to get fancy

应该使用 SAN。
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/puor14/1638258551706-2b7a5b62-a093-4b12-8b34-7c6b9eefe49b.png)
注意，Web PKI 允许在一个证书内 bind 多个 name，name 也允许通配符。也就是说，

- 一个证书可以有多个 SNA，也可以有类似 \*.smallstep.com 这样的 SAN。
- 这对有多个域名的的网站来说很有用。

# 证书扩展名与编码

通常，为了便于传输，需要为证书进行编码。就好比传输 JSON 数据时，也需要对其进行编码，收到后再解码。

一般情况下，X.509 格式证书的原始数据，会使用 ASN.1 的 DER 进行编码，将编码后的二进制数据再使用根据 PEM 格式使用 Base64 编码，然后就会生成 PEM 格式的证书数据，实际上，所谓的 X.509 格式的文件，其实就是具有 CERTIFICATE 标志的 PEM 格式文件。

- 可以将 ASN.1 理解成 X.509 的 JSON
- 但是实际上，更像是 protobuf、thrift 或 SQL DDL。说白了就是通过一种算法，将人类可读的明文的证书编码成另一种便于传输的格式。

## OID

ASN.1 除了有常见的数据类型，如整形、字符串、集合、列表等， 还有一个**不常见但很重要的类型：OID**（object identifier，**对象标识符**）。

- OID **与 URI 有些像**，但比 URI 要怪。
- OID （在设计上）是**全球唯一标识符**。
- 在结构上，OID 是在一个 hierarchical namespace 中的一个整数序列（例如 2.5.4.3）。

可以用 OID 来 tag 一段数据的类型。例如，一个 string 本来只是一个 string，但可 以 tag 一个 OID 2.5.4.3，然后就**变成了一个特殊 string**：这是 **X.509 的通用名字（common name）** 字段。
![oids.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/puor14/1638343160689-8e109cf9-cb84-4a14-94fb-99421dab444c.png)

## 证书的扩展名

X.509 有多种常用的扩展名。不过其中的一些还用于其它用途，就是说具有这个扩展名的文件可能并不是证书，比如说可能只是保存了私钥。

- .pem # 将使用 DER 格式编码的内容，再通过 PEM 进行 Base64 编码，得出来的数据存放在"-----BEGIN CERTIFICATE-----"和"-----END CERTIFICATE-----"之中
- .cer, .crt, .der – 通常是 DER 二进制格式的，但 Base64 编码后也很常见。
- .p7b, .p7c – PKCS#7 SignedData structure without data, just certificate(s) or CRL(s)
- .p12 – PKCS#12 格式，包含证书的同时可能还有带密码保护的私钥
- .pfx – PFX，PKCS#12 之前的格式（通常用 PKCS#12 格式，比如那些由 IIS 产生的 PFX 文件）

更多的编码信息见：
[密钥/证书 的编码](✏IT 学习笔记/🔐7.信息安全/Cryptography(密码学)/公开密钥加密/密钥\_证书%20 的编码.md 的编码.md)

# 证书示例

### 终端实体证书

这是 wikipedia.org 和其他几个维基百科网站使用的解码 X.509 证书的示例。它由 [GlobalSign](https://en.wikipedia.org/wiki/GlobalSign) 颁发，如 `Issuer` 字段中所述。它的 `Subject` 字段将维基百科描述为一个组织，它的 `Subject Alternative Name` 字段描述了可以使用它的域名。主题公钥信息字段包含一个[ECDSA](https://en.wikipedia.org/wiki/ECDSA)公钥，而底部的签名是由 GlobalSign 的[RSA](<https://en.wikipedia.org/wiki/RSA_(cryptosystem)>)私钥生成的。

```bash
Certificate:
    Data:
        # 证书版本
        Version: 3 (0x2)
        # 序列号
        Serial Number:
            10:e6:fc:62:b7:41:8a:d5:00:5e:45:b6
        # 证书的签名算法
        Signature Algorithm: sha256WithRSAEncryption
        # 证书的颁发者信息。CA 证书的 Issuer 与 Subject 相同
        Issuer: C=BE, O=GlobalSign nv-sa, CN=GlobalSign Organization Validation CA - SHA256 - G2
        # 证书有效期
        Validity
            Not Before: Nov 21 08:00:00 2016 GMT
            Not After : Nov 22 07:59:59 2017 GMT
        # 证书主体信息。i.e.该证书颁发给谁
        Subject: C=US, ST=California, L=San Francisco, O=Wikimedia Foundation, Inc., CN=*.wikipedia.org
        # 证书公钥信息
        Subject Public Key Info:
            # 证书主体的公钥算法
            Public Key Algorithm: id-ecPublicKey
                Public-Key: (256 bit)
            pub:
                    00:c9:22:69:31:8a:d6:6c:ea:da:c3:7f:2c:ac:a5:
                    af:c0:02:ea:81:cb:65:b9:fd:0c:6d:46:5b:c9:1e:
                    9d:3b:ef
                ASN1 OID: prime256v1
                NIST CURVE: P-256
        # X509 v3版本的扩展信息
        X509v3 extensions:
            # 密钥用法：critical级别。包括数字签名、密钥加密
            X509v3 Key Usage: critical
                Digital Signature, Key Agreement
            Authority Information Access:
                CA Issuers - URI:http://secure.globalsign.com/cacert/gsorganizationvalsha2g2r1.crt
                OCSP - URI:http://ocsp2.globalsign.com/gsorganizationvalsha2g2
            X509v3 Certificate Policies:
                Policy: 1.3.6.1.4.1.4146.1.20
                  CPS: https://www.globalsign.com/repository/
                Policy: 2.23.140.1.2.2
            X509v3 Basic Constraints:
                CA:FALSE
            X509v3 CRL Distribution Points:
                Full Name:
                  URI:http://crl.globalsign.com/gs/gsorganizationvalsha2g2.crl
            # 证书主体名称的替代名称。i.e.别名
            X509v3 Subject Alternative Name:
                DNS:*.wikipedia.org, DNS:*.m.mediawiki.org, DNS:*.m.wikibooks.org, DNS:*.m.wikidata.org, DNS:*.m.wikimedia.org, DNS:*.m.wikimediafoundation.org, DNS:*.m.wikinews.org, DNS:*.m.wikipedia.org, DNS:*.m.wikiquote.org, DNS:*.m.wikisource.org, DNS:*.m.wikiversity.org, DNS:*.m.wikivoyage.org, DNS:*.m.wiktionary.org, DNS:*.mediawiki.org, DNS:*.planet.wikimedia.org, DNS:*.wikibooks.org, DNS:*.wikidata.org, DNS:*.wikimedia.org, DNS:*.wikimediafoundation.org, DNS:*.wikinews.org, DNS:*.wikiquote.org, DNS:*.wikisource.org, DNS:*.wikiversity.org, DNS:*.wikivoyage.org, DNS:*.wiktionary.org, DNS:*.wmfusercontent.org, DNS:*.zero.wikipedia.org, DNS:mediawiki.org, DNS:w.wiki, DNS:wikibooks.org, DNS:wikidata.org, DNS:wikimedia.org, DNS:wikimediafoundation.org, DNS:wikinews.org, DNS:wikiquote.org, DNS:wikisource.org, DNS:wikiversity.org, DNS:wikivoyage.org, DNS:wiktionary.org, DNS:wmfusercontent.org, DNS:wikipedia.org
            # 扩展密钥用法
            X509v3 Extended Key Usage:
                TLS Web Server Authentication, TLS Web Client Authentication
            X509v3 Subject Key Identifier:
                28:2A:26:2A:57:8B:3B:CE:B4:D6:AB:54:EF:D7:38:21:2C:49:5C:36
            X509v3 Authority Key Identifier:
                keyid:96:DE:61:F1:BD:1C:16:29:53:1C:C0:CC:7D:3B:83:00:40:E6:1A:7C
    # 证书的签名算法及其标识符
    Signature Algorithm: sha256WithRSAEncryption
         8b:c3:ed:d1:9d:39:6f:af:40:72:bd:1e:18:5e:30:54:23:35:
         ...
```

### 中级证书

```bash
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            04:00:00:00:00:01:44:4e:f0:42:47
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: C=BE, O=GlobalSign nv-sa, OU=Root CA, CN=GlobalSign Root CA
        Validity
            Not Before: Feb 20 10:00:00 2014 GMT
            Not After : Feb 20 10:00:00 2024 GMT
        Subject: C=BE, O=GlobalSign nv-sa, CN=GlobalSign Organization Validation CA - SHA256 - G2
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                Public-Key: (2048 bit)
                Modulus:
                    00:c7:0e:6c:3f:23:93:7f:cc:70:a5:9d:20:c3:0e:
                    ...
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Key Usage: critical
                Certificate Sign, CRL Sign
            X509v3 Basic Constraints: critical
                CA:TRUE, pathlen:0
            X509v3 Subject Key Identifier:
                96:DE:61:F1:BD:1C:16:29:53:1C:C0:CC:7D:3B:83:00:40:E6:1A:7C
            X509v3 Certificate Policies:
                Policy: X509v3 Any Policy
                  CPS: https://www.globalsign.com/repository/

            X509v3 CRL Distribution Points:

                Full Name:
                  URI:http://crl.globalsign.net/root.crl

            Authority Information Access:
                OCSP - URI:http://ocsp.globalsign.com/rootr1

            X509v3 Authority Key Identifier:
                keyid:60:7B:66:1A:45:0D:97:CA:89:50:2F:7D:04:CD:34:A8:FF:FC:FD:4B

    Signature Algorithm: sha256WithRSAEncryption
         46:2a:ee:5e:bd:ae:01:60:37:31:11:86:71:74:b6:46:49:c8:
         ...
```

### 根证书

```bash
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            04:00:00:00:00:01:15:4b:5a:c3:94
        Signature Algorithm: sha1WithRSAEncryption
        Issuer: C=BE, O=GlobalSign nv-sa, OU=Root CA, CN=GlobalSign Root CA
        Validity
            Not Before: Sep  1 12:00:00 1998 GMT
            Not After : Jan 28 12:00:00 2028 GMT
        Subject: C=BE, O=GlobalSign nv-sa, OU=Root CA, CN=GlobalSign Root CA
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                Public-Key: (2048 bit)
                Modulus:
                    00:da:0e:e6:99:8d:ce:a3:e3:4f:8a:7e:fb:f1:8b:
                    ...
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Key Usage: critical
                Certificate Sign, CRL Sign
            X509v3 Basic Constraints: critical
                CA:TRUE
            X509v3 Subject Key Identifier:
                60:7B:66:1A:45:0D:97:CA:89:50:2F:7D:04:CD:34:A8:FF:FC:FD:4B
    Signature Algorithm: sha1WithRSAEncryption
         d6:73:e7:7c:4f:76:d0:8d:bf:ec:ba:a2:be:34:c5:28:32:b5:
         ...
```

# X.509 格式证书数据示例

## PEM 格式数据

从 `-----BEGIN CERTIFICATE-----` 开始到 `-----END CERTIFICATE-----` 为止是证书的明文格式经过 ASN.1 编码再经过 Base64 编码得到的。

对于私钥文件，真正的私钥是包含在字符串 `-----BEGIN PRIVATE KEY-----` 和字符串 `-----END PRIVATE KEY-----` 之间。

## 原始数据

```bash
[root@master-1 pki]# openssl x509 -text -noout -in apiserver.crt
Certificate:
    Data:
        Version: 3 (0x2)                            # 证书版本
        Serial Number: 0 (0x0)                      # 序列号
    Signature Algorithm: sha256WithRSAEncryption    # 证书的签名算法
        Issuer: CN=kubernetes                       # 证书的颁发者信息。CA证书的Issuer与Subject相同
        Validity                                    # 证书有效期
            Not Before: Nov 20 08:45:23 2019 GMT
            Not After : Nov 17 08:45:23 2029 GMT
        Subject: CN=kube-apiserver                  # 证书主体信息。i.e.该证书颁发给谁
        Subject Public Key Info:                    # 证书公钥信息
            Public Key Algorithm: rsaEncryption     # 证书主体的公钥算法
                Public-Key: (2048 bit)
                Modulus:
                    00:c7:2e:02:61:db:b0:24:db:22:aa:46:94:de:7e:
                    .......
                Exponent: 65537 (0x10001)
        X509v3 extensions:                          # X509 v3版本的扩展信息
            X509v3 Key Usage: critical              # 密钥用法：critical级别。包括数字签名、密钥加密
                Digital Signature, Key Encipherment
            X509v3 Extended Key Usage:              # 扩展密钥用法
                TLS Web Server Authentication
            X509v3 Subject Alternative Name:        # 证书主体名称的替代名称。i.e.别名
                DNS:master-1.k8s.cloud.tjiptv.net, DNS:kubernetes, DNS:kubernetes.default, DNS:kubernetes.default.svc, DNS:kubernetes.default.svc.cluster.local, IP Address:10.96.0.1, IP Address:10.10.9.51, IP Address:10.10.9.54
    Signature Algorithm: sha256WithRSAEncryption    # 证书的签名算法及其标识符
         47:38:42:cf:02:85:71:49:ac:19:9c:ba:3a:f3:74:c3:4b:09:
         .....
```
