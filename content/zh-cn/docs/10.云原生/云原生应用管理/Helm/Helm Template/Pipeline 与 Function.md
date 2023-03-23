---
title: Pipeline 与 Function
---

# 概述

> 参考：
>
> - 官方文档：<https://helm.sh/docs/chart_template_guide/functions_and_pipelines/>

将 .Values 对象中的字符串注入模板时，应引用这些字符串。我们可以通过在 **Template Directive(模板指令)**中调用 quota 函数来实现，比如下面这个示例，会将引入的指转换为字符串类型：

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{.Release.Name}}-configmap
data:
  myvalue: "Hello World"
  drink: {{quote .Values.favorite.drink}}
  food: {{quote .Values.favorite.food}}
```

模板函数遵循语法 `functionName arg1 arg2...`。在上面的代码片段中，`quote .Values.favorite.drink` 调用 quote 函数并将一个参数传递给它。

Helm 拥有超过 60 种可用函数。其中一些是由  [**Go template language(Go 模板语言)**](https://pkg.go.dev/text/template) 本身定义的。其他大多数都是 [**Sprig template library(Sprig 模板库)**](https://pkg.go.dev/github.com/Masterminds/sprig) 的一部分。随着示例的进行，我们将看到其中的许多例子。

注意：虽然我们将 **Helm template language(Helm 模板语言)** 视为 Helm 特有的，但它实际上是 Go 模板语言，一些额外函数和各种包装器的组合，以将某些对象暴露给模板。Go 模板上的许多资源在了解模板时可能会有所帮助。

# Pipeline 管道

> 在 Helm 的 Template 中，**Pipeline(管道)**的概念与 Go Template 中 Pipeline 的概念不同，并不是指产生数据的操作。

**Pipeline(管道)** 与 Linux 中管道的概念类似，用于连接多个 [**Template Directives(模板指令)**](https://www.teambition.com/project/5f90e312755d8a00446050eb/app/5eba5fba6a92214d420a3219/workspaces/5f90e312c800160016ea22fb/docs/5f9a633937398300016bed65?scroll-to-block=5f9a6348246f30f3eef35c3e)。换句话说，**Pipeline(管道)**是一种按顺序完成多项任务的有效方式。

本篇文章开头的示例，如果使用 Pipeline 重写，则是这样的：

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Release.Name }}-configmap
data:
  myvalue: "Hello World"
  drink: {{ .Values.favorite.drink | quote }}
  food: {{ .Values.favorite.food | quote }}
```

在此示例中，我们没有调用 `quote ARGUMENT`，而是颠倒了顺序。使用 Pipeline 符号 `|` 将参数“发送”到函数：`.Values.favorite.drink | quote`。使用 Pipeline，我们可以将多个功能链接在一起，比如：

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Release.Name }}-configmap
data:
  myvalue: "Hello World"
  drink: {{ .Values.favorite.drink | quote }}
  food: {{ .Values.favorite.food | upper | quote }}
```

# Function 函数

> 参考：
>
> - [官方文档，Chart 模板指南-模板函数列表](https://helm.sh/docs/chart_template_guide/function_list/)

**Helm Template Function(Helm 模板函数)** 可以用来丰富模板功能，通过函数，可以对传入模板的数据进行更多操作，以便让这些数据更符合我们的预期。

简单示例：

这是一个模板文件

```json
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{.Release.Name}}-configmap
data:
  drink: {{.Values.favorite.drink}}
  food: {{quote .Values.favorite.food}}
```

这是值文件

```json
favorite:
  drink: tee
  food: bread
```

渲染模板结果：

```shell
~]# helm template test .
---
# Source: test/templates/test.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-configmap
data:
  drink: tee # 未使用 quote 函数
  food: "bread" # 使用了 quote 函数
```

在模板中，使用了一个名为 quote 的函数，这个函数的功能是，可以将引用的对象的值都机上双引号。Helm 里所有可用的函数列表：<https://helm.sh/docs/chart_template_guide/function_list/>

## default 函数

模板中经常使用函数是 default，语法为：

```shell
default DEFAULT_VALUE GIVEN_VALUE
```

参数：

1. DEFAULT_VALUE    # 指定默认值
2. GIVEN_VALUE    # 默认值将会传给 GIVEN_VALUE

default 函数用来为模板内部的指令指定默认值

上面的示例修改一下：

```shell
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{.Release.Name}}-configmap
data:
  drink: {{.Values.favorite.drink}}
  food: {{quote .Values.favorite.food}}
  food_default: {{default "bread" .Values.favorite.food}}
```

此时，如果我们将 value.yaml 中 food: bread 注释掉，会得到如下结果：

```shell
[root@master-1 test]# helm template test .
---
# Source: test/templates/function_pipeline.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-configmap
data:
  drink: tee
  food:
  food_default: bread
```

## lookup 函数

lookup 函数可以用于查找 kubernetes 集群中的资源，语法：

```shell
lookup APIVERSION KIND NAMESPACE NAME
```

比如：

```json
(lookup "v1" "Namespace" "" "default").metadata.annotations
```

上述指令将会获取 default 名称空间的 metadata.annotaions 字段的信息。类似于 kubectl get ns default --template={{.metadata.annotations}} 命令，与该命令获取的值一致。

注意：当使用 helm template 命令或者使用 --dry-run 标志时，并不会与 Kubernetes 的 API Server 建立联系，所以 lookup 在上述两种情景下，无法获取具体的值，只会返回一个空的 map 。

## Operators 函数

Operators 运算符，在 helm 模板里也当作函数来看。eq、ne、lt、gt、and、or 等等，均视为函数

与正常情况不同，运算符作为函数时，需要将这个关键字放在语法开头，后面跟参数。比如：

**eq # 比较两个 Pipeline 是否相等**

```shell
# 比较 ARG1 与 ARG2 是否相等
# 也可以指定多个参数，所有参数都是与 ARG1 进行比较
eq ARG1 ARG2
```

**and # 与运算，当两个 Pipeline 都为真时，结果为真**

```shell
# 如果 .Values.fooString 值存在，并且 .Values.fooString 的值为 foo，则执行后面的指令
{{if and .Values.fooString (eq .Values.fooString "foo") }}
    {{...}}
{{end}}
```

## 字典类型数据处理函数

### [mergeOverwrite, mustMergeOverwrite](https://helm.sh/docs/chart_template_guide/function_list/#mergeoverwrite-mustmergeoverwrite)

将两个或多个字典合并为一个，右侧的优先级最高

Syntax(语法)
`$NewDict := mergeOverwrite $DEST $SOURCE1 $SOURCE2`

简单示例：

```yaml
dst:
  default: default
  overwrite: me
  key: true

src:
  overwrite: overwritten
  key: false
```

通过 `mergeOverwrite .Values.dst .Values.src` 语句我们可以得到如下结果：

```yaml
newdict:
  default: default
  overwrite: overwritten
  key: false
```
