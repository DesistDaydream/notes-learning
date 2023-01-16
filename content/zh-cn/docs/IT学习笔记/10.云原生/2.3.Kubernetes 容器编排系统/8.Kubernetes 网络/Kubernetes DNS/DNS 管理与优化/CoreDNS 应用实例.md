---
title: CoreDNS 应用实例
---

原文链接：<https://mp.weixin.qq.com/s/uzGhAHVqjmgH8QA8eIsS4w>

自从 Kubernetes1.11 之后，CoreDNS 作为集群内默认的域名解析服务，你是否对它还仅仅还停留在对 Kubernetes 的 Service 解析呢？事实上光 DNS 在 K8S 内就有很多有意思的操作，今天我们不妨来看看 CoreDNS 的各种高阶玩法。

## 1. 自定义 hosts 解析

默认情况下，Kubernetes 集群内的容器要解析外部域名时，CoreDNS 会将请求转发给`/etc/resolv.conf`文件里指定的上游 DNS 服务器。这个是由这个配置决定的。

    forward . /etc/resolv.conf

有的时候，我们如果需要在`集群内全局劫持某个域名时`，我们通常可以利用`hosts`插件来帮忙。**hosts**插件会每隔 5s 将需解析的信息重新加载到 coredns 当中，当你有任何变化时直接更新它的配置区块即可。常见的 host 有两种方法配置，分别如下：

- 定义 host

<!---->

    .:53 {
        hosts {
            1.1.1.1 test.cloudxiaobai.com
            2.2.2.2 test2.cloudxiaobai.com
            fallthrough
        }
    }

- 加载 hosts 文件

<!---->

    #直接从/etc/hosts加载host信息
    . {
        hosts {
            fallthrough
        }
    }
    #又或者,从当前目录的test.hosts文件中加载host信息
    . {
        hosts test.hosts {
            fallthrough
        }
    }

> 当被需要解析的域名不在 hosts 当中时，需要用`fallthrough`继续将请求转发给其它插件继续处理

_扩展_
如果我们只是想在 Workloads 内局部生效部分 host 信息时，那么可以借助于`HostAliases向Pod的/etc/hosts文件内添加主机信息。`我们拿 deployment 来举例，

    apiVersion: extensions/v1beta1
    kind: Deployment
    spec:
      template
        spec:
          containers:
          - image: busybox:latest
            name: nginx
          hostAliases:
          - ip: 1.1.1.1
            hostnames:
            - test1.cloudxiaobai.com
          - ip: 2.2.2.2
            hostnames:
            - test1.cloudxiaobai.com
    ...

## 2. 支持 SRV 记录

**SRV 记录**是域名系统中用于指定服务器提供服务的位置（如主机名和端口）数据。它在 DNS 记录中的是个新鲜面孔，在 RFC2082 中才对 SRV 记录进行了定义，因此有很多`老旧服务器并不支持SRV记录`。SRV 在 RFC2082 定义的标准记录格式如下：

    #英文
    _Service._Proto.Name TTL Class SRV Priority Weight Port Target
    #中文
    _服务._协议.名称. TTL 类别 SRV 优先级 权重 端口 主机.

- Service ：服务的符号名称
- Proto ：服务的传输协议，通常为 TCP 或 UDP，Proto 不区分大小写
- Name ：此 RR 所指的域名，在这个域名下 SRV RR 是唯一的
- TTL ：标准 DNS 存活时间
- CLASS ：标准 DNS 类别值（此值总为 IN）
- Priority ：目标主机的优先级，值越小越优先，范围 0-65535
- Weight ：相同优先度记录的相对权重，值越大越优先
- Port ：服务所在的 TCP 或 UDP 端口
- Target : 提供服务的规范主机名，以半角句号结尾

在 Kubernetes 里面，CoreDNS 会为`有名称的端口创建SRV记录`，这些端口可以是 svc 或 headless.svc 的一部分。对每个命名端口，SRV 记录了一个类似下列格式的记录：

    _port-name._port-protocol.my-svc.my-namespace.svc.cluster.local

在 Golang 中我们用 net.LookupSRV 来发起 SRV 记录查询

    func (r *Resolver) LookupSRV(ctx context.Context, service, proto, name string) (cname string, addrs []*SRV, err error)

net 库里对 SRV 结构体里定义了 4 个字段，分别是`Target`,`Port`，`Priority`,`Wright`。当我们使用`LookupSRV`发起 SRV 查询时，得到的返回的记录会按优先级排序，并在优先级内按权重进行随机分配。如果 service 和 proto 均为空字符串，则 LookupSRV 直接查找 name。

### 拿 thanos 的 SRV 查询举个例子

#### 1. 第一步 resolver.go 中 SRV 查询逻辑

thanos 中的 resolver.go 里面包含了处理 SRV 查询的逻辑，如下：

```go
case SRV, SRVNoA:
  _, recs, err := s.resolver.LookupSRV(ctx, "", "", host)
  if err != nil {
   return nil, errors.Wrapf(err, "lookup SRV records %q", host)
  }
  for _, rec := range recs {
   resPort := port
   if resPort == "" {
   //获取SRV返回的端口
    resPort = strconv.Itoa(int(rec.Port))
   }
   if qtype == SRVNoA {
         //如果不需要使用A或者AAAA记录查询时，则组合主机名:端口
    res = append(res, appendScheme(scheme, net.JoinHostPort(rec.Target, resPort)))
    continue
   }
   // Do A lookup for the domain in SRV answer.
   resIPs, err := s.resolver.LookupIPAddr(ctx, rec.Target)
   if err != nil {
    return nil, errors.Wrapf(err, "look IP addresses %q", rec.Target)
   }
   //根据主机名遍历出所有的ip地址，并组合成ip:port的方式
   for _, resIP := range resIPs {
    res = append(res, appendScheme(scheme, net.JoinHostPort(resIP.String(), resPort)))
   }
```

#### 第二步 创建 Kubernetes Service

CoreDNS 中对于有名称的 port，会为其创建一条对应的 SRV 记录。

    apiVersion: v1
    kind: Service
    metadata:
      labels:
        app: thanos-receiver
      name: thanos-receiver
      namespace: monitor
    spec:
      ClusterIP: None
      ports:
      - name: grpc
        port: 10901
        protocol: TCP
        targetPort: 10901

如上结构的 Service 我们就可以使用如下方式查询域名:

    _grpc._tcp.thanos-receiver.monitor.svc.cluster.local

我们用 dig 命令做一次 SRV 记录查询就可以得到响应，可以看到下列的`ANSWER SECTION`中得到了 net 库里面定义的 SRV 结构体的数据了，并且 CoreDNS 返回了三条记录：

    # dig srv _grpc._tcp.thanos-receiver.monitor.svc.cluster.local
    ...
    ;; ANSWER SECTION:
    _grpc._tcp.thanos-receiver.monitor.svc.cluster.local. 6 IN SRV 0 33 10901 10-59-155-238.thanos-receiver.monitor.svc.cluster.local.
    _grpc._tcp.thanos-receiver.monitor.svc.cluster.local. 6 IN SRV 0 33 10901 10-59-42-68.thanos-receiver.monitor.svc.cluster.local.
    _grpc._tcp.thanos-receiver.monitor.svc.cluster.local. 6 IN SRV 0 33 10901 10-59-48-162.thanos-receiver.monitor.svc.cluster.local.
    ;; ADDITIONAL SECTION:
    10-59-48-162.thanos-receiver.monitor.svc.cluster.local. 6 IN A 10.59.48.162
    10-59-42-68.thanos-receiver.monitor.svc.cluster.local. 6 IN A 10.59.42.68
    10-59-155-238.thanos-receiver.monitor.svc.cluster.local. 6 IN A 10.59.155.238
    ...

可以看到这条 SRV 记录里面，分别返回了三个服务的`IP地址`、`端口`、以及`服务的优先级`和`权重`

#### 第三步 使用 SRV 记录做服务发现

对于代码中启用了 SRV 记录的业务，只需要在业务配置里面加上需要访问的 SRV 地址即可，例如 thanos-query 需要调 thanos-receiver 的 grpc 端口做监控数据查询，如果我们集群内有多个 receiver 服务的话，我们就像如下配置，即可做到 DNS 的服务发现：

    ...
        spec:
          containers:
          - args:
            - query
            # 定义thanos-receiver服务SRV记录
            - --store=dnssrv+_grpc._tcp.thanos-receiver.monitor.svc.cluster.local
    ...

当服务正常运行后我们就可以查到 receiver 服务以及注册到 query 里面了
![image.gif](https://notes-learning.oss-cn-beijing.aliyuncs.com/pfyc62/1622810154712-9867e005-f445-4b25-a00d-4a1f623fb7ee.gif)

## 3. NodeLocal DNSCache

有很多同学经常会抱怨，在 Kubernetes 中有时候会遇到 DNS 解析间歇性 5s 超时的问题。其实这个问题社区很早意识到 DNS 的经过 Iptables 到 Conntrack 遇到竞争的问题，并给出来利用 Daemonset 在集群的每个 Node 上运行一个精简版的 CoreDNS 并监听一个虚拟 ip 地址来绕过 Conntrack，同时还能充当缓存环境 CoreDNS 压力。此举能大幅降低 DNS 查询 timeout 的频次，提升服务稳定性。
![image.gif](https://notes-learning.oss-cn-beijing.aliyuncs.com/pfyc62/1622810154942-09a54d12-aba5-46f9-a936-7ac84725e2b2.gif)**关于部署**
**node-local-dns**通过添加 iptables 规则能够接收节点上所有发往 169.254.20.10 的 dns 查询请求，把针对集群内部域名查询请求路由到 coredns。把集群外部域名请求直接通过 host 网络发往本地`/etc/resolv.conf`记录的外部 DNS 服务器中。

    # 下载部署脚本
    $ curl https://node-local-dns.oss-cn-hangzhou.aliyuncs.com/install-nodelocaldns.sh
    # 部署,确保kubectl能够连接集群
    $ bash install-nodelocaldns.sh

### 如何使用

`NodeLocal DNSCache的部署并不会直接产生效果`,通常我们有两种方式可以让集群的 pod 使用上本机 DNS 缓存。

#### 1. 定制业务容器 dnsConfig

Kubernetes 的 workload 中允许我们自定义 dns 相关的配置，其中我们需要注意以下几点：

- dnsPolicy: None，不使用 ClusterDNS。
- 配置 searches，保证集群内部域名能够被正常解析。
- 适当降低 ndots 值，当前 ACK 集群 ndots 值默认为 5,降低 ndots 值有利于加速集群外部域名访问。
- 适当调整 options 参数，避免并发请求`single-request`和分开 A 和 AAAA 请求采用的源端口`single-request-reopen`

可以参考如下

    dnsPolicy: None
    dnsConfig:
        nameservers: ["169.254.20.10"]
        searches:
        - default.svc.cluster.local
        - svc.cluster.local
        - cluster.local
        options:
        - name: ndots
          value: "2"
        - name: single-request-reopen
          value: ""
        - name: timeout
          value: "1"

#### 2. 修改 Kubelet 配置

kubelet 启动参数中可以通过参数`--cluster-dns`来指定容器的 nameserver，我们只需将它修改成`169.254.20.10`重启即可。不过容器要真正将 NodeLocal DNSCache 用起来话，还得将`Pod重启`才会生效。

## 4. 禁用 IPv6 域名解析

有时候我们 Kubernetes 集群内没有启用 IPv6 的话，可以在 CoreDNS 内禁止 IPv6 的域名解析，这个时候我们可以用 Template 这个插件来解决：

    .:53 {
        template ANY AAAA {
            rcode NXDOMAIN
        }
    ...
    }

> 这条记录会将所有的 AAAA 查询直接返回`NXDOMAIN`,并且不会被转发给其它插件处理
