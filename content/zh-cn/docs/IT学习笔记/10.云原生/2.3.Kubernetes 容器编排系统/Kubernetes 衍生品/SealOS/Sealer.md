---
title: Sealer
---

# 概述

> 参考：
> - [GitHub 项目，sealerio/sealer](https://github.com/sealerio/sealer)

由 [sealos](https://github.com/labring/sealos) 作者进入阿里后开源的一款工具，用于将应用程序的所以依赖和 Kubernetes 打包成 ClusterImage。然后通过 ClusterImage 将此应用程序分发到任何位置。

# 问题

尴尬。。。快速开始都报错了。。。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/virv3w/1622037606609-84ee7ea6-d616-4ebe-ab7e-1defcc4905ce.png)

镜像目录结构

```bash
.
├── bin
│   ├── conntrack
│   ├── containerd-rootless-setuptool.sh
│   ├── containerd-rootless.sh
│   ├── crictl
│   ├── kubeadm
│   ├── kubectl
│   ├── kubelet
│   ├── nerdctl
│   └── seautil
├── cni
│   └── calico
│       └── calico.yaml.tmpl
├── cri
│   ├── containerd
│   ├── containerd-shim
│   ├── containerd-shim-runc-v2
│   ├── ctr
│   ├── docker
│   ├── dockerd
│   ├── docker-init
│   ├── docker-proxy
│   ├── rootlesskit
│   ├── rootlesskit-docker-proxy
│   ├── runc
│   └── vpnkit
├── etc
│   ├── 10-kubeadm.conf
│   ├── Clusterfile  # image default Clusterfile
│   ├── daemon.json
│   ├── docker.service
│   ├── kubeadm-config.yaml
│   └── kubelet.service
├── images
│   └── registry.tar  # registry docker image, will load this image and run a local registry in cluster
├── Kubefile
├── Metadata
├── README.md
├── registry # will mount this dir to local registry
│   └── docker
│       └── registry
├── scripts
│   ├── clean.sh
│   ├── docker.sh
│   ├── init-kube.sh
│   ├── init-registry.sh
│   ├── init.sh
│   └── kubelet-pre-start.sh
└── statics # yaml files, sealer will render values in those files
    └── audit-policy.yml

```
