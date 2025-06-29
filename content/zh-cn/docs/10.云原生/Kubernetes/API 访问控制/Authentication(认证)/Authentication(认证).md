---
title: Authentication(认证)
linkTitle: Authentication(认证)
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，参考 - API 访问控制 - 认证](https://kubernetes.io/docs/reference/access-authn-authz/authentication/)

**Authentication(名词)/Authenticating(动词)(身份验证)**，指明客户端是否有权限访问 API Server。

就好比我们在登录一个网站时，需要输入账户和密码的概念类似。在使用 API Server 时，也是通过类似的方式，使用账户来登录 API server(虽然不是真的登录)。

## Accounts - Kubernetes 集群中的账号

Accounts 是一个在认证授权系统里的逻辑概念。Accounts 需要通过认证概念中的东西(比如证书、token、或者用户名和密码等)来建立。类似于登陆系统的账户。而在 Kubernetes 中，Accounts 分为如下两类

1. **UserAccount(用户账户，简称 User)**
2. **ServiceAccount(服务账户，简称 SA)**

> [!Tip]
> UA 与 SA 的对比在 [官方文档，参考 - API 访问控制 - 管理服务账号, User accounts 与 Service accounts](https://kubernetes.io/docs/reference/access-authn-authz/service-accounts-admin/#user-accounts-versus-service-accounts) 有提到，官方并没有对 UserAccount 进行明确的定义，偏向于一个没有实体的抽象概念，更多的时候是用 **KubeConfig** 这个词来作为 UserAccount 功能的实现。
>
> - UA 用来给人。SA 用来给运行在 pod 中的进程
> - UA 作用于全局，UA 的名字在集群的所有 namespace 中必须是唯一的。SA 作用于 namespace
> - UA 于 SA 的账户审核注意事项是不同的，UA 的凭证信息需要在使用 kubectl config 命令时候的手动指定；SA 的凭证信息在创建 SA 后会自动生成对应的 secret 并把凭证信息保存其中。

### User Account(用户账号)

详见：[User Account(KubeConfig)](/docs/10.云原生/Kubernetes/API%20访问控制/Authentication(认证)/User%20Account(KubeConfig).md)

User 不属于 K8S 中的一个资源。这类 Account 适用于：客户端访问集群时使用(比如使用 kubectl、scheduler 等访问 api)

一个 User 可以管理多个 k8s 集群、也可以多个 User 管理一个集群，权限不同。User 只有在 KubeConfig 文件中才具有实际意义。

由于 User 不属于 K8S 资源，那么则无法通过 API 调用来添加 User Account。但是任何提供了由群集的证书颁发机构(CA)签名的有效证书的用户都将被视为已认证。基于这种情况，Kubernetes 使用证书中的 subject 字段中的 Common Name(通用名称,即 CN)的值，作为用户名。接下来，基于授权概念中的 RBAC 子系统会确定用户是否有权针对某资源执行特定的操作。

如果想创建一个 User，则可以通过证书的方式来创建。比如像下面这样， 这就创建了一个名为 lch 的 User Account。

```bash
openssl genrsa -out lch.key 2048
openssl req -new -key lch.key -out lch.csr -subj "/CN=lch"
```

如果想使用 lch 这个 UA，则需要使用 kubectl config set-credentials 命令指定 lch 所需的相关凭证即可。还需要为 lch 绑定[授权概念](</docs/10.云原生/2.3.Kubernetes%20 容器编排系统/7.API%20 访问控制/2.Authorization(授权).md>>)中的 Role 以便让该用户具有某些操作权限，然后 lch 这个 UA 即可对所绑定的集群有 Role 中所指定的操作权限。其中为 -subj 选项中 CN 的值就是 User 的名称。这个值也是在后面为 User 赋予 RBAC 权限的 rolebinding 时所使用的 `subjects.name` 字段的值。

进一步的细节可参阅 [证书请求](https://kubernetes.io/docs/reference/access-authn-authz/certificate-signing-requests/#normal-user) 下普通用户主题。

### Service Account(服务账号)

详见：[Service Account](/docs/10.云原生/Kubernetes/API%20访问控制/Authentication(认证)/Service%20Account.md)

**Service Account(服务账号，简称 SA)** 属于 K8S 中的一个资源。这类 Account 适用于：Pod 访问集群时使用。

为什么需要 Service Account 呢？

SA 概念的引入是基于这样的使用场景：运行在 pod 里的进程需要调用 Kubernetes API 以及非 Kubernetes API 的其它服务。Service Account 是给 pod 里面 Container 中的进程使用的，它为 pod 提供必要的身份认证。(与用户控制 kubectl 去调用 API 一样，这里相当于 Pod 中 Container 在调用 API 的时候需要的认证)

## Accounts Group(账户组) - UserAccount 与 ServiceAccount 都有 Group

UA 与 SA 都可以属于一个或多个 Group

Group 是 Account 的集合，本身并没有操作权限，但附加于 Group 上的权限可由其内部的所有用户继承，以实现高效的授权管理机制。Kubernetes 有几个内建的用于特殊目的的 Group：

1. system:unauthenticated
2. system:authenticated
3. system:serviceaccounts
4. system:serviceaccounts:\<NameSpace>

KubeConfig 会给 UserAccount 提供与 APIServer 交互时所用的证书

Secret 会给 ServiceAccount 提供与 APIServer 交互时所用的证书

# Authentication Strategies 认证策略(i.e.Account 可用的认证方式)

https://kubernetes.io/zh/docs/reference/access-authn-authz/authentication/#authentication-strategies

Kubernetes 接受的认证方式有如下几种：

- client certificates
- bearer tokens
- an authenticating proxy
- HTTP basic auth

向 API Server 发起 HTTPS 请求时，kubernetes 通过身份验证插件对请求进行身份验证。

## X509 Client Certs(X509 客户端证书)

https://kubernetes.io/docs/reference/access-authn-authz/authentication/#x509-client-certs

通过给 API 服务器传递 --client-ca-file=SOMEFILE 选项，就可以启动客户端证书身份认证。 所引用的文件必须包含一个或者多个证书机构，用来验证向 API 服务器提供的客户端证书。 如果提供了客户端证书并且证书被验证通过，则 subject 中的公共名称（Common Name）就被 作为请求的用户名。 自 Kubernetes 1.4 开始，客户端证书还可以通过证书的 organization 字段标明用户的组成员信息。 要包含用户的多个组成员信息，可以在证书种包含多个 organization 字段。

例如，使用 openssl 命令行工具生成一个证书签名请求：

```bash
openssl req -new -key jbeda.pem -out jbeda-csr.pem -subj "/CN=jbeda/O=app1/O=app2"
```

此命令将使用用户名 jbeda 生成一个证书签名请求（CSR），且该用户属于 "app" 和 "app2" 两个用户组。

参阅[管理证书](https://kubernetes.io/docs/concepts/cluster-administration/certificates/)了解如何生成客户端证书

## Static Token File(静态令牌文件)

https://kubernetes.io/docs/reference/access-authn-authz/authentication/#static-token-file

当 API 服务器的命令行设置了 --token-auth-file=SOMEFILE 选项时，会从文件中 读取持有者令牌。目前，令牌会长期有效，并且在不重启 API 服务器的情况下 无法更改令牌列表。

令牌文件是一个 CSV 文件，包含至少 3 个列：令牌、用户名和用户的 UID。 其余列被视为可选的组名。

说明：

如果要设置的组名不止一个，则对应的列必须用双引号括起来，例如

```bash
token,user,uid,"group1,group2,group3"
```

在请求中放入持有者令牌

当使用持有者令牌来对某 HTTP 客户端执行身份认证时，API 服务器希望看到 一个名为 Authorization 的 HTTP 头，其值格式为 Bearer THETOKEN。 持有者令牌必须是一个可以放入 HTTP 头部值字段的字符序列，至多可使用 HTTP 的编码和引用机制。 例如：如果持有者令牌为 31ada4fd-adec-460c-809a-9e56ceb75269，则其 出现在 HTTP 头部时如下所示：

```bash
Authorization: Bearer 31ada4fd-adec-460c-809a-9e56ceb75269
# 比如一个 curl 请求中，可以通过 -H 参数加入请求头
curl --cacert ${CAPATH} -H "Authorization: Bearer ${TOKEN}"  https://${IP}:6443/
```
