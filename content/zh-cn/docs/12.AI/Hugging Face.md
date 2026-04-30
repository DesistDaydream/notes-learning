---
title: Hugging Face
linkTitle: Hugging Face
weight: 70
---

# 概述

> 参考：
>
> - [GitHub 组织，huggingface](https://github.com/huggingface)
> - [官网](https://huggingface.co/)
> - [Wiki, Hugging_Face](https://en.wikipedia.org/wiki/Hugging_Face)
> - [知乎，Huggingface 超详细介绍](https://zhuanlan.zhihu.com/p/535100411)

Hugging Face 即是一个工具包的集合，也是一个社区。

- 在 2017 年在 GitHub 上开源了非常著名的 [Transformer](/docs/12.AI/Machine%20learning/Transformer/Transformer.md) 库。
- 在 2019 年推出了 Hugging Face **Hub**，一个用于共享和加载 预训练模型、数据集 的平台。我们可以通过 Hub 上传/下载 模型，类似于大模型生态的 [GitHub](/docs/2.编程/Programming%20tools/SCM/GitHub/GitHub.md)
- 后续也陆续开发了其他的工具包，例如 tokenizers、datasets、huggingface_hub、accelerate 和 peft，以及一个公共的推理 API，为机器学习开发者提供了更多的便利和资源。

它提供了多种工具包，例如：

- **transformers**：一个用于构建和使用预训练语言模型的 Python 库，支持 PyTorch 和 TensorFlow。
- **tokenizers**：一个用于快速和高效地创建和使用词法分析器的 Python 库。
- **datasets**：一个用于轻松共享和加载数据集和评估指标的 Python 库。
- **huggingface_hub**：一个用于创建、删除、更新和检索仓库信息的 Python 库，也可以从仓库下载文件或将它们集成到你的库中。
- **accelerate**：一个用于轻松地训练和使用 PyTorch 模型的 Python 库，支持多 GPU、TPU 和混合精度。
- **peft**：一个用于实现参数高效微调（PEFT）方法的 Python 库，可以有效地将预训练语言模型（PLMs）适应于各种下游应用，而无需微调所有模型参数6。

使用 Hugging Face 可以帮助我们在 NLP 领域进行创新和探索。

Hugging face 起初是一家总部位于纽约的聊天机器人初创服务商，他们本来打算创业做聊天机器人，然后在github上开源了一个Transformers库，虽然聊天机器人业务没搞起来，但是他们的这个库在机器学习社区迅速大火起来。目前已经共享了超100,000个预训练模型，10,000个数据集，变成了机器学习界的github。

# 关联文件与配置

> 参考：
>
> - [官方文档，安装 - 缓存目录](https://huggingface.co/docs/transformers/main/en/installation#cache-directory)

**~/.cache/huggingface/** # Hugging Face 缓存目录。使用 `HF_HOME` 环境变量修改。

> [!Tip] **Windows** 使用 Powershell 命令 `[Environment]::SetEnvironmentVariable("HF_HOME", "D:\appdata\huggingface", "User")` 修改

- **./hub/** # 模型保存目录中，e.g. 自己训练的、从 HF 网站下载的、etc. 。使用 `HF_HUB_CACHE` 环境变量修改。

使用如下代码查看当前 HF 的缓存路径

```python
from huggingface_hub import constants
# 查看当前生效的 HF 根目录
print(f"HF Home: {constants.HF_HOME}")
# 查看模型具体的缓存目录
print(f"Hub Cache: {constants.HF_HUB_CACHE}")
```

## 环境变量

> 参考：
>
> - [官方文档，参考 - 环境变量](https://huggingface.co/docs/huggingface_hub/package_reference/environment_variables)

**HF_HOME** # Hugging face 相关生态（e.g. 模型、工具、etc.）所用的目录。

**HF_ENDPOINT** # 配置 Hub 的基础 url。可以通过修改这个环境变量，让 HuggingFace 的 CLI 使用私有的 Hub。甚至可以通过这种方式实现从镜像站下载模型（中国专属 `(～￣▽￣)～`）

> [!Tip] 这个环境变量最后一次出现在 [0.17.3 版本](https://huggingface.co/docs/huggingface_hub/v0.17.3/en/package_reference/environment_variables )，之后的版本的文档不记录该变量了，但是依然可用。

**HF_TOKEN** # 访问 HuggingFace 的 Hub 的认证信息。默认位于 `$HF_HOME/token` 中。

# CLI

> 参考：
>
> - [官方文档，参考 - CLI](https://huggingface.co/docs/huggingface_hub/package_reference/cli)

安装 Hugging Face 

```bash
uv tool install huggingface_hub
```

## download

**OPTIONS**

- **--local-dir**(STRING) # 指定模型包下载的文件储存的目录，代替默认的 `${HF_HOME}` 目录。查看 [这里](https://huggingface.co/docs/huggingface_hub/guides/download#download-files-to-a-local-folder) 了解更多详细信息。
    - Notes: 使用了这个选项后，储存目录下的文件是直观的没有连接的，不像默认的 `${HF_HOME}` 目录的组织方式，分为 blobs, refs, snapshots, etc. 多个目录，然后模型包的所有文件都通过软链接指向 blobs。而是直接在 `--local-dir` 存放相关的 XX.safetensors, config.json, tokenizer_config.json, etc. 文件

# 最佳实践

## 下载模型及镜像使用

https://zhuanlan.zhihu.com/p/663712983

使用 HF_ENDPOINT 环境变量定义下载位置

可用的镜像

- https://hf-mirror.com

从 HF 官网，找到模型全称

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ai/hf_official_model_name_demo_1.png)

下载 Qwen3.6 模型。在 CLI 中使用从官网看到的模型名称，即可将模型下载到 `${HF_HOME}` 目录中

```bash
HF_ENDPOINT=https://hf-mirror.com hf download Qwen/Qwen3.6-35B-A3B
```

- 使用 --local-dir 目录将模型下载到 `./models/qwen3.6/` 目录中

```bash
HF_ENDPOINT=https://hf-mirror.com hf download Qwen/Qwen3.6-35B-A3B --local-dir ./models/qwen3.6
```