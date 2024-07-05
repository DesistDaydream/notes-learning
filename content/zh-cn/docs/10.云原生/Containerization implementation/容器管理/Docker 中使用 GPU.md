---
title: Docker 中使用 GPU
linkTitle: Docker 中使用 GPU
date: 2024-07-05T08:42
weight: 20
---

# 概述

> 参考：
>
> -

### Docker 19.03，增加了对--gpus 选项的支持，我们在 docker 里面想读取 nvidia 显卡再也不需要额外的安装 nvidia-docker 了，下面开始实战

1. 安装 nvidia-container-runtime：

查看官网（[https://nvidia.github.io/nvidia-container-runtime](https://links.jianshu.com/go?to=https%3A%2F%2Fnvidia.github.io%2Fnvidia-container-runtime)）得知基于 RHEL 的发行版添加源的方式为：

```shell
# Centos
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-container-runtime/$distribution/nvidia-container-runtime.repo | \
  sudo tee /etc/yum.repos.d/nvidia-container-runtime.repo

# Ubuntu
curl -s -L https://nvidia.github.io/nvidia-container-runtime/gpgkey | \
  sudo apt-key add -
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia  -container-runtime/$distribution/nvidia-container-runtime.list | \
  sudo tee /etc/apt/sources.list.d/nvidia-container-runtime.list
sudo apt-get update
```

添加源后直接 yum 安装：

```shell
# centos
yum install nvidia-container-runtime

# Ubuntu
apt-get install nvidia-container-runtime
```

2. 安装 docker-19.03

在新主机上首次安装 Docker Engine-Community 之前，需要设置 Docker 存储库。之后，您可以从存储库安装和更新 Docker。

- 2.1 安装所需的软件包。yum-utils 提供了 yum-config-manager 效用，并 device-mapper-persistent-data 和 lvm2 由需要 devicemapper 存储驱动程序。

```shell
yum install -y yum-utils \
  device-mapper-persistent-data \
  lvm2
```

- 2.2 使用以下命令来设置稳定的存储库。

```csharp
yum-config-manager \
    --add-repo \
    https://download.docker.com/linux/centos/docker-ce.repo
```

- 2.3 开启 Docker 服务

```bash
systemctl start docker && systemctl enable docker
```

- 2.4 验证 docker 版本是否安装正常

```shell
$ docker version
Client: Docker Engine - Community
 Version:           19.03.3
 API version:       1.40
 Go version:        go1.12.10
 Git commit:        a872fc2f86
 Built:             Tue Oct  8 00:58:10 2019
 OS/Arch:           linux/amd64
 Experimental:      false

Server: Docker Engine - Community
 Engine:
  Version:          19.03.2
  API version:      1.40 (minimum version 1.12)
  Go version:       go1.12.8
  Git commit:       6a30dfc
  Built:            Thu Aug 29 05:27:34 2019
  OS/Arch:          linux/amd64
  Experimental:     false
 containerd:
  Version:          1.2.10
  GitCommit:        b34a5c8af56e510852c35414db4c1f4fa6172339
 runc:
  Version:          1.0.0-rc8+dev
  GitCommit:        3e425f80a8c931f88e6d94a8c831b9d5aa481657
 docker-init:
  Version:          0.18.0
  GitCommit:        fec3683
```

3. 启动容器

```bash
docker run -d  -it -p 1518:1518 --name="centos"  --gpus all nvidia/cuda:9.1-cudnn7-runtime-centos7 /bin/bash
# 启动导出器
docker run -d --name nvidia --restart always --gpus all -p 9400:9400 nvidia/dcgm-exporter:2.0.13-2.1.1-ubuntu18.04
```

进入容器

```bash
docker exec -it centos /bin/bash
```

查看显卡

```shell
$ nvidia-smi
Mon Oct 21 02:15:19 2019
+-----------------------------------------------------------------------------+
| NVIDIA-SMI 390.59                 Driver Version: 390.59                    |
|-------------------------------+----------------------+----------------------+
| GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp  Perf  Pwr:Usage/Cap|         Memory-Usage | GPU-Util  Compute M. |
|===============================+======================+======================|
|   0  GeForce GTX 108...  Off  | 00000000:00:08.0 Off |                  N/A |
| 29%   33C    P0    58W / 250W |      0MiB / 11178MiB |      0%      Default |
+-------------------------------+----------------------+----------------------+
|   1  GeForce GTX 108...  Off  | 00000000:00:09.0 Off |                  N/A |
| 29%   28C    P5    12W / 250W |      0MiB / 11178MiB |      0%      Default |
+-------------------------------+----------------------+----------------------+
|   2  GeForce GTX 108...  Off  | 00000000:00:0A.0 Off |                  N/A |
| 29%   27C    P5    12W / 250W |      0MiB / 11178MiB |      0%      Default |
+-------------------------------+----------------------+----------------------+
|   3  GeForce GTX 108...  Off  | 00000000:00:0B.0 Off |                  N/A |
| 29%   30C    P5    12W / 250W |      0MiB / 11178MiB |      0%      Default |
+-------------------------------+----------------------+----------------------+
|   4  GeForce GTX 108...  Off  | 00000000:00:0C.0 Off |                  N/A |
| 29%   31C    P0    58W / 250W |      0MiB / 11178MiB |      0%      Default |
+-------------------------------+----------------------+----------------------+
|   5  GeForce GTX 108...  Off  | 00000000:00:0D.0 Off |                  N/A |
| 29%   23C    P5    12W / 250W |      0MiB / 11178MiB |      0%      Default |
+-------------------------------+----------------------+----------------------+
|   6  GeForce GTX 108...  Off  | 00000000:00:0E.0 Off |                  N/A |
| 29%   27C    P5    12W / 250W |      0MiB / 11178MiB |      0%      Default |
+-------------------------------+----------------------+----------------------+
|   7  GeForce GTX 108...  Off  | 00000000:00:0F.0 Off |                  N/A |
| 29%   27C    P5    12W / 250W |      0MiB / 11178MiB |      3%      Default |
+-------------------------------+----------------------+----------------------+

+-----------------------------------------------------------------------------+
| Processes:                                                       GPU Memory |
|  GPU       PID   Type   Process name                             Usage      |
|=============================================================================|
|  No running processes found                                                 |
+-----------------------------------------------------------------------------+
```
