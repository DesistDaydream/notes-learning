---
title: "crypto 库"
linkTitle: "crypto 库"
date: "2023-08-08T23:19"
weight: 20
---

# 概述

> 参考：
>
> - [Go crypto 库](https://pkg.go.dev/crypto)
> - [Go crypto/rsa 库](https://pkg.go.dev/crypto/rsa)
> - [Go crypto/x509 库](https://pkg.go.dev/crypto/x509)
> - [Go encoding/pem 库](https://pkg.go.dev/encoding/pem)
> - [简书](https://www.jianshu.com/p/54c20cc008c5)

crypto 标准库中主要有 cipher、rand、aes、rsa、x509、等等用于密码学加密的子包。

- cipher 包实现了标准的的 [Block cipher](/docs/7.信息安全/Cryptography/Cipher/Block%20cipher.md)、[Stream cipher](/docs/7.信息安全/Cryptography/Cipher/Stream%20cipher.md)，可以围绕低级分组密码实现进行包装
- rand 包 # 顾名思义，用来生成随机数的
- aes 包 # 实现了 [AES](/docs/7.信息安全/Cryptography/对称密钥加密/AES.md) 加密
- rsa 包 # 实现 `PKCS #1` 和 RFC 8017中指定的 [RSA](/docs/7.信息安全/Cryptography/公开密钥加密/RSA/RSA.md) 加密
- x509 包  # 实现了 X.509 标准的一个子集

# AES包

依赖 cipher 包，需要实例化一个 [Block cipher](/docs/7.信息安全/Cryptography/Cipher/Block%20cipher.md)，然后使用这个 cipher.Block 下的方法进行加密/解密。

> 由于国际组织不推荐使用 ECB 的方式，所以该库没有实现 ECB 模式的方法，需要自己手动写函数。
>
# RSA 包

rsa 库大体分为 加密/解密 与 签名/验签 这两大类。

- **rsa.EncryptOAEP() 等** # 加密方法
- **rsa.DecryptOAEP() 等** # 解密方法
- **rsa.SignPKCS1v15() 等** # 签名方法
- **rsa.VerifyPKCS1v15() 等** # 验签方法

## 生成密钥

可以通过 **rsa.GenerateKey()** 函数来生成一个密钥对。私钥的类型为 `*rsa.PrivateKey`，其中公钥在该私钥当中，类型是 `*rsa.PublicKey`。示例如下

```go
// RSA 是公钥和私钥两个组成一组的密钥对
type RSA struct {
 rsaPrivateKey *rsa.PrivateKey
 rsaPublicKey  *rsa.PublicKey
}

// NewRSA 生成密钥对
func NewRSA() *RSA {
 // 随机生成一个给定大小的 RSA 密钥对。可以使用 crypto 包中的 rand.Reader 来随机。
 privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
 // 从私钥中，获取公钥
 publicKey := privateKey.PublicKey
 return &RSA{
  rsaPrivateKey: privateKey,
  rsaPublicKey:  &publicKey,
 }
}
```

## 加密与解密

```go
// RSAEncrypt 使用 RSA 算法，加密指定明文
func (r *RSA) RSAEncrypt(plaintext []byte) []byte {
 // 使用公钥加密 plaintext(明文，也就是准备加密的消息)。并返回 ciphertext(密文)
 // 其中 []byte("DesistDaydream") 是加密中的标签，解密时标签需与加密时的标签相同，否则解密失败
 ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, r.rsaPublicKey, plaintext, []byte("DesistDaydream"))
 if err != nil {
  panic(err)
 }
 return ciphertext
}

// RSADecrypt 使用 RSA 算法，解密指定密文
func (r *RSA) RSADecrypt(ciphertext []byte) []byte {
 // 使用私钥解密 ciphertext(密文，也就是加过密的消息)。并返回 plaintext(明文)
 // 其中 []byte("DesistDaydream") 是加密中的标签，解密时标签需与加密时的标签相同，否则解密失败
 plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, r.rsaPrivateKey, ciphertext, []byte("DesistDaydream"))
 if err != nil {
  panic(err)
 }
 return plaintext
}
```

## 签名与验签

```go
// RSASign RSA 签名
func (r *RSA) RSASign(plaintext []byte) []byte {
 // 只有小消息可以直接签名； 因此，对消息的哈希进行签名，而不能对消息本身进行签名。
 // 这要求哈希函数必须具有抗冲突性。 SHA-256是编写本文时(2016年)应使用的最低强度的哈希函数。
 hashed := sha256.Sum256(plaintext)
 // 使用私钥签名，必须要将明文hash后才可以签名，当验证时，同样需要对明文进行hash运算。签名于验签并不用于加密消息或消息传递，仅仅作为验证传递消息方的真实性。
 signature, err := rsa.SignPKCS1v15(rand.Reader, r.rsaPrivateKey, crypto.SHA256, hashed[:])
 if err != nil {
  fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
  return nil
 }
 fmt.Printf("Signature: %x\n", signature)
 return signature
}

// RSAVerify RSA 验签
func (r *RSA) RSAVerify(plaintext []byte, signature []byte) bool {
 // 与签名一样，只可以对 hash 后的消息进行验证。
 hashed := sha256.Sum256(plaintext)
 // 使用公钥、已签名的信息，验证签名的真实性
 err := rsa.VerifyPKCS1v15(r.rsaPublicKey, crypto.SHA256, hashed[:], signature)
 if err != nil {
  fmt.Fprintf(os.Stderr, "Error from verification: %s\n", err)
  return false
 }
 return true
}
```

# RSA 示例

main.go

```go
package main

import "fmt"

func main() {
 // 生成rsa的密钥对, 并且保存到磁盘文件中
 r := NewRSA(4096)

 // 该消息有两个作用：
 // 1. 使用公钥加密的的信息
 // 2. 验证签名时所用的消息。当该消息用于签名时，通常还需要将该消息，以及用私钥签名后的消息一起发送给对方。以便对方可以根据该消息验证签名的有效性。
 messages := []byte("你好 DesistDaydream！...这是一串待加密的字符串，如果你能看到，那么说明功能实现了！")

 // 使用公钥加密，私钥解密
 encryptedMessages := r.RSAEncrypt(messages)
 decryptedMessages := r.RSADecrypt(encryptedMessages)
 fmt.Printf("解密后的字符串为：%v\n", string(decryptedMessages))

 // 使用私钥签名，公钥验签
 // 注意，验证签名需要使用签名时发送的消息作为对比，只有消息一致，才算验证通过
 signature := r.RSASign(messages)
 if r.RSAVerify(messages, signature) {
  fmt.Println("验证成功")
 }
}
```

rsa_key_handler.go

```go
package main

import (
 "crypto"
 "crypto/rand"
 "crypto/sha256"
 "fmt"
 "os"

 "crypto/rsa"
)

// RSA 是公钥和私钥两个组成一组的密钥对
type RSA struct {
 rsaPrivateKey *rsa.PrivateKey
 rsaPublicKey  *rsa.PublicKey
}

// NewRSA 生成密钥对
func NewRSA(bits int) *RSA {
 // 随机生成一个给定大小的 RSA 密钥对。可以使用 crypto 包中的 rand.Reader 来随机。
 privateKey, _ := rsa.GenerateKey(rand.Reader, bits)
 // 从私钥中，获取公钥
 publicKey := privateKey.PublicKey
 return &RSA{
  rsaPrivateKey: privateKey,
  rsaPublicKey:  &publicKey,
 }
}

// RSAEncrypt 使用 RSA 算法，加密指定明文
func (r *RSA) RSAEncrypt(plaintext []byte) []byte {
 // 使用公钥加密 plaintext(明文，也就是准备加密的消息)。并返回 ciphertext(密文)
 // 其中 []byte("DesistDaydream") 是加密中的标签，解密时标签需与加密时的标签相同，否则解密失败
 ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, r.rsaPublicKey, plaintext, []byte("DesistDaydream"))
 if err != nil {
  panic(err)
 }
 return ciphertext
}

// RSADecrypt 使用 RSA 算法，解密指定密文
func (r *RSA) RSADecrypt(ciphertext []byte) []byte {
 // 使用私钥解密 ciphertext(密文，也就是加过密的消息)。并返回 plaintext(明文)
 // 其中 []byte("DesistDaydream") 是加密中的标签，解密时标签需与加密时的标签相同，否则解密失败
 plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, r.rsaPrivateKey, ciphertext, []byte("DesistDaydream"))
 if err != nil {
  panic(err)
 }
 return plaintext
}

// RSASign RSA 签名
func (r *RSA) RSASign(plaintext []byte) []byte {
 // 只有小消息可以直接签名； 因此，对消息的哈希进行签名，而不能对消息本身进行签名。
 // 这要求哈希函数必须具有抗冲突性。 SHA-256是编写本文时(2016年)应使用的最低强度的哈希函数。
 hashed := sha256.Sum256(plaintext)
 // 使用私钥签名，必须要将明文hash后才可以签名，当验证时，同样需要对明文进行hash运算。签名于验签并不用于加密消息或消息传递，仅仅作为验证传递消息方的真实性。
 signature, err := rsa.SignPKCS1v15(rand.Reader, r.rsaPrivateKey, crypto.SHA256, hashed[:])
 if err != nil {
  fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
  return nil
 }
 fmt.Printf("Signature: %x\n", signature)
 return signature
}

// RSAVerify RSA 验签
func (r *RSA) RSAVerify(plaintext []byte, signature []byte) bool {
 // 与签名一样，只可以对 hash 后的消息进行验证。
 hashed := sha256.Sum256(plaintext)
 // 使用公钥、已签名的信息，验证签名的真实性
 err := rsa.VerifyPKCS1v15(r.rsaPublicKey, crypto.SHA256, hashed[:], signature)
 if err != nil {
  fmt.Fprintf(os.Stderr, "Error from verification: %s\n", err)
  return false
 }
 return true
}
```

# 最佳实践