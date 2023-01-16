---
title: Standard commands(标准命令)
---

# 其他标准命令

## openssl passwd # 对指定的字符串生成 hash 过的密码

**openssl password \[OPTIONS] \[STRING]**
MD5 算法加密后的格式为：$算法简称$SALT$XXXXXX，算法简称为 1 或者 apr1，SALT 为指定的盐的字符串，XXXX 为生成的加密的字符串

OPTIONS

- -crypt # standard Unix password algorithm (default)
- -1 # -1 基于 MD5 的密码算法（注意：不指定 salt 的话，会使用随机的 slat）
- -salt # 在生成加密的密码中加盐(salt)。（为什么叫盐详见[https://zh.wikipedia.org/wiki/%E7%9B%90\_(%E5%AF%86%E7%A0%81%E5%AD%A6)](<https://zh.wikipedia.org/wiki/%25E7%259B%2590_(%25E5%25AF%2586%25E7%25A0%2581%25E5%25AD%25A6)>)这是密码学的一个概念）,加盐与不加盐得出的结果是不一样的

EXAMPLE

- openssl passwd -1 123456
  - 结果为：$1$ONQ8XSuX$Cv0wy2WbbbwOt/YkXuAlU/
- openssl passwd -1 -salt 123 123456
  - 结果为：$1$123$7mft0jKnzzvAdU4t0unTG1

## openssl rand # 生成伪随机数字节

EXAMPLE

- openssl rand -hex 6 # 生成随机数

# RSA 标准命令

## openssl genrsa # 生成 RSA 密钥

### Syntax(语法)

**openssl genrsa \[ OPTIONS ] \[ ARGUMENTS ]**

OPTIONS

EXAMPLE

- 在当前目录下生成一个 2048 位的名为 lch.key 的私钥（括号的作用是创建子 shell 执行命令，这样 umask 命令对当前 shell 没影响）
  - **(umask 077; openssl genrsa -out ./lch.key 2048)**

## openssl rsa # RSA 密钥管理

openssl rsa \[OPTIONS] \[ARGUMENTS]

OPTIONS

- **-noout** # 不输出密钥的编码格式
- **-text** # 除了编码后的格式，还会输出纯文本格式的内容，这些内容有公钥和私钥组件

EXAMPLE

- 从 lch.key 私钥中输出公钥信息，并将公钥信息写入到 lch.pub 文件中
  - **openssl rsa -in lch.key -pubout -out lch.pub**
- 显示 ca.key 密钥的信息
  - **openssl rsa -noout -text -in ca.key**

# 证书标准命令

## openssl ca # sample minimal CA application。CA 程序

**注意**：使用 `req` 和 `x509` 命令是非常精简的生成证书的方式。

EXAMPLE

- 在 CA 所在服务器使用 httpd.csr 的请求文件签署证书生成证书文件 httpd.crt，然后再把该证书文件，发送给请求方，整套流程就完成了。
  - 注意：如果想要执行该命令，需要注意为该服务器进行 CA 的配置，详见本章前面的"配置文件说明"
  - openssl ca -in httpd.csr -out httpd.crt -days 365
- openssl ca -revoke /etc/pki/CA/newcerts/SERIAL.pem # 吊销证书
- openssl ca -gencrl -out

## openssl req # 生成根证书或符合 PKCS #10 标准的证书请求

该命令主要以 PKCS#10 标准格式创建和处理 CSR。并且还可以创建自签名证书，以用作根 CA 证书。

### Syntax(语法)

**openssl req \[ OPTIONS ] \[ ARGUMENTS ]**
OPTIONS

- **-config <FILENAME>** #
- **-new **#生成新证书请求
- **-x509** # 生成自签证书
- **-key /PATH/FILE** #指明用于生成请求时用到的私钥文件
- **-days NUM** #指明证书的有效期
- **-text** # 以文本形式打印出证书
- **-noout** #不输出证书的编码格式内容

EXAMPLE

- 创建一个 ca 自签证书。
  - openssl req -new -x509 -key /etc/pki/CA/private/cakey.pem -days 3650 -out /etc/pki/CA/cacert.pem
  - 这个命令之后还需要输入几个身份信息，如下。输入完成后，自动生成 CA 自签证书
    - Country Name (2 letter code) \[XX]:CN # 指明国家名，两个英文缩写
    - State or Province Name (full name) \[]:Tianjin # 指明省或州名
    - Locality Name (eg, city) \[Default City]:Tianjin # 指明地点名称，例如城市
    - Organization Name (eg, company) \[Default Company Ltd]:GuanDian # 指明组织名称。例如公司
    - Organizational Unit Name (eg, section) \[]:Ops # 指明组织单位名称，例如部门
    - Common Name (eg, your name or your server's hostname) \[]:master0 # 指明通用名称，例如 CA 名或者服务器主机名。
      - Note：证书中的 CN 是很重要的标志，CN 可以使用主机名来表示，这样在使用证书来访问的时候，可以使用 CN 来验证域名是否可信。
      - 如果 CN 不使用主机名，则在签发证书的时候，需要 subjectAltName 字段来设定 DNS 别名，否则会报错提示证书对某些域名不可用。
      - 样例详见 harbor 使用私有证书部署：harbor 云原生注册中心.note
    - Email Address \[]:373406000@qq.com #指明邮箱地址
- 使用 httpd.key 这个密钥创建一个证书请求文件 httpd.csr
  - openssl req -new -key httpd.key -days 365 -out httpd.csr
  - 在输入完该命令后，同样需要输入几个身份信息以供 CA 进行验证。由于是私有 CA，所以所有信息应该保持跟 CA 的信息一样，具体信息详见上面那个命令，否则无法签署成功。后面还可以输入密码，当然密码也可以为空，密码主要是对改请求进行加密的。创建完请求后，把该请求文件 XXX.csr 发送给 CA 所在的服务器，然后由 CA 进行签署。

## openssl x509 # 证书显示或签名工具

> 参考：
> - [Manual(手册),openssl-x509(1)](https://www.openssl.org/docs/manmaster/man1/openssl-x509.html)

这是一个多用途的证书处理命令。

- 可以用于打印证书信息
- 将证书转换为各种形式
- 编辑证书信息设置
- 从头开始或从认证请求生成证书
- 等等

### Syntax(语法)

**openssl x509 \[ OPTIONS ] \[ ARGUMENTS ] **

#### Input, Output, and General Purpose OPTIONS(输入、输出、通用选项)

- **-noout** # 禁止输出证书请求文件中的编码部分
- **-pubkey **# 输出证书中的公钥
- **-modulus **# 输出证书中公钥模块部分
- **-serial** # 输出证书的序列号
- **-subject **# 输出证书中的 subject
- **-issuer **# 输出证书中的 issuer，即颁发者的 subject
- **-subject_hash** # 输出证书中 subject 的 hash 码
- **-issuer_hash** # 输出证书中 issuer(即颁发者的 subject)的 hash 码
- **-hash** # 等价于"-subject_hash"，但此项是为了向后兼容才提供的选项
- **-email **# 输出证书中的 email 地址，如果有 email 的话
- **-startdate** # 输出证书有效期的起始日期
- **-enddate** # 输出证书有效期的终止日期
- **-dates** # 输出证书有效期，等价于"startdate+enddate"
- **-fingerprint** # 输出指纹摘要信息

#### Certificate Printing OPTIONS(证书打印选项)

- **-text** # 以 text 格式输出证书内容，即以最全格式输出、包括 public key,signature algorithms,issuer 和 subject names,serial number 以及 any trust settings.

#### Certificate Checking OPTIONS(证书检查选项)

- **-days** #
- **-extfile <FILENAME>** #
- **-extensions <SectionName>** # 从 -extfile 选项指定的文件中指定 SectionName，该 Section 中的参数用来配置 X.509 证书的 `X509v3 extensions` 字段。
  - 若未指定本选项，则默认使用 “默认部分” 的 extensions 指令。

#### Certificate Output OPTIONS(证书输出选项)

#### Micro-CA OPTIONS

- **-CA \<FileName|URI>** # 指定用于签署证书的 CA 证书。
- **-CAkey \<FileName|URI>** # 指定用于签署证书的 CA 私钥。
