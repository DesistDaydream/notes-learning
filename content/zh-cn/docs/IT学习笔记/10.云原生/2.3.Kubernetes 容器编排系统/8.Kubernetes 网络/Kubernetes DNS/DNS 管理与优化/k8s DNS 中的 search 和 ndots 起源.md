---
title: k8s DNS 中的 search 和 ndots 起源
---

# 多余的 DNS 查询

一些需要解析外部 DNS 域名的应用，当运行在容器中时，如果我们在容器的 network namespace 中对 dns 报文（udp port 53）进行抓包，可能会发现在正确解析之前，还经过了若干次多余的尝试。

下面是我在容器中`ping google.com`，同时在容器的 network namespace 中抓到的包。

    sudo nsenter -t 3885 -n tcpdump -i eth0 udp port 53
    tcpdump: verbose output suppressed, use -v or -vv for full protocol decode
    listening on eth0, link-type EN10MB (Ethernet), capture size 262144 bytes
    10:09:11.917900 IP 10.244.2.202.38697 > 10.96.0.10.domain: 11858+ A? google.com.default.svc.cluster.local. (54)
    10:09:11.918847 IP 10.96.0.10.domain > 10.244.2.202.38697: 11858 NXDomain*- 0/1/0 (147)
    10:09:11.922468 IP 10.244.2.202.38697 > 10.96.0.10.domain: 15573+ AAAA? google.com.default.svc.cluster.local. (54)
    10:09:11.923001 IP 10.96.0.10.domain > 10.244.2.202.38697: 15573 NXDomain*- 0/1/0 (147)
    10:09:11.923248 IP 10.244.2.202.43230 > 10.96.0.10.domain: 62042+ A? google.com.svc.cluster.local. (46)
    10:09:11.923828 IP 10.96.0.10.domain > 10.244.2.202.43230: 62042 NXDomain*- 0/1/0 (139)
    10:09:11.924005 IP 10.244.2.202.43230 > 10.96.0.10.domain: 54769+ AAAA? google.com.svc.cluster.local. (46)
    10:09:11.924494 IP 10.96.0.10.domain > 10.244.2.202.43230: 54769 NXDomain*- 0/1/0 (139)
    10:09:11.924704 IP 10.244.2.202.36252 > 10.96.0.10.domain: 20727+ A? google.com.cluster.local. (42)
    10:09:11.925154 IP 10.96.0.10.domain > 10.244.2.202.36252: 20727 NXDomain*- 0/1/0 (135)
    10:09:11.925316 IP 10.244.2.202.36252 > 10.96.0.10.domain: 13066+ AAAA? google.com.cluster.local. (42)
    10:09:11.925758 IP 10.96.0.10.domain > 10.244.2.202.36252: 13066 NXDomain*- 0/1/0 (135)
    10:09:11.925929 IP 10.244.2.202.35582 > 10.96.0.10.domain: 38821+ A? google.com.lan. (32)
    10:09:11.927244 IP 10.244.2.202.35582 > 10.96.0.10.domain: 4430+ AAAA? google.com.lan. (32)
    10:09:11.927416 IP 10.96.0.10.domain > 10.244.2.202.35582: 38821 NXDomain 0/0/0 (32)
    10:09:11.928600 IP 10.96.0.10.domain > 10.244.2.202.35582: 4430 NXDomain 0/0/0 (32)
    10:09:11.928839 IP 10.244.2.202.45290 > 10.96.0.10.domain: 45577+ A? google.com. (28)
    10:09:11.929129 IP 10.244.2.202.45290 > 10.96.0.10.domain: 37586+ AAAA? google.com. (28)
    10:09:11.929303 IP 10.96.0.10.domain > 10.244.2.202.45290: 45577 1/0/0 A 172.217.160.78 (54)
    10:09:11.929541 IP 10.96.0.10.domain > 10.244.2.202.45290: 37586 1/0/0 AAAA 2404:6800:4008:801::200e (66)

可以看到，在最后（倒数第 3、4 行）正确解析之前，先是依次查询了下面几个域名，并且均查询了 IPv4 和 IPv6：

- `google.com.default.svc.cluster.local.`

- `google.com.svc.cluster.local.`

- `google.com.cluster.local.`

- `google.com.lan.`

但是这 8 次查询都失败了，因为并不存在这样的域名。

### kubernetes 的容器域名解析&#xA;

要想解释上面的现象，需要先从 kubernetes 的容器域名解析开始讲。

kubernetes 上运行的容器，其域名解析和一般的 Linux 一样，都是根据 `/etc/resolv.conf` 文件。如下是容器中该文件的内容。

    nameserver 10.96.0.10
    search default.svc.cluster.local svc.cluster.local cluster.local lan
    options ndots:5

nameserver 即为 kubernetes 集群中，kube-dns 的 svc IP，集群中容器的 nameserver 均设置为 kube-dns。

    kubectl get svc -n kube-system
    NAME             TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                  AGE
    kube-dns         ClusterIP   10.96.0.10      <none>        53/UDP,53/TCP,9153/TCP   236d

在容器中，自动添加了 `search` 和 `ndots`

### kubernetes 为什么使用搜索域

为什么呢？先来看看代码。

    var (
    	// The default dns opt strings.
    	defaultDNSOptions = []string{"ndots:5"}
    )


    func (c *Configurer) generateSearchesForDNSClusterFirst(hostSearch []string, pod *v1.Pod) []string {
    	if c.ClusterDomain == "" {
    		return hostSearch
    	}
    	nsSvcDomain := fmt.Sprintf("%s.svc.%s", pod.Namespace, c.ClusterDomain)
    	svcDomain := fmt.Sprintf("svc.%s", c.ClusterDomain)
    	clusterSearch := []string{nsSvcDomain, svcDomain, c.ClusterDomain}
    	return omitDuplicates(append(clusterSearch, hostSearch...))
    }


    func (c *Configurer) GetPodDNS(pod *v1.Pod) (*runtimeapi.DNSConfig, error) {
    	...
    	case podDNSCluster:
    		if len(c.clusterDNS) != 0 {
    			dnsConfig.Servers = []string{}
    			for _, ip := range c.clusterDNS {
    				dnsConfig.Servers = append(dnsConfig.Servers, ip.String())
    			}
    			dnsConfig.Searches = c.generateSearchesForDNSClusterFirst(dnsConfig.Searches, pod)
    			dnsConfig.Options = defaultDNSOptions
    			break
    		}
    	...
    	if utilfeature.DefaultFeatureGate.Enabled(features.CustomPodDNS) && pod.Spec.DNSConfig != nil {
    		dnsConfig = appendDNSConfig(dnsConfig, pod.Spec.DNSConfig)
    	}
    }

kubernetes 搜索域

从函数`generateSearchesForDNSClusterFirst`中可见，搜索域有三个：nsSvcDomain、svcDomain、clusterDomain。

kubernetes 之所以要设置搜索域，目的是为了方便用户访问 service。

例如，default namespace 下的 Pod a，如果访问同 namespace 下的 service b，直接使用 b 就可以访问了，而这个功能，就是通过 nsSvcDomain 搜索域`default.svc.cluster.local`完成的。类似的，对于不同 namespace 下的 service，可以用 `${service name}.${namespace name}` 来访问，是通过 svcDomain 搜索域完成的。clusterDomain 设计的目的，是为了方便同域中非 kubernetes 上的域名访问，例如设置 kubernetes 的 domain 为`ieevee.com`，那么对于`s.ieevee.com`域名，直接使用`s`来访问就可以了，当然前提是当前 namespace 中没有一个叫做`s`的 svc。（是的。搜索域是有优先级的）

ndots 默认值

`ndots`默认值是写死的，5。

为什么是 5 呢？

thockin 在[ issue 33554](https://github.com/kubernetes/kubernetes/issues/33554) 中做了解释，概况来说：

1. kubernetes 需要支持同 namespace 下 service 快速访问，例如`name`，因此 ndots>=1，对应搜索域`$namespace.svc.$zone`

2. kubernetes 需要支持跨 namespace 下 service 快速访问，例如`kubernetes.default`，因此 ndots>=2，对应搜索域`svc.$zone`

3. kubernetes 需要支持同 namespace、跨 namespace 下，非 service 名称的快速访问，例如`name.namespace.svc`，因此 ndots>=3，对应搜索域`$zone`

4. kubernetes 需要支持 statefulset 中的每个 pod 的访问，例如`mysql-0.mysql.default.svc`，因此 ndots>=4

5. kubernetes 需要支持 SRV records（`_$port._$proto.$service.$namespace.svc.$zone`），因此 ndots>=5

不过呢，如果你的使用情况并不像上面这么复杂，这个值可能并不适合你。

比如说，我们只会用到同 namespace 下(形如 `a`)、跨 namespace 下的 service 访问（形如 `a.b`），因此，ndots 默认值为 2 更合适，但该值是写死在代码中的，不支持定制化，但可以通过下面的方法修改。

ndots 修改

ndots 是可以被修改的，可以通过`pod.Spec.DNSConfig`改写。

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
            value: "2"

通过上面的修改，我们再来看容器中的 DNS 报文，就只有下面几条了。

    sudo nsenter -t 3885 -n tcpdump -i eth0 udp port 53
    tcpdump: verbose output suppressed, use -v or -vv for full protocol decode
    listening on eth0, link-type EN10MB (Ethernet), capture size 262144 bytes
    10:30:35.917282 IP 10.244.2.202.39480 > 10.96.0.10.domain: 60870+ A? google.com. (28)
    10:30:35.919194 IP 10.244.2.202.39480 > 10.96.0.10.domain: 9908+ AAAA? google.com. (28)
    10:30:35.927047 IP 10.96.0.10.domain > 10.244.2.202.39480: 60870 1/0/0 A 216.58.200.238 (54)
    10:30:35.929089 IP 10.96.0.10.domain > 10.244.2.202.39480: 9908 1/0/0 AAAA 2404:6800:4008:801::200e (66)

这对于应用的性能会有一定的提升，具体可以参见 Ref。

### dns cache

我们观察一些 python 类型的容器应用，会发现它们会发出非常多的 DNS 请求，几乎每次涉及到域名都需要发出 DNS 请求；但是对于一些 Java 应用，会发现它们发出的 DNS 请求非常规律，一般是每 30 秒一个。

这是因为，不同语言对于 dns cache 的处理不同了。

目前来看，只有 Java 做了 dns cache。从 JDK 1.6 开始，Java 默认会对 DNS 做缓存，主要是以下两个配置：

- networkaddress.cache.ttl：域名解析成功后，DNS 缓存时间，默认是 30 秒

- networkaddress.cache.negative.ttl：域名解析失败后，冷却时间，默认是 10 秒

其他语言可以支持，但都需要一定 hack。

- go

- python

### 总结&#xA;

本文从一个多余的 DNS 查询现象开始，介绍了 FQDN 和 DNS 的搜索域，回答了 kubernetes 为什么需要搜索域这个问题，并且提出了一个解决多余 DNS 查询的方案，以及介绍了不同语言对于 dns cache 的处理情况。

Ref：

- Kubernetes pods /etc/resolv.conf ndots:5 option and why it may negatively affect your application performances

- About fully qualified domain names (FQDNs)

- Kubernetes pods /etc/resolv.conf ndots:5 option and why it may negatively affect your application performances
