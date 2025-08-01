---
title: Alertmanager 扩展
linkTitle: Alertmanager 扩展
weight: 20
---

# 概述

> 参考：
>
> - 

Alertmanager 自带一个 UI 界面，可以用来查看报警和静默管理。但是告警发送目标、历史告警、etc. 个性化功能还比较缺，有很多项目可以补充这些能力。

[GitHub 项目，prymitive/karma](https://github.com/prymitive/karma)

- [公众号-k8s 技术圈，超漂亮的 Alertmanager 可视化面板 - karma](https://mp.weixin.qq.com/s/uHSlzuVBb51-qgX92pEnLQ)
- 比如报警历史记录等等

[GitHub 项目，feiyu563/PrometheusAlert](https://github.com/feiyu563/PrometheusAlert) # 可以提供更多的通知功能，将告警发送到各种地方。

- 利用 template.FuncMap 函数在 go tmpl 中加入了一些自定义函数，e.g. toUpper、etc.
- 该程序在使用 [Alertmanager 数据结构](/docs/6.可观测性/Metrics/Alertmanager/Alertmanager%20数据结构.md) 中的 webhook 推送的数据结构时，<font color="#ff0000">虽然设计了 struct，但是在模板中调用 struct 中的属性时，开头字母要是小写，这种模板跟很多扩展都不通用</font>。

https://github.com/opsre/WatchAlert # 多数据源监控告警引擎

[GitHub 项目，timonwong/prometheus-webhook-dingtalk](https://github.com/timonwong/prometheus-webhook-dingtalk) # 对接钉钉。提供 pr，增加过 feature

- 利用 template.FuncMap 函数在 go tmpl 中加入了一些自定义函数，e.g. upper、etc.

https://github.com/rea1shane/a2w # 对接企业微信

- 利用 template.FuncMap 函数在 go tmpl 中加入了一些自定义函数，e.g. timeFormat、etc.

---

[GitHub 项目，kubesphere/notification-manager](https://github.com/kubesphere/notification-manager) # kubesphere 出的，只有 k8s 的（2025-07-31 kubesphere 闭源）。

# notification-manager

> 参考：
>
> - [GitHub 项目，kubesphere/notification-manager](https://github.com/kubesphere/notification-manager)

功能测试

```bash
wget https://raw.githubusercontent.com/kubesphere/notification-manager/master/config/ci/alerts.json
curl -XPOST http://localhost:19093/api/v2/alerts -d @./alerts.json
```

## NotificationManager CRD

### 接收器与配置相关字段

**receivers**(OBJECT) #

- **globalReceiverSelector**(OBJECT) #
  - 该字段内容详见 [LabelSelector](/docs/10.云原生/Kubernetes/API%20Resource%20与%20Object/API%20参考/Common%20Definitions(通用定义)/LabelSelector.md)
- **tenantReceiverSelector**(OBJECT) #
  - 该字段内容详见 [LabelSelector](/docs/10.云原生/Kubernetes/API%20Resource%20与%20Object/API%20参考/Common%20Definitions(通用定义)/LabelSelector.md)
- **tenantKey**(STRING) #

示例:

```yaml
receivers:
  # 具有 type: global 标签的 Receiver 将会被设置为全局 Receiver
  globalReceiverSelector:
    matchLabels:
      type: global
  # 具有 type: tenant 标签的 Receiver 将会被设置为租户 Receiver
  tenantReceiverSelector:
    matchLabels:
      type: tenant
  # notification-manager 通过 tenantKey 的值识别 Receiver 的租户名称。i.e. 租户类型的 Receiver 通过 key 为 user 的标签值识别租户名称
  tenantKey: user
```

### 通知管理器的 Webhook 与 Dispatcher 相关字段

**args**(\[]TYPE) # 设定 NotificationManager Webhook 的启动参数。

**batchMaxSize**(INT) # 从缓存中获取数据时最大的告警数量。`默认值：100`

**batchMaxWait**(DURATION) # 从缓存中获取数据的等待时间。`默认值：1m`。即每隔一分钟获取一次数据

> batchMaxSize 与 batchMaxWait 说明：Notification-Manager 接收到的告警数据首先会被推送到缓存中，再从缓存中批量取出数据并行处理。所以可以通过 `batchMaxSize` 与 `batchMaxWait` 两个字段来配置每次从缓存中取出多少数据与时间间隔。详见 从[缓存中获取告警](#moaPC)的代码。所以我们会发现，每次 Notification-Manager 收到告警后，将会等待 1 分钟之后才会开始处理这些告警。

**routePolicy**(STRING) # 路由策略，定义将收到的告警信息路由给哪个 Receiver。`默认值：All`。

- All # 通知信息将会被路由到所有通过 Router 匹配到的 Receiver 上，并且同时路由到到默认的全局 Receiver
- RouterFirst # 通知信息在被路由到 Router 匹配到的 Receiver 上之后，不在路由给默认的全局 Receiver
- RouterOnly # 通知信息只会被路有道 Router 匹配到的 Receiver 上。

### 生成通知信息与组织通知信息相关字段

### 其他字段

## Router CRD

**alertSelector**(OBJECT) # 告警标签选择器。与 K8S 的 LabelSelector 的功能完全一样

- 该字段内容详见[ LabelSelector](/docs/10.云原生/2.3.Kubernetes%20 容器编排系统/1.API、Resource(资源)、Object(对象)/API%20 参考/Common%20Definitions(通用定义)/LabelSelector%20 详解.md 容器编排系统/1.API、Resource(资源)、Object(对象)/API 参考/Common Definitions(通用定义)/LabelSelector 详解.md)。注意一点：多个匹配条件之间的关键是 AND。如果想要使用 OR 的逻辑，以根据多个条件匹配多条告警，需要使用多个 Router，详见 [Issue #153](https://github.com/kubesphere/notification-manager/issues/153)

## Receiver CRD

## 通用定义

## 代码分析

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/fy8tkv/1652004437657-71dec4c2-c104-4c34-a43c-74ffb8d4545b.png)

### 入口与监听

程序开始主要做了如下几件事：

- 实例化告警存储器，用以缓存接收到的告警消息。告警存储器称之为 Provider。
- 带着告警存储器实例化 Webhook 并运行，用以接受告警消息后将告警缓存起来(执行 Provider.Push() 方法)
- 带着告警存储器实例化调度员并运行，用以获取缓存中的告警消息(执行 Provider.Pull() 方法)

Cache 默认为 Memory，在内存中存储各个地方推送过来的告警
代码：`cmd/notification-manager/main.go`

```go
var (
 storeType = kingpin.Flag(
  "store.type",
  "Type of store which used to cache the alerts",
 ).Default("memory").String()
)
func Main() int {
    // 实例化告警存储器，默认内存
    alerts := store.NewAlertStore(*storeType)

    // 带着存储器实例化一个 Webhook，并启动监听程序，默认监听在 19093
 webhook := wh.New(
  alerts,
    )
 srvCh := make(chan error, 1)
 go func() {
  srvCh <- webhook.Run(ctxHttp)
 }()

    // 带着告警存储器实例化一个 Dispatcher，用以从告警存储器中 pull 下来告警后发送出去
    disp := dispatcher.New(logger, ctl, alerts, *webhookTimeout, *wkrTimeout, *wkrQueue)
 go func() {
  dispCh <- disp.Run()
 }()
}
```

告警存储器都实现了 Provider 接口
代码：`pkg/store/provider/interface.go`

```go
type Provider interface {
 Push(alert *template.Alert) error
 Pull(batchSize int, batchWait time.Duration) ([]*template.Alert, error)
 Close() error
}
```

### 接收告警并推送到缓存

而想要 Pull 到数据，则需要先通过告警存储器中的 Provider Push 到存储中，首先通过 /api/v2/alerts 端点接收告警
代码：`pkg/webhook/webhook.go`

```go
func New(logger log.Logger, notifierCtl *controller.Controller, alerts *store.AlertStore, o *Options) *Webhook {
 h := &Webhook{
  Options: o,
  logger:  logger,
 }

 h.router.Post("/api/v2/alerts", h.handler.Alert)
}
```

通过 Provider.Push() 方法推送的告警将进入 Channel，由 Dispatcher 的通知阶段代码进行消费 Channel 中的告警信息以发送给 Receiver
代码：`pkg/webhook/v1/handler.go`

```go
func (h *HttpHandler) Alert(w http.ResponseWriter, r *http.Request) {
 data := template.Data{}

 if err := utils.JsonDecode(r.Body, &data); err != nil {
 }

 for _, alert := range data.Alerts {
        // 推送告警
  if err := h.alerts.Push(alert); err != nil {
   _ = level.Error(h.logger).Log("msg", "push alert error", "error", err.Error())
  }
 }

 h.handle(w, &response{http.StatusOK, "Notification request accepted"})
}

```

### 从缓存中获取告警

Dispatcher 中通过 Pull() 方法从存储中获取告警，并通过 Dispatcher.processAlerts() 方法处理他们以便发送。
代码：`pkg/dispatcher/dispatcher.go`

```go
func (d *Dispatcher) Run() error {

 for {
  // err is not nil means the store had closed, dispatcher should process remaining alerts, then exit.
        // BatchMaxSize 定义了每次从缓存中可以获取的最大告警条数
        // BatchMaxWait 定义了每次执行 Pull() 的间隔时间
        // 默认情况下，每隔 1 分钟会 PUll 100 条告警以进一步处理
  if alerts, err := d.alerts.Pull(d.notifierCtl.GetBatchMaxSize(), d.notifierCtl.GetBatchMaxWait()); err == nil {
   go d.processAlerts(alerts)
  } else {
   d.processAlerts(alerts)
   return nil
  }
 }
}
```

Dispatcher.processAlerts() -> Dispatcher.worker() 将会执行[告警处理阶段](#oxPq5)

### 执行告警处理阶段

代码：`pkg/dispatcher/dispatcher.go`

```go
func (d *Dispatcher) worker(ctx context.Context, data interface{}, stopCh chan struct{}) {
 pipeline := stage.MultiStage{}
 // Global silence stage
 pipeline = append(pipeline, silence.NewStage(d.notifierCtl))
 // Route stage
 pipeline = append(pipeline, route.NewStage(d.notifierCtl))
 // Tenant silence stage
 pipeline = append(pipeline, filter.NewStage(d.notifierCtl))
 // Aggregation stage
 pipeline = append(pipeline, aggregation.NewStage(d.notifierCtl))
 // Notify stage
 pipeline = append(pipeline, notify.NewStage(d.notifierCtl))
 // History stage
 pipeline = append(pipeline, history.NewStage(d.notifierCtl))

 if _, _, err := pipeline.Exec(ctx, d.l, data); err != nil {
 }

 stopCh <- struct{}{}
}
```

通过 MultiStage 按顺序执行一系列阶段，最后执行 MultiStage.Exec()，MultiStage 实现了 Stage 接口

```go
type Stage interface {
 Exec(ctx context.Context, l log.Logger, data interface{}) (context.Context, interface{}, error)
}
```

同时，所有对告警信息需要执行的操作(上图中 Cache 右边的部分)都实现了该接口：

```go
// 告警静音 pkg/silence/silence.go
type silenceStage struct {
 notifierCtl *controller.Controller
}
// 告警路由 pkg/route/router.go
type routeStage struct {
 notifierCtl *controller.Controller
}
// 告警过滤 pkg/filter/filter.go
type filterStage struct {
 notifierCtl *controller.Controller
}
// 告警聚合 pkg/aggregation/aggregation.go
type aggregationStage struct {
 notifierCtl *controller.Controller
}
// 告警通知 pkg/notify/notify.go
type notifyStage struct {
 notifierCtl *controller.Controller
}
// 告警历史 pkg/history/history.go
type historyStage struct {
 notifierCtl *controller.Controller
}
```

告警的每个处理阶段，均由上述操作的 Exec() 方法实现

#### 告警通知阶段

代码：

```go
func (s *notifyStage) Exec(ctx context.Context, l log.Logger, data interface{}) (context.Context, interface{}, error) {

 if reflect2.IsNil(data) {
  return ctx, nil, nil
 }

 _ = level.Debug(l).Log("msg", "Start notify stage", "seq", ctx.Value("seq"))

 group := async.NewGroup(ctx)

    // Receiver 是告警的接受者，即推送目标
    // []*template.Data 是需要推送的告警列表
 alertMap := data.(map[internal.Receiver][]*template.Data)

 for k, v := range alertMap {
  receiver := k
  ds := v
        // 获取推送目标，比如 钉钉、微信 等
  nf, err := factories[receiver.GetType()](l, receiver, s.notifierCtl)

        //
  for _, d := range ds {
   alert := d
   group.Add(func(stopCh chan interface{}) {
                // 使用对应的 Receiver 的 Notify() 方法发送通知
    stopCh <- nf.Notify(ctx, alert)
   })
  }

 }

 return ctx, data, group.Wait()
}
```

所有 Receiver 都实现了 Notifier 接口

代码：`pkg/notify/notifier/interface.go`

```go
type Notifier interface {
 Notify(ctx context.Context, data *template.Data) error
}
```

代码：`pkg/notify/notifier/${RECEIVER}/${RECEIVER}.go`

以 钉钉(dingtalk) 为例

代码：`pkg/notify/notifier/dingtalk/dingtalk.go`

```go
func (n *Notifier) Notify(ctx context.Context, data *template.Data) error {

 group := async.NewGroup(ctx)
 if n.receiver.ChatBot != nil {
  group.Add(func(stopCh chan interface{}) {
   stopCh <- n.sendToChatBot(ctx, data)
  })
 }

 if len(n.receiver.ChatIDs) > 0 {
  group.Add(func(stopCh chan interface{}) {
   stopCh <- n.sendToConversation(ctx, data)
  })
 }

 return group.Wait()
}
```
