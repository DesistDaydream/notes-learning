---
title: OpenClaw
linkTitle: OpenClaw
created: 2026-03-25T21:24
weight: 100
---
不熟悉 TS，而且感觉 OpenClaw 无法成为个人助手的的最终形态，不如 Go 或 Rust 写的 Claw 可以更方便得部署在便携硬件上。另外本人也更喜欢 PicoClaw 的社区氛围

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
- [PicoClaw](#PicoClaw)
- [NanoClaw](https://github.com/qwibitai/nanoclaw)
- [IronClaw](https://github.com/nearai/ironclaw)
- [LobsterAI](https://github.com/netease-youdao/LobsterAI)
- [TinyClaw](https://github.com/TinyAGI/tinyclaw)
- [Moltis](https://github.com/moltis-org/moltis)
- [CoPaw](https://github.com/agentscope-ai/CoPaw)
- [ZeptoClaw](https://github.com/qhkm/zeptoclaw)
- [EasyClaw](https://github.com/gaoyangz77/easyclaw)

# PicoClaw

> 参考：
>
> - [GitHub 项目，sipeed/picoclaw](https://github.com/sipeed/picoclaw)

🦐 PicoClaw 是一个受 NanoBot 启发的超轻量级个人 AI 助手。它采用 Go 语言 从零重构，经历了一个"自举"过程——即由 AI Agent 自身驱动了整个架构迁移和代码优化。

⚡️ **极致轻量**：可在 **10 美元** 的硬件上运行，内存占用 **<10MB**。这意味着比 OpenClaw 节省 99% 的内存

# 运行逻辑

工具调用逻辑

当 LLM 返回的 tool_calls 字段不为空时， Agent 程序开始执行调用工具的逻辑；如果 LLM 没有返回 tool_calls 的话，则直接将 LLM 的 messages 信息返回给用户

`pkg/agent/loop.go`

```go
turnLoop:
    for 条件 {
        if len(response.ToolCalls) == 0 || gracefulTerminal {
            // 包装 LLM 的响应信息
            break
        }
        normalizedToolCalls := make([]providers.ToolCall, 0, len(response.ToolCalls))
    }
```

