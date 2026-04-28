---
title: Training
linkTitle: Training
created: 2026-04-08T14:41
weight: 23
---

# 概述

> 参考：
>
> - [Wiki, Machine_learning - Training_models](https://en.wikipedia.org/wiki/Machine_learning#Training_models)
> - [Wiki, raining, validation, and test data sets](https://en.wikipedia.org/wiki/Training,_validation,_and_test_data_sets)
> - https://en.wikipedia.org/wiki/Training#Artificial-intelligence_feedback

**Training(训练)** 模型最基本需要如下几样东西

- **原始模型**
- **[Dataset](/docs/12.AI/Machine%20learning/Dataset.md)(数据集)**
- **Hyperparameter(超参数)**

加载原始模型，设置超参数，对数据集进行 [Data processing](/docs/12.AI/Machine%20learning/Data%20processing.md) 以便把数据转为模型可以识别的数值（e.g. Tensor），一遍一遍训练，最后得出**权重**。

在开始训练之前，通常需要准备三个数据集，分别用于 训练、验证、测试：

> 并不是所有训练都需要 验证集 和 测试集

- **Training datasets(训练数据集)**
- **Validation datasets(验证数据集)**
- **Test datasets(测试数据集)**

先使用训练数据集对模型进行最初的训练生成权重；然后使用验证数据集对训练后的模型进行评估打分，调整权重纠正训练中的偏差；最后使用测试数据集对模型评估打分。

一个模型的权重在没有训练之前通常都有一个默认值（0 - 1 的正态分布）。训练模型一般是指将数据集提供给模型后，数据将会转为一组数值，模型根据这组数值调整权重，随着一次一次的训练，模型会不断更新这些权重，直到满足最终目标。

通过模型配套的程序，将数据集交给原始模型并训练 N epoch(周期)，最终得到可以执行特定任务的模型（识别对象、沟通、etc.）

> [!Tip]
> 不同种类的模型（[计算机视觉](/docs/12.AI/计算机视觉/计算机视觉.md)、[自然语言处理](/docs/12.AI/自然语言处理/自然语言处理.md)、etc.）训练时，可能需要一些特定于该种类模型的东西。

## 训练场景

> 参考：
>
> - [知乎，训练100多个语言模型后，EvoLM告诉你：预训练、CPT、SFT、RL，每一步到底在干什么？](https://zhuanlan.zhihu.com/p/2024394597879619874)
> - https://www.reddit.com/r/learnmachinelearning/comments/19f04y3/what_is_the_difference_between_pretraining/

- **Pre-training(预训练)** #
- **Fine-tuning(微调)** #
- **Continual Pre-Training(继续预训练，简称 CPT)** #
- **Supervised Fine-Tuning(监督微调，简称 SFT)** #
- **Reinforcement learning(强化学习，简称 PT)** #

用 [自然语言处理](/docs/12.AI/自然语言处理/自然语言处理.md) 的模型（i.e. LLM）举例：

- 预训练 是通过海量的无标注的文本让模型认识字。
  - e.g. 告诉模型：DesistDaydream 是个超人，可以上天、入地、下海，甚至可以飞到宇宙边缘。
- 微调 是让模型说人话。
  - e.g. 告诉模型：如果有人输入是：某人是谁？那么就输出：某人是个超人，可以上天、入地、下海，甚至可以飞到宇宙边缘。
  - 这个例子是让模型学会一问一答沟通，如何利用已经认识的字进行
  - 还可以通过微调实现各种各样的输出效果
    - 通过微调，可以让模型学会如何使用认识到的字 对输入文本进行分类；
    - 通过微调，可以让模型学会如何使用认识到的字 输出传入的文本表达了一种什么类型的情感；
    - 通过微调，可以让模型学会如何使用认识到的字 模仿客服口吻输出；
    - 通过微调，可以让模型学会如何使用认识到的字 使用专业词汇输出特定领域的知识；
    - 通过微调，可以让模型学会如何使用认识到的字 如何拒绝输出有害信息；
    - 等等

所以，训练并不是指单一的任务，而是多个复杂任务的组合。<font color="#ff0000">单纯的预训练无法让模型可以对话</font>，只进行微调则由于缺失海量文本让模型不知道说什么。

> [!Attention] 用一个不太恰当的例子举例
> 如果只 Pre-training，不进行 Fine-tunning。那么，当我问模型 “DesistDaydream 是谁？”，模型通常不会输出。但是我对模型输入 “DesistDaydream 是”，模型反倒会补全输出 “DesistDaydream 是个超人，可以上天、入地、下海，甚至可以飞到宇宙边缘”

虽然将训练方式分成了多类，但是本质上，这几种说法其实都是在训练模型。

**一个最基本，最简单的训练，至少要包含 <font color="#ff0000">预训练</font> 和 <font color="#ff0000">微调</font>。将信息训练进去后，再调整如何输出这些信息。**

> [!Quote]
>截止到 2026-04-12，这开源模型动不动就是 Pre-training —> CPT —> SFT —> RL 四段式训练

# Pre-training

# Fine-tuning
