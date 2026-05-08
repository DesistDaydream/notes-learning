---
title: vllm CLI
linkTitle: vllm CLI
created: 2026-05-08T22:25
weight: 11
---

# 概述

> 参考：
>
> - [官方文档，参考 - CLI](https://docs.vllm.ai/en/stable/cli/)

vllm 命令行工具用于运行和管理 [Model](/docs/12.AI/Machine%20learning/Model.md)(模型)

COMMAND

- **serve** # 启动 vLLM OpenAI 兼容 API 服务器。默认监听在 `0.0.0.0:8000`
- **chat** # 通过运行中的 WebAPI 开始聊天。默认连接到 `0.0.0.0:8000`。可以通过 --url 参数指定要连接的地址
- **complete** # 通过运行中的 WebAPI 进行文本补全。与对话调用 `/v1/chat/completions` 不同，这个调用 `/v1/completions`。
- **bench** # 基准测试
- **collect-env** # 收集并展示当前系统环境信息

# serve

https://docs.vllm.ai/en/stable/cli/serve/

下面三个命令效果一样，加载 Qwen3-0.6B 模型并在 8000 端口上启动 OpenAI API 兼容的 Web 服务：

- `vllm serve`
- `vllm serve Qwen/Qwen3-0.6B`
- `vllm serve --model Qwen/Qwen3-0.6B`

## OPTIONS

serve 命令的参数非常多，可以分为多个大类

- Frontend # OpenAI 兼容的 WebAPI 服务器的参数
- ModelConfig # 模型的配置
- LoadConfig # 加载模型权重的配置
- ParallelConfig # 分布式执行的配置。i.e. 使用多块 GPU/NPU/etc.

**Frontend**

- **--served-model-name**(STRING) # 用于在 API 中使用的模型名称。`默认值: 与 --model 参数的值相同`
- 工具调用相关参数
  - **--enable-auto-tool-choice** # 开启工具调用能力。不开启的话，进行工具调用时会报错: `tool choice requires --enable-auto-tool-choice and --tool-call-parser to be set`
  - **--tool-call-parser**(STRING) # 告诉 vLLM 用哪种格式解析模型输出的工具调用。
    - 可用的值参考: [官方文档，特性 - 工具调用](https://docs.vllm.ai/en/stable/features/tool_calling/)

**ModelConfig**

- **--max-model-len**(INT) # 最大上下文长度。`默认值: 使用模型中 config.json 文件的 max_position_embeddings 或 model_max_length 字段的值`。
- **--enforce-eager** # 设为 True 后，禁用 CUDA Graph（在 NPU 上是禁用对应的图优化），强制用 eager 模式执行。`默认值: False`
- **--dtype**(STRING) # 模型权重的数据类型。不指定的话 vLLM 默认用 bfloat16，昇腾不支持。`默认值: auto`
  - 可用的值有: `auto`, `bfloat16`, `float`, `float16`, `float32`, `half`

**LoadConfig**

- **--download-dir** # 用于下载和加载权重的目录。`默认值: Hugging Face 的默认缓存目录`
  - 若设置了 `VLLM_USE_MODELSCOPE=True`，则默认值为 ModelScope 的默认缓存目录

**ParallelConfig**

- **--tensor-parallel-size**(INT) # 并行使用的计算设备数量。比如同时使用 N 块 GPU/NPU。

# bench

https://docs.vllm.ai/en/stable/benchmarking/cli/