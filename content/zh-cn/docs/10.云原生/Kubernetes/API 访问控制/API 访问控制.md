---
title: API 访问控制
linkTitle: API 访问控制
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，概念 - 安全 - Kubernetes API 的访问控制](https://kubernetes.io/docs/concepts/security/controlling-access/)

**认证用于身份验证，授权用于权限检查，准入控制机制用于补充授权机制**

客户端与服务端的概念：谁向谁发请求，前者就是客户端，所在在这里，客户端与服务端没有绝对。一个服务既可以是客户端也可以是服务端，kubectl 在控制集群需要给 apiservice 发送 get，creat，delete 等指令的时候，kubectl 就是 apiservice 的客户端；而 apiservice 需要往 etcd 写入数据的时候，apiservice 就是 etcd 的客户端。

当客户端向服务端发起请求的时候，服务端需要对客户端进行认证以便确认客户端身份是否可以接入；接入后再进行授权检查，检查该身份的请求是否可以在服务端执行。所以后面介绍的 认证 与 授权 是相辅相成，不可分隔，创建完认证之后，需要为这个认证信息进行授权，才是一套完整的鉴权机制

> 比如现在有这么一个场景，张三要去商场买酱油。当张三到达商场后，保安人员首先要对张三进行认证，确认张三这个人可以进入商场；然后张三到达货柜拿走酱油去结账，收银人员进行授权检查，核验张三是否有权力购买酱油。

在 kubernetes 集群中，就是类似张三买酱油的场景。。。各个组件与资源对象之间的互相访问，在大多数时候，都需要进行认证与授权的检查。

API Server 是集群的入口，不管是对资源对象的增删改查，还是访问集群中的某些对象，不可避免得只能与 API Server 交互，虽然在访问某些管理组件的 https 端口时，也需要进行认证，但是这种访问是属于基本的 https 访问。所以，在与其说是 k8s 的认证与授权，不如说是 kubernetes API 的访问控制。因为不管是从外部(kubeclt 等)、还是内部(controller-manager、某个 pod 访问集群资源)，都逃不开与 kubernetes API，也就是 api-server 这个组件的交互。毕竟 kubernetes API 是集群的唯一入口。。。就算是在集群内部署的 pod，如果想要访问集群内的资源，也逃不开 kubernetes API~~

当然，使用 curl 命令来访问 controller、scheduler 时、或者 etcd 互相交互，都属于 认证与授权 的概念范畴~只不过这种情况不占大多数，所以就不再单独讨论了。这些认证授权方式与 API 的认证授权类似。

# Kubernetes API 访问控制

我们可以通过 kubectl、客户端库、发送 REST 请求 这几种方法访问 [Kubernetes API](https://kubernetes.io/docs/concepts/overview/kubernetes-api/)。[人类用户(User Account(KubeConfig))](/docs/10.云原生/Kubernetes/API%20访问控制/Authentication(认证)/User%20Account(KubeConfig).md) 和 [Kubernetes 的 Service Account](/docs/10.云原生/Kubernetes/API%20访问控制/Authentication(认证)/Service%20Account.md) 都可以被授权进行 API 访问。 请求到达 API Server 后会经过几个阶段，具体如下图所示：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/cvkvyz/1616118854890-e2e31942-d6ea-40a7-83d8-816abb4c136a.jpeg)

## 传输层安全

在典型的 Kubernetes 集群中，API 通过 443 端口提供服务。 API 服务器会提供一份证书。 该证书一般是私有 CA 自签名的，当然，也可以基于公信的 CA 公钥基础设施签名。

如果集群使用私有证书颁发机构，需要在客户端的  `~/.kube/config`  文件中提供该 CA 证书的副本， 以便在客户端使用程序访问 API 时，可以信任该连接并确认该连接没有被拦截。

## 认证

一旦 TLS 连接建立，HTTP 请求就进入到了认证的步骤。即图中的步骤 1 。 集群创建脚本或集群管理员会为 API 服务器配置一个或多个认证模块。 更具体的认证相关的描述详见[这里](https://kubernetes.io/docs/admin/authentication/)。

认证步骤的输入是整个 HTTP 请求，但这里通常只是检查请求头和 / 或客户端证书。

认证模块支持客户端证书，密码和 Plain Tokens， Bootstrap Tokens，以及 JWT Tokens（用于服务账户）。

（管理员）可以同时设置多种认证模块，在设置了多个认证模块的情况下，每个模块会依次尝试认证， 直到其中一个认证成功。

在 GCE 平台中，客户端证书，密码和 Plain Tokens，Bootstrap Tokens，以及 JWT Tokens 同时被启用。

如果请求认证失败，则请求被拒绝，返回 401 状态码。 如果认证成功，则被认证为具体的 username，该用户名可供随后的步骤中使用。一些认证模块还提供了用户的组成员关系，另一些则没有。

尽管 Kubernetes 使用“用户名”来进行访问控制和请求记录，但它实际上并没有 user 对象，也不存储用户名称或其他相关信息。

## 授权

当请求被认证为来自某个特定的用户后，该请求需要被授权。 即图中的步骤 2 。

请求须包含请求者的用户名，请求动作，以及该动作影响的对象。 如果存在相应策略，声明该用户具有进行相应操作的权限，则该请求会被授权。

例如，如果 Bob 有如下策略，那么他只能够读取 projectCaribou 命名空间下的 pod 资源：

```json
{
  "apiVersion": "abac.authorization.kubernetes.io/v1beta1",
  "kind": "Policy",
  "spec": {
    "user": "bob",
    "namespace": "projectCaribou",
    "resource": "pods",
    "readonly": true
  }
}
```

如果 Bob 发起以下请求，那么请求能够通过授权，因为 Bob 被允许访问 projectCaribou 命名空间下的对象：

```json
{
  "apiVersion": "authorization.k8s.io/v1beta1",
  "kind": "SubjectAccessReview",
  "spec": {
    "resourceAttributes": {
      "namespace": "projectCaribou",
      "verb": "get",
      "group": "unicorn.example.org",
      "resource": "pods"
    }
  }
}
```

如果 Bob 对 projectCaribou 命名空间下的对象发起一个写（create 或者 update）请求，那么它的授权会被拒绝。 如果 Bob 请求读取 （get）其他命名空间，例如 projectFish 下的对象，其授权也会被拒绝。

Kubernetes 的授权要求使用通用的 REST 属性与现有的组织或云服务提供商的访问控制系统进行交互。 采用 REST 格式是必要的，因为除 Kubernetes 外，这些访问控制系统还可能与其他的 API 进行交互。

Kubernetes 支持多种授权模块，例如 ABAC 模式，RBAC 模式和 Webhook 模式。 管理员创建集群时，会配置 API 服务器应用的授权模块。 如果多种授权模式同时被启用，Kubernetes 将检查所有模块，如果其中一种通过授权，则请求授权通过。 如果所有的模块全部拒绝，则请求被拒绝（HTTP 状态码 403）。

要了解更多的 Kubernetes 授权相关信息，包括使用授权模块创建策略的具体说明等，可参考[授权概述](https://kubernetes.io/docs/admin/authorization)。

## 准入控制

准入控制模块是能够修改或拒绝请求的软件模块。 作为授权模块的补充，准入控制模块会访问被创建或更新的对象的内容。 它们作用于对象的创建，删除，更新和连接（proxy）阶段，但不包括对象的读取。

可以同时配置多个准入控制器，它们会按顺序依次被调用。

即图中的步骤 3 。

与认证和授权模块不同的是，如果任一个准入控制器拒绝请求，那么整个请求会立即被拒绝。

除了拒绝请求外，准入控制器还可以为对象设置复杂的默认值。

可用的准入控制模块描述 [如下](https://kubernetes.io/docs/admin/admission-controllers/)。

一旦请求通过所有准入控制器，将使用对应 API 对象的验证流程对其进行验证，然后写入对象存储 （如步骤 4）。

# API 的端口和 IP

上述讨论适用于发送请求到 API 服务器的安全端口（典型情况）。
实际上 API 服务器可以通过两个端口提供服务，默认情况下，API 服务器在 2 个端口上提供 HTTP 服务：

- Localhost Port:
  - 用于测试和启动，以及管理节点的其他组件(scheduler, controller-manager)与 API 的交互
  - 没有 TLS
  - 默认值为 8080，可以通过 API Server 的 `--insecure-port` 命令行标志来修改。
  - 默认的 IP 地址为 localhost，可以通过 API Server 的 `--insecure-bind-address` 命令行标志来修改。
  - 请求会 **绕过** 认证和鉴权模块。
  - 请求会被准入控制模块处理。
  - 其访问需要主机访问的权限。
- Secure Port:
  - 尽可能使用该端口访问
  - 应用 TLS。 可以通过 API Server 的 `--tls-cert-file` 设置证书， `--tls-private-key-file` 设置私钥。
  - 默认值为 6443，可以通过 API Server 的 `--secure-port` 命令行标志来修改。
  - 默认 IP 是首个非本地的网络接口地址，可以通过 API Server 的 `--bind-address` 命令行标志来修改。
  - 请求会经过认证和鉴权模块处理。
  - 请求会被准入控制模块处理。
  - 要求认证和授权模块正常运行。
