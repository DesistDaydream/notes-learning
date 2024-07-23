---
title: Stages(阶段) 详解
---

# 概述

> 参考：
> 
> - [官方文档,客户端-Promtail-阶段](https://grafana.com/docs/loki/latest/clients/promtail/stages/)
> - [官方文档,客户端-Promtail-配置-pipelinie_stages](https://grafana.com/docs/loki/latest/clients/promtail/configuration/#pipeline_stages)

对于 Stages 的详解，需要配合配置文件来描述。所以，本片文章即使 Stages 详解，也是 Promtail 配置文件中 `pipeline_stages` 字段的详解。

## pipeline_stages 字段配置

`pipeline_stages` 字段用于配置转换日志条目及其标签。promtail 运行流程中的 日志发现 步骤完成后，将执行 pipeline。

在大多数情况下，可以使用 regex 或 json 阶段从日志中提取数据。提取的数据将转换为临时 map 对象。这些提取出来的数据可以被 promtail 使用(比如这些数据可以作为标签的值或作为 i 内容直接输出)。此外，除 docker 和 cri 之外的任何其他阶段都可以访问提取的数据。

```yaml
scrape_configs:
- pipeline_stages:
  - docker: {}
  - cri: {}
  - regex:
    ...
  - json:
    ....
    ...... 阶段太多，其余略
```

# Parsing stages(解析阶段)

## docker 根据标准的 docker 日志文件格式来解析每行日志，并提取数据(默认行为)

来自 docker 的每行日志，都是以 JSON 格式编写，该 JSON 格式中有下列几个 key：

1. log # 日志行的具体内容
2. stream # 该 key 的值为 stdout 或 stderr，用来指明该日志行是标准输出还是标准错误
3. time # 日志行的时间戳

docker stage 会根据上述三种 key 来解析日志并提取其中数据，通过取出来的数据将创建出具有 3 个元素的 map。这些数据将会被其他 stages 所使用，并组合成 loki 可用的一行日志。

1. **output** # 与 log 对应。output stage 把该 key 的值变为发送到 loki 的一行日志。
2. **stream** # 与 stream 对应。labels stage 将该数据作为该行日志的 label。
3. **timestamp** # 与 time 对应。timestamp stage 将该数据作为 loki 记录的时间戳，并转换时间格式为 RFC3339Nano

上述 map 中的三个元素的键值对的值，就是老 key 的值

比如现在有这样一行 docker 日志：

```json
{
  "log": "log message\n",
  "stream": "stderr",
  "time": "2020-09-20T14:02:41.8443515Z"
}
```

在提取的数据中，将创建以下键值对

```yaml
output: log message\n
stream: stderr
timestamp: 2020-09-20T14:02:41.8443515
```

当该日志进入 loki 后，从 grafana 查看该数据，效果如下

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wzw6g5/1616129604931-8bc3cbd7-6057-4f9e-bdd3-6b33192313f4.jpeg)

本质上，docker 阶段的行为是 json、labels、timestamp、output 四个阶段的集合。各阶段行为如下：

```yaml
- json:
    output: log
    stream: stream
    timestamp: time
- labels:
    stream:
- timestamp:
    source: timestamp
    format: RFC3339Nano
- output:
    source: output
```

## cri - 使用标准的 CRI 日志格式来解析每行日志，并提取数据

本质上，cri 阶段的行为是 json、labels、timestamp、output 四个阶段的集合。各阶段行为如下：

```yaml
- regex:
  expression: "^(?s)(?P<time>\\S+?) (?P<stream>stdout|stderr) (?P<flags>\\S+?) (?P<content>.*)$",
- labels:
  stream:
- timestamp:
  source: time
  format: RFC3339Nano
- output:
  source: content
```

## regex - 使用正则表达式从每行日志提取数据

使用正则表达式提取数据，在 regex 中命名的捕获组支持将数据添加到提取的 Map 映射中。配置格式如下所示：

    regex:
      # RE2 正则表达式，每个捕获组必须被命名。
      expression: <string>
      # 从指定名称中提取数据，如果为空，则使用 log 信息。
      [source: <string>]

其中的 `expression` 是一个 Google RE2 正则表达式字符串，每个捕获组将被设置为到提取的 Map 中去，每个捕获组也必须命名：`(?P<name>re)`，捕获组的名称将被用作提取的 Map 中的键。
另外需要注意，在使用双引号时，必须转义正则表达式中的所有反斜杠。例如下面的几个表达式都是有效的：

- `expression: \w*`
- `expression: '\w*'`
- `expression: "\\w*"`

但是下面的这几个是无效的表达式：

- `expression: \\w*` - 在使用双引号时才转义反斜线
- `expression: '\\w*'` - 在使用双引号时才转义反斜线
- `expression: "\w*"` - 在使用双引号的时候，反斜杠必须被转义

例如我们使用下的不带 `source` 的 pipeline 配置：

    - regex:
        expression: "^(?s)(?P<time>\\S+?) (?P<stream>stdout|stderr) (?P<flags>\\S+?) (?P<content>.*)$"

当我们要抓取的日志数据为：

    2019-01-01T01:00:00.000000001Z stderr P i'm a log message!

该 pipeline 执行后以下键值对将被添加到提取的 Map 中去：

- `time`: `2019-01-01T01:00:00.000000001Z`
- `stream`: `stderr`
- `flags`: `P`
- `content`: `i'm a log message`

如果我们使用带上 `source` 的 pipeline 配置：

    - json:
        expressions:
          time:
    - regex:
        expression: "^(?P<year>\\d+)"
        source: "time"

如果需要抓取的日志数据为：

    { "time": "2019-01-01T01:00:00.000000001Z" }

则第一阶段将把以下键值对添加到提取的 Map 中：

- `time`: `2019-01-01T01:00:00.000000001Z`

而 regex 阶段将解析提取的 Map 中的时间值，并将以下键值对追加到提取的 Map 中去：

- `year`: `2019`

## json - 根据 JSON 的格式解析一行日志，并提取数据

通过将日志行解析为 JSON 来提取数据，也可以接受 `JMESPath` 表达式来提取数据，配置格式如下所示：

    json:
      # JMESPath 表达式的键/值对集合，键将是提取的数据中的键，而表达式将是值，被评估为来自源数据的 JMESPath。
      #
      # JMESPath 表达式可以通过用双引号来包装一个键完成，然后在 YAML 中必须用单引号包装起来，这样它们就会被传递给 JMESPath 解析器进行解析。
      expressions:
        [ <string>: <string> ... ]
      [source: <string>]

该阶段使用 Golang JSON 反序列化，提取的数据可以持有非字符串值，本阶段不做任何类型转换，在下游阶段将需要对这些值进行必要的类型转换，可以参考后面的 `template` 阶段了解如何进行转换。

> 注意：如果提取的值是一个复杂的类型，比如数组或 JSON 对象，它将被转换为 JSON 字符串，然后插入到提取的数据中去。

例如我们使用如下所示的 pipeline 配置：

    - json:
        expressions:
          output: log
          stream: stream
          timestamp: time

要抓取的日志行数据为：

    {
      "log": "log message\n",
      "stream": "stderr",
      "time": "2019-04-30T02:12:41.8443515Z"
    }

在提取的数据集中，将创建以下键值对：

- `output`: `log message\n`
- `stream`: `stderr`
- `timestamp`: `2019-04-30T02:12:41.8443515`

然后我们还可以用下面的 pipeline 配置来提前数据：

    - json:
        expressions:
          output: log
          stream: stream
          timestamp: time
          extra:
    - json:
        expressions:
          user:
        source: extra

要抓取的日志行数据为：

    {
      "log": "log message\n",
      "stream": "stderr",
      "time": "2019-04-30T02:12:41.8443515Z",
      "extra": "{\"user\":\"marco\"}"
    }

第一个 json 阶段执行后将在提取的数据集中创建以下键值对：

- `output`: `log message\n`
- `stream`: `stderr`
- `timestamp`: `2019-04-30T02:12:41.8443515`
- `extra`: `{"user": "marco"}`

然后经过第二个 json 阶段执行后将把提取数据中的 extra 值解析为 JSON，并将以下键值对添加到提取的数据集中：

- `user`: `marco`

此外我们还可以使用 JMESPath 表达式来解析有特殊字符的 JSON 字段（比如 `@` 或 `.`），比如我们现在有如下所示的 pipeline 配置：

    - json:
        expressions:
          output: log
          stream: '"grpc.stream"'
          timestamp: time

需要抓取的日志数据如下所示：

    {
      "log": "log message\n",
      "grpc.stream": "stderr",
      "time": "2019-04-30T02:12:41.8443515Z"
    }

在提取的数据集中，将创建以下键值对。

- `output`: `log message\n`
- `stream`: `stderr`
- `timestamp`: `2019-04-30T02:12:41.8443515`

需要注意的是在引用 `grpc.stream` 时，如果没有用单引号包裹的双引号，将无法正常工作。

## replace - 使用正则表达式替换数据

# Transform stages(转换阶段)

转换阶段用于对之前阶段提取的数据进行转换。

## multiline - 将多行日志进行合并，然后再将其传递到 pipeline 的下一个阶段。

多行阶段将多行日志进行合并，然后再将其传递到 pipeline 的下一个阶段。
一个新的日志块由**第一行正则表达式**来识别，任何与表达式不匹配的行都被认为是前一个匹配块的一部分。配置格式如下所示：

    multiline:
      # RE2 正则表达式，如果匹配将开始一个新的多行日志块
      # 这个表达式必须被提供
      firstline: <string>
      # 解析的最大等待时间（Go duration）: https://golang.org/pkg/time/#ParseDuration.
      # 如果在这个最大的等待时间内没有新的日志，那么当前日志块将被继续发送。
      # 如果被观察的应用程序因为异常而down掉了，该参数很有用，没有新的日志出现，并且异常块会在最大等待时间过后发送
      # 默认为 3s
      max_wait_time: <duration>
      # 一个多行日志块有的最大行数，如果该块有更多的行，就会认为是新的日志行
      # 默认为 128 行
      max_lines: <integer>

比如现在我们有一个 flask 应用，下面的日志数据包含异常信息：

    [2020-12-03 11:36:20] "GET /hello HTTP/1.1" 200 -
    [2020-12-03 11:36:23] ERROR in app: Exception on /error [GET]
    Traceback (most recent call last):
      File "/home/pallets/.pyenv/versions/3.8.5/lib/python3.8/site-packages/flask/app.py", line 2447, in wsgi_app
        response = self.full_dispatch_request()
      File "/home/pallets/.pyenv/versions/3.8.5/lib/python3.8/site-packages/flask/app.py", line 1952, in full_dispatch_request
        rv = self.handle_user_exception(e)
      File "/home/pallets/.pyenv/versions/3.8.5/lib/python3.8/site-packages/flask/app.py", line 1821, in handle_user_exception
        reraise(exc_type, exc_value, tb)
      File "/home/pallets/.pyenv/versions/3.8.5/lib/python3.8/site-packages/flask/_compat.py", line 39, in reraise
        raise value
      File "/home/pallets/.pyenv/versions/3.8.5/lib/python3.8/site-packages/flask/app.py", line 1950, in full_dispatch_request
        rv = self.dispatch_request()
      File "/home/pallets/.pyenv/versions/3.8.5/lib/python3.8/site-packages/flask/app.py", line 1936, in dispatch_request
        return self.view_functions[rule.endpoint](**req.view_args)
      File "/home/pallets/src/deployment_tools/hello.py", line 10, in error
        raise Exception("Sorry, this route always breaks")
    Exception: Sorry, this route always breaks
    [2020-12-03 11:36:23] "GET /error HTTP/1.1" 500 -
    [2020-12-03 11:36:26] "GET /hello HTTP/1.1" 200 -
    [2020-12-03 11:36:27] "GET /hello HTTP/1.1" 200 -

显然我们更希望将上面的 Exception 多行日志识别为一个日志块，在这个示例中，所有的日志块都是括号包括的时间开始的，所以我们可以用 `firstline` 正则表达式：`^\[\d{4}-\d{2}-\d{2} \d{1,2}:\d{2}:\d{2}\]` 来配置一个多行阶段，这将匹配上面我们的异常日志的开头部分，但是不会匹配后面的异常行，直到 `Exception: Sorry, this route always breaks` 这一行日志，这些将被识别为单个日志块，在 Loki 中也是以一个日志条目出现的。

    multiline:
      # 识别时间戳作为多行日志的第一行，注意这里字符串应该使用单引号。
      firstline: '^\[\d{4}-\d{2}-\d{2} \d{1,2}:\d{2}:\d{2}\]'
      max_wait_time: 3s

这个示例是假设我们对日志格式没有进行控制，所以我们需要一个更复杂的正则表达式来匹配第一行日志，但是如果我们能够控制被观察的日志格式，那么我们就可以简化第一行的匹配规则。
下面的是一个简单的 `Akka HTTP` 服务的日志：

    [2021-01-07 14:17:43,494] [DEBUG] [akka.io.TcpListener] [HelloAkkaHttpServer-akka.actor.default-dispatcher-26] [akka://HelloAkkaHttpServer/system/IO-TCP/selectors/$a/0] - New connection accepted
    [2021-01-07 14:17:43,499] [ERROR] [akka.actor.ActorSystemImpl] [HelloAkkaHttpServer-akka.actor.default-dispatcher-3] [akka.actor.ActorSystemImpl(HelloAkkaHttpServer)] - Error during processing of request: 'oh no! oh is unknown'. Completing with 500 Internal Server Error response. To change default exception handling behavior, provide a custom ExceptionHandler.
    java.lang.Exception: oh no! oh is unknown
    	at com.grafana.UserRoutes.$anonfun$userRoutes$6(UserRoutes.scala:28)
    	at akka.http.scaladsl.server.Directive$.$anonfun$addByNameNullaryApply$2(Directive.scala:166)
    	at akka.http.scaladsl.server.ConjunctionMagnet$$anon$2.$anonfun$apply$3(Directive.scala:234)
    	at akka.http.scaladsl.server.directives.BasicDirectives.$anonfun$mapRouteResult$2(BasicDirectives.scala:68)
    	at akka.http.scaladsl.server.directives.BasicDirectives.$anonfun$textract$2(BasicDirectives.scala:161)
    	at akka.http.scaladsl.server.RouteConcatenation$RouteWithConcatenation.$anonfun$$tilde$2(RouteConcatenation.scala:47)
    	at akka.http.scaladsl.util.FastFuture$.strictTransform$1(FastFuture.scala:40)
      ...

简单一看和其他日志一样，我们来看看日志的格式：

    <configuration>
        <appender name="FILE" class="ch.qos.logback.core.FileAppender">
            <file>crasher.log</file>
            <append>true</append>
            <encoder>
                <pattern>&ZeroWidthSpace;[%date{ISO8601}] [%level] [%logger] [%thread] [%X{akkaSource}] - %msg%n</pattern>
            </encoder>
        </appender>
        <appender name="ASYNC" class="ch.qos.logback.classic.AsyncAppender">
            <queueSize>1024</queueSize>
            <neverBlock>true</neverBlock>
            <appender-ref ref="STDOUT" />
        </appender>
        <root level="DEBUG">
            <appender-ref ref="ASYNC"/>
        </root>
    </configuration>

对于 Logback 配置来说，没有什么特别之处，除了在每个日志行的开头有一个 `&ZeroWidthSpace;`，这是零宽度空格的 HTML 代码，它使得识别第一行变得更加简单了，这里我们使用的第一行匹配正则表达式为：`\x{200B}\[`，`200B` 是零宽度空格字符的 Unicode 编码：

    multiline:
      # 将零宽度的空格确定为多行块的第一行，注意该字符串应使用单引号。
      firstline: '^\x{200B}\['
      max_wait_time: 3s

## template - 使用 Go 模板来修改提取出来数据

`template` 阶段可以使用 Go 模板语法来操作提取的数据。模板阶段主要用于在将数据设置为标签之前对其他阶段的数据进行操作，例如用下划线替换空格，或者将大写的字符串转换为小写的字符串。模板也可以用来构建具有多个键的信息。模板阶段也可以在提取的数据中创建新的键。
配置格式如下所示：

    template:
      # 要解析的提取数据中的名称，如果提前数据中的key不存在，将为其添加一个新的值
      source: <string>
      # 使用的 Go 模板字符串。 除了正常的模板之外
      # functions, ToLower, ToUpper, Replace, Trim, TrimLeft, TrimRight,
      # TrimPrefix, TrimSuffix, and TrimSpace 都是可以使用的函数。
      template: <string>s

比如下面的 pipeline 配置：

    - template:
        source: new_key
        template: "hello world!"

假如还没有任何数据被添加到提取的数据中，这个阶段将首先在提取的数据 Map 中添加一个空白值的 `new_key`，然后它的值将被设置为 `hello world!`。
在看下面的模板阶段配置：

    - template:
        source: app
        template: "{{ .Value }}_some_suffix"

这个 pipeline 在现有提取的数据中获取键为 app 的值，并将 `_som_suffix` 附加到值后面。例如，如果提前的数据 Map 的键为 app，值为 loki，那么这个阶段将把值从 loki 修改为 `loki_som_suffix`。

    - template:
        source: app
        template: "{{ ToLower .Value }}"

这个 pipeline 从提取的数据中获取键为 app 的值，并将其值转换为小写。例如，如果提取的数据键 app 的值为 LOKI，那么这个阶段将把值转换为小写的 loki。

    - template:
        source: output_msg
        template: "{{ .level }} for app {{ ToUpper .app }}"

这个 pipeline 从提取的数据中获取 `level` 与 `app` 的值，一个新的 `output_msg` 将被添加到提取的数据中，值为上面模板的计算结果。
例如，如果提取的数据中包含键为 app，值为 loki 的数据，level 的值为 warn，那么经过该阶段后会添加一个新的数据，键为 `output_msg`，其值为 `warn for app LOKI`。
任何先前提取的键都可以在模板中使用，所有提取的键都可用于模板的扩展。

    - template:
        source: app
        template: "{{ .level }} for app {{ ToUpper .Value }} in module {{.module}}"

上面的这个 pipeline 从提取的数据中获取 level、app 和  module 值。例如，如果提取的数据包含值为 loki 的 app，level 的值为 warn，moudule 的值为 test，则这个阶段会将提取数据 app 的值更改为 `warn for app LOKI in module test`。
任何之前获取的键都可以在模板中使用，此外，如果 `source` 是可用的，它可以在模板中被称为 `.Value`，我们这里 app 被当成了 source，所以它可以在模板中通过 `.Value` 使用。

    - template:
        source: app
        template: '{{ Replace .Value "loki" "blokey" 1 }}'

这里的模板使用 Go 的 `string.Replace`函数，当模板执行时，从提取的 Map 数据中的键为 app 的全部内容将最多有 1 个 loki 的实例被改为 blokey。
另外有一个名为 `Entry` 的特殊键可以用来引用当前行，当你需要追加或预设日志行的时候，这应该会很有用。

    - template:
        source: message
        template: "{{.app }}: {{ .Entry }}"
    - output:
        source: message

例如，上面的片段会在日志行前加上应用程序的名称。

> 在 Loki2.3 中，所有的 sprig 函数都被添加到了当前的模板阶段，包括 ToLower & ToUpper、Replace、Trim、Regex、Hash 和 Sha2Hash 函数。

# Actions stages(行动阶段)

用于从以前阶段中提取数据并对其进行处理。

## timestamp - 为日志条目设置时间戳的值

该阶段可以在将日志发送到 Loki 之前更改其时间戳。如果 timestamp 阶段不存在，则日志的时间戳默认为抓取日志条目的时间。

设置日志条目的时间戳值，当时间戳阶段不存在时，日志行的时间戳默认为日志条目被抓取的时间。
配置格式如下所示：

    timestamp:
      source: <string>
      # 解析时间字符串的格式，可以只有预定义的格式有：[ANSIC UnixDate RubyDate RFC822
      # RFC822Z RFC850 RFC1123 RFC1123Z RFC3339 RFC3339Nano Unix
      # UnixMs UnixUs UnixNs].
      format: <string>
      # 如果格式无法解析，可尝试的 fallback 的格式
      [fallback_formats: []<string>]
      # IANA 时区数据库字符串
      [location: <string>]
      # 在时间戳无法提取或解析的情况下，应采取何种行动。有效值为：[skip, fudge]，默认为 fudge。
      [action_on_failure: <string>]

其中的 `format` 字段可以参考格式如下所示：

- `ANSIC`: `Mon Jan \_2 15:04:05 2006`
- `UnixDate`: `Mon Jan_2 15:04:05 MST 2006`
- `RubyDate`: `Mon Jan 02 15:04:05 -0700 2006`
- `RFC822`: `02 Jan 06 15:04 MST`
- `RFC822Z`: `02 Jan 06 15:04 -0700`
- `RFC850`: `Monday, 02-Jan-06 15:04:05 MST`
- `RFC1123`: `Mon, 02 Jan 2006 15:04:05 MST`
- `RFC1123Z`: `Mon, 02 Jan 2006 15:04:05 -0700`
- `RFC3339`: `2006-01-02T15:04:05-07:00`
- `RFC3339Nano`: `2006-01-02T15:04:05.999999999-07:00`

另外支持常见的 Unix 时间戳：

- `Unix`: 1562708916 or with fractions 1562708916.000000123
- `UnixMs`: 1562708916414
- `UnixUs`: 1562708916414123
- `UnixNs`: 1562708916000000123

自定义格式是直接传递给  GO 的 `time.Parse` 函数中的 layout 参数，如果自定义格式没有指定 year，Promtail 会认为应该使用系统时钟的当前年份。

自定义格式使用的语法是使用时间戳的每个组件的特定值来定义日期和时间（例如 Mon Jan 2 15:04:05 -0700 MST 2006），下表显示了应在自定义格式中支持的参考值。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wzw6g5/1621834787456-e2597523-491b-482b-bf40-9b76610cd36e.png)

`action_on_failure` 设置定义了在提取的数据中不存在 `source` 字段或时间戳解析失败的情况下，应该如何处理，支持的动作有：

- `fudge（默认）`：将时间戳更改为最近的已知时间戳，总计 1 纳秒（以保证日志顺序）
- `skip`：不改变时间戳，保留日志被 Promtail 抓取的时间

比如使用下面的 pipeline 配置：

    - timestamp:
        source: time
        format: RFC3339Nano

经过上面的 timestamp 阶段在提取的数据中查找一个 time 字段，并以 `RFC3339Nano` 格式化其值（例如，2006-01-02T15:04:05.9999999-07:00），所得的时间值将作为时间戳与日志行一起发送给 Loki。

## output - 设置一行日志的文本。

也就是根据该配置，将解析出来的数据中的某些内容，作为发送给 loki 的一行日志的具体内容。也就是 loki 所记录的日志内容。

设置日志行文本，配置格式如下所示：

    output:
      source: <string>

比如我们有一个如下配置的 pipeline：

    - json:
        expressions:
          user: user
          message: message
    - labels:
        user:
    - output:
        source: message

需要收集的日志为：

    { "user": "alexis", "message": "hello, world!" }

在经过第一个 json 阶段后将提前以下键值对到数据中：

- `user`: `alexis`
- `message`: `hello, world!`

然后第二个 label 阶段将把 `user=alexis` 添加到输出的日志标签集中，最后的 output 阶段将把日志数据从原来的 JSON 更改为 message 的值 `hello, world!` 输出。

## labels - 更新日志条目的标签集(默认行为)

更新日志的标签集，并一起发送给 Loki。配置格式如下所示：

    labels:
      # Key 是必须的，是将被创建的标签名称。
      # Values 是可选的，提取的数据中的名称，其值将被用于标签的值。
      # 如果是空的，值将被推断为与键相同。
      [ <string>: [<string>] ... ]

比如我们有一个如下所示的 pipeline 配置：

    - json:
        expressions:
          stream: stream
    - labels:
        stream:

需要处理的日志数据为：

    {
      "log": "log message\n",
      "stream": "stderr",
      "time": "2019-04-30T02:12:41.8443515Z"
    }

第一个 json 阶段将提取 `stream` 到 Map 数据中，其值为 `stderr`。然后在第二个 labels 阶段将把这个键值对变成一个标签，在发送到 Loki 的日志行中将包括标签 `stream`，值为 `stderr`。

## metrics - 根据提取出来的数据计算指标

根据提取的数据计算指标。需要注意的是，创建的 metrics 指标不会被推送到 Loki，而是通过 Promtail 的 `/metrics` 端点暴露出去，Prometheus 应该被配置为可以抓取 Promtail 的指标，以便能够检索这个阶段所配置的指标数据。
配置格式如下所示：

    # 一个映射，key为metric的名称，value是特定的metric类型
    metrics:
      [<string>: [ <metric_counter> | <metric_gauge> | <metric_histogram> ] ...]

- **metric_counter**：定义一个 Counter 类型的指标，其值只会不断增加。
- **metric_gauge**：定义一个 Gauge 类型的指标，其值可以增加或减少。
- **metric_histogram**：定义一个直方图指标。

比如我们有一个如下所示的 pipeline 配置用于定义一个 Counter 指标：

    - metrics:
        log_lines_total:
          type: Counter
          description: "total number of log lines"
          prefix: my_promtail_custom_
          max_idle_duration: 24h
          config:
            match_all: true
            action: inc
        log_bytes_total:
          type: Counter
          description: "total bytes of log lines"
          prefix: my_promtail_custom_
          max_idle_duration: 24h
          config:
            match_all: true
            count_entry_bytes: true
            action: add

这个流水线先创建了一个 `log_lines_total` 的 Counter，通过使用 `match_all: true` 参数为每一个接收到的日志行增加。
然后还创建了一个 `log_bytes_total` 的 Counter 指标，通过使用 `count_entry_bytes: true` 参数，将收到的每个日志行的字节大小加入到指标中。
这两个指标如果没有收到新的数据，将在 24h 后小时。另外这些阶段应该放在 pipeline 的末端，在任何标签阶段之后。

    - regex:
        expression: "^.*(?P<order_success>order successful).*$"
    - metrics:
        successful_orders_total:
          type: Counter
          description: "log lines with the message `order successful`"
          source: order_success
          config:
            action: inc

比如上面这个 pipeline 首先尝试在日志中找到成功的订单，将其提取为 `order_success` 字段，然后在 metrics 阶段创建一个名为 `successful_orders_total` 的 Counter 指标，其值是在只有提取的数据中有 `order_success` 的时候才会增加。这个 pipeline 的结果是一个指标，其值只有在 Promtail 抓取的日志中带有 `order successful` 文本的日志时才会增加。

    - regex:
        expression: "^.* order_status=(?P<order_status>.*?) .*$"
    - metrics:
        successful_orders_total:
          type: Counter
          description: "successful orders"
          source: order_status
          config:
            value: success
            action: inc
        failed_orders_total:
          type: Counter
          description: "failed orders"
          source: order_status
          config:
            value: fail
            action: inc

上面这个 pipeline 首先会尝试在日志中找到格式为 `order_status=<value>` 的文本，将 `<value>` 提取到 `order_status` 中。该指标阶段创建了 `successful_orders_total` 和 `failed_orders_total` 指标，只有当提取数据中的 `order_status` 的值分别为 `success` 或 `fail` 时才会增加。

## tenant - 设置要用于日志条目的租户 ID 值。

设置日志要使用的租户 ID 值，从提取数据中的一个字段获取，如果该字段缺失，将使用默认的 Promtail 客户端租户 ID。配置格式如下所示：

    tenant:
      # source 或 value 配置选项是必须的，但二者不能同时使用（它们是互斥的）
      [ source: <string> ]
      # 当前阶段执行时用来设置租户 ID 的值。
      # 当这个阶段被包含在一个带有 "match" 的条件管道中时非常有用。
      [ value: <string> ]

比如我们有如下所示的 pipeline 配置：

    pipeline_stages:
      - json:
          expressions:
            customer_id: customer_id
      - tenant:
          source: customer_id

需要获取的日志数据为：

    {
      "customer_id": "1",
      "log": "log message\n",
      "stream": "stderr",
      "time": "2019-04-30T02:12:41.8443515Z"
    }

第一个 json 阶段将提取 `customer_id` 的值到 Map 中，值为 1。在第二个租户阶段将把 `X-Scope-OrgID` 请求 Header 头（Loki 用来识别租户）设置为提取的 `customer_id` 的值，也就是 1.
另外一种场景是用配置的值来覆盖租户 ID，如下所示的 pipeline 配置：

    pipeline_stages:
      - json:
          expressions:
            app:
            message:
      - labels:
          app:
      - match:
          selector: '{app="api"}'
          stages:
            - tenant:
                value: "team-api"
      - output:
          source: message

需要收集的日志数据为：

    {
      "app": "api",
      "log": "log message\n",
      "stream": "stderr",
      "time": "2019-04-30T02:12:41.8443515Z"
    }

这个 pipeline 将：

- Decode JSON 日志
- 设置标签 `app="api"`
- 处理匹配阶段，检查 `{app="api"}` 选择器是否匹配，如果匹配了则执行子阶段，也就是这里的租户阶段，覆盖值为 `"team-api"` 的租户。

此外在处理阶段还有 `labeldrop` 阶段，它从标签集中删除标签，这些标签与日志条目一起被发送到 Loki。还有一个 `labelallow` 阶段，它只允许将所提供的标签包含在与日志条目一起发送给 Loki 的标签集中。

# Filtering stages(过滤阶段)

## [match](https://grafana.com/docs/loki/latest/clients/promtail/stages/match/) - 依据指定的标签，过滤日志行，只有匹配到的日志行才会继续执行其他阶段

match 阶段是一个过滤阶段，当日志条目与可配置的 LogQL 流选择器和过滤器表达式匹配时，有条件地应用一组阶段或丢弃条目。

在配置文件中，match 可以嵌套一个 pipeline_stages 字段(match 字段下名为 stages)，也就意味着，Promtail 的所有阶段，都可以基于 match 匹配结果进行

```yaml
scrape_configs:
- pipeline_stages:
  - match:
      # 日志流选择器，选择要执行下述 stages 的日志行。
      selector: <string>

      # Names the pipeline. When defined, creates an additional label in
      # the pipeline_duration_seconds histogram, where the value is
      # concatenated with job_name using an underscore.
      # 命名管道。定义后，在pipeline_duration_seconds直方图中创建一个附加标签，其中该值使用下划线与job_name连接。
      [pipeline_name: <string>]

      # 嵌套上层的 pipeline_stages 字段内容。也就意味着，下面指定的各种阶段，只会在 match 匹配到的日志行中执行。
      stages:
      - [
          <docker> |
          <cri> |
          <regex>
          <json> |
          <template> |
          <match> |
          <timestamp> |
          <output> |
          <labels> |
          <metrics>
        ]
```

比如我们现在有一个如下所示的 pipeline 配置：

    pipeline_stages:
      - json:
          expressions:
            app:
      - labels:
          app:
      - match:
          selector: '{app="loki"}'
          stages:
            - json:
                expressions:
                  msg: message
      - match:
          pipeline_name: "app2"
          selector: '{app="pokey"}'
          action: keep
          stages:
            - json:
                expressions:
                  msg: msg
      - match:
          selector: '{app="promtail"} |~ ".*noisy error.*"'
          action: drop
          drop_counter_reason: promtail_noisy_error
      - output:
          source: msg

要处理的日志数据为：

    { "time":"2012-11-01T22:08:41+00:00", "app":"loki", "component": ["parser","type"], "level" : "WARN", "message" : "app1 log line" }
    { "time":"2012-11-01T22:08:41+00:00", "app":"promtail", "component": ["parser","type"], "level" : "ERROR", "message" : "foo noisy error" }

第一个 json 阶段将在第一个日志行的提取 Map 数据中添加值 `app=loki`，然后经过第二个 labels 阶段将 `app` 转换成一个标签。对于第二行日志也遵循同样的流程，只是值变成了 `promtail`。
然后在第三个 match 阶段使用 LogQL 表达式 `{app="loki"}` 进行匹配，只有在标签 `app=loki` 的时候才会执行嵌套 json 阶段，这里合我们的第一行日志是匹配的，然后嵌套的 json 阶段将 `message` 数据提取到 Map 数据中，key 变成了 `msg`，值为 `app1 log line`。
接下来执行第四个 match 阶段，需要匹配 `app="pokey"`，很显然这里我们都不匹配，所以嵌套的 json 子阶段不会被执行。
然后执行的第五个 match 阶段，将会删掉任何具有 `app="promtail"` 标签并包括 `noisy error` 文本的日志数据，并且还将增加 `logentry_drop_lines_total` 指标，标签为 `reason="promtail_noisy_error"`。
最后的 output 输出阶段将日志行的内容改为提取数据中的 msg 的值。我们这里的示例最后输出为 `app1 log line`。

## drop - Conditionally drop log lines based on several options.

drop 阶段可以让我们根据配置来删除日志。需要注意的是，如果你提供多个选项配置，它们将被视为 `AND` 子句，其中每个选项必须为真才能删除日志。如果你想用一个 `OR`子句来删除，那么就指定多个删除阶段。配置语法格式如下所示：

    drop:
      [source: <string>]
      # RE2 正则表达式，如果提供了 source，则会尝试匹配 source
      # 如果没有提供 source，则会尝试匹配日志行数据
      # 如果提供的正则匹配了日志行或者 source，则该行日志将被删除。
      [expression: <string>]
      # 只有在指定 source 源的情况下才能指定 value 值。
      # 指定 value 与 regex 是错误的。
      # 如果提供的值与`source`完全匹配，该行将被删除。
      [value: <string>]
      # older_than 被解析为 Go duration 格式
      # 如果日志行的时间戳大于当前时间减去所提供的时间，则将被删除
      [older_than: <duration>]
      # longer_than 是一个以 bytes 为单位的值，任何超过这个值的日志行都将被删除。
      # 可以指定为整数格式的字节数：8192，或者带后缀的 8kb
      [longer_than: <string>|<int>]
      # 每当一个日志行数据被删除，指标 `logentry_dropped_lines_total` 都会增加。
      # 默认的 reason 标签是 `drop_stage`，然而你可以选择指定一个自定义值，用于该指标的 "reason" 标签。
      [drop_counter_reason: <string> | default = "drop_stage"]

比如我们有一个如下所示的简单 drop 阶段配置：

    - drop:
        expression: ".*debug.*"

该阶段将删除任何带有 `debug` 字样的日志行。
如果是下面的配置示例：

    - json:
        expressions:
          level:
          msg:
    - drop:
        source: "level"
        expression: "(error|ERROR)"

则下面的日志数据都将被删除：

    {"time":"2019-01-01T01:00:00.000000001Z", "level": "error", "msg":"11.11.11.11 - "POST /loki/api/push/ HTTP/1.1" 200 932 "-" "Mozilla/5.0 (Windows; U; Windows NT 5.1; de; rv:1.9.1.7) Gecko/20091221 Firefox/3.5.7 GTB6"}
    {"time":"2019-01-01T01:00:00.000000001Z", "level": "ERROR", "msg":"11.11.11.11 - "POST /loki/api/push/ HTTP/1.1" 200 932 "-" "Mozilla/5.0 (Windows; U; Windows NT 5.1; de; rv:1.9.1.7) Gecko/20091221 Firefox/3.5.7 GTB6"}

然后使用下面的配置来删除老的日志数据：

    - json:
        expressions:
          time:
          msg:
    - timestamp:
        source: time
        format: RFC3339
    - drop:
        older_than: 24h
        drop_counter_reason: "line_too_old"

> 需要注意的是为了让 `old_than` 发挥作用，你必须在应用 drop 阶段之前，使用时间戳阶段来设置抓取日志行的时间戳。

比如当前的摄取时间为 `2021-05-01T12:00:00Z`，当从文件中读取时，会删除这个日志行：

    {"time":"2021-05-01T12:00:00Z", "level": "error", "msg":"11.11.11.11 - "POST /loki/api/push/ HTTP/1.1" 200 932 "-" "Mozilla/5.0 (Windows; U; Windows NT 5.1; de; rv:1.9.1.7) Gecko/20091221 Firefox/3.5.7 GTB6"}

但是下面的日志数据不会被删除：

    {"time":"2021-05-03T12:00:00Z", "level": "error", "msg":"11.11.11.11 - "POST /loki/api/push/ HTTP/1.1" 200 932 "-" "Mozilla/5.0 (Windows; U; Windows NT 5.1; de; rv:1.9.1.7) Gecko/20091221 Firefox/3.5.7 GTB6"}

在这个例子中，当前时间是 \`\`2021-05-03T16:00:00Z`，`older_than`是 24h。所有时间戳超过`2021-05-02T16:00:00Z`的日志行都将被删除。<br />这个删除阶段删除的所有行也将增加`logentry_drop_lines_total`指标，并标明原因为`"line_too_old"\`。
下面是另外一个复杂点的配置：

    - json:
        expressions:
          time:
          msg:
    - timestamp:
        source: time
        format: RFC3339
    - drop:
        older_than: 24h
    - drop:
        longer_than: 8kb
    - drop:
        source: msg
        regex: ".*trace.*"

上面的 pipeline 执行后将删除掉所有超过 24 小时**或者**超过 8kb 的日志**或者** json 的 msg 值中包含 `trace` 字样的日志。
