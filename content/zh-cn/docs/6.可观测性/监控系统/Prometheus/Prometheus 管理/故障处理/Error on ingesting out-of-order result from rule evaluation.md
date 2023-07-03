---
title: Error on ingesting out-of-order result from rule evaluation
---

# 概述

Error on ingesting out-of-order result from rule evaluation

该问题通常是由于记录规则的处理结果中，包含 NaN 而产生的告警，所有 NaN 的时间序列都会被丢弃，并不会保存到数据库中。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/bprr89/1633938701153-3857ea39-4849-4d33-89e3-ad34ac5313e0.png)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/bprr89/1633938707935-150d51c1-9dec-41a9-9346-f2e62bf74a53.png)

下面是报错的具体内容，可以看到 numDropped 是记录规则查询后生成的新时间序列中，被丢弃的时间序列。

```bash
caller=manager.go:651 component="rule manager" group=kube-apiserver.rules msg="Error on ingesting out-of-order result from rule evaluation" numDropped=231

caller=manager.go:651 component="rule manager" group=kube-apiserver-availability.rules msg="Error on ingesting out-of-order result from rule evaluation" numDropped=121
```

这个错误日志，可以从 Prometheus 代码中 [./prometheus/rules/manager.go](https://github.com/prometheus/prometheus/blob/release-2.28/rules/manager.go#L651) 的 Group.Eval() 方法中看到，每次评估规则时，只要有异常值得序列，都会抛出该错误日志信息。

```go
// Eval runs a single evaluation cycle in which all rules are evaluated sequentially.
func (g *Group) Eval(ctx context.Context, ts time.Time) {
	var samplesTotal float64
	for i, rule := range g.rules {
		select {
		case <-g.done:
			return
		default:
		}

		func(i int, rule Rule) {
            ......
			for _, s := range vector {
				if _, err := app.Append(0, s.Metric, s.T, s.V); err != nil {
					rule.SetHealth(HealthBad)
					rule.SetLastError(err)

					switch errors.Cause(err) {
					case storage.ErrOutOfOrderSample:
						numOutOfOrder++
						level.Debug(g.logger).Log("msg", "Rule evaluation result discarded", "err", err, "sample", s)
					case storage.ErrDuplicateSampleForTimestamp:
						numDuplicates++
						level.Debug(g.logger).Log("msg", "Rule evaluation result discarded", "err", err, "sample", s)
					default:
						level.Warn(g.logger).Log("msg", "Rule evaluation result discarded", "err", err, "sample", s)
					}
				} else {
					seriesReturned[s.Metric.String()] = s.Metric
				}
			}
			if numOutOfOrder > 0 {
				level.Warn(g.logger).Log("msg", "Error on ingesting out-of-order result from rule evaluation", "numDropped", numOutOfOrder)
			}
			if numDuplicates > 0 {
				level.Warn(g.logger).Log("msg", "Error on ingesting results from rule evaluation with different value but same timestamp", "numDropped", numDuplicates)
			}
			......
		}(i, rule)
	}
	if g.metrics != nil {
		g.metrics.GroupSamples.WithLabelValues(GroupKey(g.File(), g.Name())).Set(samplesTotal)
	}
	g.cleanupStaleSeries(ctx, ts)
}
```
