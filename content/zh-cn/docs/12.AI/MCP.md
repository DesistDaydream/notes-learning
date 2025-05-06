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

# 规范

> 参考：
>
> - [规范](https://modelcontextprotocol.io/specification/2025-03-26)

## 架构

https://modelcontextprotocol.io/specification/2025-03-26/architecture

![800](Excalidraw/mcp.excalidraw.md)

MCP 的工程化实现本质也是一个类似 C/S 的架构，但是在其中多了 Model(模型) 的参与。一套完整的 MCP 系统通常包含如下几类组件：

> Tips: 有些模型进行了 MCP 的训练，可以正常识别 MCP 的上下文内容，并且也可以返回 MCP 标准的上下文

- **Host(主持人)** # 充当内容的容器和协调器。创建和管理多个 MCP Client 实例。协调 AI/LLM 集成和采样。管理跨客户端的上下文聚合。
    - 类似 [Web](/docs/Web/Web.md) 中的 User Agent 概念。
- **MCP Client** # 由 Host 创建并维护。与 MCP Server 建立 1:1 的连接关系。
- **MCP Server** # 通过 MCP 原语公开 Resources, Tools, Prompts，提供后端程序的功能，用于与 MCP Client 通信。
    - 类似各种可以提供 API 的 Web 应用程序。只不过使用了 MCP 协议包装了 API 或其他各种功能
- **Model** # AI 模型。为问题提供 分析、总结 能力。

一套完整的 MCP 程序的运行过程通常如下所述：

用户向 Host 询问问题，Host 拉起多个 MCP Client，每个 MCP Client 注册一个 MCP Server，并获取其可以提供的能力

1. Host 带着 所有 MCP Server 的信息及其能力、问题内容，交给 Model
2. Model 分析问题后，决定是否要调用 MCP Server 以及如何调用，返回 MCP 上下文内容
3. Host 带着 MCP 上下文通过 MCP Client 向 MCP Server 发起请求获取数据
4. MCP 返回包含响应数据的 MCP 上下文
5. Host 拿着包含响应数据的 MCP 上下文交给 Model 分析各种结果
6. Model 以人类可读的方式返回结果
7. 若一个问题过于复杂，需要拆解问题后，重复上述 1 - 6 步

这里 https://github.com/CassInfra/KubeDoor/blob/1.3.0/src/kubedoor-mcp/kubedoor-mcp.py 有一个非常简单直观的 MCP Server 示例，其中利用 MCP 库的 FastMCP 实例化了一个 mcp。后面通过 @mcp.tool() 多次定义 Tools，每个 Tool 都由一个 API 提供能能力。Tools 中的 Args、Returns 本质就是结构化的 Prompts。当 Client 与 Server 建立连接时，Server 会将所有 Tools 响应给 Client。

# MCP Server

> 参考：
>
> - [规范 - 服务端功能](https://modelcontextprotocol.io/specification/2025-03-26/server)

MCP Server 的能力被抽象为如下几类（当 Client 注册 Server 时，需要将自身所具有的能力提供给 Client）：

- **Resources(资源)** # TODO: 好像一种直接返回具体数据的实体？比如一个数据库？
- **Tools(工具)** # 每个 Tool 都是一个可以实现的具体功能。可以简单理解为一个 API、一个具体的功能、etc.
- **Prompts(提示)** # TODO: 好像是让 Client 可以填充的提示词模板？

https://github.com/mark3labs/mcp-go

https://github.com/ThinkInAIXYZ/go-mcp

- https://mp.weixin.qq.com/s/LFIUVdTznkr7tWZ4_TnXGA

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

