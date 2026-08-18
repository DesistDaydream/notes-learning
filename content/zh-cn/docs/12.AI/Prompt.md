---
title: Prompt
weight: 50
---

# 概述

> 参考：
>
> -

Prompt 解决方案

- RAG
- Tool calling
    - Function calling
    - [MCP](/docs/12.AI/MCP.md)
- [Skills](#Skills)

# 历史

**Prompt engineering(提示词工程)** # 最早期的称呼，甚至催生了提示词工程师 🤣

- [GitHub 项目，f/awesome-chatgpt-prompts](https://github.com/f/awesome-chatgpt-prompts)
- [GitHub 项目，PlexPt/awesome-chatgpt-prompts-zh](https://github.com/PlexPt/awesome-chatgpt-prompts-zh)
- [公众号-云原生小白，你应该知道的ChatGPT提示语](https://mp.weixin.qq.com/s/BcJWxvhpTRFTE20rB55Sow)

**Content engineering** # 通过各种人为定义的约束来管理 Prompt，e.g [RAG](#RAG), [MCP](/docs/12.AI/MCP.md), [Skills](/docs/12.AI/Skills.md), etc.

**Agent** # [Agent](/docs/12.AI/Agent.md) 程序通过各种机制（调用工具、从 RAG 获取信息、etc.）管理大量 Prompt

**Harness engineering** # 对 Agent 本身的管理，甚至可以加上对 Prompt 的管理。

- [【闪客】你管这破玩意叫 Harness？虚拟世界的牛马套餐！](https://www.bilibili.com/video/BV1cNdrB4Evw)
- Harness engineering 第一次出现于 2026-02-05 [mitchellh Blog, My AI Adoption Journey](https://mitchellh.com/writing/my-ai-adoption-journey#step-5-engineer-the-harness)

个人感觉，Agent 的出现是个分水岭，从早期由人手工管理 Prompt，衍化成由程序管理 Prompt。

# Tool calling

**Tool calling(工具调用)**，它使 LLM(大型语言模型) 能够以结构化的方式和外部系统（可执行程序、API、etc.）进行交互。

历史

- **Function calling(函数调用)** 在 LLM 早期，由 OpenAI 公司推出的工具调用标准
- [**MCP**](/docs/12.AI/MCP.md) 由 Anthropic 推出的工具调用标准

很多模型把这些标准作为数据集进行训练，这样不用其他信息，只需要传输标准化的内容即可让模型返回格式化的信息

最后，Tool calling 被模型作为内部的基本工作，加入到训练中，调用模型时，可以直接使用 tool 角色调用，参考 [Transformer inference](/docs/12.AI/Machine%20learning/Transformer/Transformer%20inference.md) 中相关笔记。

# RAG

> 参考：
>
> - [Wiki, Retrieval-augmented generation](https://en.wikipedia.org/wiki/Retrieval-augmented_generation)
> - [B 站，AI知识库RAG技术原理，三大痛点与进阶方案【不用编程】](https://www.bilibili.com/video/BV1NMoFYoEsb)

**Retrieval-augmented generation(检索增强生成，简称 RAG)** 是一种使 [12.AI](/docs/12.AI/12.AI.md) 模型能够检索信息的技术。它修改了与 LLM 的交互，使模型能够参考指定的一组文档来响应用户的查询，并使用这些信息来补充其预先存在训练数据中的信息。

用户把资料添加进知识库的时候，程序会先把它们拆分成很多个文本块，然后使用嵌入模型对这些文本块进行向量化（向量化指的是把切拆分后的文本），变成一个超长的数字序列。然后程序把向量以及对应的文本保存在向量数据库里面。

接下来用户开始提问，不过这个提问并非直接送达到大模型那里，而是把其本身也经过向量化处理，先变成一个 1024 维的向量。然后把用户的提问与向量数据库进行相似度匹配，这个匹配过程是基于向量的纯数学运算，最后知识库**选出匹配度最高**的几个**原文片段**，再加上用户的问题发给大模型，大模型进行最后的**归纳总结**

存在问题：

- RAG 使用的向量数据库存储的向量信息，依赖模型的 [Tokenization](/docs/12.AI/自然语言处理/Tokenization.md)（i.e. 使用分词器把知识库字符串转成向量存到数据库里）。当换了一个模型时，从向量数据库里检索出来的向量无法呗新模型正确识别。
- 切片很粗暴
- 检索不精准 # 搜索知识库时，只能找到切片，无法将搜索内容与全文进行上下文管理，只有部分切片。最后会 AI 拿到的内容是不足的，导致结果不精准。

**重排序模型**，可以把向量数据库初步检索出来的数据，使用专用的重排序模型进行更深入的语义分析。然后再按照问题的相关性进行重新的排序，把相关性最大的一些数据排到前面并且交付给大模型。这是一种先粗后细的两步检索策略，可以进一步提高检索精度

使用超长上下文，避免切片太碎，但是。。。。资源消耗非常非常高。。。。

# Skills

> 参考：
>
> - [GitHub 项目，anthropics/skills](https://github.com/anthropics/skills) - 2025-10-16 第一次提交
> - [官网](https://agentskills.io/)

## 学习资料

[公众号 - 差评，骗你的，其实AI根本不需要那么多提示词](https://mp.weixin.qq.com/s/OC0l_2M1sKGGhYy8WAgBXQ)

https://github.com/heilcheng/awesome-agent-skills

https://github.com/coreyhaines31/marketingskills
