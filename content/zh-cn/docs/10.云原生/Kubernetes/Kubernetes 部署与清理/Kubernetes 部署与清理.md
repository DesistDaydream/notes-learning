---
title: Kubernetes 部署与清理
linkTitle: Kubernetes 部署与清理
date: 2024-08-22T17:29
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，快速开始](https://kubernetes.io/docs/setup/)
> - [GitHub 项目，easzlab/kubeasz](https://github.com/easzlab/kubeasz)(ansible 部署项目)
> - [官方文档，入门-生产环境-使用部署工具安装 Kubernetes-使用 kubeadm 引导集群-使用 kubeadm 支持 IPv4 与 IPv6 双栈](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/dual-stack-support/)

注意事项：

- [不要使用 nftables](https://v1-19.docs.kubernetes.io/zh/docs/setup/production-environment/tools/kubeadm/install-kubeadm/#%E7%A1%AE%E4%BF%9D-iptables-%E5%B7%A5%E5%85%B7%E4%B8%8D%E4%BD%BF%E7%94%A8-nftables-%E5%90%8E%E7%AB%AF)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/qzanwh/1643338198264-a05be403-da81-4455-88f7-b7e613c7db20.png)

## Kubernetes 关联文件

下面这些是逐步总结的，应该是准确的，但是没有官方说明

**/etc/kubernetes/** # 系统组件运行时配置

**/var/lib/etcd/** # Etcd 数据目录

**/var/lib/kubelet/** # Kubelet 运行时配置及数据持久化目录

CNI 目录

- **/etc/cni/net.d/**# 默认配置文件保存目录
- **/opt/cni/bin/** # 默认 CNI 插件保存目录
- **/var/lib/cni/** # 默认 CNI 运行时产生的数据目录

# 部署 Kubernetes 集群

## 配置安装环境

- (可选)更新内核以解决 ipvs 的(在 k8s-1.11.0 版本)BUG
  - [参考内核部署文档](1.Linux%20Kernel.md Kernel.md)
- 关闭 iptables 和 firewalld 服务
  - 由于 kubernetes 的 kube-proxy 会接管防火墙生成相关规则，所以最好关闭系统自带的

```bash
systemctl stop firewalld.service && systemctl disable firewalld.service
systemctl stop iptables.service && systemctl disable iptables.service
```

- 关闭并禁用 SELinux

```bash
sed -i 's@^\(SELINUX=\).*@\1disabled@' /etc/selinux/config
setenforce 0
```

- 禁用 swap 设备，不禁用 swap，kubelet 会报错。关闭 swap 原因参考官方 issue。

```bash
swapoff -a
sed -i 's@^[^#]\(.*swap.*\)@#\1@g' /etc/fstab
```

- (可选)启用 ipvs 内核模块，若要使用 ipvs 模型的 proxy，各个节点需要载入 ipvs 相关的各模块

```bash
yum install ipvsadm -y
cat > /etc/modules-load.d/ip_vs.conf << EOF
ip_vs
ip_vs_lblc
ip_vs_lblcr
ip_vs_lc
ip_vs_nq
ip_vs_pe_sip
ip_vs_rr
ip_vs_sed
ip_vs_sh
ip_vs_wlc
ip_vs_wrr
EOF
for i in `cat /etc/modules-load.d/ip_vs.conf`; do modprobe ${i}; done
lsmod | grep ip_vs
```

- 正常部署完集群后如果没启用 ipvs 功能，则还需要修改 proxy 配置
  - kubectl edit configmap kube-proxy -n kube-system
    - 在 mode:字段后面加上“ipvs”

## 安装 Runtime

### 安装 Docker，不要用最新版

- [安装 docker-ce](https://developer.aliyun.com/article/110806)

```bash
yum install -y yum-utils device-mapper-persistent-data lvm2
yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
yum makecache fast
yum install -y docker-ce-19.03.11
```

- 配置 Docker
  - docker 自 1.13 版起会自动设置 iptables 的 FORWARD 默认策略为 DROP，这可能会影响 Kubernetes 集群依赖的报文转发功能，因此，需要在 docker 服务启动后，重新将 FORWARD 链的默认策略设备为 ACCEPT

```bash
sed -i "14i ExecStartPost=/usr/sbin/iptables -P FORWARD ACCEPT" /usr/lib/systemd/system/docker.service
```

- 配置 docker 参数的第三行标绿的为加速器的地址，可以自行定义，非官方推荐

```bash
mkdir /etc/docker
cat > /etc/docker/daemon.json <<EOF
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
mkdir -p /etc/systemd/system/docker.service.d
systemctl daemon-reload
systemctl restart docker && systemctl enable docker
```

- 修改内核参数，确保 docker 之上所定义的桥接网络可以正常对外通信

```bash
cat > /etc/sysctl.d/docker.conf << EOF
net.ipv4.ip_forward = 1
EOF
sysctl -p /etc/sysctl.d/*
```

### 安装 Containerd

注意：

- 如果想要让 Containerd 作为 Kubernetes 的 CRI，需要删除 Containerd 的配置文件(/etc/containerd/config.toml)，否则集群无法初始化，并报错:`getting status of runtime: rpc error: code = Unimplemented desc = unknown service runtime.v1alpha2.RuntimeService" , error: exit status 1`
- 或者让 config.toml 的配置完全符合 k8s 环境运行标准
- 查看 [Containerd 的 Versioning 与 Release](https://github.com/containerd/containerd/blob/main/RELEASES.md) 页面以获取与 Kubernetes 版本匹配的 Containerd 版本

## 安装 kubeadm、kubectl、kubelet

> 参考：
>
> - [官方文档，入门-生产环境-使用部署工具安装 Kubernetes-使用 kubeadm 引导集群-安装 kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/#installing-kubeadm-kubelet-and-kubectl)

安装 k8s 基本组件 kubelet，kubeadm，kubectl([软件源以及安装方法详见阿里源的 kubernetes 帮助](https://developer.aliyun.com/mirror/kubernetes?spm=a2c6h.13651102.0.0.3e221b11aiKpn2))

### CentOS

```bash
cat > /etc/yum.repos.d/kubernetes.repo <<EOF
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64/
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF
yum install -y kubelet kubeadm kubectl --disableexcludes=kubernetes
systemctl enable kubelet
```

### Ubuntu

```bash
sudo apt-get update && sudo apt-get install -y apt-transport-https gnupg2 curl
curl https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | apt-key add -
cat > /etc/apt/sources.list.d/kubernetes.list <<EOF
deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main
EOF
sudo apt-get update
sudo apt-get install -y kubectl
apt-mark hold kubelet kubeadm kubectl
```

apt-mark 用以保证 Ubuntu 自动更新时不更新 kubelet、kubeadm、kubectl

### 二进制文件

```bash
# 安装 CNI
export CNI_VERSION="v0.8.2"
export ARCH="amd64"
sudo mkdir -p /opt/cni/bin
curl -L "https://github.com/containernetworking/plugins/releases/download/${CNI_VERSION}/cni-plugins-linux-${ARCH}-${CNI_VERSION}.tgz" | sudo tar -C /opt/cni/bin -xz

# 安装 crictl
export CRICTL_VERSION="v1.24.1"
export ARCH="amd64"
curl -L --remote-name-all "https://github.com/kubernetes-sigs/cri-tools/releases/download/${CRICTL_VERSION}/crictl-${CRICTL_VERSION}-linux-${ARCH}.tar.gz"
tar -zxvf crictl-${CRICTL_VERSION}-linux-${ARCH}.tar.gz

# 安装 kubeadm、kubelet、kubectl
export DOWNLOAD_DIR=/usr/local/bin
sudo mkdir -p $DOWNLOAD_DIR
export RELEASE="$(curl -sSL https://dl.k8s.io/release/stable.txt)"
export ARCH="amd64"
cd $DOWNLOAD_DIR
sudo curl -L --remote-name-all https://storage.googleapis.com/kubernetes-release/release/${RELEASE}/bin/linux/${ARCH}/{kubeadm,kubelet,kubectl}

# 注意，这里面的 .sevice 文件不够准确，最好从已有机器拷贝过来
RELEASE_VERSION="v0.4.0"
curl -sSL "https://raw.githubusercontent.com/kubernetes/release/${RELEASE_VERSION}/cmd/kubepkg/templates/latest/deb/kubelet/lib/systemd/system/kubelet.service" | sed "s:/usr/bin:${DOWNLOAD_DIR}:g" | sudo tee /etc/systemd/system/kubelet.service
sudo mkdir -p /etc/systemd/system/kubelet.service.d
curl -sSL "https://raw.githubusercontent.com/kubernetes/release/${RELEASE_VERSION}/cmd/kubepkg/templates/latest/deb/kubeadm/10-kubeadm.conf" | sed "s:/usr/bin:${DOWNLOAD_DIR}:g" | sudo tee /etc/systemd/system/kubelet.service.d/10-kubeadm.conf
```

## 配置 kubelet 的 cgroup 驱动

> Note：参考[此处](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/#configure-cgroup-driver-used-by-kubelet-on-control-plane-node)，注意要与 docker 的 Cgroup Driver 一致。不同版本，设置方式不一样

1.18.0 之前的版本使用如下配置

    echo "KUBELET_EXTRA_ARGS=--cgroup-driver=systemd" > /etc/sysconfig/kubelet

## 使用 kubeadm 初始化 k8s 的 master 节点

> 参考：
>
> - [官方文档](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/)

准备初始化所需镜像

  - 由于国内没法访问国际，镜像又都在谷歌上，所以需要翻墙或者提前以某些方式下载下来。如果不知道要用哪些镜像，可以直接使用 **kubeadm config images list** 命令查看初始化时所需镜像，然后从其他镜像仓库下载下来，并使用 docker tag 命令更改镜像。可以使用阿里云的镜像仓库，地址如下：
    - registry.aliyuncs.com/k8sxio # XXX 为镜像名与版本号，zhangguanzhang 在阿里云创建的仓库

初始化集群的 master 节点

- 使用该命令进行初始化并对初始化的内容进行一些设定，如果可以没法翻墙，那么就会卡在 pull 的位置，从 google 拉不到 images(注意 IP 位置需要自己改成自己所需的 IP)
```bash
kubeadm init --kubernetes-version=v1.18.8 --pod-network-cidr=10.244.0.0/16 --image-repository="registry.aliyuncs.com/k8sxio"
```
- kubeadm 也可通过配置文件加载配置，以定制更丰富的部署选项，kubeadm-config.yaml 文件配置详见 《[kubeadm 命令行工具](/docs/10.云原生/2.3.Kubernetes%20 容器编排系统/Kubernetes%20 管理/kubeadm%20 命令行工具.md 管理/kubeadm 命令行工具.md)》
```bash
kubeadm init --config kubeadm-config.yaml
```
- 安装完成后如下所示，并获取到后续 node 加入集群的启动命令(注意保存该命令以便 node 节点加入 cluster 集群)
```bash
Your Kubernetes control-plane has initialized successfully!

To start using your cluster, you need to run the following as a regular user:

  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

You should now deploy a pod network to the cluster.
Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  https://kubernetes.io/docs/concepts/cluster-administration/addons/

You can now join any number of the control-plane node running the following command on each as root:

  kubeadm join 10.10.100.104:6443 --token 5b5a14.11sdexxuycp4rocs \
    --discovery-token-ca-cert-hash sha256:52431bdd96837cf25621123e90d3f97619715f08fec18f1f658ec4bacf8cd7ef \
    --control-plane --certificate-key 193f369208b46e494840395e04cd9fc779d46ba8288a64d621693251b81a3ae4

Please note that the certificate-key gives access to cluster sensitive data, keep it secret!
As a safeguard, uploaded-certs will be deleted in two hours; If necessary, you can use
"kubeadm init phase upload-certs --upload-certs" to reload certs afterward.

Then you can join any number of worker nodes by running the following on each as root:

kubeadm join 10.10.100.104:6443 --token 5b5a14.11sdexxuycp4rocs \
    --discovery-token-ca-cert-hash sha256:52431bdd96837cf25621123e90d3f97619715f08fec18f1f658ec4bacf8cd7ef
```

(测试环境功能)删除 master 节点的污点

- kubectl taint nodes --all node-role.kubernetes.io/master-

## 配置 kubectl 命令

不配置 kubectl 的话，在 master 节点执行 `kubectl get nodes` 命令，会反馈 `localhost:8080 connection refused` 错误,配置方法如下：

- 配置 kubectl 配置信息，并让 kubectl 支持 tab 补全命令功能

```bash
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
echo 'source <(kubectl completion bash)' >> ~/.bashrc
```

## 部署附加组件

安装网络组件 flannel

- `kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml`
## 添加节点 node

执行前面的 1,2,3 步配置环境，安装 docker 以及 k8s 相关组件。

```bash
kubeadm join 10.10.100.104:6443 --token 5b5a14.11sdexxuycp4rocs \
  --discovery-token-ca-cert-hash sha256:52431bdd96837cf25621123e90d3f97619715f08fec18f1f658ec4bacf8cd7ef
```

> [!Note]
> 若是手动安装 node，则需要在 master 上查看 csr kubectl get csr
>
> 然后再接受该请求，即可将 node 加入集群 kubectl certificate approve node-csr-An1VRgJ7FEMMF_uyy6iPjyF5ahuLx6tJMbk2SMthwLs

# 清理 Kubernetes 集群

```bash
kubeadm reset -f
modprobe -r ipip
lsmod
rm -rf ~/.kube/
rm -rf /etc/kubernetes/
rm -rf /etc/cni
rm -rf /var/lib/etcd
rm -rf /var/etcd
#最好不要清理cni的二进制文件，再次安装可能无法生成
rm -rf /opt/cni
```

# 升级 Kubernetes 集群

官方文档：<https://kubernetes.io/docs/tasks/administer-cluster/kubeadm/kubeadm-upgrade/>

## 升级所有节点 kubeadm

- yum install -y kubeadm-1.19.x-0 --disableexcludes=kubernetes

## 升级 master

升级第一个 master 节点

1. 确认版本信息是否为升级目标版本
   1. kubeadm version
2. 排空当前节点
   1. kubectl drain --ignore-daemonsets \<Node-Name>
3. 开始升级
   1. kubeadm upgrade apply v1.19.x
      1. 若想更新 kubedm-config 配置，指定 --config 参数
         1. kubeadm upgrade apply 1.19.2 --config=kubeadm-config.yaml
4. 取消不可调度转台
   1. kubectl uncordon \<Node-Name>

升级其他 master 节点

1. kubeadm upgrade node

升级所有 master 的 kubectl 与 kubelet 组件

1. yum install -y kubelet-1.19.x-0 kubectl-1.19.x-0 --disableexcludes=kubernetes
2. 逐一重启 kubelet
3. systemctl daemon-reload && systemctl restart kubelet

## 升级 node

排空节点

1. kubectl drain \<Node-Name> --ignore-daemonsets

升级 kubelet 配置

1. kubeadm upgrade node

升级 kubelet 和 kubectl

1. yum install -y kubelet-1.19.x-0 kubectl-1.19.x-0 --disableexcludes=kubernetes
2. 逐一重启 kubelet
3. systemctl daemon-reload && systemctl restart kubelet

取消不可调度转台

1. kubectl uncordon \<Node-Name>

# Dual-stack(双栈)

> 参考：
>
> - [官方文档，快速开始-生产环境-使用部署工具安装 Kubernetes-使用 kubeadm 引导集群-使用 kubeadm 支持双栈](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/dual-stack-support/)
> - [官方文档，任务-网络-验证 IPv4/IPv6 双栈](https://kubernetes.io/docs/tasks/network/validate-dual-stack)

想要创建一个双栈集群，需要注意如下几点：

- pod-network-cidr
- service-cidr
- node-ip
- 在各个 Kubernetes 系统组件配置中的上述三个参数需要填写 IPv4 与 IPv6，以逗号分割
- certSANs 中添加 IPv6

一个初始化的 kubeadm-config.yaml 效果如下：

```yaml
apiVersion: kubeadm.k8s.io/v1beta3
kind: InitConfiguration
nodeRegistration:
  criSocket: unix:///run/containerd/containerd.sock
  kubeletExtraArgs:
    node-ip: 172.19.42.231,2001:db8:0:1::231
---
apiVersion: kubeadm.k8s.io/v1beta3
kind: ClusterConfiguration
networking:
  podSubnet: 10.244.0.0/16,2001:db8:42:0::/56
  serviceSubnet: 10.96.0.0/12,2001:db8:42:1::/112
apiServer:
  certSANs:
    - localhost
    - 127.0.0.1
    - k8s-api.bj-test.datalake.cn
    - 172.19.42.230
    - 2001:db8:0:1::230
---
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
---
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
```

一个 Master 节点加入集群的 kubeadm-config.yaml 效果如下：

```yaml
apiVersion: kubeadm.k8s.io/v1beta3
kind: JoinConfiguration
controlPlane:
  certificateKey: 248c2dea22ef4598498c265311209823884e6767ef59d5b266344b0cdd03e304
  localAPIEndpoint:
    advertiseAddress: "172.19.42.232"
    bindPort: 6443
discovery:
  bootstrapToken:
    apiServerEndpoint: k8s-api.bj-test.datalake.cn:6443
    token: 27wls7.3spnpsp660kekn65
    caCertHashes:
      - sha256:2ae962435b870771435e5afc8f3ede6728710db49ce1ca5598d861879868a3c7
nodeRegistration:
  criSocket: unix:///run/containerd/containerd.sock
  kubeletExtraArgs:
    node-ip: 172.19.42.232,2001:db8:0:1::232
```

Node 节点加入集群的 kubeadm-config.yaml 文件效果如下：

```yaml
apiVersion: kubeadm.k8s.io/v1beta3
kind: JoinConfiguration
discovery:
  bootstrapToken:
    apiServerEndpoint: k8s-api.bj-test.datalake.cn:6443
    token: 27wls7.3spnpsp660kekn65
    caCertHashes:
      - sha256:2ae962435b870771435e5afc8f3ede6728710db49ce1ca5598d861879868a3c7
nodeRegistration:
  criSocket: unix:///run/containerd/containerd.sock
  kubeletExtraArgs:
    node-ip: 172.19.42.234,2001:db8:0:1::234
```

## 验证集群双栈

验证是否配置了 Pod 的双栈地址范围

```bash
~]# for i in $(kubectl get node -oname); do kubectl get $i -o go-template --template='{{range .spec.podCIDRs}}{{printf "%s\n" .}}{{end}}'; done
10.244.0.0/24
2001:db8:42::/64
10.244.2.0/24
2001:db8:42:3::/64
10.244.1.0/24
2001:db8:42:2::/64
10.244.3.0/24
2001:db8:42:4::/64
10.244.4.0/24
2001:db8:42:5::/64
10.244.5.0/24
2001:db8:42:6::/64
```

验证是否检测到双栈接口

```bash
~]# for i in $(kubectl get node -oname); do kubectl get $i -o go-template --template='{{range .status.addresses}}{{printf "%s: %s\n" .type .address}}{{end}}'; done
InternalIP: 172.19.42.231
InternalIP: 2001:db8:0:1::231
Hostname: test-node-1
InternalIP: 172.19.42.232
InternalIP: 2001:db8:0:1::232
Hostname: test-node-2
InternalIP: 172.19.42.233
InternalIP: 2001:db8:0:1::233
Hostname: test-node-3
InternalIP: 172.19.42.234
InternalIP: 2001:db8:0:1::234
Hostname: test-node-4
InternalIP: 172.19.42.235
InternalIP: 2001:db8:0:1::235
Hostname: test-node-5
InternalIP: 172.19.42.236
InternalIP: 2001:db8:0:1::236
Hostname: test-node-6
```

验证 Pod 寻址

```bash
~]# kubectl get pods pod01 -o go-template --template='{{range .status.podIPs}}{{printf "%s\n" .ip}}{{end}}'

```
