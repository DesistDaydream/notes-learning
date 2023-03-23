---
title: Pipeline 概念
---

# 概述

> 参考：
>
> - [官方文档](https://grafana.com/docs/loki/latest/clients/promtail/pipelines/)
> - [公众号,Promtail Pipeline 日志处理配置](https://mp.weixin.qq.com/s/PPNa7CYk6aaYDcvH9eTw1w)

Pipeline 用来处理 tail 到的每一行日志的内容、标签、时间戳。Pipeline 的行为在配置文件的 `.scrape_config.pipeline_stages` 字段定义。是 Promtail 处理日志必不可少的一个环节。

Pipeline 由一组 **stages(阶段)** 组成，Loki 将 Stages 分为 4 大类型：

1. **Parsing stages(解析阶段)**# 解析每行日志，并从中提取数据。提取的数据可供后面几个阶段使用

2. **Transform stages(转换阶段)** # (可省略)转换解析阶段提取到的数据

3. **Actions stages(行动阶段)**# (可省略)处理转换阶段转换后的数据。行动包括以下几种

   1. 为每行日志添加标签或修改现有标签

   2. 更改每行日志的时间戳

   3. 改变日志行内容

   4. 根据提取到的数据创建 metrics(指标)

4. **Filtering stages(过滤阶段)** # (可省略)根据指定的条件，保留或删除日志行。
   1. 注意：过滤阶段的类型中，有一个名为 **match** 的过滤阶段。match 是一个通用的阶段，不受阶段顺序影响，在处理日志行之前，match 阶段可以使用 [LogQL](/docs/6.可观测性/日志系统/Loki/LogQL.md)，来过滤要使用某些阶段进行处理的日志行。

## 各阶段类型

- **Parsing stages(解析阶段)** 类型：

  - [cri](https://grafana.com/docs/loki/latest/clients/promtail/stages/cri/) # 使用标准的 CRI 日志格式来解析每行日志，并提取数据
  - [docker](https://grafana.com/docs/loki/latest/clients/promtail/stages/docker/) # 使用标准的 docker 日志文件格式来解析每行日志，并提取数据(Pipeline 的默认行为，该阶段包括 json、labels、timestamp、output 四个阶段)
  - [regex](https://grafana.com/docs/loki/latest/clients/promtail/stages/regex/) # 使用正则表达式从每行日志提取数据

  - [json](https://grafana.com/docs/loki/latest/clients/promtail/stages/json/) # 使用 JSON 格式解析每行日志，并提取数据

  - [replace](https://grafana.com/docs/loki/latest/clients/promtail/stages/replace/) # 使用正则表达式替换数据

- **Transform stages(转换阶段)** 类型：
  - [multiline](https://grafana.com/docs/loki/latest/clients/promtail/stages/multiline/) # 多行阶段将多行日志进行合并，然后再将其传递到 pipeline 的下一个阶段。
  - [pack](https://grafana.com/docs/loki/latest/clients/promtail/stages/pack/) # Packs a log line in a JSON object allowing extracted values and labels to be placed inside the log line.
  - [template](https://grafana.com/docs/loki/latest/clients/promtail/stages/template/) # 使用 Go 模板来修改提取出来数据
- **Actions stages(行动阶段)** 类型：

  - [timestamp](https://grafana.com/docs/loki/latest/clients/promtail/stages/timestamp/) # 为一行日志设置时间戳

  - [output](https://grafana.com/docs/loki/latest/clients/promtail/stages/output/) # 设置一行日志的文本。该行为是 pipeline 阶段可以确定 loki 要展示的日志内容的唯一行为

  - [labels](https://grafana.com/docs/loki/latest/clients/promtail/stages/labels/) # 更新日志条目的标签集
  - [labelallow](https://grafana.com/docs/loki/latest/clients/promtail/stages/labelallow/) # 保留标签
  - [labeldrop](https://grafana.com/docs/loki/latest/clients/promtail/stages/labeldrop/) # 丢掉标签

  - [metrics](https://grafana.com/docs/loki/latest/clients/promtail/stages/metrics/) # 根据提取出来的数据计算指标

  - [tenant](https://grafana.com/docs/loki/latest/clients/promtail/stages/tenant/) # 设置要用于日志条目的租户 ID 值。

- **Filtering stages(过滤阶段)** 支持以下行为

  - [match](https://grafana.com/docs/loki/latest/clients/promtail/stages/match/) # 依据指定的标签，过滤日志行，只有匹配到的日志行才会继续执行其他阶段
  - [drop](https://grafana.com/docs/loki/latest/clients/promtail/stages/drop/) # 依据条件丢弃日志行

## 配置示例

一个典型的 pipeline 将从解析阶段开始（如 regex 或 json 阶段）从日志行中提取数据。然后有一系列的处理阶段配置，对提取的数据进行处理。最常见的处理阶段是一个 `labels stage` 标签阶段，将提取的数据转化为标签。

需要注意的是现在 pipeline 不能用于重复的日志，例如，Loki 将多次收到同一条日志行：

- 从同一文件中读取的两个抓取配置
- 文件中重复的日志行被发送到一个 pipeline，不会做重复数据删除

然后，Loki 会在查询时对那些具有完全相同的纳秒时间戳、标签与日志内容的日志进行一些重复数据删除。

下面的配置示例可以很好地说明我们可以通过 pipeline 来对日志行数据实现什么功能：

```yaml
scrape_configs:
  - job_name: kubernetes-pods-name
    kubernetes_sd_configs: ....
    pipeline_stages:
      # 这个阶段只有在被抓取地目标有一个标签名为 name 且值为 promtail 地时候才会执行
      - match:
          selector: '{name="promtail"}'
          stages:
            # regex 阶段解析出一个 level、timestamp 与 component，在该阶段结束时，这几个值只为 pipeline 内部设置，在以后地阶段可以使用这些值并决定如何处理他们。
            - regex:
                expression: '.*level=(?P<level>[a-zA-Z]+).*ts=(?P<timestamp>[T\d-:.Z]*).*component=(?P<component>[a-zA-Z]+)'

            # labels 阶段从前面地 regex 阶段获取 level、component 值，并将他们变成一个标签，比如 level=error 可能就是这个阶段添加地一个标签。
            - labels:
                level:
                component:

            # 最后，时间戳阶段采用从 regex 提取地 timestamp，并将其变成日志的新时间戳，并解析为 RFC3339Nano 格式。
            - timestamp:
                format: RFC3339Nano
                source: timestamp

      # 这个阶段只有在抓取的目标标签为 name，值为 nginx，并且日志行中包含 GET 字样的时候才会执行
      - match:
          selector: '{name="nginx"} |= "GET"'
          stages:
            # regex 阶段通过匹配一些值来提取一个新的 output 值。
            - regex:
                expression: \w{1,3}.\w{1,3}.\w{1,3}.\w{1,3}(?P<output>.*)

            # output 输出阶段通过将捕获的日志行设置为来自上面 regex 阶段的输出值来更改其内容。
            - output:
                source: output

      # 这个阶段只有在抓取到目标中有标签 name，值为 jaeger-agent 时才会执行。
      - match:
          selector: '{name="jaeger-agent"}'
          stages:
            # JSON 阶段将日志行作为 JSON 字符串读取，并从对象中提取 level 字段，以便在后续的阶段中使用。
            - json:
                expressions:
                  level: level

            # 将上一个阶段中的 level 值变成一个标签。
            - labels:
                level:

  - job_name: kubernetes-pods-app
    kubernetes_sd_configs: ....
    pipeline_stages:
      # 这个阶段只有在被抓取的目标的标签为 "app"，名称为grafana 或 prometheus 时才会执行。
      - match:
          selector: '{app=~"grafana|prometheus"}'
          stages:
            # regex 阶段将提取一个 level 合 componet 值，供后面的阶段使用，允许 level 被定义为 lvl=<level> 或 level=<level>，组件被定义为 logger=<component> 或 component=<component>
            - regex:
                expression: ".*(lvl|level)=(?P<level>[a-zA-Z]+).*(logger|component)=(?P<component>[a-zA-Z]+)"

            # 然后标签阶段将从上面 regex 阶段提取的 level 和 component 变为标签。
            - labels:
                level:
                component:

      # 只有当被抓取的目标有一个标签 "app"，其值为 "some-app"，并且日志行不包含 "info" 一词时，这个阶段才会执行。
      - match:
          selector: '{app="some-app"} != "info"'
          stages:
            # regex 阶段尝试通过查找日志中的 panic 来提取 panic 信息
            - regex:
                expression: ".*(?P<panic>panic: .*)"

            # metrics 阶段将增加一个 Promtail 暴露的 panic_total 指标，只有当从上面的 regex 阶段获取到 panic 值的时候，该 Counter 才会增加。
            - metrics:
                panic_total:
                  type: Counter
                  description: "total count of panic"
                  source: panic
                  config:
                    action: inc
```
