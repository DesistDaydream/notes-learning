---
title: Token 与 JWT
---

# 概述

> 参考：
> 
> - [JWT 官网](https://jwt.io/)
> - [RFC 7519](https://datatracker.ietf.org/doc/html/rfc7519)
> - [阮一峰，JSON Web Token 入门教程](http://www.ruanyifeng.com/blog/2018/07/json_web_token-tutorial.html)

# JWT

**JSON Web Token(简称 JWT)** 是目前最流行的跨域认证解决方案。JWT 是一个开放标准([RFC 7519](https://tools.ietf.org/html/rfc7519))，定义了一种紧凑且子包含的方式，用于在各方之间以 JSON 对象的形式安全地传输信息。该信息可以被验证和信任，因为它是经过数字签名的。可以通过 RSA 或 ECDSA 这种公钥/私钥对进行签名、也可以使用 HMAC 进行签名。

## 背景

### 跨域认证的问题

互联网服务离不开用户认证。一般流程是下面这样。

1. 用户向服务器发送用户名和密码。
2. 服务器验证通过后，在当前对话（session）里面保存相关数据，比如用户角色、登录时间等等。
3. 服务器向用户返回一个 session_id，写入用户的 Cookie。
4. 用户随后的每一次请求，都会通过 Cookie，将 session_id 传回服务器。
5. 服务器收到 session_id，找到前期保存的数据，由此得知用户的身份。

这种模式的问题在于，扩展性（scaling）不好。单机当然没有问题，如果是服务器集群，或者是跨域的服务导向架构，就要求 session 数据共享，每台服务器都能够读取 session。

举例来说，A 网站和 B 网站是同一家公司的关联服务。现在要求，用户只要在其中一个网站登录，再访问另一个网站就会自动登录，请问怎么实现？

一种解决方案是 session 数据持久化，写入数据库或别的持久层。各种服务收到请求后，都向持久层请求数据。这种方案的优点是架构清晰，缺点是工程量比较大。另外，持久层万一挂了，就会单点失败。

另一种方案是服务器索性不保存 session 数据了，所有数据都保存在客户端，每次请求都发回服务器。JWT 就是这种方案的一个代表。

## JWT 的原理

JWT 的原理是，服务器认证以后，生成一个 JSON 对象，发回给用户，就像下面这样。

```json
{
  "姓名": "张三",
  "角色": "管理员",
  "到期时间": "2018年7月1日0点0分"
}
```

以后，用户与服务端通信的时候，都要发回这个 JSON 对象。服务器完全只靠这个对象认定用户身份。为了防止用户篡改数据，服务器在生成这个对象的时候，会加上签名（详见后文）。

服务器就不保存任何 session 数据了，也就是说，服务器变成无状态了，从而比较容易实现扩展。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/sn00y2/1623744463253-b71ac60f-9871-4aef-8c37-8166947d5b58.webp)

## JWT 的数据结构

实际的 JWT 大概就像下面这样。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/sn00y2/1623651349414-4518a10e-dbd4-4461-a4e2-a039518e8355.jpeg)

它是一个很长的字符串，中间用点（.）分隔成三个部分。<font color="#ff0000">注意，JWT 内部是没有换行的，这里只是为了便于展示，将它写成了几行。</font>

JWT 的三个部分依次如下。

- **Header(头部)** #
- **Payload(负载)** #
- **Signature(签名)** #

写成一行，就是下面的样子。

```bash
Header.Payload.Signature
```

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/sn00y2/1623651350412-81b31181-9912-4411-a330-00e6e92e791b.jpeg)

### Header

Header 部分是一个 JSON 对象，描述 JWT 的元数据，通常是下面的样子。

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

- alg # 表示签名的算法（algorithm），默认是 HMAC SHA256(写成 HS256)
- typ # 表示这个令牌（token）的类型（type），JWT 令牌统一写为 JWT。

最后，将上面的 JSON 对象使用 Base64URL 算法（详见后文）转成字符串。

### Payload

Payload 部分也是一个 JSON 对象，用来存放实际需要传递的数据。Payload 由多个 **Claims(要求)** 组成。Claims 通常用来描述实体和附加数据，实体通常来说就是一个用户。

一共有三种类型的 Claims：

- **Registered Claims(已注册要求)**，其实就是 JWT 默认自带的。JWT 规定了 7 个默认的 Claims。
  - **iss(issuer)** # 签发人。与 x509 证书的 Issuer 概念一样
  - **sub(subject)** # 主体。与 x509 证书的 Subject 概念一样
  - **aud(audience)** # 受众，就是 JWT 的人或程序或某事物
  - **exp(expiration time)** # 过期时间。这个时间必须要大于 `iat` 字段指定的时间
  - **nbf(Not Before)** # 生效时间。定义在什么时间之前该 JWT 是不可用的。类似 X.509 中的 `Not Before` 字段。
  - **iat(Issued At)** # 签发时间。
  - **jti(JWT ID)** # JWT 的 ID，主要用来作为一次性 Token，从而回避重放攻击
- **Public Claims(公共要求)**，使用 JWT 时的自定义 Claims。
  - 但是给 Public Claims 命名时，要注意不要与 [IANA JSON Web Token Registry](https://www.iana.org/assignments/jwt/jwt.xhtml) 中定义的名称相冲突。
- **Private Claims(私有要求)**，也是一种自定义的 Claims，通常用来共享信息。JWT 的生产者和消费者可以同意使用 Private Claims，而不是其他两种类型的 Claims。

下面就是一个有效的 Payload 的示例：

```json
{
  "sub": "1234567890",
  "name": "John Doe",
  "admin": true
}
```

注意，JWT 默认是不加密的，任何人都可以读到，所以不要把敏感信息放在 JWT 中。

这个 JSON 对象也要使用 Base64URL 算法转成字符串。

### Signature

**Signature(签名)** 是对前 Header 与 Payload 的签名，防止数据篡改。

首先，需要指定一个 secret(密钥)。这个密钥只有服务器才知道，不能泄露给用户。然后，使用 Header 里面指定的签名算法（默认是 HMAC SHA256），按照下面的公式产生签名。

```text
HMACSHA256(
  base64UrlEncode(HEADER) + "." +
  base64UrlEncode(PAYLOAD),
  secret)
```

当我们要解析 JWT 时，也要使用相同 Secret，如果 Secret，则验证签名将会失败，JWT 这个 Token 也就失效了~这就是签名的重要性啦~

### 总结

算出 Signature 以后，把 Header、Payload、Signature 三个部分拼成一个字符串，每个部分之间用"点"（.）分隔，这就组成了一个字符串类型的 JWT。

**Base64URL**

前面提到，Header 和 Payload 串型化的算法是 Base64URL。这个算法跟 Base64 算法基本类似，但有一些小的不同。

JWT 作为一个令牌（token），有些场合可能会放到 URL（比如 api.example.com/?token=xxx）。Base64 有三个字符+、/和=，在 URL 里面有特殊含义，所以要被替换掉：=被省略、+替换成-，/替换成\_ 。这就是 Base64URL 算法。

## JWT 的使用方式

客户端收到服务器返回的 JWT，可以储存在 Cookie 里面，也可以储存在 localStorage。

此后，客户端每次与服务器通信，都要带上这个 JWT。你可以把它放在 Cookie 里面自动发送，但是这样不能跨域，所以更好的做法是放在 HTTP 请求的头信息 Authorization 字段里面。

`Authorization: Bearer <TOKEN>`

另一种做法是，跨域的时候，JWT 就放在 POST 请求的数据体里面。

## JWT 的几个特点

（1）JWT 默认是不加密，但也是可以加密的。生成原始 Token 以后，可以用密钥再加密一次。
（2）JWT 不加密的情况下，不能将秘密数据写入 JWT。
（3）JWT 不仅可以用于认证，也可以用于交换信息。有效使用 JWT，可以降低服务器查询数据库的次数。
（4）JWT 的最大缺点是，由于服务器不保存 session 状态，因此无法在使用过程中废止某个 token，或者更改 token 的权限。也就是说，一旦 JWT 签发了，在到期之前就会始终有效，除非服务器部署额外的逻辑。
（5）JWT 本身包含了认证信息，一旦泄露，任何人都可以获得该令牌的所有权限。为了减少盗用，JWT 的有效期应该设置得比较短。对于一些比较重要的权限，使用时应该再次对用户进行认证。
（6）为了减少盗用，JWT 不应该使用 HTTP 协议明码传输，要使用 HTTPS 协议传输。

## 参考链接

- [Introduction to JSON Web Tokens](https://jwt.io/introduction/)， by Auth0
- [Sessionless Authentication using JWTs (with Node + Express + Passport JS)](https://medium.com/@bryanmanuele/sessionless-authentication-withe-jwts-with-node-express-passport-js-69b059e4b22c), by Bryan Manuele
- [Learn how to use JSON Web Tokens](https://github.com/dwyl/learn-json-web-tokens/blob/master/README.md), by dwyl
