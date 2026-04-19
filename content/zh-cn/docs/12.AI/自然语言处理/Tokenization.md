---
title: "Tokenizer"
linkTitle: "Tokenizer"
created: "2026-04-13T13:30"
weight: 100
---

# 概述

> 参考：
>
> - https://en.wikipedia.org/wiki/Large_language_model#Tokenization

**Tokenization(标记化)** 通常也称为 **分词**。

由于 [机器学习](/docs/12.AI/机器学习/机器学习.md) 算法处理的是数字而非文本，因此文本必须先转换为数字。具体流程分为以下几步：首先确定词汇表，然后为词汇表中的每个条目任意但唯一地分配一个整数索引，最后将一个 Embedding（向量表示）与该整数索引关联起来。常见的算法包括 [BPE](/docs/8.通用技术/编码与解码/BPE.md) 和 WordPiece。

# Tokenizer

**Tokenizer(分词器)** 是 NLP 进行 训练/推理 时，将数据传入模型之前进行数据预处理的 代码库或程序。

# Tokenizer 关联文件与配置

> [!TODO] 这里好像只是针对 Transformer 架构的

- **tokenizer_config.json** # 分词器的"行为配置文件"，不存数据，存的是"怎么用这个分词器"的元信息。e.g. 如何解析模型的推理结果、如何添加特殊 Token、etc. 。
- **tokenizer.json** # 词表文件。包含了 词表、合并规则、特殊 Token、etc.
    - BPE（HF 原生，e.g. Qwen 系列模型）将整个文件拆分成两部分，merges.txt 决定"怎么切"，vocab.json 决定"切完之后每个 Token 叫什么 ID"。
        - **merges.txt** # 词表文件。包含了 Token 合并的规则，用于将子 Token 合并成一个 Token。
        - **vocab.json** # 词表文件。包含了 TokenID 到 Token字符串 的映射关系，以及其他与 Token 相关的信息。
