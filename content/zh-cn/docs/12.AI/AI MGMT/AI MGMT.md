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

- [Ollama](#Ollama)
- vLLM
    - https://github.com/vllm-project/vllm
- etc.

# LM Studio

> 参考：
>
> - [GitHub 组织，LM Studio](https://github.com/lmstudio-ai)
> - [官网](https://lmstudio.ai/)

免费，不开源？

# vLLM

> 参考：
>
> - [GitHub 项目，vllm-project/vllm](https://github.com/vllm-project/vllm)
> - [官网](https://vllm.ai/)

用于 LLM 的高吞吐量和内存高效的推理和服务引擎

# Ollama

> 参考：
>
> [GitHub 项目，ollama/ollama](https://github.com/ollama/ollama)

[带你认识本地大语言模型框架Ollama(可直接上手)](https://wiki.eryajf.net/pages/97047e/)

Ollama 模型库: https://ollama.com/library

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

# Dify

> 参考：
>
> - [GitHub 项目，langgenius/dify](https://github.com/langgenius/dify)
> - [官网](https://dify.ai/)

Dify 是一个开源的 LLM 应用开发平台。Dify 直观的界面结合了 AI 工作流、RAG 管道、代理功能、模型管理、可观测性特征等，让您能够迅速从原型转向生产。
