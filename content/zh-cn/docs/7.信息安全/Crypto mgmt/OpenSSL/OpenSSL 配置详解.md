---
title: OpenSSL 配置详解
linkTitle: OpenSSL 配置详解
weight: 20
---

# 概述

> 参考：
>
> - [Manual, 文件格式](https://docs.openssl.org/master/man5/)
>     - [Manual(手册), config(5)](https://www.openssl.org/docs/manmaster/man5/config.html)
>     - [Manual(手册), x509v3_config(5)](https://www.openssl.org/docs/manmaster/man5/x509v3_config.html)
> - [Manual(手册), openssl-req(1)](https://www.openssl.org/docs/man3.0/man1/openssl-req.html)-CONFIGURATION FILE FORMAT 部分
> - <https://www.cnblogs.com/f-ck-need-u/p/6091027.html>

OpenSSL 配置文件为 OpenSSL 库及其二进制程序提供运行时参数。这是一个 [INI](/docs/2.编程/无法分类的语言/INI.md) 格式的配置文件。

配置文件一共有三类格式：

- config # OpenSSL 通用配置格式
- fips_config # OpenSSL FIPS 配置格式
- x509v3_config # X.509 V3 证书扩展配置格式

# 配置文件格式

OpenSSL 配置文件的语法与 [INI](docs/2.编程/无法分类的语言/INI.md) 类似，但与常见的 INI 配置并不太一样（Section 用于定义场景，而不是定义某种类型的配置），并且扩展了很多能力：

- Section 本身的意义来源于 Openssl 命令行工具，各个 Section 是为各种场景服务的，甚至每个 Section 都可以有很多相同的 K/V 对。
    - 通常可以配合 -extensions 选项以便让命令读取哪个 Section 的内容
- **特定的 Section 的名字是有意义的**，比如 `[req]` Section 可以为 `openssl req` 命令提供参数，当执行 `openssl req` 命令时，会从默认配置文件的 `[req]` Section 获取配置参数，若没有，则再从 `默认` Section 获取参数。
- Section 中除了 **key/value pair(键值对)** 以外，还可以包括 **Directives(指令)**
- Section 中的 **Key/Value Pair 可以进行变量定义**，也可以引用变量。此时 Key 就是变量名，Value 就是变量的值。
  - 引用方式有 `$VAR` 或 `${VAR}` 两种，要想引用其他 Section 中的变量，则使用 `$SectionName::VAR` 或 `${SectionName::VAR}`

```bash
 # 这是默认 Section
 HOME = /temp
 configdir = $ENV::HOME/config

 [ section_one ]
 # Quotes permit leading and trailing whitespace
 any = " any variable name "
 other = A string that can \
 cover several lines \
 by including \\ characters
 message = Hello World\n

 [ section_two ]
 greeting = $section_one::message
```

- 配置文件中的第一部分是一个 `默认` Section。若要引用默认 Section 的变量，则使用 `ENV` 作为 Section 名
- 可以在 Section 中引用其他 Section。比如：

```bash
 # These must be in the default section
 config_diagnostics = 1
 # 引用 openssl_init 部分的配置
 openssl_conf = openssl_init

 [openssl_init]
 # 引用 oids 部分的配置
 oid_section = oids
 # 以此类推
 providers = providers
 alg_section = evp_properties
 ssl_conf = ssl_configuration
 engines = engines
 random = random

 [oids]
 ... new oids here ...

 [providers]
 ... provider stuff here ...

 [evp_properties]
 ... EVP properties here ...

 [ssl_configuration]
 ... SSL/TLS configuration properties here ...

 [engines]
 ... engine properties here ...

 [random]
 ... random properties here ...
```

## 使用 `@` 将多值放在独立的 Section 中

若某个字段的值有多个，可以使用 `@` 将所有值放在一个独立的 Section 中

比如：

```ini
basicConstraints = critical, CA:true, pathlen:1
```

等效于

```ini
[extensions]
basicConstraints = critical, @basic_constraints

[basic_constraints]
CA = true
pathlen = 1
```

OpenSSL 不支持在一个部分中多次出现相同字段，否则将仅识别最后一字段值。要指定多个值，需要附加一个数字标识符，如下所示：

```ini
subjectAltName = DNS:desistdaydream.it, DNS:*.desistdaydream.it
```

等效于

```ini
subjectAltName = @alt_names

[alt_names]
DNS.1 = desistdaydream.it
DNS.2 = *.desistdaydream.it
```

# 字段详解

**distinguished_name = \<SectionName>** # 生成证书或 CSR 时，如何配置 DN(专有名称)。

**req_extensions = \<SectionName>** # 要添加到 CSR 的扩展信息。

**basicConstraints = CA:FALSE** #

**keyUsage = nonRepudiation, digitalSignature, keyEncipherment** #

**extendedKeyUsage = clientAuth, serverAuth** #

## SAN 相关字段

https://docs.openssl.org/master/man5/x509v3_config/#subject-alternative-name

**subjectAltName**(STRING) # 设置 SAN，多个 SNA 以 `,` 分隔，每个 SAN 格式为 `Key:Value`。

- 也可以利用 `@` 将多值放在独立的 Section 中，此时需要遵循 OpenSSL 配置文件格式定义，详见上文 [配置文件格式](#配置文件格式)
- 可用的 SAN Key 有如下这些，分别对应 X509v3 扩展 SAN 的类型：

| 配置文件中的 Key | 对应生成 X509v3 扩展 SAN 类型 |
| :--------: | :-------------------: |
|    DNS     |       DNS Name        |
|     IP     |      IP Address       |
|    URI     |          URI          |
|   email    |     email Address     |
|    RID     |                       |
|  dirName   |                       |
| otherName  |                       |

### 简单示例

```ini
subjectAltName = @alt_names

[alt_names]
DNS.1 = desistdaydream.it
DNS.2 = *.desistdaydream.it
```

该配置文件生成的证书扩展信息如下：

```ini
        X509v3 extensions:
            X509v3 Subject Alternative Name: 
                DNS:desistdaydream.it, DNS:*.desistdaydream.it
```

# 配置示例

## 其他

```bash
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth, serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = abc
IP.1 = 1.1.1.1
```
