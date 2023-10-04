---
title: Exporter 开发实践
---

# 概述

> 参考：
>
> - [GitHub 自学代码](https://github.com/DesistDaydream/prometheus-instrumenting)
> - [默认自带的 Metrics 的实现方式](https://github.com/prometheus/client_golang/blob/master/prometheus/go_collector.go)
> - [公众号,k8s 技术圈-使用 Go 开发 Prometheus Exporter](https://mp.weixin.qq.com/s/s1nSaC-8ejvM342v5KPdxA)

在 [Instrumenting 原理解析](/docs/6.可观测性/监控系统/Prometheus/Prometheus%20开发/Instrumenting%20原理解析/Instrumenting%20原理解析.md) 中，逐一了解了实现 Exporter 的方法

- 首先，定义了一个包含 Metrics 描述符的结构体。以及实例化结构体的函数(也就是自定义一些 Metrics 的基本信息)
- 然后让该 结构体 实现 Collector 接口(i.e.为这个结构体添加 `Describe()` 与 `Collect()` 方法)
- 该 结构体 实现了 Collector 之后，就需要注册该 Metric，注册之后即可让 Prometheus 库通过 Collector 接口直接操作这个 Metric
- 而想要注册，首先需要一个新的注册器
- 创建完新的注册器之后，即可使用该注册器，将实现了 Collector 的 Metric 注册给 Prometheus 库。
- 最后，使用 HandlerFor() 将注册器作为参数传递进去，并返回一个 http.Handler，指定 访问路径，并设置监听端口
- 启动后，通过指定的访问路径，请求将会进入到 返回的 http.Handler 中，开始执行代码，最后获取完 Metric 信息，再响应给客户端

现在我将前面学习过程中零散的代码合并起来

```go
package main
import (
 "math/rand"
 "net/http"
 "sync"
 "github.com/prometheus/client_golang/prometheus"
 "github.com/prometheus/client_golang/prometheus/promhttp"
)
// HelloWorldMetrics 用来保存所有 Metrics
type HelloWorldMetrics struct {
 HelloWorldDesc *prometheus.Desc
 mutex          sync.Mutex
 // 加锁用，与主要处理逻辑无关
}
// NewHelloWorldMetrics 实例化 HelloWorldMetrics，就是为所有 Mestirs 设定一些基本信息
func NewHelloWorldMetrics() *HelloWorldMetrics {
 return &HelloWorldMetrics{
  HelloWorldDesc: prometheus.NewDesc(
   "a_hello_world_exporter",              // Metric 名称
   "Help Info for Hello World Exporter ", // Metric 的帮助信息
   []string{"name"}, nil,                 // Metric 的可变标签值的标签 与 不可变标签值的标签
  ),
 }
}
// Describe 让 HelloWorldMetrics 实现 Collector 接口。将 Metrics 的描述符传到 channel 中
func (ms *HelloWorldMetrics) Describe(ch chan<- *prometheus.Desc) {
 ch <- ms.HelloWorldDesc
}
// Collect 让 HelloWorldMetrics 实现 Collector 接口。采集 Metrics 的具体行为并设置 Metrics 的值类型,将 Metrics 的信息传到 channel 中
func (ms *HelloWorldMetrics) Collect(ch chan<- prometheus.Metric) {
 ms.mutex.Lock() // 加锁
 defer ms.mutex.Unlock()
 // 为 ms.HelloWorldDesc 这个 Metric 设置其属性的值
 // 该 Metric 值的类型为 Gauge，name 标签值为 haohao 时，Metric 的值为 1000 以内的随机数
 ch <- prometheus.MustNewConstMetric(ms.HelloWorldDesc, prometheus.GaugeValue, float64(rand.Int31n(1000)), "haohao")
 // 该 Metric 值的类型为 Gauge，name 标签值为 nana 时，Metric 的值为 100 以内的随机数
 ch <- prometheus.MustNewConstMetric(ms.HelloWorldDesc, prometheus.GaugeValue, float64(rand.Int31n(100)), "nana")
}
func main() {
 // 实例化自己定义的所有 Metrics
 m := NewHelloWorldMetrics()
 // 两种注册 Metrics 的方式
 //
 // 第一种：实例化一个新注册器，用来注册 自定义Metrics
 // 使用 HandlerFor 将自己定义的已注册的 Metrics 作为参数，以便可以通过 http 获取 metric 信息
 reg := prometheus.NewRegistry()
 reg.MustRegister(m)
 http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
 //
 // 第二种：使用自带的默认注册器，用来注册 自定义Metrics
 // prometheus.MustRegister(m)
 // http.Handle("/metrics", promhttp.Handler())
 // 让该 exporter 监听在8080 上
 http.ListenAndServe(":8080", nil)
}
/*
Export 暴露结果：
# HELP a_hello_world_exporter Help Info for Hello World Exporter
# TYPE a_hello_world_exporter gauge
a_hello_world_exporter{name="haohao"} 81
a_hello_world_exporter{name="nana"} 87
*/
```

# Prometheus 自带的 Metrics

```go
package main
import (
 "log"
 "net/http"
 "github.com/prometheus/client_golang/prometheus/promhttp"
)
func main() {
 http.Handle("/metrics", promhttp.Handler())
 log.Fatal(http.ListenAndServe(":8080", nil))
}
```

这个最简单的 Exporter 内部其实是使用了一个 **prometheus 库** 默认的采集器 `NewGoCollector()` 和 `NewProcessCollector()` 采集当前 Go 运行时的相关信息，比如 go 堆栈使用、goroutine 数据、当前程序所用资源等等。内容如下：

    root@lichenhao:~# curl localhost:8080/metrics
    # HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
    # TYPE go_gc_duration_seconds summary
    go_gc_duration_seconds{quantile="0"} 0
    go_gc_duration_seconds{quantile="0.25"} 0
    go_gc_duration_seconds{quantile="0.5"} 0
    go_gc_duration_seconds{quantile="0.75"} 0
    go_gc_duration_seconds{quantile="1"} 0
    go_gc_duration_seconds_sum 0
    go_gc_duration_seconds_count 0
    # HELP go_goroutines Number of goroutines that currently exist.
    # TYPE go_goroutines gauge
    go_goroutines 7
    # HELP go_info Information about the Go environment.
    # TYPE go_info gauge
    go_info{version="go1.15.5"} 1
    # ........后面略，还有很多

## promhttp.Handler()

`promhttp.Handler()` 使用默认的 [promhttp.HandlerOpts](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus/promhttp?utm_source=gopls#HandlerOpts) 为 [prometheus.DefaultGatherer](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?utm_source=gopls#pkg-variables) 返回一个[http.Handler](https://pkg.go.dev/net/http#Handler)。

    func Handler() http.Handler {
     return InstrumentMetricHandler(
      prometheus.DefaultRegisterer, HandlerFor(prometheus.DefaultGatherer, HandlerOpts{}),
     )
    }

`promhttp.Handler()` 会注册默认采集器并采集一些默认指标。

DefaultGatherer 是 [Gatherer](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?utm_source=gopls#Gatherer) 接口的实现；DefaultRegisterer 是 [Registerer](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?utm_source=gopls#Registerer)接口的实现。最初这两个[变量](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?utm_source=gopls#pkg-variables)都指向同一个 Registry，这个默认的 Registry 中包含 `NewGoCollector()` 和 `NewProcessCollector()` 这两个采集器

```go
var (
 defaultRegistry              = NewRegistry()
 DefaultRegisterer Registerer = defaultRegistry
 DefaultGatherer   Gatherer   = defaultRegistry
)
func init() {
 MustRegister(NewProcessCollector(ProcessCollectorOpts{}))
 MustRegister(NewGoCollector())
}
```

## InstrumentMetricHandler()

在 InstrumentMetricHandler() 函数中，使用了比 NewDesc() 更高级的 NewConterVec()、NewGauge() 来定义一个 Metric

InstrumentMetricHandler is usually used with an http.Handler returned by the HandlerFor function. It instruments the provided http.Handler with two metrics: A counter vector "promhttp_metric_handler_requests_total" to count scrapes partitioned by HTTP status code, and a gauge "promhttp_metric_handler_requests_in_flight" to track the number of simultaneous scrapes. This function idempotently registers collectors for both metrics with the provided Registerer. It panics if the registration fails. The provided metrics are useful to see how many scrapes hit the monitored target (which could be from different Prometheus servers or other scrapers), and how often they overlap (which would result in more than one scrape in flight at the same time). Note that the scrapes-in-flight gauge will contain the scrape by which it is exposed, while the scrape counter will only get incremented after the scrape is complete (as only then the status code is known). For tracking scrape durations, use the "scrape_duration_seconds" gauge created by the Prometheus server upon each scrape.

```go
func InstrumentMetricHandler(reg prometheus.Registerer, handler http.Handler) http.Handler {
 cnt := prometheus.NewCounterVec(
  prometheus.CounterOpts{
   Name: "promhttp_metric_handler_requests_total",
   Help: "Total number of scrapes by HTTP status code.",
  },
  []string{"code"},
 )
 // Initialize the most likely HTTP status codes.
 cnt.WithLabelValues("200")
 cnt.WithLabelValues("500")
 cnt.WithLabelValues("503")
 if err := reg.Register(cnt); err != nil {
  if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
   cnt = are.ExistingCollector.(*prometheus.CounterVec)
  } else {
   panic(err)
  }
 }
 gge := prometheus.NewGauge(prometheus.GaugeOpts{
  Name: "promhttp_metric_handler_requests_in_flight",
  Help: "Current number of scrapes being served.",
 })
 if err := reg.Register(gge); err != nil {
  if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
   gge = are.ExistingCollector.(prometheus.Gauge)
  } else {
   panic(err)
  }
 }
 return InstrumentHandlerCounter(cnt, InstrumentHandlerInFlight(gge, handler))
}
```

    registry.MustRegister(metrics)

1
Plain Text

这个 Collector 接口又代表什么意思呢

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ig8l2r/1616068562616-5c10af60-c810-4622-bc0f-35d331cbd2b0.png)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ig8l2r/1616068562665-009ae48a-3f65-48d1-b08f-49d7ae2089ca.png)

# 添加新的 Gatherer(采集器)
