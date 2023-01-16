---
title: K3S部署与清理
---

# 概述

> 参考：
> - [**官方文档**](https://docs.rancher.cn/docs/k3s/quick-start/_index/)

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

# 部署高可用集群

## 下载 k3s 二进制文件

```latex
wget https://github.com/rancher/k3s/releases/download/v1.19.8+k3s1/k3s -O /usr/local/bin/k3s
chmod +x /usr/local/bin/k3s
```

### 配置第一个 master 节点

Unit 文件修改为

```shell
ExecStart=/usr/local/bin/k3s \
    server \
        '--cluster-init' \
```

配置环境变量文件

```latex
cat > /etc/systemd/system/k3s.service.env <<EOF
K3S_TOKEN=secret
EOF
```

### 配置其他 master 节点

Unit 文件修改为

```shell
ExecStart=/usr/local/bin/k3s \
    server \
        '--server' \
        'https://172.19.42.207:6443' \
```

配置环境变量文件

```latex
cat > /etc/systemd/system/k3s.service.env <<EOF
K3S_TOKEN=secret
EOF
```

### 配置工作节点

Unit 文件修改为

```shell
ExecStart=/usr/local/bin/k3s \
    agent
```

配置环境变量文件

```latex
cat > /etc/systemd/system/k3s.service.env <<EOF
K3S_TOKEN=secret
K3S_URL=https://172.19.42.207:6443
EOF
```

### 为所有节点配置通用信息

```shell
for cmd in kubectl crictl ctr; do
    ln -sf k3s /usr/local/bin/${cmd}
done
```

# 清理 k3s

```shell
#!/bin/sh
set -x
[ $(id -u) -eq 0 ] || exec sudo $0 $@

/usr/local/bin/k3s-killall.sh

if which systemctl; then
    systemctl disable k3s
    systemctl reset-failed k3s
    systemctl daemon-reload
fi
if which rc-update; then
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
fi
```

# 部署 k3s 并修改参数

参考：[**基杨文章**](https://mp.weixin.qq.com/s?__biz=MzU1MzY4NzQ1OA==&mid=2247485580&idx=1&sn=b474d3736caa99bb45cca993c455bd64&chksm=fbee4841cc99c1579ae6139b5959f9a9d5cceaee77b2335e1c30c96d19e003cfb1d2dca9ecff&scene=21#wechat_redirect)
