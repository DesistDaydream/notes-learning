---
title: OpenSSL
---

# 概述

> 参考：
> - [官网](https://www.openssl.org/)
> - [GitHub 项目,openssl/openssl](https://github.com/openssl/openssl)

OpenSSL 是一个商业级且功能齐全的工具包，用于通用密码学和安全通信

OpenSSL 可以实现 **TLS(传输层安全性)** 和 **SSL(安全套接字层)** 协议的预期功能，类似于 OpenSSH 是 ssh 协议的实现

OpenSSL 主要包含两组东西：

- openssl # 多用途的命令行工具
- libraries # OpenSSL 库
  - libcrypto # 加密解密库
  - libssl # ssl 库，实现了 ssl 及 tls 的功能

# OpenSSL 关联文件

**/etc/ssl/openssl.conf** # OpenSSL 的“命令行工具”和 “库”默认使用的配置文件。

如果想要使用 CA 功能，需要进行如下配置

- touch /etc/pki/CA/index.txt
- echo 01 > /etc/pki/CA/serial

# openssl 命令行工具

> 参考：
> - [Manual(手册),openssl](https://www.openssl.org/docs/manmaster/man1/openssl.html)

openssl 程序提供了丰富的子命令，以实现 TLS/SSL 网络协议以及它们所需要的相关加密标准。

## Syntax(语法)

**openssl Command \[ OPTIONS ] \[ ARGUMENTS ]**

### Command

- Standard commands # 标准命令
  - asn1parse，ca，ciphers，cms，crl，crl2pkcs7，dgst，dh，dhparam，dsa，dsaparam，ec，ecparam，enc，engine，errstr，gendh，gendsa，genpkey，genrsa，nseq，ocsp，passwd，pkcs12，pkcs7，pkcs8，pkey，pkeyparam，pkeyutl，prime，rand，req，rsa，rsautl，s_client，s_server，s_time，sess_id，smime，speed，spkac，ts，verify，version，x509
- Message Digest commands #消息摘要命令，消息摘要算法的实现(用于单向加密)。使用 dgst 命令
  - md2，md4，md5，rmd160，sha，sha1
- Cipher commands #密码命令（其中都是各种加密算法，用于对称加密）。使用 enc 命令
  - aes-128-cbc，aes-128-ecb，aes-192-cbc，aes-192-ecb，aes-256-cbc，aes-256-ecb，base64，bf，bf-cbc，bf-cfb，bf-ecb，bf-ofb，camellia-128-cbc，camellia-128-ecb，camellia-192-cbc，camellia-192-ecb，camellia-256-cbc，camellia-256-ecb，cast，cast-cbc，cast5-cbc，cast5-cfb，cast5-ecb，cast5-ofb，des，des-cbc，des-cfb，des-ecb，des-ede，des-ede-cbc，des-ede-cfb，des-ede-ofb，des-ede3，des-ede3-cbc，des-ede3-cfb，des-ede3-ofb，des-ofb，des3，desx，idea，idea-cbc，idea-cfb，idea-ecb，idea-ofb，rc2，rc2-40-cbc，rc2-64-cbc，rc2-cbc，rc2-cfb，rc2-ecb，rc2-ofb，rc4，rc4-40，rc5，rc5-cbc，rc5-cfb，rc5-ecb，rc5-ofb，seed，seed-cbc，seed-cfb，seed-ecb，seed-ofb，zlib

### Global OPTIONS

- -in FILE # 指明使用的文件
- -out FILE # 指明输出的文件

## Standard commands #标准命令

[Standard commands(标准命令)](✏IT 学习笔记/🔐7.信息安全/管理工具/OpenSSL/Standard%20commands(标准命令).md commands(标准命令).md)

## Message Digest commands # 消息摘要命令

消息摘要算法的实现(用于单向加密)。使用 dgst 命令

## Cipher commands # 密码命令

其中都是各种加密算法，用于对称加密。使用 enc 命令

### openssl enc # 对称密钥程序，用于创建管理对称密钥

OPTIONS

- **-e** # 加密文件
- **-d** # 解密文件
- **-des3** # 使用 des3 算法进行加密或解密
- **-a** # 基于文本进行编码
- **-salt** # 加入一些盐

EXAMPLE

- openssl enc -e -des3 -a -salt -in fstab -out fstab.ciphertext # 加密 fstab 文件为 fstab.ciphertext，算法为 des3，基于文本进行编码，加入一些 salt
- openssl enc -d -des3 -a -salt -in fstab.ciphertext -out fstab # 解密 fstab.ciphertext 为 fstab 文件

# 应用实例：

## 创建自签 ca 证书

- (umask 077; openssl genrsa -out ca.key 2048)
- openssl req -new -x509 -key ca.key -days 3650 -out ca.crt

## 在 kubernetes 中生成个人证书

- 在当前目录下生成一个 2048 位的名为 lch.key 的私钥（括号的作用是创建子 shell 执行命令，这样 umask 命令对当前 shell 没影响）
  - (umask 077;openssl genrsa -out lch.key 2048)
- 使用 lck.key 进行证书申请
  - openssl req -new -key lch.key -out lch.csr -subj "/CN=lch"
- 使用 ca.key 来给 lch.crt 颁发证书，以生成 lch.crt 文件
  - openssl x509 -req -in lch.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out lch.crt -days 365
- 查看 ca.crt 证书的信息
  - openssl x509 -in lch.crt -text -noout

openssl x509 部分命令

打印出证书的内容：

openssl x509 -in cert.pem -noout -text

打印出证书的系列号

openssl x509 -in cert.pem -noout -serial

打印出证书的拥有者名字

openssl x509 -in cert.pem -noout -subject

以 RFC2253 规定的格式打印出证书的拥有者名字

openssl x509 -in cert.pem -noout -subject -nameopt RFC2253

在支持 UTF8 的终端一行过打印出证书的拥有者名字

openssl x509 -in cert.pem -noout -subject -nameopt oneline -nameopt -escmsb

打印出证书的 MD5 特征参数

openssl x509 -in cert.pem -noout -fingerprint

打印出证书的 SHA 特征参数

openssl x509 -sha1 -in cert.pem -noout -fingerprint

把 PEM 格式的证书转化成 DER 格式

openssl x509 -in cert.pem -inform PEM -out cert.der -outform DER

把一个证书转化成 CSR

openssl x509 -x509toreq -in cert.pem -out req.pem -signkey key.pem

给一个 CSR 进行处理，颁发字签名证书，增加 CA 扩展项

openssl x509 -req -in careq.pem -extfile openssl.cnf -extensions v3_ca -signkey key.pem -out cacert.pem

给一个 CSR 签名，增加用户证书扩展项

openssl x509 -req -in req.pem -extfile openssl.cnf -extensions v3_usr -CA cacert.pem -CAkey key.pem -CAcreateserial

查看 csr 文件细节：

openssl req -in my.csr -noout -text
