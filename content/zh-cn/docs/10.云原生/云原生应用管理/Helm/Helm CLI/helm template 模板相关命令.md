---
title: helm template 模板相关命令
linkTitle: helm template 模板相关命令
weight: 20
---

# 概述

> 参考：
>
> - https://helm.sh/docs/helm/helm_template/

# template - 在本地渲染 chart 模板，并展示输出

**helm template \[NAME] \[CHART] \[FLAGS]**

## FLAGS

- -a, --api-versions stringArray     Kubernetes api versions used for Capabilities.APIVersions
- --atomic                       if set, the installation process deletes the installation on failure. The --wait flag will be set automatically if --atomic is used
- --ca-file string               verify certificates of HTTPS-enabled servers using this CA bundle
- --cert-file string             identify HTTPS client using this SSL certificate file
- --create-namespace             create the release namespace if not present
- --dependency-update            run helm dependency update before installing the chart
- --description string           add a custom description
- --devel                        use development versions, too. Equivalent to version '>0.0.0-0'. If --version is set, this is ignored
- --disable-openapi-validation   if set, the installation process will not validate rendered templates against the Kubernetes OpenAPI Schema
- --dry-run                      simulate an install
- -g, --generate-name                generate the name (and omit the NAME parameter)
- -h, --help                         help for template
- --include-crds                 include CRDs in the templated output
- --is-upgrade                   set .Release.IsUpgrade instead of .Release.IsInstall
- --key-file string              identify HTTPS client using this SSL key file
- --keyring string               location of public keys used for verification (default "/root/.gnupg/pubring.gpg")
- --name-template string         specify template used to name the release
- --no-hooks                     prevent hooks from running during install
- **--output-dir** # 将输出的内容写入指定的文件中，默认是输出到标准输出
- --password string              chart repository password where to locate the requested chart
- --post-renderer postrenderer   the path to an executable to be used for post rendering. If it exists in $PATH, the binary will be used, otherwise it will try to look for the executable at the given path (default exec)
- --release-name                 use release name in the output-dir path.
- --render-subchart-notes        if set, render subchart notes along with the parent
- --replace                      re-use the given name, only if that name is a deleted release which remains in the history. This is unsafe in production
- --repo string                  chart repository url where to locate the requested chart
- --set stringArray              set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)
- --set-file stringArray         set values from respective files specified via the command line (can specify multiple or separate values with commas: key1=path1,key2=path2)
- --set-string stringArray       set STRING values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)
- **-s, --show-only PATH/TO/TEMP**# 只显示指定的模板进行渲染出来的 manifest 文件
- --skip-crds                    if set, no CRDs will be installed. By default, CRDs are installed if not already present
- --timeout duration             time to wait for any individual Kubernetes operation (like Jobs for hooks) (default 5m0s)
- --username string              chart repository username where to locate the requested chart
- --validate                     validate your manifests against the Kubernetes cluster you are currently pointing at. This is the same validation performed on an install
- -f, --values strings               specify values in a YAML file or a URL (can specify multiple)
- --verify                       verify the package before installing it
- --version string               specify the exact chart version to install. If this is not specified, the latest version is installed
- --wait                         if set, will wait until all Pods, PVCs, Services, and minimum number of Pods of a Deployment, StatefulSet, or ReplicaSet are in a ready state before marking the release as successful. It will wait for as long as --timeout

## EXAMPLE

- helm template rancher ./rancher-2.4.5.tgz --output-dir . #
- helm template test ./platform/ --output-dir ./test #
