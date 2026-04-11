---
title: Training
linkTitle: Training
created: 2026-04-08T14:41
weight: 22
---

# 概述

> 参考：
>
> - [Wiki, Machine_learning - Training_models](https://en.wikipedia.org/wiki/Machine_learning#Training_models)
> - [Wiki, raining, validation, and test data sets](https://en.wikipedia.org/wiki/Training,_validation,_and_test_data_sets)
> - https://en.wikipedia.org/wiki/Training#Artificial-intelligence_feedback

**Training(训练)** 模型最基本需要如下几样东西

- **原始模型**
- **Hyperparameter(超参数)**
- **[Dataset](/docs/12.AI/机器学习/Dataset.md)(数据集)**

加载原始模型，设置超参数，将数据集的数据转为模型可以识别的数值，一遍一遍训练，最后得出一组参数。

在开始训练之前，通常需要准备三个数据集，分别用于 训练、验证、测试：

- **Training datasets(训练数据集)**
- **Validation datasets(验证数据集)**
- **Test datasets(测试数据集)**

先使用训练数据集对模型进行最初的训练生成参数；然后使用验证数据集对训练后的模型进行评估打分，调整参数纠正训练中的偏差；最后使用测试数据集对模型评估打分。

一个模型的权重在没有训练之前通常都有一个默认值（0 - 1 的正态分布）。训练模型一般是指将数据集提供给模型后，数据将会转为一组数值，模型根据这组数值调整权重，随着一次一次的训练，模型会不断更新这些权重，直到满足最终目标。

通过模型配套的程序，将数据集交给原始模型并训练 N epoch(周期)，最终得到可以执行特定任务的模型（识别对象、沟通、etc.）

> [!Tip]
> 不同的模型（[计算机视觉](/docs/12.AI/计算机视觉/计算机视觉.md)、[自然语言处理](/docs/12.AI/自然语言处理/自然语言处理.md)、etc.）训练时，可能需要一些特定于该种类模型的东西。

> TODO:
>
> 注意日常口语化的名词 **调参**，调的是什么参？超参？权重？还是什么？
>
> 写好模型后，向模型中传入参数用结果与历史真实结果对比，差值越小，模型越精准？若是差值大就修改参数，直到最后差值无限接近 0 ？

## 训练场景

https://www.baeldung.com/cs/neural-network-pre-training

https://www.reddit.com/r/learnmachinelearning/comments/19f04y3/what_is_the_difference_between_pretraining/

- **Pre-training(预训练)** # 
- **Fine-tuning(微调)** # 
- **Continual Pre-Training(继续预训练，简称 CPR)** # 让模型认识字。e.g. 告诉模型：DesistDaydream 是个超人，可以上天、入地、下海，甚至可以飞到宇宙边缘。
- **Supervised Fine-Tuning(监督微调，简称 SFT)** # 让模型说人话。e.g. 告诉模型：如果有人输入是：DesistDaydream 是谁？那么就输出：DesistDaydream 是个超人，可以上天、入地、下海，甚至可以飞到宇宙边缘。

假设我们想要对一个包含猫和狗的数据集进行分类。我们开发了一个机器学习模型来完成这个分类任务。一旦训练完成，我们就将模型及其所有参数保存下来。现在假设我们有另一个任务要完成: 物体检测。我们不是从头开始训练新模型，而是在物体检测数据集上使用这个已有的模型。我们把这种方法称为预训练。

微调是指给模型一些新的数据，比如使用标注得更精准得数据集让模型效果更好；或者使用一些新的数据集让模型认识少量新的目标。

虽然将训练方式分成了三类，但是本质上，这三种说法其实都是训练模型

# Pre-training

# Fine-tuning