---
title: K3S 部署与清理
weight: 2
---

# 概述

> 参考：
> 
> - [官方文档，快速开始指南](https://docs.k3s.io/quick-start)
> - [公众号-CNCF，利用 kube-vip 实现 K3s 高可用部署](https://mp.weixin.qq.com/s/Qe3oImSUJ1xFCsfXsUMdmA)

`curl -sfL https://get.k3s.io | sh -` 使用该脚本，可以自动创建用于运行 k3s 二进制文件的 service 文件，并通过 systemd 启动。

注意：

- k3s 会自动部署 servicelb 服务，该服务与 kube-proxy 的 ipvs 模式冲突。

# 快速部署体验

获取安装脚本

```bash
curl -LO http://rancher-mirror.cnrancher.com/k3s/k3s-install.sh
```

第一个 master 执行

```bash
K3S_TOKEN=SECRET \
INSTALL_K3S_MIRROR=cn \
INSTALL_K3S_VERSION='v1.20.4+k3s1' \
INSTALL_K3S_EXEC='server --cluster-init' \
bash k3s-install.sh
```

其余 master 执行

```bash
K3S_TOKEN=SECRET \
INSTALL_K3S_MIRROR=cn \
INSTALL_K3S_VERSION='v1.20.4+k3s1' \
INSTALL_K3S_EXEC='server --server https://172.19.42.207:6443' \
bash k3s-install.sh
```

其余 node 执行

```bash
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
export K3S_VERSION="v1.26.2"
wget https://github.com/k3s-io/k3s/releases/download/${K3S_VERSION}%2Bk3s1/k3s -O /usr/local/bin/k3s
chmod +x /usr/local/bin/k3s
```

## 下载安装脚本

```bash
curl -Lo install.sh https://get.k3s.io/
chmod 755 install.sh
```

## 下载所需镜像

```bash
export ARCH="amd64"
wget https://github.com/k3s-io/k3s/releases/download/${K3S_VERSION}%2Bk3s1/k3s-airgap-images-${ARCH}.tar.gz
mkdir -p /var/lib/rancher/k3s/agent/images/
cp ./k3s-airgap-images-${ARCH}.tar /var/lib/rancher/k3s/agent/images/
```

## mater 节点

### (可选)使用 kube-vip 提供 VIP

> 参考：
> 
> - [kube-vip 官方文档，K3S](https://kube-vip.io/docs/usage/k3s/)

### 为所有 mater 节点生成配置文件

```bash
mkdir -p /etc/rancher/k3s

tee /etc/rancher/k3s/config.yaml > /dev/null <<EOF
service-node-port-range: "30000-40000"
disable-cloud-controller: true
disable:
  - servicelb
  - traefik
kube-controller-manager-arg:
  - bind-address=0.0.0.0
kube-scheduler-arg:
  - bind-address=0.0.0.0
kube-proxy-arg:
  - proxy-mode=ipvs
  - masquerade-all=true
  - metrics-bind-address=0.0.0.0
etcd-expose-metrics: true
EOF
```

### 初始化第一个 master 节点

```bash
INSTALL_K3S_SKIP_DOWNLOAD=true \
K3S_TOKEN=SECRET \
INSTALL_K3S_EXEC='server --cluster-init' \
./install-zh.sh
```

### 加入其他 master 节点

```bash
INSTALL_K3S_SKIP_DOWNLOAD=true \
K3S_TOKEN=SECRET \
INSTALL_K3S_EXEC='server --server https://172.38.180.216:6443' \
./install-zh.sh
```

## work 节点

### 为所有 work 节点生成配置文件

### 加入 work 节点

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

## 配置 nerdctl

重点是指定 containerd 信息

```bash
tee /etc/nerdctl/nerdctl.toml > /dev/null <<EOF
address        = "unix:///run/k3s/containerd/containerd.sock"
namespace      = "k8s.io"
EOF
```

# Systemd 管理的 Unit 文件

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
# 根据实际情况修改命令行标志。更推荐使用 /etc/rancher/k3s/config.yaml 文件而避免使用 k3s 命令行参数。
EOF
```

`EnvironmentFile` 和 `ExecStart` 这两个字段就是用来控制 k3s 运行行为的，其他字段值一般不用修改。

- `EnvironmentFile` 是 k3s 运行时读取的环境变量信息
- `ExecStart` 是 k3s 二进制文件位置，以及运行时标志

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

# Ansible 部署 K3S

> 参考： 
> 
> - [GitHub 项目，techno-tim/k3s-ansible](https://github.com/techno-tim/k3s-ansible)
>   - 起源于：
>     - https://github.com/k3s-io/k3s-ansible
>     - https://github.com/geerlingguy/turing-pi-cluster
>     - https://github.com/212850a/k3s-ansible

ansible-playbook site.yml -i inventory/my-cluster/hosts.ini

# 常见问题

1.25.7 和 1.26.2 版本的 K3S 使用的 v0.21.1 版本的 Flannel 与旧版 iptables（v1.4.21，比如 centos7~凸(艹皿艹 )） 不兼容  ，需要使用 v0.21.4 Flannel，可以升级到 1.26.3 版本解决。[issue#7096](https://github.com/k3s-io/k3s/issues/7096)