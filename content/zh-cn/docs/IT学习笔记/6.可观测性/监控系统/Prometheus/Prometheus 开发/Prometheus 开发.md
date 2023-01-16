---
title: Prometheus 开发
---

# 概述

> 参考：
> - [GitHub 组织](https://github.com/prometheus)

# Prometheus 源码目录结构

> 更新日期：
> Prometheus 的源码目录随着更新迭代，也在不断变化中

cmd/ #&#x20;
config/ # 用于处理 yaml 格式的配置文件，包含与配置文件对应内容的 struct。
console_libraries/
consoles/
discovery/
docs/
documentation/
notifier/
pkg/
prompb/
promql/
rules/
scrape/
scripts/
storage/
template/
tsdb/
util/
web/

# Prometheus 通用包

> 参考：
> - [GitHub,prometheus/common](https://github.com/prometheus/common)

该存储库包含在 Prometheus 组件和库之间共享的 Go 库。它们被认为是 Prometheus 内部的，外部使用没有任何稳定性保证。

- **config** : 通用配置文件对应的 struct
  - 很多 Prometheus 的周边都会使用该包中的内容比如 blackbox-exporter 的配置文件中，就引用了该包中的 HTTPClientConfig 结构体，作为配置文件内容的一部分。
- **expfmt** : 展示格式的解码和编码
- **model**：共享数据结构
- **promlog** : [go-kit/log 的](https://github.com/go-kit/kit/tree/master/log)日志包装器
- **route**：使用 [httprouter](https://github.com/julienschmidt/httprouter) 的路由包装器 context.Context
- **server**：普通服务器
- **version**：版本信息和指标

# Go 客户端库

> 参考：
> - [prometheus/client_golang](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?utm_source=gopls#section-documentation)
> - <https://github.com/SongLee24/prometheus-exporter>

prometheus 包是一个 Core Instrumentation(核心仪器化) 的包。这个包为监控提供了用来仪器代码的 metrics 原语
