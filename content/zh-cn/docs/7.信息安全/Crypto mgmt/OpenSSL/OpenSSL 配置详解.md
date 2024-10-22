---
title: OpenSSL 配置详解
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
default*ca = CA_default /\_The default ca section*/
\####################################################################
[ CA*default ]
dir = /etc/pki/CA /* Where everything is kept _/
 /_ #### 这是第一个 openssl 目录结构中的目录 _/
certs = $dir/certs /_ Where the issued certs are kept(已颁发的证书路径，即 CA 或自签的) _/
 /_ #### 这是第二个 openssl 目录结构中的目录，但非必须 _/
crl_dir = $dir/crl /_ Where the issued crl are kept(已颁发的 crl 存放目录) _/
 /_ #### 这是第三个 openssl 目录结构中的目录*/
database = $dir/index.txt /* database index file _/
\#unique_subject = no /_ 设置为 yes 则 database 文件中的 subject 列不能出现重复值 _/
 /_ 即不能为 subject 相同的证书或证书请求签名*/
 /* 建议设置为 no，但为了保持老版本的兼容性默认是 yes _/
new_certs_dir = $dir/newcerts /_ default place for new certs(将来颁发的证书存放路径) _/
 /_ #### 这是第四个 openssl 目录结构中的目录 _/
certificate = $dir/cacert.pem /_ The A certificate(CA 自己的证书文件) _/
serial = $dir/serial /_ The current serial number(提供序列号的文件)_/
crlnumber = $dir/crlnumber /_ the current crl number(当前 crl 序列号) _/
crl = $dir/crl.pem /_ The current CRL(当前 CRL) _/
private_key = $dir/private/cakey.pem /_ The private key(签名时需要的私钥，即 CA 自己的私钥) _/
RANDFILE = $dir/private/.rand /_ private random number file(提供随机数种子的文件) _/
x509_extensions = usr_cert /_ The extentions to add to the cert(添加到证书中的扩展项) _/
/_ 以下两行是关于证书展示格式的，虽非必须项，但推荐设置。一般就如下格式不用修改 _/
name_opt = ca_default /_ Subject Name options*/
cert_opt = ca_default /* Certificate field options _/
/_ 以下是 copy*extensions 扩展项，需谨慎使用*/
\# copy*extensions = copy /* 生成证书时扩展项的 copy 行为，可设置为 none/copy/copyall _/
 /_ 不设置该 name 时默认为 none _/
 /_ 建议简单使用时设置为 none 或不设置，且强烈建议不要设置为 copyall _/
\# crl_extensions = crl_ext
default_days = 365 /_ how long to certify for(默认的证书有效期) _/
default_crl_days= 30 /_ how long before next CRL(CRL 的有效期) _/
default_md = default /_ use public key default MD(默认摘要算法) _/
preserve = no /_ keep passed DN ordering(Distinguished Name 顺序，一般设置为 no _/
 /_ 设置为 yes 仅为了和老版本的 IE 兼容)_/
policy = policy_match /_ 证书匹配策略,此处表示引用\[ policy*match ]的策略*/
/*证书匹配策略定义了证书请求的 DN 字段(field)被 CA 签署时和 CA 证书的匹配规则*/
/*对于 CA 证书请求，这些匹配规则必须要和父 CA 完全相同*/
[ policy*match ]
countryName = match /* match 表示请求中填写的该字段信息要和 CA 证书中的匹配 _/
stateOrProvinceName = match
organizationName = match
organizationalUnitName = optional /_ optional 表示该字段信息可提供可不提供 _/
commonName = supplied /_ supplied 表示该字段信息必须提供 _/
emailAddress = optional
/_ For the 'anything' policy*/
/* At this point in time, you must list all acceptable 'object' types. _/
/_ 以下是没被引用的策略扩展，只要是没被引用的都是被忽略的 _/
[ policy_anything ]
countryName = optional
stateOrProvinceName = optional
localityName = optional
organizationName = optional
organizationalUnitName = optional
commonName = supplied
emailAddress = optional
/_ 以下是添加的扩展项 usr*cert 的内容*/
[ usr*cert ]
basicConstraints=CA:FALSE /* 基本约束，CA:FALSE 表示该证书不能作为 CA 证书，即不能给其他人颁发证书*/
/* keyUsage = critical,keyCertSign,cRLSign # 指定证书的目的，也就是限制证书的用法*/
/* 除了上面两个扩展项可能会修改下，其余的扩展项别管了，如下面的 \_/
nsComment = "OpenSSL Generated Certificate"
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid,issuer
```

## req Section

为 `openssl req` 命令提供运行时参数

```
[ req ]
default*bits = 2048 /* 生成证书请求时用到的私钥的密钥长度 _/
default_md = sha1 /_ 证书请求签名时的单向加密算法 _/
default_keyfile = privkey.pem /_ 默认新创建的私钥存放位置， _/
 /_ 如-new 选项没指定-key 时会自动创建私钥 _/
 /_ -newkey 选项也会自动创建私钥 _/
distinguished_name = req_distinguished_name /_ 可识别的字段名(常被简称为 DN) _/
 /_ 引用 req*distinguished_name 段的设置*/
x509*extensions = v3_ca /* 加入到自签证书中的扩展项 _/
\# req_extensions = v3_req /_ 加入到证书请求中的扩展项 _/
attributes = req_attributes /_ 证书请求的属性，引用 req*attributes 段的设置，可以不设置它*/
\# encrypt*key = yes | no /* 自动生成的私钥文件要加密否？一般设置 no，和-nodes 选项等价 _/
/_ 输入和输出私钥文件的密码，如果该私钥文件有密码，不写该设置则会提示输入 _/
/_ input*password = secret*/
/*output_password = secret*/
\# prompt = yes | no /*设置为 no 将不提示输入 DN field，而是直接从配置文件中读取，需要同时设置 DN 默认值，否则创建证书请求时将出错。*/
string*mask = utf8only
\[ req_distinguished_name ]
/* 以下项均可指定可不指定，但 ca 段的 policy 中指定为 match 和 supplied 一定要指定。 _/
/_ 以下选项都可以自定义，如 countryName = C，commonName = CN _/
countryName = Country Name (2 letter code) /_ 国家名(C) _/
countryName_default = XX /_ 默认的国家名 _/
countryName_min = 2 /_ 填写的国家名的最小字符长度 _/
countryName_max = 2 /_ 填写的国家名的最大字符长度 _/
stateOrProvinceName = State or Province Name (full name) /_ 省份(S) _/
/_ stateOrProvinceName*default = Default Province*/
localityName = Locality Name (eg, city) /*城市(LT)*/
localityName*default = Default City
0.organizationName = Organization Name (eg, company) /* 公司(ON) _/
0.organizationName_default = Default Company Ltd
organizationalUnitName = Organizational Unit Name (eg, section) /_ 部门(OU) _/
/_ organizationalUnitName*default =*/
/*以下的 commonName(CN)一般必须给,如果作为 CA，那么需要在 ca 的 policy 中定义 CN = supplied*/
/*CN 定义的是将要申请 SSL 证书的域名或子域名或主机名。*/
/*例如要为 zhonghua.com 申请 ssl 证书则填写 zhonghua.com，而不能填写www.zhonghua.com*/
/*要为www.zhonghua.com申请SSL则填写www.zhonghua.com*/
/*CN 必须和将要访问的网站地址一样，否则访问时就会给出警告*/
/*该项要填写正确，否则该请求被签名后证书中的 CN 与实际环境中的 CN 不对应，将无法提供证书服务*/
commonName = Common Name (eg, your name or your server's hostname) /*主机名(CN)*/
commonName*max = 64
emailAddress = Email Address /* Email 地址，很多时候不需要该项的 _/
emailAddress_max = 64
[ req_attributes ] /_ 该段是为了某些特定软件的运行需要而设定的， _/
 /_ 现在一般都不需要提供 challengepassword _/
 /_ 所以该段几乎用不上 _/
 /_ 所以不用管这段 _/
challengePassword = A challenge password
challengePassword_min = 4
challengePassword_max = 20
unstructuredName = An optional company name
[ v3_req ]
/_ Extensions to add to a certificate request _/
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
[ v3_ca ]
/_ Extensions for a typical CA _/
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid:always,issuer
basicConstraints = CA:true
\# keyUsage = cRLSign, keyCertSign /_ 典型的 CA 证书的使用方法设置，由于测试使用所以注释了 \_/
/*如果真的需要申请为 CA/*么该设置可以如此配置 \*/

可以自定义 DN(Distinguished Name)段中的字段信息，注意 ca 段中的 policy 指定的匹配规则中如果指定了 match 或这 supplied 的则 DN 中必须定义。例如下面的示例：由于只有 countryName、organizationName 和 commonName 被设定为 match 和 supplied，其余的都是 optional，所以在 DN 中可以只定义这 3 个字段，而且在 DN 中定义了自定义的名称。
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
\[ tsa ]
default_tsa = tsa_config1 # the default TSA section
\[ tsa_config1 ]
\# These are used by the TSA reply generation only.
dir = ./demoCA # TSA root directory
serial = $dir/tsaserial # The current serial number (mandatory)
crypto_device = builtin # OpenSSL engine to use for signing
signer_cert = $dir/tsacert.pem # The TSA signing certificate
 \# (optional)
certs = $dir/cacert.pem # Certificate chain to include in reply
 \# (optional)
signer_key = $dir/private/tsakey.pem # The TSA private key (optional)
signer_digest = sha256 # Signing digest to use. (Optional)
default_policy = tsa_policy1 # Policy if request did not specify it
 \# (optional)
other_policies = tsa_policy2, tsa_policy3 # acceptable policies (optional)
digests = sha1, sha256, sha384, sha512 # Acceptable message digests (mandatory)
accuracy = secs:1, millisecs:500, microsecs:100 # (optional)
clock_precision_digits = 0 # number of digits after dot. (optional)
ordering = yes # Is ordering defined for timestamps?
 \# (optional, default: no)
tsa_name = yes # Must the TSA name be included in the reply?
 \# (optional, default: no)
ess_cert_id_chain = no # Must the ESS cert id chain be included?
 \# (optional, default: no)
ess_cert_id_alg = sha1 # algorithm to compute certificate
 \# identifier (optional, default: sha1)
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
