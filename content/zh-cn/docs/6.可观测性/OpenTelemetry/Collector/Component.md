---
title: Component
linkTitle: Component
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，open-telemetry/opentelemetry-collector-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib)

**Receivers** # https://github.com/open-telemetry/opentelemetry-collector/blob/main/receiver

**Processors** # https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor

**Exporters** # https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter

**Connector** # https://github.com/open-telemetry/opentelemetry-collector/blob/main/connector

**Extension** # https://github.com/open-telemetry/opentelemetry-collector/blob/main/extension

## 组件的实现

在 [open-telemetry/opentelemetry-collector-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib) 项目中包含了很多基于核心存储库 opentelemetry-collector 之上开发的一些具体实现。e.g. 从文件采集日志的 filelogreceiver、etc.

可以简单理解为：opentelemetry-collector 是核心，opentelemetry-collector-contrib 是核心的插件。前者用来定义各个组件应该如何组织，后者用来实现实现某类组件的具体实体的逻辑。

# Receiver

## Filelog

https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/filelogreceiver

从文件读取日志