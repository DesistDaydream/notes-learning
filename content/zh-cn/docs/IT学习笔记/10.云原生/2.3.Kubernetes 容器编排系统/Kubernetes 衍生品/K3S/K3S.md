---
title: K3S
---

# 概述

> 参考：
> - [GitHub 项目，k3s-io/k3s](https://github.com/k3s-io/k3s)
> - [官方文档](https://rancher.com/docs/k3s/latest/en/)
>   - [中文官方文档](https://docs.rancher.cn/k3s/)
> - [公众号-云原生实验室，K3S 工具进阶完全指南](https://mp.weixin.qq.com/s/ARhxWGypG0wepMqwTLH0mQ)

K3S 是一个轻量的 Kubernetes，具有基本的 kubernetes 功能，将 kubernetes 的主要组件都集成在一个二进制文件中(apiserver、kubelet 等)，这个二进制文件只有不到 100m。内嵌 Containerd，可以通过 Containerd 来启动 coredns 等 kubernetes 的 addone。直接使用 k3s 的二进制文件，即可启动一个 kubernetes 的节点。

Note：K3S 的 kubelet 不支持 systemd 作为 cgroup-driver，原因详见 <https://github.com/rancher/k3s/issues/797>，说是 systemd 的类型无法放进二进制文件里。

k3s 二进制文件包含 kubelet、api-server、kube-controller-manager、kube-scheduler，然后会通过 containerd 拉起 coredns 与 flannel。

# k3s 关联文件与配置

**/etc/rancher/k3s/** #

- **./k3s.yaml** # kubeconfig 文件
- **./registries.yaml** #

**/run/k3s/** # K3S 所使用的容器 Runtime 的数据保存路径。

- **./containerd/** # 与 [Containerd](docs/IT学习笔记/10.云原生/2.2.实现容器的工具/Containerd/Containerd.md#Containerd%20关联文件与配置) 中的 /run/containerd/ 目录功能一致。

**/run/flannel/** # 与 [Flannel](docs/IT学习笔记/10.云原生/2.3.Kubernetes%20容器编排系统/8.Kubernetes%20网络/CNI/Flannel.md#Flannel%20关联文件与配置) 中 /run/flannel/ 目录功能一致。

**/var/lib/rancher/k3s/\*** # k3s 运行时数据存储保存路径

- **./agent/** # 作为 k8s 的 node 节点所需要的信息保存路径
  - 包括证书、containerd 数据目录、cni，containerd 的配置文件 等等都在此处
    - **./containerd/** # 与 [Containerd](docs/IT学习笔记/10.云原生/2.2.实现容器的工具/Containerd/Containerd.md#Containerd%20关联文件与配置) 中的 /var/lib/containerd/ 目录功能一致。
- **./data/** #
- **./server/** # 作为 k8s 的 master 节点所需要的信息保存路径
  - 包括证书、kube-system 名称空间中的 manifests 文件、etcd 数据 等等都在此处
  - **./db/** # 内嵌 Etcd 数据保存路径
  - **./manifests/** # k3s 集群启动后，kube-system 名称空间中 pod 的 manifests 文件
  - **./tls/** # Kubernetes 主要组件运行所需证书保存路径

**/var/lib/kubelet/** #

## k3s 卸载脚本中包含的文件

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

# 进入容器的文件系统

详见 [进入容器文件系统](docs/IT学习笔记/10.云原生/2.2.实现容器的工具/容器管理/容器运行时管理/进入容器文件系统.md)。在 k3s 中，如果是 containerd 的话，则是在 /run/k3s/containerd/ 目录代替 /run/containerd/ 目录

/run/k3s/containerd/io.containerd.runtime.v2.task/k8s.io/${ContainerID}/rootfs/