---
title: 创建高可用 Kubernetes 集群
---

# 概述

> 参考：
> - [官方文档，入门-生产环境-使用部署工具安装 Kubernetes-使用 kubeadm 引导集群-使用 kubeadm 创建高可用集群](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/high-availability/)
>   - <https://kubernetes.io/docs/setup/independent/high-availability/>
> - [Kubeadm GitHub 项目,高科用性注意事项](https://github.com/kubernetes/kubeadm/blob/master/docs/ha-considerations.md)
>   - [keepalived 和 haproxy](https://github.com/kubernetes/kubeadm/blob/master/docs/ha-considerations.md#keepalived-and-haproxy)
>   - [kube-vip](https://github.com/kubernetes/kubeadm/blob/master/docs/ha-considerations.md#kube-vip)
> - [kube-vip 官网](https://kube-vip.io/)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kc2p8h/1616120692032-ff2748b8-3aa5-4f51-a71c-3c00aa017184.png)
想要让 Kubernets 的 Master 可以高可用，本质上就是让 API Server 高可用，也就是为 API Server 创建负载均衡，通过 VIP 来访问 API Server。

有多种方式可以实现 Kubernetes Master 节点的高可用

- 在各个 Master 上负载均衡各个 Master
  - Keepalived
  - kube-vip
- 在每个 Node 上负载均衡各个 Master
  - sealos

## 创建负载均衡

### Keepalived 方式

官方推荐版本的 keepalived 的配置

    /etc/keepalived/keepalived.conf在所有主节点上创建以下配置文件：
    ! Configuration File for keepalived
    global_defs {
      router_id k8s-master-dr
      script_user root
      enable_script_security
    }

    vrrp_script check_apiserver {
      script "/etc/keepalived/check_apiserver.sh"
      interval 3
      weight -2
      fall 10
      rise 2
    }

    vrrp_instance VI_1 {
      state ${STATE}
      interface ${INTERFACE}
      virtual_router_id ${ROUTER_ID}
      priority ${PRIORITY}
      authentication {
        auth_type PASS
        auth_pass 4be37dc3b4c90194d1600c483e10ad1d
      }
      virtual_ipaddress {
        ${VIP}
      }
      track_script {
        check_apiserver
      }
    }

- 在该部分中 vrrp_instance VI_1，根据自身情况更改：
  - STATE # 改为是 MASTER（第一主节点上）或 BACKUP（另一主节点）。
  - INTERFACE # 改为将虚拟 IP 绑定到的现有公共接口的名称（通常是主接口）。
  - PRIORITY # 对于第一主节点应该更高，例如 101，对于其他主节点应该更低，例如 100、99。
  - PASS # 改为任何随机字符串
  - VIRTUAL-IP # 改为虚拟 IP。
- 下面是配套的健康检查脚本的配置

```bash
cat > /etc/keepalived/check_apiserver.sh << EOF
#!/bin/sh
errorExit() {
   echo "*** $*" 1>&2
   exit 1
}
curl --silent --max-time 2 --insecure https://localhost:6443/healthz -o /dev/null || errorExit "Error GET https://localhost:6443/healthz"
if ip addr | grep -q ${VIP}; then
   curl --silent --max-time 2 --insecure https://${VIP}:6443/healthz -o /dev/null || errorExit "Error GET https://${VIP}:6443/healthz"
fi
EOF
chmod +x check_apiserver.sh
```

- 替换其中的内容
  - ${VIP} # 您选择的虚拟 IP。
- 安装 keepalived(e.g.下面配置文件中 10.15 为 VIP)，并指定 RS 为 3 台 master 的 IP:PORT
  - 这样，在 Node 节点上报的自身状态给 api-server 时，可以直接使用 3 台 master 的 VIP 作为通信 IP

### kube-vip 方式

配置环境变量

```bash
export INTERFACE="eth0"
export VIP="172.19.42.230"
```

- $INTERFACE # 要生成 VIP 的网络设备
- $VIP # 要生成的 VIP

在每个 Master 节点生成静态 Pod 的 Manifest 文件

```bash
nerdctl run --rm --net host docker.io/plndr/kube-vip:v0.4.1 \
manifest pod \
--interface $INTERFACE \
--address $VIP \
--controlplane \
--services \
--arp \
--leaderElection | tee  /etc/kubernetes/manifests/kube-vip.yaml
```

## 在所有 Node 创建负载均衡

### Sealos

> 参考：
> - [GitHub 项目，fanux/sealos](https://github.com/fanux/sealos)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/kc2p8h/1659603459394-80959a05-8f38-490d-818e-8e3f6d4070cf.png)
与方法一不同，在 node 与 api-server 通信时，不再依靠 master 上的 VIP，而是在各个 Node 本地创建一个 lvs 规则，该规则使用一个 VIP，RS 为 3 个 master 的 IP，这样，每个 Node 在上报信息时，可以使用本地的 ipvs 规则，通过 VIP 来选择与哪个 mater 通信进行信息上报

后续步骤与方法一一样

详见：<https://github.com/fanux/sealos/blob/master/README_zh.md>

# 使用 kubeadm 初始化 k8s 的 master 节点 HA 模式

当负载均衡创建完成呢后，即可使用 kubeadm 开始初始化 Kubernetes 集群

```yaml
cat > kubeadm-config.yaml <<EOF
apiVersion: kubeadm.k8s.io/v1beta2
kind: InitConfiguration
bootstrapTokens:
- groups:
  - system:bootstrappers:kubeadm:default-node-token
  ttl: 0s
  usages:
  - signing
  - authentication
---
apiVersion: kubeadm.k8s.io/v1beta2
kind: ClusterConfiguration
kubernetesVersion: 1.18.8
controlPlaneEndpoint: 172.38.40.215:6443
imageRepository: registry.aliyuncs.com/k8sxio
networking:
  podSubnet: 10.244.0.0/16
---
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
mode: "ipvs"
---
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
cgroupDriver: systemd
EOF

#初始化第一个K8s的master设备，并使用flag中的--upload-certs功能来上传证书，以便后续新加入的master节点可以直接使用这些证书，而不用手动拷贝了
kubeadm init --config=kubeadm-config.yaml --upload-certs
注意：
#kubernetesVersion应该设置为Kubernetes版本使用。这个例子使用1.16.2。
#controlPlaneEndpoint 应匹配负载均衡器的地址或DNS和端口。即VIP

#在添加后续HA的master节点时，注意使用几个关键的flag
kubeadm join 172.38.40.215:6443 --token XXXX --discovery-token-ca-cert-hash sha256:XXXX --experimental-control-plane --certificate-key XXXX
```

- 应用适用于 HA 模式的网络插件
- 复制本机的相关证书至其余节点
- 在其余 master 节点上执行之前节点 init 后生成的 kubeadm join 命令，注意：命令末尾加上--experimental-control-plane 这个 flag
- 注意事项：
  - 如果无法 init，则开启 ipvsadm 服务，如果无法加入其余 master，则关闭 ipvsadm 服务
