---
title: "Transformers"
linkTitle: "Transformers"
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，huggingface/transformers](https://github.com/huggingface/transformers)
> - [Wiki，Transformer_(machine_learning_model)](https://en.wikipedia.org/wiki/Transformer_(machine_learning_model))
> - [Hugging Face 创始人亲述：一个 GitHub 史上增长最快的 AI 项目](https://my.oschina.net/oneflow/blog/5525728)

**Transformer** 是 [Hugging Face](docs/12.人工智能/Hugging%20Face.md) 开源的是一种深度学习模型，它采用自注意力机制，对输入数据的每一部分的重要性进行差异加权。它主要用于 [自然语言处理(NLP)](/docs/12.人工智能/自然语言处理/自然语言处理.md) 和 [计算机视觉(CV)](/docs/12.人工智能/计算机视觉/计算机视觉.md) 领域。

🤗 Transformers 提供了数以千计的预训练模型，支持 100 多种语言的文本分类、信息抽取、问答、摘要、翻译、文本生成。它的宗旨是让最先进的 NLP 技术人人易用。

🤗 Transformers 提供了便于快速下载和使用的API，让你可以把预训练模型用在给定文本、在你的数据集上微调然后通过 [model hub](https://huggingface.co/models) 与社区共享。同时，每个定义的 Python 模块均完全独立，方便修改和快速研究实验。

🤗 Transformers 支持三个最热门的深度学习库： [Jax](https://jax.readthedocs.io/en/latest/), [PyTorch](https://pytorch.org/) 以及 [TensorFlow](https://www.tensorflow.org/) — 并与之无缝整合。你可以直接使用一个框架训练你的模型然后用另一个加载和推理。
