---
title: Target(目标) 与 Relabeling(重新标记)
linkTitle: Target(目标) 与 Relabeling(重新标记)
date: 2023-10-31T22:25
weight: 3
---

# 概述

> 参考：
>
> - [官方文档，配置 - 配置](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config)
> - [简书大佬](https://www.jianshu.com/p/c21d399c140a)

**Targets(目标)** 是 Prometheus 核心概念的其中之一，Targets 是**一组 Label(标签) 的集合。**

**Prometheus 在采集 Targets(目标) 的指标时，会自动将 Target 的标签附加到采集到的每条时间序列上**才存储，这样是为了更好的对数据进行筛选过滤，而这些附加的新标签是怎么来的呢？。。。这就是本文所要描述的东西。

如下所示，随便找一条时间序列，就可以看到，原始的指标中没有下图红框中的标签，而通过 Prometheus Server 采集后，就附加了两个新的标签上去。

```bash
]# curl -s localhost:9090/metrics | grep build_info
# HELP prometheus_build_info A metric with a constant '1' value labeled by version, revision, branch, and goversion from which prometheus was built.
# TYPE prometheus_build_info gauge
prometheus_build_info{branch="HEAD",goversion="go1.18.6",revision="1ce2197e7f9e95089bfb95cb61762b5a89a8c0da",version="2.37.1"} 1
```

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/prometheus/1616045780717-62e9699b-0064-4ea8-93d2-c741437d0a7f.png)

这里为什么会多出来两个标签呢，这种现象又是什么功能来实现的呢？~

首先，我们在 Prometheus Server 的 web 界面的 Status 标签中的 Targets 页面和 Service Discovery 页面中，可以发现 Prometheus 把标签分为两类：

- Discovered Labels(已发现的标签)
- Target Labels(目标标签)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/prometheus/1616045780732-4d75ae11-3fe5-4eb3-9866-ec9e55506a49.png)

从 Prometheus [Data Model(数据模型)](/docs/6.可观测性/监控系统/Prometheus/Storage(存储)/Data%20Model(数据模型).md) 中，可以知道 Label 的作用就是用来标识一条唯一的时间序列。那么 Prometheus 为什么会分为 Discovered Labels 和 Target Labels 呢？

- 因为，如果只有一套 Labels，那么 Prometheus 在采集目标时所使用的 Labels 就需要完完整整全部附加到已采集的指标上。当我们不想将所有标签都写入数据库，或者想改变某些标签的时候，这一套标签就已经不够用了。
- 既然出现了两套标签，也就需要一个两套标签之间的转换功能，而这就是 **Relabeling(重新标记)** 机制。

从某种程度上来说，Prometheus 一切介标签。Prometheus 有如下基本规定：

- Discovered Labels 用来描述采集目标的属性，只有根据这些属性才能找到目标。比如目标的 IP、PORT 等等。
  - 采集目标的属性标签名一般都以 `__` 符号开头。
- Target Labels 是用来附加到采集的指标上。<font color="#ff0000">这个才是真正被用户使用，以及存储的标签，PromQL 表达式中的标签也是指的 Target Labels</font>
  - `__` 符号开头的标签是不会被添加到 Target Labels 中，也就也不会附加到指标上。
- 其中 Discovered Labels 是根据配置自动生成的，而 Target Labels 则是通过一种称为 **Relabeling(重新标记)** 的功能生成的。

## Relabeling(重新标记) 功能

**Relabeling(重新标记)** 是一种可以在抓取目标数据之前动态重写目标的标签集的功能。每个抓取配置可以配置多个 Relabeling 机制。这些行为按照在配置文件中出现的顺序，由上至下依次应用于每个 Target(目标) 的标签集上。

### 为什么需要 Relabeling 功能呢？

- 由于 Prometheus 的规定，`__` 符号开头的标签不会被写入到时间序列数据中，也就无法再使用 promQL 语句进行数据过滤中使用，所以需要有一种方法，以便将`__` 符号开头的标签的值存储到一个新的标签中去，让这些新的标签写到时间序列数据中。而将标签 Relabeling 后，非 `__` 符号开头原始标签也不会消失，只是将原始标签里的值写入的目标标签中去了。
- 自动发现机制的原始标签都是 `__` 开头的，用户可以使用 Relabeling 机制，自己决定哪些需要留下，哪些需要抛弃。留下的那些又可以以新的标签名保留。
- 另外一个作用就是比较普遍的对标签值进行重新规划了，比如 原始标签有 IP 与 PORT，而实际只需要 IP，那么就需要用 relabel 功能了~~
- 总之，Relabeling 是一个非常强大的功能，可以重新规划每一条采集到的指标，也可以在持久存储前删除抓取到的指定指标。说白了，就是在**存到数据库之前，重新规划指标的信息**。因为时间序列是一组标签的集合作为唯一标识符。所以改标签，就是给时间序列改名了~~~~~

## Relabeling 功能的两个阶段

### 阶段一：发现目标之后，采集指标之前

### 阶段二：采集指标之后，储存指标之前

### 这两个阶段的区别

这两个阶段都可以使用 Relabeling 功能，不同点在于：

- 阶段一中，主要针对 Discovered Labels 进行操作，然后将 Target Label 附加到采集到的指标上。
  - 配置文件中的 `relabel_config` 字段工作在这个阶段。
- 阶段二中，可以针对采集到的指标中自带的标签进行操作。
  - 配置文件中的 `metric_relabel_configs` 字段工作在这个阶段。

# Discovered Labels 与 Target Labels

## 数据采集流程

当 Prometheus Server 加载配置文件启动时，并不是立刻就开始抓取 Target(目标) 的 Metrics(指标)

- 首先需要根据配置文件，获取目标信息，这些目标信息就是由一系列标签组成，称为 **Discovered Labels(已发现标签)**。
- 根据配置中的 Relabeling 配置，生成 **Target Labels(目标标签)**
  - 默认生成的 Target Labels 中将会删除 \_\_ 开头的标签，并将其他标签原封不动得映射到 Target Labels。
- 然后根据 Discovered Labels 的信息，从目标开始采集 Metrics，采集到 Metrics 后，将 Target Labels 附加到这些 Metrics 中。
- **<font color="#ff0000">从某种程度上来说，Prometheus 中一切皆标签</font>**

## Discovered Labels(已发现标签)

当 Prometheus 加载完 Target 后，会自动发现一些标签，这些就是 Discovered Labels。Discovered Labels 分为两部分：

- **系统信息标签**
  - **`__address__`** # 采集目标的 IP 和 PORT
  - **`__scheme__`** # 采集目标时，要使用的协议，HTTP 或者 HTTPS
  - **`__metrics_path__`** # 采集目标时，请求的访问路径
  - **`__param_XXXX`** # 配置文件中指定 params 字段的时，将会自动生成这类标签。
  - **`__meta_XXX`** # 通过服务发现功能发现的 Target 自带的元数据标签
    - 比如 kubernetes 发现，就是 \_\_meta_kubernetes_XXX 这种格式的标签
  - 其他系统标签
  - <font color="#ff0000">注意：这些前面带 `__` (两个`_` 符号)的标签，都是系统自动生成的，是无法在使用 PromQL 语句进行筛选过滤时直接使用的。</font>
    - <font color="#ff0000">并且，这些标签标签也不会变为 Target Labels</font>
  - **job** # 配置文件中的 job_name 字段的值就是 job 标签的值。
- **用户自定的标签**
  - 一般都是写在配置文件中的 labels 字段下的内容，这些内容经过 Relabels 之后，不会从 Discovered Labels 中删除。

Discovered Labels 中的 **系统标签**会告诉 Prometheus Server 如何从 Target 中获取时间序列数据(e.g.使用什么协议、从哪个 IP 上的主机哪个路径获取、等等类似的信息)。比如 `__address__`、`__metrics_path__`、`__scheme__` 这三个标签就表明，Prometheus 采集目标为 http://localhost:9090/metrics, 相当于执行了 `curl -XGET http://localhost:9090/metrics` 这个命令。如果还有 __param_XX 标签，则该标签的值，就是这次请求 URL 中的参数部分

> 默认情况下，Prometheus 会将 `__address__` 标签重新标记为 instance 标签，job 标签原封不同。如果想要其他的 Target Labels，则需要使用配置文件中的 relabel 以及 labels 字段来定义了。
>
> 如果重新标记步骤仅需要临时存储标签值（作为后续重新标记步骤的输入），请使用\_\_tmp 标签名称前缀。保证该前缀不会被 Prometheus 自己使用。

<font color="#ff0000">注意，被发现的所有标签中，标签名中的 `-`、`.`、`/` 等等特殊符号，都会转换成 `_` 符号。</font>

### 关于 \_\_meta_XX 标签

当 Prometheus 根据其自动发现机制，来自动发现待抓取目标时，会附带一些原始标签，这些标签以 `__meta_XX` 开头。
不同的服务发现配置发现标签不同，具体详见各种各种服务发现配置的官方文档(比如[这里就是 kubernetes_sd_configs 配置](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#kubernetes_sd_config)中，所有自动发现的标签)，效果如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/prometheus/1616049623207-f0fb653e-75ef-4e86-8c15-21bcb2292f02.png)

> 像 static_configs 这种直接指定抓取目标的配置，只会发现最基本的 **address**、**schem** 等标签。
>
> 这些各种服务发现配置发现标签就是用来描述目标属性的，然后可以通过 Relabeling 机制，将这些发现的标签保留下来，以便使用 PromQL 查询时，可以有更多的过滤选项。

## Target Labels(目标标签)

Target Labels 中的所有标签都是 Relabeling 之后的标签。Target Labels 包括两部分

- Prometheus 自身默认的 Relabel。就是去掉系统信息标签后剩下的标签
  - **instance** # 由 `__address__` 标签生成。<font color="#ff0000">注意：该标签会一直存在，无法通过 Relabeling 行为删除或改名</font>
  - **job** # 由 job 标签生成
- 通过配置文件配置的 relabel 标签

# Relabeling 配置

在 [Prometheus Server 配置](/docs/6.可观测性/监控系统/Prometheus/Server%20配置.md)中，Relabeling 行为的配置通常写在在 **relabel_configs** 和 **metrics_relabel_configs** 字段下。在很多地方都可以编写 **relabel_configs** 字段下的内容，以便为各种数据实现 relabel 功能

- `scrape_config.relabel_configs` 字段中 # Prometheus Relabeling 功能体现最主要的地方
- `alert_relabel_configs` 字段中 # 用于为告警内容实现 Relabeling 功能
- 等等等等，有很多地方都可以配置 Relabeling，还包括 Loki 日志套件中 [Promtail](/docs/6.可观测性/日志系统/Log%20Clients/Promtail/Promtail.md) 程序，也可以对日志流执行同样效果的 Relabeling 功能。因为 Relabeling 功能是 Prometheus 设计哲学 “**标签即一切**” 的必备功能

注意：

- **relabel_configs** 是 Prometheus 在发现**待采集的 Target(目标)** 之后，从目标采集指标之前，这两个行为之间发生的 Relabeling 行为配置。所以主要是针对 **待采集的 Target(目标) 及其标签**进行操作，而不是针对采集后的指标或其标签进行操作。<font color="#ff0000">所以，就如前文描述的一样，Relabeling 行为是针对待采集目标的 Discovered Label 的一种行为，经过 Relabeling 后，待采集目标就会生成 Target Labels</font>。
- **metrics_relabel_configs** 是 Prometheus 在采集到 Target(目标) 的指标之后，并存储到时序数据库之前，这两个行为之间发生的 Relabeling 行为配置。所以主要是针对**采集到的指标**进行操作。
  - 该配置不适用于自动生成的指标，比如 `up` 这类。因为这类指标在启动时就存在了，不用任何目标即可获取。
- 这两个配置的配置格式一摸一样。

**后文描述的 `目标` 二字，都是指 `待采集的目标`**

## relabel_configs 字段详解

**[source_labels](#source_labels)**([]STRING) # 从现有的标签中选择将要获取值的标签作为 source_labels。source_labels 可以有多个。

**[separator](#separator)**(STRING) # 指定 source_labels 中所有值之间的分隔符。`默认值: ;`。

**[target_label](#target_label)**(STRING) # 指定一个新标签名。replacement 字段的值会作为 target_label 指定的标签的标签值。

- 通过 regex 字段匹配到的值写入的指定的 target_label 中

**[regex](#regex)**(REGEX) # 从 source_label 获取的值进行正则匹配，匹配到的值写入到 target_label 中。`默认值: (.*)`，i.e.匹配所有值

**[modulus](#modulus)**(UINT64) # 去 source_labels 值的哈希值的模数

**[replacement](#replacement)**(STRING) # 替换。指定要写入 target_label 的值，STRING 中可以引用 regex 字段的值，使用正则表达式方式引用。`默认值：$1`。

- 与 action 字段的 replace 值配合使用。

**[action](#action)**(STRING) # 对匹配到的标签要执行的动作。`默认值: replace`。

### action

指定本次 Relabeling 的具体行为。`默认值：replace`。

Relabeling 的行为主要是围绕着 **Extracted Value(提取的值)** 进行的。Relabeling 行为将会对已提取的值进行正则匹配以获取匹配结果，再根据匹配结果生成新的标签。

这些**提取出来**的**待匹配的值**，其实本质上就分为两类：**标签名称** 与 **标签值**。而正则匹配中，正则表达式的内容，则是由 **regex 字段**指定的。根据这些不同类型的待匹配的值，我们可以将 Relabeling 的具体行为分类，当前 Promethus 支持如下几种 Relabeling 行为：

- **提取标签值进行匹配的行为**。从 `source_labels` 字段指定的标签名中提取所有的标签值，作为待匹配的值。
  - <font color="#ff0000">注意</font>：只有目标的 Discovered Labels 中包含 source_labels 字段中指定的标签，且这些标签的值能被 regex 字段的正则表达式匹配上，那么这些目标才会受到下列行为的影响。
  - **replace** # 为目标添加新标签，target_label 将会作为目标的新标签名，`replacement` 字段的值将会作为目标的新标签名的值。
    - 如果 regex 字段匹配不到任何内容，则不会进行替换。
    - <font color="#ff0000">与 labelmap 行为不同，除了更改标签名之外，还可以更改标签的值</font>
  - **keep** # 保留匹配到的目标。即只采集匹配到的目标的指标。
  - **drop** # 删除匹配到的目标。即不采集匹配到的目标的指标。
- **提取标签名进行匹配的行为**。从 Discovered Labels 中提取所有标签名，作为待匹配的值。
  - 注意：只有目标的 Discovered Labels 中，标签名能被 regex 字段的正则表达式匹配上，那么这些目标才会受到下列行为的影响
  - **labelmap** # 为目标添加新标签，`replacement` 字段的值将会作为目标的新标签名。
    - <font color="#ff0000">与 replace 行为不同，无法更改标签的值</font>。
  - **labelkeep** # 保留匹配到的标签。其余的标签移除
  - **labeldrop** # 移除匹配到的标签。其余的标签保留
    - 注意：labeldrop 和 labelkeep 这两个行为与 keep 和 drop 不同，仅仅是用来删除时间序列中某些标签的。小心使用这两个行为，以确保标签被删除后，指标仍然具有唯一的标签。
- **其他行为。**
  - **hashmod** # 设置 target_label 为的 modulus 哈希值的 source_labels，通过对指定 source_labels 进行 hash 计算，的出来一个新的 hash 值写入的 target_label 中。

<font color="#ff0000">注意</font>：

- action 不同的值代表不同的行为，而不同的行为也决定 relabel_configs 下可用的字段也不同。比如 labelkeep 行为，就无需 source_labels 等字段。
- 当 action 字段的值为 `replace`, `keep`, `drop`, `labelmap`,`labeldrop`、`labelkeep` 时，regex 字段是必须存在的。The regex is anchored on both ends. To un-anchor the regex, use `.*<regex>.*`。

### source_labels

指定一个或多个标签，提取这些标签的值作为待匹配的值。

- source_labels 可以有多个，<font color="#ff0000">若指定多个标签，则多个标签的值将会组合在一起作为待匹配的值，并以 `separator` 字段指定的分隔符分隔</font>。

### separator

指定 `source_labels` 字段提取出来的所有标签值之间的分隔符。`默认值： ;`。

### regex

针对提取出来的值进行正则表达式的匹配。`默认值：(.*)`，i.e.匹配提取出来的所有值，并将这些值放入一个组里。

### replacement

指定要写入新标签的标签值。`默认值：$1`。

- STRING 中可以引用 regex 字段的值，引用方式是正则表达式中的组引用。i.e.regex 字段中的第一个 ()，与正则表达式对 () 引用的概念相同，${0} 表示所有()中的值。

### modulus

去 source_labels 值的哈希值的模数

### target_label

指定一个新标签名。replacement 字段的值会作为名为 STRING 这个标签的标签值。

## action 字段下各种行为的特殊说明

### labelmap

```yaml
- regex: __meta_kubernetes_(namespace)
  action: labelmap
  replacement: $1
```

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/prometheus/1616045780763-37c765c1-bf4b-49db-8b05-e469a00618f8.png)

这个行为跟映射这个词特别搭，就是将通过正则表达式匹配到的标签，生成新的标签。

> 注意：不能这么玩 `regex: (__meta_kubernetes_namespace)`，因为 Target Labels 是要添加给采集到的 Metrics 上的，不能使用 `__` 符号开头的标签。

## relabel 配置的工作方式概述

1. 根据 action 字段指定的 Relabeling 行为，来提取待匹配的值。
2. 使用 regex 字段指定的正则表达式对这些提取出来的值进行匹配，匹配结果将会根据不同的行为，而作用在不用的地方。
3. relabel 机制中对于发现目标之后与采集目标之间发生的场景中，一般情况，每个 target 刚创建完，Prometheus 都会自动将 **address** 标签的值写入到 instance 标签中。
4. 其实，<font color="#ff0000">没有 source_label 也是可以的</font>，<font color="#ff0000">relabel 本质是定义一个标签</font>，而不是纯粹的替换，比如我不指定 source_label 字段，只指定 target_label 和 replacement 字段，就相当于是为这个 Metric 添加了一个标签
5. 所谓的 relabel，并不是绝对的替换，更像是定义 label
6. <font color="#ff0000">所以</font>，当我们在使用 Relabeling 功能时，首先应该先想自己到底要干什么，然后先决定 action 的值，再写其他的。
   1. 比如我想修改标签名，那么可用的 action 就是 replace 或者 labelmap，选择好具体的行为，再根据改行为支持的字段，写其它的。

# 最简单的配置样例

这些原始标签可以在配置中通过 relabel_configs 的配置段进行更改，这样这些标签的值在 proemetheus 存储起来之后，就会以新的标签标示

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s
scrape_configs:
  - job_name: "prometheus"
    static_configs:
      - targets:
          - localhost: 9090
    relabel_configs:
      - action: replace
        source_labels:
          - __metrics_path__ # 指定原始标签
        regex: (.*) # 指定原始标签值中要匹配的字符串
        replacement: $1 # 将原始标签值匹配到的字符串赋值给新标签
        target_label: metricsPath # 指定原始标签替换为的目标标签
```

下图的示例，就是将原始的 **metrics_path** 标签重新配置为新的 metricsPath 标签，新标签的值就是原始标签的值

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/prometheus/1616045780743-579762c3-abaa-4adf-99f8-3574b88460c2.jpeg)

## 高级样例

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/prometheus/1616045780735-c4c388d8-5602-4ffa-8c95-c78a4ff338d1.jpeg)

### 过滤 target

#### 使用 keep 行为，保留标签值匹配 regex 的 targets

```yaml
scrape_configs:
- …
- job_name: "cephs"
   relabel_configs:
   - action: keep
     source_labels:
     -  __address__
     regex:  ceph01.*
```

relabel 结果可以在 Prometheus 网页的 status/ Service Discovery 中查看

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/prometheus/1616045780731-a742cbc1-1a16-445a-8b17-1b0561db43e7.webp)

#### 使用 drop 行为，丢弃匹配 regex 的 Targets

```yaml
scrape_configs:
- …
- job_name: "cephs"
   relabel_configs:
     - action: drop
       source_labels:
         -  __address__
       regex:  ceph01.*

```

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/prometheus/1616045780860-e6406a1b-4bc1-4ba8-921e-39949436a028.webp)

### 删除标签

#### 使用 labeldrop 行为，将标签名为 job 的标签删除

```yaml
scrape_configs:
  - …
  - job_name: "cephs"
    relabel_configs:
      - regex: job
        action: labeldrop
```

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/prometheus/1616045780745-ad5eb6aa-c37b-431c-8593-ffb5ab8ff40e.webp)

labelKeep 和 labeldrop 不操作 \_\_ 开头的标签，要操作需要先改名

### 修改 label 名

#### 使用 replace 行为，将 scheme 标签改名为 protocol

```yaml
scrape_configs:
  - …
  - job_name: "cephs"
    relabel_configs:
      - action: replace
        source_labels:
          - __scheme__
        regex: (http)
        target_label: procotol
```

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/prometheus/1616045780754-bd537921-de9b-4cf6-b665-6951e6a4329c.webp)

这里可以是多个 source_labels，只有值匹配到 regex，才会进行替换

#### 使用 labelmap 行为，将原始标签的一部分转换为 target 标签，这一功能 replace 无法实现

```yaml
scrape_configs:
  - …
  - job_name: "sd_file_mysql"
    file_sd_configs:
      - files:
          - mysql.yml
        refresh_interval: 1m
    relabel_configs:
      - action: labelmap
        regex: (.*)(address)(.*)
        replacement: ${2}
```

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/prometheus/1616045780764-cc349a0e-af4b-4b83-b39d-6188ca93da4f.webp)

### 修改 label 值

配置 k8s 服务发现

```yaml
scrape_configs:
  - …
  - job_name: "sd_k8s_nodes"
    kubernetes_sd_configs:
      - role: node
        bearer_token_file: bearer_token
        tls_config:
          ca_file: ca.crt
        namespaces:
          names:
            - default
        api_server: https://master01:6443
```

服务发现完成后，默认 node 的 port 是 10250，会无法取得数据，同通过 relabel 修改标签.

```yaml
relabel_configs:
  - source_labels:
      - __address__
    regex: (.*)\:10250
    replacement: "${1}:10255"
    target_label: __address__
```

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/prometheus/1616045780766-2db0040a-e0ba-4b36-81e7-c1e1ed4a81cd.webp)

### 多标签合并

标签合并，可以将多个源标签合并为一个目标标签，可以取源标签的值，也可以进行 hash，用户 target 分组

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/prometheus/1616045780757-35f7df16-b3a1-4666-8b4a-ee755a1d6115.jpeg)

- 将多个标签的值进行 hash，形成一个 target 标签，只要 target 标签一致，则表示源标签一致，可以用来实现 prometheus 的负载均衡

```yaml
scrape_configs:
  - …
  - job_name: "sd_file_mysql"
    file_sd_configs:
      - files:
          - mysql.yml
        refresh_interval: 1m
    relabel_configs:
      - action: hashmod
        source_labels:
          - __scheme__
          - __metrics_path__
        modulus: 64
        target_label: hash_id
```

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/prometheus/1616045780783-56416a5f-6ee9-4c58-97b8-d99731d0fc9b.webp)

## 完整案例

以下是一个完整的 relabel 案例，这个案例包括

- 根据标签值过滤 target
- 合并标签值，并进行正则匹配
- 修改标签名
- 直接添加标签名

这个案例说明源标签是可以重复使用的

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/prometheus/1616045780790-516771ee-6fe1-49bf-aeb5-ee033ba1f35f.webp)
