---
title: k8s 创建 Pod 时，背后到底发生了什么
---

# 概述

> 参考：
> - [公众号,万字长文：K8S 创建 Pod 时，背后到底发生了什么](https://mp.weixin.qq.com/s/HjoU_RKBQKPCQPEQZ_fBNA)

典型的创建 Pod 的流程为
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/zhow5n/1616119512783-67ed1273-0291-4462-8535-1ea845b176f1.png)

1. 用户通过 REST API 创建一个 Pod

2. apiserver 将其写入 etcd

3. scheduluer 检测到未绑定 Node 的 Pod，开始调度并更新 Pod 的 Node 绑定

4. kubelet 检测到有新的 Pod 调度过来，通过 container runtime 运行该 Pod

5. kubelet 通过 container runtime 取到 Pod 状态，并更新到 apiserver 中

各组件默认所占用端口号
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/zhow5n/1616119523444-d2794850-3f4c-41c8-8f75-e168a8825177.png)
