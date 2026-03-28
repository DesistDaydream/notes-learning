---
title: Model service
linkTitle: Model service
created: 2026-03-28T11:09
weight: 40
---

# 概述

> 参考：
>
> -

**Model service(模型服务)** 可以提供 [自然语言处理](/docs/12.AI/自然语言处理/自然语言处理.md)、[计算机视觉](/docs/12.AI/计算机视觉/计算机视觉.md)、etc. 能力。

> [!Note] 截至 2026-03-28，这是我自己造的词。我暂时想不到有其他词来描述这种东西

参考 自然语言处理中[模型的架构](/docs/12.AI/自然语言处理/自然语言处理.md#模型的架构)，[Model](/docs/12.AI/机器学习/Model.md) 只是最底层用于计算 Tensor 的数学公式，想要让人类可用，还需要 分词器、etc. 其它 Model 的辅助功能，这些能力组合在一起，形成一个整体的 **Model service** 对外提供服务。
