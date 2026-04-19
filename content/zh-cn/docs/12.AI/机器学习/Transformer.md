---
title: Transformer
linkTitle: Transformer
weight: 100
created: 2026-04-08T15:03
---

# 概述

> 参考：
>
> - [GitHub 项目，huggingface/transformers](https://github.com/huggingface/transformers)
> - [Wiki, Transformer_(machine_learning_model)](https://en.wikipedia.org/wiki/Transformer_(machine_learning_model))
> - [Hugging Face 创始人亲述：一个 GitHub 史上增长最快的 AI 项目](https://my.oschina.net/oneflow/blog/5525728)
> - [官方文档](https://huggingface.co/docs/transformers/index)

**Transformer** 架构由 Google 在 2017 年发表的论文 《[Attention is All You Need](https://arxiv.org/abs/1706.03762)》首次提出，它使用 Self-Attention(自注意力) 机制取代了之前在 [NLP](/docs/12.AI/自然语言处理/自然语言处理.md) 任务中常用的 RNN(循环神经网络)，使其成为训练语言模型的代表架构。

[Hugging Face](/docs/12.AI/Hugging%20Face.md) 开源 Transformer 库，提供了便于快速下载和使用的 API，让我们可以把预训练模型用在自己的数据集上微调然后通过 [model hub](https://huggingface.co/models) 与社区共享。同时，每个定义的 Python 模块均完全独立，方便修改和快速研究实验。

Transformers 支持三个最热门的深度学习库： [Jax](https://jax.readthedocs.io/en/latest/), [PyTorch](https://pytorch.org/) 以及 [TensorFlow](https://www.tensorflow.org/)，并与之无缝整合。我们可以直接使用一个框架训练模型然后用另一个加载和推理。

# 安装 Transformers

安装 Transformers 本质就是安装 Transformers 框架，使用该框架可以简单得进行模型的 [Training](/docs/12.AI/机器学习/Training.md) 与 [Inference](/docs/12.AI/机器学习/Inference.md)

Transformers 模型可以对接多种热门的深度学习库：

- [PyTorch](/docs/12.AI/机器学习/PyTorch.md)
  - 注意：安装 PyTorch 时，安装 GPU 版的。如果我们想要使用 GPU 但是却安装的 CPU 版的 PyTorch，将会报错：`Torch not compiled with CUDA enabled`。说白了就是下载的 PyTorch 不是在 CUDA 环境下编译的，无法处理 CUDA 的请求。
- TensorFlow

只安装 Transformers

```bash
uv add transformers
```

安装完 Transformers 包后，可以根据需要安装 PyTorch、TensorFlow 等深度学习的的包。

# 关联文件与配置

**~/.cache/huggingface/** # HuggingFace 缓存路径，保存模型、调用模型的代码 等。可以通过 `${HF_HOME}` 更改路径位置；也可以通过 `${XDG_CACHE_HOME}` 更改路径位置，但是需要注意，`${XDG_CACHE_HOME}` 针对的 `~/.cache/` 这部分。

- **./hub/** # 模型在本地缓存的保存路径。可以通过 `${HUGGINGFACE_HUB_CACHE}` 环境变量变更路径位置。

> 为了防止下载很多模型撑爆 C 盘，个人习惯设置 `${HF_HOME}` 变量为 `D:\appdata\huggingface`

# 模型的关联文件与配置

这些关联文件与 训练 和 推理 任务相关

**权重相关**

- **model.safetensors** # safetensors 格式的权重文件。**那些动辄几个 G 的文件就是这个**。
- **config.json** # 模型结构配置。例如层数、隐藏层大小、注意力头数及 Transformers API 的调用关系等，用于加载、配置和使用预训练模型。
  - **configuration_chatglm.py** # 是 config.json 文件的类表现形式，模型配置的 Python 类代码文件，定义了用于配置模型的 ChatGLMConfig 类。
- **generation_config.json** # 推理配置

**数据预处理相关**

- 对 [自然语言处理](/docs/12.AI/自然语言处理/自然语言处理.md) 来说，通常使用  [Tokenization](/docs/12.AI/自然语言处理/Tokenization.md#Tokenizer%20关联文件与配置)(分词器) 相关文件
- 对 [计算机视觉](/docs/12.AI/计算机视觉/计算机视觉.md) 来说，通常使用

## NLP 关联文件示例

> [!Example] 以我测试过的 Qwen 模型为例（训练与推理都用）
>
> - model.safetensors # 权重
> - config.json # 架构配置
> - generation_config.json # 推理配置
> - tokenizer.json # 分词器的词表
> - tokenizer_config.json # 分词器的配置

## CV 关联文件示例

TODO

# 快速体验

只需要几行代码，就可以在给定任务中下载和使用任何预训练模型，这里官方使用了一个情绪分析模型，用以分析指定文本的情绪是正向的还是负向的：

```python
>>> from transformers import pipeline

# 下载并缓存 pipline 使用的预训练模型
>>> classifier = pipeline('sentiment-analysis')
# 评估给定的文本
>>> classifier('We are very happy to introduce pipeline to the transformers repository.')
[{'label': 'POSITIVE', 'score': 0.9996980428695679}]
```

transformers 库会自动从 Hugging Face 中下载名为 sentiment-analysis 的模型到默认目录中。

## 高级体验

有时我们使用的模型可能会产生某些问题，此时我们可以手动下载模型，比如我们用清华开源的 chatglm-6b 模型举例，只需要先在本地目录下载模型 `git clone https://huggingface.co/THUDM/chatglm-6b-int`，然后运行如下代码即可使用 CPU 体验。其中注意要安装 chatglm-6b 项目中的 Python 依赖。

```python
from transformers import AutoTokenizer, AutoModel
tokenizer = AutoTokenizer.from_pretrained("D:\Projects\DesistDaydream\python-transformers\chatglm-6b-int4", trust_remote_code=True)
model = AutoModel.from_pretrained("D:\Projects\DesistDaydream\python-transformers\chatglm-6b-int4",trust_remote_code=True).float()
model = model.eval()
response, history = model.chat(tokenizer, "你好", history=[])
print(response)
```

代码运行后，获得回复：

```bash
~]# python demo.py
你好👋！我是人工智能助手 ChatGLM-6B，很高兴见到你，欢迎问我任何问题。
```

# Attention Is All You Need

> 参考：
>
> - [公众号，一文彻底讲透GPT架构及推理原理](https://mp.weixin.qq.com/s/moVLtn0_necwuyxdIlosSg)
> - [B 站，硬读 Transformer 经典论文！不过是硬着头皮的硬...](https://www.bilibili.com/video/BV1k4o7YqEEi)
> - [B 站 - 跟李沐学AI，Transformer论文逐段精读【论文精读】](https://www.bilibili.com/video/BV1pu411o7BE)
> - [B 站 - ，王木头学科学，从编解码和词嵌入开始，一步一步理解Transformer，注意力机制(Attention)的本质是卷积神经网络(CNN)](https://www.bilibili.com/video/BV1XH4y1T76e)

TODO:

- 幻觉的来源：预测出第一个字的概率，后面所有出现的都会基于前面的所有得出各种字的权重，如果这个第一个选择错误，后面有可能会越错越多。并且模型本身并不具备向前纠错的能力
- 是否意味着 Transformers 结构本身永远无法解决幻觉问题？
- 想要解决幻觉问题，让 AI 与现实世界接触并验证模型输出结果的权重是否准确，是否是一个有效的做法？比如利用各种 MCP 与现实世界交互。

## 编码/解码
