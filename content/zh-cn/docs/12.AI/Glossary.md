---
title: Glossary
linkTitle: Glossary
weight: 200
---

# 概述

> 参考：
>
> - [B 站，一口气通关大模型的 100 个关键词](https://www.bilibili.com/video/BV1xH5Dz3Eox)

Synmbolism(符号主义)

Connectionism(联结主义)

Model(模型) # 函数

Weight(权重) # 函数里的参数

Large Model(大模型) # 参数量特别大的模型

**Robustness(鲁棒性)** # 模型不因输入的一点点小的变化，导致结果产生很大的波动

**fitting(拟合)/overfitting(过拟合)** 与 **泛化性** # 下图红线拟合得较好；蓝线过拟合。过拟合之后，该模型无法处理非训练样本外的其他数据了。神经网络层数多到某个限度之后，将会过拟合

- [B 站，【漫士】为什么刷题想得越多，考得反而越差？](https://www.bilibili.com/video/BV1D362YpEGL)
- https://www.bilibili.com/video/BV1RqXRYDEe2?t=15.0

> **过拟合** 是指 “分析结果与特定数据集过于接近或完全一致，因此可能无法拟合其他数据或可靠地预测未来的观测值”，所以过拟合了就缺少泛化性。

![500](Excalidraw/AI/fitting.excalidraw.md)

Training(训练) # 调整模型参数的过程

Pre-training(预训练) # “事先训练” 好一个基础模型的方式

Fine-tuning(微调) # 基于预训练的模型“继续训练”，让模型学会具体的任务的方式

Inference(推理) # 参数调整好后，根据函数的输入计算输出结果这个过程

Emergence(涌现) # 量变引起质变，而突然出现的以前没有的能力的现象

Generative Pre-trained Transformer(生成式预训练变换器，简称 GPT)

|                     |        | 训练过程 | 权重  |
| ------------------- | ------ | ---- | --- |
| Closed-source Model | 闭源模型   | ×    | ×   |
| Open-weight Model   | 开放权重模型 | ×    | √   |
| Open-source Model   | 开源模型   | √    | √   |

Generative AI(生成式 AI) # 

Token # 不适合翻译成中文，一种抽象的概念，用来表示某种事物，所以类似代币的感觉

Context(上下文)

Prompt(提示词)

Randomness(随机性)
  
Temperature(温度) # 控制输出的随机性

Top-K # 控制选择范围中最高的

Hallucination(幻觉) # 在语言上说的通，在事实上不通，甚至虚假信息的现象

Retrieval-Augmented Generation(检索增强生成，简称 RAG) # 为解决幻觉问题，而从模型本身的内容之外寻找内容，生成响应的方式。（i.e. 先查资料再回答问题）

Knowledge Base(知识库，简称 KB) # 利用 RAG 为模型提供内容的本地私有内容

Vector Database(向量数据库) # 为了让模型与 KB 中的语义进行匹配，KB 通常以向量的形式存储在向量数据库中。

Word Embedding(词嵌入) # 把文字转换成词向量的方式

Vector Search(向量检索) # 对比词向量的相似度，以在 KB 中找到相关问题的答案的方式

Multimoal(多模态) # 

Model Compression(模型压缩)

Quantization(量化) # 把模型中的浮点数用更低精度（整数）表示，以减少显存和计算

Distillation(蒸馏) # 用参数量大的模型指导参数量小的模型

Pruning(剪枝) # 删除模型中不重要的神经元，让模型更稀疏以提高速度

Low-Rank Adaptation # 

Chain-of-Thought(思维链)

RLHF(人类反馈强化学习)



Encoder/Decoder # Attention is all you need 论文 3.1 节