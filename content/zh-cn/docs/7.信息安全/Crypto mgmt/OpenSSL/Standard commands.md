---
title: Standard commands
linkTitle: Standard commands
weight: 20
---

# 概述

> 参考：

# 其他标准命令

## openssl passwd - 对指定的字符串生成 hash 过的密码

**openssl password \[OPTIONS] \[STRING]**

MD5 算法加密后的格式为：`$算法简称$SALT$XXXXXX`，算法简称为 1 或者 apr1，SALT 为指定的盐的字符串，XXXX 为生成的加密的字符串

OPTIONS

- -crypt # standard Unix password algorithm (default)
- -1 # -1 基于 MD5 的密码算法（注意：不指定 salt 的话，会使用随机的 slat）
- -salt # 在生成加密的密码中加盐(salt)。（为什么叫盐详见[https://zh.wikipedia.org/wiki/%E7%9B%90\_(%E5%AF%86%E7%A0%81%E5%AD%A6)](<https://zh.wikipedia.org/wiki/%25E7%259B%2590_(%25E5%25AF%2586%25E7%25A0%2581%25E5%25AD%25A6)>)这是密码学的一个概念）,加盐与不加盐得出的结果是不一样的

EXAMPLE

- openssl passwd -1 123456
  - 结果为：$1$ONQ8XSuX$Cv0wy2WbbbwOt/YkXuAlU/
- openssl passwd -1 -salt 123 123456
  - 结果为：$1$123$7mft0jKnzzvAdU4t0unTG1

## openssl rand - 生成伪随机数字节

EXAMPLE

- openssl rand -hex 6 # 生成随机数

# RSA 标准命令

## openssl genrsa - 生成 RSA 密钥

### Syntax(语法)

**openssl genrsa \[ OPTIONS ] \[ ARGUMENTS ]**

OPTIONS

EXAMPLE

- 在当前目录下生成一个 2048 位的名为 desistdaydream.key 的私钥（括号的作用是创建子 shell 执行命令，这样 umask 命令对当前 shell 没影响）
  - **(umask 077; openssl genrsa -out ./desistdaydream.key 2048)**

## openssl rsa - RSA 密钥管理

openssl rsa \[OPTIONS] \[ARGUMENTS]

OPTIONS

- **-noout** # 不输出密钥的编码格式
- **-text** # 除了编码后的格式，还会输出纯文本格式的内容，这些内容有公钥和私钥组件

EXAMPLE

- 从 desistdaydream.key 私钥中输出公钥信息，并将公钥信息写入到 lch.pub 文件中
  - **openssl rsa -in desistdaydream.key -pubout -out desistdaydream.pub**
- 显示 ca.key 密钥的信息
  - **openssl rsa -noout -text -in ca.key**

# 证书标准命令

## openssl ca - sample minimal CA application。CA 程序

**注意**：使用 `req` 和 `x509` 命令是非常精简的生成证书的方式。

EXAMPLE

- 在 CA 所在服务器使用 httpd.csr 的请求文件签署证书生成证书文件 httpd.crt，然后再把该证书文件，发送给请求方，整套流程就完成了。
  - 注意：如果想要执行该命令，需要注意为该服务器进行 CA 的配置，详见本章前面的"配置文件说明"
  - openssl ca -in httpd.csr -out httpd.crt -days 365
- openssl ca -revoke /etc/pki/CA/newcerts/SERIAL.pem # 吊销证书
- openssl ca -gencrl -out

## openssl req - 生成根证书或符合 PKCS #10 标准的证书请求

> 参考：
>
> - [Manual, openssl-req(1)](https://docs.openssl.org/master/man1/openssl-req/)

该命令主要以 PKCS#10 标准格式创建和处理 CSR。并且还可以创建自签名证书，以用作根 CA 证书。

### Syntax(语法)

**openssl req \[ OPTIONS ] \[ ARGUMENTS ]**

OPTIONS

- **-config \<FILENAME>** #
- **-new** # 生成新证书请求
- **-x509** # 生成自签证书
- **-key /PATH/FILE** # 用于生成请求时用到的私钥文件
- **-days NUM** # 证书的有效期
- **-text** # 以文本形式打印出证书
- **-noout** # 不输出证书的编码格式内容
- **-subj** # 在命令行设定 Subject 信息。可以避免交互式输入。Subject 符合 [X.509](docs/7.信息安全/Cryptography/公开密钥加密/证书%20与%20PKI/X.509.md#DN(distinguished%20names)) 的 “DN” 格式
    - 命令行参数的值的格式像这样: `/type0=value0/type1=value1/type2=...`。e.g. `-subj "/C=CN/CN=DesistDaydream-CA"`

### EXAMPLE

创建一个自建的 CA 根证书

```bash
openssl req -new -x509 -key /etc/pki/CA/private/ca.key -days 3650 -out /etc/pki/CA/ca.pem
```

这个命令之后还需要输入几个 Subject 信息，如下（输入完成后，自动生成 CA 自签证书）：

- Country Name (2 letter code) \[XX]:CN # 国家名，两个英文缩写
- State or Province Name (full name) \[]:Tianjin # 省或州名
- Locality Name (eg, city) \[Default City]:Tianjin # 地点名称，例如城市
- Organization Name (eg, company) \[Default Company Ltd]:GuanDian # 组织名称。例如公司
- Organizational Unit Name (eg, section) \[]:Ops # 组织单位名称，例如部门
- Common Name (eg, your name or your server's hostname) \[]:master0 # 通用名称，例如 CA 名或者服务器主机名。
  - Note：证书中的 CN 是很重要的标志，CN 可以使用主机名来表示，这样在使用证书来访问的时候，可以使用 CN 来验证域名是否可信。
  - 如果 CN 不使用主机名，则在签发证书的时候，需要 subjectAltName 字段来设定 DNS 别名，否则会报错提示证书对某些域名不可用。
  - 样例详见 harbor 使用私有证书部署：harbor 云原生注册中心.note
- Email Address # 邮箱地址

> Tips: 也可以使用 -subj 选项指定这些信息，避免交互式输入。

---

使用 httpd.key 这个密钥创建一个证书请求文件 httpd.csr

- openssl req -new -key httpd.key -days 365 -out httpd.csr
- 在输入完该命令后，同样需要输入几个身份信息以供 CA 进行验证。由于是私有 CA，所以所有信息应该保持跟 CA 的信息一样，具体信息详见上面那个命令，否则无法签署成功。后面还可以输入密码，当然密码也可以为空，密码主要是对改请求进行加密的。创建完请求后，把该请求文件 XXX.csr 发送给 CA 所在的服务器，然后由 CA 进行签署。

---

显示证书请求文件 desistdaydream.it.csr 的内容

- openssl req -in desistdaydream.it.csr -text -noout
- 等效于
- openssl x509 -req -in desistdaydream.it.csr -key desistdaydream.it.key -text -noout

## openssl x509 - 证书显示或签名工具

> 参考：
>
> - [Manual(手册), openssl-x509(1)](https://www.openssl.org/docs/manmaster/man1/openssl-x509.html)

这是一个多用途的证书处理命令。

- 可以用于打印证书、证书请求信息
- 将证书转换为各种形式
- 编辑证书信息设置
- 从头开始或从认证请求生成证书
- 等等

### Syntax(语法)

**openssl x509 \[ OPTIONS ] \[ ARGUMENTS ]**

#### Input, Output, and General Purpose OPTIONS(输入、输出、通用选项)

**-in**(Filename|URL) # 读取的输入源，输入源可以是本地文件或 URL。默认读取的输入源是 ”证书“ （i.e. X509v3 格式的文件）。

- Notes: 若使用了 -req 选项，则输入源必须是 ”证书请求“ （i.e. PKCS#10 格式的 XXX.csr 文件）
- 与 -new 选项互斥

**-req** # 改变 -in 选项的逻辑，输入便必须为 ”证书请求“ （i.e. PKCS#10 格式的 XXX.csr 文件），且该 csr 必须包含正确的自签名信息

- 用白话说：改变 -in 逻辑，其实就是将打印证书信息，改为打印证书请求信息、或者利用现存证书请求创建证书。所以 `openssl x509 -req` 与 `openssl req` 逻辑可能乍一看很相似，但是并不具备生成证书请求的能力，只不过都能显示证书请求内容罢了。而且显示证书请求内容时，不像 openssl-req 似的可以省略 key 文件。
- 默认情况下不会将 csr 内容中的扩展复制到证书中，需要使用 -extfile 指定要添加的扩展内容。

**-noout** # 禁止输出证书请求文件中的编码部分

**-pubkey**# 输出证书中的公钥

**-modulus**# 输出证书中公钥模块部分

**-serial** # 输出证书的序列号

**-subject** # 输出证书中的 subject

**-issuer**# 输出证书中的 issuer，即颁发者的 subject

**-subject_hash** # 输出证书中 subject 的 hash 码

**-issuer_hash** # 输出证书中 issuer(即颁发者的 subject)的 hash 码

**-hash** # 等价于"-subject_hash"，但此项是为了向后兼容才提供的选项

**-email**# 输出证书中的 email 地址，如果有 email 的话

**-startdate** # 输出证书有效期的起始日期

**-enddate** # 输出证书有效期的终止日期

**-dates** # 输出证书有效期，等价于"startdate+enddate"

**-fingerprint** # 输出指纹摘要信息

#### Certificate Printing OPTIONS(证书打印选项)

**-text** # 以 text 格式输出证书内容，即以最全格式输出、包括 public key,signature algorithms,issuer 和 subject names,serial number 以及 any trust settings.

#### Certificate Checking OPTIONS(证书检查选项)


#### Certificate Output OPTIONS(证书输出选项)

这部分的选项与生成证书有关

证书有效期相关的选项

- **-days**(INT) # 证书的有效天数（从签署证书那一刻开始计算 N 天）。`默认值: 30`。根据该选项会自动生成证书中的 `Not Before` 和 `Not After` 字段的值。

**-subj, -set_subject**(STRING) # 指定创建证书时的 Subject 信息。若证书为自签名，则 Issuer 名称默认与 Subject 名称相同，除非指定了 -set_issuer 选项。

**-extfile**(STRING) # 用于生成证书 X.509v3 扩展部分内容的配置文件。

- 配置文件的字段会如何生成证书中字段详见 [OpenSSL 配置详解](docs/7.信息安全/Crypto%20mgmt/OpenSSL/OpenSSL%20配置详解.md)

**-extensions**(SectionName) # 从 -extfile 选项指定的文件中指定 SectionName，该 Section 中的参数用来配置 X.509 证书的 `X509v3 extensions` 字段。

- 若未指定本选项，则默认使用 “默认部分” 的 extensions 指令。

#### Micro-CA OPTIONS

**-CA**(FileName|URI) # 指定用于签署证书的 CA 证书。

- 通常与 -req 选项一起使用
- 必须与 -CAkey 一起使用。若不指定 -CAkey，则 -CA 指定的输入源中必须包含密钥

**-CAkey**(FileName|URI) # 指定用于签署证书的 CA 私钥。

- 该选项必须与 -CA 标志一起使用。若不使用该选项，则 -CA 指定的输入源中必须包含密钥

**-CAcreateserial** # TODO: 具体效果未知。但是在签署证书的官方示例中，都带着这个选项

### Example

