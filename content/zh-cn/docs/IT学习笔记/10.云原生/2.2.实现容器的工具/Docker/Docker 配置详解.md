---
title: Docker 配置详解
---

# 概述

> 参考：
> - 官方文档：<https://docs.docker.com/engine/reference/commandline/dockerd/#/linux-configuration-file>
> - <https://blog.csdn.net/u013948858/article/details/79974796>

Docker 的守护进程为 dockerd，dockerd 可以通过两种方式配置运行时行为

1. 通过配置文件 /etc/docker/daemon.json 进行配置
2. 使用 dockerd 命令的 flags 进行配置，可以将 flags 添加到 dockerd.service 中。

Note：

1. 配置文件中的配置，也可以通过 dockerd 的命令行参数(也就是 flags)指定，比如配置文件中的 data-root 字段，对应的 dockerd flags 为 --data-root STRING。

## 配置文件示例

dockerd 配置文件是 JSON 格式，基本常用的配置内容如下。

```json
{
  // 指定 docker pull 时，首先去连接的 registry。
  "registry-mirrors": [
    "http://172.38.40.180",
    "https://ac1rmo5p.mirror.aliyuncs.com"
  ],
  // 指定运行 docker 操作的不安全的 registry 的列表
  "insecure-registries": ["http://172.38.40.180"],
  // 指定 docker 运行时其他的选项，这里面指定 docker 的 cgroupdriver 为 systemd
  "exec-opts": ["native.cgroupdriver=systemd"],
  // 指定 docker 的日志驱动为 json-file
  "log-driver": "json-file",
  // 指定 docker 记录容器日志的参数，这里指定容器日志文件大小最大为100m
  "log-opts": {
    "max-size": "100m"
  },
  // 指定 docker 的存储驱动类型为 overlay2
  "storage-driver": "overlay2",
  // 指定 docker 存储驱动的其他选项
  "storage-opts": ["overlay2.override_kernel_check=true"]
}
```

# 命令行标志详解

参考[配置文件详解](/docs/IT学习笔记/10.云原生/2.2.实现容器的工具/Docker/Docker%20 配置详解.md 配置详解.md)，将配置文件中的字段转成 --OPTIONS 即可。

# 配置文件详解

## data-root: STRING

配置 docker info 命令中的 Docker Root Dir，也就是 docker 存储数据的路径。

## features: {}

一些新的特性可以通过配置该字段来启动或停止

## hosts: \[] # 指定 docker 守护进程监听的端口

可以从其他机器使用 docker -H URL 命令对该设备进行 docker 操作

## live-restore: BOOL

在 docker.service 守护程序停止期间，保持容器状态，说白了就是重启 docker 的时候 Containers 不重启。
开启该参数后，就算重启 dockerd 服务也不会更改 default-address-pools 参数执行的地址范围

## log-driver: STRING # 指定 docker 的日志驱动

## log-opts: {} # 指定 docker 记录容器日志的参数

## registry-mirrors: \[] # 指定 pull、push 镜像时候的加速器地址

## 下面是[官网给的配置文件中所有可用字段的说明](https://docs.docker.com/engine/reference/commandline/dockerd/#daemon-configuration-file)：

```json
{
  "authorization-plugins": [],
  "data-root": "",
  "dns": [],
  "dns-opts": [],
  "dns-search": [],
  "exec-opts": [],
  "exec-root": "",
  "experimental": false,
  "features": {},
  "storage-driver": "",
  "storage-opts": [],
  "labels": [],
  "live-restore": true,
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "5",
    "labels": "somelabel",
    "env": "os,customer"
  },
  "mtu": 0,
  "pidfile": "",
  "cluster-store": "",
  "cluster-store-opts": {},
  "cluster-advertise": "",
  "max-concurrent-downloads": 3,
  "max-concurrent-uploads": 5,
  "default-shm-size": "64M",
  "shutdown-timeout": 15,
  "debug": true,
  // 指定docker守护进程监听的端口，可以从其他机器使用docker -H URL命令对该设备进行docker操作
  "hosts": [],
  "log-level": "",
  "tls": true,
  "tlsverify": true,
  "tlscacert": "",
  "tlscert": "",
  "tlskey": "",
  "swarm-default-advertise-addr": "",
  "api-cors-header": "",
  "selinux-enabled": false,
  "userns-remap": "",
  "group": "",
  "cgroup-parent": "",
  "default-ulimits": {
    "nofile": {
      "Name": "nofile",
      "Hard": 64000,
      "Soft": 64000
    }
  },
  "init": false,
  "init-path": "/usr/libexec/docker-init",
  "ipv6": false,
  "iptables": false,
  "ip-forward": false,
  "ip-masq": false,
  "userland-proxy": false,
  "userland-proxy-path": "/usr/libexec/docker-proxy",
  "ip": "0.0.0.0",
  "bridge": "",
  // 指定 docker0 桥的 IP
  "bip": "",
  "fixed-cidr": "",
  "fixed-cidr-v6": "",
  "default-gateway": "",
  "default-gateway-v6": "",
  "icc": false,
  "raw-logs": false,
  "allow-nondistributable-artifacts": [],
  "registry-mirrors": [],
  "seccomp-profile": "",
  // 指定不安全仓库，docker 默认无法连接 http 协议的仓库，将仓库的 URL 添加到该字段后，docker 即可连接
  "insecure-registries": [],
  "no-new-privileges": false,
  "default-runtime": "runc",
  "oom-score-adjust": -500,
  "node-generic-resources": ["NVIDIA-GPU=UUID1", "NVIDIA-GPU=UUID2"],
  "runtimes": {
    "cc-runtime": {
      "path": "/usr/bin/cc-runtime"
    },
    "custom": {
      "path": "/usr/local/bin/my-runc-replacement",
      "runtimeArgs": ["--debug"]
    }
  },
  "default-address-pools": [
    { "base": "172.80.0.0/16", "size": 24 },
    { "base": "172.90.0.0/16", "size": 24 }
  ]
}
```
