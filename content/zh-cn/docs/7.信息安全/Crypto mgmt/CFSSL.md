---
title: CFSSL
linkTitle: CFSSL
date: 2023-09-12T08:12:00
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目,cloudflare/cfssll](https://github.com/cloudflare/cfssl)
> - [官网](https://cfssl.org/)
> - [公众号](https://mp.weixin.qq.com/s/E-aU-lbieGLokDKbjdGc3g)

cfssl 与 openssl 类似，不过是使用 go 编写，由 CloudFlare 开源的一款 PKI/TLS 工具。主要程序有 cfssl，是 CFSSL 的命令行工具，cfssljson 用来从 cfssl 程序获取 JSON 输出，并将证书，密钥，CSR 和 bundle 写入文件中。

## 使用 CFSSL 创建 CA 认证步骤

### 创建认证中心(CA)

cfssl 可以创建一个获取和操作证书的内部认证中心。运行认证中心需要一个 CA 证书和相应的 CA 私钥。任何知道私钥的人都可以充当 CA 来颁发证书。因此，私钥的保护至关重要，这里我们以 k8s 所需的证书来实践一下：

```bash
cfssl print-defaults config > config.json # 默认证书策略配置模板
cfssl print-defaults csr > csr.json #默认csr请求模板
```

结合自身的要求，修改证书请求文件`csr.json`,证书 10 年

```json
{
  "CN": "kubernetes",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "BeiJing",
      "L": "BeiJing",
      "O": "k8s",
      "OU": "System"
    }
   ],
   "ca": {
    "expiry": "87600h"
  }
}
```

知识点:

- `"CN"`：Common Name，kube-apiserver 从证书中提取该字段作为请求的用户名 (User Name)
- `"O"`：Organization，kube-apiserver 从证书中提取该字段作为请求用户所属的组 (Group)
- `C`: Country， 国家
- `L`: Locality，地区，城市
- `O`: Organization Name，组织名称，公司名称
- `OU`: Organization Unit Name，组织单位名称，公司部门
- `ST`: State，州，省

证书配置模板文件`ca-config.json`

```json
{
  "signing": {
      "default": {
        "expiry": "87600h"
   },
  "profiles": {
    "kubernetes": {
      "usages": [
        "signing",
        "key encipherment",
        "server auth",
        "client auth"
      ],
      "expiry": "87600h"
    }
   }
  }
}
```

知识点：

- `config.json`：可以定义多个 `profiles`，分别指定不同的过期时间、使用场景等参数；后续在签名证书时使用某个 `profile`；此实例只有一个 kubernetes 模板。
- `signing`：表示该证书可用于签名其它证书；生成的 ca.pem 证书中`CA=TRUE`
- `server auth`：表示 client 可以用该 CA 对 server 提供的证书进行验证；
- `client auth`：表示 server 可以用该 CA 对 client 提供的证书进行验证；
- 注意标点符号，最后一个字段一般是没有`逗号`的。

初始化创建 CA 认证中心，将会生成`ca-key.pem`（私钥）和`ca.pem`（公钥）

```bash
cfssl gencert -initca ca-csr.json | cfssljson -bare ca
```

### 创建 kubernetes 证书

创建`kubernetes-csr.json`证书请求文件

```json
{
    "CN": "kubernetes",
    "hosts": [
        "127.0.0.1",
        "10.1.20.129",
        "10.1.20.128",
        "10.1.20.126",
        "10.1.20.127",
        "10.254.0.1",
        "*.kubernetes.master",
        "localhost",
        "kubernetes",
        "kubernetes.default",
        "kubernetes.default.svc",
        "kubernetes.default.svc.cluster",
        "kubernetes.default.svc.cluster.local"
    ],
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "CN",
            "ST": "BeiJing",
            "L": "BeiJing",
            "O": "k8s",
            "OU": "System"
        }
    ]
}
```

**知识点**：

- 这个证书目前专属于 apiserver，加了一个 `*.kubernetes.master`域名以便内部私有 DNS 解析使用(可删除)；至于很多人问过 kubernetes 这几个能不能删掉，答案是**不可以**的；因为当集群创建好后，`default namespace` 下会创建一个叫 `kubenretes` 的`svc`，有一些组件会直接连接这个 svc 来跟 api 通信的，证书如果不包含可能会出现无法连接的情况；其他几个 kubernetes 开头的域名作用相同
- `hosts`包含的是授权范围，不在此范围的的节点或者服务使用此证书就会报证书不匹配错误。`10.254.0.1`是指 kube-apiserver 指定的 service-cluster-ip-range 网段的第一个 IP。

生成 kubernetes 证书和私钥

```bash
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=kubernetes kubernetes-csr.json | cfssljson -bare kubernetes
```

**知识点**：

- `-config` 引用的是模板中的默认配置文件，
- `-profiles`是指定特定的使用场景，比如 config.json 中的`kubernetes`区域

### 创建 admin 证书

创建 admin 证书请求文件`admin-csr.json`

```json
{
    "CN": "admin",
    "hosts": [],
    "key": {
    "algo": "rsa",
    "size": 2048
    },
    "names": [
    {
        "C": "CN",
        "ST": "BeiJing",
        "L": "BeiJing",
        "O": "system:masters",
        "OU": "System"
    }
    ]
}
```

生成 admin 证书和私钥

```bash
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=kubernetes admin-csr.json | cfssljson -bare admin
```

**知识点**这个 admin 证书，是将来生成管理员用的`kubeconfig` 配置文件用的，现在我们一般建议使用 RBAC 来对 kubernetes 进行角色权限控制， kubernetes 将证书中的 CN 字段作为 User， O 字段作为 Group
同样，我们也可以按照同样的方式来创建 kubernetes 中 etcd 集群的证书

### 创建 etcd 集群证书

1. 证书签署请求文件`ca-csr.json`


```json
{
    "CN": "etcd CA",
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "CN",
            "L": "Beijing",
            "ST": "Beijing"
        }
    ]
}
```

2. 为节点创建服务证书请求文件，指定授权的主机节点`etcd-server-csr.json`


```json
{
    "CN": "etcd",
    "hosts": [
        "10.1.20.129",
        "10.1.20.126",
        "10.1.20.128"
        ],
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "CN",
            "L": "BeiJing",
            "ST": "BeiJing"
        }
    ]
}
```

3. 证书配置模板文件`ca-config.json`


```json
{
  "signing": {
    "default": {
      "expiry": "87600h"
    },
    "profiles": {
      "etcd": {
         "expiry": "87600h",
         "usages": [
            "signing",
            "key encipherment",
            "server auth",
            "client auth"
        ]
      }
    }
  }
}
```

5. 生成 etcd 集群所需的证书与私钥


```bash
cfssl gencert -initca ca-csr.json | cfssljson -bare ca -
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=etcd etcd-server-csr.json | cfssljson -bare server
```

这样就完成 etcd 所需证书的申请，同时了解了 cfssl 工具的强大，写到这里，本次的实验就结束了。
