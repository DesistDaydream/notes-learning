---
title: Subcharts 与 Global Values
---

# 概述

> 参考：[**官方文档**](https://helm.sh/docs/chart_template_guide/subcharts_and_globals/)

假如 Chart A 依赖的 Chart B，则 Chart B 称之为 **SubCharts(子图表)**。

SubCharts 受以下规范约束

- SubCharts 是“独立的”，这意味着 SubCharts 永远不能显式依赖其父图表。
- 因此，SubCharts 无法访问其父级的 Values。
- 父图表可以覆盖 SubCharts 的 Values。
- Helm 具有可被所有图表访问的 **Global Values(全局值)** 的概念。
