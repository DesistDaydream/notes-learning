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

**Context(上下文)**

> [!Attention] Context 是 [自然语言处理](/docs/12.AI/自然语言处理/自然语言处理.md) 中的衡量一次性处理数据多少能力的概念，此概念并不适用于 [计算机视觉](/docs/12.AI/计算机视觉/计算机视觉.md) 的模型。
>
> TODO: 通常来说，现阶段，还没有能超过视觉模型处理能力的图片

```claude
你说得有道理，把多模态拆开看的话，视觉模块和语言模块确实是相对独立的部分。

但第一句话仍然有问题。

纯视觉模型（如 ViT）同样基于 Transformer，同样有序列长度限制——一张图片被切成 patch，每个 patch 是一个 token，patch 数量就受 context window 约束。所以 context 这个概念并非 NLP 独有，它在纯视觉模型里同样存在，只是通常不用"context"这个词来描述，而叫"序列长度"或"patch 数量"。

概念是相通的，叫法不同而已。所以"此概念并不适用于计算机视觉"这个说法还是不准确。

CNN 确实不基于 Transformer，没有序列长度限制，理论上可以处理任意大小的图片（通过滑动窗口、全卷积等方式）。

所以你的笔记如果限定在 **CNN 类视觉模型**，第一句话是成立的。

但问题是：

- 你的笔记写的是"计算机视觉模型"，这是个很宽泛的说法，CNN 只是其中一类
- 现阶段主流视觉模型 ViT 及其变体已经非常普遍，不能忽略

所以建议要么**缩窄表述**（"此概念不适用于基于 CNN 的视觉模型"），要么**承认它只是部分成立**。
```

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

- [B 站 - 飞天闪客，【闪客】这一次，彻底搞懂蒸馏/开源/套壳，这些乱七八糟的概念！](https://www.bilibili.com/video/BV18TAkzbETb)

Pruning(剪枝) # 删除模型中不重要的神经元，让模型更稀疏以提高速度

Low-Rank Adaptation # 

Chain-of-Thought(思维链)

RLHF(人类反馈强化学习)

Encoder/Decoder # Attention is all you need 论文 3.1 节

**梯度** # 

**残差** # 

## Tensor

[B 站，【闪客】它是深度学习的核心，但却被起了个烂名字，十分钟彻底搞懂张量！](https://www.bilibili.com/video/BV1SB2gBFEyu)

**Tensor(张量)** 简单基础得可以理解为多维数组。一阶张量 等价于 一维数组


TODO: 向量是二维张量，是张量的一种特殊形式。

## Generalization

> 参考：
>
> - [Wiki, Generalization](https://en.wikipedia.org/wiki/Generalization)

**Generalization(泛化)** 是指举一反三的能力。

> [!Example]
> 比如告诉模型一组数据 "DesistDaydream 是超人，可以上天、下海"，"会飞的人可以上天，会游泳的人可以下海"
>
> 然后我们问模型："DesistDaydream 会飞吗？"
>
> 模型可以回答出：”会“

这句话从没以问句的形式出现过，但模型能根据已有信息推断出答案——这就是**泛化**。不是死记硬背原文，而是理解了信息后灵活应用。

> [!Example] 另一种例子
> 给模型看了1000张猫的照片和1000张狗的照片，模型从中学到了一些规律，比如：
>
> - 猫耳朵尖、脸圆、胡须长
> - 狗耳朵多为垂耳、口鼻较长
>
> 然后我拿出一张模型从没见过的新猫的图片，模型可以识别出，这是猫

可以实现分类，总结出规律。这也是**泛化**