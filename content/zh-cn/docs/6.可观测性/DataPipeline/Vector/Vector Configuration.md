---
title: Vector Configuration
linkTitle: Vector Configuration
weight: 20
date: 2025-01-05T15:55:00
---

引言 官方文档的简单食用方式

在每个 [Source](https://vector.dev/docs/reference/configuration/sources/) 的文档中，可以在 **Output Data** 和 **Examples** 两个章节看到输出格式和输出内容。比如这种 Source 的日志都有什么元数据，输出格式是什么样的，etc. 。其中有些内嵌的隐藏字段以 `_` 开头，可以这些输出的内容可以直接或利用[模板](https://vector.dev/docs/reference/configuration/template-syntax/)使用在 sinks 的定义中。

用一个简单的 Linux 中的 [Journal](/docs/6.可观测性/Logs/Journal.md) 日志采集输出到 [Loki](/docs/6.可观测性/Logs/Loki/Loki.md) 的场景进行配置演示

```yaml
sources:
  test_journald:
    type: journald
sinks:
  test_loki:
    type: loki
    inputs:
      - test_journald
    endpoint: "http://localhost:3100"
    encoding:
      codec: raw_message
    labels:
      source_type: "{{ source_type }}"
      host: "{{ host }}"
      systemd_unit: "{{ _SYSTEMD_UNIT }}"
```

这个配置带有 `{{ }}` 的是模板语法，`_SYSTEMD_UNIT` 的值示例可以在 JournalD 文档的 [Examples](https://vector.dev/docs/reference/configuration/sources/journald/#examples) 看到；`source_type` 和 `host` 的值可以在文档的 [Output Data](https://vector.dev/docs/reference/configuration/sources/journald/#output-data) 看到

# 概述

> 参考：
>
> - [官方文档，参考 - 配置](https://vector.dev/docs/reference/configuration/)

Vector 配置文件支持 [YAML](/docs/2.编程/无法分类的语言/YAML.md), [TOML](/docs/2.编程/无法分类的语言/TOML.md), [JSON](/docs/2.编程/无法分类的语言/JSON.md) 三种格式。

顶级字段

- 全局字段 # 详见下文 [全局字段](#全局字段)。这些字段并没有上级字段，直接定义在文件顶层。
- 组件相关字段
  - **sources**(map\[SourceID][source](#sources)) # 定义一个或多个 Source 组件
  - **transforms**(map\[TransformID][transform](#transforms)) # 定义一个或多个 Transforms 组件
  - **sinks**(map\[SinkID][sink](#sinks)) # 定义一个或多个 Sink 组件

> [!Note]
> SourceID, TransformID, SinkID 是自定义的字符串，用以表示定义的对应组件的唯一性，称为 **ComponentID(组件唯一标识)**。在其地方引用已定义的组件时，也要使用 ComponentID
>
> 由于 Vector 灵活的设计，组件字段可以分散在多个文件中，并且可重复（比如两个文件中都有 sources 字段），只要 ComponentID 不重复即可。

# 全局字段

**data_dir**(STRING)

**api**(Object) # Vector 的 API 配置。开启后可以使用 `vector top` 命令尝试简单的 API 交互以验证 API 情况

- **enable**(BOOLEAN) # 是否开启 Vector API。`默认值: false`
- **address**(STRING) # API 监听地址。`默认值: 127.0.0.1:8686`

# sources

https://vector.dev/docs/reference/configuration/sources/

## Docker logs

https://vector.dev/docs/reference/configuration/sources/docker_logs

从 [Docker](docs/10.云原生/Containerization%20implementation/Docker/Docker.md) 中收集日志

**最佳实践**

```yaml
sources:
  demo_my_docker_logs:
    type: docker_logs
    include_containers:
      - nginx
```

# transforms

https://vector.dev/docs/reference/configuration/transforms/

# sinks

https://vector.dev/docs/reference/configuration/sinks/

所有 Sinks 通常都要有一个 `inputs([]STRING)` 字段，以声明要输出的数据是从哪来的。inputs 中的值可以是 SourceID 或 TransformID
