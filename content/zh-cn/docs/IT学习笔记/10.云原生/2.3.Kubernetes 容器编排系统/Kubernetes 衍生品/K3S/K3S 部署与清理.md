---
title: K3S 部署与清理
---

# 概述

> 参考：
> - [官方文档，快速开始指南](https://docs.k3s.io/quick-start)

`curl -sfL https://get.k3s.io | sh -` 使用该脚本，可以自动创建用于运行 k3s 二进制文件的 service 文件，并通过 systemd 启动。

## systemd 管理的 Unit 文件

```shell
cat > /etc/systemd/system/k3s.service <<EOF
[Unit]
Description=Lightweight Kubernetes
Documentation=https://k3s.io
Wants=network-online.target
After=network-online.target

[Install]
WantedBy=multi-user.target

[Service]
Type=notify
EnvironmentFile=/etc/systemd/system/k3s.service.env
KillMode=process
Delegate=yes
# Having non-zero Limit*s causes performance problems due to accounting overhead
# in the kernel. We recommend using cgroups to do container-local accounting.
LimitNOFILE=1048576
LimitNPROC=infinity
LimitCORE=infinity
TasksMax=infinity
TimeoutStartSec=0
Restart=always
RestartSec=5s
ExecStartPre=-/sbin/modprobe br_netfilter
ExecStartPre=-/sbin/modprobe overlay
ExecStart=/usr/local/bin/k3s \
# 根据实际情况修改命令行标志
EOF
```

`EnvironmentFile` 和 `ExecStart` 这两个字段就是用来控制 k3s 运行行为的，其他字段值一般不用修改。

- `EnvironmentFile` 是 k3s 运行时读取的环境变量信息
- `ExecStart` 是 k3s 二进制文件位置，以及运行时标志

# 快速部署体验

获取安装脚本

```latex
curl -LO http://rancher-mirror.cnrancher.com/k3s/k3s-install.sh
```

第一个 master 执行

```latex
K3S_TOKEN=SECRET \
INSTALL_K3S_MIRROR=cn \
INSTALL_K3S_VERSION='v1.20.4+k3s1' \
INSTALL_K3S_EXEC='server --cluster-init' \
bash k3s-install.sh
```

其余 master 执行

```latex
K3S_TOKEN=SECRET \
INSTALL_K3S_MIRROR=cn \
INSTALL_K3S_VERSION='v1.20.4+k3s1' \
INSTALL_K3S_EXEC='server --server https://172.19.42.207:6443' \
bash k3s-install.sh
```

其余 node 执行

```latex
K3S_TOKEN=SECRET \
INSTALL_K3S_MIRROR=cn \
INSTALL_K3S_VERSION='v1.20.4+k3s1' \
K3S_URL='https://172.19.42.207:6443'
INSTALL_K3S_EXEC='agent' \
bash k3s-install.sh
```

# 离线部署高可用集群
 
安装脚将会执行如下逻辑：
- 生成用于 Systemd 的 Unit 文件
- TODO: 待完善

## 下载 k3s 二进制文件

```bash
export K3S_VERSION="v1.26.1"
wget https://github.com/k3s-io/k3s/releases/download/${K3S_VERSION}%2Bk3s1/k3s -O /usr/local/bin/k3s
chmod +x /usr/local/bin/k3s
```

## 下载安装脚本

```bash
curl -Lo install.sh https://get.k3s.io/
```

## 下载所需镜像

```bash
export ARCH="amd64"
wget https://github.com/k3s-io/k3s/releases/download/${K3S_VERSION}%2Bk3s1/k3s-airgap-images-${ARCH}.tar.gz
mkdir -p /var/lib/rancher/k3s/agent/images/
cp ./k3s-airgap-images-${ARCH}.tar /var/lib/rancher/k3s/agent/images/
```

## 配置第一个 master 节点

```bash
INSTALL_K3S_SKIP_DOWNLOAD=true \
K3S_TOKEN=SECRET \
INSTALL_K3S_EXEC='server --cluster-init --service-node-port-range 10000-60000' \
./install-zh.sh
```

## 配置其他 master 节点

```bash
INSTALL_K3S_SKIP_DOWNLOAD=true \
K3S_TOKEN=SECRET \
INSTALL_K3S_EXEC='server --server https://172.38.180.216:6443 --service-node-port-range 10000-60000' \
./install-zh.sh
```

## 配置工作节点

```bash
INSTALL_K3S_SKIP_DOWNLOAD=true \
K3S_TOKEN=SECRET \
INSTALL_K3S_EXEC='agent --server https://172.38.180.216:6443 ' \
./install-zh.sh
```

## 为所有节点配置通用信息

```shell
for cmd in kubectl crictl ctr; do
    ln -sf k3s /usr/local/bin/${cmd}
done
```

# 清理 K3S

K3S 的清理脚本来自于“安装脚本”中。

```bash
#!/bin/sh
set -x
[ $(id -u) -eq 0 ] || exec sudo $0 $@

/usr/local/bin/k3s-killall.sh

if command -v systemctl; then
    systemctl disable k3s
    systemctl reset-failed k3s
    systemctl daemon-reload
fi
if command -v rc-update; then
    rc-update delete k3s default
fi

rm -f /etc/systemd/system/k3s.service
rm -f /etc/systemd/system/k3s.service.env

remove_uninstall() {
    rm -f /usr/local/bin/k3s-uninstall.sh
}
trap remove_uninstall EXIT

if (ls /etc/systemd/system/k3s*.service || ls /etc/init.d/k3s*) >/dev/null 2>&1; then
    set +x; echo 'Additional k3s services installed, skipping uninstall of k3s'; set -x
    exit
fi

for cmd in kubectl crictl ctr; do
    if [ -L /usr/local/bin/$cmd ]; then
        rm -f /usr/local/bin/$cmd
    fi
done

rm -rf /etc/rancher/k3s
rm -rf /run/k3s
rm -rf /run/flannel
rm -rf /var/lib/rancher/k3s
rm -rf /var/lib/kubelet
rm -f /usr/local/bin/k3s
rm -f /usr/local/bin/k3s-killall.sh

if type yum >/dev/null 2>&1; then
    yum remove -y k3s-selinux
    rm -f /etc/yum.repos.d/rancher-k3s-common*.repo
elif type zypper >/dev/null 2>&1; then
    uninstall_cmd="zypper remove -y k3s-selinux"
    if [ "${TRANSACTIONAL_UPDATE=false}" != "true" ] && [ -x /usr/sbin/transactional-update ]; then
        uninstall_cmd="transactional-update --no-selfupdate -d run $uninstall_cmd"
    fi
    $uninstall_cmd
    rm -f /etc/zypp/repos.d/rancher-k3s-common*.repo
fi
```

# 部署 k3s 并修改参数

参考：[基杨文章](https://mp.weixin.qq.com/s/xpqZyoZltRkXcMQBcHos0Q)
