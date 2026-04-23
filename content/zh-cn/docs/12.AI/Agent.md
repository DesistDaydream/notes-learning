---
title: Agent
linkTitle: Agent
created: 2026-03-14T18:30
weight: 60
---
**前言**

问题：当 LLM 想要调用工具时，有两种设计方式

一

- 把 MCP Server 中所有工具的信息（名字、描述、参数 Schema、etc.）都要交给 LLM
- LLM 直接把要使用的工具和参数一起返回

二

- 只把 MCP Server 中的所有工具的 名字 和 描述交给 LLM
- LLM 返回要使用的工具名称
- Client 获取指定工具的所有信息交给 LLM
- LLM 返回要使用工具和参数

不管怎么设计，每次交互都要带着前文一起，这 Agent 能力对上下文的占用，也太大了点吧？

在这个例子中，当 AI 设计出来一程序时，如果该程序消耗 Token 过多。在人类不知道原理的情况下，只会跟 AI 说：我的 Token 消耗太多了，帮我减少 Token 消耗。但是 AI 也没法知道什么算多什么算少，应该少到多少算**合理**。AI 能合理优化代码的前提时人类告诉 AI 如何优化。因为如果 AI 本身就知道，也就不会出现需要优化的问题，也就不会写出来更高 Token 消耗的代码

# 概述

> 参考：
>
> - 

**AI Agent**

**Agent Runtime** # 类似 MCP 的 Host，Agent 与模型交互的程序。实现了 ReAct 功能

**ReAct** # Reasoning and Acting。Agent Runtime 的核心逻辑：循环。

- https://zhuanlan.zhihu.com/p/1978741652337205325
- https://www.ibm.com/think/topics/react-agent

# 资料

[AI Agent 本地LLM推理设备部署指南](https://llmdev.guide/) # 用真实数据说话，拒绝虚标 — 社区驱动的大模型推理设备性能数据库。PicoClaw 主家出的。

[B 站 - 老戴Donald，【老戴】用了几天 Clawdbot，我最担心的事开始发生了](https://www.bilibili.com/video/BV1MnFNz5EDL?spm_id_from=333.1245.0.0)

[B 站 - 飞天闪客，【闪客】20 行代码彻底搞懂小龙虾！男女老少都看得懂哟~](https://www.bilibili.com/video/BV19hwTzwETF)

# Agent 架构

![700](Excalidraw/AI/agent-arch.excalidraw.md)

# Agent 流程

![800](Excalidraw/AI/agent-flow.excalidraw.md)

# Pi-mono

> 参考：
>
> - [GitHub 项目，badlogic/pi-mono](https://github.com/badlogic/pi-mono)
> - [公众号，从 pi-mono 到 OpenClaw：源码拆解，21 万 Star 背后的 Agent 工程减法](https://mp.weixin.qq.com/s/qr0Ch8ii79KymUiK6kV_lA)
> - [公众号，OpenClaw背后的英雄Pi-mono](https://mp.weixin.qq.com/s/XRZrOvapXWiZneDtipfHTQ)

pi-mono 只保留四个工具：

- read # 读取文件内容，支持文本和图片，可指定行范围
- write # 创建新文件或完全重写，自动创建目录
- edit # 精确替换文本，oldText 必须完全匹配
- bash # 执行命令，返回 stdout 和 stderr，可设置超时

Zechner 的逻辑很直接：编程的本质就是读代码、写代码、改代码、跑代码。这四个工具组合起来能覆盖大部分编程场景。

# Agent 痛点

故事来源：B 站，飞天闪客 的某期 Agent 视频

现阶段的 Agent 就像这样：老板问我们今天几号

我们分析：这个任务需要看手机。开始执行任务

掏出手机；发现手机没电了，需要找充电器；没找到充电器，需要去楼下超市买个充电器；到超市门口发现超市关门了，但是发现可以去另外一个超市，但是需要打车；下车之后发现手机没电没法付款，又没有现金，所以需要去银行取钱；到银行取钱发现钱不够，那就先把手机抵押给司机。

此时到了新超市，用奇妙的方式买到了充电器，发现手机没了；所以需要买个手机，但是钱不够，需要去银行贷款。

最终，花了一整天的时间，丢掉了原来的手机，欠了一屁股债买了个新手机，最终回复给我们，今天是 X 年 X 月 X 日。

# 知乎-北方的狼 智能体

**本书目录**

本书共分六大部分，20个章节，沿着 **“认知 -> 大脑 -> 手脚 -> 骨架 -> 社会”** 的逻辑脉络，为您完整复原一个智能体的诞生过程。

**第一部分：定义 —— 什么是2025年的智能体？**

[AI智能体，第1章 从 Copilot 到 Autopilot](https://zhuanlan.zhihu.com/p/1978738837611095231)

[AI智能体，第2章 解剖学：Agent 的四个象限](https://zhuanlan.zhihu.com/p/1978739569391326726)

**第二部分：大脑 —— 推理与规划的原子能力**

[AI智能体，第3章 结构化输出：JSON Mode 与 Pydantic](https://zhuanlan.zhihu.com/p/1978740200084619518)

[AI智能体，第4章 思维链的进化：从 CoT 到 Reasoning Models](https://zhuanlan.zhihu.com/p/1978740962571356067)

[AI智能体，第5章 行动的循环：手动实现 ReAct Loop](https://zhuanlan.zhihu.com/p/1978741652337205325)

[AI智能体，第6章 元认知：Self-Reflection（自省机制）](https://zhuanlan.zhihu.com/p/1978742254500861944)

[AI智能体，第7章 提示词工程的终结：DSPy 自动优化](https://zhuanlan.zhihu.com/p/1978742952189776400)

**第三部分：手脚 —— 标准化工具与MCP协议**

[AI智能体，第8章 工具调用的本质：Function Calling](https://zhuanlan.zhihu.com/p/1979215299539640724)

[AI智能体，第9章 统一连接标准：MCP (Model Context Protocol)](https://zhuanlan.zhihu.com/p/1979216040341832707)

[AI智能体，第10章 实战 MCP：构建你的第一个 Server](https://zhuanlan.zhihu.com/p/1979216479242166524)

[AI智能体，第11章 视觉行动：Computer Use 与 GUI Agent](https://zhuanlan.zhihu.com/p/1979217119179735151)

**第四部分：骨架 —— 状态管理与图编排**

[AI智能体，第12章 从链到图：LangGraph 核心概念](https://zhuanlan.zhihu.com/p/1979217687793120310)

[AI智能体，第13章 状态管理：State Schema 与 Checkpoint](https://zhuanlan.zhihu.com/p/1979218315114206098)

[AI智能体，第14章 人在回路：Human-in-the-loop](https://zhuanlan.zhihu.com/p/1979218746024408085)

[AI智能体，第15章 记忆系统：Mem0 与 GraphRAG](https://zhuanlan.zhihu.com/p/1979219186959016315)

**第五部分：社会 —— 多智能体协作 (Multi-Agent)**

**第六部分：结语 —— 迈向工业级**

——完——

[@北方的郎](https://www.zhihu.com/people/7af62e4119791a452e88718cb5ccc0be) · 专注模型与代码
