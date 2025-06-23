---
title: "升级为集群 - K8S训练营"
source: "https://www.qikqiak.com/k8strain/maintain/cluster/"
author:
published:
created: 2025-06-23
description:
tags:
  - "clippings"
---
[](https://github.com/cnych/qikqiak.com/edit/master/docs/maintain/cluster.md "编辑此页")

## 高可用

前面我们课程中的集群是单 master 的集群，对于生产环境风险太大了，非常有必要做一个 [高可用的集群](https://kubernetes.io/zh/docs/setup/production-environment/tools/kubeadm/ha-topology/) ，这里的高可用主要是针对控制面板来说的，比如 kube-apiserver、etcd、kube-controller-manager、kube-scheduler 这几个组件，其中 kube-controller-manager 于 kube-scheduler 组件是 Kubernetes 集群自己去实现的高可用，当有多个组件存在的时候，会自动选择一个作为 Leader 提供服务，所以不需要我们手动去实现高可用，apiserver 和 etcd 就需要手动去搭建高可用的集群的。高可用的架构有很多，比如典型的 haproxy + keepalived 架构，或者使用 nginx 来做代理实现。我们这里为了说明如何将单 master 升级为高可用的集群，采用相对更简单的 nginx 模式，当然这种模式也有一些缺点，但是足以说明高可用的实现方式了。架构如下图所示：

![](https://bxdc-static.oss-cn-beijing.aliyuncs.com/images/20200903180048.png)

从上面架构图上可以看出来，我们需要在所有的节点上安装一个 nginx 来代理 apiserver，这里我们准备3个节点作为控制平面节点：ydzs-master、ydzs-master2、ydzs-master3，这里我们默认所有节点都已经正常安装配置好了 Docker：

![](https://bxdc-static.oss-cn-beijing.aliyuncs.com/images/20200903165950.png)

在开始下面的操作之前，在 **所有节点** hosts 中配置如下所示的信息：

```shell
$ cat /etc/hosts
127.0.0.1 api.k8s.local
10.151.30.70 ydzs-master2
10.151.30.71 ydzs-master3
10.151.30.11 ydzs-master
10.151.30.57 ydzs-node3
10.151.30.59 ydzs-node4
10.151.30.60 ydzs-node5
10.151.30.62 ydzs-node6
10.151.30.22 ydzs-node1
10.151.30.23 ydzs-node2
```

## 更新证书

由于我们要将集群替换成高可用的集群，那么势必会想到我们会用一个负载均衡器来代理 APIServer，也就是这个负载均衡器访问 APIServer 的时候要能正常访问，所以默认安装的 APIServer 证书就需要更新，因为里面没有包含我们需要的地址，需要保证在 SAN 列表中包含一些额外的名称。

首页我们一个 kubeadm 的配置文件，如果一开始安装集群的时候你就是使用的配置文件，那么我们可以直接更新这个配置文件，但是如果你没有使用配置文件，直接使用的 `kubeadm init` 来安装的集群，那么我们可以从集群中获取 kubeadm 的配置信息来创建一个配置文件，因为 kubeadm 会将其配置写入到 kube-system 命名空间下面一个名为 `kubeadm-config` 的 ConfigMap 中。可以直接执行如下所示的命令将该配置导出：

```shell
$ kubectl -n kube-system get configmap kubeadm-config -o jsonpath='{.data.ClusterConfiguration}' > kubeadm.yaml
```

上面的命令会导出一个名为 kubeadm.yaml 的配置文件，内容如下所示：

```yaml
apiServer:
  extraArgs:
    authorization-mode: Node,RBAC
  timeoutForControlPlane: 4m0s
apiVersion: kubeadm.k8s.io/v1beta2
certificatesDir: /etc/kubernetes/pki
clusterName: kubernetes
controllerManager: {}
dns:
  type: CoreDNS
etcd:
  local:
    dataDir: /var/lib/etcd
imageRepository: registry.aliyuncs.com/k8sxio
kind: ClusterConfiguration
kubernetesVersion: v1.17.11
networking:
  dnsDomain: cluster.local
  podSubnet: 10.244.0.0/16
  serviceSubnet: 10.96.0.0/12
scheduler: {}
```

上面的配置中并没有列出额外的 SAN 信息，我们要添加一个新的数据，需要在 `apiServer` 属性下面添加一个 `certsSANs` 的列表。如果你在启动集群的使用就使用的了 kubeadm 的配置文件，可能里面就已经包含 certSANs 列表了，如果没有我们就需要添加它，比如我们这里要添加一个新的域名 `api.k8s.local` 以及 ydzs-master2 和 ydzs-master3 这两个主机名和 10.151.30.70、10.151.30.71 这两个新的 IP 地址，那么我们需要在 apiServer 下面添加如下所示的数据：

```yaml
apiServer:
  certSANs:
  - api.k8s.local
  - ydzs-master
  - ydzs-master2
  - ydzs-master3
  - 10.151.30.11
  - 10.151.30.70
  - 10.151.30.71
  extraArgs:
    authorization-mode: Node,RBAC
  timeoutForControlPlane: 4m0s
```

上面我只列出了 apiServer 下面新增的 certSANs 信息，这些信息是包括在标准的 SAN 列表之外的，所以不用担心这里没有添加 kubernetes、kubernetes.default 等等这些信息，因为这些都是标准的 SAN 列表中的。

更新完 kubeadm 配置文件后我们就可以更新证书了，首先我们移动现有的 APIServer 的证书和密钥，因为 kubeadm 检测到他们已经存在于指定的位置，它就不会创建新的了。

```shell
$ mv /etc/kubernetes/pki/apiserver.{crt,key} ~
```

然后直接使用 kubeadm 命令生成一个新的证书：

```shell
$ kubeadm init phase certs apiserver --config kubeadm.yaml
W0902 10:05:28.006627     832 validation.go:28] Cannot validate kubelet config - no validator is available
W0902 10:05:28.006754     832 validation.go:28] Cannot validate kube-proxy config - no validator is available
[certs] Generating "apiserver" certificate and key
[certs] apiserver serving cert is signed for DNS names [ydzs-master kubernetes kubernetes.default kubernetes.default.svc kubernetes.default.svc.cluster.local api.k8s.local ydzs-master2 ydzs-master3] and IPs [10.96.0.1 123.59.188.12 10.151.30.11 10.151.30.70 10.151.30.71]
```

通过上面的命令可以查看到 APIServer 签名的 DNS 和 IP 地址信息，一定要和自己的目标签名信息进行对比，如果缺失了数据就需要在上面的 certSANs 中补齐，重新生成证书。

该命令会使用上面指定的 kubeadm 配置文件为 APIServer 生成一个新的证书和密钥，由于指定的配置文件中包含了 certSANs 列表，那么 kubeadm 会在创建新证书的时候自动添加这些 SANs。

最后一步是重启 APIServer 来接收新的证书，最简单的方法是直接杀死 APIServer 的容器：

```shell
$ docker ps | grep kube-apiserver | grep -v pause
7fe227a5dd3c        aa63290ccd50                               "kube-apiserver --ad…"   14 hours ago        Up 14 hours                             k8s_kube-apiserver_kube-apiserver-ydzs-master_kube-system_6aa38ee2d66b7d9b6660a88700d00581_0
$ docker kill 7fe227a5dd3c
7fe227a5dd3c
```

容器被杀掉后，kubelet 会自动重启容器，然后容器将接收新的证书，一旦 APIServer 重启后，我们就可以使用新添加的 IP 地址或者主机名来连接它了，比如我们新添加的 `api.k8s.local` 。

## 验证证书

要验证证书是否更新我们可以直接去编辑 kubeconfig 文件中的 APIServer 地址，将其更换为新添加的 IP 地址或者主机名，然后去使用 kubectl 操作集群，查看是否可以正常工作。

当然我们可以使用 *openssl* 命令去查看生成的证书信息是否包含我们新添加的 SAN 列表数据：

```shell
$ openssl x509 -in /etc/kubernetes/pki/apiserver.crt -text
Certificate:
            ......
        Subject: CN=kube-apiserver
        ......
            X509v3 Subject Alternative Name:
                DNS:ydzs-master, DNS:kubernetes, DNS:kubernetes.default, DNS:kubernetes.default.svc, DNS:kubernetes.default.svc.cluster.local, DNS:api.k8s.local, DNS:ydzs-master2, DNS:ydzs-master3, IP Address:10.96.0.1, IP Address:123.59.188.12, IP Address:10.151.30.11, IP Address:10.151.30.70, IP Address:10.151.30.71
......
```

如果上面的操作都一切顺利，最后一步是将上面的集群配置信息保存到集群的 kubeadm-config 这个 ConfigMap 中去，这一点非常重要，这样以后当我们使用 kubeadm 来操作集群的时候，相关的数据不会丢失，比如升级的时候还是会带上 certSANs 中的数据进行签名的。

```shell
$ kubeadm config upload from-file --config kubeadm.yaml
```

使用上面的命令保存配置后，我们同样可以用下面的命令来验证是否保存成功了：

```shell
$ kubectl -n kube-system get configmap kubeadm-config -o yaml
```

更新 APIServer 证书的名称在很多场景下都会使用到，比如在控制平面前面添加一个负载均衡器，或者添加新的 DNS 名称或 IP 地址来使用控制平面的端点，所以掌握更新集群证书的方法也是非常有必要的。

## 为控制平面创建负载均衡器

接下来我们为控制平面创建一个负载平衡器。如何设置和配置负载均衡器的具体细节因解决方案不同，但是一般的方案都需要包括下面的功能：

- 使用4层负载平衡器（TCP而不是HTTP / HTTPS）
- 运行健康检查应配置为 SSL，而不是 TCP 运行状况检查

不管用哪一种方式，我们最好创建一个 DNS CNAME 条目以指向您的负载均衡器（强烈建议）。如果您需要更换或重新配置负载均衡解决方案，这将为您提供更多的灵活性，因为 DNS CNAME 保持不变，就不用再次去更新证书了。

我们这里采用的方案是在节点上使用 nginx 来作为一个负载均衡器，下面的操作需要在 **所有节点** 上操作：

```shell
$ mkdir -p /etc/kubernetes

$ cat > /etc/kubernetes/nginx.conf << EOF
error_log stderr notice;

worker_processes 2;
worker_rlimit_nofile 130048;
worker_shutdown_timeout 10s;

events {
  multi_accept on;
  use epoll;
  worker_connections 16384;
}

stream {
  upstream kube_apiserver {
    least_conn;
    server ydzs-master:6443;
    server ydzs-master2:6443;
    server ydzs-master3:6443;
  }

  server {
    listen        8443;
    proxy_pass    kube_apiserver;
    proxy_timeout 10m;
    proxy_connect_timeout 1s;
  }
}

http {
  aio threads;
  aio_write on;
  tcp_nopush on;
  tcp_nodelay on;

  keepalive_timeout 5m;
  keepalive_requests 100;
  reset_timedout_connection on;
  server_tokens off;
  autoindex off;

  server {
    listen 8081;
    location /healthz {
      access_log off;
      return 200;
    }
    location /stub_status {
      stub_status on;
      access_log off;
    }
  }
}
EOF
```

使用上面的配置启动一个 nginx 容器：

```shell
$ docker run --restart=always \
    -v /etc/kubernetes/nginx.conf:/etc/nginx/nginx.conf \
    -v /etc/localtime:/etc/localtime:ro \
    --name k8s-ha \
    --net host \
    -d \
    nginx
```

启动成功后 apiserver 的负载均衡地址就成了 `https://api.k8s.local:8443` 。然后我们将 kubeconfig 文件中的 apiserver 地址替换成负载均衡器的地址。

```shell
# 修改 kubelet 配置
$ vi /etc/kubernetes/kubelet.conf
......
    server: https://api.k8s.local:8443
  name: kubernetes
......
$ systemctl restart kubelet
# 修改 controller-manager
$ vi /etc/kubernetes/controller-manager.conf
......
    server: https://api.k8s.local:8443
  name: kubernetes
......
# 重启
$ docker kill $(docker ps | grep kube-controller-manager | \
grep -v pause | cut -d' ' -f1)
# 修改 scheduler
$ vi /etc/kubernetes/scheduler.conf
......
    server: https://api.k8s.local:8443
  name: kubernetes
......
# 重启
$ docker kill $(docker ps | grep kube-scheduler | grep -v pause | \
cut -d' ' -f1)
```

然后更新 kube-proxy

```shell
$ kubectl -n kube-system edit cm kube-proxy
......
  kubeconfig.conf: |-
    apiVersion: v1
    kind: Config
    clusters:
    - cluster:
        certificate-authority: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        server: https://api.k8s.local:8443
      name: default
......
```

当然还有 kubectl 访问集群的 `~/.kube/config` 文件也需要修改。

## 更新控制面板

由于我们现在已经在控制平面的前面添加了一个负载平衡器，因此我们需要使用正确的信息更新此 ConfigMap。（您很快就会将控制平面节点添加到集群中，因此在此ConfigMap中拥有正确的信息很重要。）

首先，使用以下命令从 ConfigMap 中获取当前配置：

```shell
$ kubectl -n kube-system get configmap kubeadm-config -o jsonpath='{.data.ClusterConfiguration}' > kubeadm.yaml
```

然后在当前配置文件里面里面添加 `controlPlaneEndpoint` 属性，用于指定控制面板的负载均衡器的地址。

```shell
$ vi kubeadm.yaml
controlPlaneEndpoint: api.k8s.local:8443  # 添加改配置
apiServer:
  certSANs:
  - api.k8s.local
  - ydzs-master
  - ydzs-master2
  - ydzs-master3
  - 10.151.30.11
  - 10.151.30.70
  - 10.151.30.71
......
```

编辑完文件后，使用以下命令将其上传回集群：

```shell
$ kubeadm config upload from-file --config kubeadm.yaml
```

然后需要在 `kube-public` 命名空间中更新 `cluster-info` 这个 ConfigMap，该命名空间包含一个Kubeconfig 文件，该文件的 `server:` 一行指向单个控制平面节点。只需使用 `kubectl -n kube-public edit cm cluster-info` 更新该 `server:` 行以指向控制平面的负载均衡器即可。

```shell
$ kubectl -n kube-public edit cm cluster-info
......
    server: https://api.k8s.local:8443
  name: ""
......
$  kubectl cluster-info
Kubernetes master is running at https://api.k8s.local:8443
KubeDNS is running at https://api.k8s.local:8443/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
KubeDNSUpstream is running at https://api.k8s.local:8443/api/v1/namespaces/kube-system/services/kube-dns-upstream:dns/proxy
Metrics-server is running at https://api.k8s.local:8443/api/v1/namespaces/kube-system/services/https:metrics-server:/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

更新完成就可以看到 cluster-info 的信息变成了负载均衡器的地址了。

## 添加控制平面

接下来我们来添加额外的控制平面节点，首先使用如下命令来将集群的证书上传到集群中，供其他控制节点使用：

```shell
$ kubeadm init phase upload-certs --upload-certs
I0903 15:13:24.192467   20533 version.go:251] remote version is much newer: v1.19.0; falling back to: stable-1.17
W0903 15:13:25.739892   20533 validation.go:28] Cannot validate kube-proxy config - no validator is available
W0903 15:13:25.739966   20533 validation.go:28] Cannot validate kubelet config - no validator is available
[upload-certs] Storing the certificates in Secret "kubeadm-certs" in the "kube-system" Namespace
[upload-certs] Using certificate key:
e71ef7ede98e49f5f094b150d604c7ad50f125279180a7320b1b14ef3ccc3a34
```

上面的命令会生成一个新的证书密钥，但是只有2小时有效期。由于我们现有的集群已经运行一段时间了，所以之前的启动 Token 也已经失效了（Token 的默认生存期为24小时），所以我们也需要创建一个新的 Token 来添加新的控制平面节点：

```shell
$ kubeadm token create --print-join-command --config kubeadm.yaml
W0903 15:29:10.958329   25049 validation.go:28] Cannot validate kube-proxy config - no validator is available
W0903 15:29:10.958457   25049 validation.go:28] Cannot validate kubelet config - no validator is available
kubeadm join api.k8s.local:8443 --token f27w7m.adelvl3waw9kqdhp     --discovery-token-ca-cert-hash sha256:6917cbf7b0e73ecfef77217e9a27e76ef9270aa379c34af30201abd0f1088c34
```

上面的命令最后给出的提示是添加 node 节点的命令，我们这里要添加控制平面节点就要使用如下所示的命令：

```shell
$ kubeadm join <DNS CNAME of load balancer>:<lb port> \
--token <bootstrap-token> \
--discovery-token-ca-cert-hash sha256:<CA certificate hash> \
--control-plane --certificate-key <certificate-key>
```

获得了上面的添加命令过后，登录到 ydzs-master2 节点进行相关的操作，在 ydzs-master2 节点上安装软件：

```shell
$ yum install -y kubeadm-1.17.11-0 kubelet-1.17.11-0 kubectl-1.17.11-0
```

要加入控制平面，我们可以先拉取相关镜像：

```shell
$ kubeadm config images pull --image-repository registry.aliyuncs.com/k8sxio
```

然后执行上面生成的 join 命令，将参数替换后如下所示：

```shell
$ kubeadm join api.k8s.local:8443 \
--token f27w7m.adelvl3waw9kqdhp \
--discovery-token-ca-cert-hash sha256:6917cbf7b0e73ecfef77217e9a27e76ef9270aa379c34af30201abd0f1088c34 \
--control-plane --certificate-key e71ef7ede98e49f5f094b150d604c7ad50f125279180a7320b1b14ef3ccc3a34
[preflight] Running pre-flight checks
[preflight] Reading configuration from the cluster...
[preflight] FYI: You can look at this config file with 'kubectl -n kube-system get cm kubeadm-config -oyaml'
[preflight] Running pre-flight checks before initializing the new control plane instance
[preflight] Pulling images required for setting up a Kubernetes cluster
[preflight] This might take a minute or two, depending on the speed of your internet connection
[preflight] You can also perform this action in beforehand using 'kubeadm config images pull'
[download-certs] Downloading the certificates in Secret "kubeadm-certs" in the "kube-system" Namespace
[certs] Using certificateDir folder "/etc/kubernetes/pki"
[certs] Generating "apiserver-kubelet-client" certificate and key
[certs] Generating "apiserver" certificate and key
[certs] apiserver serving cert is signed for DNS names [ydzs-master2 kubernetes kubernetes.default kubernetes.default.svc kubernetes.default.svc.cluster.local api.k8s.local api.k8s.local ydzs-master2 ydzs-master3] and IPs [10.96.0.1 10.151.30.70 10.151.30.11 10.151.30.70 10.151.30.71]
[certs] Generating "front-proxy-client" certificate and key
[certs] Generating "etcd/healthcheck-client" certificate and key
[certs] Generating "etcd/server" certificate and key
[certs] etcd/server serving cert is signed for DNS names [ydzs-master2 localhost] and IPs [10.151.30.70 127.0.0.1 ::1]
[certs] Generating "etcd/peer" certificate and key
[certs] etcd/peer serving cert is signed for DNS names [ydzs-master2 localhost] and IPs [10.151.30.70 127.0.0.1 ::1]
[certs] Generating "apiserver-etcd-client" certificate and key
[certs] Valid certificates and keys now exist in "/etc/kubernetes/pki"
[certs] Using the existing "sa" key
[kubeconfig] Generating kubeconfig files
[kubeconfig] Using kubeconfig folder "/etc/kubernetes"
[endpoint] WARNING: port specified in controlPlaneEndpoint overrides bindPort in the controlplane address
[kubeconfig] Writing "admin.conf" kubeconfig file
[kubeconfig] Writing "controller-manager.conf" kubeconfig file
[kubeconfig] Writing "scheduler.conf" kubeconfig file
[control-plane] Using manifest folder "/etc/kubernetes/manifests"
[control-plane] Creating static Pod manifest for "kube-apiserver"
W0903 15:55:08.444989    4353 manifests.go:214] the default kube-apiserver authorization-mode is "Node,RBAC"; using "Node,RBAC"
[control-plane] Creating static Pod manifest for "kube-controller-manager"
W0903 15:55:08.457787    4353 manifests.go:214] the default kube-apiserver authorization-mode is "Node,RBAC"; using "Node,RBAC"
[control-plane] Creating static Pod manifest for "kube-scheduler"
W0903 15:55:08.459829    4353 manifests.go:214] the default kube-apiserver authorization-mode is "Node,RBAC"; using "Node,RBAC"
[check-etcd] Checking that the etcd cluster is healthy
......
This node has joined the cluster and a new control plane instance was created:

* Certificate signing request was sent to apiserver and approval was received.
* The Kubelet was informed of the new secure connection details.
* Control plane (master) label and taint were applied to the new node.
* The Kubernetes control plane instances scaled up.
* A new etcd member was added to the local/stacked etcd cluster.

To start administering your cluster from this node, you need to run the following as a regular user:

    mkdir -p $HOME/.kube
    sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
    sudo chown $(id -u):$(id -g) $HOME/.kube/config

Run 'kubectl get nodes' to see this node join the cluster.
```

如果在 etcd 里面有残留的废弃节点数据，可以用如下命令删除：

```shell
# 列出所有成员
$ ./etcdctl --endpoints=$ENDPOINTS member list  
# 删除指定成员
$ ./etcdctl --endpoints=$ENDPOINTS member remove b9057cfdc8ff17ce
```

到这里可以看到 ydzs-master2 节点就成功加入到了控制平面中，然后根据上面的提示配置 kubeconfig 文件。然后用同样的方式添加 ydzs-master3 节点，都添加成功后，在 ydzs-master3 节点上执行如下所示的命令来验证 etcd 集群使用正常：

```shell
$ docker run --rm -it \
--net host \
-v /etc/kubernetes:/etc/kubernetes registry.aliyuncs.com/k8sxio/etcd:3.4.3-0 etcdctl \
--cert /etc/kubernetes/pki/etcd/peer.crt \
--key /etc/kubernetes/pki/etcd/peer.key \
--cacert /etc/kubernetes/pki/etcd/ca.crt \
--endpoints https://10.151.30.71:2379 endpoint health --cluster
# endpoint status --write-out=table 查看状态
https://10.151.30.70:2379 is healthy: successfully committed proposal: took = 19.410192ms
https://10.151.30.71:2379 is healthy: successfully committed proposal: took = 21.077275ms
https://10.151.30.11:2379 is healthy: successfully committed proposal: took = 31.282643ms
$ docker run --rm -it --net host -v /etc/kubernetes:/etc/kubernetes registry.aliyuncs.com/k8sxio/etcd:3.4.3-0 etcdctl --cert /etc/kubernetes/pki/etcd/peer.crt --key /etc/kubernetes/pki/etcd/peer.key --cacert /etc/kubernetes/pki/etcd/ca.crt --endpoints https://10.151.30.11:2379,https://10.151.30.70:2379,https://10.151.30.71:2379 endpoint status --write-out=table
+---------------------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------+
|         ENDPOINT          |        ID        | VERSION | DB SIZE | IS LEADER | IS LEARNER | RAFT TERM | RAFT INDEX | RAFT APPLIED INDEX | ERRORS |
+---------------------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------+
| https://10.151.30.11:2379 |  3d1fd8983aed809 |   3.4.3 |   52 MB |     false |      false |       113 |  119502257 |          119502257 |        |
| https://10.151.30.70:2379 | af2e11ae8aa72bde |   3.4.3 |   52 MB |      true |      false |       113 |  119502259 |          119502259 |        |
| https://10.151.30.71:2379 | e7a7b252880befdf |   3.4.3 |   52 MB |     false |      false |       113 |  119502259 |          119502259 |        |
+---------------------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------+
```

正常我们就可以看到 etcd 集群正常了，但是由于控制平台的3个节点是先后安装的，所以前面两个节点的 etcd 中并不包含其他 etcd 节点的信息，所以我们需要同步所有控制平面节点的 etcd 集群配置：

```shell
$ cat /etc/kubernetes/manifests/etcd.yaml
......
- --initial-cluster=ydzs-master=https://10.151.30.11:2380,ydzs-master2=https://10.151.30.70:2380,ydzs-master3=https://10.151.30.71:2380
......
```

最后执行如下所示的命令查看集群是否正常：

```shell
$ kubectl get nodes
NAME           STATUS   ROLES    AGE    VERSION
ydzs-master    Ready    master   299d   v1.17.11
ydzs-master2   Ready    master   34m    v1.17.11
ydzs-master3   Ready    master   10m    v1.17.11
ydzs-node1     Ready    <none>   299d   v1.17.11
ydzs-node2     Ready    <none>   299d   v1.17.11
ydzs-node3     Ready    <none>   297d   v1.17.11
ydzs-node4     Ready    <none>   297d   v1.17.11
ydzs-node5     Ready    <none>   225d   v1.17.11
ydzs-node6     Ready    <none>   225d   v1.17.11
```

这里我们就可以看到 ydzs-master、ydzs-master2、ydzs-master3 3个节点变成了 master 节点，我们也就完成了将单 master 升级为多 master 的高可用集群了。