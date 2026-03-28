---
title: OpenClaw
linkTitle: OpenClaw
created: 2026-03-25T21:24
weight: 100
---

> [!Tip]
> 个人不熟悉 TS，而且感觉 OpenClaw 无法成为个人助手的的最终形态，不如 Go 或 Rust 写的 Claw 可以更方便得部署在便携硬件上。
>
> 另外本人也更喜欢 PicoClaw 的社区氛围，而且 PicoClaw 早早就支持了 MCP，在我提了 Issue 后，很快就将 MCP 返回的媒体信息进行了统一处理发送给 IM。

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

> [!quote] 在我总结完本篇笔记后，我换了一种提问方式再次尝试，Claude 给出了比较令人满意的结果: https://claude.ai/chat/bf5d62f7-9ddd-435a-9966-a8b644d1fa6c

- **Channels** # 消息的 入口 与 出口
- **Providers** # 消息处理模型

PicoClaw 使用 gateway 子命令[启动](#启动)主程序。

## 启动

pkg/gateway/gateway.go

**一.** `Run()` # PicoClaw 的 Gateway 模式启动

**二.** `createStartupProvider()` # 创建并启动 [Providers](#Providers)

**三.** `bus.NewMessageBus()` # 实例化消息总线

**四.** `agent.NewAgentLoop()` # 实例化 AgentLoop。AgentLoop 是程序的核心处理功能，接受消息 - 推理 - 响应消息。

- `NewAgentRegistry()` # 设置 AgentRegistry（用于管理多个 Agent 实例），以便后续消息处理时，可以将消息路由给期望的 Agent。
- 注册所有工具\初始化 HookManager、etc.

**五.** `setupAndStartServices()` # 设置并启动所有服务。Cron, Heartbeat, Bus, Media, Channels

**六.** `agent.AgentLoop.Run()` # <font color="#ff0000">**程序运行的核心逻辑**</font>。使用 Goroutine 让 AgentLoop 运行起来。进入 [消息处理](#消息处理)

## 消息处理

`pkg/agent/loop.go`

**一.** `case msg, ok := <-al.bus.InboundChan()` # 消息入口。`InboundChan()` 接收到来自 [Channels](#Channels) 的消息后开始走这部分流程。

- `InboundChan()` 返回的是 `bus.InboundMessage` 类型的 Go 通道，名为 inbound。Channels 将消息推送到 inbound 后，由消息处理模块进行消费处理。

**二.** `AgentLoop.resolveMessageRoute()` # 路由消息。根据 Channel/Peer 决定使用哪个 `AgentInstance` 处理消息。

**三.** `AgentLoop.processMessage()` # 处理消息，获取 LLM 的推理结果。

**四.** `AgentLoop.handleCommand()` # 处理指令类消息。e.g. `/clear`, `/start`, etc.

**五.** `AgentLoop.runAgentLoop()` # 处理非指令类消息。这是 Agent 的顶层外壳，负责启动 **Turn**

- Turn 相对 AgentLoop 来说，就是一个

**六.** `AgentLoop.runTurn()` # <font color="#ff0000">Turn 的完整生命周期管理</font>，Turn 中包含 **turnLoop** 循环。<font color="#ff0000">**Turn 是 OpenClaw 处理消息获取响应的核心逻辑**</font>。接收消息 -> ReAct -> 给出回复，这个过程就叫一个 Turn。

- **组装 ModelService 可用的上下文**
    - `turnState.agnet.ContextBuilder.BuildMessages()` # 组装上下文，要发送给 LLM 的 [Conversation](/docs/12.AI/自然语言处理/自然语言处理.md#Conversation)（system role 与 user role）
        - 其中还包括之前加载的历史消息。
        - 加载 Skills
    - `turnState.agent.Tools.ToProviderDefs()` # 组装上下文，要发送给 LLM 的 [Tools](/docs/12.AI/自然语言处理/自然语言处理.md#Tools)
    - TODO: 非结构化的媒体文件怎么处理？
> [!Tip] pkg/agent/context.go 中 `ContextBuilder.getIdentity()` 是代码中内置的 system prompt，设定了几个非常核心且简单的逻辑，包括工作路径，有限使用工具等。
>
> getIdentity 之后还会从其他几个文件中读取一些信息填充到 messages。
- **turnLoop** # <font color="#ff0000">开始核心循环逻辑</font>。其中包括 调用 [Providers](#Providers) 获取推理结果、调用工具、etc. 。详见 [turnLoop](#turnLoop)

**七.** `al.bus.PublishOutbound()` # 消息出口。将推理结果发送给 [Channels](#Channels)

## turnLoop

`pkg/agent/loop.go`

**这是 PicoClaw Agent 能力的核心循环**。循环条件有：

- `ts.currentIteration() < ts.agent.MaxIterations` # 是否已达到设置的最大迭代次数
- `len(pendingMessages) > 0` # 等待处理的消息是否大于 0
- TODO: 其他

在这个循环中，会执行如下动作：

**一.** Steering，处理用户在 Agent 中途再次发送的消息（e.g. 中断任务、etc.）

**二.** `callLLM func()` # 调用 LLM

- `providers.LLMProvider.Chat()` # LLMProvider 接口的 Chat 方法。调用 [Providers](#Providers) <font color="#ff0000">获取模型服务的推理结果</font>

**三.** 若 LLM 返回了 `tool_calls`，则开始 [工具调用](#工具调用)

**四.** `maybeSummarize()` # 若满足条件，则压缩历史上下文，防止上下文超长导致 LLM 无法进行推理

**五.** 循环完成后，返回 `finalContent`

### 工具调用 

当 LLM 返回的 tool_calls 字段不为空时， Agent 程序开始执行调用工具的逻辑

```go
turnLoop:
    for 条件 {
        if len(response.ToolCalls) == 0 || gracefulTerminal {
            // 包装 LLM 的响应信息后直接跳出循环，不再执行后面的工具调用逻辑。
            break
        }
        normalizedToolCalls := make([]providers.ToolCall, 0, len(response.ToolCalls))
        toolResult := ts.agent.Tools.ExecuteWithContext()
    }
```

工具的执行结果将会添加到上下文中交给模型

# 组件

### Channels

Channels 是 发送/接收 消息的端点，通常是各类 [即时通信](/docs/Utils/即时通信/即时通信.md)。

```
pkg/channels/
├── base.go              # BaseChannel 共享抽象层
├── interfaces.go        # 可选能力接口（TypingCapable, MessageEditor, ReactionCapable, PlaceholderCapable, PlaceholderRecorder）
├── media.go             # MediaSender 可选接口
├── webhook.go           # WebhookHandler, HealthChecker 可选接口
├── errors.go            # 错误哨兵值（ErrNotRunning, ErrRateLimit, ErrTemporary, ErrSendFailed）
├── errutil.go           # 错误分类帮助函数
├── registry.go          # 工厂注册表（RegisterFactory / getFactory）
├── manager.go           # 统一编排：Worker 队列、速率限制、重试、Typing/Placeholder、共享 HTTP
├── split.go             # 长消息智能分割（保留代码块完整性）
├── telegram/            # 每个 channel 独立子包
│   ├── init.go          # 工厂注册
│   ├── telegram.go      # 实现
│   └── telegram_commands.go
├── discord/
│   ├── init.go
│   └── discord.go
├── slack/ line/ onebot/ dingtalk/ feishu/ wecom/ qq/ whatsapp/ whatsapp_native/ maixcam/ pico/
│   └── ...

pkg/bus/
├── bus.go               # MessageBus（缓冲区 64，安全关闭+排水）
├── types.go             # 结构化消息类型（Peer, SenderInfo, MediaPart, InboundMessage, OutboundMessage, OutboundMediaMessage）

pkg/media/
├── store.go             # MediaStore 接口 + FileMediaStore 实现（两阶段释放，TTL 清理）

pkg/identity/
├── identity.go          # 统一用户身份：规范 "platform:id" 格式 + 向后兼容匹配
```

- bus 信息总线
- media 媒体信息处理
- TODO

```text
┌─────────────┐      InboundMessage       ┌───────────┐      LLM + Tools      ┌────────────┐
│  Telegram   │──┐                        │           │                       │            │
│  Discord    │──┤   PublishInbound()     │           │   PublishOutbound()   │            │
│  Slack      │──┼──────────────────────▶ │ MessageBus│ ◀───────────────────  │ AgentLoop  │
│  LINE       │──┤   (buffered chan, 64)  │           │   (buffered chan, 64) │            │
│  ...        │──┘                        │           │                       │            │
└─────────────┘                           └─────┬─────┘                       └────────────┘
                                                │
                            SubscribeOutbound() │  SubscribeOutboundMedia()
                                                ▼
                                    ┌────────────────────┐
                                    │   Manager          │
                                    │   ├── dispatchOutbound()    路由到 Worker 队列
                                    │   ├── dispatchOutboundMedia()
                                    │   ├── runWorker()           消息分割 + sendWithRetry()
                                    │   ├── runMediaWorker()      sendMediaWithRetry()
                                    │   ├── preSend()             停止 Typing + 撤销 Reaction + 编辑 Placeholder
                                    │   └── runTTLJanitor()       清理过期 Typing/Placeholder
                                    └────────┬───────────┘
                                             │
                                   channel.Send() / SendMedia()
                                             │
                                             ▼
                                    ┌────────────────┐
                                    │ 各平台 API/SDK  │
                                    └────────────────┘
```

### Providers

Providers 是各种 [Model service](/docs/12.AI/Model%20service.md) 的提供者。

核心接口:

pkg/providers/types.go

```go
type LLMProvider interface{
    Chat(参数略) (LLMResponse, error)
    GetDefaultModel() string
}
```

