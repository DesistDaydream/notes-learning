---
title: AI Projects
linkTitle: AI Projects
weight: 1
---

# 概述

> 参考：
>
> -

[AI 开源方案库-传递最新 AI 应用落地解决方案｜AIGCLINK](https://d.aigclink.ai/) # 收集各种 AI 相关的开源项目

https://github.com/Stability-AI/generative-models # 补帧，通过静态图片生成动图。

[labring/FastGPT](https://github.com/labring/FastGPT) # 一个快速使用 openai api 的平台。支持一键构建 AI 知识库，支持多用户、多模型管理。

https://github.com/binary-husky/gpt_academic # 为 GPT/GLM 等 LLM 大语言模型提供实用化交互接口，特别优化论文阅读/润色/写作体验，模块化设计，支持自定义快捷按钮&函数插件，支持Python 和 C++ 等项目剖析 & 自译解功能，PDF/LaTex 论文翻译&总结功能，支持并行问询多种 LLM 模型，支持 chatglm3 等本地模型。接入通义千问, deepseek coder, 讯飞星火, 文心一言, llama2, rwkv, claude2, moss 等。

**NovelAI** # https://novelai.net/ 写故事、作图

[即时通信](/docs/Utils/即时通信/即时通信.md) 中的 Chatbot 还有很多有趣的 AI 项目

[GitHub 组织，OpenRouterTeam](https://github.com/OpenRouterTeam) # 各种 AI 模型的聚合平台

- https://openrouter.ai/

MiniMax # AI 聚合平台

- [官网](https://www.minimaxi.com/)

# IDE 工具

**Cline** # IDE中的自主编码代理，能够在每一步的每一步中使用浏览器来创建/编辑文件，使用浏览器以及更多内容。实现了 [MCP](/docs/12.AI/MCP.md)

- https://github.com/cline/cline

# 用于解决复杂任务的 AI 项目/产品

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

## 各种 Agent

https://blog.x-agent.net/blog/xagent/

虽然开创性项目（e.g., [AutoGPT](https://github.com/Significant-Gravitas/AutoGPT), [BabyAGI](https://github.com/yoheinakajima/babyagi), [CAMEL](https://github.com/camel-ai/camel), [MetaGPT](https://github.com/geekan/MetaGPT), [AutoGen](https://github.com/microsoft/autogen), [DSPy](https://github.com/stanfordnlp/dspy), [AutoAgents](https://github.com/Link-AGI/AutoAgents), [OpenAgents](https://github.com/xlang-ai/OpenAgents), [Agents](https://github.com/aiwaves-cn/agents), [AgentVerse](https://github.com/OpenBMB/AgentVerse), [ChatDev](https://github.com/OpenBMB/ChatDev)）已经展示了这个方向的潜力，但完全自主的 AI Agent 之旅仍然面临着巨大的挑战。

[GitHub 项目，OpenBMB/XAgent](https://github.com/OpenBMB/XAgent)

- OpenBMB开源社区由清华大学自然语言处理实验室和[面壁智能](https://modelbest.cn/)共同支持发起

## 各种 Claw

[GitHub 项目，openclaw/openclaw](https://github.com/openclaw/openclaw) # AI Agent 的集合。TypeScript 语言编写。对接了多个聊天工具。

- 项目名字经过两次变化: Clawdbot, Moltbot
- https://claude.ai/share/eed2502c-d752-4b4f-95b5-11ff2c7674b9
- [Pi：OpenClaw 中的最小代理](https://lucumr.pocoo.org/2026/1/31/pi/) # Pi 的思想是 OpenClaw 的精髓
    - https://github.com/badlogic/pi-mono/

[GitHub 项目，qwibitai/nanoclaw](https://github.com/qwibitai/nanoclaw) # 精简的 OpenClaw

- https://mp.weixin.qq.com/s/xnyrHG3dRs2XSTDwAIa53w

[GitHub 项目，sipeed/picoclaw](https://github.com/sipeed/picoclaw) # 中国人主导开发，Go 语言编写。受 NanoClaw 启发的可以运行在树莓派中的 Claw。

[GitHub 项目，zeroclaw-labs/zeroclaw](https://github.com/zeroclaw-labs/zeroclaw) # 英国人主导，Rust 语言编写。本地优先、生产级安全、A2A 协议

[GitHub 项目，nextlevelbuilder/goclaw](https://github.com/nextlevelbuilder/goclaw) # 越南人主导开发，Go 语言编写。与 OpenClaw, ZeroClaw, PicoClaw, GoClaw 做了对比。

新闻

[GitHub 项目，gsscsd/big_model_radar](https://github.com/gsscsd/big_model_radar) # 每天早上 08:00 CST 自动运行的 GitHub Actions 工作流。追踪主流 AI CLI 工具的 GitHub 动态、OpenClaw 及其同赛道项目的生态活动、Anthropic 和 OpenAI 官网最新资讯，并每日监测 GitHub AI 热门仓库趋势，以中英双语每日简报的形式发布为 GitHub Issues 并提交为 Markdown 文件。每周和每月自动生成汇总报告。

https://github.com/duanyytop/agents-radar

- [OpenClaw](https://github.com/openclaw/openclaw)
- [NanoBot](https://github.com/HKUDS/nanobot)
- [Zeroclaw](https://github.com/zeroclaw-labs/zeroclaw)
- [PicoClaw](https://github.com/sipeed/picoclaw)
- [NanoClaw](https://github.com/qwibitai/nanoclaw)
- [IronClaw](https://github.com/nearai/ironclaw)
- [LobsterAI](https://github.com/netease-youdao/LobsterAI)
- [TinyClaw](https://github.com/TinyAGI/tinyclaw)
- [Moltis](https://github.com/moltis-org/moltis)
- [CoPaw](https://github.com/agentscope-ai/CoPaw)
- [ZeptoClaw](https://github.com/qhkm/zeptoclaw)
- [EasyClaw](https://github.com/gaoyangz77/easyclaw)

# ChatGPT

[ChatGPT](/docs/12.AI/AI%20Projects/ChatGPT.md)

# Claude

> 参考：
>
> - [公众号 - OSC 开源社区，Anthropic 推出 “更理性的 Claude”，正面硬刚 ChatGPT](https://mp.weixin.qq.com/s/7YJ7B6JTV7U1gXeLOiZsLw)

Anthropic 公司推出的

Claude 早期可以作为 [Slack](/docs/Utils/即时通信/Slack.md) 的应用被添加到 Workspace 中并无条件使用。

# DeepSeek

> 参考：
>
> - [GitHub 项目，deepseek-ai/DeepSeek-V3](https://github.com/deepseek-ai/DeepSeek-V3)

https://github.com/deepseek-ai/DeepSeek-V3/blob/main/inference/model.py 核心

模型解释: [B 站 - 秋葉aaaki，动态](https://www.bilibili.com/opus/1027408073324494885)

# Bard

https://bard.google.com/

google

# 文心一言

https://yiyan.baidu.com/

百度

# 通义千问

https://tongyi.aliyun.com/

阿里
