---
title: Helm Template
---

# 概述

> 参考：
>
> - [官方文档，Chart 模板指南](https://helm.sh/docs/chart_template_guide/)

Helm 的 [Template(模板)](https://helm.sh/docs/chart_template_guide/getting_started/) 可以生成 manifests 文件，这些文件是 kuberntes 用于描述资源的 yaml 。

## Template 背景

helm 作为 kuberntes 的包管理器，用来在 k8s 集群中安装应用程序。众所周知，对于 k8s 来说，一个应用不应该是一个简单的 pod，应该包含该 pod 的运行方式(比如 deployment)、pod 的配置文件(configmip)、这个程序如何对外提供服务(service、ingress)等等等一系列的信息。这些信息都是通过 yaml 来描述如何工作的。可以想象，如果一个应用程序，其实是一堆 yaml 文件的话，那么 helm 本质上就是管理这些 yaml 文件的。而一个应用程序想要让用户来使用，必然还涉及到自定义的问题。比如应用程序的名字、配置文件中的内容、对外提供服务所要暴露的端口等等信息。

既然有这样的需求，那么为了让一个应用程序可以自定义，template 就应运而生。template 就是可以将这些 yaml 文件中的 value 变成一种变量的形式，然后通过其他方式(helm 命令行 --set 标志或者 value.yaml 文件等)来对这些变量进行赋值，来实现应用程序自定义的效果。

helm templete 使用 Go 语言的 [Template](/docs/2.编程/高级编程语言/Go/Go%20规范与标准库/Template.md) 来实现。而 go template 具有丰富的功能，除了可以普通赋值以外，还可以使用控制结构(比如 if...else、range 等)来将赋值的过程更具体和多次赋值。

注意：
当我们谈论“ Helm 模板语言”时，就好像它是特定于 Helm 一样，但它实际上是 Go 模板语言，一些额外功能以及各种包装程序的组合，以将某些对象暴露给模板。当您了解模板时，Go 模板上的许多资源可能会有所帮助。

## 基本示例

首先通过官方的 [hello world](https://helm.sh/docs/chart_template_guide/getting_started/) 示例来看一下 helm 模板的基本功能。通过示例，可以看到，yaml 中的 `{{ .Release.Name }}` 被替换为了 release 的名字。这就相当于把 release 名字传递给了`{{ .Release.Name }}`。

`{{ }}`符号当中的内容称为 **Template Directive(模板指令)**，`{{ }}`符号中的内容，类似于一个简单的 shell 脚本。一个 **Template Directive(模板指令)** 可以包含多种内容，比如：

1. **引用对象** # 赋值功能。用于从外部传递值到模板文件中。
2. **控制结构** # 通过类似 if...else 这种语句来丰富模板中的赋值功能
3. **使用变量** # 声明一个变量并赋值，可以在后续引用。
4. **使用函数** # 通过函数来执行高级操作，比如将值进行大小写转换等等
5. **等等等......**

Note：如果想要测试模板渲染的效果，可以使用 helm install --debug --dry-run mychart ./mychart 这种命令来渲染模板，但是不会真的安装 chart

# Objects 介绍

上面基本示例中的 **.Release.Name 指令** 就是在 **引用 objects(对象)** 。

## 关于`.`这个符号的理解

而 object 开头的`.`表示当前作用域的顶级对象(有点类似 DNS 中的根域名)，可以在顶级对象下直接引用的对象称为 [Built-in Objects(内置对象)](https://helm.sh/docs/chart_template_guide/builtin_objects/)，Release 就是一个内置对象。

在上面的 [hello world](https://helm.sh/docs/chart_template_guide/getting_started/) 示例中，从顶层对象开始，找到 Release 对象，在 Release 对象中找到 Name 对象，并将 Name 对象的值传递进模板里。（这种写法特别像 DNS 中的写法，通过点 . 来区分域名的级别，在 helm 模板中也有类似的概念，通过点 . 来区分对象之间级别，上级对象包含其内所有的下级对象）

官方对于 `.` 这个符号称之为 **scope(范围)**，比如下图，圆形就代表了一个范围，最外层就是一个 `.` 符号。一个 . 范围的内置对象有 Values、Release 等等，而 **.Values 范围**下还有各种 **.Values.XXX 范围**，直到范围内的对象下没有新的范围。**scope(范围)** 常用来约束或指定模板指令的作用范围，或作为参数传递给函数来告知函数其可操作的范围有多大。比如 [Named Templates(命名模板)](/docs/10.云原生/云原生应用管理/Helm/Helm%20Template/Named%20Templates(命名模板).md) 中，就将范围作为参数传递给 template，告知 Named Templates 内引用对象时使用的范围是多大。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nladcy/1617243339047-b1ddaebb-7a5d-41b2-a83e-5f6e9c1c52f0.png)

可以这么说，这个 `.` 符号有两重含义：1、限定范围。2、引用。

这个 scope(范围) 还可以这么理解。通过 helm 的控制结构语句 range 可以以另一种方式来理解 {{ . }} 的行为，详见：[控制结构与变量](/docs/10.云原生/云原生应用管理/Helm/Helm%20Template/控制结构与变量.md) 章节中的 《通过 range 来梳理 . 范围内的对象》段落来加深理解

实际上，说了这么多概念，其实 `.` 符号就是代表一个值。这个值可以是一个单独的字符串出或者数字、也可以是一个 map 或者 array。

根据上面模板的介绍可以看到，在模板中，对象就是起到类似变量的用作(但不是变量)，用来传递值的。object(对象) 通过 template engin(模板引擎) 将值传递到模板中。而且，在模板中，可以创建新的 object 。

Objects 可以很简单，只有一个值。也可以包含其他对象或者功能。比如 Release 对象包含 Release.Name 对象。而 Files 对象，则具有一些功能(比如将文件进行 base64 编码、将文件内容传递到模板中、遍历一个文件)

## Helm 里自带的内置 Objects

1. Release # 该 Object 用来描述 release。其内部还有如下几个对象
   1. Release.Name # 该 release 的名字
   2. Release.Namespace # 该 release 所在的 k8s 的 namesapce 名称
   3. Release.IsUpgrade # 如果当前操作是升级或回滚，则该对象的值为 true
   4. Release.IsInstall # 如果当前操作是安装，则该对象的值为 true
   5. Release.Revision # The revision number for this release. On install, this is 1, and it is incremented with each upgrade and rollback.
   6. Release.Service # The service that is rendering the present template. On Helm, this is always Helm.
2. Values # 该 Object 从 `values.yaml 文件` 和 `helm 命令的 --values和--set 标志` 中提取值，并传递到模板中。默认情况下，Values 为空
3. Chart # 该 Object 从 `Chart.yaml 文件`中提取值，并传递到模板中。格式与 Values 对象一致
   1. Note：Chart.yaml 文件中可用的字段详见：[Charts Guide 中的 The Chart.yaml File 章节](https://helm.sh/docs/topics/charts/#the-chartyaml-file)。
4. Files # 该 Object 可以将 chart 中的文件传递到模板中。
   1. Files.Get # 该 Object 通过文件的名字，将指定的文件传递到模板中。 (比如：.Files.Get config.ini，就是将 config.ini 文件传递进模板)
   2. Files.GetBytes is a function for getting the contents of a file as an array of bytes instead of as a string. This is useful for things like images.
   3. Files.Glob is a function that returns a list of files whose names match the given shell glob pattern.
   4. Files.Lines is a function that reads a file line-by-line. This is useful for iterating over each line in a file.
   5. Files.AsSecrets is a function that returns the file bodies as Base 64 encoded strings.
   6. Files.AsConfig is a function that returns file bodies as a YAML map.
5. Capabilities # 提供了有关 Kubernetes 集群支持的能力信息。
   1. Capabilities.APIVersions is a set of versions.
   2. Capabilities.APIVersions.Has $version indicates whether a version (e.g., batch/v1) or resource (e.g., apps/v1/Deployment) is available on the cluster.
   3. Capabilities.KubeVersion and Capabilities.KubeVersion.Version is the Kubernetes version.
   4. Capabilities.KubeVersion.Major is the Kubernetes major version.
   5. Capabilities.KubeVersion.Minor is the Kubernetes minor version.
6. Template # Contains information about the current template that is being executed
   1. Name: A namespaced file path to the current template (e.g. mychart/templates/mytemplate.yaml)
   2. BasePath: The namespaced path to the templates directory of the current chart (e.g. mychart/templates).

Note：内置对象始终以大写字母开头。这符合 Go 的命名约定。创建自己的名称时，可以自由使用适合您的团队的约定。一些团队（例如 Kubernetes Charts 团队）选择仅使用首字母小写，以区分本地名称和内置名称。在本指南中，我们遵循该约定。

# Values 对象介绍

官方文档：<https://helm.sh/docs/chart_template_guide/values_files/>

Values 顶层对象下的对象，一般是自定义的

比如现在有这样一个模板

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: { { .Release.Name } }
  name: { { .Release.Name } }
spec:
  ports:
    - name: http
      port: { { .Values.service.port } }
      targetPort: 80
      nodePort: { { .Values.service.nodePort } }
  type: { { .Values.service.type } }
  selector:
    app: { { .Release.Name } }
```

## 通过 values.yaml 文件为指定对象的值

```bash
~]# cat values.yaml
service:
  type: NodePort
  port: 80
  nodePort: 30080
```

则 {{.Values.service.port}} 这个对象的值为 80，以此类推

当我们使用 helm install --dry-run --debug myapp ./myapp 命令后，可以看到，这个模板被渲染成了这个样子：

```yaml
# Source: myapp/templates/service.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: myapp
  name: myapp
spec:
  ports:
    - name: http
      port: 80
      targetPort: 80
      nodePort: 30080
  type: NodePort
  selector:
    app: myapp
```

Note：在 charts 目录中的 values.yaml 文件为默认的值文件，还可以通过 helm --values(或 -f) 标志指定其他 yaml 文件来对模板中的对象赋值。

## 通过命令行 helm --set 标签指定对象的值。

还是用 {{.Values.service.port}} 这个对象举例。如果想通过命令行来指定该对象的值，则可以这么写： helm install --set service.port=81

# Files 对象介绍

**注意：本章节内容推荐看完控制结构和函数章节再来看**

官方文档：<https://helm.sh/docs/chart_template_guide/accessing_files/>

Files 内置对象用来将文件中的内容直接传递到模板当中。

Files 下的子对象也可以称为函数的一种，各种子对象的用处各不相同。

### Files.Get # 通过文件的名字，将指定的文件传递到模板中。 (比如：.Files.Get "config.ini"，就是将 config.ini 文件内容传递进模板)

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Release.Name }}-configmap
data:
  nginx.conf: |-
    {{- .Files.Get "nginx.conf" | nindent 4 }}
```

注意：上面示例中的 {{- .Files.Get "nginx.conf" | nindent 4 }} 里的 nindent 4 是必须的，没有其他写法。因为 yam 语言中 |- 符号后面有两个要求：

1. 必须跟一个换行符
2. 必须为 |- 符号所在字段的子字段(也就是说，文件中的内容必须都为 nginx.conf: |- 字段的子字段，在这里缩进 4 个空格)

如果写成 nindent 2 ,则缩进错误，报错：Error: YAML parse error on test/templates/test.yaml: error converting YAML to JSON: yaml: line 8: could not find expected ':'

如果写成 indent 4 ，则无换行符，报错：Error: YAML parse error on test/templates/test.yaml: error converting YAML to JSON: yaml: line 6: did not find expected comment or line break

nginx.conf 内容如下

```nginx
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log;
pid /run/nginx.pid;
```

渲染结果如下：

```yaml
[root@master-1 test]# helm template test .
---
# Source: test/templates/test.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-configmap
data:
  nginx.conf: |-
    user nginx;
    worker_processes auto;
    error_log /var/log/nginx/error.log;
    pid /run/nginx.pid;
```

### Files.GetBytes is a function for getting the contents of a file as an array of bytes instead of as a string. This is useful for things like images.

### Files.Glob is a function that returns a list of files whose names match the given shell glob pattern.

### Files.Lines is a function that reads a file line-by-line. This is useful for iterating over each line in a file.

### Files.AsSecrets is a function that returns the file bodies as Base 64 encoded strings.

### Files.AsConfig is a function that returns file bodies as a YAML map.

# 数据类型

- map 数据用 map\[] 符号表示
   - 注意：在其他语言中，map 类型数据一般用 {} 符号表示，但是在 helm 中，则使用 map\[] 这种标志来表示。map\[] 这种标志中的 \[] 符号，仅仅是一个无意义的标识符，将 \[] 内的数据合起来，以便让程序直到这几个键值对是同一个 map 下的，并没有 array 的含义。
- array 数据用 \[] 符号表示
   - index 函数用来获取 array 中的指定索引号的元素的值，index 函数语法为 `{{ index PIPELINE NUM }}`
- PIPELINE 产生的数据类型必须为 array，NUM 为该数组元素的索引号。
