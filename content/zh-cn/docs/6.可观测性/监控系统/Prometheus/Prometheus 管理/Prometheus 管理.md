---
title: Prometheus 管理
---

# 重大变化

## V2.39

> 参考：
>
> - <https://mp.weixin.qq.com/s/RMtjCiWgTFnKhnTBQc-WLA>

大量的资源优化。改进了 relabeling 中的内存重用，优化了 WAL 重放处理，从 TSDB head series 中删除了不必要的内存使用， 以及关闭了 head compaction 的事务隔离等。
