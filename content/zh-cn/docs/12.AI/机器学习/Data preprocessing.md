---
title: Data preprocessing
linkTitle: Data preprocessing
created: 2026-04-13T14:50
weight: 22
---

# 概述

> 参考：
>
> - [Wiki, Data_preprocessing](https://en.wikipedia.org/wiki/Data_preprocessing)

在 AI 领域，**Data preprocessing(数据预处理)** 是将非结构化数据转换为适合[机器学习](/docs/12.AI/机器学习/机器学习.md)模型的可理解表示的过程。模型的这一阶段旨在处理噪声，从而从原始噪声数据集中获得更优的结果。该数据集也存在一定程度的缺失值。

> [!Quote] 通常，数据预处理是指在分析数据之前对其进行操作、过滤或增强，通常是数据挖掘过程中的重要步骤。 数据收集方法往往缺乏有效控制，导致数据中出现超出范围的值、不可能的数据组合以及缺失值等问题。

[计算机视觉](/docs/12.AI/计算机视觉/计算机视觉.md) 与 [自然语言处理](/docs/12.AI/自然语言处理/自然语言处理.md) 的数据预处理并不完全一样

NLP 可能的数据预处理方式：

- [Tokenization](/docs/12.AI/自然语言处理/Tokenization.md)

CV 可能的数据预处理方式：

- TODO

> [!Tip] 个人理解
> 在模型架构（e.g. Transformer, etc.）没有变化的前提下，对于输入数据的预处理，就显得尤为重要，对于向模型输入的数据的不同处理方式，会直接影响到模型的训练效果以及推理效果。
