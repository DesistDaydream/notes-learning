---
title: config 子命令
---

# 概述

config 子命令用于控制 [User Account(KubeConfig)](/docs/10.云原生/Kubernetes/API%20访问控制/Authentication(认证)/User%20Account(KubeConfig).md) 的 KubeConfig 文件

# Syntax(语法)

**kubectl config SUBCOMMAND \[options]

SUBCOMMAND 包括：

- **current-context** # 显示当前上下文 Displays the current-context
- **delete-cluster** # Delete the specified cluster from the kubeconfig
- **delete-context** # 从 kubeconfig 文件中删除指定的上下文
- **get-clusters** # 显示在 kubeconfig 中已经定义的 cluster 信息。Display clusters defined in the kubeconfig
- **get-contexts** # 显示在 kubeconfig 中的上下文列表。每行的信息表示包括以\*表示当前使用的 context，context 名称，cluster 与 authinfo(认证信息即用户名)，名称空间
- **rename-context**# Renames a context from the kubeconfig file.
- **set** # 在 KubeConfig 文件中设置一个单独的值。Sets an individual value in a kubeconfig file
- **set-cluster** # 在 kubeconfig 中设定集群条目。
- **set-context**# 在 kubeconfig 中设定上下文条目。Sets a context entry in kubeconfig
- **set-credentials** # 在 kubeconfig 中设定用户凭证。
- **unset**# 取消在 KubeConfig 文件中设置的一个单独的值。Unsets an individual value in a kubeconfig file
- **use-context**# 在 kubeconfig 中设定当前上下文(即使用哪个用户操作客户端)。
- **view**# 显示已经合并的 KubeConfig 文件或一个指定的 KubeConfig 文件。Display merged kubeconfig settings or a specified kubeconfig file。

## OPTIONS

- --kubeconfig=/PATH/FILE # 指明要操作的 KubeConfig 文件

# SubCommand(子命令)

## set-cluster # 在 kubeconfig 文件中设置一个集群条目。Sets a cluster entry in kubeconfig

**kubectl config set-cluster NAME \[--server=server] \[--certificate-authority=PATH/TO/CERTIFICATE/AUTHORITY] \[--insecure-skip-tls-verify=true] \[OPTIONS]** #

OPTIONS

- **--embed-certs={false|true}** # 将--certificate-authority 中指定的证书嵌入 kubeconfig 文件中。i.e.将证书内容使用 base64 编码后存入，默认为 false，即不读取内容进行编码，而是直接将路径写到 kubeconfig 文件中

EXAMPLE

```bash
kubectl config set-cluster kubernetes \
--certificate-authority=/etc/kubernetes/pki/ca.crt \
--embed-certs=true \
--server=192.168.10.10:6443 \
--kubeconfig=./lch-config # 为 lch-config 的 kubeconfig 文件设定集群信息，指定证书为/etc/kubernetes/pki/ca.crt，开启嵌入式认证，指定集群 api-server 的 ip 和 port。
```

## set-context # 设定上下文，确立 user 与 cluster 的绑定关系与上下文的 name

**kubectl config set-context \[NAME | --current] \[--cluster=cluster_nickname] \[--user=user_nickname] \[--namespace=namespace] \[options]**

EXAMPLE

- kubectl config set-context lch@kubernetes --cluster=kubernetes --user=lch
- kubectl config set-context dashboard-admin@kubernetes --cluster=kubernetes --user=dashboard-admin --kubeconfig=/root/dashbord-admin.conf #

## set-credentials # 在 kubeconfig 中设置凭证，即设置用户的认证，以便让 kubernetes 集群认识到该用户。i.e.创建 User Account

**kubectl config set-credentials NAME \[--client-certificate=Path/to/certfile] \[--client-key=Path/to/keyfile] \[--token=bearer_token] \[--username=BasicUser] \[--password=BasicPassword] \[--auth-provider=provider_name] \[--auth-provider-arg=key=value] \[OPTIONS]**

通过证书、token、用户密码或者认证提供者键值对来设定一个用户凭证。credential 的意思就是“一个人的背景的资格，成就，个人品质或方面，通常用于表明他们适合某事”。在这里就是这个 user 的证书或者 token 等以便在与集群交互时进行验证

OPTIONS

- **--embed-certs=ture|false** # 在 kubeconfig 中嵌入证书/私钥，即变成非明文的方式储存

EXAMPLE

通过使用证书与私钥的方式设定名为 lch 的用户

- `kubectl config set-credentials lch --client-certificate=./lch.crt --client-key=./lch.key --embed-certs`

使用 ${DASH_TOKEN} 中的 token 来作为 user 的凭证

- `kubectl config set-credentials dashboard-admin --token=$DASH_TOKEN --kubeconfig=/root/dashbord-admin.conf`

## use-context # 设置当前 current-context 字段(当前所用的使用的上下文)

**kubectl config use-context CONTEXT_NAME \[OPTIONS]**
EXAMPLE

- kubectl config use-context dashboard-admin@kubernetes --kubeconfig=/root/dashbord-admin.conf #

## view # 显示一个 kubeconfig 文件的信息

**kubectl config view \[FLAGS] \[OPTIONS]**

```yaml
~]# kubectl config view
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: DATA+OMITTED
    server: https://k8s-api.bj-net.ehualu.local:6443
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: kubernetes-admin
  name: kubernetes-admin@kubernetes
current-context: kubernetes-admin@kubernetes
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data: REDACTED
    client-key-data: REDACTED
```

OPTIONS

- **--raw** # 显示原始字节数据
- **--minify** # 只显示当前 context 的信息。

EXAMPLE

- 显示用户名为 user-2c2f24ck5f 的证书数据
  - **kubectl config view -o jsonpath='{.users\[?(@.name == "user-2c2f24ck5f")].user.client-certificate-data}' --raw**
