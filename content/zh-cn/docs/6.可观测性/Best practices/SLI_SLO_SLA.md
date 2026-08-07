---
title: SLI/SLO/SLA
created: 2026-07-27T11:26
weight: 100
---

# 概述

> 参考：
>
> - [Wiki, Service level indicator](https://en.wikipedia.org/wiki/Service_level_indicator)
> - [Wiki, Service level objective](https://en.wikipedia.org/wiki/Service-level_objective)
> - [Wiki, Service level agreement](https://en.wikipedia.org/wiki/Service-level_agreement)
> - [公众号，通过 Prometheus 来做 SLI/SLO 监控展示](https://mp.weixin.qq.com/s/GNx0a0IKwvtDQ4QzEro2cA)

| 概念      | 英文全称                        | 核心含义                                        | 简单比喻（以餐厅为例）                    |
| ------- | --------------------------- | ------------------------------------------- | ------------------------------ |
| **SLI** | Service Level **Indicator** | **服务级别指标**：实际测量出的定量数据。                      | 厨房测量出的“实际平均上菜时间是25分钟”。         |
| **SLO** | Service Level **Objective** | **服务级别目标**：团队内部期望达到的技术线。比如"4 个 9"，"5 个 9"等。 | 店长给厨房定的目标：“95% 的菜品必须在30分钟内上桌”。 |
| **SLA** | Service Level **Agreement** | **服务级别协议**：对客户的法律/商业承诺，没达到要赔钱。              | 菜单上的承诺：“上菜超过40分钟，本单免费”。        |

SLI 是 SLO 的基础，SLO 是 SLA 的基础。先定 Indicator(指标)，再定该指标的 Objective(目标)，根据目标做出 Agreement(协议/承诺)。

> [!important] 病因指标 与 症状指标
> CPU 利用率单机阈值属于"病因指标"（cause metric），不是"症状指标"（symptom metric）。SRE 的建议是：**告警应该建立在症状上**——比如请求延迟升高、cgroup throttling、队列积压 —— 而不是直接对 CPU 数值本身设阈值。CPU 高不一定有实际影响（可能只是在跑批处理但服务正常），CPU 不高也不代表没问题（可能是被限流限制住了）。

# 告警治理：如何把每天 500+ 条告警降到 50 条

原文：[公众号 - 运维开发故事，告警治理：如何把每天 500+ 条告警降到 50 条](https://mp.weixin.qq.com/s/aVneK425N81S8Tm00C8kXA)

- SLO 系列前三篇
- [SLO 落地的工程实践：从 SLI 设计到燃烧率告警的体系化方法](https://mp.weixin.qq.com/s/SQCxQ59HBxOAcmhSuXeCjg)
- [用 Sloth 自动生成告警规则：从手写到自动化的 SLO 工程实践](https://mp.weixin.qq.com/s/O-Pl14_YhUhcMD5532PLdw)
- [Pyrra 实战指南：Kubernetes 原生的 SLO 监控工具](https://mp.weixin.qq.com/s/bWUbDCGFdunfqMqD_1642Q)

"告警邮件比我还积极，每天早上准时问候，但 99% 都是噪音！"

这句话来自一位 SRE 团队负责人的吐槽，但它是整个行业的缩影。Gartner 2025 年的报告指出： **85% 的运维告警不需要任何人工干预** 。PagerDuty 的统计更残酷——一个中等规模的 Kubernetes 集群，每周产生 500+ 条告警是常态。

更讽刺的是：团队搭建了完善的告警系统来"捕获每一个问题"，结果噪音淹没了一切，真正的故障反而被漏掉了。

这篇文章不讲理论，只讲实战。我会用两个真实案例（一个从 400 条/夜降到 8 条，一个从 1200 条/天降到 80 条），加上完整的 Prometheus + Alertmanager 配置，给你一套可以直接落地的告警治理方案。

> 前置阅读：本文是 SLO 系列的第四篇。前三篇讲了 SLO 体系全貌、Sloth 工具和 Pyrra 工具。这篇解决一个更基础的问题——在引入 SLO 之前或之后，怎么把每天 500+ 条告警的噪音降下来。

---

## 1 告警疲劳：不是抱怨，是系统性故障

先认清现实：告警疲劳不是一个"体验问题"，它是 **系统性故障** 。

### 三个让人坐不住的数字

- **85% — 告警无需人工干预（Gartner 2025）**
- **44% — 告警从未被查看（BigPanda 2025）**
- **#1 — 告警疲劳是 SRE 离职首因（DevOps Pulse 2025）**
- **4h+ — 告警疲劳下的平均 MTTR（对比正常 11 分钟）**

告警疲劳的杀伤链是这样的：

**告警太多 → 选择性忽略 → 真正的故障被淹没 → MTTR 暴增 → 团队失去信任 → 值班人员离职 → 剩余人员更累 → 更多告警被忽略**

这是一个正反馈恶性循环。不打破它，再好的 SLO 体系、再先进的工具都没用。

### 为什么 Kubernetes 让告警疲劳更严重

Kubernetes 的动态特性放大了告警噪音，五个典型场景：

| 场景 | 问题 | 一个故障产生多少条告警 |
| --- | --- | --- |
| Pod 重启风暴 | 一个 CrashLoopBackOff 配置错误，3 副本每 30 秒重启一次 | **360 条/小时** |
| HPA 扩缩容 | 把正常弹性伸缩当成故障告警 | 每次扩缩容 5-10 条 |
| 节点故障 | 一个节点挂了，上面所有 Pod、所有依赖服务的告警全触发 | **10-50 条** |
| 重复来源 | Prometheus + K8s Events + APM 同时告警同一个问题 | 3-4 倍重复 |
| 发布窗口 | 滚动更新期间 Pod 重启、健康检查失败、连接超时 | 20-100 条/次发布 |

一个 3 副本的 CrashLoopBackOff 就能在一小时内产生 360 条告警——这不是 360 个故障，这是 **一个故障在噪音中喊了 360 次** 。

---

## 2 诊断：你的告警系统病在哪里

治理之前先诊断。你需要回答三个问题：告警从哪来？谁在看？有没有用？

### 第一步：告警审计——导出全量告警规则

把 Prometheus 里所有的告警规则导出来，按以下四个维度分类：

| 分类 | 示例 | 处置策略 |
| --- | --- | --- |
| **可行动** | 数据库连接池耗尽 | 保留，优化阈值 |
| **信息性** | Pod 重启了一次 | 转为 Dashboard 指标 |
| **重复** | Prometheus 和 K8s Events 对同一个问题各告一次 | 去重，保留权威来源 |
| **过期** | 已下线服务的告警规则 | 直接删除 |

用这个 PromQL 查询过去 30 天触发次数最多的告警：

```yaml
# 查询过去 30 天触发最多的告警 Top 20
topk(20,
  sum by (alertname) (
    rate(ALERTS{alertstate="firing"}[30d])
  )
)
# 按触发次数排序，优先处理 Top 20
```

大多数团队做完审计后发现： **30-40% 的告警规则可以立即删除或转为 Dashboard 指标** 。

### 第二步：告警分类统计

抽样分析一周的告警数据，按来源和性质分类。如下：

| 告警类型 | 占比 | 典型示例 | 是否需要人工干预 |
| --- | --- | --- | --- |
| 重复告警 | 45% | 同一实例同一规则反复触发 | 否（去重即可） |
| 短暂抖动 | 25% | CPU 突刺，几分钟自动恢复 | 否（加 for 持续时间） |
| 发布/维护窗口 | 15% | 滚动更新期间的健康检查失败 | 否（静默窗口） |
| 无主告警 | 10% | 没有负责人的旧规则 | 否（删除或归属） |
| **真正需要处理** | **5%** | 服务不可用、错误率飙升 | **是** |

看到了吗？ **95% 的告警是噪音，只有 5% 是信号** 。这就是为什么 500+ 条告警和 50 条告警在信息量上没有区别——因为真正有价值的只有那 25 条。

### 第三步：告警健康度评估

用以下五个指标评估你的告警系统健康度：

| 指标 | 计算方式 | 健康值 | 危险值 |
| --- | --- | --- | --- |
| 告警确认率 | 被人工确认的告警 / 总告警数 | \>80% | <50% |
| 误报率 | 确认后无需处理的告警 / 总告警数 | <10% | \>40% |
| 重复率 | 同告警名 24h 内触发 >3 次的占比 | <5% | \>30% |
| 平均响应时间 | 告警触发到人工确认的中位数 | <5 分钟 | \>15 分钟 |
| 沉默率 | 告警群内 30 天无任何回复的占比 | <10% | \>50% |

> **铁律**
> 
> ：如果告警确认率低于 50%，说明你的告警系统已经失去了团队的信任。这不是需要"调优阈值"的问题，是需要推倒重建的问题。

---

## 3 第一刀：砍掉 60% 的垃圾告警

诊断完了，开始动刀。这一步不需要任何技术含量，只需要 **狠心** 。

### 告警无情日（Alert Ruthlessness Day）

来自一个英国金融服务团队的真实操作——他们在治理的第一天，运行了一条 SQL 查询：

```sql
-- 查找过去 30 天触发超过 100 次但从未导致客户影响的告警
SELECT alertname, COUNT(*) AS fires
FROM prometheus_alerts
WHERE fired_at > now() - INTERVAL '30 days'
GROUP BY alertname
HAVING COUNT(*) > 100
ORDER BY fires DESC;
```

结果：312 条告警在过去 30 天触发超过 100 次， **没有一条** 导致过客户影响的事件。

他们的操作： **第一天就删了 287 条告警规则** 。

> **判断标准**
> 
> ：如果一条告警在过去 12 个月内从未导致过客户故障，它就不是告警——它是 Dashboard 指标。

### 五类必须砍掉的告警

**1\. CPU/内存使用率阈值告警**

```yaml
# ❌ 删掉这种告警
- alert: HighCPU
  expr: 100 - (avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80
  for: 1m
  labels:
    severity: warning

# ✅ CPU 使用率是指标，不是告警
# 如果 CPU 高导致请求延迟增加，告警延迟而不是 CPU
# 如果 CPU 高导致 OOM，告警 OOM 事件而不是 CPU
```

CPU 80% 在某些服务上是正常状态（比如计算密集型任务），在另一些服务上是危险信号。直接告 CPU 使用率，无法区分这两种情况。 **告用户感知的症状，不告内部资源指标** 。

**2\. Pod 重启告警**

```yaml
# ❌ 删掉——Pod 重启是 K8s 自愈机制，不是故障
- alert: PodRestarted
  expr: increase(kube_pod_container_status_restarts_total[1h]) > 0

# ✅ 只告 CrashLoopBackOff——持续重启才是问题
- alert: PodCrashLooping
  expr: increase(kube_pod_container_status_restarts_total[1h]) > 5
  for: 10m
  labels:
    severity: warning
  annotations:
    summary: "Pod {{ $labels.pod }} 在过去1小时内重启超过5次"
    runbook_url: "https://runbook.example.com/crashloop"
```

**3\. HPA 扩缩容告警**

```yaml
# ❌ 删掉——HPA 扩缩容是设计行为，不是故障
- alert: HPA scaled
  expr: changes(kube_hpa_status_current_replicas[1h]) > 0

# ✅ 只告 HPA 无法扩容——真正的问题
- alert: HPAScalingLimited
  expr: kube_hpa_status_condition{condition="ScalingLimited", status="true"}
  for: 15m
  labels:
    severity: warning
```

**4\. 磁盘空间低告警（但阈值不合理）**

```yaml
# ❌ 删掉——阈值太低，频繁误报
- alert: DiskSpaceLow
  expr: 100 - (node_filesystem_avail_bytes / node_filesystem_size_bytes * 100) > 70
  for: 1m

# ✅ 用预测告警——给运维预留处理时间
- alert: DiskWillFillIn24h
  expr: predict_linear(node_filesystem_avail_bytes[1h], 24 * 3600) < 0
  for: 10m
  labels:
    severity: warning
  annotations:
    summary: "磁盘 {{ $labels.device }} 将在24小时内写满"
```

`predict_linear` 是 Prometheus 的线性预测函数，基于历史趋势预测未来值。比固定阈值告警好在：它不会在磁盘使用率 71% 时告警（可能一周都不会写满），只在 **真正即将写满时** 才告警。

**5\. 无 Runbook 的告警**

如果一条告警没有对应的 Runbook（处理手册），值班人员看到后不知道该做什么——这种告警等于没有告警。规则很简单：

> **没有 Runbook 的告警 = 没有告警**
> 
> 。要么补 Runbook，要么删掉。

---

## 4 从阈值告警到 SLO 燃烧率告警

砍完垃圾告警，下一步是把剩余的阈值告警 **升级为 SLO 燃烧率告警** 。这是告警治理的核心转型。

### 阈值告警 vs SLO 告警的本质区别

| 维度 | 阈值告警（噪音源） | SLO 燃烧率告警（信号源） |
| --- | --- | --- |
| 判定逻辑 | 指标 > 阈值 | SLO 违规 → 用户影响 |
| 响应心理 | "又是 CPU 告警，忽略吧" | "错误预算在烧，用户受影响" |
| 结果导向 | 关注组件状态 | 关注错误预算消耗速度 |
| 信噪比 | 极低（False Positive 频发） | 极高（每条都必须处理） |
| 维护成本 | 每个服务手动调阈值 | 一套规则适用于所有服务 |

### 实战：用燃烧率告警替换 47 条旧告警

一个真实案例：某团队的支付服务 payment-api 有 47 条阈值告警（CPU 高、内存高、5xx 多、延迟高、连接数多……），全部替换为 2 条 SLO 燃烧率告警：

```yaml
# payment-api 的 SLO 定义：可用性 99.9%
# 错误预算 = 1 - 99.9% = 0.1%
# 月预算 = 43200 秒 × 0.1% = 43.2 秒（允许宕机时间）

# Page 级告警：2% 预算在 1 小时内烧完（快燃烧）
- alert: PaymentApiSLOBurnRateFast
  expr: |
    (
      sum(rate(http_requests_total{job="payment-api", status=~"5.."}[1h]))
      /
      sum(rate(http_requests_total{job="payment-api"}[1h]))
    ) > (14.4 * 0.001)
    and
    (
      sum(rate(http_requests_total{job="payment-api", status=~"5.."}[5m]))
      /
      sum(rate(http_requests_total{job="payment-api"}[5m]))
    ) > (14.4 * 0.001)
  for: 2m
  labels:
    severity: page
    slo: payment-api-availability
  annotations:
    summary: "payment-api 错误预算快速燃烧"
    description: "2% 的月错误预算在1小时内消耗完毕，按当前速度3小时内预算耗尽"
    runbook_url: "https://runbook.example.com/payment-slo"

# Ticket 级告警：10% 预算在 3 天内烧完（慢燃烧）
- alert: PaymentApiSLOBurnRateSlow
  expr: |
    (
      sum(rate(http_requests_total{job="payment-api", status=~"5.."}[3d]))
      /
      sum(rate(http_requests_total{job="payment-api"}[3d]))
    ) > (1 * 0.001)
    and
    (
      sum(rate(http_requests_total{job="payment-api", status=~"5.."}[6h]))
      /
      sum(rate(http_requests_total{job="payment-api"}[6h]))
    ) > (1 * 0.001)
  for: 15m
  labels:
    severity: ticket
    slo: payment-api-availability
  annotations:
    summary: "payment-api 错误预算慢性消耗"
    description: "10% 的月错误预算在3天内消耗完毕，需要排查慢性问题"
    runbook_url: "https://runbook.example.com/payment-slo"
```

双窗口 AND 逻辑的含义（详见 SLO 系列第一篇）：

- **长窗口（1h/3d）判断"错误预算会不会烧光"——过滤瞬时抖动**
- **短窗口（5m/6h）判断"是否正在烧"——过滤慢性病漏报**
- 两个条件同时满足才告警——既不误报也不漏报

47 条阈值告警 → 2 条 SLO 告警。这就是为什么 SLO 系列的前三篇文章如此重要——它们给了你 **替代阈值告警的完整方法论和工具链** 。

> 💡 **不想手写燃烧率规则？** 用 Sloth 或 Pyrra 自动生成。上一篇 Pyrra 文章里展示了完整的 CRD 配置，一条 SLO 定义自动生成 19 条 Prometheus 规则，包括 4 级燃烧率告警。

---

## 5 Alertmanager 降噪三板斧

SLO 告警解决了"告什么"的问题，Alertmanager 解决"怎么告"的问题。三板斧： **分组、抑制、静默** 。

### 第一板斧：分组（Grouping）——把 N 条告警合并成 1 条

一个数据库挂了，依赖它的 50 个服务全部告警。没有分组，值班人员会收到 50 条通知。有分组，只收 1 条。

```yaml
# alertmanager.yml
route:
  group_by: ['alertname', 'cluster', 'service']
  group_wait: 30s        # 首次告警等待 30s，攒齐同一组的告警
  group_interval: 5m     # 同组新告警的合并间隔
  repeat_interval: 4h    # 重复提醒间隔
  receiver: 'oncall-slack'

  routes:
    # Critical 级别：立即通知，短间隔重发
    - matchers:
        - severity="critical"
      group_wait: 10s
      group_interval: 1m
      repeat_interval: 1h
      receiver: 'oncall-pagerduty'

    # Warning 级别：发 Slack，工作时间通知
    - matchers:
        - severity="warning"
      group_wait: 30s
      group_interval: 5m
      repeat_interval: 4h
      receiver: 'oncall-slack'

    # Info 级别：只发邮件日报
    - matchers:
        - severity="info"
      group_wait: 5m
      group_interval: 1h
      repeat_interval: 24h
      receiver: 'daily-email'
```

`group_by` 的选择很关键：

- `['alertname']`
	— 最粗粒度，同名的所有告警合为一条
- `['alertname', 'cluster', 'service']`
	— 推荐，按服务维度分组
- `['alertname', 'instance']`
	— 最细粒度，几乎不分组（不推荐）

### 第二板斧：抑制（Inhibition）——根因告警抑制症状告警

数据库挂了，不该再收到"服务 A 连接超时"、"服务 B 查询失败"、"服务 C 响应慢"这些症状告警——它们的根因都是数据库挂了。

```yaml
# alertmanager.yml — 抑制规则
inhibit_rules:
  # 规则1：数据库挂了，抑制所有依赖该数据库的服务告警
  - source_matchers:
      - alertname="DatabaseDown"
    target_matchers:
      - alertname="ServiceHighErrorRate"
    equal: ['cluster']

  # 规则2：节点挂了，抑制该节点上所有 Pod 的告警
  - source_matchers:
      - alertname="NodeDown"
    target_matchers:
      - alertname=~"PodCrashLooping|PodNotReady|ContainerOOMKilled"
    equal: ['node']

  # 规则3：Critical 告警抑制同服务的 Warning 告警
  - source_matchers:
      - severity="critical"
    target_matchers:
      - severity="warning"
    equal: ['alertname', 'cluster', 'service']

  # 规则4：集群不可达，抑制所有该集群的告警
  - source_matchers:
      - alertname="ClusterUnreachable"
    target_matchers:
      - cluster="{{ $labels.cluster }}"
    equal: ['cluster']
```

抑制规则的设计原则：

1. **根因在前：source\<em>matchers 是根因告警，target\</em>matchers 是症状告警**
2. **维度对齐：equal 字段确保只抑制同一维度的告警（不会把 A 集群的数据库告警抑制 B 集群的服务告警）**
3. **宁少勿多：抑制规则太激进会导致漏报，从最明确的根因开始配置**

### 第三板斧：静默（Silence）——维护窗口自动静默

发布、巡检、扩缩容期间的告警，90% 是噪音。但手动静默容易忘记恢复——更危险的是，静默了真实故障。

正确做法： **和发布系统联动，自动创建和回收静默规则** 。

```bash
# 发布系统在滚动更新前调用 Alertmanager API 创建静默
curl -XPOST http://alertmanager:9093/api/v2/silences \
  -H 'Content-Type: application/json' \
  -d '{
    "matchers": [
      {
        "name": "service",
        "value": "payment-api",
        "isRegex": false
      },
      {
        "name": "alertname",
        "value": "PodNotReady|ContainerOOMKilled|ProbeFailed",
        "isRegex": true
      }
    ],
    "startsAt": "2026-07-27T10:00:00Z",
    "endsAt": "2026-07-27T10:15:00Z",
    "createdBy": "deploy-bot",
    "comment": "payment-api 滚动发布，预计15分钟"
  }'

# 发布完成后，静默规则自动过期（endsAt 到期）
# 如果发布提前完成，主动删除静默规则
curl -XDELETE http://alertmanager:9093/api/v2/silence/{silence_id}
```

> ⚠️ **静默的陷阱** ： **不要静默 SLO 燃烧率告警！** 发布期间如果错误率真的飙升，燃烧率告警必须能触达值班人员。静默只适用于基础设施层面的噪音告警（Pod 重启、健康检查失败等），SLO 层面的告警永远不能被静默。

---

## 6 告警分级与路由：不是所有告警都要叫醒人

### 四级告警体系

| 级别 | 通知方式 | 响应时间 | 示例 |
| --- | --- | --- | --- |
| Critical | PagerDuty / 电话 | 立即（<5 分钟） | 服务完全不可用、SLO 快燃烧 |
| Warning | Slack / 钉钉 | 工作时间内（<1 小时） | SLO 慢燃烧、磁盘将满 |
| Info | 邮件日报 | 下一个工作日 | 接近阈值、趋势异常 |
| Silent | 仅记录到 Dashboard | 不通知 | Pod 重启、HPA 扩缩容 |

关键原则： **能发 Slack 解决的不打电话，能发邮件解决的不发 Slack，能看 Dashboard 的不发邮件** 。

### 按团队路由

```yaml
# alertmanager.yml — 按团队路由
route:
  receiver: 'default-slack'
  routes:
    # 支付团队
    - matchers:
        - team="payment"
      receiver: 'payment-slack'
      routes:
        - matchers:
            - severity="critical"
          receiver: 'payment-pagerduty'

    # 基础设施团队
    - matchers:
        - team="infra"
      receiver: 'infra-slack'
      routes:
        - matchers:
            - severity="critical"
          receiver: 'infra-pagerduty'

    # 业务团队
    - matchers:
        - team="business"
      receiver: 'business-slack'
```

每个告警规则必须带 `team` 标签。没有 `team` 标签的告警规则，CI 流水线直接拒绝合入。

---

## 7 Runbook 质量门禁：没有手册的告警等于没有告警

告警治理中最容易被忽视的一环：值班人员收到告警后 **不知道该做什么** 。

### Runbook 必须回答的四个问题

1. **客户现在受影响了吗？（Yes/No，附查询命令）**
2. **一键打开 Dashboard（直接链接）**
3. **如何验证和修复？（精确的命令，不是"检查一下日志"）**
4. **升级路径（找不到原因时该找谁）**

### 告警规则模板（强制 Runbook）

```yaml
# 每条告警必须包含以下 annotations
- alert: PaymentApiHighErrorRate
  expr: |
    sum(rate(http_requests_total{job="payment-api", status=~"5.."}[5m]))
    / sum(rate(http_requests_total{job="payment-api"}[5m])) > 0.01
  for: 5m
  labels:
    severity: page
    team: payment
    service: payment-api
  annotations:
    summary: "payment-api 5xx 错误率超过 1%"
    description: "当前错误率 {{ $value | humanizePercentage }}，持续 5 分钟"
    runbook_url: "https://runbook.example.com/payment-5xx"
    dashboard_url: "https://grafana.example.com/d/payment-api"
    quick_check: "kubectl logs -l app=payment-api --tail=50 | grep -i error"
    escalate_to: "#payment-oncall"
```

### CI 流水线门禁检查

```python
# scripts/validate-alerts.py — CI 检查脚本
import yaml, sys

def validate_alerts(rules_file):
with open(rules_file) as f:
        rules = yaml.safe_load(f)

    errors = []
    required_labels = ['severity', 'team', 'service']
    required_annotations = ['summary', 'runbook_url', 'dashboard_url']

for group in rules.get('groups', []):
for rule in group.get('rules', []):
if'alert'notin rule:
continue

            name = rule['alert']

# 检查必需 labels
for label in required_labels:
if label notin rule.get('labels', {}):
                    errors.append(f"Alert '{name}' missing label: {label}")

# 检查必需 annotations
for ann in required_annotations:
if ann notin rule.get('annotations', {}):
                    errors.append(f"Alert '{name}' missing annotation: {ann}")

# 检查 for 持续时间
if'for'notin rule:
                errors.append(f"Alert '{name}' missing 'for' duration")

if errors:
for e in errors:
            print(f"  X {e}")
        sys.exit(1)
else:
        print("  All alerts validated")

validate_alerts(sys.argv[1])
```

```yaml
# .gitlab-ci.yml 或 GitHub Actions
validate-alerts:
  stage: test
  script:
    - python scripts/validate-alerts.py alerts/rules.yaml
  # 不通过则阻止合入
```

> 💡 **Runbook 自动过期机制** ：Runbook 超过 90 天未更新 → 告警自动降级为 Silent（只记录 Dashboard，不通知）。这迫使团队保持 Runbook 的时效性。

## 8 30 天告警治理冲刺计划

理论讲完了，给你一个可以直接执行的 30 天计划。分四个阶段，每个阶段有明确的交付物和验收标准。

**Week 1：诊断与清剿**

- Day 1-2：导出全量告警规则，按四分类法审计（可行动/信息性/重复/过期）
- Day 3：运行高频告警查询（30 天触发 >100 次的），标记为待删除
- Day 4-5：删除过期和无主告警，将信息性告警转为 Dashboard 指标
- Day 6-7：为剩余告警补齐 team、severity、runbook\_url 标签
- **验收标准：告警规则数量减少 40%+，所有保留的告警都有 team 标签**

**Week 2：Alertmanager 降噪**

- Day 8-9：设计 `group_by` 策略，配置 `group_wait` / `group_interval` / `repeat_interval`
- Day 10-11：梳理服务依赖关系，配置抑制规则（根因→症状）
- Day 12：和发布系统对接，实现维护窗口自动静默
- Day 13-14：配置按 severity 和 team 的路由规则
- **验收标准：同一根因的告警合并为 1 条通知，维护窗口期间噪音告警为 0**

**Week 3：SLO 告警转型**

- Day 15-16：为核心服务定义 SLI（可用性 + 延迟），参考 SLO 系列第一篇
- Day 17-18：用 Sloth 或 Pyrra 生成燃烧率告警规则
- Day 19-20：新旧告警并行运行 48 小时，对比效果
- Day 21：确认无漏报后，删除旧的阈值告警规则
- **验收标准：核心服务的告警全部为 SLO 燃烧率告警，旧阈值告警归档**

**Week 4：度量与闭环**

- Day 22-23：搭建告警健康度 Grafana Dashboard（五组 PromQL）
- Day 24：配置 CI 门禁检查（Runbook、labels、annotations 强制校验）
- Day 25-26：建立每周告警复盘机制（Top 20 高频告警评审）
- Day 27-28：配置告警反馈闭环（每条 Page 必须回答：是否可行动 + 如何永久消除）
- Day 29-30：撰写治理报告，对比 30 天前后数据
- **验收标准：告警总量下降 80%+，误报率 <10%，MTTR 下降 50%+**

---

## 9 避坑指南：六个真实踩坑场景

**坑 1：一刀切删除所有阈值告警**

- **症状：听说 SLO 好就全删阈值告警，结果某些没有 SLI 的基础设施组件（如证书过期、DNS 解析）漏报了**
- **解法：SLO 告警覆盖请求类服务，基础设施告警用预测性阈值（predict\_linear），不要一刀切。分类处理：请求类 → SLO 告警；资源类 → 预测告警；事件类 → 日志告警**

**坑 2：抑制规则太激进导致漏报**

- **症状：配了"节点挂了抑制所有该节点的告警"，结果节点上某个独立服务的 OOM 被静默了 2 小时**
- **解法：抑制规则只抑制确定性的因果关系（数据库挂→连接超时），不确定的关系不要配抑制。宁可多收几条告警，也不要漏报**

**坑 3： `for` 持续时间设太短**

- **症状：for: 1m 的告警频繁触发又自动恢复，值班手机震个不停**
- **解法：通用建议——性能类告警 for: 5m-15m，可用性告警 for: 2m-5m，只有"服务完全不可用"这种 Critical 才用 for: 0-1m。配合 rate() 的窗口大小一起调（窗口通常取 for 的 4-5 倍）**

**坑 4：静默规则忘记回收**

- **症状：手动创建的静默规则，维护结束后忘记删除，导致后续真实故障被静默**
- **解法：静默规则必须设置 endsAt，不允许创建无过期时间的静默。和发布系统联动，自动创建+自动过期。定期检查过期但未删除的静默规则**

**坑 5：告警标签不一致导致分组失效**

- **症状：同一条告警因为 service 标签的值不一致（payment-api vs payment\_api vs PaymentAPI），无法被正确分组和抑制**
- **解法：在 CI 流水线里做标签规范化检查。所有 service、team、cluster 标签的值必须是预定义枚举列表中的值，拒绝不在列表中的值合入**

**坑 6：只减量不提质**

- **症状：告警从 500 条降到 50 条了，但值班人员还是觉得告警没用——因为剩下的 50 条里有 30 条不知道该怎么处理**
- **解法：减量只是手段，提质才是目的。每条保留的告警必须满足三个条件：可检测（Detectable）、可行动（Actionable）、有影响（Impactful）。不满足的要么补 Runbook，要么降级为 Dashboard 指标**

---

## 10 工具推荐：告警治理工具栈

| 工具                   | 用途                  | 推荐理由                                 |
| -------------------- | ------------------- | ------------------------------------ |
| Prometheus           | 告警规则引擎              | 原生支持 PromQL、recording rules、for 持续时间 |
| Alertmanager         | 告警路由与降噪             | 分组、抑制、静默三板斧，K8s 生态标配                 |
| Karma                | Alertmanager Web UI | 比原生 UI 强 10 倍，支持多 Alertmanager 聚合    |
| Sloth / Pyrra        | SLO 告警自动生成          | 声明式定义 SLO，自动生成燃烧率告警规则                |
| Grafana              | 告警健康度 Dashboard     | 可视化告警趋势、Top 高频告警、MTTR 等              |
| PagerDuty / Opsgenie | 值班调度与升级             | 排班管理、升级策略、SLA 跟踪                     |
| amtool               | Alertmanager CLI    | 命令行管理静默规则、测试路由配置                     |

### amtool 实用命令

```bash
# 查看当前所有 firing 告警
amtool --alertmanager.url=http://alertmanager:9093 alert query

# 测试告警路由（不实际发送）
amtool --alertmanager.url=http://alertmanager:9093 config routes test \
  --severity=critical --team=payment --service=payment-api

# 创建静默规则
amtool --alertmanager.url=http://alertmanager:9093 silence add \
  alertname=PodCrashLooping service=payment-api \
  --duration=15m \
  --comment="发布期间静默" \
  --author="deploy-bot"

# 查看所有静默规则
amtool --alertmanager.url=http://alertmanager:9093 silence query

# 删除静默规则
amtool --alertmanager.url=http://alertmanager:9093 silence expire <silence_id>
```

---

## 11 结语：告警治理是组织工程，不是技术工程

写到最后，必须说一个真相： **告警治理的难点从来不在技术，在组织和人** 。

Alertmanager 的分组、抑制、静默配置，半天就能学会。Sloth 和 Pyrra 的 SLO 规则生成，一天就能上手。真正难的是：

1. **让团队承认现状有问题——很多人已经把"每天 500 条告警"当成了常态，觉得"告警多说明系统在工作"。不，告警多说明监控系统在制造噪音。**
2. **让开发参与告警治理——告警不是 SRE 一个人的事。每条告警必须有 Runbook，Runbook 必须由最了解服务的开发来写。**
3. **建立持续治理的机制——不是治理一次就结束了。每周复盘 Top 20 高频告警，每月评估告警健康度，每季度审视告警策略。**
4. **把“减少告警”写进 KPI——如果团队的目标是"告警响应率 99%"，那大家会拼命加告警。如果目标是"每周 Page 不超过 2 次"，大家才会认真治理噪音。**

> **告警治理的终极目标**
> 
> ：让每一条发出来的告警都值得被看见、值得被处理、能够被闭环。当值班人员重新信任告警系统，当凌晨 3 点的电话响起时第一反应是"有真正的故障"而不是"又来烦我"——你就成功了。

2026 年，告警疲劳不是"正常现象"，是一种选择。500 条告警/天不是"系统复杂度的必然结果"，是组织债务的复利。该还债了。
