---
title: Data Pipeline
linkTitle: Data Pipeline
weight: 20
---

# 概述

> 参考：
>
> - https://www.mezmo.com/blog/what-is-an-observability-data-pipeline
> - https://vector.dev/docs/about/what-is-observability-pipelines/
> - https://www.dqlabs.ai/blog/what-is-a-data-pipeline-types-architecture-components/

**Data Pipeline(数据管道)** 是一个抽象概念，通过 Data Pipeline 对多个来源的可观测性数据进行 聚合、处理、并将其路由到各种目的地。这解决了多个问题，包括：

1. 需要将数据集中到一个位置。
2. 结构化和丰富数据的能力，以便更容易理解并从中获取价值。
3. 需要将数据发送到多个目的地或团队以实现多个用例。
4. 需要控制数据量并仅以正确的格式将正确的数据发送到正确的目的地。

![https://docs.datadoghq.com/observability_pipelines/](https://datadog-docs.imgix.net/images/observability_pipelines/op_marketecture_08232024.82939fcff194b8a049c85df485665720.png?fit=max&auto=format&w=1862&h=894)

不管是 集成在程序内部的一段代码、在外部运行的采集程序，都可以看作是一种 DataPipeline。e.g. Loki 的 [Promtail](docs/6.可观测性/DataPipeline/Promtail/Promtail.md)、Prometheus 的 [Instrumenting](/docs/6.可观测性/Metrics/Instrumenting/Instrumenting.md)、etc. 都属于一种 DataPipeline 的实现。

<font color="#ff0000">**用人话说：将数据从一个地方流到另一个地方的行为，就是 Pipeline 行为**</font>。就像 Pipeline 名字一样，管道，只不过是用来承载数据的管道，可以让数据从管道的一段流向另一端。

从逻辑上看，可观测数据从采集到入库的过程中，可以存在多级 DataPipeline，采集程序自身也可以当作 DataPipeline 的一部分；然后中间可以经过处理或路由，也可以不经过处理或路由；最后进入到数据库中。
