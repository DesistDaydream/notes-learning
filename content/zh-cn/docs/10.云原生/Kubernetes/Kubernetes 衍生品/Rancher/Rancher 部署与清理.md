---
title: Rancher 部署与清理
linkTitle: Rancher 部署与清理
weight: 20
---

# 概述

> 参考：
>
> - 常见问题：<https://www.bookstack.cn/read/rancher-2.4.4-zh/%E5%B8%B8%E8%A7%81%E9%97%AE%E9%A2%98.md>

## 快速部署体验

```shell
docker run -d --name=rancher-server --restart=unless-stopped \
  -p 60080:80 -p 60443:443 \
  -v /opt/rancher:/var/lib/rancher \
  --privileged \
  rancher/rancher:v2.5.3
```

如果想让 rancher 可以验证外部 https 的自建 CA 证书，需要在启动前将证书导入 rancher-server 中，效果如下：

参考链接： <https://rancher.com/docs/rancher/v2.x/en/installation/resources/chart-options/#additional-trusted-cas>

<https://rancher.com/docs/rancher/v2.x/en/installation/other-installation-methods/single-node-docker/advanced/#custom-ca-certificate>

```shell
docker run -d --name=rancher-server --restart=unless-stopped \
  -p 60080:80 -p 60443:443 \
  -v /opt/rancher:/var/lib/rancher \
  -v /host/certs:/container/certs \
  -e SSL_CERT_DIR="/container/certs" \
  --privileged \
  rancher/rancher:latest
```

在宿主机的 /host/certs 目录中存放要导入的证书，比如可以把 harbor 的证书与私钥放入该目录中，这样 rancher 就可以添加 https 的 harbor 仓库了

还可以传递 CATTLE_SYSTEM_DEFAULT_REGISTRY 环境变量，让 rancher 内部使用私有镜像地址。比如

    docker run -d --name=rancher-server --restart=unless-stopped \
      -p 60080:80 -p 60443:443 \
      -v /opt/rancher:/var/lib/rancher \
      --privileged \
      -e CATTLE_SYSTEM_DEFAULT_REGISTRY="registry.wx-net.ehualu.com" \
      rancher/rancher:latest

## 高可用部署

Rancher 的高可用本质上就是将 Rancher 作为 kubernetes 上的一个服务对外提供(这个 k8s 集群通常只用来运行 Rancher)。Rancher 的数据储存在 k8s 集群的后端存储中(i.e.ETCD)。由于原生 k8s 部署负责，资源需求大，不易维护等问题，Rancher 官方推出了一个简化版的 k8s，即 [k3s](https://github.com/rancher/k3s)。k3s 是一个简化版的 k8s，可以实现基本的 k8s 功能，但是部署更简单，资源需求更少，更易维护。[k3s](https://github.com/rancher/k3s) 介绍详见 k3s 介绍

### (可选)启动 k3s 集群

启动 mysql 用于 k3s 存储数据

```bash
docker run -d  --name k3s-mysql  --restart=always  \
-v /opt/k3s-cluster/mysql/conf:/etc/mysql/conf.d \
-v /opt/k3s-cluster/mysql/data:/var/lib/mysql \
-e MYSQL_ROOT_PASSWORD=root \
-p 3306:3306 \
mysql:5.7.29 --default-time-zone=+8:00 \
--character-set-server=utf8mb4 --collation-server=utf8mb4_general_ci \
--explicit_defaults_for_timestamp=true --lower_case_table_names=1 --max_allowed_packet=128M
```

启动 k3s

```bash
curl -sfL https://docs.rancher.cn/k3s/k3s-install.sh | INSTALL_K3S_MIRROR=cn sh -s - server \
--docker --datastore-endpoint="mysql://root:root@tcp(172.38.40.212:3306)/k3s"
# 配置 kubectl 的 kubeconfig 文件
mkdir ~/.kube
cp /etc/rancher/k3s/k3s.yaml /root/.kube/config
echo "source <(kubectl completion bash)" >> ~/.bashrc
```

### 部署 Rancher

创建证书参考：[**自建 CA 脚本**](https://thoughts.teambition.com/workspaces/5f90e312c800160016ea22fb/docs/5fa4f848eaa1190001257bba)，该脚本参考：<https://docs.rancher.cn/docs/rancher2/installation/options/self-signed-ssl/_index>

```bash
# 创建用于运行 Rancher 的名称空间
kubectl create namespace cattle-system
# 创建CA证书
openssl genrsa -out ca.key 4096

openssl req -x509 -new -nodes -sha512 -days 36500 \
  -subj "/C=CN/ST=Tianjin/L=Tianjin/O=eHualu/OU=Operations/CN=k3s-rancher.desistdaydream.ltd" \
  -key ca.key \
  -out ca.crt
openssl genrsa -out k3s-rancher.desistdaydream.ltd.key 4096

openssl req -sha512 -new \
  -subj "/C=CN/ST=Tianjin/L=Tianjin/O=eHualu/OU=Operations/CN=k3s-rancher.desistdaydream.ltd" \
  -key k3s-rancher.desistdaydream.ltdn.key \
  -out k3s-rancher.desistdaydream.ltd.csr

cat > v3.ext <<EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1=k3s-rancher.desistdaydream.ltd
EOF

openssl x509 -req -sha512 -days 36500 \
  -extfile v3.ext \
  -CA ca.crt -CAkey ca.key -CAcreateserial \
  -in k3s-rancher.desistdaydream.ltd.csr \
  -out k3s-rancher.desistdaydream.ltd.crt

# 将证书添加到 secret 资源，以便让 Rancher 读取
cp k3s-rancher.desistdaydream.ltd.crt tls.crt && cp k3s-rancher.desistdaydream.ltd.key tls.key
kubectl -n cattle-system create secret tls tls-rancher-ingress \
  --cert=tls.crt \
  --key=tls.key
cp ca.crt cacerts.pem
kubectl -n cattle-system create secret generic tls-ca \
  --from-file=cacerts.pem=./cacerts.pem

# 添加 rancher 的 helm 仓库
helm repo add rancher-stable http://rancher-mirror.oss-cn-beijing.aliyuncs.com/server-charts/stable
helm repo update

# 部署指定版本的 Rancher
helm install rancher rancher-stable/rancher \
  --namespace cattle-system \
  --set hostname=k3s-rancher.desistdaydream.ltd \
  --set ingress.tls.source=secret \
  --set privateCA=true
```

## 高可用离线部署

与在线部署类似，但是需要一个私有镜像仓库，所需相关的部署镜像需要先推送到私有镜像仓库中。

需要提前准备的文件列表：

在 <https://github.com/rancher/rancher/releases> 页面，下载推送镜像所需的文件，这里以 2.4.5 为例。一共需要三个文件。

rancher-images.txt rancher 镜像列表。

rancher-save-images.sh 根据镜像列表文件打包所有镜像。

rancher-load-images.sh 将打包好的镜像推送到私有仓库。

1. rancher-images.tar.gz # rancher 的镜像
   1. curl -LO <https://github.com/rancher/rancher/releases/download/v2.4.5/rancher-images.txt>
   2. curl -LO <https://github.com/rancher/rancher/releases/download/v2.4.5/rancher-save-images.sh>
   3. sort -u rancher-images.txt -o rancher-images.txt
   4. ./rancher-save-images.sh --image-list ./rancher-images.txt
   5. 保存完成后，会生成名为 rancher-images.tar.gz 的镜像打包文件。
2. rancher-load-images.sh # 推送 rancher 镜像的脚本
   1. curl -LO <https://github.com/rancher/rancher/releases/download/v2.4.5/rancher-load-images.sh>
3. mysql.tar # mysql 镜像。
4. k3s # k3s 二进制文件
5. k3s-airgap-images-amd64.tar # k3s 运行所需的镜像
   1. 从 [此处](https://github.com/rancher/k3s/releases) 下载 k3s 二进制文件以及镜像的压缩包。二进制文件名称为 [k3s](https://github.com/rancher/k3s/releases/download/v1.18.6%2Bk3s1/k3s)。镜像压缩包的文件名为 [k3s-airgap-images-amd64.tar](https://github.com/rancher/k3s/releases/download/v1.18.6%2Bk3s1/k3s-airgap-images-amd64.tar)。
6. install.sh # 部署 k3s 的脚本
   1. 从 [此处](https://get.k3s.io) 获取离线安装所需的脚本。国内用户从 [这个页面](http://mirror.cnrancher.com/)的 k3s 目录下获取脚本，脚本名为 k3s-install.sh。
   2. curl -LO https://raw.githubusercontent.com/rancher/k3s/master/install.sh
   3. curl -LO https://docs.rancher.cn/k3s/k3s-install.sh
7. rancher-2.4.5.tgz # 用于部署 rancher 的 helm chart。
   1. helm repo add rancher-stable https://releases.rancher.com/server-charts/stable
   2. helm pull rancher-stable/rancher
8. helm # helm 二进制文件
   1. 从 git 上下载 tar 包，解压获取二进制文件。
9. kubectl # kubectl 二进制文件，用于在 Rancher 创建的集群节点上操作集群。

### 部署私有镜像仓库

略。

推送用于部署 Rancher 所需的镜像，到私有镜像仓库。

    # 拷贝 rancher-images.tar.gz 文件到当前目录
    # 假如私有镜像仓库的访问路径为 http://172.38.40.180
    ./rancher-load-images.sh --image-list ./rancher-images.txt --registry 172.38.40.180

### 启动 mysql

注意修改 ${CustomRegistry} 变量为指定的私有镜像仓库的仓库名

    docker run -d  --name k3s-mysql  --restart=always  \
    -v /opt/k3s-cluster/mysql/conf:/etc/mysql/conf.d \
    -v /opt/k3s-cluster/mysql/data:/var/lib/mysql \
    -e MYSQL_ROOT_PASSWORD=root \
    -p 3306:3306 \
    ${CustomRegistry}mysql:5.7.29 --default-time-zone=+8:00 \
    --character-set-server=utf8mb4 --collation-server=utf8mb4_general_ci \
    --explicit_defaults_for_timestamp=true --lower_case_table_names=1 --max_allowed_packet=128M

    docker run -d  --name k3s-mysql  --restart=always  \
    -v /opt/k3s-cluster/mysql/conf:/etc/mysql/conf.d \
    -v /opt/k3s-cluster/mysql/data:/var/lib/mysql \
    -e MYSQL_ROOT_PASSWORD=root \
    -p 3306:3306 ${CustomRegistry}/mysql:5.7.29 --default-time-zone=+8:00

### 部署 k3s

准备部署环境
需要 k3s、k3s-airgap-images-amd64.tar、install.sh 文件，拷贝到同一个目录中。并在需要部署 k3s 节点的设备上执行如下命令。

```shell
# 将下载好的k3s镜像包放到指定目录中，k3s 启动时直接使用该目录的镜像压缩包，加载并启动容器。
mkdir -p /var/lib/rancher/k3s/agent/images/
cp ./k3s-airgap-images-amd64.tar /var/lib/rancher/k3s/agent/images/
# 将二进制文件放在每个节点的 /usr/local/bin 目录中。
chmod 755 ./k3s && cp ./k3s /usr/local/bin/
# 拷贝 helm 二进制文件到 /usr/local/bin 目录下
chmod 755 ./helm && cp ./helm /usr/local/bin
# 离线安装脚本放在任意路径下。
mv ./k3s-install.sh ./install.sh && chmod 755 install.sh
# 准备 k3s containerd 操作私有镜像仓库的配置文件
mkdir -p /etc/rancher/k3s
cat > /etc/rancher/k3s/registries.yaml << EOF
mirrors:
  registry-test.ehualu.com:
    endpoint:
    - "http://172.38.40.180"
configs:
  "registry-test.ehualu.com":
    auth:
      username: admin
      password: Harbor12345
EOF
# 配置解析以访问私有仓库
cat >> /etc/hosts << EOF
172.38.40.180 registry-test.ehualu.com
EOF
```

开始部署 k3s

```shell
# 在 install.sh 所在目录执行部署命令
INSTALL_K3S_SKIP_DOWNLOAD=true INSTALL_K3S_EXEC='server --docker --datastore-endpoint=mysql://root:root@tcp(172.38.40.212:3306)/k3s' ./install.sh
```

Note：

1. 若 k3s.service 无法启动，报错 msg="starting kubernetes: preparing server: creating storage endpoint: building kine: dial tcp: unknown network tcp"，则修改 /etc/systemd/system/k3s.service 文件。
   1. 将其中 '--datastore-endpoint=mysql://root:root@tcp(172.38.40.214:3306)/k3s' \ 这行改为 '--datastore-endpoint=mysql://root:root@tcp(172.38.40.214:3306)/k3s'  \\

配置 kubectl config 文件

虽然高版本 k3s 在 /usr/local/bin/ 目录下生成了 kubectl 的软连接，但是 kubeconfig 文件依然需要拷贝到 .kube 目录中，因为 helm 也会使用。

```shell
# 获取 kubectl 二进制文件并放入 /usr/local/bin/ 目录中
# 从别的机器 copy 一个对应 k3s 版本的 kubeclt 二进制文件
# 拷贝 kubeconfig 文件到 kubectl 配置目录，并配置命令补全功能
mkdir ~/.kube
cp /etc/rancher/k3s/k3s.yaml /root/.kube/config
echo "source <(kubectl completion bash)" >> ~/.bashrc
```

### 部署 Rancher

```shell
# 创建所需 namespaces
kubectl create namespace cattle-system
```

创建 Rancher 所需证书

```shell
# 创建证书
openssl genrsa -out ca.key 4096

openssl req -x509 -new -nodes -sha512 -days 36500 \
  -subj "/C=CN/ST=Tianjin/L=Tianjin/O=eHualu/OU=Operations/CN=rancher.ehualu.com" \
  -key ca.key \
  -out ca.crt
openssl genrsa -out ehualu.com.key 4096

openssl req -sha512 -new \
  -subj "/C=CN/ST=Tianjin/L=Tianjin/O=eHualu/OU=Operations/CN=rancher.ehualu.com" \
  -key ehualu.com.key \
  -out ehualu.com.csr

cat > v3.ext <<-EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1=rancher.ehualu.com
EOF

openssl x509 -req -sha512 -days 36500 \
  -extfile v3.ext \
  -CA ca.crt -CAkey ca.key -CAcreateserial \
  -in ehualu.com.csr \
  -out ehualu.com.crt

# 将证书添加到 secret 资源，以便让 Rancher 读取
cp ehualu.com.crt tls.crt && cp ehualu.com.key tls.key
kubectl -n cattle-system create secret tls tls-rancher-ingress \
  --cert=tls.crt \
  --key=tls.key

cp ca.crt cacerts.pem
kubectl -n cattle-system create secret generic tls-ca \
  --from-file=cacerts.pem=./cacerts.pem
```

部署

    # 通过 helm 生成部署 rancher 的 yaml。
    helm template rancher ./rancher-2.4.5.tgz --output-dir . \
    --namespace cattle-system \
    --set hostname=rancher.ehualu.com \
    --set rancherImage=registry.ehualu.com/rancher/rancher \
    --set ingress.tls.source=secret \
    --set privateCA=true \
    --set systemDefaultRegistry=registry.ehualu.com \
    --set useBundledSystemChart=true \
    --set rancherImageTag=v2.4.5
    # 部署 Rancher
    kubectl -n cattle-system apply -R -f ./rancher

# Rancher 升级

Rancher 本身的升级，就是 k8s 集群中服务的升级，使用 helm 更新即可，新版 pod 创建后销毁旧版 pod。

Rancher 管理的 k8s 集群升级参考 [官方文档](https://docs.rancher.cn/docs/rancher2/cluster-admin/upgrading-kubernetes/_index/)，在 web 页面点两下就好很简单 。

# Rancher 清理

## 清理 Rancher

官方文档：<https://rancher.com/docs/rancher/v2.x/en/faq/removing-rancher/#what-if-i-don-t-want-rancher-anymore>

通过 [rancher 系统工具的 remove 子命令](https://rancher.com/docs/rancher/v2.x/en/system-tools/#remove)来清理 k8s 集群上 rancher

## 清理通过 Rancher 创建的 k8s 集群

官方文档：<https://docs.rancher.cn/docs/rancher2/cluster-admin/cleaning-cluster-nodes/_index/>

在 Rancher web UI 上删除集群后，手动执行一些 [命令](https://rancher2.docs.rancher.cn/docs/cluster-admin/cleaning-cluster-nodes/_index#%E6%89%8B%E5%8A%A8%E4%BB%8E%E9%9B%86%E7%BE%A4%E4%B8%AD%E5%88%A0%E9%99%A4-rancher-%E7%BB%84%E4%BB%B6) 以删除在节点上生成的数据，并重启相关节点。

```shell
#！/bin/bash
# 清理所有 Docker 容器、镜像和卷：
docker rm -f $(docker ps -qa)
docker rmi -f $(docker images -q)
docker volume rm $(docker volume ls -q)
# 卸载挂载
for mount in $(mount | grep tmpfs | grep '/var/lib/kubelet' | awk '{ print $3 }') /var/lib/kubelet /var/lib/rancher; do umount $mount; done
# 清理目录及数据
rm -rf /etc/ceph \
       /etc/cni \
       /etc/kubernetes \
       /opt/cni \
       /opt/rke \
       /run/secrets/kubernetes.io \
       /run/calico \
       /run/flannel \
       /var/lib/calico \
       /var/lib/etcd \
       /var/lib/cni \
       /var/lib/kubelet \
       /var/lib/rancher/rke/log \
       /var/log/containers \
       /var/log/kube-audit \
       /var/log/pods \
       /var/run/calico
# 清理 iptables 与 网络设备，需要重启设备，也可以自己手动清理
# reboot
```
