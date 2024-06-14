---
title: API Server 配置详解
---

# 概述

> 参考：
>
> - [官方文档,参考-组件工具-kube-apiserver](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/)

API Server 现阶段只能通过命令行标志才可以改变运行时行为。暂无配置文件可用。

# kube-apiserver 命令行标志详解

**--allow-privileged \<BOOL>** # 是否允许有特权的容器。`默认值：false`。

**--basic-auth-file \<FILE>** # 配置 API Server 的基础认证。

- 该标志已于 1.19 版本彻底弃用。详见 [PR #89069](https://github.com/kubernetes/kubernetes/pull/89069)

**--insecure-port \<NUM>** # 开启不安全的端口。`默认值：0`，即不开启不安全的端口

**--insecure-bind-address \<IP>** # 不安全端口的监听地址。`默认值：127.0.0.1`。

**--runtime-config \<OBJECT>** # 启用或禁用内置的 APIs。

OBJECT 是 key=value 的键值对格式。key 为 API 组名称，value 为 true 或 false。

比如

- 要关闭 `batch/v1` 组，则设置 `--runtime-config=batch/v1=false`
- 要开启 `batch/v2alpha1` 组，则设置 `--runtime-config=batch/v2alpha1`

**--secure-port \<NUM>** # API Server 监听的安全端口。`默认值：6443`。不能以 0 关闭。

**--service-node-port-range \<PortRange>** # 指定 NodePort 类型的 service 资源可以使用的端口范围。

比如：'30000-32767'. Inclusive at both ends of the range。默认范围: 30000-32767

**--v NUM # 指定日志级别**

- 参考：
  - https://kubernetes.io/docs/reference/kubectl/cheatsheet/#kubectl-output-verbosity-and-debugging
  - https://github.com/kubernetes/community/blob/master/contributors/devel/sig-instrumentation/logging.md
- --v=0 # 通常对此有用，*始终*对运维人员可见。
- --v=1 # 如果您不想要详细程度，则为合理的默认日志级别。
- --v=2 # 有关服务的有用稳定状态信息以及可能与系统中的重大更改相关的重要日志消息。这是大多数系统的建议默认日志级别。
- --v=3 # 有关更改的扩展信息。
- --v=4 # Debug 级别。
- --v=6 # 显示请求的资源。
- --v=7 # 显示 HTTP 请求头。
- --v=8 # 显示 HTTP 请求内容。
- --v=9 # 显示 HTTP 请求内容而不截断内容。
