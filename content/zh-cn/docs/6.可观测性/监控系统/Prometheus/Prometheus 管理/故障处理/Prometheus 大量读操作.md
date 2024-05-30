---
title: Prometheus 大量读操作
---

故障现象

Prometheus 的读请求无故瞬间激增。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/mk7rw5/1616068308575-35d1480d-561e-494a-be22-fb863770bbb9.png)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/mk7rw5/1616068308611-92cbab4f-01c8-4a00-844c-3944f3924dd3.png)

# 故障排查

重启 Prometheus 后解决，后续需要跟进看是否还会继续发生

当使用 Grafana 查询 30 天的指标时，Prometheus 的读请求就会激增：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/mk7rw5/1616068308639-c7a19793-f268-4562-9f5b-1deea68be7e1.png)

怀疑可能是 Grafana 与 Prometheus 之间的连接没有中断，持续查询导致，但是暂无证据

1 月 20 日早晨进行 30 天查询后再次出现该问题，添加 netfilter 规则，阻断 grafana 与 prometheus 问题依旧；使用 docker 重启 prometheus 容器问题依旧；删除 grafana 问题依旧。

故障处理

实际上是由于每次评估规则时，有很多条规则的表达式是 30 天的范围表达式导致的。
