---
title: Prometheus 规则处理逻辑
---

# 概述

> 参考：
> 
> - <https://mp.weixin.qq.com/s/fwHfKYCy_SKJzaiNy-zh7A>

# 接口

代码：`./rules/manager.go` —— `Rule{}`

Rule 接口封装了一个向量表达式，在指定的时间间隔评估规则。Prometheus 将规则分为两类：**Recording Rule(记录规则)** 与 **Alerting Rule(告警规则)**，所以将处理这两种规则的方法统一成一个接口，如下两个结构体实现了该接口：

- `./rules/alerting.go` —— `AlertingRule{}`
- `./rules/recording.go` —— `RecordingRule{}`

```go
type Rule interface {
    // 直接返回规则的名称
	Name() string
	// Labels of the rule.
	Labels() labels.Labels
    // 评估规则(规则处理逻辑中最重要的部分)
	Eval(context.Context, time.Time, QueryFunc, *url.URL) (promql.Vector, error)
	// String returns a human-readable string representation of the rule.
	String() string
	// Query returns the rule query expression.
	Query() parser.Expr
	// SetLastErr sets the current error experienced by the rule.
	SetLastError(error)
	// LastErr returns the last error experienced by the rule.
	LastError() error
	// SetHealth sets the current health of the rule.
	SetHealth(RuleHealth)
	// Health returns the current health of the rule.
	Health() RuleHealth
	SetEvaluationDuration(time.Duration)
	// GetEvaluationDuration returns last evaluation duration.
	// NOTE: Used dynamically by rules.html template.
	GetEvaluationDuration() time.Duration
	SetEvaluationTimestamp(time.Time)
	// GetEvaluationTimestamp returns last evaluation timestamp.
	// NOTE: Used dynamically by rules.html template.
	GetEvaluationTimestamp() time.Time
	// HTMLSnippet returns a human-readable string representation of the rule,
	// decorated with HTML elements for use the web frontend.
	HTMLSnippet(pathPrefix string) html_template.HTML
}
```

这些方法中，最重要的就是 `Eval()`，用来评估规则。

# 结构体

## 规则

代码：`./rules/manager.go` — `Group{}`

这是一组具有逻辑关系的规则。`Manager.LoadGroups()` 方法将会读取配置文件中的每组规则，解析后，返回一个 `map[string]*Group`(也就是说)，以供后续代码使用。

```go
type Group struct {
	name                 string
	file                 string
	interval             time.Duration
	limit                int
	rules                []Rule
	seriesInPreviousEval []map[string]labels.Labels // One per Rule.
	staleSeries          []labels.Labels
	opts                 *ManagerOptions
	mtx                  sync.Mutex
	evaluationTime       time.Duration
	lastEvaluation       time.Time

	shouldRestore bool

	markStale   bool
	done        chan struct{}
	terminated  chan struct{}
	managerDone chan struct{}

	logger log.Logger

	metrics *Metrics
}
```

## 规则管理器

代码：`./rules/manager.go` — `Manager{}`

Manager 结构体是规则管理器。负责管理 记录规则 与 告警规则。

```go
type Manager struct {
    // 管理器的额外选项，
	opts     *ManagerOptions
    // 通过配置文件传递进来的规则
	groups   map[string]*Group
    // 该结构体的读写锁，通常在处理 groups 时会上锁
	mtx      sync.RWMutex
	block    chan struct{}
	done     chan struct{}
	restored bool

	logger log.Logger
}
```

在 Prometheus Server 的 `main()` 中，通过 `NewManager()` 函数实例化一个管理器，并在实例化时传递已实例化的 `ManagerOptions{}`

## 管理器选项

`./rules/manager.go` —— `ManagerOptions{}`

在 `main()` 中实例化 `Manager{}` 时，会将该结构体作为实参传递到 `NewManager()` 函数中。

```go
type ManagerOptions struct {
	ExternalURL     *url.URL
    //
	QueryFunc       QueryFunc
	NotifyFunc      NotifyFunc
	Context         context.Context
	Appendable      storage.Appendable
	Queryable       storage.Queryable
	Logger          log.Logger
	Registerer      prometheus.Registerer
	OutageTolerance time.Duration
	ForGracePeriod  time.Duration
	ResendDelay     time.Duration
	GroupLoader     GroupLoader

	Metrics *Metrics
}
```

# 加载并运行规则

`./rules/manager.go` — `Manager.Updata()`

在 Prometheus Server 启动时会调用 `Manager.Updata()` 方法加载规则配置文件并解析。在配置热更新时。如果加载新规则失败，则将恢复旧规则。

这个更新规则是这么个逻辑：比较每一个规则组，若一样，则将老的复制成新的，然后再加上原来没有的，若不一样，则规则组处理完成后，删除这些剩余的。

这里面最重要部分有两个

- `Manager.LoadGroups()` 方法。从配置文件中读取内容，并解析规则，实例化一个规则组(即 `Group{}` 结构体)
- 在并发中执行的 `Group.run()` 方法。该方法中包含里定期评估规则的逻辑

```go
func (m *Manager) Update(interval time.Duration, files []string, externalLabels labels.Labels, externalURL string) error {
	// 下面将会处理 Manager.groups，所以上个锁~
    m.mtx.Lock()
	defer m.mtx.Unlock()
    // 从配置文件中加载并解析规则，以实例化一个规则组变量供后续处理
	groups, errs := m.LoadGroups(interval, externalLabels, externalURL, files...)
    // 准备开始并发喽~
	var wg sync.WaitGroup
	for _, newg := range groups {
        // 检查新组与旧组标识符是否相等，若相等则跳过，若不等则停止它并等待它完成当前迭代，然后将其复制到新组中
		......
        // 并发添加计数
        wg.Add(1)
        // 开始真正的加载规则组，其中主要执行逻辑在 newg.run(m.opts.Context)
		go func(newg *Group)
        	// 接着上面检查新旧组标识符，如果相等，这里就会将老规则组中的状态信息复制到新规则组中。
			if ok {
				oldg.stop()
				newg.CopyState(oldg)
			}
			wg.Done()
			// Wait with starting evaluation until the rule manager
			// is told to run. This is necessary to avoid running
			// queries against a bootstrapping storage.
			<-m.block
            // 加载新组，周期性得执行 PromQL 语句
			newg.run(m.opts.Context)
		}(newg)
	}

	// 删除余下的旧组
	wg.Add(len(m.groups))
	for n, oldg := range m.groups {
		go func(n string, g *Group) {
            ......
		}(n, oldg)
	}

	wg.Wait()
	m.groups = groups

	return nil
}
```

# 加载规则

`./rules/manager.go` —— `Manager.LoadGroups()`

在这里，会使用 `GroupOptions{}` 实例化 `Group{}`。

```go
func (m *Manager) LoadGroups(interval time.Duration, externalLabels labels.Labels, externalURL string, filenames ...string,) (map[string]*Group, []error) {
    // 循环配置文件
	for _, fn := range filenames {
        // 通过 ./rulefmt/rulefmt.go —— ParseFile() 解析配置文件，返回 rulefmt 包下的 RuleGroups{}，这就是规则组实例
		rgs, errs := m.opts.GroupLoader.Load(fn)
        // 循环规则组实例
		for _, rg := range rgs.Groups {
            // 循环规则组中的每个规则
			for _, r := range rg.Rules {
			}

			groups[GroupKey(fn, rg.Name)] = NewGroup(GroupOptions{
				Name:          rg.Name,
				File:          fn,
				Interval:      itv,
				Limit:         rg.Limit,
				Rules:         rules,
				ShouldRestore: shouldRestore,
				Opts:          m.opts,
				done:          m.done,
			})
		}
	}
    // 循环完成后，返回实例化的 Group{}
	return groups, nil
}
```

# 运行规则

`./rules/manager.go` — `Group.run()`

上面的 `Manager.LoadGroups()` 中获得了实例化的 Group{}，通过规则组，定期执行评估行为。这里面最重要的是定期执行的 `Group.Eval()` 方法。

```go
func (g *Group) run(ctx context.Context) {
	// 等待一个当前的时间戳，以开始计时，以便根据固定的间隔，持续评估
	evalTimestamp := g.EvalTimestamp(time.Now().UnixNano()).Add(g.interval)

    // ！！！！用闭包声明一个函数变量，这里面是执行评估的最主要代码逻辑！！！！！
	iter := func() {
		g.Eval(ctx, evalTimestamp)
	}
    // 在一个循环中，每隔一段时间就执行一次 Group.Eval()
	for {
		select {
		case <-g.done:
			return
		default:
			select {
			case <-g.done:
				return
			case <-tick.C:
				iter()
			}
		}
	}
}
```

## 评估规则

`./rules/manager.go` — `Group.Eval()`

规则评估，在 `Group.run()` 方法中定期调用本方法来评估所有规则组。

```go
func (g *Group) Eval(ctx context.Context, ts time.Time) {
	for i, rule := range g.rules {
		func(i int, rule Rule) {
            // 记录规则 与 报警规则 具有各自的实现。但是都会返回一个向量。
			vector, err := rule.Eval(ctx, ts, g.opts.QueryFunc, g.opts.ExternalURL, g.Limit())
            // 如果当前评估的是报警规则，则进行断言，
			if ar, ok := rule.(*AlertingRule); ok {
				ar.sendAlerts(ctx, ts, g.opts.ResendDelay, g.interval, g.opts.NotifyFunc)
			}
		}(i, rule)
	}
}
```

这里面最重要的是运行了 Rule 接口下的 `rule.Eval()` 方法，由于 Prometheus 将规则分为了两类(记录规则 与 告警规则)，所以这里用了一个结构，不同类型的规则，其中处理细节不太一样，但是都会在这里进行评估。

在调用 Rule 接口下的 Eval() 方法时，传递了 QueryFunc，在评估规则时，将会执行 PromQL

### 评估告警规则

`./rules/alerting.go` —— `AlertingRule.Eval()`

> 注意：这里面有一个常量 `const resolvedRetention = 15 * time.Minute`，这个时间是已触发的警报保存在内存中的时间，在这个保存时间内，Prometheus Server 将会持续发送该警报(即使该警报已解决，此时发送的警报状态为 Resolved)

评估报警规则表达式，然后创建 pending 状态的报警，若满足条件，则转变为 alerting 状态，或者删除过期的 pending 状态的报警。最后，返回一个向量

```go
func (r *AlertingRule) Eval(ctx context.Context, ts time.Time, query QueryFunc, externalURL *url.URL, limit int) (promql.Vector, error) {
    // 执行 PromQL 语句，若没结果则直接返回，继续等待 Group.run() 循环中的下一个元素。
	res, err := query(ctx, r.vector.String(), ts)
	// 若有结果，则给该告警规则上锁，并开始处理
	r.mtx.Lock()
	defer r.mtx.Unlock()

	// Create pending alerts for any new vector elements in the alert expression
	// or update the expression value for existing elements.
	resultFPs := map[uint64]struct{}{}

	var vec promql.Vector
	var alerts = make(map[uint64]*Alert, len(res))
    // 循环每一个 PromQL 获取到的样本，也就是循环每一个查询语句获得的结果
	for _, smpl := range res {
        // 提取样本中的信息，保存到 Alert{} 中。相当于实例化了 Alert{}
		alerts[h] = &Alert{
			Labels:      lbs,
			Annotations: annotations,
			ActiveAt:    ts,
			State:       StatePending,
			Value:       smpl.V,
		}
	}

	for h, a := range alerts {
		// 检查标签集是否以具有 alerting 状态
		if alert, ok := r.active[h]; ok && alert.State != StateInactive {
			alert.Value = a.Value
			alert.Annotations = a.Annotations
			continue
		}

		r.active[h] = a
	}

	// Check if any pending alerts should be removed or fire now. Write out alert timeseries.
	for fp, a := range r.active {
		if _, ok := resultFPs[fp]; !ok {
			// If the alert was previously firing, keep it around for a given
			// retention time so it is reported as resolved to the AlertManager.
			if a.State == StatePending || (!a.ResolvedAt.IsZero() && ts.Sub(a.ResolvedAt) > resolvedRetention) {
				delete(r.active, fp)
			}
			if a.State != StateInactive {
				a.State = StateInactive
				a.ResolvedAt = ts
			}
			continue
		}

		if a.State == StatePending && ts.Sub(a.ActiveAt) >= r.holdDuration {
			a.State = StateFiring
			a.FiredAt = ts
		}

		if r.restored {
			vec = append(vec, r.sample(a, ts))
			vec = append(vec, r.forStateSample(a, ts, float64(a.ActiveAt.Unix())))
		}
	}

	numActive := len(r.active)
	if limit != 0 && numActive > limit {
		r.active = map[uint64]*Alert{}
		return nil, errors.Errorf("exceeded limit of %d with %d alerts", limit, numActive)
	}

	return vec, nil
}
```

### 评估记录规则

`./rules/recording.go` —— `RecordingRule.Eval()`

评估记录规则，然后相应地覆盖指标名称和标签。

```go
func (rule *RecordingRule) Eval(ctx context.Context, ts time.Time, query QueryFunc, _ *url.URL, limit int) (promql.Vector, error) {
	vector, err := query(ctx, rule.vector.String(), ts)
	if err != nil {
		return nil, err
	}
	// Override the metric name and labels.
	for i := range vector {
		sample := &vector[i]

		lb := labels.NewBuilder(sample.Metric)

		lb.Set(labels.MetricName, rule.name)

		for _, l := range rule.labels {
			lb.Set(l.Name, l.Value)
		}

		sample.Metric = lb.Labels()
	}

	// Check that the rule does not produce identical metrics after applying
	// labels.
	if vector.ContainsSameLabelset() {
		return nil, fmt.Errorf("vector contains metrics with the same labelset after applying rule labels")
	}

	numSamples := len(vector)
	if limit != 0 && numSamples > limit {
		return nil, fmt.Errorf("exceeded limit %d with %d samples", limit, numSamples)
	}

	rule.SetHealth(HealthGood)
	rule.SetLastError(err)
	return vector, nil
}
```
