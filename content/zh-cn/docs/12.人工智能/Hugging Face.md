---
title: "Hugging Face"
linkTitle: "Hugging Face"
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 组织，huggingface](https://github.com/huggingface)
> - [官网](https://huggingface.co/)
> - [WIki，Hugging_Face](https://en.wikipedia.org/wiki/Hugging_Face)
> - [知乎，Huggingface 超详细介绍](https://zhuanlan.zhihu.com/p/535100411)

Hugging Face 即是一个工具包的集合，也是一个社区。

- 在 2017 年在 GitHub 上开源了非常著名的 [Transformers](/docs/12.人工智能/机器学习/Transformers.md) 库。
- 在 2019 年推出了 Hugging Face Hub，一个用于共享和加载预训练模型和数据集的平台
- 后续也陆续开发了其他的工具包，例如 tokenizers、datasets、huggingface_hub、accelerate和peft，以及一个公共的推理API，为机器学习开发者提供了更多的便利和资源。

它提供了多种工具包，例如：

-   **transformers**：一个用于构建和使用预训练语言模型的 Python 库，支持 PyTorch 和 TensorFlow。
-   **tokenizers**：一个用于快速和高效地创建和使用词法分析器的 Python 库。
-   **datasets**：一个用于轻松共享和加载数据集和评估指标的 Python 库。
-   **huggingface_hub**：一个用于创建、删除、更新和检索仓库信息的 Python 库，也可以从仓库下载文件或将它们集成到你的库中。
-   **accelerate**：一个用于轻松地训练和使用 PyTorch 模型的 Python 库，支持多 GPU、TPU 和混合精度。
-   **peft**：一个用于实现参数高效微调（PEFT）方法的 Python 库，可以有效地将预训练语言模型（PLMs）适应于各种下游应用，而无需微调所有模型参数6。

使用 Hugging Face 可以帮助我们在 NLP 领域进行创新和探索。

Hugging face 起初是一家总部位于纽约的聊天机器人初创服务商，他们本来打算创业做聊天机器人，然后在github上开源了一个Transformers库，虽然聊天机器人业务没搞起来，但是他们的这个库在机器学习社区迅速大火起来。目前已经共享了超100,000个预训练模型，10,000个数据集，变成了机器学习界的github。