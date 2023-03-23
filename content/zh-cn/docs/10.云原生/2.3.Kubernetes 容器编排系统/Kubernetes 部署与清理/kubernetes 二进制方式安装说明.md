---
title: kubernetes 二进制方式安装说明
---

首先需要明确几个概念，kubelet负责本节点容器的生命周期管理，那么在master节点上，如果只通过二进制文件运行 apiserver、controller-manager、scheduler，则无需部署 kubelet

所以，按照大体可以划分这么几块：

1. etcd节点：etcd(可部署在master节点上)。高可用至少需要3台设备

2. master节点：apiserver、controller-manager、scheduler。高可用至少需要2台设备。apiserver 需要与 etcd 进行交互

3. node节点：kubelet、CNI插件、kube-proxy
