---
title: Error on ingesting out-of-order samples
---

# 故障现象

查看日志发现很多 `Error on ingesting out-of-order samples` Warn 信息

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/rpaa9g/1629643964887-eaa5bc33-94fb-4add-8424-f40dfd65ec02.png)

# 故障原因

> 参考：
> - <https://www.robustperception.io/debugging-out-of-order-samples>

当一个 job 中从多个 Prometheus 中采集相同指标时，就容易产生这个问题。比如，下图示例：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rpaa9g/1629687458928-d2444080-a4ff-406c-8a70-76fa687459ae.jpeg)

当采集目标是具有相同数据的多个 Prometheus，并且采集时轮流采集，就会很容易产生上述问题

# 故障处理

每个 Job 配置中，添加 `honor_timestamps: false` 配置。
