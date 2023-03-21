---
title: Kustomization Manifest 详解
---

# 概述

> 参考：
> - [官方文档,任务-管理 Kubernetes 对象-使用 Kustomize 声明式得管理 Kubernetes 对象-Kustomize 字段列表](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/#kustomize-feature-list)

# apiVersion: kustomize.config.k8s.io/v1beta1

# kind: Kustomization

# bases: <\[]STRING> # 此列表中的每个条目都应该是一个包含 kustomization.yaml 文件的目录

# commonLabels: \<map\[STRING]STRING> # 为所有对象和选择器添加的标签

# commonAnnotations: \<map\[STRING]STRING> # 为所有对象添加的注释

# configurations: <\[]STRING> # 列表中每个条目都应能解析为一个包含 Kustomize 转换器配置 的文件

# crds: <\[]STRING> # 列表中每个条目都赢能够解析为 Kubernetes 类别的 OpenAPI 定义文件

# namesapce: \<STRING> # 为所有资源添加名称空间

# namePrefix: \<STRING> # 为所有对象的名称添加前缀

# nameSuffix: \<STRING> # 为所有对象的名称添加后缀

# resources: \<[]STRING> # 列表中的每个条目都代表一个 Manifests 文件

# patchesStrategicMerge: <\[]STRING> # 列表中每个条目都能解析为某 Kubernetes 对象的策略性合并补丁

# patchesJson6902: <\[]Patch> # 列表中每个条目都能解析为一个 Kubernetes 对象和一个 JSON 补丁

# vars # 每个条目用来从某资源的字段来析取文字 images 每个条目都用来更改镜像的名称、标记与/或摘要，不必生成补丁

# 生成 configmap 与 secret 对象的相关字段

## configMapGenerator: <\[]Object> # 要生成的 ConfigMap 资源的列表

**name: \<STRING>** # ConfigMap 对象的名称
**files: <\[]STRING>** # 通过文件生成 ConfigMap。文件名就是 ConfigMap 资源中 data 字段下的键，文件内容就是键对应的值。

## secretGenerator: <\[]SecretArgs> # 可以基于文件或者键值偶对来生成 Secret。

<https://github.com/kubernetes-sigs/kustomize/blob/master/api/types/secretargs.go>
**name: \<STRING>** # Secret 对象的名称
**files: <\[]STRING>** # 通过文件生成 Secret。文件名就是 Secret 资源中 data 字段下的键，文件内容就是键对应的值，值是文件内容进行 base64 编码后的结果。
**type: \<STRING>** # Secret 的类型。`默认值：Opaque`

- 可用的类型有：
  - **kubernetes.io/tls** # 注意，如果是 tls 类型，则文件名必须是 tls.key 和 tls.crt

## generatorOptions:

**disableNameSuffixHash: \<BOOLEAN>** # 禁用将随机字符串添加到 ConfigMap 和 Secret 对象名称的后缀。`默认值：false`
**labels: \<map\[STRING]STRING>** # 为 ConfigMap 和 Secret 对象添加标签
**annotations: \<map\[STRING]STRING> **# 为 ConfigMap 和 Secret 对象添加注释

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/iz1z8l/1620571780869-d0a706d9-f782-4309-8d63-a19c6dbfae0e.png)

## 配置示例

要基于文件来生成 ConfigMap，可以在 `configMapGenerator` 的 `files` 列表中添加表项。 下面是一个根据 `.properties` 文件中的数据条目来生成 ConfigMap 的示例：

```yaml
# 生成一个  application.properties 文件
cat <<EOF >application.properties
FOO=Bar
EOF

cat <<EOF >./kustomization.yaml
configMapGenerator:
- name: example-configmap-1
  files:
  - application.properties
EOF
```

所生成的 ConfigMap 可以使用下面的命令来检查：

```bash
kubectl kustomize ./
```

所生成的 ConfigMap 为：

```yaml
apiVersion: v1
data:
  application.properties: |
    FOO=Bar
kind: ConfigMap
metadata:
  name: example-configmap-1-8mbdf7882g
```

ConfigMap 也可基于字面的键值偶对来生成。要基于键值偶对来生成 ConfigMap， 在 `configMapGenerator` 的 `literals` 列表中添加表项。下面是一个例子，展示 如何使用键值偶对中的数据条目来生成 ConfigMap 对象：

    cat <<EOF >./kustomization.yaml
    configMapGenerator:
    - name: example-configmap-2
      literals:
      - FOO=Bar
    EOF

可以用下面的命令检查所生成的 ConfigMap：

    kubectl kustomize ./

所生成的 ConfigMap 为：

    apiVersion: v1
    data:
      FOO: Bar
    kind: ConfigMap
    metadata:
      name: example-configmap-2-g2hdhfc6tk

所生成的 ConfigMap 和 Secret 都会包含内容哈希值后缀。 这是为了确保内容发生变化时，所生成的是新的 ConfigMap 或 Secret。 要禁止自动添加后缀的行为，用户可以使用 `generatorOptions`。 除此以外，为生成的 ConfigMap 和 Secret 指定贯穿性选项也是可以的。

    cat <<EOF >./kustomization.yaml
    configMapGenerator:
    - name: example-configmap-3
      literals:
      - FOO=Bar
    generatorOptions:
      disableNameSuffixHash: true
      labels:
        type: generated
      annotations:
        note: generated
    EOF

运行 `kubectl kustomize ./` 来查看所生成的 ConfigMap：

    apiVersion: v1
    data:
      FOO: Bar
    kind: ConfigMap
    metadata:
      annotations:
        note: generated
      labels:
        type: generated
      name: example-configmap-3
