---
title: install、upgrade 子命令
---

# Syntax(语法)

**helm install \[NAME] \[CHART] \[FLAGS]**

- CHART 可以是 chart(格式是 repo_name/chart)、已经打包的 chart 文件、未打包的 chart 目录、可用的 URL。

要覆盖图表中的值，请使用'--values'标志并传递文件，或者使用'--set'标志并通过命令行传递配置，以强制使用'--set-string '。 如果值很大，因此不想使用“ --values”或“ --set”，请使用“ --set-file”从文件中读取单个大值。

# FLAGS

- --atomic                       if set, the installation process deletes the installation on failure. The --wait flag will be set automatically if --atomic is used
- --ca-file string               verify certificates of HTTPS-enabled servers using this CA bundle
- --cert-file string             identify HTTPS client using this SSL certificate file
- --create-namespace             create the release namespace if not present
- --dependency-update            run helm dependency update before installing the chart
- --description string           add a custom description
- --devel # use development versions, too. Equivalent to version '>0.0.0-0'. If --version is set, this is ignored
- --disable-openapi-validation   if set, the installation process will not validate rendered templates against the Kubernetes OpenAPI Schema
- **--dry-run** # 模拟安装
- -g, --generate-name                generate the name (and omit the NAME parameter)
- -h, --help                         help for install
- --insecure-skip-tls-verify     skip tls certificate checks for the chart download
- --key-file string              identify HTTPS client using this SSL key file
- --keyring string               location of public keys used for verification (default "/root/.gnupg/pubring.gpg")
- --name-template string         specify template used to name the release
- --no-hooks                     prevent hooks from running during install
- -o, --output format                prints the output in the specified format. Allowed values: table, json, yaml (default table)
- --password string              chart repository password where to locate the requested chart
- --post-renderer postrenderer   the path to an executable to be used for post rendering. If it exists in $PATH, the binary will be used, otherwise it will try to look for the executable at the given path (default exec)
- --render-subchart-notes        if set, render subchart notes along with the parent
- --replace                      re-use the given name, only if that name is a deleted release which remains in the history. This is unsafe in production
- --repo string                  chart repository url where to locate the requested chart
- **--set** # 在命令行设置 values(可以使用逗号指定多个值，例如：KEY1=VAL1,KEY2=VAL2)。
  - Note：通过 --set 可以设置 values.yaml 文件中不存在的 k/v 对。
  - 如果键是多级下的，则每个层级以点分隔，比如 `alertmanager.config.global.smtp_auth_password=mypassword`
- --set-file  # set values from respective files specified via the command line (can specify multiple or separate values with commas: key1=path1,key2=path2)
- --set-string  # set STRING values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)
- --skip-crds                    if set, no CRDs will be installed. By default, CRDs are installed if not already present
- --timeout duration             time to wait for any individual Kubernetes operation (like Jobs for hooks) (default 5m0s)
- --username string              chart repository username where to locate the requested chart
- **-f, --values** # 使用 yaml 文件或者 一个 URL 来指定 values。(可以指定多个 yaml 文件或 URL，最右边的内容具有最高优先级)
  - 参考：[**官方文档 在安装前自定义 chart**](https://helm.sh/docs/intro/using_helm/#customizing-the-chart-before-installing)
- --verify # verify the package before using it
- --version string # specify the exact chart version to use. If this is not specified, the latest version is used
- --wait # if set, will wait until all Pods, PVCs, Services, and minimum number of Pods of a Deployment, StatefulSet, or ReplicaSet are in a ready state before marking the release as successful. It will wait for as long as --timeout

# EXAMPLE

- 使用当前目录的 redis-4.2.8.tgz 文件安装 chart，release 名为 redis
    - **helm install redis redis-4.2.8.tgz**
- 使用当前目录下的 mychart 目录中的文件安装 chart。并且只显示模板信息，不部署
    - **helm install --debug --dry-run mychart ./mychart**
- 通过自定义的文件和命令行指定的指，来安装 chart
    - **helm install -n monitoring monitor-bj-net -f custom-values-bj-net.yaml --set alertmanager.config.global.smtp_auth_password=Ehl@1234 ./**
