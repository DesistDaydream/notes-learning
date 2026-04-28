---
title: AI MGMT
linkTitle: AI MGMT
weight: 1
---

# 概述

> 参考：
>
> -

一套完整的可以在本地运行各种模型并对外提供 WebAPI 的程序

- llama.cpp
- [Ollama](#ollama) # 基于 llama.cpp 的封装
- vLLM
  - https://github.com/vllm-project/vllm
- LM Studio # 基于 llama.cpp 的闭源封装
- Jan
- koboldcpp
- ramalama
- etc.

# llama.cpp

> 参考：
>
> - [GitHub 项目，ggml-org/llama.cpp](https://github.com/ggml-org/llama.cpp)

LLaMA

llama.cpp 的主要目标是在本地和云端，以最小的设置和最先进的性能实现 LLM 推理。

llama.cpp 使用的模型格式为 **GPT-Generated Unified Format(GPT 生成的统一格式，简称 GGUF)**。这个 GGUF 文件中包含了所有必要的 元数据、分词器信息、模型权重。让 llama.cpp 使用一个单一的文件，即可进行推理。

> [!Tip] [Transformer](/docs/12.AI/Machine%20learning/Transformer/Transformer.md) 的模型通常包含 权重、配置、分词器配置、etc. 多个文件组成一个完整的 模型。

## 安装

从 [Release](https://github.com/ggml-org/llama.cpp/releases) 下载并使用即可，二进制文件直接用。

Windows 注意：有两个文件 [Windows x64 (CUDA 13)](https://github.com/ggml-org/llama.cpp/releases/download/b8840/llama-b8840-bin-win-cuda-13.1-x64.zip) 和 [CUDA 13.1 DLLs](https://github.com/ggml-org/llama.cpp/releases/download/b8840/cudart-llama-bin-win-cuda-13.1-x64.zip)，把 DLLs 文件解压到 llama.cpp 所在目录相同的目录下即可。

## 关联文件与配置

## 最佳实践

开始交互式对话

```powershell
.\llama-cli.exe -m D:\appdata\models\desistdayream.gguf
```

### 将模型包转换为 GGUF 文件

使用 llama.cpp 项目中的 convert_hf_to_gguf.py 文件将训练好的模型权重及相关文件，转换为单一的 GGUF 格式的文件

```powershell
python .\convert_hf_to_gguf.py `
  D:\appdata\models\desistdaydream\ `
  --outfile D:\appdata\models\desistdayream.gguf
```

# Ollama

> 参考：
>
> [GitHub 项目，ollama/ollama](https://github.com/ollama/ollama)
> Ollama 模型库: https://ollama.com/library

> [!Attention] 个人感受：该项目不符合开源精神，开源只是获取用户信任的工具，不是承诺。后来还推出来过闭源组件，又来又逐渐转向云上模型。并且，并不是真的运行本地模型。本地自己训练的模型用 ollama 去运行还是比较麻烦的，这么折腾下来不如直接用 llama.cpp

## Ollama 关联文件与配置

## Linux

****~/.ollama/** #

- **./models/** # 模型储存位置

## Windows

**%LOCALAPPDATA%/Programs/Ollama/** # Ollama 默认安装位置，包括 二进制程序、依赖库、CUDA、etc.

**%LOCALAPPDATA%/Ollama/** # 日志保存位置

- https://docs.ollama.com/troubleshooting 根据官方文档说明，Windows 日志要想落盘，需要使用 `ollama app.exe` 运行托盘程序才会生成。否则只能从标准输入输出查看日志，或者使用 `Start-Process -FilePath "ollama" -ArgumentList "serve" -RedirectStandardOutput "D:\log\ollama\stdout.log" -RedirectStandardError "D:\log\ollama\stderr.log" -NoNewWindow` 重定向日志
- https://github.com/ollama/ollama/pull/11552 有个 pr 想要在命令行增加查看日志的功能很久了还没合并

**%HOMEPATH%/.ollama/** # 模型与配置储存位置

- **./models/** # 模型储存位置。可以使用 `OLLAMA_MODELS` 环境变量定义新的储存位置

**%TEMP%/** # 临时可执行文件

## 生态支持

https://github.com/ollama/ollama/blob/main/README.md#community-integrations

Web 与 Desktop

https://github.com/ollama/ollama/blob/main/README.md#web--desktop

Chrome 插件

- https://github.com/n4ze3m/page-assist

# vLLM

> 参考：
>
> - [GitHub 项目，vllm-project/vllm](https://github.com/vllm-project/vllm)
> - [官网](https://vllm.ai/)

用于 LLM 的高吞吐量和内存高效的推理和服务引擎

# LM Studio

> 参考：
>
> - [GitHub 组织，LM Studio](https://github.com/lmstudio-ai)
> - [官网](https://lmstudio.ai/)

免费，不开源？

# Dify

> 参考：
>
> - [GitHub 项目，langgenius/dify](https://github.com/langgenius/dify)
> - [官网](https://dify.ai/)

Dify 是一个开源的 LLM 应用开发平台。Dify 直观的界面结合了 AI 工作流、RAG 管道、代理功能、模型管理、可观测性特征等，让您能够迅速从原型转向生产。
