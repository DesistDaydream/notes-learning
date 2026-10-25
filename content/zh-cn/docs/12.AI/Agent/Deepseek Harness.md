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
> 若为了更新，pull 新代码后执行 install 和 build 之后报错，可以尝试使用 `pnpm run clean` 清理后重新尝试。

# 安装与部署

这样装好像也挺简单，依赖都在 pnpm 的缓存里

```bash
pnpm add -g @deepseek-ai/dsh@0.2.1-rc.2
```

# 关联文件与配置

**$DSH_HOME** # Web 运行时数据保存目录。包括 会话、配置、etc. `默认值: ~/.dsh/`

- **./profiles/${APP}/** # 各类应用程序的配置。
    - **./cordis.patch.yml** # patch 配置文件。包括 DSH 配置、插件配置、etc.
- **`./sessions/${工作区ID}/session-${会话ID}/`** # 针对每个工作区的会话记录
- **./attachments/** # 会话中的附件保存目录。会话里的 图片、文件、etc. 都属于附件。附件名是内容的 SHA-256 指纹。
- **./storages/** #
- **settings.yaml** #

${APP} 可用的值有：

- **web** # Web 端的配置
- **desktop** # 桌面端的配置（桌面端是 Electron 实现，但是本质也会启动 Web 端的后台，与 Web 端交互）。

# 记忆

https://deepseek-harness.github.io/deepseek-harness/guide/mcp-memory

使用 Engram，启动时使用 --patch 指定配置文件: `pnpm dsh web --patch "$PWD/apps/cli/config/examples/mcp-memory/engram.cordis.yml"`

- Engram 默认使用 `~/.engram/` 目录存储数据。

# Plugins

若在某台电脑上安装了插件，在另一台电脑同步了 .dsh 的话，则只有插件名称，没有能运行插件的文件。可以通过下面命令强制重新安装，就在新电脑上也有插件了

```bash
dsh plugin --profile web install --force
```

## 插件推荐

https://github.com/topics/dsh-plugin # 这里的信息太多了，各种没有强关联的独立项目也加上了  dsh-plugin 这个 topics，明显是蹭热度引流的。

2026-09-26 cuitianyi 在 X 上发了[关于插件的内容](https://x.com/tianyi/status/2103534313463783831)，他惦记后面时不时个人推荐一些好用的 DSH 插件。

[GitHub 组织，omdsh-dev](https://github.com/omdsh-dev) # 围绕 DeepSeek Harness，分享社区作品、交流使用经验

https://github.com/dshworks/awesome-dsh-plugins

- https://awesome-dsh-plugin.com/

https://dshmarket.com/ # 插件市场

https://dshfind.com/ # 插件市场

---

https://github.com/deepseek-ai/dsh-libreoffice-kit # Office 文件处理。在 Node.js 中实现字体友好的 Office 转换与渲染。

https://github.com/oh-my-dsh/dsh-plugin-upgrade-skill # 帮助 DSH 插件作者升级自己插件、适配 DSH 新版本 的 skill

 - cuitianyi 推荐 https://x.com/tianyi/status/2107030904233460208

**界面定制**

- https://github.com/omdsh-dev/DSH-better-sidebar # 侧边栏增强，随意拖动，还有底边栏。

**观测系统**

- https://github.com/bowenliang123/dsh-context # 用量观测，上下文洞察。

**记忆系统**

- [Mnemon](/docs/12.AI/Prompt/Prompt.md#Mnemon)

**浏览器控制**

|                | BrowserSkill                                                    | dsh-builtin-browser                                           | dsh-browser-cdp                                                                   | dsh-browser                                                       |
| -------------- | --------------------------------------------------------------- | ------------------------------------------------------------- | --------------------------------------------------------------------------------- | ----------------------------------------------------------------- |
| 共用登录态          | 天然共享                                                            | 独立 profile，需各站登一次                                             | 独立 profile，可导入 cookie                                                             | 天然共享                                                              |
| 控制方式           | CDP + 插件连接 debugger                                             | 自起浏览器 + CDP                                                   | CDP + 自管浏览器<br>（内置运行时，可附着已有调试端点）                                                  | content script 注入<br>合成 DOM 事件                                    |
| 事件 `isTrusted` | **true**                                                        | **true**                                                      | **true**（CDP 回传）                                                                  | **false**                                                         |
| 环境要求           | dsh 插件 + 浏览器扩展 + bsk CLI                                        | dsh 插件 + Electron                                             | dsh 插件 + Chromium                                                                 | dsh 插件 + 浏览器扩展                                                    |
| 浏览器归属          | 用我的浏览器                                                          | 自己启动浏览器                                                       | 自己启动浏览器 or 用我的浏览器                                                                 | 自己启动浏览器                                                           |
| 可见痕迹           | 顶部黄条 + 页面控制浮层                                                   | 无黄条；有页内虚拟光标                                                   | 无黄条；DSH 侧有实时观察窗                                                                   | 无黄条，无浮层                                                           |
| 能力上限           | 读、操作、观测、调试                                                      | 读、操作、观测（部分）                                                   | 读、操作、观测（部分）                                                                       | 读、操作                                                              |
| 反爬角度           | 事件真实，但 CDP 可被探测                                                 | 事件真实、无扩展痕迹，但行为节奏仍可判                                           | 同左，另可导入 cookie 后走独立 profile                                                       | 最容易被识别（`isTrusted=false`）                                         |
| 维护             | 积极                                                              | 活跃（0.4.14，10-19）                                              | 活跃（1.6.1，10-24）                                                                   | 三周未提交                                                             |
| 链接             | [Tencent/BrowserSkill](https://github.com/Tencent/BrowserSkill) | [wqty123/dsh-browser](https://github.com/wqty123/dsh-browser) | [drscrewdriver/dsh-browser-cdp](https://github.com/drscrewdriver/dsh-browser-cdp) | [omdsh-dev/dsh-browser](https://github.com/omdsh-dev/dsh-browser) |
| 备注             |                                                                 |                                                               | 依赖 [ego](https://github.com/citrolabs/ego-lite) 控制 Chromium                       |                                                                   |

**页面标注**

这类插件可以快捷得标注页面中的元素，让 AI 更直接得理解想要描述得是页面哪部分的内容

https://github.com/alaliqing/dsh-annotate
