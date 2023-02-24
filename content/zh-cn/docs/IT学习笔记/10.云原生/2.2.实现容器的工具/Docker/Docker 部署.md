---
title: Docker 部署
---

# 概述

> 参考：
> - [官方文档](https://docs.docker.com/engine/install/)
> - [Centos 安装](https://docs.docker.com/engine/install/centos/)
> - [Ubuntu 安装](https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository)

# 安装 Docker 套件

## 方法 1：使用 Linux 的包管理器安装

### 使用包管理器安装

```bash
# centos
yum install -y yum-utils device-mapper-persistent-data lvm2
yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
yum install -y docker-ce

# ubuntu
sudo apt-get -y install apt-transport-https ca-certificates curl gnupg-agent software-properties-common
curl -fsSL https://repo.huaweicloud.com/docker-ce/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://repo.huaweicloud.com/docker-ce/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get -y update
sudo apt-get -y install docker-ce
```

> 若 centos8 提示无法安装 contained.io ,则执行如下指令安装
> yum install -y <https://mirrors.aliyun.com/docker-ce/linux/centos/7/x86_64/edge/Packages/containerd.io-1.2.13-3.2.el7.x86_64.rpm>

### 配置 Unit 文件

docker 自 1.13 版起会自动设置 iptables 的 FORWARD 默认策略为 DROP，这可能会影响 Kubernetes 集群依赖的报文转发功能，因此，需要在 docker 服务启动后，重新将 FORWARD 链的默认策略设备为 ACCEPT

- sed -i "14i ExecStartPost=/usr/sbin/iptables -P FORWARD ACCEPT" /usr/lib/systemd/system/docker.service

## 方法 2：直接安装二进制文件

以 20.10.6 版本为例

### 获取并部署二进制文件

在 <https://download.docker.com/linux/static/stable/x86_64/> 页面下载 [20.10.6](https://download.docker.com/linux/static/stable/x86_64/docker-20.10.6.tgz) 版本的二进制程序
解压并将二进制程序放到 $PATH 中

```bash
tar -zxvf docker-20.10.6.tgz
cp docker/* /usr/bin/
```

### 配置 Unit 文件

> 参考：
> - [官方文档，使用 systemd 配置守护进程](https://docs.docker.com/config/daemon/systemd/)

#### containerd.service

```bash
cat > /usr/lib/systemd/system/containerd.service << EOF
[Unit]
Description=containerd container runtime
Documentation=https://containerd.io
After=network.target local-fs.target

[Service]
ExecStartPre=-/sbin/modprobe overlay
ExecStart=/usr/bin/containerd

Type=notify
Delegate=yes
KillMode=process
Restart=always
RestartSec=5
# Having non-zero Limit*s causes performance problems due to accounting overhead
# in the kernel. We recommend using cgroups to do container-local accounting.
LimitNPROC=infinity
LimitCORE=infinity
LimitNOFILE=1048576
# Comment TasksMax if your systemd version does not supports it.
# Only systemd 226 and above support this version.
TasksMax=infinity
OOMScoreAdjust=-999

[Install]
WantedBy=multi-user.target
EOF
```

#### docker.service

```bash
cat > /usr/lib/systemd/system/docker.service <<EOF
[Unit]
Description=Docker Application Container Engine
Documentation=https://docs.docker.com
After=network-online.target firewalld.service containerd.service
Wants=network-online.target
Requires=containerd.service

[Service]
Type=notify
# the default is not to use systemd for cgroups because the delegate issues still
# exists and systemd currently does not support the cgroup feature set required
# for containers run by docker
ExecStartPost=/usr/sbin/iptables -P FORWARD ACCEPT
ExecStart=/usr/bin/dockerd --containerd=/run/containerd/containerd.sock
ExecReload=/bin/kill -s HUP $MAINPID
TimeoutSec=0
RestartSec=2
Restart=always

# Note that StartLimit* options were moved from "Service" to "Unit" in systemd 229.
# Both the old, and new location are accepted by systemd 229 and up, so using the old location
# to make them work for either version of systemd.
StartLimitBurst=3

# Note that StartLimitInterval was renamed to StartLimitIntervalSec in systemd 230.
# Both the old, and new name are accepted by systemd 230 and up, so using the old name to make
# this option work for either version of systemd.
StartLimitInterval=60s

# Having non-zero Limit*s causes performance problems due to accounting overhead
# in the kernel. We recommend using cgroups to do container-local accounting.
LimitNOFILE=infinity
LimitNPROC=infinity
LimitCORE=infinity

# Comment TasksMax if your systemd version does not support it.
# Only systemd 226 and above support this option.
TasksMax=infinity

# set delegate yes so that systemd does not reset the cgroups of docker containers
Delegate=yes

# kill only the docker process, not all processes in the cgroup
KillMode=process
OOMScoreAdjust=-500

[Install]
WantedBy=multi-user.target
EOF
```

### 安装 CLI 补全文件

我们可以在 [GitHub 项目，docker/cli 中的 contrib/completion](https://github.com/docker/cli/tree/master/contrib/completion ) 目录下找到各种 Shell 的 CLI 补全文件

```bash
curl https://raw.githubusercontent.com/docker/cli/master/contrib/completion/bash/docker -o /usr/share/bash-completion/completions/docker
```

# 配置并启动 Docker

## 添加 dockerd 配置文件

```bash
sudo mkdir -p /etc/docker
sudo cat > /etc/docker/daemon.json <<EOF
{
  "registry-mirrors": ["https://ac1rmo5p.mirror.aliyuncs.com"],
  "exec-opts": ["native.cgroupdriver=systemd"],
  "live-restore": true,
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "5m",
    "max-file": "5"
  },
  "storage-driver": "overlay2",
  "storage-opts": [
    "overlay2.override_kernel_check=true"
  ]
}
EOF
```

## 修改内核参数

```bash
cat > /etc/sysctl.d/docker.conf << EOF
net.ipv4.ip_forward = 1
EOF
sysctl -p /etc/sysctl.d/*
```

## 启动 docker

```bash
systemctl daemon-reload
systemctl enable docker --now
```
