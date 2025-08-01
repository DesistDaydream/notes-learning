---
title: Kubesphere
linkTitle: Kubesphere
weight: 20
---

# 概述

> 参考：
>
> - 

> [!Attention] 自本公告发布之日起，将暂停KubeSphere开源版产品的下载链接，同时停止提供免费的技术支持。我们深知可能会给部分用户的使用带来不便，对此我们深表歉意。但我们相信，通过资源集中整合，能够为有持续需求的用户提供更专业、更稳定、更全面的商业级服务与支持。
>
> https://github.com/kubesphere/kubesphere/issues/6550

kubekey 部署工具小问题

kubeadm 文件无法完全自定义，在 apis/kubekey/v1alpha1/kubernetes_types.go 中的 kubernetes 结构体只有很少的几个属性，pkg/kubernetes/tmpl/kubeadm.go 中的 kubeadm 模板文件很多都是写死的。

# 问题总结

总结于 3.1

该 产品定位是集群管理，应用管理功能非常弱。

企业空间下的项目=k8s 中的 namespace

应用仓库没法使用认证添加

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/kubesphere/1619746908417-ebf183bd-a231-4518-a5b9-62888673fc65.png)

应用信息加载不出来

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/kubesphere/1619746661900-84a85e25-d389-42d7-a11e-6c37506b0dcd.png)

监控套件输入内嵌，自定义非常弱，不太懂为啥要自己实现 Grafana

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/kubesphere/1619747100329-da9e8052-a3e9-4cdd-a1a3-54563abb58cd.png)

有的代理转发的头不支持

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/kubesphere/1619746684706-fddd0c38-e6b7-4138-b236-1ad455468715.png)

总结，用开源的东西，又想要实现生态闭环，抽象概念太多，用上的人就离不开，离开了也没法理解 k8s。

该产品目标猜测：面向企业，对企业的基础设施以上的设备及应用的全生命周期管理。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/kubesphere/1652019279389-df197850-1ba6-4c0a-aaba-ff4ea27ad936.png)
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/kubesphere/1652019250255-d576b173-f392-46a7-b08d-6bf8da5b776f.png)

已经无人维护了~~~
