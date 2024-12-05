---
title: Promtail
linkTitle: Promtail
date: 2024-06-15T08:47
weight: 1
---

# 概述

> 参考：
>
> - [官方文档](https://grafana.com/docs/loki/latest/clients/promtail/)
> - [GitHub 官方文档](https://github.com/grafana/loki/tree/main/docs/sources/clients/promtail)
> - [公众号，Promtail Pipeline 日志处理配置](https://mp.weixin.qq.com/s/PPNa7CYk6aaYDcvH9eTw1w)

Promtail 是将本地日志内容发送到私有 Loki 或 Grafana Cloud 的代理。通常将其部署到 有监控需求的应用程序 的每台机器上。

promtail 通过类似于 `tail` 命令的这种方式来采集日志文件内容，采集完成后，添加 label，然后 push 给 Loki 。

Promtail 是 Loki 官方支持的日志采集端，在需要采集日志的节点上运行采集代理，再统一发送到 Loki 进行处理。除了使用 Promtail，社区还有很多采集日志的组件，比如 fluentd、fluent bit 等，都是比较优秀的。

但是 Promtail 是运行 Kubernetes 时的首选客户端，因为你可以将其配置为自动从 Promtail 运行的同一节点上运行的 Pod 中抓取日志。Promtail 和 Prometheus 在 Kubernetes 中一起运行，还可以实现非常强大的调试功能，如果 Prometheus 和 Promtail 使用相同的标签，用户还可以使用 Grafana 根据标签集在指标和日志之间切换。

此外如果你想从日志中提取指标，比如计算某个特定信息的出现次数，Promtail 效果也是非常友好的。在 Promtail 中一个 pipeline 管道被用来转换一个单一的日志行、标签和它的时间戳。

当前，Promtail 可以从两个来源 tail 日志

1. local log 本地日志文件
2. systemd journal

Promtail 与 Filebeat 性能对比图

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zl113s/1616129582749-70275c0e-893e-413b-b489-80092db526c9.png)

# Promtail 工作流程

## 日志发现

> 参考：
>
> - [Loki 官方文档，发送数据 - promtail - Scraping](https://grafana.com/docs/loki/latest/send-data/promtail/scraping/)

Promtail 与 Prometheus 的服务发现机制相同，通过配置文件中 `scrape_configs` 字段的内容，来发现需要采集日志的目标，同时发现标签，然后通过 Relabeling 行为对 要抓取的内容、要丢弃的内容、以及要附加到日志行的标签 进行细粒度的控制。

## 标签与 parse(解析)

标签与解析步骤分为两部分

1. 根据元数据添加标签。发现日志流后，会确定元数据(pod 名称、文件名等等)，这些元数据可以作为标签附加到每行日志上。通过 relabel_configs，可以将发现的标签改变为所需的形式。
2. 解析日志内容并添加或更新标签。Promtail 会通过`scrape_config.pipeline_stages`配置段的内容，解析每行日志内容。根据解析的内容，可以为日志添加新的标签或更新现有标签。这种行为称为 [Pipelines](https://grafana.com/docs/loki/latest/clients/promtail/pipelines/)
   1. Pipeline 说明详见：[Promtail Pipeline 概念](/docs/6.可观测性/Logs/Loki/Promtail/Pipeline/Pipeline.md)

## 推送日志

Promtail 具有一组目标（i.e.要读取的内容，如文件）并且正确设置了所有标签后，它将开始跟踪（连续读取）来自目标的日志。一旦将足够的数据读取到内存中或经过可配置的超时后，就将其作为单批数据推送到 Loki。

当 Promtail 从源（文件和系统日志，如果已配置）读取数据时，它将跟踪在位置文件中读取的最后一个偏移量。默认情况下，头寸文件存储在 /var/log/positions.yaml 中。位置文件可帮助 Promtail 在 Promtail 实例重新启动的情况下从中断处继续读取。

## 总结

### Promtail 的工作流程实际上是一个循环

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zl113s/1616129582786-d73bcd27-d899-46f7-beb3-8dc84d0a690c.jpeg)
promtail 的日志发现过程是逐行发现，每发现一行处理一行，直到一定的时间或处理了足够多的日志，则将这些日志发送。然后继续下一个循环。这个过程，就好像瀑布的水流一样，哗哗得从上不断往下流水。日志文件内容就是水流，经过 promtail 哗哗往下流，流到 Loki 中。这就是 Loki 中 Log Stream 的形象描述~

经过上述处理后，在 Grafana 看到的日志最终效果就像下图一样：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zl113s/1616129582773-ccfea9e2-0ad0-4ca4-bee7-f3248542c7e7.jpeg)

# Promtail 关联文件与配置

**/etc/promtail/config.yml** # Promtail 程序运行时基本配置文件

**/run/promtail/positions.yaml** # 每个已发现的目标的日志文件路径，都会保存在该文件中。重新启动 Promtail 时需要使用该文件，以使其从日志文件中断处继续抓取日志。

## promtail 启动时的最小配置文件

- `__path__` 为必须存在的标签，promtail 根据该标签值指定的路径 tail 日志文件。
- `__path__` 标签就像时间序列中 **name** 标签的作用一样。且 promtail 会自动将该标签名转换为 filename。
- promtail 无法抓取 `__path__` 下指定路径的子路径下的文件。。。暂时还没找到解决办法

```yaml
# server 配置 promtail 程序运行时行为。如指定监听的ip、port等信息。
server:
  http_listen_port: 9080
  grpc_listen_port: 0
# positions 指定 promtail 读取日志的路径
positions:
  filename: /tmp/positions.yaml

# clients 配置 Promtail 如何连接到 Loki 的多个实例，并向每个实例发送日志。
# Note：如果其中一台远程Loki服务器无法响应或发生任何可重试的错误，这将影响将日志发送到任何其他已配置的远程Loki服务器。
# 发送是在单个线程上完成的！ 如果要发送到多个远程Loki实例，通常建议并行运行多个Promtail客户端。
clients:
  - url: http://localhost:3100/loki/api/v1/push

# scrape_configs 指定要抓取日志的目标。
scrape_configs:
  - job_name: system
    static_configs:
      - targets: # 指定抓取目标，i.e.抓取哪台设备上的文件
          - localhost
        labels: # 指定该日志流的标签
          __path__: /var/log/host/* # 指定抓取路径，该匹配标识抓取 /var/log/host 目录下的所有文件。注意：不包含子目录下的文件。
```
