---
title: etcd 问题处理案例
---



## What does the etcd warning “apply entries took too long” mean?

官方文档：

<https://etcd.io/docs/v3.4.0/faq/#what-does-the-etcd-warning-apply-entries-took-too-long-mean>

<https://github.com/etcd-io/etcd/blob/master/Documentation/faq.md#what-does-the-etcd-warning-apply-entries-took-too-long-mean>

## What does the etcd warning “failed to send out heartbeat on time” mean?

官方文档：

<https://etcd.io/docs/v3.4.0/faq/#what-does-the-etcd-warning-failed-to-send-out-heartbeat-on-time-mean>

<https://github.com/etcd-io/etcd/blob/master/Documentation/faq.md#what-does-the-etcd-warning-failed-to-send-out-heartbeat-on-time-mean>

# 其他

<http://blog.itpub.net/31559758/viewspace-2704804/>

# etcdserver: read-only range request took too long to execute

问题原因：

有可能是磁盘性能导致，当磁盘性能只有 2，3 十兆的读写速度，那么有很大机率会出现此错误。

参考：<https://github.com/kubernetes/kubernetes/issues/70082>