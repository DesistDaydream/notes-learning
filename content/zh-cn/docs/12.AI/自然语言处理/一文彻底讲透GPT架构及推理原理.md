---
title: "一文彻底讲透GPT架构及推理原理"
source: "https://mp.weixin.qq.com/s/moVLtn0_necwuyxdIlosSg"
author:
  - "[[邹佳旭]]"
published:
created: 2025-02-19
description: "从开发人员的视角，围绕着大模型的正向推理过程，对大模型的原理的系统性总结，希望对初学者有所帮助。"
tags:
  - "clippings"
---

*2025年01月19日 22:30*

导读

本篇是作者从开发人员的视角，围绕着大模型正向推理过程，对大模型的原理的系统性总结，希望对初学者有所帮助。

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU523zPFj4bsgL9DiaiaG668V2yzxVobL1VkSicW6KbuecZmL9WhL271dLZQ/640?wx_fmt=jpeg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

引言

什么是人工智能？

清华大学出版社出版的《人工智能概论》中提出，人工智能是对人的意识、思维的信息过程的模拟。人工智能不是人的智能，但它能像人那样思考，也可能超过人的智能。﻿

基于这个设想，人工智能应当能够执行通常需要人类智能的任务，如视觉感知、语音识别、决策和语言翻译等工作。就像人一样，可以看见、听见、理解和表达。这涉及了众多人工智能的分支学科，如计算机视觉（CV）、自然语言处理（NLP）、语音识别（VC）、知识图谱（KG）等。

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU518ulaDS5Zphh9bgsIdfGakrbRdexuZHQdCwQasuX4hVEMjibzLvQtSg/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

NLP语言模型的发展，引自《A Survey of Large Language Models》

NLP作为其中之一，其发展历经了多个阶段。在统计语言模型阶段，模型利用马尔科夫假设学会了如何进行预测。在神经语言模型的阶段，我们通过词嵌入的向量表示让模型开始理解词的语义，并通过神经网络模型给人工智能装上了神经元，让模型开始“思考”。在预训练语言模型阶段，我们通过预训练告诉语言模型，要先学习走路再去跑。而在如今的大语言模型阶段，我们基于扩展法则认识到了了力大砖飞的重要性，并收获了各种涌现能力的惊喜，为AGI的发展立下了一个新的里程碑。

如今大模型已经带领NLP 领域走向了全新的阶段，并逐步向其他AI领域扩展。LangChain的2024年Agent调研报告显示，51.1%的企业已经在生产中使用Agent，78.1%的企业有计划将Agent投入生产，AI的浪潮已经席卷而来。并且AI发展的步伐还在加快。模型方面，OpenAI在2024年12月的发布会上，推出了o1专业版，凭借其在后训练阶段的深度推理，使得模型在数学、编程、博士级问题上达到了人类专家的水平，相比于此前的大模型达到了断层领先，近期国内Qwen推出了QwQ模型，月之暗面的KIMI推出了KIMI探索版，也都具有深度推理能力，国内外AI上的差距正在逐步缩短。应用方面，各家都在大力发展多模态AI，Sora、可灵等，让视频创意的落地成本从数十万元降低到了几十元，VLM的加持下的端到端自动驾驶也在大放异彩。

总而言之，我们仍然处在第四次科技革命的起点，业内预测在2025年我们可能会见证AGI（通用人工智能）的落地，并且还提出了ASI（超级人工智能）的概念，认为人类可能创造出在各个方面都超越人类智能的AI，AI很有可能带领人类进入新的时代。﻿

目光放回到大模型本身，对于开发人员而言，大模型的重要性，不亚于JAVA编程语言。大模型的推理原理，就像JVM虚拟机原理一样，如果不了解，那么在使用大模型时难免按照工程化的思维去思考，这样常常会遇到困难，用不明白大模型。比如为什么模型不按照提示词的要求输出，为什么大家都在用思维链。

而本篇是作者从开发人员的视角，围绕着大模型的正向推理过程，对大模型的原理的系统性总结，希望对像我一样的初学者有所帮助。文中会引入部分论文原文、计算公式和计算过程，同时会添加大量例子和解释说明，以免内容晦涩难懂。篇幅较长难免有所纰漏，欢迎各位指正。

大语言模型架构

**Transformer**

### What is Attention

### ![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5xPUEODibToRnLHxQKY2Ia9WP7TuYBIkVzaS8NrgDy20MibPHiaktiaACDg/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

### 镇楼图，来自万物的起源《Attention is All You Need》

Transformer架构由Google在2017年发表的论文《Attention is All You Need》首次提出，它使用自注意力（Self-Attention）机制取代了之前在 NLP 任务中常用的RNN（循环神经网络），使其成为预训练语言模型阶段的代表架构。

要澄清一点，注意力机制并不是由Transformer提出。注意力机制最早起源于计算机视觉领域，其思想在上世纪九十年代就已经被提出。在2014年的论文《Neural Machine Translation by Jointly Learning to Align and Translate》中注意力机制首次应用于NLP领域，结合了RNN和注意力机制，提出了序列到序列（Seq2Seq）模型，来改进机器翻译的性能。Transformer的贡献在于，它提出注意力机制不一定要和RNN绑定，我们可以将注意力机制单独拿出来，形成一套全新的架构，这也就是论文标题《Attention Is All You Need》的来源。﻿

在正式开始前，首先需要说明，什么是注意力机制？

在NLP领域，有三个基础概念：

- 分词（Tokenization）：首先大模型会将输入内容进行分词，分割成一系列的词元（Token），形成一个词元序列。
- 词元（Token）：指将输入的文本分割成的最小单位，词元可以是一个单词、一个词组、一个标点符号、一个字符等。
- 词嵌入（Embedding）：分词后的词元将被转换为高维空间中的向量表示，向量中包含了词元的语义信息。

举个例子，当我们输入“我可以拥有一杯咖啡吗？”时，首先通过分词形成“我”、“可以”、“拥有”、“一杯”、“咖啡”、“吗？”这几个词元，然后通过词嵌入转变成高维空间中的向量。

在向量空间中，每一个点代表一个实体或者概念，我们称之为“数据点”。这些数据点可以代表着具体的单词、短语、句子，他们能够被具体地、明确地识别出来，我们称之为“实体”；也可以不直接对应于一个具体的实体，而是表达一种对事物的抽象理解，我们称之为“概念”。这些数据点在向量空间中的位置和分布反映了实体或概念之间的相似性和关系，相似的实体或概念在空间中会更接近，而不同的则相距较远。其所在空间的每一个方向（维度）都代表了数据点的某种属性或特征。这些属性和特征是通过模型训练获得的，可能包括了“情感”、“语法”、“词性”等方面，事实上由于这些属性和特征是模型通过训练进行内部学习的结果，他们的具体含义往往难以直观解释清楚。

当词元被嵌入到向量空间中，每个词元都会形成一个实体的数据点，此时这些数据点所对应的向量，就代表了词本身的语义。但是在不同的上下文语境中，相同的词也会有不同的语义。比如在“科技巨头苹果”、“美味的苹果”这两句话中，“苹果”这个词，分别代表着苹果公司和水果。因此模型需要对原始向量“苹果”，基于上下文中“科技巨头”、“美味的”等信息，来进行数据点位置的调整，丰富其语义。换句话说，词元与上下文中各个词元（包括其自己）之间具有一定程度的“依赖关系”，这种关系会影响其自身的语义。

为了进一步解释上述内容，我画了一张图帮助大家理解。如下图是一个向量空间，我们假设其只有两个维度，分别是“经济价值”和“食用价值”。“科技巨头”、“美味的”、“苹果”在词嵌入后，会在空间中分别形成自己的数据点，显然“科技巨头”相对于“美味的”的经济价值属性更明显而食用价值的属性更模糊。当我们输入“科技巨头苹果”时，“苹果”的含义会受“科技巨头”的影响，其数据点的位置会向“科技巨头”靠拢，在空间中会形成一个新的概念。

﻿![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5GwARa3tDeYf7Q16XuL5VVxwWAUux2rDnMNqTLA1CNDqRpqtkn0aynw/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

向量空间中的“科技巨头苹果”

每个词元都需要捕捉自身与序列中其他词元的依赖关系，来丰富其语义。在“我配拥有一杯咖啡吗？”中，对于“配”而言，它依赖“我”作为其主语，这是一条依赖关系。对于“咖啡”而言，“一杯”是其量词，这也是一条依赖关系。

这种关系不只会在相邻的词元之间产生，在论文中有个词叫长距离依赖关系（long-range dependencies ），它指在词元序列中，相隔较远的位置的词元之间的相互依赖或关联。比如，在“我深切地感觉到，在这段漫长而忙碌的日子里，保持清醒和集中精力非常有用，难道此时不配拥有一杯咖啡吗？”中，“我”和“配”之间相隔很远，但他们仍然具有语法层面的依赖关系。这种依赖关系可能非常长，从文章的开头，一直到文章的结尾。比如，在一篇议论文中，作者可能在文章开头提出一个论点，然后通过一系列的论据和分析来支持这个论点，直到文章结尾可能再次强调或总结这个论点。在这个过程中，文章开头的论点和结尾的总结之间就存在着长距离依赖关系，它们在语义上是紧密相连的。

不过虽然词元与各个词元之间都可能存在依赖关系，但是其依赖程度不同，为了体现词元之间的依赖程度，NLP选择引入“注意力机制”。“注意力机制”可以动态地捕捉序列中不同位置元素之间的依赖关系，分析其强弱程度，并根据这些依赖关系生成新的序列表示。其核心思想是模仿人类的注意力，即在处理大量信息时，能够聚焦于输入数据的特定部分，忽略掉那些不太重要的信息，从而更好地理解输入内容。﻿

继续以“我配拥有一杯咖啡吗？”为例，读到“拥有”这个词元时，我们会发现“我”是“拥有”的主语，“配”是对“拥有”的强调，他们都与“拥有”产生了依赖关系。这句话的核心思想，是某人认为自己有资格拥有某物，所以可能“配”相对“我”而言，对“拥有”来说更重要，那么我们在分析“拥有”这个词的语义时，会给“配”分配更多的注意力，这通常体现为分配更高的“注意力权重”。﻿

实际上，我们在日常工作中已经应用了注意力机制。比如，我们都知道思维链（COT，常用<输入, 思考, 输出>来描述）对大模型处理复杂问题时很有帮助，其本质就是将复杂的问题拆分成相对简单的小问题并分步骤处理，使模型能够聚焦于问题的特定部分，来提高输出质量和准确性。“聚焦”的这一过程，就是依赖模型的注意力机制完成。通常模型会依赖输出内容或内部推理（如o1具有内部推理过程，即慢思考）来构建思考过程，但哪怕没有这些内容，仅仅依靠注意力本身，COT也能让模型提高部分性能。

### Why Transformer

自注意力机制是Transformer的核心。在论文中，Transformer说明了三点原因，来说明为何舍弃RNN和CNN，只保留注意力机制

> Transformer论文：《Attention is All You Need》﻿
> 
> 原文：In this section we compare various aspects of self-attention layers to the recurrent and convolutional layers commonly used for mapping one variable-length sequence of symbol representations (x1, ..., xn) to another sequence of equal length (z1, ..., zn), with xi, zi ∈ Rd, such as a hidden layer in a typical sequence transduction encoder or decoder. Motivating our use of self-attention we consider three desiderata. - One is the total computational complexity per layer.
> 
> \- Another is the amount of computation that can be parallelized, as measured by the minimum number of sequential operations required.
> 
> \- The third is the path length between long-range dependencies in the network. Learning long-range dependencies is a key challenge in many sequence transduction tasks. One key factor affecting the ability to learn such dependencies is the length of the paths forward and backward signals have to traverse in the network. The shorter these paths between any combination of positions in the input and output sequences, the easier it is to learn long-range dependencies \[12\].
> 
> 译文：在这一部分中，我们比较了自注意力层与通常用于将一个可变长序列的符号表示（x1, ..., xn）映射到另一个等长序列（z1, ..., zn）的循环层和卷积层的不同方面，其中xi, zi ∈ Rd。我们考虑了三个主要因素来激励我们使用自注意力：
> 
> \- 每层的总计算复杂性
> 
> \- 可以并行化的计算量，通过所需的最小顺序操作数来衡量。
> 
> \- 网络中长距离依赖之间的路径长度。学习长距离依赖性是许多序列转换任务中的一个关键挑战。影响学习这些依赖性能力的一个关键因素是前向和后向信号在网络中必须穿越的路径长度。输入和输出序列中任意两个位置之间的路径越短，学习长距离依赖性就越容易。

其中后两点尤其值得注意：

- 并行化计算：在处理序列数据时，能够同时处理序列中的所有元素，这很重要，是模型训练时使用GPU的并行计算能力的基础。
- 长距离依赖捕捉：在注意力机制上不依赖于序列中元素的特定顺序，可以有效处理序列中元素之间相隔较远但仍相互影响的关系。

可能有些难以理解，让我们输入“我配拥有一杯咖啡？”来进行文本预测，分别看一下RNN和Transformer的处理方式。

![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5CqI1lq3x1ViavKBtw9OyibtkVotn07fbn73VLJpVr8SiadnVUJicAK3BNA/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)﻿

依赖捕捉上的差异，是Transformer可以进行并行处理的前提，而并行处理是Transformer的核心优势之一。在预训练语言模型阶段，预训练（Pretrain）+微调（Finetune）是模型训练的主要范式。Transformer的并行化计算能力大大提高了模型训练的速度，长距离依赖捕捉能力为模型打开了上下文窗口，再结合位置编码等能力，使得Transformer相对于RNN获得了显著优势。而预训练又为Transformer带来了特征学习能力（从数据中自动学习和提取特征的能力）和迁移学习能力（将在一个任务上学习到的知识应用到另一个相关任务上的能力）的提升，显著提升了模型的性能和泛化能力。

这些优势，也正是GPT选择Transformer的原因。

> GPT-1的论文《Improving Language Understanding by Generative Pre-Training》﻿
> 
> 原文：For our model architecture, we use the Transformer, which has been shown to perform strongly on various tasks such as machine translation, document generation, and syntactic parsing. This model choice provides us with a more structured memory for handling long-term dependencies in text, compared to alternatives like recurrent networks, resulting in robust transfer performance across diverse tasks.
> 
> 翻译：对于我们的模型架构，我们使用了Transformer，它在机器翻译、文档生成和句法解析等各种任务上都表现出色。与循环网络等替代方案相比，这种模型选择为我们提供了更结构化的记忆来处理文本中的长期依赖关系，从而在多样化任务中实现了稳健的迁移性能。

### Transformer

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5pyAMPNia8BEsROhWCQBSKzgvAzz8icvPqXicNrDnkRtIhxTX1hJsjqm9g/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

Transformer架构图，引自《A Survey of Large Language Models》﻿

理解了Transformer的优势后，让我们先忘记RNN，回到Transformer本身。从架构图中可知，Transformer架构分为两部分：

- 编码器：用于理解输入序列的内容。它读取输入数据，并将其转换成一个连续的表示，这个表示捕捉了输入数据的关键信息和上下文。
- 解码器：用于生成输出序列，使用编码器的输出和之前生成的输出来预测序列中的下一个元素。

以文本翻译为例，输入“我爱喝咖啡”，并要求模型将其翻译成英文。编码器首先工作，通过理解每个词元本身的含义，以及其上下文的依赖关系，形成一种向量形式的中间表示，并传递给解码器，这里面包含了整个序列的语义，即“我爱喝咖啡”这句话的完整含义。解码器结合该信息，从零开始，不断地预测下一个最可能的英文词元并生成词元，直至完成翻译。

值得注意的是，在解码器的多头注意力层（也叫交叉注意力层，或编码器-解码器自注意力层），编码器所提供的输出中，“我爱喝咖啡”这句话的含义已经明示。这对于文本翻译这种序列到序列的任务而言，可以确保生成内容的准确性，但对于预测类的任务而言，无疑是提前公布了答案，会降低预测的价值。

### Transformer to GPT

随着技术的演进，基于Transformer已经形成了三种常见架构

- 编码器-解码器架构（Encoder-Decoder Architecture），参考模型：T5
- 编码器架构（Encoder-Only Architecture），参考模型：BERT
- 解码器架构（Decoder-Only Architecture），参考模型：GPT（来自OpenAI）、Qwen（来自通义千问）、GLM（来自清华大学）

﻿其中编码器-解码器架构，适合进行序列到序列的任务，比如文本翻译、内容摘要。编码器架构，适合需要对输入内容分析但不需要生成新序列的任务，比如情感分析、文本分类。解码器架构，适合基于已有信息生成新序列的任务，比如文本生成、对话系统。

﻿解码器架构下，又有两个分支：

- 因果解码器（Causal Decoder），参考模型：GPT、Qwen
- 前缀解码器（Prefix Decoder），参考模型：GLM

二者之间的主要差异在于注意力的模式。因果解码器的特点，是在生成每个词元时，只能看到它之前的词元，而不能看到它之后的词元，这种机制通过掩码实现，确保了模型在生成当前词元时，不会利用到未来的信息，我们称之为“单向注意力”。前缀解码器对于输入（前缀）部分使用“双向注意力”进行编码，这意味着前缀中的每个词元都可以访问前缀中的所有其他词元，但这仅限于前缀部分，生成输出时，前缀解码器仍然采用单向的掩码注意力，即每个输出词元只能依赖于它之前的所有输出词元。﻿

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5531iaT535I9uXQiajT4EogXJibpff0h9OibeMGHotfQXvZLTE4eLJC7CMQ/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

不同架构的注意力模式比较，引自《A Survey of Large Language Models》﻿

可能有些晦涩，让我们参考《大语言模型概述》的例子来说明一下。该例子中，已经存在了“大家共同努力”这六个词元，模型正在思考如何产生下一个新的词元。此时，“大家共”是输入（前缀），“同努力”是模型解码已经产生的输出，蓝色代表可以前缀词元之间可以互相建立依赖关系，灰色代表掩码，无法建立依赖关系。﻿

因果解码器和前缀解码器的差异在“大家共”（前缀）所对应的3\*3的方格中，两种解码器都会去分析前缀词元之间的依赖关系。对于因果解码器而言，哪怕词元是前缀的一部分，也无法看到其之后的词元，所以对于前缀中的“家”（对应第二行第二列），它只能建立与“大”（第二行第一列）的依赖关系，无法看到“共”（第二行第三列）。而在前缀解码器中，“家”可以同时建立与“大”和“共”的依赖关系。﻿

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5qjmAibAah1cNZpfmrfjGjNJfIJFn1ibRViaZM8QbbmQQCH6YKp9XpaUew/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

最后贴一张图大家感受一下不同架构之间注意力模式上的差异。

单向注意力和双向注意力，有各自的优势，例如，对于文本生成任务，可能会优先考虑单向注意力以保持生成的连贯性。而对于需要全面理解文本的任务，可能会选择双向注意力以获取更丰富的上下文信息。如上图分别是编码器架构（BERT）、编码器-解码器架构（T5）、因果解码器架构（GPT）、前缀解码器架构（T5、GLM）的注意力模式。

这些架构在当今的大语言模型阶段都有应用，其中因果解码器架构是目前的主流架构，包括Qwen-2.5在内的众多模型采用的都是这种架构。具体可以参考下面的大语言模型架构配置图，其中类别代表架构，L 表示隐藏层层数，N 表示注意力头数，H 表示隐藏状态的大小。﻿

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5Qdx2fOAmY4mZSYm0Np2xg2HqQpw64ia7b5tOoKOqGZjL6nvQE8GYTOg/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

大语言模型架构配置表，引自《A Survey of Large Language Models》﻿

从2018年GPT-1开始，模型的基本原理确实经历了一些变化和改进，但是探讨其基本架构仍然有价值。

**GPT**

时间来到2018年，OpenAI团队的论文《Improving Language Understanding by Generative Pre-Training》横空出世，它提出可以在大规模未标注数据集上预训练一个通用的语言模型，再在特定NLP子任务上进行微调，从而将大模型的语言表征能力迁移至特定子任务中。其创新之处在于，提出了一种新的预训练-微调框架，并且特别强调了生成式预训练在语言模型中的应用。生成式，指的是通过模拟训练数据的统计特性来创造原始数据集中不存在的新样本，这使得GPT在文本生成方面具有显著的优势。﻿

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5N6s4hbdkGwqM8s24iblibQtQcQAnXH9PRBhwRZgl8vTOXqXEZ3Qo69iaQ/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

GPT模型架构，引自GPT-1的论文《Improving Language Understanding by Generative Pre-Training》﻿﻿

上图来自于GPT-1的论文，图片左侧是GPT涉及到的核心组件，这是本文的重点内容。图片右侧是针对不同NLP任务的输入格式，这些格式是为了将各种任务的输入数据转换为Transformer模型能够处理的序列形式，用于模型与训练的过程。﻿

GPT使用了Transformer的解码器部分，同时舍弃了编码器中的交叉注意力机制层，保留了其余部分。整体上模型结构分为三部分：

- 输入层（Input Layer）：将文本转换为模型可以处理的格式，涉及分词、词嵌入、位置编码等。
- 隐藏层（Hidden Layer）：由多个Transformer的解码器堆叠而成，是GPT的核心，负责模型的理解、思考的过程。
- 输出层（Output Layer）：基于隐藏层的最终输出生成为模型的最终预测，在GPT中，该过程通常是生成下一个词元的概率分布。

在隐藏层中，最核心的两个结构分别是

- 掩码多头自注意力层（Masked Multi Self Attention Layers，对应Transformer的Masked Multi-Head Attention Layers，简称MHA，也叫MSA）。
- 前置反馈网络层（Feed Forward Networks Layers，简称FFN，与MLP类似）。

MHA的功能是理解输入内容，它使模型能够在处理序列时捕捉到输入数据之间的依赖关系和上下文信息，类似于我们的大脑在接收到新的信息后进行理解的过程。FFN层会对MHA的输出进行进一步的非线性变换，以提取更高级别的特征，类似于我们的大脑在思考如何回应，进而基于通过训练获得的信息和知识，产生新的内容。﻿

举个例子，当我们输入“美国2024年总统大选胜出的是”时，MHA会理解每个词元的含义及其在序列中的位置，读懂问题的含义，并给出一种中间表示，FFN层则会对这些表示进行进一步的变换，进而从更高级别的特征中得到最相近的信息——“川普”。

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU51pmwyyE1icTR7ARiaaGT0CYFTsmU6XSQuicuZXy3unWiaPxAlR81CcDicRQ/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

MHA和FFN的多层结构，引自3Blue1Brown的视频《GPT是什么？直观解释Transformer》﻿﻿

隐藏层不只有一层，而是一种多层嵌套的结构。如上图，Attention是MHA，Multilayer Perceptron（MLP）是FFN，它们就像奥利奥饼干一样彼此交错。这是为了通过建立更深的网络结构，帮助模型在不同的抽象层次上捕捉序列内部的依赖关系，最终将整段文字的所有关键含义，以某种方式充分融合到最后的输出中。隐藏层的层数并不是越多越好，这取决于模型的设计，可以参考前文贴过的模型参数表，其中的L就代表该模型的隐藏层的层数。较新的Qwen2-72B有80层，GPT-4有120层。

目前主流的大模型，在这两层都进行了不同程度的优化。比如Qwen2使用分组查询注意力（Grouped Multi-Query Attention，简称GQA）替代MHA来提高吞吐量，并在部分模型上尝试使用混合专家模型（Mixture-of-Experts，简称MoE）来替代传统FFN。GPT-4则使用了多查询注意力（Multi-Query Attention）来减少注意力层的KV缓存的内存容量，并使用了具有16个专家的混合专家模型。关于这里提到的较新的技术，在后文中会详细阐述。﻿

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5dgz7VkN7gpR0bXwEx7Ife6iaeJYOfamEMDFIAOKL9jP1bIR597M5WIA/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

编码器-解码器架构与解码器架构，引自Llama的论文﻿

此外，模型还保留了Transformer的其他部分，包括了（参考上图右半部分，该图片的细节更多一些）

- 词嵌入 （Embedding，对应GPT论文中的Text & Position Embed）。
- 位置编码（Positional Encodings，简称PE，对应GPT论文中的Text & Position Embed，Rotary Positional Encodings是位置编码的一种技术）。
- 层归一化（Layer Norm，上图中表示为RMS Norm，通常与残差连接一起用，Layer Norm和RMS Norm是归一化的两种不同技术）。
- 线性层（Linear，负责将FFN层的输出通过线性变换，通常用于将模型的输出映射到所需的维度）。
- Softmax（Softmax层，负责生成概率分布，以便进行最终的预测）。

这些部分单独拿出来看会有些抽象，让我们尝试将一段文本输入给大模型，看一看大模型的整体处理流程

1.分词（Tokenization）：首先大模型会进行分词，将文本内容分割成一系列的词元（token）。

2.词嵌入（Embedding）：分词后的词元将被转换为高维空间中的向量表示，向量中包含了词元的语义信息。

3.位置编码（PE）：将词元的在序列中的位置信息，添加到词嵌入中，以告知模型每个单词在序列中的位置。

4.掩码多头自注意力层（MHA）：通过自注意力机制捕捉序列内部词元间的依赖关系，形成对输入内容的理解。

5.前馈反馈网络（FFN）：基于MHA的输出，在更高维度的空间中，从预训练过程中学习到的特征中提取新的特征。

6.线性层（Linear）：将FFN层的输出映射到词汇表的大小，来将特征与具体的词元关联起来，线性层的输出被称作logits。

7.Softmax：基于logits形成候选词元的概率分布，并基于解码策略选择具体的输出词元。﻿

在该流程中，我故意省略了层归一化，层归一化主要是在模型训练过程中改善训练过程，通过规范化每一层的输出，帮助模型更稳定地学习。而在使用模型时，输入数据会按照训练阶段学到的层归一化参数进行处理，但这个参数是固定的，不会动态调整，从这个意义上说，层归一化不会直接影响模型的使用过程。﻿

分词、词嵌入、位置编码，属于输入层，MHA、FFN属于隐藏层，线性层、Softmax属于输出层。可能读到这里，对于输入层的理解会相对清晰，对其他层的理解还有些模糊。没关系，它们的作用有简单的几句话很难描述清楚，请继续往下读，在读完本文的所有内容后，再回头来看会比较清楚。本文中，我们会围绕着隐藏层和输出层来详细介绍，尤其是其中GPT的核心——隐藏层。﻿

自注意力机制（Attention）﻿

**MHA（多头注意力）**

MHA，全拼Multi-Head Attention（多头注意力），在GPT等因果解码器架构下模型中，指掩码多头自注意力，全拼Masked Multi Self Attention。“掩码”赋予了GPT单向注意力的特性，这符合因果解码器的架构，“多头”使GPT可以从不同角度发现特征，自注意力中的“自”指的是模型关注的是单一词元序列，而不是不同序列之间的特征，比如RNN的循环注意力围绕的是同一序列的不同时间步。﻿

多头注意力机制的概念，来自于上古时期Google的论文《Attention Is All You Need》 。

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5sw03vicB8JuzSuDAOV8lia7F17MnbfAAStmV3WyBGzhzVD6U1kh6qoKQ/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

缩放点积注意力（左）和多头注意力（右），来自《Attention Is All You Need》

从上图可以看到，多头注意力是由多个并行运行的缩放点积注意力组成，其中缩放点积注意力的目的，是帮助模型在大量数据中快速找到并专注于那些对当前任务最有价值的信息，让我们由此讲起。﻿

### 单头注意力（缩放点积注意力）

![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5XM47ibyohV0U2Nl8qk2ZJP5hRt4OaebrhJjdn2hicRg5dZSZ13lGXPZQ/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)﻿

了解模型计算注意力的过程是很重要的，Transformer团队使用了上图中非常简洁的公式来表达缩放点积注意力的计算过程。首先要再次明确一下，注意力的计算是词元维度的，它计算的是当前词元与上下文中其他词元的依赖关系，并在此基础上调整词元本身的语义。比如在“我配拥有一杯咖啡吗？”中，会分别计算“我”、“配”、“拥有”、“一杯”、“咖啡”、“吗？”各自的注意力，并分别调整每个词元的语义。

整个公式可以看作两部分，首先是含softmax在内的注意力权重计算过程，其作用是计算“当前词元”与“其他词元”（包含当前词元自身）之间的注意力权重，来体现他们之间的依赖程度，其结果是一个总和为1的比例分布。比如输入“我爱你”时，“你”会分别计算与“我”、“爱”、“你”三个词元的注意力权重，并获得一个比例分布比如\[0.2,0.3,0.5\]。﻿

然后，这些注意力权重会分别与“其他词元”各自的![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5eaLR8iaoxz7aKvKSO4xOM2EhQbkEhBbtSwibhjx0viaIfmnJnXibbibLpEQ/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)相乘获得“当前词元”的数据点在向量空间中偏移的方向和距离。比如，我们设原本“你”的数据点的坐标是![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU50ZJvbicXSNaLlERKDhtjj9b1tTdtyTvwbj4NtzXdeZPUw6GdqHrbSibQ/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)，那么在注意力计算后![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU50ZJvbicXSNaLlERKDhtjj9b1tTdtyTvwbj4NtzXdeZPUw6GdqHrbSibQ/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)的值会变成![图片](https://mp.weixin.qq.com/s/www.w3.org/2000/svg'%20xmlns:xlink='http://www.w3.org/1999/xlink'%3E%3Ctitle%3E%3C/title%3E%3Cg%20stroke='none'%20stroke-width='1'%20fill='none'%20fill-rule='evenodd'%20fill-opacity='0'%3E%3Cg%20transform='translate(-249.000000,%20-126.000000)'%20fill='%23FFFFFF'%3E%3Crect%20x='249'%20y='126'%20width='1'%20height='1'%3E%3C/rect%3E%3C/g%3E%3C/g%3E%3C/svg%3E)，其计算的方式就是![图片](https://mp.weixin.qq.com/s/www.w3.org/2000/svg'%20xmlns:xlink='http://www.w3.org/1999/xlink'%3E%3Ctitle%3E%3C/title%3E%3Cg%20stroke='none'%20stroke-width='1'%20fill='none'%20fill-rule='evenodd'%20fill-opacity='0'%3E%3Cg%20transform='translate(-249.000000,%20-126.000000)'%20fill='%23FFFFFF'%3E%3Crect%20x='249'%20y='126'%20width='1'%20height='1'%3E%3C/rect%3E%3C/g%3E%3C/g%3E%3C/svg%3E)。如果说![图片](https://mp.weixin.qq.com/s/www.w3.org/2000/svg'%20xmlns:xlink='http://www.w3.org/1999/xlink'%3E%3Ctitle%3E%3C/title%3E%3Cg%20stroke='none'%20stroke-width='1'%20fill='none'%20fill-rule='evenodd'%20fill-opacity='0'%3E%3Cg%20transform='translate(-249.000000,%20-126.000000)'%20fill='%23FFFFFF'%3E%3Crect%20x='249'%20y='126'%20width='1'%20height='1'%3E%3C/rect%3E%3C/g%3E%3C/g%3E%3C/svg%3E)就是在前文What is Attention小节中举例的“科技巨头苹果”中“苹果（实体）”所在的位置，那么此时![图片](https://mp.weixin.qq.com/s/www.w3.org/2000/svg'%20xmlns:xlink='http://www.w3.org/1999/xlink'%3E%3Ctitle%3E%3C/title%3E%3Cg%20stroke='none'%20stroke-width='1'%20fill='none'%20fill-rule='evenodd'%20fill-opacity='0'%3E%3Cg%20transform='translate(-249.000000,%20-126.000000)'%20fill='%23FFFFFF'%3E%3Crect%20x='249'%20y='126'%20width='1'%20height='1'%3E%3C/rect%3E%3C/g%3E%3C/g%3E%3C/svg%3E)就是“苹果公司（概念）”所在的位置。

上述描述可能还不够透彻，请跟着我进一步逐步拆解其计算过程。

首先，X为输入的词元序列的嵌入矩阵，包含了词元的语义信息和位置信息，矩阵中的每一列就是一个词元的向量，列的长度就是隐藏层的参数量，比如GPT-3的隐藏层参数量是12288，那么在输入100个词元的情况下，矩阵的大小就是100 \* 12288。

﻿![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5LqQVOADT5RnvQPjsnpwA3ia8MEG4iaeJ3wDPRXFxd5Tticjkz4EamP3jQ/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)是通过训练得到的三个权重矩阵，在模型训练过程中这三个参数矩阵可以采用随机策略生成，然后通过训练不断调整其参数。输入矩阵X通过与![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5LqQVOADT5RnvQPjsnpwA3ia8MEG4iaeJ3wDPRXFxd5Tticjkz4EamP3jQ/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)相乘，可以分别得到Q、K、V三个矩阵。﻿

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5STn8KLC6RopalH95T7iawyY7fLGyfMKOeMTp72UHyibcK6PIT2YAeuMA/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

QKV计算公式，引自《A Survey of Large Language Models》﻿

那么Q、K、V分别是什么？

> 原文：An attention function can be described as mapping a query and a set of key-value pairs to an output, where the query, keys, values, and output are all vectors. The output is computed as a weighted sum of the values, where the weight assigned to each value is computed by a compatibility function of the query with the corresponding key.
> 
> 翻译：注意力函数可以被描述为将一个查询（query）和一组键值对（key-value pairs）映射到一个输出（output）的过程，其中查询、键、值和输出都是向量。输出是通过值的加权和计算得到的，每个值所分配的权重是通过查询与相应键的兼容性函数（compatibility function）计算得出的。﻿

我们可以认为

- Q：查询，通常可以认为模型代表当前的“问题”或“需求”，其目的是探寻其他词元与当前词元的相关性。
- K：键，即关键信息，这些信息用于判断词元的相关性，可能包括语义信息、语法角色或其他与任务相关的信息。
- V：值，即对于键所标识的关键信息的具体“回应”或“扩展”，可以认为它是键背后的详细信息。

Q属于当前词元，而K、V都属于其他词元，每个词元的K、V之间都是相互绑定的。就像是一本书，V就是这本书的摘要、作者、具体内容，K就是这本书的标签和分类。

每个词元，都会用自身的Q和其它词元的K，来衡量![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5Yu9np9lIJR2KJGrKibgibxJicHS0vHwyu6G2fvnP472WicSqVPPze6hmPw/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)与![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5ITQluicz1unuC9UEMeYSPe63uHx1M70KvDBwTByVmG65w6icEND0EeHA/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)之间的相关性。衡量的方式就是计算QK之间的点积，在公式中体现为![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5IuMVstQJYmPauOydEP8a1YUkibmcyrsQ08AiavD3G2azjptIWQS1jdMQ/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)。点积运算是综合衡量两个向量在各个维度上的相关性的一种方式，比如![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5Qgydve2vibAds839W35V0ncWvZtwUlkfq6eUTeVRHcttjfEhB8ZMhYw/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)，其点积的结果就是![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5BJrYmld8sHrkhXibJ4gp5jtcbs0ICrsgU6FleiapnqUS0uWHYGIlKZIg/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)，即对应分量的乘积之和。相关性的大小，就是通过点积运算获得的结果的大小体现的。﻿

举个例子，以“有一天，一位勇敢的探险家”作为输入序列给大模型

- Q：假设Q代表着“查找最相关的信息点”这个概念，那么在“探险家”这个词元上，Q所对应的这一列，可能代表“和探险家最相关的信息点是什么？”，那它寻找的可能就是“冒险”、“勇气”、“积极”等特征。
- K：在“勇敢的”这个词元上，K可能在语义上与“冒险”、“勇气”相关联，在语法上与“修饰语”相关联，在情感上与“积极”相关联。这些都属于关键信息。
- V：在“探险家”这个token上，如果Q是“和探险家最相关的信息点是什么”，那么V向量将提供与“探险家”直接相关的上下文信息，比如探险家的行为、特征或相关事件。

对于“探险家”所对应的Q的向量，会分别与序列中的其他词元对应的K向量计算点积，获得与“有一天”、“一位”、“勇敢的”、“探险家”这些词元点积的值，这代表着二者之间的语义相似度。有可能，“探险家”由于其也具有“冒险”、“勇气”等关键特征，其与“勇敢的”的点积相对更大，那么其在归一化后的占比也会更高，在GPT的理解中，“勇敢的”对于“探险家”而言，就更重要。

关于为何要使用点积进行运算，论文中也进行了分析，文中分析了加性注意力和点积注意力的这两种兼容性函数之间的差异，发现点积注意力能够捕捉序列中元素间的长距离依赖关系的同时，在计算上高效且能稳定梯度。

> 原文：The two most commonly used attention functions are additive attention \[ 2\], and dot-product (multi-plicative) attention. Dot-product attention is identical to our algorithm, except for the scaling factor. Additive attention computes the compatibility function using a feed-forward network with a single hidden layer. While the two are similar in theoretical complexity, dot-product attention is much faster and more space-efficient in practice, since it can be implemented using highly optimized matrix multiplication code.
> 
> 翻译：两种最常用的注意力函数是加性注意力\[2\]和点积（乘法）注意力。点积注意力与我们的算法相同，除了有一个缩放因子。加性注意力使用具有单个隐藏层的前馈网络来计算兼容性函数。虽然两者在理论上的复杂度相似，但在实践中，点积注意力要快得多，也更节省空间，因为它可以利用高度优化的矩阵乘法代码来实现。

点积的结果会除以![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5SibeY7WyllK8toGmibwR6NveTHk5TtCr4Dw2QCJWDXv2obz5VIa23Btw/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)的平方根，来对点积的结果进行缩放，确保数值稳定，这一段在原文中也有表述。通过将“探险家”与每一个词元都进行一次计算，就可以得到一个向量，向量中的每一个元素代表着“探险家”与对应词元的点积的值。之后进行掩码（Mask），在前文中提到，在因果解码器中，当前词元是无法看到自身之后的词元的，所以需要将当前词元之后的所有点积置为负无穷，以便使其在归一化后的占比为零。﻿

最后对点积softmax进行归一化，就可以得到一个总和为1的一个比例分布，这个分布就叫“注意力分布”。下面表格中“探险家”所对应的那一列，就是“探险家”的注意力分布，代表着从“探险家”的视角出发，每一个词元对于自身内容理解的重要程度。

<table><colgroup><col width="159"><col width="159"><col width="160"><col width="159"><col width="161.99993896484375"></colgroup><tbody><tr><td rowspan="1" colspan="1"><section><span>Key\Query&nbsp;</span></section></td><td rowspan="1" colspan="1"><section><span>Q 有一天&nbsp;</span></section></td><td rowspan="1" colspan="1"><section><span>Q 一位&nbsp;</span></section></td><td rowspan="1" colspan="1"><section><span>勇敢的&nbsp;</span></section></td><td rowspan="1" colspan="1"><section><span>探险家&nbsp;</span></section></td></tr><tr><td rowspan="1" colspan="1"><section><span>K 有一天 &nbsp; &nbsp;</span></section></td><td rowspan="1" colspan="1"><section><span>1</span></section></td><td rowspan="1" colspan="1"><section><span>0.13</span></section></td><td rowspan="1" colspan="1"><section><span>0.03</span></section></td><td rowspan="1" colspan="1"><section><span>0.02</span></section></td></tr><tr><td rowspan="1" colspan="1"><section><span>K 一位 &nbsp; &nbsp; &nbsp;</span></section></td><td rowspan="1" colspan="1"><section><span>0</span></section></td><td rowspan="1" colspan="1"><section><span>0.87</span></section></td><td rowspan="1" colspan="1"><section><span>0.1</span></section></td><td rowspan="1" colspan="1"><section><span>0.05</span></section></td></tr><tr><td rowspan="1" colspan="1"><section><span>K 勇敢的 &nbsp; &nbsp;</span></section></td><td rowspan="1" colspan="1"><section><span>0</span></section></td><td rowspan="1" colspan="1"><section><span>0</span></section></td><td rowspan="1" colspan="1"><section><span>0.87</span></section></td><td rowspan="1" colspan="1"><section><span>0.08</span></section></td></tr><tr><td rowspan="1" colspan="1"><section><span>K 探险家 &nbsp; &nbsp;</span></section></td><td rowspan="1" colspan="1"><section><span>0</span></section></td><td rowspan="1" colspan="1"><section><span>0</span></section></td><td rowspan="1" colspan="1"><section><span>0</span></section></td><td rowspan="1" colspan="1"><section><span>0.85</span></section></td></tr></tbody></table>

一个可能的Softmax后的概率分布（每一列的和都为1，灰色代表掩码）

最后，将Softmax后获得的注意力分布，分别与每一个K对应的V相乘，通过注意力权重加权求和，就可以得到一个向量，称为上下文向量。它就是![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU51eWqicBibjpdhPCU1iaJAMgV9zqpHfXJuD7O31FoTNmy6DFgAWicnibI4qg/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)求解出来的值，其中包含了序列中与当前元素最相关的信息，可以认为是模型在结合了上下文的信息后，对序列内容的理解的一种表示形式。﻿

### 多头注意力﻿

关于什么是多头注意力，为何使用多头注意力，论文中有这样一段描述，各位一看便知。

> 原文：Instead of performing a single attention function with dmodel-dimensional keys, values and queries, we found it beneficial to linearly project the queries, keys and values h times with different, learned linear projections to dk, dk and dv dimensions, respectively. On each of these projected versions of queries, keys and values we then perform the attention function in parallel, yielding dv-dimensional output values. These are concatenated and once again projected, resulting in the final values, as depicted in Figure 2. Multi-head attention allows the model to jointly attend to information from different representation subspaces at different positions. With a single attention head, averaging inhibits this.
> 
> 翻译：我们发现，与其执行一个具有dmodel维键、值和查询的单一注意力函数，不如将查询、键和值线性投影h次，使用不同的、学习得到的线性投影到dk、dk和dv维度。然后我们并行地在这些投影版本的查询、键和值上执行注意力函数，得到dv维的输出值。这些输出值被连接起来，再次被投影，得到最终的值，如图2所示。多头注意力允许模型同时关注不同表示子空间中的信息，以及不同位置的信息。而单一注意力头通过平均操作抑制了这一点。﻿

多头注意力的“多头”，指的是点积注意力函数实例有多种。不同的头，他们的![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5LqQVOADT5RnvQPjsnpwA3ia8MEG4iaeJ3wDPRXFxd5Tticjkz4EamP3jQ/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)都可能不同，这可能意味着词元会从不同角度去发问，同时表达出不同角度的特征，比如一个头可能专注于捕捉语法信息，另一个头可能更关注语义信息，还有一个头可能更关注情感分析。这对于充分捕获上下文信息，尤其是在处理复杂的序列数据时，变得更加强大和灵活。﻿

这种优化并没有使计算的复杂度升高，论文中特别提到

> 原文：In this work we employ h = 8 parallel attention layers, or heads. For each of these we use dk = dv = dmodel/h = 64. Due to the reduced dimension of each head, the total computational cost is similar to that of single-head attention with full dimensionality.

> 翻译：在这项工作中，我们使用了 h=8 个并行的注意力层，或者说是“头”（heads）。对于每一个头，我们设置 dk=dv=dm/h=64 。这里的dk和dv分别代表键（keys）和值（values）的维度，而![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU59MVAwD1iaWrZq6fRhjBGnNvHxVk825aKWkY7R4yWiaZEmojqo3KW9oYA/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)是模型的总维度。

进一步举个例子来解释一下这段内容。下面是节选的GPT-3模型配置，N代表注意力头数，H代表隐藏状态的大小（参数量）。在单头注意力的情况下，每个头都是12288维，而在多头注意力的情况下，头与头之间会均分参数量，每个头的参数量只有12288/96 =128维，并且不同头的注意力计算都是并行的。﻿

![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5X1o95AjkUiavfsicHeMrmwHfQvr8Tia0abpRoBNPO35kibMvSb2mXpmy8Q/640?wx_fmt=other&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

GPT-3模型配置，引自《A Survey of Large Language Models》

如今，这种设计随着技术的发展有所演进，Q的头数在标准的MHA下，通常与KV的头数相同，然后目前主流的大模型都进行KV缓存的优化，Q和KV的头数可能并不相同。比如我们常用的Qwen2-72B，其隐藏层有8192个参数，有64个Q和8个KV头，每个头的参数量是128。数据来自于Qwen2的技术报告（如下图），具体技术细节在后续GQA部分会有详细说明。﻿

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5G0qQwCaBfoBERIq2SGJfN4fJah8x6V5LXlC2Zqkx2L7HibuSmY1chibQ/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

Qwen2系列模型参数，引自《QWEN2 TECHNICAL REPORT》﻿

回到计算过程中，多头注意力，会在每个头都按照缩放点积注意力的方式进行运算后，将他们产生的上下文向量进行连接，基于输出投影矩阵![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5tlicCGSicSloF4MOMcPjicREn6EgYOqRY6qq3nibDRxXr4sUnOFVZgXYgw/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)进行变换。这个过程是为了综合不同注意力头所提供的信息，就像是综合考虑96个人的不同意见，并形成最终的结论。

这种综合的过程并不是简单地求平均值，而是通过![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5tlicCGSicSloF4MOMcPjicREn6EgYOqRY6qq3nibDRxXr4sUnOFVZgXYgw/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)进行的连接操作。

![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5CBHP8ce57fmFE7tMXBPT5R4mFJZtAOUneP3oY0PiaRsXlY439maZnYA/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)﻿

具体来说，首先多头注意力会获取每个单头注意力所提供的上下文向量，并在特征维度上进行连接，形成一个更长的向量，对应公式中的![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5ok2esF9mAqbu7iaGr8RuQn5FcAPXZKXWtn2G0396oaQcIMt43reXmzw/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)，其中h是注意力头数。拼接后的矩阵的维度是隐藏层的维度，比如GPT-3有96个头，每个头有128\*12288维，那么拼接后形成的就是一个12288\*12288维的矩阵。![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5tlicCGSicSloF4MOMcPjicREn6EgYOqRY6qq3nibDRxXr4sUnOFVZgXYgw/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)是模型在训练过程中学习到的关键组成部分，将拼接后的矩阵向量基于该矩阵做一次线性变换，有助于模型在多头注意力的基础上进一步优化特征表示，提高模型的整体性能。

至此，标注的MHA的相关内容就结束了。接下来MHA层的输出，会传递到下一层，可能是FFN，也可能是MoE，取决于具体的模型。在此之前，先来看一下大模型对注意力层的优化。

**KV Cache**

前文提到，因果解码器的特点，是在生成每个词元时，只能看到它之前的词元，而不能看到它之后的词元。也就是说，无论模型在自回归过程中生成多少词元，此前已经生成的词元对上下文内容的理解，都不会发生任何改变。因此我们在自回归过程中，不需要在生成后续词元时重新计算已经生成的词元的注意力。﻿

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5RUbxibBIrVMt45FhNiadr9bmEO7eRlakkBlglLCnEkhW05OpeWlTJJLQ/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

我是真的爱喝咖啡﻿

但是，新生成的词元的注意力需要计算，这会涉及新生成的词元的Q与其它词元的K计算点积，并使用其它词元的V生成上下文向量。而此前生成的词元K、V，实际上始终不会改变，因此我们可以将他们缓存起来，在新生成的词元计算注意力的时候直接使用，避免重复计算，这就是KV缓存。如上图，已经生成的词元“我”、“爱”、“喝”都不会重新计算注意力，但是新生成的“咖啡”需要计算注意力，期间我们需要用到的是“咖啡”的 Q，和“我”、“爱”、“喝”的K、V。

KV缓存的核心思想是：

- 缓存不变性：在自回归生成过程中，已经生成的词元的键（Key，K）和值（Value，V）不会改变。
- 避免重复计算：由于K和V不变，模型在生成新词元时，不需要重新计算这些已生成词元的K和V。
- 动态更新：当新词元生成时，它的查询（Query，Q）会与缓存的K进行点积计算，以确定其与之前所有词元的关联。同时，新词元的K和V会被计算并添加到缓存中，以便用于下一个词元的生成。

在使用KV Cache的情况下，大模型的推理过程常被分为两个阶段

- 预填充阶段（Prefill）：模型处理输入序列，计算它们的注意力，并存储K和V矩阵到KV Cache中，为后续的自回归过程做准备。
- 解码阶段（Decode）：模型使用KV缓存中的信息，逐个生成输出新词元，计算其注意力，并将其K、V添加到KV Cache中。

其中预填充阶段是计算密集型的，因为其涉及到了矩阵乘法的计算，而解码阶段是内存密集型的，因为它涉及到了大量对缓存的访问。缓存使用的是GPU的显存，因此我们下一个面临的问题是，如何减少KV Cache的显存占用。﻿

**MQA**

2019年，Google团队发布了论文《Fast Transformer Decoding: One Write-Head is All You Need》，并提出了多查询注意力的这一MHA的架构变种，其全拼是Multi-Query Attention，简称MQA，GPT-4模型就是采用的MQA来实现其注意力层。﻿

Google为何要提出，论文中提到

> 原文1：Transformer relies on attention layers to communicate information between and across sequences. One major challenge with Transformer is the speed of incremental inference. As we will discuss, the speed of incremental Transformer inference on modern computing hardware is limited by the memory bandwidth necessary to reload the large "keys" and "values" tensors which encode the state of the attention layers.
> 
> 原文2：We propose a variant called multi-query attention, where the keys and values are shared across all of the different attention "heads", greatly reducing the size of these tensors and hence the memory bandwidth requirements of incremental decoding.
> 
> 翻译1：Transformer依赖于注意力层来在序列之间和内部传递信息。Transformer面临的一个主要挑战是增量推理的速度。正如我们将要讨论的，现代计算硬件上增量Transformer推理的速度受到重新加载注意力层状态所需的大型“键”和“值”张量内存带宽的限制。
> 
> 翻译2：我们提出了一种变体，称为多查询注意力（Multi-Query Attention），其中键（keys）和值（values）在所有不同的注意力“头”（heads）之间共享，大大减少了这些张量的大小，从而降低了增量解码的内存带宽需求。

增量推理（Incremental Inference）是指在处理序列数据时，模型逐步生成输出结果的过程。张量其实就是多维数组，在注意力层主要指的是各个与注意力有关的权重矩阵。不难看出，Google团队注意到了K、V所带来的巨大内存带宽占用，通过MQA将K、V在不同注意力头之间共享，提高了模型的性能。

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5OAYyUqCeT3kCEztEMxW5vMc16g7pnEB6OIFBRiaAws4kTrKzS2TQbvg/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

MHA、GQA、MQA的比较，引自《GQA: Training Generalized Multi-Query Transformer Models from Multi-Head Checkpoints》

我们用GPT-3举例，它有96个自注意力头。那么在传统MHA中，生成一个新的词元，都需要重新计算96个Q、K、V矩阵。而在MQA中，只需要计算96个Q矩阵，再计算1次K、V矩阵，再将其由96个头共享。每次Q、K、V的计算都需要消耗内存带宽，通过降低K、V的计算次数，可以有效优化模型的解码速度。

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU52XsyficcGKeINZwPSkV0grYv0euwbuE20kDxTomUUDfqiaNDibOQozKSQ/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

在WMT14英德（English to German）翻译任务上的性能比较，来自《Fast Transformer Decoding: One Write-Head is All You Need》

BLEU是一种评估机器翻译质量的自动化指标，分数越高表示翻译质量越好。根据论文中对性能比较的结果，MQA确实相对于MHA，在翻译效果上的性能有所下降，但是相对于其他减少注意力头数量等替代方案而言，效果仍然很好。

实际上由于KV缓存的使用，MQA降低的主要资源消耗，并不是内存带宽，而是内存占用，也就是KV缓存的大小。﻿

**GQA**

GQA，来自于Google团队的2023年的论文《GQA: Training Generalized Multi-Query Transformer Models from Multi-Head Checkpoints》，GQA的全拼是Grouped Query Attention（分组查询注意力），被包括Llama3、Qwen2在内的众多主流模型广泛采用。

论文中提到

> 原文：However, multi-query attention (MQA) can lead to quality degradation and training instability, and it may not be feasible to train separate models optimized for quality and inference. Moreover, while some language models already use multiquery attention, such as PaLM (Chowdhery et al., 2022), many do not, including publicly available language models such as T5 (Raffel et al., 2020) and LLaMA (Touvron et al., 2023).
> 
> 翻译：然而，多查询注意力（MQA）可能导致质量下降和训练不稳定性，并且可能不切实际去训练分别针对质量和推理优化的独立模型。此外，虽然一些语言模型已经采用了多查询注意力，例如PaLM（Chowdhery等人，2022年），但许多模型并没有采用，包括公开可用的语言模型，如T5（Raffel等人，2020年）和LLaMA（Touvron等人，2023年）。

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5P6a8sliasRoiadHLyHteY2rc2MOU5otuLyyzP2Lhm0xpRic2LKhHJcA0A/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

MHA、MQA、GQA的性能比较，引自《GQA: Training Generalized Multi-Query Transformer Models from Multi-Head Checkpoints》﻿

它的本质其实是对MHA、MQA的一种折中，在显存占用和推理性能上的一种平衡。上图是对MQA、GQA、MHA三种注意力模式下模型性能的比较。（XXL代表Extra Extra Large，超大型模型，具备最多的参数量，Large代表大型模型，其参数量在标准模型和XXL之间）。

前馈神经网络

FFN，全拼Feed-Forward Network，前馈神经网络。FFN层，通过对自注意力层提供的充分结合了上下文信息的输出进行处理，在高维空间中进行结合训练获得的特征和知识，获得新的特征。FFN因其简单性，在深度学习领域备受欢迎，其原理也相对更容易解释。﻿

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5tIZnuTOTxbAd1BEY5W2iae0fs7q9zEBHSfz07y9eNP8ISD4NTebzWbA/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

FFN层的处理过程，引自3Blue1Brown的视频《直观解释大语言模型如何储存事实》﻿

在Transformer中，FFN层由两个线性变换和一个激活函数构成，它的处理过程是词元维度的，每个词元都会并行地进行计算，如上图，因此在学习FFN层的处理过程时，我们只需要分析单个词元的处理过程。

![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5yfu9RhkQ4RAUwSicrsUficibIVWuUzgUEkQRdChrMBmT8N7NqYR4SGTew/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)﻿

这个过程同样可以用一个简洁的公式来表示，如上图，让我们来逐步解读一下。

首先，X是输入向量，代表了已经充分结合上下文信息的单个词元，它由自注意力层提供，其维度就是隐藏层的维度，比如GPT-3中是12288。我们将接收到的词元，首先与通过模型训练获得的矩阵![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5LjQlcFGkEUOiajG2FqicwA8NmgicCNWr67ZuLRrbzZ21fd3PWx95dUUyQ/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)相乘，进行一次线性变换，进而将输入向量映射到一个更高维度的空间，这个空间通常被称为FFN的“隐藏层”。这个过程可以认为是对输入的“扩展”，目的是使FFN的隐藏层能够表达更多的特征，这一层的维度通常比自注意力层提供的输出维度大得多，比如GPT-3中它的维度是49152，刚好是自注意力层输出维度的四倍。

举个例子，假设在自注意力层产生的输入中，模型只能了解到词元的语法特征、语义特征，比如“勇敢的”，模型能感知到它是“形容词”、代表“勇敢”。那么在经过这次线性变换后，模型通过扩充维度，就能感知到其“情感特征”，比如“正向”、“积极”。

b代表bias，中文意思是偏置、偏见、倾向性，它也是通过模型训练获得的，在模型的正向推理过程中可以视为一个常数。神经网络中的神经元可以通过公式![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5jMyZ4QEJxqWrLY6BG3rgmFdDlOEQJN5Ca2f2yz8YEOX9gxcM6sOYQQ/640?wx_fmt=jpeg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)来表示，b在其中可以控制函数到原点的距离，也叫函数的截距。通过引入bias，可以避免模型训练过程中的过拟合，增强其泛化性，以更好地适应不同的数据分布，进而提高预测的准确性。

![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU57L0sp8dZKkblnWb9XfJNAfGkM9BJ42mFfO8loT3UblnXvDfsRot9Ig/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)，西格玛，代表激活函数。激活函数的作用，是为模型引入非线性的因素，作为一个“开关”或者“调节器”，来控制信息在神经网络中的传递方式，即某些特征是否应当被传递到下一层。这种非线性的因素，使得模型能够学习和模拟复杂的函数映射关系。并通过让模型专注于那些对当前任务更有帮助的正向特征，来让模型能够更好的选择和组合特征。

举个例子，我们通过线性变换，获得了关于输入内容的大量特征信息，但其中一部分信息相对没那么重要或毫不相关，我们需要将他们去掉，避免对后续的推理产生影响。比如“我”这个词，自身的语法特征很清晰，我们要保留，但是其并没有什么情感特征，因此我们要将与“我”的情感特征相关的信息去除。

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU57wEplcibCCiaVbIxicc75XjEfZCBppOJmK7VibLzF95cxbSMtZnC3G1qTw/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

不同函数的曲线对比，引自《A Survey of Large Language Models》

当然，“去除”这个词其实并不准确，应当叫做“抑制”。ReLU作为一种激活函数，会将所有相乘后结果为零的部分去除，只保留所有结果为正的信息，我们可以认为是“去除”。不过ReLU在目前主流的大模型中并不常用，比如Qwen、Llama等模型选择使用SwiGLU，GPT选择GeLU，他们的曲线相对更加平滑，如上图。不同激活函数的选择，是一种对于模型的非线性特性和模型性能之间的权衡，类似于ReLU这种函数可能会导致关闭的神经元过多，导致模型能够感知到的特征过少，变得稀疏。

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5eJhO1yoIp5swODpQjsnwmH4DYGxCdRNMNYrtGjnd12Fn4bdoNCr1fQ/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

FFN层的处理过程，引自3Blue1Brown的视频《直观解释大语言模型如何储存事实》﻿

在经过激活函数进行非线性变换处理后的向量，会再次通过矩阵![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5ibkec7ApZNK05Ij27juUoSCxiaMcOiaibgbEsDeoHpticLSd8icGp7CkhRsw/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)进行第二次线性变换，将高维空间中的向量，重新映射回原始维度，并追加第二个偏置![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5FsY0kcSqeJRV7iczTD6sXzaib6ADNPfPEicBXFbkNVZdABQ6W5Ks9s1cA/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)，来形成模型的输出。这种降维操作，一方面使得FFN层的输出能够与下一层（自注意力层或输出层）的输入维度相匹配，保持模型的深度不变，减少后续层的计算量，另一方面模型可以对升维后学习到的特征进行选择和聚焦，只保留最重要的信息，这有助于提高模型的泛化能力。到此单次FFN层的执行过程就讲述完毕了，整体过程可以参考上图。﻿

总结一下，向上的线性变换，使得词元能够表达出更多的特征，激活函数通过非线性因素，来增强模型对特征的表达能力，向下的线性变换，会将这些特征进行组合，这就是FFN层中模型的“思考”过程。

另外要说明，向量的特征，并不会像我们前面举例的那样简单地可以概括为“语义”、“语法”、“情感”特征。比如在模型训练过程中，模型可能学习到“美国总统”和“川普”之间具有关联性，“哈利”与“波特”之间具有关联性，“唱”、“跳”、“Rap”与“篮球”之间具有关联性，这些关联性很难用简单的语言来表达清楚，但它们也实实在在地形成了“川普”、“波特”、“篮球”的某些特征。

有人会说“训练大模型的过程就像炼丹”，这其实是在描述模型内部的黑盒性。而我们使用大模型时，也要避免工程化的思维，认为大模型一定会按照预设的规则去执行，这其实并不尊重模型本身的特性。因为模型的推理过程，不仅仅受到输入（包括提示词以及模型自回归过程中不断产生的输出）的影响，还会受到训练数据、模型架构、以及训练过程中的超参数的影响。但我们可以在理解了注意力机制后通过设计良好的提示词，在理解了模型的思考过程后通过进行模型的微调或增强学习，来驯化大模型。﻿

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5rws9gM09KTFDPVxOWt4fia4XWMVrETKUSYya8E5ib9f9ibB1wiaiaht56lA/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

FFN层的神经网络结构图

最后，我想讲一个概念。在各种技术报告中，我们常会看到一个词——“稠密模型”，它指的是模型在处理任务时，模型的每个神经元都彼此相连，所有参数都共同参与计算的模型。如上图，FFN层无论是在向上的线性变换还是向下的线性变换的过程中，每一个神经元都彼此相连，因此这两层线性变换其实就是FFN层的两层稠密层，FFN层也就可以视为稠密模型的一种形式。﻿

稠密模型由于其参数量很大，能够捕捉更丰富的特征和复杂的模式，但这也导致其较高的训练和推理成本，且在数据集规模较少时，尝试去拟合那些不具有普遍性的噪声，导致模型的过拟合，降低模型的泛化性。﻿

与“稠密模型”相对应的是“稀疏模型”，其核心思想是利用数据的稀疏性，即数据中只有少部分特征是重要的，大部分特征都是冗余或者噪声。MoE就是一种典型的稀疏模型，目前在GPT-4，以及Qwen2的部分模型等众多大语言模型上，被用于替代FFN层。

输出层

在最后的最后，模型在经过多轮的隐藏层的计算后，获得了最终的隐藏状态。输出层，则负责将这个隐藏状态，转换为下一个词元的概率分布。

![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5eBLoslApMxjPTMkD2hiciaRsGkVN8VjbaRxzlavO9tPzHJQLicrX2GTYw/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

X是隐藏层的最终输出，其维度是![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5Rc10H95MZMv5lbgh1a0C17wN2KaR5iagAY7GdBwoQPtWIrnibqg9SMxg/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)，![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5OPnJSlMKzhohI4p1uRnhuyLW5LJS07BTtJZgR1NfVdhM2udaI2D7AQ/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)是输入序列的长度，![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5LicW1ERf8cC19q9XfwnbneGDVryE7JHhAEC23picRvBia343EXxlUoRicg/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)是隐藏层的维度。而![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5ib7kQssYIMwOrM7CNdXGJiax9Qx0Jpgcorw03JAnicnmVDtZ5EeAcKCRw/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)是通过训练获得的权重矩阵，其维度是![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5Nr3w4hpvR29rso7VHVvkjYtRLicaK0bgOCw37MdExHABrnjEOibzrW3w/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)，![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5lmRRJFLVWnTYpIPHXQg9zZfgqnp8Jibz3N3dJzvib7X0ViaZgFmJwnhDg/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)是词汇表的大小，比如Qwen2-72B的词汇表大小是151646。通过矩阵相乘，再与偏置bias相加，就可以将隐藏状态转换为输出词元的分数，也就是![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5MYMVTYwocU0APjPJAShIBYdyUEaoSM8N8OQkjZOSAovUwp23luMULQ/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)，它代表了模型经过“思考”后的特征，用哪个“词元”来形容更合适。

![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5ic9MJh03CTJK3MwVrLZezrw7tYpDIEyg7fht2rjor3Cjhl5VENFBYibg/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

然后将![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5MYMVTYwocU0APjPJAShIBYdyUEaoSM8N8OQkjZOSAovUwp23luMULQ/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)通过归一化操作（Softmax），转换为词元的概率，在此基础上结合解码策略，就可以选择具体的下一个词元进行输出。这个公式非常简单，其中![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5IicXtic1CYI8P3ibDenuSQv0yO958N3gx8ranicFjqA9AqJ7Hs5J7bZ6dw/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)就是某一个词元的分数。通过将指数函数应用于![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5IicXtic1CYI8P3ibDenuSQv0yO958N3gx8ranicFjqA9AqJ7Hs5J7bZ6dw/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)形成![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5d2Z4xgQHlRVW9RghQLoK74bPOeZl8fYpQZqRyh2MgoWxtjaa6dwFcg/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)，不仅可以确保词元分数的数值为正（便于转换为概率），还能增加不同分数之间的差异性。最后将单个词元的分数，与所有词元的分数之和相除，就能得到单个词元的概率，如此就能获得词汇表中每个词元的概率分布。﻿

此时，你可能会好奇，那个![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5CLKVXrOQqCIVKjF4FA6LOmVftJtaJjFh2SOExtVlt10ictuv2MWBzIQ/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)是什么？![图片](https://mmbiz.qpic.cn/mmbiz_png/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5CLKVXrOQqCIVKjF4FA6LOmVftJtaJjFh2SOExtVlt10ictuv2MWBzIQ/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)指的是temperature，没错就是我们在神机平台上常常见到的那个大模型节点的参数，它通过影响softmax后的概率，来影响最终输出概率分布的平滑程度。类似作用的还有top\_p参数，这也是在模型的输出层起作用。这些参数很重要，并且我们在开发过程中会常常遇到。

至此，模型的推理过程核心内容就都讲述清楚了，当然还有很多内容，比如分词、词嵌入、位置编码、残差连接、层归一化、以及FFN之后又兴起的MoE架构等内容，篇幅所限，此处不再展开。

写在最后﻿

**提示工程**

Prompt Engineering（提示工程），则是指针对特定任务设计合适的任务提示（Prompt）的过程。﻿

在大模型的开发和性能优化的过程中，OpenAI建议将提示工程作为大模型应用的起点，从上下文优化、大模型优化两个角度思考，这两种角度对应了两个方向：提示工程、微调。

- 提示工程：通过精心设计的输入提示（Prompt）来引导预训练模型生成期望的输出，而无需对模型的权重进行调整。
- 微调：微调是在预训练模型的基础上，使用特定任务的数据进一步训练模型，以调整模型权重，使其更好地适应特定任务。

提示工程侧重于优化输入数据的形式和内容，以激发模型的潜在能力，来提高输出的准确性和相关性。它可以快速适应新任务。而微调可以使模型更深入地理解特定领域的知识和语言模式，进而显著提高模型在特定任务上的性能，但其在灵活性上相对较弱，训练依赖于计算资源和高质量的标注数据。对于高德的业务团队而言，会相对更侧重灵活性来适应快速变化的业务需求，因此基于提示工程进行优化的方法已经成为我们使用大语言模型解决下游任务的主要途径。﻿

实际上不仅仅对于业务团队，对于任何团队，在使用大模型时，都应从提示工程开始。OpenAI针对需要提供给大模型额外知识的场景提供了一份合理的优化路线图（如下图）：从基础Prompt开始，通过提示工程优化Prompt，接入简单RAG，进行模型微调，接入高级ARG，最后带着RAG样本进行模型微调。﻿

﻿﻿![图片](https://mmbiz.qpic.cn/mmbiz_jpg/Z6bicxIx5naJW6JVtVFFSliboxWAjbPLU5wCdUCRarXviaVnq6Sibco4b9TKeYZdlAILJpe4hjY0nzatwEoIjibqdLw/640?wx_fmt=jpeg&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

OpenAI: How to maximize LLM performance（https://humanloop.com/blog/optimizing-llms）

通过提示工程我们可以应对大部分的业务场景，如果性能不够，第一件事是要考虑提示工程，再考虑其他手段。RAG的功能是上下文增强，其输出结果是提示词的一部分，因此RAG也可以视为提示工程的一部分。而对于有垂类场景模型需求的场景，也需要通过提示工程来获取高质量的用例数据，来进行模型微调，再基于微调后模型产出的更高质量用例，正向迭代来进一步优化性能。

对于开发人员而言，如果希望提升驯化大模型的能力，我建议从提示工程开始。这包括了提示词结构化（LangGPT等）、提示设计方法（如OpenAI提出的六大原则）、提示框架（ReACT等）、提示技术（COT、Few-Shot、RAG等）、Agent的概念和架构。﻿

**推荐阅读**

首先推荐大家观看网上的教程，包括吴恩达的深度学习、台湾教授李宏毅的课程，以及3Blue1Brown的视频（可视化做得非常棒），各位可以在B站找到相应的免费学习资源。

然后推荐一本书，《大语言模型综述》，该书来自中国人民大学，该书是文章《A Survey of Large Language Models》整理成册出版的中文版本，系统性地梳理了大型语言模型的研究进展与核心技术，并讨论了大量的相关工作，对于系统性地学习大语言模型有很强的指导意义。

直接阅读各类论文也很有帮助，这包括经典论文以及各类模型的技术报告等。

最后，AI是学习的好帮手，在学习期间请频繁的、频繁的、频繁的针对性对它提问，甚至直接将截图或者论文抛过去让它帮你分析，对于学习非常有帮助。同时这些Agent都会进行联网搜索，并且给出搜索过程中找到的原文链接，可以方便我们快速找到学习内容。

---