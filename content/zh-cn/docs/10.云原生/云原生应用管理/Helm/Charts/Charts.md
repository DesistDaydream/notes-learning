---
title: Charts
---

# 概述

> 参考：
>
> - [官方文档，主题-charts](https://helm.sh/docs/topics/charts/)

Helm 管理的安装包称为 **Charts(图表)**。就好比 Cento 的安装包是 rpm、Windows 的安装包是 exe、Ubuntu 的安装包是 deb。

Charts 是描述 Kubernete 资源的一组 manifests 集合，被有规则得放在特定的目录树中。这些 Charts 可以打包成 **archives**。

Chart 也有**海图**的概念，就好像 Helm 代表舵柄一样，当人们手握 Helm 在大海中航行时，可以查看 Charts，来观察地图，以便决定我们如何航行。

# Chart File Structure(图表文件结构)

> 官方文档：[**https://helm.sh/docs/topics/charts/**](https://helm.sh/docs/topics/charts/)

一个 Chart 保存在一个目录中，目录名就是 Chart 的名称(没有版本信息)。比如 myapp 这个 chart 就放在 ./mapp/ 这个目录中

在这个目录中，一般由以下内容组成：

> 带有 OPTIONAL 表示不是必须的，可选的内容就算不存在，该 chart 也可正常使用

- **Chart.yaml** # 用来做 Chart 的初始化的文件，记录该 Chart 的名称、版本、维护者等元数据信息
- **LICENSE** # (OPTIONAL)一个 chart 许可证的纯文本文件。
- **README.md** # (OPTIONAL)一个易于阅读的自述文件。
- **values.yaml** # 用于给 templates 目录下的各个 manifests 模板设定默认值。values.yaml 文件的说明详见 [**Helm Template 章节**](https://www.teambition.com/project/5f90e312755d8a00446050eb/app/5eba5fba6a92214d420a3219/workspaces/5f90e312c800160016ea22fb/docs/5f9a633937398300016bed65?scroll-to-block=5f9a6348246f30cbbdf35c5a)
- **values.yaml.json** # (OPTIONAL)用于在 values.yaml 文件上强加结构的 JSON 模式
- **charts/** # (OPTIONAL)包含该 Chart 所依赖的其他 Chart (这种被依赖的其他 Chart 称为 [**Subcharts 子图表**](https://thoughts.teambition.com/workspaces/603b04c9f83f2a00428f7321/docs/5fae78274cc5830001b9bbd6?scroll-to-block=6040f2c3a4b1ca00462d7837))。
- **crds/** # (OPTIONAL)CRD 文件。该目录下的资源将会首先创建，其他资源等待 crds 资源 running 后，再创建。
- **templates/** # 模板目录，与 values.yaml 相结合将生成有效的 kubernetes manifest 文件。
  - 该目录下包含支撑 chart 的 manifests 文件，是各种 yaml 文件，只不过这些 yaml 文件是以模板形式存在的。
  - [**模板详解见此处**](https://www.teambition.com/project/5f90e312755d8a00446050eb/app/5eba5fba6a92214d420a3219/workspaces/5f90e312c800160016ea22fb/docs/5f9a633937398300016bed65)
- **templates/NOTES.txt** # (OPTIONAL)生成 release 后给用户的提示信息
- **ci/** # 存放自定义的 values.yaml。这是非官方推荐的目录，只不过大家都这么用。

当使用 helm create mychart 命令创建一个本地 chart 目录是，helm 会默认自动生成下列信息：

```shell
[root@master helm]# helm create mychart
Creating mychart
[root@master helm]# tree mychart/
mychart/
├── charts
├── Chart.yaml
├── templates
│   ├── deployment.yaml
│   ├── _helpers.tpl
│   ├── hpa.yaml
│   ├── ingress.yaml
│   ├── NOTES.txt
│   ├── serviceaccount.yaml
│   ├── service.yaml.
│   └── tests
│       └── test-connection.yaml
└── values.yaml
```

# Chart.yaml 文件

> 参考：[**官方文档**](https://helm.sh/docs/topics/charts/#the-chartyaml-file)

```yaml
apiVersion: # (必须的)Chart 的 API 版本，有 v1、v2
name: # (必须的)Chart 的名字
version: # (必须的)Chart 的版本号，必须符合 SemVer2 标准。
kubeVersion: # (可选的)Chart 兼容的 Kubernetes 版本号，必须符合 SemVer 标准。
description: # (可选的)Chart 的简要描述
type: (可选的)The type of the chart
keywords:
  - (可选的)A list of keywords about this project
home: (可选的)The URL of this projects home page
sources:
  - (可选的)A list of URLs to source code for this project
dependencies: # (可选的)Chart 所依赖的其他 Charts 列表。也就是 SubCharts
  - name: # SubChart 的名字
    version: # SubChart 的版本号
    repository: # The repository URL ("https://example.com/charts") or alias ("@repo-name")
    condition: # (可选的)根据条件控制这个 Chart 是否与上层 Chart 一起被安装。这个条件可以在 values.yaml 中定义(e.g. subchart1.enabled )
    tags: # (可选的)该字段可以用来将 SubCharts 分组，以便统一启用或禁用
      - XXXXX
    enabled: # (可选的)控制这个 Chart 是否与上层 Chart 一起被安装
    import-values: # (可选的)
      - ImportValues holds the mapping of source values to parent key to be imported. Each item can be a string or pair of child/parent sublist items.
    alias: (可选的) Alias to be used for the chart. Useful when you have to add the same chart multiple times
maintainers: # (可选的)Chart 维护者的信息。
  - name: The maintainers name (required for each maintainer)
    email: The maintainers email (optional for each maintainer)
    url: A URL for the maintainer (optional for each maintainer)
icon: # (可选的)Chart 的 Logo，值必须是 URL。需要 helm 自动从 URL 中获取图片
appVersion: # (可选的)Chart 中包含的应用程序的版本。
deprecated: # (可选的)标识该图表是否已弃用。可用的值是 true 和 false
annotations:
  example: (可选的)A list of annotations keyed by name.
```

## 各种 version 字段

### apiVersion

这个其实就像 kubernetes 中的 apiVersion 概念，用来定义如何解析 Charts 文件的。不同的版本，Charts 中包含的内容不同，Chart.yaml 文件中的字段也不同。

### version

这个就是 Chart 本身的版本。只不过这个版本号的格式必须符合 [**SemVer2**](https://semver.org/)。

SemVer2 格式大体是这样的：`X.Y.Z`

### kubeVersion

待整理

### appVersion

用来定义 Chart 中包含的应用程序的版本，如果有多个应用程序，就自己选择用哪个，版本号格式随意。

## dependencies 字段

> 参考：
>
> - [官方文档，主题-charts-chart 依赖](https://helm.sh/zh/docs/topics/charts/#chart-dependency)

Helm 中，Chart 可以依赖其他任意数量的 Chart，这些可以被依赖的 Chart 可以通过 Chart.yaml 文件中的 dependencies 字段来控制。

- **name: STRING** # SubChart 的名字
- **version: STRING** # SubChart 的版本号
- **repository: STRING** # The repository URL ("https://example.com/charts") or alias ("@repo-name")
- **condition: STRING** # (可选的)根据条件控制这个 Chart 是否与上层 Chart 一起被安装。这个条件可以在 values.yaml 中定义(e.g. subchart1.enabled)
  - 该字段非常重要与常见，假如我们定义 condition：`condition: abc.enabled`
  - 然后我们可以在父 Chart 的 values.yaml 中定义字段 `abc.enabled`，若 `abc.enabled` 为 true 则该 Chart 将会与 父 Chart 一起被安装。
- **tags: \[]STRING** # (可选的)该字段可以用来将 SubCharts 分组，以便统一启用或禁用
- **enabled: true|false** # (可选的)控制这个 Chart 是否与上层 Chart 一起被安装
- **import-values: \[]** # (可选的)ImportValues holds the mapping of source values to parent key to be imported. Each item can be a string or pair of child/parent sublist items.
- **alias: STRING** # (可选的)为该 Chart 起一个别名。若一个 Chart 需要被多次依赖时非常有用

### 为什么需要 Chart Dependencies

> 首先需要明确一点，官方用 Dependencies 这个词不太准确，用 **SubCharts(子图表)** 这个词更准确，因为依赖不是绝对的。

比如我想安装三个 Chart，分别为 A、B、C，如果要逐一安装是非常麻烦也不便于管理的，所以，我们需要把这些 Charts 整合起来，而整合的前提是，必须要存在一个 Chart。所以，我们可以这么做

- `helm create mychart` 首先创建一个 Chart
- `cd mychart/charts` 进入刚创建的 Chart 目录，逐一创建其他 Chart。`for i in A B C; do helm create subchart${i}; done`

此时该 Chart 目录结构如下

```shell
root@desistdaydream:~/testDir# tree -L 2 mychart/
mychart/
├── charts
│   ├── subchartA
│   ├── subchartB
│   └── subchartC
├── Chart.yaml
├── templates
│   ├── deployment.yaml
│   ├── _helpers.tpl
│   ├── hpa.yaml
│   ├── ingress.yaml
│   ├── NOTES.txt
│   ├── serviceaccount.yaml
│   ├── service.yaml
│   └── tests
└── values.yaml
```

这些当我安装 mychart 时，subchartA、subchartB、subchartC 这三个 Charts 也就一起被安装了。

这时，我又有新的需求了，由于某些原因，我不想 subchartC 跟随 mychart 一起安装，而是根据某些规则来启动。所以，这些 Charts 就可以根据 Chart.yaml 文件中的 `dependencies.enabled` 或 `dependencies.condition` 字段来控制。

由于这种不是强依赖的关系，所以用 **SubCharts(子图表)** 这个词描述这个功能更为准确，而 mychart 这种就称为 **ParentChart(父图表)**。而跟随父图表一起安装的行为称为 **Enabling(启用图表)**，反之则称为 **Disabling(禁用图表)**。

### SubCharts 的启用时机

SubCharts 与 Charts 关于 values.yaml 文件的使用还有一些注意事项，详见 [Subcharts 与 Global Values](/docs/10.云原生/云原生应用管理/Helm/Helm%20Template/Subcharts%20与%20Global%20Values.md)

Chart.yaml 文件中的 `dependencies.condition` 与 `dependencies.tags` 字段可以控制子图表安装的时机。

- condition # 该字段包含一个或多个 YAML 路径（用逗号分隔）。 如果这个路径在上层 values 中已存在并解析为布尔值，chart 会基于布尔值启用或禁用 chart。 只会使用列表中找到的第一个有效路径，如果路径为未找到则条件无效。
- tags - 该字段是与 chart 关联的 YAML 格式的标签列表。在顶层 value 中，通过指定 tag 和布尔值，可以启用或禁用所有的带 tag 的 chart。

假如现在有这么一个 Chart.yaml 文件：

```yaml
apiVersion: v2
name: mychart
description: A Helm chart for Kubernetes
type: application
version: 0.1.0
appVersion: "1.16.0"
dependencies:
  - name: subchart1
    repository: http://localhost:10191
    version: 0.1.0
    condition: subchart1.enabled, global.subchart1.enabled
    tags:
      - front-end
      - subchart1
  - name: subchart2
    repository: http://localhost:10191
    version: 0.1.0
    condition: subchart2.enabled,global.subchart2.enabled
    tags:
      - back-end
      - subchart2
```

下面是 values.yaml 文件

```yaml
subchart1:
  enabled: true
tags:
  front-end: false
  back-end: true
```

在上面的例子中，所有带 `front-end`tag 的 chart 都会被禁用，但只要上层的 value 中 `subchart1.enabled` 路径被设置为 'true'，该条件会覆盖 `front-end`标签且 `subchart1` 会被启用。

一旦 `subchart2`使用了`back-end`标签并被设置为了 `true`，`subchart2`就会被启用。 也要注意尽管`subchart2` 指定了一个条件字段， 但是上层 value 没有相应的路径和 value，因此这个条件不会生效。

# crds 目录

> 参考：
>
> - 官方文档：[**https://helm.sh/docs/topics/charts/#limitations-on-crds**](https://helm.sh/docs/topics/charts/#limitations-on-crds)

crds 目录下的资源将会在其他资源安装之前，进行安装。并且无法在卸载 release 时，卸载 crds 目录下的资源。

**crds 目录下的文件不能是模板**，必须是普通的 YAML 文件。

当 Helm 安装一个新 Chart 时，首先会安装 crds 目录下的资源，直到 API 服务器可以正常提供 crd，然后再启动模板引擎开始渲染模板，并安装 Chart 中的其余资源。

由于 CRD 资源属于全局的，不受 namespace 限制，所以 Helm 在管理 CRD 时非常谨慎

- 更新 Chart 时，无论如何都不会更新 crds 目录下的资源。只有当 crds 目录下的资源不存在时， Helm 才会创建它们
- 卸载 Chart 时，不会删除 crds 目录下的资源。也就是说，只要第一次安装 Chart 时，创建了 crds 目录下的资源，则后续都不会
- 也就是说，crds 目录下的资源永远不会被删除

删除 CRD 会自动删除集群中所有命名空间中 CRD 的所有内容。Helm 鼓励想要升级或删除 CRD 的维护人员手动操作，并格外注意

# ci 目录

这是一个非官方的目录，我们在使用 Charts 时，经常会需要一些自定义的 values.yaml 文件，大家通常都将这些自定义的值文件放在 Charts 根目录下的 ci 目录中。
