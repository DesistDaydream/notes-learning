---
title: HTTPS 和 Authentication(认证)
---

# 概述

> 参考：
> - [官方文档,Prometheus-配置-HTTPS 和 认证](https://prometheus.io/docs/prometheus/latest/configuration/https/)
> - [官方文档,指南-基础认证](https://prometheus.io/docs/guides/basic-auth/)
> - [Ngxin Ingress Controller 官方文档,认证-基础认证](https://kubernetes.github.io/ingress-nginx/examples/auth/basic/)
> - [知乎](https://www.jianshu.com/p/4d5aa1995de3)

认证功能的发展：

- Prometheus 从 2.24 版本开始，才支持基本认证，截止 2021 年 8 月 25 日官方已经提供了实验性的 HTTPS 与 认证配置，详见[此处](https://prometheus.io/docs/prometheus/latest/configuration/https/)。
- 截止 2021 年 8 月 25 日，Thanos 的 Sidecar 还不支持向 Prometheus 发起请求是携带认证信息，但已经有 [issue #3975](https://github.com/thanos-io/thanos/issues/3975) 提出来该问题，并将在未来 [PR #4104](https://github.com/thanos-io/thanos/pull/4104) 实现

现阶段在 Prometheus 前面添加代理(比如 Nginx)，只暴露 Nginx 端口，所有访问 Prometheus 的请求都经过代理，并在代理上添加认证，这样可以为 Prometheus 的 web 端添加一个基本的基于用户名和密码的认证。

在 kubernetes 中，可以通过 ingress 来实现。其他环境可以直接配置 ngxin 来实现。

# 通过 ingress controller 配置认证，普通的 nginx 同理。

首先需要安装 htpasswd 二进制文件，通过 htpasswd 命令行工具生成保存用户名密码的文件，然后通过该文件创建一个 secret 对象，并在 ingress 引用该 secret 对象

**通过 htpasswd 生成一个“auth”文件;用来存取我们创建的用户及加密之后的密码**

```bash
root@lichenhao:~# htpasswd -c auth admin
New password:
Re-type new password:
Adding password for user admin
# 查看这个文件，可以看到密码是加密之后的字符串
root@lichenhao:~# cat auth
admin:$apr1$8NSwCSR3$s5G25cvkaUDAoxEFtaGZ11
# 密码：ehl1234
```

**创建 kubernetes secret 来存储 auth 文件中的用户名和密码**

```yaml
root@lichenhao:~# kubectl create -n monitoring secret generic basic-auth --from-file=auth
secret "basic-auth" created

root@lichenhao:~# kubectl get secrets -n monitoring basic-auth -oyaml | neat
apiVersion: v1
data:
  auth: YWRtaW46JGFwcjEkOE5Td0NTUjMkczVHMjVjdmthVURBb3hFRnRhR1oxMQo=
kind: Secret
metadata:
  annotations:
    meta.helm.sh/release-name: monitor-bj-net
    meta.helm.sh/release-namespace: monitoring
  labels:
    app.kubernetes.io/managed-by: Helm
  name: basic-auth
  namespace: monitoring
type: Opaque
```

**在 ingress 资源中添加注释**

```yaml
root@lichenhao:~# kubectl get ingress -n monitoring  monitor-bj-net-k8s-prometheus -oyaml | neat
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    meta.helm.sh/release-name: monitor-bj-net
    meta.helm.sh/release-namespace: monitoring
    # 添加如下三行注释内容，
    # 指定认证类型
    nginx.ingress.kubernetes.io/auth-type: basic
    # 指定包含用户名与密码的 secret 资源的名称
    nginx.ingress.kubernetes.io/auth-secret: basic-auth
    # 消息，显示一个适当的上下文，说明为什么需要身份验证。最后 admin 就是指明应该要使用的用户名
    nginx.ingress.kubernetes.io/auth-realm: Authentication Required - admin
  labels:
    app: k8s-prometheus
    app.kubernetes.io/managed-by: Helm
    chart: kube-prometheus-stack-12.12.1
    heritage: Helm
    release: monitor-bj-net
  name: monitor-bj-net-k8s-prometheus
  namespace: monitoring
spec:
  ingressClassName: nginx
  rules:
  - host: prometheus.desistdaydream.ltd
    http:
      paths:
      - backend:
          serviceName: monitor-bj-net-k8s-prometheus
          servicePort: 9090
        path: /
        pathType: ImplementationSpecific
```

# 为 Prometheus 添加基础认证

Prometheus 通过 `--web.config.file` 命令行标志来开启 **TLS** 或者 **基本认证** 功能。Prometheus 将会读取该标志指定的文件，通过该文件的配置来为 9090 端口添加 TLS 或者 基本认证。让人们在访问 9090 端口的 Web UI 或者使用 API 时，必须进行认证才可以操作。

## web-config 文件详解

该文件有 3 个主要字段

- tls_server_config
- http_server_config
- basic_auth_users

如果没有任何配置，则不开启任何 TLS 或 认证，只要配置了某个字段，就默认开启相关功能。

### tls_server_config: <Object> # 为 Prometheus 开启 TLS

- **cert_file: <FileName>** # TLS 所需的证书文件
- **key_file: <FileName>** # TLS 所需的私钥文件

  Server policy for client authentication. Maps to ClientAuth Policies.
  &#x20; \# For more detail on clientAuth options:
  &#x20; \# <https://golang.org/pkg/crypto/tls/#ClientAuthType>
  &#x20; **client_auth_type: <STRING>** # `默认值：NoClientCert`
  ============================================================

  CA certificate for client certificate authentication to the server.
  &#x20; \[ client_ca_file: <filename> ]
  ========================================

  Minimum TLS version that is acceptable.
  &#x20; \[ min_version: <string> | default = "TLS12" ]
  ======================================================

  Maximum TLS version that is acceptable.
  &#x20; \[ max_version: <string> | default = "TLS13" ]
  ======================================================

  List of supported cipher suites for TLS versions up to TLS 1.2. If empty,
  &#x20; \# Go default cipher suites are used. Available cipher suites are documented
  &#x20; \# in the go documentation:
  &#x20; \# <https://golang.org/pkg/crypto/tls/#pkg-constants>
  &#x20; \[ cipher_suites:
  &#x20; \[ - <string> ] ]
  ==========================

  prefer_server_cipher_suites controls whether the server selects the
  &#x20; \# client's most preferred ciphersuite, or the server's most preferred
  &#x20; \# ciphersuite. If true then the server's preference, as expressed in
  &#x20; \# the order of elements in cipher_suites, is used.
  &#x20; \[ prefer_server_cipher_suites: <bool> | default = true ]
  ===================================================================

  Elliptic curves that will be used in an ECDHE handshake, in preference
  &#x20; \# order. Available curves are documented in the go documentation:
  &#x20; \# <https://golang.org/pkg/crypto/tls/#CurveID>
  &#x20; \[ curve_preferences:
  &#x20; \[ - <string> ] ]
  ==========================

### http_server_config: <Object> # 为 Prometheus 开启 HTTP/2。注意，HTTP/2 仅支持 TLS

- **http2: <BOOLEAN>** # `默认值：true`

Usernames and hashed passwords that have full access to the web
\# server via basic authentication. If empty, no basic authentication is
\# required. Passwords are hashed with bcrypt.
==============================================

### basic_auth_users: \<map\[STRING]STRING> # 为 Prometheus Server 开启基本认证

- **<KEY>: <VALUE>** # KEY 是用户名，VALUE 是密码
  - 注意：密码必须是经过 hash 的字符串，可以通过[这个网站](https://bcrypt-generator.com/)在线获取 hash 过的字符串

## 配置示例

在[这里](https://bcrypt-generator.com/)生成密码的 hash 值，比如我使用 `Prometheus@lichenhao` 这个密码，生成的 hash 为 `$2a$12$twJp6N9kL5aEf08Ja8XRAOImHOjCTBQvb485Uuz7hJLEX1XT4iVDm`

在 /etc/prometheus/config_out/ 目录中创建一个 web-config.yml 文件

```bash
basic_auth_users:
  Prometheus: $2a$12$twJp6N9kL5aEf08Ja8XRAOImHOjCTBQvb485Uuz7hJLEX1XT4iVDm
```

为 Promethues Server 添加命令行标志 `--web.config.file=/etc/prometheus/config_out/web-config.yml`

启动 Prometheus 后，将会需要认证信息，效果如下：
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/bx144g/1628063307526-21ac3e4b-150d-4e77-9a7c-5069ad006369.png)
