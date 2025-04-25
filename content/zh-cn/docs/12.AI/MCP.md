---
title: MCP
linkTitle: MCP
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 组织，modelcontextprotocol](https://github.com/modelcontextprotocol)
> - [官网](https://modelcontextprotocol.io/introduction)

**Model Context Protocol(模型上下文协议，简称 MCP)** 是一种开放协议，能够实现 LLM 应用程序与外部数据源和工具之间的无缝集成。无论是在构建 AI 驱动的 IDE、增强聊天界面，还是创建自定义 AI 工作流，MCP 都提供了一种标准化的方式，将 LLM 连接到所需的上下文。

https://github.com/ThinkInAIXYZ/go-mcp

- https://mp.weixin.qq.com/s/LFIUVdTznkr7tWZ4_TnXGA

MCP 的工程化实现本质也是一个类似 C/S 的架构，但是在其中多了一个 Model(模型的参与)。有些模型进行了 MCP 的训练，可以正常识别 MCP 的上下文内容，并且也返回 MCP 标准的上下文

1. MCP Client 带着工具列表，问题列表，向 Model 发起问题
2. Model 返回 MCP 上下文内容
3. MCP Client 带着 MCP 上下文向 MCP Server 发起请求获取数据
4. MCP 返回包含响应数据的 MCP 上下文
5. MCP Client 拿着包含响应数据的 MCP 上下文交给 Model 分析各种结果
6. Model 以人类可读的方式返回结果
7. 若一个问题过于复杂，需要拆解需求，则重复上述 1 - 6 步

![800](Excalidraw/mcp.excalidraw.md)

# 历史

## LangChain

> 参考：
> 
> - [GitHub 项目，langchain-ai/langchain](https://github.com/langchain-ai/langchain)
>   - 仓库最早在 hwchase17/langchain
> - [官网](https://langchain.com/)
> - [B 站，AI必学知识点！Langchain到底是什么？开源干货详细解析！赚钱机会和未来展望！](https://www.bilibili.com/video/BV1GL411e7K4)
> - [B 站，用自己的PDF文件定制Chatgpt！langchain代码实例详解！](https://www.bilibili.com/video/BV1xX4y1B7Vt)
> - https://github.com/liaokongVFX/LangChain-Chinese-Getting-Started-Guide
> - [公众号 - 阿里云开发者，LangChain: 大语言模型的新篇章](https://mp.weixin.qq.com/s/P94AvHvQcget9OqblrmD6g)

LangChain 是一个用于开发由语言模型驱动的应用程序的框架。我们相信，最强大和差异化的应用程序不仅会通过 API 调用语言模型，而且还会：

- 数据感知：将语言模型连接到其他数据源 
- Be agentic：允许语言模型与其环境交互
 
模块

- Models(模型)
- Prompts(提示词)
- Indexes
- Memory
- Chains
- Agents

## Function calling

