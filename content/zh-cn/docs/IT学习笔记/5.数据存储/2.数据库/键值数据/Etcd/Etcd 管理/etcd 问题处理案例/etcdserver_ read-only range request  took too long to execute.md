---
title: "etcdserver: read-only range request  took too long to execute"
---

#

错误信息： etcdserver: read-only range request took too long to execute

问题原因：

有可能是磁盘性能导致，当磁盘性能只有 2，3 十兆的读写速度，那么有很大机率会出现此错误。

参考：

<https://github.com/kubernetes/kubernetes/issues/70082>
