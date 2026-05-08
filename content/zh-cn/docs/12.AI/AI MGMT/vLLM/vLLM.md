---
title: vLLM
linkTitle: vLLM
created: 2026-05-07T17:03
weight: 1
---

# 概述

> 参考：
>
> - [GitHub 项目，vllm-project/vllm](https://github.com/vllm-project/vllm)
> - [官网](https://vllm.ai/)

用于 LLM 的高吞吐量和内存高效的推理和服务引擎

支持的模型: https://docs.vllm.ai/en/stable/models/supported_models/

# 部署

部署什么？

- [vllm CLI](/docs/12.AI/AI%20MGMT/vLLM/vllm%20CLI.md)？
- vllm 第三方库供 Python 导入？

感觉都可以，也挺简单的。

vllm 可以直接用容器镜像拉起来

Python 库可以用 uv 装

主要是后面如何用 vllm 使用模型进行推理。以及国产化，尤其是 [Ansible Plugins](/docs/9.运维/Ansible/Ansible%20Plugins/Ansible%20Plugins.md) 比较麻烦

## 容器化部署

TODO

# 使用 vLLM

使用 vLLM 意味着通过 vLLM 启动一个对外提供的推理服务

使用 `vllm serve` 命令，在 8000 端口上启动一个兼容 OpenAI API 的 WebAPI，提供了 Qwen3-0.6B 模型的推理服务：

```bash
vllm serve --model Qwen/Qwen3-0.6B
```

> [!Tip] 默认情况下，vllm 会从 [Hugging Face](/docs/12.AI/Hugging%20Face.md) 下载指定的模型。或者也可以使用 HF 镜像站下载模型后，手动指定本地模型路径

之后，我们可以通过 `http://localhost:8000/v1` 使用 OpenAI 兼容的 API 获取推理结果。

vllm 自身也提供了基本的聊天能力，使用 `vllm chat` 命令即可连接到 localhost:8000 并进入聊天模式，就像这样：

```bash
~]# vllm chat
Using model: Qwen/Qwen3-0.6B
Please enter a message for the chat model:
>
```

---

https://docs.vllm.ai/en/stable/models/supported_models/?h=vllm_use_modelscope#modelscope

如果在国内访问 [Hugging Face](/docs/12.AI/Hugging%20Face.md) Hub 困难，可以使用使用 https://www.modelscope.cn/ 代替。只需要设置环境变量即可：

```bash
export VLLM_USE_MODELSCOPE=True
```

之后，vllm 会到 ModelScope 默认的缓存目录中寻找 模型、数据集、etc.

# Benchmark

> 参考：
>
> - [官方文档，基准测试套件](https://docs.vllm.ai/en/stable/benchmarking/)

vLLM 提供全面的基准测试工具，用于性能测试和评估：

- **Benchmark CLI** # 用于交互式性能测试的 [vllm CLI](/docs/12.AI/AI%20MGMT/vLLM/vllm%20CLI.md) 的 bench 子命令 和 专用基准测试脚本。
  - 早期使用的是 benchmark_throughput.py 脚本，后来（TODO: 时间）将基准测试功能合并到 vllm CLI 中。
- **参数扫描** # 自动运行多个配置的 `vllm bench` ，有助于[优化和调优](https://docs.vllm.ai/en/stable/configuration/optimization/) 。
  - https://docs.vllm.ai/en/stable/benchmarking/sweeps/
- **性能仪表盘** # 自动化 CI，每次提交都会发布基准测试结果。
  - https://docs.vllm.ai/en/stable/benchmarking/dashboard/

# Plugins

> 参考：
>
> - [官方文档，设计 - 插件系统](https://docs.vllm.ai/en/stable/design/plugin_system/)

社区经常要求为 vLLM 添加自定义功能。为了方便用户实现这一目标，vLLM 内置了一个插件系统，允许用户在不修改 vLLM 代码库的情况下添加自定义功能。

- [Plugin Ascend](/docs/12.AI/AI%20MGMT/vLLM/Plugin%20Ascend.md) # 让 vLLM 可以运行在华为昇腾系列硬件上的插件
- etc.

