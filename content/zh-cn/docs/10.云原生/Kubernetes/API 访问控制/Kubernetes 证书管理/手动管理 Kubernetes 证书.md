---
title: 手动管理 Kubernetes 证书
---

# 概述

> 参考：
>
> - [官方文档，任务-管理集群-证书](https://kubernetes.io/docs/tasks/administer-cluster/certificates/)

在使用客户端证书认证的场景下，你可以通过 `easyrsa`、`openssl` 或 `cfssl` 等工具以手工方式生成证书。

### easyrsa

**easyrsa** 支持以手工方式为你的集群生成证书。

1. 下载、解压、初始化打过补丁的 easyrsa3。

```bash
curl -LO https://storage.googleapis.com/kubernetes-release/easy-rsa/easy-rsa.tar.gz
tar xzf easy-rsa.tar.gz
cd easy-rsa-master/easyrsa3
./easyrsa init-pki
```

1. 生成新的证书颁发机构（CA）。参数 `--batch` 用于设置自动模式； 参数 `--req-cn` 用于设置新的根证书的通用名称（CN）。

    ./easyrsa --batch "--req-cn=${MASTER_IP}@`date +%s`" build-ca nopass

1. 生成服务器证书和秘钥。 参数 `--subject-alt-name` 设置 API 服务器的 IP 和 DNS 名称。 `MASTER_CLUSTER_IP` 用于 API 服务器和控制管理器，通常取 CIDR 的第一个 IP，由 `--service-cluster-ip-range` 的参数提供。 参数 `--days` 用于设置证书的过期时间。 下面的示例假定你的默认 DNS 域名为 `cluster.local`。

```bash
./easyrsa --subject-alt-name="IP:${MASTER_IP},"\
"IP:${MASTER_CLUSTER_IP},"\
"DNS:kubernetes,"\
"DNS:kubernetes.default,"\
"DNS:kubernetes.default.svc,"\
"DNS:kubernetes.default.svc.cluster,"\
"DNS:kubernetes.default.svc.cluster.local" \
--days=10000 \
build-server-full server nopass
```

1. 拷贝文件 `pki/ca.crt`、`pki/issued/server.crt` 和 `pki/private/server.key` 到你的目录中。

2. 在 API 服务器的启动参数中添加以下参数：

```bash
--client-ca-file=/yourdirectory/ca.crt
--tls-cert-file=/yourdirectory/server.crt
--tls-private-key-file=/yourdirectory/server.key
```

### openssl

**openssl** 支持以手工方式为你的集群生成证书。

1. 生成一个 2048 位的 ca.key 文件

    openssl genrsa -out ca.key 2048

1. 在 ca.key 文件的基础上，生成 ca.crt 文件（用参数 -days 设置证书有效期）

    openssl req -x509 -new -nodes -key ca.key -subj "/CN=${MASTER_IP}" -days 10000 -out ca.crt

1. 生成一个 2048 位的 server.key 文件：

    openssl genrsa -out server.key 2048

1. 创建一个用于生成证书签名请求（CSR）的配置文件。 保存文件（例如：`csr.conf`）前，记得用真实值替换掉尖括号中的值（例如：`<MASTER_IP>`）。 注意：`MASTER_CLUSTER_IP` 就像前一小节所述，它的值是 API 服务器的服务集群 IP。 下面的例子假定你的默认 DNS 域名为 `cluster.local`。

```ini
[ req ]
default_bits = 2048
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn
[ dn ]
C = <country>
ST = <state>
L = <city>
O = <organization>
OU = <organization unit>
CN = <MASTER_IP>
[ req_ext ]
subjectAltName = @alt_names
[ alt_names ]
DNS.1 = kubernetes
DNS.2 = kubernetes.default
DNS.3 = kubernetes.default.svc
DNS.4 = kubernetes.default.svc.cluster
DNS.5 = kubernetes.default.svc.cluster.local
IP.1 = <MASTER_IP>
IP.2 = <MASTER_CLUSTER_IP>
[ v3_ext ]
authorityKeyIdentifier=keyid,issuer:always
basicConstraints=CA:FALSE
keyUsage=keyEncipherment,dataEncipherment
extendedKeyUsage=serverAuth,clientAuth
subjectAltName=@alt_names
```

1. 基于上面的配置文件生成证书签名请求：

    openssl req -new -key server.key -out server.csr -config csr.conf

1. 基于 ca.key、ca.key 和 server.csr 等三个文件生成服务端证书：

```bash
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key \
-CAcreateserial -out server.crt -days 10000 \
-extensions v3_ext -extfile csr.conf
```

1. 查看证书：

    openssl x509  -noout -text -in ./server.crt

最后，为 API 服务器添加相同的启动参数。

### cfssl

**cfssl** 是另一个用于生成证书的工具。

1. 下载、解压并准备如下所示的命令行工具。 注意：你可能需要根据所用的硬件体系架构和 cfssl 版本调整示例命令。

```bash
curl -L https://github.com/cloudflare/cfssl/releases/download/v1.5.0/cfssl_1.5.0_linux_amd64 -o cfssl
chmod +x cfssl
curl -L https://github.com/cloudflare/cfssl/releases/download/v1.5.0/cfssljson_1.5.0_linux_amd64 -o cfssljson
chmod +x cfssljson
curl -L https://github.com/cloudflare/cfssl/releases/download/v1.5.0/cfssl-certinfo_1.5.0_linux_amd64 -o cfssl-certinfo
chmod +x cfssl-certinfo
```

1. 创建一个目录，用它保存所生成的构件和初始化 cfssl：

```bash
mkdir cert
cd cert
../cfssl print-defaults config > config.json
../cfssl print-defaults csr > csr.json
```

1. 创建一个 JSON 配置文件来生成 CA 文件，例如：`ca-config.json`：

```json
{
  "signing": {
    "default": {
      "expiry": "8760h"
    },
    "profiles": {
      "kubernetes": {
        "usages": [
          "signing",
          "key encipherment",
          "server auth",
          "client auth"
        ],
        "expiry": "8760h"
      }
    }
  }
}
```

1. 创建一个 JSON 配置文件，用于 CA 证书签名请求（CSR），例如：`ca-csr.json`。 确认用你需要的值替换掉尖括号中的值。

```json
{
  "CN": "kubernetes",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names":[{
    "C": "<country>",
    "ST": "<state>",
    "L": "<city>",
    "O": "<organization>",
    "OU": "<organization unit>"
  }]
}
```

1. 生成 CA 秘钥文件（`ca-key.pem`）和证书文件（`ca.pem`）：

```bash
../cfssl gencert -initca ca-csr.json | ../cfssljson -bare ca
```

1. 创建一个 JSON 配置文件，用来为 API 服务器生成秘钥和证书，例如：`server-csr.json`。 确认用你需要的值替换掉尖括号中的值。`MASTER_CLUSTER_IP` 是为 API 服务器 指定的服务集群 IP，就像前面小节描述的那样。 以下示例假定你的默认 DSN 域名为`cluster.local`。

```json
{
  "CN": "kubernetes",
  "hosts": [
    "127.0.0.1",
    "<MASTER_IP>",
    "<MASTER_CLUSTER_IP>",
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
  "names": [{
    "C": "<country>",
    "ST": "<state>",
    "L": "<city>",
    "O": "<organization>",
    "OU": "<organization unit>"
  }]
}
```

1. 为 API 服务器生成秘钥和证书，默认会分别存储为`server-key.pem` 和 `server.pem` 两个文件。

```bash
../cfssl gencert -ca=ca.pem -ca-key=ca-key.pem \
--config=ca-config.json -profile=kubernetes \
server-csr.json | ../cfssljson -bare server
```

## 分发自签名的 CA 证书

客户端节点可能不认可自签名 CA 证书的有效性。 对于非生产环境，或者运行在公司防火墙后的环境，你可以分发自签名的 CA 证书到所有客户节点，并刷新本地列表以使证书生效。
在每一个客户节点，执行以下操作：

```bash
sudo cp ca.crt /usr/local/share/ca-certificates/kubernetes.crt
sudo update-ca-certificates


Updating certificates in /etc/ssl/certs...
1 added, 0 removed; done.
Running hooks in /etc/ca-certificates/update.d....
done.
```

## 证书 API

你可以通过 `certificates.k8s.io` API 提供 x509 证书，用来做身份验证， 如[本文档](https://kubernetes.io/zh/docs/tasks/tls/managing-tls-in-a-cluster)所述。
