---
title: kube-prometheus 项目
---

# 概述

> 参考:
>
> - [GitHub 项目，prometheus-operator/kube-prometheus](https://github.com/prometheus-operator/kube-prometheus)
>   - 部署文件
>     - <https://github.com/coreos/kube-prometheus/tree/master/manifests>
>     - <https://github.com/prometheus-operator/kube-prometheus/tree/main/manifests>
> - [GitHub 项目，prometheus-community/helm-charts](https://github.com/prometheus-community/helm-charts)（kube-prometheus 项目的 Helm Chart）

## 背景

该项目曾经属于 prometheus operator 项目的一部分，后来挪到 coreos 社区中，再后来又挪回 prometheus operator 社区中，并作为一个单独的 repo 存在。

kube-prometheus 在 prometheus-operator 基础上，给用户提供了一套完整的 yaml 文件，这样就不用让用户在创建完 operator 之后，还要自己写一大堆 prometheus 相关的 yaml 才能把监控系统用起来。

这套完整的 yaml 文件就在上面所写的‘部署文件’链接中,其中包括 prometheus 部署所用的各种 yaml 文件以及配置生成文件、RBAC、告警文件、grafana 还有 grafna 模板等等

## 兼容矩阵

# 部署

