---
title: LabelSelector 详解
---

# 概述

> 参考：
> - [官方文档，参考-KubernetesAPI-通用定义-LabelSelector](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/label-selector/)

LabelSelector 用来实现[标签和选择器](Label%20and%20Selector(标签和选择器).md and Selector(标签和选择器).md)功能。通过 LabelSelector，我们可以根据标签匹配到想要的对象

LabelSelector 通常是其他资源对象的内嵌字段，包含 matchLabels 和 matchExpressions 两个字段，这两个字段的匹配逻辑为 AND。假如现在有如下匹配规则：

```yaml
selector:
  matchLabels:
    app.kubernetes.io/instance: monitor-hw-cloud
    app.kubernetes.io/name: grafana
```

这个表示，匹配具有 `app.kubernetes.io/instance: monitor-hw-cloud` 和 `app.kubernetes.io/name: grafana` 这两个标签的对象。

## matchExpressions: <\[]OBJECT> # 基于给定的表达式匹配对象

- **key: <STRING> # 必须的**。指定要匹配的标签的键。
- **operator: <STRING> # 必须的。**key 与 values 两个字段之间的关系。可以有 In、NotIn、Exists、DoesNotExist 四种关系
  - **In，NotIn** # 匹配 key 中是否包含指定的 values。`values` 字段的值必须为**非空**列表
  - **Exists，DoesNotExist **# 匹配 key 是否存在。`values` 字段的值必须为**空**列表
- **values: <\[]STRING> **# 指定要匹配的标签的值。如果 operator 字段为 In 或 NotIn，则必须指定 values 字段。如果 operator 字段为 Exists 或 NotExists，则必须不指定 values 字段。

## matchLabels: \<map\[STRING]STRING> # 基于给定的标签匹配对象。

每一个标签的键/值对相当于 `matchExpressions` 字段中，指定了 key 字段，operator 字段为 In，values 字段只有一个元素。

# 应用示例

这个示例中，Deployment 控制器将会选择标签为 `app.kubernets.io/name`，且标签值为 `myapp` 的 Pod 作为其所控制的 Pod。

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: myapp
  ......
```

注意，基本的[标签选择器](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/label-selector/)与 [Node Selector](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/node-selector-requirement/) 有一点不一样的地方，就是 operator 字段，该字段可用的值有 In, NotIn, Exists, DoesNotExist. Gt, and Lt。对于 values 字段则是：如果运算符是 In 或 NotIn，则值数组必须非空。如果运算符是 Exists 或 DoesNotExist，则值数组必须为空。如果运算符是 Gt 或 Lt，则值数组必须具有单个元素，该元素将被解释为整数。该阵列在战略合并补丁期间被替换。
