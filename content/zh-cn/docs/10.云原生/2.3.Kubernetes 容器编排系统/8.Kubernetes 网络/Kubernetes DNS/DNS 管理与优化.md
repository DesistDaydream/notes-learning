---
title: DNS 管理与优化
---

优化方法

禁用 ipv6 解析，以提高解析速度

参考文章：[https://yuerblog.cc/2019/09/13/k8s-coredns%E7%A6%81%E7%94%A8ipv6%E8%A7%A3%E6%9E%90/](https://yuerblog.cc/2019/09/13/k8s-coredns%25E7%25A6%2581%25E7%2594%25A8ipv6%25E8%25A7%25A3%25E6%259E%2590/)

如果 K8S 集群宿主机没有关闭 IPV6 内核模块的话，容器请求 coredns 时的默认行为是同时发起 IPV4 和 IPV6 解析。

由于我们通常只使用 IPV4 地址，所以此时如果我们仅仅在 coredns 中配置 DOMAIN -> IPV4 地址的解析的话，当 coredns 收到 IPV6 解析请求的时候就会因为本地找不到配置而 foward 到 upstream DNS 服务器解析，从而导致容器的 DNS 解析请求变慢。

coredns 提供了一种 plugin 叫做 template，经过配置后可以给所有的 IPV6 请求立即返回一个空结果的应答，避免请求 forward 到上游 DNS。

使用方法

template 插件的官方文档地址：<https://github.com/coredns/coredns/tree/master/plugin/template>，coredns 默认已携带此插件，大家只需要配置即可。

```json
template ANY AAAA {
  rcode NXDOMAIN
}
```

该配置添加到 forward 下面即可

AAAA 表示 IPV6 解析请求，rcode 控制应答返回 NXDOMAIN，即表示没有解析结果。

修改每个 pod 中 /etc/resolv.conf 中 ndots 的值

由于每个 pod 启动后，其 resolv.conf 文件中，会自动加入 options: ndots=5 的参数，所以此 pod 想要通过域名与其他服务交互，很有可能会搜索 search ，产生很多不必要的解析请求。所以修改 ndots 可以省掉很多解析时间。

```yaml
apiVersion: v1
kind: Pod
metadata:
  namespace: default
  name: dns-example
spec:
  containers:
  - name: test
    image: nginx
  dnsConfig:
    options:
    - name: ndots
        value: "1"
```

尝试修改 coredns 的配置来实现，如下（需要验证）：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qb8zlf/1616118141230-181857cf-18fc-4ccf-a1a9-d73735230cc7.png)

附：一则 Corefile 的 configmap 的配置：

```json
# kubectl  get configmap -n kube-system coredns -o jsonpath='{.data}'
map[Corefile:.:53 {
    log
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
       pods insecure
       upstream
       fallthrough in-addr.arpa ip6.arpa
    }
    hosts /etc/hostfile {
       fallthrough
    }
    proxy domain.com 10.111.xxx.1:53 10.111.xxx.2:53 {
       policy round_robin
    }
    proxy domain.cn 10.111.xxx.1:53 10.111.xxx.2:53 {
       policy round_robin
    }
    proxy consul 10.111.xxx.1:53 10.111.xxx.2:53 {
       policy round_robin
    }
    prometheus :9153
    proxy . /etc/resolv.conf {
       policy first
    }
    cache 120
    loop
    reload  //加上这个配置之后修改配置会自动reload pod的配置，可能存在pod不会重启
    loadbalance
}
]
```

使用：

    # kubectl edit deployment -n kube-system coredns

编辑修改，然后修改对应的 deployment 重启 pod 或者重新发布 pod 生效。

    error: 错误记录到stdout
    health ：CoreDNS的运行状况报告为 http：// localhost：8080 / health
    kubernetes ：CoreDNS将根据Kubernetes服务和pod的IP回复DNS查询
    prometheus ：CoreDNS的度量标准可以在 http://localhost:9153/ Prometheus格式的 指标 中找到
    proxy ：任何不在Kubernetes集群域内的查询都将转发到预定义的解析器（/etc/resolv.conf）
    cache ：启用前端缓存
    loop ：检测简单的转发循环，如果找到循环则停止CoreDNS进程
    reload ：允许自动重新加载已更改的Corefile。编辑ConfigMap配置后，请等待两分钟以使更改生效
    loadbalance ：这是一个循环DNS负载均衡器，可以在答案中随机化A，AAAA和MX记录的顺序

使用 .spec.HostAliases 字段将解析条目添加到容器内的 /etc /hosts 文件中

官方文档：<https://kubernetes.io/docs/concepts/services-networking/add-entries-to-pod-etc-hosts-with-host-aliases/>

为什么通过 Kubelet 管理 hosts 文件？

kubelet 管理 Pod 中每个容器的 hosts 文件，避免 Docker 在容器已经启动之后去 修改 该文件。

因为该文件是托管性质的文件，无论容器重启或 Pod 重新调度，用户修改该 hosts 文件的任何内容，都会在 Kubelet 重新安装后被覆盖。因此，不建议修改该文件的内容。

示例 manifests

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: hostaliases-pod
spec:
  restartPolicy: Never
  hostAliases:
    - ip: "127.0.0.1"
      hostnames:
        - "foo.local"
        - "bar.local"
    - ip: "10.1.2.3"
      hostnames:
        - "foo.remote"
        - "bar.remote"
  containers:
    - name: cat-hosts
      image: busybox
      command:
        - cat
      args:
        - "/etc/hosts"
```

hostaliases-pod 这个 pod 的容器内的 /etc/hosts 文件的内容看起来类似如下这样

```bash
# Kubernetes-managed hosts file.
127.0.0.1	localhost
::1	localhost ip6-localhost ip6-loopback
fe00::0	ip6-localnet
fe00::0	ip6-mcastprefix
fe00::1	ip6-allnodes
fe00::2	ip6-allrouters
10.200.0.5	hostaliases-pod
# Entries added by HostAliases.
127.0.0.1	foo.local	bar.local
10.1.2.3	foo.remote	bar.remote
```
