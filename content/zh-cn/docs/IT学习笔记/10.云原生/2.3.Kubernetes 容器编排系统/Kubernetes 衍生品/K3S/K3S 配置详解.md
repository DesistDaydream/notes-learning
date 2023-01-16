---
title: K3S 配置详解
---

# 概述

/etc/rancher/k3s/registries.yaml

k3s 通过 containerd 来控制容器，在 pull 镜像时，会默认指定 docker.io 为 registry 且不可改。

该配置用于将 docker.io 这个域名解析到指定的 私有镜像仓库地址，这样在 crictl pull IMAGE 时，会去私有镜像仓库拉取镜像。这其中需要指定登录私有镜像仓库的用户名和密码。

mirrors: docker.io: endpoint: - "http://172.38.40.180"configs: "172.38.40.180": auth: username: admin password: Harbor12345

可以将 docker.io 修改为自己设定的域名，只不过这样需要在 pull 的时候，选择镜像的时候加上这个域名。

# k3s server 命令行参数详解

-v value (logging) Number for the log level verbosity (default: 0)
\--vmodule value (logging) Comma-separated list of pattern=N settings for file-filtered logging
\--log value, -l value (logging) Log to file
\--alsologtostderr (logging) Log to standard error as well as file (if set)
\--bind-address value (listener) k3s bind address (default: 0.0.0.0)
\--https-listen-port value (listener) HTTPS listen port (default: 6443)
\--advertise-address value (listener) IP address that apiserver uses to advertise to members of the cluster (default: node-external-ip/node-ip)
\--advertise-port value (listener) Port that apiserver uses to advertise to members of the cluster (default: listen-port) (default: 0)
\--tls-san value (listener) Add additional hostname or IP as a Subject Alternative Name in the TLS cert
\--data-dir value, -d value (data) Folder to hold state default /var/lib/rancher/k3s or ${HOME}/.rancher/k3s if not root
\--cluster-cidr value                       (networking) Network CIDR to use for pod IPs (default: "10.42.0.0/16")
\--service-cidr value                       (networking) Network CIDR to use for services IPs (default: "10.43.0.0/16")
\--cluster-dns value                        (networking) Cluster IP for coredns service. Should be in your service-cidr range (default: 10.43.0.10)
\--cluster-domain value                     (networking) Cluster Domain (default: "cluster.local")
\--flannel-backend value                    (networking) One of 'none', 'vxlan', 'ipsec', or 'flannel' (default: "vxlan")
\--token value, -t value                    (cluster) Shared secret used to join a server or agent to a cluster \[$K3S_TOKEN]
\--token-file value (cluster) File containing the cluster-secret/token \[$K3S\_TOKEN\_FILE]
\--write-kubeconfig value, -o value         (client) Write kubeconfig for admin client to this file \[$K3S_KUBECONFIG_OUTPUT]
\--write-kubeconfig-mode value (client) Write kubeconfig with this mode \[$K3S_KUBECONFIG_MODE]
**--kubelet-arg value **# 指定 kubelet 的运行时 flags。

- 比如 k3s server --docker --kubelet-arg cgroup-driver=systemd
- Note：在指定 kubelet 参数时，不用加 -- ，k3s 会自动添加。如果加了-- ，那么就会变成 ----cgroup-driver=systemd。这样 kubelet 是无法启动的。
- 上面仅仅举个例子，k3s 内嵌的 kubelet 不支持 systemd 类型的 cgroup-driver

**--kube-apiserver-arg VALUE **# (flags)自定义 kube-apiserver 进程的命令行 Flags。
**--kube-controller-manager-arg value** # (flags)自定义 kube-controller-manager 进程的命令行 Flags。
**--kube-scheduler-arg value **# (flags)自定义 kube-scheduler 进程的命令行 Flags。
\--kube-proxy-arg value (agent/flags) Customized flag for kube-proxy process
\--kube-cloud-controller-manager-arg value (flags) Customized flag for kube-cloud-controller-manager process
\--datastore-endpoint value (db) Specify etcd, Mysql, Postgres, or Sqlite (default) data source name \[$K3S\_DATASTORE\_ENDPOINT]
\--datastore-cafile value                   (db) TLS Certificate Authority file used to secure datastore backend communication \[$K3S_DATASTORE_CAFILE]
\--datastore-certfile value (db) TLS certification file used to secure datastore backend communication \[$K3S\_DATASTORE\_CERTFILE]
\--datastore-keyfile value                  (db) TLS key file used to secure datastore backend communication \[$K3S_DATASTORE_KEYFILE]
\--default-local-storage-path value (storage) Default local storage path for local provisioner storage class
\--no-deploy value (components) Do not deploy packaged components (valid items: coredns, servicelb, traefik, local-storage, metrics-server)
\--disable-scheduler (components) Disable Kubernetes default scheduler
\--disable-cloud-controller (components) Disable k3s default cloud controller manager
\--disable-network-policy (components) Disable k3s default network policy controller
\--node-name value (agent/node) Node name \[$K3S\_NODE\_NAME]
\--with-node-id                             (agent/node) Append id to node name
\--node-label value                         (agent/node) Registering kubelet with set of labels
\--node-taint value                         (agent/node) Registering kubelet with set of taints
**--docker** # 指定 k3s 的 CRI 为 docker。默认为 containerd。
\--container-runtime-endpoint value         (agent/runtime) Disable embedded containerd and use alternative CRI implementation
\--pause-image value                        (agent/runtime) Customized pause image for containerd sandbox
\--private-registry value                   (agent/runtime) Private registry configuration file (default: "/etc/rancher/k3s/registries.yaml")
\--node-ip value, -i value                  (agent/networking) IP address to advertise for node
\--node-external-ip value                   (agent/networking) External IP address to advertise for node
\--resolv-conf value                        (agent/networking) Kubelet resolv.conf file \[$K3S_RESOLV_CONF]
\--flannel-iface value (agent/networking) Override default flannel interface
\--flannel-conf value (agent/networking) Override default flannel config file
\--rootless (experimental) Run rootless
\--agent-token value (experimental/cluster) Shared secret used to join agents to the cluster, but not servers \[$K3S\_AGENT\_TOKEN]
\--agent-token-file value                   (experimental/cluster) File containing the agent secret \[$K3S_AGENT_TOKEN_FILE]
\--server value, -s value (experimental/cluster) Server to connect to, used to join a cluster \[$K3S\_URL]
\--cluster-init                             (experimental/cluster) Initialize new cluster master \[$K3S_CLUSTER_INIT]
\--cluster-reset (experimental/cluster) Forget all peers and become a single cluster new cluster master \[$K3S\_CLUSTER\_RESET]
\--no-flannel                               (deprecated) use --flannel-backend=none
\--cluster-secret value                     (deprecated) use --token \[$K3S_CLUSTER_SECRET]
