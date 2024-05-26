---
title: Instrumenting 原理解析
---

# 概述

> 参考：
> - 根据源码一步一步推到自学
> - [prometheus 默认自带的 Metrics 的实现方式](https://github.com/prometheus/client_golang/blob/master/prometheus/go_collector.go)
> - [prometheus 库](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus)
> - [prometheus/promhttp 库](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus/promhttp)

Instrumenting 的实现主要依靠以下几种类型：

- **Desc(描述符)** # 结构体。定义一个 Metric
- **Registerer(注册器)** # 接口。根据 Metrics 注册一个 Collector(采集器)
- **Collector(采集器)** # 接口。采集 Metrics 的具体实现
- **Gatherer(聚集器)** # 接口。将采集到的 Metrics 聚集在一起

其中 Collector(采集器) 就像其名字一样，是定义采集 Metrics 的主要行为。在代码中，Collector(采集器) 表现为一个接口。这个接口有两个方法，`Describe()` 与 `Collect()`，其中在 `**Collect()**`** 这个方法中，定义主要的采集 Metrics 行为**

# Desc(描述符) - 用来描述 Metric 的结构体

https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#Desc

在 Prometheus 中，使用 **Desc 结构体** 来描述一个 Metric。Desc 是所有事物的基础，没有 Desc 也就无从采集 Metric，同时管理 Metric 也是通过 Desc

```go
type Desc struct {
	// 完全限定名称。也就是 Metric 的名字，fqName 由 Namespace、Subsystem、Name 三部分组成
	fqName string
	// Metric 的帮助信息
	help string
	// constLabelPairs(常量标签对) 包含基于常量标签的预先计算的 DTO标签对。
	constLabelPairs []*dto.LabelPair
	// variableLabels 包含 metrics 为其维护变量值的标签名称。
	variableLabels []string
	// id 是 ConstLabels 与 fqName 两个值的 hash。
	// 它在所有已注册的 Desc(描述符) 中必须是唯一的，因此可以用作 Desc(描述符) 的 ID(标识符)。
	id uint64
	// dimHash 是标签名称(预设和变量) 和 Help 这两个值的 hash。
	// 具有相同 fqName 的每个 Desc 必须具有相同的 dimHash
	dimHash uint64
	// err 是在每次构建过程中发生的错误。这个错误信息会报告注册 Desc 的时间
	err error
}
```

Desc 是每个 **Metric 所使用的 Descriptor(描述符)**。Desc 本质上是其对应的 Metric 的不可变的元数据。prometheus 包中那些常规 Metric 的实现，是用来在底层管理其对应的 Desc。只需要处理 Desc 即可使用诸 ExpvarCollector 或者 CustomCollectors 和 Metrics 之类的高级功能。

> 这个描述符就类似于 Linux 中的 FD(文件描述符) 的作用。我打开一个文件，就可以通过该文件的 FD 对该文件进行读写。Prometheus 的 Desc 同理，实例化这个 Desc 结构体，就等于是打开了一个 Metric，可以通过 Metric 的 描述符 来对该 Metric 进行读写操作。

Descriptors registered with the same registry have to fulfill certain consistency and uniqueness criteria if they share the same fully-qualified name: They must have the same help string and the same label names (aka label dimensions) in each, constLabels and variableLabels, but they must differ in the values of the constLabels.

Descriptors that share the same fully-qualified names and the same label values of their constLabels are considered equal.

使用 `NewDesc()` 函数来创建一个新的 Desc 实例。

## Desc 结构体中的属性(也就是一个 Metric 的元数据)详解

参考：[源码注释](https://github.com/prometheus/client_golang/blob/master/prometheus/desc.go#L46)

> 注意，源码注释参考是 GitHub 上的，不是 pkg.go 中的，GitHub 上代码所在的行很可能会变，可以通过 这里 直接点击 Desc 的连接跳转过去

在源码中，我们可以看到有如下 7 个属性

1. fqName string
2. help string
3. constLabelPairs \[]\*dto.LabelPair
4. variableLabels \[]string
5. id uint64
6. dimHash uint64
7. err error

假如现在有这么一个指标：

    # HELP go_info Information about the Go environment.
    # TYPE go_info gauge
    go_info{version="go1.15.5"} 1

那么在 Prometheus 代码中 Metric 描述符对应的信息应该这么看：

fqName、help、 constLabelPairs、variableLabels 这四个属性的值，将会响应给客户端

**fqName **# 该 Metric 的名字

> 就是上述指标中的 go_info

fqName 由 Namespace、Subsystem、Name 三部分组成

**help string** # 该 Metric 的帮助信息

> 就是上述指标中的 Information about the Go environment

**constLabelPairs \[]\*dto.LabelPair** # 该 Metric 中，标签值不可变的标签对列表

> 在这个 Metirc 中没有 constLabelParis。常量标签对是在实例化 Desc 时设置的，标签值是不可变的。
> LabelPair(标签对) 是一组键值对的组合。key 为 标签名称，value 为 标签值

constLabelPairs\*\* \*\*包含基于常量标签的预先计算的 DTO 标签对。

constLabels 属性中的标签值是常量，不变的。因为可以在 Desc 中直接设置标签对的信息。

这些不可变的指标由于是在代码中直接设置的，所以常出现在 histogram 类型的指标中，比如 etcd 中的 etcd_debugging_mvcc_db_compaction_total_duration_milliseconds 指标中的 le 标签的值就是一个常量，并不会随所在环境而改变。

    # HELP etcd_debugging_mvcc_db_compaction_total_duration_milliseconds Bucketed histogram of db compaction total duration.
    # TYPE etcd_debugging_mvcc_db_compaction_total_duration_milliseconds histogram
    etcd_debugging_mvcc_db_compaction_total_duration_milliseconds_bucket{le="100"} 11331
    etcd_debugging_mvcc_db_compaction_total_duration_milliseconds_bucket{le="200"} 17071
    ....

**variableLabels \[]string** # 该 Metric 中，标签值可变的标签名称列表

> 就是上述指标中的 version="go1.15.5" 中的 version，表示该 Metric 的标签集的 key，这个标签的值是可变的。

variableLabels\*\* \*\*包含 [Metric 接口](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?utm_source=gopls#Metric) 所维护标签值的标签名称。

variableLabels\*\* \*\*所定义的标签名称对应的标签值是可变的，因此这些值不是 Desc 中的一个属性(它们通过 [Metric 接口](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?utm_source=gopls#Metric) 管理)。

比如 node_exporter 中采集到的 node_ipvs_backend_connections_active 这个指标，其中的的所有标签都是 variableLabels，标签值是采集器根据当前环境设定的。就像这个 Metric 中的 address、port 等等，都是可变的。

    # HELP node_ipvs_backend_connections_active The current active connections by local and remote address.
    # TYPE node_ipvs_backend_connections_active gauge
    node_ipvs_backend_connections_active{local_address="10.100.121.107",local_mark="",local_port="9098",proto="TCP",remote_address="10.244.5.209",remote_port="9098"} 0
    node_ipvs_backend_connections_active{local_address="10.100.180.246",local_mark="",local_port="8500",proto="TCP",remote_address="10.244.3.23",remote_port="8500"} 0
    ...

**id uint64** # 不会响应给 scrape 请求

**dimHash uint64** # 不会响应给 scrape 请求

**err error** # 不会响应给 scrape 请求

## [NewDesc()](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#NewDesc)

参考：[代码注释](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?utm_source=gopls#NewDesc)
`NewDesc()` 用来实例化 Desc 结构体，设置属性的值并初始化一个新的 Desc。

实例化时产生的错误也会被记录在 Desc 中，并在注册该 Desc 时报告。

在实例化时，variableLabels 和 constLabels 可以为 nil，fqName 不能为空，`NewDesc()` 将会为其他未明确设置的属性设置默认值。

### 高级 NewDesc()

看完后面的章节再回来看这个，不理解 Collector、Registerer、Gatherer 的知识，无法理解这部分内容。参见[后文](#HNvsE)

## Exposed Metric(暴露指标)

所谓 Exposed Metrics 就是让指标可以展示在 http 服务器上。也就是说通过 Exporter 监听的端口来获取 Metrics 的信息。

想要暴露 Metrics，那么必须要先注册这些 Metrics，注册 Metrics 是通过 Collector 进行的。一般使用 prometheus.MustRegister() 注册一个 Metric。

注册完成之后想要参见下文 Exposed Mestrics 来让 metrics 可访问

## 基本示例

```go
// Metrics 用来保存所有 Metrics
type Metrics struct {
	HelloWorldDesc *prometheus.Desc
	mutex          sync.Mutex
}

// NewMetrics 实例化所有的 Metrics,并为 Mestirs 设定一些基本信息
func NewMetrics() *Metrics {
	return &Metrics{
		HelloWorldDesc: prometheus.NewDesc(
			"exporter_hello_world",               // Metric 名称
			"Help Info for exporter hello world", // Metric 的帮助信息
			[]string{"name"}, nil,                // Metric 的可变标签值的标签 与 不可变标签值的标签
		),
	}
}
```

# Collector(采集器)

https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#Collector

在上文 Metric 部分，我们定义个了一个 Metric，下一步就是要为这个 Metric 采集值、设置标签值、设置值类型。此时就要使用一个名为 Collector 的接口，只要让这个 Metric 实现该接口，即可实现采集行为。

**Collector(采集器**) 是由 Prometheus 可以用来采集 Metrics 的任何东西实现的接口。 一个 Collector(采集器) 必须使用 [Registerer.Register()](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?utm_source=gopls#Registerer) 进行注册后，才可以进行采集 Metric 的工作。

> Registerer.Register() 指的是 [Registerer](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?utm_source=gopls#Registerer) 接口中的 Register() 方法

也就是说，任何实现了 Collector 接口的代码，都可以通过 Registerer 注册进来，并在监听的端口上暴露出采集到的 Metrics 数据。

Collector 包含两个方法：

- **Describe()** # Metric 的描述符
- **Collect()** # 采集 Metric 的真正行为

可以看到，Collector 可以算包含了 Metrics，也就是说，注册了 Collector 也就等于相当于注册了 Metrics。

## Describe()

Describe sends the super-set of all possible descriptors of metrics collected by this Collector to the provided channel and returns once the last descriptor has been sent. The sent descriptors fulfill the consistency and uniqueness requirements described in the Desc documentation.

It is valid if one and the same Collector sends duplicate descriptors. Those duplicates are simply ignored. However, two different Collectors must not send duplicate descriptors.

Sending no descriptor at all marks the Collector as “unchecked”, i.e. no checks will be performed at registration time, and the Collector may yield any Metric it sees fit in its Collect method.

This method idempotently sends the same descriptors throughout the lifetime of the Collector. It may be called concurrently and therefore must be implemented in a concurrency safe way.

If a Collector encounters an error while executing this method, it must send an invalid descriptor (created with NewInvalidDesc) to signal the error to the registry.

Describe 将此收集器收集的指标的所有可能描述符的超集发送到提供的通道，并在发送完最后一个描述符后返回。发送的描述符满足 Desc 文档中描述的一致性和唯一性要求。

如果一个收集器和同一收集器发送重复的描述符，则该方法有效。这些重复项将被忽略。但是，两个不同的收集器不得发送重复的描述符。

完全不发送任何描述符会将收集器标记为“未检查”，即在注册时将不执行任何检查，并且收集器可以在其 Collect 方法中产生它认为合适的任何度量标准。

此方法在收集器的整个生命周期中均等地发送相同的描述符。它可以被同时调用，因此必须以并发安全的方式实现。

如果收集器在执行此方法时遇到错误，则它必须发送一个无效的描述符（使用 NewInvalidDesc 创建）以向注册表发送错误信号

## Collect()

Collect is called by the Prometheus registry when collecting metrics. The implementation sends each collected metric via the provided channel and returns once the last metric has been sent. The descriptor of each sent metric is one of those returned by Describe (unless the Collector is unchecked, see above). Returned metrics that share the same descriptor must differ in their variable label values.

This method may be called concurrently and must therefore be implemented in a concurrency safe way. Blocking occurs at the expense of total performance of rendering all registered metrics. Ideally, Collector implementations support concurrent readers.

收集指标时，Prometheus 注册表会调用“收集”。 该实现通过提供的通道发送每个收集的度量，并在发送完最后一个度量后返回。 每个发送的指标的描述符都是 Describe 返回的指标之一（除非未选中收集器，请参见上文）。 共享相同描述符的返回指标必须在其变量标签值上有所不同。

可以同时调用此方法，因此必须以并发安全的方式实现。 发生阻塞会以呈现所有已注册指标的总体性能为代价。 理想情况下，收集器实现支持并发读取器。

## 基本示例

```go
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
}
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
}
```

# Registerer(注册器)

https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#Registerer

Registerer is the interface for the part of a registry in charge of registering and unregistering. Users of custom registries should use Registerer as type for registration purposes (rather than the Registry type directly). In that way, they are free to use custom Registerer implementation(e.g. for testing purposes).

## Register()

Register registers a new Collector to be included in metrics collection. It returns an error if the descriptors provided by the Collector are invalid or if they — in combination with descriptors of already registered Collectors — do not fulfill the consistency and uniqueness criteria described in the documentation of metric.Desc.

If the provided Collector is equal to a Collector already registered (which includes the case of re-registering the same Collector), the returned error is an instance of AlreadyRegisteredError, which contains the previously registered Collector.

A Collector whose Describe method does not yield any Desc is treated as unchecked. Registration will always succeed. No check for re-registering (see previous paragraph) is performed. Thus, the caller is responsible for not double-registering the same unchecked Collector, and for providing a Collector that will not cause inconsistent metrics on collection. (This would lead to scrape errors.)

## MustRegister()

`MustRegister()` 的工作方式与 `Register()` 相同，只不过可以注册**任意数量**的**采集器**，注册过程产生的任何错误都会直接 panic。代码逻辑非常简单，就是下面这样：

```go
// 其实就是一个循环，逐一为每个 Collector 执行 Register()
func (r *Registry) MustRegister(cs ...Collector) {
	for _, c := range cs {
		if err := r.Register(c); err != nil {
			panic(err)
		}
	}
}
```

## Unregister()

Unregister unregisters the Collector that equals the Collector passed in as an argument. (Two Collectors are considered equal if their Describe method yields the same set of descriptors.) The function returns whether a Collector was unregistered. Note that an unchecked Collector cannot be unregistered (as its Describe method does not yield any descriptor).

Note that even after unregistering, it will not be possible to register a new Collector that is inconsistent with the unregistered Collector, e.g. a Collector collecting metrics with the same name but a different help string. The rationale here is that the same registry instance must only collect consistent metrics throughout its lifetime.

# Gatherer(聚集器)

https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#Gatherer

参考：[Gatherer 接口代码注释](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?utm_source=gopls#Gatherer)

Gatherer is the interface for the part of a registry in charge of gathering the collected metrics into a number of MetricFamilies. The Gatherer interface comes with the same general implication as described for the Registerer interface.

## Gather()

Gather calls the Collect method of the registered Collectors and then gathers the collected metrics into a lexicographically sorted slice of uniquely named MetricFamily protobufs. Gather ensures that the returned slice is valid and self-consistent so that it can be used for valid exposition. As an exception to the strict consistency requirements described for metric.Desc, Gather will tolerate different sets of label names for metrics of the same metric family.

Even if an error occurs, Gather attempts to gather as many metrics as possible. Hence, if a non-nil error is returned, the returned MetricFamily slice could be nil (in case of a fatal error that prevented any meaningful metric collection) or contain a number of MetricFamily protobufs, some of which might be incomplete, and some might be missing altogether. The returned error (which might be a MultiError) explains the details. Note that this is mostly useful for debugging purposes. If the gathered protobufs are to be used for exposition in actual monitoring, it is almost always better to not expose an incomplete result and instead disregard the returned MetricFamily protobufs in case the returned error is non-nil.

# 其他

## Exposed Metric(暴露指标)

Prometheus 暴露指标是通过 [promhttp.HandleFor()](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus/promhttp#HandlerFor) 函数实现的

## 高级版 NewDesc()

除了 NewDesc() 这种基本实例化 Desc 的方式，还有更高级的，可以直接定义指定类型的 Metric，并且通过这种方式实例化的 Desc 可以直接被 prometheus.MustRegister() 函数注册，并且相关的 struct 都自带了一些采集行为相关的方法，只需要再自己实现具体逻辑即可。

比如

```go
// gauge 结构体实现了 gauge 接口
type gauge struct {
	// valBits contains the bits of the represented float64 value. It has
	// to go first in the struct to guarantee alignment for atomic
	// operations.  http://golang.org/pkg/sync/atomic/#pkg-note-BUG
	valBits uint64
	selfCollector
	desc       *Desc
	labelPairs []*dto.LabelPair
}
```

gauge 类型的 Metric 结构体，包含了 \*Desc 的属性。并且 NewGauge() 实例化就是调用的 NewDesc()

> NewGauge() 返回的是一个 Gauge 的接口

```go
func NewGauge(opts GaugeOpts) Gauge {
	desc := NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		nil,
		opts.ConstLabels,
	)
	result := &gauge{desc: desc, labelPairs: desc.constLabelPairs}
	result.init(result) // Init self-collection.
	return result
}
```

### 应用示例

```go
var	cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{Name: "a_cpu_temperature_celsius",Help: "Current temperature of the CPU."})
prometheus.MustRegister(cpuTemp)
cpuTemp.Set(65.3)
```

这里通过 prometheus.NewGauge() 直接定义了一个 Metric，并通过 prometheus.MustRegister(cpuTemp) 注册这个 Metirc，并且可以通过现成的 Set() 方法为该 Metric 设定一个值。

注意：prometheus.MustRegister() 可以同时注册多个 Metircs：

```go
func MustRegister(cs ...Collector) {
	DefaultRegisterer.MustRegister(cs...)
}
```

在 prometheus 包中，有如下几个高级的 NewDesc()(待补充)：

- [NewCounter()](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?utm_source=gopls#NewCounter)
- [NewCounterVec()](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?utm_source=gopls#NewCounterVec)
- [NewGauge()](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?utm_source=gopls#NewGauge)
- [NewGaugeVec()](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?utm_source=gopls#NewGaugeVec)
- [NewHistogram()](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?utm_source=gopls#NewHistogram)
- [NewHistogramVec()](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?utm_source=gopls#HistogramVec)
- [NewSummary()](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?utm_source=gopls#NewSummary)
- [NewSummaryVec()](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?utm_source=gopls#SummaryVec)
- 。。。等等

带 Vec 的是可以为 Metric 设置标签的函数
