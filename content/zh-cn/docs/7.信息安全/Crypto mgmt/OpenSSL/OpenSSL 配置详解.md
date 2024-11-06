---
title: OpenSSL 配置详解
linkTitle: OpenSSL 配置详解
date: 2020-11-06T09:10:00
weight: 20
---

# 概述

> 参考：
>
> - [Manual(手册), config(5)](https://www.openssl.org/docs/manmaster/man5/config.html)
> - [Manual(手册), x509v3_config(5)](https://www.openssl.org/docs/manmaster/man5/x509v3_config.html)
> - [Manual(手册), openssl-req(1)](https://www.openssl.org/docs/man3.0/man1/openssl-req.html)-CONFIGURATION FILE FORMAT 部分
> - <https://www.cnblogs.com/f-ck-need-u/p/6091027.html>

OpenSSL 配置文件为 OpenSSL 库及其二进制程序提供运行时参数。这是一个 [INI](/docs/2.编程/无法分类的语言/INI.md) 格式的配置文件。

配置文件一共有三类格式：

- config # OpenSSL 通用配置格式
- fips_config # OpenSSL FIPS 配置格式
- x509v3_config # X.509 V3 证书扩展配置格式

OpenSSL 配置文件为 INI 格式的配置扩展了很多功能，并规定了一些新的规则：

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

- **特定的 Section 的名字是有意义的**，比如 `[req]` Section 可以为 `openssl req` 命令提供参数，当执行 `openssl req` 命令时，会从默认配置文件的 `[req]` Section 获取配置参数，若没有，则再从 `默认` Section 获取参数

# \[默认]

# \[req]

**distinguished_name = \<SectionName>** # 生成证书或 CSR 时，如何配置 DN(专有名称)。

**req_extensions = \<SectionName>** # 要添加到 CSR 的扩展信息。

## \[Distinguished_Name]

## \[Req_Extensions ]

> 关于 CSR 的扩展信息的格式，详见 [Manual(手册),x509v3_config(5)](https://www.openssl.org/docs/manmaster/man5/x509v3_config.html)

**basicConstraints = CA:FALSE** #

**keyUsage = nonRepudiation, digitalSignature, keyEncipherment** #

**extendedKeyUsage = clientAuth, serverAuth** #

**subjectAltName = \<SectionName>**#


### \[SubjectAltName]

DNS.1 = abc

IP.1 = 1.1.1.1

# 默认配置文件详解

默认的配置路径在 /etc/pki/tls/openssl.cnf，该文件主要设置了 CSR、签名、crl 相关的配置。为 `ca`、`req` 子命令提供信息。

该文件默认自带 4 个 Section：默认、ca、req、tsa

## 默认 Section

## ca Section

```
[ ca ]
default_ca  = CA_default        /*The default ca section*/
####################################################################
[ CA_default ]
dir     = /etc/pki/CA    /* Where everything is kept */
                         /*  #### 这是第一个openssl目录结构中的目录 */
certs       = $dir/certs /* Where the issued certs are kept(已颁发的证书路径，即CA或自签的) */
                         /* #### 这是第二个openssl目录结构中的目录，但非必须 */
crl_dir     = $dir/crl   /* Where the issued crl are kept(已颁发的crl存放目录) */
                         /*  #### 这是第三个openssl目录结构中的目录*/
database    = $dir/index.txt /* database index file */
#unique_subject = no     /* 设置为yes则database文件中的subject列不能出现重复值 */
                         /* 即不能为subject相同的证书或证书请求签名*/
                         /* 建议设置为no，但为了保持老版本的兼容性默认是yes */
new_certs_dir = $dir/newcerts /* default place for new certs(将来颁发的证书存放路径) */
                             /* #### 这是第四个openssl目录结构中的目录 */
certificate = $dir/cacert.pem  /* The A certificate(CA自己的证书文件) */
serial      = $dir/serial      /* The current serial number(提供序列号的文件)*/
crlnumber   = $dir/crlnumber   /* the current crl number(当前crl序列号) */
crl     = $dir/crl.pem         /* The current CRL(当前CRL) */
private_key = $dir/private/cakey.pem  /* The private key(签名时需要的私钥，即CA自己的私钥) */
RANDFILE    = $dir/private/.rand      /* private random number file(提供随机数种子的文件) */
x509_extensions = usr_cert  /* The extentions to add to the cert(添加到证书中的扩展项) */
/* 以下两行是关于证书展示格式的，虽非必须项，但推荐设置。一般就如下格式不用修改 */
name_opt    = ca_default        /* Subject Name options*/
cert_opt    = ca_default        /* Certificate field options */
/* 以下是copy_extensions扩展项，需谨慎使用 */
# copy_extensions = copy  /* 生成证书时扩展项的copy行为，可设置为none/copy/copyall */
                          /* 不设置该name时默认为none */
                          /* 建议简单使用时设置为none或不设置，且强烈建议不要设置为copyall */
# crl_extensions    = crl_ext
default_days    = 365   /* how long to certify for(默认的证书有效期) */
default_crl_days= 30    /* how long before next CRL(CRL的有效期) */
default_md  = default   /* use public key default MD(默认摘要算法) */
preserve    = no        /* keep passed DN ordering(Distinguished Name顺序，一般设置为no */
                        /* 设置为yes仅为了和老版本的IE兼容)*/
policy      = policy_match /* 证书匹配策略,此处表示引用[ policy_match ]的策略 */
/* 证书匹配策略定义了证书请求的DN字段(field)被CA签署时和CA证书的匹配规则 */
/* 对于CA证书请求，这些匹配规则必须要和父CA完全相同 */
[ policy_match ]
countryName = match     /* match表示请求中填写的该字段信息要和CA证书中的匹配 */
stateOrProvinceName = match
organizationName    = match
organizationalUnitName  = optional  /* optional表示该字段信息可提供可不提供 */
commonName      = supplied    /* supplied表示该字段信息必须提供 */
emailAddress        = optional
/* For the 'anything' policy*/
/* At this point in time, you must list all acceptable 'object' types. */
/* 以下是没被引用的策略扩展，只要是没被引用的都是被忽略的 */
[ policy_anything ]
countryName     = optional
stateOrProvinceName = optional
localityName        = optional
organizationName    = optional
organizationalUnitName  = optional
commonName      = supplied
emailAddress        = optional 
/* 以下是添加的扩展项usr_cert的内容*/
[ usr_cert ]
basicConstraints=CA:FALSE   /* 基本约束，CA:FALSE表示该证书不能作为CA证书，即不能给其他人颁发证书*/
/* keyUsage = critical,keyCertSign,cRLSign  # 指定证书的目的，也就是限制证书的用法*/
/* 除了上面两个扩展项可能会修改下，其余的扩展项别管了，如下面的 */
nsComment  = "OpenSSL Generated Certificate"
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid,issuer
```

## req Section

为 `openssl req` 命令提供运行时参数

```
[ req ]
default_bits    = 2048     /* 生成证书请求时用到的私钥的密钥长度 */
default_md      = sha1     /* 证书请求签名时的单向加密算法 */
default_keyfile = privkey.pem  /* 默认新创建的私钥存放位置， */
                               /* 如-new选项没指定-key时会自动创建私钥 */
                               /* -newkey选项也会自动创建私钥 */
distinguished_name  = req_distinguished_name /* 可识别的字段名(常被简称为DN) */
                                             /* 引用req_distinguished_name段的设置 */
x509_extensions = v3_ca       /* 加入到自签证书中的扩展项 */
# req_extensions = v3_req     /* 加入到证书请求中的扩展项 */
attributes  = req_attributes  /* 证书请求的属性，引用req_attributes段的设置，可以不设置它 */
# encrypt_key = yes | no /* 自动生成的私钥文件要加密否？一般设置no，和-nodes选项等价 */
/* 输入和输出私钥文件的密码，如果该私钥文件有密码，不写该设置则会提示输入 */
/* input_password = secret */
/* output_password = secret */
# prompt = yes | no /* 设置为no将不提示输入DN field，而是直接从配置文件中读取，需要同时设置DN默认值，否则创建证书请求时将出错。 */
string_mask = utf8only
[ req_distinguished_name ]
/* 以下项均可指定可不指定，但ca段的policy中指定为match和supplied一定要指定。 */
/* 以下选项都可以自定义，如countryName = C，commonName = CN */
countryName             = Country Name (2 letter code) /* 国家名(C) */
countryName_default     = XX /* 默认的国家名 */
countryName_min         = 2  /* 填写的国家名的最小字符长度 */
countryName_max         = 2  /* 填写的国家名的最大字符长度 */
stateOrProvinceName = State or Province Name (full name) /* 省份(S) */
/* stateOrProvinceName_default = Default Province */
localityName = Locality Name (eg, city) /* 城市(LT) */
localityName_default = Default City
0.organizationName  = Organization Name (eg, company) /* 公司(ON) */
0.organizationName_default  = Default Company Ltd
organizationalUnitName      = Organizational Unit Name (eg, section) /* 部门(OU) */
/* organizationalUnitName_default = */
/* 以下的commonName(CN)一般必须给,如果作为CA，那么需要在ca的policy中定义CN = supplied */
/* CN定义的是将要申请SSL证书的域名或子域名或主机名。 */
/* 例如要为zhonghua.com申请ssl证书则填写zhonghua.com，而不能填写www.zhonghua.com */
/* 要为www.zhonghua.com申请SSL则填写www.zhonghua.com */
/* CN必须和将要访问的网站地址一样，否则访问时就会给出警告 */
/* 该项要填写正确，否则该请求被签名后证书中的CN与实际环境中的CN不对应，将无法提供证书服务 */
commonName  = Common Name (eg, your name or your server\'s hostname) /* 主机名(CN) */
commonName_max  = 64
emailAddress            = Email Address /* Email地址，很多时候不需要该项的 */
emailAddress_max        = 64
[ req_attributes ] /* 该段是为了某些特定软件的运行需要而设定的， */
                   /* 现在一般都不需要提供challengepassword */
                   /* 所以该段几乎用不上 */
                   /* 所以不用管这段 */
challengePassword       = A challenge password
challengePassword_min   = 4
challengePassword_max   = 20
unstructuredName        = An optional company name
[ v3_req ]
/* Extensions to add to a certificate request */
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
[ v3_ca ]
/* Extensions for a typical CA */
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid:always,issuer
basicConstraints = CA:true
# keyUsage = cRLSign, keyCertSign  /* 典型的CA证书的使用方法设置，由于测试使用所以注释了 */
/* 如果真的需要申请为CA/*么该设置可以如此配置 */

可以自定义DN(Distinguished Name)段中的字段信息，注意ca段中的policy指定的匹配规则中如果指定了match或这supplied的则DN中必须定义。例如下面的示例：由于只有countryName、organizationName和commonName被设定为match和supplied，其余的都是optional，所以在DN中可以只定义这3个字段，而且在DN中定义了自定义的名称。
[policy_to_match]
countryName = match
stateOrProvinceName = optional
organizationName = match
organizationalUnitName = optional
commonName = supplied
emailAddress = optional
[DN]
countryName = "C"
organizationName = "O"
commonName = "Root CA"
```

## tas Section

```
[ tsa ]
default_tsa = tsa_config1   # the default TSA section
[ tsa_config1 ]
# These are used by the TSA reply generation only.
dir     = ./demoCA      # TSA root directory
serial      = $dir/tsaserial    # The current serial number (mandatory)
crypto_device   = builtin       # OpenSSL engine to use for signing
signer_cert = $dir/tsacert.pem  # The TSA signing certificate
                    # (optional)
certs       = $dir/cacert.pem   # Certificate chain to include in reply
                    # (optional)
signer_key  = $dir/private/tsakey.pem # The TSA private key (optional)
signer_digest  = sha256         # Signing digest to use. (Optional)
default_policy  = tsa_policy1       # Policy if request did not specify it
                    # (optional)
other_policies  = tsa_policy2, tsa_policy3  # acceptable policies (optional)
digests     = sha1, sha256, sha384, sha512  # Acceptable message digests (mandatory)
accuracy    = secs:1, millisecs:500, microsecs:100  # (optional)
clock_precision_digits  = 0 # number of digits after dot. (optional)
ordering        = yes   # Is ordering defined for timestamps?
                # (optional, default: no)
tsa_name        = yes   # Must the TSA name be included in the reply?
                # (optional, default: no)
ess_cert_id_chain   = no    # Must the ESS cert id chain be included?
                # (optional, default: no)
ess_cert_id_alg     = sha1  # algorithm to compute certificate
                # identifier (optional, default: sha1)
```

# 配置示例

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
