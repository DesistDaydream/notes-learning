---
title: "Deepseek Harness"
created: "2026-09-26T09:42"
weight: 100
---

# 概述

> 参考：
>
> - [GitHub 项目，deepseek-ai/deepseek-harness](https://github.com/deepseek-ai/deepseek-harness)

DeepSeek Harness 基于 [Cordis](https://github.com/cordiverse/cordis) 插件系统构建。Harness 让 Agent 在真实场景中持续工作

- Cordis 内核 # Cordis 内核只负责插件的加载、卸载和依赖关系，不承载 Agent 的具体能力。
- 插件提供能力 # 模型、工具、技能、会话、沙箱、存储、循环、调度、UI 等所有 Agent 能力均由插件提供，并通过 Cordis 服务与事件彼此协作。
- 配置层自由组合 # 开发者无需改动源码，即可在配置层选择、替换或扩展任一能力。

> [!TODO]
> ~/.dsh/ 目录 和 工作区中，都可以放 AGENTS.md 文件。~/.dsh/AGENTS.md 文件可以被在任何工作区的任务读取。核心的个人偏好逻辑可以放在这里。
>
> 记忆可以使用 Go 语言开发的 [engram](https://github.com/Gentleman-Programming/engram)，任何记忆程序的使用方式可以放在 AGENTS.md 中，因为 AGENTS.md 是任务开始时第一批 Prompt。

# 关联文件与配置

**~/.dsh/** # Web 运行时数据保存目录。包括 会话、配置、etc.

- **./profiles/** # 各类应用程序的配置
    - **./web/** # Web 端的配置
    - **./desktop/** # 桌面端的配置（桌面端是 Electron 实现，但是本质也会启动 Web 端的后台，与 Web 端交互）。
- **`./sessions/${工作区ID}/session-${会话ID}/`** # 针对每个工作区的会话记录
- **./storages/** # 
- **settings.yaml** # 

# 记忆

https://deepseek-harness.github.io/deepseek-harness/guide/mcp-memory

使用 Engram，启动时使用 --patch 指定配置文件: `pnpm dsh web --patch "$PWD/apps/cli/config/examples/mcp-memory/engram.cordis.yml"`

- Engram 默认使用 `~/.engram/` 目录存储数据。

# Plugins

