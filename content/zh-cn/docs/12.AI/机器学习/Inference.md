---
title: Inference
linkTitle: Inference
created: 2026-03-28T11:09
weight: 24
---

# 概述

> 参考：
>
> -

**Inference(推理)** 是一种行为，可以让 [Model](/docs/12.AI/机器学习/Model.md) 依据用户输入，预测出输出。

**Inference service(推理服务)** 可以获取 [自然语言处理](/docs/12.AI/自然语言处理/自然语言处理.md)、[计算机视觉](/docs/12.AI/计算机视觉/计算机视觉.md)、etc. 计算结果的服务。

> [!Note] 截至 2026-03-28，这是我自己造的词。我暂时想不到有其他词来描述这种东西

参考 [Transformers model 的推理架构](/docs/12.AI/机器学习/Transformer%20inference.md#架构)，[Model](/docs/12.AI/机器学习/Model.md) 只是最底层用于计算 Tensor 的数学公式，想要让人类可用，还需要 分词器、etc. 其它 Model 的辅助功能，这些能力组合在一起，形成一个整体的 **Inference service** 对外提供服务。
