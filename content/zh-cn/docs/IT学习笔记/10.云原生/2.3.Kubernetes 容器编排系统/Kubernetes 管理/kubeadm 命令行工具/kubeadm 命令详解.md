---
title: kubeadm 命令详解
---

# 概述

# Syntax(语法)

**kubeadm \[command]**
Command 包括：

1. alpha Kubeadm experimental sub-commands
2. completion Output shell completion code for the specified shell (bash or zsh).
3. config Manage configuration for a kubeadm cluster persisted in a ConfigMap in the cluster.
4. help Help about any command
5. init Run this command in order to set up the Kubernetes master.
6. join Run this on any machine you wish to join an existing cluster
7. reset Run this to revert any changes made to this host by 'kubeadm init' or 'kubeadm join'.
8. token Manage bootstrap tokens.
9. upgrade Upgrade your cluster smoothly to a newer version with this command.
10. version Print the version of kubeadm

## kubeadm alpha

续订运行控制计划所需的所有已知证书。无论到期日期如何，都会无条件地续订续订。续订也可以单独运行以获得更多控制权。

EXAMPLE

1. kubeadm alpha certs renew all #更新集群中所有证书

## kubeadm token \[COMMAND]

COMMAND 包括

1. create Create bootstrap tokens on the server.
2. delete Delete bootstrap tokens on the server.
3. generate Generate and print a bootstrap token, but do not create it on the server.
4. list List bootstrap tokens on the server.

kubeadm token create \[token]

EXAMPLE

- kubeadm token create --print-join-command #创建 node 节点加入 master 命令
- kubeadm token list #列出所有可以引导的令牌（i.e.join 时所用的 token）
  - 可以通过以下命令来获取 master 上 CA 证书的 hash 值(i.e.join 时所用的--discovery-token-ca-cert-hash 的值)，然后根据 list 列出的 token，与 ca 的 hash 值合在一起，就可以得到 join 时所用的相关参数
  - openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex

## kubeadm config \[COMMAND]

COMMAND 包括：

- images Interact with container images used by kubeadm.
- migrate Read an older version of the kubeadm configuration API types from a file, and output the similar config object for the newer version.
- print-default Print the default values for a kubeadm configuration object.
- upload Upload configuration about the current state, so that 'kubeadm upgrade' can later know how to configure the upgraded cluster.
- view View the kubeadm configuration stored inside the cluster.

kubeadm config images \[list | pull] \[flags]

EXAMPLE

1. kubeadm config images list --kubernetes-version=1.12.1 #列出 k8s 的 1.12.1 版本所需的所有 images 以及版本号

## kubeadm init \[Command] \[Flags]

Available Commands

1. phase # 使用该子命令来执行 kubeadm 初始化流程的单个阶段

Flags

1. \--kub -network-cidr IP #用于指定分 Pod 分配使用的网络地址，它通常应该与要部署使用的网络插件（例如 flannel、calico 等）的默认设定保持一致，10.244.0.0/16 是 flannel 默认使用的网络
2. \--service-cidr IP #用于指定为 Service 分配使用的网络地址(i.e.cluster-ip)，它由 kubernetes 管理，默认即为 10.96.0.0/12
3. \--ignore-preflight-errors=Swap #仅应该在未禁用 Swap 设备的状态下使用。

Example

- kubeadm init --config=kubeadm-config.yaml #使用 kubeadm-config.yaml 文件初始化 k8s 集群的 master
  - 以下是个符合前述命令设定方式的使用示例，不过，它明确定义了 kubeProxy 的模式为 ipvs，并支持通过修改 imageRepository 的值修改获取系统镜像时使用的镜像仓库。需要建立 FileName.conf 文件,内容说明详见 k8s 官网<https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init/>
- kubeadm init phase kubeconfig all #该命令可以生成所有 kubeconfig 文件。在误删除了 admin.conf 文件后，可以用该命令工薪生成

## kubeadm join \[Command] \[Flags]

FLAGS

- \--experimental-control-plane #在此节点上创建新的控制平面实例，等同于把该节点加入到 Cluster 后成为 master。
- \--token #
- \--discovery-token-ca-cert-hash #

EXAPMLE

- kubeadm join 192.168.10.10:6443 --token j04n3m.octy8zely83cy2ts --discovery-token-ca-cert-hash sha256:84938d2a22203a8e56a787ec0c6ddad7bc7dbd52ebabc62fd5f4dbea72b14d1f --experimental-control-plane #等同于把该节点加入到 Cluster 后成为 master
