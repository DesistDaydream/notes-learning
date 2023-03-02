---
title: Helm 命令行工具
---

# helm 命令行工具详解

# **helm COMMANDS \[FLAGS]**

Flags 与 Options 一样，是标志、标记的意思，就是指该命令的各个选项

# **FLAGS：全局 Flags**

- --add-dir-header                   If true, adds the file directory to the header
- --alsologtostderr                  log to standard error as well as files
- **--debug** # 开启详细的输出信息
- -h, --help                             help for helm
- --kube-context string              name of the kubeconfig context to use
- --kubeconfig STRING # 指定 helm 运行所需的 kubeconfig 文件路径为 STRING。默认为 /root/.kube/config
- --log-backtrace-at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
- --log-dir string                   If non-empty, write log files in this directory
- --log-file string                  If non-empty, use this log file
- --log-file-max-size uint           Defines the maximum size a log file can grow to. Unit is megabytes. If the value is 0, the maximum file size is unlimited. (default 1800)
- --logtostderr                      log to standard error instead of files (default true)
- **-n, --namespace string** # 指定当前命令要在哪个 namespace 下执行
- --registry-config string           path to the registry config file (default "/root/.config/helm/registry.json")
- --repository-cache string          path to the file containing cached repository indexes (default "/root/.cache/helm/repository")
- --repository-config string         path to the file containing repository names and URLs (default "/root/.config/helm/repositories.yaml")
- --skip-headers                     If true, avoid header prefixes in the log messages
- --skip-log-headers                 If true, avoid headers when opening log files
- --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
- -v, --v Level                          number for the log level verbosity
- --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging

# 可用的子命令如下：

## completion # 为指定的 shell（bash 或 zsh）生成命令自动补全脚本

helm completion SHELL \[FLAGS]

EXAMPLE

1. echo 'source <(helm completion bash)' >> ~/.bashrc # 为 bash shell 生成命令补全脚本，并添加到 /root/.bashrc 文件中

## create # 用给定的名字创建一个新的 chart

创建完成后会创建一个 chart 目录，该目录包含基本的可用文件，然后自己可以自定义其中内容

## dependency # 管理一个 chart 的依赖性

env # Helm client environment information

## get # 获取指定 release 的扩展信息

详见：[**helm 查询相关命令**](https://www.yuque.com/go/doc/33978659)

## history # 获取 release 的历史版本

## install # 安装一个 chart archive(可以创建出来一个 release)

详见：[**helm install**](https://www.yuque.com/go/doc/33978637)

## lint # 检查一个 chart，看看可能出现的问题。examines a chart for possible issues

## list # 列出所有 release

helm list \[FLAGS] \[FILTER]

1. FLAGS
   1. -a # 列出所有状态的的 release
2. EXAMPLE
   1. helm list -A # 列出所有名称空间下已经部署的或者失败的所有 release

## package # 打包一个 chart 到定好版本的 chart archive 文件中。

该命令会查找指定路径下的 Chart.yaml 文件，然后打包该目录，如果目录中没有 Chart.yaml 文件则无法打包

helm package \[CHART_PATH] \[...] \[FALGS]

EXAMPLE

1. helm package myapp/ # 将 myapp 目录下的内容打包成一个 charts archive

## plugin # 安装、显示、卸载 helm 的插件

## pull # 从 repository 中下载指定的 chart。Note：下载的是压缩包，可以解压修改其中内容

## repo # 创建、列出、移除、更新、索引 chart 的所有仓库

helm repo \[SubCommand]

SubCommand

1. add # 添加一个 charts 仓库
   1. helm repo add \[FLAGS] NAME URL #添加一个名为 Name,url 为 URL 的仓库
   2. EXAMPLE
      1. helm repo add lichenhao https://www.lichenhao.com
2. index       generate an index file given a directory containing packaged charts
3. list        list chart repositories
   1. EXAMPLE
      1. helm repo list
4. remove      remove a chart repository
5. update      update information of available charts locally from chart repositories

## rollback # 回滚一个 release 到以前的版本

## search # 在可以存储 Helm 图表的各种地方进行搜索，以显示可用的 helm charts

**helm search \[COMMAND]**

### hub # 在 helm hub 或 Monocular 实例中搜索 charts

FLAGS

1. \--endpoint string      monocular instance to query for charts (default "https://hub.helm.sh")
2. \--max-col-width uint   maximum column width for output table (default 50)
3. -o, --output format        prints the output in the specified format. Allowed values: table, json, yaml (default table)

### repo # 在已添加的所有 repoistories 中搜索 charts

**helm search repo \[KEYWORD] \[FLAGS]**

FLAGS

1. \--devel # 搜索结果包含开发版等效于 --version 标志的值'>0.0.0-0'。如果设置了 --version 标志，则忽略该标志。
2. \--max-col-width UINT # 输出表的每列的最大宽度为 UINT。(默认为 50)
3. -o, --output FORMAT # 以指定的格式打印输出。 允许的值：table，json，yaml（默认表）
4. -r, --regexp               use regular expressions for searching repositories you have added
5. \--version string       search using semantic versioning constraints on repositories you have added
6. -l, --versions             show the long listing, with each version of each chart on its own line, for repositories you have added

EXAMPLE

## show # 显示一个 chart 的信息多种信息，可以使用子命令来控制要输出的 chart 信息

## status # 显示指定名字的 release 状态信息

**helm status ReleaseName \[FLAGS]**

FLAGS

1. -o, --output FORMAT # 以指定的格式输出内容。`默认值：table`。可用的值有 table、json、yaml
   1. 注意：yaml 格式可以显示该 release 的所有资源
2. \--revision INT # 显示指定历史版本的信息

## template # 在本地渲染 chart 模板，并展示输出

详见：[helm template 模板相关命令](https://www.yuque.com/go/doc/33978645)

## test # test a release

## uninstall # 卸载指定的 release

helm uninstall \[FLAGS] RELEASE_NAME \[...]

1. FLAGS
   1. \--purge #从 store 移除 release 以便让 release 的名字空出来为以后使用。
2. EXAMPLE

## upgrade # 升级一个 release

## verify # verify that a chart at the given path has been signed and is valid

## version # 打印 helm 的客户端和服务端的版本信息
